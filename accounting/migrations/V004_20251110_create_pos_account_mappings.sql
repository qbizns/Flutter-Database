-- ============================================================================
-- Migration: V004 - Create POS Account Mappings
-- Description: Map POS entities (products, payment methods, discounts) to GL accounts
--              for automated posting from Go layer
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- POS ACCOUNT MAPPINGS TABLE
-- Description: Generic mapping table for POS → GL account distribution
-- ============================================================================

CREATE TABLE IF NOT EXISTS pos_account_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Source Entity (polymorphic reference)
    source_type VARCHAR(50) NOT NULL CHECK (source_type IN (
        'product',           -- Individual product
        'category',          -- Product category (fallback)
        'payment_method',    -- Payment method (cash, card, etc.)
        'sales_channel',     -- Sales channel
        'discount',          -- Discount type
        'rounding',          -- Rounding adjustment
        'tax',               -- Tax type
        'service_charge',    -- Service charges
        'shipping',          -- Shipping/delivery
        'gift_card',         -- Gift card
        'store_credit',      -- Store credit
        'loyalty_redemption',-- Loyalty point redemption
        'default'            -- Organization-level defaults
    )),
    source_id UUID,  -- NULL for 'default' type or when using source_code
    source_code VARCHAR(100), -- Alternative to source_id (e.g., payment method code)

    -- Mapping Purpose
    purpose VARCHAR(50) NOT NULL CHECK (purpose IN (
        'revenue',           -- Sales revenue account
        'cogs',              -- Cost of goods sold
        'inventory',         -- Inventory asset account
        'expense',           -- Expense account
        'liability',         -- Liability account (gift cards, store credit)
        'asset',             -- Asset account (cash, bank)
        'discount_expense',  -- Discount expense
        'discount_contra',   -- Discount contra-revenue
        'tax_liability',     -- Tax payable/receivable
        'rounding',          -- Rounding differences
        'clearing'           -- Clearing/suspense account
    )),

    -- GL Account Reference
    account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,

    -- Configuration
    is_default BOOLEAN DEFAULT FALSE,  -- Default mapping for this type
    is_active BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,  -- Higher priority = used first

    -- Conditional Logic (optional, for complex rules)
    conditions JSONB DEFAULT '{}',  -- e.g., {"amount_threshold": 1000, "customer_type": "wholesale"}

    -- Effective Dates
    effective_from DATE,
    effective_to DATE,

    -- Notes
    description TEXT,
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT valid_source_reference CHECK (
        (source_type = 'default' AND source_id IS NULL AND source_code IS NULL) OR
        (source_type != 'default' AND (source_id IS NOT NULL OR source_code IS NOT NULL))
    ),
    CONSTRAINT valid_effective_dates CHECK (
        effective_from IS NULL OR effective_to IS NULL OR effective_from <= effective_to
    )
);

-- Indexes
CREATE INDEX idx_pos_account_mappings_organization_id ON pos_account_mappings(organization_id);
CREATE INDEX idx_pos_account_mappings_source ON pos_account_mappings(source_type, source_id) WHERE source_id IS NOT NULL;
CREATE INDEX idx_pos_account_mappings_source_code ON pos_account_mappings(source_type, source_code) WHERE source_code IS NOT NULL;
CREATE INDEX idx_pos_account_mappings_purpose ON pos_account_mappings(purpose);
CREATE INDEX idx_pos_account_mappings_account_id ON pos_account_mappings(account_id);
CREATE INDEX idx_pos_account_mappings_is_default ON pos_account_mappings(source_type, purpose, is_default) WHERE is_default = TRUE;
CREATE INDEX idx_pos_account_mappings_active ON pos_account_mappings(is_active, effective_from, effective_to) WHERE is_active = TRUE;

-- Trigger
CREATE TRIGGER update_pos_account_mappings_updated_at
    BEFORE UPDATE ON pos_account_mappings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE pos_account_mappings IS 'Maps POS entities to GL accounts for automated journal entry posting';
COMMENT ON COLUMN pos_account_mappings.source_type IS 'Type of POS entity being mapped (product, payment_method, discount, etc.)';
COMMENT ON COLUMN pos_account_mappings.source_id IS 'UUID of the source entity (products.id, categories.id, etc.) - NULL for defaults';
COMMENT ON COLUMN pos_account_mappings.source_code IS 'Alternative string identifier (e.g., payment method code like "CASH", "CARD")';
COMMENT ON COLUMN pos_account_mappings.purpose IS 'What this account is used for (revenue, cogs, inventory, liability, etc.)';
COMMENT ON COLUMN pos_account_mappings.account_id IS 'Target GL account from chart_of_accounts';
COMMENT ON COLUMN pos_account_mappings.is_default IS 'Default mapping when no specific mapping exists';
COMMENT ON COLUMN pos_account_mappings.priority IS 'Priority when multiple mappings match (higher = preferred)';
COMMENT ON COLUMN pos_account_mappings.conditions IS 'Optional JSONB for complex conditional logic';

-- ============================================================================
-- HELPER FUNCTION: Get GL Account for POS Entity
-- ============================================================================

CREATE OR REPLACE FUNCTION get_pos_gl_account(
    p_organization_id UUID,
    p_source_type VARCHAR,
    p_source_id UUID,
    p_source_code VARCHAR,
    p_purpose VARCHAR,
    p_date DATE DEFAULT CURRENT_DATE
) RETURNS UUID AS $$
DECLARE
    v_account_id UUID;
BEGIN
    -- Try to find specific mapping first (by ID or code)
    SELECT account_id INTO v_account_id
    FROM pos_account_mappings
    WHERE organization_id = p_organization_id
      AND source_type = p_source_type
      AND (
          (source_id = p_source_id AND p_source_id IS NOT NULL) OR
          (source_code = p_source_code AND p_source_code IS NOT NULL)
      )
      AND purpose = p_purpose
      AND is_active = TRUE
      AND (effective_from IS NULL OR effective_from <= p_date)
      AND (effective_to IS NULL OR effective_to >= p_date)
      AND deleted_at IS NULL
    ORDER BY priority DESC, created_at DESC
    LIMIT 1;

    -- If not found, try default mapping for this type + purpose
    IF v_account_id IS NULL THEN
        SELECT account_id INTO v_account_id
        FROM pos_account_mappings
        WHERE organization_id = p_organization_id
          AND source_type = p_source_type
          AND is_default = TRUE
          AND purpose = p_purpose
          AND is_active = TRUE
          AND (effective_from IS NULL OR effective_from <= p_date)
          AND (effective_to IS NULL OR effective_to >= p_date)
          AND deleted_at IS NULL
        ORDER BY priority DESC, created_at DESC
        LIMIT 1;
    END IF;

    -- If still not found, try organization-wide default
    IF v_account_id IS NULL THEN
        SELECT account_id INTO v_account_id
        FROM pos_account_mappings
        WHERE organization_id = p_organization_id
          AND source_type = 'default'
          AND purpose = p_purpose
          AND is_active = TRUE
          AND (effective_from IS NULL OR effective_from <= p_date)
          AND (effective_to IS NULL OR effective_to >= p_date)
          AND deleted_at IS NULL
        ORDER BY priority DESC, created_at DESC
        LIMIT 1;
    END IF;

    RETURN v_account_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_pos_gl_account IS 'Resolves GL account for a POS entity with fallback logic: specific → type default → org default';

-- ============================================================================
-- HELPER VIEW: Account Mappings with Account Details
-- ============================================================================

CREATE OR REPLACE VIEW view_pos_account_mappings AS
SELECT
    pam.id,
    pam.organization_id,
    o.name as organization_name,
    pam.source_type,
    pam.source_id,
    pam.source_code,
    pam.purpose,
    pam.account_id,
    coa.account_code,
    coa.account_name,
    coa.account_number,
    at.type_name as account_type,
    pam.is_default,
    pam.is_active,
    pam.priority,
    pam.effective_from,
    pam.effective_to,
    pam.description,
    pam.created_at,
    pam.updated_at
FROM pos_account_mappings pam
INNER JOIN organizations o ON pam.organization_id = o.id
INNER JOIN chart_of_accounts coa ON pam.account_id = coa.id
INNER JOIN account_types at ON coa.account_type_id = at.id
WHERE pam.deleted_at IS NULL;

COMMENT ON VIEW view_pos_account_mappings IS 'POS account mappings with account details for easy reporting';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

ALTER TABLE pos_account_mappings ENABLE ROW LEVEL SECURITY;

CREATE POLICY pos_account_mappings_tenant_isolation ON pos_account_mappings
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY pos_account_mappings_tenant_isolation ON pos_account_mappings
    IS 'Ensure users can only access account mappings for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V004 completed successfully!';
    RAISE NOTICE 'POS Account Mappings created.';
    RAISE NOTICE ' - pos_account_mappings table';
    RAISE NOTICE ' - get_pos_gl_account() function';
    RAISE NOTICE ' - view_pos_account_mappings view';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

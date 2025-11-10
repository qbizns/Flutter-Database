-- ============================================================================
-- Migration: V008 - Create POS Tax Mappings
-- Description: Bridge POS tax usage with accounting tax system
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- POS TAX MAPPINGS TABLE
-- Description: Map POS tax codes/categories to accounting tax definitions
-- ============================================================================

CREATE TABLE IF NOT EXISTS pos_tax_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- POS Tax Reference
    pos_tax_code VARCHAR(50),              -- e.g., 'VAT_15', 'SALES_TAX', 'EXEMPT'
    tax_category_code VARCHAR(50),         -- From e-invoicing: 'S', 'Z', 'E', 'O'
    pos_tax_rate NUMERIC(5, 2),            -- Tax rate used in POS (for matching)

    -- Accounting Tax Reference
    accounting_tax_id UUID REFERENCES taxes(id) ON DELETE RESTRICT,

    -- Default GL Accounts
    default_tax_account_id UUID REFERENCES chart_of_accounts(id),  -- Tax liability/receivable account
    default_tax_expense_account_id UUID REFERENCES chart_of_accounts(id),  -- Tax expense (for non-recoverable VAT)

    -- Configuration
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,

    -- Tax Behavior
    is_inclusive BOOLEAN DEFAULT FALSE,    -- Whether tax is included in price
    applies_to_sales BOOLEAN DEFAULT TRUE,
    applies_to_purchases BOOLEAN DEFAULT FALSE,

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
    CONSTRAINT valid_effective_dates CHECK (
        effective_from IS NULL OR effective_to IS NULL OR effective_from <= effective_to
    )
);

-- Indexes
CREATE INDEX idx_pos_tax_mappings_organization_id ON pos_tax_mappings(organization_id);
CREATE INDEX idx_pos_tax_mappings_pos_tax_code ON pos_tax_mappings(pos_tax_code) WHERE pos_tax_code IS NOT NULL;
CREATE INDEX idx_pos_tax_mappings_tax_category_code ON pos_tax_mappings(tax_category_code) WHERE tax_category_code IS NOT NULL;
CREATE INDEX idx_pos_tax_mappings_accounting_tax_id ON pos_tax_mappings(accounting_tax_id);
CREATE INDEX idx_pos_tax_mappings_is_default ON pos_tax_mappings(organization_id, is_default) WHERE is_default = TRUE;
CREATE INDEX idx_pos_tax_mappings_active ON pos_tax_mappings(is_active, effective_from, effective_to) WHERE is_active = TRUE;

-- Trigger
CREATE TRIGGER update_pos_tax_mappings_updated_at
    BEFORE UPDATE ON pos_tax_mappings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE pos_tax_mappings IS 'Maps POS tax codes and categories to accounting tax definitions';
COMMENT ON COLUMN pos_tax_mappings.pos_tax_code IS 'POS-specific tax code (VAT_15, SALES_TAX, etc.)';
COMMENT ON COLUMN pos_tax_mappings.tax_category_code IS 'E-invoicing tax category (S=Standard, Z=Zero, E=Exempt, O=Out of scope)';
COMMENT ON COLUMN pos_tax_mappings.accounting_tax_id IS 'Reference to accounting.taxes table';
COMMENT ON COLUMN pos_tax_mappings.default_tax_account_id IS 'GL account for tax liability/receivable';

-- ============================================================================
-- HELPER FUNCTION: Get Accounting Tax for POS Tax
-- ============================================================================

CREATE OR REPLACE FUNCTION get_accounting_tax_for_pos(
    p_organization_id UUID,
    p_pos_tax_code VARCHAR DEFAULT NULL,
    p_tax_category_code VARCHAR DEFAULT NULL,
    p_pos_tax_rate NUMERIC DEFAULT NULL,
    p_date DATE DEFAULT CURRENT_DATE
) RETURNS UUID AS $$
DECLARE
    v_accounting_tax_id UUID;
BEGIN
    -- Try to match by POS tax code first
    IF p_pos_tax_code IS NOT NULL THEN
        SELECT accounting_tax_id INTO v_accounting_tax_id
        FROM pos_tax_mappings
        WHERE organization_id = p_organization_id
          AND pos_tax_code = p_pos_tax_code
          AND is_active = TRUE
          AND (effective_from IS NULL OR effective_from <= p_date)
          AND (effective_to IS NULL OR effective_to >= p_date)
          AND deleted_at IS NULL
        ORDER BY priority DESC, created_at DESC
        LIMIT 1;

        IF v_accounting_tax_id IS NOT NULL THEN
            RETURN v_accounting_tax_id;
        END IF;
    END IF;

    -- Try to match by tax category code
    IF p_tax_category_code IS NOT NULL THEN
        SELECT accounting_tax_id INTO v_accounting_tax_id
        FROM pos_tax_mappings
        WHERE organization_id = p_organization_id
          AND tax_category_code = p_tax_category_code
          AND is_active = TRUE
          AND (effective_from IS NULL OR effective_from <= p_date)
          AND (effective_to IS NULL OR effective_to >= p_date)
          AND deleted_at IS NULL
        ORDER BY priority DESC, created_at DESC
        LIMIT 1;

        IF v_accounting_tax_id IS NOT NULL THEN
            RETURN v_accounting_tax_id;
        END IF;
    END IF;

    -- Try to match by tax rate (less reliable but useful)
    IF p_pos_tax_rate IS NOT NULL THEN
        SELECT accounting_tax_id INTO v_accounting_tax_id
        FROM pos_tax_mappings
        WHERE organization_id = p_organization_id
          AND ABS(pos_tax_rate - p_pos_tax_rate) < 0.01
          AND is_active = TRUE
          AND (effective_from IS NULL OR effective_from <= p_date)
          AND (effective_to IS NULL OR effective_to >= p_date)
          AND deleted_at IS NULL
        ORDER BY priority DESC, created_at DESC
        LIMIT 1;

        IF v_accounting_tax_id IS NOT NULL THEN
            RETURN v_accounting_tax_id;
        END IF;
    END IF;

    -- Return default mapping if no specific match found
    SELECT accounting_tax_id INTO v_accounting_tax_id
    FROM pos_tax_mappings
    WHERE organization_id = p_organization_id
      AND is_default = TRUE
      AND is_active = TRUE
      AND (effective_from IS NULL OR effective_from <= p_date)
      AND (effective_to IS NULL OR effective_to >= p_date)
      AND deleted_at IS NULL
    ORDER BY priority DESC, created_at DESC
    LIMIT 1;

    RETURN v_accounting_tax_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_accounting_tax_for_pos IS 'Resolve accounting tax ID from POS tax code, category, or rate';

-- ============================================================================
-- HELPER VIEW: Tax Mappings with Details
-- ============================================================================

CREATE OR REPLACE VIEW view_pos_tax_mappings AS
SELECT
    ptm.id,
    ptm.organization_id,
    o.name as organization_name,
    ptm.pos_tax_code,
    ptm.tax_category_code,
    ptm.pos_tax_rate,
    ptm.accounting_tax_id,
    t.tax_code as accounting_tax_code,
    t.tax_name,
    t.tax_rate as accounting_tax_rate,
    t.tax_scope,
    ptm.default_tax_account_id,
    ta.account_code as tax_account_code,
    ta.account_name as tax_account_name,
    ptm.is_default,
    ptm.is_active,
    ptm.is_inclusive,
    ptm.applies_to_sales,
    ptm.applies_to_purchases,
    ptm.effective_from,
    ptm.effective_to,
    ptm.description
FROM pos_tax_mappings ptm
INNER JOIN organizations o ON ptm.organization_id = o.id
LEFT JOIN taxes t ON ptm.accounting_tax_id = t.id
LEFT JOIN chart_of_accounts ta ON ptm.default_tax_account_id = ta.id
WHERE ptm.deleted_at IS NULL;

COMMENT ON VIEW view_pos_tax_mappings IS 'POS tax mappings with accounting tax and account details';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

ALTER TABLE pos_tax_mappings ENABLE ROW LEVEL SECURITY;

CREATE POLICY pos_tax_mappings_tenant_isolation ON pos_tax_mappings
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY pos_tax_mappings_tenant_isolation ON pos_tax_mappings
    IS 'Ensure users can only access tax mappings for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V008 completed successfully!';
    RAISE NOTICE 'POS Tax Mappings created.';
    RAISE NOTICE ' - pos_tax_mappings table';
    RAISE NOTICE ' - get_accounting_tax_for_pos() function';
    RAISE NOTICE ' - view_pos_tax_mappings view';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

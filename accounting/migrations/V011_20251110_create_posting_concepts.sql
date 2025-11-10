-- ============================================================================
-- Migration: V011 - Create Posting Concepts Layer
-- Description: Human-friendly concept layer for posting engine (REVENUE, AR, CASH, etc.)
--              Bridge between business language and chart of accounts
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- POSTING CONCEPTS TABLE
-- Description: Logical posting concepts - organization-agnostic vocabulary
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_concepts (
    concept_key VARCHAR(50) PRIMARY KEY,

    -- User-friendly labels
    default_label VARCHAR(150) NOT NULL,
    default_description TEXT,

    -- Expected account characteristics
    expected_account_type_id UUID REFERENCES account_types(id),
    normal_side VARCHAR(10) CHECK (normal_side IN ('debit', 'credit')),

    -- Examples for documentation
    example_code VARCHAR(50),
    example_account_name VARCHAR(255),

    -- System vs custom concepts
    is_system BOOLEAN NOT NULL DEFAULT TRUE,

    -- Category for grouping in UI
    concept_category VARCHAR(50), -- 'asset', 'liability', 'equity', 'revenue', 'expense', 'cogs'

    -- Display order in UI
    sort_order INTEGER DEFAULT 0,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_posting_concepts_category ON posting_concepts(concept_category) WHERE deleted_at IS NULL;
CREATE INDEX idx_posting_concepts_is_system ON posting_concepts(is_system);
CREATE INDEX idx_posting_concepts_sort_order ON posting_concepts(sort_order);

-- Trigger
CREATE TRIGGER update_posting_concepts_updated_at
    BEFORE UPDATE ON posting_concepts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_concepts IS 'Logical posting concepts - vocabulary for posting engine (REVENUE, AR, CASH, etc.)';
COMMENT ON COLUMN posting_concepts.concept_key IS 'Unique concept identifier (AR, REVENUE, CASH, PAYROLL_EXPENSE, etc.)';
COMMENT ON COLUMN posting_concepts.default_label IS 'Default user-friendly label (e.g., "Customer Receivables")';
COMMENT ON COLUMN posting_concepts.normal_side IS 'Expected normal balance side (debit for assets, credit for liabilities/revenue)';
COMMENT ON COLUMN posting_concepts.is_system IS 'System concepts cannot be deleted; custom concepts can be added by users';
COMMENT ON COLUMN posting_concepts.concept_category IS 'Category for UI grouping (asset, liability, equity, revenue, expense, cogs)';

-- ============================================================================
-- POSTING CONCEPT OVERRIDES TABLE
-- Description: Organization-specific labels and descriptions for concepts
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_concept_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Concept being overridden
    concept_key VARCHAR(50) NOT NULL REFERENCES posting_concepts(concept_key) ON DELETE CASCADE,

    -- Override values
    label VARCHAR(150) NOT NULL,
    description TEXT,

    -- Display
    is_active BOOLEAN DEFAULT TRUE,

    -- Notes
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
    UNIQUE(organization_id, concept_key)
);

-- Indexes
CREATE INDEX idx_posting_concept_overrides_organization_id ON posting_concept_overrides(organization_id);
CREATE INDEX idx_posting_concept_overrides_concept_key ON posting_concept_overrides(concept_key);

-- Trigger
CREATE TRIGGER update_posting_concept_overrides_updated_at
    BEFORE UPDATE ON posting_concept_overrides
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_concept_overrides IS 'Organization-specific labels and descriptions for posting concepts';
COMMENT ON COLUMN posting_concept_overrides.label IS 'Organization-specific label (e.g., "Sales Income" instead of "Revenue")';
COMMENT ON COLUMN posting_concept_overrides.description IS 'Non-accountant friendly explanation for this organization';

-- ============================================================================
-- ENHANCE EXISTING POSTING_ACCOUNT_MAPPINGS WITH CONCEPT_KEY
-- ============================================================================

-- Add concept_key column to existing posting_account_mappings
ALTER TABLE pos_account_mappings
    ADD COLUMN IF NOT EXISTS concept_key VARCHAR(50);

-- Add foreign key constraint
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_pos_account_mappings_concept'
    ) THEN
        ALTER TABLE pos_account_mappings
            ADD CONSTRAINT fk_pos_account_mappings_concept
            FOREIGN KEY (concept_key)
            REFERENCES posting_concepts(concept_key);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_pos_account_mappings_concept_key ON pos_account_mappings(concept_key) WHERE concept_key IS NOT NULL;

COMMENT ON COLUMN pos_account_mappings.concept_key IS 'Logical concept this mapping resolves (AR, REVENUE, CASH, etc.)';

-- ============================================================================
-- HELPER FUNCTION: Get Concept Label
-- ============================================================================

CREATE OR REPLACE FUNCTION get_concept_label(
    p_concept_key VARCHAR,
    p_organization_id UUID DEFAULT NULL
) RETURNS VARCHAR AS $$
DECLARE
    v_label VARCHAR;
BEGIN
    -- Try organization-specific override first
    IF p_organization_id IS NOT NULL THEN
        SELECT label INTO v_label
        FROM posting_concept_overrides
        WHERE organization_id = p_organization_id
          AND concept_key = p_concept_key
          AND is_active = TRUE
          AND deleted_at IS NULL
        LIMIT 1;

        IF v_label IS NOT NULL THEN
            RETURN v_label;
        END IF;
    END IF;

    -- Fall back to default label
    SELECT default_label INTO v_label
    FROM posting_concepts
    WHERE concept_key = p_concept_key
      AND deleted_at IS NULL;

    RETURN COALESCE(v_label, p_concept_key);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_concept_label IS 'Get display label for a posting concept (org override or default)';

-- ============================================================================
-- HELPER FUNCTION: Get Concept Description
-- ============================================================================

CREATE OR REPLACE FUNCTION get_concept_description(
    p_concept_key VARCHAR,
    p_organization_id UUID DEFAULT NULL
) RETURNS TEXT AS $$
DECLARE
    v_description TEXT;
BEGIN
    -- Try organization-specific override first
    IF p_organization_id IS NOT NULL THEN
        SELECT description INTO v_description
        FROM posting_concept_overrides
        WHERE organization_id = p_organization_id
          AND concept_key = p_concept_key
          AND is_active = TRUE
          AND deleted_at IS NULL
        LIMIT 1;

        IF v_description IS NOT NULL THEN
            RETURN v_description;
        END IF;
    END IF;

    -- Fall back to default description
    SELECT default_description INTO v_description
    FROM posting_concepts
    WHERE concept_key = p_concept_key
      AND deleted_at IS NULL;

    RETURN v_description;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_concept_description IS 'Get description for a posting concept (org override or default)';

-- ============================================================================
-- HELPER VIEW: Concepts with Organization Overrides
-- ============================================================================

CREATE OR REPLACE VIEW view_posting_concepts_with_overrides AS
SELECT
    pc.concept_key,
    pc.default_label,
    pc.default_description,
    pc.expected_account_type_id,
    at.type_name as expected_account_type,
    pc.normal_side,
    pc.example_code,
    pc.example_account_name,
    pc.is_system,
    pc.concept_category,
    pc.sort_order,
    -- Include override information
    o.id as organization_id,
    COALESCE(pco.label, pc.default_label) as display_label,
    COALESCE(pco.description, pc.default_description) as display_description,
    pco.id as override_id,
    CASE WHEN pco.id IS NOT NULL THEN TRUE ELSE FALSE END as has_override
FROM posting_concepts pc
LEFT JOIN account_types at ON pc.expected_account_type_id = at.id
CROSS JOIN organizations o
LEFT JOIN posting_concept_overrides pco
    ON pco.organization_id = o.id
    AND pco.concept_key = pc.concept_key
    AND pco.deleted_at IS NULL
WHERE pc.deleted_at IS NULL
  AND o.deleted_at IS NULL;

COMMENT ON VIEW view_posting_concepts_with_overrides IS 'Posting concepts with organization-specific overrides applied';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Posting concepts are globally visible (no RLS needed)

-- Posting concept overrides use RLS
ALTER TABLE posting_concept_overrides ENABLE ROW LEVEL SECURITY;

CREATE POLICY posting_concept_overrides_tenant_isolation ON posting_concept_overrides
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY posting_concept_overrides_tenant_isolation ON posting_concept_overrides
    IS 'Ensure users can only access concept overrides for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V011 completed successfully!';
    RAISE NOTICE 'Posting Concepts Layer created:';
    RAISE NOTICE ' - posting_concepts table';
    RAISE NOTICE ' - posting_concept_overrides table';
    RAISE NOTICE ' - Enhanced pos_account_mappings with concept_key';
    RAISE NOTICE ' - get_concept_label() function';
    RAISE NOTICE ' - get_concept_description() function';
    RAISE NOTICE ' - view_posting_concepts_with_overrides';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

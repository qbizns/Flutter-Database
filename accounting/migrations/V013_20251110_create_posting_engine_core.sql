-- ============================================================================
-- Migration: V013 - Create Posting Engine Core Tables
-- Description: Core posting engine with profiles, document types, rules, and rule lines
--              Configuration-driven double-entry posting system
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- ENUMS FOR POSTING ENGINE
-- ============================================================================

-- Posting event (when to post)
CREATE TYPE posting_event AS ENUM (
    'on_post',         -- When document is posted/approved
    'on_reverse',      -- When document is reversed/cancelled
    'on_pay',          -- When payment is made/received
    'on_clear',        -- When cleared/reconciled
    'on_settle'        -- When settled
);

-- Posting side
CREATE TYPE posting_side AS ENUM ('debit', 'credit');

-- Posting level (header vs line)
CREATE TYPE posting_level AS ENUM (
    'header',          -- One entry per document (uses header fields)
    'line'             -- One entry per document line (iterates line items)
);

-- Account source (how to determine the account)
CREATE TYPE posting_account_source AS ENUM (
    'from_mapping',    -- Resolve from posting_account_mappings using concept_key
    'from_document',   -- Get account ID directly from document field
    'fixed',           -- Use a fixed account_id specified in rule
    'expression'       -- Evaluate an expression to get account_id
);

-- Amount source (where to get the amount)
CREATE TYPE posting_amount_source AS ENUM (
    'document_field',  -- Get from document field (e.g., doc.total_amount)
    'line_field',      -- Get from line field (e.g., line.amount)
    'expression'       -- Evaluate expression (e.g., doc.subtotal * 0.15)
);

-- ============================================================================
-- POSTING PROFILES TABLE
-- Description: Organization-specific posting configuration profiles
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Profile identification
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,

    -- Configuration
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Default fiscal settings
    default_fiscal_year_id UUID REFERENCES fiscal_years(id),

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
    UNIQUE(organization_id, code)
);

-- Indexes
CREATE INDEX idx_posting_profiles_organization_id ON posting_profiles(organization_id);
CREATE INDEX idx_posting_profiles_is_default ON posting_profiles(organization_id, is_default) WHERE is_default = TRUE;

-- Trigger
CREATE TRIGGER update_posting_profiles_updated_at
    BEFORE UPDATE ON posting_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_profiles IS 'Organization-specific posting configuration profiles';
COMMENT ON COLUMN posting_profiles.code IS 'Unique profile code per organization (e.g., DEFAULT_PROFILE, RETAIL_PROFILE)';
COMMENT ON COLUMN posting_profiles.is_default IS 'Default profile for organization if not explicitly specified';

-- ============================================================================
-- POSTING DOCUMENT TYPES TABLE
-- Description: Business document types that can be posted
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_document_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Document type identification
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description TEXT,

    -- Source table information
    source_schema VARCHAR(50) NOT NULL,  -- e.g., 'public', 'accounting'
    source_table VARCHAR(100) NOT NULL,  -- e.g., 'sales', 'vendor_bills', 'payroll_runs'
    source_pk_column VARCHAR(50) NOT NULL DEFAULT 'id',

    -- Document category
    category VARCHAR(50),  -- 'sales', 'purchases', 'payroll', 'banking', 'expenses'

    -- Configuration
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_system BOOLEAN NOT NULL DEFAULT TRUE,

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
CREATE INDEX idx_posting_document_types_code ON posting_document_types(code);
CREATE INDEX idx_posting_document_types_category ON posting_document_types(category) WHERE deleted_at IS NULL;
CREATE INDEX idx_posting_document_types_source ON posting_document_types(source_schema, source_table);

-- Trigger
CREATE TRIGGER update_posting_document_types_updated_at
    BEFORE UPDATE ON posting_document_types
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_document_types IS 'Business document types that can be posted (POS_SALE, VENDOR_BILL, PAYROLL_RUN, etc.)';
COMMENT ON COLUMN posting_document_types.code IS 'Unique document type code (POS_SALE, VENDOR_BILL, PAYROLL_RUN, CHECK_ISSUE, etc.)';
COMMENT ON COLUMN posting_document_types.source_table IS 'Source table containing the business documents';

-- ============================================================================
-- POSTING PROFILE DOCUMENTS TABLE
-- Description: Link posting profiles to document types they handle
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_profile_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Profile and document type
    posting_profile_id UUID NOT NULL REFERENCES posting_profiles(id) ON DELETE CASCADE,
    posting_document_type_id UUID NOT NULL REFERENCES posting_document_types(id) ON DELETE CASCADE,

    -- Configuration
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    UNIQUE(posting_profile_id, posting_document_type_id)
);

-- Indexes
CREATE INDEX idx_posting_profile_documents_profile ON posting_profile_documents(posting_profile_id);
CREATE INDEX idx_posting_profile_documents_doctype ON posting_profile_documents(posting_document_type_id);

-- Trigger
CREATE TRIGGER update_posting_profile_documents_updated_at
    BEFORE UPDATE ON posting_profile_documents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_profile_documents IS 'Links posting profiles to document types they handle';

-- ============================================================================
-- POSTING RULES TABLE
-- Description: Posting templates for document types + events
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Profile-document link
    posting_profile_document_id UUID NOT NULL REFERENCES posting_profile_documents(id) ON DELETE CASCADE,

    -- Rule identification
    rule_code VARCHAR(50) NOT NULL,
    rule_name VARCHAR(150) NOT NULL,
    description TEXT,

    -- When to apply this rule
    event posting_event NOT NULL DEFAULT 'on_post',

    -- Header or line level
    level posting_level NOT NULL DEFAULT 'header',

    -- Priority (lower number = higher priority)
    priority INTEGER NOT NULL DEFAULT 100,

    -- Condition expression (optional - if specified, only apply if true)
    -- DSL evaluated by Go, e.g., "doc.payment_type == 'CASH'"
    condition_expression TEXT,

    -- Configuration
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_posting_rules_profile_doc ON posting_rules(posting_profile_document_id);
CREATE INDEX idx_posting_rules_event ON posting_rules(event);
CREATE INDEX idx_posting_rules_level ON posting_rules(level);
CREATE INDEX idx_posting_rules_priority ON posting_rules(priority, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_posting_rules_code ON posting_rules(rule_code);

-- Trigger
CREATE TRIGGER update_posting_rules_updated_at
    BEFORE UPDATE ON posting_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_rules IS 'Posting templates for document types + events';
COMMENT ON COLUMN posting_rules.rule_code IS 'Unique rule code (e.g., POS_SALE_CASH, PAYROLL_ACCRUAL, CHECK_ISSUE)';
COMMENT ON COLUMN posting_rules.event IS 'When to apply: on_post, on_reverse, on_pay, on_clear, on_settle';
COMMENT ON COLUMN posting_rules.level IS 'header = once per document, line = once per document line';
COMMENT ON COLUMN posting_rules.condition_expression IS 'Optional condition (DSL) - only apply if true';

-- ============================================================================
-- POSTING RULE LINES TABLE
-- Description: Journal entry line templates for each posting rule
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_rule_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Parent posting rule
    posting_rule_id UUID NOT NULL REFERENCES posting_rules(id) ON DELETE CASCADE,

    -- Line ordering
    line_no INTEGER NOT NULL,

    -- Debit or credit
    side posting_side NOT NULL,

    -- Concept-based account resolution
    concept_key VARCHAR(50) REFERENCES posting_concepts(concept_key),

    -- Account determination
    account_source posting_account_source NOT NULL DEFAULT 'from_mapping',
    fixed_account_id UUID REFERENCES chart_of_accounts(id),  -- Used if account_source = 'fixed'
    account_field_path VARCHAR(255),  -- e.g., 'doc.cash_account_id' if account_source = 'from_document'
    account_expression TEXT,  -- DSL if account_source = 'expression'

    -- Amount determination
    amount_source posting_amount_source NOT NULL DEFAULT 'document_field',
    amount_field_path VARCHAR(255),  -- e.g., 'doc.total_amount', 'line.amount'
    amount_expression TEXT,  -- DSL e.g., 'doc.subtotal * 0.15' or 'line.quantity * line.unit_price'

    -- Additional mapping context (for from_mapping)
    mapping_context JSONB DEFAULT '{}',  -- e.g., {"product_id": "line.product_id", "payment_method": "doc.payment_method"}

    -- Description for journal entry line
    description_template VARCHAR(255),  -- e.g., 'Sale {doc.sale_number}', 'Payroll {doc.period_name}'

    -- Configuration
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    UNIQUE(posting_rule_id, line_no),
    CHECK (
        (account_source = 'from_mapping' AND concept_key IS NOT NULL) OR
        (account_source = 'fixed' AND fixed_account_id IS NOT NULL) OR
        (account_source = 'from_document' AND account_field_path IS NOT NULL) OR
        (account_source = 'expression' AND account_expression IS NOT NULL)
    ),
    CHECK (
        (amount_source = 'document_field' AND amount_field_path IS NOT NULL) OR
        (amount_source = 'line_field' AND amount_field_path IS NOT NULL) OR
        (amount_source = 'expression' AND amount_expression IS NOT NULL)
    )
);

-- Indexes
CREATE INDEX idx_posting_rule_lines_rule ON posting_rule_lines(posting_rule_id);
CREATE INDEX idx_posting_rule_lines_line_no ON posting_rule_lines(posting_rule_id, line_no);
CREATE INDEX idx_posting_rule_lines_concept ON posting_rule_lines(concept_key) WHERE concept_key IS NOT NULL;

-- Trigger
CREATE TRIGGER update_posting_rule_lines_updated_at
    BEFORE UPDATE ON posting_rule_lines
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_rule_lines IS 'Journal entry line templates for posting rules';
COMMENT ON COLUMN posting_rule_lines.concept_key IS 'Logical concept (AR, REVENUE, CASH, etc.) - used with account_source=from_mapping';
COMMENT ON COLUMN posting_rule_lines.account_source IS 'How to determine GL account: from_mapping (via concept), fixed, from_document, or expression';
COMMENT ON COLUMN posting_rule_lines.amount_source IS 'Where to get amount: document_field, line_field, or expression';
COMMENT ON COLUMN posting_rule_lines.mapping_context IS 'Additional context for account mapping (product_id, payment_method, etc.)';

-- ============================================================================
-- HELPER FUNCTION: Get Posting Rules for Document
-- ============================================================================

CREATE OR REPLACE FUNCTION get_posting_rules_for_document(
    p_organization_id UUID,
    p_document_type_code VARCHAR,
    p_event posting_event
) RETURNS TABLE (
    rule_id UUID,
    rule_code VARCHAR,
    rule_name VARCHAR,
    level posting_level,
    priority INTEGER,
    condition_expression TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        pr.id as rule_id,
        pr.rule_code,
        pr.rule_name,
        pr.level,
        pr.priority,
        pr.condition_expression
    FROM posting_rules pr
    INNER JOIN posting_profile_documents ppd ON pr.posting_profile_document_id = ppd.id
    INNER JOIN posting_profiles pp ON ppd.posting_profile_id = pp.id
    INNER JOIN posting_document_types pdt ON ppd.posting_document_type_id = pdt.id
    WHERE pp.organization_id = p_organization_id
      AND pdt.code = p_document_type_code
      AND pr.event = p_event
      AND pr.is_active = TRUE
      AND ppd.is_active = TRUE
      AND pp.is_active = TRUE
      AND pr.deleted_at IS NULL
      AND ppd.deleted_at IS NULL
      AND pp.deleted_at IS NULL
      AND pdt.deleted_at IS NULL
    ORDER BY pr.priority ASC, pr.created_at ASC;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_posting_rules_for_document IS 'Get all active posting rules for a document type and event';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS on tenant-scoped tables
ALTER TABLE posting_profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE posting_profile_documents ENABLE ROW LEVEL SECURITY;

CREATE POLICY posting_profiles_tenant_isolation ON posting_profiles
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

CREATE POLICY posting_profile_documents_tenant_isolation ON posting_profile_documents
    USING (posting_profile_id IN (
        SELECT id FROM posting_profiles
        WHERE organization_id IN (
            SELECT organization_id
            FROM user_organizations
            WHERE user_id = auth.uid()
        )
    ));

-- posting_document_types, posting_rules, posting_rule_lines are globally visible (indirect access via profile)

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V013 completed successfully!';
    RAISE NOTICE 'Posting Engine Core Tables created:';
    RAISE NOTICE ' - Enums: posting_event, posting_side, posting_level, posting_account_source, posting_amount_source';
    RAISE NOTICE ' - posting_profiles table';
    RAISE NOTICE ' - posting_document_types table';
    RAISE NOTICE ' - posting_profile_documents table';
    RAISE NOTICE ' - posting_rules table';
    RAISE NOTICE ' - posting_rule_lines table';
    RAISE NOTICE ' - get_posting_rules_for_document() function';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

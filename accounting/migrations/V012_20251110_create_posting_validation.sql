-- ============================================================================
-- Migration: V012 - Create Posting Validation Layer
-- Description: Configuration-driven validation rules for posting engine
--              Validates balance, periods, account types, fiscal year locks, etc.
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- ENUMS FOR VALIDATION
-- ============================================================================

-- Validation target (what we're validating)
CREATE TYPE validation_target AS ENUM (
    'document',        -- Business document (sale, bill, payroll, etc.)
    'journal_entry',   -- Journal entry header
    'journal_line'     -- Individual journal entry line
);

-- Validation severity
CREATE TYPE validation_severity AS ENUM (
    'error',           -- Blocks posting
    'warning',         -- Shows warning but allows posting
    'info'             -- Informational only
);

-- ============================================================================
-- POSTING VALIDATION RULES TABLE
-- Description: Configuration-driven validation rules
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_validation_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Scope
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,  -- NULL = global rule
    document_type_code VARCHAR(50),  -- NULL = applies to any document type
    event VARCHAR(50),               -- 'on_post', 'on_reverse', 'on_pay', 'on_clear', NULL = any

    -- Validation specification
    target validation_target NOT NULL DEFAULT 'journal_entry',
    code VARCHAR(50) NOT NULL,       -- e.g., 'BALANCED_ENTRY', 'OPEN_PERIOD', 'ACCOUNT_TYPE_MATCH'
    name VARCHAR(150) NOT NULL,
    description TEXT,

    -- Expression to evaluate (DSL evaluated by Go)
    -- Context variables available: doc, je, line, period, org, account, etc.
    expression TEXT NOT NULL,

    -- Expression examples:
    -- 'abs(je.total_debit - je.total_credit) <= 0.005'
    -- 'period.status == "open"'
    -- 'account.account_type.type_category == line.expected_category'
    -- 'fiscal_year.status != "closed"'

    -- Severity and behavior
    severity validation_severity NOT NULL DEFAULT 'error',
    is_blocking BOOLEAN NOT NULL DEFAULT TRUE,  -- If true and severity=error, blocks posting
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Error message template (can use placeholders)
    message_template TEXT,  -- e.g., 'Cannot post: accounting period {period.name} is {period.status}'

    -- Priority (lower number = evaluated first)
    priority INTEGER DEFAULT 100,

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
    CONSTRAINT uq_validation_rules UNIQUE (
        COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'::uuid),
        COALESCE(document_type_code, ''),
        COALESCE(event, ''),
        code
    )
);

-- Indexes
CREATE INDEX idx_validation_rules_lookup
    ON posting_validation_rules (organization_id, document_type_code, event, is_active, severity)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_validation_rules_code ON posting_validation_rules(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_validation_rules_priority ON posting_validation_rules(priority, is_active) WHERE is_active = TRUE;

-- Trigger
CREATE TRIGGER update_posting_validation_rules_updated_at
    BEFORE UPDATE ON posting_validation_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE posting_validation_rules IS 'Configuration-driven validation rules for posting engine';
COMMENT ON COLUMN posting_validation_rules.organization_id IS 'Organization scope (NULL = global/system rule)';
COMMENT ON COLUMN posting_validation_rules.document_type_code IS 'Document type scope (NULL = all document types)';
COMMENT ON COLUMN posting_validation_rules.event IS 'Event scope (on_post, on_reverse, etc., NULL = all events)';
COMMENT ON COLUMN posting_validation_rules.target IS 'What to validate: document, journal_entry, or journal_line';
COMMENT ON COLUMN posting_validation_rules.expression IS 'Boolean DSL expression evaluated by Go (returns true = valid, false = invalid)';
COMMENT ON COLUMN posting_validation_rules.is_blocking IS 'If true and severity=error, blocks posting';
COMMENT ON COLUMN posting_validation_rules.message_template IS 'Error message template with placeholders for context variables';

-- ============================================================================
-- POSTING VALIDATION RESULTS TABLE
-- Description: Log of validation results (failures and warnings)
-- ============================================================================

CREATE TABLE IF NOT EXISTS posting_validation_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Document context
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    document_type_code VARCHAR(50) NOT NULL,
    document_id UUID NOT NULL,
    event VARCHAR(50) NOT NULL,

    -- Journal entry reference (if created)
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,

    -- Validation rule that failed/warned
    validation_rule_id UUID REFERENCES posting_validation_rules(id) ON DELETE SET NULL,

    -- Result details
    severity validation_severity NOT NULL,
    message_code VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    is_blocking BOOLEAN NOT NULL DEFAULT TRUE,

    -- Context data (for debugging and display)
    context JSONB DEFAULT '{}',  -- e.g., {"line_index": 2, "account_code": "9999", "expected": "REVENUE", "actual": "EXPENSE"}

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by UUID REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_validation_results_document
    ON posting_validation_results (organization_id, document_type_code, document_id);
CREATE INDEX idx_validation_results_journal_entry
    ON posting_validation_results (journal_entry_id) WHERE journal_entry_id IS NOT NULL;
CREATE INDEX idx_validation_results_severity
    ON posting_validation_results (severity, is_blocking);
CREATE INDEX idx_validation_results_message_code
    ON posting_validation_results (message_code);
CREATE INDEX idx_validation_results_created_at
    ON posting_validation_results (created_at);

-- Comments
COMMENT ON TABLE posting_validation_results IS 'Log of validation results (failures and warnings) from posting engine';
COMMENT ON COLUMN posting_validation_results.document_id IS 'ID of the source business document (sale, bill, payroll, etc.)';
COMMENT ON COLUMN posting_validation_results.journal_entry_id IS 'Journal entry ID if created (may be NULL if validation failed before JE creation)';
COMMENT ON COLUMN posting_validation_results.is_blocking IS 'Whether this validation result blocked the posting';
COMMENT ON COLUMN posting_validation_results.context IS 'Additional context data for debugging and user-friendly error messages';

-- ============================================================================
-- HELPER FUNCTION: Get Active Validation Rules
-- ============================================================================

CREATE OR REPLACE FUNCTION get_active_validation_rules(
    p_organization_id UUID,
    p_document_type_code VARCHAR,
    p_event VARCHAR
) RETURNS TABLE (
    id UUID,
    code VARCHAR,
    name VARCHAR,
    target validation_target,
    expression TEXT,
    severity validation_severity,
    is_blocking BOOLEAN,
    message_template TEXT,
    priority INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        vr.id,
        vr.code,
        vr.name,
        vr.target,
        vr.expression,
        vr.severity,
        vr.is_blocking,
        vr.message_template,
        vr.priority
    FROM posting_validation_rules vr
    WHERE vr.is_active = TRUE
      AND vr.deleted_at IS NULL
      -- Match organization (specific org or global)
      AND (vr.organization_id = p_organization_id OR vr.organization_id IS NULL)
      -- Match document type (specific type or any)
      AND (vr.document_type_code = p_document_type_code OR vr.document_type_code IS NULL)
      -- Match event (specific event or any)
      AND (vr.event = p_event OR vr.event IS NULL)
    ORDER BY vr.priority ASC, vr.created_at ASC;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_active_validation_rules IS 'Get all active validation rules for a posting operation';

-- ============================================================================
-- HELPER FUNCTION: Check if Posting is Blocked
-- ============================================================================

CREATE OR REPLACE FUNCTION is_posting_blocked(
    p_organization_id UUID,
    p_document_type_code VARCHAR,
    p_document_id UUID,
    p_event VARCHAR
) RETURNS BOOLEAN AS $$
DECLARE
    v_has_blocking_errors BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM posting_validation_results
        WHERE organization_id = p_organization_id
          AND document_type_code = p_document_type_code
          AND document_id = p_document_id
          AND event = p_event
          AND severity = 'error'
          AND is_blocking = TRUE
    ) INTO v_has_blocking_errors;

    RETURN v_has_blocking_errors;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION is_posting_blocked IS 'Check if a posting operation is blocked by validation errors';

-- ============================================================================
-- HELPER FUNCTION: Clear Validation Results for Document
-- ============================================================================

CREATE OR REPLACE FUNCTION clear_validation_results(
    p_organization_id UUID,
    p_document_type_code VARCHAR,
    p_document_id UUID,
    p_event VARCHAR DEFAULT NULL
) RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER;
BEGIN
    DELETE FROM posting_validation_results
    WHERE organization_id = p_organization_id
      AND document_type_code = p_document_type_code
      AND document_id = p_document_id
      AND (p_event IS NULL OR event = p_event);

    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;
    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION clear_validation_results IS 'Clear validation results for a document (useful before re-validation)';

-- ============================================================================
-- HELPER VIEW: Validation Results with Rule Details
-- ============================================================================

CREATE OR REPLACE VIEW view_posting_validation_results AS
SELECT
    pvr.id,
    pvr.organization_id,
    o.name as organization_name,
    pvr.document_type_code,
    pvr.document_id,
    pvr.event,
    pvr.journal_entry_id,
    pvr.validation_rule_id,
    pvrule.code as rule_code,
    pvrule.name as rule_name,
    pvr.severity,
    pvr.message_code,
    pvr.message,
    pvr.is_blocking,
    pvr.context,
    pvr.created_at,
    u.email as created_by_email
FROM posting_validation_results pvr
INNER JOIN organizations o ON pvr.organization_id = o.id
LEFT JOIN posting_validation_rules pvrule ON pvr.validation_rule_id = pvrule.id
LEFT JOIN users u ON pvr.created_by = u.id
ORDER BY pvr.created_at DESC;

COMMENT ON VIEW view_posting_validation_results IS 'Validation results with rule and user details';

-- ============================================================================
-- HELPER VIEW: Blocked Postings Summary
-- ============================================================================

CREATE OR REPLACE VIEW view_blocked_postings AS
SELECT
    pvr.organization_id,
    o.name as organization_name,
    pvr.document_type_code,
    pvr.document_id,
    pvr.event,
    COUNT(*) as error_count,
    COUNT(DISTINCT pvr.validation_rule_id) as failed_rule_count,
    MAX(pvr.created_at) as last_validation_at,
    json_agg(
        json_build_object(
            'rule_code', pvrule.code,
            'message', pvr.message,
            'severity', pvr.severity
        ) ORDER BY pvr.created_at
    ) as errors
FROM posting_validation_results pvr
INNER JOIN organizations o ON pvr.organization_id = o.id
LEFT JOIN posting_validation_rules pvrule ON pvr.validation_rule_id = pvrule.id
WHERE pvr.severity = 'error'
  AND pvr.is_blocking = TRUE
GROUP BY pvr.organization_id, o.name, pvr.document_type_code, pvr.document_id, pvr.event
ORDER BY last_validation_at DESC;

COMMENT ON VIEW view_blocked_postings IS 'Summary of blocked postings with error details';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Validation rules are globally visible (no RLS needed on rules)

-- Validation results use RLS
ALTER TABLE posting_validation_results ENABLE ROW LEVEL SECURITY;

CREATE POLICY posting_validation_results_tenant_isolation ON posting_validation_results
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY posting_validation_results_tenant_isolation ON posting_validation_results
    IS 'Ensure users can only access validation results for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V012 completed successfully!';
    RAISE NOTICE 'Posting Validation Layer created:';
    RAISE NOTICE ' - validation_target and validation_severity enums';
    RAISE NOTICE ' - posting_validation_rules table';
    RAISE NOTICE ' - posting_validation_results table';
    RAISE NOTICE ' - get_active_validation_rules() function';
    RAISE NOTICE ' - is_posting_blocked() function';
    RAISE NOTICE ' - clear_validation_results() function';
    RAISE NOTICE ' - view_posting_validation_results';
    RAISE NOTICE ' - view_blocked_postings';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

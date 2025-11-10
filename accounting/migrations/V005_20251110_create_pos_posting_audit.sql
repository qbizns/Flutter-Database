-- ============================================================================
-- Migration: V005 - Create POS Posting Audit
-- Description: Track which POS documents have been posted to accounting
--              and maintain audit trail of all posting attempts
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- POS POSTING AUDIT TABLE
-- Description: Audit trail for POS → Accounting posting operations
-- ============================================================================

CREATE TABLE IF NOT EXISTS pos_posting_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Source Document (polymorphic reference)
    source_table VARCHAR(100) NOT NULL,  -- e.g., 'sales', 'goods_receipts', 'pos_sessions', 'expenses'
    source_id UUID NOT NULL,             -- ID of the source document
    source_reference VARCHAR(255),       -- Human-readable reference (sale_number, receipt_number, etc.)

    -- Posting Status
    posting_status VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (posting_status IN (
        'pending',      -- Queued for posting
        'processing',   -- Currently being posted
        'posted',       -- Successfully posted
        'failed',       -- Failed to post
        'cancelled',    -- Posting cancelled
        'reversed'      -- Posted but later reversed
    )),

    -- Accounting References
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    reversal_journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,

    -- Posting Details
    posting_date DATE,                    -- Date used for posting
    posted_at TIMESTAMP WITH TIME ZONE,   -- When posting completed
    posted_by UUID REFERENCES users(id),  -- Who initiated the posting
    posting_method VARCHAR(50) DEFAULT 'automatic',  -- 'automatic', 'manual', 'batch'

    -- Error Handling
    error_code VARCHAR(50),
    error_message TEXT,
    error_details JSONB,
    retry_count INTEGER DEFAULT 0,
    last_retry_at TIMESTAMP WITH TIME ZONE,
    max_retries INTEGER DEFAULT 3,

    -- Reversal Information
    reversed_at TIMESTAMP WITH TIME ZONE,
    reversed_by UUID REFERENCES users(id),
    reversal_reason TEXT,

    -- Posting Summary (for quick reference without joining)
    total_debit NUMERIC(20, 4),
    total_credit NUMERIC(20, 4),
    line_count INTEGER,
    currency_code VARCHAR(3),

    -- Additional Context
    posting_context JSONB DEFAULT '{}',  -- e.g., batch_id, job_id, correlation_id
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
    CONSTRAINT valid_posted_state CHECK (
        (posting_status = 'posted' AND journal_entry_id IS NOT NULL AND posted_at IS NOT NULL) OR
        (posting_status != 'posted')
    ),
    CONSTRAINT valid_reversal CHECK (
        (posting_status = 'reversed' AND reversal_journal_entry_id IS NOT NULL AND reversed_at IS NOT NULL) OR
        (posting_status != 'reversed')
    ),
    CONSTRAINT balanced_amounts CHECK (
        (posting_status != 'posted') OR
        (ABS(COALESCE(total_debit, 0) - COALESCE(total_credit, 0)) < 0.01)
    )
);

-- Indexes
CREATE INDEX idx_pos_posting_audit_organization_id ON pos_posting_audit(organization_id);
CREATE INDEX idx_pos_posting_audit_source ON pos_posting_audit(source_table, source_id);
CREATE INDEX idx_pos_posting_audit_status ON pos_posting_audit(posting_status);
CREATE INDEX idx_pos_posting_audit_journal_entry_id ON pos_posting_audit(journal_entry_id) WHERE journal_entry_id IS NOT NULL;
CREATE INDEX idx_pos_posting_audit_posted_at ON pos_posting_audit(posted_at) WHERE posted_at IS NOT NULL;
CREATE INDEX idx_pos_posting_audit_posting_date ON pos_posting_audit(posting_date) WHERE posting_date IS NOT NULL;
CREATE INDEX idx_pos_posting_audit_failed ON pos_posting_audit(posting_status, retry_count) WHERE posting_status = 'failed';

-- Unique constraint: one active posting per source document
CREATE UNIQUE INDEX idx_pos_posting_audit_source_unique ON pos_posting_audit(source_table, source_id)
    WHERE deleted_at IS NULL AND posting_status NOT IN ('cancelled', 'reversed');

-- Trigger
CREATE TRIGGER update_pos_posting_audit_updated_at
    BEFORE UPDATE ON pos_posting_audit
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE pos_posting_audit IS 'Audit trail for POS documents posted to accounting';
COMMENT ON COLUMN pos_posting_audit.source_table IS 'Source table name (sales, goods_receipts, pos_sessions, etc.)';
COMMENT ON COLUMN pos_posting_audit.source_id IS 'UUID of the source document';
COMMENT ON COLUMN pos_posting_audit.posting_status IS 'Current posting status';
COMMENT ON COLUMN pos_posting_audit.journal_entry_id IS 'Created journal entry (if posted successfully)';
COMMENT ON COLUMN pos_posting_audit.reversal_journal_entry_id IS 'Reversal journal entry (if reversed)';
COMMENT ON COLUMN pos_posting_audit.posting_context IS 'Additional context like batch_id, job_id for async processing';

-- ============================================================================
-- HELPER FUNCTION: Check if Document is Posted
-- ============================================================================

CREATE OR REPLACE FUNCTION is_pos_document_posted(
    p_source_table VARCHAR,
    p_source_id UUID
) RETURNS BOOLEAN AS $$
DECLARE
    v_is_posted BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM pos_posting_audit
        WHERE source_table = p_source_table
          AND source_id = p_source_id
          AND posting_status = 'posted'
          AND deleted_at IS NULL
    ) INTO v_is_posted;

    RETURN v_is_posted;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION is_pos_document_posted IS 'Check if a POS document has been posted to accounting';

-- ============================================================================
-- HELPER FUNCTION: Get Posting Status
-- ============================================================================

CREATE OR REPLACE FUNCTION get_pos_posting_status(
    p_source_table VARCHAR,
    p_source_id UUID
) RETURNS TABLE (
    status VARCHAR,
    journal_entry_id UUID,
    posted_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        ppa.posting_status,
        ppa.journal_entry_id,
        ppa.posted_at,
        ppa.error_message
    FROM pos_posting_audit ppa
    WHERE ppa.source_table = p_source_table
      AND ppa.source_id = p_source_id
      AND ppa.deleted_at IS NULL
    ORDER BY ppa.created_at DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_pos_posting_status IS 'Get current posting status for a POS document';

-- ============================================================================
-- HELPER VIEWS
-- ============================================================================

-- View: Pending Postings
CREATE OR REPLACE VIEW view_pending_pos_postings AS
SELECT
    ppa.id,
    ppa.organization_id,
    o.name as organization_name,
    ppa.source_table,
    ppa.source_id,
    ppa.source_reference,
    ppa.posting_status,
    ppa.created_at,
    ppa.retry_count,
    ppa.max_retries,
    ppa.error_message,
    CASE
        WHEN ppa.posting_status = 'failed' AND ppa.retry_count < ppa.max_retries
             AND (ppa.last_retry_at IS NULL OR ppa.last_retry_at < CURRENT_TIMESTAMP - INTERVAL '5 minutes')
        THEN TRUE
        ELSE FALSE
    END as should_retry
FROM pos_posting_audit ppa
INNER JOIN organizations o ON ppa.organization_id = o.id
WHERE ppa.posting_status IN ('pending', 'failed')
  AND ppa.deleted_at IS NULL
ORDER BY ppa.created_at ASC;

COMMENT ON VIEW view_pending_pos_postings IS 'POS documents pending posting or requiring retry';

-- View: Failed Postings
CREATE OR REPLACE VIEW view_failed_pos_postings AS
SELECT
    ppa.id,
    ppa.organization_id,
    o.name as organization_name,
    ppa.source_table,
    ppa.source_id,
    ppa.source_reference,
    ppa.posting_status,
    ppa.error_code,
    ppa.error_message,
    ppa.retry_count,
    ppa.max_retries,
    ppa.last_retry_at,
    ppa.created_at,
    ppa.updated_at
FROM pos_posting_audit ppa
INNER JOIN organizations o ON ppa.organization_id = o.id
WHERE ppa.posting_status = 'failed'
  AND ppa.retry_count >= ppa.max_retries
  AND ppa.deleted_at IS NULL
ORDER BY ppa.updated_at DESC;

COMMENT ON VIEW view_failed_pos_postings IS 'Failed POS postings that have exhausted retries';

-- View: Posted Documents Summary
CREATE OR REPLACE VIEW view_posted_documents_summary AS
SELECT
    ppa.organization_id,
    o.name as organization_name,
    ppa.source_table,
    COUNT(*) as total_posted,
    SUM(ppa.total_debit) as total_debit_amount,
    SUM(ppa.total_credit) as total_credit_amount,
    MIN(ppa.posted_at) as first_posted_at,
    MAX(ppa.posted_at) as last_posted_at
FROM pos_posting_audit ppa
INNER JOIN organizations o ON ppa.organization_id = o.id
WHERE ppa.posting_status = 'posted'
  AND ppa.deleted_at IS NULL
GROUP BY ppa.organization_id, o.name, ppa.source_table
ORDER BY ppa.organization_id, ppa.source_table;

COMMENT ON VIEW view_posted_documents_summary IS 'Summary statistics for posted POS documents by type';

-- View: Posting Audit Trail (with Journal Entry Details)
CREATE OR REPLACE VIEW view_pos_posting_audit_trail AS
SELECT
    ppa.id,
    ppa.organization_id,
    o.name as organization_name,
    ppa.source_table,
    ppa.source_id,
    ppa.source_reference,
    ppa.posting_status,
    ppa.journal_entry_id,
    je.entry_number,
    je.entry_date,
    je.is_posted as je_is_posted,
    ppa.total_debit,
    ppa.total_credit,
    ppa.posted_at,
    u.email as posted_by_email,
    ppa.error_message,
    ppa.created_at
FROM pos_posting_audit ppa
INNER JOIN organizations o ON ppa.organization_id = o.id
LEFT JOIN journal_entries je ON ppa.journal_entry_id = je.id
LEFT JOIN users u ON ppa.posted_by = u.id
WHERE ppa.deleted_at IS NULL
ORDER BY ppa.created_at DESC;

COMMENT ON VIEW view_pos_posting_audit_trail IS 'Complete audit trail with journal entry details';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

ALTER TABLE pos_posting_audit ENABLE ROW LEVEL SECURITY;

CREATE POLICY pos_posting_audit_tenant_isolation ON pos_posting_audit
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY pos_posting_audit_tenant_isolation ON pos_posting_audit
    IS 'Ensure users can only access posting audit records for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V005 completed successfully!';
    RAISE NOTICE 'POS Posting Audit created.';
    RAISE NOTICE ' - pos_posting_audit table';
    RAISE NOTICE ' - is_pos_document_posted() function';
    RAISE NOTICE ' - get_pos_posting_status() function';
    RAISE NOTICE ' - 4 helper views created';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

-- ============================================================================
-- Migration: V009 - Add Immutability Triggers
-- Description: Enforce immutability rules for posted accounting data
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- GENERAL LEDGER IMMUTABILITY
-- Description: GL entries are immutable - no updates or deletes allowed
-- ============================================================================

-- Prevent any modifications to general ledger
CREATE OR REPLACE FUNCTION prevent_gl_modification() RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'general_ledger rows are immutable. Use reversal journal entries instead of modifying posted transactions. Table: general_ledger, Operation: %', TG_OP
        USING HINT = 'Create a reversal journal entry to correct accounting errors',
              ERRCODE = '23502';
END;
$$ LANGUAGE plpgsql;

-- Trigger to prevent updates
CREATE TRIGGER trg_prevent_gl_update
    BEFORE UPDATE ON general_ledger
    FOR EACH ROW
    EXECUTE FUNCTION prevent_gl_modification();

-- Trigger to prevent deletes
CREATE TRIGGER trg_prevent_gl_delete
    BEFORE DELETE ON general_ledger
    FOR EACH ROW
    EXECUTE FUNCTION prevent_gl_modification();

COMMENT ON FUNCTION prevent_gl_modification IS 'Enforces immutability of general ledger entries';

-- ============================================================================
-- JOURNAL ENTRIES IMMUTABILITY (for posted entries only)
-- Description: Posted journal entries can't be modified - only reversed
-- ============================================================================

-- Prevent modifications to posted journal entries
CREATE OR REPLACE FUNCTION prevent_posted_je_modification() RETURNS trigger AS $$
BEGIN
    -- Allow UPDATE only for specific fields or status changes
    IF TG_OP = 'UPDATE' THEN
        -- Allow changing status from 'posted' to 'reversed'
        IF OLD.is_posted = TRUE AND OLD.status = 'posted' THEN
            -- Only allow status change to reversed
            IF NEW.status = 'reversed' THEN
                RETURN NEW;
            END IF;

            -- Allow adding notes/description without changing financial data
            IF OLD.total_debit = NEW.total_debit
               AND OLD.total_credit = NEW.total_credit
               AND OLD.entry_date = NEW.entry_date
               AND OLD.posting_date = NEW.posting_date
               AND OLD.is_posted = NEW.is_posted THEN
                RETURN NEW;
            END IF;

            -- Block all other modifications
            RAISE EXCEPTION 'Cannot modify posted journal entry. Financial fields are immutable. Create a reversal entry instead. JE ID: %', OLD.id
                USING HINT = 'Use reversal_journal_entry_id to link corrections',
                      ERRCODE = '23502';
        END IF;

        -- Allow modifications if not posted
        RETURN NEW;
    END IF;

    -- Prevent deletes of posted entries
    IF TG_OP = 'DELETE' AND OLD.is_posted = TRUE THEN
        RAISE EXCEPTION 'Cannot delete posted journal entry: %. Create a reversal entry instead.', OLD.entry_number
            USING ERRCODE = '23502';
    END IF;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_posted_je_modification
    BEFORE UPDATE OR DELETE ON journal_entries
    FOR EACH ROW
    EXECUTE FUNCTION prevent_posted_je_modification();

COMMENT ON FUNCTION prevent_posted_je_modification IS 'Prevents modification of posted journal entries except status changes and notes';

-- ============================================================================
-- JOURNAL ENTRY LINES IMMUTABILITY (for posted parent JE)
-- Description: Can't modify lines when parent journal entry is posted
-- ============================================================================

CREATE OR REPLACE FUNCTION prevent_posted_jel_modification() RETURNS trigger AS $$
DECLARE
    v_is_posted BOOLEAN;
BEGIN
    -- Check if parent journal entry is posted
    SELECT is_posted INTO v_is_posted
    FROM journal_entries
    WHERE id = COALESCE(NEW.journal_entry_id, OLD.journal_entry_id);

    IF v_is_posted = TRUE THEN
        IF TG_OP = 'UPDATE' THEN
            RAISE EXCEPTION 'Cannot modify journal entry lines for posted journal entry. Create a reversal entry instead. Line ID: %', OLD.id
                USING ERRCODE = '23502';
        ELSIF TG_OP = 'DELETE' THEN
            RAISE EXCEPTION 'Cannot delete journal entry lines for posted journal entry: %. Create a reversal entry instead.', OLD.id
                USING ERRCODE = '23502';
        ELSIF TG_OP = 'INSERT' THEN
            RAISE EXCEPTION 'Cannot add journal entry lines to posted journal entry: %. Create a new journal entry instead.', NEW.journal_entry_id
                USING ERRCODE = '23502';
        END IF;
    END IF;

    -- Allow if not posted
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_posted_jel_insert
    BEFORE INSERT ON journal_entry_lines
    FOR EACH ROW
    EXECUTE FUNCTION prevent_posted_jel_modification();

CREATE TRIGGER trg_prevent_posted_jel_update
    BEFORE UPDATE ON journal_entry_lines
    FOR EACH ROW
    EXECUTE FUNCTION prevent_posted_jel_modification();

CREATE TRIGGER trg_prevent_posted_jel_delete
    BEFORE DELETE ON journal_entry_lines
    FOR EACH ROW
    EXECUTE FUNCTION prevent_posted_jel_modification();

COMMENT ON FUNCTION prevent_posted_jel_modification IS 'Prevents modification of journal entry lines when parent JE is posted';

-- ============================================================================
-- FISCAL PERIOD LOCKING
-- Description: Prevent posting to closed/locked periods
-- ============================================================================

CREATE OR REPLACE FUNCTION check_fiscal_period_lock() RETURNS trigger AS $$
DECLARE
    v_period_status VARCHAR;
    v_fy_status VARCHAR;
BEGIN
    -- Check if posting date falls in a closed/locked period
    SELECT ap.status, fy.status
    INTO v_period_status, v_fy_status
    FROM accounting_periods ap
    INNER JOIN fiscal_years fy ON ap.fiscal_year_id = fy.id
    WHERE ap.organization_id = NEW.organization_id
      AND NEW.posting_date BETWEEN ap.start_date AND ap.end_date
    LIMIT 1;

    -- Block if period is closed or locked
    IF v_period_status IN ('closed', 'locked') THEN
        RAISE EXCEPTION 'Cannot post to closed accounting period. Posting date: %', NEW.posting_date
            USING HINT = 'Contact administrator to reopen the period or use a different posting date',
                  ERRCODE = '23514';
    END IF;

    -- Block if fiscal year is closed
    IF v_fy_status IN ('closed', 'locked') THEN
        RAISE EXCEPTION 'Cannot post to closed fiscal year. Posting date: %', NEW.posting_date
            USING HINT = 'Contact administrator or use a different posting date',
                  ERRCODE = '23514';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_fiscal_period_lock_je
    BEFORE INSERT OR UPDATE ON journal_entries
    FOR EACH ROW
    WHEN (NEW.is_posted = TRUE OR NEW.status = 'posted')
    EXECUTE FUNCTION check_fiscal_period_lock();

COMMENT ON FUNCTION check_fiscal_period_lock IS 'Prevents posting to closed or locked accounting periods';

-- ============================================================================
-- AUDIT LOG FOR REVERSAL ATTEMPTS
-- Description: Log all attempts to modify immutable data
-- ============================================================================

CREATE TABLE IF NOT EXISTS immutability_violations_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    table_name VARCHAR(100) NOT NULL,
    record_id UUID,
    operation VARCHAR(20) NOT NULL,  -- 'UPDATE', 'DELETE', 'INSERT'
    attempted_by UUID REFERENCES users(id),
    attempted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    error_message TEXT,
    blocked_data JSONB,  -- Store what was attempted
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_immutability_violations_log_attempted_at ON immutability_violations_log(attempted_at);
CREATE INDEX idx_immutability_violations_log_table_name ON immutability_violations_log(table_name);

COMMENT ON TABLE immutability_violations_log IS 'Audit log of attempts to modify immutable accounting data';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V009 completed successfully!';
    RAISE NOTICE 'Immutability Triggers created:';
    RAISE NOTICE ' - General Ledger: immutable (no updates/deletes)';
    RAISE NOTICE ' - Journal Entries: posted entries immutable';
    RAISE NOTICE ' - Journal Entry Lines: immutable when parent is posted';
    RAISE NOTICE ' - Fiscal Period Locking enforced';
    RAISE NOTICE ' - immutability_violations_log table created';
    RAISE NOTICE '============================================';
END $$;

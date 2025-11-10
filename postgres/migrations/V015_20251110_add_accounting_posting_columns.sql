-- ============================================================================
-- Migration: V015 - Add Accounting Posting Columns to POS Tables
-- Description: Add lightweight tracking columns to key POS tables for accounting integration
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- SALES TABLE - Add Accounting Posting Columns
-- ============================================================================

ALTER TABLE sales
    ADD COLUMN IF NOT EXISTS posted_to_accounting_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS posted_to_accounting_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS accounting_journal_entry_id UUID,  -- Will reference accounting.journal_entries(id) when accounting is enabled
    ADD COLUMN IF NOT EXISTS accounting_posting_status VARCHAR(30) DEFAULT 'not_posted' CHECK (accounting_posting_status IN (
        'not_posted', 'pending', 'posted', 'failed', 'reversed'
    ));

CREATE INDEX IF NOT EXISTS idx_sales_posted_to_accounting_at ON sales(posted_to_accounting_at) WHERE posted_to_accounting_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sales_accounting_posting_status ON sales(accounting_posting_status) WHERE accounting_posting_status != 'not_posted';

COMMENT ON COLUMN sales.posted_to_accounting_at IS 'Timestamp when this sale was posted to accounting';
COMMENT ON COLUMN sales.posted_to_accounting_by IS 'User who posted this sale to accounting';
COMMENT ON COLUMN sales.accounting_journal_entry_id IS 'Reference to accounting journal entry (FK enforced by application)';
COMMENT ON COLUMN sales.accounting_posting_status IS 'Current accounting posting status';

-- ============================================================================
-- GOODS RECEIPTS TABLE - Add Accounting Posting Columns
-- ============================================================================

ALTER TABLE goods_receipts
    ADD COLUMN IF NOT EXISTS posted_to_accounting_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS posted_to_accounting_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS accounting_journal_entry_id UUID,
    ADD COLUMN IF NOT EXISTS accounting_posting_status VARCHAR(30) DEFAULT 'not_posted' CHECK (accounting_posting_status IN (
        'not_posted', 'pending', 'posted', 'failed', 'reversed'
    ));

CREATE INDEX IF NOT EXISTS idx_goods_receipts_posted_to_accounting_at ON goods_receipts(posted_to_accounting_at) WHERE posted_to_accounting_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_goods_receipts_accounting_posting_status ON goods_receipts(accounting_posting_status) WHERE accounting_posting_status != 'not_posted';

COMMENT ON COLUMN goods_receipts.posted_to_accounting_at IS 'Timestamp when this goods receipt was posted to accounting';
COMMENT ON COLUMN goods_receipts.accounting_journal_entry_id IS 'Reference to accounting journal entry (FK enforced by application)';

-- ============================================================================
-- POS SESSIONS TABLE - Add Accounting Posting Columns
-- ============================================================================

ALTER TABLE pos_sessions
    ADD COLUMN IF NOT EXISTS posted_to_accounting_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS posted_to_accounting_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS accounting_journal_entry_id UUID,
    ADD COLUMN IF NOT EXISTS accounting_posting_status VARCHAR(30) DEFAULT 'not_posted' CHECK (accounting_posting_status IN (
        'not_posted', 'pending', 'posted', 'failed', 'reversed'
    ));

CREATE INDEX IF NOT EXISTS idx_pos_sessions_posted_to_accounting_at ON pos_sessions(posted_to_accounting_at) WHERE posted_to_accounting_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_pos_sessions_accounting_posting_status ON pos_sessions(accounting_posting_status) WHERE accounting_posting_status != 'not_posted';

COMMENT ON COLUMN pos_sessions.posted_to_accounting_at IS 'Timestamp when this POS session was posted to accounting';
COMMENT ON COLUMN pos_sessions.accounting_journal_entry_id IS 'Reference to accounting journal entry (FK enforced by application)';

-- ============================================================================
-- INVENTORY TRANSFERS TABLE - Add Accounting Posting Columns
-- ============================================================================

ALTER TABLE inventory_transfers
    ADD COLUMN IF NOT EXISTS posted_to_accounting_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS posted_to_accounting_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS accounting_journal_entry_id UUID,
    ADD COLUMN IF NOT EXISTS accounting_posting_status VARCHAR(30) DEFAULT 'not_posted' CHECK (accounting_posting_status IN (
        'not_posted', 'pending', 'posted', 'failed', 'reversed'
    ));

CREATE INDEX IF NOT EXISTS idx_inventory_transfers_posted_to_accounting_at ON inventory_transfers(posted_to_accounting_at) WHERE posted_to_accounting_at IS NOT NULL;

COMMENT ON COLUMN inventory_transfers.posted_to_accounting_at IS 'Timestamp when this inventory transfer was posted to accounting';

-- ============================================================================
-- EXPENSES TABLE - Add Accounting Posting Columns (if not exists)
-- ============================================================================

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'expenses') THEN
        ALTER TABLE expenses
            ADD COLUMN IF NOT EXISTS posted_to_accounting_at TIMESTAMP WITH TIME ZONE,
            ADD COLUMN IF NOT EXISTS posted_to_accounting_by UUID REFERENCES users(id),
            ADD COLUMN IF NOT EXISTS accounting_journal_entry_id UUID,
            ADD COLUMN IF NOT EXISTS accounting_posting_status VARCHAR(30) DEFAULT 'not_posted' CHECK (accounting_posting_status IN (
                'not_posted', 'pending', 'posted', 'failed', 'reversed'
            ));

        CREATE INDEX IF NOT EXISTS idx_expenses_posted_to_accounting_at ON expenses(posted_to_accounting_at) WHERE posted_to_accounting_at IS NOT NULL;
    END IF;
END $$;

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V015 completed successfully!';
    RAISE NOTICE 'Added accounting posting columns to:';
    RAISE NOTICE ' - sales';
    RAISE NOTICE ' - goods_receipts';
    RAISE NOTICE ' - pos_sessions';
    RAISE NOTICE ' - inventory_transfers';
    RAISE NOTICE ' - expenses (if exists)';
    RAISE NOTICE '============================================';
END $$;

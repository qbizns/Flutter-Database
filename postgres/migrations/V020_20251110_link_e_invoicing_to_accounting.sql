-- ============================================================================
-- Migration: V020 - Link E-Invoicing to Accounting
-- Description: Add references from e-invoicing documents to accounting records
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- E-INVOICING DOCUMENTS TABLE - Add Accounting References
-- ============================================================================

ALTER TABLE e_invoicing_documents
    ADD COLUMN IF NOT EXISTS accounting_invoice_id UUID,  -- References accounting.customer_invoices(id) or accounting.vendor_bills(id)
    ADD COLUMN IF NOT EXISTS accounting_invoice_type VARCHAR(50), -- 'customer_invoice', 'vendor_bill', 'credit_note', 'debit_note'
    ADD COLUMN IF NOT EXISTS accounting_journal_entry_id UUID;    -- References accounting.journal_entries(id)

CREATE INDEX IF NOT EXISTS idx_e_invoicing_documents_accounting_invoice_id ON e_invoicing_documents(accounting_invoice_id) WHERE accounting_invoice_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_e_invoicing_documents_accounting_je_id ON e_invoicing_documents(accounting_journal_entry_id) WHERE accounting_journal_entry_id IS NOT NULL;

COMMENT ON COLUMN e_invoicing_documents.accounting_invoice_id IS 'Reference to accounting invoice/bill (FK enforced by application when accounting module enabled)';
COMMENT ON COLUMN e_invoicing_documents.accounting_invoice_type IS 'Type of accounting document (customer_invoice, vendor_bill, etc.)';
COMMENT ON COLUMN e_invoicing_documents.accounting_journal_entry_id IS 'Journal entry that posted this invoice to accounting';

-- ============================================================================
-- E-INVOICING DOCUMENTS TABLE - Add Financial Amounts Cache
-- ============================================================================

ALTER TABLE e_invoicing_documents
    ADD COLUMN IF NOT EXISTS invoice_amount NUMERIC(20, 4),      -- Total invoice amount before tax
    ADD COLUMN IF NOT EXISTS tax_amount NUMERIC(20, 4),          -- Total tax amount
    ADD COLUMN IF NOT EXISTS total_amount NUMERIC(20, 4),        -- Grand total including tax
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3);            -- Invoice currency

CREATE INDEX IF NOT EXISTS idx_e_invoicing_documents_total_amount ON e_invoicing_documents(total_amount) WHERE total_amount IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_e_invoicing_documents_currency_code ON e_invoicing_documents(currency_code) WHERE currency_code IS NOT NULL;

COMMENT ON COLUMN e_invoicing_documents.invoice_amount IS 'Invoice subtotal before tax (cached for reporting)';
COMMENT ON COLUMN e_invoicing_documents.tax_amount IS 'Total tax amount (cached for reporting)';
COMMENT ON COLUMN e_invoicing_documents.total_amount IS 'Grand total including tax (cached for reporting)';
COMMENT ON COLUMN e_invoicing_documents.currency_code IS 'Invoice currency code';

-- ============================================================================
-- HELPER VIEW: E-Invoicing Documents with Accounting Links
-- ============================================================================

CREATE OR REPLACE VIEW view_e_invoicing_accounting_links AS
SELECT
    ed.id as e_invoice_id,
    ed.organization_id,
    ed.authority,
    ed.document_number,
    ed.document_type,
    ed.status,
    ed.source_table,
    ed.source_id,
    ed.accounting_invoice_id,
    ed.accounting_invoice_type,
    ed.accounting_journal_entry_id,
    ed.invoice_amount,
    ed.tax_amount,
    ed.total_amount,
    ed.currency_code,
    ed.created_at,
    ed.submitted_at,
    ed.response_at,
    -- Sale details (if source is sales)
    CASE WHEN ed.source_table = 'sales' THEN s.sale_number ELSE NULL END as sale_number,
    CASE WHEN ed.source_table = 'sales' THEN s.total_amount ELSE NULL END as sale_total,
    CASE WHEN ed.source_table = 'sales' THEN s.accounting_posting_status ELSE NULL END as sale_posting_status
FROM e_invoicing_documents ed
LEFT JOIN sales s ON ed.source_table = 'sales' AND ed.source_id = s.id
WHERE ed.deleted_at IS NULL;

COMMENT ON VIEW view_e_invoicing_accounting_links IS 'E-invoicing documents with accounting and source document links';

-- ============================================================================
-- HELPER VIEW: E-Invoicing Revenue Recognition Status
-- ============================================================================

CREATE OR REPLACE VIEW view_e_invoice_revenue_recognition AS
SELECT
    ed.id as e_invoice_id,
    ed.organization_id,
    ed.document_number,
    ed.authority,
    ed.status as e_invoice_status,
    ed.total_amount,
    ed.currency_code,
    ed.accounting_invoice_id,
    ed.accounting_journal_entry_id,
    -- Check if posted to accounting
    CASE
        WHEN ed.accounting_journal_entry_id IS NOT NULL THEN 'posted'
        WHEN ed.accounting_invoice_id IS NOT NULL THEN 'invoice_created'
        ELSE 'not_posted'
    END as accounting_status,
    -- Check if e-invoice is accepted
    CASE
        WHEN ed.status = 'accepted' THEN TRUE
        ELSE FALSE
    END as is_e_invoice_accepted,
    -- Overall revenue recognition status
    CASE
        WHEN ed.status = 'accepted' AND ed.accounting_journal_entry_id IS NOT NULL THEN 'recognized'
        WHEN ed.status = 'accepted' AND ed.accounting_invoice_id IS NOT NULL THEN 'pending_posting'
        WHEN ed.status = 'accepted' THEN 'pending_accounting'
        WHEN ed.status IN ('submitted', 'pending') THEN 'pending_acceptance'
        WHEN ed.status = 'rejected' THEN 'blocked'
        ELSE 'draft'
    END as revenue_recognition_status,
    ed.submitted_at,
    ed.response_at
FROM e_invoicing_documents ed
WHERE ed.deleted_at IS NULL
  AND ed.document_type NOT IN ('credit_note', 'debit_note')  -- Focus on revenue documents
ORDER BY ed.created_at DESC;

COMMENT ON VIEW view_e_invoice_revenue_recognition IS 'Revenue recognition status based on e-invoicing and accounting linkage';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V020 completed successfully!';
    RAISE NOTICE 'Linked E-Invoicing to Accounting:';
    RAISE NOTICE ' - Added accounting_invoice_id to e_invoicing_documents';
    RAISE NOTICE ' - Added accounting_journal_entry_id to e_invoicing_documents';
    RAISE NOTICE ' - Added financial amounts cache';
    RAISE NOTICE ' - Created view_e_invoicing_accounting_links';
    RAISE NOTICE ' - Created view_e_invoice_revenue_recognition';
    RAISE NOTICE '============================================';
END $$;

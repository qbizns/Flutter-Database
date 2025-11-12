-- ============================================================================
-- Migration: V024 - Add Performance Indexes for Dashboard and Reporting
-- Description: Critical performance indexes for high-volume queries
-- Author: Phase 1 Security & Stability
-- Date: 2025-11-12
-- ============================================================================

BEGIN;

-- ============================================================================
-- GENERAL LEDGER PERFORMANCE INDEXES
-- ============================================================================

-- Index for date range queries with organization
CREATE INDEX IF NOT EXISTS idx_gl_org_date_range
    ON general_ledger(organization_id, transaction_date, account_id)
    WHERE organization_id IS NOT NULL;

-- Index for recent transactions (dashboard)
CREATE INDEX IF NOT EXISTS idx_gl_posting_date_desc
    ON general_ledger(posting_date DESC, organization_id)
    WHERE organization_id IS NOT NULL;

-- Index for account balance calculations
CREATE INDEX IF NOT EXISTS idx_gl_account_period_balance
    ON general_ledger(account_id, accounting_period_id, created_at)
    WHERE accounting_period_id IS NOT NULL;

-- ============================================================================
-- JOURNAL ENTRIES PERFORMANCE INDEXES
-- ============================================================================

-- Index for organization timeline view
CREATE INDEX IF NOT EXISTS idx_je_org_created
    ON journal_entries(organization_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- Index for posted entries by date
CREATE INDEX IF NOT EXISTS idx_je_org_entry_date
    ON journal_entries(organization_id, entry_date DESC)
    WHERE deleted_at IS NULL AND is_posted = true;

-- Index for pending approvals
CREATE INDEX IF NOT EXISTS idx_je_requires_approval
    ON journal_entries(organization_id, requires_approval, status)
    WHERE deleted_at IS NULL AND requires_approval = true;

-- ============================================================================
-- SALES/POS PERFORMANCE INDEXES
-- ============================================================================

-- Index for sales dashboard
CREATE INDEX IF NOT EXISTS idx_sales_org_date
    ON sales(organization_id, sale_date DESC)
    WHERE deleted_at IS NULL;

-- Index for sales by status and date
CREATE INDEX IF NOT EXISTS idx_sales_status_date
    ON sales(organization_id, status, sale_date DESC)
    WHERE deleted_at IS NULL;

-- Index for cashier performance reports
CREATE INDEX IF NOT EXISTS idx_sales_cashier_date
    ON sales(cashier_id, sale_date)
    WHERE deleted_at IS NULL AND cashier_id IS NOT NULL;

-- Index for customer purchase history
CREATE INDEX IF NOT EXISTS idx_sales_customer_date
    ON sales(customer_id, sale_date DESC)
    WHERE deleted_at IS NULL AND customer_id IS NOT NULL;

-- ============================================================================
-- ACCOUNTS RECEIVABLE (AR) INDEXES
-- ============================================================================

-- Index for overdue invoices
CREATE INDEX IF NOT EXISTS idx_customer_invoices_due
    ON customer_invoices(organization_id, due_date, status)
    WHERE deleted_at IS NULL AND status IN ('unpaid', 'partial', 'overdue');

-- Index for customer statement generation
CREATE INDEX IF NOT EXISTS idx_customer_invoices_customer_date
    ON customer_invoices(customer_id, invoice_date DESC)
    WHERE deleted_at IS NULL;

-- Index for aging reports
CREATE INDEX IF NOT EXISTS idx_customer_invoices_aging
    ON customer_invoices(organization_id, due_date, balance_due)
    WHERE deleted_at IS NULL AND balance_due > 0;

-- ============================================================================
-- ACCOUNTS PAYABLE (AP) INDEXES
-- ============================================================================

-- Index for outstanding bills
CREATE INDEX IF NOT EXISTS idx_vendor_bills_due
    ON vendor_bills(organization_id, due_date, status)
    WHERE deleted_at IS NULL AND status IN ('unpaid', 'partial', 'overdue');

-- Index for vendor history
CREATE INDEX IF NOT EXISTS idx_vendor_bills_supplier_date
    ON vendor_bills(supplier_id, bill_date DESC)
    WHERE deleted_at IS NULL;

-- Index for payment planning
CREATE INDEX IF NOT EXISTS idx_vendor_bills_aging
    ON vendor_bills(organization_id, due_date, balance_due)
    WHERE deleted_at IS NULL AND balance_due > 0;

-- ============================================================================
-- PRODUCTS & INVENTORY INDEXES
-- ============================================================================

-- Index for active products
CREATE INDEX IF NOT EXISTS idx_products_org_active
    ON products(organization_id, is_active, name)
    WHERE deleted_at IS NULL;

-- Index for low stock alerts
CREATE INDEX IF NOT EXISTS idx_products_low_stock
    ON products(organization_id, quantity_on_hand)
    WHERE deleted_at IS NULL AND is_active = true AND track_inventory = true;

-- Index for product search by SKU
CREATE INDEX IF NOT EXISTS idx_products_sku_search
    ON products(organization_id, sku)
    WHERE deleted_at IS NULL AND sku IS NOT NULL;

-- Index for category-based browsing
CREATE INDEX IF NOT EXISTS idx_products_category
    ON products(category_id, is_active, name)
    WHERE deleted_at IS NULL;

-- ============================================================================
-- CUSTOMERS & SUPPLIERS INDEXES
-- ============================================================================

-- Index for customer search
CREATE INDEX IF NOT EXISTS idx_customers_search
    ON customers(organization_id, name, email)
    WHERE deleted_at IS NULL;

-- Index for supplier search
CREATE INDEX IF NOT EXISTS idx_suppliers_search
    ON suppliers(organization_id, name, email)
    WHERE deleted_at IS NULL;

-- ============================================================================
-- POSTING ENGINE PERFORMANCE
-- ============================================================================

-- Index for posting audit trail
CREATE INDEX IF NOT EXISTS idx_pos_posting_audit_source
    ON pos_posting_audit(organization_id, source_table, source_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- Index for failed postings
CREATE INDEX IF NOT EXISTS idx_pos_posting_audit_failures
    ON pos_posting_audit(organization_id, posting_status, created_at DESC)
    WHERE deleted_at IS NULL AND posting_status IN ('failed', 'validation_blocked');

-- ============================================================================
-- CHART OF ACCOUNTS OPTIMIZATION
-- ============================================================================

-- Index for account balance updates
CREATE INDEX IF NOT EXISTS idx_coa_balance_update
    ON chart_of_accounts(organization_id, id, current_balance, last_balance_update)
    WHERE deleted_at IS NULL AND is_active = true;

-- Index for reconcilable accounts
CREATE INDEX IF NOT EXISTS idx_coa_reconcilable
    ON chart_of_accounts(organization_id, is_reconcilable, account_type_id)
    WHERE deleted_at IS NULL AND is_reconcilable = true;

-- ============================================================================
-- BANK RECONCILIATION INDEXES
-- ============================================================================

-- Index for pending reconciliations
CREATE INDEX IF NOT EXISTS idx_bank_recon_pending
    ON bank_reconciliations(organization_id, status, statement_date DESC)
    WHERE deleted_at IS NULL AND status != 'reconciled';

-- ============================================================================
-- COMPOSITE INDEXES FOR COMPLEX QUERIES
-- ============================================================================

-- Sales with payment details
CREATE INDEX IF NOT EXISTS idx_sales_payment_composite
    ON sales(organization_id, payment_status, payment_method, sale_date DESC)
    WHERE deleted_at IS NULL;

-- Journal entries for audit trail
CREATE INDEX IF NOT EXISTS idx_je_audit_composite
    ON journal_entries(organization_id, source_module, source_document_type, entry_date DESC)
    WHERE deleted_at IS NULL;

COMMIT;

-- ============================================================================
-- ANALYZE TABLES FOR QUERY PLANNER
-- ============================================================================

ANALYZE general_ledger;
ANALYZE journal_entries;
ANALYZE journal_entry_lines;
ANALYZE sales;
ANALYZE customer_invoices;
ANALYZE vendor_bills;
ANALYZE products;
ANALYZE chart_of_accounts;

-- ============================================================================
-- SUCCESS MESSAGE
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V024 completed successfully!';
    RAISE NOTICE 'Performance indexes added:';
    RAISE NOTICE ' - General Ledger: 3 indexes';
    RAISE NOTICE ' - Journal Entries: 3 indexes';
    RAISE NOTICE ' - Sales/POS: 4 indexes';
    RAISE NOTICE ' - Accounts Receivable: 3 indexes';
    RAISE NOTICE ' - Accounts Payable: 3 indexes';
    RAISE NOTICE ' - Products: 4 indexes';
    RAISE NOTICE ' - Customers/Suppliers: 2 indexes';
    RAISE NOTICE ' - Posting Engine: 2 indexes';
    RAISE NOTICE ' - Chart of Accounts: 2 indexes';
    RAISE NOTICE ' - Bank Reconciliation: 1 index';
    RAISE NOTICE ' - Composite indexes: 2 indexes';
    RAISE NOTICE 'Total: 29 performance indexes';
    RAISE NOTICE '============================================';
END $$;

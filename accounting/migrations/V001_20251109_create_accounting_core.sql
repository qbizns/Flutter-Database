-- =====================================================
-- Accounting Module - V001: Core Accounting Tables
-- Description: Chart of Accounts, General Ledger, Journal Entries
-- Based on: QuickBooks, Xero, Sage, NetSuite best practices
-- Date: 2025-11-09
-- Dependencies: POS core tables (organizations, users)
-- =====================================================

-- =====================================================
-- SECTION 1: Fiscal Years & Accounting Periods
-- =====================================================

CREATE TABLE fiscal_years (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Year Details
    fiscal_year VARCHAR(10) NOT NULL, -- e.g., "2024", "FY2024"
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    -- Status
    status VARCHAR(20) DEFAULT 'open', -- open, closed, locked
    is_current BOOLEAN DEFAULT false,

    -- Closing
    closed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    closed_at TIMESTAMP WITH TIME ZONE,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_fiscal_years_status CHECK (status IN ('open', 'closed', 'locked')),
    CONSTRAINT chk_fiscal_years_dates CHECK (end_date > start_date),
    UNIQUE(organization_id, fiscal_year)
);

CREATE INDEX idx_fiscal_years_org ON fiscal_years(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_fiscal_years_current ON fiscal_years(organization_id, is_current) WHERE deleted_at IS NULL AND is_current = true;
CREATE INDEX idx_fiscal_years_status ON fiscal_years(status) WHERE deleted_at IS NULL;

CREATE TRIGGER update_fiscal_years_updated_at
    BEFORE UPDATE ON fiscal_years
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE accounting_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    fiscal_year_id UUID NOT NULL REFERENCES fiscal_years(id) ON DELETE CASCADE,

    -- Period Details
    period_number INTEGER NOT NULL, -- 1-12 for monthly, 1-4 for quarterly
    period_name VARCHAR(50) NOT NULL, -- "January 2024", "Q1 2024"
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    -- Status
    status VARCHAR(20) DEFAULT 'open', -- open, closed, locked

    -- Closing
    closed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    closed_at TIMESTAMP WITH TIME ZONE,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_accounting_periods_status CHECK (status IN ('open', 'closed', 'locked')),
    CONSTRAINT chk_accounting_periods_dates CHECK (end_date > start_date),
    UNIQUE(organization_id, fiscal_year_id, period_number)
);

CREATE INDEX idx_accounting_periods_org ON accounting_periods(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_accounting_periods_fiscal_year ON accounting_periods(fiscal_year_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_accounting_periods_status ON accounting_periods(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_accounting_periods_dates ON accounting_periods(start_date, end_date) WHERE deleted_at IS NULL;

CREATE TRIGGER update_accounting_periods_updated_at
    BEFORE UPDATE ON accounting_periods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 2: Chart of Accounts
-- =====================================================

CREATE TABLE account_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Type Details
    type_code VARCHAR(20) NOT NULL UNIQUE, -- ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
    type_name VARCHAR(100) NOT NULL,
    type_category VARCHAR(50) NOT NULL, -- balance_sheet, income_statement
    normal_balance VARCHAR(10) NOT NULL, -- debit, credit

    -- Classification
    is_balance_sheet BOOLEAN DEFAULT false,
    is_income_statement BOOLEAN DEFAULT false,

    -- Display
    display_order INTEGER DEFAULT 0,

    -- Description
    description TEXT,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT chk_account_types_category CHECK (type_category IN ('balance_sheet', 'income_statement')),
    CONSTRAINT chk_account_types_balance CHECK (normal_balance IN ('debit', 'credit'))
);

CREATE INDEX idx_account_types_code ON account_types(type_code);
CREATE INDEX idx_account_types_category ON account_types(type_category);

-- =====================================================

CREATE TABLE account_subtypes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_type_id UUID NOT NULL REFERENCES account_types(id) ON DELETE CASCADE,

    -- Subtype Details
    subtype_code VARCHAR(50) NOT NULL UNIQUE,
    subtype_name VARCHAR(100) NOT NULL,

    -- Display
    display_order INTEGER DEFAULT 0,

    -- Description
    description TEXT,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(account_type_id, subtype_code)
);

CREATE INDEX idx_account_subtypes_type ON account_subtypes(account_type_id);
CREATE INDEX idx_account_subtypes_code ON account_subtypes(subtype_code);

-- =====================================================

CREATE TABLE chart_of_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Account Identification
    account_code VARCHAR(50) NOT NULL, -- e.g., "1000", "1100-01"
    account_number VARCHAR(50) NOT NULL, -- Full account number
    account_name VARCHAR(255) NOT NULL,

    -- Classification
    account_type_id UUID NOT NULL REFERENCES account_types(id) ON DELETE RESTRICT,
    account_subtype_id UUID REFERENCES account_subtypes(id) ON DELETE SET NULL,

    -- Hierarchy (for sub-accounts)
    parent_account_id UUID REFERENCES chart_of_accounts(id) ON DELETE SET NULL,
    account_level INTEGER DEFAULT 1,
    account_path TEXT, -- e.g., "1000/1100/1110"

    -- Properties
    is_active BOOLEAN DEFAULT true,
    is_system_account BOOLEAN DEFAULT false, -- Cannot be deleted
    is_header_account BOOLEAN DEFAULT false, -- Group/parent account, no posting
    is_bank_account BOOLEAN DEFAULT false,
    is_reconcilable BOOLEAN DEFAULT false,

    -- Default Tax
    default_tax_code VARCHAR(50),

    -- Currency (for multi-currency support)
    currency_code VARCHAR(3) DEFAULT 'USD',

    -- Opening Balance
    opening_balance NUMERIC(20, 4) DEFAULT 0,
    opening_balance_date DATE,

    -- Current Balance (cached, recalculated periodically)
    current_debit_balance NUMERIC(20, 4) DEFAULT 0,
    current_credit_balance NUMERIC(20, 4) DEFAULT 0,
    current_balance NUMERIC(20, 4) DEFAULT 0,
    last_balance_update TIMESTAMP WITH TIME ZONE,

    -- Description
    description TEXT,
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_coa_account_level CHECK (account_level > 0 AND account_level <= 5),
    UNIQUE(organization_id, account_code)
);

CREATE INDEX idx_coa_org ON chart_of_accounts(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_coa_account_code ON chart_of_accounts(organization_id, account_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_coa_account_number ON chart_of_accounts(account_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_coa_type ON chart_of_accounts(account_type_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_coa_parent ON chart_of_accounts(parent_account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_coa_active ON chart_of_accounts(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_coa_bank ON chart_of_accounts(is_bank_account) WHERE deleted_at IS NULL AND is_bank_account = true;

CREATE TRIGGER update_chart_of_accounts_updated_at
    BEFORE UPDATE ON chart_of_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 3: Journal Entries (Double-Entry Accounting)
-- =====================================================

CREATE TABLE journal_entry_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Type Details
    type_code VARCHAR(50) NOT NULL UNIQUE,
    type_name VARCHAR(100) NOT NULL,
    type_category VARCHAR(50), -- standard, adjusting, closing, reversing

    -- Prefix for numbering
    number_prefix VARCHAR(10), -- JE, AJE, CJE, RJE

    -- Description
    description TEXT,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_je_types_code ON journal_entry_types(type_code);

-- =====================================================

CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Entry Identification
    entry_number VARCHAR(50) NOT NULL,
    entry_type_id UUID NOT NULL REFERENCES journal_entry_types(id) ON DELETE RESTRICT,

    -- Dates
    entry_date DATE NOT NULL,
    posting_date DATE NOT NULL,

    -- Period
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,
    fiscal_year_id UUID REFERENCES fiscal_years(id) ON DELETE SET NULL,

    -- Status
    status VARCHAR(20) DEFAULT 'draft', -- draft, posted, approved, reversed, voided
    is_posted BOOLEAN DEFAULT false,
    is_reversed BOOLEAN DEFAULT false,
    reversal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,

    -- Source
    source_module VARCHAR(50), -- manual, sales, purchases, payroll, bank, etc.
    source_document_type VARCHAR(50), -- invoice, payment, bill, etc.
    source_document_id UUID,
    reference_number VARCHAR(100),

    -- Totals (must balance)
    total_debit NUMERIC(20, 4) DEFAULT 0,
    total_credit NUMERIC(20, 4) DEFAULT 0,

    -- Description
    description TEXT NOT NULL,
    notes TEXT,

    -- Approval
    requires_approval BOOLEAN DEFAULT false,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP WITH TIME ZONE,

    -- Posting
    posted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    posted_at TIMESTAMP WITH TIME ZONE,

    -- Attachments
    attachments JSONB DEFAULT '[]',

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_je_status CHECK (status IN ('draft', 'posted', 'approved', 'reversed', 'voided')),
    CONSTRAINT chk_je_balanced CHECK (
        (is_posted = false) OR
        (ABS(total_debit - total_credit) < 0.01)
    ),
    UNIQUE(organization_id, entry_number)
);

CREATE INDEX idx_je_org ON journal_entries(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_je_entry_number ON journal_entries(organization_id, entry_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_je_entry_date ON journal_entries(entry_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_je_posting_date ON journal_entries(posting_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_je_period ON journal_entries(accounting_period_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_je_status ON journal_entries(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_je_posted ON journal_entries(is_posted) WHERE deleted_at IS NULL;
CREATE INDEX idx_je_source ON journal_entries(source_module, source_document_type) WHERE deleted_at IS NULL;

CREATE TRIGGER update_journal_entries_updated_at
    BEFORE UPDATE ON journal_entries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE journal_entry_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Entry Reference
    journal_entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,

    -- Account
    account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,

    -- Debit/Credit
    debit_amount NUMERIC(20, 4) DEFAULT 0,
    credit_amount NUMERIC(20, 4) DEFAULT 0,

    -- Dimensions (for detailed reporting)
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    department VARCHAR(100),
    project_code VARCHAR(100),
    cost_center VARCHAR(100),

    -- Tax
    tax_code VARCHAR(50),
    tax_amount NUMERIC(20, 4) DEFAULT 0,

    -- Description
    description TEXT,
    memo TEXT,

    -- Reconciliation
    is_reconciled BOOLEAN DEFAULT false,
    reconciled_at TIMESTAMP WITH TIME ZONE,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_jel_debit_credit CHECK (
        (debit_amount > 0 AND credit_amount = 0) OR
        (credit_amount > 0 AND debit_amount = 0) OR
        (debit_amount = 0 AND credit_amount = 0)
    ),
    CONSTRAINT chk_jel_amounts CHECK (debit_amount >= 0 AND credit_amount >= 0),
    UNIQUE(journal_entry_id, line_number)
);

CREATE INDEX idx_jel_org ON journal_entry_lines(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_jel_entry ON journal_entry_lines(journal_entry_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_jel_account ON journal_entry_lines(account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_jel_location ON journal_entry_lines(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_jel_reconciled ON journal_entry_lines(is_reconciled) WHERE deleted_at IS NULL;

CREATE TRIGGER update_journal_entry_lines_updated_at
    BEFORE UPDATE ON journal_entry_lines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 4: General Ledger (Posted Transactions)
-- =====================================================

CREATE TABLE general_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Source
    journal_entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    journal_entry_line_id UUID NOT NULL REFERENCES journal_entry_lines(id) ON DELETE CASCADE,

    -- Account
    account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,

    -- Dates
    transaction_date DATE NOT NULL,
    posting_date DATE NOT NULL,

    -- Period
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,
    fiscal_year_id UUID REFERENCES fiscal_years(id) ON DELETE SET NULL,

    -- Amounts
    debit_amount NUMERIC(20, 4) DEFAULT 0,
    credit_amount NUMERIC(20, 4) DEFAULT 0,

    -- Running Balance (for account)
    running_debit_balance NUMERIC(20, 4) DEFAULT 0,
    running_credit_balance NUMERIC(20, 4) DEFAULT 0,
    running_balance NUMERIC(20, 4) DEFAULT 0,

    -- Source Reference
    source_module VARCHAR(50),
    source_document_type VARCHAR(50),
    source_document_id UUID,
    reference_number VARCHAR(100),

    -- Dimensions
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    department VARCHAR(100),
    project_code VARCHAR(100),
    cost_center VARCHAR(100),

    -- Description
    description TEXT,

    -- Status
    is_reversed BOOLEAN DEFAULT false,
    reversal_gl_id UUID REFERENCES general_ledger(id) ON DELETE SET NULL,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit (no updates allowed once posted)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,

    CONSTRAINT chk_gl_amounts CHECK (debit_amount >= 0 AND credit_amount >= 0),
    CONSTRAINT chk_gl_debit_credit CHECK (
        (debit_amount > 0 AND credit_amount = 0) OR
        (credit_amount > 0 AND debit_amount = 0)
    )
);

-- Indexes for performance (GL can have millions of records)
CREATE INDEX idx_gl_org ON general_ledger(organization_id);
CREATE INDEX idx_gl_account ON general_ledger(account_id);
CREATE INDEX idx_gl_transaction_date ON general_ledger(transaction_date);
CREATE INDEX idx_gl_posting_date ON general_ledger(posting_date);
CREATE INDEX idx_gl_period ON general_ledger(accounting_period_id);
CREATE INDEX idx_gl_fiscal_year ON general_ledger(fiscal_year_id);
CREATE INDEX idx_gl_je ON general_ledger(journal_entry_id);
CREATE INDEX idx_gl_location ON general_ledger(location_id);
CREATE INDEX idx_gl_source ON general_ledger(source_module, source_document_type);

-- Composite index for common queries
CREATE INDEX idx_gl_account_date ON general_ledger(account_id, transaction_date);
CREATE INDEX idx_gl_account_period ON general_ledger(account_id, accounting_period_id);

-- =====================================================
-- SECTION 5: Row-Level Security (RLS)
-- =====================================================

ALTER TABLE fiscal_years ENABLE ROW LEVEL SECURITY;
ALTER TABLE accounting_periods ENABLE ROW LEVEL SECURITY;
ALTER TABLE chart_of_accounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE journal_entries ENABLE ROW LEVEL SECURITY;
ALTER TABLE journal_entry_lines ENABLE ROW LEVEL SECURITY;
ALTER TABLE general_ledger ENABLE ROW LEVEL SECURITY;

-- Super admin policies
CREATE POLICY fiscal_years_super_admin ON fiscal_years FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id
           WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin'));

CREATE POLICY accounting_periods_super_admin ON accounting_periods FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id
           WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin'));

CREATE POLICY coa_super_admin ON chart_of_accounts FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id
           WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin'));

CREATE POLICY je_super_admin ON journal_entries FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id
           WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin'));

CREATE POLICY jel_super_admin ON journal_entry_lines FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id
           WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin'));

CREATE POLICY gl_super_admin ON general_ledger FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id
           WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin'));

-- Organization-scoped policies
CREATE POLICY fiscal_years_org ON fiscal_years FOR SELECT TO PUBLIC
    USING (organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY accounting_periods_org ON accounting_periods FOR SELECT TO PUBLIC
    USING (organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY coa_org ON chart_of_accounts FOR SELECT TO PUBLIC
    USING (organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY je_org ON journal_entries FOR SELECT TO PUBLIC
    USING (organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY jel_org ON journal_entry_lines FOR SELECT TO PUBLIC
    USING (organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY gl_org ON general_ledger FOR SELECT TO PUBLIC
    USING (organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

-- =====================================================
-- End of Migration V001
-- =====================================================

-- Summary: Created core accounting infrastructure
-- - Fiscal years and accounting periods
-- - Chart of accounts with hierarchy
-- - Journal entries with double-entry validation
-- - General ledger for posted transactions
-- - Complete RLS policies
-- - Optimized indexes for performance

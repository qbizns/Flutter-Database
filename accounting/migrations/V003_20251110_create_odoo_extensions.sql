-- =====================================================
-- Accounting Migration V003: Odoo-Style Extensions
-- Description: Advanced accounting features (Journals, Taxes, Multi-currency,
--              Payment Terms, Analytics, Deferrals, Bank Statements, Budgets, Localization)
-- =====================================================

\echo 'Creating Odoo-style accounting extensions...';

-- Set search path
SET search_path TO accounting, public;

-- =====================================================
-- MODULE 1: JOURNALS (Proper Journal Layer)
-- =====================================================

\echo 'Creating journals module...';

-- Journals Table
CREATE TABLE IF NOT EXISTS journals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Journal identification
    journal_code VARCHAR(10) NOT NULL,
    journal_name VARCHAR(255) NOT NULL,
    journal_type VARCHAR(30) NOT NULL CHECK (journal_type IN (
        'sale', 'purchase', 'bank', 'cash', 'general', 'miscellaneous'
    )),

    -- Configuration
    bank_account_id UUID REFERENCES bank_accounts(id),
    default_debit_account_id UUID REFERENCES chart_of_accounts(id),
    default_credit_account_id UUID REFERENCES chart_of_accounts(id),

    -- Sequencing for entry numbers
    sequence_prefix VARCHAR(20),
    sequence_number INTEGER DEFAULT 1,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT uq_journal_code_org UNIQUE(organization_id, journal_code),
    CONSTRAINT chk_bank_journal_account CHECK (
        (journal_type != 'bank') OR (bank_account_id IS NOT NULL)
    )
);

CREATE INDEX idx_journals_org_type ON journals(organization_id, journal_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_journals_active ON journals(organization_id, is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE journals IS 'Accounting journals (Sales, Purchase, Bank, Cash, General, Misc)';
COMMENT ON COLUMN journals.journal_type IS 'Controls posting logic and UI behavior';

-- Enable RLS
ALTER TABLE journals ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_journals ON journals FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Add journal_id to journal_entries (extend existing table)
ALTER TABLE journal_entries ADD COLUMN IF NOT EXISTS journal_id UUID REFERENCES journals(id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_journal ON journal_entries(journal_id) WHERE deleted_at IS NULL;

-- =====================================================
-- MODULE 2: TAX ENGINE
-- =====================================================

\echo 'Creating tax engine module...';

-- Tax Groups
CREATE TABLE IF NOT EXISTS tax_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Identification
    group_code VARCHAR(20) NOT NULL,
    group_name VARCHAR(255) NOT NULL,

    -- Configuration
    sequence INTEGER DEFAULT 10,
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_tax_group_code_org UNIQUE(organization_id, group_code)
);

CREATE INDEX idx_tax_groups_org ON tax_groups(organization_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE tax_groups IS 'Logical grouping of taxes for reporting (e.g., VAT 20%, Withholding)';

-- Enable RLS
ALTER TABLE tax_groups ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_tax_groups ON tax_groups FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Taxes
CREATE TABLE IF NOT EXISTS taxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    tax_group_id UUID REFERENCES tax_groups(id),

    -- Identification
    tax_code VARCHAR(20) NOT NULL,
    tax_name VARCHAR(255) NOT NULL,

    -- Tax configuration
    tax_rate NUMERIC(10, 4) NOT NULL, -- Percentage (e.g., 20.0000 for 20%)
    tax_scope VARCHAR(20) NOT NULL CHECK (tax_scope IN ('sales', 'purchases', 'both')),
    is_price_inclusive BOOLEAN DEFAULT false,

    -- Accounts for posting
    tax_account_id UUID NOT NULL REFERENCES chart_of_accounts(id), -- Where tax is collected/paid
    tax_refund_account_id UUID REFERENCES chart_of_accounts(id), -- For refunds

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    description TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_tax_code_org UNIQUE(organization_id, tax_code)
);

CREATE INDEX idx_taxes_org_scope ON taxes(organization_id, tax_scope) WHERE deleted_at IS NULL;
CREATE INDEX idx_taxes_group ON taxes(tax_group_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE taxes IS 'Tax definitions with rates, scope, and posting accounts';
COMMENT ON COLUMN taxes.is_price_inclusive IS 'Whether price includes tax or tax is added on top';

-- Enable RLS
ALTER TABLE taxes ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_taxes ON taxes FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Fiscal Positions (Tax remapping for regions/customers)
CREATE TABLE IF NOT EXISTS fiscal_positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Identification
    position_code VARCHAR(20) NOT NULL,
    position_name VARCHAR(255) NOT NULL,

    -- Configuration
    auto_apply BOOLEAN DEFAULT false,
    country_id VARCHAR(2), -- ISO country code
    state_province VARCHAR(100),
    zip_postal_code_range VARCHAR(100),

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_fiscal_position_code_org UNIQUE(organization_id, position_code)
);

CREATE INDEX idx_fiscal_positions_org ON fiscal_positions(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_fiscal_positions_country ON fiscal_positions(country_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE fiscal_positions IS 'Tax rules for specific regions or customer types (e.g., EU B2B, Export)';

-- Enable RLS
ALTER TABLE fiscal_positions ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_fiscal_positions ON fiscal_positions FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Fiscal Position Tax Mappings
CREATE TABLE IF NOT EXISTS fiscal_position_tax_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fiscal_position_id UUID NOT NULL REFERENCES fiscal_positions(id) ON DELETE CASCADE,

    -- Tax mapping: from_tax → to_tax
    source_tax_id UUID NOT NULL REFERENCES taxes(id),
    destination_tax_id UUID REFERENCES taxes(id), -- NULL means remove tax

    -- Metadata
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_fiscal_tax_mapping UNIQUE(fiscal_position_id, source_tax_id)
);

CREATE INDEX idx_fiscal_tax_mappings_position ON fiscal_position_tax_mappings(fiscal_position_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE fiscal_position_tax_mappings IS 'Maps one tax to another under a fiscal position';

-- Enable RLS
ALTER TABLE fiscal_position_tax_mappings ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_fiscal_tax_mappings ON fiscal_position_tax_mappings FOR ALL USING (
    fiscal_position_id IN (SELECT id FROM fiscal_positions WHERE organization_id IN (
        SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
    ))
);

-- Add tax_id to invoice/bill lines (extend existing tables)
ALTER TABLE customer_invoice_items ADD COLUMN IF NOT EXISTS tax_id UUID REFERENCES taxes(id);
ALTER TABLE vendor_bill_items ADD COLUMN IF NOT EXISTS tax_id UUID REFERENCES taxes(id);
ALTER TABLE journal_entry_lines ADD COLUMN IF NOT EXISTS tax_id UUID REFERENCES taxes(id);

CREATE INDEX IF NOT EXISTS idx_customer_invoice_items_tax ON customer_invoice_items(tax_id);
CREATE INDEX IF NOT EXISTS idx_vendor_bill_items_tax ON vendor_bill_items(tax_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_tax ON journal_entry_lines(tax_id);

-- =====================================================
-- MODULE 3: MULTI-CURRENCY SUPPORT
-- =====================================================

\echo 'Creating multi-currency module...';

-- Currencies
CREATE TABLE IF NOT EXISTS currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ISO currency info
    currency_code VARCHAR(3) NOT NULL UNIQUE, -- ISO 4217 (USD, EUR, GBP)
    currency_name VARCHAR(100) NOT NULL,
    currency_symbol VARCHAR(10),

    -- Precision
    decimal_places INTEGER DEFAULT 2,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_currencies_active ON currencies(is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE currencies IS 'Supported currencies (ISO 4217 codes)';

-- Insert common currencies
INSERT INTO currencies (currency_code, currency_name, currency_symbol, decimal_places) VALUES
    ('USD', 'US Dollar', '$', 2),
    ('EUR', 'Euro', '€', 2),
    ('GBP', 'British Pound', '£', 2),
    ('JPY', 'Japanese Yen', '¥', 0),
    ('CAD', 'Canadian Dollar', 'CA$', 2),
    ('AUD', 'Australian Dollar', 'A$', 2),
    ('CHF', 'Swiss Franc', 'CHF', 2),
    ('CNY', 'Chinese Yuan', '¥', 2),
    ('INR', 'Indian Rupee', '₹', 2),
    ('MXN', 'Mexican Peso', 'MX$', 2)
ON CONFLICT (currency_code) DO NOTHING;

-- Currency Rates
CREATE TABLE IF NOT EXISTS currency_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Rate identification
    currency_code VARCHAR(3) NOT NULL REFERENCES currencies(currency_code),
    rate_date DATE NOT NULL,

    -- Exchange rate (relative to organization base currency)
    rate NUMERIC(20, 10) NOT NULL,

    -- Source
    source VARCHAR(50) DEFAULT 'manual', -- manual, ECB, API, etc.

    -- Metadata
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_currency_rate_org_date UNIQUE(organization_id, currency_code, rate_date)
);

CREATE INDEX idx_currency_rates_org_date ON currency_rates(organization_id, rate_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_currency_rates_currency ON currency_rates(currency_code) WHERE deleted_at IS NULL;

COMMENT ON TABLE currency_rates IS 'Daily exchange rates relative to organization base currency';

-- Enable RLS
ALTER TABLE currency_rates ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_currency_rates ON currency_rates FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Add currency fields to existing tables
ALTER TABLE journal_entries ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3) REFERENCES currencies(currency_code) DEFAULT 'USD';
ALTER TABLE journal_entries ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 10) DEFAULT 1.0;

ALTER TABLE journal_entry_lines ADD COLUMN IF NOT EXISTS amount_currency NUMERIC(20, 4); -- Amount in transaction currency
ALTER TABLE journal_entry_lines ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3) REFERENCES currencies(currency_code);

ALTER TABLE customer_invoices ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3) REFERENCES currencies(currency_code) DEFAULT 'USD';
ALTER TABLE customer_invoices ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 10) DEFAULT 1.0;

ALTER TABLE vendor_bills ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3) REFERENCES currencies(currency_code) DEFAULT 'USD';
ALTER TABLE vendor_bills ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 10) DEFAULT 1.0;

-- Add base_currency to organizations (extend)
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS base_currency_code VARCHAR(3) REFERENCES currencies(currency_code) DEFAULT 'USD';

-- =====================================================
-- MODULE 4: PAYMENT TERMS AND SCHEDULES
-- =====================================================

\echo 'Creating payment terms module...';

-- Payment Terms
CREATE TABLE IF NOT EXISTS payment_terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Identification
    term_code VARCHAR(20) NOT NULL,
    term_name VARCHAR(255) NOT NULL,

    -- Configuration
    note TEXT,
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_payment_term_code_org UNIQUE(organization_id, term_code)
);

CREATE INDEX idx_payment_terms_org ON payment_terms(organization_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE payment_terms IS 'Reusable payment term definitions (Net 30, 50% now + 50% in 30 days, etc.)';

-- Enable RLS
ALTER TABLE payment_terms ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_payment_terms ON payment_terms FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Payment Term Lines
CREATE TABLE IF NOT EXISTS payment_term_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_term_id UUID NOT NULL REFERENCES payment_terms(id) ON DELETE CASCADE,

    -- Line configuration
    sequence INTEGER NOT NULL DEFAULT 10,

    -- Value specification
    value_type VARCHAR(20) NOT NULL CHECK (value_type IN ('percentage', 'fixed', 'balance')),
    value_amount NUMERIC(20, 4), -- Percentage (0-100) or fixed amount

    -- Due date calculation
    days_after INTEGER DEFAULT 0, -- Days after invoice date
    end_of_month BOOLEAN DEFAULT false,
    day_of_month INTEGER, -- Fixed day (e.g., 15th of month)

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT chk_payment_term_line_value CHECK (
        (value_type = 'balance') OR (value_amount IS NOT NULL)
    )
);

CREATE INDEX idx_payment_term_lines_term ON payment_term_lines(payment_term_id, sequence) WHERE deleted_at IS NULL;

COMMENT ON TABLE payment_term_lines IS 'Defines payment split and due date calculation rules';
COMMENT ON COLUMN payment_term_lines.value_type IS 'percentage: % of total, fixed: absolute amount, balance: remainder';

-- Invoice Payment Schedules
CREATE TABLE IF NOT EXISTS invoice_payment_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Link to source document
    source_type VARCHAR(30) NOT NULL CHECK (source_type IN ('customer_invoice', 'vendor_bill')),
    source_id UUID NOT NULL,

    -- Schedule line details
    line_number INTEGER NOT NULL,
    due_date DATE NOT NULL,
    amount_due NUMERIC(20, 4) NOT NULL,
    amount_paid NUMERIC(20, 4) DEFAULT 0,

    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'partial', 'paid', 'overdue')),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_invoice_payment_schedules_source ON invoice_payment_schedules(source_type, source_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoice_payment_schedules_due_date ON invoice_payment_schedules(organization_id, due_date) WHERE deleted_at IS NULL;

COMMENT ON TABLE invoice_payment_schedules IS 'Calculated due amounts and dates for invoices/bills based on payment terms';

-- Enable RLS
ALTER TABLE invoice_payment_schedules ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_payment_schedules ON invoice_payment_schedules FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Add payment_term_id to invoices/bills
ALTER TABLE customer_invoices ADD COLUMN IF NOT EXISTS payment_term_id UUID REFERENCES payment_terms(id);
ALTER TABLE vendor_bills ADD COLUMN IF NOT EXISTS payment_term_id UUID REFERENCES payment_terms(id);

-- =====================================================
-- MODULE 5: ANALYTIC ACCOUNTING
-- =====================================================

\echo 'Creating analytic accounting module...';

-- Analytic Plans (Dimensions)
CREATE TABLE IF NOT EXISTS analytic_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Identification
    plan_code VARCHAR(20) NOT NULL,
    plan_name VARCHAR(255) NOT NULL,

    -- Configuration
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_analytic_plan_code_org UNIQUE(organization_id, plan_code)
);

CREATE INDEX idx_analytic_plans_org ON analytic_plans(organization_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE analytic_plans IS 'Analytic dimensions (Projects, Departments, Regions, etc.)';

-- Enable RLS
ALTER TABLE analytic_plans ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_analytic_plans ON analytic_plans FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Analytic Accounts
CREATE TABLE IF NOT EXISTS analytic_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    analytic_plan_id UUID REFERENCES analytic_plans(id),

    -- Identification
    account_code VARCHAR(50) NOT NULL,
    account_name VARCHAR(255) NOT NULL,

    -- Hierarchy
    parent_account_id UUID REFERENCES analytic_accounts(id),
    account_level INTEGER DEFAULT 1,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    description TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_analytic_account_code_org UNIQUE(organization_id, account_code)
);

CREATE INDEX idx_analytic_accounts_org ON analytic_accounts(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_analytic_accounts_plan ON analytic_accounts(analytic_plan_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_analytic_accounts_parent ON analytic_accounts(parent_account_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE analytic_accounts IS 'Cost centers, projects, departments - for tracking financial dimensions';

-- Enable RLS
ALTER TABLE analytic_accounts ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_analytic_accounts ON analytic_accounts FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Add analytic_account_id to journal lines
ALTER TABLE journal_entry_lines ADD COLUMN IF NOT EXISTS analytic_account_id UUID REFERENCES analytic_accounts(id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_analytic ON journal_entry_lines(analytic_account_id);

-- Add to invoice/bill lines
ALTER TABLE customer_invoice_items ADD COLUMN IF NOT EXISTS analytic_account_id UUID REFERENCES analytic_accounts(id);
ALTER TABLE vendor_bill_items ADD COLUMN IF NOT EXISTS analytic_account_id UUID REFERENCES analytic_accounts(id);

-- =====================================================
-- MODULE 6: DEFERRED REVENUE AND EXPENSES
-- =====================================================

\echo 'Creating deferral module...';

-- Deferred Revenue Contracts
CREATE TABLE IF NOT EXISTS deferred_revenue_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Link to source
    customer_invoice_id UUID REFERENCES customer_invoices(id),
    invoice_line_id UUID REFERENCES customer_invoice_items(id),

    -- Contract details
    contract_name VARCHAR(255),
    total_deferred_amount NUMERIC(20, 4) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    -- Recognition method
    recognition_method VARCHAR(30) DEFAULT 'straight_line' CHECK (
        recognition_method IN ('straight_line', 'custom', 'milestone')
    ),

    -- Accounts
    deferred_account_id UUID NOT NULL REFERENCES chart_of_accounts(id),
    revenue_account_id UUID NOT NULL REFERENCES chart_of_accounts(id),

    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('draft', 'active', 'completed', 'canceled')),
    recognized_amount NUMERIC(20, 4) DEFAULT 0,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_deferred_revenue_org ON deferred_revenue_contracts(organization_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_deferred_revenue_invoice ON deferred_revenue_contracts(customer_invoice_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE deferred_revenue_contracts IS 'Revenue received upfront but recognized over time';

-- Enable RLS
ALTER TABLE deferred_revenue_contracts ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_deferred_revenue ON deferred_revenue_contracts FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Deferred Revenue Schedule
CREATE TABLE IF NOT EXISTS deferred_revenue_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES deferred_revenue_contracts(id) ON DELETE CASCADE,

    -- Schedule line
    line_number INTEGER NOT NULL,
    recognition_date DATE NOT NULL,
    recognition_amount NUMERIC(20, 4) NOT NULL,

    -- Posting status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'posted', 'canceled')),
    journal_entry_id UUID REFERENCES journal_entries(id),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    posted_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_deferred_rev_schedule_line UNIQUE(contract_id, line_number)
);

CREATE INDEX idx_deferred_rev_schedule_contract ON deferred_revenue_schedule(contract_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_deferred_rev_schedule_date ON deferred_revenue_schedule(recognition_date, status) WHERE deleted_at IS NULL;

COMMENT ON TABLE deferred_revenue_schedule IS 'Planned recognition entries for deferred revenue';

-- Deferred Expense Contracts
CREATE TABLE IF NOT EXISTS deferred_expense_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Link to source
    vendor_bill_id UUID REFERENCES vendor_bills(id),
    bill_line_id UUID REFERENCES vendor_bill_items(id),

    -- Contract details
    contract_name VARCHAR(255),
    total_deferred_amount NUMERIC(20, 4) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    -- Recognition method
    recognition_method VARCHAR(30) DEFAULT 'straight_line' CHECK (
        recognition_method IN ('straight_line', 'custom', 'usage_based')
    ),

    -- Accounts
    deferred_account_id UUID NOT NULL REFERENCES chart_of_accounts(id),
    expense_account_id UUID NOT NULL REFERENCES chart_of_accounts(id),

    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('draft', 'active', 'completed', 'canceled')),
    recognized_amount NUMERIC(20, 4) DEFAULT 0,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_deferred_expense_org ON deferred_expense_contracts(organization_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_deferred_expense_bill ON deferred_expense_contracts(vendor_bill_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE deferred_expense_contracts IS 'Expenses paid upfront but recognized over time';

-- Enable RLS
ALTER TABLE deferred_expense_contracts ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_deferred_expense ON deferred_expense_contracts FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Deferred Expense Schedule
CREATE TABLE IF NOT EXISTS deferred_expense_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES deferred_expense_contracts(id) ON DELETE CASCADE,

    -- Schedule line
    line_number INTEGER NOT NULL,
    recognition_date DATE NOT NULL,
    recognition_amount NUMERIC(20, 4) NOT NULL,

    -- Posting status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'posted', 'canceled')),
    journal_entry_id UUID REFERENCES journal_entries(id),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    posted_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_deferred_exp_schedule_line UNIQUE(contract_id, line_number)
);

CREATE INDEX idx_deferred_exp_schedule_contract ON deferred_expense_schedule(contract_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_deferred_exp_schedule_date ON deferred_expense_schedule(recognition_date, status) WHERE deleted_at IS NULL;

COMMENT ON TABLE deferred_expense_schedule IS 'Planned recognition entries for deferred expenses';

-- =====================================================
-- MODULE 7: BANK STATEMENTS AND RECONCILIATION
-- =====================================================

\echo 'Creating bank statement module...';

-- Bank Statements
CREATE TABLE IF NOT EXISTS bank_statements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    bank_account_id UUID NOT NULL REFERENCES bank_accounts(id),

    -- Statement identification
    statement_number VARCHAR(100),
    statement_date DATE NOT NULL,

    -- Date range
    period_start_date DATE NOT NULL,
    period_end_date DATE NOT NULL,

    -- Balances
    opening_balance NUMERIC(20, 4) NOT NULL,
    closing_balance NUMERIC(20, 4) NOT NULL,

    -- Source
    import_source VARCHAR(50) DEFAULT 'manual' CHECK (
        import_source IN ('manual', 'file_import', 'api', 'bank_feed')
    ),
    import_file_name VARCHAR(255),

    -- Status
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'reconciling', 'reconciled', 'closed')),

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_bank_statements_org ON bank_statements(organization_id, statement_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_statements_account ON bank_statements(bank_account_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE bank_statements IS 'Imported or manually created bank statements';

-- Enable RLS
ALTER TABLE bank_statements ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_bank_statements ON bank_statements FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Bank Statement Lines
CREATE TABLE IF NOT EXISTS bank_statement_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bank_statement_id UUID NOT NULL REFERENCES bank_statements(id) ON DELETE CASCADE,

    -- Line details
    line_number INTEGER NOT NULL,
    transaction_date DATE NOT NULL,
    value_date DATE,

    -- Amount
    amount NUMERIC(20, 4) NOT NULL,
    currency_code VARCHAR(3) REFERENCES currencies(currency_code) DEFAULT 'USD',

    -- Transaction details
    description TEXT,
    reference VARCHAR(255),
    counterparty_name VARCHAR(255),
    counterparty_account VARCHAR(100),

    -- Bank reference
    bank_reference VARCHAR(255),
    check_number VARCHAR(50),

    -- Reconciliation status
    status VARCHAR(20) DEFAULT 'unmatched' CHECK (
        status IN ('unmatched', 'matched', 'partial_match', 'ignored')
    ),

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_bank_stmt_line UNIQUE(bank_statement_id, line_number)
);

CREATE INDEX idx_bank_stmt_lines_statement ON bank_statement_lines(bank_statement_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_stmt_lines_date ON bank_statement_lines(transaction_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_stmt_lines_status ON bank_statement_lines(status) WHERE deleted_at IS NULL AND status = 'unmatched';

COMMENT ON TABLE bank_statement_lines IS 'Individual transactions from bank statements';

-- Bank Statement Reconciliation Links
CREATE TABLE IF NOT EXISTS bank_statement_reconciliations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Links statement line to GL entries or payments
    bank_statement_line_id UUID NOT NULL REFERENCES bank_statement_lines(id),

    -- Can link to various transaction types
    journal_entry_id UUID REFERENCES journal_entries(id),
    payment_id UUID, -- Generic payment reference (vendor or customer)

    -- Match details
    matched_amount NUMERIC(20, 4) NOT NULL,

    -- Metadata
    matched_by UUID REFERENCES users(id),
    matched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_bank_stmt_recon_line ON bank_statement_reconciliations(bank_statement_line_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_stmt_recon_je ON bank_statement_reconciliations(journal_entry_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE bank_statement_reconciliations IS 'Links bank statement lines to ledger entries/payments';

-- Enable RLS
ALTER TABLE bank_statement_reconciliations ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_bank_stmt_recon ON bank_statement_reconciliations FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Reconciliation Rule Models (Auto-match patterns)
CREATE TABLE IF NOT EXISTS reconciliation_rule_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Rule identification
    rule_name VARCHAR(255) NOT NULL,
    rule_code VARCHAR(50),
    sequence INTEGER DEFAULT 10,

    -- Match criteria
    amount_min NUMERIC(20, 4),
    amount_max NUMERIC(20, 4),
    description_pattern VARCHAR(500), -- Regex or LIKE pattern
    counterparty_pattern VARCHAR(500),
    reference_pattern VARCHAR(500),

    -- Actions
    journal_id UUID REFERENCES journals(id),
    account_id UUID REFERENCES chart_of_accounts(id),
    analytic_account_id UUID REFERENCES analytic_accounts(id),
    tax_id UUID REFERENCES taxes(id),

    -- Status
    is_active BOOLEAN DEFAULT true,
    auto_apply BOOLEAN DEFAULT false,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_recon_rules_org ON reconciliation_rule_models(organization_id, sequence) WHERE deleted_at IS NULL AND is_active = true;

COMMENT ON TABLE reconciliation_rule_models IS 'Patterns for auto-matching bank statement lines to GL entries';

-- Enable RLS
ALTER TABLE reconciliation_rule_models ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_recon_rules ON reconciliation_rule_models FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- =====================================================
-- MODULE 8: BUDGETS AND ANALYTIC BUDGETS
-- =====================================================

\echo 'Creating budget module...';

-- Budgets
CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Budget identification
    budget_code VARCHAR(50) NOT NULL,
    budget_name VARCHAR(255) NOT NULL,

    -- Period
    fiscal_year_id UUID REFERENCES fiscal_years(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    -- Type
    budget_type VARCHAR(30) DEFAULT 'operating' CHECK (
        budget_type IN ('operating', 'capital', 'cash_flow', 'project', 'departmental')
    ),

    -- Status
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'approved', 'active', 'closed')),

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_budget_code_org UNIQUE(organization_id, budget_code)
);

CREATE INDEX idx_budgets_org ON budgets(organization_id, fiscal_year_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_budgets_status ON budgets(status) WHERE deleted_at IS NULL;

COMMENT ON TABLE budgets IS 'Budget scenarios for planning vs actual comparison';

-- Enable RLS
ALTER TABLE budgets ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_budgets ON budgets FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Budget Lines
CREATE TABLE IF NOT EXISTS budget_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,

    -- Dimension
    account_id UUID REFERENCES chart_of_accounts(id),
    analytic_account_id UUID REFERENCES analytic_accounts(id),

    -- Period (optional - for monthly/quarterly breakdown)
    accounting_period_id UUID REFERENCES accounting_periods(id),
    period_start_date DATE,
    period_end_date DATE,

    -- Planned amount
    planned_amount NUMERIC(20, 4) NOT NULL,

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT chk_budget_line_dimension CHECK (
        account_id IS NOT NULL OR analytic_account_id IS NOT NULL
    )
);

CREATE INDEX idx_budget_lines_budget ON budget_lines(budget_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_budget_lines_account ON budget_lines(account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_budget_lines_analytic ON budget_lines(analytic_account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_budget_lines_period ON budget_lines(accounting_period_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE budget_lines IS 'Planned amounts by account and/or analytic dimension';
COMMENT ON COLUMN budget_lines.planned_amount IS 'Positive for revenue/income, negative for expenses';

-- =====================================================
-- MODULE 9: LOCALIZATION AND TAX REPORTING
-- =====================================================

\echo 'Creating localization module...';

-- Localization Packages
CREATE TABLE IF NOT EXISTS localization_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Package identification
    package_code VARCHAR(50) NOT NULL UNIQUE, -- e.g., 'us_gaap', 'uk_vat', 'eu_ifrs'
    package_name VARCHAR(255) NOT NULL,

    -- Jurisdiction
    country_code VARCHAR(2), -- ISO country code
    region VARCHAR(100),

    -- Content
    description TEXT,
    version VARCHAR(20),

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_localization_packages_country ON localization_packages(country_code) WHERE deleted_at IS NULL;

COMMENT ON TABLE localization_packages IS 'Country-specific templates (COA, taxes, journals, reports)';

-- Insert default localization packages
INSERT INTO localization_packages (package_code, package_name, country_code, description) VALUES
    ('us_gaap', 'United States - GAAP', 'US', 'US Generally Accepted Accounting Principles'),
    ('uk_vat', 'United Kingdom - VAT', 'GB', 'UK VAT and Companies House compliant'),
    ('eu_ifrs', 'European Union - IFRS', NULL, 'International Financial Reporting Standards'),
    ('ca_gaap', 'Canada - GAAP', 'CA', 'Canadian GAAP and CRA compliant'),
    ('au_aas', 'Australia - AAS', 'AU', 'Australian Accounting Standards')
ON CONFLICT (package_code) DO NOTHING;

-- Organization Localization Link
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS localization_package_id UUID REFERENCES localization_packages(id);

-- Tax Report Definitions
CREATE TABLE IF NOT EXISTS tax_report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE, -- NULL = global/template
    localization_package_id UUID REFERENCES localization_packages(id),

    -- Report identification
    report_code VARCHAR(50) NOT NULL,
    report_name VARCHAR(255) NOT NULL,

    -- Jurisdiction
    jurisdiction VARCHAR(100),
    authority VARCHAR(255), -- e.g., "IRS", "HMRC", "ATO"

    -- Period
    report_frequency VARCHAR(20) CHECK (report_frequency IN ('monthly', 'quarterly', 'annual', 'on_demand')),

    -- Configuration
    version VARCHAR(20),
    effective_from DATE,
    effective_to DATE,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_tax_report_defs_org ON tax_report_definitions(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tax_report_defs_package ON tax_report_definitions(localization_package_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE tax_report_definitions IS 'Declarative tax reports (VAT returns, sales tax reports, etc.)';

-- Enable RLS
ALTER TABLE tax_report_definitions ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_tax_reports ON tax_report_definitions FOR ALL USING (
    organization_id IS NULL OR organization_id IN (
        SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
    )
);

-- Tax Report Lines
CREATE TABLE IF NOT EXISTS tax_report_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tax_report_definition_id UUID NOT NULL REFERENCES tax_report_definitions(id) ON DELETE CASCADE,

    -- Line identification
    line_code VARCHAR(50) NOT NULL,
    line_name VARCHAR(255) NOT NULL,
    sequence INTEGER DEFAULT 10,

    -- Parent line (for hierarchy)
    parent_line_id UUID REFERENCES tax_report_lines(id),

    -- Calculation
    formula_type VARCHAR(30) CHECK (formula_type IN ('sum', 'detail', 'formula', 'manual')),
    formula TEXT, -- SQL expression or formula for calculation

    -- Data source filters
    tax_group_ids UUID[], -- Array of tax_group IDs
    account_ids UUID[], -- Array of chart_of_accounts IDs
    tax_ids UUID[], -- Array of tax IDs

    -- Display
    is_subtotal BOOLEAN DEFAULT false,
    is_total BOOLEAN DEFAULT false,

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_tax_report_lines_report ON tax_report_lines(tax_report_definition_id, sequence) WHERE deleted_at IS NULL;
CREATE INDEX idx_tax_report_lines_parent ON tax_report_lines(parent_line_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE tax_report_lines IS 'Lines/boxes within a tax report with data source mappings';

-- =====================================================
-- AUDIT TRIGGERS
-- =====================================================

\echo 'Creating audit triggers...';

-- Create triggers for updated_at
CREATE TRIGGER update_journals_updated_at BEFORE UPDATE ON journals
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_tax_groups_updated_at BEFORE UPDATE ON tax_groups
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_taxes_updated_at BEFORE UPDATE ON taxes
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_fiscal_positions_updated_at BEFORE UPDATE ON fiscal_positions
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_currencies_updated_at BEFORE UPDATE ON currencies
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_payment_terms_updated_at BEFORE UPDATE ON payment_terms
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_analytic_plans_updated_at BEFORE UPDATE ON analytic_plans
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_analytic_accounts_updated_at BEFORE UPDATE ON analytic_accounts
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_bank_statements_updated_at BEFORE UPDATE ON bank_statements
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

CREATE TRIGGER update_budgets_updated_at BEFORE UPDATE ON budgets
    FOR EACH ROW EXECUTE FUNCTION accounting.set_updated_at();

\echo '';
\echo '==========================================';
\echo 'Odoo-Style Extensions Created Successfully';
\echo '==========================================';
\echo 'Modules Created:';
\echo '  1. ✓ Journals (Sales, Purchase, Bank, Cash, General)';
\echo '  2. ✓ Tax Engine (Tax Groups, Taxes, Fiscal Positions)';
\echo '  3. ✓ Multi-Currency (Currencies, Exchange Rates)';
\echo '  4. ✓ Payment Terms (Terms, Lines, Schedules)';
\echo '  5. ✓ Analytic Accounting (Plans, Accounts)';
\echo '  6. ✓ Deferrals (Revenue & Expense Contracts, Schedules)';
\echo '  7. ✓ Bank Statements (Statements, Lines, Reconciliation)';
\echo '  8. ✓ Budgets (Budget Scenarios, Budget Lines)';
\echo '  9. ✓ Localization (Packages, Tax Reports)';
\echo '==========================================';
\echo '';

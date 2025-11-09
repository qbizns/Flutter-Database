-- =====================================================
-- Accounting Module - V002: AP, AR, Fixed Assets, Banks
-- Description: Accounts Payable, Accounts Receivable, Assets, Banking
-- Based on: QuickBooks, Xero, Sage best practices
-- Date: 2025-11-09
-- Dependencies: V001 (core accounting)
-- =====================================================

-- =====================================================
-- SECTION 1: Accounts Payable (AP)
-- =====================================================

CREATE TABLE vendor_bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Bill Identification
    bill_number VARCHAR(50) NOT NULL,
    vendor_bill_number VARCHAR(100), -- Vendor's invoice number

    -- Vendor
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,

    -- Dates
    bill_date DATE NOT NULL,
    due_date DATE NOT NULL,
    payment_terms VARCHAR(50), -- Net 30, Net 60, Due on Receipt

    -- Period
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,

    -- Amounts
    subtotal NUMERIC(20, 4) DEFAULT 0,
    tax_amount NUMERIC(20, 4) DEFAULT 0,
    total_amount NUMERIC(20, 4) NOT NULL,
    paid_amount NUMERIC(20, 4) DEFAULT 0,
    balance_due NUMERIC(20, 4) DEFAULT 0,

    -- Status
    status VARCHAR(30) DEFAULT 'unpaid', -- unpaid, partial, paid, overdue, void

    -- GL Posting
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    is_posted BOOLEAN DEFAULT false,

    -- Purchase Order Reference
    purchase_order_id UUID REFERENCES purchase_orders(id) ON DELETE SET NULL,

    -- Description
    description TEXT,
    notes TEXT,
    memo TEXT,

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

    CONSTRAINT chk_vendor_bills_status CHECK (status IN ('unpaid', 'partial', 'paid', 'overdue', 'void')),
    CONSTRAINT chk_vendor_bills_amounts CHECK (total_amount >= 0 AND paid_amount >= 0 AND balance_due >= 0),
    CONSTRAINT chk_vendor_bills_balance CHECK (balance_due = total_amount - paid_amount),
    UNIQUE(organization_id, bill_number)
);

CREATE INDEX idx_vendor_bills_org ON vendor_bills(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_bills_supplier ON vendor_bills(supplier_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_bills_status ON vendor_bills(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_bills_due_date ON vendor_bills(due_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_bills_period ON vendor_bills(accounting_period_id) WHERE deleted_at IS NULL;

CREATE TRIGGER update_vendor_bills_updated_at
    BEFORE UPDATE ON vendor_bills FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE vendor_bill_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Bill Reference
    vendor_bill_id UUID NOT NULL REFERENCES vendor_bills(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,

    -- Account
    expense_account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,

    -- Item Details
    description TEXT NOT NULL,
    quantity NUMERIC(10, 2) DEFAULT 1,
    unit_price NUMERIC(20, 4) NOT NULL,
    amount NUMERIC(20, 4) NOT NULL,

    -- Dimensions
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    department VARCHAR(100),
    project_code VARCHAR(100),

    -- Tax
    tax_code VARCHAR(50),
    tax_amount NUMERIC(20, 4) DEFAULT 0,

    -- Product Reference (if applicable)
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_vendor_bill_lines_amounts CHECK (quantity > 0 AND unit_price >= 0 AND amount >= 0),
    UNIQUE(vendor_bill_id, line_number)
);

CREATE INDEX idx_vendor_bill_lines_org ON vendor_bill_lines(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_bill_lines_bill ON vendor_bill_lines(vendor_bill_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_bill_lines_account ON vendor_bill_lines(expense_account_id) WHERE deleted_at IS NULL;

CREATE TRIGGER update_vendor_bill_lines_updated_at
    BEFORE UPDATE ON vendor_bill_lines FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE vendor_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Payment Identification
    payment_number VARCHAR(50) NOT NULL,

    -- Vendor
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,

    -- Payment Details
    payment_date DATE NOT NULL,
    payment_method VARCHAR(30) NOT NULL, -- check, cash, wire, ach, card, other
    reference_number VARCHAR(100), -- Check number, transaction ID

    -- Amount
    payment_amount NUMERIC(20, 4) NOT NULL,

    -- Bank Account
    bank_account_id UUID REFERENCES chart_of_accounts(id) ON DELETE SET NULL,

    -- Period
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,

    -- GL Posting
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    is_posted BOOLEAN DEFAULT false,

    -- Description
    memo TEXT,
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_vendor_payments_method CHECK (payment_method IN ('check', 'cash', 'wire', 'ach', 'card', 'other')),
    CONSTRAINT chk_vendor_payments_amount CHECK (payment_amount > 0),
    UNIQUE(organization_id, payment_number)
);

CREATE INDEX idx_vendor_payments_org ON vendor_payments(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_payments_supplier ON vendor_payments(supplier_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_payments_date ON vendor_payments(payment_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_vendor_payments_bank ON vendor_payments(bank_account_id) WHERE deleted_at IS NULL;

CREATE TRIGGER update_vendor_payments_updated_at
    BEFORE UPDATE ON vendor_payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE vendor_payment_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- References
    vendor_payment_id UUID NOT NULL REFERENCES vendor_payments(id) ON DELETE CASCADE,
    vendor_bill_id UUID NOT NULL REFERENCES vendor_bills(id) ON DELETE CASCADE,

    -- Application
    applied_amount NUMERIC(20, 4) NOT NULL,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_vendor_payment_apps_amount CHECK (applied_amount > 0),
    UNIQUE(vendor_payment_id, vendor_bill_id)
);

CREATE INDEX idx_vendor_payment_apps_payment ON vendor_payment_applications(vendor_payment_id);
CREATE INDEX idx_vendor_payment_apps_bill ON vendor_payment_applications(vendor_bill_id);

-- =====================================================
-- SECTION 2: Accounts Receivable (AR)
-- =====================================================

CREATE TABLE customer_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Invoice Identification
    invoice_number VARCHAR(50) NOT NULL,

    -- Customer
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,

    -- Dates
    invoice_date DATE NOT NULL,
    due_date DATE NOT NULL,
    payment_terms VARCHAR(50), -- Net 30, Net 60, Due on Receipt

    -- Period
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,

    -- Amounts
    subtotal NUMERIC(20, 4) DEFAULT 0,
    tax_amount NUMERIC(20, 4) DEFAULT 0,
    discount_amount NUMERIC(20, 4) DEFAULT 0,
    total_amount NUMERIC(20, 4) NOT NULL,
    paid_amount NUMERIC(20, 4) DEFAULT 0,
    balance_due NUMERIC(20, 4) DEFAULT 0,

    -- Status
    status VARCHAR(30) DEFAULT 'unpaid', -- unpaid, partial, paid, overdue, void

    -- GL Posting
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    is_posted BOOLEAN DEFAULT false,

    -- Sales Order Reference
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,

    -- Description
    description TEXT,
    notes TEXT,
    memo TEXT,

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

    CONSTRAINT chk_customer_invoices_status CHECK (status IN ('unpaid', 'partial', 'paid', 'overdue', 'void')),
    CONSTRAINT chk_customer_invoices_amounts CHECK (total_amount >= 0 AND paid_amount >= 0 AND balance_due >= 0),
    CONSTRAINT chk_customer_invoices_balance CHECK (balance_due = total_amount - paid_amount),
    UNIQUE(organization_id, invoice_number)
);

CREATE INDEX idx_customer_invoices_org ON customer_invoices(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_invoices_customer ON customer_invoices(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_invoices_status ON customer_invoices(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_invoices_due_date ON customer_invoices(due_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_invoices_period ON customer_invoices(accounting_period_id) WHERE deleted_at IS NULL;

CREATE TRIGGER update_customer_invoices_updated_at
    BEFORE UPDATE ON customer_invoices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE customer_invoice_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Invoice Reference
    customer_invoice_id UUID NOT NULL REFERENCES customer_invoices(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,

    -- Account
    revenue_account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,

    -- Item Details
    description TEXT NOT NULL,
    quantity NUMERIC(10, 2) DEFAULT 1,
    unit_price NUMERIC(20, 4) NOT NULL,
    amount NUMERIC(20, 4) NOT NULL,

    -- Dimensions
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    department VARCHAR(100),
    project_code VARCHAR(100),

    -- Tax
    tax_code VARCHAR(50),
    tax_amount NUMERIC(20, 4) DEFAULT 0,

    -- Product Reference
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_customer_invoice_lines_amounts CHECK (quantity > 0 AND unit_price >= 0 AND amount >= 0),
    UNIQUE(customer_invoice_id, line_number)
);

CREATE INDEX idx_customer_invoice_lines_org ON customer_invoice_lines(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_invoice_lines_invoice ON customer_invoice_lines(customer_invoice_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_invoice_lines_account ON customer_invoice_lines(revenue_account_id) WHERE deleted_at IS NULL;

CREATE TRIGGER update_customer_invoice_lines_updated_at
    BEFORE UPDATE ON customer_invoice_lines FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE customer_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Payment Identification
    payment_number VARCHAR(50) NOT NULL,

    -- Customer
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,

    -- Payment Details
    payment_date DATE NOT NULL,
    payment_method VARCHAR(30) NOT NULL, -- cash, check, card, wire, ach, other
    reference_number VARCHAR(100), -- Check number, transaction ID

    -- Amount
    payment_amount NUMERIC(20, 4) NOT NULL,

    -- Deposit Account
    deposit_account_id UUID REFERENCES chart_of_accounts(id) ON DELETE SET NULL,

    -- Period
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,

    -- GL Posting
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    is_posted BOOLEAN DEFAULT false,

    -- Description
    memo TEXT,
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_customer_payments_method CHECK (payment_method IN ('cash', 'check', 'card', 'wire', 'ach', 'other')),
    CONSTRAINT chk_customer_payments_amount CHECK (payment_amount > 0),
    UNIQUE(organization_id, payment_number)
);

CREATE INDEX idx_customer_payments_org ON customer_payments(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_payments_customer ON customer_payments(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_payments_date ON customer_payments(payment_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_payments_account ON customer_payments(deposit_account_id) WHERE deleted_at IS NULL;

CREATE TRIGGER update_customer_payments_updated_at
    BEFORE UPDATE ON customer_payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE customer_payment_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- References
    customer_payment_id UUID NOT NULL REFERENCES customer_payments(id) ON DELETE CASCADE,
    customer_invoice_id UUID NOT NULL REFERENCES customer_invoices(id) ON DELETE CASCADE,

    -- Application
    applied_amount NUMERIC(20, 4) NOT NULL,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_customer_payment_apps_amount CHECK (applied_amount > 0),
    UNIQUE(customer_payment_id, customer_invoice_id)
);

CREATE INDEX idx_customer_payment_apps_payment ON customer_payment_applications(customer_payment_id);
CREATE INDEX idx_customer_payment_apps_invoice ON customer_payment_applications(customer_invoice_id);

-- =====================================================
-- SECTION 3: Fixed Assets
-- =====================================================

CREATE TABLE asset_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Category Details
    category_code VARCHAR(50) NOT NULL UNIQUE,
    category_name VARCHAR(100) NOT NULL,

    -- Depreciation Defaults
    default_depreciation_method VARCHAR(30), -- straight_line, declining_balance, units_of_production
    default_useful_life_years INTEGER,
    default_salvage_value_percent NUMERIC(5, 2),

    -- Accounts
    asset_account_id UUID REFERENCES chart_of_accounts(id) ON DELETE SET NULL,
    accumulated_depreciation_account_id UUID REFERENCES chart_of_accounts(id) ON DELETE SET NULL,
    depreciation_expense_account_id UUID REFERENCES chart_of_accounts(id) ON DELETE SET NULL,

    -- Description
    description TEXT,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_asset_categories_code ON asset_categories(category_code);

-- =====================================================

CREATE TABLE fixed_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Asset Identification
    asset_number VARCHAR(50) NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
    asset_category_id UUID REFERENCES asset_categories(id) ON DELETE SET NULL,

    -- Acquisition
    acquisition_date DATE NOT NULL,
    acquisition_cost NUMERIC(20, 4) NOT NULL,
    salvage_value NUMERIC(20, 4) DEFAULT 0,

    -- Supplier/Vendor
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,
    vendor_bill_id UUID REFERENCES vendor_bills(id) ON DELETE SET NULL,

    -- Depreciation
    depreciation_method VARCHAR(30) NOT NULL, -- straight_line, declining_balance, units_of_production
    useful_life_years INTEGER NOT NULL,
    depreciation_start_date DATE NOT NULL,

    -- Accounts
    asset_account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,
    accumulated_depreciation_account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,
    depreciation_expense_account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,

    -- Current Values
    current_book_value NUMERIC(20, 4) DEFAULT 0,
    accumulated_depreciation NUMERIC(20, 4) DEFAULT 0,
    last_depreciation_date DATE,

    -- Location
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    department VARCHAR(100),

    -- Disposal
    is_disposed BOOLEAN DEFAULT false,
    disposal_date DATE,
    disposal_proceeds NUMERIC(20, 4),
    disposal_journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,

    -- Description
    description TEXT,
    serial_number VARCHAR(100),
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_fixed_assets_method CHECK (depreciation_method IN ('straight_line', 'declining_balance', 'units_of_production')),
    CONSTRAINT chk_fixed_assets_amounts CHECK (acquisition_cost >= 0 AND salvage_value >= 0 AND salvage_value <= acquisition_cost),
    CONSTRAINT chk_fixed_assets_life CHECK (useful_life_years > 0),
    UNIQUE(organization_id, asset_number)
);

CREATE INDEX idx_fixed_assets_org ON fixed_assets(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_fixed_assets_category ON fixed_assets(asset_category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_fixed_assets_location ON fixed_assets(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_fixed_assets_disposed ON fixed_assets(is_disposed) WHERE deleted_at IS NULL;

CREATE TRIGGER update_fixed_assets_updated_at
    BEFORE UPDATE ON fixed_assets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE asset_depreciation_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Asset Reference
    fixed_asset_id UUID NOT NULL REFERENCES fixed_assets(id) ON DELETE CASCADE,

    -- Period
    fiscal_year_id UUID REFERENCES fiscal_years(id) ON DELETE SET NULL,
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,
    depreciation_date DATE NOT NULL,

    -- Depreciation
    depreciation_amount NUMERIC(20, 4) NOT NULL,
    accumulated_depreciation_beginning NUMERIC(20, 4) NOT NULL,
    accumulated_depreciation_ending NUMERIC(20, 4) NOT NULL,
    book_value_beginning NUMERIC(20, 4) NOT NULL,
    book_value_ending NUMERIC(20, 4) NOT NULL,

    -- GL Posting
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    is_posted BOOLEAN DEFAULT false,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    posted_at TIMESTAMP WITH TIME ZONE,
    posted_by UUID REFERENCES users(id) ON DELETE SET NULL,

    CONSTRAINT chk_asset_depreciation_amounts CHECK (depreciation_amount >= 0),
    UNIQUE(fixed_asset_id, accounting_period_id)
);

CREATE INDEX idx_asset_depreciation_asset ON asset_depreciation_schedule(fixed_asset_id);
CREATE INDEX idx_asset_depreciation_period ON asset_depreciation_schedule(accounting_period_id);
CREATE INDEX idx_asset_depreciation_posted ON asset_depreciation_schedule(is_posted);

-- =====================================================
-- SECTION 4: Bank Accounts & Reconciliation
-- =====================================================

CREATE TABLE bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Account Reference
    chart_account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE CASCADE,

    -- Bank Details
    bank_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(100) NOT NULL,
    account_type VARCHAR(30), -- checking, savings, credit_card, line_of_credit
    routing_number VARCHAR(50),
    swift_code VARCHAR(50),

    -- Currency
    currency_code VARCHAR(3) DEFAULT 'USD',

    -- Current Balance (from bank statement)
    current_balance NUMERIC(20, 4) DEFAULT 0,
    statement_balance NUMERIC(20, 4) DEFAULT 0,
    last_statement_date DATE,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Online Banking
    online_banking_enabled BOOLEAN DEFAULT false,
    last_sync_date TIMESTAMP WITH TIME ZONE,

    -- Description
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_bank_accounts_type CHECK (account_type IN ('checking', 'savings', 'credit_card', 'line_of_credit')),
    UNIQUE(organization_id, account_number)
);

CREATE INDEX idx_bank_accounts_org ON bank_accounts(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_accounts_chart ON bank_accounts(chart_account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_accounts_active ON bank_accounts(is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER update_bank_accounts_updated_at
    BEFORE UPDATE ON bank_accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE bank_reconciliations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Bank Account
    bank_account_id UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,

    -- Statement Details
    statement_date DATE NOT NULL,
    statement_balance NUMERIC(20, 4) NOT NULL,

    -- Reconciliation
    reconciliation_date DATE,
    book_balance NUMERIC(20, 4) DEFAULT 0,
    cleared_balance NUMERIC(20, 4) DEFAULT 0,
    difference NUMERIC(20, 4) DEFAULT 0,

    -- Status
    status VARCHAR(20) DEFAULT 'in_progress', -- in_progress, reconciled, locked
    is_reconciled BOOLEAN DEFAULT false,

    -- Period
    accounting_period_id UUID REFERENCES accounting_periods(id) ON DELETE SET NULL,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reconciled_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reconciled_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_bank_recon_status CHECK (status IN ('in_progress', 'reconciled', 'locked')),
    UNIQUE(bank_account_id, statement_date)
);

CREATE INDEX idx_bank_recon_org ON bank_reconciliations(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_recon_account ON bank_reconciliations(bank_account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_recon_status ON bank_reconciliations(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_recon_date ON bank_reconciliations(statement_date) WHERE deleted_at IS NULL;

CREATE TRIGGER update_bank_reconciliations_updated_at
    BEFORE UPDATE ON bank_reconciliations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================

CREATE TABLE bank_reconciliation_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Reconciliation Reference
    bank_reconciliation_id UUID REFERENCES bank_reconciliations(id) ON DELETE CASCADE,

    -- GL Transaction
    general_ledger_id UUID NOT NULL REFERENCES general_ledger(id) ON DELETE CASCADE,
    journal_entry_line_id UUID NOT NULL REFERENCES journal_entry_lines(id) ON DELETE CASCADE,

    -- Reconciliation Status
    is_cleared BOOLEAN DEFAULT false,
    cleared_date DATE,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    cleared_by UUID REFERENCES users(id) ON DELETE SET NULL,

    UNIQUE(bank_reconciliation_id, general_ledger_id)
);

CREATE INDEX idx_bank_recon_items_recon ON bank_reconciliation_items(bank_reconciliation_id);
CREATE INDEX idx_bank_recon_items_gl ON bank_reconciliation_items(general_ledger_id);
CREATE INDEX idx_bank_recon_items_cleared ON bank_reconciliation_items(is_cleared);

-- =====================================================
-- SECTION 5: Row-Level Security (RLS)
-- =====================================================

ALTER TABLE vendor_bills ENABLE ROW LEVEL SECURITY;
ALTER TABLE vendor_bill_lines ENABLE ROW LEVEL SECURITY;
ALTER TABLE vendor_payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE vendor_payment_applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_invoice_lines ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_payment_applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE fixed_assets ENABLE ROW LEVEL SECURITY;
ALTER TABLE asset_depreciation_schedule ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_accounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_reconciliations ENABLE ROW LEVEL SECURITY;
ALTER TABLE bank_reconciliation_items ENABLE ROW LEVEL SECURITY;

-- Super admin and org-scoped policies for all tables
CREATE POLICY vendor_bills_access ON vendor_bills FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY vendor_bill_lines_access ON vendor_bill_lines FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY vendor_payments_access ON vendor_payments FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY vendor_payment_apps_access ON vendor_payment_applications FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY customer_invoices_access ON customer_invoices FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY customer_invoice_lines_access ON customer_invoice_lines FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY customer_payments_access ON customer_payments FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY customer_payment_apps_access ON customer_payment_applications FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY fixed_assets_access ON fixed_assets FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY asset_depreciation_access ON asset_depreciation_schedule FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY bank_accounts_access ON bank_accounts FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY bank_recon_access ON bank_reconciliations FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

CREATE POLICY bank_recon_items_access ON bank_reconciliation_items FOR ALL TO PUBLIC
    USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = current_user_id() AND r.role_name = 'Super Admin') OR
           organization_id IN (SELECT organization_id FROM users WHERE id = current_user_id()));

-- =====================================================
-- End of Migration V002
-- =====================================================

-- Summary: Created AP, AR, Fixed Assets, and Banking modules
-- - Complete Accounts Payable workflow
-- - Complete Accounts Receivable workflow
-- - Fixed assets with depreciation
-- - Bank accounts and reconciliation
-- - Complete RLS policies

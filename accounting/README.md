# Accounting Module

Complete double-entry accounting system designed as a **plugin module** for the POS SaaS platform. This module provides comprehensive financial management capabilities while remaining independent of the POS core system.

## 🎯 Overview

The Accounting Module is built as an **optional plugin** that:
- **Shares base tables** from the POS system (organizations, users, suppliers, customers, products)
- **Operates independently** - can exist or not exist without breaking POS functionality
- **Follows accounting standards** - Based on QuickBooks, Xero, Sage, and NetSuite best practices
- **Provides heavy test data** - 200+ realistic transactions for full year 2024 testing

## 📊 Key Features

### Core Accounting
- ✅ **Double-Entry Bookkeeping** with automatic validation
- ✅ **Chart of Accounts** (80+ standard accounts)
- ✅ **General Ledger** (immutable posted transactions)
- ✅ **Journal Entries** (draft → posted workflow)
- ✅ **Fiscal Years & Periods** (12-month periods with open/closed status)

### AP/AR Management
- ✅ **Accounts Payable** (vendor bills, payments, aging)
- ✅ **Accounts Receivable** (customer invoices, payments, aging)
- ✅ **Payment Application** (link payments to bills/invoices)
- ✅ **Aged Reports** (30/60/90 day aging buckets)

### Fixed Assets
- ✅ **Asset Tracking** (acquisition, depreciation, disposal)
- ✅ **Depreciation Methods** (straight-line, declining balance)
- ✅ **Depreciation Schedules** (automated monthly depreciation)
- ✅ **Asset Categories** (furniture, equipment, vehicles, buildings, land)

### Financial Reporting
- ✅ **Trial Balance**
- ✅ **Balance Sheet** (Statement of Financial Position)
- ✅ **Income Statement** (Profit & Loss)
- ✅ **Cash Flow Statement** (Operating, Investing, Financing)
- ✅ **Account Activity** (detailed transaction history with running balance)
- ✅ **Financial Ratios** (liquidity, profitability, leverage)

### Banking & Reconciliation
- ✅ **Bank Accounts** (multiple account support)
- ✅ **Bank Reconciliation** (statement matching)
- ✅ **Cash Management** (track all cash movements)

### POS Integration (V004-V010)
- ✅ **Account Mappings** (products, payment methods, discounts → GL accounts)
- ✅ **Tax Mappings** (POS tax codes → accounting taxes)
- ✅ **Posting Audit** (track which POS docs are posted to accounting)
- ✅ **Inventory Valuation** (FIFO, weighted average, cost layers)
- ✅ **COGS Calculation** (automatic cost recognition)
- ✅ **Multi-Currency Support** (base currency, exchange rates)
- ✅ **Document Sequences** (unified numbering system)
- ✅ **E-Invoicing Links** (connect e-invoices to accounting records)

### Data Integrity & Controls
- ✅ **Immutability Triggers** (prevent modification of posted data)
- ✅ **Fiscal Period Locking** (prevent posting to closed periods)
- ✅ **Year-End Closing** (automated closing entries to retained earnings)
- ✅ **Reversal Entries** (proper correction of posted transactions)
- ✅ **Double-Entry Validation** (debits = credits enforcement)

## 🏗️ Architecture

### Plugin Design

```
┌─────────────────────────────────────────────────────┐
│                    POS System                       │
│  (organizations, users, suppliers, customers, etc)  │
└────────────────┬────────────────────────────────────┘
                 │
                 │ Foreign Key References
                 │ (shared data, no duplication)
                 │
┌────────────────▼────────────────────────────────────┐
│              Accounting Module                      │
│  (COA, GL, JE, AP, AR, Fixed Assets, Reports)      │
└─────────────────────────────────────────────────────┘
```

### Database Schema

The accounting module uses a separate schema: `accounting`

**Key Tables:**
- `fiscal_years` - Annual fiscal periods
- `accounting_periods` - Monthly accounting periods
- `account_types` - ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE, COGS
- `chart_of_accounts` - Account master data (80+ accounts)
- `journal_entries` - Draft/Posted journal entries (header)
- `journal_entry_lines` - JE line items (debits/credits)
- `general_ledger` - Posted transactions (immutable)
- `vendor_bills` & `vendor_payments` - Accounts Payable
- `customer_invoices` & `customer_payments` - Accounts Receivable
- `fixed_assets` & `asset_depreciation_schedule` - Asset management
- `bank_accounts` & `bank_reconciliations` - Banking

## 📁 Directory Structure

```
accounting/
├── README.md                    # This file
├── migrations/                  # Database migrations (DDL)
│   ├── V001_20251109_create_accounting_core.sql
│   ├── V002_20251109_create_ap_ar_assets.sql
│   ├── V003_20251110_create_odoo_extensions.sql
│   ├── V004_20251110_create_pos_account_mappings.sql
│   ├── V005_20251110_create_pos_posting_audit.sql
│   ├── V006_20251110_create_inventory_valuation_settings.sql
│   ├── V007_20251110_create_inventory_valuation_views.sql
│   ├── V008_20251110_create_pos_tax_mappings.sql
│   ├── V009_20251110_add_immutability_triggers.sql
│   └── V010_20251110_add_closing_procedures.sql
├── seed_data/                   # Test data (DML)
│   ├── 001_seed_chart_of_accounts.sql        # 80+ accounts
│   ├── 002_seed_fiscal_year_and_transactions.sql  # Fiscal year setup
│   ├── 003_seed_heavy_transactions.sql       # 182 journal entries
│   ├── 004_seed_ap_ar_data.sql               # 48 bills, 60 invoices
│   ├── 005_seed_fixed_assets.sql             # 15 assets with depreciation
│   ├── 006_seed_odoo_extensions.sql          # Odoo-style features data
│   └── 007_seed_pos_integration.sql          # POS integration mappings
├── schemas/                     # Views and reports
│   ├── accounting_reports.sql               # 8 financial report views
│   └── odoo_reports.sql                    # 10 Odoo-style report views
└── scripts/                     # Setup scripts
    ├── init_accounting.sql                  # Schema initialization
    └── run_all.sh                           # Complete setup script
```

## 🔗 POS Integration

### Account Mapping System

The account mapping system (`pos_account_mappings` table) allows flexible configuration of GL accounts for POS transactions:

```sql
-- Map product to revenue account
INSERT INTO pos_account_mappings (organization_id, source_type, source_id, purpose, account_id)
VALUES ('org-id', 'product', 'product-id', 'revenue', 'revenue-account-id');

-- Map payment method to asset account
INSERT INTO pos_account_mappings (organization_id, source_type, source_code, purpose, account_id)
VALUES ('org-id', 'payment_method', 'CASH', 'asset', 'cash-account-id');

-- Helper function to resolve accounts
SELECT get_pos_gl_account('org-id', 'product', 'product-id', NULL, 'revenue');
```

**Mapping Hierarchy** (highest priority first):
1. Specific product/payment method mapping
2. Category/type default mapping
3. Organization-wide default mapping

### Tax Integration

Bridge POS tax codes with accounting tax definitions:

```sql
-- Map POS tax code to accounting tax
INSERT INTO pos_tax_mappings (
    organization_id, pos_tax_code, tax_category_code,
    accounting_tax_id, default_tax_account_id
) VALUES (
    'org-id', 'VAT_15', 'S', 'tax-id', 'vat-liability-account-id'
);

-- Resolve accounting tax from POS tax code
SELECT get_accounting_tax_for_pos('org-id', 'VAT_15');
```

### Posting Workflow

1. **POS Transaction Created** (e.g., sale)
2. **Go Posting Engine**:
   - Resolve GL accounts using `get_pos_gl_account()`
   - Resolve taxes using `get_accounting_tax_for_pos()`
   - Build journal entry with distribution
3. **Create Journal Entry** → **Post to General Ledger**
4. **Track in Posting Audit**:
   ```sql
   INSERT INTO pos_posting_audit (
       organization_id, source_table, source_id,
       posting_status, journal_entry_id
   ) VALUES ('org-id', 'sales', 'sale-id', 'posted', 'je-id');
   ```
5. **Update Source Document**:
   ```sql
   UPDATE sales SET
       accounting_posting_status = 'posted',
       accounting_journal_entry_id = 'je-id',
       posted_to_accounting_at = CURRENT_TIMESTAMP
   WHERE id = 'sale-id';
   ```

### Inventory Valuation & COGS

**Configuration** (`inventory_valuation_settings`):
```sql
-- Set organization valuation method
INSERT INTO inventory_valuation_settings (
    organization_id, valuation_method, cost_layer_granularity,
    cogs_recognition_timing
) VALUES (
    'org-id', 'fifo', 'product_location', 'on_sale'
);
```

**Cost Layers** (`inventory_cost_layers`):
- Tracks unit cost per batch/lot
- FIFO: oldest layers consumed first
- Weighted average: calculated from all layers

**Views**:
- `view_inventory_valuation_by_product` - Current stock value
- `view_cogs_by_period` - Cost of sales by period
- `view_inventory_turnover` - Turnover ratios
- `view_inventory_reconciliation` - GL vs. valuation comparison

## 🔒 Data Integrity & Controls

### Immutability Rules

**General Ledger**: Completely immutable once posted
```sql
-- Any UPDATE or DELETE on general_ledger will fail
-- Use reversal entries instead
```

**Journal Entries**: Posted entries are immutable
```sql
-- Only status changes and notes allowed on posted entries
-- Financial fields are locked
```

**Fiscal Period Locking**:
```sql
-- Cannot post to closed/locked periods
SELECT close_accounting_period('period-id', 'user-id');
SELECT lock_accounting_period('period-id', 'user-id');
```

### Year-End Closing

```sql
-- Close fiscal year with automated closing entries
SELECT close_fiscal_year('fiscal-year-id', 'user-id', TRUE);
```

**Closing Process**:
1. Verify all periods are closed
2. Calculate net income (Revenue - Expenses)
3. Create closing journal entry
4. Transfer net income to Retained Earnings
5. Mark fiscal year as closed

**Reopen if needed**:
```sql
SELECT reopen_accounting_period('period-id', 'user-id');
```

### Reversal Entries

Never modify posted transactions. Instead, create reversals:

```sql
-- Manual reversal pattern
INSERT INTO journal_entries (
    organization_id, entry_number, description,
    reference, reversal_of_entry_id
) VALUES (
    'org-id', 'JE-2024-REV-001', 'Reversal of JE-2024-001',
    'ERROR_CORRECTION', 'original-je-id'
);
```

## 🚀 Quick Start

### Prerequisites

1. **PostgreSQL 12+** installed and running
2. **Main POS database** must be initialized first
3. Run the main POS setup:
   ```bash
   cd postgres/scripts
   ./run_all.sh pos_saas postgres
   ```

### Installation

1. **Navigate to accounting scripts directory:**
   ```bash
   cd accounting/scripts
   ```

2. **Run the complete setup:**
   ```bash
   ./run_all.sh pos_saas postgres
   ```

   This script will:
   - ✅ Verify POS database exists
   - ✅ Create accounting schema
   - ✅ Run migrations (V001, V002)
   - ✅ Create report views (8 views)
   - ✅ Load seed data (5 files with heavy test data)

3. **Verify installation:**
   ```bash
   psql -h localhost -U postgres -d pos_saas -c "SELECT fiscal_year, COUNT(*) FROM accounting.journal_entries GROUP BY fiscal_year;"
   ```

   Expected output:
   ```
    fiscal_year | count
   -------------+-------
    2024        |  182+
   ```

## 📈 Test Data Summary

The seed data provides a **complete fiscal year 2024** with realistic transactions:

| Data Type | Count | Description |
|-----------|-------|-------------|
| **Chart of Accounts** | 80+ | Standard accounts (Assets, Liabilities, Equity, Revenue, Expenses) |
| **Journal Entries** | 182+ | Full year of realistic transactions |
| **Vendor Bills** | 48 | Inventory purchases with payment terms |
| **Vendor Payments** | ~40 | Bill payments with application tracking |
| **Customer Invoices** | 60 | Credit sales with payment terms |
| **Customer Payments** | ~50 | Invoice payments with application tracking |
| **Fixed Assets** | 15 | Furniture, computers, equipment, vehicles, buildings |
| **Depreciation Entries** | 168 | 12 months × 14 depreciable assets |

### Transaction Types Included:

1. **Monthly Sales & COGS** (24 entries)
   - Revenue recognition with proper COGS matching
   - Seasonal variations (higher in Nov/Dec)

2. **Bi-weekly Payroll** (26 entries)
   - Gross wages + employer taxes
   - Realistic payroll amounts

3. **Recurring Expenses** (36 entries)
   - Rent (monthly)
   - Utilities (monthly with seasonal variations)
   - Insurance (quarterly)

4. **Monthly Depreciation** (12 entries)
   - Automated depreciation for all assets

5. **Inventory Purchases** (24 entries)
   - Bi-monthly inventory purchases on account

6. **Vendor Payments** (24 entries)
   - Bi-monthly AP payments reducing payables

7. **Customer Receipts** (24 entries)
   - Bi-monthly AR collections reducing receivables

8. **Marketing Expenses** (12 entries)
   - Monthly advertising with holiday season increases

## 📊 Financial Reports

### Available Report Views

All reports are accessible via SQL views in the `accounting` schema:

#### 1. Trial Balance
```sql
SELECT * FROM accounting.view_trial_balance
WHERE organization_id = 'your-org-id';
```

Shows all accounts with debit/credit balances to verify books are balanced.

#### 2. Balance Sheet
```sql
SELECT * FROM accounting.view_balance_sheet
WHERE organization_id = 'your-org-id';
```

Statement of Financial Position (Assets = Liabilities + Equity).

#### 3. Income Statement
```sql
SELECT * FROM accounting.view_income_statement
WHERE organization_id = 'your-org-id'
  AND fiscal_year = '2024';
```

Profit & Loss showing Revenue, COGS, Expenses, and Net Income.

#### 4. Cash Flow Statement
```sql
SELECT * FROM accounting.view_cash_flow
WHERE organization_id = 'your-org-id'
  AND fiscal_year = '2024';
```

Cash movements from Operating, Investing, and Financing activities.

#### 5. Account Activity
```sql
SELECT * FROM accounting.view_account_activity
WHERE organization_id = 'your-org-id'
  AND account_code = '1020'  -- Checking account
ORDER BY transaction_date;
```

Detailed transaction history with running balance for any account.

#### 6. Aged Accounts Payable
```sql
SELECT * FROM accounting.view_aged_accounts_payable
WHERE organization_id = 'your-org-id';
```

Outstanding vendor bills by aging bucket (Current, 1-30, 31-60, 61-90, 90+).

#### 7. Aged Accounts Receivable
```sql
SELECT * FROM accounting.view_aged_accounts_receivable
WHERE organization_id = 'your-org-id';
```

Outstanding customer invoices by aging bucket.

#### 8. Financial Ratios
```sql
SELECT * FROM accounting.view_financial_ratios
WHERE organization_id = 'your-org-id'
  AND fiscal_year = '2024';
```

Key performance metrics:
- **Liquidity:** Current Ratio
- **Profitability:** Net Profit Margin, Gross Profit Margin, ROE, ROA
- **Leverage:** Debt-to-Assets, Debt-to-Equity

## 💡 Usage Examples

### Create a New Journal Entry

```sql
-- 1. Create journal entry header
INSERT INTO accounting.journal_entries (
    id, organization_id, entry_number, entry_date, posting_date,
    fiscal_year_id, accounting_period_id, journal_entry_type_id,
    status, total_debit, total_credit, description, created_by
)
VALUES (
    gen_random_uuid(),
    'your-org-id',
    'JE-2024-999',
    '2024-12-15',
    '2024-12-15',
    'fiscal-year-id',
    'period-id',
    'entry-type-id',
    'draft',
    5000.00,
    5000.00,
    'Equipment purchase',
    'user-id'
);

-- 2. Add journal entry lines (must balance: debit = credit)
INSERT INTO accounting.journal_entry_lines (id, journal_entry_id, line_number, account_id, debit_amount, credit_amount, description)
VALUES
    (gen_random_uuid(), 'je-id', 1, 'equipment-account-id', 5000.00, 0, 'New office equipment'),
    (gen_random_uuid(), 'je-id', 2, 'cash-account-id', 0, 5000.00, 'Payment from checking');

-- 3. Post to general ledger (after validation)
-- This happens automatically via application logic or stored procedure
```

### Record a Vendor Bill

```sql
INSERT INTO accounting.vendor_bills (
    id, organization_id, bill_number, supplier_id,
    bill_date, due_date, payment_terms,
    subtotal_amount, tax_amount, total_amount,
    balance_due, status, created_by
)
VALUES (
    gen_random_uuid(),
    'your-org-id',
    'VB-2024-999',
    'supplier-id',
    '2024-12-15',
    '2024-01-14',
    'Net 30',
    10000.00,
    800.00,
    10800.00,
    10800.00,
    'unpaid',
    'user-id'
);
```

### Generate Financial Statements

```sql
-- Income Statement Summary
SELECT
    section_name,
    SUM(amount) as total
FROM accounting.view_income_statement
WHERE organization_id = 'your-org-id'
  AND fiscal_year = '2024'
GROUP BY section_name, CASE section WHEN 'REVENUE' THEN 1 WHEN 'COGS' THEN 2 WHEN 'EXPENSE' THEN 3 END
ORDER BY 2;

-- Result:
--  section_name  |   total
-- ---------------+-----------
--  Revenue       | 2,020,000
--  COGS          | 1,212,000
--  Expenses      |   650,000
--
-- Gross Profit = 808,000 (40%)
-- Net Income = 158,000 (7.8%)
```

## 🔒 Security (RLS)

All accounting tables have **Row-Level Security (RLS)** policies enabled:

```sql
-- Users can only see data for their organization
CREATE POLICY org_isolation ON accounting.journal_entries
    FOR ALL
    USING (organization_id IN (
        SELECT organization_id FROM user_organizations
        WHERE user_id = auth.uid()
    ));
```

This ensures multi-tenant data isolation at the database level.

## 🎓 Accounting Concepts

### Double-Entry Bookkeeping

Every transaction has **equal debits and credits**:

```
Debit = Credit (ALWAYS)
```

### Account Types & Normal Balances

| Account Type | Normal Balance | Increases With | Decreases With |
|--------------|----------------|----------------|----------------|
| ASSET | Debit | Debit | Credit |
| LIABILITY | Credit | Credit | Debit |
| EQUITY | Credit | Credit | Debit |
| REVENUE | Credit | Credit | Debit |
| EXPENSE | Debit | Debit | Credit |
| COGS | Debit | Debit | Credit |

### Financial Statement Equation

```
Assets = Liabilities + Equity

Revenue - COGS = Gross Profit
Gross Profit - Expenses = Net Income
```

## 🔧 Customization

### Adding Custom Accounts

```sql
INSERT INTO accounting.chart_of_accounts (
    id, organization_id, account_code, account_number, account_name,
    account_type_id, account_level, is_active
)
SELECT
    gen_random_uuid(),
    'your-org-id',
    '6999',
    '6999',
    'Custom Expense Account',
    id,
    1,
    true
FROM accounting.account_types
WHERE type_code = 'EXPENSE';
```

### Creating Custom Reports

```sql
CREATE VIEW accounting.view_custom_report AS
SELECT
    coa.account_name,
    SUM(gl.debit_amount - gl.credit_amount) as balance
FROM accounting.general_ledger gl
JOIN accounting.chart_of_accounts coa ON gl.account_id = coa.id
WHERE coa.account_code LIKE '6%'  -- Expenses only
GROUP BY coa.account_name;
```

## 📚 Standards & References

This accounting module follows standards from:

- **QuickBooks** - Chart of accounts structure
- **Xero** - Multi-organization architecture
- **Sage** - Account classifications and reporting
- **NetSuite** - Multi-period fiscal year design
- **GAAP** - Generally Accepted Accounting Principles (US)

## 🐛 Troubleshooting

### Common Issues

**1. "Main POS tables not found"**
- Solution: Run the main POS setup first: `cd ../postgres/scripts && ./run_all.sh`

**2. "Journal entry out of balance"**
- Solution: Ensure `total_debit = total_credit` before posting

**3. "Account not found"**
- Solution: Verify account exists in chart_of_accounts for your organization

**4. "Period is closed"**
- Solution: Journal entries cannot be posted to closed periods. Use adjusting entries or reopen period.

### Validation Queries

```sql
-- Verify books are balanced
SELECT SUM(debit_amount) - SUM(credit_amount) as balance
FROM accounting.general_ledger
WHERE organization_id = 'your-org-id';
-- Result should be 0.00

-- Check for unposted journal entries
SELECT COUNT(*) FROM accounting.journal_entries
WHERE organization_id = 'your-org-id' AND status = 'draft';

-- Verify trial balance
SELECT
    SUM(CASE WHEN normal_balance = 'debit' THEN balance ELSE 0 END) as total_debits,
    SUM(CASE WHEN normal_balance = 'credit' THEN balance ELSE 0 END) as total_credits
FROM accounting.view_trial_balance
WHERE organization_id = 'your-org-id';
-- Debits should equal Credits
```

## 🤝 Integration with POS

The accounting module integrates seamlessly with POS operations:

### Automatic Journal Entries

Your POS application can automatically create journal entries for:

1. **Sales Transactions** → Revenue & COGS entries
2. **Inventory Purchases** → Inventory & AP entries
3. **Customer Payments** → Cash & AR entries
4. **Vendor Payments** → AP & Cash entries
5. **Payroll** → Expense & Cash entries

### Example Integration

```typescript
// When a POS sale is completed
async function recordSale(sale: Sale) {
    // 1. Create POS transaction (in POS tables)
    const transaction = await createTransaction(sale);

    // 2. Create accounting journal entry (in accounting tables)
    const journalEntry = await createJournalEntry({
        entryDate: sale.date,
        description: `Sale #${sale.id}`,
        lines: [
            { account: 'Cash', debit: sale.total },
            { account: 'Sales Revenue', credit: sale.subtotal },
            { account: 'Sales Tax Payable', credit: sale.tax },
            { account: 'COGS', debit: sale.cogs },
            { account: 'Inventory', credit: sale.cogs }
        ]
    });

    // 3. Link POS transaction to journal entry
    await linkTransactionToJE(transaction.id, journalEntry.id);
}
```

## 📞 Support

For issues, questions, or contributions:
- Check the [COMPREHENSIVE_SYSTEM_REVIEW.md](../COMPREHENSIVE_SYSTEM_REVIEW.md) for system architecture
- Review [SCHEMA_DOCUMENTATION.md](../postgres/schemas/SCHEMA_DOCUMENTATION.md) for table structures

## 📝 License

This accounting module is part of the POS SaaS platform and follows the same license.

---

**Built with ❤️ for accurate, reliable financial management**


---

## 🚀 Odoo-Style Advanced Features (V003)

The accounting module has been extended with **9 advanced modules** to match Odoo-grade enterprise accounting capabilities.

### 1. 📖 Journals

**Purpose**: Proper journal layer for organizing and categorizing accounting entries.

**Features**:
- **6 Journal Types**: Sales, Purchase, Bank, Cash, General, Miscellaneous
- **Automatic Numbering**: Each journal has its own sequence (SAJ-001, PUR-001, etc.)
- **Default Accounts**: Pre-configured debit/credit accounts per journal
- **Bank Integration**: Bank journals linked to specific bank accounts

**Tables**: `journals`

**Usage**:
```sql
-- View all journals
SELECT * FROM accounting.journals WHERE organization_id = 'your-org-id';

-- Journal entries by journal type
SELECT * FROM accounting.view_journal_entries_by_journal
WHERE journal_type = 'sale';
```

---

### 2. 💰 Tax Engine

**Purpose**: Structured tax management with rates, groups, and fiscal positions.

**Features**:
- **Tax Groups**: VAT, Sales Tax, Withholding Tax
- **Tax Rates**: 0%, 5%, 8%, 10%, 20% (configurable)
- **Tax Scope**: Sales, Purchases, or Both
- **Price Inclusive/Exclusive**: Support for gross vs net pricing
- **Fiscal Positions**: Tax remapping for regions (Domestic, Export, EU B2B)
- **Tax Mappings**: Automatic tax substitution (e.g., VAT 20% → VAT 0% for exports)

**Tables**: `tax_groups`, `taxes`, `fiscal_positions`, `fiscal_position_tax_mappings`

**Usage**:
```sql
-- View all taxes
SELECT * FROM accounting.taxes WHERE organization_id = 'your-org-id';

-- Tax collected vs paid report
SELECT * FROM accounting.view_tax_report
WHERE fiscal_year = 2024;

-- Apply fiscal position to customer
UPDATE customers
SET fiscal_position_id = (SELECT id FROM fiscal_positions WHERE position_code = 'EXPORT')
WHERE customer_id = 'your-customer-id';
```

---

### 3. 🌍 Multi-Currency Support

**Purpose**: Handle transactions in multiple currencies with exchange rate management.

**Features**:
- **10 Major Currencies**: USD, EUR, GBP, JPY, CAD, AUD, CHF, CNY, INR, MXN
- **Daily Exchange Rates**: Historical rates for accurate conversion
- **Transaction Currency**: Record amounts in both transaction and base currency
- **FX Gain/Loss**: Automatic calculation of unrealized FX differences

**Tables**: `currencies`, `currency_rates`

**Extended Fields**:
- `journal_entries`: `currency_code`, `exchange_rate`
- `journal_entry_lines`: `amount_currency`, `currency_code`
- `customer_invoices` / `vendor_bills`: `currency_code`, `exchange_rate`

**Usage**:
```sql
-- Add exchange rates
INSERT INTO accounting.currency_rates (organization_id, currency_code, rate_date, rate, source)
VALUES ('org-id', 'EUR', '2024-12-01', 1.0950, 'manual');

-- View multi-currency transactions
SELECT * FROM accounting.view_multi_currency_summary
WHERE transaction_currency != 'USD';
```

---

### 4. 📅 Payment Terms & Schedules

**Purpose**: Structured payment terms with automatic due date calculation.

**Features**:
- **7 Standard Terms**: Immediate, Net 15/30/60, 2/10 Net 30, 50/50 Split, End of Month
- **Multi-line Terms**: Split payments (e.g., 50% now, 50% in 30 days)
- **Auto-Calculation**: Due dates computed from invoice date + term rules
- **Payment Schedules**: Track partial payments against specific due dates

**Tables**: `payment_terms`, `payment_term_lines`, `invoice_payment_schedules`

**Usage**:
```sql
-- View payment terms
SELECT * FROM accounting.payment_terms;

-- Upcoming payment schedules
SELECT * FROM accounting.view_payment_schedules
WHERE status != 'paid'
ORDER BY due_date;
```

---

### 5. 📊 Analytic Accounting

**Purpose**: Multi-dimensional cost tracking (projects, departments, regions).

**Features**:
- **3 Analytic Plans**: Projects, Departments, Regions
- **14 Analytic Accounts**: 
  - Projects: Website Redesign, Mobile App, Store Expansion
  - Departments: Sales, Marketing, Operations, IT
  - Regions: North, South, East, West
- **Hierarchical Structure**: Parent-child relationships for rollups
- **Tag Journal Lines**: Every expense can be tagged to project/department

**Tables**: `analytic_plans`, `analytic_accounts`

**Extended Fields**:
- `journal_entry_lines`: `analytic_account_id`
- `customer_invoice_items` / `vendor_bill_items`: `analytic_account_id`

**Usage**:
```sql
-- View analytic accounts
SELECT * FROM accounting.analytic_accounts
WHERE analytic_plan_id = (SELECT id FROM analytic_plans WHERE plan_code = 'PROJ');

-- Project cost report
SELECT * FROM accounting.view_analytic_report
WHERE analytic_code LIKE 'PROJ-%'
ORDER BY net_amount DESC;
```

---

### 6. 🔄 Deferred Revenue & Expense

**Purpose**: Revenue/expense recognition over time (subscriptions, prepayments).

**Features**:
- **Deferred Revenue**: Annual subscriptions, warranties, maintenance contracts
- **Deferred Expense**: Prepaid insurance, rent, software licenses
- **Recognition Methods**: Straight-line, custom, milestone-based
- **Automated Schedules**: Monthly recognition entries
- **Status Tracking**: Pending, Posted, Completed

**Tables**: 
- `deferred_revenue_contracts`, `deferred_revenue_schedule`
- `deferred_expense_contracts`, `deferred_expense_schedule`

**Usage**:
```sql
-- View active deferrals
SELECT * FROM accounting.view_deferrals_report
WHERE status = 'active';

-- Pending recognition entries
SELECT * FROM accounting.deferred_revenue_schedule
WHERE status = 'pending'
  AND recognition_date <= CURRENT_DATE + INTERVAL '7 days';
```

---

### 7. 🏦 Bank Statements & Reconciliation

**Purpose**: Import bank statements and match to GL transactions.

**Features**:
- **Bank Statement Import**: Manual, file import, API, bank feed
- **Statement Lines**: Individual bank transactions with details
- **Auto-Matching Rules**: Pattern-based reconciliation
- **Match Status**: Matched, Unmatched, Partial Match, Ignored
- **Reconciliation Links**: Connect statement lines to journal entries/payments

**Tables**: `bank_statements`, `bank_statement_lines`, `bank_statement_reconciliations`, `reconciliation_rule_models`

**Usage**:
```sql
-- View bank statements
SELECT * FROM accounting.view_bank_reconciliation_status;

-- Unmatched transactions
SELECT * FROM accounting.bank_statement_lines
WHERE status = 'unmatched'
ORDER BY transaction_date DESC;

-- Create reconciliation rule
INSERT INTO accounting.reconciliation_rule_models
(organization_id, rule_name, description_pattern, account_id, auto_apply)
VALUES ('org-id', 'Rent Payments', '%PROPERTY MGMT%', 'rent-expense-account-id', true);
```

---

### 8. 📈 Budgets & Budget Analysis

**Purpose**: Planning and budget vs actual comparison.

**Features**:
- **Budget Types**: Operating, Capital, Cash Flow, Project, Departmental
- **Multi-Level Budgets**: Account-level and analytic-level budgets
- **Period Breakdown**: Annual, quarterly, monthly granularity
- **Budget Status**: Draft, Approved, Active, Closed
- **Variance Analysis**: Automatic calculation of budget variance

**Tables**: `budgets`, `budget_lines`

**Usage**:
```sql
-- View budgets
SELECT * FROM accounting.budgets
WHERE fiscal_year_id = (SELECT id FROM fiscal_years WHERE fiscal_year = '2024');

-- Budget vs actual analysis
SELECT * FROM accounting.view_budget_vs_actual
WHERE budget_code = 'FY2024-OP'
  AND status = 'Over Budget';
```

---

### 9. 🌐 Localization & Tax Reporting

**Purpose**: Country-specific accounting and tax compliance.

**Features**:
- **5 Localization Packages**: US GAAP, UK VAT, EU IFRS, Canada GAAP, Australia AAS
- **Tax Report Definitions**: Declarative tax reports (VAT returns, sales tax)
- **Report Lines**: Configurable formulas and mappings
- **Organization Assignment**: Each org can use a specific localization

**Tables**: `localization_packages`, `tax_report_definitions`, `tax_report_lines`

**Extended Fields**:
- `organizations`: `localization_package_id`

**Usage**:
```sql
-- View localization packages
SELECT * FROM accounting.view_localization_summary;

-- Set organization localization
UPDATE organizations
SET localization_package_id = (SELECT id FROM localization_packages WHERE package_code = 'us_gaap')
WHERE id = 'your-org-id';

-- View tax reports
SELECT * FROM accounting.tax_report_definitions
WHERE organization_id = 'your-org-id';
```

---

## 📊 Extended Reporting (18 Total Views)

### Core Reports (from V001/V002)
1. Trial Balance
2. Balance Sheet
3. Income Statement
4. Cash Flow Statement
5. Account Activity
6. Aged AP
7. Aged AR
8. Financial Ratios

### Advanced Reports (from V003)
9. **Journal Entries by Journal** - Entries grouped by journal type
10. **Tax Report** - Tax collected vs paid by period
11. **Multi-Currency Summary** - FX transactions with unrealized gains/losses
12. **Payment Schedules** - Upcoming payments with aging
13. **Analytic Report** - Costs by project/department/region
14. **Deferrals Report** - Revenue/expense recognition schedules
15. **Bank Reconciliation Status** - Statement matching progress
16. **Budget vs Actual** - Variance analysis with percentages
17. **Fiscal Position Usage** - Tax remapping statistics
18. **Localization Summary** - Package usage by organization

---

## 📁 Updated Directory Structure

```
accounting/
├── README.md                                    # This file
├── migrations/                                  # Database migrations (DDL)
│   ├── V001_20251109_create_accounting_core.sql
│   ├── V002_20251109_create_ap_ar_assets.sql
│   └── V003_20251110_create_odoo_extensions.sql  # NEW: Odoo-style features
├── seed_data/                                   # Test data (DML)
│   ├── 001_seed_chart_of_accounts.sql
│   ├── 002_seed_fiscal_year_and_transactions.sql
│   ├── 003_seed_heavy_transactions.sql
│   ├── 004_seed_ap_ar_data.sql
│   ├── 005_seed_fixed_assets.sql
│   └── 006_seed_odoo_extensions.sql             # NEW: Odoo seed data
├── schemas/                                     # Views and reports
│   ├── accounting_reports.sql                   # Core 8 reports
│   └── odoo_reports.sql                         # NEW: Advanced 10 reports
└── scripts/                                     # Setup scripts
    ├── init_accounting.sql
    └── run_all.sh                               # Auto-runs all migrations + seeds
```

---

## 🆕 What's New in V003

### Database Changes
- **+30 New Tables**: Journals, Taxes, Currencies, Payment Terms, Analytics, Deferrals, Bank Statements, Budgets, Localization
- **+10 Extended Tables**: Added fields to existing tables (journal_id, tax_id, currency_code, etc.)
- **+10 New Reports**: Advanced analytical views

### Seed Data
- **6 Journals**: Sales, Purchase, Bank, Cash, General, Miscellaneous
- **6 Taxes** with 3 tax groups and 3 fiscal positions
- **60 Currency Rates** (5 currencies × 12 months)
- **7 Payment Terms** with 15+ term lines
- **14 Analytic Accounts** across 3 dimensions
- **3 Deferral Contracts** with 34 schedules
- **2 Bank Statements** with 15 transactions
- **2 Budgets** with 38 budget lines
- **1 Tax Report Definition** with 5 lines

---

## 🔄 Migration Path

If you already have the basic accounting module (V001/V002) installed:

```bash
cd accounting/scripts

# Run V003 migration
psql -h localhost -U postgres -d pos_saas -f ../migrations/V003_20251110_create_odoo_extensions.sql

# Load V003 seed data
psql -h localhost -U postgres -d pos_saas -f ../seed_data/006_seed_odoo_extensions.sql

# Create advanced reports
psql -h localhost -U postgres -d pos_saas -f ../schemas/odoo_reports.sql
```

Or run the complete setup (includes all versions):

```bash
cd accounting/scripts
./run_all.sh pos_saas postgres
```

---

## 🎓 Odoo Comparison

This accounting module now matches Odoo's accounting capabilities:

| Feature | Odoo | This Module | Status |
|---------|------|-------------|--------|
| Journals | ✅ | ✅ | Full parity |
| Tax Engine | ✅ | ✅ | Full parity |
| Multi-Currency | ✅ | ✅ | Full parity |
| Payment Terms | ✅ | ✅ | Full parity |
| Analytic Accounting | ✅ | ✅ | Full parity |
| Deferred Revenue | ✅ | ✅ | Full parity |
| Bank Reconciliation | ✅ | ✅ | Full parity |
| Budgets | ✅ | ✅ | Full parity |
| Localization | ✅ | ✅ | Full parity |
| **Total Tables** | ~40 | **45** | ✅ More comprehensive |
| **Report Views** | ~15 | **18** | ✅ More comprehensive |

---

## 💡 Advanced Usage Examples

### Multi-Currency Invoice

```sql
-- Create invoice in EUR
INSERT INTO customer_invoices (
    organization_id, customer_id, invoice_date, currency_code, exchange_rate,
    subtotal_amount, tax_amount, total_amount
)
SELECT 
    'org-id', 'customer-id', '2024-12-01', 'EUR', 
    (SELECT rate FROM currency_rates WHERE currency_code = 'EUR' AND rate_date = '2024-12-01'),
    1000.00, 200.00, 1200.00;
```

### Project Cost Tracking

```sql
-- Record expense to specific project
INSERT INTO journal_entry_lines (
    journal_entry_id, account_id, analytic_account_id, debit_amount, description
)
VALUES (
    'je-id',
    (SELECT id FROM chart_of_accounts WHERE account_code = '6010'),
    (SELECT id FROM analytic_accounts WHERE account_code = 'PROJ-001'),
    5000.00,
    'Development costs for website redesign'
);
```

### Bank Reconciliation

```sql
-- Match statement line to journal entry
INSERT INTO bank_statement_reconciliations (
    organization_id, bank_statement_line_id, journal_entry_id, matched_amount
)
VALUES (
    'org-id', 'stmt-line-id', 'je-id', 45000.00
);

-- Update line status
UPDATE bank_statement_lines
SET status = 'matched'
WHERE id = 'stmt-line-id';
```

### Budget Monitoring

```sql
-- Check departments over budget
SELECT 
    analytic_name,
    budget_amount,
    actual_amount,
    variance,
    percentage_of_budget
FROM accounting.view_budget_vs_actual
WHERE budget_code = 'FY2024-DEPT'
  AND status = 'Over Budget'
ORDER BY variance;
```

---

## 🔐 Security & Multi-Tenancy

All new tables include:
- **Row-Level Security (RLS)** policies
- **Organization-based isolation** (multi-tenant safe)
- **Soft deletes** (deleted_at timestamp)
- **Audit trails** (created_by, updated_by, timestamps)

---

## 📞 Support & Documentation

For complete system documentation, see:
- [COMPREHENSIVE_SYSTEM_REVIEW.md](../COMPREHENSIVE_SYSTEM_REVIEW.md) - System architecture
- [SCHEMA_DOCUMENTATION.md](../postgres/schemas/SCHEMA_DOCUMENTATION.md) - POS table structures

---

**🎉 You now have a production-grade, Odoo-style accounting system with 45+ tables, 18 financial reports, and comprehensive multi-currency, multi-dimensional, multi-tenant capabilities!**


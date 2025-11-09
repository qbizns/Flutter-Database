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
│   └── V002_20251109_create_ap_ar_assets.sql
├── seed_data/                   # Test data (DML)
│   ├── 001_seed_chart_of_accounts.sql        # 80+ accounts
│   ├── 002_seed_fiscal_year_and_transactions.sql  # Fiscal year setup
│   ├── 003_seed_heavy_transactions.sql       # 182 journal entries
│   ├── 004_seed_ap_ar_data.sql               # 48 bills, 60 invoices
│   └── 005_seed_fixed_assets.sql             # 15 assets with depreciation
├── schemas/                     # Views and reports
│   └── accounting_reports.sql               # 8 financial report views
└── scripts/                     # Setup scripts
    ├── init_accounting.sql                  # Schema initialization
    └── run_all.sh                           # Complete setup script
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


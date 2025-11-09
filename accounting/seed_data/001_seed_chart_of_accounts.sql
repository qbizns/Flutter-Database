-- =====================================================
-- Accounting Seed Data: Standard Chart of Accounts
-- Description: Complete COA based on QuickBooks/Xero standards
-- Includes: Assets, Liabilities, Equity, Revenue, Expenses
-- =====================================================

\echo 'Loading Chart of Accounts seed data...'

-- =====================================================
-- SECTION 1: Account Types (Standard Classifications)
-- =====================================================

INSERT INTO account_types (id, type_code, type_name, type_category, normal_balance, is_balance_sheet, is_income_statement, display_order) VALUES
('11111111-1111-1111-1111-111111111111', 'ASSET', 'Assets', 'balance_sheet', 'debit', true, false, 1),
('22222222-2222-2222-2222-222222222222', 'LIABILITY', 'Liabilities', 'balance_sheet', 'credit', true, false, 2),
('33333333-3333-3333-3333-333333333333', 'EQUITY', 'Equity', 'balance_sheet', 'credit', true, false, 3),
('44444444-4444-4444-4444-444444444444', 'REVENUE', 'Revenue', 'income_statement', 'credit', false, true, 4),
('55555555-5555-5555-5555-555555555555', 'EXPENSE', 'Expenses', 'income_statement', 'debit', false, true, 5),
('66666666-6666-6666-6666-666666666666', 'COGS', 'Cost of Goods Sold', 'income_statement', 'debit', false, true, 6);

-- =====================================================
-- SECTION 2: Account Subtypes
-- =====================================================

INSERT INTO account_subtypes (account_type_id, subtype_code, subtype_name, display_order) VALUES
-- Assets
('11111111-1111-1111-1111-111111111111', 'CASH', 'Cash and Cash Equivalents', 1),
('11111111-1111-1111-1111-111111111111', 'AR', 'Accounts Receivable', 2),
('11111111-1111-1111-1111-111111111111', 'INVENTORY', 'Inventory', 3),
('11111111-1111-1111-1111-111111111111', 'PREPAID', 'Prepaid Expenses', 4),
('11111111-1111-1111-1111-111111111111', 'FIXED_ASSET', 'Fixed Assets', 5),
('11111111-1111-1111-1111-111111111111', 'ACC_DEPR', 'Accumulated Depreciation', 6),
('11111111-1111-1111-1111-111111111111', 'OTHER_ASSET', 'Other Assets', 7),

-- Liabilities
('22222222-2222-2222-2222-222222222222', 'AP', 'Accounts Payable', 1),
('22222222-2222-2222-2222-222222222222', 'CC', 'Credit Cards', 2),
('22222222-2222-2222-2222-222222222222', 'SALES_TAX', 'Sales Tax Payable', 3),
('22222222-2222-2222-2222-222222222222', 'PAYROLL', 'Payroll Liabilities', 4),
('22222222-2222-2222-2222-222222222222', 'ST_LOAN', 'Short-term Loans', 5),
('22222222-2222-2222-2222-222222222222', 'LT_LOAN', 'Long-term Loans', 6),
('22222222-2222-2222-2222-222222222222', 'OTHER_LIAB', 'Other Liabilities', 7),

-- Equity
('33333333-3333-3333-3333-333333333333', 'CAPITAL', 'Owner''s Capital', 1),
('33333333-3333-3333-3333-333333333333', 'RETAINED', 'Retained Earnings', 2),
('33333333-3333-3333-3333-333333333333', 'DRAWINGS', 'Owner''s Drawings', 3),

-- Revenue
('44444444-4444-4444-4444-444444444444', 'SALES', 'Sales Revenue', 1),
('44444444-4444-4444-4444-444444444444', 'SERVICE', 'Service Revenue', 2),
('44444444-4444-4444-4444-444444444444', 'OTHER_INC', 'Other Income', 3),

-- COGS
('66666666-6666-6666-6666-666666666666', 'COGS_PRODUCT', 'Cost of Products Sold', 1),
('66666666-6666-6666-6666-666666666666', 'COGS_SHIPPING', 'Shipping Costs', 2),

-- Expenses
('55555555-5555-5555-5555-555555555555', 'PAYROLL_EXP', 'Payroll Expenses', 1),
('55555555-5555-5555-5555-555555555555', 'RENT', 'Rent Expense', 2),
('55555555-5555-5555-5555-555555555555', 'UTILITIES', 'Utilities', 3),
('55555555-5555-5555-5555-555555555555', 'MARKETING', 'Marketing & Advertising', 4),
('55555555-5555-5555-5555-555555555555', 'OFFICE', 'Office Expenses', 5),
('55555555-5555-5555-5555-555555555555', 'DEPRECIATION', 'Depreciation', 6),
('55555555-5555-5555-5555-555555555555', 'INSURANCE', 'Insurance', 7),
('55555555-5555-5555-5555-555555555555', 'BANK_FEES', 'Bank Fees & Charges', 8),
('55555555-5555-5555-5555-555555555555', 'PROFESSIONAL', 'Professional Fees', 9),
('55555555-5555-5555-5555-555555555555', 'OTHER_EXP', 'Other Expenses', 10);

-- =====================================================
-- SECTION 3: Chart of Accounts for Demo Retail Store
-- =====================================================

-- Get organization ID
DO $$
DECLARE
    v_org_id UUID;
    v_asset_type_id UUID := '11111111-1111-1111-1111-111111111111';
    v_liability_type_id UUID := '22222222-2222-2222-2222-222222222222';
    v_equity_type_id UUID := '33333333-3333-3333-3333-333333333333';
    v_revenue_type_id UUID := '44444444-4444-4444-4444-444444444444';
    v_expense_type_id UUID := '55555555-5555-5555-5555-555555555555';
    v_cogs_type_id UUID := '66666666-6666-6666-6666-666666666666';
BEGIN
    -- Get organization
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- ================= ASSETS =================

    -- 1000-1999: Current Assets
    INSERT INTO chart_of_accounts (organization_id, account_code, account_number, account_name, account_type_id, account_level, is_active, is_system_account, opening_balance) VALUES
    (v_org_id, '1000', '1000', 'Cash and Bank Accounts', v_asset_type_id, 1, true, true, 0),
    (v_org_id, '1010', '1010', 'Petty Cash', v_asset_type_id, 2, true, false, 1000.00),
    (v_org_id, '1020', '1020', 'Checking Account - Main', v_asset_type_id, 2, true, true, 150000.00),
    (v_org_id, '1030', '1030', 'Savings Account', v_asset_type_id, 2, true, false, 50000.00),
    (v_org_id, '1040', '1040', 'Merchant Account', v_asset_type_id, 2, true, false, 25000.00),

    (v_org_id, '1200', '1200', 'Accounts Receivable', v_asset_type_id, 1, true, true, 85000.00),
    (v_org_id, '1210', '1210', 'Accounts Receivable - Trade', v_asset_type_id, 2, true, true, 85000.00),

    (v_org_id, '1300', '1300', 'Inventory', v_asset_type_id, 1, true, true, 0),
    (v_org_id, '1310', '1310', 'Inventory - Retail', v_asset_type_id, 2, true, true, 125000.00),
    (v_org_id, '1320', '1320', 'Inventory - Warehouse', v_asset_type_id, 2, true, false, 75000.00),

    (v_org_id, '1400', '1400', 'Prepaid Expenses', v_asset_type_id, 1, true, false, 0),
    (v_org_id, '1410', '1410', 'Prepaid Insurance', v_asset_type_id, 2, true, false, 12000.00),
    (v_org_id, '1420', '1420', 'Prepaid Rent', v_asset_type_id, 2, true, false, 15000.00),

    -- 1500-1999: Fixed Assets
    (v_org_id, '1500', '1500', 'Fixed Assets', v_asset_type_id, 1, true, true, 0),
    (v_org_id, '1510', '1510', 'Furniture & Fixtures', v_asset_type_id, 2, true, false, 45000.00),
    (v_org_id, '1520', '1520', 'Computer Equipment', v_asset_type_id, 2, true, false, 35000.00),
    (v_org_id, '1530', '1530', 'Store Equipment', v_asset_type_id, 2, true, false, 55000.00),
    (v_org_id, '1540', '1540', 'Vehicles', v_asset_type_id, 2, true, false, 60000.00),
    (v_org_id, '1550', '1550', 'Leasehold Improvements', v_asset_type_id, 2, true, false, 80000.00),

    (v_org_id, '1600', '1600', 'Accumulated Depreciation', v_asset_type_id, 1, true, true, -45000.00),
    (v_org_id, '1610', '1610', 'Accumulated Depreciation - Furniture', v_asset_type_id, 2, true, false, -10000.00),
    (v_org_id, '1620', '1620', 'Accumulated Depreciation - Computers', v_asset_type_id, 2, true, false, -15000.00),
    (v_org_id, '1630', '1630', 'Accumulated Depreciation - Store Equipment', v_asset_type_id, 2, true, false, -12000.00),
    (v_org_id, '1640', '1640', 'Accumulated Depreciation - Vehicles', v_asset_type_id, 2, true, false, -8000.00);

    -- ================= LIABILITIES =================

    -- 2000-2499: Current Liabilities
    INSERT INTO chart_of_accounts (organization_id, account_code, account_number, account_name, account_type_id, account_level, is_active, is_system_account, opening_balance) VALUES
    (v_org_id, '2000', '2000', 'Accounts Payable', v_liability_type_id, 1, true, true, 0),
    (v_org_id, '2010', '2010', 'Accounts Payable - Trade', v_liability_type_id, 2, true, true, 65000.00),

    (v_org_id, '2100', '2100', 'Credit Cards', v_liability_type_id, 1, true, false, 0),
    (v_org_id, '2110', '2110', 'Business Credit Card', v_liability_type_id, 2, true, false, 12000.00),

    (v_org_id, '2200', '2200', 'Sales Tax Payable', v_liability_type_id, 1, true, true, 0),
    (v_org_id, '2210', '2210', 'Sales Tax Payable - Local', v_liability_type_id, 2, true, true, 8500.00),

    (v_org_id, '2300', '2300', 'Payroll Liabilities', v_liability_type_id, 1, true, true, 0),
    (v_org_id, '2310', '2310', 'Salaries Payable', v_liability_type_id, 2, true, false, 18000.00),
    (v_org_id, '2320', '2320', 'Payroll Tax Payable', v_liability_type_id, 2, true, false, 5400.00),
    (v_org_id, '2330', '2330', 'Employee Benefits Payable', v_liability_type_id, 2, true, false, 2400.00),

    (v_org_id, '2400', '2400', 'Short-term Loans', v_liability_type_id, 1, true, false, 0),
    (v_org_id, '2410', '2410', 'Bank Loan - Current Portion', v_liability_type_id, 2, true, false, 15000.00),

    -- 2500-2999: Long-term Liabilities
    (v_org_id, '2500', '2500', 'Long-term Loans', v_liability_type_id, 1, true, false, 0),
    (v_org_id, '2510', '2510', 'Bank Loan - Long-term', v_liability_type_id, 2, true, false, 85000.00),
    (v_org_id, '2520', '2520', 'Equipment Financing', v_liability_type_id, 2, true, false, 42000.00);

    -- ================= EQUITY =================

    INSERT INTO chart_of_accounts (organization_id, account_code, account_number, account_name, account_type_id, account_level, is_active, is_system_account, opening_balance) VALUES
    (v_org_id, '3000', '3000', 'Owner''s Equity', v_equity_type_id, 1, true, true, 0),
    (v_org_id, '3010', '3010', 'Owner''s Capital', v_equity_type_id, 2, true, true, 500000.00),
    (v_org_id, '3020', '3020', 'Retained Earnings', v_equity_type_id, 2, true, true, 150000.00),
    (v_org_id, '3030', '3030', 'Owner''s Drawings', v_equity_type_id, 2, true, false, 0),
    (v_org_id, '3900', '3900', 'Current Year Earnings', v_equity_type_id, 2, true, true, 0);

    -- ================= REVENUE =================

    INSERT INTO chart_of_accounts (organization_id, account_code, account_number, account_name, account_type_id, account_level, is_active, is_system_account, opening_balance) VALUES
    (v_org_id, '4000', '4000', 'Sales Revenue', v_revenue_type_id, 1, true, true, 0),
    (v_org_id, '4010', '4010', 'Product Sales - Electronics', v_revenue_type_id, 2, true, false, 0),
    (v_org_id, '4020', '4020', 'Product Sales - Clothing', v_revenue_type_id, 2, true, false, 0),
    (v_org_id, '4030', '4030', 'Product Sales - Home & Garden', v_revenue_type_id, 2, true, false, 0),
    (v_org_id, '4040', '4040', 'Product Sales - Food & Beverage', v_revenue_type_id, 2, true, false, 0),
    (v_org_id, '4050', '4050', 'Service Revenue', v_revenue_type_id, 2, true, false, 0),

    (v_org_id, '4100', '4100', 'Sales Discounts', v_revenue_type_id, 1, true, false, 0),
    (v_org_id, '4110', '4110', 'Sales Discounts Given', v_revenue_type_id, 2, true, false, 0),

    (v_org_id, '4900', '4900', 'Other Income', v_revenue_type_id, 1, true, false, 0),
    (v_org_id, '4910', '4910', 'Interest Income', v_revenue_type_id, 2, true, false, 0),
    (v_org_id, '4920', '4920', 'Miscellaneous Income', v_revenue_type_id, 2, true, false, 0);

    -- ================= COST OF GOODS SOLD =================

    INSERT INTO chart_of_accounts (organization_id, account_code, account_number, account_name, account_type_id, account_level, is_active, is_system_account, opening_balance) VALUES
    (v_org_id, '5000', '5000', 'Cost of Goods Sold', v_cogs_type_id, 1, true, true, 0),
    (v_org_id, '5010', '5010', 'COGS - Electronics', v_cogs_type_id, 2, true, false, 0),
    (v_org_id, '5020', '5020', 'COGS - Clothing', v_cogs_type_id, 2, true, false, 0),
    (v_org_id, '5030', '5030', 'COGS - Home & Garden', v_cogs_type_id, 2, true, false, 0),
    (v_org_id, '5040', '5040', 'COGS - Food & Beverage', v_cogs_type_id, 2, true, false, 0),
    (v_org_id, '5100', '5100', 'Freight & Shipping Costs', v_cogs_type_id, 2, true, false, 0);

    -- ================= EXPENSES =================

    -- Operating Expenses
    INSERT INTO chart_of_accounts (organization_id, account_code, account_number, account_name, account_type_id, account_level, is_active, is_system_account, opening_balance) VALUES
    (v_org_id, '6000', '6000', 'Payroll Expenses', v_expense_type_id, 1, true, true, 0),
    (v_org_id, '6010', '6010', 'Salaries & Wages', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6020', '6020', 'Payroll Taxes', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6030', '6030', 'Employee Benefits', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6040', '6040', 'Health Insurance', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6050', '6050', 'Retirement Contributions', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6100', '6100', 'Occupancy Expenses', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6110', '6110', 'Rent Expense', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6120', '6120', 'Property Tax', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6130', '6130', 'Building Maintenance', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6200', '6200', 'Utilities', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6210', '6210', 'Electricity', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6220', '6220', 'Water & Sewer', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6230', '6230', 'Gas & Heating', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6240', '6240', 'Internet & Phone', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6300', '6300', 'Marketing & Advertising', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6310', '6310', 'Online Advertising', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6320', '6320', 'Print Advertising', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6330', '6330', 'Promotional Materials', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6340', '6340', 'Website & SEO', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6400', '6400', 'Office Expenses', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6410', '6410', 'Office Supplies', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6420', '6420', 'Postage & Shipping', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6430', '6430', 'Software Subscriptions', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6500', '6500', 'Professional Fees', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6510', '6510', 'Legal Fees', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6520', '6520', 'Accounting Fees', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6530', '6530', 'Consulting Fees', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6600', '6600', 'Insurance', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6610', '6610', 'General Liability Insurance', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6620', '6620', 'Property Insurance', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6630', '6630', 'Vehicle Insurance', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6700', '6700', 'Vehicle Expenses', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6710', '6710', 'Vehicle Fuel', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6720', '6720', 'Vehicle Maintenance', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6730', '6730', 'Vehicle Lease Payments', v_expense_type_id, 2, true, false, 0),

    (v_org_id, '6800', '6800', 'Depreciation', v_expense_type_id, 1, true, true, 0),
    (v_org_id, '6810', '6810', 'Depreciation Expense', v_expense_type_id, 2, true, true, 0),

    (v_org_id, '6900', '6900', 'Other Expenses', v_expense_type_id, 1, true, false, 0),
    (v_org_id, '6910', '6910', 'Bank Fees & Charges', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6920', '6920', 'Interest Expense', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6930', '6930', 'License & Permits', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6940', '6940', 'Bad Debt Expense', v_expense_type_id, 2, true, false, 0),
    (v_org_id, '6950', '6950', 'Miscellaneous Expense', v_expense_type_id, 2, true, false, 0);

END $$;

-- Mark system accounts
UPDATE chart_of_accounts SET
    is_bank_account = true,
    is_reconcilable = true
WHERE account_code IN ('1010', '1020', '1030', '1040');

\echo 'Chart of Accounts seed data loaded successfully!'
\echo 'Created: 80+ accounts covering all standard accounting categories'

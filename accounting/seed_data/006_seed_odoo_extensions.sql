-- =====================================================
-- Accounting Seed Data: Odoo-Style Extensions
-- Description: Seed data for journals, taxes, multi-currency, payment terms,
--              analytics, deferrals, bank statements, budgets, localization
-- =====================================================

\echo 'Loading Odoo-style extension seed data...';

SET search_path TO accounting, public;

-- =====================================================
-- SECTION 1: JOURNALS
-- =====================================================

\echo 'Seeding journals...';

DO $$
DECLARE
    v_org_id UUID;
    v_bank_account_id UUID;
    v_cash_account_id UUID;
    v_ar_account_id UUID;
    v_ap_account_id UUID;
    v_revenue_account_id UUID;
    v_expense_account_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Get account IDs
    SELECT id INTO v_bank_account_id FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1020' LIMIT 1;
    SELECT id INTO v_cash_account_id FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1010' LIMIT 1;
    SELECT id INTO v_ar_account_id FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1210' LIMIT 1;
    SELECT id INTO v_ap_account_id FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '2010' LIMIT 1;
    SELECT id INTO v_revenue_account_id FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '4010' LIMIT 1;
    SELECT id INTO v_expense_account_id FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '6010' LIMIT 1;

    -- Get bank account reference
    DECLARE v_bank_acct_id UUID;
    BEGIN
        SELECT id INTO v_bank_acct_id FROM bank_accounts WHERE chart_account_id = v_bank_account_id LIMIT 1;

        -- Create journals
        INSERT INTO journals (organization_id, journal_code, journal_name, journal_type, default_debit_account_id, default_credit_account_id, sequence_prefix, is_active, created_by)
        SELECT v_org_id, 'SAJ', 'Sales Journal', 'sale', v_ar_account_id, v_revenue_account_id, 'SAJ', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT DO NOTHING;

        INSERT INTO journals (organization_id, journal_code, journal_name, journal_type, default_debit_account_id, default_credit_account_id, sequence_prefix, is_active, created_by)
        SELECT v_org_id, 'PUR', 'Purchase Journal', 'purchase', v_expense_account_id, v_ap_account_id, 'PUR', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT DO NOTHING;

        INSERT INTO journals (organization_id, journal_code, journal_name, journal_type, bank_account_id, default_debit_account_id, sequence_prefix, is_active, created_by)
        SELECT v_org_id, 'BNK', 'Bank Journal', 'bank', v_bank_acct_id, v_bank_account_id, 'BNK', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT DO NOTHING;

        INSERT INTO journals (organization_id, journal_code, journal_name, journal_type, default_debit_account_id, sequence_prefix, is_active, created_by)
        SELECT v_org_id, 'CSH', 'Cash Journal', 'cash', v_cash_account_id, 'CSH', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT DO NOTHING;

        INSERT INTO journals (organization_id, journal_code, journal_name, journal_type, sequence_prefix, is_active, created_by)
        SELECT v_org_id, 'GEN', 'General Journal', 'general', 'GEN', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT DO NOTHING;

        INSERT INTO journals (organization_id, journal_code, journal_name, journal_type, sequence_prefix, is_active, created_by)
        SELECT v_org_id, 'MISC', 'Miscellaneous Journal', 'miscellaneous', 'MISC', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT DO NOTHING;

        RAISE NOTICE 'Created 6 journals';
    END;
END $$;

-- =====================================================
-- SECTION 2: TAX ENGINE
-- =====================================================

\echo 'Seeding tax engine...';

DO $$
DECLARE
    v_org_id UUID;
    v_tax_group_vat UUID;
    v_tax_group_sales UUID;
    v_tax_payable_account UUID;
    v_tax_expense_account UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Get tax accounts
    SELECT id INTO v_tax_payable_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '2210' LIMIT 1;
    SELECT id INTO v_tax_expense_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '6710' LIMIT 1;

    -- Tax Groups
    INSERT INTO tax_groups (id, organization_id, group_code, group_name, sequence, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'VAT', 'Value Added Tax', 10, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, group_code) DO UPDATE SET group_name = EXCLUDED.group_name
    RETURNING id INTO v_tax_group_vat;

    INSERT INTO tax_groups (id, organization_id, group_code, group_name, sequence, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'SALES', 'Sales Tax', 20, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, group_code) DO UPDATE SET group_name = EXCLUDED.group_name
    RETURNING id INTO v_tax_group_sales;

    INSERT INTO tax_groups (organization_id, group_code, group_name, sequence, is_active, created_by)
    SELECT v_org_id, 'WTAX', 'Withholding Tax', 30, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, group_code) DO NOTHING;

    -- Taxes
    INSERT INTO taxes (organization_id, tax_group_id, tax_code, tax_name, tax_rate, tax_scope, is_price_inclusive, tax_account_id, is_active, created_by)
    SELECT v_org_id, v_tax_group_sales, 'SALES10', 'Sales Tax 10%', 10.0000, 'sales', false, v_tax_payable_account, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, tax_code) DO NOTHING;

    INSERT INTO taxes (organization_id, tax_group_id, tax_code, tax_name, tax_rate, tax_scope, is_price_inclusive, tax_account_id, is_active, created_by)
    SELECT v_org_id, v_tax_group_sales, 'SALES8', 'Sales Tax 8%', 8.0000, 'sales', false, v_tax_payable_account, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, tax_code) DO NOTHING;

    INSERT INTO taxes (organization_id, tax_group_id, tax_code, tax_name, tax_rate, tax_scope, is_price_inclusive, tax_account_id, is_active, created_by)
    SELECT v_org_id, v_tax_group_vat, 'VAT20', 'VAT 20% (Standard)', 20.0000, 'both', false, v_tax_payable_account, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, tax_code) DO NOTHING;

    INSERT INTO taxes (organization_id, tax_group_id, tax_code, tax_name, tax_rate, tax_scope, is_price_inclusive, tax_account_id, is_active, created_by)
    SELECT v_org_id, v_tax_group_vat, 'VAT5', 'VAT 5% (Reduced)', 5.0000, 'both', false, v_tax_payable_account, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, tax_code) DO NOTHING;

    INSERT INTO taxes (organization_id, tax_group_id, tax_code, tax_name, tax_rate, tax_scope, is_price_inclusive, tax_account_id, is_active, created_by)
    SELECT v_org_id, v_tax_group_vat, 'VAT0', 'VAT 0% (Zero-rated)', 0.0000, 'both', false, v_tax_payable_account, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, tax_code) DO NOTHING;

    -- Fiscal Positions
    INSERT INTO fiscal_positions (organization_id, position_code, position_name, auto_apply, country_id, is_active, notes, created_by)
    SELECT v_org_id, 'DOMESTIC', 'Domestic Sales', true, 'US', true, 'Standard domestic sales with tax', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, position_code) DO NOTHING;

    INSERT INTO fiscal_positions (organization_id, position_code, position_name, auto_apply, country_id, is_active, notes, created_by)
    SELECT v_org_id, 'EXPORT', 'Export (Zero-rated)', false, NULL, true, 'International sales - zero VAT', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, position_code) DO NOTHING;

    INSERT INTO fiscal_positions (organization_id, position_code, position_name, auto_apply, country_id, is_active, notes, created_by)
    SELECT v_org_id, 'EU_B2B', 'EU B2B (Reverse Charge)', false, NULL, true, 'EU business-to-business reverse charge', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, position_code) DO NOTHING;

    -- Fiscal Position Tax Mappings (Export: VAT 20% → VAT 0%)
    INSERT INTO fiscal_position_tax_mappings (fiscal_position_id, source_tax_id, destination_tax_id, created_by)
    SELECT fp.id, t1.id, t2.id, u.id
    FROM fiscal_positions fp, taxes t1, taxes t2, users u
    WHERE fp.position_code = 'EXPORT'
        AND t1.tax_code = 'VAT20'
        AND t2.tax_code = 'VAT0'
        AND u.email = 'admin@demoretail.com'
    LIMIT 1
    ON CONFLICT (fiscal_position_id, source_tax_id) DO NOTHING;

    RAISE NOTICE 'Created tax groups, taxes, and fiscal positions';
END $$;

-- =====================================================
-- SECTION 3: CURRENCY RATES (2024 Historical)
-- =====================================================

\echo 'Seeding currency rates...';

DO $$
DECLARE
    v_org_id UUID;
    v_rate_date DATE;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Generate monthly currency rates for 2024
    FOR month_num IN 1..12 LOOP
        v_rate_date := ('2024-' || LPAD(month_num::TEXT, 2, '0') || '-01')::DATE;

        -- EUR to USD (fluctuates around 1.08-1.12)
        INSERT INTO currency_rates (organization_id, currency_code, rate_date, rate, source, created_by)
        SELECT v_org_id, 'EUR', v_rate_date, 1.0900 + (RANDOM() * 0.03), 'manual', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT (organization_id, currency_code, rate_date) DO NOTHING;

        -- GBP to USD (fluctuates around 1.25-1.30)
        INSERT INTO currency_rates (organization_id, currency_code, rate_date, rate, source, created_by)
        SELECT v_org_id, 'GBP', v_rate_date, 1.2700 + (RANDOM() * 0.04), 'manual', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT (organization_id, currency_code, rate_date) DO NOTHING;

        -- JPY to USD (1 USD = ~140-150 JPY, so rate = 0.0067-0.0071)
        INSERT INTO currency_rates (organization_id, currency_code, rate_date, rate, source, created_by)
        SELECT v_org_id, 'JPY', v_rate_date, 0.0069 + (RANDOM() * 0.0004), 'manual', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT (organization_id, currency_code, rate_date) DO NOTHING;

        -- CAD to USD (fluctuates around 0.73-0.76)
        INSERT INTO currency_rates (organization_id, currency_code, rate_date, rate, source, created_by)
        SELECT v_org_id, 'CAD', v_rate_date, 0.7400 + (RANDOM() * 0.03), 'manual', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT (organization_id, currency_code, rate_date) DO NOTHING;

        -- AUD to USD (fluctuates around 0.65-0.68)
        INSERT INTO currency_rates (organization_id, currency_code, rate_date, rate, source, created_by)
        SELECT v_org_id, 'AUD', v_rate_date, 0.6600 + (RANDOM() * 0.03), 'manual', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
        ON CONFLICT (organization_id, currency_code, rate_date) DO NOTHING;
    END LOOP;

    RAISE NOTICE 'Created 60 currency rate records (5 currencies × 12 months)';
END $$;

-- =====================================================
-- SECTION 4: PAYMENT TERMS
-- =====================================================

\echo 'Seeding payment terms...';

DO $$
DECLARE
    v_org_id UUID;
    v_term_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Payment Term: Immediate
    INSERT INTO payment_terms (id, organization_id, term_code, term_name, note, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'IMMEDIATE', 'Immediate Payment', 'Payment due immediately', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, term_code) DO UPDATE SET term_name = EXCLUDED.term_name
    RETURNING id INTO v_term_id;

    INSERT INTO payment_term_lines (payment_term_id, sequence, value_type, value_amount, days_after) VALUES
        (v_term_id, 1, 'balance', 100, 0);

    -- Payment Term: Net 15
    INSERT INTO payment_terms (id, organization_id, term_code, term_name, note, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'NET15', 'Net 15', 'Payment due in 15 days', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, term_code) DO UPDATE SET term_name = EXCLUDED.term_name
    RETURNING id INTO v_term_id;

    INSERT INTO payment_term_lines (payment_term_id, sequence, value_type, value_amount, days_after) VALUES
        (v_term_id, 1, 'balance', 100, 15);

    -- Payment Term: Net 30
    INSERT INTO payment_terms (id, organization_id, term_code, term_name, note, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'NET30', 'Net 30', 'Payment due in 30 days', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, term_code) DO UPDATE SET term_name = EXCLUDED.term_name
    RETURNING id INTO v_term_id;

    INSERT INTO payment_term_lines (payment_term_id, sequence, value_type, value_amount, days_after) VALUES
        (v_term_id, 1, 'balance', 100, 30);

    -- Payment Term: Net 60
    INSERT INTO payment_terms (id, organization_id, term_code, term_name, note, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'NET60', 'Net 60', 'Payment due in 60 days', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, term_code) DO UPDATE SET term_name = EXCLUDED.term_name
    RETURNING id INTO v_term_id;

    INSERT INTO payment_term_lines (payment_term_id, sequence, value_type, value_amount, days_after) VALUES
        (v_term_id, 1, 'balance', 100, 60);

    -- Payment Term: 2/10 Net 30 (2% discount if paid in 10 days, otherwise net 30)
    INSERT INTO payment_terms (id, organization_id, term_code, term_name, note, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, '2/10NET30', '2/10 Net 30', '2% discount if paid within 10 days, otherwise net 30', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, term_code) DO UPDATE SET term_name = EXCLUDED.term_name
    RETURNING id INTO v_term_id;

    INSERT INTO payment_term_lines (payment_term_id, sequence, value_type, value_amount, days_after) VALUES
        (v_term_id, 1, 'balance', 100, 30);

    -- Payment Term: 50% Now, 50% in 30 days
    INSERT INTO payment_terms (id, organization_id, term_code, term_name, note, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, '50/50', '50% Now, 50% in 30 Days', 'Split payment: 50% upfront, 50% in 30 days', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, term_code) DO UPDATE SET term_name = EXCLUDED.term_name
    RETURNING id INTO v_term_id;

    INSERT INTO payment_term_lines (payment_term_id, sequence, value_type, value_amount, days_after) VALUES
        (v_term_id, 1, 'percentage', 50, 0),
        (v_term_id, 2, 'percentage', 50, 30);

    -- Payment Term: End of Month
    INSERT INTO payment_terms (id, organization_id, term_code, term_name, note, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'EOM', 'End of Month', 'Payment due at end of current month', true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, term_code) DO UPDATE SET term_name = EXCLUDED.term_name
    RETURNING id INTO v_term_id;

    INSERT INTO payment_term_lines (payment_term_id, sequence, value_type, value_amount, days_after, end_of_month) VALUES
        (v_term_id, 1, 'balance', 100, 0, true);

    RAISE NOTICE 'Created 7 payment terms with lines';
END $$;

-- =====================================================
-- SECTION 5: ANALYTIC ACCOUNTING
-- =====================================================

\echo 'Seeding analytic accounting...';

DO $$
DECLARE
    v_org_id UUID;
    v_plan_projects UUID;
    v_plan_departments UUID;
    v_plan_regions UUID;
    v_parent_proj_id UUID;
    v_parent_dept_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Analytic Plans
    INSERT INTO analytic_plans (id, organization_id, plan_code, plan_name, is_active, description, created_by)
    SELECT gen_random_uuid(), v_org_id, 'PROJ', 'Projects', true, 'Project-based tracking', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, plan_code) DO UPDATE SET plan_name = EXCLUDED.plan_name
    RETURNING id INTO v_plan_projects;

    INSERT INTO analytic_plans (id, organization_id, plan_code, plan_name, is_active, description, created_by)
    SELECT gen_random_uuid(), v_org_id, 'DEPT', 'Departments', true, 'Department cost centers', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, plan_code) DO UPDATE SET plan_name = EXCLUDED.plan_name
    RETURNING id INTO v_plan_departments;

    INSERT INTO analytic_plans (id, organization_id, plan_code, plan_name, is_active, description, created_by)
    SELECT gen_random_uuid(), v_org_id, 'REGION', 'Regions', true, 'Geographic regions', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, plan_code) DO UPDATE SET plan_name = EXCLUDED.plan_name
    RETURNING id INTO v_plan_regions;

    -- Analytic Accounts: Projects
    INSERT INTO analytic_accounts (id, organization_id, analytic_plan_id, account_code, account_name, account_level, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, v_plan_projects, 'PROJ-ALL', 'All Projects', 1, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO UPDATE SET account_name = EXCLUDED.account_name
    RETURNING id INTO v_parent_proj_id;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, parent_account_id, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_projects, 'PROJ-001', 'Website Redesign', v_parent_proj_id, 2, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, parent_account_id, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_projects, 'PROJ-002', 'Mobile App Development', v_parent_proj_id, 2, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, parent_account_id, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_projects, 'PROJ-003', 'Store Expansion', v_parent_proj_id, 2, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    -- Analytic Accounts: Departments
    INSERT INTO analytic_accounts (id, organization_id, analytic_plan_id, account_code, account_name, account_level, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, v_plan_departments, 'DEPT-ALL', 'All Departments', 1, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO UPDATE SET account_name = EXCLUDED.account_name
    RETURNING id INTO v_parent_dept_id;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, parent_account_id, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_departments, 'DEPT-SALES', 'Sales Department', v_parent_dept_id, 2, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, parent_account_id, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_departments, 'DEPT-MKTG', 'Marketing Department', v_parent_dept_id, 2, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, parent_account_id, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_departments, 'DEPT-OPS', 'Operations', v_parent_dept_id, 2, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, parent_account_id, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_departments, 'DEPT-IT', 'IT Department', v_parent_dept_id, 2, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    -- Analytic Accounts: Regions
    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_regions, 'REGION-NORTH', 'North Region', 1, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_regions, 'REGION-SOUTH', 'South Region', 1, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_regions, 'REGION-EAST', 'East Region', 1, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    INSERT INTO analytic_accounts (organization_id, analytic_plan_id, account_code, account_name, account_level, is_active, created_by)
    SELECT v_org_id, v_plan_regions, 'REGION-WEST', 'West Region', 1, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, account_code) DO NOTHING;

    RAISE NOTICE 'Created 3 analytic plans and 14 analytic accounts';
END $$;

-- =====================================================
-- SECTION 6: DEFERRED REVENUE/EXPENSE CONTRACTS
-- =====================================================

\echo 'Seeding deferral contracts...';

DO $$
DECLARE
    v_org_id UUID;
    v_invoice_id UUID;
    v_bill_id UUID;
    v_deferred_rev_account UUID;
    v_revenue_account UUID;
    v_deferred_exp_account UUID;
    v_expense_account UUID;
    v_contract_id UUID;
    v_monthly_amount NUMERIC(20,4);
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Get accounts
    SELECT id INTO v_deferred_rev_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '2710' LIMIT 1;
    SELECT id INTO v_revenue_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '4010' LIMIT 1;
    SELECT id INTO v_deferred_exp_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1410' LIMIT 1;
    SELECT id INTO v_expense_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '6010' LIMIT 1;

    -- Get sample invoice and bill
    SELECT id INTO v_invoice_id FROM customer_invoices WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_bill_id FROM vendor_bills WHERE organization_id = v_org_id LIMIT 1;

    -- Deferred Revenue Contract: Annual subscription
    INSERT INTO deferred_revenue_contracts (id, organization_id, customer_invoice_id, contract_name, total_deferred_amount, start_date, end_date, recognition_method, deferred_account_id, revenue_account_id, status, created_by)
    SELECT gen_random_uuid(), v_org_id, v_invoice_id, 'Annual Support Contract - Customer A', 24000.00, '2024-01-01', '2024-12-31', 'straight_line', v_deferred_rev_account, v_revenue_account, 'active', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_contract_id;

    -- Generate monthly schedule for deferred revenue (12 months)
    v_monthly_amount := 24000.00 / 12;
    FOR month_num IN 1..12 LOOP
        INSERT INTO deferred_revenue_schedule (contract_id, line_number, recognition_date, recognition_amount, status)
        VALUES (
            v_contract_id,
            month_num,
            ('2024-' || LPAD(month_num::TEXT, 2, '0') || '-28')::DATE,
            v_monthly_amount,
            CASE WHEN month_num <= 10 THEN 'posted' ELSE 'pending' END
        );
    END LOOP;

    -- Deferred Revenue Contract: Warranty sales
    INSERT INTO deferred_revenue_contracts (id, organization_id, contract_name, total_deferred_amount, start_date, end_date, recognition_method, deferred_account_id, revenue_account_id, status, created_by)
    SELECT gen_random_uuid(), v_org_id, 'Extended Warranty Revenue', 6000.00, '2024-03-01', '2025-02-28', 'straight_line', v_deferred_rev_account, v_revenue_account, 'active', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_contract_id;

    -- Generate monthly schedule (12 months from March 2024)
    v_monthly_amount := 6000.00 / 12;
    FOR month_num IN 1..10 LOOP
        INSERT INTO deferred_revenue_schedule (contract_id, line_number, recognition_date, recognition_amount, status)
        VALUES (
            v_contract_id,
            month_num,
            ('2024-' || LPAD((month_num + 2)::TEXT, 2, '0') || '-28')::DATE,
            v_monthly_amount,
            'posted'
        );
    END LOOP;

    -- Deferred Expense Contract: Annual insurance premium
    INSERT INTO deferred_expense_contracts (id, organization_id, vendor_bill_id, contract_name, total_deferred_amount, start_date, end_date, recognition_method, deferred_account_id, expense_account_id, status, created_by)
    SELECT gen_random_uuid(), v_org_id, v_bill_id, 'Annual Insurance Premium', 18000.00, '2024-01-01', '2024-12-31', 'straight_line', v_deferred_exp_account, v_expense_account, 'active', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_contract_id;

    -- Generate monthly schedule for deferred expense (12 months)
    v_monthly_amount := 18000.00 / 12;
    FOR month_num IN 1..12 LOOP
        INSERT INTO deferred_expense_schedule (contract_id, line_number, recognition_date, recognition_amount, status)
        VALUES (
            v_contract_id,
            month_num,
            ('2024-' || LPAD(month_num::TEXT, 2, '0') || '-28')::DATE,
            v_monthly_amount,
            CASE WHEN month_num <= 10 THEN 'posted' ELSE 'pending' END
        );
    END LOOP;

    RAISE NOTICE 'Created 3 deferral contracts with 34 schedule lines';
END $$;

-- =====================================================
-- SECTION 7: BANK STATEMENTS
-- =====================================================

\echo 'Seeding bank statements...';

DO $$
DECLARE
    v_org_id UUID;
    v_bank_account_id UUID;
    v_statement_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_bank_account_id FROM bank_accounts WHERE organization_id = v_org_id LIMIT 1;

    -- Bank Statement: January 2024
    INSERT INTO bank_statements (id, organization_id, bank_account_id, statement_number, statement_date, period_start_date, period_end_date, opening_balance, closing_balance, import_source, status, created_by)
    SELECT gen_random_uuid(), v_org_id, v_bank_account_id, 'STMT-2024-01', '2024-01-31', '2024-01-01', '2024-01-31', 150000.00, 168450.00, 'file_import', 'reconciled', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_statement_id;

    -- Bank Statement Lines for January
    INSERT INTO bank_statement_lines (bank_statement_id, line_number, transaction_date, amount, description, reference, counterparty_name, status)
    VALUES
        (v_statement_id, 1, '2024-01-05', 45000.00, 'Customer payment - INV-001', 'TRX-001', 'Acme Corp', 'matched'),
        (v_statement_id, 2, '2024-01-08', -25000.00, 'Vendor payment - VB-001', 'CHK-1001', 'Tech Supplies Ltd', 'matched'),
        (v_statement_id, 3, '2024-01-15', 32500.00, 'Customer payment - INV-002', 'TRX-002', 'Global Trading', 'matched'),
        (v_statement_id, 4, '2024-01-18', -8500.00, 'Rent payment', 'ACH-001', 'Property Management Co', 'matched'),
        (v_statement_id, 5, '2024-01-22', -15000.00, 'Payroll', 'PAY-001', 'Payroll Service', 'matched'),
        (v_statement_id, 6, '2024-01-25', 28000.00, 'Sales deposit', 'DEP-001', 'POS System', 'matched'),
        (v_statement_id, 7, '2024-01-28', -18550.00, 'Vendor payments', 'CHK-1002', 'Various', 'matched');

    -- Bank Statement: February 2024
    INSERT INTO bank_statements (id, organization_id, bank_account_id, statement_number, statement_date, period_start_date, period_end_date, opening_balance, closing_balance, import_source, status, created_by)
    SELECT gen_random_uuid(), v_org_id, v_bank_account_id, 'STMT-2024-02', '2024-02-29', '2024-02-01', '2024-02-29', 168450.00, 185720.00, 'file_import', 'reconciled', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_statement_id;

    -- Bank Statement Lines for February
    INSERT INTO bank_statement_lines (bank_statement_id, line_number, transaction_date, amount, description, reference, counterparty_name, status)
    VALUES
        (v_statement_id, 1, '2024-02-03', 38000.00, 'Customer payments', 'TRX-003', 'Various customers', 'matched'),
        (v_statement_id, 2, '2024-02-07', -22000.00, 'Vendor payments', 'CHK-1003', 'Inventory suppliers', 'matched'),
        (v_statement_id, 3, '2024-02-12', 42500.00, 'Sales deposit', 'DEP-002', 'POS System', 'matched'),
        (v_statement_id, 4, '2024-02-15', -15000.00, 'Payroll', 'PAY-002', 'Payroll Service', 'matched'),
        (v_statement_id, 5, '2024-02-20', -8500.00, 'Rent', 'ACH-002', 'Property Management', 'matched'),
        (v_statement_id, 6, '2024-02-25', 30000.00, 'Customer payment', 'TRX-004', 'Enterprise Client', 'matched'),
        (v_statement_id, 7, '2024-02-28', -19730.00, 'Operating expenses', 'Various', 'Multiple vendors', 'partial_match'),
        (v_statement_id, 8, '2024-02-28', 500.00, 'Interest income', 'INT-001', 'Bank', 'unmatched');

    RAISE NOTICE 'Created 2 bank statements with 15 statement lines';
END $$;

-- =====================================================
-- SECTION 8: BUDGETS
-- =====================================================

\echo 'Seeding budgets...';

DO $$
DECLARE
    v_org_id UUID;
    v_fiscal_year_id UUID;
    v_budget_id UUID;
    v_revenue_account UUID;
    v_cogs_account UUID;
    v_expense_account UUID;
    v_dept_sales UUID;
    v_dept_mktg UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_fiscal_year_id FROM fiscal_years WHERE organization_id = v_org_id AND fiscal_year = '2024' LIMIT 1;

    -- Get accounts
    SELECT id INTO v_revenue_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '4010' LIMIT 1;
    SELECT id INTO v_cogs_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '5010' LIMIT 1;
    SELECT id INTO v_expense_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '6010' LIMIT 1;

    -- Get analytic accounts
    SELECT id INTO v_dept_sales FROM analytic_accounts WHERE organization_id = v_org_id AND account_code = 'DEPT-SALES' LIMIT 1;
    SELECT id INTO v_dept_mktg FROM analytic_accounts WHERE organization_id = v_org_id AND account_code = 'DEPT-MKTG' LIMIT 1;

    -- Budget: FY 2024 Operating Budget
    INSERT INTO budgets (id, organization_id, budget_code, budget_name, fiscal_year_id, start_date, end_date, budget_type, status, notes, created_by)
    SELECT gen_random_uuid(), v_org_id, 'FY2024-OP', 'FY 2024 Operating Budget', v_fiscal_year_id, '2024-01-01', '2024-12-31', 'operating', 'approved', 'Annual operating budget for 2024', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, budget_code) DO UPDATE SET budget_name = EXCLUDED.budget_name
    RETURNING id INTO v_budget_id;

    -- Budget Lines: Revenue (by month)
    FOR month_num IN 1..12 LOOP
        INSERT INTO budget_lines (budget_id, account_id, accounting_period_id, period_start_date, period_end_date, planned_amount)
        SELECT v_budget_id, v_revenue_account, ap.id, ap.start_date, ap.end_date, 175000.00
        FROM accounting_periods ap
        WHERE ap.fiscal_year_id = v_fiscal_year_id AND ap.period_number = month_num;
    END LOOP;

    -- Budget Lines: COGS (by month, ~60% of revenue)
    FOR month_num IN 1..12 LOOP
        INSERT INTO budget_lines (budget_id, account_id, accounting_period_id, period_start_date, period_end_date, planned_amount)
        SELECT v_budget_id, v_cogs_account, ap.id, ap.start_date, ap.end_date, 105000.00
        FROM accounting_periods ap
        WHERE ap.fiscal_year_id = v_fiscal_year_id AND ap.period_number = month_num;
    END LOOP;

    -- Budget Lines: Operating Expenses (by month)
    FOR month_num IN 1..12 LOOP
        INSERT INTO budget_lines (budget_id, account_id, accounting_period_id, period_start_date, period_end_date, planned_amount)
        SELECT v_budget_id, v_expense_account, ap.id, ap.start_date, ap.end_date, 55000.00
        FROM accounting_periods ap
        WHERE ap.fiscal_year_id = v_fiscal_year_id AND ap.period_number = month_num;
    END LOOP;

    -- Budget: Departmental Budget
    INSERT INTO budgets (id, organization_id, budget_code, budget_name, fiscal_year_id, start_date, end_date, budget_type, status, notes, created_by)
    SELECT gen_random_uuid(), v_org_id, 'FY2024-DEPT', 'FY 2024 Departmental Budget', v_fiscal_year_id, '2024-01-01', '2024-12-31', 'departmental', 'approved', 'Departmental cost allocation', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    ON CONFLICT (organization_id, budget_code) DO UPDATE SET budget_name = EXCLUDED.budget_name
    RETURNING id INTO v_budget_id;

    -- Department budgets (annual)
    INSERT INTO budget_lines (budget_id, analytic_account_id, period_start_date, period_end_date, planned_amount, notes)
    VALUES
        (v_budget_id, v_dept_sales, '2024-01-01', '2024-12-31', 350000.00, 'Sales department annual budget'),
        (v_budget_id, v_dept_mktg, '2024-01-01', '2024-12-31', 120000.00, 'Marketing department annual budget');

    RAISE NOTICE 'Created 2 budgets with 38 budget lines';
END $$;

-- =====================================================
-- SECTION 9: TAX REPORT DEFINITIONS
-- =====================================================

\echo 'Seeding tax report definitions...';

DO $$
DECLARE
    v_org_id UUID;
    v_us_package UUID;
    v_report_id UUID;
    v_sales_tax_group UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_us_package FROM localization_packages WHERE package_code = 'us_gaap' LIMIT 1;
    SELECT id INTO v_sales_tax_group FROM tax_groups WHERE organization_id = v_org_id AND group_code = 'SALES' LIMIT 1;

    -- Tax Report: US Sales Tax Return
    INSERT INTO tax_report_definitions (id, organization_id, localization_package_id, report_code, report_name, jurisdiction, authority, report_frequency, version, is_active, description, created_by)
    SELECT gen_random_uuid(), v_org_id, v_us_package, 'US-SALES-TAX', 'US Sales Tax Return', 'United States', 'State Tax Authority', 'monthly', '1.0', true, 'Monthly sales tax return for US jurisdictions', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_report_id;

    -- Tax Report Lines
    INSERT INTO tax_report_lines (tax_report_definition_id, line_code, line_name, sequence, formula_type, is_subtotal)
    VALUES
        (v_report_id, 'L1', 'Gross Sales', 10, 'sum', false),
        (v_report_id, 'L2', 'Exempt Sales', 20, 'sum', false),
        (v_report_id, 'L3', 'Taxable Sales', 30, 'formula', true),
        (v_report_id, 'L4', 'Tax Collected', 40, 'detail', false),
        (v_report_id, 'L5', 'Tax Due', 50, 'formula', true);

    RAISE NOTICE 'Created 1 tax report definition with 5 report lines';
END $$;

\echo '';
\echo '==========================================';
\echo 'Odoo Extensions Seed Data Summary';
\echo '==========================================';
\echo 'Journals:                    6 journals';
\echo 'Tax Groups:                  3 groups';
\echo 'Taxes:                       6 taxes';
\echo 'Fiscal Positions:            3 positions with mappings';
\echo 'Currency Rates:              60 rates (5 currencies × 12 months)';
\echo 'Payment Terms:               7 terms with lines';
\echo 'Analytic Plans:              3 plans';
\echo 'Analytic Accounts:           14 accounts';
\echo 'Deferred Contracts:          3 contracts with 34 schedules';
\echo 'Bank Statements:             2 statements with 15 lines';
\echo 'Budgets:                     2 budgets with 38 lines';
\echo 'Tax Report Definitions:      1 report with 5 lines';
\echo '==========================================';
\echo '';
\echo 'All Odoo-style extension data loaded successfully!';

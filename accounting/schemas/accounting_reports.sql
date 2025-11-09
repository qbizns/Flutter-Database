-- =====================================================
-- Accounting Report Views
-- Description: Financial statement and report views
-- Includes: Trial Balance, Balance Sheet, Income Statement,
--           Cash Flow, Account Activity, Aged AP/AR
-- =====================================================

\echo 'Creating accounting report views...';

-- =====================================================
-- VIEW 1: Trial Balance
-- Shows all accounts with their debit/credit balances
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_trial_balance AS
WITH account_balances AS (
    SELECT
        gl.organization_id,
        gl.account_id,
        coa.account_code,
        coa.account_number,
        coa.account_name,
        at.type_code,
        at.type_name,
        at.normal_balance,
        SUM(gl.debit_amount) as total_debits,
        SUM(gl.credit_amount) as total_credits,
        CASE
            WHEN at.normal_balance = 'debit' THEN
                SUM(gl.debit_amount) - SUM(gl.credit_amount)
            ELSE
                SUM(gl.credit_amount) - SUM(gl.debit_amount)
        END as balance
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    WHERE gl.deleted_at IS NULL
        AND coa.deleted_at IS NULL
    GROUP BY gl.organization_id, gl.account_id, coa.account_code,
             coa.account_number, coa.account_name, at.type_code,
             at.type_name, at.normal_balance
)
SELECT
    organization_id,
    account_code,
    account_number,
    account_name,
    type_code,
    type_name,
    normal_balance,
    total_debits,
    total_credits,
    CASE
        WHEN normal_balance = 'debit' THEN balance
        ELSE 0
    END as debit_balance,
    CASE
        WHEN normal_balance = 'credit' THEN balance
        ELSE 0
    END as credit_balance,
    balance as net_balance
FROM account_balances
WHERE balance != 0
ORDER BY account_code;

COMMENT ON VIEW accounting.view_trial_balance IS 'Trial Balance report showing all account balances with debits and credits';

-- =====================================================
-- VIEW 2: Balance Sheet
-- Statement of Financial Position
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_balance_sheet AS
WITH account_balances AS (
    SELECT
        gl.organization_id,
        coa.account_code,
        coa.account_name,
        at.type_code,
        at.type_name,
        CASE
            WHEN at.normal_balance = 'debit' THEN
                SUM(gl.debit_amount) - SUM(gl.credit_amount)
            ELSE
                SUM(gl.credit_amount) - SUM(gl.debit_amount)
        END as balance
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    WHERE gl.deleted_at IS NULL
        AND coa.deleted_at IS NULL
        AND at.type_code IN ('ASSET', 'LIABILITY', 'EQUITY')
    GROUP BY gl.organization_id, coa.account_code, coa.account_name,
             at.type_code, at.type_name, at.normal_balance
)
SELECT
    organization_id,
    type_code as section,
    type_name as section_name,
    account_code,
    account_name,
    balance,
    -- Calculate section totals
    SUM(balance) OVER (
        PARTITION BY organization_id, type_code
    ) as section_total
FROM account_balances
WHERE balance != 0
ORDER BY
    CASE type_code
        WHEN 'ASSET' THEN 1
        WHEN 'LIABILITY' THEN 2
        WHEN 'EQUITY' THEN 3
    END,
    account_code;

COMMENT ON VIEW accounting.view_balance_sheet IS 'Balance Sheet (Statement of Financial Position) showing Assets, Liabilities, and Equity';

-- =====================================================
-- VIEW 3: Income Statement
-- Profit & Loss Statement
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_income_statement AS
WITH account_balances AS (
    SELECT
        gl.organization_id,
        fy.fiscal_year,
        ap.period_name,
        coa.account_code,
        coa.account_name,
        at.type_code,
        at.type_name,
        CASE
            WHEN at.type_code IN ('REVENUE') THEN
                SUM(gl.credit_amount) - SUM(gl.debit_amount)
            ELSE
                SUM(gl.debit_amount) - SUM(gl.credit_amount)
        END as amount
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    INNER JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
    LEFT JOIN accounting_periods ap ON gl.accounting_period_id = ap.id
    WHERE gl.deleted_at IS NULL
        AND coa.deleted_at IS NULL
        AND at.type_code IN ('REVENUE', 'COGS', 'EXPENSE')
    GROUP BY gl.organization_id, fy.fiscal_year, ap.period_name,
             coa.account_code, coa.account_name, at.type_code, at.type_name
)
SELECT
    organization_id,
    fiscal_year,
    type_code as section,
    type_name as section_name,
    account_code,
    account_name,
    amount,
    -- Calculate section totals
    SUM(amount) OVER (
        PARTITION BY organization_id, fiscal_year, type_code
    ) as section_total,
    -- Calculate gross profit (Revenue - COGS)
    SUM(CASE WHEN type_code = 'REVENUE' THEN amount ELSE 0 END) OVER (
        PARTITION BY organization_id, fiscal_year
    ) - SUM(CASE WHEN type_code = 'COGS' THEN amount ELSE 0 END) OVER (
        PARTITION BY organization_id, fiscal_year
    ) as gross_profit,
    -- Calculate net income (Revenue - COGS - Expenses)
    SUM(CASE WHEN type_code = 'REVENUE' THEN amount ELSE 0 END) OVER (
        PARTITION BY organization_id, fiscal_year
    ) - SUM(CASE WHEN type_code IN ('COGS', 'EXPENSE') THEN amount ELSE 0 END) OVER (
        PARTITION BY organization_id, fiscal_year
    ) as net_income
FROM account_balances
WHERE amount != 0
ORDER BY
    fiscal_year,
    CASE type_code
        WHEN 'REVENUE' THEN 1
        WHEN 'COGS' THEN 2
        WHEN 'EXPENSE' THEN 3
    END,
    account_code;

COMMENT ON VIEW accounting.view_income_statement IS 'Income Statement (Profit & Loss) showing Revenue, COGS, Expenses, and Net Income';

-- =====================================================
-- VIEW 4: Account Activity Detail
-- Detailed transaction history by account
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_account_activity AS
SELECT
    gl.organization_id,
    o.organization_name,
    coa.account_code,
    coa.account_number,
    coa.account_name,
    at.type_name as account_type,
    gl.transaction_date,
    je.entry_number as journal_entry_number,
    je.description as journal_entry_description,
    jet.type_name as entry_type,
    gl.description as line_description,
    gl.debit_amount,
    gl.credit_amount,
    -- Running balance calculation
    SUM(
        CASE WHEN at.normal_balance = 'debit'
            THEN gl.debit_amount - gl.credit_amount
            ELSE gl.credit_amount - gl.debit_amount
        END
    ) OVER (
        PARTITION BY gl.organization_id, gl.account_id
        ORDER BY gl.transaction_date, gl.created_at
    ) as running_balance,
    fy.fiscal_year,
    ap.period_name as accounting_period,
    gl.created_at as posted_date,
    u.full_name as posted_by
FROM general_ledger gl
INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
INNER JOIN account_types at ON coa.account_type_id = at.id
INNER JOIN organizations o ON gl.organization_id = o.id
INNER JOIN journal_entries je ON gl.journal_entry_id = je.id
LEFT JOIN journal_entry_types jet ON je.journal_entry_type_id = jet.id
LEFT JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
LEFT JOIN accounting_periods ap ON gl.accounting_period_id = ap.id
LEFT JOIN users u ON gl.created_by = u.id
WHERE gl.deleted_at IS NULL
ORDER BY coa.account_code, gl.transaction_date, gl.created_at;

COMMENT ON VIEW accounting.view_account_activity IS 'Detailed account activity showing all transactions with running balance';

-- =====================================================
-- VIEW 5: Aged Accounts Payable Report
-- Shows outstanding vendor bills by aging period
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_aged_accounts_payable AS
SELECT
    vb.organization_id,
    o.organization_name,
    s.supplier_name,
    s.contact_email,
    s.phone_number,
    vb.bill_number,
    vb.bill_date,
    vb.due_date,
    vb.total_amount,
    vb.paid_amount,
    vb.balance_due,
    vb.status,
    CURRENT_DATE - vb.due_date as days_overdue,
    CASE
        WHEN vb.status = 'paid' THEN 0
        WHEN CURRENT_DATE <= vb.due_date THEN vb.balance_due
        ELSE 0
    END as current_amount,
    CASE
        WHEN CURRENT_DATE > vb.due_date
            AND CURRENT_DATE <= vb.due_date + INTERVAL '30 days'
            AND vb.status != 'paid'
        THEN vb.balance_due
        ELSE 0
    END as days_1_30,
    CASE
        WHEN CURRENT_DATE > vb.due_date + INTERVAL '30 days'
            AND CURRENT_DATE <= vb.due_date + INTERVAL '60 days'
            AND vb.status != 'paid'
        THEN vb.balance_due
        ELSE 0
    END as days_31_60,
    CASE
        WHEN CURRENT_DATE > vb.due_date + INTERVAL '60 days'
            AND CURRENT_DATE <= vb.due_date + INTERVAL '90 days'
            AND vb.status != 'paid'
        THEN vb.balance_due
        ELSE 0
    END as days_61_90,
    CASE
        WHEN CURRENT_DATE > vb.due_date + INTERVAL '90 days'
            AND vb.status != 'paid'
        THEN vb.balance_due
        ELSE 0
    END as days_over_90
FROM vendor_bills vb
INNER JOIN organizations o ON vb.organization_id = o.id
INNER JOIN suppliers s ON vb.supplier_id = s.id
WHERE vb.deleted_at IS NULL
    AND vb.balance_due > 0
ORDER BY s.supplier_name, vb.due_date;

COMMENT ON VIEW accounting.view_aged_accounts_payable IS 'Aged AP report showing outstanding vendor bills by aging bucket';

-- =====================================================
-- VIEW 6: Aged Accounts Receivable Report
-- Shows outstanding customer invoices by aging period
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_aged_accounts_receivable AS
SELECT
    ci.organization_id,
    o.organization_name,
    c.customer_name,
    c.email,
    c.phone,
    ci.invoice_number,
    ci.invoice_date,
    ci.due_date,
    ci.total_amount,
    ci.paid_amount,
    ci.balance_due,
    ci.status,
    CURRENT_DATE - ci.due_date as days_overdue,
    CASE
        WHEN ci.status = 'paid' THEN 0
        WHEN CURRENT_DATE <= ci.due_date THEN ci.balance_due
        ELSE 0
    END as current_amount,
    CASE
        WHEN CURRENT_DATE > ci.due_date
            AND CURRENT_DATE <= ci.due_date + INTERVAL '30 days'
            AND ci.status != 'paid'
        THEN ci.balance_due
        ELSE 0
    END as days_1_30,
    CASE
        WHEN CURRENT_DATE > ci.due_date + INTERVAL '30 days'
            AND CURRENT_DATE <= ci.due_date + INTERVAL '60 days'
            AND ci.status != 'paid'
        THEN ci.balance_due
        ELSE 0
    END as days_31_60,
    CASE
        WHEN CURRENT_DATE > ci.due_date + INTERVAL '60 days'
            AND CURRENT_DATE <= ci.due_date + INTERVAL '90 days'
            AND ci.status != 'paid'
        THEN ci.balance_due
        ELSE 0
    END as days_61_90,
    CASE
        WHEN CURRENT_DATE > ci.due_date + INTERVAL '90 days'
            AND ci.status != 'paid'
        THEN ci.balance_due
        ELSE 0
    END as days_over_90
FROM customer_invoices ci
INNER JOIN organizations o ON ci.organization_id = o.id
INNER JOIN customers c ON ci.customer_id = c.id
WHERE ci.deleted_at IS NULL
    AND ci.balance_due > 0
ORDER BY c.customer_name, ci.due_date;

COMMENT ON VIEW accounting.view_aged_accounts_receivable IS 'Aged AR report showing outstanding customer invoices by aging bucket';

-- =====================================================
-- VIEW 7: Cash Flow Statement
-- Statement of Cash Flows (Indirect Method)
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_cash_flow AS
WITH cash_accounts AS (
    SELECT id FROM chart_of_accounts
    WHERE is_bank_account = true
),
operating_activities AS (
    SELECT
        gl.organization_id,
        fy.fiscal_year,
        SUM(CASE WHEN gl.debit_amount > 0 THEN gl.debit_amount ELSE -gl.credit_amount END) as net_cash_operations
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    INNER JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
    WHERE at.type_code IN ('REVENUE', 'EXPENSE', 'COGS')
    GROUP BY gl.organization_id, fy.fiscal_year
),
investing_activities AS (
    SELECT
        fa.organization_id,
        EXTRACT(YEAR FROM fa.acquisition_date)::TEXT as fiscal_year,
        -SUM(fa.acquisition_cost) as net_cash_investing
    FROM fixed_assets fa
    WHERE fa.deleted_at IS NULL
    GROUP BY fa.organization_id, EXTRACT(YEAR FROM fa.acquisition_date)
),
financing_activities AS (
    SELECT
        gl.organization_id,
        fy.fiscal_year,
        SUM(CASE WHEN gl.credit_amount > 0 THEN gl.credit_amount ELSE -gl.debit_amount END) as net_cash_financing
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    INNER JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
    WHERE at.type_code IN ('EQUITY', 'LIABILITY')
        AND coa.account_code NOT LIKE '2%' -- Exclude operating liabilities
    GROUP BY gl.organization_id, fy.fiscal_year
),
cash_beginning_balance AS (
    SELECT
        gl.organization_id,
        fy.fiscal_year,
        SUM(gl.debit_amount - gl.credit_amount) as beginning_cash
    FROM general_ledger gl
    INNER JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
    WHERE gl.account_id IN (SELECT id FROM cash_accounts)
        AND gl.transaction_date < fy.start_date
    GROUP BY gl.organization_id, fy.fiscal_year
)
SELECT
    oa.organization_id,
    oa.fiscal_year,
    COALESCE(oa.net_cash_operations, 0) as cash_from_operations,
    COALESCE(ia.net_cash_investing, 0) as cash_from_investing,
    COALESCE(fa.net_cash_financing, 0) as cash_from_financing,
    COALESCE(oa.net_cash_operations, 0) +
    COALESCE(ia.net_cash_investing, 0) +
    COALESCE(fa.net_cash_financing, 0) as net_change_in_cash,
    COALESCE(cbb.beginning_cash, 0) as beginning_cash_balance,
    COALESCE(cbb.beginning_cash, 0) +
    COALESCE(oa.net_cash_operations, 0) +
    COALESCE(ia.net_cash_investing, 0) +
    COALESCE(fa.net_cash_financing, 0) as ending_cash_balance
FROM operating_activities oa
LEFT JOIN investing_activities ia ON oa.organization_id = ia.organization_id AND oa.fiscal_year = ia.fiscal_year
LEFT JOIN financing_activities fa ON oa.organization_id = fa.organization_id AND oa.fiscal_year = fa.fiscal_year
LEFT JOIN cash_beginning_balance cbb ON oa.organization_id = cbb.organization_id AND oa.fiscal_year = cbb.fiscal_year
ORDER BY oa.fiscal_year;

COMMENT ON VIEW accounting.view_cash_flow IS 'Cash Flow Statement showing operating, investing, and financing activities';

-- =====================================================
-- VIEW 8: Financial Ratios
-- Key financial ratios and metrics
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_financial_ratios AS
WITH balance_sheet_data AS (
    SELECT
        gl.organization_id,
        fy.fiscal_year,
        SUM(CASE WHEN at.type_code = 'ASSET' THEN
            gl.debit_amount - gl.credit_amount ELSE 0 END) as total_assets,
        SUM(CASE WHEN at.type_code = 'LIABILITY' THEN
            gl.credit_amount - gl.debit_amount ELSE 0 END) as total_liabilities,
        SUM(CASE WHEN at.type_code = 'EQUITY' THEN
            gl.credit_amount - gl.debit_amount ELSE 0 END) as total_equity,
        SUM(CASE WHEN at.type_code = 'ASSET' AND coa.account_code LIKE '1%' AND coa.account_code < '1500' THEN
            gl.debit_amount - gl.credit_amount ELSE 0 END) as current_assets,
        SUM(CASE WHEN at.type_code = 'LIABILITY' AND coa.account_code LIKE '2%' AND coa.account_code < '2500' THEN
            gl.credit_amount - gl.debit_amount ELSE 0 END) as current_liabilities
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    INNER JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
    GROUP BY gl.organization_id, fy.fiscal_year
),
income_data AS (
    SELECT
        gl.organization_id,
        fy.fiscal_year,
        SUM(CASE WHEN at.type_code = 'REVENUE' THEN
            gl.credit_amount - gl.debit_amount ELSE 0 END) as total_revenue,
        SUM(CASE WHEN at.type_code = 'COGS' THEN
            gl.debit_amount - gl.credit_amount ELSE 0 END) as total_cogs,
        SUM(CASE WHEN at.type_code = 'EXPENSE' THEN
            gl.debit_amount - gl.credit_amount ELSE 0 END) as total_expenses
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    INNER JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
    GROUP BY gl.organization_id, fy.fiscal_year
)
SELECT
    bs.organization_id,
    bs.fiscal_year,
    -- Liquidity Ratios
    CASE WHEN bs.current_liabilities > 0
        THEN ROUND(bs.current_assets / bs.current_liabilities, 2)
        ELSE NULL
    END as current_ratio,
    -- Profitability Ratios
    CASE WHEN inc.total_revenue > 0
        THEN ROUND((inc.total_revenue - inc.total_cogs - inc.total_expenses) / inc.total_revenue * 100, 2)
        ELSE NULL
    END as net_profit_margin_pct,
    CASE WHEN inc.total_revenue > 0
        THEN ROUND((inc.total_revenue - inc.total_cogs) / inc.total_revenue * 100, 2)
        ELSE NULL
    END as gross_profit_margin_pct,
    CASE WHEN bs.total_equity > 0
        THEN ROUND((inc.total_revenue - inc.total_cogs - inc.total_expenses) / bs.total_equity * 100, 2)
        ELSE NULL
    END as return_on_equity_pct,
    CASE WHEN bs.total_assets > 0
        THEN ROUND((inc.total_revenue - inc.total_cogs - inc.total_expenses) / bs.total_assets * 100, 2)
        ELSE NULL
    END as return_on_assets_pct,
    -- Leverage Ratios
    CASE WHEN bs.total_assets > 0
        THEN ROUND(bs.total_liabilities / bs.total_assets * 100, 2)
        ELSE NULL
    END as debt_to_assets_pct,
    CASE WHEN bs.total_equity > 0
        THEN ROUND(bs.total_liabilities / bs.total_equity, 2)
        ELSE NULL
    END as debt_to_equity_ratio,
    -- Raw data for reference
    bs.total_assets,
    bs.total_liabilities,
    bs.total_equity,
    bs.current_assets,
    bs.current_liabilities,
    inc.total_revenue,
    inc.total_cogs,
    inc.total_expenses,
    inc.total_revenue - inc.total_cogs - inc.total_expenses as net_income
FROM balance_sheet_data bs
INNER JOIN income_data inc ON bs.organization_id = inc.organization_id AND bs.fiscal_year = inc.fiscal_year
ORDER BY bs.fiscal_year;

COMMENT ON VIEW accounting.view_financial_ratios IS 'Key financial ratios including liquidity, profitability, and leverage metrics';

\echo 'Accounting report views created successfully!';
\echo '';
\echo 'Available Reports:';
\echo '  ✓ view_trial_balance - Trial Balance';
\echo '  ✓ view_balance_sheet - Balance Sheet';
\echo '  ✓ view_income_statement - Income Statement (P&L)';
\echo '  ✓ view_account_activity - Account Activity Detail';
\echo '  ✓ view_aged_accounts_payable - Aged AP Report';
\echo '  ✓ view_aged_accounts_receivable - Aged AR Report';
\echo '  ✓ view_cash_flow - Cash Flow Statement';
\echo '  ✓ view_financial_ratios - Financial Ratios & KPIs';
\echo '';

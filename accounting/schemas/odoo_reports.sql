-- =====================================================
-- Odoo-Style Accounting Report Views
-- Description: Advanced reporting for journals, taxes, analytics, deferrals, budgets
-- =====================================================

\echo 'Creating Odoo-style report views...';

SET search_path TO accounting, public;

-- =====================================================
-- VIEW 1: Journal Entry Report (by Journal)
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_journal_entries_by_journal AS
SELECT
    je.organization_id,
    o.organization_name,
    j.journal_code,
    j.journal_name,
    j.journal_type,
    je.entry_number,
    je.entry_date,
    je.posting_date,
    je.description,
    je.status,
    je.is_posted,
    je.total_debit,
    je.total_credit,
    fy.fiscal_year,
    ap.period_name,
    je.created_at as posted_date,
    u.full_name as created_by_name
FROM journal_entries je
INNER JOIN organizations o ON je.organization_id = o.id
LEFT JOIN journals j ON je.journal_id = j.id
LEFT JOIN fiscal_years fy ON je.fiscal_year_id = fy.id
LEFT JOIN accounting_periods ap ON je.accounting_period_id = ap.id
LEFT JOIN users u ON je.created_by = u.id
WHERE je.deleted_at IS NULL
ORDER BY j.journal_code, je.entry_date DESC, je.entry_number;

COMMENT ON VIEW accounting.view_journal_entries_by_journal IS 'Journal entries grouped by journal type (Sales, Purchase, Bank, etc.)';

-- =====================================================
-- VIEW 2: Tax Report (Tax Collected and Paid)
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_tax_report AS
WITH tax_transactions AS (
    SELECT
        jel.organization_id,
        t.id as tax_id,
        t.tax_code,
        t.tax_name,
        t.tax_rate,
        t.tax_scope,
        tg.group_name as tax_group,
        je.entry_date as transaction_date,
        EXTRACT(YEAR FROM je.entry_date) as fiscal_year,
        EXTRACT(MONTH FROM je.entry_date) as fiscal_month,
        -- Tax amount (credit = tax collected on sales, debit = tax paid on purchases)
        CASE
            WHEN t.tax_scope IN ('sales', 'both') THEN jel.credit_amount
            ELSE 0
        END as tax_collected,
        CASE
            WHEN t.tax_scope IN ('purchases', 'both') THEN jel.debit_amount
            ELSE 0
        END as tax_paid,
        -- Base amount (reverse calculate from tax)
        CASE
            WHEN t.tax_rate > 0 THEN
                CASE
                    WHEN t.tax_scope IN ('sales', 'both') THEN jel.credit_amount / (t.tax_rate / 100)
                    ELSE jel.debit_amount / (t.tax_rate / 100)
                END
            ELSE 0
        END as base_amount
    FROM journal_entry_lines jel
    INNER JOIN journal_entries je ON jel.journal_entry_id = je.id
    INNER JOIN taxes t ON jel.tax_id = t.id
    LEFT JOIN tax_groups tg ON t.tax_group_id = tg.id
    WHERE jel.deleted_at IS NULL
        AND je.is_posted = true
        AND je.deleted_at IS NULL
)
SELECT
    organization_id,
    tax_code,
    tax_name,
    tax_rate,
    tax_scope,
    tax_group,
    fiscal_year,
    fiscal_month,
    TO_DATE(fiscal_year || '-' || LPAD(fiscal_month::TEXT, 2, '0') || '-01', 'YYYY-MM-DD') as period_start,
    SUM(base_amount) as total_base_amount,
    SUM(tax_collected) as total_tax_collected,
    SUM(tax_paid) as total_tax_paid,
    SUM(tax_collected) - SUM(tax_paid) as net_tax_due,
    COUNT(*) as transaction_count
FROM tax_transactions
GROUP BY organization_id, tax_code, tax_name, tax_rate, tax_scope, tax_group, fiscal_year, fiscal_month
ORDER BY fiscal_year DESC, fiscal_month DESC, tax_code;

COMMENT ON VIEW accounting.view_tax_report IS 'Tax collected vs paid by tax type and period';

-- =====================================================
-- VIEW 3: Multi-Currency Summary
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_multi_currency_summary AS
SELECT
    je.organization_id,
    o.organization_name,
    o.base_currency_code,
    je.currency_code as transaction_currency,
    je.exchange_rate,
    je.entry_date,
    je.entry_number,
    je.description,
    SUM(jel.amount_currency) as total_foreign_amount,
    SUM(jel.debit_amount + jel.credit_amount) as total_base_amount,
    -- FX gain/loss (difference between current rate and original rate)
    SUM(jel.amount_currency * (
        COALESCE((SELECT rate FROM currency_rates cr
                  WHERE cr.organization_id = je.organization_id
                    AND cr.currency_code = je.currency_code
                    AND cr.rate_date <= je.entry_date
                  ORDER BY cr.rate_date DESC
                  LIMIT 1), je.exchange_rate)
        - je.exchange_rate
    )) as unrealized_fx_gain_loss
FROM journal_entries je
INNER JOIN organizations o ON je.organization_id = o.id
INNER JOIN journal_entry_lines jel ON je.id = jel.journal_entry_id
WHERE je.deleted_at IS NULL
    AND je.currency_code IS NOT NULL
    AND je.currency_code != o.base_currency_code
    AND jel.deleted_at IS NULL
GROUP BY je.organization_id, o.organization_name, o.base_currency_code,
         je.currency_code, je.exchange_rate, je.entry_date, je.entry_number, je.description
ORDER BY je.entry_date DESC;

COMMENT ON VIEW accounting.view_multi_currency_summary IS 'Foreign currency transactions with FX gain/loss calculation';

-- =====================================================
-- VIEW 4: Payment Schedule Report
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_payment_schedules AS
SELECT
    ips.organization_id,
    o.organization_name,
    ips.source_type,
    ips.source_id,
    CASE
        WHEN ips.source_type = 'customer_invoice' THEN ci.invoice_number
        WHEN ips.source_type = 'vendor_bill' THEN vb.bill_number
    END as document_number,
    CASE
        WHEN ips.source_type = 'customer_invoice' THEN c.customer_name
        WHEN ips.source_type = 'vendor_bill' THEN s.supplier_name
    END as partner_name,
    ips.line_number,
    ips.due_date,
    ips.amount_due,
    ips.amount_paid,
    ips.amount_due - ips.amount_paid as balance_remaining,
    ips.status,
    CURRENT_DATE - ips.due_date as days_overdue,
    CASE
        WHEN ips.status = 'paid' THEN 'Paid'
        WHEN CURRENT_DATE <= ips.due_date THEN 'Current'
        WHEN CURRENT_DATE <= ips.due_date + INTERVAL '30 days' THEN '1-30 days overdue'
        WHEN CURRENT_DATE <= ips.due_date + INTERVAL '60 days' THEN '31-60 days overdue'
        WHEN CURRENT_DATE <= ips.due_date + INTERVAL '90 days' THEN '61-90 days overdue'
        ELSE 'Over 90 days overdue'
    END as aging_bucket
FROM invoice_payment_schedules ips
INNER JOIN organizations o ON ips.organization_id = o.id
LEFT JOIN customer_invoices ci ON ips.source_type = 'customer_invoice' AND ips.source_id = ci.id
LEFT JOIN vendor_bills vb ON ips.source_type = 'vendor_bill' AND ips.source_id = vb.id
LEFT JOIN customers c ON ci.customer_id = c.id
LEFT JOIN suppliers s ON vb.supplier_id = s.id
WHERE ips.deleted_at IS NULL
    AND ips.status != 'paid'
ORDER BY ips.due_date, ips.amount_due DESC;

COMMENT ON VIEW accounting.view_payment_schedules IS 'Upcoming payment schedules from payment terms with aging';

-- =====================================================
-- VIEW 5: Analytic Report (Project/Department Costs)
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_analytic_report AS
SELECT
    jel.organization_id,
    o.organization_name,
    aa.id as analytic_account_id,
    aa.account_code as analytic_code,
    aa.account_name as analytic_name,
    ap_plan.plan_name as analytic_plan,
    coa.account_code,
    coa.account_name,
    at.type_name as account_type,
    fy.fiscal_year,
    per.period_name,
    SUM(jel.debit_amount) as total_debits,
    SUM(jel.credit_amount) as total_credits,
    SUM(jel.debit_amount) - SUM(jel.credit_amount) as net_amount,
    COUNT(DISTINCT je.id) as transaction_count
FROM journal_entry_lines jel
INNER JOIN organizations o ON jel.organization_id = o.id
INNER JOIN analytic_accounts aa ON jel.analytic_account_id = aa.id
LEFT JOIN analytic_plans ap_plan ON aa.analytic_plan_id = ap_plan.id
INNER JOIN chart_of_accounts coa ON jel.account_id = coa.id
INNER JOIN account_types at ON coa.account_type_id = at.id
INNER JOIN journal_entries je ON jel.journal_entry_id = je.id
LEFT JOIN fiscal_years fy ON je.fiscal_year_id = fy.id
LEFT JOIN accounting_periods per ON je.accounting_period_id = per.id
WHERE jel.deleted_at IS NULL
    AND je.is_posted = true
    AND je.deleted_at IS NULL
GROUP BY jel.organization_id, o.organization_name, aa.id, aa.account_code,
         aa.account_name, ap_plan.plan_name, coa.account_code, coa.account_name,
         at.type_name, fy.fiscal_year, per.period_name
ORDER BY ap_plan.plan_name, aa.account_code, fy.fiscal_year, per.period_name;

COMMENT ON VIEW accounting.view_analytic_report IS 'Costs and revenues by analytic dimension (projects, departments, regions)';

-- =====================================================
-- VIEW 6: Deferred Revenue/Expense Report
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_deferrals_report AS
WITH deferred_revenue AS (
    SELECT
        drc.organization_id,
        'revenue' as deferral_type,
        drc.contract_name,
        drc.total_deferred_amount as total_amount,
        drc.recognized_amount,
        drc.total_deferred_amount - drc.recognized_amount as remaining_amount,
        drc.start_date,
        drc.end_date,
        drc.recognition_method,
        drc.status,
        drs.recognition_date as next_recognition_date,
        drs.recognition_amount as next_recognition_amount,
        coa_def.account_name as deferred_account,
        coa_rev.account_name as recognition_account
    FROM deferred_revenue_contracts drc
    LEFT JOIN deferred_revenue_schedule drs ON drc.id = drs.contract_id AND drs.status = 'pending'
    LEFT JOIN chart_of_accounts coa_def ON drc.deferred_account_id = coa_def.id
    LEFT JOIN chart_of_accounts coa_rev ON drc.revenue_account_id = coa_rev.id
    WHERE drc.deleted_at IS NULL
        AND drc.status = 'active'
),
deferred_expense AS (
    SELECT
        dec.organization_id,
        'expense' as deferral_type,
        dec.contract_name,
        dec.total_deferred_amount as total_amount,
        dec.recognized_amount,
        dec.total_deferred_amount - dec.recognized_amount as remaining_amount,
        dec.start_date,
        dec.end_date,
        dec.recognition_method,
        dec.status,
        des.recognition_date as next_recognition_date,
        des.recognition_amount as next_recognition_amount,
        coa_def.account_name as deferred_account,
        coa_exp.account_name as recognition_account
    FROM deferred_expense_contracts dec
    LEFT JOIN deferred_expense_schedule des ON dec.id = des.contract_id AND des.status = 'pending'
    LEFT JOIN chart_of_accounts coa_def ON dec.deferred_account_id = coa_def.id
    LEFT JOIN chart_of_accounts coa_exp ON dec.expense_account_id = coa_exp.id
    WHERE dec.deleted_at IS NULL
        AND dec.status = 'active'
)
SELECT * FROM deferred_revenue
UNION ALL
SELECT * FROM deferred_expense
ORDER BY deferral_type, next_recognition_date;

COMMENT ON VIEW accounting.view_deferrals_report IS 'Deferred revenue and expense contracts with recognition schedules';

-- =====================================================
-- VIEW 7: Bank Reconciliation Status
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_bank_reconciliation_status AS
SELECT
    bs.organization_id,
    o.organization_name,
    ba.bank_name,
    ba.account_number,
    bs.statement_number,
    bs.statement_date,
    bs.period_start_date,
    bs.period_end_date,
    bs.opening_balance,
    bs.closing_balance,
    bs.status as statement_status,
    COUNT(bsl.id) as total_lines,
    COUNT(CASE WHEN bsl.status = 'matched' THEN 1 END) as matched_lines,
    COUNT(CASE WHEN bsl.status = 'unmatched' THEN 1 END) as unmatched_lines,
    COUNT(CASE WHEN bsl.status = 'partial_match' THEN 1 END) as partial_matched_lines,
    SUM(bsl.amount) as total_transactions,
    SUM(CASE WHEN bsl.status = 'matched' THEN bsl.amount ELSE 0 END) as matched_amount,
    SUM(CASE WHEN bsl.status = 'unmatched' THEN bsl.amount ELSE 0 END) as unmatched_amount,
    ROUND(
        COUNT(CASE WHEN bsl.status = 'matched' THEN 1 END)::NUMERIC /
        NULLIF(COUNT(bsl.id), 0) * 100,
        2
    ) as reconciliation_percentage
FROM bank_statements bs
INNER JOIN organizations o ON bs.organization_id = o.id
INNER JOIN bank_accounts ba ON bs.bank_account_id = ba.id
LEFT JOIN bank_statement_lines bsl ON bs.id = bsl.bank_statement_id AND bsl.deleted_at IS NULL
WHERE bs.deleted_at IS NULL
GROUP BY bs.organization_id, o.organization_name, ba.bank_name, ba.account_number,
         bs.statement_number, bs.statement_date, bs.period_start_date, bs.period_end_date,
         bs.opening_balance, bs.closing_balance, bs.status
ORDER BY bs.statement_date DESC;

COMMENT ON VIEW accounting.view_bank_reconciliation_status IS 'Bank statement reconciliation progress and status';

-- =====================================================
-- VIEW 8: Budget vs Actual Report
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_budget_vs_actual AS
WITH actuals AS (
    SELECT
        gl.organization_id,
        gl.account_id,
        gl.analytic_account_id,
        fy.fiscal_year,
        ap.id as accounting_period_id,
        ap.period_name,
        SUM(CASE
            WHEN at.normal_balance = 'debit' THEN gl.debit_amount - gl.credit_amount
            ELSE gl.credit_amount - gl.debit_amount
        END) as actual_amount
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    LEFT JOIN fiscal_years fy ON gl.fiscal_year_id = fy.id
    LEFT JOIN accounting_periods ap ON gl.accounting_period_id = ap.id
    WHERE gl.deleted_at IS NULL
    GROUP BY gl.organization_id, gl.account_id, gl.analytic_account_id, fy.fiscal_year, ap.id, ap.period_name
)
SELECT
    b.organization_id,
    o.organization_name,
    b.budget_code,
    b.budget_name,
    b.budget_type,
    b.status as budget_status,
    fy.fiscal_year,
    ap.period_name,
    coa.account_code,
    coa.account_name,
    aa.account_code as analytic_code,
    aa.account_name as analytic_name,
    bl.planned_amount as budget_amount,
    COALESCE(act.actual_amount, 0) as actual_amount,
    bl.planned_amount - COALESCE(act.actual_amount, 0) as variance,
    CASE
        WHEN bl.planned_amount != 0 THEN
            ROUND((COALESCE(act.actual_amount, 0) / bl.planned_amount) * 100, 2)
        ELSE NULL
    END as percentage_of_budget,
    CASE
        WHEN bl.planned_amount < COALESCE(act.actual_amount, 0) THEN 'Over Budget'
        WHEN bl.planned_amount > COALESCE(act.actual_amount, 0) THEN 'Under Budget'
        ELSE 'On Budget'
    END as status
FROM budget_lines bl
INNER JOIN budgets b ON bl.budget_id = b.id
INNER JOIN organizations o ON b.organization_id = o.id
LEFT JOIN fiscal_years fy ON b.fiscal_year_id = fy.id
LEFT JOIN accounting_periods ap ON bl.accounting_period_id = ap.id
LEFT JOIN chart_of_accounts coa ON bl.account_id = coa.id
LEFT JOIN analytic_accounts aa ON bl.analytic_account_id = aa.id
LEFT JOIN actuals act ON
    act.organization_id = b.organization_id
    AND (act.account_id = bl.account_id OR (act.account_id IS NULL AND bl.account_id IS NULL))
    AND (act.analytic_account_id = bl.analytic_account_id OR (act.analytic_account_id IS NULL AND bl.analytic_account_id IS NULL))
    AND (act.accounting_period_id = bl.accounting_period_id OR (act.accounting_period_id IS NULL AND bl.accounting_period_id IS NULL))
WHERE bl.deleted_at IS NULL
    AND b.deleted_at IS NULL
ORDER BY b.budget_code, fy.fiscal_year, ap.period_name, coa.account_code, aa.account_code;

COMMENT ON VIEW accounting.view_budget_vs_actual IS 'Budget vs actual comparison with variance analysis';

-- =====================================================
-- VIEW 9: Fiscal Position Usage Report
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_fiscal_position_usage AS
SELECT
    fp.organization_id,
    o.organization_name,
    fp.position_code,
    fp.position_name,
    fp.country_id,
    COUNT(DISTINCT ci.id) as customer_invoices_count,
    SUM(ci.total_amount) as total_invoice_amount,
    COUNT(DISTINCT vb.id) as vendor_bills_count,
    SUM(vb.total_amount) as total_bill_amount,
    COUNT(DISTINCT fptm.id) as tax_mappings_count
FROM fiscal_positions fp
INNER JOIN organizations o ON fp.organization_id = o.id
LEFT JOIN customers c ON c.fiscal_position_id = fp.id
LEFT JOIN customer_invoices ci ON c.id = ci.customer_id AND ci.deleted_at IS NULL
LEFT JOIN suppliers s ON s.fiscal_position_id = fp.id
LEFT JOIN vendor_bills vb ON s.id = vb.supplier_id AND vb.deleted_at IS NULL
LEFT JOIN fiscal_position_tax_mappings fptm ON fp.id = fptm.fiscal_position_id AND fptm.deleted_at IS NULL
WHERE fp.deleted_at IS NULL
GROUP BY fp.organization_id, o.organization_name, fp.position_code, fp.position_name, fp.country_id
ORDER BY fp.position_code;

COMMENT ON VIEW accounting.view_fiscal_position_usage IS 'Fiscal position usage statistics for invoices and bills';

-- =====================================================
-- VIEW 10: Localization Package Summary
-- =====================================================

CREATE OR REPLACE VIEW accounting.view_localization_summary AS
SELECT
    lp.package_code,
    lp.package_name,
    lp.country_code,
    lp.version,
    lp.is_active,
    COUNT(DISTINCT o.id) as organizations_using,
    COUNT(DISTINCT trd.id) as tax_reports_defined,
    STRING_AGG(DISTINCT o.organization_name, ', ' ORDER BY o.organization_name) as organization_names
FROM localization_packages lp
LEFT JOIN organizations o ON o.localization_package_id = lp.id
LEFT JOIN tax_report_definitions trd ON trd.localization_package_id = lp.id AND trd.deleted_at IS NULL
WHERE lp.deleted_at IS NULL
GROUP BY lp.package_code, lp.package_name, lp.country_code, lp.version, lp.is_active
ORDER BY lp.country_code, lp.package_code;

COMMENT ON VIEW accounting.view_localization_summary IS 'Localization package usage and tax report availability';

\echo '';
\echo '==========================================';
\echo 'Odoo-Style Report Views Created';
\echo '==========================================';
\echo '  1. ✓ Journal Entries by Journal';
\echo '  2. ✓ Tax Report (Collected vs Paid)';
\echo '  3. ✓ Multi-Currency Summary with FX';
\echo '  4. ✓ Payment Schedules with Aging';
\echo '  5. ✓ Analytic Report (Projects/Departments)';
\echo '  6. ✓ Deferrals Report';
\echo '  7. ✓ Bank Reconciliation Status';
\echo '  8. ✓ Budget vs Actual';
\echo '  9. ✓ Fiscal Position Usage';
\echo ' 10. ✓ Localization Summary';
\echo '==========================================';
\echo '';

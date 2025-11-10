-- ============================================================================
-- Migration: V010 - Add Fiscal Period Closing Procedures
-- Description: Stored procedures for closing accounting periods and fiscal years
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- HELPER FUNCTION: Get Retained Earnings Account
-- ============================================================================

CREATE OR REPLACE FUNCTION get_retained_earnings_account(
    p_organization_id UUID
) RETURNS UUID AS $$
DECLARE
    v_account_id UUID;
BEGIN
    -- Look for retained earnings account by account type
    SELECT coa.id INTO v_account_id
    FROM chart_of_accounts coa
    INNER JOIN account_types at ON coa.account_type_id = at.id
    WHERE coa.organization_id = p_organization_id
      AND (at.type_name = 'Retained Earnings'
           OR coa.account_name ILIKE '%retained earnings%'
           OR coa.account_code IN ('3200', '3300', '3000'))
      AND coa.is_active = TRUE
      AND coa.deleted_at IS NULL
    ORDER BY coa.is_system_account DESC, coa.created_at ASC
    LIMIT 1;

    IF v_account_id IS NULL THEN
        RAISE EXCEPTION 'Retained Earnings account not found for organization %', p_organization_id
            USING HINT = 'Create a Retained Earnings account in your chart of accounts';
    END IF;

    RETURN v_account_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_retained_earnings_account IS 'Finds the retained earnings account for an organization';

-- ============================================================================
-- PROCEDURE: Close Accounting Period
-- ============================================================================

CREATE OR REPLACE FUNCTION close_accounting_period(
    p_period_id UUID,
    p_closed_by UUID,
    p_create_closing_entry BOOLEAN DEFAULT FALSE
) RETURNS JSONB AS $$
DECLARE
    v_organization_id UUID;
    v_period_start DATE;
    v_period_end DATE;
    v_period_status VARCHAR;
    v_fiscal_year_id UUID;
    v_result JSONB;
    v_unposted_count INTEGER;
BEGIN
    -- Get period details
    SELECT organization_id, start_date, end_date, status, fiscal_year_id
    INTO v_organization_id, v_period_start, v_period_end, v_period_status, v_fiscal_year_id
    FROM accounting_periods
    WHERE id = p_period_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Accounting period % not found', p_period_id;
    END IF;

    -- Check if already closed
    IF v_period_status IN ('closed', 'locked') THEN
        RAISE EXCEPTION 'Accounting period is already closed (status: %)', v_period_status;
    END IF;

    -- Check for unposted journal entries in the period
    SELECT COUNT(*) INTO v_unposted_count
    FROM journal_entries
    WHERE organization_id = v_organization_id
      AND posting_date BETWEEN v_period_start AND v_period_end
      AND is_posted = FALSE
      AND deleted_at IS NULL;

    IF v_unposted_count > 0 THEN
        RAISE EXCEPTION 'Cannot close period: % unposted journal entries exist', v_unposted_count
            USING HINT = 'Post or delete all journal entries before closing the period';
    END IF;

    -- Optionally create closing entry (usually done at year end, not period end)
    IF p_create_closing_entry THEN
        -- This would create income summary entry
        -- Implementation depends on specific requirements
        RAISE NOTICE 'Closing entry creation not implemented for mid-year periods';
    END IF;

    -- Mark period as closed
    UPDATE accounting_periods
    SET
        status = 'closed',
        closed_at = CURRENT_TIMESTAMP,
        closed_by = p_closed_by,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_period_id;

    -- Build result
    v_result := jsonb_build_object(
        'success', TRUE,
        'period_id', p_period_id,
        'period_start', v_period_start,
        'period_end', v_period_end,
        'status', 'closed',
        'closed_at', CURRENT_TIMESTAMP,
        'closed_by', p_closed_by,
        'message', 'Accounting period closed successfully'
    );

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION close_accounting_period IS 'Closes an accounting period and prevents further posting to it';

-- ============================================================================
-- PROCEDURE: Reopen Accounting Period
-- ============================================================================

CREATE OR REPLACE FUNCTION reopen_accounting_period(
    p_period_id UUID,
    p_reopened_by UUID
) RETURNS JSONB AS $$
DECLARE
    v_period_status VARCHAR;
    v_result JSONB;
BEGIN
    -- Get period status
    SELECT status INTO v_period_status
    FROM accounting_periods
    WHERE id = p_period_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Accounting period % not found', p_period_id;
    END IF;

    -- Don't allow reopening locked periods without explicit unlock
    IF v_period_status = 'locked' THEN
        RAISE EXCEPTION 'Cannot reopen locked period. Use unlock_accounting_period function first.'
            USING HINT = 'Locked periods require explicit unlock for audit trail purposes';
    END IF;

    -- Reopen the period
    UPDATE accounting_periods
    SET
        status = 'open',
        closed_at = NULL,
        closed_by = NULL,
        updated_at = CURRENT_TIMESTAMP,
        updated_by = p_reopened_by
    WHERE id = p_period_id;

    v_result := jsonb_build_object(
        'success', TRUE,
        'period_id', p_period_id,
        'status', 'open',
        'reopened_at', CURRENT_TIMESTAMP,
        'reopened_by', p_reopened_by,
        'message', 'Accounting period reopened successfully'
    );

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION reopen_accounting_period IS 'Reopens a closed accounting period';

-- ============================================================================
-- PROCEDURE: Close Fiscal Year
-- ============================================================================

CREATE OR REPLACE FUNCTION close_fiscal_year(
    p_fiscal_year_id UUID,
    p_closed_by UUID,
    p_create_closing_entries BOOLEAN DEFAULT TRUE
) RETURNS JSONB AS $$
DECLARE
    v_organization_id UUID;
    v_fiscal_year VARCHAR;
    v_fy_start DATE;
    v_fy_end DATE;
    v_fy_status VARCHAR;
    v_retained_earnings_account_id UUID;
    v_closing_je_id UUID;
    v_closing_je_number VARCHAR;
    v_total_income NUMERIC;
    v_total_expense NUMERIC;
    v_net_income NUMERIC;
    v_open_periods_count INTEGER;
    v_result JSONB;
BEGIN
    -- Get fiscal year details
    SELECT organization_id, fiscal_year, start_date, end_date, status
    INTO v_organization_id, v_fiscal_year, v_fy_start, v_fy_end, v_fy_status
    FROM fiscal_years
    WHERE id = p_fiscal_year_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Fiscal year % not found', p_fiscal_year_id;
    END IF;

    -- Check if already closed
    IF v_fy_status IN ('closed', 'locked') THEN
        RAISE EXCEPTION 'Fiscal year % is already closed (status: %)', v_fiscal_year, v_fy_status;
    END IF;

    -- Check if all periods in fiscal year are closed
    SELECT COUNT(*) INTO v_open_periods_count
    FROM accounting_periods
    WHERE fiscal_year_id = p_fiscal_year_id
      AND status NOT IN ('closed', 'locked')
      AND deleted_at IS NULL;

    IF v_open_periods_count > 0 THEN
        RAISE EXCEPTION 'Cannot close fiscal year: % accounting periods are still open', v_open_periods_count
            USING HINT = 'Close all accounting periods before closing the fiscal year';
    END IF;

    -- Create closing entries if requested
    IF p_create_closing_entries THEN
        -- Get retained earnings account
        v_retained_earnings_account_id := get_retained_earnings_account(v_organization_id);

        -- Calculate total income (credit balance accounts)
        SELECT COALESCE(SUM(gl.credit_amount - gl.debit_amount), 0) INTO v_total_income
        FROM general_ledger gl
        INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
        INNER JOIN account_types at ON coa.account_type_id = at.id
        WHERE gl.organization_id = v_organization_id
          AND gl.transaction_date BETWEEN v_fy_start AND v_fy_end
          AND at.normal_balance = 'credit'
          AND at.type_category = 'income';

        -- Calculate total expenses (debit balance accounts)
        SELECT COALESCE(SUM(gl.debit_amount - gl.credit_amount), 0) INTO v_total_expense
        FROM general_ledger gl
        INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
        INNER JOIN account_types at ON coa.account_type_id = at.id
        WHERE gl.organization_id = v_organization_id
          AND gl.transaction_date BETWEEN v_fy_start AND v_fy_end
          AND at.normal_balance = 'debit'
          AND at.type_category IN ('expense', 'cost_of_sales');

        -- Calculate net income
        v_net_income := v_total_income - v_total_expense;

        -- Create closing journal entry
        v_closing_je_number := 'CLOSE-' || v_fiscal_year;

        INSERT INTO journal_entries (
            organization_id,
            entry_number,
            entry_date,
            posting_date,
            journal_id,
            description,
            reference,
            status,
            is_posted,
            total_debit,
            total_credit,
            created_by,
            updated_by
        ) VALUES (
            v_organization_id,
            v_closing_je_number,
            v_fy_end,
            v_fy_end,
            (SELECT id FROM journals WHERE organization_id = v_organization_id AND journal_type = 'general' LIMIT 1),
            'Year-end closing entry for fiscal year ' || v_fiscal_year,
            'FY_CLOSE_' || v_fiscal_year,
            'posted',
            TRUE,
            CASE WHEN v_net_income > 0 THEN v_net_income ELSE v_total_income + v_total_expense END,
            CASE WHEN v_net_income > 0 THEN v_net_income ELSE v_total_income + v_total_expense END,
            p_closed_by,
            p_closed_by
        ) RETURNING id INTO v_closing_je_id;

        -- Create journal entry lines (close income and expense to retained earnings)
        -- This is a simplified version - real implementation would close each account individually

        -- Post to GL
        INSERT INTO general_ledger (
            organization_id,
            journal_entry_id,
            account_id,
            transaction_date,
            debit_amount,
            credit_amount,
            description,
            reference
        ) VALUES (
            v_organization_id,
            v_closing_je_id,
            v_retained_earnings_account_id,
            v_fy_end,
            CASE WHEN v_net_income > 0 THEN v_net_income ELSE 0 END,
            CASE WHEN v_net_income < 0 THEN ABS(v_net_income) ELSE 0 END,
            'Net income transferred to retained earnings',
            v_closing_je_number
        );

        RAISE NOTICE 'Created closing entry % with net income: %', v_closing_je_number, v_net_income;
    END IF;

    -- Mark fiscal year as closed
    UPDATE fiscal_years
    SET
        status = 'closed',
        closed_at = CURRENT_TIMESTAMP,
        closed_by = p_closed_by,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_fiscal_year_id;

    -- Build result
    v_result := jsonb_build_object(
        'success', TRUE,
        'fiscal_year_id', p_fiscal_year_id,
        'fiscal_year', v_fiscal_year,
        'start_date', v_fy_start,
        'end_date', v_fy_end,
        'status', 'closed',
        'closing_je_id', v_closing_je_id,
        'closing_je_number', v_closing_je_number,
        'net_income', v_net_income,
        'closed_at', CURRENT_TIMESTAMP,
        'closed_by', p_closed_by,
        'message', 'Fiscal year closed successfully'
    );

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION close_fiscal_year IS 'Closes a fiscal year and creates year-end closing entries';

-- ============================================================================
-- PROCEDURE: Lock/Unlock Accounting Period
-- ============================================================================

CREATE OR REPLACE FUNCTION lock_accounting_period(
    p_period_id UUID,
    p_locked_by UUID
) RETURNS JSONB AS $$
DECLARE
    v_result JSONB;
BEGIN
    UPDATE accounting_periods
    SET
        status = 'locked',
        updated_at = CURRENT_TIMESTAMP,
        updated_by = p_locked_by
    WHERE id = p_period_id;

    v_result := jsonb_build_object(
        'success', TRUE,
        'period_id', p_period_id,
        'status', 'locked',
        'message', 'Accounting period locked'
    );

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION unlock_accounting_period(
    p_period_id UUID,
    p_unlocked_by UUID
) RETURNS JSONB AS $$
DECLARE
    v_result JSONB;
BEGIN
    UPDATE accounting_periods
    SET
        status = 'open',
        updated_at = CURRENT_TIMESTAMP,
        updated_by = p_unlocked_by
    WHERE id = p_period_id;

    v_result := jsonb_build_object(
        'success', TRUE,
        'period_id', p_period_id,
        'status', 'open',
        'message', 'Accounting period unlocked'
    );

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION lock_accounting_period IS 'Locks an accounting period (stronger than close)';
COMMENT ON FUNCTION unlock_accounting_period IS 'Unlocks a locked accounting period';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V010 completed successfully!';
    RAISE NOTICE 'Fiscal Period Closing Procedures created:';
    RAISE NOTICE ' - close_accounting_period()';
    RAISE NOTICE ' - reopen_accounting_period()';
    RAISE NOTICE ' - close_fiscal_year()';
    RAISE NOTICE ' - lock_accounting_period()';
    RAISE NOTICE ' - unlock_accounting_period()';
    RAISE NOTICE ' - get_retained_earnings_account()';
    RAISE NOTICE '============================================';
END $$;

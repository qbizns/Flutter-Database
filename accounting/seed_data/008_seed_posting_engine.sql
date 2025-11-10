-- ============================================================================
-- Seed Data: 008 - Posting Engine Configuration
-- Description: Proof-of-concept for configuration-driven posting engine
--              Includes concepts, profiles, document types, rules, and validators
-- Dependencies: Requires accounting V011, V012, V013
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

DO $$
DECLARE
    v_org_id UUID;
    v_user_id UUID;
    v_profile_id UUID;
    v_fiscal_year_id UUID;
    -- GL Accounts
    v_ar_account_id UUID;
    v_cash_account_id UUID;
    v_bank_account_id UUID;
    v_revenue_account_id UUID;
    v_cogs_account_id UUID;
    v_inventory_account_id UUID;
    v_discount_account_id UUID;
    v_tax_liability_account_id UUID;
    v_payroll_expense_account_id UUID;
    v_payroll_liability_account_id UUID;
    v_retained_earnings_account_id UUID;
    -- Account types
    v_asset_type_id UUID;
    v_liability_type_id UUID;
    v_revenue_type_id UUID;
    v_expense_type_id UUID;
    -- Document types
    v_pos_sale_doctype_id UUID;
    v_payroll_run_doctype_id UUID;
    v_vendor_bill_doctype_id UUID;
    v_check_issue_doctype_id UUID;
    -- Profile documents
    v_profile_doc_pos_sale_id UUID;
    v_profile_doc_payroll_id UUID;
    v_profile_doc_vendor_bill_id UUID;
    v_profile_doc_check_id UUID;
    -- Posting rules
    v_rule_pos_cash_id UUID;
    v_rule_pos_credit_id UUID;
    v_rule_payroll_id UUID;
    v_rule_vendor_bill_id UUID;
    v_rule_check_issue_id UUID;
BEGIN
    -- Get organization and user
    SELECT id INTO v_org_id FROM organizations LIMIT 1;
    SELECT id INTO v_user_id FROM users LIMIT 1;

    IF v_org_id IS NULL THEN
        RAISE EXCEPTION 'No organization found. Run core seed data first.';
    END IF;

    RAISE NOTICE 'Using Organization ID: %', v_org_id;

    -- Get active fiscal year
    SELECT id INTO v_fiscal_year_id
    FROM fiscal_years
    WHERE organization_id = v_org_id
      AND status = 'open'
    ORDER BY start_date DESC
    LIMIT 1;

    -- ========================================================================
    -- GET GL ACCOUNTS (from existing chart of accounts)
    -- ========================================================================

    RAISE NOTICE 'Fetching GL accounts from chart of accounts...';

    -- AR account (1100)
    SELECT id INTO v_ar_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '1100' OR account_name ILIKE '%accounts receivable%')
    LIMIT 1;

    -- Cash account (1000)
    SELECT id INTO v_cash_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '1000' OR account_name ILIKE '%cash%')
    LIMIT 1;

    -- Bank account (1010)
    SELECT id INTO v_bank_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '1010' OR account_code LIKE '102%' OR account_name ILIKE '%bank%')
    LIMIT 1;

    -- Revenue account (4000)
    SELECT id INTO v_revenue_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '4000' OR account_code LIKE '4%' OR account_name ILIKE '%revenue%' OR account_name ILIKE '%sales%')
    LIMIT 1;

    -- COGS account (5000)
    SELECT id INTO v_cogs_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '5000' OR account_code LIKE '5%' OR account_name ILIKE '%cost of sales%' OR account_name ILIKE '%cogs%')
    LIMIT 1;

    -- Inventory account (1300)
    SELECT id INTO v_inventory_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '1300' OR account_name ILIKE '%inventory%')
    LIMIT 1;

    -- Discount expense account (5100)
    SELECT id INTO v_discount_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '5100' OR account_name ILIKE '%discount%')
    LIMIT 1;

    -- Tax liability account (2100)
    SELECT id INTO v_tax_liability_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '2100' OR account_name ILIKE '%tax payable%' OR account_name ILIKE '%vat%')
    LIMIT 1;

    -- Payroll expense account (5200)
    SELECT id INTO v_payroll_expense_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '5200' OR account_name ILIKE '%payroll expense%' OR account_name ILIKE '%salaries%')
    LIMIT 1;

    -- Payroll liability account (2200)
    SELECT id INTO v_payroll_liability_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '2200' OR account_name ILIKE '%payroll liability%' OR account_name ILIKE '%wages payable%')
    LIMIT 1;

    -- Retained earnings account (3400)
    SELECT id INTO v_retained_earnings_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code = '3400' OR account_name ILIKE '%retained earnings%')
    LIMIT 1;

    -- Get account types
    SELECT id INTO v_asset_type_id
    FROM account_types
    WHERE type_code = 'ASSET'
    LIMIT 1;

    SELECT id INTO v_liability_type_id
    FROM account_types
    WHERE type_code = 'LIABILITY'
    LIMIT 1;

    SELECT id INTO v_revenue_type_id
    FROM account_types
    WHERE type_code = 'REVENUE'
    LIMIT 1;

    SELECT id INTO v_expense_type_id
    FROM account_types
    WHERE type_code = 'EXPENSE'
    LIMIT 1;

    -- ========================================================================
    -- SEED: Posting Concepts
    -- ========================================================================

    RAISE NOTICE 'Seeding posting concepts...';

    INSERT INTO posting_concepts (
        concept_key,
        default_label,
        default_description,
        expected_account_type_id,
        normal_side,
        concept_category,
        is_system
    ) VALUES
        -- Assets
        ('AR', 'Accounts Receivable', 'Customer balances owed to the company', v_asset_type_id, 'debit', 'asset', TRUE),
        ('CASH', 'Cash', 'Cash on hand and in transit', v_asset_type_id, 'debit', 'asset', TRUE),
        ('BANK', 'Bank Account', 'Bank account balances', v_asset_type_id, 'debit', 'asset', TRUE),
        ('INVENTORY', 'Inventory', 'Stock on hand at cost', v_asset_type_id, 'debit', 'asset', TRUE),

        -- Liabilities
        ('AP', 'Accounts Payable', 'Vendor balances owed by the company', v_liability_type_id, 'credit', 'liability', TRUE),
        ('TAX_OUTPUT', 'Tax Collected (Output VAT)', 'Sales tax or VAT collected from customers', v_liability_type_id, 'credit', 'liability', TRUE),
        ('TAX_INPUT', 'Tax Paid (Input VAT)', 'Sales tax or VAT paid to vendors', v_asset_type_id, 'debit', 'asset', TRUE),
        ('PAYROLL_LIABILITY', 'Payroll Payable', 'Wages and salaries payable to employees', v_liability_type_id, 'credit', 'liability', TRUE),
        ('GIFT_CARD_LIABILITY', 'Gift Card Liability', 'Unearned revenue from gift card sales', v_liability_type_id, 'credit', 'liability', TRUE),

        -- Revenue
        ('REVENUE', 'Sales Revenue', 'Income from sale of goods or services', v_revenue_type_id, 'credit', 'revenue', TRUE),
        ('SERVICE_REVENUE', 'Service Revenue', 'Income from services rendered', v_revenue_type_id, 'credit', 'revenue', TRUE),
        ('OTHER_REVENUE', 'Other Revenue', 'Miscellaneous income', v_revenue_type_id, 'credit', 'revenue', TRUE),

        -- Expenses
        ('COGS', 'Cost of Goods Sold', 'Cost of inventory sold during the period', v_expense_type_id, 'debit', 'expense', TRUE),
        ('PAYROLL_EXPENSE', 'Payroll Expense', 'Wages and salaries expense', v_expense_type_id, 'debit', 'expense', TRUE),
        ('DISCOUNT_EXPENSE', 'Discount Expense', 'Sales discounts given to customers', v_expense_type_id, 'debit', 'expense', TRUE),
        ('ROUNDING_EXPENSE', 'Rounding Adjustment', 'Minor rounding differences', v_expense_type_id, 'debit', 'expense', TRUE),

        -- Equity
        ('RETAINED_EARNINGS', 'Retained Earnings', 'Accumulated net income', v_liability_type_id, 'credit', 'equity', TRUE);

    RAISE NOTICE 'Created % posting concepts', (SELECT COUNT(*) FROM posting_concepts);

    -- ========================================================================
    -- SEED: Posting Account Mappings (Concept → GL Account)
    -- ========================================================================

    RAISE NOTICE 'Seeding posting account mappings...';

    INSERT INTO posting_account_mappings (
        organization_id,
        concept_key,
        account_id,
        location_id,
        product_category_id,
        customer_type_code,
        vendor_type_code,
        is_default,
        priority,
        description,
        created_by,
        updated_by
    ) VALUES
        -- Asset concepts
        (v_org_id, 'AR', v_ar_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default AR account', v_user_id, v_user_id),
        (v_org_id, 'CASH', v_cash_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default cash account', v_user_id, v_user_id),
        (v_org_id, 'BANK', v_bank_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default bank account', v_user_id, v_user_id),
        (v_org_id, 'INVENTORY', v_inventory_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default inventory account', v_user_id, v_user_id),

        -- Liability concepts
        (v_org_id, 'TAX_OUTPUT', v_tax_liability_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default tax liability account', v_user_id, v_user_id),
        (v_org_id, 'PAYROLL_LIABILITY', v_payroll_liability_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default payroll liability account', v_user_id, v_user_id),

        -- Revenue concepts
        (v_org_id, 'REVENUE', v_revenue_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default revenue account', v_user_id, v_user_id),

        -- Expense concepts
        (v_org_id, 'COGS', v_cogs_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default COGS account', v_user_id, v_user_id),
        (v_org_id, 'PAYROLL_EXPENSE', v_payroll_expense_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default payroll expense account', v_user_id, v_user_id),
        (v_org_id, 'DISCOUNT_EXPENSE', v_discount_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default discount expense account', v_user_id, v_user_id),

        -- Equity concepts
        (v_org_id, 'RETAINED_EARNINGS', v_retained_earnings_account_id, NULL, NULL, NULL, NULL, TRUE, 0, 'Default retained earnings account', v_user_id, v_user_id);

    RAISE NOTICE 'Created % posting account mappings', (SELECT COUNT(*) FROM posting_account_mappings WHERE organization_id = v_org_id);

    -- ========================================================================
    -- SEED: Posting Profile
    -- ========================================================================

    RAISE NOTICE 'Seeding posting profile...';

    INSERT INTO posting_profiles (
        organization_id,
        code,
        name,
        description,
        is_default,
        is_active,
        default_fiscal_year_id,
        notes,
        created_by,
        updated_by
    ) VALUES (
        v_org_id,
        'DEFAULT_PROFILE',
        'Default Posting Profile',
        'Standard posting configuration for all business documents',
        TRUE,
        TRUE,
        v_fiscal_year_id,
        'Automatically created posting profile with standard rules',
        v_user_id,
        v_user_id
    )
    RETURNING id INTO v_profile_id;

    RAISE NOTICE 'Created posting profile: %', v_profile_id;

    -- ========================================================================
    -- SEED: Posting Document Types
    -- ========================================================================

    RAISE NOTICE 'Seeding posting document types...';

    INSERT INTO posting_document_types (
        code,
        name,
        description,
        source_schema,
        source_table,
        source_pk_column,
        category,
        is_active,
        is_system
    ) VALUES
        ('POS_SALE', 'POS Sale', 'Point of sale transaction', 'public', 'sales', 'id', 'sales', TRUE, TRUE),
        ('PAYROLL_RUN', 'Payroll Run', 'Employee payroll batch', 'public', 'payroll_runs', 'id', 'payroll', TRUE, TRUE),
        ('VENDOR_BILL', 'Vendor Bill', 'Accounts payable invoice', 'public', 'vendor_bills', 'id', 'purchases', TRUE, TRUE),
        ('CHECK_ISSUE', 'Check Issue', 'Check payment issued', 'public', 'checks', 'id', 'banking', TRUE, TRUE),
        ('EXPENSE_CLAIM', 'Expense Claim', 'Employee expense reimbursement', 'public', 'expense_claims', 'id', 'expenses', TRUE, TRUE),
        ('INVENTORY_ADJUSTMENT', 'Inventory Adjustment', 'Stock count adjustment', 'public', 'inventory_adjustments', 'id', 'inventory', TRUE, TRUE)
    RETURNING id INTO v_pos_sale_doctype_id, v_payroll_run_doctype_id, v_vendor_bill_doctype_id,
                     v_check_issue_doctype_id;

    -- Get the IDs we just inserted
    SELECT id INTO v_pos_sale_doctype_id FROM posting_document_types WHERE code = 'POS_SALE';
    SELECT id INTO v_payroll_run_doctype_id FROM posting_document_types WHERE code = 'PAYROLL_RUN';
    SELECT id INTO v_vendor_bill_doctype_id FROM posting_document_types WHERE code = 'VENDOR_BILL';
    SELECT id INTO v_check_issue_doctype_id FROM posting_document_types WHERE code = 'CHECK_ISSUE';

    RAISE NOTICE 'Created % posting document types', (SELECT COUNT(*) FROM posting_document_types);

    -- ========================================================================
    -- SEED: Posting Profile Documents
    -- ========================================================================

    RAISE NOTICE 'Linking document types to profile...';

    INSERT INTO posting_profile_documents (
        posting_profile_id,
        posting_document_type_id,
        is_active,
        notes
    ) VALUES
        (v_profile_id, v_pos_sale_doctype_id, TRUE, 'POS sales posting'),
        (v_profile_id, v_payroll_run_doctype_id, TRUE, 'Payroll posting'),
        (v_profile_id, v_vendor_bill_doctype_id, TRUE, 'Vendor bills posting'),
        (v_profile_id, v_check_issue_doctype_id, TRUE, 'Check payments posting')
    RETURNING id INTO v_profile_doc_pos_sale_id, v_profile_doc_payroll_id,
                     v_profile_doc_vendor_bill_id, v_profile_doc_check_id;

    -- Get the IDs we just inserted
    SELECT id INTO v_profile_doc_pos_sale_id
    FROM posting_profile_documents
    WHERE posting_profile_id = v_profile_id AND posting_document_type_id = v_pos_sale_doctype_id;

    SELECT id INTO v_profile_doc_payroll_id
    FROM posting_profile_documents
    WHERE posting_profile_id = v_profile_id AND posting_document_type_id = v_payroll_run_doctype_id;

    SELECT id INTO v_profile_doc_vendor_bill_id
    FROM posting_profile_documents
    WHERE posting_profile_id = v_profile_id AND posting_document_type_id = v_vendor_bill_doctype_id;

    SELECT id INTO v_profile_doc_check_id
    FROM posting_profile_documents
    WHERE posting_profile_id = v_profile_id AND posting_document_type_id = v_check_issue_doctype_id;

    RAISE NOTICE 'Created % profile-document links', (SELECT COUNT(*) FROM posting_profile_documents WHERE posting_profile_id = v_profile_id);

    -- ========================================================================
    -- SEED: Posting Rules - POS_SALE (Cash)
    -- ========================================================================

    RAISE NOTICE 'Creating POS_SALE posting rules...';

    -- Rule: POS_SALE_CASH (when payment method is CASH)
    INSERT INTO posting_rules (
        posting_profile_document_id,
        rule_code,
        rule_name,
        description,
        event,
        level,
        priority,
        condition_expression,
        is_active,
        notes,
        created_by,
        updated_by
    ) VALUES (
        v_profile_doc_pos_sale_id,
        'POS_SALE_CASH',
        'POS Sale - Cash Payment',
        'Journal entry for cash sales',
        'on_post',
        'header',
        100,
        'doc.payment_method == "CASH"',
        TRUE,
        'Debit cash, credit revenue and tax',
        v_user_id,
        v_user_id
    )
    RETURNING id INTO v_rule_pos_cash_id;

    -- Rule lines for POS_SALE_CASH
    INSERT INTO posting_rule_lines (
        posting_rule_id,
        line_no,
        side,
        concept_key,
        account_source,
        amount_source,
        amount_field_path,
        description_template,
        is_active,
        notes
    ) VALUES
        -- Line 1: Debit CASH
        (v_rule_pos_cash_id, 1, 'debit', 'CASH', 'from_mapping', 'document_field', 'doc.total_amount', 'Cash sale {doc.sale_number}', TRUE, 'Total cash received'),
        -- Line 2: Credit REVENUE
        (v_rule_pos_cash_id, 2, 'credit', 'REVENUE', 'from_mapping', 'document_field', 'doc.subtotal_amount', 'Revenue from sale {doc.sale_number}', TRUE, 'Revenue before tax'),
        -- Line 3: Credit TAX_OUTPUT
        (v_rule_pos_cash_id, 3, 'credit', 'TAX_OUTPUT', 'from_mapping', 'document_field', 'doc.tax_amount', 'Tax on sale {doc.sale_number}', TRUE, 'Sales tax collected');

    -- Rule: POS_SALE_CREDIT (when payment method is CREDIT)
    INSERT INTO posting_rules (
        posting_profile_document_id,
        rule_code,
        rule_name,
        description,
        event,
        level,
        priority,
        condition_expression,
        is_active,
        notes,
        created_by,
        updated_by
    ) VALUES (
        v_profile_doc_pos_sale_id,
        'POS_SALE_CREDIT',
        'POS Sale - Credit/Account',
        'Journal entry for credit sales',
        'on_post',
        'header',
        100,
        'doc.payment_method == "CREDIT" || doc.payment_method == "ACCOUNT"',
        TRUE,
        'Debit AR, credit revenue and tax',
        v_user_id,
        v_user_id
    )
    RETURNING id INTO v_rule_pos_credit_id;

    -- Rule lines for POS_SALE_CREDIT
    INSERT INTO posting_rule_lines (
        posting_rule_id,
        line_no,
        side,
        concept_key,
        account_source,
        amount_source,
        amount_field_path,
        description_template,
        is_active,
        notes
    ) VALUES
        -- Line 1: Debit AR
        (v_rule_pos_credit_id, 1, 'debit', 'AR', 'from_mapping', 'document_field', 'doc.total_amount', 'Credit sale {doc.sale_number}', TRUE, 'Total amount receivable'),
        -- Line 2: Credit REVENUE
        (v_rule_pos_credit_id, 2, 'credit', 'REVENUE', 'from_mapping', 'document_field', 'doc.subtotal_amount', 'Revenue from sale {doc.sale_number}', TRUE, 'Revenue before tax'),
        -- Line 3: Credit TAX_OUTPUT
        (v_rule_pos_credit_id, 3, 'credit', 'TAX_OUTPUT', 'from_mapping', 'document_field', 'doc.tax_amount', 'Tax on sale {doc.sale_number}', TRUE, 'Sales tax collected');

    RAISE NOTICE 'Created POS_SALE posting rules with % lines',
        (SELECT COUNT(*) FROM posting_rule_lines WHERE posting_rule_id IN (v_rule_pos_cash_id, v_rule_pos_credit_id));

    -- ========================================================================
    -- SEED: Posting Rules - PAYROLL_RUN
    -- ========================================================================

    RAISE NOTICE 'Creating PAYROLL_RUN posting rules...';

    INSERT INTO posting_rules (
        posting_profile_document_id,
        rule_code,
        rule_name,
        description,
        event,
        level,
        priority,
        condition_expression,
        is_active,
        notes,
        created_by,
        updated_by
    ) VALUES (
        v_profile_doc_payroll_id,
        'PAYROLL_ACCRUAL',
        'Payroll Accrual',
        'Journal entry for payroll expense accrual',
        'on_post',
        'header',
        100,
        'doc.status == "APPROVED"',
        TRUE,
        'Debit payroll expense, credit payroll liability',
        v_user_id,
        v_user_id
    )
    RETURNING id INTO v_rule_payroll_id;

    -- Rule lines for PAYROLL_ACCRUAL
    INSERT INTO posting_rule_lines (
        posting_rule_id,
        line_no,
        side,
        concept_key,
        account_source,
        amount_source,
        amount_field_path,
        description_template,
        is_active,
        notes
    ) VALUES
        -- Line 1: Debit PAYROLL_EXPENSE
        (v_rule_payroll_id, 1, 'debit', 'PAYROLL_EXPENSE', 'from_mapping', 'document_field', 'doc.gross_salary_amount', 'Payroll expense {doc.period_name}', TRUE, 'Total gross wages'),
        -- Line 2: Credit PAYROLL_LIABILITY
        (v_rule_payroll_id, 2, 'credit', 'PAYROLL_LIABILITY', 'from_mapping', 'document_field', 'doc.gross_salary_amount', 'Wages payable {doc.period_name}', TRUE, 'Total wages owed to employees');

    RAISE NOTICE 'Created PAYROLL_RUN posting rule with % lines',
        (SELECT COUNT(*) FROM posting_rule_lines WHERE posting_rule_id = v_rule_payroll_id);

    -- ========================================================================
    -- SEED: Posting Rules - VENDOR_BILL
    -- ========================================================================

    RAISE NOTICE 'Creating VENDOR_BILL posting rules...';

    INSERT INTO posting_rules (
        posting_profile_document_id,
        rule_code,
        rule_name,
        description,
        event,
        level,
        priority,
        condition_expression,
        is_active,
        notes,
        created_by,
        updated_by
    ) VALUES (
        v_profile_doc_vendor_bill_id,
        'VENDOR_BILL_POST',
        'Vendor Bill Posting',
        'Journal entry for vendor bill',
        'on_post',
        'header',
        100,
        NULL,  -- Always apply
        TRUE,
        'Debit expense/inventory, credit AP',
        v_user_id,
        v_user_id
    )
    RETURNING id INTO v_rule_vendor_bill_id;

    -- Rule lines for VENDOR_BILL_POST
    INSERT INTO posting_rule_lines (
        posting_rule_id,
        line_no,
        side,
        concept_key,
        account_source,
        amount_source,
        amount_field_path,
        description_template,
        is_active,
        notes
    ) VALUES
        -- Line 1: Debit INVENTORY (assuming inventory purchase)
        (v_rule_vendor_bill_id, 1, 'debit', 'INVENTORY', 'from_mapping', 'document_field', 'doc.subtotal_amount', 'Inventory purchase {doc.bill_number}', TRUE, 'Goods purchased'),
        -- Line 2: Debit TAX_INPUT
        (v_rule_vendor_bill_id, 2, 'debit', 'TAX_INPUT', 'from_mapping', 'document_field', 'doc.tax_amount', 'Input tax {doc.bill_number}', TRUE, 'Tax paid to vendor'),
        -- Line 3: Credit AP
        (v_rule_vendor_bill_id, 3, 'credit', 'AP', 'from_mapping', 'document_field', 'doc.total_amount', 'Vendor bill {doc.bill_number}', TRUE, 'Amount owed to vendor');

    RAISE NOTICE 'Created VENDOR_BILL posting rule with % lines',
        (SELECT COUNT(*) FROM posting_rule_lines WHERE posting_rule_id = v_rule_vendor_bill_id);

    -- ========================================================================
    -- SEED: Posting Rules - CHECK_ISSUE
    -- ========================================================================

    RAISE NOTICE 'Creating CHECK_ISSUE posting rules...';

    INSERT INTO posting_rules (
        posting_profile_document_id,
        rule_code,
        rule_name,
        description,
        event,
        level,
        priority,
        condition_expression,
        is_active,
        notes,
        created_by,
        updated_by
    ) VALUES (
        v_profile_doc_check_id,
        'CHECK_ISSUE',
        'Check Issue',
        'Journal entry when check is issued',
        'on_post',
        'header',
        100,
        'doc.status == "ISSUED"',
        TRUE,
        'Debit AP, credit bank',
        v_user_id,
        v_user_id
    )
    RETURNING id INTO v_rule_check_issue_id;

    -- Rule lines for CHECK_ISSUE
    INSERT INTO posting_rule_lines (
        posting_rule_id,
        line_no,
        side,
        concept_key,
        account_source,
        amount_source,
        amount_field_path,
        description_template,
        is_active,
        notes
    ) VALUES
        -- Line 1: Debit AP
        (v_rule_check_issue_id, 1, 'debit', 'AP', 'from_mapping', 'document_field', 'doc.amount', 'Payment to vendor {doc.payee_name}', TRUE, 'Reduce AP balance'),
        -- Line 2: Credit BANK
        (v_rule_check_issue_id, 2, 'credit', 'BANK', 'from_mapping', 'document_field', 'doc.amount', 'Check #{doc.check_number}', TRUE, 'Bank account reduction');

    RAISE NOTICE 'Created CHECK_ISSUE posting rule with % lines',
        (SELECT COUNT(*) FROM posting_rule_lines WHERE posting_rule_id = v_rule_check_issue_id);

    -- ========================================================================
    -- SEED: Posting Validation Rules
    -- ========================================================================

    RAISE NOTICE 'Seeding posting validation rules...';

    INSERT INTO posting_validation_rules (
        organization_id,
        document_type_code,
        event,
        target,
        code,
        expression,
        severity,
        is_blocking,
        message_template,
        is_active,
        notes,
        created_by,
        updated_by
    ) VALUES
        -- Global validation: Balanced journal entry
        (NULL, NULL, NULL, 'journal_entry', 'BALANCED_ENTRY',
         'abs(je.total_debit - je.total_credit) <= 0.005',
         'error', TRUE,
         'Journal entry must balance: debits ({je.total_debit}) must equal credits ({je.total_credit})',
         TRUE, 'Ensures double-entry accounting balance', v_user_id, v_user_id),

        -- Global validation: Open period check
        (NULL, NULL, NULL, 'document', 'OPEN_PERIOD',
         'period.status == "open"',
         'error', TRUE,
         'Cannot post to period {period.period_name}: period is {period.status}',
         TRUE, 'Prevents posting to closed/locked periods', v_user_id, v_user_id),

        -- Global validation: Non-zero amount
        (NULL, NULL, NULL, 'journal_line', 'NON_ZERO_AMOUNT',
         'abs(line.amount) > 0.001',
         'error', TRUE,
         'Journal entry line amount must be greater than zero',
         TRUE, 'Prevents zero-amount journal lines', v_user_id, v_user_id),

        -- Organization-specific: POS sale minimum amount
        (v_org_id, 'POS_SALE', 'on_post', 'document', 'MIN_SALE_AMOUNT',
         'doc.total_amount >= 0.01',
         'error', TRUE,
         'Sale amount must be at least 0.01',
         TRUE, 'Prevents invalid zero/negative sales', v_user_id, v_user_id),

        -- Organization-specific: Payroll requires approval
        (v_org_id, 'PAYROLL_RUN', 'on_post', 'document', 'PAYROLL_APPROVED',
         'doc.approved_by != null && doc.approved_at != null',
         'error', TRUE,
         'Payroll run must be approved before posting',
         TRUE, 'Ensures payroll approval workflow', v_user_id, v_user_id),

        -- Warning: Large sale amount
        (v_org_id, 'POS_SALE', 'on_post', 'document', 'LARGE_SALE_AMOUNT',
         'doc.total_amount <= 10000.00',
         'warning', FALSE,
         'Sale amount ({doc.total_amount}) is unusually large',
         TRUE, 'Flags large transactions for review', v_user_id, v_user_id),

        -- Info: First-time posting
        (v_org_id, NULL, 'on_post', 'document', 'FIRST_POSTING',
         'doc.posted_to_accounting_at == null',
         'info', FALSE,
         'This is the first time this document is being posted',
         TRUE, 'Informational message for first posting', v_user_id, v_user_id);

    RAISE NOTICE 'Created % posting validation rules', (SELECT COUNT(*) FROM posting_validation_rules);

    -- ========================================================================
    -- SUMMARY
    -- ========================================================================

    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Posting Engine Seed Data Summary:';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Posting Concepts: %', (SELECT COUNT(*) FROM posting_concepts);
    RAISE NOTICE 'Posting Account Mappings: %', (SELECT COUNT(*) FROM posting_account_mappings WHERE organization_id = v_org_id);
    RAISE NOTICE 'Posting Profiles: 1 (DEFAULT_PROFILE)';
    RAISE NOTICE 'Document Types: %', (SELECT COUNT(*) FROM posting_document_types);
    RAISE NOTICE 'Profile Documents: %', (SELECT COUNT(*) FROM posting_profile_documents WHERE posting_profile_id = v_profile_id);
    RAISE NOTICE 'Posting Rules: %', (SELECT COUNT(*) FROM posting_rules WHERE posting_profile_document_id IN (
        SELECT id FROM posting_profile_documents WHERE posting_profile_id = v_profile_id
    ));
    RAISE NOTICE 'Posting Rule Lines: %', (SELECT COUNT(*) FROM posting_rule_lines WHERE posting_rule_id IN (
        SELECT id FROM posting_rules WHERE posting_profile_document_id IN (
            SELECT id FROM posting_profile_documents WHERE posting_profile_id = v_profile_id
        )
    ));
    RAISE NOTICE 'Validation Rules: %', (SELECT COUNT(*) FROM posting_validation_rules);
    RAISE NOTICE '============================================';
    RAISE NOTICE '';
    RAISE NOTICE 'Example Posting Rules Created:';
    RAISE NOTICE ' - POS_SALE_CASH: Cash sale → Debit CASH, Credit REVENUE+TAX';
    RAISE NOTICE ' - POS_SALE_CREDIT: Credit sale → Debit AR, Credit REVENUE+TAX';
    RAISE NOTICE ' - PAYROLL_ACCRUAL: Payroll run → Debit PAYROLL_EXPENSE, Credit PAYROLL_LIABILITY';
    RAISE NOTICE ' - VENDOR_BILL_POST: Vendor bill → Debit INVENTORY+TAX, Credit AP';
    RAISE NOTICE ' - CHECK_ISSUE: Check payment → Debit AP, Credit BANK';
    RAISE NOTICE '============================================';

    -- Test concept resolution
    RAISE NOTICE '';
    RAISE NOTICE 'Testing concept label resolution:';
    RAISE NOTICE 'AR label: %', get_concept_label(v_org_id, 'AR');
    RAISE NOTICE 'REVENUE label: %', get_concept_label(v_org_id, 'REVENUE');
    RAISE NOTICE 'PAYROLL_EXPENSE label: %', get_concept_label(v_org_id, 'PAYROLL_EXPENSE');

END $$;

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Seed Data 008 completed successfully!';
    RAISE NOTICE 'Posting Engine configured and ready.';
    RAISE NOTICE '============================================';
END $$;

-- =====================================================
-- Accounting Seed Data: Accounts Payable & Receivable
-- Description: Realistic vendor bills, customer invoices, and payments
-- Includes: 48 vendor bills, 60 customer invoices, full payment cycles
-- =====================================================

\echo 'Loading AP/AR data (vendor bills, customer invoices, payments)...';

-- =====================================================
-- SECTION 1: Vendor Bills (from inventory purchases)
-- =====================================================

DO $$
DECLARE
    v_org_id UUID;
    v_fiscal_year_id UUID;
    v_bill_id UUID;
    v_supplier_ids UUID[];
    v_supplier_id UUID;
    v_bill_date DATE;
    v_due_date DATE;
    v_bill_number VARCHAR(50);
    v_bill_amount NUMERIC(20,4);
    v_je_id UUID;
    v_bill_count INTEGER := 0;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_fiscal_year_id FROM fiscal_years WHERE organization_id = v_org_id AND fiscal_year = '2024' LIMIT 1;

    -- Get supplier IDs
    SELECT ARRAY_AGG(id) INTO v_supplier_ids FROM suppliers WHERE organization_id = v_org_id LIMIT 10;

    -- Create 48 vendor bills (4 per month for 12 months)
    FOR month_num IN 1..12 LOOP
        FOR bill_in_month IN 1..4 LOOP
            v_bill_count := v_bill_count + 1;
            v_bill_date := ('2024-' || LPAD(month_num::TEXT, 2, '0') || '-' || LPAD((bill_in_month * 7)::TEXT, 2, '0'))::DATE;
            v_due_date := v_bill_date + INTERVAL '30 days';
            v_bill_number := 'VB-2024-' || LPAD(v_bill_count::TEXT, 4, '0');
            v_bill_amount := 25000.00 + (RANDOM() * 15000);

            -- Rotate through suppliers
            v_supplier_id := v_supplier_ids[(v_bill_count % ARRAY_LENGTH(v_supplier_ids, 1)) + 1];

            -- Get the journal entry for this purchase (if exists)
            SELECT id INTO v_je_id
            FROM journal_entries
            WHERE organization_id = v_org_id
                AND entry_date = v_bill_date
                AND reference_type = 'PURCHASE'
            LIMIT 1;

            v_bill_id := gen_random_uuid();

            -- Create vendor bill
            INSERT INTO vendor_bills (
                id, organization_id, bill_number, supplier_id,
                bill_date, due_date, payment_terms,
                subtotal_amount, tax_amount, total_amount,
                paid_amount, balance_due, status,
                journal_entry_id, notes, created_by
            )
            SELECT
                v_bill_id,
                v_org_id,
                v_bill_number,
                v_supplier_id,
                v_bill_date,
                v_due_date,
                'Net 30',
                v_bill_amount,
                v_bill_amount * 0.08, -- 8% tax
                v_bill_amount * 1.08,
                0,
                v_bill_amount * 1.08,
                CASE
                    WHEN v_bill_date < '2024-11-01'::DATE THEN 'paid'
                    WHEN v_bill_date < '2024-11-15'::DATE THEN 'partial'
                    ELSE 'unpaid'
                END,
                v_je_id,
                'Inventory purchase for resale',
                u.id
            FROM users u
            WHERE u.email = 'admin@demoretail.com'
            LIMIT 1;

            -- Create bill line items (3 items per bill)
            INSERT INTO vendor_bill_items (
                id, vendor_bill_id, line_number, description,
                quantity, unit_price, subtotal_amount, tax_amount, total_amount
            )
            VALUES
                (gen_random_uuid(), v_bill_id, 1, 'Product Category A - Bulk Purchase',
                 100, v_bill_amount * 0.40 / 100, v_bill_amount * 0.40, v_bill_amount * 0.40 * 0.08, v_bill_amount * 0.40 * 1.08),
                (gen_random_uuid(), v_bill_id, 2, 'Product Category B - Bulk Purchase',
                 75, v_bill_amount * 0.35 / 75, v_bill_amount * 0.35, v_bill_amount * 0.35 * 0.08, v_bill_amount * 0.35 * 1.08),
                (gen_random_uuid(), v_bill_id, 3, 'Product Category C - Bulk Purchase',
                 50, v_bill_amount * 0.25 / 50, v_bill_amount * 0.25, v_bill_amount * 0.25 * 0.08, v_bill_amount * 0.25 * 1.08);

        END LOOP;
    END LOOP;

    RAISE NOTICE 'Created 48 vendor bills';
END $$;

-- =====================================================
-- SECTION 2: Vendor Payments (linked to bills)
-- =====================================================

DO $$
DECLARE
    v_org_id UUID;
    v_payment_id UUID;
    v_bill_record RECORD;
    v_payment_date DATE;
    v_payment_number VARCHAR(50);
    v_payment_amount NUMERIC(20,4);
    v_bank_account_id UUID;
    v_je_id UUID;
    v_payment_count INTEGER := 0;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Get bank account
    SELECT id INTO v_bank_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id AND account_code = '1020'
    LIMIT 1;

    -- Create payments for bills that are paid or partially paid
    FOR v_bill_record IN (
        SELECT id, bill_number, supplier_id, total_amount, bill_date, status
        FROM vendor_bills
        WHERE organization_id = v_org_id
            AND status IN ('paid', 'partial')
        ORDER BY bill_date
    ) LOOP
        v_payment_count := v_payment_count + 1;
        v_payment_date := v_bill_record.bill_date + INTERVAL '25 days'; -- Paid 5 days before due
        v_payment_number := 'VP-2024-' || LPAD(v_payment_count::TEXT, 4, '0');

        IF v_bill_record.status = 'paid' THEN
            v_payment_amount := v_bill_record.total_amount;
        ELSE
            v_payment_amount := v_bill_record.total_amount * 0.60; -- Partial payment
        END IF;

        -- Get the journal entry for this payment (if exists)
        SELECT id INTO v_je_id
        FROM journal_entries
        WHERE organization_id = v_org_id
            AND entry_date = v_payment_date
            AND reference_type = 'PAYMENT'
        LIMIT 1;

        v_payment_id := gen_random_uuid();

        -- Create vendor payment
        INSERT INTO vendor_payments (
            id, organization_id, payment_number, supplier_id,
            payment_date, payment_method, payment_amount,
            bank_account_id, journal_entry_id, notes, created_by
        )
        SELECT
            v_payment_id,
            v_org_id,
            v_payment_number,
            v_bill_record.supplier_id,
            v_payment_date,
            CASE (v_payment_count % 3)
                WHEN 0 THEN 'check'
                WHEN 1 THEN 'bank_transfer'
                ELSE 'eft'
            END,
            v_payment_amount,
            v_bank_account_id,
            v_je_id,
            'Payment for bill ' || v_bill_record.bill_number,
            u.id
        FROM users u
        WHERE u.email = 'admin@demoretail.com'
        LIMIT 1;

        -- Link payment to bill
        INSERT INTO vendor_payment_applications (
            id, vendor_payment_id, vendor_bill_id,
            applied_amount, created_by
        )
        SELECT
            gen_random_uuid(),
            v_payment_id,
            v_bill_record.id,
            v_payment_amount,
            u.id
        FROM users u
        WHERE u.email = 'admin@demoretail.com'
        LIMIT 1;

        -- Update bill status
        UPDATE vendor_bills
        SET paid_amount = v_payment_amount,
            balance_due = total_amount - v_payment_amount,
            status = CASE
                WHEN v_payment_amount >= total_amount THEN 'paid'
                WHEN v_payment_amount > 0 THEN 'partial'
                ELSE 'unpaid'
            END
        WHERE id = v_bill_record.id;

    END LOOP;

    RAISE NOTICE 'Created % vendor payments', v_payment_count;
END $$;

-- =====================================================
-- SECTION 3: Customer Invoices (from credit sales)
-- =====================================================

DO $$
DECLARE
    v_org_id UUID;
    v_fiscal_year_id UUID;
    v_invoice_id UUID;
    v_customer_ids UUID[];
    v_customer_id UUID;
    v_invoice_date DATE;
    v_due_date DATE;
    v_invoice_number VARCHAR(50);
    v_invoice_amount NUMERIC(20,4);
    v_je_id UUID;
    v_invoice_count INTEGER := 0;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_fiscal_year_id FROM fiscal_years WHERE organization_id = v_org_id AND fiscal_year = '2024' LIMIT 1;

    -- Get customer IDs
    SELECT ARRAY_AGG(id) INTO v_customer_ids FROM customers WHERE organization_id = v_org_id LIMIT 15;

    -- Create 60 customer invoices (5 per month for 12 months)
    FOR month_num IN 1..12 LOOP
        FOR inv_in_month IN 1..5 LOOP
            v_invoice_count := v_invoice_count + 1;
            v_invoice_date := ('2024-' || LPAD(month_num::TEXT, 2, '0') || '-' || LPAD((inv_in_month * 5 + 2)::TEXT, 2, '0'))::DATE;
            v_due_date := v_invoice_date + INTERVAL '15 days'; -- Net 15
            v_invoice_number := 'INV-2024-' || LPAD(v_invoice_count::TEXT, 5, '0');
            v_invoice_amount := 8000.00 + (RANDOM() * 12000);

            -- Rotate through customers
            v_customer_id := v_customer_ids[(v_invoice_count % ARRAY_LENGTH(v_customer_ids, 1)) + 1];

            -- Get the journal entry for sales (if exists)
            SELECT id INTO v_je_id
            FROM journal_entries
            WHERE organization_id = v_org_id
                AND EXTRACT(MONTH FROM entry_date) = month_num
                AND reference_type = 'SALES'
            LIMIT 1;

            v_invoice_id := gen_random_uuid();

            -- Create customer invoice
            INSERT INTO customer_invoices (
                id, organization_id, invoice_number, customer_id,
                invoice_date, due_date, payment_terms,
                subtotal_amount, tax_amount, discount_amount, total_amount,
                paid_amount, balance_due, status,
                journal_entry_id, notes, created_by
            )
            SELECT
                v_invoice_id,
                v_org_id,
                v_invoice_number,
                v_customer_id,
                v_invoice_date,
                v_due_date,
                'Net 15',
                v_invoice_amount,
                v_invoice_amount * 0.10, -- 10% tax
                v_invoice_amount * 0.02, -- 2% discount
                v_invoice_amount * 1.10 - v_invoice_amount * 0.02,
                0,
                v_invoice_amount * 1.10 - v_invoice_amount * 0.02,
                CASE
                    WHEN v_invoice_date < '2024-10-15'::DATE THEN 'paid'
                    WHEN v_invoice_date < '2024-11-01'::DATE THEN 'partial'
                    ELSE 'unpaid'
                END,
                v_je_id,
                'Sales invoice - Net 15 payment terms',
                u.id
            FROM users u
            WHERE u.email = 'admin@demoretail.com'
            LIMIT 1;

            -- Create invoice line items (4 items per invoice)
            INSERT INTO customer_invoice_items (
                id, customer_invoice_id, line_number, description,
                quantity, unit_price, discount_percentage,
                subtotal_amount, tax_amount, total_amount
            )
            VALUES
                (gen_random_uuid(), v_invoice_id, 1, 'Electronics - Premium Product Line',
                 15, v_invoice_amount * 0.35 / 15, 0, v_invoice_amount * 0.35, v_invoice_amount * 0.35 * 0.10, v_invoice_amount * 0.35 * 1.10),
                (gen_random_uuid(), v_invoice_id, 2, 'Clothing - Seasonal Collection',
                 20, v_invoice_amount * 0.30 / 20, 2, v_invoice_amount * 0.30, v_invoice_amount * 0.30 * 0.10, v_invoice_amount * 0.30 * 1.10),
                (gen_random_uuid(), v_invoice_id, 3, 'Home & Garden - Featured Items',
                 25, v_invoice_amount * 0.25 / 25, 0, v_invoice_amount * 0.25, v_invoice_amount * 0.25 * 0.10, v_invoice_amount * 0.25 * 1.10),
                (gen_random_uuid(), v_invoice_id, 4, 'Accessories & Supplies',
                 30, v_invoice_amount * 0.10 / 30, 5, v_invoice_amount * 0.10, v_invoice_amount * 0.10 * 0.10, v_invoice_amount * 0.10 * 1.10);

        END LOOP;
    END LOOP;

    RAISE NOTICE 'Created 60 customer invoices';
END $$;

-- =====================================================
-- SECTION 4: Customer Payments (linked to invoices)
-- =====================================================

DO $$
DECLARE
    v_org_id UUID;
    v_payment_id UUID;
    v_invoice_record RECORD;
    v_payment_date DATE;
    v_payment_number VARCHAR(50);
    v_payment_amount NUMERIC(20,4);
    v_bank_account_id UUID;
    v_je_id UUID;
    v_payment_count INTEGER := 0;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Get bank account
    SELECT id INTO v_bank_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id AND account_code = '1020'
    LIMIT 1;

    -- Create payments for invoices that are paid or partially paid
    FOR v_invoice_record IN (
        SELECT id, invoice_number, customer_id, total_amount, invoice_date, status
        FROM customer_invoices
        WHERE organization_id = v_org_id
            AND status IN ('paid', 'partial')
        ORDER BY invoice_date
    ) LOOP
        v_payment_count := v_payment_count + 1;
        v_payment_date := v_invoice_record.invoice_date + INTERVAL '12 days'; -- Paid 3 days before due
        v_payment_number := 'CP-2024-' || LPAD(v_payment_count::TEXT, 5, '0');

        IF v_invoice_record.status = 'paid' THEN
            v_payment_amount := v_invoice_record.total_amount;
        ELSE
            v_payment_amount := v_invoice_record.total_amount * 0.50; -- Partial payment
        END IF;

        -- Get the journal entry for this receipt (if exists)
        SELECT id INTO v_je_id
        FROM journal_entries
        WHERE organization_id = v_org_id
            AND entry_date = v_payment_date
            AND reference_type = 'RECEIPT'
        LIMIT 1;

        v_payment_id := gen_random_uuid();

        -- Create customer payment
        INSERT INTO customer_payments (
            id, organization_id, payment_number, customer_id,
            payment_date, payment_method, payment_amount,
            bank_account_id, journal_entry_id, notes, created_by
        )
        SELECT
            v_payment_id,
            v_org_id,
            v_payment_number,
            v_invoice_record.customer_id,
            v_payment_date,
            CASE (v_payment_count % 4)
                WHEN 0 THEN 'cash'
                WHEN 1 THEN 'credit_card'
                WHEN 2 THEN 'bank_transfer'
                ELSE 'check'
            END,
            v_payment_amount,
            v_bank_account_id,
            v_je_id,
            'Payment for invoice ' || v_invoice_record.invoice_number,
            u.id
        FROM users u
        WHERE u.email = 'admin@demoretail.com'
        LIMIT 1;

        -- Link payment to invoice
        INSERT INTO customer_payment_applications (
            id, customer_payment_id, customer_invoice_id,
            applied_amount, created_by
        )
        SELECT
            gen_random_uuid(),
            v_payment_id,
            v_invoice_record.id,
            v_payment_amount,
            u.id
        FROM users u
        WHERE u.email = 'admin@demoretail.com'
        LIMIT 1;

        -- Update invoice status
        UPDATE customer_invoices
        SET paid_amount = v_payment_amount,
            balance_due = total_amount - v_payment_amount,
            status = CASE
                WHEN v_payment_amount >= total_amount THEN 'paid'
                WHEN v_payment_amount > 0 THEN 'partial'
                ELSE 'unpaid'
            END
        WHERE id = v_invoice_record.id;

    END LOOP;

    RAISE NOTICE 'Created % customer payments', v_payment_count;
END $$;

\echo '';
\echo '==========================================';
\echo 'AP/AR Data Summary:';
\echo '==========================================';
\echo 'Vendor Bills:               48 bills';
\echo 'Vendor Payments:            ~40 payments';
\echo 'Customer Invoices:          60 invoices';
\echo 'Customer Payments:          ~50 payments';
\echo '==========================================';
\echo '';
\echo 'AP/AR Status:';
\echo '  ✓ Bills linked to journal entries';
\echo '  ✓ Payments applied to bills/invoices';
\echo '  ✓ Aging data available (paid, partial, unpaid)';
\echo '  ✓ Ready for AP/AR aging reports';
\echo '';
\echo 'AP/AR data loaded successfully!';

-- =====================================================
-- End of seed file 004
-- =====================================================

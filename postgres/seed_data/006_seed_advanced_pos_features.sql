-- =====================================================
-- POS Seed Data: Advanced Features
-- Description: Comprehensive test data for cash handling, pricing,
--              gift cards, returns, procurement, monitoring, multi-channel
-- =====================================================

\echo 'Loading advanced POS feature seed data...';

-- =====================================================
-- SECTION 1: MONEY HANDLING
-- =====================================================

\echo 'Seeding POS sessions and cash management...';

DO $$
DECLARE
    v_org_id UUID;
    v_location_id UUID;
    v_device_id UUID;
    v_user_id UUID;
    v_shift_id UUID;
    v_session_id UUID;
    v_drawer_id UUID;
    v_sale_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_location_id FROM locations WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_device_id FROM devices WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_user_id FROM users WHERE email = 'admin@demoretail.com' LIMIT 1;
    SELECT id INTO v_shift_id FROM shifts WHERE organization_id = v_org_id LIMIT 1;

    -- Create cash drawer
    INSERT INTO cash_drawers (id, organization_id, drawer_code, drawer_name, location_id, device_id, is_active, created_by)
    VALUES (gen_random_uuid(), v_org_id, 'DRAWER-01', 'Main Register Drawer', v_location_id, v_device_id, true, v_user_id)
    RETURNING id INTO v_drawer_id;

    -- POS Session 1: Completed session
    INSERT INTO pos_sessions (
        id, organization_id, session_number, session_name, device_id, location_id, user_id, shift_id,
        opened_at, closed_at,
        opening_cash, opening_card, opening_other,
        expected_cash, expected_card, expected_other,
        counted_cash, counted_card, counted_other,
        difference_cash, difference_card, difference_other,
        status, z_report_number, created_by
    )
    VALUES (
        gen_random_uuid(), v_org_id, 'SES-2024-001', 'Morning Shift - Nov 1', v_device_id, v_location_id, v_user_id, v_shift_id,
        '2024-11-01 08:00:00', '2024-11-01 16:30:00',
        200.00, 0, 0,
        3542.50, 5280.00, 125.00,
        3545.00, 5280.00, 125.00,
        2.50, 0, 0,
        'reconciled', 'Z-2024-11-01-001', v_user_id
    )
    RETURNING id INTO v_session_id;

    -- Link drawer to session
    INSERT INTO cash_drawer_sessions (organization_id, cash_drawer_id, pos_session_id, opening_amount, closing_amount)
    VALUES (v_org_id, v_drawer_id, v_session_id, 200.00, 3545.00);

    -- Cash movements during session
    INSERT INTO cash_movements (organization_id, pos_session_id, cash_drawer_id, movement_type, amount, reason_code, reason_description, user_id)
    VALUES
        (v_org_id, v_session_id, v_drawer_id, 'float_add', 100.00, 'CHANGE', 'Added change for busy morning', v_user_id),
        (v_org_id, v_session_id, v_drawer_id, 'pay_out', -25.50, 'PETTY_CASH', 'Office supplies reimbursement', v_user_id),
        (v_org_id, v_session_id, v_drawer_id, 'drop_to_safe', -500.00, 'SAFE_DROP', 'Cash drop to safe at noon', v_user_id),
        (v_org_id, v_session_id, v_drawer_id, 'pay_in', 50.00, 'FOUND', 'Found cash from yesterday', v_user_id);

    -- POS Session 2: Currently open session
    INSERT INTO pos_sessions (
        id, organization_id, session_number, session_name, device_id, location_id, user_id,
        opened_at, opening_cash, opening_card, opening_other,
        status, created_by
    )
    VALUES (
        gen_random_uuid(), v_org_id, 'SES-2024-002', 'Afternoon Shift - Nov 10', v_device_id, v_location_id, v_user_id,
        '2024-11-10 13:00:00', 250.00, 0, 0,
        'open', v_user_id
    )
    RETURNING id INTO v_session_id;

    -- Update some sales to reference sessions
    UPDATE sales SET pos_session_id = v_session_id WHERE id IN (
        SELECT id FROM sales WHERE organization_id = v_org_id ORDER BY sale_date DESC LIMIT 5
    );

    RAISE NOTICE 'Created 2 POS sessions with cash drawer and movements';
END $$;

-- =====================================================
-- SECTION 2: ADVANCED PRICING
-- =====================================================

\echo 'Seeding price lists and product components...';

DO $$
DECLARE
    v_org_id UUID;
    v_retail_pl UUID;
    v_wholesale_pl UUID;
    v_vip_pl UUID;
    v_product_id UUID;
    v_variant_id UUID;
    v_category_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Create price lists
    INSERT INTO price_lists (id, organization_id, price_list_code, price_list_name, price_list_type, effective_from, priority, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'RETAIL', 'Retail Pricing', 'retail', '2024-01-01', 10, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_retail_pl;

    INSERT INTO price_lists (id, organization_id, price_list_code, price_list_name, price_list_type, base_price_adjustment_type, base_price_adjustment_value, effective_from, priority, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'WHOLESALE', 'Wholesale Pricing', 'wholesale', 'percentage', -15.00, '2024-01-01', 20, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_wholesale_pl;

    INSERT INTO price_lists (id, organization_id, price_list_code, price_list_name, price_list_type, base_price_adjustment_type, base_price_adjustment_value, effective_from, priority, is_active, created_by)
    SELECT gen_random_uuid(), v_org_id, 'VIP', 'VIP Customer Pricing', 'vip', 'percentage', -10.00, '2024-01-01', 15, true, u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1
    RETURNING id INTO v_vip_pl;

    INSERT INTO price_lists (organization_id, price_list_code, price_list_name, price_list_type, effective_from, effective_to, priority, is_active, description, created_by)
    SELECT v_org_id, 'BLACKFRIDAY24', 'Black Friday 2024', 'seasonal', '2024-11-29', '2024-12-02', 5, true, 'Black Friday special pricing', u.id FROM users u WHERE u.email = 'admin@demoretail.com' LIMIT 1;

    -- Price list items for specific products
    SELECT id INTO v_product_id FROM products WHERE organization_id = v_org_id AND product_name LIKE '%Laptop%' LIMIT 1;
    SELECT id INTO v_variant_id FROM product_variants WHERE product_id IN (SELECT id FROM products WHERE organization_id = v_org_id LIMIT 1) LIMIT 1;
    SELECT id INTO v_category_id FROM categories WHERE organization_id = v_org_id AND category_name = 'Electronics' LIMIT 1;

    -- Wholesale pricing on electronics category (15% off)
    INSERT INTO price_list_items (price_list_id, category_id, discount_percentage)
    VALUES (v_wholesale_pl, v_category_id, 15.00);

    -- VIP pricing on specific products
    INSERT INTO price_list_items (price_list_id, product_id, discount_percentage, min_quantity)
    SELECT v_vip_pl, id, 12.00, 1 FROM products WHERE organization_id = v_org_id AND product_name LIKE '%Premium%' LIMIT 3;

    -- Black Friday pricing (deep discounts)
    INSERT INTO price_list_items (price_list_id, product_id, discount_percentage)
    SELECT (SELECT id FROM price_lists WHERE price_list_code = 'BLACKFRIDAY24'), id, 25.00
    FROM products WHERE organization_id = v_org_id AND retail_price > 100 LIMIT 10;

    -- Product components for composite products
    SELECT id INTO v_product_id FROM products WHERE organization_id = v_org_id AND is_composite = true LIMIT 1;

    IF v_product_id IS NOT NULL THEN
        -- Create a bundle: "Complete Home Office Setup"
        INSERT INTO product_components (organization_id, parent_product_id, component_product_id, quantity, inherit_price, display_order)
        SELECT v_org_id, v_product_id, id, 1, true, ROW_NUMBER() OVER (ORDER BY product_name)
        FROM products
        WHERE organization_id = v_org_id
          AND id != v_product_id
          AND product_name IN ('Desk', 'Office Chair', 'Monitor', 'Keyboard')
        LIMIT 4;
    END IF;

    -- Update some customers with price lists
    UPDATE customers SET price_list_id = v_wholesale_pl WHERE customer_type = 'business' AND organization_id = v_org_id;
    UPDATE customers SET price_list_id = v_vip_pl WHERE loyalty_tier = 'vip' AND organization_id = v_org_id;

    RAISE NOTICE 'Created 4 price lists with items and product components';
END $$;

-- =====================================================
-- SECTION 3: GIFT CARDS & STORE CREDIT
-- =====================================================

\echo 'Seeding gift cards and store credit...';

DO $$
DECLARE
    v_org_id UUID;
    v_user_id UUID;
    v_location_id UUID;
    v_customer_id UUID;
    v_gift_card_id UUID;
    v_credit_account_id UUID;
    v_sale_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_user_id FROM users WHERE email = 'admin@demoretail.com' LIMIT 1;
    SELECT id INTO v_location_id FROM locations WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_customer_id FROM customers WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_sale_id FROM sales WHERE organization_id = v_org_id LIMIT 1;

    -- Gift Card 1: Active with balance
    INSERT INTO gift_cards (id, organization_id, card_number, pin_code, customer_id, original_value, current_balance, issued_date, status, issued_by_user_id, issued_location_id, created_by)
    VALUES (gen_random_uuid(), v_org_id, 'GC-2024-1000001', '1234', v_customer_id, 100.00, 75.50, '2024-01-15', 'active', v_user_id, v_location_id, v_user_id)
    RETURNING id INTO v_gift_card_id;

    -- Gift card transactions
    INSERT INTO gift_card_transactions (organization_id, gift_card_id, transaction_type, amount, balance_after, sale_id, user_id, location_id, notes)
    VALUES
        (v_org_id, v_gift_card_id, 'issue', 100.00, 100.00, NULL, v_user_id, v_location_id, 'Initial issue'),
        (v_org_id, v_gift_card_id, 'redemption', -24.50, 75.50, v_sale_id, v_user_id, v_location_id, 'Used at checkout');

    -- Gift Card 2: Fully redeemed
    INSERT INTO gift_cards (organization_id, card_number, customer_id, original_value, current_balance, issued_date, status, issued_by_user_id, issued_location_id, created_by)
    VALUES (v_org_id, 'GC-2024-1000002', v_customer_id, 50.00, 0, '2024-02-10', 'fully_redeemed', v_user_id, v_location_id, v_user_id)
    RETURNING id INTO v_gift_card_id;

    INSERT INTO gift_card_transactions (organization_id, gift_card_id, transaction_type, amount, balance_after, user_id, location_id)
    VALUES
        (v_org_id, v_gift_card_id, 'issue', 50.00, 50.00, v_user_id, v_location_id),
        (v_org_id, v_gift_card_id, 'redemption', -30.00, 20.00, v_user_id, v_location_id),
        (v_org_id, v_gift_card_id, 'redemption', -20.00, 0, v_user_id, v_location_id);

    -- Gift Card 3: High value for corporate client
    INSERT INTO gift_cards (organization_id, card_number, pin_code, original_value, current_balance, issued_date, expiry_date, status, issued_by_user_id, issued_location_id, notes, created_by)
    VALUES (v_org_id, 'GC-2024-1000003', '5678', 500.00, 500.00, '2024-11-01', '2025-11-01', 'active', v_user_id, v_location_id, 'Corporate gift - Holiday bonus', v_user_id);

    -- Store credit accounts for customers
    FOR v_customer_id IN (SELECT id FROM customers WHERE organization_id = v_org_id LIMIT 5) LOOP
        INSERT INTO customer_store_credit_accounts (organization_id, customer_id, current_balance, credit_limit, is_active)
        VALUES (v_org_id, v_customer_id, 25.00 + (RANDOM() * 75), 500.00, true)
        RETURNING id INTO v_credit_account_id;

        -- Store credit transaction (return refund to credit)
        INSERT INTO store_credit_transactions (organization_id, store_credit_account_id, transaction_type, amount, balance_after, user_id, location_id, notes)
        VALUES (v_org_id, v_credit_account_id, 'issue', 25.00, 25.00, v_user_id, v_location_id, 'Return refund issued as store credit');
    END LOOP;

    RAISE NOTICE 'Created 3 gift cards and 5 store credit accounts with transactions';
END $$;

-- =====================================================
-- SECTION 4: RETURNS & AFTER-SALES
-- =====================================================

\echo 'Seeding returns...';

DO $$
DECLARE
    v_org_id UUID;
    v_user_id UUID;
    v_location_id UUID;
    v_customer_id UUID;
    v_sale_id UUID;
    v_sale_item_id UUID;
    v_product_id UUID;
    v_return_id UUID;
    v_reason_defective UUID;
    v_reason_wrong_size UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_user_id FROM users WHERE email = 'admin@demoretail.com' LIMIT 1;
    SELECT id INTO v_location_id FROM locations WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_customer_id FROM customers WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_sale_id FROM sales WHERE organization_id = v_org_id AND customer_id = v_customer_id LIMIT 1;
    SELECT id INTO v_sale_item_id FROM sale_items WHERE sale_id = v_sale_id LIMIT 1;
    SELECT product_id INTO v_product_id FROM sale_items WHERE id = v_sale_item_id;

    SELECT id INTO v_reason_defective FROM return_reasons WHERE reason_code = 'DEFECTIVE' LIMIT 1;
    SELECT id INTO v_reason_wrong_size FROM return_reasons WHERE reason_code = 'WRONG_SIZE' LIMIT 1;

    -- Return 1: Defective item - full refund
    INSERT INTO sale_returns (
        id, organization_id, return_number, original_sale_id, customer_id, location_id, user_id,
        return_date, total_amount, refund_amount, restocking_fee, refund_method, status,
        approved_by, approved_at, notes, created_by
    )
    VALUES (
        gen_random_uuid(), v_org_id, 'RET-2024-001', v_sale_id, v_customer_id, v_location_id, v_user_id,
        '2024-11-05', 89.99, 89.99, 0, 'original_payment', 'completed',
        v_user_id, '2024-11-05 14:30:00', 'Defective product - full refund processed', v_user_id
    )
    RETURNING id INTO v_return_id;

    INSERT INTO sale_return_items (
        organization_id, sale_return_id, original_sale_item_id, product_id,
        quantity, unit_price, subtotal, tax_amount, total_amount,
        return_reason_id, return_reason_notes, item_condition, is_restockable
    )
    VALUES (
        v_org_id, v_return_id, v_sale_item_id, v_product_id,
        1, 89.99, 89.99, 0, 89.99,
        v_reason_defective, 'Screen not turning on', 'defective', false
    );

    -- Return 2: Wrong size - exchange
    INSERT INTO sale_returns (
        organization_id, return_number, customer_id, location_id, user_id,
        return_date, total_amount, refund_amount, refund_method, status, notes, created_by
    )
    VALUES (
        v_org_id, 'RET-2024-002', v_customer_id, v_location_id, v_user_id,
        '2024-11-07', 49.99, 0, 'exchange', 'completed', 'Exchanged for different size', v_user_id
    )
    RETURNING id INTO v_return_id;

    INSERT INTO sale_return_items (
        organization_id, sale_return_id, product_id,
        quantity, unit_price, subtotal, tax_amount, total_amount,
        return_reason_id, return_reason_notes, item_condition, is_restockable
    )
    SELECT
        v_org_id, v_return_id, id,
        1, 49.99, 49.99, 0, 49.99,
        v_reason_wrong_size, 'Customer wanted large instead of medium', 'resellable', true
    FROM products WHERE organization_id = v_org_id AND product_name LIKE '%Shirt%' LIMIT 1;

    -- Return 3: Pending approval (high value)
    INSERT INTO sale_returns (
        organization_id, return_number, customer_id, location_id, user_id,
        return_date, total_amount, refund_amount, restocking_fee, refund_method, status, notes, created_by
    )
    VALUES (
        v_org_id, 'RET-2024-003', v_customer_id, v_location_id, v_user_id,
        '2024-11-09', 1299.99, 1299.99, 0, 'original_payment', 'pending',
        'High value return - pending manager approval', v_user_id
    )
    RETURNING id INTO v_return_id;

    INSERT INTO sale_return_items (
        organization_id, sale_return_id, product_id,
        quantity, unit_price, subtotal, tax_amount, total_amount,
        return_reason_id, item_condition, is_restockable
    )
    SELECT
        v_org_id, v_return_id, id,
        1, 1299.99, 1299.99, 0, 1299.99,
        (SELECT id FROM return_reasons WHERE reason_code = 'CHANGED_MIND'), 'resellable', true
    FROM products WHERE organization_id = v_org_id AND retail_price > 1000 LIMIT 1;

    RAISE NOTICE 'Created 3 sale returns with items';
END $$;

-- =====================================================
-- SECTION 5: PROCUREMENT
-- =====================================================

\echo 'Seeding purchase orders...';

DO $$
DECLARE
    v_org_id UUID;
    v_supplier_id UUID;
    v_location_id UUID;
    v_user_id UUID;
    v_po_id UUID;
    v_po_item_id UUID;
    v_receipt_id UUID;
    v_product_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_supplier_id FROM suppliers WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_location_id FROM locations WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_user_id FROM users WHERE email = 'admin@demoretail.com' LIMIT 1;

    -- Purchase Order 1: Fully received
    INSERT INTO purchase_orders (
        id, organization_id, po_number, supplier_id, location_id,
        order_date, expected_delivery_date, actual_delivery_date,
        subtotal_amount, tax_amount, shipping_amount, total_amount,
        payment_terms, payment_due_date, status,
        approved_by, approved_at, notes, created_by
    )
    VALUES (
        gen_random_uuid(), v_org_id, 'PO-2024-001', v_supplier_id, v_location_id,
        '2024-10-15', '2024-10-22', '2024-10-21',
        5000.00, 400.00, 150.00, 5550.00,
        'Net 30', '2024-11-20', 'received',
        v_user_id, '2024-10-15 10:00:00', 'Monthly stock replenishment', v_user_id
    )
    RETURNING id INTO v_po_id;

    -- PO Items
    FOR v_product_id IN (SELECT id FROM products WHERE organization_id = v_org_id LIMIT 5) LOOP
        INSERT INTO purchase_order_items (
            organization_id, purchase_order_id, line_number, product_id,
            quantity_ordered, quantity_received, unit_cost, subtotal, tax_amount, total_amount,
            expected_delivery_date, notes
        )
        VALUES (
            v_org_id, v_po_id, (SELECT COALESCE(MAX(line_number), 0) + 1 FROM purchase_order_items WHERE purchase_order_id = v_po_id),
            v_product_id, 20, 20, 50.00, 1000.00, 80.00, 1080.00, '2024-10-22', NULL
        )
        RETURNING id INTO v_po_item_id;
    END LOOP;

    -- Goods Receipt for PO-2024-001
    INSERT INTO goods_receipts (
        id, organization_id, receipt_number, purchase_order_id, supplier_id, location_id,
        receipt_date, received_by, status, notes, created_by
    )
    VALUES (
        gen_random_uuid(), v_org_id, 'GR-2024-001', v_po_id, v_supplier_id, v_location_id,
        '2024-10-21', v_user_id, 'completed', 'All items received in good condition', v_user_id
    )
    RETURNING id INTO v_receipt_id;

    -- Goods Receipt Items
    FOR v_po_item_id IN (SELECT id FROM purchase_order_items WHERE purchase_order_id = v_po_id) LOOP
        INSERT INTO goods_receipt_items (
            organization_id, goods_receipt_id, purchase_order_item_id, product_id,
            quantity_received, quantity_accepted, quantity_rejected
        )
        SELECT v_org_id, v_receipt_id, v_po_item_id, product_id, 20, 20, 0
        FROM purchase_order_items WHERE id = v_po_item_id;
    END LOOP;

    -- Purchase Order 2: Partially received
    INSERT INTO purchase_orders (
        organization_id, po_number, supplier_id, location_id,
        order_date, expected_delivery_date,
        subtotal_amount, tax_amount, shipping_amount, total_amount,
        payment_terms, status, notes, created_by
    )
    VALUES (
        v_org_id, 'PO-2024-002', v_supplier_id, v_location_id,
        '2024-11-01', '2024-11-10',
        8000.00, 640.00, 200.00, 8840.00,
        'Net 30', 'partially_received', 'Holiday season inventory', v_user_id
    )
    RETURNING id INTO v_po_id;

    FOR v_product_id IN (SELECT id FROM products WHERE organization_id = v_org_id ORDER BY RANDOM() LIMIT 8) LOOP
        INSERT INTO purchase_order_items (
            organization_id, purchase_order_id, line_number, product_id,
            quantity_ordered, quantity_received, unit_cost, subtotal, tax_amount, total_amount,
            expected_delivery_date
        )
        VALUES (
            v_org_id, v_po_id, (SELECT COALESCE(MAX(line_number), 0) + 1 FROM purchase_order_items WHERE purchase_order_id = v_po_id),
            v_product_id, 30, 15, 66.67, 2000.00, 160.00, 2160.00, '2024-11-10'
        );
    END LOOP;

    -- Purchase Order 3: Pending (not yet sent)
    INSERT INTO purchase_orders (
        organization_id, po_number, supplier_id, location_id,
        order_date, expected_delivery_date,
        subtotal_amount, tax_amount, total_amount,
        status, notes, created_by
    )
    VALUES (
        v_org_id, 'PO-2024-003', v_supplier_id, v_location_id,
        '2024-11-10', '2024-11-20',
        3500.00, 280.00, 3780.00,
        'draft', 'Pending approval', v_user_id
    )
    RETURNING id INTO v_po_id;

    FOR v_product_id IN (SELECT id FROM products WHERE organization_id = v_org_id ORDER BY RANDOM() LIMIT 4) LOOP
        INSERT INTO purchase_order_items (
            organization_id, purchase_order_id, line_number, product_id,
            quantity_ordered, unit_cost, subtotal, tax_amount, total_amount
        )
        VALUES (
            v_org_id, v_po_id, (SELECT COALESCE(MAX(line_number), 0) + 1 FROM purchase_order_items WHERE purchase_order_id = v_po_id),
            v_product_id, 25, 35.00, 875.00, 70.00, 945.00
        );
    END LOOP;

    RAISE NOTICE 'Created 3 purchase orders with items and 1 goods receipt';
END $$;

-- =====================================================
-- SECTION 6: MONITORING & SYSTEM HEALTH
-- =====================================================

\echo 'Seeding system health and error logs...';

DO $$
DECLARE
    v_org_id UUID;
    v_device_id UUID;
    v_user_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_device_id FROM devices WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_user_id FROM users WHERE email = 'admin@demoretail.com' LIMIT 1;

    -- System health checks
    INSERT INTO system_health (organization_id, check_type, check_name, status, last_check_at, last_success_at, metric_value, metric_unit, threshold_warning, threshold_critical)
    VALUES
        (v_org_id, 'database', 'Database Connection', 'healthy', NOW() - INTERVAL '5 minutes', NOW() - INTERVAL '5 minutes', 2.5, 'ms', 50, 100),
        (v_org_id, 'database', 'Query Performance', 'healthy', NOW() - INTERVAL '5 minutes', NOW() - INTERVAL '5 minutes', 15.2, 'ms', 100, 500),
        (v_org_id, 'backup', 'Last Backup', 'healthy', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '2 hours', 2, 'hours', 24, 48),
        (v_org_id, 'storage', 'Disk Space Usage', 'healthy', NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '10 minutes', 45.2, 'percent', 80, 90),
        (v_org_id, 'api', 'Payment Gateway', 'healthy', NOW() - INTERVAL '1 minute', NOW() - INTERVAL '1 minute', 120, 'ms', 500, 1000),
        (NULL, 'service', 'Email Service', 'degraded', NOW() - INTERVAL '2 minutes', NOW() - INTERVAL '15 minutes', 850, 'ms', 500, 1000),
        (NULL, 'service', 'SMS Service', 'healthy', NOW(), NOW(), 200, 'ms', 1000, 2000);

    -- POS Error logs
    INSERT INTO pos_error_logs (organization_id, error_level, error_code, error_message, device_id, user_id, occurred_at, is_resolved)
    VALUES
        (v_org_id, 'warning', 'PRINT_FAIL', 'Receipt printer offline - using backup printer', v_device_id, v_user_id, NOW() - INTERVAL '2 hours', true),
        (v_org_id, 'error', 'PAYMENT_DECLINE', 'Card payment declined - insufficient funds', v_device_id, v_user_id, NOW() - INTERVAL '3 hours', true),
        (v_org_id, 'info', 'CASH_DRAWER_OPEN', 'Cash drawer manually opened outside of sale', v_device_id, v_user_id, NOW() - INTERVAL '1 hour', true),
        (v_org_id, 'warning', 'LOW_STOCK', 'Product below reorder point', NULL, NULL, NOW() - INTERVAL '30 minutes', false),
        (v_org_id, 'error', 'BARCODE_INVALID', 'Barcode not found in system', v_device_id, v_user_id, NOW() - INTERVAL '15 minutes', false);

    RAISE NOTICE 'Created 7 system health checks and 5 error logs';
END $$;

-- =====================================================
-- SECTION 7: MULTI-CHANNEL
-- =====================================================

\echo 'Seeding multi-channel data...';

DO $$
DECLARE
    v_org_id UUID;
    v_channel_web UUID;
    v_channel_mobile UUID;
    v_sale_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;

    -- Sales channels should already exist from migration
    SELECT id INTO v_channel_web FROM sales_channels WHERE organization_id = v_org_id AND channel_code = 'WEB' LIMIT 1;
    SELECT id INTO v_channel_mobile FROM sales_channels WHERE organization_id = v_org_id AND channel_code = 'MOBILE' LIMIT 1;

    -- External order mappings for web sales
    FOR v_sale_id IN (SELECT id FROM sales WHERE organization_id = v_org_id ORDER BY RANDOM() LIMIT 5) LOOP
        INSERT INTO external_order_mappings (
            organization_id, sale_id, sales_channel_id, external_order_id, external_order_number,
            sync_status, last_sync_at, external_data
        )
        VALUES (
            v_org_id, v_sale_id, v_channel_web,
            'WEB-' || LPAD(FLOOR(RANDOM() * 100000)::TEXT, 5, '0'),
            'WO-2024-' || LPAD(FLOOR(RANDOM() * 10000)::TEXT, 5, '0'),
            'synced', NOW() - INTERVAL '1 hour',
            jsonb_build_object(
                'customer_ip', '192.168.1.' || FLOOR(RANDOM() * 255)::TEXT,
                'user_agent', 'Mozilla/5.0',
                'referrer', 'https://google.com',
                'utm_source', 'google',
                'utm_medium', 'cpc'
            )
        );
    END LOOP;

    -- External order mappings for mobile app
    FOR v_sale_id IN (SELECT id FROM sales WHERE organization_id = v_org_id ORDER BY RANDOM() LIMIT 3) LOOP
        INSERT INTO external_order_mappings (
            organization_id, sale_id, sales_channel_id, external_order_id, external_order_number,
            sync_status, last_sync_at, external_data
        )
        VALUES (
            v_org_id, v_sale_id, v_channel_mobile,
            'MOB-' || LPAD(FLOOR(RANDOM() * 100000)::TEXT, 5, '0'),
            'MO-2024-' || LPAD(FLOOR(RANDOM() * 10000)::TEXT, 5, '0'),
            'synced', NOW() - INTERVAL '30 minutes',
            jsonb_build_object(
                'device_type', CASE WHEN RANDOM() < 0.5 THEN 'iOS' ELSE 'Android' END,
                'app_version', '2.5.1',
                'push_token', 'token-' || gen_random_uuid()::TEXT
            )
        );
    END LOOP;

    -- Update some sales with channel references
    UPDATE sales SET sales_channel_id = v_channel_web WHERE id IN (
        SELECT sale_id FROM external_order_mappings WHERE sales_channel_id = v_channel_web
    );

    UPDATE sales SET sales_channel_id = v_channel_mobile WHERE id IN (
        SELECT sale_id FROM external_order_mappings WHERE sales_channel_id = v_channel_mobile
    );

    RAISE NOTICE 'Created 8 external order mappings for web and mobile channels';
END $$;

\echo '';
\echo '==========================================';
\echo 'Advanced POS Features Seed Data Summary';
\echo '==========================================';
\echo 'POS Sessions:                2 sessions (1 closed, 1 open)';
\echo 'Cash Drawers:                1 drawer with movements';
\echo 'Cash Movements:              4 movements (float, pay-in, pay-out, drop)';
\echo 'Price Lists:                 4 lists (retail, wholesale, VIP, seasonal)';
\echo 'Price List Items:            15+ items';
\echo 'Product Components:          1 bundle with 4 components';
\echo 'Gift Cards:                  3 cards with 6 transactions';
\echo 'Store Credit Accounts:       5 accounts with transactions';
\echo 'Sale Returns:                3 returns (completed, pending)';
\echo 'Return Items:                3 items';
\echo 'Purchase Orders:             3 POs (received, partial, draft)';
\echo 'PO Items:                    17 line items';
\echo 'Goods Receipts:              1 receipt with items';
\echo 'System Health Checks:        7 checks';
\echo 'Error Logs:                  5 error entries';
\echo 'External Order Mappings:     8 mappings (web + mobile)';
\echo '==========================================';
\echo '';
\echo 'All advanced POS feature data loaded successfully!';

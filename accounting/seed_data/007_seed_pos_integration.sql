-- ============================================================================
-- Seed Data: 007 - POS Integration Data
-- Description: Sample data for POS → Accounting integration features
-- Dependencies: Requires accounting V001-V010 and POS V001-V022
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

DO $$
DECLARE
    v_org_id UUID;
    v_user_id UUID;
    v_product_id_1 UUID;
    v_product_id_2 UUID;
    v_category_id_1 UUID;
    -- GL Accounts
    v_revenue_account_id UUID;
    v_cogs_account_id UUID;
    v_inventory_account_id UUID;
    v_cash_account_id UUID;
    v_ar_account_id UUID;
    v_discount_account_id UUID;
    v_vat_liability_account_id UUID;
    v_gift_card_liability_account_id UUID;
    -- Tax
    v_tax_id_vat15 UUID;
    v_tax_id_exempt UUID;
BEGIN
    -- Get organization and user
    SELECT id INTO v_org_id FROM organizations LIMIT 1;
    SELECT id INTO v_user_id FROM users LIMIT 1;

    IF v_org_id IS NULL THEN
        RAISE EXCEPTION 'No organization found. Run core seed data first.';
    END IF;

    RAISE NOTICE 'Using Organization ID: %', v_org_id;

    -- ========================================================================
    -- GET GL ACCOUNTS (from existing chart of accounts)
    -- ========================================================================

    -- Revenue account
    SELECT id INTO v_revenue_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '4%' OR account_name ILIKE '%revenue%' OR account_name ILIKE '%sales%')
    LIMIT 1;

    -- COGS account
    SELECT id INTO v_cogs_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '5%' OR account_name ILIKE '%cost of sales%' OR account_name ILIKE '%cogs%')
    LIMIT 1;

    -- Inventory account
    SELECT id INTO v_inventory_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '1%' OR account_name ILIKE '%inventory%')
    LIMIT 1;

    -- Cash account
    SELECT id INTO v_cash_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '1%' OR account_name ILIKE '%cash%')
    LIMIT 1;

    -- AR account
    SELECT id INTO v_ar_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '1%' OR account_name ILIKE '%receivable%')
    LIMIT 1;

    -- Discount expense account
    SELECT id INTO v_discount_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '5%' OR account_name ILIKE '%discount%')
    LIMIT 1;

    -- VAT liability account
    SELECT id INTO v_vat_liability_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '2%' OR account_name ILIKE '%vat%' OR account_name ILIKE '%tax payable%')
    LIMIT 1;

    -- Gift card liability account
    SELECT id INTO v_gift_card_liability_account_id
    FROM chart_of_accounts
    WHERE organization_id = v_org_id
      AND (account_code LIKE '2%' OR account_name ILIKE '%gift card%' OR account_name ILIKE '%deferred%')
    LIMIT 1;

    -- Get some products and categories
    SELECT id INTO v_product_id_1 FROM products WHERE organization_id = v_org_id LIMIT 1;
    SELECT id INTO v_product_id_2 FROM products WHERE organization_id = v_org_id LIMIT 1 OFFSET 1;
    SELECT id INTO v_category_id_1 FROM categories WHERE organization_id = v_org_id LIMIT 1;

    -- Get taxes
    SELECT id INTO v_tax_id_vat15 FROM taxes WHERE organization_id = v_org_id AND tax_rate = 15 LIMIT 1;
    SELECT id INTO v_tax_id_exempt FROM taxes WHERE organization_id = v_org_id AND tax_rate = 0 LIMIT 1;

    -- ========================================================================
    -- SEED: POS Account Mappings
    -- ========================================================================

    RAISE NOTICE 'Seeding POS account mappings...';

    -- Product-specific mappings
    IF v_product_id_1 IS NOT NULL AND v_revenue_account_id IS NOT NULL THEN
        INSERT INTO pos_account_mappings (organization_id, source_type, source_id, purpose, account_id, is_default, priority, description, created_by, updated_by)
        VALUES
            (v_org_id, 'product', v_product_id_1, 'revenue', v_revenue_account_id, FALSE, 100, 'Product 1 revenue account', v_user_id, v_user_id),
            (v_org_id, 'product', v_product_id_1, 'cogs', v_cogs_account_id, FALSE, 100, 'Product 1 COGS account', v_user_id, v_user_id),
            (v_org_id, 'product', v_product_id_1, 'inventory', v_inventory_account_id, FALSE, 100, 'Product 1 inventory account', v_user_id, v_user_id);
    END IF;

    -- Category default mappings
    IF v_category_id_1 IS NOT NULL AND v_revenue_account_id IS NOT NULL THEN
        INSERT INTO pos_account_mappings (organization_id, source_type, source_id, purpose, account_id, is_default, priority, description, created_by, updated_by)
        VALUES
            (v_org_id, 'category', v_category_id_1, 'revenue', v_revenue_account_id, FALSE, 50, 'Category revenue fallback', v_user_id, v_user_id),
            (v_org_id, 'category', v_category_id_1, 'cogs', v_cogs_account_id, FALSE, 50, 'Category COGS fallback', v_user_id, v_user_id);
    END IF;

    -- Payment method mappings
    INSERT INTO pos_account_mappings (organization_id, source_type, source_code, purpose, account_id, is_default, priority, description, created_by, updated_by)
    VALUES
        (v_org_id, 'payment_method', 'CASH', 'asset', v_cash_account_id, FALSE, 100, 'Cash payments to cash account', v_user_id, v_user_id),
        (v_org_id, 'payment_method', 'CARD', 'asset', v_cash_account_id, FALSE, 100, 'Card payments to cash/bank account', v_user_id, v_user_id),
        (v_org_id, 'payment_method', 'CREDIT', 'asset', v_ar_account_id, FALSE, 100, 'Credit sales to AR account', v_user_id, v_user_id);

    -- Discount mapping
    INSERT INTO pos_account_mappings (organization_id, source_type, source_code, purpose, account_id, is_default, priority, description, created_by, updated_by)
    VALUES
        (v_org_id, 'discount', 'SALES_DISCOUNT', 'discount_expense', v_discount_account_id, FALSE, 100, 'Sales discounts expense', v_user_id, v_user_id);

    -- Gift card mapping
    IF v_gift_card_liability_account_id IS NOT NULL THEN
        INSERT INTO pos_account_mappings (organization_id, source_type, source_code, purpose, account_id, is_default, priority, description, created_by, updated_by)
        VALUES
            (v_org_id, 'gift_card', 'GIFT_CARD', 'liability', v_gift_card_liability_account_id, FALSE, 100, 'Gift card liability', v_user_id, v_user_id);
    END IF;

    -- Organization-wide defaults
    INSERT INTO pos_account_mappings (organization_id, source_type, source_id, purpose, account_id, is_default, priority, description, created_by, updated_by)
    VALUES
        (v_org_id, 'default', NULL, 'revenue', v_revenue_account_id, TRUE, 0, 'Default revenue account', v_user_id, v_user_id),
        (v_org_id, 'default', NULL, 'cogs', v_cogs_account_id, TRUE, 0, 'Default COGS account', v_user_id, v_user_id),
        (v_org_id, 'default', NULL, 'inventory', v_inventory_account_id, TRUE, 0, 'Default inventory account', v_user_id, v_user_id);

    RAISE NOTICE 'Created % POS account mappings', (SELECT COUNT(*) FROM pos_account_mappings WHERE organization_id = v_org_id);

    -- ========================================================================
    -- SEED: POS Tax Mappings
    -- ========================================================================

    RAISE NOTICE 'Seeding POS tax mappings...';

    IF v_tax_id_vat15 IS NOT NULL THEN
        INSERT INTO pos_tax_mappings (organization_id, pos_tax_code, tax_category_code, pos_tax_rate, accounting_tax_id, default_tax_account_id, is_default, is_active, applies_to_sales, description, created_by, updated_by)
        VALUES
            (v_org_id, 'VAT_15', 'S', 15.00, v_tax_id_vat15, v_vat_liability_account_id, TRUE, TRUE, TRUE, 'Standard VAT 15%', v_user_id, v_user_id),
            (v_org_id, 'SALES_TAX', 'S', 15.00, v_tax_id_vat15, v_vat_liability_account_id, FALSE, TRUE, TRUE, 'Sales tax 15%', v_user_id, v_user_id);
    END IF;

    IF v_tax_id_exempt IS NOT NULL THEN
        INSERT INTO pos_tax_mappings (organization_id, pos_tax_code, tax_category_code, pos_tax_rate, accounting_tax_id, is_default, is_active, applies_to_sales, description, created_by, updated_by)
        VALUES
            (v_org_id, 'EXEMPT', 'E', 0.00, v_tax_id_exempt, FALSE, TRUE, TRUE, 'Tax exempt items', v_user_id, v_user_id),
            (v_org_id, 'ZERO_RATED', 'Z', 0.00, v_tax_id_exempt, FALSE, TRUE, TRUE, 'Zero-rated items (exports)', v_user_id, v_user_id);
    END IF;

    RAISE NOTICE 'Created % POS tax mappings', (SELECT COUNT(*) FROM pos_tax_mappings WHERE organization_id = v_org_id);

    -- ========================================================================
    -- SEED: Inventory Valuation Settings
    -- ========================================================================

    RAISE NOTICE 'Seeding inventory valuation settings...';

    INSERT INTO inventory_valuation_settings (
        organization_id,
        valuation_method,
        cost_layer_granularity,
        default_inventory_account_id,
        default_cogs_account_id,
        cogs_recognition_timing,
        allow_negative_inventory,
        revalue_on_purchase,
        revaluation_frequency,
        is_active,
        notes,
        created_by,
        updated_by
    ) VALUES (
        v_org_id,
        'fifo',
        'product_location',
        v_inventory_account_id,
        v_cogs_account_id,
        'on_sale',
        FALSE,
        TRUE,
        'real_time',
        TRUE,
        'FIFO valuation with real-time cost updates',
        v_user_id,
        v_user_id
    );

    -- Create sample cost layers for products
    IF v_product_id_1 IS NOT NULL THEN
        INSERT INTO inventory_cost_layers (
            organization_id,
            product_id,
            location_id,
            layer_date,
            unit_cost,
            original_quantity,
            remaining_quantity,
            source_transaction_type,
            source_reference
        ) VALUES
            (v_org_id, v_product_id_1, NULL, CURRENT_DATE - INTERVAL '30 days', 10.50, 100, 75, 'purchase', 'PO-2024-001'),
            (v_org_id, v_product_id_1, NULL, CURRENT_DATE - INTERVAL '15 days', 11.00, 50, 50, 'purchase', 'PO-2024-002'),
            (v_org_id, v_product_id_1, NULL, CURRENT_DATE - INTERVAL '5 days', 10.75, 75, 75, 'purchase', 'PO-2024-003');
    END IF;

    RAISE NOTICE 'Created inventory valuation settings and sample cost layers';

    -- ========================================================================
    -- No explicit commit here - will be handled by outer transaction
    -- ========================================================================

    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'POS Integration Seed Data Summary:';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'POS Account Mappings: %', (SELECT COUNT(*) FROM pos_account_mappings WHERE organization_id = v_org_id);
    RAISE NOTICE 'POS Tax Mappings: %', (SELECT COUNT(*) FROM pos_tax_mappings WHERE organization_id = v_org_id);
    RAISE NOTICE 'Inventory Valuation Settings: 1';
    RAISE NOTICE 'Inventory Cost Layers: %', (SELECT COUNT(*) FROM inventory_cost_layers WHERE organization_id = v_org_id);
    RAISE NOTICE '============================================';

END $$;

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Seed Data 007 completed successfully!';
    RAISE NOTICE 'POS Integration data created.';
    RAISE NOTICE '============================================';
END $$;

-- ============================================================================
-- Seed Data: 008 - Document Sequences and Organization Features
-- Description: Sample data for document numbering and feature flags
-- Dependencies: Requires V019 (document_sequences) and V021 (organization_features)
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

DO $$
DECLARE
    v_org_id UUID;
    v_user_id UUID;
    v_location_id UUID;
BEGIN
    -- Get organization and user
    SELECT id INTO v_org_id FROM organizations LIMIT 1;
    SELECT id INTO v_user_id FROM users LIMIT 1;
    SELECT id INTO v_location_id FROM locations WHERE organization_id = v_org_id LIMIT 1;

    IF v_org_id IS NULL THEN
        RAISE EXCEPTION 'No organization found. Run core seed data first.';
    END IF;

    RAISE NOTICE 'Using Organization ID: %', v_org_id;

    -- ========================================================================
    -- SEED: Document Sequences
    -- ========================================================================

    RAISE NOTICE 'Seeding document sequences...';

    INSERT INTO document_sequences (
        organization_id,
        document_type,
        prefix,
        suffix,
        next_number,
        padding,
        reset_frequency,
        include_date,
        date_format,
        location_id,
        is_active,
        description,
        created_by,
        updated_by
    ) VALUES
        -- Sales sequences
        (v_org_id, 'sales', 'INV-', '', 1001, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Main sales/invoice sequence', v_user_id, v_user_id),
        (v_org_id, 'sales_receipt', 'RCT-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Sales receipt sequence', v_user_id, v_user_id),

        -- Purchase sequences
        (v_org_id, 'purchase_order', 'PO-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Purchase order sequence', v_user_id, v_user_id),
        (v_org_id, 'goods_receipt', 'GR-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Goods receipt sequence', v_user_id, v_user_id),

        -- Payment sequences
        (v_org_id, 'payment', 'PAY-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Payment sequence', v_user_id, v_user_id),
        (v_org_id, 'refund', 'REF-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Refund sequence', v_user_id, v_user_id),

        -- Accounting sequences
        (v_org_id, 'journal_entry', 'JE-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Journal entry sequence', v_user_id, v_user_id),
        (v_org_id, 'customer_invoice', 'CINV-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Customer invoice (accounting)', v_user_id, v_user_id),
        (v_org_id, 'vendor_bill', 'BILL-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Vendor bill (accounting)', v_user_id, v_user_id),

        -- Inventory sequences
        (v_org_id, 'inventory_transfer', 'TRF-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Inventory transfer sequence', v_user_id, v_user_id),
        (v_org_id, 'inventory_adjustment', 'ADJ-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'Inventory adjustment sequence', v_user_id, v_user_id),

        -- POS session sequence
        (v_org_id, 'pos_session', 'SES-', '', 1, 6, 'yearly', TRUE, 'YYYY', NULL, TRUE, 'POS session sequence', v_user_id, v_user_id),

        -- Location-specific sale sequence (if location exists)
        (v_org_id, 'sales', 'INV-LOC1-', '', 1, 4, 'daily', FALSE, NULL, v_location_id, TRUE, 'Location-specific sales', v_user_id, v_user_id);

    -- Update example_number for all sequences
    UPDATE document_sequences
    SET example_number = preview_document_number(organization_id, document_type, location_id, CURRENT_DATE)
    WHERE organization_id = v_org_id;

    RAISE NOTICE 'Created % document sequences', (SELECT COUNT(*) FROM document_sequences WHERE organization_id = v_org_id);

    -- ========================================================================
    -- SEED: Organization Features (JSONB update)
    -- ========================================================================

    RAISE NOTICE 'Updating organization features (JSONB)...';

    UPDATE organizations
    SET features = '{
        "accounting": true,
        "e_invoicing": true,
        "advanced_inventory": true,
        "multi_location": true,
        "multi_currency": true,
        "loyalty_program": true,
        "delivery_management": true,
        "restaurant_mode": false,
        "table_management": false,
        "kitchen_display": false,
        "online_ordering": true,
        "api_access": true
    }'::jsonb,
    updated_at = CURRENT_TIMESTAMP
    WHERE id = v_org_id;

    -- ========================================================================
    -- SEED: Organization Features (Normalized table)
    -- ========================================================================

    RAISE NOTICE 'Seeding organization features (normalized)...';

    INSERT INTO organization_features (
        organization_id,
        feature_key,
        is_enabled,
        is_available,
        configuration,
        limits,
        enabled_at,
        created_by,
        updated_by
    ) VALUES
        -- Core features
        (v_org_id, 'accounting', TRUE, TRUE, '{"auto_posting": true, "require_approval": false}'::jsonb, '{"max_journal_entries_per_month": null}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),
        (v_org_id, 'e_invoicing', TRUE, TRUE, '{"authorities": ["ZATCA", "ETA"], "auto_submit": false}'::jsonb, '{"max_invoices_per_month": null}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),
        (v_org_id, 'multi_currency', TRUE, TRUE, '{"base_currency": "USD", "auto_update_rates": true}'::jsonb, '{"max_currencies": 10}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),

        -- Inventory features
        (v_org_id, 'advanced_inventory', TRUE, TRUE, '{"batch_tracking": true, "serial_tracking": true, "expiry_tracking": true}'::jsonb, '{}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),
        (v_org_id, 'multi_location', TRUE, TRUE, '{"auto_transfer": false, "require_approval": true}'::jsonb, '{"max_locations": 20}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),

        -- Customer features
        (v_org_id, 'loyalty_program', TRUE, TRUE, '{"points_per_dollar": 1, "redemption_ratio": 0.01}'::jsonb, '{"max_active_members": 10000}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),
        (v_org_id, 'online_ordering', TRUE, TRUE, '{"accept_online_payments": true, "auto_confirm": false}'::jsonb, '{}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),

        -- Delivery features
        (v_org_id, 'delivery_management', TRUE, TRUE, '{"track_drivers": true, "route_optimization": false}'::jsonb, '{"max_drivers": 50}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),

        -- API access
        (v_org_id, 'api_access', TRUE, TRUE, '{"rate_limit": 1000, "webhook_enabled": true}'::jsonb, '{"requests_per_hour": 1000}'::jsonb, CURRENT_TIMESTAMP, v_user_id, v_user_id),

        -- Restaurant features (disabled)
        (v_org_id, 'restaurant_mode', FALSE, TRUE, '{}'::jsonb, '{}'::jsonb, NULL, v_user_id, v_user_id),
        (v_org_id, 'table_management', FALSE, TRUE, '{}'::jsonb, '{"max_tables": 100}'::jsonb, NULL, v_user_id, v_user_id),
        (v_org_id, 'kitchen_display', FALSE, TRUE, '{}'::jsonb, '{"max_screens": 10}'::jsonb, NULL, v_user_id, v_user_id);

    RAISE NOTICE 'Created % organization features', (SELECT COUNT(*) FROM organization_features WHERE organization_id = v_org_id);

    -- ========================================================================
    -- UPDATE: Organization with base currency and display settings
    -- ========================================================================

    RAISE NOTICE 'Updating organization currency settings...';

    UPDATE organizations
    SET
        base_currency_code = 'USD',
        currency_display_format = 'symbol',
        decimal_separator = '.',
        thousands_separator = ',',
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_org_id;

    -- ========================================================================
    -- SUMMARY
    -- ========================================================================

    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Sequences and Features Seed Data Summary:';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Document Sequences: %', (SELECT COUNT(*) FROM document_sequences WHERE organization_id = v_org_id);
    RAISE NOTICE 'Organization Features: %', (SELECT COUNT(*) FROM organization_features WHERE organization_id = v_org_id);
    RAISE NOTICE 'Base Currency: USD';
    RAISE NOTICE 'Features Enabled: accounting, e_invoicing, multi_currency, advanced_inventory';
    RAISE NOTICE '============================================';

    -- Test sequence generation
    RAISE NOTICE 'Testing sequence generation:';
    RAISE NOTICE 'Next Sales Invoice: %', get_next_document_number(v_org_id, 'sales');
    RAISE NOTICE 'Next Purchase Order: %', get_next_document_number(v_org_id, 'purchase_order');
    RAISE NOTICE 'Next Journal Entry: %', get_next_document_number(v_org_id, 'journal_entry');

END $$;

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Seed Data 008 completed successfully!';
    RAISE NOTICE 'Document Sequences and Features created.';
    RAISE NOTICE '============================================';
END $$;

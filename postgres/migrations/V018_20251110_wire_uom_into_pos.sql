-- ============================================================================
-- Migration: V018 - Wire UOM into POS Tables
-- Description: Add foreign key references to units_of_measure table throughout POS
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- PRODUCTS TABLE - Add UOM Foreign Key
-- ============================================================================

-- Add base UOM reference
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS base_uom_id UUID REFERENCES units_of_measure(id);

-- Add index
CREATE INDEX IF NOT EXISTS idx_products_base_uom_id ON products(base_uom_id) WHERE base_uom_id IS NOT NULL;

COMMENT ON COLUMN products.base_uom_id IS 'Base unit of measure for this product (references units_of_measure table)';

-- ============================================================================
-- SALE_ITEMS TABLE - Add UOM Foreign Key
-- ============================================================================

ALTER TABLE sale_items
    ADD COLUMN IF NOT EXISTS uom_id UUID REFERENCES units_of_measure(id);

CREATE INDEX IF NOT EXISTS idx_sale_items_uom_id ON sale_items(uom_id) WHERE uom_id IS NOT NULL;

COMMENT ON COLUMN sale_items.uom_id IS 'Unit of measure for this sale line (may differ from product base UOM via conversion)';

-- ============================================================================
-- INVENTORY_TRANSACTIONS TABLE - Add UOM Foreign Key
-- ============================================================================

ALTER TABLE inventory_transactions
    ADD COLUMN IF NOT EXISTS uom_id UUID REFERENCES units_of_measure(id);

CREATE INDEX IF NOT EXISTS idx_inventory_transactions_uom_id ON inventory_transactions(uom_id) WHERE uom_id IS NOT NULL;

COMMENT ON COLUMN inventory_transactions.uom_id IS 'Unit of measure for this inventory transaction';

-- ============================================================================
-- PURCHASE_ORDER_ITEMS TABLE - Add UOM Foreign Key (if exists)
-- ============================================================================

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'purchase_order_items') THEN
        ALTER TABLE purchase_order_items
            ADD COLUMN IF NOT EXISTS uom_id UUID REFERENCES units_of_measure(id);

        CREATE INDEX IF NOT EXISTS idx_purchase_order_items_uom_id ON purchase_order_items(uom_id) WHERE uom_id IS NOT NULL;
    END IF;
END $$;

-- ============================================================================
-- GOODS_RECEIPT_ITEMS TABLE - Add UOM Foreign Key (if exists)
-- ============================================================================

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'goods_receipt_items') THEN
        ALTER TABLE goods_receipt_items
            ADD COLUMN IF NOT EXISTS uom_id UUID REFERENCES units_of_measure(id);

        CREATE INDEX IF NOT EXISTS idx_goods_receipt_items_uom_id ON goods_receipt_items(uom_id) WHERE uom_id IS NOT NULL;
    END IF;
END $$;

-- ============================================================================
-- PRODUCT_BATCHES TABLE - Add UOM Foreign Key (if exists)
-- ============================================================================

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'product_batches') THEN
        ALTER TABLE product_batches
            ADD COLUMN IF NOT EXISTS uom_id UUID REFERENCES units_of_measure(id);

        CREATE INDEX IF NOT EXISTS idx_product_batches_uom_id ON product_batches(uom_id) WHERE uom_id IS NOT NULL;
    END IF;
END $$;

-- ============================================================================
-- HELPER FUNCTION: Convert Quantity Between UOMs
-- ============================================================================

CREATE OR REPLACE FUNCTION convert_uom_quantity(
    p_quantity NUMERIC,
    p_from_uom_id UUID,
    p_to_uom_id UUID
) RETURNS NUMERIC AS $$
DECLARE
    v_from_uom_code VARCHAR;
    v_to_uom_code VARCHAR;
    v_conversion_factor NUMERIC;
    v_result NUMERIC;
BEGIN
    -- If same UOM, return original quantity
    IF p_from_uom_id = p_to_uom_id THEN
        RETURN p_quantity;
    END IF;

    -- Get UOM codes
    SELECT uom_code INTO v_from_uom_code FROM units_of_measure WHERE id = p_from_uom_id;
    SELECT uom_code INTO v_to_uom_code FROM units_of_measure WHERE id = p_to_uom_id;

    IF v_from_uom_code IS NULL OR v_to_uom_code IS NULL THEN
        RAISE EXCEPTION 'Invalid UOM IDs provided';
    END IF;

    -- Look up conversion factor
    SELECT conversion_factor INTO v_conversion_factor
    FROM uom_conversions
    WHERE from_uom_code = v_from_uom_code
      AND to_uom_code = v_to_uom_code
      AND is_active = TRUE
      AND deleted_at IS NULL
    LIMIT 1;

    -- If no direct conversion found, try reverse
    IF v_conversion_factor IS NULL THEN
        SELECT 1.0 / conversion_factor INTO v_conversion_factor
        FROM uom_conversions
        WHERE from_uom_code = v_to_uom_code
          AND to_uom_code = v_from_uom_code
          AND is_active = TRUE
          AND deleted_at IS NULL
        LIMIT 1;
    END IF;

    -- If still no conversion, return original (or raise error in strict mode)
    IF v_conversion_factor IS NULL THEN
        RAISE WARNING 'No UOM conversion found from % to %. Returning original quantity.', v_from_uom_code, v_to_uom_code;
        RETURN p_quantity;
    END IF;

    -- Apply conversion
    v_result := p_quantity * v_conversion_factor;

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION convert_uom_quantity IS 'Convert quantity between different units of measure using uom_conversions table';

-- ============================================================================
-- HELPER FUNCTION: Get Product Base UOM
-- ============================================================================

CREATE OR REPLACE FUNCTION get_product_base_uom(
    p_product_id UUID
) RETURNS UUID AS $$
DECLARE
    v_uom_id UUID;
BEGIN
    SELECT base_uom_id INTO v_uom_id
    FROM products
    WHERE id = p_product_id;

    RETURN v_uom_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_product_base_uom IS 'Get the base unit of measure for a product';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V018 completed successfully!';
    RAISE NOTICE 'Added UOM foreign keys to:';
    RAISE NOTICE ' - products (base_uom_id)';
    RAISE NOTICE ' - sale_items (uom_id)';
    RAISE NOTICE ' - inventory_transactions (uom_id)';
    RAISE NOTICE ' - purchase_order_items (uom_id)';
    RAISE NOTICE ' - goods_receipt_items (uom_id)';
    RAISE NOTICE ' - product_batches (uom_id)';
    RAISE NOTICE ' - Helper functions for UOM conversion';
    RAISE NOTICE '============================================';
END $$;

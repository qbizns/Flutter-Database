-- ============================================================================
-- Migration: V016 - Add POS Tax Codes
-- Description: Add explicit tax code fields to products and sales for accounting integration
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- PRODUCTS TABLE - Add POS Tax Code
-- ============================================================================

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS pos_tax_code VARCHAR(50),
    ADD COLUMN IF NOT EXISTS tax_behavior VARCHAR(50) DEFAULT 'taxable' CHECK (tax_behavior IN (
        'taxable', 'exempt', 'zero_rated', 'out_of_scope'
    ));

CREATE INDEX IF NOT EXISTS idx_products_pos_tax_code ON products(pos_tax_code) WHERE pos_tax_code IS NOT NULL;

COMMENT ON COLUMN products.pos_tax_code IS 'POS tax code for accounting integration (VAT_15, SALES_TAX, etc.)';
COMMENT ON COLUMN products.tax_behavior IS 'Tax treatment: taxable, exempt, zero_rated, out_of_scope';

-- ============================================================================
-- SALE_ITEMS TABLE - Add POS Tax Code (snapshot at time of sale)
-- ============================================================================

ALTER TABLE sale_items
    ADD COLUMN IF NOT EXISTS pos_tax_code VARCHAR(50),
    ADD COLUMN IF NOT EXISTS tax_behavior VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_sale_items_pos_tax_code ON sale_items(pos_tax_code) WHERE pos_tax_code IS NOT NULL;

COMMENT ON COLUMN sale_items.pos_tax_code IS 'Tax code snapshot at time of sale for accurate accounting posting';

-- ============================================================================
-- SALES TABLE - Add Default Tax Code
-- ============================================================================

ALTER TABLE sales
    ADD COLUMN IF NOT EXISTS default_tax_code VARCHAR(50);

COMMENT ON COLUMN sales.default_tax_code IS 'Default tax code for this sale (can be overridden at line level)';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V016 completed successfully!';
    RAISE NOTICE 'Added POS tax code fields to:';
    RAISE NOTICE ' - products (pos_tax_code, tax_behavior)';
    RAISE NOTICE ' - sale_items (pos_tax_code, tax_behavior)';
    RAISE NOTICE ' - sales (default_tax_code)';
    RAISE NOTICE '============================================';
END $$;

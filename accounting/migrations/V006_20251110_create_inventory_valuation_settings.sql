-- ============================================================================
-- Migration: V006 - Create Inventory Valuation Settings
-- Description: Configuration for inventory valuation methods (FIFO, weighted average, etc.)
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- INVENTORY VALUATION SETTINGS TABLE
-- Description: Per-organization configuration for inventory valuation
-- ============================================================================

CREATE TABLE IF NOT EXISTS inventory_valuation_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Valuation Method
    valuation_method VARCHAR(50) NOT NULL DEFAULT 'fifo' CHECK (valuation_method IN (
        'fifo',              -- First In, First Out
        'lifo',              -- Last In, First Out (rare, not allowed in IFRS)
        'weighted_average',  -- Weighted Average Cost
        'moving_average',    -- Moving Average Cost
        'standard_cost',     -- Standard/Fixed Cost
        'specific_id'        -- Specific Identification (for serialized items)
    )),

    -- Cost Layer Granularity
    cost_layer_granularity VARCHAR(50) DEFAULT 'product_location' CHECK (cost_layer_granularity IN (
        'product',              -- One cost layer per product across all locations
        'product_location',     -- Separate cost layers per product per location
        'product_location_lot', -- Per lot/batch
        'serial_number'         -- Per serial number (for specific_id method)
    )),

    -- Default GL Accounts
    default_inventory_account_id UUID REFERENCES chart_of_accounts(id),
    default_cogs_account_id UUID REFERENCES chart_of_accounts(id),
    default_inventory_adjustment_account_id UUID REFERENCES chart_of_accounts(id),
    default_inventory_variance_account_id UUID REFERENCES chart_of_accounts(id),

    -- COGS Recognition Timing
    cogs_recognition_timing VARCHAR(50) DEFAULT 'on_sale' CHECK (cogs_recognition_timing IN (
        'on_sale',      -- When sale is recorded
        'on_delivery',  -- When goods are delivered
        'on_payment'    -- When payment is received (rare)
    )),

    -- Valuation Options
    allow_negative_inventory BOOLEAN DEFAULT FALSE,  -- Allow negative stock valuation
    revalue_on_purchase BOOLEAN DEFAULT TRUE,        -- Recalculate average cost on each purchase
    round_unit_cost_to_decimals INTEGER DEFAULT 4,   -- Decimal precision for unit costs

    -- Period Settings
    revaluation_frequency VARCHAR(50) DEFAULT 'real_time' CHECK (revaluation_frequency IN (
        'real_time',  -- Immediate recalculation
        'daily',      -- Daily batch
        'monthly',    -- Monthly batch
        'manual'      -- Only on manual trigger
    )),

    -- Configuration
    is_active BOOLEAN DEFAULT TRUE,
    effective_from DATE DEFAULT CURRENT_DATE,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id)  -- One active config per organization
);

-- Indexes
CREATE INDEX idx_inventory_valuation_settings_organization_id ON inventory_valuation_settings(organization_id);
CREATE INDEX idx_inventory_valuation_settings_valuation_method ON inventory_valuation_settings(valuation_method);

-- Trigger
CREATE TRIGGER update_inventory_valuation_settings_updated_at
    BEFORE UPDATE ON inventory_valuation_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE inventory_valuation_settings IS 'Per-organization configuration for inventory valuation methods and COGS calculation';
COMMENT ON COLUMN inventory_valuation_settings.valuation_method IS 'Inventory costing method: FIFO, LIFO, weighted average, etc.';
COMMENT ON COLUMN inventory_valuation_settings.cost_layer_granularity IS 'Level of cost tracking: product, product+location, lot, or serial';
COMMENT ON COLUMN inventory_valuation_settings.cogs_recognition_timing IS 'When to recognize COGS in accounting';

-- ============================================================================
-- INVENTORY COST LAYERS TABLE
-- Description: Track cost layers for FIFO/LIFO valuation
-- ============================================================================

CREATE TABLE IF NOT EXISTS inventory_cost_layers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Product & Location
    product_id UUID NOT NULL,  -- References products(id) from POS schema
    location_id UUID,          -- References locations(id) from POS schema, NULL = company-wide
    lot_number VARCHAR(100),   -- For lot-tracked items
    serial_number VARCHAR(100), -- For serialized items

    -- Cost Layer Details
    layer_date DATE NOT NULL DEFAULT CURRENT_DATE,
    unit_cost NUMERIC(20, 6) NOT NULL,
    original_quantity NUMERIC(12, 3) NOT NULL,
    remaining_quantity NUMERIC(12, 3) NOT NULL,
    uom_code VARCHAR(10),

    -- Source Transaction
    source_transaction_type VARCHAR(50),  -- 'purchase', 'transfer_in', 'adjustment', 'production'
    source_transaction_id UUID,
    source_reference VARCHAR(255),

    -- Layer Status
    is_fully_consumed BOOLEAN DEFAULT FALSE,
    consumed_at TIMESTAMP WITH TIME ZONE,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CHECK (remaining_quantity >= 0),
    CHECK (remaining_quantity <= original_quantity)
);

-- Indexes
CREATE INDEX idx_inventory_cost_layers_organization_id ON inventory_cost_layers(organization_id);
CREATE INDEX idx_inventory_cost_layers_product_location ON inventory_cost_layers(product_id, location_id);
CREATE INDEX idx_inventory_cost_layers_layer_date ON inventory_cost_layers(layer_date);
CREATE INDEX idx_inventory_cost_layers_remaining ON inventory_cost_layers(remaining_quantity) WHERE remaining_quantity > 0;
CREATE INDEX idx_inventory_cost_layers_lot ON inventory_cost_layers(lot_number) WHERE lot_number IS NOT NULL;
CREATE INDEX idx_inventory_cost_layers_serial ON inventory_cost_layers(serial_number) WHERE serial_number IS NOT NULL;

-- Trigger
CREATE TRIGGER update_inventory_cost_layers_updated_at
    BEFORE UPDATE ON inventory_cost_layers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE inventory_cost_layers IS 'Cost layers for FIFO/LIFO inventory valuation';
COMMENT ON COLUMN inventory_cost_layers.remaining_quantity IS 'Quantity still available in this cost layer';
COMMENT ON COLUMN inventory_cost_layers.unit_cost IS 'Unit cost for this specific layer';

-- ============================================================================
-- HELPER FUNCTION: Get Average Cost for Product
-- ============================================================================

CREATE OR REPLACE FUNCTION get_product_average_cost(
    p_organization_id UUID,
    p_product_id UUID,
    p_location_id UUID DEFAULT NULL
) RETURNS NUMERIC AS $$
DECLARE
    v_total_value NUMERIC;
    v_total_qty NUMERIC;
    v_avg_cost NUMERIC;
BEGIN
    SELECT
        SUM(unit_cost * remaining_quantity),
        SUM(remaining_quantity)
    INTO v_total_value, v_total_qty
    FROM inventory_cost_layers
    WHERE organization_id = p_organization_id
      AND product_id = p_product_id
      AND (p_location_id IS NULL OR location_id = p_location_id)
      AND remaining_quantity > 0
      AND deleted_at IS NULL;

    IF v_total_qty IS NULL OR v_total_qty = 0 THEN
        RETURN 0;
    END IF;

    v_avg_cost := v_total_value / v_total_qty;
    RETURN ROUND(v_avg_cost, 6);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_product_average_cost IS 'Calculate weighted average cost for a product (optionally at a location)';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

ALTER TABLE inventory_valuation_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory_cost_layers ENABLE ROW LEVEL SECURITY;

CREATE POLICY inventory_valuation_settings_tenant_isolation ON inventory_valuation_settings
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

CREATE POLICY inventory_cost_layers_tenant_isolation ON inventory_cost_layers
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V006 completed successfully!';
    RAISE NOTICE 'Inventory Valuation Settings created.';
    RAISE NOTICE ' - inventory_valuation_settings table';
    RAISE NOTICE ' - inventory_cost_layers table';
    RAISE NOTICE ' - get_product_average_cost() function';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

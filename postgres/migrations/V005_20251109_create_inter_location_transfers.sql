-- Migration V005: Inter-Location Transfers
-- Created: 2025-11-09
-- Description: Adds inter-location inventory transfer functionality
-- Dependencies: V001, V002, V003, V004

-- ============================================================================
-- INVENTORY TRANSFERS
-- ============================================================================

-- Inventory transfer header table for moving stock between locations
CREATE TABLE IF NOT EXISTS inventory_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Transfer Information
    transfer_number VARCHAR(50) NOT NULL,
    transfer_date DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Locations
    from_location_id UUID NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
    to_location_id UUID NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,

    -- Status Workflow
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'requested', 'approved', 'in_transit', 'partially_received', 'received', 'cancelled', 'rejected')),

    -- Dates
    requested_date TIMESTAMP WITH TIME ZONE,
    approved_date TIMESTAMP WITH TIME ZONE,
    shipped_date TIMESTAMP WITH TIME ZONE,
    expected_delivery_date DATE,
    received_date TIMESTAMP WITH TIME ZONE,

    -- Shipping Information
    carrier VARCHAR(100),
    tracking_number VARCHAR(100),
    shipping_cost NUMERIC(15, 2) DEFAULT 0,

    -- Notes and Documentation
    reason TEXT, -- Reason for transfer (restock, rebalance, customer order, etc.)
    notes TEXT,
    rejection_reason TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    requested_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    shipped_by UUID REFERENCES users(id),
    received_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_transfer_number_per_org UNIQUE (organization_id, transfer_number),
    CONSTRAINT different_locations CHECK (from_location_id != to_location_id),
    CONSTRAINT positive_shipping_cost CHECK (shipping_cost >= 0),
    CONSTRAINT valid_workflow_dates CHECK (
        (approved_date IS NULL OR requested_date IS NULL OR approved_date >= requested_date) AND
        (shipped_date IS NULL OR approved_date IS NULL OR shipped_date >= approved_date) AND
        (received_date IS NULL OR shipped_date IS NULL OR received_date >= shipped_date)
    )
);

-- Indexes for inventory_transfers
CREATE INDEX idx_inventory_transfers_organization_id ON inventory_transfers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventory_transfers_from_location ON inventory_transfers(from_location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventory_transfers_to_location ON inventory_transfers(to_location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventory_transfers_transfer_number ON inventory_transfers(transfer_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventory_transfers_status ON inventory_transfers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventory_transfers_transfer_date ON inventory_transfers(transfer_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventory_transfers_tracking ON inventory_transfers(tracking_number) WHERE deleted_at IS NULL AND tracking_number IS NOT NULL;

-- Auto-update trigger
CREATE TRIGGER update_inventory_transfers_updated_at
    BEFORE UPDATE ON inventory_transfers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE inventory_transfers IS 'Inventory transfers between locations with complete workflow tracking';

-- ============================================================================
-- INVENTORY TRANSFER ITEMS
-- ============================================================================

-- Line items for inventory transfers
CREATE TABLE IF NOT EXISTS inventory_transfer_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    inventory_transfer_id UUID NOT NULL REFERENCES inventory_transfers(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    product_variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,

    -- Product Snapshot (preserved even if product is deleted)
    product_name VARCHAR(200) NOT NULL,
    product_sku VARCHAR(100),

    -- Quantities
    quantity_requested NUMERIC(10, 2) NOT NULL,
    quantity_shipped NUMERIC(10, 2) DEFAULT 0,
    quantity_received NUMERIC(10, 2) DEFAULT 0,
    unit_of_measure VARCHAR(50) DEFAULT 'unit',

    -- Costing (for valuation purposes)
    unit_cost NUMERIC(15, 2),
    total_cost NUMERIC(15, 2),

    -- Item Status
    item_status VARCHAR(20) DEFAULT 'pending' CHECK (item_status IN ('pending', 'approved', 'shipped', 'partially_received', 'received', 'cancelled')),

    -- Variance Tracking
    variance_quantity NUMERIC(10, 2) DEFAULT 0, -- difference between shipped and received
    variance_reason TEXT,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT positive_quantities CHECK (
        quantity_requested > 0 AND
        quantity_shipped >= 0 AND
        quantity_received >= 0 AND
        quantity_shipped <= quantity_requested AND
        quantity_received <= quantity_shipped
    ),
    CONSTRAINT positive_costs CHECK (
        (unit_cost IS NULL OR unit_cost >= 0) AND
        (total_cost IS NULL OR total_cost >= 0)
    )
);

-- Indexes for inventory_transfer_items
CREATE INDEX idx_inventory_transfer_items_organization_id ON inventory_transfer_items(organization_id);
CREATE INDEX idx_inventory_transfer_items_transfer_id ON inventory_transfer_items(inventory_transfer_id);
CREATE INDEX idx_inventory_transfer_items_product_id ON inventory_transfer_items(product_id);
CREATE INDEX idx_inventory_transfer_items_variant_id ON inventory_transfer_items(product_variant_id);
CREATE INDEX idx_inventory_transfer_items_status ON inventory_transfer_items(item_status);

-- Auto-update trigger
CREATE TRIGGER update_inventory_transfer_items_updated_at
    BEFORE UPDATE ON inventory_transfer_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE inventory_transfer_items IS 'Line items for inventory transfers with quantity tracking';

-- ============================================================================
-- HELPER FUNCTION: Calculate Transfer Variance
-- ============================================================================

-- Function to automatically calculate variance when receiving items
CREATE OR REPLACE FUNCTION calculate_transfer_item_variance()
RETURNS TRIGGER AS $$
BEGIN
    -- Calculate variance as the difference between shipped and received
    NEW.variance_quantity := NEW.quantity_shipped - NEW.quantity_received;

    -- Update item status based on quantities
    IF NEW.quantity_received = 0 THEN
        NEW.item_status := 'shipped';
    ELSIF NEW.quantity_received < NEW.quantity_shipped THEN
        NEW.item_status := 'partially_received';
    ELSIF NEW.quantity_received = NEW.quantity_shipped THEN
        NEW.item_status := 'received';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to calculate variance automatically
CREATE TRIGGER trigger_calculate_transfer_variance
    BEFORE UPDATE ON inventory_transfer_items
    FOR EACH ROW
    WHEN (OLD.quantity_received IS DISTINCT FROM NEW.quantity_received OR
          OLD.quantity_shipped IS DISTINCT FROM NEW.quantity_shipped)
    EXECUTE FUNCTION calculate_transfer_item_variance();

-- ============================================================================
-- HELPER FUNCTION: Update Transfer Status
-- ============================================================================

-- Function to update transfer header status based on item statuses
CREATE OR REPLACE FUNCTION update_transfer_status()
RETURNS TRIGGER AS $$
DECLARE
    total_items INTEGER;
    received_items INTEGER;
    partially_received_items INTEGER;
    transfer_rec RECORD;
BEGIN
    -- Get the transfer record
    SELECT * INTO transfer_rec FROM inventory_transfers WHERE id = NEW.inventory_transfer_id;

    -- Count items
    SELECT
        COUNT(*) as total,
        COUNT(*) FILTER (WHERE item_status = 'received') as received,
        COUNT(*) FILTER (WHERE item_status = 'partially_received') as partial
    INTO total_items, received_items, partially_received_items
    FROM inventory_transfer_items
    WHERE inventory_transfer_id = NEW.inventory_transfer_id;

    -- Update transfer status based on item statuses
    IF received_items = total_items THEN
        -- All items received
        UPDATE inventory_transfers
        SET status = 'received',
            received_date = CASE WHEN received_date IS NULL THEN CURRENT_TIMESTAMP ELSE received_date END
        WHERE id = NEW.inventory_transfer_id AND status != 'received';
    ELSIF partially_received_items > 0 OR received_items > 0 THEN
        -- Some items received
        UPDATE inventory_transfers
        SET status = 'partially_received'
        WHERE id = NEW.inventory_transfer_id AND status = 'in_transit';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update transfer status
CREATE TRIGGER trigger_update_transfer_status
    AFTER UPDATE ON inventory_transfer_items
    FOR EACH ROW
    WHEN (OLD.item_status IS DISTINCT FROM NEW.item_status)
    EXECUTE FUNCTION update_transfer_status();

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS on new tables
ALTER TABLE inventory_transfers ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory_transfer_items ENABLE ROW LEVEL SECURITY;

-- ============================================================================
-- RLS POLICIES: INVENTORY_TRANSFERS
-- ============================================================================

-- Super admin bypass
CREATE POLICY inventory_transfers_super_admin_all
    ON inventory_transfers FOR ALL TO PUBLIC
    USING (is_super_admin());

-- Organization-scoped policies
CREATE POLICY inventory_transfers_select_own_org
    ON inventory_transfers FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY inventory_transfers_insert_own_org
    ON inventory_transfers FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY inventory_transfers_update_own_org
    ON inventory_transfers FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY inventory_transfers_delete_own_org
    ON inventory_transfers FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: INVENTORY_TRANSFER_ITEMS
-- ============================================================================

CREATE POLICY inventory_transfer_items_super_admin_all
    ON inventory_transfer_items FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY inventory_transfer_items_select_own_org
    ON inventory_transfer_items FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY inventory_transfer_items_insert_own_org
    ON inventory_transfer_items FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY inventory_transfer_items_update_own_org
    ON inventory_transfer_items FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY inventory_transfer_items_delete_own_org
    ON inventory_transfer_items FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- VIEWS: Transfer Summary
-- ============================================================================

-- View for transfer summary with item counts and totals
CREATE OR REPLACE VIEW v_inventory_transfer_summary AS
SELECT
    t.id,
    t.organization_id,
    t.transfer_number,
    t.transfer_date,
    t.status,
    fl.name as from_location_name,
    tl.name as to_location_name,
    t.expected_delivery_date,
    t.tracking_number,
    COUNT(ti.id) as total_items,
    SUM(ti.quantity_requested) as total_quantity_requested,
    SUM(ti.quantity_shipped) as total_quantity_shipped,
    SUM(ti.quantity_received) as total_quantity_received,
    SUM(ti.total_cost) as total_transfer_value,
    t.created_at,
    t.updated_at,
    u.full_name as created_by_name,
    req.full_name as requested_by_name,
    app.full_name as approved_by_name,
    rec.full_name as received_by_name
FROM inventory_transfers t
LEFT JOIN locations fl ON t.from_location_id = fl.id
LEFT JOIN locations tl ON t.to_location_id = tl.id
LEFT JOIN inventory_transfer_items ti ON t.id = ti.inventory_transfer_id
LEFT JOIN users u ON t.created_by = u.id
LEFT JOIN users req ON t.requested_by = req.id
LEFT JOIN users app ON t.approved_by = app.id
LEFT JOIN users rec ON t.received_by = rec.id
WHERE t.deleted_at IS NULL
GROUP BY
    t.id, t.organization_id, t.transfer_number, t.transfer_date, t.status,
    fl.name, tl.name, t.expected_delivery_date, t.tracking_number,
    t.created_at, t.updated_at, u.full_name, req.full_name, app.full_name, rec.full_name;

COMMENT ON VIEW v_inventory_transfer_summary IS 'Summary view of inventory transfers with aggregated item data';

-- Migration completed successfully

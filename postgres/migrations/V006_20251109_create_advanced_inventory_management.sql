-- Migration V006: Advanced Inventory Management
-- Created: 2025-11-09
-- Description: Adds serial number tracking, batch/lot management, expiration tracking, and cycle counting
-- Dependencies: V001, V002, V003, V004, V005

-- ============================================================================
-- PRODUCT SERIAL NUMBERS
-- ============================================================================

-- Track individual items by serial number (for electronics, equipment, etc.)
CREATE TABLE IF NOT EXISTS product_serial_numbers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Serial Number Information
    serial_number VARCHAR(100) NOT NULL,
    status VARCHAR(20) DEFAULT 'in_stock' CHECK (status IN ('in_stock', 'sold', 'reserved', 'damaged', 'returned', 'in_repair', 'scrapped')),

    -- Purchase Information
    purchase_order_id UUID REFERENCES purchase_orders(id) ON DELETE SET NULL,
    purchase_date DATE,
    purchase_cost NUMERIC(15, 2),
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,

    -- Sale Information
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    sale_date DATE,
    sale_price NUMERIC(15, 2),
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,

    -- Warranty Information
    warranty_start_date DATE,
    warranty_end_date DATE,
    warranty_provider VARCHAR(200),
    warranty_terms TEXT,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_serial_number_per_org UNIQUE (organization_id, serial_number),
    CONSTRAINT valid_warranty_dates CHECK (warranty_end_date IS NULL OR warranty_start_date IS NULL OR warranty_end_date >= warranty_start_date)
);

-- Indexes
CREATE INDEX idx_product_serial_numbers_org_id ON product_serial_numbers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_serial_numbers_product_id ON product_serial_numbers(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_serial_numbers_serial ON product_serial_numbers(serial_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_serial_numbers_status ON product_serial_numbers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_serial_numbers_location ON product_serial_numbers(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_serial_numbers_sale ON product_serial_numbers(sale_id) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_product_serial_numbers_updated_at
    BEFORE UPDATE ON product_serial_numbers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE product_serial_numbers IS 'Track individual products by serial number for warranty and traceability';

-- ============================================================================
-- PRODUCT BATCHES / LOTS
-- ============================================================================

-- Batch/Lot tracking for products (food, pharmaceuticals, chemicals, etc.)
CREATE TABLE IF NOT EXISTS product_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Batch Information
    batch_number VARCHAR(100) NOT NULL,
    lot_number VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'recalled', 'quarantine', 'depleted')),

    -- Quantities
    initial_quantity NUMERIC(10, 2) NOT NULL,
    current_quantity NUMERIC(10, 2) NOT NULL,
    unit_of_measure VARCHAR(50) DEFAULT 'unit',

    -- Dates
    manufacturing_date DATE,
    expiration_date DATE,
    received_date DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Supplier Information
    purchase_order_id UUID REFERENCES purchase_orders(id) ON DELETE SET NULL,
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,
    supplier_batch_number VARCHAR(100),

    -- Cost Information
    unit_cost NUMERIC(15, 2),
    total_cost NUMERIC(15, 2),

    -- Quality Control
    quality_status VARCHAR(20) DEFAULT 'pending' CHECK (quality_status IN ('pending', 'passed', 'failed', 'quarantine')),
    quality_check_date DATE,
    quality_checked_by UUID REFERENCES users(id),
    quality_notes TEXT,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_batch_number_per_org UNIQUE (organization_id, batch_number),
    CONSTRAINT positive_quantities CHECK (
        initial_quantity > 0 AND
        current_quantity >= 0 AND
        current_quantity <= initial_quantity
    ),
    CONSTRAINT valid_dates CHECK (
        expiration_date IS NULL OR manufacturing_date IS NULL OR expiration_date > manufacturing_date
    )
);

-- Indexes
CREATE INDEX idx_product_batches_org_id ON product_batches(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_batches_product_id ON product_batches(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_batches_batch_number ON product_batches(batch_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_batches_status ON product_batches(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_batches_expiration ON product_batches(expiration_date) WHERE deleted_at IS NULL AND status = 'active';
CREATE INDEX idx_product_batches_location ON product_batches(location_id) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_product_batches_updated_at
    BEFORE UPDATE ON product_batches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE product_batches IS 'Batch/lot tracking for products with expiration dates and quality control';

-- ============================================================================
-- BATCH TRANSACTIONS
-- ============================================================================

-- Track all transactions against batches
CREATE TABLE IF NOT EXISTS batch_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    batch_id UUID NOT NULL REFERENCES product_batches(id) ON DELETE CASCADE,

    -- Transaction Information
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('sale', 'adjustment', 'return', 'waste', 'transfer', 'expiration')),
    quantity NUMERIC(10, 2) NOT NULL,
    balance_after NUMERIC(10, 2) NOT NULL,

    -- References
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    inventory_transfer_id UUID REFERENCES inventory_transfers(id) ON DELETE SET NULL,

    -- Transaction Details
    reason TEXT,
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT valid_balance CHECK (balance_after >= 0)
);

-- Indexes
CREATE INDEX idx_batch_transactions_org_id ON batch_transactions(organization_id);
CREATE INDEX idx_batch_transactions_batch_id ON batch_transactions(batch_id);
CREATE INDEX idx_batch_transactions_type ON batch_transactions(transaction_type);
CREATE INDEX idx_batch_transactions_date ON batch_transactions(transaction_date);

COMMENT ON TABLE batch_transactions IS 'Complete audit trail of all batch/lot transactions';

-- ============================================================================
-- CYCLE COUNTS
-- ============================================================================

-- Physical inventory cycle counting
CREATE TABLE IF NOT EXISTS cycle_counts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Count Information
    count_number VARCHAR(50) NOT NULL,
    count_date DATE NOT NULL DEFAULT CURRENT_DATE,
    count_type VARCHAR(20) DEFAULT 'cycle' CHECK (count_type IN ('cycle', 'full', 'spot', 'blind')),
    status VARCHAR(20) DEFAULT 'planned' CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled')),

    -- Scope
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    include_zero_stock BOOLEAN DEFAULT false,

    -- Statistics
    total_items_planned INTEGER DEFAULT 0,
    total_items_counted INTEGER DEFAULT 0,
    items_with_variance INTEGER DEFAULT 0,
    total_variance_value NUMERIC(15, 2) DEFAULT 0,

    -- Dates
    scheduled_date DATE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    counted_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_count_number_per_org UNIQUE (organization_id, count_number),
    CONSTRAINT positive_counts CHECK (
        total_items_planned >= 0 AND
        total_items_counted >= 0 AND
        items_with_variance >= 0
    )
);

-- Indexes
CREATE INDEX idx_cycle_counts_org_id ON cycle_counts(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cycle_counts_location_id ON cycle_counts(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cycle_counts_count_number ON cycle_counts(count_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_cycle_counts_status ON cycle_counts(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_cycle_counts_date ON cycle_counts(count_date) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_cycle_counts_updated_at
    BEFORE UPDATE ON cycle_counts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE cycle_counts IS 'Physical inventory cycle counts for stock verification';

-- ============================================================================
-- CYCLE COUNT ITEMS
-- ============================================================================

-- Individual items in a cycle count
CREATE TABLE IF NOT EXISTS cycle_count_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    cycle_count_id UUID NOT NULL REFERENCES cycle_counts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,

    -- Product Snapshot
    product_name VARCHAR(200) NOT NULL,
    product_sku VARCHAR(100),

    -- Quantities
    system_quantity NUMERIC(10, 2) NOT NULL, -- What system says
    counted_quantity NUMERIC(10, 2), -- What was physically counted
    variance_quantity NUMERIC(10, 2), -- Difference
    variance_percentage NUMERIC(5, 2), -- Percentage difference

    -- Valuation
    unit_cost NUMERIC(15, 2),
    variance_value NUMERIC(15, 2), -- Financial impact of variance

    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'counted', 'recounted', 'adjusted', 'approved')),

    -- Recount
    recount_required BOOLEAN DEFAULT false,
    recount_quantity NUMERIC(10, 2),
    recount_reason TEXT,

    -- Adjustment
    adjustment_applied BOOLEAN DEFAULT false,
    adjustment_date TIMESTAMP WITH TIME ZONE,
    adjustment_reason TEXT,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    counted_at TIMESTAMP WITH TIME ZONE,
    counted_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT non_negative_system_qty CHECK (system_quantity >= 0)
);

-- Indexes
CREATE INDEX idx_cycle_count_items_org_id ON cycle_count_items(organization_id);
CREATE INDEX idx_cycle_count_items_count_id ON cycle_count_items(cycle_count_id);
CREATE INDEX idx_cycle_count_items_product_id ON cycle_count_items(product_id);
CREATE INDEX idx_cycle_count_items_status ON cycle_count_items(status);

-- Auto-update trigger
CREATE TRIGGER update_cycle_count_items_updated_at
    BEFORE UPDATE ON cycle_count_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE cycle_count_items IS 'Individual product counts within a cycle count';

-- ============================================================================
-- STOCK ADJUSTMENT REASONS
-- ============================================================================

-- Pre-defined reasons for stock adjustments
CREATE TABLE IF NOT EXISTS stock_adjustment_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,

    -- Reason Information
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    reason_type VARCHAR(20) NOT NULL CHECK (reason_type IN ('increase', 'decrease', 'both')),

    -- System vs Custom
    is_system_reason BOOLEAN DEFAULT false,

    -- Settings
    is_active BOOLEAN DEFAULT true,
    requires_approval BOOLEAN DEFAULT false,
    requires_notes BOOLEAN DEFAULT true,

    -- Additional Information
    sort_order INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT unique_reason_code UNIQUE (organization_id, code)
);

-- Indexes
CREATE INDEX idx_stock_adjustment_reasons_org_id ON stock_adjustment_reasons(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_stock_adjustment_reasons_code ON stock_adjustment_reasons(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_stock_adjustment_reasons_active ON stock_adjustment_reasons(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_stock_adjustment_reasons_updated_at
    BEFORE UPDATE ON stock_adjustment_reasons
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE stock_adjustment_reasons IS 'Pre-defined reasons for inventory adjustments';

-- ============================================================================
-- INSERT DEFAULT STOCK ADJUSTMENT REASONS
-- ============================================================================

INSERT INTO stock_adjustment_reasons (id, organization_id, code, name, description, reason_type, is_system_reason, is_active)
VALUES
    (gen_random_uuid(), NULL, 'DAMAGE', 'Damaged Goods', 'Product was damaged and cannot be sold', 'decrease', true, true),
    (gen_random_uuid(), NULL, 'THEFT', 'Theft/Shrinkage', 'Product was lost due to theft or shrinkage', 'decrease', true, true),
    (gen_random_uuid(), NULL, 'EXPIRED', 'Expired Products', 'Product has expired and must be removed', 'decrease', true, true),
    (gen_random_uuid(), NULL, 'FOUND', 'Found Inventory', 'Product was found during count or audit', 'increase', true, true),
    (gen_random_uuid(), NULL, 'RETURN_SUPPLIER', 'Return to Supplier', 'Product returned to supplier', 'decrease', true, true),
    (gen_random_uuid(), NULL, 'SAMPLE', 'Product Sample', 'Used as sample or for demonstration', 'decrease', true, true),
    (gen_random_uuid(), NULL, 'PROMOTION', 'Promotional Giveaway', 'Given away as promotion', 'decrease', true, true),
    (gen_random_uuid(), NULL, 'COUNT_ERROR', 'Counting Error', 'System count was incorrect', 'both', true, true),
    (gen_random_uuid(), NULL, 'QUALITY_FAIL', 'Failed Quality Check', 'Product failed quality control', 'decrease', true, true),
    (gen_random_uuid(), NULL, 'RESTOCK', 'Restocking Adjustment', 'Adjustment during restocking', 'both', true, true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- HELPER FUNCTION: Calculate Cycle Count Variance
-- ============================================================================

CREATE OR REPLACE FUNCTION calculate_cycle_count_variance()
RETURNS TRIGGER AS $$
BEGIN
    -- Only calculate if counted_quantity is set
    IF NEW.counted_quantity IS NOT NULL THEN
        -- Calculate variance quantity
        NEW.variance_quantity := NEW.counted_quantity - NEW.system_quantity;

        -- Calculate variance percentage
        IF NEW.system_quantity > 0 THEN
            NEW.variance_percentage := (NEW.variance_quantity / NEW.system_quantity) * 100;
        ELSE
            NEW.variance_percentage := NULL;
        END IF;

        -- Calculate variance value
        IF NEW.unit_cost IS NOT NULL THEN
            NEW.variance_value := NEW.variance_quantity * NEW.unit_cost;
        END IF;

        -- Set counted timestamp if not set
        IF NEW.counted_at IS NULL THEN
            NEW.counted_at := CURRENT_TIMESTAMP;
        END IF;

        -- Update status to counted
        IF NEW.status = 'pending' THEN
            NEW.status := 'counted';
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for variance calculation
CREATE TRIGGER trigger_calculate_cycle_count_variance
    BEFORE UPDATE ON cycle_count_items
    FOR EACH ROW
    WHEN (OLD.counted_quantity IS DISTINCT FROM NEW.counted_quantity)
    EXECUTE FUNCTION calculate_cycle_count_variance();

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS on all new tables
ALTER TABLE product_serial_numbers ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_batches ENABLE ROW LEVEL SECURITY;
ALTER TABLE batch_transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE cycle_counts ENABLE ROW LEVEL SECURITY;
ALTER TABLE cycle_count_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE stock_adjustment_reasons ENABLE ROW LEVEL SECURITY;

-- ============================================================================
-- RLS POLICIES: PRODUCT_SERIAL_NUMBERS
-- ============================================================================

CREATE POLICY product_serial_numbers_super_admin_all
    ON product_serial_numbers FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY product_serial_numbers_select_own_org
    ON product_serial_numbers FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY product_serial_numbers_insert_own_org
    ON product_serial_numbers FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY product_serial_numbers_update_own_org
    ON product_serial_numbers FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY product_serial_numbers_delete_own_org
    ON product_serial_numbers FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: PRODUCT_BATCHES
-- ============================================================================

CREATE POLICY product_batches_super_admin_all
    ON product_batches FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY product_batches_select_own_org
    ON product_batches FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY product_batches_insert_own_org
    ON product_batches FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY product_batches_update_own_org
    ON product_batches FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY product_batches_delete_own_org
    ON product_batches FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: BATCH_TRANSACTIONS
-- ============================================================================

CREATE POLICY batch_transactions_super_admin_all
    ON batch_transactions FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY batch_transactions_select_own_org
    ON batch_transactions FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY batch_transactions_insert_own_org
    ON batch_transactions FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: CYCLE_COUNTS
-- ============================================================================

CREATE POLICY cycle_counts_super_admin_all
    ON cycle_counts FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY cycle_counts_select_own_org
    ON cycle_counts FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY cycle_counts_insert_own_org
    ON cycle_counts FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY cycle_counts_update_own_org
    ON cycle_counts FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY cycle_counts_delete_own_org
    ON cycle_counts FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: CYCLE_COUNT_ITEMS
-- ============================================================================

CREATE POLICY cycle_count_items_super_admin_all
    ON cycle_count_items FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY cycle_count_items_select_own_org
    ON cycle_count_items FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY cycle_count_items_insert_own_org
    ON cycle_count_items FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY cycle_count_items_update_own_org
    ON cycle_count_items FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY cycle_count_items_delete_own_org
    ON cycle_count_items FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: STOCK_ADJUSTMENT_REASONS
-- ============================================================================

CREATE POLICY stock_adjustment_reasons_super_admin_all
    ON stock_adjustment_reasons FOR ALL TO PUBLIC
    USING (is_super_admin());

-- System reasons (org_id = NULL) are visible to all
-- Organization-specific reasons are visible only to that organization
CREATE POLICY stock_adjustment_reasons_select
    ON stock_adjustment_reasons FOR SELECT TO PUBLIC
    USING (
        organization_id IS NULL OR
        organization_id = current_user_organization_id()
    );

CREATE POLICY stock_adjustment_reasons_insert_own_org
    ON stock_adjustment_reasons FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY stock_adjustment_reasons_update_own_org
    ON stock_adjustment_reasons FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY stock_adjustment_reasons_delete_own_org
    ON stock_adjustment_reasons FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- Migration completed successfully

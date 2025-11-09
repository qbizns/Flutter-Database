-- Migration V004: Additional Essential POS Tables
-- Created: 2025-11-09
-- Description: Adds suppliers, locations, product variants, promotions, expenses, and shifts tables
-- Dependencies: V001, V002, V003

-- ============================================================================
-- SUPPLIERS
-- ============================================================================

-- Supplier master table for managing vendors and suppliers
CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Basic Information
    supplier_code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    contact_person VARCHAR(200),
    email VARCHAR(255),
    phone VARCHAR(50),

    -- Address Information
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),

    -- Business Information
    tax_number VARCHAR(100),
    payment_terms VARCHAR(100), -- e.g., "Net 30", "Net 60", "Cash on Delivery"
    credit_limit NUMERIC(15, 2) DEFAULT 0,
    outstanding_balance NUMERIC(15, 2) DEFAULT 0,

    -- Statistics (cached for performance)
    total_purchases NUMERIC(15, 2) DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    last_order_date TIMESTAMP WITH TIME ZONE,

    -- Status and Settings
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_supplier_code_per_org UNIQUE (organization_id, supplier_code),
    CONSTRAINT positive_credit_limit CHECK (credit_limit >= 0)
);

-- Indexes for suppliers
CREATE INDEX idx_suppliers_organization_id ON suppliers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_supplier_code ON suppliers(supplier_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_name ON suppliers(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_status ON suppliers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_email ON suppliers(email) WHERE deleted_at IS NULL;

-- Auto-update trigger for suppliers
CREATE TRIGGER update_suppliers_updated_at
    BEFORE UPDATE ON suppliers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE suppliers IS 'Supplier master data for managing vendors and purchase sources';

-- ============================================================================
-- PURCHASE ORDERS
-- ============================================================================

-- Purchase orders for ordering inventory from suppliers
CREATE TABLE IF NOT EXISTS purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,

    -- Order Information
    po_number VARCHAR(50) NOT NULL,
    po_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expected_delivery_date DATE,
    actual_delivery_date DATE,

    -- Status
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'pending', 'approved', 'ordered', 'partial', 'received', 'cancelled')),

    -- Financial
    subtotal NUMERIC(15, 2) DEFAULT 0,
    tax_amount NUMERIC(15, 2) DEFAULT 0,
    discount_amount NUMERIC(15, 2) DEFAULT 0,
    shipping_cost NUMERIC(15, 2) DEFAULT 0,
    total_amount NUMERIC(15, 2) DEFAULT 0,

    -- Payment
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'partial', 'paid', 'overdue')),
    paid_amount NUMERIC(15, 2) DEFAULT 0,

    -- Additional Information
    notes TEXT,
    terms_and_conditions TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_po_number_per_org UNIQUE (organization_id, po_number),
    CONSTRAINT positive_amounts CHECK (
        subtotal >= 0 AND
        tax_amount >= 0 AND
        discount_amount >= 0 AND
        shipping_cost >= 0 AND
        total_amount >= 0 AND
        paid_amount >= 0
    )
);

-- Indexes for purchase_orders
CREATE INDEX idx_purchase_orders_organization_id ON purchase_orders(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_purchase_orders_supplier_id ON purchase_orders(supplier_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_purchase_orders_po_number ON purchase_orders(po_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_purchase_orders_po_date ON purchase_orders(po_date) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_purchase_orders_updated_at
    BEFORE UPDATE ON purchase_orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE purchase_orders IS 'Purchase orders for ordering inventory from suppliers';

-- ============================================================================
-- PURCHASE ORDER ITEMS
-- ============================================================================

-- Line items for purchase orders
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,

    -- Product snapshot (preserved even if product is deleted)
    product_name VARCHAR(200) NOT NULL,
    product_sku VARCHAR(100),

    -- Quantities
    quantity_ordered NUMERIC(10, 2) NOT NULL,
    quantity_received NUMERIC(10, 2) DEFAULT 0,
    unit_of_measure VARCHAR(50) DEFAULT 'unit',

    -- Pricing
    unit_cost NUMERIC(15, 2) NOT NULL,
    discount_percentage NUMERIC(5, 2) DEFAULT 0,
    discount_amount NUMERIC(15, 2) DEFAULT 0,
    tax_percentage NUMERIC(5, 2) DEFAULT 0,
    tax_amount NUMERIC(15, 2) DEFAULT 0,
    line_total NUMERIC(15, 2) NOT NULL,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT positive_quantities CHECK (
        quantity_ordered > 0 AND
        quantity_received >= 0 AND
        quantity_received <= quantity_ordered
    ),
    CONSTRAINT positive_pricing CHECK (
        unit_cost >= 0 AND
        discount_amount >= 0 AND
        tax_amount >= 0 AND
        line_total >= 0
    )
);

-- Indexes for purchase_order_items
CREATE INDEX idx_purchase_order_items_organization_id ON purchase_order_items(organization_id);
CREATE INDEX idx_purchase_order_items_po_id ON purchase_order_items(purchase_order_id);
CREATE INDEX idx_purchase_order_items_product_id ON purchase_order_items(product_id);

-- Auto-update trigger
CREATE TRIGGER update_purchase_order_items_updated_at
    BEFORE UPDATE ON purchase_order_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE purchase_order_items IS 'Line items for purchase orders';

-- ============================================================================
-- LOCATIONS / BRANCHES
-- ============================================================================

-- Store locations/branches for multi-location businesses
CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Basic Information
    location_code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    location_type VARCHAR(50) DEFAULT 'store' CHECK (location_type IN ('store', 'warehouse', 'headquarters', 'kiosk', 'online', 'other')),

    -- Contact Information
    phone VARCHAR(50),
    email VARCHAR(255),
    manager_user_id UUID REFERENCES users(id),

    -- Address Information
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    timezone VARCHAR(100) DEFAULT 'UTC',

    -- Business Hours (JSONB for flexibility)
    business_hours JSONB DEFAULT '{}', -- e.g., {"monday": {"open": "09:00", "close": "18:00"}, ...}

    -- Settings
    is_active BOOLEAN DEFAULT true,
    is_primary BOOLEAN DEFAULT false,
    allow_sales BOOLEAN DEFAULT true,
    allow_purchases BOOLEAN DEFAULT true,
    tax_rate NUMERIC(5, 2) DEFAULT 0,

    -- Additional Information
    notes TEXT,
    settings JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_location_code_per_org UNIQUE (organization_id, location_code)
);

-- Indexes for locations
CREATE INDEX idx_locations_organization_id ON locations(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_locations_location_code ON locations(location_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_locations_name ON locations(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_locations_is_active ON locations(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_locations_manager_user_id ON locations(manager_user_id) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_locations_updated_at
    BEFORE UPDATE ON locations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE locations IS 'Store locations/branches for multi-location businesses';

-- ============================================================================
-- PRODUCT VARIANTS
-- ============================================================================

-- Product variants for items with variations (size, color, etc.)
CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    -- Variant Information
    variant_name VARCHAR(200) NOT NULL,
    sku VARCHAR(100),
    barcode VARCHAR(100),

    -- Variation Attributes (JSONB for flexibility)
    attributes JSONB DEFAULT '{}', -- e.g., {"size": "Large", "color": "Red"}

    -- Pricing (can override parent product prices)
    cost_price NUMERIC(15, 2),
    selling_price NUMERIC(15, 2),
    compare_at_price NUMERIC(15, 2),

    -- Inventory
    current_stock NUMERIC(10, 2) DEFAULT 0,
    reorder_level NUMERIC(10, 2) DEFAULT 0,
    reorder_quantity NUMERIC(10, 2) DEFAULT 0,

    -- Physical Properties
    weight NUMERIC(10, 3),
    weight_unit VARCHAR(20),
    dimensions JSONB, -- e.g., {"length": 10, "width": 5, "height": 3, "unit": "cm"}

    -- Status
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,

    -- Display
    sort_order INTEGER DEFAULT 0,
    image_url TEXT,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_variant_sku_per_org UNIQUE (organization_id, sku),
    CONSTRAINT unique_variant_barcode_per_org UNIQUE (organization_id, barcode),
    CONSTRAINT positive_prices CHECK (
        (cost_price IS NULL OR cost_price >= 0) AND
        (selling_price IS NULL OR selling_price >= 0) AND
        (compare_at_price IS NULL OR compare_at_price >= 0)
    )
);

-- Indexes for product_variants
CREATE INDEX idx_product_variants_organization_id ON product_variants(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_variants_sku ON product_variants(sku) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_variants_barcode ON product_variants(barcode) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_variants_is_active ON product_variants(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_product_variants_updated_at
    BEFORE UPDATE ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE product_variants IS 'Product variants for items with variations (size, color, etc.)';

-- ============================================================================
-- PROMOTIONS / DISCOUNTS
-- ============================================================================

-- Promotions and discount campaigns
CREATE TABLE IF NOT EXISTS promotions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Basic Information
    promotion_code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,

    -- Type and Rules
    promotion_type VARCHAR(50) NOT NULL CHECK (promotion_type IN ('percentage', 'fixed_amount', 'buy_x_get_y', 'bundle', 'quantity_discount')),
    discount_value NUMERIC(15, 2) NOT NULL,

    -- Applicability
    applies_to VARCHAR(50) DEFAULT 'all' CHECK (applies_to IN ('all', 'specific_products', 'specific_categories', 'cart_total')),
    applicable_product_ids JSONB DEFAULT '[]', -- Array of product UUIDs
    applicable_category_ids JSONB DEFAULT '[]', -- Array of category UUIDs
    minimum_purchase_amount NUMERIC(15, 2) DEFAULT 0,
    minimum_quantity INTEGER DEFAULT 0,

    -- Buy X Get Y Rules (for that type)
    buy_quantity INTEGER,
    get_quantity INTEGER,
    get_discount_percentage NUMERIC(5, 2),

    -- Usage Limits
    max_uses_total INTEGER, -- NULL = unlimited
    max_uses_per_customer INTEGER, -- NULL = unlimited
    current_uses INTEGER DEFAULT 0,

    -- Date Range
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,

    -- Status
    is_active BOOLEAN DEFAULT true,
    is_combinable BOOLEAN DEFAULT false, -- Can be combined with other promotions

    -- Priority (lower number = higher priority)
    priority INTEGER DEFAULT 0,

    -- Additional Information
    terms_and_conditions TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_promotion_code_per_org UNIQUE (organization_id, promotion_code),
    CONSTRAINT valid_date_range CHECK (end_date IS NULL OR end_date > start_date),
    CONSTRAINT positive_discount CHECK (discount_value >= 0)
);

-- Indexes for promotions
CREATE INDEX idx_promotions_organization_id ON promotions(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_promotions_promotion_code ON promotions(promotion_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_promotions_is_active ON promotions(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_promotions_date_range ON promotions(start_date, end_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_promotions_promotion_type ON promotions(promotion_type) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_promotions_updated_at
    BEFORE UPDATE ON promotions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE promotions IS 'Promotions and discount campaigns';

-- ============================================================================
-- PROMOTION USAGE
-- ============================================================================

-- Track promotion usage per customer/transaction
CREATE TABLE IF NOT EXISTS promotion_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    promotion_id UUID NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,

    -- Usage Information
    discount_amount NUMERIC(15, 2) NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for promotion_usage
CREATE INDEX idx_promotion_usage_organization_id ON promotion_usage(organization_id);
CREATE INDEX idx_promotion_usage_promotion_id ON promotion_usage(promotion_id);
CREATE INDEX idx_promotion_usage_sale_id ON promotion_usage(sale_id);
CREATE INDEX idx_promotion_usage_customer_id ON promotion_usage(customer_id);
CREATE INDEX idx_promotion_usage_used_at ON promotion_usage(used_at);

COMMENT ON TABLE promotion_usage IS 'Track promotion usage per customer/transaction';

-- ============================================================================
-- EXPENSES
-- ============================================================================

-- Business expenses tracking
CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Basic Information
    expense_number VARCHAR(50) NOT NULL,
    expense_date DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Category
    category VARCHAR(100) NOT NULL, -- e.g., 'rent', 'utilities', 'salaries', 'supplies', 'marketing', 'other'
    subcategory VARCHAR(100),

    -- Payee Information
    payee_name VARCHAR(200) NOT NULL,
    payment_method VARCHAR(50) CHECK (payment_method IN ('cash', 'card', 'bank_transfer', 'check', 'other')),

    -- Amount
    amount NUMERIC(15, 2) NOT NULL,
    tax_amount NUMERIC(15, 2) DEFAULT 0,
    total_amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'paid', 'rejected', 'cancelled')),

    -- References
    reference_number VARCHAR(100), -- Invoice number, receipt number, etc.
    purchase_order_id UUID REFERENCES purchase_orders(id) ON DELETE SET NULL,

    -- Attachments
    receipt_url TEXT,
    attachment_urls JSONB DEFAULT '[]',

    -- Description and Notes
    description TEXT,
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_expense_number_per_org UNIQUE (organization_id, expense_number),
    CONSTRAINT positive_amounts CHECK (
        amount >= 0 AND
        tax_amount >= 0 AND
        total_amount >= 0
    )
);

-- Indexes for expenses
CREATE INDEX idx_expenses_organization_id ON expenses(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_location_id ON expenses(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_expense_number ON expenses(expense_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_expense_date ON expenses(expense_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_category ON expenses(category) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_status ON expenses(status) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_expenses_updated_at
    BEFORE UPDATE ON expenses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE expenses IS 'Business expenses tracking';

-- ============================================================================
-- SHIFTS
-- ============================================================================

-- Cashier shifts and cash register management
CREATE TABLE IF NOT EXISTS shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    -- Shift Information
    shift_number VARCHAR(50) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'closed', 'suspended')),

    -- Opening Cash
    opening_cash NUMERIC(15, 2) DEFAULT 0,
    opening_notes TEXT,

    -- Closing Cash
    expected_cash NUMERIC(15, 2),
    actual_cash NUMERIC(15, 2),
    cash_difference NUMERIC(15, 2), -- actual - expected
    closing_notes TEXT,

    -- Sales Summary (calculated)
    total_sales NUMERIC(15, 2) DEFAULT 0,
    total_transactions INTEGER DEFAULT 0,
    total_refunds NUMERIC(15, 2) DEFAULT 0,
    total_discounts NUMERIC(15, 2) DEFAULT 0,

    -- Payment Method Breakdown (JSONB for flexibility)
    payment_breakdown JSONB DEFAULT '{}', -- e.g., {"cash": 500.00, "card": 300.00, "mobile_money": 200.00}

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    closed_by UUID REFERENCES users(id),
    closed_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_shift_number_per_org UNIQUE (organization_id, shift_number),
    CONSTRAINT valid_times CHECK (end_time IS NULL OR end_time > start_time),
    CONSTRAINT positive_amounts CHECK (
        opening_cash >= 0 AND
        (expected_cash IS NULL OR expected_cash >= 0) AND
        (actual_cash IS NULL OR actual_cash >= 0)
    )
);

-- Indexes for shifts
CREATE INDEX idx_shifts_organization_id ON shifts(organization_id);
CREATE INDEX idx_shifts_location_id ON shifts(location_id);
CREATE INDEX idx_shifts_user_id ON shifts(user_id);
CREATE INDEX idx_shifts_shift_number ON shifts(shift_number);
CREATE INDEX idx_shifts_status ON shifts(status);
CREATE INDEX idx_shifts_start_time ON shifts(start_time);

-- Auto-update trigger
CREATE TRIGGER update_shifts_updated_at
    BEFORE UPDATE ON shifts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE shifts IS 'Cashier shifts and cash register management';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS on all new tables
ALTER TABLE suppliers ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE locations ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_variants ENABLE ROW LEVEL SECURITY;
ALTER TABLE promotions ENABLE ROW LEVEL SECURITY;
ALTER TABLE promotion_usage ENABLE ROW LEVEL SECURITY;
ALTER TABLE expenses ENABLE ROW LEVEL SECURITY;
ALTER TABLE shifts ENABLE ROW LEVEL SECURITY;

-- ============================================================================
-- RLS POLICIES: SUPPLIERS
-- ============================================================================

-- Super admin bypass
CREATE POLICY suppliers_super_admin_all
    ON suppliers FOR ALL TO PUBLIC
    USING (is_super_admin());

-- Organization-scoped policies
CREATE POLICY suppliers_select_own_org
    ON suppliers FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY suppliers_insert_own_org
    ON suppliers FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY suppliers_update_own_org
    ON suppliers FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY suppliers_delete_own_org
    ON suppliers FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: PURCHASE ORDERS
-- ============================================================================

CREATE POLICY purchase_orders_super_admin_all
    ON purchase_orders FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY purchase_orders_select_own_org
    ON purchase_orders FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY purchase_orders_insert_own_org
    ON purchase_orders FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY purchase_orders_update_own_org
    ON purchase_orders FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY purchase_orders_delete_own_org
    ON purchase_orders FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: PURCHASE ORDER ITEMS
-- ============================================================================

CREATE POLICY purchase_order_items_super_admin_all
    ON purchase_order_items FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY purchase_order_items_select_own_org
    ON purchase_order_items FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY purchase_order_items_insert_own_org
    ON purchase_order_items FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY purchase_order_items_update_own_org
    ON purchase_order_items FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY purchase_order_items_delete_own_org
    ON purchase_order_items FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: LOCATIONS
-- ============================================================================

CREATE POLICY locations_super_admin_all
    ON locations FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY locations_select_own_org
    ON locations FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY locations_insert_own_org
    ON locations FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY locations_update_own_org
    ON locations FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY locations_delete_own_org
    ON locations FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: PRODUCT VARIANTS
-- ============================================================================

CREATE POLICY product_variants_super_admin_all
    ON product_variants FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY product_variants_select_own_org
    ON product_variants FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY product_variants_insert_own_org
    ON product_variants FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY product_variants_update_own_org
    ON product_variants FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY product_variants_delete_own_org
    ON product_variants FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: PROMOTIONS
-- ============================================================================

CREATE POLICY promotions_super_admin_all
    ON promotions FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY promotions_select_own_org
    ON promotions FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY promotions_insert_own_org
    ON promotions FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY promotions_update_own_org
    ON promotions FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY promotions_delete_own_org
    ON promotions FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: PROMOTION USAGE
-- ============================================================================

CREATE POLICY promotion_usage_super_admin_all
    ON promotion_usage FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY promotion_usage_select_own_org
    ON promotion_usage FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY promotion_usage_insert_own_org
    ON promotion_usage FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY promotion_usage_update_own_org
    ON promotion_usage FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY promotion_usage_delete_own_org
    ON promotion_usage FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: EXPENSES
-- ============================================================================

CREATE POLICY expenses_super_admin_all
    ON expenses FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY expenses_select_own_org
    ON expenses FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY expenses_insert_own_org
    ON expenses FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY expenses_update_own_org
    ON expenses FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY expenses_delete_own_org
    ON expenses FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: SHIFTS
-- ============================================================================

CREATE POLICY shifts_super_admin_all
    ON shifts FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY shifts_select_own_org
    ON shifts FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY shifts_insert_own_org
    ON shifts FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY shifts_update_own_org
    ON shifts FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY shifts_delete_own_org
    ON shifts FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- FOREIGN KEY RELATIONSHIPS
-- ============================================================================

-- Add location references to existing tables
ALTER TABLE products ADD COLUMN IF NOT EXISTS location_id UUID REFERENCES locations(id) ON DELETE SET NULL;
ALTER TABLE sales ADD COLUMN IF NOT EXISTS location_id UUID REFERENCES locations(id) ON DELETE SET NULL;
ALTER TABLE sales ADD COLUMN IF NOT EXISTS shift_id UUID REFERENCES shifts(id) ON DELETE SET NULL;
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS product_variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL;

-- Create indexes for new foreign keys
CREATE INDEX IF NOT EXISTS idx_products_location_id ON products(location_id);
CREATE INDEX IF NOT EXISTS idx_sales_location_id ON sales(location_id);
CREATE INDEX IF NOT EXISTS idx_sales_shift_id ON sales(shift_id);
CREATE INDEX IF NOT EXISTS idx_sale_items_product_variant_id ON sale_items(product_variant_id);

-- Migration completed successfully

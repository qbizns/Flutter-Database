-- =====================================================
-- POS Migration V013: Advanced POS Features
-- Description: Cash management, advanced pricing, gift cards, returns,
--              procurement, monitoring, and multi-channel support
-- =====================================================

\echo 'Creating advanced POS features...';

-- =====================================================
-- MODULE 1: MONEY HANDLING (Cash & Sessions)
-- =====================================================

\echo 'Creating cash management module...';

-- POS Sessions
CREATE TABLE IF NOT EXISTS pos_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Session identification
    session_number VARCHAR(50) NOT NULL,
    session_name VARCHAR(255),

    -- Links
    device_id UUID REFERENCES devices(id),
    location_id UUID NOT NULL REFERENCES locations(id),
    user_id UUID NOT NULL REFERENCES users(id), -- Primary cashier
    shift_id UUID REFERENCES shifts(id),

    -- Session timing
    opened_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP,

    -- Opening float by payment method
    opening_cash NUMERIC(20, 4) DEFAULT 0,
    opening_card NUMERIC(20, 4) DEFAULT 0,
    opening_other NUMERIC(20, 4) DEFAULT 0,

    -- Expected closing (system calculated)
    expected_cash NUMERIC(20, 4) DEFAULT 0,
    expected_card NUMERIC(20, 4) DEFAULT 0,
    expected_other NUMERIC(20, 4) DEFAULT 0,

    -- Actual counted closing
    counted_cash NUMERIC(20, 4),
    counted_card NUMERIC(20, 4),
    counted_other NUMERIC(20, 4),

    -- Differences (counted - expected)
    difference_cash NUMERIC(20, 4) DEFAULT 0,
    difference_card NUMERIC(20, 4) DEFAULT 0,
    difference_other NUMERIC(20, 4) DEFAULT 0,

    -- Status
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'closing', 'closed', 'reconciled')),

    -- Z-report reference
    z_report_number VARCHAR(50),

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_session_number_org UNIQUE(organization_id, session_number)
);

CREATE INDEX idx_pos_sessions_org ON pos_sessions(organization_id, opened_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_pos_sessions_device ON pos_sessions(device_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pos_sessions_user ON pos_sessions(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pos_sessions_status ON pos_sessions(status) WHERE deleted_at IS NULL AND status = 'open';

COMMENT ON TABLE pos_sessions IS 'POS cash register sessions (opening/closing/reconciliation)';

-- Enable RLS
ALTER TABLE pos_sessions ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_pos_sessions ON pos_sessions FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Cash Drawers
CREATE TABLE IF NOT EXISTS cash_drawers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Drawer identification
    drawer_code VARCHAR(50) NOT NULL,
    drawer_name VARCHAR(255) NOT NULL,

    -- Physical location
    location_id UUID NOT NULL REFERENCES locations(id),
    device_id UUID REFERENCES devices(id), -- Optional: dedicated device

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_drawer_code_org UNIQUE(organization_id, drawer_code)
);

CREATE INDEX idx_cash_drawers_org ON cash_drawers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cash_drawers_location ON cash_drawers(location_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE cash_drawers IS 'Physical cash drawers at locations';

-- Enable RLS
ALTER TABLE cash_drawers ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_cash_drawers ON cash_drawers FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Cash Drawer Sessions (link drawer to POS session)
CREATE TABLE IF NOT EXISTS cash_drawer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    cash_drawer_id UUID NOT NULL REFERENCES cash_drawers(id),
    pos_session_id UUID NOT NULL REFERENCES pos_sessions(id),

    -- Opening/closing amounts
    opening_amount NUMERIC(20, 4) DEFAULT 0,
    closing_amount NUMERIC(20, 4),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_drawer_session UNIQUE(cash_drawer_id, pos_session_id)
);

CREATE INDEX idx_cash_drawer_sessions_drawer ON cash_drawer_sessions(cash_drawer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cash_drawer_sessions_session ON cash_drawer_sessions(pos_session_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE cash_drawer_sessions IS 'Links cash drawers to POS sessions';

-- Enable RLS
ALTER TABLE cash_drawer_sessions ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_cash_drawer_sessions ON cash_drawer_sessions FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Cash Movements (pay-ins, pay-outs)
CREATE TABLE IF NOT EXISTS cash_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Link to session
    pos_session_id UUID NOT NULL REFERENCES pos_sessions(id),
    cash_drawer_id UUID REFERENCES cash_drawers(id),

    -- Movement details
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('pay_in', 'pay_out', 'float_add', 'drop_to_safe')),
    amount NUMERIC(20, 4) NOT NULL,

    -- Reason
    reason_code VARCHAR(50),
    reason_description TEXT NOT NULL,

    -- Who did it
    user_id UUID NOT NULL REFERENCES users(id),

    -- Approval (optional)
    requires_approval BOOLEAN DEFAULT false,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT chk_movement_amount CHECK (amount != 0)
);

CREATE INDEX idx_cash_movements_session ON cash_movements(pos_session_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cash_movements_type ON cash_movements(movement_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_cash_movements_date ON cash_movements(organization_id, created_at DESC) WHERE deleted_at IS NULL;

COMMENT ON TABLE cash_movements IS 'Non-sale cash movements (pay-ins, pay-outs, drops, float additions)';

-- Enable RLS
ALTER TABLE cash_movements ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_cash_movements ON cash_movements FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Add session_id to sales table
ALTER TABLE sales ADD COLUMN IF NOT EXISTS pos_session_id UUID REFERENCES pos_sessions(id);
CREATE INDEX IF NOT EXISTS idx_sales_pos_session ON sales(pos_session_id) WHERE deleted_at IS NULL;

-- =====================================================
-- MODULE 2: ADVANCED PRICING & CATALOG
-- =====================================================

\echo 'Creating advanced pricing module...';

-- Price Lists
CREATE TABLE IF NOT EXISTS price_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Price list identification
    price_list_code VARCHAR(50) NOT NULL,
    price_list_name VARCHAR(255) NOT NULL,

    -- Type
    price_list_type VARCHAR(30) DEFAULT 'standard' CHECK (
        price_list_type IN ('standard', 'retail', 'wholesale', 'vip', 'seasonal', 'location', 'customer_group')
    ),

    -- Effective dates
    effective_from DATE,
    effective_to DATE,

    -- Base configuration
    base_price_adjustment_type VARCHAR(20) CHECK (base_price_adjustment_type IN ('percentage', 'fixed', 'none')),
    base_price_adjustment_value NUMERIC(20, 4) DEFAULT 0,

    -- Priority (lower = higher priority)
    priority INTEGER DEFAULT 10,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    description TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_price_list_code_org UNIQUE(organization_id, price_list_code)
);

CREATE INDEX idx_price_lists_org ON price_lists(organization_id, priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_price_lists_dates ON price_lists(effective_from, effective_to) WHERE deleted_at IS NULL;

COMMENT ON TABLE price_lists IS 'Price lists for tiered/special pricing (retail, wholesale, VIP, etc.)';

-- Enable RLS
ALTER TABLE price_lists ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_price_lists ON price_lists FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Price List Items
CREATE TABLE IF NOT EXISTS price_list_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price_list_id UUID NOT NULL REFERENCES price_lists(id) ON DELETE CASCADE,

    -- Product reference
    product_id UUID REFERENCES products(id),
    product_variant_id UUID REFERENCES product_variants(id),
    category_id UUID REFERENCES categories(id), -- Optional: apply to entire category

    -- Price override
    override_price NUMERIC(20, 4),

    -- Or percentage adjustment
    discount_percentage NUMERIC(10, 4),
    markup_percentage NUMERIC(10, 4),

    -- Min/max constraints
    min_price NUMERIC(20, 4),
    max_price NUMERIC(20, 4),
    min_quantity NUMERIC(20, 4),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT chk_price_list_item_target CHECK (
        product_id IS NOT NULL OR product_variant_id IS NOT NULL OR category_id IS NOT NULL
    )
);

CREATE INDEX idx_price_list_items_list ON price_list_items(price_list_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_price_list_items_product ON price_list_items(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_price_list_items_variant ON price_list_items(product_variant_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE price_list_items IS 'Individual product/variant price overrides within a price list';

-- Product Components (for composite/bundle products)
CREATE TABLE IF NOT EXISTS product_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Parent composite product
    parent_product_id UUID NOT NULL REFERENCES products(id),

    -- Component (can be product or variant)
    component_product_id UUID REFERENCES products(id),
    component_variant_id UUID REFERENCES product_variants(id),

    -- Quantity required
    quantity NUMERIC(20, 4) NOT NULL DEFAULT 1,

    -- Pricing behavior
    inherit_price BOOLEAN DEFAULT false, -- If true, component price is visible on receipt
    price_override NUMERIC(20, 4), -- Optional fixed price for this component in bundle

    -- Display
    display_order INTEGER DEFAULT 0,
    is_optional BOOLEAN DEFAULT false,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT chk_component_target CHECK (
        component_product_id IS NOT NULL OR component_variant_id IS NOT NULL
    )
);

CREATE INDEX idx_product_components_parent ON product_components(parent_product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_components_component ON product_components(component_product_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE product_components IS 'Defines what composite/bundle products consist of';

-- Enable RLS
ALTER TABLE product_components ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_product_components ON product_components FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Units of Measure
CREATE TABLE IF NOT EXISTS units_of_measure (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE, -- NULL = global

    -- UoM identification
    uom_code VARCHAR(20) NOT NULL,
    uom_name VARCHAR(100) NOT NULL,
    uom_type VARCHAR(30) NOT NULL CHECK (uom_type IN ('unit', 'weight', 'volume', 'length', 'time')),

    -- Is this a base unit?
    is_base_unit BOOLEAN DEFAULT false,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_uom_code UNIQUE(organization_id, uom_code)
);

CREATE INDEX idx_uom_org ON units_of_measure(organization_id, uom_type) WHERE deleted_at IS NULL;

COMMENT ON TABLE units_of_measure IS 'Units of measure (kg, g, lb, piece, liter, etc.)';

-- Insert common UoMs
INSERT INTO units_of_measure (uom_code, uom_name, uom_type, is_base_unit) VALUES
    ('piece', 'Piece', 'unit', true),
    ('dozen', 'Dozen', 'unit', false),
    ('kg', 'Kilogram', 'weight', true),
    ('g', 'Gram', 'weight', false),
    ('lb', 'Pound', 'weight', false),
    ('oz', 'Ounce', 'weight', false),
    ('l', 'Liter', 'volume', true),
    ('ml', 'Milliliter', 'volume', false),
    ('gal', 'Gallon', 'volume', false),
    ('m', 'Meter', 'length', true),
    ('cm', 'Centimeter', 'length', false),
    ('ft', 'Foot', 'length', false)
ON CONFLICT (organization_id, uom_code) DO NOTHING;

-- UoM Conversions
CREATE TABLE IF NOT EXISTS uom_conversions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    from_uom_id UUID NOT NULL REFERENCES units_of_measure(id),
    to_uom_id UUID NOT NULL REFERENCES units_of_measure(id),

    -- Conversion factor (from * factor = to)
    conversion_factor NUMERIC(20, 10) NOT NULL,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_uom_conversion UNIQUE(from_uom_id, to_uom_id),
    CONSTRAINT chk_conversion_factor CHECK (conversion_factor > 0)
);

CREATE INDEX idx_uom_conversions_from ON uom_conversions(from_uom_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE uom_conversions IS 'Unit of measure conversion factors (e.g., 1 dozen = 12 pieces)';

-- Insert common conversions
INSERT INTO uom_conversions (from_uom_id, to_uom_id, conversion_factor)
SELECT f.id, t.id, cf FROM (VALUES
    ('dozen', 'piece', 12),
    ('kg', 'g', 1000),
    ('lb', 'oz', 16),
    ('lb', 'g', 453.592),
    ('l', 'ml', 1000),
    ('gal', 'l', 3.78541),
    ('m', 'cm', 100),
    ('ft', 'cm', 30.48)
) AS conversions(from_code, to_code, cf)
JOIN units_of_measure f ON f.uom_code = conversions.from_code AND f.organization_id IS NULL
JOIN units_of_measure t ON t.uom_code = conversions.to_code AND t.organization_id IS NULL
ON CONFLICT (from_uom_id, to_uom_id) DO NOTHING;

-- Add price_list_id and uom references to existing tables
ALTER TABLE customers ADD COLUMN IF NOT EXISTS price_list_id UUID REFERENCES price_lists(id);
ALTER TABLE locations ADD COLUMN IF NOT EXISTS default_price_list_id UUID REFERENCES price_lists(id);
ALTER TABLE products ADD COLUMN IF NOT EXISTS uom_id UUID REFERENCES units_of_measure(id);
ALTER TABLE sales ADD COLUMN IF NOT EXISTS price_list_id UUID REFERENCES price_lists(id);

-- =====================================================
-- MODULE 3: GIFT CARDS & STORE CREDIT
-- =====================================================

\echo 'Creating gift cards module...';

-- Gift Cards
CREATE TABLE IF NOT EXISTS gift_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Card identification
    card_number VARCHAR(100) NOT NULL,
    pin_code VARCHAR(20), -- Optional PIN for security

    -- Ownership
    customer_id UUID REFERENCES customers(id),

    -- Balances
    original_value NUMERIC(20, 4) NOT NULL,
    current_balance NUMERIC(20, 4) NOT NULL,

    -- Effective dates
    issued_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expiry_date DATE,

    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (
        status IN ('active', 'inactive', 'blocked', 'expired', 'fully_redeemed')
    ),

    -- Issued info
    issued_by_user_id UUID REFERENCES users(id),
    issued_location_id UUID REFERENCES locations(id),

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_gift_card_number_org UNIQUE(organization_id, card_number),
    CONSTRAINT chk_gift_card_balance CHECK (current_balance >= 0 AND current_balance <= original_value)
);

CREATE INDEX idx_gift_cards_org ON gift_cards(organization_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_gift_cards_number ON gift_cards(card_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_gift_cards_customer ON gift_cards(customer_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE gift_cards IS 'Gift card accounts with balance tracking';

-- Enable RLS
ALTER TABLE gift_cards ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_gift_cards ON gift_cards FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Gift Card Transactions
CREATE TABLE IF NOT EXISTS gift_card_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    gift_card_id UUID NOT NULL REFERENCES gift_cards(id),

    -- Transaction details
    transaction_type VARCHAR(20) NOT NULL CHECK (
        transaction_type IN ('issue', 'load', 'redemption', 'refund', 'adjustment', 'expiry')
    ),
    amount NUMERIC(20, 4) NOT NULL,
    balance_after NUMERIC(20, 4) NOT NULL,

    -- Links
    sale_id UUID REFERENCES sales(id),
    payment_id UUID REFERENCES payments(id),

    -- Who did it
    user_id UUID REFERENCES users(id),
    location_id UUID REFERENCES locations(id),

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_gift_card_txns_card ON gift_card_transactions(gift_card_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_gift_card_txns_sale ON gift_card_transactions(sale_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE gift_card_transactions IS 'All gift card loads, redemptions, and adjustments';

-- Enable RLS
ALTER TABLE gift_card_transactions ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_gift_card_txns ON gift_card_transactions FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Customer Store Credit Accounts
CREATE TABLE IF NOT EXISTS customer_store_credit_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    customer_id UUID NOT NULL REFERENCES customers(id),

    -- Balance
    current_balance NUMERIC(20, 4) DEFAULT 0,

    -- Limits
    credit_limit NUMERIC(20, 4),

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_store_credit_customer UNIQUE(organization_id, customer_id),
    CONSTRAINT chk_store_credit_balance CHECK (current_balance >= 0)
);

CREATE INDEX idx_store_credit_org ON customer_store_credit_accounts(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_store_credit_customer ON customer_store_credit_accounts(customer_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE customer_store_credit_accounts IS 'Customer store credit balances';

-- Enable RLS
ALTER TABLE customer_store_credit_accounts ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_store_credit ON customer_store_credit_accounts FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Store Credit Transactions
CREATE TABLE IF NOT EXISTS store_credit_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    store_credit_account_id UUID NOT NULL REFERENCES customer_store_credit_accounts(id),

    -- Transaction details
    transaction_type VARCHAR(20) NOT NULL CHECK (
        transaction_type IN ('issue', 'redemption', 'refund', 'adjustment', 'expiry')
    ),
    amount NUMERIC(20, 4) NOT NULL,
    balance_after NUMERIC(20, 4) NOT NULL,

    -- Links
    sale_id UUID REFERENCES sales(id),
    payment_id UUID REFERENCES payments(id),

    -- Who did it
    user_id UUID REFERENCES users(id),
    location_id UUID REFERENCES locations(id),

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_store_credit_txns_account ON store_credit_transactions(store_credit_account_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_store_credit_txns_sale ON store_credit_transactions(sale_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE store_credit_transactions IS 'All store credit increases/decreases';

-- Enable RLS
ALTER TABLE store_credit_transactions ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_store_credit_txns ON store_credit_transactions FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- =====================================================
-- MODULE 4: RETURNS & AFTER-SALES
-- =====================================================

\echo 'Creating returns module...';

-- Return Reasons
CREATE TABLE IF NOT EXISTS return_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE, -- NULL = global

    -- Reason identification
    reason_code VARCHAR(50) NOT NULL,
    reason_name VARCHAR(255) NOT NULL,

    -- Configuration
    requires_approval BOOLEAN DEFAULT false,
    affects_inventory BOOLEAN DEFAULT true,
    is_restockable BOOLEAN DEFAULT true,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_return_reason_code UNIQUE(organization_id, reason_code)
);

CREATE INDEX idx_return_reasons_org ON return_reasons(organization_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE return_reasons IS 'Standardized return reasons (defective, wrong size, changed mind, etc.)';

-- Insert common return reasons
INSERT INTO return_reasons (reason_code, reason_name, requires_approval, is_restockable) VALUES
    ('DEFECTIVE', 'Defective/Damaged', false, false),
    ('WRONG_SIZE', 'Wrong Size/Fit', false, true),
    ('WRONG_ITEM', 'Wrong Item Ordered', false, true),
    ('CHANGED_MIND', 'Changed Mind', false, true),
    ('BETTER_PRICE', 'Found Better Price', true, true),
    ('NOT_AS_DESC', 'Not as Described', false, true),
    ('ARRIVED_LATE', 'Arrived Too Late', false, true),
    ('DUPLICATE', 'Duplicate Order', false, true),
    ('WARRANTY', 'Warranty Return', true, false),
    ('OTHER', 'Other Reason', false, true)
ON CONFLICT (organization_id, reason_code) DO NOTHING;

-- Sale Returns (header)
CREATE TABLE IF NOT EXISTS sale_returns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Return identification
    return_number VARCHAR(50) NOT NULL,

    -- Link to original sale (optional)
    original_sale_id UUID REFERENCES sales(id),

    -- Customer
    customer_id UUID REFERENCES customers(id),

    -- Location and staff
    location_id UUID NOT NULL REFERENCES locations(id),
    user_id UUID NOT NULL REFERENCES users(id),

    -- Return date
    return_date DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Financial impact
    total_amount NUMERIC(20, 4) DEFAULT 0,
    refund_amount NUMERIC(20, 4) DEFAULT 0,
    restocking_fee NUMERIC(20, 4) DEFAULT 0,

    -- Refund method
    refund_method VARCHAR(30) CHECK (refund_method IN ('original_payment', 'cash', 'store_credit', 'gift_card', 'exchange')),

    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (
        status IN ('pending', 'approved', 'refunded', 'rejected', 'completed')
    ),

    -- Approval
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_return_number_org UNIQUE(organization_id, return_number)
);

CREATE INDEX idx_sale_returns_org ON sale_returns(organization_id, return_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_sale_returns_original_sale ON sale_returns(original_sale_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sale_returns_customer ON sale_returns(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sale_returns_status ON sale_returns(status) WHERE deleted_at IS NULL;

COMMENT ON TABLE sale_returns IS 'Return headers linking to original sales';

-- Enable RLS
ALTER TABLE sale_returns ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_sale_returns ON sale_returns FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Sale Return Items
CREATE TABLE IF NOT EXISTS sale_return_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    sale_return_id UUID NOT NULL REFERENCES sale_returns(id) ON DELETE CASCADE,

    -- Link to original item (optional)
    original_sale_item_id UUID REFERENCES sale_items(id),

    -- Product details
    product_id UUID NOT NULL REFERENCES products(id),
    product_variant_id UUID REFERENCES product_variants(id),

    -- Quantity and pricing
    quantity NUMERIC(20, 4) NOT NULL,
    unit_price NUMERIC(20, 4) NOT NULL,
    subtotal NUMERIC(20, 4) NOT NULL,
    tax_amount NUMERIC(20, 4) DEFAULT 0,
    discount_amount NUMERIC(20, 4) DEFAULT 0,
    total_amount NUMERIC(20, 4) NOT NULL,

    -- Return reason
    return_reason_id UUID REFERENCES return_reasons(id),
    return_reason_notes TEXT,

    -- Condition
    item_condition VARCHAR(30) DEFAULT 'resellable' CHECK (
        item_condition IN ('resellable', 'damaged', 'defective', 'used', 'opened')
    ),
    is_restockable BOOLEAN DEFAULT true,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_sale_return_items_return ON sale_return_items(sale_return_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sale_return_items_product ON sale_return_items(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sale_return_items_original ON sale_return_items(original_sale_item_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE sale_return_items IS 'Individual items being returned with condition and reason';

-- Enable RLS
ALTER TABLE sale_return_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_sale_return_items ON sale_return_items FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- =====================================================
-- MODULE 5: PROCUREMENT (Purchase Orders)
-- =====================================================

\echo 'Creating procurement module...';

-- Purchase Orders (header)
CREATE TABLE IF NOT EXISTS purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- PO identification
    po_number VARCHAR(50) NOT NULL,

    -- Supplier
    supplier_id UUID NOT NULL REFERENCES suppliers(id),

    -- Destination
    location_id UUID NOT NULL REFERENCES locations(id),

    -- Dates
    order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expected_delivery_date DATE,
    actual_delivery_date DATE,

    -- Financial
    subtotal_amount NUMERIC(20, 4) DEFAULT 0,
    tax_amount NUMERIC(20, 4) DEFAULT 0,
    shipping_amount NUMERIC(20, 4) DEFAULT 0,
    total_amount NUMERIC(20, 4) DEFAULT 0,

    -- Payment terms
    payment_terms VARCHAR(100),
    payment_due_date DATE,

    -- Status
    status VARCHAR(20) DEFAULT 'draft' CHECK (
        status IN ('draft', 'sent', 'confirmed', 'partially_received', 'received', 'cancelled')
    ),

    -- Approval
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_po_number_org UNIQUE(organization_id, po_number)
);

CREATE INDEX idx_purchase_orders_org ON purchase_orders(organization_id, order_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_purchase_orders_supplier ON purchase_orders(supplier_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status) WHERE deleted_at IS NULL;

COMMENT ON TABLE purchase_orders IS 'Purchase order headers for inventory procurement';

-- Enable RLS
ALTER TABLE purchase_orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_purchase_orders ON purchase_orders FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Purchase Order Items
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,

    -- Product
    product_id UUID NOT NULL REFERENCES products(id),
    product_variant_id UUID REFERENCES product_variants(id),

    -- Quantity
    quantity_ordered NUMERIC(20, 4) NOT NULL,
    quantity_received NUMERIC(20, 4) DEFAULT 0,
    quantity_remaining NUMERIC(20, 4) GENERATED ALWAYS AS (quantity_ordered - quantity_received) STORED,

    -- Pricing
    unit_cost NUMERIC(20, 4) NOT NULL,
    subtotal NUMERIC(20, 4) NOT NULL,
    tax_amount NUMERIC(20, 4) DEFAULT 0,
    total_amount NUMERIC(20, 4) NOT NULL,

    -- Expected dates
    expected_delivery_date DATE,

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_po_line UNIQUE(purchase_order_id, line_number)
);

CREATE INDEX idx_po_items_po ON purchase_order_items(purchase_order_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_po_items_product ON purchase_order_items(product_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE purchase_order_items IS 'Line items within purchase orders';

-- Enable RLS
ALTER TABLE purchase_order_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_po_items ON purchase_order_items FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Goods Receipts (receiving documents)
CREATE TABLE IF NOT EXISTS goods_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Receipt identification
    receipt_number VARCHAR(50) NOT NULL,

    -- Link to PO
    purchase_order_id UUID REFERENCES purchase_orders(id),

    -- Supplier
    supplier_id UUID NOT NULL REFERENCES suppliers(id),

    -- Location
    location_id UUID NOT NULL REFERENCES locations(id),

    -- Receipt date
    receipt_date DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Receiver
    received_by UUID NOT NULL REFERENCES users(id),

    -- Status
    status VARCHAR(20) DEFAULT 'draft' CHECK (
        status IN ('draft', 'received', 'inspected', 'completed')
    ),

    -- Metadata
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_receipt_number_org UNIQUE(organization_id, receipt_number)
);

CREATE INDEX idx_goods_receipts_org ON goods_receipts(organization_id, receipt_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_goods_receipts_po ON goods_receipts(purchase_order_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_goods_receipts_supplier ON goods_receipts(supplier_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE goods_receipts IS 'Goods receiving documents';

-- Enable RLS
ALTER TABLE goods_receipts ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_goods_receipts ON goods_receipts FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Goods Receipt Items
CREATE TABLE IF NOT EXISTS goods_receipt_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    goods_receipt_id UUID NOT NULL REFERENCES goods_receipts(id) ON DELETE CASCADE,
    purchase_order_item_id UUID REFERENCES purchase_order_items(id),

    -- Product
    product_id UUID NOT NULL REFERENCES products(id),
    product_variant_id UUID REFERENCES product_variants(id),

    -- Quantity received
    quantity_received NUMERIC(20, 4) NOT NULL,

    -- Quality check
    quantity_accepted NUMERIC(20, 4),
    quantity_rejected NUMERIC(20, 4),
    rejection_reason TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_gr_items_receipt ON goods_receipt_items(goods_receipt_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_gr_items_po_item ON goods_receipt_items(purchase_order_item_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE goods_receipt_items IS 'Items received in goods receipts';

-- Enable RLS
ALTER TABLE goods_receipt_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_gr_items ON goods_receipt_items FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- =====================================================
-- MODULE 6: MONITORING & SYSTEM HEALTH
-- =====================================================

\echo 'Creating monitoring module...';

-- System Health
CREATE TABLE IF NOT EXISTS system_health (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE, -- NULL = global

    -- Health check type
    check_type VARCHAR(50) NOT NULL,
    check_name VARCHAR(255) NOT NULL,

    -- Status
    status VARCHAR(20) DEFAULT 'unknown' CHECK (status IN ('healthy', 'degraded', 'unhealthy', 'unknown')),

    -- Metrics
    last_check_at TIMESTAMP,
    last_success_at TIMESTAMP,
    last_failure_at TIMESTAMP,

    -- Details
    metric_value NUMERIC(20, 4),
    metric_unit VARCHAR(50),
    threshold_warning NUMERIC(20, 4),
    threshold_critical NUMERIC(20, 4),

    -- Metadata
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_system_health_check UNIQUE(organization_id, check_type, check_name)
);

CREATE INDEX idx_system_health_org ON system_health(organization_id, status) WHERE status != 'healthy';
CREATE INDEX idx_system_health_type ON system_health(check_type, last_check_at DESC);

COMMENT ON TABLE system_health IS 'System health monitoring (DB, backups, queues, services)';

-- POS Error Logs
CREATE TABLE IF NOT EXISTS pos_error_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,

    -- Error classification
    error_level VARCHAR(20) NOT NULL CHECK (error_level IN ('debug', 'info', 'warning', 'error', 'critical')),
    error_code VARCHAR(50),
    error_message TEXT NOT NULL,

    -- Context
    device_id UUID REFERENCES devices(id),
    user_id UUID REFERENCES users(id),
    pos_session_id UUID REFERENCES pos_sessions(id),
    sale_id UUID REFERENCES sales(id),

    -- Stack trace and details
    stack_trace TEXT,
    request_data JSONB,
    error_data JSONB,

    -- Resolution
    is_resolved BOOLEAN DEFAULT false,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP,
    resolution_notes TEXT,

    -- Metadata
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pos_errors_org ON pos_error_logs(organization_id, occurred_at DESC);
CREATE INDEX idx_pos_errors_level ON pos_error_logs(error_level, is_resolved);
CREATE INDEX idx_pos_errors_device ON pos_error_logs(device_id) WHERE device_id IS NOT NULL;
CREATE INDEX idx_pos_errors_user ON pos_error_logs(user_id) WHERE user_id IS NOT NULL;

COMMENT ON TABLE pos_error_logs IS 'Application-level error logs for debugging';

-- =====================================================
-- MODULE 7: MULTI-CHANNEL / INTEGRATIONS
-- =====================================================

\echo 'Creating multi-channel module...';

-- Sales Channels
CREATE TABLE IF NOT EXISTS sales_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Channel identification
    channel_code VARCHAR(50) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    channel_type VARCHAR(30) NOT NULL CHECK (
        channel_type IN ('in_store', 'web', 'mobile_app', 'marketplace', 'phone', 'social', 'partner')
    ),

    -- Configuration
    is_active BOOLEAN DEFAULT true,
    sync_inventory BOOLEAN DEFAULT true,
    sync_customers BOOLEAN DEFAULT true,

    -- Integration
    external_system_name VARCHAR(100),
    api_endpoint VARCHAR(500),

    -- Metadata
    settings JSONB,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_channel_code_org UNIQUE(organization_id, channel_code)
);

CREATE INDEX idx_sales_channels_org ON sales_channels(organization_id, is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE sales_channels IS 'Sales channels (in-store, web, mobile, marketplaces)';

-- Enable RLS
ALTER TABLE sales_channels ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_sales_channels ON sales_channels FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Insert default channels
DO $$
DECLARE
    v_org_id UUID;
BEGIN
    FOR v_org_id IN SELECT id FROM organizations LIMIT 10 LOOP
        INSERT INTO sales_channels (organization_id, channel_code, channel_name, channel_type)
        VALUES
            (v_org_id, 'IN_STORE', 'In-Store', 'in_store'),
            (v_org_id, 'WEB', 'Web Store', 'web'),
            (v_org_id, 'MOBILE', 'Mobile App', 'mobile_app')
        ON CONFLICT (organization_id, channel_code) DO NOTHING;
    END LOOP;
END $$;

-- External Order Mappings
CREATE TABLE IF NOT EXISTS external_order_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Internal order
    sale_id UUID NOT NULL REFERENCES sales(id),

    -- External system
    sales_channel_id UUID NOT NULL REFERENCES sales_channels(id),
    external_order_id VARCHAR(255) NOT NULL,
    external_order_number VARCHAR(100),

    -- Sync status
    sync_status VARCHAR(20) DEFAULT 'pending' CHECK (
        sync_status IN ('pending', 'synced', 'failed', 'conflict')
    ),
    last_sync_at TIMESTAMP,

    -- External data
    external_data JSONB,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT uq_external_order UNIQUE(organization_id, sales_channel_id, external_order_id)
);

CREATE INDEX idx_external_orders_sale ON external_order_mappings(sale_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_external_orders_channel ON external_order_mappings(sales_channel_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_external_orders_external_id ON external_order_mappings(external_order_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE external_order_mappings IS 'Links internal sales to external e-commerce/marketplace orders';

-- Enable RLS
ALTER TABLE external_order_mappings ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_external_orders ON external_order_mappings FOR ALL USING (organization_id IN (
    SELECT organization_id FROM user_organizations WHERE user_id = auth.uid()
));

-- Add sales_channel_id to sales table
ALTER TABLE sales ADD COLUMN IF NOT EXISTS sales_channel_id UUID REFERENCES sales_channels(id);
CREATE INDEX IF NOT EXISTS idx_sales_channel ON sales(sales_channel_id) WHERE deleted_at IS NULL;

-- =====================================================
-- AUDIT TRIGGERS
-- =====================================================

\echo 'Creating audit triggers...';

CREATE TRIGGER update_pos_sessions_updated_at BEFORE UPDATE ON pos_sessions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_cash_drawers_updated_at BEFORE UPDATE ON cash_drawers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_price_lists_updated_at BEFORE UPDATE ON price_lists
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_gift_cards_updated_at BEFORE UPDATE ON gift_cards
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_sale_returns_updated_at BEFORE UPDATE ON sale_returns
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_purchase_orders_updated_at BEFORE UPDATE ON purchase_orders
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_goods_receipts_updated_at BEFORE UPDATE ON goods_receipts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

\echo '';
\echo '==========================================';
\echo 'Advanced POS Features Created Successfully';
\echo '==========================================';
\echo 'Modules Created:';
\echo '  1. ✓ Money Handling (POS Sessions, Cash Drawers, Movements)';
\echo '  2. ✓ Advanced Pricing (Price Lists, Components, UoM)';
\echo '  3. ✓ Gift Cards & Store Credit';
\echo '  4. ✓ Returns & After-Sales (RMA, Reasons)';
\echo '  5. ✓ Procurement (Purchase Orders, Goods Receipts)';
\echo '  6. ✓ Monitoring (System Health, Error Logs)';
\echo '  7. ✓ Multi-Channel (Sales Channels, External Orders)';
\echo '==========================================';
\echo '';

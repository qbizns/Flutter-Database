-- =====================================================
-- Migration: V010 - Kitchen Operations & Order Fulfillment
-- Description: Kitchen Display System (KDS), orders, order items, kitchen stations, tickets, and order fulfillment workflow
-- Author: POS Database System
-- Date: 2025-11-09
-- Dependencies: V009 (restaurant table management)
-- =====================================================

-- =====================================================
-- SECTION 1: Kitchen Stations
-- =====================================================
-- Description: Define kitchen stations/prep areas (grill, salad bar, pizza oven, bar, dessert, etc.)
-- Purpose: Route order items to appropriate kitchen stations for KDS display

CREATE TABLE kitchen_stations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Station Details
    station_name VARCHAR(100) NOT NULL,
    station_code VARCHAR(20) NOT NULL, -- GRILL, SALAD, PIZZA, BAR, DESSERT, FRY, etc.
    station_type VARCHAR(50) DEFAULT 'kitchen', -- kitchen, bar, dessert, prep
    description TEXT,

    -- Display Configuration
    display_order INTEGER DEFAULT 0,
    color_code VARCHAR(20), -- For KDS color coding
    printer_id UUID, -- Reference to device/printer for this station

    -- Station Settings
    is_active BOOLEAN DEFAULT true,
    auto_print_tickets BOOLEAN DEFAULT true,
    alert_sound_enabled BOOLEAN DEFAULT true,
    display_config JSONB DEFAULT '{}', -- KDS display preferences

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_kitchen_stations_org ON kitchen_stations(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_stations_location ON kitchen_stations(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_stations_active ON kitchen_stations(is_active) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_kitchen_stations_code_unique ON kitchen_stations(organization_id, location_id, station_code) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_kitchen_stations_updated_at
    BEFORE UPDATE ON kitchen_stations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 2: Orders (In-Progress Restaurant Orders)
-- =====================================================
-- Description: Open/in-progress restaurant orders (different from sales which are completed transactions)
-- Purpose: Track dine-in orders from when waiter takes order until payment is completed

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Order Identification
    order_number VARCHAR(50) NOT NULL,
    display_number INTEGER, -- Sequential number shown to customers (resets daily)
    order_type VARCHAR(30) DEFAULT 'dine_in', -- dine_in, takeout, delivery, online

    -- Restaurant Context
    table_id UUID REFERENCES restaurant_tables(id) ON DELETE SET NULL,
    reservation_id UUID REFERENCES reservations(id) ON DELETE SET NULL,
    covers INTEGER DEFAULT 1, -- Number of guests

    -- Order Details
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    waiter_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Order Status
    status VARCHAR(30) DEFAULT 'draft',
    -- draft, submitted, sent_to_kitchen, preparing, ready, served, completed, cancelled, on_hold

    -- Timing
    order_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    submitted_at TIMESTAMP WITH TIME ZONE, -- When order was sent to kitchen
    kitchen_received_at TIMESTAMP WITH TIME ZONE, -- When kitchen acknowledged
    ready_at TIMESTAMP WITH TIME ZONE, -- When all items ready
    served_at TIMESTAMP WITH TIME ZONE, -- When delivered to table
    completed_at TIMESTAMP WITH TIME ZONE, -- When payment completed (converts to sale)

    -- Financial
    subtotal NUMERIC(15, 2) DEFAULT 0,
    tax_amount NUMERIC(15, 2) DEFAULT 0,
    discount_amount NUMERIC(15, 2) DEFAULT 0,
    service_charge NUMERIC(15, 2) DEFAULT 0,
    total_amount NUMERIC(15, 2) DEFAULT 0,

    -- References
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL, -- Link to completed sale
    shift_id UUID REFERENCES shifts(id) ON DELETE SET NULL,

    -- Notes
    customer_notes TEXT, -- Special requests, allergies
    kitchen_notes TEXT, -- Kitchen instructions
    internal_notes TEXT, -- Staff notes

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_orders_status CHECK (status IN (
        'draft', 'submitted', 'sent_to_kitchen', 'preparing', 'ready',
        'served', 'completed', 'cancelled', 'on_hold'
    )),
    CONSTRAINT chk_orders_type CHECK (order_type IN (
        'dine_in', 'takeout', 'delivery', 'online'
    )),
    CONSTRAINT chk_orders_covers CHECK (covers > 0),
    CONSTRAINT chk_orders_amounts CHECK (
        subtotal >= 0 AND tax_amount >= 0 AND
        discount_amount >= 0 AND service_charge >= 0 AND total_amount >= 0
    )
);

-- Indexes
CREATE INDEX idx_orders_org ON orders(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_location ON orders(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_status ON orders(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_table ON orders(table_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_waiter ON orders(waiter_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_customer ON orders(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_date ON orders(order_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_orders_type ON orders(order_type) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_orders_number_unique ON orders(organization_id, location_id, order_number) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 3: Order Items
-- =====================================================
-- Description: Line items for orders with modifiers and KDS routing
-- Purpose: Track individual items ordered, their status, and route to appropriate kitchen stations

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Order Reference
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

    -- Product Reference
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,

    -- Item Details
    item_name VARCHAR(255) NOT NULL, -- Snapshot at order time
    quantity NUMERIC(10, 2) NOT NULL DEFAULT 1,
    unit_price NUMERIC(15, 2) NOT NULL,

    -- Course and Timing
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    course_position INTEGER DEFAULT 0, -- Order within the course
    fire_time TIMESTAMP WITH TIME ZONE, -- When to start preparing (course timing)

    -- Kitchen Routing
    kitchen_station_id UUID REFERENCES kitchen_stations(id) ON DELETE SET NULL,
    kitchen_ticket_id UUID, -- Reference to kitchen ticket (added after kitchen_tickets table)

    -- Item Status
    status VARCHAR(30) DEFAULT 'pending',
    -- pending, fired, acknowledged, preparing, ready, served, cancelled, on_hold, voided

    fired_at TIMESTAMP WITH TIME ZONE,
    acknowledged_at TIMESTAMP WITH TIME ZONE, -- Kitchen acknowledged
    started_preparing_at TIMESTAMP WITH TIME ZONE,
    ready_at TIMESTAMP WITH TIME ZONE,
    served_at TIMESTAMP WITH TIME ZONE,

    -- Pricing
    modifiers_total NUMERIC(15, 2) DEFAULT 0, -- Total modifier price adjustments
    discount_amount NUMERIC(15, 2) DEFAULT 0,
    line_total NUMERIC(15, 2) NOT NULL, -- (unit_price + modifiers_total) * quantity - discount

    -- Special Instructions
    special_instructions TEXT,
    customer_notes TEXT, -- Visible to customer
    kitchen_notes TEXT, -- Visible to kitchen only

    -- Seat Assignment (for multi-guest tables)
    seat_number INTEGER,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_order_items_status CHECK (status IN (
        'pending', 'fired', 'acknowledged', 'preparing', 'ready',
        'served', 'cancelled', 'on_hold', 'voided'
    )),
    CONSTRAINT chk_order_items_quantity CHECK (quantity > 0),
    CONSTRAINT chk_order_items_prices CHECK (
        unit_price >= 0 AND modifiers_total >= 0 AND
        discount_amount >= 0 AND line_total >= 0
    )
);

-- Indexes
CREATE INDEX idx_order_items_org ON order_items(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_order ON order_items(order_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_product ON order_items(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_status ON order_items(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_course ON order_items(course_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_station ON order_items(kitchen_station_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_items_fire_time ON order_items(fire_time) WHERE deleted_at IS NULL AND status = 'pending';

-- Auto-update trigger
CREATE TRIGGER update_order_items_updated_at
    BEFORE UPDATE ON order_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 4: Order Item Modifiers
-- =====================================================
-- Description: Selected modifiers for each order item (extra cheese, no onions, well done, etc.)
-- Purpose: Track customizations and price adjustments for order items

CREATE TABLE order_item_modifiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- References
    order_item_id UUID NOT NULL REFERENCES order_items(id) ON DELETE CASCADE,
    modifier_id UUID NOT NULL REFERENCES modifiers(id) ON DELETE RESTRICT,
    modifier_group_id UUID REFERENCES modifier_groups(id) ON DELETE SET NULL,

    -- Modifier Details (snapshot at order time)
    modifier_name VARCHAR(200) NOT NULL,
    quantity INTEGER DEFAULT 1,
    price_adjustment NUMERIC(15, 2) DEFAULT 0,

    -- Display
    display_order INTEGER DEFAULT 0,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_order_item_modifiers_quantity CHECK (quantity > 0)
);

-- Indexes
CREATE INDEX idx_order_item_modifiers_org ON order_item_modifiers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_item_modifiers_item ON order_item_modifiers(order_item_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_order_item_modifiers_modifier ON order_item_modifiers(modifier_id) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_order_item_modifiers_updated_at
    BEFORE UPDATE ON order_item_modifiers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 5: Kitchen Tickets (KDS Display)
-- =====================================================
-- Description: Kitchen Display System tickets grouping order items by station
-- Purpose: Organize order items into tickets displayed on kitchen screens, organized by station and course

CREATE TABLE kitchen_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Ticket Identification
    ticket_number VARCHAR(50) NOT NULL,
    display_sequence INTEGER, -- Sequential display number for the day

    -- References
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    kitchen_station_id UUID NOT NULL REFERENCES kitchen_stations(id) ON DELETE RESTRICT,
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,

    -- Ticket Details
    ticket_type VARCHAR(30) DEFAULT 'normal', -- normal, rush, remake, special
    priority INTEGER DEFAULT 0, -- Higher = more urgent

    -- Status
    status VARCHAR(30) DEFAULT 'new',
    -- new, acknowledged, preparing, ready, served, completed, cancelled, bumped

    -- Timing
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    fired_at TIMESTAMP WITH TIME ZONE,
    acknowledged_at TIMESTAMP WITH TIME ZONE, -- Kitchen staff acknowledged
    started_at TIMESTAMP WITH TIME ZONE, -- Started preparing
    ready_at TIMESTAMP WITH TIME ZONE,
    bumped_at TIMESTAMP WITH TIME ZONE, -- Removed from KDS screen
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Time Tracking
    prep_time_minutes INTEGER, -- Actual prep time
    target_prep_time INTEGER, -- Expected prep time

    -- Table/Order Info (for display)
    table_number VARCHAR(50),
    order_type VARCHAR(30),
    covers INTEGER,
    waiter_name VARCHAR(200),

    -- Notes
    special_instructions TEXT,
    kitchen_notes TEXT,

    -- Display Configuration
    display_config JSONB DEFAULT '{}', -- Color, alerts, etc.

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_kitchen_tickets_status CHECK (status IN (
        'new', 'acknowledged', 'preparing', 'ready', 'served',
        'completed', 'cancelled', 'bumped'
    )),
    CONSTRAINT chk_kitchen_tickets_type CHECK (ticket_type IN (
        'normal', 'rush', 'remake', 'special'
    ))
);

-- Indexes
CREATE INDEX idx_kitchen_tickets_org ON kitchen_tickets(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_tickets_location ON kitchen_tickets(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_tickets_order ON kitchen_tickets(order_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_tickets_station ON kitchen_tickets(kitchen_station_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_tickets_status ON kitchen_tickets(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_tickets_created ON kitchen_tickets(created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_kitchen_tickets_active ON kitchen_tickets(kitchen_station_id, status)
    WHERE deleted_at IS NULL AND status IN ('new', 'acknowledged', 'preparing');
CREATE UNIQUE INDEX idx_kitchen_tickets_number_unique ON kitchen_tickets(organization_id, location_id, ticket_number) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_kitchen_tickets_updated_at
    BEFORE UPDATE ON kitchen_tickets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add foreign key to order_items now that kitchen_tickets exists
ALTER TABLE order_items ADD CONSTRAINT fk_order_items_kitchen_ticket
    FOREIGN KEY (kitchen_ticket_id) REFERENCES kitchen_tickets(id) ON DELETE SET NULL;

CREATE INDEX idx_order_items_ticket ON order_items(kitchen_ticket_id) WHERE deleted_at IS NULL;

-- =====================================================
-- SECTION 6: Row-Level Security (RLS) Policies
-- =====================================================

-- Enable RLS
ALTER TABLE kitchen_stations ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_item_modifiers ENABLE ROW LEVEL SECURITY;
ALTER TABLE kitchen_tickets ENABLE ROW LEVEL SECURITY;

-- =====================================================
-- RLS Policies: kitchen_stations
-- =====================================================

-- Super admins can see all kitchen stations
CREATE POLICY kitchen_stations_super_admin_all ON kitchen_stations
    FOR ALL
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL
            AND r.deleted_at IS NULL
        )
    );

-- Users can view kitchen stations in their organization
CREATE POLICY kitchen_stations_select ON kitchen_stations
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users with permissions can insert kitchen stations
CREATE POLICY kitchen_stations_insert ON kitchen_stations
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users with permissions can update kitchen stations
CREATE POLICY kitchen_stations_update ON kitchen_stations
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users with permissions can delete kitchen stations
CREATE POLICY kitchen_stations_delete ON kitchen_stations
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: orders
-- =====================================================

-- Super admins can see all orders
CREATE POLICY orders_super_admin_all ON orders
    FOR ALL
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL
            AND r.deleted_at IS NULL
        )
    );

-- Users can view orders in their organization
CREATE POLICY orders_select ON orders
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can insert orders in their organization
CREATE POLICY orders_insert ON orders
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can update orders in their organization
CREATE POLICY orders_update ON orders
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can delete orders in their organization
CREATE POLICY orders_delete ON orders
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: order_items
-- =====================================================

-- Super admins can see all order items
CREATE POLICY order_items_super_admin_all ON order_items
    FOR ALL
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL
            AND r.deleted_at IS NULL
        )
    );

-- Users can view order items in their organization
CREATE POLICY order_items_select ON order_items
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can insert order items in their organization
CREATE POLICY order_items_insert ON order_items
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can update order items in their organization
CREATE POLICY order_items_update ON order_items
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can delete order items in their organization
CREATE POLICY order_items_delete ON order_items
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: order_item_modifiers
-- =====================================================

-- Super admins can see all order item modifiers
CREATE POLICY order_item_modifiers_super_admin_all ON order_item_modifiers
    FOR ALL
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL
            AND r.deleted_at IS NULL
        )
    );

-- Users can view order item modifiers in their organization
CREATE POLICY order_item_modifiers_select ON order_item_modifiers
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can insert order item modifiers in their organization
CREATE POLICY order_item_modifiers_insert ON order_item_modifiers
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can update order item modifiers in their organization
CREATE POLICY order_item_modifiers_update ON order_item_modifiers
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can delete order item modifiers in their organization
CREATE POLICY order_item_modifiers_delete ON order_item_modifiers
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: kitchen_tickets
-- =====================================================

-- Super admins can see all kitchen tickets
CREATE POLICY kitchen_tickets_super_admin_all ON kitchen_tickets
    FOR ALL
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL
            AND r.deleted_at IS NULL
        )
    );

-- Users can view kitchen tickets in their organization
CREATE POLICY kitchen_tickets_select ON kitchen_tickets
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can insert kitchen tickets in their organization
CREATE POLICY kitchen_tickets_insert ON kitchen_tickets
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can update kitchen tickets in their organization
CREATE POLICY kitchen_tickets_update ON kitchen_tickets
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- Users can delete kitchen tickets in their organization
CREATE POLICY kitchen_tickets_delete ON kitchen_tickets
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- SECTION 7: Helper Views
-- =====================================================

-- Active Kitchen Tickets View (for KDS displays)
CREATE OR REPLACE VIEW v_active_kitchen_tickets AS
SELECT
    kt.id,
    kt.ticket_number,
    kt.display_sequence,
    kt.order_id,
    kt.kitchen_station_id,
    ks.station_name,
    ks.station_code,
    ks.color_code,
    kt.status,
    kt.ticket_type,
    kt.priority,
    kt.table_number,
    kt.order_type,
    kt.covers,
    kt.waiter_name,
    kt.created_at,
    kt.fired_at,
    kt.acknowledged_at,
    kt.started_at,
    kt.target_prep_time,
    EXTRACT(EPOCH FROM (NOW() - kt.fired_at))/60 AS elapsed_minutes,
    kt.special_instructions,
    kt.kitchen_notes,
    -- Count of items
    COUNT(oi.id) AS item_count,
    -- Aggregate item details
    json_agg(
        json_build_object(
            'item_id', oi.id,
            'item_name', oi.item_name,
            'quantity', oi.quantity,
            'status', oi.status,
            'special_instructions', oi.special_instructions,
            'seat_number', oi.seat_number,
            'modifiers', (
                SELECT json_agg(
                    json_build_object(
                        'modifier_name', oim.modifier_name,
                        'quantity', oim.quantity
                    )
                )
                FROM order_item_modifiers oim
                WHERE oim.order_item_id = oi.id
                AND oim.deleted_at IS NULL
            )
        ) ORDER BY oi.course_position, oi.created_at
    ) AS items
FROM kitchen_tickets kt
JOIN kitchen_stations ks ON kt.kitchen_station_id = ks.id
LEFT JOIN order_items oi ON oi.kitchen_ticket_id = kt.id AND oi.deleted_at IS NULL
WHERE kt.deleted_at IS NULL
AND kt.status IN ('new', 'acknowledged', 'preparing')
GROUP BY
    kt.id, kt.ticket_number, kt.display_sequence, kt.order_id,
    kt.kitchen_station_id, ks.station_name, ks.station_code, ks.color_code,
    kt.status, kt.ticket_type, kt.priority, kt.table_number, kt.order_type,
    kt.covers, kt.waiter_name, kt.created_at, kt.fired_at, kt.acknowledged_at,
    kt.started_at, kt.target_prep_time, kt.special_instructions, kt.kitchen_notes;

-- Open Orders Summary View
CREATE OR REPLACE VIEW v_open_orders_summary AS
SELECT
    o.id,
    o.order_number,
    o.display_number,
    o.order_type,
    o.status,
    o.table_id,
    rt.table_number,
    o.covers,
    o.waiter_id,
    u.full_name AS waiter_name,
    o.customer_id,
    c.full_name AS customer_name,
    o.order_date,
    o.submitted_at,
    o.total_amount,
    -- Item counts
    COUNT(DISTINCT oi.id) AS total_items,
    COUNT(DISTINCT oi.id) FILTER (WHERE oi.status = 'pending') AS pending_items,
    COUNT(DISTINCT oi.id) FILTER (WHERE oi.status = 'fired') AS fired_items,
    COUNT(DISTINCT oi.id) FILTER (WHERE oi.status = 'preparing') AS preparing_items,
    COUNT(DISTINCT oi.id) FILTER (WHERE oi.status = 'ready') AS ready_items,
    COUNT(DISTINCT oi.id) FILTER (WHERE oi.status = 'served') AS served_items,
    -- Timing
    EXTRACT(EPOCH FROM (NOW() - o.order_date))/60 AS elapsed_minutes,
    o.customer_notes,
    o.kitchen_notes
FROM orders o
LEFT JOIN restaurant_tables rt ON o.table_id = rt.id
LEFT JOIN users u ON o.waiter_id = u.id
LEFT JOIN customers c ON o.customer_id = c.id
LEFT JOIN order_items oi ON o.id = oi.order_id AND oi.deleted_at IS NULL
WHERE o.deleted_at IS NULL
AND o.status NOT IN ('completed', 'cancelled')
GROUP BY
    o.id, o.order_number, o.display_number, o.order_type, o.status,
    o.table_id, rt.table_number, o.covers, o.waiter_id, u.full_name,
    o.customer_id, c.full_name, o.order_date, o.submitted_at,
    o.total_amount, o.customer_notes, o.kitchen_notes;

-- =====================================================
-- SECTION 8: Business Logic Functions
-- =====================================================

-- Function: Automatic ticket prep time calculation
CREATE OR REPLACE FUNCTION calculate_ticket_prep_time()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.bumped_at IS NOT NULL AND NEW.started_at IS NOT NULL THEN
        NEW.prep_time_minutes := EXTRACT(EPOCH FROM (NEW.bumped_at - NEW.started_at))/60;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_calculate_ticket_prep_time
    BEFORE UPDATE ON kitchen_tickets
    FOR EACH ROW
    EXECUTE FUNCTION calculate_ticket_prep_time();

-- Function: Auto-update order status based on item statuses
CREATE OR REPLACE FUNCTION update_order_status_from_items()
RETURNS TRIGGER AS $$
DECLARE
    v_order_id UUID;
    v_total_items INTEGER;
    v_ready_items INTEGER;
    v_served_items INTEGER;
BEGIN
    -- Get the order_id from the NEW or OLD record
    IF TG_OP = 'DELETE' THEN
        v_order_id := OLD.order_id;
    ELSE
        v_order_id := NEW.order_id;
    END IF;

    -- Count items
    SELECT
        COUNT(*),
        COUNT(*) FILTER (WHERE status = 'ready'),
        COUNT(*) FILTER (WHERE status = 'served')
    INTO v_total_items, v_ready_items, v_served_items
    FROM order_items
    WHERE order_id = v_order_id
    AND deleted_at IS NULL
    AND status NOT IN ('cancelled', 'voided');

    -- Update order status
    IF v_total_items = 0 THEN
        -- No items, keep current status
        RETURN NEW;
    ELSIF v_served_items = v_total_items THEN
        -- All items served
        UPDATE orders
        SET status = 'served', served_at = NOW()
        WHERE id = v_order_id AND status != 'served';
    ELSIF v_ready_items = v_total_items THEN
        -- All items ready
        UPDATE orders
        SET status = 'ready', ready_at = NOW()
        WHERE id = v_order_id AND status != 'ready';
    ELSIF v_ready_items > 0 THEN
        -- Some items ready
        UPDATE orders
        SET status = 'preparing'
        WHERE id = v_order_id AND status NOT IN ('ready', 'served');
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_order_status_from_items
    AFTER INSERT OR UPDATE OR DELETE ON order_items
    FOR EACH ROW
    EXECUTE FUNCTION update_order_status_from_items();

-- =====================================================
-- End of Migration V010
-- =====================================================

-- Summary of changes:
-- - Added 5 new tables: kitchen_stations, orders, order_items, order_item_modifiers, kitchen_tickets
-- - Created complete RLS policies for all new tables
-- - Added 2 helper views: v_active_kitchen_tickets, v_open_orders_summary
-- - Created business logic triggers for automatic status updates
-- - Full indexing strategy for performance
-- - Support for complete KDS workflow and order fulfillment

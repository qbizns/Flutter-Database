-- =====================================================
-- Migration: V011 - Delivery & Online Ordering
-- Description: Delivery zones, drivers, assignments, customer addresses, order tracking
-- Author: POS Database System
-- Date: 2025-11-09
-- Dependencies: V010 (kitchen operations)
-- =====================================================

-- =====================================================
-- SECTION 1: Delivery Zones
-- =====================================================
-- Description: Geographic areas for delivery with fees and time estimates
-- Purpose: Define delivery coverage areas, pricing, and service parameters

CREATE TABLE delivery_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Zone Details
    zone_name VARCHAR(200) NOT NULL,
    zone_code VARCHAR(50),
    description TEXT,

    -- Geographic Definition
    geofence JSONB, -- GeoJSON polygon or circle defining the zone
    postal_codes TEXT[], -- Array of postal codes in this zone
    coverage_notes TEXT,

    -- Pricing
    base_delivery_fee NUMERIC(10, 2) DEFAULT 0,
    fee_type VARCHAR(20) DEFAULT 'flat', -- flat, per_km, tiered
    minimum_order_amount NUMERIC(15, 2) DEFAULT 0, -- Minimum order for delivery
    free_delivery_threshold NUMERIC(15, 2), -- Free delivery above this amount

    -- Service Parameters
    estimated_delivery_time_minutes INTEGER DEFAULT 30,
    max_delivery_time_minutes INTEGER DEFAULT 60,
    priority INTEGER DEFAULT 0, -- Higher priority zones served first

    -- Status
    is_active BOOLEAN DEFAULT true,
    active_hours JSONB DEFAULT '{}', -- Operating hours per day

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_delivery_zones_fee_type CHECK (fee_type IN ('flat', 'per_km', 'tiered')),
    CONSTRAINT chk_delivery_zones_amounts CHECK (
        base_delivery_fee >= 0 AND minimum_order_amount >= 0
    )
);

-- Indexes
CREATE INDEX idx_delivery_zones_org ON delivery_zones(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_zones_location ON delivery_zones(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_zones_active ON delivery_zones(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_zones_postal_codes ON delivery_zones USING GIN(postal_codes) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_delivery_zones_updated_at
    BEFORE UPDATE ON delivery_zones
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 2: Delivery Drivers
-- =====================================================
-- Description: Delivery driver master data with vehicle and contact information
-- Purpose: Manage driver profiles, availability, and performance tracking

CREATE TABLE delivery_drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Driver Identity (may or may not be a system user)
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    driver_code VARCHAR(50) NOT NULL,
    full_name VARCHAR(200) NOT NULL,

    -- Contact Information
    phone VARCHAR(20),
    email VARCHAR(255),
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(20),

    -- Vehicle Information
    vehicle_type VARCHAR(50), -- bike, scooter, motorcycle, car, van
    vehicle_make VARCHAR(100),
    vehicle_model VARCHAR(100),
    vehicle_year INTEGER,
    vehicle_color VARCHAR(50),
    license_plate VARCHAR(50),

    -- Insurance & License
    drivers_license_number VARCHAR(100),
    license_expiry_date DATE,
    insurance_policy_number VARCHAR(100),
    insurance_expiry_date DATE,

    -- Work Details
    hire_date DATE,
    employment_type VARCHAR(30) DEFAULT 'full_time', -- full_time, part_time, contractor
    status VARCHAR(30) DEFAULT 'active', -- active, inactive, on_break, suspended, terminated

    -- Performance Tracking
    total_deliveries INTEGER DEFAULT 0,
    successful_deliveries INTEGER DEFAULT 0,
    rating NUMERIC(3, 2), -- Average rating from customers
    rating_count INTEGER DEFAULT 0,

    -- Current Status
    current_location JSONB, -- Real-time GPS coordinates
    is_available BOOLEAN DEFAULT false,
    last_location_update TIMESTAMP WITH TIME ZONE,

    -- Financial
    commission_rate NUMERIC(5, 2), -- Percentage of delivery fee
    payment_method VARCHAR(30), -- bank_transfer, cash, mobile_money

    -- Documents
    documents JSONB DEFAULT '{}', -- Links to uploaded documents (license, insurance, etc.)

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_delivery_drivers_status CHECK (status IN (
        'active', 'inactive', 'on_break', 'suspended', 'terminated'
    )),
    CONSTRAINT chk_delivery_drivers_employment CHECK (employment_type IN (
        'full_time', 'part_time', 'contractor'
    )),
    CONSTRAINT chk_delivery_drivers_vehicle CHECK (vehicle_type IN (
        'bike', 'scooter', 'motorcycle', 'car', 'van', 'truck', 'other'
    )),
    CONSTRAINT chk_delivery_drivers_rating CHECK (rating IS NULL OR (rating >= 0 AND rating <= 5))
);

-- Indexes
CREATE INDEX idx_delivery_drivers_org ON delivery_drivers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_drivers_user ON delivery_drivers(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_drivers_status ON delivery_drivers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_drivers_available ON delivery_drivers(is_available) WHERE deleted_at IS NULL AND status = 'active';
CREATE UNIQUE INDEX idx_delivery_drivers_code_unique ON delivery_drivers(organization_id, driver_code) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_delivery_drivers_updated_at
    BEFORE UPDATE ON delivery_drivers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 3: Driver Shifts
-- =====================================================
-- Description: Driver work schedules and availability
-- Purpose: Track when drivers are on duty and available for deliveries

CREATE TABLE driver_shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Driver Reference
    driver_id UUID NOT NULL REFERENCES delivery_drivers(id) ON DELETE CASCADE,

    -- Shift Details
    shift_date DATE NOT NULL,
    scheduled_start_time TIME,
    scheduled_end_time TIME,
    actual_start_time TIMESTAMP WITH TIME ZONE,
    actual_end_time TIMESTAMP WITH TIME ZONE,

    -- Status
    status VARCHAR(30) DEFAULT 'scheduled',
    -- scheduled, started, on_break, ended, cancelled, no_show

    -- Break Tracking
    total_break_minutes INTEGER DEFAULT 0,

    -- Performance Summary
    total_deliveries INTEGER DEFAULT 0,
    total_distance_km NUMERIC(10, 2) DEFAULT 0,
    total_earnings NUMERIC(15, 2) DEFAULT 0,

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_driver_shifts_status CHECK (status IN (
        'scheduled', 'started', 'on_break', 'ended', 'cancelled', 'no_show'
    ))
);

-- Indexes
CREATE INDEX idx_driver_shifts_org ON driver_shifts(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_driver_shifts_driver ON driver_shifts(driver_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_driver_shifts_date ON driver_shifts(shift_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_driver_shifts_status ON driver_shifts(status) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_driver_shifts_updated_at
    BEFORE UPDATE ON driver_shifts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 4: Customer Addresses
-- =====================================================
-- Description: Customer delivery addresses with labels and preferences
-- Purpose: Store multiple addresses per customer for delivery orders

CREATE TABLE customer_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Customer Reference
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,

    -- Address Details
    address_label VARCHAR(100), -- Home, Office, Apartment, etc.
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state_province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'Egypt',

    -- Location
    latitude NUMERIC(10, 8),
    longitude NUMERIC(11, 8),
    location_notes TEXT, -- Delivery instructions (gate code, landmarks, etc.)

    -- Zone Reference
    delivery_zone_id UUID REFERENCES delivery_zones(id) ON DELETE SET NULL,

    -- Preferences
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_customer_addresses_org ON customer_addresses(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_addresses_customer ON customer_addresses(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_addresses_zone ON customer_addresses(delivery_zone_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_addresses_default ON customer_addresses(customer_id, is_default) WHERE deleted_at IS NULL AND is_default = true;
CREATE INDEX idx_customer_addresses_postal ON customer_addresses(postal_code) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_customer_addresses_updated_at
    BEFORE UPDATE ON customer_addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 5: Delivery Assignments
-- =====================================================
-- Description: Assign delivery orders to drivers with routing and tracking
-- Purpose: Manage driver assignments, routes, and delivery execution

CREATE TABLE delivery_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- References
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    driver_id UUID NOT NULL REFERENCES delivery_drivers(id) ON DELETE RESTRICT,
    driver_shift_id UUID REFERENCES driver_shifts(id) ON DELETE SET NULL,
    delivery_zone_id UUID REFERENCES delivery_zones(id) ON DELETE SET NULL,

    -- Delivery Address
    customer_address_id UUID REFERENCES customer_addresses(id) ON DELETE SET NULL,
    delivery_address TEXT NOT NULL, -- Full address snapshot
    delivery_location JSONB, -- GPS coordinates

    -- Assignment Details
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Status
    status VARCHAR(30) DEFAULT 'assigned',
    -- assigned, accepted, picked_up, in_transit, arrived, delivered, failed, cancelled

    -- Timing
    accepted_at TIMESTAMP WITH TIME ZONE, -- Driver accepted
    picked_up_at TIMESTAMP WITH TIME ZONE, -- Order picked up from restaurant
    dispatched_at TIMESTAMP WITH TIME ZONE, -- Left for delivery
    arrived_at TIMESTAMP WITH TIME ZONE, -- Arrived at customer location
    delivered_at TIMESTAMP WITH TIME ZONE, -- Successfully delivered
    failed_at TIMESTAMP WITH TIME ZONE, -- Failed delivery

    -- Estimated Times
    estimated_pickup_time TIMESTAMP WITH TIME ZONE,
    estimated_delivery_time TIMESTAMP WITH TIME ZONE,

    -- Distance & Route
    distance_km NUMERIC(10, 2),
    route_info JSONB, -- Route coordinates, waypoints

    -- Delivery Details
    delivery_fee NUMERIC(10, 2) DEFAULT 0,
    driver_commission NUMERIC(10, 2) DEFAULT 0,
    payment_method VARCHAR(30), -- cash, card, online, cod (cash on delivery)
    cash_collected NUMERIC(15, 2) DEFAULT 0, -- For COD orders

    -- Proof of Delivery
    signature_image_url TEXT,
    delivery_photo_url TEXT,
    recipient_name VARCHAR(200),
    delivery_notes TEXT,

    -- Failure Tracking
    failure_reason VARCHAR(100), -- customer_not_home, wrong_address, cancelled, etc.
    failure_notes TEXT,
    retry_count INTEGER DEFAULT 0,

    -- Customer Feedback
    customer_rating INTEGER, -- 1-5 stars
    customer_feedback TEXT,
    driver_notes TEXT, -- Notes from driver

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_delivery_assignments_status CHECK (status IN (
        'assigned', 'accepted', 'picked_up', 'in_transit', 'arrived',
        'delivered', 'failed', 'cancelled'
    )),
    CONSTRAINT chk_delivery_assignments_rating CHECK (
        customer_rating IS NULL OR (customer_rating >= 1 AND customer_rating <= 5)
    ),
    CONSTRAINT chk_delivery_assignments_amounts CHECK (
        delivery_fee >= 0 AND driver_commission >= 0 AND cash_collected >= 0
    )
);

-- Indexes
CREATE INDEX idx_delivery_assignments_org ON delivery_assignments(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_assignments_order ON delivery_assignments(order_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_assignments_driver ON delivery_assignments(driver_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_assignments_status ON delivery_assignments(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_assignments_shift ON delivery_assignments(driver_shift_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_assignments_zone ON delivery_assignments(delivery_zone_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_assignments_assigned_date ON delivery_assignments(assigned_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_delivery_assignments_active ON delivery_assignments(driver_id, status)
    WHERE deleted_at IS NULL AND status IN ('assigned', 'accepted', 'picked_up', 'in_transit', 'arrived');

-- Auto-update trigger
CREATE TRIGGER update_delivery_assignments_updated_at
    BEFORE UPDATE ON delivery_assignments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 6: Order Tracking Events
-- =====================================================
-- Description: Real-time tracking events for order and delivery status updates
-- Purpose: Provide detailed timeline of order progress for customers and operations

CREATE TABLE order_tracking_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- References
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    delivery_assignment_id UUID REFERENCES delivery_assignments(id) ON DELETE SET NULL,

    -- Event Details
    event_type VARCHAR(50) NOT NULL,
    -- order_placed, order_confirmed, preparing, ready_for_pickup, picked_up,
    -- out_for_delivery, nearby, arrived, delivered, cancelled, delayed

    event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    event_message TEXT, -- User-friendly message for customer display

    -- Location (for delivery tracking)
    location JSONB, -- GPS coordinates at time of event
    location_name VARCHAR(200), -- Human-readable location

    -- Actor
    actor_type VARCHAR(30), -- system, customer, staff, driver
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_name VARCHAR(200),

    -- Additional Data
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Constraints
    CONSTRAINT chk_order_tracking_events_type CHECK (event_type IN (
        'order_placed', 'order_confirmed', 'payment_received', 'preparing',
        'ready_for_pickup', 'picked_up', 'out_for_delivery', 'nearby',
        'arrived', 'delivered', 'cancelled', 'delayed', 'failed',
        'rescheduled', 'customer_contacted'
    )),
    CONSTRAINT chk_order_tracking_events_actor CHECK (actor_type IN (
        'system', 'customer', 'staff', 'driver', 'admin'
    ))
);

-- Indexes
CREATE INDEX idx_order_tracking_events_org ON order_tracking_events(organization_id);
CREATE INDEX idx_order_tracking_events_order ON order_tracking_events(order_id);
CREATE INDEX idx_order_tracking_events_delivery ON order_tracking_events(delivery_assignment_id);
CREATE INDEX idx_order_tracking_events_timestamp ON order_tracking_events(event_timestamp);
CREATE INDEX idx_order_tracking_events_type ON order_tracking_events(event_type);

-- =====================================================
-- SECTION 7: Extend Orders Table for Delivery
-- =====================================================
-- Add delivery-specific fields to orders table

ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_zone_id UUID REFERENCES delivery_zones(id) ON DELETE SET NULL;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_address_id UUID REFERENCES customer_addresses(id) ON DELETE SET NULL;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_address TEXT;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_instructions TEXT;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS delivery_fee NUMERIC(10, 2) DEFAULT 0;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS estimated_delivery_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS is_asap BOOLEAN DEFAULT true;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS scheduled_for TIMESTAMP WITH TIME ZONE;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS source VARCHAR(30) DEFAULT 'pos';
-- pos, website, mobile_app, phone, third_party

-- Add indexes for new fields
CREATE INDEX IF NOT EXISTS idx_orders_delivery_zone ON orders(delivery_zone_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_customer_address ON orders(customer_address_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_scheduled ON orders(scheduled_for) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_orders_source ON orders(source) WHERE deleted_at IS NULL;

-- =====================================================
-- SECTION 8: Row-Level Security (RLS) Policies
-- =====================================================

-- Enable RLS
ALTER TABLE delivery_zones ENABLE ROW LEVEL SECURITY;
ALTER TABLE delivery_drivers ENABLE ROW LEVEL SECURITY;
ALTER TABLE driver_shifts ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_addresses ENABLE ROW LEVEL SECURITY;
ALTER TABLE delivery_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_tracking_events ENABLE ROW LEVEL SECURITY;

-- =====================================================
-- RLS Policies: delivery_zones
-- =====================================================

CREATE POLICY delivery_zones_super_admin_all ON delivery_zones
    FOR ALL TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
        )
    );

CREATE POLICY delivery_zones_select ON delivery_zones
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_zones_insert ON delivery_zones
    FOR INSERT TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_zones_update ON delivery_zones
    FOR UPDATE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_zones_delete ON delivery_zones
    FOR DELETE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: delivery_drivers
-- =====================================================

CREATE POLICY delivery_drivers_super_admin_all ON delivery_drivers
    FOR ALL TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
        )
    );

CREATE POLICY delivery_drivers_select ON delivery_drivers
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_drivers_insert ON delivery_drivers
    FOR INSERT TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_drivers_update ON delivery_drivers
    FOR UPDATE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_drivers_delete ON delivery_drivers
    FOR DELETE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: driver_shifts
-- =====================================================

CREATE POLICY driver_shifts_super_admin_all ON driver_shifts
    FOR ALL TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
        )
    );

CREATE POLICY driver_shifts_select ON driver_shifts
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY driver_shifts_insert ON driver_shifts
    FOR INSERT TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY driver_shifts_update ON driver_shifts
    FOR UPDATE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY driver_shifts_delete ON driver_shifts
    FOR DELETE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: customer_addresses
-- =====================================================

CREATE POLICY customer_addresses_super_admin_all ON customer_addresses
    FOR ALL TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
        )
    );

CREATE POLICY customer_addresses_select ON customer_addresses
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY customer_addresses_insert ON customer_addresses
    FOR INSERT TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY customer_addresses_update ON customer_addresses
    FOR UPDATE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY customer_addresses_delete ON customer_addresses
    FOR DELETE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: delivery_assignments
-- =====================================================

CREATE POLICY delivery_assignments_super_admin_all ON delivery_assignments
    FOR ALL TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
        )
    );

CREATE POLICY delivery_assignments_select ON delivery_assignments
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_assignments_insert ON delivery_assignments
    FOR INSERT TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_assignments_update ON delivery_assignments
    FOR UPDATE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY delivery_assignments_delete ON delivery_assignments
    FOR DELETE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: order_tracking_events
-- =====================================================

CREATE POLICY order_tracking_events_super_admin_all ON order_tracking_events
    FOR ALL TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM user_roles ur
            JOIN roles r ON ur.role_id = r.id
            WHERE ur.user_id = auth.uid()
            AND r.role_name = 'Super Admin'
            AND ur.deleted_at IS NULL AND r.deleted_at IS NULL
        )
    );

CREATE POLICY order_tracking_events_select ON order_tracking_events
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY order_tracking_events_insert ON order_tracking_events
    FOR INSERT TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY order_tracking_events_update ON order_tracking_events
    FOR UPDATE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY order_tracking_events_delete ON order_tracking_events
    FOR DELETE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- SECTION 9: Helper Views
-- =====================================================

-- Active Deliveries View (for dispatch dashboard)
CREATE OR REPLACE VIEW v_active_deliveries AS
SELECT
    da.id,
    da.order_id,
    o.order_number,
    o.display_number,
    da.driver_id,
    dd.full_name AS driver_name,
    dd.phone AS driver_phone,
    dd.vehicle_type,
    dd.current_location AS driver_location,
    da.status,
    da.delivery_address,
    da.delivery_location,
    da.customer_address_id,
    ca.location_notes,
    da.delivery_zone_id,
    dz.zone_name,
    da.assigned_at,
    da.estimated_delivery_time,
    da.distance_km,
    o.customer_id,
    c.full_name AS customer_name,
    c.phone AS customer_phone,
    o.total_amount,
    da.delivery_fee,
    da.payment_method,
    da.cash_collected,
    EXTRACT(EPOCH FROM (NOW() - da.assigned_at))/60 AS elapsed_minutes,
    CASE
        WHEN da.estimated_delivery_time IS NOT NULL THEN
            EXTRACT(EPOCH FROM (da.estimated_delivery_time - NOW()))/60
        ELSE NULL
    END AS remaining_minutes
FROM delivery_assignments da
JOIN orders o ON da.order_id = o.id
JOIN delivery_drivers dd ON da.driver_id = dd.id
LEFT JOIN customers c ON o.customer_id = c.id
LEFT JOIN customer_addresses ca ON da.customer_address_id = ca.id
LEFT JOIN delivery_zones dz ON da.delivery_zone_id = dz.id
WHERE da.deleted_at IS NULL
AND da.status IN ('assigned', 'accepted', 'picked_up', 'in_transit', 'arrived');

-- Driver Performance View
CREATE OR REPLACE VIEW v_driver_performance AS
SELECT
    dd.id,
    dd.driver_code,
    dd.full_name,
    dd.status,
    dd.total_deliveries,
    dd.successful_deliveries,
    dd.rating,
    dd.rating_count,
    CASE
        WHEN dd.total_deliveries > 0 THEN
            ROUND((dd.successful_deliveries::NUMERIC / dd.total_deliveries * 100), 2)
        ELSE 0
    END AS success_rate_percent,
    -- Today's stats
    COUNT(da.id) FILTER (WHERE DATE(da.assigned_at) = CURRENT_DATE) AS today_deliveries,
    COUNT(da.id) FILTER (WHERE DATE(da.assigned_at) = CURRENT_DATE AND da.status = 'delivered') AS today_completed,
    SUM(da.driver_commission) FILTER (WHERE DATE(da.assigned_at) = CURRENT_DATE) AS today_earnings,
    -- This week's stats
    COUNT(da.id) FILTER (WHERE da.assigned_at >= DATE_TRUNC('week', NOW())) AS week_deliveries,
    SUM(da.driver_commission) FILTER (WHERE da.assigned_at >= DATE_TRUNC('week', NOW())) AS week_earnings,
    -- Average delivery time (in minutes)
    AVG(
        EXTRACT(EPOCH FROM (da.delivered_at - da.picked_up_at))/60
    ) FILTER (WHERE da.status = 'delivered' AND da.picked_up_at IS NOT NULL) AS avg_delivery_time_minutes,
    -- Current active delivery
    MAX(da.id) FILTER (WHERE da.status IN ('assigned', 'accepted', 'picked_up', 'in_transit', 'arrived')) AS current_delivery_id
FROM delivery_drivers dd
LEFT JOIN delivery_assignments da ON dd.id = da.driver_id AND da.deleted_at IS NULL
WHERE dd.deleted_at IS NULL
GROUP BY dd.id, dd.driver_code, dd.full_name, dd.status, dd.total_deliveries,
         dd.successful_deliveries, dd.rating, dd.rating_count;

-- =====================================================
-- SECTION 10: Business Logic Functions
-- =====================================================

-- Function: Auto-create tracking event when delivery status changes
CREATE OR REPLACE FUNCTION create_delivery_tracking_event()
RETURNS TRIGGER AS $$
DECLARE
    v_event_type VARCHAR(50);
    v_event_message TEXT;
BEGIN
    -- Determine event type and message based on status change
    IF NEW.status = 'accepted' AND OLD.status = 'assigned' THEN
        v_event_type := 'picked_up';
        v_event_message := 'Driver has accepted your delivery';
    ELSIF NEW.status = 'picked_up' THEN
        v_event_type := 'picked_up';
        v_event_message := 'Order has been picked up and is on the way';
    ELSIF NEW.status = 'in_transit' THEN
        v_event_type := 'out_for_delivery';
        v_event_message := 'Your order is out for delivery';
    ELSIF NEW.status = 'arrived' THEN
        v_event_type := 'arrived';
        v_event_message := 'Driver has arrived at your location';
    ELSIF NEW.status = 'delivered' THEN
        v_event_type := 'delivered';
        v_event_message := 'Your order has been delivered successfully';
    ELSIF NEW.status = 'failed' THEN
        v_event_type := 'failed';
        v_event_message := 'Delivery attempt failed: ' || COALESCE(NEW.failure_reason, 'Unknown reason');
    ELSIF NEW.status = 'cancelled' THEN
        v_event_type := 'cancelled';
        v_event_message := 'Delivery has been cancelled';
    ELSE
        -- No event for other status changes
        RETURN NEW;
    END IF;

    -- Insert tracking event
    INSERT INTO order_tracking_events (
        organization_id,
        order_id,
        delivery_assignment_id,
        event_type,
        event_timestamp,
        event_message,
        actor_type,
        actor_id
    ) VALUES (
        NEW.organization_id,
        NEW.order_id,
        NEW.id,
        v_event_type,
        NOW(),
        v_event_message,
        'driver',
        (SELECT user_id FROM delivery_drivers WHERE id = NEW.driver_id)
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_delivery_tracking_event
    AFTER UPDATE ON delivery_assignments
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status)
    EXECUTE FUNCTION create_delivery_tracking_event();

-- Function: Update driver statistics on delivery completion
CREATE OR REPLACE FUNCTION update_driver_stats_on_delivery()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'delivered' AND OLD.status != 'delivered' THEN
        -- Update driver stats
        UPDATE delivery_drivers
        SET
            total_deliveries = total_deliveries + 1,
            successful_deliveries = successful_deliveries + 1,
            updated_at = NOW()
        WHERE id = NEW.driver_id;

        -- Update driver shift stats
        UPDATE driver_shifts
        SET
            total_deliveries = total_deliveries + 1,
            total_distance_km = total_distance_km + COALESCE(NEW.distance_km, 0),
            total_earnings = total_earnings + COALESCE(NEW.driver_commission, 0),
            updated_at = NOW()
        WHERE id = NEW.driver_shift_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_driver_stats_on_delivery
    AFTER UPDATE ON delivery_assignments
    FOR EACH ROW
    WHEN (NEW.status = 'delivered' AND OLD.status IS DISTINCT FROM NEW.status)
    EXECUTE FUNCTION update_driver_stats_on_delivery();

-- Function: Update driver rating
CREATE OR REPLACE FUNCTION update_driver_rating()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.customer_rating IS NOT NULL AND OLD.customer_rating IS NULL THEN
        UPDATE delivery_drivers
        SET
            rating = (
                (COALESCE(rating, 0) * rating_count + NEW.customer_rating) /
                (rating_count + 1)
            ),
            rating_count = rating_count + 1,
            updated_at = NOW()
        WHERE id = NEW.driver_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_driver_rating
    AFTER UPDATE ON delivery_assignments
    FOR EACH ROW
    WHEN (NEW.customer_rating IS NOT NULL AND OLD.customer_rating IS NULL)
    EXECUTE FUNCTION update_driver_rating();

-- =====================================================
-- End of Migration V011
-- =====================================================

-- Summary of changes:
-- - Added 6 new tables: delivery_zones, delivery_drivers, driver_shifts, customer_addresses, delivery_assignments, order_tracking_events
-- - Extended orders table with delivery-specific fields
-- - Created complete RLS policies for all new tables
-- - Added 2 helper views: v_active_deliveries, v_driver_performance
-- - Created business logic triggers for automatic tracking and statistics
-- - Full support for delivery management, driver operations, and real-time order tracking

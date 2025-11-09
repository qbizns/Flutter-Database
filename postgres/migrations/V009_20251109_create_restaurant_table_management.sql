-- Migration V009: Restaurant & Table Management
-- Created: 2025-11-09
-- Description: Adds restaurant-specific features including floor plans, tables, reservations, modifiers, and courses
-- Dependencies: V001, V002, V003, V004

-- ============================================================================
-- FLOOR PLANS
-- ============================================================================

-- Floor plan layouts for restaurants
CREATE TABLE IF NOT EXISTS floor_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,

    -- Floor Plan Information
    floor_name VARCHAR(100) NOT NULL,
    floor_level INTEGER DEFAULT 1, -- 1 = ground floor, 2 = second floor, etc.
    display_order INTEGER DEFAULT 0,

    -- Layout Configuration (JSONB for flexibility)
    layout_config JSONB DEFAULT '{}', -- Grid size, dimensions, background image URL, etc.

    -- Status
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,

    -- Additional Information
    description TEXT,
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_floor_name_per_location UNIQUE (location_id, floor_name)
);

-- Indexes
CREATE INDEX idx_floor_plans_org_id ON floor_plans(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_floor_plans_location_id ON floor_plans(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_floor_plans_active ON floor_plans(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_floor_plans_updated_at
    BEFORE UPDATE ON floor_plans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE floor_plans IS 'Floor plan layouts for restaurant dining areas';

-- ============================================================================
-- TABLE SECTIONS
-- ============================================================================

-- Sections/zones within floor plans (smoking/non-smoking, outdoor, VIP, etc.)
CREATE TABLE IF NOT EXISTS table_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    floor_plan_id UUID REFERENCES floor_plans(id) ON DELETE CASCADE,

    -- Section Information
    section_name VARCHAR(100) NOT NULL,
    section_code VARCHAR(50),
    section_type VARCHAR(50) DEFAULT 'regular' CHECK (section_type IN ('regular', 'vip', 'outdoor', 'bar', 'smoking', 'non_smoking', 'private', 'other')),

    -- Display
    color_code VARCHAR(20),
    icon VARCHAR(50),
    display_order INTEGER DEFAULT 0,

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Additional Information
    description TEXT,
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_section_name_per_floor UNIQUE (floor_plan_id, section_name)
);

-- Indexes
CREATE INDEX idx_table_sections_org_id ON table_sections(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_table_sections_location_id ON table_sections(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_table_sections_floor_plan_id ON table_sections(floor_plan_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_table_sections_active ON table_sections(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_table_sections_updated_at
    BEFORE UPDATE ON table_sections
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE table_sections IS 'Sections or zones within restaurant floor plans';

-- ============================================================================
-- RESTAURANT TABLES
-- ============================================================================

-- Physical tables in the restaurant
CREATE TABLE IF NOT EXISTS restaurant_tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    floor_plan_id UUID REFERENCES floor_plans(id) ON DELETE SET NULL,
    section_id UUID REFERENCES table_sections(id) ON DELETE SET NULL,

    -- Table Information
    table_number VARCHAR(50) NOT NULL,
    table_name VARCHAR(100),

    -- Capacity
    min_capacity INTEGER DEFAULT 1,
    max_capacity INTEGER NOT NULL,

    -- Table Type
    table_shape VARCHAR(20) DEFAULT 'square' CHECK (table_shape IN ('square', 'round', 'rectangle', 'oval', 'custom')),
    is_combinable BOOLEAN DEFAULT false, -- Can be combined with other tables

    -- Position (for floor plan display)
    position_x NUMERIC(10, 2),
    position_y NUMERIC(10, 2),
    rotation INTEGER DEFAULT 0, -- Degrees 0-360

    -- Current Status
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved', 'cleaning', 'maintenance', 'unavailable')),

    -- Occupancy
    current_covers INTEGER DEFAULT 0, -- Current number of guests
    seated_at TIMESTAMP WITH TIME ZONE,
    current_waiter_id UUID REFERENCES users(id),

    -- Settings
    is_active BOOLEAN DEFAULT true,
    allow_online_reservation BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,

    -- Visual
    color_code VARCHAR(20),
    icon VARCHAR(50),

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_table_number_per_location UNIQUE (location_id, table_number),
    CONSTRAINT valid_capacity CHECK (max_capacity >= min_capacity AND min_capacity > 0),
    CONSTRAINT valid_covers CHECK (current_covers >= 0 AND current_covers <= max_capacity)
);

-- Indexes
CREATE INDEX idx_restaurant_tables_org_id ON restaurant_tables(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurant_tables_location_id ON restaurant_tables(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurant_tables_floor_plan_id ON restaurant_tables(floor_plan_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurant_tables_section_id ON restaurant_tables(section_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurant_tables_status ON restaurant_tables(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurant_tables_waiter ON restaurant_tables(current_waiter_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_restaurant_tables_table_number ON restaurant_tables(table_number) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_restaurant_tables_updated_at
    BEFORE UPDATE ON restaurant_tables
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE restaurant_tables IS 'Physical restaurant tables with status and position tracking';

-- ============================================================================
-- RESERVATIONS
-- ============================================================================

-- Table reservations
CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    table_id UUID REFERENCES restaurant_tables(id) ON DELETE SET NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,

    -- Reservation Information
    reservation_number VARCHAR(50) NOT NULL,
    reservation_date DATE NOT NULL,
    reservation_time TIME NOT NULL,
    duration_minutes INTEGER DEFAULT 90,

    -- Party Information
    party_size INTEGER NOT NULL,
    customer_name VARCHAR(200) NOT NULL,
    customer_phone VARCHAR(50),
    customer_email VARCHAR(255),

    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'seated', 'completed', 'cancelled', 'no_show')),

    -- Assignment
    assigned_waiter_id UUID REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE,
    seated_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Special Requests
    special_requests TEXT,
    occasion VARCHAR(100), -- birthday, anniversary, business, etc.
    dietary_restrictions TEXT,

    -- Confirmation
    confirmation_code VARCHAR(50),
    confirmed_at TIMESTAMP WITH TIME ZONE,
    confirmed_by UUID REFERENCES users(id),

    -- Notifications
    reminder_sent_at TIMESTAMP WITH TIME ZONE,
    notification_preferences JSONB DEFAULT '{}',

    -- Cancellation
    cancelled_at TIMESTAMP WITH TIME ZONE,
    cancelled_by UUID REFERENCES users(id),
    cancellation_reason TEXT,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_reservation_number_per_org UNIQUE (organization_id, reservation_number),
    CONSTRAINT positive_party_size CHECK (party_size > 0),
    CONSTRAINT positive_duration CHECK (duration_minutes > 0)
);

-- Indexes
CREATE INDEX idx_reservations_org_id ON reservations(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_location_id ON reservations(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_table_id ON reservations(table_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_customer_id ON reservations(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_status ON reservations(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_date ON reservations(reservation_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_date_time ON reservations(reservation_date, reservation_time) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_waiter ON reservations(assigned_waiter_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reservations_phone ON reservations(customer_phone) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_reservations_updated_at
    BEFORE UPDATE ON reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE reservations IS 'Restaurant table reservations with complete workflow';

-- ============================================================================
-- MODIFIER GROUPS
-- ============================================================================

-- Groups of modifiers (e.g., "Size", "Extras", "Toppings")
CREATE TABLE IF NOT EXISTS modifier_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Group Information
    group_name VARCHAR(200) NOT NULL,
    group_code VARCHAR(50),
    display_name VARCHAR(200),

    -- Selection Rules
    selection_type VARCHAR(20) DEFAULT 'single' CHECK (selection_type IN ('single', 'multiple', 'exact')),
    min_selections INTEGER DEFAULT 0,
    max_selections INTEGER,
    exact_selections INTEGER,

    -- Pricing
    is_required BOOLEAN DEFAULT false,
    affects_price BOOLEAN DEFAULT true,

    -- Display
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,

    -- Additional Information
    description TEXT,
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_group_name_per_org UNIQUE (organization_id, group_name),
    CONSTRAINT valid_selections CHECK (
        min_selections >= 0 AND
        (max_selections IS NULL OR max_selections >= min_selections) AND
        (exact_selections IS NULL OR exact_selections > 0)
    )
);

-- Indexes
CREATE INDEX idx_modifier_groups_org_id ON modifier_groups(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_modifier_groups_active ON modifier_groups(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_modifier_groups_updated_at
    BEFORE UPDATE ON modifier_groups
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE modifier_groups IS 'Groups of modifiers with selection rules (Size, Extras, Toppings)';

-- ============================================================================
-- MODIFIERS
-- ============================================================================

-- Individual modifiers within groups
CREATE TABLE IF NOT EXISTS modifiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    modifier_group_id UUID REFERENCES modifier_groups(id) ON DELETE CASCADE,

    -- Modifier Information
    modifier_name VARCHAR(200) NOT NULL,
    modifier_code VARCHAR(50),
    display_name VARCHAR(200),

    -- Pricing
    price_adjustment NUMERIC(15, 2) DEFAULT 0, -- Can be negative for discounts
    price_type VARCHAR(20) DEFAULT 'add' CHECK (price_type IN ('add', 'multiply', 'replace')),

    -- Availability
    is_available BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,

    -- Stock (optional)
    track_inventory BOOLEAN DEFAULT false,
    current_stock NUMERIC(10, 2) DEFAULT 0,
    low_stock_threshold NUMERIC(10, 2) DEFAULT 0,

    -- Display
    display_order INTEGER DEFAULT 0,
    image_url TEXT,

    -- Additional Information
    description TEXT,
    allergen_info TEXT,
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_modifier_name_per_group UNIQUE (modifier_group_id, modifier_name)
);

-- Indexes
CREATE INDEX idx_modifiers_org_id ON modifiers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_modifiers_group_id ON modifiers(modifier_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_modifiers_available ON modifiers(is_available) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_modifiers_updated_at
    BEFORE UPDATE ON modifiers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE modifiers IS 'Individual modifiers within modifier groups';

-- ============================================================================
-- PRODUCT MODIFIER GROUPS
-- ============================================================================

-- Link products to modifier groups
CREATE TABLE IF NOT EXISTS product_modifier_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    modifier_group_id UUID NOT NULL REFERENCES modifier_groups(id) ON DELETE CASCADE,

    -- Settings
    is_required BOOLEAN DEFAULT false,
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,

    -- Override group settings (optional)
    override_min_selections INTEGER,
    override_max_selections INTEGER,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_product_modifier_group UNIQUE (product_id, modifier_group_id)
);

-- Indexes
CREATE INDEX idx_product_modifier_groups_org_id ON product_modifier_groups(organization_id);
CREATE INDEX idx_product_modifier_groups_product_id ON product_modifier_groups(product_id);
CREATE INDEX idx_product_modifier_groups_group_id ON product_modifier_groups(modifier_group_id);
CREATE INDEX idx_product_modifier_groups_active ON product_modifier_groups(is_active);

-- Auto-update trigger
CREATE TRIGGER update_product_modifier_groups_updated_at
    BEFORE UPDATE ON product_modifier_groups
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE product_modifier_groups IS 'Link products to available modifier groups';

-- ============================================================================
-- COURSES
-- ============================================================================

-- Course definitions (appetizer, main, dessert, etc.)
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Course Information
    course_name VARCHAR(100) NOT NULL,
    course_code VARCHAR(50),
    course_type VARCHAR(50) DEFAULT 'main' CHECK (course_type IN ('appetizer', 'soup', 'salad', 'main', 'side', 'dessert', 'beverage', 'other')),

    -- Timing
    typical_duration_minutes INTEGER DEFAULT 15,
    fire_delay_minutes INTEGER DEFAULT 0, -- Delay before firing to kitchen

    -- Display
    display_order INTEGER DEFAULT 0,
    color_code VARCHAR(20),
    icon VARCHAR(50),

    -- Status
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,

    -- Additional Information
    description TEXT,
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT unique_course_name_per_org UNIQUE (organization_id, course_name)
);

-- Indexes
CREATE INDEX idx_courses_org_id ON courses(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_courses_course_type ON courses(course_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_courses_active ON courses(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_courses_updated_at
    BEFORE UPDATE ON courses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE courses IS 'Course definitions for meal sequencing (appetizer, main, dessert)';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS on all new tables
ALTER TABLE floor_plans ENABLE ROW LEVEL SECURITY;
ALTER TABLE table_sections ENABLE ROW LEVEL SECURITY;
ALTER TABLE restaurant_tables ENABLE ROW LEVEL SECURITY;
ALTER TABLE reservations ENABLE ROW LEVEL SECURITY;
ALTER TABLE modifier_groups ENABLE ROW LEVEL SECURITY;
ALTER TABLE modifiers ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_modifier_groups ENABLE ROW LEVEL SECURITY;
ALTER TABLE courses ENABLE ROW LEVEL SECURITY;

-- RLS Policies (following same pattern as previous migrations)
-- FLOOR_PLANS
CREATE POLICY floor_plans_super_admin_all ON floor_plans FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY floor_plans_select_own_org ON floor_plans FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY floor_plans_insert_own_org ON floor_plans FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY floor_plans_update_own_org ON floor_plans FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY floor_plans_delete_own_org ON floor_plans FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- TABLE_SECTIONS
CREATE POLICY table_sections_super_admin_all ON table_sections FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY table_sections_select_own_org ON table_sections FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY table_sections_insert_own_org ON table_sections FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY table_sections_update_own_org ON table_sections FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY table_sections_delete_own_org ON table_sections FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- RESTAURANT_TABLES
CREATE POLICY restaurant_tables_super_admin_all ON restaurant_tables FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY restaurant_tables_select_own_org ON restaurant_tables FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY restaurant_tables_insert_own_org ON restaurant_tables FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY restaurant_tables_update_own_org ON restaurant_tables FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY restaurant_tables_delete_own_org ON restaurant_tables FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- RESERVATIONS
CREATE POLICY reservations_super_admin_all ON reservations FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY reservations_select_own_org ON reservations FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY reservations_insert_own_org ON reservations FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY reservations_update_own_org ON reservations FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY reservations_delete_own_org ON reservations FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- MODIFIER_GROUPS
CREATE POLICY modifier_groups_super_admin_all ON modifier_groups FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY modifier_groups_select_own_org ON modifier_groups FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY modifier_groups_insert_own_org ON modifier_groups FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY modifier_groups_update_own_org ON modifier_groups FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY modifier_groups_delete_own_org ON modifier_groups FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- MODIFIERS
CREATE POLICY modifiers_super_admin_all ON modifiers FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY modifiers_select_own_org ON modifiers FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY modifiers_insert_own_org ON modifiers FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY modifiers_update_own_org ON modifiers FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY modifiers_delete_own_org ON modifiers FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- PRODUCT_MODIFIER_GROUPS
CREATE POLICY product_modifier_groups_super_admin_all ON product_modifier_groups FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY product_modifier_groups_select_own_org ON product_modifier_groups FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY product_modifier_groups_insert_own_org ON product_modifier_groups FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY product_modifier_groups_update_own_org ON product_modifier_groups FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY product_modifier_groups_delete_own_org ON product_modifier_groups FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- COURSES
CREATE POLICY courses_super_admin_all ON courses FOR ALL TO PUBLIC USING (is_super_admin());
CREATE POLICY courses_select_own_org ON courses FOR SELECT TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY courses_insert_own_org ON courses FOR INSERT TO PUBLIC WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY courses_update_own_org ON courses FOR UPDATE TO PUBLIC USING (organization_id = current_user_organization_id());
CREATE POLICY courses_delete_own_org ON courses FOR DELETE TO PUBLIC USING (organization_id = current_user_organization_id());

-- Migration completed successfully

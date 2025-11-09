-- =====================================================
-- Migration: V012 - Staff & Device Management
-- Description: Employee schedules, time & attendance, device registration, printer configuration, tips & commissions
-- Author: POS Database System
-- Date: 2025-11-09
-- Dependencies: V011 (delivery & online ordering)
-- =====================================================

-- =====================================================
-- SECTION 1: Employee Schedules
-- =====================================================
-- Description: Staff work schedules and shift planning
-- Purpose: Manage employee work schedules, availability, and labor planning

CREATE TABLE employee_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Employee Reference
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Schedule Details
    schedule_date DATE NOT NULL,
    shift_type VARCHAR(30) DEFAULT 'regular', -- regular, overtime, on_call, training
    position VARCHAR(100), -- Cashier, Waiter, Chef, Manager, etc.

    -- Timing
    scheduled_start_time TIME NOT NULL,
    scheduled_end_time TIME NOT NULL,
    break_duration_minutes INTEGER DEFAULT 0,

    -- Status
    status VARCHAR(30) DEFAULT 'scheduled',
    -- scheduled, confirmed, started, completed, cancelled, no_show, late

    -- Approval
    requires_approval BOOLEAN DEFAULT false,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP WITH TIME ZONE,

    -- Notes
    notes TEXT,
    cancellation_reason TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_employee_schedules_status CHECK (status IN (
        'scheduled', 'confirmed', 'started', 'completed', 'cancelled', 'no_show', 'late'
    )),
    CONSTRAINT chk_employee_schedules_shift_type CHECK (shift_type IN (
        'regular', 'overtime', 'on_call', 'training', 'holiday'
    ))
);

-- Indexes
CREATE INDEX idx_employee_schedules_org ON employee_schedules(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employee_schedules_location ON employee_schedules(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employee_schedules_employee ON employee_schedules(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employee_schedules_date ON employee_schedules(schedule_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_employee_schedules_status ON employee_schedules(status) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_employee_schedules_updated_at
    BEFORE UPDATE ON employee_schedules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 2: Time Clock Entries
-- =====================================================
-- Description: Employee clock in/out tracking for time & attendance
-- Purpose: Track actual working hours, breaks, and overtime

CREATE TABLE time_clock_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Employee Reference
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    schedule_id UUID REFERENCES employee_schedules(id) ON DELETE SET NULL,

    -- Entry Type
    entry_type VARCHAR(30) NOT NULL,
    -- clock_in, clock_out, break_start, break_end, meal_start, meal_end

    -- Timestamp
    entry_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    scheduled_timestamp TIMESTAMP WITH TIME ZONE,

    -- Location & Device
    device_id UUID, -- Reference to device used (added later after devices table)
    gps_location JSONB, -- GPS coordinates for mobile clock-in
    ip_address VARCHAR(50),

    -- Status
    is_late BOOLEAN DEFAULT false,
    is_early BOOLEAN DEFAULT false,
    variance_minutes INTEGER, -- Difference from scheduled time

    -- Approval/Correction
    requires_approval BOOLEAN DEFAULT false,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP WITH TIME ZONE,
    is_manual_entry BOOLEAN DEFAULT false,
    correction_notes TEXT,

    -- Photo/Verification
    photo_url TEXT, -- Photo captured at clock-in (optional biometric verification)

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
    CONSTRAINT chk_time_clock_entries_type CHECK (entry_type IN (
        'clock_in', 'clock_out', 'break_start', 'break_end',
        'meal_start', 'meal_end', 'override'
    ))
);

-- Indexes
CREATE INDEX idx_time_clock_entries_org ON time_clock_entries(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_time_clock_entries_location ON time_clock_entries(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_time_clock_entries_employee ON time_clock_entries(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_time_clock_entries_timestamp ON time_clock_entries(entry_timestamp) WHERE deleted_at IS NULL;
CREATE INDEX idx_time_clock_entries_type ON time_clock_entries(entry_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_time_clock_entries_schedule ON time_clock_entries(schedule_id) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_time_clock_entries_updated_at
    BEFORE UPDATE ON time_clock_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 3: Devices
-- =====================================================
-- Description: Device registration for tablets, terminals, printers, KDS displays, handhelds, etc.
-- Purpose: Track and manage all hardware devices used in the POS ecosystem

CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Device Identity
    device_code VARCHAR(50) NOT NULL,
    device_name VARCHAR(200) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    -- pos_terminal, tablet, kds_display, customer_display, printer, handheld_scanner,
    -- price_checker, kiosk, kitchen_printer, receipt_printer, label_printer

    -- Hardware Details
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    serial_number VARCHAR(200),
    mac_address VARCHAR(50),
    ip_address VARCHAR(50),

    -- Device Configuration
    device_config JSONB DEFAULT '{}', -- Device-specific settings
    screen_resolution VARCHAR(20),
    os_version VARCHAR(50),

    -- Connection
    connection_type VARCHAR(30) DEFAULT 'network', -- network, bluetooth, usb, serial
    connection_string TEXT, -- Connection details (IP, COM port, etc.)

    -- Status
    status VARCHAR(30) DEFAULT 'inactive',
    -- active, inactive, maintenance, offline, error

    last_online_at TIMESTAMP WITH TIME ZONE,
    last_heartbeat_at TIMESTAMP WITH TIME ZONE,

    -- Assignment
    assigned_to_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_to_station_id UUID REFERENCES kitchen_stations(id) ON DELETE SET NULL,

    -- License/Warranty
    purchase_date DATE,
    warranty_expiry_date DATE,
    license_key VARCHAR(255),
    license_expiry_date DATE,

    -- Notes
    installation_notes TEXT,
    maintenance_notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_devices_status CHECK (status IN (
        'active', 'inactive', 'maintenance', 'offline', 'error'
    )),
    CONSTRAINT chk_devices_type CHECK (device_type IN (
        'pos_terminal', 'tablet', 'kds_display', 'customer_display',
        'printer', 'handheld_scanner', 'price_checker', 'kiosk',
        'kitchen_printer', 'receipt_printer', 'label_printer', 'scale', 'other'
    )),
    CONSTRAINT chk_devices_connection CHECK (connection_type IN (
        'network', 'bluetooth', 'usb', 'serial', 'wifi'
    ))
);

-- Indexes
CREATE INDEX idx_devices_org ON devices(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_location ON devices(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_type ON devices(device_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_status ON devices(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_assigned_user ON devices(assigned_to_user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_assigned_station ON devices(assigned_to_station_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_devices_code_unique ON devices(organization_id, device_code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_devices_serial_unique ON devices(organization_id, serial_number) WHERE deleted_at IS NULL AND serial_number IS NOT NULL;

-- Auto-update trigger
CREATE TRIGGER update_devices_updated_at
    BEFORE UPDATE ON devices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add foreign key to time_clock_entries now that devices table exists
ALTER TABLE time_clock_entries ADD CONSTRAINT fk_time_clock_entries_device
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE SET NULL;

CREATE INDEX idx_time_clock_entries_device ON time_clock_entries(device_id) WHERE deleted_at IS NULL;

-- =====================================================
-- SECTION 4: Printer Configurations
-- =====================================================
-- Description: Printer routing rules and configurations
-- Purpose: Define which documents print to which printers based on rules

CREATE TABLE printer_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Printer Reference
    printer_device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,

    -- Routing Rules
    document_type VARCHAR(50) NOT NULL,
    -- receipt, kitchen_ticket, order_ticket, invoice, label, report, barcode

    -- Filters (when to use this printer)
    filter_order_type VARCHAR(30), -- dine_in, takeout, delivery (NULL = all)
    filter_kitchen_station_id UUID REFERENCES kitchen_stations(id) ON DELETE SET NULL,
    filter_product_category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    filter_course_id UUID REFERENCES courses(id) ON DELETE SET NULL,

    -- Print Settings
    number_of_copies INTEGER DEFAULT 1,
    auto_print BOOLEAN DEFAULT true,
    print_priority INTEGER DEFAULT 0, -- Higher priority prints first

    -- Template
    template_config JSONB DEFAULT '{}', -- Print template configuration
    paper_size VARCHAR(20) DEFAULT 'thermal_80mm', -- thermal_58mm, thermal_80mm, a4, letter
    print_orientation VARCHAR(20) DEFAULT 'portrait', -- portrait, landscape

    -- Status
    is_active BOOLEAN DEFAULT true,

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
    CONSTRAINT chk_printer_configurations_doc_type CHECK (document_type IN (
        'receipt', 'kitchen_ticket', 'order_ticket', 'invoice', 'label',
        'report', 'barcode', 'packing_slip', 'delivery_label'
    )),
    CONSTRAINT chk_printer_configurations_copies CHECK (number_of_copies > 0)
);

-- Indexes
CREATE INDEX idx_printer_configurations_org ON printer_configurations(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_printer_configurations_location ON printer_configurations(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_printer_configurations_printer ON printer_configurations(printer_device_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_printer_configurations_doc_type ON printer_configurations(document_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_printer_configurations_active ON printer_configurations(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_printer_configurations_station ON printer_configurations(filter_kitchen_station_id) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_printer_configurations_updated_at
    BEFORE UPDATE ON printer_configurations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 5: Tip Pools
-- =====================================================
-- Description: Tip pooling configurations and rules
-- Purpose: Define how tips are collected and distributed among staff

CREATE TABLE tip_pools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Pool Details
    pool_name VARCHAR(200) NOT NULL,
    pool_type VARCHAR(30) DEFAULT 'daily', -- daily, weekly, shift, individual
    description TEXT,

    -- Distribution Rules
    distribution_method VARCHAR(30) DEFAULT 'equal', -- equal, hours_worked, points, custom
    distribution_config JSONB DEFAULT '{}', -- Custom distribution rules

    -- Eligible Positions
    eligible_positions TEXT[], -- Array of position names eligible for this pool

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_tip_pools_type CHECK (pool_type IN (
        'daily', 'weekly', 'shift', 'individual'
    )),
    CONSTRAINT chk_tip_pools_method CHECK (distribution_method IN (
        'equal', 'hours_worked', 'points', 'sales_based', 'custom'
    ))
);

-- Indexes
CREATE INDEX idx_tip_pools_org ON tip_pools(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_pools_location ON tip_pools(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_pools_active ON tip_pools(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_tip_pools_updated_at
    BEFORE UPDATE ON tip_pools
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 6: Tip Distributions
-- =====================================================
-- Description: Tip tracking and distribution to staff members
-- Purpose: Record tips collected and distributed to employees

CREATE TABLE tip_distributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Pool Reference
    tip_pool_id UUID REFERENCES tip_pools(id) ON DELETE SET NULL,

    -- Period
    distribution_date DATE NOT NULL,
    period_start TIMESTAMP WITH TIME ZONE,
    period_end TIMESTAMP WITH TIME ZONE,

    -- References
    shift_id UUID REFERENCES shifts(id) ON DELETE SET NULL,
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Source
    source_type VARCHAR(30) NOT NULL, -- sale, cash, card, delivery, pool
    source_sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    source_order_id UUID REFERENCES orders(id) ON DELETE SET NULL,

    -- Amounts
    tip_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    distribution_amount NUMERIC(15, 2) NOT NULL DEFAULT 0, -- Actual amount received by employee
    distribution_percentage NUMERIC(5, 2), -- Percentage of pool

    -- Payment
    payment_status VARCHAR(30) DEFAULT 'pending', -- pending, paid, voided
    payment_method VARCHAR(30), -- cash, payroll, bank_transfer
    paid_at TIMESTAMP WITH TIME ZONE,
    paid_by UUID REFERENCES users(id) ON DELETE SET NULL,

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
    CONSTRAINT chk_tip_distributions_source CHECK (source_type IN (
        'sale', 'cash', 'card', 'delivery', 'pool', 'adjustment'
    )),
    CONSTRAINT chk_tip_distributions_status CHECK (payment_status IN (
        'pending', 'paid', 'voided'
    )),
    CONSTRAINT chk_tip_distributions_amounts CHECK (
        tip_amount >= 0 AND distribution_amount >= 0
    )
);

-- Indexes
CREATE INDEX idx_tip_distributions_org ON tip_distributions(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_distributions_location ON tip_distributions(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_distributions_employee ON tip_distributions(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_distributions_date ON tip_distributions(distribution_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_distributions_pool ON tip_distributions(tip_pool_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_distributions_shift ON tip_distributions(shift_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_tip_distributions_status ON tip_distributions(payment_status) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_tip_distributions_updated_at
    BEFORE UPDATE ON tip_distributions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 7: Staff Commissions
-- =====================================================
-- Description: Commission tracking for sales staff
-- Purpose: Track and calculate commissions earned by sales employees

CREATE TABLE staff_commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,

    -- Employee Reference
    employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Period
    commission_date DATE NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,

    -- Source
    source_type VARCHAR(30) NOT NULL, -- sale, order, product, target
    source_sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    source_order_id UUID REFERENCES orders(id) ON DELETE SET NULL,

    -- Commission Details
    commission_type VARCHAR(30) DEFAULT 'percentage', -- percentage, fixed, tiered
    commission_rate NUMERIC(5, 2), -- Percentage rate
    sales_amount NUMERIC(15, 2) DEFAULT 0, -- Sales amount used for calculation
    commission_amount NUMERIC(15, 2) NOT NULL DEFAULT 0, -- Commission earned

    -- Status
    status VARCHAR(30) DEFAULT 'pending', -- pending, approved, paid, voided
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP WITH TIME ZONE,

    -- Payment
    payment_date DATE,
    payment_method VARCHAR(30), -- payroll, bank_transfer, cash
    paid_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Notes
    notes TEXT,
    calculation_notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_staff_commissions_source CHECK (source_type IN (
        'sale', 'order', 'product', 'category', 'target', 'bonus'
    )),
    CONSTRAINT chk_staff_commissions_type CHECK (commission_type IN (
        'percentage', 'fixed', 'tiered', 'bonus'
    )),
    CONSTRAINT chk_staff_commissions_status CHECK (status IN (
        'pending', 'approved', 'paid', 'voided', 'disputed'
    )),
    CONSTRAINT chk_staff_commissions_amounts CHECK (
        sales_amount >= 0 AND commission_amount >= 0
    )
);

-- Indexes
CREATE INDEX idx_staff_commissions_org ON staff_commissions(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_commissions_location ON staff_commissions(location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_commissions_employee ON staff_commissions(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_commissions_date ON staff_commissions(commission_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_commissions_status ON staff_commissions(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_commissions_sale ON staff_commissions(source_sale_id) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_staff_commissions_updated_at
    BEFORE UPDATE ON staff_commissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- SECTION 8: Row-Level Security (RLS) Policies
-- =====================================================

-- Enable RLS
ALTER TABLE employee_schedules ENABLE ROW LEVEL SECURITY;
ALTER TABLE time_clock_entries ENABLE ROW LEVEL SECURITY;
ALTER TABLE devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE printer_configurations ENABLE ROW LEVEL SECURITY;
ALTER TABLE tip_pools ENABLE ROW LEVEL SECURITY;
ALTER TABLE tip_distributions ENABLE ROW LEVEL SECURITY;
ALTER TABLE staff_commissions ENABLE ROW LEVEL SECURITY;

-- =====================================================
-- RLS Policies: employee_schedules
-- =====================================================

CREATE POLICY employee_schedules_super_admin_all ON employee_schedules
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

CREATE POLICY employee_schedules_select ON employee_schedules
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY employee_schedules_insert ON employee_schedules
    FOR INSERT TO PUBLIC
    WITH CHECK (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY employee_schedules_update ON employee_schedules
    FOR UPDATE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

CREATE POLICY employee_schedules_delete ON employee_schedules
    FOR DELETE TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );

-- =====================================================
-- RLS Policies: time_clock_entries (similar pattern for all remaining tables)
-- =====================================================

CREATE POLICY time_clock_entries_super_admin_all ON time_clock_entries
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

CREATE POLICY time_clock_entries_select ON time_clock_entries FOR SELECT TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY time_clock_entries_insert ON time_clock_entries FOR INSERT TO PUBLIC WITH CHECK (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY time_clock_entries_update ON time_clock_entries FOR UPDATE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY time_clock_entries_delete ON time_clock_entries FOR DELETE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));

-- devices
CREATE POLICY devices_super_admin_all ON devices FOR ALL TO PUBLIC USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = auth.uid() AND r.role_name = 'Super Admin' AND ur.deleted_at IS NULL AND r.deleted_at IS NULL));
CREATE POLICY devices_select ON devices FOR SELECT TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY devices_insert ON devices FOR INSERT TO PUBLIC WITH CHECK (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY devices_update ON devices FOR UPDATE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY devices_delete ON devices FOR DELETE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));

-- printer_configurations
CREATE POLICY printer_configurations_super_admin_all ON printer_configurations FOR ALL TO PUBLIC USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = auth.uid() AND r.role_name = 'Super Admin' AND ur.deleted_at IS NULL AND r.deleted_at IS NULL));
CREATE POLICY printer_configurations_select ON printer_configurations FOR SELECT TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY printer_configurations_insert ON printer_configurations FOR INSERT TO PUBLIC WITH CHECK (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY printer_configurations_update ON printer_configurations FOR UPDATE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY printer_configurations_delete ON printer_configurations FOR DELETE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));

-- tip_pools
CREATE POLICY tip_pools_super_admin_all ON tip_pools FOR ALL TO PUBLIC USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = auth.uid() AND r.role_name = 'Super Admin' AND ur.deleted_at IS NULL AND r.deleted_at IS NULL));
CREATE POLICY tip_pools_select ON tip_pools FOR SELECT TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY tip_pools_insert ON tip_pools FOR INSERT TO PUBLIC WITH CHECK (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY tip_pools_update ON tip_pools FOR UPDATE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY tip_pools_delete ON tip_pools FOR DELETE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));

-- tip_distributions
CREATE POLICY tip_distributions_super_admin_all ON tip_distributions FOR ALL TO PUBLIC USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = auth.uid() AND r.role_name = 'Super Admin' AND ur.deleted_at IS NULL AND r.deleted_at IS NULL));
CREATE POLICY tip_distributions_select ON tip_distributions FOR SELECT TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY tip_distributions_insert ON tip_distributions FOR INSERT TO PUBLIC WITH CHECK (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY tip_distributions_update ON tip_distributions FOR UPDATE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY tip_distributions_delete ON tip_distributions FOR DELETE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));

-- staff_commissions
CREATE POLICY staff_commissions_super_admin_all ON staff_commissions FOR ALL TO PUBLIC USING (EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = auth.uid() AND r.role_name = 'Super Admin' AND ur.deleted_at IS NULL AND r.deleted_at IS NULL));
CREATE POLICY staff_commissions_select ON staff_commissions FOR SELECT TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY staff_commissions_insert ON staff_commissions FOR INSERT TO PUBLIC WITH CHECK (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY staff_commissions_update ON staff_commissions FOR UPDATE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));
CREATE POLICY staff_commissions_delete ON staff_commissions FOR DELETE TO PUBLIC USING (organization_id IN (SELECT organization_id FROM users WHERE id = auth.uid() AND deleted_at IS NULL));

-- =====================================================
-- SECTION 9: Helper Views
-- =====================================================

-- Employee Work Hours Summary
CREATE OR REPLACE VIEW v_employee_work_hours AS
SELECT
    tce_in.employee_id,
    u.full_name AS employee_name,
    tce_in.location_id,
    DATE(tce_in.entry_timestamp) AS work_date,
    tce_in.entry_timestamp AS clock_in_time,
    tce_out.entry_timestamp AS clock_out_time,
    EXTRACT(EPOCH FROM (tce_out.entry_timestamp - tce_in.entry_timestamp))/3600 AS hours_worked,
    tce_in.schedule_id,
    es.scheduled_start_time,
    es.scheduled_end_time
FROM time_clock_entries tce_in
JOIN users u ON tce_in.employee_id = u.id
LEFT JOIN time_clock_entries tce_out ON
    tce_out.employee_id = tce_in.employee_id
    AND tce_out.entry_type = 'clock_out'
    AND DATE(tce_out.entry_timestamp) = DATE(tce_in.entry_timestamp)
    AND tce_out.entry_timestamp > tce_in.entry_timestamp
    AND tce_out.deleted_at IS NULL
LEFT JOIN employee_schedules es ON tce_in.schedule_id = es.id
WHERE tce_in.entry_type = 'clock_in'
AND tce_in.deleted_at IS NULL;

-- Active Devices View
CREATE OR REPLACE VIEW v_active_devices AS
SELECT
    d.id,
    d.device_code,
    d.device_name,
    d.device_type,
    d.status,
    d.location_id,
    l.location_name,
    d.assigned_to_user_id,
    u.full_name AS assigned_user_name,
    d.assigned_to_station_id,
    ks.station_name AS assigned_station_name,
    d.ip_address,
    d.last_online_at,
    d.last_heartbeat_at,
    EXTRACT(EPOCH FROM (NOW() - d.last_heartbeat_at))/60 AS minutes_since_heartbeat
FROM devices d
LEFT JOIN locations l ON d.location_id = l.id
LEFT JOIN users u ON d.assigned_to_user_id = u.id
LEFT JOIN kitchen_stations ks ON d.assigned_to_station_id = ks.id
WHERE d.deleted_at IS NULL
AND d.status IN ('active', 'offline');

-- =====================================================
-- SECTION 10: Business Logic Functions
-- =====================================================

-- Function: Auto-calculate work hours variance
CREATE OR REPLACE FUNCTION calculate_time_variance()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.scheduled_timestamp IS NOT NULL THEN
        NEW.variance_minutes := EXTRACT(EPOCH FROM (NEW.entry_timestamp - NEW.scheduled_timestamp))/60;

        IF NEW.variance_minutes > 5 THEN
            NEW.is_late := true;
        ELSIF NEW.variance_minutes < -5 THEN
            NEW.is_early := true;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_calculate_time_variance
    BEFORE INSERT ON time_clock_entries
    FOR EACH ROW
    EXECUTE FUNCTION calculate_time_variance();

-- Function: Update device heartbeat
CREATE OR REPLACE FUNCTION update_device_heartbeat(p_device_id UUID)
RETURNS VOID AS $$
BEGIN
    UPDATE devices
    SET
        last_heartbeat_at = NOW(),
        last_online_at = NOW(),
        status = 'active',
        updated_at = NOW()
    WHERE id = p_device_id;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- End of Migration V012
-- =====================================================

-- Summary of changes:
-- - Added 7 new tables: employee_schedules, time_clock_entries, devices, printer_configurations, tip_pools, tip_distributions, staff_commissions
-- - Created complete RLS policies for all new tables
-- - Added 2 helper views: v_employee_work_hours, v_active_devices
-- - Created business logic triggers and functions
-- - Full support for staff management, time & attendance, device management, and commission tracking
-- - Complete POS ecosystem database now ready with 50+ tables supporting 16+ applications

-- ============================================================================
-- Seed Data: Core System Data
-- Description: Seeds organizations, users, roles, and permissions
-- ============================================================================

BEGIN;

-- ============================================================================
-- ORGANIZATIONS
-- ============================================================================

INSERT INTO organizations (id, name, slug, email, phone, address, city, state, country, postal_code, status, plan, max_users, max_products, max_locations)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Demo Retail Store', 'demo-retail', 'admin@demoretail.com', '+1234567890', '123 Main Street', 'New York', 'NY', 'USA', '10001', 'active', 'premium', 20, 10000, 5),
    ('22222222-2222-2222-2222-222222222222', 'Coffee Corner', 'coffee-corner', 'info@coffeecorner.com', '+1234567891', '456 Oak Avenue', 'San Francisco', 'CA', 'USA', '94102', 'active', 'basic', 5, 1000, 1),
    ('33333333-3333-3333-3333-333333333333', 'Tech Gadgets Pro', 'tech-gadgets', 'contact@techgadgets.com', '+1234567892', '789 Silicon Drive', 'Austin', 'TX', 'USA', '73301', 'trial', 'basic', 5, 1000, 1)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PERMISSIONS
-- ============================================================================

INSERT INTO permissions (name, slug, description, resource, action, category)
VALUES
    -- Product Permissions
    ('View Products', 'products.read', 'View product catalog', 'products', 'read', 'inventory'),
    ('Create Products', 'products.create', 'Create new products', 'products', 'create', 'inventory'),
    ('Update Products', 'products.update', 'Update existing products', 'products', 'update', 'inventory'),
    ('Delete Products', 'products.delete', 'Delete products', 'products', 'delete', 'inventory'),

    -- Sales Permissions
    ('View Sales', 'sales.read', 'View sales transactions', 'sales', 'read', 'sales'),
    ('Create Sales', 'sales.create', 'Create new sales', 'sales', 'create', 'sales'),
    ('Update Sales', 'sales.update', 'Update sales transactions', 'sales', 'update', 'sales'),
    ('Delete Sales', 'sales.delete', 'Delete sales transactions', 'sales', 'delete', 'sales'),
    ('Refund Sales', 'sales.refund', 'Process refunds', 'sales', 'refund', 'sales'),

    -- Customer Permissions
    ('View Customers', 'customers.read', 'View customer data', 'customers', 'read', 'customers'),
    ('Create Customers', 'customers.create', 'Create new customers', 'customers', 'create', 'customers'),
    ('Update Customers', 'customers.update', 'Update customer data', 'customers', 'update', 'customers'),
    ('Delete Customers', 'customers.delete', 'Delete customers', 'customers', 'delete', 'customers'),

    -- User Management Permissions
    ('View Users', 'users.read', 'View users', 'users', 'read', 'admin'),
    ('Create Users', 'users.create', 'Create new users', 'users', 'create', 'admin'),
    ('Update Users', 'users.update', 'Update users', 'users', 'update', 'admin'),
    ('Delete Users', 'users.delete', 'Delete users', 'users', 'delete', 'admin'),

    -- Reports Permissions
    ('View Reports', 'reports.read', 'View all reports', 'reports', 'read', 'reports'),
    ('Export Reports', 'reports.export', 'Export reports', 'reports', 'export', 'reports'),

    -- Settings Permissions
    ('View Settings', 'settings.read', 'View settings', 'settings', 'read', 'admin'),
    ('Update Settings', 'settings.update', 'Update settings', 'settings', 'update', 'admin'),

    -- Inventory Permissions
    ('Manage Inventory', 'inventory.manage', 'Manage inventory levels', 'inventory', 'manage', 'inventory'),
    ('View Inventory', 'inventory.read', 'View inventory levels', 'inventory', 'read', 'inventory')
ON CONFLICT (slug) DO NOTHING;

-- ============================================================================
-- ROLES
-- ============================================================================

-- System-wide roles
INSERT INTO roles (id, organization_id, name, slug, description, is_system_role, is_default)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NULL, 'Super Admin', 'super-admin', 'Full system access', TRUE, FALSE),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NULL, 'Admin', 'admin', 'Organization administrator', TRUE, FALSE),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', NULL, 'Manager', 'manager', 'Store manager', TRUE, FALSE),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', NULL, 'Cashier', 'cashier', 'Cashier/Sales person', TRUE, TRUE),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', NULL, 'Viewer', 'viewer', 'Read-only access', TRUE, FALSE)
ON CONFLICT (organization_id, slug) DO NOTHING;

-- Organization-specific roles for Demo Retail Store
INSERT INTO roles (organization_id, name, slug, description, is_system_role, is_default)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Store Manager', 'store-manager', 'Custom store manager role', FALSE, FALSE),
    ('11111111-1111-1111-1111-111111111111', 'Sales Associate', 'sales-associate', 'Custom sales associate role', FALSE, FALSE)
ON CONFLICT (organization_id, slug) DO NOTHING;

-- ============================================================================
-- ROLE PERMISSIONS MAPPING
-- ============================================================================

-- Admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    id
FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Manager permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    id
FROM permissions
WHERE slug IN (
    'products.read', 'products.create', 'products.update',
    'sales.read', 'sales.create', 'sales.update', 'sales.refund',
    'customers.read', 'customers.create', 'customers.update',
    'inventory.read', 'inventory.manage',
    'reports.read', 'reports.export'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Cashier permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    id
FROM permissions
WHERE slug IN (
    'products.read',
    'sales.read', 'sales.create',
    'customers.read', 'customers.create',
    'inventory.read'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Viewer permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    id
FROM permissions
WHERE slug IN (
    'products.read',
    'sales.read',
    'customers.read',
    'inventory.read',
    'reports.read'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================================================
-- USERS
-- ============================================================================

-- Demo Retail Store Users
INSERT INTO users (id, organization_id, email, password_hash, first_name, last_name, phone, status, email_verified)
VALUES
    ('10000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'admin@demoretail.com', '$2a$10$XQz8qvKY8jZJ5ZqJZqJZqJ', 'John', 'Admin', '+1234567890', 'active', TRUE),
    ('10000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'manager@demoretail.com', '$2a$10$XQz8qvKY8jZJ5ZqJZqJZqJ', 'Sarah', 'Manager', '+1234567891', 'active', TRUE),
    ('10000000-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'cashier1@demoretail.com', '$2a$10$XQz8qvKY8jZJ5ZqJZqJZqJ', 'Mike', 'Johnson', '+1234567892', 'active', TRUE),
    ('10000000-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'cashier2@demoretail.com', '$2a$10$XQz8qvKY8jZJ5ZqJZqJZqJ', 'Emma', 'Wilson', '+1234567893', 'active', TRUE),

-- Coffee Corner Users
    ('20000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'owner@coffeecorner.com', '$2a$10$XQz8qvKY8jZJ5ZqJZqJZqJ', 'Lisa', 'Chen', '+1234567894', 'active', TRUE),
    ('20000000-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'barista@coffeecorner.com', '$2a$10$XQz8qvKY8jZJ5ZqJZqJZqJ', 'Tom', 'Brown', '+1234567895', 'active', TRUE),

-- Tech Gadgets Pro Users
    ('30000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'admin@techgadgets.com', '$2a$10$XQz8qvKY8jZJ5ZqJZqJZqJ', 'David', 'Tech', '+1234567896', 'active', TRUE)
ON CONFLICT (organization_id, email) DO NOTHING;

-- ============================================================================
-- USER ROLES MAPPING
-- ============================================================================

INSERT INTO user_roles (user_id, role_id)
VALUES
    -- Demo Retail Store
    ('10000000-0000-0000-0000-000000000001', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'), -- Admin
    ('10000000-0000-0000-0000-000000000002', 'cccccccc-cccc-cccc-cccc-cccccccccccc'), -- Manager
    ('10000000-0000-0000-0000-000000000003', 'dddddddd-dddd-dddd-dddd-dddddddddddd'), -- Cashier
    ('10000000-0000-0000-0000-000000000004', 'dddddddd-dddd-dddd-dddd-dddddddddddd'), -- Cashier

    -- Coffee Corner
    ('20000000-0000-0000-0000-000000000001', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'), -- Admin
    ('20000000-0000-0000-0000-000000000002', 'dddddddd-dddd-dddd-dddd-dddddddddddd'), -- Cashier

    -- Tech Gadgets Pro
    ('30000000-0000-0000-0000-000000000001', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb')  -- Admin
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Core seed data inserted successfully!';
    RAISE NOTICE 'Organizations: 3';
    RAISE NOTICE 'Users: 7';
    RAISE NOTICE 'Roles: 5 system + 2 custom';
    RAISE NOTICE 'Permissions: 23';
    RAISE NOTICE '============================================';
    RAISE NOTICE '';
    RAISE NOTICE 'Test Credentials (all passwords are hashed):';
    RAISE NOTICE 'admin@demoretail.com';
    RAISE NOTICE 'manager@demoretail.com';
    RAISE NOTICE 'cashier1@demoretail.com';
    RAISE NOTICE '============================================';
END $$;

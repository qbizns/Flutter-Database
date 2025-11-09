-- ============================================================================
-- Migration: V001 - Create Core Tenant Tables
-- Description: Creates the fundamental multi-tenancy tables for SAAS architecture
-- Author: System
-- Date: 2025-11-09
-- ============================================================================

BEGIN;

-- ============================================================================
-- ORGANIZATIONS TABLE
-- Description: Stores tenant/company information
-- ============================================================================

CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Basic Information
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE, -- Used for schema naming
    description TEXT,

    -- Contact Information
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),

    -- Subscription Information
    status organization_status DEFAULT 'trial' NOT NULL,
    plan VARCHAR(50) DEFAULT 'basic',
    trial_ends_at TIMESTAMP WITH TIME ZONE,
    subscription_starts_at TIMESTAMP WITH TIME ZONE,
    subscription_ends_at TIMESTAMP WITH TIME ZONE,

    -- Limits and Quotas
    max_users INTEGER DEFAULT 5,
    max_products INTEGER DEFAULT 1000,
    max_locations INTEGER DEFAULT 1,

    -- Settings (JSON for flexibility)
    settings JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID,
    updated_by UUID,

    -- Constraints
    CONSTRAINT slug_format CHECK (slug ~ '^[a-z0-9_-]+$')
);

-- Indexes
CREATE INDEX idx_organizations_slug ON organizations(slug);
CREATE INDEX idx_organizations_status ON organizations(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_organizations_deleted_at ON organizations(deleted_at);

-- Trigger for updated_at
CREATE TRIGGER update_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE organizations IS 'Tenant organizations in the SAAS system';
COMMENT ON COLUMN organizations.slug IS 'URL-friendly identifier used for tenant schema naming';
COMMENT ON COLUMN organizations.settings IS 'Organization-specific settings and configurations';
COMMENT ON COLUMN organizations.metadata IS 'Additional metadata for extensibility';

-- ============================================================================
-- ROLES TABLE
-- Description: System and custom role definitions
-- ============================================================================

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (NULL for system-wide roles)
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,

    -- Role Information
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,

    -- Role Type
    is_system_role BOOLEAN DEFAULT FALSE, -- System roles cannot be deleted
    is_default BOOLEAN DEFAULT FALSE, -- Assigned to new users by default

    -- Settings
    settings JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    UNIQUE(organization_id, slug),
    CONSTRAINT slug_format CHECK (slug ~ '^[a-z0-9_-]+$')
);

-- Indexes
CREATE INDEX idx_roles_organization_id ON roles(organization_id);
CREATE INDEX idx_roles_slug ON roles(slug);
CREATE INDEX idx_roles_is_system ON roles(is_system_role);

-- Trigger
CREATE TRIGGER update_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE roles IS 'Role definitions for RBAC (Role-Based Access Control)';
COMMENT ON COLUMN roles.is_system_role IS 'System roles are predefined and cannot be deleted';
COMMENT ON COLUMN roles.is_default IS 'Default role assigned to new users';

-- ============================================================================
-- PERMISSIONS TABLE
-- Description: System permissions
-- ============================================================================

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Permission Information
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    resource VARCHAR(100) NOT NULL, -- e.g., 'products', 'sales', 'users'
    action VARCHAR(50) NOT NULL, -- e.g., 'create', 'read', 'update', 'delete'

    -- Categorization
    category VARCHAR(100), -- e.g., 'inventory', 'sales', 'admin'

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Constraints
    CONSTRAINT slug_format CHECK (slug ~ '^[a-z0-9_.-]+$')
);

-- Indexes
CREATE INDEX idx_permissions_slug ON permissions(slug);
CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_category ON permissions(category);

-- Trigger
CREATE TRIGGER update_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE permissions IS 'System-wide permission definitions';
COMMENT ON COLUMN permissions.resource IS 'The resource this permission applies to';
COMMENT ON COLUMN permissions.action IS 'The action allowed by this permission';

-- ============================================================================
-- ROLE_PERMISSIONS TABLE
-- Description: Maps permissions to roles
-- ============================================================================

CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- References
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Constraints
    UNIQUE(role_id, permission_id)
);

-- Indexes
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Comments
COMMENT ON TABLE role_permissions IS 'Junction table mapping roles to permissions';

-- ============================================================================
-- USERS TABLE
-- Description: System users with organization mapping
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Authentication
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255), -- NULL for OAuth/SSO users

    -- Personal Information
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    full_name VARCHAR(255) GENERATED ALWAYS AS (first_name || ' ' || last_name) STORED,
    phone VARCHAR(50),
    avatar_url TEXT,

    -- Status
    status user_status DEFAULT 'pending' NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP WITH TIME ZONE,

    -- Session Management
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip INET,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,

    -- Security
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    two_factor_secret VARCHAR(255),

    -- Settings
    settings JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID,
    updated_by UUID,

    -- Constraints
    UNIQUE(organization_id, email)
);

-- Indexes
CREATE INDEX idx_users_organization_id ON users(organization_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_full_name ON users(full_name);

-- Trigger
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE users IS 'System users with multi-tenant support';
COMMENT ON COLUMN users.full_name IS 'Auto-generated from first_name and last_name';
COMMENT ON COLUMN users.settings IS 'User preferences and settings';

-- ============================================================================
-- USER_ROLES TABLE
-- Description: Maps users to roles
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- References
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    assigned_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(user_id, role_id)
);

-- Indexes
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

-- Comments
COMMENT ON TABLE user_roles IS 'Junction table mapping users to roles';

-- ============================================================================
-- AUDIT_LOGS TABLE
-- Description: System-wide audit trail
-- ============================================================================

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Context
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Action Details
    action VARCHAR(100) NOT NULL, -- e.g., 'created', 'updated', 'deleted'
    resource_type VARCHAR(100) NOT NULL, -- e.g., 'product', 'sale', 'user'
    resource_id UUID,

    -- Change Details
    old_values JSONB,
    new_values JSONB,
    changes JSONB, -- Specific fields that changed

    -- Request Context
    ip_address INET,
    user_agent TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexes
CREATE INDEX idx_audit_logs_organization_id ON audit_logs(organization_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Partitioning hint: Consider partitioning by created_at for large datasets
-- Comments
COMMENT ON TABLE audit_logs IS 'System-wide audit trail for all important actions';
COMMENT ON COLUMN audit_logs.changes IS 'Specific fields that were modified';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V001 completed successfully!';
    RAISE NOTICE 'Core tenant tables created.';
    RAISE NOTICE '============================================';
END $$;

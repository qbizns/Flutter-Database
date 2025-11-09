-- ============================================================================
-- Migration: V003 - Implement Row-Level Security (RLS)
-- Description: Enables RLS on all tables with organization-based policies
-- Author: System
-- Date: 2025-11-09
-- ============================================================================

-- IMPORTANT: This migration implements comprehensive Row-Level Security
-- All future tables MUST include RLS policies following the patterns here

BEGIN;

-- ============================================================================
-- HELPER FUNCTIONS FOR RLS
-- ============================================================================

-- Function to get current user's organization_id from session
CREATE OR REPLACE FUNCTION current_user_organization_id()
RETURNS UUID AS $$
BEGIN
    RETURN NULLIF(current_setting('app.current_organization_id', TRUE), '')::UUID;
EXCEPTION
    WHEN OTHERS THEN
        RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION current_user_organization_id() IS
'Returns the current user organization ID from session variable app.current_organization_id';

-- Function to get current user's ID from session
CREATE OR REPLACE FUNCTION current_user_id()
RETURNS UUID AS $$
BEGIN
    RETURN NULLIF(current_setting('app.current_user_id', TRUE), '')::UUID;
EXCEPTION
    WHEN OTHERS THEN
        RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION current_user_id() IS
'Returns the current user ID from session variable app.current_user_id';

-- Function to check if current user is super admin
CREATE OR REPLACE FUNCTION is_super_admin()
RETURNS BOOLEAN AS $$
BEGIN
    RETURN COALESCE(
        current_setting('app.is_super_admin', TRUE)::BOOLEAN,
        FALSE
    );
EXCEPTION
    WHEN OTHERS THEN
        RETURN FALSE;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION is_super_admin() IS
'Returns TRUE if current user is a super admin (bypasses RLS)';

-- Function to set user context (call this at the start of each session)
CREATE OR REPLACE FUNCTION set_user_context(
    p_user_id UUID,
    p_organization_id UUID,
    p_is_super_admin BOOLEAN DEFAULT FALSE
)
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.current_user_id', p_user_id::TEXT, FALSE);
    PERFORM set_config('app.current_organization_id', p_organization_id::TEXT, FALSE);
    PERFORM set_config('app.is_super_admin', p_is_super_admin::TEXT, FALSE);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION set_user_context(UUID, UUID, BOOLEAN) IS
'Sets the current user context for RLS. Call this at the start of each session.';

-- Function to clear user context
CREATE OR REPLACE FUNCTION clear_user_context()
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.current_user_id', '', FALSE);
    PERFORM set_config('app.current_organization_id', '', FALSE);
    PERFORM set_config('app.is_super_admin', '', FALSE);
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- RLS POLICIES FOR ORGANIZATIONS TABLE
-- ============================================================================

ALTER TABLE organizations ENABLE ROW LEVEL SECURITY;

-- Super admins can see all organizations
CREATE POLICY organizations_super_admin_all
    ON organizations
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only see their own organization
CREATE POLICY organizations_select_own
    ON organizations
    FOR SELECT
    TO PUBLIC
    USING (
        id = current_user_organization_id()
    );

-- Only super admins can insert organizations
CREATE POLICY organizations_insert_super_admin
    ON organizations
    FOR INSERT
    TO PUBLIC
    WITH CHECK (is_super_admin());

-- Users can update their own organization (if they have permission)
CREATE POLICY organizations_update_own
    ON organizations
    FOR UPDATE
    TO PUBLIC
    USING (
        id = current_user_organization_id()
    );

-- Only super admins can delete organizations
CREATE POLICY organizations_delete_super_admin
    ON organizations
    FOR DELETE
    TO PUBLIC
    USING (is_super_admin());

-- ============================================================================
-- RLS POLICIES FOR USERS TABLE
-- ============================================================================

ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Super admins can see all users
CREATE POLICY users_super_admin_all
    ON users
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can see users in their organization
CREATE POLICY users_select_own_org
    ON users
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- Users can insert users in their organization
CREATE POLICY users_insert_own_org
    ON users
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

-- Users can update users in their organization
CREATE POLICY users_update_own_org
    ON users
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- Users can update themselves
CREATE POLICY users_update_self
    ON users
    FOR UPDATE
    TO PUBLIC
    USING (
        id = current_user_id()
    );

-- Users can delete users in their organization
CREATE POLICY users_delete_own_org
    ON users
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- ============================================================================
-- RLS POLICIES FOR ROLES TABLE
-- ============================================================================

ALTER TABLE roles ENABLE ROW LEVEL SECURITY;

-- Super admins can see all roles
CREATE POLICY roles_super_admin_all
    ON roles
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can see system roles (organization_id IS NULL)
CREATE POLICY roles_select_system
    ON roles
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id IS NULL
    );

-- Users can see roles in their organization
CREATE POLICY roles_select_own_org
    ON roles
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- Users can insert roles in their organization (not system roles)
CREATE POLICY roles_insert_own_org
    ON roles
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
        AND organization_id IS NOT NULL
    );

-- Users can update their organization's custom roles (not system roles)
CREATE POLICY roles_update_own_org
    ON roles
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
        AND is_system_role = FALSE
    );

-- Users can delete their organization's custom roles (not system roles)
CREATE POLICY roles_delete_own_org
    ON roles
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
        AND is_system_role = FALSE
    );

-- ============================================================================
-- RLS POLICIES FOR PERMISSIONS TABLE
-- ============================================================================

ALTER TABLE permissions ENABLE ROW LEVEL SECURITY;

-- Permissions are read-only for all users (needed for permission checks)
CREATE POLICY permissions_select_all
    ON permissions
    FOR SELECT
    TO PUBLIC
    USING (TRUE);

-- Only super admins can modify permissions
CREATE POLICY permissions_all_super_admin
    ON permissions
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- ============================================================================
-- RLS POLICIES FOR ROLE_PERMISSIONS TABLE
-- ============================================================================

ALTER TABLE role_permissions ENABLE ROW LEVEL SECURITY;

-- Super admins can see all role permissions
CREATE POLICY role_permissions_super_admin_all
    ON role_permissions
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can see role permissions for roles they can see
CREATE POLICY role_permissions_select_visible_roles
    ON role_permissions
    FOR SELECT
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM roles
            WHERE roles.id = role_permissions.role_id
            AND (
                roles.organization_id = current_user_organization_id()
                OR roles.organization_id IS NULL
            )
        )
    );

-- Users can manage role permissions for their organization's custom roles
CREATE POLICY role_permissions_insert_own_org
    ON role_permissions
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM roles
            WHERE roles.id = role_permissions.role_id
            AND roles.organization_id = current_user_organization_id()
            AND roles.is_system_role = FALSE
        )
    );

CREATE POLICY role_permissions_delete_own_org
    ON role_permissions
    FOR DELETE
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM roles
            WHERE roles.id = role_permissions.role_id
            AND roles.organization_id = current_user_organization_id()
            AND roles.is_system_role = FALSE
        )
    );

-- ============================================================================
-- RLS POLICIES FOR USER_ROLES TABLE
-- ============================================================================

ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;

-- Super admins can see all user roles
CREATE POLICY user_roles_super_admin_all
    ON user_roles
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can see user roles in their organization
CREATE POLICY user_roles_select_own_org
    ON user_roles
    FOR SELECT
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = user_roles.user_id
            AND users.organization_id = current_user_organization_id()
        )
    );

-- Users can manage user roles in their organization
CREATE POLICY user_roles_insert_own_org
    ON user_roles
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = user_roles.user_id
            AND users.organization_id = current_user_organization_id()
        )
    );

CREATE POLICY user_roles_delete_own_org
    ON user_roles
    FOR DELETE
    TO PUBLIC
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE users.id = user_roles.user_id
            AND users.organization_id = current_user_organization_id()
        )
    );

-- ============================================================================
-- RLS POLICIES FOR AUDIT_LOGS TABLE
-- ============================================================================

ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- Super admins can see all audit logs
CREATE POLICY audit_logs_super_admin_all
    ON audit_logs
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can see audit logs for their organization
CREATE POLICY audit_logs_select_own_org
    ON audit_logs
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- Anyone can insert audit logs for their organization
CREATE POLICY audit_logs_insert_own_org
    ON audit_logs
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

-- Audit logs cannot be updated or deleted (immutable)
-- No UPDATE or DELETE policies = no one can update/delete

-- ============================================================================
-- RLS POLICIES FOR CATEGORIES TABLE
-- ============================================================================

ALTER TABLE categories ENABLE ROW LEVEL SECURITY;

-- Super admins can see all categories
CREATE POLICY categories_super_admin_all
    ON categories
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only access categories in their organization
CREATE POLICY categories_select_own_org
    ON categories
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY categories_insert_own_org
    ON categories
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

CREATE POLICY categories_update_own_org
    ON categories
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY categories_delete_own_org
    ON categories
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- ============================================================================
-- RLS POLICIES FOR PRODUCTS TABLE
-- ============================================================================

ALTER TABLE products ENABLE ROW LEVEL SECURITY;

-- Super admins can see all products
CREATE POLICY products_super_admin_all
    ON products
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only access products in their organization
CREATE POLICY products_select_own_org
    ON products
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY products_insert_own_org
    ON products
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

CREATE POLICY products_update_own_org
    ON products
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY products_delete_own_org
    ON products
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- ============================================================================
-- RLS POLICIES FOR CUSTOMERS TABLE
-- ============================================================================

ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

-- Super admins can see all customers
CREATE POLICY customers_super_admin_all
    ON customers
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only access customers in their organization
CREATE POLICY customers_select_own_org
    ON customers
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY customers_insert_own_org
    ON customers
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

CREATE POLICY customers_update_own_org
    ON customers
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY customers_delete_own_org
    ON customers
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- ============================================================================
-- RLS POLICIES FOR SALES TABLE
-- ============================================================================

ALTER TABLE sales ENABLE ROW LEVEL SECURITY;

-- Super admins can see all sales
CREATE POLICY sales_super_admin_all
    ON sales
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only access sales in their organization
CREATE POLICY sales_select_own_org
    ON sales
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY sales_insert_own_org
    ON sales
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

CREATE POLICY sales_update_own_org
    ON sales
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY sales_delete_own_org
    ON sales
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- ============================================================================
-- RLS POLICIES FOR SALE_ITEMS TABLE
-- ============================================================================

ALTER TABLE sale_items ENABLE ROW LEVEL SECURITY;

-- Super admins can see all sale items
CREATE POLICY sale_items_super_admin_all
    ON sale_items
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only access sale items in their organization
CREATE POLICY sale_items_select_own_org
    ON sale_items
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY sale_items_insert_own_org
    ON sale_items
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

CREATE POLICY sale_items_update_own_org
    ON sale_items
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY sale_items_delete_own_org
    ON sale_items
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- ============================================================================
-- RLS POLICIES FOR PAYMENTS TABLE
-- ============================================================================

ALTER TABLE payments ENABLE ROW LEVEL SECURITY;

-- Super admins can see all payments
CREATE POLICY payments_super_admin_all
    ON payments
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only access payments in their organization
CREATE POLICY payments_select_own_org
    ON payments
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY payments_insert_own_org
    ON payments
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

CREATE POLICY payments_update_own_org
    ON payments
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY payments_delete_own_org
    ON payments
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- ============================================================================
-- RLS POLICIES FOR INVENTORY_TRANSACTIONS TABLE
-- ============================================================================

ALTER TABLE inventory_transactions ENABLE ROW LEVEL SECURITY;

-- Super admins can see all inventory transactions
CREATE POLICY inventory_transactions_super_admin_all
    ON inventory_transactions
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can only access inventory transactions in their organization
CREATE POLICY inventory_transactions_select_own_org
    ON inventory_transactions
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

CREATE POLICY inventory_transactions_insert_own_org
    ON inventory_transactions
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

-- Inventory transactions are immutable (no UPDATE or DELETE)
-- No UPDATE or DELETE policies = no one can update/delete

-- ============================================================================
-- CREATE RLS BYPASS ROLE (for migrations and system operations)
-- ============================================================================

-- Create a role that bypasses RLS
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'rls_bypass') THEN
        CREATE ROLE rls_bypass BYPASSRLS;
    END IF;
END
$$;

COMMENT ON ROLE rls_bypass IS 'Role that bypasses RLS for system operations and migrations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V003 completed successfully!';
    RAISE NOTICE 'Row-Level Security (RLS) enabled on all tables.';
    RAISE NOTICE '';
    RAISE NOTICE 'IMPORTANT: Set user context before queries:';
    RAISE NOTICE 'SELECT set_user_context(';
    RAISE NOTICE '    ''user_id''::UUID,';
    RAISE NOTICE '    ''organization_id''::UUID,';
    RAISE NOTICE '    FALSE  -- is_super_admin';
    RAISE NOTICE ');';
    RAISE NOTICE '';
    RAISE NOTICE 'All future tables MUST include RLS policies!';
    RAISE NOTICE '============================================';
END $$;

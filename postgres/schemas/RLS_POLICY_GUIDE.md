# Row-Level Security (RLS) Policy Guide

## Overview

Row-Level Security (RLS) is **MANDATORY** for all tables in this database. This document provides templates, patterns, and guidelines for implementing RLS on all existing and future tables.

## Why RLS is Critical

In a multi-tenant SAAS system, RLS provides:

1. **Data Isolation**: Tenants cannot access each other's data
2. **Defense in Depth**: Even if application logic fails, database enforces isolation
3. **Compliance**: Meets data privacy and security regulations
4. **Audit Trail**: All access is controlled and can be audited

## How RLS Works

### Session Variables

The system uses PostgreSQL session variables to track user context:

```sql
-- Set at the start of each database session
app.current_user_id          -- UUID of the current user
app.current_organization_id  -- UUID of the user's organization
app.is_super_admin          -- Boolean, TRUE for system admins
```

### Helper Functions

```sql
-- Get current user's organization ID
current_user_organization_id()

-- Get current user ID
current_user_id()

-- Check if user is super admin
is_super_admin()

-- Set user context (call at session start)
set_user_context(user_id UUID, organization_id UUID, is_super_admin BOOLEAN)

-- Clear user context
clear_user_context()
```

## Setting User Context

### Application Code Example

```sql
-- At the start of each user session:
SELECT set_user_context(
    '10000000-0000-0000-0000-000000000001'::UUID,  -- user_id
    '11111111-1111-1111-1111-111111111111'::UUID,  -- organization_id
    FALSE                                          -- is_super_admin
);

-- Now all queries will be filtered by this organization
SELECT * FROM products;  -- Only returns products for organization 11111...

-- Clear context when session ends
SELECT clear_user_context();
```

### For Super Admin

```sql
-- Super admin can see all data
SELECT set_user_context(
    'admin_user_id'::UUID,
    NULL,  -- organization_id not needed
    TRUE   -- is_super_admin = TRUE
);

-- Now can see ALL organizations' data
SELECT * FROM products;  -- Returns all products from all organizations
```

## RLS Policy Patterns

### Pattern 1: Standard Organization-Scoped Table

Use this for most tables that have `organization_id`:

```sql
-- Enable RLS
ALTER TABLE your_table ENABLE ROW LEVEL SECURITY;

-- Super admin bypass (always include this first)
CREATE POLICY your_table_super_admin_all
    ON your_table
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- SELECT: View own organization's data
CREATE POLICY your_table_select_own_org
    ON your_table
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- INSERT: Only insert into own organization
CREATE POLICY your_table_insert_own_org
    ON your_table
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

-- UPDATE: Only update own organization's data
CREATE POLICY your_table_update_own_org
    ON your_table
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- DELETE: Only delete own organization's data
CREATE POLICY your_table_delete_own_org
    ON your_table
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );
```

### Pattern 2: Read-Only Reference Table

Use for tables like `permissions` that everyone needs to read:

```sql
ALTER TABLE your_reference_table ENABLE ROW LEVEL SECURITY;

-- Everyone can read
CREATE POLICY your_reference_table_select_all
    ON your_reference_table
    FOR SELECT
    TO PUBLIC
    USING (TRUE);

-- Only super admins can modify
CREATE POLICY your_reference_table_all_super_admin
    ON your_reference_table
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());
```

### Pattern 3: User Self-Access

Use when users should access their own records:

```sql
ALTER TABLE your_table ENABLE ROW LEVEL SECURITY;

-- Super admin bypass
CREATE POLICY your_table_super_admin_all
    ON your_table
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can see themselves
CREATE POLICY your_table_select_self
    ON your_table
    FOR SELECT
    TO PUBLIC
    USING (
        user_id = current_user_id()
    );

-- Users can update themselves
CREATE POLICY your_table_update_self
    ON your_table
    FOR UPDATE
    TO PUBLIC
    USING (
        user_id = current_user_id()
    );
```

### Pattern 4: Immutable Audit Table

Use for audit logs and transaction history:

```sql
ALTER TABLE audit_table ENABLE ROW LEVEL SECURITY;

-- Super admin can see all
CREATE POLICY audit_table_super_admin_all
    ON audit_table
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Users can see their org's audit logs
CREATE POLICY audit_table_select_own_org
    ON audit_table
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

-- Anyone can insert (for logging)
CREATE POLICY audit_table_insert_own_org
    ON audit_table
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

-- NO UPDATE OR DELETE POLICIES = immutable
```

### Pattern 5: Junction Table with Related Entity Check

Use for many-to-many junction tables:

```sql
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;

-- Super admin bypass
CREATE POLICY user_roles_super_admin_all
    ON user_roles
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Can see user_roles for users in same organization
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

-- Can insert user_roles for users in same organization
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
```

## RLS Checklist for New Tables

When creating a new table, **ALWAYS** follow this checklist:

### ✅ Pre-Migration Checklist

- [ ] Table has `organization_id UUID` column (unless it's a system-wide table)
- [ ] Table has proper indexes on `organization_id`
- [ ] Reviewed which RLS pattern applies to this table

### ✅ Migration File Checklist

- [ ] Table creation SQL is complete
- [ ] `ALTER TABLE table_name ENABLE ROW LEVEL SECURITY;` is included
- [ ] Super admin bypass policy is created (FOR ALL)
- [ ] SELECT policy is created
- [ ] INSERT policy with `WITH CHECK` is created
- [ ] UPDATE policy is created (if applicable)
- [ ] DELETE policy is created (if applicable)
- [ ] Policies are tested with different user contexts

### ✅ Testing Checklist

Test each policy with these scenarios:

```sql
-- Test 1: Super admin sees everything
SELECT set_user_context('admin_id'::UUID, NULL, TRUE);
SELECT COUNT(*) FROM your_table;  -- Should see all records

-- Test 2: Regular user sees only their org
SELECT set_user_context('user_id'::UUID, 'org1_id'::UUID, FALSE);
SELECT COUNT(*) FROM your_table;  -- Should see only org1 records

-- Test 3: User from different org sees nothing
SELECT set_user_context('user2_id'::UUID, 'org2_id'::UUID, FALSE);
SELECT COUNT(*) FROM your_table WHERE organization_id = 'org1_id'::UUID;
-- Should return 0 (RLS blocks it)

-- Test 4: Cannot insert into another org
SELECT set_user_context('user_id'::UUID, 'org1_id'::UUID, FALSE);
INSERT INTO your_table (organization_id, ...)
VALUES ('org2_id'::UUID, ...);  -- Should FAIL with RLS violation

-- Test 5: Can insert into own org
INSERT INTO your_table (organization_id, ...)
VALUES ('org1_id'::UUID, ...);  -- Should SUCCEED
```

## Migration Template

```sql
-- ============================================================================
-- Migration: VXXX - Your Migration Name
-- Description: What this migration does
-- Author: Your Name
-- Date: YYYY-MM-DD
-- ============================================================================

BEGIN;

-- ============================================================================
-- CREATE TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS your_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- REQUIRED: Organization for multi-tenancy
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Your columns here
    name VARCHAR(255) NOT NULL,

    -- Standard audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_your_table_organization_id ON your_table(organization_id);

-- Trigger
CREATE TRIGGER update_your_table_updated_at
    BEFORE UPDATE ON your_table
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- ROW LEVEL SECURITY (MANDATORY)
-- ============================================================================

ALTER TABLE your_table ENABLE ROW LEVEL SECURITY;

-- Super admin bypass
CREATE POLICY your_table_super_admin_all
    ON your_table
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

-- Organization-scoped policies
CREATE POLICY your_table_select_own_org
    ON your_table FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY your_table_insert_own_org
    ON your_table FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY your_table_update_own_org
    ON your_table FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY your_table_delete_own_org
    ON your_table FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;
```

## Common Mistakes to Avoid

### ❌ Mistake 1: Forgetting to Enable RLS

```sql
-- WRONG: Table created but RLS not enabled
CREATE TABLE products (...);
-- Data is now accessible to everyone!
```

```sql
-- CORRECT: Always enable RLS
CREATE TABLE products (...);
ALTER TABLE products ENABLE ROW LEVEL SECURITY;
CREATE POLICY ...;
```

### ❌ Mistake 2: Missing Super Admin Bypass

```sql
-- WRONG: No super admin policy
CREATE POLICY products_select_own_org
    ON products FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());
-- Super admins can't see all data for debugging/support!
```

```sql
-- CORRECT: Always add super admin bypass FIRST
CREATE POLICY products_super_admin_all
    ON products FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY products_select_own_org
    ON products FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());
```

### ❌ Mistake 3: Wrong Policy Order

```sql
-- WRONG: Specific policies before bypass
CREATE POLICY products_select_own_org ON products FOR SELECT ...;
CREATE POLICY products_super_admin_all ON products FOR ALL ...;
-- Less readable, harder to debug
```

```sql
-- CORRECT: Super admin bypass always first
CREATE POLICY products_super_admin_all ON products FOR ALL ...;
CREATE POLICY products_select_own_org ON products FOR SELECT ...;
```

### ❌ Mistake 4: Using USING for INSERT

```sql
-- WRONG: INSERT with USING instead of WITH CHECK
CREATE POLICY products_insert_own_org
    ON products FOR INSERT TO PUBLIC
    USING (organization_id = current_user_organization_id());
-- This doesn't work correctly for INSERT!
```

```sql
-- CORRECT: INSERT uses WITH CHECK
CREATE POLICY products_insert_own_org
    ON products FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());
```

### ❌ Mistake 5: Forgetting to Set User Context

```sql
-- WRONG: Querying without setting context
SELECT * FROM products;
-- Returns nothing because current_user_organization_id() is NULL!
```

```sql
-- CORRECT: Always set context first
SELECT set_user_context('user_id'::UUID, 'org_id'::UUID, FALSE);
SELECT * FROM products;  -- Now returns correct data
```

## Debugging RLS Issues

### Check if RLS is Enabled

```sql
SELECT schemaname, tablename, rowsecurity
FROM pg_tables
WHERE schemaname = 'public';
```

### List All Policies

```sql
SELECT schemaname, tablename, policyname, cmd, qual, with_check
FROM pg_policies
WHERE schemaname = 'public'
ORDER BY tablename, policyname;
```

### Test User Context

```sql
-- Check what context is set
SELECT
    current_setting('app.current_user_id', TRUE) AS user_id,
    current_setting('app.current_organization_id', TRUE) AS org_id,
    current_setting('app.is_super_admin', TRUE) AS is_super_admin;
```

### Test Policy Effect

```sql
-- Enable query logging to see what RLS adds
SET log_statement = 'all';

-- Run query and check execution plan
EXPLAIN (VERBOSE) SELECT * FROM products;
-- You'll see RLS filters in the query plan
```

### Bypass RLS for Testing (Use with Caution!)

```sql
-- Temporarily disable RLS (as superuser only)
SET SESSION AUTHORIZATION postgres;  -- or another superuser
SET row_security = OFF;

-- Run your query
SELECT * FROM products;

-- Re-enable
SET row_security = ON;
```

## Performance Considerations

### Index on organization_id

**ALWAYS** create an index on `organization_id`:

```sql
CREATE INDEX idx_your_table_organization_id ON your_table(organization_id);
```

Without this index, every query will do a full table scan!

### Avoid Complex Subqueries

```sql
-- SLOWER: Complex subquery in every row check
CREATE POLICY slow_policy ON table1 FOR SELECT
USING (
    EXISTS (
        SELECT 1 FROM table2
        JOIN table3 ON table2.id = table3.id
        WHERE table2.org_id = current_user_organization_id()
        AND table1.foreign_id = table3.id
    )
);

-- FASTER: Simple direct check
CREATE POLICY fast_policy ON table1 FOR SELECT
USING (organization_id = current_user_organization_id());
```

### Use Partial Indexes

```sql
-- Index only active records
CREATE INDEX idx_products_active
ON products(organization_id)
WHERE deleted_at IS NULL;
```

## Security Best Practices

1. **Always set user context** - Never run queries without context
2. **Use prepared statements** - Prevent SQL injection
3. **Validate at application layer too** - Defense in depth
4. **Audit RLS policy changes** - Track who modifies policies
5. **Test thoroughly** - Test each policy with different user scenarios
6. **Monitor performance** - Watch for slow queries caused by RLS
7. **Regular security reviews** - Audit policies quarterly
8. **Document exceptions** - If a table doesn't need RLS, document why

## System Operations

### Migrations and Seed Data

For migrations and seed data, bypass RLS:

```sql
-- Option 1: Set super admin context
SELECT set_user_context(NULL, NULL, TRUE);

-- Option 2: Use rls_bypass role (in psql)
SET ROLE rls_bypass;

-- Option 3: Disable RLS for session (superuser only)
SET row_security = OFF;
```

### Backup and Restore

```bash
# Backups automatically include RLS policies
pg_dump -U postgres pos_saas > backup.sql

# Restore includes RLS policies
psql -U postgres pos_saas < backup.sql
```

## Summary

### Remember: Every New Table MUST Have:

1. ✅ `organization_id` column (unless system-wide)
2. ✅ Index on `organization_id`
3. ✅ `ALTER TABLE x ENABLE ROW LEVEL SECURITY;`
4. ✅ Super admin bypass policy (FOR ALL)
5. ✅ Organization-scoped policies (SELECT, INSERT, UPDATE, DELETE)
6. ✅ Testing with different user contexts
7. ✅ Documentation of any special cases

### Quick Reference Card

```sql
-- Enable RLS
ALTER TABLE t ENABLE ROW LEVEL SECURITY;

-- Super admin bypass (always first!)
CREATE POLICY t_super ON t FOR ALL USING (is_super_admin());

-- Standard org policies
CREATE POLICY t_sel ON t FOR SELECT USING (organization_id = current_user_organization_id());
CREATE POLICY t_ins ON t FOR INSERT WITH CHECK (organization_id = current_user_organization_id());
CREATE POLICY t_upd ON t FOR UPDATE USING (organization_id = current_user_organization_id());
CREATE POLICY t_del ON t FOR DELETE USING (organization_id = current_user_organization_id());
```

---

**This is a living document. Update it when new RLS patterns are discovered or when requirements change.**

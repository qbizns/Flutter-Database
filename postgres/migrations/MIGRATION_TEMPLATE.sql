-- ============================================================================
-- Migration: VXXX - [Your Migration Title]
-- Description: [Describe what this migration does]
-- Author: [Your Name]
-- Date: [YYYY-MM-DD]
-- ============================================================================

-- IMPORTANT CHECKLIST BEFORE RUNNING:
-- [ ] Migration number (VXXX) is sequential
-- [ ] All tables have organization_id (unless system-wide)
-- [ ] All tables have RLS policies defined
-- [ ] All tables have proper indexes
-- [ ] All tables have updated_at trigger
-- [ ] Migration tested in development
-- [ ] Rollback strategy documented

BEGIN;

-- ============================================================================
-- CREATE TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS your_table_name (
    -- Primary Key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (REQUIRED for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Your Business Columns
    name VARCHAR(255) NOT NULL,
    description TEXT,
    -- Add your columns here...

    -- Settings & Metadata (for extensibility)
    settings JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',

    -- Standard Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_name_per_org UNIQUE(organization_id, name)
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Organization index (REQUIRED for RLS performance)
CREATE INDEX idx_your_table_organization_id ON your_table_name(organization_id);

-- Soft delete index
CREATE INDEX idx_your_table_deleted_at ON your_table_name(deleted_at);

-- Add indexes for frequently queried columns
-- CREATE INDEX idx_your_table_column_name ON your_table_name(column_name);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Auto-update updated_at timestamp
CREATE TRIGGER update_your_table_updated_at
    BEFORE UPDATE ON your_table_name
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE your_table_name IS 'Description of what this table stores';
COMMENT ON COLUMN your_table_name.organization_id IS 'Links to tenant organization';
COMMENT ON COLUMN your_table_name.settings IS 'Flexible settings storage';

-- ============================================================================
-- ROW LEVEL SECURITY (MANDATORY!)
-- ============================================================================

-- Enable RLS on the table
ALTER TABLE your_table_name ENABLE ROW LEVEL SECURITY;

-- Policy 1: Super admin bypass (ALWAYS include this first!)
CREATE POLICY your_table_super_admin_all
    ON your_table_name
    FOR ALL
    TO PUBLIC
    USING (is_super_admin());

COMMENT ON POLICY your_table_super_admin_all ON your_table_name IS
'Allows super admins to bypass RLS and access all records';

-- Policy 2: SELECT - Users can view their organization's records
CREATE POLICY your_table_select_own_org
    ON your_table_name
    FOR SELECT
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

COMMENT ON POLICY your_table_select_own_org ON your_table_name IS
'Users can only SELECT records from their organization';

-- Policy 3: INSERT - Users can only insert into their organization
CREATE POLICY your_table_insert_own_org
    ON your_table_name
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        organization_id = current_user_organization_id()
    );

COMMENT ON POLICY your_table_insert_own_org ON your_table_name IS
'Users can only INSERT records for their organization';

-- Policy 4: UPDATE - Users can only update their organization's records
CREATE POLICY your_table_update_own_org
    ON your_table_name
    FOR UPDATE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

COMMENT ON POLICY your_table_update_own_org ON your_table_name IS
'Users can only UPDATE records from their organization';

-- Policy 5: DELETE - Users can only delete their organization's records
CREATE POLICY your_table_delete_own_org
    ON your_table_name
    FOR DELETE
    TO PUBLIC
    USING (
        organization_id = current_user_organization_id()
    );

COMMENT ON POLICY your_table_delete_own_org ON your_table_name IS
'Users can only DELETE records from their organization';

-- ============================================================================
-- SPECIAL CASES (uncomment if needed)
-- ============================================================================

-- IMMUTABLE TABLE (audit logs, transactions)
-- Remove UPDATE and DELETE policies to make table immutable
-- DROP POLICY your_table_update_own_org ON your_table_name;
-- DROP POLICY your_table_delete_own_org ON your_table_name;

-- USER SELF-ACCESS POLICY (for user preferences, profiles)
-- CREATE POLICY your_table_select_self
--     ON your_table_name
--     FOR SELECT
--     TO PUBLIC
--     USING (user_id = current_user_id());

-- JUNCTION TABLE WITH RELATED ENTITY CHECK
-- CREATE POLICY your_table_select_related
--     ON your_table_name
--     FOR SELECT
--     TO PUBLIC
--     USING (
--         EXISTS (
--             SELECT 1 FROM related_table
--             WHERE related_table.id = your_table_name.related_id
--             AND related_table.organization_id = current_user_organization_id()
--         )
--     );

-- ============================================================================
-- DATA MIGRATION (if updating existing table)
-- ============================================================================

-- Example: Add default values to existing records
-- UPDATE your_table_name
-- SET new_column = 'default_value'
-- WHERE new_column IS NULL;

-- ============================================================================
-- VALIDATION & TESTING
-- ============================================================================

-- Verify table was created
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'your_table_name') THEN
        RAISE EXCEPTION 'Table your_table_name was not created!';
    END IF;
END $$;

-- Verify RLS is enabled
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_tables
        WHERE tablename = 'your_table_name'
        AND rowsecurity = true
    ) THEN
        RAISE EXCEPTION 'RLS is not enabled on your_table_name!';
    END IF;
END $$;

-- Verify policies were created
DO $$
DECLARE
    policy_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO policy_count
    FROM pg_policies
    WHERE tablename = 'your_table_name';

    IF policy_count < 5 THEN
        RAISE WARNING 'Only % policies found. Expected at least 5 (super_admin + 4 org policies)', policy_count;
    END IF;
END $$;

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration VXXX completed successfully!';
    RAISE NOTICE 'Table: your_table_name';
    RAISE NOTICE 'RLS: Enabled with % policies', (SELECT COUNT(*) FROM pg_policies WHERE tablename = 'your_table_name');
    RAISE NOTICE '============================================';
END $$;

-- ============================================================================
-- ROLLBACK SCRIPT (Keep this commented, use only if needed)
-- ============================================================================

/*
BEGIN;

-- Drop policies first
DROP POLICY IF EXISTS your_table_delete_own_org ON your_table_name;
DROP POLICY IF EXISTS your_table_update_own_org ON your_table_name;
DROP POLICY IF EXISTS your_table_insert_own_org ON your_table_name;
DROP POLICY IF EXISTS your_table_select_own_org ON your_table_name;
DROP POLICY IF EXISTS your_table_super_admin_all ON your_table_name;

-- Drop table
DROP TABLE IF EXISTS your_table_name CASCADE;

COMMIT;
*/

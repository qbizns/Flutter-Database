-- =====================================================
-- Accounting Module Initialization
-- Description: Creates accounting schema and extensions
-- Usage: Run AFTER the main POS database is initialized
-- =====================================================

\echo 'Initializing Accounting Module...';

-- =====================================================
-- Create Accounting Schema
-- =====================================================

-- Create accounting schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS accounting;

\echo '✓ Accounting schema created';

-- =====================================================
-- Set Search Path
-- =====================================================

SET search_path TO accounting, public;

\echo '✓ Search path configured';

-- =====================================================
-- Enable Required Extensions (if not already enabled)
-- =====================================================

-- These extensions should already be enabled by the main POS init
-- But we'll ensure they're available just in case

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\echo '✓ Extensions verified';

-- =====================================================
-- Grant Permissions
-- =====================================================

-- Grant usage on accounting schema to appropriate roles
-- Note: Adjust role names based on your environment

-- Grant to postgres user (development)
GRANT USAGE ON SCHEMA accounting TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA accounting TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA accounting TO postgres;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA accounting TO postgres;

\echo '✓ Permissions granted';

-- =====================================================
-- Create Accounting Audit Function
-- =====================================================

CREATE OR REPLACE FUNCTION accounting.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION accounting.set_updated_at() IS 'Automatically updates the updated_at timestamp';

\echo '✓ Audit function created';

-- =====================================================
-- Initialization Complete
-- =====================================================

\echo '';
\echo '==========================================';
\echo 'Accounting Module Initialized';
\echo '==========================================';
\echo 'Schema: accounting';
\echo 'Status: Ready for migrations';
\echo '';
\echo 'Next Steps:';
\echo '  1. Run migrations: V001, V002';
\echo '  2. Load seed data: 001-005';
\echo '  3. Create report views';
\echo '';

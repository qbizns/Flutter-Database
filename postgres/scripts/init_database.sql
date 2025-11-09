-- ============================================================================
-- Database Initialization Script
-- Description: Sets up the database with required extensions and schemas
-- ============================================================================

-- Create database (run this separately as postgres superuser)
-- CREATE DATABASE pos_saas;

-- Connect to the database
\c pos_saas;

-- ============================================================================
-- EXTENSIONS
-- ============================================================================

-- UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Full text search
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Additional useful extensions
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- ============================================================================
-- SCHEMAS
-- ============================================================================

-- Public schema is already created, we'll use it for shared/system tables
-- Tenant-specific schemas will be created dynamically per organization

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Function to create tenant schema
CREATE OR REPLACE FUNCTION create_tenant_schema(tenant_slug VARCHAR)
RETURNS VOID AS $$
BEGIN
    EXECUTE format('CREATE SCHEMA IF NOT EXISTS %I', tenant_slug);
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- TYPES
-- ============================================================================

-- Enum for user status
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending');

-- Enum for organization status
CREATE TYPE organization_status AS ENUM ('trial', 'active', 'suspended', 'cancelled');

-- Enum for payment status
CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed', 'refunded', 'cancelled');

-- Enum for payment method
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'mobile_money', 'bank_transfer', 'other');

-- Enum for transaction type
CREATE TYPE transaction_type AS ENUM ('sale', 'return', 'exchange', 'void');

-- Enum for inventory transaction type
CREATE TYPE inventory_transaction_type AS ENUM ('purchase', 'sale', 'adjustment', 'return', 'transfer', 'waste');

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON SCHEMA public IS 'Shared schema for system-wide tables and multi-tenant management';
COMMENT ON FUNCTION update_updated_at_column() IS 'Trigger function to automatically update the updated_at timestamp';
COMMENT ON FUNCTION create_tenant_schema(VARCHAR) IS 'Creates a new schema for a tenant organization';

-- ============================================================================
-- COMPLETION MESSAGE
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Database initialization completed!';
    RAISE NOTICE 'Extensions, functions, and types created.';
    RAISE NOTICE 'Ready for migrations.';
    RAISE NOTICE '============================================';
END $$;

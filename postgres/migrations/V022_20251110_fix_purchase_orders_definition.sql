-- ============================================================================
-- Migration: V022 - Fix Purchase Orders Duplication
-- Description: Clean up duplicate purchase_orders definition and ensure all fields exist
-- Note: This migration is idempotent and safe to run multiple times
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- ANALYSIS: purchase_orders was defined in both V004 and V013
-- - V004_20251109_create_additional_pos_tables.sql (original, basic version)
-- - V013_20251110_create_advanced_pos_features.sql (extended version with procurement fields)
--
-- APPROACH: Ensure all fields from both definitions exist without dropping the table
-- ============================================================================

-- ============================================================================
-- ENSURE ALL FIELDS EXIST (from V004 original definition)
-- ============================================================================

-- Basic fields (most should already exist from V004)
ALTER TABLE purchase_orders
    ADD COLUMN IF NOT EXISTS id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ADD COLUMN IF NOT EXISTS organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS po_number VARCHAR(50) NOT NULL,
    ADD COLUMN IF NOT EXISTS supplier_id UUID NOT NULL REFERENCES suppliers(id),
    ADD COLUMN IF NOT EXISTS location_id UUID NOT NULL REFERENCES locations(id),
    ADD COLUMN IF NOT EXISTS order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    ADD COLUMN IF NOT EXISTS expected_delivery_date DATE,
    ADD COLUMN IF NOT EXISTS total_amount NUMERIC(12, 2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS notes TEXT;

-- ============================================================================
-- ADD FIELDS FROM V013 (Advanced POS Features)
-- ============================================================================

-- Advanced procurement fields that may not exist if V013 was blocked by duplication
ALTER TABLE purchase_orders
    ADD COLUMN IF NOT EXISTS reference_number VARCHAR(100),
    ADD COLUMN IF NOT EXISTS payment_terms VARCHAR(100),
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3),  -- Added in V017, but ensure it exists
    ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 8),  -- Added in V017
    ADD COLUMN IF NOT EXISTS incoterms VARCHAR(50),  -- International Commercial Terms
    ADD COLUMN IF NOT EXISTS shipping_method VARCHAR(100),
    ADD COLUMN IF NOT EXISTS shipping_cost NUMERIC(12, 2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tax_amount NUMERIC(12, 2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS discount_amount NUMERIC(12, 2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS grand_total NUMERIC(12, 2) DEFAULT 0;

-- Status field (ensure proper check constraint)
DO $$
BEGIN
    -- Add status column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'purchase_orders' AND column_name = 'status'
    ) THEN
        ALTER TABLE purchase_orders
            ADD COLUMN status VARCHAR(20) DEFAULT 'draft' CHECK (
                status IN ('draft', 'sent', 'confirmed', 'partially_received', 'received', 'cancelled')
            );
    END IF;
END $$;

-- Approval workflow fields
ALTER TABLE purchase_orders
    ADD COLUMN IF NOT EXISTS requires_approval BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS approved_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP WITH TIME ZONE;

-- Receiving tracking
ALTER TABLE purchase_orders
    ADD COLUMN IF NOT EXISTS received_amount NUMERIC(12, 2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS outstanding_amount NUMERIC(12, 2);

-- Metadata
ALTER TABLE purchase_orders
    ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

-- ============================================================================
-- ENSURE ALL STANDARD AUDIT FIELDS EXIST
-- ============================================================================

ALTER TABLE purchase_orders
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS updated_by UUID REFERENCES users(id);

-- ============================================================================
-- ENSURE ALL INDEXES EXIST
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_purchase_orders_organization_id ON purchase_orders(organization_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_po_number ON purchase_orders(po_number);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier_id ON purchase_orders(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_location_id ON purchase_orders(location_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_order_date ON purchase_orders(order_date);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_created_at ON purchase_orders(created_at);

-- ============================================================================
-- ENSURE UNIQUE CONSTRAINT EXISTS
-- ============================================================================

DO $$
BEGIN
    -- Add unique constraint if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'purchase_orders_organization_id_po_number_key'
    ) THEN
        ALTER TABLE purchase_orders
            ADD CONSTRAINT purchase_orders_organization_id_po_number_key
            UNIQUE (organization_id, po_number);
    END IF;
END $$;

-- ============================================================================
-- ENSURE TRIGGER EXISTS
-- ============================================================================

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'update_purchase_orders_updated_at'
    ) THEN
        CREATE TRIGGER update_purchase_orders_updated_at
            BEFORE UPDATE ON purchase_orders
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- ============================================================================
-- ADD COMMENTS
-- ============================================================================

COMMENT ON TABLE purchase_orders IS 'Purchase orders for procurement and inventory replenishment';
COMMENT ON COLUMN purchase_orders.po_number IS 'Unique purchase order number per organization';
COMMENT ON COLUMN purchase_orders.status IS 'PO status: draft, sent, confirmed, partially_received, received, cancelled';
COMMENT ON COLUMN purchase_orders.incoterms IS 'International Commercial Terms (FOB, CIF, EXW, etc.)';
COMMENT ON COLUMN purchase_orders.grand_total IS 'Final total including shipping, tax, and discounts';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V022 completed successfully!';
    RAISE NOTICE 'Purchase Orders definition consolidated:';
    RAISE NOTICE ' - Ensured all fields from V004 exist';
    RAISE NOTICE ' - Ensured all fields from V013 exist';
    RAISE NOTICE ' - Added missing currency/approval/receiving fields';
    RAISE NOTICE ' - Ensured all indexes exist';
    RAISE NOTICE ' - Ensured constraints and triggers exist';
    RAISE NOTICE 'Note: No data was lost in this migration';
    RAISE NOTICE '============================================';
END $$;

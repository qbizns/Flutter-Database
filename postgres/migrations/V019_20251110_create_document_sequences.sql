-- ============================================================================
-- Migration: V019 - Create Document Sequences
-- Description: Unified document numbering system with prefixes, padding, and reset frequency
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- DOCUMENT SEQUENCES TABLE
-- Description: Manage document number sequences per organization and document type
-- ============================================================================

CREATE TABLE IF NOT EXISTS document_sequences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Document Type
    document_type VARCHAR(100) NOT NULL,  -- 'sales', 'purchase_order', 'goods_receipt', 'invoice', 'payment', etc.

    -- Sequence Configuration
    prefix VARCHAR(20) DEFAULT '',         -- e.g., 'INV-', 'PO-', 'GR-'
    suffix VARCHAR(20) DEFAULT '',         -- e.g., '-2024', '-LOC1'
    next_number BIGINT NOT NULL DEFAULT 1,
    padding INTEGER DEFAULT 6,              -- Zero-padding (e.g., 6 -> '000001')
    increment_by INTEGER DEFAULT 1,

    -- Reset Configuration
    reset_frequency VARCHAR(20) DEFAULT 'never' CHECK (reset_frequency IN (
        'never',   -- Continuous numbering
        'daily',   -- Reset each day
        'monthly', -- Reset each month
        'yearly',  -- Reset each year
        'manual'   -- Manual reset only
    )),
    last_reset_at TIMESTAMP WITH TIME ZONE,
    last_reset_value BIGINT DEFAULT 0,

    -- Date Format in Sequence (optional)
    include_date BOOLEAN DEFAULT FALSE,
    date_format VARCHAR(50),  -- e.g., 'YYYY', 'YYYYMM', 'YYYYMMDD'

    -- Location/Branch Specific (optional)
    location_id UUID,  -- NULL = company-wide

    -- Configuration
    is_active BOOLEAN DEFAULT TRUE,
    allow_manual_override BOOLEAN DEFAULT FALSE,  -- Allow users to manually set document numbers

    -- Example Preview
    example_number VARCHAR(255),  -- Computed example: "INV-2024-000001"

    -- Notes
    description TEXT,
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id, document_type, location_id)
);

-- Indexes
CREATE INDEX idx_document_sequences_organization_id ON document_sequences(organization_id);
CREATE INDEX idx_document_sequences_document_type ON document_sequences(document_type);
CREATE INDEX idx_document_sequences_location_id ON document_sequences(location_id) WHERE location_id IS NOT NULL;
CREATE INDEX idx_document_sequences_active ON document_sequences(is_active) WHERE is_active = TRUE;

-- Trigger
CREATE TRIGGER update_document_sequences_updated_at
    BEFORE UPDATE ON document_sequences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE document_sequences IS 'Manages document numbering sequences per organization and document type';
COMMENT ON COLUMN document_sequences.document_type IS 'Type of document (sales, purchase_order, goods_receipt, invoice, etc.)';
COMMENT ON COLUMN document_sequences.next_number IS 'Next number to be issued';
COMMENT ON COLUMN document_sequences.padding IS 'Number of digits for zero-padding (6 = "000001")';
COMMENT ON COLUMN document_sequences.reset_frequency IS 'How often to reset the sequence (never, daily, monthly, yearly, manual)';
COMMENT ON COLUMN document_sequences.location_id IS 'Optional location-specific sequence (NULL = company-wide)';

-- ============================================================================
-- FUNCTION: Get Next Document Number
-- ============================================================================

CREATE OR REPLACE FUNCTION get_next_document_number(
    p_organization_id UUID,
    p_document_type VARCHAR,
    p_location_id UUID DEFAULT NULL,
    p_date DATE DEFAULT CURRENT_DATE
) RETURNS VARCHAR AS $$
DECLARE
    v_sequence_id UUID;
    v_prefix VARCHAR;
    v_suffix VARCHAR;
    v_next_number BIGINT;
    v_padding INTEGER;
    v_reset_frequency VARCHAR;
    v_last_reset_at TIMESTAMP;
    v_include_date BOOLEAN;
    v_date_format VARCHAR;
    v_should_reset BOOLEAN := FALSE;
    v_date_part VARCHAR;
    v_padded_number VARCHAR;
    v_document_number VARCHAR;
BEGIN
    -- Get sequence configuration (with row lock for concurrency)
    SELECT id, prefix, suffix, next_number, padding, reset_frequency, last_reset_at, include_date, date_format
    INTO v_sequence_id, v_prefix, v_suffix, v_next_number, v_padding, v_reset_frequency, v_last_reset_at, v_include_date, v_date_format
    FROM document_sequences
    WHERE organization_id = p_organization_id
      AND document_type = p_document_type
      AND (location_id = p_location_id OR (location_id IS NULL AND p_location_id IS NULL))
      AND is_active = TRUE
      AND deleted_at IS NULL
    FOR UPDATE  -- Lock row for concurrent access
    LIMIT 1;

    IF v_sequence_id IS NULL THEN
        RAISE EXCEPTION 'No active document sequence found for organization % and document type %', p_organization_id, p_document_type
            USING HINT = 'Create a document sequence in document_sequences table';
    END IF;

    -- Check if reset is needed
    IF v_reset_frequency = 'daily' AND (v_last_reset_at IS NULL OR DATE(v_last_reset_at) < p_date) THEN
        v_should_reset := TRUE;
    ELSIF v_reset_frequency = 'monthly' AND (v_last_reset_at IS NULL OR DATE_TRUNC('month', v_last_reset_at) < DATE_TRUNC('month', p_date)) THEN
        v_should_reset := TRUE;
    ELSIF v_reset_frequency = 'yearly' AND (v_last_reset_at IS NULL OR DATE_TRUNC('year', v_last_reset_at) < DATE_TRUNC('year', p_date)) THEN
        v_should_reset := TRUE;
    END IF;

    -- Reset if needed
    IF v_should_reset THEN
        v_next_number := 1;
        UPDATE document_sequences
        SET next_number = 2,  -- Increment for next time
            last_reset_at = CURRENT_TIMESTAMP,
            last_reset_value = v_next_number,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_sequence_id;
    ELSE
        -- Increment sequence
        UPDATE document_sequences
        SET next_number = next_number + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_sequence_id;
    END IF;

    -- Build date part if needed
    IF v_include_date AND v_date_format IS NOT NULL THEN
        v_date_part := TO_CHAR(p_date, v_date_format) || '-';
    ELSE
        v_date_part := '';
    END IF;

    -- Pad number
    v_padded_number := LPAD(v_next_number::TEXT, v_padding, '0');

    -- Build final document number
    v_document_number := COALESCE(v_prefix, '') || v_date_part || v_padded_number || COALESCE(v_suffix, '');

    RETURN v_document_number;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_next_document_number IS 'Thread-safe function to get next document number from sequence';

-- ============================================================================
-- FUNCTION: Preview Document Number (without incrementing)
-- ============================================================================

CREATE OR REPLACE FUNCTION preview_document_number(
    p_organization_id UUID,
    p_document_type VARCHAR,
    p_location_id UUID DEFAULT NULL,
    p_date DATE DEFAULT CURRENT_DATE
) RETURNS VARCHAR AS $$
DECLARE
    v_prefix VARCHAR;
    v_suffix VARCHAR;
    v_next_number BIGINT;
    v_padding INTEGER;
    v_include_date BOOLEAN;
    v_date_format VARCHAR;
    v_date_part VARCHAR;
    v_padded_number VARCHAR;
    v_preview VARCHAR;
BEGIN
    -- Get sequence configuration (no lock, read-only)
    SELECT prefix, suffix, next_number, padding, include_date, date_format
    INTO v_prefix, v_suffix, v_next_number, v_padding, v_include_date, v_date_format
    FROM document_sequences
    WHERE organization_id = p_organization_id
      AND document_type = p_document_type
      AND (location_id = p_location_id OR (location_id IS NULL AND p_location_id IS NULL))
      AND is_active = TRUE
      AND deleted_at IS NULL
    LIMIT 1;

    IF NOT FOUND THEN
        RETURN 'NO-SEQUENCE-CONFIGURED';
    END IF;

    -- Build date part if needed
    IF v_include_date AND v_date_format IS NOT NULL THEN
        v_date_part := TO_CHAR(p_date, v_date_format) || '-';
    ELSE
        v_date_part := '';
    END IF;

    -- Pad number
    v_padded_number := LPAD(v_next_number::TEXT, v_padding, '0');

    -- Build preview
    v_preview := COALESCE(v_prefix, '') || v_date_part || v_padded_number || COALESCE(v_suffix, '');

    RETURN v_preview;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION preview_document_number IS 'Preview next document number without incrementing sequence';

-- ============================================================================
-- FUNCTION: Reset Document Sequence
-- ============================================================================

CREATE OR REPLACE FUNCTION reset_document_sequence(
    p_organization_id UUID,
    p_document_type VARCHAR,
    p_reset_to BIGINT DEFAULT 1,
    p_location_id UUID DEFAULT NULL
) RETURNS JSONB AS $$
DECLARE
    v_sequence_id UUID;
    v_old_next_number BIGINT;
    v_result JSONB;
BEGIN
    -- Get sequence
    SELECT id, next_number INTO v_sequence_id, v_old_next_number
    FROM document_sequences
    WHERE organization_id = p_organization_id
      AND document_type = p_document_type
      AND (location_id = p_location_id OR (location_id IS NULL AND p_location_id IS NULL))
      AND deleted_at IS NULL
    FOR UPDATE;

    IF v_sequence_id IS NULL THEN
        RAISE EXCEPTION 'Document sequence not found for organization % and document type %', p_organization_id, p_document_type;
    END IF;

    -- Reset sequence
    UPDATE document_sequences
    SET next_number = p_reset_to,
        last_reset_at = CURRENT_TIMESTAMP,
        last_reset_value = v_old_next_number,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_sequence_id;

    v_result := jsonb_build_object(
        'success', TRUE,
        'sequence_id', v_sequence_id,
        'document_type', p_document_type,
        'old_next_number', v_old_next_number,
        'new_next_number', p_reset_to,
        'reset_at', CURRENT_TIMESTAMP
    );

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION reset_document_sequence IS 'Manually reset a document sequence to a specific number';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

ALTER TABLE document_sequences ENABLE ROW LEVEL SECURITY;

CREATE POLICY document_sequences_tenant_isolation ON document_sequences
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY document_sequences_tenant_isolation ON document_sequences
    IS 'Ensure users can only access document sequences for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V019 completed successfully!';
    RAISE NOTICE 'Document Sequences created:';
    RAISE NOTICE ' - document_sequences table';
    RAISE NOTICE ' - get_next_document_number() function';
    RAISE NOTICE ' - preview_document_number() function';
    RAISE NOTICE ' - reset_document_sequence() function';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

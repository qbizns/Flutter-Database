-- ============================================================================
-- Migration: V014 - Create E-Invoicing Integration (ZATCA & ETA)
-- Description: Implements Saudi Arabia (ZATCA) and Egypt (ETA) e-invoicing support
--              with authority-agnostic design for extensibility
-- Note: Non-destructive, additive changes only
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- E-INVOICING DOCUMENTS TABLE
-- Description: Central table for all e-invoicing documents (ZATCA, ETA, future authorities)
-- ============================================================================

CREATE TABLE IF NOT EXISTS e_invoicing_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Generic Source Link (polymorphic reference)
    source_table VARCHAR(100) NOT NULL, -- e.g., 'sales', 'customer_invoices'
    source_id UUID NOT NULL, -- ID of the source document

    -- Authority Information
    authority VARCHAR(20) NOT NULL CHECK (authority IN ('ZATCA', 'ETA', 'OTHER')),
    country_code VARCHAR(3) NOT NULL, -- ISO 3166-1 alpha-3 (SAU, EGY, etc.)

    -- Document Identification
    document_uuid UUID NOT NULL DEFAULT gen_random_uuid(), -- Universal unique ID
    document_type VARCHAR(50) NOT NULL, -- 'standard', 'simplified', 'credit_note', 'debit_note'
    document_number VARCHAR(100) NOT NULL, -- Authority-specific invoice number format
    internal_reference VARCHAR(100), -- Link to original sale_number

    -- Status Workflow
    status VARCHAR(30) NOT NULL DEFAULT 'draft' CHECK (status IN (
        'draft',           -- Being prepared
        'pending',         -- Ready to submit
        'submitted',       -- Sent to authority
        'accepted',        -- Successfully validated
        'rejected',        -- Failed validation
        'cancelled',       -- Cancelled
        'error'            -- Technical error
    )),

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE,
    response_at TIMESTAMP WITH TIME ZONE,

    -- Payloads (stored for audit and resubmission)
    request_payload JSONB, -- What we send
    response_payload JSONB, -- What we receive

    -- Error Handling
    error_code VARCHAR(50),
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    last_retry_at TIMESTAMP WITH TIME ZONE,

    -- ZATCA-Specific Fields
    zatca_hash_value VARCHAR(255), -- SHA256 hash of invoice
    zatca_previous_hash_value VARCHAR(255), -- Chain to previous invoice
    zatca_invoice_counter_value INTEGER, -- Sequential counter (ICV)
    zatca_cryptographic_stamp TEXT, -- Digital signature
    zatca_qr_code_payload TEXT, -- TLV-encoded QR code
    zatca_compliance_invoice_number VARCHAR(100), -- PIH format

    -- ETA-Specific Fields
    eta_document_type_version VARCHAR(10), -- Document type version
    eta_submission_uuid UUID, -- ETA's UUID for this submission
    eta_long_id VARCHAR(255), -- ETA's long ID after acceptance
    eta_internal_id VARCHAR(100), -- Company's internal ID
    eta_digital_signature TEXT, -- CAdES-BES signature
    eta_signature_algorithm VARCHAR(50) DEFAULT 'CAdES-BES',

    -- Additional Metadata
    submission_format VARCHAR(20) CHECK (submission_format IN ('XML', 'JSON', 'UBL')),
    metadata JSONB DEFAULT '{}', -- Extensible field for authority-specific data

    -- Audit Fields
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id, authority, document_number),
    CONSTRAINT valid_source_reference CHECK (source_table <> '' AND source_id IS NOT NULL)
);

-- Indexes
CREATE INDEX idx_e_invoicing_documents_organization_id ON e_invoicing_documents(organization_id);
CREATE INDEX idx_e_invoicing_documents_source ON e_invoicing_documents(source_table, source_id);
CREATE INDEX idx_e_invoicing_documents_authority ON e_invoicing_documents(authority);
CREATE INDEX idx_e_invoicing_documents_status ON e_invoicing_documents(status);
CREATE INDEX idx_e_invoicing_documents_document_uuid ON e_invoicing_documents(document_uuid);
CREATE INDEX idx_e_invoicing_documents_document_number ON e_invoicing_documents(document_number);
CREATE INDEX idx_e_invoicing_documents_created_at ON e_invoicing_documents(created_at);
CREATE INDEX idx_e_invoicing_documents_submitted_at ON e_invoicing_documents(submitted_at) WHERE submitted_at IS NOT NULL;

-- Trigger
CREATE TRIGGER update_e_invoicing_documents_updated_at
    BEFORE UPDATE ON e_invoicing_documents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE e_invoicing_documents IS 'Central repository for all e-invoicing documents across authorities (ZATCA, ETA, etc.)';
COMMENT ON COLUMN e_invoicing_documents.source_table IS 'Polymorphic reference to source document table (sales, customer_invoices, etc.)';
COMMENT ON COLUMN e_invoicing_documents.source_id IS 'ID of the source document in the referenced table';
COMMENT ON COLUMN e_invoicing_documents.zatca_invoice_counter_value IS 'ZATCA Invoice Counter Value (ICV) - sequential across all invoices';
COMMENT ON COLUMN e_invoicing_documents.zatca_previous_hash_value IS 'Hash of previous invoice for blockchain-style chaining';
COMMENT ON COLUMN e_invoicing_documents.eta_long_id IS 'ETA-assigned long ID returned after successful acceptance';

-- ============================================================================
-- E-INVOICING DOCUMENT EVENTS TABLE
-- Description: Audit trail for all state changes and API interactions
-- ============================================================================

CREATE TABLE IF NOT EXISTS e_invoicing_document_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Document Reference
    e_invoicing_document_id UUID NOT NULL REFERENCES e_invoicing_documents(id) ON DELETE CASCADE,

    -- Event Details
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN (
        'created',         -- Document created
        'validated',       -- Local validation passed
        'submitted',       -- Sent to authority
        'accepted',        -- Authority accepted
        'rejected',        -- Authority rejected
        'cancelled',       -- Document cancelled
        'error',           -- Technical error occurred
        'retry',           -- Retry attempt
        'status_check'     -- Status verification
    )),

    event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Previous and New State
    previous_status VARCHAR(30),
    new_status VARCHAR(30),

    -- Event Details
    event_description TEXT,
    event_data JSONB DEFAULT '{}', -- Additional event context

    -- HTTP/API Details (if applicable)
    http_status_code INTEGER,
    http_method VARCHAR(10),
    api_endpoint TEXT,
    request_headers JSONB,
    response_headers JSONB,

    -- Error Details
    error_code VARCHAR(50),
    error_message TEXT,
    error_details JSONB,

    -- User/System Context
    triggered_by VARCHAR(50) NOT NULL DEFAULT 'system' CHECK (triggered_by IN ('user', 'system', 'scheduled', 'webhook')),
    user_id UUID REFERENCES users(id),

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexes
CREATE INDEX idx_e_invoicing_document_events_organization_id ON e_invoicing_document_events(organization_id);
CREATE INDEX idx_e_invoicing_document_events_document_id ON e_invoicing_document_events(e_invoicing_document_id);
CREATE INDEX idx_e_invoicing_document_events_event_type ON e_invoicing_document_events(event_type);
CREATE INDEX idx_e_invoicing_document_events_event_timestamp ON e_invoicing_document_events(event_timestamp);
CREATE INDEX idx_e_invoicing_document_events_new_status ON e_invoicing_document_events(new_status);

-- Comments
COMMENT ON TABLE e_invoicing_document_events IS 'Complete audit trail of all e-invoicing document events and state changes';
COMMENT ON COLUMN e_invoicing_document_events.event_data IS 'Additional context data for the event (flexible JSONB)';
COMMENT ON COLUMN e_invoicing_document_events.triggered_by IS 'Source of the event: user action, system automation, scheduled job, or webhook';

-- ============================================================================
-- EXTEND CUSTOMERS TABLE - E-INVOICING FIELDS
-- Description: Add e-invoicing specific customer fields (non-destructive)
-- ============================================================================

-- Tax Registration Fields
ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS tax_registration_number VARCHAR(100),
    ADD COLUMN IF NOT EXISTS tax_registration_type VARCHAR(50), -- 'TIN', 'VAT', 'CRN', 'MOMRAH', 'MISA', '700', 'SAG', 'NAT', 'GCC', 'IQA', 'OTHER'
    ADD COLUMN IF NOT EXISTS tax_registration_country VARCHAR(3), -- ISO 3166-1 alpha-3
    ADD COLUMN IF NOT EXISTS is_tax_registered BOOLEAN DEFAULT false;

-- ZATCA-Specific Customer Fields
ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS zatca_building_number VARCHAR(20),
    ADD COLUMN IF NOT EXISTS zatca_street_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS zatca_district VARCHAR(100),
    ADD COLUMN IF NOT EXISTS zatca_city_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS zatca_postal_zone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS zatca_country_subentity VARCHAR(100), -- Province/State
    ADD COLUMN IF NOT EXISTS zatca_country_code VARCHAR(3) DEFAULT 'SAU'; -- ISO 3166-1 alpha-3

-- ETA-Specific Customer Fields
ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS eta_receiver_type VARCHAR(50), -- 'B' (Business), 'P' (Person), 'F' (Foreigner)
    ADD COLUMN IF NOT EXISTS eta_receiver_id VARCHAR(100), -- National ID, Passport, Unified ID
    ADD COLUMN IF NOT EXISTS eta_receiver_id_type VARCHAR(50), -- 'NationalID', 'Passport', 'Unified', 'IQA'
    ADD COLUMN IF NOT EXISTS eta_governorate VARCHAR(100),
    ADD COLUMN IF NOT EXISTS eta_region_city VARCHAR(100);

-- Structured Address JSON (for flexibility across authorities)
ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS structured_address JSONB DEFAULT '{}';

-- Create indexes for new fields
CREATE INDEX IF NOT EXISTS idx_customers_tax_registration_number ON customers(tax_registration_number) WHERE tax_registration_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_customers_is_tax_registered ON customers(is_tax_registered) WHERE is_tax_registered = true;

-- Comments
COMMENT ON COLUMN customers.tax_registration_number IS 'Customer tax registration/VAT number for e-invoicing';
COMMENT ON COLUMN customers.tax_registration_type IS 'Type of tax registration (TIN, VAT, CRN, etc.) - varies by country';
COMMENT ON COLUMN customers.zatca_building_number IS 'ZATCA-required building number for Saudi addresses';
COMMENT ON COLUMN customers.eta_receiver_type IS 'ETA receiver classification: Business, Person, or Foreigner';
COMMENT ON COLUMN customers.structured_address IS 'Flexible JSONB field for authority-specific address requirements';

-- ============================================================================
-- EXTEND PRODUCTS TABLE - E-INVOICING FIELDS
-- Description: Add e-invoicing specific product fields (non-destructive)
-- ============================================================================

-- Standard Item Codes (GS1, EGS, etc.)
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS standard_item_code VARCHAR(100),
    ADD COLUMN IF NOT EXISTS standard_item_code_type VARCHAR(50), -- 'GS1', 'EGS', 'GTIN', 'UNSPSC', 'HS_CODE'
    ADD COLUMN IF NOT EXISTS harmonized_system_code VARCHAR(50), -- HS Code for customs/tax
    ADD COLUMN IF NOT EXISTS customs_tariff_code VARCHAR(50);

-- Standard Unit of Measure Codes
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS standard_uom_code VARCHAR(10), -- UN/ECE Recommendation 20 codes
    ADD COLUMN IF NOT EXISTS standard_uom_name VARCHAR(100); -- 'PCE' (Piece), 'KGM' (Kilogram), 'MTR' (Meter)

-- Tax Classification
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS tax_category_code VARCHAR(50), -- 'S' (Standard), 'Z' (Zero), 'E' (Exempt), 'O' (Out of scope)
    ADD COLUMN IF NOT EXISTS default_vat_rate NUMERIC(5, 2), -- Default VAT percentage for this product
    ADD COLUMN IF NOT EXISTS tax_exemption_reason VARCHAR(255),
    ADD COLUMN IF NOT EXISTS tax_exemption_code VARCHAR(50);

-- ZATCA-Specific Product Fields
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS zatca_item_classification VARCHAR(50), -- Product classification code
    ADD COLUMN IF NOT EXISTS zatca_is_excise_taxable BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS zatca_excise_tax_rate NUMERIC(5, 2) DEFAULT 0;

-- ETA-Specific Product Fields
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS eta_gs1_code VARCHAR(100), -- GS1 barcode for ETA
    ADD COLUMN IF NOT EXISTS eta_egs_code VARCHAR(100), -- Egyptian EGS code
    ADD COLUMN IF NOT EXISTS eta_item_type VARCHAR(50), -- 'GS1', 'EGS'
    ADD COLUMN IF NOT EXISTS eta_internal_code VARCHAR(100); -- Company internal code

-- Product E-Invoicing Metadata
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS e_invoicing_metadata JSONB DEFAULT '{}';

-- Create indexes for new fields
CREATE INDEX IF NOT EXISTS idx_products_standard_item_code ON products(standard_item_code) WHERE standard_item_code IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_tax_category_code ON products(tax_category_code) WHERE tax_category_code IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_harmonized_system_code ON products(harmonized_system_code) WHERE harmonized_system_code IS NOT NULL;

-- Comments
COMMENT ON COLUMN products.standard_item_code IS 'Standardized item code (GS1, EGS, etc.) for e-invoicing';
COMMENT ON COLUMN products.standard_item_code_type IS 'Type of standard item code (GS1, EGS, GTIN, UNSPSC, HS_CODE)';
COMMENT ON COLUMN products.standard_uom_code IS 'UN/ECE Recommendation 20 standard unit of measure code';
COMMENT ON COLUMN products.tax_category_code IS 'Tax category code for e-invoicing (S=Standard, Z=Zero, E=Exempt, O=Out of scope)';
COMMENT ON COLUMN products.zatca_is_excise_taxable IS 'Whether product is subject to ZATCA excise tax (tobacco, soft drinks, etc.)';
COMMENT ON COLUMN products.eta_gs1_code IS 'GS1 Global Trade Item Number (GTIN) for ETA e-invoicing';

-- ============================================================================
-- EXTEND SALES TABLE - E-INVOICING FIELDS (OPTIONAL)
-- Description: Add e-invoicing tracking to sales transactions
-- ============================================================================

-- E-Invoicing Status on Sales
ALTER TABLE sales
    ADD COLUMN IF NOT EXISTS is_e_invoice_required BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS e_invoice_status VARCHAR(30) DEFAULT 'not_applicable' CHECK (e_invoice_status IN (
        'not_applicable',
        'pending',
        'submitted',
        'accepted',
        'rejected',
        'cancelled'
    )),
    ADD COLUMN IF NOT EXISTS e_invoice_document_id UUID REFERENCES e_invoicing_documents(id);

-- E-Invoicing Metadata on Sales
ALTER TABLE sales
    ADD COLUMN IF NOT EXISTS e_invoicing_metadata JSONB DEFAULT '{}';

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sales_is_e_invoice_required ON sales(is_e_invoice_required) WHERE is_e_invoice_required = true;
CREATE INDEX IF NOT EXISTS idx_sales_e_invoice_status ON sales(e_invoice_status) WHERE e_invoice_status != 'not_applicable';
CREATE INDEX IF NOT EXISTS idx_sales_e_invoice_document_id ON sales(e_invoice_document_id) WHERE e_invoice_document_id IS NOT NULL;

-- Comments
COMMENT ON COLUMN sales.is_e_invoice_required IS 'Whether this sale requires e-invoicing submission';
COMMENT ON COLUMN sales.e_invoice_status IS 'Current e-invoicing status for this sale';
COMMENT ON COLUMN sales.e_invoice_document_id IS 'Reference to the e-invoicing document if created';

-- ============================================================================
-- HELPER VIEWS FOR E-INVOICING
-- ============================================================================

-- View: E-Invoices with Source Details
CREATE OR REPLACE VIEW view_e_invoices_with_source AS
SELECT
    ed.id,
    ed.organization_id,
    ed.authority,
    ed.country_code,
    ed.document_uuid,
    ed.document_type,
    ed.document_number,
    ed.internal_reference,
    ed.status,
    ed.created_at,
    ed.submitted_at,
    ed.response_at,
    ed.source_table,
    ed.source_id,
    -- Attempt to join with sales table (if source is sales)
    CASE
        WHEN ed.source_table = 'sales' THEN s.sale_number
        ELSE NULL
    END as sale_number,
    CASE
        WHEN ed.source_table = 'sales' THEN s.total_amount
        ELSE NULL
    END as sale_total_amount,
    CASE
        WHEN ed.source_table = 'sales' THEN s.customer_id
        ELSE NULL
    END as customer_id
FROM e_invoicing_documents ed
LEFT JOIN sales s ON ed.source_table = 'sales' AND ed.source_id = s.id
WHERE ed.deleted_at IS NULL;

COMMENT ON VIEW view_e_invoices_with_source IS 'E-invoicing documents with resolved source details';

-- View: E-Invoicing Document Event History
CREATE OR REPLACE VIEW view_e_invoice_event_history AS
SELECT
    ed.organization_id,
    ed.document_number,
    ed.authority,
    ed.status as current_status,
    ev.event_type,
    ev.event_timestamp,
    ev.previous_status,
    ev.new_status,
    ev.event_description,
    ev.error_code,
    ev.error_message,
    ev.triggered_by,
    u.email as user_email
FROM e_invoicing_documents ed
INNER JOIN e_invoicing_document_events ev ON ed.id = ev.e_invoicing_document_id
LEFT JOIN users u ON ev.user_id = u.id
WHERE ed.deleted_at IS NULL
ORDER BY ev.event_timestamp DESC;

COMMENT ON VIEW view_e_invoice_event_history IS 'Complete event history for e-invoicing documents with user details';

-- View: ZATCA Invoices Summary
CREATE OR REPLACE VIEW view_zatca_invoices AS
SELECT
    ed.id,
    ed.organization_id,
    ed.document_uuid,
    ed.document_number,
    ed.internal_reference,
    ed.status,
    ed.zatca_invoice_counter_value,
    ed.zatca_hash_value,
    ed.zatca_previous_hash_value,
    ed.zatca_qr_code_payload,
    ed.created_at,
    ed.submitted_at,
    ed.response_at,
    ed.error_code,
    ed.error_message
FROM e_invoicing_documents ed
WHERE ed.authority = 'ZATCA'
  AND ed.deleted_at IS NULL
ORDER BY ed.zatca_invoice_counter_value DESC;

COMMENT ON VIEW view_zatca_invoices IS 'ZATCA-specific e-invoicing documents with ZATCA fields';

-- View: ETA Invoices Summary
CREATE OR REPLACE VIEW view_eta_invoices AS
SELECT
    ed.id,
    ed.organization_id,
    ed.document_uuid,
    ed.eta_submission_uuid,
    ed.eta_long_id,
    ed.document_number,
    ed.eta_internal_id,
    ed.document_type,
    ed.eta_document_type_version,
    ed.status,
    ed.created_at,
    ed.submitted_at,
    ed.response_at,
    ed.error_code,
    ed.error_message
FROM e_invoicing_documents ed
WHERE ed.authority = 'ETA'
  AND ed.deleted_at IS NULL
ORDER BY ed.created_at DESC;

COMMENT ON VIEW view_eta_invoices IS 'ETA-specific e-invoicing documents with ETA fields';

-- View: Failed E-Invoices Requiring Attention
CREATE OR REPLACE VIEW view_failed_e_invoices AS
SELECT
    ed.id,
    ed.organization_id,
    ed.authority,
    ed.document_number,
    ed.internal_reference,
    ed.status,
    ed.error_code,
    ed.error_message,
    ed.retry_count,
    ed.last_retry_at,
    ed.created_at,
    -- Calculate if retry is needed
    CASE
        WHEN ed.retry_count < 3 AND (ed.last_retry_at IS NULL OR ed.last_retry_at < CURRENT_TIMESTAMP - INTERVAL '1 hour')
        THEN true
        ELSE false
    END as should_retry,
    -- Latest event
    (SELECT ev.event_description
     FROM e_invoicing_document_events ev
     WHERE ev.e_invoicing_document_id = ed.id
     ORDER BY ev.event_timestamp DESC
     LIMIT 1) as latest_event
FROM e_invoicing_documents ed
WHERE ed.status IN ('rejected', 'error')
  AND ed.deleted_at IS NULL
ORDER BY ed.created_at DESC;

COMMENT ON VIEW view_failed_e_invoices IS 'E-invoices that failed submission and may need retry or investigation';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS
ALTER TABLE e_invoicing_documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE e_invoicing_document_events ENABLE ROW LEVEL SECURITY;

-- RLS Policy for e_invoicing_documents
CREATE POLICY e_invoicing_documents_tenant_isolation ON e_invoicing_documents
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

-- RLS Policy for e_invoicing_document_events
CREATE POLICY e_invoicing_document_events_tenant_isolation ON e_invoicing_document_events
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY e_invoicing_documents_tenant_isolation ON e_invoicing_documents
    IS 'Ensure users can only access e-invoicing documents for their organizations';
COMMENT ON POLICY e_invoicing_document_events_tenant_isolation ON e_invoicing_document_events
    IS 'Ensure users can only access e-invoicing events for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V014 completed successfully!';
    RAISE NOTICE 'E-Invoicing Integration (ZATCA & ETA) created.';
    RAISE NOTICE ' - 2 new tables created';
    RAISE NOTICE ' - Extended customers table with e-invoicing fields';
    RAISE NOTICE ' - Extended products table with e-invoicing fields';
    RAISE NOTICE ' - Extended sales table with e-invoicing status';
    RAISE NOTICE ' - 5 helper views created';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;

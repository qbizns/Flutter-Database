-- ============================================================================
-- Seed Data: 007 - E-Invoicing Integration Data
-- Description: Sample data for ZATCA (Saudi Arabia) and ETA (Egypt) e-invoicing
-- Dependencies: Requires 001 (core), 002 (pos), and V014 migration
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- GET ORGANIZATION AND USER IDS
-- ============================================================================

DO $$
DECLARE
    v_org_id UUID;
    v_user_id UUID;
    v_customer_id_1 UUID;
    v_customer_id_2 UUID;
    v_customer_id_3 UUID;
    v_product_id_1 UUID;
    v_product_id_2 UUID;
    v_product_id_3 UUID;
    v_sale_id_1 UUID;
    v_sale_id_2 UUID;
    v_sale_id_3 UUID;
    v_sale_id_4 UUID;
    v_einv_doc_id_1 UUID;
    v_einv_doc_id_2 UUID;
    v_einv_doc_id_3 UUID;
    v_einv_doc_id_4 UUID;
BEGIN
    -- Get the first organization
    SELECT id INTO v_org_id FROM organizations LIMIT 1;

    -- Get the first user
    SELECT id INTO v_user_id FROM users LIMIT 1;

    -- If no organization exists, raise error
    IF v_org_id IS NULL THEN
        RAISE EXCEPTION 'No organization found. Please run 001_seed_core_data.sql first.';
    END IF;

    RAISE NOTICE 'Using Organization ID: %', v_org_id;
    RAISE NOTICE 'Using User ID: %', v_user_id;

    -- ============================================================================
    -- UPDATE EXISTING CUSTOMERS WITH E-INVOICING DATA
    -- ============================================================================

    RAISE NOTICE 'Updating customers with e-invoicing data...';

    -- Customer 1: Saudi Arabian B2B Customer
    SELECT id INTO v_customer_id_1 FROM customers WHERE organization_id = v_org_id LIMIT 1 OFFSET 0;

    IF v_customer_id_1 IS NOT NULL THEN
        UPDATE customers
        SET
            -- Tax Registration
            tax_registration_number = '300000000000003',
            tax_registration_type = 'TIN',
            tax_registration_country = 'SAU',
            is_tax_registered = true,

            -- ZATCA Fields
            zatca_building_number = '1234',
            zatca_street_name = 'King Fahd Road',
            zatca_district = 'Al Olaya',
            zatca_city_name = 'Riyadh',
            zatca_postal_zone = '12214',
            zatca_country_subentity = 'Riyadh Region',
            zatca_country_code = 'SAU',

            -- Structured Address
            structured_address = jsonb_build_object(
                'building_number', '1234',
                'street', 'King Fahd Road',
                'district', 'Al Olaya',
                'city', 'Riyadh',
                'postal_code', '12214',
                'country', 'SAU',
                'additional_number', '5678'
            ),

            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_customer_id_1;

        RAISE NOTICE 'Updated customer % with ZATCA data', v_customer_id_1;
    END IF;

    -- Customer 2: Egyptian B2B Customer
    SELECT id INTO v_customer_id_2 FROM customers WHERE organization_id = v_org_id LIMIT 1 OFFSET 1;

    IF v_customer_id_2 IS NOT NULL THEN
        UPDATE customers
        SET
            -- Tax Registration
            tax_registration_number = '100-200-300',
            tax_registration_type = 'VAT',
            tax_registration_country = 'EGY',
            is_tax_registered = true,

            -- ETA Fields
            eta_receiver_type = 'B',  -- Business
            eta_receiver_id = '100200300',
            eta_receiver_id_type = 'TIN',
            eta_governorate = 'Cairo',
            eta_region_city = 'Nasr City',

            -- Structured Address
            structured_address = jsonb_build_object(
                'country', 'EG',
                'governorate', 'Cairo',
                'regionCity', 'Nasr City',
                'street', 'Abbas El Akkad Street',
                'buildingNumber', '45',
                'postalCode', '11371',
                'floor', '3',
                'room', '301'
            ),

            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_customer_id_2;

        RAISE NOTICE 'Updated customer % with ETA data', v_customer_id_2;
    END IF;

    -- Customer 3: Saudi Arabian B2C Customer (Simplified Invoice)
    SELECT id INTO v_customer_id_3 FROM customers WHERE organization_id = v_org_id LIMIT 1 OFFSET 2;

    IF v_customer_id_3 IS NOT NULL THEN
        UPDATE customers
        SET
            -- Tax Registration (B2C may not have)
            tax_registration_number = NULL,
            tax_registration_type = NULL,
            tax_registration_country = 'SAU',
            is_tax_registered = false,

            -- ZATCA Fields (minimal for B2C)
            zatca_city_name = 'Jeddah',
            zatca_country_code = 'SAU',

            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_customer_id_3;

        RAISE NOTICE 'Updated customer % with ZATCA B2C data', v_customer_id_3;
    END IF;

    -- ============================================================================
    -- UPDATE EXISTING PRODUCTS WITH E-INVOICING DATA
    -- ============================================================================

    RAISE NOTICE 'Updating products with e-invoicing data...';

    -- Product 1: Standard VAT Product
    SELECT id INTO v_product_id_1 FROM products WHERE organization_id = v_org_id LIMIT 1 OFFSET 0;

    IF v_product_id_1 IS NOT NULL THEN
        UPDATE products
        SET
            -- Standard Item Codes
            standard_item_code = '6922266980129',
            standard_item_code_type = 'GS1',
            harmonized_system_code = '8471.30.00.00',

            -- Standard UOM
            standard_uom_code = 'PCE',
            standard_uom_name = 'Piece',

            -- Tax Classification
            tax_category_code = 'S',  -- Standard rated
            default_vat_rate = 15.00,

            -- ZATCA Fields
            zatca_item_classification = 'Electronics',
            zatca_is_excise_taxable = false,
            zatca_excise_tax_rate = 0,

            -- ETA Fields
            eta_gs1_code = '6922266980129',
            eta_item_type = 'GS1',
            eta_internal_code = 'PROD-001',

            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_product_id_1;

        RAISE NOTICE 'Updated product % with e-invoicing data', v_product_id_1;
    END IF;

    -- Product 2: Zero-rated Product
    SELECT id INTO v_product_id_2 FROM products WHERE organization_id = v_org_id LIMIT 1 OFFSET 1;

    IF v_product_id_2 IS NOT NULL THEN
        UPDATE products
        SET
            -- Standard Item Codes
            standard_item_code = '5901234123457',
            standard_item_code_type = 'GTIN',
            harmonized_system_code = '1001.90.00.00',

            -- Standard UOM
            standard_uom_code = 'KGM',
            standard_uom_name = 'Kilogram',

            -- Tax Classification
            tax_category_code = 'Z',  -- Zero rated
            default_vat_rate = 0.00,
            tax_exemption_reason = 'Export - Zero rated',
            tax_exemption_code = 'VATEX-SA-32',

            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_product_id_2;

        RAISE NOTICE 'Updated product % with zero-rated data', v_product_id_2;
    END IF;

    -- Product 3: Excise taxable product (Soft Drink)
    SELECT id INTO v_product_id_3 FROM products WHERE organization_id = v_org_id LIMIT 1 OFFSET 2;

    IF v_product_id_3 IS NOT NULL THEN
        UPDATE products
        SET
            -- Standard Item Codes
            standard_item_code = '5449000000996',
            standard_item_code_type = 'GS1',

            -- Standard UOM
            standard_uom_code = 'BX',
            standard_uom_name = 'Box',

            -- Tax Classification
            tax_category_code = 'S',  -- Standard rated
            default_vat_rate = 15.00,

            -- ZATCA Fields (Excise Tax on Soft Drinks)
            zatca_item_classification = 'Beverages',
            zatca_is_excise_taxable = true,
            zatca_excise_tax_rate = 50.00,  -- 50% excise on soft drinks in Saudi

            updated_at = CURRENT_TIMESTAMP
        WHERE id = v_product_id_3;

        RAISE NOTICE 'Updated product % with excise tax data', v_product_id_3;
    END IF;

    -- ============================================================================
    -- CREATE SALES WITH E-INVOICING ENABLED
    -- ============================================================================

    RAISE NOTICE 'Creating sales with e-invoicing enabled...';

    -- Sale 1: ZATCA Standard B2B Invoice
    v_sale_id_1 := gen_random_uuid();
    INSERT INTO sales (
        id,
        organization_id,
        sale_number,
        reference_number,
        transaction_type,
        customer_id,
        cashier_id,
        subtotal,
        tax_amount,
        discount_amount,
        total_amount,
        paid_amount,
        change_amount,
        outstanding_amount,
        payment_status,
        transaction_date,
        completed_at,
        is_e_invoice_required,
        e_invoice_status,
        created_by,
        updated_by
    ) VALUES (
        v_sale_id_1,
        v_org_id,
        'INV-2024-001',
        'REF-ZATCA-001',
        'sale',
        v_customer_id_1,
        v_user_id,
        1000.00,
        150.00,
        0.00,
        1150.00,
        1150.00,
        0.00,
        0.00,
        'paid',
        '2024-11-01 10:30:00+03',
        '2024-11-01 10:30:00+03',
        true,
        'pending',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created sale % for ZATCA B2B invoice', v_sale_id_1;

    -- Sale 2: ZATCA Simplified B2C Invoice
    v_sale_id_2 := gen_random_uuid();
    INSERT INTO sales (
        id,
        organization_id,
        sale_number,
        reference_number,
        transaction_type,
        customer_id,
        cashier_id,
        subtotal,
        tax_amount,
        discount_amount,
        total_amount,
        paid_amount,
        change_amount,
        outstanding_amount,
        payment_status,
        transaction_date,
        completed_at,
        is_e_invoice_required,
        e_invoice_status,
        created_by,
        updated_by
    ) VALUES (
        v_sale_id_2,
        v_org_id,
        'INV-2024-002',
        'REF-ZATCA-002',
        'sale',
        v_customer_id_3,
        v_user_id,
        250.00,
        37.50,
        0.00,
        287.50,
        287.50,
        0.00,
        0.00,
        'paid',
        '2024-11-01 14:15:00+03',
        '2024-11-01 14:15:00+03',
        true,
        'pending',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created sale % for ZATCA B2C simplified invoice', v_sale_id_2;

    -- Sale 3: ETA Egyptian Invoice
    v_sale_id_3 := gen_random_uuid();
    INSERT INTO sales (
        id,
        organization_id,
        sale_number,
        reference_number,
        transaction_type,
        customer_id,
        cashier_id,
        subtotal,
        tax_amount,
        discount_amount,
        total_amount,
        paid_amount,
        change_amount,
        outstanding_amount,
        payment_status,
        transaction_date,
        completed_at,
        is_e_invoice_required,
        e_invoice_status,
        created_by,
        updated_by
    ) VALUES (
        v_sale_id_3,
        v_org_id,
        'INV-2024-003',
        'REF-ETA-001',
        'sale',
        v_customer_id_2,
        v_user_id,
        5000.00,
        700.00,  -- 14% VAT in Egypt
        0.00,
        5700.00,
        5700.00,
        0.00,
        0.00,
        'paid',
        '2024-11-02 11:00:00+02',
        '2024-11-02 11:00:00+02',
        true,
        'pending',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created sale % for ETA Egyptian invoice', v_sale_id_3;

    -- Sale 4: ZATCA Credit Note
    v_sale_id_4 := gen_random_uuid();
    INSERT INTO sales (
        id,
        organization_id,
        sale_number,
        reference_number,
        transaction_type,
        customer_id,
        cashier_id,
        subtotal,
        tax_amount,
        discount_amount,
        total_amount,
        paid_amount,
        change_amount,
        outstanding_amount,
        payment_status,
        transaction_date,
        completed_at,
        is_e_invoice_required,
        e_invoice_status,
        created_by,
        updated_by
    ) VALUES (
        v_sale_id_4,
        v_org_id,
        'CN-2024-001',
        'REF-ZATCA-CN-001',
        'return',
        v_customer_id_1,
        v_user_id,
        -200.00,
        -30.00,
        0.00,
        -230.00,
        -230.00,
        0.00,
        0.00,
        'refunded',
        '2024-11-03 09:45:00+03',
        '2024-11-03 09:45:00+03',
        true,
        'pending',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created sale % for ZATCA credit note', v_sale_id_4;

    -- ============================================================================
    -- CREATE E-INVOICING DOCUMENTS
    -- ============================================================================

    RAISE NOTICE 'Creating e-invoicing documents...';

    -- Document 1: ZATCA Standard B2B Invoice (Accepted)
    v_einv_doc_id_1 := gen_random_uuid();
    INSERT INTO e_invoicing_documents (
        id,
        organization_id,
        source_table,
        source_id,
        authority,
        country_code,
        document_uuid,
        document_type,
        document_number,
        internal_reference,
        status,
        created_at,
        submitted_at,
        response_at,
        request_payload,
        response_payload,
        zatca_hash_value,
        zatca_previous_hash_value,
        zatca_invoice_counter_value,
        zatca_cryptographic_stamp,
        zatca_qr_code_payload,
        zatca_compliance_invoice_number,
        submission_format,
        created_by,
        updated_by
    ) VALUES (
        v_einv_doc_id_1,
        v_org_id,
        'sales',
        v_sale_id_1,
        'ZATCA',
        'SAU',
        gen_random_uuid(),
        'standard',
        'SME00123-2024-001',
        'INV-2024-001',
        'accepted',
        '2024-11-01 10:30:00+03',
        '2024-11-01 10:31:00+03',
        '2024-11-01 10:31:05+03',
        jsonb_build_object(
            'InvoiceNumber', 'INV-2024-001',
            'IssueDate', '2024-11-01',
            'IssueTime', '10:30:00',
            'InvoiceTypeCode', '388',
            'DocumentCurrencyCode', 'SAR',
            'TaxCurrencyCode', 'SAR',
            'LineCountNumeric', 1
        ),
        jsonb_build_object(
            'reportingStatus', 'REPORTED',
            'clearanceStatus', 'CLEARED',
            'qrCode', 'BASE64_QR_HERE',
            'validationResults', jsonb_build_object('status', 'PASS')
        ),
        'NWZlY2ViNjZmZmM4NmYzOGQ5YjIyYmI1M2QwYTE1Y2NhMTQ2ZjI2NDFmZGNlYmU0NjM5NDcyNTM5YzgxYzRhMQ==',
        'MjM0NWY3YjhkNWFiYzM0NTY3ODkwMWNkZWYyMzQ1Njc4OTBhYmNkZWYxMjM0NTY3ODkwYWJjZGVmMTIzNDU2Nzg=',
        1,
        'TUlJRFFRWUpLb1pJaHZjTkFRY0NvSUlETWpDQ0F5NENBUUl4Q3pBSkJnVXJEZ01DR0VVQU1Bc0dDU3FHU0liM0RRRUhBYUNDQWl3d2dn',
        'AQFTTUUwMDEyMwIFMzAwMDMABDE1MDAEFDIwMjQtMTEtMDFUMTA6MzA6MDAFFjExNTAuMDAGLjVmZWNlYjY2ZmZjODZmMzhkOWIyMmJiNTNkMGExNWNjYTE0NmYyNjQxZmRjZWJlNDYzOTQ3MjUzOWM4MWM0YTE=',
        'PIH-2024-001',
        'XML',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created e-invoicing document % for ZATCA standard invoice', v_einv_doc_id_1;

    -- Update sale with e-invoice reference
    UPDATE sales
    SET
        e_invoice_document_id = v_einv_doc_id_1,
        e_invoice_status = 'accepted',
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_sale_id_1;

    -- Document 2: ZATCA Simplified B2C Invoice (Accepted)
    v_einv_doc_id_2 := gen_random_uuid();
    INSERT INTO e_invoicing_documents (
        id,
        organization_id,
        source_table,
        source_id,
        authority,
        country_code,
        document_uuid,
        document_type,
        document_number,
        internal_reference,
        status,
        created_at,
        submitted_at,
        response_at,
        request_payload,
        response_payload,
        zatca_hash_value,
        zatca_previous_hash_value,
        zatca_invoice_counter_value,
        zatca_qr_code_payload,
        submission_format,
        created_by,
        updated_by
    ) VALUES (
        v_einv_doc_id_2,
        v_org_id,
        'sales',
        v_sale_id_2,
        'ZATCA',
        'SAU',
        gen_random_uuid(),
        'simplified',
        'SME00123-2024-002',
        'INV-2024-002',
        'accepted',
        '2024-11-01 14:15:00+03',
        '2024-11-01 14:15:30+03',
        '2024-11-01 14:15:33+03',
        jsonb_build_object(
            'InvoiceNumber', 'INV-2024-002',
            'IssueDate', '2024-11-01',
            'IssueTime', '14:15:00',
            'InvoiceTypeCode', '388',
            'DocumentCurrencyCode', 'SAR'
        ),
        jsonb_build_object(
            'reportingStatus', 'REPORTED',
            'clearanceStatus', 'NOT_REQUIRED',
            'qrCode', 'BASE64_QR_HERE'
        ),
        'YWJjZGVmMTIzNDU2Nzg5MGFiY2RlZjEyMzQ1Njc4OTBhYmNkZWYxMjM0NTY3ODkwYWJjZGVmMTIzNDU2Nzg5MA==',
        'NWZlY2ViNjZmZmM4NmYzOGQ5YjIyYmI1M2QwYTE1Y2NhMTQ2ZjI2NDFmZGNlYmU0NjM5NDcyNTM5YzgxYzRhMQ==',
        2,
        'AQFTTUUwMDEyMwIFMDAwMDADADQEFDIwMjQtMTEtMDFUMTQ6MTU6MDAFFDI4Ny41MAYuYWJjZGVmMTIzNDU2Nzg5MGFiY2RlZjEyMzQ1Njc4OTBhYmNkZWYxMjM0NTY3ODkwYWJjZGVmMTIzNDU2Nzg5MA==',
        'XML',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created e-invoicing document % for ZATCA simplified invoice', v_einv_doc_id_2;

    -- Update sale with e-invoice reference
    UPDATE sales
    SET
        e_invoice_document_id = v_einv_doc_id_2,
        e_invoice_status = 'accepted',
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_sale_id_2;

    -- Document 3: ETA Egyptian Invoice (Accepted)
    v_einv_doc_id_3 := gen_random_uuid();
    INSERT INTO e_invoicing_documents (
        id,
        organization_id,
        source_table,
        source_id,
        authority,
        country_code,
        document_uuid,
        document_type,
        document_number,
        internal_reference,
        status,
        created_at,
        submitted_at,
        response_at,
        request_payload,
        response_payload,
        eta_document_type_version,
        eta_submission_uuid,
        eta_long_id,
        eta_internal_id,
        eta_digital_signature,
        eta_signature_algorithm,
        submission_format,
        created_by,
        updated_by
    ) VALUES (
        v_einv_doc_id_3,
        v_org_id,
        'sales',
        v_sale_id_3,
        'ETA',
        'EGY',
        gen_random_uuid(),
        'I',  -- Invoice
        'EGS-2024-001',
        'INV-2024-003',
        'accepted',
        '2024-11-02 11:00:00+02',
        '2024-11-02 11:01:00+02',
        '2024-11-02 11:01:08+02',
        jsonb_build_object(
            'issuer', jsonb_build_object(
                'type', 'B',
                'id', '100200300',
                'name', 'My Company Ltd.'
            ),
            'receiver', jsonb_build_object(
                'type', 'B',
                'id', '100200300',
                'name', 'Customer Company Ltd.'
            ),
            'documentType', 'I',
            'documentTypeVersion', '1.0',
            'dateTimeIssued', '2024-11-02T11:00:00Z',
            'taxpayerActivityCode', '6201',
            'internalID', 'INV-2024-003',
            'totalSalesAmount', 5000.00,
            'totalAmount', 5700.00,
            'taxTotals', jsonb_build_array(
                jsonb_build_object('taxType', 'T1', 'amount', 700.00)
            )
        ),
        jsonb_build_object(
            'submissionId', gen_random_uuid(),
            'acceptedDocuments', jsonb_build_array(
                jsonb_build_object(
                    'uuid', gen_random_uuid(),
                    'longId', '7RWP5QHWTHBK5LGNHQWNVCY5PBWTXMNT42FZIZUQXQ5NDVXSCFUA',
                    'internalId', 'INV-2024-003',
                    'hashKey', 'abc123def456'
                )
            ),
            'validationSteps', jsonb_build_object('status', 'Valid')
        ),
        '1.0',
        '123e4567-e89b-12d3-a456-426614174000'::uuid,
        '7RWP5QHWTHBK5LGNHQWNVCY5PBWTXMNT42FZIZUQXQ5NDVXSCFUA',
        'INV-2024-003',
        'MIIGHgYJKoZIhvcNAQcCoIIGDzCCBgsCAQExDTALBglghkgBZQMEAgEwCwYJKoZIhvcNAQcBoIIDfDCCA3gwggJgoAMCAQICEDM',
        'CAdES-BES',
        'JSON',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created e-invoicing document % for ETA invoice', v_einv_doc_id_3;

    -- Update sale with e-invoice reference
    UPDATE sales
    SET
        e_invoice_document_id = v_einv_doc_id_3,
        e_invoice_status = 'accepted',
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_sale_id_3;

    -- Document 4: ZATCA Credit Note (Rejected - for testing)
    v_einv_doc_id_4 := gen_random_uuid();
    INSERT INTO e_invoicing_documents (
        id,
        organization_id,
        source_table,
        source_id,
        authority,
        country_code,
        document_uuid,
        document_type,
        document_number,
        internal_reference,
        status,
        created_at,
        submitted_at,
        response_at,
        request_payload,
        response_payload,
        error_code,
        error_message,
        retry_count,
        last_retry_at,
        zatca_hash_value,
        zatca_previous_hash_value,
        zatca_invoice_counter_value,
        submission_format,
        created_by,
        updated_by
    ) VALUES (
        v_einv_doc_id_4,
        v_org_id,
        'sales',
        v_sale_id_4,
        'ZATCA',
        'SAU',
        gen_random_uuid(),
        'credit_note',
        'SME00123-2024-CN-001',
        'CN-2024-001',
        'rejected',
        '2024-11-03 09:45:00+03',
        '2024-11-03 09:46:00+03',
        '2024-11-03 09:46:03+03',
        jsonb_build_object(
            'InvoiceNumber', 'CN-2024-001',
            'IssueDate', '2024-11-03',
            'IssueTime', '09:45:00',
            'InvoiceTypeCode', '381'
        ),
        jsonb_build_object(
            'reportingStatus', 'NOT_REPORTED',
            'clearanceStatus', 'REJECTED',
            'validationResults', jsonb_build_object(
                'status', 'FAIL',
                'errorMessages', jsonb_build_array(
                    jsonb_build_object(
                        'code', 'XSD_ZATCA_INVALID',
                        'category', 'ZATCA_VALIDATION',
                        'message', 'Original invoice reference is missing',
                        'severity', 'ERROR'
                    )
                )
            )
        ),
        'XSD_ZATCA_INVALID',
        'Original invoice reference is missing. Credit notes must reference the original invoice UUID.',
        1,
        '2024-11-03 10:00:00+03',
        'Y3JlZGl0bm90ZWhhc2gxMjM0NTY3ODkwYWJjZGVmMTIzNDU2Nzg5MGFiY2RlZjEyMzQ1Njc4OTBhYmNkZWYxMjM=',
        'YWJjZGVmMTIzNDU2Nzg5MGFiY2RlZjEyMzQ1Njc4OTBhYmNkZWYxMjM0NTY3ODkwYWJjZGVmMTIzNDU2Nzg5MA==',
        3,
        'XML',
        v_user_id,
        v_user_id
    );
    RAISE NOTICE 'Created e-invoicing document % for ZATCA credit note (rejected)', v_einv_doc_id_4;

    -- Update sale with e-invoice reference
    UPDATE sales
    SET
        e_invoice_document_id = v_einv_doc_id_4,
        e_invoice_status = 'rejected',
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_sale_id_4;

    -- ============================================================================
    -- CREATE E-INVOICING DOCUMENT EVENTS
    -- ============================================================================

    RAISE NOTICE 'Creating e-invoicing document events...';

    -- Events for Document 1: ZATCA Standard B2B Invoice
    INSERT INTO e_invoicing_document_events (
        organization_id,
        e_invoicing_document_id,
        event_type,
        event_timestamp,
        previous_status,
        new_status,
        event_description,
        triggered_by,
        user_id
    ) VALUES
    (
        v_org_id,
        v_einv_doc_id_1,
        'created',
        '2024-11-01 10:30:00+03',
        NULL,
        'draft',
        'E-invoice document created from sale INV-2024-001',
        'system',
        v_user_id
    ),
    (
        v_org_id,
        v_einv_doc_id_1,
        'validated',
        '2024-11-01 10:30:30+03',
        'draft',
        'pending',
        'Local validation passed, ready for submission',
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_1,
        'submitted',
        '2024-11-01 10:31:00+03',
        'pending',
        'submitted',
        'Document submitted to ZATCA for clearance',
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_1,
        'accepted',
        '2024-11-01 10:31:05+03',
        'submitted',
        'accepted',
        'Document accepted and cleared by ZATCA',
        'system',
        NULL
    );

    -- Events for Document 2: ZATCA Simplified B2C Invoice
    INSERT INTO e_invoicing_document_events (
        organization_id,
        e_invoicing_document_id,
        event_type,
        event_timestamp,
        previous_status,
        new_status,
        event_description,
        triggered_by,
        user_id
    ) VALUES
    (
        v_org_id,
        v_einv_doc_id_2,
        'created',
        '2024-11-01 14:15:00+03',
        NULL,
        'draft',
        'E-invoice document created from sale INV-2024-002',
        'system',
        v_user_id
    ),
    (
        v_org_id,
        v_einv_doc_id_2,
        'validated',
        '2024-11-01 14:15:15+03',
        'draft',
        'pending',
        'Local validation passed, ready for reporting',
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_2,
        'submitted',
        '2024-11-01 14:15:30+03',
        'pending',
        'submitted',
        'Document submitted to ZATCA for reporting (simplified invoice - no clearance required)',
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_2,
        'accepted',
        '2024-11-01 14:15:33+03',
        'submitted',
        'accepted',
        'Document accepted by ZATCA',
        'system',
        NULL
    );

    -- Events for Document 3: ETA Egyptian Invoice
    INSERT INTO e_invoicing_document_events (
        organization_id,
        e_invoicing_document_id,
        event_type,
        event_timestamp,
        previous_status,
        new_status,
        event_description,
        triggered_by,
        user_id
    ) VALUES
    (
        v_org_id,
        v_einv_doc_id_3,
        'created',
        '2024-11-02 11:00:00+02',
        NULL,
        'draft',
        'E-invoice document created from sale INV-2024-003',
        'system',
        v_user_id
    ),
    (
        v_org_id,
        v_einv_doc_id_3,
        'validated',
        '2024-11-02 11:00:45+02',
        'draft',
        'pending',
        'Local validation passed, ready for ETA submission',
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_3,
        'submitted',
        '2024-11-02 11:01:00+02',
        'pending',
        'submitted',
        'Document submitted to ETA e-invoicing portal',
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_3,
        'accepted',
        '2024-11-02 11:01:08+02',
        'submitted',
        'accepted',
        'Document accepted by ETA. Long ID assigned: 7RWP5QHWTHBK5LGNHQWNVCY5PBWTXMNT42FZIZUQXQ5NDVXSCFUA',
        'system',
        NULL
    );

    -- Events for Document 4: ZATCA Credit Note (Rejected)
    INSERT INTO e_invoicing_document_events (
        organization_id,
        e_invoicing_document_id,
        event_type,
        event_timestamp,
        previous_status,
        new_status,
        event_description,
        event_data,
        http_status_code,
        error_code,
        error_message,
        triggered_by,
        user_id
    ) VALUES
    (
        v_org_id,
        v_einv_doc_id_4,
        'created',
        '2024-11-03 09:45:00+03',
        NULL,
        'draft',
        'E-invoice document created from sale CN-2024-001',
        NULL,
        NULL,
        NULL,
        NULL,
        'system',
        v_user_id
    ),
    (
        v_org_id,
        v_einv_doc_id_4,
        'validated',
        '2024-11-03 09:45:45+03',
        'draft',
        'pending',
        'Local validation passed, ready for submission',
        NULL,
        NULL,
        NULL,
        NULL,
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_4,
        'submitted',
        '2024-11-03 09:46:00+03',
        'pending',
        'submitted',
        'Document submitted to ZATCA for clearance',
        NULL,
        NULL,
        NULL,
        NULL,
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_4,
        'rejected',
        '2024-11-03 09:46:03+03',
        'submitted',
        'rejected',
        'Document rejected by ZATCA due to validation errors',
        jsonb_build_object(
            'validationErrors', jsonb_build_array(
                jsonb_build_object(
                    'code', 'XSD_ZATCA_INVALID',
                    'category', 'ZATCA_VALIDATION',
                    'message', 'Original invoice reference is missing',
                    'severity', 'ERROR',
                    'xpath', '/Invoice/BillingReference/InvoiceDocumentReference/ID'
                )
            )
        ),
        400,
        'XSD_ZATCA_INVALID',
        'Original invoice reference is missing. Credit notes must reference the original invoice UUID.',
        'system',
        NULL
    ),
    (
        v_org_id,
        v_einv_doc_id_4,
        'retry',
        '2024-11-03 10:00:00+03',
        'rejected',
        'rejected',
        'Automatic retry attempt #1 failed - same validation error',
        jsonb_build_object('retryCount', 1),
        400,
        'XSD_ZATCA_INVALID',
        'Original invoice reference is missing. Credit notes must reference the original invoice UUID.',
        'scheduled',
        NULL
    );

    RAISE NOTICE 'Created % e-invoicing document events',
        (SELECT COUNT(*) FROM e_invoicing_document_events WHERE organization_id = v_org_id);

    -- ============================================================================
    -- SUMMARY
    -- ============================================================================

    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'E-Invoicing Seed Data Summary:';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Updated Customers: 3 (1 ZATCA B2B, 1 ETA, 1 ZATCA B2C)';
    RAISE NOTICE 'Updated Products: 3 (1 standard, 1 zero-rated, 1 excise)';
    RAISE NOTICE 'Created Sales: 4 (2 ZATCA, 1 ETA, 1 credit note)';
    RAISE NOTICE 'Created E-Invoicing Documents: 4';
    RAISE NOTICE '  - ZATCA Standard (Accepted): 1';
    RAISE NOTICE '  - ZATCA Simplified (Accepted): 1';
    RAISE NOTICE '  - ETA Invoice (Accepted): 1';
    RAISE NOTICE '  - ZATCA Credit Note (Rejected): 1';
    RAISE NOTICE 'Created Document Events: %',
        (SELECT COUNT(*) FROM e_invoicing_document_events WHERE organization_id = v_org_id);
    RAISE NOTICE '============================================';

END $$;

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Seed Data 007 completed successfully!';
    RAISE NOTICE 'E-Invoicing data created for ZATCA and ETA.';
    RAISE NOTICE '============================================';
END $$;

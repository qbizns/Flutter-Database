package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/einvoicing"
)

// EInvoicingRepository handles e-invoicing data access
type EInvoicingRepository struct {
	db *DB
}

// NewEInvoicingRepository creates a new e-invoicing repository
func NewEInvoicingRepository(db *DB) *EInvoicingRepository {
	return &EInvoicingRepository{db: db}
}

// DOCUMENT CRUD OPERATIONS

// List retrieves e-invoicing documents
func (r *EInvoicingRepository) List(ctx context.Context, orgID uuid.UUID, filters einvoicing.EInvoicingDocumentFilters) ([]einvoicing.EInvoicingDocument, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, authority, country_code,
		       document_uuid, document_type, document_number, internal_reference, status,
		       created_at, submitted_at, response_at, request_payload, response_payload,
		       error_code, error_message, retry_count, last_retry_at,
		       zatca_hash_value, zatca_previous_hash_value, zatca_invoice_counter_value,
		       zatca_cryptographic_stamp, zatca_qr_code_payload, zatca_compliance_invoice_number,
		       eta_document_type_version, eta_submission_uuid, eta_long_id, eta_internal_id,
		       eta_digital_signature, eta_signature_algorithm, submission_format, metadata,
		       updated_at, deleted_at, created_by, updated_by
		FROM e_invoicing_documents
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (document_number ILIKE $%d OR internal_reference ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Authority != nil {
		argCount++
		query += fmt.Sprintf(" AND authority = $%d", argCount)
		args = append(args, *filters.Authority)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.DocumentType != nil {
		argCount++
		query += fmt.Sprintf(" AND document_type = $%d", argCount)
		args = append(args, *filters.DocumentType)
	}

	if filters.SourceTable != nil {
		argCount++
		query += fmt.Sprintf(" AND source_table = $%d", argCount)
		args = append(args, *filters.SourceTable)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []einvoicing.EInvoicingDocument
	for rows.Next() {
		doc, err := scanEInvoicingDocument(rows)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, rows.Err()
}

// Count counts e-invoicing documents
func (r *EInvoicingRepository) Count(ctx context.Context, orgID uuid.UUID, filters einvoicing.EInvoicingDocumentFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM e_invoicing_documents WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Authority != nil {
		argCount++
		query += fmt.Sprintf(" AND authority = $%d", argCount)
		args = append(args, *filters.Authority)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// Create creates an e-invoicing document
func (r *EInvoicingRepository) Create(ctx context.Context, doc *einvoicing.EInvoicingDocument) error {
	if err := r.db.SetOrganizationContext(ctx, doc.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO e_invoicing_documents (
			id, organization_id, source_table, source_id, authority, country_code,
			document_uuid, document_type, document_number, internal_reference, status,
			created_at, submitted_at, response_at, request_payload, response_payload,
			error_code, error_message, retry_count, last_retry_at,
			zatca_hash_value, zatca_previous_hash_value, zatca_invoice_counter_value,
			zatca_cryptographic_stamp, zatca_qr_code_payload, zatca_compliance_invoice_number,
			eta_document_type_version, eta_submission_uuid, eta_long_id, eta_internal_id,
			eta_digital_signature, eta_signature_algorithm, submission_format, metadata,
			updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37
		)
	`

	reqPayload, _ := json.Marshal(doc.RequestPayload)
	respPayload, _ := json.Marshal(doc.ResponsePayload)
	metadata, _ := json.Marshal(doc.Metadata)

	_, err := r.db.Pool.Exec(ctx, query,
		doc.ID, doc.OrganizationID, doc.SourceTable, doc.SourceID, doc.Authority, doc.CountryCode,
		doc.DocumentUUID, doc.DocumentType, doc.DocumentNumber, doc.InternalReference, doc.Status,
		doc.CreatedAt, doc.SubmittedAt, doc.ResponseAt, reqPayload, respPayload,
		doc.ErrorCode, doc.ErrorMessage, doc.RetryCount, doc.LastRetryAt,
		doc.ZATCAHashValue, doc.ZATCAPreviousHashValue, doc.ZATCAInvoiceCounterValue,
		doc.ZATCACryptographicStamp, doc.ZATCAQRCodePayload, doc.ZATCAComplianceInvoiceNumber,
		doc.ETADocumentTypeVersion, doc.ETASubmissionUUID, doc.ETALongID, doc.ETAInternalID,
		doc.ETADigitalSignature, doc.ETASignatureAlgorithm, doc.SubmissionFormat, metadata,
		doc.UpdatedAt, doc.CreatedBy, doc.UpdatedBy,
	)

	return err
}

// Get retrieves an e-invoicing document by ID
func (r *EInvoicingRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*einvoicing.EInvoicingDocument, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, authority, country_code,
		       document_uuid, document_type, document_number, internal_reference, status,
		       created_at, submitted_at, response_at, request_payload, response_payload,
		       error_code, error_message, retry_count, last_retry_at,
		       zatca_hash_value, zatca_previous_hash_value, zatca_invoice_counter_value,
		       zatca_cryptographic_stamp, zatca_qr_code_payload, zatca_compliance_invoice_number,
		       eta_document_type_version, eta_submission_uuid, eta_long_id, eta_internal_id,
		       eta_digital_signature, eta_signature_algorithm, submission_format, metadata,
		       updated_at, deleted_at, created_by, updated_by
		FROM e_invoicing_documents
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	row := r.db.Pool.QueryRow(ctx, query, orgID, id)
	return scanEInvoicingDocumentRow(row)
}

// GetByDocumentNumber retrieves document by document number
func (r *EInvoicingRepository) GetByDocumentNumber(ctx context.Context, orgID uuid.UUID, authority string, documentNumber string) (*einvoicing.EInvoicingDocument, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, authority, country_code,
		       document_uuid, document_type, document_number, internal_reference, status,
		       created_at, submitted_at, response_at, request_payload, response_payload,
		       error_code, error_message, retry_count, last_retry_at,
		       zatca_hash_value, zatca_previous_hash_value, zatca_invoice_counter_value,
		       zatca_cryptographic_stamp, zatca_qr_code_payload, zatca_compliance_invoice_number,
		       eta_document_type_version, eta_submission_uuid, eta_long_id, eta_internal_id,
		       eta_digital_signature, eta_signature_algorithm, submission_format, metadata,
		       updated_at, deleted_at, created_by, updated_by
		FROM e_invoicing_documents
		WHERE organization_id = $1 AND authority = $2 AND document_number = $3 AND deleted_at IS NULL
	`

	row := r.db.Pool.QueryRow(ctx, query, orgID, authority, documentNumber)
	return scanEInvoicingDocumentRow(row)
}

// GetByDocumentUUID retrieves document by document UUID
func (r *EInvoicingRepository) GetByDocumentUUID(ctx context.Context, orgID uuid.UUID, documentUUID uuid.UUID) (*einvoicing.EInvoicingDocument, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, authority, country_code,
		       document_uuid, document_type, document_number, internal_reference, status,
		       created_at, submitted_at, response_at, request_payload, response_payload,
		       error_code, error_message, retry_count, last_retry_at,
		       zatca_hash_value, zatca_previous_hash_value, zatca_invoice_counter_value,
		       zatca_cryptographic_stamp, zatca_qr_code_payload, zatca_compliance_invoice_number,
		       eta_document_type_version, eta_submission_uuid, eta_long_id, eta_internal_id,
		       eta_digital_signature, eta_signature_algorithm, submission_format, metadata,
		       updated_at, deleted_at, created_by, updated_by
		FROM e_invoicing_documents
		WHERE organization_id = $1 AND document_uuid = $2 AND deleted_at IS NULL
	`

	row := r.db.Pool.QueryRow(ctx, query, orgID, documentUUID)
	return scanEInvoicingDocumentRow(row)
}

// Update updates an e-invoicing document
func (r *EInvoicingRepository) Update(ctx context.Context, doc *einvoicing.EInvoicingDocument) error {
	if err := r.db.SetOrganizationContext(ctx, doc.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE e_invoicing_documents SET
			status = $2, submitted_at = $3, response_at = $4,
			error_code = $5, error_message = $6, retry_count = $7, last_retry_at = $8,
			zatca_hash_value = $9, zatca_previous_hash_value = $10, zatca_invoice_counter_value = $11,
			zatca_cryptographic_stamp = $12, zatca_qr_code_payload = $13,
			eta_submission_uuid = $14, eta_long_id = $15, eta_digital_signature = $16,
			metadata = $17, updated_at = $18, updated_by = $19
		WHERE id = $1 AND organization_id = $20
	`

	metadata, _ := json.Marshal(doc.Metadata)

	_, err := r.db.Pool.Exec(ctx, query,
		doc.ID, doc.Status, doc.SubmittedAt, doc.ResponseAt,
		doc.ErrorCode, doc.ErrorMessage, doc.RetryCount, doc.LastRetryAt,
		doc.ZATCAHashValue, doc.ZATCAPreviousHashValue, doc.ZATCAInvoiceCounterValue,
		doc.ZATCACryptographicStamp, doc.ZATCAQRCodePayload,
		doc.ETASubmissionUUID, doc.ETALongID, doc.ETADigitalSignature,
		metadata, doc.UpdatedAt, doc.UpdatedBy, doc.OrganizationID,
	)

	return err
}

// Delete deletes an e-invoicing document (soft delete)
func (r *EInvoicingRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `UPDATE e_invoicing_documents SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND organization_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, id, orgID)
	return err
}

// GetBySourceReference retrieves documents by source reference
func (r *EInvoicingRepository) GetBySourceReference(ctx context.Context, orgID uuid.UUID, sourceTable string, sourceID uuid.UUID) ([]einvoicing.EInvoicingDocument, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, authority, country_code,
		       document_uuid, document_type, document_number, internal_reference, status,
		       created_at, submitted_at, response_at, request_payload, response_payload,
		       error_code, error_message, retry_count, last_retry_at,
		       zatca_hash_value, zatca_previous_hash_value, zatca_invoice_counter_value,
		       zatca_cryptographic_stamp, zatca_qr_code_payload, zatca_compliance_invoice_number,
		       eta_document_type_version, eta_submission_uuid, eta_long_id, eta_internal_id,
		       eta_digital_signature, eta_signature_algorithm, submission_format, metadata,
		       updated_at, deleted_at, created_by, updated_by
		FROM e_invoicing_documents
		WHERE organization_id = $1 AND source_table = $2 AND source_id = $3 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, sourceTable, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []einvoicing.EInvoicingDocument
	for rows.Next() {
		doc, err := scanEInvoicingDocument(rows)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, rows.Err()
}

// UpdateStatus updates only the status field
func (r *EInvoicingRepository) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `UPDATE e_invoicing_documents SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND organization_id = $3`
	_, err := r.db.Pool.Exec(ctx, query, id, status, orgID)
	return err
}

// EVENT OPERATIONS

// CreateEvent creates an event record
func (r *EInvoicingRepository) CreateEvent(ctx context.Context, event *einvoicing.EInvoicingDocumentEvent) error {
	if err := r.db.SetOrganizationContext(ctx, event.OrganizationID.String()); err != nil {
		return err
	}

	eventData, _ := json.Marshal(event.EventData)
	reqHeaders, _ := json.Marshal(event.RequestHeaders)
	respHeaders, _ := json.Marshal(event.ResponseHeaders)
	errorDetails, _ := json.Marshal(event.ErrorDetails)

	query := `
		INSERT INTO e_invoicing_document_events (
			id, organization_id, e_invoicing_document_id, event_type, event_timestamp,
			previous_status, new_status, event_description, event_data,
			http_status_code, http_method, api_endpoint, request_headers, response_headers,
			error_code, error_message, error_details, triggered_by, user_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		event.ID, event.OrganizationID, event.EInvoicingDocumentID, event.EventType, event.EventTimestamp,
		event.PreviousStatus, event.NewStatus, event.EventDescription, eventData,
		event.HTTPStatusCode, event.HTTPMethod, event.APIEndpoint, reqHeaders, respHeaders,
		event.ErrorCode, event.ErrorMessage, errorDetails, event.TriggeredBy, event.UserID, event.CreatedAt,
	)

	return err
}

// ListEvents retrieves events
func (r *EInvoicingRepository) ListEvents(ctx context.Context, orgID uuid.UUID, filters einvoicing.EInvoicingDocumentEventFilters) ([]einvoicing.EInvoicingDocumentEvent, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, e_invoicing_document_id, event_type, event_timestamp,
		       previous_status, new_status, event_description, event_data,
		       http_status_code, http_method, api_endpoint, request_headers, response_headers,
		       error_code, error_message, error_details, triggered_by, user_id, created_at
		FROM e_invoicing_document_events
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.EInvoicingDocID != nil {
		argCount++
		query += fmt.Sprintf(" AND e_invoicing_document_id = $%d", argCount)
		args = append(args, *filters.EInvoicingDocID)
	}

	if filters.EventType != nil {
		argCount++
		query += fmt.Sprintf(" AND event_type = $%d", argCount)
		args = append(args, *filters.EventType)
	}

	query += " ORDER BY event_timestamp DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []einvoicing.EInvoicingDocumentEvent
	for rows.Next() {
		event, err := scanEInvoicingDocumentEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

// CountEvents counts events
func (r *EInvoicingRepository) CountEvents(ctx context.Context, orgID uuid.UUID, filters einvoicing.EInvoicingDocumentEventFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM e_invoicing_document_events WHERE organization_id = $1"
	args := []interface{}{orgID}
	argCount := 1

	if filters.EInvoicingDocID != nil {
		argCount++
		query += fmt.Sprintf(" AND e_invoicing_document_id = $%d", argCount)
		args = append(args, *filters.EInvoicingDocID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetEvent retrieves a specific event
func (r *EInvoicingRepository) GetEvent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*einvoicing.EInvoicingDocumentEvent, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, e_invoicing_document_id, event_type, event_timestamp,
		       previous_status, new_status, event_description, event_data,
		       http_status_code, http_method, api_endpoint, request_headers, response_headers,
		       error_code, error_message, error_details, triggered_by, user_id, created_at
		FROM e_invoicing_document_events
		WHERE id = $1 AND organization_id = $2
	`

	row := r.db.Pool.QueryRow(ctx, query, id, orgID)

	var event einvoicing.EInvoicingDocumentEvent
	var eventData, reqHeaders, respHeaders, errorDetails []byte

	err := row.Scan(
		&event.ID, &event.OrganizationID, &event.EInvoicingDocumentID, &event.EventType, &event.EventTimestamp,
		&event.PreviousStatus, &event.NewStatus, &event.EventDescription, &eventData,
		&event.HTTPStatusCode, &event.HTTPMethod, &event.APIEndpoint, &reqHeaders, &respHeaders,
		&event.ErrorCode, &event.ErrorMessage, &errorDetails, &event.TriggeredBy, &event.UserID, &event.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(eventData, &event.EventData)
	_ = json.Unmarshal(reqHeaders, &event.RequestHeaders)
	_ = json.Unmarshal(respHeaders, &event.ResponseHeaders)
	_ = json.Unmarshal(errorDetails, &event.ErrorDetails)

	return &event, nil
}

// DOCUMENT SEQUENCE OPERATIONS

// CreateSequence creates a document sequence
func (r *EInvoicingRepository) CreateSequence(ctx context.Context, seq *einvoicing.DocumentSequence) error {
	if err := r.db.SetOrganizationContext(ctx, seq.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO document_sequences (
			id, organization_id, document_type, prefix, suffix, next_number, padding,
			increment_by, reset_frequency, last_reset_at, last_reset_value, include_date,
			date_format, location_id, is_active, allow_manual_override, example_number,
			description, notes, metadata, created_at, updated_at, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`

	metadata, _ := json.Marshal(seq.Metadata)

	_, err := r.db.Pool.Exec(ctx, query,
		seq.ID, seq.OrganizationID, seq.DocumentType, seq.Prefix, seq.Suffix, seq.NextNumber, seq.Padding,
		seq.IncrementBy, seq.ResetFrequency, seq.LastResetAt, seq.LastResetValue, seq.IncludeDate,
		seq.DateFormat, seq.LocationID, seq.IsActive, seq.AllowManualOverride, seq.ExampleNumber,
		seq.Description, seq.Notes, metadata, seq.CreatedAt, seq.UpdatedAt, seq.CreatedBy, seq.UpdatedBy,
	)

	return err
}

// GetSequence retrieves a sequence by ID
func (r *EInvoicingRepository) GetSequence(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*einvoicing.DocumentSequence, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, document_type, prefix, suffix, next_number, padding,
		       increment_by, reset_frequency, last_reset_at, last_reset_value, include_date,
		       date_format, location_id, is_active, allow_manual_override, example_number,
		       description, notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM document_sequences
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	return scanDocumentSequenceRow(r.db.Pool.QueryRow(ctx, query, id, orgID))
}

// GetByDocumentType retrieves a sequence by document type
func (r *EInvoicingRepository) GetByDocumentType(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (*einvoicing.DocumentSequence, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, document_type, prefix, suffix, next_number, padding,
		       increment_by, reset_frequency, last_reset_at, last_reset_value, include_date,
		       date_format, location_id, is_active, allow_manual_override, example_number,
		       description, notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM document_sequences
		WHERE organization_id = $1 AND document_type = $2 AND deleted_at IS NULL
	`

	if locationID != nil {
		query += " AND location_id = $3"
		return scanDocumentSequenceRow(r.db.Pool.QueryRow(ctx, query, orgID, documentType, locationID))
	}

	query += " AND location_id IS NULL"
	return scanDocumentSequenceRow(r.db.Pool.QueryRow(ctx, query, orgID, documentType))
}

// ListSequences retrieves sequences with filters
func (r *EInvoicingRepository) ListSequences(ctx context.Context, orgID uuid.UUID, filters einvoicing.DocumentSequenceFilters) ([]einvoicing.DocumentSequence, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, document_type, prefix, suffix, next_number, padding,
		       increment_by, reset_frequency, last_reset_at, last_reset_value, include_date,
		       date_format, location_id, is_active, allow_manual_override, example_number,
		       description, notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM document_sequences
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.DocumentType != nil {
		argCount++
		query += fmt.Sprintf(" AND document_type = $%d", argCount)
		args = append(args, *filters.DocumentType)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY document_type ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seqs []einvoicing.DocumentSequence
	for rows.Next() {
		seq, err := scanDocumentSequence(rows)
		if err != nil {
			return nil, err
		}
		seqs = append(seqs, seq)
	}

	return seqs, rows.Err()
}

// UpdateSequence updates a sequence
func (r *EInvoicingRepository) UpdateSequence(ctx context.Context, seq *einvoicing.DocumentSequence) error {
	if err := r.db.SetOrganizationContext(ctx, seq.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE document_sequences SET
			next_number = $2, reset_frequency = $3, last_reset_at = $4, last_reset_value = $5,
			is_active = $6, allow_manual_override = $7, example_number = $8, description = $9,
			notes = $10, metadata = $11, updated_at = $12, updated_by = $13
		WHERE id = $1 AND organization_id = $14
	`

	metadata, _ := json.Marshal(seq.Metadata)

	_, err := r.db.Pool.Exec(ctx, query,
		seq.ID, seq.NextNumber, seq.ResetFrequency, seq.LastResetAt, seq.LastResetValue,
		seq.IsActive, seq.AllowManualOverride, seq.ExampleNumber, seq.Description,
		seq.Notes, metadata, seq.UpdatedAt, seq.UpdatedBy, seq.OrganizationID,
	)

	return err
}

// DeleteSequence deletes a sequence
func (r *EInvoicingRepository) DeleteSequence(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `UPDATE document_sequences SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND organization_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, id, orgID)
	return err
}

// GetNextNumber gets next document number via PostgreSQL function
func (r *EInvoicingRepository) GetNextNumber(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (string, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return "", err
	}

	query := "SELECT get_next_document_number($1, $2, $3)"

	var result string
	err := r.db.Pool.QueryRow(ctx, query, orgID, documentType, locationID).Scan(&result)
	return result, err
}

// PreviewNumber previews next document number via PostgreSQL function
func (r *EInvoicingRepository) PreviewNumber(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (string, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return "", err
	}

	query := "SELECT preview_document_number($1, $2, $3)"

	var result string
	err := r.db.Pool.QueryRow(ctx, query, orgID, documentType, locationID).Scan(&result)
	return result, err
}

// ResetSequence resets a sequence via PostgreSQL function
func (r *EInvoicingRepository) ResetSequence(ctx context.Context, orgID uuid.UUID, documentType string, resetTo int64, locationID *uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "SELECT reset_document_sequence($1, $2, $3, $4)"
	_, err := r.db.Pool.Exec(ctx, query, orgID, documentType, resetTo, locationID)
	return err
}

// Helper functions

func scanEInvoicingDocument(rows interface{ Scan(...interface{}) error }) (einvoicing.EInvoicingDocument, error) {
	doc := einvoicing.EInvoicingDocument{}
	var reqPayload, respPayload, metadata []byte

	err := rows.Scan(
		&doc.ID, &doc.OrganizationID, &doc.SourceTable, &doc.SourceID, &doc.Authority, &doc.CountryCode,
		&doc.DocumentUUID, &doc.DocumentType, &doc.DocumentNumber, &doc.InternalReference, &doc.Status,
		&doc.CreatedAt, &doc.SubmittedAt, &doc.ResponseAt, &reqPayload, &respPayload,
		&doc.ErrorCode, &doc.ErrorMessage, &doc.RetryCount, &doc.LastRetryAt,
		&doc.ZATCAHashValue, &doc.ZATCAPreviousHashValue, &doc.ZATCAInvoiceCounterValue,
		&doc.ZATCACryptographicStamp, &doc.ZATCAQRCodePayload, &doc.ZATCAComplianceInvoiceNumber,
		&doc.ETADocumentTypeVersion, &doc.ETASubmissionUUID, &doc.ETALongID, &doc.ETAInternalID,
		&doc.ETADigitalSignature, &doc.ETASignatureAlgorithm, &doc.SubmissionFormat, &metadata,
		&doc.UpdatedAt, &doc.DeletedAt, &doc.CreatedBy, &doc.UpdatedBy,
	)

	if err == nil {
		_ = json.Unmarshal(reqPayload, &doc.RequestPayload)
		_ = json.Unmarshal(respPayload, &doc.ResponsePayload)
		_ = json.Unmarshal(metadata, &doc.Metadata)
	}

	return doc, err
}

func scanEInvoicingDocumentRow(row interface{ Scan(...interface{}) error }) (*einvoicing.EInvoicingDocument, error) {
	doc := &einvoicing.EInvoicingDocument{}
	var reqPayload, respPayload, metadata []byte

	err := row.Scan(
		&doc.ID, &doc.OrganizationID, &doc.SourceTable, &doc.SourceID, &doc.Authority, &doc.CountryCode,
		&doc.DocumentUUID, &doc.DocumentType, &doc.DocumentNumber, &doc.InternalReference, &doc.Status,
		&doc.CreatedAt, &doc.SubmittedAt, &doc.ResponseAt, &reqPayload, &respPayload,
		&doc.ErrorCode, &doc.ErrorMessage, &doc.RetryCount, &doc.LastRetryAt,
		&doc.ZATCAHashValue, &doc.ZATCAPreviousHashValue, &doc.ZATCAInvoiceCounterValue,
		&doc.ZATCACryptographicStamp, &doc.ZATCAQRCodePayload, &doc.ZATCAComplianceInvoiceNumber,
		&doc.ETADocumentTypeVersion, &doc.ETASubmissionUUID, &doc.ETALongID, &doc.ETAInternalID,
		&doc.ETADigitalSignature, &doc.ETASignatureAlgorithm, &doc.SubmissionFormat, &metadata,
		&doc.UpdatedAt, &doc.DeletedAt, &doc.CreatedBy, &doc.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(reqPayload, &doc.RequestPayload)
	_ = json.Unmarshal(respPayload, &doc.ResponsePayload)
	_ = json.Unmarshal(metadata, &doc.Metadata)

	return doc, nil
}

func scanEInvoicingDocumentEvent(rows interface{ Scan(...interface{}) error }) (einvoicing.EInvoicingDocumentEvent, error) {
	event := einvoicing.EInvoicingDocumentEvent{}
	var eventData, reqHeaders, respHeaders, errorDetails []byte

	err := rows.Scan(
		&event.ID, &event.OrganizationID, &event.EInvoicingDocumentID, &event.EventType, &event.EventTimestamp,
		&event.PreviousStatus, &event.NewStatus, &event.EventDescription, &eventData,
		&event.HTTPStatusCode, &event.HTTPMethod, &event.APIEndpoint, &reqHeaders, &respHeaders,
		&event.ErrorCode, &event.ErrorMessage, &errorDetails, &event.TriggeredBy, &event.UserID, &event.CreatedAt,
	)

	if err == nil {
		_ = json.Unmarshal(eventData, &event.EventData)
		_ = json.Unmarshal(reqHeaders, &event.RequestHeaders)
		_ = json.Unmarshal(respHeaders, &event.ResponseHeaders)
		_ = json.Unmarshal(errorDetails, &event.ErrorDetails)
	}

	return event, err
}

func scanDocumentSequence(rows interface{ Scan(...interface{}) error }) (einvoicing.DocumentSequence, error) {
	seq := einvoicing.DocumentSequence{}
	var metadata []byte

	err := rows.Scan(
		&seq.ID, &seq.OrganizationID, &seq.DocumentType, &seq.Prefix, &seq.Suffix, &seq.NextNumber, &seq.Padding,
		&seq.IncrementBy, &seq.ResetFrequency, &seq.LastResetAt, &seq.LastResetValue, &seq.IncludeDate,
		&seq.DateFormat, &seq.LocationID, &seq.IsActive, &seq.AllowManualOverride, &seq.ExampleNumber,
		&seq.Description, &seq.Notes, &metadata, &seq.CreatedAt, &seq.UpdatedAt, &seq.DeletedAt, &seq.CreatedBy, &seq.UpdatedBy,
	)

	if err == nil {
		_ = json.Unmarshal(metadata, &seq.Metadata)
	}

	return seq, err
}

func scanDocumentSequenceRow(row interface{ Scan(...interface{}) error }) (*einvoicing.DocumentSequence, error) {
	seq := &einvoicing.DocumentSequence{}
	var metadata []byte

	err := row.Scan(
		&seq.ID, &seq.OrganizationID, &seq.DocumentType, &seq.Prefix, &seq.Suffix, &seq.NextNumber, &seq.Padding,
		&seq.IncrementBy, &seq.ResetFrequency, &seq.LastResetAt, &seq.LastResetValue, &seq.IncludeDate,
		&seq.DateFormat, &seq.LocationID, &seq.IsActive, &seq.AllowManualOverride, &seq.ExampleNumber,
		&seq.Description, &seq.Notes, &metadata, &seq.CreatedAt, &seq.UpdatedAt, &seq.DeletedAt, &seq.CreatedBy, &seq.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(metadata, &seq.Metadata)
	return seq, nil
}


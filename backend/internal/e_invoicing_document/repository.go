package e_invoicing_document

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for EInvoicingDocuments
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new EInvoicingDocuments repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// EInvoicingDocuments represents a e_invoicing_documents entity
type EInvoicingDocuments struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SourceTable string `json:"source_table" db:"source_table"`
	SourceId uuid.UUID `json:"source_id" db:"source_id"`
	Authority string `json:"authority" db:"authority"`
	CountryCode string `json:"country_code" db:"country_code"`
	DocumentUuid uuid.UUID `json:"document_uuid" db:"document_uuid"`
	DocumentType string `json:"document_type" db:"document_type"`
	DocumentNumber string `json:"document_number" db:"document_number"`
	InternalReference *string `json:"internal_reference" db:"internal_reference"`
	Status string `json:"status" db:"status"`
	'draft', *string `json:"'draft'," db:"'draft',"`
	'pending', *string `json:"'pending'," db:"'pending',"`
	'submitted', *string `json:"'submitted'," db:"'submitted',"`
	'accepted', *string `json:"'accepted'," db:"'accepted',"`
	'rejected', *string `json:"'rejected'," db:"'rejected',"`
	'cancelled', *string `json:"'cancelled'," db:"'cancelled',"`
	'error' *string `json:"'error'" db:"'error'"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	SubmittedAt *time.Time `json:"submitted_at" db:"submitted_at"`
	ResponseAt *time.Time `json:"response_at" db:"response_at"`
	RequestPayload json.RawMessage `json:"request_payload" db:"request_payload"`
	ResponsePayload json.RawMessage `json:"response_payload" db:"response_payload"`
	ErrorCode *string `json:"error_code" db:"error_code"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	RetryCount *int64 `json:"retry_count" db:"retry_count"`
	LastRetryAt *time.Time `json:"last_retry_at" db:"last_retry_at"`
	ZatcaHashValue *string `json:"zatca_hash_value" db:"zatca_hash_value"`
	ZatcaPreviousHashValue *string `json:"zatca_previous_hash_value" db:"zatca_previous_hash_value"`
	ZatcaInvoiceCounterValue *int64 `json:"zatca_invoice_counter_value" db:"zatca_invoice_counter_value"`
	ZatcaCryptographicStamp *string `json:"zatca_cryptographic_stamp" db:"zatca_cryptographic_stamp"`
	ZatcaQrCodePayload *string `json:"zatca_qr_code_payload" db:"zatca_qr_code_payload"`
	ZatcaComplianceInvoiceNumber *string `json:"zatca_compliance_invoice_number" db:"zatca_compliance_invoice_number"`
	EtaDocumentTypeVersion *string `json:"eta_document_type_version" db:"eta_document_type_version"`
	EtaSubmissionUuid *uuid.UUID `json:"eta_submission_uuid" db:"eta_submission_uuid"`
	EtaLongId *string `json:"eta_long_id" db:"eta_long_id"`
	EtaInternalId *string `json:"eta_internal_id" db:"eta_internal_id"`
	EtaDigitalSignature *string `json:"eta_digital_signature" db:"eta_digital_signature"`
	EtaSignatureAlgorithm *string `json:"eta_signature_algorithm" db:"eta_signature_algorithm"`
	SubmissionFormat *string `json:"submission_format" db:"submission_format"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new e_invoicing_documents record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocuments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "e_invoicing_documents", duration, nil)
	}()

	query := `
		INSERT INTO e_invoicing_documents (
			, organization_id
			, source_table
			, source_id
			, authority
			, country_code
			, document_uuid
			, document_type
			, document_number
			, internal_reference
			, status
			, 'draft',
			, 'pending',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error'
			, submitted_at
			, response_at
			, request_payload
			, response_payload
			, error_code
			, error_message
			, retry_count
			, last_retry_at
			, zatca_hash_value
			, zatca_previous_hash_value
			, zatca_invoice_counter_value
			, zatca_cryptographic_stamp
			, zatca_qr_code_payload
			, zatca_compliance_invoice_number
			, eta_document_type_version
			, eta_submission_uuid
			, eta_long_id
			, eta_internal_id
			, eta_digital_signature
			, eta_signature_algorithm
			, submission_format
			, metadata
			, deleted_at
			, created_by
			, updated_by
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $10
			, $11
			, $12
			, $13
			, $14
			, $15
			, $16
			, $17
			, $18
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $27
			, $28
			, $29
			, $30
			, $31
			, $32
			, $33
			, $34
			, $35
			, $36
			, $37
			, $38
			, $39
			, $40
			, $41
			, $43
			, $44
			, $45
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SourceTable,
		entity.SourceId,
		entity.Authority,
		entity.CountryCode,
		entity.DocumentUuid,
		entity.DocumentType,
		entity.DocumentNumber,
		entity.InternalReference,
		entity.Status,
		entity.'draft',,
		entity.'pending',,
		entity.'submitted',,
		entity.'accepted',,
		entity.'rejected',,
		entity.'cancelled',,
		entity.'error',
		entity.SubmittedAt,
		entity.ResponseAt,
		entity.RequestPayload,
		entity.ResponsePayload,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.RetryCount,
		entity.LastRetryAt,
		entity.ZatcaHashValue,
		entity.ZatcaPreviousHashValue,
		entity.ZatcaInvoiceCounterValue,
		entity.ZatcaCryptographicStamp,
		entity.ZatcaQrCodePayload,
		entity.ZatcaComplianceInvoiceNumber,
		entity.EtaDocumentTypeVersion,
		entity.EtaSubmissionUuid,
		entity.EtaLongId,
		entity.EtaInternalId,
		entity.EtaDigitalSignature,
		entity.EtaSignatureAlgorithm,
		entity.SubmissionFormat,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create e_invoicing_documents", zap.Error(err))
		return fmt.Errorf("failed to create e_invoicing_documents: %w", err)
	}

	r.logger.Info("created e_invoicing_documents",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a e_invoicing_documents by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*EInvoicingDocuments, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "e_invoicing_documents", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, source_table
			, source_id
			, authority
			, country_code
			, document_uuid
			, document_type
			, document_number
			, internal_reference
			, status
			, 'draft',
			, 'pending',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error'
			, created_at
			, submitted_at
			, response_at
			, request_payload
			, response_payload
			, error_code
			, error_message
			, retry_count
			, last_retry_at
			, zatca_hash_value
			, zatca_previous_hash_value
			, zatca_invoice_counter_value
			, zatca_cryptographic_stamp
			, zatca_qr_code_payload
			, zatca_compliance_invoice_number
			, eta_document_type_version
			, eta_submission_uuid
			, eta_long_id
			, eta_internal_id
			, eta_digital_signature
			, eta_signature_algorithm
			, submission_format
			, metadata
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM e_invoicing_documents
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity EInvoicingDocuments
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SourceTable,
		&entity.SourceId,
		&entity.Authority,
		&entity.CountryCode,
		&entity.DocumentUuid,
		&entity.DocumentType,
		&entity.DocumentNumber,
		&entity.InternalReference,
		&entity.Status,
		&entity.'draft',,
		&entity.'pending',,
		&entity.'submitted',,
		&entity.'accepted',,
		&entity.'rejected',,
		&entity.'cancelled',,
		&entity.'error',
		&entity.CreatedAt,
		&entity.SubmittedAt,
		&entity.ResponseAt,
		&entity.RequestPayload,
		&entity.ResponsePayload,
		&entity.ErrorCode,
		&entity.ErrorMessage,
		&entity.RetryCount,
		&entity.LastRetryAt,
		&entity.ZatcaHashValue,
		&entity.ZatcaPreviousHashValue,
		&entity.ZatcaInvoiceCounterValue,
		&entity.ZatcaCryptographicStamp,
		&entity.ZatcaQrCodePayload,
		&entity.ZatcaComplianceInvoiceNumber,
		&entity.EtaDocumentTypeVersion,
		&entity.EtaSubmissionUuid,
		&entity.EtaLongId,
		&entity.EtaInternalId,
		&entity.EtaDigitalSignature,
		&entity.EtaSignatureAlgorithm,
		&entity.SubmissionFormat,
		&entity.Metadata,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("e_invoicing_documents not found")
	}

	if err != nil {
		r.logger.Error("failed to get e_invoicing_documents", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get e_invoicing_documents: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of e_invoicing_documents records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*EInvoicingDocuments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "e_invoicing_documents", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM e_invoicing_documents
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count e_invoicing_documents records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_table
			, source_id
			, authority
			, country_code
			, document_uuid
			, document_type
			, document_number
			, internal_reference
			, status
			, 'draft',
			, 'pending',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error'
			, created_at
			, submitted_at
			, response_at
			, request_payload
			, response_payload
			, error_code
			, error_message
			, retry_count
			, last_retry_at
			, zatca_hash_value
			, zatca_previous_hash_value
			, zatca_invoice_counter_value
			, zatca_cryptographic_stamp
			, zatca_qr_code_payload
			, zatca_compliance_invoice_number
			, eta_document_type_version
			, eta_submission_uuid
			, eta_long_id
			, eta_internal_id
			, eta_digital_signature
			, eta_signature_algorithm
			, submission_format
			, metadata
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM e_invoicing_documents
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list e_invoicing_documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list e_invoicing_documents: %w", err)
	}
	defer rows.Close()

	var entities []*EInvoicingDocuments
	for rows.Next() {
		var entity EInvoicingDocuments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceTable,
			&entity.SourceId,
			&entity.Authority,
			&entity.CountryCode,
			&entity.DocumentUuid,
			&entity.DocumentType,
			&entity.DocumentNumber,
			&entity.InternalReference,
			&entity.Status,
			&entity.'draft',,
			&entity.'pending',,
			&entity.'submitted',,
			&entity.'accepted',,
			&entity.'rejected',,
			&entity.'cancelled',,
			&entity.'error',
			&entity.CreatedAt,
			&entity.SubmittedAt,
			&entity.ResponseAt,
			&entity.RequestPayload,
			&entity.ResponsePayload,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.RetryCount,
			&entity.LastRetryAt,
			&entity.ZatcaHashValue,
			&entity.ZatcaPreviousHashValue,
			&entity.ZatcaInvoiceCounterValue,
			&entity.ZatcaCryptographicStamp,
			&entity.ZatcaQrCodePayload,
			&entity.ZatcaComplianceInvoiceNumber,
			&entity.EtaDocumentTypeVersion,
			&entity.EtaSubmissionUuid,
			&entity.EtaLongId,
			&entity.EtaInternalId,
			&entity.EtaDigitalSignature,
			&entity.EtaSignatureAlgorithm,
			&entity.SubmissionFormat,
			&entity.Metadata,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan e_invoicing_documents: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating e_invoicing_documents rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing e_invoicing_documents record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocuments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "e_invoicing_documents", duration, nil)
	}()

	query := `
		UPDATE e_invoicing_documents
		SET
			, organization_id = $2
			, source_table = $3
			, source_id = $4
			, authority = $5
			, country_code = $6
			, document_uuid = $7
			, document_type = $8
			, document_number = $9
			, internal_reference = $10
			, status = $11
			, 'draft', = $12
			, 'pending', = $13
			, 'submitted', = $14
			, 'accepted', = $15
			, 'rejected', = $16
			, 'cancelled', = $17
			, 'error' = $18
			, submitted_at = $20
			, response_at = $21
			, request_payload = $22
			, response_payload = $23
			, error_code = $24
			, error_message = $25
			, retry_count = $26
			, last_retry_at = $27
			, zatca_hash_value = $28
			, zatca_previous_hash_value = $29
			, zatca_invoice_counter_value = $30
			, zatca_cryptographic_stamp = $31
			, zatca_qr_code_payload = $32
			, zatca_compliance_invoice_number = $33
			, eta_document_type_version = $34
			, eta_submission_uuid = $35
			, eta_long_id = $36
			, eta_internal_id = $37
			, eta_digital_signature = $38
			, eta_signature_algorithm = $39
			, submission_format = $40
			, metadata = $41
			, updated_at = $42
			, deleted_at = $43
			, created_by = $44
			, updated_by = $45
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $46
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SourceTable,
		entity.SourceId,
		entity.Authority,
		entity.CountryCode,
		entity.DocumentUuid,
		entity.DocumentType,
		entity.DocumentNumber,
		entity.InternalReference,
		entity.Status,
		entity.'draft',,
		entity.'pending',,
		entity.'submitted',,
		entity.'accepted',,
		entity.'rejected',,
		entity.'cancelled',,
		entity.'error',
		entity.SubmittedAt,
		entity.ResponseAt,
		entity.RequestPayload,
		entity.ResponsePayload,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.RetryCount,
		entity.LastRetryAt,
		entity.ZatcaHashValue,
		entity.ZatcaPreviousHashValue,
		entity.ZatcaInvoiceCounterValue,
		entity.ZatcaCryptographicStamp,
		entity.ZatcaQrCodePayload,
		entity.ZatcaComplianceInvoiceNumber,
		entity.EtaDocumentTypeVersion,
		entity.EtaSubmissionUuid,
		entity.EtaLongId,
		entity.EtaInternalId,
		entity.EtaDigitalSignature,
		entity.EtaSignatureAlgorithm,
		entity.SubmissionFormat,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update e_invoicing_documents", zap.Error(err))
		return fmt.Errorf("failed to update e_invoicing_documents: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("e_invoicing_documents not found or already deleted")
	}

	r.logger.Info("updated e_invoicing_documents",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a e_invoicing_documents record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "e_invoicing_documents", duration, nil)
	}()

	query := `
		UPDATE e_invoicing_documents
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete e_invoicing_documents", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete e_invoicing_documents: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("e_invoicing_documents not found or already deleted")
	}

	r.logger.Info("deleted e_invoicing_documents", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves e_invoicing_documents records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*EInvoicingDocuments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "e_invoicing_documents", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM e_invoicing_documents
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count e_invoicing_documents records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_table
			, source_id
			, authority
			, country_code
			, document_uuid
			, document_type
			, document_number
			, internal_reference
			, status
			, 'draft',
			, 'pending',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error'
			, created_at
			, submitted_at
			, response_at
			, request_payload
			, response_payload
			, error_code
			, error_message
			, retry_count
			, last_retry_at
			, zatca_hash_value
			, zatca_previous_hash_value
			, zatca_invoice_counter_value
			, zatca_cryptographic_stamp
			, zatca_qr_code_payload
			, zatca_compliance_invoice_number
			, eta_document_type_version
			, eta_submission_uuid
			, eta_long_id
			, eta_internal_id
			, eta_digital_signature
			, eta_signature_algorithm
			, submission_format
			, metadata
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM e_invoicing_documents
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list e_invoicing_documents by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list e_invoicing_documents: %w", err)
	}
	defer rows.Close()

	var entities []*EInvoicingDocuments
	for rows.Next() {
		var entity EInvoicingDocuments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceTable,
			&entity.SourceId,
			&entity.Authority,
			&entity.CountryCode,
			&entity.DocumentUuid,
			&entity.DocumentType,
			&entity.DocumentNumber,
			&entity.InternalReference,
			&entity.Status,
			&entity.'draft',,
			&entity.'pending',,
			&entity.'submitted',,
			&entity.'accepted',,
			&entity.'rejected',,
			&entity.'cancelled',,
			&entity.'error',
			&entity.CreatedAt,
			&entity.SubmittedAt,
			&entity.ResponseAt,
			&entity.RequestPayload,
			&entity.ResponsePayload,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.RetryCount,
			&entity.LastRetryAt,
			&entity.ZatcaHashValue,
			&entity.ZatcaPreviousHashValue,
			&entity.ZatcaInvoiceCounterValue,
			&entity.ZatcaCryptographicStamp,
			&entity.ZatcaQrCodePayload,
			&entity.ZatcaComplianceInvoiceNumber,
			&entity.EtaDocumentTypeVersion,
			&entity.EtaSubmissionUuid,
			&entity.EtaLongId,
			&entity.EtaInternalId,
			&entity.EtaDigitalSignature,
			&entity.EtaSignatureAlgorithm,
			&entity.SubmissionFormat,
			&entity.Metadata,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan e_invoicing_documents: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


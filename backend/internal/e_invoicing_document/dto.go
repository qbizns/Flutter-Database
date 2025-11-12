package e_invoicing_document

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EInvoicingDocumentsResponse represents a e_invoicing_documents response
type EInvoicingDocumentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SourceTable string `json:"source_table"`
	
	SourceId uuid.UUID `json:"source_id"`
	
	Authority string `json:"authority"`
	
	CountryCode string `json:"country_code"`
	
	DocumentUuid uuid.UUID `json:"document_uuid"`
	
	DocumentType string `json:"document_type"`
	
	DocumentNumber string `json:"document_number"`
	
	InternalReference *string `json:"internal_reference"`
	
	Status string `json:"status"`
	
	'draft', *string `json:"'draft',"`
	
	'pending', *string `json:"'pending',"`
	
	'submitted', *string `json:"'submitted',"`
	
	'accepted', *string `json:"'accepted',"`
	
	'rejected', *string `json:"'rejected',"`
	
	'cancelled', *string `json:"'cancelled',"`
	
	'error' *string `json:"'error'"`
	
	CreatedAt time.Time `json:"created_at"`
	
	SubmittedAt *time.Time `json:"submitted_at"`
	
	ResponseAt *time.Time `json:"response_at"`
	
	RequestPayload json.RawMessage `json:"request_payload"`
	
	ResponsePayload json.RawMessage `json:"response_payload"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage *string `json:"error_message"`
	
	RetryCount *int64 `json:"retry_count"`
	
	LastRetryAt *time.Time `json:"last_retry_at"`
	
	ZatcaHashValue *string `json:"zatca_hash_value"`
	
	ZatcaPreviousHashValue *string `json:"zatca_previous_hash_value"`
	
	ZatcaInvoiceCounterValue *int64 `json:"zatca_invoice_counter_value"`
	
	ZatcaCryptographicStamp *string `json:"zatca_cryptographic_stamp"`
	
	ZatcaQrCodePayload *string `json:"zatca_qr_code_payload"`
	
	ZatcaComplianceInvoiceNumber *string `json:"zatca_compliance_invoice_number"`
	
	EtaDocumentTypeVersion *string `json:"eta_document_type_version"`
	
	EtaSubmissionUuid *uuid.UUID `json:"eta_submission_uuid"`
	
	EtaLongId *string `json:"eta_long_id"`
	
	EtaInternalId *string `json:"eta_internal_id"`
	
	EtaDigitalSignature *string `json:"eta_digital_signature"`
	
	EtaSignatureAlgorithm *string `json:"eta_signature_algorithm"`
	
	SubmissionFormat *string `json:"submission_format"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateEInvoicingDocumentsRequest represents a request to create a e_invoicing_documents
type CreateEInvoicingDocumentsRequest struct {
	
	SourceTable string `json:"source_table" validate:"required"`
	
	SourceId uuid.UUID `json:"source_id" validate:"required"`
	
	Authority string `json:"authority" validate:"required"`
	
	CountryCode string `json:"country_code" validate:"required"`
	
	DocumentUuid uuid.UUID `json:"document_uuid" validate:"required"`
	
	DocumentType string `json:"document_type" validate:"required"`
	
	DocumentNumber string `json:"document_number" validate:"required"`
	
	InternalReference *string `json:"internal_reference"`
	
	Status string `json:"status" validate:"required"`
	
	'draft', *string `json:"'draft',"`
	
	'pending', *string `json:"'pending',"`
	
	'submitted', *string `json:"'submitted',"`
	
	'accepted', *string `json:"'accepted',"`
	
	'rejected', *string `json:"'rejected',"`
	
	'cancelled', *string `json:"'cancelled',"`
	
	'error' *string `json:"'error'"`
	
	SubmittedAt *time.Time `json:"submitted_at"`
	
	ResponseAt *time.Time `json:"response_at"`
	
	RequestPayload json.RawMessage `json:"request_payload"`
	
	ResponsePayload json.RawMessage `json:"response_payload"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage *string `json:"error_message"`
	
	RetryCount *int64 `json:"retry_count"`
	
	LastRetryAt *time.Time `json:"last_retry_at"`
	
	ZatcaHashValue *string `json:"zatca_hash_value"`
	
	ZatcaPreviousHashValue *string `json:"zatca_previous_hash_value"`
	
	ZatcaInvoiceCounterValue *int64 `json:"zatca_invoice_counter_value"`
	
	ZatcaCryptographicStamp *string `json:"zatca_cryptographic_stamp"`
	
	ZatcaQrCodePayload *string `json:"zatca_qr_code_payload"`
	
	ZatcaComplianceInvoiceNumber *string `json:"zatca_compliance_invoice_number"`
	
	EtaDocumentTypeVersion *string `json:"eta_document_type_version"`
	
	EtaSubmissionUuid *uuid.UUID `json:"eta_submission_uuid"`
	
	EtaLongId *string `json:"eta_long_id"`
	
	EtaInternalId *string `json:"eta_internal_id"`
	
	EtaDigitalSignature *string `json:"eta_digital_signature"`
	
	EtaSignatureAlgorithm *string `json:"eta_signature_algorithm"`
	
	SubmissionFormat *string `json:"submission_format"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateEInvoicingDocumentsRequest) Validate() error {
	
	if r.SourceTable == "" {
		return fmt.Errorf("source_table is required")
	}
	
	if r.SourceId == uuid.Nil {
		return fmt.Errorf("source_id is required")
	}
	
	if r.Authority == "" {
		return fmt.Errorf("authority is required")
	}
	
	if r.CountryCode == "" {
		return fmt.Errorf("country_code is required")
	}
	
	if r.DocumentUuid == uuid.Nil {
		return fmt.Errorf("document_uuid is required")
	}
	
	if r.DocumentType == "" {
		return fmt.Errorf("document_type is required")
	}
	
	if r.DocumentNumber == "" {
		return fmt.Errorf("document_number is required")
	}
	
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateEInvoicingDocumentsRequest represents a request to update a e_invoicing_documents
type UpdateEInvoicingDocumentsRequest struct {
	
	SourceTable *string `json:"source_table,omitempty" validate:"omitempty,required"`
	
	SourceId *uuid.UUID `json:"source_id,omitempty" validate:"omitempty,required"`
	
	Authority *string `json:"authority,omitempty" validate:"omitempty,required"`
	
	CountryCode *string `json:"country_code,omitempty" validate:"omitempty,required"`
	
	DocumentUuid *uuid.UUID `json:"document_uuid,omitempty" validate:"omitempty,required"`
	
	DocumentType *string `json:"document_type,omitempty" validate:"omitempty,required"`
	
	DocumentNumber *string `json:"document_number,omitempty" validate:"omitempty,required"`
	
	InternalReference *string `json:"internal_reference,omitempty"`
	
	Status *string `json:"status,omitempty" validate:"omitempty,required"`
	
	'draft', *string `json:"'draft',,omitempty"`
	
	'pending', *string `json:"'pending',,omitempty"`
	
	'submitted', *string `json:"'submitted',,omitempty"`
	
	'accepted', *string `json:"'accepted',,omitempty"`
	
	'rejected', *string `json:"'rejected',,omitempty"`
	
	'cancelled', *string `json:"'cancelled',,omitempty"`
	
	'error' *string `json:"'error',omitempty"`
	
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	
	ResponseAt *time.Time `json:"response_at,omitempty"`
	
	RequestPayload *json.RawMessage `json:"request_payload,omitempty"`
	
	ResponsePayload *json.RawMessage `json:"response_payload,omitempty"`
	
	ErrorCode *string `json:"error_code,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	RetryCount *int64 `json:"retry_count,omitempty"`
	
	LastRetryAt *time.Time `json:"last_retry_at,omitempty"`
	
	ZatcaHashValue *string `json:"zatca_hash_value,omitempty"`
	
	ZatcaPreviousHashValue *string `json:"zatca_previous_hash_value,omitempty"`
	
	ZatcaInvoiceCounterValue *int64 `json:"zatca_invoice_counter_value,omitempty"`
	
	ZatcaCryptographicStamp *string `json:"zatca_cryptographic_stamp,omitempty"`
	
	ZatcaQrCodePayload *string `json:"zatca_qr_code_payload,omitempty"`
	
	ZatcaComplianceInvoiceNumber *string `json:"zatca_compliance_invoice_number,omitempty"`
	
	EtaDocumentTypeVersion *string `json:"eta_document_type_version,omitempty"`
	
	EtaSubmissionUuid *uuid.UUID `json:"eta_submission_uuid,omitempty"`
	
	EtaLongId *string `json:"eta_long_id,omitempty"`
	
	EtaInternalId *string `json:"eta_internal_id,omitempty"`
	
	EtaDigitalSignature *string `json:"eta_digital_signature,omitempty"`
	
	EtaSignatureAlgorithm *string `json:"eta_signature_algorithm,omitempty"`
	
	SubmissionFormat *string `json:"submission_format,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateEInvoicingDocumentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SourceTable != nil {
		hasUpdate = true
	}
	
	if r.SourceId != nil {
		hasUpdate = true
	}
	
	if r.Authority != nil {
		hasUpdate = true
	}
	
	if r.CountryCode != nil {
		hasUpdate = true
	}
	
	if r.DocumentUuid != nil {
		hasUpdate = true
	}
	
	if r.DocumentType != nil {
		hasUpdate = true
	}
	
	if r.DocumentNumber != nil {
		hasUpdate = true
	}
	
	if r.InternalReference != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.'draft', != nil {
		hasUpdate = true
	}
	
	if r.'pending', != nil {
		hasUpdate = true
	}
	
	if r.'submitted', != nil {
		hasUpdate = true
	}
	
	if r.'accepted', != nil {
		hasUpdate = true
	}
	
	if r.'rejected', != nil {
		hasUpdate = true
	}
	
	if r.'cancelled', != nil {
		hasUpdate = true
	}
	
	if r.'error' != nil {
		hasUpdate = true
	}
	
	if r.SubmittedAt != nil {
		hasUpdate = true
	}
	
	if r.ResponseAt != nil {
		hasUpdate = true
	}
	
	if r.RequestPayload != nil {
		hasUpdate = true
	}
	
	if r.ResponsePayload != nil {
		hasUpdate = true
	}
	
	if r.ErrorCode != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.RetryCount != nil {
		hasUpdate = true
	}
	
	if r.LastRetryAt != nil {
		hasUpdate = true
	}
	
	if r.ZatcaHashValue != nil {
		hasUpdate = true
	}
	
	if r.ZatcaPreviousHashValue != nil {
		hasUpdate = true
	}
	
	if r.ZatcaInvoiceCounterValue != nil {
		hasUpdate = true
	}
	
	if r.ZatcaCryptographicStamp != nil {
		hasUpdate = true
	}
	
	if r.ZatcaQrCodePayload != nil {
		hasUpdate = true
	}
	
	if r.ZatcaComplianceInvoiceNumber != nil {
		hasUpdate = true
	}
	
	if r.EtaDocumentTypeVersion != nil {
		hasUpdate = true
	}
	
	if r.EtaSubmissionUuid != nil {
		hasUpdate = true
	}
	
	if r.EtaLongId != nil {
		hasUpdate = true
	}
	
	if r.EtaInternalId != nil {
		hasUpdate = true
	}
	
	if r.EtaDigitalSignature != nil {
		hasUpdate = true
	}
	
	if r.EtaSignatureAlgorithm != nil {
		hasUpdate = true
	}
	
	if r.SubmissionFormat != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// EInvoicingDocumentsListResponse represents a paginated list of e_invoicing_documents records
type EInvoicingDocumentsListResponse struct {
	Items      []*EInvoicingDocumentsResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

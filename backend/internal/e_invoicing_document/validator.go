package e_invoicing_document

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles EInvoicingDocuments validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new EInvoicingDocuments validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateEInvoicingDocumentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate SourceTable
	
	if err := v.validateSourceTable(req.SourceTable); err != nil {
		return err
	}
	
	
	
	// Validate SourceId
	
	
	if err := v.validateSourceIdExists(ctx, tx, req.SourceId); err != nil {
		return err
	}
	
	
	// Validate Authority
	
	if err := v.validateAuthority(req.Authority); err != nil {
		return err
	}
	
	
	
	// Validate CountryCode
	
	if err := v.validateCountryCode(req.CountryCode); err != nil {
		return err
	}
	
	
	
	// Validate DocumentUuid
	
	
	
	// Validate DocumentType
	
	if err := v.validateDocumentType(req.DocumentType); err != nil {
		return err
	}
	
	
	
	// Validate DocumentNumber
	
	if err := v.validateDocumentNumber(req.DocumentNumber); err != nil {
		return err
	}
	
	
	
	// Validate InternalReference
	
	
	
	// Validate Status
	
	if err := v.validateStatus(req.Status); err != nil {
		return err
	}
	
	
	
	// Validate 'draft',
	
	
	
	// Validate 'pending',
	
	
	
	// Validate 'submitted',
	
	
	
	// Validate 'accepted',
	
	
	
	// Validate 'rejected',
	
	
	
	// Validate 'cancelled',
	
	
	
	// Validate 'error'
	
	
	
	// Validate SubmittedAt
	
	
	
	// Validate ResponseAt
	
	
	
	// Validate RequestPayload
	
	
	
	// Validate ResponsePayload
	
	
	
	// Validate ErrorCode
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate RetryCount
	
	
	
	// Validate LastRetryAt
	
	
	
	// Validate ZatcaHashValue
	
	
	
	// Validate ZatcaPreviousHashValue
	
	
	
	// Validate ZatcaInvoiceCounterValue
	
	
	
	// Validate ZatcaCryptographicStamp
	
	
	
	// Validate ZatcaQrCodePayload
	
	
	
	// Validate ZatcaComplianceInvoiceNumber
	
	
	
	// Validate EtaDocumentTypeVersion
	
	
	
	// Validate EtaSubmissionUuid
	
	
	
	// Validate EtaLongId
	
	
	if err := v.validateEtaLongIdExists(ctx, tx, req.EtaLongId); err != nil {
		return err
	}
	
	
	// Validate EtaInternalId
	
	
	if err := v.validateEtaInternalIdExists(ctx, tx, req.EtaInternalId); err != nil {
		return err
	}
	
	
	// Validate EtaDigitalSignature
	
	
	
	// Validate EtaSignatureAlgorithm
	
	
	
	// Validate SubmissionFormat
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	

	// Cross-field validation
	if err := v.validateCrossFields(ctx, tx, req); err != nil {
		return err
	}

	// Business rules validation
	if err := v.validateBusinessRules(ctx, tx, req); err != nil {
		return err
	}

	return nil
}

// ValidateUpdate validates an update request
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateEInvoicingDocumentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate SourceTable if provided
	
	if req.SourceTable != nil {
		if err := v.validateSourceTable(*req.SourceTable); err != nil {
			return err
		}
	}
	
	
	
	// Validate SourceId if provided
	
	
	if req.SourceId != nil {
		if err := v.validateSourceIdExists(ctx, tx, *req.SourceId); err != nil {
			return err
		}
	}
	
	
	// Validate Authority if provided
	
	if req.Authority != nil {
		if err := v.validateAuthority(*req.Authority); err != nil {
			return err
		}
	}
	
	
	
	// Validate CountryCode if provided
	
	if req.CountryCode != nil {
		if err := v.validateCountryCode(*req.CountryCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate DocumentUuid if provided
	
	
	
	// Validate DocumentType if provided
	
	if req.DocumentType != nil {
		if err := v.validateDocumentType(*req.DocumentType); err != nil {
			return err
		}
	}
	
	
	
	// Validate DocumentNumber if provided
	
	if req.DocumentNumber != nil {
		if err := v.validateDocumentNumber(*req.DocumentNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate InternalReference if provided
	
	
	
	// Validate Status if provided
	
	if req.Status != nil {
		if err := v.validateStatus(*req.Status); err != nil {
			return err
		}
	}
	
	
	
	// Validate 'draft', if provided
	
	
	
	// Validate 'pending', if provided
	
	
	
	// Validate 'submitted', if provided
	
	
	
	// Validate 'accepted', if provided
	
	
	
	// Validate 'rejected', if provided
	
	
	
	// Validate 'cancelled', if provided
	
	
	
	// Validate 'error' if provided
	
	
	
	// Validate SubmittedAt if provided
	
	
	
	// Validate ResponseAt if provided
	
	
	
	// Validate RequestPayload if provided
	
	
	
	// Validate ResponsePayload if provided
	
	
	
	// Validate ErrorCode if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate RetryCount if provided
	
	
	
	// Validate LastRetryAt if provided
	
	
	
	// Validate ZatcaHashValue if provided
	
	
	
	// Validate ZatcaPreviousHashValue if provided
	
	
	
	// Validate ZatcaInvoiceCounterValue if provided
	
	
	
	// Validate ZatcaCryptographicStamp if provided
	
	
	
	// Validate ZatcaQrCodePayload if provided
	
	
	
	// Validate ZatcaComplianceInvoiceNumber if provided
	
	
	
	// Validate EtaDocumentTypeVersion if provided
	
	
	
	// Validate EtaSubmissionUuid if provided
	
	
	
	// Validate EtaLongId if provided
	
	
	if req.EtaLongId != nil {
		if err := v.validateEtaLongIdExists(ctx, tx, *req.EtaLongId); err != nil {
			return err
		}
	}
	
	
	// Validate EtaInternalId if provided
	
	
	if req.EtaInternalId != nil {
		if err := v.validateEtaInternalIdExists(ctx, tx, *req.EtaInternalId); err != nil {
			return err
		}
	}
	
	
	// Validate EtaDigitalSignature if provided
	
	
	
	// Validate EtaSignatureAlgorithm if provided
	
	
	
	// Validate SubmissionFormat if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	

	// Business rules validation
	if err := v.validateUpdateBusinessRules(ctx, tx, existing, req); err != nil {
		return err
	}

	return nil
}

// ValidateDelete validates a delete request
func (v *Validator) ValidateDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	// Check if entity can be deleted (no foreign key constraints)
	if err := v.validateCanDelete(ctx, tx, existing); err != nil {
		return err
	}

	return nil
}



// validateSourceTable validates source_table field
func (v *Validator) validateSourceTable(value string) error {
	
	// Add custom validation for source_table
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("source_table cannot be empty")
	}
	
	return nil
}







// validateSourceIdExists validates that source_id exists
func (v *Validator) validateSourceIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for source
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM source WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check source existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("source with id %s does not exist", id)
	}
	return nil
}



// validateAuthority validates authority field
func (v *Validator) validateAuthority(value string) error {
	
	// Add custom validation for authority
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("authority cannot be empty")
	}
	
	return nil
}





// validateCountryCode validates country_code field
func (v *Validator) validateCountryCode(value string) error {
	
	// Add custom validation for country_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("country_code cannot be empty")
	}
	
	return nil
}









// validateDocumentType validates document_type field
func (v *Validator) validateDocumentType(value string) error {
	
	// Add custom validation for document_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("document_type cannot be empty")
	}
	
	return nil
}





// validateDocumentNumber validates document_number field
func (v *Validator) validateDocumentNumber(value string) error {
	
	// Add custom validation for document_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("document_number cannot be empty")
	}
	
	return nil
}









// validateStatus validates status field
func (v *Validator) validateStatus(value string) error {
	
	// Add custom validation for status
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("status cannot be empty")
	}
	
	return nil
}



































































































// validateEtaLongIdExists validates that eta_long_id exists
func (v *Validator) validateEtaLongIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for eta_long
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM eta_long WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check eta_long existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("eta_long with id %s does not exist", id)
	}
	return nil
}





// validateEtaInternalIdExists validates that eta_internal_id exists
func (v *Validator) validateEtaInternalIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for eta_internal
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM eta_internal WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check eta_internal existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("eta_internal with id %s does not exist", id)
	}
	return nil
}



























// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateEInvoicingDocumentsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateEInvoicingDocumentsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *EInvoicingDocuments, req *UpdateEInvoicingDocumentsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocuments) error {
	// Add delete validation here
	// Example: check for dependent records in other tables
	// Example: prevent deletion of active/in-use entities
	return nil
}

// Helper validation functions

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`) // E.164 format
	urlRegex   = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
)

// isValidEmail validates email format
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// isValidPhone validates phone number format (E.164)
func isValidPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// isValidURL validates URL format
func isValidURL(url string) bool {
	return urlRegex.MatchString(url)
}

// isValidUUID validates UUID format
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// isValidDateRange validates date range
func isValidDateRange(start, end time.Time) bool {
	return start.Before(end)
}

// isPositive validates positive numbers
func isPositive(value float64) bool {
	return value > 0
}

// isNonNegative validates non-negative numbers
func isNonNegative(value float64) bool {
	return value >= 0
}

// isWithinRange validates value is within range
func isWithinRange(value, min, max float64) bool {
	return value >= min && value <= max
}

// isValidLength validates string length
func isValidLength(value string, min, max int) bool {
	length := len(value)
	return length >= min && length <= max
}

package e_invoicing_document_event

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles EInvoicingDocumentEvents validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new EInvoicingDocumentEvents validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateEInvoicingDocumentEventsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate EInvoicingDocumentId
	
	
	if err := v.validateEInvoicingDocumentIdExists(ctx, tx, req.EInvoicingDocumentId); err != nil {
		return err
	}
	
	
	// Validate EventType
	
	if err := v.validateEventType(req.EventType); err != nil {
		return err
	}
	
	
	
	// Validate 'created',
	
	
	
	// Validate 'validated',
	
	
	
	// Validate 'submitted',
	
	
	
	// Validate 'accepted',
	
	
	
	// Validate 'rejected',
	
	
	
	// Validate 'cancelled',
	
	
	
	// Validate 'error',
	
	
	
	// Validate 'retry',
	
	
	
	// Validate 'statusCheck'
	
	
	
	// Validate EventTimestamp
	
	
	
	// Validate PreviousStatus
	
	
	
	// Validate NewStatus
	
	
	
	// Validate EventDescription
	
	
	
	// Validate EventData
	
	
	
	// Validate HttpStatusCode
	
	
	
	// Validate HttpMethod
	
	
	
	// Validate ApiEndpoint
	
	
	
	// Validate RequestHeaders
	
	
	
	// Validate ResponseHeaders
	
	
	
	// Validate ErrorCode
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate ErrorDetails
	
	
	
	// Validate TriggeredBy
	
	if err := v.validateTriggeredBy(req.TriggeredBy); err != nil {
		return err
	}
	
	
	
	// Validate UserId
	
	
	if err := v.validateUserIdExists(ctx, tx, req.UserId); err != nil {
		return err
	}
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateEInvoicingDocumentEventsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate EInvoicingDocumentId if provided
	
	
	if req.EInvoicingDocumentId != nil {
		if err := v.validateEInvoicingDocumentIdExists(ctx, tx, *req.EInvoicingDocumentId); err != nil {
			return err
		}
	}
	
	
	// Validate EventType if provided
	
	if req.EventType != nil {
		if err := v.validateEventType(*req.EventType); err != nil {
			return err
		}
	}
	
	
	
	// Validate 'created', if provided
	
	
	
	// Validate 'validated', if provided
	
	
	
	// Validate 'submitted', if provided
	
	
	
	// Validate 'accepted', if provided
	
	
	
	// Validate 'rejected', if provided
	
	
	
	// Validate 'cancelled', if provided
	
	
	
	// Validate 'error', if provided
	
	
	
	// Validate 'retry', if provided
	
	
	
	// Validate 'statusCheck' if provided
	
	
	
	// Validate EventTimestamp if provided
	
	
	
	// Validate PreviousStatus if provided
	
	
	
	// Validate NewStatus if provided
	
	
	
	// Validate EventDescription if provided
	
	
	
	// Validate EventData if provided
	
	
	
	// Validate HttpStatusCode if provided
	
	
	
	// Validate HttpMethod if provided
	
	
	
	// Validate ApiEndpoint if provided
	
	
	
	// Validate RequestHeaders if provided
	
	
	
	// Validate ResponseHeaders if provided
	
	
	
	// Validate ErrorCode if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate ErrorDetails if provided
	
	
	
	// Validate TriggeredBy if provided
	
	if req.TriggeredBy != nil {
		if err := v.validateTriggeredBy(*req.TriggeredBy); err != nil {
			return err
		}
	}
	
	
	
	// Validate UserId if provided
	
	
	if req.UserId != nil {
		if err := v.validateUserIdExists(ctx, tx, *req.UserId); err != nil {
			return err
		}
	}
	
	

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





// validateEInvoicingDocumentIdExists validates that e_invoicing_document_id exists
func (v *Validator) validateEInvoicingDocumentIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for e_invoicing_document
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM e_invoicing_document WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check e_invoicing_document existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("e_invoicing_document with id %s does not exist", id)
	}
	return nil
}



// validateEventType validates event_type field
func (v *Validator) validateEventType(value string) error {
	
	// Add custom validation for event_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("event_type cannot be empty")
	}
	
	return nil
}





























































































// validateTriggeredBy validates triggered_by field
func (v *Validator) validateTriggeredBy(value string) error {
	
	// Add custom validation for triggered_by
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("triggered_by cannot be empty")
	}
	
	return nil
}







// validateUserIdExists validates that user_id exists
func (v *Validator) validateUserIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for user
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM user WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %s does not exist", id)
	}
	return nil
}



// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateEInvoicingDocumentEventsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateEInvoicingDocumentEventsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *EInvoicingDocumentEvents, req *UpdateEInvoicingDocumentEventsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocumentEvents) error {
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

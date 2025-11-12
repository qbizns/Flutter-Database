package posting_validation_result

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles PostingValidationResults validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new PostingValidationResults validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreatePostingValidationResultsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate DocumentTypeCode
	
	if err := v.validateDocumentTypeCode(req.DocumentTypeCode); err != nil {
		return err
	}
	
	
	
	// Validate DocumentId
	
	
	if err := v.validateDocumentIdExists(ctx, tx, req.DocumentId); err != nil {
		return err
	}
	
	
	// Validate Event
	
	if err := v.validateEvent(req.Event); err != nil {
		return err
	}
	
	
	
	// Validate JournalEntryId
	
	
	if err := v.validateJournalEntryIdExists(ctx, tx, req.JournalEntryId); err != nil {
		return err
	}
	
	
	// Validate ValidationRuleId
	
	
	if err := v.validateValidationRuleIdExists(ctx, tx, req.ValidationRuleId); err != nil {
		return err
	}
	
	
	// Validate Severity
	
	if err := v.validateSeverity(req.Severity); err != nil {
		return err
	}
	
	
	
	// Validate MessageCode
	
	if err := v.validateMessageCode(req.MessageCode); err != nil {
		return err
	}
	
	
	
	// Validate Message
	
	if err := v.validateMessage(req.Message); err != nil {
		return err
	}
	
	
	
	// Validate IsBlocking
	
	
	
	// Validate Context
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdatePostingValidationResultsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate DocumentTypeCode if provided
	
	if req.DocumentTypeCode != nil {
		if err := v.validateDocumentTypeCode(*req.DocumentTypeCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate DocumentId if provided
	
	
	if req.DocumentId != nil {
		if err := v.validateDocumentIdExists(ctx, tx, *req.DocumentId); err != nil {
			return err
		}
	}
	
	
	// Validate Event if provided
	
	if req.Event != nil {
		if err := v.validateEvent(*req.Event); err != nil {
			return err
		}
	}
	
	
	
	// Validate JournalEntryId if provided
	
	
	if req.JournalEntryId != nil {
		if err := v.validateJournalEntryIdExists(ctx, tx, *req.JournalEntryId); err != nil {
			return err
		}
	}
	
	
	// Validate ValidationRuleId if provided
	
	
	if req.ValidationRuleId != nil {
		if err := v.validateValidationRuleIdExists(ctx, tx, *req.ValidationRuleId); err != nil {
			return err
		}
	}
	
	
	// Validate Severity if provided
	
	if req.Severity != nil {
		if err := v.validateSeverity(*req.Severity); err != nil {
			return err
		}
	}
	
	
	
	// Validate MessageCode if provided
	
	if req.MessageCode != nil {
		if err := v.validateMessageCode(*req.MessageCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate Message if provided
	
	if req.Message != nil {
		if err := v.validateMessage(*req.Message); err != nil {
			return err
		}
	}
	
	
	
	// Validate IsBlocking if provided
	
	
	
	// Validate Context if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	

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



// validateDocumentTypeCode validates document_type_code field
func (v *Validator) validateDocumentTypeCode(value string) error {
	
	// Add custom validation for document_type_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("document_type_code cannot be empty")
	}
	
	return nil
}







// validateDocumentIdExists validates that document_id exists
func (v *Validator) validateDocumentIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for document
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM document WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check document existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("document with id %s does not exist", id)
	}
	return nil
}



// validateEvent validates event field
func (v *Validator) validateEvent(value string) error {
	
	// Add custom validation for event
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("event cannot be empty")
	}
	
	return nil
}







// validateJournalEntryIdExists validates that journal_entry_id exists
func (v *Validator) validateJournalEntryIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for journal_entry
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM journal_entry WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check journal_entry existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("journal_entry with id %s does not exist", id)
	}
	return nil
}





// validateValidationRuleIdExists validates that validation_rule_id exists
func (v *Validator) validateValidationRuleIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for validation_rule
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM validation_rule WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check validation_rule existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("validation_rule with id %s does not exist", id)
	}
	return nil
}



// validateSeverity validates severity field
func (v *Validator) validateSeverity(value string) error {
	
	// Add custom validation for severity
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("severity cannot be empty")
	}
	
	return nil
}





// validateMessageCode validates message_code field
func (v *Validator) validateMessageCode(value string) error {
	
	// Add custom validation for message_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("message_code cannot be empty")
	}
	
	return nil
}





// validateMessage validates message field
func (v *Validator) validateMessage(value string) error {
	
	// Add custom validation for message
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("message cannot be empty")
	}
	
	return nil
}





















// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreatePostingValidationResultsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreatePostingValidationResultsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *PostingValidationResults, req *UpdatePostingValidationResultsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *PostingValidationResults) error {
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

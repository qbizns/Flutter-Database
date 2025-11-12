package immutability_violations_log

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles ImmutabilityViolationsLog validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ImmutabilityViolationsLog validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateImmutabilityViolationsLogRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate TableName
	
	if err := v.validateTableName(req.TableName); err != nil {
		return err
	}
	
	
	
	// Validate RecordId
	
	
	if err := v.validateRecordIdExists(ctx, tx, req.RecordId); err != nil {
		return err
	}
	
	
	// Validate Operation
	
	if err := v.validateOperation(req.Operation); err != nil {
		return err
	}
	
	
	
	// Validate AttemptedBy
	
	
	
	// Validate AttemptedAt
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate BlockedData
	
	
	
	// Validate Metadata
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateImmutabilityViolationsLogRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate TableName if provided
	
	if req.TableName != nil {
		if err := v.validateTableName(*req.TableName); err != nil {
			return err
		}
	}
	
	
	
	// Validate RecordId if provided
	
	
	if req.RecordId != nil {
		if err := v.validateRecordIdExists(ctx, tx, *req.RecordId); err != nil {
			return err
		}
	}
	
	
	// Validate Operation if provided
	
	if req.Operation != nil {
		if err := v.validateOperation(*req.Operation); err != nil {
			return err
		}
	}
	
	
	
	// Validate AttemptedBy if provided
	
	
	
	// Validate AttemptedAt if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate BlockedData if provided
	
	
	
	// Validate Metadata if provided
	
	
	

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



// validateTableName validates table_name field
func (v *Validator) validateTableName(value string) error {
	
	// Add custom validation for table_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("table_name cannot be empty")
	}
	
	return nil
}







// validateRecordIdExists validates that record_id exists
func (v *Validator) validateRecordIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for record
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM record WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check record existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("record with id %s does not exist", id)
	}
	return nil
}



// validateOperation validates operation field
func (v *Validator) validateOperation(value string) error {
	
	// Add custom validation for operation
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("operation cannot be empty")
	}
	
	return nil
}

























// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateImmutabilityViolationsLogRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateImmutabilityViolationsLogRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ImmutabilityViolationsLog, req *UpdateImmutabilityViolationsLogRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ImmutabilityViolationsLog) error {
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

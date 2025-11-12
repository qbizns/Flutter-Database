package fiscal_position_tax_mapping

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles FiscalPositionTaxMappings validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new FiscalPositionTaxMappings validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateFiscalPositionTaxMappingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate FiscalPositionId
	
	
	if err := v.validateFiscalPositionIdExists(ctx, tx, req.FiscalPositionId); err != nil {
		return err
	}
	
	
	// Validate SourceTaxId
	
	
	if err := v.validateSourceTaxIdExists(ctx, tx, req.SourceTaxId); err != nil {
		return err
	}
	
	
	// Validate DestinationTaxId
	
	
	if err := v.validateDestinationTaxIdExists(ctx, tx, req.DestinationTaxId); err != nil {
		return err
	}
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateFiscalPositionTaxMappingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate FiscalPositionId if provided
	
	
	if req.FiscalPositionId != nil {
		if err := v.validateFiscalPositionIdExists(ctx, tx, *req.FiscalPositionId); err != nil {
			return err
		}
	}
	
	
	// Validate SourceTaxId if provided
	
	
	if req.SourceTaxId != nil {
		if err := v.validateSourceTaxIdExists(ctx, tx, *req.SourceTaxId); err != nil {
			return err
		}
	}
	
	
	// Validate DestinationTaxId if provided
	
	
	if req.DestinationTaxId != nil {
		if err := v.validateDestinationTaxIdExists(ctx, tx, *req.DestinationTaxId); err != nil {
			return err
		}
	}
	
	
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





// validateFiscalPositionIdExists validates that fiscal_position_id exists
func (v *Validator) validateFiscalPositionIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for fiscal_position
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM fiscal_position WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check fiscal_position existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("fiscal_position with id %s does not exist", id)
	}
	return nil
}





// validateSourceTaxIdExists validates that source_tax_id exists
func (v *Validator) validateSourceTaxIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for source_tax
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM source_tax WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check source_tax existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("source_tax with id %s does not exist", id)
	}
	return nil
}





// validateDestinationTaxIdExists validates that destination_tax_id exists
func (v *Validator) validateDestinationTaxIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for destination_tax
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM destination_tax WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check destination_tax existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("destination_tax with id %s does not exist", id)
	}
	return nil
}







// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateFiscalPositionTaxMappingsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateFiscalPositionTaxMappingsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *FiscalPositionTaxMappings, req *UpdateFiscalPositionTaxMappingsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *FiscalPositionTaxMappings) error {
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

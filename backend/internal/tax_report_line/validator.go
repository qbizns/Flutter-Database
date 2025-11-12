package tax_report_line

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles TaxReportLines validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new TaxReportLines validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateTaxReportLinesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate TaxReportDefinitionId
	
	
	if err := v.validateTaxReportDefinitionIdExists(ctx, tx, req.TaxReportDefinitionId); err != nil {
		return err
	}
	
	
	// Validate LineCode
	
	if err := v.validateLineCode(req.LineCode); err != nil {
		return err
	}
	
	
	
	// Validate LineName
	
	if err := v.validateLineName(req.LineName); err != nil {
		return err
	}
	
	
	
	// Validate Sequence
	
	
	
	// Validate ParentLineId
	
	
	if err := v.validateParentLineIdExists(ctx, tx, req.ParentLineId); err != nil {
		return err
	}
	
	
	// Validate FormulaType
	
	
	
	// Validate Formula
	
	
	
	// Validate TaxGroupIds
	
	
	
	// Validate AccountIds
	
	
	
	// Validate TaxIds
	
	
	
	// Validate IsSubtotal
	
	
	
	// Validate IsTotal
	
	
	
	// Validate Notes
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateTaxReportLinesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate TaxReportDefinitionId if provided
	
	
	if req.TaxReportDefinitionId != nil {
		if err := v.validateTaxReportDefinitionIdExists(ctx, tx, *req.TaxReportDefinitionId); err != nil {
			return err
		}
	}
	
	
	// Validate LineCode if provided
	
	if req.LineCode != nil {
		if err := v.validateLineCode(*req.LineCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate LineName if provided
	
	if req.LineName != nil {
		if err := v.validateLineName(*req.LineName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Sequence if provided
	
	
	
	// Validate ParentLineId if provided
	
	
	if req.ParentLineId != nil {
		if err := v.validateParentLineIdExists(ctx, tx, *req.ParentLineId); err != nil {
			return err
		}
	}
	
	
	// Validate FormulaType if provided
	
	
	
	// Validate Formula if provided
	
	
	
	// Validate TaxGroupIds if provided
	
	
	
	// Validate AccountIds if provided
	
	
	
	// Validate TaxIds if provided
	
	
	
	// Validate IsSubtotal if provided
	
	
	
	// Validate IsTotal if provided
	
	
	
	// Validate Notes if provided
	
	
	

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





// validateTaxReportDefinitionIdExists validates that tax_report_definition_id exists
func (v *Validator) validateTaxReportDefinitionIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for tax_report_definition
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_report_definition WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check tax_report_definition existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("tax_report_definition with id %s does not exist", id)
	}
	return nil
}



// validateLineCode validates line_code field
func (v *Validator) validateLineCode(value string) error {
	
	// Add custom validation for line_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("line_code cannot be empty")
	}
	
	return nil
}





// validateLineName validates line_name field
func (v *Validator) validateLineName(value string) error {
	
	// Add custom validation for line_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("line_name cannot be empty")
	}
	
	return nil
}











// validateParentLineIdExists validates that parent_line_id exists
func (v *Validator) validateParentLineIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for parent_line
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM parent_line WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check parent_line existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("parent_line with id %s does not exist", id)
	}
	return nil
}



































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateTaxReportLinesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateTaxReportLinesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *TaxReportLines, req *UpdateTaxReportLinesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *TaxReportLines) error {
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

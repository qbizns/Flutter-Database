package budget

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles Budgets validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Budgets validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateBudgetsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate BudgetCode
	
	if err := v.validateBudgetCode(req.BudgetCode); err != nil {
		return err
	}
	
	
	
	// Validate BudgetName
	
	if err := v.validateBudgetName(req.BudgetName); err != nil {
		return err
	}
	
	
	
	// Validate FiscalYearId
	
	
	if err := v.validateFiscalYearIdExists(ctx, tx, req.FiscalYearId); err != nil {
		return err
	}
	
	
	// Validate StartDate
	
	
	
	// Validate EndDate
	
	
	
	// Validate BudgetType
	
	
	
	// Validate BudgetType
	
	
	
	// Validate Status
	
	
	
	// Validate Notes
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateBudgetsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate BudgetCode if provided
	
	if req.BudgetCode != nil {
		if err := v.validateBudgetCode(*req.BudgetCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate BudgetName if provided
	
	if req.BudgetName != nil {
		if err := v.validateBudgetName(*req.BudgetName); err != nil {
			return err
		}
	}
	
	
	
	// Validate FiscalYearId if provided
	
	
	if req.FiscalYearId != nil {
		if err := v.validateFiscalYearIdExists(ctx, tx, *req.FiscalYearId); err != nil {
			return err
		}
	}
	
	
	// Validate StartDate if provided
	
	
	
	// Validate EndDate if provided
	
	
	
	// Validate BudgetType if provided
	
	
	
	// Validate BudgetType if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate Notes if provided
	
	
	
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



// validateBudgetCode validates budget_code field
func (v *Validator) validateBudgetCode(value string) error {
	
	// Add custom validation for budget_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("budget_code cannot be empty")
	}
	
	return nil
}





// validateBudgetName validates budget_name field
func (v *Validator) validateBudgetName(value string) error {
	
	// Add custom validation for budget_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("budget_name cannot be empty")
	}
	
	return nil
}







// validateFiscalYearIdExists validates that fiscal_year_id exists
func (v *Validator) validateFiscalYearIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for fiscal_year
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM fiscal_year WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check fiscal_year existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("fiscal_year with id %s does not exist", id)
	}
	return nil
}



































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateBudgetsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateBudgetsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Budgets, req *UpdateBudgetsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Budgets) error {
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

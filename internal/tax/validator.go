package tax

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto"
)

// Validator handles Taxes validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Taxes validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateTaxesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate TaxGroupId
	
	
	if err := v.validateTaxGroupIdExists(ctx, tx, req.TaxGroupId); err != nil {
		return err
	}
	
	
	// Validate TaxCode
	
	if err := v.validateTaxCode(req.TaxCode); err != nil {
		return err
	}
	
	
	
	// Validate TaxName
	
	if err := v.validateTaxName(req.TaxName); err != nil {
		return err
	}
	
	
	
	// Validate TaxRate
	
	
	
	// Validate TaxScope
	
	if err := v.validateTaxScope(req.TaxScope); err != nil {
		return err
	}
	
	
	
	// Validate IsPriceInclusive
	
	
	
	// Validate TaxAccountId
	
	
	if err := v.validateTaxAccountIdExists(ctx, tx, req.TaxAccountId); err != nil {
		return err
	}
	
	
	// Validate TaxRefundAccountId
	
	
	if err := v.validateTaxRefundAccountIdExists(ctx, tx, req.TaxRefundAccountId); err != nil {
		return err
	}
	
	
	// Validate IsActive
	
	
	
	// Validate Description
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateTaxesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate TaxGroupId if provided
	
	
	if req.TaxGroupId != nil {
		if err := v.validateTaxGroupIdExists(ctx, tx, *req.TaxGroupId); err != nil {
			return err
		}
	}
	
	
	// Validate TaxCode if provided
	
	if req.TaxCode != nil {
		if err := v.validateTaxCode(*req.TaxCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate TaxName if provided
	
	if req.TaxName != nil {
		if err := v.validateTaxName(*req.TaxName); err != nil {
			return err
		}
	}
	
	
	
	// Validate TaxRate if provided
	
	
	
	// Validate TaxScope if provided
	
	if req.TaxScope != nil {
		if err := v.validateTaxScope(*req.TaxScope); err != nil {
			return err
		}
	}
	
	
	
	// Validate IsPriceInclusive if provided
	
	
	
	// Validate TaxAccountId if provided
	
	
	if req.TaxAccountId != nil {
		if err := v.validateTaxAccountIdExists(ctx, tx, *req.TaxAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate TaxRefundAccountId if provided
	
	
	if req.TaxRefundAccountId != nil {
		if err := v.validateTaxRefundAccountIdExists(ctx, tx, *req.TaxRefundAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate IsActive if provided
	
	
	
	// Validate Description if provided
	
	
	
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





// validateTaxGroupIdExists validates that tax_group_id exists
func (v *Validator) validateTaxGroupIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for tax_group
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_group WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check tax_group existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("tax_group with id %s does not exist", id)
	}
	return nil
}



// validateTaxCode validates tax_code field
func (v *Validator) validateTaxCode(value string) error {
	
	// Add custom validation for tax_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("tax_code cannot be empty")
	}
	
	return nil
}





// validateTaxName validates tax_name field
func (v *Validator) validateTaxName(value string) error {
	
	// Add custom validation for tax_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("tax_name cannot be empty")
	}
	
	return nil
}









// validateTaxScope validates tax_scope field
func (v *Validator) validateTaxScope(value string) error {
	
	// Add custom validation for tax_scope
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("tax_scope cannot be empty")
	}
	
	return nil
}











// validateTaxAccountIdExists validates that tax_account_id exists
func (v *Validator) validateTaxAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for tax_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check tax_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("tax_account with id %s does not exist", id)
	}
	return nil
}





// validateTaxRefundAccountIdExists validates that tax_refund_account_id exists
func (v *Validator) validateTaxRefundAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for tax_refund_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax_refund_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check tax_refund_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("tax_refund_account with id %s does not exist", id)
	}
	return nil
}



















// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateTaxesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateTaxesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Taxes, req *dto.UpdateTaxesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Taxes) error {
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

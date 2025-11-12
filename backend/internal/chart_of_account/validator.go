package chart_of_account

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles ChartOfAccounts validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ChartOfAccounts validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateChartOfAccountsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate AccountCode
	
	if err := v.validateAccountCode(req.AccountCode); err != nil {
		return err
	}
	
	
	
	// Validate AccountNumber
	
	if err := v.validateAccountNumber(req.AccountNumber); err != nil {
		return err
	}
	
	
	
	// Validate AccountName
	
	if err := v.validateAccountName(req.AccountName); err != nil {
		return err
	}
	
	
	
	// Validate AccountTypeId
	
	
	if err := v.validateAccountTypeIdExists(ctx, tx, req.AccountTypeId); err != nil {
		return err
	}
	
	
	// Validate AccountSubtypeId
	
	
	if err := v.validateAccountSubtypeIdExists(ctx, tx, req.AccountSubtypeId); err != nil {
		return err
	}
	
	
	// Validate ParentAccountId
	
	
	if err := v.validateParentAccountIdExists(ctx, tx, req.ParentAccountId); err != nil {
		return err
	}
	
	
	// Validate AccountLevel
	
	
	
	// Validate AccountPath
	
	
	
	// Validate IsActive
	
	
	
	// Validate IsSystemAccount
	
	
	
	// Validate IsHeaderAccount
	
	
	
	// Validate IsBankAccount
	
	
	
	// Validate IsReconcilable
	
	
	
	// Validate DefaultTaxCode
	
	
	
	// Validate CurrencyCode
	
	
	
	// Validate OpeningBalance
	
	
	
	// Validate OpeningBalanceDate
	
	
	
	// Validate CurrentDebitBalance
	
	
	
	// Validate CurrentCreditBalance
	
	
	
	// Validate CurrentBalance
	
	
	
	// Validate LastBalanceUpdate
	
	
	
	// Validate Description
	
	
	
	// Validate Notes
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateChartOfAccountsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate AccountCode if provided
	
	if req.AccountCode != nil {
		if err := v.validateAccountCode(*req.AccountCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate AccountNumber if provided
	
	if req.AccountNumber != nil {
		if err := v.validateAccountNumber(*req.AccountNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate AccountName if provided
	
	if req.AccountName != nil {
		if err := v.validateAccountName(*req.AccountName); err != nil {
			return err
		}
	}
	
	
	
	// Validate AccountTypeId if provided
	
	
	if req.AccountTypeId != nil {
		if err := v.validateAccountTypeIdExists(ctx, tx, *req.AccountTypeId); err != nil {
			return err
		}
	}
	
	
	// Validate AccountSubtypeId if provided
	
	
	if req.AccountSubtypeId != nil {
		if err := v.validateAccountSubtypeIdExists(ctx, tx, *req.AccountSubtypeId); err != nil {
			return err
		}
	}
	
	
	// Validate ParentAccountId if provided
	
	
	if req.ParentAccountId != nil {
		if err := v.validateParentAccountIdExists(ctx, tx, *req.ParentAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate AccountLevel if provided
	
	
	
	// Validate AccountPath if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate IsSystemAccount if provided
	
	
	
	// Validate IsHeaderAccount if provided
	
	
	
	// Validate IsBankAccount if provided
	
	
	
	// Validate IsReconcilable if provided
	
	
	
	// Validate DefaultTaxCode if provided
	
	
	
	// Validate CurrencyCode if provided
	
	
	
	// Validate OpeningBalance if provided
	
	
	
	// Validate OpeningBalanceDate if provided
	
	
	
	// Validate CurrentDebitBalance if provided
	
	
	
	// Validate CurrentCreditBalance if provided
	
	
	
	// Validate CurrentBalance if provided
	
	
	
	// Validate LastBalanceUpdate if provided
	
	
	
	// Validate Description if provided
	
	
	
	// Validate Notes if provided
	
	
	
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



// validateAccountCode validates account_code field
func (v *Validator) validateAccountCode(value string) error {
	
	// Add custom validation for account_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("account_code cannot be empty")
	}
	
	return nil
}





// validateAccountNumber validates account_number field
func (v *Validator) validateAccountNumber(value string) error {
	
	// Add custom validation for account_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("account_number cannot be empty")
	}
	
	return nil
}





// validateAccountName validates account_name field
func (v *Validator) validateAccountName(value string) error {
	
	// Add custom validation for account_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("account_name cannot be empty")
	}
	
	return nil
}







// validateAccountTypeIdExists validates that account_type_id exists
func (v *Validator) validateAccountTypeIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for account_type
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM account_type WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check account_type existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("account_type with id %s does not exist", id)
	}
	return nil
}





// validateAccountSubtypeIdExists validates that account_subtype_id exists
func (v *Validator) validateAccountSubtypeIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for account_subtype
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM account_subtype WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check account_subtype existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("account_subtype with id %s does not exist", id)
	}
	return nil
}





// validateParentAccountIdExists validates that parent_account_id exists
func (v *Validator) validateParentAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for parent_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM parent_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check parent_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("parent_account with id %s does not exist", id)
	}
	return nil
}



















































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateChartOfAccountsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateChartOfAccountsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ChartOfAccounts, req *UpdateChartOfAccountsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ChartOfAccounts) error {
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

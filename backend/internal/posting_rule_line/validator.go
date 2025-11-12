package posting_rule_line

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles PostingRuleLines validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new PostingRuleLines validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreatePostingRuleLinesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate PostingRuleId
	
	
	if err := v.validatePostingRuleIdExists(ctx, tx, req.PostingRuleId); err != nil {
		return err
	}
	
	
	// Validate LineNo
	
	
	
	// Validate Side
	
	if err := v.validateSide(req.Side); err != nil {
		return err
	}
	
	
	
	// Validate ConceptKey
	
	
	
	// Validate AccountSource
	
	if err := v.validateAccountSource(req.AccountSource); err != nil {
		return err
	}
	
	
	
	// Validate FixedAccountId
	
	
	if err := v.validateFixedAccountIdExists(ctx, tx, req.FixedAccountId); err != nil {
		return err
	}
	
	
	// Validate AccountFieldPath
	
	
	
	// Validate AccountExpression
	
	
	
	// Validate AmountSource
	
	if err := v.validateAmountSource(req.AmountSource); err != nil {
		return err
	}
	
	
	
	// Validate AmountFieldPath
	
	
	
	// Validate AmountExpression
	
	
	
	// Validate MappingContext
	
	
	
	// Validate DescriptionTemplate
	
	
	
	// Validate IsActive
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate (accountSource
	
	if err := v.validate(accountSource(req.(accountSource); err != nil {
		return err
	}
	
	
	
	// Validate (accountSource
	
	if err := v.validate(accountSource(req.(accountSource); err != nil {
		return err
	}
	
	
	
	// Validate (accountSource
	
	if err := v.validate(accountSource(req.(accountSource); err != nil {
		return err
	}
	
	
	
	// Validate (accountSource
	
	if err := v.validate(accountSource(req.(accountSource); err != nil {
		return err
	}
	
	
	
	// Validate (amountSource
	
	if err := v.validate(amountSource(req.(amountSource); err != nil {
		return err
	}
	
	
	
	// Validate (amountSource
	
	if err := v.validate(amountSource(req.(amountSource); err != nil {
		return err
	}
	
	
	
	// Validate (amountSource
	
	if err := v.validate(amountSource(req.(amountSource); err != nil {
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdatePostingRuleLinesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate PostingRuleId if provided
	
	
	if req.PostingRuleId != nil {
		if err := v.validatePostingRuleIdExists(ctx, tx, *req.PostingRuleId); err != nil {
			return err
		}
	}
	
	
	// Validate LineNo if provided
	
	
	
	// Validate Side if provided
	
	if req.Side != nil {
		if err := v.validateSide(*req.Side); err != nil {
			return err
		}
	}
	
	
	
	// Validate ConceptKey if provided
	
	
	
	// Validate AccountSource if provided
	
	if req.AccountSource != nil {
		if err := v.validateAccountSource(*req.AccountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate FixedAccountId if provided
	
	
	if req.FixedAccountId != nil {
		if err := v.validateFixedAccountIdExists(ctx, tx, *req.FixedAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate AccountFieldPath if provided
	
	
	
	// Validate AccountExpression if provided
	
	
	
	// Validate AmountSource if provided
	
	if req.AmountSource != nil {
		if err := v.validateAmountSource(*req.AmountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate AmountFieldPath if provided
	
	
	
	// Validate AmountExpression if provided
	
	
	
	// Validate MappingContext if provided
	
	
	
	// Validate DescriptionTemplate if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate (accountSource if provided
	
	if req.(accountSource != nil {
		if err := v.validate(accountSource(*req.(accountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate (accountSource if provided
	
	if req.(accountSource != nil {
		if err := v.validate(accountSource(*req.(accountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate (accountSource if provided
	
	if req.(accountSource != nil {
		if err := v.validate(accountSource(*req.(accountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate (accountSource if provided
	
	if req.(accountSource != nil {
		if err := v.validate(accountSource(*req.(accountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate (amountSource if provided
	
	if req.(amountSource != nil {
		if err := v.validate(amountSource(*req.(amountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate (amountSource if provided
	
	if req.(amountSource != nil {
		if err := v.validate(amountSource(*req.(amountSource); err != nil {
			return err
		}
	}
	
	
	
	// Validate (amountSource if provided
	
	if req.(amountSource != nil {
		if err := v.validate(amountSource(*req.(amountSource); err != nil {
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





// validatePostingRuleIdExists validates that posting_rule_id exists
func (v *Validator) validatePostingRuleIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for posting_rule
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM posting_rule WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check posting_rule existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("posting_rule with id %s does not exist", id)
	}
	return nil
}







// validateSide validates side field
func (v *Validator) validateSide(value string) error {
	
	// Add custom validation for side
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("side cannot be empty")
	}
	
	return nil
}









// validateAccountSource validates account_source field
func (v *Validator) validateAccountSource(value string) error {
	
	// Add custom validation for account_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("account_source cannot be empty")
	}
	
	return nil
}







// validateFixedAccountIdExists validates that fixed_account_id exists
func (v *Validator) validateFixedAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for fixed_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM fixed_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check fixed_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("fixed_account with id %s does not exist", id)
	}
	return nil
}











// validateAmountSource validates amount_source field
func (v *Validator) validateAmountSource(value string) error {
	
	// Add custom validation for amount_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("amount_source cannot be empty")
	}
	
	return nil
}

































// validate(accountSource validates (account_source field
func (v *Validator) validate(accountSource(value string) error {
	
	// Add custom validation for (account_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(account_source cannot be empty")
	}
	
	return nil
}





// validate(accountSource validates (account_source field
func (v *Validator) validate(accountSource(value string) error {
	
	// Add custom validation for (account_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(account_source cannot be empty")
	}
	
	return nil
}





// validate(accountSource validates (account_source field
func (v *Validator) validate(accountSource(value string) error {
	
	// Add custom validation for (account_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(account_source cannot be empty")
	}
	
	return nil
}





// validate(accountSource validates (account_source field
func (v *Validator) validate(accountSource(value string) error {
	
	// Add custom validation for (account_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(account_source cannot be empty")
	}
	
	return nil
}





// validate(amountSource validates (amount_source field
func (v *Validator) validate(amountSource(value string) error {
	
	// Add custom validation for (amount_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(amount_source cannot be empty")
	}
	
	return nil
}





// validate(amountSource validates (amount_source field
func (v *Validator) validate(amountSource(value string) error {
	
	// Add custom validation for (amount_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(amount_source cannot be empty")
	}
	
	return nil
}





// validate(amountSource validates (amount_source field
func (v *Validator) validate(amountSource(value string) error {
	
	// Add custom validation for (amount_source
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(amount_source cannot be empty")
	}
	
	return nil
}





// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreatePostingRuleLinesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreatePostingRuleLinesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *PostingRuleLines, req *UpdatePostingRuleLinesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *PostingRuleLines) error {
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

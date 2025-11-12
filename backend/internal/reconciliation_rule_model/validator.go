package reconciliation_rule_model

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

// Validator handles ReconciliationRuleModels validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ReconciliationRuleModels validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateReconciliationRuleModelsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate RuleName
	
	if err := v.validateRuleName(req.RuleName); err != nil {
		return err
	}
	
	
	
	// Validate RuleCode
	
	
	
	// Validate Sequence
	
	
	
	// Validate AmountMin
	
	
	
	// Validate AmountMax
	
	
	
	// Validate DescriptionPattern
	
	
	
	// Validate CounterpartyPattern
	
	
	
	// Validate ReferencePattern
	
	
	
	// Validate JournalId
	
	
	if err := v.validateJournalIdExists(ctx, tx, req.JournalId); err != nil {
		return err
	}
	
	
	// Validate AccountId
	
	
	if err := v.validateAccountIdExists(ctx, tx, req.AccountId); err != nil {
		return err
	}
	
	
	// Validate AnalyticAccountId
	
	
	if err := v.validateAnalyticAccountIdExists(ctx, tx, req.AnalyticAccountId); err != nil {
		return err
	}
	
	
	// Validate TaxId
	
	
	if err := v.validateTaxIdExists(ctx, tx, req.TaxId); err != nil {
		return err
	}
	
	
	// Validate IsActive
	
	
	
	// Validate AutoApply
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateReconciliationRuleModelsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate RuleName if provided
	
	if req.RuleName != nil {
		if err := v.validateRuleName(*req.RuleName); err != nil {
			return err
		}
	}
	
	
	
	// Validate RuleCode if provided
	
	
	
	// Validate Sequence if provided
	
	
	
	// Validate AmountMin if provided
	
	
	
	// Validate AmountMax if provided
	
	
	
	// Validate DescriptionPattern if provided
	
	
	
	// Validate CounterpartyPattern if provided
	
	
	
	// Validate ReferencePattern if provided
	
	
	
	// Validate JournalId if provided
	
	
	if req.JournalId != nil {
		if err := v.validateJournalIdExists(ctx, tx, *req.JournalId); err != nil {
			return err
		}
	}
	
	
	// Validate AccountId if provided
	
	
	if req.AccountId != nil {
		if err := v.validateAccountIdExists(ctx, tx, *req.AccountId); err != nil {
			return err
		}
	}
	
	
	// Validate AnalyticAccountId if provided
	
	
	if req.AnalyticAccountId != nil {
		if err := v.validateAnalyticAccountIdExists(ctx, tx, *req.AnalyticAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate TaxId if provided
	
	
	if req.TaxId != nil {
		if err := v.validateTaxIdExists(ctx, tx, *req.TaxId); err != nil {
			return err
		}
	}
	
	
	// Validate IsActive if provided
	
	
	
	// Validate AutoApply if provided
	
	
	
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



// validateRuleName validates rule_name field
func (v *Validator) validateRuleName(value string) error {
	
	// Add custom validation for rule_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("rule_name cannot be empty")
	}
	
	return nil
}



































// validateJournalIdExists validates that journal_id exists
func (v *Validator) validateJournalIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for journal
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM journal WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check journal existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("journal with id %s does not exist", id)
	}
	return nil
}





// validateAccountIdExists validates that account_id exists
func (v *Validator) validateAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("account with id %s does not exist", id)
	}
	return nil
}





// validateAnalyticAccountIdExists validates that analytic_account_id exists
func (v *Validator) validateAnalyticAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for analytic_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM analytic_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check analytic_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("analytic_account with id %s does not exist", id)
	}
	return nil
}





// validateTaxIdExists validates that tax_id exists
func (v *Validator) validateTaxIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for tax
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tax WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check tax existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("tax with id %s does not exist", id)
	}
	return nil
}























// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateReconciliationRuleModelsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateReconciliationRuleModelsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ReconciliationRuleModels, req *dto.UpdateReconciliationRuleModelsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ReconciliationRuleModels) error {
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

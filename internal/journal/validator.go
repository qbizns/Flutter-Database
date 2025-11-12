package journal

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

// Validator handles Journals validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Journals validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateJournalsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate JournalCode
	
	if err := v.validateJournalCode(req.JournalCode); err != nil {
		return err
	}
	
	
	
	// Validate JournalName
	
	if err := v.validateJournalName(req.JournalName); err != nil {
		return err
	}
	
	
	
	// Validate JournalType
	
	if err := v.validateJournalType(req.JournalType); err != nil {
		return err
	}
	
	
	
	// Validate 'sale',
	
	
	
	// Validate BankAccountId
	
	
	if err := v.validateBankAccountIdExists(ctx, tx, req.BankAccountId); err != nil {
		return err
	}
	
	
	// Validate DefaultDebitAccountId
	
	
	if err := v.validateDefaultDebitAccountIdExists(ctx, tx, req.DefaultDebitAccountId); err != nil {
		return err
	}
	
	
	// Validate DefaultCreditAccountId
	
	
	if err := v.validateDefaultCreditAccountIdExists(ctx, tx, req.DefaultCreditAccountId); err != nil {
		return err
	}
	
	
	// Validate SequencePrefix
	
	
	
	// Validate SequenceNumber
	
	
	
	// Validate IsActive
	
	
	
	// Validate Notes
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate (journalType
	
	if err := v.validate(journalType(req.(journalType); err != nil {
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateJournalsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate JournalCode if provided
	
	if req.JournalCode != nil {
		if err := v.validateJournalCode(*req.JournalCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate JournalName if provided
	
	if req.JournalName != nil {
		if err := v.validateJournalName(*req.JournalName); err != nil {
			return err
		}
	}
	
	
	
	// Validate JournalType if provided
	
	if req.JournalType != nil {
		if err := v.validateJournalType(*req.JournalType); err != nil {
			return err
		}
	}
	
	
	
	// Validate 'sale', if provided
	
	
	
	// Validate BankAccountId if provided
	
	
	if req.BankAccountId != nil {
		if err := v.validateBankAccountIdExists(ctx, tx, *req.BankAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate DefaultDebitAccountId if provided
	
	
	if req.DefaultDebitAccountId != nil {
		if err := v.validateDefaultDebitAccountIdExists(ctx, tx, *req.DefaultDebitAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate DefaultCreditAccountId if provided
	
	
	if req.DefaultCreditAccountId != nil {
		if err := v.validateDefaultCreditAccountIdExists(ctx, tx, *req.DefaultCreditAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate SequencePrefix if provided
	
	
	
	// Validate SequenceNumber if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate (journalType if provided
	
	if req.(journalType != nil {
		if err := v.validate(journalType(*req.(journalType); err != nil {
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



// validateJournalCode validates journal_code field
func (v *Validator) validateJournalCode(value string) error {
	
	// Add custom validation for journal_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("journal_code cannot be empty")
	}
	
	return nil
}





// validateJournalName validates journal_name field
func (v *Validator) validateJournalName(value string) error {
	
	// Add custom validation for journal_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("journal_name cannot be empty")
	}
	
	return nil
}





// validateJournalType validates journal_type field
func (v *Validator) validateJournalType(value string) error {
	
	// Add custom validation for journal_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("journal_type cannot be empty")
	}
	
	return nil
}











// validateBankAccountIdExists validates that bank_account_id exists
func (v *Validator) validateBankAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for bank_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM bank_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check bank_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bank_account with id %s does not exist", id)
	}
	return nil
}





// validateDefaultDebitAccountIdExists validates that default_debit_account_id exists
func (v *Validator) validateDefaultDebitAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for default_debit_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM default_debit_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check default_debit_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("default_debit_account with id %s does not exist", id)
	}
	return nil
}





// validateDefaultCreditAccountIdExists validates that default_credit_account_id exists
func (v *Validator) validateDefaultCreditAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for default_credit_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM default_credit_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check default_credit_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("default_credit_account with id %s does not exist", id)
	}
	return nil
}



























// validate(journalType validates (journal_type field
func (v *Validator) validate(journalType(value string) error {
	
	// Add custom validation for (journal_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(journal_type cannot be empty")
	}
	
	return nil
}





// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateJournalsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateJournalsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Journals, req *dto.UpdateJournalsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Journals) error {
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

package general_ledger

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles GeneralLedger validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new GeneralLedger validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateGeneralLedgerRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate JournalEntryId
	
	
	if err := v.validateJournalEntryIdExists(ctx, tx, req.JournalEntryId); err != nil {
		return err
	}
	
	
	// Validate JournalEntryLineId
	
	
	if err := v.validateJournalEntryLineIdExists(ctx, tx, req.JournalEntryLineId); err != nil {
		return err
	}
	
	
	// Validate AccountId
	
	
	if err := v.validateAccountIdExists(ctx, tx, req.AccountId); err != nil {
		return err
	}
	
	
	// Validate TransactionDate
	
	
	
	// Validate PostingDate
	
	
	
	// Validate AccountingPeriodId
	
	
	if err := v.validateAccountingPeriodIdExists(ctx, tx, req.AccountingPeriodId); err != nil {
		return err
	}
	
	
	// Validate FiscalYearId
	
	
	if err := v.validateFiscalYearIdExists(ctx, tx, req.FiscalYearId); err != nil {
		return err
	}
	
	
	// Validate DebitAmount
	
	
	
	// Validate CreditAmount
	
	
	
	// Validate RunningDebitBalance
	
	
	
	// Validate RunningCreditBalance
	
	
	
	// Validate RunningBalance
	
	
	
	// Validate SourceModule
	
	
	
	// Validate SourceDocumentType
	
	
	
	// Validate SourceDocumentId
	
	
	if err := v.validateSourceDocumentIdExists(ctx, tx, req.SourceDocumentId); err != nil {
		return err
	}
	
	
	// Validate ReferenceNumber
	
	
	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate Department
	
	
	
	// Validate ProjectCode
	
	
	
	// Validate CostCenter
	
	
	
	// Validate Description
	
	
	
	// Validate IsReversed
	
	
	
	// Validate ReversalGlId
	
	
	if err := v.validateReversalGlIdExists(ctx, tx, req.ReversalGlId); err != nil {
		return err
	}
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate (debitAmount
	
	
	
	// Validate (creditAmount
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateGeneralLedgerRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate JournalEntryId if provided
	
	
	if req.JournalEntryId != nil {
		if err := v.validateJournalEntryIdExists(ctx, tx, *req.JournalEntryId); err != nil {
			return err
		}
	}
	
	
	// Validate JournalEntryLineId if provided
	
	
	if req.JournalEntryLineId != nil {
		if err := v.validateJournalEntryLineIdExists(ctx, tx, *req.JournalEntryLineId); err != nil {
			return err
		}
	}
	
	
	// Validate AccountId if provided
	
	
	if req.AccountId != nil {
		if err := v.validateAccountIdExists(ctx, tx, *req.AccountId); err != nil {
			return err
		}
	}
	
	
	// Validate TransactionDate if provided
	
	
	
	// Validate PostingDate if provided
	
	
	
	// Validate AccountingPeriodId if provided
	
	
	if req.AccountingPeriodId != nil {
		if err := v.validateAccountingPeriodIdExists(ctx, tx, *req.AccountingPeriodId); err != nil {
			return err
		}
	}
	
	
	// Validate FiscalYearId if provided
	
	
	if req.FiscalYearId != nil {
		if err := v.validateFiscalYearIdExists(ctx, tx, *req.FiscalYearId); err != nil {
			return err
		}
	}
	
	
	// Validate DebitAmount if provided
	
	
	
	// Validate CreditAmount if provided
	
	
	
	// Validate RunningDebitBalance if provided
	
	
	
	// Validate RunningCreditBalance if provided
	
	
	
	// Validate RunningBalance if provided
	
	
	
	// Validate SourceModule if provided
	
	
	
	// Validate SourceDocumentType if provided
	
	
	
	// Validate SourceDocumentId if provided
	
	
	if req.SourceDocumentId != nil {
		if err := v.validateSourceDocumentIdExists(ctx, tx, *req.SourceDocumentId); err != nil {
			return err
		}
	}
	
	
	// Validate ReferenceNumber if provided
	
	
	
	// Validate LocationId if provided
	
	
	if req.LocationId != nil {
		if err := v.validateLocationIdExists(ctx, tx, *req.LocationId); err != nil {
			return err
		}
	}
	
	
	// Validate Department if provided
	
	
	
	// Validate ProjectCode if provided
	
	
	
	// Validate CostCenter if provided
	
	
	
	// Validate Description if provided
	
	
	
	// Validate IsReversed if provided
	
	
	
	// Validate ReversalGlId if provided
	
	
	if req.ReversalGlId != nil {
		if err := v.validateReversalGlIdExists(ctx, tx, *req.ReversalGlId); err != nil {
			return err
		}
	}
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate (debitAmount if provided
	
	
	
	// Validate (creditAmount if provided
	
	
	

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





// validateJournalEntryLineIdExists validates that journal_entry_line_id exists
func (v *Validator) validateJournalEntryLineIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for journal_entry_line
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM journal_entry_line WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check journal_entry_line existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("journal_entry_line with id %s does not exist", id)
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













// validateAccountingPeriodIdExists validates that accounting_period_id exists
func (v *Validator) validateAccountingPeriodIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for accounting_period
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounting_period WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check accounting_period existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("accounting_period with id %s does not exist", id)
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

































// validateSourceDocumentIdExists validates that source_document_id exists
func (v *Validator) validateSourceDocumentIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for source_document
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM source_document WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check source_document existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("source_document with id %s does not exist", id)
	}
	return nil
}









// validateLocationIdExists validates that location_id exists
func (v *Validator) validateLocationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for location
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM location WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check location existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("location with id %s does not exist", id)
	}
	return nil
}

























// validateReversalGlIdExists validates that reversal_gl_id exists
func (v *Validator) validateReversalGlIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for reversal_gl
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM reversal_gl WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check reversal_gl existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("reversal_gl with id %s does not exist", id)
	}
	return nil
}



















// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateGeneralLedgerRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateGeneralLedgerRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *GeneralLedger, req *UpdateGeneralLedgerRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *GeneralLedger) error {
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

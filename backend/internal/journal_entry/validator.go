package journal_entry

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

// Validator handles JournalEntries validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new JournalEntries validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateJournalEntriesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate EntryNumber
	
	if err := v.validateEntryNumber(req.EntryNumber); err != nil {
		return err
	}
	
	
	
	// Validate EntryTypeId
	
	
	if err := v.validateEntryTypeIdExists(ctx, tx, req.EntryTypeId); err != nil {
		return err
	}
	
	
	// Validate EntryDate
	
	
	
	// Validate PostingDate
	
	
	
	// Validate AccountingPeriodId
	
	
	if err := v.validateAccountingPeriodIdExists(ctx, tx, req.AccountingPeriodId); err != nil {
		return err
	}
	
	
	// Validate FiscalYearId
	
	
	if err := v.validateFiscalYearIdExists(ctx, tx, req.FiscalYearId); err != nil {
		return err
	}
	
	
	// Validate Status
	
	
	
	// Validate IsPosted
	
	
	
	// Validate IsReversed
	
	
	
	// Validate ReversalEntryId
	
	
	if err := v.validateReversalEntryIdExists(ctx, tx, req.ReversalEntryId); err != nil {
		return err
	}
	
	
	// Validate SourceModule
	
	
	
	// Validate SourceDocumentType
	
	
	
	// Validate SourceDocumentId
	
	
	if err := v.validateSourceDocumentIdExists(ctx, tx, req.SourceDocumentId); err != nil {
		return err
	}
	
	
	// Validate ReferenceNumber
	
	
	
	// Validate TotalDebit
	
	
	
	// Validate TotalCredit
	
	
	
	// Validate Description
	
	if err := v.validateDescription(req.Description); err != nil {
		return err
	}
	
	
	
	// Validate Notes
	
	
	
	// Validate RequiresApproval
	
	
	
	// Validate ApprovedBy
	
	
	
	// Validate ApprovedAt
	
	
	
	// Validate PostedBy
	
	
	
	// Validate PostedAt
	
	
	
	// Validate Attachments
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate (isPosted
	
	
	
	// Validate (ABS(totalDebit
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateJournalEntriesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate EntryNumber if provided
	
	if req.EntryNumber != nil {
		if err := v.validateEntryNumber(*req.EntryNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate EntryTypeId if provided
	
	
	if req.EntryTypeId != nil {
		if err := v.validateEntryTypeIdExists(ctx, tx, *req.EntryTypeId); err != nil {
			return err
		}
	}
	
	
	// Validate EntryDate if provided
	
	
	
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
	
	
	// Validate Status if provided
	
	
	
	// Validate IsPosted if provided
	
	
	
	// Validate IsReversed if provided
	
	
	
	// Validate ReversalEntryId if provided
	
	
	if req.ReversalEntryId != nil {
		if err := v.validateReversalEntryIdExists(ctx, tx, *req.ReversalEntryId); err != nil {
			return err
		}
	}
	
	
	// Validate SourceModule if provided
	
	
	
	// Validate SourceDocumentType if provided
	
	
	
	// Validate SourceDocumentId if provided
	
	
	if req.SourceDocumentId != nil {
		if err := v.validateSourceDocumentIdExists(ctx, tx, *req.SourceDocumentId); err != nil {
			return err
		}
	}
	
	
	// Validate ReferenceNumber if provided
	
	
	
	// Validate TotalDebit if provided
	
	
	
	// Validate TotalCredit if provided
	
	
	
	// Validate Description if provided
	
	if req.Description != nil {
		if err := v.validateDescription(*req.Description); err != nil {
			return err
		}
	}
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate RequiresApproval if provided
	
	
	
	// Validate ApprovedBy if provided
	
	
	
	// Validate ApprovedAt if provided
	
	
	
	// Validate PostedBy if provided
	
	
	
	// Validate PostedAt if provided
	
	
	
	// Validate Attachments if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate (isPosted if provided
	
	
	
	// Validate (ABS(totalDebit if provided
	
	
	

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



// validateEntryNumber validates entry_number field
func (v *Validator) validateEntryNumber(value string) error {
	
	// Add custom validation for entry_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("entry_number cannot be empty")
	}
	
	return nil
}







// validateEntryTypeIdExists validates that entry_type_id exists
func (v *Validator) validateEntryTypeIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for entry_type
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM entry_type WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check entry_type existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("entry_type with id %s does not exist", id)
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

















// validateReversalEntryIdExists validates that reversal_entry_id exists
func (v *Validator) validateReversalEntryIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for reversal_entry
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM reversal_entry WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check reversal_entry existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("reversal_entry with id %s does not exist", id)
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















// validateDescription validates description field
func (v *Validator) validateDescription(value string) error {
	
	// Add custom validation for description
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("description cannot be empty")
	}
	
	return nil
}





















































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateJournalEntriesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateJournalEntriesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *JournalEntries, req *dto.UpdateJournalEntriesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *JournalEntries) error {
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

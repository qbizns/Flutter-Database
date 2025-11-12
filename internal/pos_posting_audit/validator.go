package pos_posting_audit

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

// Validator handles PosPostingAudit validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new PosPostingAudit validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreatePosPostingAuditRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate SourceTable
	
	if err := v.validateSourceTable(req.SourceTable); err != nil {
		return err
	}
	
	
	
	// Validate SourceId
	
	
	if err := v.validateSourceIdExists(ctx, tx, req.SourceId); err != nil {
		return err
	}
	
	
	// Validate SourceReference
	
	
	
	// Validate PostingStatus
	
	if err := v.validatePostingStatus(req.PostingStatus); err != nil {
		return err
	}
	
	
	
	// Validate 'pending',
	
	
	
	// Validate 'processing',
	
	
	
	// Validate 'posted',
	
	
	
	// Validate 'failed',
	
	
	
	// Validate 'cancelled',
	
	
	
	// Validate 'reversed'
	
	
	
	// Validate JournalEntryId
	
	
	if err := v.validateJournalEntryIdExists(ctx, tx, req.JournalEntryId); err != nil {
		return err
	}
	
	
	// Validate ReversalJournalEntryId
	
	
	if err := v.validateReversalJournalEntryIdExists(ctx, tx, req.ReversalJournalEntryId); err != nil {
		return err
	}
	
	
	// Validate PostingDate
	
	
	
	// Validate PostedAt
	
	
	
	// Validate PostedBy
	
	
	
	// Validate PostingMethod
	
	
	
	// Validate ErrorCode
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate ErrorDetails
	
	
	
	// Validate RetryCount
	
	
	
	// Validate LastRetryAt
	
	
	
	// Validate MaxRetries
	
	
	
	// Validate ReversedAt
	
	
	
	// Validate ReversedBy
	
	
	
	// Validate ReversalReason
	
	
	
	// Validate TotalDebit
	
	
	
	// Validate TotalCredit
	
	
	
	// Validate LineCount
	
	
	
	// Validate CurrencyCode
	
	
	
	// Validate PostingContext
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate (postingStatus
	
	if err := v.validate(postingStatus(req.(postingStatus); err != nil {
		return err
	}
	
	
	
	// Validate (postingStatus
	
	
	
	// Validate (postingStatus
	
	if err := v.validate(postingStatus(req.(postingStatus); err != nil {
		return err
	}
	
	
	
	// Validate (postingStatus
	
	
	
	// Validate (postingStatus
	
	
	
	// Validate (ABS(COALESCE(totalDebit,
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdatePosPostingAuditRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate SourceTable if provided
	
	if req.SourceTable != nil {
		if err := v.validateSourceTable(*req.SourceTable); err != nil {
			return err
		}
	}
	
	
	
	// Validate SourceId if provided
	
	
	if req.SourceId != nil {
		if err := v.validateSourceIdExists(ctx, tx, *req.SourceId); err != nil {
			return err
		}
	}
	
	
	// Validate SourceReference if provided
	
	
	
	// Validate PostingStatus if provided
	
	if req.PostingStatus != nil {
		if err := v.validatePostingStatus(*req.PostingStatus); err != nil {
			return err
		}
	}
	
	
	
	// Validate 'pending', if provided
	
	
	
	// Validate 'processing', if provided
	
	
	
	// Validate 'posted', if provided
	
	
	
	// Validate 'failed', if provided
	
	
	
	// Validate 'cancelled', if provided
	
	
	
	// Validate 'reversed' if provided
	
	
	
	// Validate JournalEntryId if provided
	
	
	if req.JournalEntryId != nil {
		if err := v.validateJournalEntryIdExists(ctx, tx, *req.JournalEntryId); err != nil {
			return err
		}
	}
	
	
	// Validate ReversalJournalEntryId if provided
	
	
	if req.ReversalJournalEntryId != nil {
		if err := v.validateReversalJournalEntryIdExists(ctx, tx, *req.ReversalJournalEntryId); err != nil {
			return err
		}
	}
	
	
	// Validate PostingDate if provided
	
	
	
	// Validate PostedAt if provided
	
	
	
	// Validate PostedBy if provided
	
	
	
	// Validate PostingMethod if provided
	
	
	
	// Validate ErrorCode if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate ErrorDetails if provided
	
	
	
	// Validate RetryCount if provided
	
	
	
	// Validate LastRetryAt if provided
	
	
	
	// Validate MaxRetries if provided
	
	
	
	// Validate ReversedAt if provided
	
	
	
	// Validate ReversedBy if provided
	
	
	
	// Validate ReversalReason if provided
	
	
	
	// Validate TotalDebit if provided
	
	
	
	// Validate TotalCredit if provided
	
	
	
	// Validate LineCount if provided
	
	
	
	// Validate CurrencyCode if provided
	
	
	
	// Validate PostingContext if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate (postingStatus if provided
	
	if req.(postingStatus != nil {
		if err := v.validate(postingStatus(*req.(postingStatus); err != nil {
			return err
		}
	}
	
	
	
	// Validate (postingStatus if provided
	
	
	
	// Validate (postingStatus if provided
	
	if req.(postingStatus != nil {
		if err := v.validate(postingStatus(*req.(postingStatus); err != nil {
			return err
		}
	}
	
	
	
	// Validate (postingStatus if provided
	
	
	
	// Validate (postingStatus if provided
	
	
	
	// Validate (ABS(COALESCE(totalDebit, if provided
	
	
	

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



// validateSourceTable validates source_table field
func (v *Validator) validateSourceTable(value string) error {
	
	// Add custom validation for source_table
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("source_table cannot be empty")
	}
	
	return nil
}







// validateSourceIdExists validates that source_id exists
func (v *Validator) validateSourceIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for source
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM source WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check source existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("source with id %s does not exist", id)
	}
	return nil
}







// validatePostingStatus validates posting_status field
func (v *Validator) validatePostingStatus(value string) error {
	
	// Add custom validation for posting_status
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("posting_status cannot be empty")
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





// validateReversalJournalEntryIdExists validates that reversal_journal_entry_id exists
func (v *Validator) validateReversalJournalEntryIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for reversal_journal_entry
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM reversal_journal_entry WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check reversal_journal_entry existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("reversal_journal_entry with id %s does not exist", id)
	}
	return nil
}



























































































// validate(postingStatus validates (posting_status field
func (v *Validator) validate(postingStatus(value string) error {
	
	// Add custom validation for (posting_status
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(posting_status cannot be empty")
	}
	
	return nil
}









// validate(postingStatus validates (posting_status field
func (v *Validator) validate(postingStatus(value string) error {
	
	// Add custom validation for (posting_status
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(posting_status cannot be empty")
	}
	
	return nil
}

















// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreatePosPostingAuditRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreatePosPostingAuditRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *PosPostingAudit, req *dto.UpdatePosPostingAuditRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *PosPostingAudit) error {
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

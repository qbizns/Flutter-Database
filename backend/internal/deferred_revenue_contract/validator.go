package deferred_revenue_contract

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles DeferredRevenueContracts validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new DeferredRevenueContracts validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateDeferredRevenueContractsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate CustomerInvoiceId
	
	
	if err := v.validateCustomerInvoiceIdExists(ctx, tx, req.CustomerInvoiceId); err != nil {
		return err
	}
	
	
	// Validate InvoiceLineId
	
	
	if err := v.validateInvoiceLineIdExists(ctx, tx, req.InvoiceLineId); err != nil {
		return err
	}
	
	
	// Validate ContractName
	
	
	
	// Validate TotalDeferredAmount
	
	
	
	// Validate StartDate
	
	
	
	// Validate EndDate
	
	
	
	// Validate RecognitionMethod
	
	
	
	// Validate RecognitionMethod
	
	
	
	// Validate DeferredAccountId
	
	
	if err := v.validateDeferredAccountIdExists(ctx, tx, req.DeferredAccountId); err != nil {
		return err
	}
	
	
	// Validate RevenueAccountId
	
	
	if err := v.validateRevenueAccountIdExists(ctx, tx, req.RevenueAccountId); err != nil {
		return err
	}
	
	
	// Validate Status
	
	
	
	// Validate RecognizedAmount
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateDeferredRevenueContractsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate CustomerInvoiceId if provided
	
	
	if req.CustomerInvoiceId != nil {
		if err := v.validateCustomerInvoiceIdExists(ctx, tx, *req.CustomerInvoiceId); err != nil {
			return err
		}
	}
	
	
	// Validate InvoiceLineId if provided
	
	
	if req.InvoiceLineId != nil {
		if err := v.validateInvoiceLineIdExists(ctx, tx, *req.InvoiceLineId); err != nil {
			return err
		}
	}
	
	
	// Validate ContractName if provided
	
	
	
	// Validate TotalDeferredAmount if provided
	
	
	
	// Validate StartDate if provided
	
	
	
	// Validate EndDate if provided
	
	
	
	// Validate RecognitionMethod if provided
	
	
	
	// Validate RecognitionMethod if provided
	
	
	
	// Validate DeferredAccountId if provided
	
	
	if req.DeferredAccountId != nil {
		if err := v.validateDeferredAccountIdExists(ctx, tx, *req.DeferredAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate RevenueAccountId if provided
	
	
	if req.RevenueAccountId != nil {
		if err := v.validateRevenueAccountIdExists(ctx, tx, *req.RevenueAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate Status if provided
	
	
	
	// Validate RecognizedAmount if provided
	
	
	
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





// validateCustomerInvoiceIdExists validates that customer_invoice_id exists
func (v *Validator) validateCustomerInvoiceIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for customer_invoice
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM customer_invoice WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check customer_invoice existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("customer_invoice with id %s does not exist", id)
	}
	return nil
}





// validateInvoiceLineIdExists validates that invoice_line_id exists
func (v *Validator) validateInvoiceLineIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for invoice_line
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM invoice_line WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check invoice_line existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("invoice_line with id %s does not exist", id)
	}
	return nil
}





























// validateDeferredAccountIdExists validates that deferred_account_id exists
func (v *Validator) validateDeferredAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for deferred_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM deferred_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check deferred_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("deferred_account with id %s does not exist", id)
	}
	return nil
}





// validateRevenueAccountIdExists validates that revenue_account_id exists
func (v *Validator) validateRevenueAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for revenue_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM revenue_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check revenue_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("revenue_account with id %s does not exist", id)
	}
	return nil
}























// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateDeferredRevenueContractsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateDeferredRevenueContractsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *DeferredRevenueContracts, req *UpdateDeferredRevenueContractsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *DeferredRevenueContracts) error {
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

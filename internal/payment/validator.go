package payment

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

// Validator handles Payments validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Payments validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreatePaymentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate SaleId
	
	
	if err := v.validateSaleIdExists(ctx, tx, req.SaleId); err != nil {
		return err
	}
	
	
	// Validate PaymentMethod
	
	if err := v.validatePaymentMethod(req.PaymentMethod); err != nil {
		return err
	}
	
	
	
	// Validate PaymentStatus
	
	if err := v.validatePaymentStatus(req.PaymentStatus); err != nil {
		return err
	}
	
	
	
	// Validate Amount
	
	
	
	// Validate CardLastFour
	
	
	
	// Validate CardType
	
	
	
	// Validate TransactionId
	
	
	if err := v.validateTransactionIdExists(ctx, tx, req.TransactionId); err != nil {
		return err
	}
	
	
	// Validate ReferenceNumber
	
	
	
	// Validate AccountNumber
	
	
	
	// Validate AccountName
	
	
	
	// Validate PaymentDate
	
	
	
	// Validate ProcessedAt
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdatePaymentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate SaleId if provided
	
	
	if req.SaleId != nil {
		if err := v.validateSaleIdExists(ctx, tx, *req.SaleId); err != nil {
			return err
		}
	}
	
	
	// Validate PaymentMethod if provided
	
	if req.PaymentMethod != nil {
		if err := v.validatePaymentMethod(*req.PaymentMethod); err != nil {
			return err
		}
	}
	
	
	
	// Validate PaymentStatus if provided
	
	if req.PaymentStatus != nil {
		if err := v.validatePaymentStatus(*req.PaymentStatus); err != nil {
			return err
		}
	}
	
	
	
	// Validate Amount if provided
	
	
	
	// Validate CardLastFour if provided
	
	
	
	// Validate CardType if provided
	
	
	
	// Validate TransactionId if provided
	
	
	if req.TransactionId != nil {
		if err := v.validateTransactionIdExists(ctx, tx, *req.TransactionId); err != nil {
			return err
		}
	}
	
	
	// Validate ReferenceNumber if provided
	
	
	
	// Validate AccountNumber if provided
	
	
	
	// Validate AccountName if provided
	
	
	
	// Validate PaymentDate if provided
	
	
	
	// Validate ProcessedAt if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	

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





// validateSaleIdExists validates that sale_id exists
func (v *Validator) validateSaleIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for sale
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sale WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check sale existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("sale with id %s does not exist", id)
	}
	return nil
}



// validatePaymentMethod validates payment_method field
func (v *Validator) validatePaymentMethod(value string) error {
	
	// Add custom validation for payment_method
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("payment_method cannot be empty")
	}
	
	return nil
}





// validatePaymentStatus validates payment_status field
func (v *Validator) validatePaymentStatus(value string) error {
	
	// Add custom validation for payment_status
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("payment_status cannot be empty")
	}
	
	return nil
}



















// validateTransactionIdExists validates that transaction_id exists
func (v *Validator) validateTransactionIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for transaction
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM transaction WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check transaction existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("transaction with id %s does not exist", id)
	}
	return nil
}



































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreatePaymentsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreatePaymentsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Payments, req *dto.UpdatePaymentsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Payments) error {
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

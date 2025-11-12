package sale

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

// Validator handles Sales validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Sales validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateSalesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate SaleNumber
	
	if err := v.validateSaleNumber(req.SaleNumber); err != nil {
		return err
	}
	
	
	
	// Validate ReferenceNumber
	
	
	
	// Validate TransactionType
	
	if err := v.validateTransactionType(req.TransactionType); err != nil {
		return err
	}
	
	
	
	// Validate CustomerId
	
	
	if err := v.validateCustomerIdExists(ctx, tx, req.CustomerId); err != nil {
		return err
	}
	
	
	// Validate CashierId
	
	
	if err := v.validateCashierIdExists(ctx, tx, req.CashierId); err != nil {
		return err
	}
	
	
	// Validate Subtotal
	
	
	
	// Validate TaxAmount
	
	
	
	// Validate DiscountAmount
	
	
	
	// Validate TotalAmount
	
	
	
	// Validate PaidAmount
	
	
	
	// Validate ChangeAmount
	
	
	
	// Validate OutstandingAmount
	
	
	
	// Validate PaymentStatus
	
	if err := v.validatePaymentStatus(req.PaymentStatus); err != nil {
		return err
	}
	
	
	
	// Validate DiscountType
	
	
	
	// Validate DiscountValue
	
	
	
	// Validate DiscountReason
	
	
	
	// Validate TransactionDate
	
	
	
	// Validate CompletedAt
	
	
	
	// Validate Notes
	
	
	
	// Validate InternalNotes
	
	
	
	// Validate CustomFields
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateSalesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate SaleNumber if provided
	
	if req.SaleNumber != nil {
		if err := v.validateSaleNumber(*req.SaleNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate ReferenceNumber if provided
	
	
	
	// Validate TransactionType if provided
	
	if req.TransactionType != nil {
		if err := v.validateTransactionType(*req.TransactionType); err != nil {
			return err
		}
	}
	
	
	
	// Validate CustomerId if provided
	
	
	if req.CustomerId != nil {
		if err := v.validateCustomerIdExists(ctx, tx, *req.CustomerId); err != nil {
			return err
		}
	}
	
	
	// Validate CashierId if provided
	
	
	if req.CashierId != nil {
		if err := v.validateCashierIdExists(ctx, tx, *req.CashierId); err != nil {
			return err
		}
	}
	
	
	// Validate Subtotal if provided
	
	
	
	// Validate TaxAmount if provided
	
	
	
	// Validate DiscountAmount if provided
	
	
	
	// Validate TotalAmount if provided
	
	
	
	// Validate PaidAmount if provided
	
	
	
	// Validate ChangeAmount if provided
	
	
	
	// Validate OutstandingAmount if provided
	
	
	
	// Validate PaymentStatus if provided
	
	if req.PaymentStatus != nil {
		if err := v.validatePaymentStatus(*req.PaymentStatus); err != nil {
			return err
		}
	}
	
	
	
	// Validate DiscountType if provided
	
	
	
	// Validate DiscountValue if provided
	
	
	
	// Validate DiscountReason if provided
	
	
	
	// Validate TransactionDate if provided
	
	
	
	// Validate CompletedAt if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate InternalNotes if provided
	
	
	
	// Validate CustomFields if provided
	
	
	
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



// validateSaleNumber validates sale_number field
func (v *Validator) validateSaleNumber(value string) error {
	
	// Add custom validation for sale_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("sale_number cannot be empty")
	}
	
	return nil
}









// validateTransactionType validates transaction_type field
func (v *Validator) validateTransactionType(value string) error {
	
	// Add custom validation for transaction_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("transaction_type cannot be empty")
	}
	
	return nil
}







// validateCustomerIdExists validates that customer_id exists
func (v *Validator) validateCustomerIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for customer
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM customer WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check customer existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("customer with id %s does not exist", id)
	}
	return nil
}





// validateCashierIdExists validates that cashier_id exists
func (v *Validator) validateCashierIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for cashier
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM cashier WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check cashier existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("cashier with id %s does not exist", id)
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

















































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateSalesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateSalesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Sales, req *dto.UpdateSalesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Sales) error {
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

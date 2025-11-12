package expens

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

// Validator handles Expenses validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Expenses validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateExpensesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate ExpenseNumber
	
	if err := v.validateExpenseNumber(req.ExpenseNumber); err != nil {
		return err
	}
	
	
	
	// Validate ExpenseDate
	
	
	
	// Validate Category
	
	if err := v.validateCategory(req.Category); err != nil {
		return err
	}
	
	
	
	// Validate Subcategory
	
	
	
	// Validate PayeeName
	
	if err := v.validatePayeeName(req.PayeeName); err != nil {
		return err
	}
	
	
	
	// Validate PaymentMethod
	
	
	
	// Validate Amount
	
	
	
	// Validate TaxAmount
	
	
	
	// Validate TotalAmount
	
	
	
	// Validate Currency
	
	
	
	// Validate Status
	
	
	
	// Validate ReferenceNumber
	
	
	
	// Validate PurchaseOrderId
	
	
	if err := v.validatePurchaseOrderIdExists(ctx, tx, req.PurchaseOrderId); err != nil {
		return err
	}
	
	
	// Validate ReceiptUrl
	
	
	
	// Validate AttachmentUrls
	
	
	
	// Validate Description
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate ApprovedBy
	
	
	
	// Validate ApprovedAt
	
	
	
	// Validate Amount
	
	
	
	// Validate TaxAmount
	
	
	
	// Validate TotalAmount
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateExpensesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate LocationId if provided
	
	
	if req.LocationId != nil {
		if err := v.validateLocationIdExists(ctx, tx, *req.LocationId); err != nil {
			return err
		}
	}
	
	
	// Validate ExpenseNumber if provided
	
	if req.ExpenseNumber != nil {
		if err := v.validateExpenseNumber(*req.ExpenseNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate ExpenseDate if provided
	
	
	
	// Validate Category if provided
	
	if req.Category != nil {
		if err := v.validateCategory(*req.Category); err != nil {
			return err
		}
	}
	
	
	
	// Validate Subcategory if provided
	
	
	
	// Validate PayeeName if provided
	
	if req.PayeeName != nil {
		if err := v.validatePayeeName(*req.PayeeName); err != nil {
			return err
		}
	}
	
	
	
	// Validate PaymentMethod if provided
	
	
	
	// Validate Amount if provided
	
	
	
	// Validate TaxAmount if provided
	
	
	
	// Validate TotalAmount if provided
	
	
	
	// Validate Currency if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate ReferenceNumber if provided
	
	
	
	// Validate PurchaseOrderId if provided
	
	
	if req.PurchaseOrderId != nil {
		if err := v.validatePurchaseOrderIdExists(ctx, tx, *req.PurchaseOrderId); err != nil {
			return err
		}
	}
	
	
	// Validate ReceiptUrl if provided
	
	
	
	// Validate AttachmentUrls if provided
	
	
	
	// Validate Description if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate ApprovedBy if provided
	
	
	
	// Validate ApprovedAt if provided
	
	
	
	// Validate Amount if provided
	
	
	
	// Validate TaxAmount if provided
	
	
	
	// Validate TotalAmount if provided
	
	
	

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



// validateExpenseNumber validates expense_number field
func (v *Validator) validateExpenseNumber(value string) error {
	
	// Add custom validation for expense_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("expense_number cannot be empty")
	}
	
	return nil
}









// validateCategory validates category field
func (v *Validator) validateCategory(value string) error {
	
	// Add custom validation for category
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("category cannot be empty")
	}
	
	return nil
}









// validatePayeeName validates payee_name field
func (v *Validator) validatePayeeName(value string) error {
	
	// Add custom validation for payee_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("payee_name cannot be empty")
	}
	
	return nil
}



































// validatePurchaseOrderIdExists validates that purchase_order_id exists
func (v *Validator) validatePurchaseOrderIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for purchase_order
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM purchase_order WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check purchase_order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("purchase_order with id %s does not exist", id)
	}
	return nil
}



















































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateExpensesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateExpensesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Expenses, req *dto.UpdateExpensesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Expenses) error {
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

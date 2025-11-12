package sale_return_item

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles SaleReturnItems validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new SaleReturnItems validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateSaleReturnItemsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate SaleReturnId
	
	
	if err := v.validateSaleReturnIdExists(ctx, tx, req.SaleReturnId); err != nil {
		return err
	}
	
	
	// Validate OriginalSaleItemId
	
	
	if err := v.validateOriginalSaleItemIdExists(ctx, tx, req.OriginalSaleItemId); err != nil {
		return err
	}
	
	
	// Validate ProductId
	
	
	if err := v.validateProductIdExists(ctx, tx, req.ProductId); err != nil {
		return err
	}
	
	
	// Validate ProductVariantId
	
	
	if err := v.validateProductVariantIdExists(ctx, tx, req.ProductVariantId); err != nil {
		return err
	}
	
	
	// Validate Quantity
	
	
	
	// Validate UnitPrice
	
	
	
	// Validate Subtotal
	
	
	
	// Validate TaxAmount
	
	
	
	// Validate DiscountAmount
	
	
	
	// Validate TotalAmount
	
	
	
	// Validate ReturnReasonId
	
	
	if err := v.validateReturnReasonIdExists(ctx, tx, req.ReturnReasonId); err != nil {
		return err
	}
	
	
	// Validate ReturnReasonNotes
	
	
	
	// Validate ItemCondition
	
	
	
	// Validate ItemCondition
	
	
	
	// Validate IsRestockable
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateSaleReturnItemsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate SaleReturnId if provided
	
	
	if req.SaleReturnId != nil {
		if err := v.validateSaleReturnIdExists(ctx, tx, *req.SaleReturnId); err != nil {
			return err
		}
	}
	
	
	// Validate OriginalSaleItemId if provided
	
	
	if req.OriginalSaleItemId != nil {
		if err := v.validateOriginalSaleItemIdExists(ctx, tx, *req.OriginalSaleItemId); err != nil {
			return err
		}
	}
	
	
	// Validate ProductId if provided
	
	
	if req.ProductId != nil {
		if err := v.validateProductIdExists(ctx, tx, *req.ProductId); err != nil {
			return err
		}
	}
	
	
	// Validate ProductVariantId if provided
	
	
	if req.ProductVariantId != nil {
		if err := v.validateProductVariantIdExists(ctx, tx, *req.ProductVariantId); err != nil {
			return err
		}
	}
	
	
	// Validate Quantity if provided
	
	
	
	// Validate UnitPrice if provided
	
	
	
	// Validate Subtotal if provided
	
	
	
	// Validate TaxAmount if provided
	
	
	
	// Validate DiscountAmount if provided
	
	
	
	// Validate TotalAmount if provided
	
	
	
	// Validate ReturnReasonId if provided
	
	
	if req.ReturnReasonId != nil {
		if err := v.validateReturnReasonIdExists(ctx, tx, *req.ReturnReasonId); err != nil {
			return err
		}
	}
	
	
	// Validate ReturnReasonNotes if provided
	
	
	
	// Validate ItemCondition if provided
	
	
	
	// Validate ItemCondition if provided
	
	
	
	// Validate IsRestockable if provided
	
	
	

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





// validateSaleReturnIdExists validates that sale_return_id exists
func (v *Validator) validateSaleReturnIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for sale_return
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sale_return WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check sale_return existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("sale_return with id %s does not exist", id)
	}
	return nil
}





// validateOriginalSaleItemIdExists validates that original_sale_item_id exists
func (v *Validator) validateOriginalSaleItemIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for original_sale_item
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM original_sale_item WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check original_sale_item existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("original_sale_item with id %s does not exist", id)
	}
	return nil
}





// validateProductIdExists validates that product_id exists
func (v *Validator) validateProductIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for product
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check product existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("product with id %s does not exist", id)
	}
	return nil
}





// validateProductVariantIdExists validates that product_variant_id exists
func (v *Validator) validateProductVariantIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for product_variant
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM product_variant WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check product_variant existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("product_variant with id %s does not exist", id)
	}
	return nil
}





























// validateReturnReasonIdExists validates that return_reason_id exists
func (v *Validator) validateReturnReasonIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for return_reason
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM return_reason WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check return_reason existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("return_reason with id %s does not exist", id)
	}
	return nil
}



















// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateSaleReturnItemsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateSaleReturnItemsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *SaleReturnItems, req *UpdateSaleReturnItemsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *SaleReturnItems) error {
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

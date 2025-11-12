package goods_receipt_item

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

// Validator handles GoodsReceiptItems validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new GoodsReceiptItems validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateGoodsReceiptItemsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate GoodsReceiptId
	
	
	if err := v.validateGoodsReceiptIdExists(ctx, tx, req.GoodsReceiptId); err != nil {
		return err
	}
	
	
	// Validate PurchaseOrderItemId
	
	
	if err := v.validatePurchaseOrderItemIdExists(ctx, tx, req.PurchaseOrderItemId); err != nil {
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
	
	
	// Validate QuantityReceived
	
	
	
	// Validate QuantityAccepted
	
	
	
	// Validate QuantityRejected
	
	
	
	// Validate RejectionReason
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateGoodsReceiptItemsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate GoodsReceiptId if provided
	
	
	if req.GoodsReceiptId != nil {
		if err := v.validateGoodsReceiptIdExists(ctx, tx, *req.GoodsReceiptId); err != nil {
			return err
		}
	}
	
	
	// Validate PurchaseOrderItemId if provided
	
	
	if req.PurchaseOrderItemId != nil {
		if err := v.validatePurchaseOrderItemIdExists(ctx, tx, *req.PurchaseOrderItemId); err != nil {
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
	
	
	// Validate QuantityReceived if provided
	
	
	
	// Validate QuantityAccepted if provided
	
	
	
	// Validate QuantityRejected if provided
	
	
	
	// Validate RejectionReason if provided
	
	
	

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





// validateGoodsReceiptIdExists validates that goods_receipt_id exists
func (v *Validator) validateGoodsReceiptIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for goods_receipt
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM goods_receipt WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check goods_receipt existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("goods_receipt with id %s does not exist", id)
	}
	return nil
}





// validatePurchaseOrderItemIdExists validates that purchase_order_item_id exists
func (v *Validator) validatePurchaseOrderItemIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for purchase_order_item
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM purchase_order_item WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check purchase_order_item existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("purchase_order_item with id %s does not exist", id)
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



















// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateGoodsReceiptItemsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateGoodsReceiptItemsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *GoodsReceiptItems, req *dto.UpdateGoodsReceiptItemsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *GoodsReceiptItems) error {
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

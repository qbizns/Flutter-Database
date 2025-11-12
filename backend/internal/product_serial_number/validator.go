package product_serial_number

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

// Validator handles ProductSerialNumbers validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ProductSerialNumbers validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateProductSerialNumbersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ProductId
	
	
	if err := v.validateProductIdExists(ctx, tx, req.ProductId); err != nil {
		return err
	}
	
	
	// Validate ProductVariantId
	
	
	if err := v.validateProductVariantIdExists(ctx, tx, req.ProductVariantId); err != nil {
		return err
	}
	
	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate SerialNumber
	
	if err := v.validateSerialNumber(req.SerialNumber); err != nil {
		return err
	}
	
	
	
	// Validate Status
	
	
	
	// Validate PurchaseOrderId
	
	
	if err := v.validatePurchaseOrderIdExists(ctx, tx, req.PurchaseOrderId); err != nil {
		return err
	}
	
	
	// Validate PurchaseDate
	
	
	
	// Validate PurchaseCost
	
	
	
	// Validate SupplierId
	
	
	if err := v.validateSupplierIdExists(ctx, tx, req.SupplierId); err != nil {
		return err
	}
	
	
	// Validate SaleId
	
	
	if err := v.validateSaleIdExists(ctx, tx, req.SaleId); err != nil {
		return err
	}
	
	
	// Validate SaleDate
	
	
	
	// Validate SalePrice
	
	
	
	// Validate CustomerId
	
	
	if err := v.validateCustomerIdExists(ctx, tx, req.CustomerId); err != nil {
		return err
	}
	
	
	// Validate WarrantyStartDate
	
	
	
	// Validate WarrantyEndDate
	
	
	
	// Validate WarrantyProvider
	
	
	
	// Validate WarrantyTerms
	
	
	
	// Validate Notes
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateProductSerialNumbersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
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
	
	
	// Validate LocationId if provided
	
	
	if req.LocationId != nil {
		if err := v.validateLocationIdExists(ctx, tx, *req.LocationId); err != nil {
			return err
		}
	}
	
	
	// Validate SerialNumber if provided
	
	if req.SerialNumber != nil {
		if err := v.validateSerialNumber(*req.SerialNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate Status if provided
	
	
	
	// Validate PurchaseOrderId if provided
	
	
	if req.PurchaseOrderId != nil {
		if err := v.validatePurchaseOrderIdExists(ctx, tx, *req.PurchaseOrderId); err != nil {
			return err
		}
	}
	
	
	// Validate PurchaseDate if provided
	
	
	
	// Validate PurchaseCost if provided
	
	
	
	// Validate SupplierId if provided
	
	
	if req.SupplierId != nil {
		if err := v.validateSupplierIdExists(ctx, tx, *req.SupplierId); err != nil {
			return err
		}
	}
	
	
	// Validate SaleId if provided
	
	
	if req.SaleId != nil {
		if err := v.validateSaleIdExists(ctx, tx, *req.SaleId); err != nil {
			return err
		}
	}
	
	
	// Validate SaleDate if provided
	
	
	
	// Validate SalePrice if provided
	
	
	
	// Validate CustomerId if provided
	
	
	if req.CustomerId != nil {
		if err := v.validateCustomerIdExists(ctx, tx, *req.CustomerId); err != nil {
			return err
		}
	}
	
	
	// Validate WarrantyStartDate if provided
	
	
	
	// Validate WarrantyEndDate if provided
	
	
	
	// Validate WarrantyProvider if provided
	
	
	
	// Validate WarrantyTerms if provided
	
	
	
	// Validate Notes if provided
	
	
	
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



// validateSerialNumber validates serial_number field
func (v *Validator) validateSerialNumber(value string) error {
	
	// Add custom validation for serial_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("serial_number cannot be empty")
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













// validateSupplierIdExists validates that supplier_id exists
func (v *Validator) validateSupplierIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for supplier
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM supplier WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check supplier existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("supplier with id %s does not exist", id)
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



































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateProductSerialNumbersRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateProductSerialNumbersRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ProductSerialNumbers, req *dto.UpdateProductSerialNumbersRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ProductSerialNumbers) error {
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

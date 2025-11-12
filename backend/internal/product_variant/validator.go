package product_variant

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

// Validator handles ProductVariants validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ProductVariants validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateProductVariantsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ProductId
	
	
	if err := v.validateProductIdExists(ctx, tx, req.ProductId); err != nil {
		return err
	}
	
	
	// Validate VariantName
	
	if err := v.validateVariantName(req.VariantName); err != nil {
		return err
	}
	
	
	
	// Validate Sku
	
	
	
	// Validate Barcode
	
	
	
	// Validate Attributes
	
	
	
	// Validate CostPrice
	
	
	
	// Validate SellingPrice
	
	
	
	// Validate CompareAtPrice
	
	
	
	// Validate CurrentStock
	
	
	
	// Validate ReorderLevel
	
	
	
	// Validate ReorderQuantity
	
	
	
	// Validate Weight
	
	
	
	// Validate WeightUnit
	
	
	
	// Validate Dimensions
	
	
	
	// Validate IsActive
	
	
	
	// Validate IsDefault
	
	
	
	// Validate SortOrder
	
	
	
	// Validate ImageUrl
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate (costPrice
	
	
	
	// Validate (sellingPrice
	
	
	
	// Validate (compareAtPrice
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateProductVariantsRequest) error {
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
	
	
	// Validate VariantName if provided
	
	if req.VariantName != nil {
		if err := v.validateVariantName(*req.VariantName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Sku if provided
	
	
	
	// Validate Barcode if provided
	
	
	
	// Validate Attributes if provided
	
	
	
	// Validate CostPrice if provided
	
	
	
	// Validate SellingPrice if provided
	
	
	
	// Validate CompareAtPrice if provided
	
	
	
	// Validate CurrentStock if provided
	
	
	
	// Validate ReorderLevel if provided
	
	
	
	// Validate ReorderQuantity if provided
	
	
	
	// Validate Weight if provided
	
	
	
	// Validate WeightUnit if provided
	
	
	
	// Validate Dimensions if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate IsDefault if provided
	
	
	
	// Validate SortOrder if provided
	
	
	
	// Validate ImageUrl if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate (costPrice if provided
	
	
	
	// Validate (sellingPrice if provided
	
	
	
	// Validate (compareAtPrice if provided
	
	
	

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



// validateVariantName validates variant_name field
func (v *Validator) validateVariantName(value string) error {
	
	// Add custom validation for variant_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("variant_name cannot be empty")
	}
	
	return nil
}

































































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateProductVariantsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateProductVariantsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ProductVariants, req *dto.UpdateProductVariantsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ProductVariants) error {
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

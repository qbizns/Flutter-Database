package product

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

// Validator handles Products validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Products validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateProductsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate Sku
	
	
	
	// Validate Barcode
	
	
	
	// Validate Name
	
	if err := v.validateName(req.Name); err != nil {
		return err
	}
	
	
	
	// Validate Description
	
	
	
	// Validate CategoryId
	
	
	if err := v.validateCategoryIdExists(ctx, tx, req.CategoryId); err != nil {
		return err
	}
	
	
	// Validate CostPrice
	
	
	
	// Validate SellingPrice
	
	
	
	// Validate CompareAtPrice
	
	
	
	// Validate TaxRate
	
	
	
	// Validate IsTaxInclusive
	
	
	
	// Validate TrackInventory
	
	
	
	// Validate CurrentStock
	
	
	
	// Validate LowStockThreshold
	
	
	
	// Validate Unit
	
	
	
	// Validate IsService
	
	
	
	// Validate IsComposite
	
	
	
	// Validate HasVariants
	
	
	
	// Validate ImageUrl
	
	
	
	// Validate Images
	
	
	
	// Validate SortOrder
	
	
	
	// Validate IsActive
	
	
	
	// Validate IsFeatured
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateProductsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate Sku if provided
	
	
	
	// Validate Barcode if provided
	
	
	
	// Validate Name if provided
	
	if req.Name != nil {
		if err := v.validateName(*req.Name); err != nil {
			return err
		}
	}
	
	
	
	// Validate Description if provided
	
	
	
	// Validate CategoryId if provided
	
	
	if req.CategoryId != nil {
		if err := v.validateCategoryIdExists(ctx, tx, *req.CategoryId); err != nil {
			return err
		}
	}
	
	
	// Validate CostPrice if provided
	
	
	
	// Validate SellingPrice if provided
	
	
	
	// Validate CompareAtPrice if provided
	
	
	
	// Validate TaxRate if provided
	
	
	
	// Validate IsTaxInclusive if provided
	
	
	
	// Validate TrackInventory if provided
	
	
	
	// Validate CurrentStock if provided
	
	
	
	// Validate LowStockThreshold if provided
	
	
	
	// Validate Unit if provided
	
	
	
	// Validate IsService if provided
	
	
	
	// Validate IsComposite if provided
	
	
	
	// Validate HasVariants if provided
	
	
	
	// Validate ImageUrl if provided
	
	
	
	// Validate Images if provided
	
	
	
	// Validate SortOrder if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate IsFeatured if provided
	
	
	
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











// validateName validates name field
func (v *Validator) validateName(value string) error {
	
	// Add custom validation for name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	
	return nil
}











// validateCategoryIdExists validates that category_id exists
func (v *Validator) validateCategoryIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for category
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM category WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check category existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("category with id %s does not exist", id)
	}
	return nil
}























































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateProductsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateProductsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Products, req *dto.UpdateProductsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Products) error {
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

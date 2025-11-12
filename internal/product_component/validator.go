package product_component

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

// Validator handles ProductComponents validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ProductComponents validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateProductComponentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ParentProductId
	
	
	if err := v.validateParentProductIdExists(ctx, tx, req.ParentProductId); err != nil {
		return err
	}
	
	
	// Validate ComponentProductId
	
	
	if err := v.validateComponentProductIdExists(ctx, tx, req.ComponentProductId); err != nil {
		return err
	}
	
	
	// Validate ComponentVariantId
	
	
	if err := v.validateComponentVariantIdExists(ctx, tx, req.ComponentVariantId); err != nil {
		return err
	}
	
	
	// Validate Quantity
	
	
	
	// Validate InheritPrice
	
	
	
	// Validate PriceOverride
	
	
	
	// Validate DisplayOrder
	
	
	
	// Validate IsOptional
	
	
	
	// Validate ComponentProductId
	
	if err := v.validateComponentProductId(req.ComponentProductId); err != nil {
		return err
	}
	
	
	if err := v.validateComponentProductIdExists(ctx, tx, req.ComponentProductId); err != nil {
		return err
	}
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateProductComponentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate ParentProductId if provided
	
	
	if req.ParentProductId != nil {
		if err := v.validateParentProductIdExists(ctx, tx, *req.ParentProductId); err != nil {
			return err
		}
	}
	
	
	// Validate ComponentProductId if provided
	
	
	if req.ComponentProductId != nil {
		if err := v.validateComponentProductIdExists(ctx, tx, *req.ComponentProductId); err != nil {
			return err
		}
	}
	
	
	// Validate ComponentVariantId if provided
	
	
	if req.ComponentVariantId != nil {
		if err := v.validateComponentVariantIdExists(ctx, tx, *req.ComponentVariantId); err != nil {
			return err
		}
	}
	
	
	// Validate Quantity if provided
	
	
	
	// Validate InheritPrice if provided
	
	
	
	// Validate PriceOverride if provided
	
	
	
	// Validate DisplayOrder if provided
	
	
	
	// Validate IsOptional if provided
	
	
	
	// Validate ComponentProductId if provided
	
	if req.ComponentProductId != nil {
		if err := v.validateComponentProductId(*req.ComponentProductId); err != nil {
			return err
		}
	}
	
	
	if req.ComponentProductId != nil {
		if err := v.validateComponentProductIdExists(ctx, tx, *req.ComponentProductId); err != nil {
			return err
		}
	}
	
	

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





// validateParentProductIdExists validates that parent_product_id exists
func (v *Validator) validateParentProductIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for parent_product
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM parent_product WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check parent_product existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("parent_product with id %s does not exist", id)
	}
	return nil
}





// validateComponentProductIdExists validates that component_product_id exists
func (v *Validator) validateComponentProductIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for component_product
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM component_product WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check component_product existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("component_product with id %s does not exist", id)
	}
	return nil
}





// validateComponentVariantIdExists validates that component_variant_id exists
func (v *Validator) validateComponentVariantIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for component_variant
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM component_variant WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check component_variant existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("component_variant with id %s does not exist", id)
	}
	return nil
}























// validateComponentProductId validates component_product_id field
func (v *Validator) validateComponentProductId(value string) error {
	
	// Add custom validation for component_product_id
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("component_product_id cannot be empty")
	}
	
	return nil
}



// validateComponentProductIdExists validates that component_product_id exists
func (v *Validator) validateComponentProductIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for component_product
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM component_product WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check component_product existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("component_product with id %s does not exist", id)
	}
	return nil
}



// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateProductComponentsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateProductComponentsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ProductComponents, req *dto.UpdateProductComponentsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ProductComponents) error {
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

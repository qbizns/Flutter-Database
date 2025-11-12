package order_item_modifier

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles OrderItemModifiers validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new OrderItemModifiers validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateOrderItemModifiersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate OrderItemId
	
	
	if err := v.validateOrderItemIdExists(ctx, tx, req.OrderItemId); err != nil {
		return err
	}
	
	
	// Validate ModifierId
	
	
	if err := v.validateModifierIdExists(ctx, tx, req.ModifierId); err != nil {
		return err
	}
	
	
	// Validate ModifierGroupId
	
	
	if err := v.validateModifierGroupIdExists(ctx, tx, req.ModifierGroupId); err != nil {
		return err
	}
	
	
	// Validate ModifierName
	
	if err := v.validateModifierName(req.ModifierName); err != nil {
		return err
	}
	
	
	
	// Validate Quantity
	
	
	
	// Validate PriceAdjustment
	
	
	
	// Validate DisplayOrder
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateOrderItemModifiersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate OrderItemId if provided
	
	
	if req.OrderItemId != nil {
		if err := v.validateOrderItemIdExists(ctx, tx, *req.OrderItemId); err != nil {
			return err
		}
	}
	
	
	// Validate ModifierId if provided
	
	
	if req.ModifierId != nil {
		if err := v.validateModifierIdExists(ctx, tx, *req.ModifierId); err != nil {
			return err
		}
	}
	
	
	// Validate ModifierGroupId if provided
	
	
	if req.ModifierGroupId != nil {
		if err := v.validateModifierGroupIdExists(ctx, tx, *req.ModifierGroupId); err != nil {
			return err
		}
	}
	
	
	// Validate ModifierName if provided
	
	if req.ModifierName != nil {
		if err := v.validateModifierName(*req.ModifierName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Quantity if provided
	
	
	
	// Validate PriceAdjustment if provided
	
	
	
	// Validate DisplayOrder if provided
	
	
	
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





// validateOrderItemIdExists validates that order_item_id exists
func (v *Validator) validateOrderItemIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for order_item
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM order_item WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check order_item existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("order_item with id %s does not exist", id)
	}
	return nil
}





// validateModifierIdExists validates that modifier_id exists
func (v *Validator) validateModifierIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for modifier
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM modifier WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check modifier existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("modifier with id %s does not exist", id)
	}
	return nil
}





// validateModifierGroupIdExists validates that modifier_group_id exists
func (v *Validator) validateModifierGroupIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for modifier_group
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM modifier_group WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check modifier_group existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("modifier_group with id %s does not exist", id)
	}
	return nil
}



// validateModifierName validates modifier_name field
func (v *Validator) validateModifierName(value string) error {
	
	// Add custom validation for modifier_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("modifier_name cannot be empty")
	}
	
	return nil
}





























// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateOrderItemModifiersRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateOrderItemModifiersRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *OrderItemModifiers, req *UpdateOrderItemModifiersRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *OrderItemModifiers) error {
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

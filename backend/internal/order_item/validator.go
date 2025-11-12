package order_item

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles OrderItems validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new OrderItems validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateOrderItemsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate OrderId
	
	
	if err := v.validateOrderIdExists(ctx, tx, req.OrderId); err != nil {
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
	
	
	// Validate ItemName
	
	if err := v.validateItemName(req.ItemName); err != nil {
		return err
	}
	
	
	
	// Validate Quantity
	
	
	
	// Validate UnitPrice
	
	
	
	// Validate CourseId
	
	
	if err := v.validateCourseIdExists(ctx, tx, req.CourseId); err != nil {
		return err
	}
	
	
	// Validate CoursePosition
	
	
	
	// Validate FireTime
	
	
	
	// Validate KitchenStationId
	
	
	if err := v.validateKitchenStationIdExists(ctx, tx, req.KitchenStationId); err != nil {
		return err
	}
	
	
	// Validate KitchenTicketId
	
	
	if err := v.validateKitchenTicketIdExists(ctx, tx, req.KitchenTicketId); err != nil {
		return err
	}
	
	
	// Validate Status
	
	
	
	// Validate FiredAt
	
	
	
	// Validate AcknowledgedAt
	
	
	
	// Validate StartedPreparingAt
	
	
	
	// Validate ReadyAt
	
	
	
	// Validate ServedAt
	
	
	
	// Validate ModifiersTotal
	
	
	
	// Validate DiscountAmount
	
	
	
	// Validate LineTotal
	
	
	
	// Validate SpecialInstructions
	
	
	
	// Validate CustomerNotes
	
	
	
	// Validate KitchenNotes
	
	
	
	// Validate SeatNumber
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'pending',
	
	
	
	// Validate 'served',
	
	
	
	// Validate UnitPrice
	
	
	
	// Validate DiscountAmount
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateOrderItemsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate OrderId if provided
	
	
	if req.OrderId != nil {
		if err := v.validateOrderIdExists(ctx, tx, *req.OrderId); err != nil {
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
	
	
	// Validate ItemName if provided
	
	if req.ItemName != nil {
		if err := v.validateItemName(*req.ItemName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Quantity if provided
	
	
	
	// Validate UnitPrice if provided
	
	
	
	// Validate CourseId if provided
	
	
	if req.CourseId != nil {
		if err := v.validateCourseIdExists(ctx, tx, *req.CourseId); err != nil {
			return err
		}
	}
	
	
	// Validate CoursePosition if provided
	
	
	
	// Validate FireTime if provided
	
	
	
	// Validate KitchenStationId if provided
	
	
	if req.KitchenStationId != nil {
		if err := v.validateKitchenStationIdExists(ctx, tx, *req.KitchenStationId); err != nil {
			return err
		}
	}
	
	
	// Validate KitchenTicketId if provided
	
	
	if req.KitchenTicketId != nil {
		if err := v.validateKitchenTicketIdExists(ctx, tx, *req.KitchenTicketId); err != nil {
			return err
		}
	}
	
	
	// Validate Status if provided
	
	
	
	// Validate FiredAt if provided
	
	
	
	// Validate AcknowledgedAt if provided
	
	
	
	// Validate StartedPreparingAt if provided
	
	
	
	// Validate ReadyAt if provided
	
	
	
	// Validate ServedAt if provided
	
	
	
	// Validate ModifiersTotal if provided
	
	
	
	// Validate DiscountAmount if provided
	
	
	
	// Validate LineTotal if provided
	
	
	
	// Validate SpecialInstructions if provided
	
	
	
	// Validate CustomerNotes if provided
	
	
	
	// Validate KitchenNotes if provided
	
	
	
	// Validate SeatNumber if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'pending', if provided
	
	
	
	// Validate 'served', if provided
	
	
	
	// Validate UnitPrice if provided
	
	
	
	// Validate DiscountAmount if provided
	
	
	

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





// validateOrderIdExists validates that order_id exists
func (v *Validator) validateOrderIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for order
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM order WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("order with id %s does not exist", id)
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



// validateItemName validates item_name field
func (v *Validator) validateItemName(value string) error {
	
	// Add custom validation for item_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("item_name cannot be empty")
	}
	
	return nil
}















// validateCourseIdExists validates that course_id exists
func (v *Validator) validateCourseIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for course
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM course WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check course existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("course with id %s does not exist", id)
	}
	return nil
}













// validateKitchenStationIdExists validates that kitchen_station_id exists
func (v *Validator) validateKitchenStationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for kitchen_station
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM kitchen_station WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check kitchen_station existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("kitchen_station with id %s does not exist", id)
	}
	return nil
}





// validateKitchenTicketIdExists validates that kitchen_ticket_id exists
func (v *Validator) validateKitchenTicketIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for kitchen_ticket
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM kitchen_ticket WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check kitchen_ticket existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("kitchen_ticket with id %s does not exist", id)
	}
	return nil
}



















































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateOrderItemsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateOrderItemsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *OrderItems, req *UpdateOrderItemsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *OrderItems) error {
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

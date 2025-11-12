package order_tracking_event

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles OrderTrackingEvents validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new OrderTrackingEvents validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateOrderTrackingEventsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate OrderId
	
	
	if err := v.validateOrderIdExists(ctx, tx, req.OrderId); err != nil {
		return err
	}
	
	
	// Validate DeliveryAssignmentId
	
	
	if err := v.validateDeliveryAssignmentIdExists(ctx, tx, req.DeliveryAssignmentId); err != nil {
		return err
	}
	
	
	// Validate EventType
	
	if err := v.validateEventType(req.EventType); err != nil {
		return err
	}
	
	
	
	// Validate EventTimestamp
	
	
	
	// Validate EventMessage
	
	
	
	// Validate Location
	
	
	
	// Validate LocationName
	
	
	
	// Validate ActorType
	
	
	
	// Validate ActorId
	
	
	if err := v.validateActorIdExists(ctx, tx, req.ActorId); err != nil {
		return err
	}
	
	
	// Validate ActorName
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate 'orderPlaced',
	
	
	
	// Validate 'readyForPickup',
	
	
	
	// Validate 'arrived',
	
	
	
	// Validate 'rescheduled',
	
	
	
	// Validate 'system',
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateOrderTrackingEventsRequest) error {
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
	
	
	// Validate DeliveryAssignmentId if provided
	
	
	if req.DeliveryAssignmentId != nil {
		if err := v.validateDeliveryAssignmentIdExists(ctx, tx, *req.DeliveryAssignmentId); err != nil {
			return err
		}
	}
	
	
	// Validate EventType if provided
	
	if req.EventType != nil {
		if err := v.validateEventType(*req.EventType); err != nil {
			return err
		}
	}
	
	
	
	// Validate EventTimestamp if provided
	
	
	
	// Validate EventMessage if provided
	
	
	
	// Validate Location if provided
	
	
	
	// Validate LocationName if provided
	
	
	
	// Validate ActorType if provided
	
	
	
	// Validate ActorId if provided
	
	
	if req.ActorId != nil {
		if err := v.validateActorIdExists(ctx, tx, *req.ActorId); err != nil {
			return err
		}
	}
	
	
	// Validate ActorName if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate 'orderPlaced', if provided
	
	
	
	// Validate 'readyForPickup', if provided
	
	
	
	// Validate 'arrived', if provided
	
	
	
	// Validate 'rescheduled', if provided
	
	
	
	// Validate 'system', if provided
	
	
	

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





// validateDeliveryAssignmentIdExists validates that delivery_assignment_id exists
func (v *Validator) validateDeliveryAssignmentIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for delivery_assignment
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM delivery_assignment WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check delivery_assignment existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("delivery_assignment with id %s does not exist", id)
	}
	return nil
}



// validateEventType validates event_type field
func (v *Validator) validateEventType(value string) error {
	
	// Add custom validation for event_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("event_type cannot be empty")
	}
	
	return nil
}



























// validateActorIdExists validates that actor_id exists
func (v *Validator) validateActorIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for actor
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM actor WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check actor existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("actor with id %s does not exist", id)
	}
	return nil
}



































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateOrderTrackingEventsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateOrderTrackingEventsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *OrderTrackingEvents, req *UpdateOrderTrackingEventsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *OrderTrackingEvents) error {
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

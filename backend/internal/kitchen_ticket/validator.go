package kitchen_ticket

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

// Validator handles KitchenTickets validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new KitchenTickets validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateKitchenTicketsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate TicketNumber
	
	if err := v.validateTicketNumber(req.TicketNumber); err != nil {
		return err
	}
	
	
	
	// Validate DisplaySequence
	
	
	
	// Validate OrderId
	
	
	if err := v.validateOrderIdExists(ctx, tx, req.OrderId); err != nil {
		return err
	}
	
	
	// Validate KitchenStationId
	
	
	if err := v.validateKitchenStationIdExists(ctx, tx, req.KitchenStationId); err != nil {
		return err
	}
	
	
	// Validate CourseId
	
	
	if err := v.validateCourseIdExists(ctx, tx, req.CourseId); err != nil {
		return err
	}
	
	
	// Validate TicketType
	
	
	
	// Validate Priority
	
	
	
	// Validate Status
	
	
	
	// Validate FiredAt
	
	
	
	// Validate AcknowledgedAt
	
	
	
	// Validate StartedAt
	
	
	
	// Validate ReadyAt
	
	
	
	// Validate BumpedAt
	
	
	
	// Validate CompletedAt
	
	
	
	// Validate PrepTimeMinutes
	
	
	
	// Validate TargetPrepTime
	
	
	
	// Validate TableNumber
	
	
	
	// Validate OrderType
	
	
	
	// Validate Covers
	
	
	
	// Validate WaiterName
	
	
	
	// Validate SpecialInstructions
	
	
	
	// Validate KitchenNotes
	
	
	
	// Validate DisplayConfig
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'new',
	
	
	
	// Validate 'completed',
	
	
	
	// Validate 'normal',
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateKitchenTicketsRequest) error {
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
	
	
	// Validate TicketNumber if provided
	
	if req.TicketNumber != nil {
		if err := v.validateTicketNumber(*req.TicketNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate DisplaySequence if provided
	
	
	
	// Validate OrderId if provided
	
	
	if req.OrderId != nil {
		if err := v.validateOrderIdExists(ctx, tx, *req.OrderId); err != nil {
			return err
		}
	}
	
	
	// Validate KitchenStationId if provided
	
	
	if req.KitchenStationId != nil {
		if err := v.validateKitchenStationIdExists(ctx, tx, *req.KitchenStationId); err != nil {
			return err
		}
	}
	
	
	// Validate CourseId if provided
	
	
	if req.CourseId != nil {
		if err := v.validateCourseIdExists(ctx, tx, *req.CourseId); err != nil {
			return err
		}
	}
	
	
	// Validate TicketType if provided
	
	
	
	// Validate Priority if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate FiredAt if provided
	
	
	
	// Validate AcknowledgedAt if provided
	
	
	
	// Validate StartedAt if provided
	
	
	
	// Validate ReadyAt if provided
	
	
	
	// Validate BumpedAt if provided
	
	
	
	// Validate CompletedAt if provided
	
	
	
	// Validate PrepTimeMinutes if provided
	
	
	
	// Validate TargetPrepTime if provided
	
	
	
	// Validate TableNumber if provided
	
	
	
	// Validate OrderType if provided
	
	
	
	// Validate Covers if provided
	
	
	
	// Validate WaiterName if provided
	
	
	
	// Validate SpecialInstructions if provided
	
	
	
	// Validate KitchenNotes if provided
	
	
	
	// Validate DisplayConfig if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'new', if provided
	
	
	
	// Validate 'completed', if provided
	
	
	
	// Validate 'normal', if provided
	
	
	

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



// validateTicketNumber validates ticket_number field
func (v *Validator) validateTicketNumber(value string) error {
	
	// Add custom validation for ticket_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("ticket_number cannot be empty")
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



































































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateKitchenTicketsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateKitchenTicketsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *KitchenTickets, req *dto.UpdateKitchenTicketsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *KitchenTickets) error {
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

package reservation

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

// Validator handles Reservations validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Reservations validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateReservationsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate TableId
	
	
	if err := v.validateTableIdExists(ctx, tx, req.TableId); err != nil {
		return err
	}
	
	
	// Validate CustomerId
	
	
	if err := v.validateCustomerIdExists(ctx, tx, req.CustomerId); err != nil {
		return err
	}
	
	
	// Validate ReservationNumber
	
	if err := v.validateReservationNumber(req.ReservationNumber); err != nil {
		return err
	}
	
	
	
	// Validate ReservationDate
	
	
	
	// Validate ReservationTime
	
	if err := v.validateReservationTime(req.ReservationTime); err != nil {
		return err
	}
	
	
	
	// Validate DurationMinutes
	
	
	
	// Validate PartySize
	
	
	
	// Validate CustomerName
	
	if err := v.validateCustomerName(req.CustomerName); err != nil {
		return err
	}
	
	
	
	// Validate CustomerPhone
	
	
	
	// Validate CustomerEmail
	
	
	
	// Validate Status
	
	
	
	// Validate AssignedWaiterId
	
	
	if err := v.validateAssignedWaiterIdExists(ctx, tx, req.AssignedWaiterId); err != nil {
		return err
	}
	
	
	// Validate AssignedAt
	
	
	
	// Validate SeatedAt
	
	
	
	// Validate CompletedAt
	
	
	
	// Validate SpecialRequests
	
	
	
	// Validate Occasion
	
	
	
	// Validate DietaryRestrictions
	
	
	
	// Validate ConfirmationCode
	
	
	
	// Validate ConfirmedAt
	
	
	
	// Validate ConfirmedBy
	
	
	
	// Validate ReminderSentAt
	
	
	
	// Validate NotificationPreferences
	
	
	
	// Validate CancelledAt
	
	
	
	// Validate CancelledBy
	
	
	
	// Validate CancellationReason
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateReservationsRequest) error {
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
	
	
	// Validate TableId if provided
	
	
	if req.TableId != nil {
		if err := v.validateTableIdExists(ctx, tx, *req.TableId); err != nil {
			return err
		}
	}
	
	
	// Validate CustomerId if provided
	
	
	if req.CustomerId != nil {
		if err := v.validateCustomerIdExists(ctx, tx, *req.CustomerId); err != nil {
			return err
		}
	}
	
	
	// Validate ReservationNumber if provided
	
	if req.ReservationNumber != nil {
		if err := v.validateReservationNumber(*req.ReservationNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate ReservationDate if provided
	
	
	
	// Validate ReservationTime if provided
	
	if req.ReservationTime != nil {
		if err := v.validateReservationTime(*req.ReservationTime); err != nil {
			return err
		}
	}
	
	
	
	// Validate DurationMinutes if provided
	
	
	
	// Validate PartySize if provided
	
	
	
	// Validate CustomerName if provided
	
	if req.CustomerName != nil {
		if err := v.validateCustomerName(*req.CustomerName); err != nil {
			return err
		}
	}
	
	
	
	// Validate CustomerPhone if provided
	
	
	
	// Validate CustomerEmail if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate AssignedWaiterId if provided
	
	
	if req.AssignedWaiterId != nil {
		if err := v.validateAssignedWaiterIdExists(ctx, tx, *req.AssignedWaiterId); err != nil {
			return err
		}
	}
	
	
	// Validate AssignedAt if provided
	
	
	
	// Validate SeatedAt if provided
	
	
	
	// Validate CompletedAt if provided
	
	
	
	// Validate SpecialRequests if provided
	
	
	
	// Validate Occasion if provided
	
	
	
	// Validate DietaryRestrictions if provided
	
	
	
	// Validate ConfirmationCode if provided
	
	
	
	// Validate ConfirmedAt if provided
	
	
	
	// Validate ConfirmedBy if provided
	
	
	
	// Validate ReminderSentAt if provided
	
	
	
	// Validate NotificationPreferences if provided
	
	
	
	// Validate CancelledAt if provided
	
	
	
	// Validate CancelledBy if provided
	
	
	
	// Validate CancellationReason if provided
	
	
	
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





// validateTableIdExists validates that table_id exists
func (v *Validator) validateTableIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for table
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM table WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("table with id %s does not exist", id)
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



// validateReservationNumber validates reservation_number field
func (v *Validator) validateReservationNumber(value string) error {
	
	// Add custom validation for reservation_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("reservation_number cannot be empty")
	}
	
	return nil
}









// validateReservationTime validates reservation_time field
func (v *Validator) validateReservationTime(value string) error {
	
	// Add custom validation for reservation_time
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("reservation_time cannot be empty")
	}
	
	return nil
}













// validateCustomerName validates customer_name field
func (v *Validator) validateCustomerName(value string) error {
	
	// Add custom validation for customer_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("customer_name cannot be empty")
	}
	
	return nil
}



















// validateAssignedWaiterIdExists validates that assigned_waiter_id exists
func (v *Validator) validateAssignedWaiterIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for assigned_waiter
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM assigned_waiter WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check assigned_waiter existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("assigned_waiter with id %s does not exist", id)
	}
	return nil
}











































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateReservationsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateReservationsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Reservations, req *dto.UpdateReservationsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Reservations) error {
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

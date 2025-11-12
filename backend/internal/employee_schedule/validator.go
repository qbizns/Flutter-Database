package employee_schedule

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

// Validator handles EmployeeSchedules validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new EmployeeSchedules validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateEmployeeSchedulesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate EmployeeId
	
	
	if err := v.validateEmployeeIdExists(ctx, tx, req.EmployeeId); err != nil {
		return err
	}
	
	
	// Validate ScheduleDate
	
	
	
	// Validate ShiftType
	
	
	
	// Validate Position
	
	
	
	// Validate ScheduledStartTime
	
	if err := v.validateScheduledStartTime(req.ScheduledStartTime); err != nil {
		return err
	}
	
	
	
	// Validate ScheduledEndTime
	
	if err := v.validateScheduledEndTime(req.ScheduledEndTime); err != nil {
		return err
	}
	
	
	
	// Validate BreakDurationMinutes
	
	
	
	// Validate Status
	
	
	
	// Validate RequiresApproval
	
	
	
	// Validate ApprovedBy
	
	
	
	// Validate ApprovedAt
	
	
	
	// Validate Notes
	
	
	
	// Validate CancellationReason
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'scheduled',
	
	
	
	// Validate 'regular',
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateEmployeeSchedulesRequest) error {
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
	
	
	// Validate EmployeeId if provided
	
	
	if req.EmployeeId != nil {
		if err := v.validateEmployeeIdExists(ctx, tx, *req.EmployeeId); err != nil {
			return err
		}
	}
	
	
	// Validate ScheduleDate if provided
	
	
	
	// Validate ShiftType if provided
	
	
	
	// Validate Position if provided
	
	
	
	// Validate ScheduledStartTime if provided
	
	if req.ScheduledStartTime != nil {
		if err := v.validateScheduledStartTime(*req.ScheduledStartTime); err != nil {
			return err
		}
	}
	
	
	
	// Validate ScheduledEndTime if provided
	
	if req.ScheduledEndTime != nil {
		if err := v.validateScheduledEndTime(*req.ScheduledEndTime); err != nil {
			return err
		}
	}
	
	
	
	// Validate BreakDurationMinutes if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate RequiresApproval if provided
	
	
	
	// Validate ApprovedBy if provided
	
	
	
	// Validate ApprovedAt if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate CancellationReason if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'scheduled', if provided
	
	
	
	// Validate 'regular', if provided
	
	
	

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





// validateEmployeeIdExists validates that employee_id exists
func (v *Validator) validateEmployeeIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for employee
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM employee WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check employee existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("employee with id %s does not exist", id)
	}
	return nil
}















// validateScheduledStartTime validates scheduled_start_time field
func (v *Validator) validateScheduledStartTime(value string) error {
	
	// Add custom validation for scheduled_start_time
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("scheduled_start_time cannot be empty")
	}
	
	return nil
}





// validateScheduledEndTime validates scheduled_end_time field
func (v *Validator) validateScheduledEndTime(value string) error {
	
	// Add custom validation for scheduled_end_time
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("scheduled_end_time cannot be empty")
	}
	
	return nil
}





















































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateEmployeeSchedulesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateEmployeeSchedulesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *EmployeeSchedules, req *dto.UpdateEmployeeSchedulesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *EmployeeSchedules) error {
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

package time_clock_entry

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

// Validator handles TimeClockEntries validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new TimeClockEntries validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateTimeClockEntriesRequest) error {
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
	
	
	// Validate ScheduleId
	
	
	if err := v.validateScheduleIdExists(ctx, tx, req.ScheduleId); err != nil {
		return err
	}
	
	
	// Validate EntryType
	
	if err := v.validateEntryType(req.EntryType); err != nil {
		return err
	}
	
	
	
	// Validate EntryTimestamp
	
	
	
	// Validate ScheduledTimestamp
	
	
	
	// Validate DeviceId
	
	
	if err := v.validateDeviceIdExists(ctx, tx, req.DeviceId); err != nil {
		return err
	}
	
	
	// Validate GpsLocation
	
	
	
	// Validate IpAddress
	
	
	
	// Validate IsLate
	
	
	
	// Validate IsEarly
	
	
	
	// Validate VarianceMinutes
	
	
	
	// Validate RequiresApproval
	
	
	
	// Validate ApprovedBy
	
	
	
	// Validate ApprovedAt
	
	
	
	// Validate IsManualEntry
	
	
	
	// Validate CorrectionNotes
	
	
	
	// Validate PhotoUrl
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'clockIn',
	
	
	
	// Validate 'mealStart',
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateTimeClockEntriesRequest) error {
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
	
	
	// Validate ScheduleId if provided
	
	
	if req.ScheduleId != nil {
		if err := v.validateScheduleIdExists(ctx, tx, *req.ScheduleId); err != nil {
			return err
		}
	}
	
	
	// Validate EntryType if provided
	
	if req.EntryType != nil {
		if err := v.validateEntryType(*req.EntryType); err != nil {
			return err
		}
	}
	
	
	
	// Validate EntryTimestamp if provided
	
	
	
	// Validate ScheduledTimestamp if provided
	
	
	
	// Validate DeviceId if provided
	
	
	if req.DeviceId != nil {
		if err := v.validateDeviceIdExists(ctx, tx, *req.DeviceId); err != nil {
			return err
		}
	}
	
	
	// Validate GpsLocation if provided
	
	
	
	// Validate IpAddress if provided
	
	
	
	// Validate IsLate if provided
	
	
	
	// Validate IsEarly if provided
	
	
	
	// Validate VarianceMinutes if provided
	
	
	
	// Validate RequiresApproval if provided
	
	
	
	// Validate ApprovedBy if provided
	
	
	
	// Validate ApprovedAt if provided
	
	
	
	// Validate IsManualEntry if provided
	
	
	
	// Validate CorrectionNotes if provided
	
	
	
	// Validate PhotoUrl if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'clockIn', if provided
	
	
	
	// Validate 'mealStart', if provided
	
	
	

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





// validateScheduleIdExists validates that schedule_id exists
func (v *Validator) validateScheduleIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for schedule
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schedule WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check schedule existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("schedule with id %s does not exist", id)
	}
	return nil
}



// validateEntryType validates entry_type field
func (v *Validator) validateEntryType(value string) error {
	
	// Add custom validation for entry_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("entry_type cannot be empty")
	}
	
	return nil
}















// validateDeviceIdExists validates that device_id exists
func (v *Validator) validateDeviceIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for device
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM device WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check device existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("device with id %s does not exist", id)
	}
	return nil
}







































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateTimeClockEntriesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateTimeClockEntriesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *TimeClockEntries, req *dto.UpdateTimeClockEntriesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *TimeClockEntries) error {
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

package device

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles Devices validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Devices validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateDevicesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate DeviceCode
	
	if err := v.validateDeviceCode(req.DeviceCode); err != nil {
		return err
	}
	
	
	
	// Validate DeviceName
	
	if err := v.validateDeviceName(req.DeviceName); err != nil {
		return err
	}
	
	
	
	// Validate DeviceType
	
	if err := v.validateDeviceType(req.DeviceType); err != nil {
		return err
	}
	
	
	
	// Validate Manufacturer
	
	
	
	// Validate Model
	
	
	
	// Validate SerialNumber
	
	
	
	// Validate MacAddress
	
	
	
	// Validate IpAddress
	
	
	
	// Validate DeviceConfig
	
	
	
	// Validate ScreenResolution
	
	
	
	// Validate OsVersion
	
	
	
	// Validate ConnectionType
	
	
	
	// Validate ConnectionString
	
	
	
	// Validate Status
	
	
	
	// Validate LastOnlineAt
	
	
	
	// Validate LastHeartbeatAt
	
	
	
	// Validate AssignedToUserId
	
	
	if err := v.validateAssignedToUserIdExists(ctx, tx, req.AssignedToUserId); err != nil {
		return err
	}
	
	
	// Validate AssignedToStationId
	
	
	if err := v.validateAssignedToStationIdExists(ctx, tx, req.AssignedToStationId); err != nil {
		return err
	}
	
	
	// Validate PurchaseDate
	
	
	
	// Validate WarrantyExpiryDate
	
	
	
	// Validate LicenseKey
	
	
	
	// Validate LicenseExpiryDate
	
	
	
	// Validate InstallationNotes
	
	
	
	// Validate MaintenanceNotes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'active',
	
	
	
	// Validate 'posTerminal',
	
	
	
	// Validate 'printer',
	
	
	
	// Validate 'kitchenPrinter',
	
	
	
	// Validate 'network',
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateDevicesRequest) error {
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
	
	
	// Validate DeviceCode if provided
	
	if req.DeviceCode != nil {
		if err := v.validateDeviceCode(*req.DeviceCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate DeviceName if provided
	
	if req.DeviceName != nil {
		if err := v.validateDeviceName(*req.DeviceName); err != nil {
			return err
		}
	}
	
	
	
	// Validate DeviceType if provided
	
	if req.DeviceType != nil {
		if err := v.validateDeviceType(*req.DeviceType); err != nil {
			return err
		}
	}
	
	
	
	// Validate Manufacturer if provided
	
	
	
	// Validate Model if provided
	
	
	
	// Validate SerialNumber if provided
	
	
	
	// Validate MacAddress if provided
	
	
	
	// Validate IpAddress if provided
	
	
	
	// Validate DeviceConfig if provided
	
	
	
	// Validate ScreenResolution if provided
	
	
	
	// Validate OsVersion if provided
	
	
	
	// Validate ConnectionType if provided
	
	
	
	// Validate ConnectionString if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate LastOnlineAt if provided
	
	
	
	// Validate LastHeartbeatAt if provided
	
	
	
	// Validate AssignedToUserId if provided
	
	
	if req.AssignedToUserId != nil {
		if err := v.validateAssignedToUserIdExists(ctx, tx, *req.AssignedToUserId); err != nil {
			return err
		}
	}
	
	
	// Validate AssignedToStationId if provided
	
	
	if req.AssignedToStationId != nil {
		if err := v.validateAssignedToStationIdExists(ctx, tx, *req.AssignedToStationId); err != nil {
			return err
		}
	}
	
	
	// Validate PurchaseDate if provided
	
	
	
	// Validate WarrantyExpiryDate if provided
	
	
	
	// Validate LicenseKey if provided
	
	
	
	// Validate LicenseExpiryDate if provided
	
	
	
	// Validate InstallationNotes if provided
	
	
	
	// Validate MaintenanceNotes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'active', if provided
	
	
	
	// Validate 'posTerminal', if provided
	
	
	
	// Validate 'printer', if provided
	
	
	
	// Validate 'kitchenPrinter', if provided
	
	
	
	// Validate 'network', if provided
	
	
	

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



// validateDeviceCode validates device_code field
func (v *Validator) validateDeviceCode(value string) error {
	
	// Add custom validation for device_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("device_code cannot be empty")
	}
	
	return nil
}





// validateDeviceName validates device_name field
func (v *Validator) validateDeviceName(value string) error {
	
	// Add custom validation for device_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("device_name cannot be empty")
	}
	
	return nil
}





// validateDeviceType validates device_type field
func (v *Validator) validateDeviceType(value string) error {
	
	// Add custom validation for device_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("device_type cannot be empty")
	}
	
	return nil
}



























































// validateAssignedToUserIdExists validates that assigned_to_user_id exists
func (v *Validator) validateAssignedToUserIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for assigned_to_user
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM assigned_to_user WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check assigned_to_user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("assigned_to_user with id %s does not exist", id)
	}
	return nil
}





// validateAssignedToStationIdExists validates that assigned_to_station_id exists
func (v *Validator) validateAssignedToStationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for assigned_to_station
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM assigned_to_station WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check assigned_to_station existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("assigned_to_station with id %s does not exist", id)
	}
	return nil
}



























































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateDevicesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateDevicesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Devices, req *UpdateDevicesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Devices) error {
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

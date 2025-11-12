package delivery_driver

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

// Validator handles DeliveryDrivers validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new DeliveryDrivers validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateDeliveryDriversRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate UserId
	
	
	if err := v.validateUserIdExists(ctx, tx, req.UserId); err != nil {
		return err
	}
	
	
	// Validate DriverCode
	
	if err := v.validateDriverCode(req.DriverCode); err != nil {
		return err
	}
	
	
	
	// Validate FullName
	
	if err := v.validateFullName(req.FullName); err != nil {
		return err
	}
	
	
	
	// Validate Phone
	
	
	
	// Validate Email
	
	
	
	// Validate EmergencyContactName
	
	
	
	// Validate EmergencyContactPhone
	
	
	
	// Validate VehicleType
	
	
	
	// Validate VehicleMake
	
	
	
	// Validate VehicleModel
	
	
	
	// Validate VehicleYear
	
	
	
	// Validate VehicleColor
	
	
	
	// Validate LicensePlate
	
	
	
	// Validate DriversLicenseNumber
	
	
	
	// Validate LicenseExpiryDate
	
	
	
	// Validate InsurancePolicyNumber
	
	
	
	// Validate InsuranceExpiryDate
	
	
	
	// Validate HireDate
	
	
	
	// Validate EmploymentType
	
	
	
	// Validate Status
	
	
	
	// Validate TotalDeliveries
	
	
	
	// Validate SuccessfulDeliveries
	
	
	
	// Validate Rating
	
	
	
	// Validate RatingCount
	
	
	
	// Validate CurrentLocation
	
	
	
	// Validate IsAvailable
	
	
	
	// Validate LastLocationUpdate
	
	
	
	// Validate CommissionRate
	
	
	
	// Validate PaymentMethod
	
	
	
	// Validate Documents
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'active',
	
	
	
	// Validate 'fullTime',
	
	
	
	// Validate 'bike',
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateDeliveryDriversRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate UserId if provided
	
	
	if req.UserId != nil {
		if err := v.validateUserIdExists(ctx, tx, *req.UserId); err != nil {
			return err
		}
	}
	
	
	// Validate DriverCode if provided
	
	if req.DriverCode != nil {
		if err := v.validateDriverCode(*req.DriverCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate FullName if provided
	
	if req.FullName != nil {
		if err := v.validateFullName(*req.FullName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Phone if provided
	
	
	
	// Validate Email if provided
	
	
	
	// Validate EmergencyContactName if provided
	
	
	
	// Validate EmergencyContactPhone if provided
	
	
	
	// Validate VehicleType if provided
	
	
	
	// Validate VehicleMake if provided
	
	
	
	// Validate VehicleModel if provided
	
	
	
	// Validate VehicleYear if provided
	
	
	
	// Validate VehicleColor if provided
	
	
	
	// Validate LicensePlate if provided
	
	
	
	// Validate DriversLicenseNumber if provided
	
	
	
	// Validate LicenseExpiryDate if provided
	
	
	
	// Validate InsurancePolicyNumber if provided
	
	
	
	// Validate InsuranceExpiryDate if provided
	
	
	
	// Validate HireDate if provided
	
	
	
	// Validate EmploymentType if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate TotalDeliveries if provided
	
	
	
	// Validate SuccessfulDeliveries if provided
	
	
	
	// Validate Rating if provided
	
	
	
	// Validate RatingCount if provided
	
	
	
	// Validate CurrentLocation if provided
	
	
	
	// Validate IsAvailable if provided
	
	
	
	// Validate LastLocationUpdate if provided
	
	
	
	// Validate CommissionRate if provided
	
	
	
	// Validate PaymentMethod if provided
	
	
	
	// Validate Documents if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'active', if provided
	
	
	
	// Validate 'fullTime', if provided
	
	
	
	// Validate 'bike', if provided
	
	
	

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





// validateUserIdExists validates that user_id exists
func (v *Validator) validateUserIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for user
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM user WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %s does not exist", id)
	}
	return nil
}



// validateDriverCode validates driver_code field
func (v *Validator) validateDriverCode(value string) error {
	
	// Add custom validation for driver_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("driver_code cannot be empty")
	}
	
	return nil
}





// validateFullName validates full_name field
func (v *Validator) validateFullName(value string) error {
	
	// Add custom validation for full_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("full_name cannot be empty")
	}
	
	return nil
}









































































































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateDeliveryDriversRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateDeliveryDriversRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *DeliveryDrivers, req *dto.UpdateDeliveryDriversRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *DeliveryDrivers) error {
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

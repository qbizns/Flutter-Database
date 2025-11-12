package location

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles Locations validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Locations validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateLocationsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationCode
	
	if err := v.validateLocationCode(req.LocationCode); err != nil {
		return err
	}
	
	
	
	// Validate Name
	
	if err := v.validateName(req.Name); err != nil {
		return err
	}
	
	
	
	// Validate LocationType
	
	
	
	// Validate Phone
	
	
	
	// Validate Email
	
	
	
	// Validate ManagerUserId
	
	
	if err := v.validateManagerUserIdExists(ctx, tx, req.ManagerUserId); err != nil {
		return err
	}
	
	
	// Validate AddressLine1
	
	
	
	// Validate AddressLine2
	
	
	
	// Validate City
	
	
	
	// Validate State
	
	
	
	// Validate Country
	
	
	
	// Validate PostalCode
	
	
	
	// Validate Timezone
	
	
	
	// Validate BusinessHours
	
	
	
	// Validate IsActive
	
	
	
	// Validate IsPrimary
	
	
	
	// Validate AllowSales
	
	
	
	// Validate AllowPurchases
	
	
	
	// Validate TaxRate
	
	
	
	// Validate Notes
	
	
	
	// Validate Settings
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateLocationsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate LocationCode if provided
	
	if req.LocationCode != nil {
		if err := v.validateLocationCode(*req.LocationCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate Name if provided
	
	if req.Name != nil {
		if err := v.validateName(*req.Name); err != nil {
			return err
		}
	}
	
	
	
	// Validate LocationType if provided
	
	
	
	// Validate Phone if provided
	
	
	
	// Validate Email if provided
	
	
	
	// Validate ManagerUserId if provided
	
	
	if req.ManagerUserId != nil {
		if err := v.validateManagerUserIdExists(ctx, tx, *req.ManagerUserId); err != nil {
			return err
		}
	}
	
	
	// Validate AddressLine1 if provided
	
	
	
	// Validate AddressLine2 if provided
	
	
	
	// Validate City if provided
	
	
	
	// Validate State if provided
	
	
	
	// Validate Country if provided
	
	
	
	// Validate PostalCode if provided
	
	
	
	// Validate Timezone if provided
	
	
	
	// Validate BusinessHours if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate IsPrimary if provided
	
	
	
	// Validate AllowSales if provided
	
	
	
	// Validate AllowPurchases if provided
	
	
	
	// Validate TaxRate if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Settings if provided
	
	
	
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



// validateLocationCode validates location_code field
func (v *Validator) validateLocationCode(value string) error {
	
	// Add custom validation for location_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("location_code cannot be empty")
	}
	
	return nil
}





// validateName validates name field
func (v *Validator) validateName(value string) error {
	
	// Add custom validation for name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	
	return nil
}



















// validateManagerUserIdExists validates that manager_user_id exists
func (v *Validator) validateManagerUserIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for manager_user
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM manager_user WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check manager_user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("manager_user with id %s does not exist", id)
	}
	return nil
}











































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateLocationsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateLocationsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Locations, req *UpdateLocationsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Locations) error {
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

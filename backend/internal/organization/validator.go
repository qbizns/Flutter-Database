package organization

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles Organizations validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Organizations validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateOrganizationsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate Name
	
	if err := v.validateName(req.Name); err != nil {
		return err
	}
	
	
	
	// Validate Slug
	
	if err := v.validateSlug(req.Slug); err != nil {
		return err
	}
	
	
	
	// Validate Description
	
	
	
	// Validate Email
	
	
	
	// Validate Phone
	
	
	
	// Validate Address
	
	
	
	// Validate City
	
	
	
	// Validate State
	
	
	
	// Validate Country
	
	
	
	// Validate PostalCode
	
	
	
	// Validate Status
	
	if err := v.validateStatus(req.Status); err != nil {
		return err
	}
	
	
	
	// Validate Plan
	
	
	
	// Validate TrialEndsAt
	
	
	
	// Validate SubscriptionStartsAt
	
	
	
	// Validate SubscriptionEndsAt
	
	
	
	// Validate MaxUsers
	
	
	
	// Validate MaxProducts
	
	
	
	// Validate MaxLocations
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateOrganizationsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate Name if provided
	
	if req.Name != nil {
		if err := v.validateName(*req.Name); err != nil {
			return err
		}
	}
	
	
	
	// Validate Slug if provided
	
	if req.Slug != nil {
		if err := v.validateSlug(*req.Slug); err != nil {
			return err
		}
	}
	
	
	
	// Validate Description if provided
	
	
	
	// Validate Email if provided
	
	
	
	// Validate Phone if provided
	
	
	
	// Validate Address if provided
	
	
	
	// Validate City if provided
	
	
	
	// Validate State if provided
	
	
	
	// Validate Country if provided
	
	
	
	// Validate PostalCode if provided
	
	
	
	// Validate Status if provided
	
	if req.Status != nil {
		if err := v.validateStatus(*req.Status); err != nil {
			return err
		}
	}
	
	
	
	// Validate Plan if provided
	
	
	
	// Validate TrialEndsAt if provided
	
	
	
	// Validate SubscriptionStartsAt if provided
	
	
	
	// Validate SubscriptionEndsAt if provided
	
	
	
	// Validate MaxUsers if provided
	
	
	
	// Validate MaxProducts if provided
	
	
	
	// Validate MaxLocations if provided
	
	
	
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



// validateName validates name field
func (v *Validator) validateName(value string) error {
	
	// Add custom validation for name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	
	return nil
}





// validateSlug validates slug field
func (v *Validator) validateSlug(value string) error {
	
	// Add custom validation for slug
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("slug cannot be empty")
	}
	
	return nil
}





































// validateStatus validates status field
func (v *Validator) validateStatus(value string) error {
	
	// Add custom validation for status
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("status cannot be empty")
	}
	
	return nil
}

















































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateOrganizationsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateOrganizationsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Organizations, req *UpdateOrganizationsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Organizations) error {
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

package user

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

// Validator handles Users validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Users validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateUsersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate Email
	
	if err := v.validateEmail(req.Email); err != nil {
		return err
	}
	
	
	
	// Validate PasswordHash
	
	
	
	// Validate FirstName
	
	if err := v.validateFirstName(req.FirstName); err != nil {
		return err
	}
	
	
	
	// Validate LastName
	
	if err := v.validateLastName(req.LastName); err != nil {
		return err
	}
	
	
	
	// Validate Phone
	
	
	
	// Validate AvatarUrl
	
	
	
	// Validate Status
	
	if err := v.validateStatus(req.Status); err != nil {
		return err
	}
	
	
	
	// Validate EmailVerified
	
	
	
	// Validate EmailVerifiedAt
	
	
	
	// Validate LastLoginAt
	
	
	
	// Validate LastLoginIp
	
	
	
	// Validate FailedLoginAttempts
	
	
	
	// Validate LockedUntil
	
	
	
	// Validate TwoFactorEnabled
	
	
	
	// Validate TwoFactorSecret
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateUsersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate Email if provided
	
	if req.Email != nil {
		if err := v.validateEmail(*req.Email); err != nil {
			return err
		}
	}
	
	
	
	// Validate PasswordHash if provided
	
	
	
	// Validate FirstName if provided
	
	if req.FirstName != nil {
		if err := v.validateFirstName(*req.FirstName); err != nil {
			return err
		}
	}
	
	
	
	// Validate LastName if provided
	
	if req.LastName != nil {
		if err := v.validateLastName(*req.LastName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Phone if provided
	
	
	
	// Validate AvatarUrl if provided
	
	
	
	// Validate Status if provided
	
	if req.Status != nil {
		if err := v.validateStatus(*req.Status); err != nil {
			return err
		}
	}
	
	
	
	// Validate EmailVerified if provided
	
	
	
	// Validate EmailVerifiedAt if provided
	
	
	
	// Validate LastLoginAt if provided
	
	
	
	// Validate LastLoginIp if provided
	
	
	
	// Validate FailedLoginAttempts if provided
	
	
	
	// Validate LockedUntil if provided
	
	
	
	// Validate TwoFactorEnabled if provided
	
	
	
	// Validate TwoFactorSecret if provided
	
	
	
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



// validateEmail validates email field
func (v *Validator) validateEmail(value string) error {
	
	if !isValidEmail(value) {
		return fmt.Errorf("invalid email format for email")
	}
	
	return nil
}









// validateFirstName validates first_name field
func (v *Validator) validateFirstName(value string) error {
	
	// Add custom validation for first_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("first_name cannot be empty")
	}
	
	return nil
}





// validateLastName validates last_name field
func (v *Validator) validateLastName(value string) error {
	
	// Add custom validation for last_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("last_name cannot be empty")
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
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateUsersRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateUsersRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Users, req *dto.UpdateUsersRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Users) error {
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

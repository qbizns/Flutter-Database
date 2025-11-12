package api_key

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles ApiKeys validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ApiKeys validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateApiKeysRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate KeyName
	
	if err := v.validateKeyName(req.KeyName); err != nil {
		return err
	}
	
	
	
	// Validate KeyPrefix
	
	if err := v.validateKeyPrefix(req.KeyPrefix); err != nil {
		return err
	}
	
	
	
	// Validate KeyHash
	
	if err := v.validateKeyHash(req.KeyHash); err != nil {
		return err
	}
	
	
	
	// Validate Scopes
	
	
	
	// Validate AllowedIps
	
	
	
	// Validate IsActive
	
	
	
	// Validate LastUsedAt
	
	
	
	// Validate UsageCount
	
	
	
	// Validate RateLimitPerMinute
	
	
	
	// Validate RateLimitPerHour
	
	
	
	// Validate ExpiresAt
	
	
	
	// Validate CreatedBy
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateApiKeysRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate KeyName if provided
	
	if req.KeyName != nil {
		if err := v.validateKeyName(*req.KeyName); err != nil {
			return err
		}
	}
	
	
	
	// Validate KeyPrefix if provided
	
	if req.KeyPrefix != nil {
		if err := v.validateKeyPrefix(*req.KeyPrefix); err != nil {
			return err
		}
	}
	
	
	
	// Validate KeyHash if provided
	
	if req.KeyHash != nil {
		if err := v.validateKeyHash(*req.KeyHash); err != nil {
			return err
		}
	}
	
	
	
	// Validate Scopes if provided
	
	
	
	// Validate AllowedIps if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate LastUsedAt if provided
	
	
	
	// Validate UsageCount if provided
	
	
	
	// Validate RateLimitPerMinute if provided
	
	
	
	// Validate RateLimitPerHour if provided
	
	
	
	// Validate ExpiresAt if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	

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



// validateKeyName validates key_name field
func (v *Validator) validateKeyName(value string) error {
	
	// Add custom validation for key_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("key_name cannot be empty")
	}
	
	return nil
}





// validateKeyPrefix validates key_prefix field
func (v *Validator) validateKeyPrefix(value string) error {
	
	// Add custom validation for key_prefix
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("key_prefix cannot be empty")
	}
	
	return nil
}





// validateKeyHash validates key_hash field
func (v *Validator) validateKeyHash(value string) error {
	
	// Add custom validation for key_hash
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("key_hash cannot be empty")
	}
	
	return nil
}









































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateApiKeysRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateApiKeysRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ApiKeys, req *UpdateApiKeysRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ApiKeys) error {
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

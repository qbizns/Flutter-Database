package webhook

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles Webhooks validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Webhooks validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateWebhooksRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate WebhookName
	
	if err := v.validateWebhookName(req.WebhookName); err != nil {
		return err
	}
	
	
	
	// Validate Url
	
	if err := v.validateUrl(req.Url); err != nil {
		return err
	}
	
	
	
	// Validate Secret
	
	
	
	// Validate Events
	
	if err := v.validateEvents(req.Events); err != nil {
		return err
	}
	
	
	
	// Validate HttpMethod
	
	
	
	// Validate Headers
	
	
	
	// Validate TimeoutSeconds
	
	
	
	// Validate MaxRetries
	
	
	
	// Validate RetryBackoffSeconds
	
	
	
	// Validate IsActive
	
	
	
	// Validate IsVerified
	
	
	
	// Validate TotalDeliveries
	
	
	
	// Validate SuccessfulDeliveries
	
	
	
	// Validate FailedDeliveries
	
	
	
	// Validate LastDeliveryAt
	
	
	
	// Validate LastSuccessAt
	
	
	
	// Validate LastFailureAt
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateWebhooksRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate WebhookName if provided
	
	if req.WebhookName != nil {
		if err := v.validateWebhookName(*req.WebhookName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Url if provided
	
	if req.Url != nil {
		if err := v.validateUrl(*req.Url); err != nil {
			return err
		}
	}
	
	
	
	// Validate Secret if provided
	
	
	
	// Validate Events if provided
	
	if req.Events != nil {
		if err := v.validateEvents(*req.Events); err != nil {
			return err
		}
	}
	
	
	
	// Validate HttpMethod if provided
	
	
	
	// Validate Headers if provided
	
	
	
	// Validate TimeoutSeconds if provided
	
	
	
	// Validate MaxRetries if provided
	
	
	
	// Validate RetryBackoffSeconds if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate IsVerified if provided
	
	
	
	// Validate TotalDeliveries if provided
	
	
	
	// Validate SuccessfulDeliveries if provided
	
	
	
	// Validate FailedDeliveries if provided
	
	
	
	// Validate LastDeliveryAt if provided
	
	
	
	// Validate LastSuccessAt if provided
	
	
	
	// Validate LastFailureAt if provided
	
	
	
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



// validateWebhookName validates webhook_name field
func (v *Validator) validateWebhookName(value string) error {
	
	// Add custom validation for webhook_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("webhook_name cannot be empty")
	}
	
	return nil
}





// validateUrl validates url field
func (v *Validator) validateUrl(value string) error {
	
	if !isValidURL(value) {
		return fmt.Errorf("invalid URL format for url")
	}
	
	return nil
}









// validateEvents validates events field
func (v *Validator) validateEvents(value string) error {
	
	// Add custom validation for events
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("events cannot be empty")
	}
	
	return nil
}





























































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateWebhooksRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateWebhooksRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Webhooks, req *UpdateWebhooksRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Webhooks) error {
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

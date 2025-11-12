package webhook_delivery

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles WebhookDeliveries validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new WebhookDeliveries validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateWebhookDeliveriesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate WebhookId
	
	
	if err := v.validateWebhookIdExists(ctx, tx, req.WebhookId); err != nil {
		return err
	}
	
	
	// Validate EventType
	
	if err := v.validateEventType(req.EventType); err != nil {
		return err
	}
	
	
	
	// Validate EventId
	
	
	if err := v.validateEventIdExists(ctx, tx, req.EventId); err != nil {
		return err
	}
	
	
	// Validate Status
	
	
	
	// Validate Status
	
	
	
	// Validate RequestUrl
	
	if err := v.validateRequestUrl(req.RequestUrl); err != nil {
		return err
	}
	
	
	
	// Validate RequestMethod
	
	if err := v.validateRequestMethod(req.RequestMethod); err != nil {
		return err
	}
	
	
	
	// Validate RequestHeaders
	
	
	
	// Validate RequestBody
	
	
	
	// Validate ResponseStatusCode
	
	
	
	// Validate ResponseHeaders
	
	
	
	// Validate ResponseBody
	
	
	
	// Validate AttemptNumber
	
	
	
	// Validate DurationMs
	
	
	
	// Validate NextRetryAt
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate DeliveredAt
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateWebhookDeliveriesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate WebhookId if provided
	
	
	if req.WebhookId != nil {
		if err := v.validateWebhookIdExists(ctx, tx, *req.WebhookId); err != nil {
			return err
		}
	}
	
	
	// Validate EventType if provided
	
	if req.EventType != nil {
		if err := v.validateEventType(*req.EventType); err != nil {
			return err
		}
	}
	
	
	
	// Validate EventId if provided
	
	
	if req.EventId != nil {
		if err := v.validateEventIdExists(ctx, tx, *req.EventId); err != nil {
			return err
		}
	}
	
	
	// Validate Status if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate RequestUrl if provided
	
	if req.RequestUrl != nil {
		if err := v.validateRequestUrl(*req.RequestUrl); err != nil {
			return err
		}
	}
	
	
	
	// Validate RequestMethod if provided
	
	if req.RequestMethod != nil {
		if err := v.validateRequestMethod(*req.RequestMethod); err != nil {
			return err
		}
	}
	
	
	
	// Validate RequestHeaders if provided
	
	
	
	// Validate RequestBody if provided
	
	
	
	// Validate ResponseStatusCode if provided
	
	
	
	// Validate ResponseHeaders if provided
	
	
	
	// Validate ResponseBody if provided
	
	
	
	// Validate AttemptNumber if provided
	
	
	
	// Validate DurationMs if provided
	
	
	
	// Validate NextRetryAt if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate DeliveredAt if provided
	
	
	

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





// validateWebhookIdExists validates that webhook_id exists
func (v *Validator) validateWebhookIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for webhook
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM webhook WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check webhook existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("webhook with id %s does not exist", id)
	}
	return nil
}



// validateEventType validates event_type field
func (v *Validator) validateEventType(value string) error {
	
	// Add custom validation for event_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("event_type cannot be empty")
	}
	
	return nil
}







// validateEventIdExists validates that event_id exists
func (v *Validator) validateEventIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for event
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM event WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check event existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("event with id %s does not exist", id)
	}
	return nil
}











// validateRequestUrl validates request_url field
func (v *Validator) validateRequestUrl(value string) error {
	
	if !isValidURL(value) {
		return fmt.Errorf("invalid URL format for request_url")
	}
	
	return nil
}





// validateRequestMethod validates request_method field
func (v *Validator) validateRequestMethod(value string) error {
	
	// Add custom validation for request_method
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("request_method cannot be empty")
	}
	
	return nil
}













































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateWebhookDeliveriesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateWebhookDeliveriesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *WebhookDeliveries, req *UpdateWebhookDeliveriesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *WebhookDeliveries) error {
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

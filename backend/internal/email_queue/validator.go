package email_queue

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles EmailQueue validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new EmailQueue validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateEmailQueueRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ToAddresses
	
	if err := v.validateToAddresses(req.ToAddresses); err != nil {
		return err
	}
	
	
	
	// Validate CcAddresses
	
	
	
	// Validate BccAddresses
	
	
	
	// Validate FromAddress
	
	
	
	// Validate ReplyTo
	
	
	
	// Validate Subject
	
	if err := v.validateSubject(req.Subject); err != nil {
		return err
	}
	
	
	
	// Validate BodyHtml
	
	
	
	// Validate BodyText
	
	
	
	// Validate AttachmentIds
	
	
	
	// Validate TemplateName
	
	
	
	// Validate TemplateData
	
	
	
	// Validate Status
	
	
	
	// Validate Status
	
	
	
	// Validate Provider
	
	
	
	// Validate ProviderMessageId
	
	
	if err := v.validateProviderMessageIdExists(ctx, tx, req.ProviderMessageId); err != nil {
		return err
	}
	
	
	// Validate Attempts
	
	
	
	// Validate MaxAttempts
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate Priority
	
	
	
	// Validate ScheduledAt
	
	
	
	// Validate SentAt
	
	
	
	// Validate FailedAt
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateEmailQueueRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate ToAddresses if provided
	
	if req.ToAddresses != nil {
		if err := v.validateToAddresses(*req.ToAddresses); err != nil {
			return err
		}
	}
	
	
	
	// Validate CcAddresses if provided
	
	
	
	// Validate BccAddresses if provided
	
	
	
	// Validate FromAddress if provided
	
	
	
	// Validate ReplyTo if provided
	
	
	
	// Validate Subject if provided
	
	if req.Subject != nil {
		if err := v.validateSubject(*req.Subject); err != nil {
			return err
		}
	}
	
	
	
	// Validate BodyHtml if provided
	
	
	
	// Validate BodyText if provided
	
	
	
	// Validate AttachmentIds if provided
	
	
	
	// Validate TemplateName if provided
	
	
	
	// Validate TemplateData if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate Provider if provided
	
	
	
	// Validate ProviderMessageId if provided
	
	
	if req.ProviderMessageId != nil {
		if err := v.validateProviderMessageIdExists(ctx, tx, *req.ProviderMessageId); err != nil {
			return err
		}
	}
	
	
	// Validate Attempts if provided
	
	
	
	// Validate MaxAttempts if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate Priority if provided
	
	
	
	// Validate ScheduledAt if provided
	
	
	
	// Validate SentAt if provided
	
	
	
	// Validate FailedAt if provided
	
	
	

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



// validateToAddresses validates to_addresses field
func (v *Validator) validateToAddresses(value string) error {
	
	// Add custom validation for to_addresses
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("to_addresses cannot be empty")
	}
	
	return nil
}





















// validateSubject validates subject field
func (v *Validator) validateSubject(value string) error {
	
	// Add custom validation for subject
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("subject cannot be empty")
	}
	
	return nil
}







































// validateProviderMessageIdExists validates that provider_message_id exists
func (v *Validator) validateProviderMessageIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for provider_message
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM provider_message WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check provider_message existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("provider_message with id %s does not exist", id)
	}
	return nil
}































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateEmailQueueRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateEmailQueueRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *EmailQueue, req *UpdateEmailQueueRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *EmailQueue) error {
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

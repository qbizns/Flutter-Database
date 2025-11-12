package notification

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles Notifications validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Notifications validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateNotificationsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate UserId
	
	
	if err := v.validateUserIdExists(ctx, tx, req.UserId); err != nil {
		return err
	}
	
	
	// Validate NotificationType
	
	if err := v.validateNotificationType(req.NotificationType); err != nil {
		return err
	}
	
	
	
	// Validate Category
	
	if err := v.validateCategory(req.Category); err != nil {
		return err
	}
	
	
	
	// Validate Title
	
	if err := v.validateTitle(req.Title); err != nil {
		return err
	}
	
	
	
	// Validate Message
	
	if err := v.validateMessage(req.Message); err != nil {
		return err
	}
	
	
	
	// Validate ActionUrl
	
	
	
	// Validate ActionLabel
	
	
	
	// Validate Channels
	
	
	
	// Validate IsRead
	
	
	
	// Validate ReadAt
	
	
	
	// Validate RelatedEntityType
	
	
	
	// Validate RelatedEntityId
	
	
	if err := v.validateRelatedEntityIdExists(ctx, tx, req.RelatedEntityId); err != nil {
		return err
	}
	
	
	// Validate Priority
	
	
	
	// Validate ExpiresAt
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateNotificationsRequest) error {
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
	
	
	// Validate NotificationType if provided
	
	if req.NotificationType != nil {
		if err := v.validateNotificationType(*req.NotificationType); err != nil {
			return err
		}
	}
	
	
	
	// Validate Category if provided
	
	if req.Category != nil {
		if err := v.validateCategory(*req.Category); err != nil {
			return err
		}
	}
	
	
	
	// Validate Title if provided
	
	if req.Title != nil {
		if err := v.validateTitle(*req.Title); err != nil {
			return err
		}
	}
	
	
	
	// Validate Message if provided
	
	if req.Message != nil {
		if err := v.validateMessage(*req.Message); err != nil {
			return err
		}
	}
	
	
	
	// Validate ActionUrl if provided
	
	
	
	// Validate ActionLabel if provided
	
	
	
	// Validate Channels if provided
	
	
	
	// Validate IsRead if provided
	
	
	
	// Validate ReadAt if provided
	
	
	
	// Validate RelatedEntityType if provided
	
	
	
	// Validate RelatedEntityId if provided
	
	
	if req.RelatedEntityId != nil {
		if err := v.validateRelatedEntityIdExists(ctx, tx, *req.RelatedEntityId); err != nil {
			return err
		}
	}
	
	
	// Validate Priority if provided
	
	
	
	// Validate ExpiresAt if provided
	
	
	

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



// validateNotificationType validates notification_type field
func (v *Validator) validateNotificationType(value string) error {
	
	// Add custom validation for notification_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("notification_type cannot be empty")
	}
	
	return nil
}





// validateCategory validates category field
func (v *Validator) validateCategory(value string) error {
	
	// Add custom validation for category
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("category cannot be empty")
	}
	
	return nil
}





// validateTitle validates title field
func (v *Validator) validateTitle(value string) error {
	
	// Add custom validation for title
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	
	return nil
}





// validateMessage validates message field
func (v *Validator) validateMessage(value string) error {
	
	// Add custom validation for message
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("message cannot be empty")
	}
	
	return nil
}































// validateRelatedEntityIdExists validates that related_entity_id exists
func (v *Validator) validateRelatedEntityIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for related_entity
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM related_entity WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check related_entity existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("related_entity with id %s does not exist", id)
	}
	return nil
}











// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateNotificationsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateNotificationsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Notifications, req *UpdateNotificationsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Notifications) error {
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

package sales_channel

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles SalesChannels validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new SalesChannels validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateSalesChannelsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ChannelCode
	
	if err := v.validateChannelCode(req.ChannelCode); err != nil {
		return err
	}
	
	
	
	// Validate ChannelName
	
	if err := v.validateChannelName(req.ChannelName); err != nil {
		return err
	}
	
	
	
	// Validate ChannelType
	
	if err := v.validateChannelType(req.ChannelType); err != nil {
		return err
	}
	
	
	
	// Validate ChannelType
	
	
	
	// Validate IsActive
	
	
	
	// Validate SyncInventory
	
	
	
	// Validate SyncCustomers
	
	
	
	// Validate ExternalSystemName
	
	
	
	// Validate ApiEndpoint
	
	
	
	// Validate Settings
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateSalesChannelsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate ChannelCode if provided
	
	if req.ChannelCode != nil {
		if err := v.validateChannelCode(*req.ChannelCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate ChannelName if provided
	
	if req.ChannelName != nil {
		if err := v.validateChannelName(*req.ChannelName); err != nil {
			return err
		}
	}
	
	
	
	// Validate ChannelType if provided
	
	if req.ChannelType != nil {
		if err := v.validateChannelType(*req.ChannelType); err != nil {
			return err
		}
	}
	
	
	
	// Validate ChannelType if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate SyncInventory if provided
	
	
	
	// Validate SyncCustomers if provided
	
	
	
	// Validate ExternalSystemName if provided
	
	
	
	// Validate ApiEndpoint if provided
	
	
	
	// Validate Settings if provided
	
	
	
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



// validateChannelCode validates channel_code field
func (v *Validator) validateChannelCode(value string) error {
	
	// Add custom validation for channel_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("channel_code cannot be empty")
	}
	
	return nil
}





// validateChannelName validates channel_name field
func (v *Validator) validateChannelName(value string) error {
	
	// Add custom validation for channel_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("channel_name cannot be empty")
	}
	
	return nil
}





// validateChannelType validates channel_type field
func (v *Validator) validateChannelType(value string) error {
	
	// Add custom validation for channel_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("channel_type cannot be empty")
	}
	
	return nil
}





































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateSalesChannelsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateSalesChannelsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *SalesChannels, req *UpdateSalesChannelsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *SalesChannels) error {
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

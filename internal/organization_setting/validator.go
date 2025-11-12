package organization_setting

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

// Validator handles OrganizationSettings validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new OrganizationSettings validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateOrganizationSettingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate Timezone
	
	
	
	// Validate DateFormat
	
	
	
	// Validate TimeFormat
	
	
	
	// Validate NumberFormat
	
	
	
	// Validate DefaultCurrency
	
	
	
	// Validate DefaultLanguage
	
	
	
	// Validate BusinessType
	
	
	
	// Validate FiscalYearStart
	
	
	
	// Validate AutoPrintReceipts
	
	
	
	// Validate AllowNegativeInventory
	
	
	
	// Validate RequireCustomerForSale
	
	
	
	// Validate EnablePriceOverride
	
	
	
	// Validate AutoPostSales
	
	
	
	// Validate AutoPostPayments
	
	
	
	// Validate PostingFrequency
	
	
	
	// Validate SmtpHost
	
	
	
	// Validate SmtpPort
	
	
	
	// Validate SmtpUsername
	
	
	
	// Validate SmtpUseTls
	
	
	
	// Validate EmailFromAddress
	
	
	
	// Validate EmailFromName
	
	
	
	// Validate EnableEmailNotifications
	
	
	
	// Validate EnableSmsNotifications
	
	
	
	// Validate Require2fa
	
	
	
	// Validate SessionTimeoutMinutes
	
	
	
	// Validate PasswordMinLength
	
	
	
	// Validate PasswordRequireSpecial
	
	
	
	// Validate ApiEnabled
	
	
	
	// Validate ApiRateLimitPerMinute
	
	
	
	// Validate WebhookRetryMaxAttempts
	
	
	
	// Validate Features
	
	
	
	// Validate CustomSettings
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateOrganizationSettingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate Timezone if provided
	
	
	
	// Validate DateFormat if provided
	
	
	
	// Validate TimeFormat if provided
	
	
	
	// Validate NumberFormat if provided
	
	
	
	// Validate DefaultCurrency if provided
	
	
	
	// Validate DefaultLanguage if provided
	
	
	
	// Validate BusinessType if provided
	
	
	
	// Validate FiscalYearStart if provided
	
	
	
	// Validate AutoPrintReceipts if provided
	
	
	
	// Validate AllowNegativeInventory if provided
	
	
	
	// Validate RequireCustomerForSale if provided
	
	
	
	// Validate EnablePriceOverride if provided
	
	
	
	// Validate AutoPostSales if provided
	
	
	
	// Validate AutoPostPayments if provided
	
	
	
	// Validate PostingFrequency if provided
	
	
	
	// Validate SmtpHost if provided
	
	
	
	// Validate SmtpPort if provided
	
	
	
	// Validate SmtpUsername if provided
	
	
	
	// Validate SmtpUseTls if provided
	
	
	
	// Validate EmailFromAddress if provided
	
	
	
	// Validate EmailFromName if provided
	
	
	
	// Validate EnableEmailNotifications if provided
	
	
	
	// Validate EnableSmsNotifications if provided
	
	
	
	// Validate Require2fa if provided
	
	
	
	// Validate SessionTimeoutMinutes if provided
	
	
	
	// Validate PasswordMinLength if provided
	
	
	
	// Validate PasswordRequireSpecial if provided
	
	
	
	// Validate ApiEnabled if provided
	
	
	
	// Validate ApiRateLimitPerMinute if provided
	
	
	
	// Validate WebhookRetryMaxAttempts if provided
	
	
	
	// Validate Features if provided
	
	
	
	// Validate CustomSettings if provided
	
	
	
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







































































































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateOrganizationSettingsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateOrganizationSettingsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *OrganizationSettings, req *dto.UpdateOrganizationSettingsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *OrganizationSettings) error {
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

package customer

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles Customers validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Customers validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateCustomersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate CustomerCode
	
	
	
	// Validate FirstName
	
	
	
	// Validate LastName
	
	
	
	// Validate CompanyName
	
	
	
	// Validate Email
	
	
	
	// Validate Phone
	
	
	
	// Validate AlternatePhone
	
	
	
	// Validate AddressLine1
	
	
	
	// Validate AddressLine2
	
	
	
	// Validate City
	
	
	
	// Validate State
	
	
	
	// Validate Country
	
	
	
	// Validate PostalCode
	
	
	
	// Validate DateOfBirth
	
	
	
	// Validate Gender
	
	
	
	// Validate TaxNumber
	
	
	
	// Validate LoyaltyPoints
	
	
	
	// Validate LoyaltyTier
	
	
	
	// Validate CreditLimit
	
	
	
	// Validate OutstandingBalance
	
	
	
	// Validate TotalPurchases
	
	
	
	// Validate TotalOrders
	
	
	
	// Validate LastPurchaseAt
	
	
	
	// Validate IsActive
	
	
	
	// Validate Notes
	
	
	
	// Validate CustomFields
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateCustomersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate CustomerCode if provided
	
	
	
	// Validate FirstName if provided
	
	
	
	// Validate LastName if provided
	
	
	
	// Validate CompanyName if provided
	
	
	
	// Validate Email if provided
	
	
	
	// Validate Phone if provided
	
	
	
	// Validate AlternatePhone if provided
	
	
	
	// Validate AddressLine1 if provided
	
	
	
	// Validate AddressLine2 if provided
	
	
	
	// Validate City if provided
	
	
	
	// Validate State if provided
	
	
	
	// Validate Country if provided
	
	
	
	// Validate PostalCode if provided
	
	
	
	// Validate DateOfBirth if provided
	
	
	
	// Validate Gender if provided
	
	
	
	// Validate TaxNumber if provided
	
	
	
	// Validate LoyaltyPoints if provided
	
	
	
	// Validate LoyaltyTier if provided
	
	
	
	// Validate CreditLimit if provided
	
	
	
	// Validate OutstandingBalance if provided
	
	
	
	// Validate TotalPurchases if provided
	
	
	
	// Validate TotalOrders if provided
	
	
	
	// Validate LastPurchaseAt if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate CustomFields if provided
	
	
	
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























































































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateCustomersRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateCustomersRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Customers, req *UpdateCustomersRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Customers) error {
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

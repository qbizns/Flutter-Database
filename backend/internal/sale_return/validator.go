package sale_return

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles SaleReturns validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new SaleReturns validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateSaleReturnsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ReturnNumber
	
	if err := v.validateReturnNumber(req.ReturnNumber); err != nil {
		return err
	}
	
	
	
	// Validate OriginalSaleId
	
	
	if err := v.validateOriginalSaleIdExists(ctx, tx, req.OriginalSaleId); err != nil {
		return err
	}
	
	
	// Validate CustomerId
	
	
	if err := v.validateCustomerIdExists(ctx, tx, req.CustomerId); err != nil {
		return err
	}
	
	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate UserId
	
	
	if err := v.validateUserIdExists(ctx, tx, req.UserId); err != nil {
		return err
	}
	
	
	// Validate ReturnDate
	
	
	
	// Validate TotalAmount
	
	
	
	// Validate RefundAmount
	
	
	
	// Validate RestockingFee
	
	
	
	// Validate RefundMethod
	
	
	
	// Validate Status
	
	
	
	// Validate Status
	
	
	
	// Validate ApprovedBy
	
	
	
	// Validate ApprovedAt
	
	
	
	// Validate Notes
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateSaleReturnsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate ReturnNumber if provided
	
	if req.ReturnNumber != nil {
		if err := v.validateReturnNumber(*req.ReturnNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate OriginalSaleId if provided
	
	
	if req.OriginalSaleId != nil {
		if err := v.validateOriginalSaleIdExists(ctx, tx, *req.OriginalSaleId); err != nil {
			return err
		}
	}
	
	
	// Validate CustomerId if provided
	
	
	if req.CustomerId != nil {
		if err := v.validateCustomerIdExists(ctx, tx, *req.CustomerId); err != nil {
			return err
		}
	}
	
	
	// Validate LocationId if provided
	
	
	if req.LocationId != nil {
		if err := v.validateLocationIdExists(ctx, tx, *req.LocationId); err != nil {
			return err
		}
	}
	
	
	// Validate UserId if provided
	
	
	if req.UserId != nil {
		if err := v.validateUserIdExists(ctx, tx, *req.UserId); err != nil {
			return err
		}
	}
	
	
	// Validate ReturnDate if provided
	
	
	
	// Validate TotalAmount if provided
	
	
	
	// Validate RefundAmount if provided
	
	
	
	// Validate RestockingFee if provided
	
	
	
	// Validate RefundMethod if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate ApprovedBy if provided
	
	
	
	// Validate ApprovedAt if provided
	
	
	
	// Validate Notes if provided
	
	
	
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



// validateReturnNumber validates return_number field
func (v *Validator) validateReturnNumber(value string) error {
	
	// Add custom validation for return_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("return_number cannot be empty")
	}
	
	return nil
}







// validateOriginalSaleIdExists validates that original_sale_id exists
func (v *Validator) validateOriginalSaleIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for original_sale
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM original_sale WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check original_sale existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("original_sale with id %s does not exist", id)
	}
	return nil
}





// validateCustomerIdExists validates that customer_id exists
func (v *Validator) validateCustomerIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for customer
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM customer WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check customer existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("customer with id %s does not exist", id)
	}
	return nil
}





// validateLocationIdExists validates that location_id exists
func (v *Validator) validateLocationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for location
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM location WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check location existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("location with id %s does not exist", id)
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















































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateSaleReturnsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateSaleReturnsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *SaleReturns, req *UpdateSaleReturnsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *SaleReturns) error {
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

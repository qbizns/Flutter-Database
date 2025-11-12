package promotion

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

// Validator handles Promotions validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Promotions validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreatePromotionsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate PromotionCode
	
	if err := v.validatePromotionCode(req.PromotionCode); err != nil {
		return err
	}
	
	
	
	// Validate Name
	
	if err := v.validateName(req.Name); err != nil {
		return err
	}
	
	
	
	// Validate Description
	
	
	
	// Validate PromotionType
	
	if err := v.validatePromotionType(req.PromotionType); err != nil {
		return err
	}
	
	
	
	// Validate DiscountValue
	
	
	
	// Validate AppliesTo
	
	
	
	// Validate ApplicableProductIds
	
	
	
	// Validate ApplicableCategoryIds
	
	
	
	// Validate MinimumPurchaseAmount
	
	
	
	// Validate MinimumQuantity
	
	
	
	// Validate BuyQuantity
	
	
	
	// Validate GetQuantity
	
	
	
	// Validate GetDiscountPercentage
	
	
	
	// Validate MaxUsesTotal
	
	
	
	// Validate MaxUsesPerCustomer
	
	
	
	// Validate CurrentUses
	
	
	
	// Validate StartDate
	
	
	
	// Validate EndDate
	
	
	
	// Validate IsActive
	
	
	
	// Validate IsCombinable
	
	
	
	// Validate Priority
	
	
	
	// Validate TermsAndConditions
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdatePromotionsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate PromotionCode if provided
	
	if req.PromotionCode != nil {
		if err := v.validatePromotionCode(*req.PromotionCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate Name if provided
	
	if req.Name != nil {
		if err := v.validateName(*req.Name); err != nil {
			return err
		}
	}
	
	
	
	// Validate Description if provided
	
	
	
	// Validate PromotionType if provided
	
	if req.PromotionType != nil {
		if err := v.validatePromotionType(*req.PromotionType); err != nil {
			return err
		}
	}
	
	
	
	// Validate DiscountValue if provided
	
	
	
	// Validate AppliesTo if provided
	
	
	
	// Validate ApplicableProductIds if provided
	
	
	
	// Validate ApplicableCategoryIds if provided
	
	
	
	// Validate MinimumPurchaseAmount if provided
	
	
	
	// Validate MinimumQuantity if provided
	
	
	
	// Validate BuyQuantity if provided
	
	
	
	// Validate GetQuantity if provided
	
	
	
	// Validate GetDiscountPercentage if provided
	
	
	
	// Validate MaxUsesTotal if provided
	
	
	
	// Validate MaxUsesPerCustomer if provided
	
	
	
	// Validate CurrentUses if provided
	
	
	
	// Validate StartDate if provided
	
	
	
	// Validate EndDate if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate IsCombinable if provided
	
	
	
	// Validate Priority if provided
	
	
	
	// Validate TermsAndConditions if provided
	
	
	
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



// validatePromotionCode validates promotion_code field
func (v *Validator) validatePromotionCode(value string) error {
	
	// Add custom validation for promotion_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("promotion_code cannot be empty")
	}
	
	return nil
}





// validateName validates name field
func (v *Validator) validateName(value string) error {
	
	// Add custom validation for name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	
	return nil
}









// validatePromotionType validates promotion_type field
func (v *Validator) validatePromotionType(value string) error {
	
	// Add custom validation for promotion_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("promotion_type cannot be empty")
	}
	
	return nil
}

























































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreatePromotionsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreatePromotionsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Promotions, req *dto.UpdatePromotionsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Promotions) error {
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

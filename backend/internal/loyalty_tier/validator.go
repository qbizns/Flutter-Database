package loyalty_tier

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles LoyaltyTiers validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new LoyaltyTiers validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyTiersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate TierCode
	
	if err := v.validateTierCode(req.TierCode); err != nil {
		return err
	}
	
	
	
	// Validate TierName
	
	if err := v.validateTierName(req.TierName); err != nil {
		return err
	}
	
	
	
	// Validate TierLevel
	
	
	
	// Validate Description
	
	
	
	// Validate PointsThreshold
	
	
	
	// Validate AnnualSpendThreshold
	
	
	
	// Validate PurchaseCountThreshold
	
	
	
	// Validate PointsMultiplier
	
	
	
	// Validate DiscountPercentage
	
	
	
	// Validate TierColor
	
	
	
	// Validate TierIcon
	
	
	
	// Validate BadgeImageUrl
	
	
	
	// Validate IsActive
	
	
	
	// Validate IsDefault
	
	
	
	// Validate SortOrder
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate PointsThreshold
	
	
	
	// Validate (annualSpendThreshold
	
	
	
	// Validate (purchaseCountThreshold
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateLoyaltyTiersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate TierCode if provided
	
	if req.TierCode != nil {
		if err := v.validateTierCode(*req.TierCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate TierName if provided
	
	if req.TierName != nil {
		if err := v.validateTierName(*req.TierName); err != nil {
			return err
		}
	}
	
	
	
	// Validate TierLevel if provided
	
	
	
	// Validate Description if provided
	
	
	
	// Validate PointsThreshold if provided
	
	
	
	// Validate AnnualSpendThreshold if provided
	
	
	
	// Validate PurchaseCountThreshold if provided
	
	
	
	// Validate PointsMultiplier if provided
	
	
	
	// Validate DiscountPercentage if provided
	
	
	
	// Validate TierColor if provided
	
	
	
	// Validate TierIcon if provided
	
	
	
	// Validate BadgeImageUrl if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate IsDefault if provided
	
	
	
	// Validate SortOrder if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate PointsThreshold if provided
	
	
	
	// Validate (annualSpendThreshold if provided
	
	
	
	// Validate (purchaseCountThreshold if provided
	
	
	

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



// validateTierCode validates tier_code field
func (v *Validator) validateTierCode(value string) error {
	
	// Add custom validation for tier_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("tier_code cannot be empty")
	}
	
	return nil
}





// validateTierName validates tier_name field
func (v *Validator) validateTierName(value string) error {
	
	// Add custom validation for tier_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("tier_name cannot be empty")
	}
	
	return nil
}

















































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyTiersRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyTiersRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *LoyaltyTiers, req *UpdateLoyaltyTiersRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *LoyaltyTiers) error {
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

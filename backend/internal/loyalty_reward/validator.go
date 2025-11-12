package loyalty_reward

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles LoyaltyRewards validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new LoyaltyRewards validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyRewardsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate RewardCode
	
	if err := v.validateRewardCode(req.RewardCode); err != nil {
		return err
	}
	
	
	
	// Validate RewardName
	
	if err := v.validateRewardName(req.RewardName); err != nil {
		return err
	}
	
	
	
	// Validate Description
	
	
	
	// Validate RewardType
	
	if err := v.validateRewardType(req.RewardType); err != nil {
		return err
	}
	
	
	
	// Validate PointsCost
	
	
	
	// Validate RewardValue
	
	
	
	// Validate DiscountPercentage
	
	
	
	// Validate DiscountAmount
	
	
	
	// Validate ProductId
	
	
	if err := v.validateProductIdExists(ctx, tx, req.ProductId); err != nil {
		return err
	}
	
	
	// Validate ProductVariantId
	
	
	if err := v.validateProductVariantIdExists(ctx, tx, req.ProductVariantId); err != nil {
		return err
	}
	
	
	// Validate IsActive
	
	
	
	// Validate AvailableFrom
	
	
	
	// Validate AvailableTo
	
	
	
	// Validate TotalAvailable
	
	
	
	// Validate TotalRedeemed
	
	
	
	// Validate MaxRedemptionsPerCustomer
	
	
	
	// Validate MinimumTierLevel
	
	
	
	// Validate TierIds
	
	
	
	// Validate ImageUrl
	
	
	
	// Validate ThumbnailUrl
	
	
	
	// Validate Featured
	
	
	
	// Validate SortOrder
	
	
	
	// Validate IsFeatured
	
	
	
	// Validate TermsAndConditions
	
	
	
	// Validate RedemptionInstructions
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate AvailableTo
	
	
	
	// Validate TotalRedeemed
	
	
	
	// Validate (totalAvailable
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateLoyaltyRewardsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate RewardCode if provided
	
	if req.RewardCode != nil {
		if err := v.validateRewardCode(*req.RewardCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate RewardName if provided
	
	if req.RewardName != nil {
		if err := v.validateRewardName(*req.RewardName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Description if provided
	
	
	
	// Validate RewardType if provided
	
	if req.RewardType != nil {
		if err := v.validateRewardType(*req.RewardType); err != nil {
			return err
		}
	}
	
	
	
	// Validate PointsCost if provided
	
	
	
	// Validate RewardValue if provided
	
	
	
	// Validate DiscountPercentage if provided
	
	
	
	// Validate DiscountAmount if provided
	
	
	
	// Validate ProductId if provided
	
	
	if req.ProductId != nil {
		if err := v.validateProductIdExists(ctx, tx, *req.ProductId); err != nil {
			return err
		}
	}
	
	
	// Validate ProductVariantId if provided
	
	
	if req.ProductVariantId != nil {
		if err := v.validateProductVariantIdExists(ctx, tx, *req.ProductVariantId); err != nil {
			return err
		}
	}
	
	
	// Validate IsActive if provided
	
	
	
	// Validate AvailableFrom if provided
	
	
	
	// Validate AvailableTo if provided
	
	
	
	// Validate TotalAvailable if provided
	
	
	
	// Validate TotalRedeemed if provided
	
	
	
	// Validate MaxRedemptionsPerCustomer if provided
	
	
	
	// Validate MinimumTierLevel if provided
	
	
	
	// Validate TierIds if provided
	
	
	
	// Validate ImageUrl if provided
	
	
	
	// Validate ThumbnailUrl if provided
	
	
	
	// Validate Featured if provided
	
	
	
	// Validate SortOrder if provided
	
	
	
	// Validate IsFeatured if provided
	
	
	
	// Validate TermsAndConditions if provided
	
	
	
	// Validate RedemptionInstructions if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate AvailableTo if provided
	
	
	
	// Validate TotalRedeemed if provided
	
	
	
	// Validate (totalAvailable if provided
	
	
	

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



// validateRewardCode validates reward_code field
func (v *Validator) validateRewardCode(value string) error {
	
	// Add custom validation for reward_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("reward_code cannot be empty")
	}
	
	return nil
}





// validateRewardName validates reward_name field
func (v *Validator) validateRewardName(value string) error {
	
	// Add custom validation for reward_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("reward_name cannot be empty")
	}
	
	return nil
}









// validateRewardType validates reward_type field
func (v *Validator) validateRewardType(value string) error {
	
	// Add custom validation for reward_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("reward_type cannot be empty")
	}
	
	return nil
}























// validateProductIdExists validates that product_id exists
func (v *Validator) validateProductIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for product
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check product existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("product with id %s does not exist", id)
	}
	return nil
}





// validateProductVariantIdExists validates that product_variant_id exists
func (v *Validator) validateProductVariantIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for product_variant
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM product_variant WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check product_variant existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("product_variant with id %s does not exist", id)
	}
	return nil
}























































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyRewardsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyRewardsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *LoyaltyRewards, req *UpdateLoyaltyRewardsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *LoyaltyRewards) error {
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

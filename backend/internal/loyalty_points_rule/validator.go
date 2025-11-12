package loyalty_points_rule

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles LoyaltyPointsRules validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new LoyaltyPointsRules validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyPointsRulesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate RuleCode
	
	if err := v.validateRuleCode(req.RuleCode); err != nil {
		return err
	}
	
	
	
	// Validate RuleName
	
	if err := v.validateRuleName(req.RuleName); err != nil {
		return err
	}
	
	
	
	// Validate Description
	
	
	
	// Validate RuleType
	
	if err := v.validateRuleType(req.RuleType); err != nil {
		return err
	}
	
	
	
	// Validate PointsPerAmount
	
	
	
	// Validate FixedPoints
	
	
	
	// Validate Multiplier
	
	
	
	// Validate AppliesTo
	
	
	
	// Validate ApplicableProductIds
	
	
	
	// Validate ApplicableCategoryIds
	
	
	
	// Validate ApplicableTierIds
	
	
	
	// Validate MinimumPurchaseAmount
	
	
	
	// Validate MaximumPointsPerTransaction
	
	
	
	// Validate MaximumPointsPerDay
	
	
	
	// Validate MaximumPointsPerMonth
	
	
	
	// Validate StartDate
	
	
	
	// Validate EndDate
	
	
	
	// Validate IsActive
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateLoyaltyPointsRulesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate RuleCode if provided
	
	if req.RuleCode != nil {
		if err := v.validateRuleCode(*req.RuleCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate RuleName if provided
	
	if req.RuleName != nil {
		if err := v.validateRuleName(*req.RuleName); err != nil {
			return err
		}
	}
	
	
	
	// Validate Description if provided
	
	
	
	// Validate RuleType if provided
	
	if req.RuleType != nil {
		if err := v.validateRuleType(*req.RuleType); err != nil {
			return err
		}
	}
	
	
	
	// Validate PointsPerAmount if provided
	
	
	
	// Validate FixedPoints if provided
	
	
	
	// Validate Multiplier if provided
	
	
	
	// Validate AppliesTo if provided
	
	
	
	// Validate ApplicableProductIds if provided
	
	
	
	// Validate ApplicableCategoryIds if provided
	
	
	
	// Validate ApplicableTierIds if provided
	
	
	
	// Validate MinimumPurchaseAmount if provided
	
	
	
	// Validate MaximumPointsPerTransaction if provided
	
	
	
	// Validate MaximumPointsPerDay if provided
	
	
	
	// Validate MaximumPointsPerMonth if provided
	
	
	
	// Validate StartDate if provided
	
	
	
	// Validate EndDate if provided
	
	
	
	// Validate IsActive if provided
	
	
	
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



// validateRuleCode validates rule_code field
func (v *Validator) validateRuleCode(value string) error {
	
	// Add custom validation for rule_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("rule_code cannot be empty")
	}
	
	return nil
}





// validateRuleName validates rule_name field
func (v *Validator) validateRuleName(value string) error {
	
	// Add custom validation for rule_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("rule_name cannot be empty")
	}
	
	return nil
}









// validateRuleType validates rule_type field
func (v *Validator) validateRuleType(value string) error {
	
	// Add custom validation for rule_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("rule_type cannot be empty")
	}
	
	return nil
}

















































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyPointsRulesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateLoyaltyPointsRulesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *LoyaltyPointsRules, req *UpdateLoyaltyPointsRulesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *LoyaltyPointsRules) error {
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

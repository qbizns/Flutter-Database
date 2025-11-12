package pos_account_mapping

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles PosAccountMappings validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new PosAccountMappings validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreatePosAccountMappingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate SourceType
	
	if err := v.validateSourceType(req.SourceType); err != nil {
		return err
	}
	
	
	
	// Validate 'product',
	
	
	
	// Validate 'category',
	
	
	
	// Validate 'paymentMethod',
	
	
	
	// Validate 'salesChannel',
	
	
	
	// Validate 'discount',
	
	
	
	// Validate 'rounding',
	
	
	
	// Validate 'tax',
	
	
	
	// Validate 'serviceCharge',
	
	
	
	// Validate 'shipping',
	
	
	
	// Validate 'giftCard',
	
	
	
	// Validate 'storeCredit',
	
	
	
	// Validate 'loyaltyRedemption',--
	
	
	
	// Validate 'default'
	
	
	
	// Validate SourceId
	
	
	if err := v.validateSourceIdExists(ctx, tx, req.SourceId); err != nil {
		return err
	}
	
	
	// Validate SourceCode
	
	
	
	// Validate Purpose
	
	if err := v.validatePurpose(req.Purpose); err != nil {
		return err
	}
	
	
	
	// Validate 'revenue',
	
	
	
	// Validate 'cogs',
	
	
	
	// Validate 'inventory',
	
	
	
	// Validate 'expense',
	
	
	
	// Validate 'liability',
	
	
	
	// Validate 'asset',
	
	
	
	// Validate 'discountExpense',
	
	
	
	// Validate 'discountContra',
	
	
	
	// Validate 'taxLiability',
	
	
	
	// Validate 'rounding',
	
	
	
	// Validate 'clearing'
	
	
	
	// Validate AccountId
	
	
	if err := v.validateAccountIdExists(ctx, tx, req.AccountId); err != nil {
		return err
	}
	
	
	// Validate IsDefault
	
	
	
	// Validate IsActive
	
	
	
	// Validate Priority
	
	
	
	// Validate Conditions
	
	
	
	// Validate EffectiveFrom
	
	
	
	// Validate EffectiveTo
	
	
	
	// Validate Description
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate (sourceType
	
	
	
	// Validate (sourceType
	
	if err := v.validate(sourceType(req.(sourceType); err != nil {
		return err
	}
	
	
	
	// Validate EffectiveFrom
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdatePosAccountMappingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate SourceType if provided
	
	if req.SourceType != nil {
		if err := v.validateSourceType(*req.SourceType); err != nil {
			return err
		}
	}
	
	
	
	// Validate 'product', if provided
	
	
	
	// Validate 'category', if provided
	
	
	
	// Validate 'paymentMethod', if provided
	
	
	
	// Validate 'salesChannel', if provided
	
	
	
	// Validate 'discount', if provided
	
	
	
	// Validate 'rounding', if provided
	
	
	
	// Validate 'tax', if provided
	
	
	
	// Validate 'serviceCharge', if provided
	
	
	
	// Validate 'shipping', if provided
	
	
	
	// Validate 'giftCard', if provided
	
	
	
	// Validate 'storeCredit', if provided
	
	
	
	// Validate 'loyaltyRedemption',-- if provided
	
	
	
	// Validate 'default' if provided
	
	
	
	// Validate SourceId if provided
	
	
	if req.SourceId != nil {
		if err := v.validateSourceIdExists(ctx, tx, *req.SourceId); err != nil {
			return err
		}
	}
	
	
	// Validate SourceCode if provided
	
	
	
	// Validate Purpose if provided
	
	if req.Purpose != nil {
		if err := v.validatePurpose(*req.Purpose); err != nil {
			return err
		}
	}
	
	
	
	// Validate 'revenue', if provided
	
	
	
	// Validate 'cogs', if provided
	
	
	
	// Validate 'inventory', if provided
	
	
	
	// Validate 'expense', if provided
	
	
	
	// Validate 'liability', if provided
	
	
	
	// Validate 'asset', if provided
	
	
	
	// Validate 'discountExpense', if provided
	
	
	
	// Validate 'discountContra', if provided
	
	
	
	// Validate 'taxLiability', if provided
	
	
	
	// Validate 'rounding', if provided
	
	
	
	// Validate 'clearing' if provided
	
	
	
	// Validate AccountId if provided
	
	
	if req.AccountId != nil {
		if err := v.validateAccountIdExists(ctx, tx, *req.AccountId); err != nil {
			return err
		}
	}
	
	
	// Validate IsDefault if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate Priority if provided
	
	
	
	// Validate Conditions if provided
	
	
	
	// Validate EffectiveFrom if provided
	
	
	
	// Validate EffectiveTo if provided
	
	
	
	// Validate Description if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate (sourceType if provided
	
	
	
	// Validate (sourceType if provided
	
	if req.(sourceType != nil {
		if err := v.validate(sourceType(*req.(sourceType); err != nil {
			return err
		}
	}
	
	
	
	// Validate EffectiveFrom if provided
	
	
	

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



// validateSourceType validates source_type field
func (v *Validator) validateSourceType(value string) error {
	
	// Add custom validation for source_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("source_type cannot be empty")
	}
	
	return nil
}



























































// validateSourceIdExists validates that source_id exists
func (v *Validator) validateSourceIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for source
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM source WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check source existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("source with id %s does not exist", id)
	}
	return nil
}







// validatePurpose validates purpose field
func (v *Validator) validatePurpose(value string) error {
	
	// Add custom validation for purpose
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("purpose cannot be empty")
	}
	
	return nil
}



















































// validateAccountIdExists validates that account_id exists
func (v *Validator) validateAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("account with id %s does not exist", id)
	}
	return nil
}



















































// validate(sourceType validates (source_type field
func (v *Validator) validate(sourceType(value string) error {
	
	// Add custom validation for (source_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("(source_type cannot be empty")
	}
	
	return nil
}









// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreatePosAccountMappingsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreatePosAccountMappingsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *PosAccountMappings, req *UpdatePosAccountMappingsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *PosAccountMappings) error {
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

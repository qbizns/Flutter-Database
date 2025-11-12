package asset_category

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

// Validator handles AssetCategories validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new AssetCategories validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateAssetCategoriesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate CategoryCode
	
	if err := v.validateCategoryCode(req.CategoryCode); err != nil {
		return err
	}
	
	
	
	// Validate CategoryName
	
	if err := v.validateCategoryName(req.CategoryName); err != nil {
		return err
	}
	
	
	
	// Validate DefaultDepreciationMethod
	
	
	
	// Validate DefaultUsefulLifeYears
	
	
	
	// Validate DefaultSalvageValuePercent
	
	
	
	// Validate AssetAccountId
	
	
	if err := v.validateAssetAccountIdExists(ctx, tx, req.AssetAccountId); err != nil {
		return err
	}
	
	
	// Validate AccumulatedDepreciationAccountId
	
	
	if err := v.validateAccumulatedDepreciationAccountIdExists(ctx, tx, req.AccumulatedDepreciationAccountId); err != nil {
		return err
	}
	
	
	// Validate DepreciationExpenseAccountId
	
	
	if err := v.validateDepreciationExpenseAccountIdExists(ctx, tx, req.DepreciationExpenseAccountId); err != nil {
		return err
	}
	
	
	// Validate Description
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateAssetCategoriesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate CategoryCode if provided
	
	if req.CategoryCode != nil {
		if err := v.validateCategoryCode(*req.CategoryCode); err != nil {
			return err
		}
	}
	
	
	
	// Validate CategoryName if provided
	
	if req.CategoryName != nil {
		if err := v.validateCategoryName(*req.CategoryName); err != nil {
			return err
		}
	}
	
	
	
	// Validate DefaultDepreciationMethod if provided
	
	
	
	// Validate DefaultUsefulLifeYears if provided
	
	
	
	// Validate DefaultSalvageValuePercent if provided
	
	
	
	// Validate AssetAccountId if provided
	
	
	if req.AssetAccountId != nil {
		if err := v.validateAssetAccountIdExists(ctx, tx, *req.AssetAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate AccumulatedDepreciationAccountId if provided
	
	
	if req.AccumulatedDepreciationAccountId != nil {
		if err := v.validateAccumulatedDepreciationAccountIdExists(ctx, tx, *req.AccumulatedDepreciationAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate DepreciationExpenseAccountId if provided
	
	
	if req.DepreciationExpenseAccountId != nil {
		if err := v.validateDepreciationExpenseAccountIdExists(ctx, tx, *req.DepreciationExpenseAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate Description if provided
	
	
	

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



// validateCategoryCode validates category_code field
func (v *Validator) validateCategoryCode(value string) error {
	
	// Add custom validation for category_code
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("category_code cannot be empty")
	}
	
	return nil
}





// validateCategoryName validates category_name field
func (v *Validator) validateCategoryName(value string) error {
	
	// Add custom validation for category_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("category_name cannot be empty")
	}
	
	return nil
}



















// validateAssetAccountIdExists validates that asset_account_id exists
func (v *Validator) validateAssetAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for asset_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM asset_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check asset_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("asset_account with id %s does not exist", id)
	}
	return nil
}





// validateAccumulatedDepreciationAccountIdExists validates that accumulated_depreciation_account_id exists
func (v *Validator) validateAccumulatedDepreciationAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for accumulated_depreciation_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accumulated_depreciation_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check accumulated_depreciation_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("accumulated_depreciation_account with id %s does not exist", id)
	}
	return nil
}





// validateDepreciationExpenseAccountIdExists validates that depreciation_expense_account_id exists
func (v *Validator) validateDepreciationExpenseAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for depreciation_expense_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM depreciation_expense_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check depreciation_expense_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("depreciation_expense_account with id %s does not exist", id)
	}
	return nil
}







// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateAssetCategoriesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateAssetCategoriesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *AssetCategories, req *dto.UpdateAssetCategoriesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *AssetCategories) error {
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

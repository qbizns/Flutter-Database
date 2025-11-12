package inventory_valuation_setting

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles InventoryValuationSettings validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new InventoryValuationSettings validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateInventoryValuationSettingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ValuationMethod
	
	if err := v.validateValuationMethod(req.ValuationMethod); err != nil {
		return err
	}
	
	
	
	// Validate 'fifo',
	
	
	
	// Validate 'lifo',
	
	
	
	// Validate 'weightedAverage',
	
	
	
	// Validate 'movingAverage',
	
	
	
	// Validate 'standardCost',
	
	
	
	// Validate 'specificId'
	
	
	
	// Validate CostLayerGranularity
	
	
	
	// Validate 'product',
	
	
	
	// Validate 'productLocation',
	
	
	
	// Validate 'productLocationLot',
	
	
	
	// Validate 'serialNumber'
	
	
	
	// Validate DefaultInventoryAccountId
	
	
	if err := v.validateDefaultInventoryAccountIdExists(ctx, tx, req.DefaultInventoryAccountId); err != nil {
		return err
	}
	
	
	// Validate DefaultCogsAccountId
	
	
	if err := v.validateDefaultCogsAccountIdExists(ctx, tx, req.DefaultCogsAccountId); err != nil {
		return err
	}
	
	
	// Validate DefaultInventoryAdjustmentAccountId
	
	
	if err := v.validateDefaultInventoryAdjustmentAccountIdExists(ctx, tx, req.DefaultInventoryAdjustmentAccountId); err != nil {
		return err
	}
	
	
	// Validate DefaultInventoryVarianceAccountId
	
	
	if err := v.validateDefaultInventoryVarianceAccountIdExists(ctx, tx, req.DefaultInventoryVarianceAccountId); err != nil {
		return err
	}
	
	
	// Validate CogsRecognitionTiming
	
	
	
	// Validate 'onSale',
	
	
	
	// Validate 'onDelivery',
	
	
	
	// Validate 'onPayment'
	
	
	
	// Validate AllowNegativeInventory
	
	
	
	// Validate RevalueOnPurchase
	
	
	
	// Validate RoundUnitCostToDecimals
	
	
	
	// Validate RevaluationFrequency
	
	
	
	// Validate 'realTime',
	
	
	
	// Validate 'daily',
	
	
	
	// Validate 'monthly',
	
	
	
	// Validate 'manual'
	
	
	
	// Validate IsActive
	
	
	
	// Validate EffectiveFrom
	
	
	
	// Validate Notes
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateInventoryValuationSettingsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate ValuationMethod if provided
	
	if req.ValuationMethod != nil {
		if err := v.validateValuationMethod(*req.ValuationMethod); err != nil {
			return err
		}
	}
	
	
	
	// Validate 'fifo', if provided
	
	
	
	// Validate 'lifo', if provided
	
	
	
	// Validate 'weightedAverage', if provided
	
	
	
	// Validate 'movingAverage', if provided
	
	
	
	// Validate 'standardCost', if provided
	
	
	
	// Validate 'specificId' if provided
	
	
	
	// Validate CostLayerGranularity if provided
	
	
	
	// Validate 'product', if provided
	
	
	
	// Validate 'productLocation', if provided
	
	
	
	// Validate 'productLocationLot', if provided
	
	
	
	// Validate 'serialNumber' if provided
	
	
	
	// Validate DefaultInventoryAccountId if provided
	
	
	if req.DefaultInventoryAccountId != nil {
		if err := v.validateDefaultInventoryAccountIdExists(ctx, tx, *req.DefaultInventoryAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate DefaultCogsAccountId if provided
	
	
	if req.DefaultCogsAccountId != nil {
		if err := v.validateDefaultCogsAccountIdExists(ctx, tx, *req.DefaultCogsAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate DefaultInventoryAdjustmentAccountId if provided
	
	
	if req.DefaultInventoryAdjustmentAccountId != nil {
		if err := v.validateDefaultInventoryAdjustmentAccountIdExists(ctx, tx, *req.DefaultInventoryAdjustmentAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate DefaultInventoryVarianceAccountId if provided
	
	
	if req.DefaultInventoryVarianceAccountId != nil {
		if err := v.validateDefaultInventoryVarianceAccountIdExists(ctx, tx, *req.DefaultInventoryVarianceAccountId); err != nil {
			return err
		}
	}
	
	
	// Validate CogsRecognitionTiming if provided
	
	
	
	// Validate 'onSale', if provided
	
	
	
	// Validate 'onDelivery', if provided
	
	
	
	// Validate 'onPayment' if provided
	
	
	
	// Validate AllowNegativeInventory if provided
	
	
	
	// Validate RevalueOnPurchase if provided
	
	
	
	// Validate RoundUnitCostToDecimals if provided
	
	
	
	// Validate RevaluationFrequency if provided
	
	
	
	// Validate 'realTime', if provided
	
	
	
	// Validate 'daily', if provided
	
	
	
	// Validate 'monthly', if provided
	
	
	
	// Validate 'manual' if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate EffectiveFrom if provided
	
	
	
	// Validate Notes if provided
	
	
	
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



// validateValuationMethod validates valuation_method field
func (v *Validator) validateValuationMethod(value string) error {
	
	// Add custom validation for valuation_method
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("valuation_method cannot be empty")
	}
	
	return nil
}



















































// validateDefaultInventoryAccountIdExists validates that default_inventory_account_id exists
func (v *Validator) validateDefaultInventoryAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for default_inventory_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM default_inventory_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check default_inventory_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("default_inventory_account with id %s does not exist", id)
	}
	return nil
}





// validateDefaultCogsAccountIdExists validates that default_cogs_account_id exists
func (v *Validator) validateDefaultCogsAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for default_cogs_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM default_cogs_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check default_cogs_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("default_cogs_account with id %s does not exist", id)
	}
	return nil
}





// validateDefaultInventoryAdjustmentAccountIdExists validates that default_inventory_adjustment_account_id exists
func (v *Validator) validateDefaultInventoryAdjustmentAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for default_inventory_adjustment_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM default_inventory_adjustment_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check default_inventory_adjustment_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("default_inventory_adjustment_account with id %s does not exist", id)
	}
	return nil
}





// validateDefaultInventoryVarianceAccountIdExists validates that default_inventory_variance_account_id exists
func (v *Validator) validateDefaultInventoryVarianceAccountIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for default_inventory_variance_account
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM default_inventory_variance_account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check default_inventory_variance_account existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("default_inventory_variance_account with id %s does not exist", id)
	}
	return nil
}











































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateInventoryValuationSettingsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateInventoryValuationSettingsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *InventoryValuationSettings, req *UpdateInventoryValuationSettingsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *InventoryValuationSettings) error {
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

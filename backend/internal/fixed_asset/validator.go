package fixed_asset

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

// Validator handles FixedAssets validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new FixedAssets validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateFixedAssetsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate AssetNumber
	
	if err := v.validateAssetNumber(req.AssetNumber); err != nil {
		return err
	}
	
	
	
	// Validate AssetName
	
	if err := v.validateAssetName(req.AssetName); err != nil {
		return err
	}
	
	
	
	// Validate AssetCategoryId
	
	
	if err := v.validateAssetCategoryIdExists(ctx, tx, req.AssetCategoryId); err != nil {
		return err
	}
	
	
	// Validate AcquisitionDate
	
	
	
	// Validate AcquisitionCost
	
	
	
	// Validate SalvageValue
	
	
	
	// Validate SupplierId
	
	
	if err := v.validateSupplierIdExists(ctx, tx, req.SupplierId); err != nil {
		return err
	}
	
	
	// Validate VendorBillId
	
	
	if err := v.validateVendorBillIdExists(ctx, tx, req.VendorBillId); err != nil {
		return err
	}
	
	
	// Validate DepreciationMethod
	
	if err := v.validateDepreciationMethod(req.DepreciationMethod); err != nil {
		return err
	}
	
	
	
	// Validate UsefulLifeYears
	
	
	
	// Validate DepreciationStartDate
	
	
	
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
	
	
	// Validate CurrentBookValue
	
	
	
	// Validate AccumulatedDepreciation
	
	
	
	// Validate LastDepreciationDate
	
	
	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate Department
	
	
	
	// Validate IsDisposed
	
	
	
	// Validate DisposalDate
	
	
	
	// Validate DisposalProceeds
	
	
	
	// Validate DisposalJournalEntryId
	
	
	if err := v.validateDisposalJournalEntryIdExists(ctx, tx, req.DisposalJournalEntryId); err != nil {
		return err
	}
	
	
	// Validate Description
	
	
	
	// Validate SerialNumber
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateFixedAssetsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate AssetNumber if provided
	
	if req.AssetNumber != nil {
		if err := v.validateAssetNumber(*req.AssetNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate AssetName if provided
	
	if req.AssetName != nil {
		if err := v.validateAssetName(*req.AssetName); err != nil {
			return err
		}
	}
	
	
	
	// Validate AssetCategoryId if provided
	
	
	if req.AssetCategoryId != nil {
		if err := v.validateAssetCategoryIdExists(ctx, tx, *req.AssetCategoryId); err != nil {
			return err
		}
	}
	
	
	// Validate AcquisitionDate if provided
	
	
	
	// Validate AcquisitionCost if provided
	
	
	
	// Validate SalvageValue if provided
	
	
	
	// Validate SupplierId if provided
	
	
	if req.SupplierId != nil {
		if err := v.validateSupplierIdExists(ctx, tx, *req.SupplierId); err != nil {
			return err
		}
	}
	
	
	// Validate VendorBillId if provided
	
	
	if req.VendorBillId != nil {
		if err := v.validateVendorBillIdExists(ctx, tx, *req.VendorBillId); err != nil {
			return err
		}
	}
	
	
	// Validate DepreciationMethod if provided
	
	if req.DepreciationMethod != nil {
		if err := v.validateDepreciationMethod(*req.DepreciationMethod); err != nil {
			return err
		}
	}
	
	
	
	// Validate UsefulLifeYears if provided
	
	
	
	// Validate DepreciationStartDate if provided
	
	
	
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
	
	
	// Validate CurrentBookValue if provided
	
	
	
	// Validate AccumulatedDepreciation if provided
	
	
	
	// Validate LastDepreciationDate if provided
	
	
	
	// Validate LocationId if provided
	
	
	if req.LocationId != nil {
		if err := v.validateLocationIdExists(ctx, tx, *req.LocationId); err != nil {
			return err
		}
	}
	
	
	// Validate Department if provided
	
	
	
	// Validate IsDisposed if provided
	
	
	
	// Validate DisposalDate if provided
	
	
	
	// Validate DisposalProceeds if provided
	
	
	
	// Validate DisposalJournalEntryId if provided
	
	
	if req.DisposalJournalEntryId != nil {
		if err := v.validateDisposalJournalEntryIdExists(ctx, tx, *req.DisposalJournalEntryId); err != nil {
			return err
		}
	}
	
	
	// Validate Description if provided
	
	
	
	// Validate SerialNumber if provided
	
	
	
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



// validateAssetNumber validates asset_number field
func (v *Validator) validateAssetNumber(value string) error {
	
	// Add custom validation for asset_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("asset_number cannot be empty")
	}
	
	return nil
}





// validateAssetName validates asset_name field
func (v *Validator) validateAssetName(value string) error {
	
	// Add custom validation for asset_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("asset_name cannot be empty")
	}
	
	return nil
}







// validateAssetCategoryIdExists validates that asset_category_id exists
func (v *Validator) validateAssetCategoryIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for asset_category
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM asset_category WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check asset_category existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("asset_category with id %s does not exist", id)
	}
	return nil
}

















// validateSupplierIdExists validates that supplier_id exists
func (v *Validator) validateSupplierIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for supplier
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM supplier WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check supplier existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("supplier with id %s does not exist", id)
	}
	return nil
}





// validateVendorBillIdExists validates that vendor_bill_id exists
func (v *Validator) validateVendorBillIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for vendor_bill
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM vendor_bill WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check vendor_bill existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("vendor_bill with id %s does not exist", id)
	}
	return nil
}



// validateDepreciationMethod validates depreciation_method field
func (v *Validator) validateDepreciationMethod(value string) error {
	
	// Add custom validation for depreciation_method
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("depreciation_method cannot be empty")
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





















// validateDisposalJournalEntryIdExists validates that disposal_journal_entry_id exists
func (v *Validator) validateDisposalJournalEntryIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for disposal_journal_entry
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM disposal_journal_entry WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check disposal_journal_entry existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("disposal_journal_entry with id %s does not exist", id)
	}
	return nil
}



























// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateFixedAssetsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateFixedAssetsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *FixedAssets, req *dto.UpdateFixedAssetsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *FixedAssets) error {
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

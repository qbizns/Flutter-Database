package fixed_asset

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FixedAssetsResponse represents a fixed_assets response
type FixedAssetsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	AssetNumber string `json:"asset_number"`
	
	AssetName string `json:"asset_name"`
	
	AssetCategoryId *uuid.UUID `json:"asset_category_id"`
	
	AcquisitionDate time.Time `json:"acquisition_date"`
	
	AcquisitionCost float64 `json:"acquisition_cost"`
	
	SalvageValue *float64 `json:"salvage_value"`
	
	SupplierId *uuid.UUID `json:"supplier_id"`
	
	VendorBillId *uuid.UUID `json:"vendor_bill_id"`
	
	DepreciationMethod string `json:"depreciation_method"`
	
	UsefulLifeYears int64 `json:"useful_life_years"`
	
	DepreciationStartDate time.Time `json:"depreciation_start_date"`
	
	AssetAccountId uuid.UUID `json:"asset_account_id"`
	
	AccumulatedDepreciationAccountId uuid.UUID `json:"accumulated_depreciation_account_id"`
	
	DepreciationExpenseAccountId uuid.UUID `json:"depreciation_expense_account_id"`
	
	CurrentBookValue *float64 `json:"current_book_value"`
	
	AccumulatedDepreciation *float64 `json:"accumulated_depreciation"`
	
	LastDepreciationDate *time.Time `json:"last_depreciation_date"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	IsDisposed *bool `json:"is_disposed"`
	
	DisposalDate *time.Time `json:"disposal_date"`
	
	DisposalProceeds *float64 `json:"disposal_proceeds"`
	
	DisposalJournalEntryId *uuid.UUID `json:"disposal_journal_entry_id"`
	
	Description *string `json:"description"`
	
	SerialNumber *string `json:"serial_number"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateFixedAssetsRequest represents a request to create a fixed_assets
type CreateFixedAssetsRequest struct {
	
	AssetNumber string `json:"asset_number" validate:"required"`
	
	AssetName string `json:"asset_name" validate:"required"`
	
	AssetCategoryId *uuid.UUID `json:"asset_category_id"`
	
	AcquisitionDate time.Time `json:"acquisition_date" validate:"required"`
	
	AcquisitionCost float64 `json:"acquisition_cost" validate:"required"`
	
	SalvageValue *float64 `json:"salvage_value"`
	
	SupplierId *uuid.UUID `json:"supplier_id"`
	
	VendorBillId *uuid.UUID `json:"vendor_bill_id"`
	
	DepreciationMethod string `json:"depreciation_method" validate:"required"`
	
	UsefulLifeYears int64 `json:"useful_life_years" validate:"required"`
	
	DepreciationStartDate time.Time `json:"depreciation_start_date" validate:"required"`
	
	AssetAccountId uuid.UUID `json:"asset_account_id" validate:"required"`
	
	AccumulatedDepreciationAccountId uuid.UUID `json:"accumulated_depreciation_account_id" validate:"required"`
	
	DepreciationExpenseAccountId uuid.UUID `json:"depreciation_expense_account_id" validate:"required"`
	
	CurrentBookValue *float64 `json:"current_book_value"`
	
	AccumulatedDepreciation *float64 `json:"accumulated_depreciation"`
	
	LastDepreciationDate *time.Time `json:"last_depreciation_date"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	IsDisposed *bool `json:"is_disposed"`
	
	DisposalDate *time.Time `json:"disposal_date"`
	
	DisposalProceeds *float64 `json:"disposal_proceeds"`
	
	DisposalJournalEntryId *uuid.UUID `json:"disposal_journal_entry_id"`
	
	Description *string `json:"description"`
	
	SerialNumber *string `json:"serial_number"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateFixedAssetsRequest) Validate() error {
	
	if r.AssetNumber == "" {
		return fmt.Errorf("asset_number is required")
	}
	
	if r.AssetName == "" {
		return fmt.Errorf("asset_name is required")
	}
	
	if r.AcquisitionDate.IsZero() {
		return fmt.Errorf("acquisition_date is required")
	}
	
	if r.AcquisitionCost == nil {
		return fmt.Errorf("acquisition_cost is required")
	}
	
	if r.DepreciationMethod == "" {
		return fmt.Errorf("depreciation_method is required")
	}
	
	if r.UsefulLifeYears == 0 {
		return fmt.Errorf("useful_life_years is required")
	}
	
	if r.DepreciationStartDate.IsZero() {
		return fmt.Errorf("depreciation_start_date is required")
	}
	
	if r.AssetAccountId == uuid.Nil {
		return fmt.Errorf("asset_account_id is required")
	}
	
	if r.AccumulatedDepreciationAccountId == uuid.Nil {
		return fmt.Errorf("accumulated_depreciation_account_id is required")
	}
	
	if r.DepreciationExpenseAccountId == uuid.Nil {
		return fmt.Errorf("depreciation_expense_account_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateFixedAssetsRequest represents a request to update a fixed_assets
type UpdateFixedAssetsRequest struct {
	
	AssetNumber *string `json:"asset_number,omitempty" validate:"omitempty,required"`
	
	AssetName *string `json:"asset_name,omitempty" validate:"omitempty,required"`
	
	AssetCategoryId *uuid.UUID `json:"asset_category_id,omitempty"`
	
	AcquisitionDate *time.Time `json:"acquisition_date,omitempty" validate:"omitempty,required"`
	
	AcquisitionCost *float64 `json:"acquisition_cost,omitempty" validate:"omitempty,required"`
	
	SalvageValue *float64 `json:"salvage_value,omitempty"`
	
	SupplierId *uuid.UUID `json:"supplier_id,omitempty"`
	
	VendorBillId *uuid.UUID `json:"vendor_bill_id,omitempty"`
	
	DepreciationMethod *string `json:"depreciation_method,omitempty" validate:"omitempty,required"`
	
	UsefulLifeYears *int64 `json:"useful_life_years,omitempty" validate:"omitempty,required"`
	
	DepreciationStartDate *time.Time `json:"depreciation_start_date,omitempty" validate:"omitempty,required"`
	
	AssetAccountId *uuid.UUID `json:"asset_account_id,omitempty" validate:"omitempty,required"`
	
	AccumulatedDepreciationAccountId *uuid.UUID `json:"accumulated_depreciation_account_id,omitempty" validate:"omitempty,required"`
	
	DepreciationExpenseAccountId *uuid.UUID `json:"depreciation_expense_account_id,omitempty" validate:"omitempty,required"`
	
	CurrentBookValue *float64 `json:"current_book_value,omitempty"`
	
	AccumulatedDepreciation *float64 `json:"accumulated_depreciation,omitempty"`
	
	LastDepreciationDate *time.Time `json:"last_depreciation_date,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	Department *string `json:"department,omitempty"`
	
	IsDisposed *bool `json:"is_disposed,omitempty"`
	
	DisposalDate *time.Time `json:"disposal_date,omitempty"`
	
	DisposalProceeds *float64 `json:"disposal_proceeds,omitempty"`
	
	DisposalJournalEntryId *uuid.UUID `json:"disposal_journal_entry_id,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	SerialNumber *string `json:"serial_number,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateFixedAssetsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.AssetNumber != nil {
		hasUpdate = true
	}
	
	if r.AssetName != nil {
		hasUpdate = true
	}
	
	if r.AssetCategoryId != nil {
		hasUpdate = true
	}
	
	if r.AcquisitionDate != nil {
		hasUpdate = true
	}
	
	if r.AcquisitionCost != nil {
		hasUpdate = true
	}
	
	if r.SalvageValue != nil {
		hasUpdate = true
	}
	
	if r.SupplierId != nil {
		hasUpdate = true
	}
	
	if r.VendorBillId != nil {
		hasUpdate = true
	}
	
	if r.DepreciationMethod != nil {
		hasUpdate = true
	}
	
	if r.UsefulLifeYears != nil {
		hasUpdate = true
	}
	
	if r.DepreciationStartDate != nil {
		hasUpdate = true
	}
	
	if r.AssetAccountId != nil {
		hasUpdate = true
	}
	
	if r.AccumulatedDepreciationAccountId != nil {
		hasUpdate = true
	}
	
	if r.DepreciationExpenseAccountId != nil {
		hasUpdate = true
	}
	
	if r.CurrentBookValue != nil {
		hasUpdate = true
	}
	
	if r.AccumulatedDepreciation != nil {
		hasUpdate = true
	}
	
	if r.LastDepreciationDate != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.Department != nil {
		hasUpdate = true
	}
	
	if r.IsDisposed != nil {
		hasUpdate = true
	}
	
	if r.DisposalDate != nil {
		hasUpdate = true
	}
	
	if r.DisposalProceeds != nil {
		hasUpdate = true
	}
	
	if r.DisposalJournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.SerialNumber != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// FixedAssetsListResponse represents a paginated list of fixed_assets records
type FixedAssetsListResponse struct {
	Items      []*FixedAssetsResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

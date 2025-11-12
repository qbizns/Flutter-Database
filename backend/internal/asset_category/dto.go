package asset_category

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AssetCategoriesResponse represents a asset_categories response
type AssetCategoriesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	CategoryCode string `json:"category_code"`
	
	CategoryName string `json:"category_name"`
	
	DefaultDepreciationMethod *string `json:"default_depreciation_method"`
	
	DefaultUsefulLifeYears *int64 `json:"default_useful_life_years"`
	
	DefaultSalvageValuePercent *float64 `json:"default_salvage_value_percent"`
	
	AssetAccountId *uuid.UUID `json:"asset_account_id"`
	
	AccumulatedDepreciationAccountId *uuid.UUID `json:"accumulated_depreciation_account_id"`
	
	DepreciationExpenseAccountId *uuid.UUID `json:"depreciation_expense_account_id"`
	
	Description *string `json:"description"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateAssetCategoriesRequest represents a request to create a asset_categories
type CreateAssetCategoriesRequest struct {
	
	CategoryCode string `json:"category_code" validate:"required"`
	
	CategoryName string `json:"category_name" validate:"required"`
	
	DefaultDepreciationMethod *string `json:"default_depreciation_method"`
	
	DefaultUsefulLifeYears *int64 `json:"default_useful_life_years"`
	
	DefaultSalvageValuePercent *float64 `json:"default_salvage_value_percent"`
	
	AssetAccountId *uuid.UUID `json:"asset_account_id"`
	
	AccumulatedDepreciationAccountId *uuid.UUID `json:"accumulated_depreciation_account_id"`
	
	DepreciationExpenseAccountId *uuid.UUID `json:"depreciation_expense_account_id"`
	
	Description *string `json:"description"`
	
}

// Validate validates the create request
func (r *CreateAssetCategoriesRequest) Validate() error {
	
	if r.CategoryCode == "" {
		return fmt.Errorf("category_code is required")
	}
	
	if r.CategoryName == "" {
		return fmt.Errorf("category_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateAssetCategoriesRequest represents a request to update a asset_categories
type UpdateAssetCategoriesRequest struct {
	
	CategoryCode *string `json:"category_code,omitempty" validate:"omitempty,required"`
	
	CategoryName *string `json:"category_name,omitempty" validate:"omitempty,required"`
	
	DefaultDepreciationMethod *string `json:"default_depreciation_method,omitempty"`
	
	DefaultUsefulLifeYears *int64 `json:"default_useful_life_years,omitempty"`
	
	DefaultSalvageValuePercent *float64 `json:"default_salvage_value_percent,omitempty"`
	
	AssetAccountId *uuid.UUID `json:"asset_account_id,omitempty"`
	
	AccumulatedDepreciationAccountId *uuid.UUID `json:"accumulated_depreciation_account_id,omitempty"`
	
	DepreciationExpenseAccountId *uuid.UUID `json:"depreciation_expense_account_id,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAssetCategoriesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CategoryCode != nil {
		hasUpdate = true
	}
	
	if r.CategoryName != nil {
		hasUpdate = true
	}
	
	if r.DefaultDepreciationMethod != nil {
		hasUpdate = true
	}
	
	if r.DefaultUsefulLifeYears != nil {
		hasUpdate = true
	}
	
	if r.DefaultSalvageValuePercent != nil {
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
	
	if r.Description != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// AssetCategoriesListResponse represents a paginated list of asset_categories records
type AssetCategoriesListResponse struct {
	Items      []*AssetCategoriesResponse `json:"items"`
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

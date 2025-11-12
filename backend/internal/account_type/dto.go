package account_type

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccountTypesResponse represents a account_types response
type AccountTypesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	TypeCode string `json:"type_code"`
	
	TypeName string `json:"type_name"`
	
	TypeCategory string `json:"type_category"`
	
	NormalBalance string `json:"normal_balance"`
	
	IsBalanceSheet *bool `json:"is_balance_sheet"`
	
	IsIncomeStatement *bool `json:"is_income_statement"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	Description *string `json:"description"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateAccountTypesRequest represents a request to create a account_types
type CreateAccountTypesRequest struct {
	
	TypeCode string `json:"type_code" validate:"required"`
	
	TypeName string `json:"type_name" validate:"required"`
	
	TypeCategory string `json:"type_category" validate:"required"`
	
	NormalBalance string `json:"normal_balance" validate:"required"`
	
	IsBalanceSheet *bool `json:"is_balance_sheet"`
	
	IsIncomeStatement *bool `json:"is_income_statement"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	Description *string `json:"description"`
	
}

// Validate validates the create request
func (r *CreateAccountTypesRequest) Validate() error {
	
	if r.TypeCode == "" {
		return fmt.Errorf("type_code is required")
	}
	
	if r.TypeName == "" {
		return fmt.Errorf("type_name is required")
	}
	
	if r.TypeCategory == "" {
		return fmt.Errorf("type_category is required")
	}
	
	if r.NormalBalance == "" {
		return fmt.Errorf("normal_balance is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateAccountTypesRequest represents a request to update a account_types
type UpdateAccountTypesRequest struct {
	
	TypeCode *string `json:"type_code,omitempty" validate:"omitempty,required"`
	
	TypeName *string `json:"type_name,omitempty" validate:"omitempty,required"`
	
	TypeCategory *string `json:"type_category,omitempty" validate:"omitempty,required"`
	
	NormalBalance *string `json:"normal_balance,omitempty" validate:"omitempty,required"`
	
	IsBalanceSheet *bool `json:"is_balance_sheet,omitempty"`
	
	IsIncomeStatement *bool `json:"is_income_statement,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAccountTypesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TypeCode != nil {
		hasUpdate = true
	}
	
	if r.TypeName != nil {
		hasUpdate = true
	}
	
	if r.TypeCategory != nil {
		hasUpdate = true
	}
	
	if r.NormalBalance != nil {
		hasUpdate = true
	}
	
	if r.IsBalanceSheet != nil {
		hasUpdate = true
	}
	
	if r.IsIncomeStatement != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
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

// AccountTypesListResponse represents a paginated list of account_types records
type AccountTypesListResponse struct {
	Items      []*AccountTypesResponse `json:"items"`
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

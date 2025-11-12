package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccountSubtypesResponse represents a account_subtypes response
type AccountSubtypesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	AccountTypeId uuid.UUID `json:"account_type_id"`
	
	SubtypeCode string `json:"subtype_code"`
	
	SubtypeName string `json:"subtype_name"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	Description *string `json:"description"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateAccountSubtypesRequest represents a request to create a account_subtypes
type CreateAccountSubtypesRequest struct {
	
	AccountTypeId uuid.UUID `json:"account_type_id" validate:"required"`
	
	SubtypeCode string `json:"subtype_code" validate:"required"`
	
	SubtypeName string `json:"subtype_name" validate:"required"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	Description *string `json:"description"`
	
}

// Validate validates the create request
func (r *CreateAccountSubtypesRequest) Validate() error {
	
	if r.AccountTypeId == uuid.Nil {
		return fmt.Errorf("account_type_id is required")
	}
	
	if r.SubtypeCode == "" {
		return fmt.Errorf("subtype_code is required")
	}
	
	if r.SubtypeName == "" {
		return fmt.Errorf("subtype_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateAccountSubtypesRequest represents a request to update a account_subtypes
type UpdateAccountSubtypesRequest struct {
	
	AccountTypeId *uuid.UUID `json:"account_type_id,omitempty" validate:"omitempty,required"`
	
	SubtypeCode *string `json:"subtype_code,omitempty" validate:"omitempty,required"`
	
	SubtypeName *string `json:"subtype_name,omitempty" validate:"omitempty,required"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAccountSubtypesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.AccountTypeId != nil {
		hasUpdate = true
	}
	
	if r.SubtypeCode != nil {
		hasUpdate = true
	}
	
	if r.SubtypeName != nil {
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

// AccountSubtypesListResponse represents a paginated list of account_subtypes records
type AccountSubtypesListResponse struct {
	Items      []*AccountSubtypesResponse `json:"items"`
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

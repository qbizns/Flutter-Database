package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingProfilesResponse represents a posting_profiles response
type PostingProfilesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	Code string `json:"code"`
	
	Name string `json:"name"`
	
	Description *string `json:"description"`
	
	IsDefault bool `json:"is_default"`
	
	IsActive bool `json:"is_active"`
	
	DefaultFiscalYearId *uuid.UUID `json:"default_fiscal_year_id"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreatePostingProfilesRequest represents a request to create a posting_profiles
type CreatePostingProfilesRequest struct {
	
	Code string `json:"code" validate:"required"`
	
	Name string `json:"name" validate:"required"`
	
	Description *string `json:"description"`
	
	IsDefault bool `json:"is_default" validate:"required"`
	
	IsActive bool `json:"is_active" validate:"required"`
	
	DefaultFiscalYearId *uuid.UUID `json:"default_fiscal_year_id"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePostingProfilesRequest) Validate() error {
	
	if r.Code == "" {
		return fmt.Errorf("code is required")
	}
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.IsDefault == nil {
		return fmt.Errorf("is_default is required")
	}
	
	if r.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingProfilesRequest represents a request to update a posting_profiles
type UpdatePostingProfilesRequest struct {
	
	Code *string `json:"code,omitempty" validate:"omitempty,required"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty" validate:"omitempty,required"`
	
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty,required"`
	
	DefaultFiscalYearId *uuid.UUID `json:"default_fiscal_year_id,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingProfilesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Code != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.DefaultFiscalYearId != nil {
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

// PostingProfilesListResponse represents a paginated list of posting_profiles records
type PostingProfilesListResponse struct {
	Items      []*PostingProfilesResponse `json:"items"`
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

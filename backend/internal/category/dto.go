package category

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CategoriesResponse represents a categories response
type CategoriesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	Name string `json:"name"`
	
	Slug string `json:"slug"`
	
	Description *string `json:"description"`
	
	ParentId *uuid.UUID `json:"parent_id"`
	
	Level *int64 `json:"level"`
	
	Path *string `json:"path"`
	
	ImageUrl *string `json:"image_url"`
	
	Icon *string `json:"icon"`
	
	Color *string `json:"color"`
	
	SortOrder *int64 `json:"sort_order"`
	
	IsActive *bool `json:"is_active"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateCategoriesRequest represents a request to create a categories
type CreateCategoriesRequest struct {
	
	Name string `json:"name" validate:"required"`
	
	Slug string `json:"slug" validate:"required"`
	
	Description *string `json:"description"`
	
	ParentId *uuid.UUID `json:"parent_id"`
	
	Level *int64 `json:"level"`
	
	Path *string `json:"path"`
	
	ImageUrl *string `json:"image_url" validate:"url"`
	
	Icon *string `json:"icon"`
	
	Color *string `json:"color"`
	
	SortOrder *int64 `json:"sort_order"`
	
	IsActive *bool `json:"is_active"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateCategoriesRequest) Validate() error {
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCategoriesRequest represents a request to update a categories
type UpdateCategoriesRequest struct {
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Slug *string `json:"slug,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	ParentId *uuid.UUID `json:"parent_id,omitempty"`
	
	Level *int64 `json:"level,omitempty"`
	
	Path *string `json:"path,omitempty"`
	
	ImageUrl *string `json:"image_url,omitempty" validate:"omitempty,url"`
	
	Icon *string `json:"icon,omitempty"`
	
	Color *string `json:"color,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCategoriesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Slug != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.ParentId != nil {
		hasUpdate = true
	}
	
	if r.Level != nil {
		hasUpdate = true
	}
	
	if r.Path != nil {
		hasUpdate = true
	}
	
	if r.ImageUrl != nil {
		hasUpdate = true
	}
	
	if r.Icon != nil {
		hasUpdate = true
	}
	
	if r.Color != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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

// CategoriesListResponse represents a paginated list of categories records
type CategoriesListResponse struct {
	Items      []*CategoriesResponse `json:"items"`
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

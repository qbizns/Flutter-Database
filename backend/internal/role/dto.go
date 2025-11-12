package role

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RolesResponse represents a roles response
type RolesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	Name string `json:"name"`
	
	Slug string `json:"slug"`
	
	Description *string `json:"description"`
	
	IsSystemRole *bool `json:"is_system_role"`
	
	IsDefault *bool `json:"is_default"`
	
	Settings json.RawMessage `json:"settings"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateRolesRequest represents a request to create a roles
type CreateRolesRequest struct {
	
	Name string `json:"name" validate:"required"`
	
	Slug string `json:"slug" validate:"required"`
	
	Description *string `json:"description"`
	
	IsSystemRole *bool `json:"is_system_role"`
	
	IsDefault *bool `json:"is_default"`
	
	Settings json.RawMessage `json:"settings"`
	
}

// Validate validates the create request
func (r *CreateRolesRequest) Validate() error {
	
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

// UpdateRolesRequest represents a request to update a roles
type UpdateRolesRequest struct {
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Slug *string `json:"slug,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	IsSystemRole *bool `json:"is_system_role,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	Settings *json.RawMessage `json:"settings,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateRolesRequest) Validate() error {
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
	
	if r.IsSystemRole != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.Settings != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// RolesListResponse represents a paginated list of roles records
type RolesListResponse struct {
	Items      []*RolesResponse `json:"items"`
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

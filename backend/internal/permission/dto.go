package permission

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PermissionsResponse represents a permissions response
type PermissionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	Name string `json:"name"`
	
	Slug string `json:"slug"`
	
	Description *string `json:"description"`
	
	Resource string `json:"resource"`
	
	Action string `json:"action"`
	
	Category *string `json:"category"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
}

// CreatePermissionsRequest represents a request to create a permissions
type CreatePermissionsRequest struct {
	
	Name string `json:"name" validate:"required"`
	
	Slug string `json:"slug" validate:"required"`
	
	Description *string `json:"description"`
	
	Resource string `json:"resource" validate:"required"`
	
	Action string `json:"action" validate:"required"`
	
	Category *string `json:"category"`
	
}

// Validate validates the create request
func (r *CreatePermissionsRequest) Validate() error {
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	
	if r.Resource == "" {
		return fmt.Errorf("resource is required")
	}
	
	if r.Action == "" {
		return fmt.Errorf("action is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePermissionsRequest represents a request to update a permissions
type UpdatePermissionsRequest struct {
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Slug *string `json:"slug,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	Resource *string `json:"resource,omitempty" validate:"omitempty,required"`
	
	Action *string `json:"action,omitempty" validate:"omitempty,required"`
	
	Category *string `json:"category,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePermissionsRequest) Validate() error {
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
	
	if r.Resource != nil {
		hasUpdate = true
	}
	
	if r.Action != nil {
		hasUpdate = true
	}
	
	if r.Category != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PermissionsListResponse represents a paginated list of permissions records
type PermissionsListResponse struct {
	Items      []*PermissionsResponse `json:"items"`
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

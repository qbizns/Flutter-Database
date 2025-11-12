package role_permission

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RolePermissionsResponse represents a role_permissions response
type RolePermissionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	RoleId uuid.UUID `json:"role_id"`
	
	PermissionId uuid.UUID `json:"permission_id"`
	
	CreatedAt time.Time `json:"created_at"`
	
}

// CreateRolePermissionsRequest represents a request to create a role_permissions
type CreateRolePermissionsRequest struct {
	
	RoleId uuid.UUID `json:"role_id" validate:"required"`
	
	PermissionId uuid.UUID `json:"permission_id" validate:"required"`
	
}

// Validate validates the create request
func (r *CreateRolePermissionsRequest) Validate() error {
	
	if r.RoleId == uuid.Nil {
		return fmt.Errorf("role_id is required")
	}
	
	if r.PermissionId == uuid.Nil {
		return fmt.Errorf("permission_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateRolePermissionsRequest represents a request to update a role_permissions
type UpdateRolePermissionsRequest struct {
	
	RoleId *uuid.UUID `json:"role_id,omitempty" validate:"omitempty,required"`
	
	PermissionId *uuid.UUID `json:"permission_id,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdateRolePermissionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.RoleId != nil {
		hasUpdate = true
	}
	
	if r.PermissionId != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// RolePermissionsListResponse represents a paginated list of role_permissions records
type RolePermissionsListResponse struct {
	Items      []*RolePermissionsResponse `json:"items"`
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

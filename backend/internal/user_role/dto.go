package user_role

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserRolesResponse represents a user_roles response
type UserRolesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	UserId uuid.UUID `json:"user_id"`
	
	RoleId uuid.UUID `json:"role_id"`
	
	CreatedAt time.Time `json:"created_at"`
	
	AssignedBy *uuid.UUID `json:"assigned_by"`
	
}

// CreateUserRolesRequest represents a request to create a user_roles
type CreateUserRolesRequest struct {
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	RoleId uuid.UUID `json:"role_id" validate:"required"`
	
	AssignedBy *uuid.UUID `json:"assigned_by"`
	
}

// Validate validates the create request
func (r *CreateUserRolesRequest) Validate() error {
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	
	if r.RoleId == uuid.Nil {
		return fmt.Errorf("role_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateUserRolesRequest represents a request to update a user_roles
type UpdateUserRolesRequest struct {
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	RoleId *uuid.UUID `json:"role_id,omitempty" validate:"omitempty,required"`
	
	AssignedBy *uuid.UUID `json:"assigned_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateUserRolesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.RoleId != nil {
		hasUpdate = true
	}
	
	if r.AssignedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UserRolesListResponse represents a paginated list of user_roles records
type UserRolesListResponse struct {
	Items      []*UserRolesResponse `json:"items"`
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

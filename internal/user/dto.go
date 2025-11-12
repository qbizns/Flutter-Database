package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UsersResponse represents a users response
type UsersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	Email string `json:"email"`
	
	PasswordHash *string `json:"password_hash"`
	
	FirstName string `json:"first_name"`
	
	LastName string `json:"last_name"`
	
	Phone *string `json:"phone"`
	
	AvatarUrl *string `json:"avatar_url"`
	
	Status string `json:"status"`
	
	EmailVerified *bool `json:"email_verified"`
	
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	
	LastLoginAt *time.Time `json:"last_login_at"`
	
	LastLoginIp *string `json:"last_login_ip"`
	
	FailedLoginAttempts *int64 `json:"failed_login_attempts"`
	
	LockedUntil *time.Time `json:"locked_until"`
	
	TwoFactorEnabled *bool `json:"two_factor_enabled"`
	
	TwoFactorSecret *string `json:"two_factor_secret"`
	
	Settings json.RawMessage `json:"settings"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateUsersRequest represents a request to create a users
type CreateUsersRequest struct {
	
	Email string `json:"email" validate:"required,email"`
	
	PasswordHash *string `json:"password_hash"`
	
	FirstName string `json:"first_name" validate:"required"`
	
	LastName string `json:"last_name" validate:"required"`
	
	Phone *string `json:"phone" validate:"e164"`
	
	AvatarUrl *string `json:"avatar_url" validate:"url"`
	
	Status string `json:"status" validate:"required"`
	
	EmailVerified *bool `json:"email_verified"`
	
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	
	LastLoginAt *time.Time `json:"last_login_at"`
	
	LastLoginIp *string `json:"last_login_ip"`
	
	FailedLoginAttempts *int64 `json:"failed_login_attempts"`
	
	LockedUntil *time.Time `json:"locked_until"`
	
	TwoFactorEnabled *bool `json:"two_factor_enabled"`
	
	TwoFactorSecret *string `json:"two_factor_secret"`
	
	Settings json.RawMessage `json:"settings"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateUsersRequest) Validate() error {
	
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	
	if r.FirstName == "" {
		return fmt.Errorf("first_name is required")
	}
	
	if r.LastName == "" {
		return fmt.Errorf("last_name is required")
	}
	
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateUsersRequest represents a request to update a users
type UpdateUsersRequest struct {
	
	Email *string `json:"email,omitempty" validate:"omitempty,required,email"`
	
	PasswordHash *string `json:"password_hash,omitempty"`
	
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,required"`
	
	LastName *string `json:"last_name,omitempty" validate:"omitempty,required"`
	
	Phone *string `json:"phone,omitempty" validate:"omitempty,e164"`
	
	AvatarUrl *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	
	Status *string `json:"status,omitempty" validate:"omitempty,required"`
	
	EmailVerified *bool `json:"email_verified,omitempty"`
	
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	
	LastLoginIp *string `json:"last_login_ip,omitempty"`
	
	FailedLoginAttempts *int64 `json:"failed_login_attempts,omitempty"`
	
	LockedUntil *time.Time `json:"locked_until,omitempty"`
	
	TwoFactorEnabled *bool `json:"two_factor_enabled,omitempty"`
	
	TwoFactorSecret *string `json:"two_factor_secret,omitempty"`
	
	Settings *json.RawMessage `json:"settings,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateUsersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Email != nil {
		hasUpdate = true
	}
	
	if r.PasswordHash != nil {
		hasUpdate = true
	}
	
	if r.FirstName != nil {
		hasUpdate = true
	}
	
	if r.LastName != nil {
		hasUpdate = true
	}
	
	if r.Phone != nil {
		hasUpdate = true
	}
	
	if r.AvatarUrl != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.EmailVerified != nil {
		hasUpdate = true
	}
	
	if r.EmailVerifiedAt != nil {
		hasUpdate = true
	}
	
	if r.LastLoginAt != nil {
		hasUpdate = true
	}
	
	if r.LastLoginIp != nil {
		hasUpdate = true
	}
	
	if r.FailedLoginAttempts != nil {
		hasUpdate = true
	}
	
	if r.LockedUntil != nil {
		hasUpdate = true
	}
	
	if r.TwoFactorEnabled != nil {
		hasUpdate = true
	}
	
	if r.TwoFactorSecret != nil {
		hasUpdate = true
	}
	
	if r.Settings != nil {
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

// UsersListResponse represents a paginated list of users records
type UsersListResponse struct {
	Items      []*UsersResponse `json:"items"`
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

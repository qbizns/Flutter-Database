package user_session

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserSessionsResponse represents a user_sessions response
type UserSessionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	UserId uuid.UUID `json:"user_id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	SessionToken string `json:"session_token"`
	
	RefreshToken *string `json:"refresh_token"`
	
	UserAgent *string `json:"user_agent"`
	
	IpAddress *string `json:"ip_address"`
	
	DeviceType *string `json:"device_type"`
	
	DeviceName *string `json:"device_name"`
	
	Browser *string `json:"browser"`
	
	Os *string `json:"os"`
	
	CountryCode *string `json:"country_code"`
	
	City *string `json:"city"`
	
	IsActive *bool `json:"is_active"`
	
	LastActivityAt *time.Time `json:"last_activity_at"`
	
	ExpiresAt time.Time `json:"expires_at"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	RevokedAt *time.Time `json:"revoked_at"`
	
}

// CreateUserSessionsRequest represents a request to create a user_sessions
type CreateUserSessionsRequest struct {
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	SessionToken string `json:"session_token" validate:"required"`
	
	RefreshToken *string `json:"refresh_token"`
	
	UserAgent *string `json:"user_agent"`
	
	IpAddress *string `json:"ip_address"`
	
	DeviceType *string `json:"device_type"`
	
	DeviceName *string `json:"device_name"`
	
	Browser *string `json:"browser"`
	
	Os *string `json:"os"`
	
	CountryCode *string `json:"country_code"`
	
	City *string `json:"city"`
	
	IsActive *bool `json:"is_active"`
	
	LastActivityAt *time.Time `json:"last_activity_at"`
	
	// Duplicate removed: ExpiresAt time.Time `json:"expires_at"`
	
	RevokedAt *time.Time `json:"revoked_at"`
	
}

// Validate validates the create request
func (r *CreateUserSessionsRequest) Validate() error {
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	
	if r.SessionToken == "" {
		return fmt.Errorf("session_token is required")
	}
	
	if r.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateUserSessionsRequest represents a request to update a user_sessions
type UpdateUserSessionsRequest struct {
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	SessionToken *string `json:"session_token,omitempty" validate:"omitempty,required"`
	
	RefreshToken *string `json:"refresh_token,omitempty"`
	
	UserAgent *string `json:"user_agent,omitempty"`
	
	IpAddress *string `json:"ip_address,omitempty"`
	
	DeviceType *string `json:"device_type,omitempty"`
	
	DeviceName *string `json:"device_name,omitempty"`
	
	Browser *string `json:"browser,omitempty"`
	
	Os *string `json:"os,omitempty"`
	
	CountryCode *string `json:"country_code,omitempty"`
	
	City *string `json:"city,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	LastActivityAt *time.Time `json:"last_activity_at,omitempty"`
	
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateUserSessionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.SessionToken != nil {
		hasUpdate = true
	}
	
	if r.RefreshToken != nil {
		hasUpdate = true
	}
	
	if r.UserAgent != nil {
		hasUpdate = true
	}
	
	if r.IpAddress != nil {
		hasUpdate = true
	}
	
	if r.DeviceType != nil {
		hasUpdate = true
	}
	
	if r.DeviceName != nil {
		hasUpdate = true
	}
	
	if r.Browser != nil {
		hasUpdate = true
	}
	
	if r.Os != nil {
		hasUpdate = true
	}
	
	if r.CountryCode != nil {
		hasUpdate = true
	}
	
	if r.City != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.LastActivityAt != nil {
		hasUpdate = true
	}
	
	if r.ExpiresAt != nil {
		hasUpdate = true
	}
	
	if r.RevokedAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UserSessionsListResponse represents a paginated list of user_sessions records
type UserSessionsListResponse struct {
	Items      []*UserSessionsResponse `json:"items"`
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

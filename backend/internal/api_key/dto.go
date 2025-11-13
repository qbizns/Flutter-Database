package api_key

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ApiKeysResponse represents a api_keys response
type ApiKeysResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	KeyName string `json:"key_name"`
	
	KeyPrefix string `json:"key_prefix"`
	
	KeyHash string `json:"key_hash"`
	
	Scopes json.RawMessage `json:"scopes"`
	
	AllowedIps *string `json:"allowed_ips"`
	
	IsActive *bool `json:"is_active"`
	
	LastUsedAt *time.Time `json:"last_used_at"`
	
	UsageCount *int64 `json:"usage_count"`
	
	RateLimitPerMinute *int64 `json:"rate_limit_per_minute"`
	
	RateLimitPerHour *int64 `json:"rate_limit_per_hour"`
	
	ExpiresAt *time.Time `json:"expires_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateApiKeysRequest represents a request to create a api_keys
type CreateApiKeysRequest struct {
	
	KeyName string `json:"key_name" validate:"required"`
	
	KeyPrefix string `json:"key_prefix" validate:"required"`
	
	KeyHash string `json:"key_hash" validate:"required"`
	
	// Duplicate removed: Scopes json.RawMessage `json:"scopes"`
	
	AllowedIps *string `json:"allowed_ips"`
	
	IsActive *bool `json:"is_active"`
	
	LastUsedAt *time.Time `json:"last_used_at"`
	
	UsageCount *int64 `json:"usage_count"`
	
	RateLimitPerMinute *int64 `json:"rate_limit_per_minute"`
	
	RateLimitPerHour *int64 `json:"rate_limit_per_hour"`
	
	ExpiresAt *time.Time `json:"expires_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateApiKeysRequest) Validate() error {
	
	if r.KeyName == "" {
		return fmt.Errorf("key_name is required")
	}
	
	if r.KeyPrefix == "" {
		return fmt.Errorf("key_prefix is required")
	}
	
	if r.KeyHash == "" {
		return fmt.Errorf("key_hash is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateApiKeysRequest represents a request to update a api_keys
type UpdateApiKeysRequest struct {
	
	KeyName *string `json:"key_name,omitempty" validate:"omitempty,required"`
	
	KeyPrefix *string `json:"key_prefix,omitempty" validate:"omitempty,required"`
	
	KeyHash *string `json:"key_hash,omitempty" validate:"omitempty,required"`
	
	Scopes *json.RawMessage `json:"scopes,omitempty"`
	
	AllowedIps *string `json:"allowed_ips,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	
	UsageCount *int64 `json:"usage_count,omitempty"`
	
	RateLimitPerMinute *int64 `json:"rate_limit_per_minute,omitempty"`
	
	RateLimitPerHour *int64 `json:"rate_limit_per_hour,omitempty"`
	
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateApiKeysRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.KeyName != nil {
		hasUpdate = true
	}
	
	if r.KeyPrefix != nil {
		hasUpdate = true
	}
	
	if r.KeyHash != nil {
		hasUpdate = true
	}
	
	if r.Scopes != nil {
		hasUpdate = true
	}
	
	if r.AllowedIps != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.LastUsedAt != nil {
		hasUpdate = true
	}
	
	if r.UsageCount != nil {
		hasUpdate = true
	}
	
	if r.RateLimitPerMinute != nil {
		hasUpdate = true
	}
	
	if r.RateLimitPerHour != nil {
		hasUpdate = true
	}
	
	if r.ExpiresAt != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ApiKeysListResponse represents a paginated list of api_keys records
type ApiKeysListResponse struct {
	Items      []*ApiKeysResponse `json:"items"`
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

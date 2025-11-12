package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrganizationFeaturesResponse represents a organization_features response
type OrganizationFeaturesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	FeatureKey string `json:"feature_key"`
	
	IsEnabled *bool `json:"is_enabled"`
	
	IsAvailable *bool `json:"is_available"`
	
	Configuration json.RawMessage `json:"configuration"`
	
	Limits json.RawMessage `json:"limits"`
	
	EnabledAt *time.Time `json:"enabled_at"`
	
	DisabledAt *time.Time `json:"disabled_at"`
	
	ExpiresAt *time.Time `json:"expires_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateOrganizationFeaturesRequest represents a request to create a organization_features
type CreateOrganizationFeaturesRequest struct {
	
	FeatureKey string `json:"feature_key" validate:"required"`
	
	IsEnabled *bool `json:"is_enabled"`
	
	IsAvailable *bool `json:"is_available"`
	
	Configuration json.RawMessage `json:"configuration"`
	
	Limits json.RawMessage `json:"limits"`
	
	EnabledAt *time.Time `json:"enabled_at"`
	
	DisabledAt *time.Time `json:"disabled_at"`
	
	ExpiresAt *time.Time `json:"expires_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateOrganizationFeaturesRequest) Validate() error {
	
	if r.FeatureKey == "" {
		return fmt.Errorf("feature_key is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateOrganizationFeaturesRequest represents a request to update a organization_features
type UpdateOrganizationFeaturesRequest struct {
	
	FeatureKey *string `json:"feature_key,omitempty" validate:"omitempty,required"`
	
	IsEnabled *bool `json:"is_enabled,omitempty"`
	
	IsAvailable *bool `json:"is_available,omitempty"`
	
	Configuration *json.RawMessage `json:"configuration,omitempty"`
	
	Limits *json.RawMessage `json:"limits,omitempty"`
	
	EnabledAt *time.Time `json:"enabled_at,omitempty"`
	
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
	
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateOrganizationFeaturesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.FeatureKey != nil {
		hasUpdate = true
	}
	
	if r.IsEnabled != nil {
		hasUpdate = true
	}
	
	if r.IsAvailable != nil {
		hasUpdate = true
	}
	
	if r.Configuration != nil {
		hasUpdate = true
	}
	
	if r.Limits != nil {
		hasUpdate = true
	}
	
	if r.EnabledAt != nil {
		hasUpdate = true
	}
	
	if r.DisabledAt != nil {
		hasUpdate = true
	}
	
	if r.ExpiresAt != nil {
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

// OrganizationFeaturesListResponse represents a paginated list of organization_features records
type OrganizationFeaturesListResponse struct {
	Items      []*OrganizationFeaturesResponse `json:"items"`
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

package organization

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrganizationsResponse represents a organizations response
type OrganizationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	Name string `json:"name"`
	
	Slug string `json:"slug"`
	
	Description *string `json:"description"`
	
	Email *string `json:"email"`
	
	Phone *string `json:"phone"`
	
	Address *string `json:"address"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	Status string `json:"status"`
	
	Plan *string `json:"plan"`
	
	TrialEndsAt *time.Time `json:"trial_ends_at"`
	
	SubscriptionStartsAt *time.Time `json:"subscription_starts_at"`
	
	SubscriptionEndsAt *time.Time `json:"subscription_ends_at"`
	
	MaxUsers *int64 `json:"max_users"`
	
	MaxProducts *int64 `json:"max_products"`
	
	MaxLocations *int64 `json:"max_locations"`
	
	Settings json.RawMessage `json:"settings"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateOrganizationsRequest represents a request to create a organizations
type CreateOrganizationsRequest struct {
	
	Name string `json:"name" validate:"required"`
	
	Slug string `json:"slug" validate:"required"`
	
	Description *string `json:"description"`
	
	Email *string `json:"email" validate:"email"`
	
	Phone *string `json:"phone" validate:"e164"`
	
	Address *string `json:"address"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	// 	Status string `json:"status" validate:"required"`
	
	Plan *string `json:"plan"`
	
	TrialEndsAt *time.Time `json:"trial_ends_at"`
	
	SubscriptionStartsAt *time.Time `json:"subscription_starts_at"`
	
	SubscriptionEndsAt *time.Time `json:"subscription_ends_at"`
	
	MaxUsers *int64 `json:"max_users"`
	
	MaxProducts *int64 `json:"max_products"`
	
	MaxLocations *int64 `json:"max_locations"`
	
	// Duplicate removed: Settings json.RawMessage `json:"settings"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateOrganizationsRequest) Validate() error {
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateOrganizationsRequest represents a request to update a organizations
type UpdateOrganizationsRequest struct {
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Slug *string `json:"slug,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	
	Phone *string `json:"phone,omitempty" validate:"omitempty,e164"`
	
	Address *string `json:"address,omitempty"`
	
	City *string `json:"city,omitempty"`
	
	State *string `json:"state,omitempty"`
	
	Country *string `json:"country,omitempty"`
	
	PostalCode *string `json:"postal_code,omitempty"`
	
	// 	Status *string `json:"status,omitempty" validate:"omitempty,required"`
	
	Plan *string `json:"plan,omitempty"`
	
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
	
	SubscriptionStartsAt *time.Time `json:"subscription_starts_at,omitempty"`
	
	SubscriptionEndsAt *time.Time `json:"subscription_ends_at,omitempty"`
	
	MaxUsers *int64 `json:"max_users,omitempty"`
	
	MaxProducts *int64 `json:"max_products,omitempty"`
	
	MaxLocations *int64 `json:"max_locations,omitempty"`
	
	Settings *json.RawMessage `json:"settings,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateOrganizationsRequest) Validate() error {
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
	
	if r.Email != nil {
		hasUpdate = true
	}
	
	if r.Phone != nil {
		hasUpdate = true
	}
	
	if r.Address != nil {
		hasUpdate = true
	}
	
	if r.City != nil {
		hasUpdate = true
	}
	
	if r.State != nil {
		hasUpdate = true
	}
	
	if r.Country != nil {
		hasUpdate = true
	}
	
	if r.PostalCode != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Plan != nil {
		hasUpdate = true
	}
	
	if r.TrialEndsAt != nil {
		hasUpdate = true
	}
	
	if r.SubscriptionStartsAt != nil {
		hasUpdate = true
	}
	
	if r.SubscriptionEndsAt != nil {
		hasUpdate = true
	}
	
	if r.MaxUsers != nil {
		hasUpdate = true
	}
	
	if r.MaxProducts != nil {
		hasUpdate = true
	}
	
	if r.MaxLocations != nil {
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

// OrganizationsListResponse represents a paginated list of organizations records
type OrganizationsListResponse struct {
	Items      []*OrganizationsResponse `json:"items"`
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

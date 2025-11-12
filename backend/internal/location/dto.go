package location

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LocationsResponse represents a locations response
type LocationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationCode string `json:"location_code"`
	
	Name string `json:"name"`
	
	LocationType *string `json:"location_type"`
	
	Phone *string `json:"phone"`
	
	Email *string `json:"email"`
	
	ManagerUserId *uuid.UUID `json:"manager_user_id"`
	
	AddressLine1 *string `json:"address_line1"`
	
	AddressLine2 *string `json:"address_line2"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	Timezone *string `json:"timezone"`
	
	BusinessHours json.RawMessage `json:"business_hours"`
	
	IsActive *bool `json:"is_active"`
	
	IsPrimary *bool `json:"is_primary"`
	
	AllowSales *bool `json:"allow_sales"`
	
	AllowPurchases *bool `json:"allow_purchases"`
	
	TaxRate *float64 `json:"tax_rate"`
	
	Notes *string `json:"notes"`
	
	Settings json.RawMessage `json:"settings"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateLocationsRequest represents a request to create a locations
type CreateLocationsRequest struct {
	
	LocationCode string `json:"location_code" validate:"required"`
	
	Name string `json:"name" validate:"required"`
	
	LocationType *string `json:"location_type"`
	
	Phone *string `json:"phone" validate:"e164"`
	
	Email *string `json:"email" validate:"email"`
	
	ManagerUserId *uuid.UUID `json:"manager_user_id"`
	
	AddressLine1 *string `json:"address_line1"`
	
	AddressLine2 *string `json:"address_line2"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	Timezone *string `json:"timezone"`
	
	BusinessHours json.RawMessage `json:"business_hours"`
	
	IsActive *bool `json:"is_active"`
	
	IsPrimary *bool `json:"is_primary"`
	
	AllowSales *bool `json:"allow_sales"`
	
	AllowPurchases *bool `json:"allow_purchases"`
	
	TaxRate *float64 `json:"tax_rate"`
	
	Notes *string `json:"notes"`
	
	Settings json.RawMessage `json:"settings"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateLocationsRequest) Validate() error {
	
	if r.LocationCode == "" {
		return fmt.Errorf("location_code is required")
	}
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLocationsRequest represents a request to update a locations
type UpdateLocationsRequest struct {
	
	LocationCode *string `json:"location_code,omitempty" validate:"omitempty,required"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	LocationType *string `json:"location_type,omitempty"`
	
	Phone *string `json:"phone,omitempty" validate:"omitempty,e164"`
	
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	
	ManagerUserId *uuid.UUID `json:"manager_user_id,omitempty"`
	
	AddressLine1 *string `json:"address_line1,omitempty"`
	
	AddressLine2 *string `json:"address_line2,omitempty"`
	
	City *string `json:"city,omitempty"`
	
	State *string `json:"state,omitempty"`
	
	Country *string `json:"country,omitempty"`
	
	PostalCode *string `json:"postal_code,omitempty"`
	
	Timezone *string `json:"timezone,omitempty"`
	
	BusinessHours *json.RawMessage `json:"business_hours,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsPrimary *bool `json:"is_primary,omitempty"`
	
	AllowSales *bool `json:"allow_sales,omitempty"`
	
	AllowPurchases *bool `json:"allow_purchases,omitempty"`
	
	TaxRate *float64 `json:"tax_rate,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Settings *json.RawMessage `json:"settings,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLocationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationCode != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.LocationType != nil {
		hasUpdate = true
	}
	
	if r.Phone != nil {
		hasUpdate = true
	}
	
	if r.Email != nil {
		hasUpdate = true
	}
	
	if r.ManagerUserId != nil {
		hasUpdate = true
	}
	
	if r.AddressLine1 != nil {
		hasUpdate = true
	}
	
	if r.AddressLine2 != nil {
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
	
	if r.Timezone != nil {
		hasUpdate = true
	}
	
	if r.BusinessHours != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsPrimary != nil {
		hasUpdate = true
	}
	
	if r.AllowSales != nil {
		hasUpdate = true
	}
	
	if r.AllowPurchases != nil {
		hasUpdate = true
	}
	
	if r.TaxRate != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
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

// LocationsListResponse represents a paginated list of locations records
type LocationsListResponse struct {
	Items      []*LocationsResponse `json:"items"`
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

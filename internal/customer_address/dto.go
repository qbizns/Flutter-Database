package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomerAddressesResponse represents a customer_addresses response
type CustomerAddressesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerId uuid.UUID `json:"customer_id"`
	
	AddressLabel *string `json:"address_label"`
	
	AddressLine1 string `json:"address_line1"`
	
	AddressLine2 *string `json:"address_line2"`
	
	City *string `json:"city"`
	
	StateProvince *string `json:"state_province"`
	
	PostalCode *string `json:"postal_code"`
	
	Country *string `json:"country"`
	
	Latitude *float64 `json:"latitude"`
	
	Longitude *float64 `json:"longitude"`
	
	LocationNotes *string `json:"location_notes"`
	
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id"`
	
	IsDefault *bool `json:"is_default"`
	
	IsActive *bool `json:"is_active"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCustomerAddressesRequest represents a request to create a customer_addresses
type CreateCustomerAddressesRequest struct {
	
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	
	AddressLabel *string `json:"address_label"`
	
	AddressLine1 string `json:"address_line1" validate:"required"`
	
	AddressLine2 *string `json:"address_line2"`
	
	City *string `json:"city"`
	
	StateProvince *string `json:"state_province"`
	
	PostalCode *string `json:"postal_code"`
	
	Country *string `json:"country"`
	
	Latitude *float64 `json:"latitude"`
	
	Longitude *float64 `json:"longitude"`
	
	LocationNotes *string `json:"location_notes"`
	
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id"`
	
	IsDefault *bool `json:"is_default"`
	
	IsActive *bool `json:"is_active"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateCustomerAddressesRequest) Validate() error {
	
	if r.CustomerId == uuid.Nil {
		return fmt.Errorf("customer_id is required")
	}
	
	if r.AddressLine1 == "" {
		return fmt.Errorf("address_line1 is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCustomerAddressesRequest represents a request to update a customer_addresses
type UpdateCustomerAddressesRequest struct {
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty" validate:"omitempty,required"`
	
	AddressLabel *string `json:"address_label,omitempty"`
	
	AddressLine1 *string `json:"address_line1,omitempty" validate:"omitempty,required"`
	
	AddressLine2 *string `json:"address_line2,omitempty"`
	
	City *string `json:"city,omitempty"`
	
	StateProvince *string `json:"state_province,omitempty"`
	
	PostalCode *string `json:"postal_code,omitempty"`
	
	Country *string `json:"country,omitempty"`
	
	Latitude *float64 `json:"latitude,omitempty"`
	
	Longitude *float64 `json:"longitude,omitempty"`
	
	LocationNotes *string `json:"location_notes,omitempty"`
	
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCustomerAddressesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.AddressLabel != nil {
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
	
	if r.StateProvince != nil {
		hasUpdate = true
	}
	
	if r.PostalCode != nil {
		hasUpdate = true
	}
	
	if r.Country != nil {
		hasUpdate = true
	}
	
	if r.Latitude != nil {
		hasUpdate = true
	}
	
	if r.Longitude != nil {
		hasUpdate = true
	}
	
	if r.LocationNotes != nil {
		hasUpdate = true
	}
	
	if r.DeliveryZoneId != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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

// CustomerAddressesListResponse represents a paginated list of customer_addresses records
type CustomerAddressesListResponse struct {
	Items      []*CustomerAddressesResponse `json:"items"`
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

package delivery_zone

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeliveryZonesResponse represents a delivery_zones response
type DeliveryZonesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	ZoneName string `json:"zone_name"`
	
	ZoneCode *string `json:"zone_code"`
	
	Description *string `json:"description"`
	
	Geofence json.RawMessage `json:"geofence"`
	
	PostalCodes *string `json:"postal_codes"`
	
	CoverageNotes *string `json:"coverage_notes"`
	
	BaseDeliveryFee *float64 `json:"base_delivery_fee"`
	
	FeeType *string `json:"fee_type"`
	
	MinimumOrderAmount *float64 `json:"minimum_order_amount"`
	
	FreeDeliveryThreshold *float64 `json:"free_delivery_threshold"`
	
	EstimatedDeliveryTimeMinutes *int64 `json:"estimated_delivery_time_minutes"`
	
	MaxDeliveryTimeMinutes *int64 `json:"max_delivery_time_minutes"`
	
	Priority *int64 `json:"priority"`
	
	IsActive *bool `json:"is_active"`
	
	ActiveHours json.RawMessage `json:"active_hours"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	BaseDeliveryFee *string `json:"base_delivery_fee"`
	
}

// CreateDeliveryZonesRequest represents a request to create a delivery_zones
type CreateDeliveryZonesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	ZoneName string `json:"zone_name" validate:"required"`
	
	ZoneCode *string `json:"zone_code"`
	
	Description *string `json:"description"`
	
	// Duplicate removed: Geofence json.RawMessage `json:"geofence"`
	
	PostalCodes *string `json:"postal_codes"`
	
	CoverageNotes *string `json:"coverage_notes"`
	
	BaseDeliveryFee *float64 `json:"base_delivery_fee"`
	
	FeeType *string `json:"fee_type"`
	
	MinimumOrderAmount *float64 `json:"minimum_order_amount"`
	
	FreeDeliveryThreshold *float64 `json:"free_delivery_threshold"`
	
	EstimatedDeliveryTimeMinutes *int64 `json:"estimated_delivery_time_minutes"`
	
	MaxDeliveryTimeMinutes *int64 `json:"max_delivery_time_minutes"`
	
	Priority *int64 `json:"priority"`
	
	IsActive *bool `json:"is_active"`
	
	// Duplicate removed: ActiveHours json.RawMessage `json:"active_hours"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	BaseDeliveryFee *string `json:"base_delivery_fee"`
	
}

// Validate validates the create request
func (r *CreateDeliveryZonesRequest) Validate() error {
	
	if r.ZoneName == "" {
		return fmt.Errorf("zone_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDeliveryZonesRequest represents a request to update a delivery_zones
type UpdateDeliveryZonesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	ZoneName *string `json:"zone_name,omitempty" validate:"omitempty,required"`
	
	ZoneCode *string `json:"zone_code,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Geofence *json.RawMessage `json:"geofence,omitempty"`
	
	PostalCodes *string `json:"postal_codes,omitempty"`
	
	CoverageNotes *string `json:"coverage_notes,omitempty"`
	
	BaseDeliveryFee *float64 `json:"base_delivery_fee,omitempty"`
	
	FeeType *string `json:"fee_type,omitempty"`
	
	MinimumOrderAmount *float64 `json:"minimum_order_amount,omitempty"`
	
	FreeDeliveryThreshold *float64 `json:"free_delivery_threshold,omitempty"`
	
	EstimatedDeliveryTimeMinutes *int64 `json:"estimated_delivery_time_minutes,omitempty"`
	
	MaxDeliveryTimeMinutes *int64 `json:"max_delivery_time_minutes,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	ActiveHours *json.RawMessage `json:"active_hours,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	BaseDeliveryFee *string `json:"base_delivery_fee,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDeliveryZonesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.ZoneName != nil {
		hasUpdate = true
	}
	
	if r.ZoneCode != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Geofence != nil {
		hasUpdate = true
	}
	
	if r.PostalCodes != nil {
		hasUpdate = true
	}
	
	if r.CoverageNotes != nil {
		hasUpdate = true
	}
	
	if r.BaseDeliveryFee != nil {
		hasUpdate = true
	}
	
	if r.FeeType != nil {
		hasUpdate = true
	}
	
	if r.MinimumOrderAmount != nil {
		hasUpdate = true
	}
	
	if r.FreeDeliveryThreshold != nil {
		hasUpdate = true
	}
	
	if r.EstimatedDeliveryTimeMinutes != nil {
		hasUpdate = true
	}
	
	if r.MaxDeliveryTimeMinutes != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.ActiveHours != nil {
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
	
	if r.BaseDeliveryFee != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// DeliveryZonesListResponse represents a paginated list of delivery_zones records
type DeliveryZonesListResponse struct {
	Items      []*DeliveryZonesResponse `json:"items"`
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

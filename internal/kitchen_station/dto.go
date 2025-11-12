package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// KitchenStationsResponse represents a kitchen_stations response
type KitchenStationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	StationName string `json:"station_name"`
	
	StationCode string `json:"station_code"`
	
	StationType *string `json:"station_type"`
	
	Description *string `json:"description"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ColorCode *string `json:"color_code"`
	
	PrinterId *uuid.UUID `json:"printer_id"`
	
	IsActive *bool `json:"is_active"`
	
	AutoPrintTickets *bool `json:"auto_print_tickets"`
	
	AlertSoundEnabled *bool `json:"alert_sound_enabled"`
	
	DisplayConfig json.RawMessage `json:"display_config"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateKitchenStationsRequest represents a request to create a kitchen_stations
type CreateKitchenStationsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	StationName string `json:"station_name" validate:"required"`
	
	StationCode string `json:"station_code" validate:"required"`
	
	StationType *string `json:"station_type"`
	
	Description *string `json:"description"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ColorCode *string `json:"color_code"`
	
	PrinterId *uuid.UUID `json:"printer_id"`
	
	IsActive *bool `json:"is_active"`
	
	AutoPrintTickets *bool `json:"auto_print_tickets"`
	
	AlertSoundEnabled *bool `json:"alert_sound_enabled"`
	
	DisplayConfig json.RawMessage `json:"display_config"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateKitchenStationsRequest) Validate() error {
	
	if r.StationName == "" {
		return fmt.Errorf("station_name is required")
	}
	
	if r.StationCode == "" {
		return fmt.Errorf("station_code is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateKitchenStationsRequest represents a request to update a kitchen_stations
type UpdateKitchenStationsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	StationName *string `json:"station_name,omitempty" validate:"omitempty,required"`
	
	StationCode *string `json:"station_code,omitempty" validate:"omitempty,required"`
	
	StationType *string `json:"station_type,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	ColorCode *string `json:"color_code,omitempty"`
	
	PrinterId *uuid.UUID `json:"printer_id,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	AutoPrintTickets *bool `json:"auto_print_tickets,omitempty"`
	
	AlertSoundEnabled *bool `json:"alert_sound_enabled,omitempty"`
	
	DisplayConfig *json.RawMessage `json:"display_config,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateKitchenStationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.StationName != nil {
		hasUpdate = true
	}
	
	if r.StationCode != nil {
		hasUpdate = true
	}
	
	if r.StationType != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.ColorCode != nil {
		hasUpdate = true
	}
	
	if r.PrinterId != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.AutoPrintTickets != nil {
		hasUpdate = true
	}
	
	if r.AlertSoundEnabled != nil {
		hasUpdate = true
	}
	
	if r.DisplayConfig != nil {
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

// KitchenStationsListResponse represents a paginated list of kitchen_stations records
type KitchenStationsListResponse struct {
	Items      []*KitchenStationsResponse `json:"items"`
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

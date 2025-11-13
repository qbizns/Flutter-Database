package floor_plan

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FloorPlansResponse represents a floor_plans response
type FloorPlansResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	FloorName string `json:"floor_name"`
	
	FloorLevel *int64 `json:"floor_level"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	LayoutConfig json.RawMessage `json:"layout_config"`
	
	IsActive *bool `json:"is_active"`
	
	IsDefault *bool `json:"is_default"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateFloorPlansRequest represents a request to create a floor_plans
type CreateFloorPlansRequest struct {
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	FloorName string `json:"floor_name" validate:"required"`
	
	FloorLevel *int64 `json:"floor_level"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	// Duplicate removed: LayoutConfig json.RawMessage `json:"layout_config"`
	
	IsActive *bool `json:"is_active"`
	
	IsDefault *bool `json:"is_default"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateFloorPlansRequest) Validate() error {
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.FloorName == "" {
		return fmt.Errorf("floor_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateFloorPlansRequest represents a request to update a floor_plans
type UpdateFloorPlansRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	FloorName *string `json:"floor_name,omitempty" validate:"omitempty,required"`
	
	FloorLevel *int64 `json:"floor_level,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	LayoutConfig *json.RawMessage `json:"layout_config,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateFloorPlansRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.FloorName != nil {
		hasUpdate = true
	}
	
	if r.FloorLevel != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.LayoutConfig != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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

// FloorPlansListResponse represents a paginated list of floor_plans records
type FloorPlansListResponse struct {
	Items      []*FloorPlansResponse `json:"items"`
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

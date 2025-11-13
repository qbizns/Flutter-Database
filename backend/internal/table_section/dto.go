package table_section

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TableSectionsResponse represents a table_sections response
type TableSectionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	FloorPlanId *uuid.UUID `json:"floor_plan_id"`
	
	SectionName string `json:"section_name"`
	
	SectionCode *string `json:"section_code"`
	
	SectionType *string `json:"section_type"`
	
	ColorCode *string `json:"color_code"`
	
	Icon *string `json:"icon"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateTableSectionsRequest represents a request to create a table_sections
type CreateTableSectionsRequest struct {
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	FloorPlanId *uuid.UUID `json:"floor_plan_id"`
	
	SectionName string `json:"section_name" validate:"required"`
	
	SectionCode *string `json:"section_code"`
	
	SectionType *string `json:"section_type"`
	
	ColorCode *string `json:"color_code"`
	
	Icon *string `json:"icon"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateTableSectionsRequest) Validate() error {
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.SectionName == "" {
		return fmt.Errorf("section_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTableSectionsRequest represents a request to update a table_sections
type UpdateTableSectionsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	FloorPlanId *uuid.UUID `json:"floor_plan_id,omitempty"`
	
	SectionName *string `json:"section_name,omitempty" validate:"omitempty,required"`
	
	SectionCode *string `json:"section_code,omitempty"`
	
	SectionType *string `json:"section_type,omitempty"`
	
	ColorCode *string `json:"color_code,omitempty"`
	
	Icon *string `json:"icon,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTableSectionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.FloorPlanId != nil {
		hasUpdate = true
	}
	
	if r.SectionName != nil {
		hasUpdate = true
	}
	
	if r.SectionCode != nil {
		hasUpdate = true
	}
	
	if r.SectionType != nil {
		hasUpdate = true
	}
	
	if r.ColorCode != nil {
		hasUpdate = true
	}
	
	if r.Icon != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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

// TableSectionsListResponse represents a paginated list of table_sections records
type TableSectionsListResponse struct {
	Items      []*TableSectionsResponse `json:"items"`
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

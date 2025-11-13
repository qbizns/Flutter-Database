package cours

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CoursesResponse represents a courses response
type CoursesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CourseName string `json:"course_name"`
	
	CourseCode *string `json:"course_code"`
	
	CourseType *string `json:"course_type"`
	
	TypicalDurationMinutes *int64 `json:"typical_duration_minutes"`
	
	FireDelayMinutes *int64 `json:"fire_delay_minutes"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ColorCode *string `json:"color_code"`
	
	Icon *string `json:"icon"`
	
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

// CreateCoursesRequest represents a request to create a courses
type CreateCoursesRequest struct {
	
	CourseName string `json:"course_name" validate:"required"`
	
	CourseCode *string `json:"course_code"`
	
	CourseType *string `json:"course_type"`
	
	TypicalDurationMinutes *int64 `json:"typical_duration_minutes"`
	
	FireDelayMinutes *int64 `json:"fire_delay_minutes"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ColorCode *string `json:"color_code"`
	
	Icon *string `json:"icon"`
	
	IsActive *bool `json:"is_active"`
	
	IsDefault *bool `json:"is_default"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateCoursesRequest) Validate() error {
	
	if r.CourseName == "" {
		return fmt.Errorf("course_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCoursesRequest represents a request to update a courses
type UpdateCoursesRequest struct {
	
	CourseName *string `json:"course_name,omitempty" validate:"omitempty,required"`
	
	CourseCode *string `json:"course_code,omitempty"`
	
	CourseType *string `json:"course_type,omitempty"`
	
	TypicalDurationMinutes *int64 `json:"typical_duration_minutes,omitempty"`
	
	FireDelayMinutes *int64 `json:"fire_delay_minutes,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	ColorCode *string `json:"color_code,omitempty"`
	
	Icon *string `json:"icon,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCoursesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CourseName != nil {
		hasUpdate = true
	}
	
	if r.CourseCode != nil {
		hasUpdate = true
	}
	
	if r.CourseType != nil {
		hasUpdate = true
	}
	
	if r.TypicalDurationMinutes != nil {
		hasUpdate = true
	}
	
	if r.FireDelayMinutes != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.ColorCode != nil {
		hasUpdate = true
	}
	
	if r.Icon != nil {
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

// CoursesListResponse represents a paginated list of courses records
type CoursesListResponse struct {
	Items      []*CoursesResponse `json:"items"`
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

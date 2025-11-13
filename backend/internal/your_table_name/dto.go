package your_table_name

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// YourTableNameResponse represents a your_table_name response
type YourTableNameResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	Name string `json:"name"`
	
	Description *string `json:"description"`
	
	Settings json.RawMessage `json:"settings"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateYourTableNameRequest represents a request to create a your_table_name
type CreateYourTableNameRequest struct {
	
	Name string `json:"name" validate:"required"`
	
	Description *string `json:"description"`
	
	// Duplicate removed: Settings json.RawMessage `json:"settings"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateYourTableNameRequest) Validate() error {
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateYourTableNameRequest represents a request to update a your_table_name
type UpdateYourTableNameRequest struct {
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	Settings *json.RawMessage `json:"settings,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateYourTableNameRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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

// YourTableNameListResponse represents a paginated list of your_table_name records
type YourTableNameListResponse struct {
	Items      []*YourTableNameResponse `json:"items"`
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

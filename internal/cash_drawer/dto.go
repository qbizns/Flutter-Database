package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CashDrawersResponse represents a cash_drawers response
type CashDrawersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	DrawerCode string `json:"drawer_code"`
	
	DrawerName string `json:"drawer_name"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCashDrawersRequest represents a request to create a cash_drawers
type CreateCashDrawersRequest struct {
	
	DrawerCode string `json:"drawer_code" validate:"required"`
	
	DrawerName string `json:"drawer_name" validate:"required"`
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateCashDrawersRequest) Validate() error {
	
	if r.DrawerCode == "" {
		return fmt.Errorf("drawer_code is required")
	}
	
	if r.DrawerName == "" {
		return fmt.Errorf("drawer_name is required")
	}
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCashDrawersRequest represents a request to update a cash_drawers
type UpdateCashDrawersRequest struct {
	
	DrawerCode *string `json:"drawer_code,omitempty" validate:"omitempty,required"`
	
	DrawerName *string `json:"drawer_name,omitempty" validate:"omitempty,required"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	DeviceId *uuid.UUID `json:"device_id,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCashDrawersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.DrawerCode != nil {
		hasUpdate = true
	}
	
	if r.DrawerName != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.DeviceId != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CashDrawersListResponse represents a paginated list of cash_drawers records
type CashDrawersListResponse struct {
	Items      []*CashDrawersResponse `json:"items"`
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

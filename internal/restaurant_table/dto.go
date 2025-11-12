package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RestaurantTablesResponse represents a restaurant_tables response
type RestaurantTablesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	FloorPlanId *uuid.UUID `json:"floor_plan_id"`
	
	SectionId *uuid.UUID `json:"section_id"`
	
	TableNumber string `json:"table_number"`
	
	TableName *string `json:"table_name"`
	
	MinCapacity *int64 `json:"min_capacity"`
	
	MaxCapacity int64 `json:"max_capacity"`
	
	TableShape *string `json:"table_shape"`
	
	IsCombinable *bool `json:"is_combinable"`
	
	PositionX *float64 `json:"position_x"`
	
	PositionY *float64 `json:"position_y"`
	
	Rotation *int64 `json:"rotation"`
	
	Status *string `json:"status"`
	
	CurrentCovers *int64 `json:"current_covers"`
	
	SeatedAt *time.Time `json:"seated_at"`
	
	CurrentWaiterId *uuid.UUID `json:"current_waiter_id"`
	
	IsActive *bool `json:"is_active"`
	
	AllowOnlineReservation *bool `json:"allow_online_reservation"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ColorCode *string `json:"color_code"`
	
	Icon *string `json:"icon"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateRestaurantTablesRequest represents a request to create a restaurant_tables
type CreateRestaurantTablesRequest struct {
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	FloorPlanId *uuid.UUID `json:"floor_plan_id"`
	
	SectionId *uuid.UUID `json:"section_id"`
	
	TableNumber string `json:"table_number" validate:"required"`
	
	TableName *string `json:"table_name"`
	
	MinCapacity *int64 `json:"min_capacity"`
	
	MaxCapacity int64 `json:"max_capacity" validate:"required"`
	
	TableShape *string `json:"table_shape"`
	
	IsCombinable *bool `json:"is_combinable"`
	
	PositionX *float64 `json:"position_x"`
	
	PositionY *float64 `json:"position_y"`
	
	Rotation *int64 `json:"rotation"`
	
	Status *string `json:"status"`
	
	CurrentCovers *int64 `json:"current_covers"`
	
	SeatedAt *time.Time `json:"seated_at"`
	
	CurrentWaiterId *uuid.UUID `json:"current_waiter_id"`
	
	IsActive *bool `json:"is_active"`
	
	AllowOnlineReservation *bool `json:"allow_online_reservation"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ColorCode *string `json:"color_code"`
	
	Icon *string `json:"icon"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateRestaurantTablesRequest) Validate() error {
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.TableNumber == "" {
		return fmt.Errorf("table_number is required")
	}
	
	if r.MaxCapacity == 0 {
		return fmt.Errorf("max_capacity is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateRestaurantTablesRequest represents a request to update a restaurant_tables
type UpdateRestaurantTablesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	FloorPlanId *uuid.UUID `json:"floor_plan_id,omitempty"`
	
	SectionId *uuid.UUID `json:"section_id,omitempty"`
	
	TableNumber *string `json:"table_number,omitempty" validate:"omitempty,required"`
	
	TableName *string `json:"table_name,omitempty"`
	
	MinCapacity *int64 `json:"min_capacity,omitempty"`
	
	MaxCapacity *int64 `json:"max_capacity,omitempty" validate:"omitempty,required"`
	
	TableShape *string `json:"table_shape,omitempty"`
	
	IsCombinable *bool `json:"is_combinable,omitempty"`
	
	PositionX *float64 `json:"position_x,omitempty"`
	
	PositionY *float64 `json:"position_y,omitempty"`
	
	Rotation *int64 `json:"rotation,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	CurrentCovers *int64 `json:"current_covers,omitempty"`
	
	SeatedAt *time.Time `json:"seated_at,omitempty"`
	
	CurrentWaiterId *uuid.UUID `json:"current_waiter_id,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	AllowOnlineReservation *bool `json:"allow_online_reservation,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	ColorCode *string `json:"color_code,omitempty"`
	
	Icon *string `json:"icon,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateRestaurantTablesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.FloorPlanId != nil {
		hasUpdate = true
	}
	
	if r.SectionId != nil {
		hasUpdate = true
	}
	
	if r.TableNumber != nil {
		hasUpdate = true
	}
	
	if r.TableName != nil {
		hasUpdate = true
	}
	
	if r.MinCapacity != nil {
		hasUpdate = true
	}
	
	if r.MaxCapacity != nil {
		hasUpdate = true
	}
	
	if r.TableShape != nil {
		hasUpdate = true
	}
	
	if r.IsCombinable != nil {
		hasUpdate = true
	}
	
	if r.PositionX != nil {
		hasUpdate = true
	}
	
	if r.PositionY != nil {
		hasUpdate = true
	}
	
	if r.Rotation != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.CurrentCovers != nil {
		hasUpdate = true
	}
	
	if r.SeatedAt != nil {
		hasUpdate = true
	}
	
	if r.CurrentWaiterId != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.AllowOnlineReservation != nil {
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

// RestaurantTablesListResponse represents a paginated list of restaurant_tables records
type RestaurantTablesListResponse struct {
	Items      []*RestaurantTablesResponse `json:"items"`
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

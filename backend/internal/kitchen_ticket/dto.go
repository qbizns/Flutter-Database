package kitchen_ticket

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// KitchenTicketsResponse represents a kitchen_tickets response
type KitchenTicketsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	TicketNumber string `json:"ticket_number"`
	
	DisplaySequence *int64 `json:"display_sequence"`
	
	OrderId uuid.UUID `json:"order_id"`
	
	KitchenStationId uuid.UUID `json:"kitchen_station_id"`
	
	CourseId *uuid.UUID `json:"course_id"`
	
	TicketType *string `json:"ticket_type"`
	
	Priority *int64 `json:"priority"`
	
	Status *string `json:"status"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	FiredAt *time.Time `json:"fired_at"`
	
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	
	StartedAt *time.Time `json:"started_at"`
	
	ReadyAt *time.Time `json:"ready_at"`
	
	BumpedAt *time.Time `json:"bumped_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	PrepTimeMinutes *int64 `json:"prep_time_minutes"`
	
	TargetPrepTime *int64 `json:"target_prep_time"`
	
	TableNumber *string `json:"table_number"`
	
	OrderType *string `json:"order_type"`
	
	Covers *int64 `json:"covers"`
	
	WaiterName *string `json:"waiter_name"`
	
	SpecialInstructions *string `json:"special_instructions"`
	
	KitchenNotes *string `json:"kitchen_notes"`
	
	DisplayConfig json.RawMessage `json:"display_config"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'new', *string `json:"'new',"`
	
	'completed', *string `json:"'completed',"`
	
	'normal', *string `json:"'normal',"`
	
}

// CreateKitchenTicketsRequest represents a request to create a kitchen_tickets
type CreateKitchenTicketsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	TicketNumber string `json:"ticket_number" validate:"required"`
	
	DisplaySequence *int64 `json:"display_sequence"`
	
	OrderId uuid.UUID `json:"order_id" validate:"required"`
	
	KitchenStationId uuid.UUID `json:"kitchen_station_id" validate:"required"`
	
	CourseId *uuid.UUID `json:"course_id"`
	
	TicketType *string `json:"ticket_type"`
	
	Priority *int64 `json:"priority"`
	
	Status *string `json:"status"`
	
	FiredAt *time.Time `json:"fired_at"`
	
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	
	StartedAt *time.Time `json:"started_at"`
	
	ReadyAt *time.Time `json:"ready_at"`
	
	BumpedAt *time.Time `json:"bumped_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	PrepTimeMinutes *int64 `json:"prep_time_minutes"`
	
	TargetPrepTime *int64 `json:"target_prep_time"`
	
	TableNumber *string `json:"table_number"`
	
	OrderType *string `json:"order_type"`
	
	Covers *int64 `json:"covers"`
	
	WaiterName *string `json:"waiter_name"`
	
	SpecialInstructions *string `json:"special_instructions"`
	
	KitchenNotes *string `json:"kitchen_notes"`
	
	DisplayConfig json.RawMessage `json:"display_config"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'new', *string `json:"'new',"`
	
	'completed', *string `json:"'completed',"`
	
	'normal', *string `json:"'normal',"`
	
}

// Validate validates the create request
func (r *CreateKitchenTicketsRequest) Validate() error {
	
	if r.TicketNumber == "" {
		return fmt.Errorf("ticket_number is required")
	}
	
	if r.OrderId == uuid.Nil {
		return fmt.Errorf("order_id is required")
	}
	
	if r.KitchenStationId == uuid.Nil {
		return fmt.Errorf("kitchen_station_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateKitchenTicketsRequest represents a request to update a kitchen_tickets
type UpdateKitchenTicketsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	TicketNumber *string `json:"ticket_number,omitempty" validate:"omitempty,required"`
	
	DisplaySequence *int64 `json:"display_sequence,omitempty"`
	
	OrderId *uuid.UUID `json:"order_id,omitempty" validate:"omitempty,required"`
	
	KitchenStationId *uuid.UUID `json:"kitchen_station_id,omitempty" validate:"omitempty,required"`
	
	CourseId *uuid.UUID `json:"course_id,omitempty"`
	
	TicketType *string `json:"ticket_type,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	FiredAt *time.Time `json:"fired_at,omitempty"`
	
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	
	StartedAt *time.Time `json:"started_at,omitempty"`
	
	ReadyAt *time.Time `json:"ready_at,omitempty"`
	
	BumpedAt *time.Time `json:"bumped_at,omitempty"`
	
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	PrepTimeMinutes *int64 `json:"prep_time_minutes,omitempty"`
	
	TargetPrepTime *int64 `json:"target_prep_time,omitempty"`
	
	TableNumber *string `json:"table_number,omitempty"`
	
	OrderType *string `json:"order_type,omitempty"`
	
	Covers *int64 `json:"covers,omitempty"`
	
	WaiterName *string `json:"waiter_name,omitempty"`
	
	SpecialInstructions *string `json:"special_instructions,omitempty"`
	
	KitchenNotes *string `json:"kitchen_notes,omitempty"`
	
	DisplayConfig *json.RawMessage `json:"display_config,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'new', *string `json:"'new',,omitempty"`
	
	'completed', *string `json:"'completed',,omitempty"`
	
	'normal', *string `json:"'normal',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateKitchenTicketsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.TicketNumber != nil {
		hasUpdate = true
	}
	
	if r.DisplaySequence != nil {
		hasUpdate = true
	}
	
	if r.OrderId != nil {
		hasUpdate = true
	}
	
	if r.KitchenStationId != nil {
		hasUpdate = true
	}
	
	if r.CourseId != nil {
		hasUpdate = true
	}
	
	if r.TicketType != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.FiredAt != nil {
		hasUpdate = true
	}
	
	if r.AcknowledgedAt != nil {
		hasUpdate = true
	}
	
	if r.StartedAt != nil {
		hasUpdate = true
	}
	
	if r.ReadyAt != nil {
		hasUpdate = true
	}
	
	if r.BumpedAt != nil {
		hasUpdate = true
	}
	
	if r.CompletedAt != nil {
		hasUpdate = true
	}
	
	if r.PrepTimeMinutes != nil {
		hasUpdate = true
	}
	
	if r.TargetPrepTime != nil {
		hasUpdate = true
	}
	
	if r.TableNumber != nil {
		hasUpdate = true
	}
	
	if r.OrderType != nil {
		hasUpdate = true
	}
	
	if r.Covers != nil {
		hasUpdate = true
	}
	
	if r.WaiterName != nil {
		hasUpdate = true
	}
	
	if r.SpecialInstructions != nil {
		hasUpdate = true
	}
	
	if r.KitchenNotes != nil {
		hasUpdate = true
	}
	
	if r.DisplayConfig != nil {
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
	
	if r.'new', != nil {
		hasUpdate = true
	}
	
	if r.'completed', != nil {
		hasUpdate = true
	}
	
	if r.'normal', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// KitchenTicketsListResponse represents a paginated list of kitchen_tickets records
type KitchenTicketsListResponse struct {
	Items      []*KitchenTicketsResponse `json:"items"`
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

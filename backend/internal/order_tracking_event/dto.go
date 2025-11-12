package order_tracking_event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderTrackingEventsResponse represents a order_tracking_events response
type OrderTrackingEventsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	OrderId uuid.UUID `json:"order_id"`
	
	DeliveryAssignmentId *uuid.UUID `json:"delivery_assignment_id"`
	
	EventType string `json:"event_type"`
	
	EventTimestamp *time.Time `json:"event_timestamp"`
	
	EventMessage *string `json:"event_message"`
	
	Location json.RawMessage `json:"location"`
	
	LocationName *string `json:"location_name"`
	
	ActorType *string `json:"actor_type"`
	
	ActorId *uuid.UUID `json:"actor_id"`
	
	ActorName *string `json:"actor_name"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	'orderPlaced', *string `json:"'order_placed',"`
	
	'readyForPickup', *string `json:"'ready_for_pickup',"`
	
	'arrived', *string `json:"'arrived',"`
	
	'rescheduled', *string `json:"'rescheduled',"`
	
	'system', *string `json:"'system',"`
	
}

// CreateOrderTrackingEventsRequest represents a request to create a order_tracking_events
type CreateOrderTrackingEventsRequest struct {
	
	OrderId uuid.UUID `json:"order_id" validate:"required"`
	
	DeliveryAssignmentId *uuid.UUID `json:"delivery_assignment_id"`
	
	EventType string `json:"event_type" validate:"required"`
	
	EventTimestamp *time.Time `json:"event_timestamp"`
	
	EventMessage *string `json:"event_message"`
	
	Location json.RawMessage `json:"location"`
	
	LocationName *string `json:"location_name"`
	
	ActorType *string `json:"actor_type"`
	
	ActorId *uuid.UUID `json:"actor_id"`
	
	ActorName *string `json:"actor_name"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	'orderPlaced', *string `json:"'order_placed',"`
	
	'readyForPickup', *string `json:"'ready_for_pickup',"`
	
	'arrived', *string `json:"'arrived',"`
	
	'rescheduled', *string `json:"'rescheduled',"`
	
	'system', *string `json:"'system',"`
	
}

// Validate validates the create request
func (r *CreateOrderTrackingEventsRequest) Validate() error {
	
	if r.OrderId == uuid.Nil {
		return fmt.Errorf("order_id is required")
	}
	
	if r.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateOrderTrackingEventsRequest represents a request to update a order_tracking_events
type UpdateOrderTrackingEventsRequest struct {
	
	OrderId *uuid.UUID `json:"order_id,omitempty" validate:"omitempty,required"`
	
	DeliveryAssignmentId *uuid.UUID `json:"delivery_assignment_id,omitempty"`
	
	EventType *string `json:"event_type,omitempty" validate:"omitempty,required"`
	
	EventTimestamp *time.Time `json:"event_timestamp,omitempty"`
	
	EventMessage *string `json:"event_message,omitempty"`
	
	Location *json.RawMessage `json:"location,omitempty"`
	
	LocationName *string `json:"location_name,omitempty"`
	
	ActorType *string `json:"actor_type,omitempty"`
	
	ActorId *uuid.UUID `json:"actor_id,omitempty"`
	
	ActorName *string `json:"actor_name,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	'orderPlaced', *string `json:"'order_placed',,omitempty"`
	
	'readyForPickup', *string `json:"'ready_for_pickup',,omitempty"`
	
	'arrived', *string `json:"'arrived',,omitempty"`
	
	'rescheduled', *string `json:"'rescheduled',,omitempty"`
	
	'system', *string `json:"'system',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateOrderTrackingEventsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.OrderId != nil {
		hasUpdate = true
	}
	
	if r.DeliveryAssignmentId != nil {
		hasUpdate = true
	}
	
	if r.EventType != nil {
		hasUpdate = true
	}
	
	if r.EventTimestamp != nil {
		hasUpdate = true
	}
	
	if r.EventMessage != nil {
		hasUpdate = true
	}
	
	if r.Location != nil {
		hasUpdate = true
	}
	
	if r.LocationName != nil {
		hasUpdate = true
	}
	
	if r.ActorType != nil {
		hasUpdate = true
	}
	
	if r.ActorId != nil {
		hasUpdate = true
	}
	
	if r.ActorName != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.'orderPlaced', != nil {
		hasUpdate = true
	}
	
	if r.'readyForPickup', != nil {
		hasUpdate = true
	}
	
	if r.'arrived', != nil {
		hasUpdate = true
	}
	
	if r.'rescheduled', != nil {
		hasUpdate = true
	}
	
	if r.'system', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// OrderTrackingEventsListResponse represents a paginated list of order_tracking_events records
type OrderTrackingEventsListResponse struct {
	Items      []*OrderTrackingEventsResponse `json:"items"`
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

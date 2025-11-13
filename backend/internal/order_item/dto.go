package order_item

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderItemsResponse represents a order_items response
type OrderItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	OrderId uuid.UUID `json:"order_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	ItemName string `json:"item_name"`
	
	Quantity float64 `json:"quantity"`
	
	UnitPrice float64 `json:"unit_price"`
	
	CourseId *uuid.UUID `json:"course_id"`
	
	CoursePosition *int64 `json:"course_position"`
	
	FireTime *time.Time `json:"fire_time"`
	
	KitchenStationId *uuid.UUID `json:"kitchen_station_id"`
	
	KitchenTicketId *uuid.UUID `json:"kitchen_ticket_id"`
	
	Status *string `json:"status"`
	
	FiredAt *time.Time `json:"fired_at"`
	
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	
	StartedPreparingAt *time.Time `json:"started_preparing_at"`
	
	ReadyAt *time.Time `json:"ready_at"`
	
	ServedAt *time.Time `json:"served_at"`
	
	ModifiersTotal *float64 `json:"modifiers_total"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	LineTotal float64 `json:"line_total"`
	
	SpecialInstructions *string `json:"special_instructions"`
	
	CustomerNotes *string `json:"customer_notes"`
	
	KitchenNotes *string `json:"kitchen_notes"`
	
	SeatNumber *int64 `json:"seat_number"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'pending', *string `json:"'pending',"`
	
	'served', *string `json:"'served',"`
	
	UnitPrice *string `json:"unit_price"`
	
	DiscountAmount *string `json:"discount_amount"`
	
}

// CreateOrderItemsRequest represents a request to create a order_items
type CreateOrderItemsRequest struct {
	
	OrderId uuid.UUID `json:"order_id" validate:"required"`
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	ItemName string `json:"item_name" validate:"required"`
	
	Quantity float64 `json:"quantity" validate:"required"`
	
	UnitPrice float64 `json:"unit_price" validate:"required"`
	
	CourseId *uuid.UUID `json:"course_id"`
	
	CoursePosition *int64 `json:"course_position"`
	
	FireTime *time.Time `json:"fire_time"`
	
	KitchenStationId *uuid.UUID `json:"kitchen_station_id"`
	
	KitchenTicketId *uuid.UUID `json:"kitchen_ticket_id"`
	
	// 	Status *string `json:"status"`
	
	FiredAt *time.Time `json:"fired_at"`
	
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	
	StartedPreparingAt *time.Time `json:"started_preparing_at"`
	
	ReadyAt *time.Time `json:"ready_at"`
	
	ServedAt *time.Time `json:"served_at"`
	
	ModifiersTotal *float64 `json:"modifiers_total"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	LineTotal float64 `json:"line_total" validate:"required"`
	
	SpecialInstructions *string `json:"special_instructions"`
	
	CustomerNotes *string `json:"customer_notes"`
	
	KitchenNotes *string `json:"kitchen_notes"`
	
	SeatNumber *int64 `json:"seat_number"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'pending', *string `json:"'pending',"`
	
	'served', *string `json:"'served',"`
	
	UnitPrice *string `json:"unit_price"`
	
	DiscountAmount *string `json:"discount_amount"`
	
}

// Validate validates the create request
func (r *CreateOrderItemsRequest) Validate() error {
	
	if r.OrderId == uuid.Nil {
		return fmt.Errorf("order_id is required")
	}
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.ItemName == "" {
		return fmt.Errorf("item_name is required")
	}
	
	if r.Quantity == nil {
		return fmt.Errorf("quantity is required")
	}
	
	if r.UnitPrice == nil {
		return fmt.Errorf("unit_price is required")
	}
	
	if r.LineTotal == nil {
		return fmt.Errorf("line_total is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateOrderItemsRequest represents a request to update a order_items
type UpdateOrderItemsRequest struct {
	
	OrderId *uuid.UUID `json:"order_id,omitempty" validate:"omitempty,required"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	ItemName *string `json:"item_name,omitempty" validate:"omitempty,required"`
	
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,required"`
	
	UnitPrice *float64 `json:"unit_price,omitempty" validate:"omitempty,required"`
	
	CourseId *uuid.UUID `json:"course_id,omitempty"`
	
	CoursePosition *int64 `json:"course_position,omitempty"`
	
	FireTime *time.Time `json:"fire_time,omitempty"`
	
	KitchenStationId *uuid.UUID `json:"kitchen_station_id,omitempty"`
	
	KitchenTicketId *uuid.UUID `json:"kitchen_ticket_id,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	FiredAt *time.Time `json:"fired_at,omitempty"`
	
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	
	StartedPreparingAt *time.Time `json:"started_preparing_at,omitempty"`
	
	ReadyAt *time.Time `json:"ready_at,omitempty"`
	
	ServedAt *time.Time `json:"served_at,omitempty"`
	
	ModifiersTotal *float64 `json:"modifiers_total,omitempty"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	
	LineTotal *float64 `json:"line_total,omitempty" validate:"omitempty,required"`
	
	SpecialInstructions *string `json:"special_instructions,omitempty"`
	
	CustomerNotes *string `json:"customer_notes,omitempty"`
	
	KitchenNotes *string `json:"kitchen_notes,omitempty"`
	
	SeatNumber *int64 `json:"seat_number,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'pending', *string `json:"'pending',,omitempty"`
	
	'served', *string `json:"'served',,omitempty"`
	
	UnitPrice *string `json:"unit_price,omitempty"`
	
	DiscountAmount *string `json:"discount_amount,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateOrderItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.OrderId != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.ItemName != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.UnitPrice != nil {
		hasUpdate = true
	}
	
	if r.CourseId != nil {
		hasUpdate = true
	}
	
	if r.CoursePosition != nil {
		hasUpdate = true
	}
	
	if r.FireTime != nil {
		hasUpdate = true
	}
	
	if r.KitchenStationId != nil {
		hasUpdate = true
	}
	
	if r.KitchenTicketId != nil {
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
	
	if r.StartedPreparingAt != nil {
		hasUpdate = true
	}
	
	if r.ReadyAt != nil {
		hasUpdate = true
	}
	
	if r.ServedAt != nil {
		hasUpdate = true
	}
	
	if r.ModifiersTotal != nil {
		hasUpdate = true
	}
	
	if r.DiscountAmount != nil {
		hasUpdate = true
	}
	
	if r.LineTotal != nil {
		hasUpdate = true
	}
	
	if r.SpecialInstructions != nil {
		hasUpdate = true
	}
	
	if r.CustomerNotes != nil {
		hasUpdate = true
	}
	
	if r.KitchenNotes != nil {
		hasUpdate = true
	}
	
	if r.SeatNumber != nil {
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
	
	if r.'pending', != nil {
		hasUpdate = true
	}
	
	if r.'served', != nil {
		hasUpdate = true
	}
	
	if r.UnitPrice != nil {
		hasUpdate = true
	}
	
	if r.DiscountAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// OrderItemsListResponse represents a paginated list of order_items records
type OrderItemsListResponse struct {
	Items      []*OrderItemsResponse `json:"items"`
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

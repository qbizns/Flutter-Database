package order

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrdersResponse represents a orders response
type OrdersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	OrderNumber string `json:"order_number"`
	
	DisplayNumber *int64 `json:"display_number"`
	
	OrderType *string `json:"order_type"`
	
	TableId *uuid.UUID `json:"table_id"`
	
	ReservationId *uuid.UUID `json:"reservation_id"`
	
	Covers *int64 `json:"covers"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	WaiterId *uuid.UUID `json:"waiter_id"`
	
	Status *string `json:"status"`
	
	OrderDate *time.Time `json:"order_date"`
	
	SubmittedAt *time.Time `json:"submitted_at"`
	
	KitchenReceivedAt *time.Time `json:"kitchen_received_at"`
	
	ReadyAt *time.Time `json:"ready_at"`
	
	ServedAt *time.Time `json:"served_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	Subtotal *float64 `json:"subtotal"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	ServiceCharge *float64 `json:"service_charge"`
	
	TotalAmount *float64 `json:"total_amount"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	ShiftId *uuid.UUID `json:"shift_id"`
	
	CustomerNotes *string `json:"customer_notes"`
	
	KitchenNotes *string `json:"kitchen_notes"`
	
	InternalNotes *string `json:"internal_notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`

}

// CreateOrdersRequest represents a request to create a orders
type CreateOrdersRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	OrderNumber string `json:"order_number" validate:"required"`
	
	DisplayNumber *int64 `json:"display_number"`
	
	OrderType *string `json:"order_type"`
	
	TableId *uuid.UUID `json:"table_id"`
	
	ReservationId *uuid.UUID `json:"reservation_id"`
	
	Covers *int64 `json:"covers"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	WaiterId *uuid.UUID `json:"waiter_id"`
	
	Status *string `json:"status"`
	
	OrderDate *time.Time `json:"order_date"`
	
	SubmittedAt *time.Time `json:"submitted_at"`
	
	KitchenReceivedAt *time.Time `json:"kitchen_received_at"`
	
	ReadyAt *time.Time `json:"ready_at"`
	
	ServedAt *time.Time `json:"served_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	Subtotal *float64 `json:"subtotal"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	ServiceCharge *float64 `json:"service_charge"`
	
	TotalAmount *float64 `json:"total_amount"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	ShiftId *uuid.UUID `json:"shift_id"`
	
	CustomerNotes *string `json:"customer_notes"`
	
	KitchenNotes *string `json:"kitchen_notes"`
	
	InternalNotes *string `json:"internal_notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`

}

// Validate validates the create request
func (r *CreateOrdersRequest) Validate() error {
	
	if r.OrderNumber == "" {
		return fmt.Errorf("order_number is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateOrdersRequest represents a request to update a orders
type UpdateOrdersRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	OrderNumber *string `json:"order_number,omitempty" validate:"omitempty,required"`
	
	DisplayNumber *int64 `json:"display_number,omitempty"`
	
	OrderType *string `json:"order_type,omitempty"`
	
	TableId *uuid.UUID `json:"table_id,omitempty"`
	
	ReservationId *uuid.UUID `json:"reservation_id,omitempty"`
	
	Covers *int64 `json:"covers,omitempty"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty"`
	
	WaiterId *uuid.UUID `json:"waiter_id,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	OrderDate *time.Time `json:"order_date,omitempty"`
	
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	
	KitchenReceivedAt *time.Time `json:"kitchen_received_at,omitempty"`
	
	ReadyAt *time.Time `json:"ready_at,omitempty"`
	
	ServedAt *time.Time `json:"served_at,omitempty"`
	
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	Subtotal *float64 `json:"subtotal,omitempty"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	
	ServiceCharge *float64 `json:"service_charge,omitempty"`
	
	TotalAmount *float64 `json:"total_amount,omitempty"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	ShiftId *uuid.UUID `json:"shift_id,omitempty"`
	
	CustomerNotes *string `json:"customer_notes,omitempty"`
	
	KitchenNotes *string `json:"kitchen_notes,omitempty"`
	
	InternalNotes *string `json:"internal_notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`

}

// Validate validates the update request
func (r *UpdateOrdersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.OrderNumber != nil {
		hasUpdate = true
	}
	
	if r.DisplayNumber != nil {
		hasUpdate = true
	}
	
	if r.OrderType != nil {
		hasUpdate = true
	}
	
	if r.TableId != nil {
		hasUpdate = true
	}
	
	if r.ReservationId != nil {
		hasUpdate = true
	}
	
	if r.Covers != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.WaiterId != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.OrderDate != nil {
		hasUpdate = true
	}
	
	if r.SubmittedAt != nil {
		hasUpdate = true
	}
	
	if r.KitchenReceivedAt != nil {
		hasUpdate = true
	}
	
	if r.ReadyAt != nil {
		hasUpdate = true
	}
	
	if r.ServedAt != nil {
		hasUpdate = true
	}
	
	if r.CompletedAt != nil {
		hasUpdate = true
	}
	
	if r.Subtotal != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.DiscountAmount != nil {
		hasUpdate = true
	}
	
	if r.ServiceCharge != nil {
		hasUpdate = true
	}
	
	if r.TotalAmount != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.ShiftId != nil {
		hasUpdate = true
	}
	
	if r.CustomerNotes != nil {
		hasUpdate = true
	}
	
	if r.KitchenNotes != nil {
		hasUpdate = true
	}
	
	if r.InternalNotes != nil {
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

// OrdersListResponse represents a paginated list of orders records
type OrdersListResponse struct {
	Items      []*OrdersResponse `json:"items"`
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

// UpdateOrderStatusRequest represents a request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

// Validate validates the update status request
func (r *UpdateOrderStatusRequest) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	return nil
}

// CancelOrderRequest represents a request to cancel an order
type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

// OrderStatisticsResponse represents order statistics
type OrderStatisticsResponse struct {
	TotalOrders               int     `json:"total_orders"`
	CompletedOrders           int     `json:"completed_orders"`
	CancelledOrders           int     `json:"cancelled_orders"`
	ActiveOrders              int     `json:"active_orders"`
	TotalRevenue              float64 `json:"total_revenue"`
	AverageOrderValue         float64 `json:"average_order_value"`
	AveragePreparationMinutes int     `json:"average_preparation_minutes"`
}

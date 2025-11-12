package delivery_assignment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeliveryAssignmentsResponse represents a delivery_assignments response
type DeliveryAssignmentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	OrderId uuid.UUID `json:"order_id"`
	
	DriverId uuid.UUID `json:"driver_id"`
	
	DriverShiftId *uuid.UUID `json:"driver_shift_id"`
	
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id"`
	
	CustomerAddressId *uuid.UUID `json:"customer_address_id"`
	
	DeliveryAddress string `json:"delivery_address"`
	
	DeliveryLocation json.RawMessage `json:"delivery_location"`
	
	AssignedAt *time.Time `json:"assigned_at"`
	
	AssignedBy *uuid.UUID `json:"assigned_by"`
	
	Status *string `json:"status"`
	
	AcceptedAt *time.Time `json:"accepted_at"`
	
	PickedUpAt *time.Time `json:"picked_up_at"`
	
	DispatchedAt *time.Time `json:"dispatched_at"`
	
	ArrivedAt *time.Time `json:"arrived_at"`
	
	DeliveredAt *time.Time `json:"delivered_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
	EstimatedPickupTime *time.Time `json:"estimated_pickup_time"`
	
	EstimatedDeliveryTime *time.Time `json:"estimated_delivery_time"`
	
	DistanceKm *float64 `json:"distance_km"`
	
	RouteInfo json.RawMessage `json:"route_info"`
	
	DeliveryFee *float64 `json:"delivery_fee"`
	
	DriverCommission *float64 `json:"driver_commission"`
	
	PaymentMethod *string `json:"payment_method"`
	
	CashCollected *float64 `json:"cash_collected"`
	
	SignatureImageUrl *string `json:"signature_image_url"`
	
	DeliveryPhotoUrl *string `json:"delivery_photo_url"`
	
	RecipientName *string `json:"recipient_name"`
	
	DeliveryNotes *string `json:"delivery_notes"`
	
	FailureReason *string `json:"failure_reason"`
	
	FailureNotes *string `json:"failure_notes"`
	
	RetryCount *int64 `json:"retry_count"`
	
	CustomerRating *int64 `json:"customer_rating"`
	
	CustomerFeedback *string `json:"customer_feedback"`
	
	DriverNotes *string `json:"driver_notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'assigned', *string `json:"'assigned',"`
	
	'delivered', *string `json:"'delivered',"`
	
	CustomerRating *string `json:"customer_rating"`
	
	DeliveryFee *string `json:"delivery_fee"`
	
}

// CreateDeliveryAssignmentsRequest represents a request to create a delivery_assignments
type CreateDeliveryAssignmentsRequest struct {
	
	OrderId uuid.UUID `json:"order_id" validate:"required"`
	
	DriverId uuid.UUID `json:"driver_id" validate:"required"`
	
	DriverShiftId *uuid.UUID `json:"driver_shift_id"`
	
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id"`
	
	CustomerAddressId *uuid.UUID `json:"customer_address_id"`
	
	DeliveryAddress string `json:"delivery_address" validate:"required"`
	
	DeliveryLocation json.RawMessage `json:"delivery_location"`
	
	AssignedAt *time.Time `json:"assigned_at"`
	
	AssignedBy *uuid.UUID `json:"assigned_by"`
	
	Status *string `json:"status"`
	
	AcceptedAt *time.Time `json:"accepted_at"`
	
	PickedUpAt *time.Time `json:"picked_up_at"`
	
	DispatchedAt *time.Time `json:"dispatched_at"`
	
	ArrivedAt *time.Time `json:"arrived_at"`
	
	DeliveredAt *time.Time `json:"delivered_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
	EstimatedPickupTime *time.Time `json:"estimated_pickup_time"`
	
	EstimatedDeliveryTime *time.Time `json:"estimated_delivery_time"`
	
	DistanceKm *float64 `json:"distance_km"`
	
	RouteInfo json.RawMessage `json:"route_info"`
	
	DeliveryFee *float64 `json:"delivery_fee"`
	
	DriverCommission *float64 `json:"driver_commission"`
	
	PaymentMethod *string `json:"payment_method"`
	
	CashCollected *float64 `json:"cash_collected"`
	
	SignatureImageUrl *string `json:"signature_image_url" validate:"url"`
	
	DeliveryPhotoUrl *string `json:"delivery_photo_url" validate:"url"`
	
	RecipientName *string `json:"recipient_name"`
	
	DeliveryNotes *string `json:"delivery_notes"`
	
	FailureReason *string `json:"failure_reason"`
	
	FailureNotes *string `json:"failure_notes"`
	
	RetryCount *int64 `json:"retry_count"`
	
	CustomerRating *int64 `json:"customer_rating"`
	
	CustomerFeedback *string `json:"customer_feedback"`
	
	DriverNotes *string `json:"driver_notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'assigned', *string `json:"'assigned',"`
	
	'delivered', *string `json:"'delivered',"`
	
	CustomerRating *string `json:"customer_rating"`
	
	DeliveryFee *string `json:"delivery_fee"`
	
}

// Validate validates the create request
func (r *CreateDeliveryAssignmentsRequest) Validate() error {
	
	if r.OrderId == uuid.Nil {
		return fmt.Errorf("order_id is required")
	}
	
	if r.DriverId == uuid.Nil {
		return fmt.Errorf("driver_id is required")
	}
	
	if r.DeliveryAddress == "" {
		return fmt.Errorf("delivery_address is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDeliveryAssignmentsRequest represents a request to update a delivery_assignments
type UpdateDeliveryAssignmentsRequest struct {
	
	OrderId *uuid.UUID `json:"order_id,omitempty" validate:"omitempty,required"`
	
	DriverId *uuid.UUID `json:"driver_id,omitempty" validate:"omitempty,required"`
	
	DriverShiftId *uuid.UUID `json:"driver_shift_id,omitempty"`
	
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id,omitempty"`
	
	CustomerAddressId *uuid.UUID `json:"customer_address_id,omitempty"`
	
	DeliveryAddress *string `json:"delivery_address,omitempty" validate:"omitempty,required"`
	
	DeliveryLocation *json.RawMessage `json:"delivery_location,omitempty"`
	
	AssignedAt *time.Time `json:"assigned_at,omitempty"`
	
	AssignedBy *uuid.UUID `json:"assigned_by,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	
	PickedUpAt *time.Time `json:"picked_up_at,omitempty"`
	
	DispatchedAt *time.Time `json:"dispatched_at,omitempty"`
	
	ArrivedAt *time.Time `json:"arrived_at,omitempty"`
	
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	
	FailedAt *time.Time `json:"failed_at,omitempty"`
	
	EstimatedPickupTime *time.Time `json:"estimated_pickup_time,omitempty"`
	
	EstimatedDeliveryTime *time.Time `json:"estimated_delivery_time,omitempty"`
	
	DistanceKm *float64 `json:"distance_km,omitempty"`
	
	RouteInfo *json.RawMessage `json:"route_info,omitempty"`
	
	DeliveryFee *float64 `json:"delivery_fee,omitempty"`
	
	DriverCommission *float64 `json:"driver_commission,omitempty"`
	
	PaymentMethod *string `json:"payment_method,omitempty"`
	
	CashCollected *float64 `json:"cash_collected,omitempty"`
	
	SignatureImageUrl *string `json:"signature_image_url,omitempty" validate:"omitempty,url"`
	
	DeliveryPhotoUrl *string `json:"delivery_photo_url,omitempty" validate:"omitempty,url"`
	
	RecipientName *string `json:"recipient_name,omitempty"`
	
	DeliveryNotes *string `json:"delivery_notes,omitempty"`
	
	FailureReason *string `json:"failure_reason,omitempty"`
	
	FailureNotes *string `json:"failure_notes,omitempty"`
	
	RetryCount *int64 `json:"retry_count,omitempty"`
	
	CustomerRating *int64 `json:"customer_rating,omitempty"`
	
	CustomerFeedback *string `json:"customer_feedback,omitempty"`
	
	DriverNotes *string `json:"driver_notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'assigned', *string `json:"'assigned',,omitempty"`
	
	'delivered', *string `json:"'delivered',,omitempty"`
	
	CustomerRating *string `json:"customer_rating,omitempty"`
	
	DeliveryFee *string `json:"delivery_fee,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDeliveryAssignmentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.OrderId != nil {
		hasUpdate = true
	}
	
	if r.DriverId != nil {
		hasUpdate = true
	}
	
	if r.DriverShiftId != nil {
		hasUpdate = true
	}
	
	if r.DeliveryZoneId != nil {
		hasUpdate = true
	}
	
	if r.CustomerAddressId != nil {
		hasUpdate = true
	}
	
	if r.DeliveryAddress != nil {
		hasUpdate = true
	}
	
	if r.DeliveryLocation != nil {
		hasUpdate = true
	}
	
	if r.AssignedAt != nil {
		hasUpdate = true
	}
	
	if r.AssignedBy != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.AcceptedAt != nil {
		hasUpdate = true
	}
	
	if r.PickedUpAt != nil {
		hasUpdate = true
	}
	
	if r.DispatchedAt != nil {
		hasUpdate = true
	}
	
	if r.ArrivedAt != nil {
		hasUpdate = true
	}
	
	if r.DeliveredAt != nil {
		hasUpdate = true
	}
	
	if r.FailedAt != nil {
		hasUpdate = true
	}
	
	if r.EstimatedPickupTime != nil {
		hasUpdate = true
	}
	
	if r.EstimatedDeliveryTime != nil {
		hasUpdate = true
	}
	
	if r.DistanceKm != nil {
		hasUpdate = true
	}
	
	if r.RouteInfo != nil {
		hasUpdate = true
	}
	
	if r.DeliveryFee != nil {
		hasUpdate = true
	}
	
	if r.DriverCommission != nil {
		hasUpdate = true
	}
	
	if r.PaymentMethod != nil {
		hasUpdate = true
	}
	
	if r.CashCollected != nil {
		hasUpdate = true
	}
	
	if r.SignatureImageUrl != nil {
		hasUpdate = true
	}
	
	if r.DeliveryPhotoUrl != nil {
		hasUpdate = true
	}
	
	if r.RecipientName != nil {
		hasUpdate = true
	}
	
	if r.DeliveryNotes != nil {
		hasUpdate = true
	}
	
	if r.FailureReason != nil {
		hasUpdate = true
	}
	
	if r.FailureNotes != nil {
		hasUpdate = true
	}
	
	if r.RetryCount != nil {
		hasUpdate = true
	}
	
	if r.CustomerRating != nil {
		hasUpdate = true
	}
	
	if r.CustomerFeedback != nil {
		hasUpdate = true
	}
	
	if r.DriverNotes != nil {
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
	
	if r.'assigned', != nil {
		hasUpdate = true
	}
	
	if r.'delivered', != nil {
		hasUpdate = true
	}
	
	if r.CustomerRating != nil {
		hasUpdate = true
	}
	
	if r.DeliveryFee != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// DeliveryAssignmentsListResponse represents a paginated list of delivery_assignments records
type DeliveryAssignmentsListResponse struct {
	Items      []*DeliveryAssignmentsResponse `json:"items"`
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

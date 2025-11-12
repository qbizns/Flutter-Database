package inventory_transfer

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InventoryTransfersResponse represents a inventory_transfers response
type InventoryTransfersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	TransferNumber string `json:"transfer_number"`
	
	TransferDate time.Time `json:"transfer_date"`
	
	FromLocationId uuid.UUID `json:"from_location_id"`
	
	ToLocationId uuid.UUID `json:"to_location_id"`
	
	Status *string `json:"status"`
	
	RequestedDate *time.Time `json:"requested_date"`
	
	ApprovedDate *time.Time `json:"approved_date"`
	
	ShippedDate *time.Time `json:"shipped_date"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date"`
	
	ReceivedDate *time.Time `json:"received_date"`
	
	Carrier *string `json:"carrier"`
	
	TrackingNumber *string `json:"tracking_number"`
	
	ShippingCost *float64 `json:"shipping_cost"`
	
	Reason *string `json:"reason"`
	
	Notes *string `json:"notes"`
	
	RejectionReason *string `json:"rejection_reason"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	RequestedBy *uuid.UUID `json:"requested_by"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ShippedBy *uuid.UUID `json:"shipped_by"`
	
	ReceivedBy *uuid.UUID `json:"received_by"`
	
	(approvedDate *string `json:"(approved_date"`
	
	(shippedDate *string `json:"(shipped_date"`
	
	(receivedDate *string `json:"(received_date"`
	
}

// CreateInventoryTransfersRequest represents a request to create a inventory_transfers
type CreateInventoryTransfersRequest struct {
	
	TransferNumber string `json:"transfer_number" validate:"required"`
	
	TransferDate time.Time `json:"transfer_date" validate:"required"`
	
	FromLocationId uuid.UUID `json:"from_location_id" validate:"required"`
	
	ToLocationId uuid.UUID `json:"to_location_id" validate:"required"`
	
	Status *string `json:"status"`
	
	RequestedDate *time.Time `json:"requested_date"`
	
	ApprovedDate *time.Time `json:"approved_date"`
	
	ShippedDate *time.Time `json:"shipped_date"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date"`
	
	ReceivedDate *time.Time `json:"received_date"`
	
	Carrier *string `json:"carrier"`
	
	TrackingNumber *string `json:"tracking_number"`
	
	ShippingCost *float64 `json:"shipping_cost"`
	
	Reason *string `json:"reason"`
	
	Notes *string `json:"notes"`
	
	RejectionReason *string `json:"rejection_reason"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	RequestedBy *uuid.UUID `json:"requested_by"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ShippedBy *uuid.UUID `json:"shipped_by"`
	
	ReceivedBy *uuid.UUID `json:"received_by"`
	
	(approvedDate *string `json:"(approved_date"`
	
	(shippedDate *string `json:"(shipped_date"`
	
	(receivedDate *string `json:"(received_date"`
	
}

// Validate validates the create request
func (r *CreateInventoryTransfersRequest) Validate() error {
	
	if r.TransferNumber == "" {
		return fmt.Errorf("transfer_number is required")
	}
	
	if r.TransferDate == nil {
		return fmt.Errorf("transfer_date is required")
	}
	
	if r.FromLocationId == uuid.Nil {
		return fmt.Errorf("from_location_id is required")
	}
	
	if r.ToLocationId == uuid.Nil {
		return fmt.Errorf("to_location_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateInventoryTransfersRequest represents a request to update a inventory_transfers
type UpdateInventoryTransfersRequest struct {
	
	TransferNumber *string `json:"transfer_number,omitempty" validate:"omitempty,required"`
	
	TransferDate *time.Time `json:"transfer_date,omitempty" validate:"omitempty,required"`
	
	FromLocationId *uuid.UUID `json:"from_location_id,omitempty" validate:"omitempty,required"`
	
	ToLocationId *uuid.UUID `json:"to_location_id,omitempty" validate:"omitempty,required"`
	
	Status *string `json:"status,omitempty"`
	
	RequestedDate *time.Time `json:"requested_date,omitempty"`
	
	ApprovedDate *time.Time `json:"approved_date,omitempty"`
	
	ShippedDate *time.Time `json:"shipped_date,omitempty"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty"`
	
	ReceivedDate *time.Time `json:"received_date,omitempty"`
	
	Carrier *string `json:"carrier,omitempty"`
	
	TrackingNumber *string `json:"tracking_number,omitempty"`
	
	ShippingCost *float64 `json:"shipping_cost,omitempty"`
	
	Reason *string `json:"reason,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	RejectionReason *string `json:"rejection_reason,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	RequestedBy *uuid.UUID `json:"requested_by,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ShippedBy *uuid.UUID `json:"shipped_by,omitempty"`
	
	ReceivedBy *uuid.UUID `json:"received_by,omitempty"`
	
	(approvedDate *string `json:"(approved_date,omitempty"`
	
	(shippedDate *string `json:"(shipped_date,omitempty"`
	
	(receivedDate *string `json:"(received_date,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateInventoryTransfersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TransferNumber != nil {
		hasUpdate = true
	}
	
	if r.TransferDate != nil {
		hasUpdate = true
	}
	
	if r.FromLocationId != nil {
		hasUpdate = true
	}
	
	if r.ToLocationId != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.RequestedDate != nil {
		hasUpdate = true
	}
	
	if r.ApprovedDate != nil {
		hasUpdate = true
	}
	
	if r.ShippedDate != nil {
		hasUpdate = true
	}
	
	if r.ExpectedDeliveryDate != nil {
		hasUpdate = true
	}
	
	if r.ReceivedDate != nil {
		hasUpdate = true
	}
	
	if r.Carrier != nil {
		hasUpdate = true
	}
	
	if r.TrackingNumber != nil {
		hasUpdate = true
	}
	
	if r.ShippingCost != nil {
		hasUpdate = true
	}
	
	if r.Reason != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.RejectionReason != nil {
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
	
	if r.RequestedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ShippedBy != nil {
		hasUpdate = true
	}
	
	if r.ReceivedBy != nil {
		hasUpdate = true
	}
	
	if r.(approvedDate != nil {
		hasUpdate = true
	}
	
	if r.(shippedDate != nil {
		hasUpdate = true
	}
	
	if r.(receivedDate != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// InventoryTransfersListResponse represents a paginated list of inventory_transfers records
type InventoryTransfersListResponse struct {
	Items      []*InventoryTransfersResponse `json:"items"`
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

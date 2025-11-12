package goods_receipt

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GoodsReceiptsResponse represents a goods_receipts response
type GoodsReceiptsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ReceiptNumber string `json:"receipt_number"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	SupplierId uuid.UUID `json:"supplier_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	ReceiptDate time.Time `json:"receipt_date"`
	
	ReceivedBy uuid.UUID `json:"received_by"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateGoodsReceiptsRequest represents a request to create a goods_receipts
type CreateGoodsReceiptsRequest struct {
	
	ReceiptNumber string `json:"receipt_number" validate:"required"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	ReceiptDate time.Time `json:"receipt_date" validate:"required"`
	
	ReceivedBy uuid.UUID `json:"received_by" validate:"required"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateGoodsReceiptsRequest) Validate() error {
	
	if r.ReceiptNumber == "" {
		return fmt.Errorf("receipt_number is required")
	}
	
	if r.SupplierId == uuid.Nil {
		return fmt.Errorf("supplier_id is required")
	}
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.ReceiptDate == nil {
		return fmt.Errorf("receipt_date is required")
	}
	
	if r.ReceivedBy == uuid.Nil {
		return fmt.Errorf("received_by is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateGoodsReceiptsRequest represents a request to update a goods_receipts
type UpdateGoodsReceiptsRequest struct {
	
	ReceiptNumber *string `json:"receipt_number,omitempty" validate:"omitempty,required"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id,omitempty"`
	
	SupplierId *uuid.UUID `json:"supplier_id,omitempty" validate:"omitempty,required"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	ReceiptDate *time.Time `json:"receipt_date,omitempty" validate:"omitempty,required"`
	
	ReceivedBy *uuid.UUID `json:"received_by,omitempty" validate:"omitempty,required"`
	
	Status *string `json:"status,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateGoodsReceiptsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ReceiptNumber != nil {
		hasUpdate = true
	}
	
	if r.PurchaseOrderId != nil {
		hasUpdate = true
	}
	
	if r.SupplierId != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.ReceiptDate != nil {
		hasUpdate = true
	}
	
	if r.ReceivedBy != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
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

// GoodsReceiptsListResponse represents a paginated list of goods_receipts records
type GoodsReceiptsListResponse struct {
	Items      []*GoodsReceiptsResponse `json:"items"`
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

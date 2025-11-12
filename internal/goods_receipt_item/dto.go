package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GoodsReceiptItemsResponse represents a goods_receipt_items response
type GoodsReceiptItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	GoodsReceiptId uuid.UUID `json:"goods_receipt_id"`
	
	PurchaseOrderItemId *uuid.UUID `json:"purchase_order_item_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	QuantityReceived float64 `json:"quantity_received"`
	
	QuantityAccepted *float64 `json:"quantity_accepted"`
	
	QuantityRejected *float64 `json:"quantity_rejected"`
	
	RejectionReason *string `json:"rejection_reason"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateGoodsReceiptItemsRequest represents a request to create a goods_receipt_items
type CreateGoodsReceiptItemsRequest struct {
	
	GoodsReceiptId uuid.UUID `json:"goods_receipt_id" validate:"required"`
	
	PurchaseOrderItemId *uuid.UUID `json:"purchase_order_item_id"`
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	QuantityReceived float64 `json:"quantity_received" validate:"required"`
	
	QuantityAccepted *float64 `json:"quantity_accepted"`
	
	QuantityRejected *float64 `json:"quantity_rejected"`
	
	RejectionReason *string `json:"rejection_reason"`
	
}

// Validate validates the create request
func (r *CreateGoodsReceiptItemsRequest) Validate() error {
	
	if r.GoodsReceiptId == uuid.Nil {
		return fmt.Errorf("goods_receipt_id is required")
	}
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.QuantityReceived == nil {
		return fmt.Errorf("quantity_received is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateGoodsReceiptItemsRequest represents a request to update a goods_receipt_items
type UpdateGoodsReceiptItemsRequest struct {
	
	GoodsReceiptId *uuid.UUID `json:"goods_receipt_id,omitempty" validate:"omitempty,required"`
	
	PurchaseOrderItemId *uuid.UUID `json:"purchase_order_item_id,omitempty"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	QuantityReceived *float64 `json:"quantity_received,omitempty" validate:"omitempty,required"`
	
	QuantityAccepted *float64 `json:"quantity_accepted,omitempty"`
	
	QuantityRejected *float64 `json:"quantity_rejected,omitempty"`
	
	RejectionReason *string `json:"rejection_reason,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateGoodsReceiptItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.GoodsReceiptId != nil {
		hasUpdate = true
	}
	
	if r.PurchaseOrderItemId != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.QuantityReceived != nil {
		hasUpdate = true
	}
	
	if r.QuantityAccepted != nil {
		hasUpdate = true
	}
	
	if r.QuantityRejected != nil {
		hasUpdate = true
	}
	
	if r.RejectionReason != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// GoodsReceiptItemsListResponse represents a paginated list of goods_receipt_items records
type GoodsReceiptItemsListResponse struct {
	Items      []*GoodsReceiptItemsResponse `json:"items"`
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

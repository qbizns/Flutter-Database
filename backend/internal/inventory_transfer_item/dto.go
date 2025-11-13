package inventory_transfer_item

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InventoryTransferItemsResponse represents a inventory_transfer_items response
type InventoryTransferItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	InventoryTransferId uuid.UUID `json:"inventory_transfer_id"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	ProductName string `json:"product_name"`
	
	ProductSku *string `json:"product_sku"`
	
	QuantityRequested float64 `json:"quantity_requested"`
	
	QuantityShipped *float64 `json:"quantity_shipped"`
	
	QuantityReceived *float64 `json:"quantity_received"`
	
	UnitOfMeasure *string `json:"unit_of_measure"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	TotalCost *float64 `json:"total_cost"`
	
	ItemStatus *string `json:"item_status"`
	
	VarianceQuantity *float64 `json:"variance_quantity"`
	
	VarianceReason *string `json:"variance_reason"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	QuantityRequested *string `json:"quantity_requested"`
	
	QuantityShipped *string `json:"quantity_shipped"`
	
	QuantityReceived *string `json:"quantity_received"`
	
	QuantityShipped *string `json:"quantity_shipped"`
	
	QuantityReceived *string `json:"quantity_received"`
	
	(unitCost *string `json:"(unit_cost"`
	
	(totalCost *string `json:"(total_cost"`
	
}

// CreateInventoryTransferItemsRequest represents a request to create a inventory_transfer_items
type CreateInventoryTransferItemsRequest struct {
	
	InventoryTransferId uuid.UUID `json:"inventory_transfer_id" validate:"required"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	ProductName string `json:"product_name" validate:"required"`
	
	ProductSku *string `json:"product_sku"`
	
	QuantityRequested float64 `json:"quantity_requested" validate:"required"`
	
	QuantityShipped *float64 `json:"quantity_shipped"`
	
	QuantityReceived *float64 `json:"quantity_received"`
	
	UnitOfMeasure *string `json:"unit_of_measure"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	TotalCost *float64 `json:"total_cost"`
	
	// 	ItemStatus *string `json:"item_status"`
	
	VarianceQuantity *float64 `json:"variance_quantity"`
	
	VarianceReason *string `json:"variance_reason"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	QuantityRequested *string `json:"quantity_requested"`
	
	QuantityShipped *string `json:"quantity_shipped"`
	
	QuantityReceived *string `json:"quantity_received"`
	
	QuantityShipped *string `json:"quantity_shipped"`
	
	QuantityReceived *string `json:"quantity_received"`
	
	(unitCost *string `json:"(unit_cost"`
	
	(totalCost *string `json:"(total_cost"`
	
}

// Validate validates the create request
func (r *CreateInventoryTransferItemsRequest) Validate() error {
	
	if r.InventoryTransferId == uuid.Nil {
		return fmt.Errorf("inventory_transfer_id is required")
	}
	
	if r.ProductName == "" {
		return fmt.Errorf("product_name is required")
	}
	
	if r.QuantityRequested == nil {
		return fmt.Errorf("quantity_requested is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateInventoryTransferItemsRequest represents a request to update a inventory_transfer_items
type UpdateInventoryTransferItemsRequest struct {
	
	InventoryTransferId *uuid.UUID `json:"inventory_transfer_id,omitempty" validate:"omitempty,required"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	ProductName *string `json:"product_name,omitempty" validate:"omitempty,required"`
	
	ProductSku *string `json:"product_sku,omitempty"`
	
	QuantityRequested *float64 `json:"quantity_requested,omitempty" validate:"omitempty,required"`
	
	QuantityShipped *float64 `json:"quantity_shipped,omitempty"`
	
	QuantityReceived *float64 `json:"quantity_received,omitempty"`
	
	UnitOfMeasure *string `json:"unit_of_measure,omitempty"`
	
	UnitCost *float64 `json:"unit_cost,omitempty"`
	
	TotalCost *float64 `json:"total_cost,omitempty"`
	
	// 	ItemStatus *string `json:"item_status,omitempty"`
	
	VarianceQuantity *float64 `json:"variance_quantity,omitempty"`
	
	VarianceReason *string `json:"variance_reason,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	QuantityRequested *string `json:"quantity_requested,omitempty"`
	
	QuantityShipped *string `json:"quantity_shipped,omitempty"`
	
	QuantityReceived *string `json:"quantity_received,omitempty"`
	
	QuantityShipped *string `json:"quantity_shipped,omitempty"`
	
	QuantityReceived *string `json:"quantity_received,omitempty"`
	
	(unitCost *string `json:"(unit_cost,omitempty"`
	
	(totalCost *string `json:"(total_cost,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateInventoryTransferItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.InventoryTransferId != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.ProductName != nil {
		hasUpdate = true
	}
	
	if r.ProductSku != nil {
		hasUpdate = true
	}
	
	if r.QuantityRequested != nil {
		hasUpdate = true
	}
	
	if r.QuantityShipped != nil {
		hasUpdate = true
	}
	
	if r.QuantityReceived != nil {
		hasUpdate = true
	}
	
	if r.UnitOfMeasure != nil {
		hasUpdate = true
	}
	
	if r.UnitCost != nil {
		hasUpdate = true
	}
	
	if r.TotalCost != nil {
		hasUpdate = true
	}
	
	if r.ItemStatus != nil {
		hasUpdate = true
	}
	
	if r.VarianceQuantity != nil {
		hasUpdate = true
	}
	
	if r.VarianceReason != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.QuantityRequested != nil {
		hasUpdate = true
	}
	
	if r.QuantityShipped != nil {
		hasUpdate = true
	}
	
	if r.QuantityReceived != nil {
		hasUpdate = true
	}
	
	if r.QuantityShipped != nil {
		hasUpdate = true
	}
	
	if r.QuantityReceived != nil {
		hasUpdate = true
	}
	
	if r.(unitCost != nil {
		hasUpdate = true
	}
	
	if r.(totalCost != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// InventoryTransferItemsListResponse represents a paginated list of inventory_transfer_items records
type InventoryTransferItemsListResponse struct {
	Items      []*InventoryTransferItemsResponse `json:"items"`
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

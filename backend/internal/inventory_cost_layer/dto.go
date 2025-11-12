package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InventoryCostLayersResponse represents a inventory_cost_layers response
type InventoryCostLayersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	LotNumber *string `json:"lot_number"`
	
	SerialNumber *string `json:"serial_number"`
	
	LayerDate time.Time `json:"layer_date"`
	
	UnitCost float64 `json:"unit_cost"`
	
	OriginalQuantity float64 `json:"original_quantity"`
	
	RemainingQuantity float64 `json:"remaining_quantity"`
	
	UomCode *string `json:"uom_code"`
	
	SourceTransactionType *string `json:"source_transaction_type"`
	
	SourceTransactionId *uuid.UUID `json:"source_transaction_id"`
	
	SourceReference *string `json:"source_reference"`
	
	IsFullyConsumed *bool `json:"is_fully_consumed"`
	
	ConsumedAt *time.Time `json:"consumed_at"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateInventoryCostLayersRequest represents a request to create a inventory_cost_layers
type CreateInventoryCostLayersRequest struct {
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	LotNumber *string `json:"lot_number"`
	
	SerialNumber *string `json:"serial_number"`
	
	LayerDate time.Time `json:"layer_date" validate:"required"`
	
	UnitCost float64 `json:"unit_cost" validate:"required"`
	
	OriginalQuantity float64 `json:"original_quantity" validate:"required"`
	
	RemainingQuantity float64 `json:"remaining_quantity" validate:"required"`
	
	UomCode *string `json:"uom_code"`
	
	SourceTransactionType *string `json:"source_transaction_type"`
	
	SourceTransactionId *uuid.UUID `json:"source_transaction_id"`
	
	SourceReference *string `json:"source_reference"`
	
	IsFullyConsumed *bool `json:"is_fully_consumed"`
	
	ConsumedAt *time.Time `json:"consumed_at"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateInventoryCostLayersRequest) Validate() error {
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.LayerDate == nil {
		return fmt.Errorf("layer_date is required")
	}
	
	if r.UnitCost == nil {
		return fmt.Errorf("unit_cost is required")
	}
	
	if r.OriginalQuantity == nil {
		return fmt.Errorf("original_quantity is required")
	}
	
	if r.RemainingQuantity == nil {
		return fmt.Errorf("remaining_quantity is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateInventoryCostLayersRequest represents a request to update a inventory_cost_layers
type UpdateInventoryCostLayersRequest struct {
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	LotNumber *string `json:"lot_number,omitempty"`
	
	SerialNumber *string `json:"serial_number,omitempty"`
	
	LayerDate *time.Time `json:"layer_date,omitempty" validate:"omitempty,required"`
	
	UnitCost *float64 `json:"unit_cost,omitempty" validate:"omitempty,required"`
	
	OriginalQuantity *float64 `json:"original_quantity,omitempty" validate:"omitempty,required"`
	
	RemainingQuantity *float64 `json:"remaining_quantity,omitempty" validate:"omitempty,required"`
	
	UomCode *string `json:"uom_code,omitempty"`
	
	SourceTransactionType *string `json:"source_transaction_type,omitempty"`
	
	SourceTransactionId *uuid.UUID `json:"source_transaction_id,omitempty"`
	
	SourceReference *string `json:"source_reference,omitempty"`
	
	IsFullyConsumed *bool `json:"is_fully_consumed,omitempty"`
	
	ConsumedAt *time.Time `json:"consumed_at,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateInventoryCostLayersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.LotNumber != nil {
		hasUpdate = true
	}
	
	if r.SerialNumber != nil {
		hasUpdate = true
	}
	
	if r.LayerDate != nil {
		hasUpdate = true
	}
	
	if r.UnitCost != nil {
		hasUpdate = true
	}
	
	if r.OriginalQuantity != nil {
		hasUpdate = true
	}
	
	if r.RemainingQuantity != nil {
		hasUpdate = true
	}
	
	if r.UomCode != nil {
		hasUpdate = true
	}
	
	if r.SourceTransactionType != nil {
		hasUpdate = true
	}
	
	if r.SourceTransactionId != nil {
		hasUpdate = true
	}
	
	if r.SourceReference != nil {
		hasUpdate = true
	}
	
	if r.IsFullyConsumed != nil {
		hasUpdate = true
	}
	
	if r.ConsumedAt != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// InventoryCostLayersListResponse represents a paginated list of inventory_cost_layers records
type InventoryCostLayersListResponse struct {
	Items      []*InventoryCostLayersResponse `json:"items"`
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

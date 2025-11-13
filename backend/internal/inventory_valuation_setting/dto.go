package inventory_valuation_setting

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InventoryValuationSettingsResponse represents a inventory_valuation_settings response
type InventoryValuationSettingsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ValuationMethod string `json:"valuation_method"`
	
	'fifo', *string `json:"'fifo',"`
	
	'lifo', *string `json:"'lifo',"`
	
	'weightedAverage', *string `json:"'weighted_average',"`
	
	'movingAverage', *string `json:"'moving_average',"`
	
	'standardCost', *string `json:"'standard_cost',"`
	
	'specificId' *string `json:"'specific_id'"`
	
	CostLayerGranularity *string `json:"cost_layer_granularity"`
	
	'product', *string `json:"'product',"`
	
	'productLocation', *string `json:"'product_location',"`
	
	'productLocationLot', *string `json:"'product_location_lot',"`
	
	'serialNumber' *string `json:"'serial_number'"`
	
	DefaultInventoryAccountId *uuid.UUID `json:"default_inventory_account_id"`
	
	DefaultCogsAccountId *uuid.UUID `json:"default_cogs_account_id"`
	
	DefaultInventoryAdjustmentAccountId *uuid.UUID `json:"default_inventory_adjustment_account_id"`
	
	DefaultInventoryVarianceAccountId *uuid.UUID `json:"default_inventory_variance_account_id"`
	
	CogsRecognitionTiming *string `json:"cogs_recognition_timing"`
	
	'onSale', *string `json:"'on_sale',"`
	
	'onDelivery', *string `json:"'on_delivery',"`
	
	'onPayment' *string `json:"'on_payment'"`
	
	AllowNegativeInventory *bool `json:"allow_negative_inventory"`
	
	RevalueOnPurchase *bool `json:"revalue_on_purchase"`
	
	RoundUnitCostToDecimals *int64 `json:"round_unit_cost_to_decimals"`
	
	RevaluationFrequency *string `json:"revaluation_frequency"`
	
	'realTime', *string `json:"'real_time',"`
	
	'daily', *string `json:"'daily',"`
	
	'monthly', *string `json:"'monthly',"`
	
	'manual' *string `json:"'manual'"`
	
	IsActive *bool `json:"is_active"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateInventoryValuationSettingsRequest represents a request to create a inventory_valuation_settings
type CreateInventoryValuationSettingsRequest struct {
	
	ValuationMethod string `json:"valuation_method" validate:"required"`
	
	'fifo', *string `json:"'fifo',"`
	
	'lifo', *string `json:"'lifo',"`
	
	'weightedAverage', *string `json:"'weighted_average',"`
	
	'movingAverage', *string `json:"'moving_average',"`
	
	'standardCost', *string `json:"'standard_cost',"`
	
	'specificId' *string `json:"'specific_id'"`
	
	CostLayerGranularity *string `json:"cost_layer_granularity"`
	
	'product', *string `json:"'product',"`
	
	'productLocation', *string `json:"'product_location',"`
	
	'productLocationLot', *string `json:"'product_location_lot',"`
	
	'serialNumber' *string `json:"'serial_number'"`
	
	DefaultInventoryAccountId *uuid.UUID `json:"default_inventory_account_id"`
	
	DefaultCogsAccountId *uuid.UUID `json:"default_cogs_account_id"`
	
	DefaultInventoryAdjustmentAccountId *uuid.UUID `json:"default_inventory_adjustment_account_id"`
	
	DefaultInventoryVarianceAccountId *uuid.UUID `json:"default_inventory_variance_account_id"`
	
	CogsRecognitionTiming *string `json:"cogs_recognition_timing"`
	
	'onSale', *string `json:"'on_sale',"`
	
	'onDelivery', *string `json:"'on_delivery',"`
	
	'onPayment' *string `json:"'on_payment'"`
	
	AllowNegativeInventory *bool `json:"allow_negative_inventory"`
	
	RevalueOnPurchase *bool `json:"revalue_on_purchase"`
	
	RoundUnitCostToDecimals *int64 `json:"round_unit_cost_to_decimals"`
	
	RevaluationFrequency *string `json:"revaluation_frequency"`
	
	'realTime', *string `json:"'real_time',"`
	
	'daily', *string `json:"'daily',"`
	
	'monthly', *string `json:"'monthly',"`
	
	'manual' *string `json:"'manual'"`
	
	IsActive *bool `json:"is_active"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateInventoryValuationSettingsRequest) Validate() error {
	
	if r.ValuationMethod == "" {
		return fmt.Errorf("valuation_method is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateInventoryValuationSettingsRequest represents a request to update a inventory_valuation_settings
type UpdateInventoryValuationSettingsRequest struct {
	
	ValuationMethod *string `json:"valuation_method,omitempty" validate:"omitempty,required"`
	
	'fifo', *string `json:"'fifo',,omitempty"`
	
	'lifo', *string `json:"'lifo',,omitempty"`
	
	'weightedAverage', *string `json:"'weighted_average',,omitempty"`
	
	'movingAverage', *string `json:"'moving_average',,omitempty"`
	
	'standardCost', *string `json:"'standard_cost',,omitempty"`
	
	'specificId' *string `json:"'specific_id',omitempty"`
	
	CostLayerGranularity *string `json:"cost_layer_granularity,omitempty"`
	
	'product', *string `json:"'product',,omitempty"`
	
	'productLocation', *string `json:"'product_location',,omitempty"`
	
	'productLocationLot', *string `json:"'product_location_lot',,omitempty"`
	
	'serialNumber' *string `json:"'serial_number',omitempty"`
	
	DefaultInventoryAccountId *uuid.UUID `json:"default_inventory_account_id,omitempty"`
	
	DefaultCogsAccountId *uuid.UUID `json:"default_cogs_account_id,omitempty"`
	
	DefaultInventoryAdjustmentAccountId *uuid.UUID `json:"default_inventory_adjustment_account_id,omitempty"`
	
	DefaultInventoryVarianceAccountId *uuid.UUID `json:"default_inventory_variance_account_id,omitempty"`
	
	CogsRecognitionTiming *string `json:"cogs_recognition_timing,omitempty"`
	
	'onSale', *string `json:"'on_sale',,omitempty"`
	
	'onDelivery', *string `json:"'on_delivery',,omitempty"`
	
	'onPayment' *string `json:"'on_payment',omitempty"`
	
	AllowNegativeInventory *bool `json:"allow_negative_inventory,omitempty"`
	
	RevalueOnPurchase *bool `json:"revalue_on_purchase,omitempty"`
	
	RoundUnitCostToDecimals *int64 `json:"round_unit_cost_to_decimals,omitempty"`
	
	RevaluationFrequency *string `json:"revaluation_frequency,omitempty"`
	
	'realTime', *string `json:"'real_time',,omitempty"`
	
	'daily', *string `json:"'daily',,omitempty"`
	
	'monthly', *string `json:"'monthly',,omitempty"`
	
	'manual' *string `json:"'manual',omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateInventoryValuationSettingsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ValuationMethod != nil {
		hasUpdate = true
	}
	
	if r.'fifo', != nil {
		hasUpdate = true
	}
	
	if r.'lifo', != nil {
		hasUpdate = true
	}
	
	if r.'weightedAverage', != nil {
		hasUpdate = true
	}
	
	if r.'movingAverage', != nil {
		hasUpdate = true
	}
	
	if r.'standardCost', != nil {
		hasUpdate = true
	}
	
	if r.'specificId' != nil {
		hasUpdate = true
	}
	
	if r.CostLayerGranularity != nil {
		hasUpdate = true
	}
	
	if r.'product', != nil {
		hasUpdate = true
	}
	
	if r.'productLocation', != nil {
		hasUpdate = true
	}
	
	if r.'productLocationLot', != nil {
		hasUpdate = true
	}
	
	if r.'serialNumber' != nil {
		hasUpdate = true
	}
	
	if r.DefaultInventoryAccountId != nil {
		hasUpdate = true
	}
	
	if r.DefaultCogsAccountId != nil {
		hasUpdate = true
	}
	
	if r.DefaultInventoryAdjustmentAccountId != nil {
		hasUpdate = true
	}
	
	if r.DefaultInventoryVarianceAccountId != nil {
		hasUpdate = true
	}
	
	if r.CogsRecognitionTiming != nil {
		hasUpdate = true
	}
	
	if r.'onSale', != nil {
		hasUpdate = true
	}
	
	if r.'onDelivery', != nil {
		hasUpdate = true
	}
	
	if r.'onPayment' != nil {
		hasUpdate = true
	}
	
	if r.AllowNegativeInventory != nil {
		hasUpdate = true
	}
	
	if r.RevalueOnPurchase != nil {
		hasUpdate = true
	}
	
	if r.RoundUnitCostToDecimals != nil {
		hasUpdate = true
	}
	
	if r.RevaluationFrequency != nil {
		hasUpdate = true
	}
	
	if r.'realTime', != nil {
		hasUpdate = true
	}
	
	if r.'daily', != nil {
		hasUpdate = true
	}
	
	if r.'monthly', != nil {
		hasUpdate = true
	}
	
	if r.'manual' != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.EffectiveFrom != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
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

// InventoryValuationSettingsListResponse represents a paginated list of inventory_valuation_settings records
type InventoryValuationSettingsListResponse struct {
	Items      []*InventoryValuationSettingsResponse `json:"items"`
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

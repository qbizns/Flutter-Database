package cycle_count_item

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CycleCountItemsResponse represents a cycle_count_items response
type CycleCountItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CycleCountId uuid.UUID `json:"cycle_count_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	ProductName string `json:"product_name"`
	
	ProductSku *string `json:"product_sku"`
	
	SystemQuantity float64 `json:"system_quantity"`
	
	CountedQuantity *float64 `json:"counted_quantity"`
	
	VarianceQuantity *float64 `json:"variance_quantity"`
	
	VariancePercentage *float64 `json:"variance_percentage"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	VarianceValue *float64 `json:"variance_value"`
	
	Status *string `json:"status"`
	
	RecountRequired *bool `json:"recount_required"`
	
	RecountQuantity *float64 `json:"recount_quantity"`
	
	RecountReason *string `json:"recount_reason"`
	
	AdjustmentApplied *bool `json:"adjustment_applied"`
	
	AdjustmentDate *time.Time `json:"adjustment_date"`
	
	AdjustmentReason *string `json:"adjustment_reason"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CountedAt *time.Time `json:"counted_at"`
	
	CountedBy *uuid.UUID `json:"counted_by"`
	
}

// CreateCycleCountItemsRequest represents a request to create a cycle_count_items
type CreateCycleCountItemsRequest struct {
	
	CycleCountId uuid.UUID `json:"cycle_count_id" validate:"required"`
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	ProductName string `json:"product_name" validate:"required"`
	
	ProductSku *string `json:"product_sku"`
	
	SystemQuantity float64 `json:"system_quantity" validate:"required"`
	
	CountedQuantity *float64 `json:"counted_quantity"`
	
	VarianceQuantity *float64 `json:"variance_quantity"`
	
	VariancePercentage *float64 `json:"variance_percentage"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	VarianceValue *float64 `json:"variance_value"`
	
	Status *string `json:"status"`
	
	RecountRequired *bool `json:"recount_required"`
	
	RecountQuantity *float64 `json:"recount_quantity"`
	
	RecountReason *string `json:"recount_reason"`
	
	AdjustmentApplied *bool `json:"adjustment_applied"`
	
	AdjustmentDate *time.Time `json:"adjustment_date"`
	
	AdjustmentReason *string `json:"adjustment_reason"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CountedAt *time.Time `json:"counted_at"`
	
	CountedBy *uuid.UUID `json:"counted_by"`
	
}

// Validate validates the create request
func (r *CreateCycleCountItemsRequest) Validate() error {
	
	if r.CycleCountId == uuid.Nil {
		return fmt.Errorf("cycle_count_id is required")
	}
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.ProductName == "" {
		return fmt.Errorf("product_name is required")
	}
	
	if r.SystemQuantity == nil {
		return fmt.Errorf("system_quantity is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCycleCountItemsRequest represents a request to update a cycle_count_items
type UpdateCycleCountItemsRequest struct {
	
	CycleCountId *uuid.UUID `json:"cycle_count_id,omitempty" validate:"omitempty,required"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	ProductName *string `json:"product_name,omitempty" validate:"omitempty,required"`
	
	ProductSku *string `json:"product_sku,omitempty"`
	
	SystemQuantity *float64 `json:"system_quantity,omitempty" validate:"omitempty,required"`
	
	CountedQuantity *float64 `json:"counted_quantity,omitempty"`
	
	VarianceQuantity *float64 `json:"variance_quantity,omitempty"`
	
	VariancePercentage *float64 `json:"variance_percentage,omitempty"`
	
	UnitCost *float64 `json:"unit_cost,omitempty"`
	
	VarianceValue *float64 `json:"variance_value,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	RecountRequired *bool `json:"recount_required,omitempty"`
	
	RecountQuantity *float64 `json:"recount_quantity,omitempty"`
	
	RecountReason *string `json:"recount_reason,omitempty"`
	
	AdjustmentApplied *bool `json:"adjustment_applied,omitempty"`
	
	AdjustmentDate *time.Time `json:"adjustment_date,omitempty"`
	
	AdjustmentReason *string `json:"adjustment_reason,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CountedAt *time.Time `json:"counted_at,omitempty"`
	
	CountedBy *uuid.UUID `json:"counted_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCycleCountItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CycleCountId != nil {
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
	
	if r.SystemQuantity != nil {
		hasUpdate = true
	}
	
	if r.CountedQuantity != nil {
		hasUpdate = true
	}
	
	if r.VarianceQuantity != nil {
		hasUpdate = true
	}
	
	if r.VariancePercentage != nil {
		hasUpdate = true
	}
	
	if r.UnitCost != nil {
		hasUpdate = true
	}
	
	if r.VarianceValue != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.RecountRequired != nil {
		hasUpdate = true
	}
	
	if r.RecountQuantity != nil {
		hasUpdate = true
	}
	
	if r.RecountReason != nil {
		hasUpdate = true
	}
	
	if r.AdjustmentApplied != nil {
		hasUpdate = true
	}
	
	if r.AdjustmentDate != nil {
		hasUpdate = true
	}
	
	if r.AdjustmentReason != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CountedAt != nil {
		hasUpdate = true
	}
	
	if r.CountedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CycleCountItemsListResponse represents a paginated list of cycle_count_items records
type CycleCountItemsListResponse struct {
	Items      []*CycleCountItemsResponse `json:"items"`
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

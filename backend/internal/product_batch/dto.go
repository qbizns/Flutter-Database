package product_batch

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductBatchesResponse represents a product_batches response
type ProductBatchesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	BatchNumber string `json:"batch_number"`
	
	LotNumber *string `json:"lot_number"`
	
	Status *string `json:"status"`
	
	InitialQuantity float64 `json:"initial_quantity"`
	
	CurrentQuantity float64 `json:"current_quantity"`
	
	UnitOfMeasure *string `json:"unit_of_measure"`
	
	ManufacturingDate *time.Time `json:"manufacturing_date"`
	
	ExpirationDate *time.Time `json:"expiration_date"`
	
	ReceivedDate time.Time `json:"received_date"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	SupplierId *uuid.UUID `json:"supplier_id"`
	
	SupplierBatchNumber *string `json:"supplier_batch_number"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	TotalCost *float64 `json:"total_cost"`
	
	// 	QualityStatus *string `json:"quality_status"`
	
	QualityCheckDate *time.Time `json:"quality_check_date"`
	
	QualityCheckedBy *uuid.UUID `json:"quality_checked_by"`
	
	QualityNotes *string `json:"quality_notes"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	InitialQuantity *string `json:"initial_quantity"`
	
	CurrentQuantity *string `json:"current_quantity"`
	
	CurrentQuantity *string `json:"current_quantity"`
	
	ExpirationDate *string `json:"expiration_date"`
	
}

// CreateProductBatchesRequest represents a request to create a product_batches
type CreateProductBatchesRequest struct {
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	BatchNumber string `json:"batch_number" validate:"required"`
	
	LotNumber *string `json:"lot_number"`
	
	// 	Status *string `json:"status"`
	
	InitialQuantity float64 `json:"initial_quantity" validate:"required"`
	
	CurrentQuantity float64 `json:"current_quantity" validate:"required"`
	
	UnitOfMeasure *string `json:"unit_of_measure"`
	
	ManufacturingDate *time.Time `json:"manufacturing_date"`
	
	ExpirationDate *time.Time `json:"expiration_date"`
	
	ReceivedDate time.Time `json:"received_date" validate:"required"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	SupplierId *uuid.UUID `json:"supplier_id"`
	
	SupplierBatchNumber *string `json:"supplier_batch_number"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	TotalCost *float64 `json:"total_cost"`
	
	// 	QualityStatus *string `json:"quality_status"`
	
	QualityCheckDate *time.Time `json:"quality_check_date"`
	
	QualityCheckedBy *uuid.UUID `json:"quality_checked_by"`
	
	QualityNotes *string `json:"quality_notes"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	InitialQuantity *string `json:"initial_quantity"`
	
	CurrentQuantity *string `json:"current_quantity"`
	
	CurrentQuantity *string `json:"current_quantity"`
	
	ExpirationDate *string `json:"expiration_date"`
	
}

// Validate validates the create request
func (r *CreateProductBatchesRequest) Validate() error {
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.BatchNumber == "" {
		return fmt.Errorf("batch_number is required")
	}
	
	if r.InitialQuantity == nil {
		return fmt.Errorf("initial_quantity is required")
	}
	
	if r.CurrentQuantity == nil {
		return fmt.Errorf("current_quantity is required")
	}
	
	if r.ReceivedDate.IsZero() {
		return fmt.Errorf("received_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateProductBatchesRequest represents a request to update a product_batches
type UpdateProductBatchesRequest struct {
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	BatchNumber *string `json:"batch_number,omitempty" validate:"omitempty,required"`
	
	LotNumber *string `json:"lot_number,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	InitialQuantity *float64 `json:"initial_quantity,omitempty" validate:"omitempty,required"`
	
	CurrentQuantity *float64 `json:"current_quantity,omitempty" validate:"omitempty,required"`
	
	UnitOfMeasure *string `json:"unit_of_measure,omitempty"`
	
	ManufacturingDate *time.Time `json:"manufacturing_date,omitempty"`
	
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	
	ReceivedDate *time.Time `json:"received_date,omitempty" validate:"omitempty,required"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id,omitempty"`
	
	SupplierId *uuid.UUID `json:"supplier_id,omitempty"`
	
	SupplierBatchNumber *string `json:"supplier_batch_number,omitempty"`
	
	UnitCost *float64 `json:"unit_cost,omitempty"`
	
	TotalCost *float64 `json:"total_cost,omitempty"`
	
	// 	QualityStatus *string `json:"quality_status,omitempty"`
	
	QualityCheckDate *time.Time `json:"quality_check_date,omitempty"`
	
	QualityCheckedBy *uuid.UUID `json:"quality_checked_by,omitempty"`
	
	QualityNotes *string `json:"quality_notes,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	InitialQuantity *string `json:"initial_quantity,omitempty"`
	
	CurrentQuantity *string `json:"current_quantity,omitempty"`
	
	CurrentQuantity *string `json:"current_quantity,omitempty"`
	
	ExpirationDate *string `json:"expiration_date,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateProductBatchesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.BatchNumber != nil {
		hasUpdate = true
	}
	
	if r.LotNumber != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.InitialQuantity != nil {
		hasUpdate = true
	}
	
	if r.CurrentQuantity != nil {
		hasUpdate = true
	}
	
	if r.UnitOfMeasure != nil {
		hasUpdate = true
	}
	
	if r.ManufacturingDate != nil {
		hasUpdate = true
	}
	
	if r.ExpirationDate != nil {
		hasUpdate = true
	}
	
	if r.ReceivedDate != nil {
		hasUpdate = true
	}
	
	if r.PurchaseOrderId != nil {
		hasUpdate = true
	}
	
	if r.SupplierId != nil {
		hasUpdate = true
	}
	
	if r.SupplierBatchNumber != nil {
		hasUpdate = true
	}
	
	if r.UnitCost != nil {
		hasUpdate = true
	}
	
	if r.TotalCost != nil {
		hasUpdate = true
	}
	
	if r.QualityStatus != nil {
		hasUpdate = true
	}
	
	if r.QualityCheckDate != nil {
		hasUpdate = true
	}
	
	if r.QualityCheckedBy != nil {
		hasUpdate = true
	}
	
	if r.QualityNotes != nil {
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
	
	if r.InitialQuantity != nil {
		hasUpdate = true
	}
	
	if r.CurrentQuantity != nil {
		hasUpdate = true
	}
	
	if r.CurrentQuantity != nil {
		hasUpdate = true
	}
	
	if r.ExpirationDate != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ProductBatchesListResponse represents a paginated list of product_batches records
type ProductBatchesListResponse struct {
	Items      []*ProductBatchesResponse `json:"items"`
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

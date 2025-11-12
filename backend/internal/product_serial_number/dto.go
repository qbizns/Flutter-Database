package product_serial_number

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductSerialNumbersResponse represents a product_serial_numbers response
type ProductSerialNumbersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	SerialNumber string `json:"serial_number"`
	
	Status *string `json:"status"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	PurchaseDate *time.Time `json:"purchase_date"`
	
	PurchaseCost *float64 `json:"purchase_cost"`
	
	SupplierId *uuid.UUID `json:"supplier_id"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	SaleDate *time.Time `json:"sale_date"`
	
	SalePrice *float64 `json:"sale_price"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	WarrantyStartDate *time.Time `json:"warranty_start_date"`
	
	WarrantyEndDate *time.Time `json:"warranty_end_date"`
	
	WarrantyProvider *string `json:"warranty_provider"`
	
	WarrantyTerms *string `json:"warranty_terms"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateProductSerialNumbersRequest represents a request to create a product_serial_numbers
type CreateProductSerialNumbersRequest struct {
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	SerialNumber string `json:"serial_number" validate:"required"`
	
	Status *string `json:"status"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	PurchaseDate *time.Time `json:"purchase_date"`
	
	PurchaseCost *float64 `json:"purchase_cost"`
	
	SupplierId *uuid.UUID `json:"supplier_id"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	SaleDate *time.Time `json:"sale_date"`
	
	SalePrice *float64 `json:"sale_price"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	WarrantyStartDate *time.Time `json:"warranty_start_date"`
	
	WarrantyEndDate *time.Time `json:"warranty_end_date"`
	
	WarrantyProvider *string `json:"warranty_provider"`
	
	WarrantyTerms *string `json:"warranty_terms"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateProductSerialNumbersRequest) Validate() error {
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.SerialNumber == "" {
		return fmt.Errorf("serial_number is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateProductSerialNumbersRequest represents a request to update a product_serial_numbers
type UpdateProductSerialNumbersRequest struct {
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	SerialNumber *string `json:"serial_number,omitempty" validate:"omitempty,required"`
	
	Status *string `json:"status,omitempty"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id,omitempty"`
	
	PurchaseDate *time.Time `json:"purchase_date,omitempty"`
	
	PurchaseCost *float64 `json:"purchase_cost,omitempty"`
	
	SupplierId *uuid.UUID `json:"supplier_id,omitempty"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	SaleDate *time.Time `json:"sale_date,omitempty"`
	
	SalePrice *float64 `json:"sale_price,omitempty"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty"`
	
	WarrantyStartDate *time.Time `json:"warranty_start_date,omitempty"`
	
	WarrantyEndDate *time.Time `json:"warranty_end_date,omitempty"`
	
	WarrantyProvider *string `json:"warranty_provider,omitempty"`
	
	WarrantyTerms *string `json:"warranty_terms,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateProductSerialNumbersRequest) Validate() error {
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
	
	if r.SerialNumber != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.PurchaseOrderId != nil {
		hasUpdate = true
	}
	
	if r.PurchaseDate != nil {
		hasUpdate = true
	}
	
	if r.PurchaseCost != nil {
		hasUpdate = true
	}
	
	if r.SupplierId != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.SaleDate != nil {
		hasUpdate = true
	}
	
	if r.SalePrice != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.WarrantyStartDate != nil {
		hasUpdate = true
	}
	
	if r.WarrantyEndDate != nil {
		hasUpdate = true
	}
	
	if r.WarrantyProvider != nil {
		hasUpdate = true
	}
	
	if r.WarrantyTerms != nil {
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

// ProductSerialNumbersListResponse represents a paginated list of product_serial_numbers records
type ProductSerialNumbersListResponse struct {
	Items      []*ProductSerialNumbersResponse `json:"items"`
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

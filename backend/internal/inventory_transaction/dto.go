package inventory_transaction

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InventoryTransactionsResponse represents a inventory_transactions response
type InventoryTransactionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	TransactionType string `json:"transaction_type"`
	
	Quantity float64 `json:"quantity"`
	
	Unit *string `json:"unit"`
	
	BalanceAfter float64 `json:"balance_after"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	TotalCost *float64 `json:"total_cost"`
	
	TransactionDate time.Time `json:"transaction_date"`
	
	Notes *string `json:"notes"`
	
	Reason *string `json:"reason"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// CreateInventoryTransactionsRequest represents a request to create a inventory_transactions
type CreateInventoryTransactionsRequest struct {
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	TransactionType string `json:"transaction_type" validate:"required"`
	
	Quantity float64 `json:"quantity" validate:"required"`
	
	Unit *string `json:"unit"`
	
	BalanceAfter float64 `json:"balance_after" validate:"required"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	UnitCost *float64 `json:"unit_cost"`
	
	TotalCost *float64 `json:"total_cost"`
	
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
	
	Notes *string `json:"notes"`
	
	Reason *string `json:"reason"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateInventoryTransactionsRequest) Validate() error {
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.TransactionType == "" {
		return fmt.Errorf("transaction_type is required")
	}
	
	if r.Quantity == nil {
		return fmt.Errorf("quantity is required")
	}
	
	if r.BalanceAfter == nil {
		return fmt.Errorf("balance_after is required")
	}
	
	if r.TransactionDate.IsZero() {
		return fmt.Errorf("transaction_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateInventoryTransactionsRequest represents a request to update a inventory_transactions
type UpdateInventoryTransactionsRequest struct {
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,required"`
	
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,required"`
	
	Unit *string `json:"unit,omitempty"`
	
	BalanceAfter *float64 `json:"balance_after,omitempty" validate:"omitempty,required"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	UnitCost *float64 `json:"unit_cost,omitempty"`
	
	TotalCost *float64 `json:"total_cost,omitempty"`
	
	TransactionDate *time.Time `json:"transaction_date,omitempty" validate:"omitempty,required"`
	
	Notes *string `json:"notes,omitempty"`
	
	Reason *string `json:"reason,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateInventoryTransactionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.TransactionType != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.Unit != nil {
		hasUpdate = true
	}
	
	if r.BalanceAfter != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.ReferenceNumber != nil {
		hasUpdate = true
	}
	
	if r.UnitCost != nil {
		hasUpdate = true
	}
	
	if r.TotalCost != nil {
		hasUpdate = true
	}
	
	if r.TransactionDate != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Reason != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
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

// InventoryTransactionsListResponse represents a paginated list of inventory_transactions records
type InventoryTransactionsListResponse struct {
	Items      []*InventoryTransactionsResponse `json:"items"`
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

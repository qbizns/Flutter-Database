package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BatchTransactionsResponse represents a batch_transactions response
type BatchTransactionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	BatchId uuid.UUID `json:"batch_id"`
	
	TransactionType string `json:"transaction_type"`
	
	Quantity float64 `json:"quantity"`
	
	BalanceAfter float64 `json:"balance_after"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	InventoryTransferId *uuid.UUID `json:"inventory_transfer_id"`
	
	Reason *string `json:"reason"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	TransactionDate *time.Time `json:"transaction_date"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// CreateBatchTransactionsRequest represents a request to create a batch_transactions
type CreateBatchTransactionsRequest struct {
	
	BatchId uuid.UUID `json:"batch_id" validate:"required"`
	
	TransactionType string `json:"transaction_type" validate:"required"`
	
	Quantity float64 `json:"quantity" validate:"required"`
	
	BalanceAfter float64 `json:"balance_after" validate:"required"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	InventoryTransferId *uuid.UUID `json:"inventory_transfer_id"`
	
	Reason *string `json:"reason"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	TransactionDate *time.Time `json:"transaction_date"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateBatchTransactionsRequest) Validate() error {
	
	if r.BatchId == uuid.Nil {
		return fmt.Errorf("batch_id is required")
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
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBatchTransactionsRequest represents a request to update a batch_transactions
type UpdateBatchTransactionsRequest struct {
	
	BatchId *uuid.UUID `json:"batch_id,omitempty" validate:"omitempty,required"`
	
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,required"`
	
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,required"`
	
	BalanceAfter *float64 `json:"balance_after,omitempty" validate:"omitempty,required"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	InventoryTransferId *uuid.UUID `json:"inventory_transfer_id,omitempty"`
	
	Reason *string `json:"reason,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	TransactionDate *time.Time `json:"transaction_date,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBatchTransactionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BatchId != nil {
		hasUpdate = true
	}
	
	if r.TransactionType != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.BalanceAfter != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.InventoryTransferId != nil {
		hasUpdate = true
	}
	
	if r.Reason != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.TransactionDate != nil {
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

// BatchTransactionsListResponse represents a paginated list of batch_transactions records
type BatchTransactionsListResponse struct {
	Items      []*BatchTransactionsResponse `json:"items"`
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

package deferred_expense_contract

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeferredExpenseContractsResponse represents a deferred_expense_contracts response
type DeferredExpenseContractsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	VendorBillId *uuid.UUID `json:"vendor_bill_id"`
	
	BillLineId *uuid.UUID `json:"bill_line_id"`
	
	ContractName *string `json:"contract_name"`
	
	TotalDeferredAmount float64 `json:"total_deferred_amount"`
	
	StartDate time.Time `json:"start_date"`
	
	EndDate time.Time `json:"end_date"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	DeferredAccountId uuid.UUID `json:"deferred_account_id"`
	
	ExpenseAccountId uuid.UUID `json:"expense_account_id"`
	
	Status *string `json:"status"`
	
	RecognizedAmount *float64 `json:"recognized_amount"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateDeferredExpenseContractsRequest represents a request to create a deferred_expense_contracts
type CreateDeferredExpenseContractsRequest struct {
	
	VendorBillId *uuid.UUID `json:"vendor_bill_id"`
	
	BillLineId *uuid.UUID `json:"bill_line_id"`
	
	ContractName *string `json:"contract_name"`
	
	TotalDeferredAmount float64 `json:"total_deferred_amount" validate:"required"`
	
	StartDate time.Time `json:"start_date" validate:"required"`
	
	EndDate time.Time `json:"end_date" validate:"required"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	DeferredAccountId uuid.UUID `json:"deferred_account_id" validate:"required"`
	
	ExpenseAccountId uuid.UUID `json:"expense_account_id" validate:"required"`
	
	Status *string `json:"status"`
	
	RecognizedAmount *float64 `json:"recognized_amount"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateDeferredExpenseContractsRequest) Validate() error {
	
	if r.TotalDeferredAmount == nil {
		return fmt.Errorf("total_deferred_amount is required")
	}
	
	if r.StartDate == nil {
		return fmt.Errorf("start_date is required")
	}
	
	if r.EndDate == nil {
		return fmt.Errorf("end_date is required")
	}
	
	if r.DeferredAccountId == uuid.Nil {
		return fmt.Errorf("deferred_account_id is required")
	}
	
	if r.ExpenseAccountId == uuid.Nil {
		return fmt.Errorf("expense_account_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDeferredExpenseContractsRequest represents a request to update a deferred_expense_contracts
type UpdateDeferredExpenseContractsRequest struct {
	
	VendorBillId *uuid.UUID `json:"vendor_bill_id,omitempty"`
	
	BillLineId *uuid.UUID `json:"bill_line_id,omitempty"`
	
	ContractName *string `json:"contract_name,omitempty"`
	
	TotalDeferredAmount *float64 `json:"total_deferred_amount,omitempty" validate:"omitempty,required"`
	
	StartDate *time.Time `json:"start_date,omitempty" validate:"omitempty,required"`
	
	EndDate *time.Time `json:"end_date,omitempty" validate:"omitempty,required"`
	
	RecognitionMethod *string `json:"recognition_method,omitempty"`
	
	RecognitionMethod *string `json:"recognition_method,omitempty"`
	
	DeferredAccountId *uuid.UUID `json:"deferred_account_id,omitempty" validate:"omitempty,required"`
	
	ExpenseAccountId *uuid.UUID `json:"expense_account_id,omitempty" validate:"omitempty,required"`
	
	Status *string `json:"status,omitempty"`
	
	RecognizedAmount *float64 `json:"recognized_amount,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDeferredExpenseContractsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.VendorBillId != nil {
		hasUpdate = true
	}
	
	if r.BillLineId != nil {
		hasUpdate = true
	}
	
	if r.ContractName != nil {
		hasUpdate = true
	}
	
	if r.TotalDeferredAmount != nil {
		hasUpdate = true
	}
	
	if r.StartDate != nil {
		hasUpdate = true
	}
	
	if r.EndDate != nil {
		hasUpdate = true
	}
	
	if r.RecognitionMethod != nil {
		hasUpdate = true
	}
	
	if r.RecognitionMethod != nil {
		hasUpdate = true
	}
	
	if r.DeferredAccountId != nil {
		hasUpdate = true
	}
	
	if r.ExpenseAccountId != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.RecognizedAmount != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
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

// DeferredExpenseContractsListResponse represents a paginated list of deferred_expense_contracts records
type DeferredExpenseContractsListResponse struct {
	Items      []*DeferredExpenseContractsResponse `json:"items"`
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

package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomerStoreCreditAccountsResponse represents a customer_store_credit_accounts response
type CustomerStoreCreditAccountsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerId uuid.UUID `json:"customer_id"`
	
	CurrentBalance *float64 `json:"current_balance"`
	
	CreditLimit *float64 `json:"credit_limit"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCustomerStoreCreditAccountsRequest represents a request to create a customer_store_credit_accounts
type CreateCustomerStoreCreditAccountsRequest struct {
	
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	
	CurrentBalance *float64 `json:"current_balance"`
	
	CreditLimit *float64 `json:"credit_limit"`
	
	IsActive *bool `json:"is_active"`
	
}

// Validate validates the create request
func (r *CreateCustomerStoreCreditAccountsRequest) Validate() error {
	
	if r.CustomerId == uuid.Nil {
		return fmt.Errorf("customer_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCustomerStoreCreditAccountsRequest represents a request to update a customer_store_credit_accounts
type UpdateCustomerStoreCreditAccountsRequest struct {
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty" validate:"omitempty,required"`
	
	CurrentBalance *float64 `json:"current_balance,omitempty"`
	
	CreditLimit *float64 `json:"credit_limit,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCustomerStoreCreditAccountsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.CurrentBalance != nil {
		hasUpdate = true
	}
	
	if r.CreditLimit != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CustomerStoreCreditAccountsListResponse represents a paginated list of customer_store_credit_accounts records
type CustomerStoreCreditAccountsListResponse struct {
	Items      []*CustomerStoreCreditAccountsResponse `json:"items"`
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

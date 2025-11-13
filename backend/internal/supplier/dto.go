package supplier

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SuppliersResponse represents a suppliers response
type SuppliersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SupplierCode string `json:"supplier_code"`
	
	Name string `json:"name"`
	
	ContactPerson *string `json:"contact_person"`
	
	Email *string `json:"email"`
	
	Phone *string `json:"phone"`
	
	Address *string `json:"address"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	TaxNumber *string `json:"tax_number"`
	
	PaymentTerms *string `json:"payment_terms"`
	
	CreditLimit *float64 `json:"credit_limit"`
	
	OutstandingBalance *float64 `json:"outstanding_balance"`
	
	TotalPurchases *float64 `json:"total_purchases"`
	
	TotalOrders *int64 `json:"total_orders"`
	
	LastOrderDate *time.Time `json:"last_order_date"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateSuppliersRequest represents a request to create a suppliers
type CreateSuppliersRequest struct {
	
	SupplierCode string `json:"supplier_code" validate:"required"`
	
	Name string `json:"name" validate:"required"`
	
	ContactPerson *string `json:"contact_person"`
	
	Email *string `json:"email" validate:"email"`
	
	Phone *string `json:"phone" validate:"e164"`
	
	Address *string `json:"address"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	TaxNumber *string `json:"tax_number"`
	
	PaymentTerms *string `json:"payment_terms"`
	
	CreditLimit *float64 `json:"credit_limit"`
	
	OutstandingBalance *float64 `json:"outstanding_balance"`
	
	TotalPurchases *float64 `json:"total_purchases"`
	
	TotalOrders *int64 `json:"total_orders"`
	
	LastOrderDate *time.Time `json:"last_order_date"`
	
	// 	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateSuppliersRequest) Validate() error {
	
	if r.SupplierCode == "" {
		return fmt.Errorf("supplier_code is required")
	}
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSuppliersRequest represents a request to update a suppliers
type UpdateSuppliersRequest struct {
	
	SupplierCode *string `json:"supplier_code,omitempty" validate:"omitempty,required"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	ContactPerson *string `json:"contact_person,omitempty"`
	
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	
	Phone *string `json:"phone,omitempty" validate:"omitempty,e164"`
	
	Address *string `json:"address,omitempty"`
	
	City *string `json:"city,omitempty"`
	
	State *string `json:"state,omitempty"`
	
	Country *string `json:"country,omitempty"`
	
	PostalCode *string `json:"postal_code,omitempty"`
	
	TaxNumber *string `json:"tax_number,omitempty"`
	
	PaymentTerms *string `json:"payment_terms,omitempty"`
	
	CreditLimit *float64 `json:"credit_limit,omitempty"`
	
	OutstandingBalance *float64 `json:"outstanding_balance,omitempty"`
	
	TotalPurchases *float64 `json:"total_purchases,omitempty"`
	
	TotalOrders *int64 `json:"total_orders,omitempty"`
	
	LastOrderDate *time.Time `json:"last_order_date,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSuppliersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SupplierCode != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.ContactPerson != nil {
		hasUpdate = true
	}
	
	if r.Email != nil {
		hasUpdate = true
	}
	
	if r.Phone != nil {
		hasUpdate = true
	}
	
	if r.Address != nil {
		hasUpdate = true
	}
	
	if r.City != nil {
		hasUpdate = true
	}
	
	if r.State != nil {
		hasUpdate = true
	}
	
	if r.Country != nil {
		hasUpdate = true
	}
	
	if r.PostalCode != nil {
		hasUpdate = true
	}
	
	if r.TaxNumber != nil {
		hasUpdate = true
	}
	
	if r.PaymentTerms != nil {
		hasUpdate = true
	}
	
	if r.CreditLimit != nil {
		hasUpdate = true
	}
	
	if r.OutstandingBalance != nil {
		hasUpdate = true
	}
	
	if r.TotalPurchases != nil {
		hasUpdate = true
	}
	
	if r.TotalOrders != nil {
		hasUpdate = true
	}
	
	if r.LastOrderDate != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
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

// SuppliersListResponse represents a paginated list of suppliers records
type SuppliersListResponse struct {
	Items      []*SuppliersResponse `json:"items"`
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

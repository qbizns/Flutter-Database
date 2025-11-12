package customer

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomersResponse represents a customers response
type CustomersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerCode *string `json:"customer_code"`
	
	FirstName *string `json:"first_name"`
	
	LastName *string `json:"last_name"`
	
	CompanyName *string `json:"company_name"`
	
	Email *string `json:"email"`
	
	Phone *string `json:"phone"`
	
	AlternatePhone *string `json:"alternate_phone"`
	
	AddressLine1 *string `json:"address_line1"`
	
	AddressLine2 *string `json:"address_line2"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	DateOfBirth *time.Time `json:"date_of_birth"`
	
	Gender *string `json:"gender"`
	
	TaxNumber *string `json:"tax_number"`
	
	LoyaltyPoints *int64 `json:"loyalty_points"`
	
	LoyaltyTier *string `json:"loyalty_tier"`
	
	CreditLimit *float64 `json:"credit_limit"`
	
	OutstandingBalance *float64 `json:"outstanding_balance"`
	
	TotalPurchases *float64 `json:"total_purchases"`
	
	TotalOrders *int64 `json:"total_orders"`
	
	LastPurchaseAt *time.Time `json:"last_purchase_at"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CustomFields json.RawMessage `json:"custom_fields"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateCustomersRequest represents a request to create a customers
type CreateCustomersRequest struct {
	
	CustomerCode *string `json:"customer_code"`
	
	FirstName *string `json:"first_name"`
	
	LastName *string `json:"last_name"`
	
	CompanyName *string `json:"company_name"`
	
	Email *string `json:"email" validate:"email"`
	
	Phone *string `json:"phone" validate:"e164"`
	
	AlternatePhone *string `json:"alternate_phone" validate:"e164"`
	
	AddressLine1 *string `json:"address_line1"`
	
	AddressLine2 *string `json:"address_line2"`
	
	City *string `json:"city"`
	
	State *string `json:"state"`
	
	Country *string `json:"country"`
	
	PostalCode *string `json:"postal_code"`
	
	DateOfBirth *time.Time `json:"date_of_birth"`
	
	Gender *string `json:"gender"`
	
	TaxNumber *string `json:"tax_number"`
	
	LoyaltyPoints *int64 `json:"loyalty_points"`
	
	LoyaltyTier *string `json:"loyalty_tier"`
	
	CreditLimit *float64 `json:"credit_limit"`
	
	OutstandingBalance *float64 `json:"outstanding_balance"`
	
	TotalPurchases *float64 `json:"total_purchases"`
	
	TotalOrders *int64 `json:"total_orders"`
	
	LastPurchaseAt *time.Time `json:"last_purchase_at"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CustomFields json.RawMessage `json:"custom_fields"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateCustomersRequest) Validate() error {
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCustomersRequest represents a request to update a customers
type UpdateCustomersRequest struct {
	
	CustomerCode *string `json:"customer_code,omitempty"`
	
	FirstName *string `json:"first_name,omitempty"`
	
	LastName *string `json:"last_name,omitempty"`
	
	CompanyName *string `json:"company_name,omitempty"`
	
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	
	Phone *string `json:"phone,omitempty" validate:"omitempty,e164"`
	
	AlternatePhone *string `json:"alternate_phone,omitempty" validate:"omitempty,e164"`
	
	AddressLine1 *string `json:"address_line1,omitempty"`
	
	AddressLine2 *string `json:"address_line2,omitempty"`
	
	City *string `json:"city,omitempty"`
	
	State *string `json:"state,omitempty"`
	
	Country *string `json:"country,omitempty"`
	
	PostalCode *string `json:"postal_code,omitempty"`
	
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	
	Gender *string `json:"gender,omitempty"`
	
	TaxNumber *string `json:"tax_number,omitempty"`
	
	LoyaltyPoints *int64 `json:"loyalty_points,omitempty"`
	
	LoyaltyTier *string `json:"loyalty_tier,omitempty"`
	
	CreditLimit *float64 `json:"credit_limit,omitempty"`
	
	OutstandingBalance *float64 `json:"outstanding_balance,omitempty"`
	
	TotalPurchases *float64 `json:"total_purchases,omitempty"`
	
	TotalOrders *int64 `json:"total_orders,omitempty"`
	
	LastPurchaseAt *time.Time `json:"last_purchase_at,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CustomFields *json.RawMessage `json:"custom_fields,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCustomersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerCode != nil {
		hasUpdate = true
	}
	
	if r.FirstName != nil {
		hasUpdate = true
	}
	
	if r.LastName != nil {
		hasUpdate = true
	}
	
	if r.CompanyName != nil {
		hasUpdate = true
	}
	
	if r.Email != nil {
		hasUpdate = true
	}
	
	if r.Phone != nil {
		hasUpdate = true
	}
	
	if r.AlternatePhone != nil {
		hasUpdate = true
	}
	
	if r.AddressLine1 != nil {
		hasUpdate = true
	}
	
	if r.AddressLine2 != nil {
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
	
	if r.DateOfBirth != nil {
		hasUpdate = true
	}
	
	if r.Gender != nil {
		hasUpdate = true
	}
	
	if r.TaxNumber != nil {
		hasUpdate = true
	}
	
	if r.LoyaltyPoints != nil {
		hasUpdate = true
	}
	
	if r.LoyaltyTier != nil {
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
	
	if r.LastPurchaseAt != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CustomFields != nil {
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

// CustomersListResponse represents a paginated list of customers records
type CustomersListResponse struct {
	Items      []*CustomersResponse `json:"items"`
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

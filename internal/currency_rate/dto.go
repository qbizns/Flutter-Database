package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CurrencyRatesResponse represents a currency_rates response
type CurrencyRatesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CurrencyCode string `json:"currency_code"`
	
	RateDate time.Time `json:"rate_date"`
	
	Rate float64 `json:"rate"`
	
	Source *string `json:"source"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCurrencyRatesRequest represents a request to create a currency_rates
type CreateCurrencyRatesRequest struct {
	
	CurrencyCode string `json:"currency_code" validate:"required"`
	
	RateDate time.Time `json:"rate_date" validate:"required"`
	
	Rate float64 `json:"rate" validate:"required"`
	
	Source *string `json:"source"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateCurrencyRatesRequest) Validate() error {
	
	if r.CurrencyCode == "" {
		return fmt.Errorf("currency_code is required")
	}
	
	if r.RateDate == nil {
		return fmt.Errorf("rate_date is required")
	}
	
	if r.Rate == nil {
		return fmt.Errorf("rate is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCurrencyRatesRequest represents a request to update a currency_rates
type UpdateCurrencyRatesRequest struct {
	
	CurrencyCode *string `json:"currency_code,omitempty" validate:"omitempty,required"`
	
	RateDate *time.Time `json:"rate_date,omitempty" validate:"omitempty,required"`
	
	Rate *float64 `json:"rate,omitempty" validate:"omitempty,required"`
	
	Source *string `json:"source,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCurrencyRatesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CurrencyCode != nil {
		hasUpdate = true
	}
	
	if r.RateDate != nil {
		hasUpdate = true
	}
	
	if r.Rate != nil {
		hasUpdate = true
	}
	
	if r.Source != nil {
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

// CurrencyRatesListResponse represents a paginated list of currency_rates records
type CurrencyRatesListResponse struct {
	Items      []*CurrencyRatesResponse `json:"items"`
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

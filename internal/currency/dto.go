package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CurrenciesResponse represents a currencies response
type CurrenciesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	CurrencyCode string `json:"currency_code"`
	
	CurrencyName string `json:"currency_name"`
	
	CurrencySymbol *string `json:"currency_symbol"`
	
	DecimalPlaces *int64 `json:"decimal_places"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCurrenciesRequest represents a request to create a currencies
type CreateCurrenciesRequest struct {
	
	CurrencyCode string `json:"currency_code" validate:"required"`
	
	CurrencyName string `json:"currency_name" validate:"required"`
	
	CurrencySymbol *string `json:"currency_symbol"`
	
	DecimalPlaces *int64 `json:"decimal_places"`
	
	IsActive *bool `json:"is_active"`
	
}

// Validate validates the create request
func (r *CreateCurrenciesRequest) Validate() error {
	
	if r.CurrencyCode == "" {
		return fmt.Errorf("currency_code is required")
	}
	
	if r.CurrencyName == "" {
		return fmt.Errorf("currency_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCurrenciesRequest represents a request to update a currencies
type UpdateCurrenciesRequest struct {
	
	CurrencyCode *string `json:"currency_code,omitempty" validate:"omitempty,required"`
	
	CurrencyName *string `json:"currency_name,omitempty" validate:"omitempty,required"`
	
	CurrencySymbol *string `json:"currency_symbol,omitempty"`
	
	DecimalPlaces *int64 `json:"decimal_places,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCurrenciesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CurrencyCode != nil {
		hasUpdate = true
	}
	
	if r.CurrencyName != nil {
		hasUpdate = true
	}
	
	if r.CurrencySymbol != nil {
		hasUpdate = true
	}
	
	if r.DecimalPlaces != nil {
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

// CurrenciesListResponse represents a paginated list of currencies records
type CurrenciesListResponse struct {
	Items      []*CurrenciesResponse `json:"items"`
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

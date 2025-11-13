package gift_card

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GiftCardsResponse represents a gift_cards response
type GiftCardsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CardNumber string `json:"card_number"`
	
	PinCode *string `json:"pin_code"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	OriginalValue float64 `json:"original_value"`
	
	CurrentBalance float64 `json:"current_balance"`
	
	IssuedDate time.Time `json:"issued_date"`
	
	ExpiryDate *time.Time `json:"expiry_date"`
	
	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	IssuedByUserId *uuid.UUID `json:"issued_by_user_id"`
	
	IssuedLocationId *uuid.UUID `json:"issued_location_id"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateGiftCardsRequest represents a request to create a gift_cards
type CreateGiftCardsRequest struct {
	
	CardNumber string `json:"card_number" validate:"required"`
	
	PinCode *string `json:"pin_code"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	OriginalValue float64 `json:"original_value" validate:"required"`
	
	CurrentBalance float64 `json:"current_balance" validate:"required"`
	
	IssuedDate time.Time `json:"issued_date" validate:"required"`
	
	ExpiryDate *time.Time `json:"expiry_date"`
	
	// 	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	IssuedByUserId *uuid.UUID `json:"issued_by_user_id"`
	
	IssuedLocationId *uuid.UUID `json:"issued_location_id"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateGiftCardsRequest) Validate() error {
	
	if r.CardNumber == "" {
		return fmt.Errorf("card_number is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	
	if r.IssuedDate.IsZero() {
		return fmt.Errorf("issued_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateGiftCardsRequest represents a request to update a gift_cards
type UpdateGiftCardsRequest struct {
	
	CardNumber *string `json:"card_number,omitempty" validate:"omitempty,required"`
	
	PinCode *string `json:"pin_code,omitempty"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty"`
	
	OriginalValue *float64 `json:"original_value,omitempty" validate:"omitempty,required"`
	
	CurrentBalance *float64 `json:"current_balance,omitempty" validate:"omitempty,required"`
	
	IssuedDate *time.Time `json:"issued_date,omitempty" validate:"omitempty,required"`
	
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	IssuedByUserId *uuid.UUID `json:"issued_by_user_id,omitempty"`
	
	IssuedLocationId *uuid.UUID `json:"issued_location_id,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateGiftCardsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CardNumber != nil {
		hasUpdate = true
	}
	
	if r.PinCode != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.OriginalValue != nil {
		hasUpdate = true
	}
	
	if r.CurrentBalance != nil {
		hasUpdate = true
	}
	
	if r.IssuedDate != nil {
		hasUpdate = true
	}
	
	if r.ExpiryDate != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.IssuedByUserId != nil {
		hasUpdate = true
	}
	
	if r.IssuedLocationId != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
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

// GiftCardsListResponse represents a paginated list of gift_cards records
type GiftCardsListResponse struct {
	Items      []*GiftCardsResponse `json:"items"`
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

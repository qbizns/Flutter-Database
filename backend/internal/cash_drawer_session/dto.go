package cash_drawer_session

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CashDrawerSessionsResponse represents a cash_drawer_sessions response
type CashDrawerSessionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CashDrawerId uuid.UUID `json:"cash_drawer_id"`
	
	PosSessionId uuid.UUID `json:"pos_session_id"`
	
	OpeningAmount *float64 `json:"opening_amount"`
	
	ClosingAmount *float64 `json:"closing_amount"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCashDrawerSessionsRequest represents a request to create a cash_drawer_sessions
type CreateCashDrawerSessionsRequest struct {
	
	CashDrawerId uuid.UUID `json:"cash_drawer_id" validate:"required"`
	
	PosSessionId uuid.UUID `json:"pos_session_id" validate:"required"`
	
	OpeningAmount *float64 `json:"opening_amount"`
	
	ClosingAmount *float64 `json:"closing_amount"`
	
}

// Validate validates the create request
func (r *CreateCashDrawerSessionsRequest) Validate() error {
	
	if r.CashDrawerId == uuid.Nil {
		return fmt.Errorf("cash_drawer_id is required")
	}
	
	if r.PosSessionId == uuid.Nil {
		return fmt.Errorf("pos_session_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCashDrawerSessionsRequest represents a request to update a cash_drawer_sessions
type UpdateCashDrawerSessionsRequest struct {
	
	CashDrawerId *uuid.UUID `json:"cash_drawer_id,omitempty" validate:"omitempty,required"`
	
	PosSessionId *uuid.UUID `json:"pos_session_id,omitempty" validate:"omitempty,required"`
	
	OpeningAmount *float64 `json:"opening_amount,omitempty"`
	
	ClosingAmount *float64 `json:"closing_amount,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCashDrawerSessionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CashDrawerId != nil {
		hasUpdate = true
	}
	
	if r.PosSessionId != nil {
		hasUpdate = true
	}
	
	if r.OpeningAmount != nil {
		hasUpdate = true
	}
	
	if r.ClosingAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CashDrawerSessionsListResponse represents a paginated list of cash_drawer_sessions records
type CashDrawerSessionsListResponse struct {
	Items      []*CashDrawerSessionsResponse `json:"items"`
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

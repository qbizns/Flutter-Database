package cash_movement

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CashMovementsResponse represents a cash_movements response
type CashMovementsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PosSessionId uuid.UUID `json:"pos_session_id"`
	
	CashDrawerId *uuid.UUID `json:"cash_drawer_id"`
	
	MovementType string `json:"movement_type"`
	
	Amount float64 `json:"amount"`
	
	ReasonCode *string `json:"reason_code"`
	
	ReasonDescription string `json:"reason_description"`
	
	UserId uuid.UUID `json:"user_id"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCashMovementsRequest represents a request to create a cash_movements
type CreateCashMovementsRequest struct {
	
	PosSessionId uuid.UUID `json:"pos_session_id" validate:"required"`
	
	CashDrawerId *uuid.UUID `json:"cash_drawer_id"`
	
	MovementType string `json:"movement_type" validate:"required"`
	
	Amount float64 `json:"amount" validate:"required"`
	
	ReasonCode *string `json:"reason_code"`
	
	ReasonDescription string `json:"reason_description" validate:"required"`
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
}

// Validate validates the create request
func (r *CreateCashMovementsRequest) Validate() error {
	
	if r.PosSessionId == uuid.Nil {
		return fmt.Errorf("pos_session_id is required")
	}
	
	if r.MovementType == "" {
		return fmt.Errorf("movement_type is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	
	if r.ReasonDescription == "" {
		return fmt.Errorf("reason_description is required")
	}
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCashMovementsRequest represents a request to update a cash_movements
type UpdateCashMovementsRequest struct {
	
	PosSessionId *uuid.UUID `json:"pos_session_id,omitempty" validate:"omitempty,required"`
	
	CashDrawerId *uuid.UUID `json:"cash_drawer_id,omitempty"`
	
	MovementType *string `json:"movement_type,omitempty" validate:"omitempty,required"`
	
	Amount *float64 `json:"amount,omitempty" validate:"omitempty,required"`
	
	ReasonCode *string `json:"reason_code,omitempty"`
	
	ReasonDescription *string `json:"reason_description,omitempty" validate:"omitempty,required"`
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	RequiresApproval *bool `json:"requires_approval,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCashMovementsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PosSessionId != nil {
		hasUpdate = true
	}
	
	if r.CashDrawerId != nil {
		hasUpdate = true
	}
	
	if r.MovementType != nil {
		hasUpdate = true
	}
	
	if r.Amount != nil {
		hasUpdate = true
	}
	
	if r.ReasonCode != nil {
		hasUpdate = true
	}
	
	if r.ReasonDescription != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.RequiresApproval != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedAt != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CashMovementsListResponse represents a paginated list of cash_movements records
type CashMovementsListResponse struct {
	Items      []*CashMovementsResponse `json:"items"`
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

package stock_adjustment_reason

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StockAdjustmentReasonsResponse represents a stock_adjustment_reasons response
type StockAdjustmentReasonsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	Code string `json:"code"`
	
	Name string `json:"name"`
	
	Description *string `json:"description"`
	
	ReasonType string `json:"reason_type"`
	
	IsSystemReason *bool `json:"is_system_reason"`
	
	IsActive *bool `json:"is_active"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	RequiresNotes *bool `json:"requires_notes"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateStockAdjustmentReasonsRequest represents a request to create a stock_adjustment_reasons
type CreateStockAdjustmentReasonsRequest struct {
	
	Code string `json:"code" validate:"required"`
	
	Name string `json:"name" validate:"required"`
	
	Description *string `json:"description"`
	
	ReasonType string `json:"reason_type" validate:"required"`
	
	IsSystemReason *bool `json:"is_system_reason"`
	
	IsActive *bool `json:"is_active"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	RequiresNotes *bool `json:"requires_notes"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateStockAdjustmentReasonsRequest) Validate() error {
	
	if r.Code == "" {
		return fmt.Errorf("code is required")
	}
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.ReasonType == "" {
		return fmt.Errorf("reason_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateStockAdjustmentReasonsRequest represents a request to update a stock_adjustment_reasons
type UpdateStockAdjustmentReasonsRequest struct {
	
	Code *string `json:"code,omitempty" validate:"omitempty,required"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	ReasonType *string `json:"reason_type,omitempty" validate:"omitempty,required"`
	
	IsSystemReason *bool `json:"is_system_reason,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	RequiresApproval *bool `json:"requires_approval,omitempty"`
	
	RequiresNotes *bool `json:"requires_notes,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateStockAdjustmentReasonsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Code != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.ReasonType != nil {
		hasUpdate = true
	}
	
	if r.IsSystemReason != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.RequiresApproval != nil {
		hasUpdate = true
	}
	
	if r.RequiresNotes != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// StockAdjustmentReasonsListResponse represents a paginated list of stock_adjustment_reasons records
type StockAdjustmentReasonsListResponse struct {
	Items      []*StockAdjustmentReasonsResponse `json:"items"`
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

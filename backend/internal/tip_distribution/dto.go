package tip_distribution

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TipDistributionsResponse represents a tip_distributions response
type TipDistributionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	TipPoolId *uuid.UUID `json:"tip_pool_id"`
	
	DistributionDate time.Time `json:"distribution_date"`
	
	PeriodStart *time.Time `json:"period_start"`
	
	PeriodEnd *time.Time `json:"period_end"`
	
	ShiftId *uuid.UUID `json:"shift_id"`
	
	EmployeeId uuid.UUID `json:"employee_id"`
	
	SourceType string `json:"source_type"`
	
	SourceSaleId *uuid.UUID `json:"source_sale_id"`
	
	SourceOrderId *uuid.UUID `json:"source_order_id"`
	
	TipAmount float64 `json:"tip_amount"`
	
	DistributionAmount float64 `json:"distribution_amount"`
	
	DistributionPercentage *float64 `json:"distribution_percentage"`
	
	PaymentStatus *string `json:"payment_status"`
	
	PaymentMethod *string `json:"payment_method"`
	
	PaidAt *time.Time `json:"paid_at"`
	
	PaidBy *uuid.UUID `json:"paid_by"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'sale', *string `json:"'sale',"`
	
	'pending', *string `json:"'pending',"`
	
	TipAmount *string `json:"tip_amount"`
	
}

// CreateTipDistributionsRequest represents a request to create a tip_distributions
type CreateTipDistributionsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	TipPoolId *uuid.UUID `json:"tip_pool_id"`
	
	DistributionDate time.Time `json:"distribution_date" validate:"required"`
	
	PeriodStart *time.Time `json:"period_start"`
	
	PeriodEnd *time.Time `json:"period_end"`
	
	ShiftId *uuid.UUID `json:"shift_id"`
	
	EmployeeId uuid.UUID `json:"employee_id" validate:"required"`
	
	SourceType string `json:"source_type" validate:"required"`
	
	SourceSaleId *uuid.UUID `json:"source_sale_id"`
	
	SourceOrderId *uuid.UUID `json:"source_order_id"`
	
	TipAmount float64 `json:"tip_amount" validate:"required"`
	
	DistributionAmount float64 `json:"distribution_amount" validate:"required"`
	
	DistributionPercentage *float64 `json:"distribution_percentage"`
	
	// 	PaymentStatus *string `json:"payment_status"`
	
	PaymentMethod *string `json:"payment_method"`
	
	PaidAt *time.Time `json:"paid_at"`
	
	PaidBy *uuid.UUID `json:"paid_by"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'sale', *string `json:"'sale',"`
	
	'pending', *string `json:"'pending',"`
	
	TipAmount *string `json:"tip_amount"`
	
}

// Validate validates the create request
func (r *CreateTipDistributionsRequest) Validate() error {
	
	if r.DistributionDate.IsZero() {
		return fmt.Errorf("distribution_date is required")
	}
	
	if r.EmployeeId == uuid.Nil {
		return fmt.Errorf("employee_id is required")
	}
	
	if r.SourceType == "" {
		return fmt.Errorf("source_type is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTipDistributionsRequest represents a request to update a tip_distributions
type UpdateTipDistributionsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	TipPoolId *uuid.UUID `json:"tip_pool_id,omitempty"`
	
	DistributionDate *time.Time `json:"distribution_date,omitempty" validate:"omitempty,required"`
	
	PeriodStart *time.Time `json:"period_start,omitempty"`
	
	PeriodEnd *time.Time `json:"period_end,omitempty"`
	
	ShiftId *uuid.UUID `json:"shift_id,omitempty"`
	
	EmployeeId *uuid.UUID `json:"employee_id,omitempty" validate:"omitempty,required"`
	
	SourceType *string `json:"source_type,omitempty" validate:"omitempty,required"`
	
	SourceSaleId *uuid.UUID `json:"source_sale_id,omitempty"`
	
	SourceOrderId *uuid.UUID `json:"source_order_id,omitempty"`
	
	TipAmount *float64 `json:"tip_amount,omitempty" validate:"omitempty,required"`
	
	DistributionAmount *float64 `json:"distribution_amount,omitempty" validate:"omitempty,required"`
	
	DistributionPercentage *float64 `json:"distribution_percentage,omitempty"`
	
	// 	PaymentStatus *string `json:"payment_status,omitempty"`
	
	PaymentMethod *string `json:"payment_method,omitempty"`
	
	PaidAt *time.Time `json:"paid_at,omitempty"`
	
	PaidBy *uuid.UUID `json:"paid_by,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'sale', *string `json:"'sale',,omitempty"`
	
	'pending', *string `json:"'pending',,omitempty"`
	
	TipAmount *string `json:"tip_amount,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTipDistributionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.TipPoolId != nil {
		hasUpdate = true
	}
	
	if r.DistributionDate != nil {
		hasUpdate = true
	}
	
	if r.PeriodStart != nil {
		hasUpdate = true
	}
	
	if r.PeriodEnd != nil {
		hasUpdate = true
	}
	
	if r.ShiftId != nil {
		hasUpdate = true
	}
	
	if r.EmployeeId != nil {
		hasUpdate = true
	}
	
	if r.SourceType != nil {
		hasUpdate = true
	}
	
	if r.SourceSaleId != nil {
		hasUpdate = true
	}
	
	if r.SourceOrderId != nil {
		hasUpdate = true
	}
	
	if r.TipAmount != nil {
		hasUpdate = true
	}
	
	if r.DistributionAmount != nil {
		hasUpdate = true
	}
	
	if r.DistributionPercentage != nil {
		hasUpdate = true
	}
	
	if r.PaymentStatus != nil {
		hasUpdate = true
	}
	
	if r.PaymentMethod != nil {
		hasUpdate = true
	}
	
	if r.PaidAt != nil {
		hasUpdate = true
	}
	
	if r.PaidBy != nil {
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
	
	if r.'sale', != nil {
		hasUpdate = true
	}
	
	if r.'pending', != nil {
		hasUpdate = true
	}
	
	if r.TipAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// TipDistributionsListResponse represents a paginated list of tip_distributions records
type TipDistributionsListResponse struct {
	Items      []*TipDistributionsResponse `json:"items"`
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

package staff_commission

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StaffCommissionsResponse represents a staff_commissions response
type StaffCommissionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	EmployeeId uuid.UUID `json:"employee_id"`
	
	CommissionDate time.Time `json:"commission_date"`
	
	PeriodStart time.Time `json:"period_start"`
	
	PeriodEnd time.Time `json:"period_end"`
	
	SourceType string `json:"source_type"`
	
	SourceSaleId *uuid.UUID `json:"source_sale_id"`
	
	SourceOrderId *uuid.UUID `json:"source_order_id"`
	
	CommissionType *string `json:"commission_type"`
	
	CommissionRate *float64 `json:"commission_rate"`
	
	SalesAmount *float64 `json:"sales_amount"`
	
	CommissionAmount float64 `json:"commission_amount"`
	
	Status *string `json:"status"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	PaymentDate *time.Time `json:"payment_date"`
	
	PaymentMethod *string `json:"payment_method"`
	
	PaidBy *uuid.UUID `json:"paid_by"`
	
	Notes *string `json:"notes"`
	
	CalculationNotes *string `json:"calculation_notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'sale', *string `json:"'sale',"`
	
	'percentage', *string `json:"'percentage',"`
	
	'pending', *string `json:"'pending',"`
	
	SalesAmount *string `json:"sales_amount"`
	
}

// CreateStaffCommissionsRequest represents a request to create a staff_commissions
type CreateStaffCommissionsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	EmployeeId uuid.UUID `json:"employee_id" validate:"required"`
	
	CommissionDate time.Time `json:"commission_date" validate:"required"`
	
	PeriodStart time.Time `json:"period_start" validate:"required"`
	
	PeriodEnd time.Time `json:"period_end" validate:"required"`
	
	SourceType string `json:"source_type" validate:"required"`
	
	SourceSaleId *uuid.UUID `json:"source_sale_id"`
	
	SourceOrderId *uuid.UUID `json:"source_order_id"`
	
	CommissionType *string `json:"commission_type"`
	
	CommissionRate *float64 `json:"commission_rate"`
	
	SalesAmount *float64 `json:"sales_amount"`
	
	CommissionAmount float64 `json:"commission_amount" validate:"required"`
	
	// 	Status *string `json:"status"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	PaymentDate *time.Time `json:"payment_date"`
	
	PaymentMethod *string `json:"payment_method"`
	
	PaidBy *uuid.UUID `json:"paid_by"`
	
	Notes *string `json:"notes"`
	
	CalculationNotes *string `json:"calculation_notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'sale', *string `json:"'sale',"`
	
	'percentage', *string `json:"'percentage',"`
	
	'pending', *string `json:"'pending',"`
	
	SalesAmount *string `json:"sales_amount"`
	
}

// Validate validates the create request
func (r *CreateStaffCommissionsRequest) Validate() error {
	
	if r.EmployeeId == uuid.Nil {
		return fmt.Errorf("employee_id is required")
	}
	
	if r.CommissionDate.IsZero() {
		return fmt.Errorf("commission_date is required")
	}
	
	if r.PeriodStart == nil {
		return fmt.Errorf("period_start is required")
	}
	
	if r.PeriodEnd == nil {
		return fmt.Errorf("period_end is required")
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

// UpdateStaffCommissionsRequest represents a request to update a staff_commissions
type UpdateStaffCommissionsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	EmployeeId *uuid.UUID `json:"employee_id,omitempty" validate:"omitempty,required"`
	
	CommissionDate *time.Time `json:"commission_date,omitempty" validate:"omitempty,required"`
	
	PeriodStart *time.Time `json:"period_start,omitempty" validate:"omitempty,required"`
	
	PeriodEnd *time.Time `json:"period_end,omitempty" validate:"omitempty,required"`
	
	SourceType *string `json:"source_type,omitempty" validate:"omitempty,required"`
	
	SourceSaleId *uuid.UUID `json:"source_sale_id,omitempty"`
	
	SourceOrderId *uuid.UUID `json:"source_order_id,omitempty"`
	
	CommissionType *string `json:"commission_type,omitempty"`
	
	CommissionRate *float64 `json:"commission_rate,omitempty"`
	
	SalesAmount *float64 `json:"sales_amount,omitempty"`
	
	CommissionAmount *float64 `json:"commission_amount,omitempty" validate:"omitempty,required"`
	
	// 	Status *string `json:"status,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	PaymentDate *time.Time `json:"payment_date,omitempty"`
	
	PaymentMethod *string `json:"payment_method,omitempty"`
	
	PaidBy *uuid.UUID `json:"paid_by,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CalculationNotes *string `json:"calculation_notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'sale', *string `json:"'sale',,omitempty"`
	
	'percentage', *string `json:"'percentage',,omitempty"`
	
	'pending', *string `json:"'pending',,omitempty"`
	
	SalesAmount *string `json:"sales_amount,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateStaffCommissionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.EmployeeId != nil {
		hasUpdate = true
	}
	
	if r.CommissionDate != nil {
		hasUpdate = true
	}
	
	if r.PeriodStart != nil {
		hasUpdate = true
	}
	
	if r.PeriodEnd != nil {
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
	
	if r.CommissionType != nil {
		hasUpdate = true
	}
	
	if r.CommissionRate != nil {
		hasUpdate = true
	}
	
	if r.SalesAmount != nil {
		hasUpdate = true
	}
	
	if r.CommissionAmount != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedAt != nil {
		hasUpdate = true
	}
	
	if r.PaymentDate != nil {
		hasUpdate = true
	}
	
	if r.PaymentMethod != nil {
		hasUpdate = true
	}
	
	if r.PaidBy != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CalculationNotes != nil {
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
	
	if r.'percentage', != nil {
		hasUpdate = true
	}
	
	if r.'pending', != nil {
		hasUpdate = true
	}
	
	if r.SalesAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// StaffCommissionsListResponse represents a paginated list of staff_commissions records
type StaffCommissionsListResponse struct {
	Items      []*StaffCommissionsResponse `json:"items"`
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

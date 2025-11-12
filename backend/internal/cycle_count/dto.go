package cycle_count

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CycleCountsResponse represents a cycle_counts response
type CycleCountsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	CountNumber string `json:"count_number"`
	
	CountDate time.Time `json:"count_date"`
	
	CountType *string `json:"count_type"`
	
	Status *string `json:"status"`
	
	CategoryId *uuid.UUID `json:"category_id"`
	
	IncludeZeroStock *bool `json:"include_zero_stock"`
	
	TotalItemsPlanned *int64 `json:"total_items_planned"`
	
	TotalItemsCounted *int64 `json:"total_items_counted"`
	
	ItemsWithVariance *int64 `json:"items_with_variance"`
	
	TotalVarianceValue *float64 `json:"total_variance_value"`
	
	ScheduledDate *time.Time `json:"scheduled_date"`
	
	StartedAt *time.Time `json:"started_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CountedBy *uuid.UUID `json:"counted_by"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	TotalItemsPlanned *string `json:"total_items_planned"`
	
	TotalItemsCounted *string `json:"total_items_counted"`
	
	ItemsWithVariance *string `json:"items_with_variance"`
	
}

// CreateCycleCountsRequest represents a request to create a cycle_counts
type CreateCycleCountsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	CountNumber string `json:"count_number" validate:"required"`
	
	CountDate time.Time `json:"count_date" validate:"required"`
	
	CountType *string `json:"count_type"`
	
	Status *string `json:"status"`
	
	CategoryId *uuid.UUID `json:"category_id"`
	
	IncludeZeroStock *bool `json:"include_zero_stock"`
	
	TotalItemsPlanned *int64 `json:"total_items_planned"`
	
	TotalItemsCounted *int64 `json:"total_items_counted"`
	
	ItemsWithVariance *int64 `json:"items_with_variance"`
	
	TotalVarianceValue *float64 `json:"total_variance_value"`
	
	ScheduledDate *time.Time `json:"scheduled_date"`
	
	StartedAt *time.Time `json:"started_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CountedBy *uuid.UUID `json:"counted_by"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	TotalItemsPlanned *string `json:"total_items_planned"`
	
	TotalItemsCounted *string `json:"total_items_counted"`
	
	ItemsWithVariance *string `json:"items_with_variance"`
	
}

// Validate validates the create request
func (r *CreateCycleCountsRequest) Validate() error {
	
	if r.CountNumber == "" {
		return fmt.Errorf("count_number is required")
	}
	
	if r.CountDate == nil {
		return fmt.Errorf("count_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCycleCountsRequest represents a request to update a cycle_counts
type UpdateCycleCountsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	CountNumber *string `json:"count_number,omitempty" validate:"omitempty,required"`
	
	CountDate *time.Time `json:"count_date,omitempty" validate:"omitempty,required"`
	
	CountType *string `json:"count_type,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	CategoryId *uuid.UUID `json:"category_id,omitempty"`
	
	IncludeZeroStock *bool `json:"include_zero_stock,omitempty"`
	
	TotalItemsPlanned *int64 `json:"total_items_planned,omitempty"`
	
	TotalItemsCounted *int64 `json:"total_items_counted,omitempty"`
	
	ItemsWithVariance *int64 `json:"items_with_variance,omitempty"`
	
	TotalVarianceValue *float64 `json:"total_variance_value,omitempty"`
	
	ScheduledDate *time.Time `json:"scheduled_date,omitempty"`
	
	StartedAt *time.Time `json:"started_at,omitempty"`
	
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	CountedBy *uuid.UUID `json:"counted_by,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	TotalItemsPlanned *string `json:"total_items_planned,omitempty"`
	
	TotalItemsCounted *string `json:"total_items_counted,omitempty"`
	
	ItemsWithVariance *string `json:"items_with_variance,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCycleCountsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.CountNumber != nil {
		hasUpdate = true
	}
	
	if r.CountDate != nil {
		hasUpdate = true
	}
	
	if r.CountType != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.CategoryId != nil {
		hasUpdate = true
	}
	
	if r.IncludeZeroStock != nil {
		hasUpdate = true
	}
	
	if r.TotalItemsPlanned != nil {
		hasUpdate = true
	}
	
	if r.TotalItemsCounted != nil {
		hasUpdate = true
	}
	
	if r.ItemsWithVariance != nil {
		hasUpdate = true
	}
	
	if r.TotalVarianceValue != nil {
		hasUpdate = true
	}
	
	if r.ScheduledDate != nil {
		hasUpdate = true
	}
	
	if r.StartedAt != nil {
		hasUpdate = true
	}
	
	if r.CompletedAt != nil {
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
	
	if r.CountedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.TotalItemsPlanned != nil {
		hasUpdate = true
	}
	
	if r.TotalItemsCounted != nil {
		hasUpdate = true
	}
	
	if r.ItemsWithVariance != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CycleCountsListResponse represents a paginated list of cycle_counts records
type CycleCountsListResponse struct {
	Items      []*CycleCountsResponse `json:"items"`
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

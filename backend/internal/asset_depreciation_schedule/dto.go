package asset_depreciation_schedule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AssetDepreciationScheduleResponse represents a asset_depreciation_schedule response
type AssetDepreciationScheduleResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	FixedAssetId uuid.UUID `json:"fixed_asset_id"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	DepreciationDate time.Time `json:"depreciation_date"`
	
	DepreciationAmount float64 `json:"depreciation_amount"`
	
	AccumulatedDepreciationBeginning float64 `json:"accumulated_depreciation_beginning"`
	
	AccumulatedDepreciationEnding float64 `json:"accumulated_depreciation_ending"`
	
	BookValueBeginning float64 `json:"book_value_beginning"`
	
	BookValueEnding float64 `json:"book_value_ending"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	IsPosted *bool `json:"is_posted"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	PostedAt *time.Time `json:"posted_at"`
	
	PostedBy *uuid.UUID `json:"posted_by"`
	
}

// CreateAssetDepreciationScheduleRequest represents a request to create a asset_depreciation_schedule
type CreateAssetDepreciationScheduleRequest struct {
	
	FixedAssetId uuid.UUID `json:"fixed_asset_id" validate:"required"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	DepreciationDate time.Time `json:"depreciation_date" validate:"required"`
	
	DepreciationAmount float64 `json:"depreciation_amount" validate:"required"`
	
	AccumulatedDepreciationBeginning float64 `json:"accumulated_depreciation_beginning" validate:"required"`
	
	AccumulatedDepreciationEnding float64 `json:"accumulated_depreciation_ending" validate:"required"`
	
	BookValueBeginning float64 `json:"book_value_beginning" validate:"required"`
	
	BookValueEnding float64 `json:"book_value_ending" validate:"required"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	IsPosted *bool `json:"is_posted"`
	
	PostedAt *time.Time `json:"posted_at"`
	
	PostedBy *uuid.UUID `json:"posted_by"`
	
}

// Validate validates the create request
func (r *CreateAssetDepreciationScheduleRequest) Validate() error {
	
	if r.FixedAssetId == uuid.Nil {
		return fmt.Errorf("fixed_asset_id is required")
	}
	
	if r.DepreciationDate == nil {
		return fmt.Errorf("depreciation_date is required")
	}
	
	if r.DepreciationAmount == nil {
		return fmt.Errorf("depreciation_amount is required")
	}
	
	if r.AccumulatedDepreciationBeginning == nil {
		return fmt.Errorf("accumulated_depreciation_beginning is required")
	}
	
	if r.AccumulatedDepreciationEnding == nil {
		return fmt.Errorf("accumulated_depreciation_ending is required")
	}
	
	if r.BookValueBeginning == nil {
		return fmt.Errorf("book_value_beginning is required")
	}
	
	if r.BookValueEnding == nil {
		return fmt.Errorf("book_value_ending is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateAssetDepreciationScheduleRequest represents a request to update a asset_depreciation_schedule
type UpdateAssetDepreciationScheduleRequest struct {
	
	FixedAssetId *uuid.UUID `json:"fixed_asset_id,omitempty" validate:"omitempty,required"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id,omitempty"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id,omitempty"`
	
	DepreciationDate *time.Time `json:"depreciation_date,omitempty" validate:"omitempty,required"`
	
	DepreciationAmount *float64 `json:"depreciation_amount,omitempty" validate:"omitempty,required"`
	
	AccumulatedDepreciationBeginning *float64 `json:"accumulated_depreciation_beginning,omitempty" validate:"omitempty,required"`
	
	AccumulatedDepreciationEnding *float64 `json:"accumulated_depreciation_ending,omitempty" validate:"omitempty,required"`
	
	BookValueBeginning *float64 `json:"book_value_beginning,omitempty" validate:"omitempty,required"`
	
	BookValueEnding *float64 `json:"book_value_ending,omitempty" validate:"omitempty,required"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty"`
	
	IsPosted *bool `json:"is_posted,omitempty"`
	
	PostedAt *time.Time `json:"posted_at,omitempty"`
	
	PostedBy *uuid.UUID `json:"posted_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAssetDepreciationScheduleRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.FixedAssetId != nil {
		hasUpdate = true
	}
	
	if r.FiscalYearId != nil {
		hasUpdate = true
	}
	
	if r.AccountingPeriodId != nil {
		hasUpdate = true
	}
	
	if r.DepreciationDate != nil {
		hasUpdate = true
	}
	
	if r.DepreciationAmount != nil {
		hasUpdate = true
	}
	
	if r.AccumulatedDepreciationBeginning != nil {
		hasUpdate = true
	}
	
	if r.AccumulatedDepreciationEnding != nil {
		hasUpdate = true
	}
	
	if r.BookValueBeginning != nil {
		hasUpdate = true
	}
	
	if r.BookValueEnding != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.IsPosted != nil {
		hasUpdate = true
	}
	
	if r.PostedAt != nil {
		hasUpdate = true
	}
	
	if r.PostedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// AssetDepreciationScheduleListResponse represents a paginated list of asset_depreciation_schedule records
type AssetDepreciationScheduleListResponse struct {
	Items      []*AssetDepreciationScheduleResponse `json:"items"`
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

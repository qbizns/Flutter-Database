package assets

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DepreciationMethod defines how an asset depreciates
type DepreciationMethod string

const (
	DepreciationMethodStraightLine      DepreciationMethod = "straight_line"
	DepreciationMethodDecliningBalance  DepreciationMethod = "declining_balance"
	DepreciationMethodUnitsOfProduction DepreciationMethod = "units_of_production"
)

// AssetCategory represents the classification of fixed assets
type AssetCategory struct {
	ID                                  uuid.UUID
	CategoryCode                        string
	CategoryName                        string
	DefaultDepreciationMethod           *DepreciationMethod
	DefaultUsefulLifeYears              *int
	DefaultSalvageValuePercent          *float64
	AssetAccountID                      *uuid.UUID
	AccumulatedDepreciationAccountID    *uuid.UUID
	DepreciationExpenseAccountID        *uuid.UUID
	Description                         *string
	CreatedAt                           time.Time
	UpdatedAt                           time.Time
}

// FixedAsset represents a tangible asset owned by the organization
type FixedAsset struct {
	ID                                  uuid.UUID
	OrganizationID                      uuid.UUID
	AssetNumber                         string
	AssetName                           string
	AssetCategoryID                     *uuid.UUID
	AcquisitionDate                     time.Time
	AcquisitionCost                     float64
	SalvageValue                        float64
	SupplierID                          *uuid.UUID
	VendorBillID                        *uuid.UUID
	DepreciationMethod                  DepreciationMethod
	UsefulLifeYears                     int
	DepreciationStartDate               time.Time
	AssetAccountID                      uuid.UUID
	AccumulatedDepreciationAccountID    uuid.UUID
	DepreciationExpenseAccountID        uuid.UUID
	CurrentBookValue                    float64
	AccumulatedDepreciation             float64
	LastDepreciationDate                *time.Time
	LocationID                          *uuid.UUID
	Department                          *string
	IsDisposed                          bool
	DisposalDate                        *time.Time
	DisposalProceeds                    *float64
	DisposalJournalEntryID              *uuid.UUID
	Description                         *string
	SerialNumber                        *string
	Notes                               *string
	Metadata                            json.RawMessage
	CreatedAt                           time.Time
	UpdatedAt                           time.Time
	CreatedBy                           *uuid.UUID
	UpdatedBy                           *uuid.UUID
	DeletedAt                           *time.Time
}

// CreateFixedAssetRequest contains input for creating a fixed asset
type CreateFixedAssetRequest struct {
	AssetNumber              string              `json:"asset_number" validate:"required,max=50"`
	AssetName                string              `json:"asset_name" validate:"required,max=255"`
	AssetCategoryID          *uuid.UUID          `json:"asset_category_id"`
	AcquisitionDate          time.Time           `json:"acquisition_date" validate:"required"`
	AcquisitionCost          float64             `json:"acquisition_cost" validate:"required,gt=0"`
	SalvageValue             float64             `json:"salvage_value" validate:"gte=0"`
	SupplierID               *uuid.UUID          `json:"supplier_id"`
	DepreciationMethod       DepreciationMethod  `json:"depreciation_method" validate:"required,oneof=straight_line declining_balance units_of_production"`
	UsefulLifeYears          int                 `json:"useful_life_years" validate:"required,gt=0"`
	DepreciationStartDate    time.Time           `json:"depreciation_start_date" validate:"required"`
	AssetAccountID           uuid.UUID           `json:"asset_account_id" validate:"required"`
	AccumulatedDepreciationAccountID uuid.UUID  `json:"accumulated_depreciation_account_id" validate:"required"`
	DepreciationExpenseAccountID uuid.UUID      `json:"depreciation_expense_account_id" validate:"required"`
	LocationID               *uuid.UUID          `json:"location_id"`
	Department               *string             `json:"department"`
	Description              *string             `json:"description"`
	SerialNumber             *string             `json:"serial_number"`
	Notes                    *string             `json:"notes"`
	Metadata                 json.RawMessage     `json:"metadata"`
}

// UpdateFixedAssetRequest contains input for updating a fixed asset
type UpdateFixedAssetRequest struct {
	AssetName                *string             `json:"asset_name"`
	AssetCategoryID          *uuid.UUID          `json:"asset_category_id"`
	SalvageValue             *float64            `json:"salvage_value"`
	UsefulLifeYears          *int                `json:"useful_life_years"`
	LocationID               *uuid.UUID          `json:"location_id"`
	Department               *string             `json:"department"`
	Description              *string             `json:"description"`
	Notes                    *string             `json:"notes"`
	Metadata                 json.RawMessage     `json:"metadata"`
}

// DisposeAssetRequest contains input for disposing of an asset
type DisposeAssetRequest struct {
	DisposalDate       time.Time   `json:"disposal_date" validate:"required"`
	DisposalProceeds   float64     `json:"disposal_proceeds" validate:"gte=0"`
	Notes              *string     `json:"notes"`
}

// AssetDepreciationSchedule tracks depreciation for an asset in a period
type AssetDepreciationSchedule struct {
	ID                                  uuid.UUID
	OrganizationID                      uuid.UUID
	FixedAssetID                        uuid.UUID
	FiscalYearID                        *uuid.UUID
	AccountingPeriodID                  *uuid.UUID
	DepreciationDate                    time.Time
	DepreciationAmount                  float64
	AccumulatedDepreciationBeginning    float64
	AccumulatedDepreciationEnding       float64
	BookValueBeginning                  float64
	BookValueEnding                     float64
	JournalEntryID                      *uuid.UUID
	IsPosted                            bool
	CreatedAt                           time.Time
	PostedAt                            *time.Time
	PostedBy                            *uuid.UUID
	DeletedAt                           *time.Time
}

// CreateDepreciationScheduleRequest contains input for creating depreciation
type CreateDepreciationScheduleRequest struct {
	FixedAssetID            uuid.UUID   `json:"fixed_asset_id" validate:"required"`
	FiscalYearID            *uuid.UUID  `json:"fiscal_year_id"`
	AccountingPeriodID      *uuid.UUID  `json:"accounting_period_id"`
	DepreciationDate        time.Time   `json:"depreciation_date" validate:"required"`
	DepreciationAmount      float64     `json:"depreciation_amount" validate:"gte=0"`
}

// UpdateDepreciationScheduleRequest contains input for updating depreciation
type UpdateDepreciationScheduleRequest struct {
	DepreciationAmount *float64 `json:"depreciation_amount"`
}

// AssetQuery filters for fetching assets
type AssetQuery struct {
	OrganizationID  uuid.UUID
	AssetCategoryID *uuid.UUID
	LocationID      *uuid.UUID
	IsDisposed      *bool
	Status          *string // "active", "disposed"
	Limit           int
	Offset          int
}

// DepreciationQuery filters for fetching depreciation schedules
type DepreciationQuery struct {
	OrganizationID     uuid.UUID
	FixedAssetID       *uuid.UUID
	FiscalYearID       *uuid.UUID
	AccountingPeriodID *uuid.UUID
	IsPosted           *bool
	DateFrom           *time.Time
	DateTo             *time.Time
	Limit              int
	Offset             int
}

// AssetStatistics aggregates asset data
type AssetStatistics struct {
	TotalAssets                 int64
	TotalAcquisitionCost        float64
	TotalAccumulatedDepreciation float64
	TotalBookValue              float64
	ActiveAssets                int64
	DisposedAssets              int64
	DepreciatingAssets          int64
}

// CalculateDepreciationInput contains parameters for depreciation calculation
type CalculateDepreciationInput struct {
	FixedAssetID       uuid.UUID
	TargetDate         time.Time
	DepreciationMethod DepreciationMethod
	AcquisitionCost    float64
	SalvageValue       float64
	UsefulLifeYears    int
	DepreciationStart  time.Time
	LastDepreciation   *time.Time
}

// DepreciationResult contains calculated depreciation information
type DepreciationResult struct {
	DepreciationAmount              float64
	AccumulatedDepreciationToDate   float64
	BookValueAtDate                 float64
	RemainingUsefulLife             int
	HasFullyDepreciated             bool
}

// Validate checks if AssetCategory is valid
func (ac *AssetCategory) Validate() error {
	if ac.CategoryCode == "" {
		return ErrCategoryCodeRequired
	}
	if ac.CategoryName == "" {
		return ErrCategoryNameRequired
	}
	return nil
}

// Validate checks if FixedAsset is valid
func (fa *FixedAsset) Validate() error {
	if fa.AssetNumber == "" {
		return ErrAssetNumberRequired
	}
	if fa.AssetName == "" {
		return ErrAssetNameRequired
	}
	if fa.AcquisitionCost < 0 {
		return ErrInvalidAcquisitionCost
	}
	if fa.SalvageValue > fa.AcquisitionCost {
		return ErrSalvageValueExceedsAcquisitionCost
	}
	if fa.UsefulLifeYears <= 0 {
		return ErrInvalidUsefulLife
	}
	if !IsValidDepreciationMethod(fa.DepreciationMethod) {
		return ErrInvalidDepreciationMethod
	}
	return nil
}

// IsValidDepreciationMethod checks if depreciation method is valid
func IsValidDepreciationMethod(method DepreciationMethod) bool {
	switch method {
	case DepreciationMethodStraightLine, DepreciationMethodDecliningBalance, DepreciationMethodUnitsOfProduction:
		return true
	}
	return false
}

// Scan implements sql.Scanner interface
func (dm *DepreciationMethod) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*dm = DepreciationMethod(value.(string))
	return nil
}

// Value implements driver.Valuer interface
func (dm DepreciationMethod) Value() (driver.Value, error) {
	return string(dm), nil
}

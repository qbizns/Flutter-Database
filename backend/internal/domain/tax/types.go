package tax

import (
	"time"

	"github.com/google/uuid"
)

// ========================
// TAX GROUPS
// ========================

type TaxGroup struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	OrganizationID uuid.UUID  `db:"organization_id" json:"organization_id"`
	GroupCode      string     `db:"group_code" json:"group_code"`
	GroupName      string     `db:"group_name" json:"group_name"`
	Sequence       int        `db:"sequence" json:"sequence"`
	IsActive       bool       `db:"is_active" json:"is_active"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateTaxGroupRequest struct {
	GroupCode string `json:"group_code" validate:"required,min=1,max=20"`
	GroupName string `json:"group_name" validate:"required,min=1,max=255"`
	Sequence  int    `json:"sequence" validate:"min=1"`
	IsActive  bool   `json:"is_active"`
}

type UpdateTaxGroupRequest struct {
	GroupCode *string `json:"group_code" validate:"omitempty,min=1,max=20"`
	GroupName *string `json:"group_name" validate:"omitempty,min=1,max=255"`
	Sequence  *int    `json:"sequence" validate:"omitempty,min=1"`
	IsActive  *bool   `json:"is_active"`
}

// ========================
// TAXES
// ========================

type TaxScope string

const (
	TaxScopeSales     TaxScope = "sales"
	TaxScopePurchases TaxScope = "purchases"
	TaxScopeBoth      TaxScope = "both"
)

func (s TaxScope) String() string {
	return string(s)
}

type Tax struct {
	ID                   uuid.UUID  `db:"id" json:"id"`
	OrganizationID       uuid.UUID  `db:"organization_id" json:"organization_id"`
	TaxGroupID           *uuid.UUID `db:"tax_group_id" json:"tax_group_id,omitempty"`
	TaxCode              string     `db:"tax_code" json:"tax_code"`
	TaxName              string     `db:"tax_name" json:"tax_name"`
	TaxRate              float64    `db:"tax_rate" json:"tax_rate"`
	TaxScope             TaxScope   `db:"tax_scope" json:"tax_scope"`
	IsPriceInclusive     bool       `db:"is_price_inclusive" json:"is_price_inclusive"`
	TaxAccountID         uuid.UUID  `db:"tax_account_id" json:"tax_account_id"`
	TaxRefundAccountID   *uuid.UUID `db:"tax_refund_account_id" json:"tax_refund_account_id,omitempty"`
	IsActive             bool       `db:"is_active" json:"is_active"`
	Description          *string    `db:"description" json:"description,omitempty"`
	CreatedBy            *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt            *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateTaxRequest struct {
	TaxCode            string     `json:"tax_code" validate:"required,min=1,max=20"`
	TaxName            string     `json:"tax_name" validate:"required,min=1,max=255"`
	TaxRate            float64    `json:"tax_rate" validate:"required,min=0,max=100"`
	TaxScope           TaxScope   `json:"tax_scope" validate:"required,oneof=sales purchases both"`
	IsPriceInclusive   bool       `json:"is_price_inclusive"`
	TaxAccountID       uuid.UUID  `json:"tax_account_id" validate:"required"`
	TaxRefundAccountID *uuid.UUID `json:"tax_refund_account_id"`
	TaxGroupID         *uuid.UUID `json:"tax_group_id"`
	IsActive           bool       `json:"is_active"`
	Description        *string    `json:"description" validate:"omitempty,max=500"`
}

type UpdateTaxRequest struct {
	TaxCode            *string    `json:"tax_code" validate:"omitempty,min=1,max=20"`
	TaxName            *string    `json:"tax_name" validate:"omitempty,min=1,max=255"`
	TaxRate            *float64   `json:"tax_rate" validate:"omitempty,min=0,max=100"`
	TaxScope           *TaxScope  `json:"tax_scope" validate:"omitempty,oneof=sales purchases both"`
	IsPriceInclusive   *bool      `json:"is_price_inclusive"`
	TaxAccountID       *uuid.UUID `json:"tax_account_id"`
	TaxRefundAccountID *uuid.UUID `json:"tax_refund_account_id"`
	TaxGroupID         *uuid.UUID `json:"tax_group_id"`
	IsActive           *bool      `json:"is_active"`
	Description        *string    `json:"description" validate:"omitempty,max=500"`
}

// ========================
// FISCAL POSITIONS
// ========================

type FiscalPosition struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	OrganizationID     uuid.UUID  `db:"organization_id" json:"organization_id"`
	PositionCode       string     `db:"position_code" json:"position_code"`
	PositionName       string     `db:"position_name" json:"position_name"`
	AutoApply          bool       `db:"auto_apply" json:"auto_apply"`
	CountryID          *string    `db:"country_id" json:"country_id,omitempty"`
	StateProvince      *string    `db:"state_province" json:"state_province,omitempty"`
	ZipPostalCodeRange *string    `db:"zip_postal_code_range" json:"zip_postal_code_range,omitempty"`
	IsActive           bool       `db:"is_active" json:"is_active"`
	Notes              *string    `db:"notes" json:"notes,omitempty"`
	CreatedBy          *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateFiscalPositionRequest struct {
	PositionCode       string  `json:"position_code" validate:"required,min=1,max=20"`
	PositionName       string  `json:"position_name" validate:"required,min=1,max=255"`
	AutoApply          bool    `json:"auto_apply"`
	CountryID          *string `json:"country_id" validate:"omitempty,len=2"`
	StateProvince      *string `json:"state_province" validate:"omitempty,max=100"`
	ZipPostalCodeRange *string `json:"zip_postal_code_range" validate:"omitempty,max=100"`
	IsActive           bool    `json:"is_active"`
	Notes              *string `json:"notes" validate:"omitempty,max=500"`
}

type UpdateFiscalPositionRequest struct {
	PositionCode       *string `json:"position_code" validate:"omitempty,min=1,max=20"`
	PositionName       *string `json:"position_name" validate:"omitempty,min=1,max=255"`
	AutoApply          *bool   `json:"auto_apply"`
	CountryID          *string `json:"country_id" validate:"omitempty,len=2"`
	StateProvince      *string `json:"state_province" validate:"omitempty,max=100"`
	ZipPostalCodeRange *string `json:"zip_postal_code_range" validate:"omitempty,max=100"`
	IsActive           *bool   `json:"is_active"`
	Notes              *string `json:"notes" validate:"omitempty,max=500"`
}

// ========================
// FISCAL POSITION TAX MAPPINGS
// ========================

type FiscalPositionTaxMapping struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	FiscalPositionID   uuid.UUID  `db:"fiscal_position_id" json:"fiscal_position_id"`
	SourceTaxID        uuid.UUID  `db:"source_tax_id" json:"source_tax_id"`
	DestinationTaxID   *uuid.UUID `db:"destination_tax_id" json:"destination_tax_id,omitempty"`
	CreatedBy          *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateFiscalPositionTaxMappingRequest struct {
	SourceTaxID      uuid.UUID  `json:"source_tax_id" validate:"required"`
	DestinationTaxID *uuid.UUID `json:"destination_tax_id"`
}

type UpdateFiscalPositionTaxMappingRequest struct {
	DestinationTaxID *uuid.UUID `json:"destination_tax_id"`
}

// ========================
// TAX REPORT DEFINITIONS
// ========================

type TaxReportFrequency string

const (
	TaxReportFrequencyMonthly   TaxReportFrequency = "monthly"
	TaxReportFrequencyQuarterly TaxReportFrequency = "quarterly"
	TaxReportFrequencyAnnual    TaxReportFrequency = "annual"
	TaxReportFrequencyOnDemand  TaxReportFrequency = "on_demand"
)

func (f TaxReportFrequency) String() string {
	return string(f)
}

type TaxReportDefinition struct {
	ID                       uuid.UUID           `db:"id" json:"id"`
	OrganizationID           *uuid.UUID          `db:"organization_id" json:"organization_id,omitempty"`
	LocalizationPackageID    *uuid.UUID          `db:"localization_package_id" json:"localization_package_id,omitempty"`
	ReportCode               string              `db:"report_code" json:"report_code"`
	ReportName               string              `db:"report_name" json:"report_name"`
	Jurisdiction             *string             `db:"jurisdiction" json:"jurisdiction,omitempty"`
	Authority                *string             `db:"authority" json:"authority,omitempty"`
	ReportFrequency          *TaxReportFrequency `db:"report_frequency" json:"report_frequency,omitempty"`
	Version                  *string             `db:"version" json:"version,omitempty"`
	EffectiveFrom            *time.Time          `db:"effective_from" json:"effective_from,omitempty"`
	EffectiveTo              *time.Time          `db:"effective_to" json:"effective_to,omitempty"`
	IsActive                 bool                `db:"is_active" json:"is_active"`
	Description              *string             `db:"description" json:"description,omitempty"`
	CreatedBy                *uuid.UUID          `db:"created_by" json:"created_by,omitempty"`
	CreatedAt                time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time           `db:"updated_at" json:"updated_at"`
	DeletedAt                *time.Time          `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateTaxReportDefinitionRequest struct {
	ReportCode            string              `json:"report_code" validate:"required,min=1,max=50"`
	ReportName            string              `json:"report_name" validate:"required,min=1,max=255"`
	Jurisdiction          *string             `json:"jurisdiction" validate:"omitempty,max=100"`
	Authority             *string             `json:"authority" validate:"omitempty,max=255"`
	ReportFrequency       *TaxReportFrequency `json:"report_frequency" validate:"omitempty,oneof=monthly quarterly annual on_demand"`
	Version               *string             `json:"version" validate:"omitempty,max=20"`
	EffectiveFrom         *time.Time          `json:"effective_from"`
	EffectiveTo           *time.Time          `json:"effective_to"`
	IsActive              bool                `json:"is_active"`
	Description           *string             `json:"description" validate:"omitempty,max=500"`
	LocalizationPackageID *uuid.UUID          `json:"localization_package_id"`
}

type UpdateTaxReportDefinitionRequest struct {
	ReportCode            *string             `json:"report_code" validate:"omitempty,min=1,max=50"`
	ReportName            *string             `json:"report_name" validate:"omitempty,min=1,max=255"`
	Jurisdiction          *string             `json:"jurisdiction" validate:"omitempty,max=100"`
	Authority             *string             `json:"authority" validate:"omitempty,max=255"`
	ReportFrequency       *TaxReportFrequency `json:"report_frequency" validate:"omitempty,oneof=monthly quarterly annual on_demand"`
	Version               *string             `json:"version" validate:"omitempty,max=20"`
	EffectiveFrom         *time.Time          `json:"effective_from"`
	EffectiveTo           *time.Time          `json:"effective_to"`
	IsActive              *bool               `json:"is_active"`
	Description           *string             `json:"description" validate:"omitempty,max=500"`
	LocalizationPackageID *uuid.UUID          `json:"localization_package_id"`
}

// ========================
// TAX REPORT LINES
// ========================

type TaxReportLineFormulaType string

const (
	TaxReportLineFormulaTypeSum     TaxReportLineFormulaType = "sum"
	TaxReportLineFormulaTypeDetail  TaxReportLineFormulaType = "detail"
	TaxReportLineFormulaTypeFormula TaxReportLineFormulaType = "formula"
	TaxReportLineFormulaTypeManual  TaxReportLineFormulaType = "manual"
)

func (f TaxReportLineFormulaType) String() string {
	return string(f)
}

// UUIDArray is a custom type for []uuid.UUID
// TODO: Implement proper Value/Scan methods if database serialization is needed
/*
type UUIDArray []uuid.UUID

func (ua UUIDArray) Value() (driver.Value, error) {
	return []uuid.UUID(ua).Value()
}

func (ua *UUIDArray) Scan(value interface{}) error {
	return (*[]uuid.UUID)(ua).Scan(value)
}
*/

type TaxReportLine struct {
	ID                   uuid.UUID                  `db:"id" json:"id"`
	TaxReportDefinitionID uuid.UUID                  `db:"tax_report_definition_id" json:"tax_report_definition_id"`
	LineCode             string                     `db:"line_code" json:"line_code"`
	LineName             string                     `db:"line_name" json:"line_name"`
	Sequence             int                        `db:"sequence" json:"sequence"`
	ParentLineID         *uuid.UUID                 `db:"parent_line_id" json:"parent_line_id,omitempty"`
	FormulaType          *TaxReportLineFormulaType  `db:"formula_type" json:"formula_type,omitempty"`
	Formula              *string                    `db:"formula" json:"formula,omitempty"`
	TaxGroupIDs          []uuid.UUID               `db:"tax_group_ids" json:"tax_group_ids"`
	AccountIDs           []uuid.UUID               `db:"account_ids" json:"account_ids"`
	TaxIDs               []uuid.UUID               `db:"tax_ids" json:"tax_ids"`
	IsSubtotal           bool                       `db:"is_subtotal" json:"is_subtotal"`
	IsTotal              bool                       `db:"is_total" json:"is_total"`
	Notes                *string                    `db:"notes" json:"notes,omitempty"`
	CreatedAt            time.Time                  `db:"created_at" json:"created_at"`
	DeletedAt            *time.Time                 `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateTaxReportLineRequest struct {
	LineCode        string                    `json:"line_code" validate:"required,min=1,max=50"`
	LineName        string                    `json:"line_name" validate:"required,min=1,max=255"`
	Sequence        int                       `json:"sequence" validate:"min=1"`
	ParentLineID    *uuid.UUID                `json:"parent_line_id"`
	FormulaType     *TaxReportLineFormulaType `json:"formula_type" validate:"omitempty,oneof=sum detail formula manual"`
	Formula         *string                   `json:"formula" validate:"omitempty,max=1000"`
	TaxGroupIDs     []uuid.UUID               `json:"tax_group_ids"`
	AccountIDs      []uuid.UUID               `json:"account_ids"`
	TaxIDs          []uuid.UUID               `json:"tax_ids"`
	IsSubtotal      bool                      `json:"is_subtotal"`
	IsTotal         bool                      `json:"is_total"`
	Notes           *string                   `json:"notes" validate:"omitempty,max=500"`
}

type UpdateTaxReportLineRequest struct {
	LineCode        *string                   `json:"line_code" validate:"omitempty,min=1,max=50"`
	LineName        *string                   `json:"line_name" validate:"omitempty,min=1,max=255"`
	Sequence        *int                      `json:"sequence" validate:"omitempty,min=1"`
	ParentLineID    *uuid.UUID                `json:"parent_line_id"`
	FormulaType     *TaxReportLineFormulaType `json:"formula_type" validate:"omitempty,oneof=sum detail formula manual"`
	Formula         *string                   `json:"formula" validate:"omitempty,max=1000"`
	TaxGroupIDs     []uuid.UUID               `json:"tax_group_ids"`
	AccountIDs      []uuid.UUID               `json:"account_ids"`
	TaxIDs          []uuid.UUID               `json:"tax_ids"`
	IsSubtotal      *bool                     `json:"is_subtotal"`
	IsTotal         *bool                     `json:"is_total"`
	Notes           *string                   `json:"notes" validate:"omitempty,max=500"`
}

// ========================
// TAX CALCULATION RESULTS
// ========================

type TaxCalculationResult struct {
	TaxID              uuid.UUID `json:"tax_id"`
	TaxCode            string    `json:"tax_code"`
	TaxName            string    `json:"tax_name"`
	TaxRate            float64   `json:"tax_rate"`
	BaseAmount         float64   `json:"base_amount"`
	TaxAmount          float64   `json:"tax_amount"`
	IsPriceInclusive   bool      `json:"is_price_inclusive"`
	TaxAccountID       uuid.UUID `json:"tax_account_id"`
	TaxRefundAccountID *uuid.UUID `json:"tax_refund_account_id,omitempty"`
}

type TaxReport struct {
	ReportID        uuid.UUID           `json:"report_id"`
	ReportCode      string              `json:"report_code"`
	ReportName      string              `json:"report_name"`
	ReportDate      time.Time           `json:"report_date"`
	PeriodStart     time.Time           `json:"period_start"`
	PeriodEnd       time.Time           `json:"period_end"`
	Lines           []TaxReportLineData `json:"lines"`
	TotalTaxAmount  float64             `json:"total_tax_amount"`
	GeneratedAt     time.Time           `json:"generated_at"`
}

type TaxReportLineData struct {
	LineCode       string  `json:"line_code"`
	LineName       string  `json:"line_name"`
	Sequence       int     `json:"sequence"`
	Amount         float64 `json:"amount"`
	IsSubtotal     bool    `json:"is_subtotal"`
	IsTotal        bool    `json:"is_total"`
	ChildLines     []TaxReportLineData `json:"child_lines,omitempty"`
}

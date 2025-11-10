package tax

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTaxNotFound                 = errors.New("tax not found")
	ErrTaxGroupNotFound            = errors.New("tax group not found")
	ErrFiscalPositionNotFound      = errors.New("fiscal position not found")
	ErrTaxMappingNotFound          = errors.New("tax mapping not found")
	ErrTaxReportDefinitionNotFound = errors.New("tax report definition not found")
	ErrTaxReportLineNotFound       = errors.New("tax report line not found")
	ErrInvalidTaxRate              = errors.New("invalid tax rate")
	ErrDuplicateTaxCode            = errors.New("tax code already exists")
	ErrDuplicateGroupCode          = errors.New("group code already exists")
	ErrInvalidDateRange            = errors.New("invalid date range")
)

// TaxRepository defines the interface for tax database operations
type TaxRepository interface {
	// Tax Groups
	CreateTaxGroup(ctx context.Context, tg *TaxGroup) error
	GetTaxGroup(ctx context.Context, id, organizationID uuid.UUID) (*TaxGroup, error)
	GetTaxGroupByCode(ctx context.Context, organizationID uuid.UUID, code string) (*TaxGroup, error)
	ListTaxGroups(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*TaxGroup, error)
	UpdateTaxGroup(ctx context.Context, tg *TaxGroup) error
	DeleteTaxGroup(ctx context.Context, id, organizationID uuid.UUID) error

	// Taxes
	CreateTax(ctx context.Context, tax *Tax) error
	GetTax(ctx context.Context, id, organizationID uuid.UUID) (*Tax, error)
	GetTaxByCode(ctx context.Context, organizationID uuid.UUID, code string) (*Tax, error)
	ListTaxes(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*Tax, error)
	UpdateTax(ctx context.Context, tax *Tax) error
	DeleteTax(ctx context.Context, id, organizationID uuid.UUID) error

	// Fiscal Positions
	CreateFiscalPosition(ctx context.Context, fp *FiscalPosition) error
	GetFiscalPosition(ctx context.Context, id, organizationID uuid.UUID) (*FiscalPosition, error)
	GetFiscalPositionByCode(ctx context.Context, organizationID uuid.UUID, code string) (*FiscalPosition, error)
	ListFiscalPositions(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*FiscalPosition, error)
	UpdateFiscalPosition(ctx context.Context, fp *FiscalPosition) error
	DeleteFiscalPosition(ctx context.Context, id, organizationID uuid.UUID) error
	GetApplicableFiscalPositions(ctx context.Context, organizationID uuid.UUID, countryCode *string, stateProvince *string) ([]*FiscalPosition, error)

	// Fiscal Position Tax Mappings
	CreateFiscalPositionTaxMapping(ctx context.Context, ftm *FiscalPositionTaxMapping) error
	GetFiscalPositionTaxMapping(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) (*FiscalPositionTaxMapping, error)
	ListFiscalPositionTaxMappings(ctx context.Context, fiscalPositionID uuid.UUID) ([]*FiscalPositionTaxMapping, error)
	UpdateFiscalPositionTaxMapping(ctx context.Context, ftm *FiscalPositionTaxMapping) error
	DeleteFiscalPositionTaxMapping(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) error
	GetMappedTax(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) (*uuid.UUID, error)

	// Tax Report Definitions
	CreateTaxReportDefinition(ctx context.Context, trd *TaxReportDefinition) error
	GetTaxReportDefinition(ctx context.Context, id uuid.UUID) (*TaxReportDefinition, error)
	GetTaxReportDefinitionByCode(ctx context.Context, reportCode string, organizationID *uuid.UUID) (*TaxReportDefinition, error)
	ListTaxReportDefinitions(ctx context.Context, organizationID *uuid.UUID, filter map[string]interface{}) ([]*TaxReportDefinition, error)
	UpdateTaxReportDefinition(ctx context.Context, trd *TaxReportDefinition) error
	DeleteTaxReportDefinition(ctx context.Context, id uuid.UUID) error

	// Tax Report Lines
	CreateTaxReportLine(ctx context.Context, trl *TaxReportLine) error
	GetTaxReportLine(ctx context.Context, id uuid.UUID) (*TaxReportLine, error)
	ListTaxReportLines(ctx context.Context, reportDefinitionID uuid.UUID) ([]*TaxReportLine, error)
	UpdateTaxReportLine(ctx context.Context, trl *TaxReportLine) error
	DeleteTaxReportLine(ctx context.Context, id uuid.UUID) error
}

// Service handles tax business logic
type Service struct {
	repo TaxRepository
}

// NewService creates a new tax service
func NewService(repo TaxRepository) *Service {
	return &Service{repo: repo}
}

// ========================
// TAX GROUP OPERATIONS
// ========================

// CreateTaxGroup creates a new tax group
func (s *Service) CreateTaxGroup(ctx context.Context, orgID uuid.UUID, req *CreateTaxGroupRequest, createdBy uuid.UUID) (*TaxGroup, error) {
	if req.GroupCode == "" || req.GroupName == "" {
		return nil, errors.New("group code and name are required")
	}

	tg := &TaxGroup{
		ID:             uuid.New(),
		OrganizationID: orgID,
		GroupCode:      req.GroupCode,
		GroupName:      req.GroupName,
		Sequence:       req.Sequence,
		IsActive:       req.IsActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      &createdBy,
		UpdatedBy:      &createdBy,
	}

	if err := s.repo.CreateTaxGroup(ctx, tg); err != nil {
		return nil, err
	}

	return tg, nil
}

// GetTaxGroup retrieves a tax group
func (s *Service) GetTaxGroup(ctx context.Context, id, orgID uuid.UUID) (*TaxGroup, error) {
	return s.repo.GetTaxGroup(ctx, id, orgID)
}

// UpdateTaxGroup updates a tax group
func (s *Service) UpdateTaxGroup(ctx context.Context, id, orgID uuid.UUID, req *UpdateTaxGroupRequest, updatedBy uuid.UUID) (*TaxGroup, error) {
	tg, err := s.repo.GetTaxGroup(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	if req.GroupCode != nil {
		tg.GroupCode = *req.GroupCode
	}
	if req.GroupName != nil {
		tg.GroupName = *req.GroupName
	}
	if req.Sequence != nil {
		tg.Sequence = *req.Sequence
	}
	if req.IsActive != nil {
		tg.IsActive = *req.IsActive
	}

	tg.UpdatedAt = time.Now()
	tg.UpdatedBy = &updatedBy

	if err := s.repo.UpdateTaxGroup(ctx, tg); err != nil {
		return nil, err
	}

	return tg, nil
}

// DeleteTaxGroup deletes a tax group
func (s *Service) DeleteTaxGroup(ctx context.Context, id, orgID uuid.UUID) error {
	return s.repo.DeleteTaxGroup(ctx, id, orgID)
}

// ListTaxGroups lists tax groups
func (s *Service) ListTaxGroups(ctx context.Context, orgID uuid.UUID, filter map[string]interface{}) ([]*TaxGroup, error) {
	return s.repo.ListTaxGroups(ctx, orgID, filter)
}

// ========================
// TAX OPERATIONS
// ========================

// CreateTax creates a new tax
func (s *Service) CreateTax(ctx context.Context, orgID uuid.UUID, req *CreateTaxRequest, createdBy uuid.UUID) (*Tax, error) {
	if req.TaxRate < 0 || req.TaxRate > 100 {
		return nil, ErrInvalidTaxRate
	}

	tax := &Tax{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		TaxGroupID:         req.TaxGroupID,
		TaxCode:            req.TaxCode,
		TaxName:            req.TaxName,
		TaxRate:            req.TaxRate,
		TaxScope:           req.TaxScope,
		IsPriceInclusive:   req.IsPriceInclusive,
		TaxAccountID:       req.TaxAccountID,
		TaxRefundAccountID: req.TaxRefundAccountID,
		IsActive:           req.IsActive,
		Description:        req.Description,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          &createdBy,
		UpdatedBy:          &createdBy,
	}

	if err := s.repo.CreateTax(ctx, tax); err != nil {
		return nil, err
	}

	return tax, nil
}

// GetTax retrieves a tax
func (s *Service) GetTax(ctx context.Context, id, orgID uuid.UUID) (*Tax, error) {
	return s.repo.GetTax(ctx, id, orgID)
}

// UpdateTax updates a tax
func (s *Service) UpdateTax(ctx context.Context, id, orgID uuid.UUID, req *UpdateTaxRequest, updatedBy uuid.UUID) (*Tax, error) {
	tax, err := s.repo.GetTax(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	if req.TaxCode != nil {
		tax.TaxCode = *req.TaxCode
	}
	if req.TaxName != nil {
		tax.TaxName = *req.TaxName
	}
	if req.TaxRate != nil {
		if *req.TaxRate < 0 || *req.TaxRate > 100 {
			return nil, ErrInvalidTaxRate
		}
		tax.TaxRate = *req.TaxRate
	}
	if req.TaxScope != nil {
		tax.TaxScope = *req.TaxScope
	}
	if req.IsPriceInclusive != nil {
		tax.IsPriceInclusive = *req.IsPriceInclusive
	}
	if req.TaxAccountID != nil {
		tax.TaxAccountID = *req.TaxAccountID
	}
	if req.TaxRefundAccountID != nil {
		tax.TaxRefundAccountID = req.TaxRefundAccountID
	}
	if req.TaxGroupID != nil {
		tax.TaxGroupID = req.TaxGroupID
	}
	if req.IsActive != nil {
		tax.IsActive = *req.IsActive
	}
	if req.Description != nil {
		tax.Description = req.Description
	}

	tax.UpdatedAt = time.Now()
	tax.UpdatedBy = &updatedBy

	if err := s.repo.UpdateTax(ctx, tax); err != nil {
		return nil, err
	}

	return tax, nil
}

// DeleteTax deletes a tax
func (s *Service) DeleteTax(ctx context.Context, id, orgID uuid.UUID) error {
	return s.repo.DeleteTax(ctx, id, orgID)
}

// ListTaxes lists taxes
func (s *Service) ListTaxes(ctx context.Context, orgID uuid.UUID, filter map[string]interface{}) ([]*Tax, error) {
	return s.repo.ListTaxes(ctx, orgID, filter)
}

// ========================
// FISCAL POSITION OPERATIONS
// ========================

// CreateFiscalPosition creates a new fiscal position
func (s *Service) CreateFiscalPosition(ctx context.Context, orgID uuid.UUID, req *CreateFiscalPositionRequest, createdBy uuid.UUID) (*FiscalPosition, error) {
	if req.PositionCode == "" || req.PositionName == "" {
		return nil, errors.New("position code and name are required")
	}

	fp := &FiscalPosition{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		PositionCode:       req.PositionCode,
		PositionName:       req.PositionName,
		AutoApply:          req.AutoApply,
		CountryID:          req.CountryID,
		StateProvince:      req.StateProvince,
		ZipPostalCodeRange: req.ZipPostalCodeRange,
		IsActive:           req.IsActive,
		Notes:              req.Notes,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          &createdBy,
		UpdatedBy:          &createdBy,
	}

	if err := s.repo.CreateFiscalPosition(ctx, fp); err != nil {
		return nil, err
	}

	return fp, nil
}

// GetFiscalPosition retrieves a fiscal position
func (s *Service) GetFiscalPosition(ctx context.Context, id, orgID uuid.UUID) (*FiscalPosition, error) {
	return s.repo.GetFiscalPosition(ctx, id, orgID)
}

// UpdateFiscalPosition updates a fiscal position
func (s *Service) UpdateFiscalPosition(ctx context.Context, id, orgID uuid.UUID, req *UpdateFiscalPositionRequest, updatedBy uuid.UUID) (*FiscalPosition, error) {
	fp, err := s.repo.GetFiscalPosition(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	if req.PositionCode != nil {
		fp.PositionCode = *req.PositionCode
	}
	if req.PositionName != nil {
		fp.PositionName = *req.PositionName
	}
	if req.AutoApply != nil {
		fp.AutoApply = *req.AutoApply
	}
	if req.CountryID != nil {
		fp.CountryID = req.CountryID
	}
	if req.StateProvince != nil {
		fp.StateProvince = req.StateProvince
	}
	if req.ZipPostalCodeRange != nil {
		fp.ZipPostalCodeRange = req.ZipPostalCodeRange
	}
	if req.IsActive != nil {
		fp.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		fp.Notes = req.Notes
	}

	fp.UpdatedAt = time.Now()
	fp.UpdatedBy = &updatedBy

	if err := s.repo.UpdateFiscalPosition(ctx, fp); err != nil {
		return nil, err
	}

	return fp, nil
}

// DeleteFiscalPosition deletes a fiscal position
func (s *Service) DeleteFiscalPosition(ctx context.Context, id, orgID uuid.UUID) error {
	return s.repo.DeleteFiscalPosition(ctx, id, orgID)
}

// ListFiscalPositions lists fiscal positions
func (s *Service) ListFiscalPositions(ctx context.Context, orgID uuid.UUID, filter map[string]interface{}) ([]*FiscalPosition, error) {
	return s.repo.ListFiscalPositions(ctx, orgID, filter)
}

// GetApplicableFiscalPositions gets fiscal positions for a location
func (s *Service) GetApplicableFiscalPositions(ctx context.Context, orgID uuid.UUID, countryCode *string, stateProvince *string) ([]*FiscalPosition, error) {
	return s.repo.GetApplicableFiscalPositions(ctx, orgID, countryCode, stateProvince)
}

// ========================
// FISCAL POSITION TAX MAPPING OPERATIONS
// ========================

// CreateFiscalPositionTaxMapping creates a tax mapping
func (s *Service) CreateFiscalPositionTaxMapping(ctx context.Context, fiscalPositionID uuid.UUID, req *CreateFiscalPositionTaxMappingRequest, createdBy uuid.UUID) (*FiscalPositionTaxMapping, error) {
	ftm := &FiscalPositionTaxMapping{
		ID:               uuid.New(),
		FiscalPositionID: fiscalPositionID,
		SourceTaxID:      req.SourceTaxID,
		DestinationTaxID: req.DestinationTaxID,
		CreatedAt:        time.Now(),
		CreatedBy:        &createdBy,
	}

	if err := s.repo.CreateFiscalPositionTaxMapping(ctx, ftm); err != nil {
		return nil, err
	}

	return ftm, nil
}

// GetMappedTax gets the mapped tax for a fiscal position
func (s *Service) GetMappedTax(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) (*uuid.UUID, error) {
	return s.repo.GetMappedTax(ctx, fiscalPositionID, sourceTaxID)
}

// ListFiscalPositionTaxMappings lists tax mappings
func (s *Service) ListFiscalPositionTaxMappings(ctx context.Context, fiscalPositionID uuid.UUID) ([]*FiscalPositionTaxMapping, error) {
	return s.repo.ListFiscalPositionTaxMappings(ctx, fiscalPositionID)
}

// UpdateFiscalPositionTaxMapping updates a tax mapping
func (s *Service) UpdateFiscalPositionTaxMapping(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID, req *UpdateFiscalPositionTaxMappingRequest) (*FiscalPositionTaxMapping, error) {
	ftm, err := s.repo.GetFiscalPositionTaxMapping(ctx, fiscalPositionID, sourceTaxID)
	if err != nil {
		return nil, err
	}

	if req.DestinationTaxID != nil {
		ftm.DestinationTaxID = req.DestinationTaxID
	}

	if err := s.repo.UpdateFiscalPositionTaxMapping(ctx, ftm); err != nil {
		return nil, err
	}

	return ftm, nil
}

// DeleteFiscalPositionTaxMapping deletes a tax mapping
func (s *Service) DeleteFiscalPositionTaxMapping(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) error {
	return s.repo.DeleteFiscalPositionTaxMapping(ctx, fiscalPositionID, sourceTaxID)
}

// ========================
// TAX REPORT DEFINITION OPERATIONS
// ========================

// CreateTaxReportDefinition creates a new tax report definition
func (s *Service) CreateTaxReportDefinition(ctx context.Context, orgID *uuid.UUID, req *CreateTaxReportDefinitionRequest, createdBy uuid.UUID) (*TaxReportDefinition, error) {
	if req.ReportCode == "" || req.ReportName == "" {
		return nil, errors.New("report code and name are required")
	}

	if req.EffectiveFrom != nil && req.EffectiveTo != nil {
		if req.EffectiveTo.Before(*req.EffectiveFrom) {
			return nil, ErrInvalidDateRange
		}
	}

	trd := &TaxReportDefinition{
		ID:                    uuid.New(),
		OrganizationID:        orgID,
		LocalizationPackageID: req.LocalizationPackageID,
		ReportCode:            req.ReportCode,
		ReportName:            req.ReportName,
		Jurisdiction:          req.Jurisdiction,
		Authority:             req.Authority,
		ReportFrequency:       req.ReportFrequency,
		Version:               req.Version,
		EffectiveFrom:         req.EffectiveFrom,
		EffectiveTo:           req.EffectiveTo,
		IsActive:              req.IsActive,
		Description:           req.Description,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		CreatedBy:             &createdBy,
	}

	if err := s.repo.CreateTaxReportDefinition(ctx, trd); err != nil {
		return nil, err
	}

	return trd, nil
}

// GetTaxReportDefinition retrieves a tax report definition
func (s *Service) GetTaxReportDefinition(ctx context.Context, id uuid.UUID) (*TaxReportDefinition, error) {
	return s.repo.GetTaxReportDefinition(ctx, id)
}

// UpdateTaxReportDefinition updates a tax report definition
func (s *Service) UpdateTaxReportDefinition(ctx context.Context, id uuid.UUID, req *UpdateTaxReportDefinitionRequest, updatedBy uuid.UUID) (*TaxReportDefinition, error) {
	trd, err := s.repo.GetTaxReportDefinition(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ReportCode != nil {
		trd.ReportCode = *req.ReportCode
	}
	if req.ReportName != nil {
		trd.ReportName = *req.ReportName
	}
	if req.Jurisdiction != nil {
		trd.Jurisdiction = req.Jurisdiction
	}
	if req.Authority != nil {
		trd.Authority = req.Authority
	}
	if req.ReportFrequency != nil {
		trd.ReportFrequency = req.ReportFrequency
	}
	if req.Version != nil {
		trd.Version = req.Version
	}
	if req.EffectiveFrom != nil {
		trd.EffectiveFrom = req.EffectiveFrom
	}
	if req.EffectiveTo != nil {
		trd.EffectiveTo = req.EffectiveTo
	}
	if req.IsActive != nil {
		trd.IsActive = *req.IsActive
	}
	if req.Description != nil {
		trd.Description = req.Description
	}
	if req.LocalizationPackageID != nil {
		trd.LocalizationPackageID = req.LocalizationPackageID
	}

	trd.UpdatedAt = time.Now()

	if err := s.repo.UpdateTaxReportDefinition(ctx, trd); err != nil {
		return nil, err
	}

	return trd, nil
}

// DeleteTaxReportDefinition deletes a tax report definition
func (s *Service) DeleteTaxReportDefinition(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTaxReportDefinition(ctx, id)
}

// ListTaxReportDefinitions lists tax report definitions
func (s *Service) ListTaxReportDefinitions(ctx context.Context, orgID *uuid.UUID, filter map[string]interface{}) ([]*TaxReportDefinition, error) {
	return s.repo.ListTaxReportDefinitions(ctx, orgID, filter)
}

// ========================
// TAX REPORT LINE OPERATIONS
// ========================

// CreateTaxReportLine creates a new tax report line
func (s *Service) CreateTaxReportLine(ctx context.Context, reportDefinitionID uuid.UUID, req *CreateTaxReportLineRequest) (*TaxReportLine, error) {
	if req.LineCode == "" || req.LineName == "" {
		return nil, errors.New("line code and name are required")
	}

	trl := &TaxReportLine{
		ID:                    uuid.New(),
		TaxReportDefinitionID: reportDefinitionID,
		LineCode:              req.LineCode,
		LineName:              req.LineName,
		Sequence:              req.Sequence,
		ParentLineID:          req.ParentLineID,
		FormulaType:           req.FormulaType,
		Formula:               req.Formula,
		TaxGroupIDs:           req.TaxGroupIDs,
		AccountIDs:            req.AccountIDs,
		TaxIDs:                req.TaxIDs,
		IsSubtotal:            req.IsSubtotal,
		IsTotal:               req.IsTotal,
		Notes:                 req.Notes,
		CreatedAt:             time.Now(),
	}

	if err := s.repo.CreateTaxReportLine(ctx, trl); err != nil {
		return nil, err
	}

	return trl, nil
}

// GetTaxReportLine retrieves a tax report line
func (s *Service) GetTaxReportLine(ctx context.Context, id uuid.UUID) (*TaxReportLine, error) {
	return s.repo.GetTaxReportLine(ctx, id)
}

// UpdateTaxReportLine updates a tax report line
func (s *Service) UpdateTaxReportLine(ctx context.Context, id uuid.UUID, req *UpdateTaxReportLineRequest) (*TaxReportLine, error) {
	trl, err := s.repo.GetTaxReportLine(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LineCode != nil {
		trl.LineCode = *req.LineCode
	}
	if req.LineName != nil {
		trl.LineName = *req.LineName
	}
	if req.Sequence != nil {
		trl.Sequence = *req.Sequence
	}
	if req.ParentLineID != nil {
		trl.ParentLineID = req.ParentLineID
	}
	if req.FormulaType != nil {
		trl.FormulaType = req.FormulaType
	}
	if req.Formula != nil {
		trl.Formula = req.Formula
	}
	if len(req.TaxGroupIDs) > 0 {
		trl.TaxGroupIDs = req.TaxGroupIDs
	}
	if len(req.AccountIDs) > 0 {
		trl.AccountIDs = req.AccountIDs
	}
	if len(req.TaxIDs) > 0 {
		trl.TaxIDs = req.TaxIDs
	}
	if req.IsSubtotal != nil {
		trl.IsSubtotal = *req.IsSubtotal
	}
	if req.IsTotal != nil {
		trl.IsTotal = *req.IsTotal
	}
	if req.Notes != nil {
		trl.Notes = req.Notes
	}

	if err := s.repo.UpdateTaxReportLine(ctx, trl); err != nil {
		return nil, err
	}

	return trl, nil
}

// DeleteTaxReportLine deletes a tax report line
func (s *Service) DeleteTaxReportLine(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTaxReportLine(ctx, id)
}

// ListTaxReportLines lists tax report lines
func (s *Service) ListTaxReportLines(ctx context.Context, reportDefinitionID uuid.UUID) ([]*TaxReportLine, error) {
	return s.repo.ListTaxReportLines(ctx, reportDefinitionID)
}

// ========================
// TAX CALCULATION
// ========================

// CalculateTax calculates tax amount for a given base amount
func (s *Service) CalculateTax(ctx context.Context, taxID uuid.UUID, baseAmount float64, orgID uuid.UUID) (*TaxCalculationResult, error) {
	tax, err := s.repo.GetTax(ctx, taxID, orgID)
	if err != nil {
		return nil, err
	}

	var taxAmount float64
	if tax.IsPriceInclusive {
		// Tax is already included in base amount
		taxAmount = (baseAmount / (100 + tax.TaxRate)) * tax.TaxRate
	} else {
		// Tax is added on top of base amount
		taxAmount = (baseAmount / 100) * tax.TaxRate
	}

	return &TaxCalculationResult{
		TaxID:              tax.ID,
		TaxCode:            tax.TaxCode,
		TaxName:            tax.TaxName,
		TaxRate:            tax.TaxRate,
		BaseAmount:         baseAmount,
		TaxAmount:          taxAmount,
		IsPriceInclusive:   tax.IsPriceInclusive,
		TaxAccountID:       tax.TaxAccountID,
		TaxRefundAccountID: tax.TaxRefundAccountID,
	}, nil
}

// ApplyFiscalPosition applies tax mapping based on fiscal position
func (s *Service) ApplyFiscalPosition(ctx context.Context, sourceTaxID uuid.UUID, fiscalPositionID uuid.UUID) (*uuid.UUID, error) {
	return s.repo.GetMappedTax(ctx, fiscalPositionID, sourceTaxID)
}

// CalculateTaxes calculates multiple taxes for a base amount
func (s *Service) CalculateTaxes(ctx context.Context, baseAmount float64, taxIDs []uuid.UUID, orgID uuid.UUID) ([]*TaxCalculationResult, error) {
	if len(taxIDs) == 0 {
		return []*TaxCalculationResult{}, nil
	}

	var results []*TaxCalculationResult
	for _, taxID := range taxIDs {
		result, err := s.CalculateTax(ctx, taxID, baseAmount, orgID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax %s: %w", taxID, err)
		}
		results = append(results, result)
	}

	return results, nil
}

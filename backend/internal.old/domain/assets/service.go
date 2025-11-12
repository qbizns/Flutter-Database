package assets

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// Repository defines data access operations for assets
type Repository interface {
	// Asset Categories
	CreateAssetCategory(ctx context.Context, category *AssetCategory) error
	GetAssetCategory(ctx context.Context, id uuid.UUID) (*AssetCategory, error)
	GetAssetCategoryByCode(ctx context.Context, code string) (*AssetCategory, error)
	ListAssetCategories(ctx context.Context, organizationID *uuid.UUID) ([]*AssetCategory, error)
	UpdateAssetCategory(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteAssetCategory(ctx context.Context, id uuid.UUID) error

	// Fixed Assets
	CreateFixedAsset(ctx context.Context, asset *FixedAsset) error
	GetFixedAsset(ctx context.Context, id uuid.UUID) (*FixedAsset, error)
	GetFixedAssetByNumber(ctx context.Context, organizationID uuid.UUID, assetNumber string) (*FixedAsset, error)
	ListFixedAssets(ctx context.Context, query AssetQuery) ([]*FixedAsset, int64, error)
	UpdateFixedAsset(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteFixedAsset(ctx context.Context, id uuid.UUID) error

	// Depreciation Schedules
	CreateDepreciationSchedule(ctx context.Context, schedule *AssetDepreciationSchedule) error
	GetDepreciationSchedule(ctx context.Context, id uuid.UUID) (*AssetDepreciationSchedule, error)
	ListDepreciationSchedules(ctx context.Context, query DepreciationQuery) ([]*AssetDepreciationSchedule, int64, error)
	UpdateDepreciationSchedule(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteDepreciationSchedule(ctx context.Context, id uuid.UUID) error
	GetDepreciationScheduleByAssetAndPeriod(ctx context.Context, assetID, periodID uuid.UUID) (*AssetDepreciationSchedule, error)

	// Statistics
	GetAssetStatistics(ctx context.Context, organizationID uuid.UUID) (*AssetStatistics, error)
	GetAssetByCategory(ctx context.Context, organizationID uuid.UUID, categoryID uuid.UUID) ([]*FixedAsset, error)
}

// Service provides business logic for asset management
type Service struct {
	repo Repository
}

// NewService creates a new asset service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateAssetCategory creates a new asset category
func (s *Service) CreateAssetCategory(ctx context.Context, req *AssetCategory) (*AssetCategory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	req.ID = uuid.New()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := s.repo.CreateAssetCategory(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to create asset category: %w", err)
	}

	return req, nil
}

// GetAssetCategory retrieves an asset category by ID
func (s *Service) GetAssetCategory(ctx context.Context, id uuid.UUID) (*AssetCategory, error) {
	category, err := s.repo.GetAssetCategory(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset category: %w", err)
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

// ListAssetCategories retrieves all asset categories
func (s *Service) ListAssetCategories(ctx context.Context, organizationID *uuid.UUID) ([]*AssetCategory, error) {
	categories, err := s.repo.ListAssetCategories(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list asset categories: %w", err)
	}
	return categories, nil
}

// UpdateAssetCategory updates an asset category
func (s *Service) UpdateAssetCategory(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	category, err := s.GetAssetCategory(ctx, id)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateAssetCategory(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update asset category: %w", err)
	}

	return nil
}

// DeleteAssetCategory deletes an asset category
func (s *Service) DeleteAssetCategory(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetAssetCategory(ctx, id); err != nil {
		return err
	}

	if err := s.repo.DeleteAssetCategory(ctx, id); err != nil {
		return fmt.Errorf("failed to delete asset category: %w", err)
	}

	return nil
}

// CreateFixedAsset creates a new fixed asset
func (s *Service) CreateFixedAsset(ctx context.Context, req *CreateFixedAssetRequest, organizationID uuid.UUID) (*FixedAsset, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	// Validate request
	if req.AssetNumber == "" {
		return nil, ErrAssetNumberRequired
	}
	if req.AssetName == "" {
		return nil, ErrAssetNameRequired
	}
	if req.AcquisitionCost < 0 {
		return nil, ErrInvalidAcquisitionCost
	}
	if req.SalvageValue > req.AcquisitionCost {
		return nil, ErrSalvageValueExceedsAcquisitionCost
	}
	if req.UsefulLifeYears <= 0 {
		return nil, ErrInvalidUsefulLife
	}
	if !IsValidDepreciationMethod(req.DepreciationMethod) {
		return nil, ErrInvalidDepreciationMethod
	}
	if req.DepreciationStartDate.Before(req.AcquisitionDate) {
		return nil, ErrInvalidDepreciationDate
	}

	// Check for duplicate asset number
	existing, _ := s.repo.GetFixedAssetByNumber(ctx, organizationID, req.AssetNumber)
	if existing != nil {
		return nil, ErrDuplicateAssetNumber
	}

	asset := &FixedAsset{
		ID:                                uuid.New(),
		OrganizationID:                    organizationID,
		AssetNumber:                       req.AssetNumber,
		AssetName:                         req.AssetName,
		AssetCategoryID:                   req.AssetCategoryID,
		AcquisitionDate:                   req.AcquisitionDate,
		AcquisitionCost:                   req.AcquisitionCost,
		SalvageValue:                      req.SalvageValue,
		SupplierID:                        req.SupplierID,
		DepreciationMethod:                req.DepreciationMethod,
		UsefulLifeYears:                   req.UsefulLifeYears,
		DepreciationStartDate:             req.DepreciationStartDate,
		AssetAccountID:                    req.AssetAccountID,
		AccumulatedDepreciationAccountID:  req.AccumulatedDepreciationAccountID,
		DepreciationExpenseAccountID:      req.DepreciationExpenseAccountID,
		CurrentBookValue:                  req.AcquisitionCost,
		AccumulatedDepreciation:           0,
		LocationID:                        req.LocationID,
		Department:                        req.Department,
		IsDisposed:                        false,
		Description:                       req.Description,
		SerialNumber:                      req.SerialNumber,
		Notes:                             req.Notes,
		Metadata:                          req.Metadata,
		CreatedAt:                         time.Now(),
		UpdatedAt:                         time.Now(),
	}

	if err := asset.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateFixedAsset(ctx, asset); err != nil {
		return nil, fmt.Errorf("failed to create fixed asset: %w", err)
	}

	return asset, nil
}

// GetFixedAsset retrieves a fixed asset by ID
func (s *Service) GetFixedAsset(ctx context.Context, id uuid.UUID) (*FixedAsset, error) {
	asset, err := s.repo.GetFixedAsset(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed asset: %w", err)
	}
	if asset == nil {
		return nil, ErrAssetNotFound
	}
	return asset, nil
}

// ListFixedAssets retrieves fixed assets based on query filters
func (s *Service) ListFixedAssets(ctx context.Context, query AssetQuery) ([]*FixedAsset, int64, error) {
	if query.OrganizationID == uuid.Nil {
		return nil, 0, ErrOrganizationIDRequired
	}

	assets, total, err := s.repo.ListFixedAssets(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list fixed assets: %w", err)
	}

	return assets, total, nil
}

// UpdateFixedAsset updates a fixed asset
func (s *Service) UpdateFixedAsset(ctx context.Context, id uuid.UUID, req *UpdateFixedAssetRequest) error {
	asset, err := s.GetFixedAsset(ctx, id)
	if err != nil {
		return err
	}

	if asset.IsDisposed {
		return ErrAssetAlreadyDisposed
	}

	updates := make(map[string]interface{})
	if req.AssetName != nil {
		updates["asset_name"] = *req.AssetName
	}
	if req.AssetCategoryID != nil {
		updates["asset_category_id"] = *req.AssetCategoryID
	}
	if req.SalvageValue != nil {
		if *req.SalvageValue > asset.AcquisitionCost {
			return ErrSalvageValueExceedsAcquisitionCost
		}
		updates["salvage_value"] = *req.SalvageValue
	}
	if req.UsefulLifeYears != nil {
		if *req.UsefulLifeYears <= 0 {
			return ErrInvalidUsefulLife
		}
		updates["useful_life_years"] = *req.UsefulLifeYears
	}
	if req.LocationID != nil {
		updates["location_id"] = *req.LocationID
	}
	if req.Department != nil {
		updates["department"] = *req.Department
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateFixedAsset(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update fixed asset: %w", err)
	}

	return nil
}

// DisposeAsset marks an asset as disposed
func (s *Service) DisposeAsset(ctx context.Context, id uuid.UUID, req *DisposeAssetRequest) error {
	asset, err := s.GetFixedAsset(ctx, id)
	if err != nil {
		return err
	}

	if asset.IsDisposed {
		return ErrAssetAlreadyDisposed
	}

	updates := map[string]interface{}{
		"is_disposed":      true,
		"disposal_date":    req.DisposalDate,
		"disposal_proceeds": req.DisposalProceeds,
		"updated_at":       time.Now(),
	}

	if err := s.repo.UpdateFixedAsset(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to dispose asset: %w", err)
	}

	return nil
}

// CalculateDepreciation calculates depreciation based on method
func (s *Service) CalculateDepreciation(input *CalculateDepreciationInput) (*DepreciationResult, error) {
	if !IsValidDepreciationMethod(input.DepreciationMethod) {
		return nil, ErrInvalidDepreciationMethod
	}

	depreciableBase := input.AcquisitionCost - input.SalvageValue
	if depreciableBase < 0 {
		depreciableBase = 0
	}

	switch input.DepreciationMethod {
	case DepreciationMethodStraightLine:
		return s.calculateStraightLineDepreciation(input, depreciableBase)
	case DepreciationMethodDecliningBalance:
		return s.calculateDecliningBalanceDepreciation(input, depreciableBase)
	case DepreciationMethodUnitsOfProduction:
		// For UOP, we'd need production units data
		// Placeholder implementation using straight-line
		return s.calculateStraightLineDepreciation(input, depreciableBase)
	default:
		return nil, ErrInvalidDepreciationMethod
	}
}

func (s *Service) calculateStraightLineDepreciation(input *CalculateDepreciationInput, depreciableBase float64) (*DepreciationResult, error) {
	annualDepreciation := depreciableBase / float64(input.UsefulLifeYears)
	monthsElapsed := s.calculateMonthsBetween(input.DepreciationStart, input.TargetDate)
	accumulatedDepreciation := (annualDepreciation / 12) * float64(monthsElapsed)

	if accumulatedDepreciation > depreciableBase {
		accumulatedDepreciation = depreciableBase
	}

	bookValue := input.AcquisitionCost - accumulatedDepreciation
	if bookValue < input.SalvageValue {
		bookValue = input.SalvageValue
	}

	monthsDepreciated := monthsElapsed
	remainingMonths := (input.UsefulLifeYears * 12) - monthsDepreciated
	remainingYears := int(math.Ceil(float64(remainingMonths) / 12.0))

	return &DepreciationResult{
		DepreciationAmount:            annualDepreciation,
		AccumulatedDepreciationToDate: accumulatedDepreciation,
		BookValueAtDate:               bookValue,
		RemainingUsefulLife:           remainingYears,
		HasFullyDepreciated:           accumulatedDepreciation >= depreciableBase,
	}, nil
}

func (s *Service) calculateDecliningBalanceDepreciation(input *CalculateDepreciationInput, depreciableBase float64) (*DepreciationResult, error) {
	straightLineRate := 1.0 / float64(input.UsefulLifeYears)
	decliningRate := straightLineRate * 2.0 // Double declining balance

	accumulatedDepreciation := 0.0
	bookValue := input.AcquisitionCost

	monthsElapsed := s.calculateMonthsBetween(input.DepreciationStart, input.TargetDate)
	yearsElapsed := float64(monthsElapsed) / 12.0

	for year := 0.0; year < yearsElapsed; year++ {
		yearDepreciation := bookValue * decliningRate
		if yearDepreciation > depreciableBase-accumulatedDepreciation {
			yearDepreciation = depreciableBase - accumulatedDepreciation
		}
		accumulatedDepreciation += yearDepreciation
		bookValue -= yearDepreciation
	}

	if bookValue < input.SalvageValue {
		bookValue = input.SalvageValue
	}

	yearlyDepreciation := bookValue * decliningRate
	if accumulatedDepreciation+yearlyDepreciation > depreciableBase {
		yearlyDepreciation = depreciableBase - accumulatedDepreciation
	}

	remainingYears := int(math.Ceil((1.0 / decliningRate) - yearsElapsed))

	return &DepreciationResult{
		DepreciationAmount:            yearlyDepreciation,
		AccumulatedDepreciationToDate: accumulatedDepreciation,
		BookValueAtDate:               bookValue,
		RemainingUsefulLife:           remainingYears,
		HasFullyDepreciated:           accumulatedDepreciation >= depreciableBase,
	}, nil
}

// CalculateMonthsBetween calculates the number of months between two dates
func (s *Service) calculateMonthsBetween(start, end time.Time) int {
	months := (end.Year() - start.Year()) * 12
	months += int(end.Month() - start.Month())
	return months
}

// CreateDepreciationSchedule creates a depreciation schedule entry
func (s *Service) CreateDepreciationSchedule(ctx context.Context, req *CreateDepreciationScheduleRequest, organizationID uuid.UUID) (*AssetDepreciationSchedule, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	asset, err := s.GetFixedAsset(ctx, req.FixedAssetID)
	if err != nil {
		return nil, err
	}

	if asset.IsDisposed {
		return nil, ErrAssetAlreadyDisposed
	}

	depResult, err := s.CalculateDepreciation(&CalculateDepreciationInput{
		FixedAssetID:       req.FixedAssetID,
		TargetDate:         req.DepreciationDate,
		DepreciationMethod: asset.DepreciationMethod,
		AcquisitionCost:    asset.AcquisitionCost,
		SalvageValue:       asset.SalvageValue,
		UsefulLifeYears:    asset.UsefulLifeYears,
		DepreciationStart:  asset.DepreciationStartDate,
		LastDepreciation:   asset.LastDepreciationDate,
	})
	if err != nil {
		return nil, err
	}

	schedule := &AssetDepreciationSchedule{
		ID:                                uuid.New(),
		OrganizationID:                    organizationID,
		FixedAssetID:                      req.FixedAssetID,
		FiscalYearID:                      req.FiscalYearID,
		AccountingPeriodID:                req.AccountingPeriodID,
		DepreciationDate:                  req.DepreciationDate,
		DepreciationAmount:                depResult.DepreciationAmount,
		AccumulatedDepreciationBeginning:  asset.AccumulatedDepreciation,
		AccumulatedDepreciationEnding:     depResult.AccumulatedDepreciationToDate,
		BookValueBeginning:                asset.CurrentBookValue,
		BookValueEnding:                   depResult.BookValueAtDate,
		IsPosted:                          false,
		CreatedAt:                         time.Now(),
	}

	if err := s.repo.CreateDepreciationSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create depreciation schedule: %w", err)
	}

	return schedule, nil
}

// GetDepreciationSchedule retrieves a depreciation schedule by ID
func (s *Service) GetDepreciationSchedule(ctx context.Context, id uuid.UUID) (*AssetDepreciationSchedule, error) {
	schedule, err := s.repo.GetDepreciationSchedule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get depreciation schedule: %w", err)
	}
	if schedule == nil {
		return nil, ErrDepreciationScheduleNotFound
	}
	return schedule, nil
}

// ListDepreciationSchedules retrieves depreciation schedules based on query
func (s *Service) ListDepreciationSchedules(ctx context.Context, query DepreciationQuery) ([]*AssetDepreciationSchedule, int64, error) {
	if query.OrganizationID == uuid.Nil {
		return nil, 0, ErrOrganizationIDRequired
	}

	schedules, total, err := s.repo.ListDepreciationSchedules(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list depreciation schedules: %w", err)
	}

	return schedules, total, nil
}

// UpdateDepreciationSchedule updates a depreciation schedule
func (s *Service) UpdateDepreciationSchedule(ctx context.Context, id uuid.UUID, req *UpdateDepreciationScheduleRequest) error {
	schedule, err := s.GetDepreciationSchedule(ctx, id)
	if err != nil {
		return err
	}

	if schedule.IsPosted {
		return fmt.Errorf("cannot update posted depreciation schedule")
	}

	updates := make(map[string]interface{})
	if req.DepreciationAmount != nil {
		updates["depreciation_amount"] = *req.DepreciationAmount
		updates["accumulated_depreciation_ending"] = schedule.AccumulatedDepreciationBeginning + *req.DepreciationAmount
		updates["book_value_ending"] = schedule.BookValueBeginning - *req.DepreciationAmount
	}

	if err := s.repo.UpdateDepreciationSchedule(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update depreciation schedule: %w", err)
	}

	return nil
}

// PostDepreciationSchedule marks schedule as posted to GL
func (s *Service) PostDepreciationSchedule(ctx context.Context, id uuid.UUID, journalEntryID uuid.UUID) error {
	schedule, err := s.GetDepreciationSchedule(ctx, id)
	if err != nil {
		return err
	}

	if schedule.IsPosted {
		return fmt.Errorf("depreciation schedule already posted")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"is_posted":       true,
		"journal_entry_id": journalEntryID,
		"posted_at":       now,
	}

	if err := s.repo.UpdateDepreciationSchedule(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to post depreciation schedule: %w", err)
	}

	return nil
}

// GetAssetStatistics retrieves aggregated asset statistics
func (s *Service) GetAssetStatistics(ctx context.Context, organizationID uuid.UUID) (*AssetStatistics, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	stats, err := s.repo.GetAssetStatistics(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset statistics: %w", err)
	}

	return stats, nil
}

// GetAssetsByCategory retrieves assets within a category
func (s *Service) GetAssetsByCategory(ctx context.Context, organizationID uuid.UUID, categoryID uuid.UUID) ([]*FixedAsset, error) {
	assets, err := s.repo.GetAssetByCategory(ctx, organizationID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets by category: %w", err)
	}

	return assets, nil
}

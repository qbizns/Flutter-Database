package assets

import "errors"

var (
	ErrCategoryCodeRequired              = errors.New("category code is required")
	ErrCategoryNameRequired              = errors.New("category name is required")
	ErrAssetNumberRequired               = errors.New("asset number is required")
	ErrAssetNameRequired                 = errors.New("asset name is required")
	ErrInvalidAcquisitionCost            = errors.New("acquisition cost must be non-negative")
	ErrSalvageValueExceedsAcquisitionCost = errors.New("salvage value cannot exceed acquisition cost")
	ErrInvalidUsefulLife                 = errors.New("useful life years must be positive")
	ErrInvalidDepreciationMethod         = errors.New("invalid depreciation method")
	ErrAssetNotFound                     = errors.New("asset not found")
	ErrCategoryNotFound                  = errors.New("asset category not found")
	ErrDepreciationScheduleNotFound      = errors.New("depreciation schedule not found")
	ErrAssetAlreadyDisposed              = errors.New("asset has already been disposed")
	ErrAssetNotDisposed                  = errors.New("asset has not been disposed")
	ErrDuplicateAssetNumber              = errors.New("asset number already exists for this organization")
	ErrInvalidDepreciationDate           = errors.New("depreciation date must be on or after acquisition date")
	ErrOrganizationIDRequired            = errors.New("organization_id is required")
)

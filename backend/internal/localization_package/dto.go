package localization_package

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LocalizationPackagesResponse represents a localization_packages response
type LocalizationPackagesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	PackageCode string `json:"package_code"`
	
	PackageName string `json:"package_name"`
	
	CountryCode *string `json:"country_code"`
	
	Region *string `json:"region"`
	
	Description *string `json:"description"`
	
	Version *string `json:"version"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateLocalizationPackagesRequest represents a request to create a localization_packages
type CreateLocalizationPackagesRequest struct {
	
	PackageCode string `json:"package_code" validate:"required"`
	
	PackageName string `json:"package_name" validate:"required"`
	
	CountryCode *string `json:"country_code"`
	
	Region *string `json:"region"`
	
	Description *string `json:"description"`
	
	Version *string `json:"version"`
	
	IsActive *bool `json:"is_active"`
	
}

// Validate validates the create request
func (r *CreateLocalizationPackagesRequest) Validate() error {
	
	if r.PackageCode == "" {
		return fmt.Errorf("package_code is required")
	}
	
	if r.PackageName == "" {
		return fmt.Errorf("package_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLocalizationPackagesRequest represents a request to update a localization_packages
type UpdateLocalizationPackagesRequest struct {
	
	PackageCode *string `json:"package_code,omitempty" validate:"omitempty,required"`
	
	PackageName *string `json:"package_name,omitempty" validate:"omitempty,required"`
	
	CountryCode *string `json:"country_code,omitempty"`
	
	Region *string `json:"region,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Version *string `json:"version,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLocalizationPackagesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PackageCode != nil {
		hasUpdate = true
	}
	
	if r.PackageName != nil {
		hasUpdate = true
	}
	
	if r.CountryCode != nil {
		hasUpdate = true
	}
	
	if r.Region != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Version != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// LocalizationPackagesListResponse represents a paginated list of localization_packages records
type LocalizationPackagesListResponse struct {
	Items      []*LocalizationPackagesResponse `json:"items"`
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

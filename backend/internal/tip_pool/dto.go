package tip_pool

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TipPoolsResponse represents a tip_pools response
type TipPoolsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	PoolName string `json:"pool_name"`
	
	PoolType *string `json:"pool_type"`
	
	Description *string `json:"description"`
	
	DistributionMethod *string `json:"distribution_method"`
	
	DistributionConfig json.RawMessage `json:"distribution_config"`
	
	EligiblePositions *string `json:"eligible_positions"`
	
	IsActive *bool `json:"is_active"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'daily', *string `json:"'daily',"`
	
	'equal', *string `json:"'equal',"`
	
}

// CreateTipPoolsRequest represents a request to create a tip_pools
type CreateTipPoolsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	PoolName string `json:"pool_name" validate:"required"`
	
	PoolType *string `json:"pool_type"`
	
	Description *string `json:"description"`
	
	DistributionMethod *string `json:"distribution_method"`
	
	DistributionConfig json.RawMessage `json:"distribution_config"`
	
	EligiblePositions *string `json:"eligible_positions"`
	
	IsActive *bool `json:"is_active"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'daily', *string `json:"'daily',"`
	
	'equal', *string `json:"'equal',"`
	
}

// Validate validates the create request
func (r *CreateTipPoolsRequest) Validate() error {
	
	if r.PoolName == "" {
		return fmt.Errorf("pool_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTipPoolsRequest represents a request to update a tip_pools
type UpdateTipPoolsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	PoolName *string `json:"pool_name,omitempty" validate:"omitempty,required"`
	
	PoolType *string `json:"pool_type,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	DistributionMethod *string `json:"distribution_method,omitempty"`
	
	DistributionConfig *json.RawMessage `json:"distribution_config,omitempty"`
	
	EligiblePositions *string `json:"eligible_positions,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'daily', *string `json:"'daily',,omitempty"`
	
	'equal', *string `json:"'equal',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTipPoolsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.PoolName != nil {
		hasUpdate = true
	}
	
	if r.PoolType != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.DistributionMethod != nil {
		hasUpdate = true
	}
	
	if r.DistributionConfig != nil {
		hasUpdate = true
	}
	
	if r.EligiblePositions != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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
	
	if r.'daily', != nil {
		hasUpdate = true
	}
	
	if r.'equal', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// TipPoolsListResponse represents a paginated list of tip_pools records
type TipPoolsListResponse struct {
	Items      []*TipPoolsResponse `json:"items"`
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

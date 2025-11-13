package fixed_asset

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for FixedAssets
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new FixedAssets service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new fixed_assets
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateFixedAssetsRequest) (*FixedAssetsResponse, error) {
	s.logger.Info("creating fixed_assets",
		zap.String("organization_id", orgID.String()),
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Convert DTO to entity
	entity := &FixedAssets{
		OrganizationID: orgID,
		
		AssetNumber: req.AssetNumber,
		
		AssetName: req.AssetName,
		
		AssetCategoryId: req.AssetCategoryId,
		
		AcquisitionDate: req.AcquisitionDate,
		
		AcquisitionCost: req.AcquisitionCost,
		
		SalvageValue: req.SalvageValue,
		
		SupplierId: req.SupplierId,
		
		VendorBillId: req.VendorBillId,
		
		DepreciationMethod: req.DepreciationMethod,
		
		UsefulLifeYears: req.UsefulLifeYears,
		
		DepreciationStartDate: req.DepreciationStartDate,
		
		AssetAccountId: req.AssetAccountId,
		
		AccumulatedDepreciationAccountId: req.AccumulatedDepreciationAccountId,
		
		DepreciationExpenseAccountId: req.DepreciationExpenseAccountId,
		
		CurrentBookValue: req.CurrentBookValue,
		
		AccumulatedDepreciation: req.AccumulatedDepreciation,
		
		LastDepreciationDate: req.LastDepreciationDate,
		
		LocationId: req.LocationId,
		
		Department: req.Department,
		
		IsDisposed: req.IsDisposed,
		
		DisposalDate: req.DisposalDate,
		
		DisposalProceeds: req.DisposalProceeds,
		
		DisposalJournalEntryId: req.DisposalJournalEntryId,
		
		Description: req.Description,
		
		SerialNumber: req.SerialNumber,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create fixed_assets: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created fixed_assets",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a fixed_assets by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*FixedAssetsResponse, error) {
	s.logger.Debug("getting fixed_assets",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed_assets: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("fixed_assets not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of fixed_assets records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*FixedAssetsListResponse, error) {
	s.logger.Debug("listing fixed_assets",
		zap.String("organization_id", orgID.String()),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	// Get from database (organization-scoped)
	entities, total, err := s.repo.ListByOrganization(ctx, tx, orgID, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list fixed_assets: %w", err)
	}

	// Convert to response
	items := make([]*FixedAssetsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &FixedAssetsListResponse{
		Items: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing fixed_assets
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateFixedAssetsRequest) (*FixedAssetsResponse, error) {
	s.logger.Info("updating fixed_assets",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed_assets: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("fixed_assets not found or access denied")
	}
	

	// Update fields
	
	if req.AssetNumber != nil {
		entity.AssetNumber = req.AssetNumber
	}
	
	if req.AssetName != nil {
		entity.AssetName = req.AssetName
	}
	
	if req.AssetCategoryId != nil {
		entity.AssetCategoryId = req.AssetCategoryId
	}
	
	if req.AcquisitionDate != nil {
		entity.AcquisitionDate = req.AcquisitionDate
	}
	
	if req.AcquisitionCost != nil {
		entity.AcquisitionCost = req.AcquisitionCost
	}
	
	if req.SalvageValue != nil {
		entity.SalvageValue = req.SalvageValue
	}
	
	if req.SupplierId != nil {
		entity.SupplierId = req.SupplierId
	}
	
	if req.VendorBillId != nil {
		entity.VendorBillId = req.VendorBillId
	}
	
	if req.DepreciationMethod != nil {
		entity.DepreciationMethod = req.DepreciationMethod
	}
	
	if req.UsefulLifeYears != nil {
		entity.UsefulLifeYears = req.UsefulLifeYears
	}
	
	if req.DepreciationStartDate != nil {
		entity.DepreciationStartDate = req.DepreciationStartDate
	}
	
	if req.AssetAccountId != nil {
		entity.AssetAccountId = req.AssetAccountId
	}
	
	if req.AccumulatedDepreciationAccountId != nil {
		entity.AccumulatedDepreciationAccountId = req.AccumulatedDepreciationAccountId
	}
	
	if req.DepreciationExpenseAccountId != nil {
		entity.DepreciationExpenseAccountId = req.DepreciationExpenseAccountId
	}
	
	if req.CurrentBookValue != nil {
		entity.CurrentBookValue = req.CurrentBookValue
	}
	
	if req.AccumulatedDepreciation != nil {
		entity.AccumulatedDepreciation = req.AccumulatedDepreciation
	}
	
	if req.LastDepreciationDate != nil {
		entity.LastDepreciationDate = req.LastDepreciationDate
	}
	
	if req.LocationId != nil {
		entity.LocationId = req.LocationId
	}
	
	if req.Department != nil {
		entity.Department = req.Department
	}
	
	if req.IsDisposed != nil {
		entity.IsDisposed = req.IsDisposed
	}
	
	if req.DisposalDate != nil {
		entity.DisposalDate = req.DisposalDate
	}
	
	if req.DisposalProceeds != nil {
		entity.DisposalProceeds = req.DisposalProceeds
	}
	
	if req.DisposalJournalEntryId != nil {
		entity.DisposalJournalEntryId = req.DisposalJournalEntryId
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.SerialNumber != nil {
		entity.SerialNumber = req.SerialNumber
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update fixed_assets: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated fixed_assets",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a fixed_assets
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting fixed_assets",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return err
	}

	// Verify ownership
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get fixed_assets: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("fixed_assets not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete fixed_assets: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted fixed_assets",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *FixedAssets) *FixedAssetsResponse {
	return &FixedAssetsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		AssetNumber: entity.AssetNumber,
		
		AssetName: entity.AssetName,
		
		AssetCategoryId: entity.AssetCategoryId,
		
		AcquisitionDate: entity.AcquisitionDate,
		
		AcquisitionCost: entity.AcquisitionCost,
		
		SalvageValue: entity.SalvageValue,
		
		SupplierId: entity.SupplierId,
		
		VendorBillId: entity.VendorBillId,
		
		DepreciationMethod: entity.DepreciationMethod,
		
		UsefulLifeYears: entity.UsefulLifeYears,
		
		DepreciationStartDate: entity.DepreciationStartDate,
		
		AssetAccountId: entity.AssetAccountId,
		
		AccumulatedDepreciationAccountId: entity.AccumulatedDepreciationAccountId,
		
		DepreciationExpenseAccountId: entity.DepreciationExpenseAccountId,
		
		CurrentBookValue: entity.CurrentBookValue,
		
		AccumulatedDepreciation: entity.AccumulatedDepreciation,
		
		LastDepreciationDate: entity.LastDepreciationDate,
		
		LocationId: entity.LocationId,
		
		Department: entity.Department,
		
		IsDisposed: entity.IsDisposed,
		
		DisposalDate: entity.DisposalDate,
		
		DisposalProceeds: entity.DisposalProceeds,
		
		DisposalJournalEntryId: entity.DisposalJournalEntryId,
		
		Description: entity.Description,
		
		SerialNumber: entity.SerialNumber,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
	}
}


// setOrganizationContext sets the organization context for RLS
func (s *Service) setOrganizationContext(ctx context.Context, tx pgx.Tx, orgID uuid.UUID) error {
	_, err := tx.Exec(ctx, "SET LOCAL app.current_organization_id = $1", orgID)
	if err != nil {
		return fmt.Errorf("failed to set organization context: %w", err)
	}
	return nil
}


// validateBusinessRules validates business rules for fixed_assets
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *FixedAssets) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a fixed_assets can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

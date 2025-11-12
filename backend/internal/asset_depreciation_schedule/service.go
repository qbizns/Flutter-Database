package asset_depreciation_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for AssetDepreciationSchedule
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new AssetDepreciationSchedule service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new asset_depreciation_schedule
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateAssetDepreciationScheduleRequest) (*AssetDepreciationScheduleResponse, error) {
	s.logger.Info("creating asset_depreciation_schedule",
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
	entity := &AssetDepreciationSchedule{
		OrganizationID: orgID,
		
		FixedAssetId: req.FixedAssetId,
		
		FiscalYearId: req.FiscalYearId,
		
		AccountingPeriodId: req.AccountingPeriodId,
		
		DepreciationDate: req.DepreciationDate,
		
		DepreciationAmount: req.DepreciationAmount,
		
		AccumulatedDepreciationBeginning: req.AccumulatedDepreciationBeginning,
		
		AccumulatedDepreciationEnding: req.AccumulatedDepreciationEnding,
		
		BookValueBeginning: req.BookValueBeginning,
		
		BookValueEnding: req.BookValueEnding,
		
		JournalEntryId: req.JournalEntryId,
		
		IsPosted: req.IsPosted,
		
		PostedAt: req.PostedAt,
		
		PostedBy: req.PostedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create asset_depreciation_schedule: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created asset_depreciation_schedule",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a asset_depreciation_schedule by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*AssetDepreciationScheduleResponse, error) {
	s.logger.Debug("getting asset_depreciation_schedule",
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
		return nil, fmt.Errorf("failed to get asset_depreciation_schedule: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("asset_depreciation_schedule not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of asset_depreciation_schedule records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*AssetDepreciationScheduleListResponse, error) {
	s.logger.Debug("listing asset_depreciation_schedule",
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
		return nil, fmt.Errorf("failed to list asset_depreciation_schedule: %w", err)
	}

	// Convert to response
	items := make([]*AssetDepreciationScheduleResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &AssetDepreciationScheduleListResponse{
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

// Update updates an existing asset_depreciation_schedule
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateAssetDepreciationScheduleRequest) (*AssetDepreciationScheduleResponse, error) {
	s.logger.Info("updating asset_depreciation_schedule",
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
		return nil, fmt.Errorf("failed to get asset_depreciation_schedule: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("asset_depreciation_schedule not found or access denied")
	}
	

	// Update fields
	
	if req.FixedAssetId != nil {
		entity.FixedAssetId = *req.FixedAssetId
	}
	
	if req.FiscalYearId != nil {
		entity.FiscalYearId = *req.FiscalYearId
	}
	
	if req.AccountingPeriodId != nil {
		entity.AccountingPeriodId = *req.AccountingPeriodId
	}
	
	if req.DepreciationDate != nil {
		entity.DepreciationDate = *req.DepreciationDate
	}
	
	if req.DepreciationAmount != nil {
		entity.DepreciationAmount = *req.DepreciationAmount
	}
	
	if req.AccumulatedDepreciationBeginning != nil {
		entity.AccumulatedDepreciationBeginning = *req.AccumulatedDepreciationBeginning
	}
	
	if req.AccumulatedDepreciationEnding != nil {
		entity.AccumulatedDepreciationEnding = *req.AccumulatedDepreciationEnding
	}
	
	if req.BookValueBeginning != nil {
		entity.BookValueBeginning = *req.BookValueBeginning
	}
	
	if req.BookValueEnding != nil {
		entity.BookValueEnding = *req.BookValueEnding
	}
	
	if req.JournalEntryId != nil {
		entity.JournalEntryId = *req.JournalEntryId
	}
	
	if req.IsPosted != nil {
		entity.IsPosted = *req.IsPosted
	}
	
	if req.PostedAt != nil {
		entity.PostedAt = *req.PostedAt
	}
	
	if req.PostedBy != nil {
		entity.PostedBy = *req.PostedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update asset_depreciation_schedule: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated asset_depreciation_schedule",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a asset_depreciation_schedule
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting asset_depreciation_schedule",
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
		return fmt.Errorf("failed to get asset_depreciation_schedule: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("asset_depreciation_schedule not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete asset_depreciation_schedule: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted asset_depreciation_schedule",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *AssetDepreciationSchedule) *AssetDepreciationScheduleResponse {
	return &AssetDepreciationScheduleResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		FixedAssetId: entity.FixedAssetId,
		
		FiscalYearId: entity.FiscalYearId,
		
		AccountingPeriodId: entity.AccountingPeriodId,
		
		DepreciationDate: entity.DepreciationDate,
		
		DepreciationAmount: entity.DepreciationAmount,
		
		AccumulatedDepreciationBeginning: entity.AccumulatedDepreciationBeginning,
		
		AccumulatedDepreciationEnding: entity.AccumulatedDepreciationEnding,
		
		BookValueBeginning: entity.BookValueBeginning,
		
		BookValueEnding: entity.BookValueEnding,
		
		JournalEntryId: entity.JournalEntryId,
		
		IsPosted: entity.IsPosted,
		
		CreatedAt: entity.CreatedAt,
		
		PostedAt: entity.PostedAt,
		
		PostedBy: entity.PostedBy,
		
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


// validateBusinessRules validates business rules for asset_depreciation_schedule
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *AssetDepreciationSchedule) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a asset_depreciation_schedule can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

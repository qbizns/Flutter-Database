package asset_category

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/asset_category"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/asset_category"
	"go.uber.org/zap"
)

// Service handles business logic for AssetCategories
type Service struct {
	repo   *asset_category.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new AssetCategories service
func NewService(repo *asset_category.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new asset_categories
func (s *Service) Create(ctx context.Context, req *dto.CreateAssetCategoriesRequest) (*dto.AssetCategoriesResponse, error) {
	s.logger.Info("creating asset_categories",
		
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

	

	// Convert DTO to entity
	entity := &asset_category.AssetCategories{
		
		
		CategoryCode: req.CategoryCode,
		
		CategoryName: req.CategoryName,
		
		DefaultDepreciationMethod: req.DefaultDepreciationMethod,
		
		DefaultUsefulLifeYears: req.DefaultUsefulLifeYears,
		
		DefaultSalvageValuePercent: req.DefaultSalvageValuePercent,
		
		AssetAccountId: req.AssetAccountId,
		
		AccumulatedDepreciationAccountId: req.AccumulatedDepreciationAccountId,
		
		DepreciationExpenseAccountId: req.DepreciationExpenseAccountId,
		
		Description: req.Description,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create asset_categories: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created asset_categories",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a asset_categories by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*dto.AssetCategoriesResponse, error) {
	s.logger.Debug("getting asset_categories",
		zap.String("id", id.String()),
		
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset_categories: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of asset_categories records
func (s *Service) List(ctx context.Context, page, limit int) (*dto.AssetCategoriesListResponse, error) {
	s.logger.Debug("listing asset_categories",
		
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

	
	// Get from database
	entities, total, err := s.repo.List(ctx, tx, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list asset_categories: %w", err)
	}

	// Convert to response
	items := make([]*dto.AssetCategoriesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.AssetCategoriesListResponse{
		Items: items,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing asset_categories
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateAssetCategoriesRequest) (*dto.AssetCategoriesResponse, error) {
	s.logger.Info("updating asset_categories",
		zap.String("id", id.String()),
		
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

	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset_categories: %w", err)
	}

	

	// Update fields
	
	if req.CategoryCode != nil {
		entity.CategoryCode = *req.CategoryCode
	}
	
	if req.CategoryName != nil {
		entity.CategoryName = *req.CategoryName
	}
	
	if req.DefaultDepreciationMethod != nil {
		entity.DefaultDepreciationMethod = *req.DefaultDepreciationMethod
	}
	
	if req.DefaultUsefulLifeYears != nil {
		entity.DefaultUsefulLifeYears = *req.DefaultUsefulLifeYears
	}
	
	if req.DefaultSalvageValuePercent != nil {
		entity.DefaultSalvageValuePercent = *req.DefaultSalvageValuePercent
	}
	
	if req.AssetAccountId != nil {
		entity.AssetAccountId = *req.AssetAccountId
	}
	
	if req.AccumulatedDepreciationAccountId != nil {
		entity.AccumulatedDepreciationAccountId = *req.AccumulatedDepreciationAccountId
	}
	
	if req.DepreciationExpenseAccountId != nil {
		entity.DepreciationExpenseAccountId = *req.DepreciationExpenseAccountId
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update asset_categories: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated asset_categories",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a asset_categories
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting asset_categories",
		zap.String("id", id.String()),
		
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete asset_categories: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted asset_categories",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *asset_category.AssetCategories) *dto.AssetCategoriesResponse {
	return &dto.AssetCategoriesResponse{
		
		Id: entity.Id,
		
		CategoryCode: entity.CategoryCode,
		
		CategoryName: entity.CategoryName,
		
		DefaultDepreciationMethod: entity.DefaultDepreciationMethod,
		
		DefaultUsefulLifeYears: entity.DefaultUsefulLifeYears,
		
		DefaultSalvageValuePercent: entity.DefaultSalvageValuePercent,
		
		AssetAccountId: entity.AssetAccountId,
		
		AccumulatedDepreciationAccountId: entity.AccumulatedDepreciationAccountId,
		
		DepreciationExpenseAccountId: entity.DepreciationExpenseAccountId,
		
		Description: entity.Description,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
	}
}



// validateBusinessRules validates business rules for asset_categories
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *asset_category.AssetCategories) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a asset_categories can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

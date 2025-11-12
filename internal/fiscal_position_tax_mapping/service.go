package fiscal_position_tax_mapping

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/fiscal_position_tax_mapping"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/fiscal_position_tax_mapping"
	"go.uber.org/zap"
)

// Service handles business logic for FiscalPositionTaxMappings
type Service struct {
	repo   *fiscal_position_tax_mapping.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new FiscalPositionTaxMappings service
func NewService(repo *fiscal_position_tax_mapping.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new fiscal_position_tax_mappings
func (s *Service) Create(ctx context.Context, req *dto.CreateFiscalPositionTaxMappingsRequest) (*dto.FiscalPositionTaxMappingsResponse, error) {
	s.logger.Info("creating fiscal_position_tax_mappings",
		
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
	entity := &fiscal_position_tax_mapping.FiscalPositionTaxMappings{
		
		
		FiscalPositionId: req.FiscalPositionId,
		
		SourceTaxId: req.SourceTaxId,
		
		DestinationTaxId: req.DestinationTaxId,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create fiscal_position_tax_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created fiscal_position_tax_mappings",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a fiscal_position_tax_mappings by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*dto.FiscalPositionTaxMappingsResponse, error) {
	s.logger.Debug("getting fiscal_position_tax_mappings",
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
		return nil, fmt.Errorf("failed to get fiscal_position_tax_mappings: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of fiscal_position_tax_mappings records
func (s *Service) List(ctx context.Context, page, limit int) (*dto.FiscalPositionTaxMappingsListResponse, error) {
	s.logger.Debug("listing fiscal_position_tax_mappings",
		
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
		return nil, fmt.Errorf("failed to list fiscal_position_tax_mappings: %w", err)
	}

	// Convert to response
	items := make([]*dto.FiscalPositionTaxMappingsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.FiscalPositionTaxMappingsListResponse{
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

// Update updates an existing fiscal_position_tax_mappings
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateFiscalPositionTaxMappingsRequest) (*dto.FiscalPositionTaxMappingsResponse, error) {
	s.logger.Info("updating fiscal_position_tax_mappings",
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
		return nil, fmt.Errorf("failed to get fiscal_position_tax_mappings: %w", err)
	}

	

	// Update fields
	
	if req.FiscalPositionId != nil {
		entity.FiscalPositionId = *req.FiscalPositionId
	}
	
	if req.SourceTaxId != nil {
		entity.SourceTaxId = *req.SourceTaxId
	}
	
	if req.DestinationTaxId != nil {
		entity.DestinationTaxId = *req.DestinationTaxId
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update fiscal_position_tax_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated fiscal_position_tax_mappings",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a fiscal_position_tax_mappings
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting fiscal_position_tax_mappings",
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
		return fmt.Errorf("failed to delete fiscal_position_tax_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted fiscal_position_tax_mappings",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *fiscal_position_tax_mapping.FiscalPositionTaxMappings) *dto.FiscalPositionTaxMappingsResponse {
	return &dto.FiscalPositionTaxMappingsResponse{
		
		Id: entity.Id,
		
		FiscalPositionId: entity.FiscalPositionId,
		
		SourceTaxId: entity.SourceTaxId,
		
		DestinationTaxId: entity.DestinationTaxId,
		
		CreatedBy: entity.CreatedBy,
		
		CreatedAt: entity.CreatedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for fiscal_position_tax_mappings
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *fiscal_position_tax_mapping.FiscalPositionTaxMappings) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a fiscal_position_tax_mappings can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

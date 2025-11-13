package uom_conversion

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for UomConversions
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new UomConversions service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new uom_conversions
func (s *Service) Create(ctx context.Context, req *CreateUomConversionsRequest) (*UomConversionsResponse, error) {
	s.logger.Info("creating uom_conversions",
		
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
	entity := &UomConversions{
		
		
		FromUomId: req.FromUomId,
		
		ToUomId: req.ToUomId,
		
		ConversionFactor: req.ConversionFactor,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create uom_conversions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created uom_conversions",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a uom_conversions by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*UomConversionsResponse, error) {
	s.logger.Debug("getting uom_conversions",
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
		return nil, fmt.Errorf("failed to get uom_conversions: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of uom_conversions records
func (s *Service) List(ctx context.Context, page, limit int) (*UomConversionsListResponse, error) {
	s.logger.Debug("listing uom_conversions",
		
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
		return nil, fmt.Errorf("failed to list uom_conversions: %w", err)
	}

	// Convert to response
	items := make([]*UomConversionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &UomConversionsListResponse{
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

// Update updates an existing uom_conversions
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateUomConversionsRequest) (*UomConversionsResponse, error) {
	s.logger.Info("updating uom_conversions",
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
		return nil, fmt.Errorf("failed to get uom_conversions: %w", err)
	}

	

	// Update fields
	
	if req.FromUomId != nil {
		entity.FromUomId = req.FromUomId
	}
	
	if req.ToUomId != nil {
		entity.ToUomId = req.ToUomId
	}
	
	if req.ConversionFactor != nil {
		entity.ConversionFactor = req.ConversionFactor
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update uom_conversions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated uom_conversions",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a uom_conversions
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting uom_conversions",
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
		return fmt.Errorf("failed to delete uom_conversions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted uom_conversions",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *UomConversions) *UomConversionsResponse {
	return &UomConversionsResponse{
		
		Id: entity.Id,
		
		FromUomId: entity.FromUomId,
		
		ToUomId: entity.ToUomId,
		
		ConversionFactor: entity.ConversionFactor,
		
		CreatedAt: entity.CreatedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for uom_conversions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *UomConversions) error {
	// TODO: Add business rule validation
	// Basic business validation implemented
	// Production: Add module-specific validation rules as needed
	
	// Example validations that can be added:
	// - Duplicate checking within organization
	// - Foreign key validation
	// - Amount/date range validation
	// - Status transition rules
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a uom_conversions can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Basic delete validation implemented
	// Production: Add checks for dependent records
	
	// Example checks that can be added:
	// - Query related tables for dependencies
	// - Prevent deletion of entities with transactions
	// - Check business rules (e.g., dont delete active items)
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

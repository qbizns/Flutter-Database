package stock_adjustment_reason

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for StockAdjustmentReasons
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new StockAdjustmentReasons service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new stock_adjustment_reasons
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateStockAdjustmentReasonsRequest) (*StockAdjustmentReasonsResponse, error) {
	s.logger.Info("creating stock_adjustment_reasons",
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
	entity := &StockAdjustmentReasons{
		OrganizationId: orgID,
		
		Code: req.Code,
		
		Name: req.Name,
		
		Description: req.Description,
		
		ReasonType: req.ReasonType,
		
		IsSystemReason: req.IsSystemReason,
		
		IsActive: req.IsActive,
		
		RequiresApproval: req.RequiresApproval,
		
		RequiresNotes: req.RequiresNotes,
		
		SortOrder: req.SortOrder,
		
		Metadata: req.Metadata,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create stock_adjustment_reasons: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created stock_adjustment_reasons",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a stock_adjustment_reasons by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*StockAdjustmentReasonsResponse, error) {
	s.logger.Debug("getting stock_adjustment_reasons",
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
		return nil, fmt.Errorf("failed to get stock_adjustment_reasons: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("stock_adjustment_reasons not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of stock_adjustment_reasons records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*StockAdjustmentReasonsListResponse, error) {
	s.logger.Debug("listing stock_adjustment_reasons",
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
		return nil, fmt.Errorf("failed to list stock_adjustment_reasons: %w", err)
	}

	// Convert to response
	items := make([]*StockAdjustmentReasonsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &StockAdjustmentReasonsListResponse{
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

// Update updates an existing stock_adjustment_reasons
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateStockAdjustmentReasonsRequest) (*StockAdjustmentReasonsResponse, error) {
	s.logger.Info("updating stock_adjustment_reasons",
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
		return nil, fmt.Errorf("failed to get stock_adjustment_reasons: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("stock_adjustment_reasons not found or access denied")
	}
	

	// Update fields
	
	if req.Code != nil {
		entity.Code = req.Code
	}
	
	if req.Name != nil {
		entity.Name = req.Name
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.ReasonType != nil {
		entity.ReasonType = req.ReasonType
	}
	
	if req.IsSystemReason != nil {
		entity.IsSystemReason = req.IsSystemReason
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.RequiresApproval != nil {
		entity.RequiresApproval = req.RequiresApproval
	}
	
	if req.RequiresNotes != nil {
		entity.RequiresNotes = req.RequiresNotes
	}
	
	if req.SortOrder != nil {
		entity.SortOrder = req.SortOrder
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update stock_adjustment_reasons: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated stock_adjustment_reasons",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a stock_adjustment_reasons
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting stock_adjustment_reasons",
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
		return fmt.Errorf("failed to get stock_adjustment_reasons: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("stock_adjustment_reasons not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete stock_adjustment_reasons: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted stock_adjustment_reasons",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *StockAdjustmentReasons) *StockAdjustmentReasonsResponse {
	return &StockAdjustmentReasonsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		Code: entity.Code,
		
		Name: entity.Name,
		
		Description: entity.Description,
		
		ReasonType: entity.ReasonType,
		
		IsSystemReason: entity.IsSystemReason,
		
		IsActive: entity.IsActive,
		
		RequiresApproval: entity.RequiresApproval,
		
		RequiresNotes: entity.RequiresNotes,
		
		SortOrder: entity.SortOrder,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
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


// validateBusinessRules validates business rules for stock_adjustment_reasons
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *StockAdjustmentReasons) error {
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

// canDelete checks if a stock_adjustment_reasons can be deleted
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

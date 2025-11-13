package cycle_count_item

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for CycleCountItems
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new CycleCountItems service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new cycle_count_items
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateCycleCountItemsRequest) (*CycleCountItemsResponse, error) {
	s.logger.Info("creating cycle_count_items",
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
	entity := &CycleCountItems{
		OrganizationId: orgID,
		
		CycleCountId: req.CycleCountId,
		
		ProductId: req.ProductId,
		
		ProductVariantId: req.ProductVariantId,
		
		ProductName: req.ProductName,
		
		ProductSku: req.ProductSku,
		
		SystemQuantity: req.SystemQuantity,
		
		CountedQuantity: req.CountedQuantity,
		
		VarianceQuantity: req.VarianceQuantity,
		
		VariancePercentage: req.VariancePercentage,
		
		UnitCost: req.UnitCost,
		
		VarianceValue: req.VarianceValue,
		
		Status: req.Status,
		
		RecountRequired: req.RecountRequired,
		
		RecountQuantity: req.RecountQuantity,
		
		RecountReason: req.RecountReason,
		
		AdjustmentApplied: req.AdjustmentApplied,
		
		AdjustmentDate: req.AdjustmentDate,
		
		AdjustmentReason: req.AdjustmentReason,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CountedAt: req.CountedAt,
		
		CountedBy: req.CountedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create cycle_count_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created cycle_count_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a cycle_count_items by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CycleCountItemsResponse, error) {
	s.logger.Debug("getting cycle_count_items",
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
		return nil, fmt.Errorf("failed to get cycle_count_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("cycle_count_items not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of cycle_count_items records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*CycleCountItemsListResponse, error) {
	s.logger.Debug("listing cycle_count_items",
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
		return nil, fmt.Errorf("failed to list cycle_count_items: %w", err)
	}

	// Convert to response
	items := make([]*CycleCountItemsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &CycleCountItemsListResponse{
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

// Update updates an existing cycle_count_items
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateCycleCountItemsRequest) (*CycleCountItemsResponse, error) {
	s.logger.Info("updating cycle_count_items",
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
		return nil, fmt.Errorf("failed to get cycle_count_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("cycle_count_items not found or access denied")
	}
	

	// Update fields
	
	if req.CycleCountId != nil {
		entity.CycleCountId = *req.CycleCountId
	}
	
	if req.ProductId != nil {
		entity.ProductId = *req.ProductId
	}
	
	if req.ProductVariantId != nil {
		entity.ProductVariantId = *req.ProductVariantId
	}
	
	if req.ProductName != nil {
		entity.ProductName = *req.ProductName
	}
	
	if req.ProductSku != nil {
		entity.ProductSku = *req.ProductSku
	}
	
	if req.SystemQuantity != nil {
		entity.SystemQuantity = *req.SystemQuantity
	}
	
	if req.CountedQuantity != nil {
		entity.CountedQuantity = *req.CountedQuantity
	}
	
	if req.VarianceQuantity != nil {
		entity.VarianceQuantity = *req.VarianceQuantity
	}
	
	if req.VariancePercentage != nil {
		entity.VariancePercentage = *req.VariancePercentage
	}
	
	if req.UnitCost != nil {
		entity.UnitCost = *req.UnitCost
	}
	
	if req.VarianceValue != nil {
		entity.VarianceValue = *req.VarianceValue
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.RecountRequired != nil {
		entity.RecountRequired = *req.RecountRequired
	}
	
	if req.RecountQuantity != nil {
		entity.RecountQuantity = *req.RecountQuantity
	}
	
	if req.RecountReason != nil {
		entity.RecountReason = *req.RecountReason
	}
	
	if req.AdjustmentApplied != nil {
		entity.AdjustmentApplied = *req.AdjustmentApplied
	}
	
	if req.AdjustmentDate != nil {
		entity.AdjustmentDate = *req.AdjustmentDate
	}
	
	if req.AdjustmentReason != nil {
		entity.AdjustmentReason = *req.AdjustmentReason
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	
	if req.CountedAt != nil {
		entity.CountedAt = req.CountedAt
	}
	
	if req.CountedBy != nil {
		entity.CountedBy = *req.CountedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update cycle_count_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated cycle_count_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a cycle_count_items
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting cycle_count_items",
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
		return fmt.Errorf("failed to get cycle_count_items: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("cycle_count_items not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete cycle_count_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted cycle_count_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *CycleCountItems) *CycleCountItemsResponse {
	return &CycleCountItemsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		CycleCountId: entity.CycleCountId,
		
		ProductId: entity.ProductId,
		
		ProductVariantId: entity.ProductVariantId,
		
		ProductName: entity.ProductName,
		
		ProductSku: entity.ProductSku,
		
		SystemQuantity: entity.SystemQuantity,
		
		CountedQuantity: entity.CountedQuantity,
		
		VarianceQuantity: entity.VarianceQuantity,
		
		VariancePercentage: entity.VariancePercentage,
		
		UnitCost: entity.UnitCost,
		
		VarianceValue: entity.VarianceValue,
		
		Status: entity.Status,
		
		RecountRequired: entity.RecountRequired,
		
		RecountQuantity: entity.RecountQuantity,
		
		RecountReason: entity.RecountReason,
		
		AdjustmentApplied: entity.AdjustmentApplied,
		
		AdjustmentDate: entity.AdjustmentDate,
		
		AdjustmentReason: entity.AdjustmentReason,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CountedAt: entity.CountedAt,
		
		CountedBy: entity.CountedBy,
		
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


// validateBusinessRules validates business rules for cycle_count_items
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *CycleCountItems) error {
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

// canDelete checks if a cycle_count_items can be deleted
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

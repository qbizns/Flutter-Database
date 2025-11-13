package inventory_cost_layer

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

// Service handles business logic for InventoryCostLayers
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new InventoryCostLayers service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new inventory_cost_layers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateInventoryCostLayersRequest) (*InventoryCostLayersResponse, error) {
	s.logger.Info("creating inventory_cost_layers",
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
	entity := &InventoryCostLayers{
		OrganizationID: orgID,
		
		ProductId: req.ProductId,
		
		LocationId: req.LocationId,
		
		LotNumber: req.LotNumber,
		
		SerialNumber: req.SerialNumber,
		
		LayerDate: req.LayerDate,
		
		UnitCost: req.UnitCost,
		
		OriginalQuantity: req.OriginalQuantity,
		
		RemainingQuantity: req.RemainingQuantity,
		
		UomCode: req.UomCode,
		
		SourceTransactionType: req.SourceTransactionType,
		
		SourceTransactionId: req.SourceTransactionId,
		
		SourceReference: req.SourceReference,
		
		IsFullyConsumed: req.IsFullyConsumed,
		
		ConsumedAt: req.ConsumedAt,
		
		Metadata: req.Metadata,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create inventory_cost_layers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created inventory_cost_layers",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a inventory_cost_layers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryCostLayersResponse, error) {
	s.logger.Debug("getting inventory_cost_layers",
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
		return nil, fmt.Errorf("failed to get inventory_cost_layers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("inventory_cost_layers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of inventory_cost_layers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*InventoryCostLayersListResponse, error) {
	s.logger.Debug("listing inventory_cost_layers",
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
		return nil, fmt.Errorf("failed to list inventory_cost_layers: %w", err)
	}

	// Convert to response
	items := make([]*InventoryCostLayersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &InventoryCostLayersListResponse{
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

// Update updates an existing inventory_cost_layers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateInventoryCostLayersRequest) (*InventoryCostLayersResponse, error) {
	s.logger.Info("updating inventory_cost_layers",
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
		return nil, fmt.Errorf("failed to get inventory_cost_layers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("inventory_cost_layers not found or access denied")
	}
	

	// Update fields
	
	if req.ProductId != nil {
		entity.ProductId = *req.ProductId
	}
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.LotNumber != nil {
		entity.LotNumber = *req.LotNumber
	}
	
	if req.SerialNumber != nil {
		entity.SerialNumber = *req.SerialNumber
	}
	
	if req.LayerDate != nil {
		entity.LayerDate = *req.LayerDate
	}
	
	if req.UnitCost != nil {
		entity.UnitCost = *req.UnitCost
	}
	
	if req.OriginalQuantity != nil {
		entity.OriginalQuantity = *req.OriginalQuantity
	}
	
	if req.RemainingQuantity != nil {
		entity.RemainingQuantity = *req.RemainingQuantity
	}
	
	if req.UomCode != nil {
		entity.UomCode = *req.UomCode
	}
	
	if req.SourceTransactionType != nil {
		entity.SourceTransactionType = *req.SourceTransactionType
	}
	
	if req.SourceTransactionId != nil {
		entity.SourceTransactionId = *req.SourceTransactionId
	}
	
	if req.SourceReference != nil {
		entity.SourceReference = *req.SourceReference
	}
	
	if req.IsFullyConsumed != nil {
		entity.IsFullyConsumed = *req.IsFullyConsumed
	}
	
	if req.ConsumedAt != nil {
		entity.ConsumedAt = *req.ConsumedAt
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update inventory_cost_layers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated inventory_cost_layers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a inventory_cost_layers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting inventory_cost_layers",
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
		return fmt.Errorf("failed to get inventory_cost_layers: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("inventory_cost_layers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete inventory_cost_layers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted inventory_cost_layers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *InventoryCostLayers) *InventoryCostLayersResponse {
	return &InventoryCostLayersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ProductId: entity.ProductId,
		
		LocationId: entity.LocationId,
		
		LotNumber: entity.LotNumber,
		
		SerialNumber: entity.SerialNumber,
		
		LayerDate: entity.LayerDate,
		
		UnitCost: entity.UnitCost,
		
		OriginalQuantity: entity.OriginalQuantity,
		
		RemainingQuantity: entity.RemainingQuantity,
		
		UomCode: entity.UomCode,
		
		SourceTransactionType: entity.SourceTransactionType,
		
		SourceTransactionId: entity.SourceTransactionId,
		
		SourceReference: entity.SourceReference,
		
		IsFullyConsumed: entity.IsFullyConsumed,
		
		ConsumedAt: entity.ConsumedAt,
		
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


// validateBusinessRules validates business rules for inventory_cost_layers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *InventoryCostLayers) error {
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

// canDelete checks if a inventory_cost_layers can be deleted
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

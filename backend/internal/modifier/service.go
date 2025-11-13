package modifier

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for Modifiers
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Modifiers service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new modifiers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateModifiersRequest) (*ModifiersResponse, error) {
	s.logger.Info("creating modifiers",
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
	entity := &Modifiers{
		OrganizationID: orgID,
		
		ModifierGroupId: req.ModifierGroupId,
		
		ModifierName: req.ModifierName,
		
		ModifierCode: req.ModifierCode,
		
		DisplayName: req.DisplayName,
		
		PriceAdjustment: req.PriceAdjustment,
		
		PriceType: req.PriceType,
		
		IsAvailable: req.IsAvailable,
		
		IsDefault: req.IsDefault,
		
		TrackInventory: req.TrackInventory,
		
		CurrentStock: req.CurrentStock,
		
		LowStockThreshold: req.LowStockThreshold,
		
		DisplayOrder: req.DisplayOrder,
		
		ImageUrl: req.ImageUrl,
		
		Description: req.Description,
		
		AllergenInfo: req.AllergenInfo,
		
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
		return nil, fmt.Errorf("failed to create modifiers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created modifiers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a modifiers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ModifiersResponse, error) {
	s.logger.Debug("getting modifiers",
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
		return nil, fmt.Errorf("failed to get modifiers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("modifiers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of modifiers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ModifiersListResponse, error) {
	s.logger.Debug("listing modifiers",
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
		return nil, fmt.Errorf("failed to list modifiers: %w", err)
	}

	// Convert to response
	items := make([]*ModifiersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ModifiersListResponse{
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

// Update updates an existing modifiers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateModifiersRequest) (*ModifiersResponse, error) {
	s.logger.Info("updating modifiers",
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
		return nil, fmt.Errorf("failed to get modifiers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("modifiers not found or access denied")
	}
	

	// Update fields
	
	if req.ModifierGroupId != nil {
		entity.ModifierGroupId = req.ModifierGroupId
	}
	
	if req.ModifierName != nil {
		entity.ModifierName = req.ModifierName
	}
	
	if req.ModifierCode != nil {
		entity.ModifierCode = req.ModifierCode
	}
	
	if req.DisplayName != nil {
		entity.DisplayName = req.DisplayName
	}
	
	if req.PriceAdjustment != nil {
		entity.PriceAdjustment = req.PriceAdjustment
	}
	
	if req.PriceType != nil {
		entity.PriceType = req.PriceType
	}
	
	if req.IsAvailable != nil {
		entity.IsAvailable = req.IsAvailable
	}
	
	if req.IsDefault != nil {
		entity.IsDefault = req.IsDefault
	}
	
	if req.TrackInventory != nil {
		entity.TrackInventory = req.TrackInventory
	}
	
	if req.CurrentStock != nil {
		entity.CurrentStock = req.CurrentStock
	}
	
	if req.LowStockThreshold != nil {
		entity.LowStockThreshold = req.LowStockThreshold
	}
	
	if req.DisplayOrder != nil {
		entity.DisplayOrder = req.DisplayOrder
	}
	
	if req.ImageUrl != nil {
		entity.ImageUrl = req.ImageUrl
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.AllergenInfo != nil {
		entity.AllergenInfo = req.AllergenInfo
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
		return nil, fmt.Errorf("failed to update modifiers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated modifiers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a modifiers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting modifiers",
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
		return fmt.Errorf("failed to get modifiers: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("modifiers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete modifiers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted modifiers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Modifiers) *ModifiersResponse {
	return &ModifiersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ModifierGroupId: entity.ModifierGroupId,
		
		ModifierName: entity.ModifierName,
		
		ModifierCode: entity.ModifierCode,
		
		DisplayName: entity.DisplayName,
		
		PriceAdjustment: entity.PriceAdjustment,
		
		PriceType: entity.PriceType,
		
		IsAvailable: entity.IsAvailable,
		
		IsDefault: entity.IsDefault,
		
		TrackInventory: entity.TrackInventory,
		
		CurrentStock: entity.CurrentStock,
		
		LowStockThreshold: entity.LowStockThreshold,
		
		DisplayOrder: entity.DisplayOrder,
		
		ImageUrl: entity.ImageUrl,
		
		Description: entity.Description,
		
		AllergenInfo: entity.AllergenInfo,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for modifiers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Modifiers) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a modifiers can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

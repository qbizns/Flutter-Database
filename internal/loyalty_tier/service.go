package loyalty_tier

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/loyalty_tier"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/loyalty_tier"
	"go.uber.org/zap"
)

// Service handles business logic for LoyaltyTiers
type Service struct {
	repo   *loyalty_tier.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new LoyaltyTiers service
func NewService(repo *loyalty_tier.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new loyalty_tiers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateLoyaltyTiersRequest) (*dto.LoyaltyTiersResponse, error) {
	s.logger.Info("creating loyalty_tiers",
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
	entity := &loyalty_tier.LoyaltyTiers{
		OrganizationID: orgID,
		
		TierCode: req.TierCode,
		
		TierName: req.TierName,
		
		TierLevel: req.TierLevel,
		
		Description: req.Description,
		
		PointsThreshold: req.PointsThreshold,
		
		AnnualSpendThreshold: req.AnnualSpendThreshold,
		
		PurchaseCountThreshold: req.PurchaseCountThreshold,
		
		PointsMultiplier: req.PointsMultiplier,
		
		DiscountPercentage: req.DiscountPercentage,
		
		TierColor: req.TierColor,
		
		TierIcon: req.TierIcon,
		
		BadgeImageUrl: req.BadgeImageUrl,
		
		IsActive: req.IsActive,
		
		IsDefault: req.IsDefault,
		
		SortOrder: req.SortOrder,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		PointsThreshold: req.PointsThreshold,
		
		(annualSpendThreshold: req.(annualSpendThreshold,
		
		(purchaseCountThreshold: req.(purchaseCountThreshold,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create loyalty_tiers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created loyalty_tiers",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a loyalty_tiers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.LoyaltyTiersResponse, error) {
	s.logger.Debug("getting loyalty_tiers",
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
		return nil, fmt.Errorf("failed to get loyalty_tiers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_tiers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of loyalty_tiers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.LoyaltyTiersListResponse, error) {
	s.logger.Debug("listing loyalty_tiers",
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
		return nil, fmt.Errorf("failed to list loyalty_tiers: %w", err)
	}

	// Convert to response
	items := make([]*dto.LoyaltyTiersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.LoyaltyTiersListResponse{
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

// Update updates an existing loyalty_tiers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateLoyaltyTiersRequest) (*dto.LoyaltyTiersResponse, error) {
	s.logger.Info("updating loyalty_tiers",
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
		return nil, fmt.Errorf("failed to get loyalty_tiers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_tiers not found or access denied")
	}
	

	// Update fields
	
	if req.TierCode != nil {
		entity.TierCode = *req.TierCode
	}
	
	if req.TierName != nil {
		entity.TierName = *req.TierName
	}
	
	if req.TierLevel != nil {
		entity.TierLevel = *req.TierLevel
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.PointsThreshold != nil {
		entity.PointsThreshold = *req.PointsThreshold
	}
	
	if req.AnnualSpendThreshold != nil {
		entity.AnnualSpendThreshold = *req.AnnualSpendThreshold
	}
	
	if req.PurchaseCountThreshold != nil {
		entity.PurchaseCountThreshold = *req.PurchaseCountThreshold
	}
	
	if req.PointsMultiplier != nil {
		entity.PointsMultiplier = *req.PointsMultiplier
	}
	
	if req.DiscountPercentage != nil {
		entity.DiscountPercentage = *req.DiscountPercentage
	}
	
	if req.TierColor != nil {
		entity.TierColor = *req.TierColor
	}
	
	if req.TierIcon != nil {
		entity.TierIcon = *req.TierIcon
	}
	
	if req.BadgeImageUrl != nil {
		entity.BadgeImageUrl = *req.BadgeImageUrl
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.IsDefault != nil {
		entity.IsDefault = *req.IsDefault
	}
	
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = *req.UpdatedBy
	}
	
	if req.PointsThreshold != nil {
		entity.PointsThreshold = *req.PointsThreshold
	}
	
	if req.(annualSpendThreshold != nil {
		entity.(annualSpendThreshold = *req.(annualSpendThreshold
	}
	
	if req.(purchaseCountThreshold != nil {
		entity.(purchaseCountThreshold = *req.(purchaseCountThreshold
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update loyalty_tiers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated loyalty_tiers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a loyalty_tiers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting loyalty_tiers",
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
		return fmt.Errorf("failed to get loyalty_tiers: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("loyalty_tiers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete loyalty_tiers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted loyalty_tiers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *loyalty_tier.LoyaltyTiers) *dto.LoyaltyTiersResponse {
	return &dto.LoyaltyTiersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		TierCode: entity.TierCode,
		
		TierName: entity.TierName,
		
		TierLevel: entity.TierLevel,
		
		Description: entity.Description,
		
		PointsThreshold: entity.PointsThreshold,
		
		AnnualSpendThreshold: entity.AnnualSpendThreshold,
		
		PurchaseCountThreshold: entity.PurchaseCountThreshold,
		
		PointsMultiplier: entity.PointsMultiplier,
		
		DiscountPercentage: entity.DiscountPercentage,
		
		TierColor: entity.TierColor,
		
		TierIcon: entity.TierIcon,
		
		BadgeImageUrl: entity.BadgeImageUrl,
		
		IsActive: entity.IsActive,
		
		IsDefault: entity.IsDefault,
		
		SortOrder: entity.SortOrder,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		PointsThreshold: entity.PointsThreshold,
		
		(annualSpendThreshold: entity.(annualSpendThreshold,
		
		(purchaseCountThreshold: entity.(purchaseCountThreshold,
		
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


// validateBusinessRules validates business rules for loyalty_tiers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *loyalty_tier.LoyaltyTiers) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a loyalty_tiers can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

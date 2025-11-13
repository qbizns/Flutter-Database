package loyalty_reward

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for LoyaltyRewards
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new LoyaltyRewards service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new loyalty_rewards
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateLoyaltyRewardsRequest) (*LoyaltyRewardsResponse, error) {
	s.logger.Info("creating loyalty_rewards",
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
	entity := &LoyaltyRewards{
		OrganizationID: orgID,
		
		RewardCode: req.RewardCode,
		
		RewardName: req.RewardName,
		
		Description: req.Description,
		
		RewardType: req.RewardType,
		
		PointsCost: req.PointsCost,
		
		RewardValue: req.RewardValue,
		
		DiscountPercentage: req.DiscountPercentage,
		
		DiscountAmount: req.DiscountAmount,
		
		ProductId: req.ProductId,
		
		ProductVariantId: req.ProductVariantId,
		
		IsActive: req.IsActive,
		
		AvailableFrom: req.AvailableFrom,
		
		AvailableTo: req.AvailableTo,
		
		TotalAvailable: req.TotalAvailable,
		
		TotalRedeemed: req.TotalRedeemed,
		
		MaxRedemptionsPerCustomer: req.MaxRedemptionsPerCustomer,
		
		MinimumTierLevel: req.MinimumTierLevel,
		
		TierIds: req.TierIds,
		
		ImageUrl: req.ImageUrl,
		
		ThumbnailUrl: req.ThumbnailUrl,
		
		Featured: req.Featured,
		
		SortOrder: req.SortOrder,
		
		IsFeatured: req.IsFeatured,
		
		TermsAndConditions: req.TermsAndConditions,
		
		RedemptionInstructions: req.RedemptionInstructions,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		AvailableTo: req.AvailableTo,
		
		TotalRedeemed: req.TotalRedeemed,
		
		(totalAvailable: req.(totalAvailable,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create loyalty_rewards: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created loyalty_rewards",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a loyalty_rewards by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyRewardsResponse, error) {
	s.logger.Debug("getting loyalty_rewards",
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
		return nil, fmt.Errorf("failed to get loyalty_rewards: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("loyalty_rewards not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of loyalty_rewards records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*LoyaltyRewardsListResponse, error) {
	s.logger.Debug("listing loyalty_rewards",
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
		return nil, fmt.Errorf("failed to list loyalty_rewards: %w", err)
	}

	// Convert to response
	items := make([]*LoyaltyRewardsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &LoyaltyRewardsListResponse{
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

// Update updates an existing loyalty_rewards
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateLoyaltyRewardsRequest) (*LoyaltyRewardsResponse, error) {
	s.logger.Info("updating loyalty_rewards",
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
		return nil, fmt.Errorf("failed to get loyalty_rewards: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("loyalty_rewards not found or access denied")
	}
	

	// Update fields
	
	if req.RewardCode != nil {
		entity.RewardCode = req.RewardCode
	}
	
	if req.RewardName != nil {
		entity.RewardName = req.RewardName
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.RewardType != nil {
		entity.RewardType = req.RewardType
	}
	
	if req.PointsCost != nil {
		entity.PointsCost = req.PointsCost
	}
	
	if req.RewardValue != nil {
		entity.RewardValue = req.RewardValue
	}
	
	if req.DiscountPercentage != nil {
		entity.DiscountPercentage = req.DiscountPercentage
	}
	
	if req.DiscountAmount != nil {
		entity.DiscountAmount = req.DiscountAmount
	}
	
	if req.ProductId != nil {
		entity.ProductId = req.ProductId
	}
	
	if req.ProductVariantId != nil {
		entity.ProductVariantId = req.ProductVariantId
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.AvailableFrom != nil {
		entity.AvailableFrom = req.AvailableFrom
	}
	
	if req.AvailableTo != nil {
		entity.AvailableTo = req.AvailableTo
	}
	
	if req.TotalAvailable != nil {
		entity.TotalAvailable = req.TotalAvailable
	}
	
	if req.TotalRedeemed != nil {
		entity.TotalRedeemed = req.TotalRedeemed
	}
	
	if req.MaxRedemptionsPerCustomer != nil {
		entity.MaxRedemptionsPerCustomer = req.MaxRedemptionsPerCustomer
	}
	
	if req.MinimumTierLevel != nil {
		entity.MinimumTierLevel = req.MinimumTierLevel
	}
	
	if req.TierIds != nil {
		entity.TierIds = req.TierIds
	}
	
	if req.ImageUrl != nil {
		entity.ImageUrl = req.ImageUrl
	}
	
	if req.ThumbnailUrl != nil {
		entity.ThumbnailUrl = req.ThumbnailUrl
	}
	
	if req.Featured != nil {
		entity.Featured = req.Featured
	}
	
	if req.SortOrder != nil {
		entity.SortOrder = req.SortOrder
	}
	
	if req.IsFeatured != nil {
		entity.IsFeatured = req.IsFeatured
	}
	
	if req.TermsAndConditions != nil {
		entity.TermsAndConditions = req.TermsAndConditions
	}
	
	if req.RedemptionInstructions != nil {
		entity.RedemptionInstructions = req.RedemptionInstructions
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
	
	if req.AvailableTo != nil {
		entity.AvailableTo = req.AvailableTo
	}
	
	if req.TotalRedeemed != nil {
		entity.TotalRedeemed = req.TotalRedeemed
	}
	
	if req.(totalAvailable != nil {
		entity.(totalAvailable = *req.(totalAvailable
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update loyalty_rewards: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated loyalty_rewards",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a loyalty_rewards
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting loyalty_rewards",
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
		return fmt.Errorf("failed to get loyalty_rewards: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("loyalty_rewards not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete loyalty_rewards: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted loyalty_rewards",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *LoyaltyRewards) *LoyaltyRewardsResponse {
	return &LoyaltyRewardsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		RewardCode: entity.RewardCode,
		
		RewardName: entity.RewardName,
		
		Description: entity.Description,
		
		RewardType: entity.RewardType,
		
		PointsCost: entity.PointsCost,
		
		RewardValue: entity.RewardValue,
		
		DiscountPercentage: entity.DiscountPercentage,
		
		DiscountAmount: entity.DiscountAmount,
		
		ProductId: entity.ProductId,
		
		ProductVariantId: entity.ProductVariantId,
		
		IsActive: entity.IsActive,
		
		AvailableFrom: entity.AvailableFrom,
		
		AvailableTo: entity.AvailableTo,
		
		TotalAvailable: entity.TotalAvailable,
		
		TotalRedeemed: entity.TotalRedeemed,
		
		MaxRedemptionsPerCustomer: entity.MaxRedemptionsPerCustomer,
		
		MinimumTierLevel: entity.MinimumTierLevel,
		
		TierIds: entity.TierIds,
		
		ImageUrl: entity.ImageUrl,
		
		ThumbnailUrl: entity.ThumbnailUrl,
		
		Featured: entity.Featured,
		
		SortOrder: entity.SortOrder,
		
		IsFeatured: entity.IsFeatured,
		
		TermsAndConditions: entity.TermsAndConditions,
		
		RedemptionInstructions: entity.RedemptionInstructions,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		AvailableTo: entity.AvailableTo,
		
		TotalRedeemed: entity.TotalRedeemed,
		
		(totalAvailable: entity.(totalAvailable,
		
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


// validateBusinessRules validates business rules for loyalty_rewards
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *LoyaltyRewards) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a loyalty_rewards can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

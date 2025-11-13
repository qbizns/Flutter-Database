package promotion

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for Promotions
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Promotions service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new promotions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreatePromotionsRequest) (*PromotionsResponse, error) {
	s.logger.Info("creating promotions",
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
	entity := &Promotions{
		OrganizationID: orgID,
		
		PromotionCode: req.PromotionCode,
		
		Name: req.Name,
		
		Description: req.Description,
		
		PromotionType: req.PromotionType,
		
		DiscountValue: req.DiscountValue,
		
		AppliesTo: req.AppliesTo,
		
		ApplicableProductIds: req.ApplicableProductIds,
		
		ApplicableCategoryIds: req.ApplicableCategoryIds,
		
		MinimumPurchaseAmount: req.MinimumPurchaseAmount,
		
		MinimumQuantity: req.MinimumQuantity,
		
		BuyQuantity: req.BuyQuantity,
		
		GetQuantity: req.GetQuantity,
		
		GetDiscountPercentage: req.GetDiscountPercentage,
		
		MaxUsesTotal: req.MaxUsesTotal,
		
		MaxUsesPerCustomer: req.MaxUsesPerCustomer,
		
		CurrentUses: req.CurrentUses,
		
		StartDate: req.StartDate,
		
		EndDate: req.EndDate,
		
		IsActive: req.IsActive,
		
		IsCombinable: req.IsCombinable,
		
		Priority: req.Priority,
		
		TermsAndConditions: req.TermsAndConditions,
		
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
		return nil, fmt.Errorf("failed to create promotions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created promotions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a promotions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PromotionsResponse, error) {
	s.logger.Debug("getting promotions",
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
		return nil, fmt.Errorf("failed to get promotions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("promotions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of promotions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*PromotionsListResponse, error) {
	s.logger.Debug("listing promotions",
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
		return nil, fmt.Errorf("failed to list promotions: %w", err)
	}

	// Convert to response
	items := make([]*PromotionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PromotionsListResponse{
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

// Update updates an existing promotions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdatePromotionsRequest) (*PromotionsResponse, error) {
	s.logger.Info("updating promotions",
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
		return nil, fmt.Errorf("failed to get promotions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("promotions not found or access denied")
	}
	

	// Update fields
	
	if req.PromotionCode != nil {
		entity.PromotionCode = req.PromotionCode
	}
	
	if req.Name != nil {
		entity.Name = req.Name
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.PromotionType != nil {
		entity.PromotionType = req.PromotionType
	}
	
	if req.DiscountValue != nil {
		entity.DiscountValue = req.DiscountValue
	}
	
	if req.AppliesTo != nil {
		entity.AppliesTo = req.AppliesTo
	}
	
	if req.ApplicableProductIds != nil {
		entity.ApplicableProductIds = req.ApplicableProductIds
	}
	
	if req.ApplicableCategoryIds != nil {
		entity.ApplicableCategoryIds = req.ApplicableCategoryIds
	}
	
	if req.MinimumPurchaseAmount != nil {
		entity.MinimumPurchaseAmount = req.MinimumPurchaseAmount
	}
	
	if req.MinimumQuantity != nil {
		entity.MinimumQuantity = req.MinimumQuantity
	}
	
	if req.BuyQuantity != nil {
		entity.BuyQuantity = req.BuyQuantity
	}
	
	if req.GetQuantity != nil {
		entity.GetQuantity = req.GetQuantity
	}
	
	if req.GetDiscountPercentage != nil {
		entity.GetDiscountPercentage = req.GetDiscountPercentage
	}
	
	if req.MaxUsesTotal != nil {
		entity.MaxUsesTotal = req.MaxUsesTotal
	}
	
	if req.MaxUsesPerCustomer != nil {
		entity.MaxUsesPerCustomer = req.MaxUsesPerCustomer
	}
	
	if req.CurrentUses != nil {
		entity.CurrentUses = req.CurrentUses
	}
	
	if req.StartDate != nil {
		entity.StartDate = req.StartDate
	}
	
	if req.EndDate != nil {
		entity.EndDate = req.EndDate
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.IsCombinable != nil {
		entity.IsCombinable = req.IsCombinable
	}
	
	if req.Priority != nil {
		entity.Priority = req.Priority
	}
	
	if req.TermsAndConditions != nil {
		entity.TermsAndConditions = req.TermsAndConditions
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
		return nil, fmt.Errorf("failed to update promotions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated promotions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a promotions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting promotions",
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
		return fmt.Errorf("failed to get promotions: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("promotions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete promotions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted promotions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Promotions) *PromotionsResponse {
	return &PromotionsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		PromotionCode: entity.PromotionCode,
		
		Name: entity.Name,
		
		Description: entity.Description,
		
		PromotionType: entity.PromotionType,
		
		DiscountValue: entity.DiscountValue,
		
		AppliesTo: entity.AppliesTo,
		
		ApplicableProductIds: entity.ApplicableProductIds,
		
		ApplicableCategoryIds: entity.ApplicableCategoryIds,
		
		MinimumPurchaseAmount: entity.MinimumPurchaseAmount,
		
		MinimumQuantity: entity.MinimumQuantity,
		
		BuyQuantity: entity.BuyQuantity,
		
		GetQuantity: entity.GetQuantity,
		
		GetDiscountPercentage: entity.GetDiscountPercentage,
		
		MaxUsesTotal: entity.MaxUsesTotal,
		
		MaxUsesPerCustomer: entity.MaxUsesPerCustomer,
		
		CurrentUses: entity.CurrentUses,
		
		StartDate: entity.StartDate,
		
		EndDate: entity.EndDate,
		
		IsActive: entity.IsActive,
		
		IsCombinable: entity.IsCombinable,
		
		Priority: entity.Priority,
		
		TermsAndConditions: entity.TermsAndConditions,
		
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


// validateBusinessRules validates business rules for promotions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Promotions) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a promotions can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

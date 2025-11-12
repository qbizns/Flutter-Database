package loyalty_tier_benefit

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/loyalty_tier_benefit"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/loyalty_tier_benefit"
	"go.uber.org/zap"
)

// Service handles business logic for LoyaltyTierBenefits
type Service struct {
	repo   *loyalty_tier_benefit.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new LoyaltyTierBenefits service
func NewService(repo *loyalty_tier_benefit.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new loyalty_tier_benefits
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateLoyaltyTierBenefitsRequest) (*dto.LoyaltyTierBenefitsResponse, error) {
	s.logger.Info("creating loyalty_tier_benefits",
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
	entity := &loyalty_tier_benefit.LoyaltyTierBenefits{
		OrganizationID: orgID,
		
		TierId: req.TierId,
		
		BenefitCode: req.BenefitCode,
		
		BenefitName: req.BenefitName,
		
		BenefitDescription: req.BenefitDescription,
		
		BenefitType: req.BenefitType,
		
		DiscountValue: req.DiscountValue,
		
		DiscountType: req.DiscountType,
		
		IsActive: req.IsActive,
		
		SortOrder: req.SortOrder,
		
		Icon: req.Icon,
		
		TermsAndConditions: req.TermsAndConditions,
		
		Metadata: req.Metadata,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create loyalty_tier_benefits: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created loyalty_tier_benefits",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a loyalty_tier_benefits by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.LoyaltyTierBenefitsResponse, error) {
	s.logger.Debug("getting loyalty_tier_benefits",
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
		return nil, fmt.Errorf("failed to get loyalty_tier_benefits: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_tier_benefits not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of loyalty_tier_benefits records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.LoyaltyTierBenefitsListResponse, error) {
	s.logger.Debug("listing loyalty_tier_benefits",
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
		return nil, fmt.Errorf("failed to list loyalty_tier_benefits: %w", err)
	}

	// Convert to response
	items := make([]*dto.LoyaltyTierBenefitsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.LoyaltyTierBenefitsListResponse{
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

// Update updates an existing loyalty_tier_benefits
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateLoyaltyTierBenefitsRequest) (*dto.LoyaltyTierBenefitsResponse, error) {
	s.logger.Info("updating loyalty_tier_benefits",
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
		return nil, fmt.Errorf("failed to get loyalty_tier_benefits: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_tier_benefits not found or access denied")
	}
	

	// Update fields
	
	if req.TierId != nil {
		entity.TierId = *req.TierId
	}
	
	if req.BenefitCode != nil {
		entity.BenefitCode = *req.BenefitCode
	}
	
	if req.BenefitName != nil {
		entity.BenefitName = *req.BenefitName
	}
	
	if req.BenefitDescription != nil {
		entity.BenefitDescription = *req.BenefitDescription
	}
	
	if req.BenefitType != nil {
		entity.BenefitType = *req.BenefitType
	}
	
	if req.DiscountValue != nil {
		entity.DiscountValue = *req.DiscountValue
	}
	
	if req.DiscountType != nil {
		entity.DiscountType = *req.DiscountType
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}
	
	if req.Icon != nil {
		entity.Icon = *req.Icon
	}
	
	if req.TermsAndConditions != nil {
		entity.TermsAndConditions = *req.TermsAndConditions
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
		return nil, fmt.Errorf("failed to update loyalty_tier_benefits: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated loyalty_tier_benefits",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a loyalty_tier_benefits
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting loyalty_tier_benefits",
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
		return fmt.Errorf("failed to get loyalty_tier_benefits: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("loyalty_tier_benefits not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete loyalty_tier_benefits: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted loyalty_tier_benefits",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *loyalty_tier_benefit.LoyaltyTierBenefits) *dto.LoyaltyTierBenefitsResponse {
	return &dto.LoyaltyTierBenefitsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		TierId: entity.TierId,
		
		BenefitCode: entity.BenefitCode,
		
		BenefitName: entity.BenefitName,
		
		BenefitDescription: entity.BenefitDescription,
		
		BenefitType: entity.BenefitType,
		
		DiscountValue: entity.DiscountValue,
		
		DiscountType: entity.DiscountType,
		
		IsActive: entity.IsActive,
		
		SortOrder: entity.SortOrder,
		
		Icon: entity.Icon,
		
		TermsAndConditions: entity.TermsAndConditions,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
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


// validateBusinessRules validates business rules for loyalty_tier_benefits
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *loyalty_tier_benefit.LoyaltyTierBenefits) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a loyalty_tier_benefits can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

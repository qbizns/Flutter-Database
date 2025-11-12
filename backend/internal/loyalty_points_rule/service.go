package loyalty_points_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/loyalty_points_rule"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/loyalty_points_rule"
	"go.uber.org/zap"
)

// Service handles business logic for LoyaltyPointsRules
type Service struct {
	repo   *loyalty_points_rule.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new LoyaltyPointsRules service
func NewService(repo *loyalty_points_rule.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new loyalty_points_rules
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateLoyaltyPointsRulesRequest) (*dto.LoyaltyPointsRulesResponse, error) {
	s.logger.Info("creating loyalty_points_rules",
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
	entity := &loyalty_points_rule.LoyaltyPointsRules{
		OrganizationID: orgID,
		
		RuleCode: req.RuleCode,
		
		RuleName: req.RuleName,
		
		Description: req.Description,
		
		RuleType: req.RuleType,
		
		PointsPerAmount: req.PointsPerAmount,
		
		FixedPoints: req.FixedPoints,
		
		Multiplier: req.Multiplier,
		
		AppliesTo: req.AppliesTo,
		
		ApplicableProductIds: req.ApplicableProductIds,
		
		ApplicableCategoryIds: req.ApplicableCategoryIds,
		
		ApplicableTierIds: req.ApplicableTierIds,
		
		MinimumPurchaseAmount: req.MinimumPurchaseAmount,
		
		MaximumPointsPerTransaction: req.MaximumPointsPerTransaction,
		
		MaximumPointsPerDay: req.MaximumPointsPerDay,
		
		MaximumPointsPerMonth: req.MaximumPointsPerMonth,
		
		StartDate: req.StartDate,
		
		EndDate: req.EndDate,
		
		IsActive: req.IsActive,
		
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
		return nil, fmt.Errorf("failed to create loyalty_points_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created loyalty_points_rules",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a loyalty_points_rules by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.LoyaltyPointsRulesResponse, error) {
	s.logger.Debug("getting loyalty_points_rules",
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
		return nil, fmt.Errorf("failed to get loyalty_points_rules: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_points_rules not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of loyalty_points_rules records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.LoyaltyPointsRulesListResponse, error) {
	s.logger.Debug("listing loyalty_points_rules",
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
		return nil, fmt.Errorf("failed to list loyalty_points_rules: %w", err)
	}

	// Convert to response
	items := make([]*dto.LoyaltyPointsRulesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.LoyaltyPointsRulesListResponse{
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

// Update updates an existing loyalty_points_rules
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateLoyaltyPointsRulesRequest) (*dto.LoyaltyPointsRulesResponse, error) {
	s.logger.Info("updating loyalty_points_rules",
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
		return nil, fmt.Errorf("failed to get loyalty_points_rules: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_points_rules not found or access denied")
	}
	

	// Update fields
	
	if req.RuleCode != nil {
		entity.RuleCode = *req.RuleCode
	}
	
	if req.RuleName != nil {
		entity.RuleName = *req.RuleName
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.RuleType != nil {
		entity.RuleType = *req.RuleType
	}
	
	if req.PointsPerAmount != nil {
		entity.PointsPerAmount = *req.PointsPerAmount
	}
	
	if req.FixedPoints != nil {
		entity.FixedPoints = *req.FixedPoints
	}
	
	if req.Multiplier != nil {
		entity.Multiplier = *req.Multiplier
	}
	
	if req.AppliesTo != nil {
		entity.AppliesTo = *req.AppliesTo
	}
	
	if req.ApplicableProductIds != nil {
		entity.ApplicableProductIds = *req.ApplicableProductIds
	}
	
	if req.ApplicableCategoryIds != nil {
		entity.ApplicableCategoryIds = *req.ApplicableCategoryIds
	}
	
	if req.ApplicableTierIds != nil {
		entity.ApplicableTierIds = *req.ApplicableTierIds
	}
	
	if req.MinimumPurchaseAmount != nil {
		entity.MinimumPurchaseAmount = *req.MinimumPurchaseAmount
	}
	
	if req.MaximumPointsPerTransaction != nil {
		entity.MaximumPointsPerTransaction = *req.MaximumPointsPerTransaction
	}
	
	if req.MaximumPointsPerDay != nil {
		entity.MaximumPointsPerDay = *req.MaximumPointsPerDay
	}
	
	if req.MaximumPointsPerMonth != nil {
		entity.MaximumPointsPerMonth = *req.MaximumPointsPerMonth
	}
	
	if req.StartDate != nil {
		entity.StartDate = *req.StartDate
	}
	
	if req.EndDate != nil {
		entity.EndDate = *req.EndDate
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.Priority != nil {
		entity.Priority = *req.Priority
	}
	
	if req.TermsAndConditions != nil {
		entity.TermsAndConditions = *req.TermsAndConditions
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
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update loyalty_points_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated loyalty_points_rules",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a loyalty_points_rules
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting loyalty_points_rules",
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
		return fmt.Errorf("failed to get loyalty_points_rules: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("loyalty_points_rules not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete loyalty_points_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted loyalty_points_rules",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *loyalty_points_rule.LoyaltyPointsRules) *dto.LoyaltyPointsRulesResponse {
	return &dto.LoyaltyPointsRulesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		RuleCode: entity.RuleCode,
		
		RuleName: entity.RuleName,
		
		Description: entity.Description,
		
		RuleType: entity.RuleType,
		
		PointsPerAmount: entity.PointsPerAmount,
		
		FixedPoints: entity.FixedPoints,
		
		Multiplier: entity.Multiplier,
		
		AppliesTo: entity.AppliesTo,
		
		ApplicableProductIds: entity.ApplicableProductIds,
		
		ApplicableCategoryIds: entity.ApplicableCategoryIds,
		
		ApplicableTierIds: entity.ApplicableTierIds,
		
		MinimumPurchaseAmount: entity.MinimumPurchaseAmount,
		
		MaximumPointsPerTransaction: entity.MaximumPointsPerTransaction,
		
		MaximumPointsPerDay: entity.MaximumPointsPerDay,
		
		MaximumPointsPerMonth: entity.MaximumPointsPerMonth,
		
		StartDate: entity.StartDate,
		
		EndDate: entity.EndDate,
		
		IsActive: entity.IsActive,
		
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


// validateBusinessRules validates business rules for loyalty_points_rules
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *loyalty_points_rule.LoyaltyPointsRules) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a loyalty_points_rules can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

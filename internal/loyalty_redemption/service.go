package loyalty_redemption

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/loyalty_redemption"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/loyalty_redemption"
	"go.uber.org/zap"
)

// Service handles business logic for LoyaltyRedemptions
type Service struct {
	repo   *loyalty_redemption.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new LoyaltyRedemptions service
func NewService(repo *loyalty_redemption.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new loyalty_redemptions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateLoyaltyRedemptionsRequest) (*dto.LoyaltyRedemptionsResponse, error) {
	s.logger.Info("creating loyalty_redemptions",
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
	entity := &loyalty_redemption.LoyaltyRedemptions{
		OrganizationID: orgID,
		
		CustomerId: req.CustomerId,
		
		RewardId: req.RewardId,
		
		RedemptionNumber: req.RedemptionNumber,
		
		RedemptionDate: req.RedemptionDate,
		
		PointsRedeemed: req.PointsRedeemed,
		
		Status: req.Status,
		
		SaleId: req.SaleId,
		
		UsedDate: req.UsedDate,
		
		ExpiryDate: req.ExpiryDate,
		
		FulfillmentStatus: req.FulfillmentStatus,
		
		FulfillmentNotes: req.FulfillmentNotes,
		
		FulfilledBy: req.FulfilledBy,
		
		FulfilledAt: req.FulfilledAt,
		
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
		return nil, fmt.Errorf("failed to create loyalty_redemptions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created loyalty_redemptions",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a loyalty_redemptions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.LoyaltyRedemptionsResponse, error) {
	s.logger.Debug("getting loyalty_redemptions",
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
		return nil, fmt.Errorf("failed to get loyalty_redemptions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_redemptions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of loyalty_redemptions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.LoyaltyRedemptionsListResponse, error) {
	s.logger.Debug("listing loyalty_redemptions",
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
		return nil, fmt.Errorf("failed to list loyalty_redemptions: %w", err)
	}

	// Convert to response
	items := make([]*dto.LoyaltyRedemptionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.LoyaltyRedemptionsListResponse{
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

// Update updates an existing loyalty_redemptions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateLoyaltyRedemptionsRequest) (*dto.LoyaltyRedemptionsResponse, error) {
	s.logger.Info("updating loyalty_redemptions",
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
		return nil, fmt.Errorf("failed to get loyalty_redemptions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("loyalty_redemptions not found or access denied")
	}
	

	// Update fields
	
	if req.CustomerId != nil {
		entity.CustomerId = *req.CustomerId
	}
	
	if req.RewardId != nil {
		entity.RewardId = *req.RewardId
	}
	
	if req.RedemptionNumber != nil {
		entity.RedemptionNumber = *req.RedemptionNumber
	}
	
	if req.RedemptionDate != nil {
		entity.RedemptionDate = *req.RedemptionDate
	}
	
	if req.PointsRedeemed != nil {
		entity.PointsRedeemed = *req.PointsRedeemed
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.SaleId != nil {
		entity.SaleId = *req.SaleId
	}
	
	if req.UsedDate != nil {
		entity.UsedDate = *req.UsedDate
	}
	
	if req.ExpiryDate != nil {
		entity.ExpiryDate = *req.ExpiryDate
	}
	
	if req.FulfillmentStatus != nil {
		entity.FulfillmentStatus = *req.FulfillmentStatus
	}
	
	if req.FulfillmentNotes != nil {
		entity.FulfillmentNotes = *req.FulfillmentNotes
	}
	
	if req.FulfilledBy != nil {
		entity.FulfilledBy = *req.FulfilledBy
	}
	
	if req.FulfilledAt != nil {
		entity.FulfilledAt = *req.FulfilledAt
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
		return nil, fmt.Errorf("failed to update loyalty_redemptions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated loyalty_redemptions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a loyalty_redemptions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting loyalty_redemptions",
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
		return fmt.Errorf("failed to get loyalty_redemptions: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("loyalty_redemptions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete loyalty_redemptions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted loyalty_redemptions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *loyalty_redemption.LoyaltyRedemptions) *dto.LoyaltyRedemptionsResponse {
	return &dto.LoyaltyRedemptionsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		CustomerId: entity.CustomerId,
		
		RewardId: entity.RewardId,
		
		RedemptionNumber: entity.RedemptionNumber,
		
		RedemptionDate: entity.RedemptionDate,
		
		PointsRedeemed: entity.PointsRedeemed,
		
		Status: entity.Status,
		
		SaleId: entity.SaleId,
		
		UsedDate: entity.UsedDate,
		
		ExpiryDate: entity.ExpiryDate,
		
		FulfillmentStatus: entity.FulfillmentStatus,
		
		FulfillmentNotes: entity.FulfillmentNotes,
		
		FulfilledBy: entity.FulfilledBy,
		
		FulfilledAt: entity.FulfilledAt,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
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


// validateBusinessRules validates business rules for loyalty_redemptions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *loyalty_redemption.LoyaltyRedemptions) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a loyalty_redemptions can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

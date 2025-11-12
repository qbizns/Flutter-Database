package tip_distribution

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/tip_distribution"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/tip_distribution"
	"go.uber.org/zap"
)

// Service handles business logic for TipDistributions
type Service struct {
	repo   *tip_distribution.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new TipDistributions service
func NewService(repo *tip_distribution.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new tip_distributions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateTipDistributionsRequest) (*dto.TipDistributionsResponse, error) {
	s.logger.Info("creating tip_distributions",
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
	entity := &tip_distribution.TipDistributions{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		TipPoolId: req.TipPoolId,
		
		DistributionDate: req.DistributionDate,
		
		PeriodStart: req.PeriodStart,
		
		PeriodEnd: req.PeriodEnd,
		
		ShiftId: req.ShiftId,
		
		EmployeeId: req.EmployeeId,
		
		SourceType: req.SourceType,
		
		SourceSaleId: req.SourceSaleId,
		
		SourceOrderId: req.SourceOrderId,
		
		TipAmount: req.TipAmount,
		
		DistributionAmount: req.DistributionAmount,
		
		DistributionPercentage: req.DistributionPercentage,
		
		PaymentStatus: req.PaymentStatus,
		
		PaymentMethod: req.PaymentMethod,
		
		PaidAt: req.PaidAt,
		
		PaidBy: req.PaidBy,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'sale',: req.'sale',,
		
		'pending',: req.'pending',,
		
		TipAmount: req.TipAmount,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create tip_distributions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created tip_distributions",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a tip_distributions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.TipDistributionsResponse, error) {
	s.logger.Debug("getting tip_distributions",
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
		return nil, fmt.Errorf("failed to get tip_distributions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("tip_distributions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of tip_distributions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.TipDistributionsListResponse, error) {
	s.logger.Debug("listing tip_distributions",
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
		return nil, fmt.Errorf("failed to list tip_distributions: %w", err)
	}

	// Convert to response
	items := make([]*dto.TipDistributionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.TipDistributionsListResponse{
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

// Update updates an existing tip_distributions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateTipDistributionsRequest) (*dto.TipDistributionsResponse, error) {
	s.logger.Info("updating tip_distributions",
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
		return nil, fmt.Errorf("failed to get tip_distributions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("tip_distributions not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.TipPoolId != nil {
		entity.TipPoolId = *req.TipPoolId
	}
	
	if req.DistributionDate != nil {
		entity.DistributionDate = *req.DistributionDate
	}
	
	if req.PeriodStart != nil {
		entity.PeriodStart = *req.PeriodStart
	}
	
	if req.PeriodEnd != nil {
		entity.PeriodEnd = *req.PeriodEnd
	}
	
	if req.ShiftId != nil {
		entity.ShiftId = *req.ShiftId
	}
	
	if req.EmployeeId != nil {
		entity.EmployeeId = *req.EmployeeId
	}
	
	if req.SourceType != nil {
		entity.SourceType = *req.SourceType
	}
	
	if req.SourceSaleId != nil {
		entity.SourceSaleId = *req.SourceSaleId
	}
	
	if req.SourceOrderId != nil {
		entity.SourceOrderId = *req.SourceOrderId
	}
	
	if req.TipAmount != nil {
		entity.TipAmount = *req.TipAmount
	}
	
	if req.DistributionAmount != nil {
		entity.DistributionAmount = *req.DistributionAmount
	}
	
	if req.DistributionPercentage != nil {
		entity.DistributionPercentage = *req.DistributionPercentage
	}
	
	if req.PaymentStatus != nil {
		entity.PaymentStatus = *req.PaymentStatus
	}
	
	if req.PaymentMethod != nil {
		entity.PaymentMethod = *req.PaymentMethod
	}
	
	if req.PaidAt != nil {
		entity.PaidAt = *req.PaidAt
	}
	
	if req.PaidBy != nil {
		entity.PaidBy = *req.PaidBy
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
	
	if req.'sale', != nil {
		entity.'sale', = *req.'sale',
	}
	
	if req.'pending', != nil {
		entity.'pending', = *req.'pending',
	}
	
	if req.TipAmount != nil {
		entity.TipAmount = *req.TipAmount
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update tip_distributions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated tip_distributions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a tip_distributions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting tip_distributions",
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
		return fmt.Errorf("failed to get tip_distributions: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("tip_distributions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete tip_distributions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted tip_distributions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *tip_distribution.TipDistributions) *dto.TipDistributionsResponse {
	return &dto.TipDistributionsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		TipPoolId: entity.TipPoolId,
		
		DistributionDate: entity.DistributionDate,
		
		PeriodStart: entity.PeriodStart,
		
		PeriodEnd: entity.PeriodEnd,
		
		ShiftId: entity.ShiftId,
		
		EmployeeId: entity.EmployeeId,
		
		SourceType: entity.SourceType,
		
		SourceSaleId: entity.SourceSaleId,
		
		SourceOrderId: entity.SourceOrderId,
		
		TipAmount: entity.TipAmount,
		
		DistributionAmount: entity.DistributionAmount,
		
		DistributionPercentage: entity.DistributionPercentage,
		
		PaymentStatus: entity.PaymentStatus,
		
		PaymentMethod: entity.PaymentMethod,
		
		PaidAt: entity.PaidAt,
		
		PaidBy: entity.PaidBy,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'sale',: entity.'sale',,
		
		'pending',: entity.'pending',,
		
		TipAmount: entity.TipAmount,
		
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


// validateBusinessRules validates business rules for tip_distributions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *tip_distribution.TipDistributions) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a tip_distributions can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

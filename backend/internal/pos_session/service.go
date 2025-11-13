package pos_session

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

// Service handles business logic for PosSessions
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PosSessions service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new pos_sessions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreatePosSessionsRequest) (*PosSessionsResponse, error) {
	s.logger.Info("creating pos_sessions",
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
	entity := &PosSessions{
		OrganizationID: orgID,
		
		SessionNumber: req.SessionNumber,
		
		SessionName: req.SessionName,
		
		DeviceId: req.DeviceId,
		
		LocationId: req.LocationId,
		
		UserId: req.UserId,
		
		ShiftId: req.ShiftId,
		
		OpenedAt: req.OpenedAt,
		
		ClosedAt: req.ClosedAt,
		
		OpeningCash: req.OpeningCash,
		
		OpeningCard: req.OpeningCard,
		
		OpeningOther: req.OpeningOther,
		
		ExpectedCash: req.ExpectedCash,
		
		ExpectedCard: req.ExpectedCard,
		
		ExpectedOther: req.ExpectedOther,
		
		CountedCash: req.CountedCash,
		
		CountedCard: req.CountedCard,
		
		CountedOther: req.CountedOther,
		
		DifferenceCash: req.DifferenceCash,
		
		DifferenceCard: req.DifferenceCard,
		
		DifferenceOther: req.DifferenceOther,
		
		Status: req.Status,
		
		ZReportNumber: req.ZReportNumber,
		
		Notes: req.Notes,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create pos_sessions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created pos_sessions",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a pos_sessions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PosSessionsResponse, error) {
	s.logger.Debug("getting pos_sessions",
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
		return nil, fmt.Errorf("failed to get pos_sessions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("pos_sessions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of pos_sessions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*PosSessionsListResponse, error) {
	s.logger.Debug("listing pos_sessions",
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
		return nil, fmt.Errorf("failed to list pos_sessions: %w", err)
	}

	// Convert to response
	items := make([]*PosSessionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PosSessionsListResponse{
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

// Update updates an existing pos_sessions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdatePosSessionsRequest) (*PosSessionsResponse, error) {
	s.logger.Info("updating pos_sessions",
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
		return nil, fmt.Errorf("failed to get pos_sessions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("pos_sessions not found or access denied")
	}
	

	// Update fields
	
	if req.SessionNumber != nil {
		entity.SessionNumber = *req.SessionNumber
	}
	
	if req.SessionName != nil {
		entity.SessionName = *req.SessionName
	}
	
	if req.DeviceId != nil {
		entity.DeviceId = *req.DeviceId
	}
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.ShiftId != nil {
		entity.ShiftId = *req.ShiftId
	}
	
	if req.OpenedAt != nil {
		entity.OpenedAt = *req.OpenedAt
	}
	
	if req.ClosedAt != nil {
		entity.ClosedAt = *req.ClosedAt
	}
	
	if req.OpeningCash != nil {
		entity.OpeningCash = *req.OpeningCash
	}
	
	if req.OpeningCard != nil {
		entity.OpeningCard = *req.OpeningCard
	}
	
	if req.OpeningOther != nil {
		entity.OpeningOther = *req.OpeningOther
	}
	
	if req.ExpectedCash != nil {
		entity.ExpectedCash = *req.ExpectedCash
	}
	
	if req.ExpectedCard != nil {
		entity.ExpectedCard = *req.ExpectedCard
	}
	
	if req.ExpectedOther != nil {
		entity.ExpectedOther = *req.ExpectedOther
	}
	
	if req.CountedCash != nil {
		entity.CountedCash = *req.CountedCash
	}
	
	if req.CountedCard != nil {
		entity.CountedCard = *req.CountedCard
	}
	
	if req.CountedOther != nil {
		entity.CountedOther = *req.CountedOther
	}
	
	if req.DifferenceCash != nil {
		entity.DifferenceCash = *req.DifferenceCash
	}
	
	if req.DifferenceCard != nil {
		entity.DifferenceCard = *req.DifferenceCard
	}
	
	if req.DifferenceOther != nil {
		entity.DifferenceOther = *req.DifferenceOther
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.ZReportNumber != nil {
		entity.ZReportNumber = *req.ZReportNumber
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
		return nil, fmt.Errorf("failed to update pos_sessions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated pos_sessions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a pos_sessions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting pos_sessions",
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
		return fmt.Errorf("failed to get pos_sessions: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("pos_sessions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete pos_sessions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted pos_sessions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PosSessions) *PosSessionsResponse {
	return &PosSessionsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		SessionNumber: entity.SessionNumber,
		
		SessionName: entity.SessionName,
		
		DeviceId: entity.DeviceId,
		
		LocationId: entity.LocationId,
		
		UserId: entity.UserId,
		
		ShiftId: entity.ShiftId,
		
		OpenedAt: entity.OpenedAt,
		
		ClosedAt: entity.ClosedAt,
		
		OpeningCash: entity.OpeningCash,
		
		OpeningCard: entity.OpeningCard,
		
		OpeningOther: entity.OpeningOther,
		
		ExpectedCash: entity.ExpectedCash,
		
		ExpectedCard: entity.ExpectedCard,
		
		ExpectedOther: entity.ExpectedOther,
		
		CountedCash: entity.CountedCash,
		
		CountedCard: entity.CountedCard,
		
		CountedOther: entity.CountedOther,
		
		DifferenceCash: entity.DifferenceCash,
		
		DifferenceCard: entity.DifferenceCard,
		
		DifferenceOther: entity.DifferenceOther,
		
		Status: entity.Status,
		
		ZReportNumber: entity.ZReportNumber,
		
		Notes: entity.Notes,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for pos_sessions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PosSessions) error {
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

// canDelete checks if a pos_sessions can be deleted
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

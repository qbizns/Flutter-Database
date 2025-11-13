package shift

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/pos-backend/internal/logging"

	"go.uber.org/zap"
)

// Service handles business logic for Shifts
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Shifts service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new shifts
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateShiftsRequest) (*ShiftsResponse, error) {
	s.logger.Info("creating shifts",
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
	entity := &Shifts{
		OrganizationId: orgID,

		LocationId: req.LocationId,

		UserId: req.UserId,

		ShiftNumber: req.ShiftNumber,

		StartTime: req.StartTime,

		EndTime: req.EndTime,

		Status: req.Status,

		OpeningCash: req.OpeningCash,

		OpeningNotes: req.OpeningNotes,

		ExpectedCash: req.ExpectedCash,

		ActualCash: req.ActualCash,

		CashDifference: req.CashDifference,

		ClosingNotes: req.ClosingNotes,

		TotalSales: req.TotalSales,

		TotalTransactions: req.TotalTransactions,

		TotalRefunds: req.TotalRefunds,

		TotalDiscounts: req.TotalDiscounts,

		PaymentBreakdown: req.PaymentBreakdown,

		Notes: req.Notes,

		Metadata: req.Metadata,

		ClosedBy: req.ClosedBy,

		ClosedAt: req.ClosedAt,

		CreatedBy: req.CreatedBy,

		UpdatedBy: req.UpdatedBy,

	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create shifts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created shifts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a shifts by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ShiftsResponse, error) {
	s.logger.Debug("getting shifts",
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
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("shifts not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of shifts records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ShiftsListResponse, error) {
	s.logger.Debug("listing shifts",
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
		return nil, fmt.Errorf("failed to list shifts: %w", err)
	}

	// Convert to response
	items := make([]*ShiftsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ShiftsListResponse{
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

// Update updates an existing shifts
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateShiftsRequest) (*ShiftsResponse, error) {
	s.logger.Info("updating shifts",
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
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("shifts not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = req.LocationId
	}
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}

	if req.ShiftNumber != nil {
		entity.ShiftNumber = *req.ShiftNumber
	}

	if req.StartTime != nil {
		entity.StartTime = *req.StartTime
	}
	
	if req.EndTime != nil {
		entity.EndTime = req.EndTime
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.OpeningCash != nil {
		entity.OpeningCash = req.OpeningCash
	}
	
	if req.OpeningNotes != nil {
		entity.OpeningNotes = req.OpeningNotes
	}
	
	if req.ExpectedCash != nil {
		entity.ExpectedCash = req.ExpectedCash
	}
	
	if req.ActualCash != nil {
		entity.ActualCash = req.ActualCash
	}
	
	if req.CashDifference != nil {
		entity.CashDifference = req.CashDifference
	}
	
	if req.ClosingNotes != nil {
		entity.ClosingNotes = req.ClosingNotes
	}
	
	if req.TotalSales != nil {
		entity.TotalSales = req.TotalSales
	}
	
	if req.TotalTransactions != nil {
		entity.TotalTransactions = req.TotalTransactions
	}
	
	if req.TotalRefunds != nil {
		entity.TotalRefunds = req.TotalRefunds
	}
	
	if req.TotalDiscounts != nil {
		entity.TotalDiscounts = req.TotalDiscounts
	}
	
	if req.PaymentBreakdown != nil {
		entity.PaymentBreakdown = *req.PaymentBreakdown
	}

	if req.Notes != nil {
		entity.Notes = req.Notes
	}

	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.ClosedBy != nil {
		entity.ClosedBy = req.ClosedBy
	}
	
	if req.ClosedAt != nil {
		entity.ClosedAt = req.ClosedAt
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
		return nil, fmt.Errorf("failed to update shifts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated shifts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a shifts
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting shifts",
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
		return fmt.Errorf("failed to get shifts: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("shifts not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete shifts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted shifts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Shifts) *ShiftsResponse {
	return &ShiftsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		UserId: entity.UserId,
		
		ShiftNumber: entity.ShiftNumber,
		
		StartTime: entity.StartTime,
		
		EndTime: entity.EndTime,
		
		Status: entity.Status,
		
		OpeningCash: entity.OpeningCash,
		
		OpeningNotes: entity.OpeningNotes,
		
		ExpectedCash: entity.ExpectedCash,
		
		ActualCash: entity.ActualCash,
		
		CashDifference: entity.CashDifference,
		
		ClosingNotes: entity.ClosingNotes,
		
		TotalSales: entity.TotalSales,
		
		TotalTransactions: entity.TotalTransactions,
		
		TotalRefunds: entity.TotalRefunds,
		
		TotalDiscounts: entity.TotalDiscounts,
		
		PaymentBreakdown: entity.PaymentBreakdown,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		ClosedBy: entity.ClosedBy,

		ClosedAt: entity.ClosedAt,

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


// validateBusinessRules validates business rules for shifts
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Shifts) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a shifts can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

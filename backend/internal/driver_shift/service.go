package driver_shift

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for DriverShifts
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new DriverShifts service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new driver_shifts
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateDriverShiftsRequest) (*DriverShiftsResponse, error) {
	s.logger.Info("creating driver_shifts",
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
	entity := &DriverShifts{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		DriverId: req.DriverId,
		
		ShiftDate: req.ShiftDate,
		
		ScheduledStartTime: req.ScheduledStartTime,
		
		ScheduledEndTime: req.ScheduledEndTime,
		
		ActualStartTime: req.ActualStartTime,
		
		ActualEndTime: req.ActualEndTime,
		
		Status: req.Status,
		
		TotalBreakMinutes: req.TotalBreakMinutes,
		
		TotalDeliveries: req.TotalDeliveries,
		
		TotalDistanceKm: req.TotalDistanceKm,
		
		TotalEarnings: req.TotalEarnings,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'scheduled',: req.'scheduled',,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create driver_shifts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created driver_shifts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a driver_shifts by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DriverShiftsResponse, error) {
	s.logger.Debug("getting driver_shifts",
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
		return nil, fmt.Errorf("failed to get driver_shifts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("driver_shifts not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of driver_shifts records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*DriverShiftsListResponse, error) {
	s.logger.Debug("listing driver_shifts",
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
		return nil, fmt.Errorf("failed to list driver_shifts: %w", err)
	}

	// Convert to response
	items := make([]*DriverShiftsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &DriverShiftsListResponse{
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

// Update updates an existing driver_shifts
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateDriverShiftsRequest) (*DriverShiftsResponse, error) {
	s.logger.Info("updating driver_shifts",
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
		return nil, fmt.Errorf("failed to get driver_shifts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("driver_shifts not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = req.LocationId
	}
	
	if req.DriverId != nil {
		entity.DriverId = req.DriverId
	}
	
	if req.ShiftDate != nil {
		entity.ShiftDate = req.ShiftDate
	}
	
	if req.ScheduledStartTime != nil {
		entity.ScheduledStartTime = req.ScheduledStartTime
	}
	
	if req.ScheduledEndTime != nil {
		entity.ScheduledEndTime = req.ScheduledEndTime
	}
	
	if req.ActualStartTime != nil {
		entity.ActualStartTime = req.ActualStartTime
	}
	
	if req.ActualEndTime != nil {
		entity.ActualEndTime = req.ActualEndTime
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.TotalBreakMinutes != nil {
		entity.TotalBreakMinutes = req.TotalBreakMinutes
	}
	
	if req.TotalDeliveries != nil {
		entity.TotalDeliveries = req.TotalDeliveries
	}
	
	if req.TotalDistanceKm != nil {
		entity.TotalDistanceKm = req.TotalDistanceKm
	}
	
	if req.TotalEarnings != nil {
		entity.TotalEarnings = req.TotalEarnings
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
	
	if req.'scheduled', != nil {
		entity.'scheduled', = *req.'scheduled',
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update driver_shifts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated driver_shifts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a driver_shifts
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting driver_shifts",
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
		return fmt.Errorf("failed to get driver_shifts: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("driver_shifts not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete driver_shifts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted driver_shifts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *DriverShifts) *DriverShiftsResponse {
	return &DriverShiftsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		DriverId: entity.DriverId,
		
		ShiftDate: entity.ShiftDate,
		
		ScheduledStartTime: entity.ScheduledStartTime,
		
		ScheduledEndTime: entity.ScheduledEndTime,
		
		ActualStartTime: entity.ActualStartTime,
		
		ActualEndTime: entity.ActualEndTime,
		
		Status: entity.Status,
		
		TotalBreakMinutes: entity.TotalBreakMinutes,
		
		TotalDeliveries: entity.TotalDeliveries,
		
		TotalDistanceKm: entity.TotalDistanceKm,
		
		TotalEarnings: entity.TotalEarnings,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'scheduled',: entity.'scheduled',,
		
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


// validateBusinessRules validates business rules for driver_shifts
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *DriverShifts) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a driver_shifts can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

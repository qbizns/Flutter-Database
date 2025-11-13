package time_clock_entry

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for TimeClockEntries
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new TimeClockEntries service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new time_clock_entries
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateTimeClockEntriesRequest) (*TimeClockEntriesResponse, error) {
	s.logger.Info("creating time_clock_entries",
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
	entity := &TimeClockEntries{
		OrganizationId: orgID,
		
		LocationId: req.LocationId,
		
		EmployeeId: req.EmployeeId,
		
		ScheduleId: req.ScheduleId,
		
		EntryType: req.EntryType,
		
		EntryTimestamp: req.EntryTimestamp,
		
		ScheduledTimestamp: req.ScheduledTimestamp,
		
		DeviceId: req.DeviceId,
		
		GpsLocation: req.GpsLocation,
		
		IpAddress: req.IpAddress,
		
		IsLate: req.IsLate,
		
		IsEarly: req.IsEarly,
		
		VarianceMinutes: req.VarianceMinutes,
		
		RequiresApproval: req.RequiresApproval,
		
		ApprovedBy: req.ApprovedBy,
		
		ApprovedAt: req.ApprovedAt,
		
		IsManualEntry: req.IsManualEntry,
		
		CorrectionNotes: req.CorrectionNotes,
		
		PhotoUrl: req.PhotoUrl,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'clockIn',: req.'clockIn',,
		
		'mealStart',: req.'mealStart',,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create time_clock_entries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created time_clock_entries",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a time_clock_entries by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*TimeClockEntriesResponse, error) {
	s.logger.Debug("getting time_clock_entries",
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
		return nil, fmt.Errorf("failed to get time_clock_entries: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("time_clock_entries not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of time_clock_entries records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*TimeClockEntriesListResponse, error) {
	s.logger.Debug("listing time_clock_entries",
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
		return nil, fmt.Errorf("failed to list time_clock_entries: %w", err)
	}

	// Convert to response
	items := make([]*TimeClockEntriesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &TimeClockEntriesListResponse{
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

// Update updates an existing time_clock_entries
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateTimeClockEntriesRequest) (*TimeClockEntriesResponse, error) {
	s.logger.Info("updating time_clock_entries",
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
		return nil, fmt.Errorf("failed to get time_clock_entries: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("time_clock_entries not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.EmployeeId != nil {
		entity.EmployeeId = *req.EmployeeId
	}
	
	if req.ScheduleId != nil {
		entity.ScheduleId = *req.ScheduleId
	}
	
	if req.EntryType != nil {
		entity.EntryType = *req.EntryType
	}
	
	if req.EntryTimestamp != nil {
		entity.EntryTimestamp = *req.EntryTimestamp
	}
	
	if req.ScheduledTimestamp != nil {
		entity.ScheduledTimestamp = *req.ScheduledTimestamp
	}
	
	if req.DeviceId != nil {
		entity.DeviceId = *req.DeviceId
	}
	
	if req.GpsLocation != nil {
		entity.GpsLocation = *req.GpsLocation
	}
	
	if req.IpAddress != nil {
		entity.IpAddress = req.IpAddress
	}
	
	if req.IsLate != nil {
		entity.IsLate = req.IsLate
	}
	
	if req.IsEarly != nil {
		entity.IsEarly = req.IsEarly
	}
	
	if req.VarianceMinutes != nil {
		entity.VarianceMinutes = *req.VarianceMinutes
	}
	
	if req.RequiresApproval != nil {
		entity.RequiresApproval = *req.RequiresApproval
	}
	
	if req.ApprovedBy != nil {
		entity.ApprovedBy = *req.ApprovedBy
	}
	
	if req.ApprovedAt != nil {
		entity.ApprovedAt = req.ApprovedAt
	}
	
	if req.IsManualEntry != nil {
		entity.IsManualEntry = req.IsManualEntry
	}
	
	if req.CorrectionNotes != nil {
		entity.CorrectionNotes = *req.CorrectionNotes
	}
	
	if req.PhotoUrl != nil {
		entity.PhotoUrl = *req.PhotoUrl
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
	
	if req.'clockIn', != nil {
		entity.'clockIn', = *req.'clockIn',
	}
	
	if req.'mealStart', != nil {
		entity.'mealStart', = *req.'mealStart',
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update time_clock_entries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated time_clock_entries",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a time_clock_entries
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting time_clock_entries",
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
		return fmt.Errorf("failed to get time_clock_entries: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("time_clock_entries not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete time_clock_entries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted time_clock_entries",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *TimeClockEntries) *TimeClockEntriesResponse {
	return &TimeClockEntriesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		EmployeeId: entity.EmployeeId,
		
		ScheduleId: entity.ScheduleId,
		
		EntryType: entity.EntryType,
		
		EntryTimestamp: entity.EntryTimestamp,
		
		ScheduledTimestamp: entity.ScheduledTimestamp,
		
		DeviceId: entity.DeviceId,
		
		GpsLocation: entity.GpsLocation,
		
		IpAddress: entity.IpAddress,
		
		IsLate: entity.IsLate,
		
		IsEarly: entity.IsEarly,
		
		VarianceMinutes: entity.VarianceMinutes,
		
		RequiresApproval: entity.RequiresApproval,
		
		ApprovedBy: entity.ApprovedBy,
		
		ApprovedAt: entity.ApprovedAt,
		
		IsManualEntry: entity.IsManualEntry,
		
		CorrectionNotes: entity.CorrectionNotes,
		
		PhotoUrl: entity.PhotoUrl,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'clockIn',: entity.'clockIn',,
		
		'mealStart',: entity.'mealStart',,
		
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


// validateBusinessRules validates business rules for time_clock_entries
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *TimeClockEntries) error {
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

// canDelete checks if a time_clock_entries can be deleted
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

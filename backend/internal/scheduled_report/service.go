package scheduled_report

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

// Service handles business logic for ScheduledReports
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new ScheduledReports service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new scheduled_reports
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateScheduledReportsRequest) (*ScheduledReportsResponse, error) {
	s.logger.Info("creating scheduled_reports",
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
	entity := &ScheduledReports{
		OrganizationID: orgID,
		
		ReportName: req.ReportName,
		
		ReportType: req.ReportType,
		
		ScheduleFrequency: req.ScheduleFrequency,
		
		ScheduleDayOfWeek: req.ScheduleDayOfWeek,
		
		ScheduleDayOfMonth: req.ScheduleDayOfMonth,
		
		ScheduleTime: req.ScheduleTime,
		
		ScheduleTimezone: req.ScheduleTimezone,
		
		ReportParameters: req.ReportParameters,
		
		DeliveryMethod: req.DeliveryMethod,
		
		DeliveryRecipients: req.DeliveryRecipients,
		
		OutputFormat: req.OutputFormat,
		
		IsActive: req.IsActive,
		
		LastRunAt: req.LastRunAt,
		
		LastRunStatus: req.LastRunStatus,
		
		NextRunAt: req.NextRunAt,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create scheduled_reports: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created scheduled_reports",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a scheduled_reports by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ScheduledReportsResponse, error) {
	s.logger.Debug("getting scheduled_reports",
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
		return nil, fmt.Errorf("failed to get scheduled_reports: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("scheduled_reports not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of scheduled_reports records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ScheduledReportsListResponse, error) {
	s.logger.Debug("listing scheduled_reports",
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
		return nil, fmt.Errorf("failed to list scheduled_reports: %w", err)
	}

	// Convert to response
	items := make([]*ScheduledReportsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ScheduledReportsListResponse{
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

// Update updates an existing scheduled_reports
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateScheduledReportsRequest) (*ScheduledReportsResponse, error) {
	s.logger.Info("updating scheduled_reports",
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
		return nil, fmt.Errorf("failed to get scheduled_reports: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("scheduled_reports not found or access denied")
	}
	

	// Update fields
	
	if req.ReportName != nil {
		entity.ReportName = *req.ReportName
	}
	
	if req.ReportType != nil {
		entity.ReportType = *req.ReportType
	}
	
	if req.ScheduleFrequency != nil {
		entity.ScheduleFrequency = *req.ScheduleFrequency
	}
	
	if req.ScheduleDayOfWeek != nil {
		entity.ScheduleDayOfWeek = *req.ScheduleDayOfWeek
	}
	
	if req.ScheduleDayOfMonth != nil {
		entity.ScheduleDayOfMonth = *req.ScheduleDayOfMonth
	}
	
	if req.ScheduleTime != nil {
		entity.ScheduleTime = *req.ScheduleTime
	}
	
	if req.ScheduleTimezone != nil {
		entity.ScheduleTimezone = *req.ScheduleTimezone
	}
	
	if req.ReportParameters != nil {
		entity.ReportParameters = *req.ReportParameters
	}
	
	if req.DeliveryMethod != nil {
		entity.DeliveryMethod = *req.DeliveryMethod
	}
	
	if req.DeliveryRecipients != nil {
		entity.DeliveryRecipients = *req.DeliveryRecipients
	}
	
	if req.OutputFormat != nil {
		entity.OutputFormat = *req.OutputFormat
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.LastRunAt != nil {
		entity.LastRunAt = *req.LastRunAt
	}
	
	if req.LastRunStatus != nil {
		entity.LastRunStatus = *req.LastRunStatus
	}
	
	if req.NextRunAt != nil {
		entity.NextRunAt = *req.NextRunAt
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update scheduled_reports: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated scheduled_reports",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a scheduled_reports
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting scheduled_reports",
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
		return fmt.Errorf("failed to get scheduled_reports: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("scheduled_reports not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete scheduled_reports: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted scheduled_reports",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *ScheduledReports) *ScheduledReportsResponse {
	return &ScheduledReportsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ReportName: entity.ReportName,
		
		ReportType: entity.ReportType,
		
		ScheduleFrequency: entity.ScheduleFrequency,
		
		ScheduleDayOfWeek: entity.ScheduleDayOfWeek,
		
		ScheduleDayOfMonth: entity.ScheduleDayOfMonth,
		
		ScheduleTime: entity.ScheduleTime,
		
		ScheduleTimezone: entity.ScheduleTimezone,
		
		ReportParameters: entity.ReportParameters,
		
		DeliveryMethod: entity.DeliveryMethod,
		
		DeliveryRecipients: entity.DeliveryRecipients,
		
		OutputFormat: entity.OutputFormat,
		
		IsActive: entity.IsActive,
		
		LastRunAt: entity.LastRunAt,
		
		LastRunStatus: entity.LastRunStatus,
		
		NextRunAt: entity.NextRunAt,
		
		CreatedBy: entity.CreatedBy,
		
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


// validateBusinessRules validates business rules for scheduled_reports
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *ScheduledReports) error {
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

// canDelete checks if a scheduled_reports can be deleted
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

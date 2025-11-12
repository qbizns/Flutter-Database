package pos_error_log

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/pos_error_log"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/pos_error_log"
	"go.uber.org/zap"
)

// Service handles business logic for PosErrorLogs
type Service struct {
	repo   *pos_error_log.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PosErrorLogs service
func NewService(repo *pos_error_log.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new pos_error_logs
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreatePosErrorLogsRequest) (*dto.PosErrorLogsResponse, error) {
	s.logger.Info("creating pos_error_logs",
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
	entity := &pos_error_log.PosErrorLogs{
		OrganizationID: orgID,
		
		ErrorLevel: req.ErrorLevel,
		
		ErrorCode: req.ErrorCode,
		
		ErrorMessage: req.ErrorMessage,
		
		DeviceId: req.DeviceId,
		
		UserId: req.UserId,
		
		PosSessionId: req.PosSessionId,
		
		SaleId: req.SaleId,
		
		StackTrace: req.StackTrace,
		
		RequestData: req.RequestData,
		
		ErrorData: req.ErrorData,
		
		IsResolved: req.IsResolved,
		
		ResolvedBy: req.ResolvedBy,
		
		ResolvedAt: req.ResolvedAt,
		
		ResolutionNotes: req.ResolutionNotes,
		
		OccurredAt: req.OccurredAt,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create pos_error_logs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created pos_error_logs",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a pos_error_logs by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.PosErrorLogsResponse, error) {
	s.logger.Debug("getting pos_error_logs",
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
		return nil, fmt.Errorf("failed to get pos_error_logs: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("pos_error_logs not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of pos_error_logs records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.PosErrorLogsListResponse, error) {
	s.logger.Debug("listing pos_error_logs",
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
		return nil, fmt.Errorf("failed to list pos_error_logs: %w", err)
	}

	// Convert to response
	items := make([]*dto.PosErrorLogsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.PosErrorLogsListResponse{
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

// Update updates an existing pos_error_logs
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdatePosErrorLogsRequest) (*dto.PosErrorLogsResponse, error) {
	s.logger.Info("updating pos_error_logs",
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
		return nil, fmt.Errorf("failed to get pos_error_logs: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("pos_error_logs not found or access denied")
	}
	

	// Update fields
	
	if req.ErrorLevel != nil {
		entity.ErrorLevel = *req.ErrorLevel
	}
	
	if req.ErrorCode != nil {
		entity.ErrorCode = *req.ErrorCode
	}
	
	if req.ErrorMessage != nil {
		entity.ErrorMessage = *req.ErrorMessage
	}
	
	if req.DeviceId != nil {
		entity.DeviceId = *req.DeviceId
	}
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.PosSessionId != nil {
		entity.PosSessionId = *req.PosSessionId
	}
	
	if req.SaleId != nil {
		entity.SaleId = *req.SaleId
	}
	
	if req.StackTrace != nil {
		entity.StackTrace = *req.StackTrace
	}
	
	if req.RequestData != nil {
		entity.RequestData = *req.RequestData
	}
	
	if req.ErrorData != nil {
		entity.ErrorData = *req.ErrorData
	}
	
	if req.IsResolved != nil {
		entity.IsResolved = *req.IsResolved
	}
	
	if req.ResolvedBy != nil {
		entity.ResolvedBy = *req.ResolvedBy
	}
	
	if req.ResolvedAt != nil {
		entity.ResolvedAt = *req.ResolvedAt
	}
	
	if req.ResolutionNotes != nil {
		entity.ResolutionNotes = *req.ResolutionNotes
	}
	
	if req.OccurredAt != nil {
		entity.OccurredAt = *req.OccurredAt
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update pos_error_logs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated pos_error_logs",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a pos_error_logs
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting pos_error_logs",
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
		return fmt.Errorf("failed to get pos_error_logs: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("pos_error_logs not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete pos_error_logs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted pos_error_logs",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *pos_error_log.PosErrorLogs) *dto.PosErrorLogsResponse {
	return &dto.PosErrorLogsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ErrorLevel: entity.ErrorLevel,
		
		ErrorCode: entity.ErrorCode,
		
		ErrorMessage: entity.ErrorMessage,
		
		DeviceId: entity.DeviceId,
		
		UserId: entity.UserId,
		
		PosSessionId: entity.PosSessionId,
		
		SaleId: entity.SaleId,
		
		StackTrace: entity.StackTrace,
		
		RequestData: entity.RequestData,
		
		ErrorData: entity.ErrorData,
		
		IsResolved: entity.IsResolved,
		
		ResolvedBy: entity.ResolvedBy,
		
		ResolvedAt: entity.ResolvedAt,
		
		ResolutionNotes: entity.ResolutionNotes,
		
		OccurredAt: entity.OccurredAt,
		
		CreatedAt: entity.CreatedAt,
		
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


// validateBusinessRules validates business rules for pos_error_logs
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *pos_error_log.PosErrorLogs) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a pos_error_logs can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

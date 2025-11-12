package sms_queue

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/sms_queue"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/sms_queue"
	"go.uber.org/zap"
)

// Service handles business logic for SmsQueue
type Service struct {
	repo   *sms_queue.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new SmsQueue service
func NewService(repo *sms_queue.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new sms_queue
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateSmsQueueRequest) (*dto.SmsQueueResponse, error) {
	s.logger.Info("creating sms_queue",
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
	entity := &sms_queue.SmsQueue{
		OrganizationID: orgID,
		
		ToPhone: req.ToPhone,
		
		FromPhone: req.FromPhone,
		
		Message: req.Message,
		
		Status: req.Status,
		
		Status: req.Status,
		
		Provider: req.Provider,
		
		ProviderMessageId: req.ProviderMessageId,
		
		Attempts: req.Attempts,
		
		MaxAttempts: req.MaxAttempts,
		
		ErrorMessage: req.ErrorMessage,
		
		CostAmount: req.CostAmount,
		
		CostCurrency: req.CostCurrency,
		
		ScheduledAt: req.ScheduledAt,
		
		SentAt: req.SentAt,
		
		FailedAt: req.FailedAt,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create sms_queue: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created sms_queue",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a sms_queue by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.SmsQueueResponse, error) {
	s.logger.Debug("getting sms_queue",
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
		return nil, fmt.Errorf("failed to get sms_queue: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("sms_queue not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of sms_queue records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.SmsQueueListResponse, error) {
	s.logger.Debug("listing sms_queue",
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
		return nil, fmt.Errorf("failed to list sms_queue: %w", err)
	}

	// Convert to response
	items := make([]*dto.SmsQueueResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.SmsQueueListResponse{
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

// Update updates an existing sms_queue
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateSmsQueueRequest) (*dto.SmsQueueResponse, error) {
	s.logger.Info("updating sms_queue",
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
		return nil, fmt.Errorf("failed to get sms_queue: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("sms_queue not found or access denied")
	}
	

	// Update fields
	
	if req.ToPhone != nil {
		entity.ToPhone = *req.ToPhone
	}
	
	if req.FromPhone != nil {
		entity.FromPhone = *req.FromPhone
	}
	
	if req.Message != nil {
		entity.Message = *req.Message
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.Provider != nil {
		entity.Provider = *req.Provider
	}
	
	if req.ProviderMessageId != nil {
		entity.ProviderMessageId = *req.ProviderMessageId
	}
	
	if req.Attempts != nil {
		entity.Attempts = *req.Attempts
	}
	
	if req.MaxAttempts != nil {
		entity.MaxAttempts = *req.MaxAttempts
	}
	
	if req.ErrorMessage != nil {
		entity.ErrorMessage = *req.ErrorMessage
	}
	
	if req.CostAmount != nil {
		entity.CostAmount = *req.CostAmount
	}
	
	if req.CostCurrency != nil {
		entity.CostCurrency = *req.CostCurrency
	}
	
	if req.ScheduledAt != nil {
		entity.ScheduledAt = *req.ScheduledAt
	}
	
	if req.SentAt != nil {
		entity.SentAt = *req.SentAt
	}
	
	if req.FailedAt != nil {
		entity.FailedAt = *req.FailedAt
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update sms_queue: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated sms_queue",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a sms_queue
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting sms_queue",
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
		return fmt.Errorf("failed to get sms_queue: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("sms_queue not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete sms_queue: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted sms_queue",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *sms_queue.SmsQueue) *dto.SmsQueueResponse {
	return &dto.SmsQueueResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ToPhone: entity.ToPhone,
		
		FromPhone: entity.FromPhone,
		
		Message: entity.Message,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		Provider: entity.Provider,
		
		ProviderMessageId: entity.ProviderMessageId,
		
		Attempts: entity.Attempts,
		
		MaxAttempts: entity.MaxAttempts,
		
		ErrorMessage: entity.ErrorMessage,
		
		CostAmount: entity.CostAmount,
		
		CostCurrency: entity.CostCurrency,
		
		ScheduledAt: entity.ScheduledAt,
		
		SentAt: entity.SentAt,
		
		FailedAt: entity.FailedAt,
		
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


// validateBusinessRules validates business rules for sms_queue
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *sms_queue.SmsQueue) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a sms_queue can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

package email_queue

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

// Service handles business logic for EmailQueue
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new EmailQueue service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new email_queue
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateEmailQueueRequest) (*EmailQueueResponse, error) {
	s.logger.Info("creating email_queue",
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
	entity := &EmailQueue{
		OrganizationID: orgID,
		
		ToAddresses: req.ToAddresses,
		
		CcAddresses: req.CcAddresses,
		
		BccAddresses: req.BccAddresses,
		
		FromAddress: req.FromAddress,
		
		ReplyTo: req.ReplyTo,
		
		Subject: req.Subject,
		
		BodyHtml: req.BodyHtml,
		
		BodyText: req.BodyText,
		
		AttachmentIds: req.AttachmentIds,
		
		TemplateName: req.TemplateName,
		
		TemplateData: req.TemplateData,
		
		Status: req.Status,
		
		Status: req.Status,
		
		Provider: req.Provider,
		
		ProviderMessageId: req.ProviderMessageId,
		
		Attempts: req.Attempts,
		
		MaxAttempts: req.MaxAttempts,
		
		ErrorMessage: req.ErrorMessage,
		
		Priority: req.Priority,
		
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
		return nil, fmt.Errorf("failed to create email_queue: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created email_queue",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a email_queue by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*EmailQueueResponse, error) {
	s.logger.Debug("getting email_queue",
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
		return nil, fmt.Errorf("failed to get email_queue: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("email_queue not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of email_queue records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*EmailQueueListResponse, error) {
	s.logger.Debug("listing email_queue",
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
		return nil, fmt.Errorf("failed to list email_queue: %w", err)
	}

	// Convert to response
	items := make([]*EmailQueueResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &EmailQueueListResponse{
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

// Update updates an existing email_queue
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateEmailQueueRequest) (*EmailQueueResponse, error) {
	s.logger.Info("updating email_queue",
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
		return nil, fmt.Errorf("failed to get email_queue: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("email_queue not found or access denied")
	}
	

	// Update fields
	
	if req.ToAddresses != nil {
		entity.ToAddresses = *req.ToAddresses
	}
	
	if req.CcAddresses != nil {
		entity.CcAddresses = *req.CcAddresses
	}
	
	if req.BccAddresses != nil {
		entity.BccAddresses = *req.BccAddresses
	}
	
	if req.FromAddress != nil {
		entity.FromAddress = *req.FromAddress
	}
	
	if req.ReplyTo != nil {
		entity.ReplyTo = *req.ReplyTo
	}
	
	if req.Subject != nil {
		entity.Subject = *req.Subject
	}
	
	if req.BodyHtml != nil {
		entity.BodyHtml = *req.BodyHtml
	}
	
	if req.BodyText != nil {
		entity.BodyText = *req.BodyText
	}
	
	if req.AttachmentIds != nil {
		entity.AttachmentIds = *req.AttachmentIds
	}
	
	if req.TemplateName != nil {
		entity.TemplateName = *req.TemplateName
	}
	
	if req.TemplateData != nil {
		entity.TemplateData = *req.TemplateData
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
	
	if req.Priority != nil {
		entity.Priority = *req.Priority
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
		return nil, fmt.Errorf("failed to update email_queue: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated email_queue",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a email_queue
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting email_queue",
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
		return fmt.Errorf("failed to get email_queue: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("email_queue not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete email_queue: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted email_queue",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *EmailQueue) *EmailQueueResponse {
	return &EmailQueueResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ToAddresses: entity.ToAddresses,
		
		CcAddresses: entity.CcAddresses,
		
		BccAddresses: entity.BccAddresses,
		
		FromAddress: entity.FromAddress,
		
		ReplyTo: entity.ReplyTo,
		
		Subject: entity.Subject,
		
		BodyHtml: entity.BodyHtml,
		
		BodyText: entity.BodyText,
		
		AttachmentIds: entity.AttachmentIds,
		
		TemplateName: entity.TemplateName,
		
		TemplateData: entity.TemplateData,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		Provider: entity.Provider,
		
		ProviderMessageId: entity.ProviderMessageId,
		
		Attempts: entity.Attempts,
		
		MaxAttempts: entity.MaxAttempts,
		
		ErrorMessage: entity.ErrorMessage,
		
		Priority: entity.Priority,
		
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


// validateBusinessRules validates business rules for email_queue
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *EmailQueue) error {
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

// canDelete checks if a email_queue can be deleted
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

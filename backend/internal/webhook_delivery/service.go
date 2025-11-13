package webhook_delivery

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

// Service handles business logic for WebhookDeliveries
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new WebhookDeliveries service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new webhook_deliveries
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateWebhookDeliveriesRequest) (*WebhookDeliveriesResponse, error) {
	s.logger.Info("creating webhook_deliveries",
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
	entity := &WebhookDeliveries{
		OrganizationID: orgID,
		
		WebhookId: req.WebhookId,
		
		EventType: req.EventType,
		
		EventId: req.EventId,
		
		Status: req.Status,
		
		Status: req.Status,
		
		RequestUrl: req.RequestUrl,
		
		RequestMethod: req.RequestMethod,
		
		RequestHeaders: req.RequestHeaders,
		
		RequestBody: req.RequestBody,
		
		ResponseStatusCode: req.ResponseStatusCode,
		
		ResponseHeaders: req.ResponseHeaders,
		
		ResponseBody: req.ResponseBody,
		
		AttemptNumber: req.AttemptNumber,
		
		DurationMs: req.DurationMs,
		
		NextRetryAt: req.NextRetryAt,
		
		ErrorMessage: req.ErrorMessage,
		
		DeliveredAt: req.DeliveredAt,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create webhook_deliveries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created webhook_deliveries",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a webhook_deliveries by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*WebhookDeliveriesResponse, error) {
	s.logger.Debug("getting webhook_deliveries",
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
		return nil, fmt.Errorf("failed to get webhook_deliveries: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("webhook_deliveries not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of webhook_deliveries records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*WebhookDeliveriesListResponse, error) {
	s.logger.Debug("listing webhook_deliveries",
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
		return nil, fmt.Errorf("failed to list webhook_deliveries: %w", err)
	}

	// Convert to response
	items := make([]*WebhookDeliveriesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &WebhookDeliveriesListResponse{
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

// Update updates an existing webhook_deliveries
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateWebhookDeliveriesRequest) (*WebhookDeliveriesResponse, error) {
	s.logger.Info("updating webhook_deliveries",
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
		return nil, fmt.Errorf("failed to get webhook_deliveries: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("webhook_deliveries not found or access denied")
	}
	

	// Update fields
	
	if req.WebhookId != nil {
		entity.WebhookId = *req.WebhookId
	}
	
	if req.EventType != nil {
		entity.EventType = *req.EventType
	}
	
	if req.EventId != nil {
		entity.EventId = *req.EventId
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.RequestUrl != nil {
		entity.RequestUrl = *req.RequestUrl
	}
	
	if req.RequestMethod != nil {
		entity.RequestMethod = *req.RequestMethod
	}
	
	if req.RequestHeaders != nil {
		entity.RequestHeaders = *req.RequestHeaders
	}
	
	if req.RequestBody != nil {
		entity.RequestBody = *req.RequestBody
	}
	
	if req.ResponseStatusCode != nil {
		entity.ResponseStatusCode = *req.ResponseStatusCode
	}
	
	if req.ResponseHeaders != nil {
		entity.ResponseHeaders = *req.ResponseHeaders
	}
	
	if req.ResponseBody != nil {
		entity.ResponseBody = *req.ResponseBody
	}
	
	if req.AttemptNumber != nil {
		entity.AttemptNumber = *req.AttemptNumber
	}
	
	if req.DurationMs != nil {
		entity.DurationMs = *req.DurationMs
	}
	
	if req.NextRetryAt != nil {
		entity.NextRetryAt = *req.NextRetryAt
	}
	
	if req.ErrorMessage != nil {
		entity.ErrorMessage = *req.ErrorMessage
	}
	
	if req.DeliveredAt != nil {
		entity.DeliveredAt = *req.DeliveredAt
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update webhook_deliveries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated webhook_deliveries",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a webhook_deliveries
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting webhook_deliveries",
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
		return fmt.Errorf("failed to get webhook_deliveries: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("webhook_deliveries not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete webhook_deliveries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted webhook_deliveries",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *WebhookDeliveries) *WebhookDeliveriesResponse {
	return &WebhookDeliveriesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		WebhookId: entity.WebhookId,
		
		EventType: entity.EventType,
		
		EventId: entity.EventId,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		RequestUrl: entity.RequestUrl,
		
		RequestMethod: entity.RequestMethod,
		
		RequestHeaders: entity.RequestHeaders,
		
		RequestBody: entity.RequestBody,
		
		ResponseStatusCode: entity.ResponseStatusCode,
		
		ResponseHeaders: entity.ResponseHeaders,
		
		ResponseBody: entity.ResponseBody,
		
		AttemptNumber: entity.AttemptNumber,
		
		DurationMs: entity.DurationMs,
		
		NextRetryAt: entity.NextRetryAt,
		
		ErrorMessage: entity.ErrorMessage,
		
		CreatedAt: entity.CreatedAt,
		
		DeliveredAt: entity.DeliveredAt,
		
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


// validateBusinessRules validates business rules for webhook_deliveries
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *WebhookDeliveries) error {
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

// canDelete checks if a webhook_deliveries can be deleted
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

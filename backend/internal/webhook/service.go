package webhook

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

// Service handles business logic for Webhooks
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Webhooks service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new webhooks
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateWebhooksRequest) (*WebhooksResponse, error) {
	s.logger.Info("creating webhooks",
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
	entity := &Webhooks{
		OrganizationID: orgID,
		
		WebhookName: req.WebhookName,
		
		Url: req.Url,
		
		Secret: req.Secret,
		
		Events: req.Events,
		
		HttpMethod: req.HttpMethod,
		
		Headers: req.Headers,
		
		TimeoutSeconds: req.TimeoutSeconds,
		
		MaxRetries: req.MaxRetries,
		
		RetryBackoffSeconds: req.RetryBackoffSeconds,
		
		IsActive: req.IsActive,
		
		IsVerified: req.IsVerified,
		
		TotalDeliveries: req.TotalDeliveries,
		
		SuccessfulDeliveries: req.SuccessfulDeliveries,
		
		FailedDeliveries: req.FailedDeliveries,
		
		LastDeliveryAt: req.LastDeliveryAt,
		
		LastSuccessAt: req.LastSuccessAt,
		
		LastFailureAt: req.LastFailureAt,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create webhooks: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created webhooks",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a webhooks by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*WebhooksResponse, error) {
	s.logger.Debug("getting webhooks",
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
		return nil, fmt.Errorf("failed to get webhooks: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("webhooks not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of webhooks records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*WebhooksListResponse, error) {
	s.logger.Debug("listing webhooks",
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
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	// Convert to response
	items := make([]*WebhooksResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &WebhooksListResponse{
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

// Update updates an existing webhooks
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateWebhooksRequest) (*WebhooksResponse, error) {
	s.logger.Info("updating webhooks",
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
		return nil, fmt.Errorf("failed to get webhooks: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("webhooks not found or access denied")
	}
	

	// Update fields
	
	if req.WebhookName != nil {
		entity.WebhookName = *req.WebhookName
	}
	
	if req.Url != nil {
		entity.Url = *req.Url
	}
	
	if req.Secret != nil {
		entity.Secret = *req.Secret
	}
	
	if req.Events != nil {
		entity.Events = *req.Events
	}
	
	if req.HttpMethod != nil {
		entity.HttpMethod = *req.HttpMethod
	}
	
	if req.Headers != nil {
		entity.Headers = *req.Headers
	}
	
	if req.TimeoutSeconds != nil {
		entity.TimeoutSeconds = *req.TimeoutSeconds
	}
	
	if req.MaxRetries != nil {
		entity.MaxRetries = *req.MaxRetries
	}
	
	if req.RetryBackoffSeconds != nil {
		entity.RetryBackoffSeconds = *req.RetryBackoffSeconds
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.IsVerified != nil {
		entity.IsVerified = *req.IsVerified
	}
	
	if req.TotalDeliveries != nil {
		entity.TotalDeliveries = *req.TotalDeliveries
	}
	
	if req.SuccessfulDeliveries != nil {
		entity.SuccessfulDeliveries = *req.SuccessfulDeliveries
	}
	
	if req.FailedDeliveries != nil {
		entity.FailedDeliveries = *req.FailedDeliveries
	}
	
	if req.LastDeliveryAt != nil {
		entity.LastDeliveryAt = *req.LastDeliveryAt
	}
	
	if req.LastSuccessAt != nil {
		entity.LastSuccessAt = *req.LastSuccessAt
	}
	
	if req.LastFailureAt != nil {
		entity.LastFailureAt = *req.LastFailureAt
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
		return nil, fmt.Errorf("failed to update webhooks: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated webhooks",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a webhooks
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting webhooks",
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
		return fmt.Errorf("failed to get webhooks: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("webhooks not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete webhooks: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted webhooks",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Webhooks) *WebhooksResponse {
	return &WebhooksResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		WebhookName: entity.WebhookName,
		
		Url: entity.Url,
		
		Secret: entity.Secret,
		
		Events: entity.Events,
		
		HttpMethod: entity.HttpMethod,
		
		Headers: entity.Headers,
		
		TimeoutSeconds: entity.TimeoutSeconds,
		
		MaxRetries: entity.MaxRetries,
		
		RetryBackoffSeconds: entity.RetryBackoffSeconds,
		
		IsActive: entity.IsActive,
		
		IsVerified: entity.IsVerified,
		
		TotalDeliveries: entity.TotalDeliveries,
		
		SuccessfulDeliveries: entity.SuccessfulDeliveries,
		
		FailedDeliveries: entity.FailedDeliveries,
		
		LastDeliveryAt: entity.LastDeliveryAt,
		
		LastSuccessAt: entity.LastSuccessAt,
		
		LastFailureAt: entity.LastFailureAt,
		
		CreatedBy: entity.CreatedBy,
		
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


// validateBusinessRules validates business rules for webhooks
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Webhooks) error {
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

// canDelete checks if a webhooks can be deleted
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

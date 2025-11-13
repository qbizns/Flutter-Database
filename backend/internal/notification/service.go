package notification

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for Notifications
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Notifications service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new notifications
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateNotificationsRequest) (*NotificationsResponse, error) {
	s.logger.Info("creating notifications",
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
	entity := &Notifications{
		OrganizationId: orgID,
		
		UserId: req.UserId,
		
		NotificationType: req.NotificationType,
		
		Category: req.Category,
		
		Title: req.Title,
		
		Message: req.Message,
		
		ActionUrl: req.ActionUrl,
		
		ActionLabel: req.ActionLabel,
		
		Channels: req.Channels,
		
		IsRead: req.IsRead,
		
		ReadAt: req.ReadAt,
		
		RelatedEntityType: req.RelatedEntityType,
		
		RelatedEntityId: req.RelatedEntityId,
		
		Priority: req.Priority,
		
		ExpiresAt: req.ExpiresAt,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create notifications: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created notifications",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a notifications by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*NotificationsResponse, error) {
	s.logger.Debug("getting notifications",
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
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("notifications not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of notifications records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*NotificationsListResponse, error) {
	s.logger.Debug("listing notifications",
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
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	// Convert to response
	items := make([]*NotificationsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &NotificationsListResponse{
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

// Update updates an existing notifications
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateNotificationsRequest) (*NotificationsResponse, error) {
	s.logger.Info("updating notifications",
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
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("notifications not found or access denied")
	}
	

	// Update fields
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.NotificationType != nil {
		entity.NotificationType = *req.NotificationType
	}
	
	if req.Category != nil {
		entity.Category = *req.Category
	}
	
	if req.Title != nil {
		entity.Title = *req.Title
	}
	
	if req.Message != nil {
		entity.Message = *req.Message
	}
	
	if req.ActionUrl != nil {
		entity.ActionUrl = *req.ActionUrl
	}
	
	if req.ActionLabel != nil {
		entity.ActionLabel = *req.ActionLabel
	}
	
	if req.Channels != nil {
		entity.Channels = *req.Channels
	}
	
	if req.IsRead != nil {
		entity.IsRead = req.IsRead
	}
	
	if req.ReadAt != nil {
		entity.ReadAt = req.ReadAt
	}
	
	if req.RelatedEntityType != nil {
		entity.RelatedEntityType = *req.RelatedEntityType
	}
	
	if req.RelatedEntityId != nil {
		entity.RelatedEntityId = *req.RelatedEntityId
	}
	
	if req.Priority != nil {
		entity.Priority = req.Priority
	}
	
	if req.ExpiresAt != nil {
		entity.ExpiresAt = req.ExpiresAt
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update notifications: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated notifications",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a notifications
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting notifications",
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
		return fmt.Errorf("failed to get notifications: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("notifications not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted notifications",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Notifications) *NotificationsResponse {
	return &NotificationsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		UserId: entity.UserId,
		
		NotificationType: entity.NotificationType,
		
		Category: entity.Category,
		
		Title: entity.Title,
		
		Message: entity.Message,
		
		ActionUrl: entity.ActionUrl,
		
		ActionLabel: entity.ActionLabel,
		
		Channels: entity.Channels,
		
		IsRead: entity.IsRead,
		
		ReadAt: entity.ReadAt,
		
		RelatedEntityType: entity.RelatedEntityType,
		
		RelatedEntityId: entity.RelatedEntityId,
		
		Priority: entity.Priority,
		
		ExpiresAt: entity.ExpiresAt,
		
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


// validateBusinessRules validates business rules for notifications
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Notifications) error {
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

// canDelete checks if a notifications can be deleted
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

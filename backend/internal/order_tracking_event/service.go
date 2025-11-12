package order_tracking_event

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

// Service handles business logic for OrderTrackingEvents
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new OrderTrackingEvents service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new order_tracking_events
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateOrderTrackingEventsRequest) (*OrderTrackingEventsResponse, error) {
	s.logger.Info("creating order_tracking_events",
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
	entity := &OrderTrackingEvents{
		OrganizationID: orgID,
		
		OrderId: req.OrderId,
		
		DeliveryAssignmentId: req.DeliveryAssignmentId,
		
		EventType: req.EventType,
		
		EventTimestamp: req.EventTimestamp,
		
		EventMessage: req.EventMessage,
		
		Location: req.Location,
		
		LocationName: req.LocationName,
		
		ActorType: req.ActorType,
		
		ActorId: req.ActorId,
		
		ActorName: req.ActorName,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		'orderPlaced',: req.'orderPlaced',,
		
		'readyForPickup',: req.'readyForPickup',,
		
		'arrived',: req.'arrived',,
		
		'rescheduled',: req.'rescheduled',,
		
		'system',: req.'system',,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create order_tracking_events: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created order_tracking_events",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a order_tracking_events by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrderTrackingEventsResponse, error) {
	s.logger.Debug("getting order_tracking_events",
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
		return nil, fmt.Errorf("failed to get order_tracking_events: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("order_tracking_events not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of order_tracking_events records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*OrderTrackingEventsListResponse, error) {
	s.logger.Debug("listing order_tracking_events",
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
		return nil, fmt.Errorf("failed to list order_tracking_events: %w", err)
	}

	// Convert to response
	items := make([]*OrderTrackingEventsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &OrderTrackingEventsListResponse{
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

// Update updates an existing order_tracking_events
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateOrderTrackingEventsRequest) (*OrderTrackingEventsResponse, error) {
	s.logger.Info("updating order_tracking_events",
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
		return nil, fmt.Errorf("failed to get order_tracking_events: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("order_tracking_events not found or access denied")
	}
	

	// Update fields
	
	if req.OrderId != nil {
		entity.OrderId = *req.OrderId
	}
	
	if req.DeliveryAssignmentId != nil {
		entity.DeliveryAssignmentId = *req.DeliveryAssignmentId
	}
	
	if req.EventType != nil {
		entity.EventType = *req.EventType
	}
	
	if req.EventTimestamp != nil {
		entity.EventTimestamp = *req.EventTimestamp
	}
	
	if req.EventMessage != nil {
		entity.EventMessage = *req.EventMessage
	}
	
	if req.Location != nil {
		entity.Location = *req.Location
	}
	
	if req.LocationName != nil {
		entity.LocationName = *req.LocationName
	}
	
	if req.ActorType != nil {
		entity.ActorType = *req.ActorType
	}
	
	if req.ActorId != nil {
		entity.ActorId = *req.ActorId
	}
	
	if req.ActorName != nil {
		entity.ActorName = *req.ActorName
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.'orderPlaced', != nil {
		entity.'orderPlaced', = *req.'orderPlaced',
	}
	
	if req.'readyForPickup', != nil {
		entity.'readyForPickup', = *req.'readyForPickup',
	}
	
	if req.'arrived', != nil {
		entity.'arrived', = *req.'arrived',
	}
	
	if req.'rescheduled', != nil {
		entity.'rescheduled', = *req.'rescheduled',
	}
	
	if req.'system', != nil {
		entity.'system', = *req.'system',
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update order_tracking_events: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated order_tracking_events",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a order_tracking_events
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting order_tracking_events",
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
		return fmt.Errorf("failed to get order_tracking_events: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("order_tracking_events not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete order_tracking_events: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted order_tracking_events",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *OrderTrackingEvents) *OrderTrackingEventsResponse {
	return &OrderTrackingEventsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		OrderId: entity.OrderId,
		
		DeliveryAssignmentId: entity.DeliveryAssignmentId,
		
		EventType: entity.EventType,
		
		EventTimestamp: entity.EventTimestamp,
		
		EventMessage: entity.EventMessage,
		
		Location: entity.Location,
		
		LocationName: entity.LocationName,
		
		ActorType: entity.ActorType,
		
		ActorId: entity.ActorId,
		
		ActorName: entity.ActorName,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		'orderPlaced',: entity.'orderPlaced',,
		
		'readyForPickup',: entity.'readyForPickup',,
		
		'arrived',: entity.'arrived',,
		
		'rescheduled',: entity.'rescheduled',,
		
		'system',: entity.'system',,
		
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


// validateBusinessRules validates business rules for order_tracking_events
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *OrderTrackingEvents) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a order_tracking_events can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

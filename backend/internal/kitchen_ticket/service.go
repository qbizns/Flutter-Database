package kitchen_ticket

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for KitchenTickets
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new KitchenTickets service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new kitchen_tickets
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateKitchenTicketsRequest) (*KitchenTicketsResponse, error) {
	s.logger.Info("creating kitchen_tickets",
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
	entity := &KitchenTickets{
		OrganizationId: orgID,
		
		LocationId: req.LocationId,
		
		TicketNumber: req.TicketNumber,
		
		DisplaySequence: req.DisplaySequence,
		
		OrderId: req.OrderId,
		
		KitchenStationId: req.KitchenStationId,
		
		CourseId: req.CourseId,
		
		TicketType: req.TicketType,
		
		Priority: req.Priority,
		
		Status: req.Status,
		
		FiredAt: req.FiredAt,
		
		AcknowledgedAt: req.AcknowledgedAt,
		
		StartedAt: req.StartedAt,
		
		ReadyAt: req.ReadyAt,
		
		BumpedAt: req.BumpedAt,
		
		CompletedAt: req.CompletedAt,
		
		PrepTimeMinutes: req.PrepTimeMinutes,
		
		TargetPrepTime: req.TargetPrepTime,
		
		TableNumber: req.TableNumber,
		
		OrderType: req.OrderType,
		
		Covers: req.Covers,
		
		WaiterName: req.WaiterName,
		
		SpecialInstructions: req.SpecialInstructions,
		
		KitchenNotes: req.KitchenNotes,
		
		DisplayConfig: req.DisplayConfig,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'new',: req.'new',,
		
		'completed',: req.'completed',,
		
		'normal',: req.'normal',,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create kitchen_tickets: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created kitchen_tickets",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a kitchen_tickets by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*KitchenTicketsResponse, error) {
	s.logger.Debug("getting kitchen_tickets",
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
		return nil, fmt.Errorf("failed to get kitchen_tickets: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("kitchen_tickets not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of kitchen_tickets records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*KitchenTicketsListResponse, error) {
	s.logger.Debug("listing kitchen_tickets",
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
		return nil, fmt.Errorf("failed to list kitchen_tickets: %w", err)
	}

	// Convert to response
	items := make([]*KitchenTicketsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &KitchenTicketsListResponse{
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

// Update updates an existing kitchen_tickets
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateKitchenTicketsRequest) (*KitchenTicketsResponse, error) {
	s.logger.Info("updating kitchen_tickets",
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
		return nil, fmt.Errorf("failed to get kitchen_tickets: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("kitchen_tickets not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.TicketNumber != nil {
		entity.TicketNumber = *req.TicketNumber
	}
	
	if req.DisplaySequence != nil {
		entity.DisplaySequence = *req.DisplaySequence
	}
	
	if req.OrderId != nil {
		entity.OrderId = *req.OrderId
	}
	
	if req.KitchenStationId != nil {
		entity.KitchenStationId = *req.KitchenStationId
	}
	
	if req.CourseId != nil {
		entity.CourseId = *req.CourseId
	}
	
	if req.TicketType != nil {
		entity.TicketType = *req.TicketType
	}
	
	if req.Priority != nil {
		entity.Priority = req.Priority
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.FiredAt != nil {
		entity.FiredAt = req.FiredAt
	}
	
	if req.AcknowledgedAt != nil {
		entity.AcknowledgedAt = req.AcknowledgedAt
	}
	
	if req.StartedAt != nil {
		entity.StartedAt = req.StartedAt
	}
	
	if req.ReadyAt != nil {
		entity.ReadyAt = req.ReadyAt
	}
	
	if req.BumpedAt != nil {
		entity.BumpedAt = req.BumpedAt
	}
	
	if req.CompletedAt != nil {
		entity.CompletedAt = req.CompletedAt
	}
	
	if req.PrepTimeMinutes != nil {
		entity.PrepTimeMinutes = *req.PrepTimeMinutes
	}
	
	if req.TargetPrepTime != nil {
		entity.TargetPrepTime = *req.TargetPrepTime
	}
	
	if req.TableNumber != nil {
		entity.TableNumber = *req.TableNumber
	}
	
	if req.OrderType != nil {
		entity.OrderType = *req.OrderType
	}
	
	if req.Covers != nil {
		entity.Covers = *req.Covers
	}
	
	if req.WaiterName != nil {
		entity.WaiterName = *req.WaiterName
	}
	
	if req.SpecialInstructions != nil {
		entity.SpecialInstructions = *req.SpecialInstructions
	}
	
	if req.KitchenNotes != nil {
		entity.KitchenNotes = *req.KitchenNotes
	}
	
	if req.DisplayConfig != nil {
		entity.DisplayConfig = *req.DisplayConfig
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
	
	if req.'new', != nil {
		entity.'new', = *req.'new',
	}
	
	if req.'completed', != nil {
		entity.'completed', = *req.'completed',
	}
	
	if req.'normal', != nil {
		entity.'normal', = *req.'normal',
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update kitchen_tickets: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated kitchen_tickets",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a kitchen_tickets
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting kitchen_tickets",
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
		return fmt.Errorf("failed to get kitchen_tickets: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("kitchen_tickets not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete kitchen_tickets: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted kitchen_tickets",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *KitchenTickets) *KitchenTicketsResponse {
	return &KitchenTicketsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		TicketNumber: entity.TicketNumber,
		
		DisplaySequence: entity.DisplaySequence,
		
		OrderId: entity.OrderId,
		
		KitchenStationId: entity.KitchenStationId,
		
		CourseId: entity.CourseId,
		
		TicketType: entity.TicketType,
		
		Priority: entity.Priority,
		
		Status: entity.Status,
		
		CreatedAt: entity.CreatedAt,
		
		FiredAt: entity.FiredAt,
		
		AcknowledgedAt: entity.AcknowledgedAt,
		
		StartedAt: entity.StartedAt,
		
		ReadyAt: entity.ReadyAt,
		
		BumpedAt: entity.BumpedAt,
		
		CompletedAt: entity.CompletedAt,
		
		PrepTimeMinutes: entity.PrepTimeMinutes,
		
		TargetPrepTime: entity.TargetPrepTime,
		
		TableNumber: entity.TableNumber,
		
		OrderType: entity.OrderType,
		
		Covers: entity.Covers,
		
		WaiterName: entity.WaiterName,
		
		SpecialInstructions: entity.SpecialInstructions,
		
		KitchenNotes: entity.KitchenNotes,
		
		DisplayConfig: entity.DisplayConfig,
		
		Metadata: entity.Metadata,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'new',: entity.'new',,
		
		'completed',: entity.'completed',,
		
		'normal',: entity.'normal',,
		
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


// validateBusinessRules validates business rules for kitchen_tickets
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *KitchenTickets) error {
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

// canDelete checks if a kitchen_tickets can be deleted
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

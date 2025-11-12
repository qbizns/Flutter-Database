package order_item

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

// Service handles business logic for OrderItems
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new OrderItems service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new order_items
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateOrderItemsRequest) (*OrderItemsResponse, error) {
	s.logger.Info("creating order_items",
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
	entity := &OrderItems{
		OrganizationID: orgID,
		
		OrderId: req.OrderId,
		
		ProductId: req.ProductId,
		
		ProductVariantId: req.ProductVariantId,
		
		ItemName: req.ItemName,
		
		Quantity: req.Quantity,
		
		UnitPrice: req.UnitPrice,
		
		CourseId: req.CourseId,
		
		CoursePosition: req.CoursePosition,
		
		FireTime: req.FireTime,
		
		KitchenStationId: req.KitchenStationId,
		
		KitchenTicketId: req.KitchenTicketId,
		
		Status: req.Status,
		
		FiredAt: req.FiredAt,
		
		AcknowledgedAt: req.AcknowledgedAt,
		
		StartedPreparingAt: req.StartedPreparingAt,
		
		ReadyAt: req.ReadyAt,
		
		ServedAt: req.ServedAt,
		
		ModifiersTotal: req.ModifiersTotal,
		
		DiscountAmount: req.DiscountAmount,
		
		LineTotal: req.LineTotal,
		
		SpecialInstructions: req.SpecialInstructions,
		
		CustomerNotes: req.CustomerNotes,
		
		KitchenNotes: req.KitchenNotes,
		
		SeatNumber: req.SeatNumber,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'pending',: req.'pending',,
		
		'served',: req.'served',,
		
		UnitPrice: req.UnitPrice,
		
		DiscountAmount: req.DiscountAmount,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create order_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created order_items",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a order_items by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrderItemsResponse, error) {
	s.logger.Debug("getting order_items",
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
		return nil, fmt.Errorf("failed to get order_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("order_items not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of order_items records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*OrderItemsListResponse, error) {
	s.logger.Debug("listing order_items",
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
		return nil, fmt.Errorf("failed to list order_items: %w", err)
	}

	// Convert to response
	items := make([]*OrderItemsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &OrderItemsListResponse{
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

// Update updates an existing order_items
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateOrderItemsRequest) (*OrderItemsResponse, error) {
	s.logger.Info("updating order_items",
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
		return nil, fmt.Errorf("failed to get order_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("order_items not found or access denied")
	}
	

	// Update fields
	
	if req.OrderId != nil {
		entity.OrderId = *req.OrderId
	}
	
	if req.ProductId != nil {
		entity.ProductId = *req.ProductId
	}
	
	if req.ProductVariantId != nil {
		entity.ProductVariantId = *req.ProductVariantId
	}
	
	if req.ItemName != nil {
		entity.ItemName = *req.ItemName
	}
	
	if req.Quantity != nil {
		entity.Quantity = *req.Quantity
	}
	
	if req.UnitPrice != nil {
		entity.UnitPrice = *req.UnitPrice
	}
	
	if req.CourseId != nil {
		entity.CourseId = *req.CourseId
	}
	
	if req.CoursePosition != nil {
		entity.CoursePosition = *req.CoursePosition
	}
	
	if req.FireTime != nil {
		entity.FireTime = *req.FireTime
	}
	
	if req.KitchenStationId != nil {
		entity.KitchenStationId = *req.KitchenStationId
	}
	
	if req.KitchenTicketId != nil {
		entity.KitchenTicketId = *req.KitchenTicketId
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.FiredAt != nil {
		entity.FiredAt = *req.FiredAt
	}
	
	if req.AcknowledgedAt != nil {
		entity.AcknowledgedAt = *req.AcknowledgedAt
	}
	
	if req.StartedPreparingAt != nil {
		entity.StartedPreparingAt = *req.StartedPreparingAt
	}
	
	if req.ReadyAt != nil {
		entity.ReadyAt = *req.ReadyAt
	}
	
	if req.ServedAt != nil {
		entity.ServedAt = *req.ServedAt
	}
	
	if req.ModifiersTotal != nil {
		entity.ModifiersTotal = *req.ModifiersTotal
	}
	
	if req.DiscountAmount != nil {
		entity.DiscountAmount = *req.DiscountAmount
	}
	
	if req.LineTotal != nil {
		entity.LineTotal = *req.LineTotal
	}
	
	if req.SpecialInstructions != nil {
		entity.SpecialInstructions = *req.SpecialInstructions
	}
	
	if req.CustomerNotes != nil {
		entity.CustomerNotes = *req.CustomerNotes
	}
	
	if req.KitchenNotes != nil {
		entity.KitchenNotes = *req.KitchenNotes
	}
	
	if req.SeatNumber != nil {
		entity.SeatNumber = *req.SeatNumber
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = *req.UpdatedBy
	}
	
	if req.'pending', != nil {
		entity.'pending', = *req.'pending',
	}
	
	if req.'served', != nil {
		entity.'served', = *req.'served',
	}
	
	if req.UnitPrice != nil {
		entity.UnitPrice = *req.UnitPrice
	}
	
	if req.DiscountAmount != nil {
		entity.DiscountAmount = *req.DiscountAmount
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update order_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated order_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a order_items
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting order_items",
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
		return fmt.Errorf("failed to get order_items: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("order_items not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete order_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted order_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *OrderItems) *OrderItemsResponse {
	return &OrderItemsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		OrderId: entity.OrderId,
		
		ProductId: entity.ProductId,
		
		ProductVariantId: entity.ProductVariantId,
		
		ItemName: entity.ItemName,
		
		Quantity: entity.Quantity,
		
		UnitPrice: entity.UnitPrice,
		
		CourseId: entity.CourseId,
		
		CoursePosition: entity.CoursePosition,
		
		FireTime: entity.FireTime,
		
		KitchenStationId: entity.KitchenStationId,
		
		KitchenTicketId: entity.KitchenTicketId,
		
		Status: entity.Status,
		
		FiredAt: entity.FiredAt,
		
		AcknowledgedAt: entity.AcknowledgedAt,
		
		StartedPreparingAt: entity.StartedPreparingAt,
		
		ReadyAt: entity.ReadyAt,
		
		ServedAt: entity.ServedAt,
		
		ModifiersTotal: entity.ModifiersTotal,
		
		DiscountAmount: entity.DiscountAmount,
		
		LineTotal: entity.LineTotal,
		
		SpecialInstructions: entity.SpecialInstructions,
		
		CustomerNotes: entity.CustomerNotes,
		
		KitchenNotes: entity.KitchenNotes,
		
		SeatNumber: entity.SeatNumber,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'pending',: entity.'pending',,
		
		'served',: entity.'served',,
		
		UnitPrice: entity.UnitPrice,
		
		DiscountAmount: entity.DiscountAmount,
		
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


// validateBusinessRules validates business rules for order_items
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *OrderItems) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a order_items can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

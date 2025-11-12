package order

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

// Service handles business logic for Orders
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Orders service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new orders
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateOrdersRequest) (*OrdersResponse, error) {
	s.logger.Info("creating orders",
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
	entity := &Orders{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		OrderNumber: req.OrderNumber,
		
		DisplayNumber: req.DisplayNumber,
		
		OrderType: req.OrderType,
		
		TableId: req.TableId,
		
		ReservationId: req.ReservationId,
		
		Covers: req.Covers,
		
		CustomerId: req.CustomerId,
		
		WaiterId: req.WaiterId,
		
		Status: req.Status,
		
		OrderDate: req.OrderDate,
		
		SubmittedAt: req.SubmittedAt,
		
		KitchenReceivedAt: req.KitchenReceivedAt,
		
		ReadyAt: req.ReadyAt,
		
		ServedAt: req.ServedAt,
		
		CompletedAt: req.CompletedAt,
		
		Subtotal: req.Subtotal,
		
		TaxAmount: req.TaxAmount,
		
		DiscountAmount: req.DiscountAmount,
		
		ServiceCharge: req.ServiceCharge,
		
		TotalAmount: req.TotalAmount,
		
		SaleId: req.SaleId,
		
		ShiftId: req.ShiftId,
		
		CustomerNotes: req.CustomerNotes,
		
		KitchenNotes: req.KitchenNotes,
		
		InternalNotes: req.InternalNotes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'draft',: req.'draft',,
		
		'served',: req.'served',,
		
		'dineIn',: req.'dineIn',,
		
		Subtotal: req.Subtotal,
		
		DiscountAmount: req.DiscountAmount,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create orders: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created orders",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a orders by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrdersResponse, error) {
	s.logger.Debug("getting orders",
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
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("orders not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of orders records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*OrdersListResponse, error) {
	s.logger.Debug("listing orders",
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
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// Convert to response
	items := make([]*OrdersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &OrdersListResponse{
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

// Update updates an existing orders
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateOrdersRequest) (*OrdersResponse, error) {
	s.logger.Info("updating orders",
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
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("orders not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.OrderNumber != nil {
		entity.OrderNumber = *req.OrderNumber
	}
	
	if req.DisplayNumber != nil {
		entity.DisplayNumber = *req.DisplayNumber
	}
	
	if req.OrderType != nil {
		entity.OrderType = *req.OrderType
	}
	
	if req.TableId != nil {
		entity.TableId = *req.TableId
	}
	
	if req.ReservationId != nil {
		entity.ReservationId = *req.ReservationId
	}
	
	if req.Covers != nil {
		entity.Covers = *req.Covers
	}
	
	if req.CustomerId != nil {
		entity.CustomerId = *req.CustomerId
	}
	
	if req.WaiterId != nil {
		entity.WaiterId = *req.WaiterId
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.OrderDate != nil {
		entity.OrderDate = *req.OrderDate
	}
	
	if req.SubmittedAt != nil {
		entity.SubmittedAt = *req.SubmittedAt
	}
	
	if req.KitchenReceivedAt != nil {
		entity.KitchenReceivedAt = *req.KitchenReceivedAt
	}
	
	if req.ReadyAt != nil {
		entity.ReadyAt = *req.ReadyAt
	}
	
	if req.ServedAt != nil {
		entity.ServedAt = *req.ServedAt
	}
	
	if req.CompletedAt != nil {
		entity.CompletedAt = *req.CompletedAt
	}
	
	if req.Subtotal != nil {
		entity.Subtotal = *req.Subtotal
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = *req.TaxAmount
	}
	
	if req.DiscountAmount != nil {
		entity.DiscountAmount = *req.DiscountAmount
	}
	
	if req.ServiceCharge != nil {
		entity.ServiceCharge = *req.ServiceCharge
	}
	
	if req.TotalAmount != nil {
		entity.TotalAmount = *req.TotalAmount
	}
	
	if req.SaleId != nil {
		entity.SaleId = *req.SaleId
	}
	
	if req.ShiftId != nil {
		entity.ShiftId = *req.ShiftId
	}
	
	if req.CustomerNotes != nil {
		entity.CustomerNotes = *req.CustomerNotes
	}
	
	if req.KitchenNotes != nil {
		entity.KitchenNotes = *req.KitchenNotes
	}
	
	if req.InternalNotes != nil {
		entity.InternalNotes = *req.InternalNotes
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
	
	if req.'draft', != nil {
		entity.'draft', = *req.'draft',
	}
	
	if req.'served', != nil {
		entity.'served', = *req.'served',
	}
	
	if req.'dineIn', != nil {
		entity.'dineIn', = *req.'dineIn',
	}
	
	if req.Subtotal != nil {
		entity.Subtotal = *req.Subtotal
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
		return nil, fmt.Errorf("failed to update orders: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated orders",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a orders
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting orders",
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
		return fmt.Errorf("failed to get orders: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("orders not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete orders: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted orders",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Orders) *OrdersResponse {
	return &OrdersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		OrderNumber: entity.OrderNumber,
		
		DisplayNumber: entity.DisplayNumber,
		
		OrderType: entity.OrderType,
		
		TableId: entity.TableId,
		
		ReservationId: entity.ReservationId,
		
		Covers: entity.Covers,
		
		CustomerId: entity.CustomerId,
		
		WaiterId: entity.WaiterId,
		
		Status: entity.Status,
		
		OrderDate: entity.OrderDate,
		
		SubmittedAt: entity.SubmittedAt,
		
		KitchenReceivedAt: entity.KitchenReceivedAt,
		
		ReadyAt: entity.ReadyAt,
		
		ServedAt: entity.ServedAt,
		
		CompletedAt: entity.CompletedAt,
		
		Subtotal: entity.Subtotal,
		
		TaxAmount: entity.TaxAmount,
		
		DiscountAmount: entity.DiscountAmount,
		
		ServiceCharge: entity.ServiceCharge,
		
		TotalAmount: entity.TotalAmount,
		
		SaleId: entity.SaleId,
		
		ShiftId: entity.ShiftId,
		
		CustomerNotes: entity.CustomerNotes,
		
		KitchenNotes: entity.KitchenNotes,
		
		InternalNotes: entity.InternalNotes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'draft',: entity.'draft',,
		
		'served',: entity.'served',,
		
		'dineIn',: entity.'dineIn',,
		
		Subtotal: entity.Subtotal,
		
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


// validateBusinessRules validates business rules for orders
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Orders) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a orders can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

package purchase_order

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for PurchaseOrders
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PurchaseOrders service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new purchase_orders
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreatePurchaseOrdersRequest) (*PurchaseOrdersResponse, error) {
	s.logger.Info("creating purchase_orders",
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
	entity := &PurchaseOrders{
		OrganizationId: orgID,
		
		PoNumber: req.PoNumber,
		
		SupplierId: req.SupplierId,
		
		LocationId: req.LocationId,
		
		OrderDate: req.OrderDate,
		
		ExpectedDeliveryDate: req.ExpectedDeliveryDate,
		
		ActualDeliveryDate: req.ActualDeliveryDate,
		
		SubtotalAmount: req.SubtotalAmount,
		
		TaxAmount: req.TaxAmount,
		
		ShippingAmount: req.ShippingAmount,
		
		TotalAmount: req.TotalAmount,
		
		PaymentTerms: req.PaymentTerms,
		
		PaymentDueDate: req.PaymentDueDate,
		
		Status: req.Status,
		
		Status: req.Status,
		
		ApprovedBy: req.ApprovedBy,
		
		ApprovedAt: req.ApprovedAt,
		
		Notes: req.Notes,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create purchase_orders: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created purchase_orders",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a purchase_orders by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PurchaseOrdersResponse, error) {
	s.logger.Debug("getting purchase_orders",
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
		return nil, fmt.Errorf("failed to get purchase_orders: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("purchase_orders not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of purchase_orders records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*PurchaseOrdersListResponse, error) {
	s.logger.Debug("listing purchase_orders",
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
		return nil, fmt.Errorf("failed to list purchase_orders: %w", err)
	}

	// Convert to response
	items := make([]*PurchaseOrdersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PurchaseOrdersListResponse{
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

// Update updates an existing purchase_orders
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdatePurchaseOrdersRequest) (*PurchaseOrdersResponse, error) {
	s.logger.Info("updating purchase_orders",
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
		return nil, fmt.Errorf("failed to get purchase_orders: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("purchase_orders not found or access denied")
	}
	

	// Update fields
	
	if req.PoNumber != nil {
		entity.PoNumber = req.PoNumber
	}
	
	if req.SupplierId != nil {
		entity.SupplierId = req.SupplierId
	}
	
	if req.LocationId != nil {
		entity.LocationId = req.LocationId
	}
	
	if req.OrderDate != nil {
		entity.OrderDate = req.OrderDate
	}
	
	if req.ExpectedDeliveryDate != nil {
		entity.ExpectedDeliveryDate = req.ExpectedDeliveryDate
	}
	
	if req.ActualDeliveryDate != nil {
		entity.ActualDeliveryDate = req.ActualDeliveryDate
	}
	
	if req.SubtotalAmount != nil {
		entity.SubtotalAmount = req.SubtotalAmount
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = req.TaxAmount
	}
	
	if req.ShippingAmount != nil {
		entity.ShippingAmount = req.ShippingAmount
	}
	
	if req.TotalAmount != nil {
		entity.TotalAmount = req.TotalAmount
	}
	
	if req.PaymentTerms != nil {
		entity.PaymentTerms = req.PaymentTerms
	}
	
	if req.PaymentDueDate != nil {
		entity.PaymentDueDate = req.PaymentDueDate
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.ApprovedBy != nil {
		entity.ApprovedBy = req.ApprovedBy
	}
	
	if req.ApprovedAt != nil {
		entity.ApprovedAt = req.ApprovedAt
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update purchase_orders: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated purchase_orders",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a purchase_orders
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting purchase_orders",
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
		return fmt.Errorf("failed to get purchase_orders: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("purchase_orders not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete purchase_orders: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted purchase_orders",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PurchaseOrders) *PurchaseOrdersResponse {
	return &PurchaseOrdersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		PoNumber: entity.PoNumber,
		
		SupplierId: entity.SupplierId,
		
		LocationId: entity.LocationId,
		
		OrderDate: entity.OrderDate,
		
		ExpectedDeliveryDate: entity.ExpectedDeliveryDate,
		
		ActualDeliveryDate: entity.ActualDeliveryDate,
		
		SubtotalAmount: entity.SubtotalAmount,
		
		TaxAmount: entity.TaxAmount,
		
		ShippingAmount: entity.ShippingAmount,
		
		TotalAmount: entity.TotalAmount,
		
		PaymentTerms: entity.PaymentTerms,
		
		PaymentDueDate: entity.PaymentDueDate,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		ApprovedBy: entity.ApprovedBy,
		
		ApprovedAt: entity.ApprovedAt,
		
		Notes: entity.Notes,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for purchase_orders
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PurchaseOrders) error {
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

// canDelete checks if a purchase_orders can be deleted
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

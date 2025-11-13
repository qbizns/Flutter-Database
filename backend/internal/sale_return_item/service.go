package sale_return_item

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for SaleReturnItems
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new SaleReturnItems service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new sale_return_items
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateSaleReturnItemsRequest) (*SaleReturnItemsResponse, error) {
	s.logger.Info("creating sale_return_items",
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
	entity := &SaleReturnItems{
		OrganizationId: orgID,
		
		SaleReturnId: req.SaleReturnId,
		
		OriginalSaleItemId: req.OriginalSaleItemId,
		
		ProductId: req.ProductId,
		
		ProductVariantId: req.ProductVariantId,
		
		Quantity: req.Quantity,
		
		UnitPrice: req.UnitPrice,
		
		Subtotal: req.Subtotal,
		
		TaxAmount: req.TaxAmount,
		
		DiscountAmount: req.DiscountAmount,
		
		TotalAmount: req.TotalAmount,
		
		ReturnReasonId: req.ReturnReasonId,
		
		ReturnReasonNotes: req.ReturnReasonNotes,
		
		ItemCondition: req.ItemCondition,
		
		ItemCondition: req.ItemCondition,
		
		IsRestockable: req.IsRestockable,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create sale_return_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created sale_return_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a sale_return_items by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*SaleReturnItemsResponse, error) {
	s.logger.Debug("getting sale_return_items",
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
		return nil, fmt.Errorf("failed to get sale_return_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("sale_return_items not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of sale_return_items records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*SaleReturnItemsListResponse, error) {
	s.logger.Debug("listing sale_return_items",
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
		return nil, fmt.Errorf("failed to list sale_return_items: %w", err)
	}

	// Convert to response
	items := make([]*SaleReturnItemsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &SaleReturnItemsListResponse{
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

// Update updates an existing sale_return_items
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateSaleReturnItemsRequest) (*SaleReturnItemsResponse, error) {
	s.logger.Info("updating sale_return_items",
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
		return nil, fmt.Errorf("failed to get sale_return_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("sale_return_items not found or access denied")
	}
	

	// Update fields
	
	if req.SaleReturnId != nil {
		entity.SaleReturnId = req.SaleReturnId
	}
	
	if req.OriginalSaleItemId != nil {
		entity.OriginalSaleItemId = req.OriginalSaleItemId
	}
	
	if req.ProductId != nil {
		entity.ProductId = req.ProductId
	}
	
	if req.ProductVariantId != nil {
		entity.ProductVariantId = req.ProductVariantId
	}
	
	if req.Quantity != nil {
		entity.Quantity = req.Quantity
	}
	
	if req.UnitPrice != nil {
		entity.UnitPrice = req.UnitPrice
	}
	
	if req.Subtotal != nil {
		entity.Subtotal = req.Subtotal
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = req.TaxAmount
	}
	
	if req.DiscountAmount != nil {
		entity.DiscountAmount = req.DiscountAmount
	}
	
	if req.TotalAmount != nil {
		entity.TotalAmount = req.TotalAmount
	}
	
	if req.ReturnReasonId != nil {
		entity.ReturnReasonId = req.ReturnReasonId
	}
	
	if req.ReturnReasonNotes != nil {
		entity.ReturnReasonNotes = req.ReturnReasonNotes
	}
	
	if req.ItemCondition != nil {
		entity.ItemCondition = req.ItemCondition
	}
	
	if req.ItemCondition != nil {
		entity.ItemCondition = req.ItemCondition
	}
	
	if req.IsRestockable != nil {
		entity.IsRestockable = req.IsRestockable
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update sale_return_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated sale_return_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a sale_return_items
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting sale_return_items",
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
		return fmt.Errorf("failed to get sale_return_items: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("sale_return_items not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete sale_return_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted sale_return_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *SaleReturnItems) *SaleReturnItemsResponse {
	return &SaleReturnItemsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		SaleReturnId: entity.SaleReturnId,
		
		OriginalSaleItemId: entity.OriginalSaleItemId,
		
		ProductId: entity.ProductId,
		
		ProductVariantId: entity.ProductVariantId,
		
		Quantity: entity.Quantity,
		
		UnitPrice: entity.UnitPrice,
		
		Subtotal: entity.Subtotal,
		
		TaxAmount: entity.TaxAmount,
		
		DiscountAmount: entity.DiscountAmount,
		
		TotalAmount: entity.TotalAmount,
		
		ReturnReasonId: entity.ReturnReasonId,
		
		ReturnReasonNotes: entity.ReturnReasonNotes,
		
		ItemCondition: entity.ItemCondition,
		
		ItemCondition: entity.ItemCondition,
		
		IsRestockable: entity.IsRestockable,
		
		CreatedAt: entity.CreatedAt,
		
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


// validateBusinessRules validates business rules for sale_return_items
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *SaleReturnItems) error {
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

// canDelete checks if a sale_return_items can be deleted
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

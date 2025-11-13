package sale_item

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for SaleItems
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new SaleItems service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new sale_items
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateSaleItemsRequest) (*SaleItemsResponse, error) {
	s.logger.Info("creating sale_items",
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
	entity := &SaleItems{
		OrganizationID: orgID,
		
		SaleId: req.SaleId,
		
		ProductId: req.ProductId,
		
		ProductName: req.ProductName,
		
		ProductSku: req.ProductSku,
		
		Quantity: req.Quantity,
		
		Unit: req.Unit,
		
		UnitPrice: req.UnitPrice,
		
		CostPrice: req.CostPrice,
		
		Subtotal: req.Subtotal,
		
		TaxRate: req.TaxRate,
		
		TaxAmount: req.TaxAmount,
		
		DiscountAmount: req.DiscountAmount,
		
		Total: req.Total,
		
		DiscountType: req.DiscountType,
		
		DiscountValue: req.DiscountValue,
		
		Notes: req.Notes,
		
		CustomFields: req.CustomFields,
		
		Metadata: req.Metadata,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create sale_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created sale_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a sale_items by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*SaleItemsResponse, error) {
	s.logger.Debug("getting sale_items",
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
		return nil, fmt.Errorf("failed to get sale_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("sale_items not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of sale_items records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*SaleItemsListResponse, error) {
	s.logger.Debug("listing sale_items",
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
		return nil, fmt.Errorf("failed to list sale_items: %w", err)
	}

	// Convert to response
	items := make([]*SaleItemsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &SaleItemsListResponse{
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

// Update updates an existing sale_items
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateSaleItemsRequest) (*SaleItemsResponse, error) {
	s.logger.Info("updating sale_items",
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
		return nil, fmt.Errorf("failed to get sale_items: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("sale_items not found or access denied")
	}
	

	// Update fields
	
	if req.SaleId != nil {
		entity.SaleId = req.SaleId
	}
	
	if req.ProductId != nil {
		entity.ProductId = req.ProductId
	}
	
	if req.ProductName != nil {
		entity.ProductName = req.ProductName
	}
	
	if req.ProductSku != nil {
		entity.ProductSku = req.ProductSku
	}
	
	if req.Quantity != nil {
		entity.Quantity = req.Quantity
	}
	
	if req.Unit != nil {
		entity.Unit = req.Unit
	}
	
	if req.UnitPrice != nil {
		entity.UnitPrice = req.UnitPrice
	}
	
	if req.CostPrice != nil {
		entity.CostPrice = req.CostPrice
	}
	
	if req.Subtotal != nil {
		entity.Subtotal = req.Subtotal
	}
	
	if req.TaxRate != nil {
		entity.TaxRate = req.TaxRate
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = req.TaxAmount
	}
	
	if req.DiscountAmount != nil {
		entity.DiscountAmount = req.DiscountAmount
	}
	
	if req.Total != nil {
		entity.Total = req.Total
	}
	
	if req.DiscountType != nil {
		entity.DiscountType = req.DiscountType
	}
	
	if req.DiscountValue != nil {
		entity.DiscountValue = req.DiscountValue
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.CustomFields != nil {
		entity.CustomFields = req.CustomFields
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update sale_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated sale_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a sale_items
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting sale_items",
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
		return fmt.Errorf("failed to get sale_items: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("sale_items not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete sale_items: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted sale_items",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *SaleItems) *SaleItemsResponse {
	return &SaleItemsResponse{
		
		Id: entity.Id,
		
		SaleId: entity.SaleId,
		
		OrganizationId: entity.OrganizationId,
		
		ProductId: entity.ProductId,
		
		ProductName: entity.ProductName,
		
		ProductSku: entity.ProductSku,
		
		Quantity: entity.Quantity,
		
		Unit: entity.Unit,
		
		UnitPrice: entity.UnitPrice,
		
		CostPrice: entity.CostPrice,
		
		Subtotal: entity.Subtotal,
		
		TaxRate: entity.TaxRate,
		
		TaxAmount: entity.TaxAmount,
		
		DiscountAmount: entity.DiscountAmount,
		
		Total: entity.Total,
		
		DiscountType: entity.DiscountType,
		
		DiscountValue: entity.DiscountValue,
		
		Notes: entity.Notes,
		
		CustomFields: entity.CustomFields,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
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


// validateBusinessRules validates business rules for sale_items
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *SaleItems) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a sale_items can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

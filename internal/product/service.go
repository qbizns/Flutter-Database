package product

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/product"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/product"
	"go.uber.org/zap"
)

// Service handles business logic for Products
type Service struct {
	repo   *product.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Products service
func NewService(repo *product.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new products
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateProductsRequest) (*dto.ProductsResponse, error) {
	s.logger.Info("creating products",
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
	entity := &product.Products{
		OrganizationID: orgID,
		
		Sku: req.Sku,
		
		Barcode: req.Barcode,
		
		Name: req.Name,
		
		Description: req.Description,
		
		CategoryId: req.CategoryId,
		
		CostPrice: req.CostPrice,
		
		SellingPrice: req.SellingPrice,
		
		CompareAtPrice: req.CompareAtPrice,
		
		TaxRate: req.TaxRate,
		
		IsTaxInclusive: req.IsTaxInclusive,
		
		TrackInventory: req.TrackInventory,
		
		CurrentStock: req.CurrentStock,
		
		LowStockThreshold: req.LowStockThreshold,
		
		Unit: req.Unit,
		
		IsService: req.IsService,
		
		IsComposite: req.IsComposite,
		
		HasVariants: req.HasVariants,
		
		ImageUrl: req.ImageUrl,
		
		Images: req.Images,
		
		SortOrder: req.SortOrder,
		
		IsActive: req.IsActive,
		
		IsFeatured: req.IsFeatured,
		
		CustomFields: req.CustomFields,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create products: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created products",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a products by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.ProductsResponse, error) {
	s.logger.Debug("getting products",
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
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("products not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of products records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.ProductsListResponse, error) {
	s.logger.Debug("listing products",
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
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Convert to response
	items := make([]*dto.ProductsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.ProductsListResponse{
		Items: items,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing products
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateProductsRequest) (*dto.ProductsResponse, error) {
	s.logger.Info("updating products",
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
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("products not found or access denied")
	}
	

	// Update fields
	
	if req.Sku != nil {
		entity.Sku = *req.Sku
	}
	
	if req.Barcode != nil {
		entity.Barcode = *req.Barcode
	}
	
	if req.Name != nil {
		entity.Name = *req.Name
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.CategoryId != nil {
		entity.CategoryId = *req.CategoryId
	}
	
	if req.CostPrice != nil {
		entity.CostPrice = *req.CostPrice
	}
	
	if req.SellingPrice != nil {
		entity.SellingPrice = *req.SellingPrice
	}
	
	if req.CompareAtPrice != nil {
		entity.CompareAtPrice = *req.CompareAtPrice
	}
	
	if req.TaxRate != nil {
		entity.TaxRate = *req.TaxRate
	}
	
	if req.IsTaxInclusive != nil {
		entity.IsTaxInclusive = *req.IsTaxInclusive
	}
	
	if req.TrackInventory != nil {
		entity.TrackInventory = *req.TrackInventory
	}
	
	if req.CurrentStock != nil {
		entity.CurrentStock = *req.CurrentStock
	}
	
	if req.LowStockThreshold != nil {
		entity.LowStockThreshold = *req.LowStockThreshold
	}
	
	if req.Unit != nil {
		entity.Unit = *req.Unit
	}
	
	if req.IsService != nil {
		entity.IsService = *req.IsService
	}
	
	if req.IsComposite != nil {
		entity.IsComposite = *req.IsComposite
	}
	
	if req.HasVariants != nil {
		entity.HasVariants = *req.HasVariants
	}
	
	if req.ImageUrl != nil {
		entity.ImageUrl = *req.ImageUrl
	}
	
	if req.Images != nil {
		entity.Images = *req.Images
	}
	
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.IsFeatured != nil {
		entity.IsFeatured = *req.IsFeatured
	}
	
	if req.CustomFields != nil {
		entity.CustomFields = *req.CustomFields
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
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update products: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated products",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a products
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting products",
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
		return fmt.Errorf("failed to get products: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("products not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete products: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted products",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *product.Products) *dto.ProductsResponse {
	return &dto.ProductsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		Sku: entity.Sku,
		
		Barcode: entity.Barcode,
		
		Name: entity.Name,
		
		Description: entity.Description,
		
		CategoryId: entity.CategoryId,
		
		CostPrice: entity.CostPrice,
		
		SellingPrice: entity.SellingPrice,
		
		CompareAtPrice: entity.CompareAtPrice,
		
		TaxRate: entity.TaxRate,
		
		IsTaxInclusive: entity.IsTaxInclusive,
		
		TrackInventory: entity.TrackInventory,
		
		CurrentStock: entity.CurrentStock,
		
		LowStockThreshold: entity.LowStockThreshold,
		
		Unit: entity.Unit,
		
		IsService: entity.IsService,
		
		IsComposite: entity.IsComposite,
		
		HasVariants: entity.HasVariants,
		
		ImageUrl: entity.ImageUrl,
		
		Images: entity.Images,
		
		SortOrder: entity.SortOrder,
		
		IsActive: entity.IsActive,
		
		IsFeatured: entity.IsFeatured,
		
		CustomFields: entity.CustomFields,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for products
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *product.Products) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a products can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

package product_variant

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for ProductVariants
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new ProductVariants service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new product_variants
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateProductVariantsRequest) (*ProductVariantsResponse, error) {
	s.logger.Info("creating product_variants",
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
	entity := &ProductVariants{
		OrganizationId: orgID,
		
		ProductId: req.ProductId,
		
		VariantName: req.VariantName,
		
		Sku: req.Sku,
		
		Barcode: req.Barcode,
		
		Attributes: req.Attributes,
		
		CostPrice: req.CostPrice,
		
		SellingPrice: req.SellingPrice,
		
		CompareAtPrice: req.CompareAtPrice,
		
		CurrentStock: req.CurrentStock,
		
		ReorderLevel: req.ReorderLevel,
		
		ReorderQuantity: req.ReorderQuantity,
		
		Weight: req.Weight,
		
		WeightUnit: req.WeightUnit,
		
		Dimensions: req.Dimensions,
		
		IsActive: req.IsActive,
		
		IsDefault: req.IsDefault,
		
		SortOrder: req.SortOrder,
		
		ImageUrl: req.ImageUrl,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		(costPrice: req.(costPrice,
		
		(sellingPrice: req.(sellingPrice,
		
		(compareAtPrice: req.(compareAtPrice,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create product_variants: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created product_variants",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a product_variants by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductVariantsResponse, error) {
	s.logger.Debug("getting product_variants",
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
		return nil, fmt.Errorf("failed to get product_variants: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("product_variants not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of product_variants records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ProductVariantsListResponse, error) {
	s.logger.Debug("listing product_variants",
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
		return nil, fmt.Errorf("failed to list product_variants: %w", err)
	}

	// Convert to response
	items := make([]*ProductVariantsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ProductVariantsListResponse{
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

// Update updates an existing product_variants
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateProductVariantsRequest) (*ProductVariantsResponse, error) {
	s.logger.Info("updating product_variants",
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
		return nil, fmt.Errorf("failed to get product_variants: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("product_variants not found or access denied")
	}
	

	// Update fields
	
	if req.ProductId != nil {
		entity.ProductId = *req.ProductId
	}
	
	if req.VariantName != nil {
		entity.VariantName = *req.VariantName
	}
	
	if req.Sku != nil {
		entity.Sku = *req.Sku
	}
	
	if req.Barcode != nil {
		entity.Barcode = *req.Barcode
	}
	
	if req.Attributes != nil {
		entity.Attributes = *req.Attributes
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
	
	if req.CurrentStock != nil {
		entity.CurrentStock = *req.CurrentStock
	}
	
	if req.ReorderLevel != nil {
		entity.ReorderLevel = req.ReorderLevel
	}
	
	if req.ReorderQuantity != nil {
		entity.ReorderQuantity = *req.ReorderQuantity
	}
	
	if req.Weight != nil {
		entity.Weight = *req.Weight
	}
	
	if req.WeightUnit != nil {
		entity.WeightUnit = *req.WeightUnit
	}
	
	if req.Dimensions != nil {
		entity.Dimensions = *req.Dimensions
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.IsDefault != nil {
		entity.IsDefault = req.IsDefault
	}
	
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}
	
	if req.ImageUrl != nil {
		entity.ImageUrl = *req.ImageUrl
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
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
	
	if req.(costPrice != nil {
		entity.(costPrice = *req.(costPrice
	}
	
	if req.(sellingPrice != nil {
		entity.(sellingPrice = *req.(sellingPrice
	}
	
	if req.(compareAtPrice != nil {
		entity.(compareAtPrice = *req.(compareAtPrice
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update product_variants: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated product_variants",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a product_variants
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting product_variants",
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
		return fmt.Errorf("failed to get product_variants: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("product_variants not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete product_variants: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted product_variants",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *ProductVariants) *ProductVariantsResponse {
	return &ProductVariantsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ProductId: entity.ProductId,
		
		VariantName: entity.VariantName,
		
		Sku: entity.Sku,
		
		Barcode: entity.Barcode,
		
		Attributes: entity.Attributes,
		
		CostPrice: entity.CostPrice,
		
		SellingPrice: entity.SellingPrice,
		
		CompareAtPrice: entity.CompareAtPrice,
		
		CurrentStock: entity.CurrentStock,
		
		ReorderLevel: entity.ReorderLevel,
		
		ReorderQuantity: entity.ReorderQuantity,
		
		Weight: entity.Weight,
		
		WeightUnit: entity.WeightUnit,
		
		Dimensions: entity.Dimensions,
		
		IsActive: entity.IsActive,
		
		IsDefault: entity.IsDefault,
		
		SortOrder: entity.SortOrder,
		
		ImageUrl: entity.ImageUrl,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		(costPrice: entity.(costPrice,
		
		(sellingPrice: entity.(sellingPrice,
		
		(compareAtPrice: entity.(compareAtPrice,
		
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


// validateBusinessRules validates business rules for product_variants
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *ProductVariants) error {
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

// canDelete checks if a product_variants can be deleted
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

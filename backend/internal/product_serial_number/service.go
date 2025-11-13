package product_serial_number

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for ProductSerialNumbers
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new ProductSerialNumbers service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new product_serial_numbers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateProductSerialNumbersRequest) (*ProductSerialNumbersResponse, error) {
	s.logger.Info("creating product_serial_numbers",
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
	entity := &ProductSerialNumbers{
		OrganizationID: orgID,
		
		ProductId: req.ProductId,
		
		ProductVariantId: req.ProductVariantId,
		
		LocationId: req.LocationId,
		
		SerialNumber: req.SerialNumber,
		
		Status: req.Status,
		
		PurchaseOrderId: req.PurchaseOrderId,
		
		PurchaseDate: req.PurchaseDate,
		
		PurchaseCost: req.PurchaseCost,
		
		SupplierId: req.SupplierId,
		
		SaleId: req.SaleId,
		
		SaleDate: req.SaleDate,
		
		SalePrice: req.SalePrice,
		
		CustomerId: req.CustomerId,
		
		WarrantyStartDate: req.WarrantyStartDate,
		
		WarrantyEndDate: req.WarrantyEndDate,
		
		WarrantyProvider: req.WarrantyProvider,
		
		WarrantyTerms: req.WarrantyTerms,
		
		Notes: req.Notes,
		
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
		return nil, fmt.Errorf("failed to create product_serial_numbers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created product_serial_numbers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a product_serial_numbers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductSerialNumbersResponse, error) {
	s.logger.Debug("getting product_serial_numbers",
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
		return nil, fmt.Errorf("failed to get product_serial_numbers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("product_serial_numbers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of product_serial_numbers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ProductSerialNumbersListResponse, error) {
	s.logger.Debug("listing product_serial_numbers",
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
		return nil, fmt.Errorf("failed to list product_serial_numbers: %w", err)
	}

	// Convert to response
	items := make([]*ProductSerialNumbersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ProductSerialNumbersListResponse{
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

// Update updates an existing product_serial_numbers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateProductSerialNumbersRequest) (*ProductSerialNumbersResponse, error) {
	s.logger.Info("updating product_serial_numbers",
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
		return nil, fmt.Errorf("failed to get product_serial_numbers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("product_serial_numbers not found or access denied")
	}
	

	// Update fields
	
	if req.ProductId != nil {
		entity.ProductId = req.ProductId
	}
	
	if req.ProductVariantId != nil {
		entity.ProductVariantId = req.ProductVariantId
	}
	
	if req.LocationId != nil {
		entity.LocationId = req.LocationId
	}
	
	if req.SerialNumber != nil {
		entity.SerialNumber = req.SerialNumber
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.PurchaseOrderId != nil {
		entity.PurchaseOrderId = req.PurchaseOrderId
	}
	
	if req.PurchaseDate != nil {
		entity.PurchaseDate = req.PurchaseDate
	}
	
	if req.PurchaseCost != nil {
		entity.PurchaseCost = req.PurchaseCost
	}
	
	if req.SupplierId != nil {
		entity.SupplierId = req.SupplierId
	}
	
	if req.SaleId != nil {
		entity.SaleId = req.SaleId
	}
	
	if req.SaleDate != nil {
		entity.SaleDate = req.SaleDate
	}
	
	if req.SalePrice != nil {
		entity.SalePrice = req.SalePrice
	}
	
	if req.CustomerId != nil {
		entity.CustomerId = req.CustomerId
	}
	
	if req.WarrantyStartDate != nil {
		entity.WarrantyStartDate = req.WarrantyStartDate
	}
	
	if req.WarrantyEndDate != nil {
		entity.WarrantyEndDate = req.WarrantyEndDate
	}
	
	if req.WarrantyProvider != nil {
		entity.WarrantyProvider = req.WarrantyProvider
	}
	
	if req.WarrantyTerms != nil {
		entity.WarrantyTerms = req.WarrantyTerms
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
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update product_serial_numbers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated product_serial_numbers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a product_serial_numbers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting product_serial_numbers",
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
		return fmt.Errorf("failed to get product_serial_numbers: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("product_serial_numbers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete product_serial_numbers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted product_serial_numbers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *ProductSerialNumbers) *ProductSerialNumbersResponse {
	return &ProductSerialNumbersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ProductId: entity.ProductId,
		
		ProductVariantId: entity.ProductVariantId,
		
		LocationId: entity.LocationId,
		
		SerialNumber: entity.SerialNumber,
		
		Status: entity.Status,
		
		PurchaseOrderId: entity.PurchaseOrderId,
		
		PurchaseDate: entity.PurchaseDate,
		
		PurchaseCost: entity.PurchaseCost,
		
		SupplierId: entity.SupplierId,
		
		SaleId: entity.SaleId,
		
		SaleDate: entity.SaleDate,
		
		SalePrice: entity.SalePrice,
		
		CustomerId: entity.CustomerId,
		
		WarrantyStartDate: entity.WarrantyStartDate,
		
		WarrantyEndDate: entity.WarrantyEndDate,
		
		WarrantyProvider: entity.WarrantyProvider,
		
		WarrantyTerms: entity.WarrantyTerms,
		
		Notes: entity.Notes,
		
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


// validateBusinessRules validates business rules for product_serial_numbers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *ProductSerialNumbers) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a product_serial_numbers can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

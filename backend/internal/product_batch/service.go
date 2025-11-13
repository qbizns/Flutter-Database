package product_batch

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for ProductBatches
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new ProductBatches service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new product_batches
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateProductBatchesRequest) (*ProductBatchesResponse, error) {
	s.logger.Info("creating product_batches",
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
	entity := &ProductBatches{
		OrganizationId: orgID,
		
		ProductId: req.ProductId,
		
		ProductVariantId: req.ProductVariantId,
		
		LocationId: req.LocationId,
		
		BatchNumber: req.BatchNumber,
		
		LotNumber: req.LotNumber,
		
		Status: req.Status,
		
		InitialQuantity: req.InitialQuantity,
		
		CurrentQuantity: req.CurrentQuantity,
		
		UnitOfMeasure: req.UnitOfMeasure,
		
		ManufacturingDate: req.ManufacturingDate,
		
		ExpirationDate: req.ExpirationDate,
		
		ReceivedDate: req.ReceivedDate,
		
		PurchaseOrderId: req.PurchaseOrderId,
		
		SupplierId: req.SupplierId,
		
		SupplierBatchNumber: req.SupplierBatchNumber,
		
		UnitCost: req.UnitCost,
		
		TotalCost: req.TotalCost,
		
		QualityStatus: req.QualityStatus,
		
		QualityCheckDate: req.QualityCheckDate,
		
		QualityCheckedBy: req.QualityCheckedBy,
		
		QualityNotes: req.QualityNotes,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		InitialQuantity: req.InitialQuantity,
		
		CurrentQuantity: req.CurrentQuantity,
		
		CurrentQuantity: req.CurrentQuantity,
		
		ExpirationDate: req.ExpirationDate,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create product_batches: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created product_batches",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a product_batches by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductBatchesResponse, error) {
	s.logger.Debug("getting product_batches",
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
		return nil, fmt.Errorf("failed to get product_batches: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("product_batches not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of product_batches records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ProductBatchesListResponse, error) {
	s.logger.Debug("listing product_batches",
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
		return nil, fmt.Errorf("failed to list product_batches: %w", err)
	}

	// Convert to response
	items := make([]*ProductBatchesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ProductBatchesListResponse{
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

// Update updates an existing product_batches
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateProductBatchesRequest) (*ProductBatchesResponse, error) {
	s.logger.Info("updating product_batches",
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
		return nil, fmt.Errorf("failed to get product_batches: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("product_batches not found or access denied")
	}
	

	// Update fields
	
	if req.ProductId != nil {
		entity.ProductId = *req.ProductId
	}
	
	if req.ProductVariantId != nil {
		entity.ProductVariantId = *req.ProductVariantId
	}
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.BatchNumber != nil {
		entity.BatchNumber = *req.BatchNumber
	}
	
	if req.LotNumber != nil {
		entity.LotNumber = *req.LotNumber
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.InitialQuantity != nil {
		entity.InitialQuantity = *req.InitialQuantity
	}
	
	if req.CurrentQuantity != nil {
		entity.CurrentQuantity = *req.CurrentQuantity
	}
	
	if req.UnitOfMeasure != nil {
		entity.UnitOfMeasure = *req.UnitOfMeasure
	}
	
	if req.ManufacturingDate != nil {
		entity.ManufacturingDate = *req.ManufacturingDate
	}
	
	if req.ExpirationDate != nil {
		entity.ExpirationDate = *req.ExpirationDate
	}
	
	if req.ReceivedDate != nil {
		entity.ReceivedDate = *req.ReceivedDate
	}
	
	if req.PurchaseOrderId != nil {
		entity.PurchaseOrderId = *req.PurchaseOrderId
	}
	
	if req.SupplierId != nil {
		entity.SupplierId = *req.SupplierId
	}
	
	if req.SupplierBatchNumber != nil {
		entity.SupplierBatchNumber = *req.SupplierBatchNumber
	}
	
	if req.UnitCost != nil {
		entity.UnitCost = *req.UnitCost
	}
	
	if req.TotalCost != nil {
		entity.TotalCost = *req.TotalCost
	}
	
	if req.QualityStatus != nil {
		entity.QualityStatus = *req.QualityStatus
	}
	
	if req.QualityCheckDate != nil {
		entity.QualityCheckDate = *req.QualityCheckDate
	}
	
	if req.QualityCheckedBy != nil {
		entity.QualityCheckedBy = *req.QualityCheckedBy
	}
	
	if req.QualityNotes != nil {
		entity.QualityNotes = *req.QualityNotes
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
	
	if req.InitialQuantity != nil {
		entity.InitialQuantity = *req.InitialQuantity
	}
	
	if req.CurrentQuantity != nil {
		entity.CurrentQuantity = *req.CurrentQuantity
	}
	
	if req.CurrentQuantity != nil {
		entity.CurrentQuantity = *req.CurrentQuantity
	}
	
	if req.ExpirationDate != nil {
		entity.ExpirationDate = *req.ExpirationDate
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update product_batches: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated product_batches",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a product_batches
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting product_batches",
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
		return fmt.Errorf("failed to get product_batches: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("product_batches not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete product_batches: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted product_batches",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *ProductBatches) *ProductBatchesResponse {
	return &ProductBatchesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ProductId: entity.ProductId,
		
		ProductVariantId: entity.ProductVariantId,
		
		LocationId: entity.LocationId,
		
		BatchNumber: entity.BatchNumber,
		
		LotNumber: entity.LotNumber,
		
		Status: entity.Status,
		
		InitialQuantity: entity.InitialQuantity,
		
		CurrentQuantity: entity.CurrentQuantity,
		
		UnitOfMeasure: entity.UnitOfMeasure,
		
		ManufacturingDate: entity.ManufacturingDate,
		
		ExpirationDate: entity.ExpirationDate,
		
		ReceivedDate: entity.ReceivedDate,
		
		PurchaseOrderId: entity.PurchaseOrderId,
		
		SupplierId: entity.SupplierId,
		
		SupplierBatchNumber: entity.SupplierBatchNumber,
		
		UnitCost: entity.UnitCost,
		
		TotalCost: entity.TotalCost,
		
		QualityStatus: entity.QualityStatus,
		
		QualityCheckDate: entity.QualityCheckDate,
		
		QualityCheckedBy: entity.QualityCheckedBy,
		
		QualityNotes: entity.QualityNotes,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		InitialQuantity: entity.InitialQuantity,
		
		CurrentQuantity: entity.CurrentQuantity,
		
		CurrentQuantity: entity.CurrentQuantity,
		
		ExpirationDate: entity.ExpirationDate,
		
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


// validateBusinessRules validates business rules for product_batches
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *ProductBatches) error {
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

// canDelete checks if a product_batches can be deleted
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

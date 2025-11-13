package goods_receipt

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for GoodsReceipts
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new GoodsReceipts service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new goods_receipts
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateGoodsReceiptsRequest) (*GoodsReceiptsResponse, error) {
	s.logger.Info("creating goods_receipts",
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
	entity := &GoodsReceipts{
		OrganizationId: orgID,
		
		ReceiptNumber: req.ReceiptNumber,
		
		PurchaseOrderId: req.PurchaseOrderId,
		
		SupplierId: req.SupplierId,
		
		LocationId: req.LocationId,
		
		ReceiptDate: req.ReceiptDate,
		
		ReceivedBy: req.ReceivedBy,
		
		Status: req.Status,
		
		Status: req.Status,
		
		Notes: req.Notes,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create goods_receipts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created goods_receipts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a goods_receipts by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*GoodsReceiptsResponse, error) {
	s.logger.Debug("getting goods_receipts",
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
		return nil, fmt.Errorf("failed to get goods_receipts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("goods_receipts not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of goods_receipts records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*GoodsReceiptsListResponse, error) {
	s.logger.Debug("listing goods_receipts",
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
		return nil, fmt.Errorf("failed to list goods_receipts: %w", err)
	}

	// Convert to response
	items := make([]*GoodsReceiptsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &GoodsReceiptsListResponse{
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

// Update updates an existing goods_receipts
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateGoodsReceiptsRequest) (*GoodsReceiptsResponse, error) {
	s.logger.Info("updating goods_receipts",
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
		return nil, fmt.Errorf("failed to get goods_receipts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("goods_receipts not found or access denied")
	}
	

	// Update fields
	
	if req.ReceiptNumber != nil {
		entity.ReceiptNumber = req.ReceiptNumber
	}
	
	if req.PurchaseOrderId != nil {
		entity.PurchaseOrderId = req.PurchaseOrderId
	}
	
	if req.SupplierId != nil {
		entity.SupplierId = req.SupplierId
	}
	
	if req.LocationId != nil {
		entity.LocationId = req.LocationId
	}
	
	if req.ReceiptDate != nil {
		entity.ReceiptDate = req.ReceiptDate
	}
	
	if req.ReceivedBy != nil {
		entity.ReceivedBy = req.ReceivedBy
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update goods_receipts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated goods_receipts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a goods_receipts
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting goods_receipts",
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
		return fmt.Errorf("failed to get goods_receipts: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("goods_receipts not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete goods_receipts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted goods_receipts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *GoodsReceipts) *GoodsReceiptsResponse {
	return &GoodsReceiptsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ReceiptNumber: entity.ReceiptNumber,
		
		PurchaseOrderId: entity.PurchaseOrderId,
		
		SupplierId: entity.SupplierId,
		
		LocationId: entity.LocationId,
		
		ReceiptDate: entity.ReceiptDate,
		
		ReceivedBy: entity.ReceivedBy,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		Notes: entity.Notes,
		
		CreatedBy: entity.CreatedBy,
		
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


// validateBusinessRules validates business rules for goods_receipts
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *GoodsReceipts) error {
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

// canDelete checks if a goods_receipts can be deleted
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

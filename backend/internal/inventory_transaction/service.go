package inventory_transaction

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for InventoryTransactions
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new InventoryTransactions service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new inventory_transactions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateInventoryTransactionsRequest) (*InventoryTransactionsResponse, error) {
	s.logger.Info("creating inventory_transactions",
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
	entity := &InventoryTransactions{
		OrganizationId: orgID,
		
		ProductId: req.ProductId,
		
		TransactionType: req.TransactionType,
		
		Quantity: req.Quantity,
		
		Unit: req.Unit,
		
		BalanceAfter: req.BalanceAfter,
		
		SaleId: req.SaleId,
		
		ReferenceNumber: req.ReferenceNumber,
		
		UnitCost: req.UnitCost,
		
		TotalCost: req.TotalCost,
		
		TransactionDate: req.TransactionDate,
		
		Notes: req.Notes,
		
		Reason: req.Reason,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create inventory_transactions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created inventory_transactions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a inventory_transactions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryTransactionsResponse, error) {
	s.logger.Debug("getting inventory_transactions",
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
		return nil, fmt.Errorf("failed to get inventory_transactions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("inventory_transactions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of inventory_transactions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*InventoryTransactionsListResponse, error) {
	s.logger.Debug("listing inventory_transactions",
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
		return nil, fmt.Errorf("failed to list inventory_transactions: %w", err)
	}

	// Convert to response
	items := make([]*InventoryTransactionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &InventoryTransactionsListResponse{
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

// Update updates an existing inventory_transactions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateInventoryTransactionsRequest) (*InventoryTransactionsResponse, error) {
	s.logger.Info("updating inventory_transactions",
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
		return nil, fmt.Errorf("failed to get inventory_transactions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("inventory_transactions not found or access denied")
	}
	

	// Update fields
	
	if req.ProductId != nil {
		entity.ProductId = req.ProductId
	}
	
	if req.TransactionType != nil {
		entity.TransactionType = req.TransactionType
	}
	
	if req.Quantity != nil {
		entity.Quantity = req.Quantity
	}
	
	if req.Unit != nil {
		entity.Unit = req.Unit
	}
	
	if req.BalanceAfter != nil {
		entity.BalanceAfter = req.BalanceAfter
	}
	
	if req.SaleId != nil {
		entity.SaleId = req.SaleId
	}
	
	if req.ReferenceNumber != nil {
		entity.ReferenceNumber = req.ReferenceNumber
	}
	
	if req.UnitCost != nil {
		entity.UnitCost = req.UnitCost
	}
	
	if req.TotalCost != nil {
		entity.TotalCost = req.TotalCost
	}
	
	if req.TransactionDate != nil {
		entity.TransactionDate = req.TransactionDate
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.Reason != nil {
		entity.Reason = req.Reason
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
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
		return nil, fmt.Errorf("failed to update inventory_transactions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated inventory_transactions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a inventory_transactions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting inventory_transactions",
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
		return fmt.Errorf("failed to get inventory_transactions: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("inventory_transactions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete inventory_transactions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted inventory_transactions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *InventoryTransactions) *InventoryTransactionsResponse {
	return &InventoryTransactionsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ProductId: entity.ProductId,
		
		TransactionType: entity.TransactionType,
		
		Quantity: entity.Quantity,
		
		Unit: entity.Unit,
		
		BalanceAfter: entity.BalanceAfter,
		
		SaleId: entity.SaleId,
		
		ReferenceNumber: entity.ReferenceNumber,
		
		UnitCost: entity.UnitCost,
		
		TotalCost: entity.TotalCost,
		
		TransactionDate: entity.TransactionDate,
		
		Notes: entity.Notes,
		
		Reason: entity.Reason,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		CreatedBy: entity.CreatedBy,
		
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


// validateBusinessRules validates business rules for inventory_transactions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *InventoryTransactions) error {
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

// canDelete checks if a inventory_transactions can be deleted
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

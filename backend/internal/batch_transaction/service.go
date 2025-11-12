package batch_transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/batch_transaction"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/batch_transaction"
	"go.uber.org/zap"
)

// Service handles business logic for BatchTransactions
type Service struct {
	repo   *batch_transaction.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new BatchTransactions service
func NewService(repo *batch_transaction.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new batch_transactions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateBatchTransactionsRequest) (*dto.BatchTransactionsResponse, error) {
	s.logger.Info("creating batch_transactions",
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
	entity := &batch_transaction.BatchTransactions{
		OrganizationID: orgID,
		
		BatchId: req.BatchId,
		
		TransactionType: req.TransactionType,
		
		Quantity: req.Quantity,
		
		BalanceAfter: req.BalanceAfter,
		
		SaleId: req.SaleId,
		
		InventoryTransferId: req.InventoryTransferId,
		
		Reason: req.Reason,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		TransactionDate: req.TransactionDate,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create batch_transactions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created batch_transactions",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a batch_transactions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.BatchTransactionsResponse, error) {
	s.logger.Debug("getting batch_transactions",
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
		return nil, fmt.Errorf("failed to get batch_transactions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("batch_transactions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of batch_transactions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.BatchTransactionsListResponse, error) {
	s.logger.Debug("listing batch_transactions",
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
		return nil, fmt.Errorf("failed to list batch_transactions: %w", err)
	}

	// Convert to response
	items := make([]*dto.BatchTransactionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.BatchTransactionsListResponse{
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

// Update updates an existing batch_transactions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateBatchTransactionsRequest) (*dto.BatchTransactionsResponse, error) {
	s.logger.Info("updating batch_transactions",
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
		return nil, fmt.Errorf("failed to get batch_transactions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("batch_transactions not found or access denied")
	}
	

	// Update fields
	
	if req.BatchId != nil {
		entity.BatchId = *req.BatchId
	}
	
	if req.TransactionType != nil {
		entity.TransactionType = *req.TransactionType
	}
	
	if req.Quantity != nil {
		entity.Quantity = *req.Quantity
	}
	
	if req.BalanceAfter != nil {
		entity.BalanceAfter = *req.BalanceAfter
	}
	
	if req.SaleId != nil {
		entity.SaleId = *req.SaleId
	}
	
	if req.InventoryTransferId != nil {
		entity.InventoryTransferId = *req.InventoryTransferId
	}
	
	if req.Reason != nil {
		entity.Reason = *req.Reason
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.TransactionDate != nil {
		entity.TransactionDate = *req.TransactionDate
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update batch_transactions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated batch_transactions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a batch_transactions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting batch_transactions",
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
		return fmt.Errorf("failed to get batch_transactions: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("batch_transactions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete batch_transactions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted batch_transactions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *batch_transaction.BatchTransactions) *dto.BatchTransactionsResponse {
	return &dto.BatchTransactionsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		BatchId: entity.BatchId,
		
		TransactionType: entity.TransactionType,
		
		Quantity: entity.Quantity,
		
		BalanceAfter: entity.BalanceAfter,
		
		SaleId: entity.SaleId,
		
		InventoryTransferId: entity.InventoryTransferId,
		
		Reason: entity.Reason,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		TransactionDate: entity.TransactionDate,
		
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


// validateBusinessRules validates business rules for batch_transactions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *batch_transaction.BatchTransactions) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a batch_transactions can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

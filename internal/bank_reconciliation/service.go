package bank_reconciliation

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/bank_reconciliation"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/bank_reconciliation"
	"go.uber.org/zap"
)

// Service handles business logic for BankReconciliations
type Service struct {
	repo   *bank_reconciliation.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new BankReconciliations service
func NewService(repo *bank_reconciliation.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new bank_reconciliations
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateBankReconciliationsRequest) (*dto.BankReconciliationsResponse, error) {
	s.logger.Info("creating bank_reconciliations",
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
	entity := &bank_reconciliation.BankReconciliations{
		OrganizationID: orgID,
		
		BankAccountId: req.BankAccountId,
		
		StatementDate: req.StatementDate,
		
		StatementBalance: req.StatementBalance,
		
		ReconciliationDate: req.ReconciliationDate,
		
		BookBalance: req.BookBalance,
		
		ClearedBalance: req.ClearedBalance,
		
		Difference: req.Difference,
		
		Status: req.Status,
		
		IsReconciled: req.IsReconciled,
		
		AccountingPeriodId: req.AccountingPeriodId,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		ReconciledBy: req.ReconciledBy,
		
		ReconciledAt: req.ReconciledAt,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create bank_reconciliations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created bank_reconciliations",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a bank_reconciliations by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.BankReconciliationsResponse, error) {
	s.logger.Debug("getting bank_reconciliations",
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
		return nil, fmt.Errorf("failed to get bank_reconciliations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("bank_reconciliations not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of bank_reconciliations records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.BankReconciliationsListResponse, error) {
	s.logger.Debug("listing bank_reconciliations",
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
		return nil, fmt.Errorf("failed to list bank_reconciliations: %w", err)
	}

	// Convert to response
	items := make([]*dto.BankReconciliationsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.BankReconciliationsListResponse{
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

// Update updates an existing bank_reconciliations
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateBankReconciliationsRequest) (*dto.BankReconciliationsResponse, error) {
	s.logger.Info("updating bank_reconciliations",
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
		return nil, fmt.Errorf("failed to get bank_reconciliations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("bank_reconciliations not found or access denied")
	}
	

	// Update fields
	
	if req.BankAccountId != nil {
		entity.BankAccountId = *req.BankAccountId
	}
	
	if req.StatementDate != nil {
		entity.StatementDate = *req.StatementDate
	}
	
	if req.StatementBalance != nil {
		entity.StatementBalance = *req.StatementBalance
	}
	
	if req.ReconciliationDate != nil {
		entity.ReconciliationDate = *req.ReconciliationDate
	}
	
	if req.BookBalance != nil {
		entity.BookBalance = *req.BookBalance
	}
	
	if req.ClearedBalance != nil {
		entity.ClearedBalance = *req.ClearedBalance
	}
	
	if req.Difference != nil {
		entity.Difference = *req.Difference
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.IsReconciled != nil {
		entity.IsReconciled = *req.IsReconciled
	}
	
	if req.AccountingPeriodId != nil {
		entity.AccountingPeriodId = *req.AccountingPeriodId
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.ReconciledBy != nil {
		entity.ReconciledBy = *req.ReconciledBy
	}
	
	if req.ReconciledAt != nil {
		entity.ReconciledAt = *req.ReconciledAt
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
		return nil, fmt.Errorf("failed to update bank_reconciliations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated bank_reconciliations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a bank_reconciliations
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting bank_reconciliations",
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
		return fmt.Errorf("failed to get bank_reconciliations: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("bank_reconciliations not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete bank_reconciliations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted bank_reconciliations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *bank_reconciliation.BankReconciliations) *dto.BankReconciliationsResponse {
	return &dto.BankReconciliationsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		BankAccountId: entity.BankAccountId,
		
		StatementDate: entity.StatementDate,
		
		StatementBalance: entity.StatementBalance,
		
		ReconciliationDate: entity.ReconciliationDate,
		
		BookBalance: entity.BookBalance,
		
		ClearedBalance: entity.ClearedBalance,
		
		Difference: entity.Difference,
		
		Status: entity.Status,
		
		IsReconciled: entity.IsReconciled,
		
		AccountingPeriodId: entity.AccountingPeriodId,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		ReconciledBy: entity.ReconciledBy,
		
		ReconciledAt: entity.ReconciledAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for bank_reconciliations
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *bank_reconciliation.BankReconciliations) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a bank_reconciliations can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

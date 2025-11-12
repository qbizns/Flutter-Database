package bank_statement

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for BankStatements
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new BankStatements service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new bank_statements
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateBankStatementsRequest) (*BankStatementsResponse, error) {
	s.logger.Info("creating bank_statements",
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
	entity := &BankStatements{
		OrganizationID: orgID,
		
		BankAccountId: req.BankAccountId,
		
		StatementNumber: req.StatementNumber,
		
		StatementDate: req.StatementDate,
		
		PeriodStartDate: req.PeriodStartDate,
		
		PeriodEndDate: req.PeriodEndDate,
		
		OpeningBalance: req.OpeningBalance,
		
		ClosingBalance: req.ClosingBalance,
		
		ImportSource: req.ImportSource,
		
		ImportSource: req.ImportSource,
		
		ImportFileName: req.ImportFileName,
		
		Status: req.Status,
		
		Notes: req.Notes,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create bank_statements: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created bank_statements",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a bank_statements by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*BankStatementsResponse, error) {
	s.logger.Debug("getting bank_statements",
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
		return nil, fmt.Errorf("failed to get bank_statements: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("bank_statements not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of bank_statements records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*BankStatementsListResponse, error) {
	s.logger.Debug("listing bank_statements",
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
		return nil, fmt.Errorf("failed to list bank_statements: %w", err)
	}

	// Convert to response
	items := make([]*BankStatementsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &BankStatementsListResponse{
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

// Update updates an existing bank_statements
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateBankStatementsRequest) (*BankStatementsResponse, error) {
	s.logger.Info("updating bank_statements",
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
		return nil, fmt.Errorf("failed to get bank_statements: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("bank_statements not found or access denied")
	}
	

	// Update fields
	
	if req.BankAccountId != nil {
		entity.BankAccountId = *req.BankAccountId
	}
	
	if req.StatementNumber != nil {
		entity.StatementNumber = *req.StatementNumber
	}
	
	if req.StatementDate != nil {
		entity.StatementDate = *req.StatementDate
	}
	
	if req.PeriodStartDate != nil {
		entity.PeriodStartDate = *req.PeriodStartDate
	}
	
	if req.PeriodEndDate != nil {
		entity.PeriodEndDate = *req.PeriodEndDate
	}
	
	if req.OpeningBalance != nil {
		entity.OpeningBalance = *req.OpeningBalance
	}
	
	if req.ClosingBalance != nil {
		entity.ClosingBalance = *req.ClosingBalance
	}
	
	if req.ImportSource != nil {
		entity.ImportSource = *req.ImportSource
	}
	
	if req.ImportSource != nil {
		entity.ImportSource = *req.ImportSource
	}
	
	if req.ImportFileName != nil {
		entity.ImportFileName = *req.ImportFileName
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
		return nil, fmt.Errorf("failed to update bank_statements: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated bank_statements",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a bank_statements
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting bank_statements",
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
		return fmt.Errorf("failed to get bank_statements: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("bank_statements not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete bank_statements: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted bank_statements",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *BankStatements) *BankStatementsResponse {
	return &BankStatementsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		BankAccountId: entity.BankAccountId,
		
		StatementNumber: entity.StatementNumber,
		
		StatementDate: entity.StatementDate,
		
		PeriodStartDate: entity.PeriodStartDate,
		
		PeriodEndDate: entity.PeriodEndDate,
		
		OpeningBalance: entity.OpeningBalance,
		
		ClosingBalance: entity.ClosingBalance,
		
		ImportSource: entity.ImportSource,
		
		ImportSource: entity.ImportSource,
		
		ImportFileName: entity.ImportFileName,
		
		Status: entity.Status,
		
		Notes: entity.Notes,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for bank_statements
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *BankStatements) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a bank_statements can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

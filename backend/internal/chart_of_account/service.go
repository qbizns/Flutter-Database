package chart_of_account

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

// Service handles business logic for ChartOfAccounts
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new ChartOfAccounts service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new chart_of_accounts
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateChartOfAccountsRequest) (*ChartOfAccountsResponse, error) {
	s.logger.Info("creating chart_of_accounts",
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
	entity := &ChartOfAccounts{
		OrganizationID: orgID,
		
		AccountCode: req.AccountCode,
		
		AccountNumber: req.AccountNumber,
		
		AccountName: req.AccountName,
		
		AccountTypeId: req.AccountTypeId,
		
		AccountSubtypeId: req.AccountSubtypeId,
		
		ParentAccountId: req.ParentAccountId,
		
		AccountLevel: req.AccountLevel,
		
		AccountPath: req.AccountPath,
		
		IsActive: req.IsActive,
		
		IsSystemAccount: req.IsSystemAccount,
		
		IsHeaderAccount: req.IsHeaderAccount,
		
		IsBankAccount: req.IsBankAccount,
		
		IsReconcilable: req.IsReconcilable,
		
		DefaultTaxCode: req.DefaultTaxCode,
		
		CurrencyCode: req.CurrencyCode,
		
		OpeningBalance: req.OpeningBalance,
		
		OpeningBalanceDate: req.OpeningBalanceDate,
		
		CurrentDebitBalance: req.CurrentDebitBalance,
		
		CurrentCreditBalance: req.CurrentCreditBalance,
		
		CurrentBalance: req.CurrentBalance,
		
		LastBalanceUpdate: req.LastBalanceUpdate,
		
		Description: req.Description,
		
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
		return nil, fmt.Errorf("failed to create chart_of_accounts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created chart_of_accounts",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a chart_of_accounts by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ChartOfAccountsResponse, error) {
	s.logger.Debug("getting chart_of_accounts",
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
		return nil, fmt.Errorf("failed to get chart_of_accounts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("chart_of_accounts not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of chart_of_accounts records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ChartOfAccountsListResponse, error) {
	s.logger.Debug("listing chart_of_accounts",
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
		return nil, fmt.Errorf("failed to list chart_of_accounts: %w", err)
	}

	// Convert to response
	items := make([]*ChartOfAccountsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ChartOfAccountsListResponse{
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

// Update updates an existing chart_of_accounts
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateChartOfAccountsRequest) (*ChartOfAccountsResponse, error) {
	s.logger.Info("updating chart_of_accounts",
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
		return nil, fmt.Errorf("failed to get chart_of_accounts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("chart_of_accounts not found or access denied")
	}
	

	// Update fields
	
	if req.AccountCode != nil {
		entity.AccountCode = *req.AccountCode
	}
	
	if req.AccountNumber != nil {
		entity.AccountNumber = *req.AccountNumber
	}
	
	if req.AccountName != nil {
		entity.AccountName = *req.AccountName
	}
	
	if req.AccountTypeId != nil {
		entity.AccountTypeId = *req.AccountTypeId
	}
	
	if req.AccountSubtypeId != nil {
		entity.AccountSubtypeId = *req.AccountSubtypeId
	}
	
	if req.ParentAccountId != nil {
		entity.ParentAccountId = *req.ParentAccountId
	}
	
	if req.AccountLevel != nil {
		entity.AccountLevel = *req.AccountLevel
	}
	
	if req.AccountPath != nil {
		entity.AccountPath = *req.AccountPath
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.IsSystemAccount != nil {
		entity.IsSystemAccount = *req.IsSystemAccount
	}
	
	if req.IsHeaderAccount != nil {
		entity.IsHeaderAccount = *req.IsHeaderAccount
	}
	
	if req.IsBankAccount != nil {
		entity.IsBankAccount = *req.IsBankAccount
	}
	
	if req.IsReconcilable != nil {
		entity.IsReconcilable = *req.IsReconcilable
	}
	
	if req.DefaultTaxCode != nil {
		entity.DefaultTaxCode = *req.DefaultTaxCode
	}
	
	if req.CurrencyCode != nil {
		entity.CurrencyCode = *req.CurrencyCode
	}
	
	if req.OpeningBalance != nil {
		entity.OpeningBalance = *req.OpeningBalance
	}
	
	if req.OpeningBalanceDate != nil {
		entity.OpeningBalanceDate = *req.OpeningBalanceDate
	}
	
	if req.CurrentDebitBalance != nil {
		entity.CurrentDebitBalance = *req.CurrentDebitBalance
	}
	
	if req.CurrentCreditBalance != nil {
		entity.CurrentCreditBalance = *req.CurrentCreditBalance
	}
	
	if req.CurrentBalance != nil {
		entity.CurrentBalance = *req.CurrentBalance
	}
	
	if req.LastBalanceUpdate != nil {
		entity.LastBalanceUpdate = *req.LastBalanceUpdate
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
		return nil, fmt.Errorf("failed to update chart_of_accounts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated chart_of_accounts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a chart_of_accounts
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting chart_of_accounts",
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
		return fmt.Errorf("failed to get chart_of_accounts: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("chart_of_accounts not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete chart_of_accounts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted chart_of_accounts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *ChartOfAccounts) *ChartOfAccountsResponse {
	return &ChartOfAccountsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		AccountCode: entity.AccountCode,
		
		AccountNumber: entity.AccountNumber,
		
		AccountName: entity.AccountName,
		
		AccountTypeId: entity.AccountTypeId,
		
		AccountSubtypeId: entity.AccountSubtypeId,
		
		ParentAccountId: entity.ParentAccountId,
		
		AccountLevel: entity.AccountLevel,
		
		AccountPath: entity.AccountPath,
		
		IsActive: entity.IsActive,
		
		IsSystemAccount: entity.IsSystemAccount,
		
		IsHeaderAccount: entity.IsHeaderAccount,
		
		IsBankAccount: entity.IsBankAccount,
		
		IsReconcilable: entity.IsReconcilable,
		
		DefaultTaxCode: entity.DefaultTaxCode,
		
		CurrencyCode: entity.CurrencyCode,
		
		OpeningBalance: entity.OpeningBalance,
		
		OpeningBalanceDate: entity.OpeningBalanceDate,
		
		CurrentDebitBalance: entity.CurrentDebitBalance,
		
		CurrentCreditBalance: entity.CurrentCreditBalance,
		
		CurrentBalance: entity.CurrentBalance,
		
		LastBalanceUpdate: entity.LastBalanceUpdate,
		
		Description: entity.Description,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
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


// validateBusinessRules validates business rules for chart_of_accounts
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *ChartOfAccounts) error {
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

// canDelete checks if a chart_of_accounts can be deleted
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

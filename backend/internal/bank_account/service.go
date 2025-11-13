package bank_account

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

// Service handles business logic for BankAccounts
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new BankAccounts service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new bank_accounts
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateBankAccountsRequest) (*BankAccountsResponse, error) {
	s.logger.Info("creating bank_accounts",
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
	entity := &BankAccounts{
		OrganizationID: orgID,
		
		ChartAccountId: req.ChartAccountId,
		
		BankName: req.BankName,
		
		AccountNumber: req.AccountNumber,
		
		AccountType: req.AccountType,
		
		RoutingNumber: req.RoutingNumber,
		
		SwiftCode: req.SwiftCode,
		
		CurrencyCode: req.CurrencyCode,
		
		CurrentBalance: req.CurrentBalance,
		
		StatementBalance: req.StatementBalance,
		
		LastStatementDate: req.LastStatementDate,
		
		IsActive: req.IsActive,
		
		OnlineBankingEnabled: req.OnlineBankingEnabled,
		
		LastSyncDate: req.LastSyncDate,
		
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
		return nil, fmt.Errorf("failed to create bank_accounts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created bank_accounts",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a bank_accounts by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*BankAccountsResponse, error) {
	s.logger.Debug("getting bank_accounts",
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
		return nil, fmt.Errorf("failed to get bank_accounts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("bank_accounts not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of bank_accounts records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*BankAccountsListResponse, error) {
	s.logger.Debug("listing bank_accounts",
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
		return nil, fmt.Errorf("failed to list bank_accounts: %w", err)
	}

	// Convert to response
	items := make([]*BankAccountsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &BankAccountsListResponse{
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

// Update updates an existing bank_accounts
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateBankAccountsRequest) (*BankAccountsResponse, error) {
	s.logger.Info("updating bank_accounts",
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
		return nil, fmt.Errorf("failed to get bank_accounts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("bank_accounts not found or access denied")
	}
	

	// Update fields
	
	if req.ChartAccountId != nil {
		entity.ChartAccountId = *req.ChartAccountId
	}
	
	if req.BankName != nil {
		entity.BankName = *req.BankName
	}
	
	if req.AccountNumber != nil {
		entity.AccountNumber = *req.AccountNumber
	}
	
	if req.AccountType != nil {
		entity.AccountType = *req.AccountType
	}
	
	if req.RoutingNumber != nil {
		entity.RoutingNumber = *req.RoutingNumber
	}
	
	if req.SwiftCode != nil {
		entity.SwiftCode = *req.SwiftCode
	}
	
	if req.CurrencyCode != nil {
		entity.CurrencyCode = *req.CurrencyCode
	}
	
	if req.CurrentBalance != nil {
		entity.CurrentBalance = *req.CurrentBalance
	}
	
	if req.StatementBalance != nil {
		entity.StatementBalance = *req.StatementBalance
	}
	
	if req.LastStatementDate != nil {
		entity.LastStatementDate = *req.LastStatementDate
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.OnlineBankingEnabled != nil {
		entity.OnlineBankingEnabled = *req.OnlineBankingEnabled
	}
	
	if req.LastSyncDate != nil {
		entity.LastSyncDate = *req.LastSyncDate
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
		return nil, fmt.Errorf("failed to update bank_accounts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated bank_accounts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a bank_accounts
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting bank_accounts",
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
		return fmt.Errorf("failed to get bank_accounts: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("bank_accounts not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete bank_accounts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted bank_accounts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *BankAccounts) *BankAccountsResponse {
	return &BankAccountsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ChartAccountId: entity.ChartAccountId,
		
		BankName: entity.BankName,
		
		AccountNumber: entity.AccountNumber,
		
		AccountType: entity.AccountType,
		
		RoutingNumber: entity.RoutingNumber,
		
		SwiftCode: entity.SwiftCode,
		
		CurrencyCode: entity.CurrencyCode,
		
		CurrentBalance: entity.CurrentBalance,
		
		StatementBalance: entity.StatementBalance,
		
		LastStatementDate: entity.LastStatementDate,
		
		IsActive: entity.IsActive,
		
		OnlineBankingEnabled: entity.OnlineBankingEnabled,
		
		LastSyncDate: entity.LastSyncDate,
		
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


// validateBusinessRules validates business rules for bank_accounts
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *BankAccounts) error {
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

// canDelete checks if a bank_accounts can be deleted
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

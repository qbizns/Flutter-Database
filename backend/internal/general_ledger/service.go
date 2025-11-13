package general_ledger

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

// Service handles business logic for GeneralLedger
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new GeneralLedger service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new general_ledger
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateGeneralLedgerRequest) (*GeneralLedgerResponse, error) {
	s.logger.Info("creating general_ledger",
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
	entity := &GeneralLedger{
		OrganizationID: orgID,
		
		JournalEntryId: req.JournalEntryId,
		
		JournalEntryLineId: req.JournalEntryLineId,
		
		AccountId: req.AccountId,
		
		TransactionDate: req.TransactionDate,
		
		PostingDate: req.PostingDate,
		
		AccountingPeriodId: req.AccountingPeriodId,
		
		FiscalYearId: req.FiscalYearId,
		
		DebitAmount: req.DebitAmount,
		
		CreditAmount: req.CreditAmount,
		
		RunningDebitBalance: req.RunningDebitBalance,
		
		RunningCreditBalance: req.RunningCreditBalance,
		
		RunningBalance: req.RunningBalance,
		
		SourceModule: req.SourceModule,
		
		SourceDocumentType: req.SourceDocumentType,
		
		SourceDocumentId: req.SourceDocumentId,
		
		ReferenceNumber: req.ReferenceNumber,
		
		LocationId: req.LocationId,
		
		Department: req.Department,
		
		ProjectCode: req.ProjectCode,
		
		CostCenter: req.CostCenter,
		
		Description: req.Description,
		
		IsReversed: req.IsReversed,
		
		ReversalGlId: req.ReversalGlId,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		(debitAmount: req.(debitAmount,
		
		(creditAmount: req.(creditAmount,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create general_ledger: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created general_ledger",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a general_ledger by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*GeneralLedgerResponse, error) {
	s.logger.Debug("getting general_ledger",
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
		return nil, fmt.Errorf("failed to get general_ledger: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("general_ledger not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of general_ledger records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*GeneralLedgerListResponse, error) {
	s.logger.Debug("listing general_ledger",
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
		return nil, fmt.Errorf("failed to list general_ledger: %w", err)
	}

	// Convert to response
	items := make([]*GeneralLedgerResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &GeneralLedgerListResponse{
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

// Update updates an existing general_ledger
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateGeneralLedgerRequest) (*GeneralLedgerResponse, error) {
	s.logger.Info("updating general_ledger",
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
		return nil, fmt.Errorf("failed to get general_ledger: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("general_ledger not found or access denied")
	}
	

	// Update fields
	
	if req.JournalEntryId != nil {
		entity.JournalEntryId = *req.JournalEntryId
	}
	
	if req.JournalEntryLineId != nil {
		entity.JournalEntryLineId = *req.JournalEntryLineId
	}
	
	if req.AccountId != nil {
		entity.AccountId = *req.AccountId
	}
	
	if req.TransactionDate != nil {
		entity.TransactionDate = *req.TransactionDate
	}
	
	if req.PostingDate != nil {
		entity.PostingDate = *req.PostingDate
	}
	
	if req.AccountingPeriodId != nil {
		entity.AccountingPeriodId = *req.AccountingPeriodId
	}
	
	if req.FiscalYearId != nil {
		entity.FiscalYearId = *req.FiscalYearId
	}
	
	if req.DebitAmount != nil {
		entity.DebitAmount = *req.DebitAmount
	}
	
	if req.CreditAmount != nil {
		entity.CreditAmount = *req.CreditAmount
	}
	
	if req.RunningDebitBalance != nil {
		entity.RunningDebitBalance = *req.RunningDebitBalance
	}
	
	if req.RunningCreditBalance != nil {
		entity.RunningCreditBalance = *req.RunningCreditBalance
	}
	
	if req.RunningBalance != nil {
		entity.RunningBalance = *req.RunningBalance
	}
	
	if req.SourceModule != nil {
		entity.SourceModule = *req.SourceModule
	}
	
	if req.SourceDocumentType != nil {
		entity.SourceDocumentType = *req.SourceDocumentType
	}
	
	if req.SourceDocumentId != nil {
		entity.SourceDocumentId = *req.SourceDocumentId
	}
	
	if req.ReferenceNumber != nil {
		entity.ReferenceNumber = *req.ReferenceNumber
	}
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.Department != nil {
		entity.Department = *req.Department
	}
	
	if req.ProjectCode != nil {
		entity.ProjectCode = *req.ProjectCode
	}
	
	if req.CostCenter != nil {
		entity.CostCenter = *req.CostCenter
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.IsReversed != nil {
		entity.IsReversed = *req.IsReversed
	}
	
	if req.ReversalGlId != nil {
		entity.ReversalGlId = *req.ReversalGlId
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.(debitAmount != nil {
		entity.(debitAmount = *req.(debitAmount
	}
	
	if req.(creditAmount != nil {
		entity.(creditAmount = *req.(creditAmount
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update general_ledger: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated general_ledger",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a general_ledger
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting general_ledger",
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
		return fmt.Errorf("failed to get general_ledger: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("general_ledger not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete general_ledger: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted general_ledger",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *GeneralLedger) *GeneralLedgerResponse {
	return &GeneralLedgerResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		JournalEntryId: entity.JournalEntryId,
		
		JournalEntryLineId: entity.JournalEntryLineId,
		
		AccountId: entity.AccountId,
		
		TransactionDate: entity.TransactionDate,
		
		PostingDate: entity.PostingDate,
		
		AccountingPeriodId: entity.AccountingPeriodId,
		
		FiscalYearId: entity.FiscalYearId,
		
		DebitAmount: entity.DebitAmount,
		
		CreditAmount: entity.CreditAmount,
		
		RunningDebitBalance: entity.RunningDebitBalance,
		
		RunningCreditBalance: entity.RunningCreditBalance,
		
		RunningBalance: entity.RunningBalance,
		
		SourceModule: entity.SourceModule,
		
		SourceDocumentType: entity.SourceDocumentType,
		
		SourceDocumentId: entity.SourceDocumentId,
		
		ReferenceNumber: entity.ReferenceNumber,
		
		LocationId: entity.LocationId,
		
		Department: entity.Department,
		
		ProjectCode: entity.ProjectCode,
		
		CostCenter: entity.CostCenter,
		
		Description: entity.Description,
		
		IsReversed: entity.IsReversed,
		
		ReversalGlId: entity.ReversalGlId,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		(debitAmount: entity.(debitAmount,
		
		(creditAmount: entity.(creditAmount,
		
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


// validateBusinessRules validates business rules for general_ledger
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *GeneralLedger) error {
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

// canDelete checks if a general_ledger can be deleted
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

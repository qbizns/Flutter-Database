package journal_entry_line

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

// Service handles business logic for JournalEntryLines
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new JournalEntryLines service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new journal_entry_lines
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateJournalEntryLinesRequest) (*JournalEntryLinesResponse, error) {
	s.logger.Info("creating journal_entry_lines",
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
	entity := &JournalEntryLines{
		OrganizationID: orgID,
		
		JournalEntryId: req.JournalEntryId,
		
		LineNumber: req.LineNumber,
		
		AccountId: req.AccountId,
		
		DebitAmount: req.DebitAmount,
		
		CreditAmount: req.CreditAmount,
		
		LocationId: req.LocationId,
		
		Department: req.Department,
		
		ProjectCode: req.ProjectCode,
		
		CostCenter: req.CostCenter,
		
		TaxCode: req.TaxCode,
		
		TaxAmount: req.TaxAmount,
		
		Description: req.Description,
		
		Memo: req.Memo,
		
		IsReconciled: req.IsReconciled,
		
		ReconciledAt: req.ReconciledAt,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		(debitAmount: req.(debitAmount,
		
		(creditAmount: req.(creditAmount,
		
		(debitAmount: req.(debitAmount,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create journal_entry_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created journal_entry_lines",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a journal_entry_lines by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*JournalEntryLinesResponse, error) {
	s.logger.Debug("getting journal_entry_lines",
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
		return nil, fmt.Errorf("failed to get journal_entry_lines: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("journal_entry_lines not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of journal_entry_lines records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*JournalEntryLinesListResponse, error) {
	s.logger.Debug("listing journal_entry_lines",
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
		return nil, fmt.Errorf("failed to list journal_entry_lines: %w", err)
	}

	// Convert to response
	items := make([]*JournalEntryLinesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &JournalEntryLinesListResponse{
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

// Update updates an existing journal_entry_lines
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateJournalEntryLinesRequest) (*JournalEntryLinesResponse, error) {
	s.logger.Info("updating journal_entry_lines",
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
		return nil, fmt.Errorf("failed to get journal_entry_lines: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("journal_entry_lines not found or access denied")
	}
	

	// Update fields
	
	if req.JournalEntryId != nil {
		entity.JournalEntryId = *req.JournalEntryId
	}
	
	if req.LineNumber != nil {
		entity.LineNumber = *req.LineNumber
	}
	
	if req.AccountId != nil {
		entity.AccountId = *req.AccountId
	}
	
	if req.DebitAmount != nil {
		entity.DebitAmount = *req.DebitAmount
	}
	
	if req.CreditAmount != nil {
		entity.CreditAmount = *req.CreditAmount
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
	
	if req.TaxCode != nil {
		entity.TaxCode = *req.TaxCode
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = *req.TaxAmount
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.Memo != nil {
		entity.Memo = *req.Memo
	}
	
	if req.IsReconciled != nil {
		entity.IsReconciled = *req.IsReconciled
	}
	
	if req.ReconciledAt != nil {
		entity.ReconciledAt = *req.ReconciledAt
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
	
	if req.(debitAmount != nil {
		entity.(debitAmount = *req.(debitAmount
	}
	
	if req.(creditAmount != nil {
		entity.(creditAmount = *req.(creditAmount
	}
	
	if req.(debitAmount != nil {
		entity.(debitAmount = *req.(debitAmount
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update journal_entry_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated journal_entry_lines",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a journal_entry_lines
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting journal_entry_lines",
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
		return fmt.Errorf("failed to get journal_entry_lines: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("journal_entry_lines not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete journal_entry_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted journal_entry_lines",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *JournalEntryLines) *JournalEntryLinesResponse {
	return &JournalEntryLinesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		JournalEntryId: entity.JournalEntryId,
		
		LineNumber: entity.LineNumber,
		
		AccountId: entity.AccountId,
		
		DebitAmount: entity.DebitAmount,
		
		CreditAmount: entity.CreditAmount,
		
		LocationId: entity.LocationId,
		
		Department: entity.Department,
		
		ProjectCode: entity.ProjectCode,
		
		CostCenter: entity.CostCenter,
		
		TaxCode: entity.TaxCode,
		
		TaxAmount: entity.TaxAmount,
		
		Description: entity.Description,
		
		Memo: entity.Memo,
		
		IsReconciled: entity.IsReconciled,
		
		ReconciledAt: entity.ReconciledAt,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		(debitAmount: entity.(debitAmount,
		
		(creditAmount: entity.(creditAmount,
		
		(debitAmount: entity.(debitAmount,
		
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


// validateBusinessRules validates business rules for journal_entry_lines
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *JournalEntryLines) error {
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

// canDelete checks if a journal_entry_lines can be deleted
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

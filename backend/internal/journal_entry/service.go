package journal_entry

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for JournalEntries
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new JournalEntries service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new journal_entries
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateJournalEntriesRequest) (*JournalEntriesResponse, error) {
	s.logger.Info("creating journal_entries",
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
	entity := &JournalEntries{
		OrganizationId: orgID,
		
		EntryNumber: req.EntryNumber,
		
		EntryTypeId: req.EntryTypeId,
		
		EntryDate: req.EntryDate,
		
		PostingDate: req.PostingDate,
		
		AccountingPeriodId: req.AccountingPeriodId,
		
		FiscalYearId: req.FiscalYearId,
		
		Status: req.Status,
		
		IsPosted: req.IsPosted,
		
		IsReversed: req.IsReversed,
		
		ReversalEntryId: req.ReversalEntryId,
		
		SourceModule: req.SourceModule,
		
		SourceDocumentType: req.SourceDocumentType,
		
		SourceDocumentId: req.SourceDocumentId,
		
		ReferenceNumber: req.ReferenceNumber,
		
		TotalDebit: req.TotalDebit,
		
		TotalCredit: req.TotalCredit,
		
		Description: req.Description,
		
		Notes: req.Notes,
		
		RequiresApproval: req.RequiresApproval,
		
		ApprovedBy: req.ApprovedBy,
		
		ApprovedAt: req.ApprovedAt,
		
		PostedBy: req.PostedBy,
		
		PostedAt: req.PostedAt,
		
		Attachments: req.Attachments,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		(isPosted: req.(isPosted,
		
		(ABS(totalDebit: req.(ABS(totalDebit,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create journal_entries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created journal_entries",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a journal_entries by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*JournalEntriesResponse, error) {
	s.logger.Debug("getting journal_entries",
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
		return nil, fmt.Errorf("failed to get journal_entries: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("journal_entries not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of journal_entries records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*JournalEntriesListResponse, error) {
	s.logger.Debug("listing journal_entries",
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
		return nil, fmt.Errorf("failed to list journal_entries: %w", err)
	}

	// Convert to response
	items := make([]*JournalEntriesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &JournalEntriesListResponse{
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

// Update updates an existing journal_entries
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateJournalEntriesRequest) (*JournalEntriesResponse, error) {
	s.logger.Info("updating journal_entries",
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
		return nil, fmt.Errorf("failed to get journal_entries: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("journal_entries not found or access denied")
	}
	

	// Update fields
	
	if req.EntryNumber != nil {
		entity.EntryNumber = req.EntryNumber
	}
	
	if req.EntryTypeId != nil {
		entity.EntryTypeId = req.EntryTypeId
	}
	
	if req.EntryDate != nil {
		entity.EntryDate = req.EntryDate
	}
	
	if req.PostingDate != nil {
		entity.PostingDate = req.PostingDate
	}
	
	if req.AccountingPeriodId != nil {
		entity.AccountingPeriodId = req.AccountingPeriodId
	}
	
	if req.FiscalYearId != nil {
		entity.FiscalYearId = req.FiscalYearId
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.IsPosted != nil {
		entity.IsPosted = req.IsPosted
	}
	
	if req.IsReversed != nil {
		entity.IsReversed = req.IsReversed
	}
	
	if req.ReversalEntryId != nil {
		entity.ReversalEntryId = req.ReversalEntryId
	}
	
	if req.SourceModule != nil {
		entity.SourceModule = req.SourceModule
	}
	
	if req.SourceDocumentType != nil {
		entity.SourceDocumentType = req.SourceDocumentType
	}
	
	if req.SourceDocumentId != nil {
		entity.SourceDocumentId = req.SourceDocumentId
	}
	
	if req.ReferenceNumber != nil {
		entity.ReferenceNumber = req.ReferenceNumber
	}
	
	if req.TotalDebit != nil {
		entity.TotalDebit = req.TotalDebit
	}
	
	if req.TotalCredit != nil {
		entity.TotalCredit = req.TotalCredit
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.RequiresApproval != nil {
		entity.RequiresApproval = req.RequiresApproval
	}
	
	if req.ApprovedBy != nil {
		entity.ApprovedBy = req.ApprovedBy
	}
	
	if req.ApprovedAt != nil {
		entity.ApprovedAt = req.ApprovedAt
	}
	
	if req.PostedBy != nil {
		entity.PostedBy = req.PostedBy
	}
	
	if req.PostedAt != nil {
		entity.PostedAt = req.PostedAt
	}
	
	if req.Attachments != nil {
		entity.Attachments = req.Attachments
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
	
	if req.(isPosted != nil {
		entity.(isPosted = *req.(isPosted
	}
	
	if req.(ABS(totalDebit != nil {
		entity.(ABS(totalDebit = *req.(ABS(totalDebit
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update journal_entries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated journal_entries",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a journal_entries
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting journal_entries",
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
		return fmt.Errorf("failed to get journal_entries: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("journal_entries not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete journal_entries: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted journal_entries",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *JournalEntries) *JournalEntriesResponse {
	return &JournalEntriesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		EntryNumber: entity.EntryNumber,
		
		EntryTypeId: entity.EntryTypeId,
		
		EntryDate: entity.EntryDate,
		
		PostingDate: entity.PostingDate,
		
		AccountingPeriodId: entity.AccountingPeriodId,
		
		FiscalYearId: entity.FiscalYearId,
		
		Status: entity.Status,
		
		IsPosted: entity.IsPosted,
		
		IsReversed: entity.IsReversed,
		
		ReversalEntryId: entity.ReversalEntryId,
		
		SourceModule: entity.SourceModule,
		
		SourceDocumentType: entity.SourceDocumentType,
		
		SourceDocumentId: entity.SourceDocumentId,
		
		ReferenceNumber: entity.ReferenceNumber,
		
		TotalDebit: entity.TotalDebit,
		
		TotalCredit: entity.TotalCredit,
		
		Description: entity.Description,
		
		Notes: entity.Notes,
		
		RequiresApproval: entity.RequiresApproval,
		
		ApprovedBy: entity.ApprovedBy,
		
		ApprovedAt: entity.ApprovedAt,
		
		PostedBy: entity.PostedBy,
		
		PostedAt: entity.PostedAt,
		
		Attachments: entity.Attachments,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		(isPosted: entity.(isPosted,
		
		(ABS(totalDebit: entity.(ABS(totalDebit,
		
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


// validateBusinessRules validates business rules for journal_entries
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *JournalEntries) error {
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

// canDelete checks if a journal_entries can be deleted
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

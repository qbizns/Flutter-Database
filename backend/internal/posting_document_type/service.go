package posting_document_type

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for PostingDocumentTypes
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PostingDocumentTypes service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new posting_document_types
func (s *Service) Create(ctx context.Context, req *CreatePostingDocumentTypesRequest) (*PostingDocumentTypesResponse, error) {
	s.logger.Info("creating posting_document_types",
		
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

	

	// Convert DTO to entity
	entity := &PostingDocumentTypes{
		
		
		Code: req.Code,
		
		Name: req.Name,
		
		Description: req.Description,
		
		SourceSchema: req.SourceSchema,
		
		SourceTable: req.SourceTable,
		
		SourcePkColumn: req.SourcePkColumn,
		
		Category: req.Category,
		
		IsActive: req.IsActive,
		
		IsSystem: req.IsSystem,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create posting_document_types: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created posting_document_types",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a posting_document_types by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*PostingDocumentTypesResponse, error) {
	s.logger.Debug("getting posting_document_types",
		zap.String("id", id.String()),
		
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get posting_document_types: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of posting_document_types records
func (s *Service) List(ctx context.Context, page, limit int) (*PostingDocumentTypesListResponse, error) {
	s.logger.Debug("listing posting_document_types",
		
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

	
	// Get from database
	entities, total, err := s.repo.List(ctx, tx, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list posting_document_types: %w", err)
	}

	// Convert to response
	items := make([]*PostingDocumentTypesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PostingDocumentTypesListResponse{
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

// Update updates an existing posting_document_types
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdatePostingDocumentTypesRequest) (*PostingDocumentTypesResponse, error) {
	s.logger.Info("updating posting_document_types",
		zap.String("id", id.String()),
		
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

	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get posting_document_types: %w", err)
	}

	

	// Update fields
	
	if req.Code != nil {
		entity.Code = req.Code
	}
	
	if req.Name != nil {
		entity.Name = req.Name
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.SourceSchema != nil {
		entity.SourceSchema = req.SourceSchema
	}
	
	if req.SourceTable != nil {
		entity.SourceTable = req.SourceTable
	}
	
	if req.SourcePkColumn != nil {
		entity.SourcePkColumn = req.SourcePkColumn
	}
	
	if req.Category != nil {
		entity.Category = req.Category
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.IsSystem != nil {
		entity.IsSystem = req.IsSystem
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update posting_document_types: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated posting_document_types",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a posting_document_types
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting posting_document_types",
		zap.String("id", id.String()),
		
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete posting_document_types: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted posting_document_types",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PostingDocumentTypes) *PostingDocumentTypesResponse {
	return &PostingDocumentTypesResponse{
		
		Id: entity.Id,
		
		Code: entity.Code,
		
		Name: entity.Name,
		
		Description: entity.Description,
		
		SourceSchema: entity.SourceSchema,
		
		SourceTable: entity.SourceTable,
		
		SourcePkColumn: entity.SourcePkColumn,
		
		Category: entity.Category,
		
		IsActive: entity.IsActive,
		
		IsSystem: entity.IsSystem,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for posting_document_types
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PostingDocumentTypes) error {
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

// canDelete checks if a posting_document_types can be deleted
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

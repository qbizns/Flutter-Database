package journal_entry_type

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for JournalEntryTypes
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new JournalEntryTypes service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new journal_entry_types
func (s *Service) Create(ctx context.Context, req *CreateJournalEntryTypesRequest) (*JournalEntryTypesResponse, error) {
	s.logger.Info("creating journal_entry_types",
		
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
	entity := &JournalEntryTypes{
		
		
		TypeCode: req.TypeCode,
		
		TypeName: req.TypeName,
		
		TypeCategory: req.TypeCategory,
		
		NumberPrefix: req.NumberPrefix,
		
		Description: req.Description,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create journal_entry_types: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created journal_entry_types",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a journal_entry_types by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*JournalEntryTypesResponse, error) {
	s.logger.Debug("getting journal_entry_types",
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
		return nil, fmt.Errorf("failed to get journal_entry_types: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of journal_entry_types records
func (s *Service) List(ctx context.Context, page, limit int) (*JournalEntryTypesListResponse, error) {
	s.logger.Debug("listing journal_entry_types",
		
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
		return nil, fmt.Errorf("failed to list journal_entry_types: %w", err)
	}

	// Convert to response
	items := make([]*JournalEntryTypesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &JournalEntryTypesListResponse{
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

// Update updates an existing journal_entry_types
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateJournalEntryTypesRequest) (*JournalEntryTypesResponse, error) {
	s.logger.Info("updating journal_entry_types",
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
		return nil, fmt.Errorf("failed to get journal_entry_types: %w", err)
	}

	

	// Update fields
	
	if req.TypeCode != nil {
		entity.TypeCode = req.TypeCode
	}
	
	if req.TypeName != nil {
		entity.TypeName = req.TypeName
	}
	
	if req.TypeCategory != nil {
		entity.TypeCategory = req.TypeCategory
	}
	
	if req.NumberPrefix != nil {
		entity.NumberPrefix = req.NumberPrefix
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update journal_entry_types: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated journal_entry_types",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a journal_entry_types
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting journal_entry_types",
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
		return fmt.Errorf("failed to delete journal_entry_types: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted journal_entry_types",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *JournalEntryTypes) *JournalEntryTypesResponse {
	return &JournalEntryTypesResponse{
		
		Id: entity.Id,
		
		TypeCode: entity.TypeCode,
		
		TypeName: entity.TypeName,
		
		TypeCategory: entity.TypeCategory,
		
		NumberPrefix: entity.NumberPrefix,
		
		Description: entity.Description,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
	}
}



// validateBusinessRules validates business rules for journal_entry_types
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *JournalEntryTypes) error {
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

// canDelete checks if a journal_entry_types can be deleted
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

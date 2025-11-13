package posting_concept

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for PostingConcepts
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PostingConcepts service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new posting_concepts
func (s *Service) Create(ctx context.Context, req *CreatePostingConceptsRequest) (*PostingConceptsResponse, error) {
	s.logger.Info("creating posting_concepts",
		
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
	entity := &PostingConcepts{
		
		
		ConceptKey: req.ConceptKey,
		
		DefaultLabel: req.DefaultLabel,
		
		DefaultDescription: req.DefaultDescription,
		
		ExpectedAccountTypeId: req.ExpectedAccountTypeId,
		
		NormalSide: req.NormalSide,
		
		ExampleCode: req.ExampleCode,
		
		ExampleAccountName: req.ExampleAccountName,
		
		IsSystem: req.IsSystem,
		
		ConceptCategory: req.ConceptCategory,
		
		SortOrder: req.SortOrder,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create posting_concepts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created posting_concepts",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a posting_concepts by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*PostingConceptsResponse, error) {
	s.logger.Debug("getting posting_concepts",
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
		return nil, fmt.Errorf("failed to get posting_concepts: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of posting_concepts records
func (s *Service) List(ctx context.Context, page, limit int) (*PostingConceptsListResponse, error) {
	s.logger.Debug("listing posting_concepts",
		
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
		return nil, fmt.Errorf("failed to list posting_concepts: %w", err)
	}

	// Convert to response
	items := make([]*PostingConceptsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PostingConceptsListResponse{
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

// Update updates an existing posting_concepts
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdatePostingConceptsRequest) (*PostingConceptsResponse, error) {
	s.logger.Info("updating posting_concepts",
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
		return nil, fmt.Errorf("failed to get posting_concepts: %w", err)
	}

	

	// Update fields
	
	if req.ConceptKey != nil {
		entity.ConceptKey = req.ConceptKey
	}
	
	if req.DefaultLabel != nil {
		entity.DefaultLabel = req.DefaultLabel
	}
	
	if req.DefaultDescription != nil {
		entity.DefaultDescription = req.DefaultDescription
	}
	
	if req.ExpectedAccountTypeId != nil {
		entity.ExpectedAccountTypeId = req.ExpectedAccountTypeId
	}
	
	if req.NormalSide != nil {
		entity.NormalSide = req.NormalSide
	}
	
	if req.ExampleCode != nil {
		entity.ExampleCode = req.ExampleCode
	}
	
	if req.ExampleAccountName != nil {
		entity.ExampleAccountName = req.ExampleAccountName
	}
	
	if req.IsSystem != nil {
		entity.IsSystem = req.IsSystem
	}
	
	if req.ConceptCategory != nil {
		entity.ConceptCategory = req.ConceptCategory
	}
	
	if req.SortOrder != nil {
		entity.SortOrder = req.SortOrder
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
		return nil, fmt.Errorf("failed to update posting_concepts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated posting_concepts",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a posting_concepts
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting posting_concepts",
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
		return fmt.Errorf("failed to delete posting_concepts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted posting_concepts",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PostingConcepts) *PostingConceptsResponse {
	return &PostingConceptsResponse{
		
		ConceptKey: entity.ConceptKey,
		
		DefaultLabel: entity.DefaultLabel,
		
		DefaultDescription: entity.DefaultDescription,
		
		ExpectedAccountTypeId: entity.ExpectedAccountTypeId,
		
		NormalSide: entity.NormalSide,
		
		ExampleCode: entity.ExampleCode,
		
		ExampleAccountName: entity.ExampleAccountName,
		
		IsSystem: entity.IsSystem,
		
		ConceptCategory: entity.ConceptCategory,
		
		SortOrder: entity.SortOrder,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for posting_concepts
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PostingConcepts) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a posting_concepts can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

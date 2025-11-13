package posting_rule

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for PostingRules
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PostingRules service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new posting_rules
func (s *Service) Create(ctx context.Context, req *CreatePostingRulesRequest) (*PostingRulesResponse, error) {
	s.logger.Info("creating posting_rules",
		
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
	entity := &PostingRules{
		
		
		PostingProfileDocumentId: req.PostingProfileDocumentId,
		
		RuleCode: req.RuleCode,
		
		RuleName: req.RuleName,
		
		Description: req.Description,
		
		Event: req.Event,
		
		Level: req.Level,
		
		Priority: req.Priority,
		
		ConditionExpression: req.ConditionExpression,
		
		IsActive: req.IsActive,
		
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
		return nil, fmt.Errorf("failed to create posting_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created posting_rules",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a posting_rules by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*PostingRulesResponse, error) {
	s.logger.Debug("getting posting_rules",
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
		return nil, fmt.Errorf("failed to get posting_rules: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of posting_rules records
func (s *Service) List(ctx context.Context, page, limit int) (*PostingRulesListResponse, error) {
	s.logger.Debug("listing posting_rules",
		
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
		return nil, fmt.Errorf("failed to list posting_rules: %w", err)
	}

	// Convert to response
	items := make([]*PostingRulesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PostingRulesListResponse{
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

// Update updates an existing posting_rules
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdatePostingRulesRequest) (*PostingRulesResponse, error) {
	s.logger.Info("updating posting_rules",
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
		return nil, fmt.Errorf("failed to get posting_rules: %w", err)
	}

	

	// Update fields
	
	if req.PostingProfileDocumentId != nil {
		entity.PostingProfileDocumentId = req.PostingProfileDocumentId
	}
	
	if req.RuleCode != nil {
		entity.RuleCode = req.RuleCode
	}
	
	if req.RuleName != nil {
		entity.RuleName = req.RuleName
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Event != nil {
		entity.Event = req.Event
	}
	
	if req.Level != nil {
		entity.Level = req.Level
	}
	
	if req.Priority != nil {
		entity.Priority = req.Priority
	}
	
	if req.ConditionExpression != nil {
		entity.ConditionExpression = req.ConditionExpression
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
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
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update posting_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated posting_rules",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a posting_rules
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting posting_rules",
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
		return fmt.Errorf("failed to delete posting_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted posting_rules",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PostingRules) *PostingRulesResponse {
	return &PostingRulesResponse{
		
		Id: entity.Id,
		
		PostingProfileDocumentId: entity.PostingProfileDocumentId,
		
		RuleCode: entity.RuleCode,
		
		RuleName: entity.RuleName,
		
		Description: entity.Description,
		
		Event: entity.Event,
		
		Level: entity.Level,
		
		Priority: entity.Priority,
		
		ConditionExpression: entity.ConditionExpression,
		
		IsActive: entity.IsActive,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
	}
}



// validateBusinessRules validates business rules for posting_rules
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PostingRules) error {
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

// canDelete checks if a posting_rules can be deleted
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

package posting_rule_line

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for PostingRuleLines
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PostingRuleLines service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new posting_rule_lines
func (s *Service) Create(ctx context.Context, req *CreatePostingRuleLinesRequest) (*PostingRuleLinesResponse, error) {
	s.logger.Info("creating posting_rule_lines",
		
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
	entity := &PostingRuleLines{
		
		
		PostingRuleId: req.PostingRuleId,
		
		LineNo: req.LineNo,
		
		Side: req.Side,
		
		ConceptKey: req.ConceptKey,
		
		AccountSource: req.AccountSource,
		
		FixedAccountId: req.FixedAccountId,
		
		AccountFieldPath: req.AccountFieldPath,
		
		AccountExpression: req.AccountExpression,
		
		AmountSource: req.AmountSource,
		
		AmountFieldPath: req.AmountFieldPath,
		
		AmountExpression: req.AmountExpression,
		
		MappingContext: req.MappingContext,
		
		DescriptionTemplate: req.DescriptionTemplate,
		
		IsActive: req.IsActive,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		(accountSource: req.(accountSource,
		
		(accountSource: req.(accountSource,
		
		(accountSource: req.(accountSource,
		
		(accountSource: req.(accountSource,
		
		(amountSource: req.(amountSource,
		
		(amountSource: req.(amountSource,
		
		(amountSource: req.(amountSource,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create posting_rule_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created posting_rule_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a posting_rule_lines by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*PostingRuleLinesResponse, error) {
	s.logger.Debug("getting posting_rule_lines",
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
		return nil, fmt.Errorf("failed to get posting_rule_lines: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of posting_rule_lines records
func (s *Service) List(ctx context.Context, page, limit int) (*PostingRuleLinesListResponse, error) {
	s.logger.Debug("listing posting_rule_lines",
		
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
		return nil, fmt.Errorf("failed to list posting_rule_lines: %w", err)
	}

	// Convert to response
	items := make([]*PostingRuleLinesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PostingRuleLinesListResponse{
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

// Update updates an existing posting_rule_lines
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdatePostingRuleLinesRequest) (*PostingRuleLinesResponse, error) {
	s.logger.Info("updating posting_rule_lines",
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
		return nil, fmt.Errorf("failed to get posting_rule_lines: %w", err)
	}

	

	// Update fields
	
	if req.PostingRuleId != nil {
		entity.PostingRuleId = req.PostingRuleId
	}
	
	if req.LineNo != nil {
		entity.LineNo = req.LineNo
	}
	
	if req.Side != nil {
		entity.Side = req.Side
	}
	
	if req.ConceptKey != nil {
		entity.ConceptKey = req.ConceptKey
	}
	
	if req.AccountSource != nil {
		entity.AccountSource = req.AccountSource
	}
	
	if req.FixedAccountId != nil {
		entity.FixedAccountId = req.FixedAccountId
	}
	
	if req.AccountFieldPath != nil {
		entity.AccountFieldPath = req.AccountFieldPath
	}
	
	if req.AccountExpression != nil {
		entity.AccountExpression = req.AccountExpression
	}
	
	if req.AmountSource != nil {
		entity.AmountSource = req.AmountSource
	}
	
	if req.AmountFieldPath != nil {
		entity.AmountFieldPath = req.AmountFieldPath
	}
	
	if req.AmountExpression != nil {
		entity.AmountExpression = req.AmountExpression
	}
	
	if req.MappingContext != nil {
		entity.MappingContext = req.MappingContext
	}
	
	if req.DescriptionTemplate != nil {
		entity.DescriptionTemplate = req.DescriptionTemplate
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
	
	if req.(accountSource != nil {
		entity.(accountSource = *req.(accountSource
	}
	
	if req.(accountSource != nil {
		entity.(accountSource = *req.(accountSource
	}
	
	if req.(accountSource != nil {
		entity.(accountSource = *req.(accountSource
	}
	
	if req.(accountSource != nil {
		entity.(accountSource = *req.(accountSource
	}
	
	if req.(amountSource != nil {
		entity.(amountSource = *req.(amountSource
	}
	
	if req.(amountSource != nil {
		entity.(amountSource = *req.(amountSource
	}
	
	if req.(amountSource != nil {
		entity.(amountSource = *req.(amountSource
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update posting_rule_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated posting_rule_lines",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a posting_rule_lines
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting posting_rule_lines",
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
		return fmt.Errorf("failed to delete posting_rule_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted posting_rule_lines",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PostingRuleLines) *PostingRuleLinesResponse {
	return &PostingRuleLinesResponse{
		
		Id: entity.Id,
		
		PostingRuleId: entity.PostingRuleId,
		
		LineNo: entity.LineNo,
		
		Side: entity.Side,
		
		ConceptKey: entity.ConceptKey,
		
		AccountSource: entity.AccountSource,
		
		FixedAccountId: entity.FixedAccountId,
		
		AccountFieldPath: entity.AccountFieldPath,
		
		AccountExpression: entity.AccountExpression,
		
		AmountSource: entity.AmountSource,
		
		AmountFieldPath: entity.AmountFieldPath,
		
		AmountExpression: entity.AmountExpression,
		
		MappingContext: entity.MappingContext,
		
		DescriptionTemplate: entity.DescriptionTemplate,
		
		IsActive: entity.IsActive,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		(accountSource: entity.(accountSource,
		
		(accountSource: entity.(accountSource,
		
		(accountSource: entity.(accountSource,
		
		(accountSource: entity.(accountSource,
		
		(amountSource: entity.(amountSource,
		
		(amountSource: entity.(amountSource,
		
		(amountSource: entity.(amountSource,
		
	}
}



// validateBusinessRules validates business rules for posting_rule_lines
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PostingRuleLines) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a posting_rule_lines can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

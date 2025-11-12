package budget_line

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/budget_line"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/budget_line"
	"go.uber.org/zap"
)

// Service handles business logic for BudgetLines
type Service struct {
	repo   *budget_line.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new BudgetLines service
func NewService(repo *budget_line.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new budget_lines
func (s *Service) Create(ctx context.Context, req *dto.CreateBudgetLinesRequest) (*dto.BudgetLinesResponse, error) {
	s.logger.Info("creating budget_lines",
		
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
	entity := &budget_line.BudgetLines{
		
		
		BudgetId: req.BudgetId,
		
		AccountId: req.AccountId,
		
		AnalyticAccountId: req.AnalyticAccountId,
		
		AccountingPeriodId: req.AccountingPeriodId,
		
		PeriodStartDate: req.PeriodStartDate,
		
		PeriodEndDate: req.PeriodEndDate,
		
		PlannedAmount: req.PlannedAmount,
		
		Notes: req.Notes,
		
		AccountId: req.AccountId,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create budget_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created budget_lines",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a budget_lines by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*dto.BudgetLinesResponse, error) {
	s.logger.Debug("getting budget_lines",
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
		return nil, fmt.Errorf("failed to get budget_lines: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of budget_lines records
func (s *Service) List(ctx context.Context, page, limit int) (*dto.BudgetLinesListResponse, error) {
	s.logger.Debug("listing budget_lines",
		
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
		return nil, fmt.Errorf("failed to list budget_lines: %w", err)
	}

	// Convert to response
	items := make([]*dto.BudgetLinesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.BudgetLinesListResponse{
		Items: items,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing budget_lines
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateBudgetLinesRequest) (*dto.BudgetLinesResponse, error) {
	s.logger.Info("updating budget_lines",
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
		return nil, fmt.Errorf("failed to get budget_lines: %w", err)
	}

	

	// Update fields
	
	if req.BudgetId != nil {
		entity.BudgetId = *req.BudgetId
	}
	
	if req.AccountId != nil {
		entity.AccountId = *req.AccountId
	}
	
	if req.AnalyticAccountId != nil {
		entity.AnalyticAccountId = *req.AnalyticAccountId
	}
	
	if req.AccountingPeriodId != nil {
		entity.AccountingPeriodId = *req.AccountingPeriodId
	}
	
	if req.PeriodStartDate != nil {
		entity.PeriodStartDate = *req.PeriodStartDate
	}
	
	if req.PeriodEndDate != nil {
		entity.PeriodEndDate = *req.PeriodEndDate
	}
	
	if req.PlannedAmount != nil {
		entity.PlannedAmount = *req.PlannedAmount
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.AccountId != nil {
		entity.AccountId = *req.AccountId
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update budget_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated budget_lines",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a budget_lines
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting budget_lines",
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
		return fmt.Errorf("failed to delete budget_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted budget_lines",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *budget_line.BudgetLines) *dto.BudgetLinesResponse {
	return &dto.BudgetLinesResponse{
		
		Id: entity.Id,
		
		BudgetId: entity.BudgetId,
		
		AccountId: entity.AccountId,
		
		AnalyticAccountId: entity.AnalyticAccountId,
		
		AccountingPeriodId: entity.AccountingPeriodId,
		
		PeriodStartDate: entity.PeriodStartDate,
		
		PeriodEndDate: entity.PeriodEndDate,
		
		PlannedAmount: entity.PlannedAmount,
		
		Notes: entity.Notes,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		AccountId: entity.AccountId,
		
	}
}



// validateBusinessRules validates business rules for budget_lines
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *budget_line.BudgetLines) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a budget_lines can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

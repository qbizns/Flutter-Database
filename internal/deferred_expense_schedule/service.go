package deferred_expense_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/deferred_expense_schedule"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/deferred_expense_schedule"
	"go.uber.org/zap"
)

// Service handles business logic for DeferredExpenseSchedule
type Service struct {
	repo   *deferred_expense_schedule.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new DeferredExpenseSchedule service
func NewService(repo *deferred_expense_schedule.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new deferred_expense_schedule
func (s *Service) Create(ctx context.Context, req *dto.CreateDeferredExpenseScheduleRequest) (*dto.DeferredExpenseScheduleResponse, error) {
	s.logger.Info("creating deferred_expense_schedule",
		
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
	entity := &deferred_expense_schedule.DeferredExpenseSchedule{
		
		
		ContractId: req.ContractId,
		
		LineNumber: req.LineNumber,
		
		RecognitionDate: req.RecognitionDate,
		
		RecognitionAmount: req.RecognitionAmount,
		
		Status: req.Status,
		
		JournalEntryId: req.JournalEntryId,
		
		PostedAt: req.PostedAt,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create deferred_expense_schedule: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created deferred_expense_schedule",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a deferred_expense_schedule by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*dto.DeferredExpenseScheduleResponse, error) {
	s.logger.Debug("getting deferred_expense_schedule",
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
		return nil, fmt.Errorf("failed to get deferred_expense_schedule: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of deferred_expense_schedule records
func (s *Service) List(ctx context.Context, page, limit int) (*dto.DeferredExpenseScheduleListResponse, error) {
	s.logger.Debug("listing deferred_expense_schedule",
		
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
		return nil, fmt.Errorf("failed to list deferred_expense_schedule: %w", err)
	}

	// Convert to response
	items := make([]*dto.DeferredExpenseScheduleResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.DeferredExpenseScheduleListResponse{
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

// Update updates an existing deferred_expense_schedule
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDeferredExpenseScheduleRequest) (*dto.DeferredExpenseScheduleResponse, error) {
	s.logger.Info("updating deferred_expense_schedule",
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
		return nil, fmt.Errorf("failed to get deferred_expense_schedule: %w", err)
	}

	

	// Update fields
	
	if req.ContractId != nil {
		entity.ContractId = *req.ContractId
	}
	
	if req.LineNumber != nil {
		entity.LineNumber = *req.LineNumber
	}
	
	if req.RecognitionDate != nil {
		entity.RecognitionDate = *req.RecognitionDate
	}
	
	if req.RecognitionAmount != nil {
		entity.RecognitionAmount = *req.RecognitionAmount
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.JournalEntryId != nil {
		entity.JournalEntryId = *req.JournalEntryId
	}
	
	if req.PostedAt != nil {
		entity.PostedAt = *req.PostedAt
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update deferred_expense_schedule: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated deferred_expense_schedule",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a deferred_expense_schedule
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting deferred_expense_schedule",
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
		return fmt.Errorf("failed to delete deferred_expense_schedule: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted deferred_expense_schedule",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *deferred_expense_schedule.DeferredExpenseSchedule) *dto.DeferredExpenseScheduleResponse {
	return &dto.DeferredExpenseScheduleResponse{
		
		Id: entity.Id,
		
		ContractId: entity.ContractId,
		
		LineNumber: entity.LineNumber,
		
		RecognitionDate: entity.RecognitionDate,
		
		RecognitionAmount: entity.RecognitionAmount,
		
		Status: entity.Status,
		
		JournalEntryId: entity.JournalEntryId,
		
		CreatedAt: entity.CreatedAt,
		
		PostedAt: entity.PostedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for deferred_expense_schedule
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *deferred_expense_schedule.DeferredExpenseSchedule) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a deferred_expense_schedule can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

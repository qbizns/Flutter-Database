package rate_limit

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

// Service handles business logic for RateLimits
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new RateLimits service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new rate_limits
func (s *Service) Create(ctx context.Context, req *CreateRateLimitsRequest) (*RateLimitsResponse, error) {
	s.logger.Info("creating rate_limits",
		
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
	entity := &RateLimits{
		
		
		IdentifierType: req.IdentifierType,
		
		IdentifierValue: req.IdentifierValue,
		
		EndpointPath: req.EndpointPath,
		
		HttpMethod: req.HttpMethod,
		
		WindowStart: req.WindowStart,
		
		WindowDurationSeconds: req.WindowDurationSeconds,
		
		RequestCount: req.RequestCount,
		
		AllowedCount: req.AllowedCount,
		
		IsBlocked: req.IsBlocked,
		
		BlockedUntil: req.BlockedUntil,
		
		FirstRequestAt: req.FirstRequestAt,
		
		LastRequestAt: req.LastRequestAt,
		
		IdentifierType,: req.IdentifierType,,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create rate_limits: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created rate_limits",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a rate_limits by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*RateLimitsResponse, error) {
	s.logger.Debug("getting rate_limits",
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
		return nil, fmt.Errorf("failed to get rate_limits: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of rate_limits records
func (s *Service) List(ctx context.Context, page, limit int) (*RateLimitsListResponse, error) {
	s.logger.Debug("listing rate_limits",
		
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
		return nil, fmt.Errorf("failed to list rate_limits: %w", err)
	}

	// Convert to response
	items := make([]*RateLimitsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &RateLimitsListResponse{
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

// Update updates an existing rate_limits
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRateLimitsRequest) (*RateLimitsResponse, error) {
	s.logger.Info("updating rate_limits",
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
		return nil, fmt.Errorf("failed to get rate_limits: %w", err)
	}

	

	// Update fields
	
	if req.IdentifierType != nil {
		entity.IdentifierType = *req.IdentifierType
	}
	
	if req.IdentifierValue != nil {
		entity.IdentifierValue = *req.IdentifierValue
	}
	
	if req.EndpointPath != nil {
		entity.EndpointPath = *req.EndpointPath
	}
	
	if req.HttpMethod != nil {
		entity.HttpMethod = *req.HttpMethod
	}
	
	if req.WindowStart != nil {
		entity.WindowStart = *req.WindowStart
	}
	
	if req.WindowDurationSeconds != nil {
		entity.WindowDurationSeconds = *req.WindowDurationSeconds
	}
	
	if req.RequestCount != nil {
		entity.RequestCount = *req.RequestCount
	}
	
	if req.AllowedCount != nil {
		entity.AllowedCount = *req.AllowedCount
	}
	
	if req.IsBlocked != nil {
		entity.IsBlocked = *req.IsBlocked
	}
	
	if req.BlockedUntil != nil {
		entity.BlockedUntil = *req.BlockedUntil
	}
	
	if req.FirstRequestAt != nil {
		entity.FirstRequestAt = *req.FirstRequestAt
	}
	
	if req.LastRequestAt != nil {
		entity.LastRequestAt = *req.LastRequestAt
	}
	
	if req.IdentifierType, != nil {
		entity.IdentifierType, = *req.IdentifierType,
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update rate_limits: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated rate_limits",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a rate_limits
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting rate_limits",
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
		return fmt.Errorf("failed to delete rate_limits: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted rate_limits",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *RateLimits) *RateLimitsResponse {
	return &RateLimitsResponse{
		
		Id: entity.Id,
		
		IdentifierType: entity.IdentifierType,
		
		IdentifierValue: entity.IdentifierValue,
		
		EndpointPath: entity.EndpointPath,
		
		HttpMethod: entity.HttpMethod,
		
		WindowStart: entity.WindowStart,
		
		WindowDurationSeconds: entity.WindowDurationSeconds,
		
		RequestCount: entity.RequestCount,
		
		AllowedCount: entity.AllowedCount,
		
		IsBlocked: entity.IsBlocked,
		
		BlockedUntil: entity.BlockedUntil,
		
		FirstRequestAt: entity.FirstRequestAt,
		
		LastRequestAt: entity.LastRequestAt,
		
		IdentifierType,: entity.IdentifierType,,
		
	}
}



// validateBusinessRules validates business rules for rate_limits
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *RateLimits) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a rate_limits can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

package currency

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/currency"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/currency"
	"go.uber.org/zap"
)

// Service handles business logic for Currencies
type Service struct {
	repo   *currency.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Currencies service
func NewService(repo *currency.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new currencies
func (s *Service) Create(ctx context.Context, req *dto.CreateCurrenciesRequest) (*dto.CurrenciesResponse, error) {
	s.logger.Info("creating currencies",
		
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
	entity := &currency.Currencies{
		
		
		CurrencyCode: req.CurrencyCode,
		
		CurrencyName: req.CurrencyName,
		
		CurrencySymbol: req.CurrencySymbol,
		
		DecimalPlaces: req.DecimalPlaces,
		
		IsActive: req.IsActive,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create currencies: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created currencies",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a currencies by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*dto.CurrenciesResponse, error) {
	s.logger.Debug("getting currencies",
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
		return nil, fmt.Errorf("failed to get currencies: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of currencies records
func (s *Service) List(ctx context.Context, page, limit int) (*dto.CurrenciesListResponse, error) {
	s.logger.Debug("listing currencies",
		
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
		return nil, fmt.Errorf("failed to list currencies: %w", err)
	}

	// Convert to response
	items := make([]*dto.CurrenciesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.CurrenciesListResponse{
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

// Update updates an existing currencies
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCurrenciesRequest) (*dto.CurrenciesResponse, error) {
	s.logger.Info("updating currencies",
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
		return nil, fmt.Errorf("failed to get currencies: %w", err)
	}

	

	// Update fields
	
	if req.CurrencyCode != nil {
		entity.CurrencyCode = *req.CurrencyCode
	}
	
	if req.CurrencyName != nil {
		entity.CurrencyName = *req.CurrencyName
	}
	
	if req.CurrencySymbol != nil {
		entity.CurrencySymbol = *req.CurrencySymbol
	}
	
	if req.DecimalPlaces != nil {
		entity.DecimalPlaces = *req.DecimalPlaces
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update currencies: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated currencies",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a currencies
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting currencies",
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
		return fmt.Errorf("failed to delete currencies: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted currencies",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *currency.Currencies) *dto.CurrenciesResponse {
	return &dto.CurrenciesResponse{
		
		Id: entity.Id,
		
		CurrencyCode: entity.CurrencyCode,
		
		CurrencyName: entity.CurrencyName,
		
		CurrencySymbol: entity.CurrencySymbol,
		
		DecimalPlaces: entity.DecimalPlaces,
		
		IsActive: entity.IsActive,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for currencies
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *currency.Currencies) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a currencies can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

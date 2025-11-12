package user_role

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/user_role"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/user_role"
	"go.uber.org/zap"
)

// Service handles business logic for UserRoles
type Service struct {
	repo   *user_role.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new UserRoles service
func NewService(repo *user_role.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new user_roles
func (s *Service) Create(ctx context.Context, req *dto.CreateUserRolesRequest) (*dto.UserRolesResponse, error) {
	s.logger.Info("creating user_roles",
		
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
	entity := &user_role.UserRoles{
		
		
		UserId: req.UserId,
		
		RoleId: req.RoleId,
		
		AssignedBy: req.AssignedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create user_roles: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created user_roles",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a user_roles by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserRolesResponse, error) {
	s.logger.Debug("getting user_roles",
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
		return nil, fmt.Errorf("failed to get user_roles: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of user_roles records
func (s *Service) List(ctx context.Context, page, limit int) (*dto.UserRolesListResponse, error) {
	s.logger.Debug("listing user_roles",
		
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
		return nil, fmt.Errorf("failed to list user_roles: %w", err)
	}

	// Convert to response
	items := make([]*dto.UserRolesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.UserRolesListResponse{
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

// Update updates an existing user_roles
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRolesRequest) (*dto.UserRolesResponse, error) {
	s.logger.Info("updating user_roles",
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
		return nil, fmt.Errorf("failed to get user_roles: %w", err)
	}

	

	// Update fields
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.RoleId != nil {
		entity.RoleId = *req.RoleId
	}
	
	if req.AssignedBy != nil {
		entity.AssignedBy = *req.AssignedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update user_roles: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated user_roles",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a user_roles
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting user_roles",
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
		return fmt.Errorf("failed to delete user_roles: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted user_roles",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *user_role.UserRoles) *dto.UserRolesResponse {
	return &dto.UserRolesResponse{
		
		Id: entity.Id,
		
		UserId: entity.UserId,
		
		RoleId: entity.RoleId,
		
		CreatedAt: entity.CreatedAt,
		
		AssignedBy: entity.AssignedBy,
		
	}
}



// validateBusinessRules validates business rules for user_roles
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *user_role.UserRoles) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a user_roles can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

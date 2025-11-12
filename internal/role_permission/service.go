package role_permission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/role_permission"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/role_permission"
	"go.uber.org/zap"
)

// Service handles business logic for RolePermissions
type Service struct {
	repo   *role_permission.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new RolePermissions service
func NewService(repo *role_permission.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new role_permissions
func (s *Service) Create(ctx context.Context, req *dto.CreateRolePermissionsRequest) (*dto.RolePermissionsResponse, error) {
	s.logger.Info("creating role_permissions",
		
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
	entity := &role_permission.RolePermissions{
		
		
		RoleId: req.RoleId,
		
		PermissionId: req.PermissionId,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create role_permissions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created role_permissions",
		zap.String("id", entity.ID.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a role_permissions by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*dto.RolePermissionsResponse, error) {
	s.logger.Debug("getting role_permissions",
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
		return nil, fmt.Errorf("failed to get role_permissions: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of role_permissions records
func (s *Service) List(ctx context.Context, page, limit int) (*dto.RolePermissionsListResponse, error) {
	s.logger.Debug("listing role_permissions",
		
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
		return nil, fmt.Errorf("failed to list role_permissions: %w", err)
	}

	// Convert to response
	items := make([]*dto.RolePermissionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.RolePermissionsListResponse{
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

// Update updates an existing role_permissions
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateRolePermissionsRequest) (*dto.RolePermissionsResponse, error) {
	s.logger.Info("updating role_permissions",
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
		return nil, fmt.Errorf("failed to get role_permissions: %w", err)
	}

	

	// Update fields
	
	if req.RoleId != nil {
		entity.RoleId = *req.RoleId
	}
	
	if req.PermissionId != nil {
		entity.PermissionId = *req.PermissionId
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update role_permissions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated role_permissions",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a role_permissions
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting role_permissions",
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
		return fmt.Errorf("failed to delete role_permissions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted role_permissions",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *role_permission.RolePermissions) *dto.RolePermissionsResponse {
	return &dto.RolePermissionsResponse{
		
		Id: entity.Id,
		
		RoleId: entity.RoleId,
		
		PermissionId: entity.PermissionId,
		
		CreatedAt: entity.CreatedAt,
		
	}
}



// validateBusinessRules validates business rules for role_permissions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *role_permission.RolePermissions) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a role_permissions can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

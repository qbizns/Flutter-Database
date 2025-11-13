package user

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

// Service handles business logic for Users
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Users service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new users
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateUsersRequest) (*UsersResponse, error) {
	s.logger.Info("creating users",
		zap.String("organization_id", orgID.String()),
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

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Convert DTO to entity
	entity := &Users{
		OrganizationID: orgID,
		
		Email: req.Email,
		
		PasswordHash: req.PasswordHash,
		
		FirstName: req.FirstName,
		
		LastName: req.LastName,
		
		Phone: req.Phone,
		
		AvatarUrl: req.AvatarUrl,
		
		Status: req.Status,
		
		EmailVerified: req.EmailVerified,
		
		EmailVerifiedAt: req.EmailVerifiedAt,
		
		LastLoginAt: req.LastLoginAt,
		
		LastLoginIp: req.LastLoginIp,
		
		FailedLoginAttempts: req.FailedLoginAttempts,
		
		LockedUntil: req.LockedUntil,
		
		TwoFactorEnabled: req.TwoFactorEnabled,
		
		TwoFactorSecret: req.TwoFactorSecret,
		
		Settings: req.Settings,
		
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
		return nil, fmt.Errorf("failed to create users: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created users",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a users by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*UsersResponse, error) {
	s.logger.Debug("getting users",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("users not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of users records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*UsersListResponse, error) {
	s.logger.Debug("listing users",
		zap.String("organization_id", orgID.String()),
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

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	// Get from database (organization-scoped)
	entities, total, err := s.repo.ListByOrganization(ctx, tx, orgID, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to response
	items := make([]*UsersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &UsersListResponse{
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

// Update updates an existing users
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateUsersRequest) (*UsersResponse, error) {
	s.logger.Info("updating users",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
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

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("users not found or access denied")
	}
	

	// Update fields
	
	if req.Email != nil {
		entity.Email = *req.Email
	}
	
	if req.PasswordHash != nil {
		entity.PasswordHash = *req.PasswordHash
	}
	
	if req.FirstName != nil {
		entity.FirstName = *req.FirstName
	}
	
	if req.LastName != nil {
		entity.LastName = *req.LastName
	}
	
	if req.Phone != nil {
		entity.Phone = *req.Phone
	}
	
	if req.AvatarUrl != nil {
		entity.AvatarUrl = *req.AvatarUrl
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.EmailVerified != nil {
		entity.EmailVerified = *req.EmailVerified
	}
	
	if req.EmailVerifiedAt != nil {
		entity.EmailVerifiedAt = *req.EmailVerifiedAt
	}
	
	if req.LastLoginAt != nil {
		entity.LastLoginAt = *req.LastLoginAt
	}
	
	if req.LastLoginIp != nil {
		entity.LastLoginIp = *req.LastLoginIp
	}
	
	if req.FailedLoginAttempts != nil {
		entity.FailedLoginAttempts = *req.FailedLoginAttempts
	}
	
	if req.LockedUntil != nil {
		entity.LockedUntil = *req.LockedUntil
	}
	
	if req.TwoFactorEnabled != nil {
		entity.TwoFactorEnabled = *req.TwoFactorEnabled
	}
	
	if req.TwoFactorSecret != nil {
		entity.TwoFactorSecret = *req.TwoFactorSecret
	}
	
	if req.Settings != nil {
		entity.Settings = *req.Settings
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = *req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update users: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated users",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a users
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting users",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return err
	}

	// Verify ownership
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("users not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted users",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Users) *UsersResponse {
	return &UsersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		Email: entity.Email,
		
		PasswordHash: entity.PasswordHash,
		
		FirstName: entity.FirstName,
		
		LastName: entity.LastName,
		
		Phone: entity.Phone,
		
		AvatarUrl: entity.AvatarUrl,
		
		Status: entity.Status,
		
		EmailVerified: entity.EmailVerified,
		
		EmailVerifiedAt: entity.EmailVerifiedAt,
		
		LastLoginAt: entity.LastLoginAt,
		
		LastLoginIp: entity.LastLoginIp,
		
		FailedLoginAttempts: entity.FailedLoginAttempts,
		
		LockedUntil: entity.LockedUntil,
		
		TwoFactorEnabled: entity.TwoFactorEnabled,
		
		TwoFactorSecret: entity.TwoFactorSecret,
		
		Settings: entity.Settings,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
	}
}


// setOrganizationContext sets the organization context for RLS
func (s *Service) setOrganizationContext(ctx context.Context, tx pgx.Tx, orgID uuid.UUID) error {
	_, err := tx.Exec(ctx, "SET LOCAL app.current_organization_id = $1", orgID)
	if err != nil {
		return fmt.Errorf("failed to set organization context: %w", err)
	}
	return nil
}


// validateBusinessRules validates business rules for users
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Users) error {
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

// canDelete checks if a users can be deleted
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

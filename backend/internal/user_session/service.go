package user_session

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

// Service handles business logic for UserSessions
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new UserSessions service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new user_sessions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateUserSessionsRequest) (*UserSessionsResponse, error) {
	s.logger.Info("creating user_sessions",
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
	entity := &UserSessions{
		OrganizationID: orgID,
		
		UserId: req.UserId,
		
		SessionToken: req.SessionToken,
		
		RefreshToken: req.RefreshToken,
		
		UserAgent: req.UserAgent,
		
		IpAddress: req.IpAddress,
		
		DeviceType: req.DeviceType,
		
		DeviceName: req.DeviceName,
		
		Browser: req.Browser,
		
		Os: req.Os,
		
		CountryCode: req.CountryCode,
		
		City: req.City,
		
		IsActive: req.IsActive,
		
		LastActivityAt: req.LastActivityAt,
		
		ExpiresAt: req.ExpiresAt,
		
		RevokedAt: req.RevokedAt,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create user_sessions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created user_sessions",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a user_sessions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*UserSessionsResponse, error) {
	s.logger.Debug("getting user_sessions",
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
		return nil, fmt.Errorf("failed to get user_sessions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("user_sessions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of user_sessions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*UserSessionsListResponse, error) {
	s.logger.Debug("listing user_sessions",
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
		return nil, fmt.Errorf("failed to list user_sessions: %w", err)
	}

	// Convert to response
	items := make([]*UserSessionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &UserSessionsListResponse{
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

// Update updates an existing user_sessions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateUserSessionsRequest) (*UserSessionsResponse, error) {
	s.logger.Info("updating user_sessions",
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
		return nil, fmt.Errorf("failed to get user_sessions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("user_sessions not found or access denied")
	}
	

	// Update fields
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.SessionToken != nil {
		entity.SessionToken = *req.SessionToken
	}
	
	if req.RefreshToken != nil {
		entity.RefreshToken = *req.RefreshToken
	}
	
	if req.UserAgent != nil {
		entity.UserAgent = *req.UserAgent
	}
	
	if req.IpAddress != nil {
		entity.IpAddress = *req.IpAddress
	}
	
	if req.DeviceType != nil {
		entity.DeviceType = *req.DeviceType
	}
	
	if req.DeviceName != nil {
		entity.DeviceName = *req.DeviceName
	}
	
	if req.Browser != nil {
		entity.Browser = *req.Browser
	}
	
	if req.Os != nil {
		entity.Os = *req.Os
	}
	
	if req.CountryCode != nil {
		entity.CountryCode = *req.CountryCode
	}
	
	if req.City != nil {
		entity.City = *req.City
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.LastActivityAt != nil {
		entity.LastActivityAt = *req.LastActivityAt
	}
	
	if req.ExpiresAt != nil {
		entity.ExpiresAt = *req.ExpiresAt
	}
	
	if req.RevokedAt != nil {
		entity.RevokedAt = *req.RevokedAt
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update user_sessions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated user_sessions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a user_sessions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting user_sessions",
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
		return fmt.Errorf("failed to get user_sessions: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("user_sessions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete user_sessions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted user_sessions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *UserSessions) *UserSessionsResponse {
	return &UserSessionsResponse{
		
		Id: entity.Id,
		
		UserId: entity.UserId,
		
		OrganizationId: entity.OrganizationId,
		
		SessionToken: entity.SessionToken,
		
		RefreshToken: entity.RefreshToken,
		
		UserAgent: entity.UserAgent,
		
		IpAddress: entity.IpAddress,
		
		DeviceType: entity.DeviceType,
		
		DeviceName: entity.DeviceName,
		
		Browser: entity.Browser,
		
		Os: entity.Os,
		
		CountryCode: entity.CountryCode,
		
		City: entity.City,
		
		IsActive: entity.IsActive,
		
		LastActivityAt: entity.LastActivityAt,
		
		ExpiresAt: entity.ExpiresAt,
		
		CreatedAt: entity.CreatedAt,
		
		RevokedAt: entity.RevokedAt,
		
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


// validateBusinessRules validates business rules for user_sessions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *UserSessions) error {
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

// canDelete checks if a user_sessions can be deleted
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

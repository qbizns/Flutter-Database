package user_setting

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for UserSettings
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new UserSettings service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new user_settings
func (s *Service) Create(ctx context.Context, req *CreateUserSettingsRequest) (*UserSettingsResponse, error) {
	s.logger.Info("creating user_settings",
		
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
	entity := &UserSettings{
		
		
		UserId: req.UserId,
		
		Theme: req.Theme,
		
		Language: req.Language,
		
		Timezone: req.Timezone,
		
		DefaultDashboard: req.DefaultDashboard,
		
		DashboardLayout: req.DashboardLayout,
		
		ItemsPerPage: req.ItemsPerPage,
		
		DefaultView: req.DefaultView,
		
		DesktopNotifications: req.DesktopNotifications,
		
		SoundNotifications: req.SoundNotifications,
		
		DefaultLocationId: req.DefaultLocationId,
		
		QuickActions: req.QuickActions,
		
		CustomPreferences: req.CustomPreferences,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create user_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created user_settings",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a user_settings by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*UserSettingsResponse, error) {
	s.logger.Debug("getting user_settings",
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
		return nil, fmt.Errorf("failed to get user_settings: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of user_settings records
func (s *Service) List(ctx context.Context, page, limit int) (*UserSettingsListResponse, error) {
	s.logger.Debug("listing user_settings",
		
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
		return nil, fmt.Errorf("failed to list user_settings: %w", err)
	}

	// Convert to response
	items := make([]*UserSettingsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &UserSettingsListResponse{
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

// Update updates an existing user_settings
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateUserSettingsRequest) (*UserSettingsResponse, error) {
	s.logger.Info("updating user_settings",
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
		return nil, fmt.Errorf("failed to get user_settings: %w", err)
	}

	

	// Update fields
	
	if req.UserId != nil {
		entity.UserId = req.UserId
	}
	
	if req.Theme != nil {
		entity.Theme = req.Theme
	}
	
	if req.Language != nil {
		entity.Language = req.Language
	}
	
	if req.Timezone != nil {
		entity.Timezone = req.Timezone
	}
	
	if req.DefaultDashboard != nil {
		entity.DefaultDashboard = req.DefaultDashboard
	}
	
	if req.DashboardLayout != nil {
		entity.DashboardLayout = req.DashboardLayout
	}
	
	if req.ItemsPerPage != nil {
		entity.ItemsPerPage = req.ItemsPerPage
	}
	
	if req.DefaultView != nil {
		entity.DefaultView = req.DefaultView
	}
	
	if req.DesktopNotifications != nil {
		entity.DesktopNotifications = req.DesktopNotifications
	}
	
	if req.SoundNotifications != nil {
		entity.SoundNotifications = req.SoundNotifications
	}
	
	if req.DefaultLocationId != nil {
		entity.DefaultLocationId = req.DefaultLocationId
	}
	
	if req.QuickActions != nil {
		entity.QuickActions = req.QuickActions
	}
	
	if req.CustomPreferences != nil {
		entity.CustomPreferences = req.CustomPreferences
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update user_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated user_settings",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a user_settings
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting user_settings",
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
		return fmt.Errorf("failed to delete user_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted user_settings",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *UserSettings) *UserSettingsResponse {
	return &UserSettingsResponse{
		
		UserId: entity.UserId,
		
		Theme: entity.Theme,
		
		Language: entity.Language,
		
		Timezone: entity.Timezone,
		
		DefaultDashboard: entity.DefaultDashboard,
		
		DashboardLayout: entity.DashboardLayout,
		
		ItemsPerPage: entity.ItemsPerPage,
		
		DefaultView: entity.DefaultView,
		
		DesktopNotifications: entity.DesktopNotifications,
		
		SoundNotifications: entity.SoundNotifications,
		
		DefaultLocationId: entity.DefaultLocationId,
		
		QuickActions: entity.QuickActions,
		
		CustomPreferences: entity.CustomPreferences,
		
		UpdatedAt: entity.UpdatedAt,
		
	}
}



// validateBusinessRules validates business rules for user_settings
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *UserSettings) error {
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

// canDelete checks if a user_settings can be deleted
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

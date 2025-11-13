package integration_config

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for IntegrationConfigs
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new IntegrationConfigs service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new integration_configs
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateIntegrationConfigsRequest) (*IntegrationConfigsResponse, error) {
	s.logger.Info("creating integration_configs",
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
	entity := &IntegrationConfigs{
		OrganizationId: orgID,
		
		IntegrationType: req.IntegrationType,
		
		ProviderName: req.ProviderName,
		
		Credentials: req.Credentials,
		
		Settings: req.Settings,
		
		IsActive: req.IsActive,
		
		IsConnected: req.IsConnected,
		
		ConnectionStatus: req.ConnectionStatus,
		
		LastSyncAt: req.LastSyncAt,
		
		LastSyncStatus: req.LastSyncStatus,
		
		SyncFrequency: req.SyncFrequency,
		
		WebhookUrl: req.WebhookUrl,
		
		WebhookSecret: req.WebhookSecret,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create integration_configs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created integration_configs",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a integration_configs by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*IntegrationConfigsResponse, error) {
	s.logger.Debug("getting integration_configs",
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
		return nil, fmt.Errorf("failed to get integration_configs: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("integration_configs not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of integration_configs records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*IntegrationConfigsListResponse, error) {
	s.logger.Debug("listing integration_configs",
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
		return nil, fmt.Errorf("failed to list integration_configs: %w", err)
	}

	// Convert to response
	items := make([]*IntegrationConfigsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &IntegrationConfigsListResponse{
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

// Update updates an existing integration_configs
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateIntegrationConfigsRequest) (*IntegrationConfigsResponse, error) {
	s.logger.Info("updating integration_configs",
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
		return nil, fmt.Errorf("failed to get integration_configs: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("integration_configs not found or access denied")
	}
	

	// Update fields
	
	if req.IntegrationType != nil {
		entity.IntegrationType = *req.IntegrationType
	}
	
	if req.ProviderName != nil {
		entity.ProviderName = *req.ProviderName
	}
	
	if req.Credentials != nil {
		entity.Credentials = *req.Credentials
	}
	
	if req.Settings != nil {
		entity.Settings = req.Settings
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.IsConnected != nil {
		entity.IsConnected = req.IsConnected
	}
	
	if req.ConnectionStatus != nil {
		entity.ConnectionStatus = *req.ConnectionStatus
	}
	
	if req.LastSyncAt != nil {
		entity.LastSyncAt = req.LastSyncAt
	}
	
	if req.LastSyncStatus != nil {
		entity.LastSyncStatus = *req.LastSyncStatus
	}
	
	if req.SyncFrequency != nil {
		entity.SyncFrequency = *req.SyncFrequency
	}
	
	if req.WebhookUrl != nil {
		entity.WebhookUrl = *req.WebhookUrl
	}
	
	if req.WebhookSecret != nil {
		entity.WebhookSecret = *req.WebhookSecret
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update integration_configs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated integration_configs",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a integration_configs
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting integration_configs",
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
		return fmt.Errorf("failed to get integration_configs: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("integration_configs not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete integration_configs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted integration_configs",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *IntegrationConfigs) *IntegrationConfigsResponse {
	return &IntegrationConfigsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		IntegrationType: entity.IntegrationType,
		
		ProviderName: entity.ProviderName,
		
		Credentials: entity.Credentials,
		
		Settings: entity.Settings,
		
		IsActive: entity.IsActive,
		
		IsConnected: entity.IsConnected,
		
		ConnectionStatus: entity.ConnectionStatus,
		
		LastSyncAt: entity.LastSyncAt,
		
		LastSyncStatus: entity.LastSyncStatus,
		
		SyncFrequency: entity.SyncFrequency,
		
		WebhookUrl: entity.WebhookUrl,
		
		WebhookSecret: entity.WebhookSecret,
		
		CreatedBy: entity.CreatedBy,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
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


// validateBusinessRules validates business rules for integration_configs
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *IntegrationConfigs) error {
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

// canDelete checks if a integration_configs can be deleted
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

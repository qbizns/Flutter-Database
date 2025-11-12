package kitchen_station

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

// Service handles business logic for KitchenStations
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new KitchenStations service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new kitchen_stations
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateKitchenStationsRequest) (*KitchenStationsResponse, error) {
	s.logger.Info("creating kitchen_stations",
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
	entity := &KitchenStations{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		StationName: req.StationName,
		
		StationCode: req.StationCode,
		
		StationType: req.StationType,
		
		Description: req.Description,
		
		DisplayOrder: req.DisplayOrder,
		
		ColorCode: req.ColorCode,
		
		PrinterId: req.PrinterId,
		
		IsActive: req.IsActive,
		
		AutoPrintTickets: req.AutoPrintTickets,
		
		AlertSoundEnabled: req.AlertSoundEnabled,
		
		DisplayConfig: req.DisplayConfig,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create kitchen_stations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created kitchen_stations",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a kitchen_stations by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*KitchenStationsResponse, error) {
	s.logger.Debug("getting kitchen_stations",
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
		return nil, fmt.Errorf("failed to get kitchen_stations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("kitchen_stations not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of kitchen_stations records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*KitchenStationsListResponse, error) {
	s.logger.Debug("listing kitchen_stations",
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
		return nil, fmt.Errorf("failed to list kitchen_stations: %w", err)
	}

	// Convert to response
	items := make([]*KitchenStationsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &KitchenStationsListResponse{
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

// Update updates an existing kitchen_stations
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateKitchenStationsRequest) (*KitchenStationsResponse, error) {
	s.logger.Info("updating kitchen_stations",
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
		return nil, fmt.Errorf("failed to get kitchen_stations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("kitchen_stations not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.StationName != nil {
		entity.StationName = *req.StationName
	}
	
	if req.StationCode != nil {
		entity.StationCode = *req.StationCode
	}
	
	if req.StationType != nil {
		entity.StationType = *req.StationType
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.DisplayOrder != nil {
		entity.DisplayOrder = *req.DisplayOrder
	}
	
	if req.ColorCode != nil {
		entity.ColorCode = *req.ColorCode
	}
	
	if req.PrinterId != nil {
		entity.PrinterId = *req.PrinterId
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.AutoPrintTickets != nil {
		entity.AutoPrintTickets = *req.AutoPrintTickets
	}
	
	if req.AlertSoundEnabled != nil {
		entity.AlertSoundEnabled = *req.AlertSoundEnabled
	}
	
	if req.DisplayConfig != nil {
		entity.DisplayConfig = *req.DisplayConfig
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
		return nil, fmt.Errorf("failed to update kitchen_stations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated kitchen_stations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a kitchen_stations
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting kitchen_stations",
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
		return fmt.Errorf("failed to get kitchen_stations: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("kitchen_stations not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete kitchen_stations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted kitchen_stations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *KitchenStations) *KitchenStationsResponse {
	return &KitchenStationsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		StationName: entity.StationName,
		
		StationCode: entity.StationCode,
		
		StationType: entity.StationType,
		
		Description: entity.Description,
		
		DisplayOrder: entity.DisplayOrder,
		
		ColorCode: entity.ColorCode,
		
		PrinterId: entity.PrinterId,
		
		IsActive: entity.IsActive,
		
		AutoPrintTickets: entity.AutoPrintTickets,
		
		AlertSoundEnabled: entity.AlertSoundEnabled,
		
		DisplayConfig: entity.DisplayConfig,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for kitchen_stations
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *KitchenStations) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a kitchen_stations can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

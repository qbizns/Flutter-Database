package delivery_zone

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

// Service handles business logic for DeliveryZones
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new DeliveryZones service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new delivery_zones
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateDeliveryZonesRequest) (*DeliveryZonesResponse, error) {
	s.logger.Info("creating delivery_zones",
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
	entity := &DeliveryZones{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		ZoneName: req.ZoneName,
		
		ZoneCode: req.ZoneCode,
		
		Description: req.Description,
		
		Geofence: req.Geofence,
		
		PostalCodes: req.PostalCodes,
		
		CoverageNotes: req.CoverageNotes,
		
		BaseDeliveryFee: req.BaseDeliveryFee,
		
		FeeType: req.FeeType,
		
		MinimumOrderAmount: req.MinimumOrderAmount,
		
		FreeDeliveryThreshold: req.FreeDeliveryThreshold,
		
		EstimatedDeliveryTimeMinutes: req.EstimatedDeliveryTimeMinutes,
		
		MaxDeliveryTimeMinutes: req.MaxDeliveryTimeMinutes,
		
		Priority: req.Priority,
		
		IsActive: req.IsActive,
		
		ActiveHours: req.ActiveHours,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		BaseDeliveryFee: req.BaseDeliveryFee,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create delivery_zones: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created delivery_zones",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a delivery_zones by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DeliveryZonesResponse, error) {
	s.logger.Debug("getting delivery_zones",
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
		return nil, fmt.Errorf("failed to get delivery_zones: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("delivery_zones not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of delivery_zones records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*DeliveryZonesListResponse, error) {
	s.logger.Debug("listing delivery_zones",
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
		return nil, fmt.Errorf("failed to list delivery_zones: %w", err)
	}

	// Convert to response
	items := make([]*DeliveryZonesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &DeliveryZonesListResponse{
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

// Update updates an existing delivery_zones
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateDeliveryZonesRequest) (*DeliveryZonesResponse, error) {
	s.logger.Info("updating delivery_zones",
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
		return nil, fmt.Errorf("failed to get delivery_zones: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("delivery_zones not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.ZoneName != nil {
		entity.ZoneName = *req.ZoneName
	}
	
	if req.ZoneCode != nil {
		entity.ZoneCode = *req.ZoneCode
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.Geofence != nil {
		entity.Geofence = *req.Geofence
	}
	
	if req.PostalCodes != nil {
		entity.PostalCodes = *req.PostalCodes
	}
	
	if req.CoverageNotes != nil {
		entity.CoverageNotes = *req.CoverageNotes
	}
	
	if req.BaseDeliveryFee != nil {
		entity.BaseDeliveryFee = *req.BaseDeliveryFee
	}
	
	if req.FeeType != nil {
		entity.FeeType = *req.FeeType
	}
	
	if req.MinimumOrderAmount != nil {
		entity.MinimumOrderAmount = *req.MinimumOrderAmount
	}
	
	if req.FreeDeliveryThreshold != nil {
		entity.FreeDeliveryThreshold = *req.FreeDeliveryThreshold
	}
	
	if req.EstimatedDeliveryTimeMinutes != nil {
		entity.EstimatedDeliveryTimeMinutes = *req.EstimatedDeliveryTimeMinutes
	}
	
	if req.MaxDeliveryTimeMinutes != nil {
		entity.MaxDeliveryTimeMinutes = *req.MaxDeliveryTimeMinutes
	}
	
	if req.Priority != nil {
		entity.Priority = *req.Priority
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.ActiveHours != nil {
		entity.ActiveHours = *req.ActiveHours
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
	
	if req.BaseDeliveryFee != nil {
		entity.BaseDeliveryFee = *req.BaseDeliveryFee
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update delivery_zones: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated delivery_zones",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a delivery_zones
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting delivery_zones",
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
		return fmt.Errorf("failed to get delivery_zones: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("delivery_zones not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete delivery_zones: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted delivery_zones",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *DeliveryZones) *DeliveryZonesResponse {
	return &DeliveryZonesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		ZoneName: entity.ZoneName,
		
		ZoneCode: entity.ZoneCode,
		
		Description: entity.Description,
		
		Geofence: entity.Geofence,
		
		PostalCodes: entity.PostalCodes,
		
		CoverageNotes: entity.CoverageNotes,
		
		BaseDeliveryFee: entity.BaseDeliveryFee,
		
		FeeType: entity.FeeType,
		
		MinimumOrderAmount: entity.MinimumOrderAmount,
		
		FreeDeliveryThreshold: entity.FreeDeliveryThreshold,
		
		EstimatedDeliveryTimeMinutes: entity.EstimatedDeliveryTimeMinutes,
		
		MaxDeliveryTimeMinutes: entity.MaxDeliveryTimeMinutes,
		
		Priority: entity.Priority,
		
		IsActive: entity.IsActive,
		
		ActiveHours: entity.ActiveHours,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		BaseDeliveryFee: entity.BaseDeliveryFee,
		
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


// validateBusinessRules validates business rules for delivery_zones
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *DeliveryZones) error {
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

// canDelete checks if a delivery_zones can be deleted
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

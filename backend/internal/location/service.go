package location

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

// Service handles business logic for Locations
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Locations service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new locations
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateLocationsRequest) (*LocationsResponse, error) {
	s.logger.Info("creating locations",
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
	entity := &Locations{
		OrganizationID: orgID,
		
		LocationCode: req.LocationCode,
		
		Name: req.Name,
		
		LocationType: req.LocationType,
		
		Phone: req.Phone,
		
		Email: req.Email,
		
		ManagerUserId: req.ManagerUserId,
		
		AddressLine1: req.AddressLine1,
		
		AddressLine2: req.AddressLine2,
		
		City: req.City,
		
		State: req.State,
		
		Country: req.Country,
		
		PostalCode: req.PostalCode,
		
		Timezone: req.Timezone,
		
		BusinessHours: req.BusinessHours,
		
		IsActive: req.IsActive,
		
		IsPrimary: req.IsPrimary,
		
		AllowSales: req.AllowSales,
		
		AllowPurchases: req.AllowPurchases,
		
		TaxRate: req.TaxRate,
		
		Notes: req.Notes,
		
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
		return nil, fmt.Errorf("failed to create locations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created locations",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a locations by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LocationsResponse, error) {
	s.logger.Debug("getting locations",
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
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("locations not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of locations records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*LocationsListResponse, error) {
	s.logger.Debug("listing locations",
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
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}

	// Convert to response
	items := make([]*LocationsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &LocationsListResponse{
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

// Update updates an existing locations
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateLocationsRequest) (*LocationsResponse, error) {
	s.logger.Info("updating locations",
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
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("locations not found or access denied")
	}
	

	// Update fields
	
	if req.LocationCode != nil {
		entity.LocationCode = *req.LocationCode
	}
	
	if req.Name != nil {
		entity.Name = *req.Name
	}
	
	if req.LocationType != nil {
		entity.LocationType = *req.LocationType
	}
	
	if req.Phone != nil {
		entity.Phone = *req.Phone
	}
	
	if req.Email != nil {
		entity.Email = *req.Email
	}
	
	if req.ManagerUserId != nil {
		entity.ManagerUserId = *req.ManagerUserId
	}
	
	if req.AddressLine1 != nil {
		entity.AddressLine1 = *req.AddressLine1
	}
	
	if req.AddressLine2 != nil {
		entity.AddressLine2 = *req.AddressLine2
	}
	
	if req.City != nil {
		entity.City = *req.City
	}
	
	if req.State != nil {
		entity.State = *req.State
	}
	
	if req.Country != nil {
		entity.Country = *req.Country
	}
	
	if req.PostalCode != nil {
		entity.PostalCode = *req.PostalCode
	}
	
	if req.Timezone != nil {
		entity.Timezone = *req.Timezone
	}
	
	if req.BusinessHours != nil {
		entity.BusinessHours = *req.BusinessHours
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.IsPrimary != nil {
		entity.IsPrimary = *req.IsPrimary
	}
	
	if req.AllowSales != nil {
		entity.AllowSales = *req.AllowSales
	}
	
	if req.AllowPurchases != nil {
		entity.AllowPurchases = *req.AllowPurchases
	}
	
	if req.TaxRate != nil {
		entity.TaxRate = *req.TaxRate
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
		return nil, fmt.Errorf("failed to update locations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated locations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a locations
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting locations",
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
		return fmt.Errorf("failed to get locations: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("locations not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete locations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted locations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Locations) *LocationsResponse {
	return &LocationsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationCode: entity.LocationCode,
		
		Name: entity.Name,
		
		LocationType: entity.LocationType,
		
		Phone: entity.Phone,
		
		Email: entity.Email,
		
		ManagerUserId: entity.ManagerUserId,
		
		AddressLine1: entity.AddressLine1,
		
		AddressLine2: entity.AddressLine2,
		
		City: entity.City,
		
		State: entity.State,
		
		Country: entity.Country,
		
		PostalCode: entity.PostalCode,
		
		Timezone: entity.Timezone,
		
		BusinessHours: entity.BusinessHours,
		
		IsActive: entity.IsActive,
		
		IsPrimary: entity.IsPrimary,
		
		AllowSales: entity.AllowSales,
		
		AllowPurchases: entity.AllowPurchases,
		
		TaxRate: entity.TaxRate,
		
		Notes: entity.Notes,
		
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


// validateBusinessRules validates business rules for locations
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Locations) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a locations can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

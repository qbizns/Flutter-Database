package restaurant_table

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/pos-backend/internal/logging"

	"go.uber.org/zap"
)

// Service handles business logic for RestaurantTables
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new RestaurantTables service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new restaurant_tables
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateRestaurantTablesRequest) (*RestaurantTablesResponse, error) {
	s.logger.Info("creating restaurant_tables",
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
	entity := &RestaurantTables{
		OrganizationId: orgID,
		
		LocationId: req.LocationId,
		
		FloorPlanId: req.FloorPlanId,
		
		SectionId: req.SectionId,
		
		TableNumber: req.TableNumber,
		
		TableName: req.TableName,
		
		MinCapacity: req.MinCapacity,
		
		MaxCapacity: req.MaxCapacity,
		
		TableShape: req.TableShape,
		
		IsCombinable: req.IsCombinable,
		
		PositionX: req.PositionX,
		
		PositionY: req.PositionY,
		
		Rotation: req.Rotation,
		
		Status: req.Status,
		
		CurrentCovers: req.CurrentCovers,
		
		SeatedAt: req.SeatedAt,
		
		CurrentWaiterId: req.CurrentWaiterId,
		
		IsActive: req.IsActive,
		
		AllowOnlineReservation: req.AllowOnlineReservation,
		
		DisplayOrder: req.DisplayOrder,
		
		ColorCode: req.ColorCode,
		
		Icon: req.Icon,
		
		Notes: req.Notes,
		
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
		return nil, fmt.Errorf("failed to create restaurant_tables: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created restaurant_tables",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a restaurant_tables by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*RestaurantTablesResponse, error) {
	s.logger.Debug("getting restaurant_tables",
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
		return nil, fmt.Errorf("failed to get restaurant_tables: %w", err)
	}


	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("restaurant_tables not found or access denied")
	}


	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of restaurant_tables records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*RestaurantTablesListResponse, error) {
	s.logger.Debug("listing restaurant_tables",
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
		return nil, fmt.Errorf("failed to list restaurant_tables: %w", err)
	}

	// Convert to response
	items := make([]*RestaurantTablesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &RestaurantTablesListResponse{
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

// Update updates an existing restaurant_tables
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateRestaurantTablesRequest) (*RestaurantTablesResponse, error) {
	s.logger.Info("updating restaurant_tables",
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
		return nil, fmt.Errorf("failed to get restaurant_tables: %w", err)
	}


	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("restaurant_tables not found or access denied")
	}


	// Update fields

	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.FloorPlanId != nil {
		entity.FloorPlanId = req.FloorPlanId
	}
	
	if req.SectionId != nil {
		entity.SectionId = req.SectionId
	}

	if req.TableNumber != nil {
		entity.TableNumber = *req.TableNumber
	}

	if req.TableName != nil {
		entity.TableName = req.TableName
	}

	if req.MinCapacity != nil {
		entity.MinCapacity = req.MinCapacity
	}

	if req.MaxCapacity != nil {
		entity.MaxCapacity = *req.MaxCapacity
	}

	if req.TableShape != nil {
		entity.TableShape = req.TableShape
	}

	if req.IsCombinable != nil {
		entity.IsCombinable = req.IsCombinable
	}

	if req.PositionX != nil {
		entity.PositionX = req.PositionX
	}

	if req.PositionY != nil {
		entity.PositionY = req.PositionY
	}

	if req.Rotation != nil {
		entity.Rotation = req.Rotation
	}

	if req.Status != nil {
		entity.Status = req.Status
	}

	if req.CurrentCovers != nil {
		entity.CurrentCovers = req.CurrentCovers
	}

	if req.SeatedAt != nil {
		entity.SeatedAt = req.SeatedAt
	}

	if req.CurrentWaiterId != nil {
		entity.CurrentWaiterId = req.CurrentWaiterId
	}

	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}

	if req.AllowOnlineReservation != nil {
		entity.AllowOnlineReservation = req.AllowOnlineReservation
	}

	if req.DisplayOrder != nil {
		entity.DisplayOrder = req.DisplayOrder
	}

	if req.ColorCode != nil {
		entity.ColorCode = req.ColorCode
	}

	if req.Icon != nil {
		entity.Icon = req.Icon
	}

	if req.Notes != nil {
		entity.Notes = req.Notes
	}

	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}

	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}

	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update restaurant_tables: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated restaurant_tables",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a restaurant_tables
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting restaurant_tables",
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
		return fmt.Errorf("failed to get restaurant_tables: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("restaurant_tables not found or access denied")
	}


	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete restaurant_tables: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted restaurant_tables",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *RestaurantTables) *RestaurantTablesResponse {
	return &RestaurantTablesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		FloorPlanId: entity.FloorPlanId,
		
		SectionId: entity.SectionId,
		
		TableNumber: entity.TableNumber,
		
		TableName: entity.TableName,
		
		MinCapacity: entity.MinCapacity,
		
		MaxCapacity: entity.MaxCapacity,
		
		TableShape: entity.TableShape,
		
		IsCombinable: entity.IsCombinable,
		
		PositionX: entity.PositionX,
		
		PositionY: entity.PositionY,
		
		Rotation: entity.Rotation,
		
		Status: entity.Status,
		
		CurrentCovers: entity.CurrentCovers,
		
		SeatedAt: entity.SeatedAt,
		
		CurrentWaiterId: entity.CurrentWaiterId,
		
		IsActive: entity.IsActive,
		
		AllowOnlineReservation: entity.AllowOnlineReservation,
		
		DisplayOrder: entity.DisplayOrder,
		
		ColorCode: entity.ColorCode,
		
		Icon: entity.Icon,
		
		Notes: entity.Notes,
		
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


// validateBusinessRules validates business rules for restaurant_tables
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *RestaurantTables) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a restaurant_tables can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

// UpdateStatus updates the status of a table
func (s *Service) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateTableStatusRequest) (*RestaurantTablesResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("table not found or access denied")
	}

	entity.Status = &req.Status

	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.entityToResponse(entity), nil
}

// AssignOrder assigns an order to a table
func (s *Service) AssignOrder(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *AssignOrderRequest) (*RestaurantTablesResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("table not found or access denied")
	}

	occupiedStatus := "occupied"
	entity.Status = &occupiedStatus
	// Note: CurrentOrderId, CustomerName, and WaiterId fields don't exist in RestaurantTables struct
	// Only update CurrentWaiterId if available
	if req.WaiterID != nil {
		entity.CurrentWaiterId = req.WaiterID
	}

	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to assign order to table: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.entityToResponse(entity), nil
}

// Clear clears a table
func (s *Service) Clear(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*RestaurantTablesResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("table not found or access denied")
	}

	availableStatus := "available"
	entity.Status = &availableStatus
	// Note: CurrentOrderId, CustomerName, CustomerId, and WaiterId fields don't exist in RestaurantTables struct
	// Only clear CurrentWaiterId
	entity.CurrentWaiterId = nil

	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to clear table: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.entityToResponse(entity), nil
}

// GetStatistics retrieves table statistics
func (s *Service) GetStatistics(ctx context.Context, orgID uuid.UUID) (*TableStatisticsResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	stats, err := s.repo.GetStatistics(ctx, tx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table statistics: %w", err)
	}

	return stats, nil
}

// GetCountByZone retrieves table count by zone
func (s *Service) GetCountByZone(ctx context.Context, orgID uuid.UUID) (map[string]int, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	counts, err := s.repo.GetCountByZone(ctx, tx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table count by zone: %w", err)
	}

	return counts, nil
}

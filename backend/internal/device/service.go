package device

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for Devices
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Devices service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new devices
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateDevicesRequest) (*DevicesResponse, error) {
	s.logger.Info("creating devices",
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
	entity := &Devices{
		OrganizationId: orgID,
		
		LocationId: req.LocationId,
		
		DeviceCode: req.DeviceCode,
		
		DeviceName: req.DeviceName,
		
		DeviceType: req.DeviceType,
		
		Manufacturer: req.Manufacturer,
		
		Model: req.Model,
		
		SerialNumber: req.SerialNumber,
		
		MacAddress: req.MacAddress,
		
		IpAddress: req.IpAddress,
		
		DeviceConfig: req.DeviceConfig,
		
		ScreenResolution: req.ScreenResolution,
		
		OsVersion: req.OsVersion,
		
		ConnectionType: req.ConnectionType,
		
		ConnectionString: req.ConnectionString,
		
		Status: req.Status,
		
		LastOnlineAt: req.LastOnlineAt,
		
		LastHeartbeatAt: req.LastHeartbeatAt,
		
		AssignedToUserId: req.AssignedToUserId,
		
		AssignedToStationId: req.AssignedToStationId,
		
		PurchaseDate: req.PurchaseDate,
		
		WarrantyExpiryDate: req.WarrantyExpiryDate,
		
		LicenseKey: req.LicenseKey,
		
		LicenseExpiryDate: req.LicenseExpiryDate,
		
		InstallationNotes: req.InstallationNotes,
		
		MaintenanceNotes: req.MaintenanceNotes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'active',: req.'active',,
		
		'posTerminal',: req.'posTerminal',,
		
		'printer',: req.'printer',,
		
		'kitchenPrinter',: req.'kitchenPrinter',,
		
		'network',: req.'network',,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create devices: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created devices",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a devices by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DevicesResponse, error) {
	s.logger.Debug("getting devices",
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
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("devices not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of devices records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*DevicesListResponse, error) {
	s.logger.Debug("listing devices",
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
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	// Convert to response
	items := make([]*DevicesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &DevicesListResponse{
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

// Update updates an existing devices
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateDevicesRequest) (*DevicesResponse, error) {
	s.logger.Info("updating devices",
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
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("devices not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.DeviceCode != nil {
		entity.DeviceCode = *req.DeviceCode
	}
	
	if req.DeviceName != nil {
		entity.DeviceName = *req.DeviceName
	}
	
	if req.DeviceType != nil {
		entity.DeviceType = *req.DeviceType
	}
	
	if req.Manufacturer != nil {
		entity.Manufacturer = *req.Manufacturer
	}
	
	if req.Model != nil {
		entity.Model = *req.Model
	}
	
	if req.SerialNumber != nil {
		entity.SerialNumber = *req.SerialNumber
	}
	
	if req.MacAddress != nil {
		entity.MacAddress = *req.MacAddress
	}
	
	if req.IpAddress != nil {
		entity.IpAddress = req.IpAddress
	}
	
	if req.DeviceConfig != nil {
		entity.DeviceConfig = *req.DeviceConfig
	}
	
	if req.ScreenResolution != nil {
		entity.ScreenResolution = *req.ScreenResolution
	}
	
	if req.OsVersion != nil {
		entity.OsVersion = *req.OsVersion
	}
	
	if req.ConnectionType != nil {
		entity.ConnectionType = *req.ConnectionType
	}
	
	if req.ConnectionString != nil {
		entity.ConnectionString = *req.ConnectionString
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.LastOnlineAt != nil {
		entity.LastOnlineAt = req.LastOnlineAt
	}
	
	if req.LastHeartbeatAt != nil {
		entity.LastHeartbeatAt = req.LastHeartbeatAt
	}
	
	if req.AssignedToUserId != nil {
		entity.AssignedToUserId = *req.AssignedToUserId
	}
	
	if req.AssignedToStationId != nil {
		entity.AssignedToStationId = *req.AssignedToStationId
	}
	
	if req.PurchaseDate != nil {
		entity.PurchaseDate = *req.PurchaseDate
	}
	
	if req.WarrantyExpiryDate != nil {
		entity.WarrantyExpiryDate = *req.WarrantyExpiryDate
	}
	
	if req.LicenseKey != nil {
		entity.LicenseKey = *req.LicenseKey
	}
	
	if req.LicenseExpiryDate != nil {
		entity.LicenseExpiryDate = *req.LicenseExpiryDate
	}
	
	if req.InstallationNotes != nil {
		entity.InstallationNotes = *req.InstallationNotes
	}
	
	if req.MaintenanceNotes != nil {
		entity.MaintenanceNotes = *req.MaintenanceNotes
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	
	if req.'active', != nil {
		entity.'active', = *req.'active',
	}
	
	if req.'posTerminal', != nil {
		entity.'posTerminal', = *req.'posTerminal',
	}
	
	if req.'printer', != nil {
		entity.'printer', = *req.'printer',
	}
	
	if req.'kitchenPrinter', != nil {
		entity.'kitchenPrinter', = *req.'kitchenPrinter',
	}
	
	if req.'network', != nil {
		entity.'network', = *req.'network',
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update devices: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated devices",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a devices
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting devices",
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
		return fmt.Errorf("failed to get devices: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("devices not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete devices: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted devices",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Devices) *DevicesResponse {
	return &DevicesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		DeviceCode: entity.DeviceCode,
		
		DeviceName: entity.DeviceName,
		
		DeviceType: entity.DeviceType,
		
		Manufacturer: entity.Manufacturer,
		
		Model: entity.Model,
		
		SerialNumber: entity.SerialNumber,
		
		MacAddress: entity.MacAddress,
		
		IpAddress: entity.IpAddress,
		
		DeviceConfig: entity.DeviceConfig,
		
		ScreenResolution: entity.ScreenResolution,
		
		OsVersion: entity.OsVersion,
		
		ConnectionType: entity.ConnectionType,
		
		ConnectionString: entity.ConnectionString,
		
		Status: entity.Status,
		
		LastOnlineAt: entity.LastOnlineAt,
		
		LastHeartbeatAt: entity.LastHeartbeatAt,
		
		AssignedToUserId: entity.AssignedToUserId,
		
		AssignedToStationId: entity.AssignedToStationId,
		
		PurchaseDate: entity.PurchaseDate,
		
		WarrantyExpiryDate: entity.WarrantyExpiryDate,
		
		LicenseKey: entity.LicenseKey,
		
		LicenseExpiryDate: entity.LicenseExpiryDate,
		
		InstallationNotes: entity.InstallationNotes,
		
		MaintenanceNotes: entity.MaintenanceNotes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'active',: entity.'active',,
		
		'posTerminal',: entity.'posTerminal',,
		
		'printer',: entity.'printer',,
		
		'kitchenPrinter',: entity.'kitchenPrinter',,
		
		'network',: entity.'network',,
		
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


// validateBusinessRules validates business rules for devices
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Devices) error {
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

// canDelete checks if a devices can be deleted
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

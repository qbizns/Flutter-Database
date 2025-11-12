package delivery_driver

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/delivery_driver"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/delivery_driver"
	"go.uber.org/zap"
)

// Service handles business logic for DeliveryDrivers
type Service struct {
	repo   *delivery_driver.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new DeliveryDrivers service
func NewService(repo *delivery_driver.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new delivery_drivers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateDeliveryDriversRequest) (*dto.DeliveryDriversResponse, error) {
	s.logger.Info("creating delivery_drivers",
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
	entity := &delivery_driver.DeliveryDrivers{
		OrganizationID: orgID,
		
		UserId: req.UserId,
		
		DriverCode: req.DriverCode,
		
		FullName: req.FullName,
		
		Phone: req.Phone,
		
		Email: req.Email,
		
		EmergencyContactName: req.EmergencyContactName,
		
		EmergencyContactPhone: req.EmergencyContactPhone,
		
		VehicleType: req.VehicleType,
		
		VehicleMake: req.VehicleMake,
		
		VehicleModel: req.VehicleModel,
		
		VehicleYear: req.VehicleYear,
		
		VehicleColor: req.VehicleColor,
		
		LicensePlate: req.LicensePlate,
		
		DriversLicenseNumber: req.DriversLicenseNumber,
		
		LicenseExpiryDate: req.LicenseExpiryDate,
		
		InsurancePolicyNumber: req.InsurancePolicyNumber,
		
		InsuranceExpiryDate: req.InsuranceExpiryDate,
		
		HireDate: req.HireDate,
		
		EmploymentType: req.EmploymentType,
		
		Status: req.Status,
		
		TotalDeliveries: req.TotalDeliveries,
		
		SuccessfulDeliveries: req.SuccessfulDeliveries,
		
		Rating: req.Rating,
		
		RatingCount: req.RatingCount,
		
		CurrentLocation: req.CurrentLocation,
		
		IsAvailable: req.IsAvailable,
		
		LastLocationUpdate: req.LastLocationUpdate,
		
		CommissionRate: req.CommissionRate,
		
		PaymentMethod: req.PaymentMethod,
		
		Documents: req.Documents,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'active',: req.'active',,
		
		'fullTime',: req.'fullTime',,
		
		'bike',: req.'bike',,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create delivery_drivers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created delivery_drivers",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a delivery_drivers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.DeliveryDriversResponse, error) {
	s.logger.Debug("getting delivery_drivers",
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
		return nil, fmt.Errorf("failed to get delivery_drivers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("delivery_drivers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of delivery_drivers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.DeliveryDriversListResponse, error) {
	s.logger.Debug("listing delivery_drivers",
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
		return nil, fmt.Errorf("failed to list delivery_drivers: %w", err)
	}

	// Convert to response
	items := make([]*dto.DeliveryDriversResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.DeliveryDriversListResponse{
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

// Update updates an existing delivery_drivers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateDeliveryDriversRequest) (*dto.DeliveryDriversResponse, error) {
	s.logger.Info("updating delivery_drivers",
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
		return nil, fmt.Errorf("failed to get delivery_drivers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("delivery_drivers not found or access denied")
	}
	

	// Update fields
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.DriverCode != nil {
		entity.DriverCode = *req.DriverCode
	}
	
	if req.FullName != nil {
		entity.FullName = *req.FullName
	}
	
	if req.Phone != nil {
		entity.Phone = *req.Phone
	}
	
	if req.Email != nil {
		entity.Email = *req.Email
	}
	
	if req.EmergencyContactName != nil {
		entity.EmergencyContactName = *req.EmergencyContactName
	}
	
	if req.EmergencyContactPhone != nil {
		entity.EmergencyContactPhone = *req.EmergencyContactPhone
	}
	
	if req.VehicleType != nil {
		entity.VehicleType = *req.VehicleType
	}
	
	if req.VehicleMake != nil {
		entity.VehicleMake = *req.VehicleMake
	}
	
	if req.VehicleModel != nil {
		entity.VehicleModel = *req.VehicleModel
	}
	
	if req.VehicleYear != nil {
		entity.VehicleYear = *req.VehicleYear
	}
	
	if req.VehicleColor != nil {
		entity.VehicleColor = *req.VehicleColor
	}
	
	if req.LicensePlate != nil {
		entity.LicensePlate = *req.LicensePlate
	}
	
	if req.DriversLicenseNumber != nil {
		entity.DriversLicenseNumber = *req.DriversLicenseNumber
	}
	
	if req.LicenseExpiryDate != nil {
		entity.LicenseExpiryDate = *req.LicenseExpiryDate
	}
	
	if req.InsurancePolicyNumber != nil {
		entity.InsurancePolicyNumber = *req.InsurancePolicyNumber
	}
	
	if req.InsuranceExpiryDate != nil {
		entity.InsuranceExpiryDate = *req.InsuranceExpiryDate
	}
	
	if req.HireDate != nil {
		entity.HireDate = *req.HireDate
	}
	
	if req.EmploymentType != nil {
		entity.EmploymentType = *req.EmploymentType
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.TotalDeliveries != nil {
		entity.TotalDeliveries = *req.TotalDeliveries
	}
	
	if req.SuccessfulDeliveries != nil {
		entity.SuccessfulDeliveries = *req.SuccessfulDeliveries
	}
	
	if req.Rating != nil {
		entity.Rating = *req.Rating
	}
	
	if req.RatingCount != nil {
		entity.RatingCount = *req.RatingCount
	}
	
	if req.CurrentLocation != nil {
		entity.CurrentLocation = *req.CurrentLocation
	}
	
	if req.IsAvailable != nil {
		entity.IsAvailable = *req.IsAvailable
	}
	
	if req.LastLocationUpdate != nil {
		entity.LastLocationUpdate = *req.LastLocationUpdate
	}
	
	if req.CommissionRate != nil {
		entity.CommissionRate = *req.CommissionRate
	}
	
	if req.PaymentMethod != nil {
		entity.PaymentMethod = *req.PaymentMethod
	}
	
	if req.Documents != nil {
		entity.Documents = *req.Documents
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
	
	if req.'active', != nil {
		entity.'active', = *req.'active',
	}
	
	if req.'fullTime', != nil {
		entity.'fullTime', = *req.'fullTime',
	}
	
	if req.'bike', != nil {
		entity.'bike', = *req.'bike',
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update delivery_drivers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated delivery_drivers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a delivery_drivers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting delivery_drivers",
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
		return fmt.Errorf("failed to get delivery_drivers: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("delivery_drivers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete delivery_drivers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted delivery_drivers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *delivery_driver.DeliveryDrivers) *dto.DeliveryDriversResponse {
	return &dto.DeliveryDriversResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		UserId: entity.UserId,
		
		DriverCode: entity.DriverCode,
		
		FullName: entity.FullName,
		
		Phone: entity.Phone,
		
		Email: entity.Email,
		
		EmergencyContactName: entity.EmergencyContactName,
		
		EmergencyContactPhone: entity.EmergencyContactPhone,
		
		VehicleType: entity.VehicleType,
		
		VehicleMake: entity.VehicleMake,
		
		VehicleModel: entity.VehicleModel,
		
		VehicleYear: entity.VehicleYear,
		
		VehicleColor: entity.VehicleColor,
		
		LicensePlate: entity.LicensePlate,
		
		DriversLicenseNumber: entity.DriversLicenseNumber,
		
		LicenseExpiryDate: entity.LicenseExpiryDate,
		
		InsurancePolicyNumber: entity.InsurancePolicyNumber,
		
		InsuranceExpiryDate: entity.InsuranceExpiryDate,
		
		HireDate: entity.HireDate,
		
		EmploymentType: entity.EmploymentType,
		
		Status: entity.Status,
		
		TotalDeliveries: entity.TotalDeliveries,
		
		SuccessfulDeliveries: entity.SuccessfulDeliveries,
		
		Rating: entity.Rating,
		
		RatingCount: entity.RatingCount,
		
		CurrentLocation: entity.CurrentLocation,
		
		IsAvailable: entity.IsAvailable,
		
		LastLocationUpdate: entity.LastLocationUpdate,
		
		CommissionRate: entity.CommissionRate,
		
		PaymentMethod: entity.PaymentMethod,
		
		Documents: entity.Documents,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'active',: entity.'active',,
		
		'fullTime',: entity.'fullTime',,
		
		'bike',: entity.'bike',,
		
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


// validateBusinessRules validates business rules for delivery_drivers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *delivery_driver.DeliveryDrivers) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a delivery_drivers can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

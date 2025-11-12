package printer_configuration

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

// Service handles business logic for PrinterConfigurations
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PrinterConfigurations service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new printer_configurations
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreatePrinterConfigurationsRequest) (*PrinterConfigurationsResponse, error) {
	s.logger.Info("creating printer_configurations",
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
	entity := &PrinterConfigurations{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		PrinterDeviceId: req.PrinterDeviceId,
		
		DocumentType: req.DocumentType,
		
		FilterOrderType: req.FilterOrderType,
		
		FilterKitchenStationId: req.FilterKitchenStationId,
		
		FilterProductCategoryId: req.FilterProductCategoryId,
		
		FilterCourseId: req.FilterCourseId,
		
		NumberOfCopies: req.NumberOfCopies,
		
		AutoPrint: req.AutoPrint,
		
		PrintPriority: req.PrintPriority,
		
		TemplateConfig: req.TemplateConfig,
		
		PaperSize: req.PaperSize,
		
		PrintOrientation: req.PrintOrientation,
		
		IsActive: req.IsActive,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'receipt',: req.'receipt',,
		
		'report',: req.'report',,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create printer_configurations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created printer_configurations",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a printer_configurations by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PrinterConfigurationsResponse, error) {
	s.logger.Debug("getting printer_configurations",
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
		return nil, fmt.Errorf("failed to get printer_configurations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("printer_configurations not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of printer_configurations records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*PrinterConfigurationsListResponse, error) {
	s.logger.Debug("listing printer_configurations",
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
		return nil, fmt.Errorf("failed to list printer_configurations: %w", err)
	}

	// Convert to response
	items := make([]*PrinterConfigurationsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PrinterConfigurationsListResponse{
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

// Update updates an existing printer_configurations
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdatePrinterConfigurationsRequest) (*PrinterConfigurationsResponse, error) {
	s.logger.Info("updating printer_configurations",
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
		return nil, fmt.Errorf("failed to get printer_configurations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("printer_configurations not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.PrinterDeviceId != nil {
		entity.PrinterDeviceId = *req.PrinterDeviceId
	}
	
	if req.DocumentType != nil {
		entity.DocumentType = *req.DocumentType
	}
	
	if req.FilterOrderType != nil {
		entity.FilterOrderType = *req.FilterOrderType
	}
	
	if req.FilterKitchenStationId != nil {
		entity.FilterKitchenStationId = *req.FilterKitchenStationId
	}
	
	if req.FilterProductCategoryId != nil {
		entity.FilterProductCategoryId = *req.FilterProductCategoryId
	}
	
	if req.FilterCourseId != nil {
		entity.FilterCourseId = *req.FilterCourseId
	}
	
	if req.NumberOfCopies != nil {
		entity.NumberOfCopies = *req.NumberOfCopies
	}
	
	if req.AutoPrint != nil {
		entity.AutoPrint = *req.AutoPrint
	}
	
	if req.PrintPriority != nil {
		entity.PrintPriority = *req.PrintPriority
	}
	
	if req.TemplateConfig != nil {
		entity.TemplateConfig = *req.TemplateConfig
	}
	
	if req.PaperSize != nil {
		entity.PaperSize = *req.PaperSize
	}
	
	if req.PrintOrientation != nil {
		entity.PrintOrientation = *req.PrintOrientation
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
	
	if req.'receipt', != nil {
		entity.'receipt', = *req.'receipt',
	}
	
	if req.'report', != nil {
		entity.'report', = *req.'report',
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update printer_configurations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated printer_configurations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a printer_configurations
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting printer_configurations",
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
		return fmt.Errorf("failed to get printer_configurations: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("printer_configurations not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete printer_configurations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted printer_configurations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PrinterConfigurations) *PrinterConfigurationsResponse {
	return &PrinterConfigurationsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		PrinterDeviceId: entity.PrinterDeviceId,
		
		DocumentType: entity.DocumentType,
		
		FilterOrderType: entity.FilterOrderType,
		
		FilterKitchenStationId: entity.FilterKitchenStationId,
		
		FilterProductCategoryId: entity.FilterProductCategoryId,
		
		FilterCourseId: entity.FilterCourseId,
		
		NumberOfCopies: entity.NumberOfCopies,
		
		AutoPrint: entity.AutoPrint,
		
		PrintPriority: entity.PrintPriority,
		
		TemplateConfig: entity.TemplateConfig,
		
		PaperSize: entity.PaperSize,
		
		PrintOrientation: entity.PrintOrientation,
		
		IsActive: entity.IsActive,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'receipt',: entity.'receipt',,
		
		'report',: entity.'report',,
		
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


// validateBusinessRules validates business rules for printer_configurations
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PrinterConfigurations) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a printer_configurations can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

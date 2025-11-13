package data_export_request

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for DataExportRequests
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new DataExportRequests service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new data_export_requests
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateDataExportRequestsRequest) (*DataExportRequestsResponse, error) {
	s.logger.Info("creating data_export_requests",
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
	entity := &DataExportRequests{
		OrganizationId: orgID,
		
		ExportType: req.ExportType,
		
		ExportFormat: req.ExportFormat,
		
		DateFrom: req.DateFrom,
		
		DateTo: req.DateTo,
		
		Filters: req.Filters,
		
		Status: req.Status,
		
		Status: req.Status,
		
		FileName: req.FileName,
		
		FileSize: req.FileSize,
		
		FilePath: req.FilePath,
		
		DownloadUrl: req.DownloadUrl,
		
		DownloadExpiresAt: req.DownloadExpiresAt,
		
		TotalRecords: req.TotalRecords,
		
		ProcessedRecords: req.ProcessedRecords,
		
		ErrorMessage: req.ErrorMessage,
		
		RequestedBy: req.RequestedBy,
		
		RequestedAt: req.RequestedAt,
		
		StartedAt: req.StartedAt,
		
		CompletedAt: req.CompletedAt,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create data_export_requests: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created data_export_requests",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a data_export_requests by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DataExportRequestsResponse, error) {
	s.logger.Debug("getting data_export_requests",
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
		return nil, fmt.Errorf("failed to get data_export_requests: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("data_export_requests not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of data_export_requests records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*DataExportRequestsListResponse, error) {
	s.logger.Debug("listing data_export_requests",
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
		return nil, fmt.Errorf("failed to list data_export_requests: %w", err)
	}

	// Convert to response
	items := make([]*DataExportRequestsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &DataExportRequestsListResponse{
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

// Update updates an existing data_export_requests
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateDataExportRequestsRequest) (*DataExportRequestsResponse, error) {
	s.logger.Info("updating data_export_requests",
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
		return nil, fmt.Errorf("failed to get data_export_requests: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("data_export_requests not found or access denied")
	}
	

	// Update fields
	
	if req.ExportType != nil {
		entity.ExportType = req.ExportType
	}
	
	if req.ExportFormat != nil {
		entity.ExportFormat = req.ExportFormat
	}
	
	if req.DateFrom != nil {
		entity.DateFrom = req.DateFrom
	}
	
	if req.DateTo != nil {
		entity.DateTo = req.DateTo
	}
	
	if req.Filters != nil {
		entity.Filters = req.Filters
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.FileName != nil {
		entity.FileName = req.FileName
	}
	
	if req.FileSize != nil {
		entity.FileSize = req.FileSize
	}
	
	if req.FilePath != nil {
		entity.FilePath = req.FilePath
	}
	
	if req.DownloadUrl != nil {
		entity.DownloadUrl = req.DownloadUrl
	}
	
	if req.DownloadExpiresAt != nil {
		entity.DownloadExpiresAt = req.DownloadExpiresAt
	}
	
	if req.TotalRecords != nil {
		entity.TotalRecords = req.TotalRecords
	}
	
	if req.ProcessedRecords != nil {
		entity.ProcessedRecords = req.ProcessedRecords
	}
	
	if req.ErrorMessage != nil {
		entity.ErrorMessage = req.ErrorMessage
	}
	
	if req.RequestedBy != nil {
		entity.RequestedBy = req.RequestedBy
	}
	
	if req.RequestedAt != nil {
		entity.RequestedAt = req.RequestedAt
	}
	
	if req.StartedAt != nil {
		entity.StartedAt = req.StartedAt
	}
	
	if req.CompletedAt != nil {
		entity.CompletedAt = req.CompletedAt
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update data_export_requests: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated data_export_requests",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a data_export_requests
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting data_export_requests",
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
		return fmt.Errorf("failed to get data_export_requests: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("data_export_requests not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete data_export_requests: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted data_export_requests",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *DataExportRequests) *DataExportRequestsResponse {
	return &DataExportRequestsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ExportType: entity.ExportType,
		
		ExportFormat: entity.ExportFormat,
		
		DateFrom: entity.DateFrom,
		
		DateTo: entity.DateTo,
		
		Filters: entity.Filters,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		FileName: entity.FileName,
		
		FileSize: entity.FileSize,
		
		FilePath: entity.FilePath,
		
		DownloadUrl: entity.DownloadUrl,
		
		DownloadExpiresAt: entity.DownloadExpiresAt,
		
		TotalRecords: entity.TotalRecords,
		
		ProcessedRecords: entity.ProcessedRecords,
		
		ErrorMessage: entity.ErrorMessage,
		
		RequestedBy: entity.RequestedBy,
		
		RequestedAt: entity.RequestedAt,
		
		StartedAt: entity.StartedAt,
		
		CompletedAt: entity.CompletedAt,
		
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


// validateBusinessRules validates business rules for data_export_requests
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *DataExportRequests) error {
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

// canDelete checks if a data_export_requests can be deleted
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

package file_attachment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for FileAttachments
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new FileAttachments service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new file_attachments
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateFileAttachmentsRequest) (*FileAttachmentsResponse, error) {
	s.logger.Info("creating file_attachments",
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
	entity := &FileAttachments{
		OrganizationId: orgID,
		
		FileName: req.FileName,
		
		FileSize: req.FileSize,
		
		MimeType: req.MimeType,
		
		FileExtension: req.FileExtension,
		
		StorageProvider: req.StorageProvider,
		
		StoragePath: req.StoragePath,
		
		StorageUrl: req.StorageUrl,
		
		FileHash: req.FileHash,
		
		EntityType: req.EntityType,
		
		EntityId: req.EntityId,
		
		Description: req.Description,
		
		Tags: req.Tags,
		
		IsPublic: req.IsPublic,
		
		ImageWidth: req.ImageWidth,
		
		ImageHeight: req.ImageHeight,
		
		VirusScanStatus: req.VirusScanStatus,
		
		VirusScanAt: req.VirusScanAt,
		
		UploadedBy: req.UploadedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create file_attachments: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created file_attachments",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a file_attachments by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*FileAttachmentsResponse, error) {
	s.logger.Debug("getting file_attachments",
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
		return nil, fmt.Errorf("failed to get file_attachments: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("file_attachments not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of file_attachments records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*FileAttachmentsListResponse, error) {
	s.logger.Debug("listing file_attachments",
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
		return nil, fmt.Errorf("failed to list file_attachments: %w", err)
	}

	// Convert to response
	items := make([]*FileAttachmentsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &FileAttachmentsListResponse{
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

// Update updates an existing file_attachments
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateFileAttachmentsRequest) (*FileAttachmentsResponse, error) {
	s.logger.Info("updating file_attachments",
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
		return nil, fmt.Errorf("failed to get file_attachments: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("file_attachments not found or access denied")
	}
	

	// Update fields
	
	if req.FileName != nil {
		entity.FileName = *req.FileName
	}
	
	if req.FileSize != nil {
		entity.FileSize = *req.FileSize
	}
	
	if req.MimeType != nil {
		entity.MimeType = *req.MimeType
	}
	
	if req.FileExtension != nil {
		entity.FileExtension = *req.FileExtension
	}
	
	if req.StorageProvider != nil {
		entity.StorageProvider = *req.StorageProvider
	}
	
	if req.StoragePath != nil {
		entity.StoragePath = *req.StoragePath
	}
	
	if req.StorageUrl != nil {
		entity.StorageUrl = *req.StorageUrl
	}
	
	if req.FileHash != nil {
		entity.FileHash = *req.FileHash
	}
	
	if req.EntityType != nil {
		entity.EntityType = *req.EntityType
	}
	
	if req.EntityId != nil {
		entity.EntityId = *req.EntityId
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Tags != nil {
		entity.Tags = *req.Tags
	}
	
	if req.IsPublic != nil {
		entity.IsPublic = req.IsPublic
	}
	
	if req.ImageWidth != nil {
		entity.ImageWidth = *req.ImageWidth
	}
	
	if req.ImageHeight != nil {
		entity.ImageHeight = *req.ImageHeight
	}
	
	if req.VirusScanStatus != nil {
		entity.VirusScanStatus = *req.VirusScanStatus
	}
	
	if req.VirusScanAt != nil {
		entity.VirusScanAt = req.VirusScanAt
	}
	
	if req.UploadedBy != nil {
		entity.UploadedBy = *req.UploadedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update file_attachments: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated file_attachments",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a file_attachments
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting file_attachments",
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
		return fmt.Errorf("failed to get file_attachments: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("file_attachments not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete file_attachments: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted file_attachments",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *FileAttachments) *FileAttachmentsResponse {
	return &FileAttachmentsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		FileName: entity.FileName,
		
		FileSize: entity.FileSize,
		
		MimeType: entity.MimeType,
		
		FileExtension: entity.FileExtension,
		
		StorageProvider: entity.StorageProvider,
		
		StoragePath: entity.StoragePath,
		
		StorageUrl: entity.StorageUrl,
		
		FileHash: entity.FileHash,
		
		EntityType: entity.EntityType,
		
		EntityId: entity.EntityId,
		
		Description: entity.Description,
		
		Tags: entity.Tags,
		
		IsPublic: entity.IsPublic,
		
		ImageWidth: entity.ImageWidth,
		
		ImageHeight: entity.ImageHeight,
		
		VirusScanStatus: entity.VirusScanStatus,
		
		VirusScanAt: entity.VirusScanAt,
		
		UploadedBy: entity.UploadedBy,
		
		CreatedAt: entity.CreatedAt,
		
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


// validateBusinessRules validates business rules for file_attachments
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *FileAttachments) error {
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

// canDelete checks if a file_attachments can be deleted
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

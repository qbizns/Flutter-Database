package document_sequence

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

// Service handles business logic for DocumentSequences
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new DocumentSequences service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new document_sequences
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateDocumentSequencesRequest) (*DocumentSequencesResponse, error) {
	s.logger.Info("creating document_sequences",
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
	entity := &DocumentSequences{
		OrganizationID: orgID,
		
		DocumentType: req.DocumentType,
		
		Prefix: req.Prefix,
		
		Suffix: req.Suffix,
		
		NextNumber: req.NextNumber,
		
		Padding: req.Padding,
		
		IncrementBy: req.IncrementBy,
		
		ResetFrequency: req.ResetFrequency,
		
		'never',: req.'never',,
		
		'daily',: req.'daily',,
		
		'monthly',: req.'monthly',,
		
		'yearly',: req.'yearly',,
		
		'manual': req.'manual',
		
		LastResetAt: req.LastResetAt,
		
		LastResetValue: req.LastResetValue,
		
		IncludeDate: req.IncludeDate,
		
		DateFormat: req.DateFormat,
		
		LocationId: req.LocationId,
		
		IsActive: req.IsActive,
		
		AllowManualOverride: req.AllowManualOverride,
		
		ExampleNumber: req.ExampleNumber,
		
		Description: req.Description,
		
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
		return nil, fmt.Errorf("failed to create document_sequences: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created document_sequences",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a document_sequences by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DocumentSequencesResponse, error) {
	s.logger.Debug("getting document_sequences",
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
		return nil, fmt.Errorf("failed to get document_sequences: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("document_sequences not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of document_sequences records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*DocumentSequencesListResponse, error) {
	s.logger.Debug("listing document_sequences",
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
		return nil, fmt.Errorf("failed to list document_sequences: %w", err)
	}

	// Convert to response
	items := make([]*DocumentSequencesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &DocumentSequencesListResponse{
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

// Update updates an existing document_sequences
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateDocumentSequencesRequest) (*DocumentSequencesResponse, error) {
	s.logger.Info("updating document_sequences",
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
		return nil, fmt.Errorf("failed to get document_sequences: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("document_sequences not found or access denied")
	}
	

	// Update fields
	
	if req.DocumentType != nil {
		entity.DocumentType = *req.DocumentType
	}
	
	if req.Prefix != nil {
		entity.Prefix = *req.Prefix
	}
	
	if req.Suffix != nil {
		entity.Suffix = *req.Suffix
	}
	
	if req.NextNumber != nil {
		entity.NextNumber = *req.NextNumber
	}
	
	if req.Padding != nil {
		entity.Padding = *req.Padding
	}
	
	if req.IncrementBy != nil {
		entity.IncrementBy = *req.IncrementBy
	}
	
	if req.ResetFrequency != nil {
		entity.ResetFrequency = *req.ResetFrequency
	}
	
	if req.'never', != nil {
		entity.'never', = *req.'never',
	}
	
	if req.'daily', != nil {
		entity.'daily', = *req.'daily',
	}
	
	if req.'monthly', != nil {
		entity.'monthly', = *req.'monthly',
	}
	
	if req.'yearly', != nil {
		entity.'yearly', = *req.'yearly',
	}
	
	if req.'manual' != nil {
		entity.'manual' = *req.'manual'
	}
	
	if req.LastResetAt != nil {
		entity.LastResetAt = *req.LastResetAt
	}
	
	if req.LastResetValue != nil {
		entity.LastResetValue = *req.LastResetValue
	}
	
	if req.IncludeDate != nil {
		entity.IncludeDate = *req.IncludeDate
	}
	
	if req.DateFormat != nil {
		entity.DateFormat = *req.DateFormat
	}
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.AllowManualOverride != nil {
		entity.AllowManualOverride = *req.AllowManualOverride
	}
	
	if req.ExampleNumber != nil {
		entity.ExampleNumber = *req.ExampleNumber
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
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
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update document_sequences: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated document_sequences",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a document_sequences
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting document_sequences",
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
		return fmt.Errorf("failed to get document_sequences: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("document_sequences not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete document_sequences: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted document_sequences",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *DocumentSequences) *DocumentSequencesResponse {
	return &DocumentSequencesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		DocumentType: entity.DocumentType,
		
		Prefix: entity.Prefix,
		
		Suffix: entity.Suffix,
		
		NextNumber: entity.NextNumber,
		
		Padding: entity.Padding,
		
		IncrementBy: entity.IncrementBy,
		
		ResetFrequency: entity.ResetFrequency,
		
		'never',: entity.'never',,
		
		'daily',: entity.'daily',,
		
		'monthly',: entity.'monthly',,
		
		'yearly',: entity.'yearly',,
		
		'manual': entity.'manual',
		
		LastResetAt: entity.LastResetAt,
		
		LastResetValue: entity.LastResetValue,
		
		IncludeDate: entity.IncludeDate,
		
		DateFormat: entity.DateFormat,
		
		LocationId: entity.LocationId,
		
		IsActive: entity.IsActive,
		
		AllowManualOverride: entity.AllowManualOverride,
		
		ExampleNumber: entity.ExampleNumber,
		
		Description: entity.Description,
		
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


// validateBusinessRules validates business rules for document_sequences
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *DocumentSequences) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a document_sequences can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

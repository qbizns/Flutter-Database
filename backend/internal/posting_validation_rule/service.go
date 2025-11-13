package posting_validation_rule

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for PostingValidationRules
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PostingValidationRules service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new posting_validation_rules
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreatePostingValidationRulesRequest) (*PostingValidationRulesResponse, error) {
	s.logger.Info("creating posting_validation_rules",
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
	entity := &PostingValidationRules{
		OrganizationID: orgID,
		
		DocumentTypeCode: req.DocumentTypeCode,
		
		Event: req.Event,
		
		Target: req.Target,
		
		Code: req.Code,
		
		Name: req.Name,
		
		Description: req.Description,
		
		Expression: req.Expression,
		
		Severity: req.Severity,
		
		IsBlocking: req.IsBlocking,
		
		IsActive: req.IsActive,
		
		MessageTemplate: req.MessageTemplate,
		
		Priority: req.Priority,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		COALESCE(organizationId,: req.COALESCE(organizationId,,
		
		COALESCE(documentTypeCode,: req.COALESCE(documentTypeCode,,
		
		COALESCE(event,: req.COALESCE(event,,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create posting_validation_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created posting_validation_rules",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a posting_validation_rules by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PostingValidationRulesResponse, error) {
	s.logger.Debug("getting posting_validation_rules",
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
		return nil, fmt.Errorf("failed to get posting_validation_rules: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("posting_validation_rules not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of posting_validation_rules records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*PostingValidationRulesListResponse, error) {
	s.logger.Debug("listing posting_validation_rules",
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
		return nil, fmt.Errorf("failed to list posting_validation_rules: %w", err)
	}

	// Convert to response
	items := make([]*PostingValidationRulesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PostingValidationRulesListResponse{
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

// Update updates an existing posting_validation_rules
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdatePostingValidationRulesRequest) (*PostingValidationRulesResponse, error) {
	s.logger.Info("updating posting_validation_rules",
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
		return nil, fmt.Errorf("failed to get posting_validation_rules: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("posting_validation_rules not found or access denied")
	}
	

	// Update fields
	
	if req.DocumentTypeCode != nil {
		entity.DocumentTypeCode = req.DocumentTypeCode
	}
	
	if req.Event != nil {
		entity.Event = req.Event
	}
	
	if req.Target != nil {
		entity.Target = req.Target
	}
	
	if req.Code != nil {
		entity.Code = req.Code
	}
	
	if req.Name != nil {
		entity.Name = req.Name
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Expression != nil {
		entity.Expression = req.Expression
	}
	
	if req.Severity != nil {
		entity.Severity = req.Severity
	}
	
	if req.IsBlocking != nil {
		entity.IsBlocking = req.IsBlocking
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.MessageTemplate != nil {
		entity.MessageTemplate = req.MessageTemplate
	}
	
	if req.Priority != nil {
		entity.Priority = req.Priority
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
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
	
	if req.COALESCE(organizationId, != nil {
		entity.COALESCE(organizationId, = *req.COALESCE(organizationId,
	}
	
	if req.COALESCE(documentTypeCode, != nil {
		entity.COALESCE(documentTypeCode, = *req.COALESCE(documentTypeCode,
	}
	
	if req.COALESCE(event, != nil {
		entity.COALESCE(event, = *req.COALESCE(event,
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update posting_validation_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated posting_validation_rules",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a posting_validation_rules
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting posting_validation_rules",
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
		return fmt.Errorf("failed to get posting_validation_rules: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("posting_validation_rules not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete posting_validation_rules: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted posting_validation_rules",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PostingValidationRules) *PostingValidationRulesResponse {
	return &PostingValidationRulesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		DocumentTypeCode: entity.DocumentTypeCode,
		
		Event: entity.Event,
		
		Target: entity.Target,
		
		Code: entity.Code,
		
		Name: entity.Name,
		
		Description: entity.Description,
		
		Expression: entity.Expression,
		
		Severity: entity.Severity,
		
		IsBlocking: entity.IsBlocking,
		
		IsActive: entity.IsActive,
		
		MessageTemplate: entity.MessageTemplate,
		
		Priority: entity.Priority,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		COALESCE(organizationId,: entity.COALESCE(organizationId,,
		
		COALESCE(documentTypeCode,: entity.COALESCE(documentTypeCode,,
		
		COALESCE(event,: entity.COALESCE(event,,
		
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


// validateBusinessRules validates business rules for posting_validation_rules
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PostingValidationRules) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a posting_validation_rules can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

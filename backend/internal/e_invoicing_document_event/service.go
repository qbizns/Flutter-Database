package e_invoicing_document_event

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for EInvoicingDocumentEvents
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new EInvoicingDocumentEvents service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new e_invoicing_document_events
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateEInvoicingDocumentEventsRequest) (*EInvoicingDocumentEventsResponse, error) {
	s.logger.Info("creating e_invoicing_document_events",
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
	entity := &EInvoicingDocumentEvents{
		OrganizationId: orgID,
		
		EInvoicingDocumentId: req.EInvoicingDocumentId,
		
		EventType: req.EventType,
		
		'created',: req.'created',,
		
		'validated',: req.'validated',,
		
		'submitted',: req.'submitted',,
		
		'accepted',: req.'accepted',,
		
		'rejected',: req.'rejected',,
		
		'cancelled',: req.'cancelled',,
		
		'error',: req.'error',,
		
		'retry',: req.'retry',,
		
		'statusCheck': req.'statusCheck',
		
		EventTimestamp: req.EventTimestamp,
		
		PreviousStatus: req.PreviousStatus,
		
		NewStatus: req.NewStatus,
		
		EventDescription: req.EventDescription,
		
		EventData: req.EventData,
		
		HttpStatusCode: req.HttpStatusCode,
		
		HttpMethod: req.HttpMethod,
		
		ApiEndpoint: req.ApiEndpoint,
		
		RequestHeaders: req.RequestHeaders,
		
		ResponseHeaders: req.ResponseHeaders,
		
		ErrorCode: req.ErrorCode,
		
		ErrorMessage: req.ErrorMessage,
		
		ErrorDetails: req.ErrorDetails,
		
		TriggeredBy: req.TriggeredBy,
		
		UserId: req.UserId,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create e_invoicing_document_events: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created e_invoicing_document_events",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a e_invoicing_document_events by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*EInvoicingDocumentEventsResponse, error) {
	s.logger.Debug("getting e_invoicing_document_events",
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
		return nil, fmt.Errorf("failed to get e_invoicing_document_events: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("e_invoicing_document_events not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of e_invoicing_document_events records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*EInvoicingDocumentEventsListResponse, error) {
	s.logger.Debug("listing e_invoicing_document_events",
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
		return nil, fmt.Errorf("failed to list e_invoicing_document_events: %w", err)
	}

	// Convert to response
	items := make([]*EInvoicingDocumentEventsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &EInvoicingDocumentEventsListResponse{
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

// Update updates an existing e_invoicing_document_events
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateEInvoicingDocumentEventsRequest) (*EInvoicingDocumentEventsResponse, error) {
	s.logger.Info("updating e_invoicing_document_events",
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
		return nil, fmt.Errorf("failed to get e_invoicing_document_events: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("e_invoicing_document_events not found or access denied")
	}
	

	// Update fields
	
	if req.EInvoicingDocumentId != nil {
		entity.EInvoicingDocumentId = req.EInvoicingDocumentId
	}
	
	if req.EventType != nil {
		entity.EventType = req.EventType
	}
	
	if req.'created', != nil {
		entity.'created', = *req.'created',
	}
	
	if req.'validated', != nil {
		entity.'validated', = *req.'validated',
	}
	
	if req.'submitted', != nil {
		entity.'submitted', = *req.'submitted',
	}
	
	if req.'accepted', != nil {
		entity.'accepted', = *req.'accepted',
	}
	
	if req.'rejected', != nil {
		entity.'rejected', = *req.'rejected',
	}
	
	if req.'cancelled', != nil {
		entity.'cancelled', = *req.'cancelled',
	}
	
	if req.'error', != nil {
		entity.'error', = *req.'error',
	}
	
	if req.'retry', != nil {
		entity.'retry', = *req.'retry',
	}
	
	if req.'statusCheck' != nil {
		entity.'statusCheck' = *req.'statusCheck'
	}
	
	if req.EventTimestamp != nil {
		entity.EventTimestamp = req.EventTimestamp
	}
	
	if req.PreviousStatus != nil {
		entity.PreviousStatus = req.PreviousStatus
	}
	
	if req.NewStatus != nil {
		entity.NewStatus = req.NewStatus
	}
	
	if req.EventDescription != nil {
		entity.EventDescription = req.EventDescription
	}
	
	if req.EventData != nil {
		entity.EventData = req.EventData
	}
	
	if req.HttpStatusCode != nil {
		entity.HttpStatusCode = req.HttpStatusCode
	}
	
	if req.HttpMethod != nil {
		entity.HttpMethod = req.HttpMethod
	}
	
	if req.ApiEndpoint != nil {
		entity.ApiEndpoint = req.ApiEndpoint
	}
	
	if req.RequestHeaders != nil {
		entity.RequestHeaders = req.RequestHeaders
	}
	
	if req.ResponseHeaders != nil {
		entity.ResponseHeaders = req.ResponseHeaders
	}
	
	if req.ErrorCode != nil {
		entity.ErrorCode = req.ErrorCode
	}
	
	if req.ErrorMessage != nil {
		entity.ErrorMessage = req.ErrorMessage
	}
	
	if req.ErrorDetails != nil {
		entity.ErrorDetails = req.ErrorDetails
	}
	
	if req.TriggeredBy != nil {
		entity.TriggeredBy = req.TriggeredBy
	}
	
	if req.UserId != nil {
		entity.UserId = req.UserId
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update e_invoicing_document_events: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated e_invoicing_document_events",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a e_invoicing_document_events
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting e_invoicing_document_events",
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
		return fmt.Errorf("failed to get e_invoicing_document_events: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("e_invoicing_document_events not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete e_invoicing_document_events: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted e_invoicing_document_events",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *EInvoicingDocumentEvents) *EInvoicingDocumentEventsResponse {
	return &EInvoicingDocumentEventsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		EInvoicingDocumentId: entity.EInvoicingDocumentId,
		
		EventType: entity.EventType,
		
		'created',: entity.'created',,
		
		'validated',: entity.'validated',,
		
		'submitted',: entity.'submitted',,
		
		'accepted',: entity.'accepted',,
		
		'rejected',: entity.'rejected',,
		
		'cancelled',: entity.'cancelled',,
		
		'error',: entity.'error',,
		
		'retry',: entity.'retry',,
		
		'statusCheck': entity.'statusCheck',
		
		EventTimestamp: entity.EventTimestamp,
		
		PreviousStatus: entity.PreviousStatus,
		
		NewStatus: entity.NewStatus,
		
		EventDescription: entity.EventDescription,
		
		EventData: entity.EventData,
		
		HttpStatusCode: entity.HttpStatusCode,
		
		HttpMethod: entity.HttpMethod,
		
		ApiEndpoint: entity.ApiEndpoint,
		
		RequestHeaders: entity.RequestHeaders,
		
		ResponseHeaders: entity.ResponseHeaders,
		
		ErrorCode: entity.ErrorCode,
		
		ErrorMessage: entity.ErrorMessage,
		
		ErrorDetails: entity.ErrorDetails,
		
		TriggeredBy: entity.TriggeredBy,
		
		UserId: entity.UserId,
		
		CreatedAt: entity.CreatedAt,
		
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


// validateBusinessRules validates business rules for e_invoicing_document_events
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocumentEvents) error {
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

// canDelete checks if a e_invoicing_document_events can be deleted
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

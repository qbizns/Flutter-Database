package api_request_log

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

// Service handles business logic for ApiRequestLogs
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new ApiRequestLogs service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new api_request_logs
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateApiRequestLogsRequest) (*ApiRequestLogsResponse, error) {
	s.logger.Info("creating api_request_logs",
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
	entity := &ApiRequestLogs{
		OrganizationID: orgID,
		
		RequestId: req.RequestId,
		
		Method: req.Method,
		
		Path: req.Path,
		
		QueryParams: req.QueryParams,
		
		UserId: req.UserId,
		
		ApiKeyId: req.ApiKeyId,
		
		RequestHeaders: req.RequestHeaders,
		
		RequestBody: req.RequestBody,
		
		IpAddress: req.IpAddress,
		
		UserAgent: req.UserAgent,
		
		StatusCode: req.StatusCode,
		
		ResponseHeaders: req.ResponseHeaders,
		
		ResponseBody: req.ResponseBody,
		
		DurationMs: req.DurationMs,
		
		ErrorMessage: req.ErrorMessage,
		
		ErrorStack: req.ErrorStack,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create api_request_logs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created api_request_logs",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a api_request_logs by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ApiRequestLogsResponse, error) {
	s.logger.Debug("getting api_request_logs",
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
		return nil, fmt.Errorf("failed to get api_request_logs: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("api_request_logs not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of api_request_logs records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*ApiRequestLogsListResponse, error) {
	s.logger.Debug("listing api_request_logs",
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
		return nil, fmt.Errorf("failed to list api_request_logs: %w", err)
	}

	// Convert to response
	items := make([]*ApiRequestLogsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &ApiRequestLogsListResponse{
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

// Update updates an existing api_request_logs
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateApiRequestLogsRequest) (*ApiRequestLogsResponse, error) {
	s.logger.Info("updating api_request_logs",
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
		return nil, fmt.Errorf("failed to get api_request_logs: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("api_request_logs not found or access denied")
	}
	

	// Update fields
	
	if req.RequestId != nil {
		entity.RequestId = *req.RequestId
	}
	
	if req.Method != nil {
		entity.Method = *req.Method
	}
	
	if req.Path != nil {
		entity.Path = *req.Path
	}
	
	if req.QueryParams != nil {
		entity.QueryParams = *req.QueryParams
	}
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.ApiKeyId != nil {
		entity.ApiKeyId = *req.ApiKeyId
	}
	
	if req.RequestHeaders != nil {
		entity.RequestHeaders = *req.RequestHeaders
	}
	
	if req.RequestBody != nil {
		entity.RequestBody = *req.RequestBody
	}
	
	if req.IpAddress != nil {
		entity.IpAddress = *req.IpAddress
	}
	
	if req.UserAgent != nil {
		entity.UserAgent = *req.UserAgent
	}
	
	if req.StatusCode != nil {
		entity.StatusCode = *req.StatusCode
	}
	
	if req.ResponseHeaders != nil {
		entity.ResponseHeaders = *req.ResponseHeaders
	}
	
	if req.ResponseBody != nil {
		entity.ResponseBody = *req.ResponseBody
	}
	
	if req.DurationMs != nil {
		entity.DurationMs = *req.DurationMs
	}
	
	if req.ErrorMessage != nil {
		entity.ErrorMessage = *req.ErrorMessage
	}
	
	if req.ErrorStack != nil {
		entity.ErrorStack = *req.ErrorStack
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update api_request_logs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated api_request_logs",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a api_request_logs
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting api_request_logs",
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
		return fmt.Errorf("failed to get api_request_logs: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("api_request_logs not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete api_request_logs: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted api_request_logs",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *ApiRequestLogs) *ApiRequestLogsResponse {
	return &ApiRequestLogsResponse{
		
		Id: entity.Id,
		
		RequestId: entity.RequestId,
		
		Method: entity.Method,
		
		Path: entity.Path,
		
		QueryParams: entity.QueryParams,
		
		UserId: entity.UserId,
		
		OrganizationId: entity.OrganizationId,
		
		ApiKeyId: entity.ApiKeyId,
		
		RequestHeaders: entity.RequestHeaders,
		
		RequestBody: entity.RequestBody,
		
		IpAddress: entity.IpAddress,
		
		UserAgent: entity.UserAgent,
		
		StatusCode: entity.StatusCode,
		
		ResponseHeaders: entity.ResponseHeaders,
		
		ResponseBody: entity.ResponseBody,
		
		DurationMs: entity.DurationMs,
		
		ErrorMessage: entity.ErrorMessage,
		
		ErrorStack: entity.ErrorStack,
		
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


// validateBusinessRules validates business rules for api_request_logs
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *ApiRequestLogs) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a api_request_logs can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

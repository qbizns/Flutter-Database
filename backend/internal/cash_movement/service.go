package cash_movement

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

// Service handles business logic for CashMovements
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new CashMovements service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new cash_movements
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateCashMovementsRequest) (*CashMovementsResponse, error) {
	s.logger.Info("creating cash_movements",
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
	entity := &CashMovements{
		OrganizationID: orgID,
		
		PosSessionId: req.PosSessionId,
		
		CashDrawerId: req.CashDrawerId,
		
		MovementType: req.MovementType,
		
		Amount: req.Amount,
		
		ReasonCode: req.ReasonCode,
		
		ReasonDescription: req.ReasonDescription,
		
		UserId: req.UserId,
		
		RequiresApproval: req.RequiresApproval,
		
		ApprovedBy: req.ApprovedBy,
		
		ApprovedAt: req.ApprovedAt,
		
		Notes: req.Notes,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create cash_movements: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created cash_movements",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a cash_movements by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CashMovementsResponse, error) {
	s.logger.Debug("getting cash_movements",
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
		return nil, fmt.Errorf("failed to get cash_movements: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("cash_movements not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of cash_movements records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*CashMovementsListResponse, error) {
	s.logger.Debug("listing cash_movements",
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
		return nil, fmt.Errorf("failed to list cash_movements: %w", err)
	}

	// Convert to response
	items := make([]*CashMovementsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &CashMovementsListResponse{
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

// Update updates an existing cash_movements
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateCashMovementsRequest) (*CashMovementsResponse, error) {
	s.logger.Info("updating cash_movements",
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
		return nil, fmt.Errorf("failed to get cash_movements: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("cash_movements not found or access denied")
	}
	

	// Update fields
	
	if req.PosSessionId != nil {
		entity.PosSessionId = *req.PosSessionId
	}
	
	if req.CashDrawerId != nil {
		entity.CashDrawerId = *req.CashDrawerId
	}
	
	if req.MovementType != nil {
		entity.MovementType = *req.MovementType
	}
	
	if req.Amount != nil {
		entity.Amount = *req.Amount
	}
	
	if req.ReasonCode != nil {
		entity.ReasonCode = *req.ReasonCode
	}
	
	if req.ReasonDescription != nil {
		entity.ReasonDescription = *req.ReasonDescription
	}
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.RequiresApproval != nil {
		entity.RequiresApproval = *req.RequiresApproval
	}
	
	if req.ApprovedBy != nil {
		entity.ApprovedBy = *req.ApprovedBy
	}
	
	if req.ApprovedAt != nil {
		entity.ApprovedAt = *req.ApprovedAt
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update cash_movements: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated cash_movements",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a cash_movements
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting cash_movements",
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
		return fmt.Errorf("failed to get cash_movements: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("cash_movements not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete cash_movements: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted cash_movements",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *CashMovements) *CashMovementsResponse {
	return &CashMovementsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		PosSessionId: entity.PosSessionId,
		
		CashDrawerId: entity.CashDrawerId,
		
		MovementType: entity.MovementType,
		
		Amount: entity.Amount,
		
		ReasonCode: entity.ReasonCode,
		
		ReasonDescription: entity.ReasonDescription,
		
		UserId: entity.UserId,
		
		RequiresApproval: entity.RequiresApproval,
		
		ApprovedBy: entity.ApprovedBy,
		
		ApprovedAt: entity.ApprovedAt,
		
		Notes: entity.Notes,
		
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


// validateBusinessRules validates business rules for cash_movements
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *CashMovements) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a cash_movements can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

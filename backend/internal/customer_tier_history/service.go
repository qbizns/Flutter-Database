package customer_tier_history

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

// Service handles business logic for CustomerTierHistory
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new CustomerTierHistory service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new customer_tier_history
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateCustomerTierHistoryRequest) (*CustomerTierHistoryResponse, error) {
	s.logger.Info("creating customer_tier_history",
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
	entity := &CustomerTierHistory{
		OrganizationID: orgID,
		
		CustomerId: req.CustomerId,
		
		TierId: req.TierId,
		
		PreviousTierId: req.PreviousTierId,
		
		ChangeType: req.ChangeType,
		
		ChangeReason: req.ChangeReason,
		
		QualifyingPoints: req.QualifyingPoints,
		
		QualifyingSpend: req.QualifyingSpend,
		
		QualifyingPurchases: req.QualifyingPurchases,
		
		EffectiveDate: req.EffectiveDate,
		
		ValidUntil: req.ValidUntil,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create customer_tier_history: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created customer_tier_history",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a customer_tier_history by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CustomerTierHistoryResponse, error) {
	s.logger.Debug("getting customer_tier_history",
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
		return nil, fmt.Errorf("failed to get customer_tier_history: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("customer_tier_history not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of customer_tier_history records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*CustomerTierHistoryListResponse, error) {
	s.logger.Debug("listing customer_tier_history",
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
		return nil, fmt.Errorf("failed to list customer_tier_history: %w", err)
	}

	// Convert to response
	items := make([]*CustomerTierHistoryResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &CustomerTierHistoryListResponse{
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

// Update updates an existing customer_tier_history
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateCustomerTierHistoryRequest) (*CustomerTierHistoryResponse, error) {
	s.logger.Info("updating customer_tier_history",
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
		return nil, fmt.Errorf("failed to get customer_tier_history: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("customer_tier_history not found or access denied")
	}
	

	// Update fields
	
	if req.CustomerId != nil {
		entity.CustomerId = *req.CustomerId
	}
	
	if req.TierId != nil {
		entity.TierId = *req.TierId
	}
	
	if req.PreviousTierId != nil {
		entity.PreviousTierId = *req.PreviousTierId
	}
	
	if req.ChangeType != nil {
		entity.ChangeType = *req.ChangeType
	}
	
	if req.ChangeReason != nil {
		entity.ChangeReason = *req.ChangeReason
	}
	
	if req.QualifyingPoints != nil {
		entity.QualifyingPoints = *req.QualifyingPoints
	}
	
	if req.QualifyingSpend != nil {
		entity.QualifyingSpend = *req.QualifyingSpend
	}
	
	if req.QualifyingPurchases != nil {
		entity.QualifyingPurchases = *req.QualifyingPurchases
	}
	
	if req.EffectiveDate != nil {
		entity.EffectiveDate = *req.EffectiveDate
	}
	
	if req.ValidUntil != nil {
		entity.ValidUntil = *req.ValidUntil
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
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update customer_tier_history: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated customer_tier_history",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a customer_tier_history
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting customer_tier_history",
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
		return fmt.Errorf("failed to get customer_tier_history: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("customer_tier_history not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete customer_tier_history: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted customer_tier_history",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *CustomerTierHistory) *CustomerTierHistoryResponse {
	return &CustomerTierHistoryResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		CustomerId: entity.CustomerId,
		
		TierId: entity.TierId,
		
		PreviousTierId: entity.PreviousTierId,
		
		ChangeType: entity.ChangeType,
		
		ChangeReason: entity.ChangeReason,
		
		QualifyingPoints: entity.QualifyingPoints,
		
		QualifyingSpend: entity.QualifyingSpend,
		
		QualifyingPurchases: entity.QualifyingPurchases,
		
		EffectiveDate: entity.EffectiveDate,
		
		ValidUntil: entity.ValidUntil,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		CreatedBy: entity.CreatedBy,
		
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


// validateBusinessRules validates business rules for customer_tier_history
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *CustomerTierHistory) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a customer_tier_history can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

package pos_account_mapping

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/pos_account_mapping"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/pos_account_mapping"
	"go.uber.org/zap"
)

// Service handles business logic for PosAccountMappings
type Service struct {
	repo   *pos_account_mapping.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PosAccountMappings service
func NewService(repo *pos_account_mapping.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new pos_account_mappings
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreatePosAccountMappingsRequest) (*dto.PosAccountMappingsResponse, error) {
	s.logger.Info("creating pos_account_mappings",
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
	entity := &pos_account_mapping.PosAccountMappings{
		OrganizationID: orgID,
		
		SourceType: req.SourceType,
		
		'product',: req.'product',,
		
		'category',: req.'category',,
		
		'paymentMethod',: req.'paymentMethod',,
		
		'salesChannel',: req.'salesChannel',,
		
		'discount',: req.'discount',,
		
		'rounding',: req.'rounding',,
		
		'tax',: req.'tax',,
		
		'serviceCharge',: req.'serviceCharge',,
		
		'shipping',: req.'shipping',,
		
		'giftCard',: req.'giftCard',,
		
		'storeCredit',: req.'storeCredit',,
		
		'loyaltyRedemption',--: req.'loyaltyRedemption',--,
		
		'default': req.'default',
		
		SourceId: req.SourceId,
		
		SourceCode: req.SourceCode,
		
		Purpose: req.Purpose,
		
		'revenue',: req.'revenue',,
		
		'cogs',: req.'cogs',,
		
		'inventory',: req.'inventory',,
		
		'expense',: req.'expense',,
		
		'liability',: req.'liability',,
		
		'asset',: req.'asset',,
		
		'discountExpense',: req.'discountExpense',,
		
		'discountContra',: req.'discountContra',,
		
		'taxLiability',: req.'taxLiability',,
		
		'rounding',: req.'rounding',,
		
		'clearing': req.'clearing',
		
		AccountId: req.AccountId,
		
		IsDefault: req.IsDefault,
		
		IsActive: req.IsActive,
		
		Priority: req.Priority,
		
		Conditions: req.Conditions,
		
		EffectiveFrom: req.EffectiveFrom,
		
		EffectiveTo: req.EffectiveTo,
		
		Description: req.Description,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		(sourceType: req.(sourceType,
		
		(sourceType: req.(sourceType,
		
		EffectiveFrom: req.EffectiveFrom,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create pos_account_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created pos_account_mappings",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a pos_account_mappings by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.PosAccountMappingsResponse, error) {
	s.logger.Debug("getting pos_account_mappings",
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
		return nil, fmt.Errorf("failed to get pos_account_mappings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("pos_account_mappings not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of pos_account_mappings records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.PosAccountMappingsListResponse, error) {
	s.logger.Debug("listing pos_account_mappings",
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
		return nil, fmt.Errorf("failed to list pos_account_mappings: %w", err)
	}

	// Convert to response
	items := make([]*dto.PosAccountMappingsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.PosAccountMappingsListResponse{
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

// Update updates an existing pos_account_mappings
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdatePosAccountMappingsRequest) (*dto.PosAccountMappingsResponse, error) {
	s.logger.Info("updating pos_account_mappings",
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
		return nil, fmt.Errorf("failed to get pos_account_mappings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("pos_account_mappings not found or access denied")
	}
	

	// Update fields
	
	if req.SourceType != nil {
		entity.SourceType = *req.SourceType
	}
	
	if req.'product', != nil {
		entity.'product', = *req.'product',
	}
	
	if req.'category', != nil {
		entity.'category', = *req.'category',
	}
	
	if req.'paymentMethod', != nil {
		entity.'paymentMethod', = *req.'paymentMethod',
	}
	
	if req.'salesChannel', != nil {
		entity.'salesChannel', = *req.'salesChannel',
	}
	
	if req.'discount', != nil {
		entity.'discount', = *req.'discount',
	}
	
	if req.'rounding', != nil {
		entity.'rounding', = *req.'rounding',
	}
	
	if req.'tax', != nil {
		entity.'tax', = *req.'tax',
	}
	
	if req.'serviceCharge', != nil {
		entity.'serviceCharge', = *req.'serviceCharge',
	}
	
	if req.'shipping', != nil {
		entity.'shipping', = *req.'shipping',
	}
	
	if req.'giftCard', != nil {
		entity.'giftCard', = *req.'giftCard',
	}
	
	if req.'storeCredit', != nil {
		entity.'storeCredit', = *req.'storeCredit',
	}
	
	if req.'loyaltyRedemption',-- != nil {
		entity.'loyaltyRedemption',-- = *req.'loyaltyRedemption',--
	}
	
	if req.'default' != nil {
		entity.'default' = *req.'default'
	}
	
	if req.SourceId != nil {
		entity.SourceId = *req.SourceId
	}
	
	if req.SourceCode != nil {
		entity.SourceCode = *req.SourceCode
	}
	
	if req.Purpose != nil {
		entity.Purpose = *req.Purpose
	}
	
	if req.'revenue', != nil {
		entity.'revenue', = *req.'revenue',
	}
	
	if req.'cogs', != nil {
		entity.'cogs', = *req.'cogs',
	}
	
	if req.'inventory', != nil {
		entity.'inventory', = *req.'inventory',
	}
	
	if req.'expense', != nil {
		entity.'expense', = *req.'expense',
	}
	
	if req.'liability', != nil {
		entity.'liability', = *req.'liability',
	}
	
	if req.'asset', != nil {
		entity.'asset', = *req.'asset',
	}
	
	if req.'discountExpense', != nil {
		entity.'discountExpense', = *req.'discountExpense',
	}
	
	if req.'discountContra', != nil {
		entity.'discountContra', = *req.'discountContra',
	}
	
	if req.'taxLiability', != nil {
		entity.'taxLiability', = *req.'taxLiability',
	}
	
	if req.'rounding', != nil {
		entity.'rounding', = *req.'rounding',
	}
	
	if req.'clearing' != nil {
		entity.'clearing' = *req.'clearing'
	}
	
	if req.AccountId != nil {
		entity.AccountId = *req.AccountId
	}
	
	if req.IsDefault != nil {
		entity.IsDefault = *req.IsDefault
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.Priority != nil {
		entity.Priority = *req.Priority
	}
	
	if req.Conditions != nil {
		entity.Conditions = *req.Conditions
	}
	
	if req.EffectiveFrom != nil {
		entity.EffectiveFrom = *req.EffectiveFrom
	}
	
	if req.EffectiveTo != nil {
		entity.EffectiveTo = *req.EffectiveTo
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
	
	if req.(sourceType != nil {
		entity.(sourceType = *req.(sourceType
	}
	
	if req.(sourceType != nil {
		entity.(sourceType = *req.(sourceType
	}
	
	if req.EffectiveFrom != nil {
		entity.EffectiveFrom = *req.EffectiveFrom
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update pos_account_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated pos_account_mappings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a pos_account_mappings
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting pos_account_mappings",
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
		return fmt.Errorf("failed to get pos_account_mappings: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("pos_account_mappings not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete pos_account_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted pos_account_mappings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *pos_account_mapping.PosAccountMappings) *dto.PosAccountMappingsResponse {
	return &dto.PosAccountMappingsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		SourceType: entity.SourceType,
		
		'product',: entity.'product',,
		
		'category',: entity.'category',,
		
		'paymentMethod',: entity.'paymentMethod',,
		
		'salesChannel',: entity.'salesChannel',,
		
		'discount',: entity.'discount',,
		
		'rounding',: entity.'rounding',,
		
		'tax',: entity.'tax',,
		
		'serviceCharge',: entity.'serviceCharge',,
		
		'shipping',: entity.'shipping',,
		
		'giftCard',: entity.'giftCard',,
		
		'storeCredit',: entity.'storeCredit',,
		
		'loyaltyRedemption',--: entity.'loyaltyRedemption',--,
		
		'default': entity.'default',
		
		SourceId: entity.SourceId,
		
		SourceCode: entity.SourceCode,
		
		Purpose: entity.Purpose,
		
		'revenue',: entity.'revenue',,
		
		'cogs',: entity.'cogs',,
		
		'inventory',: entity.'inventory',,
		
		'expense',: entity.'expense',,
		
		'liability',: entity.'liability',,
		
		'asset',: entity.'asset',,
		
		'discountExpense',: entity.'discountExpense',,
		
		'discountContra',: entity.'discountContra',,
		
		'taxLiability',: entity.'taxLiability',,
		
		'rounding',: entity.'rounding',,
		
		'clearing': entity.'clearing',
		
		AccountId: entity.AccountId,
		
		IsDefault: entity.IsDefault,
		
		IsActive: entity.IsActive,
		
		Priority: entity.Priority,
		
		Conditions: entity.Conditions,
		
		EffectiveFrom: entity.EffectiveFrom,
		
		EffectiveTo: entity.EffectiveTo,
		
		Description: entity.Description,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		(sourceType: entity.(sourceType,
		
		(sourceType: entity.(sourceType,
		
		EffectiveFrom: entity.EffectiveFrom,
		
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


// validateBusinessRules validates business rules for pos_account_mappings
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *pos_account_mapping.PosAccountMappings) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a pos_account_mappings can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

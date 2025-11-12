package inventory_valuation_setting

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

// Service handles business logic for InventoryValuationSettings
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new InventoryValuationSettings service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new inventory_valuation_settings
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateInventoryValuationSettingsRequest) (*InventoryValuationSettingsResponse, error) {
	s.logger.Info("creating inventory_valuation_settings",
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
	entity := &InventoryValuationSettings{
		OrganizationID: orgID,
		
		ValuationMethod: req.ValuationMethod,
		
		'fifo',: req.'fifo',,
		
		'lifo',: req.'lifo',,
		
		'weightedAverage',: req.'weightedAverage',,
		
		'movingAverage',: req.'movingAverage',,
		
		'standardCost',: req.'standardCost',,
		
		'specificId': req.'specificId',
		
		CostLayerGranularity: req.CostLayerGranularity,
		
		'product',: req.'product',,
		
		'productLocation',: req.'productLocation',,
		
		'productLocationLot',: req.'productLocationLot',,
		
		'serialNumber': req.'serialNumber',
		
		DefaultInventoryAccountId: req.DefaultInventoryAccountId,
		
		DefaultCogsAccountId: req.DefaultCogsAccountId,
		
		DefaultInventoryAdjustmentAccountId: req.DefaultInventoryAdjustmentAccountId,
		
		DefaultInventoryVarianceAccountId: req.DefaultInventoryVarianceAccountId,
		
		CogsRecognitionTiming: req.CogsRecognitionTiming,
		
		'onSale',: req.'onSale',,
		
		'onDelivery',: req.'onDelivery',,
		
		'onPayment': req.'onPayment',
		
		AllowNegativeInventory: req.AllowNegativeInventory,
		
		RevalueOnPurchase: req.RevalueOnPurchase,
		
		RoundUnitCostToDecimals: req.RoundUnitCostToDecimals,
		
		RevaluationFrequency: req.RevaluationFrequency,
		
		'realTime',: req.'realTime',,
		
		'daily',: req.'daily',,
		
		'monthly',: req.'monthly',,
		
		'manual': req.'manual',
		
		IsActive: req.IsActive,
		
		EffectiveFrom: req.EffectiveFrom,
		
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
		return nil, fmt.Errorf("failed to create inventory_valuation_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created inventory_valuation_settings",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a inventory_valuation_settings by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryValuationSettingsResponse, error) {
	s.logger.Debug("getting inventory_valuation_settings",
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
		return nil, fmt.Errorf("failed to get inventory_valuation_settings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("inventory_valuation_settings not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of inventory_valuation_settings records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*InventoryValuationSettingsListResponse, error) {
	s.logger.Debug("listing inventory_valuation_settings",
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
		return nil, fmt.Errorf("failed to list inventory_valuation_settings: %w", err)
	}

	// Convert to response
	items := make([]*InventoryValuationSettingsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &InventoryValuationSettingsListResponse{
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

// Update updates an existing inventory_valuation_settings
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateInventoryValuationSettingsRequest) (*InventoryValuationSettingsResponse, error) {
	s.logger.Info("updating inventory_valuation_settings",
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
		return nil, fmt.Errorf("failed to get inventory_valuation_settings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("inventory_valuation_settings not found or access denied")
	}
	

	// Update fields
	
	if req.ValuationMethod != nil {
		entity.ValuationMethod = *req.ValuationMethod
	}
	
	if req.'fifo', != nil {
		entity.'fifo', = *req.'fifo',
	}
	
	if req.'lifo', != nil {
		entity.'lifo', = *req.'lifo',
	}
	
	if req.'weightedAverage', != nil {
		entity.'weightedAverage', = *req.'weightedAverage',
	}
	
	if req.'movingAverage', != nil {
		entity.'movingAverage', = *req.'movingAverage',
	}
	
	if req.'standardCost', != nil {
		entity.'standardCost', = *req.'standardCost',
	}
	
	if req.'specificId' != nil {
		entity.'specificId' = *req.'specificId'
	}
	
	if req.CostLayerGranularity != nil {
		entity.CostLayerGranularity = *req.CostLayerGranularity
	}
	
	if req.'product', != nil {
		entity.'product', = *req.'product',
	}
	
	if req.'productLocation', != nil {
		entity.'productLocation', = *req.'productLocation',
	}
	
	if req.'productLocationLot', != nil {
		entity.'productLocationLot', = *req.'productLocationLot',
	}
	
	if req.'serialNumber' != nil {
		entity.'serialNumber' = *req.'serialNumber'
	}
	
	if req.DefaultInventoryAccountId != nil {
		entity.DefaultInventoryAccountId = *req.DefaultInventoryAccountId
	}
	
	if req.DefaultCogsAccountId != nil {
		entity.DefaultCogsAccountId = *req.DefaultCogsAccountId
	}
	
	if req.DefaultInventoryAdjustmentAccountId != nil {
		entity.DefaultInventoryAdjustmentAccountId = *req.DefaultInventoryAdjustmentAccountId
	}
	
	if req.DefaultInventoryVarianceAccountId != nil {
		entity.DefaultInventoryVarianceAccountId = *req.DefaultInventoryVarianceAccountId
	}
	
	if req.CogsRecognitionTiming != nil {
		entity.CogsRecognitionTiming = *req.CogsRecognitionTiming
	}
	
	if req.'onSale', != nil {
		entity.'onSale', = *req.'onSale',
	}
	
	if req.'onDelivery', != nil {
		entity.'onDelivery', = *req.'onDelivery',
	}
	
	if req.'onPayment' != nil {
		entity.'onPayment' = *req.'onPayment'
	}
	
	if req.AllowNegativeInventory != nil {
		entity.AllowNegativeInventory = *req.AllowNegativeInventory
	}
	
	if req.RevalueOnPurchase != nil {
		entity.RevalueOnPurchase = *req.RevalueOnPurchase
	}
	
	if req.RoundUnitCostToDecimals != nil {
		entity.RoundUnitCostToDecimals = *req.RoundUnitCostToDecimals
	}
	
	if req.RevaluationFrequency != nil {
		entity.RevaluationFrequency = *req.RevaluationFrequency
	}
	
	if req.'realTime', != nil {
		entity.'realTime', = *req.'realTime',
	}
	
	if req.'daily', != nil {
		entity.'daily', = *req.'daily',
	}
	
	if req.'monthly', != nil {
		entity.'monthly', = *req.'monthly',
	}
	
	if req.'manual' != nil {
		entity.'manual' = *req.'manual'
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.EffectiveFrom != nil {
		entity.EffectiveFrom = *req.EffectiveFrom
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
		return nil, fmt.Errorf("failed to update inventory_valuation_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated inventory_valuation_settings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a inventory_valuation_settings
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting inventory_valuation_settings",
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
		return fmt.Errorf("failed to get inventory_valuation_settings: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("inventory_valuation_settings not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete inventory_valuation_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted inventory_valuation_settings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *InventoryValuationSettings) *InventoryValuationSettingsResponse {
	return &InventoryValuationSettingsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ValuationMethod: entity.ValuationMethod,
		
		'fifo',: entity.'fifo',,
		
		'lifo',: entity.'lifo',,
		
		'weightedAverage',: entity.'weightedAverage',,
		
		'movingAverage',: entity.'movingAverage',,
		
		'standardCost',: entity.'standardCost',,
		
		'specificId': entity.'specificId',
		
		CostLayerGranularity: entity.CostLayerGranularity,
		
		'product',: entity.'product',,
		
		'productLocation',: entity.'productLocation',,
		
		'productLocationLot',: entity.'productLocationLot',,
		
		'serialNumber': entity.'serialNumber',
		
		DefaultInventoryAccountId: entity.DefaultInventoryAccountId,
		
		DefaultCogsAccountId: entity.DefaultCogsAccountId,
		
		DefaultInventoryAdjustmentAccountId: entity.DefaultInventoryAdjustmentAccountId,
		
		DefaultInventoryVarianceAccountId: entity.DefaultInventoryVarianceAccountId,
		
		CogsRecognitionTiming: entity.CogsRecognitionTiming,
		
		'onSale',: entity.'onSale',,
		
		'onDelivery',: entity.'onDelivery',,
		
		'onPayment': entity.'onPayment',
		
		AllowNegativeInventory: entity.AllowNegativeInventory,
		
		RevalueOnPurchase: entity.RevalueOnPurchase,
		
		RoundUnitCostToDecimals: entity.RoundUnitCostToDecimals,
		
		RevaluationFrequency: entity.RevaluationFrequency,
		
		'realTime',: entity.'realTime',,
		
		'daily',: entity.'daily',,
		
		'monthly',: entity.'monthly',,
		
		'manual': entity.'manual',
		
		IsActive: entity.IsActive,
		
		EffectiveFrom: entity.EffectiveFrom,
		
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


// validateBusinessRules validates business rules for inventory_valuation_settings
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *InventoryValuationSettings) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a inventory_valuation_settings can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

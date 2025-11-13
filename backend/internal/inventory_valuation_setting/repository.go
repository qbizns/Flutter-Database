package inventory_valuation_setting

import (
	"encoding/json"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for InventoryValuationSettings
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new InventoryValuationSettings repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// InventoryValuationSettings represents a inventory_valuation_settings entity
type InventoryValuationSettings struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ValuationMethod string `json:"valuation_method" db:"valuation_method"`
	'fifo', *string `json:"'fifo'," db:"'fifo',"`
	'lifo', *string `json:"'lifo'," db:"'lifo',"`
	'weightedAverage', *string `json:"'weighted_average'," db:"'weighted_average',"`
	'movingAverage', *string `json:"'moving_average'," db:"'moving_average',"`
	'standardCost', *string `json:"'standard_cost'," db:"'standard_cost',"`
	'specificId' *string `json:"'specific_id'" db:"'specific_id'"`
	CostLayerGranularity *string `json:"cost_layer_granularity" db:"cost_layer_granularity"`
	'product', *string `json:"'product'," db:"'product',"`
	'productLocation', *string `json:"'product_location'," db:"'product_location',"`
	'productLocationLot', *string `json:"'product_location_lot'," db:"'product_location_lot',"`
	'serialNumber' *string `json:"'serial_number'" db:"'serial_number'"`
	DefaultInventoryAccountId *uuid.UUID `json:"default_inventory_account_id" db:"default_inventory_account_id"`
	DefaultCogsAccountId *uuid.UUID `json:"default_cogs_account_id" db:"default_cogs_account_id"`
	DefaultInventoryAdjustmentAccountId *uuid.UUID `json:"default_inventory_adjustment_account_id" db:"default_inventory_adjustment_account_id"`
	DefaultInventoryVarianceAccountId *uuid.UUID `json:"default_inventory_variance_account_id" db:"default_inventory_variance_account_id"`
	CogsRecognitionTiming *string `json:"cogs_recognition_timing" db:"cogs_recognition_timing"`
	'onSale', *string `json:"'on_sale'," db:"'on_sale',"`
	'onDelivery', *string `json:"'on_delivery'," db:"'on_delivery',"`
	'onPayment' *string `json:"'on_payment'" db:"'on_payment'"`
	AllowNegativeInventory *bool `json:"allow_negative_inventory" db:"allow_negative_inventory"`
	RevalueOnPurchase *bool `json:"revalue_on_purchase" db:"revalue_on_purchase"`
	RoundUnitCostToDecimals *int64 `json:"round_unit_cost_to_decimals" db:"round_unit_cost_to_decimals"`
	RevaluationFrequency *string `json:"revaluation_frequency" db:"revaluation_frequency"`
	'realTime', *string `json:"'real_time'," db:"'real_time',"`
	'daily', *string `json:"'daily'," db:"'daily',"`
	'monthly', *string `json:"'monthly'," db:"'monthly',"`
	'manual' *string `json:"'manual'" db:"'manual'"`
	IsActive *bool `json:"is_active" db:"is_active"`
	EffectiveFrom *time.Time `json:"effective_from" db:"effective_from"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new inventory_valuation_settings record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *InventoryValuationSettings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "inventory_valuation_settings", duration, nil)
	}()

	query := `
		INSERT INTO inventory_valuation_settings (
			, organization_id
			, valuation_method
			, 'fifo',
			, 'lifo',
			, 'weighted_average',
			, 'moving_average',
			, 'standard_cost',
			, 'specific_id'
			, cost_layer_granularity
			, 'product',
			, 'product_location',
			, 'product_location_lot',
			, 'serial_number'
			, default_inventory_account_id
			, default_cogs_account_id
			, default_inventory_adjustment_account_id
			, default_inventory_variance_account_id
			, cogs_recognition_timing
			, 'on_sale',
			, 'on_delivery',
			, 'on_payment'
			, allow_negative_inventory
			, revalue_on_purchase
			, round_unit_cost_to_decimals
			, revaluation_frequency
			, 'real_time',
			, 'daily',
			, 'monthly',
			, 'manual'
			, is_active
			, effective_from
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $10
			, $11
			, $12
			, $13
			, $14
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $27
			, $28
			, $29
			, $30
			, $31
			, $32
			, $33
			, $34
			, $37
			, $38
			, $39
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ValuationMethod,
		entity.'fifo',,
		entity.'lifo',,
		entity.'weightedAverage',,
		entity.'movingAverage',,
		entity.'standardCost',,
		entity.'specificId',
		entity.CostLayerGranularity,
		entity.'product',,
		entity.'productLocation',,
		entity.'productLocationLot',,
		entity.'serialNumber',
		entity.DefaultInventoryAccountId,
		entity.DefaultCogsAccountId,
		entity.DefaultInventoryAdjustmentAccountId,
		entity.DefaultInventoryVarianceAccountId,
		entity.CogsRecognitionTiming,
		entity.'onSale',,
		entity.'onDelivery',,
		entity.'onPayment',
		entity.AllowNegativeInventory,
		entity.RevalueOnPurchase,
		entity.RoundUnitCostToDecimals,
		entity.RevaluationFrequency,
		entity.'realTime',,
		entity.'daily',,
		entity.'monthly',,
		entity.'manual',
		entity.IsActive,
		entity.EffectiveFrom,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create inventory_valuation_settings", zap.Error(err))
		return fmt.Errorf("failed to create inventory_valuation_settings: %w", err)
	}

	r.logger.Info("created inventory_valuation_settings",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a inventory_valuation_settings by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*InventoryValuationSettings, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_valuation_settings", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, valuation_method
			, 'fifo',
			, 'lifo',
			, 'weighted_average',
			, 'moving_average',
			, 'standard_cost',
			, 'specific_id'
			, cost_layer_granularity
			, 'product',
			, 'product_location',
			, 'product_location_lot',
			, 'serial_number'
			, default_inventory_account_id
			, default_cogs_account_id
			, default_inventory_adjustment_account_id
			, default_inventory_variance_account_id
			, cogs_recognition_timing
			, 'on_sale',
			, 'on_delivery',
			, 'on_payment'
			, allow_negative_inventory
			, revalue_on_purchase
			, round_unit_cost_to_decimals
			, revaluation_frequency
			, 'real_time',
			, 'daily',
			, 'monthly',
			, 'manual'
			, is_active
			, effective_from
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM inventory_valuation_settings
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity InventoryValuationSettings
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ValuationMethod,
		&entity.'fifo',,
		&entity.'lifo',,
		&entity.'weightedAverage',,
		&entity.'movingAverage',,
		&entity.'standardCost',,
		&entity.'specificId',
		&entity.CostLayerGranularity,
		&entity.'product',,
		&entity.'productLocation',,
		&entity.'productLocationLot',,
		&entity.'serialNumber',
		&entity.DefaultInventoryAccountId,
		&entity.DefaultCogsAccountId,
		&entity.DefaultInventoryAdjustmentAccountId,
		&entity.DefaultInventoryVarianceAccountId,
		&entity.CogsRecognitionTiming,
		&entity.'onSale',,
		&entity.'onDelivery',,
		&entity.'onPayment',
		&entity.AllowNegativeInventory,
		&entity.RevalueOnPurchase,
		&entity.RoundUnitCostToDecimals,
		&entity.RevaluationFrequency,
		&entity.'realTime',,
		&entity.'daily',,
		&entity.'monthly',,
		&entity.'manual',
		&entity.IsActive,
		&entity.EffectiveFrom,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("inventory_valuation_settings not found")
	}

	if err != nil {
		r.logger.Error("failed to get inventory_valuation_settings", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get inventory_valuation_settings: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of inventory_valuation_settings records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*InventoryValuationSettings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_valuation_settings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_valuation_settings
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_valuation_settings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, valuation_method
			, 'fifo',
			, 'lifo',
			, 'weighted_average',
			, 'moving_average',
			, 'standard_cost',
			, 'specific_id'
			, cost_layer_granularity
			, 'product',
			, 'product_location',
			, 'product_location_lot',
			, 'serial_number'
			, default_inventory_account_id
			, default_cogs_account_id
			, default_inventory_adjustment_account_id
			, default_inventory_variance_account_id
			, cogs_recognition_timing
			, 'on_sale',
			, 'on_delivery',
			, 'on_payment'
			, allow_negative_inventory
			, revalue_on_purchase
			, round_unit_cost_to_decimals
			, revaluation_frequency
			, 'real_time',
			, 'daily',
			, 'monthly',
			, 'manual'
			, is_active
			, effective_from
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM inventory_valuation_settings
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_valuation_settings", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_valuation_settings: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryValuationSettings
	for rows.Next() {
		var entity InventoryValuationSettings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ValuationMethod,
			&entity.'fifo',,
			&entity.'lifo',,
			&entity.'weightedAverage',,
			&entity.'movingAverage',,
			&entity.'standardCost',,
			&entity.'specificId',
			&entity.CostLayerGranularity,
			&entity.'product',,
			&entity.'productLocation',,
			&entity.'productLocationLot',,
			&entity.'serialNumber',
			&entity.DefaultInventoryAccountId,
			&entity.DefaultCogsAccountId,
			&entity.DefaultInventoryAdjustmentAccountId,
			&entity.DefaultInventoryVarianceAccountId,
			&entity.CogsRecognitionTiming,
			&entity.'onSale',,
			&entity.'onDelivery',,
			&entity.'onPayment',
			&entity.AllowNegativeInventory,
			&entity.RevalueOnPurchase,
			&entity.RoundUnitCostToDecimals,
			&entity.RevaluationFrequency,
			&entity.'realTime',,
			&entity.'daily',,
			&entity.'monthly',,
			&entity.'manual',
			&entity.IsActive,
			&entity.EffectiveFrom,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_valuation_settings: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating inventory_valuation_settings rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing inventory_valuation_settings record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *InventoryValuationSettings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_valuation_settings", duration, nil)
	}()

	query := `
		UPDATE inventory_valuation_settings
		SET
			, organization_id = $2
			, valuation_method = $3
			, 'fifo', = $4
			, 'lifo', = $5
			, 'weighted_average', = $6
			, 'moving_average', = $7
			, 'standard_cost', = $8
			, 'specific_id' = $9
			, cost_layer_granularity = $10
			, 'product', = $11
			, 'product_location', = $12
			, 'product_location_lot', = $13
			, 'serial_number' = $14
			, default_inventory_account_id = $15
			, default_cogs_account_id = $16
			, default_inventory_adjustment_account_id = $17
			, default_inventory_variance_account_id = $18
			, cogs_recognition_timing = $19
			, 'on_sale', = $20
			, 'on_delivery', = $21
			, 'on_payment' = $22
			, allow_negative_inventory = $23
			, revalue_on_purchase = $24
			, round_unit_cost_to_decimals = $25
			, revaluation_frequency = $26
			, 'real_time', = $27
			, 'daily', = $28
			, 'monthly', = $29
			, 'manual' = $30
			, is_active = $31
			, effective_from = $32
			, notes = $33
			, metadata = $34
			, updated_at = $36
			, deleted_at = $37
			, created_by = $38
			, updated_by = $39
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $40
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ValuationMethod,
		entity.'fifo',,
		entity.'lifo',,
		entity.'weightedAverage',,
		entity.'movingAverage',,
		entity.'standardCost',,
		entity.'specificId',
		entity.CostLayerGranularity,
		entity.'product',,
		entity.'productLocation',,
		entity.'productLocationLot',,
		entity.'serialNumber',
		entity.DefaultInventoryAccountId,
		entity.DefaultCogsAccountId,
		entity.DefaultInventoryAdjustmentAccountId,
		entity.DefaultInventoryVarianceAccountId,
		entity.CogsRecognitionTiming,
		entity.'onSale',,
		entity.'onDelivery',,
		entity.'onPayment',
		entity.AllowNegativeInventory,
		entity.RevalueOnPurchase,
		entity.RoundUnitCostToDecimals,
		entity.RevaluationFrequency,
		entity.'realTime',,
		entity.'daily',,
		entity.'monthly',,
		entity.'manual',
		entity.IsActive,
		entity.EffectiveFrom,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update inventory_valuation_settings", zap.Error(err))
		return fmt.Errorf("failed to update inventory_valuation_settings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_valuation_settings not found or already deleted")
	}

	r.logger.Info("updated inventory_valuation_settings",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a inventory_valuation_settings record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_valuation_settings", duration, nil)
	}()

	query := `
		UPDATE inventory_valuation_settings
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete inventory_valuation_settings", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete inventory_valuation_settings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_valuation_settings not found or already deleted")
	}

	r.logger.Info("deleted inventory_valuation_settings", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves inventory_valuation_settings records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*InventoryValuationSettings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_valuation_settings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_valuation_settings
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_valuation_settings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, valuation_method
			, 'fifo',
			, 'lifo',
			, 'weighted_average',
			, 'moving_average',
			, 'standard_cost',
			, 'specific_id'
			, cost_layer_granularity
			, 'product',
			, 'product_location',
			, 'product_location_lot',
			, 'serial_number'
			, default_inventory_account_id
			, default_cogs_account_id
			, default_inventory_adjustment_account_id
			, default_inventory_variance_account_id
			, cogs_recognition_timing
			, 'on_sale',
			, 'on_delivery',
			, 'on_payment'
			, allow_negative_inventory
			, revalue_on_purchase
			, round_unit_cost_to_decimals
			, revaluation_frequency
			, 'real_time',
			, 'daily',
			, 'monthly',
			, 'manual'
			, is_active
			, effective_from
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM inventory_valuation_settings
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_valuation_settings by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_valuation_settings: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryValuationSettings
	for rows.Next() {
		var entity InventoryValuationSettings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ValuationMethod,
			&entity.'fifo',,
			&entity.'lifo',,
			&entity.'weightedAverage',,
			&entity.'movingAverage',,
			&entity.'standardCost',,
			&entity.'specificId',
			&entity.CostLayerGranularity,
			&entity.'product',,
			&entity.'productLocation',,
			&entity.'productLocationLot',,
			&entity.'serialNumber',
			&entity.DefaultInventoryAccountId,
			&entity.DefaultCogsAccountId,
			&entity.DefaultInventoryAdjustmentAccountId,
			&entity.DefaultInventoryVarianceAccountId,
			&entity.CogsRecognitionTiming,
			&entity.'onSale',,
			&entity.'onDelivery',,
			&entity.'onPayment',
			&entity.AllowNegativeInventory,
			&entity.RevalueOnPurchase,
			&entity.RoundUnitCostToDecimals,
			&entity.RevaluationFrequency,
			&entity.'realTime',,
			&entity.'daily',,
			&entity.'monthly',,
			&entity.'manual',
			&entity.IsActive,
			&entity.EffectiveFrom,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_valuation_settings: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


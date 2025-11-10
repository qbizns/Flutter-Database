package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/inventory"
)

type InventoryRepository struct {
	db *DB
}

func NewInventoryRepository(db *DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) List(ctx context.Context, orgID uuid.UUID, filters inventory.InventoryTransactionFilters) ([]inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, filters.TransactionType)
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY transaction_date DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txList []inventory.InventoryTransaction
	for rows.Next() {
		var tx inventory.InventoryTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
			&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
			&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		tx.Metadata = metadata
		txList = append(txList, tx)
	}

	return txList, rows.Err()
}

func (r *InventoryRepository) Count(ctx context.Context, orgID uuid.UUID, filters inventory.InventoryTransactionFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM inventory_transactions WHERE organization_id = $1"
	args := []interface{}{orgID}
	argCount := 1

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, filters.TransactionType)
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *InventoryRepository) Create(ctx context.Context, transaction *inventory.InventoryTransaction) error {
	if err := r.db.SetOrganizationContext(ctx, transaction.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_transactions (
			id, organization_id, product_id, transaction_type, quantity, unit,
			balance_after, sale_id, reference_number, unit_cost, total_cost,
			transaction_date, notes, reason, metadata, created_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		transaction.ID, transaction.OrganizationID, transaction.ProductID, transaction.TransactionType, transaction.Quantity, transaction.Unit,
		transaction.BalanceAfter, transaction.SaleID, transaction.ReferenceNumber, transaction.UnitCost, transaction.TotalCost,
		transaction.TransactionDate, transaction.Notes, transaction.Reason, transaction.Metadata, transaction.CreatedAt, transaction.CreatedBy,
	)
	return err
}

func (r *InventoryRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND id = $2
	`

	var tx inventory.InventoryTransaction
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
		&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
		&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tx.Metadata = metadata
	return &tx, nil
}

func (r *InventoryRepository) GetByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, limit int, offset int) ([]inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND product_id = $2
		ORDER BY transaction_date DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txList []inventory.InventoryTransaction
	for rows.Next() {
		var tx inventory.InventoryTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
			&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
			&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		tx.Metadata = metadata
		txList = append(txList, tx)
	}

	return txList, rows.Err()
}

func (r *InventoryRepository) GetBySale(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) ([]inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND sale_id = $2
		ORDER BY transaction_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txList []inventory.InventoryTransaction
	for rows.Next() {
		var tx inventory.InventoryTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
			&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
			&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		tx.Metadata = metadata
		txList = append(txList, tx)
	}

	return txList, rows.Err()
}

func (r *InventoryRepository) Update(ctx context.Context, transaction *inventory.InventoryTransaction) error {
	if err := r.db.SetOrganizationContext(ctx, transaction.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transactions SET
			transaction_type = $3, quantity = $4, unit = $5, balance_after = $6,
			sale_id = $7, reference_number = $8, unit_cost = $9, total_cost = $10,
			transaction_date = $11, notes = $12, reason = $13, metadata = $14
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		transaction.OrganizationID, transaction.ID,
		transaction.TransactionType, transaction.Quantity, transaction.Unit, transaction.BalanceAfter,
		transaction.SaleID, transaction.ReferenceNumber, transaction.UnitCost, transaction.TotalCost,
		transaction.TransactionDate, transaction.Notes, transaction.Reason, transaction.Metadata,
	)
	return err
}

func (r *InventoryRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		DELETE FROM inventory_transactions
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

func (r *InventoryRepository) GetLatestBalanceForProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) (*inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND product_id = $2
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT 1
	`

	var tx inventory.InventoryTransaction
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, productID).Scan(
		&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
		&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
		&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tx.Metadata = metadata
	return &tx, nil
}

// ============================================================================
// INVENTORY VALUATION SETTINGS REPOSITORY
// ============================================================================

type InventoryValuationSettingsRepository struct {
	db *DB
}

func NewInventoryValuationSettingsRepository(db *DB) *InventoryValuationSettingsRepository {
	return &InventoryValuationSettingsRepository{db: db}
}

func (r *InventoryValuationSettingsRepository) Get(ctx context.Context, orgID uuid.UUID) (*inventory.InventoryValuationSettings, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, valuation_method, cost_layer_granularity,
		       default_inventory_account_id, default_cogs_account_id,
		       default_inventory_adjustment_account_id, default_inventory_variance_account_id,
		       cogs_recognition_timing, allow_negative_inventory, revalue_on_purchase,
		       round_unit_cost_to_decimals, auto_create_cost_layers_on_purchase,
		       auto_consume_cost_layers_on_sale, recalculate_inventory_value_on_adjustment,
		       created_at, updated_at
		FROM inventory_valuation_settings
		WHERE organization_id = $1
	`

	var s inventory.InventoryValuationSettings
	err := r.db.Pool.QueryRow(ctx, query, orgID).Scan(
		&s.ID, &s.OrganizationID, &s.ValuationMethod, &s.CostLayerGranularity,
		&s.DefaultInventoryAccountID, &s.DefaultCOGSAccountID,
		&s.DefaultInventoryAdjustmentAccountID, &s.DefaultInventoryVarianceAccountID,
		&s.COGSRecognitionTiming, &s.AllowNegativeInventory, &s.RevalueOnPurchase,
		&s.RoundUnitCostToDecimals, &s.AutoCreateCostLayersOnPurchase,
		&s.AutoConsumeCostLayersOnSale, &s.RecalculateInventoryValueOnAdjustment,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *InventoryValuationSettingsRepository) Create(ctx context.Context, settings *inventory.InventoryValuationSettings) error {
	if err := r.db.SetOrganizationContext(ctx, settings.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_valuation_settings (
			id, organization_id, valuation_method, cost_layer_granularity,
			default_inventory_account_id, default_cogs_account_id,
			default_inventory_adjustment_account_id, default_inventory_variance_account_id,
			cogs_recognition_timing, allow_negative_inventory, revalue_on_purchase,
			round_unit_cost_to_decimals, auto_create_cost_layers_on_purchase,
			auto_consume_cost_layers_on_sale, recalculate_inventory_value_on_adjustment,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		settings.ID, settings.OrganizationID, settings.ValuationMethod, settings.CostLayerGranularity,
		settings.DefaultInventoryAccountID, settings.DefaultCOGSAccountID,
		settings.DefaultInventoryAdjustmentAccountID, settings.DefaultInventoryVarianceAccountID,
		settings.COGSRecognitionTiming, settings.AllowNegativeInventory, settings.RevalueOnPurchase,
		settings.RoundUnitCostToDecimals, settings.AutoCreateCostLayersOnPurchase,
		settings.AutoConsumeCostLayersOnSale, settings.RecalculateInventoryValueOnAdjustment,
		settings.CreatedAt, settings.UpdatedAt,
	)
	return err
}

func (r *InventoryValuationSettingsRepository) Update(ctx context.Context, settings *inventory.InventoryValuationSettings) error {
	if err := r.db.SetOrganizationContext(ctx, settings.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_valuation_settings SET
			valuation_method = $3, cost_layer_granularity = $4,
			default_inventory_account_id = $5, default_cogs_account_id = $6,
			default_inventory_adjustment_account_id = $7, default_inventory_variance_account_id = $8,
			cogs_recognition_timing = $9, allow_negative_inventory = $10, revalue_on_purchase = $11,
			round_unit_cost_to_decimals = $12, auto_create_cost_layers_on_purchase = $13,
			auto_consume_cost_layers_on_sale = $14, recalculate_inventory_value_on_adjustment = $15,
			updated_at = $16
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		settings.OrganizationID, settings.ID,
		settings.ValuationMethod, settings.CostLayerGranularity,
		settings.DefaultInventoryAccountID, settings.DefaultCOGSAccountID,
		settings.DefaultInventoryAdjustmentAccountID, settings.DefaultInventoryVarianceAccountID,
		settings.COGSRecognitionTiming, settings.AllowNegativeInventory, settings.RevalueOnPurchase,
		settings.RoundUnitCostToDecimals, settings.AutoCreateCostLayersOnPurchase,
		settings.AutoConsumeCostLayersOnSale, settings.RecalculateInventoryValueOnAdjustment,
		settings.UpdatedAt,
	)
	return err
}

// ============================================================================
// INVENTORY COST LAYER REPOSITORY
// ============================================================================

type InventoryCostLayerRepository struct {
	db *DB
}

func NewInventoryCostLayerRepository(db *DB) *InventoryCostLayerRepository {
	return &InventoryCostLayerRepository{db: db}
}

func (r *InventoryCostLayerRepository) List(ctx context.Context, orgID uuid.UUID, productID *uuid.UUID, locationID *uuid.UUID) ([]inventory.InventoryCostLayer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, location_id, lot_number, serial_number,
		       layer_date, unit_cost, original_quantity, remaining_quantity, uom_code,
		       source_transaction_type, source_transaction_id, source_reference,
		       is_fully_consumed, consumed_at, metadata, created_at, updated_at
		FROM inventory_cost_layers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if productID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *productID)
	}

	if locationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *locationID)
	}

	query += " ORDER BY layer_date ASC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layers []inventory.InventoryCostLayer
	for rows.Next() {
		var layer inventory.InventoryCostLayer
		err := rows.Scan(
			&layer.ID, &layer.OrganizationID, &layer.ProductID, &layer.LocationID,
			&layer.LotNumber, &layer.SerialNumber, &layer.LayerDate, &layer.UnitCost,
			&layer.OriginalQuantity, &layer.RemainingQuantity, &layer.UOMCode,
			&layer.SourceTransactionType, &layer.SourceTransactionID, &layer.SourceReference,
			&layer.IsFullyConsumed, &layer.ConsumedAt, &layer.Metadata,
			&layer.CreatedAt, &layer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}

	return layers, rows.Err()
}

func (r *InventoryCostLayerRepository) Create(ctx context.Context, layer *inventory.InventoryCostLayer) error {
	if err := r.db.SetOrganizationContext(ctx, layer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_cost_layers (
			id, organization_id, product_id, location_id, lot_number, serial_number,
			layer_date, unit_cost, original_quantity, remaining_quantity, uom_code,
			source_transaction_type, source_transaction_id, source_reference,
			is_fully_consumed, consumed_at, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		layer.ID, layer.OrganizationID, layer.ProductID, layer.LocationID,
		layer.LotNumber, layer.SerialNumber, layer.LayerDate, layer.UnitCost,
		layer.OriginalQuantity, layer.RemainingQuantity, layer.UOMCode,
		layer.SourceTransactionType, layer.SourceTransactionID, layer.SourceReference,
		layer.IsFullyConsumed, layer.ConsumedAt, layer.Metadata,
		layer.CreatedAt, layer.UpdatedAt,
	)
	return err
}

func (r *InventoryCostLayerRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.InventoryCostLayer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, location_id, lot_number, serial_number,
		       layer_date, unit_cost, original_quantity, remaining_quantity, uom_code,
		       source_transaction_type, source_transaction_id, source_reference,
		       is_fully_consumed, consumed_at, metadata, created_at, updated_at
		FROM inventory_cost_layers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var layer inventory.InventoryCostLayer
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&layer.ID, &layer.OrganizationID, &layer.ProductID, &layer.LocationID,
		&layer.LotNumber, &layer.SerialNumber, &layer.LayerDate, &layer.UnitCost,
		&layer.OriginalQuantity, &layer.RemainingQuantity, &layer.UOMCode,
		&layer.SourceTransactionType, &layer.SourceTransactionID, &layer.SourceReference,
		&layer.IsFullyConsumed, &layer.ConsumedAt, &layer.Metadata,
		&layer.CreatedAt, &layer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &layer, nil
}

func (r *InventoryCostLayerRepository) Update(ctx context.Context, layer *inventory.InventoryCostLayer) error {
	if err := r.db.SetOrganizationContext(ctx, layer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_cost_layers SET
			product_id = $3, location_id = $4, lot_number = $5, serial_number = $6,
			layer_date = $7, unit_cost = $8, original_quantity = $9, remaining_quantity = $10,
			uom_code = $11, source_transaction_type = $12, source_transaction_id = $13,
			source_reference = $14, is_fully_consumed = $15, consumed_at = $16,
			metadata = $17, updated_at = $18
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		layer.OrganizationID, layer.ID,
		layer.ProductID, layer.LocationID, layer.LotNumber, layer.SerialNumber,
		layer.LayerDate, layer.UnitCost, layer.OriginalQuantity, layer.RemainingQuantity,
		layer.UOMCode, layer.SourceTransactionType, layer.SourceTransactionID,
		layer.SourceReference, layer.IsFullyConsumed, layer.ConsumedAt,
		layer.Metadata, layer.UpdatedAt,
	)
	return err
}

func (r *InventoryCostLayerRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_cost_layers
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *InventoryCostLayerRepository) GetAvailableLayers(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, locationID *uuid.UUID) ([]inventory.InventoryCostLayer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, location_id, lot_number, serial_number,
		       layer_date, unit_cost, original_quantity, remaining_quantity, uom_code,
		       source_transaction_type, source_transaction_id, source_reference,
		       is_fully_consumed, consumed_at, metadata, created_at, updated_at
		FROM inventory_cost_layers
		WHERE organization_id = $1 AND product_id = $2 AND remaining_quantity > 0 AND deleted_at IS NULL
	`

	args := []interface{}{orgID, productID}
	if locationID != nil {
		query += " AND location_id = $3"
		args = append(args, *locationID)
	}

	query += " ORDER BY layer_date ASC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layers []inventory.InventoryCostLayer
	for rows.Next() {
		var layer inventory.InventoryCostLayer
		err := rows.Scan(
			&layer.ID, &layer.OrganizationID, &layer.ProductID, &layer.LocationID,
			&layer.LotNumber, &layer.SerialNumber, &layer.LayerDate, &layer.UnitCost,
			&layer.OriginalQuantity, &layer.RemainingQuantity, &layer.UOMCode,
			&layer.SourceTransactionType, &layer.SourceTransactionID, &layer.SourceReference,
			&layer.IsFullyConsumed, &layer.ConsumedAt, &layer.Metadata,
			&layer.CreatedAt, &layer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}

	return layers, rows.Err()
}

func (r *InventoryCostLayerRepository) ConsumeQuantity(ctx context.Context, layerID uuid.UUID, quantity float64) error {
	query := `
		UPDATE inventory_cost_layers
		SET remaining_quantity = remaining_quantity - $2,
		    is_fully_consumed = (remaining_quantity - $2 <= 0),
		    consumed_at = CASE WHEN (remaining_quantity - $2 <= 0) THEN $3 ELSE consumed_at END,
		    updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, layerID, quantity, time.Now())
	return err
}

func (r *InventoryCostLayerRepository) GetOldestLayer(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, locationID *uuid.UUID) (*inventory.InventoryCostLayer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, location_id, lot_number, serial_number,
		       layer_date, unit_cost, original_quantity, remaining_quantity, uom_code,
		       source_transaction_type, source_transaction_id, source_reference,
		       is_fully_consumed, consumed_at, metadata, created_at, updated_at
		FROM inventory_cost_layers
		WHERE organization_id = $1 AND product_id = $2 AND remaining_quantity > 0 AND deleted_at IS NULL
	`

	args := []interface{}{orgID, productID}
	if locationID != nil {
		query += " AND location_id = $3"
		args = append(args, *locationID)
	}

	query += " ORDER BY layer_date ASC LIMIT 1"

	var layer inventory.InventoryCostLayer
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&layer.ID, &layer.OrganizationID, &layer.ProductID, &layer.LocationID,
		&layer.LotNumber, &layer.SerialNumber, &layer.LayerDate, &layer.UnitCost,
		&layer.OriginalQuantity, &layer.RemainingQuantity, &layer.UOMCode,
		&layer.SourceTransactionType, &layer.SourceTransactionID, &layer.SourceReference,
		&layer.IsFullyConsumed, &layer.ConsumedAt, &layer.Metadata,
		&layer.CreatedAt, &layer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &layer, nil
}

func (r *InventoryCostLayerRepository) GetNewestLayer(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, locationID *uuid.UUID) (*inventory.InventoryCostLayer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, location_id, lot_number, serial_number,
		       layer_date, unit_cost, original_quantity, remaining_quantity, uom_code,
		       source_transaction_type, source_transaction_id, source_reference,
		       is_fully_consumed, consumed_at, metadata, created_at, updated_at
		FROM inventory_cost_layers
		WHERE organization_id = $1 AND product_id = $2 AND remaining_quantity > 0 AND deleted_at IS NULL
	`

	args := []interface{}{orgID, productID}
	if locationID != nil {
		query += " AND location_id = $3"
		args = append(args, *locationID)
	}

	query += " ORDER BY layer_date DESC LIMIT 1"

	var layer inventory.InventoryCostLayer
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&layer.ID, &layer.OrganizationID, &layer.ProductID, &layer.LocationID,
		&layer.LotNumber, &layer.SerialNumber, &layer.LayerDate, &layer.UnitCost,
		&layer.OriginalQuantity, &layer.RemainingQuantity, &layer.UOMCode,
		&layer.SourceTransactionType, &layer.SourceTransactionID, &layer.SourceReference,
		&layer.IsFullyConsumed, &layer.ConsumedAt, &layer.Metadata,
		&layer.CreatedAt, &layer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &layer, nil
}

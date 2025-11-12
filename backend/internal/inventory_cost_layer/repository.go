package inventory_cost_layer

import (
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

// Repository handles database operations for InventoryCostLayers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new InventoryCostLayers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// InventoryCostLayers represents a inventory_cost_layers entity
type InventoryCostLayers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	LotNumber *string `json:"lot_number" db:"lot_number"`
	SerialNumber *string `json:"serial_number" db:"serial_number"`
	LayerDate time.Time `json:"layer_date" db:"layer_date"`
	UnitCost float64 `json:"unit_cost" db:"unit_cost"`
	OriginalQuantity float64 `json:"original_quantity" db:"original_quantity"`
	RemainingQuantity float64 `json:"remaining_quantity" db:"remaining_quantity"`
	UomCode *string `json:"uom_code" db:"uom_code"`
	SourceTransactionType *string `json:"source_transaction_type" db:"source_transaction_type"`
	SourceTransactionId *uuid.UUID `json:"source_transaction_id" db:"source_transaction_id"`
	SourceReference *string `json:"source_reference" db:"source_reference"`
	IsFullyConsumed *bool `json:"is_fully_consumed" db:"is_fully_consumed"`
	ConsumedAt *time.Time `json:"consumed_at" db:"consumed_at"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new inventory_cost_layers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *InventoryCostLayers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "inventory_cost_layers", duration, nil)
	}()

	query := `
		INSERT INTO inventory_cost_layers (
			, organization_id
			, product_id
			, location_id
			, lot_number
			, serial_number
			, layer_date
			, unit_cost
			, original_quantity
			, remaining_quantity
			, uom_code
			, source_transaction_type
			, source_transaction_id
			, source_reference
			, is_fully_consumed
			, consumed_at
			, metadata
			, deleted_at
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
			, $20
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.LocationId,
		entity.LotNumber,
		entity.SerialNumber,
		entity.LayerDate,
		entity.UnitCost,
		entity.OriginalQuantity,
		entity.RemainingQuantity,
		entity.UomCode,
		entity.SourceTransactionType,
		entity.SourceTransactionId,
		entity.SourceReference,
		entity.IsFullyConsumed,
		entity.ConsumedAt,
		entity.Metadata,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create inventory_cost_layers", zap.Error(err))
		return fmt.Errorf("failed to create inventory_cost_layers: %w", err)
	}

	r.logger.Info("created inventory_cost_layers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a inventory_cost_layers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*InventoryCostLayers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_cost_layers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, product_id
			, location_id
			, lot_number
			, serial_number
			, layer_date
			, unit_cost
			, original_quantity
			, remaining_quantity
			, uom_code
			, source_transaction_type
			, source_transaction_id
			, source_reference
			, is_fully_consumed
			, consumed_at
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM inventory_cost_layers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity InventoryCostLayers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ProductId,
		&entity.LocationId,
		&entity.LotNumber,
		&entity.SerialNumber,
		&entity.LayerDate,
		&entity.UnitCost,
		&entity.OriginalQuantity,
		&entity.RemainingQuantity,
		&entity.UomCode,
		&entity.SourceTransactionType,
		&entity.SourceTransactionId,
		&entity.SourceReference,
		&entity.IsFullyConsumed,
		&entity.ConsumedAt,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("inventory_cost_layers not found")
	}

	if err != nil {
		r.logger.Error("failed to get inventory_cost_layers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get inventory_cost_layers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of inventory_cost_layers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*InventoryCostLayers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_cost_layers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_cost_layers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_cost_layers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, location_id
			, lot_number
			, serial_number
			, layer_date
			, unit_cost
			, original_quantity
			, remaining_quantity
			, uom_code
			, source_transaction_type
			, source_transaction_id
			, source_reference
			, is_fully_consumed
			, consumed_at
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM inventory_cost_layers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_cost_layers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_cost_layers: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryCostLayers
	for rows.Next() {
		var entity InventoryCostLayers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.LocationId,
			&entity.LotNumber,
			&entity.SerialNumber,
			&entity.LayerDate,
			&entity.UnitCost,
			&entity.OriginalQuantity,
			&entity.RemainingQuantity,
			&entity.UomCode,
			&entity.SourceTransactionType,
			&entity.SourceTransactionId,
			&entity.SourceReference,
			&entity.IsFullyConsumed,
			&entity.ConsumedAt,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_cost_layers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating inventory_cost_layers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing inventory_cost_layers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *InventoryCostLayers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_cost_layers", duration, nil)
	}()

	query := `
		UPDATE inventory_cost_layers
		SET
			, organization_id = $2
			, product_id = $3
			, location_id = $4
			, lot_number = $5
			, serial_number = $6
			, layer_date = $7
			, unit_cost = $8
			, original_quantity = $9
			, remaining_quantity = $10
			, uom_code = $11
			, source_transaction_type = $12
			, source_transaction_id = $13
			, source_reference = $14
			, is_fully_consumed = $15
			, consumed_at = $16
			, metadata = $17
			, updated_at = $19
			, deleted_at = $20
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.LocationId,
		entity.LotNumber,
		entity.SerialNumber,
		entity.LayerDate,
		entity.UnitCost,
		entity.OriginalQuantity,
		entity.RemainingQuantity,
		entity.UomCode,
		entity.SourceTransactionType,
		entity.SourceTransactionId,
		entity.SourceReference,
		entity.IsFullyConsumed,
		entity.ConsumedAt,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update inventory_cost_layers", zap.Error(err))
		return fmt.Errorf("failed to update inventory_cost_layers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_cost_layers not found or already deleted")
	}

	r.logger.Info("updated inventory_cost_layers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a inventory_cost_layers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_cost_layers", duration, nil)
	}()

	query := `
		UPDATE inventory_cost_layers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete inventory_cost_layers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete inventory_cost_layers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_cost_layers not found or already deleted")
	}

	r.logger.Info("deleted inventory_cost_layers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves inventory_cost_layers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*InventoryCostLayers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_cost_layers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_cost_layers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_cost_layers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, location_id
			, lot_number
			, serial_number
			, layer_date
			, unit_cost
			, original_quantity
			, remaining_quantity
			, uom_code
			, source_transaction_type
			, source_transaction_id
			, source_reference
			, is_fully_consumed
			, consumed_at
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM inventory_cost_layers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_cost_layers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_cost_layers: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryCostLayers
	for rows.Next() {
		var entity InventoryCostLayers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.LocationId,
			&entity.LotNumber,
			&entity.SerialNumber,
			&entity.LayerDate,
			&entity.UnitCost,
			&entity.OriginalQuantity,
			&entity.RemainingQuantity,
			&entity.UomCode,
			&entity.SourceTransactionType,
			&entity.SourceTransactionId,
			&entity.SourceReference,
			&entity.IsFullyConsumed,
			&entity.ConsumedAt,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_cost_layers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


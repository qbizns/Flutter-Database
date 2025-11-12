package inventory_transfer_item

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

// Repository handles database operations for InventoryTransferItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new InventoryTransferItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// InventoryTransferItems represents a inventory_transfer_items entity
type InventoryTransferItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	InventoryTransferId uuid.UUID `json:"inventory_transfer_id" db:"inventory_transfer_id"`
	ProductId *uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	ProductName string `json:"product_name" db:"product_name"`
	ProductSku *string `json:"product_sku" db:"product_sku"`
	QuantityRequested float64 `json:"quantity_requested" db:"quantity_requested"`
	QuantityShipped *float64 `json:"quantity_shipped" db:"quantity_shipped"`
	QuantityReceived *float64 `json:"quantity_received" db:"quantity_received"`
	UnitOfMeasure *string `json:"unit_of_measure" db:"unit_of_measure"`
	UnitCost *float64 `json:"unit_cost" db:"unit_cost"`
	TotalCost *float64 `json:"total_cost" db:"total_cost"`
	ItemStatus *string `json:"item_status" db:"item_status"`
	VarianceQuantity *float64 `json:"variance_quantity" db:"variance_quantity"`
	VarianceReason *string `json:"variance_reason" db:"variance_reason"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	QuantityRequested *string `json:"quantity_requested" db:"quantity_requested"`
	QuantityShipped *string `json:"quantity_shipped" db:"quantity_shipped"`
	QuantityReceived *string `json:"quantity_received" db:"quantity_received"`
	QuantityShipped *string `json:"quantity_shipped" db:"quantity_shipped"`
	QuantityReceived *string `json:"quantity_received" db:"quantity_received"`
	(unitCost *string `json:"(unit_cost" db:"(unit_cost"`
	(totalCost *string `json:"(total_cost" db:"(total_cost"`
}

// Create inserts a new inventory_transfer_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *InventoryTransferItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "inventory_transfer_items", duration, nil)
	}()

	query := `
		INSERT INTO inventory_transfer_items (
			, organization_id
			, inventory_transfer_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, unit_of_measure
			, unit_cost
			, total_cost
			, item_status
			, variance_quantity
			, variance_reason
			, notes
			, metadata
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, quantity_shipped
			, quantity_received
			, (unit_cost
			, (total_cost
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
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $27
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.InventoryTransferId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.ProductName,
		entity.ProductSku,
		entity.QuantityRequested,
		entity.QuantityShipped,
		entity.QuantityReceived,
		entity.UnitOfMeasure,
		entity.UnitCost,
		entity.TotalCost,
		entity.ItemStatus,
		entity.VarianceQuantity,
		entity.VarianceReason,
		entity.Notes,
		entity.Metadata,
		entity.QuantityRequested,
		entity.QuantityShipped,
		entity.QuantityReceived,
		entity.QuantityShipped,
		entity.QuantityReceived,
		entity.(unitCost,
		entity.(totalCost,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create inventory_transfer_items", zap.Error(err))
		return fmt.Errorf("failed to create inventory_transfer_items: %w", err)
	}

	r.logger.Info("created inventory_transfer_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a inventory_transfer_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*InventoryTransferItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transfer_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, inventory_transfer_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, unit_of_measure
			, unit_cost
			, total_cost
			, item_status
			, variance_quantity
			, variance_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, quantity_shipped
			, quantity_received
			, (unit_cost
			, (total_cost
		FROM inventory_transfer_items
		WHERE id = $1
		
	`

	var entity InventoryTransferItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.InventoryTransferId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.ProductName,
		&entity.ProductSku,
		&entity.QuantityRequested,
		&entity.QuantityShipped,
		&entity.QuantityReceived,
		&entity.UnitOfMeasure,
		&entity.UnitCost,
		&entity.TotalCost,
		&entity.ItemStatus,
		&entity.VarianceQuantity,
		&entity.VarianceReason,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.QuantityRequested,
		&entity.QuantityShipped,
		&entity.QuantityReceived,
		&entity.QuantityShipped,
		&entity.QuantityReceived,
		&entity.(unitCost,
		&entity.(totalCost,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("inventory_transfer_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get inventory_transfer_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get inventory_transfer_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of inventory_transfer_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*InventoryTransferItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transfer_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_transfer_items
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_transfer_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, inventory_transfer_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, unit_of_measure
			, unit_cost
			, total_cost
			, item_status
			, variance_quantity
			, variance_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, quantity_shipped
			, quantity_received
			, (unit_cost
			, (total_cost
		FROM inventory_transfer_items
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_transfer_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_transfer_items: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryTransferItems
	for rows.Next() {
		var entity InventoryTransferItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.InventoryTransferId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.ProductName,
			&entity.ProductSku,
			&entity.QuantityRequested,
			&entity.QuantityShipped,
			&entity.QuantityReceived,
			&entity.UnitOfMeasure,
			&entity.UnitCost,
			&entity.TotalCost,
			&entity.ItemStatus,
			&entity.VarianceQuantity,
			&entity.VarianceReason,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.QuantityRequested,
			&entity.QuantityShipped,
			&entity.QuantityReceived,
			&entity.QuantityShipped,
			&entity.QuantityReceived,
			&entity.(unitCost,
			&entity.(totalCost,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_transfer_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating inventory_transfer_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing inventory_transfer_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *InventoryTransferItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_transfer_items", duration, nil)
	}()

	query := `
		UPDATE inventory_transfer_items
		SET
			, organization_id = $2
			, inventory_transfer_id = $3
			, product_id = $4
			, product_variant_id = $5
			, product_name = $6
			, product_sku = $7
			, quantity_requested = $8
			, quantity_shipped = $9
			, quantity_received = $10
			, unit_of_measure = $11
			, unit_cost = $12
			, total_cost = $13
			, item_status = $14
			, variance_quantity = $15
			, variance_reason = $16
			, notes = $17
			, metadata = $18
			, updated_at = $20
			, quantity_requested = $21
			, quantity_shipped = $22
			, quantity_received = $23
			, quantity_shipped = $24
			, quantity_received = $25
			, (unit_cost = $26
			, (total_cost = $27
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $28
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.InventoryTransferId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.ProductName,
		entity.ProductSku,
		entity.QuantityRequested,
		entity.QuantityShipped,
		entity.QuantityReceived,
		entity.UnitOfMeasure,
		entity.UnitCost,
		entity.TotalCost,
		entity.ItemStatus,
		entity.VarianceQuantity,
		entity.VarianceReason,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.QuantityRequested,
		entity.QuantityShipped,
		entity.QuantityReceived,
		entity.QuantityShipped,
		entity.QuantityReceived,
		entity.(unitCost,
		entity.(totalCost,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update inventory_transfer_items", zap.Error(err))
		return fmt.Errorf("failed to update inventory_transfer_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_transfer_items not found or already deleted")
	}

	r.logger.Info("updated inventory_transfer_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a inventory_transfer_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "inventory_transfer_items", duration, nil)
	}()

	query := `DELETE FROM inventory_transfer_items WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete inventory_transfer_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete inventory_transfer_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_transfer_items not found")
	}

	r.logger.Info("deleted inventory_transfer_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves inventory_transfer_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*InventoryTransferItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transfer_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_transfer_items
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_transfer_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, inventory_transfer_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, unit_of_measure
			, unit_cost
			, total_cost
			, item_status
			, variance_quantity
			, variance_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, quantity_requested
			, quantity_shipped
			, quantity_received
			, quantity_shipped
			, quantity_received
			, (unit_cost
			, (total_cost
		FROM inventory_transfer_items
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_transfer_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_transfer_items: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryTransferItems
	for rows.Next() {
		var entity InventoryTransferItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.InventoryTransferId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.ProductName,
			&entity.ProductSku,
			&entity.QuantityRequested,
			&entity.QuantityShipped,
			&entity.QuantityReceived,
			&entity.UnitOfMeasure,
			&entity.UnitCost,
			&entity.TotalCost,
			&entity.ItemStatus,
			&entity.VarianceQuantity,
			&entity.VarianceReason,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.QuantityRequested,
			&entity.QuantityShipped,
			&entity.QuantityReceived,
			&entity.QuantityShipped,
			&entity.QuantityReceived,
			&entity.(unitCost,
			&entity.(totalCost,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_transfer_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


package product_batch

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

// Repository handles database operations for ProductBatches
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ProductBatches repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ProductBatches represents a product_batches entity
type ProductBatches struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	BatchNumber string `json:"batch_number" db:"batch_number"`
	LotNumber *string `json:"lot_number" db:"lot_number"`
	Status *string `json:"status" db:"status"`
	InitialQuantity float64 `json:"initial_quantity" db:"initial_quantity"`
	CurrentQuantity float64 `json:"current_quantity" db:"current_quantity"`
	UnitOfMeasure *string `json:"unit_of_measure" db:"unit_of_measure"`
	ManufacturingDate *time.Time `json:"manufacturing_date" db:"manufacturing_date"`
	ExpirationDate *time.Time `json:"expiration_date" db:"expiration_date"`
	ReceivedDate time.Time `json:"received_date" db:"received_date"`
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id" db:"purchase_order_id"`
	SupplierId *uuid.UUID `json:"supplier_id" db:"supplier_id"`
	SupplierBatchNumber *string `json:"supplier_batch_number" db:"supplier_batch_number"`
	UnitCost *float64 `json:"unit_cost" db:"unit_cost"`
	TotalCost *float64 `json:"total_cost" db:"total_cost"`
	// 	QualityStatus *string `json:"quality_status" db:"quality_status"`
	QualityCheckDate *time.Time `json:"quality_check_date" db:"quality_check_date"`
	QualityCheckedBy *uuid.UUID `json:"quality_checked_by" db:"quality_checked_by"`
	QualityNotes *string `json:"quality_notes" db:"quality_notes"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	InitialQuantity *string `json:"initial_quantity" db:"initial_quantity"`
	CurrentQuantity *string `json:"current_quantity" db:"current_quantity"`
	CurrentQuantity *string `json:"current_quantity" db:"current_quantity"`
	ExpirationDate *string `json:"expiration_date" db:"expiration_date"`
}

// Create inserts a new product_batches record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ProductBatches) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "product_batches", duration, nil)
	}()

	query := `
		INSERT INTO product_batches (
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, batch_number
			, lot_number
			, status
			, initial_quantity
			, current_quantity
			, unit_of_measure
			, manufacturing_date
			, expiration_date
			, received_date
			, purchase_order_id
			, supplier_id
			, supplier_batch_number
			, unit_cost
			, total_cost
			, quality_status
			, quality_check_date
			, quality_checked_by
			, quality_notes
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, initial_quantity
			, current_quantity
			, current_quantity
			, expiration_date
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
			, $28
			, $29
			, $30
			, $31
			, $32
			, $33
			, $34
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.LocationId,
		entity.BatchNumber,
		entity.LotNumber,
		entity.Status,
		entity.InitialQuantity,
		entity.CurrentQuantity,
		entity.UnitOfMeasure,
		entity.ManufacturingDate,
		entity.ExpirationDate,
		entity.ReceivedDate,
		entity.PurchaseOrderId,
		entity.SupplierId,
		entity.SupplierBatchNumber,
		entity.UnitCost,
		entity.TotalCost,
		entity.QualityStatus,
		entity.QualityCheckDate,
		entity.QualityCheckedBy,
		entity.QualityNotes,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.InitialQuantity,
		entity.CurrentQuantity,
		entity.CurrentQuantity,
		entity.ExpirationDate,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create product_batches", zap.Error(err))
		return fmt.Errorf("failed to create product_batches: %w", err)
	}

	r.logger.Info("created product_batches",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a product_batches by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ProductBatches, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_batches", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, batch_number
			, lot_number
			, initial_quantity
			, current_quantity
			, unit_of_measure
			, manufacturing_date
			, expiration_date
			, received_date
			, purchase_order_id
			, supplier_id
			, supplier_batch_number
			, unit_cost
			, total_cost
			, quality_status
			, quality_check_date
			, quality_checked_by
			, quality_notes
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, initial_quantity
			, current_quantity
			, current_quantity
			, expiration_date
		FROM product_batches
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ProductBatches
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.LocationId,
		&entity.BatchNumber,
		&entity.LotNumber,
		&entity.Status,
		&entity.InitialQuantity,
		&entity.CurrentQuantity,
		&entity.UnitOfMeasure,
		&entity.ManufacturingDate,
		&entity.ExpirationDate,
		&entity.ReceivedDate,
		&entity.PurchaseOrderId,
		&entity.SupplierId,
		&entity.SupplierBatchNumber,
		&entity.UnitCost,
		&entity.TotalCost,
		&entity.QualityStatus,
		&entity.QualityCheckDate,
		&entity.QualityCheckedBy,
		&entity.QualityNotes,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.InitialQuantity,
		&entity.CurrentQuantity,
		&entity.CurrentQuantity,
		&entity.ExpirationDate,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("product_batches not found")
	}

	if err != nil {
		r.logger.Error("failed to get product_batches", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get product_batches: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of product_batches records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ProductBatches, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_batches", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_batches
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_batches records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, batch_number
			, lot_number
			, initial_quantity
			, current_quantity
			, unit_of_measure
			, manufacturing_date
			, expiration_date
			, received_date
			, purchase_order_id
			, supplier_id
			, supplier_batch_number
			, unit_cost
			, total_cost
			, quality_status
			, quality_check_date
			, quality_checked_by
			, quality_notes
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, initial_quantity
			, current_quantity
			, current_quantity
			, expiration_date
		FROM product_batches
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_batches", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_batches: %w", err)
	}
	defer rows.Close()

	var entities []*ProductBatches
	for rows.Next() {
		var entity ProductBatches
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.LocationId,
			&entity.BatchNumber,
			&entity.LotNumber,
			&entity.Status,
			&entity.InitialQuantity,
			&entity.CurrentQuantity,
			&entity.UnitOfMeasure,
			&entity.ManufacturingDate,
			&entity.ExpirationDate,
			&entity.ReceivedDate,
			&entity.PurchaseOrderId,
			&entity.SupplierId,
			&entity.SupplierBatchNumber,
			&entity.UnitCost,
			&entity.TotalCost,
			&entity.QualityStatus,
			&entity.QualityCheckDate,
			&entity.QualityCheckedBy,
			&entity.QualityNotes,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.InitialQuantity,
			&entity.CurrentQuantity,
			&entity.CurrentQuantity,
			&entity.ExpirationDate,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_batches: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating product_batches rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing product_batches record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ProductBatches) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_batches", duration, nil)
	}()

	query := `
		UPDATE product_batches
		SET
			, organization_id = $2
			, product_id = $3
			, product_variant_id = $4
			, location_id = $5
			, batch_number = $6
			, lot_number = $7
			, status = $8
			, initial_quantity = $9
			, current_quantity = $10
			, unit_of_measure = $11
			, manufacturing_date = $12
			, expiration_date = $13
			, received_date = $14
			, purchase_order_id = $15
			, supplier_id = $16
			, supplier_batch_number = $17
			, unit_cost = $18
			, total_cost = $19
			, quality_status = $20
			, quality_check_date = $21
			, quality_checked_by = $22
			, quality_notes = $23
			, notes = $24
			, metadata = $25
			, updated_at = $27
			, deleted_at = $28
			, created_by = $29
			, updated_by = $30
			, initial_quantity = $31
			, current_quantity = $32
			, current_quantity = $33
			, expiration_date = $34
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $35
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.LocationId,
		entity.BatchNumber,
		entity.LotNumber,
		entity.Status,
		entity.InitialQuantity,
		entity.CurrentQuantity,
		entity.UnitOfMeasure,
		entity.ManufacturingDate,
		entity.ExpirationDate,
		entity.ReceivedDate,
		entity.PurchaseOrderId,
		entity.SupplierId,
		entity.SupplierBatchNumber,
		entity.UnitCost,
		entity.TotalCost,
		entity.QualityStatus,
		entity.QualityCheckDate,
		entity.QualityCheckedBy,
		entity.QualityNotes,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.InitialQuantity,
		entity.CurrentQuantity,
		entity.CurrentQuantity,
		entity.ExpirationDate,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update product_batches", zap.Error(err))
		return fmt.Errorf("failed to update product_batches: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_batches not found or already deleted")
	}

	r.logger.Info("updated product_batches",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a product_batches record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_batches", duration, nil)
	}()

	query := `
		UPDATE product_batches
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete product_batches", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete product_batches: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_batches not found or already deleted")
	}

	r.logger.Info("deleted product_batches", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves product_batches records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ProductBatches, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_batches", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_batches
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_batches records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, batch_number
			, lot_number
			, status
			, initial_quantity
			, current_quantity
			, unit_of_measure
			, manufacturing_date
			, expiration_date
			, received_date
			, purchase_order_id
			, supplier_id
			, supplier_batch_number
			, unit_cost
			, total_cost
			, quality_status
			, quality_check_date
			, quality_checked_by
			, quality_notes
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, initial_quantity
			, current_quantity
			, current_quantity
			, expiration_date
		FROM product_batches
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_batches by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_batches: %w", err)
	}
	defer rows.Close()

	var entities []*ProductBatches
	for rows.Next() {
		var entity ProductBatches
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.LocationId,
			&entity.BatchNumber,
			&entity.LotNumber,
			&entity.Status,
			&entity.InitialQuantity,
			&entity.CurrentQuantity,
			&entity.UnitOfMeasure,
			&entity.ManufacturingDate,
			&entity.ExpirationDate,
			&entity.ReceivedDate,
			&entity.PurchaseOrderId,
			&entity.SupplierId,
			&entity.SupplierBatchNumber,
			&entity.UnitCost,
			&entity.TotalCost,
			&entity.QualityStatus,
			&entity.QualityCheckDate,
			&entity.QualityCheckedBy,
			&entity.QualityNotes,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.InitialQuantity,
			&entity.CurrentQuantity,
			&entity.CurrentQuantity,
			&entity.ExpirationDate,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_batches: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


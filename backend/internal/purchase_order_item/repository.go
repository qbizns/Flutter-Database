package purchase_order_item

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

// Repository handles database operations for PurchaseOrderItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PurchaseOrderItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PurchaseOrderItems represents a purchase_order_items entity
type PurchaseOrderItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PurchaseOrderId uuid.UUID `json:"purchase_order_id" db:"purchase_order_id"`
	LineNumber int64 `json:"line_number" db:"line_number"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	QuantityOrdered float64 `json:"quantity_ordered" db:"quantity_ordered"`
	QuantityReceived *float64 `json:"quantity_received" db:"quantity_received"`
	UnitCost float64 `json:"unit_cost" db:"unit_cost"`
	Subtotal float64 `json:"subtotal" db:"subtotal"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date" db:"expected_delivery_date"`
	Notes *string `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new purchase_order_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PurchaseOrderItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "purchase_order_items", duration, nil)
	}()

	query := `
		INSERT INTO purchase_order_items (
			, organization_id
			, purchase_order_id
			, line_number
			, product_id
			, product_variant_id
			, quantity_ordered
			, quantity_received
			, unit_cost
			, subtotal
			, tax_amount
			, total_amount
			, expected_delivery_date
			, notes
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
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PurchaseOrderId,
		entity.LineNumber,
		entity.ProductId,
		entity.ProductVariantId,
		entity.QuantityOrdered,
		entity.QuantityReceived,
		entity.UnitCost,
		entity.Subtotal,
		entity.TaxAmount,
		entity.TotalAmount,
		entity.ExpectedDeliveryDate,
		entity.Notes,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create purchase_order_items", zap.Error(err))
		return fmt.Errorf("failed to create purchase_order_items: %w", err)
	}

	r.logger.Info("created purchase_order_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a purchase_order_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PurchaseOrderItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "purchase_order_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, purchase_order_id
			, line_number
			, product_id
			, product_variant_id
			, quantity_ordered
			, quantity_received
			, unit_cost
			, subtotal
			, tax_amount
			, total_amount
			, expected_delivery_date
			, notes
			, created_at
			, updated_at
			, deleted_at
		FROM purchase_order_items
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PurchaseOrderItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PurchaseOrderId,
		&entity.LineNumber,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.QuantityOrdered,
		&entity.QuantityReceived,
		&entity.UnitCost,
		&entity.Subtotal,
		&entity.TaxAmount,
		&entity.TotalAmount,
		&entity.ExpectedDeliveryDate,
		&entity.Notes,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("purchase_order_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get purchase_order_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get purchase_order_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of purchase_order_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PurchaseOrderItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "purchase_order_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM purchase_order_items
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase_order_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, purchase_order_id
			, line_number
			, product_id
			, product_variant_id
			, quantity_ordered
			, quantity_received
			, unit_cost
			, subtotal
			, tax_amount
			, total_amount
			, expected_delivery_date
			, notes
			, created_at
			, updated_at
			, deleted_at
		FROM purchase_order_items
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list purchase_order_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list purchase_order_items: %w", err)
	}
	defer rows.Close()

	var entities []*PurchaseOrderItems
	for rows.Next() {
		var entity PurchaseOrderItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PurchaseOrderId,
			&entity.LineNumber,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.QuantityOrdered,
			&entity.QuantityReceived,
			&entity.UnitCost,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.TotalAmount,
			&entity.ExpectedDeliveryDate,
			&entity.Notes,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase_order_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating purchase_order_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing purchase_order_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PurchaseOrderItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "purchase_order_items", duration, nil)
	}()

	query := `
		UPDATE purchase_order_items
		SET
			, organization_id = $2
			, purchase_order_id = $3
			, line_number = $4
			, product_id = $5
			, product_variant_id = $6
			, quantity_ordered = $7
			, quantity_received = $8
			, unit_cost = $9
			, subtotal = $10
			, tax_amount = $11
			, total_amount = $12
			, expected_delivery_date = $13
			, notes = $14
			, updated_at = $16
			, deleted_at = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PurchaseOrderId,
		entity.LineNumber,
		entity.ProductId,
		entity.ProductVariantId,
		entity.QuantityOrdered,
		entity.QuantityReceived,
		entity.UnitCost,
		entity.Subtotal,
		entity.TaxAmount,
		entity.TotalAmount,
		entity.ExpectedDeliveryDate,
		entity.Notes,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update purchase_order_items", zap.Error(err))
		return fmt.Errorf("failed to update purchase_order_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("purchase_order_items not found or already deleted")
	}

	r.logger.Info("updated purchase_order_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a purchase_order_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "purchase_order_items", duration, nil)
	}()

	query := `
		UPDATE purchase_order_items
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete purchase_order_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete purchase_order_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("purchase_order_items not found or already deleted")
	}

	r.logger.Info("deleted purchase_order_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves purchase_order_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PurchaseOrderItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "purchase_order_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM purchase_order_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase_order_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, purchase_order_id
			, line_number
			, product_id
			, product_variant_id
			, quantity_ordered
			, quantity_received
			, unit_cost
			, subtotal
			, tax_amount
			, total_amount
			, expected_delivery_date
			, notes
			, created_at
			, updated_at
			, deleted_at
		FROM purchase_order_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list purchase_order_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list purchase_order_items: %w", err)
	}
	defer rows.Close()

	var entities []*PurchaseOrderItems
	for rows.Next() {
		var entity PurchaseOrderItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PurchaseOrderId,
			&entity.LineNumber,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.QuantityOrdered,
			&entity.QuantityReceived,
			&entity.UnitCost,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.TotalAmount,
			&entity.ExpectedDeliveryDate,
			&entity.Notes,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase_order_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


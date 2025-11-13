package goods_receipt_item

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

// Repository handles database operations for GoodsReceiptItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new GoodsReceiptItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// GoodsReceiptItems represents a goods_receipt_items entity
type GoodsReceiptItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	GoodsReceiptId uuid.UUID `json:"goods_receipt_id" db:"goods_receipt_id"`
	PurchaseOrderItemId *uuid.UUID `json:"purchase_order_item_id" db:"purchase_order_item_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	QuantityReceived float64 `json:"quantity_received" db:"quantity_received"`
	QuantityAccepted *float64 `json:"quantity_accepted" db:"quantity_accepted"`
	QuantityRejected *float64 `json:"quantity_rejected" db:"quantity_rejected"`
	RejectionReason *string `json:"rejection_reason" db:"rejection_reason"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new goods_receipt_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *GoodsReceiptItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "goods_receipt_items", duration, nil)
	}()

	query := `
		INSERT INTO goods_receipt_items (
			, organization_id
			, goods_receipt_id
			, purchase_order_item_id
			, product_id
			, product_variant_id
			, quantity_received
			, quantity_accepted
			, quantity_rejected
			, rejection_reason
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
			, $12
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.GoodsReceiptId,
		entity.PurchaseOrderItemId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.QuantityReceived,
		entity.QuantityAccepted,
		entity.QuantityRejected,
		entity.RejectionReason,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create goods_receipt_items", zap.Error(err))
		return fmt.Errorf("failed to create goods_receipt_items: %w", err)
	}

	r.logger.Info("created goods_receipt_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a goods_receipt_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*GoodsReceiptItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "goods_receipt_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, goods_receipt_id
			, purchase_order_item_id
			, product_id
			, product_variant_id
			, quantity_received
			, quantity_accepted
			, quantity_rejected
			, rejection_reason
			, created_at
			, deleted_at
		FROM goods_receipt_items
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity GoodsReceiptItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.GoodsReceiptId,
		&entity.PurchaseOrderItemId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.QuantityReceived,
		&entity.QuantityAccepted,
		&entity.QuantityRejected,
		&entity.RejectionReason,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("goods_receipt_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get goods_receipt_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get goods_receipt_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of goods_receipt_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*GoodsReceiptItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "goods_receipt_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM goods_receipt_items
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count goods_receipt_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, goods_receipt_id
			, purchase_order_item_id
			, product_id
			, product_variant_id
			, quantity_received
			, quantity_accepted
			, quantity_rejected
			, rejection_reason
			, created_at
			, deleted_at
		FROM goods_receipt_items
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list goods_receipt_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list goods_receipt_items: %w", err)
	}
	defer rows.Close()

	var entities []*GoodsReceiptItems
	for rows.Next() {
		var entity GoodsReceiptItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.GoodsReceiptId,
			&entity.PurchaseOrderItemId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.QuantityReceived,
			&entity.QuantityAccepted,
			&entity.QuantityRejected,
			&entity.RejectionReason,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan goods_receipt_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating goods_receipt_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing goods_receipt_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *GoodsReceiptItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "goods_receipt_items", duration, nil)
	}()

	query := `
		UPDATE goods_receipt_items
		SET
			, organization_id = $2
			, goods_receipt_id = $3
			, purchase_order_item_id = $4
			, product_id = $5
			, product_variant_id = $6
			, quantity_received = $7
			, quantity_accepted = $8
			, quantity_rejected = $9
			, rejection_reason = $10
			, deleted_at = $12
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.GoodsReceiptId,
		entity.PurchaseOrderItemId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.QuantityReceived,
		entity.QuantityAccepted,
		entity.QuantityRejected,
		entity.RejectionReason,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update goods_receipt_items", zap.Error(err))
		return fmt.Errorf("failed to update goods_receipt_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("goods_receipt_items not found or already deleted")
	}

	r.logger.Info("updated goods_receipt_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a goods_receipt_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "goods_receipt_items", duration, nil)
	}()

	query := `
		UPDATE goods_receipt_items
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete goods_receipt_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete goods_receipt_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("goods_receipt_items not found or already deleted")
	}

	r.logger.Info("deleted goods_receipt_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves goods_receipt_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*GoodsReceiptItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "goods_receipt_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM goods_receipt_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count goods_receipt_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, goods_receipt_id
			, purchase_order_item_id
			, product_id
			, product_variant_id
			, quantity_received
			, quantity_accepted
			, quantity_rejected
			, rejection_reason
			, created_at
			, deleted_at
		FROM goods_receipt_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list goods_receipt_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list goods_receipt_items: %w", err)
	}
	defer rows.Close()

	var entities []*GoodsReceiptItems
	for rows.Next() {
		var entity GoodsReceiptItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.GoodsReceiptId,
			&entity.PurchaseOrderItemId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.QuantityReceived,
			&entity.QuantityAccepted,
			&entity.QuantityRejected,
			&entity.RejectionReason,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan goods_receipt_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


package goods_receipt

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

// Repository handles database operations for GoodsReceipts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new GoodsReceipts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// GoodsReceipts represents a goods_receipts entity
type GoodsReceipts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ReceiptNumber string `json:"receipt_number" db:"receipt_number"`
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id" db:"purchase_order_id"`
	SupplierId uuid.UUID `json:"supplier_id" db:"supplier_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	ReceiptDate time.Time `json:"receipt_date" db:"receipt_date"`
	ReceivedBy uuid.UUID `json:"received_by" db:"received_by"`
	Status *string `json:"status" db:"status"`
	Status *string `json:"status" db:"status"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new goods_receipts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *GoodsReceipts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "goods_receipts", duration, nil)
	}()

	query := `
		INSERT INTO goods_receipts (
			, organization_id
			, receipt_number
			, purchase_order_id
			, supplier_id
			, location_id
			, receipt_date
			, received_by
			, status
			, status
			, notes
			, created_by
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
			, $15
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ReceiptNumber,
		entity.PurchaseOrderId,
		entity.SupplierId,
		entity.LocationId,
		entity.ReceiptDate,
		entity.ReceivedBy,
		entity.Status,
		entity.Status,
		entity.Notes,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create goods_receipts", zap.Error(err))
		return fmt.Errorf("failed to create goods_receipts: %w", err)
	}

	r.logger.Info("created goods_receipts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a goods_receipts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*GoodsReceipts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "goods_receipts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, receipt_number
			, purchase_order_id
			, supplier_id
			, location_id
			, receipt_date
			, received_by
			, status
			, status
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM goods_receipts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity GoodsReceipts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ReceiptNumber,
		&entity.PurchaseOrderId,
		&entity.SupplierId,
		&entity.LocationId,
		&entity.ReceiptDate,
		&entity.ReceivedBy,
		&entity.Status,
		&entity.Status,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("goods_receipts not found")
	}

	if err != nil {
		r.logger.Error("failed to get goods_receipts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get goods_receipts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of goods_receipts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*GoodsReceipts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "goods_receipts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM goods_receipts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count goods_receipts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, receipt_number
			, purchase_order_id
			, supplier_id
			, location_id
			, receipt_date
			, received_by
			, status
			, status
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM goods_receipts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list goods_receipts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list goods_receipts: %w", err)
	}
	defer rows.Close()

	var entities []*GoodsReceipts
	for rows.Next() {
		var entity GoodsReceipts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReceiptNumber,
			&entity.PurchaseOrderId,
			&entity.SupplierId,
			&entity.LocationId,
			&entity.ReceiptDate,
			&entity.ReceivedBy,
			&entity.Status,
			&entity.Status,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan goods_receipts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating goods_receipts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing goods_receipts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *GoodsReceipts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "goods_receipts", duration, nil)
	}()

	query := `
		UPDATE goods_receipts
		SET
			, organization_id = $2
			, receipt_number = $3
			, purchase_order_id = $4
			, supplier_id = $5
			, location_id = $6
			, receipt_date = $7
			, received_by = $8
			, status = $9
			, status = $10
			, notes = $11
			, created_by = $12
			, updated_at = $14
			, deleted_at = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ReceiptNumber,
		entity.PurchaseOrderId,
		entity.SupplierId,
		entity.LocationId,
		entity.ReceiptDate,
		entity.ReceivedBy,
		entity.Status,
		entity.Status,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update goods_receipts", zap.Error(err))
		return fmt.Errorf("failed to update goods_receipts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("goods_receipts not found or already deleted")
	}

	r.logger.Info("updated goods_receipts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a goods_receipts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "goods_receipts", duration, nil)
	}()

	query := `
		UPDATE goods_receipts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete goods_receipts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete goods_receipts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("goods_receipts not found or already deleted")
	}

	r.logger.Info("deleted goods_receipts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves goods_receipts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*GoodsReceipts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "goods_receipts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM goods_receipts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count goods_receipts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, receipt_number
			, purchase_order_id
			, supplier_id
			, location_id
			, receipt_date
			, received_by
			, status
			, status
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM goods_receipts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list goods_receipts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list goods_receipts: %w", err)
	}
	defer rows.Close()

	var entities []*GoodsReceipts
	for rows.Next() {
		var entity GoodsReceipts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReceiptNumber,
			&entity.PurchaseOrderId,
			&entity.SupplierId,
			&entity.LocationId,
			&entity.ReceiptDate,
			&entity.ReceivedBy,
			&entity.Status,
			&entity.Status,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan goods_receipts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


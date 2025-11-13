package sale_return_item

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

// Repository handles database operations for SaleReturnItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new SaleReturnItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// SaleReturnItems represents a sale_return_items entity
type SaleReturnItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SaleReturnId uuid.UUID `json:"sale_return_id" db:"sale_return_id"`
	OriginalSaleItemId *uuid.UUID `json:"original_sale_item_id" db:"original_sale_item_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	Quantity float64 `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
	Subtotal float64 `json:"subtotal" db:"subtotal"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	DiscountAmount *float64 `json:"discount_amount" db:"discount_amount"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	ReturnReasonId *uuid.UUID `json:"return_reason_id" db:"return_reason_id"`
	ReturnReasonNotes *string `json:"return_reason_notes" db:"return_reason_notes"`
	ItemCondition *string `json:"item_condition" db:"item_condition"`
	ItemCondition *string `json:"item_condition" db:"item_condition"`
	IsRestockable *bool `json:"is_restockable" db:"is_restockable"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new sale_return_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *SaleReturnItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "sale_return_items", duration, nil)
	}()

	query := `
		INSERT INTO sale_return_items (
			, organization_id
			, sale_return_id
			, original_sale_item_id
			, product_id
			, product_variant_id
			, quantity
			, unit_price
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, return_reason_id
			, return_reason_notes
			, item_condition
			, item_condition
			, is_restockable
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
			, $19
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SaleReturnId,
		entity.OriginalSaleItemId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.Quantity,
		entity.UnitPrice,
		entity.Subtotal,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.TotalAmount,
		entity.ReturnReasonId,
		entity.ReturnReasonNotes,
		entity.ItemCondition,
		entity.ItemCondition,
		entity.IsRestockable,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create sale_return_items", zap.Error(err))
		return fmt.Errorf("failed to create sale_return_items: %w", err)
	}

	r.logger.Info("created sale_return_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a sale_return_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*SaleReturnItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_return_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, sale_return_id
			, original_sale_item_id
			, product_id
			, product_variant_id
			, quantity
			, unit_price
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, return_reason_id
			, return_reason_notes
			, item_condition
			, item_condition
			, is_restockable
			, created_at
			, deleted_at
		FROM sale_return_items
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity SaleReturnItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SaleReturnId,
		&entity.OriginalSaleItemId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.Quantity,
		&entity.UnitPrice,
		&entity.Subtotal,
		&entity.TaxAmount,
		&entity.DiscountAmount,
		&entity.TotalAmount,
		&entity.ReturnReasonId,
		&entity.ReturnReasonNotes,
		&entity.ItemCondition,
		&entity.ItemCondition,
		&entity.IsRestockable,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("sale_return_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get sale_return_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get sale_return_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of sale_return_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*SaleReturnItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_return_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sale_return_items
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sale_return_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_return_id
			, original_sale_item_id
			, product_id
			, product_variant_id
			, quantity
			, unit_price
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, return_reason_id
			, return_reason_notes
			, item_condition
			, item_condition
			, is_restockable
			, created_at
			, deleted_at
		FROM sale_return_items
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sale_return_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sale_return_items: %w", err)
	}
	defer rows.Close()

	var entities []*SaleReturnItems
	for rows.Next() {
		var entity SaleReturnItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleReturnId,
			&entity.OriginalSaleItemId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.Quantity,
			&entity.UnitPrice,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.TotalAmount,
			&entity.ReturnReasonId,
			&entity.ReturnReasonNotes,
			&entity.ItemCondition,
			&entity.ItemCondition,
			&entity.IsRestockable,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sale_return_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating sale_return_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing sale_return_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *SaleReturnItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sale_return_items", duration, nil)
	}()

	query := `
		UPDATE sale_return_items
		SET
			, organization_id = $2
			, sale_return_id = $3
			, original_sale_item_id = $4
			, product_id = $5
			, product_variant_id = $6
			, quantity = $7
			, unit_price = $8
			, subtotal = $9
			, tax_amount = $10
			, discount_amount = $11
			, total_amount = $12
			, return_reason_id = $13
			, return_reason_notes = $14
			, item_condition = $15
			, item_condition = $16
			, is_restockable = $17
			, deleted_at = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SaleReturnId,
		entity.OriginalSaleItemId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.Quantity,
		entity.UnitPrice,
		entity.Subtotal,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.TotalAmount,
		entity.ReturnReasonId,
		entity.ReturnReasonNotes,
		entity.ItemCondition,
		entity.ItemCondition,
		entity.IsRestockable,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update sale_return_items", zap.Error(err))
		return fmt.Errorf("failed to update sale_return_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sale_return_items not found or already deleted")
	}

	r.logger.Info("updated sale_return_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a sale_return_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sale_return_items", duration, nil)
	}()

	query := `
		UPDATE sale_return_items
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete sale_return_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete sale_return_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sale_return_items not found or already deleted")
	}

	r.logger.Info("deleted sale_return_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves sale_return_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*SaleReturnItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_return_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sale_return_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sale_return_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_return_id
			, original_sale_item_id
			, product_id
			, product_variant_id
			, quantity
			, unit_price
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, return_reason_id
			, return_reason_notes
			, item_condition
			, item_condition
			, is_restockable
			, created_at
			, deleted_at
		FROM sale_return_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sale_return_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sale_return_items: %w", err)
	}
	defer rows.Close()

	var entities []*SaleReturnItems
	for rows.Next() {
		var entity SaleReturnItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleReturnId,
			&entity.OriginalSaleItemId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.Quantity,
			&entity.UnitPrice,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.TotalAmount,
			&entity.ReturnReasonId,
			&entity.ReturnReasonNotes,
			&entity.ItemCondition,
			&entity.ItemCondition,
			&entity.IsRestockable,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sale_return_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


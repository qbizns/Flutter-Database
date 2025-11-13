package price_list_item

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

// Repository handles database operations for PriceListItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PriceListItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PriceListItems represents a price_list_items entity
type PriceListItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	PriceListId uuid.UUID `json:"price_list_id" db:"price_list_id"`
	ProductId *uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	CategoryId *uuid.UUID `json:"category_id" db:"category_id"`
	OverridePrice *float64 `json:"override_price" db:"override_price"`
	DiscountPercentage *float64 `json:"discount_percentage" db:"discount_percentage"`
	MarkupPercentage *float64 `json:"markup_percentage" db:"markup_percentage"`
	MinPrice *float64 `json:"min_price" db:"min_price"`
	MaxPrice *float64 `json:"max_price" db:"max_price"`
	MinQuantity *float64 `json:"min_quantity" db:"min_quantity"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	ProductId string `json:"product_id" db:"product_id"`
}

// Create inserts a new price_list_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PriceListItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "price_list_items", duration, nil)
	}()

	query := `
		INSERT INTO price_list_items (
			, price_list_id
			, product_id
			, product_variant_id
			, category_id
			, override_price
			, discount_percentage
			, markup_percentage
			, min_price
			, max_price
			, min_quantity
			, deleted_at
			, product_id
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
			, $14
			, $15
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.PriceListId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.CategoryId,
		entity.OverridePrice,
		entity.DiscountPercentage,
		entity.MarkupPercentage,
		entity.MinPrice,
		entity.MaxPrice,
		entity.MinQuantity,
		entity.DeletedAt,
		entity.ProductId,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create price_list_items", zap.Error(err))
		return fmt.Errorf("failed to create price_list_items: %w", err)
	}

	r.logger.Info("created price_list_items",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a price_list_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PriceListItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "price_list_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, price_list_id
			, product_id
			, product_variant_id
			, category_id
			, override_price
			, discount_percentage
			, markup_percentage
			, min_price
			, max_price
			, min_quantity
			, created_at
			, updated_at
			, deleted_at
			, product_id
		FROM price_list_items
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PriceListItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.PriceListId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.CategoryId,
		&entity.OverridePrice,
		&entity.DiscountPercentage,
		&entity.MarkupPercentage,
		&entity.MinPrice,
		&entity.MaxPrice,
		&entity.MinQuantity,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.ProductId,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("price_list_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get price_list_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get price_list_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of price_list_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PriceListItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "price_list_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM price_list_items
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count price_list_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, price_list_id
			, product_id
			, product_variant_id
			, category_id
			, override_price
			, discount_percentage
			, markup_percentage
			, min_price
			, max_price
			, min_quantity
			, created_at
			, updated_at
			, deleted_at
			, product_id
		FROM price_list_items
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list price_list_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list price_list_items: %w", err)
	}
	defer rows.Close()

	var entities []*PriceListItems
	for rows.Next() {
		var entity PriceListItems
		err := rows.Scan(
			&entity.Id,
			&entity.PriceListId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.CategoryId,
			&entity.OverridePrice,
			&entity.DiscountPercentage,
			&entity.MarkupPercentage,
			&entity.MinPrice,
			&entity.MaxPrice,
			&entity.MinQuantity,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.ProductId,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan price_list_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating price_list_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing price_list_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PriceListItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "price_list_items", duration, nil)
	}()

	query := `
		UPDATE price_list_items
		SET
			, price_list_id = $2
			, product_id = $3
			, product_variant_id = $4
			, category_id = $5
			, override_price = $6
			, discount_percentage = $7
			, markup_percentage = $8
			, min_price = $9
			, max_price = $10
			, min_quantity = $11
			, updated_at = $13
			, deleted_at = $14
			, product_id = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.PriceListId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.CategoryId,
		entity.OverridePrice,
		entity.DiscountPercentage,
		entity.MarkupPercentage,
		entity.MinPrice,
		entity.MaxPrice,
		entity.MinQuantity,
		time.Now(),
		entity.DeletedAt,
		entity.ProductId,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update price_list_items", zap.Error(err))
		return fmt.Errorf("failed to update price_list_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("price_list_items not found or already deleted")
	}

	r.logger.Info("updated price_list_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a price_list_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "price_list_items", duration, nil)
	}()

	query := `
		UPDATE price_list_items
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete price_list_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete price_list_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("price_list_items not found or already deleted")
	}

	r.logger.Info("deleted price_list_items", zap.String("id", id.String()))
	return nil
}




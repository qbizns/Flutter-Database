package sale_item

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

// Repository handles database operations for SaleItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new SaleItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// SaleItems represents a sale_items entity
type SaleItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	SaleId uuid.UUID `json:"sale_id" db:"sale_id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ProductId *uuid.UUID `json:"product_id" db:"product_id"`
	ProductName string `json:"product_name" db:"product_name"`
	ProductSku *string `json:"product_sku" db:"product_sku"`
	Quantity float64 `json:"quantity" db:"quantity"`
	Unit *string `json:"unit" db:"unit"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
	CostPrice *float64 `json:"cost_price" db:"cost_price"`
	Subtotal float64 `json:"subtotal" db:"subtotal"`
	TaxRate *float64 `json:"tax_rate" db:"tax_rate"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	DiscountAmount *float64 `json:"discount_amount" db:"discount_amount"`
	Total float64 `json:"total" db:"total"`
	DiscountType *string `json:"discount_type" db:"discount_type"`
	DiscountValue *float64 `json:"discount_value" db:"discount_value"`
	Notes *string `json:"notes" db:"notes"`
	CustomFields json.RawMessage `json:"custom_fields" db:"custom_fields"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new sale_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *SaleItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "sale_items", duration, nil)
	}()

	query := `
		INSERT INTO sale_items (
			, sale_id
			, organization_id
			, product_id
			, product_name
			, product_sku
			, quantity
			, unit
			, unit_price
			, cost_price
			, subtotal
			, tax_rate
			, tax_amount
			, discount_amount
			, total
			, discount_type
			, discount_value
			, notes
			, custom_fields
			, metadata
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.SaleId,
		entity.OrganizationId,
		entity.ProductId,
		entity.ProductName,
		entity.ProductSku,
		entity.Quantity,
		entity.Unit,
		entity.UnitPrice,
		entity.CostPrice,
		entity.Subtotal,
		entity.TaxRate,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.Total,
		entity.DiscountType,
		entity.DiscountValue,
		entity.Notes,
		entity.CustomFields,
		entity.Metadata,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create sale_items", zap.Error(err))
		return fmt.Errorf("failed to create sale_items: %w", err)
	}

	r.logger.Info("created sale_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a sale_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*SaleItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, sale_id
			, organization_id
			, product_id
			, product_name
			, product_sku
			, quantity
			, unit
			, unit_price
			, cost_price
			, subtotal
			, tax_rate
			, tax_amount
			, discount_amount
			, total
			, discount_type
			, discount_value
			, notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
		FROM sale_items
		WHERE id = $1
		
	`

	var entity SaleItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.SaleId,
		&entity.OrganizationId,
		&entity.ProductId,
		&entity.ProductName,
		&entity.ProductSku,
		&entity.Quantity,
		&entity.Unit,
		&entity.UnitPrice,
		&entity.CostPrice,
		&entity.Subtotal,
		&entity.TaxRate,
		&entity.TaxAmount,
		&entity.DiscountAmount,
		&entity.Total,
		&entity.DiscountType,
		&entity.DiscountValue,
		&entity.Notes,
		&entity.CustomFields,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("sale_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get sale_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get sale_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of sale_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*SaleItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sale_items
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sale_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, sale_id
			, organization_id
			, product_id
			, product_name
			, product_sku
			, quantity
			, unit
			, unit_price
			, cost_price
			, subtotal
			, tax_rate
			, tax_amount
			, discount_amount
			, total
			, discount_type
			, discount_value
			, notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
		FROM sale_items
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sale_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sale_items: %w", err)
	}
	defer rows.Close()

	var entities []*SaleItems
	for rows.Next() {
		var entity SaleItems
		err := rows.Scan(
			&entity.Id,
			&entity.SaleId,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ProductName,
			&entity.ProductSku,
			&entity.Quantity,
			&entity.Unit,
			&entity.UnitPrice,
			&entity.CostPrice,
			&entity.Subtotal,
			&entity.TaxRate,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.Total,
			&entity.DiscountType,
			&entity.DiscountValue,
			&entity.Notes,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sale_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating sale_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing sale_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *SaleItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sale_items", duration, nil)
	}()

	query := `
		UPDATE sale_items
		SET
			, sale_id = $2
			, organization_id = $3
			, product_id = $4
			, product_name = $5
			, product_sku = $6
			, quantity = $7
			, unit = $8
			, unit_price = $9
			, cost_price = $10
			, subtotal = $11
			, tax_rate = $12
			, tax_amount = $13
			, discount_amount = $14
			, total = $15
			, discount_type = $16
			, discount_value = $17
			, notes = $18
			, custom_fields = $19
			, metadata = $20
			, updated_at = $22
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $23
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.SaleId,
		entity.OrganizationId,
		entity.ProductId,
		entity.ProductName,
		entity.ProductSku,
		entity.Quantity,
		entity.Unit,
		entity.UnitPrice,
		entity.CostPrice,
		entity.Subtotal,
		entity.TaxRate,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.Total,
		entity.DiscountType,
		entity.DiscountValue,
		entity.Notes,
		entity.CustomFields,
		entity.Metadata,
		time.Now(),
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update sale_items", zap.Error(err))
		return fmt.Errorf("failed to update sale_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sale_items not found or already deleted")
	}

	r.logger.Info("updated sale_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a sale_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "sale_items", duration, nil)
	}()

	query := `DELETE FROM sale_items WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete sale_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete sale_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sale_items not found")
	}

	r.logger.Info("deleted sale_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves sale_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*SaleItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sale_items
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sale_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, sale_id
			, organization_id
			, product_id
			, product_name
			, product_sku
			, quantity
			, unit
			, unit_price
			, cost_price
			, subtotal
			, tax_rate
			, tax_amount
			, discount_amount
			, total
			, discount_type
			, discount_value
			, notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
		FROM sale_items
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sale_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sale_items: %w", err)
	}
	defer rows.Close()

	var entities []*SaleItems
	for rows.Next() {
		var entity SaleItems
		err := rows.Scan(
			&entity.Id,
			&entity.SaleId,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ProductName,
			&entity.ProductSku,
			&entity.Quantity,
			&entity.Unit,
			&entity.UnitPrice,
			&entity.CostPrice,
			&entity.Subtotal,
			&entity.TaxRate,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.Total,
			&entity.DiscountType,
			&entity.DiscountValue,
			&entity.Notes,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sale_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


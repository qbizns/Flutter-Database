package cycle_count_item

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

// Repository handles database operations for CycleCountItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CycleCountItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CycleCountItems represents a cycle_count_items entity
type CycleCountItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CycleCountId uuid.UUID `json:"cycle_count_id" db:"cycle_count_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	ProductName string `json:"product_name" db:"product_name"`
	ProductSku *string `json:"product_sku" db:"product_sku"`
	SystemQuantity float64 `json:"system_quantity" db:"system_quantity"`
	CountedQuantity *float64 `json:"counted_quantity" db:"counted_quantity"`
	VarianceQuantity *float64 `json:"variance_quantity" db:"variance_quantity"`
	VariancePercentage *float64 `json:"variance_percentage" db:"variance_percentage"`
	UnitCost *float64 `json:"unit_cost" db:"unit_cost"`
	VarianceValue *float64 `json:"variance_value" db:"variance_value"`
	Status *string `json:"status" db:"status"`
	RecountRequired *bool `json:"recount_required" db:"recount_required"`
	RecountQuantity *float64 `json:"recount_quantity" db:"recount_quantity"`
	RecountReason *string `json:"recount_reason" db:"recount_reason"`
	AdjustmentApplied *bool `json:"adjustment_applied" db:"adjustment_applied"`
	AdjustmentDate *time.Time `json:"adjustment_date" db:"adjustment_date"`
	AdjustmentReason *string `json:"adjustment_reason" db:"adjustment_reason"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CountedAt *time.Time `json:"counted_at" db:"counted_at"`
	CountedBy *uuid.UUID `json:"counted_by" db:"counted_by"`
}

// Create inserts a new cycle_count_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CycleCountItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "cycle_count_items", duration, nil)
	}()

	query := `
		INSERT INTO cycle_count_items (
			, organization_id
			, cycle_count_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, system_quantity
			, counted_quantity
			, variance_quantity
			, variance_percentage
			, unit_cost
			, variance_value
			, status
			, recount_required
			, recount_quantity
			, recount_reason
			, adjustment_applied
			, adjustment_date
			, adjustment_reason
			, notes
			, metadata
			, counted_at
			, counted_by
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
			, $25
			, $26
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CycleCountId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.ProductName,
		entity.ProductSku,
		entity.SystemQuantity,
		entity.CountedQuantity,
		entity.VarianceQuantity,
		entity.VariancePercentage,
		entity.UnitCost,
		entity.VarianceValue,
		entity.Status,
		entity.RecountRequired,
		entity.RecountQuantity,
		entity.RecountReason,
		entity.AdjustmentApplied,
		entity.AdjustmentDate,
		entity.AdjustmentReason,
		entity.Notes,
		entity.Metadata,
		entity.CountedAt,
		entity.CountedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create cycle_count_items", zap.Error(err))
		return fmt.Errorf("failed to create cycle_count_items: %w", err)
	}

	r.logger.Info("created cycle_count_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a cycle_count_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CycleCountItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cycle_count_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, cycle_count_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, system_quantity
			, counted_quantity
			, variance_quantity
			, variance_percentage
			, unit_cost
			, variance_value
			, status
			, recount_required
			, recount_quantity
			, recount_reason
			, adjustment_applied
			, adjustment_date
			, adjustment_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, counted_at
			, counted_by
		FROM cycle_count_items
		WHERE id = $1
		
	`

	var entity CycleCountItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CycleCountId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.ProductName,
		&entity.ProductSku,
		&entity.SystemQuantity,
		&entity.CountedQuantity,
		&entity.VarianceQuantity,
		&entity.VariancePercentage,
		&entity.UnitCost,
		&entity.VarianceValue,
		&entity.Status,
		&entity.RecountRequired,
		&entity.RecountQuantity,
		&entity.RecountReason,
		&entity.AdjustmentApplied,
		&entity.AdjustmentDate,
		&entity.AdjustmentReason,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CountedAt,
		&entity.CountedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("cycle_count_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get cycle_count_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get cycle_count_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of cycle_count_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CycleCountItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cycle_count_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cycle_count_items
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cycle_count_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, cycle_count_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, system_quantity
			, counted_quantity
			, variance_quantity
			, variance_percentage
			, unit_cost
			, variance_value
			, status
			, recount_required
			, recount_quantity
			, recount_reason
			, adjustment_applied
			, adjustment_date
			, adjustment_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, counted_at
			, counted_by
		FROM cycle_count_items
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cycle_count_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cycle_count_items: %w", err)
	}
	defer rows.Close()

	var entities []*CycleCountItems
	for rows.Next() {
		var entity CycleCountItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CycleCountId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.ProductName,
			&entity.ProductSku,
			&entity.SystemQuantity,
			&entity.CountedQuantity,
			&entity.VarianceQuantity,
			&entity.VariancePercentage,
			&entity.UnitCost,
			&entity.VarianceValue,
			&entity.Status,
			&entity.RecountRequired,
			&entity.RecountQuantity,
			&entity.RecountReason,
			&entity.AdjustmentApplied,
			&entity.AdjustmentDate,
			&entity.AdjustmentReason,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CountedAt,
			&entity.CountedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cycle_count_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating cycle_count_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing cycle_count_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CycleCountItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cycle_count_items", duration, nil)
	}()

	query := `
		UPDATE cycle_count_items
		SET
			, organization_id = $2
			, cycle_count_id = $3
			, product_id = $4
			, product_variant_id = $5
			, product_name = $6
			, product_sku = $7
			, system_quantity = $8
			, counted_quantity = $9
			, variance_quantity = $10
			, variance_percentage = $11
			, unit_cost = $12
			, variance_value = $13
			, status = $14
			, recount_required = $15
			, recount_quantity = $16
			, recount_reason = $17
			, adjustment_applied = $18
			, adjustment_date = $19
			, adjustment_reason = $20
			, notes = $21
			, metadata = $22
			, updated_at = $24
			, counted_at = $25
			, counted_by = $26
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $27
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CycleCountId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.ProductName,
		entity.ProductSku,
		entity.SystemQuantity,
		entity.CountedQuantity,
		entity.VarianceQuantity,
		entity.VariancePercentage,
		entity.UnitCost,
		entity.VarianceValue,
		entity.Status,
		entity.RecountRequired,
		entity.RecountQuantity,
		entity.RecountReason,
		entity.AdjustmentApplied,
		entity.AdjustmentDate,
		entity.AdjustmentReason,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CountedAt,
		entity.CountedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update cycle_count_items", zap.Error(err))
		return fmt.Errorf("failed to update cycle_count_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cycle_count_items not found or already deleted")
	}

	r.logger.Info("updated cycle_count_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a cycle_count_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "cycle_count_items", duration, nil)
	}()

	query := `DELETE FROM cycle_count_items WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete cycle_count_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete cycle_count_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cycle_count_items not found")
	}

	r.logger.Info("deleted cycle_count_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves cycle_count_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CycleCountItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cycle_count_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cycle_count_items
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cycle_count_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, cycle_count_id
			, product_id
			, product_variant_id
			, product_name
			, product_sku
			, system_quantity
			, counted_quantity
			, variance_quantity
			, variance_percentage
			, unit_cost
			, variance_value
			, status
			, recount_required
			, recount_quantity
			, recount_reason
			, adjustment_applied
			, adjustment_date
			, adjustment_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, counted_at
			, counted_by
		FROM cycle_count_items
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cycle_count_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cycle_count_items: %w", err)
	}
	defer rows.Close()

	var entities []*CycleCountItems
	for rows.Next() {
		var entity CycleCountItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CycleCountId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.ProductName,
			&entity.ProductSku,
			&entity.SystemQuantity,
			&entity.CountedQuantity,
			&entity.VarianceQuantity,
			&entity.VariancePercentage,
			&entity.UnitCost,
			&entity.VarianceValue,
			&entity.Status,
			&entity.RecountRequired,
			&entity.RecountQuantity,
			&entity.RecountReason,
			&entity.AdjustmentApplied,
			&entity.AdjustmentDate,
			&entity.AdjustmentReason,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CountedAt,
			&entity.CountedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cycle_count_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


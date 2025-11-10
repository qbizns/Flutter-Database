package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/inventory"
)

// CycleCountRepository implements inventory cycle count and adjustment reason repositories
type CycleCountRepository struct {
	db *DB
}

// NewCycleCountRepository creates a new cycle count repository
func NewCycleCountRepository(db *DB) *CycleCountRepository {
	return &CycleCountRepository{db: db}
}

// ============================================================================
// Cycle Count Operations
// ============================================================================

// ListCycleCounts retrieves cycle counts with filters
func (r *CycleCountRepository) ListCycleCounts(ctx context.Context, orgID uuid.UUID, filters inventory.CycleCountFilters) ([]inventory.CycleCount, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, count_number, count_date, count_type, status,
		       category_id, include_zero_stock, total_items_planned, total_items_counted, items_with_variance,
		       total_variance_value, scheduled_date, started_at, completed_at, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by, counted_by, approved_by
		FROM cycle_counts
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.CountType != nil {
		argCount++
		query += fmt.Sprintf(" AND count_type = $%d", argCount)
		args = append(args, *filters.CountType)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
	}

	if filters.CountNumberLike != nil {
		argCount++
		query += fmt.Sprintf(" AND count_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.CountNumberLike+"%")
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND count_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND count_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	if filters.HasVariance != nil && *filters.HasVariance {
		query += " AND items_with_variance > 0"
	}

	query += " ORDER BY count_date DESC, created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []inventory.CycleCount
	for rows.Next() {
		var c inventory.CycleCount
		var metadata interface{}
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.LocationID, &c.CountNumber, &c.CountDate, &c.CountType, &c.Status,
			&c.CategoryID, &c.IncludeZeroStock, &c.TotalItemsPlanned, &c.TotalItemsCounted, &c.ItemsWithVariance,
			&c.TotalVarianceValue, &c.ScheduledDate, &c.StartedAt, &c.CompletedAt, &c.Notes, &metadata,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &c.CreatedBy, &c.UpdatedBy, &c.CountedBy, &c.ApprovedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				c.Metadata = data
			}
		}
		counts = append(counts, c)
	}

	return counts, rows.Err()
}

// CountCycleCounts counts cycle counts
func (r *CycleCountRepository) CountCycleCounts(ctx context.Context, orgID uuid.UUID, filters inventory.CycleCountFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM cycle_counts WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.CountType != nil {
		argCount++
		query += fmt.Sprintf(" AND count_type = $%d", argCount)
		args = append(args, *filters.CountType)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
	}

	if filters.CountNumberLike != nil {
		argCount++
		query += fmt.Sprintf(" AND count_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.CountNumberLike+"%")
	}

	if filters.HasVariance != nil && *filters.HasVariance {
		query += " AND items_with_variance > 0"
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetCycleCount retrieves a cycle count by ID
func (r *CycleCountRepository) GetCycleCount(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.CycleCount, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, count_number, count_date, count_type, status,
		       category_id, include_zero_stock, total_items_planned, total_items_counted, items_with_variance,
		       total_variance_value, scheduled_date, started_at, completed_at, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by, counted_by, approved_by
		FROM cycle_counts
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var c inventory.CycleCount
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&c.ID, &c.OrganizationID, &c.LocationID, &c.CountNumber, &c.CountDate, &c.CountType, &c.Status,
		&c.CategoryID, &c.IncludeZeroStock, &c.TotalItemsPlanned, &c.TotalItemsCounted, &c.ItemsWithVariance,
		&c.TotalVarianceValue, &c.ScheduledDate, &c.StartedAt, &c.CompletedAt, &c.Notes, &metadata,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &c.CreatedBy, &c.UpdatedBy, &c.CountedBy, &c.ApprovedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			c.Metadata = data
		}
	}

	return &c, nil
}

// GetCycleCountByNumber retrieves a cycle count by count number
func (r *CycleCountRepository) GetCycleCountByNumber(ctx context.Context, orgID uuid.UUID, countNumber string) (*inventory.CycleCount, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, count_number, count_date, count_type, status,
		       category_id, include_zero_stock, total_items_planned, total_items_counted, items_with_variance,
		       total_variance_value, scheduled_date, started_at, completed_at, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by, counted_by, approved_by
		FROM cycle_counts
		WHERE organization_id = $1 AND count_number = $2 AND deleted_at IS NULL
	`

	var c inventory.CycleCount
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, countNumber).Scan(
		&c.ID, &c.OrganizationID, &c.LocationID, &c.CountNumber, &c.CountDate, &c.CountType, &c.Status,
		&c.CategoryID, &c.IncludeZeroStock, &c.TotalItemsPlanned, &c.TotalItemsCounted, &c.ItemsWithVariance,
		&c.TotalVarianceValue, &c.ScheduledDate, &c.StartedAt, &c.CompletedAt, &c.Notes, &metadata,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &c.CreatedBy, &c.UpdatedBy, &c.CountedBy, &c.ApprovedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			c.Metadata = data
		}
	}

	return &c, nil
}

// CreateCycleCount creates a new cycle count
func (r *CycleCountRepository) CreateCycleCount(ctx context.Context, count *inventory.CycleCount) error {
	if err := r.db.SetOrganizationContext(ctx, count.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO cycle_counts (
			id, organization_id, location_id, count_number, count_date, count_type, status,
			category_id, include_zero_stock, total_items_planned, total_items_counted, items_with_variance,
			total_variance_value, scheduled_date, started_at, completed_at, notes, metadata,
			created_at, updated_at, created_by, updated_by, counted_by, approved_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		count.ID, count.OrganizationID, count.LocationID, count.CountNumber, count.CountDate, count.CountType, count.Status,
		count.CategoryID, count.IncludeZeroStock, count.TotalItemsPlanned, count.TotalItemsCounted, count.ItemsWithVariance,
		count.TotalVarianceValue, count.ScheduledDate, count.StartedAt, count.CompletedAt, count.Notes, count.Metadata,
		count.CreatedAt, count.UpdatedAt, count.CreatedBy, count.UpdatedBy, count.CountedBy, count.ApprovedBy,
	)

	return err
}

// UpdateCycleCount updates a cycle count
func (r *CycleCountRepository) UpdateCycleCount(ctx context.Context, count *inventory.CycleCount) error {
	if err := r.db.SetOrganizationContext(ctx, count.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_counts SET
			location_id = $3, count_number = $4, count_date = $5, count_type = $6, status = $7,
			category_id = $8, include_zero_stock = $9, total_items_planned = $10, total_items_counted = $11,
			items_with_variance = $12, total_variance_value = $13, scheduled_date = $14, started_at = $15,
			completed_at = $16, notes = $17, metadata = $18, updated_at = $19, updated_by = $20,
			counted_by = $21, approved_by = $22
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		count.OrganizationID, count.ID,
		count.LocationID, count.CountNumber, count.CountDate, count.CountType, count.Status,
		count.CategoryID, count.IncludeZeroStock, count.TotalItemsPlanned, count.TotalItemsCounted,
		count.ItemsWithVariance, count.TotalVarianceValue, count.ScheduledDate, count.StartedAt,
		count.CompletedAt, count.Notes, count.Metadata, count.UpdatedAt, count.UpdatedBy,
		count.CountedBy, count.ApprovedBy,
	)

	return err
}

// DeleteCycleCount soft deletes a cycle count
func (r *CycleCountRepository) DeleteCycleCount(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_counts SET deleted_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// UpdateCycleCountStatus updates the status of a cycle count
func (r *CycleCountRepository) UpdateCycleCountStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status inventory.CycleCountStatus, countedBy *uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_counts SET
			status = $3, counted_by = $4, started_at = CASE WHEN status = 'planned' THEN CURRENT_TIMESTAMP ELSE started_at END,
			updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status, countedBy)
	return err
}

// CompleteCycleCount marks a cycle count as completed
func (r *CycleCountRepository) CompleteCycleCount(ctx context.Context, orgID uuid.UUID, id uuid.UUID, approvedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_counts SET
			status = $3, completed_at = CURRENT_TIMESTAMP, approved_by = $4, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, inventory.CycleCountStatusCompleted, approvedBy)
	return err
}

// RecalculateCycleCountStats recalculates statistics for a cycle count
func (r *CycleCountRepository) RecalculateCycleCountStats(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_counts SET
			total_items_counted = (
				SELECT COUNT(*) FROM cycle_count_items WHERE cycle_count_id = $2 AND status != 'pending'
			),
			items_with_variance = (
				SELECT COUNT(*) FROM cycle_count_items WHERE cycle_count_id = $2 AND variance_quantity IS NOT NULL AND variance_quantity != 0
			),
			total_variance_value = (
				SELECT COALESCE(SUM(variance_value), 0) FROM cycle_count_items WHERE cycle_count_id = $2 AND variance_value IS NOT NULL
			),
			updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// ============================================================================
// Cycle Count Item Operations
// ============================================================================

// ListCycleCountItems retrieves cycle count items with filters
func (r *CycleCountRepository) ListCycleCountItems(ctx context.Context, orgID uuid.UUID, filters inventory.CycleCountItemFilters) ([]inventory.CycleCountItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, cycle_count_id, product_id, product_variant_id,
		       product_name, product_sku, system_quantity, counted_quantity, variance_quantity, variance_percentage,
		       unit_cost, variance_value, status, recount_required, recount_quantity, recount_reason,
		       adjustment_applied, adjustment_date, adjustment_reason, notes, metadata,
		       created_at, updated_at, counted_at, counted_by
		FROM cycle_count_items
		WHERE organization_id = $1 AND cycle_count_id = $2
	`

	args := []interface{}{orgID, filters.CycleCountID}
	argCount := 2

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.HasVariance != nil && *filters.HasVariance {
		query += " AND variance_quantity IS NOT NULL AND variance_quantity != 0"
	}

	if filters.RecountRequired != nil && *filters.RecountRequired {
		query += " AND recount_required = true"
	}

	query += " ORDER BY product_name ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []inventory.CycleCountItem
	for rows.Next() {
		var item inventory.CycleCountItem
		var metadata interface{}
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.CycleCountID, &item.ProductID, &item.ProductVariantID,
			&item.ProductName, &item.ProductSKU, &item.SystemQuantity, &item.CountedQuantity, &item.VarianceQuantity, &item.VariancePercentage,
			&item.UnitCost, &item.VarianceValue, &item.Status, &item.RecountRequired, &item.RecountQuantity, &item.RecountReason,
			&item.AdjustmentApplied, &item.AdjustmentDate, &item.AdjustmentReason, &item.Notes, &metadata,
			&item.CreatedAt, &item.UpdatedAt, &item.CountedAt, &item.CountedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				item.Metadata = data
			}
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// CountCycleCountItems counts cycle count items
func (r *CycleCountRepository) CountCycleCountItems(ctx context.Context, orgID uuid.UUID, filters inventory.CycleCountItemFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM cycle_count_items WHERE organization_id = $1 AND cycle_count_id = $2"
	args := []interface{}{orgID, filters.CycleCountID}
	argCount := 2

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.HasVariance != nil && *filters.HasVariance {
		query += " AND variance_quantity IS NOT NULL AND variance_quantity != 0"
	}

	if filters.RecountRequired != nil && *filters.RecountRequired {
		query += " AND recount_required = true"
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetCycleCountItem retrieves a cycle count item
func (r *CycleCountRepository) GetCycleCountItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.CycleCountItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, cycle_count_id, product_id, product_variant_id,
		       product_name, product_sku, system_quantity, counted_quantity, variance_quantity, variance_percentage,
		       unit_cost, variance_value, status, recount_required, recount_quantity, recount_reason,
		       adjustment_applied, adjustment_date, adjustment_reason, notes, metadata,
		       created_at, updated_at, counted_at, counted_by
		FROM cycle_count_items
		WHERE organization_id = $1 AND id = $2
	`

	var item inventory.CycleCountItem
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&item.ID, &item.OrganizationID, &item.CycleCountID, &item.ProductID, &item.ProductVariantID,
		&item.ProductName, &item.ProductSKU, &item.SystemQuantity, &item.CountedQuantity, &item.VarianceQuantity, &item.VariancePercentage,
		&item.UnitCost, &item.VarianceValue, &item.Status, &item.RecountRequired, &item.RecountQuantity, &item.RecountReason,
		&item.AdjustmentApplied, &item.AdjustmentDate, &item.AdjustmentReason, &item.Notes, &metadata,
		&item.CreatedAt, &item.UpdatedAt, &item.CountedAt, &item.CountedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			item.Metadata = data
		}
	}

	return &item, nil
}

// CreateCycleCountItem creates a new cycle count item
func (r *CycleCountRepository) CreateCycleCountItem(ctx context.Context, item *inventory.CycleCountItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO cycle_count_items (
			id, organization_id, cycle_count_id, product_id, product_variant_id,
			product_name, product_sku, system_quantity, counted_quantity, variance_quantity, variance_percentage,
			unit_cost, variance_value, status, recount_required, recount_quantity, recount_reason,
			adjustment_applied, adjustment_date, adjustment_reason, notes, metadata,
			created_at, updated_at, counted_at, counted_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.OrganizationID, item.CycleCountID, item.ProductID, item.ProductVariantID,
		item.ProductName, item.ProductSKU, item.SystemQuantity, item.CountedQuantity, item.VarianceQuantity, item.VariancePercentage,
		item.UnitCost, item.VarianceValue, item.Status, item.RecountRequired, item.RecountQuantity, item.RecountReason,
		item.AdjustmentApplied, item.AdjustmentDate, item.AdjustmentReason, item.Notes, item.Metadata,
		item.CreatedAt, item.UpdatedAt, item.CountedAt, item.CountedBy,
	)

	return err
}

// UpdateCycleCountItem updates a cycle count item
func (r *CycleCountRepository) UpdateCycleCountItem(ctx context.Context, item *inventory.CycleCountItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_count_items SET
			product_variant_id = $3, product_name = $4, product_sku = $5,
			counted_quantity = $6, variance_quantity = $7, variance_percentage = $8,
			unit_cost = $9, variance_value = $10, status = $11, recount_required = $12,
			recount_quantity = $13, recount_reason = $14, adjustment_applied = $15,
			adjustment_date = $16, adjustment_reason = $17, notes = $18, metadata = $19,
			updated_at = $20, counted_at = $21, counted_by = $22
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.OrganizationID, item.ID,
		item.ProductVariantID, item.ProductName, item.ProductSKU,
		item.CountedQuantity, item.VarianceQuantity, item.VariancePercentage,
		item.UnitCost, item.VarianceValue, item.Status, item.RecountRequired,
		item.RecountQuantity, item.RecountReason, item.AdjustmentApplied,
		item.AdjustmentDate, item.AdjustmentReason, item.Notes, item.Metadata,
		item.UpdatedAt, item.CountedAt, item.CountedBy,
	)

	return err
}

// DeleteCycleCountItem deletes a cycle count item
func (r *CycleCountRepository) DeleteCycleCountItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "DELETE FROM cycle_count_items WHERE organization_id = $1 AND id = $2"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// UpdateCountedQuantity updates the counted quantity for a cycle count item
func (r *CycleCountRepository) UpdateCountedQuantity(ctx context.Context, orgID uuid.UUID, id uuid.UUID, countedQuantity float64, countedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_count_items SET
			counted_quantity = $3, counted_by = $4, counted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, countedQuantity, countedBy)
	return err
}

// GetCycleCountItemsByCount retrieves all items for a cycle count
func (r *CycleCountRepository) GetCycleCountItemsByCount(ctx context.Context, orgID uuid.UUID, countID uuid.UUID) ([]inventory.CycleCountItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, cycle_count_id, product_id, product_variant_id,
		       product_name, product_sku, system_quantity, counted_quantity, variance_quantity, variance_percentage,
		       unit_cost, variance_value, status, recount_required, recount_quantity, recount_reason,
		       adjustment_applied, adjustment_date, adjustment_reason, notes, metadata,
		       created_at, updated_at, counted_at, counted_by
		FROM cycle_count_items
		WHERE organization_id = $1 AND cycle_count_id = $2
		ORDER BY product_name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, countID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []inventory.CycleCountItem
	for rows.Next() {
		var item inventory.CycleCountItem
		var metadata interface{}
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.CycleCountID, &item.ProductID, &item.ProductVariantID,
			&item.ProductName, &item.ProductSKU, &item.SystemQuantity, &item.CountedQuantity, &item.VarianceQuantity, &item.VariancePercentage,
			&item.UnitCost, &item.VarianceValue, &item.Status, &item.RecountRequired, &item.RecountQuantity, &item.RecountReason,
			&item.AdjustmentApplied, &item.AdjustmentDate, &item.AdjustmentReason, &item.Notes, &metadata,
			&item.CreatedAt, &item.UpdatedAt, &item.CountedAt, &item.CountedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				item.Metadata = data
			}
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// GetItemsWithVariance retrieves all items with variances in a cycle count
func (r *CycleCountRepository) GetItemsWithVariance(ctx context.Context, orgID uuid.UUID, countID uuid.UUID) ([]inventory.CycleCountItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, cycle_count_id, product_id, product_variant_id,
		       product_name, product_sku, system_quantity, counted_quantity, variance_quantity, variance_percentage,
		       unit_cost, variance_value, status, recount_required, recount_quantity, recount_reason,
		       adjustment_applied, adjustment_date, adjustment_reason, notes, metadata,
		       created_at, updated_at, counted_at, counted_by
		FROM cycle_count_items
		WHERE organization_id = $1 AND cycle_count_id = $2 AND variance_quantity IS NOT NULL AND variance_quantity != 0
		ORDER BY variance_value DESC NULLS LAST
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, countID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []inventory.CycleCountItem
	for rows.Next() {
		var item inventory.CycleCountItem
		var metadata interface{}
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.CycleCountID, &item.ProductID, &item.ProductVariantID,
			&item.ProductName, &item.ProductSKU, &item.SystemQuantity, &item.CountedQuantity, &item.VarianceQuantity, &item.VariancePercentage,
			&item.UnitCost, &item.VarianceValue, &item.Status, &item.RecountRequired, &item.RecountQuantity, &item.RecountReason,
			&item.AdjustmentApplied, &item.AdjustmentDate, &item.AdjustmentReason, &item.Notes, &metadata,
			&item.CreatedAt, &item.UpdatedAt, &item.CountedAt, &item.CountedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				item.Metadata = data
			}
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// ApplyAdjustment marks a cycle count item adjustment as applied
func (r *CycleCountRepository) ApplyAdjustment(ctx context.Context, orgID uuid.UUID, itemID uuid.UUID, adjustmentReason string, approvedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cycle_count_items SET
			adjustment_applied = true, adjustment_date = CURRENT_TIMESTAMP, adjustment_reason = $3,
			status = 'adjusted', updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, itemID, adjustmentReason)
	return err
}

// ============================================================================
// Stock Adjustment Reason Operations
// ============================================================================

// ListAdjustmentReasons retrieves adjustment reasons with filters
func (r *CycleCountRepository) ListAdjustmentReasons(ctx context.Context, orgID *uuid.UUID, filters inventory.AdjustmentReasonFilters) ([]inventory.StockAdjustmentReason, error) {
	query := `
		SELECT id, organization_id, code, name, description, reason_type, is_system_reason,
		       is_active, requires_approval, requires_notes, sort_order, metadata,
		       created_at, updated_at, deleted_at
		FROM stock_adjustment_reasons
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND (organization_id = $%d OR organization_id IS NULL)", argCount)
		args = append(args, *orgID)
	} else {
		query += " AND organization_id IS NULL"
	}

	if filters.Code != nil {
		argCount++
		query += fmt.Sprintf(" AND code ILIKE $%d", argCount)
		args = append(args, "%"+*filters.Code+"%")
	}

	if filters.ReasonType != nil {
		argCount++
		query += fmt.Sprintf(" AND reason_type = $%d", argCount)
		args = append(args, *filters.ReasonType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.IsSystemReason != nil {
		argCount++
		query += fmt.Sprintf(" AND is_system_reason = $%d", argCount)
		args = append(args, *filters.IsSystemReason)
	}

	query += " ORDER BY sort_order ASC, name ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reasons []inventory.StockAdjustmentReason
	for rows.Next() {
		var reason inventory.StockAdjustmentReason
		var metadata interface{}
		err := rows.Scan(
			&reason.ID, &reason.OrganizationID, &reason.Code, &reason.Name, &reason.Description, &reason.ReasonType,
			&reason.IsSystemReason, &reason.IsActive, &reason.RequiresApproval, &reason.RequiresNotes,
			&reason.SortOrder, &metadata, &reason.CreatedAt, &reason.UpdatedAt, &reason.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				reason.Metadata = data
			}
		}
		reasons = append(reasons, reason)
	}

	return reasons, rows.Err()
}

// CountAdjustmentReasons counts adjustment reasons
func (r *CycleCountRepository) CountAdjustmentReasons(ctx context.Context, orgID *uuid.UUID, filters inventory.AdjustmentReasonFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM stock_adjustment_reasons WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND (organization_id = $%d OR organization_id IS NULL)", argCount)
		args = append(args, *orgID)
	} else {
		query += " AND organization_id IS NULL"
	}

	if filters.Code != nil {
		argCount++
		query += fmt.Sprintf(" AND code ILIKE $%d", argCount)
		args = append(args, "%"+*filters.Code+"%")
	}

	if filters.ReasonType != nil {
		argCount++
		query += fmt.Sprintf(" AND reason_type = $%d", argCount)
		args = append(args, *filters.ReasonType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(context.Background(), query, args...).Scan(&count)
	return count, err
}

// GetAdjustmentReason retrieves an adjustment reason by ID
func (r *CycleCountRepository) GetAdjustmentReason(ctx context.Context, id uuid.UUID) (*inventory.StockAdjustmentReason, error) {
	query := `
		SELECT id, organization_id, code, name, description, reason_type, is_system_reason,
		       is_active, requires_approval, requires_notes, sort_order, metadata,
		       created_at, updated_at, deleted_at
		FROM stock_adjustment_reasons
		WHERE id = $1 AND deleted_at IS NULL
	`

	var reason inventory.StockAdjustmentReason
	var metadata interface{}
	err := r.db.Pool.QueryRow(context.Background(), query, id).Scan(
		&reason.ID, &reason.OrganizationID, &reason.Code, &reason.Name, &reason.Description, &reason.ReasonType,
		&reason.IsSystemReason, &reason.IsActive, &reason.RequiresApproval, &reason.RequiresNotes,
		&reason.SortOrder, &metadata, &reason.CreatedAt, &reason.UpdatedAt, &reason.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			reason.Metadata = data
		}
	}

	return &reason, nil
}

// GetAdjustmentReasonByCode retrieves an adjustment reason by code
func (r *CycleCountRepository) GetAdjustmentReasonByCode(ctx context.Context, orgID *uuid.UUID, code string) (*inventory.StockAdjustmentReason, error) {
	query := `
		SELECT id, organization_id, code, name, description, reason_type, is_system_reason,
		       is_active, requires_approval, requires_notes, sort_order, metadata,
		       created_at, updated_at, deleted_at
		FROM stock_adjustment_reasons
		WHERE code = $1 AND deleted_at IS NULL
	`

	args := []interface{}{code}

	if orgID != nil {
		query += " AND (organization_id = $2 OR organization_id IS NULL)"
		args = append(args, *orgID)
	} else {
		query += " AND organization_id IS NULL"
	}

	var reason inventory.StockAdjustmentReason
	var metadata interface{}
	err := r.db.Pool.QueryRow(context.Background(), query, args...).Scan(
		&reason.ID, &reason.OrganizationID, &reason.Code, &reason.Name, &reason.Description, &reason.ReasonType,
		&reason.IsSystemReason, &reason.IsActive, &reason.RequiresApproval, &reason.RequiresNotes,
		&reason.SortOrder, &metadata, &reason.CreatedAt, &reason.UpdatedAt, &reason.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			reason.Metadata = data
		}
	}

	return &reason, nil
}

// CreateAdjustmentReason creates a new adjustment reason
func (r *CycleCountRepository) CreateAdjustmentReason(ctx context.Context, reason *inventory.StockAdjustmentReason) error {
	query := `
		INSERT INTO stock_adjustment_reasons (
			id, organization_id, code, name, description, reason_type, is_system_reason,
			is_active, requires_approval, requires_notes, sort_order, metadata,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	_, err := r.db.Pool.Exec(context.Background(), query,
		reason.ID, reason.OrganizationID, reason.Code, reason.Name, reason.Description, reason.ReasonType,
		reason.IsSystemReason, reason.IsActive, reason.RequiresApproval, reason.RequiresNotes, reason.SortOrder,
		reason.Metadata, reason.CreatedAt, reason.UpdatedAt,
	)

	return err
}

// UpdateAdjustmentReason updates an adjustment reason
func (r *CycleCountRepository) UpdateAdjustmentReason(ctx context.Context, reason *inventory.StockAdjustmentReason) error {
	query := `
		UPDATE stock_adjustment_reasons SET
			code = $2, name = $3, description = $4, reason_type = $5, is_system_reason = $6,
			is_active = $7, requires_approval = $8, requires_notes = $9, sort_order = $10,
			metadata = $11, updated_at = $12
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(context.Background(), query,
		reason.ID, reason.Code, reason.Name, reason.Description, reason.ReasonType, reason.IsSystemReason,
		reason.IsActive, reason.RequiresApproval, reason.RequiresNotes, reason.SortOrder,
		reason.Metadata, reason.UpdatedAt,
	)

	return err
}

// DeleteAdjustmentReason soft deletes an adjustment reason
func (r *CycleCountRepository) DeleteAdjustmentReason(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE stock_adjustment_reasons SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1"
	_, err := r.db.Pool.Exec(context.Background(), query, id)
	return err
}

// GetActiveAdjustmentReasons retrieves active adjustment reasons of a specific type
func (r *CycleCountRepository) GetActiveAdjustmentReasons(ctx context.Context, orgID *uuid.UUID, reasonType inventory.AdjustmentReasonType) ([]inventory.StockAdjustmentReason, error) {
	query := `
		SELECT id, organization_id, code, name, description, reason_type, is_system_reason,
		       is_active, requires_approval, requires_notes, sort_order, metadata,
		       created_at, updated_at, deleted_at
		FROM stock_adjustment_reasons
		WHERE is_active = true AND deleted_at IS NULL
		AND (reason_type = $1 OR reason_type = $2)
	`

	args := []interface{}{reasonType, inventory.AdjustmentReasonTypeBoth}

	if orgID != nil {
		query += " AND (organization_id = $3 OR organization_id IS NULL)"
		args = append(args, *orgID)
	} else {
		query += " AND organization_id IS NULL"
	}

	query += " ORDER BY sort_order ASC, name ASC"

	rows, err := r.db.Pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reasons []inventory.StockAdjustmentReason
	for rows.Next() {
		var reason inventory.StockAdjustmentReason
		var metadata interface{}
		err := rows.Scan(
			&reason.ID, &reason.OrganizationID, &reason.Code, &reason.Name, &reason.Description, &reason.ReasonType,
			&reason.IsSystemReason, &reason.IsActive, &reason.RequiresApproval, &reason.RequiresNotes,
			&reason.SortOrder, &metadata, &reason.CreatedAt, &reason.UpdatedAt, &reason.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				reason.Metadata = data
			}
		}
		reasons = append(reasons, reason)
	}

	return reasons, rows.Err()
}

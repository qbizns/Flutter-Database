package cycle_count

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

// Repository handles database operations for CycleCounts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CycleCounts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CycleCounts represents a cycle_counts entity
type CycleCounts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	CountNumber string `json:"count_number" db:"count_number"`
	CountDate time.Time `json:"count_date" db:"count_date"`
	CountType *string `json:"count_type" db:"count_type"`
	Status *string `json:"status" db:"status"`
	CategoryId *uuid.UUID `json:"category_id" db:"category_id"`
	IncludeZeroStock *bool `json:"include_zero_stock" db:"include_zero_stock"`
	TotalItemsPlanned *int64 `json:"total_items_planned" db:"total_items_planned"`
	TotalItemsCounted *int64 `json:"total_items_counted" db:"total_items_counted"`
	ItemsWithVariance *int64 `json:"items_with_variance" db:"items_with_variance"`
	TotalVarianceValue *float64 `json:"total_variance_value" db:"total_variance_value"`
	ScheduledDate *time.Time `json:"scheduled_date" db:"scheduled_date"`
	StartedAt *time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CountedBy *uuid.UUID `json:"counted_by" db:"counted_by"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	TotalItemsPlanned *string `json:"total_items_planned" db:"total_items_planned"`
	TotalItemsCounted *string `json:"total_items_counted" db:"total_items_counted"`
	ItemsWithVariance *string `json:"items_with_variance" db:"items_with_variance"`
}

// Create inserts a new cycle_counts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CycleCounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "cycle_counts", duration, nil)
	}()

	query := `
		INSERT INTO cycle_counts (
			, organization_id
			, location_id
			, count_number
			, count_date
			, count_type
			, status
			, category_id
			, include_zero_stock
			, total_items_planned
			, total_items_counted
			, items_with_variance
			, total_variance_value
			, scheduled_date
			, started_at
			, completed_at
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, counted_by
			, approved_by
			, total_items_planned
			, total_items_counted
			, items_with_variance
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
			, $28
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.CountNumber,
		entity.CountDate,
		entity.CountType,
		entity.Status,
		entity.CategoryId,
		entity.IncludeZeroStock,
		entity.TotalItemsPlanned,
		entity.TotalItemsCounted,
		entity.ItemsWithVariance,
		entity.TotalVarianceValue,
		entity.ScheduledDate,
		entity.StartedAt,
		entity.CompletedAt,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.CountedBy,
		entity.ApprovedBy,
		entity.TotalItemsPlanned,
		entity.TotalItemsCounted,
		entity.ItemsWithVariance,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create cycle_counts", zap.Error(err))
		return fmt.Errorf("failed to create cycle_counts: %w", err)
	}

	r.logger.Info("created cycle_counts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a cycle_counts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CycleCounts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cycle_counts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, count_number
			, count_date
			, count_type
			, category_id
			, include_zero_stock
			, total_items_planned
			, total_items_counted
			, items_with_variance
			, total_variance_value
			, scheduled_date
			, started_at
			, completed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, counted_by
			, approved_by
			, total_items_planned
			, total_items_counted
			, items_with_variance
		FROM cycle_counts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CycleCounts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.CountNumber,
		&entity.CountDate,
		&entity.CountType,
		&entity.Status,
		&entity.CategoryId,
		&entity.IncludeZeroStock,
		&entity.TotalItemsPlanned,
		&entity.TotalItemsCounted,
		&entity.ItemsWithVariance,
		&entity.TotalVarianceValue,
		&entity.ScheduledDate,
		&entity.StartedAt,
		&entity.CompletedAt,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CountedBy,
		&entity.ApprovedBy,
		&entity.TotalItemsPlanned,
		&entity.TotalItemsCounted,
		&entity.ItemsWithVariance,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("cycle_counts not found")
	}

	if err != nil {
		r.logger.Error("failed to get cycle_counts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get cycle_counts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of cycle_counts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CycleCounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cycle_counts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cycle_counts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cycle_counts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, count_number
			, count_date
			, count_type
			, category_id
			, include_zero_stock
			, total_items_planned
			, total_items_counted
			, items_with_variance
			, total_variance_value
			, scheduled_date
			, started_at
			, completed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, counted_by
			, approved_by
			, total_items_planned
			, total_items_counted
			, items_with_variance
		FROM cycle_counts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cycle_counts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cycle_counts: %w", err)
	}
	defer rows.Close()

	var entities []*CycleCounts
	for rows.Next() {
		var entity CycleCounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.CountNumber,
			&entity.CountDate,
			&entity.CountType,
			&entity.Status,
			&entity.CategoryId,
			&entity.IncludeZeroStock,
			&entity.TotalItemsPlanned,
			&entity.TotalItemsCounted,
			&entity.ItemsWithVariance,
			&entity.TotalVarianceValue,
			&entity.ScheduledDate,
			&entity.StartedAt,
			&entity.CompletedAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CountedBy,
			&entity.ApprovedBy,
			&entity.TotalItemsPlanned,
			&entity.TotalItemsCounted,
			&entity.ItemsWithVariance,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cycle_counts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating cycle_counts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing cycle_counts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CycleCounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cycle_counts", duration, nil)
	}()

	query := `
		UPDATE cycle_counts
		SET
			, organization_id = $2
			, location_id = $3
			, count_number = $4
			, count_date = $5
			, count_type = $6
			, status = $7
			, category_id = $8
			, include_zero_stock = $9
			, total_items_planned = $10
			, total_items_counted = $11
			, items_with_variance = $12
			, total_variance_value = $13
			, scheduled_date = $14
			, started_at = $15
			, completed_at = $16
			, notes = $17
			, metadata = $18
			, updated_at = $20
			, deleted_at = $21
			, created_by = $22
			, updated_by = $23
			, counted_by = $24
			, approved_by = $25
			, total_items_planned = $26
			, total_items_counted = $27
			, items_with_variance = $28
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $29
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.CountNumber,
		entity.CountDate,
		entity.CountType,
		entity.Status,
		entity.CategoryId,
		entity.IncludeZeroStock,
		entity.TotalItemsPlanned,
		entity.TotalItemsCounted,
		entity.ItemsWithVariance,
		entity.TotalVarianceValue,
		entity.ScheduledDate,
		entity.StartedAt,
		entity.CompletedAt,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.CountedBy,
		entity.ApprovedBy,
		entity.TotalItemsPlanned,
		entity.TotalItemsCounted,
		entity.ItemsWithVariance,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update cycle_counts", zap.Error(err))
		return fmt.Errorf("failed to update cycle_counts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cycle_counts not found or already deleted")
	}

	r.logger.Info("updated cycle_counts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a cycle_counts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cycle_counts", duration, nil)
	}()

	query := `
		UPDATE cycle_counts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete cycle_counts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete cycle_counts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cycle_counts not found or already deleted")
	}

	r.logger.Info("deleted cycle_counts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves cycle_counts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CycleCounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cycle_counts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cycle_counts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cycle_counts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, count_number
			, count_date
			, count_type
			, status
			, category_id
			, include_zero_stock
			, total_items_planned
			, total_items_counted
			, items_with_variance
			, total_variance_value
			, scheduled_date
			, started_at
			, completed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, counted_by
			, approved_by
			, total_items_planned
			, total_items_counted
			, items_with_variance
		FROM cycle_counts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cycle_counts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cycle_counts: %w", err)
	}
	defer rows.Close()

	var entities []*CycleCounts
	for rows.Next() {
		var entity CycleCounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.CountNumber,
			&entity.CountDate,
			&entity.CountType,
			&entity.Status,
			&entity.CategoryId,
			&entity.IncludeZeroStock,
			&entity.TotalItemsPlanned,
			&entity.TotalItemsCounted,
			&entity.ItemsWithVariance,
			&entity.TotalVarianceValue,
			&entity.ScheduledDate,
			&entity.StartedAt,
			&entity.CompletedAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CountedBy,
			&entity.ApprovedBy,
			&entity.TotalItemsPlanned,
			&entity.TotalItemsCounted,
			&entity.ItemsWithVariance,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cycle_counts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


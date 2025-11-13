package stock_adjustment_reason

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

// Repository handles database operations for StockAdjustmentReasons
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new StockAdjustmentReasons repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// StockAdjustmentReasons represents a stock_adjustment_reasons entity
type StockAdjustmentReasons struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	Code string `json:"code" db:"code"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	ReasonType string `json:"reason_type" db:"reason_type"`
	IsSystemReason *bool `json:"is_system_reason" db:"is_system_reason"`
	IsActive *bool `json:"is_active" db:"is_active"`
	RequiresApproval *bool `json:"requires_approval" db:"requires_approval"`
	RequiresNotes *bool `json:"requires_notes" db:"requires_notes"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new stock_adjustment_reasons record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *StockAdjustmentReasons) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "stock_adjustment_reasons", duration, nil)
	}()

	query := `
		INSERT INTO stock_adjustment_reasons (
			, organization_id
			, code
			, name
			, description
			, reason_type
			, is_system_reason
			, is_active
			, requires_approval
			, requires_notes
			, sort_order
			, metadata
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
		entity.Code,
		entity.Name,
		entity.Description,
		entity.ReasonType,
		entity.IsSystemReason,
		entity.IsActive,
		entity.RequiresApproval,
		entity.RequiresNotes,
		entity.SortOrder,
		entity.Metadata,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create stock_adjustment_reasons", zap.Error(err))
		return fmt.Errorf("failed to create stock_adjustment_reasons: %w", err)
	}

	r.logger.Info("created stock_adjustment_reasons",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a stock_adjustment_reasons by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*StockAdjustmentReasons, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "stock_adjustment_reasons", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, code
			, name
			, description
			, reason_type
			, is_system_reason
			, is_active
			, requires_approval
			, requires_notes
			, sort_order
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM stock_adjustment_reasons
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity StockAdjustmentReasons
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Code,
		&entity.Name,
		&entity.Description,
		&entity.ReasonType,
		&entity.IsSystemReason,
		&entity.IsActive,
		&entity.RequiresApproval,
		&entity.RequiresNotes,
		&entity.SortOrder,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("stock_adjustment_reasons not found")
	}

	if err != nil {
		r.logger.Error("failed to get stock_adjustment_reasons", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get stock_adjustment_reasons: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of stock_adjustment_reasons records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*StockAdjustmentReasons, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "stock_adjustment_reasons", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM stock_adjustment_reasons
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock_adjustment_reasons records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, code
			, name
			, description
			, reason_type
			, is_system_reason
			, is_active
			, requires_approval
			, requires_notes
			, sort_order
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM stock_adjustment_reasons
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list stock_adjustment_reasons", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list stock_adjustment_reasons: %w", err)
	}
	defer rows.Close()

	var entities []*StockAdjustmentReasons
	for rows.Next() {
		var entity StockAdjustmentReasons
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Code,
			&entity.Name,
			&entity.Description,
			&entity.ReasonType,
			&entity.IsSystemReason,
			&entity.IsActive,
			&entity.RequiresApproval,
			&entity.RequiresNotes,
			&entity.SortOrder,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stock_adjustment_reasons: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating stock_adjustment_reasons rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing stock_adjustment_reasons record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *StockAdjustmentReasons) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "stock_adjustment_reasons", duration, nil)
	}()

	query := `
		UPDATE stock_adjustment_reasons
		SET
			, organization_id = $2
			, code = $3
			, name = $4
			, description = $5
			, reason_type = $6
			, is_system_reason = $7
			, is_active = $8
			, requires_approval = $9
			, requires_notes = $10
			, sort_order = $11
			, metadata = $12
			, updated_at = $14
			, deleted_at = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Code,
		entity.Name,
		entity.Description,
		entity.ReasonType,
		entity.IsSystemReason,
		entity.IsActive,
		entity.RequiresApproval,
		entity.RequiresNotes,
		entity.SortOrder,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update stock_adjustment_reasons", zap.Error(err))
		return fmt.Errorf("failed to update stock_adjustment_reasons: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("stock_adjustment_reasons not found or already deleted")
	}

	r.logger.Info("updated stock_adjustment_reasons",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a stock_adjustment_reasons record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "stock_adjustment_reasons", duration, nil)
	}()

	query := `
		UPDATE stock_adjustment_reasons
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete stock_adjustment_reasons", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete stock_adjustment_reasons: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("stock_adjustment_reasons not found or already deleted")
	}

	r.logger.Info("deleted stock_adjustment_reasons", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves stock_adjustment_reasons records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*StockAdjustmentReasons, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "stock_adjustment_reasons", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM stock_adjustment_reasons
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock_adjustment_reasons records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, code
			, name
			, description
			, reason_type
			, is_system_reason
			, is_active
			, requires_approval
			, requires_notes
			, sort_order
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM stock_adjustment_reasons
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list stock_adjustment_reasons by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list stock_adjustment_reasons: %w", err)
	}
	defer rows.Close()

	var entities []*StockAdjustmentReasons
	for rows.Next() {
		var entity StockAdjustmentReasons
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Code,
			&entity.Name,
			&entity.Description,
			&entity.ReasonType,
			&entity.IsSystemReason,
			&entity.IsActive,
			&entity.RequiresApproval,
			&entity.RequiresNotes,
			&entity.SortOrder,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stock_adjustment_reasons: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


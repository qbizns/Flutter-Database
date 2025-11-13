package return_reason

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

// Repository handles database operations for ReturnReasons
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ReturnReasons repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ReturnReasons represents a return_reasons entity
type ReturnReasons struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	ReasonCode string `json:"reason_code" db:"reason_code"`
	ReasonName string `json:"reason_name" db:"reason_name"`
	RequiresApproval *bool `json:"requires_approval" db:"requires_approval"`
	AffectsInventory *bool `json:"affects_inventory" db:"affects_inventory"`
	IsRestockable *bool `json:"is_restockable" db:"is_restockable"`
	IsActive *bool `json:"is_active" db:"is_active"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new return_reasons record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ReturnReasons) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "return_reasons", duration, nil)
	}()

	query := `
		INSERT INTO return_reasons (
			, organization_id
			, reason_code
			, reason_name
			, requires_approval
			, affects_inventory
			, is_restockable
			, is_active
			, display_order
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
			, $11
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ReasonCode,
		entity.ReasonName,
		entity.RequiresApproval,
		entity.AffectsInventory,
		entity.IsRestockable,
		entity.IsActive,
		entity.DisplayOrder,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create return_reasons", zap.Error(err))
		return fmt.Errorf("failed to create return_reasons: %w", err)
	}

	r.logger.Info("created return_reasons",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a return_reasons by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ReturnReasons, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "return_reasons", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, reason_code
			, reason_name
			, requires_approval
			, affects_inventory
			, is_restockable
			, is_active
			, display_order
			, created_at
			, deleted_at
		FROM return_reasons
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ReturnReasons
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ReasonCode,
		&entity.ReasonName,
		&entity.RequiresApproval,
		&entity.AffectsInventory,
		&entity.IsRestockable,
		&entity.IsActive,
		&entity.DisplayOrder,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("return_reasons not found")
	}

	if err != nil {
		r.logger.Error("failed to get return_reasons", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get return_reasons: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of return_reasons records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ReturnReasons, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "return_reasons", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM return_reasons
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count return_reasons records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, reason_code
			, reason_name
			, requires_approval
			, affects_inventory
			, is_restockable
			, is_active
			, display_order
			, created_at
			, deleted_at
		FROM return_reasons
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list return_reasons", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list return_reasons: %w", err)
	}
	defer rows.Close()

	var entities []*ReturnReasons
	for rows.Next() {
		var entity ReturnReasons
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReasonCode,
			&entity.ReasonName,
			&entity.RequiresApproval,
			&entity.AffectsInventory,
			&entity.IsRestockable,
			&entity.IsActive,
			&entity.DisplayOrder,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan return_reasons: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating return_reasons rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing return_reasons record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ReturnReasons) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "return_reasons", duration, nil)
	}()

	query := `
		UPDATE return_reasons
		SET
			, organization_id = $2
			, reason_code = $3
			, reason_name = $4
			, requires_approval = $5
			, affects_inventory = $6
			, is_restockable = $7
			, is_active = $8
			, display_order = $9
			, deleted_at = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ReasonCode,
		entity.ReasonName,
		entity.RequiresApproval,
		entity.AffectsInventory,
		entity.IsRestockable,
		entity.IsActive,
		entity.DisplayOrder,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update return_reasons", zap.Error(err))
		return fmt.Errorf("failed to update return_reasons: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("return_reasons not found or already deleted")
	}

	r.logger.Info("updated return_reasons",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a return_reasons record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "return_reasons", duration, nil)
	}()

	query := `
		UPDATE return_reasons
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete return_reasons", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete return_reasons: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("return_reasons not found or already deleted")
	}

	r.logger.Info("deleted return_reasons", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves return_reasons records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ReturnReasons, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "return_reasons", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM return_reasons
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count return_reasons records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, reason_code
			, reason_name
			, requires_approval
			, affects_inventory
			, is_restockable
			, is_active
			, display_order
			, created_at
			, deleted_at
		FROM return_reasons
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list return_reasons by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list return_reasons: %w", err)
	}
	defer rows.Close()

	var entities []*ReturnReasons
	for rows.Next() {
		var entity ReturnReasons
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReasonCode,
			&entity.ReasonName,
			&entity.RequiresApproval,
			&entity.AffectsInventory,
			&entity.IsRestockable,
			&entity.IsActive,
			&entity.DisplayOrder,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan return_reasons: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


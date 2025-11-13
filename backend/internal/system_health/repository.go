package system_health

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

// Repository handles database operations for SystemHealth
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new SystemHealth repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// SystemHealth represents a system_health entity
type SystemHealth struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	Status *string `json:"status" db:"status"`
	LastCheckAt *time.Time `json:"last_check_at" db:"last_check_at"`
	LastSuccessAt *time.Time `json:"last_success_at" db:"last_success_at"`
	LastFailureAt *time.Time `json:"last_failure_at" db:"last_failure_at"`
	MetricValue *float64 `json:"metric_value" db:"metric_value"`
	MetricUnit *string `json:"metric_unit" db:"metric_unit"`
	ThresholdWarning *float64 `json:"threshold_warning" db:"threshold_warning"`
	ThresholdCritical *float64 `json:"threshold_critical" db:"threshold_critical"`
	Details json.RawMessage `json:"details" db:"details"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new system_health record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *SystemHealth) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "system_health", duration, nil)
	}()

	query := `
		INSERT INTO system_health (
			, organization_id
			, status
			, last_check_at
			, last_success_at
			, last_failure_at
			, metric_value
			, metric_unit
			, threshold_warning
			, threshold_critical
			, details
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Status,
		entity.LastCheckAt,
		entity.LastSuccessAt,
		entity.LastFailureAt,
		entity.MetricValue,
		entity.MetricUnit,
		entity.ThresholdWarning,
		entity.ThresholdCritical,
		entity.Details,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create system_health", zap.Error(err))
		return fmt.Errorf("failed to create system_health: %w", err)
	}

	r.logger.Info("created system_health",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a system_health by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*SystemHealth, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "system_health", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, last_check_at
			, last_success_at
			, last_failure_at
			, metric_value
			, metric_unit
			, threshold_warning
			, threshold_critical
			, details
			, created_at
			, updated_at
		FROM system_health
		WHERE id = $1
		
	`

	var entity SystemHealth
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Status,
		&entity.LastCheckAt,
		&entity.LastSuccessAt,
		&entity.LastFailureAt,
		&entity.MetricValue,
		&entity.MetricUnit,
		&entity.ThresholdWarning,
		&entity.ThresholdCritical,
		&entity.Details,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("system_health not found")
	}

	if err != nil {
		r.logger.Error("failed to get system_health", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get system_health: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of system_health records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*SystemHealth, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "system_health", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM system_health
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count system_health records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, last_check_at
			, last_success_at
			, last_failure_at
			, metric_value
			, metric_unit
			, threshold_warning
			, threshold_critical
			, details
			, created_at
			, updated_at
		FROM system_health
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list system_health", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list system_health: %w", err)
	}
	defer rows.Close()

	var entities []*SystemHealth
	for rows.Next() {
		var entity SystemHealth
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Status,
			&entity.LastCheckAt,
			&entity.LastSuccessAt,
			&entity.LastFailureAt,
			&entity.MetricValue,
			&entity.MetricUnit,
			&entity.ThresholdWarning,
			&entity.ThresholdCritical,
			&entity.Details,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan system_health: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating system_health rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing system_health record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *SystemHealth) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "system_health", duration, nil)
	}()

	query := `
		UPDATE system_health
		SET
			, organization_id = $2
			, status = $3
			, last_check_at = $4
			, last_success_at = $5
			, last_failure_at = $6
			, metric_value = $7
			, metric_unit = $8
			, threshold_warning = $9
			, threshold_critical = $10
			, details = $11
			, updated_at = $13
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Status,
		entity.LastCheckAt,
		entity.LastSuccessAt,
		entity.LastFailureAt,
		entity.MetricValue,
		entity.MetricUnit,
		entity.ThresholdWarning,
		entity.ThresholdCritical,
		entity.Details,
		time.Now(),
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update system_health", zap.Error(err))
		return fmt.Errorf("failed to update system_health: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("system_health not found or already deleted")
	}

	r.logger.Info("updated system_health",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a system_health record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "system_health", duration, nil)
	}()

	query := `DELETE FROM system_health WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete system_health", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete system_health: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("system_health not found")
	}

	r.logger.Info("deleted system_health", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves system_health records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*SystemHealth, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "system_health", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM system_health
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count system_health records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, last_check_at
			, last_success_at
			, last_failure_at
			, metric_value
			, metric_unit
			, threshold_warning
			, threshold_critical
			, details
			, created_at
			, updated_at
		FROM system_health
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list system_health by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list system_health: %w", err)
	}
	defer rows.Close()

	var entities []*SystemHealth
	for rows.Next() {
		var entity SystemHealth
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Status,
			&entity.LastCheckAt,
			&entity.LastSuccessAt,
			&entity.LastFailureAt,
			&entity.MetricValue,
			&entity.MetricUnit,
			&entity.ThresholdWarning,
			&entity.ThresholdCritical,
			&entity.Details,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan system_health: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


package analytic_plan

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

// Repository handles database operations for AnalyticPlans
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AnalyticPlans repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AnalyticPlans represents a analytic_plans entity
type AnalyticPlans struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PlanCode string `json:"plan_code" db:"plan_code"`
	PlanName string `json:"plan_name" db:"plan_name"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Description *string `json:"description" db:"description"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new analytic_plans record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AnalyticPlans) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "analytic_plans", duration, nil)
	}()

	query := `
		INSERT INTO analytic_plans (
			, organization_id
			, plan_code
			, plan_name
			, is_active
			, description
			, created_by
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $10
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PlanCode,
		entity.PlanName,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create analytic_plans", zap.Error(err))
		return fmt.Errorf("failed to create analytic_plans: %w", err)
	}

	r.logger.Info("created analytic_plans",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a analytic_plans by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AnalyticPlans, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "analytic_plans", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, plan_code
			, plan_name
			, is_active
			, description
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM analytic_plans
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity AnalyticPlans
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PlanCode,
		&entity.PlanName,
		&entity.IsActive,
		&entity.Description,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("analytic_plans not found")
	}

	if err != nil {
		r.logger.Error("failed to get analytic_plans", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get analytic_plans: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of analytic_plans records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AnalyticPlans, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "analytic_plans", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM analytic_plans
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count analytic_plans records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, plan_code
			, plan_name
			, is_active
			, description
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM analytic_plans
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list analytic_plans", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list analytic_plans: %w", err)
	}
	defer rows.Close()

	var entities []*AnalyticPlans
	for rows.Next() {
		var entity AnalyticPlans
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PlanCode,
			&entity.PlanName,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan analytic_plans: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating analytic_plans rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing analytic_plans record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AnalyticPlans) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "analytic_plans", duration, nil)
	}()

	query := `
		UPDATE analytic_plans
		SET
			, organization_id = $2
			, plan_code = $3
			, plan_name = $4
			, is_active = $5
			, description = $6
			, created_by = $7
			, updated_at = $9
			, deleted_at = $10
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PlanCode,
		entity.PlanName,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update analytic_plans", zap.Error(err))
		return fmt.Errorf("failed to update analytic_plans: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("analytic_plans not found or already deleted")
	}

	r.logger.Info("updated analytic_plans",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a analytic_plans record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "analytic_plans", duration, nil)
	}()

	query := `
		UPDATE analytic_plans
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete analytic_plans", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete analytic_plans: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("analytic_plans not found or already deleted")
	}

	r.logger.Info("deleted analytic_plans", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves analytic_plans records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*AnalyticPlans, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "analytic_plans", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM analytic_plans
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count analytic_plans records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, plan_code
			, plan_name
			, is_active
			, description
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM analytic_plans
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list analytic_plans by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list analytic_plans: %w", err)
	}
	defer rows.Close()

	var entities []*AnalyticPlans
	for rows.Next() {
		var entity AnalyticPlans
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PlanCode,
			&entity.PlanName,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan analytic_plans: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


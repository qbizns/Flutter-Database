package tax_group

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

// Repository handles database operations for TaxGroups
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new TaxGroups repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// TaxGroups represents a tax_groups entity
type TaxGroups struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	GroupCode string `json:"group_code" db:"group_code"`
	GroupName string `json:"group_name" db:"group_name"`
	Sequence *int64 `json:"sequence" db:"sequence"`
	IsActive *bool `json:"is_active" db:"is_active"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new tax_groups record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *TaxGroups) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "tax_groups", duration, nil)
	}()

	query := `
		INSERT INTO tax_groups (
			, organization_id
			, group_code
			, group_name
			, sequence
			, is_active
			, created_by
			, updated_by
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $11
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.GroupCode,
		entity.GroupName,
		entity.Sequence,
		entity.IsActive,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create tax_groups", zap.Error(err))
		return fmt.Errorf("failed to create tax_groups: %w", err)
	}

	r.logger.Info("created tax_groups",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a tax_groups by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*TaxGroups, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_groups", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, group_code
			, group_name
			, sequence
			, is_active
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM tax_groups
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity TaxGroups
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.GroupCode,
		&entity.GroupName,
		&entity.Sequence,
		&entity.IsActive,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("tax_groups not found")
	}

	if err != nil {
		r.logger.Error("failed to get tax_groups", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get tax_groups: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of tax_groups records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*TaxGroups, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_groups", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tax_groups
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tax_groups records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, group_code
			, group_name
			, sequence
			, is_active
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM tax_groups
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tax_groups", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tax_groups: %w", err)
	}
	defer rows.Close()

	var entities []*TaxGroups
	for rows.Next() {
		var entity TaxGroups
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.GroupCode,
			&entity.GroupName,
			&entity.Sequence,
			&entity.IsActive,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tax_groups: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating tax_groups rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing tax_groups record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *TaxGroups) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tax_groups", duration, nil)
	}()

	query := `
		UPDATE tax_groups
		SET
			, organization_id = $2
			, group_code = $3
			, group_name = $4
			, sequence = $5
			, is_active = $6
			, created_by = $7
			, updated_by = $8
			, updated_at = $10
			, deleted_at = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.GroupCode,
		entity.GroupName,
		entity.Sequence,
		entity.IsActive,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update tax_groups", zap.Error(err))
		return fmt.Errorf("failed to update tax_groups: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tax_groups not found or already deleted")
	}

	r.logger.Info("updated tax_groups",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a tax_groups record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tax_groups", duration, nil)
	}()

	query := `
		UPDATE tax_groups
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete tax_groups", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete tax_groups: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tax_groups not found or already deleted")
	}

	r.logger.Info("deleted tax_groups", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves tax_groups records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*TaxGroups, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_groups", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tax_groups
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tax_groups records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, group_code
			, group_name
			, sequence
			, is_active
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM tax_groups
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tax_groups by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tax_groups: %w", err)
	}
	defer rows.Close()

	var entities []*TaxGroups
	for rows.Next() {
		var entity TaxGroups
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.GroupCode,
			&entity.GroupName,
			&entity.Sequence,
			&entity.IsActive,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tax_groups: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


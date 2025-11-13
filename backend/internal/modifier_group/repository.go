package modifier_group

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

// Repository handles database operations for ModifierGroups
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ModifierGroups repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ModifierGroups represents a modifier_groups entity
type ModifierGroups struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	GroupName string `json:"group_name" db:"group_name"`
	GroupCode *string `json:"group_code" db:"group_code"`
	DisplayName *string `json:"display_name" db:"display_name"`
	SelectionType *string `json:"selection_type" db:"selection_type"`
	MinSelections *int64 `json:"min_selections" db:"min_selections"`
	MaxSelections *int64 `json:"max_selections" db:"max_selections"`
	ExactSelections *int64 `json:"exact_selections" db:"exact_selections"`
	IsRequired *bool `json:"is_required" db:"is_required"`
	AffectsPrice *bool `json:"affects_price" db:"affects_price"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	MinSelections *string `json:"min_selections" db:"min_selections"`
	(maxSelections *string `json:"(max_selections" db:"(max_selections"`
	(exactSelections *string `json:"(exact_selections" db:"(exact_selections"`
}

// Create inserts a new modifier_groups record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ModifierGroups) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "modifier_groups", duration, nil)
	}()

	query := `
		INSERT INTO modifier_groups (
			, organization_id
			, group_name
			, group_code
			, display_name
			, selection_type
			, min_selections
			, max_selections
			, exact_selections
			, is_required
			, affects_price
			, display_order
			, is_active
			, description
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, min_selections
			, (max_selections
			, (exact_selections
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
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.GroupName,
		entity.GroupCode,
		entity.DisplayName,
		entity.SelectionType,
		entity.MinSelections,
		entity.MaxSelections,
		entity.ExactSelections,
		entity.IsRequired,
		entity.AffectsPrice,
		entity.DisplayOrder,
		entity.IsActive,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.MinSelections,
		entity.(maxSelections,
		entity.(exactSelections,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create modifier_groups", zap.Error(err))
		return fmt.Errorf("failed to create modifier_groups: %w", err)
	}

	r.logger.Info("created modifier_groups",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a modifier_groups by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ModifierGroups, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "modifier_groups", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, group_name
			, group_code
			, display_name
			, selection_type
			, min_selections
			, max_selections
			, exact_selections
			, is_required
			, affects_price
			, display_order
			, is_active
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, min_selections
			, (max_selections
			, (exact_selections
		FROM modifier_groups
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ModifierGroups
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.GroupName,
		&entity.GroupCode,
		&entity.DisplayName,
		&entity.SelectionType,
		&entity.MinSelections,
		&entity.MaxSelections,
		&entity.ExactSelections,
		&entity.IsRequired,
		&entity.AffectsPrice,
		&entity.DisplayOrder,
		&entity.IsActive,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.MinSelections,
		&entity.(maxSelections,
		&entity.(exactSelections,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("modifier_groups not found")
	}

	if err != nil {
		r.logger.Error("failed to get modifier_groups", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get modifier_groups: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of modifier_groups records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ModifierGroups, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "modifier_groups", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM modifier_groups
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count modifier_groups records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, group_name
			, group_code
			, display_name
			, selection_type
			, min_selections
			, max_selections
			, exact_selections
			, is_required
			, affects_price
			, display_order
			, is_active
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, min_selections
			, (max_selections
			, (exact_selections
		FROM modifier_groups
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list modifier_groups", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list modifier_groups: %w", err)
	}
	defer rows.Close()

	var entities []*ModifierGroups
	for rows.Next() {
		var entity ModifierGroups
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.GroupName,
			&entity.GroupCode,
			&entity.DisplayName,
			&entity.SelectionType,
			&entity.MinSelections,
			&entity.MaxSelections,
			&entity.ExactSelections,
			&entity.IsRequired,
			&entity.AffectsPrice,
			&entity.DisplayOrder,
			&entity.IsActive,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.MinSelections,
			&entity.(maxSelections,
			&entity.(exactSelections,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan modifier_groups: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating modifier_groups rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing modifier_groups record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ModifierGroups) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "modifier_groups", duration, nil)
	}()

	query := `
		UPDATE modifier_groups
		SET
			, organization_id = $2
			, group_name = $3
			, group_code = $4
			, display_name = $5
			, selection_type = $6
			, min_selections = $7
			, max_selections = $8
			, exact_selections = $9
			, is_required = $10
			, affects_price = $11
			, display_order = $12
			, is_active = $13
			, description = $14
			, notes = $15
			, metadata = $16
			, updated_at = $18
			, deleted_at = $19
			, created_by = $20
			, updated_by = $21
			, min_selections = $22
			, (max_selections = $23
			, (exact_selections = $24
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $25
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.GroupName,
		entity.GroupCode,
		entity.DisplayName,
		entity.SelectionType,
		entity.MinSelections,
		entity.MaxSelections,
		entity.ExactSelections,
		entity.IsRequired,
		entity.AffectsPrice,
		entity.DisplayOrder,
		entity.IsActive,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.MinSelections,
		entity.(maxSelections,
		entity.(exactSelections,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update modifier_groups", zap.Error(err))
		return fmt.Errorf("failed to update modifier_groups: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("modifier_groups not found or already deleted")
	}

	r.logger.Info("updated modifier_groups",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a modifier_groups record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "modifier_groups", duration, nil)
	}()

	query := `
		UPDATE modifier_groups
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete modifier_groups", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete modifier_groups: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("modifier_groups not found or already deleted")
	}

	r.logger.Info("deleted modifier_groups", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves modifier_groups records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ModifierGroups, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "modifier_groups", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM modifier_groups
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count modifier_groups records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, group_name
			, group_code
			, display_name
			, selection_type
			, min_selections
			, max_selections
			, exact_selections
			, is_required
			, affects_price
			, display_order
			, is_active
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, min_selections
			, (max_selections
			, (exact_selections
		FROM modifier_groups
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list modifier_groups by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list modifier_groups: %w", err)
	}
	defer rows.Close()

	var entities []*ModifierGroups
	for rows.Next() {
		var entity ModifierGroups
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.GroupName,
			&entity.GroupCode,
			&entity.DisplayName,
			&entity.SelectionType,
			&entity.MinSelections,
			&entity.MaxSelections,
			&entity.ExactSelections,
			&entity.IsRequired,
			&entity.AffectsPrice,
			&entity.DisplayOrder,
			&entity.IsActive,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.MinSelections,
			&entity.(maxSelections,
			&entity.(exactSelections,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan modifier_groups: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


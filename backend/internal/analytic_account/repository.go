package analytic_account

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

// Repository handles database operations for AnalyticAccounts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AnalyticAccounts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AnalyticAccounts represents a analytic_accounts entity
type AnalyticAccounts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	AnalyticPlanId *uuid.UUID `json:"analytic_plan_id" db:"analytic_plan_id"`
	AccountCode string `json:"account_code" db:"account_code"`
	AccountName string `json:"account_name" db:"account_name"`
	ParentAccountId *uuid.UUID `json:"parent_account_id" db:"parent_account_id"`
	AccountLevel *int64 `json:"account_level" db:"account_level"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Description *string `json:"description" db:"description"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new analytic_accounts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AnalyticAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "analytic_accounts", duration, nil)
	}()

	query := `
		INSERT INTO analytic_accounts (
			, organization_id
			, analytic_plan_id
			, account_code
			, account_name
			, parent_account_id
			, account_level
			, is_active
			, description
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
			, $9
			, $10
			, $11
			, $14
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.AnalyticPlanId,
		entity.AccountCode,
		entity.AccountName,
		entity.ParentAccountId,
		entity.AccountLevel,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create analytic_accounts", zap.Error(err))
		return fmt.Errorf("failed to create analytic_accounts: %w", err)
	}

	r.logger.Info("created analytic_accounts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a analytic_accounts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AnalyticAccounts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "analytic_accounts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, analytic_plan_id
			, account_code
			, account_name
			, parent_account_id
			, account_level
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM analytic_accounts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity AnalyticAccounts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.AnalyticPlanId,
		&entity.AccountCode,
		&entity.AccountName,
		&entity.ParentAccountId,
		&entity.AccountLevel,
		&entity.IsActive,
		&entity.Description,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("analytic_accounts not found")
	}

	if err != nil {
		r.logger.Error("failed to get analytic_accounts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get analytic_accounts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of analytic_accounts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AnalyticAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "analytic_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM analytic_accounts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count analytic_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, analytic_plan_id
			, account_code
			, account_name
			, parent_account_id
			, account_level
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM analytic_accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list analytic_accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list analytic_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*AnalyticAccounts
	for rows.Next() {
		var entity AnalyticAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.AnalyticPlanId,
			&entity.AccountCode,
			&entity.AccountName,
			&entity.ParentAccountId,
			&entity.AccountLevel,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan analytic_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating analytic_accounts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing analytic_accounts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AnalyticAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "analytic_accounts", duration, nil)
	}()

	query := `
		UPDATE analytic_accounts
		SET
			, organization_id = $2
			, analytic_plan_id = $3
			, account_code = $4
			, account_name = $5
			, parent_account_id = $6
			, account_level = $7
			, is_active = $8
			, description = $9
			, created_by = $10
			, updated_by = $11
			, updated_at = $13
			, deleted_at = $14
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $15
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.AnalyticPlanId,
		entity.AccountCode,
		entity.AccountName,
		entity.ParentAccountId,
		entity.AccountLevel,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedBy,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update analytic_accounts", zap.Error(err))
		return fmt.Errorf("failed to update analytic_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("analytic_accounts not found or already deleted")
	}

	r.logger.Info("updated analytic_accounts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a analytic_accounts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "analytic_accounts", duration, nil)
	}()

	query := `
		UPDATE analytic_accounts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete analytic_accounts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete analytic_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("analytic_accounts not found or already deleted")
	}

	r.logger.Info("deleted analytic_accounts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves analytic_accounts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*AnalyticAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "analytic_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM analytic_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count analytic_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, analytic_plan_id
			, account_code
			, account_name
			, parent_account_id
			, account_level
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM analytic_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list analytic_accounts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list analytic_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*AnalyticAccounts
	for rows.Next() {
		var entity AnalyticAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.AnalyticPlanId,
			&entity.AccountCode,
			&entity.AccountName,
			&entity.ParentAccountId,
			&entity.AccountLevel,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan analytic_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


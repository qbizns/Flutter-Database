package cash_drawer_session

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

// Repository handles database operations for CashDrawerSessions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CashDrawerSessions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CashDrawerSessions represents a cash_drawer_sessions entity
type CashDrawerSessions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CashDrawerId uuid.UUID `json:"cash_drawer_id" db:"cash_drawer_id"`
	PosSessionId uuid.UUID `json:"pos_session_id" db:"pos_session_id"`
	OpeningAmount *float64 `json:"opening_amount" db:"opening_amount"`
	ClosingAmount *float64 `json:"closing_amount" db:"closing_amount"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new cash_drawer_sessions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CashDrawerSessions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "cash_drawer_sessions", duration, nil)
	}()

	query := `
		INSERT INTO cash_drawer_sessions (
			, organization_id
			, cash_drawer_id
			, pos_session_id
			, opening_amount
			, closing_amount
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $8
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CashDrawerId,
		entity.PosSessionId,
		entity.OpeningAmount,
		entity.ClosingAmount,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create cash_drawer_sessions", zap.Error(err))
		return fmt.Errorf("failed to create cash_drawer_sessions: %w", err)
	}

	r.logger.Info("created cash_drawer_sessions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a cash_drawer_sessions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CashDrawerSessions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_drawer_sessions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, cash_drawer_id
			, pos_session_id
			, opening_amount
			, closing_amount
			, created_at
			, deleted_at
		FROM cash_drawer_sessions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CashDrawerSessions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CashDrawerId,
		&entity.PosSessionId,
		&entity.OpeningAmount,
		&entity.ClosingAmount,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("cash_drawer_sessions not found")
	}

	if err != nil {
		r.logger.Error("failed to get cash_drawer_sessions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get cash_drawer_sessions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of cash_drawer_sessions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CashDrawerSessions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_drawer_sessions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cash_drawer_sessions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cash_drawer_sessions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, cash_drawer_id
			, pos_session_id
			, opening_amount
			, closing_amount
			, created_at
			, deleted_at
		FROM cash_drawer_sessions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cash_drawer_sessions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cash_drawer_sessions: %w", err)
	}
	defer rows.Close()

	var entities []*CashDrawerSessions
	for rows.Next() {
		var entity CashDrawerSessions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CashDrawerId,
			&entity.PosSessionId,
			&entity.OpeningAmount,
			&entity.ClosingAmount,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cash_drawer_sessions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating cash_drawer_sessions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing cash_drawer_sessions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CashDrawerSessions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cash_drawer_sessions", duration, nil)
	}()

	query := `
		UPDATE cash_drawer_sessions
		SET
			, organization_id = $2
			, cash_drawer_id = $3
			, pos_session_id = $4
			, opening_amount = $5
			, closing_amount = $6
			, deleted_at = $8
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CashDrawerId,
		entity.PosSessionId,
		entity.OpeningAmount,
		entity.ClosingAmount,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update cash_drawer_sessions", zap.Error(err))
		return fmt.Errorf("failed to update cash_drawer_sessions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cash_drawer_sessions not found or already deleted")
	}

	r.logger.Info("updated cash_drawer_sessions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a cash_drawer_sessions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cash_drawer_sessions", duration, nil)
	}()

	query := `
		UPDATE cash_drawer_sessions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete cash_drawer_sessions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete cash_drawer_sessions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cash_drawer_sessions not found or already deleted")
	}

	r.logger.Info("deleted cash_drawer_sessions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves cash_drawer_sessions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CashDrawerSessions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_drawer_sessions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cash_drawer_sessions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cash_drawer_sessions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, cash_drawer_id
			, pos_session_id
			, opening_amount
			, closing_amount
			, created_at
			, deleted_at
		FROM cash_drawer_sessions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cash_drawer_sessions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cash_drawer_sessions: %w", err)
	}
	defer rows.Close()

	var entities []*CashDrawerSessions
	for rows.Next() {
		var entity CashDrawerSessions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CashDrawerId,
			&entity.PosSessionId,
			&entity.OpeningAmount,
			&entity.ClosingAmount,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cash_drawer_sessions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


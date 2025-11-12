package currency_rate

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

// Repository handles database operations for CurrencyRates
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CurrencyRates repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CurrencyRates represents a currency_rates entity
type CurrencyRates struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CurrencyCode string `json:"currency_code" db:"currency_code"`
	RateDate time.Time `json:"rate_date" db:"rate_date"`
	Rate float64 `json:"rate" db:"rate"`
	Source *string `json:"source" db:"source"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new currency_rates record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CurrencyRates) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "currency_rates", duration, nil)
	}()

	query := `
		INSERT INTO currency_rates (
			, organization_id
			, currency_code
			, rate_date
			, rate
			, source
			, created_by
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $9
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CurrencyCode,
		entity.RateDate,
		entity.Rate,
		entity.Source,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create currency_rates", zap.Error(err))
		return fmt.Errorf("failed to create currency_rates: %w", err)
	}

	r.logger.Info("created currency_rates",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a currency_rates by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CurrencyRates, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "currency_rates", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, currency_code
			, rate_date
			, rate
			, source
			, created_by
			, created_at
			, deleted_at
		FROM currency_rates
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CurrencyRates
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CurrencyCode,
		&entity.RateDate,
		&entity.Rate,
		&entity.Source,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("currency_rates not found")
	}

	if err != nil {
		r.logger.Error("failed to get currency_rates", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get currency_rates: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of currency_rates records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CurrencyRates, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "currency_rates", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM currency_rates
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count currency_rates records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, currency_code
			, rate_date
			, rate
			, source
			, created_by
			, created_at
			, deleted_at
		FROM currency_rates
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list currency_rates", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list currency_rates: %w", err)
	}
	defer rows.Close()

	var entities []*CurrencyRates
	for rows.Next() {
		var entity CurrencyRates
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CurrencyCode,
			&entity.RateDate,
			&entity.Rate,
			&entity.Source,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan currency_rates: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating currency_rates rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing currency_rates record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CurrencyRates) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "currency_rates", duration, nil)
	}()

	query := `
		UPDATE currency_rates
		SET
			, organization_id = $2
			, currency_code = $3
			, rate_date = $4
			, rate = $5
			, source = $6
			, created_by = $7
			, deleted_at = $9
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CurrencyCode,
		entity.RateDate,
		entity.Rate,
		entity.Source,
		entity.CreatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update currency_rates", zap.Error(err))
		return fmt.Errorf("failed to update currency_rates: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("currency_rates not found or already deleted")
	}

	r.logger.Info("updated currency_rates",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a currency_rates record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "currency_rates", duration, nil)
	}()

	query := `
		UPDATE currency_rates
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete currency_rates", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete currency_rates: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("currency_rates not found or already deleted")
	}

	r.logger.Info("deleted currency_rates", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves currency_rates records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CurrencyRates, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "currency_rates", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM currency_rates
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count currency_rates records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, currency_code
			, rate_date
			, rate
			, source
			, created_by
			, created_at
			, deleted_at
		FROM currency_rates
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list currency_rates by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list currency_rates: %w", err)
	}
	defer rows.Close()

	var entities []*CurrencyRates
	for rows.Next() {
		var entity CurrencyRates
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CurrencyCode,
			&entity.RateDate,
			&entity.Rate,
			&entity.Source,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan currency_rates: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


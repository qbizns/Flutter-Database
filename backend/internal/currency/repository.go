package currency

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

// Repository handles database operations for Currencies
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Currencies repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Currencies represents a currencies entity
type Currencies struct {
	Id *uuid.UUID `json:"id" db:"id"`
	CurrencyCode string `json:"currency_code" db:"currency_code"`
	CurrencyName string `json:"currency_name" db:"currency_name"`
	CurrencySymbol *string `json:"currency_symbol" db:"currency_symbol"`
	DecimalPlaces *int64 `json:"decimal_places" db:"decimal_places"`
	IsActive *bool `json:"is_active" db:"is_active"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new currencies record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Currencies) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "currencies", duration, nil)
	}()

	query := `
		INSERT INTO currencies (
			, currency_code
			, currency_name
			, currency_symbol
			, decimal_places
			, is_active
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $9
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.CurrencyCode,
		entity.CurrencyName,
		entity.CurrencySymbol,
		entity.DecimalPlaces,
		entity.IsActive,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create currencies", zap.Error(err))
		return fmt.Errorf("failed to create currencies: %w", err)
	}

	r.logger.Info("created currencies",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a currencies by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Currencies, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "currencies", duration, nil)
	}()

	query := `
		SELECT
			id
			, currency_code
			, currency_name
			, currency_symbol
			, decimal_places
			, is_active
			, created_at
			, updated_at
			, deleted_at
		FROM currencies
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Currencies
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.CurrencyCode,
		&entity.CurrencyName,
		&entity.CurrencySymbol,
		&entity.DecimalPlaces,
		&entity.IsActive,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("currencies not found")
	}

	if err != nil {
		r.logger.Error("failed to get currencies", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get currencies: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of currencies records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Currencies, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "currencies", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM currencies
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count currencies records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, currency_code
			, currency_name
			, currency_symbol
			, decimal_places
			, is_active
			, created_at
			, updated_at
			, deleted_at
		FROM currencies
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list currencies", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list currencies: %w", err)
	}
	defer rows.Close()

	var entities []*Currencies
	for rows.Next() {
		var entity Currencies
		err := rows.Scan(
			&entity.Id,
			&entity.CurrencyCode,
			&entity.CurrencyName,
			&entity.CurrencySymbol,
			&entity.DecimalPlaces,
			&entity.IsActive,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan currencies: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating currencies rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing currencies record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Currencies) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "currencies", duration, nil)
	}()

	query := `
		UPDATE currencies
		SET
			, currency_code = $2
			, currency_name = $3
			, currency_symbol = $4
			, decimal_places = $5
			, is_active = $6
			, updated_at = $8
			, deleted_at = $9
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.CurrencyCode,
		entity.CurrencyName,
		entity.CurrencySymbol,
		entity.DecimalPlaces,
		entity.IsActive,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update currencies", zap.Error(err))
		return fmt.Errorf("failed to update currencies: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("currencies not found or already deleted")
	}

	r.logger.Info("updated currencies",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a currencies record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "currencies", duration, nil)
	}()

	query := `
		UPDATE currencies
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete currencies", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete currencies: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("currencies not found or already deleted")
	}

	r.logger.Info("deleted currencies", zap.String("id", id.String()))
	return nil
}




package payment_term_line

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

// Repository handles database operations for PaymentTermLines
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PaymentTermLines repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PaymentTermLines represents a payment_term_lines entity
type PaymentTermLines struct {
	Id *uuid.UUID `json:"id" db:"id"`
	PaymentTermId uuid.UUID `json:"payment_term_id" db:"payment_term_id"`
	Sequence int64 `json:"sequence" db:"sequence"`
	ValueType string `json:"value_type" db:"value_type"`
	ValueAmount *float64 `json:"value_amount" db:"value_amount"`
	DaysAfter *int64 `json:"days_after" db:"days_after"`
	EndOfMonth *bool `json:"end_of_month" db:"end_of_month"`
	DayOfMonth *int64 `json:"day_of_month" db:"day_of_month"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	(valueType string `json:"(value_type" db:"(value_type"`
}

// Create inserts a new payment_term_lines record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PaymentTermLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "payment_term_lines", duration, nil)
	}()

	query := `
		INSERT INTO payment_term_lines (
			, payment_term_id
			, sequence
			, value_type
			, value_amount
			, days_after
			, end_of_month
			, day_of_month
			, deleted_at
			, (value_type
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $10
			, $11
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.PaymentTermId,
		entity.Sequence,
		entity.ValueType,
		entity.ValueAmount,
		entity.DaysAfter,
		entity.EndOfMonth,
		entity.DayOfMonth,
		entity.DeletedAt,
		entity.(valueType,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create payment_term_lines", zap.Error(err))
		return fmt.Errorf("failed to create payment_term_lines: %w", err)
	}

	r.logger.Info("created payment_term_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a payment_term_lines by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PaymentTermLines, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payment_term_lines", duration, nil)
	}()

	query := `
		SELECT
			id
			, payment_term_id
			, sequence
			, value_type
			, value_amount
			, days_after
			, end_of_month
			, day_of_month
			, created_at
			, deleted_at
			, (value_type
		FROM payment_term_lines
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PaymentTermLines
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.PaymentTermId,
		&entity.Sequence,
		&entity.ValueType,
		&entity.ValueAmount,
		&entity.DaysAfter,
		&entity.EndOfMonth,
		&entity.DayOfMonth,
		&entity.CreatedAt,
		&entity.DeletedAt,
		&entity.(valueType,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("payment_term_lines not found")
	}

	if err != nil {
		r.logger.Error("failed to get payment_term_lines", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get payment_term_lines: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of payment_term_lines records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PaymentTermLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payment_term_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM payment_term_lines
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payment_term_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, payment_term_id
			, sequence
			, value_type
			, value_amount
			, days_after
			, end_of_month
			, day_of_month
			, created_at
			, deleted_at
			, (value_type
		FROM payment_term_lines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list payment_term_lines", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list payment_term_lines: %w", err)
	}
	defer rows.Close()

	var entities []*PaymentTermLines
	for rows.Next() {
		var entity PaymentTermLines
		err := rows.Scan(
			&entity.Id,
			&entity.PaymentTermId,
			&entity.Sequence,
			&entity.ValueType,
			&entity.ValueAmount,
			&entity.DaysAfter,
			&entity.EndOfMonth,
			&entity.DayOfMonth,
			&entity.CreatedAt,
			&entity.DeletedAt,
			&entity.(valueType,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment_term_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating payment_term_lines rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing payment_term_lines record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PaymentTermLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "payment_term_lines", duration, nil)
	}()

	query := `
		UPDATE payment_term_lines
		SET
			, payment_term_id = $2
			, sequence = $3
			, value_type = $4
			, value_amount = $5
			, days_after = $6
			, end_of_month = $7
			, day_of_month = $8
			, deleted_at = $10
			, (value_type = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.PaymentTermId,
		entity.Sequence,
		entity.ValueType,
		entity.ValueAmount,
		entity.DaysAfter,
		entity.EndOfMonth,
		entity.DayOfMonth,
		entity.DeletedAt,
		entity.(valueType,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update payment_term_lines", zap.Error(err))
		return fmt.Errorf("failed to update payment_term_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payment_term_lines not found or already deleted")
	}

	r.logger.Info("updated payment_term_lines",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a payment_term_lines record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "payment_term_lines", duration, nil)
	}()

	query := `
		UPDATE payment_term_lines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete payment_term_lines", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete payment_term_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payment_term_lines not found or already deleted")
	}

	r.logger.Info("deleted payment_term_lines", zap.String("id", id.String()))
	return nil
}




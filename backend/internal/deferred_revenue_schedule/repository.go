package deferred_revenue_schedule

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

// Repository handles database operations for DeferredRevenueSchedule
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DeferredRevenueSchedule repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DeferredRevenueSchedule represents a deferred_revenue_schedule entity
type DeferredRevenueSchedule struct {
	Id *uuid.UUID `json:"id" db:"id"`
	ContractId uuid.UUID `json:"contract_id" db:"contract_id"`
	LineNumber int64 `json:"line_number" db:"line_number"`
	RecognitionDate time.Time `json:"recognition_date" db:"recognition_date"`
	RecognitionAmount float64 `json:"recognition_amount" db:"recognition_amount"`
	Status *string `json:"status" db:"status"`
	JournalEntryId *uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	PostedAt *time.Time `json:"posted_at" db:"posted_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new deferred_revenue_schedule record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DeferredRevenueSchedule) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "deferred_revenue_schedule", duration, nil)
	}()

	query := `
		INSERT INTO deferred_revenue_schedule (
			, contract_id
			, line_number
			, recognition_date
			, recognition_amount
			, status
			, journal_entry_id
			, posted_at
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $9
			, $10
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.ContractId,
		entity.LineNumber,
		entity.RecognitionDate,
		entity.RecognitionAmount,
		entity.Status,
		entity.JournalEntryId,
		entity.PostedAt,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create deferred_revenue_schedule", zap.Error(err))
		return fmt.Errorf("failed to create deferred_revenue_schedule: %w", err)
	}

	r.logger.Info("created deferred_revenue_schedule",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a deferred_revenue_schedule by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DeferredRevenueSchedule, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "deferred_revenue_schedule", duration, nil)
	}()

	query := `
		SELECT
			id
			, contract_id
			, line_number
			, recognition_date
			, recognition_amount
			, status
			, journal_entry_id
			, created_at
			, posted_at
			, deleted_at
		FROM deferred_revenue_schedule
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity DeferredRevenueSchedule
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.ContractId,
		&entity.LineNumber,
		&entity.RecognitionDate,
		&entity.RecognitionAmount,
		&entity.Status,
		&entity.JournalEntryId,
		&entity.CreatedAt,
		&entity.PostedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("deferred_revenue_schedule not found")
	}

	if err != nil {
		r.logger.Error("failed to get deferred_revenue_schedule", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get deferred_revenue_schedule: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of deferred_revenue_schedule records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DeferredRevenueSchedule, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "deferred_revenue_schedule", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM deferred_revenue_schedule
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count deferred_revenue_schedule records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, contract_id
			, line_number
			, recognition_date
			, recognition_amount
			, status
			, journal_entry_id
			, created_at
			, posted_at
			, deleted_at
		FROM deferred_revenue_schedule
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list deferred_revenue_schedule", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list deferred_revenue_schedule: %w", err)
	}
	defer rows.Close()

	var entities []*DeferredRevenueSchedule
	for rows.Next() {
		var entity DeferredRevenueSchedule
		err := rows.Scan(
			&entity.Id,
			&entity.ContractId,
			&entity.LineNumber,
			&entity.RecognitionDate,
			&entity.RecognitionAmount,
			&entity.Status,
			&entity.JournalEntryId,
			&entity.CreatedAt,
			&entity.PostedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan deferred_revenue_schedule: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating deferred_revenue_schedule rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing deferred_revenue_schedule record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DeferredRevenueSchedule) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "deferred_revenue_schedule", duration, nil)
	}()

	query := `
		UPDATE deferred_revenue_schedule
		SET
			, contract_id = $2
			, line_number = $3
			, recognition_date = $4
			, recognition_amount = $5
			, status = $6
			, journal_entry_id = $7
			, posted_at = $9
			, deleted_at = $10
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.ContractId,
		entity.LineNumber,
		entity.RecognitionDate,
		entity.RecognitionAmount,
		entity.Status,
		entity.JournalEntryId,
		entity.PostedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update deferred_revenue_schedule", zap.Error(err))
		return fmt.Errorf("failed to update deferred_revenue_schedule: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("deferred_revenue_schedule not found or already deleted")
	}

	r.logger.Info("updated deferred_revenue_schedule",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a deferred_revenue_schedule record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "deferred_revenue_schedule", duration, nil)
	}()

	query := `
		UPDATE deferred_revenue_schedule
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete deferred_revenue_schedule", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete deferred_revenue_schedule: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("deferred_revenue_schedule not found or already deleted")
	}

	r.logger.Info("deleted deferred_revenue_schedule", zap.String("id", id.String()))
	return nil
}




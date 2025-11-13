package invoice_payment_schedule

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

// Repository handles database operations for InvoicePaymentSchedules
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new InvoicePaymentSchedules repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// InvoicePaymentSchedules represents a invoice_payment_schedules entity
type InvoicePaymentSchedules struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SourceType string `json:"source_type" db:"source_type"`
	SourceId uuid.UUID `json:"source_id" db:"source_id"`
	LineNumber int64 `json:"line_number" db:"line_number"`
	DueDate time.Time `json:"due_date" db:"due_date"`
	AmountDue float64 `json:"amount_due" db:"amount_due"`
	AmountPaid *float64 `json:"amount_paid" db:"amount_paid"`
	Status *string `json:"status" db:"status"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new invoice_payment_schedules record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *InvoicePaymentSchedules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "invoice_payment_schedules", duration, nil)
	}()

	query := `
		INSERT INTO invoice_payment_schedules (
			, organization_id
			, source_type
			, source_id
			, line_number
			, due_date
			, amount_due
			, amount_paid
			, status
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
			, $12
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SourceType,
		entity.SourceId,
		entity.LineNumber,
		entity.DueDate,
		entity.AmountDue,
		entity.AmountPaid,
		entity.Status,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create invoice_payment_schedules", zap.Error(err))
		return fmt.Errorf("failed to create invoice_payment_schedules: %w", err)
	}

	r.logger.Info("created invoice_payment_schedules",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a invoice_payment_schedules by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*InvoicePaymentSchedules, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "invoice_payment_schedules", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, source_type
			, source_id
			, line_number
			, due_date
			, amount_due
			, amount_paid
			, created_at
			, updated_at
			, deleted_at
		FROM invoice_payment_schedules
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity InvoicePaymentSchedules
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SourceType,
		&entity.SourceId,
		&entity.LineNumber,
		&entity.DueDate,
		&entity.AmountDue,
		&entity.AmountPaid,
		&entity.Status,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("invoice_payment_schedules not found")
	}

	if err != nil {
		r.logger.Error("failed to get invoice_payment_schedules", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get invoice_payment_schedules: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of invoice_payment_schedules records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*InvoicePaymentSchedules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "invoice_payment_schedules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM invoice_payment_schedules
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count invoice_payment_schedules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_type
			, source_id
			, line_number
			, due_date
			, amount_due
			, amount_paid
			, created_at
			, updated_at
			, deleted_at
		FROM invoice_payment_schedules
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list invoice_payment_schedules", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list invoice_payment_schedules: %w", err)
	}
	defer rows.Close()

	var entities []*InvoicePaymentSchedules
	for rows.Next() {
		var entity InvoicePaymentSchedules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceType,
			&entity.SourceId,
			&entity.LineNumber,
			&entity.DueDate,
			&entity.AmountDue,
			&entity.AmountPaid,
			&entity.Status,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan invoice_payment_schedules: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating invoice_payment_schedules rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing invoice_payment_schedules record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *InvoicePaymentSchedules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "invoice_payment_schedules", duration, nil)
	}()

	query := `
		UPDATE invoice_payment_schedules
		SET
			, organization_id = $2
			, source_type = $3
			, source_id = $4
			, line_number = $5
			, due_date = $6
			, amount_due = $7
			, amount_paid = $8
			, status = $9
			, updated_at = $11
			, deleted_at = $12
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SourceType,
		entity.SourceId,
		entity.LineNumber,
		entity.DueDate,
		entity.AmountDue,
		entity.AmountPaid,
		entity.Status,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update invoice_payment_schedules", zap.Error(err))
		return fmt.Errorf("failed to update invoice_payment_schedules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invoice_payment_schedules not found or already deleted")
	}

	r.logger.Info("updated invoice_payment_schedules",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a invoice_payment_schedules record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "invoice_payment_schedules", duration, nil)
	}()

	query := `
		UPDATE invoice_payment_schedules
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete invoice_payment_schedules", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete invoice_payment_schedules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invoice_payment_schedules not found or already deleted")
	}

	r.logger.Info("deleted invoice_payment_schedules", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves invoice_payment_schedules records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*InvoicePaymentSchedules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "invoice_payment_schedules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM invoice_payment_schedules
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count invoice_payment_schedules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_type
			, source_id
			, line_number
			, due_date
			, amount_due
			, amount_paid
			, status
			, created_at
			, updated_at
			, deleted_at
		FROM invoice_payment_schedules
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list invoice_payment_schedules by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list invoice_payment_schedules: %w", err)
	}
	defer rows.Close()

	var entities []*InvoicePaymentSchedules
	for rows.Next() {
		var entity InvoicePaymentSchedules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceType,
			&entity.SourceId,
			&entity.LineNumber,
			&entity.DueDate,
			&entity.AmountDue,
			&entity.AmountPaid,
			&entity.Status,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan invoice_payment_schedules: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


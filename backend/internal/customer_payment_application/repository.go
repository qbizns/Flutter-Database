package customer_payment_application

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

// Repository handles database operations for CustomerPaymentApplications
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CustomerPaymentApplications repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CustomerPaymentApplications represents a customer_payment_applications entity
type CustomerPaymentApplications struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CustomerPaymentId uuid.UUID `json:"customer_payment_id" db:"customer_payment_id"`
	CustomerInvoiceId uuid.UUID `json:"customer_invoice_id" db:"customer_invoice_id"`
	AppliedAmount float64 `json:"applied_amount" db:"applied_amount"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new customer_payment_applications record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CustomerPaymentApplications) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "customer_payment_applications", duration, nil)
	}()

	query := `
		INSERT INTO customer_payment_applications (
			, organization_id
			, customer_payment_id
			, customer_invoice_id
			, applied_amount
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $7
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CustomerPaymentId,
		entity.CustomerInvoiceId,
		entity.AppliedAmount,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create customer_payment_applications", zap.Error(err))
		return fmt.Errorf("failed to create customer_payment_applications: %w", err)
	}

	r.logger.Info("created customer_payment_applications",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a customer_payment_applications by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CustomerPaymentApplications, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_payment_applications", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, customer_payment_id
			, customer_invoice_id
			, applied_amount
			, created_at
			, deleted_at
		FROM customer_payment_applications
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CustomerPaymentApplications
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CustomerPaymentId,
		&entity.CustomerInvoiceId,
		&entity.AppliedAmount,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("customer_payment_applications not found")
	}

	if err != nil {
		r.logger.Error("failed to get customer_payment_applications", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get customer_payment_applications: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of customer_payment_applications records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CustomerPaymentApplications, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_payment_applications", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_payment_applications
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_payment_applications records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_payment_id
			, customer_invoice_id
			, applied_amount
			, created_at
			, deleted_at
		FROM customer_payment_applications
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_payment_applications", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_payment_applications: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerPaymentApplications
	for rows.Next() {
		var entity CustomerPaymentApplications
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerPaymentId,
			&entity.CustomerInvoiceId,
			&entity.AppliedAmount,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_payment_applications: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating customer_payment_applications rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing customer_payment_applications record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CustomerPaymentApplications) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customer_payment_applications", duration, nil)
	}()

	query := `
		UPDATE customer_payment_applications
		SET
			, organization_id = $2
			, customer_payment_id = $3
			, customer_invoice_id = $4
			, applied_amount = $5
			, deleted_at = $7
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CustomerPaymentId,
		entity.CustomerInvoiceId,
		entity.AppliedAmount,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update customer_payment_applications", zap.Error(err))
		return fmt.Errorf("failed to update customer_payment_applications: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_payment_applications not found or already deleted")
	}

	r.logger.Info("updated customer_payment_applications",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a customer_payment_applications record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customer_payment_applications", duration, nil)
	}()

	query := `
		UPDATE customer_payment_applications
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete customer_payment_applications", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete customer_payment_applications: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_payment_applications not found or already deleted")
	}

	r.logger.Info("deleted customer_payment_applications", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves customer_payment_applications records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CustomerPaymentApplications, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_payment_applications", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_payment_applications
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_payment_applications records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_payment_id
			, customer_invoice_id
			, applied_amount
			, created_at
			, deleted_at
		FROM customer_payment_applications
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_payment_applications by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_payment_applications: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerPaymentApplications
	for rows.Next() {
		var entity CustomerPaymentApplications
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerPaymentId,
			&entity.CustomerInvoiceId,
			&entity.AppliedAmount,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_payment_applications: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


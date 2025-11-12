package payment

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

// Repository handles database operations for Payments
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Payments repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Payments represents a payments entity
type Payments struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SaleId uuid.UUID `json:"sale_id" db:"sale_id"`
	PaymentMethod string `json:"payment_method" db:"payment_method"`
	PaymentStatus string `json:"payment_status" db:"payment_status"`
	Amount float64 `json:"amount" db:"amount"`
	CardLastFour *string `json:"card_last_four" db:"card_last_four"`
	CardType *string `json:"card_type" db:"card_type"`
	TransactionId *string `json:"transaction_id" db:"transaction_id"`
	ReferenceNumber *string `json:"reference_number" db:"reference_number"`
	AccountNumber *string `json:"account_number" db:"account_number"`
	AccountName *string `json:"account_name" db:"account_name"`
	PaymentDate time.Time `json:"payment_date" db:"payment_date"`
	ProcessedAt *time.Time `json:"processed_at" db:"processed_at"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
}

// Create inserts a new payments record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Payments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "payments", duration, nil)
	}()

	query := `
		INSERT INTO payments (
			, organization_id
			, sale_id
			, payment_method
			, payment_status
			, amount
			, card_last_four
			, card_type
			, transaction_id
			, reference_number
			, account_number
			, account_name
			, payment_date
			, processed_at
			, notes
			, metadata
			, created_by
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SaleId,
		entity.PaymentMethod,
		entity.PaymentStatus,
		entity.Amount,
		entity.CardLastFour,
		entity.CardType,
		entity.TransactionId,
		entity.ReferenceNumber,
		entity.AccountNumber,
		entity.AccountName,
		entity.PaymentDate,
		entity.ProcessedAt,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create payments", zap.Error(err))
		return fmt.Errorf("failed to create payments: %w", err)
	}

	r.logger.Info("created payments",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a payments by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Payments, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payments", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, sale_id
			, payment_method
			, payment_status
			, amount
			, card_last_four
			, card_type
			, transaction_id
			, reference_number
			, account_number
			, account_name
			, payment_date
			, processed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
		FROM payments
		WHERE id = $1
		
	`

	var entity Payments
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SaleId,
		&entity.PaymentMethod,
		&entity.PaymentStatus,
		&entity.Amount,
		&entity.CardLastFour,
		&entity.CardType,
		&entity.TransactionId,
		&entity.ReferenceNumber,
		&entity.AccountNumber,
		&entity.AccountName,
		&entity.PaymentDate,
		&entity.ProcessedAt,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("payments not found")
	}

	if err != nil {
		r.logger.Error("failed to get payments", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of payments records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Payments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM payments
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_id
			, payment_method
			, payment_status
			, amount
			, card_last_four
			, card_type
			, transaction_id
			, reference_number
			, account_number
			, account_name
			, payment_date
			, processed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
		FROM payments
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list payments", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var entities []*Payments
	for rows.Next() {
		var entity Payments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleId,
			&entity.PaymentMethod,
			&entity.PaymentStatus,
			&entity.Amount,
			&entity.CardLastFour,
			&entity.CardType,
			&entity.TransactionId,
			&entity.ReferenceNumber,
			&entity.AccountNumber,
			&entity.AccountName,
			&entity.PaymentDate,
			&entity.ProcessedAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payments: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating payments rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing payments record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Payments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "payments", duration, nil)
	}()

	query := `
		UPDATE payments
		SET
			, organization_id = $2
			, sale_id = $3
			, payment_method = $4
			, payment_status = $5
			, amount = $6
			, card_last_four = $7
			, card_type = $8
			, transaction_id = $9
			, reference_number = $10
			, account_number = $11
			, account_name = $12
			, payment_date = $13
			, processed_at = $14
			, notes = $15
			, metadata = $16
			, updated_at = $18
			, created_by = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SaleId,
		entity.PaymentMethod,
		entity.PaymentStatus,
		entity.Amount,
		entity.CardLastFour,
		entity.CardType,
		entity.TransactionId,
		entity.ReferenceNumber,
		entity.AccountNumber,
		entity.AccountName,
		entity.PaymentDate,
		entity.ProcessedAt,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update payments", zap.Error(err))
		return fmt.Errorf("failed to update payments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payments not found or already deleted")
	}

	r.logger.Info("updated payments",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a payments record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "payments", duration, nil)
	}()

	query := `DELETE FROM payments WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete payments", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete payments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payments not found")
	}

	r.logger.Info("deleted payments", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves payments records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Payments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM payments
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_id
			, payment_method
			, payment_status
			, amount
			, card_last_four
			, card_type
			, transaction_id
			, reference_number
			, account_number
			, account_name
			, payment_date
			, processed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
		FROM payments
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list payments by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var entities []*Payments
	for rows.Next() {
		var entity Payments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleId,
			&entity.PaymentMethod,
			&entity.PaymentStatus,
			&entity.Amount,
			&entity.CardLastFour,
			&entity.CardType,
			&entity.TransactionId,
			&entity.ReferenceNumber,
			&entity.AccountNumber,
			&entity.AccountName,
			&entity.PaymentDate,
			&entity.ProcessedAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payments: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


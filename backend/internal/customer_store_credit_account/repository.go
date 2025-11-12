package customer_store_credit_account

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

// Repository handles database operations for CustomerStoreCreditAccounts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CustomerStoreCreditAccounts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CustomerStoreCreditAccounts represents a customer_store_credit_accounts entity
type CustomerStoreCreditAccounts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CustomerId uuid.UUID `json:"customer_id" db:"customer_id"`
	CurrentBalance *float64 `json:"current_balance" db:"current_balance"`
	CreditLimit *float64 `json:"credit_limit" db:"credit_limit"`
	IsActive *bool `json:"is_active" db:"is_active"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new customer_store_credit_accounts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CustomerStoreCreditAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "customer_store_credit_accounts", duration, nil)
	}()

	query := `
		INSERT INTO customer_store_credit_accounts (
			, organization_id
			, customer_id
			, current_balance
			, credit_limit
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
		entity.OrganizationId,
		entity.CustomerId,
		entity.CurrentBalance,
		entity.CreditLimit,
		entity.IsActive,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create customer_store_credit_accounts", zap.Error(err))
		return fmt.Errorf("failed to create customer_store_credit_accounts: %w", err)
	}

	r.logger.Info("created customer_store_credit_accounts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a customer_store_credit_accounts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CustomerStoreCreditAccounts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_store_credit_accounts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, current_balance
			, credit_limit
			, is_active
			, created_at
			, updated_at
			, deleted_at
		FROM customer_store_credit_accounts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CustomerStoreCreditAccounts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CustomerId,
		&entity.CurrentBalance,
		&entity.CreditLimit,
		&entity.IsActive,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("customer_store_credit_accounts not found")
	}

	if err != nil {
		r.logger.Error("failed to get customer_store_credit_accounts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get customer_store_credit_accounts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of customer_store_credit_accounts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CustomerStoreCreditAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_store_credit_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_store_credit_accounts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_store_credit_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, current_balance
			, credit_limit
			, is_active
			, created_at
			, updated_at
			, deleted_at
		FROM customer_store_credit_accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_store_credit_accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_store_credit_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerStoreCreditAccounts
	for rows.Next() {
		var entity CustomerStoreCreditAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.CurrentBalance,
			&entity.CreditLimit,
			&entity.IsActive,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_store_credit_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating customer_store_credit_accounts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing customer_store_credit_accounts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CustomerStoreCreditAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customer_store_credit_accounts", duration, nil)
	}()

	query := `
		UPDATE customer_store_credit_accounts
		SET
			, organization_id = $2
			, customer_id = $3
			, current_balance = $4
			, credit_limit = $5
			, is_active = $6
			, updated_at = $8
			, deleted_at = $9
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.CurrentBalance,
		entity.CreditLimit,
		entity.IsActive,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update customer_store_credit_accounts", zap.Error(err))
		return fmt.Errorf("failed to update customer_store_credit_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_store_credit_accounts not found or already deleted")
	}

	r.logger.Info("updated customer_store_credit_accounts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a customer_store_credit_accounts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customer_store_credit_accounts", duration, nil)
	}()

	query := `
		UPDATE customer_store_credit_accounts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete customer_store_credit_accounts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete customer_store_credit_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_store_credit_accounts not found or already deleted")
	}

	r.logger.Info("deleted customer_store_credit_accounts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves customer_store_credit_accounts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CustomerStoreCreditAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_store_credit_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_store_credit_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_store_credit_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, current_balance
			, credit_limit
			, is_active
			, created_at
			, updated_at
			, deleted_at
		FROM customer_store_credit_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_store_credit_accounts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_store_credit_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerStoreCreditAccounts
	for rows.Next() {
		var entity CustomerStoreCreditAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.CurrentBalance,
			&entity.CreditLimit,
			&entity.IsActive,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_store_credit_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


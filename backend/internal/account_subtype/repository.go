package account_subtype

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

// Repository handles database operations for AccountSubtypes
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AccountSubtypes repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AccountSubtypes represents a account_subtypes entity
type AccountSubtypes struct {
	Id *uuid.UUID `json:"id" db:"id"`
	AccountTypeId uuid.UUID `json:"account_type_id" db:"account_type_id"`
	SubtypeCode string `json:"subtype_code" db:"subtype_code"`
	SubtypeName string `json:"subtype_name" db:"subtype_name"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	Description *string `json:"description" db:"description"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new account_subtypes record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AccountSubtypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "account_subtypes", duration, nil)
	}()

	query := `
		INSERT INTO account_subtypes (
			, account_type_id
			, subtype_code
			, subtype_name
			, display_order
			, description
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.AccountTypeId,
		entity.SubtypeCode,
		entity.SubtypeName,
		entity.DisplayOrder,
		entity.Description,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create account_subtypes", zap.Error(err))
		return fmt.Errorf("failed to create account_subtypes: %w", err)
	}

	r.logger.Info("created account_subtypes",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a account_subtypes by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AccountSubtypes, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "account_subtypes", duration, nil)
	}()

	query := `
		SELECT
			id
			, account_type_id
			, subtype_code
			, subtype_name
			, display_order
			, description
			, created_at
			, updated_at
		FROM account_subtypes
		WHERE id = $1
		
	`

	var entity AccountSubtypes
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.AccountTypeId,
		&entity.SubtypeCode,
		&entity.SubtypeName,
		&entity.DisplayOrder,
		&entity.Description,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account_subtypes not found")
	}

	if err != nil {
		r.logger.Error("failed to get account_subtypes", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get account_subtypes: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of account_subtypes records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AccountSubtypes, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "account_subtypes", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM account_subtypes
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count account_subtypes records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, account_type_id
			, subtype_code
			, subtype_name
			, display_order
			, description
			, created_at
			, updated_at
		FROM account_subtypes
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list account_subtypes", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list account_subtypes: %w", err)
	}
	defer rows.Close()

	var entities []*AccountSubtypes
	for rows.Next() {
		var entity AccountSubtypes
		err := rows.Scan(
			&entity.Id,
			&entity.AccountTypeId,
			&entity.SubtypeCode,
			&entity.SubtypeName,
			&entity.DisplayOrder,
			&entity.Description,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan account_subtypes: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating account_subtypes rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing account_subtypes record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AccountSubtypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "account_subtypes", duration, nil)
	}()

	query := `
		UPDATE account_subtypes
		SET
			, account_type_id = $2
			, subtype_code = $3
			, subtype_name = $4
			, display_order = $5
			, description = $6
			, updated_at = $8
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.AccountTypeId,
		entity.SubtypeCode,
		entity.SubtypeName,
		entity.DisplayOrder,
		entity.Description,
		entity.UpdatedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update account_subtypes", zap.Error(err))
		return fmt.Errorf("failed to update account_subtypes: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("account_subtypes not found or already deleted")
	}

	r.logger.Info("updated account_subtypes",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a account_subtypes record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "account_subtypes", duration, nil)
	}()

	query := `DELETE FROM account_subtypes WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete account_subtypes", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete account_subtypes: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("account_subtypes not found")
	}

	r.logger.Info("deleted account_subtypes", zap.String("id", id.String()))
	return nil
}




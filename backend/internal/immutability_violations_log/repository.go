package immutability_violations_log

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

// Repository handles database operations for ImmutabilityViolationsLog
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ImmutabilityViolationsLog repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ImmutabilityViolationsLog represents a immutability_violations_log entity
type ImmutabilityViolationsLog struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	TableName string `json:"table_name" db:"table_name"`
	RecordId *uuid.UUID `json:"record_id" db:"record_id"`
	Operation string `json:"operation" db:"operation"`
	AttemptedBy *uuid.UUID `json:"attempted_by" db:"attempted_by"`
	AttemptedAt time.Time `json:"attempted_at" db:"attempted_at"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	BlockedData json.RawMessage `json:"blocked_data" db:"blocked_data"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
}

// Create inserts a new immutability_violations_log record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ImmutabilityViolationsLog) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "immutability_violations_log", duration, nil)
	}()

	query := `
		INSERT INTO immutability_violations_log (
			, organization_id
			, table_name
			, record_id
			, operation
			, attempted_by
			, attempted_at
			, error_message
			, blocked_data
			, metadata
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
		)
		RETURNING id
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.TableName,
		entity.RecordId,
		entity.Operation,
		entity.AttemptedBy,
		entity.AttemptedAt,
		entity.ErrorMessage,
		entity.BlockedData,
		entity.Metadata,
	)

	
	err := row.Scan(&entity.Id)
	

	if err != nil {
		r.logger.Error("failed to create immutability_violations_log", zap.Error(err))
		return fmt.Errorf("failed to create immutability_violations_log: %w", err)
	}

	r.logger.Info("created immutability_violations_log",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a immutability_violations_log by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ImmutabilityViolationsLog, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "immutability_violations_log", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, table_name
			, record_id
			, operation
			, attempted_by
			, attempted_at
			, error_message
			, blocked_data
			, metadata
		FROM immutability_violations_log
		WHERE id = $1
		
	`

	var entity ImmutabilityViolationsLog
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.TableName,
		&entity.RecordId,
		&entity.Operation,
		&entity.AttemptedBy,
		&entity.AttemptedAt,
		&entity.ErrorMessage,
		&entity.BlockedData,
		&entity.Metadata,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("immutability_violations_log not found")
	}

	if err != nil {
		r.logger.Error("failed to get immutability_violations_log", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get immutability_violations_log: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of immutability_violations_log records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ImmutabilityViolationsLog, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "immutability_violations_log", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM immutability_violations_log
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count immutability_violations_log records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, table_name
			, record_id
			, operation
			, attempted_by
			, attempted_at
			, error_message
			, blocked_data
			, metadata
		FROM immutability_violations_log
		
		
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list immutability_violations_log", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list immutability_violations_log: %w", err)
	}
	defer rows.Close()

	var entities []*ImmutabilityViolationsLog
	for rows.Next() {
		var entity ImmutabilityViolationsLog
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TableName,
			&entity.RecordId,
			&entity.Operation,
			&entity.AttemptedBy,
			&entity.AttemptedAt,
			&entity.ErrorMessage,
			&entity.BlockedData,
			&entity.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan immutability_violations_log: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating immutability_violations_log rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing immutability_violations_log record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ImmutabilityViolationsLog) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "immutability_violations_log", duration, nil)
	}()

	query := `
		UPDATE immutability_violations_log
		SET
			, organization_id = $2
			, table_name = $3
			, record_id = $4
			, operation = $5
			, attempted_by = $6
			, attempted_at = $7
			, error_message = $8
			, blocked_data = $9
			, metadata = $10
			
		WHERE id = $11
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.TableName,
		entity.RecordId,
		entity.Operation,
		entity.AttemptedBy,
		entity.AttemptedAt,
		entity.ErrorMessage,
		entity.BlockedData,
		entity.Metadata,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update immutability_violations_log", zap.Error(err))
		return fmt.Errorf("failed to update immutability_violations_log: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("immutability_violations_log not found or already deleted")
	}

	r.logger.Info("updated immutability_violations_log",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a immutability_violations_log record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "immutability_violations_log", duration, nil)
	}()

	query := `DELETE FROM immutability_violations_log WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete immutability_violations_log", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete immutability_violations_log: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("immutability_violations_log not found")
	}

	r.logger.Info("deleted immutability_violations_log", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves immutability_violations_log records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ImmutabilityViolationsLog, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "immutability_violations_log", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM immutability_violations_log
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count immutability_violations_log records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, table_name
			, record_id
			, operation
			, attempted_by
			, attempted_at
			, error_message
			, blocked_data
			, metadata
		FROM immutability_violations_log
		WHERE organization_id = $1
		
		
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list immutability_violations_log by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list immutability_violations_log: %w", err)
	}
	defer rows.Close()

	var entities []*ImmutabilityViolationsLog
	for rows.Next() {
		var entity ImmutabilityViolationsLog
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TableName,
			&entity.RecordId,
			&entity.Operation,
			&entity.AttemptedBy,
			&entity.AttemptedAt,
			&entity.ErrorMessage,
			&entity.BlockedData,
			&entity.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan immutability_violations_log: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


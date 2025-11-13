package pos_error_log

import (
	"encoding/json"
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

// Repository handles database operations for PosErrorLogs
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PosErrorLogs repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PosErrorLogs represents a pos_error_logs entity
type PosErrorLogs struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	ErrorLevel string `json:"error_level" db:"error_level"`
	ErrorCode *string `json:"error_code" db:"error_code"`
	ErrorMessage string `json:"error_message" db:"error_message"`
	DeviceId *uuid.UUID `json:"device_id" db:"device_id"`
	UserId *uuid.UUID `json:"user_id" db:"user_id"`
	PosSessionId *uuid.UUID `json:"pos_session_id" db:"pos_session_id"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	StackTrace *string `json:"stack_trace" db:"stack_trace"`
	RequestData json.RawMessage `json:"request_data" db:"request_data"`
	ErrorData json.RawMessage `json:"error_data" db:"error_data"`
	IsResolved *bool `json:"is_resolved" db:"is_resolved"`
	ResolvedBy *uuid.UUID `json:"resolved_by" db:"resolved_by"`
	ResolvedAt *time.Time `json:"resolved_at" db:"resolved_at"`
	ResolutionNotes *string `json:"resolution_notes" db:"resolution_notes"`
	OccurredAt *time.Time `json:"occurred_at" db:"occurred_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new pos_error_logs record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PosErrorLogs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "pos_error_logs", duration, nil)
	}()

	query := `
		INSERT INTO pos_error_logs (
			, organization_id
			, error_level
			, error_code
			, error_message
			, device_id
			, user_id
			, pos_session_id
			, sale_id
			, stack_trace
			, request_data
			, error_data
			, is_resolved
			, resolved_by
			, resolved_at
			, resolution_notes
			, occurred_at
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
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ErrorLevel,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.DeviceId,
		entity.UserId,
		entity.PosSessionId,
		entity.SaleId,
		entity.StackTrace,
		entity.RequestData,
		entity.ErrorData,
		entity.IsResolved,
		entity.ResolvedBy,
		entity.ResolvedAt,
		entity.ResolutionNotes,
		entity.OccurredAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create pos_error_logs", zap.Error(err))
		return fmt.Errorf("failed to create pos_error_logs: %w", err)
	}

	r.logger.Info("created pos_error_logs",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a pos_error_logs by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PosErrorLogs, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_error_logs", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, error_level
			, error_code
			, error_message
			, device_id
			, user_id
			, pos_session_id
			, sale_id
			, stack_trace
			, request_data
			, error_data
			, is_resolved
			, resolved_by
			, resolved_at
			, resolution_notes
			, occurred_at
			, created_at
		FROM pos_error_logs
		WHERE id = $1
		
	`

	var entity PosErrorLogs
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ErrorLevel,
		&entity.ErrorCode,
		&entity.ErrorMessage,
		&entity.DeviceId,
		&entity.UserId,
		&entity.PosSessionId,
		&entity.SaleId,
		&entity.StackTrace,
		&entity.RequestData,
		&entity.ErrorData,
		&entity.IsResolved,
		&entity.ResolvedBy,
		&entity.ResolvedAt,
		&entity.ResolutionNotes,
		&entity.OccurredAt,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("pos_error_logs not found")
	}

	if err != nil {
		r.logger.Error("failed to get pos_error_logs", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get pos_error_logs: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of pos_error_logs records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PosErrorLogs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_error_logs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_error_logs
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_error_logs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, error_level
			, error_code
			, error_message
			, device_id
			, user_id
			, pos_session_id
			, sale_id
			, stack_trace
			, request_data
			, error_data
			, is_resolved
			, resolved_by
			, resolved_at
			, resolution_notes
			, occurred_at
			, created_at
		FROM pos_error_logs
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_error_logs", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_error_logs: %w", err)
	}
	defer rows.Close()

	var entities []*PosErrorLogs
	for rows.Next() {
		var entity PosErrorLogs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ErrorLevel,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.DeviceId,
			&entity.UserId,
			&entity.PosSessionId,
			&entity.SaleId,
			&entity.StackTrace,
			&entity.RequestData,
			&entity.ErrorData,
			&entity.IsResolved,
			&entity.ResolvedBy,
			&entity.ResolvedAt,
			&entity.ResolutionNotes,
			&entity.OccurredAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_error_logs: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating pos_error_logs rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing pos_error_logs record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PosErrorLogs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_error_logs", duration, nil)
	}()

	query := `
		UPDATE pos_error_logs
		SET
			, organization_id = $2
			, error_level = $3
			, error_code = $4
			, error_message = $5
			, device_id = $6
			, user_id = $7
			, pos_session_id = $8
			, sale_id = $9
			, stack_trace = $10
			, request_data = $11
			, error_data = $12
			, is_resolved = $13
			, resolved_by = $14
			, resolved_at = $15
			, resolution_notes = $16
			, occurred_at = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ErrorLevel,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.DeviceId,
		entity.UserId,
		entity.PosSessionId,
		entity.SaleId,
		entity.StackTrace,
		entity.RequestData,
		entity.ErrorData,
		entity.IsResolved,
		entity.ResolvedBy,
		entity.ResolvedAt,
		entity.ResolutionNotes,
		entity.OccurredAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update pos_error_logs", zap.Error(err))
		return fmt.Errorf("failed to update pos_error_logs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_error_logs not found or already deleted")
	}

	r.logger.Info("updated pos_error_logs",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a pos_error_logs record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "pos_error_logs", duration, nil)
	}()

	query := `DELETE FROM pos_error_logs WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete pos_error_logs", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete pos_error_logs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_error_logs not found")
	}

	r.logger.Info("deleted pos_error_logs", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves pos_error_logs records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PosErrorLogs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_error_logs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_error_logs
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_error_logs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, error_level
			, error_code
			, error_message
			, device_id
			, user_id
			, pos_session_id
			, sale_id
			, stack_trace
			, request_data
			, error_data
			, is_resolved
			, resolved_by
			, resolved_at
			, resolution_notes
			, occurred_at
			, created_at
		FROM pos_error_logs
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_error_logs by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_error_logs: %w", err)
	}
	defer rows.Close()

	var entities []*PosErrorLogs
	for rows.Next() {
		var entity PosErrorLogs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ErrorLevel,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.DeviceId,
			&entity.UserId,
			&entity.PosSessionId,
			&entity.SaleId,
			&entity.StackTrace,
			&entity.RequestData,
			&entity.ErrorData,
			&entity.IsResolved,
			&entity.ResolvedBy,
			&entity.ResolvedAt,
			&entity.ResolutionNotes,
			&entity.OccurredAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_error_logs: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


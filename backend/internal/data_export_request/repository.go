package data_export_request

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

// Repository handles database operations for DataExportRequests
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DataExportRequests repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DataExportRequests represents a data_export_requests entity
type DataExportRequests struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ExportType string `json:"export_type" db:"export_type"`
	ExportFormat string `json:"export_format" db:"export_format"`
	DateFrom *time.Time `json:"date_from" db:"date_from"`
	DateTo *time.Time `json:"date_to" db:"date_to"`
	Filters json.RawMessage `json:"filters" db:"filters"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	FileName *string `json:"file_name" db:"file_name"`
	FileSize *int64 `json:"file_size" db:"file_size"`
	FilePath *string `json:"file_path" db:"file_path"`
	DownloadUrl *string `json:"download_url" db:"download_url"`
	DownloadExpiresAt *time.Time `json:"download_expires_at" db:"download_expires_at"`
	TotalRecords *int64 `json:"total_records" db:"total_records"`
	ProcessedRecords *int64 `json:"processed_records" db:"processed_records"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	RequestedBy *uuid.UUID `json:"requested_by" db:"requested_by"`
	RequestedAt *time.Time `json:"requested_at" db:"requested_at"`
	StartedAt *time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
}

// Create inserts a new data_export_requests record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DataExportRequests) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "data_export_requests", duration, nil)
	}()

	query := `
		INSERT INTO data_export_requests (
			, organization_id
			, export_type
			, export_format
			, date_from
			, date_to
			, filters
			, status
			, status
			, file_name
			, file_size
			, file_path
			, download_url
			, download_expires_at
			, total_records
			, processed_records
			, error_message
			, requested_by
			, requested_at
			, started_at
			, completed_at
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
			, $18
			, $19
			, $20
			, $21
		)
		RETURNING id
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ExportType,
		entity.ExportFormat,
		entity.DateFrom,
		entity.DateTo,
		entity.Filters,
		entity.Status,
		entity.Status,
		entity.FileName,
		entity.FileSize,
		entity.FilePath,
		entity.DownloadUrl,
		entity.DownloadExpiresAt,
		entity.TotalRecords,
		entity.ProcessedRecords,
		entity.ErrorMessage,
		entity.RequestedBy,
		entity.RequestedAt,
		entity.StartedAt,
		entity.CompletedAt,
	)

	
	err := row.Scan(&entity.Id)
	

	if err != nil {
		r.logger.Error("failed to create data_export_requests", zap.Error(err))
		return fmt.Errorf("failed to create data_export_requests: %w", err)
	}

	r.logger.Info("created data_export_requests",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a data_export_requests by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DataExportRequests, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "data_export_requests", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, export_type
			, export_format
			, date_from
			, date_to
			, filters
			, status
			, status
			, file_name
			, file_size
			, file_path
			, download_url
			, download_expires_at
			, total_records
			, processed_records
			, error_message
			, requested_by
			, requested_at
			, started_at
			, completed_at
		FROM data_export_requests
		WHERE id = $1
		
	`

	var entity DataExportRequests
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ExportType,
		&entity.ExportFormat,
		&entity.DateFrom,
		&entity.DateTo,
		&entity.Filters,
		&entity.Status,
		&entity.Status,
		&entity.FileName,
		&entity.FileSize,
		&entity.FilePath,
		&entity.DownloadUrl,
		&entity.DownloadExpiresAt,
		&entity.TotalRecords,
		&entity.ProcessedRecords,
		&entity.ErrorMessage,
		&entity.RequestedBy,
		&entity.RequestedAt,
		&entity.StartedAt,
		&entity.CompletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("data_export_requests not found")
	}

	if err != nil {
		r.logger.Error("failed to get data_export_requests", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get data_export_requests: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of data_export_requests records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DataExportRequests, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "data_export_requests", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM data_export_requests
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count data_export_requests records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, export_type
			, export_format
			, date_from
			, date_to
			, filters
			, status
			, status
			, file_name
			, file_size
			, file_path
			, download_url
			, download_expires_at
			, total_records
			, processed_records
			, error_message
			, requested_by
			, requested_at
			, started_at
			, completed_at
		FROM data_export_requests
		
		
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list data_export_requests", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list data_export_requests: %w", err)
	}
	defer rows.Close()

	var entities []*DataExportRequests
	for rows.Next() {
		var entity DataExportRequests
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ExportType,
			&entity.ExportFormat,
			&entity.DateFrom,
			&entity.DateTo,
			&entity.Filters,
			&entity.Status,
			&entity.Status,
			&entity.FileName,
			&entity.FileSize,
			&entity.FilePath,
			&entity.DownloadUrl,
			&entity.DownloadExpiresAt,
			&entity.TotalRecords,
			&entity.ProcessedRecords,
			&entity.ErrorMessage,
			&entity.RequestedBy,
			&entity.RequestedAt,
			&entity.StartedAt,
			&entity.CompletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan data_export_requests: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating data_export_requests rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing data_export_requests record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DataExportRequests) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "data_export_requests", duration, nil)
	}()

	query := `
		UPDATE data_export_requests
		SET
			, organization_id = $2
			, export_type = $3
			, export_format = $4
			, date_from = $5
			, date_to = $6
			, filters = $7
			, status = $8
			, status = $9
			, file_name = $10
			, file_size = $11
			, file_path = $12
			, download_url = $13
			, download_expires_at = $14
			, total_records = $15
			, processed_records = $16
			, error_message = $17
			, requested_by = $18
			, requested_at = $19
			, started_at = $20
			, completed_at = $21
			
		WHERE id = $22
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ExportType,
		entity.ExportFormat,
		entity.DateFrom,
		entity.DateTo,
		entity.Filters,
		entity.Status,
		entity.Status,
		entity.FileName,
		entity.FileSize,
		entity.FilePath,
		entity.DownloadUrl,
		entity.DownloadExpiresAt,
		entity.TotalRecords,
		entity.ProcessedRecords,
		entity.ErrorMessage,
		entity.RequestedBy,
		entity.RequestedAt,
		entity.StartedAt,
		entity.CompletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update data_export_requests", zap.Error(err))
		return fmt.Errorf("failed to update data_export_requests: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("data_export_requests not found or already deleted")
	}

	r.logger.Info("updated data_export_requests",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a data_export_requests record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "data_export_requests", duration, nil)
	}()

	query := `DELETE FROM data_export_requests WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete data_export_requests", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete data_export_requests: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("data_export_requests not found")
	}

	r.logger.Info("deleted data_export_requests", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves data_export_requests records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*DataExportRequests, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "data_export_requests", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM data_export_requests
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count data_export_requests records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, export_type
			, export_format
			, date_from
			, date_to
			, filters
			, status
			, status
			, file_name
			, file_size
			, file_path
			, download_url
			, download_expires_at
			, total_records
			, processed_records
			, error_message
			, requested_by
			, requested_at
			, started_at
			, completed_at
		FROM data_export_requests
		WHERE organization_id = $1
		
		
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list data_export_requests by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list data_export_requests: %w", err)
	}
	defer rows.Close()

	var entities []*DataExportRequests
	for rows.Next() {
		var entity DataExportRequests
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ExportType,
			&entity.ExportFormat,
			&entity.DateFrom,
			&entity.DateTo,
			&entity.Filters,
			&entity.Status,
			&entity.Status,
			&entity.FileName,
			&entity.FileSize,
			&entity.FilePath,
			&entity.DownloadUrl,
			&entity.DownloadExpiresAt,
			&entity.TotalRecords,
			&entity.ProcessedRecords,
			&entity.ErrorMessage,
			&entity.RequestedBy,
			&entity.RequestedAt,
			&entity.StartedAt,
			&entity.CompletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan data_export_requests: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


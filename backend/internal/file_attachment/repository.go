package file_attachment

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

// Repository handles database operations for FileAttachments
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new FileAttachments repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FileAttachments represents a file_attachments entity
type FileAttachments struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	FileName string `json:"file_name" db:"file_name"`
	FileSize int64 `json:"file_size" db:"file_size"`
	MimeType string `json:"mime_type" db:"mime_type"`
	FileExtension *string `json:"file_extension" db:"file_extension"`
	StorageProvider *string `json:"storage_provider" db:"storage_provider"`
	StoragePath string `json:"storage_path" db:"storage_path"`
	StorageUrl *string `json:"storage_url" db:"storage_url"`
	FileHash *string `json:"file_hash" db:"file_hash"`
	EntityType string `json:"entity_type" db:"entity_type"`
	EntityId uuid.UUID `json:"entity_id" db:"entity_id"`
	Description *string `json:"description" db:"description"`
	Tags *string `json:"tags" db:"tags"`
	IsPublic *bool `json:"is_public" db:"is_public"`
	ImageWidth *int64 `json:"image_width" db:"image_width"`
	ImageHeight *int64 `json:"image_height" db:"image_height"`
	VirusScanStatus *string `json:"virus_scan_status" db:"virus_scan_status"`
	VirusScanAt *time.Time `json:"virus_scan_at" db:"virus_scan_at"`
	UploadedBy *uuid.UUID `json:"uploaded_by" db:"uploaded_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new file_attachments record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *FileAttachments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "file_attachments", duration, nil)
	}()

	query := `
		INSERT INTO file_attachments (
			, organization_id
			, file_name
			, file_size
			, mime_type
			, file_extension
			, storage_provider
			, storage_path
			, storage_url
			, file_hash
			, entity_type
			, entity_id
			, description
			, tags
			, is_public
			, image_width
			, image_height
			, virus_scan_status
			, virus_scan_at
			, uploaded_by
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
			, $22
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.FileName,
		entity.FileSize,
		entity.MimeType,
		entity.FileExtension,
		entity.StorageProvider,
		entity.StoragePath,
		entity.StorageUrl,
		entity.FileHash,
		entity.EntityType,
		entity.EntityId,
		entity.Description,
		entity.Tags,
		entity.IsPublic,
		entity.ImageWidth,
		entity.ImageHeight,
		entity.VirusScanStatus,
		entity.VirusScanAt,
		entity.UploadedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create file_attachments", zap.Error(err))
		return fmt.Errorf("failed to create file_attachments: %w", err)
	}

	r.logger.Info("created file_attachments",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a file_attachments by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*FileAttachments, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "file_attachments", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, file_name
			, file_size
			, mime_type
			, file_extension
			, storage_provider
			, storage_path
			, storage_url
			, file_hash
			, entity_type
			, entity_id
			, description
			, tags
			, is_public
			, image_width
			, image_height
			, virus_scan_status
			, virus_scan_at
			, uploaded_by
			, created_at
			, deleted_at
		FROM file_attachments
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity FileAttachments
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.FileName,
		&entity.FileSize,
		&entity.MimeType,
		&entity.FileExtension,
		&entity.StorageProvider,
		&entity.StoragePath,
		&entity.StorageUrl,
		&entity.FileHash,
		&entity.EntityType,
		&entity.EntityId,
		&entity.Description,
		&entity.Tags,
		&entity.IsPublic,
		&entity.ImageWidth,
		&entity.ImageHeight,
		&entity.VirusScanStatus,
		&entity.VirusScanAt,
		&entity.UploadedBy,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("file_attachments not found")
	}

	if err != nil {
		r.logger.Error("failed to get file_attachments", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get file_attachments: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of file_attachments records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*FileAttachments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "file_attachments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM file_attachments
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count file_attachments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, file_name
			, file_size
			, mime_type
			, file_extension
			, storage_provider
			, storage_path
			, storage_url
			, file_hash
			, entity_type
			, entity_id
			, description
			, tags
			, is_public
			, image_width
			, image_height
			, virus_scan_status
			, virus_scan_at
			, uploaded_by
			, created_at
			, deleted_at
		FROM file_attachments
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list file_attachments", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list file_attachments: %w", err)
	}
	defer rows.Close()

	var entities []*FileAttachments
	for rows.Next() {
		var entity FileAttachments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FileName,
			&entity.FileSize,
			&entity.MimeType,
			&entity.FileExtension,
			&entity.StorageProvider,
			&entity.StoragePath,
			&entity.StorageUrl,
			&entity.FileHash,
			&entity.EntityType,
			&entity.EntityId,
			&entity.Description,
			&entity.Tags,
			&entity.IsPublic,
			&entity.ImageWidth,
			&entity.ImageHeight,
			&entity.VirusScanStatus,
			&entity.VirusScanAt,
			&entity.UploadedBy,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan file_attachments: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating file_attachments rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing file_attachments record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *FileAttachments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "file_attachments", duration, nil)
	}()

	query := `
		UPDATE file_attachments
		SET
			, organization_id = $2
			, file_name = $3
			, file_size = $4
			, mime_type = $5
			, file_extension = $6
			, storage_provider = $7
			, storage_path = $8
			, storage_url = $9
			, file_hash = $10
			, entity_type = $11
			, entity_id = $12
			, description = $13
			, tags = $14
			, is_public = $15
			, image_width = $16
			, image_height = $17
			, virus_scan_status = $18
			, virus_scan_at = $19
			, uploaded_by = $20
			, deleted_at = $22
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $23
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.FileName,
		entity.FileSize,
		entity.MimeType,
		entity.FileExtension,
		entity.StorageProvider,
		entity.StoragePath,
		entity.StorageUrl,
		entity.FileHash,
		entity.EntityType,
		entity.EntityId,
		entity.Description,
		entity.Tags,
		entity.IsPublic,
		entity.ImageWidth,
		entity.ImageHeight,
		entity.VirusScanStatus,
		entity.VirusScanAt,
		entity.UploadedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update file_attachments", zap.Error(err))
		return fmt.Errorf("failed to update file_attachments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("file_attachments not found or already deleted")
	}

	r.logger.Info("updated file_attachments",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a file_attachments record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "file_attachments", duration, nil)
	}()

	query := `
		UPDATE file_attachments
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete file_attachments", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete file_attachments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("file_attachments not found or already deleted")
	}

	r.logger.Info("deleted file_attachments", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves file_attachments records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*FileAttachments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "file_attachments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM file_attachments
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count file_attachments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, file_name
			, file_size
			, mime_type
			, file_extension
			, storage_provider
			, storage_path
			, storage_url
			, file_hash
			, entity_type
			, entity_id
			, description
			, tags
			, is_public
			, image_width
			, image_height
			, virus_scan_status
			, virus_scan_at
			, uploaded_by
			, created_at
			, deleted_at
		FROM file_attachments
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list file_attachments by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list file_attachments: %w", err)
	}
	defer rows.Close()

	var entities []*FileAttachments
	for rows.Next() {
		var entity FileAttachments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FileName,
			&entity.FileSize,
			&entity.MimeType,
			&entity.FileExtension,
			&entity.StorageProvider,
			&entity.StoragePath,
			&entity.StorageUrl,
			&entity.FileHash,
			&entity.EntityType,
			&entity.EntityId,
			&entity.Description,
			&entity.Tags,
			&entity.IsPublic,
			&entity.ImageWidth,
			&entity.ImageHeight,
			&entity.VirusScanStatus,
			&entity.VirusScanAt,
			&entity.UploadedBy,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan file_attachments: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


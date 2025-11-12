package printer_configuration

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

// Repository handles database operations for PrinterConfigurations
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PrinterConfigurations repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PrinterConfigurations represents a printer_configurations entity
type PrinterConfigurations struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	PrinterDeviceId uuid.UUID `json:"printer_device_id" db:"printer_device_id"`
	DocumentType string `json:"document_type" db:"document_type"`
	FilterOrderType *string `json:"filter_order_type" db:"filter_order_type"`
	FilterKitchenStationId *uuid.UUID `json:"filter_kitchen_station_id" db:"filter_kitchen_station_id"`
	FilterProductCategoryId *uuid.UUID `json:"filter_product_category_id" db:"filter_product_category_id"`
	FilterCourseId *uuid.UUID `json:"filter_course_id" db:"filter_course_id"`
	NumberOfCopies *int64 `json:"number_of_copies" db:"number_of_copies"`
	AutoPrint *bool `json:"auto_print" db:"auto_print"`
	PrintPriority *int64 `json:"print_priority" db:"print_priority"`
	TemplateConfig json.RawMessage `json:"template_config" db:"template_config"`
	PaperSize *string `json:"paper_size" db:"paper_size"`
	PrintOrientation *string `json:"print_orientation" db:"print_orientation"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'receipt', *string `json:"'receipt'," db:"'receipt',"`
	'report', *string `json:"'report'," db:"'report',"`
}

// Create inserts a new printer_configurations record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PrinterConfigurations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "printer_configurations", duration, nil)
	}()

	query := `
		INSERT INTO printer_configurations (
			, organization_id
			, location_id
			, printer_device_id
			, document_type
			, filter_order_type
			, filter_kitchen_station_id
			, filter_product_category_id
			, filter_course_id
			, number_of_copies
			, auto_print
			, print_priority
			, template_config
			, paper_size
			, print_orientation
			, is_active
			, notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'receipt',
			, 'report',
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
			, $21
			, $22
			, $23
			, $24
			, $25
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.PrinterDeviceId,
		entity.DocumentType,
		entity.FilterOrderType,
		entity.FilterKitchenStationId,
		entity.FilterProductCategoryId,
		entity.FilterCourseId,
		entity.NumberOfCopies,
		entity.AutoPrint,
		entity.PrintPriority,
		entity.TemplateConfig,
		entity.PaperSize,
		entity.PrintOrientation,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'receipt',,
		entity.'report',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create printer_configurations", zap.Error(err))
		return fmt.Errorf("failed to create printer_configurations: %w", err)
	}

	r.logger.Info("created printer_configurations",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a printer_configurations by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PrinterConfigurations, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "printer_configurations", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, printer_device_id
			, document_type
			, filter_order_type
			, filter_kitchen_station_id
			, filter_product_category_id
			, filter_course_id
			, number_of_copies
			, auto_print
			, print_priority
			, template_config
			, paper_size
			, print_orientation
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'receipt',
			, 'report',
		FROM printer_configurations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PrinterConfigurations
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.PrinterDeviceId,
		&entity.DocumentType,
		&entity.FilterOrderType,
		&entity.FilterKitchenStationId,
		&entity.FilterProductCategoryId,
		&entity.FilterCourseId,
		&entity.NumberOfCopies,
		&entity.AutoPrint,
		&entity.PrintPriority,
		&entity.TemplateConfig,
		&entity.PaperSize,
		&entity.PrintOrientation,
		&entity.IsActive,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'receipt',,
		&entity.'report',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("printer_configurations not found")
	}

	if err != nil {
		r.logger.Error("failed to get printer_configurations", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get printer_configurations: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of printer_configurations records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PrinterConfigurations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "printer_configurations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM printer_configurations
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count printer_configurations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, printer_device_id
			, document_type
			, filter_order_type
			, filter_kitchen_station_id
			, filter_product_category_id
			, filter_course_id
			, number_of_copies
			, auto_print
			, print_priority
			, template_config
			, paper_size
			, print_orientation
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'receipt',
			, 'report',
		FROM printer_configurations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list printer_configurations", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list printer_configurations: %w", err)
	}
	defer rows.Close()

	var entities []*PrinterConfigurations
	for rows.Next() {
		var entity PrinterConfigurations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.PrinterDeviceId,
			&entity.DocumentType,
			&entity.FilterOrderType,
			&entity.FilterKitchenStationId,
			&entity.FilterProductCategoryId,
			&entity.FilterCourseId,
			&entity.NumberOfCopies,
			&entity.AutoPrint,
			&entity.PrintPriority,
			&entity.TemplateConfig,
			&entity.PaperSize,
			&entity.PrintOrientation,
			&entity.IsActive,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'receipt',,
			&entity.'report',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan printer_configurations: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating printer_configurations rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing printer_configurations record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PrinterConfigurations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "printer_configurations", duration, nil)
	}()

	query := `
		UPDATE printer_configurations
		SET
			, organization_id = $2
			, location_id = $3
			, printer_device_id = $4
			, document_type = $5
			, filter_order_type = $6
			, filter_kitchen_station_id = $7
			, filter_product_category_id = $8
			, filter_course_id = $9
			, number_of_copies = $10
			, auto_print = $11
			, print_priority = $12
			, template_config = $13
			, paper_size = $14
			, print_orientation = $15
			, is_active = $16
			, notes = $17
			, metadata = $18
			, updated_at = $20
			, created_by = $21
			, updated_by = $22
			, deleted_at = $23
			, 'receipt', = $24
			, 'report', = $25
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $26
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.PrinterDeviceId,
		entity.DocumentType,
		entity.FilterOrderType,
		entity.FilterKitchenStationId,
		entity.FilterProductCategoryId,
		entity.FilterCourseId,
		entity.NumberOfCopies,
		entity.AutoPrint,
		entity.PrintPriority,
		entity.TemplateConfig,
		entity.PaperSize,
		entity.PrintOrientation,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'receipt',,
		entity.'report',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update printer_configurations", zap.Error(err))
		return fmt.Errorf("failed to update printer_configurations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("printer_configurations not found or already deleted")
	}

	r.logger.Info("updated printer_configurations",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a printer_configurations record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "printer_configurations", duration, nil)
	}()

	query := `
		UPDATE printer_configurations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete printer_configurations", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete printer_configurations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("printer_configurations not found or already deleted")
	}

	r.logger.Info("deleted printer_configurations", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves printer_configurations records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PrinterConfigurations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "printer_configurations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM printer_configurations
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count printer_configurations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, printer_device_id
			, document_type
			, filter_order_type
			, filter_kitchen_station_id
			, filter_product_category_id
			, filter_course_id
			, number_of_copies
			, auto_print
			, print_priority
			, template_config
			, paper_size
			, print_orientation
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'receipt',
			, 'report',
		FROM printer_configurations
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list printer_configurations by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list printer_configurations: %w", err)
	}
	defer rows.Close()

	var entities []*PrinterConfigurations
	for rows.Next() {
		var entity PrinterConfigurations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.PrinterDeviceId,
			&entity.DocumentType,
			&entity.FilterOrderType,
			&entity.FilterKitchenStationId,
			&entity.FilterProductCategoryId,
			&entity.FilterCourseId,
			&entity.NumberOfCopies,
			&entity.AutoPrint,
			&entity.PrintPriority,
			&entity.TemplateConfig,
			&entity.PaperSize,
			&entity.PrintOrientation,
			&entity.IsActive,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'receipt',,
			&entity.'report',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan printer_configurations: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


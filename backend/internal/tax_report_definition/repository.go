package tax_report_definition

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

// Repository handles database operations for TaxReportDefinitions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new TaxReportDefinitions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// TaxReportDefinitions represents a tax_report_definitions entity
type TaxReportDefinitions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	LocalizationPackageId *uuid.UUID `json:"localization_package_id" db:"localization_package_id"`
	ReportCode string `json:"report_code" db:"report_code"`
	ReportName string `json:"report_name" db:"report_name"`
	Jurisdiction *string `json:"jurisdiction" db:"jurisdiction"`
	Authority *string `json:"authority" db:"authority"`
	ReportFrequency *string `json:"report_frequency" db:"report_frequency"`
	Version *string `json:"version" db:"version"`
	EffectiveFrom *time.Time `json:"effective_from" db:"effective_from"`
	EffectiveTo *time.Time `json:"effective_to" db:"effective_to"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Description *string `json:"description" db:"description"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new tax_report_definitions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *TaxReportDefinitions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "tax_report_definitions", duration, nil)
	}()

	query := `
		INSERT INTO tax_report_definitions (
			, organization_id
			, localization_package_id
			, report_code
			, report_name
			, jurisdiction
			, authority
			, report_frequency
			, version
			, effective_from
			, effective_to
			, is_active
			, description
			, created_by
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
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocalizationPackageId,
		entity.ReportCode,
		entity.ReportName,
		entity.Jurisdiction,
		entity.Authority,
		entity.ReportFrequency,
		entity.Version,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create tax_report_definitions", zap.Error(err))
		return fmt.Errorf("failed to create tax_report_definitions: %w", err)
	}

	r.logger.Info("created tax_report_definitions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a tax_report_definitions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*TaxReportDefinitions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_report_definitions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, localization_package_id
			, report_code
			, report_name
			, jurisdiction
			, authority
			, report_frequency
			, version
			, effective_from
			, effective_to
			, is_active
			, description
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM tax_report_definitions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity TaxReportDefinitions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocalizationPackageId,
		&entity.ReportCode,
		&entity.ReportName,
		&entity.Jurisdiction,
		&entity.Authority,
		&entity.ReportFrequency,
		&entity.Version,
		&entity.EffectiveFrom,
		&entity.EffectiveTo,
		&entity.IsActive,
		&entity.Description,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("tax_report_definitions not found")
	}

	if err != nil {
		r.logger.Error("failed to get tax_report_definitions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get tax_report_definitions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of tax_report_definitions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*TaxReportDefinitions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_report_definitions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tax_report_definitions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tax_report_definitions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, localization_package_id
			, report_code
			, report_name
			, jurisdiction
			, authority
			, report_frequency
			, version
			, effective_from
			, effective_to
			, is_active
			, description
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM tax_report_definitions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tax_report_definitions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tax_report_definitions: %w", err)
	}
	defer rows.Close()

	var entities []*TaxReportDefinitions
	for rows.Next() {
		var entity TaxReportDefinitions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocalizationPackageId,
			&entity.ReportCode,
			&entity.ReportName,
			&entity.Jurisdiction,
			&entity.Authority,
			&entity.ReportFrequency,
			&entity.Version,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tax_report_definitions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating tax_report_definitions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing tax_report_definitions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *TaxReportDefinitions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tax_report_definitions", duration, nil)
	}()

	query := `
		UPDATE tax_report_definitions
		SET
			, organization_id = $2
			, localization_package_id = $3
			, report_code = $4
			, report_name = $5
			, jurisdiction = $6
			, authority = $7
			, report_frequency = $8
			, version = $9
			, effective_from = $10
			, effective_to = $11
			, is_active = $12
			, description = $13
			, created_by = $14
			, updated_at = $16
			, deleted_at = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocalizationPackageId,
		entity.ReportCode,
		entity.ReportName,
		entity.Jurisdiction,
		entity.Authority,
		entity.ReportFrequency,
		entity.Version,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update tax_report_definitions", zap.Error(err))
		return fmt.Errorf("failed to update tax_report_definitions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tax_report_definitions not found or already deleted")
	}

	r.logger.Info("updated tax_report_definitions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a tax_report_definitions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tax_report_definitions", duration, nil)
	}()

	query := `
		UPDATE tax_report_definitions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete tax_report_definitions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete tax_report_definitions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tax_report_definitions not found or already deleted")
	}

	r.logger.Info("deleted tax_report_definitions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves tax_report_definitions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*TaxReportDefinitions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_report_definitions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tax_report_definitions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tax_report_definitions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, localization_package_id
			, report_code
			, report_name
			, jurisdiction
			, authority
			, report_frequency
			, version
			, effective_from
			, effective_to
			, is_active
			, description
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM tax_report_definitions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tax_report_definitions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tax_report_definitions: %w", err)
	}
	defer rows.Close()

	var entities []*TaxReportDefinitions
	for rows.Next() {
		var entity TaxReportDefinitions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocalizationPackageId,
			&entity.ReportCode,
			&entity.ReportName,
			&entity.Jurisdiction,
			&entity.Authority,
			&entity.ReportFrequency,
			&entity.Version,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tax_report_definitions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


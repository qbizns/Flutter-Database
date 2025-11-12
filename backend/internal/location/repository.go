package location

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

// Repository handles database operations for Locations
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Locations repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Locations represents a locations entity
type Locations struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationCode string `json:"location_code" db:"location_code"`
	Name string `json:"name" db:"name"`
	LocationType *string `json:"location_type" db:"location_type"`
	Phone *string `json:"phone" db:"phone"`
	Email *string `json:"email" db:"email"`
	ManagerUserId *uuid.UUID `json:"manager_user_id" db:"manager_user_id"`
	AddressLine1 *string `json:"address_line1" db:"address_line1"`
	AddressLine2 *string `json:"address_line2" db:"address_line2"`
	City *string `json:"city" db:"city"`
	State *string `json:"state" db:"state"`
	Country *string `json:"country" db:"country"`
	PostalCode *string `json:"postal_code" db:"postal_code"`
	Timezone *string `json:"timezone" db:"timezone"`
	BusinessHours json.RawMessage `json:"business_hours" db:"business_hours"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsPrimary *bool `json:"is_primary" db:"is_primary"`
	AllowSales *bool `json:"allow_sales" db:"allow_sales"`
	AllowPurchases *bool `json:"allow_purchases" db:"allow_purchases"`
	TaxRate *float64 `json:"tax_rate" db:"tax_rate"`
	Notes *string `json:"notes" db:"notes"`
	Settings json.RawMessage `json:"settings" db:"settings"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new locations record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Locations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "locations", duration, nil)
	}()

	query := `
		INSERT INTO locations (
			, organization_id
			, location_code
			, name
			, location_type
			, phone
			, email
			, manager_user_id
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, timezone
			, business_hours
			, is_active
			, is_primary
			, allow_sales
			, allow_purchases
			, tax_rate
			, notes
			, settings
			, metadata
			, deleted_at
			, created_by
			, updated_by
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
			, $22
			, $23
			, $24
			, $27
			, $28
			, $29
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationCode,
		entity.Name,
		entity.LocationType,
		entity.Phone,
		entity.Email,
		entity.ManagerUserId,
		entity.AddressLine1,
		entity.AddressLine2,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.Timezone,
		entity.BusinessHours,
		entity.IsActive,
		entity.IsPrimary,
		entity.AllowSales,
		entity.AllowPurchases,
		entity.TaxRate,
		entity.Notes,
		entity.Settings,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create locations", zap.Error(err))
		return fmt.Errorf("failed to create locations: %w", err)
	}

	r.logger.Info("created locations",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a locations by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Locations, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "locations", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_code
			, name
			, location_type
			, phone
			, email
			, manager_user_id
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, timezone
			, business_hours
			, is_active
			, is_primary
			, allow_sales
			, allow_purchases
			, tax_rate
			, notes
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM locations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Locations
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationCode,
		&entity.Name,
		&entity.LocationType,
		&entity.Phone,
		&entity.Email,
		&entity.ManagerUserId,
		&entity.AddressLine1,
		&entity.AddressLine2,
		&entity.City,
		&entity.State,
		&entity.Country,
		&entity.PostalCode,
		&entity.Timezone,
		&entity.BusinessHours,
		&entity.IsActive,
		&entity.IsPrimary,
		&entity.AllowSales,
		&entity.AllowPurchases,
		&entity.TaxRate,
		&entity.Notes,
		&entity.Settings,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("locations not found")
	}

	if err != nil {
		r.logger.Error("failed to get locations", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of locations records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Locations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "locations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM locations
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count locations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_code
			, name
			, location_type
			, phone
			, email
			, manager_user_id
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, timezone
			, business_hours
			, is_active
			, is_primary
			, allow_sales
			, allow_purchases
			, tax_rate
			, notes
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM locations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list locations", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list locations: %w", err)
	}
	defer rows.Close()

	var entities []*Locations
	for rows.Next() {
		var entity Locations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationCode,
			&entity.Name,
			&entity.LocationType,
			&entity.Phone,
			&entity.Email,
			&entity.ManagerUserId,
			&entity.AddressLine1,
			&entity.AddressLine2,
			&entity.City,
			&entity.State,
			&entity.Country,
			&entity.PostalCode,
			&entity.Timezone,
			&entity.BusinessHours,
			&entity.IsActive,
			&entity.IsPrimary,
			&entity.AllowSales,
			&entity.AllowPurchases,
			&entity.TaxRate,
			&entity.Notes,
			&entity.Settings,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan locations: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating locations rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing locations record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Locations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "locations", duration, nil)
	}()

	query := `
		UPDATE locations
		SET
			, organization_id = $2
			, location_code = $3
			, name = $4
			, location_type = $5
			, phone = $6
			, email = $7
			, manager_user_id = $8
			, address_line1 = $9
			, address_line2 = $10
			, city = $11
			, state = $12
			, country = $13
			, postal_code = $14
			, timezone = $15
			, business_hours = $16
			, is_active = $17
			, is_primary = $18
			, allow_sales = $19
			, allow_purchases = $20
			, tax_rate = $21
			, notes = $22
			, settings = $23
			, metadata = $24
			, updated_at = $26
			, deleted_at = $27
			, created_by = $28
			, updated_by = $29
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $30
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationCode,
		entity.Name,
		entity.LocationType,
		entity.Phone,
		entity.Email,
		entity.ManagerUserId,
		entity.AddressLine1,
		entity.AddressLine2,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.Timezone,
		entity.BusinessHours,
		entity.IsActive,
		entity.IsPrimary,
		entity.AllowSales,
		entity.AllowPurchases,
		entity.TaxRate,
		entity.Notes,
		entity.Settings,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update locations", zap.Error(err))
		return fmt.Errorf("failed to update locations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("locations not found or already deleted")
	}

	r.logger.Info("updated locations",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a locations record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "locations", duration, nil)
	}()

	query := `
		UPDATE locations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete locations", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete locations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("locations not found or already deleted")
	}

	r.logger.Info("deleted locations", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves locations records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Locations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "locations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM locations
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count locations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_code
			, name
			, location_type
			, phone
			, email
			, manager_user_id
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, timezone
			, business_hours
			, is_active
			, is_primary
			, allow_sales
			, allow_purchases
			, tax_rate
			, notes
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM locations
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list locations by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list locations: %w", err)
	}
	defer rows.Close()

	var entities []*Locations
	for rows.Next() {
		var entity Locations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationCode,
			&entity.Name,
			&entity.LocationType,
			&entity.Phone,
			&entity.Email,
			&entity.ManagerUserId,
			&entity.AddressLine1,
			&entity.AddressLine2,
			&entity.City,
			&entity.State,
			&entity.Country,
			&entity.PostalCode,
			&entity.Timezone,
			&entity.BusinessHours,
			&entity.IsActive,
			&entity.IsPrimary,
			&entity.AllowSales,
			&entity.AllowPurchases,
			&entity.TaxRate,
			&entity.Notes,
			&entity.Settings,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan locations: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


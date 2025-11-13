package customer_address

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

// Repository handles database operations for CustomerAddresses
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CustomerAddresses repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CustomerAddresses represents a customer_addresses entity
type CustomerAddresses struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CustomerId uuid.UUID `json:"customer_id" db:"customer_id"`
	AddressLabel *string `json:"address_label" db:"address_label"`
	AddressLine1 string `json:"address_line1" db:"address_line1"`
	AddressLine2 *string `json:"address_line2" db:"address_line2"`
	City *string `json:"city" db:"city"`
	StateProvince *string `json:"state_province" db:"state_province"`
	PostalCode *string `json:"postal_code" db:"postal_code"`
	Country *string `json:"country" db:"country"`
	Latitude *float64 `json:"latitude" db:"latitude"`
	Longitude *float64 `json:"longitude" db:"longitude"`
	LocationNotes *string `json:"location_notes" db:"location_notes"`
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id" db:"delivery_zone_id"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new customer_addresses record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CustomerAddresses) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "customer_addresses", duration, nil)
	}()

	query := `
		INSERT INTO customer_addresses (
			, organization_id
			, customer_id
			, address_label
			, address_line1
			, address_line2
			, city
			, state_province
			, postal_code
			, country
			, latitude
			, longitude
			, location_notes
			, delivery_zone_id
			, is_default
			, is_active
			, metadata
			, created_by
			, updated_by
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
			, $20
			, $21
			, $22
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.AddressLabel,
		entity.AddressLine1,
		entity.AddressLine2,
		entity.City,
		entity.StateProvince,
		entity.PostalCode,
		entity.Country,
		entity.Latitude,
		entity.Longitude,
		entity.LocationNotes,
		entity.DeliveryZoneId,
		entity.IsDefault,
		entity.IsActive,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create customer_addresses", zap.Error(err))
		return fmt.Errorf("failed to create customer_addresses: %w", err)
	}

	r.logger.Info("created customer_addresses",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a customer_addresses by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CustomerAddresses, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_addresses", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, address_label
			, address_line1
			, address_line2
			, city
			, state_province
			, postal_code
			, country
			, latitude
			, longitude
			, location_notes
			, delivery_zone_id
			, is_default
			, is_active
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM customer_addresses
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CustomerAddresses
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CustomerId,
		&entity.AddressLabel,
		&entity.AddressLine1,
		&entity.AddressLine2,
		&entity.City,
		&entity.StateProvince,
		&entity.PostalCode,
		&entity.Country,
		&entity.Latitude,
		&entity.Longitude,
		&entity.LocationNotes,
		&entity.DeliveryZoneId,
		&entity.IsDefault,
		&entity.IsActive,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("customer_addresses not found")
	}

	if err != nil {
		r.logger.Error("failed to get customer_addresses", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get customer_addresses: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of customer_addresses records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CustomerAddresses, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_addresses", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_addresses
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_addresses records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, address_label
			, address_line1
			, address_line2
			, city
			, state_province
			, postal_code
			, country
			, latitude
			, longitude
			, location_notes
			, delivery_zone_id
			, is_default
			, is_active
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM customer_addresses
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_addresses", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_addresses: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerAddresses
	for rows.Next() {
		var entity CustomerAddresses
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.AddressLabel,
			&entity.AddressLine1,
			&entity.AddressLine2,
			&entity.City,
			&entity.StateProvince,
			&entity.PostalCode,
			&entity.Country,
			&entity.Latitude,
			&entity.Longitude,
			&entity.LocationNotes,
			&entity.DeliveryZoneId,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_addresses: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating customer_addresses rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing customer_addresses record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CustomerAddresses) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customer_addresses", duration, nil)
	}()

	query := `
		UPDATE customer_addresses
		SET
			, organization_id = $2
			, customer_id = $3
			, address_label = $4
			, address_line1 = $5
			, address_line2 = $6
			, city = $7
			, state_province = $8
			, postal_code = $9
			, country = $10
			, latitude = $11
			, longitude = $12
			, location_notes = $13
			, delivery_zone_id = $14
			, is_default = $15
			, is_active = $16
			, metadata = $17
			, updated_at = $19
			, created_by = $20
			, updated_by = $21
			, deleted_at = $22
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $23
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.AddressLabel,
		entity.AddressLine1,
		entity.AddressLine2,
		entity.City,
		entity.StateProvince,
		entity.PostalCode,
		entity.Country,
		entity.Latitude,
		entity.Longitude,
		entity.LocationNotes,
		entity.DeliveryZoneId,
		entity.IsDefault,
		entity.IsActive,
		entity.Metadata,
		time.Now(),
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update customer_addresses", zap.Error(err))
		return fmt.Errorf("failed to update customer_addresses: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_addresses not found or already deleted")
	}

	r.logger.Info("updated customer_addresses",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a customer_addresses record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customer_addresses", duration, nil)
	}()

	query := `
		UPDATE customer_addresses
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete customer_addresses", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete customer_addresses: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_addresses not found or already deleted")
	}

	r.logger.Info("deleted customer_addresses", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves customer_addresses records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CustomerAddresses, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_addresses", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_addresses
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_addresses records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, address_label
			, address_line1
			, address_line2
			, city
			, state_province
			, postal_code
			, country
			, latitude
			, longitude
			, location_notes
			, delivery_zone_id
			, is_default
			, is_active
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM customer_addresses
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_addresses by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_addresses: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerAddresses
	for rows.Next() {
		var entity CustomerAddresses
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.AddressLabel,
			&entity.AddressLine1,
			&entity.AddressLine2,
			&entity.City,
			&entity.StateProvince,
			&entity.PostalCode,
			&entity.Country,
			&entity.Latitude,
			&entity.Longitude,
			&entity.LocationNotes,
			&entity.DeliveryZoneId,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_addresses: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/locations"
)

type LocationRepository struct {
	db *DB
}

func NewLocationRepository(db *DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) List(ctx context.Context, orgID uuid.UUID, filters locations.LocationFilters) ([]locations.Location, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, location_code, name, location_type,
			phone, email, manager_user_id,
			address_line1, address_line2, city, state, country, postal_code,
			timezone, business_hours, is_active, is_primary,
			allow_sales, allow_purchases, tax_rate,
			notes, settings, metadata,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM locations
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR location_code ILIKE $%d OR city ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.LocationType != nil {
		argCount++
		query += fmt.Sprintf(" AND location_type = $%d", argCount)
		args = append(args, *filters.LocationType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.IsPrimary != nil {
		argCount++
		query += fmt.Sprintf(" AND is_primary = $%d", argCount)
		args = append(args, *filters.IsPrimary)
	}

	if filters.City != nil {
		argCount++
		query += fmt.Sprintf(" AND city = $%d", argCount)
		args = append(args, *filters.City)
	}

	if filters.State != nil {
		argCount++
		query += fmt.Sprintf(" AND state = $%d", argCount)
		args = append(args, *filters.State)
	}

	if filters.Country != nil {
		argCount++
		query += fmt.Sprintf(" AND country = $%d", argCount)
		args = append(args, *filters.Country)
	}

	query += " ORDER BY created_at DESC"

	// Pagination
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locationsList []locations.Location
	for rows.Next() {
		var loc locations.Location
		var businessHours []byte
		var settings []byte
		var metadata []byte

		err := rows.Scan(
			&loc.ID, &loc.OrganizationID, &loc.LocationCode, &loc.Name, &loc.LocationType,
			&loc.Phone, &loc.Email, &loc.ManagerUserID,
			&loc.AddressLine1, &loc.AddressLine2, &loc.City, &loc.State, &loc.Country, &loc.PostalCode,
			&loc.Timezone, &businessHours, &loc.IsActive, &loc.IsPrimary,
			&loc.AllowSales, &loc.AllowPurchases, &loc.TaxRate,
			&loc.Notes, &settings, &metadata,
			&loc.CreatedAt, &loc.UpdatedAt, &loc.CreatedBy, &loc.UpdatedBy, &loc.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB fields
		if len(businessHours) > 0 {
			json.Unmarshal(businessHours, &loc.BusinessHours)
		}
		if len(settings) > 0 {
			json.Unmarshal(settings, &loc.Settings)
		}
		if len(metadata) > 0 {
			json.Unmarshal(metadata, &loc.Metadata)
		}

		locationsList = append(locationsList, loc)
	}

	return locationsList, rows.Err()
}

func (r *LocationRepository) Count(ctx context.Context, orgID uuid.UUID, filters locations.LocationFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM locations WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	// Apply same filters as List
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR location_code ILIKE $%d OR city ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.LocationType != nil {
		argCount++
		query += fmt.Sprintf(" AND location_type = $%d", argCount)
		args = append(args, *filters.LocationType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.IsPrimary != nil {
		argCount++
		query += fmt.Sprintf(" AND is_primary = $%d", argCount)
		args = append(args, *filters.IsPrimary)
	}

	if filters.City != nil {
		argCount++
		query += fmt.Sprintf(" AND city = $%d", argCount)
		args = append(args, *filters.City)
	}

	if filters.State != nil {
		argCount++
		query += fmt.Sprintf(" AND state = $%d", argCount)
		args = append(args, *filters.State)
	}

	if filters.Country != nil {
		argCount++
		query += fmt.Sprintf(" AND country = $%d", argCount)
		args = append(args, *filters.Country)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *LocationRepository) Create(ctx context.Context, location *locations.Location) error {
	if err := r.db.SetOrganizationContext(ctx, location.OrganizationID.String()); err != nil {
		return err
	}

	// Marshal JSONB fields
	businessHours, _ := json.Marshal(location.BusinessHours)
	settings, _ := json.Marshal(location.Settings)
	metadata, _ := json.Marshal(location.Metadata)

	query := `
		INSERT INTO locations (
			id, organization_id, location_code, name, location_type,
			phone, email, manager_user_id,
			address_line1, address_line2, city, state, country, postal_code,
			timezone, business_hours, is_active, is_primary,
			allow_sales, allow_purchases, tax_rate,
			notes, settings, metadata,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18,
			$19, $20, $21,
			$22, $23, $24,
			$25, $26, $27
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		location.ID, location.OrganizationID, location.LocationCode, location.Name, location.LocationType,
		location.Phone, location.Email, location.ManagerUserID,
		location.AddressLine1, location.AddressLine2, location.City, location.State, location.Country, location.PostalCode,
		location.Timezone, businessHours, location.IsActive, location.IsPrimary,
		location.AllowSales, location.AllowPurchases, location.TaxRate,
		location.Notes, settings, metadata,
		location.CreatedAt, location.UpdatedAt, location.CreatedBy,
	)
	return err
}

func (r *LocationRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*locations.Location, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, location_code, name, location_type,
			phone, email, manager_user_id,
			address_line1, address_line2, city, state, country, postal_code,
			timezone, business_hours, is_active, is_primary,
			allow_sales, allow_purchases, tax_rate,
			notes, settings, metadata,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM locations
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var loc locations.Location
	var businessHours []byte
	var settings []byte
	var metadata []byte

	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&loc.ID, &loc.OrganizationID, &loc.LocationCode, &loc.Name, &loc.LocationType,
		&loc.Phone, &loc.Email, &loc.ManagerUserID,
		&loc.AddressLine1, &loc.AddressLine2, &loc.City, &loc.State, &loc.Country, &loc.PostalCode,
		&loc.Timezone, &businessHours, &loc.IsActive, &loc.IsPrimary,
		&loc.AllowSales, &loc.AllowPurchases, &loc.TaxRate,
		&loc.Notes, &settings, &metadata,
		&loc.CreatedAt, &loc.UpdatedAt, &loc.CreatedBy, &loc.UpdatedBy, &loc.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB fields
	if len(businessHours) > 0 {
		json.Unmarshal(businessHours, &loc.BusinessHours)
	}
	if len(settings) > 0 {
		json.Unmarshal(settings, &loc.Settings)
	}
	if len(metadata) > 0 {
		json.Unmarshal(metadata, &loc.Metadata)
	}

	return &loc, nil
}

func (r *LocationRepository) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*locations.Location, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, location_code, name, location_type,
			phone, email, manager_user_id,
			address_line1, address_line2, city, state, country, postal_code,
			timezone, business_hours, is_active, is_primary,
			allow_sales, allow_purchases, tax_rate,
			notes, settings, metadata,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM locations
		WHERE organization_id = $1 AND location_code = $2 AND deleted_at IS NULL
	`

	var loc locations.Location
	var businessHours []byte
	var settings []byte
	var metadata []byte

	err := r.db.Pool.QueryRow(ctx, query, orgID, code).Scan(
		&loc.ID, &loc.OrganizationID, &loc.LocationCode, &loc.Name, &loc.LocationType,
		&loc.Phone, &loc.Email, &loc.ManagerUserID,
		&loc.AddressLine1, &loc.AddressLine2, &loc.City, &loc.State, &loc.Country, &loc.PostalCode,
		&loc.Timezone, &businessHours, &loc.IsActive, &loc.IsPrimary,
		&loc.AllowSales, &loc.AllowPurchases, &loc.TaxRate,
		&loc.Notes, &settings, &metadata,
		&loc.CreatedAt, &loc.UpdatedAt, &loc.CreatedBy, &loc.UpdatedBy, &loc.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB fields
	if len(businessHours) > 0 {
		json.Unmarshal(businessHours, &loc.BusinessHours)
	}
	if len(settings) > 0 {
		json.Unmarshal(settings, &loc.Settings)
	}
	if len(metadata) > 0 {
		json.Unmarshal(metadata, &loc.Metadata)
	}

	return &loc, nil
}

func (r *LocationRepository) Update(ctx context.Context, location *locations.Location) error {
	if err := r.db.SetOrganizationContext(ctx, location.OrganizationID.String()); err != nil {
		return err
	}

	// Marshal JSONB fields
	businessHours, _ := json.Marshal(location.BusinessHours)
	settings, _ := json.Marshal(location.Settings)
	metadata, _ := json.Marshal(location.Metadata)

	query := `
		UPDATE locations SET
			location_code = $3,
			name = $4,
			location_type = $5,
			phone = $6,
			email = $7,
			manager_user_id = $8,
			address_line1 = $9,
			address_line2 = $10,
			city = $11,
			state = $12,
			country = $13,
			postal_code = $14,
			timezone = $15,
			business_hours = $16,
			is_active = $17,
			is_primary = $18,
			allow_sales = $19,
			allow_purchases = $20,
			tax_rate = $21,
			notes = $22,
			settings = $23,
			metadata = $24,
			updated_at = $25,
			updated_by = $26
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		location.OrganizationID, location.ID,
		location.LocationCode, location.Name, location.LocationType,
		location.Phone, location.Email, location.ManagerUserID,
		location.AddressLine1, location.AddressLine2, location.City, location.State, location.Country, location.PostalCode,
		location.Timezone, businessHours, location.IsActive, location.IsPrimary,
		location.AllowSales, location.AllowPurchases, location.TaxRate,
		location.Notes, settings, metadata,
		location.UpdatedAt, location.UpdatedBy,
	)
	return err
}

func (r *LocationRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE locations
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

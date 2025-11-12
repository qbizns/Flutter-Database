package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/products"
)

// UoMRepository implements products.UoMRepository
type UoMRepository struct {
	db *DB
}

// NewUoMRepository creates a new UoM repository
func NewUoMRepository(db *DB) *UoMRepository {
	return &UoMRepository{db: db}
}

// CreateUoM creates a new unit of measure
func (r *UoMRepository) CreateUoM(ctx context.Context, uom *products.UnitOfMeasure) error {
	query := `
		INSERT INTO units_of_measure (
			id, organization_id, uom_code, uom_name, uom_type,
			is_base_unit, is_active, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		uom.ID, uom.OrganizationID, uom.Code, uom.Name, uom.Type,
		uom.IsBaseUnit, uom.IsActive, uom.CreatedAt,
	)

	return err
}

// GetUoM retrieves a UoM by ID
func (r *UoMRepository) GetUoM(ctx context.Context, id uuid.UUID) (*products.UnitOfMeasure, error) {
	query := `
		SELECT
			id, organization_id, uom_code, uom_name, uom_type,
			is_base_unit, is_active, created_at, deleted_at
		FROM units_of_measure
		WHERE id = $1 AND deleted_at IS NULL
	`

	var uom products.UnitOfMeasure
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&uom.ID, &uom.OrganizationID, &uom.Code, &uom.Name, &uom.Type,
		&uom.IsBaseUnit, &uom.IsActive, &uom.CreatedAt, &uom.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &uom, nil
}

// GetUoMByCode retrieves a UoM by code
func (r *UoMRepository) GetUoMByCode(ctx context.Context, code string, orgID *uuid.UUID) (*products.UnitOfMeasure, error) {
	query := `
		SELECT
			id, organization_id, uom_code, uom_name, uom_type,
			is_base_unit, is_active, created_at, deleted_at
		FROM units_of_measure
		WHERE LOWER(uom_code) = LOWER($1) AND (organization_id = $2 OR (organization_id IS NULL AND $2::UUID IS NULL)) AND deleted_at IS NULL
	`

	var uom products.UnitOfMeasure
	err := r.db.Pool.QueryRow(ctx, query, code, orgID).Scan(
		&uom.ID, &uom.OrganizationID, &uom.Code, &uom.Name, &uom.Type,
		&uom.IsBaseUnit, &uom.IsActive, &uom.CreatedAt, &uom.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &uom, nil
}

// UpdateUoM updates a UoM
func (r *UoMRepository) UpdateUoM(ctx context.Context, uom *products.UnitOfMeasure) error {
	query := `
		UPDATE units_of_measure SET
			uom_code = $2, uom_name = $3, uom_type = $4,
			is_base_unit = $5, is_active = $6
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		uom.ID, uom.Code, uom.Name, uom.Type,
		uom.IsBaseUnit, uom.IsActive,
	)

	return err
}

// DeleteUoM soft-deletes a UoM
func (r *UoMRepository) DeleteUoM(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE units_of_measure
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	return err
}

// ListUoMs retrieves UoMs with filters
func (r *UoMRepository) ListUoMs(ctx context.Context, filters products.UoMFilters) ([]products.UnitOfMeasure, error) {
	query := "SELECT id, organization_id, uom_code, uom_name, uom_type, is_base_unit, is_active, created_at, deleted_at FROM units_of_measure WHERE deleted_at IS NULL"

	args := []interface{}{}
	argCount := 0

	if filters.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND uom_type = $%d", argCount)
		args = append(args, filters.Type)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (uom_code ILIKE $%d OR uom_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	query += " ORDER BY uom_type, uom_name ASC"

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

	var uoms []products.UnitOfMeasure
	for rows.Next() {
		var uom products.UnitOfMeasure
		err := rows.Scan(
			&uom.ID, &uom.OrganizationID, &uom.Code, &uom.Name, &uom.Type,
			&uom.IsBaseUnit, &uom.IsActive, &uom.CreatedAt, &uom.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		uoms = append(uoms, uom)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uoms, nil
}

// CountUoMs counts UoMs matching filters
func (r *UoMRepository) CountUoMs(ctx context.Context, filters products.UoMFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM units_of_measure WHERE deleted_at IS NULL"

	args := []interface{}{}
	argCount := 0

	if filters.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND uom_type = $%d", argCount)
		args = append(args, filters.Type)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (uom_code ILIKE $%d OR uom_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// CreateConversion creates a UoM conversion
func (r *UoMRepository) CreateConversion(ctx context.Context, conversion *products.UoMConversion) error {
	query := `
		INSERT INTO uom_conversions (
			id, from_uom_id, to_uom_id, conversion_factor, created_at
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		conversion.ID, conversion.FromUoMID, conversion.ToUoMID, conversion.ConversionFactor, conversion.CreatedAt,
	)

	return err
}

// GetConversion retrieves a conversion
func (r *UoMRepository) GetConversion(ctx context.Context, fromID uuid.UUID, toID uuid.UUID) (*products.UoMConversion, error) {
	query := `
		SELECT
			id, from_uom_id, to_uom_id, conversion_factor, created_at, deleted_at
		FROM uom_conversions
		WHERE from_uom_id = $1 AND to_uom_id = $2 AND deleted_at IS NULL
	`

	var conv products.UoMConversion
	err := r.db.Pool.QueryRow(ctx, query, fromID, toID).Scan(
		&conv.ID, &conv.FromUoMID, &conv.ToUoMID, &conv.ConversionFactor, &conv.CreatedAt, &conv.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &conv, nil
}

// DeleteConversion deletes a conversion
func (r *UoMRepository) DeleteConversion(ctx context.Context, fromID uuid.UUID, toID uuid.UUID) error {
	query := `
		UPDATE uom_conversions
		SET deleted_at = $3
		WHERE from_uom_id = $1 AND to_uom_id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, fromID, toID, time.Now())
	return err
}

// ListConversionsFrom retrieves conversions from a UoM
func (r *UoMRepository) ListConversionsFrom(ctx context.Context, fromID uuid.UUID) ([]products.UoMConversion, error) {
	query := `
		SELECT id, from_uom_id, to_uom_id, conversion_factor, created_at, deleted_at
		FROM uom_conversions
		WHERE from_uom_id = $1 AND deleted_at IS NULL
		ORDER BY to_uom_id
	`

	rows, err := r.db.Pool.Query(ctx, query, fromID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []products.UoMConversion
	for rows.Next() {
		var conv products.UoMConversion
		err := rows.Scan(
			&conv.ID, &conv.FromUoMID, &conv.ToUoMID, &conv.ConversionFactor, &conv.CreatedAt, &conv.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		conversions = append(conversions, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversions, nil
}

// ListConversionsTo retrieves conversions to a UoM
func (r *UoMRepository) ListConversionsTo(ctx context.Context, toID uuid.UUID) ([]products.UoMConversion, error) {
	query := `
		SELECT id, from_uom_id, to_uom_id, conversion_factor, created_at, deleted_at
		FROM uom_conversions
		WHERE to_uom_id = $1 AND deleted_at IS NULL
		ORDER BY from_uom_id
	`

	rows, err := r.db.Pool.Query(ctx, query, toID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []products.UoMConversion
	for rows.Next() {
		var conv products.UoMConversion
		err := rows.Scan(
			&conv.ID, &conv.FromUoMID, &conv.ToUoMID, &conv.ConversionFactor, &conv.CreatedAt, &conv.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		conversions = append(conversions, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversions, nil
}

// ListByType retrieves UoMs by type
func (r *UoMRepository) ListByType(ctx context.Context, uomType string) ([]products.UnitOfMeasure, error) {
	query := `
		SELECT id, organization_id, uom_code, uom_name, uom_type, is_base_unit, is_active, created_at, deleted_at
		FROM units_of_measure
		WHERE uom_type = $1 AND deleted_at IS NULL
		ORDER BY is_base_unit DESC, uom_name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, uomType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uoms []products.UnitOfMeasure
	for rows.Next() {
		var uom products.UnitOfMeasure
		err := rows.Scan(
			&uom.ID, &uom.OrganizationID, &uom.Code, &uom.Name, &uom.Type,
			&uom.IsBaseUnit, &uom.IsActive, &uom.CreatedAt, &uom.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		uoms = append(uoms, uom)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uoms, nil
}

// ConvertQuantity converts quantity between two UoMs
func (r *UoMRepository) ConvertQuantity(ctx context.Context, fromID uuid.UUID, toID uuid.UUID, quantity float64) (float64, error) {
	// If same UoM, no conversion needed
	if fromID == toID {
		return quantity, nil
	}

	// Get the conversion factor
	conversion, err := r.GetConversion(ctx, fromID, toID)
	if err != nil {
		return 0, err
	}

	if conversion == nil {
		return 0, fmt.Errorf("no conversion found between these units")
	}

	return quantity * conversion.ConversionFactor, nil
}

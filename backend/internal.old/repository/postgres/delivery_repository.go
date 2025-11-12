package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/delivery"
)

type DeliveryRepository struct {
	db *DB
}

func NewDeliveryRepository(db *DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

// ==================== DELIVERY ZONES ====================

func (r *DeliveryRepository) ListDeliveryZones(ctx context.Context, orgID uuid.UUID, filters delivery.DeliveryZoneFilters) ([]delivery.DeliveryZone, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, zone_name, zone_code, description,
		       geofence, postal_codes, coverage_notes, base_delivery_fee, fee_type,
		       minimum_order_amount, free_delivery_threshold, estimated_delivery_time_minutes,
		       max_delivery_time_minutes, priority, is_active, active_hours, metadata,
		       created_at, updated_at, created_by, updated_by
		FROM delivery_zones
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (zone_name ILIKE $%d OR zone_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY priority DESC, zone_name ASC"

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

	var zones []delivery.DeliveryZone
	for rows.Next() {
		var z delivery.DeliveryZone
		err := rows.Scan(
			&z.ID, &z.OrganizationID, &z.LocationID, &z.ZoneName, &z.ZoneCode, &z.Description,
			&z.Geofence, &z.PostalCodes, &z.CoverageNotes, &z.BaseDeliveryFee, &z.FeeType,
			&z.MinimumOrderAmount, &z.FreeDeliveryThreshold, &z.EstimatedDeliveryMinutes,
			&z.MaxDeliveryTimeMinutes, &z.Priority, &z.IsActive, &z.ActiveHours, &z.Metadata,
			&z.CreatedAt, &z.UpdatedAt, &z.CreatedBy, &z.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		zones = append(zones, z)
	}
	return zones, rows.Err()
}

func (r *DeliveryRepository) CountDeliveryZones(ctx context.Context, orgID uuid.UUID, filters delivery.DeliveryZoneFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM delivery_zones WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (zone_name ILIKE $%d OR zone_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *DeliveryRepository) CreateDeliveryZone(ctx context.Context, zone *delivery.DeliveryZone) error {
	if err := r.db.SetOrganizationContext(ctx, zone.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO delivery_zones (
			id, organization_id, location_id, zone_name, zone_code, description,
			geofence, postal_codes, coverage_notes, base_delivery_fee, fee_type,
			minimum_order_amount, free_delivery_threshold, estimated_delivery_time_minutes,
			max_delivery_time_minutes, priority, is_active, active_hours, metadata,
			created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		zone.ID, zone.OrganizationID, zone.LocationID, zone.ZoneName, zone.ZoneCode, zone.Description,
		zone.Geofence, zone.PostalCodes, zone.CoverageNotes, zone.BaseDeliveryFee, zone.FeeType,
		zone.MinimumOrderAmount, zone.FreeDeliveryThreshold, zone.EstimatedDeliveryMinutes,
		zone.MaxDeliveryTimeMinutes, zone.Priority, zone.IsActive, zone.ActiveHours, zone.Metadata,
		zone.CreatedAt, zone.UpdatedAt, zone.CreatedBy,
	)
	return err
}

func (r *DeliveryRepository) GetDeliveryZone(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*delivery.DeliveryZone, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, zone_name, zone_code, description,
		       geofence, postal_codes, coverage_notes, base_delivery_fee, fee_type,
		       minimum_order_amount, free_delivery_threshold, estimated_delivery_time_minutes,
		       max_delivery_time_minutes, priority, is_active, active_hours, metadata,
		       created_at, updated_at, created_by, updated_by
		FROM delivery_zones
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var z delivery.DeliveryZone
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&z.ID, &z.OrganizationID, &z.LocationID, &z.ZoneName, &z.ZoneCode, &z.Description,
		&z.Geofence, &z.PostalCodes, &z.CoverageNotes, &z.BaseDeliveryFee, &z.FeeType,
		&z.MinimumOrderAmount, &z.FreeDeliveryThreshold, &z.EstimatedDeliveryMinutes,
		&z.MaxDeliveryTimeMinutes, &z.Priority, &z.IsActive, &z.ActiveHours, &z.Metadata,
		&z.CreatedAt, &z.UpdatedAt, &z.CreatedBy, &z.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &z, err
}

func (r *DeliveryRepository) UpdateDeliveryZone(ctx context.Context, zone *delivery.DeliveryZone) error {
	if err := r.db.SetOrganizationContext(ctx, zone.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_zones SET
			location_id = $3, zone_name = $4, zone_code = $5, description = $6,
			geofence = $7, postal_codes = $8, coverage_notes = $9, base_delivery_fee = $10,
			fee_type = $11, minimum_order_amount = $12, free_delivery_threshold = $13,
			estimated_delivery_time_minutes = $14, max_delivery_time_minutes = $15,
			priority = $16, is_active = $17, active_hours = $18, metadata = $19,
			updated_at = $20, updated_by = $21
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		zone.OrganizationID, zone.ID, zone.LocationID, zone.ZoneName, zone.ZoneCode, zone.Description,
		zone.Geofence, zone.PostalCodes, zone.CoverageNotes, zone.BaseDeliveryFee, zone.FeeType,
		zone.MinimumOrderAmount, zone.FreeDeliveryThreshold, zone.EstimatedDeliveryMinutes,
		zone.MaxDeliveryTimeMinutes, zone.Priority, zone.IsActive, zone.ActiveHours, zone.Metadata,
		zone.UpdatedAt, zone.UpdatedBy,
	)
	return err
}

func (r *DeliveryRepository) DeleteDeliveryZone(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_zones
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// ==================== DELIVERY DRIVERS ====================

func (r *DeliveryRepository) ListDeliveryDrivers(ctx context.Context, orgID uuid.UUID, filters delivery.DeliveryDriverFilters) ([]delivery.DeliveryDriver, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, user_id, driver_code, full_name, phone, email,
		       emergency_contact_name, emergency_contact_phone, vehicle_type, vehicle_make,
		       vehicle_model, vehicle_year, vehicle_color, license_plate, drivers_license_number,
		       license_expiry_date, insurance_policy_number, insurance_expiry_date, hire_date,
		       employment_type, status, total_deliveries, successful_deliveries, rating, rating_count,
		       current_location, is_available, last_location_update, commission_rate, payment_method,
		       documents, metadata, created_at, updated_at, created_by, updated_by
		FROM delivery_drivers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (full_name ILIKE $%d OR driver_code ILIKE $%d OR phone ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.IsAvailable != nil {
		argCount++
		query += fmt.Sprintf(" AND is_available = $%d", argCount)
		args = append(args, *filters.IsAvailable)
	}

	query += " ORDER BY full_name ASC"

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

	var drivers []delivery.DeliveryDriver
	for rows.Next() {
		var d delivery.DeliveryDriver
		err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.UserID, &d.DriverCode, &d.FullName, &d.Phone, &d.Email,
			&d.EmergencyContactName, &d.EmergencyContactPhone, &d.VehicleType, &d.VehicleMake,
			&d.VehicleModel, &d.VehicleYear, &d.VehicleColor, &d.LicensePlate, &d.DriversLicenseNumber,
			&d.LicenseExpiryDate, &d.InsurancePolicyNumber, &d.InsuranceExpiryDate, &d.HireDate,
			&d.EmploymentType, &d.Status, &d.TotalDeliveries, &d.SuccessfulDeliveries, &d.Rating, &d.RatingCount,
			&d.CurrentLocation, &d.IsAvailable, &d.LastLocationUpdate, &d.CommissionRate, &d.PaymentMethod,
			&d.Documents, &d.Metadata, &d.CreatedAt, &d.UpdatedAt, &d.CreatedBy, &d.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, rows.Err()
}

func (r *DeliveryRepository) CountDeliveryDrivers(ctx context.Context, orgID uuid.UUID, filters delivery.DeliveryDriverFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM delivery_drivers WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (full_name ILIKE $%d OR driver_code ILIKE $%d OR phone ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.IsAvailable != nil {
		argCount++
		query += fmt.Sprintf(" AND is_available = $%d", argCount)
		args = append(args, *filters.IsAvailable)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *DeliveryRepository) CreateDeliveryDriver(ctx context.Context, driver *delivery.DeliveryDriver) error {
	if err := r.db.SetOrganizationContext(ctx, driver.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO delivery_drivers (
			id, organization_id, user_id, driver_code, full_name, phone, email,
			emergency_contact_name, emergency_contact_phone, vehicle_type, vehicle_make,
			vehicle_model, vehicle_year, vehicle_color, license_plate, drivers_license_number,
			license_expiry_date, insurance_policy_number, insurance_expiry_date, hire_date,
			employment_type, status, total_deliveries, successful_deliveries, rating, rating_count,
			current_location, is_available, last_location_update, commission_rate, payment_method,
			documents, metadata, created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		driver.ID, driver.OrganizationID, driver.UserID, driver.DriverCode, driver.FullName, driver.Phone, driver.Email,
		driver.EmergencyContactName, driver.EmergencyContactPhone, driver.VehicleType, driver.VehicleMake,
		driver.VehicleModel, driver.VehicleYear, driver.VehicleColor, driver.LicensePlate, driver.DriversLicenseNumber,
		driver.LicenseExpiryDate, driver.InsurancePolicyNumber, driver.InsuranceExpiryDate, driver.HireDate,
		driver.EmploymentType, driver.Status, driver.TotalDeliveries, driver.SuccessfulDeliveries, driver.Rating, driver.RatingCount,
		driver.CurrentLocation, driver.IsAvailable, driver.LastLocationUpdate, driver.CommissionRate, driver.PaymentMethod,
		driver.Documents, driver.Metadata, driver.CreatedAt, driver.UpdatedAt, driver.CreatedBy,
	)
	return err
}

func (r *DeliveryRepository) GetDeliveryDriver(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*delivery.DeliveryDriver, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, user_id, driver_code, full_name, phone, email,
		       emergency_contact_name, emergency_contact_phone, vehicle_type, vehicle_make,
		       vehicle_model, vehicle_year, vehicle_color, license_plate, drivers_license_number,
		       license_expiry_date, insurance_policy_number, insurance_expiry_date, hire_date,
		       employment_type, status, total_deliveries, successful_deliveries, rating, rating_count,
		       current_location, is_available, last_location_update, commission_rate, payment_method,
		       documents, metadata, created_at, updated_at, created_by, updated_by
		FROM delivery_drivers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var d delivery.DeliveryDriver
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&d.ID, &d.OrganizationID, &d.UserID, &d.DriverCode, &d.FullName, &d.Phone, &d.Email,
		&d.EmergencyContactName, &d.EmergencyContactPhone, &d.VehicleType, &d.VehicleMake,
		&d.VehicleModel, &d.VehicleYear, &d.VehicleColor, &d.LicensePlate, &d.DriversLicenseNumber,
		&d.LicenseExpiryDate, &d.InsurancePolicyNumber, &d.InsuranceExpiryDate, &d.HireDate,
		&d.EmploymentType, &d.Status, &d.TotalDeliveries, &d.SuccessfulDeliveries, &d.Rating, &d.RatingCount,
		&d.CurrentLocation, &d.IsAvailable, &d.LastLocationUpdate, &d.CommissionRate, &d.PaymentMethod,
		&d.Documents, &d.Metadata, &d.CreatedAt, &d.UpdatedAt, &d.CreatedBy, &d.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &d, err
}

func (r *DeliveryRepository) GetDeliveryDriverByCode(ctx context.Context, orgID uuid.UUID, code string) (*delivery.DeliveryDriver, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, user_id, driver_code, full_name, phone, email,
		       emergency_contact_name, emergency_contact_phone, vehicle_type, vehicle_make,
		       vehicle_model, vehicle_year, vehicle_color, license_plate, drivers_license_number,
		       license_expiry_date, insurance_policy_number, insurance_expiry_date, hire_date,
		       employment_type, status, total_deliveries, successful_deliveries, rating, rating_count,
		       current_location, is_available, last_location_update, commission_rate, payment_method,
		       documents, metadata, created_at, updated_at, created_by, updated_by
		FROM delivery_drivers
		WHERE organization_id = $1 AND driver_code = $2 AND deleted_at IS NULL
	`

	var d delivery.DeliveryDriver
	err := r.db.Pool.QueryRow(ctx, query, orgID, code).Scan(
		&d.ID, &d.OrganizationID, &d.UserID, &d.DriverCode, &d.FullName, &d.Phone, &d.Email,
		&d.EmergencyContactName, &d.EmergencyContactPhone, &d.VehicleType, &d.VehicleMake,
		&d.VehicleModel, &d.VehicleYear, &d.VehicleColor, &d.LicensePlate, &d.DriversLicenseNumber,
		&d.LicenseExpiryDate, &d.InsurancePolicyNumber, &d.InsuranceExpiryDate, &d.HireDate,
		&d.EmploymentType, &d.Status, &d.TotalDeliveries, &d.SuccessfulDeliveries, &d.Rating, &d.RatingCount,
		&d.CurrentLocation, &d.IsAvailable, &d.LastLocationUpdate, &d.CommissionRate, &d.PaymentMethod,
		&d.Documents, &d.Metadata, &d.CreatedAt, &d.UpdatedAt, &d.CreatedBy, &d.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &d, err
}

func (r *DeliveryRepository) UpdateDeliveryDriver(ctx context.Context, driver *delivery.DeliveryDriver) error {
	if err := r.db.SetOrganizationContext(ctx, driver.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_drivers SET
			user_id = $3, driver_code = $4, full_name = $5, phone = $6, email = $7,
			emergency_contact_name = $8, emergency_contact_phone = $9, vehicle_type = $10,
			vehicle_make = $11, vehicle_model = $12, vehicle_year = $13, vehicle_color = $14,
			license_plate = $15, drivers_license_number = $16, license_expiry_date = $17,
			insurance_policy_number = $18, insurance_expiry_date = $19, hire_date = $20,
			employment_type = $21, status = $22, commission_rate = $23, payment_method = $24,
			documents = $25, metadata = $26, updated_at = $27, updated_by = $28
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		driver.OrganizationID, driver.ID, driver.UserID, driver.DriverCode, driver.FullName, driver.Phone, driver.Email,
		driver.EmergencyContactName, driver.EmergencyContactPhone, driver.VehicleType, driver.VehicleMake,
		driver.VehicleModel, driver.VehicleYear, driver.VehicleColor, driver.LicensePlate, driver.DriversLicenseNumber,
		driver.LicenseExpiryDate, driver.InsurancePolicyNumber, driver.InsuranceExpiryDate, driver.HireDate,
		driver.EmploymentType, driver.Status, driver.CommissionRate, driver.PaymentMethod,
		driver.Documents, driver.Metadata, driver.UpdatedAt, driver.UpdatedBy,
	)
	return err
}

func (r *DeliveryRepository) DeleteDeliveryDriver(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_drivers
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *DeliveryRepository) UpdateDriverAvailability(ctx context.Context, orgID uuid.UUID, driverID uuid.UUID, available bool, location json.RawMessage) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_drivers
		SET is_available = $3, current_location = $4, last_location_update = $5, updated_at = $6
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, driverID, available, location, time.Now(), time.Now())
	return err
}

func (r *DeliveryRepository) UpdateDriverRating(ctx context.Context, orgID uuid.UUID, driverID uuid.UUID, rating float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_drivers
		SET rating = $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, driverID, rating, time.Now())
	return err
}

// ==================== DRIVER SHIFTS ====================

func (r *DeliveryRepository) ListDriverShifts(ctx context.Context, orgID uuid.UUID, filters delivery.DriverShiftFilters) ([]delivery.DriverShift, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, driver_id, shift_date, scheduled_start_time,
		       scheduled_end_time, actual_start_time, actual_end_time, status, total_break_minutes,
		       total_deliveries, total_distance_km, total_earnings, notes, metadata,
		       created_at, updated_at, created_by, updated_by
		FROM driver_shifts
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.DriverID != nil {
		argCount++
		query += fmt.Sprintf(" AND driver_id = $%d", argCount)
		args = append(args, *filters.DriverID)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND shift_date >= $%d", argCount)
		args = append(args, filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND shift_date <= $%d", argCount)
		args = append(args, filters.ToDate)
	}

	query += " ORDER BY shift_date DESC"

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

	var shifts []delivery.DriverShift
	for rows.Next() {
		var s delivery.DriverShift
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.LocationID, &s.DriverID, &s.ShiftDate, &s.ScheduledStartTime,
			&s.ScheduledEndTime, &s.ActualStartTime, &s.ActualEndTime, &s.Status, &s.TotalBreakMinutes,
			&s.TotalDeliveries, &s.TotalDistanceKm, &s.TotalEarnings, &s.Notes, &s.Metadata,
			&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, s)
	}
	return shifts, rows.Err()
}

func (r *DeliveryRepository) CountDriverShifts(ctx context.Context, orgID uuid.UUID, filters delivery.DriverShiftFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM driver_shifts WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.DriverID != nil {
		argCount++
		query += fmt.Sprintf(" AND driver_id = $%d", argCount)
		args = append(args, *filters.DriverID)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND shift_date >= $%d", argCount)
		args = append(args, filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND shift_date <= $%d", argCount)
		args = append(args, filters.ToDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *DeliveryRepository) CreateDriverShift(ctx context.Context, shift *delivery.DriverShift) error {
	if err := r.db.SetOrganizationContext(ctx, shift.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO driver_shifts (
			id, organization_id, location_id, driver_id, shift_date, scheduled_start_time,
			scheduled_end_time, actual_start_time, actual_end_time, status, total_break_minutes,
			total_deliveries, total_distance_km, total_earnings, notes, metadata,
			created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		shift.ID, shift.OrganizationID, shift.LocationID, shift.DriverID, shift.ShiftDate, shift.ScheduledStartTime,
		shift.ScheduledEndTime, shift.ActualStartTime, shift.ActualEndTime, shift.Status, shift.TotalBreakMinutes,
		shift.TotalDeliveries, shift.TotalDistanceKm, shift.TotalEarnings, shift.Notes, shift.Metadata,
		shift.CreatedAt, shift.UpdatedAt, shift.CreatedBy,
	)
	return err
}

func (r *DeliveryRepository) GetDriverShift(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*delivery.DriverShift, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, driver_id, shift_date, scheduled_start_time,
		       scheduled_end_time, actual_start_time, actual_end_time, status, total_break_minutes,
		       total_deliveries, total_distance_km, total_earnings, notes, metadata,
		       created_at, updated_at, created_by, updated_by
		FROM driver_shifts
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var s delivery.DriverShift
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&s.ID, &s.OrganizationID, &s.LocationID, &s.DriverID, &s.ShiftDate, &s.ScheduledStartTime,
		&s.ScheduledEndTime, &s.ActualStartTime, &s.ActualEndTime, &s.Status, &s.TotalBreakMinutes,
		&s.TotalDeliveries, &s.TotalDistanceKm, &s.TotalEarnings, &s.Notes, &s.Metadata,
		&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (r *DeliveryRepository) UpdateDriverShift(ctx context.Context, shift *delivery.DriverShift) error {
	if err := r.db.SetOrganizationContext(ctx, shift.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE driver_shifts SET
			location_id = $3, shift_date = $4, scheduled_start_time = $5, scheduled_end_time = $6,
			actual_start_time = $7, actual_end_time = $8, status = $9, total_break_minutes = $10,
			total_deliveries = $11, total_distance_km = $12, total_earnings = $13, notes = $14,
			metadata = $15, updated_at = $16, updated_by = $17
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		shift.OrganizationID, shift.ID, shift.LocationID, shift.ShiftDate, shift.ScheduledStartTime, shift.ScheduledEndTime,
		shift.ActualStartTime, shift.ActualEndTime, shift.Status, shift.TotalBreakMinutes,
		shift.TotalDeliveries, shift.TotalDistanceKm, shift.TotalEarnings, shift.Notes,
		shift.Metadata, shift.UpdatedAt, shift.UpdatedBy,
	)
	return err
}

func (r *DeliveryRepository) DeleteDriverShift(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE driver_shifts
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// ==================== CUSTOMER ADDRESSES ====================

func (r *DeliveryRepository) ListCustomerAddresses(ctx context.Context, orgID uuid.UUID, filters delivery.CustomerAddressFilters) ([]delivery.CustomerAddress, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, customer_id, address_label, address_line_1, address_line_2,
		       city, state_province, postal_code, country, latitude, longitude, location_notes,
		       delivery_zone_id, is_default, is_active, metadata, created_at, updated_at, created_by, updated_by
		FROM customer_addresses
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{orgID}
	argCount := 1

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (address_label ILIKE $%d OR city ILIKE $%d OR postal_code ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY is_default DESC, created_at DESC"

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

	var addresses []delivery.CustomerAddress
	for rows.Next() {
		var a delivery.CustomerAddress
		err := rows.Scan(
			&a.ID, &a.OrganizationID, &a.CustomerID, &a.AddressLabel, &a.AddressLine1, &a.AddressLine2,
			&a.City, &a.StateProvince, &a.PostalCode, &a.Country, &a.Latitude, &a.Longitude, &a.LocationNotes,
			&a.DeliveryZoneID, &a.IsDefault, &a.IsActive, &a.Metadata, &a.CreatedAt, &a.UpdatedAt, &a.CreatedBy, &a.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, rows.Err()
}

func (r *DeliveryRepository) CountCustomerAddresses(ctx context.Context, orgID uuid.UUID, filters delivery.CustomerAddressFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM customer_addresses WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (address_label ILIKE $%d OR city ILIKE $%d OR postal_code ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *DeliveryRepository) CreateCustomerAddress(ctx context.Context, address *delivery.CustomerAddress) error {
	if err := r.db.SetOrganizationContext(ctx, address.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO customer_addresses (
			id, organization_id, customer_id, address_label, address_line_1, address_line_2,
			city, state_province, postal_code, country, latitude, longitude, location_notes,
			delivery_zone_id, is_default, is_active, metadata, created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		address.ID, address.OrganizationID, address.CustomerID, address.AddressLabel, address.AddressLine1, address.AddressLine2,
		address.City, address.StateProvince, address.PostalCode, address.Country, address.Latitude, address.Longitude, address.LocationNotes,
		address.DeliveryZoneID, address.IsDefault, address.IsActive, address.Metadata, address.CreatedAt, address.UpdatedAt, address.CreatedBy,
	)
	return err
}

func (r *DeliveryRepository) GetCustomerAddress(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*delivery.CustomerAddress, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, customer_id, address_label, address_line_1, address_line_2,
		       city, state_province, postal_code, country, latitude, longitude, location_notes,
		       delivery_zone_id, is_default, is_active, metadata, created_at, updated_at, created_by, updated_by
		FROM customer_addresses
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var a delivery.CustomerAddress
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&a.ID, &a.OrganizationID, &a.CustomerID, &a.AddressLabel, &a.AddressLine1, &a.AddressLine2,
		&a.City, &a.StateProvince, &a.PostalCode, &a.Country, &a.Latitude, &a.Longitude, &a.LocationNotes,
		&a.DeliveryZoneID, &a.IsDefault, &a.IsActive, &a.Metadata, &a.CreatedAt, &a.UpdatedAt, &a.CreatedBy, &a.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *DeliveryRepository) UpdateCustomerAddress(ctx context.Context, address *delivery.CustomerAddress) error {
	if err := r.db.SetOrganizationContext(ctx, address.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customer_addresses SET
			address_label = $3, address_line_1 = $4, address_line_2 = $5, city = $6,
			state_province = $7, postal_code = $8, country = $9, latitude = $10,
			longitude = $11, location_notes = $12, delivery_zone_id = $13, is_default = $14,
			is_active = $15, metadata = $16, updated_at = $17, updated_by = $18
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		address.OrganizationID, address.ID, address.AddressLabel, address.AddressLine1, address.AddressLine2, address.City,
		address.StateProvince, address.PostalCode, address.Country, address.Latitude, address.Longitude, address.LocationNotes,
		address.DeliveryZoneID, address.IsDefault, address.IsActive, address.Metadata, address.UpdatedAt, address.UpdatedBy,
	)
	return err
}

func (r *DeliveryRepository) DeleteCustomerAddress(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customer_addresses
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *DeliveryRepository) SetDefaultCustomerAddress(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, addressID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customer_addresses
		SET is_default = false
		WHERE organization_id = $1 AND customer_id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, customerID)
	if err != nil {
		return err
	}

	query = `
		UPDATE customer_addresses
		SET is_default = true, updated_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err = r.db.Pool.Exec(ctx, query, orgID, addressID, time.Now())
	return err
}

// ==================== DELIVERY ASSIGNMENTS ====================

func (r *DeliveryRepository) ListDeliveryAssignments(ctx context.Context, orgID uuid.UUID, filters delivery.DeliveryAssignmentFilters) ([]delivery.DeliveryAssignment, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, driver_id, driver_shift_id, delivery_zone_id,
		       customer_address_id, delivery_address, delivery_location, assigned_at, assigned_by,
		       status, accepted_at, picked_up_at, dispatched_at, arrived_at, delivered_at, failed_at,
		       estimated_pickup_time, estimated_delivery_time, distance_km, route_info,
		       delivery_fee, driver_commission, payment_method, cash_collected, signature_image_url,
		       delivery_photo_url, recipient_name, delivery_notes, failure_reason, failure_notes,
		       retry_count, customer_rating, customer_feedback, driver_notes, metadata,
		       created_at, updated_at, created_by, updated_by
		FROM delivery_assignments
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.DriverID != nil {
		argCount++
		query += fmt.Sprintf(" AND driver_id = $%d", argCount)
		args = append(args, *filters.DriverID)
	}

	if filters.OrderID != nil {
		argCount++
		query += fmt.Sprintf(" AND order_id = $%d", argCount)
		args = append(args, *filters.OrderID)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND assigned_at >= $%d", argCount)
		args = append(args, filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND assigned_at <= $%d", argCount)
		args = append(args, filters.ToDate)
	}

	query += " ORDER BY assigned_at DESC"

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

	var assignments []delivery.DeliveryAssignment
	for rows.Next() {
		var a delivery.DeliveryAssignment
		err := rows.Scan(
			&a.ID, &a.OrganizationID, &a.OrderID, &a.DriverID, &a.DriverShiftID, &a.DeliveryZoneID,
			&a.CustomerAddressID, &a.DeliveryAddress, &a.DeliveryLocation, &a.AssignedAt, &a.AssignedBy,
			&a.Status, &a.AcceptedAt, &a.PickedUpAt, &a.DispatchedAt, &a.ArrivedAt, &a.DeliveredAt, &a.FailedAt,
			&a.EstimatedPickupTime, &a.EstimatedDeliveryTime, &a.DistanceKm, &a.RouteInfo,
			&a.DeliveryFee, &a.DriverCommission, &a.PaymentMethod, &a.CashCollected, &a.SignatureImageURL,
			&a.DeliveryPhotoURL, &a.RecipientName, &a.DeliveryNotes, &a.FailureReason, &a.FailureNotes,
			&a.RetryCount, &a.CustomerRating, &a.CustomerFeedback, &a.DriverNotes, &a.Metadata,
			&a.CreatedAt, &a.UpdatedAt, &a.CreatedBy, &a.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}

func (r *DeliveryRepository) CountDeliveryAssignments(ctx context.Context, orgID uuid.UUID, filters delivery.DeliveryAssignmentFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM delivery_assignments WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.DriverID != nil {
		argCount++
		query += fmt.Sprintf(" AND driver_id = $%d", argCount)
		args = append(args, *filters.DriverID)
	}

	if filters.OrderID != nil {
		argCount++
		query += fmt.Sprintf(" AND order_id = $%d", argCount)
		args = append(args, *filters.OrderID)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND assigned_at >= $%d", argCount)
		args = append(args, filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND assigned_at <= $%d", argCount)
		args = append(args, filters.ToDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *DeliveryRepository) CreateDeliveryAssignment(ctx context.Context, assignment *delivery.DeliveryAssignment) error {
	if err := r.db.SetOrganizationContext(ctx, assignment.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO delivery_assignments (
			id, organization_id, order_id, driver_id, driver_shift_id, delivery_zone_id,
			customer_address_id, delivery_address, delivery_location, assigned_at, assigned_by,
			status, estimated_pickup_time, estimated_delivery_time, distance_km, route_info,
			delivery_fee, driver_commission, payment_method, cash_collected,
			signature_image_url, delivery_photo_url, recipient_name, delivery_notes,
			failure_reason, failure_notes, retry_count, customer_rating, customer_feedback,
			driver_notes, metadata, created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		          $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		assignment.ID, assignment.OrganizationID, assignment.OrderID, assignment.DriverID, assignment.DriverShiftID, assignment.DeliveryZoneID,
		assignment.CustomerAddressID, assignment.DeliveryAddress, assignment.DeliveryLocation, assignment.AssignedAt, assignment.AssignedBy,
		assignment.Status, assignment.EstimatedPickupTime, assignment.EstimatedDeliveryTime, assignment.DistanceKm, assignment.RouteInfo,
		assignment.DeliveryFee, assignment.DriverCommission, assignment.PaymentMethod, assignment.CashCollected,
		assignment.SignatureImageURL, assignment.DeliveryPhotoURL, assignment.RecipientName, assignment.DeliveryNotes,
		assignment.FailureReason, assignment.FailureNotes, assignment.RetryCount, assignment.CustomerRating, assignment.CustomerFeedback,
		assignment.DriverNotes, assignment.Metadata, assignment.CreatedAt, assignment.UpdatedAt, assignment.CreatedBy,
	)
	return err
}

func (r *DeliveryRepository) GetDeliveryAssignment(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*delivery.DeliveryAssignment, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, driver_id, driver_shift_id, delivery_zone_id,
		       customer_address_id, delivery_address, delivery_location, assigned_at, assigned_by,
		       status, accepted_at, picked_up_at, dispatched_at, arrived_at, delivered_at, failed_at,
		       estimated_pickup_time, estimated_delivery_time, distance_km, route_info,
		       delivery_fee, driver_commission, payment_method, cash_collected, signature_image_url,
		       delivery_photo_url, recipient_name, delivery_notes, failure_reason, failure_notes,
		       retry_count, customer_rating, customer_feedback, driver_notes, metadata,
		       created_at, updated_at, created_by, updated_by
		FROM delivery_assignments
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var a delivery.DeliveryAssignment
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&a.ID, &a.OrganizationID, &a.OrderID, &a.DriverID, &a.DriverShiftID, &a.DeliveryZoneID,
		&a.CustomerAddressID, &a.DeliveryAddress, &a.DeliveryLocation, &a.AssignedAt, &a.AssignedBy,
		&a.Status, &a.AcceptedAt, &a.PickedUpAt, &a.DispatchedAt, &a.ArrivedAt, &a.DeliveredAt, &a.FailedAt,
		&a.EstimatedPickupTime, &a.EstimatedDeliveryTime, &a.DistanceKm, &a.RouteInfo,
		&a.DeliveryFee, &a.DriverCommission, &a.PaymentMethod, &a.CashCollected, &a.SignatureImageURL,
		&a.DeliveryPhotoURL, &a.RecipientName, &a.DeliveryNotes, &a.FailureReason, &a.FailureNotes,
		&a.RetryCount, &a.CustomerRating, &a.CustomerFeedback, &a.DriverNotes, &a.Metadata,
		&a.CreatedAt, &a.UpdatedAt, &a.CreatedBy, &a.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *DeliveryRepository) GetAssignmentByOrderID(ctx context.Context, orgID uuid.UUID, orderID uuid.UUID) (*delivery.DeliveryAssignment, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, driver_id, driver_shift_id, delivery_zone_id,
		       customer_address_id, delivery_address, delivery_location, assigned_at, assigned_by,
		       status, accepted_at, picked_up_at, dispatched_at, arrived_at, delivered_at, failed_at,
		       estimated_pickup_time, estimated_delivery_time, distance_km, route_info,
		       delivery_fee, driver_commission, payment_method, cash_collected, signature_image_url,
		       delivery_photo_url, recipient_name, delivery_notes, failure_reason, failure_notes,
		       retry_count, customer_rating, customer_feedback, driver_notes, metadata,
		       created_at, updated_at, created_by, updated_by
		FROM delivery_assignments
		WHERE organization_id = $1 AND order_id = $2 AND deleted_at IS NULL
		LIMIT 1
	`

	var a delivery.DeliveryAssignment
	err := r.db.Pool.QueryRow(ctx, query, orgID, orderID).Scan(
		&a.ID, &a.OrganizationID, &a.OrderID, &a.DriverID, &a.DriverShiftID, &a.DeliveryZoneID,
		&a.CustomerAddressID, &a.DeliveryAddress, &a.DeliveryLocation, &a.AssignedAt, &a.AssignedBy,
		&a.Status, &a.AcceptedAt, &a.PickedUpAt, &a.DispatchedAt, &a.ArrivedAt, &a.DeliveredAt, &a.FailedAt,
		&a.EstimatedPickupTime, &a.EstimatedDeliveryTime, &a.DistanceKm, &a.RouteInfo,
		&a.DeliveryFee, &a.DriverCommission, &a.PaymentMethod, &a.CashCollected, &a.SignatureImageURL,
		&a.DeliveryPhotoURL, &a.RecipientName, &a.DeliveryNotes, &a.FailureReason, &a.FailureNotes,
		&a.RetryCount, &a.CustomerRating, &a.CustomerFeedback, &a.DriverNotes, &a.Metadata,
		&a.CreatedAt, &a.UpdatedAt, &a.CreatedBy, &a.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *DeliveryRepository) UpdateDeliveryAssignment(ctx context.Context, assignment *delivery.DeliveryAssignment) error {
	if err := r.db.SetOrganizationContext(ctx, assignment.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_assignments SET
			driver_shift_id = $3, delivery_zone_id = $4, customer_address_id = $5,
			delivery_address = $6, delivery_location = $7, status = $8, accepted_at = $9,
			picked_up_at = $10, dispatched_at = $11, arrived_at = $12, delivered_at = $13,
			failed_at = $14, estimated_pickup_time = $15, estimated_delivery_time = $16,
			distance_km = $17, route_info = $18, delivery_fee = $19, driver_commission = $20,
			payment_method = $21, cash_collected = $22, signature_image_url = $23,
			delivery_photo_url = $24, recipient_name = $25, delivery_notes = $26,
			failure_reason = $27, failure_notes = $28, retry_count = $29, customer_rating = $30,
			customer_feedback = $31, driver_notes = $32, metadata = $33, updated_at = $34, updated_by = $35
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		assignment.OrganizationID, assignment.ID, assignment.DriverShiftID, assignment.DeliveryZoneID,
		assignment.CustomerAddressID, assignment.DeliveryAddress, assignment.DeliveryLocation, assignment.Status,
		assignment.AcceptedAt, assignment.PickedUpAt, assignment.DispatchedAt, assignment.ArrivedAt, assignment.DeliveredAt,
		assignment.FailedAt, assignment.EstimatedPickupTime, assignment.EstimatedDeliveryTime, assignment.DistanceKm,
		assignment.RouteInfo, assignment.DeliveryFee, assignment.DriverCommission, assignment.PaymentMethod, assignment.CashCollected,
		assignment.SignatureImageURL, assignment.DeliveryPhotoURL, assignment.RecipientName, assignment.DeliveryNotes,
		assignment.FailureReason, assignment.FailureNotes, assignment.RetryCount, assignment.CustomerRating,
		assignment.CustomerFeedback, assignment.DriverNotes, assignment.Metadata, assignment.UpdatedAt, assignment.UpdatedBy,
	)
	return err
}

func (r *DeliveryRepository) DeleteDeliveryAssignment(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_assignments
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *DeliveryRepository) UpdateAssignmentStatus(ctx context.Context, orgID uuid.UUID, assignmentID uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE delivery_assignments
		SET status = $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, orgID, assignmentID, status, time.Now())
	return err
}

// ==================== ORDER TRACKING EVENTS ====================

func (r *DeliveryRepository) ListOrderTrackingEvents(ctx context.Context, orgID uuid.UUID, filters delivery.OrderTrackingEventFilters) ([]delivery.OrderTrackingEvent, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, delivery_assignment_id, event_type, event_timestamp,
		       event_message, location, location_name, actor_type, actor_id, actor_name, metadata, created_at, created_by
		FROM order_tracking_events
		WHERE organization_id = $1
	`
	args := []interface{}{orgID}
	argCount := 1

	if filters.EventType != nil {
		argCount++
		query += fmt.Sprintf(" AND event_type = $%d", argCount)
		args = append(args, *filters.EventType)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND event_timestamp >= $%d", argCount)
		args = append(args, filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND event_timestamp <= $%d", argCount)
		args = append(args, filters.ToDate)
	}

	query += " ORDER BY event_timestamp DESC"

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

	var events []delivery.OrderTrackingEvent
	for rows.Next() {
		var e delivery.OrderTrackingEvent
		err := rows.Scan(
			&e.ID, &e.OrganizationID, &e.OrderID, &e.DeliveryAssignmentID, &e.EventType, &e.EventTimestamp,
			&e.EventMessage, &e.Location, &e.LocationName, &e.ActorType, &e.ActorID, &e.ActorName, &e.Metadata,
			&e.CreatedAt, &e.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *DeliveryRepository) CountOrderTrackingEvents(ctx context.Context, orgID uuid.UUID, filters delivery.OrderTrackingEventFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM order_tracking_events WHERE organization_id = $1"
	args := []interface{}{orgID}
	argCount := 1

	if filters.EventType != nil {
		argCount++
		query += fmt.Sprintf(" AND event_type = $%d", argCount)
		args = append(args, *filters.EventType)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND event_timestamp >= $%d", argCount)
		args = append(args, filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND event_timestamp <= $%d", argCount)
		args = append(args, filters.ToDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *DeliveryRepository) CreateOrderTrackingEvent(ctx context.Context, event *delivery.OrderTrackingEvent) error {
	if err := r.db.SetOrganizationContext(ctx, event.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO order_tracking_events (
			id, organization_id, order_id, delivery_assignment_id, event_type, event_timestamp,
			event_message, location, location_name, actor_type, actor_id, actor_name, metadata, created_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		event.ID, event.OrganizationID, event.OrderID, event.DeliveryAssignmentID, event.EventType, event.EventTimestamp,
		event.EventMessage, event.Location, event.LocationName, event.ActorType, event.ActorID, event.ActorName, event.Metadata,
		event.CreatedAt, event.CreatedBy,
	)
	return err
}

func (r *DeliveryRepository) GetOrderTrackingEvent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*delivery.OrderTrackingEvent, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, delivery_assignment_id, event_type, event_timestamp,
		       event_message, location, location_name, actor_type, actor_id, actor_name, metadata, created_at, created_by
		FROM order_tracking_events
		WHERE organization_id = $1 AND id = $2
	`

	var e delivery.OrderTrackingEvent
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&e.ID, &e.OrganizationID, &e.OrderID, &e.DeliveryAssignmentID, &e.EventType, &e.EventTimestamp,
		&e.EventMessage, &e.Location, &e.LocationName, &e.ActorType, &e.ActorID, &e.ActorName, &e.Metadata,
		&e.CreatedAt, &e.CreatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &e, err
}

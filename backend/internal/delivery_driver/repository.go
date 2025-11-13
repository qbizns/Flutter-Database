package delivery_driver

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

// Repository handles database operations for DeliveryDrivers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DeliveryDrivers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DeliveryDrivers represents a delivery_drivers entity
type DeliveryDrivers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	UserId *uuid.UUID `json:"user_id" db:"user_id"`
	DriverCode string `json:"driver_code" db:"driver_code"`
	FullName string `json:"full_name" db:"full_name"`
	Phone *string `json:"phone" db:"phone"`
	Email *string `json:"email" db:"email"`
	EmergencyContactName *string `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyContactPhone *string `json:"emergency_contact_phone" db:"emergency_contact_phone"`
	VehicleType *string `json:"vehicle_type" db:"vehicle_type"`
	VehicleMake *string `json:"vehicle_make" db:"vehicle_make"`
	VehicleModel *string `json:"vehicle_model" db:"vehicle_model"`
	VehicleYear *int64 `json:"vehicle_year" db:"vehicle_year"`
	VehicleColor *string `json:"vehicle_color" db:"vehicle_color"`
	LicensePlate *string `json:"license_plate" db:"license_plate"`
	DriversLicenseNumber *string `json:"drivers_license_number" db:"drivers_license_number"`
	LicenseExpiryDate *time.Time `json:"license_expiry_date" db:"license_expiry_date"`
	InsurancePolicyNumber *string `json:"insurance_policy_number" db:"insurance_policy_number"`
	InsuranceExpiryDate *time.Time `json:"insurance_expiry_date" db:"insurance_expiry_date"`
	HireDate *time.Time `json:"hire_date" db:"hire_date"`
	EmploymentType *string `json:"employment_type" db:"employment_type"`
	Status *string `json:"status" db:"status"`
	TotalDeliveries *int64 `json:"total_deliveries" db:"total_deliveries"`
	SuccessfulDeliveries *int64 `json:"successful_deliveries" db:"successful_deliveries"`
	Rating *float64 `json:"rating" db:"rating"`
	RatingCount *int64 `json:"rating_count" db:"rating_count"`
	CurrentLocation json.RawMessage `json:"current_location" db:"current_location"`
	IsAvailable *bool `json:"is_available" db:"is_available"`
	LastLocationUpdate *time.Time `json:"last_location_update" db:"last_location_update"`
	CommissionRate *float64 `json:"commission_rate" db:"commission_rate"`
	PaymentMethod *string `json:"payment_method" db:"payment_method"`
	Documents json.RawMessage `json:"documents" db:"documents"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'active', *string `json:"'active'," db:"'active',"`
	'fullTime', *string `json:"'full_time'," db:"'full_time',"`
	'bike', *string `json:"'bike'," db:"'bike',"`
}

// Create inserts a new delivery_drivers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DeliveryDrivers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "delivery_drivers", duration, nil)
	}()

	query := `
		INSERT INTO delivery_drivers (
			, organization_id
			, user_id
			, driver_code
			, full_name
			, phone
			, email
			, emergency_contact_name
			, emergency_contact_phone
			, vehicle_type
			, vehicle_make
			, vehicle_model
			, vehicle_year
			, vehicle_color
			, license_plate
			, drivers_license_number
			, license_expiry_date
			, insurance_policy_number
			, insurance_expiry_date
			, hire_date
			, employment_type
			, status
			, total_deliveries
			, successful_deliveries
			, rating
			, rating_count
			, current_location
			, is_available
			, last_location_update
			, commission_rate
			, payment_method
			, documents
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'full_time',
			, 'bike',
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
			, $25
			, $26
			, $27
			, $28
			, $29
			, $30
			, $31
			, $32
			, $33
			, $36
			, $37
			, $38
			, $39
			, $40
			, $41
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.DriverCode,
		entity.FullName,
		entity.Phone,
		entity.Email,
		entity.EmergencyContactName,
		entity.EmergencyContactPhone,
		entity.VehicleType,
		entity.VehicleMake,
		entity.VehicleModel,
		entity.VehicleYear,
		entity.VehicleColor,
		entity.LicensePlate,
		entity.DriversLicenseNumber,
		entity.LicenseExpiryDate,
		entity.InsurancePolicyNumber,
		entity.InsuranceExpiryDate,
		entity.HireDate,
		entity.EmploymentType,
		entity.Status,
		entity.TotalDeliveries,
		entity.SuccessfulDeliveries,
		entity.Rating,
		entity.RatingCount,
		entity.CurrentLocation,
		entity.IsAvailable,
		entity.LastLocationUpdate,
		entity.CommissionRate,
		entity.PaymentMethod,
		entity.Documents,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'active',,
		entity.'fullTime',,
		entity.'bike',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create delivery_drivers", zap.Error(err))
		return fmt.Errorf("failed to create delivery_drivers: %w", err)
	}

	r.logger.Info("created delivery_drivers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a delivery_drivers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DeliveryDrivers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_drivers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, user_id
			, driver_code
			, full_name
			, phone
			, email
			, emergency_contact_name
			, emergency_contact_phone
			, vehicle_type
			, vehicle_make
			, vehicle_model
			, vehicle_year
			, vehicle_color
			, license_plate
			, drivers_license_number
			, license_expiry_date
			, insurance_policy_number
			, insurance_expiry_date
			, hire_date
			, employment_type
			, status
			, total_deliveries
			, successful_deliveries
			, rating
			, rating_count
			, current_location
			, is_available
			, last_location_update
			, commission_rate
			, payment_method
			, documents
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'full_time',
			, 'bike',
		FROM delivery_drivers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity DeliveryDrivers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.UserId,
		&entity.DriverCode,
		&entity.FullName,
		&entity.Phone,
		&entity.Email,
		&entity.EmergencyContactName,
		&entity.EmergencyContactPhone,
		&entity.VehicleType,
		&entity.VehicleMake,
		&entity.VehicleModel,
		&entity.VehicleYear,
		&entity.VehicleColor,
		&entity.LicensePlate,
		&entity.DriversLicenseNumber,
		&entity.LicenseExpiryDate,
		&entity.InsurancePolicyNumber,
		&entity.InsuranceExpiryDate,
		&entity.HireDate,
		&entity.EmploymentType,
		&entity.Status,
		&entity.TotalDeliveries,
		&entity.SuccessfulDeliveries,
		&entity.Rating,
		&entity.RatingCount,
		&entity.CurrentLocation,
		&entity.IsAvailable,
		&entity.LastLocationUpdate,
		&entity.CommissionRate,
		&entity.PaymentMethod,
		&entity.Documents,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'active',,
		&entity.'fullTime',,
		&entity.'bike',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("delivery_drivers not found")
	}

	if err != nil {
		r.logger.Error("failed to get delivery_drivers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get delivery_drivers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of delivery_drivers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DeliveryDrivers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_drivers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM delivery_drivers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count delivery_drivers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, driver_code
			, full_name
			, phone
			, email
			, emergency_contact_name
			, emergency_contact_phone
			, vehicle_type
			, vehicle_make
			, vehicle_model
			, vehicle_year
			, vehicle_color
			, license_plate
			, drivers_license_number
			, license_expiry_date
			, insurance_policy_number
			, insurance_expiry_date
			, hire_date
			, employment_type
			, status
			, total_deliveries
			, successful_deliveries
			, rating
			, rating_count
			, current_location
			, is_available
			, last_location_update
			, commission_rate
			, payment_method
			, documents
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'full_time',
			, 'bike',
		FROM delivery_drivers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list delivery_drivers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list delivery_drivers: %w", err)
	}
	defer rows.Close()

	var entities []*DeliveryDrivers
	for rows.Next() {
		var entity DeliveryDrivers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.DriverCode,
			&entity.FullName,
			&entity.Phone,
			&entity.Email,
			&entity.EmergencyContactName,
			&entity.EmergencyContactPhone,
			&entity.VehicleType,
			&entity.VehicleMake,
			&entity.VehicleModel,
			&entity.VehicleYear,
			&entity.VehicleColor,
			&entity.LicensePlate,
			&entity.DriversLicenseNumber,
			&entity.LicenseExpiryDate,
			&entity.InsurancePolicyNumber,
			&entity.InsuranceExpiryDate,
			&entity.HireDate,
			&entity.EmploymentType,
			&entity.Status,
			&entity.TotalDeliveries,
			&entity.SuccessfulDeliveries,
			&entity.Rating,
			&entity.RatingCount,
			&entity.CurrentLocation,
			&entity.IsAvailable,
			&entity.LastLocationUpdate,
			&entity.CommissionRate,
			&entity.PaymentMethod,
			&entity.Documents,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'active',,
			&entity.'fullTime',,
			&entity.'bike',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan delivery_drivers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating delivery_drivers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing delivery_drivers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DeliveryDrivers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "delivery_drivers", duration, nil)
	}()

	query := `
		UPDATE delivery_drivers
		SET
			, organization_id = $2
			, user_id = $3
			, driver_code = $4
			, full_name = $5
			, phone = $6
			, email = $7
			, emergency_contact_name = $8
			, emergency_contact_phone = $9
			, vehicle_type = $10
			, vehicle_make = $11
			, vehicle_model = $12
			, vehicle_year = $13
			, vehicle_color = $14
			, license_plate = $15
			, drivers_license_number = $16
			, license_expiry_date = $17
			, insurance_policy_number = $18
			, insurance_expiry_date = $19
			, hire_date = $20
			, employment_type = $21
			, status = $22
			, total_deliveries = $23
			, successful_deliveries = $24
			, rating = $25
			, rating_count = $26
			, current_location = $27
			, is_available = $28
			, last_location_update = $29
			, commission_rate = $30
			, payment_method = $31
			, documents = $32
			, metadata = $33
			, updated_at = $35
			, created_by = $36
			, updated_by = $37
			, deleted_at = $38
			, 'active', = $39
			, 'full_time', = $40
			, 'bike', = $41
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $42
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.DriverCode,
		entity.FullName,
		entity.Phone,
		entity.Email,
		entity.EmergencyContactName,
		entity.EmergencyContactPhone,
		entity.VehicleType,
		entity.VehicleMake,
		entity.VehicleModel,
		entity.VehicleYear,
		entity.VehicleColor,
		entity.LicensePlate,
		entity.DriversLicenseNumber,
		entity.LicenseExpiryDate,
		entity.InsurancePolicyNumber,
		entity.InsuranceExpiryDate,
		entity.HireDate,
		entity.EmploymentType,
		entity.Status,
		entity.TotalDeliveries,
		entity.SuccessfulDeliveries,
		entity.Rating,
		entity.RatingCount,
		entity.CurrentLocation,
		entity.IsAvailable,
		entity.LastLocationUpdate,
		entity.CommissionRate,
		entity.PaymentMethod,
		entity.Documents,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'active',,
		entity.'fullTime',,
		entity.'bike',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update delivery_drivers", zap.Error(err))
		return fmt.Errorf("failed to update delivery_drivers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delivery_drivers not found or already deleted")
	}

	r.logger.Info("updated delivery_drivers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a delivery_drivers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "delivery_drivers", duration, nil)
	}()

	query := `
		UPDATE delivery_drivers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete delivery_drivers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete delivery_drivers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delivery_drivers not found or already deleted")
	}

	r.logger.Info("deleted delivery_drivers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves delivery_drivers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*DeliveryDrivers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_drivers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM delivery_drivers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count delivery_drivers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, driver_code
			, full_name
			, phone
			, email
			, emergency_contact_name
			, emergency_contact_phone
			, vehicle_type
			, vehicle_make
			, vehicle_model
			, vehicle_year
			, vehicle_color
			, license_plate
			, drivers_license_number
			, license_expiry_date
			, insurance_policy_number
			, insurance_expiry_date
			, hire_date
			, employment_type
			, status
			, total_deliveries
			, successful_deliveries
			, rating
			, rating_count
			, current_location
			, is_available
			, last_location_update
			, commission_rate
			, payment_method
			, documents
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'full_time',
			, 'bike',
		FROM delivery_drivers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list delivery_drivers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list delivery_drivers: %w", err)
	}
	defer rows.Close()

	var entities []*DeliveryDrivers
	for rows.Next() {
		var entity DeliveryDrivers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.DriverCode,
			&entity.FullName,
			&entity.Phone,
			&entity.Email,
			&entity.EmergencyContactName,
			&entity.EmergencyContactPhone,
			&entity.VehicleType,
			&entity.VehicleMake,
			&entity.VehicleModel,
			&entity.VehicleYear,
			&entity.VehicleColor,
			&entity.LicensePlate,
			&entity.DriversLicenseNumber,
			&entity.LicenseExpiryDate,
			&entity.InsurancePolicyNumber,
			&entity.InsuranceExpiryDate,
			&entity.HireDate,
			&entity.EmploymentType,
			&entity.Status,
			&entity.TotalDeliveries,
			&entity.SuccessfulDeliveries,
			&entity.Rating,
			&entity.RatingCount,
			&entity.CurrentLocation,
			&entity.IsAvailable,
			&entity.LastLocationUpdate,
			&entity.CommissionRate,
			&entity.PaymentMethod,
			&entity.Documents,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'active',,
			&entity.'fullTime',,
			&entity.'bike',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan delivery_drivers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


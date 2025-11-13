package delivery_assignment

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

// Repository handles database operations for DeliveryAssignments
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DeliveryAssignments repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DeliveryAssignments represents a delivery_assignments entity
type DeliveryAssignments struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	OrderId uuid.UUID `json:"order_id" db:"order_id"`
	DriverId uuid.UUID `json:"driver_id" db:"driver_id"`
	DriverShiftId *uuid.UUID `json:"driver_shift_id" db:"driver_shift_id"`
	DeliveryZoneId *uuid.UUID `json:"delivery_zone_id" db:"delivery_zone_id"`
	CustomerAddressId *uuid.UUID `json:"customer_address_id" db:"customer_address_id"`
	DeliveryAddress string `json:"delivery_address" db:"delivery_address"`
	DeliveryLocation json.RawMessage `json:"delivery_location" db:"delivery_location"`
	AssignedAt *time.Time `json:"assigned_at" db:"assigned_at"`
	AssignedBy *uuid.UUID `json:"assigned_by" db:"assigned_by"`
	Status *string `json:"status" db:"status"`
	AcceptedAt *time.Time `json:"accepted_at" db:"accepted_at"`
	PickedUpAt *time.Time `json:"picked_up_at" db:"picked_up_at"`
	DispatchedAt *time.Time `json:"dispatched_at" db:"dispatched_at"`
	ArrivedAt *time.Time `json:"arrived_at" db:"arrived_at"`
	DeliveredAt *time.Time `json:"delivered_at" db:"delivered_at"`
	FailedAt *time.Time `json:"failed_at" db:"failed_at"`
	EstimatedPickupTime *time.Time `json:"estimated_pickup_time" db:"estimated_pickup_time"`
	EstimatedDeliveryTime *time.Time `json:"estimated_delivery_time" db:"estimated_delivery_time"`
	DistanceKm *float64 `json:"distance_km" db:"distance_km"`
	RouteInfo json.RawMessage `json:"route_info" db:"route_info"`
	DeliveryFee *float64 `json:"delivery_fee" db:"delivery_fee"`
	DriverCommission *float64 `json:"driver_commission" db:"driver_commission"`
	PaymentMethod *string `json:"payment_method" db:"payment_method"`
	CashCollected *float64 `json:"cash_collected" db:"cash_collected"`
	SignatureImageUrl *string `json:"signature_image_url" db:"signature_image_url"`
	DeliveryPhotoUrl *string `json:"delivery_photo_url" db:"delivery_photo_url"`
	RecipientName *string `json:"recipient_name" db:"recipient_name"`
	DeliveryNotes *string `json:"delivery_notes" db:"delivery_notes"`
	FailureReason *string `json:"failure_reason" db:"failure_reason"`
	FailureNotes *string `json:"failure_notes" db:"failure_notes"`
	RetryCount *int64 `json:"retry_count" db:"retry_count"`
	CustomerRating *int64 `json:"customer_rating" db:"customer_rating"`
	CustomerFeedback *string `json:"customer_feedback" db:"customer_feedback"`
	DriverNotes *string `json:"driver_notes" db:"driver_notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'assigned', *string `json:"'assigned'," db:"'assigned',"`
	'delivered', *string `json:"'delivered'," db:"'delivered',"`
	CustomerRating *string `json:"customer_rating" db:"customer_rating"`
	DeliveryFee *string `json:"delivery_fee" db:"delivery_fee"`
}

// Create inserts a new delivery_assignments record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DeliveryAssignments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "delivery_assignments", duration, nil)
	}()

	query := `
		INSERT INTO delivery_assignments (
			, organization_id
			, order_id
			, driver_id
			, driver_shift_id
			, delivery_zone_id
			, customer_address_id
			, delivery_address
			, delivery_location
			, assigned_at
			, assigned_by
			, status
			, accepted_at
			, picked_up_at
			, dispatched_at
			, arrived_at
			, delivered_at
			, failed_at
			, estimated_pickup_time
			, estimated_delivery_time
			, distance_km
			, route_info
			, delivery_fee
			, driver_commission
			, payment_method
			, cash_collected
			, signature_image_url
			, delivery_photo_url
			, recipient_name
			, delivery_notes
			, failure_reason
			, failure_notes
			, retry_count
			, customer_rating
			, customer_feedback
			, driver_notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'assigned',
			, 'delivered',
			, customer_rating
			, delivery_fee
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
			, $34
			, $35
			, $36
			, $37
			, $40
			, $41
			, $42
			, $43
			, $44
			, $45
			, $46
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.OrderId,
		entity.DriverId,
		entity.DriverShiftId,
		entity.DeliveryZoneId,
		entity.CustomerAddressId,
		entity.DeliveryAddress,
		entity.DeliveryLocation,
		entity.AssignedAt,
		entity.AssignedBy,
		entity.Status,
		entity.AcceptedAt,
		entity.PickedUpAt,
		entity.DispatchedAt,
		entity.ArrivedAt,
		entity.DeliveredAt,
		entity.FailedAt,
		entity.EstimatedPickupTime,
		entity.EstimatedDeliveryTime,
		entity.DistanceKm,
		entity.RouteInfo,
		entity.DeliveryFee,
		entity.DriverCommission,
		entity.PaymentMethod,
		entity.CashCollected,
		entity.SignatureImageUrl,
		entity.DeliveryPhotoUrl,
		entity.RecipientName,
		entity.DeliveryNotes,
		entity.FailureReason,
		entity.FailureNotes,
		entity.RetryCount,
		entity.CustomerRating,
		entity.CustomerFeedback,
		entity.DriverNotes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'assigned',,
		entity.'delivered',,
		entity.CustomerRating,
		entity.DeliveryFee,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create delivery_assignments", zap.Error(err))
		return fmt.Errorf("failed to create delivery_assignments: %w", err)
	}

	r.logger.Info("created delivery_assignments",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a delivery_assignments by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DeliveryAssignments, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_assignments", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, order_id
			, driver_id
			, driver_shift_id
			, delivery_zone_id
			, customer_address_id
			, delivery_address
			, delivery_location
			, assigned_at
			, assigned_by
			, status
			, accepted_at
			, picked_up_at
			, dispatched_at
			, arrived_at
			, delivered_at
			, failed_at
			, estimated_pickup_time
			, estimated_delivery_time
			, distance_km
			, route_info
			, delivery_fee
			, driver_commission
			, payment_method
			, cash_collected
			, signature_image_url
			, delivery_photo_url
			, recipient_name
			, delivery_notes
			, failure_reason
			, failure_notes
			, retry_count
			, customer_rating
			, customer_feedback
			, driver_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'assigned',
			, 'delivered',
			, customer_rating
			, delivery_fee
		FROM delivery_assignments
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity DeliveryAssignments
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.OrderId,
		&entity.DriverId,
		&entity.DriverShiftId,
		&entity.DeliveryZoneId,
		&entity.CustomerAddressId,
		&entity.DeliveryAddress,
		&entity.DeliveryLocation,
		&entity.AssignedAt,
		&entity.AssignedBy,
		&entity.Status,
		&entity.AcceptedAt,
		&entity.PickedUpAt,
		&entity.DispatchedAt,
		&entity.ArrivedAt,
		&entity.DeliveredAt,
		&entity.FailedAt,
		&entity.EstimatedPickupTime,
		&entity.EstimatedDeliveryTime,
		&entity.DistanceKm,
		&entity.RouteInfo,
		&entity.DeliveryFee,
		&entity.DriverCommission,
		&entity.PaymentMethod,
		&entity.CashCollected,
		&entity.SignatureImageUrl,
		&entity.DeliveryPhotoUrl,
		&entity.RecipientName,
		&entity.DeliveryNotes,
		&entity.FailureReason,
		&entity.FailureNotes,
		&entity.RetryCount,
		&entity.CustomerRating,
		&entity.CustomerFeedback,
		&entity.DriverNotes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'assigned',,
		&entity.'delivered',,
		&entity.CustomerRating,
		&entity.DeliveryFee,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("delivery_assignments not found")
	}

	if err != nil {
		r.logger.Error("failed to get delivery_assignments", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get delivery_assignments: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of delivery_assignments records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DeliveryAssignments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_assignments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM delivery_assignments
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count delivery_assignments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_id
			, driver_id
			, driver_shift_id
			, delivery_zone_id
			, customer_address_id
			, delivery_address
			, delivery_location
			, assigned_at
			, assigned_by
			, status
			, accepted_at
			, picked_up_at
			, dispatched_at
			, arrived_at
			, delivered_at
			, failed_at
			, estimated_pickup_time
			, estimated_delivery_time
			, distance_km
			, route_info
			, delivery_fee
			, driver_commission
			, payment_method
			, cash_collected
			, signature_image_url
			, delivery_photo_url
			, recipient_name
			, delivery_notes
			, failure_reason
			, failure_notes
			, retry_count
			, customer_rating
			, customer_feedback
			, driver_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'assigned',
			, 'delivered',
			, customer_rating
			, delivery_fee
		FROM delivery_assignments
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list delivery_assignments", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list delivery_assignments: %w", err)
	}
	defer rows.Close()

	var entities []*DeliveryAssignments
	for rows.Next() {
		var entity DeliveryAssignments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderId,
			&entity.DriverId,
			&entity.DriverShiftId,
			&entity.DeliveryZoneId,
			&entity.CustomerAddressId,
			&entity.DeliveryAddress,
			&entity.DeliveryLocation,
			&entity.AssignedAt,
			&entity.AssignedBy,
			&entity.Status,
			&entity.AcceptedAt,
			&entity.PickedUpAt,
			&entity.DispatchedAt,
			&entity.ArrivedAt,
			&entity.DeliveredAt,
			&entity.FailedAt,
			&entity.EstimatedPickupTime,
			&entity.EstimatedDeliveryTime,
			&entity.DistanceKm,
			&entity.RouteInfo,
			&entity.DeliveryFee,
			&entity.DriverCommission,
			&entity.PaymentMethod,
			&entity.CashCollected,
			&entity.SignatureImageUrl,
			&entity.DeliveryPhotoUrl,
			&entity.RecipientName,
			&entity.DeliveryNotes,
			&entity.FailureReason,
			&entity.FailureNotes,
			&entity.RetryCount,
			&entity.CustomerRating,
			&entity.CustomerFeedback,
			&entity.DriverNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'assigned',,
			&entity.'delivered',,
			&entity.CustomerRating,
			&entity.DeliveryFee,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan delivery_assignments: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating delivery_assignments rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing delivery_assignments record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DeliveryAssignments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "delivery_assignments", duration, nil)
	}()

	query := `
		UPDATE delivery_assignments
		SET
			, organization_id = $2
			, order_id = $3
			, driver_id = $4
			, driver_shift_id = $5
			, delivery_zone_id = $6
			, customer_address_id = $7
			, delivery_address = $8
			, delivery_location = $9
			, assigned_at = $10
			, assigned_by = $11
			, status = $12
			, accepted_at = $13
			, picked_up_at = $14
			, dispatched_at = $15
			, arrived_at = $16
			, delivered_at = $17
			, failed_at = $18
			, estimated_pickup_time = $19
			, estimated_delivery_time = $20
			, distance_km = $21
			, route_info = $22
			, delivery_fee = $23
			, driver_commission = $24
			, payment_method = $25
			, cash_collected = $26
			, signature_image_url = $27
			, delivery_photo_url = $28
			, recipient_name = $29
			, delivery_notes = $30
			, failure_reason = $31
			, failure_notes = $32
			, retry_count = $33
			, customer_rating = $34
			, customer_feedback = $35
			, driver_notes = $36
			, metadata = $37
			, updated_at = $39
			, created_by = $40
			, updated_by = $41
			, deleted_at = $42
			, 'assigned', = $43
			, 'delivered', = $44
			, customer_rating = $45
			, delivery_fee = $46
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $47
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.OrderId,
		entity.DriverId,
		entity.DriverShiftId,
		entity.DeliveryZoneId,
		entity.CustomerAddressId,
		entity.DeliveryAddress,
		entity.DeliveryLocation,
		entity.AssignedAt,
		entity.AssignedBy,
		entity.Status,
		entity.AcceptedAt,
		entity.PickedUpAt,
		entity.DispatchedAt,
		entity.ArrivedAt,
		entity.DeliveredAt,
		entity.FailedAt,
		entity.EstimatedPickupTime,
		entity.EstimatedDeliveryTime,
		entity.DistanceKm,
		entity.RouteInfo,
		entity.DeliveryFee,
		entity.DriverCommission,
		entity.PaymentMethod,
		entity.CashCollected,
		entity.SignatureImageUrl,
		entity.DeliveryPhotoUrl,
		entity.RecipientName,
		entity.DeliveryNotes,
		entity.FailureReason,
		entity.FailureNotes,
		entity.RetryCount,
		entity.CustomerRating,
		entity.CustomerFeedback,
		entity.DriverNotes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'assigned',,
		entity.'delivered',,
		entity.CustomerRating,
		entity.DeliveryFee,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update delivery_assignments", zap.Error(err))
		return fmt.Errorf("failed to update delivery_assignments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delivery_assignments not found or already deleted")
	}

	r.logger.Info("updated delivery_assignments",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a delivery_assignments record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "delivery_assignments", duration, nil)
	}()

	query := `
		UPDATE delivery_assignments
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete delivery_assignments", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete delivery_assignments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delivery_assignments not found or already deleted")
	}

	r.logger.Info("deleted delivery_assignments", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves delivery_assignments records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*DeliveryAssignments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_assignments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM delivery_assignments
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count delivery_assignments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_id
			, driver_id
			, driver_shift_id
			, delivery_zone_id
			, customer_address_id
			, delivery_address
			, delivery_location
			, assigned_at
			, assigned_by
			, status
			, accepted_at
			, picked_up_at
			, dispatched_at
			, arrived_at
			, delivered_at
			, failed_at
			, estimated_pickup_time
			, estimated_delivery_time
			, distance_km
			, route_info
			, delivery_fee
			, driver_commission
			, payment_method
			, cash_collected
			, signature_image_url
			, delivery_photo_url
			, recipient_name
			, delivery_notes
			, failure_reason
			, failure_notes
			, retry_count
			, customer_rating
			, customer_feedback
			, driver_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'assigned',
			, 'delivered',
			, customer_rating
			, delivery_fee
		FROM delivery_assignments
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list delivery_assignments by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list delivery_assignments: %w", err)
	}
	defer rows.Close()

	var entities []*DeliveryAssignments
	for rows.Next() {
		var entity DeliveryAssignments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderId,
			&entity.DriverId,
			&entity.DriverShiftId,
			&entity.DeliveryZoneId,
			&entity.CustomerAddressId,
			&entity.DeliveryAddress,
			&entity.DeliveryLocation,
			&entity.AssignedAt,
			&entity.AssignedBy,
			&entity.Status,
			&entity.AcceptedAt,
			&entity.PickedUpAt,
			&entity.DispatchedAt,
			&entity.ArrivedAt,
			&entity.DeliveredAt,
			&entity.FailedAt,
			&entity.EstimatedPickupTime,
			&entity.EstimatedDeliveryTime,
			&entity.DistanceKm,
			&entity.RouteInfo,
			&entity.DeliveryFee,
			&entity.DriverCommission,
			&entity.PaymentMethod,
			&entity.CashCollected,
			&entity.SignatureImageUrl,
			&entity.DeliveryPhotoUrl,
			&entity.RecipientName,
			&entity.DeliveryNotes,
			&entity.FailureReason,
			&entity.FailureNotes,
			&entity.RetryCount,
			&entity.CustomerRating,
			&entity.CustomerFeedback,
			&entity.DriverNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'assigned',,
			&entity.'delivered',,
			&entity.CustomerRating,
			&entity.DeliveryFee,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan delivery_assignments: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


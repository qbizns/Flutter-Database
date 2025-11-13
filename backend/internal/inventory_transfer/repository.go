package inventory_transfer

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

// Repository handles database operations for InventoryTransfers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new InventoryTransfers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// InventoryTransfers represents a inventory_transfers entity
type InventoryTransfers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	TransferNumber string `json:"transfer_number" db:"transfer_number"`
	TransferDate time.Time `json:"transfer_date" db:"transfer_date"`
	FromLocationId uuid.UUID `json:"from_location_id" db:"from_location_id"`
	ToLocationId uuid.UUID `json:"to_location_id" db:"to_location_id"`
	Status *string `json:"status" db:"status"`
	RequestedDate *time.Time `json:"requested_date" db:"requested_date"`
	ApprovedDate *time.Time `json:"approved_date" db:"approved_date"`
	ShippedDate *time.Time `json:"shipped_date" db:"shipped_date"`
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date" db:"expected_delivery_date"`
	ReceivedDate *time.Time `json:"received_date" db:"received_date"`
	Carrier *string `json:"carrier" db:"carrier"`
	TrackingNumber *string `json:"tracking_number" db:"tracking_number"`
	ShippingCost *float64 `json:"shipping_cost" db:"shipping_cost"`
	Reason *string `json:"reason" db:"reason"`
	Notes *string `json:"notes" db:"notes"`
	RejectionReason *string `json:"rejection_reason" db:"rejection_reason"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	RequestedBy *uuid.UUID `json:"requested_by" db:"requested_by"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ShippedBy *uuid.UUID `json:"shipped_by" db:"shipped_by"`
	ReceivedBy *uuid.UUID `json:"received_by" db:"received_by"`
	(approvedDate *string `json:"(approved_date" db:"(approved_date"`
	(shippedDate *string `json:"(shipped_date" db:"(shipped_date"`
	(receivedDate *string `json:"(received_date" db:"(received_date"`
}

// Create inserts a new inventory_transfers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *InventoryTransfers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "inventory_transfers", duration, nil)
	}()

	query := `
		INSERT INTO inventory_transfers (
			, organization_id
			, transfer_number
			, transfer_date
			, from_location_id
			, to_location_id
			, status
			, requested_date
			, approved_date
			, shipped_date
			, expected_delivery_date
			, received_date
			, carrier
			, tracking_number
			, shipping_cost
			, reason
			, notes
			, rejection_reason
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, requested_by
			, approved_by
			, shipped_by
			, received_by
			, (approved_date
			, (shipped_date
			, (received_date
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.TransferNumber,
		entity.TransferDate,
		entity.FromLocationId,
		entity.ToLocationId,
		entity.Status,
		entity.RequestedDate,
		entity.ApprovedDate,
		entity.ShippedDate,
		entity.ExpectedDeliveryDate,
		entity.ReceivedDate,
		entity.Carrier,
		entity.TrackingNumber,
		entity.ShippingCost,
		entity.Reason,
		entity.Notes,
		entity.RejectionReason,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.RequestedBy,
		entity.ApprovedBy,
		entity.ShippedBy,
		entity.ReceivedBy,
		entity.(approvedDate,
		entity.(shippedDate,
		entity.(receivedDate,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create inventory_transfers", zap.Error(err))
		return fmt.Errorf("failed to create inventory_transfers: %w", err)
	}

	r.logger.Info("created inventory_transfers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a inventory_transfers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*InventoryTransfers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transfers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, transfer_number
			, transfer_date
			, from_location_id
			, to_location_id
			, status
			, requested_date
			, approved_date
			, shipped_date
			, expected_delivery_date
			, received_date
			, carrier
			, tracking_number
			, shipping_cost
			, reason
			, notes
			, rejection_reason
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, requested_by
			, approved_by
			, shipped_by
			, received_by
			, (approved_date
			, (shipped_date
			, (received_date
		FROM inventory_transfers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity InventoryTransfers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.TransferNumber,
		&entity.TransferDate,
		&entity.FromLocationId,
		&entity.ToLocationId,
		&entity.Status,
		&entity.RequestedDate,
		&entity.ApprovedDate,
		&entity.ShippedDate,
		&entity.ExpectedDeliveryDate,
		&entity.ReceivedDate,
		&entity.Carrier,
		&entity.TrackingNumber,
		&entity.ShippingCost,
		&entity.Reason,
		&entity.Notes,
		&entity.RejectionReason,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.RequestedBy,
		&entity.ApprovedBy,
		&entity.ShippedBy,
		&entity.ReceivedBy,
		&entity.(approvedDate,
		&entity.(shippedDate,
		&entity.(receivedDate,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("inventory_transfers not found")
	}

	if err != nil {
		r.logger.Error("failed to get inventory_transfers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get inventory_transfers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of inventory_transfers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*InventoryTransfers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transfers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_transfers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_transfers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, transfer_number
			, transfer_date
			, from_location_id
			, to_location_id
			, status
			, requested_date
			, approved_date
			, shipped_date
			, expected_delivery_date
			, received_date
			, carrier
			, tracking_number
			, shipping_cost
			, reason
			, notes
			, rejection_reason
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, requested_by
			, approved_by
			, shipped_by
			, received_by
			, (approved_date
			, (shipped_date
			, (received_date
		FROM inventory_transfers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_transfers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_transfers: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryTransfers
	for rows.Next() {
		var entity InventoryTransfers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TransferNumber,
			&entity.TransferDate,
			&entity.FromLocationId,
			&entity.ToLocationId,
			&entity.Status,
			&entity.RequestedDate,
			&entity.ApprovedDate,
			&entity.ShippedDate,
			&entity.ExpectedDeliveryDate,
			&entity.ReceivedDate,
			&entity.Carrier,
			&entity.TrackingNumber,
			&entity.ShippingCost,
			&entity.Reason,
			&entity.Notes,
			&entity.RejectionReason,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.RequestedBy,
			&entity.ApprovedBy,
			&entity.ShippedBy,
			&entity.ReceivedBy,
			&entity.(approvedDate,
			&entity.(shippedDate,
			&entity.(receivedDate,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_transfers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating inventory_transfers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing inventory_transfers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *InventoryTransfers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_transfers", duration, nil)
	}()

	query := `
		UPDATE inventory_transfers
		SET
			, organization_id = $2
			, transfer_number = $3
			, transfer_date = $4
			, from_location_id = $5
			, to_location_id = $6
			, status = $7
			, requested_date = $8
			, approved_date = $9
			, shipped_date = $10
			, expected_delivery_date = $11
			, received_date = $12
			, carrier = $13
			, tracking_number = $14
			, shipping_cost = $15
			, reason = $16
			, notes = $17
			, rejection_reason = $18
			, metadata = $19
			, updated_at = $21
			, deleted_at = $22
			, created_by = $23
			, updated_by = $24
			, requested_by = $25
			, approved_by = $26
			, shipped_by = $27
			, received_by = $28
			, (approved_date = $29
			, (shipped_date = $30
			, (received_date = $31
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $32
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.TransferNumber,
		entity.TransferDate,
		entity.FromLocationId,
		entity.ToLocationId,
		entity.Status,
		entity.RequestedDate,
		entity.ApprovedDate,
		entity.ShippedDate,
		entity.ExpectedDeliveryDate,
		entity.ReceivedDate,
		entity.Carrier,
		entity.TrackingNumber,
		entity.ShippingCost,
		entity.Reason,
		entity.Notes,
		entity.RejectionReason,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.RequestedBy,
		entity.ApprovedBy,
		entity.ShippedBy,
		entity.ReceivedBy,
		entity.(approvedDate,
		entity.(shippedDate,
		entity.(receivedDate,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update inventory_transfers", zap.Error(err))
		return fmt.Errorf("failed to update inventory_transfers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_transfers not found or already deleted")
	}

	r.logger.Info("updated inventory_transfers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a inventory_transfers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_transfers", duration, nil)
	}()

	query := `
		UPDATE inventory_transfers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete inventory_transfers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete inventory_transfers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_transfers not found or already deleted")
	}

	r.logger.Info("deleted inventory_transfers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves inventory_transfers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*InventoryTransfers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transfers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_transfers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_transfers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, transfer_number
			, transfer_date
			, from_location_id
			, to_location_id
			, status
			, requested_date
			, approved_date
			, shipped_date
			, expected_delivery_date
			, received_date
			, carrier
			, tracking_number
			, shipping_cost
			, reason
			, notes
			, rejection_reason
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, requested_by
			, approved_by
			, shipped_by
			, received_by
			, (approved_date
			, (shipped_date
			, (received_date
		FROM inventory_transfers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_transfers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_transfers: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryTransfers
	for rows.Next() {
		var entity InventoryTransfers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TransferNumber,
			&entity.TransferDate,
			&entity.FromLocationId,
			&entity.ToLocationId,
			&entity.Status,
			&entity.RequestedDate,
			&entity.ApprovedDate,
			&entity.ShippedDate,
			&entity.ExpectedDeliveryDate,
			&entity.ReceivedDate,
			&entity.Carrier,
			&entity.TrackingNumber,
			&entity.ShippingCost,
			&entity.Reason,
			&entity.Notes,
			&entity.RejectionReason,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.RequestedBy,
			&entity.ApprovedBy,
			&entity.ShippedBy,
			&entity.ReceivedBy,
			&entity.(approvedDate,
			&entity.(shippedDate,
			&entity.(receivedDate,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_transfers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


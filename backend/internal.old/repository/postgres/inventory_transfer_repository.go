package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/inventory"
)

// InventoryTransferRepository implements inventory.TransferRepository
type InventoryTransferRepository struct {
	db *DB
}

// NewInventoryTransferRepository creates a new inventory transfer repository
func NewInventoryTransferRepository(db *DB) *InventoryTransferRepository {
	return &InventoryTransferRepository{db: db}
}

// ListTransfers retrieves inventory transfers with filters
func (r *InventoryTransferRepository) ListTransfers(ctx context.Context, orgID uuid.UUID, filters inventory.TransferFilters) ([]inventory.InventoryTransfer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, transfer_number, transfer_date, from_location_id, to_location_id,
		       status, requested_date, approved_date, shipped_date, expected_delivery_date, received_date,
		       carrier, tracking_number, shipping_cost, reason, notes, rejection_reason, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by, requested_by, approved_by, shipped_by, received_by
		FROM inventory_transfers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.FromLocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND from_location_id = $%d", argCount)
		args = append(args, *filters.FromLocationID)
	}

	if filters.ToLocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND to_location_id = $%d", argCount)
		args = append(args, *filters.ToLocationID)
	}

	if filters.TransferNumberLike != nil {
		argCount++
		query += fmt.Sprintf(" AND transfer_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.TransferNumberLike+"%")
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transfer_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transfer_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	if filters.CreatedBy != nil {
		argCount++
		query += fmt.Sprintf(" AND created_by = $%d", argCount)
		args = append(args, *filters.CreatedBy)
	}

	if filters.HasTracking != nil && *filters.HasTracking {
		query += " AND tracking_number IS NOT NULL"
	}

	query += " ORDER BY transfer_date DESC, created_at DESC"

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

	var transfers []inventory.InventoryTransfer
	for rows.Next() {
		var t inventory.InventoryTransfer
		var metadata interface{}
		err := rows.Scan(
			&t.ID, &t.OrganizationID, &t.TransferNumber, &t.TransferDate, &t.FromLocationID, &t.ToLocationID,
			&t.Status, &t.RequestedDate, &t.ApprovedDate, &t.ShippedDate, &t.ExpectedDeliveryDate, &t.ReceivedDate,
			&t.Carrier, &t.TrackingNumber, &t.ShippingCost, &t.Reason, &t.Notes, &t.RejectionReason, &metadata,
			&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.CreatedBy, &t.UpdatedBy, &t.RequestedBy, &t.ApprovedBy, &t.ShippedBy, &t.ReceivedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				t.Metadata = data
			}
		}
		transfers = append(transfers, t)
	}

	return transfers, rows.Err()
}

// CountTransfers counts inventory transfers with filters
func (r *InventoryTransferRepository) CountTransfers(ctx context.Context, orgID uuid.UUID, filters inventory.TransferFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM inventory_transfers WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.FromLocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND from_location_id = $%d", argCount)
		args = append(args, *filters.FromLocationID)
	}

	if filters.ToLocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND to_location_id = $%d", argCount)
		args = append(args, *filters.ToLocationID)
	}

	if filters.TransferNumberLike != nil {
		argCount++
		query += fmt.Sprintf(" AND transfer_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.TransferNumberLike+"%")
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transfer_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transfer_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	if filters.CreatedBy != nil {
		argCount++
		query += fmt.Sprintf(" AND created_by = $%d", argCount)
		args = append(args, *filters.CreatedBy)
	}

	if filters.HasTracking != nil && *filters.HasTracking {
		query += " AND tracking_number IS NOT NULL"
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetTransfer retrieves a single transfer by ID
func (r *InventoryTransferRepository) GetTransfer(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.InventoryTransfer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, transfer_number, transfer_date, from_location_id, to_location_id,
		       status, requested_date, approved_date, shipped_date, expected_delivery_date, received_date,
		       carrier, tracking_number, shipping_cost, reason, notes, rejection_reason, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by, requested_by, approved_by, shipped_by, received_by
		FROM inventory_transfers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var t inventory.InventoryTransfer
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&t.ID, &t.OrganizationID, &t.TransferNumber, &t.TransferDate, &t.FromLocationID, &t.ToLocationID,
		&t.Status, &t.RequestedDate, &t.ApprovedDate, &t.ShippedDate, &t.ExpectedDeliveryDate, &t.ReceivedDate,
		&t.Carrier, &t.TrackingNumber, &t.ShippingCost, &t.Reason, &t.Notes, &t.RejectionReason, &metadata,
		&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.CreatedBy, &t.UpdatedBy, &t.RequestedBy, &t.ApprovedBy, &t.ShippedBy, &t.ReceivedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			t.Metadata = data
		}
	}

	return &t, nil
}

// GetTransferByNumber retrieves a transfer by transfer number
func (r *InventoryTransferRepository) GetTransferByNumber(ctx context.Context, orgID uuid.UUID, transferNumber string) (*inventory.InventoryTransfer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, transfer_number, transfer_date, from_location_id, to_location_id,
		       status, requested_date, approved_date, shipped_date, expected_delivery_date, received_date,
		       carrier, tracking_number, shipping_cost, reason, notes, rejection_reason, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by, requested_by, approved_by, shipped_by, received_by
		FROM inventory_transfers
		WHERE organization_id = $1 AND transfer_number = $2 AND deleted_at IS NULL
	`

	var t inventory.InventoryTransfer
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, transferNumber).Scan(
		&t.ID, &t.OrganizationID, &t.TransferNumber, &t.TransferDate, &t.FromLocationID, &t.ToLocationID,
		&t.Status, &t.RequestedDate, &t.ApprovedDate, &t.ShippedDate, &t.ExpectedDeliveryDate, &t.ReceivedDate,
		&t.Carrier, &t.TrackingNumber, &t.ShippingCost, &t.Reason, &t.Notes, &t.RejectionReason, &metadata,
		&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.CreatedBy, &t.UpdatedBy, &t.RequestedBy, &t.ApprovedBy, &t.ShippedBy, &t.ReceivedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			t.Metadata = data
		}
	}

	return &t, nil
}

// CreateTransfer creates a new transfer
func (r *InventoryTransferRepository) CreateTransfer(ctx context.Context, transfer *inventory.InventoryTransfer) error {
	if err := r.db.SetOrganizationContext(ctx, transfer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_transfers (
			id, organization_id, transfer_number, transfer_date, from_location_id, to_location_id,
			status, requested_date, approved_date, shipped_date, expected_delivery_date, received_date,
			carrier, tracking_number, shipping_cost, reason, notes, rejection_reason, metadata,
			created_at, updated_at, created_by, updated_by, requested_by, approved_by, shipped_by, received_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23, $24, $25, $26, $27
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		transfer.ID, transfer.OrganizationID, transfer.TransferNumber, transfer.TransferDate, transfer.FromLocationID, transfer.ToLocationID,
		transfer.Status, transfer.RequestedDate, transfer.ApprovedDate, transfer.ShippedDate, transfer.ExpectedDeliveryDate, transfer.ReceivedDate,
		transfer.Carrier, transfer.TrackingNumber, transfer.ShippingCost, transfer.Reason, transfer.Notes, transfer.RejectionReason, transfer.Metadata,
		transfer.CreatedAt, transfer.UpdatedAt, transfer.CreatedBy, transfer.UpdatedBy, transfer.RequestedBy, transfer.ApprovedBy, transfer.ShippedBy, transfer.ReceivedBy,
	)

	return err
}

// UpdateTransfer updates an existing transfer
func (r *InventoryTransferRepository) UpdateTransfer(ctx context.Context, transfer *inventory.InventoryTransfer) error {
	if err := r.db.SetOrganizationContext(ctx, transfer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transfers SET
			transfer_number = $3, transfer_date = $4, from_location_id = $5, to_location_id = $6,
			status = $7, requested_date = $8, approved_date = $9, shipped_date = $10, expected_delivery_date = $11, received_date = $12,
			carrier = $13, tracking_number = $14, shipping_cost = $15, reason = $16, notes = $17, rejection_reason = $18, metadata = $19,
			updated_at = $20, updated_by = $21, requested_by = $22, approved_by = $23, shipped_by = $24, received_by = $25
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		transfer.OrganizationID, transfer.ID,
		transfer.TransferNumber, transfer.TransferDate, transfer.FromLocationID, transfer.ToLocationID,
		transfer.Status, transfer.RequestedDate, transfer.ApprovedDate, transfer.ShippedDate, transfer.ExpectedDeliveryDate, transfer.ReceivedDate,
		transfer.Carrier, transfer.TrackingNumber, transfer.ShippingCost, transfer.Reason, transfer.Notes, transfer.RejectionReason, transfer.Metadata,
		transfer.UpdatedAt, transfer.UpdatedBy, transfer.RequestedBy, transfer.ApprovedBy, transfer.ShippedBy, transfer.ReceivedBy,
	)

	return err
}

// DeleteTransfer soft deletes a transfer
func (r *InventoryTransferRepository) DeleteTransfer(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transfers SET deleted_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// UpdateTransferStatus updates the status of a transfer
func (r *InventoryTransferRepository) UpdateTransferStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status inventory.TransferStatus) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transfers SET status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status)
	return err
}

// UpdateTransferApproval updates the approval status of a transfer
func (r *InventoryTransferRepository) UpdateTransferApproval(ctx context.Context, orgID uuid.UUID, id uuid.UUID, approved bool, approvedBy uuid.UUID, approvalDate time.Time) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	var status inventory.TransferStatus
	if approved {
		status = inventory.TransferStatusApproved
	} else {
		status = inventory.TransferStatusRejected
	}

	query := `
		UPDATE inventory_transfers SET
			status = $3, approved_date = $4, approved_by = $5, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status, approvalDate, approvedBy)
	return err
}

// UpdateTransferShipping updates shipping information
func (r *InventoryTransferRepository) UpdateTransferShipping(ctx context.Context, orgID uuid.UUID, id uuid.UUID, carrier, trackingNumber *string, shippedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transfers SET
			status = $3, carrier = $4, tracking_number = $5, shipped_date = CURRENT_TIMESTAMP, shipped_by = $6, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, inventory.TransferStatusInTransit, carrier, trackingNumber, shippedBy)
	return err
}

// ReceiveTransfer marks a transfer as received
func (r *InventoryTransferRepository) ReceiveTransfer(ctx context.Context, orgID uuid.UUID, id uuid.UUID, receivedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transfers SET
			status = $3, received_date = CURRENT_TIMESTAMP, received_by = $4, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, inventory.TransferStatusReceived, receivedBy)
	return err
}

// ListTransferItems retrieves transfer items with filters
func (r *InventoryTransferRepository) ListTransferItems(ctx context.Context, orgID uuid.UUID, filters inventory.TransferItemFilters) ([]inventory.InventoryTransferItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, inventory_transfer_id, product_id, product_variant_id,
		       product_name, product_sku, quantity_requested, quantity_shipped, quantity_received,
		       unit_of_measure, unit_cost, total_cost, item_status, variance_quantity, variance_reason, notes, metadata,
		       created_at, updated_at
		FROM inventory_transfer_items
		WHERE organization_id = $1 AND inventory_transfer_id = $2
	`

	args := []interface{}{orgID, filters.TransferID}
	argCount := 2

	if filters.ItemStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND item_status = $%d", argCount)
		args = append(args, *filters.ItemStatus)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.HasVariance != nil && *filters.HasVariance {
		query += " AND variance_quantity != 0"
	}

	query += " ORDER BY created_at ASC"

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

	var items []inventory.InventoryTransferItem
	for rows.Next() {
		var item inventory.InventoryTransferItem
		var metadata interface{}
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.InventoryTransferID, &item.ProductID, &item.ProductVariantID,
			&item.ProductName, &item.ProductSKU, &item.QuantityRequested, &item.QuantityShipped, &item.QuantityReceived,
			&item.UnitOfMeasure, &item.UnitCost, &item.TotalCost, &item.ItemStatus, &item.VarianceQuantity, &item.VarianceReason, &item.Notes, &metadata,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				item.Metadata = data
			}
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// CountTransferItems counts transfer items
func (r *InventoryTransferRepository) CountTransferItems(ctx context.Context, orgID uuid.UUID, filters inventory.TransferItemFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM inventory_transfer_items WHERE organization_id = $1 AND inventory_transfer_id = $2"
	args := []interface{}{orgID, filters.TransferID}
	argCount := 2

	if filters.ItemStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND item_status = $%d", argCount)
		args = append(args, *filters.ItemStatus)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.HasVariance != nil && *filters.HasVariance {
		query += " AND variance_quantity != 0"
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetTransferItem retrieves a single transfer item
func (r *InventoryTransferRepository) GetTransferItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.InventoryTransferItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, inventory_transfer_id, product_id, product_variant_id,
		       product_name, product_sku, quantity_requested, quantity_shipped, quantity_received,
		       unit_of_measure, unit_cost, total_cost, item_status, variance_quantity, variance_reason, notes, metadata,
		       created_at, updated_at
		FROM inventory_transfer_items
		WHERE organization_id = $1 AND id = $2
	`

	var item inventory.InventoryTransferItem
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&item.ID, &item.OrganizationID, &item.InventoryTransferID, &item.ProductID, &item.ProductVariantID,
		&item.ProductName, &item.ProductSKU, &item.QuantityRequested, &item.QuantityShipped, &item.QuantityReceived,
		&item.UnitOfMeasure, &item.UnitCost, &item.TotalCost, &item.ItemStatus, &item.VarianceQuantity, &item.VarianceReason, &item.Notes, &metadata,
		&item.CreatedAt, &item.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			item.Metadata = data
		}
	}

	return &item, nil
}

// CreateTransferItem creates a new transfer item
func (r *InventoryTransferRepository) CreateTransferItem(ctx context.Context, item *inventory.InventoryTransferItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_transfer_items (
			id, organization_id, inventory_transfer_id, product_id, product_variant_id,
			product_name, product_sku, quantity_requested, quantity_shipped, quantity_received,
			unit_of_measure, unit_cost, total_cost, item_status, variance_quantity, variance_reason, notes, metadata,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.OrganizationID, item.InventoryTransferID, item.ProductID, item.ProductVariantID,
		item.ProductName, item.ProductSKU, item.QuantityRequested, item.QuantityShipped, item.QuantityReceived,
		item.UnitOfMeasure, item.UnitCost, item.TotalCost, item.ItemStatus, item.VarianceQuantity, item.VarianceReason, item.Notes, item.Metadata,
		item.CreatedAt, item.UpdatedAt,
	)

	return err
}

// UpdateTransferItem updates a transfer item
func (r *InventoryTransferRepository) UpdateTransferItem(ctx context.Context, item *inventory.InventoryTransferItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transfer_items SET
			product_id = $3, product_variant_id = $4, product_name = $5, product_sku = $6,
			quantity_requested = $7, quantity_shipped = $8, quantity_received = $9,
			unit_of_measure = $10, unit_cost = $11, total_cost = $12, item_status = $13,
			variance_quantity = $14, variance_reason = $15, notes = $16, metadata = $17,
			updated_at = $18
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.OrganizationID, item.ID,
		item.ProductID, item.ProductVariantID, item.ProductName, item.ProductSKU,
		item.QuantityRequested, item.QuantityShipped, item.QuantityReceived,
		item.UnitOfMeasure, item.UnitCost, item.TotalCost, item.ItemStatus,
		item.VarianceQuantity, item.VarianceReason, item.Notes, item.Metadata,
		item.UpdatedAt,
	)

	return err
}

// DeleteTransferItem deletes a transfer item
func (r *InventoryTransferRepository) DeleteTransferItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "DELETE FROM inventory_transfer_items WHERE organization_id = $1 AND id = $2"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// UpdateTransferItemReceipt updates the receipt information for a transfer item
func (r *InventoryTransferRepository) UpdateTransferItemReceipt(ctx context.Context, orgID uuid.UUID, id uuid.UUID, quantityReceived float64, varianceReason *string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transfer_items SET
			quantity_received = $3, variance_reason = $4, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, quantityReceived, varianceReason)
	return err
}

// GetTransferItemsByTransfer retrieves all items for a transfer
func (r *InventoryTransferRepository) GetTransferItemsByTransfer(ctx context.Context, orgID uuid.UUID, transferID uuid.UUID) ([]inventory.InventoryTransferItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, inventory_transfer_id, product_id, product_variant_id,
		       product_name, product_sku, quantity_requested, quantity_shipped, quantity_received,
		       unit_of_measure, unit_cost, total_cost, item_status, variance_quantity, variance_reason, notes, metadata,
		       created_at, updated_at
		FROM inventory_transfer_items
		WHERE organization_id = $1 AND inventory_transfer_id = $2
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []inventory.InventoryTransferItem
	for rows.Next() {
		var item inventory.InventoryTransferItem
		var metadata interface{}
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.InventoryTransferID, &item.ProductID, &item.ProductVariantID,
			&item.ProductName, &item.ProductSKU, &item.QuantityRequested, &item.QuantityShipped, &item.QuantityReceived,
			&item.UnitOfMeasure, &item.UnitCost, &item.TotalCost, &item.ItemStatus, &item.VarianceQuantity, &item.VarianceReason, &item.Notes, &metadata,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				item.Metadata = data
			}
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// CalculateTransferValue calculates the total value of a transfer
func (r *InventoryTransferRepository) CalculateTransferValue(ctx context.Context, orgID uuid.UUID, transferID uuid.UUID) (float64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := `
		SELECT COALESCE(SUM(total_cost), 0)
		FROM inventory_transfer_items
		WHERE organization_id = $1 AND inventory_transfer_id = $2
	`

	var value float64
	err := r.db.Pool.QueryRow(ctx, query, orgID, transferID).Scan(&value)
	return value, err
}

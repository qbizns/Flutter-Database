package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/inventory"
)

// InventoryBatchRepository implements inventory batch and serial number repositories
type InventoryBatchRepository struct {
	db *DB
}

// NewInventoryBatchRepository creates a new inventory batch repository
func NewInventoryBatchRepository(db *DB) *InventoryBatchRepository {
	return &InventoryBatchRepository{db: db}
}

// ============================================================================
// Serial Number Operations
// ============================================================================

// ListSerialNumbers retrieves serial numbers with filters
func (r *InventoryBatchRepository) ListSerialNumbers(ctx context.Context, orgID uuid.UUID, filters inventory.SerialNumberFilters) ([]inventory.ProductSerialNumber, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       serial_number, status, purchase_order_id, purchase_date, purchase_cost, supplier_id,
		       sale_id, sale_date, sale_price, customer_id, warranty_start_date, warranty_end_date,
		       warranty_provider, warranty_terms, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_serial_numbers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.SerialLike != nil {
		argCount++
		query += fmt.Sprintf(" AND serial_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.SerialLike+"%")
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.HasWarranty != nil && *filters.HasWarranty {
		query += " AND warranty_end_date IS NOT NULL"
	}

	if filters.WarrantyExpiring != nil && *filters.WarrantyExpiring {
		query += " AND warranty_end_date IS NOT NULL AND warranty_end_date <= CURRENT_DATE"
	}

	query += " ORDER BY serial_number ASC"

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

	var serials []inventory.ProductSerialNumber
	for rows.Next() {
		var s inventory.ProductSerialNumber
		var metadata interface{}
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ProductID, &s.ProductVariantID, &s.LocationID,
			&s.SerialNumber, &s.Status, &s.PurchaseOrderID, &s.PurchaseDate, &s.PurchaseCost, &s.SupplierID,
			&s.SaleID, &s.SaleDate, &s.SalePrice, &s.CustomerID, &s.WarrantyStartDate, &s.WarrantyEndDate,
			&s.WarrantyProvider, &s.WarrantyTerms, &s.Notes, &metadata,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.CreatedBy, &s.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				s.Metadata = data
			}
		}
		serials = append(serials, s)
	}

	return serials, rows.Err()
}

// CountSerialNumbers counts serial numbers
func (r *InventoryBatchRepository) CountSerialNumbers(ctx context.Context, orgID uuid.UUID, filters inventory.SerialNumberFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM product_serial_numbers WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.SerialLike != nil {
		argCount++
		query += fmt.Sprintf(" AND serial_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.SerialLike+"%")
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.HasWarranty != nil && *filters.HasWarranty {
		query += " AND warranty_end_date IS NOT NULL"
	}

	if filters.WarrantyExpiring != nil && *filters.WarrantyExpiring {
		query += " AND warranty_end_date IS NOT NULL AND warranty_end_date <= CURRENT_DATE"
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetSerialNumber retrieves a serial number by ID
func (r *InventoryBatchRepository) GetSerialNumber(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.ProductSerialNumber, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       serial_number, status, purchase_order_id, purchase_date, purchase_cost, supplier_id,
		       sale_id, sale_date, sale_price, customer_id, warranty_start_date, warranty_end_date,
		       warranty_provider, warranty_terms, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_serial_numbers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var s inventory.ProductSerialNumber
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&s.ID, &s.OrganizationID, &s.ProductID, &s.ProductVariantID, &s.LocationID,
		&s.SerialNumber, &s.Status, &s.PurchaseOrderID, &s.PurchaseDate, &s.PurchaseCost, &s.SupplierID,
		&s.SaleID, &s.SaleDate, &s.SalePrice, &s.CustomerID, &s.WarrantyStartDate, &s.WarrantyEndDate,
		&s.WarrantyProvider, &s.WarrantyTerms, &s.Notes, &metadata,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.CreatedBy, &s.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			s.Metadata = data
		}
	}

	return &s, nil
}

// GetBySerialNumber retrieves a serial number by its serial number string
func (r *InventoryBatchRepository) GetBySerialNumber(ctx context.Context, orgID uuid.UUID, serialNumber string) (*inventory.ProductSerialNumber, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       serial_number, status, purchase_order_id, purchase_date, purchase_cost, supplier_id,
		       sale_id, sale_date, sale_price, customer_id, warranty_start_date, warranty_end_date,
		       warranty_provider, warranty_terms, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_serial_numbers
		WHERE organization_id = $1 AND serial_number = $2 AND deleted_at IS NULL
	`

	var s inventory.ProductSerialNumber
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, serialNumber).Scan(
		&s.ID, &s.OrganizationID, &s.ProductID, &s.ProductVariantID, &s.LocationID,
		&s.SerialNumber, &s.Status, &s.PurchaseOrderID, &s.PurchaseDate, &s.PurchaseCost, &s.SupplierID,
		&s.SaleID, &s.SaleDate, &s.SalePrice, &s.CustomerID, &s.WarrantyStartDate, &s.WarrantyEndDate,
		&s.WarrantyProvider, &s.WarrantyTerms, &s.Notes, &metadata,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.CreatedBy, &s.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			s.Metadata = data
		}
	}

	return &s, nil
}

// CreateSerialNumber creates a new serial number
func (r *InventoryBatchRepository) CreateSerialNumber(ctx context.Context, serial *inventory.ProductSerialNumber) error {
	if err := r.db.SetOrganizationContext(ctx, serial.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO product_serial_numbers (
			id, organization_id, product_id, product_variant_id, location_id,
			serial_number, status, purchase_order_id, purchase_date, purchase_cost, supplier_id,
			sale_id, sale_date, sale_price, customer_id, warranty_start_date, warranty_end_date,
			warranty_provider, warranty_terms, notes, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21,
			$22, $23, $24, $25
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		serial.ID, serial.OrganizationID, serial.ProductID, serial.ProductVariantID, serial.LocationID,
		serial.SerialNumber, serial.Status, serial.PurchaseOrderID, serial.PurchaseDate, serial.PurchaseCost, serial.SupplierID,
		serial.SaleID, serial.SaleDate, serial.SalePrice, serial.CustomerID, serial.WarrantyStartDate, serial.WarrantyEndDate,
		serial.WarrantyProvider, serial.WarrantyTerms, serial.Notes, serial.Metadata,
		serial.CreatedAt, serial.UpdatedAt, serial.CreatedBy, serial.UpdatedBy,
	)

	return err
}

// UpdateSerialNumber updates a serial number
func (r *InventoryBatchRepository) UpdateSerialNumber(ctx context.Context, serial *inventory.ProductSerialNumber) error {
	if err := r.db.SetOrganizationContext(ctx, serial.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_serial_numbers SET
			product_variant_id = $3, location_id = $4,
			status = $5, purchase_order_id = $6, purchase_date = $7, purchase_cost = $8, supplier_id = $9,
			sale_id = $10, sale_date = $11, sale_price = $12, customer_id = $13,
			warranty_start_date = $14, warranty_end_date = $15, warranty_provider = $16, warranty_terms = $17,
			notes = $18, metadata = $19, updated_at = $20, updated_by = $21
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		serial.OrganizationID, serial.ID,
		serial.ProductVariantID, serial.LocationID,
		serial.Status, serial.PurchaseOrderID, serial.PurchaseDate, serial.PurchaseCost, serial.SupplierID,
		serial.SaleID, serial.SaleDate, serial.SalePrice, serial.CustomerID,
		serial.WarrantyStartDate, serial.WarrantyEndDate, serial.WarrantyProvider, serial.WarrantyTerms,
		serial.Notes, serial.Metadata, serial.UpdatedAt, serial.UpdatedBy,
	)

	return err
}

// DeleteSerialNumber soft deletes a serial number
func (r *InventoryBatchRepository) DeleteSerialNumber(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_serial_numbers SET deleted_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// UpdateSerialStatus updates the status of a serial number
func (r *InventoryBatchRepository) UpdateSerialStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status inventory.SerialNumberStatus) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_serial_numbers SET status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status)
	return err
}

// GetSerialsByProduct retrieves all serials for a product
func (r *InventoryBatchRepository) GetSerialsByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]inventory.ProductSerialNumber, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       serial_number, status, purchase_order_id, purchase_date, purchase_cost, supplier_id,
		       sale_id, sale_date, sale_price, customer_id, warranty_start_date, warranty_end_date,
		       warranty_provider, warranty_terms, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_serial_numbers
		WHERE organization_id = $1 AND product_id = $2 AND deleted_at IS NULL
		ORDER BY serial_number ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serials []inventory.ProductSerialNumber
	for rows.Next() {
		var s inventory.ProductSerialNumber
		var metadata interface{}
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ProductID, &s.ProductVariantID, &s.LocationID,
			&s.SerialNumber, &s.Status, &s.PurchaseOrderID, &s.PurchaseDate, &s.PurchaseCost, &s.SupplierID,
			&s.SaleID, &s.SaleDate, &s.SalePrice, &s.CustomerID, &s.WarrantyStartDate, &s.WarrantyEndDate,
			&s.WarrantyProvider, &s.WarrantyTerms, &s.Notes, &metadata,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.CreatedBy, &s.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				s.Metadata = data
			}
		}
		serials = append(serials, s)
	}

	return serials, rows.Err()
}

// GetExpiredWarrantySerials retrieves serials with expired warranties
func (r *InventoryBatchRepository) GetExpiredWarrantySerials(ctx context.Context, orgID uuid.UUID) ([]inventory.ProductSerialNumber, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       serial_number, status, purchase_order_id, purchase_date, purchase_cost, supplier_id,
		       sale_id, sale_date, sale_price, customer_id, warranty_start_date, warranty_end_date,
		       warranty_provider, warranty_terms, notes, metadata,
		       created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_serial_numbers
		WHERE organization_id = $1 AND warranty_end_date IS NOT NULL AND warranty_end_date <= CURRENT_DATE AND deleted_at IS NULL
		ORDER BY warranty_end_date ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serials []inventory.ProductSerialNumber
	for rows.Next() {
		var s inventory.ProductSerialNumber
		var metadata interface{}
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ProductID, &s.ProductVariantID, &s.LocationID,
			&s.SerialNumber, &s.Status, &s.PurchaseOrderID, &s.PurchaseDate, &s.PurchaseCost, &s.SupplierID,
			&s.SaleID, &s.SaleDate, &s.SalePrice, &s.CustomerID, &s.WarrantyStartDate, &s.WarrantyEndDate,
			&s.WarrantyProvider, &s.WarrantyTerms, &s.Notes, &metadata,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &s.CreatedBy, &s.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				s.Metadata = data
			}
		}
		serials = append(serials, s)
	}

	return serials, rows.Err()
}

// ============================================================================
// Batch Operations
// ============================================================================

// ListBatches retrieves product batches with filters
func (r *InventoryBatchRepository) ListBatches(ctx context.Context, orgID uuid.UUID, filters inventory.BatchFilters) ([]inventory.ProductBatch, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       batch_number, lot_number, status, initial_quantity, current_quantity, unit_of_measure,
		       manufacturing_date, expiration_date, received_date, purchase_order_id, supplier_id, supplier_batch_number,
		       unit_cost, total_cost, quality_status, quality_check_date, quality_checked_by, quality_notes,
		       notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_batches
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.BatchNumberLike != nil {
		argCount++
		query += fmt.Sprintf(" AND batch_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.BatchNumberLike+"%")
	}

	if filters.QualityStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND quality_status = $%d", argCount)
		args = append(args, *filters.QualityStatus)
	}

	if filters.ExpiringBefore != nil {
		argCount++
		query += fmt.Sprintf(" AND expiration_date IS NOT NULL AND expiration_date <= $%d", argCount)
		args = append(args, *filters.ExpiringBefore)
	}

	if filters.ExpiredAfter != nil {
		argCount++
		query += fmt.Sprintf(" AND expiration_date IS NULL OR expiration_date > $%d", argCount)
		args = append(args, *filters.ExpiredAfter)
	}

	if filters.HasQualityIssues != nil && *filters.HasQualityIssues {
		query += " AND quality_status IN ('failed', 'quarantine')"
	}

	if filters.MinStockLevel != nil {
		argCount++
		query += fmt.Sprintf(" AND current_quantity < $%d", argCount)
		args = append(args, *filters.MinStockLevel)
	}

	query += " ORDER BY batch_number ASC"

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

	var batches []inventory.ProductBatch
	for rows.Next() {
		var b inventory.ProductBatch
		var metadata interface{}
		err := rows.Scan(
			&b.ID, &b.OrganizationID, &b.ProductID, &b.ProductVariantID, &b.LocationID,
			&b.BatchNumber, &b.LotNumber, &b.Status, &b.InitialQuantity, &b.CurrentQuantity, &b.UnitOfMeasure,
			&b.ManufacturingDate, &b.ExpirationDate, &b.ReceivedDate, &b.PurchaseOrderID, &b.SupplierID, &b.SupplierBatchNumber,
			&b.UnitCost, &b.TotalCost, &b.QualityStatus, &b.QualityCheckDate, &b.QualityCheckedBy, &b.QualityNotes,
			&b.Notes, &metadata, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				b.Metadata = data
			}
		}
		batches = append(batches, b)
	}

	return batches, rows.Err()
}

// CountBatches counts product batches
func (r *InventoryBatchRepository) CountBatches(ctx context.Context, orgID uuid.UUID, filters inventory.BatchFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM product_batches WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.BatchNumberLike != nil {
		argCount++
		query += fmt.Sprintf(" AND batch_number ILIKE $%d", argCount)
		args = append(args, "%"+*filters.BatchNumberLike+"%")
	}

	if filters.QualityStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND quality_status = $%d", argCount)
		args = append(args, *filters.QualityStatus)
	}

	if filters.HasQualityIssues != nil && *filters.HasQualityIssues {
		query += " AND quality_status IN ('failed', 'quarantine')"
	}

	if filters.MinStockLevel != nil {
		argCount++
		query += fmt.Sprintf(" AND current_quantity < $%d", argCount)
		args = append(args, *filters.MinStockLevel)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetBatch retrieves a batch by ID
func (r *InventoryBatchRepository) GetBatch(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.ProductBatch, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       batch_number, lot_number, status, initial_quantity, current_quantity, unit_of_measure,
		       manufacturing_date, expiration_date, received_date, purchase_order_id, supplier_id, supplier_batch_number,
		       unit_cost, total_cost, quality_status, quality_check_date, quality_checked_by, quality_notes,
		       notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_batches
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var b inventory.ProductBatch
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&b.ID, &b.OrganizationID, &b.ProductID, &b.ProductVariantID, &b.LocationID,
		&b.BatchNumber, &b.LotNumber, &b.Status, &b.InitialQuantity, &b.CurrentQuantity, &b.UnitOfMeasure,
		&b.ManufacturingDate, &b.ExpirationDate, &b.ReceivedDate, &b.PurchaseOrderID, &b.SupplierID, &b.SupplierBatchNumber,
		&b.UnitCost, &b.TotalCost, &b.QualityStatus, &b.QualityCheckDate, &b.QualityCheckedBy, &b.QualityNotes,
		&b.Notes, &metadata, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			b.Metadata = data
		}
	}

	return &b, nil
}

// GetBatchByNumber retrieves a batch by batch number
func (r *InventoryBatchRepository) GetBatchByNumber(ctx context.Context, orgID uuid.UUID, batchNumber string) (*inventory.ProductBatch, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       batch_number, lot_number, status, initial_quantity, current_quantity, unit_of_measure,
		       manufacturing_date, expiration_date, received_date, purchase_order_id, supplier_id, supplier_batch_number,
		       unit_cost, total_cost, quality_status, quality_check_date, quality_checked_by, quality_notes,
		       notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_batches
		WHERE organization_id = $1 AND batch_number = $2 AND deleted_at IS NULL
	`

	var b inventory.ProductBatch
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, batchNumber).Scan(
		&b.ID, &b.OrganizationID, &b.ProductID, &b.ProductVariantID, &b.LocationID,
		&b.BatchNumber, &b.LotNumber, &b.Status, &b.InitialQuantity, &b.CurrentQuantity, &b.UnitOfMeasure,
		&b.ManufacturingDate, &b.ExpirationDate, &b.ReceivedDate, &b.PurchaseOrderID, &b.SupplierID, &b.SupplierBatchNumber,
		&b.UnitCost, &b.TotalCost, &b.QualityStatus, &b.QualityCheckDate, &b.QualityCheckedBy, &b.QualityNotes,
		&b.Notes, &metadata, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			b.Metadata = data
		}
	}

	return &b, nil
}

// CreateBatch creates a new batch
func (r *InventoryBatchRepository) CreateBatch(ctx context.Context, batch *inventory.ProductBatch) error {
	if err := r.db.SetOrganizationContext(ctx, batch.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO product_batches (
			id, organization_id, product_id, product_variant_id, location_id,
			batch_number, lot_number, status, initial_quantity, current_quantity, unit_of_measure,
			manufacturing_date, expiration_date, received_date, purchase_order_id, supplier_id, supplier_batch_number,
			unit_cost, total_cost, quality_status, quality_check_date, quality_checked_by, quality_notes,
			notes, metadata, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23,
			$24, $25, $26, $27, $28, $29
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		batch.ID, batch.OrganizationID, batch.ProductID, batch.ProductVariantID, batch.LocationID,
		batch.BatchNumber, batch.LotNumber, batch.Status, batch.InitialQuantity, batch.CurrentQuantity, batch.UnitOfMeasure,
		batch.ManufacturingDate, batch.ExpirationDate, batch.ReceivedDate, batch.PurchaseOrderID, batch.SupplierID, batch.SupplierBatchNumber,
		batch.UnitCost, batch.TotalCost, batch.QualityStatus, batch.QualityCheckDate, batch.QualityCheckedBy, batch.QualityNotes,
		batch.Notes, batch.Metadata, batch.CreatedAt, batch.UpdatedAt, batch.CreatedBy, batch.UpdatedBy,
	)

	return err
}

// UpdateBatch updates a batch
func (r *InventoryBatchRepository) UpdateBatch(ctx context.Context, batch *inventory.ProductBatch) error {
	if err := r.db.SetOrganizationContext(ctx, batch.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_batches SET
			product_variant_id = $3, location_id = $4,
			batch_number = $5, lot_number = $6, status = $7, initial_quantity = $8, current_quantity = $9, unit_of_measure = $10,
			manufacturing_date = $11, expiration_date = $12, received_date = $13, purchase_order_id = $14, supplier_id = $15, supplier_batch_number = $16,
			unit_cost = $17, total_cost = $18, quality_status = $19, quality_check_date = $20, quality_checked_by = $21, quality_notes = $22,
			notes = $23, metadata = $24, updated_at = $25, updated_by = $26
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		batch.OrganizationID, batch.ID,
		batch.ProductVariantID, batch.LocationID,
		batch.BatchNumber, batch.LotNumber, batch.Status, batch.InitialQuantity, batch.CurrentQuantity, batch.UnitOfMeasure,
		batch.ManufacturingDate, batch.ExpirationDate, batch.ReceivedDate, batch.PurchaseOrderID, batch.SupplierID, batch.SupplierBatchNumber,
		batch.UnitCost, batch.TotalCost, batch.QualityStatus, batch.QualityCheckDate, batch.QualityCheckedBy, batch.QualityNotes,
		batch.Notes, batch.Metadata, batch.UpdatedAt, batch.UpdatedBy,
	)

	return err
}

// DeleteBatch soft deletes a batch
func (r *InventoryBatchRepository) DeleteBatch(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_batches SET deleted_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// UpdateBatchQuantity updates the current quantity of a batch
func (r *InventoryBatchRepository) UpdateBatchQuantity(ctx context.Context, orgID uuid.UUID, id uuid.UUID, quantity float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_batches SET
			current_quantity = $3, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, quantity)
	return err
}

// UpdateBatchStatus updates the status of a batch
func (r *InventoryBatchRepository) UpdateBatchStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status inventory.BatchStatus) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_batches SET
			status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status)
	return err
}

// UpdateBatchQuality updates the quality status of a batch
func (r *InventoryBatchRepository) UpdateBatchQuality(ctx context.Context, orgID uuid.UUID, id uuid.UUID, qualityStatus inventory.QualityStatus, checkedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_batches SET
			quality_status = $3, quality_check_date = CURRENT_DATE, quality_checked_by = $4, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, qualityStatus, checkedBy)
	return err
}

// GetExpiringBatches retrieves batches expiring within specified days
func (r *InventoryBatchRepository) GetExpiringBatches(ctx context.Context, orgID uuid.UUID, daysUntilExpiry int) ([]inventory.ProductBatch, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       batch_number, lot_number, status, initial_quantity, current_quantity, unit_of_measure,
		       manufacturing_date, expiration_date, received_date, purchase_order_id, supplier_id, supplier_batch_number,
		       unit_cost, total_cost, quality_status, quality_check_date, quality_checked_by, quality_notes,
		       notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_batches
		WHERE organization_id = $1 AND status = 'active' AND deleted_at IS NULL
		AND expiration_date IS NOT NULL AND expiration_date <= CURRENT_DATE + INTERVAL '1 day' * $2
		ORDER BY expiration_date ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, daysUntilExpiry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []inventory.ProductBatch
	for rows.Next() {
		var b inventory.ProductBatch
		var metadata interface{}
		err := rows.Scan(
			&b.ID, &b.OrganizationID, &b.ProductID, &b.ProductVariantID, &b.LocationID,
			&b.BatchNumber, &b.LotNumber, &b.Status, &b.InitialQuantity, &b.CurrentQuantity, &b.UnitOfMeasure,
			&b.ManufacturingDate, &b.ExpirationDate, &b.ReceivedDate, &b.PurchaseOrderID, &b.SupplierID, &b.SupplierBatchNumber,
			&b.UnitCost, &b.TotalCost, &b.QualityStatus, &b.QualityCheckDate, &b.QualityCheckedBy, &b.QualityNotes,
			&b.Notes, &metadata, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				b.Metadata = data
			}
		}
		batches = append(batches, b)
	}

	return batches, rows.Err()
}

// GetBatchesByProduct retrieves all batches for a product
func (r *InventoryBatchRepository) GetBatchesByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]inventory.ProductBatch, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, product_variant_id, location_id,
		       batch_number, lot_number, status, initial_quantity, current_quantity, unit_of_measure,
		       manufacturing_date, expiration_date, received_date, purchase_order_id, supplier_id, supplier_batch_number,
		       unit_cost, total_cost, quality_status, quality_check_date, quality_checked_by, quality_notes,
		       notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_batches
		WHERE organization_id = $1 AND product_id = $2 AND deleted_at IS NULL
		ORDER BY batch_number ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []inventory.ProductBatch
	for rows.Next() {
		var b inventory.ProductBatch
		var metadata interface{}
		err := rows.Scan(
			&b.ID, &b.OrganizationID, &b.ProductID, &b.ProductVariantID, &b.LocationID,
			&b.BatchNumber, &b.LotNumber, &b.Status, &b.InitialQuantity, &b.CurrentQuantity, &b.UnitOfMeasure,
			&b.ManufacturingDate, &b.ExpirationDate, &b.ReceivedDate, &b.PurchaseOrderID, &b.SupplierID, &b.SupplierBatchNumber,
			&b.UnitCost, &b.TotalCost, &b.QualityStatus, &b.QualityCheckDate, &b.QualityCheckedBy, &b.QualityNotes,
			&b.Notes, &metadata, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				b.Metadata = data
			}
		}
		batches = append(batches, b)
	}

	return batches, rows.Err()
}

// ============================================================================
// Batch Transaction Operations
// ============================================================================

// ListBatchTransactions retrieves batch transactions
func (r *InventoryBatchRepository) ListBatchTransactions(ctx context.Context, orgID uuid.UUID, filters inventory.BatchTransactionFilters) ([]inventory.BatchTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, batch_id, transaction_type, quantity, balance_after,
		       sale_id, inventory_transfer_id, reason, notes, metadata, transaction_date, created_by
		FROM batch_transactions
		WHERE organization_id = $1 AND batch_id = $2
	`

	args := []interface{}{orgID, filters.BatchID}
	argCount := 2

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, *filters.TransactionType)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY transaction_date DESC"

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

	var txs []inventory.BatchTransaction
	for rows.Next() {
		var tx inventory.BatchTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.BatchID, &tx.TransactionType, &tx.Quantity, &tx.BalanceAfter,
			&tx.SaleID, &tx.InventoryTransferID, &tx.Reason, &tx.Notes, &metadata, &tx.TransactionDate, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				tx.Metadata = data
			}
		}
		txs = append(txs, tx)
	}

	return txs, rows.Err()
}

// CountBatchTransactions counts batch transactions
func (r *InventoryBatchRepository) CountBatchTransactions(ctx context.Context, orgID uuid.UUID, filters inventory.BatchTransactionFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM batch_transactions WHERE organization_id = $1 AND batch_id = $2"
	args := []interface{}{orgID, filters.BatchID}
	argCount := 2

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, *filters.TransactionType)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// GetBatchTransaction retrieves a batch transaction
func (r *InventoryBatchRepository) GetBatchTransaction(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.BatchTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, batch_id, transaction_type, quantity, balance_after,
		       sale_id, inventory_transfer_id, reason, notes, metadata, transaction_date, created_by
		FROM batch_transactions
		WHERE organization_id = $1 AND id = $2
	`

	var tx inventory.BatchTransaction
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&tx.ID, &tx.OrganizationID, &tx.BatchID, &tx.TransactionType, &tx.Quantity, &tx.BalanceAfter,
		&tx.SaleID, &tx.InventoryTransferID, &tx.Reason, &tx.Notes, &metadata, &tx.TransactionDate, &tx.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		if data, ok := metadata.([]byte); ok {
			tx.Metadata = data
		}
	}

	return &tx, nil
}

// CreateBatchTransaction creates a batch transaction
func (r *InventoryBatchRepository) CreateBatchTransaction(ctx context.Context, tx *inventory.BatchTransaction) error {
	if err := r.db.SetOrganizationContext(ctx, tx.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO batch_transactions (
			id, organization_id, batch_id, transaction_type, quantity, balance_after,
			sale_id, inventory_transfer_id, reason, notes, metadata, transaction_date, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		tx.ID, tx.OrganizationID, tx.BatchID, tx.TransactionType, tx.Quantity, tx.BalanceAfter,
		tx.SaleID, tx.InventoryTransferID, tx.Reason, tx.Notes, tx.Metadata, tx.TransactionDate, tx.CreatedBy,
	)

	return err
}

// GetBatchTransactionHistory retrieves the transaction history of a batch
func (r *InventoryBatchRepository) GetBatchTransactionHistory(ctx context.Context, orgID uuid.UUID, batchID uuid.UUID) ([]inventory.BatchTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, batch_id, transaction_type, quantity, balance_after,
		       sale_id, inventory_transfer_id, reason, notes, metadata, transaction_date, created_by
		FROM batch_transactions
		WHERE organization_id = $1 AND batch_id = $2
		ORDER BY transaction_date ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, batchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []inventory.BatchTransaction
	for rows.Next() {
		var tx inventory.BatchTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.BatchID, &tx.TransactionType, &tx.Quantity, &tx.BalanceAfter,
			&tx.SaleID, &tx.InventoryTransferID, &tx.Reason, &tx.Notes, &metadata, &tx.TransactionDate, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			if data, ok := metadata.([]byte); ok {
				tx.Metadata = data
			}
		}
		txs = append(txs, tx)
	}

	return txs, rows.Err()
}

// CalculateBatchCurrentQuantity calculates the current quantity of a batch from transactions
func (r *InventoryBatchRepository) CalculateBatchCurrentQuantity(ctx context.Context, orgID uuid.UUID, batchID uuid.UUID) (float64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := `
		SELECT COALESCE(initial_quantity, 0) - COALESCE(SUM(CASE WHEN bt.transaction_type IN ('sale', 'waste', 'return', 'expiration') THEN bt.quantity ELSE 0 END), 0)
		FROM product_batches pb
		LEFT JOIN batch_transactions bt ON pb.id = bt.batch_id
		WHERE pb.organization_id = $1 AND pb.id = $2
		GROUP BY pb.initial_quantity
	`

	var quantity float64
	err := r.db.Pool.QueryRow(ctx, query, orgID, batchID).Scan(&quantity)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	return quantity, err
}

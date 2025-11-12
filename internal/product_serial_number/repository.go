package product_serial_number

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

// Repository handles database operations for ProductSerialNumbers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ProductSerialNumbers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ProductSerialNumbers represents a product_serial_numbers entity
type ProductSerialNumbers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	SerialNumber string `json:"serial_number" db:"serial_number"`
	Status *string `json:"status" db:"status"`
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id" db:"purchase_order_id"`
	PurchaseDate *time.Time `json:"purchase_date" db:"purchase_date"`
	PurchaseCost *float64 `json:"purchase_cost" db:"purchase_cost"`
	SupplierId *uuid.UUID `json:"supplier_id" db:"supplier_id"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	SaleDate *time.Time `json:"sale_date" db:"sale_date"`
	SalePrice *float64 `json:"sale_price" db:"sale_price"`
	CustomerId *uuid.UUID `json:"customer_id" db:"customer_id"`
	WarrantyStartDate *time.Time `json:"warranty_start_date" db:"warranty_start_date"`
	WarrantyEndDate *time.Time `json:"warranty_end_date" db:"warranty_end_date"`
	WarrantyProvider *string `json:"warranty_provider" db:"warranty_provider"`
	WarrantyTerms *string `json:"warranty_terms" db:"warranty_terms"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new product_serial_numbers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ProductSerialNumbers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "product_serial_numbers", duration, nil)
	}()

	query := `
		INSERT INTO product_serial_numbers (
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, serial_number
			, status
			, purchase_order_id
			, purchase_date
			, purchase_cost
			, supplier_id
			, sale_id
			, sale_date
			, sale_price
			, customer_id
			, warranty_start_date
			, warranty_end_date
			, warranty_provider
			, warranty_terms
			, notes
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
			, $24
			, $25
			, $26
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.LocationId,
		entity.SerialNumber,
		entity.Status,
		entity.PurchaseOrderId,
		entity.PurchaseDate,
		entity.PurchaseCost,
		entity.SupplierId,
		entity.SaleId,
		entity.SaleDate,
		entity.SalePrice,
		entity.CustomerId,
		entity.WarrantyStartDate,
		entity.WarrantyEndDate,
		entity.WarrantyProvider,
		entity.WarrantyTerms,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create product_serial_numbers", zap.Error(err))
		return fmt.Errorf("failed to create product_serial_numbers: %w", err)
	}

	r.logger.Info("created product_serial_numbers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a product_serial_numbers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ProductSerialNumbers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_serial_numbers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, serial_number
			, status
			, purchase_order_id
			, purchase_date
			, purchase_cost
			, supplier_id
			, sale_id
			, sale_date
			, sale_price
			, customer_id
			, warranty_start_date
			, warranty_end_date
			, warranty_provider
			, warranty_terms
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM product_serial_numbers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ProductSerialNumbers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.LocationId,
		&entity.SerialNumber,
		&entity.Status,
		&entity.PurchaseOrderId,
		&entity.PurchaseDate,
		&entity.PurchaseCost,
		&entity.SupplierId,
		&entity.SaleId,
		&entity.SaleDate,
		&entity.SalePrice,
		&entity.CustomerId,
		&entity.WarrantyStartDate,
		&entity.WarrantyEndDate,
		&entity.WarrantyProvider,
		&entity.WarrantyTerms,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("product_serial_numbers not found")
	}

	if err != nil {
		r.logger.Error("failed to get product_serial_numbers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get product_serial_numbers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of product_serial_numbers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ProductSerialNumbers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_serial_numbers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_serial_numbers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_serial_numbers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, serial_number
			, status
			, purchase_order_id
			, purchase_date
			, purchase_cost
			, supplier_id
			, sale_id
			, sale_date
			, sale_price
			, customer_id
			, warranty_start_date
			, warranty_end_date
			, warranty_provider
			, warranty_terms
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM product_serial_numbers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_serial_numbers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_serial_numbers: %w", err)
	}
	defer rows.Close()

	var entities []*ProductSerialNumbers
	for rows.Next() {
		var entity ProductSerialNumbers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.LocationId,
			&entity.SerialNumber,
			&entity.Status,
			&entity.PurchaseOrderId,
			&entity.PurchaseDate,
			&entity.PurchaseCost,
			&entity.SupplierId,
			&entity.SaleId,
			&entity.SaleDate,
			&entity.SalePrice,
			&entity.CustomerId,
			&entity.WarrantyStartDate,
			&entity.WarrantyEndDate,
			&entity.WarrantyProvider,
			&entity.WarrantyTerms,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_serial_numbers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating product_serial_numbers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing product_serial_numbers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ProductSerialNumbers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_serial_numbers", duration, nil)
	}()

	query := `
		UPDATE product_serial_numbers
		SET
			, organization_id = $2
			, product_id = $3
			, product_variant_id = $4
			, location_id = $5
			, serial_number = $6
			, status = $7
			, purchase_order_id = $8
			, purchase_date = $9
			, purchase_cost = $10
			, supplier_id = $11
			, sale_id = $12
			, sale_date = $13
			, sale_price = $14
			, customer_id = $15
			, warranty_start_date = $16
			, warranty_end_date = $17
			, warranty_provider = $18
			, warranty_terms = $19
			, notes = $20
			, metadata = $21
			, updated_at = $23
			, deleted_at = $24
			, created_by = $25
			, updated_by = $26
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $27
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.LocationId,
		entity.SerialNumber,
		entity.Status,
		entity.PurchaseOrderId,
		entity.PurchaseDate,
		entity.PurchaseCost,
		entity.SupplierId,
		entity.SaleId,
		entity.SaleDate,
		entity.SalePrice,
		entity.CustomerId,
		entity.WarrantyStartDate,
		entity.WarrantyEndDate,
		entity.WarrantyProvider,
		entity.WarrantyTerms,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update product_serial_numbers", zap.Error(err))
		return fmt.Errorf("failed to update product_serial_numbers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_serial_numbers not found or already deleted")
	}

	r.logger.Info("updated product_serial_numbers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a product_serial_numbers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_serial_numbers", duration, nil)
	}()

	query := `
		UPDATE product_serial_numbers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete product_serial_numbers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete product_serial_numbers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_serial_numbers not found or already deleted")
	}

	r.logger.Info("deleted product_serial_numbers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves product_serial_numbers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ProductSerialNumbers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_serial_numbers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_serial_numbers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_serial_numbers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, product_variant_id
			, location_id
			, serial_number
			, status
			, purchase_order_id
			, purchase_date
			, purchase_cost
			, supplier_id
			, sale_id
			, sale_date
			, sale_price
			, customer_id
			, warranty_start_date
			, warranty_end_date
			, warranty_provider
			, warranty_terms
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM product_serial_numbers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_serial_numbers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_serial_numbers: %w", err)
	}
	defer rows.Close()

	var entities []*ProductSerialNumbers
	for rows.Next() {
		var entity ProductSerialNumbers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.LocationId,
			&entity.SerialNumber,
			&entity.Status,
			&entity.PurchaseOrderId,
			&entity.PurchaseDate,
			&entity.PurchaseCost,
			&entity.SupplierId,
			&entity.SaleId,
			&entity.SaleDate,
			&entity.SalePrice,
			&entity.CustomerId,
			&entity.WarrantyStartDate,
			&entity.WarrantyEndDate,
			&entity.WarrantyProvider,
			&entity.WarrantyTerms,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_serial_numbers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


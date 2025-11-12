package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/products"
)

// ProductVariantRepository implements products.VariantRepository
type ProductVariantRepository struct {
	db *DB
}

// NewProductVariantRepository creates a new product variant repository
func NewProductVariantRepository(db *DB) *ProductVariantRepository {
	return &ProductVariantRepository{db: db}
}

// Create creates a new product variant
func (r *ProductVariantRepository) Create(ctx context.Context, variant *products.ProductVariant) error {
	if err := r.db.SetOrganizationContext(ctx, variant.OrganizationID.String()); err != nil {
		return err
	}

	attributesJSON, _ := json.Marshal(variant.Attributes)
	dimensionsJSON, _ := json.Marshal(variant.Dimensions)

	query := `
		INSERT INTO product_variants (
			id, organization_id, product_id, variant_name, sku, barcode,
			attributes, cost_price, selling_price, compare_at_price,
			current_stock, reorder_level, reorder_quantity,
			weight, weight_unit, dimensions,
			is_active, is_default, sort_order, image_url, notes,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21,
			$22, $23, $24
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		variant.ID, variant.OrganizationID, variant.ProductID, variant.VariantName, variant.SKU, variant.Barcode,
		attributesJSON, variant.CostPrice, variant.SellingPrice, variant.CompareAtPrice,
		variant.CurrentStock, variant.ReorderLevel, variant.ReorderQuantity,
		variant.Weight, variant.WeightUnit, dimensionsJSON,
		variant.IsActive, variant.IsDefault, variant.SortOrder, variant.ImageURL, variant.Notes,
		variant.CreatedAt, variant.UpdatedAt, variant.CreatedBy,
	)

	return err
}

// Get retrieves a product variant by ID
func (r *ProductVariantRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*products.ProductVariant, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, product_id, variant_name, sku, barcode,
			attributes, cost_price, selling_price, compare_at_price,
			current_stock, reorder_level, reorder_quantity,
			weight, weight_unit, dimensions,
			is_active, is_default, sort_order, image_url, notes,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_variants
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var v products.ProductVariant
	var attributesJSON []byte
	var dimensionsJSON []byte

	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&v.ID, &v.OrganizationID, &v.ProductID, &v.VariantName, &v.SKU, &v.Barcode,
		&attributesJSON, &v.CostPrice, &v.SellingPrice, &v.CompareAtPrice,
		&v.CurrentStock, &v.ReorderLevel, &v.ReorderQuantity,
		&v.Weight, &v.WeightUnit, &dimensionsJSON,
		&v.IsActive, &v.IsDefault, &v.SortOrder, &v.ImageURL, &v.Notes,
		&v.CreatedAt, &v.UpdatedAt, &v.DeletedAt, &v.CreatedBy, &v.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	if len(attributesJSON) > 0 {
		v.Attributes = products.VariantAttributes{}
		json.Unmarshal(attributesJSON, &v.Attributes)
	}
	if len(dimensionsJSON) > 0 {
		v.Dimensions = &products.VariantDimensions{}
		json.Unmarshal(dimensionsJSON, v.Dimensions)
	}

	return &v, nil
}

// GetBySKU retrieves a variant by SKU
func (r *ProductVariantRepository) GetBySKU(ctx context.Context, orgID uuid.UUID, sku string) (*products.ProductVariant, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, product_id, variant_name, sku, barcode,
			attributes, cost_price, selling_price, compare_at_price,
			current_stock, reorder_level, reorder_quantity,
			weight, weight_unit, dimensions,
			is_active, is_default, sort_order, image_url, notes,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_variants
		WHERE organization_id = $1 AND LOWER(sku) = LOWER($2) AND deleted_at IS NULL
	`

	var v products.ProductVariant
	var attributesJSON []byte
	var dimensionsJSON []byte

	err := r.db.Pool.QueryRow(ctx, query, orgID, sku).Scan(
		&v.ID, &v.OrganizationID, &v.ProductID, &v.VariantName, &v.SKU, &v.Barcode,
		&attributesJSON, &v.CostPrice, &v.SellingPrice, &v.CompareAtPrice,
		&v.CurrentStock, &v.ReorderLevel, &v.ReorderQuantity,
		&v.Weight, &v.WeightUnit, &dimensionsJSON,
		&v.IsActive, &v.IsDefault, &v.SortOrder, &v.ImageURL, &v.Notes,
		&v.CreatedAt, &v.UpdatedAt, &v.DeletedAt, &v.CreatedBy, &v.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	if len(attributesJSON) > 0 {
		v.Attributes = products.VariantAttributes{}
		json.Unmarshal(attributesJSON, &v.Attributes)
	}
	if len(dimensionsJSON) > 0 {
		v.Dimensions = &products.VariantDimensions{}
		json.Unmarshal(dimensionsJSON, v.Dimensions)
	}

	return &v, nil
}

// Update updates a product variant
func (r *ProductVariantRepository) Update(ctx context.Context, variant *products.ProductVariant) error {
	if err := r.db.SetOrganizationContext(ctx, variant.OrganizationID.String()); err != nil {
		return err
	}

	attributesJSON, _ := json.Marshal(variant.Attributes)
	dimensionsJSON, _ := json.Marshal(variant.Dimensions)

	query := `
		UPDATE product_variants SET
			variant_name = $3, sku = $4, barcode = $5,
			attributes = $6, cost_price = $7, selling_price = $8, compare_at_price = $9,
			reorder_level = $10, reorder_quantity = $11,
			weight = $12, weight_unit = $13, dimensions = $14,
			is_active = $15, is_default = $16, sort_order = $17, image_url = $18, notes = $19,
			updated_at = $20, updated_by = $21
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		variant.OrganizationID, variant.ID,
		variant.VariantName, variant.SKU, variant.Barcode,
		attributesJSON, variant.CostPrice, variant.SellingPrice, variant.CompareAtPrice,
		variant.ReorderLevel, variant.ReorderQuantity,
		variant.Weight, variant.WeightUnit, dimensionsJSON,
		variant.IsActive, variant.IsDefault, variant.SortOrder, variant.ImageURL, variant.Notes,
		variant.UpdatedAt, variant.UpdatedBy,
	)

	return err
}

// Delete soft-deletes a variant
func (r *ProductVariantRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_variants
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// List retrieves variants with filters
func (r *ProductVariantRepository) List(ctx context.Context, orgID uuid.UUID, filters products.VariantFilters) ([]products.ProductVariant, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, product_id, variant_name, sku, barcode,
			attributes, cost_price, selling_price, compare_at_price,
			current_stock, reorder_level, reorder_quantity,
			weight, weight_unit, dimensions,
			is_active, is_default, sort_order, image_url, notes,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM product_variants
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	// Apply filters
	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (variant_name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	// Add ordering
	query += " ORDER BY sort_order ASC, variant_name ASC"

	// Add pagination
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

	var variants []products.ProductVariant
	for rows.Next() {
		var v products.ProductVariant
		var attributesJSON []byte
		var dimensionsJSON []byte

		err := rows.Scan(
			&v.ID, &v.OrganizationID, &v.ProductID, &v.VariantName, &v.SKU, &v.Barcode,
			&attributesJSON, &v.CostPrice, &v.SellingPrice, &v.CompareAtPrice,
			&v.CurrentStock, &v.ReorderLevel, &v.ReorderQuantity,
			&v.Weight, &v.WeightUnit, &dimensionsJSON,
			&v.IsActive, &v.IsDefault, &v.SortOrder, &v.ImageURL, &v.Notes,
			&v.CreatedAt, &v.UpdatedAt, &v.DeletedAt, &v.CreatedBy, &v.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if len(attributesJSON) > 0 {
			v.Attributes = products.VariantAttributes{}
			json.Unmarshal(attributesJSON, &v.Attributes)
		}
		if len(dimensionsJSON) > 0 {
			v.Dimensions = &products.VariantDimensions{}
			json.Unmarshal(dimensionsJSON, v.Dimensions)
		}

		variants = append(variants, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return variants, nil
}

// Count counts variants matching filters
func (r *ProductVariantRepository) Count(ctx context.Context, orgID uuid.UUID, filters products.VariantFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM product_variants WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (variant_name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// ListByProductID retrieves all variants for a product
func (r *ProductVariantRepository) ListByProductID(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]products.ProductVariant, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	filters := products.VariantFilters{ProductID: &productID}
	return r.List(ctx, orgID, filters)
}

// DeleteByProductID deletes all variants for a product
func (r *ProductVariantRepository) DeleteByProductID(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_variants
		SET deleted_at = $3
		WHERE organization_id = $1 AND product_id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, productID, time.Now())
	return err
}

// UpdateStock sets stock quantity
func (r *ProductVariantRepository) UpdateStock(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID, quantity float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_variants
		SET current_stock = $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, variantID, quantity, time.Now())
	return err
}

// AdjustStock adjusts stock by a delta
func (r *ProductVariantRepository) AdjustStock(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID, quantity float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_variants
		SET current_stock = current_stock + $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, variantID, quantity, time.Now())
	return err
}

// SetDefault sets a variant as default for its product
func (r *ProductVariantRepository) SetDefault(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	// Start a transaction to ensure atomicity
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get the variant to find its product
	var productID uuid.UUID
	err = tx.QueryRow(ctx, "SELECT product_id FROM product_variants WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL", orgID, variantID).Scan(&productID)
	if err != nil {
		return err
	}

	// Clear default for all variants of this product
	_, err = tx.Exec(ctx, "UPDATE product_variants SET is_default = false WHERE organization_id = $1 AND product_id = $2 AND deleted_at IS NULL", orgID, productID)
	if err != nil {
		return err
	}

	// Set this variant as default
	_, err = tx.Exec(ctx, "UPDATE product_variants SET is_default = true, updated_at = $3 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL", orgID, variantID, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

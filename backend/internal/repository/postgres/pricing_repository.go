package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/products"
)

// PriceListRepository implements products.PriceListRepository
type PriceListRepository struct {
	db *DB
}

// NewPriceListRepository creates a new pricing repository
func NewPriceListRepository(db *DB) *PriceListRepository {
	return &PriceListRepository{db: db}
}

// CreatePriceList creates a new price list
func (r *PriceListRepository) CreatePriceList(ctx context.Context, pl *products.PriceList) error {
	if err := r.db.SetOrganizationContext(ctx, pl.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO price_lists (
			id, organization_id, price_list_code, price_list_name, price_list_type,
			effective_from, effective_to, base_price_adjustment_type,
			base_price_adjustment_value, priority, is_active, description,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		pl.ID, pl.OrganizationID, pl.Code, pl.Name, pl.Type,
		pl.EffectiveFrom, pl.EffectiveTo, pl.BaseAdjustmentType,
		pl.BaseAdjustmentValue, pl.Priority, pl.IsActive, pl.Description,
		pl.CreatedAt, pl.UpdatedAt, pl.CreatedBy,
	)

	return err
}

// GetPriceList retrieves a price list by ID
func (r *PriceListRepository) GetPriceList(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*products.PriceList, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, price_list_code, price_list_name, price_list_type,
			effective_from, effective_to, base_price_adjustment_type,
			base_price_adjustment_value, priority, is_active, description,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM price_lists
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var pl products.PriceList
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&pl.ID, &pl.OrganizationID, &pl.Code, &pl.Name, &pl.Type,
		&pl.EffectiveFrom, &pl.EffectiveTo, &pl.BaseAdjustmentType,
		&pl.BaseAdjustmentValue, &pl.Priority, &pl.IsActive, &pl.Description,
		&pl.CreatedAt, &pl.UpdatedAt, &pl.DeletedAt, &pl.CreatedBy, &pl.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pl, nil
}

// UpdatePriceList updates a price list
func (r *PriceListRepository) UpdatePriceList(ctx context.Context, pl *products.PriceList) error {
	if err := r.db.SetOrganizationContext(ctx, pl.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE price_lists SET
			price_list_code = $3, price_list_name = $4, price_list_type = $5,
			effective_from = $6, effective_to = $7, base_price_adjustment_type = $8,
			base_price_adjustment_value = $9, priority = $10, is_active = $11,
			description = $12, updated_at = $13, updated_by = $14
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		pl.OrganizationID, pl.ID,
		pl.Code, pl.Name, pl.Type,
		pl.EffectiveFrom, pl.EffectiveTo, pl.BaseAdjustmentType,
		pl.BaseAdjustmentValue, pl.Priority, pl.IsActive,
		pl.Description, pl.UpdatedAt, pl.UpdatedBy,
	)

	return err
}

// DeletePriceList soft-deletes a price list
func (r *PriceListRepository) DeletePriceList(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE price_lists
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// ListPriceLists retrieves price lists
func (r *PriceListRepository) ListPriceLists(ctx context.Context, orgID uuid.UUID, filters products.PriceListFilters) ([]products.PriceList, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, price_list_code, price_list_name, price_list_type,
			effective_from, effective_to, base_price_adjustment_type,
			base_price_adjustment_value, priority, is_active, description,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM price_lists
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND price_list_type = $%d", argCount)
		args = append(args, filters.Type)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (price_list_code ILIKE $%d OR price_list_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	query += " ORDER BY priority ASC, price_list_name ASC"

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

	var priceLists []products.PriceList
	for rows.Next() {
		var pl products.PriceList
		err := rows.Scan(
			&pl.ID, &pl.OrganizationID, &pl.Code, &pl.Name, &pl.Type,
			&pl.EffectiveFrom, &pl.EffectiveTo, &pl.BaseAdjustmentType,
			&pl.BaseAdjustmentValue, &pl.Priority, &pl.IsActive, &pl.Description,
			&pl.CreatedAt, &pl.UpdatedAt, &pl.DeletedAt, &pl.CreatedBy, &pl.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		priceLists = append(priceLists, pl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return priceLists, nil
}

// CountPriceLists counts price lists
func (r *PriceListRepository) CountPriceLists(ctx context.Context, orgID uuid.UUID, filters products.PriceListFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM price_lists WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND price_list_type = $%d", argCount)
		args = append(args, filters.Type)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (price_list_code ILIKE $%d OR price_list_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// CreatePriceListItem creates a price list item
func (r *PriceListRepository) CreatePriceListItem(ctx context.Context, item *products.PriceListItem) error {
	query := `
		INSERT INTO price_list_items (
			id, price_list_id, product_id, product_variant_id, category_id,
			override_price, discount_percentage, markup_percentage,
			min_price, max_price, min_quantity,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.PriceListID, item.ProductID, item.ProductVariantID, item.CategoryID,
		item.OverridePrice, item.DiscountPercentage, item.MarkupPercentage,
		item.MinPrice, item.MaxPrice, item.MinQuantity,
		item.CreatedAt, item.UpdatedAt,
	)

	return err
}

// GetPriceListItem retrieves a price list item
func (r *PriceListRepository) GetPriceListItem(ctx context.Context, id uuid.UUID) (*products.PriceListItem, error) {
	query := `
		SELECT
			id, price_list_id, product_id, product_variant_id, category_id,
			override_price, discount_percentage, markup_percentage,
			min_price, max_price, min_quantity,
			created_at, updated_at, deleted_at
		FROM price_list_items
		WHERE id = $1 AND deleted_at IS NULL
	`

	var item products.PriceListItem
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.PriceListID, &item.ProductID, &item.ProductVariantID, &item.CategoryID,
		&item.OverridePrice, &item.DiscountPercentage, &item.MarkupPercentage,
		&item.MinPrice, &item.MaxPrice, &item.MinQuantity,
		&item.CreatedAt, &item.UpdatedAt, &item.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdatePriceListItem updates a price list item
func (r *PriceListRepository) UpdatePriceListItem(ctx context.Context, item *products.PriceListItem) error {
	query := `
		UPDATE price_list_items SET
			product_id = $2, product_variant_id = $3, category_id = $4,
			override_price = $5, discount_percentage = $6, markup_percentage = $7,
			min_price = $8, max_price = $9, min_quantity = $10,
			updated_at = $11
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID,
		item.ProductID, item.ProductVariantID, item.CategoryID,
		item.OverridePrice, item.DiscountPercentage, item.MarkupPercentage,
		item.MinPrice, item.MaxPrice, item.MinQuantity,
		item.UpdatedAt,
	)

	return err
}

// DeletePriceListItem soft-deletes a price list item
func (r *PriceListRepository) DeletePriceListItem(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE price_list_items
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	return err
}

// ListPriceListItems retrieves price list items
func (r *PriceListRepository) ListPriceListItems(ctx context.Context, filters products.PriceListItemFilters) ([]products.PriceListItem, error) {
	query := `
		SELECT
			id, price_list_id, product_id, product_variant_id, category_id,
			override_price, discount_percentage, markup_percentage,
			min_price, max_price, min_quantity,
			created_at, updated_at, deleted_at
		FROM price_list_items
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 0

	if filters.PriceListID != nil {
		argCount++
		query += fmt.Sprintf(" AND price_list_id = $%d", argCount)
		args = append(args, *filters.PriceListID)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.VariantID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_variant_id = $%d", argCount)
		args = append(args, *filters.VariantID)
	}

	if filters.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
	}

	query += " ORDER BY created_at DESC"

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

	var items []products.PriceListItem
	for rows.Next() {
		var item products.PriceListItem
		err := rows.Scan(
			&item.ID, &item.PriceListID, &item.ProductID, &item.ProductVariantID, &item.CategoryID,
			&item.OverridePrice, &item.DiscountPercentage, &item.MarkupPercentage,
			&item.MinPrice, &item.MaxPrice, &item.MinQuantity,
			&item.CreatedAt, &item.UpdatedAt, &item.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// CountPriceListItems counts price list items
func (r *PriceListRepository) CountPriceListItems(ctx context.Context, filters products.PriceListItemFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM price_list_items WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 0

	if filters.PriceListID != nil {
		argCount++
		query += fmt.Sprintf(" AND price_list_id = $%d", argCount)
		args = append(args, *filters.PriceListID)
	}

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.VariantID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_variant_id = $%d", argCount)
		args = append(args, *filters.VariantID)
	}

	if filters.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// CreateComponent creates a product component
func (r *PriceListRepository) CreateComponent(ctx context.Context, component *products.ProductComponent) error {
	if err := r.db.SetOrganizationContext(ctx, component.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO product_components (
			id, organization_id, parent_product_id,
			component_product_id, component_variant_id, quantity,
			inherit_price, price_override, display_order, is_optional,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		component.ID, component.OrganizationID, component.ParentProductID,
		component.ComponentProductID, component.ComponentVariantID, component.Quantity,
		component.InheritPrice, component.PriceOverride, component.DisplayOrder, component.IsOptional,
		component.CreatedAt, component.UpdatedAt,
	)

	return err
}

// GetComponent retrieves a component
func (r *PriceListRepository) GetComponent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*products.ProductComponent, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, parent_product_id,
			component_product_id, component_variant_id, quantity,
			inherit_price, price_override, display_order, is_optional,
			created_at, updated_at, deleted_at
		FROM product_components
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var component products.ProductComponent
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&component.ID, &component.OrganizationID, &component.ParentProductID,
		&component.ComponentProductID, &component.ComponentVariantID, &component.Quantity,
		&component.InheritPrice, &component.PriceOverride, &component.DisplayOrder, &component.IsOptional,
		&component.CreatedAt, &component.UpdatedAt, &component.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &component, nil
}

// UpdateComponent updates a component
func (r *PriceListRepository) UpdateComponent(ctx context.Context, component *products.ProductComponent) error {
	if err := r.db.SetOrganizationContext(ctx, component.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_components SET
			component_product_id = $3, component_variant_id = $4, quantity = $5,
			inherit_price = $6, price_override = $7, display_order = $8, is_optional = $9,
			updated_at = $10
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		component.OrganizationID, component.ID,
		component.ComponentProductID, component.ComponentVariantID, component.Quantity,
		component.InheritPrice, component.PriceOverride, component.DisplayOrder, component.IsOptional,
		component.UpdatedAt,
	)

	return err
}

// DeleteComponent soft-deletes a component
func (r *PriceListRepository) DeleteComponent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_components
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// ListComponents retrieves components with filters
func (r *PriceListRepository) ListComponents(ctx context.Context, orgID uuid.UUID, filters products.ComponentFilters) ([]products.ProductComponent, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, parent_product_id,
			component_product_id, component_variant_id, quantity,
			inherit_price, price_override, display_order, is_optional,
			created_at, updated_at, deleted_at
		FROM product_components
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.ParentProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND parent_product_id = $%d", argCount)
		args = append(args, *filters.ParentProductID)
	}

	query += " ORDER BY display_order ASC"

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

	var components []products.ProductComponent
	for rows.Next() {
		var component products.ProductComponent
		err := rows.Scan(
			&component.ID, &component.OrganizationID, &component.ParentProductID,
			&component.ComponentProductID, &component.ComponentVariantID, &component.Quantity,
			&component.InheritPrice, &component.PriceOverride, &component.DisplayOrder, &component.IsOptional,
			&component.CreatedAt, &component.UpdatedAt, &component.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		components = append(components, component)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return components, nil
}

// ListComponentsByParent retrieves all components of a product
func (r *PriceListRepository) ListComponentsByParent(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) ([]products.ProductComponent, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	filters := products.ComponentFilters{ParentProductID: &parentID}
	return r.ListComponents(ctx, orgID, filters)
}

// DeleteComponentsByParent deletes all components of a product
func (r *PriceListRepository) DeleteComponentsByParent(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE product_components
		SET deleted_at = $3
		WHERE organization_id = $1 AND parent_product_id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, parentID, time.Now())
	return err
}

// GetApplicablePriceLists retrieves applicable price lists at a given time
func (r *PriceListRepository) GetApplicablePriceLists(ctx context.Context, orgID uuid.UUID, now time.Time) ([]products.PriceList, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, price_list_code, price_list_name, price_list_type,
			effective_from, effective_to, base_price_adjustment_type,
			base_price_adjustment_value, priority, is_active, description,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM price_lists
		WHERE organization_id = $1 AND is_active = true AND deleted_at IS NULL
		AND (effective_from IS NULL OR effective_from <= $2)
		AND (effective_to IS NULL OR effective_to >= $2)
		ORDER BY priority ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var priceLists []products.PriceList
	for rows.Next() {
		var pl products.PriceList
		err := rows.Scan(
			&pl.ID, &pl.OrganizationID, &pl.Code, &pl.Name, &pl.Type,
			&pl.EffectiveFrom, &pl.EffectiveTo, &pl.BaseAdjustmentType,
			&pl.BaseAdjustmentValue, &pl.Priority, &pl.IsActive, &pl.Description,
			&pl.CreatedAt, &pl.UpdatedAt, &pl.DeletedAt, &pl.CreatedBy, &pl.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		priceLists = append(priceLists, pl)
	}

	return priceLists, rows.Err()
}

// CalculatePrice calculates the final price for a product
func (r *PriceListRepository) CalculatePrice(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, basePrice float64) (float64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	// Get applicable price lists
	applicableLists, err := r.GetApplicablePriceLists(ctx, orgID, time.Now())
	if err != nil {
		return 0, err
	}

	finalPrice := basePrice

	// Apply price lists in priority order
	for _, pl := range applicableLists {
		// Look for matching item in this price list
		var target *uuid.UUID
		var isVariantQuery bool

		if variantID != nil {
			target = variantID
			isVariantQuery = true
		} else {
			target = &productID
			isVariantQuery = false
		}

		query := `
			SELECT override_price, discount_percentage, markup_percentage
			FROM price_list_items
			WHERE price_list_id = $1 AND deleted_at IS NULL
		`

		if isVariantQuery {
			query += " AND product_variant_id = $2"
		} else {
			query += " AND product_id = $2"
		}

		var overridePrice *float64
		var discountPct *float64
		var markupPct *float64

		err := r.db.Pool.QueryRow(ctx, query, pl.ID, target).Scan(&overridePrice, &discountPct, &markupPct)

		if err == pgx.ErrNoRows {
			// Check for category-level override if no specific product override
			continue
		} else if err != nil {
			continue
		}

		// Apply the most specific override
		if overridePrice != nil {
			finalPrice = *overridePrice
		} else if discountPct != nil {
			finalPrice = finalPrice * (1 - *discountPct/100)
		} else if markupPct != nil {
			finalPrice = finalPrice * (1 + *markupPct/100)
		}
	}

	return finalPrice, nil
}

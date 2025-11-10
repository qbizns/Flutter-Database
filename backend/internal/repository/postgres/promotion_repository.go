package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/promotions"
)

type PromotionRepository struct {
	db *DB
}

func NewPromotionRepository(db *DB) *PromotionRepository {
	return &PromotionRepository{db: db}
}

// Promotions

func (r *PromotionRepository) List(ctx context.Context, orgID uuid.UUID, filters promotions.PromotionFilters) ([]promotions.Promotion, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, promotion_code, name, description, promotion_type,
		       discount_value, applies_to, applicable_product_ids, applicable_category_ids,
		       minimum_purchase_amount, minimum_quantity, buy_quantity, get_quantity,
		       get_discount_percentage, max_uses_total, max_uses_per_customer, current_uses,
		       start_date, end_date, is_active, is_combinable, priority,
		       terms_and_conditions, created_at, updated_at, created_by, updated_by
		FROM promotions
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR promotion_code ILIKE $%d OR description ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.PromotionType != nil {
		argCount++
		query += fmt.Sprintf(" AND promotion_type = $%d", argCount)
		args = append(args, *filters.PromotionType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY priority ASC, start_date DESC"

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

	var promotionList []promotions.Promotion
	for rows.Next() {
		var p promotions.Promotion
		var productIDs, categoryIDs []byte

		err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.PromotionCode, &p.Name, &p.Description, &p.PromotionType,
			&p.DiscountValue, &p.AppliesToType, &productIDs, &categoryIDs,
			&p.MinimumPurchaseAmount, &p.MinimumQuantity, &p.BuyQuantity, &p.GetQuantity,
			&p.GetDiscountPercentage, &p.MaxUsesTotal, &p.MaxUsesPerCustomer, &p.CurrentUses,
			&p.StartDate, &p.EndDate, &p.IsActive, &p.IsCombinable, &p.Priority,
			&p.TermsAndConditions, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		if productIDs != nil {
			json.Unmarshal(productIDs, &p.ApplicableProductIDs)
		}
		if categoryIDs != nil {
			json.Unmarshal(categoryIDs, &p.ApplicableCategoryIDs)
		}

		promotionList = append(promotionList, p)
	}

	return promotionList, rows.Err()
}

func (r *PromotionRepository) Count(ctx context.Context, orgID uuid.UUID, filters promotions.PromotionFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM promotions WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR promotion_code ILIKE $%d OR description ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.PromotionType != nil {
		argCount++
		query += fmt.Sprintf(" AND promotion_type = $%d", argCount)
		args = append(args, *filters.PromotionType)
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

func (r *PromotionRepository) Create(ctx context.Context, promotion *promotions.Promotion) error {
	if err := r.db.SetOrganizationContext(ctx, promotion.OrganizationID.String()); err != nil {
		return err
	}

	productIDs, _ := json.Marshal(promotion.ApplicableProductIDs)
	categoryIDs, _ := json.Marshal(promotion.ApplicableCategoryIDs)

	query := `
		INSERT INTO promotions (
			id, organization_id, promotion_code, name, description, promotion_type,
			discount_value, applies_to, applicable_product_ids, applicable_category_ids,
			minimum_purchase_amount, minimum_quantity, buy_quantity, get_quantity,
			get_discount_percentage, max_uses_total, max_uses_per_customer, current_uses,
			start_date, end_date, is_active, is_combinable, priority,
			terms_and_conditions, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26, $27
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		promotion.ID, promotion.OrganizationID, promotion.PromotionCode, promotion.Name,
		promotion.Description, promotion.PromotionType, promotion.DiscountValue, promotion.AppliesToType,
		productIDs, categoryIDs, promotion.MinimumPurchaseAmount, promotion.MinimumQuantity,
		promotion.BuyQuantity, promotion.GetQuantity, promotion.GetDiscountPercentage,
		promotion.MaxUsesTotal, promotion.MaxUsesPerCustomer, promotion.CurrentUses,
		promotion.StartDate, promotion.EndDate, promotion.IsActive, promotion.IsCombinable,
		promotion.Priority, promotion.TermsAndConditions, promotion.CreatedAt, promotion.UpdatedAt,
		promotion.CreatedBy,
	)
	return err
}

func (r *PromotionRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*promotions.Promotion, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, promotion_code, name, description, promotion_type,
		       discount_value, applies_to, applicable_product_ids, applicable_category_ids,
		       minimum_purchase_amount, minimum_quantity, buy_quantity, get_quantity,
		       get_discount_percentage, max_uses_total, max_uses_per_customer, current_uses,
		       start_date, end_date, is_active, is_combinable, priority,
		       terms_and_conditions, created_at, updated_at, created_by, updated_by
		FROM promotions
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var p promotions.Promotion
	var productIDs, categoryIDs []byte

	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&p.ID, &p.OrganizationID, &p.PromotionCode, &p.Name, &p.Description, &p.PromotionType,
		&p.DiscountValue, &p.AppliesToType, &productIDs, &categoryIDs,
		&p.MinimumPurchaseAmount, &p.MinimumQuantity, &p.BuyQuantity, &p.GetQuantity,
		&p.GetDiscountPercentage, &p.MaxUsesTotal, &p.MaxUsesPerCustomer, &p.CurrentUses,
		&p.StartDate, &p.EndDate, &p.IsActive, &p.IsCombinable, &p.Priority,
		&p.TermsAndConditions, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if productIDs != nil {
		json.Unmarshal(productIDs, &p.ApplicableProductIDs)
	}
	if categoryIDs != nil {
		json.Unmarshal(categoryIDs, &p.ApplicableCategoryIDs)
	}

	return &p, nil
}

func (r *PromotionRepository) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*promotions.Promotion, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, promotion_code, name, description, promotion_type,
		       discount_value, applies_to, applicable_product_ids, applicable_category_ids,
		       minimum_purchase_amount, minimum_quantity, buy_quantity, get_quantity,
		       get_discount_percentage, max_uses_total, max_uses_per_customer, current_uses,
		       start_date, end_date, is_active, is_combinable, priority,
		       terms_and_conditions, created_at, updated_at, created_by, updated_by
		FROM promotions
		WHERE organization_id = $1 AND LOWER(promotion_code) = LOWER($2) AND deleted_at IS NULL
	`

	var p promotions.Promotion
	var productIDs, categoryIDs []byte

	err := r.db.Pool.QueryRow(ctx, query, orgID, code).Scan(
		&p.ID, &p.OrganizationID, &p.PromotionCode, &p.Name, &p.Description, &p.PromotionType,
		&p.DiscountValue, &p.AppliesToType, &productIDs, &categoryIDs,
		&p.MinimumPurchaseAmount, &p.MinimumQuantity, &p.BuyQuantity, &p.GetQuantity,
		&p.GetDiscountPercentage, &p.MaxUsesTotal, &p.MaxUsesPerCustomer, &p.CurrentUses,
		&p.StartDate, &p.EndDate, &p.IsActive, &p.IsCombinable, &p.Priority,
		&p.TermsAndConditions, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if productIDs != nil {
		json.Unmarshal(productIDs, &p.ApplicableProductIDs)
	}
	if categoryIDs != nil {
		json.Unmarshal(categoryIDs, &p.ApplicableCategoryIDs)
	}

	return &p, nil
}

func (r *PromotionRepository) Update(ctx context.Context, promotion *promotions.Promotion) error {
	if err := r.db.SetOrganizationContext(ctx, promotion.OrganizationID.String()); err != nil {
		return err
	}

	productIDs, _ := json.Marshal(promotion.ApplicableProductIDs)
	categoryIDs, _ := json.Marshal(promotion.ApplicableCategoryIDs)

	query := `
		UPDATE promotions SET
			promotion_code = $3, name = $4, description = $5, promotion_type = $6,
			discount_value = $7, applies_to = $8, applicable_product_ids = $9,
			applicable_category_ids = $10, minimum_purchase_amount = $11, minimum_quantity = $12,
			buy_quantity = $13, get_quantity = $14, get_discount_percentage = $15,
			max_uses_total = $16, max_uses_per_customer = $17, start_date = $18,
			end_date = $19, is_active = $20, is_combinable = $21, priority = $22,
			terms_and_conditions = $23, updated_at = $24, updated_by = $25
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		promotion.OrganizationID, promotion.ID, promotion.PromotionCode, promotion.Name,
		promotion.Description, promotion.PromotionType, promotion.DiscountValue, promotion.AppliesToType,
		productIDs, categoryIDs, promotion.MinimumPurchaseAmount, promotion.MinimumQuantity,
		promotion.BuyQuantity, promotion.GetQuantity, promotion.GetDiscountPercentage,
		promotion.MaxUsesTotal, promotion.MaxUsesPerCustomer, promotion.StartDate,
		promotion.EndDate, promotion.IsActive, promotion.IsCombinable, promotion.Priority,
		promotion.TermsAndConditions, promotion.UpdatedAt, promotion.UpdatedBy,
	)
	return err
}

func (r *PromotionRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE promotions
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// Promotion Usage

func (r *PromotionRepository) ListUsage(ctx context.Context, orgID uuid.UUID, filters promotions.PromotionUsageFilters) ([]promotions.PromotionUsage, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, promotion_id, sale_id, customer_id,
		       discount_amount, used_at, created_at
		FROM promotion_usage
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.PromotionID != nil {
		argCount++
		query += fmt.Sprintf(" AND promotion_id = $%d", argCount)
		args = append(args, *filters.PromotionID)
	}

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND used_at >= $%d", argCount)
		args = append(args, *filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND used_at <= $%d", argCount)
		args = append(args, *filters.DateTo)
	}

	query += " ORDER BY used_at DESC"

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

	var usageList []promotions.PromotionUsage
	for rows.Next() {
		var u promotions.PromotionUsage
		err := rows.Scan(
			&u.ID, &u.OrganizationID, &u.PromotionID, &u.SaleID, &u.CustomerID,
			&u.DiscountAmount, &u.UsedAt, &u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		usageList = append(usageList, u)
	}

	return usageList, rows.Err()
}

func (r *PromotionRepository) CountUsage(ctx context.Context, orgID uuid.UUID, filters promotions.PromotionUsageFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM promotion_usage WHERE organization_id = $1"
	args := []interface{}{orgID}
	argCount := 1

	if filters.PromotionID != nil {
		argCount++
		query += fmt.Sprintf(" AND promotion_id = $%d", argCount)
		args = append(args, *filters.PromotionID)
	}

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND used_at >= $%d", argCount)
		args = append(args, *filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND used_at <= $%d", argCount)
		args = append(args, *filters.DateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *PromotionRepository) RecordUsage(ctx context.Context, usage *promotions.PromotionUsage) error {
	if err := r.db.SetOrganizationContext(ctx, usage.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO promotion_usage (
			id, organization_id, promotion_id, sale_id, customer_id,
			discount_amount, used_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		usage.ID, usage.OrganizationID, usage.PromotionID, usage.SaleID, usage.CustomerID,
		usage.DiscountAmount, usage.UsedAt, usage.CreatedAt,
	)
	return err
}

func (r *PromotionRepository) GetUsageByPromotion(ctx context.Context, orgID uuid.UUID, promotionID uuid.UUID, customerID *uuid.UUID) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM promotion_usage WHERE organization_id = $1 AND promotion_id = $2"
	args := []interface{}{orgID, promotionID}

	if customerID != nil {
		query += " AND customer_id = $3"
		args = append(args, *customerID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

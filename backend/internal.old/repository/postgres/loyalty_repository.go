package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/loyalty"
)

// LoyaltyRepository implements loyalty.LoyaltyTierRepository
type LoyaltyRepository struct {
	db *DB
}

// NewLoyaltyRepository creates a new loyalty repository
func NewLoyaltyRepository(db *DB) *LoyaltyRepository {
	return &LoyaltyRepository{db: db}
}

// ========================================================================
// LOYALTY TIERS
// ========================================================================

// ListTiers retrieves loyalty tiers with filters
func (r *LoyaltyRepository) ListTiers(ctx context.Context, orgID uuid.UUID, filters loyalty.LoyaltyTierFilters) ([]loyalty.LoyaltyTier, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, tier_code, tier_name, tier_level, description,
			points_threshold, annual_spend_threshold, purchase_count_threshold,
			points_multiplier, discount_percentage, tier_color, tier_icon,
			badge_image_url, is_active, is_default, sort_order, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_tiers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (tier_name ILIKE $%d OR tier_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	// Add ordering
	query += " ORDER BY tier_level ASC, sort_order ASC"

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

	var tiers []loyalty.LoyaltyTier
	for rows.Next() {
		var t loyalty.LoyaltyTier
		err := rows.Scan(
			&t.ID, &t.OrganizationID, &t.TierCode, &t.TierName, &t.TierLevel, &t.Description,
			&t.PointsThreshold, &t.AnnualSpendThreshold, &t.PurchaseCountThreshold,
			&t.PointsMultiplier, &t.DiscountPercentage, &t.TierColor, &t.TierIcon,
			&t.BadgeImageURL, &t.IsActive, &t.IsDefault, &t.SortOrder, &t.Metadata,
			&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.CreatedBy, &t.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		tiers = append(tiers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tiers, nil
}

// CountTiers counts loyalty tiers matching filters
func (r *LoyaltyRepository) CountTiers(ctx context.Context, orgID uuid.UUID, filters loyalty.LoyaltyTierFilters) (int64, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM loyalty_tiers WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (tier_name ILIKE $%d OR tier_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CreateTier creates a new loyalty tier
func (r *LoyaltyRepository) CreateTier(ctx context.Context, tier *loyalty.LoyaltyTier) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, tier.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO loyalty_tiers (
			id, organization_id, tier_code, tier_name, tier_level, description,
			points_threshold, annual_spend_threshold, purchase_count_threshold,
			points_multiplier, discount_percentage, tier_color, tier_icon,
			badge_image_url, is_active, is_default, sort_order, metadata,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		tier.ID, tier.OrganizationID, tier.TierCode, tier.TierName, tier.TierLevel, tier.Description,
		tier.PointsThreshold, tier.AnnualSpendThreshold, tier.PurchaseCountThreshold,
		tier.PointsMultiplier, tier.DiscountPercentage, tier.TierColor, tier.TierIcon,
		tier.BadgeImageURL, tier.IsActive, tier.IsDefault, tier.SortOrder, tier.Metadata,
		tier.CreatedAt, tier.UpdatedAt, tier.CreatedBy,
	)

	return err
}

// GetTier retrieves a loyalty tier by ID
func (r *LoyaltyRepository) GetTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*loyalty.LoyaltyTier, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, tier_code, tier_name, tier_level, description,
			points_threshold, annual_spend_threshold, purchase_count_threshold,
			points_multiplier, discount_percentage, tier_color, tier_icon,
			badge_image_url, is_active, is_default, sort_order, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_tiers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var t loyalty.LoyaltyTier
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&t.ID, &t.OrganizationID, &t.TierCode, &t.TierName, &t.TierLevel, &t.Description,
		&t.PointsThreshold, &t.AnnualSpendThreshold, &t.PurchaseCountThreshold,
		&t.PointsMultiplier, &t.DiscountPercentage, &t.TierColor, &t.TierIcon,
		&t.BadgeImageURL, &t.IsActive, &t.IsDefault, &t.SortOrder, &t.Metadata,
		&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.CreatedBy, &t.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetTierByCode retrieves a loyalty tier by code
func (r *LoyaltyRepository) GetTierByCode(ctx context.Context, orgID uuid.UUID, code string) (*loyalty.LoyaltyTier, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, tier_code, tier_name, tier_level, description,
			points_threshold, annual_spend_threshold, purchase_count_threshold,
			points_multiplier, discount_percentage, tier_color, tier_icon,
			badge_image_url, is_active, is_default, sort_order, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_tiers
		WHERE organization_id = $1 AND tier_code = $2 AND deleted_at IS NULL
	`

	var t loyalty.LoyaltyTier
	err := r.db.Pool.QueryRow(ctx, query, orgID, code).Scan(
		&t.ID, &t.OrganizationID, &t.TierCode, &t.TierName, &t.TierLevel, &t.Description,
		&t.PointsThreshold, &t.AnnualSpendThreshold, &t.PurchaseCountThreshold,
		&t.PointsMultiplier, &t.DiscountPercentage, &t.TierColor, &t.TierIcon,
		&t.BadgeImageURL, &t.IsActive, &t.IsDefault, &t.SortOrder, &t.Metadata,
		&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.CreatedBy, &t.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// UpdateTier updates an existing loyalty tier
func (r *LoyaltyRepository) UpdateTier(ctx context.Context, tier *loyalty.LoyaltyTier) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, tier.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE loyalty_tiers SET
			tier_code = $3, tier_name = $4, tier_level = $5, description = $6,
			points_threshold = $7, annual_spend_threshold = $8, purchase_count_threshold = $9,
			points_multiplier = $10, discount_percentage = $11, tier_color = $12, tier_icon = $13,
			badge_image_url = $14, is_active = $15, is_default = $16, sort_order = $17, metadata = $18,
			updated_at = $19, updated_by = $20
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		tier.OrganizationID, tier.ID,
		tier.TierCode, tier.TierName, tier.TierLevel, tier.Description,
		tier.PointsThreshold, tier.AnnualSpendThreshold, tier.PurchaseCountThreshold,
		tier.PointsMultiplier, tier.DiscountPercentage, tier.TierColor, tier.TierIcon,
		tier.BadgeImageURL, tier.IsActive, tier.IsDefault, tier.SortOrder, tier.Metadata,
		tier.UpdatedAt, tier.UpdatedBy,
	)

	return err
}

// DeleteTier soft-deletes a loyalty tier
func (r *LoyaltyRepository) DeleteTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE loyalty_tiers
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// GetDefaultTier retrieves the default tier for an organization
func (r *LoyaltyRepository) GetDefaultTier(ctx context.Context, orgID uuid.UUID) (*loyalty.LoyaltyTier, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, tier_code, tier_name, tier_level, description,
			points_threshold, annual_spend_threshold, purchase_count_threshold,
			points_multiplier, discount_percentage, tier_color, tier_icon,
			badge_image_url, is_active, is_default, sort_order, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_tiers
		WHERE organization_id = $1 AND is_default = true AND deleted_at IS NULL
		LIMIT 1
	`

	var t loyalty.LoyaltyTier
	err := r.db.Pool.QueryRow(ctx, query, orgID).Scan(
		&t.ID, &t.OrganizationID, &t.TierCode, &t.TierName, &t.TierLevel, &t.Description,
		&t.PointsThreshold, &t.AnnualSpendThreshold, &t.PurchaseCountThreshold,
		&t.PointsMultiplier, &t.DiscountPercentage, &t.TierColor, &t.TierIcon,
		&t.BadgeImageURL, &t.IsActive, &t.IsDefault, &t.SortOrder, &t.Metadata,
		&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt, &t.CreatedBy, &t.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// ========================================================================
// LOYALTY TIER BENEFITS
// ========================================================================

// ListBenefits retrieves loyalty tier benefits with filters
func (r *LoyaltyRepository) ListBenefits(ctx context.Context, orgID uuid.UUID, filters loyalty.LoyaltyTierBenefitFilters) ([]loyalty.LoyaltyTierBenefit, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, tier_id, benefit_code, benefit_name, benefit_description,
			benefit_type, discount_value, discount_type, is_active, sort_order, icon,
			terms_and_conditions, metadata, created_at, updated_at
		FROM loyalty_tier_benefits
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.TierID != nil {
		argCount++
		query += fmt.Sprintf(" AND tier_id = $%d", argCount)
		args = append(args, *filters.TierID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY sort_order ASC"

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

	var benefits []loyalty.LoyaltyTierBenefit
	for rows.Next() {
		var b loyalty.LoyaltyTierBenefit
		err := rows.Scan(
			&b.ID, &b.OrganizationID, &b.TierID, &b.BenefitCode, &b.BenefitName, &b.BenefitDescription,
			&b.BenefitType, &b.DiscountValue, &b.DiscountType, &b.IsActive, &b.SortOrder, &b.Icon,
			&b.TermsAndConditions, &b.Metadata, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		benefits = append(benefits, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return benefits, nil
}

// CreateBenefit creates a new loyalty tier benefit
func (r *LoyaltyRepository) CreateBenefit(ctx context.Context, benefit *loyalty.LoyaltyTierBenefit) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, benefit.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO loyalty_tier_benefits (
			id, organization_id, tier_id, benefit_code, benefit_name, benefit_description,
			benefit_type, discount_value, discount_type, is_active, sort_order, icon,
			terms_and_conditions, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		benefit.ID, benefit.OrganizationID, benefit.TierID, benefit.BenefitCode, benefit.BenefitName,
		benefit.BenefitDescription, benefit.BenefitType, benefit.DiscountValue, benefit.DiscountType,
		benefit.IsActive, benefit.SortOrder, benefit.Icon, benefit.TermsAndConditions, benefit.Metadata,
		benefit.CreatedAt, benefit.UpdatedAt,
	)

	return err
}

// GetBenefit retrieves a loyalty tier benefit by ID
func (r *LoyaltyRepository) GetBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*loyalty.LoyaltyTierBenefit, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, tier_id, benefit_code, benefit_name, benefit_description,
			benefit_type, discount_value, discount_type, is_active, sort_order, icon,
			terms_and_conditions, metadata, created_at, updated_at
		FROM loyalty_tier_benefits
		WHERE organization_id = $1 AND id = $2
	`

	var b loyalty.LoyaltyTierBenefit
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&b.ID, &b.OrganizationID, &b.TierID, &b.BenefitCode, &b.BenefitName, &b.BenefitDescription,
		&b.BenefitType, &b.DiscountValue, &b.DiscountType, &b.IsActive, &b.SortOrder, &b.Icon,
		&b.TermsAndConditions, &b.Metadata, &b.CreatedAt, &b.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// UpdateBenefit updates an existing loyalty tier benefit
func (r *LoyaltyRepository) UpdateBenefit(ctx context.Context, benefit *loyalty.LoyaltyTierBenefit) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, benefit.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE loyalty_tier_benefits SET
			benefit_code = $3, benefit_name = $4, benefit_description = $5,
			benefit_type = $6, discount_value = $7, discount_type = $8, is_active = $9,
			sort_order = $10, icon = $11, terms_and_conditions = $12, metadata = $13,
			updated_at = $14
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		benefit.OrganizationID, benefit.ID,
		benefit.BenefitCode, benefit.BenefitName, benefit.BenefitDescription,
		benefit.BenefitType, benefit.DiscountValue, benefit.DiscountType, benefit.IsActive,
		benefit.SortOrder, benefit.Icon, benefit.TermsAndConditions, benefit.Metadata,
		benefit.UpdatedAt,
	)

	return err
}

// DeleteBenefit deletes a loyalty tier benefit
func (r *LoyaltyRepository) DeleteBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "DELETE FROM loyalty_tier_benefits WHERE organization_id = $1 AND id = $2"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// GetTierBenefits retrieves all benefits for a specific tier
func (r *LoyaltyRepository) GetTierBenefits(ctx context.Context, orgID uuid.UUID, tierID uuid.UUID) ([]loyalty.LoyaltyTierBenefit, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, tier_id, benefit_code, benefit_name, benefit_description,
			benefit_type, discount_value, discount_type, is_active, sort_order, icon,
			terms_and_conditions, metadata, created_at, updated_at
		FROM loyalty_tier_benefits
		WHERE organization_id = $1 AND tier_id = $2
		ORDER BY sort_order ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, tierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var benefits []loyalty.LoyaltyTierBenefit
	for rows.Next() {
		var b loyalty.LoyaltyTierBenefit
		err := rows.Scan(
			&b.ID, &b.OrganizationID, &b.TierID, &b.BenefitCode, &b.BenefitName, &b.BenefitDescription,
			&b.BenefitType, &b.DiscountValue, &b.DiscountType, &b.IsActive, &b.SortOrder, &b.Icon,
			&b.TermsAndConditions, &b.Metadata, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		benefits = append(benefits, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return benefits, nil
}

// ========================================================================
// LOYALTY POINTS RULES
// ========================================================================

// ListPointsRules retrieves loyalty points rules with filters
func (r *LoyaltyRepository) ListPointsRules(ctx context.Context, orgID uuid.UUID, filters loyalty.LoyaltyPointsRuleFilters) ([]loyalty.LoyaltyPointsRule, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, rule_code, rule_name, description, rule_type,
			points_per_amount, fixed_points, multiplier, applies_to,
			applicable_product_ids, applicable_category_ids, applicable_tier_ids,
			minimum_purchase_amount, maximum_points_per_transaction,
			maximum_points_per_day, maximum_points_per_month,
			start_date, end_date, is_active, priority,
			terms_and_conditions, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_points_rules
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (rule_name ILIKE $%d OR rule_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.RuleType != nil {
		argCount++
		query += fmt.Sprintf(" AND rule_type = $%d", argCount)
		args = append(args, *filters.RuleType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY priority DESC, created_at DESC"

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

	var rules []loyalty.LoyaltyPointsRule
	for rows.Next() {
		var rule loyalty.LoyaltyPointsRule
		err := rows.Scan(
			&rule.ID, &rule.OrganizationID, &rule.RuleCode, &rule.RuleName, &rule.Description, &rule.RuleType,
			&rule.PointsPerAmount, &rule.FixedPoints, &rule.Multiplier, &rule.AppliesToType,
			&rule.ApplicableProductIDs, &rule.ApplicableCategoryIDs, &rule.ApplicableTierIDs,
			&rule.MinimumPurchaseAmount, &rule.MaximumPointsPerTransaction,
			&rule.MaximumPointsPerDay, &rule.MaximumPointsPerMonth,
			&rule.StartDate, &rule.EndDate, &rule.IsActive, &rule.Priority,
			&rule.TermsAndConditions, &rule.Metadata,
			&rule.CreatedAt, &rule.UpdatedAt, &rule.DeletedAt, &rule.CreatedBy, &rule.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

// CreatePointsRule creates a new loyalty points rule
func (r *LoyaltyRepository) CreatePointsRule(ctx context.Context, rule *loyalty.LoyaltyPointsRule) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, rule.OrganizationID.String()); err != nil {
		return err
	}

	applicableProductIDsJSON, _ := json.Marshal(rule.ApplicableProductIDs)
	applicableCategoryIDsJSON, _ := json.Marshal(rule.ApplicableCategoryIDs)
	applicableTierIDsJSON, _ := json.Marshal(rule.ApplicableTierIDs)

	query := `
		INSERT INTO loyalty_points_rules (
			id, organization_id, rule_code, rule_name, description, rule_type,
			points_per_amount, fixed_points, multiplier, applies_to,
			applicable_product_ids, applicable_category_ids, applicable_tier_ids,
			minimum_purchase_amount, maximum_points_per_transaction,
			maximum_points_per_day, maximum_points_per_month,
			start_date, end_date, is_active, priority,
			terms_and_conditions, metadata,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		rule.ID, rule.OrganizationID, rule.RuleCode, rule.RuleName, rule.Description, rule.RuleType,
		rule.PointsPerAmount, rule.FixedPoints, rule.Multiplier, rule.AppliesToType,
		applicableProductIDsJSON, applicableCategoryIDsJSON, applicableTierIDsJSON,
		rule.MinimumPurchaseAmount, rule.MaximumPointsPerTransaction,
		rule.MaximumPointsPerDay, rule.MaximumPointsPerMonth,
		rule.StartDate, rule.EndDate, rule.IsActive, rule.Priority,
		rule.TermsAndConditions, rule.Metadata,
		rule.CreatedAt, rule.UpdatedAt, rule.CreatedBy,
	)

	return err
}

// GetPointsRule retrieves a loyalty points rule by ID
func (r *LoyaltyRepository) GetPointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*loyalty.LoyaltyPointsRule, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, rule_code, rule_name, description, rule_type,
			points_per_amount, fixed_points, multiplier, applies_to,
			applicable_product_ids, applicable_category_ids, applicable_tier_ids,
			minimum_purchase_amount, maximum_points_per_transaction,
			maximum_points_per_day, maximum_points_per_month,
			start_date, end_date, is_active, priority,
			terms_and_conditions, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_points_rules
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var rule loyalty.LoyaltyPointsRule
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&rule.ID, &rule.OrganizationID, &rule.RuleCode, &rule.RuleName, &rule.Description, &rule.RuleType,
		&rule.PointsPerAmount, &rule.FixedPoints, &rule.Multiplier, &rule.AppliesToType,
		&rule.ApplicableProductIDs, &rule.ApplicableCategoryIDs, &rule.ApplicableTierIDs,
		&rule.MinimumPurchaseAmount, &rule.MaximumPointsPerTransaction,
		&rule.MaximumPointsPerDay, &rule.MaximumPointsPerMonth,
		&rule.StartDate, &rule.EndDate, &rule.IsActive, &rule.Priority,
		&rule.TermsAndConditions, &rule.Metadata,
		&rule.CreatedAt, &rule.UpdatedAt, &rule.DeletedAt, &rule.CreatedBy, &rule.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// UpdatePointsRule updates an existing loyalty points rule
func (r *LoyaltyRepository) UpdatePointsRule(ctx context.Context, rule *loyalty.LoyaltyPointsRule) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, rule.OrganizationID.String()); err != nil {
		return err
	}

	applicableProductIDsJSON, _ := json.Marshal(rule.ApplicableProductIDs)
	applicableCategoryIDsJSON, _ := json.Marshal(rule.ApplicableCategoryIDs)
	applicableTierIDsJSON, _ := json.Marshal(rule.ApplicableTierIDs)

	query := `
		UPDATE loyalty_points_rules SET
			rule_code = $3, rule_name = $4, description = $5, rule_type = $6,
			points_per_amount = $7, fixed_points = $8, multiplier = $9, applies_to = $10,
			applicable_product_ids = $11, applicable_category_ids = $12, applicable_tier_ids = $13,
			minimum_purchase_amount = $14, maximum_points_per_transaction = $15,
			maximum_points_per_day = $16, maximum_points_per_month = $17,
			start_date = $18, end_date = $19, is_active = $20, priority = $21,
			terms_and_conditions = $22, metadata = $23,
			updated_at = $24, updated_by = $25
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		rule.OrganizationID, rule.ID,
		rule.RuleCode, rule.RuleName, rule.Description, rule.RuleType,
		rule.PointsPerAmount, rule.FixedPoints, rule.Multiplier, rule.AppliesToType,
		applicableProductIDsJSON, applicableCategoryIDsJSON, applicableTierIDsJSON,
		rule.MinimumPurchaseAmount, rule.MaximumPointsPerTransaction,
		rule.MaximumPointsPerDay, rule.MaximumPointsPerMonth,
		rule.StartDate, rule.EndDate, rule.IsActive, rule.Priority,
		rule.TermsAndConditions, rule.Metadata,
		rule.UpdatedAt, rule.UpdatedBy,
	)

	return err
}

// DeletePointsRule soft-deletes a loyalty points rule
func (r *LoyaltyRepository) DeletePointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE loyalty_points_rules
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// GetActivePointsRules retrieves all active points rules
func (r *LoyaltyRepository) GetActivePointsRules(ctx context.Context, orgID uuid.UUID) ([]loyalty.LoyaltyPointsRule, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, rule_code, rule_name, description, rule_type,
			points_per_amount, fixed_points, multiplier, applies_to,
			applicable_product_ids, applicable_category_ids, applicable_tier_ids,
			minimum_purchase_amount, maximum_points_per_transaction,
			maximum_points_per_day, maximum_points_per_month,
			start_date, end_date, is_active, priority,
			terms_and_conditions, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_points_rules
		WHERE organization_id = $1 AND is_active = true AND deleted_at IS NULL
		AND (start_date IS NULL OR start_date <= CURRENT_DATE)
		AND (end_date IS NULL OR end_date >= CURRENT_DATE)
		ORDER BY priority DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []loyalty.LoyaltyPointsRule
	for rows.Next() {
		var rule loyalty.LoyaltyPointsRule
		err := rows.Scan(
			&rule.ID, &rule.OrganizationID, &rule.RuleCode, &rule.RuleName, &rule.Description, &rule.RuleType,
			&rule.PointsPerAmount, &rule.FixedPoints, &rule.Multiplier, &rule.AppliesToType,
			&rule.ApplicableProductIDs, &rule.ApplicableCategoryIDs, &rule.ApplicableTierIDs,
			&rule.MinimumPurchaseAmount, &rule.MaximumPointsPerTransaction,
			&rule.MaximumPointsPerDay, &rule.MaximumPointsPerMonth,
			&rule.StartDate, &rule.EndDate, &rule.IsActive, &rule.Priority,
			&rule.TermsAndConditions, &rule.Metadata,
			&rule.CreatedAt, &rule.UpdatedAt, &rule.DeletedAt, &rule.CreatedBy, &rule.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

// ========================================================================
// LOYALTY REWARDS
// ========================================================================

// ListRewards retrieves loyalty rewards with filters
func (r *LoyaltyRepository) ListRewards(ctx context.Context, orgID uuid.UUID, filters loyalty.LoyaltyRewardFilters) ([]loyalty.LoyaltyReward, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, reward_code, reward_name, description, reward_type,
			points_cost, reward_value, discount_percentage, discount_amount,
			product_id, product_variant_id, is_active,
			available_from, available_to, total_available, total_redeemed,
			max_redemptions_per_customer, minimum_tier_level, tier_ids,
			image_url, thumbnail_url, sort_order, is_featured,
			terms_and_conditions, redemption_instructions, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_rewards
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (reward_name ILIKE $%d OR reward_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.RewardType != nil {
		argCount++
		query += fmt.Sprintf(" AND reward_type = $%d", argCount)
		args = append(args, *filters.RewardType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.IsFeatured != nil {
		argCount++
		query += fmt.Sprintf(" AND is_featured = $%d", argCount)
		args = append(args, *filters.IsFeatured)
	}

	if filters.MinPoints != nil {
		argCount++
		query += fmt.Sprintf(" AND points_cost >= $%d", argCount)
		args = append(args, *filters.MinPoints)
	}

	if filters.MaxPoints != nil {
		argCount++
		query += fmt.Sprintf(" AND points_cost <= $%d", argCount)
		args = append(args, *filters.MaxPoints)
	}

	query += " ORDER BY sort_order ASC, created_at DESC"

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

	var rewards []loyalty.LoyaltyReward
	for rows.Next() {
		var reward loyalty.LoyaltyReward
		err := rows.Scan(
			&reward.ID, &reward.OrganizationID, &reward.RewardCode, &reward.RewardName, &reward.Description, &reward.RewardType,
			&reward.PointsCost, &reward.RewardValue, &reward.DiscountPercentage, &reward.DiscountAmount,
			&reward.ProductID, &reward.ProductVariantID, &reward.IsActive,
			&reward.AvailableFrom, &reward.AvailableTo, &reward.TotalAvailable, &reward.TotalRedeemed,
			&reward.MaxRedemptionsPerCustomer, &reward.MinimumTierLevel, &reward.TierIDs,
			&reward.ImageURL, &reward.ThumbnailURL, &reward.SortOrder, &reward.IsFeatured,
			&reward.TermsAndConditions, &reward.RedemptionInstructions, &reward.Metadata,
			&reward.CreatedAt, &reward.UpdatedAt, &reward.DeletedAt, &reward.CreatedBy, &reward.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, reward)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rewards, nil
}

// CreateReward creates a new loyalty reward
func (r *LoyaltyRepository) CreateReward(ctx context.Context, reward *loyalty.LoyaltyReward) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, reward.OrganizationID.String()); err != nil {
		return err
	}

	tierIDsJSON, _ := json.Marshal(reward.TierIDs)

	query := `
		INSERT INTO loyalty_rewards (
			id, organization_id, reward_code, reward_name, description, reward_type,
			points_cost, reward_value, discount_percentage, discount_amount,
			product_id, product_variant_id, is_active,
			available_from, available_to, total_available, total_redeemed,
			max_redemptions_per_customer, minimum_tier_level, tier_ids,
			image_url, thumbnail_url, sort_order, is_featured,
			terms_and_conditions, redemption_instructions, metadata,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		reward.ID, reward.OrganizationID, reward.RewardCode, reward.RewardName, reward.Description, reward.RewardType,
		reward.PointsCost, reward.RewardValue, reward.DiscountPercentage, reward.DiscountAmount,
		reward.ProductID, reward.ProductVariantID, reward.IsActive,
		reward.AvailableFrom, reward.AvailableTo, reward.TotalAvailable, reward.TotalRedeemed,
		reward.MaxRedemptionsPerCustomer, reward.MinimumTierLevel, tierIDsJSON,
		reward.ImageURL, reward.ThumbnailURL, reward.SortOrder, reward.IsFeatured,
		reward.TermsAndConditions, reward.RedemptionInstructions, reward.Metadata,
		reward.CreatedAt, reward.UpdatedAt, reward.CreatedBy,
	)

	return err
}

// GetReward retrieves a loyalty reward by ID
func (r *LoyaltyRepository) GetReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*loyalty.LoyaltyReward, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, reward_code, reward_name, description, reward_type,
			points_cost, reward_value, discount_percentage, discount_amount,
			product_id, product_variant_id, is_active,
			available_from, available_to, total_available, total_redeemed,
			max_redemptions_per_customer, minimum_tier_level, tier_ids,
			image_url, thumbnail_url, sort_order, is_featured,
			terms_and_conditions, redemption_instructions, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_rewards
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var reward loyalty.LoyaltyReward
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&reward.ID, &reward.OrganizationID, &reward.RewardCode, &reward.RewardName, &reward.Description, &reward.RewardType,
		&reward.PointsCost, &reward.RewardValue, &reward.DiscountPercentage, &reward.DiscountAmount,
		&reward.ProductID, &reward.ProductVariantID, &reward.IsActive,
		&reward.AvailableFrom, &reward.AvailableTo, &reward.TotalAvailable, &reward.TotalRedeemed,
		&reward.MaxRedemptionsPerCustomer, &reward.MinimumTierLevel, &reward.TierIDs,
		&reward.ImageURL, &reward.ThumbnailURL, &reward.SortOrder, &reward.IsFeatured,
		&reward.TermsAndConditions, &reward.RedemptionInstructions, &reward.Metadata,
		&reward.CreatedAt, &reward.UpdatedAt, &reward.DeletedAt, &reward.CreatedBy, &reward.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &reward, nil
}

// UpdateReward updates an existing loyalty reward
func (r *LoyaltyRepository) UpdateReward(ctx context.Context, reward *loyalty.LoyaltyReward) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, reward.OrganizationID.String()); err != nil {
		return err
	}

	tierIDsJSON, _ := json.Marshal(reward.TierIDs)

	query := `
		UPDATE loyalty_rewards SET
			reward_code = $3, reward_name = $4, description = $5, reward_type = $6,
			points_cost = $7, reward_value = $8, discount_percentage = $9, discount_amount = $10,
			product_id = $11, product_variant_id = $12, is_active = $13,
			available_from = $14, available_to = $15, total_available = $16, total_redeemed = $17,
			max_redemptions_per_customer = $18, minimum_tier_level = $19, tier_ids = $20,
			image_url = $21, thumbnail_url = $22, sort_order = $23, is_featured = $24,
			terms_and_conditions = $25, redemption_instructions = $26, metadata = $27,
			updated_at = $28, updated_by = $29
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		reward.OrganizationID, reward.ID,
		reward.RewardCode, reward.RewardName, reward.Description, reward.RewardType,
		reward.PointsCost, reward.RewardValue, reward.DiscountPercentage, reward.DiscountAmount,
		reward.ProductID, reward.ProductVariantID, reward.IsActive,
		reward.AvailableFrom, reward.AvailableTo, reward.TotalAvailable, reward.TotalRedeemed,
		reward.MaxRedemptionsPerCustomer, reward.MinimumTierLevel, tierIDsJSON,
		reward.ImageURL, reward.ThumbnailURL, reward.SortOrder, reward.IsFeatured,
		reward.TermsAndConditions, reward.RedemptionInstructions, reward.Metadata,
		reward.UpdatedAt, reward.UpdatedBy,
	)

	return err
}

// DeleteReward soft-deletes a loyalty reward
func (r *LoyaltyRepository) DeleteReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE loyalty_rewards
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// GetFeaturedRewards retrieves featured rewards
func (r *LoyaltyRepository) GetFeaturedRewards(ctx context.Context, orgID uuid.UUID, limit int) ([]loyalty.LoyaltyReward, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, reward_code, reward_name, description, reward_type,
			points_cost, reward_value, discount_percentage, discount_amount,
			product_id, product_variant_id, is_active,
			available_from, available_to, total_available, total_redeemed,
			max_redemptions_per_customer, minimum_tier_level, tier_ids,
			image_url, thumbnail_url, sort_order, is_featured,
			terms_and_conditions, redemption_instructions, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM loyalty_rewards
		WHERE organization_id = $1 AND deleted_at IS NULL AND is_active = true AND is_featured = true
		ORDER BY sort_order ASC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []loyalty.LoyaltyReward
	for rows.Next() {
		var reward loyalty.LoyaltyReward
		err := rows.Scan(
			&reward.ID, &reward.OrganizationID, &reward.RewardCode, &reward.RewardName, &reward.Description, &reward.RewardType,
			&reward.PointsCost, &reward.RewardValue, &reward.DiscountPercentage, &reward.DiscountAmount,
			&reward.ProductID, &reward.ProductVariantID, &reward.IsActive,
			&reward.AvailableFrom, &reward.AvailableTo, &reward.TotalAvailable, &reward.TotalRedeemed,
			&reward.MaxRedemptionsPerCustomer, &reward.MinimumTierLevel, &reward.TierIDs,
			&reward.ImageURL, &reward.ThumbnailURL, &reward.SortOrder, &reward.IsFeatured,
			&reward.TermsAndConditions, &reward.RedemptionInstructions, &reward.Metadata,
			&reward.CreatedAt, &reward.UpdatedAt, &reward.DeletedAt, &reward.CreatedBy, &reward.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, reward)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rewards, nil
}

// ========================================================================
// LOYALTY REDEMPTIONS
// ========================================================================

// ListRedemptions retrieves loyalty redemptions with filters
func (r *LoyaltyRepository) ListRedemptions(ctx context.Context, orgID uuid.UUID, filters loyalty.LoyaltyRedemptionFilters) ([]loyalty.LoyaltyRedemption, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, reward_id, redemption_number, redemption_date,
			points_redeemed, status, sale_id, used_date, expiry_date,
			fulfillment_status, fulfillment_notes, fulfilled_by, fulfilled_at,
			notes, metadata, created_at, updated_at, created_by, updated_by
		FROM loyalty_redemptions
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.FulfillmentStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND fulfillment_status = $%d", argCount)
		args = append(args, *filters.FulfillmentStatus)
	}

	query += " ORDER BY redemption_date DESC"

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

	var redemptions []loyalty.LoyaltyRedemption
	for rows.Next() {
		var r loyalty.LoyaltyRedemption
		err := rows.Scan(
			&r.ID, &r.OrganizationID, &r.CustomerID, &r.RewardID, &r.RedemptionNumber, &r.RedemptionDate,
			&r.PointsRedeemed, &r.Status, &r.SaleID, &r.UsedDate, &r.ExpiryDate,
			&r.FulfillmentStatus, &r.FulfillmentNotes, &r.FulfilledBy, &r.FulfilledAt,
			&r.Notes, &r.Metadata, &r.CreatedAt, &r.UpdatedAt, &r.CreatedBy, &r.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return redemptions, nil
}

// CreateRedemption creates a new loyalty redemption
func (r *LoyaltyRepository) CreateRedemption(ctx context.Context, redemption *loyalty.LoyaltyRedemption) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, redemption.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO loyalty_redemptions (
			id, organization_id, customer_id, reward_id, redemption_number, redemption_date,
			points_redeemed, status, sale_id, used_date, expiry_date,
			fulfillment_status, fulfillment_notes, fulfilled_by, fulfilled_at,
			notes, metadata, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		redemption.ID, redemption.OrganizationID, redemption.CustomerID, redemption.RewardID,
		redemption.RedemptionNumber, redemption.RedemptionDate, redemption.PointsRedeemed, redemption.Status,
		redemption.SaleID, redemption.UsedDate, redemption.ExpiryDate,
		redemption.FulfillmentStatus, redemption.FulfillmentNotes, redemption.FulfilledBy, redemption.FulfilledAt,
		redemption.Notes, redemption.Metadata, redemption.CreatedAt, redemption.UpdatedAt, redemption.CreatedBy,
	)

	return err
}

// GetRedemption retrieves a loyalty redemption by ID
func (r *LoyaltyRepository) GetRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*loyalty.LoyaltyRedemption, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, reward_id, redemption_number, redemption_date,
			points_redeemed, status, sale_id, used_date, expiry_date,
			fulfillment_status, fulfillment_notes, fulfilled_by, fulfilled_at,
			notes, metadata, created_at, updated_at, created_by, updated_by
		FROM loyalty_redemptions
		WHERE organization_id = $1 AND id = $2
	`

	var redemption loyalty.LoyaltyRedemption
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&redemption.ID, &redemption.OrganizationID, &redemption.CustomerID, &redemption.RewardID,
		&redemption.RedemptionNumber, &redemption.RedemptionDate, &redemption.PointsRedeemed, &redemption.Status,
		&redemption.SaleID, &redemption.UsedDate, &redemption.ExpiryDate,
		&redemption.FulfillmentStatus, &redemption.FulfillmentNotes, &redemption.FulfilledBy, &redemption.FulfilledAt,
		&redemption.Notes, &redemption.Metadata, &redemption.CreatedAt, &redemption.UpdatedAt, &redemption.CreatedBy, &redemption.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &redemption, nil
}

// UpdateRedemption updates an existing loyalty redemption
func (r *LoyaltyRepository) UpdateRedemption(ctx context.Context, redemption *loyalty.LoyaltyRedemption) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, redemption.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE loyalty_redemptions SET
			redemption_number = $3, redemption_date = $4, points_redeemed = $5, status = $6,
			sale_id = $7, used_date = $8, expiry_date = $9,
			fulfillment_status = $10, fulfillment_notes = $11, fulfilled_by = $12, fulfilled_at = $13,
			notes = $14, metadata = $15, updated_at = $16, updated_by = $17
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		redemption.OrganizationID, redemption.ID,
		redemption.RedemptionNumber, redemption.RedemptionDate, redemption.PointsRedeemed, redemption.Status,
		redemption.SaleID, redemption.UsedDate, redemption.ExpiryDate,
		redemption.FulfillmentStatus, redemption.FulfillmentNotes, redemption.FulfilledBy, redemption.FulfilledAt,
		redemption.Notes, redemption.Metadata, redemption.UpdatedAt, redemption.UpdatedBy,
	)

	return err
}

// DeleteRedemption deletes a loyalty redemption
func (r *LoyaltyRepository) DeleteRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "DELETE FROM loyalty_redemptions WHERE organization_id = $1 AND id = $2"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// GetCustomerRedemptions retrieves all redemptions for a customer
func (r *LoyaltyRepository) GetCustomerRedemptions(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) ([]loyalty.LoyaltyRedemption, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, reward_id, redemption_number, redemption_date,
			points_redeemed, status, sale_id, used_date, expiry_date,
			fulfillment_status, fulfillment_notes, fulfilled_by, fulfilled_at,
			notes, metadata, created_at, updated_at, created_by, updated_by
		FROM loyalty_redemptions
		WHERE organization_id = $1 AND customer_id = $2
		ORDER BY redemption_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []loyalty.LoyaltyRedemption
	for rows.Next() {
		var redemption loyalty.LoyaltyRedemption
		err := rows.Scan(
			&redemption.ID, &redemption.OrganizationID, &redemption.CustomerID, &redemption.RewardID,
			&redemption.RedemptionNumber, &redemption.RedemptionDate, &redemption.PointsRedeemed, &redemption.Status,
			&redemption.SaleID, &redemption.UsedDate, &redemption.ExpiryDate,
			&redemption.FulfillmentStatus, &redemption.FulfillmentNotes, &redemption.FulfilledBy, &redemption.FulfilledAt,
			&redemption.Notes, &redemption.Metadata, &redemption.CreatedAt, &redemption.UpdatedAt, &redemption.CreatedBy, &redemption.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, redemption)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return redemptions, nil
}

// ========================================================================
// LOYALTY POINTS TRANSACTIONS
// ========================================================================

// ListPointsTransactions retrieves loyalty points transactions with filters
func (r *LoyaltyRepository) ListPointsTransactions(ctx context.Context, orgID uuid.UUID, filters loyalty.LoyaltyPointsTransactionFilters) ([]loyalty.LoyaltyPointsTransaction, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, transaction_type, points, balance_after,
			sale_id, redemption_id, points_rule_id, description, reason, notes,
			expiry_date, metadata, transaction_date, created_by
		FROM loyalty_points_transactions
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

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

	var transactions []loyalty.LoyaltyPointsTransaction
	for rows.Next() {
		var t loyalty.LoyaltyPointsTransaction
		err := rows.Scan(
			&t.ID, &t.OrganizationID, &t.CustomerID, &t.TransactionType, &t.Points, &t.BalanceAfter,
			&t.SaleID, &t.RedemptionID, &t.PointsRuleID, &t.Description, &t.Reason, &t.Notes,
			&t.ExpiryDate, &t.Metadata, &t.TransactionDate, &t.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// CreatePointsTransaction creates a new loyalty points transaction
func (r *LoyaltyRepository) CreatePointsTransaction(ctx context.Context, transaction *loyalty.LoyaltyPointsTransaction) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, transaction.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO loyalty_points_transactions (
			id, organization_id, customer_id, transaction_type, points, balance_after,
			sale_id, redemption_id, points_rule_id, description, reason, notes,
			expiry_date, metadata, transaction_date, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		transaction.ID, transaction.OrganizationID, transaction.CustomerID, transaction.TransactionType,
		transaction.Points, transaction.BalanceAfter, transaction.SaleID, transaction.RedemptionID,
		transaction.PointsRuleID, transaction.Description, transaction.Reason, transaction.Notes,
		transaction.ExpiryDate, transaction.Metadata, transaction.TransactionDate, transaction.CreatedBy,
	)

	return err
}

// GetPointsTransaction retrieves a loyalty points transaction by ID
func (r *LoyaltyRepository) GetPointsTransaction(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*loyalty.LoyaltyPointsTransaction, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, transaction_type, points, balance_after,
			sale_id, redemption_id, points_rule_id, description, reason, notes,
			expiry_date, metadata, transaction_date, created_by
		FROM loyalty_points_transactions
		WHERE organization_id = $1 AND id = $2
	`

	var t loyalty.LoyaltyPointsTransaction
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&t.ID, &t.OrganizationID, &t.CustomerID, &t.TransactionType, &t.Points, &t.BalanceAfter,
		&t.SaleID, &t.RedemptionID, &t.PointsRuleID, &t.Description, &t.Reason, &t.Notes,
		&t.ExpiryDate, &t.Metadata, &t.TransactionDate, &t.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetCustomerPointsBalance retrieves the current points balance for a customer
func (r *LoyaltyRepository) GetCustomerPointsBalance(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (int, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := `
		SELECT COALESCE(balance_after, 0)
		FROM loyalty_points_transactions
		WHERE organization_id = $1 AND customer_id = $2
		ORDER BY transaction_date DESC
		LIMIT 1
	`

	var balance int
	err := r.db.Pool.QueryRow(ctx, query, orgID, customerID).Scan(&balance)
	if err == pgx.ErrNoRows {
		return 0, nil // No transactions means 0 balance
	}
	if err != nil {
		return 0, err
	}

	return balance, nil
}

// GetCustomerPointsTransactions retrieves all points transactions for a customer
func (r *LoyaltyRepository) GetCustomerPointsTransactions(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) ([]loyalty.LoyaltyPointsTransaction, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, transaction_type, points, balance_after,
			sale_id, redemption_id, points_rule_id, description, reason, notes,
			expiry_date, metadata, transaction_date, created_by
		FROM loyalty_points_transactions
		WHERE organization_id = $1 AND customer_id = $2
		ORDER BY transaction_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []loyalty.LoyaltyPointsTransaction
	for rows.Next() {
		var t loyalty.LoyaltyPointsTransaction
		err := rows.Scan(
			&t.ID, &t.OrganizationID, &t.CustomerID, &t.TransactionType, &t.Points, &t.BalanceAfter,
			&t.SaleID, &t.RedemptionID, &t.PointsRuleID, &t.Description, &t.Reason, &t.Notes,
			&t.ExpiryDate, &t.Metadata, &t.TransactionDate, &t.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// ========================================================================
// CUSTOMER TIER HISTORY
// ========================================================================

// ListTierHistory retrieves customer tier history with filters
func (r *LoyaltyRepository) ListTierHistory(ctx context.Context, orgID uuid.UUID, filters loyalty.CustomerTierHistoryFilters) ([]loyalty.CustomerTierHistory, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, tier_id, previous_tier_id, change_type,
			change_reason, qualifying_points, qualifying_spend, qualifying_purchases,
			effective_date, valid_until, notes, metadata, created_at, created_by
		FROM customer_tier_history
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.TierID != nil {
		argCount++
		query += fmt.Sprintf(" AND tier_id = $%d", argCount)
		args = append(args, *filters.TierID)
	}

	if filters.ChangeType != nil {
		argCount++
		query += fmt.Sprintf(" AND change_type = $%d", argCount)
		args = append(args, *filters.ChangeType)
	}

	query += " ORDER BY effective_date DESC"

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

	var histories []loyalty.CustomerTierHistory
	for rows.Next() {
		var h loyalty.CustomerTierHistory
		err := rows.Scan(
			&h.ID, &h.OrganizationID, &h.CustomerID, &h.TierID, &h.PreviousTierID, &h.ChangeType,
			&h.ChangeReason, &h.QualifyingPoints, &h.QualifyingSpend, &h.QualifyingPurchases,
			&h.EffectiveDate, &h.ValidUntil, &h.Notes, &h.Metadata, &h.CreatedAt, &h.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}

// CreateTierHistory creates a new customer tier history entry
func (r *LoyaltyRepository) CreateTierHistory(ctx context.Context, history *loyalty.CustomerTierHistory) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, history.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO customer_tier_history (
			id, organization_id, customer_id, tier_id, previous_tier_id, change_type,
			change_reason, qualifying_points, qualifying_spend, qualifying_purchases,
			effective_date, valid_until, notes, metadata, created_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		history.ID, history.OrganizationID, history.CustomerID, history.TierID, history.PreviousTierID,
		history.ChangeType, history.ChangeReason, history.QualifyingPoints, history.QualifyingSpend,
		history.QualifyingPurchases, history.EffectiveDate, history.ValidUntil, history.Notes, history.Metadata,
		history.CreatedAt, history.CreatedBy,
	)

	return err
}

// GetTierHistory retrieves a customer tier history entry by ID
func (r *LoyaltyRepository) GetTierHistory(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*loyalty.CustomerTierHistory, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, tier_id, previous_tier_id, change_type,
			change_reason, qualifying_points, qualifying_spend, qualifying_purchases,
			effective_date, valid_until, notes, metadata, created_at, created_by
		FROM customer_tier_history
		WHERE organization_id = $1 AND id = $2
	`

	var h loyalty.CustomerTierHistory
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&h.ID, &h.OrganizationID, &h.CustomerID, &h.TierID, &h.PreviousTierID, &h.ChangeType,
		&h.ChangeReason, &h.QualifyingPoints, &h.QualifyingSpend, &h.QualifyingPurchases,
		&h.EffectiveDate, &h.ValidUntil, &h.Notes, &h.Metadata, &h.CreatedAt, &h.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &h, nil
}

// GetCustomerCurrentTier retrieves the current tier for a customer
func (r *LoyaltyRepository) GetCustomerCurrentTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (*loyalty.LoyaltyTier, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			t.id, t.organization_id, t.tier_code, t.tier_name, t.tier_level, t.description,
			t.points_threshold, t.annual_spend_threshold, t.purchase_count_threshold,
			t.points_multiplier, t.discount_percentage, t.tier_color, t.tier_icon,
			t.badge_image_url, t.is_active, t.is_default, t.sort_order, t.metadata,
			t.created_at, t.updated_at, t.deleted_at, t.created_by, t.updated_by
		FROM loyalty_tiers t
		JOIN customers c ON t.id = c.current_tier_id
		WHERE t.organization_id = $1 AND c.id = $2 AND t.deleted_at IS NULL
		LIMIT 1
	`

	var tier loyalty.LoyaltyTier
	err := r.db.Pool.QueryRow(ctx, query, orgID, customerID).Scan(
		&tier.ID, &tier.OrganizationID, &tier.TierCode, &tier.TierName, &tier.TierLevel, &tier.Description,
		&tier.PointsThreshold, &tier.AnnualSpendThreshold, &tier.PurchaseCountThreshold,
		&tier.PointsMultiplier, &tier.DiscountPercentage, &tier.TierColor, &tier.TierIcon,
		&tier.BadgeImageURL, &tier.IsActive, &tier.IsDefault, &tier.SortOrder, &tier.Metadata,
		&tier.CreatedAt, &tier.UpdatedAt, &tier.DeletedAt, &tier.CreatedBy, &tier.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tier, nil
}

// GetCustomerTierHistory retrieves all tier history entries for a customer
func (r *LoyaltyRepository) GetCustomerTierHistory(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) ([]loyalty.CustomerTierHistory, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, customer_id, tier_id, previous_tier_id, change_type,
			change_reason, qualifying_points, qualifying_spend, qualifying_purchases,
			effective_date, valid_until, notes, metadata, created_at, created_by
		FROM customer_tier_history
		WHERE organization_id = $1 AND customer_id = $2
		ORDER BY effective_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []loyalty.CustomerTierHistory
	for rows.Next() {
		var h loyalty.CustomerTierHistory
		err := rows.Scan(
			&h.ID, &h.OrganizationID, &h.CustomerID, &h.TierID, &h.PreviousTierID, &h.ChangeType,
			&h.ChangeReason, &h.QualifyingPoints, &h.QualifyingSpend, &h.QualifyingPurchases,
			&h.EffectiveDate, &h.ValidUntil, &h.Notes, &h.Metadata, &h.CreatedAt, &h.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}

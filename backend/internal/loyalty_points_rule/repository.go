package loyalty_points_rule

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

// Repository handles database operations for LoyaltyPointsRules
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new LoyaltyPointsRules repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// LoyaltyPointsRules represents a loyalty_points_rules entity
type LoyaltyPointsRules struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	RuleCode string `json:"rule_code" db:"rule_code"`
	RuleName string `json:"rule_name" db:"rule_name"`
	Description *string `json:"description" db:"description"`
	RuleType string `json:"rule_type" db:"rule_type"`
	PointsPerAmount *float64 `json:"points_per_amount" db:"points_per_amount"`
	FixedPoints *int64 `json:"fixed_points" db:"fixed_points"`
	Multiplier *float64 `json:"multiplier" db:"multiplier"`
	AppliesTo *string `json:"applies_to" db:"applies_to"`
	ApplicableProductIds json.RawMessage `json:"applicable_product_ids" db:"applicable_product_ids"`
	ApplicableCategoryIds json.RawMessage `json:"applicable_category_ids" db:"applicable_category_ids"`
	ApplicableTierIds json.RawMessage `json:"applicable_tier_ids" db:"applicable_tier_ids"`
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount" db:"minimum_purchase_amount"`
	MaximumPointsPerTransaction *int64 `json:"maximum_points_per_transaction" db:"maximum_points_per_transaction"`
	MaximumPointsPerDay *int64 `json:"maximum_points_per_day" db:"maximum_points_per_day"`
	MaximumPointsPerMonth *int64 `json:"maximum_points_per_month" db:"maximum_points_per_month"`
	StartDate *time.Time `json:"start_date" db:"start_date"`
	EndDate *time.Time `json:"end_date" db:"end_date"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Priority *int64 `json:"priority" db:"priority"`
	TermsAndConditions *string `json:"terms_and_conditions" db:"terms_and_conditions"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new loyalty_points_rules record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *LoyaltyPointsRules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "loyalty_points_rules", duration, nil)
	}()

	query := `
		INSERT INTO loyalty_points_rules (
			, organization_id
			, rule_code
			, rule_name
			, description
			, rule_type
			, points_per_amount
			, fixed_points
			, multiplier
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, applicable_tier_ids
			, minimum_purchase_amount
			, maximum_points_per_transaction
			, maximum_points_per_day
			, maximum_points_per_month
			, start_date
			, end_date
			, is_active
			, priority
			, terms_and_conditions
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
			, $22
			, $23
			, $26
			, $27
			, $28
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.RuleCode,
		entity.RuleName,
		entity.Description,
		entity.RuleType,
		entity.PointsPerAmount,
		entity.FixedPoints,
		entity.Multiplier,
		entity.AppliesTo,
		entity.ApplicableProductIds,
		entity.ApplicableCategoryIds,
		entity.ApplicableTierIds,
		entity.MinimumPurchaseAmount,
		entity.MaximumPointsPerTransaction,
		entity.MaximumPointsPerDay,
		entity.MaximumPointsPerMonth,
		entity.StartDate,
		entity.EndDate,
		entity.IsActive,
		entity.Priority,
		entity.TermsAndConditions,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create loyalty_points_rules", zap.Error(err))
		return fmt.Errorf("failed to create loyalty_points_rules: %w", err)
	}

	r.logger.Info("created loyalty_points_rules",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a loyalty_points_rules by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*LoyaltyPointsRules, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_points_rules", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, rule_code
			, rule_name
			, description
			, rule_type
			, points_per_amount
			, fixed_points
			, multiplier
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, applicable_tier_ids
			, minimum_purchase_amount
			, maximum_points_per_transaction
			, maximum_points_per_day
			, maximum_points_per_month
			, start_date
			, end_date
			, is_active
			, priority
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM loyalty_points_rules
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity LoyaltyPointsRules
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.RuleCode,
		&entity.RuleName,
		&entity.Description,
		&entity.RuleType,
		&entity.PointsPerAmount,
		&entity.FixedPoints,
		&entity.Multiplier,
		&entity.AppliesTo,
		&entity.ApplicableProductIds,
		&entity.ApplicableCategoryIds,
		&entity.ApplicableTierIds,
		&entity.MinimumPurchaseAmount,
		&entity.MaximumPointsPerTransaction,
		&entity.MaximumPointsPerDay,
		&entity.MaximumPointsPerMonth,
		&entity.StartDate,
		&entity.EndDate,
		&entity.IsActive,
		&entity.Priority,
		&entity.TermsAndConditions,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("loyalty_points_rules not found")
	}

	if err != nil {
		r.logger.Error("failed to get loyalty_points_rules", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get loyalty_points_rules: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of loyalty_points_rules records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*LoyaltyPointsRules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_points_rules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_points_rules
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_points_rules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, rule_code
			, rule_name
			, description
			, rule_type
			, points_per_amount
			, fixed_points
			, multiplier
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, applicable_tier_ids
			, minimum_purchase_amount
			, maximum_points_per_transaction
			, maximum_points_per_day
			, maximum_points_per_month
			, start_date
			, end_date
			, is_active
			, priority
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM loyalty_points_rules
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_points_rules", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_points_rules: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyPointsRules
	for rows.Next() {
		var entity LoyaltyPointsRules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.RuleCode,
			&entity.RuleName,
			&entity.Description,
			&entity.RuleType,
			&entity.PointsPerAmount,
			&entity.FixedPoints,
			&entity.Multiplier,
			&entity.AppliesTo,
			&entity.ApplicableProductIds,
			&entity.ApplicableCategoryIds,
			&entity.ApplicableTierIds,
			&entity.MinimumPurchaseAmount,
			&entity.MaximumPointsPerTransaction,
			&entity.MaximumPointsPerDay,
			&entity.MaximumPointsPerMonth,
			&entity.StartDate,
			&entity.EndDate,
			&entity.IsActive,
			&entity.Priority,
			&entity.TermsAndConditions,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_points_rules: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating loyalty_points_rules rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing loyalty_points_rules record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *LoyaltyPointsRules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_points_rules", duration, nil)
	}()

	query := `
		UPDATE loyalty_points_rules
		SET
			, organization_id = $2
			, rule_code = $3
			, rule_name = $4
			, description = $5
			, rule_type = $6
			, points_per_amount = $7
			, fixed_points = $8
			, multiplier = $9
			, applies_to = $10
			, applicable_product_ids = $11
			, applicable_category_ids = $12
			, applicable_tier_ids = $13
			, minimum_purchase_amount = $14
			, maximum_points_per_transaction = $15
			, maximum_points_per_day = $16
			, maximum_points_per_month = $17
			, start_date = $18
			, end_date = $19
			, is_active = $20
			, priority = $21
			, terms_and_conditions = $22
			, metadata = $23
			, updated_at = $25
			, deleted_at = $26
			, created_by = $27
			, updated_by = $28
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $29
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.RuleCode,
		entity.RuleName,
		entity.Description,
		entity.RuleType,
		entity.PointsPerAmount,
		entity.FixedPoints,
		entity.Multiplier,
		entity.AppliesTo,
		entity.ApplicableProductIds,
		entity.ApplicableCategoryIds,
		entity.ApplicableTierIds,
		entity.MinimumPurchaseAmount,
		entity.MaximumPointsPerTransaction,
		entity.MaximumPointsPerDay,
		entity.MaximumPointsPerMonth,
		entity.StartDate,
		entity.EndDate,
		entity.IsActive,
		entity.Priority,
		entity.TermsAndConditions,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update loyalty_points_rules", zap.Error(err))
		return fmt.Errorf("failed to update loyalty_points_rules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_points_rules not found or already deleted")
	}

	r.logger.Info("updated loyalty_points_rules",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a loyalty_points_rules record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_points_rules", duration, nil)
	}()

	query := `
		UPDATE loyalty_points_rules
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete loyalty_points_rules", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete loyalty_points_rules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_points_rules not found or already deleted")
	}

	r.logger.Info("deleted loyalty_points_rules", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves loyalty_points_rules records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*LoyaltyPointsRules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_points_rules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_points_rules
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_points_rules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, rule_code
			, rule_name
			, description
			, rule_type
			, points_per_amount
			, fixed_points
			, multiplier
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, applicable_tier_ids
			, minimum_purchase_amount
			, maximum_points_per_transaction
			, maximum_points_per_day
			, maximum_points_per_month
			, start_date
			, end_date
			, is_active
			, priority
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM loyalty_points_rules
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_points_rules by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_points_rules: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyPointsRules
	for rows.Next() {
		var entity LoyaltyPointsRules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.RuleCode,
			&entity.RuleName,
			&entity.Description,
			&entity.RuleType,
			&entity.PointsPerAmount,
			&entity.FixedPoints,
			&entity.Multiplier,
			&entity.AppliesTo,
			&entity.ApplicableProductIds,
			&entity.ApplicableCategoryIds,
			&entity.ApplicableTierIds,
			&entity.MinimumPurchaseAmount,
			&entity.MaximumPointsPerTransaction,
			&entity.MaximumPointsPerDay,
			&entity.MaximumPointsPerMonth,
			&entity.StartDate,
			&entity.EndDate,
			&entity.IsActive,
			&entity.Priority,
			&entity.TermsAndConditions,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_points_rules: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


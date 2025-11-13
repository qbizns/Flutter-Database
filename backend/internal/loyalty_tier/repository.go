package loyalty_tier

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

// Repository handles database operations for LoyaltyTiers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new LoyaltyTiers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// LoyaltyTiers represents a loyalty_tiers entity
type LoyaltyTiers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	TierCode string `json:"tier_code" db:"tier_code"`
	TierName string `json:"tier_name" db:"tier_name"`
	TierLevel int64 `json:"tier_level" db:"tier_level"`
	Description *string `json:"description" db:"description"`
	PointsThreshold int64 `json:"points_threshold" db:"points_threshold"`
	AnnualSpendThreshold *float64 `json:"annual_spend_threshold" db:"annual_spend_threshold"`
	PurchaseCountThreshold *int64 `json:"purchase_count_threshold" db:"purchase_count_threshold"`
	PointsMultiplier *float64 `json:"points_multiplier" db:"points_multiplier"`
	DiscountPercentage *float64 `json:"discount_percentage" db:"discount_percentage"`
	TierColor *string `json:"tier_color" db:"tier_color"`
	TierIcon *string `json:"tier_icon" db:"tier_icon"`
	BadgeImageUrl *string `json:"badge_image_url" db:"badge_image_url"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	PointsThreshold *string `json:"points_threshold" db:"points_threshold"`
	(annualSpendThreshold *string `json:"(annual_spend_threshold" db:"(annual_spend_threshold"`
	(purchaseCountThreshold *string `json:"(purchase_count_threshold" db:"(purchase_count_threshold"`
}

// Create inserts a new loyalty_tiers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *LoyaltyTiers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "loyalty_tiers", duration, nil)
	}()

	query := `
		INSERT INTO loyalty_tiers (
			, organization_id
			, tier_code
			, tier_name
			, tier_level
			, description
			, points_threshold
			, annual_spend_threshold
			, purchase_count_threshold
			, points_multiplier
			, discount_percentage
			, tier_color
			, tier_icon
			, badge_image_url
			, is_active
			, is_default
			, sort_order
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, points_threshold
			, (annual_spend_threshold
			, (purchase_count_threshold
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
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.TierCode,
		entity.TierName,
		entity.TierLevel,
		entity.Description,
		entity.PointsThreshold,
		entity.AnnualSpendThreshold,
		entity.PurchaseCountThreshold,
		entity.PointsMultiplier,
		entity.DiscountPercentage,
		entity.TierColor,
		entity.TierIcon,
		entity.BadgeImageUrl,
		entity.IsActive,
		entity.IsDefault,
		entity.SortOrder,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.PointsThreshold,
		entity.(annualSpendThreshold,
		entity.(purchaseCountThreshold,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create loyalty_tiers", zap.Error(err))
		return fmt.Errorf("failed to create loyalty_tiers: %w", err)
	}

	r.logger.Info("created loyalty_tiers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a loyalty_tiers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*LoyaltyTiers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_tiers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, tier_code
			, tier_name
			, tier_level
			, description
			, points_threshold
			, annual_spend_threshold
			, purchase_count_threshold
			, points_multiplier
			, discount_percentage
			, tier_color
			, tier_icon
			, badge_image_url
			, is_active
			, is_default
			, sort_order
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, points_threshold
			, (annual_spend_threshold
			, (purchase_count_threshold
		FROM loyalty_tiers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity LoyaltyTiers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.TierCode,
		&entity.TierName,
		&entity.TierLevel,
		&entity.Description,
		&entity.PointsThreshold,
		&entity.AnnualSpendThreshold,
		&entity.PurchaseCountThreshold,
		&entity.PointsMultiplier,
		&entity.DiscountPercentage,
		&entity.TierColor,
		&entity.TierIcon,
		&entity.BadgeImageUrl,
		&entity.IsActive,
		&entity.IsDefault,
		&entity.SortOrder,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.PointsThreshold,
		&entity.(annualSpendThreshold,
		&entity.(purchaseCountThreshold,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("loyalty_tiers not found")
	}

	if err != nil {
		r.logger.Error("failed to get loyalty_tiers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get loyalty_tiers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of loyalty_tiers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*LoyaltyTiers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_tiers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_tiers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_tiers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, tier_code
			, tier_name
			, tier_level
			, description
			, points_threshold
			, annual_spend_threshold
			, purchase_count_threshold
			, points_multiplier
			, discount_percentage
			, tier_color
			, tier_icon
			, badge_image_url
			, is_active
			, is_default
			, sort_order
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, points_threshold
			, (annual_spend_threshold
			, (purchase_count_threshold
		FROM loyalty_tiers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_tiers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_tiers: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyTiers
	for rows.Next() {
		var entity LoyaltyTiers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TierCode,
			&entity.TierName,
			&entity.TierLevel,
			&entity.Description,
			&entity.PointsThreshold,
			&entity.AnnualSpendThreshold,
			&entity.PurchaseCountThreshold,
			&entity.PointsMultiplier,
			&entity.DiscountPercentage,
			&entity.TierColor,
			&entity.TierIcon,
			&entity.BadgeImageUrl,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.SortOrder,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.PointsThreshold,
			&entity.(annualSpendThreshold,
			&entity.(purchaseCountThreshold,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_tiers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating loyalty_tiers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing loyalty_tiers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *LoyaltyTiers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_tiers", duration, nil)
	}()

	query := `
		UPDATE loyalty_tiers
		SET
			, organization_id = $2
			, tier_code = $3
			, tier_name = $4
			, tier_level = $5
			, description = $6
			, points_threshold = $7
			, annual_spend_threshold = $8
			, purchase_count_threshold = $9
			, points_multiplier = $10
			, discount_percentage = $11
			, tier_color = $12
			, tier_icon = $13
			, badge_image_url = $14
			, is_active = $15
			, is_default = $16
			, sort_order = $17
			, metadata = $18
			, updated_at = $20
			, deleted_at = $21
			, created_by = $22
			, updated_by = $23
			, points_threshold = $24
			, (annual_spend_threshold = $25
			, (purchase_count_threshold = $26
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $27
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.TierCode,
		entity.TierName,
		entity.TierLevel,
		entity.Description,
		entity.PointsThreshold,
		entity.AnnualSpendThreshold,
		entity.PurchaseCountThreshold,
		entity.PointsMultiplier,
		entity.DiscountPercentage,
		entity.TierColor,
		entity.TierIcon,
		entity.BadgeImageUrl,
		entity.IsActive,
		entity.IsDefault,
		entity.SortOrder,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.PointsThreshold,
		entity.(annualSpendThreshold,
		entity.(purchaseCountThreshold,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update loyalty_tiers", zap.Error(err))
		return fmt.Errorf("failed to update loyalty_tiers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_tiers not found or already deleted")
	}

	r.logger.Info("updated loyalty_tiers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a loyalty_tiers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_tiers", duration, nil)
	}()

	query := `
		UPDATE loyalty_tiers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete loyalty_tiers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete loyalty_tiers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_tiers not found or already deleted")
	}

	r.logger.Info("deleted loyalty_tiers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves loyalty_tiers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*LoyaltyTiers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_tiers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_tiers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_tiers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, tier_code
			, tier_name
			, tier_level
			, description
			, points_threshold
			, annual_spend_threshold
			, purchase_count_threshold
			, points_multiplier
			, discount_percentage
			, tier_color
			, tier_icon
			, badge_image_url
			, is_active
			, is_default
			, sort_order
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, points_threshold
			, (annual_spend_threshold
			, (purchase_count_threshold
		FROM loyalty_tiers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_tiers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_tiers: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyTiers
	for rows.Next() {
		var entity LoyaltyTiers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TierCode,
			&entity.TierName,
			&entity.TierLevel,
			&entity.Description,
			&entity.PointsThreshold,
			&entity.AnnualSpendThreshold,
			&entity.PurchaseCountThreshold,
			&entity.PointsMultiplier,
			&entity.DiscountPercentage,
			&entity.TierColor,
			&entity.TierIcon,
			&entity.BadgeImageUrl,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.SortOrder,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.PointsThreshold,
			&entity.(annualSpendThreshold,
			&entity.(purchaseCountThreshold,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_tiers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


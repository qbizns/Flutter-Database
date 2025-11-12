package loyalty_reward

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

// Repository handles database operations for LoyaltyRewards
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new LoyaltyRewards repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// LoyaltyRewards represents a loyalty_rewards entity
type LoyaltyRewards struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	RewardCode string `json:"reward_code" db:"reward_code"`
	RewardName string `json:"reward_name" db:"reward_name"`
	Description *string `json:"description" db:"description"`
	RewardType string `json:"reward_type" db:"reward_type"`
	PointsCost int64 `json:"points_cost" db:"points_cost"`
	RewardValue *float64 `json:"reward_value" db:"reward_value"`
	DiscountPercentage *float64 `json:"discount_percentage" db:"discount_percentage"`
	DiscountAmount *float64 `json:"discount_amount" db:"discount_amount"`
	ProductId *uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	IsActive *bool `json:"is_active" db:"is_active"`
	AvailableFrom *time.Time `json:"available_from" db:"available_from"`
	AvailableTo *time.Time `json:"available_to" db:"available_to"`
	TotalAvailable *int64 `json:"total_available" db:"total_available"`
	TotalRedeemed *int64 `json:"total_redeemed" db:"total_redeemed"`
	MaxRedemptionsPerCustomer *int64 `json:"max_redemptions_per_customer" db:"max_redemptions_per_customer"`
	MinimumTierLevel *int64 `json:"minimum_tier_level" db:"minimum_tier_level"`
	TierIds json.RawMessage `json:"tier_ids" db:"tier_ids"`
	ImageUrl *string `json:"image_url" db:"image_url"`
	ThumbnailUrl *string `json:"thumbnail_url" db:"thumbnail_url"`
	Featured *bool `json:"featured" db:"featured"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	IsFeatured *bool `json:"is_featured" db:"is_featured"`
	TermsAndConditions *string `json:"terms_and_conditions" db:"terms_and_conditions"`
	RedemptionInstructions *string `json:"redemption_instructions" db:"redemption_instructions"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	AvailableTo *string `json:"available_to" db:"available_to"`
	TotalRedeemed *string `json:"total_redeemed" db:"total_redeemed"`
	(totalAvailable *string `json:"(total_available" db:"(total_available"`
}

// Create inserts a new loyalty_rewards record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *LoyaltyRewards) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "loyalty_rewards", duration, nil)
	}()

	query := `
		INSERT INTO loyalty_rewards (
			, organization_id
			, reward_code
			, reward_name
			, description
			, reward_type
			, points_cost
			, reward_value
			, discount_percentage
			, discount_amount
			, product_id
			, product_variant_id
			, is_active
			, available_from
			, available_to
			, total_available
			, total_redeemed
			, max_redemptions_per_customer
			, minimum_tier_level
			, tier_ids
			, image_url
			, thumbnail_url
			, featured
			, sort_order
			, is_featured
			, terms_and_conditions
			, redemption_instructions
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, available_to
			, total_redeemed
			, (total_available
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
			, $24
			, $25
			, $26
			, $27
			, $28
			, $31
			, $32
			, $33
			, $34
			, $35
			, $36
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.RewardCode,
		entity.RewardName,
		entity.Description,
		entity.RewardType,
		entity.PointsCost,
		entity.RewardValue,
		entity.DiscountPercentage,
		entity.DiscountAmount,
		entity.ProductId,
		entity.ProductVariantId,
		entity.IsActive,
		entity.AvailableFrom,
		entity.AvailableTo,
		entity.TotalAvailable,
		entity.TotalRedeemed,
		entity.MaxRedemptionsPerCustomer,
		entity.MinimumTierLevel,
		entity.TierIds,
		entity.ImageUrl,
		entity.ThumbnailUrl,
		entity.Featured,
		entity.SortOrder,
		entity.IsFeatured,
		entity.TermsAndConditions,
		entity.RedemptionInstructions,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.AvailableTo,
		entity.TotalRedeemed,
		entity.(totalAvailable,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create loyalty_rewards", zap.Error(err))
		return fmt.Errorf("failed to create loyalty_rewards: %w", err)
	}

	r.logger.Info("created loyalty_rewards",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a loyalty_rewards by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*LoyaltyRewards, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_rewards", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, reward_code
			, reward_name
			, description
			, reward_type
			, points_cost
			, reward_value
			, discount_percentage
			, discount_amount
			, product_id
			, product_variant_id
			, is_active
			, available_from
			, available_to
			, total_available
			, total_redeemed
			, max_redemptions_per_customer
			, minimum_tier_level
			, tier_ids
			, image_url
			, thumbnail_url
			, featured
			, sort_order
			, is_featured
			, terms_and_conditions
			, redemption_instructions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, available_to
			, total_redeemed
			, (total_available
		FROM loyalty_rewards
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity LoyaltyRewards
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.RewardCode,
		&entity.RewardName,
		&entity.Description,
		&entity.RewardType,
		&entity.PointsCost,
		&entity.RewardValue,
		&entity.DiscountPercentage,
		&entity.DiscountAmount,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.IsActive,
		&entity.AvailableFrom,
		&entity.AvailableTo,
		&entity.TotalAvailable,
		&entity.TotalRedeemed,
		&entity.MaxRedemptionsPerCustomer,
		&entity.MinimumTierLevel,
		&entity.TierIds,
		&entity.ImageUrl,
		&entity.ThumbnailUrl,
		&entity.Featured,
		&entity.SortOrder,
		&entity.IsFeatured,
		&entity.TermsAndConditions,
		&entity.RedemptionInstructions,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.AvailableTo,
		&entity.TotalRedeemed,
		&entity.(totalAvailable,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("loyalty_rewards not found")
	}

	if err != nil {
		r.logger.Error("failed to get loyalty_rewards", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get loyalty_rewards: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of loyalty_rewards records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*LoyaltyRewards, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_rewards", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_rewards
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_rewards records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, reward_code
			, reward_name
			, description
			, reward_type
			, points_cost
			, reward_value
			, discount_percentage
			, discount_amount
			, product_id
			, product_variant_id
			, is_active
			, available_from
			, available_to
			, total_available
			, total_redeemed
			, max_redemptions_per_customer
			, minimum_tier_level
			, tier_ids
			, image_url
			, thumbnail_url
			, featured
			, sort_order
			, is_featured
			, terms_and_conditions
			, redemption_instructions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, available_to
			, total_redeemed
			, (total_available
		FROM loyalty_rewards
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_rewards", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_rewards: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyRewards
	for rows.Next() {
		var entity LoyaltyRewards
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.RewardCode,
			&entity.RewardName,
			&entity.Description,
			&entity.RewardType,
			&entity.PointsCost,
			&entity.RewardValue,
			&entity.DiscountPercentage,
			&entity.DiscountAmount,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.IsActive,
			&entity.AvailableFrom,
			&entity.AvailableTo,
			&entity.TotalAvailable,
			&entity.TotalRedeemed,
			&entity.MaxRedemptionsPerCustomer,
			&entity.MinimumTierLevel,
			&entity.TierIds,
			&entity.ImageUrl,
			&entity.ThumbnailUrl,
			&entity.Featured,
			&entity.SortOrder,
			&entity.IsFeatured,
			&entity.TermsAndConditions,
			&entity.RedemptionInstructions,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.AvailableTo,
			&entity.TotalRedeemed,
			&entity.(totalAvailable,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_rewards: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating loyalty_rewards rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing loyalty_rewards record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *LoyaltyRewards) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_rewards", duration, nil)
	}()

	query := `
		UPDATE loyalty_rewards
		SET
			, organization_id = $2
			, reward_code = $3
			, reward_name = $4
			, description = $5
			, reward_type = $6
			, points_cost = $7
			, reward_value = $8
			, discount_percentage = $9
			, discount_amount = $10
			, product_id = $11
			, product_variant_id = $12
			, is_active = $13
			, available_from = $14
			, available_to = $15
			, total_available = $16
			, total_redeemed = $17
			, max_redemptions_per_customer = $18
			, minimum_tier_level = $19
			, tier_ids = $20
			, image_url = $21
			, thumbnail_url = $22
			, featured = $23
			, sort_order = $24
			, is_featured = $25
			, terms_and_conditions = $26
			, redemption_instructions = $27
			, metadata = $28
			, updated_at = $30
			, deleted_at = $31
			, created_by = $32
			, updated_by = $33
			, available_to = $34
			, total_redeemed = $35
			, (total_available = $36
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $37
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.RewardCode,
		entity.RewardName,
		entity.Description,
		entity.RewardType,
		entity.PointsCost,
		entity.RewardValue,
		entity.DiscountPercentage,
		entity.DiscountAmount,
		entity.ProductId,
		entity.ProductVariantId,
		entity.IsActive,
		entity.AvailableFrom,
		entity.AvailableTo,
		entity.TotalAvailable,
		entity.TotalRedeemed,
		entity.MaxRedemptionsPerCustomer,
		entity.MinimumTierLevel,
		entity.TierIds,
		entity.ImageUrl,
		entity.ThumbnailUrl,
		entity.Featured,
		entity.SortOrder,
		entity.IsFeatured,
		entity.TermsAndConditions,
		entity.RedemptionInstructions,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.AvailableTo,
		entity.TotalRedeemed,
		entity.(totalAvailable,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update loyalty_rewards", zap.Error(err))
		return fmt.Errorf("failed to update loyalty_rewards: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_rewards not found or already deleted")
	}

	r.logger.Info("updated loyalty_rewards",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a loyalty_rewards record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_rewards", duration, nil)
	}()

	query := `
		UPDATE loyalty_rewards
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete loyalty_rewards", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete loyalty_rewards: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_rewards not found or already deleted")
	}

	r.logger.Info("deleted loyalty_rewards", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves loyalty_rewards records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*LoyaltyRewards, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_rewards", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_rewards
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_rewards records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, reward_code
			, reward_name
			, description
			, reward_type
			, points_cost
			, reward_value
			, discount_percentage
			, discount_amount
			, product_id
			, product_variant_id
			, is_active
			, available_from
			, available_to
			, total_available
			, total_redeemed
			, max_redemptions_per_customer
			, minimum_tier_level
			, tier_ids
			, image_url
			, thumbnail_url
			, featured
			, sort_order
			, is_featured
			, terms_and_conditions
			, redemption_instructions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, available_to
			, total_redeemed
			, (total_available
		FROM loyalty_rewards
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_rewards by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_rewards: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyRewards
	for rows.Next() {
		var entity LoyaltyRewards
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.RewardCode,
			&entity.RewardName,
			&entity.Description,
			&entity.RewardType,
			&entity.PointsCost,
			&entity.RewardValue,
			&entity.DiscountPercentage,
			&entity.DiscountAmount,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.IsActive,
			&entity.AvailableFrom,
			&entity.AvailableTo,
			&entity.TotalAvailable,
			&entity.TotalRedeemed,
			&entity.MaxRedemptionsPerCustomer,
			&entity.MinimumTierLevel,
			&entity.TierIds,
			&entity.ImageUrl,
			&entity.ThumbnailUrl,
			&entity.Featured,
			&entity.SortOrder,
			&entity.IsFeatured,
			&entity.TermsAndConditions,
			&entity.RedemptionInstructions,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.AvailableTo,
			&entity.TotalRedeemed,
			&entity.(totalAvailable,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_rewards: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


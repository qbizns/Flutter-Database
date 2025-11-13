package promotion

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

// Repository handles database operations for Promotions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Promotions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Promotions represents a promotions entity
type Promotions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PromotionCode string `json:"promotion_code" db:"promotion_code"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	PromotionType string `json:"promotion_type" db:"promotion_type"`
	DiscountValue float64 `json:"discount_value" db:"discount_value"`
	AppliesTo *string `json:"applies_to" db:"applies_to"`
	ApplicableProductIds json.RawMessage `json:"applicable_product_ids" db:"applicable_product_ids"`
	ApplicableCategoryIds json.RawMessage `json:"applicable_category_ids" db:"applicable_category_ids"`
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount" db:"minimum_purchase_amount"`
	MinimumQuantity *int64 `json:"minimum_quantity" db:"minimum_quantity"`
	BuyQuantity *int64 `json:"buy_quantity" db:"buy_quantity"`
	GetQuantity *int64 `json:"get_quantity" db:"get_quantity"`
	GetDiscountPercentage *float64 `json:"get_discount_percentage" db:"get_discount_percentage"`
	MaxUsesTotal *int64 `json:"max_uses_total" db:"max_uses_total"`
	MaxUsesPerCustomer *int64 `json:"max_uses_per_customer" db:"max_uses_per_customer"`
	CurrentUses *int64 `json:"current_uses" db:"current_uses"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate *time.Time `json:"end_date" db:"end_date"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsCombinable *bool `json:"is_combinable" db:"is_combinable"`
	Priority *int64 `json:"priority" db:"priority"`
	TermsAndConditions *string `json:"terms_and_conditions" db:"terms_and_conditions"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new promotions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Promotions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "promotions", duration, nil)
	}()

	query := `
		INSERT INTO promotions (
			, organization_id
			, promotion_code
			, name
			, description
			, promotion_type
			, discount_value
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, minimum_purchase_amount
			, minimum_quantity
			, buy_quantity
			, get_quantity
			, get_discount_percentage
			, max_uses_total
			, max_uses_per_customer
			, current_uses
			, start_date
			, end_date
			, is_active
			, is_combinable
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
			, $24
			, $25
			, $28
			, $29
			, $30
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PromotionCode,
		entity.Name,
		entity.Description,
		entity.PromotionType,
		entity.DiscountValue,
		entity.AppliesTo,
		entity.ApplicableProductIds,
		entity.ApplicableCategoryIds,
		entity.MinimumPurchaseAmount,
		entity.MinimumQuantity,
		entity.BuyQuantity,
		entity.GetQuantity,
		entity.GetDiscountPercentage,
		entity.MaxUsesTotal,
		entity.MaxUsesPerCustomer,
		entity.CurrentUses,
		entity.StartDate,
		entity.EndDate,
		entity.IsActive,
		entity.IsCombinable,
		entity.Priority,
		entity.TermsAndConditions,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create promotions", zap.Error(err))
		return fmt.Errorf("failed to create promotions: %w", err)
	}

	r.logger.Info("created promotions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a promotions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Promotions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "promotions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, promotion_code
			, name
			, description
			, promotion_type
			, discount_value
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, minimum_purchase_amount
			, minimum_quantity
			, buy_quantity
			, get_quantity
			, get_discount_percentage
			, max_uses_total
			, max_uses_per_customer
			, current_uses
			, start_date
			, end_date
			, is_active
			, is_combinable
			, priority
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM promotions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Promotions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PromotionCode,
		&entity.Name,
		&entity.Description,
		&entity.PromotionType,
		&entity.DiscountValue,
		&entity.AppliesTo,
		&entity.ApplicableProductIds,
		&entity.ApplicableCategoryIds,
		&entity.MinimumPurchaseAmount,
		&entity.MinimumQuantity,
		&entity.BuyQuantity,
		&entity.GetQuantity,
		&entity.GetDiscountPercentage,
		&entity.MaxUsesTotal,
		&entity.MaxUsesPerCustomer,
		&entity.CurrentUses,
		&entity.StartDate,
		&entity.EndDate,
		&entity.IsActive,
		&entity.IsCombinable,
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
		return nil, fmt.Errorf("promotions not found")
	}

	if err != nil {
		r.logger.Error("failed to get promotions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get promotions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of promotions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Promotions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "promotions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM promotions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count promotions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, promotion_code
			, name
			, description
			, promotion_type
			, discount_value
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, minimum_purchase_amount
			, minimum_quantity
			, buy_quantity
			, get_quantity
			, get_discount_percentage
			, max_uses_total
			, max_uses_per_customer
			, current_uses
			, start_date
			, end_date
			, is_active
			, is_combinable
			, priority
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM promotions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list promotions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list promotions: %w", err)
	}
	defer rows.Close()

	var entities []*Promotions
	for rows.Next() {
		var entity Promotions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PromotionCode,
			&entity.Name,
			&entity.Description,
			&entity.PromotionType,
			&entity.DiscountValue,
			&entity.AppliesTo,
			&entity.ApplicableProductIds,
			&entity.ApplicableCategoryIds,
			&entity.MinimumPurchaseAmount,
			&entity.MinimumQuantity,
			&entity.BuyQuantity,
			&entity.GetQuantity,
			&entity.GetDiscountPercentage,
			&entity.MaxUsesTotal,
			&entity.MaxUsesPerCustomer,
			&entity.CurrentUses,
			&entity.StartDate,
			&entity.EndDate,
			&entity.IsActive,
			&entity.IsCombinable,
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
			return nil, 0, fmt.Errorf("failed to scan promotions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating promotions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing promotions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Promotions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "promotions", duration, nil)
	}()

	query := `
		UPDATE promotions
		SET
			, organization_id = $2
			, promotion_code = $3
			, name = $4
			, description = $5
			, promotion_type = $6
			, discount_value = $7
			, applies_to = $8
			, applicable_product_ids = $9
			, applicable_category_ids = $10
			, minimum_purchase_amount = $11
			, minimum_quantity = $12
			, buy_quantity = $13
			, get_quantity = $14
			, get_discount_percentage = $15
			, max_uses_total = $16
			, max_uses_per_customer = $17
			, current_uses = $18
			, start_date = $19
			, end_date = $20
			, is_active = $21
			, is_combinable = $22
			, priority = $23
			, terms_and_conditions = $24
			, metadata = $25
			, updated_at = $27
			, deleted_at = $28
			, created_by = $29
			, updated_by = $30
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $31
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PromotionCode,
		entity.Name,
		entity.Description,
		entity.PromotionType,
		entity.DiscountValue,
		entity.AppliesTo,
		entity.ApplicableProductIds,
		entity.ApplicableCategoryIds,
		entity.MinimumPurchaseAmount,
		entity.MinimumQuantity,
		entity.BuyQuantity,
		entity.GetQuantity,
		entity.GetDiscountPercentage,
		entity.MaxUsesTotal,
		entity.MaxUsesPerCustomer,
		entity.CurrentUses,
		entity.StartDate,
		entity.EndDate,
		entity.IsActive,
		entity.IsCombinable,
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
		r.logger.Error("failed to update promotions", zap.Error(err))
		return fmt.Errorf("failed to update promotions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("promotions not found or already deleted")
	}

	r.logger.Info("updated promotions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a promotions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "promotions", duration, nil)
	}()

	query := `
		UPDATE promotions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete promotions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete promotions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("promotions not found or already deleted")
	}

	r.logger.Info("deleted promotions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves promotions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Promotions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "promotions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM promotions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count promotions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, promotion_code
			, name
			, description
			, promotion_type
			, discount_value
			, applies_to
			, applicable_product_ids
			, applicable_category_ids
			, minimum_purchase_amount
			, minimum_quantity
			, buy_quantity
			, get_quantity
			, get_discount_percentage
			, max_uses_total
			, max_uses_per_customer
			, current_uses
			, start_date
			, end_date
			, is_active
			, is_combinable
			, priority
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM promotions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list promotions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list promotions: %w", err)
	}
	defer rows.Close()

	var entities []*Promotions
	for rows.Next() {
		var entity Promotions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PromotionCode,
			&entity.Name,
			&entity.Description,
			&entity.PromotionType,
			&entity.DiscountValue,
			&entity.AppliesTo,
			&entity.ApplicableProductIds,
			&entity.ApplicableCategoryIds,
			&entity.MinimumPurchaseAmount,
			&entity.MinimumQuantity,
			&entity.BuyQuantity,
			&entity.GetQuantity,
			&entity.GetDiscountPercentage,
			&entity.MaxUsesTotal,
			&entity.MaxUsesPerCustomer,
			&entity.CurrentUses,
			&entity.StartDate,
			&entity.EndDate,
			&entity.IsActive,
			&entity.IsCombinable,
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
			return nil, 0, fmt.Errorf("failed to scan promotions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


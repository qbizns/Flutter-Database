package loyalty_tier_benefit

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

// Repository handles database operations for LoyaltyTierBenefits
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new LoyaltyTierBenefits repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// LoyaltyTierBenefits represents a loyalty_tier_benefits entity
type LoyaltyTierBenefits struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	TierId uuid.UUID `json:"tier_id" db:"tier_id"`
	BenefitCode string `json:"benefit_code" db:"benefit_code"`
	BenefitName string `json:"benefit_name" db:"benefit_name"`
	BenefitDescription *string `json:"benefit_description" db:"benefit_description"`
	BenefitType string `json:"benefit_type" db:"benefit_type"`
	DiscountValue *float64 `json:"discount_value" db:"discount_value"`
	DiscountType *string `json:"discount_type" db:"discount_type"`
	IsActive *bool `json:"is_active" db:"is_active"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	Icon *string `json:"icon" db:"icon"`
	TermsAndConditions *string `json:"terms_and_conditions" db:"terms_and_conditions"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new loyalty_tier_benefits record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *LoyaltyTierBenefits) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "loyalty_tier_benefits", duration, nil)
	}()

	query := `
		INSERT INTO loyalty_tier_benefits (
			, organization_id
			, tier_id
			, benefit_code
			, benefit_name
			, benefit_description
			, benefit_type
			, discount_value
			, discount_type
			, is_active
			, sort_order
			, icon
			, terms_and_conditions
			, metadata
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.TierId,
		entity.BenefitCode,
		entity.BenefitName,
		entity.BenefitDescription,
		entity.BenefitType,
		entity.DiscountValue,
		entity.DiscountType,
		entity.IsActive,
		entity.SortOrder,
		entity.Icon,
		entity.TermsAndConditions,
		entity.Metadata,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create loyalty_tier_benefits", zap.Error(err))
		return fmt.Errorf("failed to create loyalty_tier_benefits: %w", err)
	}

	r.logger.Info("created loyalty_tier_benefits",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a loyalty_tier_benefits by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*LoyaltyTierBenefits, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_tier_benefits", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, tier_id
			, benefit_code
			, benefit_name
			, benefit_description
			, benefit_type
			, discount_value
			, discount_type
			, is_active
			, sort_order
			, icon
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
		FROM loyalty_tier_benefits
		WHERE id = $1
		
	`

	var entity LoyaltyTierBenefits
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.TierId,
		&entity.BenefitCode,
		&entity.BenefitName,
		&entity.BenefitDescription,
		&entity.BenefitType,
		&entity.DiscountValue,
		&entity.DiscountType,
		&entity.IsActive,
		&entity.SortOrder,
		&entity.Icon,
		&entity.TermsAndConditions,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("loyalty_tier_benefits not found")
	}

	if err != nil {
		r.logger.Error("failed to get loyalty_tier_benefits", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get loyalty_tier_benefits: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of loyalty_tier_benefits records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*LoyaltyTierBenefits, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_tier_benefits", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_tier_benefits
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_tier_benefits records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, tier_id
			, benefit_code
			, benefit_name
			, benefit_description
			, benefit_type
			, discount_value
			, discount_type
			, is_active
			, sort_order
			, icon
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
		FROM loyalty_tier_benefits
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_tier_benefits", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_tier_benefits: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyTierBenefits
	for rows.Next() {
		var entity LoyaltyTierBenefits
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TierId,
			&entity.BenefitCode,
			&entity.BenefitName,
			&entity.BenefitDescription,
			&entity.BenefitType,
			&entity.DiscountValue,
			&entity.DiscountType,
			&entity.IsActive,
			&entity.SortOrder,
			&entity.Icon,
			&entity.TermsAndConditions,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_tier_benefits: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating loyalty_tier_benefits rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing loyalty_tier_benefits record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *LoyaltyTierBenefits) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_tier_benefits", duration, nil)
	}()

	query := `
		UPDATE loyalty_tier_benefits
		SET
			, organization_id = $2
			, tier_id = $3
			, benefit_code = $4
			, benefit_name = $5
			, benefit_description = $6
			, benefit_type = $7
			, discount_value = $8
			, discount_type = $9
			, is_active = $10
			, sort_order = $11
			, icon = $12
			, terms_and_conditions = $13
			, metadata = $14
			, updated_at = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $17
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.TierId,
		entity.BenefitCode,
		entity.BenefitName,
		entity.BenefitDescription,
		entity.BenefitType,
		entity.DiscountValue,
		entity.DiscountType,
		entity.IsActive,
		entity.SortOrder,
		entity.Icon,
		entity.TermsAndConditions,
		entity.Metadata,
		entity.UpdatedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update loyalty_tier_benefits", zap.Error(err))
		return fmt.Errorf("failed to update loyalty_tier_benefits: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_tier_benefits not found or already deleted")
	}

	r.logger.Info("updated loyalty_tier_benefits",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a loyalty_tier_benefits record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "loyalty_tier_benefits", duration, nil)
	}()

	query := `DELETE FROM loyalty_tier_benefits WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete loyalty_tier_benefits", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete loyalty_tier_benefits: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_tier_benefits not found")
	}

	r.logger.Info("deleted loyalty_tier_benefits", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves loyalty_tier_benefits records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*LoyaltyTierBenefits, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_tier_benefits", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_tier_benefits
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_tier_benefits records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, tier_id
			, benefit_code
			, benefit_name
			, benefit_description
			, benefit_type
			, discount_value
			, discount_type
			, is_active
			, sort_order
			, icon
			, terms_and_conditions
			, metadata
			, created_at
			, updated_at
		FROM loyalty_tier_benefits
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_tier_benefits by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_tier_benefits: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyTierBenefits
	for rows.Next() {
		var entity LoyaltyTierBenefits
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TierId,
			&entity.BenefitCode,
			&entity.BenefitName,
			&entity.BenefitDescription,
			&entity.BenefitType,
			&entity.DiscountValue,
			&entity.DiscountType,
			&entity.IsActive,
			&entity.SortOrder,
			&entity.Icon,
			&entity.TermsAndConditions,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_tier_benefits: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


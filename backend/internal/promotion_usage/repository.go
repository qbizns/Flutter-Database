package promotion_usage

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

// Repository handles database operations for PromotionUsage
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PromotionUsage repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PromotionUsage represents a promotion_usage entity
type PromotionUsage struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PromotionId uuid.UUID `json:"promotion_id" db:"promotion_id"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	CustomerId *uuid.UUID `json:"customer_id" db:"customer_id"`
	DiscountAmount float64 `json:"discount_amount" db:"discount_amount"`
	UsedAt *time.Time `json:"used_at" db:"used_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new promotion_usage record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PromotionUsage) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "promotion_usage", duration, nil)
	}()

	query := `
		INSERT INTO promotion_usage (
			, organization_id
			, promotion_id
			, sale_id
			, customer_id
			, discount_amount
			, used_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PromotionId,
		entity.SaleId,
		entity.CustomerId,
		entity.DiscountAmount,
		entity.UsedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create promotion_usage", zap.Error(err))
		return fmt.Errorf("failed to create promotion_usage: %w", err)
	}

	r.logger.Info("created promotion_usage",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a promotion_usage by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PromotionUsage, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "promotion_usage", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, promotion_id
			, sale_id
			, customer_id
			, discount_amount
			, used_at
			, created_at
		FROM promotion_usage
		WHERE id = $1
		
	`

	var entity PromotionUsage
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PromotionId,
		&entity.SaleId,
		&entity.CustomerId,
		&entity.DiscountAmount,
		&entity.UsedAt,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("promotion_usage not found")
	}

	if err != nil {
		r.logger.Error("failed to get promotion_usage", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get promotion_usage: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of promotion_usage records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PromotionUsage, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "promotion_usage", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM promotion_usage
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count promotion_usage records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, promotion_id
			, sale_id
			, customer_id
			, discount_amount
			, used_at
			, created_at
		FROM promotion_usage
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list promotion_usage", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list promotion_usage: %w", err)
	}
	defer rows.Close()

	var entities []*PromotionUsage
	for rows.Next() {
		var entity PromotionUsage
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PromotionId,
			&entity.SaleId,
			&entity.CustomerId,
			&entity.DiscountAmount,
			&entity.UsedAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan promotion_usage: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating promotion_usage rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing promotion_usage record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PromotionUsage) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "promotion_usage", duration, nil)
	}()

	query := `
		UPDATE promotion_usage
		SET
			, organization_id = $2
			, promotion_id = $3
			, sale_id = $4
			, customer_id = $5
			, discount_amount = $6
			, used_at = $7
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PromotionId,
		entity.SaleId,
		entity.CustomerId,
		entity.DiscountAmount,
		entity.UsedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update promotion_usage", zap.Error(err))
		return fmt.Errorf("failed to update promotion_usage: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("promotion_usage not found or already deleted")
	}

	r.logger.Info("updated promotion_usage",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a promotion_usage record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "promotion_usage", duration, nil)
	}()

	query := `DELETE FROM promotion_usage WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete promotion_usage", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete promotion_usage: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("promotion_usage not found")
	}

	r.logger.Info("deleted promotion_usage", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves promotion_usage records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PromotionUsage, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "promotion_usage", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM promotion_usage
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count promotion_usage records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, promotion_id
			, sale_id
			, customer_id
			, discount_amount
			, used_at
			, created_at
		FROM promotion_usage
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list promotion_usage by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list promotion_usage: %w", err)
	}
	defer rows.Close()

	var entities []*PromotionUsage
	for rows.Next() {
		var entity PromotionUsage
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PromotionId,
			&entity.SaleId,
			&entity.CustomerId,
			&entity.DiscountAmount,
			&entity.UsedAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan promotion_usage: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


package loyalty_redemption

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

// Repository handles database operations for LoyaltyRedemptions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new LoyaltyRedemptions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// LoyaltyRedemptions represents a loyalty_redemptions entity
type LoyaltyRedemptions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CustomerId uuid.UUID `json:"customer_id" db:"customer_id"`
	RewardId uuid.UUID `json:"reward_id" db:"reward_id"`
	RedemptionNumber string `json:"redemption_number" db:"redemption_number"`
	RedemptionDate *time.Time `json:"redemption_date" db:"redemption_date"`
	PointsRedeemed int64 `json:"points_redeemed" db:"points_redeemed"`
	Status *string `json:"status" db:"status"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	UsedDate *time.Time `json:"used_date" db:"used_date"`
	ExpiryDate *time.Time `json:"expiry_date" db:"expiry_date"`
	// 	FulfillmentStatus *string `json:"fulfillment_status" db:"fulfillment_status"`
	FulfillmentNotes *string `json:"fulfillment_notes" db:"fulfillment_notes"`
	FulfilledBy *uuid.UUID `json:"fulfilled_by" db:"fulfilled_by"`
	FulfilledAt *time.Time `json:"fulfilled_at" db:"fulfilled_at"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new loyalty_redemptions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *LoyaltyRedemptions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "loyalty_redemptions", duration, nil)
	}()

	query := `
		INSERT INTO loyalty_redemptions (
			, organization_id
			, customer_id
			, reward_id
			, redemption_number
			, redemption_date
			, points_redeemed
			, status
			, sale_id
			, used_date
			, expiry_date
			, fulfillment_status
			, fulfillment_notes
			, fulfilled_by
			, fulfilled_at
			, notes
			, metadata
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
			, $20
			, $21
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.RewardId,
		entity.RedemptionNumber,
		entity.RedemptionDate,
		entity.PointsRedeemed,
		entity.Status,
		entity.SaleId,
		entity.UsedDate,
		entity.ExpiryDate,
		entity.FulfillmentStatus,
		entity.FulfillmentNotes,
		entity.FulfilledBy,
		entity.FulfilledAt,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create loyalty_redemptions", zap.Error(err))
		return fmt.Errorf("failed to create loyalty_redemptions: %w", err)
	}

	r.logger.Info("created loyalty_redemptions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a loyalty_redemptions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*LoyaltyRedemptions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_redemptions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, reward_id
			, redemption_number
			, redemption_date
			, points_redeemed
			, sale_id
			, used_date
			, expiry_date
			, fulfillment_status
			, fulfillment_notes
			, fulfilled_by
			, fulfilled_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
		FROM loyalty_redemptions
		WHERE id = $1
		
	`

	var entity LoyaltyRedemptions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CustomerId,
		&entity.RewardId,
		&entity.RedemptionNumber,
		&entity.RedemptionDate,
		&entity.PointsRedeemed,
		&entity.Status,
		&entity.SaleId,
		&entity.UsedDate,
		&entity.ExpiryDate,
		&entity.FulfillmentStatus,
		&entity.FulfillmentNotes,
		&entity.FulfilledBy,
		&entity.FulfilledAt,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("loyalty_redemptions not found")
	}

	if err != nil {
		r.logger.Error("failed to get loyalty_redemptions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get loyalty_redemptions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of loyalty_redemptions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*LoyaltyRedemptions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_redemptions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_redemptions
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_redemptions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, reward_id
			, redemption_number
			, redemption_date
			, points_redeemed
			, sale_id
			, used_date
			, expiry_date
			, fulfillment_status
			, fulfillment_notes
			, fulfilled_by
			, fulfilled_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
		FROM loyalty_redemptions
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_redemptions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_redemptions: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyRedemptions
	for rows.Next() {
		var entity LoyaltyRedemptions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.RewardId,
			&entity.RedemptionNumber,
			&entity.RedemptionDate,
			&entity.PointsRedeemed,
			&entity.Status,
			&entity.SaleId,
			&entity.UsedDate,
			&entity.ExpiryDate,
			&entity.FulfillmentStatus,
			&entity.FulfillmentNotes,
			&entity.FulfilledBy,
			&entity.FulfilledAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_redemptions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating loyalty_redemptions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing loyalty_redemptions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *LoyaltyRedemptions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_redemptions", duration, nil)
	}()

	query := `
		UPDATE loyalty_redemptions
		SET
			, organization_id = $2
			, customer_id = $3
			, reward_id = $4
			, redemption_number = $5
			, redemption_date = $6
			, points_redeemed = $7
			, status = $8
			, sale_id = $9
			, used_date = $10
			, expiry_date = $11
			, fulfillment_status = $12
			, fulfillment_notes = $13
			, fulfilled_by = $14
			, fulfilled_at = $15
			, notes = $16
			, metadata = $17
			, updated_at = $19
			, created_by = $20
			, updated_by = $21
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $22
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.RewardId,
		entity.RedemptionNumber,
		entity.RedemptionDate,
		entity.PointsRedeemed,
		entity.Status,
		entity.SaleId,
		entity.UsedDate,
		entity.ExpiryDate,
		entity.FulfillmentStatus,
		entity.FulfillmentNotes,
		entity.FulfilledBy,
		entity.FulfilledAt,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update loyalty_redemptions", zap.Error(err))
		return fmt.Errorf("failed to update loyalty_redemptions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_redemptions not found or already deleted")
	}

	r.logger.Info("updated loyalty_redemptions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a loyalty_redemptions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "loyalty_redemptions", duration, nil)
	}()

	query := `DELETE FROM loyalty_redemptions WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete loyalty_redemptions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete loyalty_redemptions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_redemptions not found")
	}

	r.logger.Info("deleted loyalty_redemptions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves loyalty_redemptions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*LoyaltyRedemptions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_redemptions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_redemptions
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_redemptions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, reward_id
			, redemption_number
			, redemption_date
			, points_redeemed
			, sale_id
			, used_date
			, expiry_date
			, fulfillment_status
			, fulfillment_notes
			, fulfilled_by
			, fulfilled_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
		FROM loyalty_redemptions
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_redemptions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_redemptions: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyRedemptions
	for rows.Next() {
		var entity LoyaltyRedemptions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.RewardId,
			&entity.RedemptionNumber,
			&entity.RedemptionDate,
			&entity.PointsRedeemed,
			&entity.Status,
			&entity.SaleId,
			&entity.UsedDate,
			&entity.ExpiryDate,
			&entity.FulfillmentStatus,
			&entity.FulfillmentNotes,
			&entity.FulfilledBy,
			&entity.FulfilledAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_redemptions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


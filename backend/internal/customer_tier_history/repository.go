package customer_tier_history

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

// Repository handles database operations for CustomerTierHistory
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CustomerTierHistory repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CustomerTierHistory represents a customer_tier_history entity
type CustomerTierHistory struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CustomerId uuid.UUID `json:"customer_id" db:"customer_id"`
	TierId uuid.UUID `json:"tier_id" db:"tier_id"`
	PreviousTierId *uuid.UUID `json:"previous_tier_id" db:"previous_tier_id"`
	ChangeType string `json:"change_type" db:"change_type"`
	ChangeReason *string `json:"change_reason" db:"change_reason"`
	QualifyingPoints *int64 `json:"qualifying_points" db:"qualifying_points"`
	QualifyingSpend *float64 `json:"qualifying_spend" db:"qualifying_spend"`
	QualifyingPurchases *int64 `json:"qualifying_purchases" db:"qualifying_purchases"`
	EffectiveDate time.Time `json:"effective_date" db:"effective_date"`
	ValidUntil *time.Time `json:"valid_until" db:"valid_until"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
}

// Create inserts a new customer_tier_history record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CustomerTierHistory) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "customer_tier_history", duration, nil)
	}()

	query := `
		INSERT INTO customer_tier_history (
			, organization_id
			, customer_id
			, tier_id
			, previous_tier_id
			, change_type
			, change_reason
			, qualifying_points
			, qualifying_spend
			, qualifying_purchases
			, effective_date
			, valid_until
			, notes
			, metadata
			, created_by
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
			, $16
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.TierId,
		entity.PreviousTierId,
		entity.ChangeType,
		entity.ChangeReason,
		entity.QualifyingPoints,
		entity.QualifyingSpend,
		entity.QualifyingPurchases,
		entity.EffectiveDate,
		entity.ValidUntil,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create customer_tier_history", zap.Error(err))
		return fmt.Errorf("failed to create customer_tier_history: %w", err)
	}

	r.logger.Info("created customer_tier_history",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a customer_tier_history by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CustomerTierHistory, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_tier_history", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, tier_id
			, previous_tier_id
			, change_type
			, change_reason
			, qualifying_points
			, qualifying_spend
			, qualifying_purchases
			, effective_date
			, valid_until
			, notes
			, metadata
			, created_at
			, created_by
		FROM customer_tier_history
		WHERE id = $1
		
	`

	var entity CustomerTierHistory
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CustomerId,
		&entity.TierId,
		&entity.PreviousTierId,
		&entity.ChangeType,
		&entity.ChangeReason,
		&entity.QualifyingPoints,
		&entity.QualifyingSpend,
		&entity.QualifyingPurchases,
		&entity.EffectiveDate,
		&entity.ValidUntil,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("customer_tier_history not found")
	}

	if err != nil {
		r.logger.Error("failed to get customer_tier_history", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get customer_tier_history: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of customer_tier_history records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CustomerTierHistory, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_tier_history", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_tier_history
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_tier_history records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, tier_id
			, previous_tier_id
			, change_type
			, change_reason
			, qualifying_points
			, qualifying_spend
			, qualifying_purchases
			, effective_date
			, valid_until
			, notes
			, metadata
			, created_at
			, created_by
		FROM customer_tier_history
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_tier_history", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_tier_history: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerTierHistory
	for rows.Next() {
		var entity CustomerTierHistory
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.TierId,
			&entity.PreviousTierId,
			&entity.ChangeType,
			&entity.ChangeReason,
			&entity.QualifyingPoints,
			&entity.QualifyingSpend,
			&entity.QualifyingPurchases,
			&entity.EffectiveDate,
			&entity.ValidUntil,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_tier_history: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating customer_tier_history rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing customer_tier_history record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CustomerTierHistory) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customer_tier_history", duration, nil)
	}()

	query := `
		UPDATE customer_tier_history
		SET
			, organization_id = $2
			, customer_id = $3
			, tier_id = $4
			, previous_tier_id = $5
			, change_type = $6
			, change_reason = $7
			, qualifying_points = $8
			, qualifying_spend = $9
			, qualifying_purchases = $10
			, effective_date = $11
			, valid_until = $12
			, notes = $13
			, metadata = $14
			, created_by = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $17
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.TierId,
		entity.PreviousTierId,
		entity.ChangeType,
		entity.ChangeReason,
		entity.QualifyingPoints,
		entity.QualifyingSpend,
		entity.QualifyingPurchases,
		entity.EffectiveDate,
		entity.ValidUntil,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update customer_tier_history", zap.Error(err))
		return fmt.Errorf("failed to update customer_tier_history: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_tier_history not found or already deleted")
	}

	r.logger.Info("updated customer_tier_history",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a customer_tier_history record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "customer_tier_history", duration, nil)
	}()

	query := `DELETE FROM customer_tier_history WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete customer_tier_history", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete customer_tier_history: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer_tier_history not found")
	}

	r.logger.Info("deleted customer_tier_history", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves customer_tier_history records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CustomerTierHistory, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customer_tier_history", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customer_tier_history
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer_tier_history records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, tier_id
			, previous_tier_id
			, change_type
			, change_reason
			, qualifying_points
			, qualifying_spend
			, qualifying_purchases
			, effective_date
			, valid_until
			, notes
			, metadata
			, created_at
			, created_by
		FROM customer_tier_history
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customer_tier_history by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customer_tier_history: %w", err)
	}
	defer rows.Close()

	var entities []*CustomerTierHistory
	for rows.Next() {
		var entity CustomerTierHistory
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.TierId,
			&entity.PreviousTierId,
			&entity.ChangeType,
			&entity.ChangeReason,
			&entity.QualifyingPoints,
			&entity.QualifyingSpend,
			&entity.QualifyingPurchases,
			&entity.EffectiveDate,
			&entity.ValidUntil,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer_tier_history: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


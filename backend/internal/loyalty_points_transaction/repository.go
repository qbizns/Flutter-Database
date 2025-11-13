package loyalty_points_transaction

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

// Repository handles database operations for LoyaltyPointsTransactions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new LoyaltyPointsTransactions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// LoyaltyPointsTransactions represents a loyalty_points_transactions entity
type LoyaltyPointsTransactions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CustomerId uuid.UUID `json:"customer_id" db:"customer_id"`
	TransactionType string `json:"transaction_type" db:"transaction_type"`
	Points int64 `json:"points" db:"points"`
	BalanceAfter int64 `json:"balance_after" db:"balance_after"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	RedemptionId *uuid.UUID `json:"redemption_id" db:"redemption_id"`
	PointsRuleId *uuid.UUID `json:"points_rule_id" db:"points_rule_id"`
	Description *string `json:"description" db:"description"`
	Reason *string `json:"reason" db:"reason"`
	Notes *string `json:"notes" db:"notes"`
	ExpiryDate *time.Time `json:"expiry_date" db:"expiry_date"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	TransactionDate *time.Time `json:"transaction_date" db:"transaction_date"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
}

// Create inserts a new loyalty_points_transactions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *LoyaltyPointsTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "loyalty_points_transactions", duration, nil)
	}()

	query := `
		INSERT INTO loyalty_points_transactions (
			, organization_id
			, customer_id
			, transaction_type
			, points
			, balance_after
			, sale_id
			, redemption_id
			, points_rule_id
			, description
			, reason
			, notes
			, expiry_date
			, metadata
			, transaction_date
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
			, $15
			, $16
		)
		RETURNING id
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.TransactionType,
		entity.Points,
		entity.BalanceAfter,
		entity.SaleId,
		entity.RedemptionId,
		entity.PointsRuleId,
		entity.Description,
		entity.Reason,
		entity.Notes,
		entity.ExpiryDate,
		entity.Metadata,
		entity.TransactionDate,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id)
	

	if err != nil {
		r.logger.Error("failed to create loyalty_points_transactions", zap.Error(err))
		return fmt.Errorf("failed to create loyalty_points_transactions: %w", err)
	}

	r.logger.Info("created loyalty_points_transactions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a loyalty_points_transactions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*LoyaltyPointsTransactions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_points_transactions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, transaction_type
			, points
			, balance_after
			, sale_id
			, redemption_id
			, points_rule_id
			, description
			, reason
			, notes
			, expiry_date
			, metadata
			, transaction_date
			, created_by
		FROM loyalty_points_transactions
		WHERE id = $1
		
	`

	var entity LoyaltyPointsTransactions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CustomerId,
		&entity.TransactionType,
		&entity.Points,
		&entity.BalanceAfter,
		&entity.SaleId,
		&entity.RedemptionId,
		&entity.PointsRuleId,
		&entity.Description,
		&entity.Reason,
		&entity.Notes,
		&entity.ExpiryDate,
		&entity.Metadata,
		&entity.TransactionDate,
		&entity.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("loyalty_points_transactions not found")
	}

	if err != nil {
		r.logger.Error("failed to get loyalty_points_transactions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get loyalty_points_transactions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of loyalty_points_transactions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*LoyaltyPointsTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_points_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_points_transactions
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_points_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, transaction_type
			, points
			, balance_after
			, sale_id
			, redemption_id
			, points_rule_id
			, description
			, reason
			, notes
			, expiry_date
			, metadata
			, transaction_date
			, created_by
		FROM loyalty_points_transactions
		
		
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_points_transactions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_points_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyPointsTransactions
	for rows.Next() {
		var entity LoyaltyPointsTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.TransactionType,
			&entity.Points,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.RedemptionId,
			&entity.PointsRuleId,
			&entity.Description,
			&entity.Reason,
			&entity.Notes,
			&entity.ExpiryDate,
			&entity.Metadata,
			&entity.TransactionDate,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_points_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating loyalty_points_transactions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing loyalty_points_transactions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *LoyaltyPointsTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "loyalty_points_transactions", duration, nil)
	}()

	query := `
		UPDATE loyalty_points_transactions
		SET
			, organization_id = $2
			, customer_id = $3
			, transaction_type = $4
			, points = $5
			, balance_after = $6
			, sale_id = $7
			, redemption_id = $8
			, points_rule_id = $9
			, description = $10
			, reason = $11
			, notes = $12
			, expiry_date = $13
			, metadata = $14
			, transaction_date = $15
			, created_by = $16
			
		WHERE id = $17
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CustomerId,
		entity.TransactionType,
		entity.Points,
		entity.BalanceAfter,
		entity.SaleId,
		entity.RedemptionId,
		entity.PointsRuleId,
		entity.Description,
		entity.Reason,
		entity.Notes,
		entity.ExpiryDate,
		entity.Metadata,
		entity.TransactionDate,
		entity.CreatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update loyalty_points_transactions", zap.Error(err))
		return fmt.Errorf("failed to update loyalty_points_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_points_transactions not found or already deleted")
	}

	r.logger.Info("updated loyalty_points_transactions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a loyalty_points_transactions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "loyalty_points_transactions", duration, nil)
	}()

	query := `DELETE FROM loyalty_points_transactions WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete loyalty_points_transactions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete loyalty_points_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("loyalty_points_transactions not found")
	}

	r.logger.Info("deleted loyalty_points_transactions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves loyalty_points_transactions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*LoyaltyPointsTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "loyalty_points_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM loyalty_points_transactions
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count loyalty_points_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_id
			, transaction_type
			, points
			, balance_after
			, sale_id
			, redemption_id
			, points_rule_id
			, description
			, reason
			, notes
			, expiry_date
			, metadata
			, transaction_date
			, created_by
		FROM loyalty_points_transactions
		WHERE organization_id = $1
		
		
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list loyalty_points_transactions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list loyalty_points_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*LoyaltyPointsTransactions
	for rows.Next() {
		var entity LoyaltyPointsTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerId,
			&entity.TransactionType,
			&entity.Points,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.RedemptionId,
			&entity.PointsRuleId,
			&entity.Description,
			&entity.Reason,
			&entity.Notes,
			&entity.ExpiryDate,
			&entity.Metadata,
			&entity.TransactionDate,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan loyalty_points_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


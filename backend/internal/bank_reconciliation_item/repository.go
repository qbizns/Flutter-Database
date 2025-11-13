package bank_reconciliation_item

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

// Repository handles database operations for BankReconciliationItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BankReconciliationItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BankReconciliationItems represents a bank_reconciliation_items entity
type BankReconciliationItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	BankReconciliationId *uuid.UUID `json:"bank_reconciliation_id" db:"bank_reconciliation_id"`
	GeneralLedgerId uuid.UUID `json:"general_ledger_id" db:"general_ledger_id"`
	JournalEntryLineId uuid.UUID `json:"journal_entry_line_id" db:"journal_entry_line_id"`
	IsCleared *bool `json:"is_cleared" db:"is_cleared"`
	ClearedDate *time.Time `json:"cleared_date" db:"cleared_date"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	ClearedBy *uuid.UUID `json:"cleared_by" db:"cleared_by"`
}

// Create inserts a new bank_reconciliation_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BankReconciliationItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "bank_reconciliation_items", duration, nil)
	}()

	query := `
		INSERT INTO bank_reconciliation_items (
			, organization_id
			, bank_reconciliation_id
			, general_ledger_id
			, journal_entry_line_id
			, is_cleared
			, cleared_date
			, cleared_by
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $9
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.BankReconciliationId,
		entity.GeneralLedgerId,
		entity.JournalEntryLineId,
		entity.IsCleared,
		entity.ClearedDate,
		entity.ClearedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create bank_reconciliation_items", zap.Error(err))
		return fmt.Errorf("failed to create bank_reconciliation_items: %w", err)
	}

	r.logger.Info("created bank_reconciliation_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a bank_reconciliation_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BankReconciliationItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_reconciliation_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, bank_reconciliation_id
			, general_ledger_id
			, journal_entry_line_id
			, is_cleared
			, cleared_date
			, created_at
			, cleared_by
		FROM bank_reconciliation_items
		WHERE id = $1
		
	`

	var entity BankReconciliationItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.BankReconciliationId,
		&entity.GeneralLedgerId,
		&entity.JournalEntryLineId,
		&entity.IsCleared,
		&entity.ClearedDate,
		&entity.CreatedAt,
		&entity.ClearedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bank_reconciliation_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get bank_reconciliation_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get bank_reconciliation_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of bank_reconciliation_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BankReconciliationItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_reconciliation_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_reconciliation_items
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_reconciliation_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_reconciliation_id
			, general_ledger_id
			, journal_entry_line_id
			, is_cleared
			, cleared_date
			, created_at
			, cleared_by
		FROM bank_reconciliation_items
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_reconciliation_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_reconciliation_items: %w", err)
	}
	defer rows.Close()

	var entities []*BankReconciliationItems
	for rows.Next() {
		var entity BankReconciliationItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankReconciliationId,
			&entity.GeneralLedgerId,
			&entity.JournalEntryLineId,
			&entity.IsCleared,
			&entity.ClearedDate,
			&entity.CreatedAt,
			&entity.ClearedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_reconciliation_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating bank_reconciliation_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing bank_reconciliation_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BankReconciliationItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_reconciliation_items", duration, nil)
	}()

	query := `
		UPDATE bank_reconciliation_items
		SET
			, organization_id = $2
			, bank_reconciliation_id = $3
			, general_ledger_id = $4
			, journal_entry_line_id = $5
			, is_cleared = $6
			, cleared_date = $7
			, cleared_by = $9
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.BankReconciliationId,
		entity.GeneralLedgerId,
		entity.JournalEntryLineId,
		entity.IsCleared,
		entity.ClearedDate,
		entity.ClearedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update bank_reconciliation_items", zap.Error(err))
		return fmt.Errorf("failed to update bank_reconciliation_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_reconciliation_items not found or already deleted")
	}

	r.logger.Info("updated bank_reconciliation_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a bank_reconciliation_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "bank_reconciliation_items", duration, nil)
	}()

	query := `DELETE FROM bank_reconciliation_items WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete bank_reconciliation_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete bank_reconciliation_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_reconciliation_items not found")
	}

	r.logger.Info("deleted bank_reconciliation_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves bank_reconciliation_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*BankReconciliationItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_reconciliation_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_reconciliation_items
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_reconciliation_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_reconciliation_id
			, general_ledger_id
			, journal_entry_line_id
			, is_cleared
			, cleared_date
			, created_at
			, cleared_by
		FROM bank_reconciliation_items
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_reconciliation_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_reconciliation_items: %w", err)
	}
	defer rows.Close()

	var entities []*BankReconciliationItems
	for rows.Next() {
		var entity BankReconciliationItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankReconciliationId,
			&entity.GeneralLedgerId,
			&entity.JournalEntryLineId,
			&entity.IsCleared,
			&entity.ClearedDate,
			&entity.CreatedAt,
			&entity.ClearedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_reconciliation_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


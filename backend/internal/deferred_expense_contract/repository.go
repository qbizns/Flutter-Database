package deferred_expense_contract

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

// Repository handles database operations for DeferredExpenseContracts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DeferredExpenseContracts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DeferredExpenseContracts represents a deferred_expense_contracts entity
type DeferredExpenseContracts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	VendorBillId *uuid.UUID `json:"vendor_bill_id" db:"vendor_bill_id"`
	BillLineId *uuid.UUID `json:"bill_line_id" db:"bill_line_id"`
	ContractName *string `json:"contract_name" db:"contract_name"`
	TotalDeferredAmount float64 `json:"total_deferred_amount" db:"total_deferred_amount"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate time.Time `json:"end_date" db:"end_date"`
	RecognitionMethod *string `json:"recognition_method" db:"recognition_method"`
	RecognitionMethod *string `json:"recognition_method" db:"recognition_method"`
	DeferredAccountId uuid.UUID `json:"deferred_account_id" db:"deferred_account_id"`
	ExpenseAccountId uuid.UUID `json:"expense_account_id" db:"expense_account_id"`
	Status *string `json:"status" db:"status"`
	RecognizedAmount *float64 `json:"recognized_amount" db:"recognized_amount"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new deferred_expense_contracts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DeferredExpenseContracts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "deferred_expense_contracts", duration, nil)
	}()

	query := `
		INSERT INTO deferred_expense_contracts (
			, organization_id
			, vendor_bill_id
			, bill_line_id
			, contract_name
			, total_deferred_amount
			, start_date
			, end_date
			, recognition_method
			, recognition_method
			, deferred_account_id
			, expense_account_id
			, status
			, recognized_amount
			, notes
			, created_by
			, updated_by
			, deleted_at
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.VendorBillId,
		entity.BillLineId,
		entity.ContractName,
		entity.TotalDeferredAmount,
		entity.StartDate,
		entity.EndDate,
		entity.RecognitionMethod,
		entity.RecognitionMethod,
		entity.DeferredAccountId,
		entity.ExpenseAccountId,
		entity.Status,
		entity.RecognizedAmount,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create deferred_expense_contracts", zap.Error(err))
		return fmt.Errorf("failed to create deferred_expense_contracts: %w", err)
	}

	r.logger.Info("created deferred_expense_contracts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a deferred_expense_contracts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DeferredExpenseContracts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "deferred_expense_contracts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, vendor_bill_id
			, bill_line_id
			, contract_name
			, total_deferred_amount
			, start_date
			, end_date
			, recognition_method
			, recognition_method
			, deferred_account_id
			, expense_account_id
			, status
			, recognized_amount
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM deferred_expense_contracts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity DeferredExpenseContracts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.VendorBillId,
		&entity.BillLineId,
		&entity.ContractName,
		&entity.TotalDeferredAmount,
		&entity.StartDate,
		&entity.EndDate,
		&entity.RecognitionMethod,
		&entity.RecognitionMethod,
		&entity.DeferredAccountId,
		&entity.ExpenseAccountId,
		&entity.Status,
		&entity.RecognizedAmount,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("deferred_expense_contracts not found")
	}

	if err != nil {
		r.logger.Error("failed to get deferred_expense_contracts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get deferred_expense_contracts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of deferred_expense_contracts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DeferredExpenseContracts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "deferred_expense_contracts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM deferred_expense_contracts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count deferred_expense_contracts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, vendor_bill_id
			, bill_line_id
			, contract_name
			, total_deferred_amount
			, start_date
			, end_date
			, recognition_method
			, recognition_method
			, deferred_account_id
			, expense_account_id
			, status
			, recognized_amount
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM deferred_expense_contracts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list deferred_expense_contracts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list deferred_expense_contracts: %w", err)
	}
	defer rows.Close()

	var entities []*DeferredExpenseContracts
	for rows.Next() {
		var entity DeferredExpenseContracts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.VendorBillId,
			&entity.BillLineId,
			&entity.ContractName,
			&entity.TotalDeferredAmount,
			&entity.StartDate,
			&entity.EndDate,
			&entity.RecognitionMethod,
			&entity.RecognitionMethod,
			&entity.DeferredAccountId,
			&entity.ExpenseAccountId,
			&entity.Status,
			&entity.RecognizedAmount,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan deferred_expense_contracts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating deferred_expense_contracts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing deferred_expense_contracts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DeferredExpenseContracts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "deferred_expense_contracts", duration, nil)
	}()

	query := `
		UPDATE deferred_expense_contracts
		SET
			, organization_id = $2
			, vendor_bill_id = $3
			, bill_line_id = $4
			, contract_name = $5
			, total_deferred_amount = $6
			, start_date = $7
			, end_date = $8
			, recognition_method = $9
			, recognition_method = $10
			, deferred_account_id = $11
			, expense_account_id = $12
			, status = $13
			, recognized_amount = $14
			, notes = $15
			, created_by = $16
			, updated_by = $17
			, updated_at = $19
			, deleted_at = $20
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.VendorBillId,
		entity.BillLineId,
		entity.ContractName,
		entity.TotalDeferredAmount,
		entity.StartDate,
		entity.EndDate,
		entity.RecognitionMethod,
		entity.RecognitionMethod,
		entity.DeferredAccountId,
		entity.ExpenseAccountId,
		entity.Status,
		entity.RecognizedAmount,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update deferred_expense_contracts", zap.Error(err))
		return fmt.Errorf("failed to update deferred_expense_contracts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("deferred_expense_contracts not found or already deleted")
	}

	r.logger.Info("updated deferred_expense_contracts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a deferred_expense_contracts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "deferred_expense_contracts", duration, nil)
	}()

	query := `
		UPDATE deferred_expense_contracts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete deferred_expense_contracts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete deferred_expense_contracts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("deferred_expense_contracts not found or already deleted")
	}

	r.logger.Info("deleted deferred_expense_contracts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves deferred_expense_contracts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*DeferredExpenseContracts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "deferred_expense_contracts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM deferred_expense_contracts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count deferred_expense_contracts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, vendor_bill_id
			, bill_line_id
			, contract_name
			, total_deferred_amount
			, start_date
			, end_date
			, recognition_method
			, recognition_method
			, deferred_account_id
			, expense_account_id
			, status
			, recognized_amount
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM deferred_expense_contracts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list deferred_expense_contracts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list deferred_expense_contracts: %w", err)
	}
	defer rows.Close()

	var entities []*DeferredExpenseContracts
	for rows.Next() {
		var entity DeferredExpenseContracts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.VendorBillId,
			&entity.BillLineId,
			&entity.ContractName,
			&entity.TotalDeferredAmount,
			&entity.StartDate,
			&entity.EndDate,
			&entity.RecognitionMethod,
			&entity.RecognitionMethod,
			&entity.DeferredAccountId,
			&entity.ExpenseAccountId,
			&entity.Status,
			&entity.RecognizedAmount,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan deferred_expense_contracts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


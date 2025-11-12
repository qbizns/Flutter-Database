package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"internal/domain/analytics"
)

// AnalyticsRepository implements analytics.Repository for PostgreSQL
type AnalyticsRepository struct {
	pool *pgxpool.Pool
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{pool: pool}
}

// ==================== ANALYTIC PLANS ====================

// CreateAnalyticPlan creates a new analytic plan
func (r *AnalyticsRepository) CreateAnalyticPlan(ctx context.Context, plan *analytics.AnalyticPlan) error {
	query := `
		INSERT INTO analytic_plans (
			id, organization_id, plan_code, plan_name, is_active,
			description, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		plan.ID, plan.OrganizationID, plan.PlanCode, plan.PlanName,
		plan.IsActive, plan.Description, plan.CreatedBy,
		plan.CreatedAt, plan.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return analytics.ErrDuplicateAnalyticPlanCode
		}
		return fmt.Errorf("failed to create analytic plan: %w", err)
	}

	return nil
}

// GetAnalyticPlan retrieves an analytic plan
func (r *AnalyticsRepository) GetAnalyticPlan(ctx context.Context, id uuid.UUID) (*analytics.AnalyticPlan, error) {
	query := `
		SELECT id, organization_id, plan_code, plan_name, is_active,
		       description, created_by, created_at, updated_at, deleted_at
		FROM analytic_plans
		WHERE id = $1 AND deleted_at IS NULL
	`

	var plan analytics.AnalyticPlan
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&plan.ID, &plan.OrganizationID, &plan.PlanCode, &plan.PlanName,
		&plan.IsActive, &plan.Description, &plan.CreatedBy,
		&plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get analytic plan: %w", err)
	}

	return &plan, nil
}

// GetAnalyticPlanByCode retrieves an analytic plan by code
func (r *AnalyticsRepository) GetAnalyticPlanByCode(ctx context.Context, organizationID uuid.UUID, code string) (*analytics.AnalyticPlan, error) {
	query := `
		SELECT id, organization_id, plan_code, plan_name, is_active,
		       description, created_by, created_at, updated_at, deleted_at
		FROM analytic_plans
		WHERE organization_id = $1 AND plan_code = $2 AND deleted_at IS NULL
	`

	var plan analytics.AnalyticPlan
	err := r.pool.QueryRow(ctx, query, organizationID, code).Scan(
		&plan.ID, &plan.OrganizationID, &plan.PlanCode, &plan.PlanName,
		&plan.IsActive, &plan.Description, &plan.CreatedBy,
		&plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get analytic plan by code: %w", err)
	}

	return &plan, nil
}

// ListAnalyticPlans retrieves analytic plans
func (r *AnalyticsRepository) ListAnalyticPlans(ctx context.Context, query analytics.AnalyticQuery) ([]*analytics.AnalyticPlan, int64, error) {
	where := []string{"deleted_at IS NULL", "organization_id = $1"}
	args := []interface{}{query.OrganizationID}
	idx := 2

	if query.IsActive != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *query.IsActive)
		idx++
	}

	// Get count
	countQuery := "SELECT COUNT(*) FROM analytic_plans WHERE " + strings.Join(where, " AND ")
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count analytic plans: %w", err)
	}

	// Get records
	listQuery := `
		SELECT id, organization_id, plan_code, plan_name, is_active,
		       description, created_by, created_at, updated_at, deleted_at
		FROM analytic_plans
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY plan_code ASC LIMIT $` + fmt.Sprintf("%d", idx) +
		` OFFSET $` + fmt.Sprintf("%d", idx+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list analytic plans: %w", err)
	}
	defer rows.Close()

	var plans []*analytics.AnalyticPlan
	for rows.Next() {
		var plan analytics.AnalyticPlan
		err := rows.Scan(
			&plan.ID, &plan.OrganizationID, &plan.PlanCode, &plan.PlanName,
			&plan.IsActive, &plan.Description, &plan.CreatedBy,
			&plan.CreatedAt, &plan.UpdatedAt, &plan.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan analytic plan: %w", err)
		}
		plans = append(plans, &plan)
	}

	return plans, total, rows.Err()
}

// UpdateAnalyticPlan updates an analytic plan
func (r *AnalyticsRepository) UpdateAnalyticPlan(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE analytic_plans SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteAnalyticPlan soft deletes an analytic plan
func (r *AnalyticsRepository) DeleteAnalyticPlan(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE analytic_plans SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// ==================== ANALYTIC ACCOUNTS ====================

// CreateAnalyticAccount creates a new analytic account
func (r *AnalyticsRepository) CreateAnalyticAccount(ctx context.Context, account *analytics.AnalyticAccount) error {
	query := `
		INSERT INTO analytic_accounts (
			id, organization_id, analytic_plan_id, account_code, account_name,
			parent_account_id, account_level, is_active, description,
			created_by, updated_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.pool.Exec(ctx, query,
		account.ID, account.OrganizationID, account.AnalyticPlanID,
		account.AccountCode, account.AccountName, account.ParentAccountID,
		account.AccountLevel, account.IsActive, account.Description,
		account.CreatedBy, account.UpdatedBy, account.CreatedAt, account.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return analytics.ErrDuplicateAnalyticAccountCode
		}
		return fmt.Errorf("failed to create analytic account: %w", err)
	}

	return nil
}

// GetAnalyticAccount retrieves an analytic account
func (r *AnalyticsRepository) GetAnalyticAccount(ctx context.Context, id uuid.UUID) (*analytics.AnalyticAccount, error) {
	query := `
		SELECT id, organization_id, analytic_plan_id, account_code, account_name,
		       parent_account_id, account_level, is_active, description,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM analytic_accounts
		WHERE id = $1 AND deleted_at IS NULL
	`

	var account analytics.AnalyticAccount
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&account.ID, &account.OrganizationID, &account.AnalyticPlanID,
		&account.AccountCode, &account.AccountName, &account.ParentAccountID,
		&account.AccountLevel, &account.IsActive, &account.Description,
		&account.CreatedBy, &account.UpdatedBy, &account.CreatedAt,
		&account.UpdatedAt, &account.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get analytic account: %w", err)
	}

	return &account, nil
}

// GetAnalyticAccountByCode retrieves an analytic account by code
func (r *AnalyticsRepository) GetAnalyticAccountByCode(ctx context.Context, organizationID uuid.UUID, code string) (*analytics.AnalyticAccount, error) {
	query := `
		SELECT id, organization_id, analytic_plan_id, account_code, account_name,
		       parent_account_id, account_level, is_active, description,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM analytic_accounts
		WHERE organization_id = $1 AND account_code = $2 AND deleted_at IS NULL
	`

	var account analytics.AnalyticAccount
	err := r.pool.QueryRow(ctx, query, organizationID, code).Scan(
		&account.ID, &account.OrganizationID, &account.AnalyticPlanID,
		&account.AccountCode, &account.AccountName, &account.ParentAccountID,
		&account.AccountLevel, &account.IsActive, &account.Description,
		&account.CreatedBy, &account.UpdatedBy, &account.CreatedAt,
		&account.UpdatedAt, &account.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get analytic account by code: %w", err)
	}

	return &account, nil
}

// ListAnalyticAccounts retrieves analytic accounts
func (r *AnalyticsRepository) ListAnalyticAccounts(ctx context.Context, query analytics.AnalyticQuery) ([]*analytics.AnalyticAccount, int64, error) {
	where := []string{"deleted_at IS NULL", "organization_id = $1"}
	args := []interface{}{query.OrganizationID}
	idx := 2

	if query.AnalyticPlanID != nil {
		where = append(where, fmt.Sprintf("analytic_plan_id = $%d", idx))
		args = append(args, *query.AnalyticPlanID)
		idx++
	}

	if query.IsActive != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, *query.IsActive)
		idx++
	}

	if query.ParentAccountID != nil {
		where = append(where, fmt.Sprintf("parent_account_id = $%d", idx))
		args = append(args, *query.ParentAccountID)
		idx++
	}

	// Get count
	countQuery := "SELECT COUNT(*) FROM analytic_accounts WHERE " + strings.Join(where, " AND ")
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count analytic accounts: %w", err)
	}

	// Get records
	listQuery := `
		SELECT id, organization_id, analytic_plan_id, account_code, account_name,
		       parent_account_id, account_level, is_active, description,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM analytic_accounts
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY account_code ASC LIMIT $` + fmt.Sprintf("%d", idx) +
		` OFFSET $` + fmt.Sprintf("%d", idx+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list analytic accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*analytics.AnalyticAccount
	for rows.Next() {
		var account analytics.AnalyticAccount
		err := rows.Scan(
			&account.ID, &account.OrganizationID, &account.AnalyticPlanID,
			&account.AccountCode, &account.AccountName, &account.ParentAccountID,
			&account.AccountLevel, &account.IsActive, &account.Description,
			&account.CreatedBy, &account.UpdatedBy, &account.CreatedAt,
			&account.UpdatedAt, &account.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan analytic account: %w", err)
		}
		accounts = append(accounts, &account)
	}

	return accounts, total, rows.Err()
}

// UpdateAnalyticAccount updates an analytic account
func (r *AnalyticsRepository) UpdateAnalyticAccount(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE analytic_accounts SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteAnalyticAccount soft deletes an analytic account
func (r *AnalyticsRepository) DeleteAnalyticAccount(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE analytic_accounts SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// GetAnalyticAccountHierarchy retrieves the full hierarchy
func (r *AnalyticsRepository) GetAnalyticAccountHierarchy(ctx context.Context, organizationID uuid.UUID) ([]*analytics.AnalyticAccount, error) {
	query := `
		SELECT id, organization_id, analytic_plan_id, account_code, account_name,
		       parent_account_id, account_level, is_active, description,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM analytic_accounts
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY parent_account_id, account_code
	`

	rows, err := r.pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytic account hierarchy: %w", err)
	}
	defer rows.Close()

	var accounts []*analytics.AnalyticAccount
	for rows.Next() {
		var account analytics.AnalyticAccount
		err := rows.Scan(
			&account.ID, &account.OrganizationID, &account.AnalyticPlanID,
			&account.AccountCode, &account.AccountName, &account.ParentAccountID,
			&account.AccountLevel, &account.IsActive, &account.Description,
			&account.CreatedBy, &account.UpdatedBy, &account.CreatedAt,
			&account.UpdatedAt, &account.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analytic account: %w", err)
		}
		accounts = append(accounts, &account)
	}

	return accounts, rows.Err()
}

// ==================== DEFERRED REVENUE ====================

// CreateDeferredRevenueContract creates a deferred revenue contract
func (r *AnalyticsRepository) CreateDeferredRevenueContract(ctx context.Context, contract *analytics.DeferredRevenueContract) error {
	query := `
		INSERT INTO deferred_revenue_contracts (
			id, organization_id, customer_invoice_id, invoice_line_id,
			contract_name, total_deferred_amount, start_date, end_date,
			recognition_method, deferred_account_id, revenue_account_id,
			status, recognized_amount, notes, created_by, updated_by,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.pool.Exec(ctx, query,
		contract.ID, contract.OrganizationID, contract.CustomerInvoiceID,
		contract.InvoiceLineID, contract.ContractName, contract.TotalDeferredAmount,
		contract.StartDate, contract.EndDate, contract.RecognitionMethod,
		contract.DeferredAccountID, contract.RevenueAccountID, contract.Status,
		contract.RecognizedAmount, contract.Notes, contract.CreatedBy,
		contract.UpdatedBy, contract.CreatedAt, contract.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create deferred revenue contract: %w", err)
	}

	return nil
}

// GetDeferredRevenueContract retrieves a deferred revenue contract
func (r *AnalyticsRepository) GetDeferredRevenueContract(ctx context.Context, id uuid.UUID) (*analytics.DeferredRevenueContract, error) {
	query := `
		SELECT id, organization_id, customer_invoice_id, invoice_line_id,
		       contract_name, total_deferred_amount, start_date, end_date,
		       recognition_method, deferred_account_id, revenue_account_id,
		       status, recognized_amount, notes, created_by, updated_by,
		       created_at, updated_at, deleted_at
		FROM deferred_revenue_contracts
		WHERE id = $1 AND deleted_at IS NULL
	`

	var contract analytics.DeferredRevenueContract
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&contract.ID, &contract.OrganizationID, &contract.CustomerInvoiceID,
		&contract.InvoiceLineID, &contract.ContractName, &contract.TotalDeferredAmount,
		&contract.StartDate, &contract.EndDate, &contract.RecognitionMethod,
		&contract.DeferredAccountID, &contract.RevenueAccountID, &contract.Status,
		&contract.RecognizedAmount, &contract.Notes, &contract.CreatedBy,
		&contract.UpdatedBy, &contract.CreatedAt, &contract.UpdatedAt,
		&contract.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get deferred revenue contract: %w", err)
	}

	return &contract, nil
}

// ListDeferredRevenueContracts retrieves deferred revenue contracts
func (r *AnalyticsRepository) ListDeferredRevenueContracts(ctx context.Context, query analytics.DeferralQuery) ([]*analytics.DeferredRevenueContract, int64, error) {
	where := []string{"deleted_at IS NULL", "organization_id = $1"}
	args := []interface{}{query.OrganizationID}
	idx := 2

	if query.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, string(*query.Status))
		idx++
	}

	if query.DateFrom != nil {
		where = append(where, fmt.Sprintf("start_date >= $%d", idx))
		args = append(args, *query.DateFrom)
		idx++
	}

	if query.DateTo != nil {
		where = append(where, fmt.Sprintf("end_date <= $%d", idx))
		args = append(args, *query.DateTo)
		idx++
	}

	// Get count
	countQuery := "SELECT COUNT(*) FROM deferred_revenue_contracts WHERE " + strings.Join(where, " AND ")
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count deferred revenue contracts: %w", err)
	}

	// Get records
	listQuery := `
		SELECT id, organization_id, customer_invoice_id, invoice_line_id,
		       contract_name, total_deferred_amount, start_date, end_date,
		       recognition_method, deferred_account_id, revenue_account_id,
		       status, recognized_amount, notes, created_by, updated_by,
		       created_at, updated_at, deleted_at
		FROM deferred_revenue_contracts
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY start_date DESC LIMIT $` + fmt.Sprintf("%d", idx) +
		` OFFSET $` + fmt.Sprintf("%d", idx+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list deferred revenue contracts: %w", err)
	}
	defer rows.Close()

	var contracts []*analytics.DeferredRevenueContract
	for rows.Next() {
		var contract analytics.DeferredRevenueContract
		err := rows.Scan(
			&contract.ID, &contract.OrganizationID, &contract.CustomerInvoiceID,
			&contract.InvoiceLineID, &contract.ContractName, &contract.TotalDeferredAmount,
			&contract.StartDate, &contract.EndDate, &contract.RecognitionMethod,
			&contract.DeferredAccountID, &contract.RevenueAccountID, &contract.Status,
			&contract.RecognizedAmount, &contract.Notes, &contract.CreatedBy,
			&contract.UpdatedBy, &contract.CreatedAt, &contract.UpdatedAt,
			&contract.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan deferred revenue contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}

	return contracts, total, rows.Err()
}

// UpdateDeferredRevenueContract updates a deferred revenue contract
func (r *AnalyticsRepository) UpdateDeferredRevenueContract(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE deferred_revenue_contracts SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteDeferredRevenueContract soft deletes a deferred revenue contract
func (r *AnalyticsRepository) DeleteDeferredRevenueContract(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE deferred_revenue_contracts SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// ==================== DEFERRED REVENUE SCHEDULES ====================

// CreateDeferredRevenueSchedule creates a deferred revenue schedule entry
func (r *AnalyticsRepository) CreateDeferredRevenueSchedule(ctx context.Context, schedule *analytics.DeferredRevenueSchedule) error {
	query := `
		INSERT INTO deferred_revenue_schedule (
			id, contract_id, line_number, recognition_date,
			recognition_amount, status, journal_entry_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(ctx, query,
		schedule.ID, schedule.ContractID, schedule.LineNumber,
		schedule.RecognitionDate, schedule.RecognitionAmount, schedule.Status,
		schedule.JournalEntryID, schedule.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create deferred revenue schedule: %w", err)
	}

	return nil
}

// GetDeferredRevenueSchedule retrieves a deferred revenue schedule
func (r *AnalyticsRepository) GetDeferredRevenueSchedule(ctx context.Context, id uuid.UUID) (*analytics.DeferredRevenueSchedule, error) {
	query := `
		SELECT id, contract_id, line_number, recognition_date,
		       recognition_amount, status, journal_entry_id, created_at,
		       posted_at, deleted_at
		FROM deferred_revenue_schedule
		WHERE id = $1 AND deleted_at IS NULL
	`

	var schedule analytics.DeferredRevenueSchedule
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&schedule.ID, &schedule.ContractID, &schedule.LineNumber,
		&schedule.RecognitionDate, &schedule.RecognitionAmount, &schedule.Status,
		&schedule.JournalEntryID, &schedule.CreatedAt, &schedule.PostedAt,
		&schedule.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get deferred revenue schedule: %w", err)
	}

	return &schedule, nil
}

// ListDeferredRevenueSchedules retrieves deferred revenue schedules
func (r *AnalyticsRepository) ListDeferredRevenueSchedules(ctx context.Context, contractID uuid.UUID, limit, offset int) ([]*analytics.DeferredRevenueSchedule, int64, error) {
	countQuery := `SELECT COUNT(*) FROM deferred_revenue_schedule WHERE contract_id = $1 AND deleted_at IS NULL`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, contractID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	query := `
		SELECT id, contract_id, line_number, recognition_date,
		       recognition_amount, status, journal_entry_id, created_at,
		       posted_at, deleted_at
		FROM deferred_revenue_schedule
		WHERE contract_id = $1 AND deleted_at IS NULL
		ORDER BY line_number ASC LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, contractID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list deferred revenue schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*analytics.DeferredRevenueSchedule
	for rows.Next() {
		var schedule analytics.DeferredRevenueSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.ContractID, &schedule.LineNumber,
			&schedule.RecognitionDate, &schedule.RecognitionAmount, &schedule.Status,
			&schedule.JournalEntryID, &schedule.CreatedAt, &schedule.PostedAt,
			&schedule.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	return schedules, total, rows.Err()
}

// UpdateDeferredRevenueSchedule updates a deferred revenue schedule
func (r *AnalyticsRepository) UpdateDeferredRevenueSchedule(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE deferred_revenue_schedule SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteDeferredRevenueSchedule soft deletes a deferred revenue schedule
func (r *AnalyticsRepository) DeleteDeferredRevenueSchedule(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE deferred_revenue_schedule SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// GetPendingDeferredRevenueSchedules retrieves pending schedules before a date
func (r *AnalyticsRepository) GetPendingDeferredRevenueSchedules(ctx context.Context, organizationID uuid.UUID, beforeDate time.Time) ([]*analytics.DeferredRevenueSchedule, error) {
	query := `
		SELECT drs.id, drs.contract_id, drs.line_number, drs.recognition_date,
		       drs.recognition_amount, drs.status, drs.journal_entry_id, drs.created_at,
		       drs.posted_at, drs.deleted_at
		FROM deferred_revenue_schedule drs
		JOIN deferred_revenue_contracts drc ON drs.contract_id = drc.id
		WHERE drc.organization_id = $1 AND drs.status = 'pending' AND drs.recognition_date <= $2
		      AND drs.deleted_at IS NULL
		ORDER BY drs.recognition_date ASC
	`

	rows, err := r.pool.Query(ctx, query, organizationID, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending revenue schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*analytics.DeferredRevenueSchedule
	for rows.Next() {
		var schedule analytics.DeferredRevenueSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.ContractID, &schedule.LineNumber,
			&schedule.RecognitionDate, &schedule.RecognitionAmount, &schedule.Status,
			&schedule.JournalEntryID, &schedule.CreatedAt, &schedule.PostedAt,
			&schedule.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	return schedules, rows.Err()
}

// ==================== DEFERRED EXPENSE ====================

// CreateDeferredExpenseContract creates a deferred expense contract
func (r *AnalyticsRepository) CreateDeferredExpenseContract(ctx context.Context, contract *analytics.DeferredExpenseContract) error {
	query := `
		INSERT INTO deferred_expense_contracts (
			id, organization_id, vendor_bill_id, bill_line_id,
			contract_name, total_deferred_amount, start_date, end_date,
			recognition_method, deferred_account_id, expense_account_id,
			status, recognized_amount, notes, created_by, updated_by,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.pool.Exec(ctx, query,
		contract.ID, contract.OrganizationID, contract.VendorBillID,
		contract.BillLineID, contract.ContractName, contract.TotalDeferredAmount,
		contract.StartDate, contract.EndDate, contract.RecognitionMethod,
		contract.DeferredAccountID, contract.ExpenseAccountID, contract.Status,
		contract.RecognizedAmount, contract.Notes, contract.CreatedBy,
		contract.UpdatedBy, contract.CreatedAt, contract.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create deferred expense contract: %w", err)
	}

	return nil
}

// GetDeferredExpenseContract retrieves a deferred expense contract
func (r *AnalyticsRepository) GetDeferredExpenseContract(ctx context.Context, id uuid.UUID) (*analytics.DeferredExpenseContract, error) {
	query := `
		SELECT id, organization_id, vendor_bill_id, bill_line_id,
		       contract_name, total_deferred_amount, start_date, end_date,
		       recognition_method, deferred_account_id, expense_account_id,
		       status, recognized_amount, notes, created_by, updated_by,
		       created_at, updated_at, deleted_at
		FROM deferred_expense_contracts
		WHERE id = $1 AND deleted_at IS NULL
	`

	var contract analytics.DeferredExpenseContract
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&contract.ID, &contract.OrganizationID, &contract.VendorBillID,
		&contract.BillLineID, &contract.ContractName, &contract.TotalDeferredAmount,
		&contract.StartDate, &contract.EndDate, &contract.RecognitionMethod,
		&contract.DeferredAccountID, &contract.ExpenseAccountID, &contract.Status,
		&contract.RecognizedAmount, &contract.Notes, &contract.CreatedBy,
		&contract.UpdatedBy, &contract.CreatedAt, &contract.UpdatedAt,
		&contract.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get deferred expense contract: %w", err)
	}

	return &contract, nil
}

// ListDeferredExpenseContracts retrieves deferred expense contracts
func (r *AnalyticsRepository) ListDeferredExpenseContracts(ctx context.Context, query analytics.DeferralQuery) ([]*analytics.DeferredExpenseContract, int64, error) {
	where := []string{"deleted_at IS NULL", "organization_id = $1"}
	args := []interface{}{query.OrganizationID}
	idx := 2

	if query.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, string(*query.Status))
		idx++
	}

	// Get count
	countQuery := "SELECT COUNT(*) FROM deferred_expense_contracts WHERE " + strings.Join(where, " AND ")
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count deferred expense contracts: %w", err)
	}

	// Get records
	listQuery := `
		SELECT id, organization_id, vendor_bill_id, bill_line_id,
		       contract_name, total_deferred_amount, start_date, end_date,
		       recognition_method, deferred_account_id, expense_account_id,
		       status, recognized_amount, notes, created_by, updated_by,
		       created_at, updated_at, deleted_at
		FROM deferred_expense_contracts
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY start_date DESC LIMIT $` + fmt.Sprintf("%d", idx) +
		` OFFSET $` + fmt.Sprintf("%d", idx+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list deferred expense contracts: %w", err)
	}
	defer rows.Close()

	var contracts []*analytics.DeferredExpenseContract
	for rows.Next() {
		var contract analytics.DeferredExpenseContract
		err := rows.Scan(
			&contract.ID, &contract.OrganizationID, &contract.VendorBillID,
			&contract.BillLineID, &contract.ContractName, &contract.TotalDeferredAmount,
			&contract.StartDate, &contract.EndDate, &contract.RecognitionMethod,
			&contract.DeferredAccountID, &contract.ExpenseAccountID, &contract.Status,
			&contract.RecognizedAmount, &contract.Notes, &contract.CreatedBy,
			&contract.UpdatedBy, &contract.CreatedAt, &contract.UpdatedAt,
			&contract.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan deferred expense contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}

	return contracts, total, rows.Err()
}

// UpdateDeferredExpenseContract updates a deferred expense contract
func (r *AnalyticsRepository) UpdateDeferredExpenseContract(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE deferred_expense_contracts SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteDeferredExpenseContract soft deletes a deferred expense contract
func (r *AnalyticsRepository) DeleteDeferredExpenseContract(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE deferred_expense_contracts SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// ==================== DEFERRED EXPENSE SCHEDULES ====================

// CreateDeferredExpenseSchedule creates a deferred expense schedule entry
func (r *AnalyticsRepository) CreateDeferredExpenseSchedule(ctx context.Context, schedule *analytics.DeferredExpenseSchedule) error {
	query := `
		INSERT INTO deferred_expense_schedule (
			id, contract_id, line_number, recognition_date,
			recognition_amount, status, journal_entry_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(ctx, query,
		schedule.ID, schedule.ContractID, schedule.LineNumber,
		schedule.RecognitionDate, schedule.RecognitionAmount, schedule.Status,
		schedule.JournalEntryID, schedule.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create deferred expense schedule: %w", err)
	}

	return nil
}

// GetDeferredExpenseSchedule retrieves a deferred expense schedule
func (r *AnalyticsRepository) GetDeferredExpenseSchedule(ctx context.Context, id uuid.UUID) (*analytics.DeferredExpenseSchedule, error) {
	query := `
		SELECT id, contract_id, line_number, recognition_date,
		       recognition_amount, status, journal_entry_id, created_at,
		       posted_at, deleted_at
		FROM deferred_expense_schedule
		WHERE id = $1 AND deleted_at IS NULL
	`

	var schedule analytics.DeferredExpenseSchedule
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&schedule.ID, &schedule.ContractID, &schedule.LineNumber,
		&schedule.RecognitionDate, &schedule.RecognitionAmount, &schedule.Status,
		&schedule.JournalEntryID, &schedule.CreatedAt, &schedule.PostedAt,
		&schedule.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get deferred expense schedule: %w", err)
	}

	return &schedule, nil
}

// ListDeferredExpenseSchedules retrieves deferred expense schedules
func (r *AnalyticsRepository) ListDeferredExpenseSchedules(ctx context.Context, contractID uuid.UUID, limit, offset int) ([]*analytics.DeferredExpenseSchedule, int64, error) {
	countQuery := `SELECT COUNT(*) FROM deferred_expense_schedule WHERE contract_id = $1 AND deleted_at IS NULL`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, contractID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	query := `
		SELECT id, contract_id, line_number, recognition_date,
		       recognition_amount, status, journal_entry_id, created_at,
		       posted_at, deleted_at
		FROM deferred_expense_schedule
		WHERE contract_id = $1 AND deleted_at IS NULL
		ORDER BY line_number ASC LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, contractID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list deferred expense schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*analytics.DeferredExpenseSchedule
	for rows.Next() {
		var schedule analytics.DeferredExpenseSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.ContractID, &schedule.LineNumber,
			&schedule.RecognitionDate, &schedule.RecognitionAmount, &schedule.Status,
			&schedule.JournalEntryID, &schedule.CreatedAt, &schedule.PostedAt,
			&schedule.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	return schedules, total, rows.Err()
}

// UpdateDeferredExpenseSchedule updates a deferred expense schedule
func (r *AnalyticsRepository) UpdateDeferredExpenseSchedule(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE deferred_expense_schedule SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteDeferredExpenseSchedule soft deletes a deferred expense schedule
func (r *AnalyticsRepository) DeleteDeferredExpenseSchedule(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE deferred_expense_schedule SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// GetPendingDeferredExpenseSchedules retrieves pending schedules before a date
func (r *AnalyticsRepository) GetPendingDeferredExpenseSchedules(ctx context.Context, organizationID uuid.UUID, beforeDate time.Time) ([]*analytics.DeferredExpenseSchedule, error) {
	query := `
		SELECT des.id, des.contract_id, des.line_number, des.recognition_date,
		       des.recognition_amount, des.status, des.journal_entry_id, des.created_at,
		       des.posted_at, des.deleted_at
		FROM deferred_expense_schedule des
		JOIN deferred_expense_contracts dec ON des.contract_id = dec.id
		WHERE dec.organization_id = $1 AND des.status = 'pending' AND des.recognition_date <= $2
		      AND des.deleted_at IS NULL
		ORDER BY des.recognition_date ASC
	`

	rows, err := r.pool.Query(ctx, query, organizationID, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending expense schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*analytics.DeferredExpenseSchedule
	for rows.Next() {
		var schedule analytics.DeferredExpenseSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.ContractID, &schedule.LineNumber,
			&schedule.RecognitionDate, &schedule.RecognitionAmount, &schedule.Status,
			&schedule.JournalEntryID, &schedule.CreatedAt, &schedule.PostedAt,
			&schedule.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	return schedules, rows.Err()
}

// ==================== BUDGETS ====================

// CreateBudget creates a new budget
func (r *AnalyticsRepository) CreateBudget(ctx context.Context, budget *analytics.Budget) error {
	query := `
		INSERT INTO budgets (
			id, organization_id, budget_code, budget_name, fiscal_year_id,
			start_date, end_date, budget_type, status, notes,
			created_by, updated_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.pool.Exec(ctx, query,
		budget.ID, budget.OrganizationID, budget.BudgetCode, budget.BudgetName,
		budget.FiscalYearID, budget.StartDate, budget.EndDate, budget.BudgetType,
		budget.Status, budget.Notes, budget.CreatedBy, budget.UpdatedBy,
		budget.CreatedAt, budget.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return analytics.ErrDuplicateBudgetCode
		}
		return fmt.Errorf("failed to create budget: %w", err)
	}

	return nil
}

// GetBudget retrieves a budget
func (r *AnalyticsRepository) GetBudget(ctx context.Context, id uuid.UUID) (*analytics.Budget, error) {
	query := `
		SELECT id, organization_id, budget_code, budget_name, fiscal_year_id,
		       start_date, end_date, budget_type, status, notes,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM budgets
		WHERE id = $1 AND deleted_at IS NULL
	`

	var budget analytics.Budget
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&budget.ID, &budget.OrganizationID, &budget.BudgetCode, &budget.BudgetName,
		&budget.FiscalYearID, &budget.StartDate, &budget.EndDate, &budget.BudgetType,
		&budget.Status, &budget.Notes, &budget.CreatedBy, &budget.UpdatedBy,
		&budget.CreatedAt, &budget.UpdatedAt, &budget.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	return &budget, nil
}

// GetBudgetByCode retrieves a budget by code
func (r *AnalyticsRepository) GetBudgetByCode(ctx context.Context, organizationID uuid.UUID, code string) (*analytics.Budget, error) {
	query := `
		SELECT id, organization_id, budget_code, budget_name, fiscal_year_id,
		       start_date, end_date, budget_type, status, notes,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM budgets
		WHERE organization_id = $1 AND budget_code = $2 AND deleted_at IS NULL
	`

	var budget analytics.Budget
	err := r.pool.QueryRow(ctx, query, organizationID, code).Scan(
		&budget.ID, &budget.OrganizationID, &budget.BudgetCode, &budget.BudgetName,
		&budget.FiscalYearID, &budget.StartDate, &budget.EndDate, &budget.BudgetType,
		&budget.Status, &budget.Notes, &budget.CreatedBy, &budget.UpdatedBy,
		&budget.CreatedAt, &budget.UpdatedAt, &budget.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get budget by code: %w", err)
	}

	return &budget, nil
}

// ListBudgets retrieves budgets
func (r *AnalyticsRepository) ListBudgets(ctx context.Context, query analytics.BudgetQuery) ([]*analytics.Budget, int64, error) {
	where := []string{"deleted_at IS NULL", "organization_id = $1"}
	args := []interface{}{query.OrganizationID}
	idx := 2

	if query.FiscalYearID != nil {
		where = append(where, fmt.Sprintf("fiscal_year_id = $%d", idx))
		args = append(args, *query.FiscalYearID)
		idx++
	}

	if query.BudgetType != nil {
		where = append(where, fmt.Sprintf("budget_type = $%d", idx))
		args = append(args, string(*query.BudgetType))
		idx++
	}

	if query.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, string(*query.Status))
		idx++
	}

	// Get count
	countQuery := "SELECT COUNT(*) FROM budgets WHERE " + strings.Join(where, " AND ")
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count budgets: %w", err)
	}

	// Get records
	listQuery := `
		SELECT id, organization_id, budget_code, budget_name, fiscal_year_id,
		       start_date, end_date, budget_type, status, notes,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM budgets
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY budget_code ASC LIMIT $` + fmt.Sprintf("%d", idx) +
		` OFFSET $` + fmt.Sprintf("%d", idx+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list budgets: %w", err)
	}
	defer rows.Close()

	var budgets []*analytics.Budget
	for rows.Next() {
		var budget analytics.Budget
		err := rows.Scan(
			&budget.ID, &budget.OrganizationID, &budget.BudgetCode, &budget.BudgetName,
			&budget.FiscalYearID, &budget.StartDate, &budget.EndDate, &budget.BudgetType,
			&budget.Status, &budget.Notes, &budget.CreatedBy, &budget.UpdatedBy,
			&budget.CreatedAt, &budget.UpdatedAt, &budget.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, &budget)
	}

	return budgets, total, rows.Err()
}

// UpdateBudget updates a budget
func (r *AnalyticsRepository) UpdateBudget(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE budgets SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteBudget soft deletes a budget
func (r *AnalyticsRepository) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE budgets SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// ==================== BUDGET LINES ====================

// CreateBudgetLine creates a budget line
func (r *AnalyticsRepository) CreateBudgetLine(ctx context.Context, line *analytics.BudgetLine) error {
	query := `
		INSERT INTO budget_lines (
			id, budget_id, account_id, analytic_account_id,
			accounting_period_id, period_start_date, period_end_date,
			planned_amount, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.pool.Exec(ctx, query,
		line.ID, line.BudgetID, line.AccountID, line.AnalyticAccountID,
		line.AccountingPeriodID, line.PeriodStartDate, line.PeriodEndDate,
		line.PlannedAmount, line.Notes, line.CreatedAt, line.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create budget line: %w", err)
	}

	return nil
}

// GetBudgetLine retrieves a budget line
func (r *AnalyticsRepository) GetBudgetLine(ctx context.Context, id uuid.UUID) (*analytics.BudgetLine, error) {
	query := `
		SELECT id, budget_id, account_id, analytic_account_id,
		       accounting_period_id, period_start_date, period_end_date,
		       planned_amount, notes, created_at, updated_at, deleted_at
		FROM budget_lines
		WHERE id = $1 AND deleted_at IS NULL
	`

	var line analytics.BudgetLine
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&line.ID, &line.BudgetID, &line.AccountID, &line.AnalyticAccountID,
		&line.AccountingPeriodID, &line.PeriodStartDate, &line.PeriodEndDate,
		&line.PlannedAmount, &line.Notes, &line.CreatedAt, &line.UpdatedAt,
		&line.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get budget line: %w", err)
	}

	return &line, nil
}

// ListBudgetLines retrieves budget lines
func (r *AnalyticsRepository) ListBudgetLines(ctx context.Context, budgetID uuid.UUID, limit, offset int) ([]*analytics.BudgetLine, int64, error) {
	countQuery := `SELECT COUNT(*) FROM budget_lines WHERE budget_id = $1 AND deleted_at IS NULL`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, budgetID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count budget lines: %w", err)
	}

	query := `
		SELECT id, budget_id, account_id, analytic_account_id,
		       accounting_period_id, period_start_date, period_end_date,
		       planned_amount, notes, created_at, updated_at, deleted_at
		FROM budget_lines
		WHERE budget_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, budgetID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list budget lines: %w", err)
	}
	defer rows.Close()

	var lines []*analytics.BudgetLine
	for rows.Next() {
		var line analytics.BudgetLine
		err := rows.Scan(
			&line.ID, &line.BudgetID, &line.AccountID, &line.AnalyticAccountID,
			&line.AccountingPeriodID, &line.PeriodStartDate, &line.PeriodEndDate,
			&line.PlannedAmount, &line.Notes, &line.CreatedAt, &line.UpdatedAt,
			&line.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan budget line: %w", err)
		}
		lines = append(lines, &line)
	}

	return lines, total, rows.Err()
}

// UpdateBudgetLine updates a budget line
func (r *AnalyticsRepository) UpdateBudgetLine(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE budget_lines SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteBudgetLine soft deletes a budget line
func (r *AnalyticsRepository) DeleteBudgetLine(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE budget_lines SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

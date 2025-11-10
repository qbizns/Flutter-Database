package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/monitoring"
)

type POSErrorLogRepository struct {
	db *DB
}

func NewPOSErrorLogRepository(db *DB) *POSErrorLogRepository {
	return &POSErrorLogRepository{db: db}
}

func (r *POSErrorLogRepository) List(ctx context.Context, orgID *uuid.UUID, filters monitoring.POSErrorLogFilters) ([]monitoring.POSErrorLog, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT id, organization_id, error_level, error_code, error_message,
		       device_id, user_id, pos_session_id, sale_id, stack_trace,
		       request_data, error_data, is_resolved, resolved_by, resolved_at,
		       resolution_notes, occurred_at, created_at
		FROM pos_error_logs
		WHERE 1=1
	`

	var args []interface{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.ErrorLevel != nil {
		argCount++
		query += fmt.Sprintf(" AND error_level = $%d", argCount)
		args = append(args, *filters.ErrorLevel)
	}

	if filters.DeviceID != nil {
		argCount++
		query += fmt.Sprintf(" AND device_id = $%d", argCount)
		args = append(args, *filters.DeviceID)
	}

	if filters.UserID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *filters.UserID)
	}

	if filters.POSSessionID != nil {
		argCount++
		query += fmt.Sprintf(" AND pos_session_id = $%d", argCount)
		args = append(args, *filters.POSSessionID)
	}

	if filters.IsResolved != nil {
		argCount++
		query += fmt.Sprintf(" AND is_resolved = $%d", argCount)
		args = append(args, *filters.IsResolved)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND occurred_at >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND occurred_at <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY occurred_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var errors []monitoring.POSErrorLog
	for rows.Next() {
		var e monitoring.POSErrorLog
		err := rows.Scan(
			&e.ID, &e.OrganizationID, &e.ErrorLevel, &e.ErrorCode, &e.ErrorMessage,
			&e.DeviceID, &e.UserID, &e.POSSessionID, &e.SaleID, &e.StackTrace,
			&e.RequestData, &e.ErrorData, &e.IsResolved, &e.ResolvedBy, &e.ResolvedAt,
			&e.ResolutionNotes, &e.OccurredAt, &e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		errors = append(errors, e)
	}

	return errors, rows.Err()
}

func (r *POSErrorLogRepository) Count(ctx context.Context, orgID *uuid.UUID, filters monitoring.POSErrorLogFilters) (int64, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return 0, err
		}
	}

	query := "SELECT COUNT(*) FROM pos_error_logs WHERE 1=1"
	var args []interface{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.ErrorLevel != nil {
		argCount++
		query += fmt.Sprintf(" AND error_level = $%d", argCount)
		args = append(args, *filters.ErrorLevel)
	}

	if filters.IsResolved != nil {
		argCount++
		query += fmt.Sprintf(" AND is_resolved = $%d", argCount)
		args = append(args, *filters.IsResolved)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *POSErrorLogRepository) Create(ctx context.Context, errorLog *monitoring.POSErrorLog) error {
	if errorLog.OrganizationID != nil {
		if err := r.db.SetOrganizationContext(ctx, errorLog.OrganizationID.String()); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO pos_error_logs (
			id, organization_id, error_level, error_code, error_message,
			device_id, user_id, pos_session_id, sale_id, stack_trace,
			request_data, error_data, is_resolved, resolved_by, resolved_at,
			resolution_notes, occurred_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		errorLog.ID, errorLog.OrganizationID, errorLog.ErrorLevel, errorLog.ErrorCode, errorLog.ErrorMessage,
		errorLog.DeviceID, errorLog.UserID, errorLog.POSSessionID, errorLog.SaleID, errorLog.StackTrace,
		errorLog.RequestData, errorLog.ErrorData, errorLog.IsResolved, errorLog.ResolvedBy, errorLog.ResolvedAt,
		errorLog.ResolutionNotes, errorLog.OccurredAt, errorLog.CreatedAt,
	)
	return err
}

func (r *POSErrorLogRepository) Get(ctx context.Context, id uuid.UUID) (*monitoring.POSErrorLog, error) {
	query := `
		SELECT id, organization_id, error_level, error_code, error_message,
		       device_id, user_id, pos_session_id, sale_id, stack_trace,
		       request_data, error_data, is_resolved, resolved_by, resolved_at,
		       resolution_notes, occurred_at, created_at
		FROM pos_error_logs
		WHERE id = $1
	`

	var e monitoring.POSErrorLog
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.OrganizationID, &e.ErrorLevel, &e.ErrorCode, &e.ErrorMessage,
		&e.DeviceID, &e.UserID, &e.POSSessionID, &e.SaleID, &e.StackTrace,
		&e.RequestData, &e.ErrorData, &e.IsResolved, &e.ResolvedBy, &e.ResolvedAt,
		&e.ResolutionNotes, &e.OccurredAt, &e.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *POSErrorLogRepository) Update(ctx context.Context, errorLog *monitoring.POSErrorLog) error {
	if errorLog.OrganizationID != nil {
		if err := r.db.SetOrganizationContext(ctx, errorLog.OrganizationID.String()); err != nil {
			return err
		}
	}

	query := `
		UPDATE pos_error_logs SET
			error_level = $2, error_code = $3, error_message = $4,
			device_id = $5, user_id = $6, pos_session_id = $7, sale_id = $8,
			stack_trace = $9, request_data = $10, error_data = $11,
			is_resolved = $12, resolved_by = $13, resolved_at = $14,
			resolution_notes = $15, occurred_at = $16
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query,
		errorLog.ID, errorLog.ErrorLevel, errorLog.ErrorCode, errorLog.ErrorMessage,
		errorLog.DeviceID, errorLog.UserID, errorLog.POSSessionID, errorLog.SaleID,
		errorLog.StackTrace, errorLog.RequestData, errorLog.ErrorData,
		errorLog.IsResolved, errorLog.ResolvedBy, errorLog.ResolvedAt,
		errorLog.ResolutionNotes, errorLog.OccurredAt,
	)
	return err
}

func (r *POSErrorLogRepository) Resolve(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID, notes string) error {
	query := `
		UPDATE pos_error_logs
		SET is_resolved = true, resolved_by = $2, resolved_at = $3, resolution_notes = $4
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, id, resolvedBy, time.Now(), notes)
	return err
}

// SystemHealthRepository implementation

type SystemHealthRepository struct {
	db *DB
}

func NewSystemHealthRepository(db *DB) *SystemHealthRepository {
	return &SystemHealthRepository{db: db}
}

func (r *SystemHealthRepository) List(ctx context.Context, orgID *uuid.UUID, filters monitoring.SystemHealthFilters) ([]monitoring.SystemHealth, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT id, organization_id, check_type, check_name, status,
		       last_check_at, last_success_at, last_failure_at,
		       metric_value, metric_unit, threshold_warning, threshold_critical,
		       details, created_at, updated_at
		FROM system_health
		WHERE 1=1
	`

	var args []interface{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.CheckType != nil {
		argCount++
		query += fmt.Sprintf(" AND check_type = $%d", argCount)
		args = append(args, *filters.CheckType)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	query += " ORDER BY check_name ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []monitoring.SystemHealth
	for rows.Next() {
		var h monitoring.SystemHealth
		err := rows.Scan(
			&h.ID, &h.OrganizationID, &h.CheckType, &h.CheckName, &h.Status,
			&h.LastCheckAt, &h.LastSuccessAt, &h.LastFailureAt,
			&h.MetricValue, &h.MetricUnit, &h.ThresholdWarning, &h.ThresholdCritical,
			&h.Details, &h.CreatedAt, &h.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		checks = append(checks, h)
	}

	return checks, rows.Err()
}

func (r *SystemHealthRepository) Count(ctx context.Context, orgID *uuid.UUID, filters monitoring.SystemHealthFilters) (int64, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return 0, err
		}
	}

	query := "SELECT COUNT(*) FROM system_health WHERE 1=1"
	var args []interface{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.CheckType != nil {
		argCount++
		query += fmt.Sprintf(" AND check_type = $%d", argCount)
		args = append(args, *filters.CheckType)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *SystemHealthRepository) Upsert(ctx context.Context, health *monitoring.SystemHealth) error {
	if health.OrganizationID != nil {
		if err := r.db.SetOrganizationContext(ctx, health.OrganizationID.String()); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO system_health (
			id, organization_id, check_type, check_name, status,
			last_check_at, last_success_at, last_failure_at,
			metric_value, metric_unit, threshold_warning, threshold_critical,
			details, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (organization_id, check_type, check_name)
		DO UPDATE SET
			status = EXCLUDED.status,
			last_check_at = EXCLUDED.last_check_at,
			last_success_at = EXCLUDED.last_success_at,
			last_failure_at = EXCLUDED.last_failure_at,
			metric_value = EXCLUDED.metric_value,
			metric_unit = EXCLUDED.metric_unit,
			threshold_warning = EXCLUDED.threshold_warning,
			threshold_critical = EXCLUDED.threshold_critical,
			details = EXCLUDED.details,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Pool.Exec(ctx, query,
		health.ID, health.OrganizationID, health.CheckType, health.CheckName, health.Status,
		health.LastCheckAt, health.LastSuccessAt, health.LastFailureAt,
		health.MetricValue, health.MetricUnit, health.ThresholdWarning, health.ThresholdCritical,
		health.Details, health.CreatedAt, health.UpdatedAt,
	)
	return err
}

func (r *SystemHealthRepository) Get(ctx context.Context, orgID *uuid.UUID, checkType, checkName string) (*monitoring.SystemHealth, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT id, organization_id, check_type, check_name, status,
		       last_check_at, last_success_at, last_failure_at,
		       metric_value, metric_unit, threshold_warning, threshold_critical,
		       details, created_at, updated_at
		FROM system_health
		WHERE check_type = $1 AND check_name = $2
	`

	args := []interface{}{checkType, checkName}
	if orgID != nil {
		query += " AND organization_id = $3"
		args = append(args, *orgID)
	} else {
		query += " AND organization_id IS NULL"
	}

	var h monitoring.SystemHealth
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&h.ID, &h.OrganizationID, &h.CheckType, &h.CheckName, &h.Status,
		&h.LastCheckAt, &h.LastSuccessAt, &h.LastFailureAt,
		&h.MetricValue, &h.MetricUnit, &h.ThresholdWarning, &h.ThresholdCritical,
		&h.Details, &h.CreatedAt, &h.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *SystemHealthRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, metricValue *float64) error {
	query := `
		UPDATE system_health
		SET status = $2, metric_value = $3, last_check_at = $4, updated_at = $4
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, query, id, status, metricValue, now)
	return err
}

func (r *SystemHealthRepository) GetUnhealthyChecks(ctx context.Context, orgID *uuid.UUID) ([]monitoring.SystemHealth, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT id, organization_id, check_type, check_name, status,
		       last_check_at, last_success_at, last_failure_at,
		       metric_value, metric_unit, threshold_warning, threshold_critical,
		       details, created_at, updated_at
		FROM system_health
		WHERE status IN ('degraded', 'unhealthy')
	`

	if orgID != nil {
		query += " AND organization_id = $1"
	}

	var args []interface{}
	if orgID != nil {
		args = append(args, *orgID)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []monitoring.SystemHealth
	for rows.Next() {
		var h monitoring.SystemHealth
		err := rows.Scan(
			&h.ID, &h.OrganizationID, &h.CheckType, &h.CheckName, &h.Status,
			&h.LastCheckAt, &h.LastSuccessAt, &h.LastFailureAt,
			&h.MetricValue, &h.MetricUnit, &h.ThresholdWarning, &h.ThresholdCritical,
			&h.Details, &h.CreatedAt, &h.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		checks = append(checks, h)
	}

	return checks, rows.Err()
}

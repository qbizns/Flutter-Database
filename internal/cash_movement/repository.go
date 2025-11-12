package cash_movement

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

// Repository handles database operations for CashMovements
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CashMovements repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CashMovements represents a cash_movements entity
type CashMovements struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PosSessionId uuid.UUID `json:"pos_session_id" db:"pos_session_id"`
	CashDrawerId *uuid.UUID `json:"cash_drawer_id" db:"cash_drawer_id"`
	MovementType string `json:"movement_type" db:"movement_type"`
	Amount float64 `json:"amount" db:"amount"`
	ReasonCode *string `json:"reason_code" db:"reason_code"`
	ReasonDescription string `json:"reason_description" db:"reason_description"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	RequiresApproval *bool `json:"requires_approval" db:"requires_approval"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	Notes *string `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new cash_movements record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CashMovements) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "cash_movements", duration, nil)
	}()

	query := `
		INSERT INTO cash_movements (
			, organization_id
			, pos_session_id
			, cash_drawer_id
			, movement_type
			, amount
			, reason_code
			, reason_description
			, user_id
			, requires_approval
			, approved_by
			, approved_at
			, notes
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
			, $15
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PosSessionId,
		entity.CashDrawerId,
		entity.MovementType,
		entity.Amount,
		entity.ReasonCode,
		entity.ReasonDescription,
		entity.UserId,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create cash_movements", zap.Error(err))
		return fmt.Errorf("failed to create cash_movements: %w", err)
	}

	r.logger.Info("created cash_movements",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a cash_movements by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CashMovements, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_movements", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, pos_session_id
			, cash_drawer_id
			, movement_type
			, amount
			, reason_code
			, reason_description
			, user_id
			, requires_approval
			, approved_by
			, approved_at
			, notes
			, created_at
			, deleted_at
		FROM cash_movements
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CashMovements
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PosSessionId,
		&entity.CashDrawerId,
		&entity.MovementType,
		&entity.Amount,
		&entity.ReasonCode,
		&entity.ReasonDescription,
		&entity.UserId,
		&entity.RequiresApproval,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.Notes,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("cash_movements not found")
	}

	if err != nil {
		r.logger.Error("failed to get cash_movements", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get cash_movements: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of cash_movements records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CashMovements, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_movements", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cash_movements
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cash_movements records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, pos_session_id
			, cash_drawer_id
			, movement_type
			, amount
			, reason_code
			, reason_description
			, user_id
			, requires_approval
			, approved_by
			, approved_at
			, notes
			, created_at
			, deleted_at
		FROM cash_movements
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cash_movements", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cash_movements: %w", err)
	}
	defer rows.Close()

	var entities []*CashMovements
	for rows.Next() {
		var entity CashMovements
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PosSessionId,
			&entity.CashDrawerId,
			&entity.MovementType,
			&entity.Amount,
			&entity.ReasonCode,
			&entity.ReasonDescription,
			&entity.UserId,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cash_movements: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating cash_movements rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing cash_movements record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CashMovements) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cash_movements", duration, nil)
	}()

	query := `
		UPDATE cash_movements
		SET
			, organization_id = $2
			, pos_session_id = $3
			, cash_drawer_id = $4
			, movement_type = $5
			, amount = $6
			, reason_code = $7
			, reason_description = $8
			, user_id = $9
			, requires_approval = $10
			, approved_by = $11
			, approved_at = $12
			, notes = $13
			, deleted_at = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PosSessionId,
		entity.CashDrawerId,
		entity.MovementType,
		entity.Amount,
		entity.ReasonCode,
		entity.ReasonDescription,
		entity.UserId,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update cash_movements", zap.Error(err))
		return fmt.Errorf("failed to update cash_movements: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cash_movements not found or already deleted")
	}

	r.logger.Info("updated cash_movements",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a cash_movements record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cash_movements", duration, nil)
	}()

	query := `
		UPDATE cash_movements
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete cash_movements", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete cash_movements: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cash_movements not found or already deleted")
	}

	r.logger.Info("deleted cash_movements", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves cash_movements records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CashMovements, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_movements", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cash_movements
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cash_movements records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, pos_session_id
			, cash_drawer_id
			, movement_type
			, amount
			, reason_code
			, reason_description
			, user_id
			, requires_approval
			, approved_by
			, approved_at
			, notes
			, created_at
			, deleted_at
		FROM cash_movements
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cash_movements by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cash_movements: %w", err)
	}
	defer rows.Close()

	var entities []*CashMovements
	for rows.Next() {
		var entity CashMovements
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PosSessionId,
			&entity.CashDrawerId,
			&entity.MovementType,
			&entity.Amount,
			&entity.ReasonCode,
			&entity.ReasonDescription,
			&entity.UserId,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cash_movements: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


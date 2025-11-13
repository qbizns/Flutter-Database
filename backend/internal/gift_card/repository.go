package gift_card

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

// Repository handles database operations for GiftCards
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new GiftCards repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// GiftCards represents a gift_cards entity
type GiftCards struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CardNumber string `json:"card_number" db:"card_number"`
	PinCode *string `json:"pin_code" db:"pin_code"`
	CustomerId *uuid.UUID `json:"customer_id" db:"customer_id"`
	OriginalValue float64 `json:"original_value" db:"original_value"`
	CurrentBalance float64 `json:"current_balance" db:"current_balance"`
	IssuedDate time.Time `json:"issued_date" db:"issued_date"`
	ExpiryDate *time.Time `json:"expiry_date" db:"expiry_date"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	IssuedByUserId *uuid.UUID `json:"issued_by_user_id" db:"issued_by_user_id"`
	IssuedLocationId *uuid.UUID `json:"issued_location_id" db:"issued_location_id"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new gift_cards record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *GiftCards) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "gift_cards", duration, nil)
	}()

	query := `
		INSERT INTO gift_cards (
			, organization_id
			, card_number
			, pin_code
			, customer_id
			, original_value
			, current_balance
			, issued_date
			, expiry_date
			, status
			, status
			, issued_by_user_id
			, issued_location_id
			, notes
			, created_by
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
			, $18
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CardNumber,
		entity.PinCode,
		entity.CustomerId,
		entity.OriginalValue,
		entity.CurrentBalance,
		entity.IssuedDate,
		entity.ExpiryDate,
		entity.Status,
		entity.Status,
		entity.IssuedByUserId,
		entity.IssuedLocationId,
		entity.Notes,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create gift_cards", zap.Error(err))
		return fmt.Errorf("failed to create gift_cards: %w", err)
	}

	r.logger.Info("created gift_cards",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a gift_cards by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*GiftCards, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "gift_cards", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, card_number
			, pin_code
			, customer_id
			, original_value
			, current_balance
			, issued_date
			, expiry_date
			, status
			, status
			, issued_by_user_id
			, issued_location_id
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM gift_cards
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity GiftCards
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CardNumber,
		&entity.PinCode,
		&entity.CustomerId,
		&entity.OriginalValue,
		&entity.CurrentBalance,
		&entity.IssuedDate,
		&entity.ExpiryDate,
		&entity.Status,
		&entity.Status,
		&entity.IssuedByUserId,
		&entity.IssuedLocationId,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("gift_cards not found")
	}

	if err != nil {
		r.logger.Error("failed to get gift_cards", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get gift_cards: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of gift_cards records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*GiftCards, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "gift_cards", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM gift_cards
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count gift_cards records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, card_number
			, pin_code
			, customer_id
			, original_value
			, current_balance
			, issued_date
			, expiry_date
			, status
			, status
			, issued_by_user_id
			, issued_location_id
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM gift_cards
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list gift_cards", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list gift_cards: %w", err)
	}
	defer rows.Close()

	var entities []*GiftCards
	for rows.Next() {
		var entity GiftCards
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CardNumber,
			&entity.PinCode,
			&entity.CustomerId,
			&entity.OriginalValue,
			&entity.CurrentBalance,
			&entity.IssuedDate,
			&entity.ExpiryDate,
			&entity.Status,
			&entity.Status,
			&entity.IssuedByUserId,
			&entity.IssuedLocationId,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan gift_cards: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating gift_cards rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing gift_cards record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *GiftCards) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "gift_cards", duration, nil)
	}()

	query := `
		UPDATE gift_cards
		SET
			, organization_id = $2
			, card_number = $3
			, pin_code = $4
			, customer_id = $5
			, original_value = $6
			, current_balance = $7
			, issued_date = $8
			, expiry_date = $9
			, status = $10
			, status = $11
			, issued_by_user_id = $12
			, issued_location_id = $13
			, notes = $14
			, created_by = $15
			, updated_at = $17
			, deleted_at = $18
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CardNumber,
		entity.PinCode,
		entity.CustomerId,
		entity.OriginalValue,
		entity.CurrentBalance,
		entity.IssuedDate,
		entity.ExpiryDate,
		entity.Status,
		entity.Status,
		entity.IssuedByUserId,
		entity.IssuedLocationId,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update gift_cards", zap.Error(err))
		return fmt.Errorf("failed to update gift_cards: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("gift_cards not found or already deleted")
	}

	r.logger.Info("updated gift_cards",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a gift_cards record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "gift_cards", duration, nil)
	}()

	query := `
		UPDATE gift_cards
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete gift_cards", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete gift_cards: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("gift_cards not found or already deleted")
	}

	r.logger.Info("deleted gift_cards", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves gift_cards records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*GiftCards, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "gift_cards", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM gift_cards
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count gift_cards records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, card_number
			, pin_code
			, customer_id
			, original_value
			, current_balance
			, issued_date
			, expiry_date
			, status
			, status
			, issued_by_user_id
			, issued_location_id
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM gift_cards
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list gift_cards by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list gift_cards: %w", err)
	}
	defer rows.Close()

	var entities []*GiftCards
	for rows.Next() {
		var entity GiftCards
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CardNumber,
			&entity.PinCode,
			&entity.CustomerId,
			&entity.OriginalValue,
			&entity.CurrentBalance,
			&entity.IssuedDate,
			&entity.ExpiryDate,
			&entity.Status,
			&entity.Status,
			&entity.IssuedByUserId,
			&entity.IssuedLocationId,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan gift_cards: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


package order_item_modifier

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

// Repository handles database operations for OrderItemModifiers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new OrderItemModifiers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// OrderItemModifiers represents a order_item_modifiers entity
type OrderItemModifiers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	OrderItemId uuid.UUID `json:"order_item_id" db:"order_item_id"`
	ModifierId uuid.UUID `json:"modifier_id" db:"modifier_id"`
	ModifierGroupId *uuid.UUID `json:"modifier_group_id" db:"modifier_group_id"`
	ModifierName string `json:"modifier_name" db:"modifier_name"`
	Quantity *int64 `json:"quantity" db:"quantity"`
	PriceAdjustment *float64 `json:"price_adjustment" db:"price_adjustment"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new order_item_modifiers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *OrderItemModifiers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "order_item_modifiers", duration, nil)
	}()

	query := `
		INSERT INTO order_item_modifiers (
			, organization_id
			, order_item_id
			, modifier_id
			, modifier_group_id
			, modifier_name
			, quantity
			, price_adjustment
			, display_order
			, metadata
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
			, $13
			, $14
			, $15
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.OrderItemId,
		entity.ModifierId,
		entity.ModifierGroupId,
		entity.ModifierName,
		entity.Quantity,
		entity.PriceAdjustment,
		entity.DisplayOrder,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create order_item_modifiers", zap.Error(err))
		return fmt.Errorf("failed to create order_item_modifiers: %w", err)
	}

	r.logger.Info("created order_item_modifiers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a order_item_modifiers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*OrderItemModifiers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_item_modifiers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, order_item_id
			, modifier_id
			, modifier_group_id
			, modifier_name
			, quantity
			, price_adjustment
			, display_order
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM order_item_modifiers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity OrderItemModifiers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.OrderItemId,
		&entity.ModifierId,
		&entity.ModifierGroupId,
		&entity.ModifierName,
		&entity.Quantity,
		&entity.PriceAdjustment,
		&entity.DisplayOrder,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("order_item_modifiers not found")
	}

	if err != nil {
		r.logger.Error("failed to get order_item_modifiers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get order_item_modifiers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of order_item_modifiers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*OrderItemModifiers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_item_modifiers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM order_item_modifiers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count order_item_modifiers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_item_id
			, modifier_id
			, modifier_group_id
			, modifier_name
			, quantity
			, price_adjustment
			, display_order
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM order_item_modifiers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list order_item_modifiers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list order_item_modifiers: %w", err)
	}
	defer rows.Close()

	var entities []*OrderItemModifiers
	for rows.Next() {
		var entity OrderItemModifiers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderItemId,
			&entity.ModifierId,
			&entity.ModifierGroupId,
			&entity.ModifierName,
			&entity.Quantity,
			&entity.PriceAdjustment,
			&entity.DisplayOrder,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order_item_modifiers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating order_item_modifiers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing order_item_modifiers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *OrderItemModifiers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "order_item_modifiers", duration, nil)
	}()

	query := `
		UPDATE order_item_modifiers
		SET
			, organization_id = $2
			, order_item_id = $3
			, modifier_id = $4
			, modifier_group_id = $5
			, modifier_name = $6
			, quantity = $7
			, price_adjustment = $8
			, display_order = $9
			, metadata = $10
			, updated_at = $12
			, created_by = $13
			, updated_by = $14
			, deleted_at = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.OrderItemId,
		entity.ModifierId,
		entity.ModifierGroupId,
		entity.ModifierName,
		entity.Quantity,
		entity.PriceAdjustment,
		entity.DisplayOrder,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update order_item_modifiers", zap.Error(err))
		return fmt.Errorf("failed to update order_item_modifiers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order_item_modifiers not found or already deleted")
	}

	r.logger.Info("updated order_item_modifiers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a order_item_modifiers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "order_item_modifiers", duration, nil)
	}()

	query := `
		UPDATE order_item_modifiers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete order_item_modifiers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete order_item_modifiers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order_item_modifiers not found or already deleted")
	}

	r.logger.Info("deleted order_item_modifiers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves order_item_modifiers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*OrderItemModifiers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_item_modifiers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM order_item_modifiers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count order_item_modifiers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_item_id
			, modifier_id
			, modifier_group_id
			, modifier_name
			, quantity
			, price_adjustment
			, display_order
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM order_item_modifiers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list order_item_modifiers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list order_item_modifiers: %w", err)
	}
	defer rows.Close()

	var entities []*OrderItemModifiers
	for rows.Next() {
		var entity OrderItemModifiers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderItemId,
			&entity.ModifierId,
			&entity.ModifierGroupId,
			&entity.ModifierName,
			&entity.Quantity,
			&entity.PriceAdjustment,
			&entity.DisplayOrder,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order_item_modifiers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


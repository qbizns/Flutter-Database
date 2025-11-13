package pos_account_mapping

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

// Repository handles database operations for PosAccountMappings
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PosAccountMappings repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PosAccountMappings represents a pos_account_mappings entity
type PosAccountMappings struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SourceType string `json:"source_type" db:"source_type"`
	'product', *string `json:"'product'," db:"'product',"`
	'category', *string `json:"'category'," db:"'category',"`
	'paymentMethod', *string `json:"'payment_method'," db:"'payment_method',"`
	'salesChannel', *string `json:"'sales_channel'," db:"'sales_channel',"`
	'discount', *string `json:"'discount'," db:"'discount',"`
	'rounding', *string `json:"'rounding'," db:"'rounding',"`
	'tax', *string `json:"'tax'," db:"'tax',"`
	'serviceCharge', *string `json:"'service_charge'," db:"'service_charge',"`
	'shipping', *string `json:"'shipping'," db:"'shipping',"`
	'giftCard', *string `json:"'gift_card'," db:"'gift_card',"`
	'storeCredit', *string `json:"'store_credit'," db:"'store_credit',"`
	'loyaltyRedemption',-- *string `json:"'loyalty_redemption',--" db:"'loyalty_redemption',--"`
	'default' *string `json:"'default'" db:"'default'"`
	SourceId *uuid.UUID `json:"source_id" db:"source_id"`
	SourceCode *string `json:"source_code" db:"source_code"`
	Purpose string `json:"purpose" db:"purpose"`
	'revenue', *string `json:"'revenue'," db:"'revenue',"`
	'cogs', *string `json:"'cogs'," db:"'cogs',"`
	'inventory', *string `json:"'inventory'," db:"'inventory',"`
	'expense', *string `json:"'expense'," db:"'expense',"`
	'liability', *string `json:"'liability'," db:"'liability',"`
	'asset', *string `json:"'asset'," db:"'asset',"`
	'discountExpense', *string `json:"'discount_expense'," db:"'discount_expense',"`
	'discountContra', *string `json:"'discount_contra'," db:"'discount_contra',"`
	'taxLiability', *string `json:"'tax_liability'," db:"'tax_liability',"`
	'rounding', *string `json:"'rounding'," db:"'rounding',"`
	'clearing' *string `json:"'clearing'" db:"'clearing'"`
	AccountId uuid.UUID `json:"account_id" db:"account_id"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Priority *int64 `json:"priority" db:"priority"`
	Conditions json.RawMessage `json:"conditions" db:"conditions"`
	EffectiveFrom *time.Time `json:"effective_from" db:"effective_from"`
	EffectiveTo *time.Time `json:"effective_to" db:"effective_to"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	(sourceType *string `json:"(source_type" db:"(source_type"`
	(sourceType string `json:"(source_type" db:"(source_type"`
	EffectiveFrom *string `json:"effective_from" db:"effective_from"`
}

// Create inserts a new pos_account_mappings record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PosAccountMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "pos_account_mappings", duration, nil)
	}()

	query := `
		INSERT INTO pos_account_mappings (
			, organization_id
			, source_type
			, 'product',
			, 'category',
			, 'payment_method',
			, 'sales_channel',
			, 'discount',
			, 'rounding',
			, 'tax',
			, 'service_charge',
			, 'shipping',
			, 'gift_card',
			, 'store_credit',
			, 'loyalty_redemption',--
			, 'default'
			, source_id
			, source_code
			, purpose
			, 'revenue',
			, 'cogs',
			, 'inventory',
			, 'expense',
			, 'liability',
			, 'asset',
			, 'discount_expense',
			, 'discount_contra',
			, 'tax_liability',
			, 'rounding',
			, 'clearing'
			, account_id
			, is_default
			, is_active
			, priority
			, conditions
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, (source_type
			, (source_type
			, effective_from
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
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $27
			, $28
			, $29
			, $30
			, $31
			, $32
			, $33
			, $34
			, $35
			, $36
			, $37
			, $38
			, $39
			, $40
			, $43
			, $44
			, $45
			, $46
			, $47
			, $48
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SourceType,
		entity.'product',,
		entity.'category',,
		entity.'paymentMethod',,
		entity.'salesChannel',,
		entity.'discount',,
		entity.'rounding',,
		entity.'tax',,
		entity.'serviceCharge',,
		entity.'shipping',,
		entity.'giftCard',,
		entity.'storeCredit',,
		entity.'loyaltyRedemption',--,
		entity.'default',
		entity.SourceId,
		entity.SourceCode,
		entity.Purpose,
		entity.'revenue',,
		entity.'cogs',,
		entity.'inventory',,
		entity.'expense',,
		entity.'liability',,
		entity.'asset',,
		entity.'discountExpense',,
		entity.'discountContra',,
		entity.'taxLiability',,
		entity.'rounding',,
		entity.'clearing',
		entity.AccountId,
		entity.IsDefault,
		entity.IsActive,
		entity.Priority,
		entity.Conditions,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.(sourceType,
		entity.(sourceType,
		entity.EffectiveFrom,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create pos_account_mappings", zap.Error(err))
		return fmt.Errorf("failed to create pos_account_mappings: %w", err)
	}

	r.logger.Info("created pos_account_mappings",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a pos_account_mappings by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PosAccountMappings, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_account_mappings", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, source_type
			, 'product',
			, 'category',
			, 'payment_method',
			, 'sales_channel',
			, 'discount',
			, 'rounding',
			, 'tax',
			, 'service_charge',
			, 'shipping',
			, 'gift_card',
			, 'store_credit',
			, 'loyalty_redemption',--
			, 'default'
			, source_id
			, source_code
			, purpose
			, 'revenue',
			, 'cogs',
			, 'inventory',
			, 'expense',
			, 'liability',
			, 'asset',
			, 'discount_expense',
			, 'discount_contra',
			, 'tax_liability',
			, 'rounding',
			, 'clearing'
			, account_id
			, is_default
			, is_active
			, priority
			, conditions
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (source_type
			, (source_type
			, effective_from
		FROM pos_account_mappings
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PosAccountMappings
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SourceType,
		&entity.'product',,
		&entity.'category',,
		&entity.'paymentMethod',,
		&entity.'salesChannel',,
		&entity.'discount',,
		&entity.'rounding',,
		&entity.'tax',,
		&entity.'serviceCharge',,
		&entity.'shipping',,
		&entity.'giftCard',,
		&entity.'storeCredit',,
		&entity.'loyaltyRedemption',--,
		&entity.'default',
		&entity.SourceId,
		&entity.SourceCode,
		&entity.Purpose,
		&entity.'revenue',,
		&entity.'cogs',,
		&entity.'inventory',,
		&entity.'expense',,
		&entity.'liability',,
		&entity.'asset',,
		&entity.'discountExpense',,
		&entity.'discountContra',,
		&entity.'taxLiability',,
		&entity.'rounding',,
		&entity.'clearing',
		&entity.AccountId,
		&entity.IsDefault,
		&entity.IsActive,
		&entity.Priority,
		&entity.Conditions,
		&entity.EffectiveFrom,
		&entity.EffectiveTo,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.(sourceType,
		&entity.(sourceType,
		&entity.EffectiveFrom,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("pos_account_mappings not found")
	}

	if err != nil {
		r.logger.Error("failed to get pos_account_mappings", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get pos_account_mappings: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of pos_account_mappings records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PosAccountMappings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_account_mappings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_account_mappings
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_account_mappings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_type
			, 'product',
			, 'category',
			, 'payment_method',
			, 'sales_channel',
			, 'discount',
			, 'rounding',
			, 'tax',
			, 'service_charge',
			, 'shipping',
			, 'gift_card',
			, 'store_credit',
			, 'loyalty_redemption',--
			, 'default'
			, source_id
			, source_code
			, purpose
			, 'revenue',
			, 'cogs',
			, 'inventory',
			, 'expense',
			, 'liability',
			, 'asset',
			, 'discount_expense',
			, 'discount_contra',
			, 'tax_liability',
			, 'rounding',
			, 'clearing'
			, account_id
			, is_default
			, is_active
			, priority
			, conditions
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (source_type
			, (source_type
			, effective_from
		FROM pos_account_mappings
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_account_mappings", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_account_mappings: %w", err)
	}
	defer rows.Close()

	var entities []*PosAccountMappings
	for rows.Next() {
		var entity PosAccountMappings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceType,
			&entity.'product',,
			&entity.'category',,
			&entity.'paymentMethod',,
			&entity.'salesChannel',,
			&entity.'discount',,
			&entity.'rounding',,
			&entity.'tax',,
			&entity.'serviceCharge',,
			&entity.'shipping',,
			&entity.'giftCard',,
			&entity.'storeCredit',,
			&entity.'loyaltyRedemption',--,
			&entity.'default',
			&entity.SourceId,
			&entity.SourceCode,
			&entity.Purpose,
			&entity.'revenue',,
			&entity.'cogs',,
			&entity.'inventory',,
			&entity.'expense',,
			&entity.'liability',,
			&entity.'asset',,
			&entity.'discountExpense',,
			&entity.'discountContra',,
			&entity.'taxLiability',,
			&entity.'rounding',,
			&entity.'clearing',
			&entity.AccountId,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.Priority,
			&entity.Conditions,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.(sourceType,
			&entity.(sourceType,
			&entity.EffectiveFrom,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_account_mappings: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating pos_account_mappings rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing pos_account_mappings record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PosAccountMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_account_mappings", duration, nil)
	}()

	query := `
		UPDATE pos_account_mappings
		SET
			, organization_id = $2
			, source_type = $3
			, 'product', = $4
			, 'category', = $5
			, 'payment_method', = $6
			, 'sales_channel', = $7
			, 'discount', = $8
			, 'rounding', = $9
			, 'tax', = $10
			, 'service_charge', = $11
			, 'shipping', = $12
			, 'gift_card', = $13
			, 'store_credit', = $14
			, 'loyalty_redemption',-- = $15
			, 'default' = $16
			, source_id = $17
			, source_code = $18
			, purpose = $19
			, 'revenue', = $20
			, 'cogs', = $21
			, 'inventory', = $22
			, 'expense', = $23
			, 'liability', = $24
			, 'asset', = $25
			, 'discount_expense', = $26
			, 'discount_contra', = $27
			, 'tax_liability', = $28
			, 'rounding', = $29
			, 'clearing' = $30
			, account_id = $31
			, is_default = $32
			, is_active = $33
			, priority = $34
			, conditions = $35
			, effective_from = $36
			, effective_to = $37
			, description = $38
			, notes = $39
			, metadata = $40
			, updated_at = $42
			, deleted_at = $43
			, created_by = $44
			, updated_by = $45
			, (source_type = $46
			, (source_type = $47
			, effective_from = $48
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $49
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SourceType,
		entity.'product',,
		entity.'category',,
		entity.'paymentMethod',,
		entity.'salesChannel',,
		entity.'discount',,
		entity.'rounding',,
		entity.'tax',,
		entity.'serviceCharge',,
		entity.'shipping',,
		entity.'giftCard',,
		entity.'storeCredit',,
		entity.'loyaltyRedemption',--,
		entity.'default',
		entity.SourceId,
		entity.SourceCode,
		entity.Purpose,
		entity.'revenue',,
		entity.'cogs',,
		entity.'inventory',,
		entity.'expense',,
		entity.'liability',,
		entity.'asset',,
		entity.'discountExpense',,
		entity.'discountContra',,
		entity.'taxLiability',,
		entity.'rounding',,
		entity.'clearing',
		entity.AccountId,
		entity.IsDefault,
		entity.IsActive,
		entity.Priority,
		entity.Conditions,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.(sourceType,
		entity.(sourceType,
		entity.EffectiveFrom,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update pos_account_mappings", zap.Error(err))
		return fmt.Errorf("failed to update pos_account_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_account_mappings not found or already deleted")
	}

	r.logger.Info("updated pos_account_mappings",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a pos_account_mappings record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_account_mappings", duration, nil)
	}()

	query := `
		UPDATE pos_account_mappings
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete pos_account_mappings", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete pos_account_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_account_mappings not found or already deleted")
	}

	r.logger.Info("deleted pos_account_mappings", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves pos_account_mappings records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PosAccountMappings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_account_mappings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_account_mappings
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_account_mappings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_type
			, 'product',
			, 'category',
			, 'payment_method',
			, 'sales_channel',
			, 'discount',
			, 'rounding',
			, 'tax',
			, 'service_charge',
			, 'shipping',
			, 'gift_card',
			, 'store_credit',
			, 'loyalty_redemption',--
			, 'default'
			, source_id
			, source_code
			, purpose
			, 'revenue',
			, 'cogs',
			, 'inventory',
			, 'expense',
			, 'liability',
			, 'asset',
			, 'discount_expense',
			, 'discount_contra',
			, 'tax_liability',
			, 'rounding',
			, 'clearing'
			, account_id
			, is_default
			, is_active
			, priority
			, conditions
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (source_type
			, (source_type
			, effective_from
		FROM pos_account_mappings
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_account_mappings by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_account_mappings: %w", err)
	}
	defer rows.Close()

	var entities []*PosAccountMappings
	for rows.Next() {
		var entity PosAccountMappings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceType,
			&entity.'product',,
			&entity.'category',,
			&entity.'paymentMethod',,
			&entity.'salesChannel',,
			&entity.'discount',,
			&entity.'rounding',,
			&entity.'tax',,
			&entity.'serviceCharge',,
			&entity.'shipping',,
			&entity.'giftCard',,
			&entity.'storeCredit',,
			&entity.'loyaltyRedemption',--,
			&entity.'default',
			&entity.SourceId,
			&entity.SourceCode,
			&entity.Purpose,
			&entity.'revenue',,
			&entity.'cogs',,
			&entity.'inventory',,
			&entity.'expense',,
			&entity.'liability',,
			&entity.'asset',,
			&entity.'discountExpense',,
			&entity.'discountContra',,
			&entity.'taxLiability',,
			&entity.'rounding',,
			&entity.'clearing',
			&entity.AccountId,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.Priority,
			&entity.Conditions,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.(sourceType,
			&entity.(sourceType,
			&entity.EffectiveFrom,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_account_mappings: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


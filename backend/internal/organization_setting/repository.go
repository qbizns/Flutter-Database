package organization_setting

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

// Repository handles database operations for OrganizationSettings
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new OrganizationSettings repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// OrganizationSettings represents a organization_settings entity
type OrganizationSettings struct {
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	Timezone *string `json:"timezone" db:"timezone"`
	DateFormat *string `json:"date_format" db:"date_format"`
	TimeFormat *string `json:"time_format" db:"time_format"`
	NumberFormat *string `json:"number_format" db:"number_format"`
	DefaultCurrency *string `json:"default_currency" db:"default_currency"`
	DefaultLanguage *string `json:"default_language" db:"default_language"`
	BusinessType *string `json:"business_type" db:"business_type"`
	FiscalYearStart *string `json:"fiscal_year_start" db:"fiscal_year_start"`
	AutoPrintReceipts *bool `json:"auto_print_receipts" db:"auto_print_receipts"`
	AllowNegativeInventory *bool `json:"allow_negative_inventory" db:"allow_negative_inventory"`
	RequireCustomerForSale *bool `json:"require_customer_for_sale" db:"require_customer_for_sale"`
	EnablePriceOverride *bool `json:"enable_price_override" db:"enable_price_override"`
	AutoPostSales *bool `json:"auto_post_sales" db:"auto_post_sales"`
	AutoPostPayments *bool `json:"auto_post_payments" db:"auto_post_payments"`
	PostingFrequency *string `json:"posting_frequency" db:"posting_frequency"`
	SmtpHost *string `json:"smtp_host" db:"smtp_host"`
	SmtpPort *int64 `json:"smtp_port" db:"smtp_port"`
	SmtpUsername *string `json:"smtp_username" db:"smtp_username"`
	SmtpUseTls *bool `json:"smtp_use_tls" db:"smtp_use_tls"`
	EmailFromAddress *string `json:"email_from_address" db:"email_from_address"`
	EmailFromName *string `json:"email_from_name" db:"email_from_name"`
	EnableEmailNotifications *bool `json:"enable_email_notifications" db:"enable_email_notifications"`
	EnableSmsNotifications *bool `json:"enable_sms_notifications" db:"enable_sms_notifications"`
	Require2fa *bool `json:"require_2fa" db:"require_2fa"`
	SessionTimeoutMinutes *int64 `json:"session_timeout_minutes" db:"session_timeout_minutes"`
	PasswordMinLength *int64 `json:"password_min_length" db:"password_min_length"`
	PasswordRequireSpecial *bool `json:"password_require_special" db:"password_require_special"`
	ApiEnabled *bool `json:"api_enabled" db:"api_enabled"`
	ApiRateLimitPerMinute *int64 `json:"api_rate_limit_per_minute" db:"api_rate_limit_per_minute"`
	WebhookRetryMaxAttempts *int64 `json:"webhook_retry_max_attempts" db:"webhook_retry_max_attempts"`
	Features json.RawMessage `json:"features" db:"features"`
	CustomSettings json.RawMessage `json:"custom_settings" db:"custom_settings"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new organization_settings record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *OrganizationSettings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "organization_settings", duration, nil)
	}()

	query := `
		INSERT INTO organization_settings (
			organization_id
			, timezone
			, date_format
			, time_format
			, number_format
			, default_currency
			, default_language
			, business_type
			, fiscal_year_start
			, auto_print_receipts
			, allow_negative_inventory
			, require_customer_for_sale
			, enable_price_override
			, auto_post_sales
			, auto_post_payments
			, posting_frequency
			, smtp_host
			, smtp_port
			, smtp_username
			, smtp_use_tls
			, email_from_address
			, email_from_name
			, enable_email_notifications
			, enable_sms_notifications
			, require_2fa
			, session_timeout_minutes
			, password_min_length
			, password_require_special
			, api_enabled
			, api_rate_limit_per_minute
			, webhook_retry_max_attempts
			, features
			, custom_settings
			, updated_by
		) VALUES (
			$1
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
			, $35
		)
		RETURNING organization_id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Timezone,
		entity.DateFormat,
		entity.TimeFormat,
		entity.NumberFormat,
		entity.DefaultCurrency,
		entity.DefaultLanguage,
		entity.BusinessType,
		entity.FiscalYearStart,
		entity.AutoPrintReceipts,
		entity.AllowNegativeInventory,
		entity.RequireCustomerForSale,
		entity.EnablePriceOverride,
		entity.AutoPostSales,
		entity.AutoPostPayments,
		entity.PostingFrequency,
		entity.SmtpHost,
		entity.SmtpPort,
		entity.SmtpUsername,
		entity.SmtpUseTls,
		entity.EmailFromAddress,
		entity.EmailFromName,
		entity.EnableEmailNotifications,
		entity.EnableSmsNotifications,
		entity.Require2fa,
		entity.SessionTimeoutMinutes,
		entity.PasswordMinLength,
		entity.PasswordRequireSpecial,
		entity.ApiEnabled,
		entity.ApiRateLimitPerMinute,
		entity.WebhookRetryMaxAttempts,
		entity.Features,
		entity.CustomSettings,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.OrganizationId, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create organization_settings", zap.Error(err))
		return fmt.Errorf("failed to create organization_settings: %w", err)
	}

	r.logger.Info("created organization_settings",
		zap.String("id", entity.OrganizationId.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a organization_settings by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*OrganizationSettings, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organization_settings", duration, nil)
	}()

	query := `
		SELECT
			organization_id
			, timezone
			, date_format
			, time_format
			, number_format
			, default_currency
			, default_language
			, business_type
			, fiscal_year_start
			, auto_print_receipts
			, allow_negative_inventory
			, require_customer_for_sale
			, enable_price_override
			, auto_post_sales
			, auto_post_payments
			, posting_frequency
			, smtp_host
			, smtp_port
			, smtp_username
			, smtp_use_tls
			, email_from_address
			, email_from_name
			, enable_email_notifications
			, enable_sms_notifications
			, require_2fa
			, session_timeout_minutes
			, password_min_length
			, password_require_special
			, api_enabled
			, api_rate_limit_per_minute
			, webhook_retry_max_attempts
			, features
			, custom_settings
			, updated_at
			, updated_by
		FROM organization_settings
		WHERE organization_id = $1
		
	`

	var entity OrganizationSettings
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.OrganizationId,
		&entity.Timezone,
		&entity.DateFormat,
		&entity.TimeFormat,
		&entity.NumberFormat,
		&entity.DefaultCurrency,
		&entity.DefaultLanguage,
		&entity.BusinessType,
		&entity.FiscalYearStart,
		&entity.AutoPrintReceipts,
		&entity.AllowNegativeInventory,
		&entity.RequireCustomerForSale,
		&entity.EnablePriceOverride,
		&entity.AutoPostSales,
		&entity.AutoPostPayments,
		&entity.PostingFrequency,
		&entity.SmtpHost,
		&entity.SmtpPort,
		&entity.SmtpUsername,
		&entity.SmtpUseTls,
		&entity.EmailFromAddress,
		&entity.EmailFromName,
		&entity.EnableEmailNotifications,
		&entity.EnableSmsNotifications,
		&entity.Require2fa,
		&entity.SessionTimeoutMinutes,
		&entity.PasswordMinLength,
		&entity.PasswordRequireSpecial,
		&entity.ApiEnabled,
		&entity.ApiRateLimitPerMinute,
		&entity.WebhookRetryMaxAttempts,
		&entity.Features,
		&entity.CustomSettings,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("organization_settings not found")
	}

	if err != nil {
		r.logger.Error("failed to get organization_settings", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get organization_settings: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of organization_settings records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*OrganizationSettings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organization_settings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM organization_settings
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organization_settings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			organization_id
			, timezone
			, date_format
			, time_format
			, number_format
			, default_currency
			, default_language
			, business_type
			, fiscal_year_start
			, auto_print_receipts
			, allow_negative_inventory
			, require_customer_for_sale
			, enable_price_override
			, auto_post_sales
			, auto_post_payments
			, posting_frequency
			, smtp_host
			, smtp_port
			, smtp_username
			, smtp_use_tls
			, email_from_address
			, email_from_name
			, enable_email_notifications
			, enable_sms_notifications
			, require_2fa
			, session_timeout_minutes
			, password_min_length
			, password_require_special
			, api_enabled
			, api_rate_limit_per_minute
			, webhook_retry_max_attempts
			, features
			, custom_settings
			, updated_at
			, updated_by
		FROM organization_settings
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list organization_settings", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list organization_settings: %w", err)
	}
	defer rows.Close()

	var entities []*OrganizationSettings
	for rows.Next() {
		var entity OrganizationSettings
		err := rows.Scan(
			&entity.OrganizationId,
			&entity.Timezone,
			&entity.DateFormat,
			&entity.TimeFormat,
			&entity.NumberFormat,
			&entity.DefaultCurrency,
			&entity.DefaultLanguage,
			&entity.BusinessType,
			&entity.FiscalYearStart,
			&entity.AutoPrintReceipts,
			&entity.AllowNegativeInventory,
			&entity.RequireCustomerForSale,
			&entity.EnablePriceOverride,
			&entity.AutoPostSales,
			&entity.AutoPostPayments,
			&entity.PostingFrequency,
			&entity.SmtpHost,
			&entity.SmtpPort,
			&entity.SmtpUsername,
			&entity.SmtpUseTls,
			&entity.EmailFromAddress,
			&entity.EmailFromName,
			&entity.EnableEmailNotifications,
			&entity.EnableSmsNotifications,
			&entity.Require2fa,
			&entity.SessionTimeoutMinutes,
			&entity.PasswordMinLength,
			&entity.PasswordRequireSpecial,
			&entity.ApiEnabled,
			&entity.ApiRateLimitPerMinute,
			&entity.WebhookRetryMaxAttempts,
			&entity.Features,
			&entity.CustomSettings,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan organization_settings: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating organization_settings rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing organization_settings record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *OrganizationSettings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "organization_settings", duration, nil)
	}()

	query := `
		UPDATE organization_settings
		SET
			organization_id = $1
			, timezone = $2
			, date_format = $3
			, time_format = $4
			, number_format = $5
			, default_currency = $6
			, default_language = $7
			, business_type = $8
			, fiscal_year_start = $9
			, auto_print_receipts = $10
			, allow_negative_inventory = $11
			, require_customer_for_sale = $12
			, enable_price_override = $13
			, auto_post_sales = $14
			, auto_post_payments = $15
			, posting_frequency = $16
			, smtp_host = $17
			, smtp_port = $18
			, smtp_username = $19
			, smtp_use_tls = $20
			, email_from_address = $21
			, email_from_name = $22
			, enable_email_notifications = $23
			, enable_sms_notifications = $24
			, require_2fa = $25
			, session_timeout_minutes = $26
			, password_min_length = $27
			, password_require_special = $28
			, api_enabled = $29
			, api_rate_limit_per_minute = $30
			, webhook_retry_max_attempts = $31
			, features = $32
			, custom_settings = $33
			, updated_at = $34
			, updated_by = $35
			, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $36
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Timezone,
		entity.DateFormat,
		entity.TimeFormat,
		entity.NumberFormat,
		entity.DefaultCurrency,
		entity.DefaultLanguage,
		entity.BusinessType,
		entity.FiscalYearStart,
		entity.AutoPrintReceipts,
		entity.AllowNegativeInventory,
		entity.RequireCustomerForSale,
		entity.EnablePriceOverride,
		entity.AutoPostSales,
		entity.AutoPostPayments,
		entity.PostingFrequency,
		entity.SmtpHost,
		entity.SmtpPort,
		entity.SmtpUsername,
		entity.SmtpUseTls,
		entity.EmailFromAddress,
		entity.EmailFromName,
		entity.EnableEmailNotifications,
		entity.EnableSmsNotifications,
		entity.Require2fa,
		entity.SessionTimeoutMinutes,
		entity.PasswordMinLength,
		entity.PasswordRequireSpecial,
		entity.ApiEnabled,
		entity.ApiRateLimitPerMinute,
		entity.WebhookRetryMaxAttempts,
		entity.Features,
		entity.CustomSettings,
		time.Now(),
		entity.UpdatedBy,
		entity.OrganizationId,
	)

	if err != nil {
		r.logger.Error("failed to update organization_settings", zap.Error(err))
		return fmt.Errorf("failed to update organization_settings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("organization_settings not found or already deleted")
	}

	r.logger.Info("updated organization_settings",
		zap.String("id", entity.OrganizationId.String()),
	)

	return nil
}


// Delete permanently deletes a organization_settings record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "organization_settings", duration, nil)
	}()

	query := `DELETE FROM organization_settings WHERE organization_id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete organization_settings", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete organization_settings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("organization_settings not found")
	}

	r.logger.Info("deleted organization_settings", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves organization_settings records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*OrganizationSettings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organization_settings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM organization_settings
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organization_settings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			organization_id
			, timezone
			, date_format
			, time_format
			, number_format
			, default_currency
			, default_language
			, business_type
			, fiscal_year_start
			, auto_print_receipts
			, allow_negative_inventory
			, require_customer_for_sale
			, enable_price_override
			, auto_post_sales
			, auto_post_payments
			, posting_frequency
			, smtp_host
			, smtp_port
			, smtp_username
			, smtp_use_tls
			, email_from_address
			, email_from_name
			, enable_email_notifications
			, enable_sms_notifications
			, require_2fa
			, session_timeout_minutes
			, password_min_length
			, password_require_special
			, api_enabled
			, api_rate_limit_per_minute
			, webhook_retry_max_attempts
			, features
			, custom_settings
			, updated_at
			, updated_by
		FROM organization_settings
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list organization_settings by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list organization_settings: %w", err)
	}
	defer rows.Close()

	var entities []*OrganizationSettings
	for rows.Next() {
		var entity OrganizationSettings
		err := rows.Scan(
			&entity.OrganizationId,
			&entity.Timezone,
			&entity.DateFormat,
			&entity.TimeFormat,
			&entity.NumberFormat,
			&entity.DefaultCurrency,
			&entity.DefaultLanguage,
			&entity.BusinessType,
			&entity.FiscalYearStart,
			&entity.AutoPrintReceipts,
			&entity.AllowNegativeInventory,
			&entity.RequireCustomerForSale,
			&entity.EnablePriceOverride,
			&entity.AutoPostSales,
			&entity.AutoPostPayments,
			&entity.PostingFrequency,
			&entity.SmtpHost,
			&entity.SmtpPort,
			&entity.SmtpUsername,
			&entity.SmtpUseTls,
			&entity.EmailFromAddress,
			&entity.EmailFromName,
			&entity.EnableEmailNotifications,
			&entity.EnableSmsNotifications,
			&entity.Require2fa,
			&entity.SessionTimeoutMinutes,
			&entity.PasswordMinLength,
			&entity.PasswordRequireSpecial,
			&entity.ApiEnabled,
			&entity.ApiRateLimitPerMinute,
			&entity.WebhookRetryMaxAttempts,
			&entity.Features,
			&entity.CustomSettings,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan organization_settings: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


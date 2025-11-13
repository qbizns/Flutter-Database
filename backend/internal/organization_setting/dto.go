package organization_setting

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrganizationSettingsResponse represents a organization_settings response
type OrganizationSettingsResponse struct {
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	Timezone *string `json:"timezone"`
	
	DateFormat *string `json:"date_format"`
	
	TimeFormat *string `json:"time_format"`
	
	NumberFormat *string `json:"number_format"`
	
	DefaultCurrency *string `json:"default_currency"`
	
	DefaultLanguage *string `json:"default_language"`
	
	BusinessType *string `json:"business_type"`
	
	FiscalYearStart *string `json:"fiscal_year_start"`
	
	AutoPrintReceipts *bool `json:"auto_print_receipts"`
	
	AllowNegativeInventory *bool `json:"allow_negative_inventory"`
	
	RequireCustomerForSale *bool `json:"require_customer_for_sale"`
	
	EnablePriceOverride *bool `json:"enable_price_override"`
	
	AutoPostSales *bool `json:"auto_post_sales"`
	
	AutoPostPayments *bool `json:"auto_post_payments"`
	
	PostingFrequency *string `json:"posting_frequency"`
	
	SmtpHost *string `json:"smtp_host"`
	
	SmtpPort *int64 `json:"smtp_port"`
	
	SmtpUsername *string `json:"smtp_username"`
	
	SmtpUseTls *bool `json:"smtp_use_tls"`
	
	EmailFromAddress *string `json:"email_from_address"`
	
	EmailFromName *string `json:"email_from_name"`
	
	EnableEmailNotifications *bool `json:"enable_email_notifications"`
	
	EnableSmsNotifications *bool `json:"enable_sms_notifications"`
	
	Require2fa *bool `json:"require_2fa"`
	
	SessionTimeoutMinutes *int64 `json:"session_timeout_minutes"`
	
	PasswordMinLength *int64 `json:"password_min_length"`
	
	PasswordRequireSpecial *bool `json:"password_require_special"`
	
	ApiEnabled *bool `json:"api_enabled"`
	
	ApiRateLimitPerMinute *int64 `json:"api_rate_limit_per_minute"`
	
	WebhookRetryMaxAttempts *int64 `json:"webhook_retry_max_attempts"`
	
	Features json.RawMessage `json:"features"`
	
	CustomSettings json.RawMessage `json:"custom_settings"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateOrganizationSettingsRequest represents a request to create a organization_settings
type CreateOrganizationSettingsRequest struct {
	
	Timezone *string `json:"timezone"`
	
	DateFormat *string `json:"date_format"`
	
	TimeFormat *string `json:"time_format"`
	
	NumberFormat *string `json:"number_format"`
	
	DefaultCurrency *string `json:"default_currency"`
	
	DefaultLanguage *string `json:"default_language"`
	
	BusinessType *string `json:"business_type"`
	
	FiscalYearStart *string `json:"fiscal_year_start"`
	
	AutoPrintReceipts *bool `json:"auto_print_receipts"`
	
	AllowNegativeInventory *bool `json:"allow_negative_inventory"`
	
	RequireCustomerForSale *bool `json:"require_customer_for_sale"`
	
	EnablePriceOverride *bool `json:"enable_price_override"`
	
	AutoPostSales *bool `json:"auto_post_sales"`
	
	AutoPostPayments *bool `json:"auto_post_payments"`
	
	PostingFrequency *string `json:"posting_frequency"`
	
	SmtpHost *string `json:"smtp_host"`
	
	SmtpPort *int64 `json:"smtp_port"`
	
	SmtpUsername *string `json:"smtp_username"`
	
	SmtpUseTls *bool `json:"smtp_use_tls"`
	
	EmailFromAddress *string `json:"email_from_address" validate:"email"`
	
	EmailFromName *string `json:"email_from_name" validate:"email"`
	
	EnableEmailNotifications *bool `json:"enable_email_notifications"`
	
	EnableSmsNotifications *bool `json:"enable_sms_notifications"`
	
	Require2fa *bool `json:"require_2fa"`
	
	SessionTimeoutMinutes *int64 `json:"session_timeout_minutes"`
	
	PasswordMinLength *int64 `json:"password_min_length"`
	
	PasswordRequireSpecial *bool `json:"password_require_special"`
	
	ApiEnabled *bool `json:"api_enabled"`
	
	ApiRateLimitPerMinute *int64 `json:"api_rate_limit_per_minute"`
	
	WebhookRetryMaxAttempts *int64 `json:"webhook_retry_max_attempts"`
	
	// Duplicate removed: Features json.RawMessage `json:"features"`
	
	// Duplicate removed: CustomSettings json.RawMessage `json:"custom_settings"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateOrganizationSettingsRequest) Validate() error {
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateOrganizationSettingsRequest represents a request to update a organization_settings
type UpdateOrganizationSettingsRequest struct {
	
	Timezone *string `json:"timezone,omitempty"`
	
	DateFormat *string `json:"date_format,omitempty"`
	
	TimeFormat *string `json:"time_format,omitempty"`
	
	NumberFormat *string `json:"number_format,omitempty"`
	
	DefaultCurrency *string `json:"default_currency,omitempty"`
	
	DefaultLanguage *string `json:"default_language,omitempty"`
	
	BusinessType *string `json:"business_type,omitempty"`
	
	FiscalYearStart *string `json:"fiscal_year_start,omitempty"`
	
	AutoPrintReceipts *bool `json:"auto_print_receipts,omitempty"`
	
	AllowNegativeInventory *bool `json:"allow_negative_inventory,omitempty"`
	
	RequireCustomerForSale *bool `json:"require_customer_for_sale,omitempty"`
	
	EnablePriceOverride *bool `json:"enable_price_override,omitempty"`
	
	AutoPostSales *bool `json:"auto_post_sales,omitempty"`
	
	AutoPostPayments *bool `json:"auto_post_payments,omitempty"`
	
	PostingFrequency *string `json:"posting_frequency,omitempty"`
	
	SmtpHost *string `json:"smtp_host,omitempty"`
	
	SmtpPort *int64 `json:"smtp_port,omitempty"`
	
	SmtpUsername *string `json:"smtp_username,omitempty"`
	
	SmtpUseTls *bool `json:"smtp_use_tls,omitempty"`
	
	EmailFromAddress *string `json:"email_from_address,omitempty" validate:"omitempty,email"`
	
	EmailFromName *string `json:"email_from_name,omitempty" validate:"omitempty,email"`
	
	EnableEmailNotifications *bool `json:"enable_email_notifications,omitempty"`
	
	EnableSmsNotifications *bool `json:"enable_sms_notifications,omitempty"`
	
	Require2fa *bool `json:"require_2fa,omitempty"`
	
	SessionTimeoutMinutes *int64 `json:"session_timeout_minutes,omitempty"`
	
	PasswordMinLength *int64 `json:"password_min_length,omitempty"`
	
	PasswordRequireSpecial *bool `json:"password_require_special,omitempty"`
	
	ApiEnabled *bool `json:"api_enabled,omitempty"`
	
	ApiRateLimitPerMinute *int64 `json:"api_rate_limit_per_minute,omitempty"`
	
	WebhookRetryMaxAttempts *int64 `json:"webhook_retry_max_attempts,omitempty"`
	
	Features *json.RawMessage `json:"features,omitempty"`
	
	CustomSettings *json.RawMessage `json:"custom_settings,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateOrganizationSettingsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Timezone != nil {
		hasUpdate = true
	}
	
	if r.DateFormat != nil {
		hasUpdate = true
	}
	
	if r.TimeFormat != nil {
		hasUpdate = true
	}
	
	if r.NumberFormat != nil {
		hasUpdate = true
	}
	
	if r.DefaultCurrency != nil {
		hasUpdate = true
	}
	
	if r.DefaultLanguage != nil {
		hasUpdate = true
	}
	
	if r.BusinessType != nil {
		hasUpdate = true
	}
	
	if r.FiscalYearStart != nil {
		hasUpdate = true
	}
	
	if r.AutoPrintReceipts != nil {
		hasUpdate = true
	}
	
	if r.AllowNegativeInventory != nil {
		hasUpdate = true
	}
	
	if r.RequireCustomerForSale != nil {
		hasUpdate = true
	}
	
	if r.EnablePriceOverride != nil {
		hasUpdate = true
	}
	
	if r.AutoPostSales != nil {
		hasUpdate = true
	}
	
	if r.AutoPostPayments != nil {
		hasUpdate = true
	}
	
	if r.PostingFrequency != nil {
		hasUpdate = true
	}
	
	if r.SmtpHost != nil {
		hasUpdate = true
	}
	
	if r.SmtpPort != nil {
		hasUpdate = true
	}
	
	if r.SmtpUsername != nil {
		hasUpdate = true
	}
	
	if r.SmtpUseTls != nil {
		hasUpdate = true
	}
	
	if r.EmailFromAddress != nil {
		hasUpdate = true
	}
	
	if r.EmailFromName != nil {
		hasUpdate = true
	}
	
	if r.EnableEmailNotifications != nil {
		hasUpdate = true
	}
	
	if r.EnableSmsNotifications != nil {
		hasUpdate = true
	}
	
	if r.Require2fa != nil {
		hasUpdate = true
	}
	
	if r.SessionTimeoutMinutes != nil {
		hasUpdate = true
	}
	
	if r.PasswordMinLength != nil {
		hasUpdate = true
	}
	
	if r.PasswordRequireSpecial != nil {
		hasUpdate = true
	}
	
	if r.ApiEnabled != nil {
		hasUpdate = true
	}
	
	if r.ApiRateLimitPerMinute != nil {
		hasUpdate = true
	}
	
	if r.WebhookRetryMaxAttempts != nil {
		hasUpdate = true
	}
	
	if r.Features != nil {
		hasUpdate = true
	}
	
	if r.CustomSettings != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// OrganizationSettingsListResponse represents a paginated list of organization_settings records
type OrganizationSettingsListResponse struct {
	Items      []*OrganizationSettingsResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

-- ============================================================================
-- Migration V023: Backend Infrastructure Tables
-- Description: Essential tables for production-grade Go backend API
-- Created: 2025-11-10
-- Dependencies: V001, V002, V003
-- ============================================================================

-- This migration adds critical infrastructure tables for:
-- - Background job processing
-- - API authentication & rate limiting
-- - Webhooks & integrations
-- - Notifications & messaging
-- - File attachments
-- - Settings & preferences
-- - Audit & monitoring

BEGIN;

-- ============================================================================
-- BACKGROUND JOBS QUEUE
-- ============================================================================

-- Job queue for async processing (posting engine, reports, exports, etc.)
CREATE TABLE IF NOT EXISTS background_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,

    -- Job Identification
    job_type VARCHAR(100) NOT NULL,  -- 'posting_engine', 'report_generation', 'export', 'email', 'import'
    job_name VARCHAR(255) NOT NULL,
    queue_name VARCHAR(50) DEFAULT 'default',  -- 'default', 'high_priority', 'low_priority'

    -- Job Status
    status VARCHAR(30) DEFAULT 'pending' CHECK (
        status IN ('pending', 'processing', 'completed', 'failed', 'cancelled', 'retrying')
    ),

    -- Job Data
    payload JSONB NOT NULL DEFAULT '{}',  -- Input data
    result JSONB,                          -- Output data
    error_message TEXT,
    error_details JSONB,

    -- Processing Info
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    priority INTEGER DEFAULT 0,  -- Higher = more priority

    -- Timing
    scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,

    -- Worker Info
    worker_id VARCHAR(100),  -- Which worker picked up the job
    processing_timeout INTEGER DEFAULT 300,  -- Seconds

    -- Metadata
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_background_jobs_organization_id ON background_jobs(organization_id);
CREATE INDEX idx_background_jobs_status ON background_jobs(status) WHERE status IN ('pending', 'retrying');
CREATE INDEX idx_background_jobs_job_type ON background_jobs(job_type);
CREATE INDEX idx_background_jobs_queue_name ON background_jobs(queue_name);
CREATE INDEX idx_background_jobs_scheduled_at ON background_jobs(scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_background_jobs_created_at ON background_jobs(created_at);

COMMENT ON TABLE background_jobs IS 'Async job queue for background processing';

-- ============================================================================
-- API KEYS
-- ============================================================================

-- API keys for programmatic access
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Key Identification
    key_name VARCHAR(255) NOT NULL,
    key_prefix VARCHAR(20) NOT NULL,  -- First 8 chars of key for display (e.g., 'sk_live_')
    key_hash VARCHAR(255) NOT NULL,   -- Hashed full key

    -- Permissions
    scopes JSONB DEFAULT '["read"]',  -- ['read', 'write', 'delete', 'admin']
    allowed_ips TEXT[],                -- Whitelist IPs

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Usage Tracking
    last_used_at TIMESTAMP WITH TIME ZONE,
    usage_count INTEGER DEFAULT 0,

    -- Rate Limiting
    rate_limit_per_minute INTEGER DEFAULT 60,
    rate_limit_per_hour INTEGER DEFAULT 1000,

    -- Expiration
    expires_at TIMESTAMP WITH TIME ZONE,

    -- Audit
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_key_hash UNIQUE (key_hash)
);

CREATE INDEX idx_api_keys_organization_id ON api_keys(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_api_keys_key_prefix ON api_keys(key_prefix);
CREATE INDEX idx_api_keys_is_active ON api_keys(is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE api_keys IS 'API keys for programmatic access to the platform';

-- ============================================================================
-- WEBHOOKS
-- ============================================================================

-- Webhook configurations
CREATE TABLE IF NOT EXISTS webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Webhook Configuration
    webhook_name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    secret VARCHAR(255),  -- For signature verification

    -- Events to Subscribe
    events TEXT[] NOT NULL,  -- ['sale.created', 'payment.completed', 'invoice.posted']

    -- HTTP Configuration
    http_method VARCHAR(10) DEFAULT 'POST',
    headers JSONB DEFAULT '{}',
    timeout_seconds INTEGER DEFAULT 30,

    -- Retry Configuration
    max_retries INTEGER DEFAULT 3,
    retry_backoff_seconds INTEGER DEFAULT 60,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,

    -- Statistics
    total_deliveries INTEGER DEFAULT 0,
    successful_deliveries INTEGER DEFAULT 0,
    failed_deliveries INTEGER DEFAULT 0,
    last_delivery_at TIMESTAMP WITH TIME ZONE,
    last_success_at TIMESTAMP WITH TIME ZONE,
    last_failure_at TIMESTAMP WITH TIME ZONE,

    -- Audit
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_webhooks_organization_id ON webhooks(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_webhooks_is_active ON webhooks(is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE webhooks IS 'Outbound webhook configurations for event notifications';

-- ============================================================================
-- WEBHOOK DELIVERIES
-- ============================================================================

-- Webhook delivery log
CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    webhook_id UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,

    -- Event Info
    event_type VARCHAR(100) NOT NULL,
    event_id UUID NOT NULL,  -- ID of the source event (sale_id, payment_id, etc.)

    -- Delivery Info
    status VARCHAR(30) DEFAULT 'pending' CHECK (
        status IN ('pending', 'success', 'failed', 'retrying')
    ),

    -- Request/Response
    request_url TEXT NOT NULL,
    request_method VARCHAR(10) NOT NULL,
    request_headers JSONB,
    request_body JSONB,
    response_status_code INTEGER,
    response_headers JSONB,
    response_body TEXT,

    -- Timing
    attempt_number INTEGER DEFAULT 1,
    duration_ms INTEGER,
    next_retry_at TIMESTAMP WITH TIME ZONE,

    -- Error Info
    error_message TEXT,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_organization_id ON webhook_deliveries(organization_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_event_type ON webhook_deliveries(event_type);
CREATE INDEX idx_webhook_deliveries_created_at ON webhook_deliveries(created_at);
CREATE INDEX idx_webhook_deliveries_next_retry ON webhook_deliveries(next_retry_at) WHERE status = 'retrying';

COMMENT ON TABLE webhook_deliveries IS 'Log of webhook delivery attempts';

-- ============================================================================
-- NOTIFICATIONS
-- ============================================================================

-- User notifications (in-app, email, SMS, push)
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Notification Content
    notification_type VARCHAR(50) NOT NULL,  -- 'info', 'success', 'warning', 'error'
    category VARCHAR(50) NOT NULL,           -- 'sale', 'payment', 'inventory', 'system'
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,

    -- Actions
    action_url VARCHAR(500),
    action_label VARCHAR(100),

    -- Channel
    channels TEXT[] DEFAULT ARRAY['in_app'],  -- ['in_app', 'email', 'sms', 'push']

    -- Status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,

    -- Related Entity
    related_entity_type VARCHAR(100),  -- 'sale', 'payment', 'invoice', etc.
    related_entity_id UUID,

    -- Priority
    priority VARCHAR(20) DEFAULT 'normal',  -- 'low', 'normal', 'high', 'urgent'

    -- Expiration
    expires_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_organization_id ON notifications(organization_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_category ON notifications(category);

COMMENT ON TABLE notifications IS 'User notifications across multiple channels';

-- ============================================================================
-- NOTIFICATION PREFERENCES
-- ============================================================================

-- User notification preferences
CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Category Preferences
    category VARCHAR(50) NOT NULL,  -- 'sales', 'inventory', 'accounting', 'system'

    -- Channel Preferences
    in_app_enabled BOOLEAN DEFAULT TRUE,
    email_enabled BOOLEAN DEFAULT TRUE,
    sms_enabled BOOLEAN DEFAULT FALSE,
    push_enabled BOOLEAN DEFAULT TRUE,

    -- Frequency
    frequency VARCHAR(20) DEFAULT 'realtime',  -- 'realtime', 'daily_digest', 'weekly_digest', 'disabled'

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_user_category UNIQUE (user_id, category)
);

CREATE INDEX idx_notification_preferences_user_id ON notification_preferences(user_id);
CREATE INDEX idx_notification_preferences_organization_id ON notification_preferences(organization_id);

COMMENT ON TABLE notification_preferences IS 'User notification channel and frequency preferences';

-- ============================================================================
-- FILE ATTACHMENTS
-- ============================================================================

-- File attachments for documents (receipts, invoices, etc.)
CREATE TABLE IF NOT EXISTS file_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- File Info
    file_name VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,  -- Bytes
    mime_type VARCHAR(100) NOT NULL,
    file_extension VARCHAR(20),

    -- Storage
    storage_provider VARCHAR(50) DEFAULT 'local',  -- 'local', 's3', 'gcs', 'azure'
    storage_path TEXT NOT NULL,
    storage_url TEXT,

    -- File Hash (for deduplication)
    file_hash VARCHAR(64),  -- SHA-256

    -- Related Entity
    entity_type VARCHAR(100) NOT NULL,  -- 'sale', 'invoice', 'expense', 'product', etc.
    entity_id UUID NOT NULL,

    -- Metadata
    description TEXT,
    tags TEXT[],
    is_public BOOLEAN DEFAULT FALSE,

    -- Image Metadata (if applicable)
    image_width INTEGER,
    image_height INTEGER,

    -- Virus Scan
    virus_scan_status VARCHAR(20) DEFAULT 'pending',  -- 'pending', 'clean', 'infected', 'error'
    virus_scan_at TIMESTAMP WITH TIME ZONE,

    -- Audit
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_file_attachments_organization_id ON file_attachments(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_file_attachments_entity ON file_attachments(entity_type, entity_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_file_attachments_uploaded_by ON file_attachments(uploaded_by);
CREATE INDEX idx_file_attachments_file_hash ON file_attachments(file_hash);
CREATE INDEX idx_file_attachments_created_at ON file_attachments(created_at);

COMMENT ON TABLE file_attachments IS 'File attachments for documents and entities';

-- ============================================================================
-- EMAIL QUEUE
-- ============================================================================

-- Email sending queue
CREATE TABLE IF NOT EXISTS email_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,

    -- Recipient Info
    to_addresses TEXT[] NOT NULL,
    cc_addresses TEXT[],
    bcc_addresses TEXT[],
    from_address VARCHAR(255),
    reply_to VARCHAR(255),

    -- Content
    subject VARCHAR(998) NOT NULL,  -- RFC 2822 limit
    body_html TEXT,
    body_text TEXT,

    -- Attachments
    attachment_ids UUID[],  -- References file_attachments(id)

    -- Template (optional)
    template_name VARCHAR(100),
    template_data JSONB,

    -- Status
    status VARCHAR(30) DEFAULT 'pending' CHECK (
        status IN ('pending', 'sending', 'sent', 'failed', 'cancelled')
    ),

    -- Delivery Info
    provider VARCHAR(50),  -- 'smtp', 'sendgrid', 'ses', 'mailgun'
    provider_message_id VARCHAR(255),
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,

    -- Error Info
    error_message TEXT,

    -- Priority
    priority INTEGER DEFAULT 0,

    -- Timestamps
    scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_queue_status ON email_queue(status);
CREATE INDEX idx_email_queue_organization_id ON email_queue(organization_id);
CREATE INDEX idx_email_queue_scheduled_at ON email_queue(scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_email_queue_created_at ON email_queue(created_at);

COMMENT ON TABLE email_queue IS 'Queue for sending emails asynchronously';

-- ============================================================================
-- SMS QUEUE
-- ============================================================================

-- SMS sending queue
CREATE TABLE IF NOT EXISTS sms_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,

    -- Recipient Info
    to_phone VARCHAR(20) NOT NULL,
    from_phone VARCHAR(20),

    -- Content
    message TEXT NOT NULL,

    -- Status
    status VARCHAR(30) DEFAULT 'pending' CHECK (
        status IN ('pending', 'sending', 'sent', 'failed', 'cancelled')
    ),

    -- Delivery Info
    provider VARCHAR(50),  -- 'twilio', 'nexmo', 'sns'
    provider_message_id VARCHAR(255),
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,

    -- Error Info
    error_message TEXT,

    -- Cost
    cost_amount NUMERIC(10, 4),
    cost_currency VARCHAR(3),

    -- Timestamps
    scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sms_queue_status ON sms_queue(status);
CREATE INDEX idx_sms_queue_organization_id ON sms_queue(organization_id);
CREATE INDEX idx_sms_queue_scheduled_at ON sms_queue(scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_sms_queue_created_at ON sms_queue(created_at);

COMMENT ON TABLE sms_queue IS 'Queue for sending SMS messages asynchronously';

-- ============================================================================
-- RATE LIMITS
-- ============================================================================

-- API rate limiting tracking
CREATE TABLE IF NOT EXISTS rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Identifier (user, org, API key, or IP)
    identifier_type VARCHAR(20) NOT NULL,  -- 'user', 'organization', 'api_key', 'ip'
    identifier_value VARCHAR(255) NOT NULL,

    -- Endpoint
    endpoint_path VARCHAR(500),  -- NULL = global limit
    http_method VARCHAR(10),

    -- Window
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    window_duration_seconds INTEGER NOT NULL,  -- 60 (per minute), 3600 (per hour)

    -- Counts
    request_count INTEGER DEFAULT 1,
    allowed_count INTEGER NOT NULL,

    -- Status
    is_blocked BOOLEAN DEFAULT FALSE,
    blocked_until TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    first_request_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_request_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_rate_limit_window UNIQUE (
        identifier_type, identifier_value, endpoint_path, http_method, window_start
    )
);

CREATE INDEX idx_rate_limits_identifier ON rate_limits(identifier_type, identifier_value);
CREATE INDEX idx_rate_limits_window_start ON rate_limits(window_start);
CREATE INDEX idx_rate_limits_is_blocked ON rate_limits(is_blocked) WHERE is_blocked = TRUE;

COMMENT ON TABLE rate_limits IS 'API rate limiting tracking per user/org/key/IP';

-- ============================================================================
-- USER SESSIONS
-- ============================================================================

-- User session management (beyond JWT)
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,

    -- Session Info
    session_token VARCHAR(255) NOT NULL UNIQUE,
    refresh_token VARCHAR(255) UNIQUE,

    -- Device Info
    user_agent TEXT,
    ip_address INET,
    device_type VARCHAR(50),  -- 'desktop', 'mobile', 'tablet'
    device_name VARCHAR(255),
    browser VARCHAR(100),
    os VARCHAR(100),

    -- Location (GeoIP)
    country_code VARCHAR(2),
    city VARCHAR(100),

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Timestamps
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_session_token ON user_sessions(session_token) WHERE is_active = TRUE;
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_organization_id ON user_sessions(organization_id);

COMMENT ON TABLE user_sessions IS 'User session tracking for authentication and device management';

-- ============================================================================
-- ORGANIZATION SETTINGS
-- ============================================================================

-- Organization-specific settings
CREATE TABLE IF NOT EXISTS organization_settings (
    organization_id UUID PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,

    -- General Settings
    timezone VARCHAR(50) DEFAULT 'UTC',
    date_format VARCHAR(20) DEFAULT 'YYYY-MM-DD',
    time_format VARCHAR(20) DEFAULT 'HH:mm:ss',
    number_format VARCHAR(20) DEFAULT 'en_US',

    -- Currency & Locale
    default_currency VARCHAR(3) DEFAULT 'USD',
    default_language VARCHAR(10) DEFAULT 'en',

    -- Business Settings
    business_type VARCHAR(50),  -- 'retail', 'restaurant', 'wholesale', 'service'
    fiscal_year_start VARCHAR(5) DEFAULT '01-01',  -- MM-DD

    -- POS Settings
    auto_print_receipts BOOLEAN DEFAULT FALSE,
    allow_negative_inventory BOOLEAN DEFAULT FALSE,
    require_customer_for_sale BOOLEAN DEFAULT FALSE,
    enable_price_override BOOLEAN DEFAULT FALSE,

    -- Accounting Settings
    auto_post_sales BOOLEAN DEFAULT FALSE,
    auto_post_payments BOOLEAN DEFAULT FALSE,
    posting_frequency VARCHAR(20) DEFAULT 'realtime',  -- 'realtime', 'daily', 'manual'

    -- Email Settings
    smtp_host VARCHAR(255),
    smtp_port INTEGER,
    smtp_username VARCHAR(255),
    smtp_use_tls BOOLEAN DEFAULT TRUE,
    email_from_address VARCHAR(255),
    email_from_name VARCHAR(255),

    -- Notification Settings
    enable_email_notifications BOOLEAN DEFAULT TRUE,
    enable_sms_notifications BOOLEAN DEFAULT FALSE,

    -- Security Settings
    require_2fa BOOLEAN DEFAULT FALSE,
    session_timeout_minutes INTEGER DEFAULT 480,
    password_min_length INTEGER DEFAULT 8,
    password_require_special BOOLEAN DEFAULT TRUE,

    -- API Settings
    api_enabled BOOLEAN DEFAULT TRUE,
    api_rate_limit_per_minute INTEGER DEFAULT 60,
    webhook_retry_max_attempts INTEGER DEFAULT 3,

    -- Feature Flags (override organization_features)
    features JSONB DEFAULT '{}',

    -- Custom Settings (extensible)
    custom_settings JSONB DEFAULT '{}',

    -- Audit
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL
);

COMMENT ON TABLE organization_settings IS 'Organization-specific configuration and preferences';

-- ============================================================================
-- USER SETTINGS
-- ============================================================================

-- User-specific settings and preferences
CREATE TABLE IF NOT EXISTS user_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- Display Preferences
    theme VARCHAR(20) DEFAULT 'light',  -- 'light', 'dark', 'auto'
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50),

    -- Dashboard Preferences
    default_dashboard VARCHAR(50) DEFAULT 'overview',
    dashboard_layout JSONB DEFAULT '{}',

    -- Table/Grid Preferences
    items_per_page INTEGER DEFAULT 25,
    default_view VARCHAR(20) DEFAULT 'grid',  -- 'grid', 'list', 'table'

    -- Notification Preferences (global)
    desktop_notifications BOOLEAN DEFAULT TRUE,
    sound_notifications BOOLEAN DEFAULT TRUE,

    -- POS Preferences
    default_location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    quick_actions JSONB DEFAULT '[]',  -- Customizable quick actions

    -- Custom Preferences (extensible)
    custom_preferences JSONB DEFAULT '{}',

    -- Audit
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE user_settings IS 'User-specific preferences and settings';

-- ============================================================================
-- DATA EXPORT REQUESTS
-- ============================================================================

-- Data export job tracking
CREATE TABLE IF NOT EXISTS data_export_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Request Info
    export_type VARCHAR(50) NOT NULL,  -- 'sales', 'products', 'customers', 'full_backup'
    export_format VARCHAR(20) NOT NULL,  -- 'csv', 'xlsx', 'json', 'pdf'

    -- Filters
    date_from DATE,
    date_to DATE,
    filters JSONB DEFAULT '{}',

    -- Status
    status VARCHAR(30) DEFAULT 'pending' CHECK (
        status IN ('pending', 'processing', 'completed', 'failed', 'expired')
    ),

    -- File Info
    file_name VARCHAR(500),
    file_size BIGINT,
    file_path TEXT,
    download_url TEXT,
    download_expires_at TIMESTAMP WITH TIME ZONE,

    -- Processing Info
    total_records INTEGER,
    processed_records INTEGER,
    error_message TEXT,

    -- Timestamps
    requested_by UUID REFERENCES users(id) ON DELETE SET NULL,
    requested_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_data_export_requests_organization_id ON data_export_requests(organization_id);
CREATE INDEX idx_data_export_requests_status ON data_export_requests(status);
CREATE INDEX idx_data_export_requests_requested_by ON data_export_requests(requested_by);
CREATE INDEX idx_data_export_requests_requested_at ON data_export_requests(requested_at);

COMMENT ON TABLE data_export_requests IS 'User data export requests and job tracking';

-- ============================================================================
-- SCHEDULED REPORTS
-- ============================================================================

-- Scheduled report generation
CREATE TABLE IF NOT EXISTS scheduled_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Report Configuration
    report_name VARCHAR(255) NOT NULL,
    report_type VARCHAR(100) NOT NULL,  -- 'sales_summary', 'inventory_valuation', 'balance_sheet'

    -- Schedule
    schedule_frequency VARCHAR(20) NOT NULL,  -- 'daily', 'weekly', 'monthly', 'quarterly'
    schedule_day_of_week INTEGER,  -- 0-6 (Sunday-Saturday)
    schedule_day_of_month INTEGER,  -- 1-31
    schedule_time TIME NOT NULL DEFAULT '09:00:00',
    schedule_timezone VARCHAR(50) DEFAULT 'UTC',

    -- Filters & Parameters
    report_parameters JSONB DEFAULT '{}',

    -- Delivery
    delivery_method VARCHAR(20) DEFAULT 'email',  -- 'email', 'webhook', 'sftp'
    delivery_recipients TEXT[],

    -- Format
    output_format VARCHAR(20) DEFAULT 'pdf',  -- 'pdf', 'xlsx', 'csv'

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Last Run
    last_run_at TIMESTAMP WITH TIME ZONE,
    last_run_status VARCHAR(20),  -- 'success', 'failed'
    next_run_at TIMESTAMP WITH TIME ZONE,

    -- Audit
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_scheduled_reports_organization_id ON scheduled_reports(organization_id);
CREATE INDEX idx_scheduled_reports_next_run ON scheduled_reports(next_run_at) WHERE is_active = TRUE;
CREATE INDEX idx_scheduled_reports_is_active ON scheduled_reports(is_active);

COMMENT ON TABLE scheduled_reports IS 'Automated scheduled report generation and delivery';

-- ============================================================================
-- API REQUEST LOGS
-- ============================================================================

-- API request/response logging (for debugging and monitoring)
CREATE TABLE IF NOT EXISTS api_request_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Request Info
    request_id VARCHAR(100) UNIQUE,  -- Correlation ID
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    query_params JSONB,

    -- Authentication
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    api_key_id UUID REFERENCES api_keys(id) ON DELETE SET NULL,

    -- Request Details
    request_headers JSONB,
    request_body JSONB,
    ip_address INET,
    user_agent TEXT,

    -- Response Details
    status_code INTEGER NOT NULL,
    response_headers JSONB,
    response_body JSONB,

    -- Performance
    duration_ms INTEGER,

    -- Error Info
    error_message TEXT,
    error_stack TEXT,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Partition by month for performance
CREATE INDEX idx_api_request_logs_created_at ON api_request_logs(created_at);
CREATE INDEX idx_api_request_logs_user_id ON api_request_logs(user_id);
CREATE INDEX idx_api_request_logs_organization_id ON api_request_logs(organization_id);
CREATE INDEX idx_api_request_logs_request_id ON api_request_logs(request_id);
CREATE INDEX idx_api_request_logs_status_code ON api_request_logs(status_code);
CREATE INDEX idx_api_request_logs_path ON api_request_logs(path);

COMMENT ON TABLE api_request_logs IS 'API request/response logs for monitoring and debugging';

-- ============================================================================
-- INTEGRATION CONFIGS
-- ============================================================================

-- Third-party integration configurations
CREATE TABLE IF NOT EXISTS integration_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Integration Info
    integration_type VARCHAR(50) NOT NULL,  -- 'payment_gateway', 'shipping', 'accounting', 'crm'
    provider_name VARCHAR(100) NOT NULL,    -- 'stripe', 'paypal', 'quickbooks', 'xero'

    -- Configuration
    credentials JSONB NOT NULL,  -- Encrypted credentials
    settings JSONB DEFAULT '{}',

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_connected BOOLEAN DEFAULT FALSE,
    connection_status VARCHAR(20),  -- 'connected', 'error', 'pending'

    -- Sync Info
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_sync_status VARCHAR(20),
    sync_frequency VARCHAR(20) DEFAULT 'manual',  -- 'manual', 'hourly', 'daily'

    -- Webhooks
    webhook_url TEXT,
    webhook_secret VARCHAR(255),

    -- Audit
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_integration_configs_organization_id ON integration_configs(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_integration_configs_integration_type ON integration_configs(integration_type);
CREATE INDEX idx_integration_configs_is_active ON integration_configs(is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE integration_configs IS 'Third-party integration configurations and credentials';

-- ============================================================================
-- AUTO-UPDATE TRIGGERS
-- ============================================================================

CREATE TRIGGER update_background_jobs_updated_at
    BEFORE UPDATE ON background_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_api_keys_updated_at
    BEFORE UPDATE ON api_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhooks_updated_at
    BEFORE UPDATE ON webhooks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notification_preferences_updated_at
    BEFORE UPDATE ON notification_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_organization_settings_updated_at
    BEFORE UPDATE ON organization_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at
    BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_scheduled_reports_updated_at
    BEFORE UPDATE ON scheduled_reports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_integration_configs_updated_at
    BEFORE UPDATE ON integration_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;

-- ============================================================================
-- SUMMARY
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V023 completed successfully!';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Added 17 infrastructure tables:';
    RAISE NOTICE '  - background_jobs (async job queue)';
    RAISE NOTICE '  - api_keys (API authentication)';
    RAISE NOTICE '  - webhooks (outbound webhooks)';
    RAISE NOTICE '  - webhook_deliveries (webhook logs)';
    RAISE NOTICE '  - notifications (user notifications)';
    RAISE NOTICE '  - notification_preferences';
    RAISE NOTICE '  - file_attachments (document files)';
    RAISE NOTICE '  - email_queue (email sending)';
    RAISE NOTICE '  - sms_queue (SMS sending)';
    RAISE NOTICE '  - rate_limits (API rate limiting)';
    RAISE NOTICE '  - user_sessions (session management)';
    RAISE NOTICE '  - organization_settings';
    RAISE NOTICE '  - user_settings';
    RAISE NOTICE '  - data_export_requests';
    RAISE NOTICE '  - scheduled_reports';
    RAISE NOTICE '  - api_request_logs';
    RAISE NOTICE '  - integration_configs';
    RAISE NOTICE '============================================';
END $$;

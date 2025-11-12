package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ============================================================================
// BACKGROUND JOBS - Async job processing queue
// ============================================================================

// BackgroundJob represents a background job in the queue
type BackgroundJob struct {
	ID                uuid.UUID        `json:"id" db:"id"`
	OrganizationID    uuid.UUID        `json:"organization_id" db:"organization_id"`
	JobType           string           `json:"job_type" db:"job_type"`
	JobName           string           `json:"job_name" db:"job_name"`
	QueueName         string           `json:"queue_name" db:"queue_name"`
	Status            string           `json:"status" db:"status"`
	Payload           json.RawMessage  `json:"payload" db:"payload"`
	Result            *json.RawMessage `json:"result" db:"result"`
	ErrorMessage      *string          `json:"error_message" db:"error_message"`
	ErrorDetails      *json.RawMessage `json:"error_details" db:"error_details"`
	Attempts          int              `json:"attempts" db:"attempts"`
	MaxAttempts       int              `json:"max_attempts" db:"max_attempts"`
	Priority          int              `json:"priority" db:"priority"`
	ScheduledAt       time.Time        `json:"scheduled_at" db:"scheduled_at"`
	StartedAt         *time.Time       `json:"started_at" db:"started_at"`
	CompletedAt       *time.Time       `json:"completed_at" db:"completed_at"`
	FailedAt          *time.Time       `json:"failed_at" db:"failed_at"`
	WorkerID          *string          `json:"worker_id" db:"worker_id"`
	ProcessingTimeout int              `json:"processing_timeout" db:"processing_timeout"`
	CreatedBy         *uuid.UUID       `json:"created_by" db:"created_by"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// API KEYS - Programmatic access authentication
// ============================================================================

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID                 uuid.UUID      `json:"id" db:"id"`
	OrganizationID     uuid.UUID      `json:"organization_id" db:"organization_id"`
	KeyName            string         `json:"key_name" db:"key_name"`
	KeyPrefix          string         `json:"key_prefix" db:"key_prefix"`
	KeyHash            string         `json:"key_hash" db:"key_hash"`
	Scopes             pq.StringArray `json:"scopes" db:"scopes"`
	AllowedIPs         pq.StringArray `json:"allowed_ips" db:"allowed_ips"`
	IsActive           bool           `json:"is_active" db:"is_active"`
	LastUsedAt         *time.Time     `json:"last_used_at" db:"last_used_at"`
	UsageCount         int            `json:"usage_count" db:"usage_count"`
	RateLimitPerMinute int            `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
	RateLimitPerHour   int            `json:"rate_limit_per_hour" db:"rate_limit_per_hour"`
	ExpiresAt          *time.Time     `json:"expires_at" db:"expires_at"`
	CreatedBy          *uuid.UUID     `json:"created_by" db:"created_by"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time     `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// WEBHOOKS - Outbound webhook configurations
// ============================================================================

// Webhook represents a webhook configuration
type Webhook struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	OrganizationID       uuid.UUID       `json:"organization_id" db:"organization_id"`
	WebhookName          string          `json:"webhook_name" db:"webhook_name"`
	URL                  string          `json:"url" db:"url"`
	Secret               *string         `json:"secret" db:"secret"`
	Events               pq.StringArray  `json:"events" db:"events"`
	HTTPMethod           string          `json:"http_method" db:"http_method"`
	Headers              json.RawMessage `json:"headers" db:"headers"`
	TimeoutSeconds       int             `json:"timeout_seconds" db:"timeout_seconds"`
	MaxRetries           int             `json:"max_retries" db:"max_retries"`
	RetryBackoffSeconds  int             `json:"retry_backoff_seconds" db:"retry_backoff_seconds"`
	IsActive             bool            `json:"is_active" db:"is_active"`
	IsVerified           bool            `json:"is_verified" db:"is_verified"`
	TotalDeliveries      int             `json:"total_deliveries" db:"total_deliveries"`
	SuccessfulDeliveries int             `json:"successful_deliveries" db:"successful_deliveries"`
	FailedDeliveries     int             `json:"failed_deliveries" db:"failed_deliveries"`
	LastDeliveryAt       *time.Time      `json:"last_delivery_at" db:"last_delivery_at"`
	LastSuccessAt        *time.Time      `json:"last_success_at" db:"last_success_at"`
	LastFailureAt        *time.Time      `json:"last_failure_at" db:"last_failure_at"`
	CreatedBy            *uuid.UUID      `json:"created_by" db:"created_by"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time      `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// WEBHOOK DELIVERIES - Webhook delivery log
// ============================================================================

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID                 uuid.UUID        `json:"id" db:"id"`
	OrganizationID     uuid.UUID        `json:"organization_id" db:"organization_id"`
	WebhookID          uuid.UUID        `json:"webhook_id" db:"webhook_id"`
	EventType          string           `json:"event_type" db:"event_type"`
	EventID            uuid.UUID        `json:"event_id" db:"event_id"`
	Status             string           `json:"status" db:"status"`
	RequestURL         string           `json:"request_url" db:"request_url"`
	RequestMethod      string           `json:"request_method" db:"request_method"`
	RequestHeaders     *json.RawMessage `json:"request_headers" db:"request_headers"`
	RequestBody        *json.RawMessage `json:"request_body" db:"request_body"`
	ResponseStatusCode *int             `json:"response_status_code" db:"response_status_code"`
	ResponseHeaders    *json.RawMessage `json:"response_headers" db:"response_headers"`
	ResponseBody       *string          `json:"response_body" db:"response_body"`
	AttemptNumber      int              `json:"attempt_number" db:"attempt_number"`
	DurationMs         *int             `json:"duration_ms" db:"duration_ms"`
	NextRetryAt        *time.Time       `json:"next_retry_at" db:"next_retry_at"`
	ErrorMessage       *string          `json:"error_message" db:"error_message"`
	CreatedAt          time.Time        `json:"created_at" db:"created_at"`
	DeliveredAt        *time.Time       `json:"delivered_at" db:"delivered_at"`
}

// ============================================================================
// NOTIFICATIONS - User notifications
// ============================================================================

// Notification represents a user notification
type Notification struct {
	ID                uuid.UUID      `json:"id" db:"id"`
	OrganizationID    uuid.UUID      `json:"organization_id" db:"organization_id"`
	UserID            uuid.UUID      `json:"user_id" db:"user_id"`
	NotificationType  string         `json:"notification_type" db:"notification_type"`
	Category          string         `json:"category" db:"category"`
	Title             string         `json:"title" db:"title"`
	Message           string         `json:"message" db:"message"`
	ActionURL         *string        `json:"action_url" db:"action_url"`
	ActionLabel       *string        `json:"action_label" db:"action_label"`
	Channels          pq.StringArray `json:"channels" db:"channels"`
	IsRead            bool           `json:"is_read" db:"is_read"`
	ReadAt            *time.Time     `json:"read_at" db:"read_at"`
	RelatedEntityType *string        `json:"related_entity_type" db:"related_entity_type"`
	RelatedEntityID   *uuid.UUID     `json:"related_entity_id" db:"related_entity_id"`
	Priority          string         `json:"priority" db:"priority"`
	ExpiresAt         *time.Time     `json:"expires_at" db:"expires_at"`
	CreatedAt         time.Time      `json:"created_at" db:"created_at"`
}

// ============================================================================
// NOTIFICATION PREFERENCES - User notification settings
// ============================================================================

// NotificationPreference represents user notification preferences
type NotificationPreference struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Category       string    `json:"category" db:"category"`
	InAppEnabled   bool      `json:"in_app_enabled" db:"in_app_enabled"`
	EmailEnabled   bool      `json:"email_enabled" db:"email_enabled"`
	SMSEnabled     bool      `json:"sms_enabled" db:"sms_enabled"`
	PushEnabled    bool      `json:"push_enabled" db:"push_enabled"`
	Frequency      string    `json:"frequency" db:"frequency"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// USER SESSIONS - Session management
// ============================================================================

// UserSession represents a user session
type UserSession struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	OrganizationID *uuid.UUID `json:"organization_id" db:"organization_id"`
	SessionToken   string     `json:"session_token" db:"session_token"`
	RefreshToken   *string    `json:"refresh_token" db:"refresh_token"`
	UserAgent      *string    `json:"user_agent" db:"user_agent"`
	IPAddress      *string    `json:"ip_address" db:"ip_address"`
	DeviceType     *string    `json:"device_type" db:"device_type"`
	DeviceName     *string    `json:"device_name" db:"device_name"`
	Browser        *string    `json:"browser" db:"browser"`
	OS             *string    `json:"os" db:"os"`
	CountryCode    *string    `json:"country_code" db:"country_code"`
	City           *string    `json:"city" db:"city"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	LastActivityAt time.Time  `json:"last_activity_at" db:"last_activity_at"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	RevokedAt      *time.Time `json:"revoked_at" db:"revoked_at"`
}

// ============================================================================
// REPOSITORY INTERFACES - Core infrastructure data access
// ============================================================================

// BackgroundJobRepository defines background job data access
type BackgroundJobRepository interface {
	List(ctx context.Context, orgID uuid.UUID, status *string, limit int, offset int) ([]BackgroundJob, error)
	Create(ctx context.Context, job *BackgroundJob) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*BackgroundJob, error)
	Update(ctx context.Context, job *BackgroundJob) error
	UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error
}

// APIKeyRepository defines API key data access
type APIKeyRepository interface {
	List(ctx context.Context, orgID uuid.UUID, isActive *bool, limit int, offset int) ([]APIKey, error)
	Create(ctx context.Context, key *APIKey) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*APIKey, error)
	GetByKeyHash(ctx context.Context, keyHash string) (*APIKey, error)
	Update(ctx context.Context, key *APIKey) error
	UpdateUsage(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

// WebhookRepository defines webhook data access
type WebhookRepository interface {
	List(ctx context.Context, orgID uuid.UUID, isActive *bool, limit int, offset int) ([]Webhook, error)
	Create(ctx context.Context, webhook *Webhook) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Webhook, error)
	Update(ctx context.Context, webhook *Webhook) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateStats(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error
}

// WebhookDeliveryRepository defines webhook delivery data access
type WebhookDeliveryRepository interface {
	List(ctx context.Context, orgID uuid.UUID, webhookID *uuid.UUID, status *string, limit int, offset int) ([]WebhookDelivery, error)
	Create(ctx context.Context, delivery *WebhookDelivery) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*WebhookDelivery, error)
	Update(ctx context.Context, delivery *WebhookDelivery) error
	GetPendingDeliveries(ctx context.Context, limit int) ([]WebhookDelivery, error)
}

// NotificationRepository defines notification data access
type NotificationRepository interface {
	List(ctx context.Context, userID uuid.UUID, isRead *bool, limit int, offset int) ([]Notification, error)
	Create(ctx context.Context, notification *Notification) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Notification, error)
	Update(ctx context.Context, notification *Notification) error
	MarkAsRead(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
}

// NotificationPreferenceRepository defines notification preference data access
type NotificationPreferenceRepository interface {
	Get(ctx context.Context, userID uuid.UUID, category string) (*NotificationPreference, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]NotificationPreference, error)
	Upsert(ctx context.Context, pref *NotificationPreference) error
	Delete(ctx context.Context, userID uuid.UUID, category string) error
}

// UserSessionRepository defines user session data access
type UserSessionRepository interface {
	List(ctx context.Context, userID uuid.UUID, isActive *bool, limit int, offset int) ([]UserSession, error)
	Create(ctx context.Context, session *UserSession) error
	Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*UserSession, error)
	GetByToken(ctx context.Context, token string) (*UserSession, error)
	Update(ctx context.Context, session *UserSession) error
	Revoke(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	RevokeAll(ctx context.Context, userID uuid.UUID) error
	CleanupExpired(ctx context.Context) error
}

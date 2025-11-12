package monitoring

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// POSErrorLog represents an error log entry from POS operations
type POSErrorLog struct {
	ID              uuid.UUID   `json:"id"`
	OrganizationID  *uuid.UUID  `json:"organization_id"`
	ErrorLevel      string      `json:"error_level"` // debug, info, warning, error, critical
	ErrorCode       string      `json:"error_code"`
	ErrorMessage    string      `json:"error_message"`
	DeviceID        *uuid.UUID  `json:"device_id"`
	UserID          *uuid.UUID  `json:"user_id"`
	POSSessionID    *uuid.UUID  `json:"pos_session_id"`
	SaleID          *uuid.UUID  `json:"sale_id"`
	StackTrace      string      `json:"stack_trace"`
	RequestData     interface{} `json:"request_data"`
	ErrorData       interface{} `json:"error_data"`
	IsResolved      bool        `json:"is_resolved"`
	ResolvedBy      *uuid.UUID  `json:"resolved_by"`
	ResolvedAt      *time.Time  `json:"resolved_at"`
	ResolutionNotes string      `json:"resolution_notes"`
	OccurredAt      time.Time   `json:"occurred_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

// SystemHealth represents a system health check record
type SystemHealth struct {
	ID                uuid.UUID   `json:"id"`
	OrganizationID    *uuid.UUID  `json:"organization_id"` // NULL = global check
	CheckType         string      `json:"check_type"`
	CheckName         string      `json:"check_name"`
	Status            string      `json:"status"` // healthy, degraded, unhealthy, unknown
	LastCheckAt       *time.Time  `json:"last_check_at"`
	LastSuccessAt     *time.Time  `json:"last_success_at"`
	LastFailureAt     *time.Time  `json:"last_failure_at"`
	MetricValue       *float64    `json:"metric_value"`
	MetricUnit        string      `json:"metric_unit"`
	ThresholdWarning  *float64    `json:"threshold_warning"`
	ThresholdCritical *float64    `json:"threshold_critical"`
	Details           interface{} `json:"details"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

// POSErrorLogFilters for list operations
type POSErrorLogFilters struct {
	ErrorLevel     *string
	DeviceID       *uuid.UUID
	UserID         *uuid.UUID
	POSSessionID   *uuid.UUID
	IsResolved     *bool
	StartDate      *time.Time
	EndDate        *time.Time
	Page           int
	PageSize       int
}

// SystemHealthFilters for list operations
type SystemHealthFilters struct {
	CheckType *string
	Status    *string
	Page      int
	PageSize  int
}

// POSErrorLogRepository defines data access interface
type POSErrorLogRepository interface {
	List(ctx context.Context, orgID *uuid.UUID, filters POSErrorLogFilters) ([]POSErrorLog, error)
	Count(ctx context.Context, orgID *uuid.UUID, filters POSErrorLogFilters) (int64, error)
	Create(ctx context.Context, errorLog *POSErrorLog) error
	Get(ctx context.Context, id uuid.UUID) (*POSErrorLog, error)
	Update(ctx context.Context, errorLog *POSErrorLog) error
	Resolve(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID, notes string) error
}

// SystemHealthRepository defines data access interface
type SystemHealthRepository interface {
	List(ctx context.Context, orgID *uuid.UUID, filters SystemHealthFilters) ([]SystemHealth, error)
	Count(ctx context.Context, orgID *uuid.UUID, filters SystemHealthFilters) (int64, error)
	Upsert(ctx context.Context, health *SystemHealth) error
	Get(ctx context.Context, orgID *uuid.UUID, checkType, checkName string) (*SystemHealth, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, metricValue *float64) error
	GetUnhealthyChecks(ctx context.Context, orgID *uuid.UUID) ([]SystemHealth, error)
}

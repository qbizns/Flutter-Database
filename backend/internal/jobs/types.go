package jobs

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

// Job type constants
const (
	TypeEmailNotification = "email:notification"
	TypeReportGeneration  = "report:generation"
	TypeDataExport        = "data:export"
	TypeDatabaseBackup    = "database:backup"
	TypeInvoiceGeneration = "invoice:generation"
	TypePostingBatch      = "posting:batch"
	TypeAuditLogCleanup   = "audit:cleanup"
)

// EmailNotificationPayload represents email job payload
type EmailNotificationPayload struct {
	To      []string          `json:"to"`
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	Template string           `json:"template,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ReportGenerationPayload represents report generation job payload
type ReportGenerationPayload struct {
	OrganizationID string    `json:"organization_id"`
	ReportType     string    `json:"report_type"` // trial_balance, profit_loss, balance_sheet
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Format         string    `json:"format"` // pdf, excel, csv
	UserID         string    `json:"user_id"`
	CallbackURL    string    `json:"callback_url,omitempty"`
}

// DataExportPayload represents data export job payload
type DataExportPayload struct {
	OrganizationID string   `json:"organization_id"`
	EntityType     string   `json:"entity_type"` // customers, products, sales, etc
	Format         string   `json:"format"`      // csv, json, excel
	Filters        map[string]interface{} `json:"filters,omitempty"`
	UserID         string   `json:"user_id"`
	CallbackURL    string   `json:"callback_url,omitempty"`
}

// DatabaseBackupPayload represents database backup job payload
type DatabaseBackupPayload struct {
	BackupType      string   `json:"backup_type"` // full, incremental
	OrganizationIDs []string `json:"organization_ids,omitempty"`
	S3Bucket        string   `json:"s3_bucket"`
	RetentionDays   int      `json:"retention_days"`
}

// InvoiceGenerationPayload represents invoice generation job payload
type InvoiceGenerationPayload struct {
	OrganizationID string    `json:"organization_id"`
	InvoiceID      string    `json:"invoice_id"`
	UserID         string    `json:"user_id"`
	SendEmail      bool      `json:"send_email"`
	EmailTo        []string  `json:"email_to,omitempty"`
}

// PostingBatchPayload represents batch posting job payload
type PostingBatchPayload struct {
	OrganizationID string                   `json:"organization_id"`
	Documents      []PostingBatchDocument   `json:"documents"`
	UserID         string                   `json:"user_id"`
}

type PostingBatchDocument struct {
	DocumentType string `json:"document_type"`
	DocumentID   string `json:"document_id"`
	Event        string `json:"event"`
}

// AuditLogCleanupPayload represents audit log cleanup job payload
type AuditLogCleanupPayload struct {
	RetentionDays int       `json:"retention_days"`
	BeforeDate    time.Time `json:"before_date"`
}

// JobResult represents the result of a job execution
type JobResult struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Duration  time.Duration          `json:"duration"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
}

// NewTask creates a new asynq task from job type and payload
func NewTask(taskType string, payload interface{}) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(taskType, data), nil
}

// ParsePayload parses task payload into the specified type
func ParsePayload[T any](task *asynq.Task) (*T, error) {
	var payload T
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

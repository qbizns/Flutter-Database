package infrastructure

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// Extended Infrastructure Package Exports and Utilities
// ============================================================================

// This file re-exports key infrastructure types and provides helper functions
// for the extended infrastructure layer (file attachments, queues, settings, etc.)

// ============================================================================
// Context Helpers
// ============================================================================

// ContextWithInfrastructure adds infrastructure-related values to context
func ContextWithInfrastructure(ctx context.Context, orgID uuid.UUID) context.Context {
	return context.WithValue(ctx, "organization_id", orgID)
}

// GetOrganizationIDFromContext extracts organization ID from context
func GetOrganizationIDFromContext(ctx context.Context) *uuid.UUID {
	if val := ctx.Value("organization_id"); val != nil {
		if orgID, ok := val.(uuid.UUID); ok {
			return &orgID
		}
	}
	return nil
}

// ============================================================================
// Validation Helpers
// ============================================================================

// IsValidStorageProvider validates storage provider enum
func IsValidStorageProvider(provider string) bool {
	validProviders := map[string]bool{
		"local": true,
		"s3":    true,
		"gcs":   true,
		"azure": true,
	}
	return validProviders[provider]
}

// IsValidEmailStatus validates email queue status
func IsValidEmailStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"sending":   true,
		"sent":      true,
		"failed":    true,
		"cancelled": true,
	}
	return validStatuses[status]
}

// IsValidSMSStatus validates SMS queue status
func IsValidSMSStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"sending":   true,
		"sent":      true,
		"failed":    true,
		"cancelled": true,
	}
	return validStatuses[status]
}

// IsValidExportFormat validates export format
func IsValidExportFormat(format string) bool {
	validFormats := map[string]bool{
		"csv":  true,
		"xlsx": true,
		"json": true,
		"pdf":  true,
	}
	return validFormats[format]
}

// IsValidReportFrequency validates report schedule frequency
func IsValidReportFrequency(frequency string) bool {
	validFrequencies := map[string]bool{
		"daily":     true,
		"weekly":    true,
		"monthly":   true,
		"quarterly": true,
	}
	return validFrequencies[frequency]
}

// IsValidIntegrationType validates integration type
func IsValidIntegrationType(intType string) bool {
	validTypes := map[string]bool{
		"payment_gateway": true,
		"shipping":        true,
		"accounting":      true,
		"crm":             true,
	}
	return validTypes[intType]
}

// ============================================================================
// Constants
// ============================================================================

const (
	// Storage providers
	StorageProviderLocal = "local"
	StorageProviderS3    = "s3"
	StorageProviderGCS   = "gcs"
	StorageProviderAzure = "azure"

	// Email/SMS statuses
	StatusPending   = "pending"
	StatusSending   = "sending"
	StatusSent      = "sent"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"

	// Virus scan statuses
	VirusScanPending  = "pending"
	VirusScanClean    = "clean"
	VirusScanInfected = "infected"
	VirusScanError    = "error"

	// Report frequencies
	FrequencyDaily     = "daily"
	FrequencyWeekly    = "weekly"
	FrequencyMonthly   = "monthly"
	FrequencyQuarterly = "quarterly"

	// Export formats
	ExportFormatCSV  = "csv"
	ExportFormatXLSX = "xlsx"
	ExportFormatJSON = "json"
	ExportFormatPDF  = "pdf"

	// Rate limit identifier types
	IdentifierTypeUser         = "user"
	IdentifierTypeOrganization = "organization"
	IdentifierTypeAPIKey       = "api_key"
	IdentifierTypeIP           = "ip"
)

// ============================================================================
// Query Helpers
// ============================================================================

// BuildPaginationQuery adds LIMIT and OFFSET to a query
func BuildPaginationQuery(page, pageSize int) (string, []interface{}) {
	var query string
	var args []interface{}

	if pageSize > 0 {
		query = " LIMIT $1"
		args = append(args, pageSize)

		if page > 1 {
			offset := (page - 1) * pageSize
			query += " OFFSET $2"
			args = append(args, offset)
		}
	}

	return query, args
}

// IsExpired checks if a time is in the past
func IsExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	return expiresAt.Before(time.Now())
}

// GetEmailQueuePriority returns query priority order
func GetEmailQueuePriority() string {
	return "ORDER BY priority DESC, scheduled_at ASC, created_at ASC"
}

// GetSMSQueuePriority returns query priority order
func GetSMSQueuePriority() string {
	return "ORDER BY scheduled_at ASC, created_at ASC"
}

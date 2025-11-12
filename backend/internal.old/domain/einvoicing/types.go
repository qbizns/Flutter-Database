package einvoicing

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EInvoicingDocument represents an e-invoicing document
type EInvoicingDocument struct {
	ID                           uuid.UUID                  `json:"id"`
	OrganizationID               uuid.UUID                  `json:"organization_id"`
	SourceTable                  string                     `json:"source_table"`
	SourceID                     uuid.UUID                  `json:"source_id"`
	Authority                    string                     `json:"authority"`
	CountryCode                  string                     `json:"country_code"`
	DocumentUUID                 uuid.UUID                  `json:"document_uuid"`
	DocumentType                 string                     `json:"document_type"`
	DocumentNumber               string                     `json:"document_number"`
	InternalReference            *string                    `json:"internal_reference"`
	Status                       string                     `json:"status"`
	CreatedAt                    time.Time                  `json:"created_at"`
	SubmittedAt                  *time.Time                 `json:"submitted_at"`
	ResponseAt                   *time.Time                 `json:"response_at"`
	RequestPayload               JSONB                      `json:"request_payload"`
	ResponsePayload              JSONB                      `json:"response_payload"`
	ErrorCode                    *string                    `json:"error_code"`
	ErrorMessage                 *string                    `json:"error_message"`
	RetryCount                   int                        `json:"retry_count"`
	LastRetryAt                  *time.Time                 `json:"last_retry_at"`
	ZATCAHashValue               *string                    `json:"zatca_hash_value"`
	ZATCAPreviousHashValue       *string                    `json:"zatca_previous_hash_value"`
	ZATCAInvoiceCounterValue     *int                       `json:"zatca_invoice_counter_value"`
	ZATCACryptographicStamp      *string                    `json:"zatca_cryptographic_stamp"`
	ZATCAQRCodePayload           *string                    `json:"zatca_qr_code_payload"`
	ZATCAComplianceInvoiceNumber *string                    `json:"zatca_compliance_invoice_number"`
	ETADocumentTypeVersion       *string                    `json:"eta_document_type_version"`
	ETASubmissionUUID            *uuid.UUID                 `json:"eta_submission_uuid"`
	ETALongID                    *string                    `json:"eta_long_id"`
	ETAInternalID                *string                    `json:"eta_internal_id"`
	ETADigitalSignature          *string                    `json:"eta_digital_signature"`
	ETASignatureAlgorithm        *string                    `json:"eta_signature_algorithm"`
	SubmissionFormat             *string                    `json:"submission_format"`
	Metadata                     JSONB                      `json:"metadata"`
	UpdatedAt                    time.Time                  `json:"updated_at"`
	DeletedAt                    *time.Time                 `json:"deleted_at"`
	CreatedBy                    *uuid.UUID                 `json:"created_by"`
	UpdatedBy                    *uuid.UUID                 `json:"updated_by"`
}

// EInvoicingDocumentEvent represents an event in the e-invoicing workflow
type EInvoicingDocumentEvent struct {
	ID                      uuid.UUID      `json:"id"`
	OrganizationID          uuid.UUID      `json:"organization_id"`
	EInvoicingDocumentID    uuid.UUID      `json:"e_invoicing_document_id"`
	EventType               string         `json:"event_type"`
	EventTimestamp          time.Time      `json:"event_timestamp"`
	PreviousStatus          *string        `json:"previous_status"`
	NewStatus               *string        `json:"new_status"`
	EventDescription        *string        `json:"event_description"`
	EventData               JSONB          `json:"event_data"`
	HTTPStatusCode          *int           `json:"http_status_code"`
	HTTPMethod              *string        `json:"http_method"`
	APIEndpoint             *string        `json:"api_endpoint"`
	RequestHeaders          JSONB          `json:"request_headers"`
	ResponseHeaders         JSONB          `json:"response_headers"`
	ErrorCode               *string        `json:"error_code"`
	ErrorMessage            *string        `json:"error_message"`
	ErrorDetails            JSONB          `json:"error_details"`
	TriggeredBy             string         `json:"triggered_by"`
	UserID                  *uuid.UUID     `json:"user_id"`
	CreatedAt               time.Time      `json:"created_at"`
}

// DocumentSequence represents a document numbering sequence
type DocumentSequence struct {
	ID                  uuid.UUID  `json:"id"`
	OrganizationID      uuid.UUID  `json:"organization_id"`
	DocumentType        string     `json:"document_type"`
	Prefix              string     `json:"prefix"`
	Suffix              string     `json:"suffix"`
	NextNumber          int64      `json:"next_number"`
	Padding             int        `json:"padding"`
	IncrementBy         int        `json:"increment_by"`
	ResetFrequency      string     `json:"reset_frequency"`
	LastResetAt         *time.Time `json:"last_reset_at"`
	LastResetValue      int64      `json:"last_reset_value"`
	IncludeDate         bool       `json:"include_date"`
	DateFormat          *string    `json:"date_format"`
	LocationID          *uuid.UUID `json:"location_id"`
	IsActive            bool       `json:"is_active"`
	AllowManualOverride bool       `json:"allow_manual_override"`
	ExampleNumber       *string    `json:"example_number"`
	Description         *string    `json:"description"`
	Notes               *string    `json:"notes"`
	Metadata            JSONB      `json:"metadata"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
	CreatedBy           *uuid.UUID `json:"created_by"`
	UpdatedBy           *uuid.UUID `json:"updated_by"`
}

// Filters
type EInvoicingDocumentFilters struct {
	Search          string
	Authority       *string
	Status          *string
	DocumentType    *string
	StartDate       *time.Time
	EndDate         *time.Time
	SourceTable     *string
	Page            int
	PageSize        int
}

type EInvoicingDocumentEventFilters struct {
	EventType          *string
	EInvoicingDocID    *uuid.UUID
	StartDate          *time.Time
	EndDate            *time.Time
	Page               int
	PageSize           int
}

type DocumentSequenceFilters struct {
	DocumentType *string
	LocationID   *uuid.UUID
	IsActive     *bool
	Page         int
	PageSize     int
}

// JSONB is a custom type for JSON data
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &j)
}

// Repository interfaces
type EInvoicingDocumentRepository interface {
	// CRUD operations
	List(ctx context.Context, orgID uuid.UUID, filters EInvoicingDocumentFilters) ([]EInvoicingDocument, error)
	Count(ctx context.Context, orgID uuid.UUID, filters EInvoicingDocumentFilters) (int64, error)
	Create(ctx context.Context, doc *EInvoicingDocument) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*EInvoicingDocument, error)
	GetByDocumentNumber(ctx context.Context, orgID uuid.UUID, authority string, documentNumber string) (*EInvoicingDocument, error)
	GetByDocumentUUID(ctx context.Context, orgID uuid.UUID, documentUUID uuid.UUID) (*EInvoicingDocument, error)
	Update(ctx context.Context, doc *EInvoicingDocument) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Events
	CreateEvent(ctx context.Context, event *EInvoicingDocumentEvent) error
	ListEvents(ctx context.Context, orgID uuid.UUID, filters EInvoicingDocumentEventFilters) ([]EInvoicingDocumentEvent, error)
	CountEvents(ctx context.Context, orgID uuid.UUID, filters EInvoicingDocumentEventFilters) (int64, error)
	GetEvent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*EInvoicingDocumentEvent, error)

	// Bulk operations
	GetBySourceReference(ctx context.Context, orgID uuid.UUID, sourceTable string, sourceID uuid.UUID) ([]EInvoicingDocument, error)
	UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error
}

type DocumentSequenceRepository interface {
	Create(ctx context.Context, seq *DocumentSequence) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DocumentSequence, error)
	GetByDocumentType(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (*DocumentSequence, error)
	List(ctx context.Context, orgID uuid.UUID, filters DocumentSequenceFilters) ([]DocumentSequence, error)
	Count(ctx context.Context, orgID uuid.UUID, filters DocumentSequenceFilters) (int64, error)
	Update(ctx context.Context, seq *DocumentSequence) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetNextNumber(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (string, error)
	PreviewNumber(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (string, error)
	ResetSequence(ctx context.Context, orgID uuid.UUID, documentType string, resetTo int64, locationID *uuid.UUID) error
}

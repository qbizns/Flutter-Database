package posting

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PostingRule represents a posting rule configuration
type PostingRule struct {
	ID                  uuid.UUID
	RuleCode            string
	RuleName            string
	Event               string
	Level               string
	Priority            int
	ConditionExpression string
	IsActive            bool
}

// PostingRuleLine represents a journal entry line template
type PostingRuleLine struct {
	ID                  uuid.UUID
	PostingRuleID       uuid.UUID
	LineNo              int
	Side                string // 'debit' or 'credit'
	ConceptKey          string
	AccountSource       string
	AmountSource        string
	AmountFieldPath     string
	AmountExpression    string
	DescriptionTemplate string
}

// JournalEntry represents a journal entry
type JournalEntry struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	EntryNumber    string
	EntryDate      time.Time
	PostingDate    time.Time
	Description    string
	Status         string
	TotalDebit     float64
	TotalCredit    float64
	Lines          []JournalEntryLine
	CreatedAt      time.Time
	CreatedBy      uuid.UUID
}

// JournalEntryLine represents a journal entry line
type JournalEntryLine struct {
	ID           uuid.UUID
	AccountID    uuid.UUID
	DebitAmount  float64
	CreditAmount float64
	Description  string
}

// PostingAudit represents posting audit log
type PostingAudit struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	SourceTable    string
	SourceID       uuid.UUID
	JournalEntryID uuid.UUID
	PostingStatus  string
	Event          string
	ErrorMessage   string
	PostedAt       time.Time
}

// ValidationRule represents a posting validation rule
type ValidationRule struct {
	ID              uuid.UUID
	Code            string
	Expression      string
	Severity        string
	IsBlocking      bool
	MessageTemplate string
}

// ValidationResult represents a validation result
type ValidationResult struct {
	OrganizationID   uuid.UUID
	DocumentType     string
	DocumentID       uuid.UUID
	Event            string
	ValidationRuleID uuid.UUID
	Severity         string
	Message          string
	IsBlocking       bool
}

// Repository defines the posting engine data access interface
type Repository interface {
	// Document operations
	LoadDocument(ctx context.Context, documentType string, documentID uuid.UUID) (map[string]interface{}, error)
	UpdateDocumentPostingStatus(ctx context.Context, documentType string, documentID uuid.UUID, journalEntryID uuid.UUID, status string) error

	// Posting rules
	GetPostingRules(ctx context.Context, orgID uuid.UUID, documentType, event string) ([]PostingRule, error)
	GetPostingRuleLines(ctx context.Context, ruleID uuid.UUID) ([]PostingRuleLine, error)

	// Account resolution
	ResolveAccountFromConcept(ctx context.Context, orgID uuid.UUID, conceptKey string) (uuid.UUID, error)

	// Journal entry operations
	CreateJournalEntry(ctx context.Context, je *JournalEntry) error

	// Validation
	GetValidationRules(ctx context.Context, orgID uuid.UUID, documentType, event string) ([]ValidationRule, error)
	LogValidationResult(ctx context.Context, result ValidationResult) error

	// Audit
	LogPostingAudit(ctx context.Context, audit PostingAudit) error
	GetPostingAuditLogs(ctx context.Context, orgID uuid.UUID, filter map[string]interface{}) ([]PostingAudit, error)
}

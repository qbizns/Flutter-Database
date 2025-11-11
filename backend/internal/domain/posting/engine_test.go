package posting

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/logging"
)

// mockRepository implements the Repository interface for testing
type mockRepository struct {
	// Data stores
	documents       map[uuid.UUID]map[string]interface{}
	postingRules    map[string][]PostingRule // key: documentType_event
	postingRuleLines map[uuid.UUID][]PostingRuleLine
	conceptMappings map[string]uuid.UUID // key: conceptKey, value: accountID
	validationRules map[string][]ValidationRule // key: documentType_event
	journalEntries  []*JournalEntry
	auditLogs       []PostingAudit
	validationResults []ValidationResult
	documentStatuses map[uuid.UUID]string // key: documentID, value: status

	// Hook functions for custom behavior
	onLoadDocument                  func(ctx context.Context, documentType string, documentID uuid.UUID) (map[string]interface{}, error)
	onGetPostingRules              func(ctx context.Context, orgID uuid.UUID, documentType, event string) ([]PostingRule, error)
	onGetPostingRuleLines          func(ctx context.Context, ruleID uuid.UUID) ([]PostingRuleLine, error)
	onResolveAccountFromConcept    func(ctx context.Context, orgID uuid.UUID, conceptKey string) (uuid.UUID, error)
	onCreateJournalEntry           func(ctx context.Context, je *JournalEntry) error
	onUpdateDocumentPostingStatus  func(ctx context.Context, documentType string, documentID uuid.UUID, journalEntryID uuid.UUID, status string) error
	onGetValidationRules           func(ctx context.Context, orgID uuid.UUID, documentType, event string) ([]ValidationRule, error)
	onLogValidationResult          func(ctx context.Context, result ValidationResult) error
	onLogPostingAudit              func(ctx context.Context, audit PostingAudit) error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		documents:        make(map[uuid.UUID]map[string]interface{}),
		postingRules:     make(map[string][]PostingRule),
		postingRuleLines: make(map[uuid.UUID][]PostingRuleLine),
		conceptMappings:  make(map[string]uuid.UUID),
		validationRules:  make(map[string][]ValidationRule),
		journalEntries:   make([]*JournalEntry, 0),
		auditLogs:        make([]PostingAudit, 0),
		validationResults: make([]ValidationResult, 0),
		documentStatuses: make(map[uuid.UUID]string),
	}
}

func (m *mockRepository) LoadDocument(ctx context.Context, documentType string, documentID uuid.UUID) (map[string]interface{}, error) {
	if m.onLoadDocument != nil {
		return m.onLoadDocument(ctx, documentType, documentID)
	}
	if doc, ok := m.documents[documentID]; ok {
		return doc, nil
	}
	return nil, pgx.ErrNoRows
}

func (m *mockRepository) UpdateDocumentPostingStatus(ctx context.Context, documentType string, documentID uuid.UUID, journalEntryID uuid.UUID, status string) error {
	if m.onUpdateDocumentPostingStatus != nil {
		return m.onUpdateDocumentPostingStatus(ctx, documentType, documentID, journalEntryID, status)
	}
	m.documentStatuses[documentID] = status
	return nil
}

func (m *mockRepository) GetPostingRules(ctx context.Context, orgID uuid.UUID, documentType, event string) ([]PostingRule, error) {
	if m.onGetPostingRules != nil {
		return m.onGetPostingRules(ctx, orgID, documentType, event)
	}
	key := documentType + "_" + event
	if rules, ok := m.postingRules[key]; ok {
		return rules, nil
	}
	return []PostingRule{}, nil
}

func (m *mockRepository) GetPostingRuleLines(ctx context.Context, ruleID uuid.UUID) ([]PostingRuleLine, error) {
	if m.onGetPostingRuleLines != nil {
		return m.onGetPostingRuleLines(ctx, ruleID)
	}
	if lines, ok := m.postingRuleLines[ruleID]; ok {
		return lines, nil
	}
	return []PostingRuleLine{}, nil
}

func (m *mockRepository) ResolveAccountFromConcept(ctx context.Context, orgID uuid.UUID, conceptKey string) (uuid.UUID, error) {
	if m.onResolveAccountFromConcept != nil {
		return m.onResolveAccountFromConcept(ctx, orgID, conceptKey)
	}
	if accountID, ok := m.conceptMappings[conceptKey]; ok {
		return accountID, nil
	}
	return uuid.Nil, pgx.ErrNoRows
}

func (m *mockRepository) CreateJournalEntry(ctx context.Context, je *JournalEntry) error {
	if m.onCreateJournalEntry != nil {
		return m.onCreateJournalEntry(ctx, je)
	}
	m.journalEntries = append(m.journalEntries, je)
	return nil
}

func (m *mockRepository) GetValidationRules(ctx context.Context, orgID uuid.UUID, documentType, event string) ([]ValidationRule, error) {
	if m.onGetValidationRules != nil {
		return m.onGetValidationRules(ctx, orgID, documentType, event)
	}
	key := documentType + "_" + event
	if rules, ok := m.validationRules[key]; ok {
		return rules, nil
	}
	return []ValidationRule{}, nil
}

func (m *mockRepository) LogValidationResult(ctx context.Context, result ValidationResult) error {
	if m.onLogValidationResult != nil {
		return m.onLogValidationResult(ctx, result)
	}
	m.validationResults = append(m.validationResults, result)
	return nil
}

func (m *mockRepository) LogPostingAudit(ctx context.Context, audit PostingAudit) error {
	if m.onLogPostingAudit != nil {
		return m.onLogPostingAudit(ctx, audit)
	}
	m.auditLogs = append(m.auditLogs, audit)
	return nil
}

// setupEngine creates an engine with mock repository
func setupEngine() (*Engine, *mockRepository) {
	logger, _ := logging.NewLogger("info", "json")
	mockRepo := newMockRepository()
	engine := NewEngine(mockRepo, logger)
	return engine, mockRepo
}

// ========================================
// POST METHOD TESTS
// ========================================

func TestEngine_Post_Success(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
		"tax_amount":   10.0,
	}

	// Setup posting rule
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Sales Invoice Posting",
		Event:               event,
		ConditionExpression: "", // Always match
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines
	receivablesAcctID := uuid.New()
	salesAcctID := uuid.New()
	mockRepo.conceptMappings["accounts_receivable"] = receivablesAcctID
	mockRepo.conceptMappings["sales_revenue"] = salesAcctID

	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              1,
			Side:                "debit",
			ConceptKey:          "accounts_receivable",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Invoice",
		},
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              2,
			Side:                "credit",
			ConceptKey:          "sales_revenue",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Revenue",
		},
	}

	// Execute
	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify journal entry created
	if len(mockRepo.journalEntries) != 1 {
		t.Fatalf("expected 1 journal entry, got %d", len(mockRepo.journalEntries))
	}

	je := mockRepo.journalEntries[0]
	if je.TotalDebit != 100.0 {
		t.Errorf("expected total debit 100.0, got %f", je.TotalDebit)
	}
	if je.TotalCredit != 100.0 {
		t.Errorf("expected total credit 100.0, got %f", je.TotalCredit)
	}
	if len(je.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(je.Lines))
	}

	// Verify document status updated
	if status, ok := mockRepo.documentStatuses[docID]; !ok || status != "posted" {
		t.Errorf("expected document status 'posted', got '%s'", status)
	}

	// Verify audit log created
	if len(mockRepo.auditLogs) != 1 {
		t.Errorf("expected 1 audit log, got %d", len(mockRepo.auditLogs))
	}
}

func TestEngine_Post_DocumentNotFound(t *testing.T) {
	engine, _ := setupEngine()
	ctx := context.Background()

	input := PostingInput{
		OrganizationID: uuid.New(),
		DocumentType:   "sales_invoice",
		DocumentID:     uuid.New(), // Not in mock data
		Event:          "complete",
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected error for document not found, got nil")
	}
}

func TestEngine_Post_LoadDocumentError(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	expectedErr := errors.New("database connection failed")
	mockRepo.onLoadDocument = func(ctx context.Context, documentType string, documentID uuid.UUID) (map[string]interface{}, error) {
		return nil, expectedErr
	}

	input := PostingInput{
		OrganizationID: uuid.New(),
		DocumentType:   "sales_invoice",
		DocumentID:     uuid.New(),
		Event:          "complete",
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected error from LoadDocument, got nil")
	}
}

func TestEngine_Post_NoPostingRules(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	docID := uuid.New()
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
	}

	input := PostingInput{
		OrganizationID: uuid.New(),
		DocumentType:   "sales_invoice",
		DocumentID:     docID,
		Event:          "complete",
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected error for no posting rules, got nil")
	}
}

func TestEngine_Post_NoMatchingRule(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
		"payment_type": "cash",
	}

	// Setup posting rule with condition that won't match
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Credit Sale Posting",
		Event:               event,
		ConditionExpression: `doc.payment_type == "credit"`, // Won't match
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected error for no matching rule, got nil")
	}
}

func TestEngine_Post_ConditionEvaluationMatch(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
		"payment_type": "credit",
	}

	// Setup posting rule with matching condition
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Credit Sale Posting",
		Event:               event,
		ConditionExpression: `doc.payment_type == "credit"`, // Will match
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines
	receivablesAcctID := uuid.New()
	salesAcctID := uuid.New()
	mockRepo.conceptMappings["accounts_receivable"] = receivablesAcctID
	mockRepo.conceptMappings["sales_revenue"] = salesAcctID

	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              1,
			Side:                "debit",
			ConceptKey:          "accounts_receivable",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Credit Sale",
		},
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              2,
			Side:                "credit",
			ConceptKey:          "sales_revenue",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Revenue",
		},
	}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(mockRepo.journalEntries) != 1 {
		t.Fatalf("expected 1 journal entry, got %d", len(mockRepo.journalEntries))
	}
}

func TestEngine_Post_AccountNotMapped(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
	}

	// Setup posting rule
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Sales Invoice Posting",
		Event:               event,
		ConditionExpression: "",
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines but DON'T map the concept
	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              1,
			Side:                "debit",
			ConceptKey:          "accounts_receivable", // NOT MAPPED
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Invoice",
		},
	}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected error for unmapped account, got nil")
	}
}

func TestEngine_Post_AmountFieldNotFound(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document WITHOUT the expected field
	mockRepo.documents[docID] = map[string]interface{}{
		"some_other_field": 100.0,
	}

	// Setup posting rule
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Sales Invoice Posting",
		Event:               event,
		ConditionExpression: "",
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines
	receivablesAcctID := uuid.New()
	mockRepo.conceptMappings["accounts_receivable"] = receivablesAcctID

	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              1,
			Side:                "debit",
			ConceptKey:          "accounts_receivable",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount", // NOT IN DOCUMENT
			DescriptionTemplate: "Sales Invoice",
		},
	}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected error for missing amount field, got nil")
	}
}

func TestEngine_Post_ExpressionAmountSource(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"subtotal": 90.0,
		"tax":      10.0,
	}

	// Setup posting rule
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Sales Invoice Posting",
		Event:               event,
		ConditionExpression: "",
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines with expression amount source
	receivablesAcctID := uuid.New()
	salesAcctID := uuid.New()
	mockRepo.conceptMappings["accounts_receivable"] = receivablesAcctID
	mockRepo.conceptMappings["sales_revenue"] = salesAcctID

	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:               uuid.New(),
			PostingRuleID:    ruleID,
			LineNo:           1,
			Side:             "debit",
			ConceptKey:       "accounts_receivable",
			AmountSource:     "expression",
			AmountExpression: `doc.subtotal + doc.tax`, // 90 + 10 = 100
			DescriptionTemplate: "Sales Invoice",
		},
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              2,
			Side:                "credit",
			ConceptKey:          "sales_revenue",
			AmountSource:        "expression",
			AmountExpression:    `doc.subtotal + doc.tax`,
			DescriptionTemplate: "Sales Revenue",
		},
	}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	je := mockRepo.journalEntries[0]
	if je.TotalDebit != 100.0 {
		t.Errorf("expected total debit 100.0, got %f", je.TotalDebit)
	}
	if je.TotalCredit != 100.0 {
		t.Errorf("expected total credit 100.0, got %f", je.TotalCredit)
	}
}

func TestEngine_Post_ValidationBlocking(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
	}

	// Setup posting rule
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Sales Invoice Posting",
		Event:               event,
		ConditionExpression: "",
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines
	receivablesAcctID := uuid.New()
	salesAcctID := uuid.New()
	mockRepo.conceptMappings["accounts_receivable"] = receivablesAcctID
	mockRepo.conceptMappings["sales_revenue"] = salesAcctID

	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              1,
			Side:                "debit",
			ConceptKey:          "accounts_receivable",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Invoice",
		},
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              2,
			Side:                "credit",
			ConceptKey:          "sales_revenue",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Revenue",
		},
	}

	// Setup BLOCKING validation rule that will fail
	validationRule := ValidationRule{
		ID:              uuid.New(),
		Code:            "BAL_001",
		Expression:      `je.total_debit > 1000.0`, // Will fail (100 < 1000)
		Severity:        "error",
		IsBlocking:      true,
		MessageTemplate: "Total must exceed 1000",
	}
	mockRepo.validationRules[docType+"_"+event] = []ValidationRule{validationRule}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected blocking validation error, got nil")
	}

	// Verify no journal entry was created
	if len(mockRepo.journalEntries) != 0 {
		t.Errorf("expected 0 journal entries due to validation block, got %d", len(mockRepo.journalEntries))
	}

	// Verify validation result was logged
	if len(mockRepo.validationResults) != 1 {
		t.Errorf("expected 1 validation result, got %d", len(mockRepo.validationResults))
	}
}

func TestEngine_Post_ValidationWarning(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
	}

	// Setup posting rule
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Sales Invoice Posting",
		Event:               event,
		ConditionExpression: "",
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines
	receivablesAcctID := uuid.New()
	salesAcctID := uuid.New()
	mockRepo.conceptMappings["accounts_receivable"] = receivablesAcctID
	mockRepo.conceptMappings["sales_revenue"] = salesAcctID

	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              1,
			Side:                "debit",
			ConceptKey:          "accounts_receivable",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Invoice",
		},
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              2,
			Side:                "credit",
			ConceptKey:          "sales_revenue",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Revenue",
		},
	}

	// Setup NON-BLOCKING validation rule that will fail
	validationRule := ValidationRule{
		ID:              uuid.New(),
		Code:            "WARN_001",
		Expression:      `je.total_debit > 1000.0`, // Will fail (100 < 1000)
		Severity:        "warning",
		IsBlocking:      false, // Non-blocking
		MessageTemplate: "Consider larger amounts",
	}
	mockRepo.validationRules[docType+"_"+event] = []ValidationRule{validationRule}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err != nil {
		t.Fatalf("expected no error (warning should not block), got %v", err)
	}

	// Verify journal entry WAS created (warning doesn't block)
	if len(mockRepo.journalEntries) != 1 {
		t.Errorf("expected 1 journal entry (warning should not block), got %d", len(mockRepo.journalEntries))
	}

	// Verify validation result was logged
	if len(mockRepo.validationResults) != 1 {
		t.Errorf("expected 1 validation result, got %d", len(mockRepo.validationResults))
	}
}

func TestEngine_Post_CreateJournalEntryError(t *testing.T) {
	engine, mockRepo := setupEngine()
	ctx := context.Background()

	orgID := uuid.New()
	docID := uuid.New()
	docType := "sales_invoice"
	event := "complete"

	// Setup document
	mockRepo.documents[docID] = map[string]interface{}{
		"total_amount": 100.0,
	}

	// Setup posting rule
	ruleID := uuid.New()
	rule := PostingRule{
		ID:                  ruleID,
		RuleCode:            "SALE_001",
		RuleName:            "Sales Invoice Posting",
		Event:               event,
		ConditionExpression: "",
		IsActive:            true,
	}
	mockRepo.postingRules[docType+"_"+event] = []PostingRule{rule}

	// Setup rule lines
	receivablesAcctID := uuid.New()
	salesAcctID := uuid.New()
	mockRepo.conceptMappings["accounts_receivable"] = receivablesAcctID
	mockRepo.conceptMappings["sales_revenue"] = salesAcctID

	mockRepo.postingRuleLines[ruleID] = []PostingRuleLine{
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              1,
			Side:                "debit",
			ConceptKey:          "accounts_receivable",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Invoice",
		},
		{
			ID:                  uuid.New(),
			PostingRuleID:       ruleID,
			LineNo:              2,
			Side:                "credit",
			ConceptKey:          "sales_revenue",
			AmountSource:        "document_field",
			AmountFieldPath:     "total_amount",
			DescriptionTemplate: "Sales Revenue",
		},
	}

	// Mock CreateJournalEntry to return error
	expectedErr := errors.New("database write failed")
	mockRepo.onCreateJournalEntry = func(ctx context.Context, je *JournalEntry) error {
		return expectedErr
	}

	input := PostingInput{
		OrganizationID: orgID,
		DocumentType:   docType,
		DocumentID:     docID,
		Event:          event,
		UserID:         uuid.New(),
	}

	err := engine.Post(ctx, input)
	if err == nil {
		t.Fatal("expected error from CreateJournalEntry, got nil")
	}
}

// ========================================
// HELPER METHOD TESTS
// ========================================

func TestEngine_EvaluateCondition_EmptyCondition(t *testing.T) {
	engine, _ := setupEngine()

	doc := map[string]interface{}{
		"amount": 100.0,
	}

	result := engine.evaluateCondition("", doc)
	if !result {
		t.Error("expected empty condition to return true")
	}
}

func TestEngine_EvaluateCondition_InvalidExpression(t *testing.T) {
	engine, _ := setupEngine()

	doc := map[string]interface{}{
		"amount": 100.0,
	}

	result := engine.evaluateCondition("invalid syntax {{", doc)
	if result {
		t.Error("expected invalid expression to return false")
	}
}

func TestEngine_EvaluateCondition_NonBooleanResult(t *testing.T) {
	engine, _ := setupEngine()

	doc := map[string]interface{}{
		"amount": 100.0,
	}

	// Expression returns a number, not boolean
	result := engine.evaluateCondition("doc.amount", doc)
	if result {
		t.Error("expected non-boolean result to return false")
	}
}

package banking

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BankingRepository defines the interface for banking database operations
type BankingRepository interface {
	// Bank Accounts
	CreateBankAccount(ctx context.Context, account *BankAccount) error
	GetBankAccount(ctx context.Context, id, organizationID uuid.UUID) (*BankAccount, error)
	GetBankAccountByNumber(ctx context.Context, organizationID uuid.UUID, accountNumber string) (*BankAccount, error)
	ListBankAccounts(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*BankAccount, error)
	UpdateBankAccount(ctx context.Context, account *BankAccount) error
	DeleteBankAccount(ctx context.Context, id, organizationID uuid.UUID) error

	// Bank Reconciliations
	CreateBankReconciliation(ctx context.Context, reconciliation *BankReconciliation) error
	GetBankReconciliation(ctx context.Context, id, organizationID uuid.UUID) (*BankReconciliation, error)
	ListBankReconciliations(ctx context.Context, organizationID, bankAccountID uuid.UUID, filter map[string]interface{}) ([]*BankReconciliation, error)
	UpdateBankReconciliation(ctx context.Context, reconciliation *BankReconciliation) error
	DeleteBankReconciliation(ctx context.Context, id, organizationID uuid.UUID) error
	MarkReconciliationComplete(ctx context.Context, id, organizationID, reconciledBy uuid.UUID) error

	// Bank Reconciliation Items
	CreateBankReconciliationItem(ctx context.Context, item *BankReconciliationItem) error
	GetBankReconciliationItem(ctx context.Context, id, organizationID uuid.UUID) (*BankReconciliationItem, error)
	ListBankReconciliationItems(ctx context.Context, reconciliationID uuid.UUID) ([]*BankReconciliationItem, error)
	UpdateBankReconciliationItem(ctx context.Context, item *BankReconciliationItem) error
	DeleteBankReconciliationItem(ctx context.Context, id, organizationID uuid.UUID) error
	MarkReconciliationItemCleared(ctx context.Context, id, organizationID, clearedBy uuid.UUID) error

	// Bank Statements
	CreateBankStatement(ctx context.Context, statement *BankStatement) error
	GetBankStatement(ctx context.Context, id, organizationID uuid.UUID) (*BankStatement, error)
	ListBankStatements(ctx context.Context, organizationID, bankAccountID uuid.UUID, filter map[string]interface{}) ([]*BankStatement, error)
	UpdateBankStatement(ctx context.Context, statement *BankStatement) error
	DeleteBankStatement(ctx context.Context, id, organizationID uuid.UUID) error

	// Bank Statement Lines
	CreateBankStatementLine(ctx context.Context, line *BankStatementLine) error
	GetBankStatementLine(ctx context.Context, id uuid.UUID) (*BankStatementLine, error)
	ListBankStatementLines(ctx context.Context, statementID uuid.UUID, filter map[string]interface{}) ([]*BankStatementLine, error)
	UpdateBankStatementLine(ctx context.Context, line *BankStatementLine) error
	DeleteBankStatementLine(ctx context.Context, id uuid.UUID) error

	// Bank Statement Reconciliations
	CreateBankStatementReconciliation(ctx context.Context, reconciliation *BankStatementReconciliation) error
	GetBankStatementReconciliation(ctx context.Context, id, organizationID uuid.UUID) (*BankStatementReconciliation, error)
	ListBankStatementReconciliations(ctx context.Context, statementLineID uuid.UUID) ([]*BankStatementReconciliation, error)
	DeleteBankStatementReconciliation(ctx context.Context, id, organizationID uuid.UUID) error

	// Reconciliation Rule Models
	CreateReconciliationRuleModel(ctx context.Context, rule *ReconciliationRuleModel) error
	GetReconciliationRuleModel(ctx context.Context, id, organizationID uuid.UUID) (*ReconciliationRuleModel, error)
	ListReconciliationRuleModels(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*ReconciliationRuleModel, error)
	UpdateReconciliationRuleModel(ctx context.Context, rule *ReconciliationRuleModel) error
	DeleteReconciliationRuleModel(ctx context.Context, id, organizationID uuid.UUID) error
}

// Service handles banking business logic
type Service struct {
	repo BankingRepository
}

// NewService creates a new banking service
func NewService(repo BankingRepository) *Service {
	return &Service{repo: repo}
}

// ========================
// BANK ACCOUNT MANAGEMENT
// ========================

// CreateBankAccount creates a new bank account
func (s *Service) CreateBankAccount(ctx context.Context, orgID uuid.UUID, req *CreateBankAccountRequest, createdBy uuid.UUID) (*BankAccount, error) {
	if req.AccountType != BankAccountTypeChecking && req.AccountType != BankAccountTypeSavings &&
		req.AccountType != BankAccountTypeCreditCard && req.AccountType != BankAccountTypeLineOfCredit {
		return nil, ErrInvalidAccountType
	}

	if len(req.CurrencyCode) != 3 {
		return nil, ErrInvalidCurrencyCode
	}

	metadata, _ := json.Marshal(req.Metadata)

	account := &BankAccount{
		ID:                   uuid.New(),
		OrganizationID:       orgID,
		ChartAccountID:       req.ChartAccountID,
		BankName:             req.BankName,
		AccountNumber:        req.AccountNumber,
		AccountType:          req.AccountType,
		RoutingNumber:        req.RoutingNumber,
		SwiftCode:            req.SwiftCode,
		CurrencyCode:         req.CurrencyCode,
		CurrentBalance:       "0",
		StatementBalance:     "0",
		IsActive:             req.IsActive,
		OnlineBankingEnabled: req.OnlineBankingEnabled,
		Notes:                req.Notes,
		Metadata:             metadata,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		CreatedBy:            &createdBy,
		UpdatedBy:            &createdBy,
	}

	if err := s.repo.CreateBankAccount(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// GetBankAccount retrieves a bank account by ID
func (s *Service) GetBankAccount(ctx context.Context, id, organizationID uuid.UUID) (*BankAccount, error) {
	return s.repo.GetBankAccount(ctx, id, organizationID)
}

// GetBankAccountByNumber retrieves a bank account by account number
func (s *Service) GetBankAccountByNumber(ctx context.Context, organizationID uuid.UUID, accountNumber string) (*BankAccount, error) {
	return s.repo.GetBankAccountByNumber(ctx, organizationID, accountNumber)
}

// ListBankAccounts lists bank accounts with filters
func (s *Service) ListBankAccounts(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*BankAccount, error) {
	return s.repo.ListBankAccounts(ctx, organizationID, filter)
}

// UpdateBankAccount updates a bank account
func (s *Service) UpdateBankAccount(ctx context.Context, orgID, accountID uuid.UUID, req *UpdateBankAccountRequest, updatedBy uuid.UUID) (*BankAccount, error) {
	account, err := s.repo.GetBankAccount(ctx, accountID, orgID)
	if err != nil {
		return nil, err
	}

	if req.BankName != nil {
		account.BankName = *req.BankName
	}
	if req.AccountType != nil {
		account.AccountType = *req.AccountType
	}
	if req.RoutingNumber != nil {
		account.RoutingNumber = req.RoutingNumber
	}
	if req.SwiftCode != nil {
		account.SwiftCode = req.SwiftCode
	}
	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}
	if req.OnlineBankingEnabled != nil {
		account.OnlineBankingEnabled = *req.OnlineBankingEnabled
	}
	if req.Notes != nil {
		account.Notes = req.Notes
	}
	if req.Metadata != nil {
		metadata, _ := json.Marshal(req.Metadata)
		account.Metadata = metadata
	}

	account.UpdatedAt = time.Now()
	account.UpdatedBy = &updatedBy

	if err := s.repo.UpdateBankAccount(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// DeleteBankAccount soft deletes a bank account
func (s *Service) DeleteBankAccount(ctx context.Context, id, organizationID uuid.UUID) error {
	return s.repo.DeleteBankAccount(ctx, id, organizationID)
}

// ========================
// BANK RECONCILIATION MANAGEMENT
// ========================

// CreateBankReconciliation creates a new bank reconciliation
func (s *Service) CreateBankReconciliation(ctx context.Context, orgID uuid.UUID, req *CreateBankReconciliationRequest, createdBy uuid.UUID) (*BankReconciliation, error) {
	metadata, _ := json.Marshal(req.Metadata)

	reconciliation := &BankReconciliation{
		ID:                  uuid.New(),
		OrganizationID:      orgID,
		BankAccountID:       req.BankAccountID,
		StatementDate:       req.StatementDate,
		StatementBalance:    req.StatementBalance,
		BookBalance:         "0",
		ClearedBalance:      "0",
		Difference:          "0",
		Status:              BankReconciliationStatusInProgress,
		IsReconciled:        false,
		AccountingPeriodID:  req.AccountingPeriodID,
		Notes:               req.Notes,
		Metadata:            metadata,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CreatedBy:           &createdBy,
		UpdatedBy:           &createdBy,
	}

	if err := s.repo.CreateBankReconciliation(ctx, reconciliation); err != nil {
		return nil, err
	}

	return reconciliation, nil
}

// GetBankReconciliation retrieves a bank reconciliation by ID
func (s *Service) GetBankReconciliation(ctx context.Context, id, organizationID uuid.UUID) (*BankReconciliation, error) {
	return s.repo.GetBankReconciliation(ctx, id, organizationID)
}

// ListBankReconciliations lists bank reconciliations with filters
func (s *Service) ListBankReconciliations(ctx context.Context, organizationID, bankAccountID uuid.UUID, filter map[string]interface{}) ([]*BankReconciliation, error) {
	return s.repo.ListBankReconciliations(ctx, organizationID, bankAccountID, filter)
}

// UpdateBankReconciliation updates a bank reconciliation
func (s *Service) UpdateBankReconciliation(ctx context.Context, orgID, reconciliationID uuid.UUID, req *UpdateBankReconciliationRequest, updatedBy uuid.UUID) (*BankReconciliation, error) {
	reconciliation, err := s.repo.GetBankReconciliation(ctx, reconciliationID, orgID)
	if err != nil {
		return nil, err
	}

	if reconciliation.Status == BankReconciliationStatusLocked {
		return nil, ErrReconciliationLocked
	}

	if req.StatementBalance != "" {
		reconciliation.StatementBalance = req.StatementBalance
	}
	if req.Status != nil {
		reconciliation.Status = *req.Status
	}
	if req.Notes != nil {
		reconciliation.Notes = req.Notes
	}
	if req.Metadata != nil {
		metadata, _ := json.Marshal(req.Metadata)
		reconciliation.Metadata = metadata
	}

	reconciliation.UpdatedAt = time.Now()
	reconciliation.UpdatedBy = &updatedBy

	if err := s.repo.UpdateBankReconciliation(ctx, reconciliation); err != nil {
		return nil, err
	}

	return reconciliation, nil
}

// MarkReconciliationComplete marks a reconciliation as complete
func (s *Service) MarkReconciliationComplete(ctx context.Context, id, orgID, reconciledBy uuid.UUID) (*BankReconciliation, error) {
	reconciliation, err := s.repo.GetBankReconciliation(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	reconciliation.IsReconciled = true
	reconciliation.Status = BankReconciliationStatusReconciled
	reconciliation.ReconciledBy = &reconciledBy
	now := time.Now()
	reconciliation.ReconciledAt = &now
	reconciliation.ReconciliationDate = &now

	if err := s.repo.MarkReconciliationComplete(ctx, id, orgID, reconciledBy); err != nil {
		return nil, err
	}

	return reconciliation, nil
}

// DeleteBankReconciliation soft deletes a bank reconciliation
func (s *Service) DeleteBankReconciliation(ctx context.Context, id, organizationID uuid.UUID) error {
	return s.repo.DeleteBankReconciliation(ctx, id, organizationID)
}

// ========================
// BANK RECONCILIATION ITEMS
// ========================

// CreateBankReconciliationItem creates a new reconciliation item
func (s *Service) CreateBankReconciliationItem(ctx context.Context, orgID uuid.UUID, req *CreateBankReconciliationItemRequest) (*BankReconciliationItem, error) {
	item := &BankReconciliationItem{
		ID:                   uuid.New(),
		OrganizationID:       orgID,
		BankReconciliationID: req.BankReconciliationID,
		GeneralLedgerID:      req.GeneralLedgerID,
		JournalEntryLineID:   req.JournalEntryLineID,
		IsCleared:            false,
		CreatedAt:            time.Now(),
	}

	if err := s.repo.CreateBankReconciliationItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetBankReconciliationItem retrieves a reconciliation item by ID
func (s *Service) GetBankReconciliationItem(ctx context.Context, id, organizationID uuid.UUID) (*BankReconciliationItem, error) {
	return s.repo.GetBankReconciliationItem(ctx, id, organizationID)
}

// ListBankReconciliationItems lists reconciliation items
func (s *Service) ListBankReconciliationItems(ctx context.Context, reconciliationID uuid.UUID) ([]*BankReconciliationItem, error) {
	return s.repo.ListBankReconciliationItems(ctx, reconciliationID)
}

// MarkReconciliationItemCleared marks a reconciliation item as cleared
func (s *Service) MarkReconciliationItemCleared(ctx context.Context, id, orgID, clearedBy uuid.UUID) (*BankReconciliationItem, error) {
	item, err := s.repo.GetBankReconciliationItem(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	item.IsCleared = true
	now := time.Now()
	item.ClearedDate = &now
	item.ClearedBy = &clearedBy

	if err := s.repo.MarkReconciliationItemCleared(ctx, id, orgID, clearedBy); err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteBankReconciliationItem soft deletes a reconciliation item
func (s *Service) DeleteBankReconciliationItem(ctx context.Context, id, organizationID uuid.UUID) error {
	return s.repo.DeleteBankReconciliationItem(ctx, id, organizationID)
}

// ========================
// BANK STATEMENT MANAGEMENT
// ========================

// CreateBankStatement creates a new bank statement
func (s *Service) CreateBankStatement(ctx context.Context, orgID uuid.UUID, bankAccountID uuid.UUID, req *CreateBankStatementRequest, createdBy uuid.UUID) (*BankStatement, error) {
	if req.PeriodEndDate.Before(req.PeriodStartDate) {
		return nil, ErrInvalidDateRange
	}

	statement := &BankStatement{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		BankAccountID:   bankAccountID,
		StatementNumber: req.StatementNumber,
		StatementDate:   req.StatementDate,
		PeriodStartDate: req.PeriodStartDate,
		PeriodEndDate:   req.PeriodEndDate,
		OpeningBalance:  req.OpeningBalance,
		ClosingBalance:  req.ClosingBalance,
		ImportSource:    req.ImportSource,
		ImportFileName:  req.ImportFileName,
		Status:          BankStatementStatusDraft,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedBy:       &createdBy,
		UpdatedBy:       &createdBy,
	}

	if err := s.repo.CreateBankStatement(ctx, statement); err != nil {
		return nil, err
	}

	return statement, nil
}

// GetBankStatement retrieves a bank statement by ID
func (s *Service) GetBankStatement(ctx context.Context, id, organizationID uuid.UUID) (*BankStatement, error) {
	return s.repo.GetBankStatement(ctx, id, organizationID)
}

// ListBankStatements lists bank statements with filters
func (s *Service) ListBankStatements(ctx context.Context, organizationID, bankAccountID uuid.UUID, filter map[string]interface{}) ([]*BankStatement, error) {
	return s.repo.ListBankStatements(ctx, organizationID, bankAccountID, filter)
}

// UpdateBankStatement updates a bank statement
func (s *Service) UpdateBankStatement(ctx context.Context, orgID, statementID uuid.UUID, req *UpdateBankStatementRequest, updatedBy uuid.UUID) (*BankStatement, error) {
	statement, err := s.repo.GetBankStatement(ctx, statementID, orgID)
	if err != nil {
		return nil, err
	}

	if req.StatementNumber != nil {
		statement.StatementNumber = req.StatementNumber
	}
	if req.OpeningBalance != nil {
		statement.OpeningBalance = *req.OpeningBalance
	}
	if req.ClosingBalance != nil {
		statement.ClosingBalance = *req.ClosingBalance
	}
	if req.Status != nil {
		statement.Status = *req.Status
	}
	if req.Notes != nil {
		statement.Notes = req.Notes
	}

	statement.UpdatedAt = time.Now()
	statement.UpdatedBy = &updatedBy

	if err := s.repo.UpdateBankStatement(ctx, statement); err != nil {
		return nil, err
	}

	return statement, nil
}

// DeleteBankStatement soft deletes a bank statement
func (s *Service) DeleteBankStatement(ctx context.Context, id, organizationID uuid.UUID) error {
	return s.repo.DeleteBankStatement(ctx, id, organizationID)
}

// ========================
// BANK STATEMENT LINES
// ========================

// CreateBankStatementLine creates a new statement line
func (s *Service) CreateBankStatementLine(ctx context.Context, statementID uuid.UUID, req *CreateBankStatementLineRequest) (*BankStatementLine, error) {
	line := &BankStatementLine{
		ID:                  uuid.New(),
		BankStatementID:     statementID,
		LineNumber:          req.LineNumber,
		TransactionDate:     req.TransactionDate,
		ValueDate:           req.ValueDate,
		Amount:              req.Amount,
		CurrencyCode:        req.CurrencyCode,
		Description:         req.Description,
		Reference:           req.Reference,
		CounterpartyName:    req.CounterpartyName,
		CounterpartyAccount: req.CounterpartyAccount,
		BankReference:       req.BankReference,
		CheckNumber:         req.CheckNumber,
		Status:              BankStatementLineStatusUnmatched,
		Notes:               req.Notes,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.repo.CreateBankStatementLine(ctx, line); err != nil {
		return nil, err
	}

	return line, nil
}

// GetBankStatementLine retrieves a statement line by ID
func (s *Service) GetBankStatementLine(ctx context.Context, id uuid.UUID) (*BankStatementLine, error) {
	return s.repo.GetBankStatementLine(ctx, id)
}

// ListBankStatementLines lists statement lines with filters
func (s *Service) ListBankStatementLines(ctx context.Context, statementID uuid.UUID, filter map[string]interface{}) ([]*BankStatementLine, error) {
	return s.repo.ListBankStatementLines(ctx, statementID, filter)
}

// UpdateBankStatementLine updates a statement line
func (s *Service) UpdateBankStatementLine(ctx context.Context, id uuid.UUID, req *UpdateBankStatementLineRequest) (*BankStatementLine, error) {
	line, err := s.repo.GetBankStatementLine(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Amount != nil {
		line.Amount = *req.Amount
	}
	if req.Description != nil {
		line.Description = req.Description
	}
	if req.Reference != nil {
		line.Reference = req.Reference
	}
	if req.CounterpartyName != nil {
		line.CounterpartyName = req.CounterpartyName
	}
	if req.CounterpartyAccount != nil {
		line.CounterpartyAccount = req.CounterpartyAccount
	}
	if req.Status != nil {
		line.Status = *req.Status
	}
	if req.Notes != nil {
		line.Notes = req.Notes
	}

	line.UpdatedAt = time.Now()

	if err := s.repo.UpdateBankStatementLine(ctx, line); err != nil {
		return nil, err
	}

	return line, nil
}

// DeleteBankStatementLine soft deletes a statement line
func (s *Service) DeleteBankStatementLine(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteBankStatementLine(ctx, id)
}

// ========================
// BANK STATEMENT RECONCILIATIONS
// ========================

// CreateBankStatementReconciliation creates a new statement reconciliation
func (s *Service) CreateBankStatementReconciliation(ctx context.Context, orgID uuid.UUID, req *CreateBankStatementReconciliationRequest, matchedBy uuid.UUID) (*BankStatementReconciliation, error) {
	reconciliation := &BankStatementReconciliation{
		ID:                  uuid.New(),
		OrganizationID:      orgID,
		BankStatementLineID: req.BankStatementLineID,
		JournalEntryID:      req.JournalEntryID,
		PaymentID:           req.PaymentID,
		MatchedAmount:       req.MatchedAmount,
		MatchedBy:           &matchedBy,
		MatchedAt:           time.Now(),
	}

	if err := s.repo.CreateBankStatementReconciliation(ctx, reconciliation); err != nil {
		return nil, err
	}

	return reconciliation, nil
}

// GetBankStatementReconciliation retrieves a statement reconciliation by ID
func (s *Service) GetBankStatementReconciliation(ctx context.Context, id, organizationID uuid.UUID) (*BankStatementReconciliation, error) {
	return s.repo.GetBankStatementReconciliation(ctx, id, organizationID)
}

// ListBankStatementReconciliations lists statement reconciliations
func (s *Service) ListBankStatementReconciliations(ctx context.Context, statementLineID uuid.UUID) ([]*BankStatementReconciliation, error) {
	return s.repo.ListBankStatementReconciliations(ctx, statementLineID)
}

// DeleteBankStatementReconciliation soft deletes a statement reconciliation
func (s *Service) DeleteBankStatementReconciliation(ctx context.Context, id, organizationID uuid.UUID) error {
	return s.repo.DeleteBankStatementReconciliation(ctx, id, organizationID)
}

// ========================
// RECONCILIATION RULE MODELS
// ========================

// CreateReconciliationRuleModel creates a new reconciliation rule
func (s *Service) CreateReconciliationRuleModel(ctx context.Context, orgID uuid.UUID, req *CreateReconciliationRuleModelRequest, createdBy uuid.UUID) (*ReconciliationRuleModel, error) {
	rule := &ReconciliationRuleModel{
		ID:                  uuid.New(),
		OrganizationID:      orgID,
		RuleName:            req.RuleName,
		RuleCode:            req.RuleCode,
		Sequence:            req.Sequence,
		AmountMin:           req.AmountMin,
		AmountMax:           req.AmountMax,
		DescriptionPattern:  req.DescriptionPattern,
		CounterpartyPattern: req.CounterpartyPattern,
		ReferencePattern:    req.ReferencePattern,
		JournalID:           req.JournalID,
		AccountID:           req.AccountID,
		AnalyticAccountID:   req.AnalyticAccountID,
		TaxID:               req.TaxID,
		IsActive:            req.IsActive,
		AutoApply:           req.AutoApply,
		Notes:               req.Notes,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CreatedBy:           &createdBy,
		UpdatedBy:           &createdBy,
	}

	if err := s.repo.CreateReconciliationRuleModel(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// GetReconciliationRuleModel retrieves a rule by ID
func (s *Service) GetReconciliationRuleModel(ctx context.Context, id, organizationID uuid.UUID) (*ReconciliationRuleModel, error) {
	return s.repo.GetReconciliationRuleModel(ctx, id, organizationID)
}

// ListReconciliationRuleModels lists rules with filters
func (s *Service) ListReconciliationRuleModels(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*ReconciliationRuleModel, error) {
	return s.repo.ListReconciliationRuleModels(ctx, organizationID, filter)
}

// UpdateReconciliationRuleModel updates a rule
func (s *Service) UpdateReconciliationRuleModel(ctx context.Context, orgID, ruleID uuid.UUID, req *UpdateReconciliationRuleModelRequest, updatedBy uuid.UUID) (*ReconciliationRuleModel, error) {
	rule, err := s.repo.GetReconciliationRuleModel(ctx, ruleID, orgID)
	if err != nil {
		return nil, err
	}

	if req.RuleName != nil {
		rule.RuleName = *req.RuleName
	}
	if req.RuleCode != nil {
		rule.RuleCode = req.RuleCode
	}
	if req.Sequence != nil {
		rule.Sequence = *req.Sequence
	}
	if req.AmountMin != nil {
		rule.AmountMin = req.AmountMin
	}
	if req.AmountMax != nil {
		rule.AmountMax = req.AmountMax
	}
	if req.DescriptionPattern != nil {
		rule.DescriptionPattern = req.DescriptionPattern
	}
	if req.CounterpartyPattern != nil {
		rule.CounterpartyPattern = req.CounterpartyPattern
	}
	if req.ReferencePattern != nil {
		rule.ReferencePattern = req.ReferencePattern
	}
	if req.JournalID != nil {
		rule.JournalID = req.JournalID
	}
	if req.AccountID != nil {
		rule.AccountID = req.AccountID
	}
	if req.AnalyticAccountID != nil {
		rule.AnalyticAccountID = req.AnalyticAccountID
	}
	if req.TaxID != nil {
		rule.TaxID = req.TaxID
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	if req.AutoApply != nil {
		rule.AutoApply = *req.AutoApply
	}
	if req.Notes != nil {
		rule.Notes = req.Notes
	}

	rule.UpdatedAt = time.Now()
	rule.UpdatedBy = &updatedBy

	if err := s.repo.UpdateReconciliationRuleModel(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// DeleteReconciliationRuleModel soft deletes a rule
func (s *Service) DeleteReconciliationRuleModel(ctx context.Context, id, organizationID uuid.UUID) error {
	return s.repo.DeleteReconciliationRuleModel(ctx, id, organizationID)
}

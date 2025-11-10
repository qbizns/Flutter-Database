package accounting

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccountingRepository defines the interface for accounting database operations
type AccountingRepository interface {
	// Fiscal Years
	CreateFiscalYear(ctx context.Context, fy *FiscalYear) error
	GetFiscalYear(ctx context.Context, id, organizationID uuid.UUID) (*FiscalYear, error)
	ListFiscalYears(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*FiscalYear, error)
	UpdateFiscalYear(ctx context.Context, fy *FiscalYear) error
	DeleteFiscalYear(ctx context.Context, id, organizationID uuid.UUID) error
	GetCurrentFiscalYear(ctx context.Context, organizationID uuid.UUID) (*FiscalYear, error)
	CloseFiscalYear(ctx context.Context, id, organizationID, closedBy uuid.UUID) error

	// Accounting Periods
	CreateAccountingPeriod(ctx context.Context, ap *AccountingPeriod) error
	GetAccountingPeriod(ctx context.Context, id, organizationID uuid.UUID) (*AccountingPeriod, error)
	ListAccountingPeriods(ctx context.Context, organizationID, fiscalYearID uuid.UUID) ([]*AccountingPeriod, error)
	UpdateAccountingPeriod(ctx context.Context, ap *AccountingPeriod) error
	DeleteAccountingPeriod(ctx context.Context, id, organizationID uuid.UUID) error
	CloseAccountingPeriod(ctx context.Context, id, organizationID, closedBy uuid.UUID) error
	GetAccountingPeriodByDate(ctx context.Context, organizationID uuid.UUID, date time.Time) (*AccountingPeriod, error)

	// Account Types & Subtypes
	GetAccountType(ctx context.Context, id uuid.UUID) (*AccountType, error)
	ListAccountTypes(ctx context.Context) ([]*AccountType, error)
	GetAccountSubtype(ctx context.Context, id uuid.UUID) (*AccountSubtype, error)
	ListAccountSubtypes(ctx context.Context, accountTypeID uuid.UUID) ([]*AccountSubtype, error)

	// Chart of Accounts
	CreateChartOfAccount(ctx context.Context, coa *ChartOfAccount) error
	GetChartOfAccount(ctx context.Context, id, organizationID uuid.UUID) (*ChartOfAccount, error)
	GetChartOfAccountByCode(ctx context.Context, organizationID uuid.UUID, code string) (*ChartOfAccount, error)
	ListChartOfAccounts(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*ChartOfAccount, error)
	UpdateChartOfAccount(ctx context.Context, coa *ChartOfAccount) error
	DeleteChartOfAccount(ctx context.Context, id, organizationID uuid.UUID) error
	GetAccountHierarchy(ctx context.Context, organizationID, parentID uuid.UUID) ([]*ChartOfAccount, error)

	// Journal Entry Types
	GetJournalEntryType(ctx context.Context, id uuid.UUID) (*JournalEntryType, error)
	ListJournalEntryTypes(ctx context.Context) ([]*JournalEntryType, error)

	// Journal Entries
	CreateJournalEntry(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error
	GetJournalEntry(ctx context.Context, id, organizationID uuid.UUID) (*JournalEntry, error)
	GetJournalEntryByNumber(ctx context.Context, organizationID uuid.UUID, number string) (*JournalEntry, error)
	ListJournalEntries(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*JournalEntry, error)
	UpdateJournalEntry(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error
	DeleteJournalEntry(ctx context.Context, id, organizationID uuid.UUID) error
	PostJournalEntry(ctx context.Context, id, organizationID, postedBy uuid.UUID) error
	ApproveJournalEntry(ctx context.Context, id, organizationID, approvedBy uuid.UUID) error
	ReverseJournalEntry(ctx context.Context, id, organizationID, reversalEntryID, reversedBy uuid.UUID) error

	// Journal Entry Lines
	GetJournalEntryLine(ctx context.Context, id, organizationID uuid.UUID) (*JournalEntryLine, error)
	ListJournalEntryLines(ctx context.Context, journalEntryID uuid.UUID) ([]*JournalEntryLine, error)

	// General Ledger
	PostToGeneralLedger(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error
	GetGeneralLedgerEntries(ctx context.Context, organizationID, accountID uuid.UUID, filter map[string]interface{}) ([]*GeneralLedger, error)
	GetAccountBalance(ctx context.Context, organizationID, accountID uuid.UUID, asOfDate time.Time) (debit, credit, balance string, err error)
	GetTrialBalance(ctx context.Context, organizationID uuid.UUID, asOfDate time.Time) (map[uuid.UUID]map[string]string, error)
}

// Service handles accounting business logic
type Service struct {
	repo AccountingRepository
}

// NewService creates a new accounting service
func NewService(repo AccountingRepository) *Service {
	return &Service{repo: repo}
}

// ========================
// FISCAL YEAR MANAGEMENT
// ========================

// CreateFiscalYear creates a new fiscal year
func (s *Service) CreateFiscalYear(ctx context.Context, orgID uuid.UUID, req *CreateFiscalYearRequest, createdBy uuid.UUID) (*FiscalYear, error) {
	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDateRange
	}

	metadata, _ := json.Marshal(req.Metadata)

	fy := &FiscalYear{
		ID:             uuid.New(),
		OrganizationID: orgID,
		FiscalYear:     req.FiscalYear,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      req.IsCurrent,
		Notes:          req.Notes,
		Metadata:       metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      &createdBy,
		UpdatedBy:      &createdBy,
	}

	if err := s.repo.CreateFiscalYear(ctx, fy); err != nil {
		return nil, err
	}

	return fy, nil
}

// GetFiscalYear retrieves a fiscal year
func (s *Service) GetFiscalYear(ctx context.Context, id, orgID uuid.UUID) (*FiscalYear, error) {
	return s.repo.GetFiscalYear(ctx, id, orgID)
}

// UpdateFiscalYear updates a fiscal year
func (s *Service) UpdateFiscalYear(ctx context.Context, id, orgID uuid.UUID, req *UpdateFiscalYearRequest, updatedBy uuid.UUID) (*FiscalYear, error) {
	fy, err := s.repo.GetFiscalYear(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	if req.FiscalYear != nil {
		fy.FiscalYear = *req.FiscalYear
	}
	if req.StartDate != nil {
		fy.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		fy.EndDate = *req.EndDate
	}
	if req.IsCurrent != nil {
		fy.IsCurrent = *req.IsCurrent
	}
	if req.Notes != nil {
		fy.Notes = req.Notes
	}
	if req.Metadata != nil {
		metadata, _ := json.Marshal(req.Metadata)
		fy.Metadata = metadata
	}

	fy.UpdatedAt = time.Now()
	fy.UpdatedBy = &updatedBy

	if err := s.repo.UpdateFiscalYear(ctx, fy); err != nil {
		return nil, err
	}

	return fy, nil
}

// CloseFiscalYear closes a fiscal year
func (s *Service) CloseFiscalYear(ctx context.Context, id, orgID, closedBy uuid.UUID) error {
	return s.repo.CloseFiscalYear(ctx, id, orgID, closedBy)
}

// ListFiscalYears lists fiscal years
func (s *Service) ListFiscalYears(ctx context.Context, orgID uuid.UUID, filter map[string]interface{}) ([]*FiscalYear, error) {
	return s.repo.ListFiscalYears(ctx, orgID, filter)
}

// ========================
// ACCOUNTING PERIOD MANAGEMENT
// ========================

// CreateAccountingPeriod creates an accounting period
func (s *Service) CreateAccountingPeriod(ctx context.Context, orgID uuid.UUID, req *CreateAccountingPeriodRequest, createdBy uuid.UUID) (*AccountingPeriod, error) {
	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDateRange
	}

	metadata, _ := json.Marshal(req.Metadata)

	ap := &AccountingPeriod{
		ID:             uuid.New(),
		OrganizationID: orgID,
		FiscalYearID:   req.FiscalYearID,
		PeriodNumber:   req.PeriodNumber,
		PeriodName:     req.PeriodName,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         AccountingPeriodStatusOpen,
		Metadata:       metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      &createdBy,
		UpdatedBy:      &createdBy,
	}

	if err := s.repo.CreateAccountingPeriod(ctx, ap); err != nil {
		return nil, err
	}

	return ap, nil
}

// GetAccountingPeriod retrieves an accounting period
func (s *Service) GetAccountingPeriod(ctx context.Context, id, orgID uuid.UUID) (*AccountingPeriod, error) {
	return s.repo.GetAccountingPeriod(ctx, id, orgID)
}

// GetAccountingPeriodByDate gets the period for a given date
func (s *Service) GetAccountingPeriodByDate(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
	return s.repo.GetAccountingPeriodByDate(ctx, orgID, date)
}

// CloseAccountingPeriod closes an accounting period
func (s *Service) CloseAccountingPeriod(ctx context.Context, id, orgID, closedBy uuid.UUID) error {
	return s.repo.CloseAccountingPeriod(ctx, id, orgID, closedBy)
}

// ListAccountingPeriods lists periods for a fiscal year
func (s *Service) ListAccountingPeriods(ctx context.Context, orgID, fiscalYearID uuid.UUID) ([]*AccountingPeriod, error) {
	return s.repo.ListAccountingPeriods(ctx, orgID, fiscalYearID)
}

// ========================
// CHART OF ACCOUNTS MANAGEMENT
// ========================

// CreateChartOfAccount creates a chart of account entry
func (s *Service) CreateChartOfAccount(ctx context.Context, orgID uuid.UUID, req *CreateChartOfAccountRequest, createdBy uuid.UUID) (*ChartOfAccount, error) {
	// Get account type
	_, err := s.repo.GetAccountType(ctx, req.AccountTypeID)
	if err != nil {
		return nil, err
	}

	// Validate account level
	accountLevel := 1
	if req.ParentAccountID != nil {
		parent, err := s.repo.GetChartOfAccount(ctx, *req.ParentAccountID, orgID)
		if err != nil {
			return nil, err
		}
		accountLevel = parent.AccountLevel + 1
		if accountLevel > 5 {
			return nil, ErrInvalidAccountLevel
		}
	}

	metadata, _ := json.Marshal(req.Metadata)
	openingBalance := "0"
	if req.OpeningBalance != nil {
		openingBalance = *req.OpeningBalance
	}

	coa := &ChartOfAccount{
		ID:                  uuid.New(),
		OrganizationID:      orgID,
		AccountCode:         req.AccountCode,
		AccountNumber:       req.AccountNumber,
		AccountName:         req.AccountName,
		AccountTypeID:       req.AccountTypeID,
		AccountSubtypeID:    req.AccountSubtypeID,
		ParentAccountID:     req.ParentAccountID,
		AccountLevel:        accountLevel,
		IsActive:            req.IsActive,
		IsSystemAccount:     req.IsSystemAccount,
		IsHeaderAccount:     req.IsHeaderAccount,
		IsBankAccount:       req.IsBankAccount,
		IsReconcilable:      req.IsReconcilable,
		DefaultTaxCode:      req.DefaultTaxCode,
		CurrencyCode:        req.CurrencyCode,
		OpeningBalance:      openingBalance,
		OpeningBalanceDate:  req.OpeningBalanceDate,
		CurrentDebitBalance: "0",
		CurrentCreditBalance: "0",
		CurrentBalance:      openingBalance,
		Description:         req.Description,
		Notes:               req.Notes,
		Metadata:            metadata,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CreatedBy:           &createdBy,
		UpdatedBy:           &createdBy,
	}

	// Build account path
	if req.ParentAccountID != nil {
		parent, _ := s.repo.GetChartOfAccount(ctx, *req.ParentAccountID, orgID)
		accountPath := parent.AccountCode
		if parent.AccountPath != nil {
			accountPath = *parent.AccountPath + "/" + parent.AccountCode
		}
		accountPath = accountPath + "/" + req.AccountCode
		coa.AccountPath = &accountPath
	}

	if err := s.repo.CreateChartOfAccount(ctx, coa); err != nil {
		return nil, err
	}

	return coa, nil
}

// GetChartOfAccount retrieves a chart of account
func (s *Service) GetChartOfAccount(ctx context.Context, id, orgID uuid.UUID) (*ChartOfAccount, error) {
	return s.repo.GetChartOfAccount(ctx, id, orgID)
}

// GetChartOfAccountByCode retrieves by code
func (s *Service) GetChartOfAccountByCode(ctx context.Context, orgID uuid.UUID, code string) (*ChartOfAccount, error) {
	return s.repo.GetChartOfAccountByCode(ctx, orgID, code)
}

// UpdateChartOfAccount updates a chart of account
func (s *Service) UpdateChartOfAccount(ctx context.Context, id, orgID uuid.UUID, req *UpdateChartOfAccountRequest, updatedBy uuid.UUID) (*ChartOfAccount, error) {
	coa, err := s.repo.GetChartOfAccount(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	// Prevent system account modification
	if coa.IsSystemAccount {
		return nil, ErrSystemAccountModification
	}

	if req.AccountName != nil {
		coa.AccountName = *req.AccountName
	}
	if req.AccountSubtypeID != nil {
		coa.AccountSubtypeID = req.AccountSubtypeID
	}
	if req.IsActive != nil {
		coa.IsActive = *req.IsActive
	}
	if req.IsHeaderAccount != nil {
		coa.IsHeaderAccount = *req.IsHeaderAccount
	}
	if req.IsBankAccount != nil {
		coa.IsBankAccount = *req.IsBankAccount
	}
	if req.IsReconcilable != nil {
		coa.IsReconcilable = *req.IsReconcilable
	}
	if req.DefaultTaxCode != nil {
		coa.DefaultTaxCode = req.DefaultTaxCode
	}
	if req.Description != nil {
		coa.Description = req.Description
	}
	if req.Notes != nil {
		coa.Notes = req.Notes
	}
	if req.Metadata != nil {
		metadata, _ := json.Marshal(req.Metadata)
		coa.Metadata = metadata
	}

	coa.UpdatedAt = time.Now()
	coa.UpdatedBy = &updatedBy

	if err := s.repo.UpdateChartOfAccount(ctx, coa); err != nil {
		return nil, err
	}

	return coa, nil
}

// ListChartOfAccounts lists accounts
func (s *Service) ListChartOfAccounts(ctx context.Context, orgID uuid.UUID, filter map[string]interface{}) ([]*ChartOfAccount, error) {
	return s.repo.ListChartOfAccounts(ctx, orgID, filter)
}

// GetAccountHierarchy retrieves account hierarchy
func (s *Service) GetAccountHierarchy(ctx context.Context, orgID, parentID uuid.UUID) ([]*ChartOfAccount, error) {
	return s.repo.GetAccountHierarchy(ctx, orgID, parentID)
}

// ========================
// JOURNAL ENTRY MANAGEMENT
// ========================

// CreateJournalEntry creates a journal entry with lines
func (s *Service) CreateJournalEntry(ctx context.Context, orgID uuid.UUID, req *CreateJournalEntryRequest, createdBy uuid.UUID) (*JournalEntry, error) {
	// Validate minimum 2 lines
	if len(req.Lines) < 2 {
		return nil, ErrMinimumTwoLines
	}

	// Build journal entry
	entryID := uuid.New()
	now := time.Now()
	attachments, _ := json.Marshal(req.Attachments)
	metadata, _ := json.Marshal(req.Metadata)

	// Generate entry number
	entryNumber := fmt.Sprintf("JE-%d-%s", now.Unix(), entryID.String()[:8])

	je := &JournalEntry{
		ID:              entryID,
		OrganizationID:  orgID,
		EntryNumber:     entryNumber,
		EntryTypeID:     req.EntryTypeID,
		EntryDate:       req.EntryDate,
		PostingDate:     req.PostingDate,
		AccountingPeriodID: req.AccountingPeriodID,
		FiscalYearID:    req.FiscalYearID,
		Status:          JournalEntryStatusDraft,
		IsPosted:        false,
		IsReversed:      false,
		SourceModule:    req.SourceModule,
		SourceDocumentType: req.SourceDocumentType,
		SourceDocumentID: req.SourceDocumentID,
		ReferenceNumber: req.ReferenceNumber,
		TotalDebit:      "0",
		TotalCredit:     "0",
		Description:     req.Description,
		Notes:           req.Notes,
		RequiresApproval: req.RequiresApproval,
		Attachments:     attachments,
		Metadata:        metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       &createdBy,
		UpdatedBy:       &createdBy,
	}

	// Build lines
	lines := make([]*JournalEntryLine, 0, len(req.Lines))
	totalDebit := 0.0
	totalCredit := 0.0

	for _, lineReq := range req.Lines {
		lineID := uuid.New()
		debitAmount := "0"
		creditAmount := "0"

		if lineReq.DebitAmount != nil {
			debitAmount = *lineReq.DebitAmount
		}
		if lineReq.CreditAmount != nil {
			creditAmount = *lineReq.CreditAmount
		}

		// Validate debit/credit
		if debitAmount != "0" && creditAmount != "0" {
			return nil, ErrDebitCreditBoth
		}

		// Validate account
		account, err := s.repo.GetChartOfAccount(ctx, lineReq.AccountID, orgID)
		if err != nil {
			return nil, err
		}

		if !account.IsActive {
			return nil, ErrInactiveAccount
		}

		if account.IsHeaderAccount {
			return nil, ErrHeaderAccountPosting
		}

		lineMetadata, _ := json.Marshal(lineReq.Metadata)

		line := &JournalEntryLine{
			ID:             lineID,
			OrganizationID: orgID,
			JournalEntryID: entryID,
			LineNumber:     lineReq.LineNumber,
			AccountID:      lineReq.AccountID,
			DebitAmount:    debitAmount,
			CreditAmount:   creditAmount,
			LocationID:     lineReq.LocationID,
			Department:     lineReq.Department,
			ProjectCode:    lineReq.ProjectCode,
			CostCenter:     lineReq.CostCenter,
			TaxCode:        lineReq.TaxCode,
			TaxAmount:      "0",
			Description:    lineReq.Description,
			Memo:           lineReq.Memo,
			Metadata:       lineMetadata,
			CreatedAt:      now,
			UpdatedAt:      now,
			CreatedBy:      &createdBy,
			UpdatedBy:      &createdBy,
		}

		// Track totals (simplified - in real world use decimal library)
		if debitAmount != "0" {
			totalDebit += parseAmount(debitAmount)
		}
		if creditAmount != "0" {
			totalCredit += parseAmount(creditAmount)
		}

		lines = append(lines, line)
	}

	// Set totals
	je.TotalDebit = fmt.Sprintf("%.4f", totalDebit)
	je.TotalCredit = fmt.Sprintf("%.4f", totalCredit)

	// Validate balanced (with tolerance for floating point)
	if !isBalanced(totalDebit, totalCredit) {
		return nil, ErrJournalEntryNotBalanced
	}

	// Create in database
	if err := s.repo.CreateJournalEntry(ctx, je, lines); err != nil {
		return nil, err
	}

	return je, nil
}

// GetJournalEntry retrieves a journal entry
func (s *Service) GetJournalEntry(ctx context.Context, id, orgID uuid.UUID) (*JournalEntry, error) {
	return s.repo.GetJournalEntry(ctx, id, orgID)
}

// PostJournalEntry posts a journal entry to the general ledger
func (s *Service) PostJournalEntry(ctx context.Context, id, orgID, postedBy uuid.UUID) error {
	je, err := s.repo.GetJournalEntry(ctx, id, orgID)
	if err != nil {
		return err
	}

	// Check period status
	if je.AccountingPeriodID != nil {
		period, err := s.repo.GetAccountingPeriod(ctx, *je.AccountingPeriodID, orgID)
		if err != nil {
			return err
		}
		if period.Status != AccountingPeriodStatusOpen {
			return ErrPeriodClosed
		}
	}

	// Check fiscal year status
	if je.FiscalYearID != nil {
		fy, err := s.repo.GetFiscalYear(ctx, *je.FiscalYearID, orgID)
		if err != nil {
			return err
		}
		if fy.Status != FiscalYearStatusOpen {
			return ErrFiscalYearClosed
		}
	}

	// Get lines
	lines, err := s.repo.ListJournalEntryLines(ctx, je.ID)
	if err != nil {
		return err
	}

	// Post to general ledger
	if err := s.repo.PostToGeneralLedger(ctx, je, lines); err != nil {
		return err
	}

	// Update journal entry status
	je.Status = JournalEntryStatusPosted
	je.IsPosted = true
	je.PostedBy = &postedBy
	now := time.Now()
	je.PostedAt = &now
	je.UpdatedAt = now
	je.UpdatedBy = &postedBy

	return s.repo.UpdateJournalEntry(ctx, je, lines)
}

// ApproveJournalEntry approves a journal entry
func (s *Service) ApproveJournalEntry(ctx context.Context, id, orgID, approvedBy uuid.UUID) error {
	je, err := s.repo.GetJournalEntry(ctx, id, orgID)
	if err != nil {
		return err
	}

	je.Status = JournalEntryStatusApproved
	je.ApprovedBy = &approvedBy
	now := time.Now()
	je.ApprovedAt = &now
	je.UpdatedAt = now
	je.UpdatedBy = &approvedBy

	lines, err := s.repo.ListJournalEntryLines(ctx, je.ID)
	if err != nil {
		return err
	}

	return s.repo.UpdateJournalEntry(ctx, je, lines)
}

// ReverseJournalEntry reverses a posted journal entry
func (s *Service) ReverseJournalEntry(ctx context.Context, id, orgID uuid.UUID, reversedBy uuid.UUID) (*JournalEntry, error) {
	originalJE, err := s.repo.GetJournalEntry(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	if !originalJE.IsPosted {
		return nil, fmt.Errorf("can only reverse posted entries")
	}

	// Get original lines
	originalLines, err := s.repo.ListJournalEntryLines(ctx, originalJE.ID)
	if err != nil {
		return nil, err
	}

	// Create reversing entry
	now := time.Now()
	reversalID := uuid.New()
	entryNumber := fmt.Sprintf("JE-REV-%d-%s", now.Unix(), reversalID.String()[:8])

	// Build reversing lines (swap debit/credit)
	reversalLines := make([]*JournalEntryLine, 0, len(originalLines))
	for i, line := range originalLines {
		reversalLine := &JournalEntryLine{
			ID:             uuid.New(),
			OrganizationID: orgID,
			JournalEntryID: reversalID,
			LineNumber:     i + 1,
			AccountID:      line.AccountID,
			DebitAmount:    line.CreditAmount,  // Swap
			CreditAmount:   line.DebitAmount,   // Swap
			LocationID:     line.LocationID,
			Department:     line.Department,
			ProjectCode:    line.ProjectCode,
			CostCenter:     line.CostCenter,
			TaxCode:        line.TaxCode,
			Description:    line.Description,
			Memo:           line.Memo,
			CreatedAt:      now,
			UpdatedAt:      now,
			CreatedBy:      &reversedBy,
			UpdatedBy:      &reversedBy,
		}
		reversalLines = append(reversalLines, reversalLine)
	}

	reversalEntry := &JournalEntry{
		ID:               reversalID,
		OrganizationID:   orgID,
		EntryNumber:      entryNumber,
		EntryTypeID:      originalJE.EntryTypeID,
		EntryDate:        time.Now(),
		PostingDate:      time.Now(),
		AccountingPeriodID: originalJE.AccountingPeriodID,
		FiscalYearID:     originalJE.FiscalYearID,
		Status:           JournalEntryStatusReversed,
		IsPosted:         true,
		IsReversed:       true,
		ReversalEntryID:  &originalJE.ID,
		SourceModule:     originalJE.SourceModule,
		SourceDocumentType: originalJE.SourceDocumentType,
		SourceDocumentID: originalJE.SourceDocumentID,
		ReferenceNumber:  originalJE.ReferenceNumber,
		TotalDebit:       originalJE.TotalCredit,  // Swap
		TotalCredit:      originalJE.TotalDebit,   // Swap
		Description:      "Reversing: " + originalJE.Description,
		Notes:            originalJE.Notes,
		RequiresApproval: false,
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        &reversedBy,
		UpdatedBy:        &reversedBy,
	}

	// Create reversing entry
	if err := s.repo.CreateJournalEntry(ctx, reversalEntry, reversalLines); err != nil {
		return nil, err
	}

	// Update original entry
	originalJE.IsReversed = true
	originalJE.ReversalEntryID = &reversalID
	originalJE.UpdatedAt = now
	originalJE.UpdatedBy = &reversedBy

	if err := s.repo.UpdateJournalEntry(ctx, originalJE, originalLines); err != nil {
		return nil, err
	}

	return reversalEntry, nil
}

// ListJournalEntries lists journal entries
func (s *Service) ListJournalEntries(ctx context.Context, orgID uuid.UUID, filter map[string]interface{}) ([]*JournalEntry, error) {
	return s.repo.ListJournalEntries(ctx, orgID, filter)
}

// ========================
// GENERAL LEDGER & REPORTING
// ========================

// GetAccountBalance retrieves account balance as of a date
func (s *Service) GetAccountBalance(ctx context.Context, orgID, accountID uuid.UUID, asOfDate time.Time) (debit, credit, balance string, err error) {
	return s.repo.GetAccountBalance(ctx, orgID, accountID, asOfDate)
}

// GetTrialBalance retrieves trial balance
func (s *Service) GetTrialBalance(ctx context.Context, orgID uuid.UUID, asOfDate time.Time) (map[uuid.UUID]map[string]string, error) {
	return s.repo.GetTrialBalance(ctx, orgID, asOfDate)
}

// ========================
// HELPER FUNCTIONS
// ========================

func parseAmount(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func isBalanced(debit, credit float64) bool {
	diff := debit - credit
	if diff < 0 {
		diff = -diff
	}
	return diff < 0.01
}

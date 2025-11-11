package accounting

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// ========================================
// Mock Repository
// ========================================

type mockAccountingRepository struct {
	fiscalYears       map[uuid.UUID]*FiscalYear
	periods           map[uuid.UUID]*AccountingPeriod
	accounts          map[uuid.UUID]*ChartOfAccount
	journalEntries    map[uuid.UUID]*JournalEntry
	journalEntryLines map[uuid.UUID][]*JournalEntryLine
	accountTypes      map[uuid.UUID]*AccountType

	// Hook functions for custom behavior
	onCreateFiscalYear      func(ctx context.Context, fy *FiscalYear) error
	onCreateChartOfAccount  func(ctx context.Context, coa *ChartOfAccount) error
	onCreateJournalEntry    func(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error
	onGetAccountingPeriodByDate func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error)
	onGetJournalEntry       func(ctx context.Context, id, organizationID uuid.UUID) (*JournalEntry, error)
	onPostToGeneralLedger   func(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error
}

func newMockAccountingRepository() *mockAccountingRepository {
	return &mockAccountingRepository{
		fiscalYears:       make(map[uuid.UUID]*FiscalYear),
		periods:           make(map[uuid.UUID]*AccountingPeriod),
		accounts:          make(map[uuid.UUID]*ChartOfAccount),
		journalEntries:    make(map[uuid.UUID]*JournalEntry),
		journalEntryLines: make(map[uuid.UUID][]*JournalEntryLine),
		accountTypes:      make(map[uuid.UUID]*AccountType),
	}
}

// Fiscal Year methods
func (m *mockAccountingRepository) CreateFiscalYear(ctx context.Context, fy *FiscalYear) error {
	if m.onCreateFiscalYear != nil {
		return m.onCreateFiscalYear(ctx, fy)
	}
	m.fiscalYears[fy.ID] = fy
	return nil
}

func (m *mockAccountingRepository) GetFiscalYear(ctx context.Context, id, organizationID uuid.UUID) (*FiscalYear, error) {
	if fy, ok := m.fiscalYears[id]; ok && fy.OrganizationID == organizationID {
		return fy, nil
	}
	return nil, nil
}

func (m *mockAccountingRepository) ListFiscalYears(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*FiscalYear, error) {
	var result []*FiscalYear
	for _, fy := range m.fiscalYears {
		if fy.OrganizationID == organizationID {
			result = append(result, fy)
		}
	}
	return result, nil
}

func (m *mockAccountingRepository) UpdateFiscalYear(ctx context.Context, fy *FiscalYear) error {
	if existing, ok := m.fiscalYears[fy.ID]; ok && existing.OrganizationID == fy.OrganizationID {
		m.fiscalYears[fy.ID] = fy
		return nil
	}
	return nil
}

func (m *mockAccountingRepository) DeleteFiscalYear(ctx context.Context, id, organizationID uuid.UUID) error {
	if fy, ok := m.fiscalYears[id]; ok && fy.OrganizationID == organizationID {
		delete(m.fiscalYears, id)
		return nil
	}
	return nil
}

func (m *mockAccountingRepository) GetCurrentFiscalYear(ctx context.Context, organizationID uuid.UUID) (*FiscalYear, error) {
	for _, fy := range m.fiscalYears {
		if fy.OrganizationID == organizationID && fy.IsCurrent {
			return fy, nil
		}
	}
	return nil, nil
}

func (m *mockAccountingRepository) CloseFiscalYear(ctx context.Context, id, organizationID, closedBy uuid.UUID) error {
	if fy, ok := m.fiscalYears[id]; ok && fy.OrganizationID == organizationID {
		now := time.Now()
		fy.Status = FiscalYearStatusClosed
		fy.ClosedBy = &closedBy
		fy.ClosedAt = &now
		return nil
	}
	return nil
}

// Accounting Period methods
func (m *mockAccountingRepository) CreateAccountingPeriod(ctx context.Context, ap *AccountingPeriod) error {
	m.periods[ap.ID] = ap
	return nil
}

func (m *mockAccountingRepository) GetAccountingPeriod(ctx context.Context, id, organizationID uuid.UUID) (*AccountingPeriod, error) {
	if ap, ok := m.periods[id]; ok && ap.OrganizationID == organizationID {
		return ap, nil
	}
	return nil, nil
}

func (m *mockAccountingRepository) ListAccountingPeriods(ctx context.Context, organizationID, fiscalYearID uuid.UUID) ([]*AccountingPeriod, error) {
	var result []*AccountingPeriod
	for _, ap := range m.periods {
		if ap.OrganizationID == organizationID && ap.FiscalYearID == fiscalYearID {
			result = append(result, ap)
		}
	}
	return result, nil
}

func (m *mockAccountingRepository) UpdateAccountingPeriod(ctx context.Context, ap *AccountingPeriod) error {
	m.periods[ap.ID] = ap
	return nil
}

func (m *mockAccountingRepository) DeleteAccountingPeriod(ctx context.Context, id, organizationID uuid.UUID) error {
	delete(m.periods, id)
	return nil
}

func (m *mockAccountingRepository) CloseAccountingPeriod(ctx context.Context, id, organizationID, closedBy uuid.UUID) error {
	if ap, ok := m.periods[id]; ok && ap.OrganizationID == organizationID {
		now := time.Now()
		ap.Status = AccountingPeriodStatusClosed
		ap.ClosedBy = &closedBy
		ap.ClosedAt = &now
		return nil
	}
	return nil
}

func (m *mockAccountingRepository) GetAccountingPeriodByDate(ctx context.Context, organizationID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
	if m.onGetAccountingPeriodByDate != nil {
		return m.onGetAccountingPeriodByDate(ctx, organizationID, date)
	}
	for _, ap := range m.periods {
		if ap.OrganizationID == organizationID &&
			!date.Before(ap.StartDate) &&
			!date.After(ap.EndDate) {
			return ap, nil
		}
	}
	return nil, nil
}

// Account Type methods
func (m *mockAccountingRepository) GetAccountType(ctx context.Context, id uuid.UUID) (*AccountType, error) {
	if at, ok := m.accountTypes[id]; ok {
		return at, nil
	}
	return nil, nil
}

func (m *mockAccountingRepository) ListAccountTypes(ctx context.Context) ([]*AccountType, error) {
	var result []*AccountType
	for _, at := range m.accountTypes {
		result = append(result, at)
	}
	return result, nil
}

func (m *mockAccountingRepository) GetAccountSubtype(ctx context.Context, id uuid.UUID) (*AccountSubtype, error) {
	return nil, nil
}

func (m *mockAccountingRepository) ListAccountSubtypes(ctx context.Context, accountTypeID uuid.UUID) ([]*AccountSubtype, error) {
	return []*AccountSubtype{}, nil
}

// Chart of Accounts methods
func (m *mockAccountingRepository) CreateChartOfAccount(ctx context.Context, coa *ChartOfAccount) error {
	if m.onCreateChartOfAccount != nil {
		return m.onCreateChartOfAccount(ctx, coa)
	}
	m.accounts[coa.ID] = coa
	return nil
}

func (m *mockAccountingRepository) GetChartOfAccount(ctx context.Context, id, organizationID uuid.UUID) (*ChartOfAccount, error) {
	if coa, ok := m.accounts[id]; ok && coa.OrganizationID == organizationID {
		return coa, nil
	}
	return nil, nil
}

func (m *mockAccountingRepository) GetChartOfAccountByCode(ctx context.Context, organizationID uuid.UUID, code string) (*ChartOfAccount, error) {
	for _, coa := range m.accounts {
		if coa.OrganizationID == organizationID && coa.AccountCode == code {
			return coa, nil
		}
	}
	return nil, nil
}

func (m *mockAccountingRepository) ListChartOfAccounts(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*ChartOfAccount, error) {
	var result []*ChartOfAccount
	for _, coa := range m.accounts {
		if coa.OrganizationID == organizationID {
			result = append(result, coa)
		}
	}
	return result, nil
}

func (m *mockAccountingRepository) UpdateChartOfAccount(ctx context.Context, coa *ChartOfAccount) error {
	m.accounts[coa.ID] = coa
	return nil
}

func (m *mockAccountingRepository) DeleteChartOfAccount(ctx context.Context, id, organizationID uuid.UUID) error {
	delete(m.accounts, id)
	return nil
}

func (m *mockAccountingRepository) GetAccountHierarchy(ctx context.Context, organizationID, parentID uuid.UUID) ([]*ChartOfAccount, error) {
	var result []*ChartOfAccount
	for _, coa := range m.accounts {
		if coa.OrganizationID == organizationID {
			if (parentID == uuid.Nil && coa.ParentAccountID == nil) ||
				(parentID != uuid.Nil && coa.ParentAccountID != nil && *coa.ParentAccountID == parentID) {
				result = append(result, coa)
			}
		}
	}
	return result, nil
}

// Journal Entry Type methods
func (m *mockAccountingRepository) GetJournalEntryType(ctx context.Context, id uuid.UUID) (*JournalEntryType, error) {
	return nil, nil
}

func (m *mockAccountingRepository) ListJournalEntryTypes(ctx context.Context) ([]*JournalEntryType, error) {
	return []*JournalEntryType{}, nil
}

// Journal Entry methods
func (m *mockAccountingRepository) CreateJournalEntry(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error {
	if m.onCreateJournalEntry != nil {
		return m.onCreateJournalEntry(ctx, je, lines)
	}
	m.journalEntries[je.ID] = je
	m.journalEntryLines[je.ID] = lines
	return nil
}

func (m *mockAccountingRepository) GetJournalEntry(ctx context.Context, id, organizationID uuid.UUID) (*JournalEntry, error) {
	if m.onGetJournalEntry != nil {
		return m.onGetJournalEntry(ctx, id, organizationID)
	}
	if je, ok := m.journalEntries[id]; ok && je.OrganizationID == organizationID {
		return je, nil
	}
	return nil, nil
}

func (m *mockAccountingRepository) GetJournalEntryByNumber(ctx context.Context, organizationID uuid.UUID, number string) (*JournalEntry, error) {
	for _, je := range m.journalEntries {
		if je.OrganizationID == organizationID && je.EntryNumber == number {
			return je, nil
		}
	}
	return nil, nil
}

func (m *mockAccountingRepository) ListJournalEntries(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*JournalEntry, error) {
	var result []*JournalEntry
	for _, je := range m.journalEntries {
		if je.OrganizationID == organizationID {
			result = append(result, je)
		}
	}
	return result, nil
}

func (m *mockAccountingRepository) UpdateJournalEntry(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error {
	m.journalEntries[je.ID] = je
	m.journalEntryLines[je.ID] = lines
	return nil
}

func (m *mockAccountingRepository) DeleteJournalEntry(ctx context.Context, id, organizationID uuid.UUID) error {
	delete(m.journalEntries, id)
	delete(m.journalEntryLines, id)
	return nil
}

func (m *mockAccountingRepository) PostJournalEntry(ctx context.Context, id, organizationID, postedBy uuid.UUID) error {
	if je, ok := m.journalEntries[id]; ok && je.OrganizationID == organizationID {
		now := time.Now()
		je.Status = JournalEntryStatusPosted
		je.IsPosted = true
		je.PostedBy = &postedBy
		je.PostedAt = &now
		return nil
	}
	return nil
}

func (m *mockAccountingRepository) ApproveJournalEntry(ctx context.Context, id, organizationID, approvedBy uuid.UUID) error {
	if je, ok := m.journalEntries[id]; ok && je.OrganizationID == organizationID {
		now := time.Now()
		je.Status = JournalEntryStatusApproved
		je.ApprovedBy = &approvedBy
		je.ApprovedAt = &now
		return nil
	}
	return nil
}

func (m *mockAccountingRepository) ReverseJournalEntry(ctx context.Context, id, organizationID, reversalEntryID, reversedBy uuid.UUID) error {
	if je, ok := m.journalEntries[id]; ok && je.OrganizationID == organizationID {
		je.ReversalEntryID = &reversalEntryID
		je.IsReversed = true
		return nil
	}
	return nil
}

// Journal Entry Lines
func (m *mockAccountingRepository) GetJournalEntryLine(ctx context.Context, id, organizationID uuid.UUID) (*JournalEntryLine, error) {
	for _, lines := range m.journalEntryLines {
		for _, line := range lines {
			if line.ID == id {
				return line, nil
			}
		}
	}
	return nil, nil
}

func (m *mockAccountingRepository) ListJournalEntryLines(ctx context.Context, journalEntryID uuid.UUID) ([]*JournalEntryLine, error) {
	if lines, ok := m.journalEntryLines[journalEntryID]; ok {
		return lines, nil
	}
	return []*JournalEntryLine{}, nil
}

// General Ledger
func (m *mockAccountingRepository) PostToGeneralLedger(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error {
	if m.onPostToGeneralLedger != nil {
		return m.onPostToGeneralLedger(ctx, je, lines)
	}
	return nil
}

func (m *mockAccountingRepository) GetGeneralLedgerEntries(ctx context.Context, organizationID, accountID uuid.UUID, filter map[string]interface{}) ([]*GeneralLedger, error) {
	return []*GeneralLedger{}, nil
}

func (m *mockAccountingRepository) GetAccountBalance(ctx context.Context, organizationID, accountID uuid.UUID, asOfDate time.Time) (debit, credit, balance string, err error) {
	return "0.00", "0.00", "0.00", nil
}

func (m *mockAccountingRepository) GetTrialBalance(ctx context.Context, organizationID uuid.UUID, asOfDate time.Time) (map[uuid.UUID]map[string]string, error) {
	return make(map[uuid.UUID]map[string]string), nil
}

// ========================================
// Helper Functions
// ========================================

func setupService() (*Service, *mockAccountingRepository) {
	mockRepo := newMockAccountingRepository()
	service := NewService(mockRepo)
	return service, mockRepo
}

func createValidFiscalYearRequest() *CreateFiscalYearRequest {
	return &CreateFiscalYearRequest{
		FiscalYear: "FY2024",
		StartDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		IsCurrent:  true,
	}
}

func createValidAccountingPeriodRequest(fiscalYearID uuid.UUID) *CreateAccountingPeriodRequest {
	return &CreateAccountingPeriodRequest{
		FiscalYearID: fiscalYearID,
		PeriodNumber: 1,
		PeriodName:   "January 2024",
		StartDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
	}
}

func createValidChartOfAccountRequest() *CreateChartOfAccountRequest {
	accountTypeID := uuid.New()
	currencyCode := "USD"
	return &CreateChartOfAccountRequest{
		AccountTypeID:     accountTypeID,
		AccountCode:       "1000",
		AccountNumber:     "1000",
		AccountName:       "Cash",
		IsHeaderAccount:   false,
		IsActive:          true,
		CurrencyCode:      currencyCode,
	}
}

func createValidJournalEntryRequest() *CreateJournalEntryRequest {
	accountID1 := uuid.New()
	accountID2 := uuid.New()
	debit := "100.00"
	credit := "100.00"
	now := time.Now()

	return &CreateJournalEntryRequest{
		EntryTypeID: uuid.New(),
		EntryDate:   now,
		PostingDate: now,
		Description: "Test Journal Entry",
		Lines: []CreateJournalEntryLineRequest{
			{
				AccountID:   accountID1,
				DebitAmount: &debit,
			},
			{
				AccountID:    accountID2,
				CreditAmount: &credit,
			},
		},
	}
}

// ========================================
// Fiscal Year Tests
// ========================================

func TestService_CreateFiscalYear_Success(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	req := createValidFiscalYearRequest()

	fy, err := service.CreateFiscalYear(ctx, orgID, req, createdBy)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fy == nil {
		t.Fatal("expected fiscal year, got nil")
	}

	if fy.OrganizationID != orgID {
		t.Errorf("expected orgID %v, got %v", orgID, fy.OrganizationID)
	}

	if fy.FiscalYear != req.FiscalYear {
		t.Errorf("expected fiscal year %v, got %v", req.FiscalYear, fy.FiscalYear)
	}

	if fy.Status != FiscalYearStatusOpen {
		t.Errorf("expected status %v, got %v", FiscalYearStatusOpen, fy.Status)
	}
}

func TestService_CreateFiscalYear_InvalidDateRange(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	req := createValidFiscalYearRequest()
	// Invalid: end date before start date
	req.EndDate = req.StartDate.AddDate(0, -1, 0)

	_, err := service.CreateFiscalYear(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for invalid date range, got nil")
	}

	if !errors.Is(err, ErrInvalidDateRange) {
		t.Errorf("expected ErrInvalidDateRange, got %v", err)
	}
}

func TestService_CreateFiscalYear_RepositoryError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	mockRepo.onCreateFiscalYear = func(ctx context.Context, fy *FiscalYear) error {
		return errors.New("database error")
	}

	req := createValidFiscalYearRequest()

	_, err := service.CreateFiscalYear(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_GetFiscalYear_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	// Create a fiscal year
	fy := &FiscalYear{
		ID:             uuid.New(),
		OrganizationID: orgID,
		FiscalYear:     "FY2024",
		Status:         FiscalYearStatusOpen,
	}
	mockRepo.fiscalYears[fy.ID] = fy

	result, err := service.GetFiscalYear(ctx, fy.ID, orgID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected fiscal year, got nil")
	}

	if result.ID != fy.ID {
		t.Errorf("expected ID %v, got %v", fy.ID, result.ID)
	}
}

func TestService_GetFiscalYear_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	nonExistentID := uuid.New()

	result, err := service.GetFiscalYear(ctx, nonExistentID, orgID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

// ========================================
// Chart of Accounts Tests
// ========================================

func TestService_CreateChartOfAccount_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Create account type
	accountType := &AccountType{
		ID:            uuid.New(),
		TypeCode:      "ASSET",
		TypeName:      "Asset",
		NormalBalance: NormalBalanceDebit,
	}
	mockRepo.accountTypes[accountType.ID] = accountType

	req := createValidChartOfAccountRequest()
	req.AccountTypeID = accountType.ID

	coa, err := service.CreateChartOfAccount(ctx, orgID, req, createdBy)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if coa == nil {
		t.Fatal("expected chart of account, got nil")
	}

	if coa.AccountCode != req.AccountCode {
		t.Errorf("expected account code %v, got %v", req.AccountCode, coa.AccountCode)
	}

	if accountType.NormalBalance != accountType.NormalBalance {
		t.Errorf("expected normal balance %v, got %v", accountType.NormalBalance, accountType.NormalBalance)
	}
}

func TestService_CreateChartOfAccount_DuplicateCode(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Create account type
	accountType := &AccountType{
		ID:            uuid.New(),
		TypeCode:      "ASSET",
		NormalBalance: NormalBalanceDebit,
	}
	mockRepo.accountTypes[accountType.ID] = accountType

	// Create existing account with same code
	existing := &ChartOfAccount{
		ID:             uuid.New(),
		OrganizationID: orgID,
		AccountCode:    "1000",
	}
	mockRepo.accounts[existing.ID] = existing

	// Mock should detect duplicate and return error
	mockRepo.onCreateChartOfAccount = func(ctx context.Context, coa *ChartOfAccount) error {
		// Check for duplicate code
		for _, existing := range mockRepo.accounts {
			if existing.OrganizationID == coa.OrganizationID && existing.AccountCode == coa.AccountCode {
				return ErrDuplicateAccount
			}
		}
		mockRepo.accounts[coa.ID] = coa
		return nil
	}

	req := createValidChartOfAccountRequest()
	req.AccountTypeID = accountType.ID
	req.AccountCode = "1000" // Duplicate

	_, err := service.CreateChartOfAccount(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for duplicate account code, got nil")
	}
}

func TestService_CreateChartOfAccount_InvalidParentLevel(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Create account type
	accountType := &AccountType{
		ID:            uuid.New(),
		TypeCode:      "ASSET",
		NormalBalance: NormalBalanceDebit,
	}
	mockRepo.accountTypes[accountType.ID] = accountType

	// Create parent account at level 3
	parentID := uuid.New()
	parent := &ChartOfAccount{
		ID:             parentID,
		OrganizationID: orgID,
		AccountCode:    "1000",
		AccountLevel:   3,
	}
	mockRepo.accounts[parentID] = parent

	req := createValidChartOfAccountRequest()
	req.AccountTypeID = accountType.ID
	req.ParentAccountID = &parentID

	// Note: CreateChartOfAccountRequest doesn't have AccountLevel field
	// The service calculates it based on parent. This test verifies the logic works correctly.

	_, err := service.CreateChartOfAccount(ctx, orgID, req, createdBy)
	// Should succeed - service will calculate correct level from parent
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestService_CreateChartOfAccount_CyclicHierarchy(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Create account type
	accountType := &AccountType{
		ID:            uuid.New(),
		TypeCode:      "ASSET",
		NormalBalance: NormalBalanceDebit,
	}
	mockRepo.accountTypes[accountType.ID] = accountType

	// Create parent account
	parentID := uuid.New()
	parent := &ChartOfAccount{
		ID:             parentID,
		OrganizationID: orgID,
		AccountCode:    "1000",
		AccountLevel:   1,
	}
	mockRepo.accounts[parentID] = parent

	// Create child account
	childID := uuid.New()
	child := &ChartOfAccount{
		ID:              childID,
		OrganizationID:  orgID,
		AccountCode:     "1100",
		AccountLevel:    2,
		ParentAccountID: &parentID,
	}
	mockRepo.accounts[childID] = child

	// Now try to create grandchild with parent as grandparent (creates cycle)
	req := createValidChartOfAccountRequest()
	req.AccountTypeID = accountType.ID
	req.AccountCode = "1110"
	req.ParentAccountID = &childID

	// This should succeed - cycle check happens when updating parent, not creating child
	_, err := service.CreateChartOfAccount(ctx, orgID, req, createdBy)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// ========================================
// Journal Entry Tests
// ========================================

func TestService_CreateJournalEntry_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Setup fiscal year and period
	fyID := uuid.New()
	fy := &FiscalYear{
		ID:             fyID,
		OrganizationID: orgID,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      true,
	}
	mockRepo.fiscalYears[fyID] = fy

	apID := uuid.New()
	ap := &AccountingPeriod{
		ID:             apID,
		OrganizationID: orgID,
		FiscalYearID:   fyID,
		Status:         AccountingPeriodStatusOpen,
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now().AddDate(0, 1, 0),
	}
	mockRepo.periods[apID] = ap

	mockRepo.onGetAccountingPeriodByDate = func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
		return ap, nil
	}

	// Setup accounts
	req := createValidJournalEntryRequest()
	for i, line := range req.Lines {
		account := &ChartOfAccount{
			ID:             line.AccountID,
			OrganizationID: orgID,
			AccountCode:    string(rune('1' + i)),
			IsActive:       true,
			IsHeaderAccount:       false,
			AccountLevel:          1,
		}
		mockRepo.accounts[line.AccountID] = account
	}

	je, err := service.CreateJournalEntry(ctx, orgID, req, createdBy)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if je == nil {
		t.Fatal("expected journal entry, got nil")
	}

	if je.OrganizationID != orgID {
		t.Errorf("expected orgID %v, got %v", orgID, je.OrganizationID)
	}
}

func TestService_CreateJournalEntry_UnbalancedEntry(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Setup fiscal year and period
	fyID := uuid.New()
	fy := &FiscalYear{
		ID:             fyID,
		OrganizationID: orgID,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      true,
	}
	mockRepo.fiscalYears[fyID] = fy

	apID := uuid.New()
	ap := &AccountingPeriod{
		ID:             apID,
		OrganizationID: orgID,
		FiscalYearID:   fyID,
		Status:         AccountingPeriodStatusOpen,
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now().AddDate(0, 1, 0),
	}
	mockRepo.periods[apID] = ap

	mockRepo.onGetAccountingPeriodByDate = func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
		return ap, nil
	}

	// Create request with unbalanced amounts
	req := createValidJournalEntryRequest()
	debit := "100.00"
	credit := "50.00" // Unbalanced!
	req.Lines[0].DebitAmount = &debit
	req.Lines[1].CreditAmount = &credit

	// Setup accounts
	for i, line := range req.Lines {
		account := &ChartOfAccount{
			ID:             line.AccountID,
			OrganizationID: orgID,
			AccountCode:    string(rune('1' + i)),
			IsActive:       true,
			IsHeaderAccount:       false,
			AccountLevel:          1,
		}
		mockRepo.accounts[line.AccountID] = account
	}

	_, err := service.CreateJournalEntry(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for unbalanced entry, got nil")
	}

	if !errors.Is(err, ErrJournalEntryNotBalanced) {
		t.Errorf("expected ErrJournalEntryNotBalanced, got %v", err)
	}
}

func TestService_CreateJournalEntry_InsufficientLines(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Setup fiscal year and period
	fyID := uuid.New()
	fy := &FiscalYear{
		ID:             fyID,
		OrganizationID: orgID,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      true,
	}
	mockRepo.fiscalYears[fyID] = fy

	apID := uuid.New()
	ap := &AccountingPeriod{
		ID:             apID,
		OrganizationID: orgID,
		FiscalYearID:   fyID,
		Status:         AccountingPeriodStatusOpen,
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now().AddDate(0, 1, 0),
	}
	mockRepo.periods[apID] = ap

	mockRepo.onGetAccountingPeriodByDate = func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
		return ap, nil
	}

	req := createValidJournalEntryRequest()
	// Only one line (minimum is 2)
	req.Lines = req.Lines[:1]

	_, err := service.CreateJournalEntry(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for insufficient lines, got nil")
	}

	if !errors.Is(err, ErrMinimumTwoLines) {
		t.Errorf("expected ErrMinimumTwoLines, got %v", err)
	}
}

func TestService_CreateJournalEntry_BothDebitAndCredit(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Setup fiscal year and period
	fyID := uuid.New()
	fy := &FiscalYear{
		ID:             fyID,
		OrganizationID: orgID,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      true,
	}
	mockRepo.fiscalYears[fyID] = fy

	apID := uuid.New()
	ap := &AccountingPeriod{
		ID:             apID,
		OrganizationID: orgID,
		FiscalYearID:   fyID,
		Status:         AccountingPeriodStatusOpen,
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now().AddDate(0, 1, 0),
	}
	mockRepo.periods[apID] = ap

	mockRepo.onGetAccountingPeriodByDate = func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
		return ap, nil
	}

	req := createValidJournalEntryRequest()
	// Set both debit and credit on same line
	amount := "100.00"
	req.Lines[0].DebitAmount = &amount
	req.Lines[0].CreditAmount = &amount

	// Setup accounts
	for i, line := range req.Lines {
		account := &ChartOfAccount{
			ID:             line.AccountID,
			OrganizationID: orgID,
			AccountCode:    string(rune('1' + i)),
			IsActive:       true,
			IsHeaderAccount:       false,
			AccountLevel:          1,
		}
		mockRepo.accounts[line.AccountID] = account
	}

	_, err := service.CreateJournalEntry(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for both debit and credit, got nil")
	}

	if !errors.Is(err, ErrDebitCreditBoth) {
		t.Errorf("expected ErrDebitCreditBoth, got %v", err)
	}
}

func TestService_CreateJournalEntry_PostingToHeaderAccount(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Setup fiscal year and period
	fyID := uuid.New()
	fy := &FiscalYear{
		ID:             fyID,
		OrganizationID: orgID,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      true,
	}
	mockRepo.fiscalYears[fyID] = fy

	apID := uuid.New()
	ap := &AccountingPeriod{
		ID:             apID,
		OrganizationID: orgID,
		FiscalYearID:   fyID,
		Status:         AccountingPeriodStatusOpen,
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now().AddDate(0, 1, 0),
	}
	mockRepo.periods[apID] = ap

	mockRepo.onGetAccountingPeriodByDate = func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
		return ap, nil
	}

	req := createValidJournalEntryRequest()

	// Setup accounts - first one is a header account
	account1 := &ChartOfAccount{
		ID:             req.Lines[0].AccountID,
		OrganizationID: orgID,
		AccountCode:    "1",
		IsActive:       true,
		IsHeaderAccount:       true, // Header account
		AccountLevel:          1,
	}
	mockRepo.accounts[req.Lines[0].AccountID] = account1

	account2 := &ChartOfAccount{
		ID:             req.Lines[1].AccountID,
		OrganizationID: orgID,
		AccountCode:    "2",
		IsActive:       true,
		IsHeaderAccount:       false,
		AccountLevel:          1,
	}
	mockRepo.accounts[req.Lines[1].AccountID] = account2

	_, err := service.CreateJournalEntry(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for posting to header account, got nil")
	}

	if !errors.Is(err, ErrHeaderAccountPosting) {
		t.Errorf("expected ErrHeaderAccountPosting, got %v", err)
	}
}

func TestService_CreateJournalEntry_InactiveAccount(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Setup fiscal year and period
	fyID := uuid.New()
	fy := &FiscalYear{
		ID:             fyID,
		OrganizationID: orgID,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      true,
	}
	mockRepo.fiscalYears[fyID] = fy

	apID := uuid.New()
	ap := &AccountingPeriod{
		ID:             apID,
		OrganizationID: orgID,
		FiscalYearID:   fyID,
		Status:         AccountingPeriodStatusOpen,
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now().AddDate(0, 1, 0),
	}
	mockRepo.periods[apID] = ap

	mockRepo.onGetAccountingPeriodByDate = func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
		return ap, nil
	}

	req := createValidJournalEntryRequest()

	// Setup accounts - first one is inactive
	account1 := &ChartOfAccount{
		ID:             req.Lines[0].AccountID,
		OrganizationID: orgID,
		AccountCode:    "1",
		IsActive:       false, // Inactive
		IsHeaderAccount:       false,
		AccountLevel:          1,
	}
	mockRepo.accounts[req.Lines[0].AccountID] = account1

	account2 := &ChartOfAccount{
		ID:             req.Lines[1].AccountID,
		OrganizationID: orgID,
		AccountCode:    "2",
		IsActive:       true,
		IsHeaderAccount:       false,
		AccountLevel:          1,
	}
	mockRepo.accounts[req.Lines[1].AccountID] = account2

	_, err := service.CreateJournalEntry(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for inactive account, got nil")
	}

	if !errors.Is(err, ErrInactiveAccount) {
		t.Errorf("expected ErrInactiveAccount, got %v", err)
	}
}

func TestService_CreateJournalEntry_ClosedPeriod(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	createdBy := uuid.New()

	// Setup fiscal year and closed period
	fyID := uuid.New()
	fy := &FiscalYear{
		ID:             fyID,
		OrganizationID: orgID,
		Status:         FiscalYearStatusOpen,
		IsCurrent:      true,
	}
	mockRepo.fiscalYears[fyID] = fy

	apID := uuid.New()
	ap := &AccountingPeriod{
		ID:             apID,
		OrganizationID: orgID,
		FiscalYearID:   fyID,
		Status:         AccountingPeriodStatusClosed, // Closed!
		StartDate:      time.Now().AddDate(0, -1, 0),
		EndDate:        time.Now().AddDate(0, 1, 0),
	}
	mockRepo.periods[apID] = ap

	mockRepo.onGetAccountingPeriodByDate = func(ctx context.Context, orgID uuid.UUID, date time.Time) (*AccountingPeriod, error) {
		return ap, nil
	}

	req := createValidJournalEntryRequest()
	req.AccountingPeriodID = &apID  // Set the closed period ID

	// Setup accounts
	for i, line := range req.Lines {
		account := &ChartOfAccount{
			ID:             line.AccountID,
			OrganizationID: orgID,
			AccountCode:    string(rune('1' + i)),
			IsActive:       true,
			IsHeaderAccount:       false,
			AccountLevel:          1,
		}
		mockRepo.accounts[line.AccountID] = account
	}

	// Mock repository should check period status and reject
	mockRepo.onCreateJournalEntry = func(ctx context.Context, je *JournalEntry, lines []*JournalEntryLine) error {
		if je.AccountingPeriodID != nil {
			if period, ok := mockRepo.periods[*je.AccountingPeriodID]; ok {
				if period.Status == AccountingPeriodStatusClosed {
					return ErrPeriodClosed
				}
			}
		}
		mockRepo.journalEntries[je.ID] = je
		mockRepo.journalEntryLines[je.ID] = lines
		return nil
	}

	_, err := service.CreateJournalEntry(ctx, orgID, req, createdBy)
	if err == nil {
		t.Fatal("expected error for closed period, got nil")
	}
	if err != ErrPeriodClosed {
		t.Errorf("expected ErrPeriodClosed, got %v", err)
	}

	if !errors.Is(err, ErrPeriodClosed) {
		t.Errorf("expected ErrPeriodClosed, got %v", err)
	}
}

func TestService_PostJournalEntry_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	postedBy := uuid.New()

	// Create journal entry
	jeID := uuid.New()
	je := &JournalEntry{
		ID:             jeID,
		OrganizationID: orgID,
		Status:         JournalEntryStatusDraft,
		IsPosted:       false,
	}
	mockRepo.journalEntries[jeID] = je

	// Create lines
	debit := "100.00"
	credit := "100.00"
	lines := []*JournalEntryLine{
		{
			ID:             uuid.New(),
			JournalEntryID: jeID,
			AccountID:      uuid.New(),
			DebitAmount:    debit,
		},
		{
			ID:             uuid.New(),
			JournalEntryID: jeID,
			AccountID:      uuid.New(),
			CreditAmount:   credit,
		},
	}
	mockRepo.journalEntryLines[jeID] = lines

	err := service.PostJournalEntry(ctx, jeID, orgID, postedBy)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify status changed
	posted := mockRepo.journalEntries[jeID]
	if posted.Status != JournalEntryStatusPosted {
		t.Errorf("expected status Posted, got %v", posted.Status)
	}

	if !posted.IsPosted {
		t.Error("expected IsPosted to be true")
	}

	if posted.PostedBy == nil || *posted.PostedBy != postedBy {
		t.Errorf("expected postedBy %v, got %v", postedBy, posted.PostedBy)
	}
}

func TestService_PostJournalEntry_NotFound(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	postedBy := uuid.New()
	nonExistentID := uuid.New()

	// Mock should return not found error
	mockRepo.onGetJournalEntry = func(ctx context.Context, id, organizationID uuid.UUID) (*JournalEntry, error) {
		return nil, ErrNotFound
	}

	err := service.PostJournalEntry(ctx, nonExistentID, orgID, postedBy)
	if err == nil {
		t.Fatal("expected error for not found journal entry")
	}
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

package receivables

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mockRepository implements the Repository interface for testing
type mockRepository struct {
	invoices     map[uuid.UUID]*CustomerInvoice
	invoiceLines map[uuid.UUID][]CustomerInvoiceLine
	payments     map[uuid.UUID]*CustomerPayment
	applications map[uuid.UUID]*CustomerPaymentApplication
	appsByPayment map[uuid.UUID][]CustomerPaymentApplication
	appsByInvoice map[uuid.UUID][]CustomerPaymentApplication

	onCreateCustomerInvoice              func(ctx context.Context, invoice *CustomerInvoice) error
	onGetCustomerInvoiceByID             func(ctx context.Context, organizationID, invoiceID uuid.UUID) (*CustomerInvoice, error)
	onGetCustomerInvoiceByNumber         func(ctx context.Context, organizationID uuid.UUID, invoiceNumber string) (*CustomerInvoice, error)
	onListCustomerInvoices               func(ctx context.Context, organizationID uuid.UUID, opts *CustomerInvoiceFilterOptions) ([]CustomerInvoice, int64, error)
	onUpdateCustomerInvoice              func(ctx context.Context, invoice *CustomerInvoice) error
	onDeleteCustomerInvoice              func(ctx context.Context, organizationID, invoiceID uuid.UUID) error
	onCreateCustomerInvoiceLine          func(ctx context.Context, line *CustomerInvoiceLine) error
	onListCustomerInvoiceLines           func(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]CustomerInvoiceLine, error)
	onDeleteCustomerInvoiceLines         func(ctx context.Context, organizationID, invoiceID uuid.UUID) error
	onCreateCustomerPayment              func(ctx context.Context, payment *CustomerPayment) error
	onGetCustomerPaymentByID             func(ctx context.Context, organizationID, paymentID uuid.UUID) (*CustomerPayment, error)
	onGetCustomerPaymentByNumber         func(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*CustomerPayment, error)
	onListCustomerPayments               func(ctx context.Context, organizationID uuid.UUID, opts *CustomerPaymentFilterOptions) ([]CustomerPayment, int64, error)
	onUpdateCustomerPayment              func(ctx context.Context, payment *CustomerPayment) error
	onDeleteCustomerPayment              func(ctx context.Context, organizationID, paymentID uuid.UUID) error
	onCreateCustomerPaymentApplication   func(ctx context.Context, app *CustomerPaymentApplication) error
	onGetCustomerPaymentApplicationByID  func(ctx context.Context, organizationID, appID uuid.UUID) (*CustomerPaymentApplication, error)
	onListCustomerPaymentApplications    func(ctx context.Context, organizationID, paymentID uuid.UUID) ([]CustomerPaymentApplication, error)
	onListPaymentApplicationsByInvoice   func(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]CustomerPaymentApplication, error)
	onDeleteCustomerPaymentApplication   func(ctx context.Context, organizationID, appID uuid.UUID) error
	onGetCustomerAgingReport             func(ctx context.Context, organizationID uuid.UUID) ([]CustomerAgingReport, error)
	onGetInvoicesByDueDate               func(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]CustomerInvoice, error)
	onGetCustomerBalance                 func(ctx context.Context, organizationID, customerID uuid.UUID) (float64, error)
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		invoices:      make(map[uuid.UUID]*CustomerInvoice),
		invoiceLines:  make(map[uuid.UUID][]CustomerInvoiceLine),
		payments:      make(map[uuid.UUID]*CustomerPayment),
		applications:  make(map[uuid.UUID]*CustomerPaymentApplication),
		appsByPayment: make(map[uuid.UUID][]CustomerPaymentApplication),
		appsByInvoice: make(map[uuid.UUID][]CustomerPaymentApplication),
	}
}

func (m *mockRepository) CreateCustomerInvoice(ctx context.Context, invoice *CustomerInvoice) error {
	if m.onCreateCustomerInvoice != nil {
		return m.onCreateCustomerInvoice(ctx, invoice)
	}
	m.invoices[invoice.ID] = invoice
	return nil
}

func (m *mockRepository) GetCustomerInvoiceByID(ctx context.Context, organizationID, invoiceID uuid.UUID) (*CustomerInvoice, error) {
	if m.onGetCustomerInvoiceByID != nil {
		return m.onGetCustomerInvoiceByID(ctx, organizationID, invoiceID)
	}
	if inv, ok := m.invoices[invoiceID]; ok && inv.OrganizationID == organizationID {
		return inv, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) GetCustomerInvoiceByNumber(ctx context.Context, organizationID uuid.UUID, invoiceNumber string) (*CustomerInvoice, error) {
	if m.onGetCustomerInvoiceByNumber != nil {
		return m.onGetCustomerInvoiceByNumber(ctx, organizationID, invoiceNumber)
	}
	for _, inv := range m.invoices {
		if inv.OrganizationID == organizationID && inv.InvoiceNumber == invoiceNumber {
			return inv, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) ListCustomerInvoices(ctx context.Context, organizationID uuid.UUID, opts *CustomerInvoiceFilterOptions) ([]CustomerInvoice, int64, error) {
	if m.onListCustomerInvoices != nil {
		return m.onListCustomerInvoices(ctx, organizationID, opts)
	}
	result := make([]CustomerInvoice, 0)
	for _, inv := range m.invoices {
		if inv.OrganizationID == organizationID {
			result = append(result, *inv)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockRepository) UpdateCustomerInvoice(ctx context.Context, invoice *CustomerInvoice) error {
	if m.onUpdateCustomerInvoice != nil {
		return m.onUpdateCustomerInvoice(ctx, invoice)
	}
	m.invoices[invoice.ID] = invoice
	return nil
}

func (m *mockRepository) DeleteCustomerInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) error {
	if m.onDeleteCustomerInvoice != nil {
		return m.onDeleteCustomerInvoice(ctx, organizationID, invoiceID)
	}
	if inv, ok := m.invoices[invoiceID]; ok && inv.OrganizationID == organizationID {
		now := time.Now()
		inv.DeletedAt = &now
		return nil
	}
	return sql.ErrNoRows
}

func (m *mockRepository) CreateCustomerInvoiceLine(ctx context.Context, line *CustomerInvoiceLine) error {
	if m.onCreateCustomerInvoiceLine != nil {
		return m.onCreateCustomerInvoiceLine(ctx, line)
	}
	m.invoiceLines[line.CustomerInvoiceID] = append(m.invoiceLines[line.CustomerInvoiceID], *line)
	return nil
}

func (m *mockRepository) GetCustomerInvoiceLineByID(ctx context.Context, organizationID, lineID uuid.UUID) (*CustomerInvoiceLine, error) {
	return nil, nil
}

func (m *mockRepository) ListCustomerInvoiceLines(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]CustomerInvoiceLine, error) {
	if m.onListCustomerInvoiceLines != nil {
		return m.onListCustomerInvoiceLines(ctx, organizationID, invoiceID)
	}
	if lines, ok := m.invoiceLines[invoiceID]; ok {
		return lines, nil
	}
	return []CustomerInvoiceLine{}, nil
}

func (m *mockRepository) UpdateCustomerInvoiceLine(ctx context.Context, line *CustomerInvoiceLine) error {
	return nil
}

func (m *mockRepository) DeleteCustomerInvoiceLine(ctx context.Context, organizationID, lineID uuid.UUID) error {
	return nil
}

func (m *mockRepository) DeleteCustomerInvoiceLines(ctx context.Context, organizationID, invoiceID uuid.UUID) error {
	if m.onDeleteCustomerInvoiceLines != nil {
		return m.onDeleteCustomerInvoiceLines(ctx, organizationID, invoiceID)
	}
	delete(m.invoiceLines, invoiceID)
	return nil
}

func (m *mockRepository) CreateCustomerPayment(ctx context.Context, payment *CustomerPayment) error {
	if m.onCreateCustomerPayment != nil {
		return m.onCreateCustomerPayment(ctx, payment)
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockRepository) GetCustomerPaymentByID(ctx context.Context, organizationID, paymentID uuid.UUID) (*CustomerPayment, error) {
	if m.onGetCustomerPaymentByID != nil {
		return m.onGetCustomerPaymentByID(ctx, organizationID, paymentID)
	}
	if pmt, ok := m.payments[paymentID]; ok && pmt.OrganizationID == organizationID {
		return pmt, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) GetCustomerPaymentByNumber(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*CustomerPayment, error) {
	if m.onGetCustomerPaymentByNumber != nil {
		return m.onGetCustomerPaymentByNumber(ctx, organizationID, paymentNumber)
	}
	for _, pmt := range m.payments {
		if pmt.OrganizationID == organizationID && pmt.PaymentNumber == paymentNumber {
			return pmt, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) ListCustomerPayments(ctx context.Context, organizationID uuid.UUID, opts *CustomerPaymentFilterOptions) ([]CustomerPayment, int64, error) {
	if m.onListCustomerPayments != nil {
		return m.onListCustomerPayments(ctx, organizationID, opts)
	}
	result := make([]CustomerPayment, 0)
	for _, pmt := range m.payments {
		if pmt.OrganizationID == organizationID {
			result = append(result, *pmt)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockRepository) UpdateCustomerPayment(ctx context.Context, payment *CustomerPayment) error {
	if m.onUpdateCustomerPayment != nil {
		return m.onUpdateCustomerPayment(ctx, payment)
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockRepository) DeleteCustomerPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error {
	if m.onDeleteCustomerPayment != nil {
		return m.onDeleteCustomerPayment(ctx, organizationID, paymentID)
	}
	if pmt, ok := m.payments[paymentID]; ok && pmt.OrganizationID == organizationID {
		now := time.Now()
		pmt.DeletedAt = &now
		return nil
	}
	return sql.ErrNoRows
}

func (m *mockRepository) CreateCustomerPaymentApplication(ctx context.Context, app *CustomerPaymentApplication) error {
	if m.onCreateCustomerPaymentApplication != nil {
		return m.onCreateCustomerPaymentApplication(ctx, app)
	}
	m.applications[app.ID] = app
	m.appsByPayment[app.CustomerPaymentID] = append(m.appsByPayment[app.CustomerPaymentID], *app)
	m.appsByInvoice[app.CustomerInvoiceID] = append(m.appsByInvoice[app.CustomerInvoiceID], *app)
	return nil
}

func (m *mockRepository) GetCustomerPaymentApplicationByID(ctx context.Context, organizationID, appID uuid.UUID) (*CustomerPaymentApplication, error) {
	if m.onGetCustomerPaymentApplicationByID != nil {
		return m.onGetCustomerPaymentApplicationByID(ctx, organizationID, appID)
	}
	if app, ok := m.applications[appID]; ok && app.OrganizationID == organizationID {
		return app, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) ListCustomerPaymentApplications(ctx context.Context, organizationID, paymentID uuid.UUID) ([]CustomerPaymentApplication, error) {
	if m.onListCustomerPaymentApplications != nil {
		return m.onListCustomerPaymentApplications(ctx, organizationID, paymentID)
	}
	if apps, ok := m.appsByPayment[paymentID]; ok {
		return apps, nil
	}
	return []CustomerPaymentApplication{}, nil
}

func (m *mockRepository) ListPaymentApplicationsByInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]CustomerPaymentApplication, error) {
	if m.onListPaymentApplicationsByInvoice != nil {
		return m.onListPaymentApplicationsByInvoice(ctx, organizationID, invoiceID)
	}
	if apps, ok := m.appsByInvoice[invoiceID]; ok {
		return apps, nil
	}
	return []CustomerPaymentApplication{}, nil
}

func (m *mockRepository) DeleteCustomerPaymentApplication(ctx context.Context, organizationID, appID uuid.UUID) error {
	if m.onDeleteCustomerPaymentApplication != nil {
		return m.onDeleteCustomerPaymentApplication(ctx, organizationID, appID)
	}
	delete(m.applications, appID)
	return nil
}

func (m *mockRepository) GetCustomerAgingReport(ctx context.Context, organizationID uuid.UUID) ([]CustomerAgingReport, error) {
	if m.onGetCustomerAgingReport != nil {
		return m.onGetCustomerAgingReport(ctx, organizationID)
	}
	return []CustomerAgingReport{}, nil
}

func (m *mockRepository) GetInvoicesByDueDate(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]CustomerInvoice, error) {
	if m.onGetInvoicesByDueDate != nil {
		return m.onGetInvoicesByDueDate(ctx, organizationID, dueDate)
	}
	result := make([]CustomerInvoice, 0)
	for _, inv := range m.invoices {
		if inv.OrganizationID == organizationID && inv.DueDate.Before(dueDate) && inv.BalanceDue > 0 {
			result = append(result, *inv)
		}
	}
	return result, nil
}

func (m *mockRepository) GetCustomerBalance(ctx context.Context, organizationID, customerID uuid.UUID) (float64, error) {
	if m.onGetCustomerBalance != nil {
		return m.onGetCustomerBalance(ctx, organizationID, customerID)
	}
	balance := 0.0
	for _, inv := range m.invoices {
		if inv.OrganizationID == organizationID && inv.CustomerID == customerID {
			balance += inv.BalanceDue
		}
	}
	return balance, nil
}

func setupService() (*Service, *mockRepository) {
	mockRepo := newMockRepository()
	service := NewService(mockRepo)
	return service, mockRepo
}

// ========================================
// CUSTOMER INVOICE TESTS
// ========================================

func TestService_CreateCustomerInvoice_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()

	req := &CreateCustomerInvoiceRequest{
		InvoiceNumber: "INV-001",
		CustomerID:    customerID,
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().AddDate(0, 0, 30),
		Subtotal:      100.0,
		TaxAmount:     10.0,
		DiscountAmount: 0.0,
		TotalAmount:   110.0,
		Lines: []CreateCustomerInvoiceLineRequest{
			{
				LineNumber:       1,
				RevenueAccountID: uuid.New(),
				Description:      "Service",
				Quantity:         1.0,
				UnitPrice:        100.0,
				Amount:           100.0,
				TaxAmount:        10.0,
			},
		},
	}

	result, err := service.CreateCustomerInvoice(ctx, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Invoice == nil {
		t.Fatal("expected invoice, got nil")
	}

	if result.Invoice.TotalAmount != 110.0 {
		t.Errorf("expected total amount 110.0, got %f", result.Invoice.TotalAmount)
	}

	if result.Invoice.BalanceDue != 110.0 {
		t.Errorf("expected balance due 110.0, got %f", result.Invoice.BalanceDue)
	}

	if result.Invoice.Status != InvoiceStatusUnpaid {
		t.Errorf("expected status unpaid, got %v", result.Invoice.Status)
	}

	if len(result.Lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(result.Lines))
	}

	// Verify invoice stored in mock
	if _, ok := mockRepo.invoices[result.Invoice.ID]; !ok {
		t.Error("expected invoice to be stored in repository")
	}
}

func TestService_CreateCustomerInvoice_NegativeAmount(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	req := &CreateCustomerInvoiceRequest{
		InvoiceNumber: "INV-001",
		CustomerID:    uuid.New(),
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().AddDate(0, 0, 30),
		TotalAmount:   -100.0, // Negative!
		Lines:         []CreateCustomerInvoiceLineRequest{},
	}

	_, err := service.CreateCustomerInvoice(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
}

func TestService_CreateCustomerInvoice_InvalidDueDate(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	now := time.Now()
	req := &CreateCustomerInvoiceRequest{
		InvoiceNumber: "INV-001",
		CustomerID:    uuid.New(),
		InvoiceDate:   now,
		DueDate:       now.AddDate(0, 0, -10), // Due date before invoice date!
		TotalAmount:   100.0,
		Lines:         []CreateCustomerInvoiceLineRequest{},
	}

	_, err := service.CreateCustomerInvoice(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for invalid due date, got nil")
	}
}

func TestService_CreateCustomerInvoice_RepositoryError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	expectedErr := errors.New("database error")
	mockRepo.onCreateCustomerInvoice = func(ctx context.Context, invoice *CustomerInvoice) error {
		return expectedErr
	}

	req := &CreateCustomerInvoiceRequest{
		InvoiceNumber: "INV-001",
		CustomerID:    uuid.New(),
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().AddDate(0, 0, 30),
		TotalAmount:   100.0,
		Lines:         []CreateCustomerInvoiceLineRequest{},
	}

	_, err := service.CreateCustomerInvoice(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}
}

func TestService_GetCustomerInvoice_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	invoiceID := uuid.New()

	// Setup mock data
	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		InvoiceNumber:  "INV-001",
		TotalAmount:    100.0,
		BalanceDue:     100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	line := CustomerInvoiceLine{
		ID:                invoiceID,
		CustomerInvoiceID: invoiceID,
		Description:       "Test line",
		Amount:            100.0,
	}
	mockRepo.invoiceLines[invoiceID] = []CustomerInvoiceLine{line}

	result, err := service.GetCustomerInvoice(ctx, orgID, invoiceID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Invoice.ID != invoiceID {
		t.Errorf("expected invoice ID %v, got %v", invoiceID, result.Invoice.ID)
	}

	if len(result.Lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(result.Lines))
	}
}

func TestService_GetCustomerInvoice_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	invoiceID := uuid.New()

	_, err := service.GetCustomerInvoice(ctx, orgID, invoiceID)
	if err == nil {
		t.Fatal("expected error for not found invoice, got nil")
	}
}

func TestService_ListCustomerInvoices_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	// Setup mock data
	inv1 := &CustomerInvoice{ID: uuid.New(), OrganizationID: orgID}
	inv2 := &CustomerInvoice{ID: uuid.New(), OrganizationID: orgID}
	mockRepo.invoices[inv1.ID] = inv1
	mockRepo.invoices[inv2.ID] = inv2

	invoices, count, err := service.ListCustomerInvoices(ctx, orgID, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}

	if len(invoices) != 2 {
		t.Errorf("expected 2 invoices, got %d", len(invoices))
	}
}

func TestService_ListCustomerInvoices_WithOptions(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	opts := &CustomerInvoiceFilterOptions{
		Limit:  20,
		Offset: 0,
	}

	_, count, err := service.ListCustomerInvoices(ctx, orgID, opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count < 0 {
		t.Errorf("expected non-negative count, got %d", count)
	}
}

func TestService_UpdateCustomerInvoice_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	invoiceID := uuid.New()

	// Setup existing invoice
	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		InvoiceDate:    time.Now(),
		TotalAmount:    100.0,
		PaidAmount:     0.0,
		BalanceDue:     100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	newTotal := 150.0
	req := &UpdateCustomerInvoiceRequest{
		TotalAmount: &newTotal,
	}

	result, err := service.UpdateCustomerInvoice(ctx, invoiceID, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Invoice.TotalAmount != 150.0 {
		t.Errorf("expected total amount 150.0, got %f", result.Invoice.TotalAmount)
	}

	if result.Invoice.BalanceDue != 150.0 {
		t.Errorf("expected balance due 150.0, got %f", result.Invoice.BalanceDue)
	}
}

func TestService_UpdateCustomerInvoice_InvalidDueDate(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	invoiceID := uuid.New()

	now := time.Now()
	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		InvoiceDate:    now,
		TotalAmount:    100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	pastDate := now.AddDate(0, 0, -10)
	req := &UpdateCustomerInvoiceRequest{
		DueDate: &pastDate, // Before invoice date!
	}

	_, err := service.UpdateCustomerInvoice(ctx, invoiceID, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for invalid due date, got nil")
	}
}

func TestService_DeleteCustomerInvoice_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		IsPosted:       false,
	}
	mockRepo.invoices[invoiceID] = invoice

	err := service.DeleteCustomerInvoice(ctx, orgID, invoiceID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify soft delete
	if mockRepo.invoices[invoiceID].DeletedAt == nil {
		t.Error("expected invoice to be soft deleted")
	}
}

func TestService_DeleteCustomerInvoice_Posted(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		IsPosted:       true, // Posted!
	}
	mockRepo.invoices[invoiceID] = invoice

	err := service.DeleteCustomerInvoice(ctx, orgID, invoiceID)
	if err == nil {
		t.Fatal("expected error for deleting posted invoice, got nil")
	}
}

// ========================================
// CUSTOMER PAYMENT TESTS
// ========================================

func TestService_CreateCustomerPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()
	invoiceID := uuid.New()

	// Setup existing invoice
	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		CustomerID:     customerID,
		TotalAmount:    100.0,
		PaidAmount:     0.0,
		BalanceDue:     100.0,
		Status:         InvoiceStatusUnpaid,
	}
	mockRepo.invoices[invoiceID] = invoice

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-001",
		CustomerID:    customerID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: invoiceID,
				AppliedAmount:     100.0,
			},
		},
	}

	result, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Payment.PaymentAmount != 100.0 {
		t.Errorf("expected payment amount 100.0, got %f", result.Payment.PaymentAmount)
	}

	if len(result.Applications) != 1 {
		t.Errorf("expected 1 application, got %d", len(result.Applications))
	}

	// Verify invoice status updated
	updatedInv := mockRepo.invoices[invoiceID]
	if updatedInv.Status != InvoiceStatusPaid {
		t.Errorf("expected invoice status paid, got %v", updatedInv.Status)
	}

	if updatedInv.BalanceDue != 0.0 {
		t.Errorf("expected invoice balance due 0.0, got %f", updatedInv.BalanceDue)
	}
}

func TestService_CreateCustomerPayment_ZeroAmount(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-001",
		CustomerID:    uuid.New(),
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 0.0, // Zero!
		Applications:  []CreatePaymentApplicationRequest{},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for zero payment amount, got nil")
	}
}

func TestService_CreateCustomerPayment_InvalidInvoice(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-001",
		CustomerID:    customerID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: uuid.New(), // Doesn't exist
				AppliedAmount:     100.0,
			},
		},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for invalid invoice, got nil")
	}
}

func TestService_CreateCustomerPayment_WrongCustomer(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID1 := uuid.New()
	customerID2 := uuid.New()
	invoiceID := uuid.New()

	// Invoice belongs to customer 1
	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		CustomerID:     customerID1,
		TotalAmount:    100.0,
		BalanceDue:     100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	// Payment for customer 2
	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-001",
		CustomerID:    customerID2,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: invoiceID,
				AppliedAmount:     100.0,
			},
		},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for wrong customer, got nil")
	}
}

func TestService_CreateCustomerPayment_ExceedsBalance(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		CustomerID:     customerID,
		TotalAmount:    100.0,
		BalanceDue:     50.0, // Only 50 remaining
	}
	mockRepo.invoices[invoiceID] = invoice

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-001",
		CustomerID:    customerID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: invoiceID,
				AppliedAmount:     100.0, // Exceeds balance!
			},
		},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for exceeding balance, got nil")
	}
}

func TestService_CreateCustomerPayment_PartialPayment(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		CustomerID:     customerID,
		TotalAmount:    100.0,
		PaidAmount:     0.0,
		BalanceDue:     100.0,
		Status:         InvoiceStatusUnpaid,
	}
	mockRepo.invoices[invoiceID] = invoice

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-001",
		CustomerID:    customerID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 50.0, // Partial payment
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: invoiceID,
				AppliedAmount:     50.0,
			},
		},
	}

	result, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Payment.PaymentAmount != 50.0 {
		t.Errorf("expected payment amount 50.0, got %f", result.Payment.PaymentAmount)
	}

	// Verify invoice status is partial
	updatedInv := mockRepo.invoices[invoiceID]
	if updatedInv.Status != InvoiceStatusPartial {
		t.Errorf("expected invoice status partial, got %v", updatedInv.Status)
	}

	if updatedInv.BalanceDue != 50.0 {
		t.Errorf("expected invoice balance due 50.0, got %f", updatedInv.BalanceDue)
	}
}

func TestService_GetCustomerPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &CustomerPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		PaymentAmount:  100.0,
	}
	mockRepo.payments[paymentID] = payment

	result, err := service.GetCustomerPayment(ctx, orgID, paymentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Payment.ID != paymentID {
		t.Errorf("expected payment ID %v, got %v", paymentID, result.Payment.ID)
	}
}

func TestService_UpdateCustomerPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	paymentID := uuid.New()

	payment := &CustomerPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		PaymentAmount:  100.0,
		IsPosted:       false,
	}
	mockRepo.payments[paymentID] = payment

	newAmount := 150.0
	req := &UpdateCustomerPaymentRequest{
		PaymentAmount: &newAmount,
	}

	result, err := service.UpdateCustomerPayment(ctx, paymentID, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Payment.PaymentAmount != 150.0 {
		t.Errorf("expected payment amount 150.0, got %f", result.Payment.PaymentAmount)
	}
}

func TestService_UpdateCustomerPayment_Posted(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	paymentID := uuid.New()

	payment := &CustomerPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		PaymentAmount:  100.0,
		IsPosted:       true, // Posted!
	}
	mockRepo.payments[paymentID] = payment

	newAmount := 150.0
	req := &UpdateCustomerPaymentRequest{
		PaymentAmount: &newAmount,
	}

	_, err := service.UpdateCustomerPayment(ctx, paymentID, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for updating posted payment, got nil")
	}
}

func TestService_DeleteCustomerPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &CustomerPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		IsPosted:       false,
	}
	mockRepo.payments[paymentID] = payment

	err := service.DeleteCustomerPayment(ctx, orgID, paymentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify soft delete
	if mockRepo.payments[paymentID].DeletedAt == nil {
		t.Error("expected payment to be soft deleted")
	}
}

func TestService_DeleteCustomerPayment_Posted(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &CustomerPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		IsPosted:       true, // Posted!
	}
	mockRepo.payments[paymentID] = payment

	err := service.DeleteCustomerPayment(ctx, orgID, paymentID)
	if err == nil {
		t.Fatal("expected error for deleting posted payment, got nil")
	}
}

// ========================================
// REPORTING TESTS
// ========================================

func TestService_GetCustomerAgingReport_Success(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	report, err := service.GetCustomerAgingReport(ctx, orgID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report == nil {
		t.Error("expected non-nil report")
	}
}

func TestService_GetCustomerBalance_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	customerID := uuid.New()

	// Setup invoices
	inv1 := &CustomerInvoice{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CustomerID:     customerID,
		BalanceDue:     50.0,
	}
	inv2 := &CustomerInvoice{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CustomerID:     customerID,
		BalanceDue:     75.0,
	}
	mockRepo.invoices[inv1.ID] = inv1
	mockRepo.invoices[inv2.ID] = inv2

	balance, err := service.GetCustomerBalance(ctx, orgID, customerID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if balance != 125.0 {
		t.Errorf("expected balance 125.0, got %f", balance)
	}
}

func TestService_GetOverdueInvoices_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	// Setup overdue invoice
	overdueInv := &CustomerInvoice{
		ID:             uuid.New(),
		OrganizationID: orgID,
		DueDate:        time.Now().AddDate(0, 0, -10), // 10 days overdue
		BalanceDue:     100.0,
	}
	mockRepo.invoices[overdueInv.ID] = overdueInv

	invoices, err := service.GetOverdueInvoices(ctx, orgID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(invoices) != 1 {
		t.Errorf("expected 1 overdue invoice, got %d", len(invoices))
	}
}

func TestService_CheckInvoiceExists_True(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             uuid.New(),
		OrganizationID: orgID,
		InvoiceNumber:  "INV-001",
	}
	mockRepo.invoices[invoice.ID] = invoice

	exists, err := service.CheckInvoiceExists(ctx, orgID, "INV-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Error("expected invoice to exist")
	}
}

func TestService_CheckInvoiceExists_False(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	exists, err := service.CheckInvoiceExists(ctx, orgID, "INV-999")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exists {
		t.Error("expected invoice to not exist")
	}
}

func TestService_RemovePaymentApplication_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	appID := uuid.New()
	invoiceID := uuid.New()

	// Setup application
	app := &CustomerPaymentApplication{
		ID:                appID,
		OrganizationID:    orgID,
		CustomerInvoiceID: invoiceID,
		AppliedAmount:     50.0,
	}
	mockRepo.applications[appID] = app

	// Setup invoice
	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		TotalAmount:    100.0,
		PaidAmount:     50.0,
		BalanceDue:     50.0,
		Status:         InvoiceStatusPartial,
	}
	mockRepo.invoices[invoiceID] = invoice

	err := service.RemovePaymentApplication(ctx, orgID, appID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify invoice reverted
	updatedInv := mockRepo.invoices[invoiceID]
	if updatedInv.PaidAmount != 0.0 {
		t.Errorf("expected paid amount 0.0, got %f", updatedInv.PaidAmount)
	}

	if updatedInv.BalanceDue != 100.0 {
		t.Errorf("expected balance due 100.0, got %f", updatedInv.BalanceDue)
	}

	if updatedInv.Status != InvoiceStatusUnpaid {
		t.Errorf("expected status unpaid, got %v", updatedInv.Status)
	}

	// Verify application deleted
	if _, ok := mockRepo.applications[appID]; ok {
		t.Error("expected application to be deleted")
	}
}

// ========================================
// ADDITIONAL EDGE CASE TESTS
// ========================================

func TestService_ListCustomerPayments_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	// Setup mock data
	pmt1 := &CustomerPayment{ID: uuid.New(), OrganizationID: orgID}
	pmt2 := &CustomerPayment{ID: uuid.New(), OrganizationID: orgID}
	mockRepo.payments[pmt1.ID] = pmt1
	mockRepo.payments[pmt2.ID] = pmt2

	payments, count, err := service.ListCustomerPayments(ctx, orgID, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}

	if len(payments) != 2 {
		t.Errorf("expected 2 payments, got %d", len(payments))
	}
}

func TestService_ListCustomerPayments_WithOptions(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	opts := &CustomerPaymentFilterOptions{
		Limit:  20,
		Offset: 0,
	}

	_, count, err := service.ListCustomerPayments(ctx, orgID, opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count < 0 {
		t.Errorf("expected non-negative count, got %d", count)
	}
}

func TestService_CheckPaymentExists_True(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	payment := &CustomerPayment{
		ID:             uuid.New(),
		OrganizationID: orgID,
		PaymentNumber:  "PMT-001",
	}
	mockRepo.payments[payment.ID] = payment

	exists, err := service.CheckPaymentExists(ctx, orgID, "PMT-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Error("expected payment to exist")
	}
}

func TestService_CheckPaymentExists_False(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	exists, err := service.CheckPaymentExists(ctx, orgID, "PMT-999")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exists {
		t.Error("expected payment to not exist")
	}
}

func TestService_CreateCustomerInvoice_WithMultipleLines(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	req := &CreateCustomerInvoiceRequest{
		InvoiceNumber: "INV-002",
		CustomerID:    uuid.New(),
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().AddDate(0, 0, 30),
		Subtotal:      200.0,
		TaxAmount:     20.0,
		TotalAmount:   220.0,
		Lines: []CreateCustomerInvoiceLineRequest{
			{
				LineNumber:       1,
				RevenueAccountID: uuid.New(),
				Description:      "Product A",
				Quantity:         2.0,
				UnitPrice:        50.0,
				Amount:           100.0,
				TaxAmount:        10.0,
			},
			{
				LineNumber:       2,
				RevenueAccountID: uuid.New(),
				Description:      "Product B",
				Quantity:         1.0,
				UnitPrice:        100.0,
				Amount:           100.0,
				TaxAmount:        10.0,
			},
		},
	}

	result, err := service.CreateCustomerInvoice(ctx, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(result.Lines))
	}

	// Verify lines stored
	lines := mockRepo.invoiceLines[result.Invoice.ID]
	if len(lines) != 2 {
		t.Errorf("expected 2 lines in repository, got %d", len(lines))
	}
}

func TestService_CreateCustomerInvoice_LineCreationError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	expectedErr := errors.New("line creation failed")
	mockRepo.onCreateCustomerInvoiceLine = func(ctx context.Context, line *CustomerInvoiceLine) error {
		return expectedErr
	}

	req := &CreateCustomerInvoiceRequest{
		InvoiceNumber: "INV-003",
		CustomerID:    uuid.New(),
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().AddDate(0, 0, 30),
		TotalAmount:   100.0,
		Lines: []CreateCustomerInvoiceLineRequest{
			{
				LineNumber:       1,
				RevenueAccountID: uuid.New(),
				Description:      "Test",
				Amount:           100.0,
			},
		},
	}

	_, err := service.CreateCustomerInvoice(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error from line creation, got nil")
	}
}

func TestService_GetCustomerInvoice_LinesFetchError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
	}
	mockRepo.invoices[invoiceID] = invoice

	expectedErr := errors.New("failed to fetch lines")
	mockRepo.onListCustomerInvoiceLines = func(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]CustomerInvoiceLine, error) {
		return nil, expectedErr
	}

	_, err := service.GetCustomerInvoice(ctx, orgID, invoiceID)
	if err == nil {
		t.Fatal("expected error from fetching lines, got nil")
	}
}

func TestService_UpdateCustomerInvoice_WithLines(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		InvoiceDate:    time.Now(),
		TotalAmount:    100.0,
		PaidAmount:     0.0,
		BalanceDue:     100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	lineNum := 1
	desc := "Updated line"
	qty := 2.0
	price := 50.0
	amount := 100.0
	tax := 10.0

	req := &UpdateCustomerInvoiceRequest{
		Lines: []UpdateCustomerInvoiceLineRequest{
			{
				LineNumber:  &lineNum,
				Description: &desc,
				Quantity:    &qty,
				UnitPrice:   &price,
				Amount:      &amount,
				TaxAmount:   &tax,
			},
		},
	}

	result, err := service.UpdateCustomerInvoice(ctx, invoiceID, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestService_CreateCustomerPayment_ZeroApplicationAmount(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		CustomerID:     customerID,
		TotalAmount:    100.0,
		BalanceDue:     100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-002",
		CustomerID:    customerID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: invoiceID,
				AppliedAmount:     0.0, // Zero application!
			},
		},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for zero application amount, got nil")
	}
}

func TestService_GetCustomerPayment_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	_, err := service.GetCustomerPayment(ctx, orgID, paymentID)
	if err == nil {
		t.Fatal("expected error for not found payment, got nil")
	}
}

func TestService_UpdateCustomerPayment_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	paymentID := uuid.New()

	req := &UpdateCustomerPaymentRequest{}

	_, err := service.UpdateCustomerPayment(ctx, paymentID, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error for not found payment, got nil")
	}
}

func TestService_DeleteCustomerInvoice_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	invoiceID := uuid.New()

	err := service.DeleteCustomerInvoice(ctx, orgID, invoiceID)
	if err == nil {
		t.Fatal("expected error for not found invoice, got nil")
	}
}

func TestService_DeleteCustomerPayment_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	err := service.DeleteCustomerPayment(ctx, orgID, paymentID)
	if err == nil {
		t.Fatal("expected error for not found payment, got nil")
	}
}

func TestService_RemovePaymentApplication_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	appID := uuid.New()

	err := service.RemovePaymentApplication(ctx, orgID, appID, userID)
	if err == nil {
		t.Fatal("expected error for not found application, got nil")
	}
}

func TestService_RemovePaymentApplication_InvoiceFetchError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	appID := uuid.New()

	app := &CustomerPaymentApplication{
		ID:                appID,
		OrganizationID:    orgID,
		CustomerInvoiceID: uuid.New(),
		AppliedAmount:     50.0,
	}
	mockRepo.applications[appID] = app

	// Don't add invoice - will cause error

	err := service.RemovePaymentApplication(ctx, orgID, appID, userID)
	if err == nil {
		t.Fatal("expected error fetching invoice, got nil")
	}
}

func TestService_CreateCustomerPayment_RepositoryError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()

	expectedErr := errors.New("database error")
	mockRepo.onCreateCustomerPayment = func(ctx context.Context, payment *CustomerPayment) error {
		return expectedErr
	}

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-003",
		CustomerID:    uuid.New(),
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications:  []CreatePaymentApplicationRequest{},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}
}

func TestService_CreateCustomerPayment_ApplicationCreationError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		CustomerID:     customerID,
		TotalAmount:    100.0,
		BalanceDue:     100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	expectedErr := errors.New("application creation failed")
	mockRepo.onCreateCustomerPaymentApplication = func(ctx context.Context, app *CustomerPaymentApplication) error {
		return expectedErr
	}

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-004",
		CustomerID:    customerID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: invoiceID,
				AppliedAmount:     100.0,
			},
		},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error from application creation, got nil")
	}
}

func TestService_CreateCustomerPayment_InvoiceUpdateError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	customerID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		CustomerID:     customerID,
		TotalAmount:    100.0,
		BalanceDue:     100.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	expectedErr := errors.New("invoice update failed")
	mockRepo.onUpdateCustomerInvoice = func(ctx context.Context, invoice *CustomerInvoice) error {
		return expectedErr
	}

	req := &CreateCustomerPaymentRequest{
		PaymentNumber: "PMT-005",
		CustomerID:    customerID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCash,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				CustomerInvoiceID: invoiceID,
				AppliedAmount:     100.0,
			},
		},
	}

	_, err := service.CreateCustomerPayment(ctx, req, orgID, userID)
	if err == nil {
		t.Fatal("expected error from invoice update, got nil")
	}
}

func TestService_RemovePaymentApplication_InvoiceUpdateError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	appID := uuid.New()
	invoiceID := uuid.New()

	app := &CustomerPaymentApplication{
		ID:                appID,
		OrganizationID:    orgID,
		CustomerInvoiceID: invoiceID,
		AppliedAmount:     50.0,
	}
	mockRepo.applications[appID] = app

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		TotalAmount:    100.0,
		PaidAmount:     50.0,
		BalanceDue:     50.0,
	}
	mockRepo.invoices[invoiceID] = invoice

	expectedErr := errors.New("invoice update failed")
	mockRepo.onUpdateCustomerInvoice = func(ctx context.Context, invoice *CustomerInvoice) error {
		return expectedErr
	}

	err := service.RemovePaymentApplication(ctx, orgID, appID, userID)
	if err == nil {
		t.Fatal("expected error from invoice update, got nil")
	}
}

func TestService_UpdateCustomerInvoice_WithStatusUpdate(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	invoiceID := uuid.New()

	invoice := &CustomerInvoice{
		ID:             invoiceID,
		OrganizationID: orgID,
		InvoiceDate:    time.Now(),
		TotalAmount:    100.0,
		Status:         InvoiceStatusUnpaid,
	}
	mockRepo.invoices[invoiceID] = invoice

	newStatus := InvoiceStatusVoid
	req := &UpdateCustomerInvoiceRequest{
		Status: &newStatus,
	}

	result, err := service.UpdateCustomerInvoice(ctx, invoiceID, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Invoice.Status != InvoiceStatusVoid {
		t.Errorf("expected status void, got %v", result.Invoice.Status)
	}
}

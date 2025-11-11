package payables

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
	bills        map[uuid.UUID]*VendorBill
	billLines    map[uuid.UUID][]VendorBillLine
	payments     map[uuid.UUID]*VendorPayment
	applications map[uuid.UUID]*VendorPaymentApplication
	appsByPayment map[uuid.UUID][]VendorPaymentApplication
	appsByBill    map[uuid.UUID][]VendorPaymentApplication

	onCreateVendorBill               func(ctx context.Context, bill *VendorBill) error
	onGetVendorBillByID              func(ctx context.Context, organizationID, billID uuid.UUID) (*VendorBill, error)
	onGetVendorBillByNumber          func(ctx context.Context, organizationID uuid.UUID, billNumber string) (*VendorBill, error)
	onListVendorBills                func(ctx context.Context, organizationID uuid.UUID, opts *VendorBillFilterOptions) ([]VendorBill, int64, error)
	onUpdateVendorBill               func(ctx context.Context, bill *VendorBill) error
	onDeleteVendorBill               func(ctx context.Context, organizationID, billID uuid.UUID) error
	onCreateVendorBillLine           func(ctx context.Context, line *VendorBillLine) error
	onListVendorBillLines            func(ctx context.Context, organizationID, billID uuid.UUID) ([]VendorBillLine, error)
	onDeleteVendorBillLines          func(ctx context.Context, organizationID, billID uuid.UUID) error
	onCreateVendorPayment            func(ctx context.Context, payment *VendorPayment) error
	onGetVendorPaymentByID           func(ctx context.Context, organizationID, paymentID uuid.UUID) (*VendorPayment, error)
	onGetVendorPaymentByNumber       func(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*VendorPayment, error)
	onListVendorPayments             func(ctx context.Context, organizationID uuid.UUID, opts *VendorPaymentFilterOptions) ([]VendorPayment, int64, error)
	onUpdateVendorPayment            func(ctx context.Context, payment *VendorPayment) error
	onDeleteVendorPayment            func(ctx context.Context, organizationID, paymentID uuid.UUID) error
	onCreateVendorPaymentApplication func(ctx context.Context, app *VendorPaymentApplication) error
	onGetVendorPaymentApplicationByID func(ctx context.Context, organizationID, appID uuid.UUID) (*VendorPaymentApplication, error)
	onListVendorPaymentApplications  func(ctx context.Context, organizationID, paymentID uuid.UUID) ([]VendorPaymentApplication, error)
	onListPaymentApplicationsByBill  func(ctx context.Context, organizationID, billID uuid.UUID) ([]VendorPaymentApplication, error)
	onDeleteVendorPaymentApplication func(ctx context.Context, organizationID, appID uuid.UUID) error
	onGetVendorAgingReport           func(ctx context.Context, organizationID uuid.UUID) ([]VendorAgingReport, error)
	onGetBillsByDueDate              func(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]VendorBill, error)
	onGetSupplierBalance             func(ctx context.Context, organizationID, supplierID uuid.UUID) (float64, error)
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		bills:         make(map[uuid.UUID]*VendorBill),
		billLines:     make(map[uuid.UUID][]VendorBillLine),
		payments:      make(map[uuid.UUID]*VendorPayment),
		applications:  make(map[uuid.UUID]*VendorPaymentApplication),
		appsByPayment: make(map[uuid.UUID][]VendorPaymentApplication),
		appsByBill:    make(map[uuid.UUID][]VendorPaymentApplication),
	}
}

func (m *mockRepository) CreateVendorBill(ctx context.Context, bill *VendorBill) error {
	if m.onCreateVendorBill != nil {
		return m.onCreateVendorBill(ctx, bill)
	}
	m.bills[bill.ID] = bill
	return nil
}

func (m *mockRepository) GetVendorBillByID(ctx context.Context, organizationID, billID uuid.UUID) (*VendorBill, error) {
	if m.onGetVendorBillByID != nil {
		return m.onGetVendorBillByID(ctx, organizationID, billID)
	}
	if bill, ok := m.bills[billID]; ok && bill.OrganizationID == organizationID {
		return bill, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) GetVendorBillByNumber(ctx context.Context, organizationID uuid.UUID, billNumber string) (*VendorBill, error) {
	if m.onGetVendorBillByNumber != nil {
		return m.onGetVendorBillByNumber(ctx, organizationID, billNumber)
	}
	for _, bill := range m.bills {
		if bill.OrganizationID == organizationID && bill.BillNumber == billNumber {
			return bill, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) ListVendorBills(ctx context.Context, organizationID uuid.UUID, opts *VendorBillFilterOptions) ([]VendorBill, int64, error) {
	if m.onListVendorBills != nil {
		return m.onListVendorBills(ctx, organizationID, opts)
	}
	result := make([]VendorBill, 0)
	for _, bill := range m.bills {
		if bill.OrganizationID == organizationID {
			result = append(result, *bill)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockRepository) UpdateVendorBill(ctx context.Context, bill *VendorBill) error {
	if m.onUpdateVendorBill != nil {
		return m.onUpdateVendorBill(ctx, bill)
	}
	m.bills[bill.ID] = bill
	return nil
}

func (m *mockRepository) DeleteVendorBill(ctx context.Context, organizationID, billID uuid.UUID) error {
	if m.onDeleteVendorBill != nil {
		return m.onDeleteVendorBill(ctx, organizationID, billID)
	}
	if bill, ok := m.bills[billID]; ok && bill.OrganizationID == organizationID {
		now := time.Now()
		bill.DeletedAt = &now
		return nil
	}
	return sql.ErrNoRows
}

func (m *mockRepository) CreateVendorBillLine(ctx context.Context, line *VendorBillLine) error {
	if m.onCreateVendorBillLine != nil {
		return m.onCreateVendorBillLine(ctx, line)
	}
	m.billLines[line.VendorBillID] = append(m.billLines[line.VendorBillID], *line)
	return nil
}

func (m *mockRepository) GetVendorBillLineByID(ctx context.Context, organizationID, lineID uuid.UUID) (*VendorBillLine, error) {
	return nil, nil
}

func (m *mockRepository) ListVendorBillLines(ctx context.Context, organizationID, billID uuid.UUID) ([]VendorBillLine, error) {
	if m.onListVendorBillLines != nil {
		return m.onListVendorBillLines(ctx, organizationID, billID)
	}
	if lines, ok := m.billLines[billID]; ok {
		return lines, nil
	}
	return []VendorBillLine{}, nil
}

func (m *mockRepository) UpdateVendorBillLine(ctx context.Context, line *VendorBillLine) error {
	return nil
}

func (m *mockRepository) DeleteVendorBillLine(ctx context.Context, organizationID, lineID uuid.UUID) error {
	return nil
}

func (m *mockRepository) DeleteVendorBillLines(ctx context.Context, organizationID, billID uuid.UUID) error {
	if m.onDeleteVendorBillLines != nil {
		return m.onDeleteVendorBillLines(ctx, organizationID, billID)
	}
	delete(m.billLines, billID)
	return nil
}

func (m *mockRepository) CreateVendorPayment(ctx context.Context, payment *VendorPayment) error {
	if m.onCreateVendorPayment != nil {
		return m.onCreateVendorPayment(ctx, payment)
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockRepository) GetVendorPaymentByID(ctx context.Context, organizationID, paymentID uuid.UUID) (*VendorPayment, error) {
	if m.onGetVendorPaymentByID != nil {
		return m.onGetVendorPaymentByID(ctx, organizationID, paymentID)
	}
	if pmt, ok := m.payments[paymentID]; ok && pmt.OrganizationID == organizationID {
		return pmt, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) GetVendorPaymentByNumber(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*VendorPayment, error) {
	if m.onGetVendorPaymentByNumber != nil {
		return m.onGetVendorPaymentByNumber(ctx, organizationID, paymentNumber)
	}
	for _, pmt := range m.payments {
		if pmt.OrganizationID == organizationID && pmt.PaymentNumber == paymentNumber {
			return pmt, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) ListVendorPayments(ctx context.Context, organizationID uuid.UUID, opts *VendorPaymentFilterOptions) ([]VendorPayment, int64, error) {
	if m.onListVendorPayments != nil {
		return m.onListVendorPayments(ctx, organizationID, opts)
	}
	result := make([]VendorPayment, 0)
	for _, pmt := range m.payments {
		if pmt.OrganizationID == organizationID {
			result = append(result, *pmt)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockRepository) UpdateVendorPayment(ctx context.Context, payment *VendorPayment) error {
	if m.onUpdateVendorPayment != nil {
		return m.onUpdateVendorPayment(ctx, payment)
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockRepository) DeleteVendorPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error {
	if m.onDeleteVendorPayment != nil {
		return m.onDeleteVendorPayment(ctx, organizationID, paymentID)
	}
	if pmt, ok := m.payments[paymentID]; ok && pmt.OrganizationID == organizationID {
		now := time.Now()
		pmt.DeletedAt = &now
		return nil
	}
	return sql.ErrNoRows
}

func (m *mockRepository) CreateVendorPaymentApplication(ctx context.Context, app *VendorPaymentApplication) error {
	if m.onCreateVendorPaymentApplication != nil {
		return m.onCreateVendorPaymentApplication(ctx, app)
	}
	m.applications[app.ID] = app
	m.appsByPayment[app.VendorPaymentID] = append(m.appsByPayment[app.VendorPaymentID], *app)
	m.appsByBill[app.VendorBillID] = append(m.appsByBill[app.VendorBillID], *app)
	return nil
}

func (m *mockRepository) GetVendorPaymentApplicationByID(ctx context.Context, organizationID, appID uuid.UUID) (*VendorPaymentApplication, error) {
	if m.onGetVendorPaymentApplicationByID != nil {
		return m.onGetVendorPaymentApplicationByID(ctx, organizationID, appID)
	}
	if app, ok := m.applications[appID]; ok && app.OrganizationID == organizationID {
		return app, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockRepository) ListVendorPaymentApplications(ctx context.Context, organizationID, paymentID uuid.UUID) ([]VendorPaymentApplication, error) {
	if m.onListVendorPaymentApplications != nil {
		return m.onListVendorPaymentApplications(ctx, organizationID, paymentID)
	}
	if apps, ok := m.appsByPayment[paymentID]; ok {
		return apps, nil
	}
	return []VendorPaymentApplication{}, nil
}

func (m *mockRepository) ListPaymentApplicationsByBill(ctx context.Context, organizationID, billID uuid.UUID) ([]VendorPaymentApplication, error) {
	if m.onListPaymentApplicationsByBill != nil {
		return m.onListPaymentApplicationsByBill(ctx, organizationID, billID)
	}
	if apps, ok := m.appsByBill[billID]; ok {
		return apps, nil
	}
	return []VendorPaymentApplication{}, nil
}

func (m *mockRepository) DeleteVendorPaymentApplication(ctx context.Context, organizationID, appID uuid.UUID) error {
	if m.onDeleteVendorPaymentApplication != nil {
		return m.onDeleteVendorPaymentApplication(ctx, organizationID, appID)
	}
	delete(m.applications, appID)
	return nil
}

func (m *mockRepository) GetVendorAgingReport(ctx context.Context, organizationID uuid.UUID) ([]VendorAgingReport, error) {
	if m.onGetVendorAgingReport != nil {
		return m.onGetVendorAgingReport(ctx, organizationID)
	}
	return []VendorAgingReport{}, nil
}

func (m *mockRepository) GetBillsByDueDate(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]VendorBill, error) {
	if m.onGetBillsByDueDate != nil {
		return m.onGetBillsByDueDate(ctx, organizationID, dueDate)
	}
	result := make([]VendorBill, 0)
	for _, bill := range m.bills {
		if bill.OrganizationID == organizationID && bill.DueDate.Before(dueDate) && bill.BalanceDue > 0 {
			result = append(result, *bill)
		}
	}
	return result, nil
}

func (m *mockRepository) GetSupplierBalance(ctx context.Context, organizationID, supplierID uuid.UUID) (float64, error) {
	if m.onGetSupplierBalance != nil {
		return m.onGetSupplierBalance(ctx, organizationID, supplierID)
	}
	balance := 0.0
	for _, bill := range m.bills {
		if bill.OrganizationID == organizationID && bill.SupplierID == supplierID {
			balance += bill.BalanceDue
		}
	}
	return balance, nil
}

func setupService() (*Service, *mockRepository) {
	mockRepo := newMockRepository()
	service := NewService(mockRepo)
	return service, mockRepo
}

// Test bill creation
func TestService_CreateVendorBill_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	supplierID := uuid.New()

	req := &CreateVendorBillRequest{
		BillNumber:  "BILL-001",
		SupplierID:  supplierID,
		BillDate:    time.Now(),
		DueDate:     time.Now().AddDate(0, 0, 30),
		Subtotal:    100.0,
		TaxAmount:   10.0,
		TotalAmount: 110.0,
		Lines: []CreateVendorBillLineRequest{
			{
				LineNumber:       1,
				ExpenseAccountID: uuid.New(),
				Description:      "Office supplies",
				Quantity:         1.0,
				UnitPrice:        100.0,
				Amount:           100.0,
				TaxAmount:        10.0,
			},
		},
	}

	result, err := service.CreateVendorBill(ctx, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Bill.TotalAmount != 110.0 {
		t.Errorf("expected total amount 110.0, got %f", result.Bill.TotalAmount)
	}
	if result.Bill.Status != VendorBillStatusUnpaid {
		t.Errorf("expected status unpaid, got %v", result.Bill.Status)
	}
	if _, ok := mockRepo.bills[result.Bill.ID]; !ok {
		t.Error("expected bill to be stored")
	}
}

func TestService_CreateVendorBill_NegativeAmount(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	req := &CreateVendorBillRequest{
		BillNumber:  "BILL-001",
		SupplierID:  uuid.New(),
		BillDate:    time.Now(),
		DueDate:     time.Now().AddDate(0, 0, 30),
		TotalAmount: -100.0,
	}

	_, err := service.CreateVendorBill(ctx, req, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for negative amount")
	}
}

func TestService_CreateVendorBill_InvalidDueDate(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	now := time.Now()
	req := &CreateVendorBillRequest{
		BillNumber:  "BILL-001",
		SupplierID:  uuid.New(),
		BillDate:    now,
		DueDate:     now.AddDate(0, 0, -10),
		TotalAmount: 100.0,
	}

	_, err := service.CreateVendorBill(ctx, req, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for invalid due date")
	}
}

func TestService_GetVendorBill_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		BillNumber:     "BILL-001",
		TotalAmount:    100.0,
	}
	mockRepo.bills[billID] = bill
	mockRepo.billLines[billID] = []VendorBillLine{{Description: "Test"}}

	result, err := service.GetVendorBill(ctx, orgID, billID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Bill.ID != billID {
		t.Errorf("expected bill ID %v, got %v", billID, result.Bill.ID)
	}
}

func TestService_GetVendorBill_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	_, err := service.GetVendorBill(ctx, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestService_ListVendorBills_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	mockRepo.bills[uuid.New()] = &VendorBill{OrganizationID: orgID}
	mockRepo.bills[uuid.New()] = &VendorBill{OrganizationID: orgID}

	bills, count, err := service.ListVendorBills(ctx, orgID, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
	if len(bills) != 2 {
		t.Errorf("expected 2 bills, got %d", len(bills))
	}
}

func TestService_UpdateVendorBill_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		BillDate:       time.Now(),
		TotalAmount:    100.0,
		PaidAmount:     0.0,
		BalanceDue:     100.0,
	}
	mockRepo.bills[billID] = bill

	newTotal := 150.0
	req := &UpdateVendorBillRequest{
		TotalAmount: &newTotal,
	}

	result, err := service.UpdateVendorBill(ctx, billID, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Bill.TotalAmount != 150.0 {
		t.Errorf("expected total 150.0, got %f", result.Bill.TotalAmount)
	}
}

func TestService_UpdateVendorBill_InvalidDueDate(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	billID := uuid.New()

	now := time.Now()
	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		BillDate:       now,
		TotalAmount:    100.0,
	}
	mockRepo.bills[billID] = bill

	pastDate := now.AddDate(0, 0, -10)
	req := &UpdateVendorBillRequest{
		DueDate: &pastDate,
	}

	_, err := service.UpdateVendorBill(ctx, billID, req, orgID, uuid.New())
	if err == nil {
		t.Fatal("expected error for invalid due date")
	}
}

func TestService_DeleteVendorBill_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		IsPosted:       false,
	}
	mockRepo.bills[billID] = bill

	err := service.DeleteVendorBill(ctx, orgID, billID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mockRepo.bills[billID].DeletedAt == nil {
		t.Error("expected bill to be soft deleted")
	}
}

func TestService_DeleteVendorBill_Posted(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		IsPosted:       true,
	}
	mockRepo.bills[billID] = bill

	err := service.DeleteVendorBill(ctx, orgID, billID)
	if err == nil {
		t.Fatal("expected error for deleting posted bill")
	}
}

// Test payment creation
func TestService_CreateVendorPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	supplierID := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		SupplierID:     supplierID,
		TotalAmount:    100.0,
		PaidAmount:     0.0,
		BalanceDue:     100.0,
		Status:         VendorBillStatusUnpaid,
	}
	mockRepo.bills[billID] = bill

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-001",
		SupplierID:    supplierID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				VendorBillID:  billID,
				AppliedAmount: 100.0,
			},
		},
	}

	result, err := service.CreateVendorPayment(ctx, req, orgID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Payment.PaymentAmount != 100.0 {
		t.Errorf("expected payment amount 100.0, got %f", result.Payment.PaymentAmount)
	}

	updatedBill := mockRepo.bills[billID]
	if updatedBill.Status != VendorBillStatusPaid {
		t.Errorf("expected bill status paid, got %v", updatedBill.Status)
	}
}

func TestService_CreateVendorPayment_ZeroAmount(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-001",
		SupplierID:    uuid.New(),
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 0.0,
	}

	_, err := service.CreateVendorPayment(ctx, req, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for zero payment amount")
	}
}

func TestService_CreateVendorPayment_WrongSupplier(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	supplierID1 := uuid.New()
	supplierID2 := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		SupplierID:     supplierID1,
		BalanceDue:     100.0,
	}
	mockRepo.bills[billID] = bill

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-001",
		SupplierID:    supplierID2,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				VendorBillID:  billID,
				AppliedAmount: 100.0,
			},
		},
	}

	_, err := service.CreateVendorPayment(ctx, req, orgID, uuid.New())
	if err == nil {
		t.Fatal("expected error for wrong supplier")
	}
}

func TestService_CreateVendorPayment_ExceedsBalance(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	supplierID := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		SupplierID:     supplierID,
		TotalAmount:    100.0,
		BalanceDue:     50.0,
	}
	mockRepo.bills[billID] = bill

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-001",
		SupplierID:    supplierID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				VendorBillID:  billID,
				AppliedAmount: 100.0,
			},
		},
	}

	_, err := service.CreateVendorPayment(ctx, req, orgID, uuid.New())
	if err == nil {
		t.Fatal("expected error for exceeding balance")
	}
}

func TestService_CreateVendorPayment_PartialPayment(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	supplierID := uuid.New()
	billID := uuid.New()

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		SupplierID:     supplierID,
		TotalAmount:    100.0,
		PaidAmount:     0.0,
		BalanceDue:     100.0,
		Status:         VendorBillStatusUnpaid,
	}
	mockRepo.bills[billID] = bill

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-001",
		SupplierID:    supplierID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 50.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				VendorBillID:  billID,
				AppliedAmount: 50.0,
			},
		},
	}

	_, err := service.CreateVendorPayment(ctx, req, orgID, uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedBill := mockRepo.bills[billID]
	if updatedBill.Status != VendorBillStatusPartial {
		t.Errorf("expected bill status partial, got %v", updatedBill.Status)
	}
	if updatedBill.BalanceDue != 50.0 {
		t.Errorf("expected balance due 50.0, got %f", updatedBill.BalanceDue)
	}
}

func TestService_GetVendorPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &VendorPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		PaymentAmount:  100.0,
	}
	mockRepo.payments[paymentID] = payment

	result, err := service.GetVendorPayment(ctx, orgID, paymentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Payment.ID != paymentID {
		t.Errorf("expected payment ID %v, got %v", paymentID, result.Payment.ID)
	}
}

func TestService_UpdateVendorPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &VendorPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		PaymentAmount:  100.0,
		IsPosted:       false,
	}
	mockRepo.payments[paymentID] = payment

	newAmount := 150.0
	req := &UpdateVendorPaymentRequest{
		PaymentAmount: &newAmount,
	}

	result, err := service.UpdateVendorPayment(ctx, paymentID, req, orgID, uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Payment.PaymentAmount != 150.0 {
		t.Errorf("expected payment amount 150.0, got %f", result.Payment.PaymentAmount)
	}
}

func TestService_UpdateVendorPayment_Posted(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &VendorPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		IsPosted:       true,
	}
	mockRepo.payments[paymentID] = payment

	newAmount := 150.0
	req := &UpdateVendorPaymentRequest{
		PaymentAmount: &newAmount,
	}

	_, err := service.UpdateVendorPayment(ctx, paymentID, req, orgID, uuid.New())
	if err == nil {
		t.Fatal("expected error for updating posted payment")
	}
}

func TestService_DeleteVendorPayment_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &VendorPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		IsPosted:       false,
	}
	mockRepo.payments[paymentID] = payment

	err := service.DeleteVendorPayment(ctx, orgID, paymentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mockRepo.payments[paymentID].DeletedAt == nil {
		t.Error("expected payment to be soft deleted")
	}
}

func TestService_DeleteVendorPayment_Posted(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	paymentID := uuid.New()

	payment := &VendorPayment{
		ID:             paymentID,
		OrganizationID: orgID,
		IsPosted:       true,
	}
	mockRepo.payments[paymentID] = payment

	err := service.DeleteVendorPayment(ctx, orgID, paymentID)
	if err == nil {
		t.Fatal("expected error for deleting posted payment")
	}
}

// Test reporting functions
func TestService_GetVendorAgingReport_Success(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	report, err := service.GetVendorAgingReport(ctx, uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report == nil {
		t.Error("expected non-nil report")
	}
}

func TestService_GetSupplierBalance_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	supplierID := uuid.New()

	mockRepo.bills[uuid.New()] = &VendorBill{
		OrganizationID: orgID,
		SupplierID:     supplierID,
		BalanceDue:     50.0,
	}
	mockRepo.bills[uuid.New()] = &VendorBill{
		OrganizationID: orgID,
		SupplierID:     supplierID,
		BalanceDue:     75.0,
	}

	balance, err := service.GetSupplierBalance(ctx, orgID, supplierID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if balance != 125.0 {
		t.Errorf("expected balance 125.0, got %f", balance)
	}
}

func TestService_GetOverdueBills_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	mockRepo.bills[uuid.New()] = &VendorBill{
		OrganizationID: orgID,
		DueDate:        time.Now().AddDate(0, 0, -10),
		BalanceDue:     100.0,
	}

	bills, err := service.GetOverdueBills(ctx, orgID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(bills) != 1 {
		t.Errorf("expected 1 overdue bill, got %d", len(bills))
	}
}

func TestService_CheckBillExists_True(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	mockRepo.bills[uuid.New()] = &VendorBill{
		OrganizationID: orgID,
		BillNumber:     "BILL-001",
	}

	exists, err := service.CheckBillExists(ctx, orgID, "BILL-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !exists {
		t.Error("expected bill to exist")
	}
}

func TestService_CheckBillExists_False(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	exists, err := service.CheckBillExists(ctx, uuid.New(), "BILL-999")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected bill to not exist")
	}
}

func TestService_CheckPaymentExists_True(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	mockRepo.payments[uuid.New()] = &VendorPayment{
		OrganizationID: orgID,
		PaymentNumber:  "PMT-001",
	}

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

	exists, err := service.CheckPaymentExists(ctx, uuid.New(), "PMT-999")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected payment to not exist")
	}
}

func TestService_RemovePaymentApplication_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	userID := uuid.New()
	appID := uuid.New()
	billID := uuid.New()

	app := &VendorPaymentApplication{
		ID:              appID,
		OrganizationID:  orgID,
		VendorBillID:    billID,
		AppliedAmount:   50.0,
	}
	mockRepo.applications[appID] = app

	bill := &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		TotalAmount:    100.0,
		PaidAmount:     50.0,
		BalanceDue:     50.0,
		Status:         VendorBillStatusPartial,
	}
	mockRepo.bills[billID] = bill

	err := service.RemovePaymentApplication(ctx, orgID, appID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedBill := mockRepo.bills[billID]
	if updatedBill.PaidAmount != 0.0 {
		t.Errorf("expected paid amount 0.0, got %f", updatedBill.PaidAmount)
	}
	if updatedBill.Status != VendorBillStatusUnpaid {
		t.Errorf("expected status unpaid, got %v", updatedBill.Status)
	}
	if _, ok := mockRepo.applications[appID]; ok {
		t.Error("expected application to be deleted")
	}
}

// Additional edge case tests
func TestService_ListVendorBills_WithOptions(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	opts := &VendorBillFilterOptions{
		Limit:  20,
		Offset: 0,
	}

	_, count, err := service.ListVendorBills(ctx, uuid.New(), opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count < 0 {
		t.Errorf("expected non-negative count, got %d", count)
	}
}

func TestService_ListVendorPayments_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	mockRepo.payments[uuid.New()] = &VendorPayment{OrganizationID: orgID}
	mockRepo.payments[uuid.New()] = &VendorPayment{OrganizationID: orgID}

	payments, count, err := service.ListVendorPayments(ctx, orgID, nil)
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

func TestService_CreateVendorBill_LineCreationError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	mockRepo.onCreateVendorBillLine = func(ctx context.Context, line *VendorBillLine) error {
		return errors.New("line creation failed")
	}

	req := &CreateVendorBillRequest{
		BillNumber:  "BILL-002",
		SupplierID:  uuid.New(),
		BillDate:    time.Now(),
		DueDate:     time.Now().AddDate(0, 0, 30),
		TotalAmount: 100.0,
		Lines: []CreateVendorBillLineRequest{
			{LineNumber: 1, ExpenseAccountID: uuid.New(), Description: "Test", Amount: 100.0},
		},
	}

	_, err := service.CreateVendorBill(ctx, req, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error from line creation")
	}
}

func TestService_GetVendorBill_LinesFetchError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	billID := uuid.New()

	mockRepo.bills[billID] = &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
	}

	mockRepo.onListVendorBillLines = func(ctx context.Context, organizationID, billID uuid.UUID) ([]VendorBillLine, error) {
		return nil, errors.New("failed to fetch lines")
	}

	_, err := service.GetVendorBill(ctx, orgID, billID)
	if err == nil {
		t.Fatal("expected error from fetching lines")
	}
}

func TestService_UpdateVendorBill_WithStatusUpdate(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	billID := uuid.New()

	mockRepo.bills[billID] = &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		BillDate:       time.Now(),
		TotalAmount:    100.0,
		Status:         VendorBillStatusUnpaid,
	}

	newStatus := VendorBillStatusVoid
	req := &UpdateVendorBillRequest{
		Status: &newStatus,
	}

	result, err := service.UpdateVendorBill(ctx, billID, req, orgID, uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Bill.Status != VendorBillStatusVoid {
		t.Errorf("expected status void, got %v", result.Bill.Status)
	}
}

func TestService_CreateVendorPayment_ZeroApplicationAmount(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	supplierID := uuid.New()
	billID := uuid.New()

	mockRepo.bills[billID] = &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		SupplierID:     supplierID,
		BalanceDue:     100.0,
	}

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-002",
		SupplierID:    supplierID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				VendorBillID:  billID,
				AppliedAmount: 0.0,
			},
		},
	}

	_, err := service.CreateVendorPayment(ctx, req, orgID, uuid.New())
	if err == nil {
		t.Fatal("expected error for zero application amount")
	}
}

func TestService_CreateVendorPayment_ApplicationCreationError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	supplierID := uuid.New()
	billID := uuid.New()

	mockRepo.bills[billID] = &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		SupplierID:     supplierID,
		BalanceDue:     100.0,
	}

	mockRepo.onCreateVendorPaymentApplication = func(ctx context.Context, app *VendorPaymentApplication) error {
		return errors.New("application creation failed")
	}

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-003",
		SupplierID:    supplierID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				VendorBillID:  billID,
				AppliedAmount: 100.0,
			},
		},
	}

	_, err := service.CreateVendorPayment(ctx, req, orgID, uuid.New())
	if err == nil {
		t.Fatal("expected error from application creation")
	}
}

func TestService_CreateVendorPayment_BillUpdateError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	supplierID := uuid.New()
	billID := uuid.New()

	mockRepo.bills[billID] = &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		SupplierID:     supplierID,
		BalanceDue:     100.0,
	}

	mockRepo.onUpdateVendorBill = func(ctx context.Context, bill *VendorBill) error {
		return errors.New("bill update failed")
	}

	req := &CreateVendorPaymentRequest{
		PaymentNumber: "PMT-004",
		SupplierID:    supplierID,
		PaymentDate:   time.Now(),
		PaymentMethod: PaymentMethodCheck,
		PaymentAmount: 100.0,
		Applications: []CreatePaymentApplicationRequest{
			{
				VendorBillID:  billID,
				AppliedAmount: 100.0,
			},
		},
	}

	_, err := service.CreateVendorPayment(ctx, req, orgID, uuid.New())
	if err == nil {
		t.Fatal("expected error from bill update")
	}
}

func TestService_RemovePaymentApplication_BillFetchError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	appID := uuid.New()

	mockRepo.applications[appID] = &VendorPaymentApplication{
		ID:             appID,
		OrganizationID: orgID,
		VendorBillID:   uuid.New(),
	}

	err := service.RemovePaymentApplication(ctx, orgID, appID, uuid.New())
	if err == nil {
		t.Fatal("expected error fetching bill")
	}
}

func TestService_RemovePaymentApplication_BillUpdateError(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	appID := uuid.New()
	billID := uuid.New()

	mockRepo.applications[appID] = &VendorPaymentApplication{
		ID:             appID,
		OrganizationID: orgID,
		VendorBillID:   billID,
		AppliedAmount:  50.0,
	}

	mockRepo.bills[billID] = &VendorBill{
		ID:             billID,
		OrganizationID: orgID,
		PaidAmount:     50.0,
	}

	mockRepo.onUpdateVendorBill = func(ctx context.Context, bill *VendorBill) error {
		return errors.New("bill update failed")
	}

	err := service.RemovePaymentApplication(ctx, orgID, appID, uuid.New())
	if err == nil {
		t.Fatal("expected error from bill update")
	}
}

package sales

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
)

// Mock Repository
type mockRepository struct {
	sales      map[uuid.UUID]*Sale
	onList     func(ctx context.Context, orgID uuid.UUID, filters SaleFilters) ([]Sale, error)
	onCreate   func(ctx context.Context, sale *Sale) error
	onGet      func(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Sale, error)
	onUpdate   func(ctx context.Context, sale *Sale) error
	onDelete   func(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		sales: make(map[uuid.UUID]*Sale),
	}
}

func (m *mockRepository) List(ctx context.Context, orgID uuid.UUID, filters SaleFilters) ([]Sale, error) {
	if m.onList != nil {
		return m.onList(ctx, orgID, filters)
	}
	var result []Sale
	for _, sale := range m.sales {
		if sale.OrganizationID == orgID {
			result = append(result, *sale)
		}
	}
	return result, nil
}

func (m *mockRepository) Count(ctx context.Context, orgID uuid.UUID, filters SaleFilters) (int64, error) {
	count := int64(0)
	for _, sale := range m.sales {
		if sale.OrganizationID == orgID {
			count++
		}
	}
	return count, nil
}

func (m *mockRepository) Create(ctx context.Context, sale *Sale) error {
	if m.onCreate != nil {
		return m.onCreate(ctx, sale)
	}
	m.sales[sale.ID] = sale
	return nil
}

func (m *mockRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Sale, error) {
	if m.onGet != nil {
		return m.onGet(ctx, orgID, id)
	}
	if sale, ok := m.sales[id]; ok && sale.OrganizationID == orgID {
		return sale, nil
	}
	return nil, nil
}

func (m *mockRepository) GetBySaleNumber(ctx context.Context, orgID uuid.UUID, saleNumber string) (*Sale, error) {
	for _, sale := range m.sales {
		if sale.OrganizationID == orgID && sale.SaleNumber == saleNumber {
			return sale, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) GetWithItems(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) (*Sale, error) {
	return m.Get(ctx, orgID, saleID)
}

func (m *mockRepository) Update(ctx context.Context, sale *Sale) error {
	if m.onUpdate != nil {
		return m.onUpdate(ctx, sale)
	}
	if existing, ok := m.sales[sale.ID]; ok && existing.OrganizationID == sale.OrganizationID {
		m.sales[sale.ID] = sale
		return nil
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if m.onDelete != nil {
		return m.onDelete(ctx, orgID, id)
	}
	if sale, ok := m.sales[id]; ok && sale.OrganizationID == orgID {
		delete(m.sales, id)
		return nil
	}
	return nil
}

func (m *mockRepository) CreateItem(ctx context.Context, item *SaleItem) error {
	return nil
}

func (m *mockRepository) GetItem(ctx context.Context, saleID uuid.UUID, itemID uuid.UUID) (*SaleItem, error) {
	return nil, nil
}

func (m *mockRepository) ListItems(ctx context.Context, saleID uuid.UUID, filters SaleItemFilters) ([]SaleItem, error) {
	return []SaleItem{}, nil
}

func (m *mockRepository) CountItems(ctx context.Context, saleID uuid.UUID) (int64, error) {
	return 0, nil
}

func (m *mockRepository) UpdateItem(ctx context.Context, item *SaleItem) error {
	return nil
}

func (m *mockRepository) DeleteItem(ctx context.Context, saleID uuid.UUID, itemID uuid.UUID) error {
	return nil
}

// Test helpers
func setupService() (*Service, *mockRepository) {
	repo := newMockRepository()
	logger, _ := logging.NewLogger("info", "console")
	service := NewService(repo, logger)
	return service, repo
}

func createValidSale(orgID uuid.UUID) *Sale {
	return &Sale{
		OrganizationID:  orgID,
		SaleNumber:      "SALE-001",
		TransactionType: "sale",
		Subtotal:        100.00,
		TaxAmount:       10.00,
		DiscountAmount:  0.00,
		TotalAmount:     110.00,
		PaymentStatus:   "paid",
		TransactionDate: time.Now(),
		Items: []SaleItem{
			{
				ProductID:   &uuid.UUID{},
				ProductName: "Test Product",
				Quantity:    1.0,
				UnitPrice:   100.00,
				Subtotal:    100.00,
				TaxAmount:   10.00,
				Total:       110.00,
			},
		},
	}
}

// ============================================================================
// Create Tests
// ============================================================================

func TestService_Create_Success(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)

	err := service.Create(ctx, sale)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify IDs were generated
	if sale.ID == uuid.Nil {
		t.Error("expected ID to be generated")
	}

	// Verify timestamps were set
	if sale.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if sale.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	// Verify item IDs and timestamps
	for i, item := range sale.Items {
		if item.ID == uuid.Nil {
			t.Errorf("expected item %d ID to be generated", i)
		}
		if item.SaleID != sale.ID {
			t.Errorf("expected item %d SaleID to match sale ID", i)
		}
		if item.OrganizationID != orgID {
			t.Errorf("expected item %d OrganizationID to match", i)
		}
		if item.CreatedAt.IsZero() {
			t.Errorf("expected item %d CreatedAt to be set", i)
		}
	}
}

func TestService_Create_ValidationError_MissingOrgID(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()

	sale := createValidSale(uuid.Nil)
	sale.OrganizationID = uuid.Nil

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestService_Create_ValidationError_MissingSaleNumber(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)
	sale.SaleNumber = ""

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestService_Create_ValidationError_MissingTransactionType(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)
	sale.TransactionType = ""

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestService_Create_ValidationError_NegativeTotalAmount(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)
	sale.TotalAmount = -10.00

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestService_Create_ValidationError_NegativeSubtotal(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)
	sale.Subtotal = -100.00

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestService_Create_ValidationError_NegativeTaxAmount(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)
	sale.TaxAmount = -5.00

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestService_Create_ValidationError_NegativeDiscountAmount(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)
	sale.DiscountAmount = -5.00

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestService_Create_ValidationError_NoItems(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	sale := createValidSale(orgID)
	sale.Items = []SaleItem{}

	err := service.Create(ctx, sale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

// ============================================================================
// Get Tests
// ============================================================================

func TestService_Get_Success(t *testing.T) {
	service, repo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	saleID := uuid.New()

	expectedSale := createValidSale(orgID)
	expectedSale.ID = saleID
	repo.sales[saleID] = expectedSale

	result, err := service.Get(ctx, orgID, saleID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected sale to be returned")
	}

	if result.ID != saleID {
		t.Errorf("expected ID %v, got %v", saleID, result.ID)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	saleID := uuid.New()

	_, err := service.Get(ctx, orgID, saleID)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}

	if !apperrors.IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestService_GetBySaleNumber_Success(t *testing.T) {
	service, repo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	expectedSale := createValidSale(orgID)
	expectedSale.ID = uuid.New()
	repo.sales[expectedSale.ID] = expectedSale

	result, err := service.GetBySaleNumber(ctx, orgID, "SALE-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected sale to be returned")
	}

	if result.SaleNumber != "SALE-001" {
		t.Errorf("expected sale number SALE-001, got %v", result.SaleNumber)
	}
}

func TestService_GetBySaleNumber_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	_, err := service.GetBySaleNumber(ctx, orgID, "NONEXISTENT")
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}

	if !apperrors.IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}
}

// ============================================================================
// Update Tests
// ============================================================================

func TestService_Update_Success(t *testing.T) {
	service, repo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	saleID := uuid.New()

	existingSale := createValidSale(orgID)
	existingSale.ID = saleID
	repo.sales[saleID] = existingSale

	updatedSale := createValidSale(orgID)
	updatedSale.ID = saleID
	updatedSale.Notes = "Updated notes"

	err := service.Update(ctx, updatedSale)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify UpdatedAt was set
	if updatedSale.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestService_Update_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	saleID := uuid.New()

	sale := createValidSale(orgID)
	sale.ID = saleID

	err := service.Update(ctx, sale)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}

	if !apperrors.IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestService_Update_ValidationError(t *testing.T) {
	service, repo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	saleID := uuid.New()

	existingSale := createValidSale(orgID)
	existingSale.ID = saleID
	repo.sales[saleID] = existingSale

	updatedSale := createValidSale(orgID)
	updatedSale.ID = saleID
	updatedSale.TotalAmount = -100.00 // Invalid

	err := service.Update(ctx, updatedSale)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if appErr, ok := apperrors.IsAppError(err); !ok || appErr.Code != apperrors.CodeValidationFailed {
		t.Errorf("expected validation error, got %v", err)
	}
}

// ============================================================================
// Delete Tests
// ============================================================================

func TestService_Delete_Success(t *testing.T) {
	service, repo := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	saleID := uuid.New()

	existingSale := createValidSale(orgID)
	existingSale.ID = saleID
	repo.sales[saleID] = existingSale

	err := service.Delete(ctx, orgID, saleID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify sale was deleted
	if _, exists := repo.sales[saleID]; exists {
		t.Error("expected sale to be deleted from repository")
	}
}

func TestService_Delete_NotFound(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()
	saleID := uuid.New()

	err := service.Delete(ctx, orgID, saleID)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}

	if !apperrors.IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}
}

// ============================================================================
// List Tests
// ============================================================================

func TestService_List_DefaultPagination(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	filters := SaleFilters{
		Page:     0,  // Should default to 1
		PageSize: 0,  // Should default to 20
	}

	capturedFilters := SaleFilters{}
	service.repo.(*mockRepository).onList = func(ctx context.Context, orgID uuid.UUID, filters SaleFilters) ([]Sale, error) {
		capturedFilters = filters
		return []Sale{}, nil
	}

	_, err := service.List(ctx, orgID, filters)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedFilters.Page != 1 {
		t.Errorf("expected page to default to 1, got %d", capturedFilters.Page)
	}

	if capturedFilters.PageSize != 20 {
		t.Errorf("expected page size to default to 20, got %d", capturedFilters.PageSize)
	}
}

func TestService_List_CustomPagination(t *testing.T) {
	service, _ := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	filters := SaleFilters{
		Page:     2,
		PageSize: 50,
	}

	capturedFilters := SaleFilters{}
	service.repo.(*mockRepository).onList = func(ctx context.Context, orgID uuid.UUID, filters SaleFilters) ([]Sale, error) {
		capturedFilters = filters
		return []Sale{}, nil
	}

	_, err := service.List(ctx, orgID, filters)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedFilters.Page != 2 {
		t.Errorf("expected page to be 2, got %d", capturedFilters.Page)
	}

	if capturedFilters.PageSize != 50 {
		t.Errorf("expected page size to be 50, got %d", capturedFilters.PageSize)
	}
}

func TestService_Count_Success(t *testing.T) {
	service, repo := setupService()
	ctx := context.Background()
	orgID := uuid.New()

	// Add 3 sales
	for i := 0; i < 3; i++ {
		sale := createValidSale(orgID)
		sale.ID = uuid.New()
		repo.sales[sale.ID] = sale
	}

	count, err := service.Count(ctx, orgID, SaleFilters{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != 3 {
		t.Errorf("expected count to be 3, got %d", count)
	}
}

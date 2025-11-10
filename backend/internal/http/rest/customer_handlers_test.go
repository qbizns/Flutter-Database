package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"github.com/your-org/pos-backend/internal/testhelpers"
)

func TestListCustomersHandler(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	// Create test organization
	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")

	// Create test user
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("success - empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListCustomersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("success - with pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers?page=1&page_size=10", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListCustomersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("success - with filters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers?search=john&customer_type=regular&is_active=true", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListCustomersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		rec := httptest.NewRecorder()

		handler := ListCustomersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestCreateCustomerHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("success", func(t *testing.T) {
		customer := CreateCustomerRequest{
			FirstName:       "John",
			LastName:        "Doe",
			Email:           "john.doe@example.com",
			Phone:           "+1234567890",
			AddressLine1:    "123 Main St",
			City:            "New York",
			State:           "NY",
			PostalCode:      "10001",
			Country:         "US",
			CustomerType:    "regular",
			PaymentTermDays: 30,
			CreditLimit:     1000.00,
			IsActive:        true,
		}

		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify response contains customer ID
		var response map[string]interface{}
		json.NewDecoder(rec.Body).Decode(&response)

		if response["id"] == nil {
			t.Error("response should contain customer ID")
		}
	})

	t.Run("failure - missing required fields", func(t *testing.T) {
		customer := CreateCustomerRequest{
			// Missing FirstName and Email
			LastName: "Doe",
		}

		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid email format", func(t *testing.T) {
		customer := CreateCustomerRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid-email",
		}

		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - negative credit limit", func(t *testing.T) {
		customer := CreateCustomerRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john@example.com",
			CreditLimit: -100.00,
		}

		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		customer := CreateCustomerRequest{FirstName: "John", Email: "john@example.com"}
		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		handler := CreateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("failure - missing user context", func(t *testing.T) {
		customer := CreateCustomerRequest{FirstName: "John", Email: "john@example.com"}
		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))

		rec := httptest.NewRecorder()

		handler := CreateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestGetCustomerHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - customer not found", func(t *testing.T) {
		customerID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/customers/"+customerID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		// Add chi context for URL params
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid customer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		customerID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/customers/"+customerID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestUpdateCustomerHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - customer not found", func(t *testing.T) {
		customerID := uuid.New()

		newFirstName := "Jane"
		update := UpdateCustomerRequest{
			FirstName: &newFirstName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/customers/"+customerID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		customerID := uuid.New()

		req := httptest.NewRequest(http.MethodPatch, "/customers/"+customerID.String(), bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid email format", func(t *testing.T) {
		customerID := uuid.New()

		invalidEmail := "not-an-email"
		update := UpdateCustomerRequest{
			Email: &invalidEmail,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/customers/"+customerID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid customer ID", func(t *testing.T) {
		newFirstName := "Jane"
		update := UpdateCustomerRequest{
			FirstName: &newFirstName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/customers/invalid-uuid", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		customerID := uuid.New()
		newFirstName := "Jane"
		update := UpdateCustomerRequest{FirstName: &newFirstName}
		body, _ := json.Marshal(update)

		req := httptest.NewRequest(http.MethodPatch, "/customers/"+customerID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("failure - missing user context", func(t *testing.T) {
		customerID := uuid.New()
		newFirstName := "Jane"
		update := UpdateCustomerRequest{FirstName: &newFirstName}
		body, _ := json.Marshal(update)

		req := httptest.NewRequest(http.MethodPatch, "/customers/"+customerID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestDeleteCustomerHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger(&config.Config{
		Server: config.ServerConfig{Env: "test"},
	})

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - customer not found", func(t *testing.T) {
		customerID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/customers/"+customerID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid customer ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/customers/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		customerID := uuid.New()
		req := httptest.NewRequest(http.MethodDelete, "/customers/"+customerID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", customerID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteCustomerHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

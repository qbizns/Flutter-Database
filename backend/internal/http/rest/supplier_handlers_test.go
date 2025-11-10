package rest

import (
	"context"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	"github.com/your-org/pos-backend/internal/testhelpers"
)

func TestListSuppliersHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("success - empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListSuppliersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("success - with pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers?page=1&page_size=10", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListSuppliersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("success - with filters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers?search=acme&status=active", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListSuppliersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		rec := httptest.NewRecorder()

		handler := ListSuppliersHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestCreateSupplierHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("success", func(t *testing.T) {
		supplier := CreateSupplierRequest{
			Name:          "ACME Corp",
			ContactPerson: "John Doe",
			Email:         "contact@acme.com",
			Phone:         "+1234567890",
			AddressLine1:  "123 Business St",
			City:          "New York",
			State:         "NY",
			PostalCode:    "10001",
			Country:       "US",
			TaxID:         "TAX123456",
			CreditLimit:   10000.00,
		}

		body, _ := json.Marshal(supplier)
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify response contains supplier ID
		var response map[string]interface{}
		json.NewDecoder(rec.Body).Decode(&response)

		if response["id"] == nil {
			t.Error("response should contain supplier ID")
		}
	})

	t.Run("failure - missing required fields", func(t *testing.T) {
		supplier := CreateSupplierRequest{
			// Missing Name and Email
			ContactPerson: "John Doe",
		}

		body, _ := json.Marshal(supplier)
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid email format", func(t *testing.T) {
		supplier := CreateSupplierRequest{
			Name:  "ACME Corp",
			Email: "not-an-email",
		}

		body, _ := json.Marshal(supplier)
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - negative credit limit", func(t *testing.T) {
		supplier := CreateSupplierRequest{
			Name:        "ACME Corp",
			Email:       "contact@acme.com",
			CreditLimit: -500.00,
		}

		body, _ := json.Marshal(supplier)
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		supplier := CreateSupplierRequest{Name: "ACME", Email: "contact@acme.com"}
		body, _ := json.Marshal(supplier)
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		handler := CreateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("failure - missing user context", func(t *testing.T) {
		supplier := CreateSupplierRequest{Name: "ACME", Email: "contact@acme.com"}
		body, _ := json.Marshal(supplier)
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))

		rec := httptest.NewRecorder()

		handler := CreateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestGetSupplierHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - supplier not found", func(t *testing.T) {
		supplierID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/"+supplierID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", supplierID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid supplier ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		supplierID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/suppliers/"+supplierID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", supplierID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestUpdateSupplierHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - supplier not found", func(t *testing.T) {
		supplierID := uuid.New()

		newName := "Updated ACME Corp"
		update := UpdateSupplierRequest{
			Name: &newName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/suppliers/"+supplierID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", supplierID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		supplierID := uuid.New()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/"+supplierID.String(), bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", supplierID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid email format", func(t *testing.T) {
		supplierID := uuid.New()

		invalidEmail := "not-an-email"
		update := UpdateSupplierRequest{
			Email: &invalidEmail,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/suppliers/"+supplierID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", supplierID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid supplier ID", func(t *testing.T) {
		newName := "Updated"
		update := UpdateSupplierRequest{
			Name: &newName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/suppliers/invalid-uuid", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})
}

func TestDeleteSupplierHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - supplier not found", func(t *testing.T) {
		supplierID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/"+supplierID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", supplierID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid supplier ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/suppliers/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		supplierID := uuid.New()
		req := httptest.NewRequest(http.MethodDelete, "/suppliers/"+supplierID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", supplierID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteSupplierHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

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

func TestListProductsHandler(t *testing.T) {
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
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListProductsHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()

		handler := ListProductsHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestCreateProductHandler(t *testing.T) {
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
		product := CreateProductRequest{
			Name:              "Test Product",
			SKU:               "TEST-001",
			Barcode:           "1234567890",
			Description:       "A test product",
			UnitPrice:         9.99,
			Cost:              5.00,
			TaxRate:           0.08,
			Unit:              "piece",
			MinStockLevel:     10,
			MaxStockLevel:     100,
			IsActive:          true,
			IsTrackInventory:  true,
			AllowNegativeStock: false,
			ProductType:       "simple",
		}

		body, _ := json.Marshal(product)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify response contains product ID
		var response map[string]interface{}
		json.NewDecoder(rec.Body).Decode(&response)

		if response["id"] == nil {
			t.Error("response should contain product ID")
		}
	})

	t.Run("failure - missing required fields", func(t *testing.T) {
		product := CreateProductRequest{
			// Missing name
			SKU: "TEST-002",
		}

		body, _ := json.Marshal(product)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		product := CreateProductRequest{Name: "Test"}
		body, _ := json.Marshal(product)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		handler := CreateProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("failure - missing user context", func(t *testing.T) {
		product := CreateProductRequest{Name: "Test"}
		body, _ := json.Marshal(product)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))

		rec := httptest.NewRecorder()

		handler := CreateProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestGetProductHandler(t *testing.T) {
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

	t.Run("failure - product not found", func(t *testing.T) {
		productID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		// Add chi context for URL params
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", productID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid product ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		productID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", productID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestUpdateProductHandler(t *testing.T) {
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

	t.Run("failure - product not found", func(t *testing.T) {
		productID := uuid.New()

		newName := "Updated Name"
		update := UpdateProductRequest{
			Name: &newName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/products/"+productID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", productID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		productID := uuid.New()

		req := httptest.NewRequest(http.MethodPatch, "/products/"+productID.String(), bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", productID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})
}

func TestDeleteProductHandler(t *testing.T) {
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

	t.Run("failure - product not found", func(t *testing.T) {
		productID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", productID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid product ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteProductHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})
}

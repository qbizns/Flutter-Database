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

func TestListCategoriesHandler(t *testing.T) {
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

	t.Run("success - empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListCategoriesHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("success - with pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories?page=1&page_size=10", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListCategoriesHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("success - with filters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories?search=test&is_active=true", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListCategoriesHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()

		handler := ListCategoriesHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestCreateCategoryHandler(t *testing.T) {
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
		category := CreateCategoryRequest{
			Name:         "Electronics",
			Description:  "Electronic items and accessories",
			Color:        "#FF5733",
			IconName:     "electronics",
			DisplayOrder: 1,
			IsActive:     true,
		}

		body, _ := json.Marshal(category)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify response contains category ID
		var response map[string]interface{}
		json.NewDecoder(rec.Body).Decode(&response)

		if response["id"] == nil {
			t.Error("response should contain category ID")
		}
	})

	t.Run("failure - missing required name", func(t *testing.T) {
		category := CreateCategoryRequest{
			// Missing Name
			Description: "Test category",
		}

		body, _ := json.Marshal(category)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - name too short", func(t *testing.T) {
		category := CreateCategoryRequest{
			Name: "AB", // Less than 3 characters
		}

		body, _ := json.Marshal(category)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		category := CreateCategoryRequest{Name: "Test Category"}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		handler := CreateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("failure - missing user context", func(t *testing.T) {
		category := CreateCategoryRequest{Name: "Test Category"}
		body, _ := json.Marshal(category)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))

		rec := httptest.NewRecorder()

		handler := CreateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestGetCategoryHandler(t *testing.T) {
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

	t.Run("failure - category not found", func(t *testing.T) {
		categoryID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/categories/"+categoryID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", categoryID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		categoryID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/categories/"+categoryID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", categoryID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestUpdateCategoryHandler(t *testing.T) {
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

	t.Run("failure - category not found", func(t *testing.T) {
		categoryID := uuid.New()

		newName := "Updated Electronics"
		update := UpdateCategoryRequest{
			Name: &newName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/categories/"+categoryID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", categoryID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		categoryID := uuid.New()

		req := httptest.NewRequest(http.MethodPatch, "/categories/"+categoryID.String(), bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", categoryID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - name too short", func(t *testing.T) {
		categoryID := uuid.New()

		shortName := "AB" // Less than 3 characters
		update := UpdateCategoryRequest{
			Name: &shortName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/categories/"+categoryID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", categoryID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid category ID", func(t *testing.T) {
		newName := "Updated"
		update := UpdateCategoryRequest{
			Name: &newName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/categories/invalid-uuid", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})
}

func TestDeleteCategoryHandler(t *testing.T) {
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

	t.Run("failure - category not found", func(t *testing.T) {
		categoryID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", categoryID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		categoryID := uuid.New()
		req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", categoryID.String())
		req = req.WithContext(chi.NewRouteContext.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteCategoryHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

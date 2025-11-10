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

func TestListLocationsHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("success - empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/locations", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListLocationsHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("success - with filters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/locations?search=warehouse&location_type=warehouse&is_active=true", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := ListLocationsHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/locations", nil)
		rec := httptest.NewRecorder()

		handler := ListLocationsHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestCreateLocationHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("success", func(t *testing.T) {
		location := CreateLocationRequest{
			Name:         "Main Warehouse",
			Code:         "WH-001",
			LocationType: "warehouse",
			AddressLine1: "123 Storage St",
			City:         "Los Angeles",
			State:        "CA",
			PostalCode:   "90001",
			Country:      "US",
			Phone:        "+1234567890",
			Email:        "warehouse@example.com",
			IsActive:     true,
		}

		body, _ := json.Marshal(location)
		req := httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
		}

		// Verify response contains location ID
		var response map[string]interface{}
		json.NewDecoder(rec.Body).Decode(&response)

		if response["id"] == nil {
			t.Error("response should contain location ID")
		}
	})

	t.Run("failure - missing required fields", func(t *testing.T) {
		location := CreateLocationRequest{
			// Missing Name and Code
			LocationType: "warehouse",
		}

		body, _ := json.Marshal(location)
		req := httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rec := httptest.NewRecorder()

		handler := CreateLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		location := CreateLocationRequest{Name: "Test", Code: "TEST"}
		body, _ := json.Marshal(location)
		req := httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		handler := CreateLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("failure - missing user context", func(t *testing.T) {
		location := CreateLocationRequest{Name: "Test", Code: "TEST"}
		body, _ := json.Marshal(location)
		req := httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))

		rec := httptest.NewRecorder()

		handler := CreateLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestGetLocationHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - location not found", func(t *testing.T) {
		locationID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/locations/"+locationID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", locationID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid location ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/locations/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		locationID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/locations/"+locationID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", locationID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := GetLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestUpdateLocationHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - location not found", func(t *testing.T) {
		locationID := uuid.New()

		newName := "Updated Warehouse"
		update := UpdateLocationRequest{
			Name: &newName,
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest(http.MethodPatch, "/locations/"+locationID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", locationID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid JSON", func(t *testing.T) {
		locationID := uuid.New()

		req := httptest.NewRequest(http.MethodPatch, "/locations/"+locationID.String(), bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", locationID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := UpdateLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})
}

func TestDeleteLocationHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testhelpers.SetupTestDB(t)
	logger, _ := logging.NewLogger("info", "console")

	ctx := testDB.Context()
	orgID := testDB.CreateTestOrganization(ctx, "Test Org")
	userID := testDB.CreateTestUser(ctx, orgID, "test@example.com")

	t.Run("failure - location not found", func(t *testing.T) {
		locationID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/locations/"+locationID.String(), nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", locationID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("failure - invalid location ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/locations/invalid-uuid", nil)
		req = req.WithContext(appctx.WithOrganizationID(req.Context(), orgID))
		req = req.WithContext(appctx.WithUserID(req.Context(), userID))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("failure - missing organization context", func(t *testing.T) {
		locationID := uuid.New()
		req := httptest.NewRequest(http.MethodDelete, "/locations/"+locationID.String(), nil)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", locationID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		rec := httptest.NewRecorder()

		handler := DeleteLocationHandler(testDB.DB, logger)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

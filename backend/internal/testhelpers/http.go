package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
)

// HTTPTestHelper provides utilities for HTTP handler testing
type HTTPTestHelper struct {
	t *testing.T
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper(t *testing.T) *HTTPTestHelper {
	return &HTTPTestHelper{t: t}
}

// NewRequest creates a new HTTP test request
func (h *HTTPTestHelper) NewRequest(method, path string, body interface{}) *http.Request {
	h.t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			h.t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

// NewAuthenticatedRequest creates a request with auth context
func (h *HTTPTestHelper) NewAuthenticatedRequest(method, path string, body interface{}, orgID, userID uuid.UUID) *http.Request {
	h.t.Helper()

	req := h.NewRequest(method, path, body)

	// Add organization and user to context
	ctx := req.Context()
	ctx = appctx.WithOrganizationID(ctx, orgID)
	ctx = appctx.WithUserID(ctx, userID)

	return req.WithContext(ctx)
}

// NewRecorder creates a new response recorder
func (h *HTTPTestHelper) NewRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// ExecuteRequest executes a request and returns the response
func (h *HTTPTestHelper) ExecuteRequest(handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	h.t.Helper()

	rec := h.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec
}

// AssertStatusCode asserts the response status code
func (h *HTTPTestHelper) AssertStatusCode(rec *httptest.ResponseRecorder, expected int) {
	h.t.Helper()

	if rec.Code != expected {
		h.t.Errorf("expected status code %d, got %d\nResponse body: %s",
			expected, rec.Code, rec.Body.String())
	}
}

// AssertJSONResponse decodes and returns the JSON response
func (h *HTTPTestHelper) AssertJSONResponse(rec *httptest.ResponseRecorder, target interface{}) {
	h.t.Helper()

	if err := json.NewDecoder(rec.Body).Decode(target); err != nil {
		h.t.Fatalf("failed to decode JSON response: %v\nBody: %s", err, rec.Body.String())
	}
}

// AssertErrorResponse asserts an error response
func (h *HTTPTestHelper) AssertErrorResponse(rec *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	h.t.Helper()

	h.AssertStatusCode(rec, expectedStatus)

	var errorResp map[string]interface{}
	h.AssertJSONResponse(rec, &errorResp)

	if msg, ok := errorResp["error"].(string); !ok || msg != expectedMessage {
		h.t.Errorf("expected error message %q, got %q", expectedMessage, msg)
	}
}

// AssertHeader asserts a response header value
func (h *HTTPTestHelper) AssertHeader(rec *httptest.ResponseRecorder, key, expected string) {
	h.t.Helper()

	actual := rec.Header().Get(key)
	if actual != expected {
		h.t.Errorf("expected header %s=%q, got %q", key, expected, actual)
	}
}

// AssertHeaderContains asserts a response header contains a value
func (h *HTTPTestHelper) AssertHeaderContains(rec *httptest.ResponseRecorder, key, contains string) {
	h.t.Helper()

	actual := rec.Header().Get(key)
	if !containsString(actual, contains) {
		h.t.Errorf("expected header %s to contain %q, got %q", key, contains, actual)
	}
}

// AssertBodyContains asserts the response body contains a string
func (h *HTTPTestHelper) AssertBodyContains(rec *httptest.ResponseRecorder, contains string) {
	h.t.Helper()

	body := rec.Body.String()
	if !containsString(body, contains) {
		h.t.Errorf("expected body to contain %q, got:\n%s", contains, body)
	}
}

// WithContext adds values to request context
func (h *HTTPTestHelper) WithContext(req *http.Request, key, value interface{}) *http.Request {
	ctx := context.WithValue(req.Context(), key, value)
	return req.WithContext(ctx)
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

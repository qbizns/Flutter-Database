package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(CodeBadRequest, "test message", http.StatusBadRequest)

	if err.Code != CodeBadRequest {
		t.Errorf("Code = %v, want %v", err.Code, CodeBadRequest)
	}

	if err.Message != "test message" {
		t.Errorf("Message = %v, want %v", err.Message, "test message")
	}

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusBadRequest)
	}

	if err.Err != nil {
		t.Error("Err should be nil")
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := Wrap(originalErr, CodeInternalError, "wrapped message", http.StatusInternalServerError)

	if err.Code != CodeInternalError {
		t.Errorf("Code = %v, want %v", err.Code, CodeInternalError)
	}

	if err.Message != "wrapped message" {
		t.Errorf("Message = %v, want %v", err.Message, "wrapped message")
	}

	if err.Err != originalErr {
		t.Error("Err should be the original error")
	}
}

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *AppError
		want    string
	}{
		{
			name: "without wrapped error",
			err: &AppError{
				Code:    CodeBadRequest,
				Message: "test message",
			},
			want: "BAD_REQUEST: test message",
		},
		{
			name: "with wrapped error",
			err: &AppError{
				Code:    CodeInternalError,
				Message: "test message",
				Err:     errors.New("wrapped error"),
			},
			want: "INTERNAL_ERROR: test message (wrapped error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	originalErr := errors.New("original")
	err := &AppError{
		Code:    CodeInternalError,
		Message: "test",
		Err:     originalErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, originalErr)
	}

	// Test unwrap without error
	err2 := &AppError{
		Code:    CodeBadRequest,
		Message: "test",
	}

	unwrapped2 := err2.Unwrap()
	if unwrapped2 != nil {
		t.Error("Unwrap() should return nil when Err is nil")
	}
}

func TestAppError_WithDetails(t *testing.T) {
	err := New(CodeValidationFailed, "validation error", http.StatusUnprocessableEntity)

	details := map[string]interface{}{
		"field": "email",
		"error": "invalid format",
	}

	err = err.WithDetails(details)

	if err.Details == nil {
		t.Fatal("Details should not be nil")
	}

	if err.Details["field"] != "email" {
		t.Errorf("Details[field] = %v, want %v", err.Details["field"], "email")
	}

	if err.Details["error"] != "invalid format" {
		t.Errorf("Details[error] = %v, want %v", err.Details["error"], "invalid format")
	}
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("invalid input")

	if err.Code != CodeBadRequest {
		t.Errorf("Code = %v, want %v", err.Code, CodeBadRequest)
	}

	if err.Message != "invalid input" {
		t.Errorf("Message = %v, want %v", err.Message, "invalid input")
	}

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusBadRequest)
	}
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("authentication required")

	if err.Code != CodeUnauthorized {
		t.Errorf("Code = %v, want %v", err.Code, CodeUnauthorized)
	}

	if err.Message != "authentication required" {
		t.Errorf("Message = %v, want %v", err.Message, "authentication required")
	}

	if err.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusUnauthorized)
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("access denied")

	if err.Code != CodeForbidden {
		t.Errorf("Code = %v, want %v", err.Code, CodeForbidden)
	}

	if err.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusForbidden)
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("User")

	if err.Code != CodeNotFound {
		t.Errorf("Code = %v, want %v", err.Code, CodeNotFound)
	}

	if err.Message != "User not found" {
		t.Errorf("Message = %v, want %v", err.Message, "User not found")
	}

	if err.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusNotFound)
	}
}

func TestConflict(t *testing.T) {
	err := Conflict("resource conflict")

	if err.Code != CodeConflict {
		t.Errorf("Code = %v, want %v", err.Code, CodeConflict)
	}

	if err.StatusCode != http.StatusConflict {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusConflict)
	}
}

func TestValidationFailed(t *testing.T) {
	err := ValidationFailed("validation error")

	if err.Code != CodeValidationFailed {
		t.Errorf("Code = %v, want %v", err.Code, CodeValidationFailed)
	}

	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestInternalError(t *testing.T) {
	originalErr := errors.New("database connection failed")
	err := InternalError(originalErr)

	if err.Code != CodeInternalError {
		t.Errorf("Code = %v, want %v", err.Code, CodeInternalError)
	}

	if err.Message != "Internal server error" {
		t.Errorf("Message = %v, want %v", err.Message, "Internal server error")
	}

	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusInternalServerError)
	}

	if err.Err != originalErr {
		t.Error("Err should be the original error")
	}
}

func TestDatabaseError(t *testing.T) {
	originalErr := errors.New("query failed")
	err := DatabaseError(originalErr)

	if err.Code != CodeDatabaseError {
		t.Errorf("Code = %v, want %v", err.Code, CodeDatabaseError)
	}

	if err.Message != "Database error occurred" {
		t.Errorf("Message = %v, want %v", err.Message, "Database error occurred")
	}

	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusInternalServerError)
	}

	if err.Err != originalErr {
		t.Error("Err should be the original error")
	}
}

func TestPostingFailed(t *testing.T) {
	err := PostingFailed("posting validation failed")

	if err.Code != CodePostingFailed {
		t.Errorf("Code = %v, want %v", err.Code, CodePostingFailed)
	}

	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestValidationBlocked(t *testing.T) {
	details := map[string]interface{}{
		"reason": "insufficient balance",
		"amount": 100.50,
	}

	err := ValidationBlocked("transaction blocked", details)

	if err.Code != CodeValidationBlocked {
		t.Errorf("Code = %v, want %v", err.Code, CodeValidationBlocked)
	}

	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %v, want %v", err.StatusCode, http.StatusUnprocessableEntity)
	}

	if err.Details == nil {
		t.Fatal("Details should not be nil")
	}

	if err.Details["reason"] != "insufficient balance" {
		t.Error("Details should contain the provided values")
	}
}

func TestAlreadyExists(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		err := AlreadyExists("User", "email already registered")

		if err.Code != CodeConflict {
			t.Errorf("Code = %v, want %v", err.Code, CodeConflict)
		}

		if err.Message != "email already registered" {
			t.Errorf("Message = %v, want %v", err.Message, "email already registered")
		}
	})

	t.Run("with default message", func(t *testing.T) {
		err := AlreadyExists("Product", "")

		if err.Message != "Product already exists" {
			t.Errorf("Message = %v, want %v", err.Message, "Product already exists")
		}
	})
}

func TestIsAppError(t *testing.T) {
	t.Run("with AppError", func(t *testing.T) {
		appErr := BadRequest("test")

		err, ok := IsAppError(appErr)
		if !ok {
			t.Error("IsAppError() should return true for AppError")
		}

		if err != appErr {
			t.Error("IsAppError() should return the same error")
		}
	})

	t.Run("with standard error", func(t *testing.T) {
		stdErr := errors.New("standard error")

		_, ok := IsAppError(stdErr)
		if ok {
			t.Error("IsAppError() should return false for standard error")
		}
	})

	t.Run("with nil error", func(t *testing.T) {
		_, ok := IsAppError(nil)
		if ok {
			t.Error("IsAppError() should return false for nil")
		}
	})
}

func TestIsNotFound(t *testing.T) {
	t.Run("with NotFound error", func(t *testing.T) {
		err := NotFound("User")

		if !IsNotFound(err) {
			t.Error("IsNotFound() should return true for NotFound error")
		}
	})

	t.Run("with different AppError", func(t *testing.T) {
		err := BadRequest("test")

		if IsNotFound(err) {
			t.Error("IsNotFound() should return false for BadRequest")
		}
	})

	t.Run("with standard error", func(t *testing.T) {
		err := errors.New("standard error")

		if IsNotFound(err) {
			t.Error("IsNotFound() should return false for standard error")
		}
	})
}

func TestIsConflict(t *testing.T) {
	t.Run("with Conflict error", func(t *testing.T) {
		err := Conflict("test")

		if !IsConflict(err) {
			t.Error("IsConflict() should return true for Conflict error")
		}
	})

	t.Run("with AlreadyExists error", func(t *testing.T) {
		err := AlreadyExists("User", "")

		if !IsConflict(err) {
			t.Error("IsConflict() should return true for AlreadyExists")
		}
	})

	t.Run("with different AppError", func(t *testing.T) {
		err := BadRequest("test")

		if IsConflict(err) {
			t.Error("IsConflict() should return false for BadRequest")
		}
	})

	t.Run("with standard error", func(t *testing.T) {
		err := errors.New("standard error")

		if IsConflict(err) {
			t.Error("IsConflict() should return false for standard error")
		}
	})
}

func TestErrorChaining(t *testing.T) {
	// Test that errors.Is works with wrapped errors
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, CodeInternalError, "wrapped", http.StatusInternalServerError)

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is should work with wrapped AppError")
	}
}

func TestAllErrorCodes(t *testing.T) {
	codes := []string{
		CodeBadRequest,
		CodeUnauthorized,
		CodeForbidden,
		CodeNotFound,
		CodeConflict,
		CodeValidationFailed,
		CodeInternalError,
		CodeDatabaseError,
		CodePostingFailed,
		CodeValidationBlocked,
	}

	// Ensure all codes are unique
	seen := make(map[string]bool)
	for _, code := range codes {
		if seen[code] {
			t.Errorf("Duplicate error code: %s", code)
		}
		seen[code] = true
	}

	// Ensure codes are not empty
	for _, code := range codes {
		if code == "" {
			t.Error("Error code should not be empty")
		}
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(CodeBadRequest, "test message", http.StatusBadRequest)
	}
}

func BenchmarkBadRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BadRequest("test message")
	}
}

func BenchmarkIsAppError(b *testing.B) {
	err := BadRequest("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IsAppError(err)
	}
}

func BenchmarkErrorString(b *testing.B) {
	err := BadRequest("test message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

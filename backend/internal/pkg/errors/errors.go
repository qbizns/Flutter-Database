package errors

import (
	"fmt"
	"net/http"
)

// Error codes
const (
	CodeBadRequest          = "BAD_REQUEST"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeNotFound            = "NOT_FOUND"
	CodeConflict            = "CONFLICT"
	CodeValidationFailed    = "VALIDATION_FAILED"
	CodeInternalError       = "INTERNAL_ERROR"
	CodeDatabaseError       = "DATABASE_ERROR"
	CodePostingFailed       = "POSTING_FAILED"
	CodeValidationBlocked   = "VALIDATION_BLOCKED"
)

// AppError represents an application error
type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an error with AppError
func Wrap(err error, code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common error constructors

func BadRequest(message string) *AppError {
	return New(CodeBadRequest, message, http.StatusBadRequest)
}

func Unauthorized(message string) *AppError {
	return New(CodeUnauthorized, message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return New(CodeForbidden, message, http.StatusForbidden)
}

func NotFound(resource string) *AppError {
	return New(CodeNotFound, fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

func Conflict(message string) *AppError {
	return New(CodeConflict, message, http.StatusConflict)
}

func ValidationFailed(message string) *AppError {
	return New(CodeValidationFailed, message, http.StatusUnprocessableEntity)
}

func InternalError(err error) *AppError {
	return Wrap(err, CodeInternalError, "Internal server error", http.StatusInternalServerError)
}

func DatabaseError(err error) *AppError {
	return Wrap(err, CodeDatabaseError, "Database error occurred", http.StatusInternalServerError)
}

func PostingFailed(message string) *AppError {
	return New(CodePostingFailed, message, http.StatusUnprocessableEntity)
}

func ValidationBlocked(message string, details map[string]interface{}) *AppError {
	return New(CodeValidationBlocked, message, http.StatusUnprocessableEntity).WithDetails(details)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

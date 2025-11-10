package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Global validator instance
var validate = validator.New()

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// If encoding fails, we can't send a proper error response
			// since headers are already written
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, logger *logging.Logger, err error) {
	// Check if it's an AppError
	if appErr, ok := err.(*apperrors.AppError); ok {
		logger.Warn("request error",
			zap.String("code", appErr.Code),
			zap.String("message", appErr.Message),
			zap.Int("status_code", appErr.StatusCode),
		)

		respondJSON(w, appErr.StatusCode, ErrorResponse{
			Error:   appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		})
		return
	}

	// Generic error
	logger.Error("internal server error", zap.Error(err))
	respondJSON(w, http.StatusInternalServerError, ErrorResponse{
		Error:   "INTERNAL_ERROR",
		Message: "An internal error occurred",
	})
}

// parseUUIDQuery parses a UUID from query parameters
func parseUUIDQuery(r *http.Request, key string) *uuid.UUID {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return nil
	}

	return &id
}

// parseBoolQuery parses a boolean from query parameters
func parseBoolQuery(r *http.Request, key string) *bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}

	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}

	return &boolVal
}

// parseIntQuery parses an integer from query parameters
func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intVal
}

// parseFloatQuery parses a float from query parameters
func parseFloatQuery(r *http.Request, key string) *float64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}

	return &floatVal
}

// getPaginationParams extracts pagination parameters from request
func getPaginationParams(r *http.Request) (page, pageSize int) {
	page = parseIntQuery(r, "page", 1)
	pageSize = parseIntQuery(r, "page_size", 20)

	// Validate bounds
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return page, pageSize
}

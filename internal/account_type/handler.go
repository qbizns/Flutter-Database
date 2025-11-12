package account_type

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/dto/account_type"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/service/account_type"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for AccountTypes
type Handler struct {
	service *account_type.Service
	logger  *logging.Logger
}

// NewHandler creates a new AccountTypes handler
func NewHandler(service *account_type.Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/account_types
// @Summary Create account_types
// @Description Create a new account_types record
// @Tags account_types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body dto.CreateAccountTypesRequest true "AccountTypes data"
// @Success 201 {object} dto.AccountTypesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/account_types [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Parse request body
	var req dto.CreateAccountTypesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, &req)
	
	if err != nil {
		h.logger.Error("failed to create account_types", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create account_types", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/account_types/{id}
// @Summary Get account_types by ID
// @Description Retrieve a account_types by its ID
// @Tags account_types
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "AccountTypes ID"
// @Success 200 {object} dto.AccountTypesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/account_types/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid account_types ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, id)
	
	if err != nil {
		h.logger.Error("failed to get account_types", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "account_types not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/account_types
// @Summary List account_types
// @Description Retrieve a paginated list of account_types records
// @Tags account_types
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.AccountTypesListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/account_types [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Call service
	
	result, err := h.service.List(ctx, page, limit)
	
	if err != nil {
		h.logger.Error("failed to list account_types", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list account_types", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/account_types/{id}
// @Summary Update account_types
// @Description Update an existing account_types record
// @Tags account_types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "AccountTypes ID"
// @Param request body dto.UpdateAccountTypesRequest true "AccountTypes data"
// @Success 200 {object} dto.AccountTypesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/account_types/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid account_types ID", err)
		return
	}

	// Parse request body
	var req dto.UpdateAccountTypesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update account_types", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update account_types", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/account_types/{id}
// @Summary Delete account_types
// @Description Delete a account_types record
// @Tags account_types
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "AccountTypes ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/account_types/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid account_types ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, id)
	
	if err != nil {
		h.logger.Error("failed to delete account_types", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete account_types", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondJSON writes a JSON response
func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.logger.Error("failed to encode response", zap.Error(err))
		}
	}
}

// respondError writes an error response
func (h *Handler) respondError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Error:   message,
		Message: err.Error(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

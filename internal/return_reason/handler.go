package return_reason

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/dto/return_reason"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/service/return_reason"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for ReturnReasons
type Handler struct {
	service *return_reason.Service
	logger  *logging.Logger
}

// NewHandler creates a new ReturnReasons handler
func NewHandler(service *return_reason.Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/return_reasons
// @Summary Create return_reasons
// @Description Create a new return_reasons record
// @Tags return_reasons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body dto.CreateReturnReasonsRequest true "ReturnReasons data"
// @Success 201 {object} dto.ReturnReasonsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/return_reasons [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Parse request body
	var req dto.CreateReturnReasonsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, orgID, &req)
	
	if err != nil {
		h.logger.Error("failed to create return_reasons", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create return_reasons", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/return_reasons/{id}
// @Summary Get return_reasons by ID
// @Description Retrieve a return_reasons by its ID
// @Tags return_reasons
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "ReturnReasons ID"
// @Success 200 {object} dto.ReturnReasonsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/return_reasons/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid return_reasons ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to get return_reasons", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "return_reasons not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/return_reasons
// @Summary List return_reasons
// @Description Retrieve a paginated list of return_reasons records
// @Tags return_reasons
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.ReturnReasonsListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/return_reasons [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

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
	
	result, err := h.service.List(ctx, orgID, page, limit)
	
	if err != nil {
		h.logger.Error("failed to list return_reasons", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list return_reasons", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/return_reasons/{id}
// @Summary Update return_reasons
// @Description Update an existing return_reasons record
// @Tags return_reasons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "ReturnReasons ID"
// @Param request body dto.UpdateReturnReasonsRequest true "ReturnReasons data"
// @Success 200 {object} dto.ReturnReasonsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/return_reasons/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid return_reasons ID", err)
		return
	}

	// Parse request body
	var req dto.UpdateReturnReasonsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, orgID, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update return_reasons", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update return_reasons", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/return_reasons/{id}
// @Summary Delete return_reasons
// @Description Delete a return_reasons record
// @Tags return_reasons
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "ReturnReasons ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/return_reasons/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid return_reasons ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to delete return_reasons", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete return_reasons", err)
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

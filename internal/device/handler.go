package device

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/dto/device"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/service/device"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for Devices
type Handler struct {
	service *device.Service
	logger  *logging.Logger
}

// NewHandler creates a new Devices handler
func NewHandler(service *device.Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/devices
// @Summary Create devices
// @Description Create a new devices record
// @Tags devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body dto.CreateDevicesRequest true "Devices data"
// @Success 201 {object} dto.DevicesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/devices [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Parse request body
	var req dto.CreateDevicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, orgID, &req)
	
	if err != nil {
		h.logger.Error("failed to create devices", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create devices", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/devices/{id}
// @Summary Get devices by ID
// @Description Retrieve a devices by its ID
// @Tags devices
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Devices ID"
// @Success 200 {object} dto.DevicesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/devices/{id} [get]
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
		h.respondError(w, http.StatusBadRequest, "invalid devices ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to get devices", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "devices not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/devices
// @Summary List devices
// @Description Retrieve a paginated list of devices records
// @Tags devices
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.DevicesListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/devices [get]
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
		h.logger.Error("failed to list devices", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list devices", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/devices/{id}
// @Summary Update devices
// @Description Update an existing devices record
// @Tags devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Devices ID"
// @Param request body dto.UpdateDevicesRequest true "Devices data"
// @Success 200 {object} dto.DevicesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/devices/{id} [put]
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
		h.respondError(w, http.StatusBadRequest, "invalid devices ID", err)
		return
	}

	// Parse request body
	var req dto.UpdateDevicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, orgID, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update devices", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update devices", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/devices/{id}
// @Summary Delete devices
// @Description Delete a devices record
// @Tags devices
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Devices ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/devices/{id} [delete]
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
		h.respondError(w, http.StatusBadRequest, "invalid devices ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to delete devices", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete devices", err)
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

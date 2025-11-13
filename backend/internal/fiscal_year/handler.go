package fiscal_year

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for FiscalYears
type Handler struct {
	service *Service
	logger  *logging.Logger
}

// NewHandler creates a new FiscalYears handler
func NewHandler(service *Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/fiscal_years
// @Summary Create fiscal_years
// @Description Create a new fiscal_years record
// @Tags fiscal_years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body CreateFiscalYearsRequest true "FiscalYears data"
// @Success 201 {object} FiscalYearsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/fiscal_years [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Parse request body
	var req CreateFiscalYearsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, orgID, &req)
	
	if err != nil {
		h.logger.Error("failed to create fiscal_years", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create fiscal_years", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/fiscal_years/{id}
// @Summary Get fiscal_years by ID
// @Description Retrieve a fiscal_years by its ID
// @Tags fiscal_years
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "FiscalYears ID"
// @Success 200 {object} FiscalYearsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/fiscal_years/{id} [get]
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
		h.respondError(w, http.StatusBadRequest, "invalid fiscal_years ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to get fiscal_years", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "fiscal_years not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/fiscal_years
// @Summary List fiscal_years
// @Description Retrieve a paginated list of fiscal_years records
// @Tags fiscal_years
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} FiscalYearsListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/fiscal_years [get]
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
		h.logger.Error("failed to list fiscal_years", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list fiscal_years", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/fiscal_years/{id}
// @Summary Update fiscal_years
// @Description Update an existing fiscal_years record
// @Tags fiscal_years
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "FiscalYears ID"
// @Param request body UpdateFiscalYearsRequest true "FiscalYears data"
// @Success 200 {object} FiscalYearsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/fiscal_years/{id} [put]
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
		h.respondError(w, http.StatusBadRequest, "invalid fiscal_years ID", err)
		return
	}

	// Parse request body
	var req UpdateFiscalYearsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, orgID, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update fiscal_years", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update fiscal_years", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/fiscal_years/{id}
// @Summary Delete fiscal_years
// @Description Delete a fiscal_years record
// @Tags fiscal_years
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "FiscalYears ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/fiscal_years/{id} [delete]
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
		h.respondError(w, http.StatusBadRequest, "invalid fiscal_years ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to delete fiscal_years", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete fiscal_years", err)
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

// Admin handler methods

// AdminList handles GET /api/v1/admin/{module}/
func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	// Reuse standard List with potential admin-only filters
	h.List(w, r)
}

// GetStats handles GET /api/v1/admin/{module}/stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement statistics gathering
	stats := map[string]interface{}{
		"total":   0,
		"active":  0,
		"deleted": 0,
	}
	h.respondJSON(w, http.StatusOK, stats)
}

// Export handles POST /api/v1/admin/{module}/export
func (h *Handler) Export(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement data export (CSV, JSON, Excel)
	h.respondError(w, http.StatusNotImplemented, "export not yet implemented", nil)
}

// Import handles POST /api/v1/admin/{module}/import
func (h *Handler) Import(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement data import with validation
	h.respondError(w, http.StatusNotImplemented, "import not yet implemented", nil)
}

// ListDeleted handles GET /api/v1/admin/{module}/deleted
func (h *Handler) ListDeleted(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement listing soft-deleted records
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"items": []interface{}{},
		"pagination": map[string]int{
			"page":  1,
			"limit": 20,
			"total": 0,
		},
	})
}

// Restore handles POST /api/v1/admin/{module}/{id}/restore
func (h *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	// TODO: Implement restore functionality
	h.respondError(w, http.StatusNotImplemented, "restore not yet implemented", nil)
}

// PermanentDelete handles DELETE /api/v1/admin/{module}/{id}/permanent
func (h *Handler) PermanentDelete(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	// TODO: Implement permanent deletion (hard delete)
	// WARNING: This cannot be undone!
	h.respondError(w, http.StatusNotImplemented, "permanent delete not yet implemented", nil)
}

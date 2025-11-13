package payment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for Payments
type Handler struct {
	service *Service
	logger  *logging.Logger
}

// NewHandler creates a new Payments handler
func NewHandler(service *Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/payments
// @Summary Create payments
// @Description Create a new payments record
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body CreatePaymentsRequest true "Payments data"
// @Success 201 {object} PaymentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/payments [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Parse request body
	var req CreatePaymentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, orgID, &req)
	
	if err != nil {
		h.logger.Error("failed to create payments", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create payments", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/payments/{id}
// @Summary Get payments by ID
// @Description Retrieve a payments by its ID
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Payments ID"
// @Success 200 {object} PaymentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/payments/{id} [get]
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
		h.respondError(w, http.StatusBadRequest, "invalid payments ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to get payments", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "payments not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/payments
// @Summary List payments
// @Description Retrieve a paginated list of payments records
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} PaymentsListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/payments [get]
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
		h.logger.Error("failed to list payments", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list payments", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/payments/{id}
// @Summary Update payments
// @Description Update an existing payments record
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Payments ID"
// @Param request body UpdatePaymentsRequest true "Payments data"
// @Success 200 {object} PaymentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/payments/{id} [put]
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
		h.respondError(w, http.StatusBadRequest, "invalid payments ID", err)
		return
	}

	// Parse request body
	var req UpdatePaymentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, orgID, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update payments", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update payments", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/payments/{id}
// @Summary Delete payments
// @Description Delete a payments record
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Payments ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/payments/{id} [delete]
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
		h.respondError(w, http.StatusBadRequest, "invalid payments ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to delete payments", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete payments", err)
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

// Cancel handles POST /api/v1/organizations/{orgID}/payments/{id}/cancel
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid payment ID", err)
		return
	}

	var req CancelPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	result, err := h.service.Cancel(ctx, orgID, id, &req)
	if err != nil {
		h.logger.Error("failed to cancel payment", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to cancel payment", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// GetStatistics handles GET /api/v1/organizations/{orgID}/payments/statistics
func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	fromDate := r.URL.Query().Get("from_date")
	toDate := r.URL.Query().Get("to_date")

	result, err := h.service.GetStatistics(ctx, orgID, fromDate, toDate)
	if err != nil {
		h.logger.Error("failed to get payment statistics", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to get payment statistics", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// CreateRefund handles POST /api/v1/organizations/{orgID}/refunds
func (h *Handler) CreateRefund(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	var req CreateRefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	result, err := h.service.CreateRefund(ctx, orgID, &req)
	if err != nil {
		h.logger.Error("failed to create refund", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create refund", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// ListRefunds handles GET /api/v1/organizations/{orgID}/refunds
func (h *Handler) ListRefunds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	paymentID := r.URL.Query().Get("payment_id")
	orderID := r.URL.Query().Get("order_id")

	result, err := h.service.ListRefunds(ctx, orgID, paymentID, orderID)
	if err != nil {
		h.logger.Error("failed to list refunds", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list refunds", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// GetRefund handles GET /api/v1/organizations/{orgID}/refunds/{id}
func (h *Handler) GetRefund(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid refund ID", err)
		return
	}

	result, err := h.service.GetRefund(ctx, orgID, id)
	if err != nil {
		h.logger.Error("failed to get refund", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "refund not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

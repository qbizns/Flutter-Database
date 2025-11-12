package gift_card_transaction

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/dto/gift_card_transaction"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/service/gift_card_transaction"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for GiftCardTransactions
type Handler struct {
	service *gift_card_transaction.Service
	logger  *logging.Logger
}

// NewHandler creates a new GiftCardTransactions handler
func NewHandler(service *gift_card_transaction.Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/gift_card_transactions
// @Summary Create gift_card_transactions
// @Description Create a new gift_card_transactions record
// @Tags gift_card_transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body dto.CreateGiftCardTransactionsRequest true "GiftCardTransactions data"
// @Success 201 {object} dto.GiftCardTransactionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/gift_card_transactions [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Parse request body
	var req dto.CreateGiftCardTransactionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, orgID, &req)
	
	if err != nil {
		h.logger.Error("failed to create gift_card_transactions", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create gift_card_transactions", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/gift_card_transactions/{id}
// @Summary Get gift_card_transactions by ID
// @Description Retrieve a gift_card_transactions by its ID
// @Tags gift_card_transactions
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "GiftCardTransactions ID"
// @Success 200 {object} dto.GiftCardTransactionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/gift_card_transactions/{id} [get]
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
		h.respondError(w, http.StatusBadRequest, "invalid gift_card_transactions ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to get gift_card_transactions", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "gift_card_transactions not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/gift_card_transactions
// @Summary List gift_card_transactions
// @Description Retrieve a paginated list of gift_card_transactions records
// @Tags gift_card_transactions
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.GiftCardTransactionsListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/gift_card_transactions [get]
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
		h.logger.Error("failed to list gift_card_transactions", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list gift_card_transactions", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/gift_card_transactions/{id}
// @Summary Update gift_card_transactions
// @Description Update an existing gift_card_transactions record
// @Tags gift_card_transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "GiftCardTransactions ID"
// @Param request body dto.UpdateGiftCardTransactionsRequest true "GiftCardTransactions data"
// @Success 200 {object} dto.GiftCardTransactionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/gift_card_transactions/{id} [put]
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
		h.respondError(w, http.StatusBadRequest, "invalid gift_card_transactions ID", err)
		return
	}

	// Parse request body
	var req dto.UpdateGiftCardTransactionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, orgID, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update gift_card_transactions", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update gift_card_transactions", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/gift_card_transactions/{id}
// @Summary Delete gift_card_transactions
// @Description Delete a gift_card_transactions record
// @Tags gift_card_transactions
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "GiftCardTransactions ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/gift_card_transactions/{id} [delete]
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
		h.respondError(w, http.StatusBadRequest, "invalid gift_card_transactions ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to delete gift_card_transactions", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete gift_card_transactions", err)
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

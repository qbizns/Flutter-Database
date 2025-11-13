package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for Products
type Handler struct {
	service *Service
	logger  *logging.Logger
}

// NewHandler creates a new Products handler
func NewHandler(service *Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/products
// @Summary Create products
// @Description Create a new products record
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body CreateProductsRequest true "Products data"
// @Success 201 {object} ProductsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/products [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	
	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "orgID"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}
	

	// Parse request body
	var req CreateProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, orgID, &req)
	
	if err != nil {
		h.logger.Error("failed to create products", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create products", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/products/{id}
// @Summary Get products by ID
// @Description Retrieve a products by its ID
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Products ID"
// @Success 200 {object} ProductsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/products/{id} [get]
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
		h.respondError(w, http.StatusBadRequest, "invalid products ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to get products", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "products not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/products
// @Summary List products
// @Description Retrieve a paginated list of products records
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} ProductsListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/products [get]
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
		h.logger.Error("failed to list products", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list products", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/products/{id}
// @Summary Update products
// @Description Update an existing products record
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Products ID"
// @Param request body UpdateProductsRequest true "Products data"
// @Success 200 {object} ProductsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/products/{id} [put]
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
		h.respondError(w, http.StatusBadRequest, "invalid products ID", err)
		return
	}

	// Parse request body
	var req UpdateProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, orgID, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update products", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update products", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/products/{id}
// @Summary Delete products
// @Description Delete a products record
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "Products ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/products/{id} [delete]
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
		h.respondError(w, http.StatusBadRequest, "invalid products ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, orgID, id)
	
	if err != nil {
		h.logger.Error("failed to delete products", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete products", err)
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

// Search handles GET /api/v1/organizations/{orgID}/products/search
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	// Get query parameters
	query := r.URL.Query().Get("q")
	categoryID := r.URL.Query().Get("category_id")

	// Call service
	result, err := h.service.Search(ctx, orgID, query, categoryID)
	if err != nil {
		h.logger.Error("failed to search products", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to search products", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// GetBatch handles POST /api/v1/organizations/{orgID}/products/batch
func (h *Handler) GetBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	// Parse request body
	var req GetBatchProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	result, err := h.service.GetBatch(ctx, orgID, &req)
	if err != nil {
		h.logger.Error("failed to get batch products", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to get batch products", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// GetFeatured handles GET /api/v1/organizations/{orgID}/products/featured
func (h *Handler) GetFeatured(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	// Get query parameter
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Call service
	result, err := h.service.GetFeatured(ctx, orgID, limit)
	if err != nil {
		h.logger.Error("failed to get featured products", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to get featured products", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// GetLowStock handles GET /api/v1/organizations/{orgID}/products/low-stock
func (h *Handler) GetLowStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	// Call service
	result, err := h.service.GetLowStock(ctx, orgID)
	if err != nil {
		h.logger.Error("failed to get low stock products", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to get low stock products", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// UpdateStock handles PATCH /api/v1/organizations/{orgID}/products/{id}/stock
func (h *Handler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	// Parse request body
	var req UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	result, err := h.service.UpdateStock(ctx, orgID, id, &req)
	if err != nil {
		h.logger.Error("failed to update product stock", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update product stock", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// UpdateAvailability handles PATCH /api/v1/organizations/{orgID}/products/{id}/availability
func (h *Handler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get organization ID from URL
	orgID, err := uuid.Parse(chi.URLParam(r, "org_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid organization ID", err)
		return
	}

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	// Parse request body
	var req UpdateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	result, err := h.service.UpdateAvailability(ctx, orgID, id, &req)
	if err != nil {
		h.logger.Error("failed to update product availability", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update product availability", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

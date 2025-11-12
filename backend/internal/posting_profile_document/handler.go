package posting_profile_document

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/dto/posting_profile_document"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/service/posting_profile_document"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for PostingProfileDocuments
type Handler struct {
	service *posting_profile_document.Service
	logger  *logging.Logger
}

// NewHandler creates a new PostingProfileDocuments handler
func NewHandler(service *posting_profile_document.Service, logger *logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /api/v1/organizations/{orgID}/posting_profile_documents
// @Summary Create posting_profile_documents
// @Description Create a new posting_profile_documents record
// @Tags posting_profile_documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param request body dto.CreatePostingProfileDocumentsRequest true "PostingProfileDocuments data"
// @Success 201 {object} dto.PostingProfileDocumentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/posting_profile_documents [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Parse request body
	var req dto.CreatePostingProfileDocumentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Create(ctx, &req)
	
	if err != nil {
		h.logger.Error("failed to create posting_profile_documents", zap.Error(err))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to create posting_profile_documents", err)
		return
	}

	h.respondJSON(w, http.StatusCreated, result)
}

// GetByID handles GET /api/v1/organizations/{orgID}/posting_profile_documents/{id}
// @Summary Get posting_profile_documents by ID
// @Description Retrieve a posting_profile_documents by its ID
// @Tags posting_profile_documents
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "PostingProfileDocuments ID"
// @Success 200 {object} dto.PostingProfileDocumentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/posting_profile_documents/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid posting_profile_documents ID", err)
		return
	}

	// Call service
	
	result, err := h.service.GetByID(ctx, id)
	
	if err != nil {
		h.logger.Error("failed to get posting_profile_documents", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusNotFound, "posting_profile_documents not found", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// List handles GET /api/v1/organizations/{orgID}/posting_profile_documents
// @Summary List posting_profile_documents
// @Description Retrieve a paginated list of posting_profile_documents records
// @Tags posting_profile_documents
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.PostingProfileDocumentsListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/posting_profile_documents [get]
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
		h.logger.Error("failed to list posting_profile_documents", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "failed to list posting_profile_documents", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Update handles PUT /api/v1/organizations/{orgID}/posting_profile_documents/{id}
// @Summary Update posting_profile_documents
// @Description Update an existing posting_profile_documents record
// @Tags posting_profile_documents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "PostingProfileDocuments ID"
// @Param request body dto.UpdatePostingProfileDocumentsRequest true "PostingProfileDocuments data"
// @Success 200 {object} dto.PostingProfileDocumentsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/posting_profile_documents/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid posting_profile_documents ID", err)
		return
	}

	// Parse request body
	var req dto.UpdatePostingProfileDocumentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Call service
	
	result, err := h.service.Update(ctx, id, &req)
	
	if err != nil {
		h.logger.Error("failed to update posting_profile_documents", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to update posting_profile_documents", err)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/organizations/{orgID}/posting_profile_documents/{id}
// @Summary Delete posting_profile_documents
// @Description Delete a posting_profile_documents record
// @Tags posting_profile_documents
// @Produce json
// @Security BearerAuth
// @Param orgID path string true "Organization ID"
// @Param id path string true "PostingProfileDocuments ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /organizations/{orgID}/posting_profile_documents/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	

	// Get ID from URL
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid posting_profile_documents ID", err)
		return
	}

	// Call service
	
	err = h.service.Delete(ctx, id)
	
	if err != nil {
		h.logger.Error("failed to delete posting_profile_documents", zap.Error(err), zap.String("id", id.String()))
		h.respondError(w, http.StatusUnprocessableEntity, "failed to delete posting_profile_documents", err)
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

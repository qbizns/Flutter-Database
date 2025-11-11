package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/posting"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

// PostDocumentHandler handles POST /api/v1/organizations/{org_id}/posting/post
func PostDocumentHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		userID, err := getUserID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse request body
		var req PostDocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and engine
		repo := postgres.NewPostingRepository(db)
		engine := posting.NewEngine(repo, logger)

		// Build posting input
		input := posting.PostingInput{
			OrganizationID: orgID,
			DocumentType:   req.DocumentType,
			DocumentID:     req.DocumentID,
			Event:          req.Event,
			UserID:         userID,
		}

		// Post document
		if err := engine.Post(ctx, input); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("document posted successfully",
			zap.String("document_type", req.DocumentType),
			zap.String("document_id", req.DocumentID.String()),
		)

		// Return success response
		response := PostDocumentResponse{
			JournalEntryID: uuid.Nil, // Will be populated by engine in future
			Status:         "posted",
			Message:        "Document posted successfully",
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// GetPostingRulesHandler handles GET /api/v1/organizations/{org_id}/posting/rules
func GetPostingRulesHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse query parameters
		documentType := r.URL.Query().Get("document_type")
		event := r.URL.Query().Get("event")

		// Create repository
		repo := postgres.NewPostingRepository(db)

		// Get posting rules
		rules, err := repo.GetPostingRules(ctx, orgID, documentType, event)
		if err != nil {
			respondError(w, logger, apperrors.DatabaseError(err))
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"rules": rules,
			"count": len(rules),
		})
	}
}

// GetPostingAuditHandler handles GET /api/v1/organizations/{org_id}/posting/audit
func GetPostingAuditHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse query parameters
		documentType := r.URL.Query().Get("document_type")
		documentID := parseUUIDQuery(r, "document_id")

		// Get pagination params
		page, pageSize := getPaginationParams(r)

		// Create repository
		repo := postgres.NewPostingRepository(db)

		// Build filter
		filter := make(map[string]interface{})
		if documentType != "" {
			filter["source_table"] = documentType
		}
		if documentID != nil {
			filter["source_id"] = documentID
		}
		filter["limit"] = pageSize
		filter["offset"] = (page - 1) * pageSize

		// Get audit logs
		auditLogs, err := repo.GetPostingAuditLogs(ctx, orgID, filter)
		if err != nil {
			respondError(w, logger, apperrors.DatabaseError(err))
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"audit_logs": auditLogs,
			"count":      len(auditLogs),
			"page":       page,
			"page_size":  pageSize,
		})
	}
}

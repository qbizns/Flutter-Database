package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/accounting"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

// ListJournalEntriesHandler handles GET /api/v1/organizations/{org_id}/journal-entries
func ListJournalEntriesHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Parse query parameters
		filter := make(map[string]interface{})
		if status := r.URL.Query().Get("status"); status != "" {
			filter["status"] = status
		}
		if fiscalYearID := parseUUIDQuery(r, "fiscal_year_id"); fiscalYearID != nil {
			filter["fiscal_year_id"] = fiscalYearID
		}
		if periodID := parseUUIDQuery(r, "accounting_period_id"); periodID != nil {
			filter["accounting_period_id"] = periodID
		}
		if entryTypeID := parseUUIDQuery(r, "entry_type_id"); entryTypeID != nil {
			filter["entry_type_id"] = entryTypeID
		}
		if sourceModule := r.URL.Query().Get("source_module"); sourceModule != "" {
			filter["source_module"] = sourceModule
		}
		if isPosted := parseBoolQuery(r, "is_posted"); isPosted != nil {
			filter["is_posted"] = *isPosted
		}

		// Get pagination params
		page, pageSize := getPaginationParams(r)
		filter["limit"] = pageSize
		filter["offset"] = (page - 1) * pageSize

		// Get journal entries
		entries, err := service.ListJournalEntries(ctx, orgID, filter)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"journal_entries": entries,
			"count":           len(entries),
			"page":            page,
			"page_size":       pageSize,
		})
	}
}

// CreateJournalEntryHandler handles POST /api/v1/organizations/{org_id}/journal-entries
func CreateJournalEntryHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req accounting.CreateJournalEntryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Create journal entry
		entry, err := service.CreateJournalEntry(ctx, orgID, &req, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("journal entry created",
			zap.String("entry_id", entry.ID.String()),
			zap.String("entry_number", entry.EntryNumber),
		)

		respondJSON(w, http.StatusCreated, entry)
	}
}

// GetJournalEntryHandler handles GET /api/v1/organizations/{org_id}/journal-entries/{id}
func GetJournalEntryHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse entry ID from URL
		entryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid journal entry ID"))
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Get journal entry
		entry, err := service.GetJournalEntry(ctx, entryID, orgID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Get lines
		lines, err := repo.ListJournalEntryLines(ctx, entryID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"journal_entry": entry,
			"lines":         lines,
		})
	}
}

// PostJournalEntryHandler handles POST /api/v1/organizations/{org_id}/journal-entries/{id}/post
func PostJournalEntryHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Parse entry ID from URL
		entryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid journal entry ID"))
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Post journal entry
		if err := service.PostJournalEntry(ctx, entryID, orgID, userID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("journal entry posted",
			zap.String("entry_id", entryID.String()),
			zap.String("posted_by", userID.String()),
		)

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Journal entry posted successfully",
			"status":  "posted",
		})
	}
}

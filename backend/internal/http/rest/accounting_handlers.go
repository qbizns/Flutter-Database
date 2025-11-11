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

// ============================================================================
// CHART OF ACCOUNTS HANDLERS
// ============================================================================

// ListChartOfAccountsHandler handles GET /api/v1/organizations/{org_id}/chart-of-accounts
func ListChartOfAccountsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		if accountTypeID := parseUUIDQuery(r, "account_type_id"); accountTypeID != nil {
			filter["account_type_id"] = accountTypeID
		}
		if isActive := parseBoolQuery(r, "is_active"); isActive != nil {
			filter["is_active"] = *isActive
		}
		if isHeaderAccount := parseBoolQuery(r, "is_header_account"); isHeaderAccount != nil {
			filter["is_header_account"] = *isHeaderAccount
		}
		if search := r.URL.Query().Get("search"); search != "" {
			filter["search"] = search
		}

		// Get accounts
		accounts, err := service.ListChartOfAccounts(ctx, orgID, filter)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"accounts": accounts,
			"count":    len(accounts),
		})
	}
}

// CreateChartOfAccountHandler handles POST /api/v1/organizations/{org_id}/chart-of-accounts
func CreateChartOfAccountHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req accounting.CreateChartOfAccountRequest
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

		// Create account
		account, err := service.CreateChartOfAccount(ctx, orgID, &req, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("chart of account created",
			zap.String("account_id", account.ID.String()),
			zap.String("account_code", account.AccountCode),
		)

		respondJSON(w, http.StatusCreated, account)
	}
}

// GetChartOfAccountHandler handles GET /api/v1/organizations/{org_id}/chart-of-accounts/{id}
func GetChartOfAccountHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse account ID from URL
		accountID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid account ID"))
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Get account
		account, err := service.GetChartOfAccount(ctx, accountID, orgID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, account)
	}
}

// UpdateChartOfAccountHandler handles PATCH /api/v1/organizations/{org_id}/chart-of-accounts/{id}
func UpdateChartOfAccountHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Parse account ID from URL
		accountID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid account ID"))
			return
		}

		// Parse request body
		var req accounting.UpdateChartOfAccountRequest
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

		// Update account
		account, err := service.UpdateChartOfAccount(ctx, accountID, orgID, &req, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("chart of account updated",
			zap.String("account_id", account.ID.String()),
			zap.String("account_code", account.AccountCode),
		)

		respondJSON(w, http.StatusOK, account)
	}
}

// DeleteChartOfAccountHandler handles DELETE /api/v1/organizations/{org_id}/chart-of-accounts/{id}
func DeleteChartOfAccountHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse account ID from URL
		accountID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid account ID"))
			return
		}

		// Create repository and service
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())

		// Delete account
		if err := repo.DeleteChartOfAccount(ctx, accountID, orgID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("chart of account deleted",
			zap.String("account_id", accountID.String()),
		)

		w.WriteHeader(http.StatusNoContent)
	}
}

// ============================================================================
// FISCAL YEAR HANDLERS
// ============================================================================

// ListFiscalYearsHandler handles GET /api/v1/organizations/{org_id}/fiscal-years
func ListFiscalYearsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		if isCurrent := parseBoolQuery(r, "is_current"); isCurrent != nil {
			filter["is_current"] = *isCurrent
		}

		// Get fiscal years
		fiscalYears, err := service.ListFiscalYears(ctx, orgID, filter)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"fiscal_years": fiscalYears,
			"count":        len(fiscalYears),
		})
	}
}

// CreateFiscalYearHandler handles POST /api/v1/organizations/{org_id}/fiscal-years
func CreateFiscalYearHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req accounting.CreateFiscalYearRequest
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

		// Create fiscal year
		fiscalYear, err := service.CreateFiscalYear(ctx, orgID, &req, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("fiscal year created",
			zap.String("fiscal_year_id", fiscalYear.ID.String()),
			zap.String("fiscal_year", fiscalYear.FiscalYear),
		)

		respondJSON(w, http.StatusCreated, fiscalYear)
	}
}

// GetFiscalYearHandler handles GET /api/v1/organizations/{org_id}/fiscal-years/{id}
func GetFiscalYearHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse fiscal year ID from URL
		fiscalYearID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid fiscal year ID"))
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Get fiscal year
		fiscalYear, err := service.GetFiscalYear(ctx, fiscalYearID, orgID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, fiscalYear)
	}
}

// UpdateFiscalYearHandler handles PATCH /api/v1/organizations/{org_id}/fiscal-years/{id}
func UpdateFiscalYearHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Parse fiscal year ID from URL
		fiscalYearID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid fiscal year ID"))
			return
		}

		// Parse request body
		var req accounting.UpdateFiscalYearRequest
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

		// Update fiscal year
		fiscalYear, err := service.UpdateFiscalYear(ctx, fiscalYearID, orgID, &req, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("fiscal year updated",
			zap.String("fiscal_year_id", fiscalYear.ID.String()),
			zap.String("fiscal_year", fiscalYear.FiscalYear),
		)

		respondJSON(w, http.StatusOK, fiscalYear)
	}
}

// CloseFiscalYearHandler handles POST /api/v1/organizations/{org_id}/fiscal-years/{id}/close
func CloseFiscalYearHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Parse fiscal year ID from URL
		fiscalYearID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid fiscal year ID"))
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Close fiscal year
		if err := service.CloseFiscalYear(ctx, fiscalYearID, orgID, userID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("fiscal year closed",
			zap.String("fiscal_year_id", fiscalYearID.String()),
			zap.String("closed_by", userID.String()),
		)

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Fiscal year closed successfully",
			"status":  "closed",
		})
	}
}

// ============================================================================
// ACCOUNTING PERIOD HANDLERS
// ============================================================================

// ListAccountingPeriodsHandler handles GET /api/v1/organizations/{org_id}/accounting-periods
func ListAccountingPeriodsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse fiscal year ID
		fiscalYearID := parseUUIDQuery(r, "fiscal_year_id")
		if fiscalYearID == nil {
			respondError(w, logger, apperrors.BadRequest("fiscal_year_id is required"))
			return
		}

		// Create repository and service using sqlx compatibility layer
		repo := postgres.NewAccountingRepository(db.GetSQLXDB())
		service := accounting.NewService(repo)

		// Get accounting periods
		periods, err := service.ListAccountingPeriods(ctx, orgID, *fiscalYearID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"accounting_periods": periods,
			"count":              len(periods),
		})
	}
}

// CreateAccountingPeriodHandler handles POST /api/v1/organizations/{org_id}/accounting-periods
func CreateAccountingPeriodHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req accounting.CreateAccountingPeriodRequest
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

		// Create accounting period
		period, err := service.CreateAccountingPeriod(ctx, orgID, &req, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("accounting period created",
			zap.String("period_id", period.ID.String()),
			zap.String("period_name", period.PeriodName),
		)

		respondJSON(w, http.StatusCreated, period)
	}
}

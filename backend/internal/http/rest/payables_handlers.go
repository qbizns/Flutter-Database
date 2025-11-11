package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/payables"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

// ============================================================================
// VENDOR BILL HANDLERS
// ============================================================================

// ListVendorBillsHandler handles GET /api/v1/organizations/{org_id}/vendor-bills
func ListVendorBillsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Parse query parameters
		opts := &payables.VendorBillFilterOptions{}
		if supplierID := parseUUIDQuery(r, "supplier_id"); supplierID != nil {
			opts.SupplierID = supplierID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			opts.Status = (*payables.VendorBillStatus)(&status)
		}
		if isPosted := parseBoolQuery(r, "is_posted"); isPosted != nil {
			opts.IsPosted = isPosted
		}

		// Get pagination params
		page, pageSize := getPaginationParams(r)
		opts.Limit = pageSize
		opts.Offset = (page - 1) * pageSize

		// Get bills
		bills, total, err := service.ListVendorBills(ctx, orgID, opts)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"bills":       bills,
			"total_count": total,
			"page":        page,
			"page_size":   pageSize,
		})
	}
}

// CreateVendorBillHandler handles POST /api/v1/organizations/{org_id}/vendor-bills
func CreateVendorBillHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req payables.CreateVendorBillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Create bill
		bill, err := service.CreateVendorBill(ctx, &req, orgID, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("vendor bill created",
			zap.String("bill_id", bill.Bill.ID.String()),
			zap.String("bill_number", bill.Bill.BillNumber),
		)

		respondJSON(w, http.StatusCreated, bill)
	}
}

// GetVendorBillHandler handles GET /api/v1/organizations/{org_id}/vendor-bills/{id}
func GetVendorBillHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse bill ID from URL
		billID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid bill ID"))
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Get bill
		bill, err := service.GetVendorBill(ctx, orgID, billID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, bill)
	}
}

// UpdateVendorBillHandler handles PATCH /api/v1/organizations/{org_id}/vendor-bills/{id}
func UpdateVendorBillHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Parse bill ID from URL
		billID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid bill ID"))
			return
		}

		// Parse request body
		var req payables.UpdateVendorBillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Update bill
		bill, err := service.UpdateVendorBill(ctx, billID, &req, orgID, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("vendor bill updated",
			zap.String("bill_id", billID.String()),
		)

		respondJSON(w, http.StatusOK, bill)
	}
}

// DeleteVendorBillHandler handles DELETE /api/v1/organizations/{org_id}/vendor-bills/{id}
func DeleteVendorBillHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse bill ID from URL
		billID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid bill ID"))
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Delete bill
		if err := service.DeleteVendorBill(ctx, orgID, billID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("vendor bill deleted",
			zap.String("bill_id", billID.String()),
		)

		w.WriteHeader(http.StatusNoContent)
	}
}

// ============================================================================
// VENDOR PAYMENT HANDLERS
// ============================================================================

// ListVendorPaymentsHandler handles GET /api/v1/organizations/{org_id}/vendor-payments
func ListVendorPaymentsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Parse query parameters
		opts := &payables.VendorPaymentFilterOptions{}
		if supplierID := parseUUIDQuery(r, "supplier_id"); supplierID != nil {
			opts.SupplierID = supplierID
		}
		if isPosted := parseBoolQuery(r, "is_posted"); isPosted != nil {
			opts.IsPosted = isPosted
		}

		// Get pagination params
		page, pageSize := getPaginationParams(r)
		opts.Limit = pageSize
		opts.Offset = (page - 1) * pageSize

		// Get payments
		payments, total, err := service.ListVendorPayments(ctx, orgID, opts)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"payments":    payments,
			"total_count": total,
			"page":        page,
			"page_size":   pageSize,
		})
	}
}

// CreateVendorPaymentHandler handles POST /api/v1/organizations/{org_id}/vendor-payments
func CreateVendorPaymentHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req payables.CreateVendorPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Create payment
		payment, err := service.CreateVendorPayment(ctx, &req, orgID, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("vendor payment created",
			zap.String("payment_id", payment.Payment.ID.String()),
			zap.String("payment_number", payment.Payment.PaymentNumber),
		)

		respondJSON(w, http.StatusCreated, payment)
	}
}

// GetVendorPaymentHandler handles GET /api/v1/organizations/{org_id}/vendor-payments/{id}
func GetVendorPaymentHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse payment ID from URL
		paymentID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid payment ID"))
			return
		}

		// Create repository and service
		repo := postgres.NewPayablesRepository(db)
		service := payables.NewService(repo)

		// Get payment
		payment, err := service.GetVendorPayment(ctx, orgID, paymentID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, payment)
	}
}

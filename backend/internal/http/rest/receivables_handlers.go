package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/receivables"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

// ============================================================================
// CUSTOMER INVOICE HANDLERS
// ============================================================================

// ListCustomerInvoicesHandler handles GET /api/v1/organizations/{org_id}/customer-invoices
func ListCustomerInvoicesHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Parse query parameters
		opts := &receivables.CustomerInvoiceFilterOptions{}
		if customerID := parseUUIDQuery(r, "customer_id"); customerID != nil {
			opts.CustomerID = customerID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			opts.Status = (*receivables.InvoiceStatus)(&status)
		}
		if isPosted := parseBoolQuery(r, "is_posted"); isPosted != nil {
			opts.IsPosted = isPosted
		}

		// Get pagination params
		page, pageSize := getPaginationParams(r)
		opts.Limit = pageSize
		opts.Offset = (page - 1) * pageSize

		// Get invoices
		invoices, total, err := service.ListCustomerInvoices(ctx, orgID, opts)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"invoices":    invoices,
			"total_count": total,
			"page":        page,
			"page_size":   pageSize,
		})
	}
}

// CreateCustomerInvoiceHandler handles POST /api/v1/organizations/{org_id}/customer-invoices
func CreateCustomerInvoiceHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req receivables.CreateCustomerInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Create invoice
		invoice, err := service.CreateCustomerInvoice(ctx, &req, orgID, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("customer invoice created",
			zap.String("invoice_id", invoice.Invoice.ID.String()),
			zap.String("invoice_number", invoice.Invoice.InvoiceNumber),
		)

		respondJSON(w, http.StatusCreated, invoice)
	}
}

// GetCustomerInvoiceHandler handles GET /api/v1/organizations/{org_id}/customer-invoices/{id}
func GetCustomerInvoiceHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse invoice ID from URL
		invoiceID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid invoice ID"))
			return
		}

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Get invoice
		invoice, err := service.GetCustomerInvoice(ctx, orgID, invoiceID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, invoice)
	}
}

// UpdateCustomerInvoiceHandler handles PATCH /api/v1/organizations/{org_id}/customer-invoices/{id}
func UpdateCustomerInvoiceHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Parse invoice ID from URL
		invoiceID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid invoice ID"))
			return
		}

		// Parse request body
		var req receivables.UpdateCustomerInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Update invoice
		invoice, err := service.UpdateCustomerInvoice(ctx, invoiceID, &req, orgID, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("customer invoice updated",
			zap.String("invoice_id", invoiceID.String()),
		)

		respondJSON(w, http.StatusOK, invoice)
	}
}

// DeleteCustomerInvoiceHandler handles DELETE /api/v1/organizations/{org_id}/customer-invoices/{id}
func DeleteCustomerInvoiceHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse invoice ID from URL
		invoiceID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid invoice ID"))
			return
		}

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Delete invoice
		if err := service.DeleteCustomerInvoice(ctx, orgID, invoiceID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("customer invoice deleted",
			zap.String("invoice_id", invoiceID.String()),
		)

		w.WriteHeader(http.StatusNoContent)
	}
}

// ============================================================================
// CUSTOMER PAYMENT HANDLERS
// ============================================================================

// ListCustomerPaymentsHandler handles GET /api/v1/organizations/{org_id}/customer-payments
func ListCustomerPaymentsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Parse query parameters
		opts := &receivables.CustomerPaymentFilterOptions{}
		if customerID := parseUUIDQuery(r, "customer_id"); customerID != nil {
			opts.CustomerID = customerID
		}
		if isPosted := parseBoolQuery(r, "is_posted"); isPosted != nil {
			opts.IsPosted = isPosted
		}

		// Get pagination params
		page, pageSize := getPaginationParams(r)
		opts.Limit = pageSize
		opts.Offset = (page - 1) * pageSize

		// Get payments
		payments, total, err := service.ListCustomerPayments(ctx, orgID, opts)
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

// CreateCustomerPaymentHandler handles POST /api/v1/organizations/{org_id}/customer-payments
func CreateCustomerPaymentHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req receivables.CreateCustomerPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Create payment
		payment, err := service.CreateCustomerPayment(ctx, &req, orgID, userID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("customer payment created",
			zap.String("payment_id", payment.Payment.ID.String()),
			zap.String("payment_number", payment.Payment.PaymentNumber),
		)

		respondJSON(w, http.StatusCreated, payment)
	}
}

// GetCustomerPaymentHandler handles GET /api/v1/organizations/{org_id}/customer-payments/{id}
func GetCustomerPaymentHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Create repository and service using stdlib compatibility layer
		repo := postgres.NewReceivablesRepository(db.GetStdlibDB())
		service := receivables.NewService(repo)

		// Get payment
		payment, err := service.GetCustomerPayment(ctx, orgID, paymentID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, payment)
	}
}

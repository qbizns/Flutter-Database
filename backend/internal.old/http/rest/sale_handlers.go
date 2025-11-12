package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/sales"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

func ListSalesHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		repo := postgres.NewSaleRepository(db)
		service := sales.NewService(repo, logger)

		var customerID, cashierID *uuid.UUID
		if customerIDStr := r.URL.Query().Get("customer_id"); customerIDStr != "" {
			if parsed, err := uuid.Parse(customerIDStr); err == nil {
				customerID = &parsed
			}
		}
		if cashierIDStr := r.URL.Query().Get("cashier_id"); cashierIDStr != "" {
			if parsed, err := uuid.Parse(cashierIDStr); err == nil {
				cashierID = &parsed
			}
		}

		var startDate, endDate *time.Time
		if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
			if parsed, err := time.Parse(time.RFC3339, startDateStr); err == nil {
				startDate = &parsed
			}
		}
		if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
			if parsed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
				endDate = &parsed
			}
		}

		paymentStatus := r.URL.Query().Get("payment_status")
		transactionType := r.URL.Query().Get("transaction_type")

		filters := sales.SaleFilters{
			Search:          r.URL.Query().Get("search"),
			CustomerID:      customerID,
			CashierID:       cashierID,
			PaymentStatus:   stringPtr(paymentStatus),
			TransactionType: stringPtr(transactionType),
			StartDate:       startDate,
			EndDate:         endDate,
		}
		filters.Page, filters.PageSize = getPaginationParams(r)

		saleList, err := service.List(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Get total count for pagination
		count, err := service.Count(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		totalPages := int(count) / filters.PageSize
		if int(count)%filters.PageSize > 0 {
			totalPages++
		}

		response := PaginatedResponse{
			Data:       saleList,
			Page:       filters.Page,
			PageSize:   filters.PageSize,
			TotalCount: count,
			TotalPages: totalPages,
		}

		respondJSON(w, http.StatusOK, response)
	}
}

func CreateSaleHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)
		userID := appctx.MustGetUserID(ctx)

		var req CreateSaleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewSaleRepository(db)
		service := sales.NewService(repo, logger)

		// Generate sale number if not provided
		saleNumber := req.SaleNumber
		if saleNumber == "" {
			saleNumber = generateSaleNumber()
		}

		// Convert sale items
		var saleItems []sales.SaleItem
		for _, item := range req.Items {
			saleItems = append(saleItems, sales.SaleItem{
				ProductID:      &item.ProductID,
				Quantity:       item.Quantity,
				UnitPrice:      item.UnitPrice,
				TaxRate:        item.TaxRate,
				DiscountAmount: item.DiscountAmount,
				Subtotal:       item.UnitPrice * item.Quantity,
				TaxAmount:      (item.UnitPrice * item.Quantity) * (item.TaxRate / 100),
				Total:          item.LineTotal,
			})
		}

		sale := &sales.Sale{
			OrganizationID:    orgID,
			SaleNumber:        saleNumber,
			TransactionType:   "sale",
			CustomerID:        req.CustomerID,
			CashierID:         &userID,
			Subtotal:          req.Subtotal,
			TaxAmount:         req.TaxAmount,
			DiscountAmount:    req.DiscountAmount,
			TotalAmount:       req.TotalAmount,
			PaidAmount:        req.TotalAmount,
			ChangeAmount:      0,
			OutstandingAmount: 0,
			PaymentStatus:     "paid",
			TransactionDate:   req.SaleDate,
			Notes:             req.Notes,
			Items:             saleItems,
			CreatedBy:         userID,
		}

		if err := service.Create(ctx, sale); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("sale created", zap.String("sale_id", sale.ID.String()), zap.String("sale_number", sale.SaleNumber))
		respondJSON(w, http.StatusCreated, sale)
	}
}

func GetSaleHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		saleID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid sale ID"))
			return
		}

		repo := postgres.NewSaleRepository(db)
		service := sales.NewService(repo, logger)

		// Get sale with items
		sale, err := service.GetWithItems(ctx, orgID, saleID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, sale)
	}
}

func UpdateSaleHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)
		userID := appctx.MustGetUserID(ctx)

		saleID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid sale ID"))
			return
		}

		var req CreateSaleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewSaleRepository(db)
		service := sales.NewService(repo, logger)

		sale, err := service.Get(ctx, orgID, saleID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Update fields
		sale.CustomerID = req.CustomerID
		sale.Subtotal = req.Subtotal
		sale.TaxAmount = req.TaxAmount
		sale.DiscountAmount = req.DiscountAmount
		sale.TotalAmount = req.TotalAmount
		sale.TransactionDate = req.SaleDate
		sale.Notes = req.Notes
		sale.UpdatedBy = &userID

		if err := service.Update(ctx, sale); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("sale updated", zap.String("sale_id", saleID.String()))
		respondJSON(w, http.StatusOK, sale)
	}
}

func DeleteSaleHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		saleID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid sale ID"))
			return
		}

		repo := postgres.NewSaleRepository(db)
		service := sales.NewService(repo, logger)

		if err := service.Delete(ctx, orgID, saleID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("sale deleted", zap.String("sale_id", saleID.String()))
		w.WriteHeader(http.StatusNoContent)
	}
}

// Helper function to generate sale number
func generateSaleNumber() string {
	return "SALE-" + time.Now().Format("20060102-150405") + "-" + uuid.New().String()[:8]
}

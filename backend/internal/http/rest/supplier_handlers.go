package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/suppliers"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

func ListSuppliersHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		repo := postgres.NewSupplierRepository(db)
		service := suppliers.NewService(repo, logger)

		status := r.URL.Query().Get("status")
		filters := suppliers.SupplierFilters{
			Search: r.URL.Query().Get("search"),
			Status: stringPtr(status),
		}
		filters.Page, filters.PageSize = getPaginationParams(r)

		supplierList, err := service.List(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, supplierList)
	}
}

func CreateSupplierHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		var req CreateSupplierRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewSupplierRepository(db)
		service := suppliers.NewService(repo, logger)

		supplier := &suppliers.Supplier{
			OrganizationID: orgID,
			Name:           req.Name,
			ContactPerson:  req.ContactPerson,
			Email:          req.Email,
			Phone:          req.Phone,
			Address:        req.AddressLine1,
			City:           req.City,
			State:          req.State,
			Country:        req.Country,
			PostalCode:     req.PostalCode,
			TaxNumber:      req.TaxID,
			PaymentTerms:   "",
			CreditLimit:    req.CreditLimit,
			Status:         "active",
			CreatedBy:      userID,
		}

		if err := service.Create(ctx, supplier); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("supplier created", zap.String("supplier_id", supplier.ID.String()))
		respondJSON(w, http.StatusCreated, supplier)
	}
}

func GetSupplierHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		supplierID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid supplier ID"))
			return
		}

		repo := postgres.NewSupplierRepository(db)
		service := suppliers.NewService(repo, logger)

		supplier, err := service.Get(ctx, orgID, supplierID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, supplier)
	}
}

func UpdateSupplierHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		supplierID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid supplier ID"))
			return
		}

		var req UpdateSupplierRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewSupplierRepository(db)
		service := suppliers.NewService(repo, logger)

		supplier, err := service.Get(ctx, orgID, supplierID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Apply updates
		if req.Name != nil {
			supplier.Name = *req.Name
		}
		if req.ContactPerson != nil {
			supplier.ContactPerson = *req.ContactPerson
		}
		if req.Email != nil {
			supplier.Email = *req.Email
		}
		if req.Phone != nil {
			supplier.Phone = *req.Phone
		}
		if req.AddressLine1 != nil {
			supplier.Address = *req.AddressLine1
		}
		if req.City != nil {
			supplier.City = *req.City
		}
		if req.State != nil {
			supplier.State = *req.State
		}
		if req.PostalCode != nil {
			supplier.PostalCode = *req.PostalCode
		}
		if req.Country != nil {
			supplier.Country = *req.Country
		}
		if req.TaxID != nil {
			supplier.TaxNumber = *req.TaxID
		}
		if req.CreditLimit != nil {
			supplier.CreditLimit = *req.CreditLimit
		}
		if req.IsActive != nil {
			if *req.IsActive {
				supplier.Status = "active"
			} else {
				supplier.Status = "inactive"
			}
		}
		supplier.UpdatedBy = &userID

		if err := service.Update(ctx, supplier); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("supplier updated", zap.String("supplier_id", supplierID.String()))
		respondJSON(w, http.StatusOK, supplier)
	}
}

func DeleteSupplierHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		supplierID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid supplier ID"))
			return
		}

		repo := postgres.NewSupplierRepository(db)
		service := suppliers.NewService(repo, logger)

		if err := service.Delete(ctx, orgID, supplierID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("supplier deleted", zap.String("supplier_id", supplierID.String()))
		w.WriteHeader(http.StatusNoContent)
	}
}

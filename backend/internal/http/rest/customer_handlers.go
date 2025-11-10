package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/customers"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

func ListCustomersHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		repo := postgres.NewCustomerRepository(db)
		service := customers.NewService(repo, logger)

		customerType := r.URL.Query().Get("customer_type")
		filters := customers.CustomerFilters{
			Search:       r.URL.Query().Get("search"),
			CustomerType: stringPtr(customerType),
			IsActive:     parseBoolQuery(r, "is_active"),
		}
		filters.Page, filters.PageSize = getPaginationParams(r)

		customerList, err := service.List(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, customerList)
	}
}

func CreateCustomerHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		var req CreateCustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewCustomerRepository(db)
		service := customers.NewService(repo, logger)

		customer := &customers.Customer{
			OrganizationID:      orgID,
			FirstName:           req.FirstName,
			LastName:            req.LastName,
			Email:               req.Email,
			Phone:               req.Phone,
			DateOfBirth:         req.DateOfBirth,
			Gender:              req.Gender,
			AddressLine1:        req.AddressLine1,
			AddressLine2:        req.AddressLine2,
			City:                req.City,
			State:               req.State,
			PostalCode:          req.PostalCode,
			Country:             req.Country,
			CustomerType:        req.CustomerType,
			TaxID:               req.TaxID,
			PaymentTermDays:     req.PaymentTermDays,
			CreditLimit:         req.CreditLimit,
			IsActive:            req.IsActive,
			CreatedBy:           userID,
		}

		if err := service.Create(ctx, customer); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("customer created", zap.String("customer_id", customer.ID.String()))
		respondJSON(w, http.StatusCreated, customer)
	}
}

func GetCustomerHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		customerID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid customer ID"))
			return
		}

		repo := postgres.NewCustomerRepository(db)
		service := customers.NewService(repo, logger)

		customer, err := service.Get(ctx, orgID, customerID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, customer)
	}
}

func UpdateCustomerHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		customerID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid customer ID"))
			return
		}

		var req UpdateCustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewCustomerRepository(db)
		service := customers.NewService(repo, logger)

		customer, err := service.Get(ctx, orgID, customerID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Apply updates
		if req.FirstName != nil {
			customer.FirstName = *req.FirstName
		}
		if req.LastName != nil {
			customer.LastName = *req.LastName
		}
		if req.Email != nil {
			customer.Email = *req.Email
		}
		if req.Phone != nil {
			customer.Phone = *req.Phone
		}
		if req.DateOfBirth != nil {
			customer.DateOfBirth = req.DateOfBirth
		}
		if req.Gender != nil {
			customer.Gender = *req.Gender
		}
		if req.AddressLine1 != nil {
			customer.AddressLine1 = *req.AddressLine1
		}
		if req.AddressLine2 != nil {
			customer.AddressLine2 = *req.AddressLine2
		}
		if req.City != nil {
			customer.City = *req.City
		}
		if req.State != nil {
			customer.State = *req.State
		}
		if req.PostalCode != nil {
			customer.PostalCode = *req.PostalCode
		}
		if req.Country != nil {
			customer.Country = *req.Country
		}
		if req.CustomerType != nil {
			customer.CustomerType = *req.CustomerType
		}
		if req.TaxID != nil {
			customer.TaxID = *req.TaxID
		}
		if req.PaymentTermDays != nil {
			customer.PaymentTermDays = *req.PaymentTermDays
		}
		if req.CreditLimit != nil {
			customer.CreditLimit = *req.CreditLimit
		}
		if req.IsActive != nil {
			customer.IsActive = *req.IsActive
		}
		customer.UpdatedBy = &userID

		if err := service.Update(ctx, customer); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("customer updated", zap.String("customer_id", customerID.String()))
		respondJSON(w, http.StatusOK, customer)
	}
}

func DeleteCustomerHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		customerID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid customer ID"))
			return
		}

		repo := postgres.NewCustomerRepository(db)
		service := customers.NewService(repo, logger)

		if err := service.Delete(ctx, orgID, customerID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("customer deleted", zap.String("customer_id", customerID.String()))
		w.WriteHeader(http.StatusNoContent)
	}
}

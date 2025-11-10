package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/locations"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

func ListLocationsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		repo := postgres.NewLocationRepository(db)
		service := locations.NewService(repo, logger)

		locationType := r.URL.Query().Get("location_type")
		filters := locations.LocationFilters{
			Search:       r.URL.Query().Get("search"),
			LocationType: stringPtr(locationType),
			IsActive:     parseBoolQuery(r, "is_active"),
		}
		filters.Page, filters.PageSize = getPaginationParams(r)

		locationList, err := service.List(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, locationList)
	}
}

func CreateLocationHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		var req CreateLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewLocationRepository(db)
		service := locations.NewService(repo, logger)

		location := &locations.Location{
			OrganizationID: orgID,
			Name:           req.Name,
			Code:           req.Code,
			LocationType:   req.LocationType,
			AddressLine1:   req.AddressLine1,
			AddressLine2:   req.AddressLine2,
			City:           req.City,
			State:          req.State,
			PostalCode:     req.PostalCode,
			Country:        req.Country,
			Phone:          req.Phone,
			Email:          req.Email,
			IsActive:       req.IsActive,
			CreatedBy:      userID,
		}

		if err := service.Create(ctx, location); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("location created", zap.String("location_id", location.ID.String()))
		respondJSON(w, http.StatusCreated, location)
	}
}

func GetLocationHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		locationID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid location ID"))
			return
		}

		repo := postgres.NewLocationRepository(db)
		service := locations.NewService(repo, logger)

		location, err := service.Get(ctx, orgID, locationID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, location)
	}
}

func UpdateLocationHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		locationID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid location ID"))
			return
		}

		var req UpdateLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewLocationRepository(db)
		service := locations.NewService(repo, logger)

		location, err := service.Get(ctx, orgID, locationID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Apply updates
		if req.Name != nil {
			location.Name = *req.Name
		}
		if req.Code != nil {
			location.Code = *req.Code
		}
		if req.LocationType != nil {
			location.LocationType = *req.LocationType
		}
		if req.AddressLine1 != nil {
			location.AddressLine1 = *req.AddressLine1
		}
		if req.AddressLine2 != nil {
			location.AddressLine2 = *req.AddressLine2
		}
		if req.City != nil {
			location.City = *req.City
		}
		if req.State != nil {
			location.State = *req.State
		}
		if req.PostalCode != nil {
			location.PostalCode = *req.PostalCode
		}
		if req.Country != nil {
			location.Country = *req.Country
		}
		if req.Phone != nil {
			location.Phone = *req.Phone
		}
		if req.Email != nil {
			location.Email = *req.Email
		}
		if req.IsActive != nil {
			location.IsActive = *req.IsActive
		}
		location.UpdatedBy = &userID

		if err := service.Update(ctx, location); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("location updated", zap.String("location_id", locationID.String()))
		respondJSON(w, http.StatusOK, location)
	}
}

func DeleteLocationHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		locationID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid location ID"))
			return
		}

		repo := postgres.NewLocationRepository(db)
		service := locations.NewService(repo, logger)

		if err := service.Delete(ctx, orgID, locationID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("location deleted", zap.String("location_id", locationID.String()))
		w.WriteHeader(http.StatusNoContent)
	}
}

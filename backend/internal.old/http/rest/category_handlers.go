package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/categories"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

func ListCategoriesHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		repo := postgres.NewCategoryRepository(db)
		service := categories.NewService(repo, logger)

		var parentID *uuid.UUID
		if parentIDStr := r.URL.Query().Get("parent_id"); parentIDStr != "" {
			if parsed, err := uuid.Parse(parentIDStr); err == nil {
				parentID = &parsed
			}
		}

		filters := categories.CategoryFilters{
			Search:   r.URL.Query().Get("search"),
			ParentID: parentID,
			IsActive: parseBoolQuery(r, "is_active"),
		}
		filters.Page, filters.PageSize = getPaginationParams(r)

		categoryList, err := service.List(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, categoryList)
	}
}

func CreateCategoryHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)
		userID := appctx.MustGetUserID(ctx)

		var req CreateCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewCategoryRepository(db)
		service := categories.NewService(repo, logger)

		category := &categories.Category{
			OrganizationID: orgID,
			Name:           req.Name,
			Description:    req.Description,
			ParentID:       req.ParentID,
			Color:          req.Color,
			IconName:       req.IconName,
			DisplayOrder:   req.DisplayOrder,
			IsActive:       req.IsActive,
			CreatedBy:      userID,
		}

		if err := service.Create(ctx, category); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("category created", zap.String("category_id", category.ID.String()))
		respondJSON(w, http.StatusCreated, category)
	}
}

func GetCategoryHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		categoryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid category ID"))
			return
		}

		repo := postgres.NewCategoryRepository(db)
		service := categories.NewService(repo, logger)

		category, err := service.Get(ctx, orgID, categoryID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, category)
	}
}

func UpdateCategoryHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)
		userID := appctx.MustGetUserID(ctx)

		categoryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid category ID"))
			return
		}

		var req UpdateCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.NewCategoryRepository(db)
		service := categories.NewService(repo, logger)

		category, err := service.Get(ctx, orgID, categoryID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Apply updates
		if req.Name != nil {
			category.Name = *req.Name
		}
		if req.Description != nil {
			category.Description = *req.Description
		}
		if req.ParentID != nil {
			category.ParentID = req.ParentID
		}
		if req.Color != nil {
			category.Color = *req.Color
		}
		if req.IconName != nil {
			category.IconName = *req.IconName
		}
		if req.DisplayOrder != nil {
			category.DisplayOrder = *req.DisplayOrder
		}
		if req.IsActive != nil {
			category.IsActive = *req.IsActive
		}
		category.UpdatedBy = &userID

		if err := service.Update(ctx, category); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("category updated", zap.String("category_id", categoryID.String()))
		respondJSON(w, http.StatusOK, category)
	}
}

func DeleteCategoryHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		categoryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid category ID"))
			return
		}

		repo := postgres.NewCategoryRepository(db)
		service := categories.NewService(repo, logger)

		if err := service.Delete(ctx, orgID, categoryID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("category deleted", zap.String("category_id", categoryID.String()))
		w.WriteHeader(http.StatusNoContent)
	}
}

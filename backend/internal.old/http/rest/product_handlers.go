package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/products"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

// ListProductsHandler handles GET /api/v1/organizations/{org_id}/products
func ListProductsHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Create repository and service
		repo := postgres.NewProductRepository(db)
		service := products.NewService(repo, logger)

		// Parse query parameters
		filters := products.ProductFilters{
			Search:     r.URL.Query().Get("search"),
			CategoryID: parseUUIDQuery(r, "category_id"),
			IsActive:   parseBoolQuery(r, "is_active"),
		}

		// Get products
		productList, err := service.List(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, productList)
	}
}

// CreateProductHandler handles POST /api/v1/organizations/{org_id}/products
func CreateProductHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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
		var req CreateProductRequest
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
		repo := postgres.NewProductRepository(db)
		service := products.NewService(repo, logger)

		// Create product
		product := &products.Product{
			OrganizationID:      orgID,
			Name:                req.Name,
			SKU:                 req.SKU,
			Barcode:             req.Barcode,
			Description:         req.Description,
			CategoryID:          req.CategoryID,
			UnitPrice:           req.UnitPrice,
			Cost:                req.Cost,
			TaxRate:             req.TaxRate,
			Unit:                req.Unit,
			MinStockLevel:       req.MinStockLevel,
			MaxStockLevel:       req.MaxStockLevel,
			IsActive:            req.IsActive,
			IsTrackInventory:    req.IsTrackInventory,
			AllowNegativeStock:  req.AllowNegativeStock,
			ImageURL:            req.ImageURL,
			ProductType:         req.ProductType,
			AccountingAccountID: req.AccountingAccountID,
			CreatedBy:           userID,
		}

		if err := service.Create(ctx, product); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("product created",
			zap.String("product_id", product.ID.String()),
			zap.String("sku", product.SKU),
		)

		respondJSON(w, http.StatusCreated, product)
	}
}

// GetProductHandler handles GET /api/v1/organizations/{org_id}/products/{id}
func GetProductHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse product ID from URL
		productID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid product ID"))
			return
		}

		// Create repository and service
		repo := postgres.NewProductRepository(db)
		service := products.NewService(repo, logger)

		// Get product
		product, err := service.Get(ctx, orgID, productID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, product)
	}
}

// UpdateProductHandler handles PATCH /api/v1/organizations/{org_id}/products/{id}
func UpdateProductHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
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

		// Parse product ID from URL
		productID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid product ID"))
			return
		}

		// Parse request body
		var req UpdateProductRequest
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
		repo := postgres.NewProductRepository(db)
		service := products.NewService(repo, logger)

		// Get existing product
		product, err := service.Get(ctx, orgID, productID)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Apply updates
		if req.Name != nil {
			product.Name = *req.Name
		}
		if req.SKU != nil {
			product.SKU = *req.SKU
		}
		if req.Barcode != nil {
			product.Barcode = *req.Barcode
		}
		if req.Description != nil {
			product.Description = *req.Description
		}
		if req.CategoryID != nil {
			product.CategoryID = req.CategoryID
		}
		if req.UnitPrice != nil {
			product.UnitPrice = *req.UnitPrice
		}
		if req.Cost != nil {
			product.Cost = *req.Cost
		}
		if req.TaxRate != nil {
			product.TaxRate = *req.TaxRate
		}
		if req.Unit != nil {
			product.Unit = *req.Unit
		}
		if req.MinStockLevel != nil {
			product.MinStockLevel = *req.MinStockLevel
		}
		if req.MaxStockLevel != nil {
			product.MaxStockLevel = *req.MaxStockLevel
		}
		if req.IsActive != nil {
			product.IsActive = *req.IsActive
		}
		if req.IsTrackInventory != nil {
			product.IsTrackInventory = *req.IsTrackInventory
		}
		if req.AllowNegativeStock != nil {
			product.AllowNegativeStock = *req.AllowNegativeStock
		}
		if req.ImageURL != nil {
			product.ImageURL = *req.ImageURL
		}
		if req.ProductType != nil {
			product.ProductType = *req.ProductType
		}
		if req.AccountingAccountID != nil {
			product.AccountingAccountID = req.AccountingAccountID
		}
		product.UpdatedBy = &userID

		// Update product
		if err := service.Update(ctx, product); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("product updated",
			zap.String("product_id", product.ID.String()),
			zap.String("sku", product.SKU),
		)

		respondJSON(w, http.StatusOK, product)
	}
}

// DeleteProductHandler handles DELETE /api/v1/organizations/{org_id}/products/{id}
func DeleteProductHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse product ID from URL
		productID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid product ID"))
			return
		}

		// Create repository and service
		repo := postgres.NewProductRepository(db)
		service := products.NewService(repo, logger)

		// Delete product (soft delete)
		if err := service.Delete(ctx, orgID, productID); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("product deleted",
			zap.String("product_id", productID.String()),
		)

		w.WriteHeader(http.StatusNoContent)
	}
}

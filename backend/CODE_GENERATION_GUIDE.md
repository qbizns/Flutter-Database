# Code Generation Guide for CRUD APIs

## Overview

This guide shows the exact pattern to implement full CRUD APIs for all 172 tables in the database. Each table requires 4-5 files following the clean architecture pattern.

## Files Required Per Table

For each table (e.g., `products`):

1. **Domain Types**: `internal/domain/products/types.go`
2. **Domain Service**: `internal/domain/products/service.go`
3. **Repository**: `internal/repository/postgres/product_repository.go`
4. **HTTP Handlers**: `internal/http/rest/product_handlers.go`
5. **Request/Response Types**: Add to `internal/http/rest/types.go`

## Step-by-Step Implementation Pattern

### Step 1: Domain Types (`internal/domain/{entity}/types.go`)

```go
package {entity}

import (
	"context"
	"time"
	"github.com/google/uuid"
)

// {Entity} represents the entity
type {Entity} struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	// ... add all fields from database schema
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	UpdatedBy      *uuid.UUID `json:"updated_by"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// {Entity}Filters represents filters for listing
type {Entity}Filters struct {
	Search   string
	IsActive *bool
	// Add relevant filters
	Page     int
	PageSize int
}

// Repository defines the data access interface
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters {Entity}Filters) ([]{Entity}, error)
	Count(ctx context.Context, orgID uuid.UUID, filters {Entity}Filters) (int64, error)
	Create(ctx context.Context, entity *{Entity}) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*{Entity}, error)
	Update(ctx context.Context, entity *{Entity}) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}
```

### Step 2: Domain Service (`internal/domain/{entity}/service.go`)

```go
package {entity}

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	repo   Repository
	logger *logging.Logger
}

func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters {Entity}Filters) ([]{Entity}, error) {
	entities, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list {entity}", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return entities, nil
}

func (s *Service) Create(ctx context.Context, entity *{Entity}) error {
	// Validate
	if err := s.validate(entity); err != nil {
		return err
	}

	// Set defaults
	entity.ID = uuid.New()
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()

	// Create
	if err := s.repo.Create(ctx, entity); err != nil {
		s.logger.Error("failed to create {entity}", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*{Entity}, error) {
	entity, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if entity == nil {
		return nil, apperrors.NotFound("{entity}")
	}
	return entity, nil
}

func (s *Service) Update(ctx context.Context, entity *{Entity}) error {
	// Validate
	if err := s.validate(entity); err != nil {
		return err
	}

	// Check exists
	existing, err := s.repo.Get(ctx, entity.OrganizationID, entity.ID)
	if err != nil {
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("{entity}")
	}

	entity.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, entity); err != nil {
		s.logger.Error("failed to update {entity}", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	existing, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("{entity}")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete {entity}", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) validate(entity *{Entity}) error {
	// Add validation rules
	if entity.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	return nil
}
```

### Step 3: Repository (`internal/repository/postgres/{entity}_repository.go`)

```go
package postgres

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/{entity}"
)

type {Entity}Repository struct {
	db *DB
}

func New{Entity}Repository(db *DB) *{Entity}Repository {
	return &{Entity}Repository{db: db}
}

func (r *{Entity}Repository) List(ctx context.Context, orgID uuid.UUID, filters {entity}.{Entity}Filters) ([]{entity}.{Entity}, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, /* list all columns */
		FROM {table_name}
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	query += " ORDER BY created_at DESC"

	// Pagination
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []{entity}.{Entity}
	for rows.Next() {
		var e {entity}.{Entity}
		err := rows.Scan(
			&e.ID, &e.OrganizationID, /* scan all fields */
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, rows.Err()
}

func (r *{Entity}Repository) Count(ctx context.Context, orgID uuid.UUID, filters {entity}.{Entity}Filters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM {table_name} WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}

	// Apply same filters as List

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *{Entity}Repository) Create(ctx context.Context, entity *{entity}.{Entity}) error {
	if err := r.db.SetOrganizationContext(ctx, entity.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO {table_name} (
			id, organization_id, /* all columns */
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, /* $3, $4, ... */ $N, $N+1, $N+2
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		entity.ID, entity.OrganizationID, /* all values */
		entity.CreatedAt, entity.UpdatedAt, entity.CreatedBy,
	)
	return err
}

func (r *{Entity}Repository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*{entity}.{Entity}, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, /* all columns */
		FROM {table_name}
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var e {entity}.{Entity}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&e.ID, &e.OrganizationID, /* scan all fields */
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *{Entity}Repository) Update(ctx context.Context, entity *{entity}.{Entity}) error {
	if err := r.db.SetOrganizationContext(ctx, entity.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE {table_name} SET
			name = $3, /* field = $N, ... */
			updated_at = $N, updated_by = $N+1
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		entity.OrganizationID, entity.ID,
		/* all update values */
		entity.UpdatedAt, entity.UpdatedBy,
	)
	return err
}

func (r *{Entity}Repository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE {table_name}
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}
```

### Step 4: HTTP Handlers (`internal/http/rest/{entity}_handlers.go`)

```go
package rest

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/{entity}"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

func List{Entity}Handler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		repo := postgres.New{Entity}Repository(db)
		service := {entity}.NewService(repo, logger)

		filters := {entity}.{Entity}Filters{
			Search: r.URL.Query().Get("search"),
			IsActive: parseBoolQuery(r, "is_active"),
		}

		entities, err := service.List(ctx, orgID, filters)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, entities)
	}
}

func Create{Entity}Handler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)
		userID := appctx.MustGetUserID(ctx)

		var req Create{Entity}Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.New{Entity}Repository(db)
		service := {entity}.NewService(repo, logger)

		entity := &{entity}.{Entity}{
			OrganizationID: orgID,
			Name: req.Name,
			// Map all request fields
			CreatedBy: userID,
		}

		if err := service.Create(ctx, entity); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("{entity} created", zap.String("id", entity.ID.String()))
		respondJSON(w, http.StatusCreated, entity)
	}
}

func Get{Entity}Handler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid ID"))
			return
		}

		repo := postgres.New{Entity}Repository(db)
		service := {entity}.NewService(repo, logger)

		entity, err := service.Get(ctx, orgID, id)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		respondJSON(w, http.StatusOK, entity)
	}
}

func Update{Entity}Handler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)
		userID := appctx.MustGetUserID(ctx)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid ID"))
			return
		}

		var req Update{Entity}Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid request body"))
			return
		}

		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		repo := postgres.New{Entity}Repository(db)
		service := {entity}.NewService(repo, logger)

		entity, err := service.Get(ctx, orgID, id)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Apply updates (check if field is not nil, then update)
		if req.Name != nil {
			entity.Name = *req.Name
		}
		// ... for all fields

		entity.UpdatedBy = &userID

		if err := service.Update(ctx, entity); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("{entity} updated", zap.String("id", id.String()))
		respondJSON(w, http.StatusOK, entity)
	}
}

func Delete{Entity}Handler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := appctx.MustGetOrganizationID(ctx)

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, logger, apperrors.BadRequest("Invalid ID"))
			return
		}

		repo := postgres.New{Entity}Repository(db)
		service := {entity}.NewService(repo, logger)

		if err := service.Delete(ctx, orgID, id); err != nil {
			respondError(w, logger, err)
			return
		}

		logger.Info("{entity} deleted", zap.String("id", id.String()))
		w.WriteHeader(http.StatusNoContent)
	}
}
```

### Step 5: Request/Response Types

Add to `internal/http/rest/types.go`:

```go
// {ENTITY}
type Create{Entity}Request struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	// ... all required fields
}

type Update{Entity}Request struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=255"`
	// ... all optional fields (use pointers)
}
```

### Step 6: Register Routes

Add to `cmd/api/main.go`:

```go
// {Entity}
r.Get("/{entity}", rest.List{Entity}Handler(db, logger))
r.Post("/{entity}", rest.Create{Entity}Handler(db, logger))
r.Get("/{entity}/{id}", rest.Get{Entity}Handler(db, logger))
r.Patch("/{entity}/{id}", rest.Update{Entity}Handler(db, logger))
r.Delete("/{entity}/{id}", rest.Delete{Entity}Handler(db, logger))
```

## Quick Reference

### Naming Conventions
- **Entity**: `Product` (singular, capitalized)
- **Package**: `products` (plural, lowercase)
- **Table**: `products` (plural, lowercase)
- **File**: `product_repository.go` (singular)

### Common Patterns

**UUID Fields**: Use `*uuid.UUID` for optional foreign keys
**Timestamps**: Always include `created_at`, `updated_at`, `deleted_at`
**Soft Delete**: Set `deleted_at` instead of hard delete
**RLS**: Always call `SetOrganizationContext` before queries
**Filters**: Support search, pagination, and entity-specific filters
**Validation**: Validate in service layer, not handlers

## Automation

For faster implementation, consider creating code generation scripts using this template pattern. See `scripts/generate-crud.sh` for automation.

## Next Steps

1. Implement remaining Tier 1 tables (11 total)
2. Implement Tier 2 & 3 (core operations + posting engine)
3. Continue through all tiers systematically
4. Write integration tests for each tier
5. Add OpenAPI documentation

## Reference Implementations

Complete reference implementations:
- **Products**: Full CRUD with stock management
- **Customers**: Full CRUD with loyalty integration

Use these as templates for all other tables.

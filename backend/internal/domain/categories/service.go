package categories

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles category business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new category service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of categories with filters
func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters CategoryFilters) ([]Category, error) {
	categories, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list categories", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return categories, nil
}

// Create creates a new category
func (s *Service) Create(ctx context.Context, category *Category) error {
	// Validate category
	if err := s.validate(category); err != nil {
		return err
	}

	// Auto-generate slug if not provided
	if category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}

	// Check for duplicate slug
	existing, err := s.repo.GetBySlug(ctx, category.OrganizationID, category.Slug)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate slug", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("category", "slug already exists")
	}

	// If this is a child category, calculate level and path
	if category.ParentID != nil {
		parent, err := s.repo.Get(ctx, category.OrganizationID, *category.ParentID)
		if err != nil {
			s.logger.Error("failed to get parent category", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if parent == nil {
			return apperrors.NotFound("parent category")
		}

		category.Level = parent.Level + 1
		if parent.Path != "" {
			category.Path = parent.Path + "/" + category.Slug
		} else {
			category.Path = parent.Slug + "/" + category.Slug
		}
	} else {
		// Root category
		category.Level = 0
		category.Path = category.Slug
	}

	// Set defaults
	category.ID = uuid.New()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	if category.Metadata == nil {
		category.Metadata = []byte("{}")
	}

	// Create category
	if err := s.repo.Create(ctx, category); err != nil {
		s.logger.Error("failed to create category", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a category by ID
func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Category, error) {
	category, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get category", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if category == nil {
		return nil, apperrors.NotFound("category")
	}

	return category, nil
}

// GetChildren retrieves all child categories of a parent
func (s *Service) GetChildren(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) ([]Category, error) {
	children, err := s.repo.GetChildren(ctx, orgID, parentID)
	if err != nil {
		s.logger.Error("failed to get child categories", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return children, nil
}

// Update updates an existing category
func (s *Service) Update(ctx context.Context, category *Category) error {
	// Validate category
	if err := s.validate(category); err != nil {
		return err
	}

	// Check if category exists
	existing, err := s.repo.Get(ctx, category.OrganizationID, category.ID)
	if err != nil {
		s.logger.Error("failed to get category", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("category")
	}

	// Prevent circular parent reference
	if category.ParentID != nil && *category.ParentID == category.ID {
		return apperrors.ValidationFailed("category cannot be its own parent")
	}

	// Check for duplicate slug (if slug changed)
	if category.Slug != existing.Slug {
		duplicate, err := s.repo.GetBySlug(ctx, category.OrganizationID, category.Slug)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate slug", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != category.ID {
			return apperrors.AlreadyExists("category", "slug already exists")
		}
	}

	// Recalculate level and path if parent changed
	if (category.ParentID == nil && existing.ParentID != nil) ||
		(category.ParentID != nil && existing.ParentID == nil) ||
		(category.ParentID != nil && existing.ParentID != nil && *category.ParentID != *existing.ParentID) {

		if category.ParentID != nil {
			parent, err := s.repo.Get(ctx, category.OrganizationID, *category.ParentID)
			if err != nil {
				s.logger.Error("failed to get parent category", zap.Error(err))
				return apperrors.DatabaseError(err)
			}
			if parent == nil {
				return apperrors.NotFound("parent category")
			}

			// Check for circular reference in hierarchy
			if err := s.checkCircularReference(ctx, category.OrganizationID, category.ID, *category.ParentID); err != nil {
				return err
			}

			category.Level = parent.Level + 1
			if parent.Path != "" {
				category.Path = parent.Path + "/" + category.Slug
			} else {
				category.Path = parent.Slug + "/" + category.Slug
			}
		} else {
			// Changed to root category
			category.Level = 0
			category.Path = category.Slug
		}
	} else if category.Slug != existing.Slug {
		// Slug changed but parent didn't - update path
		if category.ParentID != nil {
			parent, err := s.repo.Get(ctx, category.OrganizationID, *category.ParentID)
			if err != nil {
				s.logger.Error("failed to get parent category", zap.Error(err))
				return apperrors.DatabaseError(err)
			}
			if parent == nil {
				return apperrors.NotFound("parent category")
			}

			if parent.Path != "" {
				category.Path = parent.Path + "/" + category.Slug
			} else {
				category.Path = parent.Slug + "/" + category.Slug
			}
		} else {
			category.Path = category.Slug
		}
	}

	// Update timestamp
	category.UpdatedAt = time.Now()

	// Update category
	if err := s.repo.Update(ctx, category); err != nil {
		s.logger.Error("failed to update category", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete soft-deletes a category
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if category exists
	category, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get category", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if category == nil {
		return apperrors.NotFound("category")
	}

	// Check if category has children
	children, err := s.repo.GetChildren(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to check for child categories", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if len(children) > 0 {
		return apperrors.ValidationFailed(fmt.Sprintf("cannot delete category with %d child categories", len(children)))
	}

	// Soft delete category
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete category", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validate validates a category
func (s *Service) validate(category *Category) error {
	if category.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	if len(category.Name) > 255 {
		return apperrors.ValidationFailed("name must not exceed 255 characters")
	}
	if category.Slug != "" && !isValidSlug(category.Slug) {
		return apperrors.ValidationFailed("slug must contain only lowercase letters, numbers, and hyphens")
	}
	if len(category.Slug) > 255 {
		return apperrors.ValidationFailed("slug must not exceed 255 characters")
	}

	return nil
}

// generateSlug generates a URL-friendly slug from a name
func (s *Service) generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

// isValidSlug checks if a slug is valid (lowercase letters, numbers, and hyphens only)
func isValidSlug(slug string) bool {
	match, _ := regexp.MatchString("^[a-z0-9-]+$", slug)
	return match
}

// checkCircularReference checks if setting a parent would create a circular reference
func (s *Service) checkCircularReference(ctx context.Context, orgID uuid.UUID, categoryID uuid.UUID, parentID uuid.UUID) error {
	// Walk up the parent chain to ensure we don't encounter categoryID
	currentID := parentID
	visited := make(map[uuid.UUID]bool)
	maxDepth := 100 // Prevent infinite loops

	for i := 0; i < maxDepth; i++ {
		// Check if we've seen this ID before (circular reference)
		if visited[currentID] {
			return apperrors.ValidationFailed("circular reference detected in category hierarchy")
		}
		visited[currentID] = true

		// If we encounter the category we're updating, it's a circular reference
		if currentID == categoryID {
			return apperrors.ValidationFailed("circular reference detected in category hierarchy")
		}

		// Get the current category's parent
		current, err := s.repo.Get(ctx, orgID, currentID)
		if err != nil {
			return apperrors.DatabaseError(err)
		}
		if current == nil {
			return apperrors.NotFound("parent category in hierarchy")
		}

		// If no parent, we've reached the root
		if current.ParentID == nil {
			break
		}

		currentID = *current.ParentID
	}

	return nil
}

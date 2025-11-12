package locations

import (
	"context"
	"strings"
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
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters LocationFilters) ([]Location, error) {
	locations, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list locations", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return locations, nil
}

func (s *Service) Count(ctx context.Context, orgID uuid.UUID, filters LocationFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count locations", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *Service) Create(ctx context.Context, location *Location) error {
	// Validate
	if err := s.validate(location); err != nil {
		return err
	}

	// Check for duplicate location code
	existing, err := s.repo.GetByCode(ctx, location.OrganizationID, location.LocationCode)
	if err != nil {
		s.logger.Error("failed to check duplicate location code", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.ValidationFailed("location code already exists for this organization")
	}

	// Set defaults
	location.ID = uuid.New()
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()

	// Default values if not set
	if location.Timezone == "" {
		location.Timezone = "UTC"
	}
	if location.LocationType == "" {
		location.LocationType = "store"
	}

	// Create
	if err := s.repo.Create(ctx, location); err != nil {
		s.logger.Error("failed to create location", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("location created", zap.String("id", location.ID.String()), zap.String("code", location.LocationCode))
	return nil
}

func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Location, error) {
	location, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get location", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if location == nil {
		return nil, apperrors.NotFound("location")
	}
	return location, nil
}

func (s *Service) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*Location, error) {
	location, err := s.repo.GetByCode(ctx, orgID, code)
	if err != nil {
		s.logger.Error("failed to get location by code", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if location == nil {
		return nil, apperrors.NotFound("location")
	}
	return location, nil
}

func (s *Service) Update(ctx context.Context, location *Location) error {
	// Validate
	if err := s.validate(location); err != nil {
		return err
	}

	// Check exists
	existing, err := s.repo.Get(ctx, location.OrganizationID, location.ID)
	if err != nil {
		s.logger.Error("failed to check existing location", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("location")
	}

	// Check for duplicate location code (if changed)
	if existing.LocationCode != location.LocationCode {
		duplicate, err := s.repo.GetByCode(ctx, location.OrganizationID, location.LocationCode)
		if err != nil {
			s.logger.Error("failed to check duplicate location code", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil {
			return apperrors.ValidationFailed("location code already exists for this organization")
		}
	}

	location.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, location); err != nil {
		s.logger.Error("failed to update location", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("location updated", zap.String("id", location.ID.String()))
	return nil
}

func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	existing, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to check existing location", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("location")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete location", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("location deleted", zap.String("id", id.String()))
	return nil
}

func (s *Service) validate(location *Location) error {
	// Required fields
	if strings.TrimSpace(location.Name) == "" {
		return apperrors.ValidationFailed("name is required")
	}

	if strings.TrimSpace(location.LocationCode) == "" {
		return apperrors.ValidationFailed("location code is required")
	}

	// Validate location type
	validTypes := map[string]bool{
		"store":        true,
		"warehouse":    true,
		"headquarters": true,
		"kiosk":        true,
		"online":       true,
		"other":        true,
	}
	if !validTypes[location.LocationType] {
		return apperrors.ValidationFailed("invalid location type")
	}

	// Validate tax rate
	if location.TaxRate < 0 || location.TaxRate > 100 {
		return apperrors.ValidationFailed("tax rate must be between 0 and 100")
	}

	// Validate email format if provided
	if location.Email != nil && *location.Email != "" {
		email := strings.TrimSpace(*location.Email)
		if !strings.Contains(email, "@") {
			return apperrors.ValidationFailed("invalid email format")
		}
	}

	return nil
}

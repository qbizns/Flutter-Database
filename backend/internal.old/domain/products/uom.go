package products

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UnitOfMeasure represents a unit of measure (kg, g, piece, liter, etc.)
type UnitOfMeasure struct {
	ID         uuid.UUID  `json:"id"`
	OrganizationID *uuid.UUID `json:"organization_id"` // NULL for global UoMs
	Code       string     `json:"code"`
	Name       string     `json:"name"`
	Type       string     `json:"type"` // unit, weight, volume, length, time
	IsBaseUnit bool       `json:"is_base_unit"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// UoMConversion represents conversion between two units of measure
type UoMConversion struct {
	ID                 uuid.UUID  `json:"id"`
	FromUoMID          uuid.UUID  `json:"from_uom_id"`
	ToUoMID            uuid.UUID  `json:"to_uom_id"`
	ConversionFactor   float64    `json:"conversion_factor"`
	CreatedAt          time.Time  `json:"created_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// UoMFilters for querying UoMs
type UoMFilters struct {
	Type      string
	IsActive  *bool
	Search    string
	Page      int
	PageSize  int
}

// UoMRepository defines UoM data access interface
type UoMRepository interface {
	// UoM operations
	CreateUoM(ctx context.Context, uom *UnitOfMeasure) error
	GetUoM(ctx context.Context, id uuid.UUID) (*UnitOfMeasure, error)
	GetUoMByCode(ctx context.Context, code string, orgID *uuid.UUID) (*UnitOfMeasure, error)
	UpdateUoM(ctx context.Context, uom *UnitOfMeasure) error
	DeleteUoM(ctx context.Context, id uuid.UUID) error
	ListUoMs(ctx context.Context, filters UoMFilters) ([]UnitOfMeasure, error)
	CountUoMs(ctx context.Context, filters UoMFilters) (int64, error)

	// Conversion operations
	CreateConversion(ctx context.Context, conversion *UoMConversion) error
	GetConversion(ctx context.Context, fromID uuid.UUID, toID uuid.UUID) (*UoMConversion, error)
	DeleteConversion(ctx context.Context, fromID uuid.UUID, toID uuid.UUID) error
	ListConversionsFrom(ctx context.Context, fromID uuid.UUID) ([]UoMConversion, error)
	ListConversionsTo(ctx context.Context, toID uuid.UUID) ([]UoMConversion, error)

	// Batch operations
	ListByType(ctx context.Context, uomType string) ([]UnitOfMeasure, error)
	ConvertQuantity(ctx context.Context, fromID uuid.UUID, toID uuid.UUID, quantity float64) (float64, error)
}

// UoMService handles UoM business logic
type UoMService struct {
	repo UoMRepository
}

// NewUoMService creates a new UoM service
func NewUoMService(repo UoMRepository) *UoMService {
	return &UoMService{repo: repo}
}

// CreateUoM creates a new unit of measure
func (s *UoMService) CreateUoM(ctx context.Context, uom *UnitOfMeasure) error {
	if err := s.validateUoM(uom); err != nil {
		return err
	}

	uom.ID = uuid.New()
	uom.CreatedAt = time.Now()

	return s.repo.CreateUoM(ctx, uom)
}

// GetUoM retrieves a UoM by ID
func (s *UoMService) GetUoM(ctx context.Context, id uuid.UUID) (*UnitOfMeasure, error) {
	return s.repo.GetUoM(ctx, id)
}

// GetUoMByCode retrieves a UoM by code
func (s *UoMService) GetUoMByCode(ctx context.Context, code string, orgID *uuid.UUID) (*UnitOfMeasure, error) {
	return s.repo.GetUoMByCode(ctx, code, orgID)
}

// UpdateUoM updates an existing UoM
func (s *UoMService) UpdateUoM(ctx context.Context, uom *UnitOfMeasure) error {
	if err := s.validateUoM(uom); err != nil {
		return err
	}

	return s.repo.UpdateUoM(ctx, uom)
}

// DeleteUoM deletes a UoM
func (s *UoMService) DeleteUoM(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUoM(ctx, id)
}

// ListUoMs retrieves UoMs with filters
func (s *UoMService) ListUoMs(ctx context.Context, filters UoMFilters) ([]UnitOfMeasure, error) {
	return s.repo.ListUoMs(ctx, filters)
}

// ListByType retrieves UoMs by type
func (s *UoMService) ListByType(ctx context.Context, uomType string) ([]UnitOfMeasure, error) {
	return s.repo.ListByType(ctx, uomType)
}

// CreateConversion creates a conversion between two UoMs
func (s *UoMService) CreateConversion(ctx context.Context, conversion *UoMConversion) error {
	if err := s.validateConversion(conversion); err != nil {
		return err
	}

	conversion.ID = uuid.New()
	conversion.CreatedAt = time.Now()

	return s.repo.CreateConversion(ctx, conversion)
}

// GetConversion retrieves a conversion
func (s *UoMService) GetConversion(ctx context.Context, fromID uuid.UUID, toID uuid.UUID) (*UoMConversion, error) {
	return s.repo.GetConversion(ctx, fromID, toID)
}

// DeleteConversion deletes a conversion
func (s *UoMService) DeleteConversion(ctx context.Context, fromID uuid.UUID, toID uuid.UUID) error {
	return s.repo.DeleteConversion(ctx, fromID, toID)
}

// ListConversionsFrom retrieves all conversions from a UoM
func (s *UoMService) ListConversionsFrom(ctx context.Context, fromID uuid.UUID) ([]UoMConversion, error) {
	return s.repo.ListConversionsFrom(ctx, fromID)
}

// ConvertQuantity converts quantity from one UoM to another
func (s *UoMService) ConvertQuantity(ctx context.Context, fromID uuid.UUID, toID uuid.UUID, quantity float64) (float64, error) {
	if quantity < 0 {
		return 0, NewValidationError("quantity cannot be negative")
	}
	return s.repo.ConvertQuantity(ctx, fromID, toID, quantity)
}

// validateUoM validates UoM data
func (s *UoMService) validateUoM(uom *UnitOfMeasure) error {
	if uom.Code == "" {
		return NewValidationError("code is required")
	}
	if uom.Name == "" {
		return NewValidationError("name is required")
	}
	if uom.Type == "" {
		return NewValidationError("type is required")
	}

	validTypes := map[string]bool{
		"unit":   true,
		"weight": true,
		"volume": true,
		"length": true,
		"time":   true,
	}
	if !validTypes[uom.Type] {
		return NewValidationError("invalid type, must be one of: unit, weight, volume, length, time")
	}

	return nil
}

// validateConversion validates conversion data
func (s *UoMService) validateConversion(conversion *UoMConversion) error {
	if conversion.FromUoMID == uuid.Nil {
		return NewValidationError("from_uom_id is required")
	}
	if conversion.ToUoMID == uuid.Nil {
		return NewValidationError("to_uom_id is required")
	}
	if conversion.FromUoMID == conversion.ToUoMID {
		return NewValidationError("cannot convert a unit to itself")
	}
	if conversion.ConversionFactor <= 0 {
		return NewValidationError("conversion_factor must be positive")
	}

	return nil
}

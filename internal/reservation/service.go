package reservation

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/reservation"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/reservation"
	"go.uber.org/zap"
)

// Service handles business logic for Reservations
type Service struct {
	repo   *reservation.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Reservations service
func NewService(repo *reservation.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new reservations
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateReservationsRequest) (*dto.ReservationsResponse, error) {
	s.logger.Info("creating reservations",
		zap.String("organization_id", orgID.String()),
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Convert DTO to entity
	entity := &reservation.Reservations{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		TableId: req.TableId,
		
		CustomerId: req.CustomerId,
		
		ReservationNumber: req.ReservationNumber,
		
		ReservationDate: req.ReservationDate,
		
		ReservationTime: req.ReservationTime,
		
		DurationMinutes: req.DurationMinutes,
		
		PartySize: req.PartySize,
		
		CustomerName: req.CustomerName,
		
		CustomerPhone: req.CustomerPhone,
		
		CustomerEmail: req.CustomerEmail,
		
		Status: req.Status,
		
		AssignedWaiterId: req.AssignedWaiterId,
		
		AssignedAt: req.AssignedAt,
		
		SeatedAt: req.SeatedAt,
		
		CompletedAt: req.CompletedAt,
		
		SpecialRequests: req.SpecialRequests,
		
		Occasion: req.Occasion,
		
		DietaryRestrictions: req.DietaryRestrictions,
		
		ConfirmationCode: req.ConfirmationCode,
		
		ConfirmedAt: req.ConfirmedAt,
		
		ConfirmedBy: req.ConfirmedBy,
		
		ReminderSentAt: req.ReminderSentAt,
		
		NotificationPreferences: req.NotificationPreferences,
		
		CancelledAt: req.CancelledAt,
		
		CancelledBy: req.CancelledBy,
		
		CancellationReason: req.CancellationReason,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create reservations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created reservations",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a reservations by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.ReservationsResponse, error) {
	s.logger.Debug("getting reservations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("reservations not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of reservations records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.ReservationsListResponse, error) {
	s.logger.Debug("listing reservations",
		zap.String("organization_id", orgID.String()),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	// Get from database (organization-scoped)
	entities, total, err := s.repo.ListByOrganization(ctx, tx, orgID, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list reservations: %w", err)
	}

	// Convert to response
	items := make([]*dto.ReservationsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.ReservationsListResponse{
		Items: items,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing reservations
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateReservationsRequest) (*dto.ReservationsResponse, error) {
	s.logger.Info("updating reservations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("reservations not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.TableId != nil {
		entity.TableId = *req.TableId
	}
	
	if req.CustomerId != nil {
		entity.CustomerId = *req.CustomerId
	}
	
	if req.ReservationNumber != nil {
		entity.ReservationNumber = *req.ReservationNumber
	}
	
	if req.ReservationDate != nil {
		entity.ReservationDate = *req.ReservationDate
	}
	
	if req.ReservationTime != nil {
		entity.ReservationTime = *req.ReservationTime
	}
	
	if req.DurationMinutes != nil {
		entity.DurationMinutes = *req.DurationMinutes
	}
	
	if req.PartySize != nil {
		entity.PartySize = *req.PartySize
	}
	
	if req.CustomerName != nil {
		entity.CustomerName = *req.CustomerName
	}
	
	if req.CustomerPhone != nil {
		entity.CustomerPhone = *req.CustomerPhone
	}
	
	if req.CustomerEmail != nil {
		entity.CustomerEmail = *req.CustomerEmail
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.AssignedWaiterId != nil {
		entity.AssignedWaiterId = *req.AssignedWaiterId
	}
	
	if req.AssignedAt != nil {
		entity.AssignedAt = *req.AssignedAt
	}
	
	if req.SeatedAt != nil {
		entity.SeatedAt = *req.SeatedAt
	}
	
	if req.CompletedAt != nil {
		entity.CompletedAt = *req.CompletedAt
	}
	
	if req.SpecialRequests != nil {
		entity.SpecialRequests = *req.SpecialRequests
	}
	
	if req.Occasion != nil {
		entity.Occasion = *req.Occasion
	}
	
	if req.DietaryRestrictions != nil {
		entity.DietaryRestrictions = *req.DietaryRestrictions
	}
	
	if req.ConfirmationCode != nil {
		entity.ConfirmationCode = *req.ConfirmationCode
	}
	
	if req.ConfirmedAt != nil {
		entity.ConfirmedAt = *req.ConfirmedAt
	}
	
	if req.ConfirmedBy != nil {
		entity.ConfirmedBy = *req.ConfirmedBy
	}
	
	if req.ReminderSentAt != nil {
		entity.ReminderSentAt = *req.ReminderSentAt
	}
	
	if req.NotificationPreferences != nil {
		entity.NotificationPreferences = *req.NotificationPreferences
	}
	
	if req.CancelledAt != nil {
		entity.CancelledAt = *req.CancelledAt
	}
	
	if req.CancelledBy != nil {
		entity.CancelledBy = *req.CancelledBy
	}
	
	if req.CancellationReason != nil {
		entity.CancellationReason = *req.CancellationReason
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = *req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update reservations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated reservations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a reservations
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting reservations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return err
	}

	// Verify ownership
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get reservations: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("reservations not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete reservations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted reservations",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *reservation.Reservations) *dto.ReservationsResponse {
	return &dto.ReservationsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		TableId: entity.TableId,
		
		CustomerId: entity.CustomerId,
		
		ReservationNumber: entity.ReservationNumber,
		
		ReservationDate: entity.ReservationDate,
		
		ReservationTime: entity.ReservationTime,
		
		DurationMinutes: entity.DurationMinutes,
		
		PartySize: entity.PartySize,
		
		CustomerName: entity.CustomerName,
		
		CustomerPhone: entity.CustomerPhone,
		
		CustomerEmail: entity.CustomerEmail,
		
		Status: entity.Status,
		
		AssignedWaiterId: entity.AssignedWaiterId,
		
		AssignedAt: entity.AssignedAt,
		
		SeatedAt: entity.SeatedAt,
		
		CompletedAt: entity.CompletedAt,
		
		SpecialRequests: entity.SpecialRequests,
		
		Occasion: entity.Occasion,
		
		DietaryRestrictions: entity.DietaryRestrictions,
		
		ConfirmationCode: entity.ConfirmationCode,
		
		ConfirmedAt: entity.ConfirmedAt,
		
		ConfirmedBy: entity.ConfirmedBy,
		
		ReminderSentAt: entity.ReminderSentAt,
		
		NotificationPreferences: entity.NotificationPreferences,
		
		CancelledAt: entity.CancelledAt,
		
		CancelledBy: entity.CancelledBy,
		
		CancellationReason: entity.CancellationReason,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
	}
}


// setOrganizationContext sets the organization context for RLS
func (s *Service) setOrganizationContext(ctx context.Context, tx pgx.Tx, orgID uuid.UUID) error {
	_, err := tx.Exec(ctx, "SET LOCAL app.current_organization_id = $1", orgID)
	if err != nil {
		return fmt.Errorf("failed to set organization context: %w", err)
	}
	return nil
}


// validateBusinessRules validates business rules for reservations
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *reservation.Reservations) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a reservations can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

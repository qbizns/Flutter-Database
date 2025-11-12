package pos

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// POSSessionService handles POS session business logic
type POSSessionService struct {
	repo   POSSessionRepository
	logger *logging.Logger
}

// NewPOSSessionService creates a new POS session service
func NewPOSSessionService(repo POSSessionRepository, logger *logging.Logger) *POSSessionService {
	return &POSSessionService{
		repo:   repo,
		logger: logger,
	}
}

// OpenSession opens a new POS session
func (s *POSSessionService) OpenSession(ctx context.Context, session *POSSession) error {
	// Validate session
	if err := s.validateSession(session); err != nil {
		return err
	}

	// Check for existing open session on device
	if session.DeviceID != nil {
		existing, err := s.repo.GetOpenSession(ctx, session.OrganizationID, *session.DeviceID)
		if err != nil {
			s.logger.Error("failed to check for open session", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if existing != nil {
			return apperrors.ValidationFailed("device already has an open session")
		}
	}

	// Generate session number if not provided
	if session.SessionNumber == "" {
		session.SessionNumber = generateSessionNumber()
	}

	// Initialize session
	session.ID = uuid.New()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	session.Status = SessionStatusOpen
	session.OpenedAt = time.Now()
	session.ExpectedCash = session.OpeningCash
	session.ExpectedCard = session.OpeningCard
	session.ExpectedOther = session.OpeningOther
	session.DifferenceCash = 0
	session.DifferenceCard = 0
	session.DifferenceOther = 0

	// Create session
	if err := s.repo.Create(ctx, session); err != nil {
		s.logger.Error("failed to open session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("POS session opened", zap.String("session_id", session.ID.String()))
	return nil
}

// CloseSession closes a POS session
func (s *POSSessionService) CloseSession(ctx context.Context, orgID, sessionID uuid.UUID, countedCash, countedCard, countedOther float64) error {
	// Get session
	session, err := s.repo.Get(ctx, orgID, sessionID)
	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if session == nil {
		return apperrors.NotFound("session")
	}

	// Validate session status
	if session.Status != SessionStatusOpen {
		return apperrors.ValidationFailed("session is not open")
	}

	// Update session with counted amounts
	session.CountedCash = &countedCash
	session.CountedCard = &countedCard
	session.CountedOther = &countedOther

	// Calculate differences
	session.DifferenceCash = countedCash - session.ExpectedCash
	session.DifferenceCard = countedCard - session.ExpectedCard
	session.DifferenceOther = countedOther - session.ExpectedOther

	// Update status
	session.Status = SessionStatusClosing
	now := time.Now()
	session.ClosedAt = &now
	session.UpdatedAt = now

	if err := s.repo.Update(ctx, session); err != nil {
		s.logger.Error("failed to close session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("POS session closed",
		zap.String("session_id", sessionID.String()),
		zap.Float64("difference_cash", session.DifferenceCash),
		zap.Float64("difference_card", session.DifferenceCard),
	)

	return nil
}

// ReconcileSession reconciles a closed session (after discrepancies are handled)
func (s *POSSessionService) ReconcileSession(ctx context.Context, orgID, sessionID uuid.UUID) error {
	// Get session
	session, err := s.repo.Get(ctx, orgID, sessionID)
	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if session == nil {
		return apperrors.NotFound("session")
	}

	// Validate session status
	if session.Status != SessionStatusClosing {
		return apperrors.ValidationFailed("session must be in closing status to reconcile")
	}

	// Update status
	session.Status = SessionStatusReconciled
	session.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, session); err != nil {
		s.logger.Error("failed to reconcile session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("POS session reconciled", zap.String("session_id", sessionID.String()))
	return nil
}

// Get retrieves a session by ID
func (s *POSSessionService) Get(ctx context.Context, orgID, sessionID uuid.UUID) (*POSSession, error) {
	session, err := s.repo.Get(ctx, orgID, sessionID)
	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if session == nil {
		return nil, apperrors.NotFound("session")
	}
	return session, nil
}

// List retrieves sessions with filters
func (s *POSSessionService) List(ctx context.Context, orgID uuid.UUID, filters POSSessionFilters) ([]POSSession, error) {
	sessions, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list sessions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return sessions, nil
}

// Count counts sessions with filters
func (s *POSSessionService) Count(ctx context.Context, orgID uuid.UUID, filters POSSessionFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count sessions", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

// validateSession validates a POS session
func (s *POSSessionService) validateSession(session *POSSession) error {
	if session.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}
	if session.LocationID == uuid.Nil {
		return apperrors.ValidationFailed("location_id is required")
	}
	if session.UserID == uuid.Nil {
		return apperrors.ValidationFailed("user_id is required")
	}
	if session.OpeningCash < 0 || session.OpeningCard < 0 || session.OpeningOther < 0 {
		return apperrors.ValidationFailed("opening amounts must be non-negative")
	}
	return nil
}

// generateSessionNumber generates a unique session number
func generateSessionNumber() string {
	return fmt.Sprintf("POS_%d", time.Now().UnixNano()/1000000)
}

// CashDrawerService handles cash drawer business logic
type CashDrawerService struct {
	repo   CashDrawerRepository
	logger *logging.Logger
}

// NewCashDrawerService creates a new cash drawer service
func NewCashDrawerService(repo CashDrawerRepository, logger *logging.Logger) *CashDrawerService {
	return &CashDrawerService{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new cash drawer
func (s *CashDrawerService) Create(ctx context.Context, drawer *CashDrawer) error {
	// Validate drawer
	if err := s.validateDrawer(drawer); err != nil {
		return err
	}

	// Check for duplicate code
	existing, err := s.repo.GetByCode(ctx, drawer.OrganizationID, drawer.DrawerCode)
	if err != nil {
		s.logger.Error("failed to check duplicate drawer code", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("cash_drawer", "Drawer code already exists")
	}

	// Initialize
	drawer.ID = uuid.New()
	drawer.CreatedAt = time.Now()
	drawer.UpdatedAt = time.Now()
	drawer.IsActive = true

	if err := s.repo.Create(ctx, drawer); err != nil {
		s.logger.Error("failed to create cash drawer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("Cash drawer created", zap.String("drawer_id", drawer.ID.String()))
	return nil
}

// Get retrieves a cash drawer by ID
func (s *CashDrawerService) Get(ctx context.Context, orgID, drawerID uuid.UUID) (*CashDrawer, error) {
	drawer, err := s.repo.Get(ctx, orgID, drawerID)
	if err != nil {
		s.logger.Error("failed to get cash drawer", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if drawer == nil {
		return nil, apperrors.NotFound("cash_drawer")
	}
	return drawer, nil
}

// List retrieves cash drawers with filters
func (s *CashDrawerService) List(ctx context.Context, orgID uuid.UUID, filters CashDrawerFilters) ([]CashDrawer, error) {
	drawers, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list cash drawers", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return drawers, nil
}

// Update updates a cash drawer
func (s *CashDrawerService) Update(ctx context.Context, drawer *CashDrawer) error {
	// Validate drawer
	if err := s.validateDrawer(drawer); err != nil {
		return err
	}

	// Check if drawer exists
	existing, err := s.repo.Get(ctx, drawer.OrganizationID, drawer.ID)
	if err != nil {
		s.logger.Error("failed to get cash drawer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("cash_drawer")
	}

	// Check for duplicate code if code changed
	if drawer.DrawerCode != existing.DrawerCode {
		dup, err := s.repo.GetByCode(ctx, drawer.OrganizationID, drawer.DrawerCode)
		if err != nil {
			s.logger.Error("failed to check duplicate drawer code", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if dup != nil && dup.ID != drawer.ID {
			return apperrors.AlreadyExists("cash_drawer", "Drawer code already exists")
		}
	}

	drawer.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, drawer); err != nil {
		s.logger.Error("failed to update cash drawer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete soft-deletes a cash drawer
func (s *CashDrawerService) Delete(ctx context.Context, orgID, drawerID uuid.UUID) error {
	drawer, err := s.repo.Get(ctx, orgID, drawerID)
	if err != nil {
		s.logger.Error("failed to get cash drawer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if drawer == nil {
		return apperrors.NotFound("cash_drawer")
	}

	if err := s.repo.Delete(ctx, orgID, drawerID); err != nil {
		s.logger.Error("failed to delete cash drawer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validateDrawer validates a cash drawer
func (s *CashDrawerService) validateDrawer(drawer *CashDrawer) error {
	if drawer.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}
	if drawer.DrawerCode == "" {
		return apperrors.ValidationFailed("drawer_code is required")
	}
	if drawer.DrawerName == "" {
		return apperrors.ValidationFailed("drawer_name is required")
	}
	if drawer.LocationID == uuid.Nil {
		return apperrors.ValidationFailed("location_id is required")
	}
	return nil
}

// CashMovementService handles cash movement business logic
type CashMovementService struct {
	movementRepo CashMovementRepository
	sessionRepo  POSSessionRepository
	logger       *logging.Logger
}

// NewCashMovementService creates a new cash movement service
func NewCashMovementService(movementRepo CashMovementRepository, sessionRepo POSSessionRepository, logger *logging.Logger) *CashMovementService {
	return &CashMovementService{
		movementRepo: movementRepo,
		sessionRepo:  sessionRepo,
		logger:       logger,
	}
}

// RecordMovement records a cash movement
func (s *CashMovementService) RecordMovement(ctx context.Context, movement *CashMovement) error {
	// Validate movement
	if err := s.validateMovement(movement); err != nil {
		return err
	}

	// Verify session exists and is open
	session, err := s.sessionRepo.Get(ctx, movement.OrganizationID, movement.POSSessionID)
	if err != nil {
		s.logger.Error("failed to get session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if session == nil {
		return apperrors.NotFound("session")
	}
	if session.Status != SessionStatusOpen {
		return apperrors.ValidationFailed("session must be open to record movements")
	}

	// Initialize movement
	movement.ID = uuid.New()
	movement.CreatedAt = time.Now()

	// Determine if approval is required based on movement type and amount
	if movement.Amount > 1000 && (movement.MovementType == MovementTypePayOut || movement.MovementType == MovementTypeDropToSafe) {
		movement.RequiresApproval = true
	}

	if err := s.movementRepo.Create(ctx, movement); err != nil {
		s.logger.Error("failed to record cash movement", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("Cash movement recorded",
		zap.String("movement_id", movement.ID.String()),
		zap.String("type", string(movement.MovementType)),
		zap.Float64("amount", movement.Amount),
	)

	return nil
}

// ApproveMovement approves a pending cash movement
func (s *CashMovementService) ApproveMovement(ctx context.Context, orgID, movementID, approvedBy uuid.UUID) error {
	movement, err := s.movementRepo.Get(ctx, orgID, movementID)
	if err != nil {
		s.logger.Error("failed to get movement", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if movement == nil {
		return apperrors.NotFound("movement")
	}

	if !movement.RequiresApproval || movement.ApprovedBy != nil {
		return apperrors.ValidationFailed("movement does not require approval or is already approved")
	}

	if err := s.movementRepo.ApproveMovement(ctx, orgID, movementID, approvedBy); err != nil {
		s.logger.Error("failed to approve movement", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("Cash movement approved", zap.String("movement_id", movementID.String()))
	return nil
}

// GetMovements retrieves movements for a session
func (s *CashMovementService) GetMovements(ctx context.Context, orgID, sessionID uuid.UUID) ([]CashMovement, error) {
	movements, err := s.movementRepo.GetForSession(ctx, orgID, sessionID)
	if err != nil {
		s.logger.Error("failed to get movements", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return movements, nil
}

// GetSessionBalance calculates the total balance for a session
func (s *CashMovementService) GetSessionBalance(ctx context.Context, orgID, sessionID uuid.UUID) (float64, error) {
	balance, err := s.movementRepo.CalculateSessionBalance(ctx, orgID, sessionID)
	if err != nil {
		s.logger.Error("failed to calculate session balance", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return balance, nil
}

// List retrieves movements with filters
func (s *CashMovementService) List(ctx context.Context, orgID uuid.UUID, filters CashMovementFilters) ([]CashMovement, error) {
	movements, err := s.movementRepo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list movements", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return movements, nil
}

// validateMovement validates a cash movement
func (s *CashMovementService) validateMovement(movement *CashMovement) error {
	if movement.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}
	if movement.POSSessionID == uuid.Nil {
		return apperrors.ValidationFailed("pos_session_id is required")
	}
	if movement.UserID == uuid.Nil {
		return apperrors.ValidationFailed("user_id is required")
	}
	if movement.Amount == 0 {
		return apperrors.ValidationFailed("amount must not be zero")
	}
	if movement.ReasonDescription == "" {
		return apperrors.ValidationFailed("reason_description is required")
	}

	switch movement.MovementType {
	case MovementTypePayIn, MovementTypePayOut, MovementTypeFloatAdd, MovementTypeDropToSafe:
		// Valid
	default:
		return apperrors.ValidationFailed("invalid movement type")
	}

	return nil
}

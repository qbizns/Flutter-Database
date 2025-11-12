package staff

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles staff/HR business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new staff service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// ==================== EmployeeSchedule Operations ====================

// CreateSchedule creates a new employee schedule
func (s *Service) CreateSchedule(ctx context.Context, schedule *EmployeeSchedule) error {
	if err := s.validateSchedule(schedule); err != nil {
		return err
	}

	schedule.ID = uuid.New()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()
	schedule.Status = "scheduled"

	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		s.logger.Error("failed to create schedule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetSchedule retrieves an employee schedule by ID
func (s *Service) GetSchedule(ctx context.Context, orgID, scheduleID uuid.UUID) (*EmployeeSchedule, error) {
	schedule, err := s.repo.GetSchedule(ctx, orgID, scheduleID)
	if err != nil {
		s.logger.Error("failed to get schedule", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if schedule == nil {
		return nil, apperrors.NotFound("schedule")
	}

	return schedule, nil
}

// ListSchedules lists employee schedules with filters
func (s *Service) ListSchedules(ctx context.Context, orgID uuid.UUID, filters EmployeeScheduleFilters) ([]EmployeeSchedule, error) {
	schedules, err := s.repo.ListSchedules(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list schedules", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return schedules, nil
}

// CountSchedules counts employee schedules
func (s *Service) CountSchedules(ctx context.Context, orgID uuid.UUID, filters EmployeeScheduleFilters) (int64, error) {
	count, err := s.repo.CountSchedules(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count schedules", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateSchedule updates an employee schedule
func (s *Service) UpdateSchedule(ctx context.Context, schedule *EmployeeSchedule) error {
	if err := s.validateSchedule(schedule); err != nil {
		return err
	}

	// Check if schedule exists
	existing, err := s.repo.GetSchedule(ctx, schedule.OrganizationID, schedule.ID)
	if err != nil {
		s.logger.Error("failed to get schedule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("schedule")
	}

	schedule.UpdatedAt = time.Now()

	if err := s.repo.UpdateSchedule(ctx, schedule); err != nil {
		s.logger.Error("failed to update schedule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteSchedule deletes an employee schedule
func (s *Service) DeleteSchedule(ctx context.Context, orgID, scheduleID uuid.UUID) error {
	if err := s.repo.DeleteSchedule(ctx, orgID, scheduleID); err != nil {
		s.logger.Error("failed to delete schedule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== TimeClockEntry Operations ====================

// CreateClockEntry creates a new time clock entry
func (s *Service) CreateClockEntry(ctx context.Context, entry *TimeClockEntry) error {
	if err := s.validateClockEntry(entry); err != nil {
		return err
	}

	entry.ID = uuid.New()
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	if entry.EntryTimestamp.IsZero() {
		entry.EntryTimestamp = time.Now()
	}

	if err := s.repo.CreateClockEntry(ctx, entry); err != nil {
		s.logger.Error("failed to create clock entry", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetClockEntry retrieves a time clock entry
func (s *Service) GetClockEntry(ctx context.Context, orgID, entryID uuid.UUID) (*TimeClockEntry, error) {
	entry, err := s.repo.GetClockEntry(ctx, orgID, entryID)
	if err != nil {
		s.logger.Error("failed to get clock entry", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if entry == nil {
		return nil, apperrors.NotFound("clock entry")
	}

	return entry, nil
}

// ListClockEntries lists time clock entries
func (s *Service) ListClockEntries(ctx context.Context, orgID uuid.UUID, filters TimeClockEntryFilters) ([]TimeClockEntry, error) {
	entries, err := s.repo.ListClockEntries(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list clock entries", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return entries, nil
}

// CountClockEntries counts time clock entries
func (s *Service) CountClockEntries(ctx context.Context, orgID uuid.UUID, filters TimeClockEntryFilters) (int64, error) {
	count, err := s.repo.CountClockEntries(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count clock entries", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateClockEntry updates a time clock entry
func (s *Service) UpdateClockEntry(ctx context.Context, entry *TimeClockEntry) error {
	if err := s.validateClockEntry(entry); err != nil {
		return err
	}

	existing, err := s.repo.GetClockEntry(ctx, entry.OrganizationID, entry.ID)
	if err != nil {
		s.logger.Error("failed to get clock entry", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("clock entry")
	}

	entry.UpdatedAt = time.Now()

	if err := s.repo.UpdateClockEntry(ctx, entry); err != nil {
		s.logger.Error("failed to update clock entry", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteClockEntry deletes a time clock entry
func (s *Service) DeleteClockEntry(ctx context.Context, orgID, entryID uuid.UUID) error {
	if err := s.repo.DeleteClockEntry(ctx, orgID, entryID); err != nil {
		s.logger.Error("failed to delete clock entry", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== Device Operations ====================

// CreateDevice creates a new device
func (s *Service) CreateDevice(ctx context.Context, device *Device) error {
	if err := s.validateDevice(device); err != nil {
		return err
	}

	device.ID = uuid.New()
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	device.Status = "inactive"

	if err := s.repo.CreateDevice(ctx, device); err != nil {
		s.logger.Error("failed to create device", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetDevice retrieves a device by ID
func (s *Service) GetDevice(ctx context.Context, orgID, deviceID uuid.UUID) (*Device, error) {
	device, err := s.repo.GetDevice(ctx, orgID, deviceID)
	if err != nil {
		s.logger.Error("failed to get device", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if device == nil {
		return nil, apperrors.NotFound("device")
	}

	return device, nil
}

// ListDevices lists devices
func (s *Service) ListDevices(ctx context.Context, orgID uuid.UUID, filters DeviceFilters) ([]Device, error) {
	devices, err := s.repo.ListDevices(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list devices", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return devices, nil
}

// CountDevices counts devices
func (s *Service) CountDevices(ctx context.Context, orgID uuid.UUID, filters DeviceFilters) (int64, error) {
	count, err := s.repo.CountDevices(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count devices", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateDevice updates a device
func (s *Service) UpdateDevice(ctx context.Context, device *Device) error {
	if err := s.validateDevice(device); err != nil {
		return err
	}

	existing, err := s.repo.GetDevice(ctx, device.OrganizationID, device.ID)
	if err != nil {
		s.logger.Error("failed to get device", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("device")
	}

	device.UpdatedAt = time.Now()

	if err := s.repo.UpdateDevice(ctx, device); err != nil {
		s.logger.Error("failed to update device", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteDevice deletes a device
func (s *Service) DeleteDevice(ctx context.Context, orgID, deviceID uuid.UUID) error {
	if err := s.repo.DeleteDevice(ctx, orgID, deviceID); err != nil {
		s.logger.Error("failed to delete device", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// UpdateDeviceHeartbeat updates device heartbeat
func (s *Service) UpdateDeviceHeartbeat(ctx context.Context, deviceID uuid.UUID) error {
	if err := s.repo.UpdateDeviceHeartbeat(ctx, deviceID); err != nil {
		s.logger.Error("failed to update device heartbeat", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== PrinterConfiguration Operations ====================

// CreatePrinterConfig creates a new printer configuration
func (s *Service) CreatePrinterConfig(ctx context.Context, config *PrinterConfiguration) error {
	if err := s.validatePrinterConfig(config); err != nil {
		return err
	}

	config.ID = uuid.New()
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	if err := s.repo.CreatePrinterConfig(ctx, config); err != nil {
		s.logger.Error("failed to create printer config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetPrinterConfig retrieves a printer configuration
func (s *Service) GetPrinterConfig(ctx context.Context, orgID, configID uuid.UUID) (*PrinterConfiguration, error) {
	config, err := s.repo.GetPrinterConfig(ctx, orgID, configID)
	if err != nil {
		s.logger.Error("failed to get printer config", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if config == nil {
		return nil, apperrors.NotFound("printer configuration")
	}

	return config, nil
}

// ListPrinterConfigs lists printer configurations
func (s *Service) ListPrinterConfigs(ctx context.Context, orgID uuid.UUID, filters PrinterConfigurationFilters) ([]PrinterConfiguration, error) {
	configs, err := s.repo.ListPrinterConfigs(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list printer configs", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return configs, nil
}

// CountPrinterConfigs counts printer configurations
func (s *Service) CountPrinterConfigs(ctx context.Context, orgID uuid.UUID, filters PrinterConfigurationFilters) (int64, error) {
	count, err := s.repo.CountPrinterConfigs(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count printer configs", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdatePrinterConfig updates a printer configuration
func (s *Service) UpdatePrinterConfig(ctx context.Context, config *PrinterConfiguration) error {
	if err := s.validatePrinterConfig(config); err != nil {
		return err
	}

	existing, err := s.repo.GetPrinterConfig(ctx, config.OrganizationID, config.ID)
	if err != nil {
		s.logger.Error("failed to get printer config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("printer configuration")
	}

	config.UpdatedAt = time.Now()

	if err := s.repo.UpdatePrinterConfig(ctx, config); err != nil {
		s.logger.Error("failed to update printer config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeletePrinterConfig deletes a printer configuration
func (s *Service) DeletePrinterConfig(ctx context.Context, orgID, configID uuid.UUID) error {
	if err := s.repo.DeletePrinterConfig(ctx, orgID, configID); err != nil {
		s.logger.Error("failed to delete printer config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== TipPool Operations ====================

// CreateTipPool creates a new tip pool
func (s *Service) CreateTipPool(ctx context.Context, pool *TipPool) error {
	if err := s.validateTipPool(pool); err != nil {
		return err
	}

	pool.ID = uuid.New()
	pool.CreatedAt = time.Now()
	pool.UpdatedAt = time.Now()
	pool.IsActive = true

	if err := s.repo.CreateTipPool(ctx, pool); err != nil {
		s.logger.Error("failed to create tip pool", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetTipPool retrieves a tip pool
func (s *Service) GetTipPool(ctx context.Context, orgID, poolID uuid.UUID) (*TipPool, error) {
	pool, err := s.repo.GetTipPool(ctx, orgID, poolID)
	if err != nil {
		s.logger.Error("failed to get tip pool", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if pool == nil {
		return nil, apperrors.NotFound("tip pool")
	}

	return pool, nil
}

// ListTipPools lists tip pools
func (s *Service) ListTipPools(ctx context.Context, orgID uuid.UUID, filters TipPoolFilters) ([]TipPool, error) {
	pools, err := s.repo.ListTipPools(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list tip pools", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return pools, nil
}

// CountTipPools counts tip pools
func (s *Service) CountTipPools(ctx context.Context, orgID uuid.UUID, filters TipPoolFilters) (int64, error) {
	count, err := s.repo.CountTipPools(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count tip pools", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateTipPool updates a tip pool
func (s *Service) UpdateTipPool(ctx context.Context, pool *TipPool) error {
	if err := s.validateTipPool(pool); err != nil {
		return err
	}

	existing, err := s.repo.GetTipPool(ctx, pool.OrganizationID, pool.ID)
	if err != nil {
		s.logger.Error("failed to get tip pool", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("tip pool")
	}

	pool.UpdatedAt = time.Now()

	if err := s.repo.UpdateTipPool(ctx, pool); err != nil {
		s.logger.Error("failed to update tip pool", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteTipPool deletes a tip pool
func (s *Service) DeleteTipPool(ctx context.Context, orgID, poolID uuid.UUID) error {
	if err := s.repo.DeleteTipPool(ctx, orgID, poolID); err != nil {
		s.logger.Error("failed to delete tip pool", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== TipDistribution Operations ====================

// CreateTipDistribution creates a new tip distribution
func (s *Service) CreateTipDistribution(ctx context.Context, dist *TipDistribution) error {
	if err := s.validateTipDistribution(dist); err != nil {
		return err
	}

	dist.ID = uuid.New()
	dist.CreatedAt = time.Now()
	dist.UpdatedAt = time.Now()
	dist.PaymentStatus = "pending"

	if dist.DistributionDate.IsZero() {
		dist.DistributionDate = time.Now()
	}

	if err := s.repo.CreateTipDistribution(ctx, dist); err != nil {
		s.logger.Error("failed to create tip distribution", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetTipDistribution retrieves a tip distribution
func (s *Service) GetTipDistribution(ctx context.Context, orgID, distID uuid.UUID) (*TipDistribution, error) {
	dist, err := s.repo.GetTipDistribution(ctx, orgID, distID)
	if err != nil {
		s.logger.Error("failed to get tip distribution", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if dist == nil {
		return nil, apperrors.NotFound("tip distribution")
	}

	return dist, nil
}

// ListTipDistributions lists tip distributions
func (s *Service) ListTipDistributions(ctx context.Context, orgID uuid.UUID, filters TipDistributionFilters) ([]TipDistribution, error) {
	dists, err := s.repo.ListTipDistributions(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list tip distributions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return dists, nil
}

// CountTipDistributions counts tip distributions
func (s *Service) CountTipDistributions(ctx context.Context, orgID uuid.UUID, filters TipDistributionFilters) (int64, error) {
	count, err := s.repo.CountTipDistributions(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count tip distributions", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateTipDistribution updates a tip distribution
func (s *Service) UpdateTipDistribution(ctx context.Context, dist *TipDistribution) error {
	if err := s.validateTipDistribution(dist); err != nil {
		return err
	}

	existing, err := s.repo.GetTipDistribution(ctx, dist.OrganizationID, dist.ID)
	if err != nil {
		s.logger.Error("failed to get tip distribution", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("tip distribution")
	}

	dist.UpdatedAt = time.Now()

	if err := s.repo.UpdateTipDistribution(ctx, dist); err != nil {
		s.logger.Error("failed to update tip distribution", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteTipDistribution deletes a tip distribution
func (s *Service) DeleteTipDistribution(ctx context.Context, orgID, distID uuid.UUID) error {
	if err := s.repo.DeleteTipDistribution(ctx, orgID, distID); err != nil {
		s.logger.Error("failed to delete tip distribution", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== StaffCommission Operations ====================

// CreateCommission creates a new staff commission
func (s *Service) CreateCommission(ctx context.Context, comm *StaffCommission) error {
	if err := s.validateCommission(comm); err != nil {
		return err
	}

	comm.ID = uuid.New()
	comm.CreatedAt = time.Now()
	comm.UpdatedAt = time.Now()
	comm.Status = "pending"

	if comm.CommissionDate.IsZero() {
		comm.CommissionDate = time.Now()
	}

	if err := s.repo.CreateCommission(ctx, comm); err != nil {
		s.logger.Error("failed to create commission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetCommission retrieves a staff commission
func (s *Service) GetCommission(ctx context.Context, orgID, commID uuid.UUID) (*StaffCommission, error) {
	comm, err := s.repo.GetCommission(ctx, orgID, commID)
	if err != nil {
		s.logger.Error("failed to get commission", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if comm == nil {
		return nil, apperrors.NotFound("commission")
	}

	return comm, nil
}

// ListCommissions lists staff commissions
func (s *Service) ListCommissions(ctx context.Context, orgID uuid.UUID, filters StaffCommissionFilters) ([]StaffCommission, error) {
	comms, err := s.repo.ListCommissions(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list commissions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return comms, nil
}

// CountCommissions counts staff commissions
func (s *Service) CountCommissions(ctx context.Context, orgID uuid.UUID, filters StaffCommissionFilters) (int64, error) {
	count, err := s.repo.CountCommissions(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count commissions", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateCommission updates a staff commission
func (s *Service) UpdateCommission(ctx context.Context, comm *StaffCommission) error {
	if err := s.validateCommission(comm); err != nil {
		return err
	}

	existing, err := s.repo.GetCommission(ctx, comm.OrganizationID, comm.ID)
	if err != nil {
		s.logger.Error("failed to get commission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("commission")
	}

	comm.UpdatedAt = time.Now()

	if err := s.repo.UpdateCommission(ctx, comm); err != nil {
		s.logger.Error("failed to update commission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteCommission deletes a staff commission
func (s *Service) DeleteCommission(ctx context.Context, orgID, commID uuid.UUID) error {
	if err := s.repo.DeleteCommission(ctx, orgID, commID); err != nil {
		s.logger.Error("failed to delete commission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== Shift Operations ====================

// CreateShift creates a new shift
func (s *Service) CreateShift(ctx context.Context, shift *Shift) error {
	if err := s.validateShift(shift); err != nil {
		return err
	}

	shift.ID = uuid.New()
	shift.CreatedAt = time.Now()
	shift.UpdatedAt = time.Now()
	shift.Status = "open"
	shift.TotalSales = 0
	shift.TotalTransactions = 0
	shift.TotalRefunds = 0
	shift.TotalDiscounts = 0

	if shift.StartTime.IsZero() {
		shift.StartTime = time.Now()
	}

	if err := s.repo.CreateShift(ctx, shift); err != nil {
		s.logger.Error("failed to create shift", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetShift retrieves a shift
func (s *Service) GetShift(ctx context.Context, orgID, shiftID uuid.UUID) (*Shift, error) {
	shift, err := s.repo.GetShift(ctx, orgID, shiftID)
	if err != nil {
		s.logger.Error("failed to get shift", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if shift == nil {
		return nil, apperrors.NotFound("shift")
	}

	return shift, nil
}

// ListShifts lists shifts
func (s *Service) ListShifts(ctx context.Context, orgID uuid.UUID, filters ShiftFilters) ([]Shift, error) {
	shifts, err := s.repo.ListShifts(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list shifts", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return shifts, nil
}

// CountShifts counts shifts
func (s *Service) CountShifts(ctx context.Context, orgID uuid.UUID, filters ShiftFilters) (int64, error) {
	count, err := s.repo.CountShifts(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count shifts", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateShift updates a shift
func (s *Service) UpdateShift(ctx context.Context, shift *Shift) error {
	if err := s.validateShift(shift); err != nil {
		return err
	}

	existing, err := s.repo.GetShift(ctx, shift.OrganizationID, shift.ID)
	if err != nil {
		s.logger.Error("failed to get shift", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("shift")
	}

	shift.UpdatedAt = time.Now()

	if err := s.repo.UpdateShift(ctx, shift); err != nil {
		s.logger.Error("failed to update shift", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteShift deletes a shift
func (s *Service) DeleteShift(ctx context.Context, orgID, shiftID uuid.UUID) error {
	if err := s.repo.DeleteShift(ctx, orgID, shiftID); err != nil {
		s.logger.Error("failed to delete shift", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== Expense Operations ====================

// CreateExpense creates a new expense
func (s *Service) CreateExpense(ctx context.Context, expense *Expense) error {
	if err := s.validateExpense(expense); err != nil {
		return err
	}

	expense.ID = uuid.New()
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()
	expense.Status = "pending"

	if expense.ExpenseDate.IsZero() {
		expense.ExpenseDate = time.Now()
	}

	if expense.Currency == "" {
		expense.Currency = "USD"
	}

	if err := s.repo.CreateExpense(ctx, expense); err != nil {
		s.logger.Error("failed to create expense", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetExpense retrieves an expense
func (s *Service) GetExpense(ctx context.Context, orgID, expenseID uuid.UUID) (*Expense, error) {
	expense, err := s.repo.GetExpense(ctx, orgID, expenseID)
	if err != nil {
		s.logger.Error("failed to get expense", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if expense == nil {
		return nil, apperrors.NotFound("expense")
	}

	return expense, nil
}

// ListExpenses lists expenses
func (s *Service) ListExpenses(ctx context.Context, orgID uuid.UUID, filters ExpenseFilters) ([]Expense, error) {
	expenses, err := s.repo.ListExpenses(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list expenses", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return expenses, nil
}

// CountExpenses counts expenses
func (s *Service) CountExpenses(ctx context.Context, orgID uuid.UUID, filters ExpenseFilters) (int64, error) {
	count, err := s.repo.CountExpenses(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count expenses", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// UpdateExpense updates an expense
func (s *Service) UpdateExpense(ctx context.Context, expense *Expense) error {
	if err := s.validateExpense(expense); err != nil {
		return err
	}

	existing, err := s.repo.GetExpense(ctx, expense.OrganizationID, expense.ID)
	if err != nil {
		s.logger.Error("failed to get expense", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if existing == nil {
		return apperrors.NotFound("expense")
	}

	expense.UpdatedAt = time.Now()

	if err := s.repo.UpdateExpense(ctx, expense); err != nil {
		s.logger.Error("failed to update expense", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteExpense deletes an expense
func (s *Service) DeleteExpense(ctx context.Context, orgID, expenseID uuid.UUID) error {
	if err := s.repo.DeleteExpense(ctx, orgID, expenseID); err != nil {
		s.logger.Error("failed to delete expense", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ==================== Validation Methods ====================

func (s *Service) validateSchedule(schedule *EmployeeSchedule) error {
	if schedule.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if schedule.EmployeeID == uuid.Nil {
		return apperrors.ValidationError("employee_id is required")
	}

	if schedule.ScheduleDate.IsZero() {
		return apperrors.ValidationError("schedule_date is required")
	}

	if schedule.ScheduledStartTime == "" {
		return apperrors.ValidationError("scheduled_start_time is required")
	}

	if schedule.ScheduledEndTime == "" {
		return apperrors.ValidationError("scheduled_end_time is required")
	}

	return nil
}

func (s *Service) validateClockEntry(entry *TimeClockEntry) error {
	if entry.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if entry.EmployeeID == uuid.Nil {
		return apperrors.ValidationError("employee_id is required")
	}

	if entry.EntryType == "" {
		return apperrors.ValidationError("entry_type is required")
	}

	return nil
}

func (s *Service) validateDevice(device *Device) error {
	if device.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if device.DeviceCode == "" {
		return apperrors.ValidationError("device_code is required")
	}

	if device.DeviceName == "" {
		return apperrors.ValidationError("device_name is required")
	}

	if device.DeviceType == "" {
		return apperrors.ValidationError("device_type is required")
	}

	return nil
}

func (s *Service) validatePrinterConfig(config *PrinterConfiguration) error {
	if config.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if config.PrinterDeviceID == uuid.Nil {
		return apperrors.ValidationError("printer_device_id is required")
	}

	if config.DocumentType == "" {
		return apperrors.ValidationError("document_type is required")
	}

	return nil
}

func (s *Service) validateTipPool(pool *TipPool) error {
	if pool.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if pool.PoolName == "" {
		return apperrors.ValidationError("pool_name is required")
	}

	if pool.PoolType == "" {
		return apperrors.ValidationError("pool_type is required")
	}

	if pool.DistributionMethod == "" {
		return apperrors.ValidationError("distribution_method is required")
	}

	return nil
}

func (s *Service) validateTipDistribution(dist *TipDistribution) error {
	if dist.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if dist.EmployeeID == uuid.Nil {
		return apperrors.ValidationError("employee_id is required")
	}

	if dist.SourceType == "" {
		return apperrors.ValidationError("source_type is required")
	}

	if dist.TipAmount < 0 {
		return apperrors.ValidationError("tip_amount must be greater than or equal to 0")
	}

	return nil
}

func (s *Service) validateCommission(comm *StaffCommission) error {
	if comm.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if comm.EmployeeID == uuid.Nil {
		return apperrors.ValidationError("employee_id is required")
	}

	if comm.SourceType == "" {
		return apperrors.ValidationError("source_type is required")
	}

	if comm.CommissionAmount < 0 {
		return apperrors.ValidationError("commission_amount must be greater than or equal to 0")
	}

	return nil
}

func (s *Service) validateShift(shift *Shift) error {
	if shift.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if shift.UserID == uuid.Nil {
		return apperrors.ValidationError("user_id is required")
	}

	if shift.ShiftNumber == "" {
		return apperrors.ValidationError("shift_number is required")
	}

	if shift.OpeningCash < 0 {
		return apperrors.ValidationError("opening_cash must be greater than or equal to 0")
	}

	return nil
}

func (s *Service) validateExpense(expense *Expense) error {
	if expense.OrganizationID == uuid.Nil {
		return apperrors.ValidationError("organization_id is required")
	}

	if expense.ExpenseNumber == "" {
		return apperrors.ValidationError("expense_number is required")
	}

	if expense.Category == "" {
		return apperrors.ValidationError("category is required")
	}

	if expense.PayeeName == "" {
		return apperrors.ValidationError("payee_name is required")
	}

	if expense.Amount < 0 {
		return apperrors.ValidationError("amount must be greater than or equal to 0")
	}

	if expense.TotalAmount < 0 {
		return apperrors.ValidationError("total_amount must be greater than or equal to 0")
	}

	return nil
}

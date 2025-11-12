package restaurant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// FloorPlanService handles floor plan operations
type FloorPlanService struct {
	repo   FloorPlanRepository
	logger *logging.Logger
}

func NewFloorPlanService(repo FloorPlanRepository, logger *logging.Logger) *FloorPlanService {
	return &FloorPlanService{repo: repo, logger: logger}
}

func (s *FloorPlanService) List(ctx context.Context, orgID uuid.UUID, filters FloorPlanFilters) ([]FloorPlan, error) {
	floorPlans, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list floor plans", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return floorPlans, nil
}

func (s *FloorPlanService) Count(ctx context.Context, orgID uuid.UUID, filters FloorPlanFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count floor plans", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *FloorPlanService) Create(ctx context.Context, floorPlan *FloorPlan) error {
	if err := s.validateFloorPlan(floorPlan); err != nil {
		return err
	}

	floorPlan.ID = uuid.New()
	floorPlan.CreatedAt = time.Now()
	floorPlan.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, floorPlan); err != nil {
		s.logger.Error("failed to create floor plan", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("floor plan created", zap.String("id", floorPlan.ID.String()), zap.String("name", floorPlan.FloorName))
	return nil
}

func (s *FloorPlanService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*FloorPlan, error) {
	floorPlan, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get floor plan", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if floorPlan == nil {
		return nil, apperrors.NotFound("floor plan")
	}
	return floorPlan, nil
}

func (s *FloorPlanService) Update(ctx context.Context, floorPlan *FloorPlan) error {
	if err := s.validateFloorPlan(floorPlan); err != nil {
		return err
	}

	floorPlan.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, floorPlan); err != nil {
		s.logger.Error("failed to update floor plan", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("floor plan updated", zap.String("id", floorPlan.ID.String()))
	return nil
}

func (s *FloorPlanService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete floor plan", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("floor plan deleted", zap.String("id", id.String()))
	return nil
}

func (s *FloorPlanService) validateFloorPlan(floorPlan *FloorPlan) error {
	if strings.TrimSpace(floorPlan.FloorName) == "" {
		return apperrors.ValidationFailed("floor name is required")
	}
	if floorPlan.FloorLevel < 1 {
		return apperrors.ValidationFailed("floor level must be at least 1")
	}
	return nil
}

// TableSectionService handles table section operations
type TableSectionService struct {
	repo   TableSectionRepository
	logger *logging.Logger
}

func NewTableSectionService(repo TableSectionRepository, logger *logging.Logger) *TableSectionService {
	return &TableSectionService{repo: repo, logger: logger}
}

func (s *TableSectionService) List(ctx context.Context, orgID uuid.UUID, filters TableSectionFilters) ([]TableSection, error) {
	sections, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list table sections", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return sections, nil
}

func (s *TableSectionService) Count(ctx context.Context, orgID uuid.UUID, filters TableSectionFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count table sections", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *TableSectionService) Create(ctx context.Context, section *TableSection) error {
	if err := s.validateSection(section); err != nil {
		return err
	}

	section.ID = uuid.New()
	section.CreatedAt = time.Now()
	section.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, section); err != nil {
		s.logger.Error("failed to create table section", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("table section created", zap.String("id", section.ID.String()), zap.String("name", section.SectionName))
	return nil
}

func (s *TableSectionService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*TableSection, error) {
	section, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get table section", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if section == nil {
		return nil, apperrors.NotFound("table section")
	}
	return section, nil
}

func (s *TableSectionService) Update(ctx context.Context, section *TableSection) error {
	if err := s.validateSection(section); err != nil {
		return err
	}

	section.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, section); err != nil {
		s.logger.Error("failed to update table section", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("table section updated", zap.String("id", section.ID.String()))
	return nil
}

func (s *TableSectionService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete table section", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("table section deleted", zap.String("id", id.String()))
	return nil
}

func (s *TableSectionService) validateSection(section *TableSection) error {
	if strings.TrimSpace(section.SectionName) == "" {
		return apperrors.ValidationFailed("section name is required")
	}
	validTypes := map[string]bool{"regular": true, "vip": true, "outdoor": true, "bar": true, "smoking": true, "non_smoking": true, "private": true, "other": true}
	if !validTypes[section.SectionType] {
		return apperrors.ValidationFailed("invalid section type")
	}
	return nil
}

// RestaurantTableService handles restaurant table operations
type RestaurantTableService struct {
	repo   RestaurantTableRepository
	logger *logging.Logger
}

func NewRestaurantTableService(repo RestaurantTableRepository, logger *logging.Logger) *RestaurantTableService {
	return &RestaurantTableService{repo: repo, logger: logger}
}

func (s *RestaurantTableService) List(ctx context.Context, orgID uuid.UUID, filters RestaurantTableFilters) ([]RestaurantTable, error) {
	tables, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list restaurant tables", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return tables, nil
}

func (s *RestaurantTableService) Count(ctx context.Context, orgID uuid.UUID, filters RestaurantTableFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count restaurant tables", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *RestaurantTableService) Create(ctx context.Context, table *RestaurantTable) error {
	if err := s.validateTable(table); err != nil {
		return err
	}

	table.ID = uuid.New()
	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()
	if table.Status == "" {
		table.Status = "available"
	}

	if err := s.repo.Create(ctx, table); err != nil {
		s.logger.Error("failed to create restaurant table", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("restaurant table created", zap.String("id", table.ID.String()), zap.String("number", table.TableNumber))
	return nil
}

func (s *RestaurantTableService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*RestaurantTable, error) {
	table, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get restaurant table", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if table == nil {
		return nil, apperrors.NotFound("restaurant table")
	}
	return table, nil
}

func (s *RestaurantTableService) GetByNumber(ctx context.Context, orgID uuid.UUID, locationID uuid.UUID, tableNumber string) (*RestaurantTable, error) {
	table, err := s.repo.GetByNumber(ctx, orgID, locationID, tableNumber)
	if err != nil {
		s.logger.Error("failed to get restaurant table by number", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if table == nil {
		return nil, apperrors.NotFound("restaurant table")
	}
	return table, nil
}

func (s *RestaurantTableService) Update(ctx context.Context, table *RestaurantTable) error {
	if err := s.validateTable(table); err != nil {
		return err
	}

	table.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, table); err != nil {
		s.logger.Error("failed to update restaurant table", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("restaurant table updated", zap.String("id", table.ID.String()))
	return nil
}

func (s *RestaurantTableService) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	validStatuses := map[string]bool{"available": true, "occupied": true, "reserved": true, "cleaning": true, "maintenance": true, "unavailable": true}
	if !validStatuses[status] {
		return apperrors.ValidationFailed("invalid table status")
	}

	if err := s.repo.UpdateStatus(ctx, orgID, id, status); err != nil {
		s.logger.Error("failed to update restaurant table status", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("restaurant table status updated", zap.String("id", id.String()), zap.String("status", status))
	return nil
}

func (s *RestaurantTableService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete restaurant table", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("restaurant table deleted", zap.String("id", id.String()))
	return nil
}

func (s *RestaurantTableService) validateTable(table *RestaurantTable) error {
	if strings.TrimSpace(table.TableNumber) == "" {
		return apperrors.ValidationFailed("table number is required")
	}
	if table.MaxCapacity < 1 {
		return apperrors.ValidationFailed("max capacity must be at least 1")
	}
	if table.MinCapacity > table.MaxCapacity {
		return apperrors.ValidationFailed("min capacity cannot exceed max capacity")
	}
	if table.CurrentCovers < 0 || table.CurrentCovers > table.MaxCapacity {
		return apperrors.ValidationFailed("current covers must be between 0 and max capacity")
	}
	validShapes := map[string]bool{"square": true, "round": true, "rectangle": true, "oval": true, "custom": true}
	if !validShapes[table.TableShape] {
		return apperrors.ValidationFailed("invalid table shape")
	}
	validStatuses := map[string]bool{"available": true, "occupied": true, "reserved": true, "cleaning": true, "maintenance": true, "unavailable": true}
	if !validStatuses[table.Status] {
		return apperrors.ValidationFailed("invalid table status")
	}
	return nil
}

// ReservationService handles reservation operations
type ReservationService struct {
	repo   ReservationRepository
	logger *logging.Logger
}

func NewReservationService(repo ReservationRepository, logger *logging.Logger) *ReservationService {
	return &ReservationService{repo: repo, logger: logger}
}

func (s *ReservationService) List(ctx context.Context, orgID uuid.UUID, filters ReservationFilters) ([]Reservation, error) {
	reservations, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list reservations", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return reservations, nil
}

func (s *ReservationService) Count(ctx context.Context, orgID uuid.UUID, filters ReservationFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count reservations", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *ReservationService) Create(ctx context.Context, reservation *Reservation) error {
	if err := s.validateReservation(reservation); err != nil {
		return err
	}

	reservation.ID = uuid.New()
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	if reservation.Status == "" {
		reservation.Status = "pending"
	}

	if err := s.repo.Create(ctx, reservation); err != nil {
		s.logger.Error("failed to create reservation", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("reservation created", zap.String("id", reservation.ID.String()), zap.String("number", reservation.ReservationNumber))
	return nil
}

func (s *ReservationService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Reservation, error) {
	reservation, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get reservation", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if reservation == nil {
		return nil, apperrors.NotFound("reservation")
	}
	return reservation, nil
}

func (s *ReservationService) Update(ctx context.Context, reservation *Reservation) error {
	if err := s.validateReservation(reservation); err != nil {
		return err
	}

	reservation.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, reservation); err != nil {
		s.logger.Error("failed to update reservation", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("reservation updated", zap.String("id", reservation.ID.String()))
	return nil
}

func (s *ReservationService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete reservation", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("reservation deleted", zap.String("id", id.String()))
	return nil
}

func (s *ReservationService) validateReservation(reservation *Reservation) error {
	if strings.TrimSpace(reservation.CustomerName) == "" {
		return apperrors.ValidationFailed("customer name is required")
	}
	if reservation.PartySize < 1 {
		return apperrors.ValidationFailed("party size must be at least 1")
	}
	if reservation.DurationMinutes < 1 {
		return apperrors.ValidationFailed("duration must be at least 1 minute")
	}
	validStatuses := map[string]bool{"pending": true, "confirmed": true, "seated": true, "completed": true, "cancelled": true, "no_show": true}
	if !validStatuses[reservation.Status] {
		return apperrors.ValidationFailed("invalid reservation status")
	}
	return nil
}

// ModifierGroupService handles modifier group operations
type ModifierGroupService struct {
	repo   ModifierGroupRepository
	logger *logging.Logger
}

func NewModifierGroupService(repo ModifierGroupRepository, logger *logging.Logger) *ModifierGroupService {
	return &ModifierGroupService{repo: repo, logger: logger}
}

func (s *ModifierGroupService) List(ctx context.Context, orgID uuid.UUID, filters ModifierGroupFilters) ([]ModifierGroup, error) {
	groups, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list modifier groups", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return groups, nil
}

func (s *ModifierGroupService) Count(ctx context.Context, orgID uuid.UUID, filters ModifierGroupFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count modifier groups", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *ModifierGroupService) Create(ctx context.Context, group *ModifierGroup) error {
	if err := s.validateModifierGroup(group); err != nil {
		return err
	}

	group.ID = uuid.New()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, group); err != nil {
		s.logger.Error("failed to create modifier group", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("modifier group created", zap.String("id", group.ID.String()), zap.String("name", group.GroupName))
	return nil
}

func (s *ModifierGroupService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ModifierGroup, error) {
	group, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get modifier group", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if group == nil {
		return nil, apperrors.NotFound("modifier group")
	}
	return group, nil
}

func (s *ModifierGroupService) Update(ctx context.Context, group *ModifierGroup) error {
	if err := s.validateModifierGroup(group); err != nil {
		return err
	}

	group.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, group); err != nil {
		s.logger.Error("failed to update modifier group", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("modifier group updated", zap.String("id", group.ID.String()))
	return nil
}

func (s *ModifierGroupService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete modifier group", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("modifier group deleted", zap.String("id", id.String()))
	return nil
}

func (s *ModifierGroupService) validateModifierGroup(group *ModifierGroup) error {
	if strings.TrimSpace(group.GroupName) == "" {
		return apperrors.ValidationFailed("group name is required")
	}
	validTypes := map[string]bool{"single": true, "multiple": true, "exact": true}
	if !validTypes[group.SelectionType] {
		return apperrors.ValidationFailed("invalid selection type")
	}
	if group.MinSelections < 0 {
		return apperrors.ValidationFailed("min selections cannot be negative")
	}
	if group.MaxSelections != nil && *group.MaxSelections < group.MinSelections {
		return apperrors.ValidationFailed("max selections must be greater than or equal to min selections")
	}
	return nil
}

// ModifierService handles modifier operations
type ModifierService struct {
	repo   ModifierRepository
	logger *logging.Logger
}

func NewModifierService(repo ModifierRepository, logger *logging.Logger) *ModifierService {
	return &ModifierService{repo: repo, logger: logger}
}

func (s *ModifierService) List(ctx context.Context, orgID uuid.UUID, filters ModifierFilters) ([]Modifier, error) {
	modifiers, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list modifiers", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return modifiers, nil
}

func (s *ModifierService) Count(ctx context.Context, orgID uuid.UUID, filters ModifierFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count modifiers", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *ModifierService) Create(ctx context.Context, modifier *Modifier) error {
	if err := s.validateModifier(modifier); err != nil {
		return err
	}

	modifier.ID = uuid.New()
	modifier.CreatedAt = time.Now()
	modifier.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, modifier); err != nil {
		s.logger.Error("failed to create modifier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("modifier created", zap.String("id", modifier.ID.String()), zap.String("name", modifier.ModifierName))
	return nil
}

func (s *ModifierService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Modifier, error) {
	modifier, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get modifier", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if modifier == nil {
		return nil, apperrors.NotFound("modifier")
	}
	return modifier, nil
}

func (s *ModifierService) Update(ctx context.Context, modifier *Modifier) error {
	if err := s.validateModifier(modifier); err != nil {
		return err
	}

	modifier.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, modifier); err != nil {
		s.logger.Error("failed to update modifier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("modifier updated", zap.String("id", modifier.ID.String()))
	return nil
}

func (s *ModifierService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete modifier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("modifier deleted", zap.String("id", id.String()))
	return nil
}

func (s *ModifierService) validateModifier(modifier *Modifier) error {
	if strings.TrimSpace(modifier.ModifierName) == "" {
		return apperrors.ValidationFailed("modifier name is required")
	}
	validTypes := map[string]bool{"add": true, "multiply": true, "replace": true}
	if !validTypes[modifier.PriceType] {
		return apperrors.ValidationFailed("invalid price type")
	}
	return nil
}

// CourseService handles course operations
type CourseService struct {
	repo   CourseRepository
	logger *logging.Logger
}

func NewCourseService(repo CourseRepository, logger *logging.Logger) *CourseService {
	return &CourseService{repo: repo, logger: logger}
}

func (s *CourseService) List(ctx context.Context, orgID uuid.UUID, filters CourseFilters) ([]Course, error) {
	courses, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list courses", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return courses, nil
}

func (s *CourseService) Count(ctx context.Context, orgID uuid.UUID, filters CourseFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count courses", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *CourseService) Create(ctx context.Context, course *Course) error {
	if err := s.validateCourse(course); err != nil {
		return err
	}

	course.ID = uuid.New()
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, course); err != nil {
		s.logger.Error("failed to create course", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("course created", zap.String("id", course.ID.String()), zap.String("name", course.CourseName))
	return nil
}

func (s *CourseService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Course, error) {
	course, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get course", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if course == nil {
		return nil, apperrors.NotFound("course")
	}
	return course, nil
}

func (s *CourseService) Update(ctx context.Context, course *Course) error {
	if err := s.validateCourse(course); err != nil {
		return err
	}

	course.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, course); err != nil {
		s.logger.Error("failed to update course", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("course updated", zap.String("id", course.ID.String()))
	return nil
}

func (s *CourseService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete course", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("course deleted", zap.String("id", id.String()))
	return nil
}

func (s *CourseService) validateCourse(course *Course) error {
	if strings.TrimSpace(course.CourseName) == "" {
		return apperrors.ValidationFailed("course name is required")
	}
	validTypes := map[string]bool{"appetizer": true, "soup": true, "salad": true, "main": true, "side": true, "dessert": true, "beverage": true, "other": true}
	if !validTypes[course.CourseType] {
		return apperrors.ValidationFailed("invalid course type")
	}
	return nil
}

// KitchenStationService handles kitchen station operations
type KitchenStationService struct {
	repo   KitchenStationRepository
	logger *logging.Logger
}

func NewKitchenStationService(repo KitchenStationRepository, logger *logging.Logger) *KitchenStationService {
	return &KitchenStationService{repo: repo, logger: logger}
}

func (s *KitchenStationService) List(ctx context.Context, orgID uuid.UUID, filters KitchenStationFilters) ([]KitchenStation, error) {
	stations, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list kitchen stations", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return stations, nil
}

func (s *KitchenStationService) Count(ctx context.Context, orgID uuid.UUID, filters KitchenStationFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count kitchen stations", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *KitchenStationService) Create(ctx context.Context, station *KitchenStation) error {
	if err := s.validateStation(station); err != nil {
		return err
	}

	station.ID = uuid.New()
	station.CreatedAt = time.Now()
	station.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, station); err != nil {
		s.logger.Error("failed to create kitchen station", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("kitchen station created", zap.String("id", station.ID.String()), zap.String("code", station.StationCode))
	return nil
}

func (s *KitchenStationService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*KitchenStation, error) {
	station, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get kitchen station", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if station == nil {
		return nil, apperrors.NotFound("kitchen station")
	}
	return station, nil
}

func (s *KitchenStationService) Update(ctx context.Context, station *KitchenStation) error {
	if err := s.validateStation(station); err != nil {
		return err
	}

	station.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, station); err != nil {
		s.logger.Error("failed to update kitchen station", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("kitchen station updated", zap.String("id", station.ID.String()))
	return nil
}

func (s *KitchenStationService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete kitchen station", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("kitchen station deleted", zap.String("id", id.String()))
	return nil
}

func (s *KitchenStationService) validateStation(station *KitchenStation) error {
	if strings.TrimSpace(station.StationName) == "" {
		return apperrors.ValidationFailed("station name is required")
	}
	if strings.TrimSpace(station.StationCode) == "" {
		return apperrors.ValidationFailed("station code is required")
	}
	validTypes := map[string]bool{"kitchen": true, "bar": true, "dessert": true, "prep": true}
	if !validTypes[station.StationType] {
		return apperrors.ValidationFailed("invalid station type")
	}
	return nil
}

// OrderService handles order operations
type OrderService struct {
	repo   OrderRepository
	logger *logging.Logger
}

func NewOrderService(repo OrderRepository, logger *logging.Logger) *OrderService {
	return &OrderService{repo: repo, logger: logger}
}

func (s *OrderService) List(ctx context.Context, orgID uuid.UUID, filters OrderFilters) ([]Order, error) {
	orders, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list orders", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return orders, nil
}

func (s *OrderService) Count(ctx context.Context, orgID uuid.UUID, filters OrderFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count orders", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *OrderService) Create(ctx context.Context, order *Order) error {
	if err := s.validateOrder(order); err != nil {
		return err
	}

	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if order.Status == "" {
		order.Status = "draft"
	}
	if order.OrderType == "" {
		order.OrderType = "dine_in"
	}

	if err := s.repo.Create(ctx, order); err != nil {
		s.logger.Error("failed to create order", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order created", zap.String("id", order.ID.String()), zap.String("number", order.OrderNumber))
	return nil
}

func (s *OrderService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Order, error) {
	order, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get order", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if order == nil {
		return nil, apperrors.NotFound("order")
	}
	return order, nil
}

func (s *OrderService) Update(ctx context.Context, order *Order) error {
	if err := s.validateOrder(order); err != nil {
		return err
	}

	order.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, order); err != nil {
		s.logger.Error("failed to update order", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order updated", zap.String("id", order.ID.String()))
	return nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	validStatuses := map[string]bool{"draft": true, "submitted": true, "sent_to_kitchen": true, "preparing": true, "ready": true, "served": true, "completed": true, "cancelled": true, "on_hold": true}
	if !validStatuses[status] {
		return apperrors.ValidationFailed("invalid order status")
	}

	if err := s.repo.UpdateStatus(ctx, orgID, id, status); err != nil {
		s.logger.Error("failed to update order status", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order status updated", zap.String("id", id.String()), zap.String("status", status))
	return nil
}

func (s *OrderService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete order", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order deleted", zap.String("id", id.String()))
	return nil
}

func (s *OrderService) validateOrder(order *Order) error {
	if strings.TrimSpace(order.OrderNumber) == "" {
		return apperrors.ValidationFailed("order number is required")
	}
	if order.Covers < 1 {
		return apperrors.ValidationFailed("covers must be at least 1")
	}
	if order.TotalAmount < 0 || order.Subtotal < 0 || order.TaxAmount < 0 || order.DiscountAmount < 0 || order.ServiceCharge < 0 {
		return apperrors.ValidationFailed("all amounts must be non-negative")
	}
	validTypes := map[string]bool{"dine_in": true, "takeout": true, "delivery": true, "online": true}
	if !validTypes[order.OrderType] {
		return apperrors.ValidationFailed("invalid order type")
	}
	validStatuses := map[string]bool{"draft": true, "submitted": true, "sent_to_kitchen": true, "preparing": true, "ready": true, "served": true, "completed": true, "cancelled": true, "on_hold": true}
	if !validStatuses[order.Status] {
		return apperrors.ValidationFailed("invalid order status")
	}
	return nil
}

// OrderItemService handles order item operations
type OrderItemService struct {
	repo   OrderItemRepository
	logger *logging.Logger
}

func NewOrderItemService(repo OrderItemRepository, logger *logging.Logger) *OrderItemService {
	return &OrderItemService{repo: repo, logger: logger}
}

func (s *OrderItemService) List(ctx context.Context, orgID uuid.UUID, filters OrderItemFilters) ([]OrderItem, error) {
	items, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list order items", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return items, nil
}

func (s *OrderItemService) Count(ctx context.Context, orgID uuid.UUID, filters OrderItemFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count order items", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *OrderItemService) Create(ctx context.Context, item *OrderItem) error {
	if err := s.validateOrderItem(item); err != nil {
		return err
	}

	item.ID = uuid.New()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	if item.Status == "" {
		item.Status = "pending"
	}

	if err := s.repo.Create(ctx, item); err != nil {
		s.logger.Error("failed to create order item", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order item created", zap.String("id", item.ID.String()), zap.String("name", item.ItemName))
	return nil
}

func (s *OrderItemService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrderItem, error) {
	item, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get order item", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if item == nil {
		return nil, apperrors.NotFound("order item")
	}
	return item, nil
}

func (s *OrderItemService) Update(ctx context.Context, item *OrderItem) error {
	if err := s.validateOrderItem(item); err != nil {
		return err
	}

	item.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, item); err != nil {
		s.logger.Error("failed to update order item", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order item updated", zap.String("id", item.ID.String()))
	return nil
}

func (s *OrderItemService) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	validStatuses := map[string]bool{"pending": true, "fired": true, "acknowledged": true, "preparing": true, "ready": true, "served": true, "cancelled": true, "on_hold": true, "voided": true}
	if !validStatuses[status] {
		return apperrors.ValidationFailed("invalid order item status")
	}

	if err := s.repo.UpdateStatus(ctx, orgID, id, status); err != nil {
		s.logger.Error("failed to update order item status", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order item status updated", zap.String("id", id.String()), zap.String("status", status))
	return nil
}

func (s *OrderItemService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete order item", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("order item deleted", zap.String("id", id.String()))
	return nil
}

func (s *OrderItemService) validateOrderItem(item *OrderItem) error {
	if strings.TrimSpace(item.ItemName) == "" {
		return apperrors.ValidationFailed("item name is required")
	}
	if item.Quantity <= 0 {
		return apperrors.ValidationFailed("quantity must be greater than 0")
	}
	if item.UnitPrice < 0 || item.LineTotal < 0 || item.ModifiersTotal < 0 || item.DiscountAmount < 0 {
		return apperrors.ValidationFailed("all amounts must be non-negative")
	}
	validStatuses := map[string]bool{"pending": true, "fired": true, "acknowledged": true, "preparing": true, "ready": true, "served": true, "cancelled": true, "on_hold": true, "voided": true}
	if !validStatuses[item.Status] {
		return apperrors.ValidationFailed("invalid order item status")
	}
	return nil
}

// KitchenTicketService handles kitchen ticket operations
type KitchenTicketService struct {
	repo   KitchenTicketRepository
	logger *logging.Logger
}

func NewKitchenTicketService(repo KitchenTicketRepository, logger *logging.Logger) *KitchenTicketService {
	return &KitchenTicketService{repo: repo, logger: logger}
}

func (s *KitchenTicketService) List(ctx context.Context, orgID uuid.UUID, filters KitchenTicketFilters) ([]KitchenTicket, error) {
	tickets, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list kitchen tickets", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return tickets, nil
}

func (s *KitchenTicketService) Count(ctx context.Context, orgID uuid.UUID, filters KitchenTicketFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count kitchen tickets", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *KitchenTicketService) Create(ctx context.Context, ticket *KitchenTicket) error {
	if err := s.validateKitchenTicket(ticket); err != nil {
		return err
	}

	ticket.ID = uuid.New()
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()
	if ticket.Status == "" {
		ticket.Status = "new"
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		s.logger.Error("failed to create kitchen ticket", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("kitchen ticket created", zap.String("id", ticket.ID.String()), zap.String("number", ticket.TicketNumber))
	return nil
}

func (s *KitchenTicketService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*KitchenTicket, error) {
	ticket, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get kitchen ticket", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if ticket == nil {
		return nil, apperrors.NotFound("kitchen ticket")
	}
	return ticket, nil
}

func (s *KitchenTicketService) Update(ctx context.Context, ticket *KitchenTicket) error {
	if err := s.validateKitchenTicket(ticket); err != nil {
		return err
	}

	ticket.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, ticket); err != nil {
		s.logger.Error("failed to update kitchen ticket", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("kitchen ticket updated", zap.String("id", ticket.ID.String()))
	return nil
}

func (s *KitchenTicketService) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	validStatuses := map[string]bool{"new": true, "acknowledged": true, "preparing": true, "ready": true, "served": true, "completed": true, "cancelled": true, "bumped": true}
	if !validStatuses[status] {
		return apperrors.ValidationFailed("invalid kitchen ticket status")
	}

	if err := s.repo.UpdateStatus(ctx, orgID, id, status); err != nil {
		s.logger.Error("failed to update kitchen ticket status", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("kitchen ticket status updated", zap.String("id", id.String()), zap.String("status", status))
	return nil
}

func (s *KitchenTicketService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete kitchen ticket", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	s.logger.Info("kitchen ticket deleted", zap.String("id", id.String()))
	return nil
}

func (s *KitchenTicketService) validateKitchenTicket(ticket *KitchenTicket) error {
	if strings.TrimSpace(ticket.TicketNumber) == "" {
		return apperrors.ValidationFailed("ticket number is required")
	}
	validTypes := map[string]bool{"normal": true, "rush": true, "remake": true, "special": true}
	if !validTypes[ticket.TicketType] {
		return apperrors.ValidationFailed("invalid ticket type")
	}
	validStatuses := map[string]bool{"new": true, "acknowledged": true, "preparing": true, "ready": true, "served": true, "completed": true, "cancelled": true, "bumped": true}
	if !validStatuses[ticket.Status] {
		return apperrors.ValidationFailed("invalid kitchen ticket status")
	}
	return nil
}

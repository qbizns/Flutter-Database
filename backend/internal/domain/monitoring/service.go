package monitoring

import (
	"fmt"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	errorRepo  POSErrorLogRepository
	healthRepo SystemHealthRepository
	logger     *logging.Logger
}

func NewService(errorRepo POSErrorLogRepository, healthRepo SystemHealthRepository, logger *logging.Logger) *Service {
	return &Service{
		errorRepo:  errorRepo,
		healthRepo: healthRepo,
		logger:     logger,
	}
}

// POS Error Log operations
func (s *Service) ListErrors(ctx context.Context, orgID *uuid.UUID, filters POSErrorLogFilters) ([]POSErrorLog, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 50
	}
	return s.errorRepo.List(ctx, orgID, filters)
}

func (s *Service) LogError(ctx context.Context, errorLog *POSErrorLog) error {
	if err := s.validateErrorLog(errorLog); err != nil {
		return err
	}

	if errorLog.ID == uuid.Nil {
		errorLog.ID = uuid.New()
	}

	now := time.Now()
	if errorLog.OccurredAt.IsZero() {
		errorLog.OccurredAt = now
	}
	errorLog.CreatedAt = now

	if err := s.errorRepo.Create(ctx, errorLog); err != nil {
		s.logger.Error("failed to create error log", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to create error log"))
	}

	return nil
}

func (s *Service) GetError(ctx context.Context, id uuid.UUID) (*POSErrorLog, error) {
	errorLog, err := s.errorRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if errorLog == nil {
		return nil, apperrors.NotFound("error_log")
	}
	return errorLog, nil
}

func (s *Service) ResolveError(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID, notes string) error {
	if err := s.errorRepo.Resolve(ctx, id, resolvedBy, notes); err != nil {
		s.logger.Error("failed to resolve error", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to resolve error"))
	}
	return nil
}

// System Health operations
func (s *Service) ListHealthChecks(ctx context.Context, orgID *uuid.UUID, filters SystemHealthFilters) ([]SystemHealth, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 50
	}
	return s.healthRepo.List(ctx, orgID, filters)
}

func (s *Service) RecordHealthCheck(ctx context.Context, health *SystemHealth) error {
	if err := s.validateHealthCheck(health); err != nil {
		return err
	}

	if health.ID == uuid.Nil {
		health.ID = uuid.New()
	}

	now := time.Now()
	health.CreatedAt = now
	health.UpdatedAt = now
	health.LastCheckAt = &now

	// Determine status based on metric value and thresholds
	if health.MetricValue != nil && health.ThresholdCritical != nil && *health.MetricValue >= *health.ThresholdCritical {
		health.Status = "unhealthy"
		health.LastFailureAt = &now
	} else if health.MetricValue != nil && health.ThresholdWarning != nil && *health.MetricValue >= *health.ThresholdWarning {
		health.Status = "degraded"
	} else {
		health.Status = "healthy"
		health.LastSuccessAt = &now
	}

	if err := s.healthRepo.Upsert(ctx, health); err != nil {
		s.logger.Error("failed to record health check", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to record health check"))
	}

	return nil
}

func (s *Service) GetHealthCheck(ctx context.Context, orgID *uuid.UUID, checkType, checkName string) (*SystemHealth, error) {
	health, err := s.healthRepo.Get(ctx, orgID, checkType, checkName)
	if err != nil {
		return nil, err
	}
	if health == nil {
		return nil, apperrors.NotFound("health_check")
	}
	return health, nil
}

func (s *Service) GetUnhealthyChecks(ctx context.Context, orgID *uuid.UUID) ([]SystemHealth, error) {
	return s.healthRepo.GetUnhealthyChecks(ctx, orgID)
}

// Validation helpers
func (s *Service) validateErrorLog(errorLog *POSErrorLog) error {
	if errorLog.ErrorLevel == "" {
		return apperrors.ValidationFailed("error_level is required")
	}
	if errorLog.ErrorMessage == "" {
		return apperrors.ValidationFailed("error_message is required")
	}

	validLevels := map[string]bool{
		"debug":    true,
		"info":     true,
		"warning":  true,
		"error":    true,
		"critical": true,
	}
	if !validLevels[errorLog.ErrorLevel] {
		return apperrors.ValidationFailed("invalid error_level")
	}

	return nil
}

func (s *Service) validateHealthCheck(health *SystemHealth) error {
	if health.CheckType == "" {
		return apperrors.ValidationFailed("check_type is required")
	}
	if health.CheckName == "" {
		return apperrors.ValidationFailed("check_name is required")
	}

	validStatuses := map[string]bool{
		"healthy":   true,
		"degraded":  true,
		"unhealthy": true,
		"unknown":   true,
	}
	if health.Status != "" && !validStatuses[health.Status] {
		return apperrors.ValidationFailed("invalid status")
	}

	return nil
}

package integrations

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
	channelRepo SalesChannelRepository
	mappingRepo ExternalOrderMappingRepository
	logger      *logging.Logger
}

func NewService(channelRepo SalesChannelRepository, mappingRepo ExternalOrderMappingRepository, logger *logging.Logger) *Service {
	return &Service{
		channelRepo: channelRepo,
		mappingRepo: mappingRepo,
		logger:      logger,
	}
}

// Sales Channel operations
func (s *Service) ListChannels(ctx context.Context, orgID uuid.UUID, filters SalesChannelFilters) ([]SalesChannel, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}
	return s.channelRepo.List(ctx, orgID, filters)
}

func (s *Service) CreateChannel(ctx context.Context, channel *SalesChannel) error {
	if err := s.validateChannel(channel); err != nil {
		return err
	}

	// Check if channel code already exists
	existing, err := s.channelRepo.GetByCode(ctx, channel.OrganizationID, channel.ChannelCode)
	if err != nil {
		return err
	}
	if existing != nil {
		return apperrors.AlreadyExists("sales_channel", "channel_code already exists")
	}

	if channel.ID == uuid.Nil {
		channel.ID = uuid.New()
	}

	now := time.Now()
	channel.CreatedAt = now
	channel.UpdatedAt = now

	if err := s.channelRepo.Create(ctx, channel); err != nil {
		s.logger.Error("failed to create sales channel", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to create sales channel"))
	}

	return nil
}

func (s *Service) GetChannel(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*SalesChannel, error) {
	channel, err := s.channelRepo.Get(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if channel == nil {
		return nil, apperrors.NotFound("sales_channel")
	}
	return channel, nil
}

func (s *Service) UpdateChannel(ctx context.Context, channel *SalesChannel) error {
	if err := s.validateChannel(channel); err != nil {
		return err
	}

	existing, err := s.channelRepo.Get(ctx, channel.OrganizationID, channel.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.NotFound("sales_channel")
	}

	channel.UpdatedAt = time.Now()

	if err := s.channelRepo.Update(ctx, channel); err != nil {
		s.logger.Error("failed to update sales channel", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to update sales channel"))
	}

	return nil
}

func (s *Service) DeleteChannel(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	existing, err := s.channelRepo.Get(ctx, orgID, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.NotFound("sales_channel")
	}

	if err := s.channelRepo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete sales channel", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to delete sales channel"))
	}

	return nil
}

// External Order Mapping operations
func (s *Service) ListMappings(ctx context.Context, orgID uuid.UUID, filters ExternalOrderMappingFilters) ([]ExternalOrderMapping, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}
	return s.mappingRepo.List(ctx, orgID, filters)
}

func (s *Service) CreateMapping(ctx context.Context, mapping *ExternalOrderMapping) error {
	if err := s.validateMapping(mapping); err != nil {
		return err
	}

	// Check if mapping already exists
	existing, err := s.mappingRepo.GetByExternalOrderID(ctx, mapping.OrganizationID, mapping.SalesChannelID, mapping.ExternalOrderID)
	if err != nil {
		return err
	}
	if existing != nil {
		return apperrors.AlreadyExists("external_order_mapping", "external order already mapped")
	}

	if mapping.ID == uuid.Nil {
		mapping.ID = uuid.New()
	}

	now := time.Now()
	mapping.CreatedAt = now
	mapping.UpdatedAt = now

	if err := s.mappingRepo.Create(ctx, mapping); err != nil {
		s.logger.Error("failed to create external order mapping", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to create external order mapping"))
	}

	return nil
}

func (s *Service) GetMapping(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ExternalOrderMapping, error) {
	mapping, err := s.mappingRepo.Get(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if mapping == nil {
		return nil, apperrors.NotFound("external_order_mapping")
	}
	return mapping, nil
}

func (s *Service) GetMappingBySaleID(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) (*ExternalOrderMapping, error) {
	mapping, err := s.mappingRepo.GetBySaleID(ctx, orgID, saleID)
	if err != nil {
		return nil, err
	}
	if mapping == nil {
		return nil, apperrors.NotFound("external_order_mapping")
	}
	return mapping, nil
}

func (s *Service) UpdateMappingSyncStatus(ctx context.Context, id uuid.UUID, status string, lastSyncAt time.Time) error {
	if err := s.mappingRepo.UpdateSyncStatus(ctx, id, status, lastSyncAt); err != nil {
		s.logger.Error("failed to update sync status", zap.Error(err))
		return apperrors.InternalError(fmt.Errorf("failed to update sync status"))
	}
	return nil
}

// Validation helpers
func (s *Service) validateChannel(channel *SalesChannel) error {
	if channel.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}
	if channel.ChannelCode == "" {
		return apperrors.ValidationFailed("channel_code is required")
	}
	if channel.ChannelName == "" {
		return apperrors.ValidationFailed("channel_name is required")
	}
	if channel.ChannelType == "" {
		return apperrors.ValidationFailed("channel_type is required")
	}

	validTypes := map[string]bool{
		"in_store":    true,
		"web":         true,
		"mobile_app":  true,
		"marketplace": true,
		"phone":       true,
		"social":      true,
		"partner":     true,
	}
	if !validTypes[channel.ChannelType] {
		return apperrors.ValidationFailed("invalid channel_type")
	}

	return nil
}

func (s *Service) validateMapping(mapping *ExternalOrderMapping) error {
	if mapping.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}
	if mapping.SaleID == uuid.Nil {
		return apperrors.ValidationFailed("sale_id is required")
	}
	if mapping.SalesChannelID == uuid.Nil {
		return apperrors.ValidationFailed("sales_channel_id is required")
	}
	if mapping.ExternalOrderID == "" {
		return apperrors.ValidationFailed("external_order_id is required")
	}

	validStatuses := map[string]bool{
		"pending":  true,
		"synced":   true,
		"failed":   true,
		"conflict": true,
	}
	if mapping.SyncStatus != "" && !validStatuses[mapping.SyncStatus] {
		return apperrors.ValidationFailed("invalid sync_status")
	}

	return nil
}

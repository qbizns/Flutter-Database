package organization_setting

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for OrganizationSettings
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new OrganizationSettings service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new organization_settings
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateOrganizationSettingsRequest) (*OrganizationSettingsResponse, error) {
	s.logger.Info("creating organization_settings",
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
	entity := &OrganizationSettings{
		OrganizationId: orgID,
		
		Timezone: req.Timezone,
		
		DateFormat: req.DateFormat,
		
		TimeFormat: req.TimeFormat,
		
		NumberFormat: req.NumberFormat,
		
		DefaultCurrency: req.DefaultCurrency,
		
		DefaultLanguage: req.DefaultLanguage,
		
		BusinessType: req.BusinessType,
		
		FiscalYearStart: req.FiscalYearStart,
		
		AutoPrintReceipts: req.AutoPrintReceipts,
		
		AllowNegativeInventory: req.AllowNegativeInventory,
		
		RequireCustomerForSale: req.RequireCustomerForSale,
		
		EnablePriceOverride: req.EnablePriceOverride,
		
		AutoPostSales: req.AutoPostSales,
		
		AutoPostPayments: req.AutoPostPayments,
		
		PostingFrequency: req.PostingFrequency,
		
		SmtpHost: req.SmtpHost,
		
		SmtpPort: req.SmtpPort,
		
		SmtpUsername: req.SmtpUsername,
		
		SmtpUseTls: req.SmtpUseTls,
		
		EmailFromAddress: req.EmailFromAddress,
		
		EmailFromName: req.EmailFromName,
		
		EnableEmailNotifications: req.EnableEmailNotifications,
		
		EnableSmsNotifications: req.EnableSmsNotifications,
		
		Require2fa: req.Require2fa,
		
		SessionTimeoutMinutes: req.SessionTimeoutMinutes,
		
		PasswordMinLength: req.PasswordMinLength,
		
		PasswordRequireSpecial: req.PasswordRequireSpecial,
		
		ApiEnabled: req.ApiEnabled,
		
		ApiRateLimitPerMinute: req.ApiRateLimitPerMinute,
		
		WebhookRetryMaxAttempts: req.WebhookRetryMaxAttempts,
		
		Features: req.Features,
		
		CustomSettings: req.CustomSettings,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create organization_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created organization_settings",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a organization_settings by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrganizationSettingsResponse, error) {
	s.logger.Debug("getting organization_settings",
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
		return nil, fmt.Errorf("failed to get organization_settings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("organization_settings not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of organization_settings records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*OrganizationSettingsListResponse, error) {
	s.logger.Debug("listing organization_settings",
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
		return nil, fmt.Errorf("failed to list organization_settings: %w", err)
	}

	// Convert to response
	items := make([]*OrganizationSettingsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &OrganizationSettingsListResponse{
		Items: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing organization_settings
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateOrganizationSettingsRequest) (*OrganizationSettingsResponse, error) {
	s.logger.Info("updating organization_settings",
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
		return nil, fmt.Errorf("failed to get organization_settings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("organization_settings not found or access denied")
	}
	

	// Update fields
	
	if req.Timezone != nil {
		entity.Timezone = *req.Timezone
	}
	
	if req.DateFormat != nil {
		entity.DateFormat = *req.DateFormat
	}
	
	if req.TimeFormat != nil {
		entity.TimeFormat = *req.TimeFormat
	}
	
	if req.NumberFormat != nil {
		entity.NumberFormat = *req.NumberFormat
	}
	
	if req.DefaultCurrency != nil {
		entity.DefaultCurrency = *req.DefaultCurrency
	}
	
	if req.DefaultLanguage != nil {
		entity.DefaultLanguage = *req.DefaultLanguage
	}
	
	if req.BusinessType != nil {
		entity.BusinessType = *req.BusinessType
	}
	
	if req.FiscalYearStart != nil {
		entity.FiscalYearStart = *req.FiscalYearStart
	}
	
	if req.AutoPrintReceipts != nil {
		entity.AutoPrintReceipts = *req.AutoPrintReceipts
	}
	
	if req.AllowNegativeInventory != nil {
		entity.AllowNegativeInventory = *req.AllowNegativeInventory
	}
	
	if req.RequireCustomerForSale != nil {
		entity.RequireCustomerForSale = *req.RequireCustomerForSale
	}
	
	if req.EnablePriceOverride != nil {
		entity.EnablePriceOverride = *req.EnablePriceOverride
	}
	
	if req.AutoPostSales != nil {
		entity.AutoPostSales = *req.AutoPostSales
	}
	
	if req.AutoPostPayments != nil {
		entity.AutoPostPayments = *req.AutoPostPayments
	}
	
	if req.PostingFrequency != nil {
		entity.PostingFrequency = *req.PostingFrequency
	}
	
	if req.SmtpHost != nil {
		entity.SmtpHost = *req.SmtpHost
	}
	
	if req.SmtpPort != nil {
		entity.SmtpPort = *req.SmtpPort
	}
	
	if req.SmtpUsername != nil {
		entity.SmtpUsername = *req.SmtpUsername
	}
	
	if req.SmtpUseTls != nil {
		entity.SmtpUseTls = *req.SmtpUseTls
	}
	
	if req.EmailFromAddress != nil {
		entity.EmailFromAddress = *req.EmailFromAddress
	}
	
	if req.EmailFromName != nil {
		entity.EmailFromName = *req.EmailFromName
	}
	
	if req.EnableEmailNotifications != nil {
		entity.EnableEmailNotifications = *req.EnableEmailNotifications
	}
	
	if req.EnableSmsNotifications != nil {
		entity.EnableSmsNotifications = *req.EnableSmsNotifications
	}
	
	if req.Require2fa != nil {
		entity.Require2fa = *req.Require2fa
	}
	
	if req.SessionTimeoutMinutes != nil {
		entity.SessionTimeoutMinutes = *req.SessionTimeoutMinutes
	}
	
	if req.PasswordMinLength != nil {
		entity.PasswordMinLength = *req.PasswordMinLength
	}
	
	if req.PasswordRequireSpecial != nil {
		entity.PasswordRequireSpecial = *req.PasswordRequireSpecial
	}
	
	if req.ApiEnabled != nil {
		entity.ApiEnabled = *req.ApiEnabled
	}
	
	if req.ApiRateLimitPerMinute != nil {
		entity.ApiRateLimitPerMinute = *req.ApiRateLimitPerMinute
	}
	
	if req.WebhookRetryMaxAttempts != nil {
		entity.WebhookRetryMaxAttempts = *req.WebhookRetryMaxAttempts
	}
	
	if req.Features != nil {
		entity.Features = *req.Features
	}
	
	if req.CustomSettings != nil {
		entity.CustomSettings = *req.CustomSettings
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update organization_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated organization_settings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a organization_settings
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting organization_settings",
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
		return fmt.Errorf("failed to get organization_settings: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("organization_settings not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete organization_settings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted organization_settings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *OrganizationSettings) *OrganizationSettingsResponse {
	return &OrganizationSettingsResponse{
		
		OrganizationId: entity.OrganizationId,
		
		Timezone: entity.Timezone,
		
		DateFormat: entity.DateFormat,
		
		TimeFormat: entity.TimeFormat,
		
		NumberFormat: entity.NumberFormat,
		
		DefaultCurrency: entity.DefaultCurrency,
		
		DefaultLanguage: entity.DefaultLanguage,
		
		BusinessType: entity.BusinessType,
		
		FiscalYearStart: entity.FiscalYearStart,
		
		AutoPrintReceipts: entity.AutoPrintReceipts,
		
		AllowNegativeInventory: entity.AllowNegativeInventory,
		
		RequireCustomerForSale: entity.RequireCustomerForSale,
		
		EnablePriceOverride: entity.EnablePriceOverride,
		
		AutoPostSales: entity.AutoPostSales,
		
		AutoPostPayments: entity.AutoPostPayments,
		
		PostingFrequency: entity.PostingFrequency,
		
		SmtpHost: entity.SmtpHost,
		
		SmtpPort: entity.SmtpPort,
		
		SmtpUsername: entity.SmtpUsername,
		
		SmtpUseTls: entity.SmtpUseTls,
		
		EmailFromAddress: entity.EmailFromAddress,
		
		EmailFromName: entity.EmailFromName,
		
		EnableEmailNotifications: entity.EnableEmailNotifications,
		
		EnableSmsNotifications: entity.EnableSmsNotifications,
		
		Require2fa: entity.Require2fa,
		
		SessionTimeoutMinutes: entity.SessionTimeoutMinutes,
		
		PasswordMinLength: entity.PasswordMinLength,
		
		PasswordRequireSpecial: entity.PasswordRequireSpecial,
		
		ApiEnabled: entity.ApiEnabled,
		
		ApiRateLimitPerMinute: entity.ApiRateLimitPerMinute,
		
		WebhookRetryMaxAttempts: entity.WebhookRetryMaxAttempts,
		
		Features: entity.Features,
		
		CustomSettings: entity.CustomSettings,
		
		UpdatedAt: entity.UpdatedAt,
		
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


// validateBusinessRules validates business rules for organization_settings
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *OrganizationSettings) error {
	// TODO: Add business rule validation
	// Basic business validation implemented
	// Production: Add module-specific validation rules as needed
	
	// Example validations that can be added:
	// - Duplicate checking within organization
	// - Foreign key validation
	// - Amount/date range validation
	// - Status transition rules
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a organization_settings can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Basic delete validation implemented
	// Production: Add checks for dependent records
	
	// Example checks that can be added:
	// - Query related tables for dependencies
	// - Prevent deletion of entities with transactions
	// - Check business rules (e.g., dont delete active items)
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}

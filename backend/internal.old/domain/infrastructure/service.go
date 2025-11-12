package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// ============================================================================
// Background Job Service
// ============================================================================

// BackgroundJobService handles background job business logic
type BackgroundJobService struct {
	repo   BackgroundJobRepository
	logger *logging.Logger
}

// NewBackgroundJobService creates a new background job service
func NewBackgroundJobService(repo BackgroundJobRepository, logger *logging.Logger) *BackgroundJobService {
	return &BackgroundJobService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of background jobs
func (s *BackgroundJobService) List(ctx context.Context, orgID uuid.UUID, filters BackgroundJobFilters) ([]BackgroundJob, error) {
	jobs, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list background jobs", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return jobs, nil
}

// Create creates a new background job
func (s *BackgroundJobService) Create(ctx context.Context, job *BackgroundJob) error {
	if err := s.validateJob(job); err != nil {
		return err
	}

	job.ID = uuid.New()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	job.Status = "pending"
	job.Attempts = 0

	if err := s.repo.Create(ctx, job); err != nil {
		s.logger.Error("failed to create background job", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a background job by ID
func (s *BackgroundJobService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*BackgroundJob, error) {
	job, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get background job", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if job == nil {
		return nil, apperrors.NotFound("background job")
	}
	return job, nil
}

// Update updates a background job
func (s *BackgroundJobService) Update(ctx context.Context, job *BackgroundJob) error {
	existing, err := s.repo.Get(ctx, job.OrganizationID, job.ID)
	if err != nil {
		s.logger.Error("failed to get background job", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("background job")
	}

	job.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, job); err != nil {
		s.logger.Error("failed to update background job", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes a background job
func (s *BackgroundJobService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	job, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get background job", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if job == nil {
		return apperrors.NotFound("background job")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete background job", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// UpdateStatus updates job status
func (s *BackgroundJobService) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	if !isValidJobStatus(status) {
		return apperrors.ValidationFailed("invalid job status")
	}

	if err := s.repo.UpdateStatus(ctx, orgID, id, status); err != nil {
		s.logger.Error("failed to update job status", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *BackgroundJobService) validateJob(job *BackgroundJob) error {
	if job.JobType == "" {
		return apperrors.ValidationFailed("job_type is required")
	}
	if job.JobName == "" {
		return apperrors.ValidationFailed("job_name is required")
	}
	if len(job.Payload) == 0 {
		return apperrors.ValidationFailed("payload is required")
	}
	return nil
}

func isValidJobStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"completed":  true,
		"failed":     true,
		"cancelled":  true,
		"retrying":   true,
	}
	return validStatuses[status]
}

// ============================================================================
// API Key Service
// ============================================================================

// APIKeyService handles API key business logic
type APIKeyService struct {
	repo   APIKeyRepository
	logger *logging.Logger
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(repo APIKeyRepository, logger *logging.Logger) *APIKeyService {
	return &APIKeyService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves API keys
func (s *APIKeyService) List(ctx context.Context, orgID uuid.UUID, filters APIKeyFilters) ([]APIKey, error) {
	keys, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list API keys", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return keys, nil
}

// Create creates a new API key
func (s *APIKeyService) Create(ctx context.Context, key *APIKey) (string, error) {
	if err := s.validateKey(key); err != nil {
		return "", err
	}

	// Generate key hash
	fullKey := uuid.New().String() + ":" + uuid.New().String()
	hash := sha256.Sum256([]byte(fullKey))
	key.KeyHash = hex.EncodeToString(hash[:])
	key.KeyPrefix = fullKey[:8]

	key.ID = uuid.New()
	key.CreatedAt = time.Now()
	key.UpdatedAt = time.Now()
	key.IsActive = true

	if err := s.repo.Create(ctx, key); err != nil {
		s.logger.Error("failed to create API key", zap.Error(err))
		return "", apperrors.DatabaseError(err)
	}

	return fullKey, nil
}

// Get retrieves an API key
func (s *APIKeyService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*APIKey, error) {
	key, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get API key", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if key == nil {
		return nil, apperrors.NotFound("API key")
	}
	return key, nil
}

// VerifyKey verifies an API key
func (s *APIKeyService) VerifyKey(ctx context.Context, fullKey string) (*APIKey, error) {
	hash := sha256.Sum256([]byte(fullKey))
	keyHash := hex.EncodeToString(hash[:])

	key, err := s.repo.GetByKeyHash(ctx, keyHash)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}
	if key == nil {
		return nil, apperrors.NotFound("API key")
	}

	if !key.IsActive {
		return nil, apperrors.ValidationFailed("API key is inactive")
	}

	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ValidationFailed("API key has expired")
	}

	return key, nil
}

// Update updates an API key
func (s *APIKeyService) Update(ctx context.Context, key *APIKey) error {
	existing, err := s.repo.Get(ctx, key.OrganizationID, key.ID)
	if err != nil {
		s.logger.Error("failed to get API key", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("API key")
	}

	key.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, key); err != nil {
		s.logger.Error("failed to update API key", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes an API key
func (s *APIKeyService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	key, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get API key", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if key == nil {
		return apperrors.NotFound("API key")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete API key", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *APIKeyService) validateKey(key *APIKey) error {
	if key.KeyName == "" {
		return apperrors.ValidationFailed("key_name is required")
	}
	return nil
}

// ============================================================================
// Webhook Service
// ============================================================================

// WebhookService handles webhook business logic
type WebhookService struct {
	repo   WebhookRepository
	logger *logging.Logger
}

// NewWebhookService creates a new webhook service
func NewWebhookService(repo WebhookRepository, logger *logging.Logger) *WebhookService {
	return &WebhookService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves webhooks
func (s *WebhookService) List(ctx context.Context, orgID uuid.UUID, filters WebhookFilters) ([]Webhook, error) {
	webhooks, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list webhooks", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return webhooks, nil
}

// Create creates a new webhook
func (s *WebhookService) Create(ctx context.Context, webhook *Webhook) error {
	if err := s.validateWebhook(webhook); err != nil {
		return err
	}

	webhook.ID = uuid.New()
	webhook.CreatedAt = time.Now()
	webhook.UpdatedAt = time.Now()
	webhook.IsActive = true

	if err := s.repo.Create(ctx, webhook); err != nil {
		s.logger.Error("failed to create webhook", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a webhook
func (s *WebhookService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Webhook, error) {
	webhook, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get webhook", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if webhook == nil {
		return nil, apperrors.NotFound("webhook")
	}
	return webhook, nil
}

// Update updates a webhook
func (s *WebhookService) Update(ctx context.Context, webhook *Webhook) error {
	if err := s.validateWebhook(webhook); err != nil {
		return err
	}

	existing, err := s.repo.Get(ctx, webhook.OrganizationID, webhook.ID)
	if err != nil {
		s.logger.Error("failed to get webhook", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("webhook")
	}

	webhook.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, webhook); err != nil {
		s.logger.Error("failed to update webhook", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes a webhook
func (s *WebhookService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	webhook, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get webhook", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if webhook == nil {
		return apperrors.NotFound("webhook")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete webhook", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *WebhookService) validateWebhook(webhook *Webhook) error {
	if webhook.WebhookName == "" {
		return apperrors.ValidationFailed("webhook_name is required")
	}
	if webhook.URL == "" {
		return apperrors.ValidationFailed("url is required")
	}
	if len(webhook.Events) == 0 {
		return apperrors.ValidationFailed("at least one event is required")
	}
	return nil
}

// ============================================================================
// Notification Service
// ============================================================================

// NotificationService handles notification business logic
type NotificationService struct {
	repo   NotificationRepository
	logger *logging.Logger
}

// NewNotificationService creates a new notification service
func NewNotificationService(repo NotificationRepository, logger *logging.Logger) *NotificationService {
	return &NotificationService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves notifications
func (s *NotificationService) List(ctx context.Context, filters NotificationFilters) ([]Notification, error) {
	notifications, err := s.repo.List(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list notifications", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return notifications, nil
}

// Create creates a new notification
func (s *NotificationService) Create(ctx context.Context, notification *Notification) error {
	if err := s.validateNotification(notification); err != nil {
		return err
	}

	notification.ID = uuid.New()
	notification.CreatedAt = time.Now()
	notification.IsRead = false

	if err := s.repo.Create(ctx, notification); err != nil {
		s.logger.Error("failed to create notification", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a notification
func (s *NotificationService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Notification, error) {
	notification, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get notification", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if notification == nil {
		return nil, apperrors.NotFound("notification")
	}
	return notification, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.MarkAsRead(ctx, orgID, id); err != nil {
		s.logger.Error("failed to mark notification as read", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.MarkAllAsRead(ctx, userID); err != nil {
		s.logger.Error("failed to mark all notifications as read", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *NotificationService) validateNotification(notification *Notification) error {
	if notification.Title == "" {
		return apperrors.ValidationFailed("title is required")
	}
	if notification.Message == "" {
		return apperrors.ValidationFailed("message is required")
	}
	if notification.Category == "" {
		return apperrors.ValidationFailed("category is required")
	}
	return nil
}

// ============================================================================
// File Attachment Service
// ============================================================================

// FileAttachmentService handles file attachment business logic
type FileAttachmentService struct {
	repo   FileAttachmentRepository
	logger *logging.Logger
}

// NewFileAttachmentService creates a new file attachment service
func NewFileAttachmentService(repo FileAttachmentRepository, logger *logging.Logger) *FileAttachmentService {
	return &FileAttachmentService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves file attachments
func (s *FileAttachmentService) List(ctx context.Context, orgID uuid.UUID, filters FileAttachmentFilters) ([]FileAttachment, error) {
	attachments, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list file attachments", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return attachments, nil
}

// Create creates a new file attachment
func (s *FileAttachmentService) Create(ctx context.Context, attachment *FileAttachment) error {
	if err := s.validateAttachment(attachment); err != nil {
		return err
	}

	attachment.ID = uuid.New()
	attachment.CreatedAt = time.Now()

	if err := s.repo.Create(ctx, attachment); err != nil {
		s.logger.Error("failed to create file attachment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a file attachment
func (s *FileAttachmentService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*FileAttachment, error) {
	attachment, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get file attachment", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if attachment == nil {
		return nil, apperrors.NotFound("file attachment")
	}
	return attachment, nil
}

// Update updates a file attachment
func (s *FileAttachmentService) Update(ctx context.Context, attachment *FileAttachment) error {
	existing, err := s.repo.Get(ctx, attachment.OrganizationID, attachment.ID)
	if err != nil {
		s.logger.Error("failed to get file attachment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("file attachment")
	}

	if err := s.repo.Update(ctx, attachment); err != nil {
		s.logger.Error("failed to update file attachment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes a file attachment
func (s *FileAttachmentService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	attachment, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get file attachment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if attachment == nil {
		return apperrors.NotFound("file attachment")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete file attachment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *FileAttachmentService) validateAttachment(attachment *FileAttachment) error {
	if attachment.FileName == "" {
		return apperrors.ValidationFailed("file_name is required")
	}
	if attachment.FileSize <= 0 {
		return apperrors.ValidationFailed("file_size must be greater than 0")
	}
	if attachment.MimeType == "" {
		return apperrors.ValidationFailed("mime_type is required")
	}
	if attachment.StoragePath == "" {
		return apperrors.ValidationFailed("storage_path is required")
	}
	if attachment.EntityType == "" {
		return apperrors.ValidationFailed("entity_type is required")
	}
	return nil
}

// ============================================================================
// Email Queue Service
// ============================================================================

// EmailQueueService handles email queue business logic
type EmailQueueService struct {
	repo   EmailQueueRepository
	logger *logging.Logger
}

// NewEmailQueueService creates a new email queue service
func NewEmailQueueService(repo EmailQueueRepository, logger *logging.Logger) *EmailQueueService {
	return &EmailQueueService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves emails in queue
func (s *EmailQueueService) List(ctx context.Context, orgID *uuid.UUID, filters EmailQueueFilters) ([]EmailQueue, error) {
	emails, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list email queue", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return emails, nil
}

// Create creates a new email queue item
func (s *EmailQueueService) Create(ctx context.Context, email *EmailQueue) error {
	if err := s.validateEmail(email); err != nil {
		return err
	}

	email.ID = uuid.New()
	email.CreatedAt = time.Now()
	email.Status = "pending"
	if email.ScheduledAt.IsZero() {
		email.ScheduledAt = time.Now()
	}

	if err := s.repo.Create(ctx, email); err != nil {
		s.logger.Error("failed to create email queue item", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves an email
func (s *EmailQueueService) Get(ctx context.Context, id uuid.UUID) (*EmailQueue, error) {
	email, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get email", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if email == nil {
		return nil, apperrors.NotFound("email")
	}
	return email, nil
}

// Update updates an email
func (s *EmailQueueService) Update(ctx context.Context, email *EmailQueue) error {
	existing, err := s.repo.Get(ctx, email.ID)
	if err != nil {
		s.logger.Error("failed to get email", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("email")
	}

	if err := s.repo.Update(ctx, email); err != nil {
		s.logger.Error("failed to update email", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *EmailQueueService) validateEmail(email *EmailQueue) error {
	if len(email.ToAddresses) == 0 {
		return apperrors.ValidationFailed("at least one recipient is required")
	}
	if email.Subject == "" {
		return apperrors.ValidationFailed("subject is required")
	}
	if email.BodyHTML == nil && email.BodyText == nil {
		return apperrors.ValidationFailed("body_html or body_text is required")
	}
	return nil
}

// ============================================================================
// SMS Queue Service
// ============================================================================

// SMSQueueService handles SMS queue business logic
type SMSQueueService struct {
	repo   SMSQueueRepository
	logger *logging.Logger
}

// NewSMSQueueService creates a new SMS queue service
func NewSMSQueueService(repo SMSQueueRepository, logger *logging.Logger) *SMSQueueService {
	return &SMSQueueService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves SMS in queue
func (s *SMSQueueService) List(ctx context.Context, orgID *uuid.UUID, filters SMSQueueFilters) ([]SMSQueue, error) {
	messages, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list SMS queue", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return messages, nil
}

// Create creates a new SMS queue item
func (s *SMSQueueService) Create(ctx context.Context, sms *SMSQueue) error {
	if err := s.validateSMS(sms); err != nil {
		return err
	}

	sms.ID = uuid.New()
	sms.CreatedAt = time.Now()
	sms.Status = "pending"
	if sms.ScheduledAt.IsZero() {
		sms.ScheduledAt = time.Now()
	}

	if err := s.repo.Create(ctx, sms); err != nil {
		s.logger.Error("failed to create SMS queue item", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves an SMS
func (s *SMSQueueService) Get(ctx context.Context, id uuid.UUID) (*SMSQueue, error) {
	sms, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get SMS", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if sms == nil {
		return nil, apperrors.NotFound("SMS")
	}
	return sms, nil
}

// Update updates an SMS
func (s *SMSQueueService) Update(ctx context.Context, sms *SMSQueue) error {
	existing, err := s.repo.Get(ctx, sms.ID)
	if err != nil {
		s.logger.Error("failed to get SMS", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("SMS")
	}

	if err := s.repo.Update(ctx, sms); err != nil {
		s.logger.Error("failed to update SMS", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *SMSQueueService) validateSMS(sms *SMSQueue) error {
	if sms.ToPhone == "" {
		return apperrors.ValidationFailed("to_phone is required")
	}
	if sms.Message == "" {
		return apperrors.ValidationFailed("message is required")
	}
	return nil
}

// ============================================================================
// User Session Service
// ============================================================================

// UserSessionService handles user session business logic
type UserSessionService struct {
	repo   UserSessionRepository
	logger *logging.Logger
}

// NewUserSessionService creates a new user session service
func NewUserSessionService(repo UserSessionRepository, logger *logging.Logger) *UserSessionService {
	return &UserSessionService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves user sessions
func (s *UserSessionService) List(ctx context.Context, filters UserSessionFilters) ([]UserSession, error) {
	sessions, err := s.repo.List(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list user sessions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return sessions, nil
}

// Create creates a new user session
func (s *UserSessionService) Create(ctx context.Context, session *UserSession) error {
	if err := s.validateSession(session); err != nil {
		return err
	}

	session.ID = uuid.New()
	session.CreatedAt = time.Now()
	session.LastActivityAt = time.Now()
	session.IsActive = true

	if err := s.repo.Create(ctx, session); err != nil {
		s.logger.Error("failed to create user session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a user session
func (s *UserSessionService) Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*UserSession, error) {
	session, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		s.logger.Error("failed to get user session", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if session == nil {
		return nil, apperrors.NotFound("user session")
	}
	return session, nil
}

// GetByToken retrieves a session by token
func (s *UserSessionService) GetByToken(ctx context.Context, token string) (*UserSession, error) {
	session, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		s.logger.Error("failed to get session by token", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if session == nil {
		return nil, apperrors.NotFound("user session")
	}
	return session, nil
}

// Revoke revokes a session
func (s *UserSessionService) Revoke(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if err := s.repo.Revoke(ctx, userID, id); err != nil {
		s.logger.Error("failed to revoke session", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// RevokeAll revokes all sessions for a user
func (s *UserSessionService) RevokeAll(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.RevokeAll(ctx, userID); err != nil {
		s.logger.Error("failed to revoke all sessions", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *UserSessionService) validateSession(session *UserSession) error {
	if session.SessionToken == "" {
		return apperrors.ValidationFailed("session_token is required")
	}
	if session.ExpiresAt.Before(time.Now()) {
		return apperrors.ValidationFailed("session expiration time must be in the future")
	}
	return nil
}

// ============================================================================
// Organization Settings Service
// ============================================================================

// OrganizationSettingsService handles organization settings business logic
type OrganizationSettingsService struct {
	repo   OrganizationSettingsRepository
	logger *logging.Logger
}

// NewOrganizationSettingsService creates a new organization settings service
func NewOrganizationSettingsService(repo OrganizationSettingsRepository, logger *logging.Logger) *OrganizationSettingsService {
	return &OrganizationSettingsService{
		repo:   repo,
		logger: logger,
	}
}

// Get retrieves organization settings
func (s *OrganizationSettingsService) Get(ctx context.Context, orgID uuid.UUID) (*OrganizationSettings, error) {
	settings, err := s.repo.Get(ctx, orgID)
	if err != nil {
		s.logger.Error("failed to get organization settings", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return settings, nil
}

// Create creates organization settings
func (s *OrganizationSettingsService) Create(ctx context.Context, settings *OrganizationSettings) error {
	if err := s.validateSettings(settings); err != nil {
		return err
	}

	settings.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, settings); err != nil {
		s.logger.Error("failed to create organization settings", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Update updates organization settings
func (s *OrganizationSettingsService) Update(ctx context.Context, settings *OrganizationSettings) error {
	existing, err := s.repo.Get(ctx, settings.OrganizationID)
	if err != nil {
		s.logger.Error("failed to get organization settings", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("organization settings")
	}

	settings.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, settings); err != nil {
		s.logger.Error("failed to update organization settings", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *OrganizationSettingsService) validateSettings(settings *OrganizationSettings) error {
	if settings.Timezone == "" {
		return apperrors.ValidationFailed("timezone is required")
	}
	return nil
}

// ============================================================================
// User Settings Service
// ============================================================================

// UserSettingsService handles user settings business logic
type UserSettingsService struct {
	repo   UserSettingsRepository
	logger *logging.Logger
}

// NewUserSettingsService creates a new user settings service
func NewUserSettingsService(repo UserSettingsRepository, logger *logging.Logger) *UserSettingsService {
	return &UserSettingsService{
		repo:   repo,
		logger: logger,
	}
}

// Get retrieves user settings
func (s *UserSettingsService) Get(ctx context.Context, userID uuid.UUID) (*UserSettings, error) {
	settings, err := s.repo.Get(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user settings", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return settings, nil
}

// Create creates user settings
func (s *UserSettingsService) Create(ctx context.Context, settings *UserSettings) error {
	if err := s.validateUserSettings(settings); err != nil {
		return err
	}

	settings.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, settings); err != nil {
		s.logger.Error("failed to create user settings", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Update updates user settings
func (s *UserSettingsService) Update(ctx context.Context, settings *UserSettings) error {
	existing, err := s.repo.Get(ctx, settings.UserID)
	if err != nil {
		s.logger.Error("failed to get user settings", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("user settings")
	}

	settings.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, settings); err != nil {
		s.logger.Error("failed to update user settings", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

func (s *UserSettingsService) validateUserSettings(settings *UserSettings) error {
	if settings.Language == "" {
		return apperrors.ValidationFailed("language is required")
	}
	return nil
}

// ============================================================================
// Integration Config Service
// ============================================================================

// IntegrationConfigService handles integration config business logic
type IntegrationConfigService struct {
	repo   IntegrationConfigRepository
	logger *logging.Logger
}

// NewIntegrationConfigService creates a new integration config service
func NewIntegrationConfigService(repo IntegrationConfigRepository, logger *logging.Logger) *IntegrationConfigService {
	return &IntegrationConfigService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves integration configs
func (s *IntegrationConfigService) List(ctx context.Context, orgID uuid.UUID, filters IntegrationConfigFilters) ([]IntegrationConfig, error) {
	configs, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list integration configs", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return configs, nil
}

// Create creates an integration config
func (s *IntegrationConfigService) Create(ctx context.Context, config *IntegrationConfig) error {
	if err := s.validateConfig(config); err != nil {
		return err
	}

	config.ID = uuid.New()
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, config); err != nil {
		s.logger.Error("failed to create integration config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves an integration config
func (s *IntegrationConfigService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*IntegrationConfig, error) {
	config, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get integration config", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if config == nil {
		return nil, apperrors.NotFound("integration config")
	}
	return config, nil
}

// Update updates an integration config
func (s *IntegrationConfigService) Update(ctx context.Context, config *IntegrationConfig) error {
	if err := s.validateConfig(config); err != nil {
		return err
	}

	existing, err := s.repo.Get(ctx, config.OrganizationID, config.ID)
	if err != nil {
		s.logger.Error("failed to get integration config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("integration config")
	}

	config.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, config); err != nil {
		s.logger.Error("failed to update integration config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes an integration config
func (s *IntegrationConfigService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	config, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get integration config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if config == nil {
		return apperrors.NotFound("integration config")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete integration config", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetActive retrieves active integration configs by type
func (s *IntegrationConfigService) GetActive(ctx context.Context, orgID uuid.UUID, integrationType string) ([]IntegrationConfig, error) {
	configs, err := s.repo.GetActive(ctx, orgID, integrationType)
	if err != nil {
		s.logger.Error("failed to get active integration configs", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return configs, nil
}

func (s *IntegrationConfigService) validateConfig(config *IntegrationConfig) error {
	if config.IntegrationType == "" {
		return apperrors.ValidationFailed("integration_type is required")
	}
	if config.ProviderName == "" {
		return apperrors.ValidationFailed("provider_name is required")
	}
	if len(config.Credentials) == 0 {
		return apperrors.ValidationFailed("credentials are required")
	}
	return nil
}

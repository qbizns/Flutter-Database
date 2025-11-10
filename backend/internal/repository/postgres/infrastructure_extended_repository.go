package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/your-org/pos-backend/internal/domain/infrastructure"
)

// ============================================================================
// FILE ATTACHMENTS REPOSITORY
// ============================================================================

type FileAttachmentRepositoryImpl struct {
	db *DB
}

func NewFileAttachmentRepository(db *DB) infrastructure.FileAttachmentRepository {
	return &FileAttachmentRepositoryImpl{db: db}
}

func (r *FileAttachmentRepositoryImpl) List(ctx context.Context, orgID uuid.UUID, filters infrastructure.FileAttachmentFilters) ([]infrastructure.FileAttachment, error) {
	query := `
		SELECT id, organization_id, file_name, file_size, mime_type, file_extension,
		       storage_provider, storage_path, storage_url, file_hash, entity_type, entity_id,
		       description, tags, is_public, image_width, image_height,
		       virus_scan_status, virus_scan_at, uploaded_by, created_at, deleted_at
		FROM file_attachments
		WHERE organization_id = $1 AND deleted_at IS NULL`

	args := []interface{}{orgID}
	argCount := 1

	if filters.EntityType != nil {
		argCount++
		query += fmt.Sprintf(" AND entity_type = $%d", argCount)
		args = append(args, *filters.EntityType)
	}

	if filters.EntityID != nil {
		argCount++
		query += fmt.Sprintf(" AND entity_id = $%d", argCount)
		args = append(args, *filters.EntityID)
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []infrastructure.FileAttachment
	for rows.Next() {
		var a infrastructure.FileAttachment
		var tags pq.StringArray
		err := rows.Scan(
			&a.ID, &a.OrganizationID, &a.FileName, &a.FileSize, &a.MimeType, &a.FileExtension,
			&a.StorageProvider, &a.StoragePath, &a.StorageURL, &a.FileHash, &a.EntityType, &a.EntityID,
			&a.Description, &tags, &a.IsPublic, &a.ImageWidth, &a.ImageHeight,
			&a.VirusScanStatus, &a.VirusScanAt, &a.UploadedBy, &a.CreatedAt, &a.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		a.Tags = tags
		attachments = append(attachments, a)
	}

	return attachments, rows.Err()
}

func (r *FileAttachmentRepositoryImpl) Count(ctx context.Context, orgID uuid.UUID, filters infrastructure.FileAttachmentFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM file_attachments WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}

	if filters.EntityType != nil {
		query += " AND entity_type = $2"
		args = append(args, *filters.EntityType)
	}
	if filters.EntityID != nil {
		idx := 3
		if filters.EntityType == nil {
			idx = 2
		}
		query += fmt.Sprintf(" AND entity_id = $%d", idx)
		args = append(args, *filters.EntityID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *FileAttachmentRepositoryImpl) Create(ctx context.Context, attachment *infrastructure.FileAttachment) error {
	attachment.ID = uuid.New()
	attachment.CreatedAt = time.Now()

	query := `
		INSERT INTO file_attachments (
			id, organization_id, file_name, file_size, mime_type, file_extension,
			storage_provider, storage_path, storage_url, file_hash, entity_type, entity_id,
			description, tags, is_public, image_width, image_height,
			virus_scan_status, uploaded_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	tags := pq.Array(attachment.Tags)
	err := r.db.Pool.QueryRow(ctx, query,
		attachment.ID, attachment.OrganizationID, attachment.FileName, attachment.FileSize,
		attachment.MimeType, attachment.FileExtension, attachment.StorageProvider,
		attachment.StoragePath, attachment.StorageURL, attachment.FileHash, attachment.EntityType,
		attachment.EntityID, attachment.Description, tags, attachment.IsPublic,
		attachment.ImageWidth, attachment.ImageHeight, infrastructure.VirusScanPending,
		attachment.UploadedBy, attachment.CreatedAt,
	).Scan(&attachment.ID)

	return err
}

func (r *FileAttachmentRepositoryImpl) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.FileAttachment, error) {
	query := `
		SELECT id, organization_id, file_name, file_size, mime_type, file_extension,
		       storage_provider, storage_path, storage_url, file_hash, entity_type, entity_id,
		       description, tags, is_public, image_width, image_height,
		       virus_scan_status, virus_scan_at, uploaded_by, created_at, deleted_at
		FROM file_attachments
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var a infrastructure.FileAttachment
	var tags pq.StringArray

	err := r.db.Pool.QueryRow(ctx, query, id, orgID).Scan(
		&a.ID, &a.OrganizationID, &a.FileName, &a.FileSize, &a.MimeType, &a.FileExtension,
		&a.StorageProvider, &a.StoragePath, &a.StorageURL, &a.FileHash, &a.EntityType, &a.EntityID,
		&a.Description, &tags, &a.IsPublic, &a.ImageWidth, &a.ImageHeight,
		&a.VirusScanStatus, &a.VirusScanAt, &a.UploadedBy, &a.CreatedAt, &a.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	a.Tags = tags
	return &a, nil
}

func (r *FileAttachmentRepositoryImpl) Update(ctx context.Context, attachment *infrastructure.FileAttachment) error {
	query := `
		UPDATE file_attachments
		SET file_name = $1, virus_scan_status = $2, virus_scan_at = $3
		WHERE id = $4 AND organization_id = $5
	`

	result, err := r.db.Pool.Exec(ctx, query,
		attachment.FileName, attachment.VirusScanStatus, attachment.VirusScanAt,
		attachment.ID, attachment.OrganizationID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *FileAttachmentRepositoryImpl) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	query := "UPDATE file_attachments SET deleted_at = $1 WHERE id = $2 AND organization_id = $3"
	result, err := r.db.Pool.Exec(ctx, query, time.Now(), id, orgID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *FileAttachmentRepositoryImpl) GetByHash(ctx context.Context, orgID uuid.UUID, hash string) (*infrastructure.FileAttachment, error) {
	query := `
		SELECT id, organization_id, file_name, file_size, mime_type, file_extension,
		       storage_provider, storage_path, storage_url, file_hash, entity_type, entity_id,
		       description, tags, is_public, image_width, image_height,
		       virus_scan_status, virus_scan_at, uploaded_by, created_at, deleted_at
		FROM file_attachments
		WHERE file_hash = $1 AND organization_id = $2 AND deleted_at IS NULL
		LIMIT 1
	`

	var a infrastructure.FileAttachment
	var tags pq.StringArray

	err := r.db.Pool.QueryRow(ctx, query, hash, orgID).Scan(
		&a.ID, &a.OrganizationID, &a.FileName, &a.FileSize, &a.MimeType, &a.FileExtension,
		&a.StorageProvider, &a.StoragePath, &a.StorageURL, &a.FileHash, &a.EntityType, &a.EntityID,
		&a.Description, &tags, &a.IsPublic, &a.ImageWidth, &a.ImageHeight,
		&a.VirusScanStatus, &a.VirusScanAt, &a.UploadedBy, &a.CreatedAt, &a.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	a.Tags = tags
	return &a, nil
}

// ============================================================================
// EMAIL QUEUE REPOSITORY
// ============================================================================

type EmailQueueRepositoryImpl struct {
	db *DB
}

func NewEmailQueueRepository(db *DB) infrastructure.EmailQueueRepository {
	return &EmailQueueRepositoryImpl{db: db}
}

func (r *EmailQueueRepositoryImpl) List(ctx context.Context, orgID *uuid.UUID, filters infrastructure.EmailQueueFilters) ([]infrastructure.EmailQueue, error) {
	query := `
		SELECT id, organization_id, to_addresses, cc_addresses, bcc_addresses, from_address, reply_to,
		       subject, body_html, body_text, attachment_ids, template_name, template_data,
		       status, provider, provider_message_id, attempts, max_attempts, error_message, priority,
		       scheduled_at, sent_at, failed_at, created_at
		FROM email_queue
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	query += " ORDER BY priority DESC, scheduled_at ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []infrastructure.EmailQueue
	for rows.Next() {
		var e infrastructure.EmailQueue
		var toAddrs, ccAddrs, bccAddrs pq.StringArray
		var attachmentIds []uuid.UUID
		var templateData *json.RawMessage

		err := rows.Scan(
			&e.ID, &e.OrganizationID, &toAddrs, &ccAddrs, &bccAddrs, &e.FromAddress, &e.ReplyTo,
			&e.Subject, &e.BodyHTML, &e.BodyText, &attachmentIds, &e.TemplateName, &templateData,
			&e.Status, &e.Provider, &e.ProviderMessageID, &e.Attempts, &e.MaxAttempts, &e.ErrorMessage,
			&e.Priority, &e.ScheduledAt, &e.SentAt, &e.FailedAt, &e.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		e.ToAddresses = toAddrs
		e.CCAddresses = ccAddrs
		e.BCCAddresses = bccAddrs
		e.AttachmentIds = attachmentIds
		e.TemplateData = templateData
		emails = append(emails, e)
	}

	return emails, rows.Err()
}

func (r *EmailQueueRepositoryImpl) Count(ctx context.Context, orgID *uuid.UUID, filters infrastructure.EmailQueueFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM email_queue WHERE 1=1"
	args := []interface{}{}

	if orgID != nil {
		query += " AND organization_id = $1"
		args = append(args, *orgID)
	}

	if filters.Status != nil {
		idx := len(args) + 1
		query += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, *filters.Status)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *EmailQueueRepositoryImpl) Create(ctx context.Context, email *infrastructure.EmailQueue) error {
	email.ID = uuid.New()
	email.CreatedAt = time.Now()
	if email.ScheduledAt.IsZero() {
		email.ScheduledAt = time.Now()
	}

	query := `
		INSERT INTO email_queue (
			id, organization_id, to_addresses, cc_addresses, bcc_addresses, from_address, reply_to,
			subject, body_html, body_text, attachment_ids, template_name, template_data,
			status, provider, provider_message_id, attempts, max_attempts, error_message, priority,
			scheduled_at, sent_at, failed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`

	toAddrs := pq.Array(email.ToAddresses)
	ccAddrs := pq.Array(email.CCAddresses)
	bccAddrs := pq.Array(email.BCCAddresses)
	attachmentIds := pq.Array(email.AttachmentIds)

	err := r.db.Pool.QueryRow(ctx, query,
		email.ID, email.OrganizationID, toAddrs, ccAddrs, bccAddrs, email.FromAddress, email.ReplyTo,
		email.Subject, email.BodyHTML, email.BodyText, attachmentIds, email.TemplateName, email.TemplateData,
		infrastructure.StatusPending, email.Provider, email.ProviderMessageID, 0, email.MaxAttempts,
		email.ErrorMessage, email.Priority, email.ScheduledAt, email.SentAt, email.FailedAt, email.CreatedAt,
	).Scan(&email.ID)

	return err
}

func (r *EmailQueueRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (*infrastructure.EmailQueue, error) {
	query := `
		SELECT id, organization_id, to_addresses, cc_addresses, bcc_addresses, from_address, reply_to,
		       subject, body_html, body_text, attachment_ids, template_name, template_data,
		       status, provider, provider_message_id, attempts, max_attempts, error_message, priority,
		       scheduled_at, sent_at, failed_at, created_at
		FROM email_queue
		WHERE id = $1
	`

	var e infrastructure.EmailQueue
	var toAddrs, ccAddrs, bccAddrs pq.StringArray
	var attachmentIds []uuid.UUID
	var templateData *json.RawMessage

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.OrganizationID, &toAddrs, &ccAddrs, &bccAddrs, &e.FromAddress, &e.ReplyTo,
		&e.Subject, &e.BodyHTML, &e.BodyText, &attachmentIds, &e.TemplateName, &templateData,
		&e.Status, &e.Provider, &e.ProviderMessageID, &e.Attempts, &e.MaxAttempts, &e.ErrorMessage,
		&e.Priority, &e.ScheduledAt, &e.SentAt, &e.FailedAt, &e.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	e.ToAddresses = toAddrs
	e.CCAddresses = ccAddrs
	e.BCCAddresses = bccAddrs
	e.AttachmentIds = attachmentIds
	e.TemplateData = templateData
	return &e, nil
}

func (r *EmailQueueRepositoryImpl) Update(ctx context.Context, email *infrastructure.EmailQueue) error {
	query := `
		UPDATE email_queue
		SET status = $1, attempts = $2, error_message = $3, provider_message_id = $4, sent_at = $5, failed_at = $6
		WHERE id = $7
	`

	result, err := r.db.Pool.Exec(ctx, query,
		email.Status, email.Attempts, email.ErrorMessage, email.ProviderMessageID,
		email.SentAt, email.FailedAt, email.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *EmailQueueRepositoryImpl) GetPendingEmails(ctx context.Context, limit int) ([]infrastructure.EmailQueue, error) {
	query := `
		SELECT id, organization_id, to_addresses, cc_addresses, bcc_addresses, from_address, reply_to,
		       subject, body_html, body_text, attachment_ids, template_name, template_data,
		       status, provider, provider_message_id, attempts, max_attempts, error_message, priority,
		       scheduled_at, sent_at, failed_at, created_at
		FROM email_queue
		WHERE status = 'pending' AND scheduled_at <= NOW()
		ORDER BY priority DESC, scheduled_at ASC
		LIMIT $1
	`

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []infrastructure.EmailQueue
	for rows.Next() {
		var e infrastructure.EmailQueue
		var toAddrs, ccAddrs, bccAddrs pq.StringArray
		var attachmentIds []uuid.UUID
		var templateData *json.RawMessage

		err := rows.Scan(
			&e.ID, &e.OrganizationID, &toAddrs, &ccAddrs, &bccAddrs, &e.FromAddress, &e.ReplyTo,
			&e.Subject, &e.BodyHTML, &e.BodyText, &attachmentIds, &e.TemplateName, &templateData,
			&e.Status, &e.Provider, &e.ProviderMessageID, &e.Attempts, &e.MaxAttempts, &e.ErrorMessage,
			&e.Priority, &e.ScheduledAt, &e.SentAt, &e.FailedAt, &e.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		e.ToAddresses = toAddrs
		e.CCAddresses = ccAddrs
		e.BCCAddresses = bccAddrs
		e.AttachmentIds = attachmentIds
		e.TemplateData = templateData
		emails = append(emails, e)
	}

	return emails, rows.Err()
}

// ============================================================================
// SMS QUEUE REPOSITORY
// ============================================================================

type SMSQueueRepositoryImpl struct {
	db *DB
}

func NewSMSQueueRepository(db *DB) infrastructure.SMSQueueRepository {
	return &SMSQueueRepositoryImpl{db: db}
}

func (r *SMSQueueRepositoryImpl) List(ctx context.Context, orgID *uuid.UUID, filters infrastructure.SMSQueueFilters) ([]infrastructure.SMSQueue, error) {
	query := `
		SELECT id, organization_id, to_phone, from_phone, message,
		       status, provider, provider_message_id, attempts, max_attempts, error_message,
		       cost_amount, cost_currency, scheduled_at, sent_at, failed_at, created_at
		FROM sms_queue
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	query += " ORDER BY scheduled_at ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var smsList []infrastructure.SMSQueue
	for rows.Next() {
		var s infrastructure.SMSQueue
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ToPhone, &s.FromPhone, &s.Message,
			&s.Status, &s.Provider, &s.ProviderMessageID, &s.Attempts, &s.MaxAttempts, &s.ErrorMessage,
			&s.CostAmount, &s.CostCurrency, &s.ScheduledAt, &s.SentAt, &s.FailedAt, &s.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		smsList = append(smsList, s)
	}

	return smsList, rows.Err()
}

func (r *SMSQueueRepositoryImpl) Count(ctx context.Context, orgID *uuid.UUID, filters infrastructure.SMSQueueFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM sms_queue WHERE 1=1"
	args := []interface{}{}

	if orgID != nil {
		query += " AND organization_id = $1"
		args = append(args, *orgID)
	}

	if filters.Status != nil {
		idx := len(args) + 1
		query += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, *filters.Status)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *SMSQueueRepositoryImpl) Create(ctx context.Context, sms *infrastructure.SMSQueue) error {
	sms.ID = uuid.New()
	sms.CreatedAt = time.Now()
	if sms.ScheduledAt.IsZero() {
		sms.ScheduledAt = time.Now()
	}

	query := `
		INSERT INTO sms_queue (
			id, organization_id, to_phone, from_phone, message,
			status, provider, provider_message_id, attempts, max_attempts, error_message,
			cost_amount, cost_currency, scheduled_at, sent_at, failed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	err := r.db.Pool.QueryRow(ctx, query,
		sms.ID, sms.OrganizationID, sms.ToPhone, sms.FromPhone, sms.Message,
		infrastructure.StatusPending, sms.Provider, sms.ProviderMessageID, 0, sms.MaxAttempts,
		sms.ErrorMessage, sms.CostAmount, sms.CostCurrency, sms.ScheduledAt, sms.SentAt, sms.FailedAt, sms.CreatedAt,
	).Scan(&sms.ID)

	return err
}

func (r *SMSQueueRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (*infrastructure.SMSQueue, error) {
	query := `
		SELECT id, organization_id, to_phone, from_phone, message,
		       status, provider, provider_message_id, attempts, max_attempts, error_message,
		       cost_amount, cost_currency, scheduled_at, sent_at, failed_at, created_at
		FROM sms_queue
		WHERE id = $1
	`

	var s infrastructure.SMSQueue
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.OrganizationID, &s.ToPhone, &s.FromPhone, &s.Message,
		&s.Status, &s.Provider, &s.ProviderMessageID, &s.Attempts, &s.MaxAttempts, &s.ErrorMessage,
		&s.CostAmount, &s.CostCurrency, &s.ScheduledAt, &s.SentAt, &s.FailedAt, &s.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SMSQueueRepositoryImpl) Update(ctx context.Context, sms *infrastructure.SMSQueue) error {
	query := `
		UPDATE sms_queue
		SET status = $1, attempts = $2, error_message = $3, provider_message_id = $4, sent_at = $5, failed_at = $6
		WHERE id = $7
	`

	result, err := r.db.Pool.Exec(ctx, query,
		sms.Status, sms.Attempts, sms.ErrorMessage, sms.ProviderMessageID,
		sms.SentAt, sms.FailedAt, sms.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SMSQueueRepositoryImpl) GetPendingSMS(ctx context.Context, limit int) ([]infrastructure.SMSQueue, error) {
	query := `
		SELECT id, organization_id, to_phone, from_phone, message,
		       status, provider, provider_message_id, attempts, max_attempts, error_message,
		       cost_amount, cost_currency, scheduled_at, sent_at, failed_at, created_at
		FROM sms_queue
		WHERE status = 'pending' AND scheduled_at <= NOW()
		ORDER BY scheduled_at ASC
		LIMIT $1
	`

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var smsList []infrastructure.SMSQueue
	for rows.Next() {
		var s infrastructure.SMSQueue
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ToPhone, &s.FromPhone, &s.Message,
			&s.Status, &s.Provider, &s.ProviderMessageID, &s.Attempts, &s.MaxAttempts, &s.ErrorMessage,
			&s.CostAmount, &s.CostCurrency, &s.ScheduledAt, &s.SentAt, &s.FailedAt, &s.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		smsList = append(smsList, s)
	}

	return smsList, rows.Err()
}

// ============================================================================
// RATE LIMITS REPOSITORY
// ============================================================================

type RateLimitRepositoryImpl struct {
	db *DB
}

func NewRateLimitRepository(db *DB) infrastructure.RateLimitRepository {
	return &RateLimitRepositoryImpl{db: db}
}

func (r *RateLimitRepositoryImpl) Get(ctx context.Context, identifierType, identifierValue, endpointPath, httpMethod string, windowStart time.Time) (*infrastructure.RateLimit, error) {
	query := `
		SELECT id, identifier_type, identifier_value, endpoint_path, http_method,
		       window_start, window_duration_seconds, request_count, allowed_count,
		       is_blocked, blocked_until, first_request_at, last_request_at
		FROM rate_limits
		WHERE identifier_type = $1 AND identifier_value = $2
		  AND endpoint_path = $3 AND http_method = $4 AND window_start = $5
	`

	var rl infrastructure.RateLimit
	err := r.db.Pool.QueryRow(ctx, query, identifierType, identifierValue, endpointPath, httpMethod, windowStart).Scan(
		&rl.ID, &rl.IdentifierType, &rl.IdentifierValue, &rl.EndpointPath, &rl.HTTPMethod,
		&rl.WindowStart, &rl.WindowDurationSeconds, &rl.RequestCount, &rl.AllowedCount,
		&rl.IsBlocked, &rl.BlockedUntil, &rl.FirstRequestAt, &rl.LastRequestAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &rl, nil
}

func (r *RateLimitRepositoryImpl) Create(ctx context.Context, limit *infrastructure.RateLimit) error {
	limit.ID = uuid.New()
	if limit.FirstRequestAt.IsZero() {
		limit.FirstRequestAt = time.Now()
	}
	if limit.LastRequestAt.IsZero() {
		limit.LastRequestAt = time.Now()
	}

	query := `
		INSERT INTO rate_limits (
			id, identifier_type, identifier_value, endpoint_path, http_method,
			window_start, window_duration_seconds, request_count, allowed_count,
			is_blocked, blocked_until, first_request_at, last_request_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	err := r.db.Pool.QueryRow(ctx, query,
		limit.ID, limit.IdentifierType, limit.IdentifierValue, limit.EndpointPath, limit.HTTPMethod,
		limit.WindowStart, limit.WindowDurationSeconds, limit.RequestCount, limit.AllowedCount,
		limit.IsBlocked, limit.BlockedUntil, limit.FirstRequestAt, limit.LastRequestAt,
	).Scan(&limit.ID)

	return err
}

func (r *RateLimitRepositoryImpl) Update(ctx context.Context, limit *infrastructure.RateLimit) error {
	query := `
		UPDATE rate_limits
		SET request_count = $1, is_blocked = $2, blocked_until = $3, last_request_at = $4
		WHERE id = $5
	`

	limit.LastRequestAt = time.Now()

	result, err := r.db.Pool.Exec(ctx, query,
		limit.RequestCount, limit.IsBlocked, limit.BlockedUntil, limit.LastRequestAt, limit.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RateLimitRepositoryImpl) Cleanup(ctx context.Context, before time.Time) error {
	query := "DELETE FROM rate_limits WHERE window_start < $1"
	_, err := r.db.Pool.Exec(ctx, query, before)
	return err
}

// ============================================================================
// ORGANIZATION SETTINGS REPOSITORY
// ============================================================================

type OrganizationSettingsRepositoryImpl struct {
	db *DB
}

func NewOrganizationSettingsRepository(db *DB) infrastructure.OrganizationSettingsRepository {
	return &OrganizationSettingsRepositoryImpl{db: db}
}

func (r *OrganizationSettingsRepositoryImpl) Get(ctx context.Context, orgID uuid.UUID) (*infrastructure.OrganizationSettings, error) {
	query := `
		SELECT organization_id, timezone, date_format, time_format, number_format,
		       default_currency, default_language, business_type, fiscal_year_start,
		       auto_print_receipts, allow_negative_inventory, require_customer_for_sale, enable_price_override,
		       auto_post_sales, auto_post_payments, posting_frequency,
		       smtp_host, smtp_port, smtp_username, smtp_use_tls, email_from_address, email_from_name,
		       enable_email_notifications, enable_sms_notifications,
		       require_2fa, session_timeout_minutes, password_min_length, password_require_special,
		       api_enabled, api_rate_limit_per_minute, webhook_retry_max_attempts,
		       features, custom_settings, updated_at, updated_by
		FROM organization_settings
		WHERE organization_id = $1
	`

	var s infrastructure.OrganizationSettings
	var features, customSettings json.RawMessage

	err := r.db.Pool.QueryRow(ctx, query, orgID).Scan(
		&s.OrganizationID, &s.Timezone, &s.DateFormat, &s.TimeFormat, &s.NumberFormat,
		&s.DefaultCurrency, &s.DefaultLanguage, &s.BusinessType, &s.FiscalYearStart,
		&s.AutoPrintReceipts, &s.AllowNegativeInventory, &s.RequireCustomerForSale, &s.EnablePriceOverride,
		&s.AutoPostSales, &s.AutoPostPayments, &s.PostingFrequency,
		&s.SMTPHost, &s.SMTPPort, &s.SMTPUsername, &s.SMTPUseTLS, &s.EmailFromAddress, &s.EmailFromName,
		&s.EnableEmailNotifications, &s.EnableSMSNotifications,
		&s.Require2FA, &s.SessionTimeoutMinutes, &s.PasswordMinLength, &s.PasswordRequireSpecial,
		&s.APIEnabled, &s.APIRateLimitPerMinute, &s.WebhookRetryMaxAttempts,
		&features, &customSettings, &s.UpdatedAt, &s.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s.Features = features
	s.CustomSettings = customSettings
	return &s, nil
}

func (r *OrganizationSettingsRepositoryImpl) Create(ctx context.Context, settings *infrastructure.OrganizationSettings) error {
	settings.UpdatedAt = time.Now()

	query := `
		INSERT INTO organization_settings (
			organization_id, timezone, date_format, time_format, number_format,
			default_currency, default_language, business_type, fiscal_year_start,
			auto_print_receipts, allow_negative_inventory, require_customer_for_sale, enable_price_override,
			auto_post_sales, auto_post_payments, posting_frequency,
			smtp_host, smtp_port, smtp_username, smtp_use_tls, email_from_address, email_from_name,
			enable_email_notifications, enable_sms_notifications,
			require_2fa, session_timeout_minutes, password_min_length, password_require_special,
			api_enabled, api_rate_limit_per_minute, webhook_retry_max_attempts,
			features, custom_settings, updated_at, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35)
	`

	err := r.db.Pool.QueryRow(ctx, query,
		settings.OrganizationID, settings.Timezone, settings.DateFormat, settings.TimeFormat, settings.NumberFormat,
		settings.DefaultCurrency, settings.DefaultLanguage, settings.BusinessType, settings.FiscalYearStart,
		settings.AutoPrintReceipts, settings.AllowNegativeInventory, settings.RequireCustomerForSale, settings.EnablePriceOverride,
		settings.AutoPostSales, settings.AutoPostPayments, settings.PostingFrequency,
		settings.SMTPHost, settings.SMTPPort, settings.SMTPUsername, settings.SMTPUseTLS, settings.EmailFromAddress, settings.EmailFromName,
		settings.EnableEmailNotifications, settings.EnableSMSNotifications,
		settings.Require2FA, settings.SessionTimeoutMinutes, settings.PasswordMinLength, settings.PasswordRequireSpecial,
		settings.APIEnabled, settings.APIRateLimitPerMinute, settings.WebhookRetryMaxAttempts,
		settings.Features, settings.CustomSettings, settings.UpdatedAt, settings.UpdatedBy,
	).Scan(&settings.OrganizationID)

	return err
}

func (r *OrganizationSettingsRepositoryImpl) Update(ctx context.Context, settings *infrastructure.OrganizationSettings) error {
	settings.UpdatedAt = time.Now()

	query := `
		UPDATE organization_settings
		SET timezone = $1, date_format = $2, time_format = $3, number_format = $4,
		    default_currency = $5, default_language = $6, business_type = $7, fiscal_year_start = $8,
		    auto_print_receipts = $9, allow_negative_inventory = $10, require_customer_for_sale = $11, enable_price_override = $12,
		    auto_post_sales = $13, auto_post_payments = $14, posting_frequency = $15,
		    smtp_host = $16, smtp_port = $17, smtp_username = $18, smtp_use_tls = $19, email_from_address = $20, email_from_name = $21,
		    enable_email_notifications = $22, enable_sms_notifications = $23,
		    require_2fa = $24, session_timeout_minutes = $25, password_min_length = $26, password_require_special = $27,
		    api_enabled = $28, api_rate_limit_per_minute = $29, webhook_retry_max_attempts = $30,
		    features = $31, custom_settings = $32, updated_at = $33, updated_by = $34
		WHERE organization_id = $35
	`

	result, err := r.db.Pool.Exec(ctx, query,
		settings.Timezone, settings.DateFormat, settings.TimeFormat, settings.NumberFormat,
		settings.DefaultCurrency, settings.DefaultLanguage, settings.BusinessType, settings.FiscalYearStart,
		settings.AutoPrintReceipts, settings.AllowNegativeInventory, settings.RequireCustomerForSale, settings.EnablePriceOverride,
		settings.AutoPostSales, settings.AutoPostPayments, settings.PostingFrequency,
		settings.SMTPHost, settings.SMTPPort, settings.SMTPUsername, settings.SMTPUseTLS, settings.EmailFromAddress, settings.EmailFromName,
		settings.EnableEmailNotifications, settings.EnableSMSNotifications,
		settings.Require2FA, settings.SessionTimeoutMinutes, settings.PasswordMinLength, settings.PasswordRequireSpecial,
		settings.APIEnabled, settings.APIRateLimitPerMinute, settings.WebhookRetryMaxAttempts,
		settings.Features, settings.CustomSettings, settings.UpdatedAt, settings.UpdatedBy,
		settings.OrganizationID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ============================================================================
// USER SETTINGS REPOSITORY
// ============================================================================

type UserSettingsRepositoryImpl struct {
	db *DB
}

func NewUserSettingsRepository(db *DB) infrastructure.UserSettingsRepository {
	return &UserSettingsRepositoryImpl{db: db}
}

func (r *UserSettingsRepositoryImpl) Get(ctx context.Context, userID uuid.UUID) (*infrastructure.UserSettings, error) {
	query := `
		SELECT user_id, theme, language, timezone, default_dashboard, dashboard_layout,
		       items_per_page, default_view, desktop_notifications, sound_notifications,
		       default_location_id, quick_actions, custom_preferences, updated_at
		FROM user_settings
		WHERE user_id = $1
	`

	var us infrastructure.UserSettings
	var dashboardLayout, quickActions, customPreferences json.RawMessage

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&us.UserID, &us.Theme, &us.Language, &us.Timezone, &us.DefaultDashboard, &dashboardLayout,
		&us.ItemsPerPage, &us.DefaultView, &us.DesktopNotifications, &us.SoundNotifications,
		&us.DefaultLocationID, &quickActions, &customPreferences, &us.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	us.DashboardLayout = dashboardLayout
	us.QuickActions = quickActions
	us.CustomPreferences = customPreferences
	return &us, nil
}

func (r *UserSettingsRepositoryImpl) Create(ctx context.Context, settings *infrastructure.UserSettings) error {
	settings.UpdatedAt = time.Now()

	query := `
		INSERT INTO user_settings (
			user_id, theme, language, timezone, default_dashboard, dashboard_layout,
			items_per_page, default_view, desktop_notifications, sound_notifications,
			default_location_id, quick_actions, custom_preferences, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	err := r.db.Pool.QueryRow(ctx, query,
		settings.UserID, settings.Theme, settings.Language, settings.Timezone, settings.DefaultDashboard, settings.DashboardLayout,
		settings.ItemsPerPage, settings.DefaultView, settings.DesktopNotifications, settings.SoundNotifications,
		settings.DefaultLocationID, settings.QuickActions, settings.CustomPreferences, settings.UpdatedAt,
	).Scan(&settings.UserID)

	return err
}

func (r *UserSettingsRepositoryImpl) Update(ctx context.Context, settings *infrastructure.UserSettings) error {
	settings.UpdatedAt = time.Now()

	query := `
		UPDATE user_settings
		SET theme = $1, language = $2, timezone = $3, default_dashboard = $4, dashboard_layout = $5,
		    items_per_page = $6, default_view = $7, desktop_notifications = $8, sound_notifications = $9,
		    default_location_id = $10, quick_actions = $11, custom_preferences = $12, updated_at = $13
		WHERE user_id = $14
	`

	result, err := r.db.Pool.Exec(ctx, query,
		settings.Theme, settings.Language, settings.Timezone, settings.DefaultDashboard, settings.DashboardLayout,
		settings.ItemsPerPage, settings.DefaultView, settings.DesktopNotifications, settings.SoundNotifications,
		settings.DefaultLocationID, settings.QuickActions, settings.CustomPreferences, settings.UpdatedAt,
		settings.UserID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ============================================================================
// DATA EXPORT REQUESTS REPOSITORY
// ============================================================================

type DataExportRequestRepositoryImpl struct {
	db *DB
}

func NewDataExportRequestRepository(db *DB) infrastructure.DataExportRequestRepository {
	return &DataExportRequestRepositoryImpl{db: db}
}

func (r *DataExportRequestRepositoryImpl) List(ctx context.Context, orgID uuid.UUID, filters infrastructure.DataExportRequestFilters) ([]infrastructure.DataExportRequest, error) {
	query := `
		SELECT id, organization_id, export_type, export_format, date_from, date_to, filters,
		       status, file_name, file_size, file_path, download_url, download_expires_at,
		       total_records, processed_records, error_message, requested_by, requested_at, started_at, completed_at
		FROM data_export_requests
		WHERE organization_id = $1`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	query += " ORDER BY requested_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []infrastructure.DataExportRequest
	for rows.Next() {
		var der infrastructure.DataExportRequest
		var filters json.RawMessage

		err := rows.Scan(
			&der.ID, &der.OrganizationID, &der.ExportType, &der.ExportFormat, &der.DateFrom, &der.DateTo, &filters,
			&der.Status, &der.FileName, &der.FileSize, &der.FilePath, &der.DownloadURL, &der.DownloadExpiresAt,
			&der.TotalRecords, &der.ProcessedRecords, &der.ErrorMessage, &der.RequestedBy, &der.RequestedAt, &der.StartedAt, &der.CompletedAt,
		)

		if err != nil {
			return nil, err
		}

		der.Filters = filters
		requests = append(requests, der)
	}

	return requests, rows.Err()
}

func (r *DataExportRequestRepositoryImpl) Count(ctx context.Context, orgID uuid.UUID, filters infrastructure.DataExportRequestFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM data_export_requests WHERE organization_id = $1"
	args := []interface{}{orgID}

	if filters.Status != nil {
		query += " AND status = $2"
		args = append(args, *filters.Status)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *DataExportRequestRepositoryImpl) Create(ctx context.Context, request *infrastructure.DataExportRequest) error {
	request.ID = uuid.New()
	request.RequestedAt = time.Now()

	query := `
		INSERT INTO data_export_requests (
			id, organization_id, export_type, export_format, date_from, date_to, filters,
			status, file_name, file_size, file_path, download_url, download_expires_at,
			total_records, processed_records, error_message, requested_by, requested_at, started_at, completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	err := r.db.Pool.QueryRow(ctx, query,
		request.ID, request.OrganizationID, request.ExportType, request.ExportFormat, request.DateFrom, request.DateTo, request.Filters,
		infrastructure.StatusPending, request.FileName, request.FileSize, request.FilePath, request.DownloadURL, request.DownloadExpiresAt,
		request.TotalRecords, request.ProcessedRecords, request.ErrorMessage, request.RequestedBy, request.RequestedAt, request.StartedAt, request.CompletedAt,
	).Scan(&request.ID)

	return err
}

func (r *DataExportRequestRepositoryImpl) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.DataExportRequest, error) {
	query := `
		SELECT id, organization_id, export_type, export_format, date_from, date_to, filters,
		       status, file_name, file_size, file_path, download_url, download_expires_at,
		       total_records, processed_records, error_message, requested_by, requested_at, started_at, completed_at
		FROM data_export_requests
		WHERE id = $1 AND organization_id = $2
	`

	var der infrastructure.DataExportRequest
	var filters json.RawMessage

	err := r.db.Pool.QueryRow(ctx, query, id, orgID).Scan(
		&der.ID, &der.OrganizationID, &der.ExportType, &der.ExportFormat, &der.DateFrom, &der.DateTo, &filters,
		&der.Status, &der.FileName, &der.FileSize, &der.FilePath, &der.DownloadURL, &der.DownloadExpiresAt,
		&der.TotalRecords, &der.ProcessedRecords, &der.ErrorMessage, &der.RequestedBy, &der.RequestedAt, &der.StartedAt, &der.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	der.Filters = filters
	return &der, nil
}

func (r *DataExportRequestRepositoryImpl) Update(ctx context.Context, request *infrastructure.DataExportRequest) error {
	query := `
		UPDATE data_export_requests
		SET status = $1, file_name = $2, file_size = $3, file_path = $4, download_url = $5,
		    download_expires_at = $6, total_records = $7, processed_records = $8, error_message = $9,
		    started_at = $10, completed_at = $11
		WHERE id = $12 AND organization_id = $13
	`

	result, err := r.db.Pool.Exec(ctx, query,
		request.Status, request.FileName, request.FileSize, request.FilePath, request.DownloadURL,
		request.DownloadExpiresAt, request.TotalRecords, request.ProcessedRecords, request.ErrorMessage,
		request.StartedAt, request.CompletedAt, request.ID, request.OrganizationID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *DataExportRequestRepositoryImpl) GetPendingExports(ctx context.Context) ([]infrastructure.DataExportRequest, error) {
	query := `
		SELECT id, organization_id, export_type, export_format, date_from, date_to, filters,
		       status, file_name, file_size, file_path, download_url, download_expires_at,
		       total_records, processed_records, error_message, requested_by, requested_at, started_at, completed_at
		FROM data_export_requests
		WHERE status IN ('pending', 'processing')
		ORDER BY requested_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []infrastructure.DataExportRequest
	for rows.Next() {
		var der infrastructure.DataExportRequest
		var filters json.RawMessage

		err := rows.Scan(
			&der.ID, &der.OrganizationID, &der.ExportType, &der.ExportFormat, &der.DateFrom, &der.DateTo, &filters,
			&der.Status, &der.FileName, &der.FileSize, &der.FilePath, &der.DownloadURL, &der.DownloadExpiresAt,
			&der.TotalRecords, &der.ProcessedRecords, &der.ErrorMessage, &der.RequestedBy, &der.RequestedAt, &der.StartedAt, &der.CompletedAt,
		)

		if err != nil {
			return nil, err
		}

		der.Filters = filters
		requests = append(requests, der)
	}

	return requests, rows.Err()
}

// ============================================================================
// SCHEDULED REPORTS REPOSITORY
// ============================================================================

type ScheduledReportRepositoryImpl struct {
	db *DB
}

func NewScheduledReportRepository(db *DB) infrastructure.ScheduledReportRepository {
	return &ScheduledReportRepositoryImpl{db: db}
}

func (r *ScheduledReportRepositoryImpl) List(ctx context.Context, orgID uuid.UUID, filters infrastructure.ScheduledReportFilters) ([]infrastructure.ScheduledReport, error) {
	query := `
		SELECT id, organization_id, report_name, report_type, schedule_frequency, schedule_day_of_week,
		       schedule_day_of_month, schedule_time, schedule_timezone, report_parameters, delivery_method,
		       delivery_recipients, output_format, is_active, last_run_at, last_run_status, next_run_at,
		       created_by, created_at, updated_at
		FROM scheduled_reports
		WHERE organization_id = $1`

	args := []interface{}{orgID}
	argCount := 1

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY next_run_at ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []infrastructure.ScheduledReport
	for rows.Next() {
		var sr infrastructure.ScheduledReport
		var reportParams json.RawMessage
		var deliveryRecipients pq.StringArray

		err := rows.Scan(
			&sr.ID, &sr.OrganizationID, &sr.ReportName, &sr.ReportType, &sr.ScheduleFrequency, &sr.ScheduleDayOfWeek,
			&sr.ScheduleDayOfMonth, &sr.ScheduleTime, &sr.ScheduleTimezone, &reportParams, &sr.DeliveryMethod,
			&deliveryRecipients, &sr.OutputFormat, &sr.IsActive, &sr.LastRunAt, &sr.LastRunStatus, &sr.NextRunAt,
			&sr.CreatedBy, &sr.CreatedAt, &sr.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		sr.ReportParameters = reportParams
		sr.DeliveryRecipients = deliveryRecipients
		reports = append(reports, sr)
	}

	return reports, rows.Err()
}

func (r *ScheduledReportRepositoryImpl) Count(ctx context.Context, orgID uuid.UUID, filters infrastructure.ScheduledReportFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM scheduled_reports WHERE organization_id = $1"
	args := []interface{}{orgID}

	if filters.IsActive != nil {
		query += " AND is_active = $2"
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *ScheduledReportRepositoryImpl) Create(ctx context.Context, report *infrastructure.ScheduledReport) error {
	report.ID = uuid.New()
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	query := `
		INSERT INTO scheduled_reports (
			id, organization_id, report_name, report_type, schedule_frequency, schedule_day_of_week,
			schedule_day_of_month, schedule_time, schedule_timezone, report_parameters, delivery_method,
			delivery_recipients, output_format, is_active, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	deliveryRecipients := pq.Array(report.DeliveryRecipients)
	err := r.db.Pool.QueryRow(ctx, query,
		report.ID, report.OrganizationID, report.ReportName, report.ReportType, report.ScheduleFrequency,
		report.ScheduleDayOfWeek, report.ScheduleDayOfMonth, report.ScheduleTime, report.ScheduleTimezone,
		report.ReportParameters, report.DeliveryMethod, deliveryRecipients, report.OutputFormat, report.IsActive,
		report.CreatedBy, report.CreatedAt, report.UpdatedAt,
	).Scan(&report.ID)

	return err
}

func (r *ScheduledReportRepositoryImpl) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.ScheduledReport, error) {
	query := `
		SELECT id, organization_id, report_name, report_type, schedule_frequency, schedule_day_of_week,
		       schedule_day_of_month, schedule_time, schedule_timezone, report_parameters, delivery_method,
		       delivery_recipients, output_format, is_active, last_run_at, last_run_status, next_run_at,
		       created_by, created_at, updated_at
		FROM scheduled_reports
		WHERE id = $1 AND organization_id = $2
	`

	var sr infrastructure.ScheduledReport
	var reportParams json.RawMessage
	var deliveryRecipients pq.StringArray

	err := r.db.Pool.QueryRow(ctx, query, id, orgID).Scan(
		&sr.ID, &sr.OrganizationID, &sr.ReportName, &sr.ReportType, &sr.ScheduleFrequency, &sr.ScheduleDayOfWeek,
		&sr.ScheduleDayOfMonth, &sr.ScheduleTime, &sr.ScheduleTimezone, &reportParams, &sr.DeliveryMethod,
		&deliveryRecipients, &sr.OutputFormat, &sr.IsActive, &sr.LastRunAt, &sr.LastRunStatus, &sr.NextRunAt,
		&sr.CreatedBy, &sr.CreatedAt, &sr.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	sr.ReportParameters = reportParams
	sr.DeliveryRecipients = deliveryRecipients
	return &sr, nil
}

func (r *ScheduledReportRepositoryImpl) Update(ctx context.Context, report *infrastructure.ScheduledReport) error {
	report.UpdatedAt = time.Now()

	query := `
		UPDATE scheduled_reports
		SET report_name = $1, schedule_frequency = $2, schedule_day_of_week = $3,
		    schedule_day_of_month = $4, schedule_time = $5, report_parameters = $6,
		    delivery_method = $7, delivery_recipients = $8, output_format = $9, is_active = $10,
		    last_run_at = $11, last_run_status = $12, next_run_at = $13, updated_at = $14
		WHERE id = $15 AND organization_id = $16
	`

	deliveryRecipients := pq.Array(report.DeliveryRecipients)
	result, err := r.db.Pool.Exec(ctx, query,
		report.ReportName, report.ScheduleFrequency, report.ScheduleDayOfWeek,
		report.ScheduleDayOfMonth, report.ScheduleTime, report.ReportParameters,
		report.DeliveryMethod, deliveryRecipients, report.OutputFormat, report.IsActive,
		report.LastRunAt, report.LastRunStatus, report.NextRunAt, report.UpdatedAt,
		report.ID, report.OrganizationID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ScheduledReportRepositoryImpl) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	query := "DELETE FROM scheduled_reports WHERE id = $1 AND organization_id = $2"
	result, err := r.db.Pool.Exec(ctx, query, id, orgID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ScheduledReportRepositoryImpl) GetDueReports(ctx context.Context) ([]infrastructure.ScheduledReport, error) {
	query := `
		SELECT id, organization_id, report_name, report_type, schedule_frequency, schedule_day_of_week,
		       schedule_day_of_month, schedule_time, schedule_timezone, report_parameters, delivery_method,
		       delivery_recipients, output_format, is_active, last_run_at, last_run_status, next_run_at,
		       created_by, created_at, updated_at
		FROM scheduled_reports
		WHERE is_active = true AND next_run_at <= NOW()
		ORDER BY next_run_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []infrastructure.ScheduledReport
	for rows.Next() {
		var sr infrastructure.ScheduledReport
		var reportParams json.RawMessage
		var deliveryRecipients pq.StringArray

		err := rows.Scan(
			&sr.ID, &sr.OrganizationID, &sr.ReportName, &sr.ReportType, &sr.ScheduleFrequency, &sr.ScheduleDayOfWeek,
			&sr.ScheduleDayOfMonth, &sr.ScheduleTime, &sr.ScheduleTimezone, &reportParams, &sr.DeliveryMethod,
			&deliveryRecipients, &sr.OutputFormat, &sr.IsActive, &sr.LastRunAt, &sr.LastRunStatus, &sr.NextRunAt,
			&sr.CreatedBy, &sr.CreatedAt, &sr.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		sr.ReportParameters = reportParams
		sr.DeliveryRecipients = deliveryRecipients
		reports = append(reports, sr)
	}

	return reports, rows.Err()
}

// ============================================================================
// API REQUEST LOGS REPOSITORY
// ============================================================================

type APIRequestLogRepositoryImpl struct {
	db *DB
}

func NewAPIRequestLogRepository(db *DB) infrastructure.APIRequestLogRepository {
	return &APIRequestLogRepositoryImpl{db: db}
}

func (r *APIRequestLogRepositoryImpl) Create(ctx context.Context, log *infrastructure.APIRequestLog) error {
	log.ID = uuid.New()
	log.CreatedAt = time.Now()

	query := `
		INSERT INTO api_request_logs (
			id, request_id, method, path, query_params, user_id, organization_id, api_key_id,
			request_headers, request_body, ip_address, user_agent,
			status_code, response_headers, response_body, duration_ms, error_message, error_stack, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	err := r.db.Pool.QueryRow(ctx, query,
		log.ID, log.RequestID, log.Method, log.Path, log.QueryParams, log.UserID, log.OrganizationID, log.APIKeyID,
		log.RequestHeaders, log.RequestBody, log.IPAddress, log.UserAgent,
		log.StatusCode, log.ResponseHeaders, log.ResponseBody, log.DurationMs, log.ErrorMessage, log.ErrorStack, log.CreatedAt,
	).Scan(&log.ID)

	return err
}

func (r *APIRequestLogRepositoryImpl) List(ctx context.Context, filters infrastructure.APIRequestLogFilters) ([]infrastructure.APIRequestLog, error) {
	query := `
		SELECT id, request_id, method, path, query_params, user_id, organization_id, api_key_id,
		       request_headers, request_body, ip_address, user_agent,
		       status_code, response_headers, response_body, duration_ms, error_message, error_stack, created_at
		FROM api_request_logs
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if filters.Path != nil {
		argCount++
		query += fmt.Sprintf(" AND path LIKE $%d", argCount)
		args = append(args, "%"+*filters.Path+"%")
	}

	if filters.Method != nil {
		argCount++
		query += fmt.Sprintf(" AND method = $%d", argCount)
		args = append(args, *filters.Method)
	}

	if filters.StatusCode != nil {
		argCount++
		query += fmt.Sprintf(" AND status_code = $%d", argCount)
		args = append(args, *filters.StatusCode)
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []infrastructure.APIRequestLog
	for rows.Next() {
		var arl infrastructure.APIRequestLog
		err := rows.Scan(
			&arl.ID, &arl.RequestID, &arl.Method, &arl.Path, &arl.QueryParams, &arl.UserID, &arl.OrganizationID, &arl.APIKeyID,
			&arl.RequestHeaders, &arl.RequestBody, &arl.IPAddress, &arl.UserAgent,
			&arl.StatusCode, &arl.ResponseHeaders, &arl.ResponseBody, &arl.DurationMs, &arl.ErrorMessage, &arl.ErrorStack, &arl.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		logs = append(logs, arl)
	}

	return logs, rows.Err()
}

func (r *APIRequestLogRepositoryImpl) Count(ctx context.Context, filters infrastructure.APIRequestLogFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM api_request_logs WHERE 1=1"
	args := []interface{}{}

	if filters.Path != nil {
		query += " AND path LIKE $1"
		args = append(args, "%"+*filters.Path+"%")
	}

	if filters.Method != nil {
		idx := len(args) + 1
		query += fmt.Sprintf(" AND method = $%d", idx)
		args = append(args, *filters.Method)
	}

	if filters.StatusCode != nil {
		idx := len(args) + 1
		query += fmt.Sprintf(" AND status_code = $%d", idx)
		args = append(args, *filters.StatusCode)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *APIRequestLogRepositoryImpl) GetByRequestID(ctx context.Context, requestID string) (*infrastructure.APIRequestLog, error) {
	query := `
		SELECT id, request_id, method, path, query_params, user_id, organization_id, api_key_id,
		       request_headers, request_body, ip_address, user_agent,
		       status_code, response_headers, response_body, duration_ms, error_message, error_stack, created_at
		FROM api_request_logs
		WHERE request_id = $1
		LIMIT 1
	`

	var arl infrastructure.APIRequestLog
	err := r.db.Pool.QueryRow(ctx, query, requestID).Scan(
		&arl.ID, &arl.RequestID, &arl.Method, &arl.Path, &arl.QueryParams, &arl.UserID, &arl.OrganizationID, &arl.APIKeyID,
		&arl.RequestHeaders, &arl.RequestBody, &arl.IPAddress, &arl.UserAgent,
		&arl.StatusCode, &arl.ResponseHeaders, &arl.ResponseBody, &arl.DurationMs, &arl.ErrorMessage, &arl.ErrorStack, &arl.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &arl, nil
}

func (r *APIRequestLogRepositoryImpl) Cleanup(ctx context.Context, before time.Time) error {
	query := "DELETE FROM api_request_logs WHERE created_at < $1"
	_, err := r.db.Pool.Exec(ctx, query, before)
	return err
}

// ============================================================================
// INTEGRATION CONFIGS REPOSITORY
// ============================================================================

type IntegrationConfigRepositoryImpl struct {
	db *DB
}

func NewIntegrationConfigRepository(db *DB) infrastructure.IntegrationConfigRepository {
	return &IntegrationConfigRepositoryImpl{db: db}
}

func (r *IntegrationConfigRepositoryImpl) List(ctx context.Context, orgID uuid.UUID, filters infrastructure.IntegrationConfigFilters) ([]infrastructure.IntegrationConfig, error) {
	query := `
		SELECT id, organization_id, integration_type, provider_name, credentials, settings,
		       is_active, is_connected, connection_status, last_sync_at, last_sync_status, sync_frequency,
		       webhook_url, webhook_secret, created_by, created_at, updated_at, deleted_at
		FROM integration_configs
		WHERE organization_id = $1 AND deleted_at IS NULL`

	args := []interface{}{orgID}
	argCount := 1

	if filters.IntegrationType != nil {
		argCount++
		query += fmt.Sprintf(" AND integration_type = $%d", argCount)
		args = append(args, *filters.IntegrationType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []infrastructure.IntegrationConfig
	for rows.Next() {
		var ic infrastructure.IntegrationConfig
		var credentials, settings json.RawMessage

		err := rows.Scan(
			&ic.ID, &ic.OrganizationID, &ic.IntegrationType, &ic.ProviderName, &credentials, &settings,
			&ic.IsActive, &ic.IsConnected, &ic.ConnectionStatus, &ic.LastSyncAt, &ic.LastSyncStatus, &ic.SyncFrequency,
			&ic.WebhookURL, &ic.WebhookSecret, &ic.CreatedBy, &ic.CreatedAt, &ic.UpdatedAt, &ic.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		ic.Credentials = credentials
		ic.Settings = settings
		configs = append(configs, ic)
	}

	return configs, rows.Err()
}

func (r *IntegrationConfigRepositoryImpl) Count(ctx context.Context, orgID uuid.UUID, filters infrastructure.IntegrationConfigFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM integration_configs WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}

	if filters.IntegrationType != nil {
		query += " AND integration_type = $2"
		args = append(args, *filters.IntegrationType)
	}

	if filters.IsActive != nil {
		idx := 3
		if filters.IntegrationType == nil {
			idx = 2
		}
		query += fmt.Sprintf(" AND is_active = $%d", idx)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *IntegrationConfigRepositoryImpl) Create(ctx context.Context, config *infrastructure.IntegrationConfig) error {
	config.ID = uuid.New()
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	query := `
		INSERT INTO integration_configs (
			id, organization_id, integration_type, provider_name, credentials, settings,
			is_active, is_connected, connection_status, sync_frequency,
			webhook_url, webhook_secret, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	err := r.db.Pool.QueryRow(ctx, query,
		config.ID, config.OrganizationID, config.IntegrationType, config.ProviderName, config.Credentials, config.Settings,
		config.IsActive, config.IsConnected, config.ConnectionStatus, config.SyncFrequency,
		config.WebhookURL, config.WebhookSecret, config.CreatedBy, config.CreatedAt, config.UpdatedAt,
	).Scan(&config.ID)

	return err
}

func (r *IntegrationConfigRepositoryImpl) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.IntegrationConfig, error) {
	query := `
		SELECT id, organization_id, integration_type, provider_name, credentials, settings,
		       is_active, is_connected, connection_status, last_sync_at, last_sync_status, sync_frequency,
		       webhook_url, webhook_secret, created_by, created_at, updated_at, deleted_at
		FROM integration_configs
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var ic infrastructure.IntegrationConfig
	var credentials, settings json.RawMessage

	err := r.db.Pool.QueryRow(ctx, query, id, orgID).Scan(
		&ic.ID, &ic.OrganizationID, &ic.IntegrationType, &ic.ProviderName, &credentials, &settings,
		&ic.IsActive, &ic.IsConnected, &ic.ConnectionStatus, &ic.LastSyncAt, &ic.LastSyncStatus, &ic.SyncFrequency,
		&ic.WebhookURL, &ic.WebhookSecret, &ic.CreatedBy, &ic.CreatedAt, &ic.UpdatedAt, &ic.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ic.Credentials = credentials
	ic.Settings = settings
	return &ic, nil
}

func (r *IntegrationConfigRepositoryImpl) Update(ctx context.Context, config *infrastructure.IntegrationConfig) error {
	config.UpdatedAt = time.Now()

	query := `
		UPDATE integration_configs
		SET integration_type = $1, provider_name = $2, credentials = $3, settings = $4,
		    is_active = $5, is_connected = $6, connection_status = $7,
		    last_sync_at = $8, last_sync_status = $9, sync_frequency = $10,
		    webhook_url = $11, webhook_secret = $12, updated_at = $13
		WHERE id = $14 AND organization_id = $15
	`

	result, err := r.db.Pool.Exec(ctx, query,
		config.IntegrationType, config.ProviderName, config.Credentials, config.Settings,
		config.IsActive, config.IsConnected, config.ConnectionStatus,
		config.LastSyncAt, config.LastSyncStatus, config.SyncFrequency,
		config.WebhookURL, config.WebhookSecret, config.UpdatedAt,
		config.ID, config.OrganizationID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *IntegrationConfigRepositoryImpl) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	query := "UPDATE integration_configs SET deleted_at = $1 WHERE id = $2 AND organization_id = $3"
	result, err := r.db.Pool.Exec(ctx, query, time.Now(), id, orgID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *IntegrationConfigRepositoryImpl) GetActive(ctx context.Context, orgID uuid.UUID, integrationType string) ([]infrastructure.IntegrationConfig, error) {
	query := `
		SELECT id, organization_id, integration_type, provider_name, credentials, settings,
		       is_active, is_connected, connection_status, last_sync_at, last_sync_status, sync_frequency,
		       webhook_url, webhook_secret, created_by, created_at, updated_at, deleted_at
		FROM integration_configs
		WHERE organization_id = $1 AND integration_type = $2 AND is_active = true AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, integrationType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []infrastructure.IntegrationConfig
	for rows.Next() {
		var ic infrastructure.IntegrationConfig
		var credentials, settings json.RawMessage

		err := rows.Scan(
			&ic.ID, &ic.OrganizationID, &ic.IntegrationType, &ic.ProviderName, &credentials, &settings,
			&ic.IsActive, &ic.IsConnected, &ic.ConnectionStatus, &ic.LastSyncAt, &ic.LastSyncStatus, &ic.SyncFrequency,
			&ic.WebhookURL, &ic.WebhookSecret, &ic.CreatedBy, &ic.CreatedAt, &ic.UpdatedAt, &ic.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		ic.Credentials = credentials
		ic.Settings = settings
		configs = append(configs, ic)
	}

	return configs, rows.Err()
}

// ============================================================================
// IMMUTABILITY VIOLATIONS LOG REPOSITORY
// ============================================================================

type ImmutabilityViolationLogRepository struct {
	db *DB
}

func NewImmutabilityViolationLogRepository(db *DB) *ImmutabilityViolationLogRepository {
	return &ImmutabilityViolationLogRepository{db: db}
}

func (r *ImmutabilityViolationLogRepository) List(ctx context.Context, orgID *uuid.UUID, filters infrastructure.ImmutabilityViolationLogFilters) ([]infrastructure.ImmutabilityViolationLog, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT id, organization_id, table_name, record_id, operation,
		       attempted_by, attempted_at, error_message, blocked_data, metadata
		FROM immutability_violations_log
		WHERE 1=1
	`

	var args []interface{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.TableName != nil {
		argCount++
		query += fmt.Sprintf(" AND table_name = $%d", argCount)
		args = append(args, *filters.TableName)
	}

	if filters.Operation != nil {
		argCount++
		query += fmt.Sprintf(" AND operation = $%d", argCount)
		args = append(args, *filters.Operation)
	}

	if filters.AttemptedBy != nil {
		argCount++
		query += fmt.Sprintf(" AND attempted_by = $%d", argCount)
		args = append(args, *filters.AttemptedBy)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND attempted_at >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND attempted_at <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY attempted_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []infrastructure.ImmutabilityViolationLog
	for rows.Next() {
		var log infrastructure.ImmutabilityViolationLog
		err := rows.Scan(
			&log.ID, &log.OrganizationID, &log.TableName, &log.RecordID, &log.Operation,
			&log.AttemptedBy, &log.AttemptedAt, &log.ErrorMessage, &log.BlockedData, &log.Metadata,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

func (r *ImmutabilityViolationLogRepository) Count(ctx context.Context, orgID *uuid.UUID, filters infrastructure.ImmutabilityViolationLogFilters) (int64, error) {
	if orgID != nil {
		if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
			return 0, err
		}
	}

	query := "SELECT COUNT(*) FROM immutability_violations_log WHERE 1=1"
	var args []interface{}
	argCount := 0

	if orgID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *orgID)
	}

	if filters.TableName != nil {
		argCount++
		query += fmt.Sprintf(" AND table_name = $%d", argCount)
		args = append(args, *filters.TableName)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *ImmutabilityViolationLogRepository) Create(ctx context.Context, log *infrastructure.ImmutabilityViolationLog) error {
	if log.OrganizationID != nil {
		if err := r.db.SetOrganizationContext(ctx, log.OrganizationID.String()); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO immutability_violations_log (
			id, organization_id, table_name, record_id, operation,
			attempted_by, attempted_at, error_message, blocked_data, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		log.ID, log.OrganizationID, log.TableName, log.RecordID, log.Operation,
		log.AttemptedBy, log.AttemptedAt, log.ErrorMessage, log.BlockedData, log.Metadata,
	)
	return err
}

func (r *ImmutabilityViolationLogRepository) Get(ctx context.Context, id uuid.UUID) (*infrastructure.ImmutabilityViolationLog, error) {
	query := `
		SELECT id, organization_id, table_name, record_id, operation,
		       attempted_by, attempted_at, error_message, blocked_data, metadata
		FROM immutability_violations_log
		WHERE id = $1
	`

	var log infrastructure.ImmutabilityViolationLog
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&log.ID, &log.OrganizationID, &log.TableName, &log.RecordID, &log.Operation,
		&log.AttemptedBy, &log.AttemptedAt, &log.ErrorMessage, &log.BlockedData, &log.Metadata,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *ImmutabilityViolationLogRepository) GetByTableAndRecord(ctx context.Context, tableName string, recordID uuid.UUID) ([]infrastructure.ImmutabilityViolationLog, error) {
	query := `
		SELECT id, organization_id, table_name, record_id, operation,
		       attempted_by, attempted_at, error_message, blocked_data, metadata
		FROM immutability_violations_log
		WHERE table_name = $1 AND record_id = $2
		ORDER BY attempted_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, tableName, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []infrastructure.ImmutabilityViolationLog
	for rows.Next() {
		var log infrastructure.ImmutabilityViolationLog
		err := rows.Scan(
			&log.ID, &log.OrganizationID, &log.TableName, &log.RecordID, &log.Operation,
			&log.AttemptedBy, &log.AttemptedAt, &log.ErrorMessage, &log.BlockedData, &log.Metadata,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

func (r *ImmutabilityViolationLogRepository) Cleanup(ctx context.Context, before time.Time) error {
	query := "DELETE FROM immutability_violations_log WHERE attempted_at < $1"
	_, err := r.db.Pool.Exec(ctx, query, before)
	return err
}

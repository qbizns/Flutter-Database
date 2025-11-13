package email_queue

import (
	"encoding/json"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for EmailQueue
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new EmailQueue repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// EmailQueue represents a email_queue entity
type EmailQueue struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	ToAddresses string `json:"to_addresses" db:"to_addresses"`
	CcAddresses *string `json:"cc_addresses" db:"cc_addresses"`
	BccAddresses *string `json:"bcc_addresses" db:"bcc_addresses"`
	FromAddress *string `json:"from_address" db:"from_address"`
	ReplyTo *string `json:"reply_to" db:"reply_to"`
	Subject string `json:"subject" db:"subject"`
	BodyHtml *string `json:"body_html" db:"body_html"`
	BodyText *string `json:"body_text" db:"body_text"`
	AttachmentIds *uuid.UUID `json:"attachment_ids" db:"attachment_ids"`
	TemplateName *string `json:"template_name" db:"template_name"`
	TemplateData json.RawMessage `json:"template_data" db:"template_data"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	Provider *string `json:"provider" db:"provider"`
	ProviderMessageId *string `json:"provider_message_id" db:"provider_message_id"`
	Attempts *int64 `json:"attempts" db:"attempts"`
	MaxAttempts *int64 `json:"max_attempts" db:"max_attempts"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	Priority *int64 `json:"priority" db:"priority"`
	ScheduledAt *time.Time `json:"scheduled_at" db:"scheduled_at"`
	SentAt *time.Time `json:"sent_at" db:"sent_at"`
	FailedAt *time.Time `json:"failed_at" db:"failed_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new email_queue record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *EmailQueue) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "email_queue", duration, nil)
	}()

	query := `
		INSERT INTO email_queue (
			, organization_id
			, to_addresses
			, cc_addresses
			, bcc_addresses
			, from_address
			, reply_to
			, subject
			, body_html
			, body_text
			, attachment_ids
			, template_name
			, template_data
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, priority
			, scheduled_at
			, sent_at
			, failed_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $10
			, $11
			, $12
			, $13
			, $14
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ToAddresses,
		entity.CcAddresses,
		entity.BccAddresses,
		entity.FromAddress,
		entity.ReplyTo,
		entity.Subject,
		entity.BodyHtml,
		entity.BodyText,
		entity.AttachmentIds,
		entity.TemplateName,
		entity.TemplateData,
		entity.Status,
		entity.Status,
		entity.Provider,
		entity.ProviderMessageId,
		entity.Attempts,
		entity.MaxAttempts,
		entity.ErrorMessage,
		entity.Priority,
		entity.ScheduledAt,
		entity.SentAt,
		entity.FailedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create email_queue", zap.Error(err))
		return fmt.Errorf("failed to create email_queue: %w", err)
	}

	r.logger.Info("created email_queue",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a email_queue by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*EmailQueue, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "email_queue", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, to_addresses
			, cc_addresses
			, bcc_addresses
			, from_address
			, reply_to
			, subject
			, body_html
			, body_text
			, attachment_ids
			, template_name
			, template_data
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, priority
			, scheduled_at
			, sent_at
			, failed_at
			, created_at
		FROM email_queue
		WHERE id = $1
		
	`

	var entity EmailQueue
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ToAddresses,
		&entity.CcAddresses,
		&entity.BccAddresses,
		&entity.FromAddress,
		&entity.ReplyTo,
		&entity.Subject,
		&entity.BodyHtml,
		&entity.BodyText,
		&entity.AttachmentIds,
		&entity.TemplateName,
		&entity.TemplateData,
		&entity.Status,
		&entity.Status,
		&entity.Provider,
		&entity.ProviderMessageId,
		&entity.Attempts,
		&entity.MaxAttempts,
		&entity.ErrorMessage,
		&entity.Priority,
		&entity.ScheduledAt,
		&entity.SentAt,
		&entity.FailedAt,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("email_queue not found")
	}

	if err != nil {
		r.logger.Error("failed to get email_queue", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get email_queue: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of email_queue records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*EmailQueue, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "email_queue", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM email_queue
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count email_queue records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, to_addresses
			, cc_addresses
			, bcc_addresses
			, from_address
			, reply_to
			, subject
			, body_html
			, body_text
			, attachment_ids
			, template_name
			, template_data
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, priority
			, scheduled_at
			, sent_at
			, failed_at
			, created_at
		FROM email_queue
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list email_queue", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list email_queue: %w", err)
	}
	defer rows.Close()

	var entities []*EmailQueue
	for rows.Next() {
		var entity EmailQueue
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ToAddresses,
			&entity.CcAddresses,
			&entity.BccAddresses,
			&entity.FromAddress,
			&entity.ReplyTo,
			&entity.Subject,
			&entity.BodyHtml,
			&entity.BodyText,
			&entity.AttachmentIds,
			&entity.TemplateName,
			&entity.TemplateData,
			&entity.Status,
			&entity.Status,
			&entity.Provider,
			&entity.ProviderMessageId,
			&entity.Attempts,
			&entity.MaxAttempts,
			&entity.ErrorMessage,
			&entity.Priority,
			&entity.ScheduledAt,
			&entity.SentAt,
			&entity.FailedAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan email_queue: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating email_queue rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing email_queue record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *EmailQueue) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "email_queue", duration, nil)
	}()

	query := `
		UPDATE email_queue
		SET
			, organization_id = $2
			, to_addresses = $3
			, cc_addresses = $4
			, bcc_addresses = $5
			, from_address = $6
			, reply_to = $7
			, subject = $8
			, body_html = $9
			, body_text = $10
			, attachment_ids = $11
			, template_name = $12
			, template_data = $13
			, status = $14
			, status = $15
			, provider = $16
			, provider_message_id = $17
			, attempts = $18
			, max_attempts = $19
			, error_message = $20
			, priority = $21
			, scheduled_at = $22
			, sent_at = $23
			, failed_at = $24
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $26
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ToAddresses,
		entity.CcAddresses,
		entity.BccAddresses,
		entity.FromAddress,
		entity.ReplyTo,
		entity.Subject,
		entity.BodyHtml,
		entity.BodyText,
		entity.AttachmentIds,
		entity.TemplateName,
		entity.TemplateData,
		entity.Status,
		entity.Status,
		entity.Provider,
		entity.ProviderMessageId,
		entity.Attempts,
		entity.MaxAttempts,
		entity.ErrorMessage,
		entity.Priority,
		entity.ScheduledAt,
		entity.SentAt,
		entity.FailedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update email_queue", zap.Error(err))
		return fmt.Errorf("failed to update email_queue: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("email_queue not found or already deleted")
	}

	r.logger.Info("updated email_queue",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a email_queue record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "email_queue", duration, nil)
	}()

	query := `DELETE FROM email_queue WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete email_queue", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete email_queue: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("email_queue not found")
	}

	r.logger.Info("deleted email_queue", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves email_queue records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*EmailQueue, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "email_queue", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM email_queue
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count email_queue records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, to_addresses
			, cc_addresses
			, bcc_addresses
			, from_address
			, reply_to
			, subject
			, body_html
			, body_text
			, attachment_ids
			, template_name
			, template_data
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, priority
			, scheduled_at
			, sent_at
			, failed_at
			, created_at
		FROM email_queue
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list email_queue by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list email_queue: %w", err)
	}
	defer rows.Close()

	var entities []*EmailQueue
	for rows.Next() {
		var entity EmailQueue
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ToAddresses,
			&entity.CcAddresses,
			&entity.BccAddresses,
			&entity.FromAddress,
			&entity.ReplyTo,
			&entity.Subject,
			&entity.BodyHtml,
			&entity.BodyText,
			&entity.AttachmentIds,
			&entity.TemplateName,
			&entity.TemplateData,
			&entity.Status,
			&entity.Status,
			&entity.Provider,
			&entity.ProviderMessageId,
			&entity.Attempts,
			&entity.MaxAttempts,
			&entity.ErrorMessage,
			&entity.Priority,
			&entity.ScheduledAt,
			&entity.SentAt,
			&entity.FailedAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan email_queue: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


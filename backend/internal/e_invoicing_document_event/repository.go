package e_invoicing_document_event

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

// Repository handles database operations for EInvoicingDocumentEvents
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new EInvoicingDocumentEvents repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// EInvoicingDocumentEvents represents a e_invoicing_document_events entity
type EInvoicingDocumentEvents struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	EInvoicingDocumentId uuid.UUID `json:"e_invoicing_document_id" db:"e_invoicing_document_id"`
	EventType string `json:"event_type" db:"event_type"`
	'created', *string `json:"'created'," db:"'created',"`
	'validated', *string `json:"'validated'," db:"'validated',"`
	'submitted', *string `json:"'submitted'," db:"'submitted',"`
	'accepted', *string `json:"'accepted'," db:"'accepted',"`
	'rejected', *string `json:"'rejected'," db:"'rejected',"`
	'cancelled', *string `json:"'cancelled'," db:"'cancelled',"`
	'error', *string `json:"'error'," db:"'error',"`
	'retry', *string `json:"'retry'," db:"'retry',"`
	'statusCheck' *string `json:"'status_check'" db:"'status_check'"`
	EventTimestamp time.Time `json:"event_timestamp" db:"event_timestamp"`
	PreviousStatus *string `json:"previous_status" db:"previous_status"`
	// 	NewStatus *string `json:"new_status" db:"new_status"`
	EventDescription *string `json:"event_description" db:"event_description"`
	EventData json.RawMessage `json:"event_data" db:"event_data"`
	HttpStatusCode *int64 `json:"http_status_code" db:"http_status_code"`
	HttpMethod *string `json:"http_method" db:"http_method"`
	ApiEndpoint *string `json:"api_endpoint" db:"api_endpoint"`
	RequestHeaders json.RawMessage `json:"request_headers" db:"request_headers"`
	ResponseHeaders json.RawMessage `json:"response_headers" db:"response_headers"`
	ErrorCode *string `json:"error_code" db:"error_code"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	ErrorDetails json.RawMessage `json:"error_details" db:"error_details"`
	TriggeredBy string `json:"triggered_by" db:"triggered_by"`
	UserId *uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new e_invoicing_document_events record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocumentEvents) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "e_invoicing_document_events", duration, nil)
	}()

	query := `
		INSERT INTO e_invoicing_document_events (
			, organization_id
			, e_invoicing_document_id
			, event_type
			, 'created',
			, 'validated',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error',
			, 'retry',
			, 'status_check'
			, event_timestamp
			, previous_status
			, new_status
			, event_description
			, event_data
			, http_status_code
			, http_method
			, api_endpoint
			, request_headers
			, response_headers
			, error_code
			, error_message
			, error_details
			, triggered_by
			, user_id
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
			, $25
			, $26
			, $27
			, $28
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.EInvoicingDocumentId,
		entity.EventType,
		entity.'created',,
		entity.'validated',,
		entity.'submitted',,
		entity.'accepted',,
		entity.'rejected',,
		entity.'cancelled',,
		entity.'error',,
		entity.'retry',,
		entity.'statusCheck',
		entity.EventTimestamp,
		entity.PreviousStatus,
		entity.NewStatus,
		entity.EventDescription,
		entity.EventData,
		entity.HttpStatusCode,
		entity.HttpMethod,
		entity.ApiEndpoint,
		entity.RequestHeaders,
		entity.ResponseHeaders,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.ErrorDetails,
		entity.TriggeredBy,
		entity.UserId,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create e_invoicing_document_events", zap.Error(err))
		return fmt.Errorf("failed to create e_invoicing_document_events: %w", err)
	}

	r.logger.Info("created e_invoicing_document_events",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a e_invoicing_document_events by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*EInvoicingDocumentEvents, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "e_invoicing_document_events", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, e_invoicing_document_id
			, event_type
			, 'created',
			, 'validated',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error',
			, 'retry',
			, 'status_check'
			, event_timestamp
			, previous_status
			, new_status
			, event_description
			, event_data
			, http_status_code
			, http_method
			, api_endpoint
			, request_headers
			, response_headers
			, error_code
			, error_message
			, error_details
			, triggered_by
			, user_id
			, created_at
		FROM e_invoicing_document_events
		WHERE id = $1
		
	`

	var entity EInvoicingDocumentEvents
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.EInvoicingDocumentId,
		&entity.EventType,
		&entity.'created',,
		&entity.'validated',,
		&entity.'submitted',,
		&entity.'accepted',,
		&entity.'rejected',,
		&entity.'cancelled',,
		&entity.'error',,
		&entity.'retry',,
		&entity.'statusCheck',
		&entity.EventTimestamp,
		&entity.PreviousStatus,
		&entity.NewStatus,
		&entity.EventDescription,
		&entity.EventData,
		&entity.HttpStatusCode,
		&entity.HttpMethod,
		&entity.ApiEndpoint,
		&entity.RequestHeaders,
		&entity.ResponseHeaders,
		&entity.ErrorCode,
		&entity.ErrorMessage,
		&entity.ErrorDetails,
		&entity.TriggeredBy,
		&entity.UserId,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("e_invoicing_document_events not found")
	}

	if err != nil {
		r.logger.Error("failed to get e_invoicing_document_events", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get e_invoicing_document_events: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of e_invoicing_document_events records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*EInvoicingDocumentEvents, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "e_invoicing_document_events", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM e_invoicing_document_events
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count e_invoicing_document_events records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, e_invoicing_document_id
			, event_type
			, 'created',
			, 'validated',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error',
			, 'retry',
			, 'status_check'
			, event_timestamp
			, previous_status
			, new_status
			, event_description
			, event_data
			, http_status_code
			, http_method
			, api_endpoint
			, request_headers
			, response_headers
			, error_code
			, error_message
			, error_details
			, triggered_by
			, user_id
			, created_at
		FROM e_invoicing_document_events
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list e_invoicing_document_events", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list e_invoicing_document_events: %w", err)
	}
	defer rows.Close()

	var entities []*EInvoicingDocumentEvents
	for rows.Next() {
		var entity EInvoicingDocumentEvents
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.EInvoicingDocumentId,
			&entity.EventType,
			&entity.'created',,
			&entity.'validated',,
			&entity.'submitted',,
			&entity.'accepted',,
			&entity.'rejected',,
			&entity.'cancelled',,
			&entity.'error',,
			&entity.'retry',,
			&entity.'statusCheck',
			&entity.EventTimestamp,
			&entity.PreviousStatus,
			&entity.NewStatus,
			&entity.EventDescription,
			&entity.EventData,
			&entity.HttpStatusCode,
			&entity.HttpMethod,
			&entity.ApiEndpoint,
			&entity.RequestHeaders,
			&entity.ResponseHeaders,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.ErrorDetails,
			&entity.TriggeredBy,
			&entity.UserId,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan e_invoicing_document_events: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating e_invoicing_document_events rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing e_invoicing_document_events record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocumentEvents) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "e_invoicing_document_events", duration, nil)
	}()

	query := `
		UPDATE e_invoicing_document_events
		SET
			, organization_id = $2
			, e_invoicing_document_id = $3
			, event_type = $4
			, 'created', = $5
			, 'validated', = $6
			, 'submitted', = $7
			, 'accepted', = $8
			, 'rejected', = $9
			, 'cancelled', = $10
			, 'error', = $11
			, 'retry', = $12
			, 'status_check' = $13
			, event_timestamp = $14
			, previous_status = $15
			, new_status = $16
			, event_description = $17
			, event_data = $18
			, http_status_code = $19
			, http_method = $20
			, api_endpoint = $21
			, request_headers = $22
			, response_headers = $23
			, error_code = $24
			, error_message = $25
			, error_details = $26
			, triggered_by = $27
			, user_id = $28
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $30
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.EInvoicingDocumentId,
		entity.EventType,
		entity.'created',,
		entity.'validated',,
		entity.'submitted',,
		entity.'accepted',,
		entity.'rejected',,
		entity.'cancelled',,
		entity.'error',,
		entity.'retry',,
		entity.'statusCheck',
		entity.EventTimestamp,
		entity.PreviousStatus,
		entity.NewStatus,
		entity.EventDescription,
		entity.EventData,
		entity.HttpStatusCode,
		entity.HttpMethod,
		entity.ApiEndpoint,
		entity.RequestHeaders,
		entity.ResponseHeaders,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.ErrorDetails,
		entity.TriggeredBy,
		entity.UserId,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update e_invoicing_document_events", zap.Error(err))
		return fmt.Errorf("failed to update e_invoicing_document_events: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("e_invoicing_document_events not found or already deleted")
	}

	r.logger.Info("updated e_invoicing_document_events",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a e_invoicing_document_events record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "e_invoicing_document_events", duration, nil)
	}()

	query := `DELETE FROM e_invoicing_document_events WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete e_invoicing_document_events", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete e_invoicing_document_events: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("e_invoicing_document_events not found")
	}

	r.logger.Info("deleted e_invoicing_document_events", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves e_invoicing_document_events records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*EInvoicingDocumentEvents, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "e_invoicing_document_events", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM e_invoicing_document_events
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count e_invoicing_document_events records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, e_invoicing_document_id
			, event_type
			, 'created',
			, 'validated',
			, 'submitted',
			, 'accepted',
			, 'rejected',
			, 'cancelled',
			, 'error',
			, 'retry',
			, 'status_check'
			, event_timestamp
			, previous_status
			, new_status
			, event_description
			, event_data
			, http_status_code
			, http_method
			, api_endpoint
			, request_headers
			, response_headers
			, error_code
			, error_message
			, error_details
			, triggered_by
			, user_id
			, created_at
		FROM e_invoicing_document_events
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list e_invoicing_document_events by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list e_invoicing_document_events: %w", err)
	}
	defer rows.Close()

	var entities []*EInvoicingDocumentEvents
	for rows.Next() {
		var entity EInvoicingDocumentEvents
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.EInvoicingDocumentId,
			&entity.EventType,
			&entity.'created',,
			&entity.'validated',,
			&entity.'submitted',,
			&entity.'accepted',,
			&entity.'rejected',,
			&entity.'cancelled',,
			&entity.'error',,
			&entity.'retry',,
			&entity.'statusCheck',
			&entity.EventTimestamp,
			&entity.PreviousStatus,
			&entity.NewStatus,
			&entity.EventDescription,
			&entity.EventData,
			&entity.HttpStatusCode,
			&entity.HttpMethod,
			&entity.ApiEndpoint,
			&entity.RequestHeaders,
			&entity.ResponseHeaders,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.ErrorDetails,
			&entity.TriggeredBy,
			&entity.UserId,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan e_invoicing_document_events: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


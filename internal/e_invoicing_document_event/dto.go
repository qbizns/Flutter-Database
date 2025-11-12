package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EInvoicingDocumentEventsResponse represents a e_invoicing_document_events response
type EInvoicingDocumentEventsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	EInvoicingDocumentId uuid.UUID `json:"e_invoicing_document_id"`
	
	EventType string `json:"event_type"`
	
	'created', *string `json:"'created',"`
	
	'validated', *string `json:"'validated',"`
	
	'submitted', *string `json:"'submitted',"`
	
	'accepted', *string `json:"'accepted',"`
	
	'rejected', *string `json:"'rejected',"`
	
	'cancelled', *string `json:"'cancelled',"`
	
	'error', *string `json:"'error',"`
	
	'retry', *string `json:"'retry',"`
	
	'statusCheck' *string `json:"'status_check'"`
	
	EventTimestamp time.Time `json:"event_timestamp"`
	
	PreviousStatus *string `json:"previous_status"`
	
	NewStatus *string `json:"new_status"`
	
	EventDescription *string `json:"event_description"`
	
	EventData json.RawMessage `json:"event_data"`
	
	HttpStatusCode *int64 `json:"http_status_code"`
	
	HttpMethod *string `json:"http_method"`
	
	ApiEndpoint *string `json:"api_endpoint"`
	
	RequestHeaders json.RawMessage `json:"request_headers"`
	
	ResponseHeaders json.RawMessage `json:"response_headers"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorDetails json.RawMessage `json:"error_details"`
	
	TriggeredBy string `json:"triggered_by"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	CreatedAt time.Time `json:"created_at"`
	
}

// CreateEInvoicingDocumentEventsRequest represents a request to create a e_invoicing_document_events
type CreateEInvoicingDocumentEventsRequest struct {
	
	EInvoicingDocumentId uuid.UUID `json:"e_invoicing_document_id" validate:"required"`
	
	EventType string `json:"event_type" validate:"required"`
	
	'created', *string `json:"'created',"`
	
	'validated', *string `json:"'validated',"`
	
	'submitted', *string `json:"'submitted',"`
	
	'accepted', *string `json:"'accepted',"`
	
	'rejected', *string `json:"'rejected',"`
	
	'cancelled', *string `json:"'cancelled',"`
	
	'error', *string `json:"'error',"`
	
	'retry', *string `json:"'retry',"`
	
	'statusCheck' *string `json:"'status_check'"`
	
	EventTimestamp time.Time `json:"event_timestamp" validate:"required"`
	
	PreviousStatus *string `json:"previous_status"`
	
	NewStatus *string `json:"new_status"`
	
	EventDescription *string `json:"event_description"`
	
	EventData json.RawMessage `json:"event_data"`
	
	HttpStatusCode *int64 `json:"http_status_code"`
	
	HttpMethod *string `json:"http_method"`
	
	ApiEndpoint *string `json:"api_endpoint"`
	
	RequestHeaders json.RawMessage `json:"request_headers"`
	
	ResponseHeaders json.RawMessage `json:"response_headers"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorDetails json.RawMessage `json:"error_details"`
	
	TriggeredBy string `json:"triggered_by" validate:"required"`
	
	UserId *uuid.UUID `json:"user_id"`
	
}

// Validate validates the create request
func (r *CreateEInvoicingDocumentEventsRequest) Validate() error {
	
	if r.EInvoicingDocumentId == uuid.Nil {
		return fmt.Errorf("e_invoicing_document_id is required")
	}
	
	if r.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	
	if r.EventTimestamp == nil {
		return fmt.Errorf("event_timestamp is required")
	}
	
	if r.TriggeredBy == "" {
		return fmt.Errorf("triggered_by is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateEInvoicingDocumentEventsRequest represents a request to update a e_invoicing_document_events
type UpdateEInvoicingDocumentEventsRequest struct {
	
	EInvoicingDocumentId *uuid.UUID `json:"e_invoicing_document_id,omitempty" validate:"omitempty,required"`
	
	EventType *string `json:"event_type,omitempty" validate:"omitempty,required"`
	
	'created', *string `json:"'created',,omitempty"`
	
	'validated', *string `json:"'validated',,omitempty"`
	
	'submitted', *string `json:"'submitted',,omitempty"`
	
	'accepted', *string `json:"'accepted',,omitempty"`
	
	'rejected', *string `json:"'rejected',,omitempty"`
	
	'cancelled', *string `json:"'cancelled',,omitempty"`
	
	'error', *string `json:"'error',,omitempty"`
	
	'retry', *string `json:"'retry',,omitempty"`
	
	'statusCheck' *string `json:"'status_check',omitempty"`
	
	EventTimestamp *time.Time `json:"event_timestamp,omitempty" validate:"omitempty,required"`
	
	PreviousStatus *string `json:"previous_status,omitempty"`
	
	NewStatus *string `json:"new_status,omitempty"`
	
	EventDescription *string `json:"event_description,omitempty"`
	
	EventData *json.RawMessage `json:"event_data,omitempty"`
	
	HttpStatusCode *int64 `json:"http_status_code,omitempty"`
	
	HttpMethod *string `json:"http_method,omitempty"`
	
	ApiEndpoint *string `json:"api_endpoint,omitempty"`
	
	RequestHeaders *json.RawMessage `json:"request_headers,omitempty"`
	
	ResponseHeaders *json.RawMessage `json:"response_headers,omitempty"`
	
	ErrorCode *string `json:"error_code,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	ErrorDetails *json.RawMessage `json:"error_details,omitempty"`
	
	TriggeredBy *string `json:"triggered_by,omitempty" validate:"omitempty,required"`
	
	UserId *uuid.UUID `json:"user_id,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateEInvoicingDocumentEventsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.EInvoicingDocumentId != nil {
		hasUpdate = true
	}
	
	if r.EventType != nil {
		hasUpdate = true
	}
	
	if r.'created', != nil {
		hasUpdate = true
	}
	
	if r.'validated', != nil {
		hasUpdate = true
	}
	
	if r.'submitted', != nil {
		hasUpdate = true
	}
	
	if r.'accepted', != nil {
		hasUpdate = true
	}
	
	if r.'rejected', != nil {
		hasUpdate = true
	}
	
	if r.'cancelled', != nil {
		hasUpdate = true
	}
	
	if r.'error', != nil {
		hasUpdate = true
	}
	
	if r.'retry', != nil {
		hasUpdate = true
	}
	
	if r.'statusCheck' != nil {
		hasUpdate = true
	}
	
	if r.EventTimestamp != nil {
		hasUpdate = true
	}
	
	if r.PreviousStatus != nil {
		hasUpdate = true
	}
	
	if r.NewStatus != nil {
		hasUpdate = true
	}
	
	if r.EventDescription != nil {
		hasUpdate = true
	}
	
	if r.EventData != nil {
		hasUpdate = true
	}
	
	if r.HttpStatusCode != nil {
		hasUpdate = true
	}
	
	if r.HttpMethod != nil {
		hasUpdate = true
	}
	
	if r.ApiEndpoint != nil {
		hasUpdate = true
	}
	
	if r.RequestHeaders != nil {
		hasUpdate = true
	}
	
	if r.ResponseHeaders != nil {
		hasUpdate = true
	}
	
	if r.ErrorCode != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.ErrorDetails != nil {
		hasUpdate = true
	}
	
	if r.TriggeredBy != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// EInvoicingDocumentEventsListResponse represents a paginated list of e_invoicing_document_events records
type EInvoicingDocumentEventsListResponse struct {
	Items      []*EInvoicingDocumentEventsResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

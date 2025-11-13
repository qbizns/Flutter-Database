package email_queue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EmailQueueResponse represents a email_queue response
type EmailQueueResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	ToAddresses string `json:"to_addresses"`
	
	CcAddresses *string `json:"cc_addresses"`
	
	BccAddresses *string `json:"bcc_addresses"`
	
	FromAddress *string `json:"from_address"`
	
	ReplyTo *string `json:"reply_to"`
	
	Subject string `json:"subject"`
	
	BodyHtml *string `json:"body_html"`
	
	BodyText *string `json:"body_text"`
	
	AttachmentIds *uuid.UUID `json:"attachment_ids"`
	
	TemplateName *string `json:"template_name"`
	
	TemplateData json.RawMessage `json:"template_data"`
	
	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	Provider *string `json:"provider"`
	
	ProviderMessageId *string `json:"provider_message_id"`
	
	Attempts *int64 `json:"attempts"`
	
	MaxAttempts *int64 `json:"max_attempts"`
	
	ErrorMessage *string `json:"error_message"`
	
	Priority *int64 `json:"priority"`
	
	ScheduledAt *time.Time `json:"scheduled_at"`
	
	SentAt *time.Time `json:"sent_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
	CreatedAt *time.Time `json:"created_at"`
	
}

// CreateEmailQueueRequest represents a request to create a email_queue
type CreateEmailQueueRequest struct {
	
	ToAddresses string `json:"to_addresses" validate:"required"`
	
	CcAddresses *string `json:"cc_addresses"`
	
	BccAddresses *string `json:"bcc_addresses"`
	
	FromAddress *string `json:"from_address"`
	
	ReplyTo *string `json:"reply_to"`
	
	Subject string `json:"subject" validate:"required"`
	
	BodyHtml *string `json:"body_html"`
	
	BodyText *string `json:"body_text"`
	
	AttachmentIds *uuid.UUID `json:"attachment_ids"`
	
	TemplateName *string `json:"template_name"`
	
	// Duplicate removed: TemplateData json.RawMessage `json:"template_data"`
	
	// 	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	Provider *string `json:"provider"`
	
	ProviderMessageId *string `json:"provider_message_id"`
	
	Attempts *int64 `json:"attempts"`
	
	MaxAttempts *int64 `json:"max_attempts"`
	
	ErrorMessage *string `json:"error_message"`
	
	Priority *int64 `json:"priority"`
	
	ScheduledAt *time.Time `json:"scheduled_at"`
	
	SentAt *time.Time `json:"sent_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
}

// Validate validates the create request
func (r *CreateEmailQueueRequest) Validate() error {
	
	if r.ToAddresses == "" {
		return fmt.Errorf("to_addresses is required")
	}
	
	if r.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateEmailQueueRequest represents a request to update a email_queue
type UpdateEmailQueueRequest struct {
	
	ToAddresses *string `json:"to_addresses,omitempty" validate:"omitempty,required"`
	
	CcAddresses *string `json:"cc_addresses,omitempty"`
	
	BccAddresses *string `json:"bcc_addresses,omitempty"`
	
	FromAddress *string `json:"from_address,omitempty"`
	
	ReplyTo *string `json:"reply_to,omitempty"`
	
	Subject *string `json:"subject,omitempty" validate:"omitempty,required"`
	
	BodyHtml *string `json:"body_html,omitempty"`
	
	BodyText *string `json:"body_text,omitempty"`
	
	AttachmentIds *uuid.UUID `json:"attachment_ids,omitempty"`
	
	TemplateName *string `json:"template_name,omitempty"`
	
	TemplateData *json.RawMessage `json:"template_data,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	Provider *string `json:"provider,omitempty"`
	
	ProviderMessageId *string `json:"provider_message_id,omitempty"`
	
	Attempts *int64 `json:"attempts,omitempty"`
	
	MaxAttempts *int64 `json:"max_attempts,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	
	SentAt *time.Time `json:"sent_at,omitempty"`
	
	FailedAt *time.Time `json:"failed_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateEmailQueueRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ToAddresses != nil {
		hasUpdate = true
	}
	
	if r.CcAddresses != nil {
		hasUpdate = true
	}
	
	if r.BccAddresses != nil {
		hasUpdate = true
	}
	
	if r.FromAddress != nil {
		hasUpdate = true
	}
	
	if r.ReplyTo != nil {
		hasUpdate = true
	}
	
	if r.Subject != nil {
		hasUpdate = true
	}
	
	if r.BodyHtml != nil {
		hasUpdate = true
	}
	
	if r.BodyText != nil {
		hasUpdate = true
	}
	
	if r.AttachmentIds != nil {
		hasUpdate = true
	}
	
	if r.TemplateName != nil {
		hasUpdate = true
	}
	
	if r.TemplateData != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Provider != nil {
		hasUpdate = true
	}
	
	if r.ProviderMessageId != nil {
		hasUpdate = true
	}
	
	if r.Attempts != nil {
		hasUpdate = true
	}
	
	if r.MaxAttempts != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.ScheduledAt != nil {
		hasUpdate = true
	}
	
	if r.SentAt != nil {
		hasUpdate = true
	}
	
	if r.FailedAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// EmailQueueListResponse represents a paginated list of email_queue records
type EmailQueueListResponse struct {
	Items      []*EmailQueueResponse `json:"items"`
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

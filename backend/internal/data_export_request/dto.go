package data_export_request

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DataExportRequestsResponse represents a data_export_requests response
type DataExportRequestsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ExportType string `json:"export_type"`
	
	ExportFormat string `json:"export_format"`
	
	DateFrom *time.Time `json:"date_from"`
	
	DateTo *time.Time `json:"date_to"`
	
	Filters json.RawMessage `json:"filters"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	FileName *string `json:"file_name"`
	
	FileSize *int64 `json:"file_size"`
	
	FilePath *string `json:"file_path"`
	
	DownloadUrl *string `json:"download_url"`
	
	DownloadExpiresAt *time.Time `json:"download_expires_at"`
	
	TotalRecords *int64 `json:"total_records"`
	
	ProcessedRecords *int64 `json:"processed_records"`
	
	ErrorMessage *string `json:"error_message"`
	
	RequestedBy *uuid.UUID `json:"requested_by"`
	
	RequestedAt *time.Time `json:"requested_at"`
	
	StartedAt *time.Time `json:"started_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
}

// CreateDataExportRequestsRequest represents a request to create a data_export_requests
type CreateDataExportRequestsRequest struct {
	
	ExportType string `json:"export_type" validate:"required"`
	
	ExportFormat string `json:"export_format" validate:"required"`
	
	DateFrom *time.Time `json:"date_from"`
	
	DateTo *time.Time `json:"date_to"`
	
	Filters json.RawMessage `json:"filters"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	FileName *string `json:"file_name"`
	
	FileSize *int64 `json:"file_size"`
	
	FilePath *string `json:"file_path"`
	
	DownloadUrl *string `json:"download_url" validate:"url"`
	
	DownloadExpiresAt *time.Time `json:"download_expires_at"`
	
	TotalRecords *int64 `json:"total_records"`
	
	ProcessedRecords *int64 `json:"processed_records"`
	
	ErrorMessage *string `json:"error_message"`
	
	RequestedBy *uuid.UUID `json:"requested_by"`
	
	RequestedAt *time.Time `json:"requested_at"`
	
	StartedAt *time.Time `json:"started_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
}

// Validate validates the create request
func (r *CreateDataExportRequestsRequest) Validate() error {
	
	if r.ExportType == "" {
		return fmt.Errorf("export_type is required")
	}
	
	if r.ExportFormat == "" {
		return fmt.Errorf("export_format is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDataExportRequestsRequest represents a request to update a data_export_requests
type UpdateDataExportRequestsRequest struct {
	
	ExportType *string `json:"export_type,omitempty" validate:"omitempty,required"`
	
	ExportFormat *string `json:"export_format,omitempty" validate:"omitempty,required"`
	
	DateFrom *time.Time `json:"date_from,omitempty"`
	
	DateTo *time.Time `json:"date_to,omitempty"`
	
	Filters *json.RawMessage `json:"filters,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	FileName *string `json:"file_name,omitempty"`
	
	FileSize *int64 `json:"file_size,omitempty"`
	
	FilePath *string `json:"file_path,omitempty"`
	
	DownloadUrl *string `json:"download_url,omitempty" validate:"omitempty,url"`
	
	DownloadExpiresAt *time.Time `json:"download_expires_at,omitempty"`
	
	TotalRecords *int64 `json:"total_records,omitempty"`
	
	ProcessedRecords *int64 `json:"processed_records,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	RequestedBy *uuid.UUID `json:"requested_by,omitempty"`
	
	RequestedAt *time.Time `json:"requested_at,omitempty"`
	
	StartedAt *time.Time `json:"started_at,omitempty"`
	
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDataExportRequestsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ExportType != nil {
		hasUpdate = true
	}
	
	if r.ExportFormat != nil {
		hasUpdate = true
	}
	
	if r.DateFrom != nil {
		hasUpdate = true
	}
	
	if r.DateTo != nil {
		hasUpdate = true
	}
	
	if r.Filters != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.FileName != nil {
		hasUpdate = true
	}
	
	if r.FileSize != nil {
		hasUpdate = true
	}
	
	if r.FilePath != nil {
		hasUpdate = true
	}
	
	if r.DownloadUrl != nil {
		hasUpdate = true
	}
	
	if r.DownloadExpiresAt != nil {
		hasUpdate = true
	}
	
	if r.TotalRecords != nil {
		hasUpdate = true
	}
	
	if r.ProcessedRecords != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.RequestedBy != nil {
		hasUpdate = true
	}
	
	if r.RequestedAt != nil {
		hasUpdate = true
	}
	
	if r.StartedAt != nil {
		hasUpdate = true
	}
	
	if r.CompletedAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// DataExportRequestsListResponse represents a paginated list of data_export_requests records
type DataExportRequestsListResponse struct {
	Items      []*DataExportRequestsResponse `json:"items"`
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

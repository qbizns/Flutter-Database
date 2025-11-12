package file_attachment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FileAttachmentsResponse represents a file_attachments response
type FileAttachmentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	FileName string `json:"file_name"`
	
	FileSize int64 `json:"file_size"`
	
	MimeType string `json:"mime_type"`
	
	FileExtension *string `json:"file_extension"`
	
	StorageProvider *string `json:"storage_provider"`
	
	StoragePath string `json:"storage_path"`
	
	StorageUrl *string `json:"storage_url"`
	
	FileHash *string `json:"file_hash"`
	
	EntityType string `json:"entity_type"`
	
	EntityId uuid.UUID `json:"entity_id"`
	
	Description *string `json:"description"`
	
	Tags *string `json:"tags"`
	
	IsPublic *bool `json:"is_public"`
	
	ImageWidth *int64 `json:"image_width"`
	
	ImageHeight *int64 `json:"image_height"`
	
	VirusScanStatus *string `json:"virus_scan_status"`
	
	VirusScanAt *time.Time `json:"virus_scan_at"`
	
	UploadedBy *uuid.UUID `json:"uploaded_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateFileAttachmentsRequest represents a request to create a file_attachments
type CreateFileAttachmentsRequest struct {
	
	FileName string `json:"file_name" validate:"required"`
	
	FileSize int64 `json:"file_size" validate:"required"`
	
	MimeType string `json:"mime_type" validate:"required"`
	
	FileExtension *string `json:"file_extension"`
	
	StorageProvider *string `json:"storage_provider"`
	
	StoragePath string `json:"storage_path" validate:"required"`
	
	StorageUrl *string `json:"storage_url" validate:"url"`
	
	FileHash *string `json:"file_hash"`
	
	EntityType string `json:"entity_type" validate:"required"`
	
	EntityId uuid.UUID `json:"entity_id" validate:"required"`
	
	Description *string `json:"description"`
	
	Tags *string `json:"tags"`
	
	IsPublic *bool `json:"is_public"`
	
	ImageWidth *int64 `json:"image_width"`
	
	ImageHeight *int64 `json:"image_height"`
	
	VirusScanStatus *string `json:"virus_scan_status"`
	
	VirusScanAt *time.Time `json:"virus_scan_at"`
	
	UploadedBy *uuid.UUID `json:"uploaded_by"`
	
}

// Validate validates the create request
func (r *CreateFileAttachmentsRequest) Validate() error {
	
	if r.FileName == "" {
		return fmt.Errorf("file_name is required")
	}
	
	if r.FileSize == 0 {
		return fmt.Errorf("file_size is required")
	}
	
	if r.MimeType == "" {
		return fmt.Errorf("mime_type is required")
	}
	
	if r.StoragePath == "" {
		return fmt.Errorf("storage_path is required")
	}
	
	if r.EntityType == "" {
		return fmt.Errorf("entity_type is required")
	}
	
	if r.EntityId == uuid.Nil {
		return fmt.Errorf("entity_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateFileAttachmentsRequest represents a request to update a file_attachments
type UpdateFileAttachmentsRequest struct {
	
	FileName *string `json:"file_name,omitempty" validate:"omitempty,required"`
	
	FileSize *int64 `json:"file_size,omitempty" validate:"omitempty,required"`
	
	MimeType *string `json:"mime_type,omitempty" validate:"omitempty,required"`
	
	FileExtension *string `json:"file_extension,omitempty"`
	
	StorageProvider *string `json:"storage_provider,omitempty"`
	
	StoragePath *string `json:"storage_path,omitempty" validate:"omitempty,required"`
	
	StorageUrl *string `json:"storage_url,omitempty" validate:"omitempty,url"`
	
	FileHash *string `json:"file_hash,omitempty"`
	
	EntityType *string `json:"entity_type,omitempty" validate:"omitempty,required"`
	
	EntityId *uuid.UUID `json:"entity_id,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	Tags *string `json:"tags,omitempty"`
	
	IsPublic *bool `json:"is_public,omitempty"`
	
	ImageWidth *int64 `json:"image_width,omitempty"`
	
	ImageHeight *int64 `json:"image_height,omitempty"`
	
	VirusScanStatus *string `json:"virus_scan_status,omitempty"`
	
	VirusScanAt *time.Time `json:"virus_scan_at,omitempty"`
	
	UploadedBy *uuid.UUID `json:"uploaded_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateFileAttachmentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.FileName != nil {
		hasUpdate = true
	}
	
	if r.FileSize != nil {
		hasUpdate = true
	}
	
	if r.MimeType != nil {
		hasUpdate = true
	}
	
	if r.FileExtension != nil {
		hasUpdate = true
	}
	
	if r.StorageProvider != nil {
		hasUpdate = true
	}
	
	if r.StoragePath != nil {
		hasUpdate = true
	}
	
	if r.StorageUrl != nil {
		hasUpdate = true
	}
	
	if r.FileHash != nil {
		hasUpdate = true
	}
	
	if r.EntityType != nil {
		hasUpdate = true
	}
	
	if r.EntityId != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Tags != nil {
		hasUpdate = true
	}
	
	if r.IsPublic != nil {
		hasUpdate = true
	}
	
	if r.ImageWidth != nil {
		hasUpdate = true
	}
	
	if r.ImageHeight != nil {
		hasUpdate = true
	}
	
	if r.VirusScanStatus != nil {
		hasUpdate = true
	}
	
	if r.VirusScanAt != nil {
		hasUpdate = true
	}
	
	if r.UploadedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// FileAttachmentsListResponse represents a paginated list of file_attachments records
type FileAttachmentsListResponse struct {
	Items      []*FileAttachmentsResponse `json:"items"`
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

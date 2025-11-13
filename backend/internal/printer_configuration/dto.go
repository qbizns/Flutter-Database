package printer_configuration

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PrinterConfigurationsResponse represents a printer_configurations response
type PrinterConfigurationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	PrinterDeviceId uuid.UUID `json:"printer_device_id"`
	
	DocumentType string `json:"document_type"`
	
	FilterOrderType *string `json:"filter_order_type"`
	
	FilterKitchenStationId *uuid.UUID `json:"filter_kitchen_station_id"`
	
	FilterProductCategoryId *uuid.UUID `json:"filter_product_category_id"`
	
	FilterCourseId *uuid.UUID `json:"filter_course_id"`
	
	NumberOfCopies *int64 `json:"number_of_copies"`
	
	AutoPrint *bool `json:"auto_print"`
	
	PrintPriority *int64 `json:"print_priority"`
	
	TemplateConfig json.RawMessage `json:"template_config"`
	
	PaperSize *string `json:"paper_size"`
	
	PrintOrientation *string `json:"print_orientation"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'receipt', *string `json:"'receipt',"`
	
	'report', *string `json:"'report',"`
	
}

// CreatePrinterConfigurationsRequest represents a request to create a printer_configurations
type CreatePrinterConfigurationsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	PrinterDeviceId uuid.UUID `json:"printer_device_id" validate:"required"`
	
	DocumentType string `json:"document_type" validate:"required"`
	
	FilterOrderType *string `json:"filter_order_type"`
	
	FilterKitchenStationId *uuid.UUID `json:"filter_kitchen_station_id"`
	
	FilterProductCategoryId *uuid.UUID `json:"filter_product_category_id"`
	
	FilterCourseId *uuid.UUID `json:"filter_course_id"`
	
	NumberOfCopies *int64 `json:"number_of_copies"`
	
	AutoPrint *bool `json:"auto_print"`
	
	PrintPriority *int64 `json:"print_priority"`
	
	// Duplicate removed: TemplateConfig json.RawMessage `json:"template_config"`
	
	PaperSize *string `json:"paper_size"`
	
	PrintOrientation *string `json:"print_orientation"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'receipt', *string `json:"'receipt',"`
	
	'report', *string `json:"'report',"`
	
}

// Validate validates the create request
func (r *CreatePrinterConfigurationsRequest) Validate() error {
	
	if r.PrinterDeviceId == uuid.Nil {
		return fmt.Errorf("printer_device_id is required")
	}
	
	if r.DocumentType == "" {
		return fmt.Errorf("document_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePrinterConfigurationsRequest represents a request to update a printer_configurations
type UpdatePrinterConfigurationsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	PrinterDeviceId *uuid.UUID `json:"printer_device_id,omitempty" validate:"omitempty,required"`
	
	DocumentType *string `json:"document_type,omitempty" validate:"omitempty,required"`
	
	FilterOrderType *string `json:"filter_order_type,omitempty"`
	
	FilterKitchenStationId *uuid.UUID `json:"filter_kitchen_station_id,omitempty"`
	
	FilterProductCategoryId *uuid.UUID `json:"filter_product_category_id,omitempty"`
	
	FilterCourseId *uuid.UUID `json:"filter_course_id,omitempty"`
	
	NumberOfCopies *int64 `json:"number_of_copies,omitempty"`
	
	AutoPrint *bool `json:"auto_print,omitempty"`
	
	PrintPriority *int64 `json:"print_priority,omitempty"`
	
	TemplateConfig *json.RawMessage `json:"template_config,omitempty"`
	
	PaperSize *string `json:"paper_size,omitempty"`
	
	PrintOrientation *string `json:"print_orientation,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'receipt', *string `json:"'receipt',,omitempty"`
	
	'report', *string `json:"'report',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePrinterConfigurationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.PrinterDeviceId != nil {
		hasUpdate = true
	}
	
	if r.DocumentType != nil {
		hasUpdate = true
	}
	
	if r.FilterOrderType != nil {
		hasUpdate = true
	}
	
	if r.FilterKitchenStationId != nil {
		hasUpdate = true
	}
	
	if r.FilterProductCategoryId != nil {
		hasUpdate = true
	}
	
	if r.FilterCourseId != nil {
		hasUpdate = true
	}
	
	if r.NumberOfCopies != nil {
		hasUpdate = true
	}
	
	if r.AutoPrint != nil {
		hasUpdate = true
	}
	
	if r.PrintPriority != nil {
		hasUpdate = true
	}
	
	if r.TemplateConfig != nil {
		hasUpdate = true
	}
	
	if r.PaperSize != nil {
		hasUpdate = true
	}
	
	if r.PrintOrientation != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	
	if r.'receipt', != nil {
		hasUpdate = true
	}
	
	if r.'report', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PrinterConfigurationsListResponse represents a paginated list of printer_configurations records
type PrinterConfigurationsListResponse struct {
	Items      []*PrinterConfigurationsResponse `json:"items"`
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

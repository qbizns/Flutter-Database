package pos_account_mapping

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PosAccountMappingsResponse represents a pos_account_mappings response
type PosAccountMappingsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SourceType string `json:"source_type"`
	
	'product', *string `json:"'product',"`
	
	'category', *string `json:"'category',"`
	
	'paymentMethod', *string `json:"'payment_method',"`
	
	'salesChannel', *string `json:"'sales_channel',"`
	
	'discount', *string `json:"'discount',"`
	
	'rounding', *string `json:"'rounding',"`
	
	'tax', *string `json:"'tax',"`
	
	'serviceCharge', *string `json:"'service_charge',"`
	
	'shipping', *string `json:"'shipping',"`
	
	'giftCard', *string `json:"'gift_card',"`
	
	'storeCredit', *string `json:"'store_credit',"`
	
	'loyaltyRedemption',-- *string `json:"'loyalty_redemption',--"`
	
	'default' *string `json:"'default'"`
	
	SourceId *uuid.UUID `json:"source_id"`
	
	SourceCode *string `json:"source_code"`
	
	Purpose string `json:"purpose"`
	
	'revenue', *string `json:"'revenue',"`
	
	'cogs', *string `json:"'cogs',"`
	
	'inventory', *string `json:"'inventory',"`
	
	'expense', *string `json:"'expense',"`
	
	'liability', *string `json:"'liability',"`
	
	'asset', *string `json:"'asset',"`
	
	'discountExpense', *string `json:"'discount_expense',"`
	
	'discountContra', *string `json:"'discount_contra',"`
	
	'taxLiability', *string `json:"'tax_liability',"`
	
	'rounding', *string `json:"'rounding',"`
	
	'clearing' *string `json:"'clearing'"`
	
	AccountId uuid.UUID `json:"account_id"`
	
	IsDefault *bool `json:"is_default"`
	
	IsActive *bool `json:"is_active"`
	
	Priority *int64 `json:"priority"`
	
	Conditions json.RawMessage `json:"conditions"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(sourceType *string `json:"(source_type"`
	
	(sourceType string `json:"(source_type"`
	
	EffectiveFrom *string `json:"effective_from"`
	
}

// CreatePosAccountMappingsRequest represents a request to create a pos_account_mappings
type CreatePosAccountMappingsRequest struct {
	
	SourceType string `json:"source_type" validate:"required"`
	
	'product', *string `json:"'product',"`
	
	'category', *string `json:"'category',"`
	
	'paymentMethod', *string `json:"'payment_method',"`
	
	'salesChannel', *string `json:"'sales_channel',"`
	
	'discount', *string `json:"'discount',"`
	
	'rounding', *string `json:"'rounding',"`
	
	'tax', *string `json:"'tax',"`
	
	'serviceCharge', *string `json:"'service_charge',"`
	
	'shipping', *string `json:"'shipping',"`
	
	'giftCard', *string `json:"'gift_card',"`
	
	'storeCredit', *string `json:"'store_credit',"`
	
	'loyaltyRedemption',-- *string `json:"'loyalty_redemption',--"`
	
	'default' *string `json:"'default'"`
	
	SourceId *uuid.UUID `json:"source_id"`
	
	SourceCode *string `json:"source_code"`
	
	Purpose string `json:"purpose" validate:"required"`
	
	'revenue', *string `json:"'revenue',"`
	
	'cogs', *string `json:"'cogs',"`
	
	'inventory', *string `json:"'inventory',"`
	
	'expense', *string `json:"'expense',"`
	
	'liability', *string `json:"'liability',"`
	
	'asset', *string `json:"'asset',"`
	
	'discountExpense', *string `json:"'discount_expense',"`
	
	'discountContra', *string `json:"'discount_contra',"`
	
	'taxLiability', *string `json:"'tax_liability',"`
	
	'rounding', *string `json:"'rounding',"`
	
	'clearing' *string `json:"'clearing'"`
	
	AccountId uuid.UUID `json:"account_id" validate:"required"`
	
	IsDefault *bool `json:"is_default"`
	
	IsActive *bool `json:"is_active"`
	
	Priority *int64 `json:"priority"`
	
	// Duplicate removed: Conditions json.RawMessage `json:"conditions"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(sourceType *string `json:"(source_type"`
	
	(sourceType string `json:"(source_type" validate:"required"`
	
	EffectiveFrom *string `json:"effective_from"`
	
}

// Validate validates the create request
func (r *CreatePosAccountMappingsRequest) Validate() error {
	
	if r.SourceType == "" {
		return fmt.Errorf("source_type is required")
	}
	
	if r.Purpose == "" {
		return fmt.Errorf("purpose is required")
	}
	
	if r.AccountId == uuid.Nil {
		return fmt.Errorf("account_id is required")
	}
	
	if r.(sourceType == "" {
		return fmt.Errorf("(source_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePosAccountMappingsRequest represents a request to update a pos_account_mappings
type UpdatePosAccountMappingsRequest struct {
	
	SourceType *string `json:"source_type,omitempty" validate:"omitempty,required"`
	
	'product', *string `json:"'product',,omitempty"`
	
	'category', *string `json:"'category',,omitempty"`
	
	'paymentMethod', *string `json:"'payment_method',,omitempty"`
	
	'salesChannel', *string `json:"'sales_channel',,omitempty"`
	
	'discount', *string `json:"'discount',,omitempty"`
	
	'rounding', *string `json:"'rounding',,omitempty"`
	
	'tax', *string `json:"'tax',,omitempty"`
	
	'serviceCharge', *string `json:"'service_charge',,omitempty"`
	
	'shipping', *string `json:"'shipping',,omitempty"`
	
	'giftCard', *string `json:"'gift_card',,omitempty"`
	
	'storeCredit', *string `json:"'store_credit',,omitempty"`
	
	'loyaltyRedemption',-- *string `json:"'loyalty_redemption',--,omitempty"`
	
	'default' *string `json:"'default',omitempty"`
	
	SourceId *uuid.UUID `json:"source_id,omitempty"`
	
	SourceCode *string `json:"source_code,omitempty"`
	
	Purpose *string `json:"purpose,omitempty" validate:"omitempty,required"`
	
	'revenue', *string `json:"'revenue',,omitempty"`
	
	'cogs', *string `json:"'cogs',,omitempty"`
	
	'inventory', *string `json:"'inventory',,omitempty"`
	
	'expense', *string `json:"'expense',,omitempty"`
	
	'liability', *string `json:"'liability',,omitempty"`
	
	'asset', *string `json:"'asset',,omitempty"`
	
	'discountExpense', *string `json:"'discount_expense',,omitempty"`
	
	'discountContra', *string `json:"'discount_contra',,omitempty"`
	
	'taxLiability', *string `json:"'tax_liability',,omitempty"`
	
	'rounding', *string `json:"'rounding',,omitempty"`
	
	'clearing' *string `json:"'clearing',omitempty"`
	
	AccountId *uuid.UUID `json:"account_id,omitempty" validate:"omitempty,required"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	Conditions *json.RawMessage `json:"conditions,omitempty"`
	
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	
	EffectiveTo *time.Time `json:"effective_to,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	(sourceType *string `json:"(source_type,omitempty"`
	
	(sourceType *string `json:"(source_type,omitempty" validate:"omitempty,required"`
	
	EffectiveFrom *string `json:"effective_from,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePosAccountMappingsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SourceType != nil {
		hasUpdate = true
	}
	
	if r.'product', != nil {
		hasUpdate = true
	}
	
	if r.'category', != nil {
		hasUpdate = true
	}
	
	if r.'paymentMethod', != nil {
		hasUpdate = true
	}
	
	if r.'salesChannel', != nil {
		hasUpdate = true
	}
	
	if r.'discount', != nil {
		hasUpdate = true
	}
	
	if r.'rounding', != nil {
		hasUpdate = true
	}
	
	if r.'tax', != nil {
		hasUpdate = true
	}
	
	if r.'serviceCharge', != nil {
		hasUpdate = true
	}
	
	if r.'shipping', != nil {
		hasUpdate = true
	}
	
	if r.'giftCard', != nil {
		hasUpdate = true
	}
	
	if r.'storeCredit', != nil {
		hasUpdate = true
	}
	
	if r.'loyaltyRedemption',-- != nil {
		hasUpdate = true
	}
	
	if r.'default' != nil {
		hasUpdate = true
	}
	
	if r.SourceId != nil {
		hasUpdate = true
	}
	
	if r.SourceCode != nil {
		hasUpdate = true
	}
	
	if r.Purpose != nil {
		hasUpdate = true
	}
	
	if r.'revenue', != nil {
		hasUpdate = true
	}
	
	if r.'cogs', != nil {
		hasUpdate = true
	}
	
	if r.'inventory', != nil {
		hasUpdate = true
	}
	
	if r.'expense', != nil {
		hasUpdate = true
	}
	
	if r.'liability', != nil {
		hasUpdate = true
	}
	
	if r.'asset', != nil {
		hasUpdate = true
	}
	
	if r.'discountExpense', != nil {
		hasUpdate = true
	}
	
	if r.'discountContra', != nil {
		hasUpdate = true
	}
	
	if r.'taxLiability', != nil {
		hasUpdate = true
	}
	
	if r.'rounding', != nil {
		hasUpdate = true
	}
	
	if r.'clearing' != nil {
		hasUpdate = true
	}
	
	if r.AccountId != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.Conditions != nil {
		hasUpdate = true
	}
	
	if r.EffectiveFrom != nil {
		hasUpdate = true
	}
	
	if r.EffectiveTo != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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
	
	if r.(sourceType != nil {
		hasUpdate = true
	}
	
	if r.(sourceType != nil {
		hasUpdate = true
	}
	
	if r.EffectiveFrom != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PosAccountMappingsListResponse represents a paginated list of pos_account_mappings records
type PosAccountMappingsListResponse struct {
	Items      []*PosAccountMappingsResponse `json:"items"`
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

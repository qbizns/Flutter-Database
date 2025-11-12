package reconciliation_rule_model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReconciliationRuleModelsResponse represents a reconciliation_rule_models response
type ReconciliationRuleModelsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	RuleName string `json:"rule_name"`
	
	RuleCode *string `json:"rule_code"`
	
	Sequence *int64 `json:"sequence"`
	
	AmountMin *float64 `json:"amount_min"`
	
	AmountMax *float64 `json:"amount_max"`
	
	DescriptionPattern *string `json:"description_pattern"`
	
	CounterpartyPattern *string `json:"counterparty_pattern"`
	
	ReferencePattern *string `json:"reference_pattern"`
	
	JournalId *uuid.UUID `json:"journal_id"`
	
	AccountId *uuid.UUID `json:"account_id"`
	
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id"`
	
	TaxId *uuid.UUID `json:"tax_id"`
	
	IsActive *bool `json:"is_active"`
	
	AutoApply *bool `json:"auto_apply"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateReconciliationRuleModelsRequest represents a request to create a reconciliation_rule_models
type CreateReconciliationRuleModelsRequest struct {
	
	RuleName string `json:"rule_name" validate:"required"`
	
	RuleCode *string `json:"rule_code"`
	
	Sequence *int64 `json:"sequence"`
	
	AmountMin *float64 `json:"amount_min"`
	
	AmountMax *float64 `json:"amount_max"`
	
	DescriptionPattern *string `json:"description_pattern"`
	
	CounterpartyPattern *string `json:"counterparty_pattern"`
	
	ReferencePattern *string `json:"reference_pattern"`
	
	JournalId *uuid.UUID `json:"journal_id"`
	
	AccountId *uuid.UUID `json:"account_id"`
	
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id"`
	
	TaxId *uuid.UUID `json:"tax_id"`
	
	IsActive *bool `json:"is_active"`
	
	AutoApply *bool `json:"auto_apply"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateReconciliationRuleModelsRequest) Validate() error {
	
	if r.RuleName == "" {
		return fmt.Errorf("rule_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateReconciliationRuleModelsRequest represents a request to update a reconciliation_rule_models
type UpdateReconciliationRuleModelsRequest struct {
	
	RuleName *string `json:"rule_name,omitempty" validate:"omitempty,required"`
	
	RuleCode *string `json:"rule_code,omitempty"`
	
	Sequence *int64 `json:"sequence,omitempty"`
	
	AmountMin *float64 `json:"amount_min,omitempty"`
	
	AmountMax *float64 `json:"amount_max,omitempty"`
	
	DescriptionPattern *string `json:"description_pattern,omitempty"`
	
	CounterpartyPattern *string `json:"counterparty_pattern,omitempty"`
	
	ReferencePattern *string `json:"reference_pattern,omitempty"`
	
	JournalId *uuid.UUID `json:"journal_id,omitempty"`
	
	AccountId *uuid.UUID `json:"account_id,omitempty"`
	
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id,omitempty"`
	
	TaxId *uuid.UUID `json:"tax_id,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	AutoApply *bool `json:"auto_apply,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateReconciliationRuleModelsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.RuleName != nil {
		hasUpdate = true
	}
	
	if r.RuleCode != nil {
		hasUpdate = true
	}
	
	if r.Sequence != nil {
		hasUpdate = true
	}
	
	if r.AmountMin != nil {
		hasUpdate = true
	}
	
	if r.AmountMax != nil {
		hasUpdate = true
	}
	
	if r.DescriptionPattern != nil {
		hasUpdate = true
	}
	
	if r.CounterpartyPattern != nil {
		hasUpdate = true
	}
	
	if r.ReferencePattern != nil {
		hasUpdate = true
	}
	
	if r.JournalId != nil {
		hasUpdate = true
	}
	
	if r.AccountId != nil {
		hasUpdate = true
	}
	
	if r.AnalyticAccountId != nil {
		hasUpdate = true
	}
	
	if r.TaxId != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.AutoApply != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ReconciliationRuleModelsListResponse represents a paginated list of reconciliation_rule_models records
type ReconciliationRuleModelsListResponse struct {
	Items      []*ReconciliationRuleModelsResponse `json:"items"`
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

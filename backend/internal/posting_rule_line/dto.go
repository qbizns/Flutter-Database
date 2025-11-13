package posting_rule_line

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingRuleLinesResponse represents a posting_rule_lines response
type PostingRuleLinesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	PostingRuleId uuid.UUID `json:"posting_rule_id"`
	
	LineNo int64 `json:"line_no"`
	
	Side string `json:"side"`
	
	ConceptKey *string `json:"concept_key"`
	
	AccountSource string `json:"account_source"`
	
	FixedAccountId *uuid.UUID `json:"fixed_account_id"`
	
	AccountFieldPath *string `json:"account_field_path"`
	
	AccountExpression *string `json:"account_expression"`
	
	AmountSource string `json:"amount_source"`
	
	AmountFieldPath *string `json:"amount_field_path"`
	
	AmountExpression *string `json:"amount_expression"`
	
	MappingContext json.RawMessage `json:"mapping_context"`
	
	DescriptionTemplate *string `json:"description_template"`
	
	IsActive bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	(accountSource string `json:"(account_source"`
	
	(accountSource string `json:"(account_source"`
	
	(accountSource string `json:"(account_source"`
	
	(accountSource string `json:"(account_source"`
	
	(amountSource string `json:"(amount_source"`
	
	(amountSource string `json:"(amount_source"`
	
	(amountSource string `json:"(amount_source"`
	
}

// CreatePostingRuleLinesRequest represents a request to create a posting_rule_lines
type CreatePostingRuleLinesRequest struct {
	
	PostingRuleId uuid.UUID `json:"posting_rule_id" validate:"required"`
	
	LineNo int64 `json:"line_no" validate:"required"`
	
	Side string `json:"side" validate:"required"`
	
	ConceptKey *string `json:"concept_key"`
	
	AccountSource string `json:"account_source" validate:"required"`
	
	FixedAccountId *uuid.UUID `json:"fixed_account_id"`
	
	AccountFieldPath *string `json:"account_field_path"`
	
	AccountExpression *string `json:"account_expression"`
	
	AmountSource string `json:"amount_source" validate:"required"`
	
	AmountFieldPath *string `json:"amount_field_path"`
	
	AmountExpression *string `json:"amount_expression"`
	
	// Duplicate removed: MappingContext json.RawMessage `json:"mapping_context"`
	
	DescriptionTemplate *string `json:"description_template"`
	
	IsActive bool `json:"is_active" validate:"required"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	(accountSource string `json:"(account_source" validate:"required"`
	
	(accountSource string `json:"(account_source" validate:"required"`
	
	(accountSource string `json:"(account_source" validate:"required"`
	
	(accountSource string `json:"(account_source" validate:"required"`
	
	(amountSource string `json:"(amount_source" validate:"required"`
	
	(amountSource string `json:"(amount_source" validate:"required"`
	
	(amountSource string `json:"(amount_source" validate:"required"`
	
}

// Validate validates the create request
func (r *CreatePostingRuleLinesRequest) Validate() error {
	
	if r.PostingRuleId == uuid.Nil {
		return fmt.Errorf("posting_rule_id is required")
	}
	
	if r.LineNo == 0 {
		return fmt.Errorf("line_no is required")
	}
	
	if r.Side == "" {
		return fmt.Errorf("side is required")
	}
	
	if r.AccountSource == "" {
		return fmt.Errorf("account_source is required")
	}
	
	if r.AmountSource == "" {
		return fmt.Errorf("amount_source is required")
	}
	
	if r.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	
	if r.(accountSource == "" {
		return fmt.Errorf("(account_source is required")
	}
	
	if r.(accountSource == "" {
		return fmt.Errorf("(account_source is required")
	}
	
	if r.(accountSource == "" {
		return fmt.Errorf("(account_source is required")
	}
	
	if r.(accountSource == "" {
		return fmt.Errorf("(account_source is required")
	}
	
	if r.(amountSource == "" {
		return fmt.Errorf("(amount_source is required")
	}
	
	if r.(amountSource == "" {
		return fmt.Errorf("(amount_source is required")
	}
	
	if r.(amountSource == "" {
		return fmt.Errorf("(amount_source is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingRuleLinesRequest represents a request to update a posting_rule_lines
type UpdatePostingRuleLinesRequest struct {
	
	PostingRuleId *uuid.UUID `json:"posting_rule_id,omitempty" validate:"omitempty,required"`
	
	LineNo *int64 `json:"line_no,omitempty" validate:"omitempty,required"`
	
	Side *string `json:"side,omitempty" validate:"omitempty,required"`
	
	ConceptKey *string `json:"concept_key,omitempty"`
	
	AccountSource *string `json:"account_source,omitempty" validate:"omitempty,required"`
	
	FixedAccountId *uuid.UUID `json:"fixed_account_id,omitempty"`
	
	AccountFieldPath *string `json:"account_field_path,omitempty"`
	
	AccountExpression *string `json:"account_expression,omitempty"`
	
	AmountSource *string `json:"amount_source,omitempty" validate:"omitempty,required"`
	
	AmountFieldPath *string `json:"amount_field_path,omitempty"`
	
	AmountExpression *string `json:"amount_expression,omitempty"`
	
	MappingContext *json.RawMessage `json:"mapping_context,omitempty"`
	
	DescriptionTemplate *string `json:"description_template,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty,required"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	(accountSource *string `json:"(account_source,omitempty" validate:"omitempty,required"`
	
	(accountSource *string `json:"(account_source,omitempty" validate:"omitempty,required"`
	
	(accountSource *string `json:"(account_source,omitempty" validate:"omitempty,required"`
	
	(accountSource *string `json:"(account_source,omitempty" validate:"omitempty,required"`
	
	(amountSource *string `json:"(amount_source,omitempty" validate:"omitempty,required"`
	
	(amountSource *string `json:"(amount_source,omitempty" validate:"omitempty,required"`
	
	(amountSource *string `json:"(amount_source,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdatePostingRuleLinesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PostingRuleId != nil {
		hasUpdate = true
	}
	
	if r.LineNo != nil {
		hasUpdate = true
	}
	
	if r.Side != nil {
		hasUpdate = true
	}
	
	if r.ConceptKey != nil {
		hasUpdate = true
	}
	
	if r.AccountSource != nil {
		hasUpdate = true
	}
	
	if r.FixedAccountId != nil {
		hasUpdate = true
	}
	
	if r.AccountFieldPath != nil {
		hasUpdate = true
	}
	
	if r.AccountExpression != nil {
		hasUpdate = true
	}
	
	if r.AmountSource != nil {
		hasUpdate = true
	}
	
	if r.AmountFieldPath != nil {
		hasUpdate = true
	}
	
	if r.AmountExpression != nil {
		hasUpdate = true
	}
	
	if r.MappingContext != nil {
		hasUpdate = true
	}
	
	if r.DescriptionTemplate != nil {
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
	
	if r.(accountSource != nil {
		hasUpdate = true
	}
	
	if r.(accountSource != nil {
		hasUpdate = true
	}
	
	if r.(accountSource != nil {
		hasUpdate = true
	}
	
	if r.(accountSource != nil {
		hasUpdate = true
	}
	
	if r.(amountSource != nil {
		hasUpdate = true
	}
	
	if r.(amountSource != nil {
		hasUpdate = true
	}
	
	if r.(amountSource != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PostingRuleLinesListResponse represents a paginated list of posting_rule_lines records
type PostingRuleLinesListResponse struct {
	Items      []*PostingRuleLinesResponse `json:"items"`
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

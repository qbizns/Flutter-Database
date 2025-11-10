package posting

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Engine is the posting engine service
type Engine struct {
	repo   Repository
	logger *logging.Logger
}

// NewEngine creates a new posting engine
func NewEngine(repo Repository, logger *logging.Logger) *Engine {
	return &Engine{
		repo:   repo,
		logger: logger,
	}
}

// Post processes a document and creates journal entries
func (e *Engine) Post(ctx context.Context, input PostingInput) error {
	e.logger.Info("posting document",
		zap.String("document_type", input.DocumentType),
		zap.String("document_id", input.DocumentID.String()),
		zap.String("event", input.Event),
	)

	// 1. Load document data
	doc, err := e.repo.LoadDocument(ctx, input.DocumentType, input.DocumentID)
	if err != nil {
		return apperrors.DatabaseError(err)
	}

	// 2. Get posting rules
	rules, err := e.repo.GetPostingRules(ctx, input.OrganizationID, input.DocumentType, input.Event)
	if err != nil {
		return apperrors.DatabaseError(err)
	}

	if len(rules) == 0 {
		return apperrors.PostingFailed("no posting rules found for document type")
	}

	// 3. Find matching rule
	var matchedRule *PostingRule
	for i := range rules {
		if e.evaluateCondition(rules[i].ConditionExpression, doc) {
			matchedRule = &rules[i]
			break
		}
	}

	if matchedRule == nil {
		return apperrors.PostingFailed("no matching posting rule found")
	}

	e.logger.Info("matched posting rule",
		zap.String("rule_code", matchedRule.RuleCode),
	)

	// 4. Build journal entry
	je, err := e.buildJournalEntry(ctx, *matchedRule, doc, input)
	if err != nil {
		return err
	}

	// 5. Validate journal entry
	if err := e.validate(ctx, input, je); err != nil {
		return err
	}

	// 6. Create journal entry + post to GL
	if err := e.repo.CreateJournalEntry(ctx, je); err != nil {
		return apperrors.DatabaseError(err)
	}

	// 7. Update document posting status
	if err := e.repo.UpdateDocumentPostingStatus(ctx, input.DocumentType, input.DocumentID, je.ID, "posted"); err != nil {
		return apperrors.DatabaseError(err)
	}

	// 8. Log to posting audit
	if err := e.repo.LogPostingAudit(ctx, PostingAudit{
		OrganizationID:   input.OrganizationID,
		SourceTable:      input.DocumentType,
		SourceID:         input.DocumentID,
		JournalEntryID:   je.ID,
		PostingStatus:    "posted",
		Event:            input.Event,
	}); err != nil {
		e.logger.Error("failed to log posting audit", zap.Error(err))
		// Don't fail the entire posting for audit log failures
	}

	e.logger.Info("document posted successfully",
		zap.String("journal_entry_id", je.ID.String()),
	)

	return nil
}

// evaluateCondition evaluates a DSL expression
func (e *Engine) evaluateCondition(condition string, doc map[string]interface{}) bool {
	if condition == "" {
		return true // No condition = always match
	}

	env := map[string]interface{}{
		"doc": doc,
	}

	program, err := expr.Compile(condition, expr.Env(env))
	if err != nil {
		e.logger.Error("failed to compile condition", zap.Error(err), zap.String("condition", condition))
		return false
	}

	result, err := expr.Run(program, env)
	if err != nil {
		e.logger.Error("failed to evaluate condition", zap.Error(err), zap.String("condition", condition))
		return false
	}

	matched, ok := result.(bool)
	if !ok {
		e.logger.Error("condition did not return boolean", zap.String("condition", condition))
		return false
	}

	return matched
}

// buildJournalEntry builds a journal entry from posting rule
func (e *Engine) buildJournalEntry(ctx context.Context, rule PostingRule, doc map[string]interface{}, input PostingInput) (*JournalEntry, error) {
	je := &JournalEntry{
		ID:             uuid.New(),
		OrganizationID: input.OrganizationID,
		Description:    fmt.Sprintf("Auto-posted from %s", input.DocumentType),
		Status:         "posted",
		Lines:          []JournalEntryLine{},
	}

	// Get rule lines
	lines, err := e.repo.GetPostingRuleLines(ctx, rule.ID)
	if err != nil {
		return nil, apperrors.DatabaseError(err)
	}

	var totalDebit, totalCredit float64

	for _, ruleLine := range lines {
		// Resolve account
		accountID, err := e.resolveAccount(ctx, input.OrganizationID, ruleLine.ConceptKey, doc)
		if err != nil {
			return nil, err
		}

		// Get amount
		amount, err := e.getAmount(ruleLine, doc)
		if err != nil {
			return nil, err
		}

		// Create line
		line := JournalEntryLine{
			ID:          uuid.New(),
			AccountID:   accountID,
			Description: e.interpolateDescription(ruleLine.DescriptionTemplate, doc),
		}

		if ruleLine.Side == "debit" {
			line.DebitAmount = amount
			totalDebit += amount
		} else {
			line.CreditAmount = amount
			totalCredit += amount
		}

		je.Lines = append(je.Lines, line)
	}

	je.TotalDebit = totalDebit
	je.TotalCredit = totalCredit

	return je, nil
}

// resolveAccount resolves a GL account from concept mapping
func (e *Engine) resolveAccount(ctx context.Context, orgID uuid.UUID, conceptKey string, doc map[string]interface{}) (uuid.UUID, error) {
	// Try to resolve account via concept mapping
	accountID, err := e.repo.ResolveAccountFromConcept(ctx, orgID, conceptKey)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, apperrors.PostingFailed(fmt.Sprintf("no account mapped for concept: %s", conceptKey))
		}
		return uuid.Nil, apperrors.DatabaseError(err)
	}

	return accountID, nil
}

// getAmount extracts amount from document based on rule
func (e *Engine) getAmount(ruleLine PostingRuleLine, doc map[string]interface{}) (float64, error) {
	switch ruleLine.AmountSource {
	case "document_field":
		// Extract field from doc
		fieldPath := ruleLine.AmountFieldPath
		value, ok := doc[fieldPath]
		if !ok {
			return 0, apperrors.PostingFailed(fmt.Sprintf("field not found: %s", fieldPath))
		}

		// Convert to float64
		switch v := value.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case json.Number:
			f, err := v.Float64()
			if err != nil {
				return 0, apperrors.PostingFailed(fmt.Sprintf("invalid number format for field: %s", fieldPath))
			}
			return f, nil
		default:
			return 0, apperrors.PostingFailed(fmt.Sprintf("invalid type for field: %s", fieldPath))
		}

	case "expression":
		// Evaluate expression (DSL)
		env := map[string]interface{}{
			"doc": doc,
		}

		program, err := expr.Compile(ruleLine.AmountExpression, expr.Env(env))
		if err != nil {
			return 0, apperrors.PostingFailed(fmt.Sprintf("failed to compile amount expression: %v", err))
		}

		result, err := expr.Run(program, env)
		if err != nil {
			return 0, apperrors.PostingFailed(fmt.Sprintf("failed to evaluate amount expression: %v", err))
		}

		amount, ok := result.(float64)
		if !ok {
			return 0, apperrors.PostingFailed("amount expression did not return a number")
		}

		return amount, nil

	default:
		return 0, apperrors.PostingFailed(fmt.Sprintf("unsupported amount source: %s", ruleLine.AmountSource))
	}
}

// interpolateDescription replaces placeholders in description template
func (e *Engine) interpolateDescription(template string, doc map[string]interface{}) string {
	// Simple placeholder replacement (can be enhanced)
	// For now, just return template
	// TODO: Implement proper template interpolation
	return template
}

// validate validates the journal entry
func (e *Engine) validate(ctx context.Context, input PostingInput, je *JournalEntry) error {
	// Get validation rules
	rules, err := e.repo.GetValidationRules(ctx, input.OrganizationID, input.DocumentType, input.Event)
	if err != nil {
		return apperrors.DatabaseError(err)
	}

	var blockingErrors []string
	var warnings []string

	for _, rule := range rules {
		// Evaluate validation expression
		env := map[string]interface{}{
			"je": map[string]interface{}{
				"total_debit":  je.TotalDebit,
				"total_credit": je.TotalCredit,
			},
		}

		program, err := expr.Compile(rule.Expression, expr.Env(env))
		if err != nil {
			e.logger.Error("failed to compile validation expression", zap.Error(err))
			continue
		}

		result, err := expr.Run(program, env)
		if err != nil {
			e.logger.Error("failed to evaluate validation expression", zap.Error(err))
			continue
		}

		passed, ok := result.(bool)
		if !ok {
			e.logger.Error("validation expression did not return boolean")
			continue
		}

		if !passed {
			message := fmt.Sprintf("Validation failed: %s", rule.MessageTemplate)
			if rule.IsBlocking {
				blockingErrors = append(blockingErrors, message)
			} else if rule.Severity == "warning" {
				warnings = append(warnings, message)
			}

			// Log validation result
			_ = e.repo.LogValidationResult(ctx, ValidationResult{
				OrganizationID:    input.OrganizationID,
				DocumentType:      input.DocumentType,
				DocumentID:        input.DocumentID,
				Event:             input.Event,
				ValidationRuleID:  rule.ID,
				Severity:          rule.Severity,
				Message:           message,
				IsBlocking:        rule.IsBlocking,
			})
		}
	}

	// Check for blocking errors
	if len(blockingErrors) > 0 {
		return apperrors.ValidationBlocked(
			"Validation failed",
			map[string]interface{}{
				"errors":   blockingErrors,
				"warnings": warnings,
			},
		)
	}

	// Log warnings
	if len(warnings) > 0 {
		e.logger.Warn("posting validation warnings", zap.Strings("warnings", warnings))
	}

	return nil
}

// PostingInput represents input for posting
type PostingInput struct {
	OrganizationID uuid.UUID
	DocumentType   string
	DocumentID     uuid.UUID
	Event          string
	UserID         uuid.UUID
}

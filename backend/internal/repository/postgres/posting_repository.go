package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/posting"
)

// ============================================================================
// POSTING ENGINE REPOSITORY
// ============================================================================

// PostingRepository implements posting.Repository interface
type PostingRepository struct {
	db *DB
}

// NewPostingRepository creates a new posting repository
func NewPostingRepository(db *DB) posting.Repository {
	return &PostingRepository{db: db}
}

// ============================================================================
// POSTING CONCEPT REPOSITORY
// ============================================================================

type PostingConceptRepository struct {
	db *DB
}

func NewPostingConceptRepository(db *DB) *PostingConceptRepository {
	return &PostingConceptRepository{db: db}
}

type PostingConcept struct {
	ConceptKey              string
	DefaultLabel            string
	DefaultDescription      *string
	ExpectedAccountTypeID   *uuid.UUID
	NormalSide              *string
	ExampleCode             *string
	ExampleAccountName      *string
	IsSystem                bool
	ConceptCategory         *string
	SortOrder               int
	Notes                   *string
	Metadata                json.RawMessage
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}

type PostingConceptOverride struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	ConceptKey     string
	Label          string
	Description    *string
	IsActive       bool
	Notes          *string
	Metadata       json.RawMessage
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
	CreatedBy      *uuid.UUID
	UpdatedBy      *uuid.UUID
}

func (r *PostingConceptRepository) GetConcept(ctx context.Context, conceptKey string) (*PostingConcept, error) {
	query := `
		SELECT concept_key, default_label, default_description, expected_account_type_id,
		       normal_side, example_code, example_account_name, is_system, concept_category,
		       sort_order, notes, metadata, created_at, updated_at, deleted_at
		FROM posting_concepts
		WHERE concept_key = $1 AND deleted_at IS NULL
	`

	var pc PostingConcept
	err := r.db.Pool.QueryRow(ctx, query, conceptKey).Scan(
		&pc.ConceptKey, &pc.DefaultLabel, &pc.DefaultDescription, &pc.ExpectedAccountTypeID,
		&pc.NormalSide, &pc.ExampleCode, &pc.ExampleAccountName, &pc.IsSystem, &pc.ConceptCategory,
		&pc.SortOrder, &pc.Notes, &pc.Metadata, &pc.CreatedAt, &pc.UpdatedAt, &pc.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

func (r *PostingConceptRepository) ListConcepts(ctx context.Context) ([]PostingConcept, error) {
	query := `
		SELECT concept_key, default_label, default_description, expected_account_type_id,
		       normal_side, example_code, example_account_name, is_system, concept_category,
		       sort_order, notes, metadata, created_at, updated_at, deleted_at
		FROM posting_concepts
		WHERE deleted_at IS NULL
		ORDER BY sort_order, concept_key
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var concepts []PostingConcept
	for rows.Next() {
		var pc PostingConcept
		err := rows.Scan(
			&pc.ConceptKey, &pc.DefaultLabel, &pc.DefaultDescription, &pc.ExpectedAccountTypeID,
			&pc.NormalSide, &pc.ExampleCode, &pc.ExampleAccountName, &pc.IsSystem, &pc.ConceptCategory,
			&pc.SortOrder, &pc.Notes, &pc.Metadata, &pc.CreatedAt, &pc.UpdatedAt, &pc.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		concepts = append(concepts, pc)
	}

	return concepts, rows.Err()
}

func (r *PostingConceptRepository) GetConceptOverride(ctx context.Context, orgID uuid.UUID, conceptKey string) (*PostingConceptOverride, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, concept_key, label, description, is_active,
		       notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM posting_concept_overrides
		WHERE organization_id = $1 AND concept_key = $2 AND deleted_at IS NULL
	`

	var pco PostingConceptOverride
	err := r.db.Pool.QueryRow(ctx, query, orgID, conceptKey).Scan(
		&pco.ID, &pco.OrganizationID, &pco.ConceptKey, &pco.Label, &pco.Description, &pco.IsActive,
		&pco.Notes, &pco.Metadata, &pco.CreatedAt, &pco.UpdatedAt, &pco.DeletedAt, &pco.CreatedBy, &pco.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pco, nil
}

func (r *PostingConceptRepository) CreateConceptOverride(ctx context.Context, override *PostingConceptOverride) error {
	if err := r.db.SetOrganizationContext(ctx, override.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO posting_concept_overrides
		(id, organization_id, concept_key, label, description, is_active, notes, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		override.ID, override.OrganizationID, override.ConceptKey, override.Label,
		override.Description, override.IsActive, override.Notes, override.Metadata,
		override.CreatedAt, override.UpdatedAt, override.CreatedBy,
	)
	return err
}

func (r *PostingConceptRepository) UpdateConceptOverride(ctx context.Context, override *PostingConceptOverride) error {
	if err := r.db.SetOrganizationContext(ctx, override.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE posting_concept_overrides
		SET label = $1, description = $2, is_active = $3, notes = $4, metadata = $5, updated_at = $6, updated_by = $7
		WHERE id = $8 AND organization_id = $9 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		override.Label, override.Description, override.IsActive, override.Notes, override.Metadata,
		time.Now(), override.UpdatedBy, override.ID, override.OrganizationID,
	)
	return err
}

// ============================================================================
// POSTING RULE REPOSITORY
// ============================================================================

type PostingRuleRepository struct {
	db *DB
}

func NewPostingRuleRepository(db *DB) *PostingRuleRepository {
	return &PostingRuleRepository{db: db}
}

type PostingRuleDB struct {
	ID                      uuid.UUID
	PostingProfileDocumentID uuid.UUID
	RuleCode                string
	RuleName                string
	Description             *string
	Event                   string // posting_event enum
	Level                   string // posting_level enum
	Priority                int
	ConditionExpression     *string
	IsActive                bool
	Notes                   *string
	Metadata                json.RawMessage
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
	CreatedBy               *uuid.UUID
	UpdatedBy               *uuid.UUID
}

type PostingRuleLineDB struct {
	ID                  uuid.UUID
	PostingRuleID       uuid.UUID
	LineNo              int
	Side                string // posting_side enum
	ConceptKey          *string
	AccountSource       string // posting_account_source enum
	FixedAccountID      *uuid.UUID
	AccountFieldPath    *string
	AccountExpression   *string
	AmountSource        string // posting_amount_source enum
	AmountFieldPath     *string
	AmountExpression    *string
	MappingContext      json.RawMessage
	DescriptionTemplate *string
	IsActive            bool
	Notes               *string
	Metadata            json.RawMessage
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

func (r *PostingRuleRepository) GetRule(ctx context.Context, ruleID uuid.UUID) (*PostingRuleDB, error) {
	query := `
		SELECT id, posting_profile_document_id, rule_code, rule_name, description, event, level,
		       priority, condition_expression, is_active, notes, metadata, created_at, updated_at,
		       deleted_at, created_by, updated_by
		FROM posting_rules
		WHERE id = $1 AND deleted_at IS NULL
	`

	var rule PostingRuleDB
	err := r.db.Pool.QueryRow(ctx, query, ruleID).Scan(
		&rule.ID, &rule.PostingProfileDocumentID, &rule.RuleCode, &rule.RuleName, &rule.Description,
		&rule.Event, &rule.Level, &rule.Priority, &rule.ConditionExpression, &rule.IsActive,
		&rule.Notes, &rule.Metadata, &rule.CreatedAt, &rule.UpdatedAt, &rule.DeletedAt,
		&rule.CreatedBy, &rule.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *PostingRuleRepository) GetRulesForDocument(ctx context.Context, orgID uuid.UUID, documentType string, event string) ([]PostingRuleDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT pr.id, pr.posting_profile_document_id, pr.rule_code, pr.rule_name, pr.description,
		       pr.event, pr.level, pr.priority, pr.condition_expression, pr.is_active,
		       pr.notes, pr.metadata, pr.created_at, pr.updated_at, pr.deleted_at, pr.created_by, pr.updated_by
		FROM posting_rules pr
		INNER JOIN posting_profile_documents ppd ON pr.posting_profile_document_id = ppd.id
		INNER JOIN posting_profiles pp ON ppd.posting_profile_id = pp.id
		INNER JOIN posting_document_types pdt ON ppd.posting_document_type_id = pdt.id
		WHERE pp.organization_id = $1 AND pdt.code = $2 AND pr.event = $3
		  AND pr.is_active = TRUE AND ppd.is_active = TRUE AND pp.is_active = TRUE
		  AND pr.deleted_at IS NULL AND ppd.deleted_at IS NULL AND pp.deleted_at IS NULL
		ORDER BY pr.priority ASC, pr.created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, documentType, event)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []PostingRuleDB
	for rows.Next() {
		var rule PostingRuleDB
		err := rows.Scan(
			&rule.ID, &rule.PostingProfileDocumentID, &rule.RuleCode, &rule.RuleName, &rule.Description,
			&rule.Event, &rule.Level, &rule.Priority, &rule.ConditionExpression, &rule.IsActive,
			&rule.Notes, &rule.Metadata, &rule.CreatedAt, &rule.UpdatedAt, &rule.DeletedAt,
			&rule.CreatedBy, &rule.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, rows.Err()
}

func (r *PostingRuleRepository) CreateRule(ctx context.Context, rule *PostingRuleDB) error {
	query := `
		INSERT INTO posting_rules
		(id, posting_profile_document_id, rule_code, rule_name, description, event, level,
		 priority, condition_expression, is_active, notes, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		rule.ID, rule.PostingProfileDocumentID, rule.RuleCode, rule.RuleName, rule.Description,
		rule.Event, rule.Level, rule.Priority, rule.ConditionExpression, rule.IsActive,
		rule.Notes, rule.Metadata, rule.CreatedAt, rule.UpdatedAt, rule.CreatedBy,
	)
	return err
}

func (r *PostingRuleRepository) GetRuleLines(ctx context.Context, ruleID uuid.UUID) ([]PostingRuleLineDB, error) {
	query := `
		SELECT id, posting_rule_id, line_no, side, concept_key, account_source,
		       fixed_account_id, account_field_path, account_expression, amount_source,
		       amount_field_path, amount_expression, mapping_context, description_template,
		       is_active, notes, metadata, created_at, updated_at, deleted_at
		FROM posting_rule_lines
		WHERE posting_rule_id = $1 AND deleted_at IS NULL
		ORDER BY line_no ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []PostingRuleLineDB
	for rows.Next() {
		var line PostingRuleLineDB
		err := rows.Scan(
			&line.ID, &line.PostingRuleID, &line.LineNo, &line.Side, &line.ConceptKey,
			&line.AccountSource, &line.FixedAccountID, &line.AccountFieldPath,
			&line.AccountExpression, &line.AmountSource, &line.AmountFieldPath,
			&line.AmountExpression, &line.MappingContext, &line.DescriptionTemplate,
			&line.IsActive, &line.Notes, &line.Metadata, &line.CreatedAt,
			&line.UpdatedAt, &line.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, rows.Err()
}

func (r *PostingRuleRepository) CreateRuleLine(ctx context.Context, line *PostingRuleLineDB) error {
	query := `
		INSERT INTO posting_rule_lines
		(id, posting_rule_id, line_no, side, concept_key, account_source, fixed_account_id,
		 account_field_path, account_expression, amount_source, amount_field_path,
		 amount_expression, mapping_context, description_template, is_active, notes,
		 metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		line.ID, line.PostingRuleID, line.LineNo, line.Side, line.ConceptKey, line.AccountSource,
		line.FixedAccountID, line.AccountFieldPath, line.AccountExpression, line.AmountSource,
		line.AmountFieldPath, line.AmountExpression, line.MappingContext, line.DescriptionTemplate,
		line.IsActive, line.Notes, line.Metadata, line.CreatedAt, line.UpdatedAt,
	)
	return err
}

// ============================================================================
// POSTING PROFILE REPOSITORY
// ============================================================================

type PostingProfileRepository struct {
	db *DB
}

func NewPostingProfileRepository(db *DB) *PostingProfileRepository {
	return &PostingProfileRepository{db: db}
}

type PostingProfileDB struct {
	ID                   uuid.UUID
	OrganizationID       uuid.UUID
	Code                 string
	Name                 string
	Description          *string
	IsDefault            bool
	IsActive             bool
	DefaultFiscalYearID  *uuid.UUID
	Notes                *string
	Metadata             json.RawMessage
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
	CreatedBy            *uuid.UUID
	UpdatedBy            *uuid.UUID
}

type PostingDocumentTypeDB struct {
	ID              uuid.UUID
	Code            string
	Name            string
	Description     *string
	SourceSchema    string
	SourceTable     string
	SourcePKColumn  string
	Category        *string
	IsActive        bool
	IsSystem        bool
	Notes           *string
	Metadata        json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type PostingProfileDocumentDB struct {
	ID                        uuid.UUID
	PostingProfileID          uuid.UUID
	PostingDocumentTypeID     uuid.UUID
	IsActive                  bool
	Notes                     *string
	Metadata                  json.RawMessage
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 *time.Time
}

func (r *PostingProfileRepository) GetProfile(ctx context.Context, orgID uuid.UUID, profileID uuid.UUID) (*PostingProfileDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, code, name, description, is_default, is_active,
		       default_fiscal_year_id, notes, metadata, created_at, updated_at, deleted_at,
		       created_by, updated_by
		FROM posting_profiles
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var profile PostingProfileDB
	err := r.db.Pool.QueryRow(ctx, query, orgID, profileID).Scan(
		&profile.ID, &profile.OrganizationID, &profile.Code, &profile.Name, &profile.Description,
		&profile.IsDefault, &profile.IsActive, &profile.DefaultFiscalYearID, &profile.Notes,
		&profile.Metadata, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt,
		&profile.CreatedBy, &profile.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *PostingProfileRepository) GetProfileByCode(ctx context.Context, orgID uuid.UUID, code string) (*PostingProfileDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, code, name, description, is_default, is_active,
		       default_fiscal_year_id, notes, metadata, created_at, updated_at, deleted_at,
		       created_by, updated_by
		FROM posting_profiles
		WHERE organization_id = $1 AND code = $2 AND deleted_at IS NULL
	`

	var profile PostingProfileDB
	err := r.db.Pool.QueryRow(ctx, query, orgID, code).Scan(
		&profile.ID, &profile.OrganizationID, &profile.Code, &profile.Name, &profile.Description,
		&profile.IsDefault, &profile.IsActive, &profile.DefaultFiscalYearID, &profile.Notes,
		&profile.Metadata, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt,
		&profile.CreatedBy, &profile.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *PostingProfileRepository) ListProfiles(ctx context.Context, orgID uuid.UUID) ([]PostingProfileDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, code, name, description, is_default, is_active,
		       default_fiscal_year_id, notes, metadata, created_at, updated_at, deleted_at,
		       created_by, updated_by
		FROM posting_profiles
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY is_default DESC, code ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []PostingProfileDB
	for rows.Next() {
		var p PostingProfileDB
		err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.Code, &p.Name, &p.Description, &p.IsDefault,
			&p.IsActive, &p.DefaultFiscalYearID, &p.Notes, &p.Metadata, &p.CreatedAt,
			&p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	return profiles, rows.Err()
}

func (r *PostingProfileRepository) CreateProfile(ctx context.Context, profile *PostingProfileDB) error {
	if err := r.db.SetOrganizationContext(ctx, profile.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO posting_profiles
		(id, organization_id, code, name, description, is_default, is_active,
		 default_fiscal_year_id, notes, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		profile.ID, profile.OrganizationID, profile.Code, profile.Name, profile.Description,
		profile.IsDefault, profile.IsActive, profile.DefaultFiscalYearID, profile.Notes,
		profile.Metadata, profile.CreatedAt, profile.UpdatedAt, profile.CreatedBy,
	)
	return err
}

func (r *PostingProfileRepository) UpdateProfile(ctx context.Context, profile *PostingProfileDB) error {
	if err := r.db.SetOrganizationContext(ctx, profile.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE posting_profiles
		SET name = $1, description = $2, is_default = $3, is_active = $4,
		    default_fiscal_year_id = $5, notes = $6, metadata = $7,
		    updated_at = $8, updated_by = $9
		WHERE id = $10 AND organization_id = $11 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		profile.Name, profile.Description, profile.IsDefault, profile.IsActive,
		profile.DefaultFiscalYearID, profile.Notes, profile.Metadata,
		time.Now(), profile.UpdatedBy, profile.ID, profile.OrganizationID,
	)
	return err
}

func (r *PostingProfileRepository) GetDocumentType(ctx context.Context, code string) (*PostingDocumentTypeDB, error) {
	query := `
		SELECT id, code, name, description, source_schema, source_table, source_pk_column,
		       category, is_active, is_system, notes, metadata, created_at, updated_at, deleted_at
		FROM posting_document_types
		WHERE code = $1 AND deleted_at IS NULL
	`

	var dt PostingDocumentTypeDB
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&dt.ID, &dt.Code, &dt.Name, &dt.Description, &dt.SourceSchema, &dt.SourceTable,
		&dt.SourcePKColumn, &dt.Category, &dt.IsActive, &dt.IsSystem, &dt.Notes,
		&dt.Metadata, &dt.CreatedAt, &dt.UpdatedAt, &dt.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dt, nil
}

func (r *PostingProfileRepository) CreateProfileDocument(ctx context.Context, ppd *PostingProfileDocumentDB) error {
	query := `
		INSERT INTO posting_profile_documents
		(id, posting_profile_id, posting_document_type_id, is_active, notes, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		ppd.ID, ppd.PostingProfileID, ppd.PostingDocumentTypeID, ppd.IsActive,
		ppd.Notes, ppd.Metadata, ppd.CreatedAt, ppd.UpdatedAt,
	)
	return err
}

// ============================================================================
// POSTING VALIDATION REPOSITORY
// ============================================================================

type PostingValidationRepository struct {
	db *DB
}

func NewPostingValidationRepository(db *DB) *PostingValidationRepository {
	return &PostingValidationRepository{db: db}
}

type PostingValidationRuleDB struct {
	ID                uuid.UUID
	OrganizationID    *uuid.UUID
	DocumentTypeCode  *string
	Event             *string
	Target            string // validation_target enum
	Code              string
	Name              string
	Description       *string
	Expression        string
	Severity          string // validation_severity enum
	IsBlocking        bool
	IsActive          bool
	MessageTemplate   *string
	Priority          int
	Notes             *string
	Metadata          json.RawMessage
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
	CreatedBy         *uuid.UUID
	UpdatedBy         *uuid.UUID
}

type PostingValidationResultDB struct {
	ID                  uuid.UUID
	OrganizationID      uuid.UUID
	DocumentTypeCode    string
	DocumentID          uuid.UUID
	Event               string
	JournalEntryID      *uuid.UUID
	ValidationRuleID    *uuid.UUID
	Severity            string // validation_severity enum
	MessageCode         string
	Message             string
	IsBlocking          bool
	Context             json.RawMessage
	Metadata            json.RawMessage
	CreatedAt           time.Time
	CreatedBy           *uuid.UUID
}

func (r *PostingValidationRepository) GetValidationRules(ctx context.Context, orgID uuid.UUID, documentType string, event string) ([]PostingValidationRuleDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, document_type_code, event, target, code, name,
		       description, expression, severity, is_blocking, is_active, message_template,
		       priority, notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM posting_validation_rules
		WHERE is_active = TRUE AND deleted_at IS NULL
		  AND (organization_id = $1 OR organization_id IS NULL)
		  AND (document_type_code = $2 OR document_type_code IS NULL)
		  AND (event = $3 OR event IS NULL)
		ORDER BY priority ASC, created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, documentType, event)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []PostingValidationRuleDB
	for rows.Next() {
		var rule PostingValidationRuleDB
		err := rows.Scan(
			&rule.ID, &rule.OrganizationID, &rule.DocumentTypeCode, &rule.Event, &rule.Target,
			&rule.Code, &rule.Name, &rule.Description, &rule.Expression, &rule.Severity,
			&rule.IsBlocking, &rule.IsActive, &rule.MessageTemplate, &rule.Priority, &rule.Notes,
			&rule.Metadata, &rule.CreatedAt, &rule.UpdatedAt, &rule.DeletedAt, &rule.CreatedBy, &rule.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, rows.Err()
}

func (r *PostingValidationRepository) CreateValidationRule(ctx context.Context, rule *PostingValidationRuleDB) error {
	if rule.OrganizationID != nil {
		if err := r.db.SetOrganizationContext(ctx, rule.OrganizationID.String()); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO posting_validation_rules
		(id, organization_id, document_type_code, event, target, code, name, description,
		 expression, severity, is_blocking, is_active, message_template, priority, notes,
		 metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		rule.ID, rule.OrganizationID, rule.DocumentTypeCode, rule.Event, rule.Target,
		rule.Code, rule.Name, rule.Description, rule.Expression, rule.Severity,
		rule.IsBlocking, rule.IsActive, rule.MessageTemplate, rule.Priority, rule.Notes,
		rule.Metadata, rule.CreatedAt, rule.UpdatedAt, rule.CreatedBy,
	)
	return err
}

func (r *PostingValidationRepository) LogValidationResult(ctx context.Context, result *PostingValidationResultDB) error {
	if err := r.db.SetOrganizationContext(ctx, result.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO posting_validation_results
		(id, organization_id, document_type_code, document_id, event, journal_entry_id,
		 validation_rule_id, severity, message_code, message, is_blocking, context,
		 metadata, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		result.ID, result.OrganizationID, result.DocumentTypeCode, result.DocumentID,
		result.Event, result.JournalEntryID, result.ValidationRuleID, result.Severity,
		result.MessageCode, result.Message, result.IsBlocking, result.Context,
		result.Metadata, result.CreatedAt, result.CreatedBy,
	)
	return err
}

// ============================================================================
// POS MAPPING REPOSITORY
// ============================================================================

type POSMappingRepository struct {
	db *DB
}

func NewPOSMappingRepository(db *DB) *POSMappingRepository {
	return &POSMappingRepository{db: db}
}

type POSAccountMappingDB struct {
	ID               uuid.UUID
	OrganizationID   uuid.UUID
	SourceType       string
	SourceID         *uuid.UUID
	SourceCode       *string
	Purpose          string
	AccountID        uuid.UUID
	IsDefault        bool
	IsActive         bool
	Priority         int
	Conditions       json.RawMessage
	EffectiveFrom    *time.Time
	EffectiveTo      *time.Time
	Description      *string
	Notes            *string
	Metadata         json.RawMessage
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	CreatedBy        *uuid.UUID
	UpdatedBy        *uuid.UUID
}

type POSPostingAuditDB struct {
	ID                        uuid.UUID
	OrganizationID            uuid.UUID
	SourceTable               string
	SourceID                  uuid.UUID
	SourceReference           *string
	PostingStatus             string
	JournalEntryID            *uuid.UUID
	ReversalJournalEntryID    *uuid.UUID
	PostingDate               *time.Time
	PostedAt                  *time.Time
	PostedBy                  *uuid.UUID
	PostingMethod             string
	ErrorCode                 *string
	ErrorMessage              *string
	ErrorDetails              json.RawMessage
	RetryCount                int
	LastRetryAt               *time.Time
	MaxRetries                int
	ReversedAt                *time.Time
	ReversedBy                *uuid.UUID
	ReversalReason            *string
	TotalDebit                *float64
	TotalCredit               *float64
	LineCount                 *int
	CurrencyCode              *string
	PostingContext            json.RawMessage
	Notes                     *string
	Metadata                  json.RawMessage
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 *time.Time
	CreatedBy                 *uuid.UUID
	UpdatedBy                 *uuid.UUID
}

type POSTaxMappingDB struct {
	ID                           uuid.UUID
	OrganizationID               uuid.UUID
	POSTaxCode                   *string
	TaxCategoryCode              *string
	POSTaxRate                   *float64
	AccountingTaxID              *uuid.UUID
	DefaultTaxAccountID          *uuid.UUID
	DefaultTaxExpenseAccountID   *uuid.UUID
	IsDefault                    bool
	IsActive                     bool
	Priority                     int
	IsInclusive                  bool
	AppliesToSales               bool
	AppliesToPurchases           bool
	EffectiveFrom                *time.Time
	EffectiveTo                  *time.Time
	Description                  *string
	Notes                        *string
	Metadata                     json.RawMessage
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	DeletedAt                    *time.Time
	CreatedBy                    *uuid.UUID
	UpdatedBy                    *uuid.UUID
}

func (r *POSMappingRepository) GetAccountMapping(ctx context.Context, orgID uuid.UUID, mappingID uuid.UUID) (*POSAccountMappingDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_type, source_id, source_code, purpose, account_id,
		       is_default, is_active, priority, conditions, effective_from, effective_to,
		       description, notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM pos_account_mappings
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var m POSAccountMappingDB
	err := r.db.Pool.QueryRow(ctx, query, mappingID, orgID).Scan(
		&m.ID, &m.OrganizationID, &m.SourceType, &m.SourceID, &m.SourceCode, &m.Purpose,
		&m.AccountID, &m.IsDefault, &m.IsActive, &m.Priority, &m.Conditions, &m.EffectiveFrom,
		&m.EffectiveTo, &m.Description, &m.Notes, &m.Metadata, &m.CreatedAt, &m.UpdatedAt,
		&m.DeletedAt, &m.CreatedBy, &m.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *POSMappingRepository) ListAccountMappings(ctx context.Context, orgID uuid.UUID, sourceType string, purpose string) ([]POSAccountMappingDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_type, source_id, source_code, purpose, account_id,
		       is_default, is_active, priority, conditions, effective_from, effective_to,
		       description, notes, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM pos_account_mappings
		WHERE organization_id = $1 AND source_type = $2 AND purpose = $3 AND deleted_at IS NULL
		  AND is_active = TRUE
		ORDER BY is_default DESC, priority DESC, created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, sourceType, purpose)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []POSAccountMappingDB
	for rows.Next() {
		var m POSAccountMappingDB
		err := rows.Scan(
			&m.ID, &m.OrganizationID, &m.SourceType, &m.SourceID, &m.SourceCode, &m.Purpose,
			&m.AccountID, &m.IsDefault, &m.IsActive, &m.Priority, &m.Conditions, &m.EffectiveFrom,
			&m.EffectiveTo, &m.Description, &m.Notes, &m.Metadata, &m.CreatedAt, &m.UpdatedAt,
			&m.DeletedAt, &m.CreatedBy, &m.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, m)
	}

	return mappings, rows.Err()
}

func (r *POSMappingRepository) CreateAccountMapping(ctx context.Context, mapping *POSAccountMappingDB) error {
	if err := r.db.SetOrganizationContext(ctx, mapping.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO pos_account_mappings
		(id, organization_id, source_type, source_id, source_code, purpose, account_id,
		 is_default, is_active, priority, conditions, effective_from, effective_to,
		 description, notes, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		mapping.ID, mapping.OrganizationID, mapping.SourceType, mapping.SourceID, mapping.SourceCode,
		mapping.Purpose, mapping.AccountID, mapping.IsDefault, mapping.IsActive, mapping.Priority,
		mapping.Conditions, mapping.EffectiveFrom, mapping.EffectiveTo, mapping.Description,
		mapping.Notes, mapping.Metadata, mapping.CreatedAt, mapping.UpdatedAt, mapping.CreatedBy,
	)
	return err
}

func (r *POSMappingRepository) UpdateAccountMapping(ctx context.Context, mapping *POSAccountMappingDB) error {
	if err := r.db.SetOrganizationContext(ctx, mapping.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_account_mappings
		SET account_id = $1, is_default = $2, is_active = $3, priority = $4,
		    conditions = $5, effective_from = $6, effective_to = $7,
		    description = $8, notes = $9, metadata = $10, updated_at = $11, updated_by = $12
		WHERE id = $13 AND organization_id = $14 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		mapping.AccountID, mapping.IsDefault, mapping.IsActive, mapping.Priority,
		mapping.Conditions, mapping.EffectiveFrom, mapping.EffectiveTo,
		mapping.Description, mapping.Notes, mapping.Metadata, time.Now(), mapping.UpdatedBy,
		mapping.ID, mapping.OrganizationID,
	)
	return err
}

func (r *POSMappingRepository) GetPostingAudit(ctx context.Context, orgID uuid.UUID, auditID uuid.UUID) (*POSPostingAuditDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, source_reference, posting_status,
		       journal_entry_id, reversal_journal_entry_id, posting_date, posted_at, posted_by,
		       posting_method, error_code, error_message, error_details, retry_count, last_retry_at,
		       max_retries, reversed_at, reversed_by, reversal_reason, total_debit, total_credit,
		       line_count, currency_code, posting_context, notes, metadata, created_at, updated_at,
		       deleted_at, created_by, updated_by
		FROM pos_posting_audit
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var audit POSPostingAuditDB
	err := r.db.Pool.QueryRow(ctx, query, auditID, orgID).Scan(
		&audit.ID, &audit.OrganizationID, &audit.SourceTable, &audit.SourceID, &audit.SourceReference,
		&audit.PostingStatus, &audit.JournalEntryID, &audit.ReversalJournalEntryID, &audit.PostingDate,
		&audit.PostedAt, &audit.PostedBy, &audit.PostingMethod, &audit.ErrorCode, &audit.ErrorMessage,
		&audit.ErrorDetails, &audit.RetryCount, &audit.LastRetryAt, &audit.MaxRetries, &audit.ReversedAt,
		&audit.ReversedBy, &audit.ReversalReason, &audit.TotalDebit, &audit.TotalCredit,
		&audit.LineCount, &audit.CurrencyCode, &audit.PostingContext, &audit.Notes, &audit.Metadata,
		&audit.CreatedAt, &audit.UpdatedAt, &audit.DeletedAt, &audit.CreatedBy, &audit.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &audit, nil
}

func (r *POSMappingRepository) GetPostingAuditBySource(ctx context.Context, orgID uuid.UUID, sourceTable string, sourceID uuid.UUID) (*POSPostingAuditDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, source_reference, posting_status,
		       journal_entry_id, reversal_journal_entry_id, posting_date, posted_at, posted_by,
		       posting_method, error_code, error_message, error_details, retry_count, last_retry_at,
		       max_retries, reversed_at, reversed_by, reversal_reason, total_debit, total_credit,
		       line_count, currency_code, posting_context, notes, metadata, created_at, updated_at,
		       deleted_at, created_by, updated_by
		FROM pos_posting_audit
		WHERE organization_id = $1 AND source_table = $2 AND source_id = $3
		  AND posting_status NOT IN ('cancelled', 'reversed') AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	var audit POSPostingAuditDB
	err := r.db.Pool.QueryRow(ctx, query, orgID, sourceTable, sourceID).Scan(
		&audit.ID, &audit.OrganizationID, &audit.SourceTable, &audit.SourceID, &audit.SourceReference,
		&audit.PostingStatus, &audit.JournalEntryID, &audit.ReversalJournalEntryID, &audit.PostingDate,
		&audit.PostedAt, &audit.PostedBy, &audit.PostingMethod, &audit.ErrorCode, &audit.ErrorMessage,
		&audit.ErrorDetails, &audit.RetryCount, &audit.LastRetryAt, &audit.MaxRetries, &audit.ReversedAt,
		&audit.ReversedBy, &audit.ReversalReason, &audit.TotalDebit, &audit.TotalCredit,
		&audit.LineCount, &audit.CurrencyCode, &audit.PostingContext, &audit.Notes, &audit.Metadata,
		&audit.CreatedAt, &audit.UpdatedAt, &audit.DeletedAt, &audit.CreatedBy, &audit.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &audit, nil
}

func (r *POSMappingRepository) CreatePostingAudit(ctx context.Context, audit *POSPostingAuditDB) error {
	if err := r.db.SetOrganizationContext(ctx, audit.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO pos_posting_audit
		(id, organization_id, source_table, source_id, source_reference, posting_status,
		 journal_entry_id, reversal_journal_entry_id, posting_date, posted_at, posted_by,
		 posting_method, error_code, error_message, error_details, retry_count, last_retry_at,
		 max_retries, reversed_at, reversed_by, reversal_reason, total_debit, total_credit,
		 line_count, currency_code, posting_context, notes, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
		        $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		audit.ID, audit.OrganizationID, audit.SourceTable, audit.SourceID, audit.SourceReference,
		audit.PostingStatus, audit.JournalEntryID, audit.ReversalJournalEntryID, audit.PostingDate,
		audit.PostedAt, audit.PostedBy, audit.PostingMethod, audit.ErrorCode, audit.ErrorMessage,
		audit.ErrorDetails, audit.RetryCount, audit.LastRetryAt, audit.MaxRetries, audit.ReversedAt,
		audit.ReversedBy, audit.ReversalReason, audit.TotalDebit, audit.TotalCredit,
		audit.LineCount, audit.CurrencyCode, audit.PostingContext, audit.Notes, audit.Metadata,
		audit.CreatedAt, audit.UpdatedAt, audit.CreatedBy,
	)
	return err
}

func (r *POSMappingRepository) UpdatePostingAudit(ctx context.Context, audit *POSPostingAuditDB) error {
	if err := r.db.SetOrganizationContext(ctx, audit.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_posting_audit
		SET posting_status = $1, journal_entry_id = $2, reversal_journal_entry_id = $3,
		    posted_at = $4, posted_by = $5, error_code = $6, error_message = $7,
		    error_details = $8, retry_count = $9, last_retry_at = $10, reversed_at = $11,
		    reversed_by = $12, reversal_reason = $13, total_debit = $14, total_credit = $15,
		    line_count = $16, posting_context = $17, notes = $18, metadata = $19,
		    updated_at = $20, updated_by = $21
		WHERE id = $22 AND organization_id = $23 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		audit.PostingStatus, audit.JournalEntryID, audit.ReversalJournalEntryID,
		audit.PostedAt, audit.PostedBy, audit.ErrorCode, audit.ErrorMessage,
		audit.ErrorDetails, audit.RetryCount, audit.LastRetryAt, audit.ReversedAt,
		audit.ReversedBy, audit.ReversalReason, audit.TotalDebit, audit.TotalCredit,
		audit.LineCount, audit.PostingContext, audit.Notes, audit.Metadata,
		time.Now(), audit.UpdatedBy, audit.ID, audit.OrganizationID,
	)
	return err
}

func (r *POSMappingRepository) ListFailedPostings(ctx context.Context, orgID uuid.UUID) ([]POSPostingAuditDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, source_table, source_id, source_reference, posting_status,
		       journal_entry_id, reversal_journal_entry_id, posting_date, posted_at, posted_by,
		       posting_method, error_code, error_message, error_details, retry_count, last_retry_at,
		       max_retries, reversed_at, reversed_by, reversal_reason, total_debit, total_credit,
		       line_count, currency_code, posting_context, notes, metadata, created_at, updated_at,
		       deleted_at, created_by, updated_by
		FROM pos_posting_audit
		WHERE organization_id = $1 AND posting_status = 'failed'
		  AND retry_count >= max_retries AND deleted_at IS NULL
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []POSPostingAuditDB
	for rows.Next() {
		var audit POSPostingAuditDB
		err := rows.Scan(
			&audit.ID, &audit.OrganizationID, &audit.SourceTable, &audit.SourceID, &audit.SourceReference,
			&audit.PostingStatus, &audit.JournalEntryID, &audit.ReversalJournalEntryID, &audit.PostingDate,
			&audit.PostedAt, &audit.PostedBy, &audit.PostingMethod, &audit.ErrorCode, &audit.ErrorMessage,
			&audit.ErrorDetails, &audit.RetryCount, &audit.LastRetryAt, &audit.MaxRetries, &audit.ReversedAt,
			&audit.ReversedBy, &audit.ReversalReason, &audit.TotalDebit, &audit.TotalCredit,
			&audit.LineCount, &audit.CurrencyCode, &audit.PostingContext, &audit.Notes, &audit.Metadata,
			&audit.CreatedAt, &audit.UpdatedAt, &audit.DeletedAt, &audit.CreatedBy, &audit.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}

	return audits, rows.Err()
}

func (r *POSMappingRepository) GetTaxMapping(ctx context.Context, orgID uuid.UUID, mappingID uuid.UUID) (*POSTaxMappingDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_tax_code, tax_category_code, pos_tax_rate,
		       accounting_tax_id, default_tax_account_id, default_tax_expense_account_id,
		       is_default, is_active, priority, is_inclusive, applies_to_sales, applies_to_purchases,
		       effective_from, effective_to, description, notes, metadata, created_at, updated_at,
		       deleted_at, created_by, updated_by
		FROM pos_tax_mappings
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var tm POSTaxMappingDB
	err := r.db.Pool.QueryRow(ctx, query, mappingID, orgID).Scan(
		&tm.ID, &tm.OrganizationID, &tm.POSTaxCode, &tm.TaxCategoryCode, &tm.POSTaxRate,
		&tm.AccountingTaxID, &tm.DefaultTaxAccountID, &tm.DefaultTaxExpenseAccountID,
		&tm.IsDefault, &tm.IsActive, &tm.Priority, &tm.IsInclusive, &tm.AppliesToSales,
		&tm.AppliesToPurchases, &tm.EffectiveFrom, &tm.EffectiveTo, &tm.Description,
		&tm.Notes, &tm.Metadata, &tm.CreatedAt, &tm.UpdatedAt, &tm.DeletedAt, &tm.CreatedBy, &tm.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tm, nil
}

func (r *POSMappingRepository) GetTaxMappingByCode(ctx context.Context, orgID uuid.UUID, posTaxCode string) (*POSTaxMappingDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_tax_code, tax_category_code, pos_tax_rate,
		       accounting_tax_id, default_tax_account_id, default_tax_expense_account_id,
		       is_default, is_active, priority, is_inclusive, applies_to_sales, applies_to_purchases,
		       effective_from, effective_to, description, notes, metadata, created_at, updated_at,
		       deleted_at, created_by, updated_by
		FROM pos_tax_mappings
		WHERE organization_id = $1 AND pos_tax_code = $2 AND is_active = TRUE AND deleted_at IS NULL
		ORDER BY priority DESC, created_at DESC
		LIMIT 1
	`

	var tm POSTaxMappingDB
	err := r.db.Pool.QueryRow(ctx, query, orgID, posTaxCode).Scan(
		&tm.ID, &tm.OrganizationID, &tm.POSTaxCode, &tm.TaxCategoryCode, &tm.POSTaxRate,
		&tm.AccountingTaxID, &tm.DefaultTaxAccountID, &tm.DefaultTaxExpenseAccountID,
		&tm.IsDefault, &tm.IsActive, &tm.Priority, &tm.IsInclusive, &tm.AppliesToSales,
		&tm.AppliesToPurchases, &tm.EffectiveFrom, &tm.EffectiveTo, &tm.Description,
		&tm.Notes, &tm.Metadata, &tm.CreatedAt, &tm.UpdatedAt, &tm.DeletedAt, &tm.CreatedBy, &tm.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tm, nil
}

func (r *POSMappingRepository) ListTaxMappings(ctx context.Context, orgID uuid.UUID) ([]POSTaxMappingDB, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_tax_code, tax_category_code, pos_tax_rate,
		       accounting_tax_id, default_tax_account_id, default_tax_expense_account_id,
		       is_default, is_active, priority, is_inclusive, applies_to_sales, applies_to_purchases,
		       effective_from, effective_to, description, notes, metadata, created_at, updated_at,
		       deleted_at, created_by, updated_by
		FROM pos_tax_mappings
		WHERE organization_id = $1 AND deleted_at IS NULL AND is_active = TRUE
		ORDER BY priority DESC, pos_tax_code ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []POSTaxMappingDB
	for rows.Next() {
		var tm POSTaxMappingDB
		err := rows.Scan(
			&tm.ID, &tm.OrganizationID, &tm.POSTaxCode, &tm.TaxCategoryCode, &tm.POSTaxRate,
			&tm.AccountingTaxID, &tm.DefaultTaxAccountID, &tm.DefaultTaxExpenseAccountID,
			&tm.IsDefault, &tm.IsActive, &tm.Priority, &tm.IsInclusive, &tm.AppliesToSales,
			&tm.AppliesToPurchases, &tm.EffectiveFrom, &tm.EffectiveTo, &tm.Description,
			&tm.Notes, &tm.Metadata, &tm.CreatedAt, &tm.UpdatedAt, &tm.DeletedAt, &tm.CreatedBy, &tm.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, tm)
	}

	return mappings, rows.Err()
}

func (r *POSMappingRepository) CreateTaxMapping(ctx context.Context, mapping *POSTaxMappingDB) error {
	if err := r.db.SetOrganizationContext(ctx, mapping.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO pos_tax_mappings
		(id, organization_id, pos_tax_code, tax_category_code, pos_tax_rate, accounting_tax_id,
		 default_tax_account_id, default_tax_expense_account_id, is_default, is_active, priority,
		 is_inclusive, applies_to_sales, applies_to_purchases, effective_from, effective_to,
		 description, notes, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		mapping.ID, mapping.OrganizationID, mapping.POSTaxCode, mapping.TaxCategoryCode,
		mapping.POSTaxRate, mapping.AccountingTaxID, mapping.DefaultTaxAccountID,
		mapping.DefaultTaxExpenseAccountID, mapping.IsDefault, mapping.IsActive, mapping.Priority,
		mapping.IsInclusive, mapping.AppliesToSales, mapping.AppliesToPurchases,
		mapping.EffectiveFrom, mapping.EffectiveTo, mapping.Description, mapping.Notes,
		mapping.Metadata, mapping.CreatedAt, mapping.UpdatedAt, mapping.CreatedBy,
	)
	return err
}

func (r *POSMappingRepository) UpdateTaxMapping(ctx context.Context, mapping *POSTaxMappingDB) error {
	if err := r.db.SetOrganizationContext(ctx, mapping.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_tax_mappings
		SET accounting_tax_id = $1, default_tax_account_id = $2, default_tax_expense_account_id = $3,
		    is_default = $4, is_active = $5, priority = $6, is_inclusive = $7,
		    applies_to_sales = $8, applies_to_purchases = $9, effective_from = $10, effective_to = $11,
		    description = $12, notes = $13, metadata = $14, updated_at = $15, updated_by = $16
		WHERE id = $17 AND organization_id = $18 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		mapping.AccountingTaxID, mapping.DefaultTaxAccountID, mapping.DefaultTaxExpenseAccountID,
		mapping.IsDefault, mapping.IsActive, mapping.Priority, mapping.IsInclusive,
		mapping.AppliesToSales, mapping.AppliesToPurchases, mapping.EffectiveFrom, mapping.EffectiveTo,
		mapping.Description, mapping.Notes, mapping.Metadata, time.Now(), mapping.UpdatedBy,
		mapping.ID, mapping.OrganizationID,
	)
	return err
}

// ============================================================================
// CONTEXT BINDING FOR DOMAIN LAYER
// ============================================================================

// Bind contextual repositories for the domain posting.Repository interface
// These methods map database models to domain models

func (r *PostingRuleRepository) convertToDomainRule(dbRule *PostingRuleDB, lines []PostingRuleLineDB) posting.PostingRule {
	return posting.PostingRule{
		ID:                  dbRule.ID,
		RuleCode:            dbRule.RuleCode,
		RuleName:            dbRule.RuleName,
		Event:               dbRule.Event,
		Level:               dbRule.Level,
		Priority:            dbRule.Priority,
		ConditionExpression: derefString(dbRule.ConditionExpression),
		IsActive:            dbRule.IsActive,
	}
}

func (r *PostingRuleRepository) convertToDomainRuleLine(dbLine PostingRuleLineDB) posting.PostingRuleLine {
	return posting.PostingRuleLine{
		ID:                  dbLine.ID,
		PostingRuleID:       dbLine.PostingRuleID,
		LineNo:              dbLine.LineNo,
		Side:                dbLine.Side,
		ConceptKey:          derefString(dbLine.ConceptKey),
		AccountSource:       dbLine.AccountSource,
		AmountSource:        dbLine.AmountSource,
		AmountFieldPath:     derefString(dbLine.AmountFieldPath),
		AmountExpression:    derefString(dbLine.AmountExpression),
		DescriptionTemplate: derefString(dbLine.DescriptionTemplate),
	}
}

// Helper functions
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefUUID(u *uuid.UUID) uuid.UUID {
	if u == nil {
		return uuid.Nil
	}
	return *u
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func derefFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

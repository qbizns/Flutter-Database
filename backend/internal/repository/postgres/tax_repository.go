package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/your-org/pos-backend/internal/domain/tax"
)

// TaxRepository implements tax.TaxRepository
type TaxRepository struct {
	db *sqlx.DB
}

// NewTaxRepository creates a new tax repository
func NewTaxRepository(db *sqlx.DB) tax.TaxRepository {
	return &TaxRepository{db: db}
}

// ========================
// TAX GROUP OPERATIONS
// ========================

// CreateTaxGroup creates a new tax group
func (r *TaxRepository) CreateTaxGroup(ctx context.Context, tg *tax.TaxGroup) error {
	query := `
		INSERT INTO tax_groups (
			id, organization_id, group_code, group_name, sequence,
			is_active, created_by, updated_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		tg.ID, tg.OrganizationID, tg.GroupCode, tg.GroupName, tg.Sequence,
		tg.IsActive, tg.CreatedBy, tg.UpdatedBy, tg.CreatedAt, tg.UpdatedAt,
	)

	return err
}

// GetTaxGroup retrieves a tax group
func (r *TaxRepository) GetTaxGroup(ctx context.Context, id, organizationID uuid.UUID) (*tax.TaxGroup, error) {
	query := `
		SELECT id, organization_id, group_code, group_name, sequence,
		       is_active, created_by, updated_by, created_at, updated_at, deleted_at
		FROM tax_groups
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	tg := &tax.TaxGroup{}
	err := r.db.GetContext(ctx, tg, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxGroupNotFound
		}
		return nil, err
	}

	return tg, nil
}

// GetTaxGroupByCode retrieves a tax group by code
func (r *TaxRepository) GetTaxGroupByCode(ctx context.Context, organizationID uuid.UUID, code string) (*tax.TaxGroup, error) {
	query := `
		SELECT id, organization_id, group_code, group_name, sequence,
		       is_active, created_by, updated_by, created_at, updated_at, deleted_at
		FROM tax_groups
		WHERE organization_id = $1 AND group_code = $2 AND deleted_at IS NULL
	`

	tg := &tax.TaxGroup{}
	err := r.db.GetContext(ctx, tg, query, organizationID, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxGroupNotFound
		}
		return nil, err
	}

	return tg, nil
}

// ListTaxGroups lists tax groups
func (r *TaxRepository) ListTaxGroups(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*tax.TaxGroup, error) {
	query := `
		SELECT id, organization_id, group_code, group_name, sequence,
		       is_active, created_by, updated_by, created_at, updated_at, deleted_at
		FROM tax_groups
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	query += " ORDER BY sequence ASC, group_name ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*tax.TaxGroup
	for rows.Next() {
		tg := &tax.TaxGroup{}
		if err := rows.StructScan(tg); err != nil {
			return nil, err
		}
		groups = append(groups, tg)
	}

	return groups, rows.Err()
}

// UpdateTaxGroup updates a tax group
func (r *TaxRepository) UpdateTaxGroup(ctx context.Context, tg *tax.TaxGroup) error {
	query := `
		UPDATE tax_groups
		SET group_code = $1, group_name = $2, sequence = $3,
		    is_active = $4, updated_by = $5, updated_at = $6
		WHERE id = $7 AND organization_id = $8 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		tg.GroupCode, tg.GroupName, tg.Sequence,
		tg.IsActive, tg.UpdatedBy, tg.UpdatedAt,
		tg.ID, tg.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxGroupNotFound
	}

	return nil
}

// DeleteTaxGroup soft deletes a tax group
func (r *TaxRepository) DeleteTaxGroup(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE tax_groups
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxGroupNotFound
	}

	return nil
}

// ========================
// TAX OPERATIONS
// ========================

// CreateTax creates a new tax
func (r *TaxRepository) CreateTax(ctx context.Context, t *tax.Tax) error {
	query := `
		INSERT INTO taxes (
			id, organization_id, tax_group_id, tax_code, tax_name,
			tax_rate, tax_scope, is_price_inclusive, tax_account_id,
			tax_refund_account_id, is_active, description, created_by, updated_by,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.OrganizationID, t.TaxGroupID, t.TaxCode, t.TaxName,
		t.TaxRate, t.TaxScope, t.IsPriceInclusive, t.TaxAccountID,
		t.TaxRefundAccountID, t.IsActive, t.Description, t.CreatedBy, t.UpdatedBy,
		t.CreatedAt, t.UpdatedAt,
	)

	return err
}

// GetTax retrieves a tax
func (r *TaxRepository) GetTax(ctx context.Context, id, organizationID uuid.UUID) (*tax.Tax, error) {
	query := `
		SELECT id, organization_id, tax_group_id, tax_code, tax_name,
		       tax_rate, tax_scope, is_price_inclusive, tax_account_id,
		       tax_refund_account_id, is_active, description, created_by, updated_by,
		       created_at, updated_at, deleted_at
		FROM taxes
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	t := &tax.Tax{}
	err := r.db.GetContext(ctx, t, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxNotFound
		}
		return nil, err
	}

	return t, nil
}

// GetTaxByCode retrieves a tax by code
func (r *TaxRepository) GetTaxByCode(ctx context.Context, organizationID uuid.UUID, code string) (*tax.Tax, error) {
	query := `
		SELECT id, organization_id, tax_group_id, tax_code, tax_name,
		       tax_rate, tax_scope, is_price_inclusive, tax_account_id,
		       tax_refund_account_id, is_active, description, created_by, updated_by,
		       created_at, updated_at, deleted_at
		FROM taxes
		WHERE organization_id = $1 AND tax_code = $2 AND deleted_at IS NULL
	`

	t := &tax.Tax{}
	err := r.db.GetContext(ctx, t, query, organizationID, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxNotFound
		}
		return nil, err
	}

	return t, nil
}

// ListTaxes lists taxes
func (r *TaxRepository) ListTaxes(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*tax.Tax, error) {
	query := `
		SELECT id, organization_id, tax_group_id, tax_code, tax_name,
		       tax_rate, tax_scope, is_price_inclusive, tax_account_id,
		       tax_refund_account_id, is_active, description, created_by, updated_by,
		       created_at, updated_at, deleted_at
		FROM taxes
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if taxScope, ok := filter["tax_scope"].(string); ok {
		query += fmt.Sprintf(" AND tax_scope = $%d", argIndex)
		args = append(args, taxScope)
		argIndex++
	}

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	if groupID, ok := filter["tax_group_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND tax_group_id = $%d", argIndex)
		args = append(args, groupID)
		argIndex++
	}

	query += " ORDER BY tax_code ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxes []*tax.Tax
	for rows.Next() {
		t := &tax.Tax{}
		if err := rows.StructScan(t); err != nil {
			return nil, err
		}
		taxes = append(taxes, t)
	}

	return taxes, rows.Err()
}

// UpdateTax updates a tax
func (r *TaxRepository) UpdateTax(ctx context.Context, t *tax.Tax) error {
	query := `
		UPDATE taxes
		SET tax_code = $1, tax_name = $2, tax_rate = $3,
		    tax_scope = $4, is_price_inclusive = $5, tax_account_id = $6,
		    tax_refund_account_id = $7, tax_group_id = $8, is_active = $9,
		    description = $10, updated_by = $11, updated_at = $12
		WHERE id = $13 AND organization_id = $14 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		t.TaxCode, t.TaxName, t.TaxRate,
		t.TaxScope, t.IsPriceInclusive, t.TaxAccountID,
		t.TaxRefundAccountID, t.TaxGroupID, t.IsActive,
		t.Description, t.UpdatedBy, t.UpdatedAt,
		t.ID, t.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxNotFound
	}

	return nil
}

// DeleteTax soft deletes a tax
func (r *TaxRepository) DeleteTax(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE taxes
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxNotFound
	}

	return nil
}

// ========================
// FISCAL POSITION OPERATIONS
// ========================

// CreateFiscalPosition creates a new fiscal position
func (r *TaxRepository) CreateFiscalPosition(ctx context.Context, fp *tax.FiscalPosition) error {
	query := `
		INSERT INTO fiscal_positions (
			id, organization_id, position_code, position_name, auto_apply,
			country_id, state_province, zip_postal_code_range, is_active,
			notes, created_by, updated_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		fp.ID, fp.OrganizationID, fp.PositionCode, fp.PositionName, fp.AutoApply,
		fp.CountryID, fp.StateProvince, fp.ZipPostalCodeRange, fp.IsActive,
		fp.Notes, fp.CreatedBy, fp.UpdatedBy, fp.CreatedAt, fp.UpdatedAt,
	)

	return err
}

// GetFiscalPosition retrieves a fiscal position
func (r *TaxRepository) GetFiscalPosition(ctx context.Context, id, organizationID uuid.UUID) (*tax.FiscalPosition, error) {
	query := `
		SELECT id, organization_id, position_code, position_name, auto_apply,
		       country_id, state_province, zip_postal_code_range, is_active,
		       notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM fiscal_positions
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	fp := &tax.FiscalPosition{}
	err := r.db.GetContext(ctx, fp, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrFiscalPositionNotFound
		}
		return nil, err
	}

	return fp, nil
}

// GetFiscalPositionByCode retrieves a fiscal position by code
func (r *TaxRepository) GetFiscalPositionByCode(ctx context.Context, organizationID uuid.UUID, code string) (*tax.FiscalPosition, error) {
	query := `
		SELECT id, organization_id, position_code, position_name, auto_apply,
		       country_id, state_province, zip_postal_code_range, is_active,
		       notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM fiscal_positions
		WHERE organization_id = $1 AND position_code = $2 AND deleted_at IS NULL
	`

	fp := &tax.FiscalPosition{}
	err := r.db.GetContext(ctx, fp, query, organizationID, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrFiscalPositionNotFound
		}
		return nil, err
	}

	return fp, nil
}

// ListFiscalPositions lists fiscal positions
func (r *TaxRepository) ListFiscalPositions(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*tax.FiscalPosition, error) {
	query := `
		SELECT id, organization_id, position_code, position_name, auto_apply,
		       country_id, state_province, zip_postal_code_range, is_active,
		       notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM fiscal_positions
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if countryID, ok := filter["country_id"].(string); ok {
		query += fmt.Sprintf(" AND country_id = $%d", argIndex)
		args = append(args, countryID)
		argIndex++
	}

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	query += " ORDER BY position_code ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []*tax.FiscalPosition
	for rows.Next() {
		fp := &tax.FiscalPosition{}
		if err := rows.StructScan(fp); err != nil {
			return nil, err
		}
		positions = append(positions, fp)
	}

	return positions, rows.Err()
}

// UpdateFiscalPosition updates a fiscal position
func (r *TaxRepository) UpdateFiscalPosition(ctx context.Context, fp *tax.FiscalPosition) error {
	query := `
		UPDATE fiscal_positions
		SET position_code = $1, position_name = $2, auto_apply = $3,
		    country_id = $4, state_province = $5, zip_postal_code_range = $6,
		    is_active = $7, notes = $8, updated_by = $9, updated_at = $10
		WHERE id = $11 AND organization_id = $12 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		fp.PositionCode, fp.PositionName, fp.AutoApply,
		fp.CountryID, fp.StateProvince, fp.ZipPostalCodeRange,
		fp.IsActive, fp.Notes, fp.UpdatedBy, fp.UpdatedAt,
		fp.ID, fp.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrFiscalPositionNotFound
	}

	return nil
}

// DeleteFiscalPosition soft deletes a fiscal position
func (r *TaxRepository) DeleteFiscalPosition(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE fiscal_positions
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrFiscalPositionNotFound
	}

	return nil
}

// GetApplicableFiscalPositions gets fiscal positions matching location criteria
func (r *TaxRepository) GetApplicableFiscalPositions(ctx context.Context, organizationID uuid.UUID, countryCode *string, stateProvince *string) ([]*tax.FiscalPosition, error) {
	query := `
		SELECT id, organization_id, position_code, position_name, auto_apply,
		       country_id, state_province, zip_postal_code_range, is_active,
		       notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM fiscal_positions
		WHERE organization_id = $1 AND is_active = true AND auto_apply = true
		  AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if countryCode != nil {
		query += fmt.Sprintf(" AND (country_id IS NULL OR country_id = $%d)", argIndex)
		args = append(args, *countryCode)
		argIndex++
	} else {
		query += " AND country_id IS NULL"
	}

	if stateProvince != nil && *stateProvince != "" {
		query += fmt.Sprintf(" AND (state_province IS NULL OR state_province = $%d)", argIndex)
		args = append(args, *stateProvince)
	}

	query += " ORDER BY country_id DESC, state_province DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []*tax.FiscalPosition
	for rows.Next() {
		fp := &tax.FiscalPosition{}
		if err := rows.StructScan(fp); err != nil {
			return nil, err
		}
		positions = append(positions, fp)
	}

	return positions, rows.Err()
}

// ========================
// FISCAL POSITION TAX MAPPING OPERATIONS
// ========================

// CreateFiscalPositionTaxMapping creates a tax mapping
func (r *TaxRepository) CreateFiscalPositionTaxMapping(ctx context.Context, ftm *tax.FiscalPositionTaxMapping) error {
	query := `
		INSERT INTO fiscal_position_tax_mappings (
			id, fiscal_position_id, source_tax_id, destination_tax_id,
			created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		ftm.ID, ftm.FiscalPositionID, ftm.SourceTaxID, ftm.DestinationTaxID,
		ftm.CreatedBy, ftm.CreatedAt,
	)

	return err
}

// GetFiscalPositionTaxMapping retrieves a tax mapping
func (r *TaxRepository) GetFiscalPositionTaxMapping(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) (*tax.FiscalPositionTaxMapping, error) {
	query := `
		SELECT id, fiscal_position_id, source_tax_id, destination_tax_id,
		       created_by, created_at, deleted_at
		FROM fiscal_position_tax_mappings
		WHERE fiscal_position_id = $1 AND source_tax_id = $2 AND deleted_at IS NULL
	`

	ftm := &tax.FiscalPositionTaxMapping{}
	err := r.db.GetContext(ctx, ftm, query, fiscalPositionID, sourceTaxID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxMappingNotFound
		}
		return nil, err
	}

	return ftm, nil
}

// ListFiscalPositionTaxMappings lists tax mappings
func (r *TaxRepository) ListFiscalPositionTaxMappings(ctx context.Context, fiscalPositionID uuid.UUID) ([]*tax.FiscalPositionTaxMapping, error) {
	query := `
		SELECT id, fiscal_position_id, source_tax_id, destination_tax_id,
		       created_by, created_at, deleted_at
		FROM fiscal_position_tax_mappings
		WHERE fiscal_position_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryxContext(ctx, query, fiscalPositionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []*tax.FiscalPositionTaxMapping
	for rows.Next() {
		ftm := &tax.FiscalPositionTaxMapping{}
		if err := rows.StructScan(ftm); err != nil {
			return nil, err
		}
		mappings = append(mappings, ftm)
	}

	return mappings, rows.Err()
}

// UpdateFiscalPositionTaxMapping updates a tax mapping
func (r *TaxRepository) UpdateFiscalPositionTaxMapping(ctx context.Context, ftm *tax.FiscalPositionTaxMapping) error {
	query := `
		UPDATE fiscal_position_tax_mappings
		SET destination_tax_id = $1
		WHERE fiscal_position_id = $2 AND source_tax_id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		ftm.DestinationTaxID, ftm.FiscalPositionID, ftm.SourceTaxID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxMappingNotFound
	}

	return nil
}

// DeleteFiscalPositionTaxMapping soft deletes a tax mapping
func (r *TaxRepository) DeleteFiscalPositionTaxMapping(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) error {
	query := `
		UPDATE fiscal_position_tax_mappings
		SET deleted_at = NOW()
		WHERE fiscal_position_id = $1 AND source_tax_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, fiscalPositionID, sourceTaxID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxMappingNotFound
	}

	return nil
}

// GetMappedTax retrieves the mapped destination tax
func (r *TaxRepository) GetMappedTax(ctx context.Context, fiscalPositionID, sourceTaxID uuid.UUID) (*uuid.UUID, error) {
	query := `
		SELECT destination_tax_id
		FROM fiscal_position_tax_mappings
		WHERE fiscal_position_id = $1 AND source_tax_id = $2 AND deleted_at IS NULL
	`

	var destinationTaxID *uuid.UUID
	err := r.db.GetContext(ctx, &destinationTaxID, query, fiscalPositionID, sourceTaxID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return destinationTaxID, nil
}

// ========================
// TAX REPORT DEFINITION OPERATIONS
// ========================

// CreateTaxReportDefinition creates a new tax report definition
func (r *TaxRepository) CreateTaxReportDefinition(ctx context.Context, trd *tax.TaxReportDefinition) error {
	query := `
		INSERT INTO tax_report_definitions (
			id, organization_id, localization_package_id, report_code, report_name,
			jurisdiction, authority, report_frequency, version, effective_from,
			effective_to, is_active, description, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		trd.ID, trd.OrganizationID, trd.LocalizationPackageID, trd.ReportCode, trd.ReportName,
		trd.Jurisdiction, trd.Authority, trd.ReportFrequency, trd.Version, trd.EffectiveFrom,
		trd.EffectiveTo, trd.IsActive, trd.Description, trd.CreatedBy, trd.CreatedAt, trd.UpdatedAt,
	)

	return err
}

// GetTaxReportDefinition retrieves a tax report definition
func (r *TaxRepository) GetTaxReportDefinition(ctx context.Context, id uuid.UUID) (*tax.TaxReportDefinition, error) {
	query := `
		SELECT id, organization_id, localization_package_id, report_code, report_name,
		       jurisdiction, authority, report_frequency, version, effective_from,
		       effective_to, is_active, description, created_by, created_at, updated_at, deleted_at
		FROM tax_report_definitions
		WHERE id = $1 AND deleted_at IS NULL
	`

	trd := &tax.TaxReportDefinition{}
	err := r.db.GetContext(ctx, trd, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxReportDefinitionNotFound
		}
		return nil, err
	}

	return trd, nil
}

// GetTaxReportDefinitionByCode retrieves a tax report definition by code
func (r *TaxRepository) GetTaxReportDefinitionByCode(ctx context.Context, reportCode string, organizationID *uuid.UUID) (*tax.TaxReportDefinition, error) {
	query := `
		SELECT id, organization_id, localization_package_id, report_code, report_name,
		       jurisdiction, authority, report_frequency, version, effective_from,
		       effective_to, is_active, description, created_by, created_at, updated_at, deleted_at
		FROM tax_report_definitions
		WHERE report_code = $1 AND (organization_id IS NULL OR organization_id = $2)
		  AND deleted_at IS NULL
	`

	trd := &tax.TaxReportDefinition{}
	err := r.db.GetContext(ctx, trd, query, reportCode, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxReportDefinitionNotFound
		}
		return nil, err
	}

	return trd, nil
}

// ListTaxReportDefinitions lists tax report definitions
func (r *TaxRepository) ListTaxReportDefinitions(ctx context.Context, organizationID *uuid.UUID, filter map[string]interface{}) ([]*tax.TaxReportDefinition, error) {
	query := `
		SELECT id, organization_id, localization_package_id, report_code, report_name,
		       jurisdiction, authority, report_frequency, version, effective_from,
		       effective_to, is_active, description, created_by, created_at, updated_at, deleted_at
		FROM tax_report_definitions
		WHERE (organization_id IS NULL OR organization_id = $1) AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	if reportFrequency, ok := filter["report_frequency"].(string); ok {
		query += fmt.Sprintf(" AND report_frequency = $%d", argIndex)
		args = append(args, reportFrequency)
		argIndex++
	}

	query += " ORDER BY report_code ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var definitions []*tax.TaxReportDefinition
	for rows.Next() {
		trd := &tax.TaxReportDefinition{}
		if err := rows.StructScan(trd); err != nil {
			return nil, err
		}
		definitions = append(definitions, trd)
	}

	return definitions, rows.Err()
}

// UpdateTaxReportDefinition updates a tax report definition
func (r *TaxRepository) UpdateTaxReportDefinition(ctx context.Context, trd *tax.TaxReportDefinition) error {
	query := `
		UPDATE tax_report_definitions
		SET report_code = $1, report_name = $2, jurisdiction = $3,
		    authority = $4, report_frequency = $5, version = $6,
		    effective_from = $7, effective_to = $8, is_active = $9,
		    description = $10, localization_package_id = $11, updated_at = $12
		WHERE id = $13 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		trd.ReportCode, trd.ReportName, trd.Jurisdiction,
		trd.Authority, trd.ReportFrequency, trd.Version,
		trd.EffectiveFrom, trd.EffectiveTo, trd.IsActive,
		trd.Description, trd.LocalizationPackageID, trd.UpdatedAt,
		trd.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxReportDefinitionNotFound
	}

	return nil
}

// DeleteTaxReportDefinition soft deletes a tax report definition
func (r *TaxRepository) DeleteTaxReportDefinition(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tax_report_definitions
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxReportDefinitionNotFound
	}

	return nil
}

// ========================
// TAX REPORT LINE OPERATIONS
// ========================

// CreateTaxReportLine creates a new tax report line
func (r *TaxRepository) CreateTaxReportLine(ctx context.Context, trl *tax.TaxReportLine) error {
	query := `
		INSERT INTO tax_report_lines (
			id, tax_report_definition_id, line_code, line_name, sequence,
			parent_line_id, formula_type, formula, tax_group_ids, account_ids,
			tax_ids, is_subtotal, is_total, notes, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		trl.ID, trl.TaxReportDefinitionID, trl.LineCode, trl.LineName, trl.Sequence,
		trl.ParentLineID, trl.FormulaType, trl.Formula, pq.Array(trl.TaxGroupIDs),
		pq.Array(trl.AccountIDs), pq.Array(trl.TaxIDs), trl.IsSubtotal, trl.IsTotal,
		trl.Notes, trl.CreatedAt,
	)

	return err
}

// GetTaxReportLine retrieves a tax report line
func (r *TaxRepository) GetTaxReportLine(ctx context.Context, id uuid.UUID) (*tax.TaxReportLine, error) {
	query := `
		SELECT id, tax_report_definition_id, line_code, line_name, sequence,
		       parent_line_id, formula_type, formula, tax_group_ids, account_ids,
		       tax_ids, is_subtotal, is_total, notes, created_at, deleted_at
		FROM tax_report_lines
		WHERE id = $1 AND deleted_at IS NULL
	`

	trl := &tax.TaxReportLine{}
	err := r.db.GetContext(ctx, trl, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tax.ErrTaxReportLineNotFound
		}
		return nil, err
	}

	return trl, nil
}

// ListTaxReportLines lists tax report lines
func (r *TaxRepository) ListTaxReportLines(ctx context.Context, reportDefinitionID uuid.UUID) ([]*tax.TaxReportLine, error) {
	query := `
		SELECT id, tax_report_definition_id, line_code, line_name, sequence,
		       parent_line_id, formula_type, formula, tax_group_ids, account_ids,
		       tax_ids, is_subtotal, is_total, notes, created_at, deleted_at
		FROM tax_report_lines
		WHERE tax_report_definition_id = $1 AND deleted_at IS NULL
		ORDER BY sequence ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, reportDefinitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []*tax.TaxReportLine
	for rows.Next() {
		trl := &tax.TaxReportLine{}
		if err := rows.StructScan(trl); err != nil {
			return nil, err
		}
		lines = append(lines, trl)
	}

	return lines, rows.Err()
}

// UpdateTaxReportLine updates a tax report line
func (r *TaxRepository) UpdateTaxReportLine(ctx context.Context, trl *tax.TaxReportLine) error {
	query := `
		UPDATE tax_report_lines
		SET line_code = $1, line_name = $2, sequence = $3,
		    parent_line_id = $4, formula_type = $5, formula = $6,
		    tax_group_ids = $7, account_ids = $8, tax_ids = $9,
		    is_subtotal = $10, is_total = $11, notes = $12
		WHERE id = $13 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		trl.LineCode, trl.LineName, trl.Sequence,
		trl.ParentLineID, trl.FormulaType, trl.Formula,
		pq.Array(trl.TaxGroupIDs), pq.Array(trl.AccountIDs), pq.Array(trl.TaxIDs),
		trl.IsSubtotal, trl.IsTotal, trl.Notes,
		trl.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxReportLineNotFound
	}

	return nil
}

// DeleteTaxReportLine soft deletes a tax report line
func (r *TaxRepository) DeleteTaxReportLine(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE tax_report_lines
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return tax.ErrTaxReportLineNotFound
	}

	return nil
}

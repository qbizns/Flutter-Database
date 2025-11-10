package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/auth"
)

// UserRepository implements auth.UserRepository
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// List retrieves users with filters
func (r *UserRepository) List(ctx context.Context, filters auth.UserFilters) ([]auth.User, error) {
	query := `
		SELECT
			id, organization_id, email, password_hash, first_name, last_name,
			full_name, phone, avatar_url, status, email_verified, email_verified_at,
			last_login_at, last_login_ip, failed_login_attempts, locked_until,
			two_factor_enabled, two_factor_secret, settings, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM users
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{filters.OrganizationID}
	argCount := 1

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d OR full_name ILIKE $%d)",
			argCount, argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, string(*filters.Status))
	}

	if filters.RoleID != nil {
		query += fmt.Sprintf(` AND id IN (SELECT user_id FROM user_roles WHERE role_id = $%d)`, argCount+1)
		argCount++
		args = append(args, *filters.RoleID)
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Add pagination
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	var userList []auth.User
	for rows.Next() {
		var u auth.User
		err := rows.Scan(
			&u.ID, &u.OrganizationID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
			&u.FullName, &u.Phone, &u.AvatarURL, &u.Status, &u.EmailVerified, &u.EmailVerifiedAt,
			&u.LastLoginAt, &u.LastLoginIP, &u.FailedLoginAttempts, &u.LockedUntil,
			&u.TwoFactorEnabled, &u.TwoFactorSecret, &u.Settings, &u.Metadata,
			&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.CreatedBy, &u.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		userList = append(userList, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}

// Count counts users matching filters
func (r *UserRepository) Count(ctx context.Context, filters auth.UserFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM users WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{filters.OrganizationID}
	argCount := 1

	// Apply filters (same as List)
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d OR full_name ILIKE $%d)",
			argCount, argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, string(*filters.Status))
	}

	if filters.RoleID != nil {
		query += fmt.Sprintf(` AND id IN (SELECT user_id FROM user_roles WHERE role_id = $%d)`, argCount+1)
		args = append(args, *filters.RoleID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *auth.User) error {
	query := `
		INSERT INTO users (
			id, organization_id, email, password_hash, first_name, last_name,
			phone, avatar_url, status, email_verified, email_verified_at,
			last_login_at, last_login_ip, failed_login_attempts, locked_until,
			two_factor_enabled, two_factor_secret, settings, metadata,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
	`

	settings := user.Settings
	if settings == nil {
		settings = json.RawMessage("{}")
	}
	metadata := user.Metadata
	if metadata == nil {
		metadata = json.RawMessage("{}")
	}

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.OrganizationID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.AvatarURL, user.Status, user.EmailVerified, user.EmailVerifiedAt,
		user.LastLoginAt, user.LastLoginIP, user.FailedLoginAttempts, user.LockedUntil,
		user.TwoFactorEnabled, user.TwoFactorSecret, settings, metadata,
		user.CreatedAt, user.UpdatedAt, user.CreatedBy,
	)

	return err
}

// Get retrieves a user by ID
func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	query := `
		SELECT
			id, organization_id, email, password_hash, first_name, last_name,
			full_name, phone, avatar_url, status, email_verified, email_verified_at,
			last_login_at, last_login_ip, failed_login_attempts, locked_until,
			two_factor_enabled, two_factor_secret, settings, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var u auth.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.OrganizationID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&u.FullName, &u.Phone, &u.AvatarURL, &u.Status, &u.EmailVerified, &u.EmailVerifiedAt,
		&u.LastLoginAt, &u.LastLoginIP, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.TwoFactorEnabled, &u.TwoFactorSecret, &u.Settings, &u.Metadata,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.CreatedBy, &u.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetByEmail retrieves a user by email within an organization
func (r *UserRepository) GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*auth.User, error) {
	query := `
		SELECT
			id, organization_id, email, password_hash, first_name, last_name,
			full_name, phone, avatar_url, status, email_verified, email_verified_at,
			last_login_at, last_login_ip, failed_login_attempts, locked_until,
			two_factor_enabled, two_factor_secret, settings, metadata,
			created_at, updated_at, deleted_at, created_by, updated_by
		FROM users
		WHERE organization_id = $1 AND LOWER(email) = LOWER($2) AND deleted_at IS NULL
	`

	var u auth.User
	err := r.db.Pool.QueryRow(ctx, query, orgID, email).Scan(
		&u.ID, &u.OrganizationID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&u.FullName, &u.Phone, &u.AvatarURL, &u.Status, &u.EmailVerified, &u.EmailVerifiedAt,
		&u.LastLoginAt, &u.LastLoginIP, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.TwoFactorEnabled, &u.TwoFactorSecret, &u.Settings, &u.Metadata,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.CreatedBy, &u.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *auth.User) error {
	query := `
		UPDATE users SET
			email = $3, password_hash = $4, first_name = $5, last_name = $6,
			phone = $7, avatar_url = $8, status = $9, email_verified = $10,
			email_verified_at = $11, last_login_at = $12, last_login_ip = $13,
			failed_login_attempts = $14, locked_until = $15,
			two_factor_enabled = $16, two_factor_secret = $17, settings = $18,
			metadata = $19, updated_at = $20, updated_by = $21
		WHERE id = $1 AND deleted_at IS NULL
	`

	settings := user.Settings
	if settings == nil {
		settings = json.RawMessage("{}")
	}
	metadata := user.Metadata
	if metadata == nil {
		metadata = json.RawMessage("{}")
	}

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID,
		user.CreatedAt,
		user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.AvatarURL, user.Status, user.EmailVerified,
		user.EmailVerifiedAt, user.LastLoginAt, user.LastLoginIP,
		user.FailedLoginAttempts, user.LockedUntil,
		user.TwoFactorEnabled, user.TwoFactorSecret, settings,
		metadata, user.UpdatedAt, user.UpdatedBy,
	)

	return err
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	return err
}

// UpdateLastLogin updates user's last login timestamp and IP
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, ip string) error {
	query := `
		UPDATE users
		SET last_login_at = $2, last_login_ip = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, time.Now(), ip, time.Now())
	return err
}

// IncrementFailedLoginAttempts increments failed login attempts counter
func (r *UserRepository) IncrementFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET failed_login_attempts = failed_login_attempts + 1, updated_at = $2
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, time.Now())
	return err
}

// ResetFailedLoginAttempts resets failed login attempts counter
func (r *UserRepository) ResetFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET failed_login_attempts = 0, locked_until = NULL, updated_at = $2
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, time.Now())
	return err
}

// LockUser locks a user account until a specified time
func (r *UserRepository) LockUser(ctx context.Context, userID uuid.UUID, until time.Time) error {
	query := `
		UPDATE users
		SET status = $2, locked_until = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, auth.StatusLocked, until, time.Now())
	return err
}

// RoleRepository implements auth.RoleRepository
type RoleRepository struct {
	db *DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// List retrieves roles with filters
func (r *RoleRepository) List(ctx context.Context, filters auth.RoleFilters) ([]auth.Role, error) {
	query := `
		SELECT
			id, organization_id, name, slug, description, is_system_role, is_default,
			settings, created_at, updated_at, deleted_at
		FROM roles
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 0

	// Apply filters
	if filters.OrganizationID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *filters.OrganizationID)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsSystemRole != nil {
		argCount++
		query += fmt.Sprintf(" AND is_system_role = $%d", argCount)
		args = append(args, *filters.IsSystemRole)
	}

	if filters.IsDefault != nil {
		argCount++
		query += fmt.Sprintf(" AND is_default = $%d", argCount)
		args = append(args, *filters.IsDefault)
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Add pagination
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	var roleList []auth.Role
	for rows.Next() {
		var role auth.Role
		err := rows.Scan(
			&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
			&role.IsSystemRole, &role.IsDefault, &role.Settings, &role.CreatedAt,
			&role.UpdatedAt, &role.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		roleList = append(roleList, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roleList, nil
}

// Count counts roles matching filters
func (r *RoleRepository) Count(ctx context.Context, filters auth.RoleFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 0

	// Apply filters (same as List)
	if filters.OrganizationID != nil {
		argCount++
		query += fmt.Sprintf(" AND organization_id = $%d", argCount)
		args = append(args, *filters.OrganizationID)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.IsSystemRole != nil {
		argCount++
		query += fmt.Sprintf(" AND is_system_role = $%d", argCount)
		args = append(args, *filters.IsSystemRole)
	}

	if filters.IsDefault != nil {
		argCount++
		query += fmt.Sprintf(" AND is_default = $%d", argCount)
		args = append(args, *filters.IsDefault)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create creates a new role
func (r *RoleRepository) Create(ctx context.Context, role *auth.Role) error {
	query := `
		INSERT INTO roles (
			id, organization_id, name, slug, description, is_system_role, is_default,
			settings, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	settings := role.Settings
	if settings == nil {
		settings = json.RawMessage("{}")
	}

	_, err := r.db.Pool.Exec(ctx, query,
		role.ID, role.OrganizationID, role.Name, role.Slug, role.Description,
		role.IsSystemRole, role.IsDefault, settings, role.CreatedAt, role.UpdatedAt,
	)

	return err
}

// Get retrieves a role by ID
func (r *RoleRepository) Get(ctx context.Context, id uuid.UUID) (*auth.Role, error) {
	query := `
		SELECT
			id, organization_id, name, slug, description, is_system_role, is_default,
			settings, created_at, updated_at, deleted_at
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL
	`

	var role auth.Role
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
		&role.IsSystemRole, &role.IsDefault, &role.Settings, &role.CreatedAt,
		&role.UpdatedAt, &role.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetBySlug retrieves a role by slug
func (r *RoleRepository) GetBySlug(ctx context.Context, slug string) (*auth.Role, error) {
	query := `
		SELECT
			id, organization_id, name, slug, description, is_system_role, is_default,
			settings, created_at, updated_at, deleted_at
		FROM roles
		WHERE LOWER(slug) = LOWER($1) AND deleted_at IS NULL
	`

	var role auth.Role
	err := r.db.Pool.QueryRow(ctx, query, slug).Scan(
		&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
		&role.IsSystemRole, &role.IsDefault, &role.Settings, &role.CreatedAt,
		&role.UpdatedAt, &role.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// Update updates an existing role
func (r *RoleRepository) Update(ctx context.Context, role *auth.Role) error {
	query := `
		UPDATE roles SET
			name = $3, slug = $4, description = $5, is_default = $6,
			settings = $7, updated_at = $8
		WHERE id = $1 AND deleted_at IS NULL
	`

	settings := role.Settings
	if settings == nil {
		settings = json.RawMessage("{}")
	}

	_, err := r.db.Pool.Exec(ctx, query,
		role.ID,
		role.CreatedAt,
		role.Name, role.Slug, role.Description, role.IsDefault,
		settings, role.UpdatedAt,
	)

	return err
}

// Delete soft-deletes a role
func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE roles
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	return err
}

// GetDefaultRole retrieves the default role for an organization
func (r *RoleRepository) GetDefaultRole(ctx context.Context, orgID uuid.UUID) (*auth.Role, error) {
	query := `
		SELECT
			id, organization_id, name, slug, description, is_system_role, is_default,
			settings, created_at, updated_at, deleted_at
		FROM roles
		WHERE organization_id = $1 AND is_default = true AND deleted_at IS NULL
		LIMIT 1
	`

	var role auth.Role
	err := r.db.Pool.QueryRow(ctx, query, orgID).Scan(
		&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
		&role.IsSystemRole, &role.IsDefault, &role.Settings, &role.CreatedAt,
		&role.UpdatedAt, &role.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// PermissionRepository implements auth.PermissionRepository
type PermissionRepository struct {
	db *DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// List retrieves permissions with filters
func (r *PermissionRepository) List(ctx context.Context, filters auth.PermissionFilters) ([]auth.Permission, error) {
	query := `
		SELECT
			id, name, slug, description, resource, action, category,
			created_at, updated_at
		FROM permissions
	`

	args := []interface{}{}
	argCount := 0

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" WHERE name ILIKE $%d OR slug ILIKE $%d", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	} else {
		query += " WHERE 1=1"
	}

	if filters.Resource != "" {
		argCount++
		query += fmt.Sprintf(" AND resource = $%d", argCount)
		args = append(args, filters.Resource)
	}

	if filters.Action != "" {
		argCount++
		query += fmt.Sprintf(" AND action = $%d", argCount)
		args = append(args, filters.Action)
	}

	if filters.Category != "" {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filters.Category)
	}

	// Add ordering
	query += " ORDER BY resource, action ASC"

	// Add pagination
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	var permList []auth.Permission
	for rows.Next() {
		var p auth.Permission
		err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.Resource, &p.Action,
			&p.Category, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permList = append(permList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permList, nil
}

// Count counts permissions matching filters
func (r *PermissionRepository) Count(ctx context.Context, filters auth.PermissionFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM permissions WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	// Apply filters (same as List)
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Resource != "" {
		argCount++
		query += fmt.Sprintf(" AND resource = $%d", argCount)
		args = append(args, filters.Resource)
	}

	if filters.Action != "" {
		argCount++
		query += fmt.Sprintf(" AND action = $%d", argCount)
		args = append(args, filters.Action)
	}

	if filters.Category != "" {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filters.Category)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create creates a new permission
func (r *PermissionRepository) Create(ctx context.Context, permission *auth.Permission) error {
	query := `
		INSERT INTO permissions (
			id, name, slug, description, resource, action, category,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		permission.ID, permission.Name, permission.Slug, permission.Description,
		permission.Resource, permission.Action, permission.Category,
		permission.CreatedAt, permission.UpdatedAt,
	)

	return err
}

// Get retrieves a permission by ID
func (r *PermissionRepository) Get(ctx context.Context, id uuid.UUID) (*auth.Permission, error) {
	query := `
		SELECT
			id, name, slug, description, resource, action, category,
			created_at, updated_at
		FROM permissions
		WHERE id = $1
	`

	var p auth.Permission
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.Resource, &p.Action,
		&p.Category, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetBySlug retrieves a permission by slug
func (r *PermissionRepository) GetBySlug(ctx context.Context, slug string) (*auth.Permission, error) {
	query := `
		SELECT
			id, name, slug, description, resource, action, category,
			created_at, updated_at
		FROM permissions
		WHERE LOWER(slug) = LOWER($1)
	`

	var p auth.Permission
	err := r.db.Pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.Resource, &p.Action,
		&p.Category, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Update updates an existing permission
func (r *PermissionRepository) Update(ctx context.Context, permission *auth.Permission) error {
	query := `
		UPDATE permissions SET
			name = $3, slug = $4, description = $5, resource = $6,
			action = $7, category = $8, updated_at = $9
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query,
		permission.ID,
		permission.CreatedAt,
		permission.Name, permission.Slug, permission.Description,
		permission.Resource, permission.Action, permission.Category,
		permission.UpdatedAt,
	)

	return err
}

// Delete deletes a permission
func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM permissions WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// UserRoleRepository implements auth.UserRoleRepository
type UserRoleRepository struct {
	db *DB
}

// NewUserRoleRepository creates a new user-role repository
func NewUserRoleRepository(db *DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// Assign assigns a role to a user
func (r *UserRoleRepository) Assign(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, assignedBy *uuid.UUID) error {
	query := `
		INSERT INTO user_roles (
			id, user_id, role_id, created_at, assigned_by
		) VALUES (
			$1, $2, $3, $4, $5
		)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.Pool.Exec(ctx, query,
		uuid.New(), userID, roleID, time.Now(), assignedBy,
	)

	return err
}

// Revoke revokes a role from a user
func (r *UserRoleRepository) Revoke(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, roleID)
	return err
}

// GetUserRoles retrieves all roles for a user
func (r *UserRoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]auth.Role, error) {
	query := `
		SELECT
			r.id, r.organization_id, r.name, r.slug, r.description, r.is_system_role,
			r.is_default, r.settings, r.created_at, r.updated_at, r.deleted_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.deleted_at IS NULL
		ORDER BY r.name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []auth.Role
	for rows.Next() {
		var role auth.Role
		err := rows.Scan(
			&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
			&role.IsSystemRole, &role.IsDefault, &role.Settings, &role.CreatedAt,
			&role.UpdatedAt, &role.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRoleUsers retrieves all users with a specific role
func (r *UserRoleRepository) GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]auth.User, error) {
	query := `
		SELECT
			u.id, u.organization_id, u.email, u.password_hash, u.first_name, u.last_name,
			u.full_name, u.phone, u.avatar_url, u.status, u.email_verified, u.email_verified_at,
			u.last_login_at, u.last_login_ip, u.failed_login_attempts, u.locked_until,
			u.two_factor_enabled, u.two_factor_secret, u.settings, u.metadata,
			u.created_at, u.updated_at, u.deleted_at, u.created_by, u.updated_by
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.role_id = $1 AND u.deleted_at IS NULL
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []auth.User
	for rows.Next() {
		var u auth.User
		err := rows.Scan(
			&u.ID, &u.OrganizationID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
			&u.FullName, &u.Phone, &u.AvatarURL, &u.Status, &u.EmailVerified, &u.EmailVerifiedAt,
			&u.LastLoginAt, &u.LastLoginIP, &u.FailedLoginAttempts, &u.LockedUntil,
			&u.TwoFactorEnabled, &u.TwoFactorSecret, &u.Settings, &u.Metadata,
			&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.CreatedBy, &u.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// HasRole checks if a user has a specific role
func (r *UserRoleRepository) HasRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_roles
			WHERE user_id = $1 AND role_id = $2
		)
	`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, userID, roleID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteAllUserRoles deletes all roles for a user
func (r *UserRoleRepository) DeleteAllUserRoles(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

// RolePermissionRepository implements auth.RolePermissionRepository
type RolePermissionRepository struct {
	db *DB
}

// NewRolePermissionRepository creates a new role-permission repository
func NewRolePermissionRepository(db *DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

// Assign assigns a permission to a role
func (r *RolePermissionRepository) Assign(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	query := `
		INSERT INTO role_permissions (
			id, role_id, permission_id, created_at
		) VALUES (
			$1, $2, $3, $4
		)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`

	_, err := r.db.Pool.Exec(ctx, query,
		uuid.New(), roleID, permissionID, time.Now(),
	)

	return err
}

// Revoke revokes a permission from a role
func (r *RolePermissionRepository) Revoke(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	query := `
		DELETE FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, roleID, permissionID)
	return err
}

// GetRolePermissions retrieves all permissions for a role
func (r *RolePermissionRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	query := `
		SELECT
			p.id, p.name, p.slug, p.description, p.resource, p.action, p.category,
			p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []auth.Permission
	for rows.Next() {
		var p auth.Permission
		err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.Resource, &p.Action,
			&p.Category, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPermissionRoles retrieves all roles with a specific permission
func (r *RolePermissionRepository) GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]auth.Role, error) {
	query := `
		SELECT
			r.id, r.organization_id, r.name, r.slug, r.description, r.is_system_role,
			r.is_default, r.settings, r.created_at, r.updated_at, r.deleted_at
		FROM roles r
		JOIN role_permissions rp ON r.id = rp.role_id
		WHERE rp.permission_id = $1 AND r.deleted_at IS NULL
		ORDER BY r.name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, permissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []auth.Role
	for rows.Next() {
		var role auth.Role
		err := rows.Scan(
			&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
			&role.IsSystemRole, &role.IsDefault, &role.Settings, &role.CreatedAt,
			&role.UpdatedAt, &role.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// HasPermission checks if a role has a specific permission
func (r *RolePermissionRepository) HasPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM role_permissions
			WHERE role_id = $1 AND permission_id = $2
		)
	`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, roleID, permissionID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// DeleteAllRolePermissions deletes all permissions for a role
func (r *RolePermissionRepository) DeleteAllRolePermissions(ctx context.Context, roleID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, roleID)
	return err
}

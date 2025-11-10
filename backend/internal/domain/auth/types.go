package auth

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	StatusPending  UserStatus = "pending"
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusLocked   UserStatus = "locked"
)

// User represents a system user
type User struct {
	ID                  uuid.UUID  `json:"id"`
	OrganizationID      uuid.UUID  `json:"organization_id"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"password_hash,omitempty"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	FullName            string     `json:"full_name"`
	Phone               string     `json:"phone"`
	AvatarURL           string     `json:"avatar_url"`
	Status              UserStatus `json:"status"`
	EmailVerified       bool       `json:"email_verified"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at"`
	LastLoginAt         *time.Time `json:"last_login_at"`
	LastLoginIP         string     `json:"last_login_ip"`
	FailedLoginAttempts int        `json:"failed_login_attempts"`
	LockedUntil         *time.Time `json:"locked_until"`
	TwoFactorEnabled    bool       `json:"two_factor_enabled"`
	TwoFactorSecret     string     `json:"two_factor_secret,omitempty"`
	Settings            json.RawMessage `json:"settings"`
	Metadata            json.RawMessage `json:"metadata"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
	CreatedBy           *uuid.UUID `json:"created_by"`
	UpdatedBy           *uuid.UUID `json:"updated_by"`
	Roles               []Role     `json:"roles,omitempty"`
}

// Role represents a system role
type Role struct {
	ID           uuid.UUID           `json:"id"`
	OrganizationID *uuid.UUID        `json:"organization_id"`
	Name         string              `json:"name"`
	Slug         string              `json:"slug"`
	Description  string              `json:"description"`
	IsSystemRole bool                `json:"is_system_role"`
	IsDefault    bool                `json:"is_default"`
	Settings     json.RawMessage     `json:"settings"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    *time.Time          `json:"deleted_at,omitempty"`
	Permissions  []Permission        `json:"permissions,omitempty"`
}

// Permission represents a system permission
type Permission struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Resource    string          `json:"resource"`
	Action      string          `json:"action"`
	Category    string          `json:"category"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// UserRole represents a user-role mapping
type UserRole struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	AssignedBy *uuid.UUID `json:"assigned_by"`
}

// RolePermission represents a role-permission mapping
type RolePermission struct {
	ID           uuid.UUID `json:"id"`
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserFilters represents filters for listing users
type UserFilters struct {
	Search       string
	OrganizationID uuid.UUID
	Status       *UserStatus
	RoleID       *uuid.UUID
	Page         int
	PageSize     int
}

// RoleFilters represents filters for listing roles
type RoleFilters struct {
	Search       string
	OrganizationID *uuid.UUID
	IsSystemRole *bool
	IsDefault    *bool
	Page         int
	PageSize     int
}

// PermissionFilters represents filters for listing permissions
type PermissionFilters struct {
	Search   string
	Resource string
	Action   string
	Category string
	Page     int
	PageSize int
}

// UserRepository defines the user data access interface
type UserRepository interface {
	List(ctx context.Context, filters UserFilters) ([]User, error)
	Count(ctx context.Context, filters UserFilters) (int64, error)
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID, ip string) error
	IncrementFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error
	ResetFailedLoginAttempts(ctx context.Context, userID uuid.UUID) error
	LockUser(ctx context.Context, userID uuid.UUID, until time.Time) error
}

// RoleRepository defines the role data access interface
type RoleRepository interface {
	List(ctx context.Context, filters RoleFilters) ([]Role, error)
	Count(ctx context.Context, filters RoleFilters) (int64, error)
	Create(ctx context.Context, role *Role) error
	Get(ctx context.Context, id uuid.UUID) (*Role, error)
	GetBySlug(ctx context.Context, slug string) (*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDefaultRole(ctx context.Context, orgID uuid.UUID) (*Role, error)
}

// PermissionRepository defines the permission data access interface
type PermissionRepository interface {
	List(ctx context.Context, filters PermissionFilters) ([]Permission, error)
	Count(ctx context.Context, filters PermissionFilters) (int64, error)
	Create(ctx context.Context, permission *Permission) error
	Get(ctx context.Context, id uuid.UUID) (*Permission, error)
	GetBySlug(ctx context.Context, slug string) (*Permission, error)
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserRoleRepository defines the user-role mapping data access interface
type UserRoleRepository interface {
	Assign(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, assignedBy *uuid.UUID) error
	Revoke(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)
	GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]User, error)
	HasRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (bool, error)
	DeleteAllUserRoles(ctx context.Context, userID uuid.UUID) error
}

// RolePermissionRepository defines the role-permission mapping data access interface
type RolePermissionRepository interface {
	Assign(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error
	Revoke(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error)
	GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]Role, error)
	HasPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (bool, error)
	DeleteAllRolePermissions(ctx context.Context, roleID uuid.UUID) error
}

// Scan implements the Scanner interface for UserStatus
func (s *UserStatus) Scan(value interface{}) error {
	*s = UserStatus(value.(string))
	return nil
}

// Value implements the driver Valuer interface for UserStatus
func (s UserStatus) Value() (driver.Value, error) {
	return string(s), nil
}

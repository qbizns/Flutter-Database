package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost           = 12
	MaxFailedLoginAttempts = 5
	LockoutDuration      = 15 * time.Minute
)

// UserService handles user business logic
type UserService struct {
	userRepo      UserRepository
	roleRepo      RoleRepository
	userRoleRepo  UserRoleRepository
	logger        *logging.Logger
}

// NewUserService creates a new user service
func NewUserService(
	userRepo UserRepository,
	roleRepo RoleRepository,
	userRoleRepo UserRoleRepository,
	logger *logging.Logger,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		logger:       logger,
	}
}

// List retrieves a list of users with filters
func (s *UserService) List(ctx context.Context, filters UserFilters) ([]User, error) {
	users, err := s.userRepo.List(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return users, nil
}

// Count counts users matching filters
func (s *UserService) Count(ctx context.Context, filters UserFilters) (int64, error) {
	count, err := s.userRepo.Count(ctx, filters)
	if err != nil {
		s.logger.Error("failed to count users", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, user *User, password string) error {
	// Validate user
	if err := s.validateUser(user); err != nil {
		return err
	}

	// Hash password
	if password != "" {
		hashedPassword, err := s.HashPassword(password)
		if err != nil {
			s.logger.Error("failed to hash password", zap.Error(err))
			return apperrors.InternalError(err)
		}
		user.PasswordHash = hashedPassword
	}

	// Check for duplicate email within organization
	existing, err := s.userRepo.GetByEmail(ctx, user.OrganizationID, user.Email)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate email", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("user", "email already exists in this organization")
	}

	// Set defaults
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.Status == "" {
		user.Status = StatusPending
	}
	user.FailedLoginAttempts = 0

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	// Assign default role if exists
	defaultRole, err := s.roleRepo.GetDefaultRole(ctx, user.OrganizationID)
	if err == nil && defaultRole != nil {
		if err := s.userRoleRepo.Assign(ctx, user.ID, defaultRole.ID, nil); err != nil {
			s.logger.Error("failed to assign default role", zap.Error(err))
			// Don't fail user creation if default role assignment fails
		}
	}

	return nil
}

// Get retrieves a user by ID
func (s *UserService) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if user == nil {
		return nil, apperrors.NotFound("user")
	}

	// Get user roles
	roles, err := s.userRoleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		s.logger.Error("failed to get user roles", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	user.Roles = roles

	return user, nil
}

// GetByEmail retrieves a user by email within an organization
func (s *UserService) GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*User, error) {
	user, err := s.userRepo.GetByEmail(ctx, orgID, email)
	if err != nil {
		s.logger.Error("failed to get user by email", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if user == nil {
		return nil, apperrors.NotFound("user")
	}

	// Get user roles
	roles, err := s.userRoleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		s.logger.Error("failed to get user roles", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	user.Roles = roles

	return user, nil
}

// Update updates an existing user
func (s *UserService) Update(ctx context.Context, user *User) error {
	// Validate user
	if err := s.validateUser(user); err != nil {
		return err
	}

	// Check if user exists
	existing, err := s.userRepo.Get(ctx, user.ID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("user")
	}

	// Check for duplicate email (if email changed)
	if user.Email != existing.Email {
		duplicate, err := s.userRepo.GetByEmail(ctx, user.OrganizationID, user.Email)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate email", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != user.ID {
			return apperrors.AlreadyExists("user", "email already exists in this organization")
		}
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete soft-deletes a user
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if user == nil {
		return apperrors.NotFound("user")
	}

	// Soft delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Get user
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if user == nil {
		return apperrors.NotFound("user")
	}

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return apperrors.InternalError(err)
	}

	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update user password", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// VerifyPassword checks if the provided password matches the user's stored hash
func (s *UserService) VerifyPassword(userPassword, storedHash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(userPassword)); err != nil {
		return apperrors.Unauthorized("invalid credentials")
	}
	return nil
}

// HashPassword hashes a password using bcrypt
func (s *UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// RecordLoginAttempt records a user login attempt
func (s *UserService) RecordLoginAttempt(ctx context.Context, userID uuid.UUID, ip string, success bool) error {
	if success {
		// Reset failed attempts and update last login
		if err := s.userRepo.ResetFailedLoginAttempts(ctx, userID); err != nil {
			s.logger.Error("failed to reset failed login attempts", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if err := s.userRepo.UpdateLastLogin(ctx, userID, ip); err != nil {
			s.logger.Error("failed to update last login", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
	} else {
		// Increment failed attempts
		if err := s.userRepo.IncrementFailedLoginAttempts(ctx, userID); err != nil {
			s.logger.Error("failed to increment failed login attempts", zap.Error(err))
			return apperrors.DatabaseError(err)
		}

		// Get user to check failed attempts
		user, err := s.userRepo.Get(ctx, userID)
		if err != nil {
			s.logger.Error("failed to get user", zap.Error(err))
			return apperrors.DatabaseError(err)
		}

		// Lock user if too many failed attempts
		if user != nil && user.FailedLoginAttempts >= MaxFailedLoginAttempts {
			if err := s.userRepo.LockUser(ctx, userID, time.Now().Add(LockoutDuration)); err != nil {
				s.logger.Error("failed to lock user", zap.Error(err))
				return apperrors.DatabaseError(err)
			}
		}
	}

	return nil
}

// AssignRole assigns a role to a user
func (s *UserService) AssignRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, assignedBy *uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if user == nil {
		return apperrors.NotFound("user")
	}

	// Check if role exists
	role, err := s.roleRepo.Get(ctx, roleID)
	if err != nil {
		s.logger.Error("failed to get role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if role == nil {
		return apperrors.NotFound("role")
	}

	// Assign role
	if err := s.userRoleRepo.Assign(ctx, userID, roleID, assignedBy); err != nil {
		s.logger.Error("failed to assign role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// RevokeRole revokes a role from a user
func (s *UserService) RevokeRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if user == nil {
		return apperrors.NotFound("user")
	}

	// Revoke role
	if err := s.userRoleRepo.Revoke(ctx, userID, roleID); err != nil {
		s.logger.Error("failed to revoke role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validateUser validates a user
func (s *UserService) validateUser(user *User) error {
	if user.Email == "" {
		return apperrors.ValidationFailed("email is required")
	}
	if len(user.Email) > 255 {
		return apperrors.ValidationFailed("email cannot exceed 255 characters")
	}
	if user.FirstName == "" {
		return apperrors.ValidationFailed("first name is required")
	}
	if len(user.FirstName) > 100 {
		return apperrors.ValidationFailed("first name cannot exceed 100 characters")
	}
	if user.LastName == "" {
		return apperrors.ValidationFailed("last name is required")
	}
	if len(user.LastName) > 100 {
		return apperrors.ValidationFailed("last name cannot exceed 100 characters")
	}

	return nil
}

// RoleService handles role business logic
type RoleService struct {
	roleRepo           RoleRepository
	permissionRepo     PermissionRepository
	rolePermissionRepo RolePermissionRepository
	logger             *logging.Logger
}

// NewRoleService creates a new role service
func NewRoleService(
	roleRepo RoleRepository,
	permissionRepo PermissionRepository,
	rolePermissionRepo RolePermissionRepository,
	logger *logging.Logger,
) *RoleService {
	return &RoleService{
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
		logger:             logger,
	}
}

// List retrieves a list of roles with filters
func (s *RoleService) List(ctx context.Context, filters RoleFilters) ([]Role, error) {
	roles, err := s.roleRepo.List(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list roles", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return roles, nil
}

// Count counts roles matching filters
func (s *RoleService) Count(ctx context.Context, filters RoleFilters) (int64, error) {
	count, err := s.roleRepo.Count(ctx, filters)
	if err != nil {
		s.logger.Error("failed to count roles", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// Create creates a new role
func (s *RoleService) Create(ctx context.Context, role *Role) error {
	// Validate role
	if err := s.validateRole(role); err != nil {
		return err
	}

	// Check for duplicate slug
	existing, err := s.roleRepo.GetBySlug(ctx, role.Slug)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate slug", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("role", "slug already exists")
	}

	// Set defaults
	role.ID = uuid.New()
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	// Create role
	if err := s.roleRepo.Create(ctx, role); err != nil {
		s.logger.Error("failed to create role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a role by ID
func (s *RoleService) Get(ctx context.Context, id uuid.UUID) (*Role, error) {
	role, err := s.roleRepo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get role", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if role == nil {
		return nil, apperrors.NotFound("role")
	}

	// Get role permissions
	permissions, err := s.rolePermissionRepo.GetRolePermissions(ctx, role.ID)
	if err != nil {
		s.logger.Error("failed to get role permissions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	role.Permissions = permissions

	return role, nil
}

// Update updates an existing role
func (s *RoleService) Update(ctx context.Context, role *Role) error {
	// Validate role
	if err := s.validateRole(role); err != nil {
		return err
	}

	// Check if role exists
	existing, err := s.roleRepo.Get(ctx, role.ID)
	if err != nil {
		s.logger.Error("failed to get role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("role")
	}

	// System roles cannot be updated
	if existing.IsSystemRole {
		return apperrors.Forbidden("system roles cannot be updated")
	}

	// Check for duplicate slug (if slug changed)
	if role.Slug != existing.Slug {
		duplicate, err := s.roleRepo.GetBySlug(ctx, role.Slug)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate slug", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != role.ID {
			return apperrors.AlreadyExists("role", "slug already exists")
		}
	}

	// Update timestamp
	role.UpdatedAt = time.Now()

	// Update role
	if err := s.roleRepo.Update(ctx, role); err != nil {
		s.logger.Error("failed to update role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete soft-deletes a role
func (s *RoleService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if role exists
	role, err := s.roleRepo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if role == nil {
		return apperrors.NotFound("role")
	}

	// System roles cannot be deleted
	if role.IsSystemRole {
		return apperrors.Forbidden("system roles cannot be deleted")
	}

	// Soft delete role
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// AssignPermission assigns a permission to a role
func (s *RoleService) AssignPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	// Check if role exists
	role, err := s.roleRepo.Get(ctx, roleID)
	if err != nil {
		s.logger.Error("failed to get role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if role == nil {
		return apperrors.NotFound("role")
	}

	// Check if permission exists
	permission, err := s.permissionRepo.Get(ctx, permissionID)
	if err != nil {
		s.logger.Error("failed to get permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if permission == nil {
		return apperrors.NotFound("permission")
	}

	// Assign permission
	if err := s.rolePermissionRepo.Assign(ctx, roleID, permissionID); err != nil {
		s.logger.Error("failed to assign permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// RevokePermission revokes a permission from a role
func (s *RoleService) RevokePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	// Check if role exists
	role, err := s.roleRepo.Get(ctx, roleID)
	if err != nil {
		s.logger.Error("failed to get role", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if role == nil {
		return apperrors.NotFound("role")
	}

	// Revoke permission
	if err := s.rolePermissionRepo.Revoke(ctx, roleID, permissionID); err != nil {
		s.logger.Error("failed to revoke permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validateRole validates a role
func (s *RoleService) validateRole(role *Role) error {
	if role.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	if len(role.Name) > 100 {
		return apperrors.ValidationFailed("name cannot exceed 100 characters")
	}
	if role.Slug == "" {
		return apperrors.ValidationFailed("slug is required")
	}
	if len(role.Slug) > 100 {
		return apperrors.ValidationFailed("slug cannot exceed 100 characters")
	}

	return nil
}

// PermissionService handles permission business logic
type PermissionService struct {
	repo   PermissionRepository
	logger *logging.Logger
}

// NewPermissionService creates a new permission service
func NewPermissionService(repo PermissionRepository, logger *logging.Logger) *PermissionService {
	return &PermissionService{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of permissions with filters
func (s *PermissionService) List(ctx context.Context, filters PermissionFilters) ([]Permission, error) {
	permissions, err := s.repo.List(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list permissions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return permissions, nil
}

// Count counts permissions matching filters
func (s *PermissionService) Count(ctx context.Context, filters PermissionFilters) (int64, error) {
	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		s.logger.Error("failed to count permissions", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// Create creates a new permission
func (s *PermissionService) Create(ctx context.Context, permission *Permission) error {
	// Validate permission
	if err := s.validatePermission(permission); err != nil {
		return err
	}

	// Check for duplicate slug
	existing, err := s.repo.GetBySlug(ctx, permission.Slug)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate slug", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("permission", "slug already exists")
	}

	// Set defaults
	permission.ID = uuid.New()
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()

	// Create permission
	if err := s.repo.Create(ctx, permission); err != nil {
		s.logger.Error("failed to create permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a permission by ID
func (s *PermissionService) Get(ctx context.Context, id uuid.UUID) (*Permission, error) {
	permission, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get permission", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if permission == nil {
		return nil, apperrors.NotFound("permission")
	}

	return permission, nil
}

// Update updates an existing permission
func (s *PermissionService) Update(ctx context.Context, permission *Permission) error {
	// Validate permission
	if err := s.validatePermission(permission); err != nil {
		return err
	}

	// Check if permission exists
	existing, err := s.repo.Get(ctx, permission.ID)
	if err != nil {
		s.logger.Error("failed to get permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("permission")
	}

	// Check for duplicate slug (if slug changed)
	if permission.Slug != existing.Slug {
		duplicate, err := s.repo.GetBySlug(ctx, permission.Slug)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate slug", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != permission.ID {
			return apperrors.AlreadyExists("permission", "slug already exists")
		}
	}

	// Update timestamp
	permission.UpdatedAt = time.Now()

	// Update permission
	if err := s.repo.Update(ctx, permission); err != nil {
		s.logger.Error("failed to update permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes a permission
func (s *PermissionService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if permission exists
	permission, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if permission == nil {
		return apperrors.NotFound("permission")
	}

	// Delete permission
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete permission", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validatePermission validates a permission
func (s *PermissionService) validatePermission(permission *Permission) error {
	if permission.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	if len(permission.Name) > 100 {
		return apperrors.ValidationFailed("name cannot exceed 100 characters")
	}
	if permission.Slug == "" {
		return apperrors.ValidationFailed("slug is required")
	}
	if len(permission.Slug) > 100 {
		return apperrors.ValidationFailed("slug cannot exceed 100 characters")
	}
	if permission.Resource == "" {
		return apperrors.ValidationFailed("resource is required")
	}
	if len(permission.Resource) > 100 {
		return apperrors.ValidationFailed("resource cannot exceed 100 characters")
	}
	if permission.Action == "" {
		return apperrors.ValidationFailed("action is required")
	}
	if len(permission.Action) > 50 {
		return apperrors.ValidationFailed("action cannot exceed 50 characters")
	}

	return nil
}

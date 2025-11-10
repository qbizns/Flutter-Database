package context

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestWithUserID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	ctx = WithUserID(ctx, userID)

	got, ok := GetUserID(ctx)
	if !ok {
		t.Error("GetUserID() returned false, want true")
	}

	if got != userID {
		t.Errorf("GetUserID() = %v, want %v", got, userID)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := GetUserID(ctx)
	if ok {
		t.Error("GetUserID() returned true for empty context, want false")
	}
}

func TestGetUserIDOrError_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	ctx = WithUserID(ctx, userID)

	got, err := GetUserIDOrError(ctx)
	if err != nil {
		t.Errorf("GetUserIDOrError() unexpected error = %v", err)
	}

	if got != userID {
		t.Errorf("GetUserIDOrError() = %v, want %v", got, userID)
	}
}

func TestGetUserIDOrError_NotSet(t *testing.T) {
	ctx := context.Background()

	_, err := GetUserIDOrError(ctx)
	if err == nil {
		t.Error("GetUserIDOrError() expected error for empty context, got nil")
	}

	if err != ErrUserIDNotFound {
		t.Errorf("GetUserIDOrError() error = %v, want %v", err, ErrUserIDNotFound)
	}
}

func TestMustGetUserID_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	ctx = WithUserID(ctx, userID)

	got := MustGetUserID(ctx)
	if got != userID {
		t.Errorf("MustGetUserID() = %v, want %v", got, userID)
	}
}

func TestMustGetUserID_Panic(t *testing.T) {
	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGetUserID() did not panic")
		}
	}()

	MustGetUserID(ctx)
}

func TestWithOrganizationID(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	ctx = WithOrganizationID(ctx, orgID)

	got, ok := GetOrganizationID(ctx)
	if !ok {
		t.Error("GetOrganizationID() returned false, want true")
	}

	if got != orgID {
		t.Errorf("GetOrganizationID() = %v, want %v", got, orgID)
	}
}

func TestGetOrganizationID_NotSet(t *testing.T) {
	ctx := context.Background()

	_, ok := GetOrganizationID(ctx)
	if ok {
		t.Error("GetOrganizationID() returned true for empty context, want false")
	}
}

func TestGetOrganizationIDOrError_Success(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	ctx = WithOrganizationID(ctx, orgID)

	got, err := GetOrganizationIDOrError(ctx)
	if err != nil {
		t.Errorf("GetOrganizationIDOrError() unexpected error = %v", err)
	}

	if got != orgID {
		t.Errorf("GetOrganizationIDOrError() = %v, want %v", got, orgID)
	}
}

func TestGetOrganizationIDOrError_NotSet(t *testing.T) {
	ctx := context.Background()

	_, err := GetOrganizationIDOrError(ctx)
	if err == nil {
		t.Error("GetOrganizationIDOrError() expected error for empty context, got nil")
	}

	if err != ErrOrganizationIDNotFound {
		t.Errorf("GetOrganizationIDOrError() error = %v, want %v", err, ErrOrganizationIDNotFound)
	}
}

func TestMustGetOrganizationID_Success(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	ctx = WithOrganizationID(ctx, orgID)

	got := MustGetOrganizationID(ctx)
	if got != orgID {
		t.Errorf("MustGetOrganizationID() = %v, want %v", got, orgID)
	}
}

func TestMustGetOrganizationID_Panic(t *testing.T) {
	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGetOrganizationID() did not panic")
		}
	}()

	MustGetOrganizationID(ctx)
}

func TestMultipleContextValues(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	orgID := uuid.New()

	// Set both values
	ctx = WithUserID(ctx, userID)
	ctx = WithOrganizationID(ctx, orgID)

	// Verify both are retrievable
	gotUserID, ok := GetUserID(ctx)
	if !ok {
		t.Error("GetUserID() returned false")
	}
	if gotUserID != userID {
		t.Errorf("GetUserID() = %v, want %v", gotUserID, userID)
	}

	gotOrgID, ok := GetOrganizationID(ctx)
	if !ok {
		t.Error("GetOrganizationID() returned false")
	}
	if gotOrgID != orgID {
		t.Errorf("GetOrganizationID() = %v, want %v", gotOrgID, orgID)
	}
}

func TestContextInheritance(t *testing.T) {
	// Create parent context with user ID
	parent := context.Background()
	userID := uuid.New()
	parent = WithUserID(parent, userID)

	// Create child context with organization ID
	child := WithOrganizationID(parent, uuid.New())

	// Verify child has access to both values
	gotUserID, ok := GetUserID(child)
	if !ok {
		t.Error("GetUserID() on child context returned false")
	}
	if gotUserID != userID {
		t.Errorf("GetUserID() on child = %v, want %v", gotUserID, userID)
	}

	_, ok = GetOrganizationID(child)
	if !ok {
		t.Error("GetOrganizationID() on child context returned false")
	}
}

func TestNilUUID(t *testing.T) {
	ctx := context.Background()

	// Test with nil UUID
	ctx = WithUserID(ctx, uuid.Nil)

	got, ok := GetUserID(ctx)
	if !ok {
		t.Error("GetUserID() returned false for nil UUID")
	}
	if got != uuid.Nil {
		t.Errorf("GetUserID() = %v, want %v", got, uuid.Nil)
	}
}

func TestErrorMessages(t *testing.T) {
	if ErrUserIDNotFound.Error() != "user_id not found in context" {
		t.Errorf("ErrUserIDNotFound message = %v, want 'user_id not found in context'",
			ErrUserIDNotFound.Error())
	}

	if ErrOrganizationIDNotFound.Error() != "organization_id not found in context" {
		t.Errorf("ErrOrganizationIDNotFound message = %v, want 'organization_id not found in context'",
			ErrOrganizationIDNotFound.Error())
	}
}

// Benchmark tests
func BenchmarkWithUserID(b *testing.B) {
	ctx := context.Background()
	userID := uuid.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithUserID(ctx, userID)
	}
}

func BenchmarkGetUserID(b *testing.B) {
	ctx := context.Background()
	userID := uuid.New()
	ctx = WithUserID(ctx, userID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetUserID(ctx)
	}
}

func BenchmarkGetUserIDOrError(b *testing.B) {
	ctx := context.Background()
	userID := uuid.New()
	ctx = WithUserID(ctx, userID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetUserIDOrError(ctx)
	}
}

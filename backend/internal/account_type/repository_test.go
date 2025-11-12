package account_type_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/pos-backend/internal/account_type"
	
	"github.com/your-org/pos-backend/internal/logging"
// 	"github.com/your-org/pos-backend/internal/testutil"
)

var (
	testDB     *pgxpool.Pool
	testLogger *logging.Logger
	
)

// TestMain sets up and tears down test fixtures
func TestMain(m *testing.M) {
	// Setup
	var err error
	testDB, err = // testutil.SetupTestDB()
	if err != nil {
		panic(err)
	}
	testLogger = // testutil.NewTestLogger()
	

	// Run tests
	code := m.Run()

	// Teardown
	// testutil.TeardownTestDB(testDB)

	// Exit
	// testutil.Exit(code)
}

// =============================================================================
// Repository Tests
// =============================================================================

func TestRepository_Create(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tests := []struct {
		name    string
		entity  *account_type.AccountTypes
		wantErr bool
	}{
		{
			name:    "valid entity",
			entity:  newTestAccountTypes(),
			wantErr: false,
		},
		{
			name:    "duplicate primary key",
			entity:  newTestAccountTypesWithID(uuid.New()),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			err = repo.Create(ctx, tx, tt.entity)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, tt.entity.ID)
				
				assert.NotZero(t, tt.entity.CreatedAt)
				assert.NotZero(t, tt.entity.UpdatedAt)
				
			}
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	// Create test entity
	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	entity := newTestAccountTypes()
	err = repo.Create(ctx, tx, entity)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing entity",
			id:      entity.ID,
			wantErr: false,
		},
		{
			name:    "non-existent entity",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(ctx, tx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
			}
		})
	}
}

func TestRepository_List(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entities
	for i := 0; i < 5; i++ {
		entity := newTestAccountTypes()
		err := repo.Create(ctx, tx, entity)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
		wantTotal int
	}{
		{
			name:      "first page",
			limit:     2,
			offset:    0,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "second page",
			limit:     2,
			offset:    2,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "last page",
			limit:     2,
			offset:    4,
			wantCount: 1,
			wantTotal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entities, total, err := repo.List(ctx, tx, tt.limit, tt.offset)
			assert.NoError(t, err)
			assert.Len(t, entities, tt.wantCount)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entity
	entity := newTestAccountTypes()
	err = repo.Create(ctx, tx, entity)
	require.NoError(t, err)

	// Modify entity
	
	
	entity.TypeCode = "updated_value"
	
	

	err = repo.Update(ctx, tx, entity)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(ctx, tx, entity.ID)
	require.NoError(t, err)
	
	
	assert.Equal(t, "updated_value", updated.TypeCode)
	
	
}

func TestRepository_Delete(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entity
	entity := newTestAccountTypes()
	err = repo.Create(ctx, tx, entity)
	require.NoError(t, err)

	// Delete entity
	err = repo.Delete(ctx, tx, entity.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, tx, entity.ID)
	assert.Error(t, err)
}



// =============================================================================
// Service Tests
// =============================================================================

func TestService_Create(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	service := account_type.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	req := &CreateAccountTypesRequest{
		
		TypeCode: "test_value",
		
		
		
		
		TypeName: "test_value",
		
		
		
		
		TypeCategory: "test_value",
		
		
		
		
		NormalBalance: "test_value",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	}

	
	resp, err := service.Create(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEqual(t, uuid.Nil, resp.ID)
}

func TestService_GetByID(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	service := account_type.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	// Create test entity
	req := newTestAccountTypesRequest()
	
	created, err := service.Create(ctx, req)
	
	require.NoError(t, err)

	// Get by ID
	
	resp, err := service.GetByID(ctx, created.ID)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, created.ID, resp.ID)
}

func TestService_Update(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	service := account_type.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	// Create test entity
	createReq := newTestAccountTypesRequest()
	
	created, err := service.Create(ctx, createReq)
	
	require.NoError(t, err)

	// Update entity
	
	
	updatedValue := "updated_value"
	updateReq := &UpdateAccountTypesRequest{
		TypeCode: &updatedValue,
	}
	
	

	
	resp, err := service.Update(ctx, created.ID, updateReq)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// =============================================================================
// Handler Tests
// =============================================================================

func TestHandler_Create(t *testing.T) {
	repo := account_type.NewRepository(testDB, testLogger)
	service := account_type.NewService(testDB, repo, testLogger)
	handler := account_type.NewHandler(service, testLogger)

	reqBody := &CreateAccountTypesRequest{
		
		TypeCode: "test_value",
		
		TypeName: "test_value",
		
		TypeCategory: "test_value",
		
		NormalBalance: "test_value",
		
		
		
		
		
		
		
		
		
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	
	req := httptest.NewRequest(http.MethodPost, "/api/v1/account_types", strings.NewReader(string(body)))
	
	req.Header.Set("Content-Type", "application/json")
	

	w := httptest.NewRecorder()
	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp AccountTypesResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, resp.ID)
}

// =============================================================================
// Test Fixtures and Helpers
// =============================================================================

// newTestAccountTypes creates a test AccountTypes entity
func newTestAccountTypes() *account_type.AccountTypes {
	return &account_type.AccountTypes{
		
		
		TypeCode: "test_type_code",
		
		
		
		
		
		
		
		TypeName: "test_type_name",
		
		
		
		
		
		
		
		TypeCategory: "test_type_category",
		
		
		
		
		
		
		
		NormalBalance: "test_normal_balance",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	}
}

// newTestAccountTypesWithID creates a test AccountTypes entity with specific ID
func newTestAccountTypesWithID(id uuid.UUID) *account_type.AccountTypes {
	entity := newTestAccountTypes()
	entity.ID = id
	return entity
}

// newTestAccountTypesRequest creates a test create request
func newTestAccountTypesRequest() *CreateAccountTypesRequest {
	return &CreateAccountTypesRequest{
		
		TypeCode: "test_type_code",
		
		
		
		
		TypeName: "test_type_name",
		
		
		
		
		TypeCategory: "test_type_category",
		
		
		
		
		NormalBalance: "test_normal_balance",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	}
}

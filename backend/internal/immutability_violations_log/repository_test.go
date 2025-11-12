package immutability_violations_log_test

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
	"github.com/your-org/pos-backend/internal/immutability_violations_log"
	"github.com/your-org/pos-backend/internal/dto"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/testutil"
)

var (
	testDB     *pgxpool.Pool
	testLogger *logging.Logger
	testOrgID  uuid.UUID
)

// TestMain sets up and tears down test fixtures
func TestMain(m *testing.M) {
	// Setup
	var err error
	testDB, err = testutil.SetupTestDB()
	if err != nil {
		panic(err)
	}
	testLogger = testutil.NewTestLogger()
	
	testOrgID = testutil.CreateTestOrganization(testDB)
	

	// Run tests
	code := m.Run()

	// Teardown
	testutil.TeardownTestDB(testDB)

	// Exit
	testutil.Exit(code)
}

// =============================================================================
// Repository Tests
// =============================================================================

func TestRepository_Create(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tests := []struct {
		name    string
		entity  *immutability_violations_log.ImmutabilityViolationsLog
		wantErr bool
	}{
		{
			name:    "valid entity",
			entity:  newTestImmutabilityViolationsLog(testOrgID),
			wantErr: false,
		},
		{
			name:    "duplicate primary key",
			entity:  newTestImmutabilityViolationsLogWithID(uuid.New(), testOrgID),
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
				
			}
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	ctx := context.Background()

	// Create test entity
	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	entity := newTestImmutabilityViolationsLog(testOrgID)
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
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entities
	for i := 0; i < 5; i++ {
		entity := newTestImmutabilityViolationsLog(testOrgID)
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
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entity
	entity := newTestImmutabilityViolationsLog(testOrgID)
	err = repo.Create(ctx, tx, entity)
	require.NoError(t, err)

	// Modify entity
	
	
	entity.TableName = "updated_value"
	
	

	err = repo.Update(ctx, tx, entity)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(ctx, tx, entity.ID)
	require.NoError(t, err)
	
	
	assert.Equal(t, "updated_value", updated.TableName)
	
	
}

func TestRepository_Delete(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entity
	entity := newTestImmutabilityViolationsLog(testOrgID)
	err = repo.Create(ctx, tx, entity)
	require.NoError(t, err)

	// Delete entity
	err = repo.Delete(ctx, tx, entity.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, tx, entity.ID)
	assert.Error(t, err)
}


func TestRepository_ListByOrganization(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create entities for test org
	for i := 0; i < 3; i++ {
		entity := newTestImmutabilityViolationsLog(testOrgID)
		err := repo.Create(ctx, tx, entity)
		require.NoError(t, err)
	}

	// Create entities for different org
	otherOrgID := uuid.New()
	for i := 0; i < 2; i++ {
		entity := newTestImmutabilityViolationsLog(otherOrgID)
		err := repo.Create(ctx, tx, entity)
		require.NoError(t, err)
	}

	// Query by organization
	entities, total, err := repo.ListByOrganization(ctx, tx, testOrgID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, entities, 3)
	assert.Equal(t, 3, total)

	// Verify all belong to test org
	for _, entity := range entities {
		assert.Equal(t, testOrgID, entity.OrganizationID)
	}
}


// =============================================================================
// Service Tests
// =============================================================================

func TestService_Create(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	service := immutability_violations_log.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	req := &dto.CreateImmutabilityViolationsLogRequest{
		
		TableName: "test_value",
		
		
		
		
		
		
		
		
		
		Operation: "test_value",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	}

	
	resp, err := service.Create(ctx, testOrgID, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEqual(t, uuid.Nil, resp.ID)
}

func TestService_GetByID(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	service := immutability_violations_log.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	// Create test entity
	req := newTestImmutabilityViolationsLogRequest(testOrgID)
	
	created, err := service.Create(ctx, testOrgID, req)
	
	require.NoError(t, err)

	// Get by ID
	
	resp, err := service.GetByID(ctx, testOrgID, created.ID)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, created.ID, resp.ID)
}

func TestService_Update(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	service := immutability_violations_log.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	// Create test entity
	createReq := newTestImmutabilityViolationsLogRequest(testOrgID)
	
	created, err := service.Create(ctx, testOrgID, createReq)
	
	require.NoError(t, err)

	// Update entity
	
	
	updatedValue := "updated_value"
	updateReq := &dto.UpdateImmutabilityViolationsLogRequest{
		TableName: &updatedValue,
	}
	
	

	
	resp, err := service.Update(ctx, testOrgID, created.ID, updateReq)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// =============================================================================
// Handler Tests
// =============================================================================

func TestHandler_Create(t *testing.T) {
	repo := immutability_violations_log.NewRepository(testDB, testLogger)
	service := immutability_violations_log.NewService(testDB, repo, testLogger)
	handler := immutability_violations_log.NewHandler(service, testLogger)

	reqBody := &dto.CreateImmutabilityViolationsLogRequest{
		
		TableName: "test_value",
		
		
		
		Operation: "test_value",
		
		
		
		
		
		
		
		
		
		
		
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	
	req := httptest.NewRequest(http.MethodPost, "/api/v1/organizations/"+testOrgID.String()+"/immutability_violations_log", strings.NewReader(string(body)))
	
	req.Header.Set("Content-Type", "application/json")
	
	req = testutil.WithOrgID(req, testOrgID)
	

	w := httptest.NewRecorder()
	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.ImmutabilityViolationsLogResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, resp.ID)
}

// =============================================================================
// Test Fixtures and Helpers
// =============================================================================

// newTestImmutabilityViolationsLog creates a test ImmutabilityViolationsLog entity
func newTestImmutabilityViolationsLog(orgID uuid.UUID) *immutability_violations_log.ImmutabilityViolationsLog {
	return &immutability_violations_log.ImmutabilityViolationsLog{
		
		
		OrganizationID: orgID,
		
		
		
		TableName: "test_table_name",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		Operation: "test_operation",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		AttemptedAt: time.Now(),
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	}
}

// newTestImmutabilityViolationsLogWithID creates a test ImmutabilityViolationsLog entity with specific ID
func newTestImmutabilityViolationsLogWithID(id uuid.UUID, orgID uuid.UUID) *immutability_violations_log.ImmutabilityViolationsLog {
	entity := newTestImmutabilityViolationsLog(orgID)
	entity.ID = id
	return entity
}

// newTestImmutabilityViolationsLogRequest creates a test create request
func newTestImmutabilityViolationsLogRequest(orgID uuid.UUID) *dto.CreateImmutabilityViolationsLogRequest {
	return &dto.CreateImmutabilityViolationsLogRequest{
		
		TableName: "test_table_name",
		
		
		
		
		
		
		
		
		
		Operation: "test_operation",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	}
}

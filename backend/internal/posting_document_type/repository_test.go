package posting_document_type_test

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
	"github.com/your-org/pos-backend/internal/posting_document_type"
	
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
	repo := posting_document_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tests := []struct {
		name    string
		entity  *posting_document_type.PostingDocumentTypes
		wantErr bool
	}{
		{
			name:    "valid entity",
			entity:  newTestPostingDocumentTypes(),
			wantErr: false,
		},
		{
			name:    "duplicate primary key",
			entity:  newTestPostingDocumentTypesWithID(uuid.New()),
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
	repo := posting_document_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	// Create test entity
	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	entity := newTestPostingDocumentTypes()
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
	repo := posting_document_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entities
	for i := 0; i < 5; i++ {
		entity := newTestPostingDocumentTypes()
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
	repo := posting_document_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entity
	entity := newTestPostingDocumentTypes()
	err = repo.Create(ctx, tx, entity)
	require.NoError(t, err)

	// Modify entity
	
	
	entity.Code = "updated_value"
	
	

	err = repo.Update(ctx, tx, entity)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(ctx, tx, entity.ID)
	require.NoError(t, err)
	
	
	assert.Equal(t, "updated_value", updated.Code)
	
	
}

func TestRepository_Delete(t *testing.T) {
	repo := posting_document_type.NewRepository(testDB, testLogger)
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Create test entity
	entity := newTestPostingDocumentTypes()
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
	repo := posting_document_type.NewRepository(testDB, testLogger)
	service := posting_document_type.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	req := &CreatePostingDocumentTypesRequest{
		
		Code: "test_value",
		
		
		
		
		Name: "test_value",
		
		
		
		
		
		
		
		
		
		SourceSchema: "test_value",
		
		
		
		
		SourceTable: "test_value",
		
		
		
		
		SourcePkColumn: "test_value",
		
		
		
		
		
		
		
		
		
		
		
		IsActive: true,
		
		
		
		
		IsSystem: true,
		
		
		
		
		
		
		
		
		
		
		
		
	}

	
	resp, err := service.Create(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEqual(t, uuid.Nil, resp.ID)
}

func TestService_GetByID(t *testing.T) {
	repo := posting_document_type.NewRepository(testDB, testLogger)
	service := posting_document_type.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	// Create test entity
	req := newTestPostingDocumentTypesRequest()
	
	created, err := service.Create(ctx, req)
	
	require.NoError(t, err)

	// Get by ID
	
	resp, err := service.GetByID(ctx, created.ID)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, created.ID, resp.ID)
}

func TestService_Update(t *testing.T) {
	repo := posting_document_type.NewRepository(testDB, testLogger)
	service := posting_document_type.NewService(testDB, repo, testLogger)
	ctx := context.Background()

	// Create test entity
	createReq := newTestPostingDocumentTypesRequest()
	
	created, err := service.Create(ctx, createReq)
	
	require.NoError(t, err)

	// Update entity
	
	
	updatedValue := "updated_value"
	updateReq := &UpdatePostingDocumentTypesRequest{
		Code: &updatedValue,
	}
	
	

	
	resp, err := service.Update(ctx, created.ID, updateReq)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// =============================================================================
// Handler Tests
// =============================================================================

func TestHandler_Create(t *testing.T) {
	repo := posting_document_type.NewRepository(testDB, testLogger)
	service := posting_document_type.NewService(testDB, repo, testLogger)
	handler := posting_document_type.NewHandler(service, testLogger)

	reqBody := &CreatePostingDocumentTypesRequest{
		
		Code: "test_value",
		
		Name: "test_value",
		
		
		
		SourceSchema: "test_value",
		
		SourceTable: "test_value",
		
		SourcePkColumn: "test_value",
		
		
		
		
		
		
		
		
		
		
		
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posting_document_types", strings.NewReader(string(body)))
	
	req.Header.Set("Content-Type", "application/json")
	

	w := httptest.NewRecorder()
	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp PostingDocumentTypesResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, resp.ID)
}

// =============================================================================
// Test Fixtures and Helpers
// =============================================================================

// newTestPostingDocumentTypes creates a test PostingDocumentTypes entity
func newTestPostingDocumentTypes() *posting_document_type.PostingDocumentTypes {
	return &posting_document_type.PostingDocumentTypes{
		
		
		Code: "test_code",
		
		
		
		
		
		
		
		Name: "test_name",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		SourceSchema: "test_source_schema",
		
		
		
		
		
		
		
		SourceTable: "test_source_table",
		
		
		
		
		
		
		
		SourcePkColumn: "test_source_pk_column",
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		IsActive: true,
		
		
		
		
		
		
		
		IsSystem: true,
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	}
}

// newTestPostingDocumentTypesWithID creates a test PostingDocumentTypes entity with specific ID
func newTestPostingDocumentTypesWithID(id uuid.UUID) *posting_document_type.PostingDocumentTypes {
	entity := newTestPostingDocumentTypes()
	entity.ID = id
	return entity
}

// newTestPostingDocumentTypesRequest creates a test create request
func newTestPostingDocumentTypesRequest() *CreatePostingDocumentTypesRequest {
	return &CreatePostingDocumentTypesRequest{
		
		Code: "test_code",
		
		
		
		
		Name: "test_name",
		
		
		
		
		
		
		
		
		
		SourceSchema: "test_source_schema",
		
		
		
		
		SourceTable: "test_source_table",
		
		
		
		
		SourcePkColumn: "test_source_pk_column",
		
		
		
		
		
		
		
		
		
		
		
		IsActive: true,
		
		
		
		
		IsSystem: true,
		
		
		
		
		
		
		
		
		
		
		
		
	}
}

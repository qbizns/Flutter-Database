package testhelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/postgres"
)

// TestDB wraps a database connection for testing
type TestDB struct {
	*postgres.DB
	pool *pgxpool.Pool
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Load test configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}

	// Override with test values
	cfg.Database.DBName = fmt.Sprintf("test_%s_%d", t.Name(), time.Now().Unix())
	cfg.Server.Env = "test"

	// Create logger
	logger, err := logging.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Create test database
	ctx := context.Background()
	if err := createTestDatabase(ctx, cfg); err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Connect to database
	db, err := postgres.New(cfg, logger)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	testDB := &TestDB{
		DB:   db,
		pool: db.Pool,
	}

	// Register cleanup
	t.Cleanup(func() {
		testDB.Teardown(t)
	})

	return testDB
}

// createTestDatabase creates an isolated test database
func createTestDatabase(ctx context.Context, cfg *config.Config) error {
	// Connect to default postgres database to create test database
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer pool.Close()

	// Create test database
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName))
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}

	return nil
}

// Teardown cleans up the test database
func (tdb *TestDB) Teardown(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	dbName := tdb.DB.Pool.Config().ConnConfig.Database

	// Close connection
	tdb.Close()

	// Wait a bit for connections to close
	time.Sleep(100 * time.Millisecond)

	// Connect to default database to drop test database
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		"localhost", "5433", "pos_test_user", "pos_test_password",
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Logf("warning: failed to connect for cleanup: %v", err)
		return
	}
	defer pool.Close()

	// Drop test database
	_, err = pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", dbName))
	if err != nil {
		t.Logf("warning: failed to drop test database: %v", err)
	}
}

// CreateTestOrganization creates a test organization
func (tdb *TestDB) CreateTestOrganization(ctx context.Context, name string) uuid.UUID {
	orgID := uuid.New()

	query := `
		INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`

	err := tdb.Pool.QueryRow(ctx, query, orgID, name).Scan(&orgID)
	if err != nil {
		panic(fmt.Sprintf("failed to create test organization: %v", err))
	}

	return orgID
}

// CreateTestUser creates a test user
func (tdb *TestDB) CreateTestUser(ctx context.Context, orgID uuid.UUID, email string) uuid.UUID {
	userID := uuid.New()

	query := `
		INSERT INTO users (id, organization_id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	// Use bcrypt hash of "password123"
	passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	err := tdb.Pool.QueryRow(ctx, query, userID, orgID, email, passwordHash).Scan(&userID)
	if err != nil {
		panic(fmt.Sprintf("failed to create test user: %v", err))
	}

	return userID
}

// Truncate truncates specified tables
func (tdb *TestDB) Truncate(ctx context.Context, tables ...string) error {
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		_, err := tdb.Pool.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}
	return nil
}

// BeginTestTx starts a test transaction that will be rolled back
func (tdb *TestDB) BeginTestTx(ctx context.Context, t *testing.T) context.Context {
	t.Helper()

	tx, err := tdb.Pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin test transaction: %v", err)
	}

	// Register rollback on cleanup
	t.Cleanup(func() {
		_ = tx.Rollback(context.Background())
	})

	// Return context with transaction
	return context.WithValue(ctx, "tx", tx)
}

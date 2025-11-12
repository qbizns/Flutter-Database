package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/postgres"
)

// TestRLSIsolation verifies that RLS properly isolates tenant data
func TestRLSIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg, logger := setupTest(t)
	db, err := postgres.New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()

	// Create two test organizations
	org1ID := createTestOrganization(t, db, ctx, "Org1")
	org2ID := createTestOrganization(t, db, ctx, "Org2")

	// Create test data for org1
	err = db.WithOrgContext(ctx, org1ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO products (id, organization_id, name, sku, price)
			VALUES (gen_random_uuid(), $1, 'Org1 Product', 'SKU1', 10.00)
		`, org1ID)
		return err
	})
	require.NoError(t, err)

	// Create test data for org2
	err = db.WithOrgContext(ctx, org2ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO products (id, organization_id, name, sku, price)
			VALUES (gen_random_uuid(), $2, 'Org2 Product', 'SKU2', 20.00)
		`, org2ID)
		return err
	})
	require.NoError(t, err)

	// Test: Org1 should only see their products
	var count int
	err = db.WithOrgContext(ctx, org1ID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL").Scan(&count)
	})
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Org1 should see exactly 1 product")

	// Test: Org2 should only see their products
	err = db.WithOrgContext(ctx, org2ID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL").Scan(&count)
	})
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Org2 should see exactly 1 product")

	// Cleanup
	cleanupTestOrganization(t, db, ctx, org1ID)
	cleanupTestOrganization(t, db, ctx, org2ID)
}

// TestQueryTimeout verifies that query timeouts are enforced
func TestQueryTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg, logger := setupTest(t)
	db, err := postgres.New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()

	err = db.WithTx(ctx, func(tx pgx.Tx) error {
		// Simulate a slow query (pg_sleep)
		_, err := db.ExecWithTimeout(ctx, tx, 100*time.Millisecond, "SELECT pg_sleep(1)")
		return err
	})

	// Should return a timeout error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout", "Expected timeout error")
}

// TestInvalidOrganizationContext verifies that invalid org IDs are rejected
func TestInvalidOrganizationContext(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg, logger := setupTest(t)
	db, err := postgres.New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()

	// Try to use a non-existent organization ID
	err = db.WithOrgContext(ctx, "00000000-0000-0000-0000-000000000000", func(tx pgx.Tx) error {
		return nil
	})

	// Should return an error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "organization not found", "Expected organization not found error")
}

// TestConcurrentRLSRequests verifies RLS works correctly under concurrent load
func TestConcurrentRLSRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg, logger := setupTest(t)
	db, err := postgres.New(cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()

	// Create two test organizations
	org1ID := createTestOrganization(t, db, ctx, "ConcurrentOrg1")
	org2ID := createTestOrganization(t, db, ctx, "ConcurrentOrg2")

	// Run concurrent queries
	done := make(chan bool, 20)
	errors := make(chan error, 20)

	for i := 0; i < 10; i++ {
		// Org1 queries
		go func() {
			err := db.WithOrgContext(ctx, org1ID, func(tx pgx.Tx) error {
				var count int
				return tx.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
			})
			if err != nil {
				errors <- err
			}
			done <- true
		}()

		// Org2 queries
		go func() {
			err := db.WithOrgContext(ctx, org2ID, func(tx pgx.Tx) error {
				var count int
				return tx.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
			})
			if err != nil {
				errors <- err
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Check for errors
	close(errors)
	for err := range errors {
		t.Errorf("Concurrent request failed: %v", err)
	}

	// Cleanup
	cleanupTestOrganization(t, db, ctx, org1ID)
	cleanupTestOrganization(t, db, ctx, org2ID)
}

// Helper functions

func setupTest(t *testing.T) (*config.Config, *logging.Logger) {
	cfg, err := config.Load()
	require.NoError(t, err)

	// Use test database
	cfg.Database.DBName = cfg.Database.DBName + "_test"

	logger, err := logging.NewLogger("info", "json")
	require.NoError(t, err)

	return cfg, logger
}

func createTestOrganization(t *testing.T, db *postgres.DB, ctx context.Context, name string) string {
	var orgID string
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO organizations (id, name, is_active)
		VALUES (gen_random_uuid(), $1, true)
		RETURNING id
	`, name).Scan(&orgID)
	require.NoError(t, err)
	return orgID
}

func cleanupTestOrganization(t *testing.T, db *postgres.DB, ctx context.Context, orgID string) {
	_, err := db.Pool.Exec(ctx, "DELETE FROM organizations WHERE id = $1", orgID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test organization %s: %v", orgID, err)
	}
}

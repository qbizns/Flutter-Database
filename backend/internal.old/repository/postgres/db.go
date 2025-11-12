package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/config"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// DB wraps pgxpool for database operations
type DB struct {
	Pool   *pgxpool.Pool
	logger *logging.Logger
}

// New creates a new database connection pool
func New(cfg *config.Config, logger *logging.Logger) (*DB, error) {
	// Build connection string with SSL support
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Add SSL certificates if provided
	if cfg.Database.SSLRootCert != "" {
		connString += fmt.Sprintf(" sslrootcert=%s", cfg.Database.SSLRootCert)
	}
	if cfg.Database.SSLCert != "" {
		connString += fmt.Sprintf(" sslcert=%s", cfg.Database.SSLCert)
	}
	if cfg.Database.SSLKey != "" {
		connString += fmt.Sprintf(" sslkey=%s", cfg.Database.SSLKey)
	}

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Performance tuning
	poolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.Database.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.Database.ConnMaxIdleTime
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// Create connection pool with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection with timeout
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established",
		zap.String("host", cfg.Database.Host),
		zap.String("database", cfg.Database.DBName),
		zap.String("ssl_mode", cfg.Database.SSLMode),
		zap.Int("max_conns", cfg.Database.MaxOpenConns),
		zap.Int("min_conns", cfg.Database.MaxIdleConns),
	)

	return &DB{
		Pool:   pool,
		logger: logger,
	}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	db.Pool.Close()
	db.logger.Info("database connection closed")
}

// SetOrganizationContext sets the organization_id for RLS within a transaction
// SECURITY: Uses SET LOCAL to ensure context is transaction-scoped, preventing context leakage
// between pooled connections. Must be called within a transaction.
func (db *DB) SetOrganizationContext(ctx context.Context, tx pgx.Tx, orgID string) error {
	// Validate organization exists first - prevents RLS bypass
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM organizations WHERE id = $1 AND deleted_at IS NULL)`,
		orgID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate organization: %w", err)
	}
	if !exists {
		return fmt.Errorf("organization not found or inactive: %s", orgID)
	}

	// Use SET LOCAL instead of SET (transaction-scoped, not session-scoped)
	query := `SET LOCAL app.current_organization_id = $1`
	_, err = tx.Exec(ctx, query, orgID)
	if err != nil {
		return fmt.Errorf("failed to set organization context: %w", err)
	}
	return nil
}

// Health checks database health
func (db *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return db.Pool.Ping(ctx)
}

// Stats returns pool statistics
func (db *DB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}

// BeginTx starts a new database transaction
func (db *DB) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// WithTx executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (db *DB) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			_ = tx.Rollback(ctx)
			db.logger.Error("panic in transaction, rolled back", zap.Any("panic", p))
			panic(p) // re-throw panic after rollback
		}
	}()

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			db.logger.Error("failed to rollback transaction", zap.Error(rbErr), zap.NamedError("original_error", err))
			return fmt.Errorf("rollback failed: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExecInTx executes a query within the current transaction or connection pool
type QueryExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgx.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// GetExecutor returns either a transaction or the pool based on context
func (db *DB) GetExecutor(tx pgx.Tx) QueryExecutor {
	if tx != nil {
		return tx
	}
	return db.Pool
}

// WithOrgContext executes a function within a database transaction with organization context set
// This is the RECOMMENDED way to execute queries with RLS protection
// Example:
//   err := db.WithOrgContext(ctx, orgID, func(tx pgx.Tx) error {
//       return saleRepo.Create(ctx, tx, sale)
//   })
func (db *DB) WithOrgContext(ctx context.Context, orgID string, fn func(pgx.Tx) error) error {
	return db.WithTx(ctx, func(tx pgx.Tx) error {
		// Set organization context for RLS
		if err := db.SetOrganizationContext(ctx, tx, orgID); err != nil {
			return err
		}
		// Execute the provided function
		return fn(tx)
	})
}

// ExecWithTimeout executes a query with a timeout
func (db *DB) ExecWithTimeout(ctx context.Context, tx pgx.Tx, timeout time.Duration, query string, args ...interface{}) (pgx.CommandTag, error) {
	queryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tag, err := tx.Exec(queryCtx, query, args...)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			db.logger.Warn("query timeout exceeded",
				zap.Duration("timeout", timeout),
				zap.String("query", query[:min(100, len(query))]),
			)
			return pgx.CommandTag{}, fmt.Errorf("query timeout after %v: %w", timeout, err)
		}
		return pgx.CommandTag{}, err
	}
	return tag, nil
}

// QueryWithTimeout executes a query with timeout and returns rows
func (db *DB) QueryWithTimeout(ctx context.Context, tx pgx.Tx, timeout time.Duration, query string, args ...interface{}) (pgx.Rows, error) {
	queryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rows, err := tx.Query(queryCtx, query, args...)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			db.logger.Warn("query timeout exceeded",
				zap.Duration("timeout", timeout),
				zap.String("query", query[:min(100, len(query))]),
			)
			return nil, fmt.Errorf("query timeout after %v: %w", timeout, err)
		}
		return nil, err
	}
	return rows, nil
}

// QueryRowWithTimeout executes a query with timeout and returns a single row
func (db *DB) QueryRowWithTimeout(ctx context.Context, tx pgx.Tx, timeout time.Duration, query string, args ...interface{}) pgx.Row {
	queryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return tx.QueryRow(queryCtx, query, args...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

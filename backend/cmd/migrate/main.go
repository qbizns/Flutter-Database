package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/your-org/pos-backend/internal/config"
)

func main() {
	var command string
	var steps int
	var version int

	flag.StringVar(&command, "cmd", "up", "Migration command: up, down, version, force, steps")
	flag.IntVar(&steps, "steps", 0, "Number of steps for up/down migration")
	flag.IntVar(&version, "version", 0, "Version number for force command")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Build connection string with SSL support
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Add SSL certificates if provided
	if cfg.Database.SSLRootCert != "" {
		dbURL += fmt.Sprintf("&sslrootcert=%s", cfg.Database.SSLRootCert)
	}
	if cfg.Database.SSLCert != "" {
		dbURL += fmt.Sprintf("&sslcert=%s", cfg.Database.SSLCert)
	}
	if cfg.Database.SSLKey != "" {
		dbURL += fmt.Sprintf("&sslkey=%s", cfg.Database.SSLKey)
	}

	// Create migration instances for both postgres and accounting schemas
	postgresMigrate, err := migrate.New(
		"file://../../postgres/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("Failed to create postgres migration instance:", err)
	}
	defer postgresMigrate.Close()

	accountingMigrate, err := migrate.New(
		"file://../../accounting/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("Failed to create accounting migration instance:", err)
	}
	defer accountingMigrate.Close()

	// Execute command
	switch command {
	case "up":
		fmt.Println("Running migrations UP...")

		// Run postgres migrations first
		fmt.Println("→ Running postgres schema migrations...")
		if steps > 0 {
			err = postgresMigrate.Steps(steps)
		} else {
			err = postgresMigrate.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Postgres migrations failed:", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("  No changes")
		} else {
			fmt.Println("  ✓ Postgres migrations applied")
		}

		// Run accounting migrations
		fmt.Println("→ Running accounting schema migrations...")
		if steps > 0 {
			err = accountingMigrate.Steps(steps)
		} else {
			err = accountingMigrate.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Accounting migrations failed:", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("  No changes")
		} else {
			fmt.Println("  ✓ Accounting migrations applied")
		}

		fmt.Println("✓ All migrations completed successfully")

	case "down":
		fmt.Println("Rolling back migrations...")

		// Roll back accounting first (reverse order)
		fmt.Println("→ Rolling back accounting migrations...")
		if steps > 0 {
			err = accountingMigrate.Steps(-steps)
		} else {
			err = accountingMigrate.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Accounting rollback failed:", err)
		}

		// Roll back postgres
		fmt.Println("→ Rolling back postgres migrations...")
		if steps > 0 {
			err = postgresMigrate.Steps(-steps)
		} else {
			err = postgresMigrate.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Postgres rollback failed:", err)
		}

		fmt.Println("✓ Migrations rolled back successfully")

	case "version":
		postgresVer, dirty, err := postgresMigrate.Version()
		if err != nil {
			log.Fatal("Failed to get postgres version:", err)
		}
		fmt.Printf("Postgres schema version: %d (dirty: %v)\n", postgresVer, dirty)

		accountingVer, dirty, err := accountingMigrate.Version()
		if err != nil {
			log.Fatal("Failed to get accounting version:", err)
		}
		fmt.Printf("Accounting schema version: %d (dirty: %v)\n", accountingVer, dirty)

	case "force":
		if version == 0 {
			log.Fatal("force command requires --version flag")
		}

		fmt.Printf("Forcing postgres version to %d...\n", version)
		if err := postgresMigrate.Force(version); err != nil {
			log.Fatal("Failed to force postgres version:", err)
		}

		fmt.Printf("Forcing accounting version to %d...\n", version)
		if err := accountingMigrate.Force(version); err != nil {
			log.Fatal("Failed to force accounting version:", err)
		}

		fmt.Println("✓ Version forced successfully")

	case "create":
		log.Fatal("Use 'create' command is not implemented. Create migrations manually in migrations/ directories")

	default:
		log.Fatal("Unknown command:", command)
	}
}

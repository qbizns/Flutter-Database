#!/bin/bash

# Migration runner script for POS Backend
# This script provides convenient migration management

set -e

# Configuration
MIGRATE_BIN="${HOME}/go/bin/migrate"
MIGRATIONS_POSTGRES="./db/migrations-postgres"
MIGRATIONS_ACCOUNTING="./db/migrations-accounting"
DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/pos_saas?sslmode=disable}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

check_migrate() {
    if [ ! -f "$MIGRATE_BIN" ]; then
        print_error "golang-migrate not found at $MIGRATE_BIN"
        echo "Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
    print_success "golang-migrate found"
}

check_database() {
    if ! pg_isready -d "$DATABASE_URL" > /dev/null 2>&1; then
        print_error "Cannot connect to database"
        echo "DATABASE_URL: $DATABASE_URL"
        echo "Make sure PostgreSQL is running: make docker-up"
        exit 1
    fi
    print_success "Database connection OK"
}

migrate_up() {
    echo "==================== MIGRATION UP ===================="

    print_warning "Applying PostgreSQL migrations..."
    $MIGRATE_BIN -path "$MIGRATIONS_POSTGRES" -database "$DATABASE_URL" up
    print_success "PostgreSQL migrations applied"

    print_warning "Applying Accounting migrations..."
    $MIGRATE_BIN -path "$MIGRATIONS_ACCOUNTING" -database "$DATABASE_URL" up
    print_success "Accounting migrations applied"

    echo "======================================================"
    print_success "All migrations applied successfully!"
}

migrate_down() {
    local steps=${1:-1}
    echo "==================== MIGRATION DOWN ===================="
    print_warning "Rolling back $steps step(s)..."

    $MIGRATE_BIN -path "$MIGRATIONS_ACCOUNTING" -database "$DATABASE_URL" down "$steps"
    print_success "Accounting migrations rolled back"

    $MIGRATE_BIN -path "$MIGRATIONS_POSTGRES" -database "$DATABASE_URL" down "$steps"
    print_success "PostgreSQL migrations rolled back"

    echo "======================================================="
}

migrate_status() {
    echo "==================== MIGRATION STATUS ===================="

    echo -e "\n${YELLOW}PostgreSQL migrations:${NC}"
    $MIGRATE_BIN -path "$MIGRATIONS_POSTGRES" -database "$DATABASE_URL" version 2>&1 || echo "No version set"

    echo -e "\n${YELLOW}Accounting migrations:${NC}"
    $MIGRATE_BIN -path "$MIGRATIONS_ACCOUNTING" -database "$DATABASE_URL" version 2>&1 || echo "No version set"

    echo "==========================================================="
}

migrate_force() {
    local version=$1
    if [ -z "$version" ]; then
        print_error "Version required: ./scripts/migrate.sh force <version>"
        exit 1
    fi

    echo "==================== FORCE VERSION ===================="
    print_warning "Forcing migration version to: $version"

    $MIGRATE_BIN -path "$MIGRATIONS_POSTGRES" -database "$DATABASE_URL" force "$version"
    $MIGRATE_BIN -path "$MIGRATIONS_ACCOUNTING" -database "$DATABASE_URL" force "$version"

    print_success "Version forced to: $version"
    echo "======================================================="
}

show_help() {
    cat << EOF
Migration Runner for POS Backend

Usage: ./scripts/migrate.sh <command> [options]

Commands:
  up              Apply all pending migrations
  down [N]        Rollback N migrations (default: 1)
  status          Show current migration version
  force <version> Force migration version (use when dirty)
  check           Check migration tool and database connection
  help            Show this help message

Environment Variables:
  DATABASE_URL    PostgreSQL connection string
                  Default: postgres://postgres:postgres@localhost:5432/pos_saas?sslmode=disable

Examples:
  ./scripts/migrate.sh up              # Apply all migrations
  ./scripts/migrate.sh down            # Rollback last migration
  ./scripts/migrate.sh down 3          # Rollback last 3 migrations
  ./scripts/migrate.sh status          # Check current version
  ./scripts/migrate.sh force 5         # Force version to 5

  # With custom database
  DATABASE_URL="postgres://user:pass@host:5432/db" ./scripts/migrate.sh up

EOF
}

# Main
case "${1:-help}" in
    up)
        check_migrate
        check_database
        migrate_up
        ;;
    down)
        check_migrate
        check_database
        migrate_down "${2:-1}"
        ;;
    status)
        check_migrate
        check_database
        migrate_status
        ;;
    force)
        check_migrate
        check_database
        migrate_force "$2"
        ;;
    check)
        check_migrate
        check_database
        print_success "All checks passed!"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac

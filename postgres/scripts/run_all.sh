#!/bin/bash
# ============================================================================
# Database Setup Script
# Description: Runs initialization, migrations, and seed data
# Usage: ./run_all.sh [database_name] [username]
# ============================================================================

# Default values
DB_NAME=${1:-pos_saas}
DB_USER=${2:-postgres}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}POS SAAS Database Setup${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "Database: ${GREEN}$DB_NAME${NC}"
echo -e "User: ${GREEN}$DB_USER${NC}"
echo -e "Host: ${GREEN}$DB_HOST${NC}"
echo -e "Port: ${GREEN}$DB_PORT${NC}"
echo ""

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
POSTGRES_DIR="$(dirname "$SCRIPT_DIR")"

# Function to run SQL file
run_sql_file() {
    local file=$1
    local description=$2

    echo -e "${YELLOW}Running:${NC} $description"

    if PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} Success: $description"
        return 0
    else
        echo -e "${RED}✗${NC} Failed: $description"
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
        return 1
    fi
}

# Check if database exists, create if not
echo -e "${YELLOW}Checking database...${NC}"
if PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
    echo -e "${GREEN}✓${NC} Database '$DB_NAME' exists"
else
    echo -e "${YELLOW}Creating database '$DB_NAME'...${NC}"
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;"
    echo -e "${GREEN}✓${NC} Database created"
fi

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Step 1: Database Initialization${NC}"
echo -e "${BLUE}============================================${NC}"

run_sql_file "$POSTGRES_DIR/scripts/init_database.sql" "Database initialization (extensions, functions, types)"

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Step 2: Running Migrations${NC}"
echo -e "${BLUE}============================================${NC}"

# Run migrations in order
for migration in "$POSTGRES_DIR/migrations"/*.sql; do
    if [ -f "$migration" ]; then
        filename=$(basename "$migration")
        run_sql_file "$migration" "Migration: $filename"
    fi
done

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Step 3: Seeding Data${NC}"
echo -e "${BLUE}============================================${NC}"

# Run seed data in order
for seed in "$POSTGRES_DIR/seed_data"/*.sql; do
    if [ -f "$seed" ]; then
        filename=$(basename "$seed")
        run_sql_file "$seed" "Seed data: $filename"
    fi
done

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}Setup Complete!${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "Connection string:"
echo -e "${GREEN}postgresql://$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME${NC}"
echo ""
echo -e "Test credentials:"
echo -e "  Email: ${GREEN}admin@demoretail.com${NC}"
echo -e "  Email: ${GREEN}manager@demoretail.com${NC}"
echo -e "  Email: ${GREEN}cashier1@demoretail.com${NC}"
echo ""
echo -e "${YELLOW}Note: Password hashes are for demonstration only${NC}"
echo ""

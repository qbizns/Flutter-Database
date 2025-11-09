#!/bin/bash
# ============================================================================
# Accounting Module Setup Script
# Description: Runs accounting migrations, seed data, and report views
# Usage: ./run_all.sh [database_name] [username]
# IMPORTANT: Run AFTER the main POS database is set up
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
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Accounting Module Setup${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "Database: ${GREEN}$DB_NAME${NC}"
echo -e "User: ${GREEN}$DB_USER${NC}"
echo -e "Host: ${GREEN}$DB_HOST${NC}"
echo -e "Port: ${GREEN}$DB_PORT${NC}"
echo ""

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ACCOUNTING_DIR="$(dirname "$SCRIPT_DIR")"

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

# Check if database exists
echo -e "${YELLOW}Checking database...${NC}"
if PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
    echo -e "${GREEN}✓${NC} Database '$DB_NAME' exists"
else
    echo -e "${RED}✗${NC} Database '$DB_NAME' not found"
    echo -e "${RED}ERROR: Main POS database must be created first!${NC}"
    echo -e "${YELLOW}Please run the main POS setup script first:${NC}"
    echo -e "  cd ../postgres/scripts && ./run_all.sh"
    exit 1
fi

# Check if main POS tables exist
echo -e "${YELLOW}Checking POS prerequisites...${NC}"
if PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1 FROM organizations LIMIT 1;" > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} POS base tables found"
else
    echo -e "${RED}✗${NC} POS base tables not found"
    echo -e "${RED}ERROR: Main POS system must be initialized first!${NC}"
    echo -e "${YELLOW}Please run the main POS setup script first:${NC}"
    echo -e "  cd ../postgres/scripts && ./run_all.sh"
    exit 1
fi

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Step 1: Initialize Accounting Module${NC}"
echo -e "${BLUE}============================================${NC}"

run_sql_file "$SCRIPT_DIR/init_accounting.sql" "Accounting module initialization"

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Step 2: Running Accounting Migrations${NC}"
echo -e "${BLUE}============================================${NC}"

# Run migrations in order
for migration in "$ACCOUNTING_DIR/migrations"/*.sql; do
    if [ -f "$migration" ]; then
        filename=$(basename "$migration")
        run_sql_file "$migration" "Migration: $filename"
    fi
done

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Step 3: Creating Report Views${NC}"
echo -e "${BLUE}============================================${NC}"

# Run report views
for schema_file in "$ACCOUNTING_DIR/schemas"/*.sql; do
    if [ -f "$schema_file" ]; then
        filename=$(basename "$schema_file")
        run_sql_file "$schema_file" "Schema: $filename"
    fi
done

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Step 4: Loading Seed Data${NC}"
echo -e "${BLUE}============================================${NC}"

# Run seed data in order
for seed in "$ACCOUNTING_DIR/seed_data"/*.sql; do
    if [ -f "$seed" ]; then
        filename=$(basename "$seed")
        run_sql_file "$seed" "Seed data: $filename"
    fi
done

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}Accounting Module Setup Complete!${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "${CYAN}Summary of Data Loaded:${NC}"
echo -e "  ${GREEN}✓${NC} 80+ Chart of Accounts"
echo -e "  ${GREEN}✓${NC} Fiscal Year 2024 with 12 periods"
echo -e "  ${GREEN}✓${NC} 182+ Journal Entries (full year)"
echo -e "  ${GREEN}✓${NC} 48 Vendor Bills with payments"
echo -e "  ${GREEN}✓${NC} 60 Customer Invoices with payments"
echo -e "  ${GREEN}✓${NC} 15 Fixed Assets with depreciation"
echo -e "  ${GREEN}✓${NC} 8 Financial Report Views"
echo ""
echo -e "${CYAN}Available Reports:${NC}"
echo -e "  1. Trial Balance:        ${YELLOW}SELECT * FROM accounting.view_trial_balance;${NC}"
echo -e "  2. Balance Sheet:        ${YELLOW}SELECT * FROM accounting.view_balance_sheet;${NC}"
echo -e "  3. Income Statement:     ${YELLOW}SELECT * FROM accounting.view_income_statement;${NC}"
echo -e "  4. Cash Flow:            ${YELLOW}SELECT * FROM accounting.view_cash_flow;${NC}"
echo -e "  5. Account Activity:     ${YELLOW}SELECT * FROM accounting.view_account_activity WHERE account_code = '1020';${NC}"
echo -e "  6. Aged AP:              ${YELLOW}SELECT * FROM accounting.view_aged_accounts_payable;${NC}"
echo -e "  7. Aged AR:              ${YELLOW}SELECT * FROM accounting.view_aged_accounts_receivable;${NC}"
echo -e "  8. Financial Ratios:     ${YELLOW}SELECT * FROM accounting.view_financial_ratios;${NC}"
echo ""
echo -e "${CYAN}Quick Test Query:${NC}"
echo -e "${YELLOW}psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \"SELECT * FROM accounting.view_income_statement WHERE fiscal_year = '2024' ORDER BY account_code;\"${NC}"
echo ""

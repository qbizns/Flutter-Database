#!/bin/bash

# Migration Validation Script
# Validates migration files for common issues before running them

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

MIGRATIONS_POSTGRES="../postgres/migrations"
MIGRATIONS_ACCOUNTING="../accounting/migrations"

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_info() { echo -e "  $1"; }

ERRORS=0

echo "==================== MIGRATION VALIDATION ===================="

# Check if directories exist
echo -e "\n${YELLOW}1. Checking migration directories...${NC}"
if [ -d "$MIGRATIONS_POSTGRES" ]; then
    POSTGRES_COUNT=$(find "$MIGRATIONS_POSTGRES" -name "V*.sql" | wc -l)
    print_success "PostgreSQL migrations found: $POSTGRES_COUNT files"
else
    print_error "PostgreSQL migrations directory not found"
    ERRORS=$((ERRORS + 1))
fi

if [ -d "$MIGRATIONS_ACCOUNTING" ]; then
    ACCOUNTING_COUNT=$(find "$MIGRATIONS_ACCOUNTING" -name "V*.sql" | wc -l)
    print_success "Accounting migrations found: $ACCOUNTING_COUNT files"
else
    print_error "Accounting migrations directory not found"
    ERRORS=$((ERRORS + 1))
fi

# Check file naming convention
echo -e "\n${YELLOW}2. Checking file naming conventions...${NC}"
for dir in "$MIGRATIONS_POSTGRES" "$MIGRATIONS_ACCOUNTING"; do
    if [ -d "$dir" ]; then
        invalid_files=$(find "$dir" -type f -name "*.sql" ! -name "V*_*.sql" ! -name "MIGRATION_TEMPLATE.sql")
        if [ -n "$invalid_files" ]; then
            print_error "Invalid file names found in $(basename "$dir"):"
            echo "$invalid_files" | while read -r file; do
                print_info "  - $(basename "$file")"
            done
            ERRORS=$((ERRORS + 1))
        else
            print_success "All files in $(basename "$dir") follow naming convention"
        fi
    fi
done

# Check for duplicate version numbers
echo -e "\n${YELLOW}3. Checking for duplicate version numbers...${NC}"
for dir in "$MIGRATIONS_POSTGRES" "$MIGRATIONS_ACCOUNTING"; do
    if [ -d "$dir" ]; then
        duplicates=$(find "$dir" -name "V*.sql" -type f | sed 's/.*V\([0-9]*\)_.*/\1/' | sort | uniq -d)
        if [ -n "$duplicates" ]; then
            print_error "Duplicate version numbers in $(basename "$dir"): $duplicates"
            ERRORS=$((ERRORS + 1))
        else
            print_success "No duplicate versions in $(basename "$dir")"
        fi
    fi
done

# Check file sizes (empty files)
echo -e "\n${YELLOW}4. Checking for empty migration files...${NC}"
for dir in "$MIGRATIONS_POSTGRES" "$MIGRATIONS_ACCOUNTING"; do
    if [ -d "$dir" ]; then
        empty_files=$(find "$dir" -name "V*.sql" -type f -empty)
        if [ -n "$empty_files" ]; then
            print_error "Empty migration files found in $(basename "$dir"):"
            echo "$empty_files" | while read -r file; do
                print_info "  - $(basename "$file")"
            done
            ERRORS=$((ERRORS + 1))
        else
            print_success "No empty files in $(basename "$dir")"
        fi
    fi
done

# Check for common SQL issues
echo -e "\n${YELLOW}5. Checking for common SQL issues...${NC}"
check_sql_issues() {
    local dir=$1
    local dir_name=$(basename "$dir")
    local issues=0

    if [ ! -d "$dir" ]; then
        return
    fi

    # Check for missing semicolons at end of statements
    files_without_semicolon=$(find "$dir" -name "V*.sql" -type f -exec sh -c '
        last_line=$(grep -v "^--" "$1" | grep -v "^$" | tail -1)
        if [ -n "$last_line" ] && ! echo "$last_line" | grep -q ";$"; then
            echo "$1"
        fi
    ' _ {} \;)

    if [ -n "$files_without_semicolon" ]; then
        print_warning "Files in $dir_name may be missing trailing semicolons:"
        echo "$files_without_semicolon" | while read -r file; do
            print_info "  - $(basename "$file")"
        done
        issues=$((issues + 1))
    fi

    # Check for DROP TABLE without IF EXISTS
    files_with_unsafe_drop=$(grep -l "DROP TABLE [^I]" "$dir"/V*.sql 2>/dev/null || true)
    if [ -n "$files_with_unsafe_drop" ]; then
        print_warning "Files in $dir_name contain DROP TABLE without IF EXISTS:"
        echo "$files_with_unsafe_drop" | while read -r file; do
            print_info "  - $(basename "$file")"
        done
        issues=$((issues + 1))
    fi

    if [ $issues -eq 0 ]; then
        print_success "No common SQL issues in $dir_name"
    fi
}

check_sql_issues "$MIGRATIONS_POSTGRES"
check_sql_issues "$MIGRATIONS_ACCOUNTING"

# Check migration sequence
echo -e "\n${YELLOW}6. Checking migration sequence...${NC}"
check_sequence() {
    local dir=$1
    local dir_name=$(basename "$dir")

    if [ ! -d "$dir" ]; then
        return
    fi

    versions=$(find "$dir" -name "V*.sql" -type f | sed 's/.*V\([0-9]*\)_.*/\1/' | sort -n)
    expected=1
    gaps=""

    for ver in $versions; do
        ver_num=$((10#$ver))  # Convert to decimal (remove leading zeros)
        if [ $ver_num -ne $expected ]; then
            gaps="$gaps $expected"
        fi
        expected=$((ver_num + 1))
    done

    if [ -n "$gaps" ]; then
        print_warning "Version gaps in $dir_name: missing versions$gaps"
    else
        print_success "Migration sequence is continuous in $dir_name"
    fi
}

check_sequence "$MIGRATIONS_POSTGRES"
check_sequence "$MIGRATIONS_ACCOUNTING"

# Summary
echo -e "\n============================================================"
if [ $ERRORS -eq 0 ]; then
    print_success "All validation checks passed!"
    echo -e "\nTotal migrations: $((POSTGRES_COUNT + ACCOUNTING_COUNT))"
    echo -e "  - PostgreSQL: $POSTGRES_COUNT"
    echo -e "  - Accounting: $ACCOUNTING_COUNT"
    exit 0
else
    print_error "Validation failed with $ERRORS critical error(s)"
    exit 1
fi

#!/bin/bash

#################################################################
# POS Backend Database Restore Script
# Usage: ./scripts/restore.sh [backup_file or 'latest']
#################################################################

set -euo pipefail

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKUP_DIR="$PROJECT_ROOT/backups"
BACKUP_FILE_ARG="${1:-}"

# Load environment
ENV_FILE="${ENV_FILE:-.env.production}"
if [ -f "$PROJECT_ROOT/$ENV_FILE" ]; then
    set -a
    source "$PROJECT_ROOT/$ENV_FILE"
    set +a
fi

# Database configuration
DB_HOST="${DATABASE_HOST:-postgres}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME:-pos_db}"
DB_USER="${DATABASE_USER:-postgres}"
DB_PASSWORD="${DATABASE_PASSWORD:-postgres}"

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Find backup file
find_backup_file() {
    if [ -z "$BACKUP_FILE_ARG" ]; then
        error "No backup file specified"
        echo "Usage: $0 [backup_file or 'latest']"
        exit 1
    fi

    if [ "$BACKUP_FILE_ARG" = "latest" ]; then
        BACKUP_FILE=$(find "$BACKUP_DIR" -name "pos_backup_*.sql.gz" -type f -printf '%T@ %p\n' | sort -rn | head -1 | cut -d' ' -f2-)

        if [ -z "$BACKUP_FILE" ]; then
            error "No backup files found in $BACKUP_DIR"
            exit 1
        fi

        log "Found latest backup: $(basename "$BACKUP_FILE")"
    else
        if [ -f "$BACKUP_FILE_ARG" ]; then
            BACKUP_FILE="$BACKUP_FILE_ARG"
        elif [ -f "$BACKUP_DIR/$BACKUP_FILE_ARG" ]; then
            BACKUP_FILE="$BACKUP_DIR/$BACKUP_FILE_ARG"
        else
            error "Backup file not found: $BACKUP_FILE_ARG"
            exit 1
        fi
    fi

    log "Using backup file: $(basename "$BACKUP_FILE")"
}

# Verify backup file
verify_backup() {
    log "Verifying backup file..."

    if [ ! -f "$BACKUP_FILE" ]; then
        error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi

    if ! gzip -t "$BACKUP_FILE" 2> /dev/null; then
        error "Backup file is corrupted"
        exit 1
    fi

    log "✅ Backup file verified"
}

# Check database connectivity
check_database() {
    log "Checking database connectivity..."

    if docker-compose -f docker-compose.prod.yml exec -T postgres pg_isready -U "$DB_USER" > /dev/null 2>&1; then
        log "✅ Database is accessible"
    else
        error "Database is not accessible"
        exit 1
    fi
}

# Create pre-restore backup
create_pre_restore_backup() {
    log "Creating pre-restore backup..."

    local pre_restore_backup="$BACKUP_DIR/pre_restore_$(date +%Y%m%d_%H%M%S).sql.gz"

    if docker-compose -f docker-compose.prod.yml exec -T postgres \
        pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "$pre_restore_backup"; then

        log "✅ Pre-restore backup created: $(basename "$pre_restore_backup")"
    else
        warn "Failed to create pre-restore backup"
    fi
}

# Terminate active connections
terminate_connections() {
    log "Terminating active database connections..."

    docker-compose -f docker-compose.prod.yml exec -T postgres psql -U "$DB_USER" -d postgres <<EOF
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = '$DB_NAME'
  AND pid <> pg_backend_pid();
EOF

    log "✅ Active connections terminated"
}

# Drop and recreate database
recreate_database() {
    log "Dropping and recreating database..."

    docker-compose -f docker-compose.prod.yml exec -T postgres psql -U "$DB_USER" -d postgres <<EOF
DROP DATABASE IF EXISTS "$DB_NAME";
CREATE DATABASE "$DB_NAME";
EOF

    log "✅ Database recreated"
}

# Restore backup
restore_backup() {
    log "Restoring backup..."

    if gunzip < "$BACKUP_FILE" | docker-compose -f docker-compose.prod.yml exec -T postgres \
        psql -U "$DB_USER" -d "$DB_NAME" > /dev/null; then

        log "✅ Backup restored successfully"
    else
        error "Restore failed"
        exit 1
    fi
}

# Verify restore
verify_restore() {
    log "Verifying restore..."

    # Check if database exists and is accessible
    if ! docker-compose -f docker-compose.prod.yml exec -T postgres \
        psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" > /dev/null 2>&1; then

        error "Database verification failed"
        exit 1
    fi

    # Check for key tables
    local table_count=$(docker-compose -f docker-compose.prod.yml exec -T postgres \
        psql -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public'")

    log "✅ Restore verified ($table_count tables found)"
}

# Update sequences (if needed)
update_sequences() {
    log "Updating sequences..."

    docker-compose -f docker-compose.prod.yml exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" <<'EOF'
DO $$
DECLARE
    seq_record RECORD;
    max_val bigint;
BEGIN
    FOR seq_record IN
        SELECT sequence_schema, sequence_name
        FROM information_schema.sequences
        WHERE sequence_schema = 'public'
    LOOP
        EXECUTE format('SELECT COALESCE(MAX(id), 0) FROM %I.%I',
                      seq_record.sequence_schema,
                      regexp_replace(seq_record.sequence_name, '_id_seq$', ''))
        INTO max_val;

        EXECUTE format('SELECT setval(%L, %s)',
                      seq_record.sequence_schema || '.' || seq_record.sequence_name,
                      max_val + 1);
    END LOOP;
END $$;
EOF

    log "✅ Sequences updated"
}

# Restart services
restart_services() {
    log "Restarting services..."

    docker-compose -f docker-compose.prod.yml restart backend worker

    log "✅ Services restarted"
}

# Main restore process
main() {
    log "========================================="
    log "POS Database Restore"
    log "========================================="

    # Confirmation prompt
    warn "⚠️  WARNING: This will overwrite the current database!"
    read -p "Are you sure you want to continue? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        log "Restore cancelled"
        exit 0
    fi

    find_backup_file
    verify_backup
    check_database
    create_pre_restore_backup
    terminate_connections
    recreate_database
    restore_backup
    verify_restore
    update_sequences
    restart_services

    log "========================================="
    log "✅ Restore completed successfully!"
    log "Backup file: $(basename "$BACKUP_FILE")"
    log "========================================="
}

# Run main
main

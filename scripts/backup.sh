#!/bin/bash

#################################################################
# POS Backend Database Backup Script
# Usage: ./scripts/backup.sh [backup_type]
# Types: full (default), schema-only, data-only
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
BACKUP_TYPE="${1:-full}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/pos_backup_${BACKUP_TYPE}_${TIMESTAMP}.sql.gz"
RETENTION_DAYS=30

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

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Check if database is accessible
check_database() {
    log "Checking database connectivity..."

    if docker-compose -f docker-compose.prod.yml exec -T postgres pg_isready -U "$DB_USER" > /dev/null 2>&1; then
        log "✅ Database is accessible"
    else
        error "Database is not accessible"
        exit 1
    fi
}

# Create backup
create_backup() {
    log "Creating $BACKUP_TYPE backup..."
    log "Backup file: $BACKUP_FILE"

    local pg_dump_opts=""

    case "$BACKUP_TYPE" in
        full)
            pg_dump_opts=""
            ;;
        schema-only)
            pg_dump_opts="--schema-only"
            ;;
        data-only)
            pg_dump_opts="--data-only"
            ;;
        *)
            error "Unknown backup type: $BACKUP_TYPE"
            exit 1
            ;;
    esac

    # Create backup using docker-compose
    if docker-compose -f docker-compose.prod.yml exec -T postgres \
        pg_dump -U "$DB_USER" $pg_dump_opts "$DB_NAME" | gzip > "$BACKUP_FILE"; then

        local size=$(du -h "$BACKUP_FILE" | cut -f1)
        log "✅ Backup created successfully ($size)"
    else
        error "Backup failed"
        rm -f "$BACKUP_FILE"
        exit 1
    fi
}

# Verify backup
verify_backup() {
    log "Verifying backup..."

    if [ ! -f "$BACKUP_FILE" ]; then
        error "Backup file not found"
        exit 1
    fi

    if ! gzip -t "$BACKUP_FILE" 2> /dev/null; then
        error "Backup file is corrupted"
        exit 1
    fi

    log "✅ Backup verification passed"
}

# Upload to S3 (optional)
upload_to_s3() {
    if [ -n "${S3_BACKUP_BUCKET:-}" ]; then
        log "Uploading backup to S3..."

        if command -v aws &> /dev/null; then
            if aws s3 cp "$BACKUP_FILE" "s3://$S3_BACKUP_BUCKET/backups/$(basename "$BACKUP_FILE")"; then
                log "✅ Backup uploaded to S3"
            else
                warn "Failed to upload backup to S3"
            fi
        else
            warn "AWS CLI not installed, skipping S3 upload"
        fi
    fi
}

# Cleanup old backups
cleanup_old_backups() {
    log "Cleaning up old backups (retention: $RETENTION_DAYS days)..."

    find "$BACKUP_DIR" -name "pos_backup_*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete

    local count=$(find "$BACKUP_DIR" -name "pos_backup_*.sql.gz" -type f | wc -l)
    log "✅ Cleanup completed ($count backups remaining)"
}

# Create backup manifest
create_manifest() {
    local manifest_file="$BACKUP_DIR/backup_manifest.txt"

    cat > "$manifest_file" <<EOF
Backup Information
==================
Type: $BACKUP_TYPE
Timestamp: $TIMESTAMP
Date: $(date)
Database: $DB_NAME
File: $(basename "$BACKUP_FILE")
Size: $(du -h "$BACKUP_FILE" | cut -f1)
Checksum: $(sha256sum "$BACKUP_FILE" | cut -d' ' -f1)
EOF

    log "✅ Backup manifest created"
}

# Main backup process
main() {
    log "========================================="
    log "POS Database Backup"
    log "Type: $BACKUP_TYPE"
    log "========================================="

    check_database
    create_backup
    verify_backup
    upload_to_s3
    create_manifest
    cleanup_old_backups

    log "========================================="
    log "✅ Backup completed successfully!"
    log "File: $BACKUP_FILE"
    log "========================================="
}

# Run main
main

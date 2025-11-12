#!/bin/bash

#################################################################
# POS Backend Deployment Script
# Usage: ./scripts/deploy.sh [environment]
# Environments: staging, production
#################################################################

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default environment
ENVIRONMENT="${1:-staging}"

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKUP_DIR="$PROJECT_ROOT/backups"
LOG_FILE="$PROJECT_ROOT/logs/deploy-$(date +%Y%m%d-%H%M%S).log"

# Create directories
mkdir -p "$BACKUP_DIR"
mkdir -p "$(dirname "$LOG_FILE")"

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

# Check if running as correct user
check_user() {
    if [ "$EUID" -eq 0 ]; then
        error "Do not run this script as root"
        exit 1
    fi
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."

    local missing_tools=()

    for tool in docker docker-compose git curl; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done

    if [ ${#missing_tools[@]} -gt 0 ]; then
        error "Missing required tools: ${missing_tools[*]}"
        error "Please install missing tools and try again"
        exit 1
    fi

    log "✅ Prerequisites check passed"
}

# Load environment variables
load_environment() {
    log "Loading environment configuration..."

    local env_file=".env.$ENVIRONMENT"

    if [ ! -f "$PROJECT_ROOT/$env_file" ]; then
        error "Environment file not found: $env_file"
        exit 1
    fi

    # Export environment variables
    set -a
    source "$PROJECT_ROOT/$env_file"
    set +a

    log "✅ Environment loaded: $ENVIRONMENT"
}

# Pull latest code
pull_code() {
    log "Pulling latest code..."

    cd "$PROJECT_ROOT"

    # Check if git repo
    if [ ! -d ".git" ]; then
        error "Not a git repository"
        exit 1
    fi

    # Stash any local changes
    if ! git diff-index --quiet HEAD --; then
        warn "Local changes detected, stashing..."
        git stash
    fi

    # Pull latest
    git pull origin main || {
        error "Failed to pull latest code"
        exit 1
    }

    log "✅ Code updated"
}

# Backup database
backup_database() {
    log "Creating database backup..."

    "$SCRIPT_DIR/backup.sh" || {
        error "Database backup failed"
        exit 1
    }

    log "✅ Database backed up"
}

# Pull Docker images
pull_images() {
    log "Pulling Docker images..."

    cd "$PROJECT_ROOT"

    docker-compose -f docker-compose.prod.yml pull || {
        error "Failed to pull Docker images"
        exit 1
    }

    log "✅ Images pulled"
}

# Run database migrations
run_migrations() {
    log "Running database migrations..."

    cd "$PROJECT_ROOT"

    # Run postgres migrations
    docker-compose -f docker-compose.prod.yml run --rm backend \
        /app/migrate -cmd up -path /app/postgres/migrations || {
        error "PostgreSQL migrations failed"
        return 1
    }

    # Run accounting migrations
    docker-compose -f docker-compose.prod.yml run --rm backend \
        /app/migrate -cmd up -path /app/accounting/migrations || {
        error "Accounting migrations failed"
        return 1
    }

    log "✅ Migrations completed"
}

# Deploy services
deploy_services() {
    log "Deploying services..."

    cd "$PROJECT_ROOT"

    # Start services with zero downtime
    docker-compose -f docker-compose.prod.yml up -d --no-deps --build || {
        error "Failed to deploy services"
        return 1
    }

    log "✅ Services deployed"
}

# Wait for health check
wait_for_health() {
    log "Waiting for services to be healthy..."

    local max_attempts=30
    local attempt=0
    local backend_url="${BACKEND_URL:-http://localhost:8080}"

    while [ $attempt -lt $max_attempts ]; do
        if curl -sf "$backend_url/health" > /dev/null 2>&1; then
            log "✅ Backend is healthy"
            return 0
        fi

        attempt=$((attempt + 1))
        log "Waiting for backend... ($attempt/$max_attempts)"
        sleep 2
    done

    error "Backend health check failed after $max_attempts attempts"
    return 1
}

# Run smoke tests
run_smoke_tests() {
    log "Running smoke tests..."

    local backend_url="${BACKEND_URL:-http://localhost:8080}"

    # Health check
    if ! curl -sf "$backend_url/health" > /dev/null; then
        error "Health check failed"
        return 1
    fi

    # Metrics endpoint
    if ! curl -sf "$backend_url/metrics" > /dev/null; then
        warn "Metrics endpoint not accessible"
    fi

    log "✅ Smoke tests passed"
}

# Cleanup old images
cleanup() {
    log "Cleaning up old Docker images..."

    docker image prune -f || {
        warn "Failed to prune Docker images"
    }

    log "✅ Cleanup completed"
}

# Rollback function
rollback() {
    error "Deployment failed, initiating rollback..."

    # Stop new containers
    docker-compose -f docker-compose.prod.yml down

    # Restore from backup
    "$SCRIPT_DIR/restore.sh" latest || {
        error "Rollback failed - manual intervention required"
        exit 1
    }

    # Restart services
    docker-compose -f docker-compose.prod.yml up -d

    error "Rollback completed"
    exit 1
}

# Main deployment flow
main() {
    log "========================================="
    log "POS Backend Deployment"
    log "Environment: $ENVIRONMENT"
    log "========================================="

    check_user
    check_prerequisites
    load_environment

    # Prompt for confirmation in production
    if [ "$ENVIRONMENT" = "production" ]; then
        read -p "⚠️  You are about to deploy to PRODUCTION. Continue? (yes/no): " confirm
        if [ "$confirm" != "yes" ]; then
            log "Deployment cancelled"
            exit 0
        fi
    fi

    # Execute deployment steps
    pull_code || rollback
    backup_database || rollback
    pull_images || rollback
    run_migrations || rollback
    deploy_services || rollback
    wait_for_health || rollback
    run_smoke_tests || warn "Smoke tests failed, but deployment continues"
    cleanup

    log "========================================="
    log "✅ Deployment completed successfully!"
    log "Environment: $ENVIRONMENT"
    log "Log file: $LOG_FILE"
    log "========================================="
}

# Run main function
main

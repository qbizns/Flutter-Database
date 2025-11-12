#!/bin/bash

#################################################################
# POS Backend - Linux Installation Script
# Supports: Ubuntu, Debian, CentOS, RHEL, Fedora, Arch
#################################################################

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# Configuration
APP_NAME="pos-backend"
WORKER_NAME="pos-worker"
MIGRATE_NAME="pos-migrate"
INSTALL_DIR="/opt/pos-backend"
CONFIG_DIR="/etc/pos-backend"
DATA_DIR="/var/lib/pos-backend"
LOG_DIR="/var/log/pos-backend"
SERVICE_USER="pos"
SERVICE_GROUP="pos"

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        BUILD_ARCH="amd64"
        ;;
    aarch64|arm64)
        BUILD_ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

BIN_SOURCE="build/bin/linux-$BUILD_ARCH"

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

info() {
    echo -e "${CYAN}[INFO]${NC} $1"
}

header() {
    echo -e "${BOLD}${BLUE}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║          POS Backend - Linux Installation                 ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."

    # Check if binaries exist
    if [ ! -f "$BIN_SOURCE/$APP_NAME" ]; then
        error "Binary not found: $BIN_SOURCE/$APP_NAME"
        error "Run 'make build-linux' first"
        exit 1
    fi

    log "✓ Prerequisites check passed"
}

# Create system user
create_user() {
    log "Creating system user and group..."

    if ! getent group "$SERVICE_GROUP" > /dev/null 2>&1; then
        groupadd --system "$SERVICE_GROUP"
        log "Created group: $SERVICE_GROUP"
    else
        info "Group already exists: $SERVICE_GROUP"
    fi

    if ! getent passwd "$SERVICE_USER" > /dev/null 2>&1; then
        useradd --system \
            --gid "$SERVICE_GROUP" \
            --no-create-home \
            --home-dir "$INSTALL_DIR" \
            --shell /usr/sbin/nologin \
            --comment "POS Backend Service User" \
            "$SERVICE_USER"
        log "Created user: $SERVICE_USER"
    else
        info "User already exists: $SERVICE_USER"
    fi
}

# Create directories
create_directories() {
    log "Creating directories..."

    mkdir -p "$INSTALL_DIR"/{bin,config,scripts}
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$DATA_DIR"
    mkdir -p "$LOG_DIR"

    log "✓ Directories created"
}

# Install binaries
install_binaries() {
    log "Installing binaries..."

    cp "$BIN_SOURCE/$APP_NAME" "$INSTALL_DIR/bin/"
    cp "$BIN_SOURCE/$WORKER_NAME" "$INSTALL_DIR/bin/"
    cp "$BIN_SOURCE/$MIGRATE_NAME" "$INSTALL_DIR/bin/"

    chmod +x "$INSTALL_DIR/bin/"*

    # Create symlinks
    ln -sf "$INSTALL_DIR/bin/$APP_NAME" /usr/local/bin/$APP_NAME
    ln -sf "$INSTALL_DIR/bin/$WORKER_NAME" /usr/local/bin/$WORKER_NAME
    ln -sf "$INSTALL_DIR/bin/$MIGRATE_NAME" /usr/local/bin/$MIGRATE_NAME

    log "✓ Binaries installed"
}

# Install configuration
install_config() {
    log "Installing configuration..."

    if [ ! -f "$CONFIG_DIR/config.env" ]; then
        cat > "$CONFIG_DIR/config.env" <<'EOF'
# POS Backend Configuration

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=pos_db
DATABASE_USER=postgres
DATABASE_PASSWORD=changeme
DATABASE_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=please-change-this-to-a-secure-random-string-minimum-64-characters

# Environment
ENVIRONMENT=production

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Metrics
ENABLE_METRICS=true

# Tracing
ENABLE_TRACING=false
JAEGER_ENDPOINT=http://localhost:14268/api/traces

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100
EOF

        log "✓ Default configuration created at $CONFIG_DIR/config.env"
        warn "⚠️  IMPORTANT: Edit $CONFIG_DIR/config.env and set secure values!"
    else
        info "Configuration already exists"
    fi
}

# Install systemd service
install_service() {
    log "Installing systemd service..."

    # API service
    cat > /etc/systemd/system/pos-backend.service <<EOF
[Unit]
Description=POS Backend API Server
Documentation=https://github.com/your-org/pos-backend
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_GROUP
WorkingDirectory=$INSTALL_DIR
EnvironmentFile=$CONFIG_DIR/config.env
ExecStart=$INSTALL_DIR/bin/$APP_NAME
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=10
StandardOutput=append:$LOG_DIR/backend.log
StandardError=append:$LOG_DIR/backend.error.log

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$DATA_DIR $LOG_DIR
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
EOF

    # Worker service
    cat > /etc/systemd/system/pos-worker.service <<EOF
[Unit]
Description=POS Backend Background Worker
Documentation=https://github.com/your-org/pos-backend
After=network.target postgresql.service redis.service pos-backend.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_GROUP
WorkingDirectory=$INSTALL_DIR
EnvironmentFile=$CONFIG_DIR/config.env
Environment="WORKER_MODE=true"
Environment="WORKER_CONCURRENCY=10"
ExecStart=$INSTALL_DIR/bin/$WORKER_NAME
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=10
StandardOutput=append:$LOG_DIR/worker.log
StandardError=append:$LOG_DIR/worker.error.log

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$DATA_DIR $LOG_DIR
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload

    log "✓ Systemd services installed"
}

# Set permissions
set_permissions() {
    log "Setting permissions..."

    chown -R "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR"
    chown -R "$SERVICE_USER:$SERVICE_GROUP" "$CONFIG_DIR"
    chown -R "$SERVICE_USER:$SERVICE_GROUP" "$DATA_DIR"
    chown -R "$SERVICE_USER:$SERVICE_GROUP" "$LOG_DIR"

    chmod 750 "$INSTALL_DIR"
    chmod 750 "$CONFIG_DIR"
    chmod 750 "$DATA_DIR"
    chmod 750 "$LOG_DIR"
    chmod 640 "$CONFIG_DIR/config.env"

    log "✓ Permissions set"
}

# Print post-installation instructions
post_install() {
    echo ""
    echo -e "${BOLD}${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${GREEN}║          Installation Complete! 🎉                        ║${NC}"
    echo -e "${BOLD}${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BOLD}Installation Details:${NC}"
    echo "  • Binaries:       $INSTALL_DIR/bin/"
    echo "  • Configuration:  $CONFIG_DIR/config.env"
    echo "  • Data:           $DATA_DIR"
    echo "  • Logs:           $LOG_DIR"
    echo "  • User:           $SERVICE_USER"
    echo ""
    echo -e "${BOLD}${YELLOW}⚠️  IMPORTANT NEXT STEPS:${NC}"
    echo ""
    echo "1. Edit configuration:"
    echo -e "   ${CYAN}sudo nano $CONFIG_DIR/config.env${NC}"
    echo ""
    echo "2. Set secure JWT secret (64+ characters)"
    echo ""
    echo "3. Configure database connection"
    echo ""
    echo "4. Run database migrations:"
    echo -e "   ${CYAN}sudo -u $SERVICE_USER $MIGRATE_NAME -cmd up${NC}"
    echo ""
    echo "5. Enable and start services:"
    echo -e "   ${CYAN}sudo systemctl enable pos-backend pos-worker${NC}"
    echo -e "   ${CYAN}sudo systemctl start pos-backend pos-worker${NC}"
    echo ""
    echo "6. Check service status:"
    echo -e "   ${CYAN}sudo systemctl status pos-backend${NC}"
    echo -e "   ${CYAN}sudo systemctl status pos-worker${NC}"
    echo ""
    echo "7. View logs:"
    echo -e "   ${CYAN}sudo journalctl -u pos-backend -f${NC}"
    echo -e "   ${CYAN}tail -f $LOG_DIR/backend.log${NC}"
    echo ""
    echo -e "${BOLD}Useful Commands:${NC}"
    echo "  • Start:   ${CYAN}sudo systemctl start pos-backend${NC}"
    echo "  • Stop:    ${CYAN}sudo systemctl stop pos-backend${NC}"
    echo "  • Restart: ${CYAN}sudo systemctl restart pos-backend${NC}"
    echo "  • Status:  ${CYAN}sudo systemctl status pos-backend${NC}"
    echo "  • Logs:    ${CYAN}sudo journalctl -u pos-backend -f${NC}"
    echo ""
    echo -e "${GREEN}API will be available at: ${BOLD}http://localhost:8080${NC}"
    echo ""
}

# Main installation
main() {
    header
    check_root
    check_prerequisites
    create_user
    create_directories
    install_binaries
    install_config
    install_service
    set_permissions
    post_install
}

main

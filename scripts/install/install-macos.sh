#!/bin/bash

#################################################################
# POS Backend - macOS Installation Script
# Supports: macOS 10.15+ (Catalina and later)
# Architecture: Intel (amd64) and Apple Silicon (arm64)
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
INSTALL_DIR="/usr/local/opt/pos-backend"
CONFIG_DIR="/usr/local/etc/pos-backend"
DATA_DIR="/usr/local/var/pos-backend"
LOG_DIR="/usr/local/var/log/pos-backend"
LAUNCH_AGENTS_DIR="$HOME/Library/LaunchAgents"
LAUNCH_DAEMONS_DIR="/Library/LaunchDaemons"

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        BUILD_ARCH="amd64"
        CHIP_TYPE="Intel"
        ;;
    arm64)
        BUILD_ARCH="arm64"
        CHIP_TYPE="Apple Silicon (M1/M2)"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

BIN_SOURCE="build/bin/darwin-$BUILD_ARCH"

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
    echo "║          POS Backend - macOS Installation                 ║"
    echo "║          Architecture: $CHIP_TYPE                    ║"
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
        error "Run 'make build-macos' first"
        exit 1
    fi

    # Check Homebrew (optional but recommended)
    if ! command -v brew &> /dev/null; then
        warn "Homebrew not installed. Install from: https://brew.sh"
        warn "PostgreSQL and Redis can be installed via Homebrew"
    fi

    log "✓ Prerequisites check passed"
}

# Create directories
create_directories() {
    log "Creating directories..."

    mkdir -p "$INSTALL_DIR"/{bin,config,scripts}
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$DATA_DIR"
    mkdir -p "$LOG_DIR"
    mkdir -p "$LAUNCH_DAEMONS_DIR"

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

    # Code sign binaries (if certificate available)
    if command -v codesign &> /dev/null; then
        codesign --force --sign - "$INSTALL_DIR/bin/"* 2>/dev/null || true
        info "Binaries signed"
    fi

    log "✓ Binaries installed"
}

# Install configuration
install_config() {
    log "Installing configuration..."

    if [ ! -f "$CONFIG_DIR/config.env" ]; then
        cat > "$CONFIG_DIR/config.env" <<'EOF'
# POS Backend Configuration

# Server
SERVER_HOST=127.0.0.1
SERVER_PORT=8080

# Database (PostgreSQL via Homebrew: brew install postgresql)
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=pos_db
DATABASE_USER=postgres
DATABASE_PASSWORD=changeme
DATABASE_SSLMODE=disable

# Redis (via Homebrew: brew install redis)
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

# Install launchd service (backend API)
install_launchd_backend() {
    log "Installing launchd service for backend API..."

    cat > "$LAUNCH_DAEMONS_DIR/com.posbackend.api.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.posbackend.api</string>

    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/bin/$APP_NAME</string>
    </array>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>

    <key>WorkingDirectory</key>
    <string>$INSTALL_DIR</string>

    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>

    <key>StandardOutPath</key>
    <string>$LOG_DIR/backend.log</string>

    <key>StandardErrorPath</key>
    <string>$LOG_DIR/backend.error.log</string>

    <key>ThrottleInterval</key>
    <integer>10</integer>

    <key>ProcessType</key>
    <string>Interactive</string>
</dict>
</plist>
EOF

    chmod 644 "$LAUNCH_DAEMONS_DIR/com.posbackend.api.plist"
    log "✓ Backend launchd service installed"
}

# Install launchd service (worker)
install_launchd_worker() {
    log "Installing launchd service for worker..."

    cat > "$LAUNCH_DAEMONS_DIR/com.posbackend.worker.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.posbackend.worker</string>

    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/bin/$WORKER_NAME</string>
    </array>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>

    <key>WorkingDirectory</key>
    <string>$INSTALL_DIR</string>

    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
        <key>WORKER_MODE</key>
        <string>true</string>
        <key>WORKER_CONCURRENCY</key>
        <string>10</string>
    </dict>

    <key>StandardOutPath</key>
    <string>$LOG_DIR/worker.log</string>

    <key>StandardErrorPath</key>
    <string>$LOG_DIR/worker.error.log</string>

    <key>ThrottleInterval</key>
    <integer>10</integer>

    <key>ProcessType</key>
    <string>Background</string>
</dict>
</plist>
EOF

    chmod 644 "$LAUNCH_DAEMONS_DIR/com.posbackend.worker.plist"
    log "✓ Worker launchd service installed"
}

# Set permissions
set_permissions() {
    log "Setting permissions..."

    chmod -R 755 "$INSTALL_DIR"
    chmod -R 755 "$CONFIG_DIR"
    chmod -R 755 "$DATA_DIR"
    chmod -R 755 "$LOG_DIR"
    chmod 600 "$CONFIG_DIR/config.env"

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
    echo "  • Architecture:   $CHIP_TYPE ($BUILD_ARCH)"
    echo "  • Binaries:       $INSTALL_DIR/bin/"
    echo "  • Configuration:  $CONFIG_DIR/config.env"
    echo "  • Data:           $DATA_DIR"
    echo "  • Logs:           $LOG_DIR"
    echo ""
    echo -e "${BOLD}${YELLOW}⚠️  IMPORTANT NEXT STEPS:${NC}"
    echo ""
    echo "1. Install dependencies (if not already installed):"
    echo -e "   ${CYAN}brew install postgresql@15${NC}"
    echo -e "   ${CYAN}brew install redis${NC}"
    echo -e "   ${CYAN}brew services start postgresql@15${NC}"
    echo -e "   ${CYAN}brew services start redis${NC}"
    echo ""
    echo "2. Edit configuration:"
    echo -e "   ${CYAN}sudo nano $CONFIG_DIR/config.env${NC}"
    echo ""
    echo "3. Set secure JWT secret (64+ characters)"
    echo ""
    echo "4. Run database migrations:"
    echo -e "   ${CYAN}$MIGRATE_NAME -cmd up${NC}"
    echo ""
    echo "5. Load and start services:"
    echo -e "   ${CYAN}sudo launchctl load $LAUNCH_DAEMONS_DIR/com.posbackend.api.plist${NC}"
    echo -e "   ${CYAN}sudo launchctl load $LAUNCH_DAEMONS_DIR/com.posbackend.worker.plist${NC}"
    echo ""
    echo "6. Check service status:"
    echo -e "   ${CYAN}sudo launchctl list | grep posbackend${NC}"
    echo ""
    echo "7. View logs:"
    echo -e "   ${CYAN}tail -f $LOG_DIR/backend.log${NC}"
    echo -e "   ${CYAN}tail -f $LOG_DIR/worker.log${NC}"
    echo ""
    echo -e "${BOLD}Useful Commands:${NC}"
    echo "  • Start:   ${CYAN}sudo launchctl load $LAUNCH_DAEMONS_DIR/com.posbackend.api.plist${NC}"
    echo "  • Stop:    ${CYAN}sudo launchctl unload $LAUNCH_DAEMONS_DIR/com.posbackend.api.plist${NC}"
    echo "  • Restart: ${CYAN}sudo launchctl kickstart -k system/com.posbackend.api${NC}"
    echo "  • Status:  ${CYAN}sudo launchctl list | grep posbackend${NC}"
    echo "  • Logs:    ${CYAN}tail -f $LOG_DIR/backend.log${NC}"
    echo ""
    echo -e "${GREEN}API will be available at: ${BOLD}http://localhost:8080${NC}"
    echo ""
}

# Main installation
main() {
    header
    check_root
    check_prerequisites
    create_directories
    install_binaries
    install_config
    install_launchd_backend
    install_launchd_worker
    set_permissions
    post_install
}

main

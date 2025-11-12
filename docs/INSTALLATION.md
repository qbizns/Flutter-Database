# POS Backend - Installation Guide

Complete installation guide for Linux, macOS, and Windows.

---

## 📋 Table of Contents

- [System Requirements](#system-requirements)
- [Quick Start](#quick-start)
- [Linux Installation](#linux-installation)
- [macOS Installation](#macos-installation)
- [Windows Installation](#windows-installation)
- [Docker Installation](#docker-installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [First Run](#first-run)
- [Troubleshooting](#troubleshooting)
- [Uninstallation](#uninstallation)

---

## System Requirements

### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 512MB
- **Disk**: 500MB for application + space for data
- **Network**: Internet connection for dependencies

### Recommended Requirements
- **CPU**: 4+ cores
- **RAM**: 2GB+
- **Disk**: 10GB+ (SSD recommended)
- **Network**: Stable internet connection

### Dependencies
- **PostgreSQL**: 15.x or later
- **Redis**: 7.x or later
- **Go**: 1.21+ (only for building from source)

---

## Quick Start

### Option 1: Using Pre-built Binaries (Recommended)

**Download the latest release** for your platform:
- **Linux**: `pos-backend-v1.0.0-linux-amd64.tar.gz`
- **macOS Intel**: `pos-backend-v1.0.0-macos-amd64.tar.gz`
- **macOS Apple Silicon**: `pos-backend-v1.0.0-macos-arm64.tar.gz`
- **Windows**: `pos-backend-v1.0.0-windows-amd64.zip`

**Extract and run:**

```bash
# Linux/macOS
tar -xzf pos-backend-v1.0.0-*.tar.gz
cd pos-backend-v1.0.0-*/
./pos-backend --version

# Windows (PowerShell)
Expand-Archive pos-backend-v1.0.0-windows-amd64.zip
cd pos-backend-v1.0.0-windows-amd64
.\pos-backend.exe --version
```

### Option 2: Building from Source

```bash
# Clone repository
git clone https://github.com/your-org/pos-backend.git
cd pos-backend

# Build for your platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Install
sudo make install
```

---

## Linux Installation

### Supported Distributions
- ✅ Ubuntu 20.04, 22.04, 24.04
- ✅ Debian 11, 12
- ✅ CentOS 8, 9
- ✅ RHEL 8, 9
- ✅ Fedora 38+
- ✅ Arch Linux

### Method 1: Using Installation Script (Recommended)

**1. Download and build:**

```bash
# Clone repository
git clone https://github.com/your-org/pos-backend.git
cd pos-backend

# Build for Linux
make build-linux
```

**2. Run installation script:**

```bash
# Make script executable
chmod +x scripts/install/install-linux.sh

# Run installation (requires sudo)
sudo ./scripts/install/install-linux.sh
```

**3. Configure:**

```bash
# Edit configuration
sudo nano /etc/pos-backend/config.env

# Set your database credentials, JWT secret, etc.
```

**4. Setup database:**

```bash
# Install PostgreSQL and Redis
sudo apt update
sudo apt install postgresql postgresql-contrib redis-server

# Create database
sudo -u postgres createdb pos_db
sudo -u postgres createuser pos_user
sudo -u postgres psql -c "ALTER USER pos_user PASSWORD 'your_password';"

# Run migrations
sudo -u pos pos-migrate -cmd up
```

**5. Start services:**

```bash
# Enable services to start on boot
sudo systemctl enable pos-backend pos-worker

# Start services
sudo systemctl start pos-backend pos-worker

# Check status
sudo systemctl status pos-backend
sudo systemctl status pos-worker
```

### Method 2: Manual Installation

**1. Create directories:**

```bash
sudo mkdir -p /opt/pos-backend/bin
sudo mkdir -p /etc/pos-backend
sudo mkdir -p /var/lib/pos-backend
sudo mkdir -p /var/log/pos-backend
```

**2. Copy binaries:**

```bash
sudo cp build/bin/linux-amd64/* /opt/pos-backend/bin/
sudo chmod +x /opt/pos-backend/bin/*
sudo ln -s /opt/pos-backend/bin/pos-backend /usr/local/bin/
```

**3. Create service user:**

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin pos
```

**4. Create systemd service:**

```bash
sudo nano /etc/systemd/system/pos-backend.service
```

Copy the service file from `scripts/install/install-linux.sh` or see [systemd configuration](#systemd-configuration).

**5. Reload and start:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable pos-backend
sudo systemctl start pos-backend
```

### systemd Configuration

Service file location: `/etc/systemd/system/pos-backend.service`

```ini
[Unit]
Description=POS Backend API Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=pos
Group=pos
WorkingDirectory=/opt/pos-backend
EnvironmentFile=/etc/pos-backend/config.env
ExecStart=/opt/pos-backend/bin/pos-backend
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Useful Commands

```bash
# Start service
sudo systemctl start pos-backend

# Stop service
sudo systemctl stop pos-backend

# Restart service
sudo systemctl restart pos-backend

# View status
sudo systemctl status pos-backend

# View logs (real-time)
sudo journalctl -u pos-backend -f

# View logs (last 100 lines)
sudo journalctl -u pos-backend -n 100

# Check log file
tail -f /var/log/pos-backend/backend.log
```

---

## macOS Installation

### Supported Versions
- ✅ macOS 10.15 Catalina or later
- ✅ macOS 11 Big Sur
- ✅ macOS 12 Monterey
- ✅ macOS 13 Ventura
- ✅ macOS 14 Sonoma

### Architecture Support
- ✅ Intel (x86_64 / amd64)
- ✅ Apple Silicon (arm64 / M1/M2/M3)

### Method 1: Using Installation Script (Recommended)

**1. Install Homebrew (if not installed):**

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

**2. Install dependencies:**

```bash
brew install postgresql@15 redis
brew services start postgresql@15
brew services start redis
```

**3. Download and build:**

```bash
# Clone repository
git clone https://github.com/your-org/pos-backend.git
cd pos-backend

# Build for macOS
make build-macos
```

**4. Run installation script:**

```bash
# Make script executable
chmod +x scripts/install/install-macos.sh

# Run installation (requires sudo)
sudo ./scripts/install/install-macos.sh
```

**5. Configure:**

```bash
# Edit configuration
sudo nano /usr/local/etc/pos-backend/config.env

# Set your database credentials, JWT secret, etc.
```

**6. Setup database:**

```bash
# Create database
createdb pos_db

# Run migrations
pos-migrate -cmd up
```

**7. Start services:**

```bash
# Load and start backend API
sudo launchctl load /Library/LaunchDaemons/com.posbackend.api.plist

# Load and start worker
sudo launchctl load /Library/LaunchDaemons/com.posbackend.worker.plist

# Check status
sudo launchctl list | grep posbackend
```

### Method 2: Manual Installation

**1. Create directories:**

```bash
sudo mkdir -p /usr/local/opt/pos-backend/bin
sudo mkdir -p /usr/local/etc/pos-backend
sudo mkdir -p /usr/local/var/pos-backend
sudo mkdir -p /usr/local/var/log/pos-backend
```

**2. Copy binaries:**

```bash
# For Intel Macs
sudo cp build/bin/darwin-amd64/* /usr/local/opt/pos-backend/bin/

# For Apple Silicon Macs
sudo cp build/bin/darwin-arm64/* /usr/local/opt/pos-backend/bin/

sudo chmod +x /usr/local/opt/pos-backend/bin/*
sudo ln -s /usr/local/opt/pos-backend/bin/pos-backend /usr/local/bin/
```

**3. Create launchd plist:**

See `scripts/install/install-macos.sh` for complete plist configuration.

### launchd Configuration

Service file location: `/Library/LaunchDaemons/com.posbackend.api.plist`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.posbackend.api</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/opt/pos-backend/bin/pos-backend</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/usr/local/var/log/pos-backend/backend.log</string>
    <key>StandardErrorPath</key>
    <string>/usr/local/var/log/pos-backend/backend.error.log</string>
</dict>
</plist>
```

### Useful Commands

```bash
# Start service
sudo launchctl load /Library/LaunchDaemons/com.posbackend.api.plist

# Stop service
sudo launchctl unload /Library/LaunchDaemons/com.posbackend.api.plist

# Restart service
sudo launchctl kickstart -k system/com.posbackend.api

# Check status
sudo launchctl list | grep posbackend

# View logs
tail -f /usr/local/var/log/pos-backend/backend.log
```

### Apple Silicon Notes

On Apple Silicon Macs (M1/M2/M3), the system automatically runs the correct architecture:
- Use `darwin-arm64` binaries
- PostgreSQL and Redis from Homebrew are native ARM64
- No Rosetta 2 translation needed
- Better performance and battery life

---

## Windows Installation

### Supported Versions
- ✅ Windows 10 (64-bit)
- ✅ Windows 11 (64-bit)
- ✅ Windows Server 2016 or later

### Method 1: Using Installation Script (Recommended)

**1. Install dependencies:**

Download and install:
- **PostgreSQL**: https://www.postgresql.org/download/windows/
- **Redis**: https://github.com/microsoftarchive/redis/releases (or use Memurai: https://www.memurai.com/)

**2. Download and build:**

```powershell
# Clone repository
git clone https://github.com/your-org/pos-backend.git
cd pos-backend

# Build for Windows
make build-windows
```

**3. Run installation script (as Administrator):**

```powershell
# Right-click PowerShell and select "Run as Administrator"

# Set execution policy (if needed)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Run installation
.\scripts\install\install-windows.ps1
```

**4. Configure:**

```powershell
# Edit configuration
notepad "C:\ProgramData\POS Backend\config.env"

# Set your database credentials, JWT secret, etc.
```

**5. Setup database:**

```powershell
# Using psql command (add PostgreSQL bin to PATH)
createdb pos_db

# Run migrations
pos-migrate -cmd up
```

**6. Start services:**

```powershell
# Start backend API service
Start-Service POSBackend

# Start worker service
Start-Service POSWorker

# Check status
Get-Service POSBackend, POSWorker
```

### Method 2: Manual Installation

**1. Create directories:**

```powershell
New-Item -ItemType Directory -Force -Path "C:\Program Files\POS Backend\bin"
New-Item -ItemType Directory -Force -Path "C:\ProgramData\POS Backend"
New-Item -ItemType Directory -Force -Path "C:\ProgramData\POS Backend\logs"
```

**2. Copy binaries:**

```powershell
Copy-Item build\bin\windows-amd64\* "C:\Program Files\POS Backend\bin\"
```

**3. Add to PATH:**

```powershell
$path = [Environment]::GetEnvironmentVariable("Path", "Machine")
$newPath = "$path;C:\Program Files\POS Backend\bin"
[Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
```

**4. Install as Windows Service:**

Use NSSM (Non-Sucking Service Manager):

```powershell
# Download NSSM
Invoke-WebRequest -Uri "https://nssm.cc/release/nssm-2.24.zip" -OutFile "nssm.zip"

# Extract and install service
.\nssm.exe install POSBackend "C:\Program Files\POS Backend\bin\pos-backend.exe"
.\nssm.exe set POSBackend AppDirectory "C:\Program Files\POS Backend"
.\nssm.exe set POSBackend DisplayName "POS Backend API Server"
.\nssm.exe start POSBackend
```

### Windows Service Configuration

Services are managed using NSSM (Non-Sucking Service Manager):

- **Service Name**: `POSBackend` (API), `POSWorker` (Worker)
- **Start Type**: Automatic
- **Log Rotation**: Enabled (10MB per file)
- **Restart**: Automatic on failure

### Useful Commands

```powershell
# Start service
Start-Service POSBackend

# Stop service
Stop-Service POSBackend

# Restart service
Restart-Service POSBackend

# Check status
Get-Service POSBackend

# View logs (real-time)
Get-Content "C:\ProgramData\POS Backend\logs\backend.log" -Tail 50 -Wait

# View service configuration
sc.exe qc POSBackend
```

### Windows Firewall

The installation script automatically adds a firewall rule. To add manually:

```powershell
New-NetFirewallRule -DisplayName "POS Backend API" `
    -Direction Inbound `
    -Program "C:\Program Files\POS Backend\bin\pos-backend.exe" `
    -Action Allow `
    -Profile Domain,Private
```

---

## Docker Installation

### Using Docker Compose (Recommended)

**1. Clone repository:**

```bash
git clone https://github.com/your-org/pos-backend.git
cd pos-backend
```

**2. Configure environment:**

```bash
# Copy example environment file
cp .env.example .env.production

# Edit configuration
nano .env.production
```

**3. Start services:**

```bash
# Start all services
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Check status
docker-compose -f docker-compose.prod.yml ps
```

**4. Run migrations:**

```bash
docker-compose -f docker-compose.prod.yml exec backend pos-migrate -cmd up
```

### Using Docker Only

**1. Pull image:**

```bash
docker pull ghcr.io/your-org/pos-backend:latest
```

**2. Create network:**

```bash
docker network create pos-network
```

**3. Start PostgreSQL:**

```bash
docker run -d \
  --name pos-postgres \
  --network pos-network \
  -e POSTGRES_DB=pos_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=changeme \
  -v pos-data:/var/lib/postgresql/data \
  postgres:15-alpine
```

**4. Start Redis:**

```bash
docker run -d \
  --name pos-redis \
  --network pos-network \
  -v pos-redis:/data \
  redis:7-alpine
```

**5. Start POS Backend:**

```bash
docker run -d \
  --name pos-backend \
  --network pos-network \
  -p 8080:8080 \
  -e DATABASE_HOST=pos-postgres \
  -e DATABASE_NAME=pos_db \
  -e DATABASE_USER=postgres \
  -e DATABASE_PASSWORD=changeme \
  -e REDIS_HOST=pos-redis \
  -e JWT_SECRET=your-secure-secret-key-minimum-64-characters-long \
  ghcr.io/your-org/pos-backend:latest
```

---

## Configuration

Configuration file location varies by platform:
- **Linux**: `/etc/pos-backend/config.env`
- **macOS**: `/usr/local/etc/pos-backend/config.env`
- **Windows**: `C:\ProgramData\POS Backend\config.env`

### Required Configuration

```bash
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=pos_db
DATABASE_USER=postgres
DATABASE_PASSWORD=your_secure_password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT (MUST be at least 64 characters)
JWT_SECRET=your-very-long-secure-random-string-minimum-64-characters-for-production

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=production
```

### Optional Configuration

```bash
# Logging
LOG_LEVEL=info                    # debug, info, warn, error
LOG_FORMAT=json                   # json, text

# Metrics
ENABLE_METRICS=true
METRICS_PORT=8080

# Tracing
ENABLE_TRACING=false
JAEGER_ENDPOINT=http://localhost:14268/api/traces

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# Background Jobs
WORKER_CONCURRENCY=10
JOB_RETENTION_DAYS=30

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
```

### Generate Secure JWT Secret

```bash
# Linux/macOS
openssl rand -base64 64

# Windows (PowerShell)
[Convert]::ToBase64String((1..64 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))

# Go
go run -e 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 64); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

---

## Database Setup

### PostgreSQL Installation

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**macOS:**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Windows:**
Download installer from: https://www.postgresql.org/download/windows/

### Create Database

```bash
# Connect as postgres user
sudo -u postgres psql

# Create database and user
CREATE DATABASE pos_db;
CREATE USER pos_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE pos_db TO pos_user;

# Exit
\q
```

### Run Migrations

```bash
# Linux/macOS
pos-migrate -cmd up

# Windows
pos-migrate.exe -cmd up

# With Docker
docker-compose exec backend pos-migrate -cmd up
```

### Verify Database

```bash
# Connect to database
psql -h localhost -U pos_user -d pos_db

# List tables
\dt

# Check table structure
\d organizations
\d users
\d products
```

---

## First Run

### 1. Start Services

**Linux:**
```bash
sudo systemctl start pos-backend pos-worker
```

**macOS:**
```bash
sudo launchctl load /Library/LaunchDaemons/com.posbackend.api.plist
```

**Windows:**
```powershell
Start-Service POSBackend, POSWorker
```

### 2. Verify Installation

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected response:
{
  "status": "ok",
  "timestamp": "2025-11-12T10:00:00Z",
  "version": "1.0.0"
}

# Check detailed health
curl http://localhost:8080/health/detailed
```

### 3. Create First User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!",
    "name": "Admin User",
    "organization_name": "My Company"
  }'
```

### 4. Test API

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"SecurePassword123!"}' \
  | jq -r '.token')

# List products
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/organizations/{org_id}/products
```

---

## Troubleshooting

### Common Issues

#### 1. Service Won't Start

**Check logs:**
```bash
# Linux
sudo journalctl -u pos-backend -n 50

# macOS
tail -f /usr/local/var/log/pos-backend/backend.log

# Windows
Get-Content "C:\ProgramData\POS Backend\logs\backend.log" -Tail 50
```

**Common causes:**
- Database connection failed
- Redis not running
- Port already in use
- Invalid configuration

#### 2. Database Connection Error

```bash
# Check PostgreSQL is running
# Linux
sudo systemctl status postgresql

# macOS
brew services list | grep postgresql

# Windows
Get-Service *postgresql*

# Test connection
psql -h localhost -U postgres -d pos_db
```

#### 3. Permission Denied

**Linux:**
```bash
# Fix ownership
sudo chown -R pos:pos /opt/pos-backend
sudo chown -R pos:pos /var/lib/pos-backend
sudo chown -R pos:pos /var/log/pos-backend
```

**macOS:**
```bash
# Fix permissions
sudo chmod -R 755 /usr/local/opt/pos-backend
sudo chmod -R 755 /usr/local/var/pos-backend
```

**Windows:**
Run PowerShell as Administrator

#### 4. Port Already in Use

```bash
# Linux/macOS
# Find process using port 8080
sudo lsof -i :8080
# or
sudo netstat -tlnp | grep 8080

# Kill process
sudo kill -9 <PID>

# Windows
# Find process
netstat -ano | findstr :8080

# Kill process
taskkill /PID <PID> /F
```

#### 5. JWT Secret Too Short

Error: "JWT secret must be at least 64 characters"

**Solution:** Generate a secure secret:
```bash
openssl rand -base64 64
```

### Debug Mode

Enable debug logging:

```bash
# Edit config
LOG_LEVEL=debug

# Restart service
# Linux
sudo systemctl restart pos-backend

# macOS
sudo launchctl kickstart -k system/com.posbackend.api

# Windows
Restart-Service POSBackend
```

### Getting Help

1. **Check Logs**: Always check logs first
2. **Documentation**: Read API documentation at `/docs/API_DOCUMENTATION.md`
3. **GitHub Issues**: https://github.com/your-org/pos-backend/issues
4. **Support Email**: support@example.com

---

## Uninstallation

### Linux

```bash
# Stop services
sudo systemctl stop pos-backend pos-worker
sudo systemctl disable pos-backend pos-worker

# Remove services
sudo rm /etc/systemd/system/pos-backend.service
sudo rm /etc/systemd/system/pos-worker.service
sudo systemctl daemon-reload

# Remove files
sudo rm -rf /opt/pos-backend
sudo rm -rf /etc/pos-backend
sudo rm -rf /var/lib/pos-backend
sudo rm -rf /var/log/pos-backend

# Remove user
sudo userdel pos

# Remove symlinks
sudo rm /usr/local/bin/pos-backend
sudo rm /usr/local/bin/pos-worker
sudo rm /usr/local/bin/pos-migrate
```

### macOS

```bash
# Stop services
sudo launchctl unload /Library/LaunchDaemons/com.posbackend.api.plist
sudo launchctl unload /Library/LaunchDaemons/com.posbackend.worker.plist

# Remove services
sudo rm /Library/LaunchDaemons/com.posbackend.*.plist

# Remove files
sudo rm -rf /usr/local/opt/pos-backend
sudo rm -rf /usr/local/etc/pos-backend
sudo rm -rf /usr/local/var/pos-backend

# Remove symlinks
sudo rm /usr/local/bin/pos-backend
sudo rm /usr/local/bin/pos-worker
sudo rm /usr/local/bin/pos-migrate
```

### Windows

```powershell
# Stop and remove services
Stop-Service POSBackend, POSWorker
sc.exe delete POSBackend
sc.exe delete POSWorker

# Remove files
Remove-Item -Recurse -Force "C:\Program Files\POS Backend"
Remove-Item -Recurse -Force "C:\ProgramData\POS Backend"

# Remove from PATH (optional)
# Edit system environment variables manually
```

### Database Cleanup (Optional)

⚠️ **WARNING**: This will delete all data!

```bash
# Backup first!
pg_dump -U postgres pos_db > backup.sql

# Drop database
dropdb pos_db

# Drop user
psql -U postgres -c "DROP USER pos_user;"
```

---

## Next Steps

After installation:

1. ✅ **Configure**: Set secure passwords and secrets
2. ✅ **Setup Database**: Run migrations
3. ✅ **Start Services**: Enable and start backend and worker
4. ✅ **Verify**: Test health endpoint
5. ✅ **Create User**: Register first admin user
6. ✅ **Monitor**: Check logs and metrics
7. ✅ **Secure**: Configure firewall and HTTPS
8. ✅ **Backup**: Set up automated backups

---

**Documentation Version**: 1.0.0
**Last Updated**: 2025-11-12
**Support**: support@example.com

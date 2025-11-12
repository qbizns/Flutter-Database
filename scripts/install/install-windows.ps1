#################################################################
# POS Backend - Windows Installation Script
# Supports: Windows 10, Windows 11, Windows Server 2016+
# Architecture: x64 (amd64)
#################################################################

#Requires -RunAsAdministrator

# Configuration
$AppName = "pos-backend"
$WorkerName = "pos-worker"
$MigrateName = "pos-migrate"
$InstallDir = "C:\Program Files\POS Backend"
$ConfigDir = "C:\ProgramData\POS Backend"
$DataDir = "C:\ProgramData\POS Backend\data"
$LogDir = "C:\ProgramData\POS Backend\logs"
$ServiceName = "POSBackend"
$WorkerServiceName = "POSWorker"
$BinSource = "build\bin\windows-amd64"

# Colors for output
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Log {
    param($Message)
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Write-ColorOutput Green "[$timestamp] $Message"
}

function Error {
    param($Message)
    Write-ColorOutput Red "[ERROR] $Message"
}

function Warn {
    param($Message)
    Write-ColorOutput Yellow "[WARNING] $Message"
}

function Info {
    param($Message)
    Write-ColorOutput Cyan "[INFO] $Message"
}

function Header {
    Write-ColorOutput Blue @"
╔════════════════════════════════════════════════════════════╗
║          POS Backend - Windows Installation                ║
║          Architecture: x64 (64-bit)                        ║
╚════════════════════════════════════════════════════════════╝
"@
}

# Check if running as Administrator
function Check-Administrator {
    $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        Error "This script must be run as Administrator"
        Error "Right-click PowerShell and select 'Run as Administrator'"
        exit 1
    }
}

# Check prerequisites
function Check-Prerequisites {
    Log "Checking prerequisites..."

    # Check if binaries exist
    if (-not (Test-Path "$BinSource\$AppName.exe")) {
        Error "Binary not found: $BinSource\$AppName.exe"
        Error "Run 'make build-windows' first"
        exit 1
    }

    # Check .NET Framework (for NSSM service wrapper)
    $netVersion = (Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full" -ErrorAction SilentlyContinue).Release
    if ($netVersion -lt 378389) {
        Warn ".NET Framework 4.5 or higher recommended"
    }

    Log "✓ Prerequisites check passed"
}

# Create directories
function Create-Directories {
    Log "Creating directories..."

    New-Item -ItemType Directory -Force -Path "$InstallDir\bin" | Out-Null
    New-Item -ItemType Directory -Force -Path "$InstallDir\config" | Out-Null
    New-Item -ItemType Directory -Force -Path "$InstallDir\scripts" | Out-Null
    New-Item -ItemType Directory -Force -Path $ConfigDir | Out-Null
    New-Item -ItemType Directory -Force -Path $DataDir | Out-Null
    New-Item -ItemType Directory -Force -Path $LogDir | Out-Null

    Log "✓ Directories created"
}

# Install binaries
function Install-Binaries {
    Log "Installing binaries..."

    Copy-Item "$BinSource\$AppName.exe" "$InstallDir\bin\" -Force
    Copy-Item "$BinSource\$WorkerName.exe" "$InstallDir\bin\" -Force
    Copy-Item "$BinSource\$MigrateName.exe" "$InstallDir\bin\" -Force

    # Add to PATH
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    if ($currentPath -notlike "*$InstallDir\bin*") {
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$currentPath;$InstallDir\bin",
            "Machine"
        )
        Log "Added to system PATH"
    }

    Log "✓ Binaries installed"
}

# Install configuration
function Install-Configuration {
    Log "Installing configuration..."

    $configFile = "$ConfigDir\config.env"
    if (-not (Test-Path $configFile)) {
        $config = @"
# POS Backend Configuration

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database (PostgreSQL)
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
"@

        Set-Content -Path $configFile -Value $config
        Log "✓ Default configuration created at $configFile"
        Warn "⚠️  IMPORTANT: Edit $configFile and set secure values!"
    } else {
        Info "Configuration already exists"
    }
}

# Download and install NSSM (Non-Sucking Service Manager)
function Install-NSSM {
    Log "Checking for NSSM (service wrapper)..."

    $nssmPath = "$InstallDir\bin\nssm.exe"

    if (-not (Test-Path $nssmPath)) {
        Log "Downloading NSSM..."
        $nssmUrl = "https://nssm.cc/release/nssm-2.24.zip"
        $nssmZip = "$env:TEMP\nssm.zip"

        try {
            Invoke-WebRequest -Uri $nssmUrl -OutFile $nssmZip -UseBasicParsing
            Expand-Archive -Path $nssmZip -DestinationPath $env:TEMP -Force

            # Copy the appropriate architecture
            if ([Environment]::Is64BitOperatingSystem) {
                Copy-Item "$env:TEMP\nssm-2.24\win64\nssm.exe" $nssmPath -Force
            } else {
                Copy-Item "$env:TEMP\nssm-2.24\win32\nssm.exe" $nssmPath -Force
            }

            # Cleanup
            Remove-Item $nssmZip -Force
            Remove-Item "$env:TEMP\nssm-2.24" -Recurse -Force

            Log "✓ NSSM installed"
        } catch {
            Error "Failed to download NSSM: $_"
            Error "Please download manually from https://nssm.cc"
            exit 1
        }
    } else {
        Info "NSSM already installed"
    }
}

# Install Windows Service (Backend API)
function Install-BackendService {
    Log "Installing backend API service..."

    $nssmPath = "$InstallDir\bin\nssm.exe"
    $exePath = "$InstallDir\bin\$AppName.exe"
    $configFile = "$ConfigDir\config.env"

    # Check if service already exists
    $existingService = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($existingService) {
        Log "Removing existing service..."
        & $nssmPath stop $ServiceName
        & $nssmPath remove $ServiceName confirm
    }

    # Install service
    & $nssmPath install $ServiceName $exePath

    # Configure service
    & $nssmPath set $ServiceName DisplayName "POS Backend API Server"
    & $nssmPath set $ServiceName Description "POS Backend API Server for Point of Sale and Accounting"
    & $nssmPath set $ServiceName Start SERVICE_AUTO_START
    & $nssmPath set $ServiceName AppDirectory "$InstallDir"
    & $nssmPath set $ServiceName AppExit Default Restart
    & $nssmPath set $ServiceName AppRestartDelay 10000
    & $nssmPath set $ServiceName AppStdout "$LogDir\backend.log"
    & $nssmPath set $ServiceName AppStderr "$LogDir\backend.error.log"
    & $nssmPath set $ServiceName AppRotateFiles 1
    & $nssmPath set $ServiceName AppRotateOnline 1
    & $nssmPath set $ServiceName AppRotateBytes 10485760

    # Set environment from config file
    if (Test-Path $configFile) {
        Get-Content $configFile | Where-Object { $_ -match '=' -and $_ -notmatch '^#' } | ForEach-Object {
            $parts = $_ -split '=', 2
            $key = $parts[0].Trim()
            $value = $parts[1].Trim()
            & $nssmPath set $ServiceName AppEnvironmentExtra "$key=$value"
        }
    }

    Log "✓ Backend API service installed"
}

# Install Windows Service (Worker)
function Install-WorkerService {
    Log "Installing worker service..."

    $nssmPath = "$InstallDir\bin\nssm.exe"
    $exePath = "$InstallDir\bin\$WorkerName.exe"
    $configFile = "$ConfigDir\config.env"

    # Check if service already exists
    $existingService = Get-Service -Name $WorkerServiceName -ErrorAction SilentlyContinue
    if ($existingService) {
        Log "Removing existing service..."
        & $nssmPath stop $WorkerServiceName
        & $nssmPath remove $WorkerServiceName confirm
    }

    # Install service
    & $nssmPath install $WorkerServiceName $exePath

    # Configure service
    & $nssmPath set $WorkerServiceName DisplayName "POS Backend Worker"
    & $nssmPath set $WorkerServiceName Description "POS Backend Background Job Worker"
    & $nssmPath set $WorkerServiceName Start SERVICE_AUTO_START
    & $nssmPath set $WorkerServiceName AppDirectory "$InstallDir"
    & $nssmPath set $WorkerServiceName AppExit Default Restart
    & $nssmPath set $WorkerServiceName AppRestartDelay 10000
    & $nssmPath set $WorkerServiceName AppStdout "$LogDir\worker.log"
    & $nssmPath set $WorkerServiceName AppStderr "$LogDir\worker.error.log"
    & $nssmPath set $WorkerServiceName AppRotateFiles 1
    & $nssmPath set $WorkerServiceName AppRotateOnline 1
    & $nssmPath set $WorkerServiceName AppRotateBytes 10485760

    # Add worker-specific environment
    & $nssmPath set $WorkerServiceName AppEnvironmentExtra "WORKER_MODE=true"
    & $nssmPath set $WorkerServiceName AppEnvironmentExtra "WORKER_CONCURRENCY=10"

    # Set environment from config file
    if (Test-Path $configFile) {
        Get-Content $configFile | Where-Object { $_ -match '=' -and $_ -notmatch '^#' } | ForEach-Object {
            $parts = $_ -split '=', 2
            $key = $parts[0].Trim()
            $value = $parts[1].Trim()
            & $nssmPath set $WorkerServiceName AppEnvironmentExtra "$key=$value"
        }
    }

    Log "✓ Worker service installed"
}

# Configure Windows Firewall
function Configure-Firewall {
    Log "Configuring Windows Firewall..."

    $ruleName = "POS Backend API"
    $existingRule = Get-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue

    if (-not $existingRule) {
        New-NetFirewallRule -DisplayName $ruleName `
            -Direction Inbound `
            -Program "$InstallDir\bin\$AppName.exe" `
            -Action Allow `
            -Profile Domain,Private `
            -Description "Allow incoming connections to POS Backend API"
        Log "✓ Firewall rule created"
    } else {
        Info "Firewall rule already exists"
    }
}

# Print post-installation instructions
function Post-Install {
    Write-Output ""
    Write-ColorOutput Green @"
╔════════════════════════════════════════════════════════════╗
║          Installation Complete! 🎉                        ║
╚════════════════════════════════════════════════════════════╝
"@
    Write-Output ""
    Write-Output "Installation Details:"
    Write-Output "  • Binaries:       $InstallDir\bin\"
    Write-Output "  • Configuration:  $ConfigDir\config.env"
    Write-Output "  • Data:           $DataDir"
    Write-Output "  • Logs:           $LogDir"
    Write-Output ""
    Write-ColorOutput Yellow "⚠️  IMPORTANT NEXT STEPS:"
    Write-Output ""
    Write-Output "1. Install dependencies:"
    Write-Output "   • PostgreSQL: https://www.postgresql.org/download/windows/"
    Write-Output "   • Redis: https://github.com/microsoftarchive/redis/releases"
    Write-Output ""
    Write-Output "2. Edit configuration:"
    Write-Output "   notepad `"$ConfigDir\config.env`""
    Write-Output ""
    Write-Output "3. Set secure JWT secret (64+ characters)"
    Write-Output ""
    Write-Output "4. Run database migrations:"
    Write-Output "   pos-migrate -cmd up"
    Write-Output ""
    Write-Output "5. Start services:"
    Write-Output "   Start-Service $ServiceName"
    Write-Output "   Start-Service $WorkerServiceName"
    Write-Output ""
    Write-Output "6. Check service status:"
    Write-Output "   Get-Service $ServiceName"
    Write-Output "   Get-Service $WorkerServiceName"
    Write-Output ""
    Write-Output "7. View logs:"
    Write-Output "   Get-Content $LogDir\backend.log -Tail 50 -Wait"
    Write-Output ""
    Write-Output "Useful Commands:"
    Write-Output "  • Start:   Start-Service $ServiceName"
    Write-Output "  • Stop:    Stop-Service $ServiceName"
    Write-Output "  • Restart: Restart-Service $ServiceName"
    Write-Output "  • Status:  Get-Service $ServiceName"
    Write-Output "  • Logs:    Get-Content $LogDir\backend.log -Tail 50 -Wait"
    Write-Output ""
    Write-ColorOutput Green "API will be available at: http://localhost:8080"
    Write-Output ""
    Write-Output "NOTE: Restart your terminal to use the 'pos-backend' command"
    Write-Output ""
}

# Main installation
function Main {
    try {
        Header
        Check-Administrator
        Check-Prerequisites
        Create-Directories
        Install-Binaries
        Install-Configuration
        Install-NSSM
        Install-BackendService
        Install-WorkerService
        Configure-Firewall
        Post-Install
    } catch {
        Error "Installation failed: $_"
        Error $_.ScriptStackTrace
        exit 1
    }
}

# Run installation
Main

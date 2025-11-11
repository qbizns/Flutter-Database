# Deployment Guide

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Database Setup](#database-setup)
4. [Application Deployment](#application-deployment)
5. [Post-Deployment Verification](#post-deployment-verification)
6. [Monitoring & Alerting](#monitoring--alerting)
7. [Backup & Recovery](#backup--recovery)
8. [Troubleshooting](#troubleshooting)
9. [Rollback Procedures](#rollback-procedures)

## Prerequisites

### Infrastructure Requirements

**Production Environment:**
- **Compute**: 2 vCPUs, 4GB RAM minimum (recommended: 4 vCPUs, 8GB RAM)
- **Database**: PostgreSQL 16+ with 20GB storage minimum
- **Network**: HTTPS support, load balancer (optional but recommended)
- **Monitoring**: Prometheus + Grafana stack
- **Backup**: Automated daily backups to S3-compatible storage

**Software Dependencies:**
- Go 1.22+
- PostgreSQL 16+
- Docker (optional, for containerized deployment)
- nginx or similar reverse proxy
- systemd or equivalent init system

### Access Requirements

- SSH access to production servers
- Database admin credentials
- Docker registry access (if using containers)
- SSL/TLS certificates for HTTPS
- Environment variable management system (e.g., Vault, AWS Secrets Manager)

## Environment Setup

### 1. Environment Variables

Create `.env` file or configure environment variables:

```bash
# Server Configuration
SERVER_ENV=production
SERVER_API_PORT=8080
SERVER_METRICS_PORT=9090
SERVER_HOST=0.0.0.0

# Database Configuration
DATABASE_HOST=your-postgres-host.example.com
DATABASE_PORT=5432
DATABASE_NAME=pos_production
DATABASE_USER=pos_app_user
DATABASE_PASSWORD=<secure-password>
DATABASE_SSLMODE=require
DATABASE_MAX_OPEN_CONNS=25
DATABASE_MAX_IDLE_CONNS=5
DATABASE_CONN_MAX_LIFETIME=5m

# JWT Configuration
JWT_SECRET=<generate-secure-random-string>
JWT_ACCESS_TOKEN_EXPIRY=1h
JWT_REFRESH_TOKEN_EXPIRY=168h

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://app.yourcompany.com,https://admin.yourcompany.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-CSRF-Token

# Logging
LOGGING_LEVEL=info
LOGGING_FORMAT=json

# Metrics
METRICS_METRICS_PORT=9090
METRICS_METRICS_PATH=/metrics
```

### 2. Generate Secure Secrets

```bash
# Generate JWT secret (256-bit)
openssl rand -base64 32

# Generate database password
openssl rand -base64 24
```

### 3. SSL/TLS Certificates

Place certificates in secure location:

```bash
/etc/ssl/private/yourcompany.key
/etc/ssl/certs/yourcompany.crt
/etc/ssl/certs/yourcompany-chain.crt
```

## Database Setup

### 1. Create Database and User

```sql
-- Connect as postgres superuser
psql -U postgres

-- Create database
CREATE DATABASE pos_production
    WITH ENCODING 'UTF8'
    LC_COLLATE='en_US.UTF-8'
    LC_CTYPE='en_US.UTF-8'
    TEMPLATE=template0;

-- Create application user
CREATE USER pos_app_user WITH PASSWORD '<secure-password>';

-- Grant privileges
GRANT CONNECT ON DATABASE pos_production TO pos_app_user;

-- Connect to database
\c pos_production

-- Grant schema privileges
GRANT USAGE, CREATE ON SCHEMA public TO pos_app_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pos_app_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO pos_app_user;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO pos_app_user;

-- Set default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON TABLES TO pos_app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON SEQUENCES TO pos_app_user;

-- Enable Row-Level Security
ALTER DATABASE pos_production SET row_security = on;
```

### 2. Run Migrations

```bash
# Download migration binary
wget https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz
tar xvf migrate.linux-amd64.tar.gz
sudo mv migrate /usr/local/bin/

# Run migrations
migrate -path ./db/migrations \
    -database "postgres://pos_app_user:<password>@<host>:5432/pos_production?sslmode=require" \
    up

# Verify migration version
migrate -path ./db/migrations \
    -database "postgres://pos_app_user:<password>@<host>:5432/pos_production?sslmode=require" \
    version
```

### 3. Seed Initial Data

```bash
# Seed system roles and permissions
psql -U pos_app_user -d pos_production -f db/seeds/001_permissions.sql
psql -U pos_app_user -d pos_production -f db/seeds/002_roles.sql
psql -U pos_app_user -d pos_production -f db/seeds/003_role_permissions.sql

# Create first organization and admin user (run application with seed flag)
./pos-backend seed --admin-email=admin@yourcompany.com --admin-password=<secure-password>
```

### 4. Configure Database Performance

```sql
-- Optimize PostgreSQL for production
ALTER SYSTEM SET shared_buffers = '2GB';
ALTER SYSTEM SET effective_cache_size = '6GB';
ALTER SYSTEM SET maintenance_work_mem = '512MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;
ALTER SYSTEM SET work_mem = '20MB';
ALTER SYSTEM SET min_wal_size = '1GB';
ALTER SYSTEM SET max_wal_size = '4GB';

-- Reload configuration
SELECT pg_reload_conf();
```

## Application Deployment

### Method 1: Systemd Service (Recommended)

#### 1. Build Application

```bash
# On build server
cd /path/to/pos-backend
go build -o pos-backend-v1.0.0 -ldflags="-s -w" cmd/api/main.go

# Verify binary
./pos-backend-v1.0.0 --version
```

#### 2. Deploy Binary

```bash
# Copy to production server
scp pos-backend-v1.0.0 production-server:/opt/pos-backend/

# On production server
cd /opt/pos-backend
sudo mv pos-backend-v1.0.0 pos-backend
sudo chmod +x pos-backend
```

#### 3. Create Systemd Service

```bash
# Create service file
sudo nano /etc/systemd/system/pos-backend.service
```

```ini
[Unit]
Description=POS Backend API Server
After=network.target postgresql.service
Wants=network-online.target

[Service]
Type=simple
User=pos
Group=pos
WorkingDirectory=/opt/pos-backend
EnvironmentFile=/opt/pos-backend/.env
ExecStart=/opt/pos-backend/pos-backend
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=pos-backend

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/pos-backend/logs

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

#### 4. Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable pos-backend

# Start service
sudo systemctl start pos-backend

# Check status
sudo systemctl status pos-backend

# View logs
sudo journalctl -u pos-backend -f
```

### Method 2: Docker Deployment

#### 1. Build Docker Image

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080 9090
CMD ["./main"]
```

```bash
# Build image
docker build -t pos-backend:v1.0.0 .

# Tag for registry
docker tag pos-backend:v1.0.0 your-registry.com/pos-backend:v1.0.0

# Push to registry
docker push your-registry.com/pos-backend:v1.0.0
```

#### 2. Deploy with Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  pos-backend:
    image: your-registry.com/pos-backend:v1.0.0
    container_name: pos-backend
    restart: always
    ports:
      - "8080:8080"
      - "9090:9090"
    environment:
      - SERVER_ENV=production
      - DATABASE_HOST=postgres
      - DATABASE_NAME=pos_production
      - DATABASE_USER=pos_app_user
      - DATABASE_PASSWORD=${DB_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
    env_file:
      - .env
    depends_on:
      - postgres
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  postgres:
    image: postgres:16-alpine
    container_name: postgres
    restart: always
    environment:
      - POSTGRES_DB=pos_production
      - POSTGRES_USER=pos_app_user
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U pos_app_user"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres-data:
```

```bash
# Deploy
docker-compose up -d

# Check logs
docker-compose logs -f pos-backend

# Check status
docker-compose ps
```

### Method 3: Kubernetes Deployment

#### 1. Create Kubernetes Manifests

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pos-backend
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pos-backend
  template:
    metadata:
      labels:
        app: pos-backend
    spec:
      containers:
      - name: pos-backend
        image: your-registry.com/pos-backend:v1.0.0
        ports:
        - containerPort: 8080
          name: api
        - containerPort: 9090
          name: metrics
        env:
        - name: SERVER_ENV
          value: "production"
        - name: DATABASE_HOST
          valueFrom:
            secretKeyRef:
              name: pos-secrets
              key: db-host
        - name: DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: pos-secrets
              key: db-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: pos-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: pos-backend
  namespace: production
spec:
  selector:
    app: pos-backend
  ports:
  - name: api
    port: 80
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pos-backend-ingress
  namespace: production
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.yourcompany.com
    secretName: pos-backend-tls
  rules:
  - host: api.yourcompany.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: pos-backend
            port:
              number: 80
```

```bash
# Apply manifests
kubectl apply -f deployment.yaml

# Check deployment
kubectl get deployments -n production
kubectl get pods -n production
kubectl logs -f deployment/pos-backend -n production
```

## Post-Deployment Verification

### 1. Health Check

```bash
# Check API health
curl https://api.yourcompany.com/health

# Expected response
{
  "status": "ok",
  "timestamp": "2025-11-11T10:00:00Z",
  "version": "1.0.0"
}
```

### 2. Metrics Endpoint

```bash
# Check Prometheus metrics
curl https://api.yourcompany.com/metrics

# Should see metrics like:
# pos_backend_http_requests_total
# pos_backend_http_request_duration_seconds
# pos_backend_active_connections
```

### 3. Database Connectivity

```bash
# Test database connection
psql -h <host> -U pos_app_user -d pos_production -c "SELECT COUNT(*) FROM migrations;"

# Run application smoke test
curl -X POST https://api.yourcompany.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@yourcompany.com","password":"<admin-password>"}'
```

### 4. Authentication Flow

```bash
# Login
TOKEN=$(curl -s -X POST https://api.yourcompany.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@yourcompany.com","password":"<password>"}' | jq -r '.access_token')

# Test authenticated endpoint
curl -H "Authorization: Bearer $TOKEN" https://api.yourcompany.com/api/users/me
```

### 5. Performance Baseline

```bash
# Run load test
ab -n 1000 -c 10 https://api.yourcompany.com/health

# Expected: <100ms response time, 0% errors
```

## Monitoring & Alerting

### 1. Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'pos-backend'
    static_configs:
      - targets: ['localhost:9090']
    metrics_path: '/metrics'
```

### 2. Key Metrics to Monitor

- `pos_backend_http_requests_total` - Total HTTP requests
- `pos_backend_http_request_duration_seconds` - Request latency
- `pos_backend_active_connections` - Active database connections
- `pos_backend_db_query_duration_seconds` - Database query performance
- `process_cpu_seconds_total` - CPU usage
- `process_resident_memory_bytes` - Memory usage

### 3. Alert Rules

```yaml
# alerts.yml
groups:
  - name: pos_backend_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(pos_backend_http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(pos_backend_http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "95th percentile latency > 1s"

      - alert: DatabaseConnectionsHigh
        expr: pos_backend_active_connections > 20
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Database connection pool near capacity"
```

## Backup & Recovery

### 1. Database Backups

```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/backups/postgres"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/pos_production_$TIMESTAMP.sql.gz"

pg_dump -h <host> -U pos_app_user -d pos_production | gzip > $BACKUP_FILE

# Upload to S3
aws s3 cp $BACKUP_FILE s3://your-backup-bucket/postgres/

# Keep only last 30 days locally
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
```

Add to crontab:
```bash
0 2 * * * /opt/scripts/postgres-backup.sh >> /var/log/postgres-backup.log 2>&1
```

### 2. Application State Backup

```bash
# Backup uploaded files, logs, etc.
tar -czf /backups/app-state-$(date +%Y%m%d).tar.gz \
  /opt/pos-backend/uploads \
  /opt/pos-backend/logs

aws s3 cp /backups/app-state-$(date +%Y%m%d).tar.gz \
  s3://your-backup-bucket/app-state/
```

### 3. Recovery Procedure

```bash
# Download latest backup
aws s3 cp s3://your-backup-bucket/postgres/pos_production_latest.sql.gz /tmp/

# Restore database
gunzip < /tmp/pos_production_latest.sql.gz | psql -h <host> -U pos_app_user -d pos_production

# Restart application
sudo systemctl restart pos-backend

# Verify
curl https://api.yourcompany.com/health
```

## Troubleshooting

### Issue: Application Won't Start

```bash
# Check logs
sudo journalctl -u pos-backend -n 100

# Common issues:
# 1. Database connection failure
#    - Verify DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD
#    - Check network connectivity: telnet <db-host> 5432

# 2. Port already in use
#    - Check: sudo lsof -i :8080
#    - Kill process or change port

# 3. Permission issues
#    - Check file ownership: ls -la /opt/pos-backend
#    - Fix: sudo chown -R pos:pos /opt/pos-backend
```

### Issue: High CPU Usage

```bash
# Check CPU usage
top -p $(pgrep pos-backend)

# Profile application
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof

# Common causes:
# 1. N+1 queries - Check slow query log
# 2. Infinite loops - Check recent code changes
# 3. Memory leaks - Check heap profile
```

### Issue: Database Connection Pool Exhausted

```bash
# Check active connections
SELECT count(*) FROM pg_stat_activity WHERE datname = 'pos_production';

# Kill idle connections
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'pos_production'
AND state = 'idle'
AND state_change < NOW() - INTERVAL '10 minutes';

# Increase pool size (temporarily)
# Edit .env: DATABASE_MAX_OPEN_CONNS=50
sudo systemctl restart pos-backend
```

### Issue: Slow API Response Times

```bash
# Check database query performance
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

# Enable query logging (temporarily)
ALTER SYSTEM SET log_min_duration_statement = 100;
SELECT pg_reload_conf();

# Check application metrics
curl http://localhost:9090/metrics | grep http_request_duration
```

## Rollback Procedures

### 1. Application Rollback

```bash
# Stop current version
sudo systemctl stop pos-backend

# Restore previous binary
sudo mv /opt/pos-backend/pos-backend /opt/pos-backend/pos-backend.failed
sudo mv /opt/pos-backend/pos-backend.previous /opt/pos-backend/pos-backend

# Start previous version
sudo systemctl start pos-backend

# Verify
curl https://api.yourcompany.com/health
```

### 2. Database Rollback

```bash
# Rollback last migration
migrate -path ./db/migrations \
    -database "postgres://pos_app_user:<password>@<host>:5432/pos_production?sslmode=require" \
    down 1

# Verify migration version
migrate -path ./db/migrations \
    -database "postgres://pos_app_user:<password>@<host>:5432/pos_production?sslmode=require" \
    version
```

### 3. Full Rollback

```bash
# 1. Rollback application
sudo systemctl stop pos-backend
sudo mv /opt/pos-backend/pos-backend.previous /opt/pos-backend/pos-backend

# 2. Rollback database
migrate down 1

# 3. Restart
sudo systemctl start pos-backend

# 4. Verify
curl https://api.yourcompany.com/health
```

## Maintenance

### Regular Tasks

**Daily:**
- Check application logs for errors
- Verify backups completed successfully
- Monitor key metrics (error rate, latency, connections)

**Weekly:**
- Review slow query log
- Check disk space usage
- Update OS security patches

**Monthly:**
- Rotate logs
- Review and optimize database indexes
- Update dependencies (security patches only)

### Maintenance Window Procedure

```bash
# 1. Announce maintenance (30 min before)
# Send notification to users

# 2. Enable maintenance mode
# Create maintenance page in nginx/load balancer

# 3. Stop application
sudo systemctl stop pos-backend

# 4. Perform maintenance
#    - Run migrations
#    - Update application
#    - Database maintenance (VACUUM, ANALYZE)

# 5. Start application
sudo systemctl start pos-backend

# 6. Verify
curl https://api.yourcompany.com/health

# 7. Disable maintenance mode
# Remove maintenance page

# 8. Announce completion
```

## Security Checklist

- [ ] JWT secret is strong and secure
- [ ] Database passwords are rotated regularly
- [ ] SSL/TLS certificates are valid and up-to-date
- [ ] Firewall rules limit access to necessary ports only
- [ ] Database backups are encrypted
- [ ] Application logs don't contain sensitive data
- [ ] Rate limiting is enabled
- [ ] CSRF protection is active
- [ ] Security headers are configured
- [ ] Dependency vulnerabilities are monitored

## Support

For production issues:
- **Critical Issues**: Contact on-call engineer (PagerDuty)
- **Non-Critical Issues**: Create ticket in issue tracker
- **Questions**: Refer to this guide and AUTHORIZATION_GUIDE.md

---

**Document Version**: 1.0.0
**Last Updated**: 2025-11-11
**Maintained By**: DevOps Team

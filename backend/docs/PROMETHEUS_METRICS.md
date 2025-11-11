# Prometheus Metrics Documentation

## Overview

The POS Backend exposes Prometheus metrics for monitoring application health, performance, and business operations.

## Metrics Endpoint

**Endpoint:** `GET /metrics`
**Format:** Prometheus text-based exposition format
**Authentication:** None (should be restricted via network/firewall rules)

## Available Metrics

### HTTP Request Metrics

All HTTP metrics are prefixed with `pos_backend_http_`.

#### `pos_backend_http_requests_total`
- **Type:** Counter
- **Labels:** `method`, `path`, `status`
- **Description:** Total number of HTTP requests processed
- **Example:**
  ```
  pos_backend_http_requests_total{method="GET",path="/api/v1/products",status="200"} 1543
  pos_backend_http_requests_total{method="POST",path="/api/v1/auth/login",status="401"} 23
  ```

#### `pos_backend_http_request_duration_seconds`
- **Type:** Histogram
- **Labels:** `method`, `path`, `status`
- **Buckets:** 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 seconds
- **Description:** HTTP request latency in seconds
- **Example:**
  ```
  pos_backend_http_request_duration_seconds_bucket{method="GET",path="/api/v1/sales",status="200",le="0.1"} 1234
  pos_backend_http_request_duration_seconds_sum{method="GET",path="/api/v1/sales",status="200"} 45.6
  pos_backend_http_request_duration_seconds_count{method="GET",path="/api/v1/sales",status="200"} 1234
  ```

#### `pos_backend_http_request_size_bytes`
- **Type:** Histogram
- **Labels:** `method`, `path`
- **Buckets:** Exponential from 100 bytes to ~100MB
- **Description:** HTTP request body size in bytes

#### `pos_backend_http_response_size_bytes`
- **Type:** Histogram
- **Labels:** `method`, `path`
- **Buckets:** Exponential from 100 bytes to ~100MB
- **Description:** HTTP response body size in bytes

#### `pos_backend_http_requests_in_flight`
- **Type:** Gauge
- **Description:** Current number of HTTP requests being processed
- **Example:**
  ```
  pos_backend_http_requests_in_flight 12
  ```

### Database Connection Pool Metrics

All database metrics are prefixed with `pos_backend_db_pool_`.

#### `pos_backend_db_pool_acquired_conns`
- **Type:** Gauge
- **Description:** Number of currently acquired connections from the pool
- **Example:** `pos_backend_db_pool_acquired_conns 8`

#### `pos_backend_db_pool_idle_conns`
- **Type:** Gauge
- **Description:** Number of idle connections in the pool
- **Example:** `pos_backend_db_pool_idle_conns 17`

#### `pos_backend_db_pool_total_conns`
- **Type:** Gauge
- **Description:** Total connections in the pool (acquired + idle + constructing)
- **Example:** `pos_backend_db_pool_total_conns 25`

#### `pos_backend_db_pool_max_conns`
- **Type:** Gauge
- **Description:** Maximum number of connections configured for the pool
- **Example:** `pos_backend_db_pool_max_conns 25`

#### `pos_backend_db_pool_constructing_conns`
- **Type:** Gauge
- **Description:** Number of connections currently being established
- **Example:** `pos_backend_db_pool_constructing_conns 0`

#### `pos_backend_db_pool_acquire_count`
- **Type:** Counter
- **Labels:** `status` (success/failure)
- **Description:** Cumulative count of connection acquire operations
- **Example:**
  ```
  pos_backend_db_pool_acquire_count{status="success"} 15234
  pos_backend_db_pool_acquire_count{status="failure"} 3
  ```

#### `pos_backend_db_pool_acquire_duration_seconds`
- **Type:** Histogram
- **Buckets:** 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1 second
- **Description:** Time taken to acquire a connection from the pool
- **Example:**
  ```
  pos_backend_db_pool_acquire_duration_seconds_bucket{le="0.01"} 14523
  pos_backend_db_pool_acquire_duration_seconds_sum 234.5
  pos_backend_db_pool_acquire_duration_seconds_count 15234
  ```

#### `pos_backend_db_pool_canceled_acquires_total`
- **Type:** Counter
- **Description:** Number of connection acquires canceled by context
- **Example:** `pos_backend_db_pool_canceled_acquires_total 12`

#### `pos_backend_db_pool_empty_acquires_total`
- **Type:** Counter
- **Description:** Number of acquires that had to wait for a connection
- **Example:** `pos_backend_db_pool_empty_acquires_total 456`

#### `pos_backend_db_pool_max_lifetime_destroy_total`
- **Type:** Counter
- **Description:** Connections destroyed for exceeding max lifetime
- **Example:** `pos_backend_db_pool_max_lifetime_destroy_total 89`

#### `pos_backend_db_pool_max_idle_destroy_total`
- **Type:** Counter
- **Description:** Connections destroyed for exceeding max idle time
- **Example:** `pos_backend_db_pool_max_idle_destroy_total 234`

### Business Metrics (Future)

Additional business metrics are available via `internal/metrics/collector.go`:
- Sales transactions and amounts
- Authentication attempts and failures
- Posting operations
- Errors by type

These will be instrumented in domain services in future phases.

## Configuration

Metrics are configured via environment variables:

```bash
METRICS_ENABLED=true
METRICS_PORT=9091  # Default port for metrics endpoint
```

## Grafana Dashboard

### Sample Queries

**Request Rate:**
```promql
rate(pos_backend_http_requests_total[5m])
```

**Request Latency (p95):**
```promql
histogram_quantile(0.95, rate(pos_backend_http_request_duration_seconds_bucket[5m]))
```

**Error Rate:**
```promql
rate(pos_backend_http_requests_total{status=~"5.."}[5m])
```

**Database Connection Pool Usage:**
```promql
pos_backend_db_pool_acquired_conns / pos_backend_db_pool_max_conns * 100
```

**Connection Wait Time (p99):**
```promql
histogram_quantile(0.99, rate(pos_backend_db_pool_acquire_duration_seconds_bucket[5m]))
```

### Alerts

**High Error Rate:**
```yaml
- alert: HighErrorRate
  expr: rate(pos_backend_http_requests_total{status=~"5.."}[5m]) > 0.05
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "High error rate detected"
    description: "Error rate is {{ $value }} errors/sec"
```

**Database Pool Exhaustion:**
```yaml
- alert: DatabasePoolExhausted
  expr: pos_backend_db_pool_acquired_conns / pos_backend_db_pool_max_conns > 0.9
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Database connection pool nearly exhausted"
    description: "Pool usage is {{ $value | humanizePercentage }}"
```

**Slow Requests:**
```yaml
- alert: SlowRequests
  expr: histogram_quantile(0.95, rate(pos_backend_http_request_duration_seconds_bucket[5m])) > 1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "p95 request latency exceeds 1 second"
    description: "p95 latency is {{ $value }}s"
```

## Prometheus Configuration

Add the following scrape configuration to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'pos-backend'
    scrape_interval: 15s
    static_configs:
      - targets: ['localhost:8080']  # Replace with your API server address
    metrics_path: '/metrics'
```

## Deployment Considerations

### Security

1. **Network Isolation:** The `/metrics` endpoint should not be publicly accessible
2. **Firewall Rules:** Restrict access to Prometheus server IPs only
3. **Reverse Proxy:** Use nginx/traefik to add authentication if needed

Example nginx config:
```nginx
location /metrics {
    allow 10.0.0.0/8;  # Prometheus server network
    deny all;
    proxy_pass http://backend:8080/metrics;
}
```

### Performance

- Metrics collection adds minimal overhead (~1-2ms per request)
- Database pool metrics collected every 15 seconds (configurable)
- High-cardinality labels avoided (no user IDs, transaction IDs in labels)

### Cardinality Management

To prevent cardinality explosion:
- URL paths use patterns (`/users/{id}`) not actual values
- Limited label values (method, status code)
- No unbounded label values (like email addresses)

## Testing

```bash
# Start the server
go run cmd/api/main.go

# Fetch metrics
curl http://localhost:8080/metrics

# Verify specific metric
curl -s http://localhost:8080/metrics | grep pos_backend_http_requests_total
```

## Architecture

```
┌─────────────────┐
│  HTTP Request   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│  Metrics Middleware     │◄─── Records HTTP metrics
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  Business Logic         │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  Database Layer         │◄─── DB pool metrics collected
└─────────────────────────┘
         │
         ▼
┌─────────────────────────┐
│  /metrics Endpoint      │──►  Prometheus scrapes here
└─────────────────────────┘
```

## Files

- `internal/middleware/metrics.go` - HTTP request metrics middleware
- `internal/metrics/collector.go` - Business metrics collector
- `internal/metrics/db_collector.go` - Database pool metrics collector
- `cmd/api/main.go` - Metrics initialization

## Next Steps

1. **Grafana Dashboard:** Create comprehensive dashboard with all metrics
2. **Business Metrics:** Instrument domain services with business metrics
3. **Alerting:** Configure Prometheus alerting rules
4. **Distributed Tracing:** Add OpenTelemetry for request tracing
5. **Custom Metrics:** Add application-specific metrics as needed

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [pgxpool Stats](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#Stat)

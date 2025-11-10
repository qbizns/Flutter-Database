# Production Readiness Plan - 3 Phases

**Goal:** Achieve 100% acceptance criteria compliance
**Current Status:** 4/7 Mandatory (57%) → Target: 7/7 Mandatory (100%)
**Total Estimated Time:** 74-119 hours → **Organized into 3 phases**

---

## Phase 3A: Critical Blockers (PRIORITY 1)
**Duration:** 24-36 hours
**Focus:** Remove production deployment blockers

### 🎯 Objectives
1. ✅ Enable clean database deployments
2. ✅ Automate quality checks
3. ✅ Fix immediate test failures

### Tasks

#### Task 3A.1: Database Migrations Setup (16-24 hours)
**Status:** ✅ **MIGRATIONS FOUND!**
- **Location:** `/postgres/migrations/` (23 files) + `/accounting/migrations/` (13 files)
- **Total:** 36 migration files already exist

**Sub-tasks:**
- [x] ~~Create migration structure~~ → Already exists!
- [ ] Link migrations to backend project (1 hour)
- [ ] Install migration tool (golang-migrate) (1 hour)
- [ ] Create migration runner scripts (2-3 hours)
- [ ] Test migrations on clean database (2-3 hours)
- [ ] Create rollback procedures (2-3 hours)
- [ ] Document migration usage (2-3 hours)
- [ ] Update Docker Compose with migration paths (1 hour)
- [ ] Create Makefile migration commands (1 hour)
- [ ] Verify all 36 migrations execute successfully (2-4 hours)

**Acceptance Criteria:**
- ✅ All 36 migrations execute without errors
- ✅ Clean database → Full schema in < 5 minutes
- ✅ Documented rollback process
- ✅ Single command execution (`make migrate-up`)

**Estimated:** 14-20 hours (reduced from 16-24 due to existing migrations)

---

#### Task 3A.2: Fix Broken Tests (4-6 hours)
**Status:** ❌ Critical - Blocking test execution

**Issues to Fix:**
1. **HTTP Handler Tests** - `category_handlers_test.go`:
   ```go
   // ERRORS:
   - logging.NewLogger(cfg) → logging.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
   - testDB.Context undefined → Add Context() method to TestDB
   - IconName/DisplayOrder → Icon/SortOrder
   - chi.NewRouteContext API incorrect
   ```

2. **Apply to Other Tests:**
   - customer_handlers_test.go
   - location_handlers_test.go
   - product_handlers_test.go
   - sale_handlers_test.go

**Sub-tasks:**
- [ ] Fix logging.NewLogger calls in all test files (1 hour)
- [ ] Add TestDB.Context() method (30 min)
- [ ] Fix field name mismatches (30 min)
- [ ] Fix chi router context usage (1 hour)
- [ ] Run tests and verify they execute (1-2 hours)
- [ ] Fix any remaining compilation errors (1-2 hours)

**Acceptance Criteria:**
- ✅ All tests compile without errors
- ✅ Tests execute (pass/fail acceptable at this stage)
- ✅ go test ./... completes without build failures

**Estimated:** 4-6 hours

---

#### Task 3A.3: CI/CD Pipeline Setup (4-6 hours)
**Status:** ❌ Critical - No automation

**Implementation:**
Create `.github/workflows/ci.yml`:

```yaml
name: Backend CI

on:
  push:
    branches: [ main, develop, claude/* ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    name: Test & Build
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: pos_test
          POSTGRES_PASSWORD: test_password
          POSTGRES_DB: pos_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

      redis:
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          cache: true

      - name: Install golang-migrate
        run: |
          curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
          sudo mv migrate /usr/local/bin/

      - name: Run Migrations
        working-directory: backend
        run: |
          make migrate-up
        env:
          DATABASE_URL: postgres://pos_test:test_password@localhost:5432/pos_test?sslmode=disable

      - name: Download Dependencies
        working-directory: backend
        run: go mod download

      - name: Build
        working-directory: backend
        run: go build -v ./...

      - name: Run Tests
        working-directory: backend
        run: go test -v -race -coverprofile=coverage.out ./...
        env:
          DATABASE_HOST: localhost
          DATABASE_PORT: 5432
          DATABASE_USER: pos_test
          DATABASE_PASSWORD: test_password
          DATABASE_NAME: pos_test
          DATABASE_SSL_MODE: disable
          REDIS_ADDR: localhost:6379

      - name: Coverage Check
        working-directory: backend
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Total coverage: $coverage%"
          if (( $(echo "$coverage < 50" | bc -l) )); then
            echo "❌ Coverage $coverage% is below 50% threshold"
            exit 1
          fi
          echo "✅ Coverage $coverage% meets 50% threshold"

      - name: Lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          working-directory: backend
          args: --timeout=5m

      - name: Security Scan
        uses: securego/gosec@master
        with:
          args: '-exclude-dir=testhelpers ./...'
          working-directory: backend

      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          files: backend/coverage.out
          flags: unittests

  docker:
    name: Docker Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build Docker Image
        run: docker-compose -f docker-compose.yml build backend

      - name: Test Docker Compose
        run: |
          docker-compose -f docker-compose.yml up -d
          sleep 10
          docker-compose -f docker-compose.yml ps
          docker-compose -f docker-compose.yml down
```

**Sub-tasks:**
- [ ] Create `.github/workflows/` directory (5 min)
- [ ] Write CI workflow file (1-2 hours)
- [ ] Configure quality gates (30 min)
- [ ] Test CI locally with act (optional) (1 hour)
- [ ] Push and verify CI runs (30 min)
- [ ] Add status badge to README (15 min)
- [ ] Fix any CI failures (1-2 hours)

**Acceptance Criteria:**
- ✅ CI runs on every push/PR
- ✅ Build step passes
- ✅ Test step executes
- ✅ Linting enforced
- ✅ Security scanning active
- ✅ Coverage check enforced (50% threshold)

**Estimated:** 4-6 hours

---

### Phase 3A Summary
**Total Time:** 22-32 hours
**Deliverables:**
- ✅ Working database migrations (36 files)
- ✅ Tests compile and run
- ✅ Automated CI/CD pipeline
- ✅ Quality gates enforced

**Acceptance Items Completed:**
- Item 6: Database Migrations ✅
- Item 10: CI Pipeline ✅
- Item 1: Testing (compilation fixed, ready for Phase 3B)

---

## Phase 3B: Testing & Quality (PRIORITY 2)
**Duration:** 38-54 hours
**Focus:** Achieve 50%+ test coverage and ensure data integrity

### 🎯 Objectives
1. ✅ Achieve 50% line coverage
2. ✅ Verify transaction safety
3. ✅ Test critical business logic

### Tasks

#### Task 3B.1: Domain Service Tests (16-24 hours)

**Critical Domains to Test:**

**1. Sales Domain (6-8 hours):**
```go
// Test files to create:
- sales/service_test.go
- sales/calculations_test.go
- sales/validation_test.go

// Test scenarios:
- CreateSale with multiple line items
- Calculate totals (subtotal, tax, discount, total)
- Apply discounts (percentage, fixed, item-level)
- Tax calculations (inclusive, exclusive, multiple rates)
- Inventory deduction on sale completion
- Refund processing
- Split payments
```

**2. Inventory Domain (4-6 hours):**
```go
// Test files:
- inventory/service_test.go
- inventory/costing_test.go
- inventory/transfers_test.go

// Test scenarios:
- Stock adjustments (increase, decrease)
- Transfer between locations
- FIFO/LIFO cost calculations
- Low stock alerts
- Batch/serial tracking
```

**3. Accounting/Posting Domain (6-10 hours):**
```go
// Test files:
- accounting/service_test.go
- posting/engine_test.go
- accounting/validation_test.go

// Test scenarios:
- Journal entry creation
- Double-entry validation (debits = credits)
- Posting to general ledger
- Period closing
- Trial balance calculation
- Chart of accounts operations
```

**Sub-tasks:**
- [ ] Set up test database helpers (2 hours)
- [ ] Write sales domain tests (6-8 hours)
- [ ] Write inventory domain tests (4-6 hours)
- [ ] Write accounting domain tests (6-10 hours)
- [ ] Run tests and fix failures (2-4 hours)

**Acceptance Criteria:**
- ✅ All critical domains have test files
- ✅ Tests cover happy paths and error cases
- ✅ Domain coverage > 60%

**Estimated:** 16-24 hours

---

#### Task 3B.2: Repository Integration Tests (12-16 hours)

**Repositories to Test:**
1. **Product Repository** (3-4 hours)
2. **Sale Repository** (3-4 hours)
3. **Inventory Repository** (3-4 hours)
4. **Accounting Repository** (3-4 hours)

**Test Template:**
```go
func TestProductRepository_CRUD(t *testing.T) {
    testDB := testhelpers.SetupTestDB(t)
    repo := postgres.NewProductRepository(testDB.DB)

    // Test Create
    product := &products.Product{...}
    err := repo.Create(context.Background(), product)
    require.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, product.ID)

    // Test Get
    retrieved, err := repo.Get(context.Background(), product.OrganizationID, product.ID)
    require.NoError(t, err)
    assert.Equal(t, product.Name, retrieved.Name)

    // Test Update
    product.Name = "Updated"
    err = repo.Update(context.Background(), product)
    require.NoError(t, err)

    // Test List with filters
    products, err := repo.List(context.Background(), orgID, filters)
    require.NoError(t, err)
    assert.Greater(t, len(products), 0)

    // Test Delete (soft delete)
    err = repo.Delete(context.Background(), orgID, product.ID)
    require.NoError(t, err)
}
```

**Sub-tasks:**
- [ ] Set up Docker test database (1 hour)
- [ ] Create test data fixtures (2 hours)
- [ ] Write product repository tests (3-4 hours)
- [ ] Write sale repository tests (3-4 hours)
- [ ] Write inventory repository tests (3-4 hours)
- [ ] Write accounting repository tests (3-4 hours)

**Acceptance Criteria:**
- ✅ Integration tests use real PostgreSQL
- ✅ Tests run in Docker environment
- ✅ Repository coverage > 50%

**Estimated:** 12-16 hours

---

#### Task 3B.3: Transaction Safety Audit (10-14 hours)

**Critical Workflows to Audit:**

**1. Sales Transactions (3-4 hours):**
```go
// Verify these are wrapped in transactions:
- CreateSale + CreateSaleLines + DeductInventory
- ProcessPayment + UpdateSaleStatus + GenerateReceipt
- ProcessRefund + RestoreInventory + CreateCreditNote
```

**2. Inventory Transactions (3-4 hours):**
```go
// Verify:
- TransferStock: Deduct from source + Add to destination
- AdjustStock: Update quantity + Create cost layer + Update valuation
- ReceivePurchase: Add stock + Update cost + Create GL posting
```

**3. Accounting Transactions (4-6 hours):**
```go
// Verify:
- CreateJournalEntry + CreateJournalLines (already verified ✓)
- PostToGL: Update GL + Mark as posted
- PeriodClose: Lock period + Generate reports + Create closing entries
```

**Audit Process:**
1. Read service code
2. Identify multi-step write operations
3. Verify transaction wrapping
4. Add transactions where missing
5. Write tests to verify rollback on error

**Sub-tasks:**
- [ ] Audit sales service (3-4 hours)
- [ ] Audit inventory service (3-4 hours)
- [ ] Audit accounting service (2-3 hours)
- [ ] Add missing transactions (2-3 hours)
- [ ] Write transaction rollback tests (2-3 hours)

**Acceptance Criteria:**
- ✅ All multi-step operations use transactions
- ✅ Error handling triggers rollback
- ✅ Tests verify rollback behavior

**Estimated:** 10-14 hours

---

### Phase 3B Summary
**Total Time:** 38-54 hours
**Deliverables:**
- ✅ Domain service tests (3 critical domains)
- ✅ Repository integration tests (4 repositories)
- ✅ Transaction safety verified
- ✅ 50%+ test coverage achieved

**Acceptance Items Completed:**
- Item 1: Testing Baseline ✅ (50% coverage achieved)
- Item 7: Transaction Safety ✅ (audit complete)

---

## Phase 3C: Production Polish (PRIORITY 3)
**Duration:** 24-35 hours
**Focus:** Authorization, observability, and operational readiness

### 🎯 Objectives
1. ✅ Complete RBAC implementation
2. ✅ Add observability endpoints
3. ✅ Finalize documentation

### Tasks

#### Task 3C.1: RBAC Audit & Tests (12-17 hours)

**1. Document Role Model (2-3 hours):**
```markdown
# Role-Based Access Control (RBAC)

## Roles
- **admin**: Full system access (all operations)
- **manager**: Location management, reporting, user management
- **cashier**: POS operations, sales, refunds
- **inventory_clerk**: Inventory management, stock adjustments
- **accountant**: Financial reports, journal entries (read-only)
- **viewer**: Read-only access to all modules

## Permission Matrix
| Endpoint | Admin | Manager | Cashier | Inventory | Accountant | Viewer |
|----------|-------|---------|---------|-----------|------------|--------|
| POST /sales | ✓ | ✓ | ✓ | ✗ | ✗ | ✗ |
| GET /sales | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ |
| POST /products | ✓ | ✓ | ✗ | ✓ | ✗ | ✗ |
| DELETE /users | ✓ | ✓ | ✗ | ✗ | ✗ | ✗ |
| ...
```

**2. Audit Route Protection (4-6 hours):**
Review all HTTP handlers and apply appropriate middleware:
```go
// Example:
r.Route("/api/v1/sales", func(r chi.Router) {
    r.Use(auth.RequireRole("cashier")) // or "manager" or "admin"
    r.Post("/", CreateSaleHandler(db, logger))
    r.Get("/", ListSalesHandler(db, logger))

    r.Route("/{id}", func(r chi.Router) {
        r.Get("/", GetSaleHandler(db, logger))
        r.With(auth.RequireRole("manager")).Delete("/", DeleteSaleHandler(db, logger))
    })
})
```

**3. Write Authorization Tests (6-8 hours):**
```go
func TestSaleEndpoints_Authorization(t *testing.T) {
    tests := []struct {
        name       string
        method     string
        path       string
        role       string
        wantStatus int
    }{
        {"Admin can create sale", "POST", "/sales", "admin", 201},
        {"Cashier can create sale", "POST", "/sales", "cashier", 201},
        {"Viewer cannot create sale", "POST", "/sales", "viewer", 403},
        {"Unauthenticated denied", "POST", "/sales", "", 401},
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := createRequestWithRole(tt.method, tt.path, tt.role)
            resp := executeRequest(req)
            assert.Equal(t, tt.wantStatus, resp.StatusCode)
        })
    }
}
```

**Sub-tasks:**
- [ ] Document role model and permissions (2-3 hours)
- [ ] Audit all HTTP routes (2-3 hours)
- [ ] Apply RequireRole middleware (2-3 hours)
- [ ] Write authorization tests (6-8 hours)
- [ ] Verify 401/403 responses (1-2 hours)

**Acceptance Criteria:**
- ✅ Role model documented
- ✅ All sensitive endpoints protected
- ✅ Tests verify 401/403 responses
- ✅ Permission matrix complete

**Estimated:** 12-17 hours

---

#### Task 3C.2: Observability Implementation (8-12 hours)

**1. Prometheus Metrics (4-6 hours):**
```go
// metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HttpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests by method, path, and status",
        },
        []string{"method", "path", "status"},
    )

    HttpDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request latency",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    DbConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_connections_active",
            Help: "Active database connections",
        },
    )
)
```

**Middleware:**
```go
func MetricsMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

            next.ServeHTTP(ww, r)

            duration := time.Since(start).Seconds()
            metrics.HttpRequestsTotal.WithLabelValues(
                r.Method,
                r.URL.Path,
                strconv.Itoa(ww.Status()),
            ).Inc()

            metrics.HttpDuration.WithLabelValues(
                r.Method,
                r.URL.Path,
            ).Observe(duration)
        })
    }
}
```

**2. Health Endpoints (2-3 hours):**
```go
// GET /healthz - Liveness probe
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

// GET /readyz - Readiness probe
func ReadyzHandler(db *postgres.DB, redis *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        health := map[string]interface{}{
            "status": "ok",
            "checks": map[string]string{},
        }

        // Check database
        if err := db.Pool.Ping(ctx); err != nil {
            health["status"] = "degraded"
            health["checks"].(map[string]string)["database"] = "unhealthy"
            w.WriteHeader(http.StatusServiceUnavailable)
        } else {
            health["checks"].(map[string]string)["database"] = "healthy"
        }

        // Check Redis
        if _, err := redis.Ping(ctx).Result(); err != nil {
            health["status"] = "degraded"
            health["checks"].(map[string]string)["redis"] = "unhealthy"
            w.WriteHeader(http.StatusServiceUnavailable)
        } else {
            health["checks"].(map[string]string)["redis"] = "healthy"
        }

        if health["status"] == "ok" {
            w.WriteHeader(http.StatusOK)
        }

        json.NewEncoder(w).Encode(health)
    }
}

// GET /metrics - Prometheus metrics
func MetricsHandler() http.HandlerFunc {
    return promhttp.Handler().ServeHTTP
}
```

**3. Configure Routes (1 hour):**
```go
// In router setup:
r.Get("/healthz", HealthzHandler)
r.Get("/readyz", ReadyzHandler(db, redis))
r.Get("/metrics", MetricsHandler())
```

**Sub-tasks:**
- [ ] Create metrics package (1 hour)
- [ ] Implement metrics middleware (2-3 hours)
- [ ] Add database/Redis metrics (1-2 hours)
- [ ] Create health endpoints (2-3 hours)
- [ ] Add metrics endpoint (1 hour)
- [ ] Test observability endpoints (1-2 hours)
- [ ] Document Prometheus setup (1 hour)

**Acceptance Criteria:**
- ✅ /metrics endpoint returns Prometheus format
- ✅ HTTP metrics (requests, duration) collected
- ✅ /healthz returns 200 OK
- ✅ /readyz checks dependencies
- ✅ Database connection metrics

**Estimated:** 8-12 hours

---

#### Task 3C.3: Documentation & Deployment (4-6 hours)

**1. Create Deployment Guide (2-3 hours):**
```markdown
# Deployment Guide

## Prerequisites
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+
- Go 1.21+

## Environment Variables
[Complete list with descriptions]

## Database Setup
1. Run migrations: `make migrate-up`
2. Seed initial data: `make seed`

## Running the Application
1. Development: `make dev`
2. Production: `docker-compose up -d`

## Health Checks
- Liveness: curl http://localhost:8080/healthz
- Readiness: curl http://localhost:8080/readyz
- Metrics: curl http://localhost:8080/metrics
```

**2. Update README (1-2 hours):**
- Add badges (build status, coverage, go report)
- Quick start guide
- API documentation links
- Contributing guidelines

**3. Create Migration Guide (1 hour):**
```markdown
# Database Migrations

## Running Migrations
```bash
# Up
make migrate-up

# Down
make migrate-down

# Status
make migrate-status
```

## Creating New Migrations
```bash
make migrate-create name=add_feature
```

## Rollback
```bash
make migrate-down  # Rollback last migration
```
```

**Sub-tasks:**
- [ ] Write deployment guide (2-3 hours)
- [ ] Update README with badges and guides (1-2 hours)
- [ ] Document migration procedures (1 hour)
- [ ] Create environment variable reference (30 min)

**Acceptance Criteria:**
- ✅ Deployment guide complete
- ✅ README updated
- ✅ Migration documentation clear
- ✅ Environment variables documented

**Estimated:** 4-6 hours

---

### Phase 3C Summary
**Total Time:** 24-35 hours
**Deliverables:**
- ✅ RBAC fully implemented and tested
- ✅ Observability endpoints (metrics, health, readiness)
- ✅ Complete deployment documentation
- ✅ Operational readiness

**Acceptance Items Completed:**
- Item 9: Authorization/RBAC ✅
- Item 8: Observability ✅ (Recommended)

---

## Master Timeline

| Phase | Duration | Working Days | Calendar Days | Key Deliverables |
|-------|----------|--------------|---------------|------------------|
| **3A: Critical Blockers** | 22-32 hrs | 3-4 days | 4-5 days | Migrations, CI/CD, Tests Fixed |
| **3B: Testing & Quality** | 38-54 hrs | 5-7 days | 6-9 days | 50% Coverage, Transaction Safety |
| **3C: Production Polish** | 24-35 hrs | 3-4 days | 4-6 days | RBAC, Observability, Docs |
| **TOTAL** | **84-121 hrs** | **11-15 days** | **14-20 days** | **Full Production Readiness** |

*Assumes single developer, 8-hour days. Team collaboration can reduce timeline.*

---

## Success Criteria

### Phase 3A Complete When:
- [x] Migrations discovered (36 files found!)
- [ ] Migration tooling installed and tested
- [ ] All 36 migrations run successfully
- [ ] Tests compile and execute
- [ ] CI pipeline runs automatically
- [ ] Quality gates enforced

### Phase 3B Complete When:
- [ ] Domain tests written (Sales, Inventory, Accounting)
- [ ] Repository integration tests complete
- [ ] Transaction safety audit done
- [ ] Overall test coverage ≥ 50%
- [ ] All tests passing in CI

### Phase 3C Complete When:
- [ ] RBAC documented and tested
- [ ] Observability endpoints working
- [ ] Deployment guide complete
- [ ] All 10 acceptance criteria met
- [ ] Ready for production deployment

---

## Acceptance Criteria Mapping

| Item | Phase | Status After Completion |
|------|-------|------------------------|
| 1. Testing Baseline | 3A + 3B | ✅ PASS |
| 2. Rate Limiting | ✅ Done | ✅ PASS |
| 3. Context Safety | ✅ Done | ✅ PASS |
| 4. Config Parsing | ✅ Done | ✅ PASS |
| 5. Database SSL | ✅ Done | ✅ PASS |
| 6. Migrations | 3A | ✅ PASS |
| 7. Transactions | 3B | ✅ PASS |
| 8. Observability | 3C | ✅ PASS |
| 9. Authorization | 3C | ✅ PASS |
| 10. CI Pipeline | 3A | ✅ PASS |

**Final Result: 10/10 Mandatory + Recommended ✅**

---

## Payment Milestones

| Milestone | Deliverable | Value | Cumulative |
|-----------|-------------|-------|------------|
| Phase 2 Complete | Compilation Success | $4,000 | $4,000 |
| Phase 3A Complete | Critical Blockers Fixed | $2,000 | $6,000 |
| Phase 3B Complete | Testing & Quality | $2,500 | $8,500 |
| Phase 3C Complete | Production Ready | $1,500 | $10,000 |

---

## Risk Mitigation

### Phase 3A Risks:
- **Migration failures:** Migrations already exist, just need setup
- **CI configuration issues:** Use proven templates, test locally first
- **Test fixes more complex:** Isolated changes, one test file at a time

### Phase 3B Risks:
- **Low coverage achievement:** Focus on critical paths first
- **Transaction changes break code:** Comprehensive testing required
- **Integration tests flaky:** Use Docker, clean state between tests

### Phase 3C Risks:
- **RBAC changes require route refactoring:** Minimal middleware addition
- **Observability overhead:** Use proven libraries (Prometheus)
- **Documentation incomplete:** Use templates, focus on essentials

---

## Next Steps

**Immediate Action - Starting Phase 3A:**

1. **Task 3A.1.1:** Link migrations to backend (NEXT)
   ```bash
   cd backend
   ln -s ../postgres/migrations db/migrations-postgres
   ln -s ../accounting/migrations db/migrations-accounting
   ```

2. **Task 3A.1.2:** Install golang-migrate
   ```bash
   brew install golang-migrate  # Mac
   # or
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

3. **Task 3A.1.3:** Create Makefile commands
   ```makefile
   migrate-up:
       migrate -path db/migrations-postgres -database $(DATABASE_URL) up
       migrate -path db/migrations-accounting -database $(DATABASE_URL) up
   ```

**Let's begin! 🚀**

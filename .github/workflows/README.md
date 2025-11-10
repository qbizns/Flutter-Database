# CI/CD Pipeline Documentation

This directory contains GitHub Actions workflows for the POS Backend project.

## Workflows

### `ci.yml` - Main CI/CD Pipeline

Runs on every push to `main`, `develop`, and `claude/**` branches, and on pull requests.

#### Jobs

1. **Build & Test** (`build-and-test`)
   - Sets up Go 1.21
   - Starts PostgreSQL and Redis services
   - Runs database migrations
   - Builds all packages
   - Runs tests with race detection
   - Generates coverage report
   - Checks coverage threshold (50%)
   - Uploads coverage to Codecov (optional)

2. **Lint** (`lint`)
   - Runs golangci-lint with comprehensive checks
   - Configuration: `backend/.golangci.yml`
   - Checks code style, complexity, and common issues

3. **Security Scan** (`security`)
   - Runs Gosec security scanner
   - Uploads SARIF results to GitHub Security
   - Detects common security vulnerabilities

4. **Dependency Check** (`dependency-check`)
   - Runs govulncheck for known vulnerabilities
   - Checks Go module dependencies
   - Fails on critical vulnerabilities

5. **Migration Validation** (`migration-check`)
   - Validates migration file structure
   - Checks for duplicate versions
   - Verifies naming conventions
   - Runs `backend/scripts/validate-migrations.sh`

6. **Quality Gate** (`quality-gate`)
   - Final summary job
   - Requires all previous jobs to pass
   - Displays overall status

## Environment Variables

- `GO_VERSION`: Go version to use (default: `1.21`)
- `COVERAGE_THRESHOLD`: Minimum test coverage percentage (default: `50`)

## Services

### PostgreSQL
- Image: `postgres:15-alpine`
- Port: `5432`
- Database: `pos_test`
- User: `postgres`
- Password: `postgres`

### Redis
- Image: `redis:7-alpine`
- Port: `6379`

## Quality Gates

All jobs must pass for the quality gate to succeed:

| Check | Requirement |
|-------|-------------|
| Build | Must compile successfully |
| Tests | All tests must pass |
| Coverage | ≥50% (warning only for now) |
| Lint | No critical issues |
| Security | No high/critical vulnerabilities |
| Dependencies | No known vulnerabilities |
| Migrations | Valid structure and naming |

## Local Testing

### Run CI checks locally:

```bash
# Build
cd backend
go build ./...

# Tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Lint (requires golangci-lint)
golangci-lint run ./...

# Security scan (requires gosec)
gosec -no-fail ./...

# Dependency check (requires govulncheck)
govulncheck ./...

# Migration validation
./scripts/validate-migrations.sh
```

### Install required tools:

```bash
# golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Badges (Add to main README)

```markdown
![CI Status](https://github.com/your-org/pos-backend/workflows/Backend%20CI%2FCD%20Pipeline/badge.svg)
![Coverage](https://codecov.io/gh/your-org/pos-backend/branch/main/graph/badge.svg)
```

## Troubleshooting

### Coverage below threshold
- Currently showing as warning only
- Will be enforced once we reach 50%
- Add tests to increase coverage

### Lint failures
- Run `golangci-lint run ./...` locally
- Check `.golangci.yml` for configuration
- Some checks are stricter in CI

### Migration validation failures
- Run `./scripts/validate-migrations.sh` locally
- Check migration file naming: `VXXX_YYYYMMDD_description.sql`
- Ensure no duplicate version numbers

### Security scan failures
- Review Gosec findings in GitHub Security tab
- Some warnings can be suppressed in `.golangci.yml`
- Address critical/high severity issues immediately

## Future Enhancements

- [ ] Add deployment job for staging/production
- [ ] Add integration test job with full database
- [ ] Add performance benchmarking
- [ ] Add Docker image building and publishing
- [ ] Add automatic changelog generation
- [ ] Add Slack/Discord notifications
- [ ] Add manual approval for production deploys

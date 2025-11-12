# API Documentation

Complete documentation for the POS Backend API.

## Quick Start

### View API Documentation

1. **OpenAPI/Swagger Specification**
   - File: `docs/swagger.yaml`
   - Format: OpenAPI 3.0.3

2. **View in Swagger UI** (Recommended)
   ```bash
   # Using Docker
   docker run -p 8081:8080 -e SWAGGER_JSON=/swagger/swagger.yaml \
     -v $(pwd)/docs:/swagger swaggerapi/swagger-ui

   # Open browser at: http://localhost:8081
   ```

3. **View in Swagger Editor**
   ```bash
   # Using Docker
   docker run -p 8082:8080 swaggerapi/swagger-editor

   # Open browser at: http://localhost:8082
   # Then: File > Import File > Select swagger.yaml
   ```

4. **Generate Client SDKs**
   ```bash
   # Install OpenAPI Generator
   npm install @openapitools/openapi-generator-cli -g

   # Generate TypeScript/JavaScript client
   openapi-generator-cli generate -i docs/swagger.yaml \
     -g typescript-axios -o clients/typescript

   # Generate Python client
   openapi-generator-cli generate -i docs/swagger.yaml \
     -g python -o clients/python

   # Generate Go client
   openapi-generator-cli generate -i docs/swagger.yaml \
     -g go -o clients/go
   ```

---

## API Overview

### Base URLs

| Environment | URL |
|-------------|-----|
| Local Development | `http://localhost:8080/api/v1` |
| Staging | `https://api-staging.example.com/api/v1` |
| Production | `https://api.example.com/api/v1` |

### Authentication

All API endpoints (except `/health` and `/auth/*`) require Bearer token authentication.

#### Get a Token

**1. Register a new account:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "name": "John Doe",
    "organization_name": "My Company Inc"
  }'
```

Response:
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "organization_id": "660e8400-e29b-41d4-a716-446655440001",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**2. Login to existing account:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "organization_id": "660e8400-e29b-41d4-a716-446655440001",
  "expires_at": "2025-11-13T10:00:00Z"
}
```

**3. Use the token in requests:**
```bash
curl -X GET http://localhost:8080/api/v1/organizations/{org_id}/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Common API Patterns

### Pagination

Most list endpoints support pagination:

```bash
GET /api/v1/organizations/{org_id}/products?page=1&limit=20
```

Response includes pagination metadata:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_pages": 5,
    "total_items": 97,
    "has_next": true,
    "has_prev": false
  }
}
```

### Filtering and Search

Many endpoints support filtering:

```bash
# Search products
GET /api/v1/organizations/{org_id}/products?search=laptop

# Filter by category
GET /api/v1/organizations/{org_id}/products?category=electronics

# Multiple filters
GET /api/v1/organizations/{org_id}/sales?status=completed&payment_method=CASH
```

### Sorting

Use the `sort` parameter:

```bash
# Sort by created date (descending)
GET /api/v1/organizations/{org_id}/sales?sort=-created_at

# Sort by name (ascending)
GET /api/v1/organizations/{org_id}/products?sort=name

# Multiple sort fields
GET /api/v1/organizations/{org_id}/invoices?sort=-due_date,status
```

---

## API Endpoints by Category

### 1. Health & Status
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics (if enabled)

### 2. Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - User login
- `POST /auth/logout` - Logout (invalidate token)
- `POST /auth/refresh` - Refresh token

### 3. Organizations
- `GET /organizations/{org_id}` - Get organization details
- `PUT /organizations/{org_id}` - Update organization
- `GET /organizations/{org_id}/users` - List organization users
- `POST /organizations/{org_id}/users` - Invite user

### 4. Products & Inventory
- `GET /organizations/{org_id}/products` - List products
- `POST /organizations/{org_id}/products` - Create product
- `GET /organizations/{org_id}/products/{id}` - Get product
- `PUT /organizations/{org_id}/products/{id}` - Update product
- `DELETE /organizations/{org_id}/products/{id}` - Delete product
- `GET /organizations/{org_id}/inventory/movements` - Inventory movements

### 5. Customers
- `GET /organizations/{org_id}/customers` - List customers
- `POST /organizations/{org_id}/customers` - Create customer
- `GET /organizations/{org_id}/customers/{id}` - Get customer
- `PUT /organizations/{org_id}/customers/{id}` - Update customer
- `DELETE /organizations/{org_id}/customers/{id}` - Delete customer

### 6. Sales
- `GET /organizations/{org_id}/sales` - List sales
- `POST /organizations/{org_id}/sales` - Create sale
- `GET /organizations/{org_id}/sales/{id}` - Get sale
- `PUT /organizations/{org_id}/sales/{id}` - Update sale
- `DELETE /organizations/{org_id}/sales/{id}` - Delete sale
- `POST /organizations/{org_id}/sales/{id}/return` - Process return

### 7. Invoices
- `GET /organizations/{org_id}/invoices` - List invoices
- `POST /organizations/{org_id}/invoices` - Create invoice
- `GET /organizations/{org_id}/invoices/{id}` - Get invoice
- `PUT /organizations/{org_id}/invoices/{id}` - Update invoice
- `POST /organizations/{org_id}/invoices/{id}/send` - Send invoice
- `POST /organizations/{org_id}/invoices/{id}/payments` - Record payment

### 8. Accounting
- `GET /organizations/{org_id}/accounts` - List chart of accounts
- `POST /organizations/{org_id}/accounts` - Create account
- `GET /organizations/{org_id}/journal-entries` - List journal entries
- `POST /organizations/{org_id}/journal-entries` - Create journal entry
- `GET /organizations/{org_id}/general-ledger` - Query general ledger

### 9. Posting Engine
- `POST /organizations/{org_id}/posting/post` - Execute posting
- `POST /organizations/{org_id}/posting/batch` - Batch posting
- `GET /organizations/{org_id}/posting/configurations` - List configurations
- `POST /organizations/{org_id}/posting/configurations` - Create configuration

### 10. Reports
- `GET /organizations/{org_id}/reports/trial-balance` - Trial balance
- `GET /organizations/{org_id}/reports/profit-loss` - P&L statement
- `GET /organizations/{org_id}/reports/balance-sheet` - Balance sheet
- `GET /organizations/{org_id}/reports/cash-flow` - Cash flow statement
- `GET /organizations/{org_id}/reports/sales-summary` - Sales summary

### 11. Background Jobs
- `POST /organizations/{org_id}/jobs/export` - Enqueue data export
- `POST /organizations/{org_id}/jobs/report` - Enqueue report generation
- `GET /organizations/{org_id}/jobs/{id}` - Get job status

---

## Example Workflows

### Complete Sales Transaction

```bash
# 1. Create a sale
curl -X POST http://localhost:8080/api/v1/organizations/{org_id}/sales \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "...",
    "sale_date": "2025-11-12T10:00:00Z",
    "payment_method": "CASH",
    "items": [
      {
        "product_id": "...",
        "quantity": 2,
        "unit_price": 49.99,
        "discount": 0
      }
    ]
  }'

# 2. Post to accounting (automatic via posting engine)
# This happens automatically in the background

# 3. Generate invoice (if needed)
curl -X POST http://localhost:8080/api/v1/organizations/{org_id}/invoices \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "sale_id": "...",
    "send_email": true
  }'
```

### Generate Financial Report

```bash
# 1. Request trial balance report
curl -X GET "http://localhost:8080/api/v1/organizations/{org_id}/reports/trial-balance?start_date=2025-01-01&end_date=2025-11-12" \
  -H "Authorization: Bearer {token}"

# 2. Async report generation (for large reports)
curl -X POST http://localhost:8080/api/v1/organizations/{org_id}/jobs/report \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "report_type": "profit_loss",
    "start_date": "2025-01-01",
    "end_date": "2025-12-31",
    "format": "pdf"
  }'

# Response
{
  "job_id": "...",
  "status": "queued"
}

# 3. Check job status
curl -X GET http://localhost:8080/api/v1/organizations/{org_id}/jobs/{job_id} \
  -H "Authorization: Bearer {token}"
```

---

## Error Handling

### Error Response Format

All errors return a consistent JSON format:

```json
{
  "error": "Human-readable error message",
  "code": "MACHINE_READABLE_CODE",
  "details": {
    "field": "additional_info"
  }
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| `BAD_REQUEST` | 400 | Invalid request format or parameters |
| `UNAUTHORIZED` | 401 | Missing or invalid authentication token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 422 | Business rule violation |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_SERVER_ERROR` | 500 | Server error |

### Example Error Responses

**Validation Error:**
```json
{
  "error": "Validation failed",
  "code": "VALIDATION_ERROR",
  "details": {
    "field": "email",
    "message": "Email already exists"
  }
}
```

**Rate Limit Exceeded:**
```json
{
  "error": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED",
  "details": {
    "retry_after": 60
  }
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

| Endpoint Type | Limit | Window |
|---------------|-------|--------|
| Auth endpoints | 10 requests | per minute |
| Standard endpoints | 100 requests | per minute |
| Burst allowance | 150 requests | per minute |

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1636732800
```

---

## Webhooks (Future)

The API will support webhooks for real-time event notifications:

- `sale.created`
- `invoice.paid`
- `inventory.low_stock`
- `posting.completed`
- `report.ready`

---

## Versioning

The API uses URL versioning: `/api/v1/...`

Breaking changes will be introduced in new versions (`v2`, `v3`, etc.) while maintaining backward compatibility for previous versions for at least 12 months.

---

## Support

- API Documentation: `/docs/swagger.yaml`
- Metrics Endpoint: `/metrics` (Prometheus format)
- Health Check: `/health`
- Support Email: support@example.com

---

## Testing

See `tests/load/api_load_test.js` for k6 load testing examples.

```bash
# Run load test
k6 run tests/load/api_load_test.js

# Run with custom options
k6 run --vus 50 --duration 5m tests/load/api_load_test.js
```

---

**Last Updated**: 2025-11-12
**API Version**: 1.0.0

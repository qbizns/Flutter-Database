# 📡 Complete API Endpoints Reference

**Backend:** Flutter-Database POS System
**Total APIs:** 171 tables × 5 endpoints = **855+ REST endpoints**
**Base URL:** `http://localhost:8080`
**Authentication:** Bearer Token (JWT)

---

## 🔐 Authentication

All endpoints require authentication except `/health` and `/metrics`.

**Headers Required:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

---

## 📋 Endpoint Pattern

Each table follows this pattern (example: `customers`):

```
POST   /api/v1/organizations/{orgID}/customers           # Create
GET    /api/v1/organizations/{orgID}/customers           # List (paginated)
GET    /api/v1/organizations/{orgID}/customers/{id}      # Get by ID
PUT    /api/v1/organizations/{orgID}/customers/{id}      # Update
DELETE /api/v1/organizations/{orgID}/customers/{id}      # Delete (soft)
```

**Query Parameters for List:**
- `page` (int, default: 1)
- `limit` (int, default: 20, max: 100)

---

## 📚 Complete Table List (Alphabetical)

### A-C

| # | Table | Endpoint Base Path |
|---|-------|-------------------|
| 1 | account_subtype | `/api/v1/organizations/{orgID}/account_subtype` |
| 2 | account_type | `/api/v1/organizations/{orgID}/account_type` |
| 3 | accounting_period | `/api/v1/organizations/{orgID}/accounting_period` |
| 4 | analytic_account | `/api/v1/organizations/{orgID}/analytic_account` |
| 5 | analytic_plan | `/api/v1/organizations/{orgID}/analytic_plan` |
| 6 | api_key | `/api/v1/organizations/{orgID}/api_key` |
| 7 | api_request_log | `/api/v1/organizations/{orgID}/api_request_log` |
| 8 | asset_category | `/api/v1/organizations/{orgID}/asset_category` |
| 9 | asset_depreciation_schedule | `/api/v1/organizations/{orgID}/asset_depreciation_schedule` |
| 10 | audit_log | `/api/v1/organizations/{orgID}/audit_log` |
| 11 | background_job | `/api/v1/organizations/{orgID}/background_job` |
| 12 | bank_account | `/api/v1/organizations/{orgID}/bank_account` |
| 13 | bank_reconciliation | `/api/v1/organizations/{orgID}/bank_reconciliation` |
| 14 | bank_reconciliation_item | `/api/v1/organizations/{orgID}/bank_reconciliation_item` |
| 15 | bank_statement | `/api/v1/organizations/{orgID}/bank_statement` |
| 16 | bank_statement_line | `/api/v1/organizations/{orgID}/bank_statement_line` |
| 17 | bank_statement_reconciliation | `/api/v1/organizations/{orgID}/bank_statement_reconciliation` |
| 18 | batch_transaction | `/api/v1/organizations/{orgID}/batch_transaction` |
| 19 | budget | `/api/v1/organizations/{orgID}/budget` |
| 20 | budget_line | `/api/v1/organizations/{orgID}/budget_line` |
| 21 | cash_drawer | `/api/v1/organizations/{orgID}/cash_drawer` |
| 22 | cash_drawer_session | `/api/v1/organizations/{orgID}/cash_drawer_session` |
| 23 | cash_movement | `/api/v1/organizations/{orgID}/cash_movement` |
| 24 | category | `/api/v1/organizations/{orgID}/category` |
| 25 | chart_of_account | `/api/v1/organizations/{orgID}/chart_of_account` |
| 26 | cours | `/api/v1/organizations/{orgID}/cours` |
| 27 | currency | `/api/v1/organizations/{orgID}/currency` |
| 28 | currency_rate | `/api/v1/organizations/{orgID}/currency_rate` |
| 29 | customer | `/api/v1/organizations/{orgID}/customer` |
| 30 | customer_address | `/api/v1/organizations/{orgID}/customer_address` |
| 31 | customer_invoice | `/api/v1/organizations/{orgID}/customer_invoice` |
| 32 | customer_invoice_line | `/api/v1/organizations/{orgID}/customer_invoice_line` |
| 33 | customer_payment | `/api/v1/organizations/{orgID}/customer_payment` |
| 34 | customer_payment_application | `/api/v1/organizations/{orgID}/customer_payment_application` |
| 35 | customer_store_credit_account | `/api/v1/organizations/{orgID}/customer_store_credit_account` |
| 36 | customer_tier_history | `/api/v1/organizations/{orgID}/customer_tier_history` |
| 37 | cycle_count | `/api/v1/organizations/{orgID}/cycle_count` |
| 38 | cycle_count_item | `/api/v1/organizations/{orgID}/cycle_count_item` |

### D-G

| # | Table | Endpoint Base Path |
|---|-------|-------------------|
| 39 | data_export_request | `/api/v1/organizations/{orgID}/data_export_request` |
| 40 | deferred_expense_contract | `/api/v1/organizations/{orgID}/deferred_expense_contract` |
| 41 | deferred_expense_schedule | `/api/v1/organizations/{orgID}/deferred_expense_schedule` |
| 42 | deferred_revenue_contract | `/api/v1/organizations/{orgID}/deferred_revenue_contract` |
| 43 | deferred_revenue_schedule | `/api/v1/organizations/{orgID}/deferred_revenue_schedule` |
| 44 | delivery_assignment | `/api/v1/organizations/{orgID}/delivery_assignment` |
| 45 | delivery_driver | `/api/v1/organizations/{orgID}/delivery_driver` |
| 46 | delivery_zone | `/api/v1/organizations/{orgID}/delivery_zone` |
| 47 | device | `/api/v1/organizations/{orgID}/device` |
| 48 | document_sequence | `/api/v1/organizations/{orgID}/document_sequence` |
| 49 | driver_shift | `/api/v1/organizations/{orgID}/driver_shift` |
| 50 | e_invoicing_document | `/api/v1/organizations/{orgID}/e_invoicing_document` |
| 51 | e_invoicing_document_event | `/api/v1/organizations/{orgID}/e_invoicing_document_event` |
| 52 | email_queue | `/api/v1/organizations/{orgID}/email_queue` |
| 53 | employee_schedule | `/api/v1/organizations/{orgID}/employee_schedule` |
| 54 | expens | `/api/v1/organizations/{orgID}/expens` |
| 55 | file_attachment | `/api/v1/organizations/{orgID}/file_attachment` |
| 56 | fiscal_year | `/api/v1/organizations/{orgID}/fiscal_year` |
| 57 | fixed_asset | `/api/v1/organizations/{orgID}/fixed_asset` |
| 58 | general_ledger | `/api/v1/organizations/{orgID}/general_ledger` |
| 59 | general_ledger_account | `/api/v1/organizations/{orgID}/general_ledger_account` |
| 60 | gift_card | `/api/v1/organizations/{orgID}/gift_card` |
| 61 | gift_card_transaction | `/api/v1/organizations/{orgID}/gift_card_transaction` |
| 62 | goods_receipt | `/api/v1/organizations/{orgID}/goods_receipt` |
| 63 | goods_receipt_item | `/api/v1/organizations/{orgID}/goods_receipt_item` |

### I-O

| # | Table | Endpoint Base Path |
|---|-------|-------------------|
| 64 | inventory_adjustment | `/api/v1/organizations/{orgID}/inventory_adjustment` |
| 65 | inventory_transaction | `/api/v1/organizations/{orgID}/inventory_transaction` |
| 66 | inventory_transfer | `/api/v1/organizations/{orgID}/inventory_transfer` |
| 67 | inventory_transfer_item | `/api/v1/organizations/{orgID}/inventory_transfer_item` |
| 68 | inventory_valuation | `/api/v1/organizations/{orgID}/inventory_valuation` |
| 69 | journal_entry | `/api/v1/organizations/{orgID}/journal_entry` |
| 70 | journal_entry_line | `/api/v1/organizations/{orgID}/journal_entry_line` |
| 71 | kitchen_order | `/api/v1/organizations/{orgID}/kitchen_order` |
| 72 | kitchen_station | `/api/v1/organizations/{orgID}/kitchen_station` |
| 73 | location | `/api/v1/organizations/{orgID}/location` |
| 74 | loyalty_point | `/api/v1/organizations/{orgID}/loyalty_point` |
| 75 | loyalty_program | `/api/v1/organizations/{orgID}/loyalty_program` |
| 76 | loyalty_reward | `/api/v1/organizations/{orgID}/loyalty_reward` |
| 77 | loyalty_tier | `/api/v1/organizations/{orgID}/loyalty_tier` |
| 78 | loyalty_tier_benefit | `/api/v1/organizations/{orgID}/loyalty_tier_benefit` |
| 79 | modifier | `/api/v1/organizations/{orgID}/modifier` |
| 80 | modifier_group | `/api/v1/organizations/{orgID}/modifier_group` |
| 81 | notification_preference | `/api/v1/organizations/{orgID}/notification_preference` |
| 82 | order | `/api/v1/organizations/{orgID}/order` |
| 83 | order_item | `/api/v1/organizations/{orgID}/order_item` |
| 84 | order_item_modifier | `/api/v1/organizations/{orgID}/order_item_modifier` |
| 85 | organization | `/api/v1/organizations/{orgID}/organization` |
| 86 | organization_feature | `/api/v1/organizations/{orgID}/organization_feature` |

### P-R

| # | Table | Endpoint Base Path |
|---|-------|-------------------|
| 87 | payment_method | `/api/v1/organizations/{orgID}/payment_method` |
| 88 | payment_term | `/api/v1/organizations/{orgID}/payment_term` |
| 89 | period_lock | `/api/v1/organizations/{orgID}/period_lock` |
| 90 | permission | `/api/v1/organizations/{orgID}/permission` |
| 91 | pos_session | `/api/v1/organizations/{orgID}/pos_session` |
| 92 | posting_rule | `/api/v1/organizations/{orgID}/posting_rule` |
| 93 | posting_validation_rule | `/api/v1/organizations/{orgID}/posting_validation_rule` |
| 94 | price_list | `/api/v1/organizations/{orgID}/price_list` |
| 95 | price_list_item | `/api/v1/organizations/{orgID}/price_list_item` |
| 96 | product | `/api/v1/organizations/{orgID}/product` |
| 97 | product_barcode | `/api/v1/organizations/{orgID}/product_barcode` |
| 98 | product_image | `/api/v1/organizations/{orgID}/product_image` |
| 99 | product_modifier | `/api/v1/organizations/{orgID}/product_modifier` |
| 100 | product_modifier_group | `/api/v1/organizations/{orgID}/product_modifier_group` |
| 101 | product_serial_number | `/api/v1/organizations/{orgID}/product_serial_number` |
| 102 | product_variant | `/api/v1/organizations/{orgID}/product_variant` |
| 103 | promotion | `/api/v1/organizations/{orgID}/promotion` |
| 104 | promotion_usage | `/api/v1/organizations/{orgID}/promotion_usage` |
| 105 | purchase_order | `/api/v1/organizations/{orgID}/purchase_order` |
| 106 | purchase_order_item | `/api/v1/organizations/{orgID}/purchase_order_item` |
| 107 | receivable | `/api/v1/organizations/{orgID}/receivable` |
| 108 | reconciliation_rule | `/api/v1/organizations/{orgID}/reconciliation_rule` |
| 109 | reservation | `/api/v1/organizations/{orgID}/reservation` |
| 110 | return_reason | `/api/v1/organizations/{orgID}/return_reason` |
| 111 | role | `/api/v1/organizations/{orgID}/role` |
| 112 | role_permission | `/api/v1/organizations/{orgID}/role_permission` |

### S-Z

| # | Table | Endpoint Base Path |
|---|-------|-------------------|
| 113 | sale | `/api/v1/organizations/{orgID}/sale` |
| 114 | sale_item | `/api/v1/organizations/{orgID}/sale_item` |
| 115 | sale_payment | `/api/v1/organizations/{orgID}/sale_payment` |
| 116 | scheduled_report | `/api/v1/organizations/{orgID}/scheduled_report` |
| 117 | serial_number_tracking | `/api/v1/organizations/{orgID}/serial_number_tracking` |
| 118 | stock_level | `/api/v1/organizations/{orgID}/stock_level` |
| 119 | stock_movement | `/api/v1/organizations/{orgID}/stock_movement` |
| 120 | supplier | `/api/v1/organizations/{orgID}/supplier` |
| 121 | system_preference | `/api/v1/organizations/{orgID}/system_preference` |
| 122 | table | `/api/v1/organizations/{orgID}/table` |
| 123 | table_section | `/api/v1/organizations/{orgID}/table_section` |
| 124 | tax | `/api/v1/organizations/{orgID}/tax` |
| 125 | tax_group | `/api/v1/organizations/{orgID}/tax_group` |
| 126 | tip | `/api/v1/organizations/{orgID}/tip` |
| 127 | tip_pool | `/api/v1/organizations/{orgID}/tip_pool` |
| 128 | uom | `/api/v1/organizations/{orgID}/uom` |
| 129 | uom_conversion | `/api/v1/organizations/{orgID}/uom_conversion` |
| 130 | user | `/api/v1/organizations/{orgID}/user` |
| 131 | vendor_bill | `/api/v1/organizations/{orgID}/vendor_bill` |
| 132 | vendor_bill_line | `/api/v1/organizations/{orgID}/vendor_bill_line` |
| 133 | vendor_payment | `/api/v1/organizations/{orgID}/vendor_payment` |
| 134 | vendor_payment_application | `/api/v1/organizations/{orgID}/vendor_payment_application` |
| 135 | warehouse | `/api/v1/organizations/{orgID}/warehouse` |
| 136-171 | ... (36 more tables) | ... |

---

## 📝 Request/Response Examples

### Create Customer
```http
POST /api/v1/organizations/550e8400-e29b-41d4-a716-446655440000/customers
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "customer_code": "CUST001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "customer_type": "retail"
}
```

**Response:** `201 Created`
```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "customer_code": "CUST001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "customer_type": "retail",
  "created_at": "2025-11-12T10:30:00Z",
  "updated_at": "2025-11-12T10:30:00Z"
}
```

### List Customers (Paginated)
```http
GET /api/v1/organizations/550e8400-e29b-41d4-a716-446655440000/customers?page=1&limit=20
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "customer_code": "CUST001",
      "first_name": "John",
      "last_name": "Doe",
      ...
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### Get Customer by ID
```http
GET /api/v1/organizations/550e8400-e29b-41d4-a716-446655440000/customers/7c9e6679-7425-40de-944b-e07fc1f90ae7
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:** `200 OK`
```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "customer_code": "CUST001",
  ...
}
```

### Update Customer
```http
PUT /api/v1/organizations/550e8400-e29b-41d4-a716-446655440000/customers/7c9e6679-7425-40de-944b-e07fc1f90ae7
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "email": "john.doe.updated@example.com",
  "phone": "+0987654321"
}
```

**Response:** `200 OK`
```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "email": "john.doe.updated@example.com",
  "phone": "+0987654321",
  "updated_at": "2025-11-12T11:00:00Z",
  ...
}
```

### Delete Customer (Soft Delete)
```http
DELETE /api/v1/organizations/550e8400-e29b-41d4-a716-446655440000/customers/7c9e6679-7425-40de-944b-e07fc1f90ae7
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:** `204 No Content`

---

## 🔧 System Endpoints

### Health Check
```http
GET /health

Response: 200 OK
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "timestamp": "2025-11-12T10:30:00Z"
}
```

### Metrics (Prometheus)
```http
GET /metrics

Response: 200 OK (Prometheus format)
# HELP api_requests_total Total API requests
# TYPE api_requests_total counter
api_requests_total{method="GET",endpoint="/customers",status="200"} 1234
...
```

### Swagger Documentation
```http
GET /swagger/index.html

Response: 200 OK (Interactive API docs)
```

---

## 🚦 HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET, PUT |
| 201 | Created | Successful POST |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Invalid request body/parameters |
| 401 | Unauthorized | Missing/invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 422 | Unprocessable Entity | Validation failed |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |

---

## 🔒 Row-Level Security (RLS)

All multi-tenant endpoints automatically filter by `organization_id`. Users can only:
- See data from their organization
- Create data in their organization
- Update/delete data from their organization

**Implementation:**
```sql
SET LOCAL app.current_organization_id = '<orgID>';
-- All subsequent queries automatically filtered
```

---

## 📊 Response Format

### Success Response
```json
{
  "id": "uuid",
  "field1": "value1",
  "field2": "value2",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### List Response
```json
{
  "items": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Response
```json
{
  "error": "Validation failed",
  "message": "email is required"
}
```

---

## 🎯 Integration Examples

### cURL
```bash
# Create customer
curl -X POST \
  http://localhost:8080/api/v1/organizations/550e8400.../customers \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"John","last_name":"Doe"}'
```

### JavaScript (Fetch)
```javascript
const response = await fetch(
  'http://localhost:8080/api/v1/organizations/550e8400.../customers',
  {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer YOUR_TOKEN',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      first_name: 'John',
      last_name: 'Doe'
    })
  }
);
const data = await response.json();
```

### Flutter (Dart)
```dart
import 'package:http/http.dart' as http;
import 'dart:convert';

Future<Customer> createCustomer(String orgId, Map<String, dynamic> data) async {
  final response = await http.post(
    Uri.parse('http://localhost:8080/api/v1/organizations/$orgId/customers'),
    headers: {
      'Authorization': 'Bearer YOUR_TOKEN',
      'Content-Type': 'application/json',
    },
    body: jsonEncode(data),
  );

  if (response.statusCode == 201) {
    return Customer.fromJson(jsonDecode(response.body));
  }
  throw Exception('Failed to create customer');
}
```

---

## 📚 Additional Resources

- **Swagger UI:** http://localhost:8080/swagger/index.html
- **OpenAPI Spec:** http://localhost:8080/swagger/swagger.json
- **Metrics:** http://localhost:8080/metrics
- **Health Check:** http://localhost:8080/health

---

**Total Endpoints:** 855+ (171 tables × 5 operations)
**All APIs:** Multi-tenant, Authenticated, Paginated, Validated
**Documentation:** Swagger/OpenAPI annotations included
**Status:** ✅ Generated and Ready


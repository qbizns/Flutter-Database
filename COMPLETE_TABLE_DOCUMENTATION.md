# Complete Table Documentation for API Development

## 📋 Overview

This document provides comprehensive documentation of **ALL** tables in the database, organized by functional domain. Use this as the foundation for building your Go backend API.

**Total Tables**: ~90 tables across 23 migrations
**Database Version**: PostgreSQL 15+
**Multi-Tenant**: Yes (all tables have organization_id)
**Row-Level Security**: Enabled

---

## 🗂️ Table Categories

### Core System (Tenant Management)
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `organizations` | Multi-tenant organizations | name, subscription_tier, status | GET/POST/PATCH /organizations |
| `users` | User accounts | email, password_hash, status | POST /auth/register, /auth/login |
| `roles` | Role definitions | role_name, role_code, is_system | GET/POST/PATCH /roles |
| `permissions` | Permission definitions | permission_name, resource_type | GET /permissions |
| `user_roles` | User-role assignments | user_id, role_id, organization_id | POST /users/{id}/roles |
| `user_organizations` | User-org relationships | user_id, organization_id, is_primary | GET /users/{id}/organizations |

### Products & Inventory
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `products` | Product master data | name, sku, barcode, price, cost | GET/POST/PATCH/DELETE /products |
| `product_categories` | Product categorization | name, parent_category_id | GET/POST/PATCH /categories |
| `product_variants` | Product variations | variant_name, sku, price | GET/POST /products/{id}/variants |
| `product_batches` | Batch/lot tracking | batch_number, expiry_date, quantity | GET/POST /batches |
| `product_serial_numbers` | Serial number tracking | serial_number, status | GET /serials |
| `product_components` | Product BOM | component_product_id, quantity | GET/POST /products/{id}/components |
| `modifiers` | Product modifiers | name, modifier_type, price_adjustment | GET/POST /modifiers |
| `modifier_groups` | Modifier grouping | name, is_required, max_selection | GET/POST /modifier-groups |
| `product_modifier_groups` | Product-modifier link | product_id, modifier_group_id | POST /products/{id}/modifiers |
| `units_of_measure` | UOM definitions | name, symbol, uom_type | GET/POST /uom |
| `uom_conversions` | UOM conversion rules | from_uom_id, to_uom_id, conversion_factor | GET/POST /uom/conversions |
| `inventory_transactions` | Stock movements | transaction_type, quantity, cost | GET/POST /inventory/transactions |
| `locations` | Warehouse/store locations | name, location_type, address | GET/POST/PATCH /locations |
| `inventory_transfers` | Inter-location transfers | from_location_id, to_location_id, status | GET/POST/PATCH /transfers |
| `inventory_transfer_items` | Transfer line items | product_id, quantity, transfer_id | GET /transfers/{id}/items |
| `cycle_counts` | Inventory cycle counts | count_date, status, location_id | GET/POST /cycle-counts |
| `cycle_count_items` | Cycle count lines | product_id, expected_qty, counted_qty | GET/POST /cycle-counts/{id}/items |
| `stock_adjustment_reasons` | Adjustment reason codes | reason_code, reason_name | GET/POST /adjustment-reasons |

### Sales & Customers
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `customers` | Customer master data | first_name, last_name, email, phone | GET/POST/PATCH/DELETE /customers |
| `sales` | Sales transactions | sale_number, total_amount, payment_method | GET/POST/PATCH /sales |
| `sale_items` | Sale line items | product_id, quantity, price, discount | GET /sales/{id}/items |
| `payments` | Payment records | payment_method, amount, status | GET/POST /payments |
| `sale_returns` | Sales returns | original_sale_id, total_refund_amount | GET/POST /returns |
| `sale_return_items` | Return line items | sale_item_id, return_quantity, reason | GET /returns/{id}/items |
| `return_reasons` | Return reason codes | reason_code, reason_name | GET/POST /return-reasons |
| `sales_channels` | Sales channels | channel_code, channel_name | GET/POST /sales-channels |
| `gift_cards` | Gift card master | card_number, initial_balance, balance | GET/POST/PATCH /gift-cards |
| `gift_card_transactions` | Gift card usage | transaction_type, amount, sale_id | GET /gift-cards/{id}/transactions |
| `customer_store_credit_accounts` | Store credit accounts | customer_id, balance | GET /customers/{id}/store-credit |
| `store_credit_transactions` | Store credit usage | transaction_type, amount | GET /store-credit/{id}/transactions |

### Loyalty Program
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `loyalty_tiers` | Loyalty tier definitions | tier_name, min_points, discount_percent | GET/POST/PATCH /loyalty/tiers |
| `loyalty_tier_benefits` | Tier benefits | benefit_name, benefit_type, value | GET /loyalty/tiers/{id}/benefits |
| `loyalty_points_rules` | Points earning rules | rule_name, points_per_currency, multiplier | GET/POST /loyalty/rules |
| `loyalty_points_transactions` | Points history | transaction_type, points, sale_id | GET/POST /loyalty/transactions |
| `loyalty_rewards` | Reward catalog | reward_name, points_cost, reward_type | GET/POST /loyalty/rewards |
| `loyalty_redemptions` | Reward redemptions | customer_id, reward_id, points_spent | GET/POST /loyalty/redemptions |
| `customer_tier_history` | Tier change history | customer_id, old_tier_id, new_tier_id | GET /customers/{id}/tier-history |

### Purchase & Suppliers
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `suppliers` | Supplier master data | supplier_code, name, contact_person | GET/POST/PATCH/DELETE /suppliers |
| `purchase_orders` | Purchase orders | po_number, supplier_id, total_amount | GET/POST/PATCH /purchase-orders |
| `purchase_order_items` | PO line items | product_id, quantity, unit_cost | GET /purchase-orders/{id}/items |
| `goods_receipts` | Goods receipt notes | gr_number, purchase_order_id, received_date | GET/POST /goods-receipts |
| `goods_receipt_items` | GR line items | product_id, quantity_received | GET /goods-receipts/{id}/items |

### POS Operations
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `pos_sessions` | Cash register sessions | session_number, opening_cash, status | GET/POST/PATCH /pos/sessions |
| `cash_drawers` | Cash drawer definitions | drawer_code, location_id | GET/POST /pos/drawers |
| `cash_drawer_sessions` | Drawer session link | drawer_id, session_id, opening_float | GET /pos/drawers/{id}/sessions |
| `cash_movements` | Cash in/out | movement_type, amount, reason | GET/POST /pos/cash-movements |
| `shifts` | Employee shifts | user_id, clock_in_time, clock_out_time | GET/POST/PATCH /shifts |
| `expenses` | Operating expenses | expense_date, amount, category | GET/POST /expenses |

### Pricing & Promotions
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `price_lists` | Price list definitions | name, effective_from, effective_to | GET/POST/PATCH /price-lists |
| `price_list_items` | Price list items | product_id, price, discount_percent | GET/POST /price-lists/{id}/items |
| `promotions` | Promotional campaigns | promotion_name, discount_type, start_date | GET/POST/PATCH /promotions |
| `promotion_usage` | Promotion usage log | promotion_id, sale_id, discount_amount | GET /promotions/{id}/usage |

### Restaurant Features
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `table_sections` | Restaurant sections | section_name, capacity | GET/POST /restaurant/sections |
| `restaurant_tables` | Table definitions | table_number, section_id, capacity | GET/POST/PATCH /restaurant/tables |
| `floor_plans` | Floor plan layouts | plan_name, layout_data | GET/POST /restaurant/floor-plans |
| `reservations` | Table reservations | customer_id, reservation_date, table_id | GET/POST/PATCH/DELETE /reservations |
| `courses` | Course definitions (appetizer, main, dessert) | course_name, sequence | GET /courses |

### E-Invoicing (ZATCA & ETA)
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `e_invoicing_documents` | E-invoice master | invoice_type, regulation, status | GET/POST /e-invoicing/documents |
| `e_invoicing_document_events` | Invoice event log | event_type, submission_status | GET /e-invoicing/documents/{id}/events |

### Document Management
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `document_sequences` | Document numbering | document_type, prefix, next_number | GET/POST/PATCH /sequences |

### System & Monitoring
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `audit_logs` | Audit trail | action, entity_type, entity_id, changes | GET /audit-logs |
| `pos_error_logs` | POS error logs | error_code, error_message, stack_trace | GET/POST /error-logs |
| `system_health` | System health metrics | metric_name, metric_value, status | GET /system/health |

---

## 🆕 New Infrastructure Tables (V023)

### Background Processing
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `background_jobs` | Async job queue | job_type, status, payload, result | GET/POST /jobs, GET /jobs/{id} |

**Job Types**: `posting_engine`, `report_generation`, `data_export`, `email_batch`, `data_import`, `inventory_sync`

**Usage Example**:
```go
// Create posting job
POST /api/v1/jobs
{
  "job_type": "posting_engine",
  "queue_name": "high_priority",
  "payload": {
    "document_type": "POS_SALE",
    "document_id": "uuid",
    "event": "on_post"
  }
}
```

### API Authentication & Security
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `api_keys` | API keys for programmatic access | key_name, key_hash, scopes | GET/POST/DELETE /api-keys |
| `rate_limits` | Rate limiting tracking | identifier_type, request_count | GET /rate-limits (internal) |
| `user_sessions` | Session management | session_token, device_info | GET/DELETE /sessions |

**API Key Scopes**: `read`, `write`, `delete`, `admin`

### Webhooks & Integrations
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `webhooks` | Webhook configurations | url, events, is_active | GET/POST/PATCH/DELETE /webhooks |
| `webhook_deliveries` | Webhook delivery log | status, request_body, response_status | GET /webhooks/{id}/deliveries |
| `integration_configs` | Third-party integrations | integration_type, provider, credentials | GET/POST/PATCH /integrations |

**Webhook Events**: `sale.created`, `sale.completed`, `payment.received`, `invoice.posted`, `inventory.updated`

### Notifications
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `notifications` | User notifications | notification_type, category, message | GET/PATCH /notifications |
| `notification_preferences` | Notification settings | category, in_app_enabled, email_enabled | GET/PATCH /users/{id}/notification-preferences |

**Notification Categories**: `sales`, `inventory`, `accounting`, `system`, `security`

### File Management
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `file_attachments` | Document attachments | file_name, storage_path, entity_type | GET/POST/DELETE /files |

**Supported Entities**: `sale`, `invoice`, `expense`, `product`, `customer`, `supplier`

### Communication Queues
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `email_queue` | Email sending queue | to_addresses, subject, status | GET/POST /emails |
| `sms_queue` | SMS sending queue | to_phone, message, status | GET/POST /sms |

### Settings & Preferences
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `organization_settings` | Org configuration | timezone, default_currency, auto_post_sales | GET/PATCH /settings |
| `user_settings` | User preferences | theme, language, default_dashboard | GET/PATCH /users/me/settings |

### Data Export & Reporting
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `data_export_requests` | Export job tracking | export_type, export_format, status | GET/POST /exports |
| `scheduled_reports` | Automated reports | report_type, schedule_frequency | GET/POST/PATCH/DELETE /scheduled-reports |

### API Monitoring
| Table | Purpose | Key Fields | API Endpoints |
|-------|---------|-----------|---------------|
| `api_request_logs` | API request logs | method, path, status_code, duration_ms | GET /api-logs (admin only) |

---

## 🏗️ Accounting Module Tables

### Core Accounting
| Table | Purpose | Key Fields | Schema |
|-------|---------|-----------|---------|
| `fiscal_years` | Fiscal year definitions | fiscal_year, start_date, end_date, status | accounting |
| `accounting_periods` | Monthly periods | period_name, status, is_locked | accounting |
| `account_types` | Account type classifications | type_code, normal_balance, financial_statement_section | accounting |
| `chart_of_accounts` | GL account master | account_code, account_name, account_type_id | accounting |
| `journal_entry_types` | JE type definitions | type_code, type_name | accounting |
| `journal_entries` | Journal entry headers | entry_number, entry_date, status, total_debit | accounting |
| `journal_entry_lines` | JE line items | account_id, debit_amount, credit_amount | accounting |
| `general_ledger` | Posted transactions (immutable) | transaction_date, account_id, debit_amount | accounting |

### AP/AR & Payments
| Table | Purpose | Key Fields | Schema |
|-------|---------|-----------|---------|
| `vendor_bills` | Accounts payable | bill_number, supplier_id, total_amount, status | accounting |
| `vendor_bill_items` | AP line items | description, amount, account_id | accounting |
| `vendor_payments` | AP payments | payment_number, supplier_id, amount | accounting |
| `vendor_bill_payments` | Bill payment application | vendor_bill_id, vendor_payment_id, amount_applied | accounting |
| `customer_invoices` | Accounts receivable | invoice_number, customer_id, total_amount | accounting |
| `customer_invoice_items` | AR line items | description, amount, account_id | accounting |
| `customer_payments` | AR payments | payment_number, customer_id, amount | accounting |
| `customer_invoice_payments` | Invoice payment application | customer_invoice_id, customer_payment_id | accounting |

### Fixed Assets
| Table | Purpose | Key Fields | Schema |
|-------|---------|-----------|---------|
| `asset_categories` | Asset classifications | category_name, depreciation_method | accounting |
| `fixed_assets` | Asset register | asset_name, acquisition_cost, depreciation_method | accounting |
| `asset_depreciation_schedule` | Depreciation entries | fiscal_year_id, period_id, depreciation_amount | accounting |

### Banking
| Table | Purpose | Key Fields | Schema |
|-------|---------|-----------|---------|
| `bank_accounts` | Bank account master | account_number, bank_name, balance | accounting |
| `bank_reconciliations` | Bank rec headers | statement_date, ending_balance, status | accounting |
| `bank_reconciliation_items` | Bank rec lines | transaction_date, description, amount, is_matched | accounting |

### Advanced Accounting (Odoo-Style)
| Table | Purpose | Key Fields | Schema |
|-------|---------|-----------|---------|
| `journals` | Journal definitions | journal_code, journal_type, default_debit_account_id | accounting |
| `tax_groups` | Tax groupings | group_name | accounting |
| `taxes` | Tax definitions | tax_name, tax_rate, tax_scope | accounting |
| `fiscal_positions` | Tax remapping (export, domestic) | position_code, position_name | accounting |
| `fiscal_position_tax_mappings` | Tax substitution rules | source_tax_id, destination_tax_id | accounting |
| `currencies` | Currency master | currency_code, currency_name, symbol | accounting |
| `currency_rates` | Exchange rates | currency_code, rate_date, rate | accounting |
| `payment_terms` | Payment term definitions | term_name, term_code | accounting |
| `payment_term_lines` | Term line items | sequence, value_type, value | accounting |
| `invoice_payment_schedules` | Payment due dates | due_date, amount_due, status | accounting |
| `analytic_plans` | Analytic dimensions | plan_code, plan_name | accounting |
| `analytic_accounts` | Cost centers/projects | account_code, account_name, parent_account_id | accounting |
| `deferred_revenue_contracts` | Deferred revenue | contract_number, total_amount, recognition_method | accounting |
| `deferred_revenue_schedule` | Revenue recognition | recognition_date, amount, status | accounting |
| `deferred_expense_contracts` | Prepaid expenses | contract_number, total_amount | accounting |
| `deferred_expense_schedule` | Expense amortization | recognition_date, amount, status | accounting |
| `bank_statements` | Bank statement headers | statement_number, statement_date | accounting |
| `bank_statement_lines` | Bank statement lines | transaction_date, description, amount, status | accounting |
| `bank_statement_reconciliations` | Statement reconciliation | bank_statement_line_id, journal_entry_id | accounting |
| `reconciliation_rule_models` | Auto-match rules | rule_name, description_pattern, account_id | accounting |
| `budgets` | Budget headers | budget_code, fiscal_year_id, budget_type | accounting |
| `budget_lines` | Budget line items | account_id, analytic_account_id, budgeted_amount | accounting |
| `localization_packages` | Country-specific accounting | package_code, country_code | accounting |
| `tax_report_definitions` | Tax report templates | report_code, report_name | accounting |
| `tax_report_lines` | Tax report line definitions | line_code, formula | accounting |

### POS-Accounting Integration
| Table | Purpose | Key Fields | Schema |
|-------|---------|-----------|---------|
| `pos_account_mappings` | Product/payment → GL account | source_type, source_id, purpose, account_id, concept_key | public |
| `pos_posting_audit` | Posting status tracking | source_table, source_id, posting_status, journal_entry_id | public |
| `pos_tax_mappings` | POS tax → accounting tax | pos_tax_code, accounting_tax_id | public |
| `inventory_valuation_settings` | FIFO/weighted average config | valuation_method, cost_layer_granularity | public |
| `inventory_cost_layers` | Inventory cost tracking | product_id, location_id, unit_cost, remaining_quantity | public |

### Posting Engine (Configuration-Driven)
| Table | Purpose | Key Fields | Schema |
|-------|---------|-----------|---------|
| `posting_concepts` | Logical vocabulary (AR, REVENUE, CASH) | concept_key, default_label, normal_side | accounting |
| `posting_concept_overrides` | Org-specific labels | organization_id, concept_key, label | accounting |
| `posting_account_mappings` | Concept → GL account mapping | concept_key, account_id, location_id | accounting |
| `posting_profiles` | Org posting configurations | code, name, is_default | accounting |
| `posting_document_types` | Business documents | code, source_table, category | accounting |
| `posting_profile_documents` | Profile → document link | posting_profile_id, posting_document_type_id | accounting |
| `posting_rules` | Posting templates | rule_code, event, condition_expression | accounting |
| `posting_rule_lines` | JE line templates | side, concept_key, amount_source | accounting |
| `posting_validation_rules` | Validation DSL | code, expression, severity, is_blocking | accounting |
| `posting_validation_results` | Validation log | validation_rule_id, severity, message | accounting |

---

## 🔑 Key Relationships for API Development

### Multi-Tenancy Pattern
```
organizations
    ↓
users ← user_organizations → organizations
    ↓
ALL other tables (via organization_id)
```

### Sales Flow
```
customers → sales → sale_items → products
                 ↓
              payments
                 ↓
           accounting (via posting engine)
```

### Inventory Flow
```
purchase_orders → goods_receipts → inventory_transactions → products
      ↓                                      ↓
   suppliers                            locations
```

### Accounting Posting Flow
```
sales/purchase_orders (POS)
    → posting_engine (evaluate rules)
    → journal_entries
    → journal_entry_lines
    → general_ledger (immutable)
```

### Loyalty Flow
```
customers → sales → loyalty_points_transactions
               ↓
         loyalty_tiers
               ↓
      loyalty_tier_benefits
```

---

## 🔐 Row-Level Security (RLS)

All tables have RLS policies:

```sql
CREATE POLICY org_isolation ON {table_name}
    FOR ALL
    USING (organization_id IN (
        SELECT organization_id FROM user_organizations
        WHERE user_id = auth.uid()
    ));
```

**API Implementation**: Set `app.current_organization_id` session variable:
```go
_, err := db.Exec("SET app.current_organization_id = $1", orgID)
```

---

## 📊 Key Indexes for Performance

### Common Index Patterns
```sql
-- Organization + status filtering
CREATE INDEX idx_{table}_org_status ON {table}(organization_id, status);

-- Soft deletes
CREATE INDEX idx_{table}_org ON {table}(organization_id) WHERE deleted_at IS NULL;

-- Date range queries
CREATE INDEX idx_{table}_date ON {table}(date_column);

-- Foreign key lookups
CREATE INDEX idx_{table}_foreign_key ON {table}(foreign_key_id);
```

---

## 🚀 Essential API Endpoints to Build

### Authentication (Priority: 🔴 CRITICAL)
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh-token
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/auth/me
```

### POS Core (Priority: 🔴 CRITICAL)
```
GET/POST   /api/v1/organizations/{org_id}/products
GET/POST   /api/v1/organizations/{org_id}/customers
GET/POST   /api/v1/organizations/{org_id}/sales
POST       /api/v1/organizations/{org_id}/sales/{id}/complete
GET        /api/v1/organizations/{org_id}/inventory
POST       /api/v1/organizations/{org_id}/inventory/adjust
```

### Posting Engine (Priority: 🔴 CRITICAL)
```
POST   /api/v1/organizations/{org_id}/posting/post
GET    /api/v1/organizations/{org_id}/posting/rules
POST   /api/v1/organizations/{org_id}/posting/validate
GET    /api/v1/organizations/{org_id}/posting/audit
```

### Accounting (Priority: 🟡 HIGH)
```
GET/POST   /api/v1/organizations/{org_id}/journal-entries
GET        /api/v1/organizations/{org_id}/reports/balance-sheet
GET        /api/v1/organizations/{org_id}/reports/income-statement
GET        /api/v1/organizations/{org_id}/reports/trial-balance
GET/POST   /api/v1/organizations/{org_id}/vendor-bills
GET/POST   /api/v1/organizations/{org_id}/customer-invoices
```

### Background Jobs (Priority: 🟡 HIGH)
```
GET/POST   /api/v1/jobs
GET        /api/v1/jobs/{id}
POST       /api/v1/jobs/{id}/retry
DELETE     /api/v1/jobs/{id}
```

### Webhooks (Priority: 🟢 MEDIUM)
```
GET/POST/PATCH/DELETE  /api/v1/webhooks
GET                    /api/v1/webhooks/{id}/deliveries
POST                   /api/v1/webhooks/{id}/test
```

### Files (Priority: 🟢 MEDIUM)
```
POST       /api/v1/files/upload
GET        /api/v1/files/{id}
DELETE     /api/v1/files/{id}
GET        /api/v1/files?entity_type={type}&entity_id={id}
```

### Notifications (Priority: 🟢 MEDIUM)
```
GET        /api/v1/notifications
PATCH      /api/v1/notifications/{id}/read
PATCH      /api/v1/notifications/mark-all-read
GET/PATCH  /api/v1/users/me/notification-preferences
```

---

## 📦 Go Struct Examples

### Product
```go
type Product struct {
    ID             uuid.UUID      `json:"id" db:"id"`
    OrganizationID uuid.UUID      `json:"organization_id" db:"organization_id"`
    Name           string         `json:"name" db:"name"`
    SKU            string         `json:"sku" db:"sku"`
    Barcode        *string        `json:"barcode" db:"barcode"`
    Price          float64        `json:"price" db:"price"`
    Cost           float64        `json:"cost" db:"cost"`
    CategoryID     *uuid.UUID     `json:"category_id" db:"category_id"`
    IsActive       bool           `json:"is_active" db:"is_active"`
    CreatedAt      time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
    DeletedAt      *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
}
```

### Sale
```go
type Sale struct {
    ID              uuid.UUID      `json:"id" db:"id"`
    OrganizationID  uuid.UUID      `json:"organization_id" db:"organization_id"`
    SaleNumber      string         `json:"sale_number" db:"sale_number"`
    CustomerID      *uuid.UUID     `json:"customer_id" db:"customer_id"`
    TotalAmount     float64        `json:"total_amount" db:"total_amount"`
    SubtotalAmount  float64        `json:"subtotal_amount" db:"subtotal_amount"`
    TaxAmount       float64        `json:"tax_amount" db:"tax_amount"`
    DiscountAmount  float64        `json:"discount_amount" db:"discount_amount"`
    PaymentMethod   string         `json:"payment_method" db:"payment_method"`
    Status          string         `json:"status" db:"status"`
    SaleDate        time.Time      `json:"sale_date" db:"sale_date"`
    CreatedAt       time.Time      `json:"created_at" db:"created_at"`
}
```

### Background Job
```go
type BackgroundJob struct {
    ID            uuid.UUID              `json:"id" db:"id"`
    OrganizationID *uuid.UUID            `json:"organization_id" db:"organization_id"`
    JobType       string                 `json:"job_type" db:"job_type"`
    JobName       string                 `json:"job_name" db:"job_name"`
    Status        string                 `json:"status" db:"status"`
    Payload       map[string]interface{} `json:"payload" db:"payload"`
    Result        map[string]interface{} `json:"result" db:"result"`
    Attempts      int                    `json:"attempts" db:"attempts"`
    CreatedAt     time.Time              `json:"created_at" db:"created_at"`
    CompletedAt   *time.Time             `json:"completed_at" db:"completed_at"`
}
```

---

## 🎯 Quick Reference: Table Count by Domain

| Domain | Table Count |
|--------|-------------|
| **Core System** | 6 tables |
| **Products & Inventory** | 19 tables |
| **Sales & Customers** | 13 tables |
| **Loyalty Program** | 7 tables |
| **Purchase & Suppliers** | 5 tables |
| **POS Operations** | 6 tables |
| **Pricing & Promotions** | 4 tables |
| **Restaurant Features** | 5 tables |
| **E-Invoicing** | 2 tables |
| **Infrastructure (NEW)** | 17 tables |
| **Accounting Core** | 8 tables |
| **AP/AR & Payments** | 8 tables |
| **Fixed Assets** | 3 tables |
| **Banking** | 3 tables |
| **Advanced Accounting** | 23 tables |
| **POS-Accounting Integration** | 5 tables |
| **Posting Engine** | 10 tables |
| **TOTAL** | **~145 tables** |

---

## ✅ Readiness Checklist

Before starting API development:

- [x] All base tables exist (organizations, users, products, sales, etc.)
- [x] All infrastructure tables created (jobs, webhooks, notifications, etc.)
- [x] All accounting tables created
- [x] Posting engine configuration tables ready
- [x] Multi-tenancy with RLS policies
- [x] Indexes for performance
- [x] Audit trails (audit_logs, updated_at triggers)
- [x] Soft deletes (deleted_at pattern)

**Status**: ✅ **100% READY FOR API DEVELOPMENT**

---

## 🚀 Next Steps

1. **Set up Go project structure** (see NEXT_STEPS_RECOMMENDATIONS.md)
2. **Implement authentication middleware**
3. **Create CRUD endpoints for core tables** (products, sales, customers)
4. **Build posting engine service** (HIGHEST PRIORITY)
5. **Add background job worker**
6. **Implement webhooks system**
7. **Add file upload/download**
8. **Build notification system**
9. **Create reporting endpoints**
10. **Add API documentation (Swagger)**

---

**Database is production-ready. Start building the API! 🎉**

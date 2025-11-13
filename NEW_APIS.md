# Backend API Requirements Documentation

**Project:** Flutter-Base POS System
**Date:** 2025-11-13
**Purpose:** Complete API endpoint specifications for backend implementation
**Backend Repository:** https://github.com/qbizns/Flutter-Database/tree/claude/pos-database-setup-011CUxJ8SiQmm5Zoj6SqGfZ9/backend

---

## 📋 Overview

This document specifies ALL required REST API endpoints for the Flutter-Base POS system. The Flutter frontend has HTTP client implementations ready and will automatically use these APIs when configured with:
- `apiBaseUrl`: Your backend API base URL
- `organizationId`: Multi-tenant organization identifier

**Architecture:** Multi-tenant with organization-based routing
**Authentication:** JWT Bearer token
**Format:** JSON (snake_case for API, camelCase in Flutter)

---

## 🔐 Authentication

### Base Headers
```http
Authorization: Bearer {jwt_token}
Content-Type: application/json
Accept: application/json
```

### JWT Token Structure
```json
{
  "user_id": "uuid",
  "organization_id": "uuid",
  "email": "user@example.com",
  "roles": ["admin", "cashier"],
  "exp": 1234567890,
  "iat": 1234567890
}
```

---

## 📦 API Endpoints by Module

## 1. ORDERS MODULE ✅ (CRITICAL - Priority 1)

**Base Path:** `/api/v1/organizations/{org_id}/orders`

### 1.1 Create Order
```http
POST /api/v1/organizations/{org_id}/orders
```

**Request Body:**
```json
{
  "order_number": "ORD-001",
  "status": "pending",
  "order_type": "dineIn",
  "items": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "product_name": "Margherita Pizza",
      "base_price": 12.99,
      "quantity": 1,
      "selected_modifiers": [
        {
          "modifier_id": "uuid",
          "modifier_name": "Extra Cheese",
          "price": 1.50
        }
      ],
      "notes": "No onions",
      "tax_percent": 8.5
    }
  ],
  "table_id": "uuid",
  "table_name": "Table 5",
  "customer_id": "uuid",
  "customer_name": "John Doe",
  "customer_phone": "+1234567890",
  "customer_email": "john@example.com",
  "notes": "Order notes",
  "subtotal": 35.47,
  "discount_amount": 0,
  "discount_percent": 0,
  "tax_amount": 3.02,
  "tax_percent": 8.5,
  "tip_amount": 5.00,
  "total": 43.49,
  "payment_method": "cash",
  "payment_status": "pending"
}
```

**Response:** 201 Created
```json
{
  "id": "uuid",
  "order_number": "ORD-001",
  "status": "pending",
  "order_type": "dineIn",
  "items": [...],
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z",
  "completed_at": null,
  ...
}
```

### 1.2 Get Orders (List with Filters)
```http
GET /api/v1/organizations/{org_id}/orders?status={status}&type={type}&table_id={table_id}&from_date={iso8601}&to_date={iso8601}
```

**Query Parameters:**
- `status` (optional): pending, preparing, ready, completed, cancelled
- `type` (optional): dineIn, takeaway, delivery
- `table_id` (optional): Filter by table UUID
- `from_date` (optional): ISO 8601 datetime
- `to_date` (optional): ISO 8601 datetime

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "order_number": "ORD-001",
      "status": "pending",
      ...
    }
  ],
  "total": 150,
  "page": 1,
  "per_page": 50
}
```

### 1.3 Get Order by ID
```http
GET /api/v1/organizations/{org_id}/orders/{order_id}
```

**Response:** 200 OK (Same structure as Create Order response)

### 1.4 Update Order
```http
PATCH /api/v1/organizations/{org_id}/orders/{order_id}
```

**Request Body:** (Same as Create Order, all fields optional)

**Response:** 200 OK (Updated order)

### 1.5 Update Order Status
```http
PATCH /api/v1/organizations/{org_id}/orders/{order_id}/status
```

**Request Body:**
```json
{
  "status": "preparing"
}
```

**Response:** 200 OK (Updated order)

### 1.6 Cancel Order
```http
POST /api/v1/organizations/{org_id}/orders/{order_id}/cancel
```

**Request Body:**
```json
{
  "reason": "Customer requested cancellation"
}
```

**Response:** 200 OK (Cancelled order)

### 1.7 Get Order Statistics
```http
GET /api/v1/organizations/{org_id}/orders/statistics?from_date={iso8601}&to_date={iso8601}
```

**Response:** 200 OK
```json
{
  "total_orders": 150,
  "completed_orders": 120,
  "cancelled_orders": 5,
  "active_orders": 25,
  "total_revenue": 15420.50,
  "average_order_value": 128.50,
  "average_preparation_minutes": 15
}
```

---

## 2. PRODUCTS MODULE ✅ (CRITICAL - Priority 1)

**Base Path:** `/api/v1/organizations/{org_id}/products`

### 2.1 Get Products (List with Filters)
```http
GET /api/v1/organizations/{org_id}/products?category_id={uuid}&active_only={boolean}
```

**Query Parameters:**
- `category_id` (optional): Filter by category UUID
- `active_only` (optional): true/false, default true

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Margherita Pizza",
      "sku": "PIZZA-001",
      "description": "Classic pizza with tomato and mozzarella",
      "category_id": "uuid",
      "category_name": "Pizzas",
      "base_price": 12.99,
      "cost_price": 5.00,
      "barcode": "1234567890123",
      "is_available": true,
      "is_featured": false,
      "track_inventory": true,
      "current_stock": 50,
      "low_stock_threshold": 10,
      "preparation_time_minutes": 15,
      "calories": 800,
      "allergens": ["gluten", "dairy"],
      "tags": ["vegetarian", "popular"],
      "image_url": "https://...",
      "prices": [
        {
          "size": "Medium",
          "price": 12.99,
          "is_default": true
        }
      ],
      "modifier_groups": [
        {
          "id": "uuid",
          "name": "Toppings",
          "selection_type": "multiple",
          "min_selections": 0,
          "max_selections": 5,
          "is_required": false,
          "modifiers": [
            {
              "id": "uuid",
              "name": "Extra Cheese",
              "price": 1.50,
              "is_available": true,
              "is_default": false
            }
          ]
        }
      ],
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 2.2 Get Product by ID
```http
GET /api/v1/organizations/{org_id}/products/{product_id}
```

**Response:** 200 OK (Same structure as above, single product)

### 2.3 Search Products
```http
GET /api/v1/organizations/{org_id}/products/search?q={query}&category_id={uuid}
```

**Query Parameters:**
- `q` (required): Search query string
- `category_id` (optional): Filter by category

**Response:** 200 OK (Array of products matching query)

### 2.4 Get Products by IDs (Batch)
```http
POST /api/v1/organizations/{org_id}/products/batch
```

**Request Body:**
```json
{
  "product_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**Response:** 200 OK (Array of products)

### 2.5 Get Featured Products
```http
GET /api/v1/organizations/{org_id}/products/featured?limit={number}
```

**Query Parameters:**
- `limit` (optional): Number of products to return, default 10

**Response:** 200 OK (Array of featured products)

### 2.6 Get Low Stock Products
```http
GET /api/v1/organizations/{org_id}/products/low-stock
```

**Response:** 200 OK (Array of products below low stock threshold)

### 2.7 Update Product Stock
```http
PATCH /api/v1/organizations/{org_id}/products/{product_id}/stock
```

**Request Body:**
```json
{
  "quantity": 100
}
```

**Response:** 200 OK

### 2.8 Update Product Availability
```http
PATCH /api/v1/organizations/{org_id}/products/{product_id}/availability
```

**Request Body:**
```json
{
  "is_available": true
}
```

**Response:** 200 OK

### 2.9 Create Product
```http
POST /api/v1/organizations/{org_id}/products
```

**Request Body:** (Same structure as Get Products response)

**Response:** 201 Created

### 2.10 Update Product
```http
PATCH /api/v1/organizations/{org_id}/products/{product_id}
```

**Request Body:** (Same as Create, all fields optional)

**Response:** 200 OK

### 2.11 Delete Product
```http
DELETE /api/v1/organizations/{org_id}/products/{product_id}
```

**Response:** 204 No Content

---

## 3. CATEGORIES MODULE ✅ (Priority 1)

**Base Path:** `/api/v1/organizations/{org_id}/categories`

### 3.1 Get Categories
```http
GET /api/v1/organizations/{org_id}/categories?parent_id={uuid}&active_only={boolean}
```

**Query Parameters:**
- `parent_id` (optional): Get subcategories of parent
- `active_only` (optional): true/false, default true

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Pizzas",
      "description": "All pizza varieties",
      "parent_id": null,
      "image_url": "https://...",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 3.2 Get Category by ID
```http
GET /api/v1/organizations/{org_id}/categories/{category_id}
```

**Response:** 200 OK

### 3.3 Get Products Count by Category
```http
GET /api/v1/organizations/{org_id}/products/count-by-category
```

**Response:** 200 OK
```json
{
  "uuid1": 25,
  "uuid2": 18,
  "uuid3": 42
}
```

### 3.4 Create Category
```http
POST /api/v1/organizations/{org_id}/categories
```

**Request Body:**
```json
{
  "name": "Pizzas",
  "description": "All pizza varieties",
  "parent_id": null,
  "image_url": "https://...",
  "sort_order": 1,
  "is_active": true
}
```

**Response:** 201 Created

### 3.5 Update Category
```http
PATCH /api/v1/organizations/{org_id}/categories/{category_id}
```

**Response:** 200 OK

### 3.6 Delete Category
```http
DELETE /api/v1/organizations/{org_id}/categories/{category_id}
```

**Response:** 204 No Content

---

## 4. TABLES MODULE ✅ (CRITICAL - Priority 1)

**Base Path:** `/api/v1/organizations/{org_id}/tables`

### 4.1 Get Tables
```http
GET /api/v1/organizations/{org_id}/tables?status={status}&zone_id={uuid}
```

**Query Parameters:**
- `status` (optional): available, occupied, reserved, cleaning
- `zone_id` (optional): Filter by zone UUID

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Table 5",
      "number": 5,
      "capacity": 4,
      "status": "available",
      "shape": "rectangle",
      "section": "Main Hall",
      "position": {
        "x": 100,
        "y": 200,
        "width": 80,
        "height": 80,
        "rotation": 0
      },
      "current_order_id": null,
      "customer_id": null,
      "customer_name": null,
      "waiter_id": null,
      "waiter_name": null,
      "seated_at": null,
      "reservation_time": null,
      "reservation_name": null,
      "notes": null,
      "zone_id": "uuid",
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 4.2 Get Table by ID
```http
GET /api/v1/organizations/{org_id}/tables/{table_id}
```

**Response:** 200 OK

### 4.3 Update Table Status
```http
PATCH /api/v1/organizations/{org_id}/tables/{table_id}/status
```

**Request Body:**
```json
{
  "status": "occupied"
}
```

**Response:** 200 OK

### 4.4 Assign Order to Table
```http
POST /api/v1/organizations/{org_id}/tables/{table_id}/assign-order
```

**Request Body:**
```json
{
  "order_id": "uuid",
  "customer_name": "John Doe",
  "waiter_id": "uuid"
}
```

**Response:** 200 OK

### 4.5 Clear Table
```http
POST /api/v1/organizations/{org_id}/tables/{table_id}/clear
```

**Response:** 200 OK

### 4.6 Create Table
```http
POST /api/v1/organizations/{org_id}/tables
```

**Request Body:** (Same structure as Get Tables response)

**Response:** 201 Created

### 4.7 Update Table
```http
PATCH /api/v1/organizations/{org_id}/tables/{table_id}
```

**Response:** 200 OK

### 4.8 Delete Table
```http
DELETE /api/v1/organizations/{org_id}/tables/{table_id}
```

**Response:** 204 No Content

### 4.9 Get Table Statistics
```http
GET /api/v1/organizations/{org_id}/tables/statistics
```

**Response:** 200 OK
```json
{
  "total_tables": 25,
  "available_tables": 15,
  "occupied_tables": 8,
  "reserved_tables": 2,
  "occupancy_rate": 40.0,
  "average_turn_time_minutes": 45
}
```

---

## 5. ZONES MODULE ✅ (Priority 2)

**Base Path:** `/api/v1/organizations/{org_id}/zones`

### 5.1 Get Zones
```http
GET /api/v1/organizations/{org_id}/zones
```

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Main Hall",
      "description": "Main dining area",
      "is_active": true,
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 5.2 Get Zone by ID
```http
GET /api/v1/organizations/{org_id}/zones/{zone_id}
```

**Response:** 200 OK

### 5.3 Get Table Count by Zone
```http
GET /api/v1/organizations/{org_id}/tables/count-by-zone
```

**Response:** 200 OK
```json
{
  "uuid1": 10,
  "uuid2": 8,
  "uuid3": 7
}
```

---

## 6. PAYMENTS MODULE ✅ (CRITICAL - Priority 1)

**Base Path:** `/api/v1/organizations/{org_id}/payments`

### 6.1 Process Payment
```http
POST /api/v1/organizations/{org_id}/payments
```

**Request Body:**
```json
{
  "order_id": "uuid",
  "amount": 43.49,
  "payment_method": "cash",
  "status": "completed",
  "transaction_id": "TXN-12345",
  "card_last_four": null,
  "card_type": null,
  "payment_processor": "internal",
  "processed_at": "2025-11-13T10:00:00Z",
  "metadata": {}
}
```

**Response:** 201 Created
```json
{
  "id": "uuid",
  "order_id": "uuid",
  "amount": 43.49,
  "payment_method": "cash",
  "status": "completed",
  ...
}
```

### 6.2 Get Payment by ID
```http
GET /api/v1/organizations/{org_id}/payments/{payment_id}
```

**Response:** 200 OK

### 6.3 Get Payments by Order
```http
GET /api/v1/organizations/{org_id}/payments?order_id={uuid}
```

**Response:** 200 OK (Array of payments)

### 6.4 Get Payments (List with Filters)
```http
GET /api/v1/organizations/{org_id}/payments?method={method}&status={status}&from_date={iso8601}&to_date={iso8601}
```

**Query Parameters:**
- `method` (optional): cash, card, digitalWallet, giftCard, check, other
- `status` (optional): pending, processing, completed, failed, refunded
- `from_date`, `to_date` (optional): ISO 8601 datetime

**Response:** 200 OK

### 6.5 Cancel Payment
```http
POST /api/v1/organizations/{org_id}/payments/{payment_id}/cancel
```

**Request Body:**
```json
{
  "reason": "Customer requested cancellation"
}
```

**Response:** 200 OK

### 6.6 Get Payment Statistics
```http
GET /api/v1/organizations/{org_id}/payments/statistics?from_date={iso8601}&to_date={iso8601}
```

**Response:** 200 OK
```json
{
  "total_amount": 125000.50,
  "payment_count": 850,
  "average_payment": 147.06,
  "payment_method_breakdown": {
    "cash": 45000.00,
    "card": 65000.00,
    "digitalWallet": 15000.50
  }
}
```

---

## 7. REFUNDS MODULE ✅ (Priority 1)

**Base Path:** `/api/v1/organizations/{org_id}/refunds`

### 7.1 Process Refund
```http
POST /api/v1/organizations/{org_id}/refunds
```

**Request Body:**
```json
{
  "payment_id": "uuid",
  "order_id": "uuid",
  "amount": 43.49,
  "reason": "customerRequest",
  "status": "pending",
  "requested_by": "user_uuid",
  "requested_at": "2025-11-13T10:00:00Z",
  "notes": "Customer not satisfied"
}
```

**Response:** 201 Created
```json
{
  "id": "uuid",
  "payment_id": "uuid",
  "order_id": "uuid",
  "amount": 43.49,
  "reason": "customerRequest",
  "status": "completed",
  "requested_by": "user_uuid",
  "requested_at": "2025-11-13T10:00:00Z",
  "processed_at": "2025-11-13T10:01:00Z",
  "processed_by": "user_uuid",
  "notes": "Customer not satisfied"
}
```

### 7.2 Get Refund by ID
```http
GET /api/v1/organizations/{org_id}/refunds/{refund_id}
```

**Response:** 200 OK

### 7.3 Get Refunds by Payment
```http
GET /api/v1/organizations/{org_id}/refunds?payment_id={uuid}
```

**Response:** 200 OK (Array of refunds)

### 7.4 Get Refunds by Order
```http
GET /api/v1/organizations/{org_id}/refunds?order_id={uuid}
```

**Response:** 200 OK (Array of refunds)

---

## 8. STAFF MODULE ✅ (Priority 2)

**Base Path:** `/api/v1/organizations/{org_id}/staff`

### 8.1 Get Staff Members
```http
GET /api/v1/organizations/{org_id}/staff?role={role}&status={status}
```

**Query Parameters:**
- `role` (optional): admin, manager, cashier, waiter, kitchen, bartender
- `status` (optional): active, inactive, onLeave, terminated

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "role": "cashier",
      "status": "active",
      "hire_date": "2024-01-15",
      "hourly_rate": 15.50,
      "avatar_url": "https://...",
      "address": "123 Main St",
      "emergency_contact": "+0987654321",
      "permissions": ["view_orders", "create_orders", "process_payments"],
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 8.2 Get Staff Member by ID
```http
GET /api/v1/organizations/{org_id}/staff/{staff_id}
```

**Response:** 200 OK

### 8.3 Create Staff Member
```http
POST /api/v1/organizations/{org_id}/staff
```

**Request Body:** (Same structure as Get Staff response)

**Response:** 201 Created

### 8.4 Update Staff Member
```http
PATCH /api/v1/organizations/{org_id}/staff/{staff_id}
```

**Response:** 200 OK

### 8.5 Delete Staff Member
```http
DELETE /api/v1/organizations/{org_id}/staff/{staff_id}
```

**Response:** 204 No Content

---

## 9. ROLES MODULE ✅ (Priority 2)

**Base Path:** `/api/v1/organizations/{org_id}/roles`

### 9.1 Get Roles
```http
GET /api/v1/organizations/{org_id}/roles
```

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Cashier",
      "permissions": [
        "view_orders",
        "create_orders",
        "process_payments"
      ],
      "is_custom": false,
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 9.2 Get Role by ID
```http
GET /api/v1/organizations/{org_id}/roles/{role_id}
```

**Response:** 200 OK

### 9.3 Create Role
```http
POST /api/v1/organizations/{org_id}/roles
```

**Request Body:**
```json
{
  "name": "Custom Role",
  "permissions": ["view_orders", "create_orders"],
  "is_custom": true
}
```

**Response:** 201 Created

### 9.4 Update Role
```http
PATCH /api/v1/organizations/{org_id}/roles/{role_id}
```

**Response:** 200 OK

### 9.5 Delete Role
```http
DELETE /api/v1/organizations/{org_id}/roles/{role_id}
```

**Response:** 204 No Content

---

## 10. SHIFTS MODULE ✅ (Priority 2)

**Base Path:** `/api/v1/organizations/{org_id}/shifts`

### 10.1 Get Shifts
```http
GET /api/v1/organizations/{org_id}/shifts?staff_id={uuid}&date={date}&from_date={date}&to_date={date}
```

**Query Parameters:**
- `staff_id` (optional): Filter by staff member
- `date` (optional): Filter by specific date (YYYY-MM-DD)
- `from_date`, `to_date` (optional): Date range

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "staff_id": "uuid",
      "staff_name": "John Doe",
      "date": "2025-11-13",
      "start_time": "09:00",
      "end_time": "17:00",
      "position": "Cashier",
      "is_confirmed": true,
      "notes": "Morning shift",
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 10.2 Get Shift by ID
```http
GET /api/v1/organizations/{org_id}/shifts/{shift_id}
```

**Response:** 200 OK

### 10.3 Create Shift
```http
POST /api/v1/organizations/{org_id}/shifts
```

**Request Body:** (Same structure as Get Shifts response)

**Response:** 201 Created

### 10.4 Update Shift
```http
PATCH /api/v1/organizations/{org_id}/shifts/{shift_id}
```

**Response:** 200 OK

### 10.5 Delete Shift
```http
DELETE /api/v1/organizations/{org_id}/shifts/{shift_id}
```

**Response:** 204 No Content

---

## 11. ACTIVITY LOGS MODULE ✅ (Priority 3)

**Base Path:** `/api/v1/organizations/{org_id}/activity-logs`

### 11.1 Get Activity Logs
```http
GET /api/v1/organizations/{org_id}/activity-logs?staff_id={uuid}&type={type}&from_date={iso8601}&to_date={iso8601}
```

**Query Parameters:**
- `staff_id` (optional): Filter by staff member
- `type` (optional): clockIn, clockOut, orderCreated, paymentProcessed, refundProcessed, other
- `from_date`, `to_date` (optional): ISO 8601 datetime

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "staff_id": "uuid",
      "staff_name": "John Doe",
      "action": "Clocked in",
      "type": "clockIn",
      "timestamp": "2025-11-13T10:00:00Z",
      "details": "Started morning shift"
    }
  ]
}
```

---

## 12. KITCHEN STATIONS MODULE ✅ (Priority 2)

**Base Path:** `/api/v1/organizations/{org_id}/kitchen-stations`

### 12.1 Get Kitchen Stations
```http
GET /api/v1/organizations/{org_id}/kitchen-stations
```

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Grill Station",
      "type": "grill",
      "is_active": true,
      "order_position": 1,
      "display_name": "GRILL",
      "printer_name": "EPSON-KITCHEN-01",
      "assigned_staff": ["uuid1", "uuid2"],
      "active_orders": 5,
      "notes": "Main grill station",
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 12.2 Get Kitchen Station by ID
```http
GET /api/v1/organizations/{org_id}/kitchen-stations/{station_id}
```

**Response:** 200 OK

### 12.3 Create Kitchen Station
```http
POST /api/v1/organizations/{org_id}/kitchen-stations
```

**Request Body:** (Same structure as Get Kitchen Stations response)

**Response:** 201 Created

### 12.4 Update Kitchen Station
```http
PATCH /api/v1/organizations/{org_id}/kitchen-stations/{station_id}
```

**Response:** 200 OK

### 12.5 Delete Kitchen Station
```http
DELETE /api/v1/organizations/{org_id}/kitchen-stations/{station_id}
```

**Response:** 204 No Content

---

## 13. CUSTOMERS MODULE ✅ (Priority 3)

**Base Path:** `/api/v1/organizations/{org_id}/customers`

### 13.1 Get Customers
```http
GET /api/v1/organizations/{org_id}/customers?search={query}
```

**Query Parameters:**
- `search` (optional): Search by name, email, or phone

**Response:** 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Jane Smith",
      "email": "jane@example.com",
      "phone": "+1234567890",
      "address": "456 Oak Ave",
      "loyalty_points": 150,
      "total_orders": 25,
      "total_spent": 850.50,
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ]
}
```

### 13.2 Get Customer by ID
```http
GET /api/v1/organizations/{org_id}/customers/{customer_id}
```

**Response:** 200 OK

### 13.3 Create Customer
```http
POST /api/v1/organizations/{org_id}/customers
```

**Request Body:** (Same structure as Get Customers response)

**Response:** 201 Created

### 13.4 Update Customer
```http
PATCH /api/v1/organizations/{org_id}/customers/{customer_id}
```

**Response:** 200 OK

### 13.5 Delete Customer
```http
DELETE /api/v1/organizations/{org_id}/customers/{customer_id}
```

**Response:** 204 No Content

---

## 📝 Implementation Notes

### Priority Levels

**Priority 1 (CRITICAL - Week 1):**
- Orders Module (Complete)
- Products Module (Complete)
- Categories Module (Complete)
- Tables Module (Complete)
- Payments Module (Complete)
- Refunds Module (Complete)

**Priority 2 (IMPORTANT - Week 2):**
- Staff Module (Complete)
- Roles Module (Complete)
- Shifts Module (Complete)
- Kitchen Stations Module (Complete)
- Zones Module (Complete)

**Priority 3 (NICE TO HAVE - Week 3):**
- Activity Logs Module (Read-only)
- Customers Module (Complete)

### Multi-Tenancy

All endpoints are organization-scoped using `{org_id}` in the path. The backend must:
1. Extract `organization_id` from JWT token
2. Validate that the user has access to the requested organization
3. Apply Row-Level Security (RLS) to ensure data isolation
4. Set PostgreSQL session variable: `SET app.current_organization_id = '{org_id}'`

### Response Format Standards

**Success Responses:**
- 200 OK: Successful GET/PATCH/POST (update)
- 201 Created: Successful POST (create)
- 204 No Content: Successful DELETE

**Error Responses:**
```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Order not found",
    "details": {
      "order_id": "uuid"
    }
  }
}
```

**Error Codes:**
- 400 Bad Request: Invalid input
- 401 Unauthorized: Missing or invalid JWT
- 403 Forbidden: Insufficient permissions
- 404 Not Found: Resource doesn't exist
- 409 Conflict: Duplicate resource
- 422 Unprocessable Entity: Validation failed
- 500 Internal Server Error: Server error

### Data Format Conventions

**Dates/Times:**
- Use ISO 8601 format: `2025-11-13T10:00:00Z`
- All times in UTC
- Date only: `2025-11-13`
- Time only: `14:30` (24-hour format)

**Field Naming:**
- API uses snake_case: `order_number`, `created_at`
- Flutter uses camelCase: `orderNumber`, `createdAt`
- Automatic conversion handled by client

**Enums:**
- Sent as strings in lowercase: `"pending"`, `"completed"`
- Case-insensitive matching recommended

### Pagination

For list endpoints, use query parameters:
```
?page=1&per_page=50&sort_by=created_at&sort_order=desc
```

Default pagination:
- `page`: 1
- `per_page`: 50
- `max_per_page`: 100

Response structure:
```json
{
  "data": [...],
  "pagination": {
    "total": 150,
    "page": 1,
    "per_page": 50,
    "total_pages": 3
  }
}
```

---

## 🚀 Quick Start for Backend Developer

### Step 1: Set Up Database
```bash
# Run migrations in order
cd /backend/postgres/migrations
psql -d pos_saas -f V001__core_tenant_tables.sql
psql -d pos_saas -f V002__pos_core_tables.sql
# ... continue through V016
```

### Step 2: Implement Authentication
```go
// Use existing middleware
// File: backend/internal/auth/middleware.go
```

### Step 3: Generate CRUD Endpoints
```bash
# Use code generation tool
cd /backend/tools/crud-generator
./generate.sh orders
./generate.sh products
# Follow: backend/CODE_GENERATION_GUIDE.md
```

### Step 4: Test with Flutter App
```dart
// In Flutter app
final config = AppConfig(
  apiBaseUrl: 'http://localhost:8080/api/v1',
  // ... other config
);

final context = AppContext(
  tenantId: 'your-org-uuid',
  // ... other context
);
```

---

## 📚 Reference Implementation

The backend already has a **complete Products module** that can be used as a reference:

**Location:** `/backend/internal/product/`

**Files:**
- `handler.go` - HTTP handlers
- `service.go` - Business logic
- `repository.go` - Database queries
- `dto.go` - Data transfer objects
- `validator.go` - Input validation
- `routes.go` - Route definitions

**Copy this pattern** for all other modules!

---

## ✅ Checklist for Backend Developer

- [ ] Review Flutter-Database backend structure
- [ ] Set up PostgreSQL database with migrations
- [ ] Configure JWT authentication
- [ ] Implement Priority 1 endpoints (Orders, Products, Tables, Payments)
- [ ] Test with Postman/curl
- [ ] Configure CORS for Flutter app
- [ ] Implement Priority 2 endpoints (Staff, Roles, Shifts, Kitchen Stations)
- [ ] Add Rate Limiting (already configured in backend)
- [ ] Set up monitoring (Prometheus/Grafana ready)
- [ ] Deploy to staging environment
- [ ] Test with Flutter app integration

---

## 📧 Contact

For questions or clarifications, please provide:
1. Specific endpoint name
2. Expected behavior
3. Current error or issue

**Flutter Client Code Location:**
- `/packages/pos_core/lib/src/features/*/data/sources/*_remote_source.dart`

All HTTP implementations are ready and waiting for these endpoints! 🚀

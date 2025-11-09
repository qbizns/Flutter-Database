# POS Database System

A comprehensive Point of Sale (POS) database system designed for Flutter applications.

## Overview

This repository contains the database schema, models, and data management layer for a Point of Sale system. It provides a robust foundation for managing retail operations including sales transactions, inventory, customers, and reporting.

## Features

- **Sales Management**: Track sales transactions, receipts, and payment methods
- **Inventory Management**: Manage products, stock levels, and categories
- **Customer Management**: Store customer information and purchase history
- **User & Authentication**: Multi-user support with role-based access
- **Reporting**: Sales reports, inventory tracking, and analytics
- **Offline-First**: Designed to work seamlessly offline with sync capabilities

## Database Structure

The system includes the following core components:

### Core Entities
- **Products**: Item catalog with pricing, categories, and stock information
- **Sales**: Transaction records with line items and payment details
- **Customers**: Customer profiles and contact information
- **Inventory**: Stock management and tracking
- **Users**: System users with roles and permissions
- **Categories**: Product categorization and organization

## Technology Stack

- **Database**: SQLite (local) / PostgreSQL (server)
- **ORM/Query Builder**: To be determined based on implementation
- **Platform**: Flutter/Dart compatible
- **Sync**: Offline-first with background synchronization

## Getting Started

### Prerequisites

- Flutter SDK (latest stable version)
- Dart SDK
- Database client (for development and testing)

### Installation

```bash
# Clone the repository
git clone https://github.com/Macber-eg/Flutter-Database.git

# Navigate to project directory
cd Flutter-Database

# Install dependencies (once Flutter project is set up)
flutter pub get
```

## PostgreSQL Database

The `postgres/` directory contains a production-ready PostgreSQL database schema for the POS SAAS system.

### Quick Start

```bash
# Navigate to postgres directory
cd postgres

# Run the complete setup (creates database, runs migrations, seeds data)
./scripts/run_all.sh pos_saas postgres

# Or run individual steps:
psql -U postgres -d pos_saas -f scripts/init_database.sql
psql -U postgres -d pos_saas -f migrations/V001_20251109_create_core_tenant_tables.sql
psql -U postgres -d pos_saas -f migrations/V002_20251109_create_pos_core_tables.sql
psql -U postgres -d pos_saas -f migrations/V003_20251109_implement_row_level_security.sql
psql -U postgres -d pos_saas -f migrations/V004_20251109_create_additional_pos_tables.sql
psql -U postgres -d pos_saas -f migrations/V005_20251109_create_inter_location_transfers.sql
psql -U postgres -d pos_saas -f migrations/V006_20251109_create_advanced_inventory_management.sql
psql -U postgres -d pos_saas -f migrations/V007_20251109_create_enhanced_loyalty_program.sql
psql -U postgres -d pos_saas -f migrations/V008_20251109_create_reporting_analytics.sql
psql -U postgres -d pos_saas -f seed_data/001_seed_core_data.sql
psql -U postgres -d pos_saas -f seed_data/002_seed_pos_data.sql
psql -U postgres -d pos_saas -f seed_data/003_seed_additional_pos_data.sql
psql -U postgres -d pos_saas -f seed_data/004_seed_enhanced_features.sql
```

### Database Structure

#### Core System Tables (Multi-Tenancy)
- **organizations** - Tenant/company management with subscription plans
- **users** - System users with organization mapping
- **roles** - RBAC role definitions (system and custom)
- **permissions** - Granular permission system
- **role_permissions** - Role-permission mapping
- **user_roles** - User-role assignments
- **audit_logs** - Complete system audit trail

#### POS Tables
- **categories** - Hierarchical product categorization
- **products** - Product catalog with pricing, inventory, and variants
- **product_variants** - Product variations (size, color, etc.)
- **customers** - Customer master with loyalty and credit management
- **sales** - Sales transaction headers
- **sale_items** - Transaction line items
- **payments** - Payment records with multiple payment methods
- **inventory_transactions** - Complete inventory movement audit trail
- **suppliers** - Supplier/vendor master data
- **purchase_orders** - Purchase orders for inventory procurement
- **purchase_order_items** - Purchase order line items
- **locations** - Store locations and branches for multi-location support
- **promotions** - Promotions and discount campaigns
- **promotion_usage** - Promotion usage tracking
- **expenses** - Business expense tracking
- **shifts** - Cashier shifts and cash register management

#### Advanced Inventory (V006)
- **product_serial_numbers** - Serial number tracking for individual items
- **product_batches** - Batch/lot tracking with expiration management
- **batch_transactions** - Batch transaction audit trail
- **cycle_counts** - Physical inventory cycle counting
- **cycle_count_items** - Individual cycle count line items
- **stock_adjustment_reasons** - Pre-defined inventory adjustment reasons

#### Inter-Location Transfers (V005)
- **inventory_transfers** - Transfer orders between locations
- **inventory_transfer_items** - Transfer line items with variance tracking

#### Enhanced Loyalty Program (V007)
- **loyalty_tiers** - Customer loyalty tiers (Bronze, Silver, Gold, Platinum)
- **loyalty_tier_benefits** - Benefits for each loyalty tier
- **loyalty_points_rules** - Rules for earning loyalty points
- **loyalty_rewards** - Catalog of rewards for redemption
- **loyalty_redemptions** - Customer reward redemptions
- **loyalty_points_transactions** - Complete points transaction history
- **customer_tier_history** - Customer tier change tracking

#### Reporting & Analytics (V008)
- **mv_daily_sales_summary** - Daily sales aggregation (materialized view)
- **mv_product_performance** - Product performance metrics (materialized view)
- **mv_customer_analytics** - Customer behavior analytics with RFM (materialized view)
- **mv_inventory_valuation** - Current inventory valuation (materialized view)
- **mv_location_performance** - Location-based metrics (materialized view)
- **mv_promotion_effectiveness** - Promotion ROI analysis (materialized view)

### Key Features

✅ **Multi-Tenant SAAS Architecture**
  - Organization-based data isolation
  - Subscription and quota management
  - Flexible per-tenant schemas

✅ **Comprehensive RBAC**
  - Role-based access control
  - Granular permissions
  - System and custom roles

✅ **Multi-Location Support**
  - Store locations and branches
  - Location-specific inventory and sales
  - Warehouse management
  - Location-based reporting

✅ **Supplier & Procurement Management**
  - Supplier master data
  - Purchase orders with approval workflow
  - Partial receives support
  - Purchase tracking and history

✅ **Product Variants**
  - Size, color, and custom variations
  - Variant-specific pricing and inventory
  - Flexible attributes (JSONB)
  - SKU and barcode per variant

✅ **Promotions & Discounts**
  - Multiple promotion types (percentage, fixed, BOGO)
  - Product/category targeting
  - Usage limits and tracking
  - Date-based campaigns

✅ **Expense Tracking**
  - Business expense management
  - Category classification
  - Approval workflow
  - Receipt attachments

✅ **Shift Management**
  - Cashier shift tracking
  - Cash reconciliation
  - Sales summary per shift
  - Payment method breakdown

✅ **Inter-Location Transfers** (V005)
  - Transfer inventory between locations
  - Complete workflow (draft → approved → shipped → received)
  - Variance tracking
  - Carrier and tracking number support

✅ **Advanced Inventory Management** (V006)
  - Serial number tracking for warranty and traceability
  - Batch/lot tracking with expiration dates
  - Cycle counting for physical inventory verification
  - Pre-defined stock adjustment reasons
  - Quality control status tracking

✅ **Enhanced Loyalty Program** (V007)
  - Multi-tier loyalty system (Bronze, Silver, Gold, Platinum)
  - Flexible points earning rules (purchase, signup, birthday, referral)
  - Rewards catalog with redemption tracking
  - Tier-based benefits and discounts
  - Complete points transaction history
  - Tier upgrade/downgrade tracking

✅ **Reporting & Analytics** (V008)
  - Materialized views for fast reporting
  - Daily sales summaries
  - Product performance analysis
  - Customer analytics with RFM segmentation
  - Inventory valuation reports
  - Location performance metrics
  - Promotion effectiveness tracking

✅ **Flexible Schema Design**
  - JSONB fields for extensibility
  - Soft deletes for data preservation
  - Audit trails on all tables

✅ **Production-Ready**
  - UUID primary keys
  - Comprehensive Row-Level Security (RLS) on all tables
  - Proper indexing strategy
  - 8 migration files with complete workflow
  - Comprehensive seed data for testing
  - Automated analytics refresh functions

### Documentation

- [Database README](postgres/README.md) - Setup and migration guide
- [Schema Documentation](postgres/schemas/SCHEMA_DOCUMENTATION.md) - Complete schema reference
- [Enhanced Features Summary](postgres/schemas/ENHANCED_FEATURES_SUMMARY.md) - V005-V008 features guide

### Test Data

The seed data includes:

**Core Data:**
- 3 organizations (Demo Retail Store, Coffee Corner, Tech Gadgets Pro)
- 7 users with different roles
- 5 system roles + custom roles
- 23 granular permissions
- 16 sample products across categories
- 12 product variants (colors, sizes)
- 10 customers with loyalty data

**Operations Data:**
- 8 complete sales transactions with payments
- 6 suppliers across organizations
- 6 store locations/branches
- 7 promotional campaigns
- 4 purchase orders with line items
- 10 business expenses
- 6 cashier shifts with reconciliation

**Enhanced Features (V005-V007):**
- 3 inter-location transfers with line items
- 8 product serial numbers (phones, laptops)
- 4 product batches with expiration tracking
- 2 cycle counts with variance analysis
- 7 loyalty tiers across organizations
- 7 loyalty rewards in catalog
- 4 reward redemptions
- 10+ loyalty points transactions
- 10 system-defined stock adjustment reasons

**Analytics (V008):**
- 6 pre-configured materialized views
- 3 helper views for common reports
- Automated refresh functions

**Total Database Size:** 38 tables + 9 views + 200+ seed records

Test credentials: `admin@demoretail.com`, `manager@demoretail.com`, `cashier1@demoretail.com`

## Database Schema

See [postgres/schemas/SCHEMA_DOCUMENTATION.md](postgres/schemas/SCHEMA_DOCUMENTATION.md) for complete schema documentation.

## Development Roadmap

### Phase 1: Database Layer ✅ COMPLETED
- [x] Define database schema
- [x] Create core migrations (V001-V002)
- [x] Implement Row-Level Security (V003)
- [x] Add essential POS tables (V004)
- [x] Implement inter-location transfers (V005)
- [x] Add advanced inventory management (V006)
- [x] Create enhanced loyalty program (V007)
- [x] Build reporting & analytics infrastructure (V008)
- [x] Add comprehensive seed data
- [x] Full documentation

### Phase 2: Application Layer (Next Steps)
- [ ] API Development
  - [ ] REST API with authentication
  - [ ] CRUD endpoints for all entities
  - [ ] Real-time updates (WebSocket)
  - [ ] API documentation (OpenAPI/Swagger)
- [ ] Flutter/Dart Integration
  - [ ] Generate Dart models from schema
  - [ ] Build data access layer
  - [ ] Implement offline sync
  - [ ] SQLite local database
- [ ] Operational Tools
  - [ ] Backup and restore scripts
  - [ ] Database utilities
  - [ ] Migration rollback scripts
  - [ ] Health monitoring

### Phase 3: Advanced Features (Future)
- [ ] Payment gateway integration
- [ ] E-commerce platform integration
- [ ] Email/SMS notifications
- [ ] Advanced reporting dashboards
- [ ] Mobile app features
- [ ] Third-party integrations

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For questions or support, please open an issue in the repository.

---

**Status**: 🚧 In Development

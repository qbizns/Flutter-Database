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
psql -U postgres -d pos_saas -f seed_data/001_seed_core_data.sql
psql -U postgres -d pos_saas -f seed_data/002_seed_pos_data.sql
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
- **customers** - Customer master with loyalty and credit management
- **sales** - Sales transaction headers
- **sale_items** - Transaction line items
- **payments** - Payment records with multiple payment methods
- **inventory_transactions** - Complete inventory movement audit trail

### Key Features

✅ **Multi-Tenant SAAS Architecture**
  - Organization-based data isolation
  - Subscription and quota management
  - Flexible per-tenant schemas

✅ **Comprehensive RBAC**
  - Role-based access control
  - Granular permissions
  - System and custom roles

✅ **Flexible Schema Design**
  - JSONB fields for extensibility
  - Soft deletes for data preservation
  - Audit trails on all tables

✅ **Production-Ready**
  - UUID primary keys
  - Proper indexing strategy
  - Migration system
  - Seed data for testing

### Documentation

- [Database README](postgres/README.md) - Setup and migration guide
- [Schema Documentation](postgres/schemas/SCHEMA_DOCUMENTATION.md) - Complete schema reference

### Test Data

The seed data includes:
- 3 organizations (Demo Retail Store, Coffee Corner, Tech Gadgets Pro)
- 7 users with different roles
- 5 system roles + custom roles
- 23 granular permissions
- 16 sample products across categories
- 10 customers
- 8 complete sales transactions with payments

Test credentials: `admin@demoretail.com`, `manager@demoretail.com`, `cashier1@demoretail.com`

## Database Schema

See [postgres/schemas/SCHEMA_DOCUMENTATION.md](postgres/schemas/SCHEMA_DOCUMENTATION.md) for complete schema documentation.

## Development Roadmap

- [x] Define database schema
- [x] Create database migrations
- [x] Add seed data for testing
- [ ] Implement data models (Dart)
- [ ] Build CRUD operations
- [ ] Add data validation
- [ ] Implement sync functionality
- [ ] Add backup and restore features
- [ ] Performance optimization
- [ ] API layer development

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For questions or support, please open an issue in the repository.

---

**Status**: 🚧 In Development

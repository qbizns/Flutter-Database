-- ============================================================================
-- Migration: V002 - Create POS Core Tables
-- Description: Creates essential POS tables for products, sales, customers, and inventory
-- Note: These tables are created in the public schema as templates
--       In production, each tenant will have these in their own schema
-- Author: System
-- Date: 2025-11-09
-- ============================================================================

BEGIN;

-- ============================================================================
-- CATEGORIES TABLE
-- Description: Product categorization hierarchy
-- ============================================================================

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Category Information
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,

    -- Hierarchy
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    level INTEGER DEFAULT 0, -- For hierarchy depth
    path VARCHAR(500), -- e.g., 'electronics/phones/smartphones'

    -- Display
    image_url TEXT,
    icon VARCHAR(100),
    color VARCHAR(20), -- Hex color code
    sort_order INTEGER DEFAULT 0,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Settings
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id, slug)
);

-- Indexes
CREATE INDEX idx_categories_organization_id ON categories(organization_id);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_is_active ON categories(is_active) WHERE deleted_at IS NULL;

-- Trigger
CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE categories IS 'Product categories with hierarchical support';
COMMENT ON COLUMN categories.path IS 'Full category path for easy hierarchy navigation';

-- ============================================================================
-- PRODUCTS TABLE
-- Description: Product/Item catalog
-- ============================================================================

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Product Identification
    sku VARCHAR(100), -- Stock Keeping Unit
    barcode VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    description TEXT,

    -- Categorization
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,

    -- Pricing
    cost_price NUMERIC(12, 2) DEFAULT 0, -- What you pay
    selling_price NUMERIC(12, 2) NOT NULL, -- What customer pays
    compare_at_price NUMERIC(12, 2), -- Original price (for showing discounts)

    -- Tax
    tax_rate NUMERIC(5, 2) DEFAULT 0, -- Percentage
    is_tax_inclusive BOOLEAN DEFAULT FALSE,

    -- Inventory Tracking
    track_inventory BOOLEAN DEFAULT TRUE,
    current_stock NUMERIC(12, 2) DEFAULT 0,
    low_stock_threshold NUMERIC(12, 2) DEFAULT 0,
    unit VARCHAR(50) DEFAULT 'unit', -- e.g., 'unit', 'kg', 'liter'

    -- Product Type
    is_service BOOLEAN DEFAULT FALSE,
    is_composite BOOLEAN DEFAULT FALSE, -- Bundle/combo product
    has_variants BOOLEAN DEFAULT FALSE, -- e.g., different sizes/colors

    -- Display
    image_url TEXT,
    images JSONB DEFAULT '[]', -- Array of image URLs
    sort_order INTEGER DEFAULT 0,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,

    -- Custom Fields (extensibility)
    custom_fields JSONB DEFAULT '{}',

    -- Settings
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id, sku),
    UNIQUE(organization_id, barcode),
    CONSTRAINT positive_prices CHECK (selling_price >= 0 AND cost_price >= 0)
);

-- Indexes
CREATE INDEX idx_products_organization_id ON products(organization_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_name ON products USING gin(to_tsvector('english', name));
CREATE INDEX idx_products_is_active ON products(is_active) WHERE deleted_at IS NULL;

-- Trigger
CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE products IS 'Product catalog with pricing and inventory tracking';
COMMENT ON COLUMN products.custom_fields IS 'Extensible JSON field for custom product attributes';

-- ============================================================================
-- CUSTOMERS TABLE
-- Description: Customer information
-- ============================================================================

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Customer Identification
    customer_code VARCHAR(50), -- Internal customer code

    -- Personal Information
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    full_name VARCHAR(255) GENERATED ALWAYS AS (
        CASE
            WHEN first_name IS NOT NULL AND last_name IS NOT NULL THEN first_name || ' ' || last_name
            WHEN first_name IS NOT NULL THEN first_name
            ELSE last_name
        END
    ) STORED,
    company_name VARCHAR(255),

    -- Contact Information
    email VARCHAR(255),
    phone VARCHAR(50),
    alternate_phone VARCHAR(50),

    -- Address
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),

    -- Customer Details
    date_of_birth DATE,
    gender VARCHAR(20),

    -- Business Information
    tax_number VARCHAR(100), -- VAT/TIN number

    -- Loyalty
    loyalty_points INTEGER DEFAULT 0,
    loyalty_tier VARCHAR(50),

    -- Credit
    credit_limit NUMERIC(12, 2) DEFAULT 0,
    outstanding_balance NUMERIC(12, 2) DEFAULT 0,

    -- Stats (can be calculated or cached)
    total_purchases NUMERIC(12, 2) DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    last_purchase_at TIMESTAMP WITH TIME ZONE,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Notes
    notes TEXT,

    -- Custom Fields
    custom_fields JSONB DEFAULT '{}',

    -- Settings
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id, customer_code),
    UNIQUE(organization_id, email) WHERE email IS NOT NULL
);

-- Indexes
CREATE INDEX idx_customers_organization_id ON customers(organization_id);
CREATE INDEX idx_customers_customer_code ON customers(customer_code);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_full_name ON customers(full_name);
CREATE INDEX idx_customers_is_active ON customers(is_active) WHERE deleted_at IS NULL;

-- Trigger
CREATE TRIGGER update_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE customers IS 'Customer master data with contact and loyalty information';

-- ============================================================================
-- SALES TABLE
-- Description: Sales transactions/receipts
-- ============================================================================

CREATE TABLE IF NOT EXISTS sales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Transaction Identification
    sale_number VARCHAR(50) NOT NULL, -- Receipt/Invoice number
    reference_number VARCHAR(100), -- External reference

    -- Transaction Type
    transaction_type transaction_type DEFAULT 'sale' NOT NULL,

    -- Customer
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,

    -- Cashier/Salesperson
    cashier_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Amounts
    subtotal NUMERIC(12, 2) DEFAULT 0 NOT NULL,
    tax_amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,
    discount_amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL,

    -- Payments
    paid_amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,
    change_amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,
    outstanding_amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,

    -- Payment Status
    payment_status payment_status DEFAULT 'pending' NOT NULL,

    -- Discount
    discount_type VARCHAR(20), -- 'percentage', 'fixed'
    discount_value NUMERIC(12, 2) DEFAULT 0,
    discount_reason TEXT,

    -- Transaction Details
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Notes
    notes TEXT,
    internal_notes TEXT, -- Not shown to customer

    -- Custom Fields
    custom_fields JSONB DEFAULT '{}',

    -- Settings
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id, sale_number),
    CONSTRAINT valid_amounts CHECK (total_amount >= 0)
);

-- Indexes
CREATE INDEX idx_sales_organization_id ON sales(organization_id);
CREATE INDEX idx_sales_sale_number ON sales(sale_number);
CREATE INDEX idx_sales_customer_id ON sales(customer_id);
CREATE INDEX idx_sales_cashier_id ON sales(cashier_id);
CREATE INDEX idx_sales_transaction_date ON sales(transaction_date);
CREATE INDEX idx_sales_payment_status ON sales(payment_status);
CREATE INDEX idx_sales_created_at ON sales(created_at);

-- Trigger
CREATE TRIGGER update_sales_updated_at
    BEFORE UPDATE ON sales
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE sales IS 'Sales transactions and receipts';
COMMENT ON COLUMN sales.outstanding_amount IS 'Remaining amount to be paid';

-- ============================================================================
-- SALE_ITEMS TABLE
-- Description: Line items for each sale
-- ============================================================================

CREATE TABLE IF NOT EXISTS sale_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Sale Reference
    sale_id UUID NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Product Reference
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    product_name VARCHAR(255) NOT NULL, -- Snapshot at time of sale
    product_sku VARCHAR(100),

    -- Quantity
    quantity NUMERIC(12, 3) NOT NULL,
    unit VARCHAR(50) DEFAULT 'unit',

    -- Pricing (snapshot at time of sale)
    unit_price NUMERIC(12, 2) NOT NULL,
    cost_price NUMERIC(12, 2) DEFAULT 0,

    -- Calculations
    subtotal NUMERIC(12, 2) NOT NULL, -- quantity * unit_price
    tax_rate NUMERIC(5, 2) DEFAULT 0,
    tax_amount NUMERIC(12, 2) DEFAULT 0,
    discount_amount NUMERIC(12, 2) DEFAULT 0,
    total NUMERIC(12, 2) NOT NULL,

    -- Discount
    discount_type VARCHAR(20), -- 'percentage', 'fixed'
    discount_value NUMERIC(12, 2) DEFAULT 0,

    -- Item Notes
    notes TEXT,

    -- Custom Fields
    custom_fields JSONB DEFAULT '{}',

    -- Settings
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Constraints
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT positive_prices CHECK (unit_price >= 0 AND cost_price >= 0)
);

-- Indexes
CREATE INDEX idx_sale_items_sale_id ON sale_items(sale_id);
CREATE INDEX idx_sale_items_product_id ON sale_items(product_id);
CREATE INDEX idx_sale_items_organization_id ON sale_items(organization_id);

-- Trigger
CREATE TRIGGER update_sale_items_updated_at
    BEFORE UPDATE ON sale_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE sale_items IS 'Line items for sales transactions';
COMMENT ON COLUMN sale_items.product_name IS 'Product name snapshot at time of sale';

-- ============================================================================
-- PAYMENTS TABLE
-- Description: Payment records for sales
-- ============================================================================

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Sale Reference
    sale_id UUID NOT NULL REFERENCES sales(id) ON DELETE CASCADE,

    -- Payment Details
    payment_method payment_method NOT NULL,
    payment_status payment_status DEFAULT 'pending' NOT NULL,

    amount NUMERIC(12, 2) NOT NULL,

    -- Payment Method Specific Details
    card_last_four VARCHAR(4),
    card_type VARCHAR(50), -- visa, mastercard, etc.
    transaction_id VARCHAR(255), -- External payment gateway transaction ID
    reference_number VARCHAR(255),

    -- Mobile Money / Bank Transfer
    account_number VARCHAR(100),
    account_name VARCHAR(255),

    -- Payment Date
    payment_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE,

    -- Notes
    notes TEXT,

    -- Custom Fields
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT positive_amount CHECK (amount > 0)
);

-- Indexes
CREATE INDEX idx_payments_organization_id ON payments(organization_id);
CREATE INDEX idx_payments_sale_id ON payments(sale_id);
CREATE INDEX idx_payments_payment_method ON payments(payment_method);
CREATE INDEX idx_payments_payment_status ON payments(payment_status);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);

-- Trigger
CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE payments IS 'Payment records for sales transactions';

-- ============================================================================
-- INVENTORY_TRANSACTIONS TABLE
-- Description: Track all inventory movements
-- ============================================================================

CREATE TABLE IF NOT EXISTS inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Product Reference
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    -- Transaction Type
    transaction_type inventory_transaction_type NOT NULL,

    -- Quantity Change
    quantity NUMERIC(12, 3) NOT NULL, -- Positive for increase, negative for decrease
    unit VARCHAR(50) DEFAULT 'unit',

    -- Balance After Transaction
    balance_after NUMERIC(12, 3) NOT NULL,

    -- Related Records
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    reference_number VARCHAR(100),

    -- Cost
    unit_cost NUMERIC(12, 2),
    total_cost NUMERIC(12, 2),

    -- Date
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Notes
    notes TEXT,
    reason TEXT,

    -- Custom Fields
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by UUID REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_inventory_transactions_organization_id ON inventory_transactions(organization_id);
CREATE INDEX idx_inventory_transactions_product_id ON inventory_transactions(product_id);
CREATE INDEX idx_inventory_transactions_transaction_type ON inventory_transactions(transaction_type);
CREATE INDEX idx_inventory_transactions_transaction_date ON inventory_transactions(transaction_date);
CREATE INDEX idx_inventory_transactions_sale_id ON inventory_transactions(sale_id);

-- Comments
COMMENT ON TABLE inventory_transactions IS 'Complete audit trail of all inventory movements';
COMMENT ON COLUMN inventory_transactions.quantity IS 'Positive for stock increase, negative for decrease';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V002 completed successfully!';
    RAISE NOTICE 'POS core tables created.';
    RAISE NOTICE '============================================';
END $$;

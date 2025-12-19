-- Migration: 003_create_reception_module.sql
-- Description: Crear tablas del módulo de recepción

-- Tabla de proveedores
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    brand brand NOT NULL,
    rfc VARCHAR(13),
    contact_name VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de productos (HU-04)
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    brand brand NOT NULL,
    category VARCHAR(50),
    barcode VARCHAR(13),
    weight_kg DECIMAL(10, 2),
    length_cm DECIMAL(10, 2),
    width_cm DECIMAL(10, 2),
    height_cm DECIMAL(10, 2),
    is_fragile BOOLEAN DEFAULT FALSE,
    unit_price DECIMAL(10, 2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Tabla de órdenes de recepción (HU-01)
CREATE TABLE reception_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(20) UNIQUE NOT NULL,
    supplier_id UUID REFERENCES suppliers(id),
    brand brand NOT NULL,
    invoice_number VARCHAR(50),
    invoice_file_url VARCHAR(500),
    status reception_status DEFAULT 'PENDIENTE',
    received_by UUID REFERENCES users(id),
    received_at TIMESTAMP,
    validated_by UUID REFERENCES users(id),
    validated_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de líneas de recepción (HU-02, HU-03)
CREATE TABLE reception_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_order_id UUID REFERENCES reception_orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    expected_quantity INT NOT NULL,
    counted_quantity INT,
    discrepancy INT GENERATED ALWAYS AS (counted_quantity - expected_quantity) STORED,
    lot_number VARCHAR(50),
    expiration_date DATE,
    condition product_condition DEFAULT 'APTO',
    counted_by UUID REFERENCES users(id),
    counted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de incidencias de recepción (HU-03)
CREATE TABLE reception_discrepancies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_line_id UUID REFERENCES reception_lines(id) ON DELETE CASCADE,
    expected_qty INT NOT NULL,
    counted_qty INT NOT NULL,
    difference INT NOT NULL,
    status discrepancy_status DEFAULT 'DETECTADA',
    resolution_notes TEXT,
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_suppliers_brand ON suppliers(brand);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_category ON products(brand, category);
CREATE INDEX idx_reception_orders_number ON reception_orders(order_number);
CREATE INDEX idx_reception_orders_supplier ON reception_orders(supplier_id);
CREATE INDEX idx_reception_orders_status ON reception_orders(status);
CREATE INDEX idx_reception_lines_order ON reception_lines(reception_order_id);
CREATE INDEX idx_reception_lines_product ON reception_lines(product_id);

-- Triggers
CREATE TRIGGER trigger_update_suppliers_updated_at
    BEFORE UPDATE ON suppliers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_reception_orders_updated_at
    BEFORE UPDATE ON reception_orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_reception_lines_updated_at
    BEFORE UPDATE ON reception_lines
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_reception_discrepancies_updated_at
    BEFORE UPDATE ON reception_discrepancies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

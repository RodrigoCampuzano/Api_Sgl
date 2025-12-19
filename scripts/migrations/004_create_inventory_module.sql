-- Migration: 004_create_inventory_module.sql
-- Description: Crear tablas del módulo de inventario

-- Tabla de inventario (HU-05, HU-06)
CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES products(id),
    lot_number VARCHAR(50),
    expiration_date DATE,
    quantity INT NOT NULL DEFAULT 0,
    status stock_status DEFAULT 'DISPONIBLE',
    warehouse_location VARCHAR(50),
    last_movement_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_inventory_lot UNIQUE (product_id, lot_number, warehouse_location),
    CONSTRAINT check_quantity_non_negative CHECK (quantity >= 0)
);

-- Tabla de movimientos de inventario (HU-13)
CREATE TABLE inventory_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_id UUID REFERENCES inventory(id),
    movement_type movement_type NOT NULL,
    quantity INT NOT NULL,
    previous_quantity INT NOT NULL,
    new_quantity INT NOT NULL,
    reference_id UUID,
    reference_type VARCHAR(50),
    reason TEXT,
    evidence_photo_url VARCHAR(500),
    performed_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de conteo cíclico (HU-15)
CREATE TABLE cycle_counts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scheduled_date DATE NOT NULL,
    location VARCHAR(50),
    product_id UUID REFERENCES products(id),
    expected_quantity INT,
    counted_quantity INT,
    variance INT GENERATED ALWAYS AS (counted_quantity - expected_quantity) STORED,
    counted_by UUID REFERENCES users(id),
    counted_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'PENDIENTE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_inventory_product ON inventory(product_id);
CREATE INDEX idx_inventory_status ON inventory(status);
CREATE INDEX idx_inventory_location ON inventory(warehouse_location);
CREATE INDEX idx_inventory_expiration ON inventory(expiration_date);
CREATE INDEX idx_inventory_movements_inventory ON inventory_movements(inventory_id);
CREATE INDEX idx_inventory_movements_type ON inventory_movements(movement_type);
CREATE INDEX idx_inventory_movements_created ON inventory_movements(created_at);
CREATE INDEX idx_cycle_counts_scheduled ON cycle_counts(scheduled_date);
CREATE INDEX idx_cycle_counts_status ON cycle_counts(status);

-- Triggers
CREATE TRIGGER trigger_update_inventory_updated_at
    BEFORE UPDATE ON inventory
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

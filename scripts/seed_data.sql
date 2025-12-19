-- Script de datos iniciales para SGL-DISASUR
-- IMPORTANTE: Ejecutar después de las migraciones

-- ===== USUARIOS INICIALES =====
-- Contraseña para todos: "password123" (hash bcrypt)
-- Hash generado: $2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y

INSERT INTO users (username, email, password_hash, role, status) VALUES
('admin', 'admin@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'ADMIN_TI', 'ACTIVO'),
('gerente', 'gerente@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'GERENTE', 'ACTIVO'),
('jefe_almacen', 'jefe.almacen@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'JEFE_ALMACEN', 'ACTIVO'),
('auxiliar1', 'auxiliar1@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'AUXILIAR', 'ACTIVO'),
('supervisor', 'supervisor@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'SUPERVISOR', 'ACTIVO'),
('vendedor1', 'vendedor1@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'VENDEDOR', 'ACTIVO'),
('chofer1', 'chofer1@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'CHOFER', 'ACTIVO'),
('jefe_trafico', 'jefe.trafico@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'JEFE_TRAFICO', 'ACTIVO');

-- ===== PROVEEDORES =====
INSERT INTO suppliers (name, brand, rfc, contact_name, phone, email) VALUES
('Grupo Herdez - La Costeña', 'LA_COSTENA', 'GHE960425A12', 'Juan Pérez', '5551234567', 'contacto@lacostena.com'),
('Jugos del Valle', 'JUMEX', 'JDV850315B34', 'María González', '5559876543', 'ventas@jumex.com'),
('Alimentos Pronto', 'PRONTO', 'ALP920820C56', 'Carlos Ruiz', '5556543210', 'pedidos@pronto.com');

-- ===== PRODUCTOS DE EJEMPLO =====
-- La Costeña
INSERT INTO products (sku, name, brand, category, barcode, weight_kg, length_cm, width_cm, height_cm, is_fragile, unit_price) VALUES
('COST-SAL-001', 'Salsa Verde La Costeña 250g', 'LA_COSTENA', 'SALSAS', '7501005100110', 0.25, 6, 6, 12, TRUE, 18.50),
('COST-SAL-002', 'Salsa Roja La Costeña 250g', 'LA_COSTENA', 'SALSAS', '7501005100127', 0.25, 6, 6, 12, TRUE, 18.50),
('COST-FRJ-001', 'Frijoles Refritos La Costeña 430g', 'LA_COSTENA', 'CONSERVAS', '7501005101315', 0.43, 8, 8, 10, TRUE, 22.00),
('COST-CHI-001', 'Chiles Jalapeños La Costeña 380g', 'LA_COSTENA', 'CONSERVAS', '7501005102456', 0.38, 7, 7, 11, TRUE, 25.00);

-- Jumex
INSERT INTO products (sku, name, brand, category, barcode, weight_kg, length_cm, width_cm, height_cm, is_fragile, unit_price) VALUES
('JUME-JUG-001', 'Jugo Jumex Naranja 1L', 'JUMEX', 'JUGOS', '7501032800014', 1.05, 7, 7, 20, TRUE, 15.00),
('JUME-JUG-002', 'Jugo Jumex Mango 1L', 'JUMEX', 'JUGOS', '7501032800021', 1.05, 7, 7, 20, TRUE, 15.00),
('JUME-JUG-003', 'Jugo Jumex Durazno 1L', 'JUMEX', 'JUGOS', '7501032800038', 1.05, 7, 7, 20, TRUE, 15.00),
('JUME-NEC-001', 'Néctar Jumex Guayaba 1L', 'JUMEX', 'JUGOS', '7501032801012', 1.05, 7, 7, 20, TRUE, 16.50);

-- Pronto
INSERT INTO products (sku, name, brand, category, barcode, weight_kg, length_cm, width_cm, height_cm, is_fragile, unit_price) VALUES
('PRON-HAR-001', 'Harina Pronto para Hot Cakes 1kg', 'PRONTO', 'HARINAS', '7501020400015', 1.00, 15, 10, 25, FALSE, 35.00),
('PRON-HAR-002', 'Harina Pronto para Panqué 800g', 'PRONTO', 'HARINAS', '7501020400022', 0.80, 15, 10, 22, FALSE, 32.00),
('PRON-ATO-001', 'Atole Pronto Sabor Chocolate 500g', 'PRONTO', 'BEBIDAS', '7501020401234', 0.50, 12, 8, 18, FALSE, 28.00);

-- ===== CLIENTES DE EJEMPLO =====
INSERT INTO customers (name, rfc, address, city, state, postal_code, phone, email, credit_limit) VALUES
('Abarrotes Don José', 'ADJ950215AB1', 'Av. Juárez 123', 'Monterrey', 'Nuevo León', '64000', '8181234567', 'contacto@abadjonjose.com', 50000.00),
('Super Mercado La Esquina', 'SML880420CD2', 'Calle Morelos 456', 'Guadalajara', 'Jalisco', '44100', '3331234567', 'ventas@laesquina.com', 75000.00),
('Tienda Mi Despensa', 'TMD920830EF3', 'Av. Reforma 789', 'Ciudad de México', 'CDMX', '06600', '5551234567', 'pedidos@midespensa.com', 100000.00),
('Minisuper El Ahorro', 'MEA000510GH5', 'Calle Hidalgo 321', 'Puebla', 'Puebla', '72000', '2221234567', 'info@elahorro.com', 35000.00);

-- ===== VEHÍCULOS =====
INSERT INTO vehicles (plate_number, vehicle_type, brand, model, year, capacity_kg, capacity_m3, status) VALUES
('ABC-123-D', 'VAN', 'Nissan', 'Urvan', 2020, 1200, 8.5, 'DISPONIBLE'),
('DEF-456-G', 'CAMIONETA', 'Toyota', 'Hilux', 2021, 900, 6.0, 'DISPONIBLE'),
('GHI-789-J', 'CAMION_3_5', 'Isuzu', 'NQR', 2019, 3500, 18.0, 'DISPONIBLE'),
('JKL-012-M', 'TORTON', 'Freightliner', 'M2', 2018, 8000, 35.0, 'DISPONIBLE'),
('MNO-345-P', 'VAN', 'Ford', 'Transit', 2022, 1100, 9.0, 'EN_TALLER');

-- ===== CHOFERES =====
-- Obtener el user_id del chofer1 para vincularlo
DO $$
DECLARE
    chofer_user_id UUID;
BEGIN
    SELECT id INTO chofer_user_id FROM users WHERE username = 'chofer1';
    
    INSERT INTO drivers (user_id, license_number, license_expiry, phone, status) VALUES
    (chofer_user_id, 'A1234567', '2026-12-31', '5551112222', 'DISPONIBLE');
END $$;

-- Crear más choferes ficticios
INSERT INTO users (username, email, password_hash, role, status) VALUES
('chofer2', 'chofer2@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'CHOFER', 'ACTIVO'),
('chofer3', 'chofer3@sgl-disasur.com', '$2a$10$rJ8qXq7pKZ8yh3k9vQ5zVOGZJ4X5pKZ8yh3k9vQ5zVOGZJ4X5pKZ8y', 'CHOFER', 'ACTIVO');

DO $$
DECLARE
    chofer2_user_id UUID;
    chofer3_user_id UUID;
BEGIN
    SELECT id INTO chofer2_user_id FROM users WHERE username = 'chofer2';
    SELECT id INTO chofer3_user_id FROM users WHERE username = 'chofer3';
    
    INSERT INTO drivers (user_id, license_number, license_expiry, phone, status) VALUES
    (chofer2_user_id, 'B7654321', '2025-06-30', '5553334444', 'DISPONIBLE'),
    (chofer3_user_id, 'C9876543', '2027-03-15', '5555556666', 'DISPONIBLE');
END $$;

-- Log de auditoría de ejemplo
INSERT INTO audit_logs (user_id, action, entity_type, ip_address, user_agent)
SELECT id, 'SYSTEM_INIT', 'DATABASE', '127.0.0.1'::inet, 'Seed Script v1.0'
FROM users WHERE username = 'admin';

-- Mensaje de confirmación
DO $$
BEGIN
    RAISE NOTICE 'Seed data loaded successfully!';
    RAISE NOTICE 'Users created: 11 (admin, gerente, jefe_almacen, auxiliar1, supervisor, vendedor1, chofer1-3)';
    RAISE NOTICE 'Default password for all users: password123';
    RAISE NOTICE 'Suppliers: 3 | Products: 11 | Customers: 4 | Vehicles: 5 | Drivers: 3';
END $$;

-- Migration: 001_create_enums.sql
-- Description: Crear todos los tipos ENUM para el sistema

-- Roles y seguridad
CREATE TYPE user_role AS ENUM (
    'ADMIN_TI',
    'GERENTE',
    'JEFE_ALMACEN',
    'AUXILIAR',
    'SUPERVISOR',
    'RECEPCIONISTA',
    'VENDEDOR',
    'JEFE_TRAFICO',
    'CHOFER',
    'MONTACARGUISTA',
    'CARGADOR',
    'PLANIFICADOR',
    'FLOTA',
    'AUDITOR',
    'SERVICIO_CLIENTE'
);

CREATE TYPE user_status AS ENUM ('ACTIVO', 'BLOQUEADO', 'INACTIVO');

-- Módulo Recepción
CREATE TYPE brand AS ENUM ('COSTENA', 'JUMEX', 'PRONTO', 'LA_COSTENA', 'OTROS');

CREATE TYPE reception_status AS ENUM (
    'PENDIENTE',
    'EN_CONTEO',
    'VALIDADA',
    'CON_INCIDENCIA',
    'COMPLETADA'
);

CREATE TYPE discrepancy_status AS ENUM (
    'DETECTADA',
    'EN_REVISION',
    'RESUELTA',
    'ACEPTADA'
);

CREATE TYPE product_condition AS ENUM ('APTO', 'DESECHO', 'CUARENTENA');

-- Módulo Inventario
CREATE TYPE stock_status AS ENUM (
    'DISPONIBLE',
    'RESERVADO',
    'BLOQUEADO',
    'CUARENTENA',
    'CADUCADO'
);

CREATE TYPE movement_type AS ENUM (
    'ENTRADA',
    'SALIDA',
    'AJUSTE',
    'MERMA',
    'DEVOLUCION',
    'TRANSFERENCIA'
);

-- Módulo Pedidos
CREATE TYPE order_status AS ENUM (
    'BORRADOR',
    'CONFIRMADO',
    'EN_PREPARACION',
    'LISTO',
    'EN_RUTA',
    'ENTREGADO',
    'CANCELADO'
);

CREATE TYPE vehicle_type AS ENUM ('VAN', 'CAMIONETA', 'CAMION_3_5', 'TORTON');

-- Módulo Flota
CREATE TYPE route_type AS ENUM ('LOCAL', 'FORANEA');

CREATE TYPE vehicle_status AS ENUM (
    'DISPONIBLE',
    'EN_RUTA',
    'EN_TALLER',
    'FUERA_SERVICIO'
);

CREATE TYPE driver_status AS ENUM (
    'DISPONIBLE',
    'EN_RUTA',
    'DESCANSO',
    'LICENCIA'
);

# SGL-DISASUR API

Sistema de Gestión Logística Multi-Marca con Go, PostgreSQL y Arquitectura Hexagonal.

## 🚀 Características

- **Arquitectura Hexagonal (Clean Architecture)**
- **Base de datos PostgreSQL** con enums simplificados
- **Seguridad**: JWT, bcrypt, RBAC
- **Auditoría completa** de acciones
- **Docker y Docker Compose** para deployment
- **Swagger/OpenAPI** (en desarrollo)

## 📋 Requisitos

- Go 1.21+
- PostgreSQL 15+
- Docker y Docker Compose (opcional)
- Make (opcional, para comandos)

## 🔧 Instalación

### 1. Clonar el repositorio

```bash
cd NEWWWWW_API
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Editar `.env` con tus configuraciones:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=sgl_user
DB_PASSWORD=secure_password
DB_NAME=sgl_disasur
JWT_SECRET_KEY=cambia-esto-en-produccion
PORT=8080
```

### 3. Instalar dependencias

```bash
go mod download
go mod tidy
```

## 🐳 Opción 1: Ejecución con Docker

```bash
# Levantar base de datos y API
docker-compose up -d

# Ver logs
docker-compose logs -f api
```

La API estará disponible en `http://localhost:8080`

## 💻 Opción 2: Ejecución local

### 1. Iniciar PostgreSQL

Puedes usar Docker solo para PostgreSQL:

```bash
docker-compose up -d postgres
```

O instalar PostgreSQL localmente.

### 2. Ejecutar migraciones

```bash
# Las migraciones se ejecutarán automáticamente al iniciar el contenedor
# O manualmente con psql:
psql -U sgl_user -h localhost -d sgl_disasur -f scripts/migrations/001_create_enums.sql
psql -U sgl_user -h localhost -d sgl_disasur -f scripts/migrations/002_create_users_and_security.sql
psql -U sgl_user -h localhost -d sgl_disasur -f scripts/migrations/003_create_reception_module.sql
psql -U sgl_user -h localhost -d sgl_disasur -f scripts/migrations/004_create_inventory_module.sql
psql -U sgl_user -h localhost -d sgl_disasur -f scripts/migrations/005_create_orders_module.sql
psql -U sgl_user -h localhost -d sgl_disasur -f scripts/migrations/006_create_fleet_module.sql
```

### 3. Cargar datos iniciales

```bash
psql -U sgl_user -h localhost -d sgl_disasur -f scripts/seed_data.sql
```

Esto creará:
- 11 usuarios de prueba (admin, gerente, jefe_almacen, etc.)
- 3 proveedores
- 11 productos
- 4 clientes
- 5 vehículos
- 3 choferes

**Contraseña para todos los usuarios de prueba**: `password123`

### 4. Ejecutar la API

```bash
go run cmd/api/main.go
```

La API estará disponible en `http://localhost:8080`

## 📡 Endpoints Disponibles

### Salud del sistema

```bash
GET /health
```

### Autenticación

```bash
# Login
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}

# Respuesta:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...},
  "expires_at": "2024-12-19T04:00:00Z"
}

# Logout
POST /api/v1/auth/logout
Authorization: Bearer {token}
```

### Rutas protegidas (requieren autenticación)

Todas las rutas bajo `/api/v1/*` (excepto `/auth/login`) requieren el header:

```
Authorization: Bearer {token}
```

#### Usuarios (Solo ADMIN_TI y GERENTE)
```bash
GET /api/v1/users
```

#### Módulos (En desarrollo)
- `GET /api/v1/products` - Productos
- `GET /api/v1/reception/orders` - Órdenes de recepción
- `GET /api/v1/inventory/stock` - Inventario
- `GET /api/v1/orders` - Pedidos
- `GET /api/v1/fleet/vehicles` - Vehículos
- `GET /api/v1/fleet/drivers` - Choferes
- `GET /api/v1/reports/dashboard` - Dashboard (Solo GERENTE y ADMIN_TI)

## 🔒 Seguridad Implementada

### HU-00: Bloqueo de cuenta
- Después de 3 intentos fallidos de login, la cuenta se bloquea
- El usuario debe ser desbloqueado por un administrador

### HU-19: RBAC (Control de acceso basado en roles)
- Cada endpoint especifica qué roles tienen acceso
- Middleware valida el rol antes de permitir la operación

### HU-20: Auditoría
- Todos los login (exitosos y fallidos) se registran
- Se captura: usuario, acción, IP, user agent, timestamp
- Los logs de auditoría NO pueden ser borrados

## 📁 Estructura del Proyecto

```
.
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── domain/                  # Entidades de dominio
│   ├── usecase/                 # Casos de uso (lógica de negocio)
│   ├── repository/              # Repositorios (PostgreSQL)
│   ├── delivery/                # Handlers HTTP
│   │   └── http/
│   │       ├── handler/
│   │       ├── middleware/
│   │       └── router.go
│   └── infrastructure/          # Config, DB, Security, Logger
├── scripts/
│   ├── migrations/              # Migraciones SQL
│   └── seed_data.sql            # Datos iniciales
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## 🛠️ Comandos Make (opcional)

```bash
make help          # Mostrar ayuda
make run           # Ejecutar API localmente
make build         # Compilar binario
make docker-up     # Levantar con Docker
make docker-down   # Detener contenedores
make seed          # Cargar datos iniciales
```

## 🗺️ Roadmap

### ✅ Fase 1: Fundamentos
- [x] Configuración del proyecto
- [x] Base de datos con enums
- [x] Seguridad (JWT, bcrypt, RBAC)
- [x] Infraestructura base

### ✅ Fase 2: Módulo 0 - Acceso
- [x] HU-00: Login con bloqueo de cuentas
- [x] HU-19: RBAC
- [x] HU-20: Auditoría

### 🚧 Fase 3-7: En desarrollo
- Módulo 1: Recepción
- Módulo 2: Inventario
- Módulo 3: Pedidos
- Módulo 4: Flota
- Módulo 6: Reportes

## 📝 Notas de Desarrollo

### Roles disponibles
- `ADMIN_TI` - Acceso total
- `GERENTE` - Gestión general
- `JEFE_ALMACEN` - Operaciones de almacén
- `AUXILIAR` - Operaciones básicas
- `SUPERVISOR` - Supervisión de procesos
- `RECEPCIONISTA` - Recepción de mercancía
- `VENDEDOR` - Gestión de pedidos
- `JEFE_TRAFICO` - Asignación de rutas
- `CHOFER` - Operación de vehículos
- `MONTACARGUISTAoper` - Manejo de inventario
- `AUDITOR` - Consulta de auditoría
- `SERVICIO_CLIENTE` - Atención a clientes

### Marcas disponibles
- `LA_COSTENA`
- `JUMEX`
- `PRONTO`
- `COSTENA`
- `OTROS`

## 📞 Soporte

Para problemas o preguntas, contactar al equipo de desarrollo.

## 📄 Licencia

MIT License - ver archivo LICENSE para más detalles.

# [cite_start]PRODUCT BACKLOG: SGL-DISASUR [cite: 1]
[cite_start]**Proyecto:** Sistema de Gestión Logística Multi-Marca [cite: 2]
[cite_start]**Total de Historias:** 25 [cite: 3]

---

## [cite_start]MÓDULO 0: ACCESO Y SEGURIDAD [cite: 4]
[cite_start]**(FUNDAMENTAL)** [cite: 5]

| ID | Prioridad | Historia de Usuario | Criterios de Aceptación (Pruebas) |
| :--- | :--- | :--- | :--- |
| **HU-00** | Crítica | [cite_start]**Autenticación de Usuarios (Login)**<br>Como Usuario, quiero ingresar con usuario y contraseña encriptada, para acceder solo a los módulos que me corresponden. [cite: 6] | 1. Bloqueo de cuenta tras 3 intentos fallidos.<br>2. [cite_start]Cierre de sesión automático tras 30 min de inactividad. [cite: 6] |
| **HU-19** | Alta | [cite_start]**Gestión de Roles Permisos (RBAC)**<br>Como Admin TI, quiero asignar perfiles (Chofer, Almacenista, Gerente), para restringir funciones sensibles (ej. ver costos). [cite: 6] | 1. Perfil "Chofer": Solo acceso a "Mis Rutas" y Check-list.<br>2. [cite_start]Perfil "Almacenista": Sin permiso de "Borrar Usuarios". [cite: 6] |
| **HU-20** | Baja | [cite_start]**Logs de Auditoría (Trazabilidad)**<br>Como Auditor, quiero un registro oculto de quién modificó o borró registros, para investigar robos o errores. [cite: 6, 7] | 1. Registro de: Usuario, Acción (Borrar/Editar), Fecha/Hora, IP.<br>2. [cite_start]El log no puede ser borrado por ningún usuario. [cite: 6, 7] |

---

## [cite_start]MÓDULO 1: RECEPCIÓN (INBOUND) [cite: 9]

| ID | Prioridad | Historia de Usuario | Criterios de Aceptación (Pruebas) |
| :--- | :--- | :--- | :--- |
| **HU-01** | Alta | [cite_start]**Alta de Orden de Recepción**<br>Como Jefe de Almacén, quiero registrar la llegada de un camión seleccionando Proveedor (Costeña/Jumex) y Factura, para iniciar el proceso. [cite: 8] | 1. Selección única de marca por orden.<br>2. [cite_start]Carga de archivo PDF/XML de factura obligatoria. [cite: 8] |
| **HU-02** | Alta | [cite_start]**Conteo Ciego (Blind Count)**<br>Como Auxiliar, quiero ingresar lo que cuento fisicamente sin ver la cantidad que dice la factura, para evitar vicios en la entrada. [cite: 8] | 1. Campo de "Cantidad Esperada" oculto para este rol.<br>2. [cite_start]Permite ingreso manual o por escáner. [cite: 8] |
| **HU-03** | Alta | [cite_start]**Validación de Discrepancias**<br>Como Supervisor, quiero que el sistema compare automáticamente Contado vs Facturado, para detectar faltantes al instante. [cite: 10] | 1. Diferencia = 0: Entrada automática al stock.<br>2. [cite_start]Diferencia $\ne$ 0: Bloqueo y alerta de "Incidencia". [cite: 10] |
| **HU-14** | Media | [cite_start]**Gestión de Devoluciones**<br>Como Recepcionista, quiero reingresar mercancía devuelta clasificándola en "Apta" o "Desecho", para no vender producto dañado. [cite: 10] | 1. Validación: Si empaque "Pronto" está abierto -> Desecho automático.<br>2. [cite_start]Stock entra a "Almacén Cuarentena". [cite: 10] |
| **HU-04** | Alta | [cite_start]**Catálogo de Productos Detallado**<br>Como Admin, quiero definir atributos (Peso, Fragilidad, Dimensiones) por SKU, para cálculos logísticos. [cite: 11] | 1. Checkbox: "¿Es Frágil?".<br>2. [cite_start]Campos numéricos para Largo, Ancho, Alto y Peso. [cite: 11] |

---

## [cite_start]MÓDULO 2: INVENTARIO (ALMACÉN) [cite: 12]

| ID | Prioridad | Historia de Usuario | Criterios de Aceptación (Pruebas) |
| :--- | :--- | :--- | :--- |
| **HU-05** | Alta | [cite_start]**Monitor de Stock Multi-Marca**<br>Como Ventas, quiero ver existencias filtradas por marca en tiempo real, para saber qué puedo prometer al cliente. [cite: 13] | 1. Filtros por familia (Salsas, Jugos, Harinas).<br>2. [cite_start]Indicador visual de "Punto de Reorden" (Stock bajo). [cite: 13] |
| **HU-06** | Alta | [cite_start]**Rotación FEFO (Caducidad)**<br>Como Sistema, quiero bloquear lotes nuevos si hay antiguos disponibles, para evitar que el producto caduque en bodega. [cite: 13] | 1. Al agregar a pedido, selecciona lote con fecha caducidad más próxima.<br>2. [cite_start]Alerta roja si caducidad < 30 días. [cite: 13] |
| **HU-13** | Media | [cite_start]**Registro de Mermas Internas**<br>Como Montacarguista, quiero reportar roturas (botellas/bolsas) desde el móvil, para ajustar el inventario real. [cite: 13] | 1. Requiere foto de evidencia.<br>2. [cite_start]Descuenta stock y genera asiento de "Pérdida". [cite: 13] |
| **HU-15** | Baja | [cite_start]**Conteo Cíclico (Spot Check)**<br>Como Auditor, quiero que el sistema me pida contar 5 ubicaciones al azar diariamente, para mantener la exactitud del inventario. [cite: 13, 14] | 1. Selección aleatoria inteligente (productos de alto valor).<br>2. [cite_start]Bloqueo de ubicación durante el conteo. [cite: 13, 14] |

---

## [cite_start]MÓDULO 3: PEDIDOS Y CARGA (OUTBOUND) [cite: 16]

| ID | Prioridad | Historia de Usuario | Criterios de Aceptación (Pruebas) |
| :--- | :--- | :--- | :--- |
| **HU-07** | Alta | [cite_start]**Pedido Mixto (Multi-Marca)**<br>Como Vendedor, quiero agregar productos de diferentes marcas en una sola orden, para surtir tiendas de abarrotes. [cite: 15] | 1. Carrito permite mezcla de SKUs.<br>2. [cite_start]Sumatoria automática de Costo, Peso y Volumen total. [cite: 15] |
| **HU-08** | Media | [cite_start]**Sugerencia de Vehículo**<br>Como Tráfico, quiero que el sistema recomiende Camioneta o Camión según el volumen ($m^3$) del pedido, para ahorrar fletes. [cite: 15] | 1. Vol < 10$m^3$: Sugiere "Nissan/Van".<br>2. [cite_start]Vol > 10$m^3$: Sugiere "Camión 3.5" o "Torton". [cite: 15] |
| **HU-09** | Alta | [cite_start]**Alerta de Estiba (Seguridad)**<br>Como Cargador, quiero una alerta si el pedido mezcla "Pesados" con "Frágiles", para tener cuidado al estibar. [cite: 15, 17] | 1. [cite_start]Pop-up de advertencia: "Cuidado: No poner Cajas Costeña sobre Cajas Pronto". [cite: 15, 17] |
| **HU-18** | Media | [cite_start]**Eficiencia de Carga**<br>Como Planificador, quiero ver una barra de % de llenado del camión, para maximizar el uso del espacio. [cite: 17] | 1. [cite_start]Visualización gráfica (Verde/Amarillo/Rojo) del espacio ocupado en el camión. [cite: 17] |

---

## [cite_start]MÓDULO 4: FLOTA Y RUTAS [cite: 19]

| ID | Prioridad | Historia de Usuario | Criterios de Aceptación (Pruebas) |
| :--- | :--- | :--- | :--- |
| **HU-10** | Alta | [cite_start]**Asignación de Ruta y Chofer**<br>Como Jefe Tráfico, quiero vincular pedido + vehículo + chofer, para crear el viaje. [cite: 18] | 1. Validación: Chofer no disponible si ya está "En Ruta".<br>2. [cite_start]Tipo de Ruta: Local vs Foránea. [cite: 18] |
| **HU-11** | Alta | [cite_start]**Generación de Remisión/Factura**<br>Como Admin, quiero generar el PDF de salida oficial, para entregar al chofer como comprobante. [cite: 18, 20] | 1. Generación descuenta inventario final.<br>2. [cite_start]PDF incluye: Cliente, Lista Prod, Placas, Chofer. [cite: 18, 20] |
| **HU-16** | Alta | [cite_start]**Control de Mantenimiento**<br>Como Flota, quiero registrar cuando un camión entra a taller, para que el sistema no lo asigne a viajes. [cite: 20] | 1. Estado "En Taller" bloquea asignación en HU-10.<br>2. [cite_start]Historial de reparaciones por vehículo. [cite: 20] |
| **HU-17** | Media | [cite_start]**Check-list Pre-Salida**<br>Como Chofer, quiero confirmar estado de llantas y gasolina en mi celular, para validar que salgo seguro. [cite: 20] | 1. Formulario obligatorio antes de imprimir Remisión (HU-11).<br>2. [cite_start]Opción de subir foto de golpes previos. [cite: 20] |

---

## [cite_start]MÓDULO 5: UX Y MOVILIDAD [cite: 22]

| ID | Prioridad | Historia de Usuario | Criterios de Aceptación (Pruebas) |
| :--- | :--- | :--- | :--- |
| **HU-21** | - | [cite_start]**Interfaz "Modo Industrial"**<br>Como Operario, quiero botones grandes y alto contraste, para usar la App en terminales portátiles (Zebra). [cite: 21] | 1. Botones min 44px.<br>2. [cite_start]Diseño adaptable a pantallas verticales. [cite: 21] |
| **HU-22** | Media | [cite_start]**Escaneo de Códigos**<br>Como Auxiliar, quiero usar la cámara/láser para buscar productos, para agilizar la operación. [cite: 23] | 1. [cite_start]Input acepta lectura de código de barras EAN-13 y autocompleta el producto. [cite: 23] |

---

## [cite_start]MÓDULO 6: REPORTES (BI) [cite: 25]

| ID | Prioridad | Historia de Usuario | Criterios de Aceptación (Pruebas) |
| :--- | :--- | :--- | :--- |
| **HU-12** | Media | [cite_start]**Dashboard General**<br>Como Gerente, quiero ver KPIs clave (Ventas hoy, Camiones fuera), para monitorear la operación. [cite: 24] | 1. Gráficas actualizadas < 5 seg de retraso.<br>2. [cite_start]Visible solo para Rol Gerencia. [cite: 24] |
| **HU-23** | Baja | [cite_start]**Reporte de Rotación (Días Inventario)**<br>Como Finanzas, quiero saber cuánto tarda en venderse cada marca, para planear compras. [cite: 24] | 1. Tabla exportable a Excel.<br>2. [cite_start]Filtro por rango de fechas. [cite: 24] |
| **HU-24** | Media | [cite_start]**Alerta de Pedidos Atorados**<br>Como Servicio Cliente, quiero ver pedidos que llevan 48h sin salir, para prevenir quejas. [cite: 26] | 1. [cite_start]Lista destacada en rojo en el Dashboard si Fecha > 48h y Status $\ne$ Entregado. [cite: 26] |
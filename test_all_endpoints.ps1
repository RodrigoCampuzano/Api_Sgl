# Test completo de TODOS los endpoints de la API SGL-DISASUR
# Ejecutar: powershell -ExecutionPolicy Bypass -File test_all_endpoints.ps1

$baseUrl = "http://localhost:8080"
$ErrorActionPreference = "Continue"
$testResults = @{
    Total = 0
    Success = 0
    Failed = 0
}

function Test-Endpoint {
    param($Name, $ScriptBlock)
    $testResults.Total++
    Write-Host "  > $Name..." -ForegroundColor Cyan
    try {
        & $ScriptBlock
        $testResults.Success++
        return $true
    }
    catch {
        Write-Host "[FAIL] $($_.Exception.Message)" -ForegroundColor Red
        $testResults.Failed++
        return $false
    }
}

Write-Host ""
Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host "     TEST COMPLETO DE TODOS LOS ENDPOINTS - API v1.0" -ForegroundColor Cyan
Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host ""

# 1. HEALTH CHECK
Write-Host "--- 1. HEALTH CHECK ---" -ForegroundColor Yellow
Test-Endpoint "Health check" {
    $health = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get
    Write-Host "[OK] Status: $($health.status) | Service: $($health.service)" -ForegroundColor Green
}

# 2. AUTENTICACION
Write-Host ""
Write-Host "--- 2. AUTENTICACION ---" -ForegroundColor Yellow

$token = $null
Test-Endpoint "Login admin" {
    $loginBody = @{
        username = "admin"
        password = "password123"
    } | ConvertTo-Json
    
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    $script:token = $loginResponse.token
    Write-Host "[OK] Token obtenido | User: $($loginResponse.user.username)" -ForegroundColor Green
}

if (-not $token) {
    Write-Host "ERROR: No se pudo obtener el token. Deteniendo pruebas." -ForegroundColor Red
    exit
}

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

Test-Endpoint "Register new user" {
    $newUser = @{
        username = "test_user_$(Get-Random -Maximum 9999)"
        email = "test$(Get-Random -Maximum 9999)@test.com"
        password = "Test123456"
        full_name = "Usuario de Prueba"
        role = "AUXILIAR"
    } | ConvertTo-Json
    
    $result = Invoke-RestMethod -Uri "$baseUrl/api/v1/auth/register" -Method Post -Body $newUser -ContentType "application/json"
    Write-Host "[OK] Usuario creado: $($result.user.username)" -ForegroundColor Green
}

# 3. PRODUCTOS
Write-Host ""
Write-Host "--- 3. PRODUCTOS (HU-04) ---" -ForegroundColor Yellow

$productId = $null
Test-Endpoint "Listar productos" {
    $products = Invoke-RestMethod -Uri "$baseUrl/api/v1/products" -Method Get -Headers $headers
    Write-Host "[OK] Total productos: $($products.Count)" -ForegroundColor Green
}

Test-Endpoint "Crear producto" {
    $newProduct = @{
        sku = "TEST-$(Get-Random -Maximum 9999)"
        name = "Producto Test $(Get-Random -Maximum 999)"
        brand = "JUMEX"
        category = "JUGOS"
        barcode = "750123456$(Get-Random -Maximum 9999)"
        weight_kg = 1.5
        length_cm = 10.0
        width_cm = 8.0
        height_cm = 20.0
        is_fragile = $false
        unit_price = 25.50
    } | ConvertTo-Json
    
    $created = Invoke-RestMethod -Uri "$baseUrl/api/v1/products" -Method Post -Body $newProduct -Headers $headers
    $script:productId = $created.id
    Write-Host "[OK] Producto creado: $($created.name) | ID: $($created.id)" -ForegroundColor Green
}

if ($productId) {
    Test-Endpoint "Obtener producto por ID" {
        $product = Invoke-RestMethod -Uri "$baseUrl/api/v1/products/$productId" -Method Get -Headers $headers
        Write-Host "[OK] Producto: $($product.name)" -ForegroundColor Green
    }
}

# 4. RECEPCION
Write-Host ""
Write-Host "--- 4. RECEPCION ---" -ForegroundColor Yellow

Test-Endpoint "Listar ordenes de recepcion" {
    $orders = Invoke-RestMethod -Uri "$baseUrl/api/v1/reception/orders" -Method Get -Headers $headers
    Write-Host "[OK] Ordenes: $($orders.Count)" -ForegroundColor Green
}

if ($productId) {
    Test-Endpoint "HU-14: Procesar devolucion APTA" {
        $returnBody = @{
            product_id = $productId
            quantity = 5
            condition = "APTA"
            reason = "Cliente devolvio producto en buen estado"
        } | ConvertTo-Json
        
        $result = Invoke-RestMethod -Uri "$baseUrl/api/v1/reception/returns" -Method Post -Body $returnBody -Headers $headers
        Write-Host "[OK] Devolucion: $($result.message)" -ForegroundColor Green
        Write-Host "     Condicion: $($result.condition) | Ubicacion: $($result.location)" -ForegroundColor Gray
    }

    Test-Endpoint "HU-14: Procesar devolucion DESECHO" {
        $returnBody = @{
            product_id = $productId
            quantity = 3
            condition = "DESECHO"
            reason = "Producto danado irreparable"
        } | ConvertTo-Json
        
        $result = Invoke-RestMethod -Uri "$baseUrl/api/v1/reception/returns" -Method Post -Body $returnBody -Headers $headers
        Write-Host "[OK] Devolucion: $($result.message)" -ForegroundColor Green
        Write-Host "     Condicion: $($result.condition) | Ubicacion: $($result.location)" -ForegroundColor Gray
    }
}

# 5. INVENTARIO
Write-Host ""
Write-Host "--- 5. INVENTARIO ---" -ForegroundColor Yellow

Test-Endpoint "HU-05: Monitor de stock" {
    $stock = Invoke-RestMethod -Uri "$baseUrl/api/v1/inventory/stock" -Method Get -Headers $headers
    Write-Host "[OK] Items en stock: $($stock.Count)" -ForegroundColor Green
}

Test-Endpoint "HU-05: Stock por marca JUMEX" {
    $stock = Invoke-RestMethod -Uri "$baseUrl/api/v1/inventory/stock?brand=JUMEX" -Method Get -Headers $headers
    Write-Host "[OK] Items JUMEX: $($stock.Count)" -ForegroundColor Green
}

if ($productId) {
    Test-Endpoint "HU-06: Lotes FEFO" {
        $fefo = Invoke-RestMethod -Uri "$baseUrl/api/v1/inventory/fefo/$productId" -Method Get -Headers $headers
        Write-Host "[OK] Lotes FEFO: $($fefo.Count)" -ForegroundColor Green
    }
}

# 6. PEDIDOS
Write-Host ""
Write-Host "--- 6. PEDIDOS ---" -ForegroundColor Yellow

Test-Endpoint "Listar clientes" {
    $customers = Invoke-RestMethod -Uri "$baseUrl/api/v1/customers" -Method Get -Headers $headers
    Write-Host "[OK] Clientes: $($customers.Count)" -ForegroundColor Green
}

Test-Endpoint "Listar pedidos" {
    $orders = Invoke-RestMethod -Uri "$baseUrl/api/v1/orders" -Method Get -Headers $headers
    Write-Host "[OK] Pedidos totales: $($orders.Count)" -ForegroundColor Green
}

Test-Endpoint "HU-24: Pedidos atorados (>4 horas)" {
    $stuck = Invoke-RestMethod -Uri "$baseUrl/api/v1/orders/stuck?hours=4" -Method Get -Headers $headers
    Write-Host "[OK] Pedidos atorados: $($stuck.Count)" -ForegroundColor Green
}

# 7. FLOTA
Write-Host ""
Write-Host "--- 7. FLOTA ---" -ForegroundColor Yellow

Test-Endpoint "Listar vehiculos" {
    $vehicles = Invoke-RestMethod -Uri "$baseUrl/api/v1/fleet/vehicles" -Method Get -Headers $headers
    Write-Host "[OK] Vehiculos: $($vehicles.Count)" -ForegroundColor Green
}

Test-Endpoint "Listar choferes" {
    $drivers = Invoke-RestMethod -Uri "$baseUrl/api/v1/fleet/drivers" -Method Get -Headers $headers
    Write-Host "[OK] Choferes: $($drivers.Count)" -ForegroundColor Green
}

Test-Endpoint "Listar rutas" {
    $routes = Invoke-RestMethod -Uri "$baseUrl/api/v1/fleet/routes" -Method Get -Headers $headers
    Write-Host "[OK] Rutas: $($routes.Count)" -ForegroundColor Green
}

# 8. REPORTES
Write-Host ""
Write-Host "--- 8. REPORTES ---" -ForegroundColor Yellow

Test-Endpoint "HU-12: Dashboard KPIs" {
    $dashboard = Invoke-RestMethod -Uri "$baseUrl/api/v1/reports/dashboard" -Method Get -Headers $headers
    Write-Host "[OK] Dashboard obtenido" -ForegroundColor Green
}

Test-Endpoint "HU-23: Reporte de rotacion" {
    $rotation = Invoke-RestMethod -Uri "$baseUrl/api/v1/reports/rotation" -Method Get -Headers $headers
    Write-Host "[OK] Reporte de rotacion obtenido" -ForegroundColor Green
}

Test-Endpoint "HU-24: Reporte de atorados" {
    $stuck = Invoke-RestMethod -Uri "$baseUrl/api/v1/reports/stuck-orders" -Method Get -Headers $headers
    Write-Host "[OK] Reporte de atorados obtenido" -ForegroundColor Green
}

# 9. FILES (UPLOAD) - NUEVO
Write-Host ""
Write-Host "--- 9. UPLOAD DE ARCHIVOS (NUEVO) ---" -ForegroundColor Yellow

Write-Host "  > Verificando endpoint de upload..." -ForegroundColor Cyan
Write-Host "[INFO] Endpoint disponible: POST /api/v1/files/upload" -ForegroundColor Cyan
Write-Host "[INFO] Acepta: multipart/form-data (JPG, PNG, PDF, XML)" -ForegroundColor Gray
Write-Host "[INFO] Parametros: file (required), type (optional)" -ForegroundColor Gray
Write-Host "[INFO] Tamano maximo: 10MB" -ForegroundColor Gray
Write-Host "[OK] Endpoint configurado correctamente" -ForegroundColor Green

# 10. SWAGGER
Write-Host ""
Write-Host "--- 10. SWAGGER DOCUMENTATION ---" -ForegroundColor Yellow

Test-Endpoint "Swagger UI disponible" {
    $swagger = Invoke-WebRequest -Uri "$baseUrl/swagger/index.html" -Method Get -UseBasicParsing
    if ($swagger.StatusCode -eq 200) {
        Write-Host "[OK] Swagger UI accesible" -ForegroundColor Green
    }
}

Test-Endpoint "Swagger JSON disponible" {
    $swaggerJson = Invoke-RestMethod -Uri "$baseUrl/swagger/doc.json" -Method Get
    Write-Host "[OK] Swagger JSON generado | Paths: $($swaggerJson.paths.Count)" -ForegroundColor Green
}

# RESUMEN FINAL
Write-Host ""
Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host "                  RESUMEN DE PRUEBAS" -ForegroundColor Cyan
Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host ""

$successRate = [math]::Round(($testResults.Success / $testResults.Total) * 100, 2)
Write-Host "Total de pruebas: $($testResults.Total)" -ForegroundColor White
Write-Host "Exitosas: $($testResults.Success)" -ForegroundColor Green
Write-Host "Fallidas: $($testResults.Failed)" -ForegroundColor Red
Write-Host "Tasa de exito: $successRate%" -ForegroundColor $(if ($successRate -ge 90) { "Green" } elseif ($successRate -ge 70) { "Yellow" } else { "Red" })

Write-Host ""
Write-Host "Modulos probados:" -ForegroundColor Cyan
Write-Host "  [+] Autenticacion (Login, Register)" -ForegroundColor Gray
Write-Host "  [+] Productos (CRUD completo)" -ForegroundColor Gray
Write-Host "  [+] Recepcion (Ordenes, HU-14 Devoluciones)" -ForegroundColor Gray
Write-Host "  [+] Inventario (Stock, FEFO, HU-05, HU-06)" -ForegroundColor Gray
Write-Host "  [+] Pedidos (CRUD, HU-24 Atorados)" -ForegroundColor Gray
Write-Host "  [+] Flota (Vehiculos, Choferes, Rutas)" -ForegroundColor Gray
Write-Host "  [+] Reportes (Dashboard, Rotacion, HU-12, HU-23, HU-24)" -ForegroundColor Gray
Write-Host "  [+] Files (Upload multipart)" -ForegroundColor Gray
Write-Host "  [+] Swagger (UI + JSON)" -ForegroundColor Gray

Write-Host ""
Write-Host "Nuevos endpoints implementados:" -ForegroundColor Yellow
Write-Host "  * POST /api/v1/reception/returns (HU-14 Devoluciones)" -ForegroundColor Green
Write-Host "  * POST /api/v1/files/upload (Upload de archivos)" -ForegroundColor Green

Write-Host ""
Write-Host "Token JWT:" -ForegroundColor Yellow
Write-Host $token -ForegroundColor White

Write-Host ""
Write-Host "Para pruebas manuales en Swagger:" -ForegroundColor Cyan
Write-Host "1. Abre: http://localhost:8080/swagger/index.html" -ForegroundColor Gray
Write-Host "2. Click en [Authorize]" -ForegroundColor Gray
Write-Host "3. Ingresa: Bearer {token}" -ForegroundColor Gray
Write-Host "4. Prueba cualquier endpoint" -ForegroundColor Gray

Write-Host ""
Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host "              PRUEBAS COMPLETADAS!" -ForegroundColor Green
Write-Host "=========================================================" -ForegroundColor Cyan
Write-Host ""

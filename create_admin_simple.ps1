# Crear usuario administrador
# Uso: powershell -ExecutionPolicy Bypass -File create_admin_simple.ps1

$apiUrl = "http://localhost:8080/api/v1/auth/register"

$jsonBody = '{
  "username": "admin",
  "email": "admin@sgl-disasur.com",
  "password": "password123",
  "role": "ADMIN_TI"
}'

Write-Host "Creando usuario administrador..." -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri $apiUrl -Method Post -Body $jsonBody -ContentType "application/json"
    Write-Host "`nUsuario creado exitosamente!" -ForegroundColor Green
    Write-Host "Username: admin" -ForegroundColor Yellow
    Write-Host "Password: password123" -ForegroundColor Yellow
    Write-Host "`nID: $($response.user.id)" -ForegroundColor Gray
}
catch {
    Write-Host "`nError al crear usuario" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}

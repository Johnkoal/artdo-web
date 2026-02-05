# Helper script to run the Artdotech website
# In Go, you must compile all files in the package together.
# This script ensures that both main.go and db.go are included.

Write-Host "Iniciando servidor Artdotech..." -ForegroundColor Cyan
go run .

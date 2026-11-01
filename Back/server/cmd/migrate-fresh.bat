@echo off
setlocal
REM Migrate Fresh para PostgreSQL
REM Uso: migrate-fresh.bat

REM Ubicarse en Back/server sin importar desde donde se ejecute
cd /d "%~dp0.."

echo.
echo ========================================
echo  EJECUTANDO MIGRATE:FRESH
echo ========================================
echo.

set "DB_PASSWORD="
set /p DB_PASSWORD=Ingresa el password de PostgreSQL: 

set "DB_CONN=host=localhost port=5432 user=postgres password=%DB_PASSWORD% dbname=FileManagerDB sslmode=disable"

where goose >nul 2>&1
if errorlevel 1 (
    echo ERROR: goose no encontrado en PATH
    pause
    exit /b 1
)

echo.
echo 1. Revirtiendo todas las migraciones...
goose postgres "%DB_CONN%" -dir=db/migrations reset
if errorlevel 1 (
    echo ERROR: Fallo al revertir migraciones
    pause
    exit /b 1
)

echo.
echo 2. Ejecutando migraciones...
goose postgres "%DB_CONN%" -dir=db/migrations up
if errorlevel 1 (
    echo ERROR: Fallo al ejecutar migraciones
    pause
    exit /b 1
)

echo.
echo ========================================
echo  MIGRATE:FRESH COMPLETADO
echo ========================================
echo.
pause
endlocal

@echo off
REM Migrate Fresh para PostgreSQL
REM Uso: migrate-fresh.bat

echo.
echo ========================================
echo  EJECUTANDO MIGRATE:FRESH
echo ========================================
echo.

REM Obtener password si no esta definido
set /p DB_PASSWORD:"Ingresa el password de PostgreSQL (presiona Enter si esta vacio): "

REM Construir connection string
set DB_CONN=postgresql://postgres:%DB_PASSWORD%@localhost:5432/FileManagerDB

echo.
echo 1. Eliminando todas las tablas...
psql %DB_CONN% -c "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Fallo al recrear esquema
    pause
    exit /b 1
)

echo.
echo 2. Ejecutando migraciones...
goose -dir=db/migrations %DB_CONN% up
if %ERRORLEVEL% NEQ 0 (
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
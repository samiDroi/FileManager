#!/bin/bash
set -e
# Migrate Fresh para PostgreSQL
# Uso: ./migrate-fresh.sh

# Ubicarse en Back/server sin importar desde donde se ejecute
cd "$(dirname "$0")/.." || exit 1

echo ""
echo "========================================"
echo "  EJECUTANDO MIGRATE:FRESH"
echo "========================================"
echo ""

read -s -p "Ingresa el password de PostgreSQL: " DB_PASSWORD
echo ""

DB_CONN="host=localhost port=5432 user=postgres password=${DB_PASSWORD} dbname=FileManagerDB sslmode=disable"

if ! command -v goose >/dev/null 2>&1; then
    echo "ERROR: goose no encontrado en PATH"
    exit 1
fi

echo ""
echo "1. Revirtiendo todas las migraciones..."
if ! goose postgres "$DB_CONN" -dir=db/migrations reset; then
    echo "ERROR: Fallo al revertir migraciones"
    exit 1
fi

echo ""
echo "2. Ejecutando migraciones..."
if ! goose postgres "$DB_CONN" -dir=db/migrations up; then
    echo "ERROR: Fallo al ejecutar migraciones"
    exit 1
fi

echo ""
echo "========================================"
echo "  MIGRATE:FRESH COMPLETADO"
echo "========================================"
echo ""

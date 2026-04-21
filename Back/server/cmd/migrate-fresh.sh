#!/bin/bash
# Migrate Fresh para PostgreSQL
# Uso: ./migrate-fresh.sh

echo ""
echo "========================================"
echo "  EJECUTANDO MIGRATE:FRESH"
echo "========================================"
echo ""

read -s -p "Ingresa el password de PostgreSQL: " DB_PASSWORD
echo ""

DB_CONN="postgresql://postgres:${DB_PASSWORD}@localhost:5432/FileManagerDB"

echo ""
echo "1. Eliminando todas las tablas..."
psql "$DB_CONN" -c "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"
if [ $? -ne 0 ]; then
    echo "ERROR: Fallo al recrear esquema"
    exit 1
fi

echo ""
echo "2. Ejecutando migraciones..."
goose -dir=db/migrations "$DB_CONN" up
if [ $? -ne 0 ]; then
    echo "ERROR: Fallo al ejecutar migraciones"
    exit 1
fi

echo ""
echo "========================================"
echo "  MIGRATE:FRESH COMPLETADO"
echo "========================================"
echo ""
package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"server/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

var Pool *pgxpool.Pool

// InitDB inicializa el pool de pgx y el cliente sqlc
func InitDB(connStr string) *pgxpool.Pool {
	// Configuración básica del pool
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("InitDB: ParseConfig failed: %v", err)
	}

	// Conexiones máximas simultáneas
	cfg.MaxConns = 25
	// Mantener al menos X conexiones listas
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 15 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	// Conexión al pool
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("InitDB: pgxpool.NewWithConfig failed: %v", err)
	}

	Pool = pool
	return Pool
}

// CloseDB closes all connections in the pool
func CloseDB() {
	if Pool != nil {
		Pool.Close()
		log.Println("Database pool closed.")
	}
}

// TestDBConnection tests the database pool connection
func TestDBConnection() error {
	if Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection from pool: %v", err)
	}
	defer conn.Release()

	if err := conn.Conn().Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Database connection test successful.
	return nil
}
func RunMigrations() {
	// 1. Obtenemos la URL del .env que ya cargamos en el main
	dbConn := config.Config.DBConn

	// 2. Goose necesita una conexión estándar de Go (sql.DB)
	// El driver "pgx" está registrado automáticamente por la importación
	sqldb, err := sql.Open("pgx", dbConn)
	if err != nil {
		log.Fatalf("Error al abrir conexión para goose: %v", err)
	}
	defer sqldb.Close()

	// 3. Probar la conexión
	if err := sqldb.Ping(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// 4. Ejecutar las migraciones que estén en la carpeta db/migrations
	log.Println("Checking for database migrations...")
	if err := goose.Up(sqldb, "db/migrations"); err != nil {
		log.Fatalf("Error en migraciones: %v", err)
	}

	log.Println("Migraciones ejecutadas exitosamente")
}

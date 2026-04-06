package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	DBConn       string
	Port         string
	SecretKey    string
	AuthKey      string
	AllowOrigins string
}

var Config Environment

func Load() {
	// Solo intentamos cargar .env si existe. Si no, Go leerá las variables del sistema.
	_ = godotenv.Load()

	Config = Environment{
		Port:         getEnv("PORT", "4100"),
		DBConn:       getEnv("DB_CONN", ""),
		SecretKey:    getEnv("SECRET_KEY", ""),
		AuthKey:      getEnv("AUTH_KEY", ""),
		AllowOrigins: getEnv("ALLOW_ORIGINS", "*"),
	}

	// Validación rápida: Si no hay base de datos, el servidor no puede funcionar.
	if Config.DBConn == "" {
		log.Fatal("[!] DB_CONN no está definida en el entorno")
	}
}

// Función auxiliar simple
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

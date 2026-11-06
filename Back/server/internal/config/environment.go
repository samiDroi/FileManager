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
	UploadPath   string
	FrontendURL  string

	MailMailer      string
	MailHost        string
	MailPort        string
	MailUsername    string
	MailPassword    string
	MailEncryption  string
	MailFromAddress string
	MailFromName    string
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
		AllowOrigins: getEnv("ALLOW_ORIGINS", "http://localhost:4200"),
		UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:4200"),

		MailMailer:      getEnv("MAIL_MAILER", "smtp"),
		MailHost:        getEnv("MAIL_HOST", ""),
		MailPort:        getEnv("MAIL_PORT", "587"),
		MailUsername:    getEnv("MAIL_USERNAME", ""),
		MailPassword:    getEnv("MAIL_PASSWORD", ""),
		MailEncryption:  getEnv("MAIL_ENCRYPTION", "tls"),
		MailFromAddress: getEnv("MAIL_FROM_ADDRESS", ""),
		MailFromName:    getEnv("MAIL_FROM_NAME", "CattFileManager"),
	}

	// Validación rápida: Si no hay base de datos, el servidor no puede funcionar.
	if Config.DBConn == "" {
		log.Fatal("[!] DB_CONN no está definida en el entorno")
	}
	// Validación de seguridad: Sin secreto, no hay JWT seguro.
	if Config.SecretKey == "" {
		log.Fatal("[!] SECRET_KEY no está definida en el entorno. Imposible firmar tokens.")
	}
}

// Función auxiliar simple
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/internal/config"
	"server/internal/routes" // Importamos tu nuevo sistema de rutas
	"server/internal/store"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Cargar configuración y DB
	config.Load()
	runMigrate := flag.Bool("migrate", false, "Ejecuta las migraciones de Goose antes de iniciar")
	flag.Parse()

	if *runMigrate {
		log.Println("🛠️ Ejecutando migraciones...")
		store.RunMigrations()
		log.Println("✅ Migraciones finalizadas.")
	}
	// Inicializamos la DB (esto nos devuelve el *pgxpool.Pool)
	db := store.InitDB(config.Config.DBConn)

	e := echo.New()
	// 🌟 ACTIVA CORS PARA QUE ANGULAR PUEDA ENTRAR
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"}, // El puerto de tu Angular
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// 2. Middlewares Globales
	e.Use(middleware.RequestLogger()) // Para ver qué pasa en la consola
	e.Use(middleware.Recover())       // Para que el servidor no muera si hay un error grave

	// 3. Conectar las Rutas (Inyectamos la DB aquí)
	// Esta función para usar paths y handlers
	routes.SetupRoutes(e, db)

	// 4. ARRANCAR EL SERVIDOR (Goroutine)
	go func() {
		port := ":" + config.Config.Port
		log.Printf("Servidor corriendo en http://localhost%s", port)

		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error crítico al arrancar: %v", err)
		}
	}()

	// 5. SEÑALES DE APAGADO
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Espera aquí hasta que presiones Ctrl+C

	// 6. APAGADO ELEGANTE (Graceful Shutdown)
	log.Println("\n⚠️  Cerrando servidor de forma segura...")

	// Cerramos la conexión a la base de datos
	store.CloseDB()

	// Damos 10 segundos para que las peticiones actuales terminen
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Error durante el apagado: %v", err)
	}

	log.Println("Servidor apagado. ¡Hasta pronto!")
}

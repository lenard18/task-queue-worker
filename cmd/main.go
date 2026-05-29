package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lenard18/task-queue-worker/internal/api"
	"github.com/lenard18/task-queue-worker/internal/config"
	"github.com/lenard18/task-queue-worker/internal/database"
	"github.com/lenard18/task-queue-worker/internal/mailer"
	"github.com/lenard18/task-queue-worker/internal/models"
	"github.com/lenard18/task-queue-worker/internal/worker"
)

func main() {
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🔄  Task Queue Worker — Leonardo Gómez")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 1. Cargar configuración
	cfg := config.Load()

	// 2. Conectar a PostgreSQL
	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	// 3. Crear repositorio y mailer
	repo := models.NewJobRepository(db)
	mail := mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)

	// 4. Arrancar pool de workers en background
	pool := worker.NewPool(cfg.WorkerCount, repo, mail, cfg.PollIntervalMs)
	pool.Start()

	// 5. Configurar API HTTP
	handler := api.NewHandler(repo)
	router := api.SetupRouter(handler)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// 6. Arrancar servidor en goroutine separada
	go func() {
		log.Printf("🌐 API escuchando en http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error en servidor: %v", err)
		}
	}()

	// 7. Esperar señal de shutdown (Ctrl+C o SIGTERM de Docker/Render)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("📴 Señal de shutdown recibida, cerrando limpiamente...")

	// Dar 10 segundos para que los jobs en proceso terminen
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️  Error en shutdown del servidor: %v", err)
	}

	pool.Stop()
	log.Println("✅ Aplicación cerrada correctamente")
}

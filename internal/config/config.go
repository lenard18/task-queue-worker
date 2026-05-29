package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	// Servidor
	Port string

	// Base de datos
	DatabaseURL string

	// Worker
	WorkerCount    int // cuántos workers corren en paralelo
	PollIntervalMs int // cada cuántos ms revisa la cola

	// SMTP (email)
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

// Load lee las variables de entorno y devuelve la configuración
func Load() *Config {
	// En desarrollo carga el .env; en producción las vars ya están seteadas
	_ = godotenv.Load()

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/taskqueue?sslmode=disable"),
		WorkerCount:    getEnvInt("WORKER_COUNT", 3),
		PollIntervalMs: getEnvInt("POLL_INTERVAL_MS", 2000),
		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnvInt("SMTP_PORT", 587),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:       getEnv("SMTP_FROM", "noreply@taskqueue.dev"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("⚠️  Config: valor inválido para %s='%s', usando %d", key, v, fallback)
		return fallback
	}
	return n
}

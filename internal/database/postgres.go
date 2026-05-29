package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // driver PostgreSQL
)

// Connect abre la conexión a PostgreSQL y corre las migraciones
func Connect(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("❌ Error abriendo base de datos: %v", err)
	}

	// Pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Reintentar conexión hasta 10 veces (útil en Docker cuando Postgres tarda en arrancar)
	for i := 1; i <= 10; i++ {
		if err = db.Ping(); err == nil {
			log.Println("✅ Conectado a PostgreSQL")
			break
		}
		log.Printf("⏳ Esperando PostgreSQL... intento %d/10", i)
		time.Sleep(2 * time.Second)
		if i == 10 {
			log.Fatalf("❌ No se pudo conectar a PostgreSQL: %v", err)
		}
	}

	runMigrations(db)
	return db
}

// runMigrations ejecuta el SQL de migraciones si existe el archivo
func runMigrations(db *sql.DB) {
	migrationFile := "migrations/001_create_jobs.sql"
	sql, err := os.ReadFile(migrationFile)
	if err != nil {
		log.Printf("⚠️  No se encontró %s, asumiendo tablas ya creadas", migrationFile)
		return
	}

	if _, err := db.Exec(string(sql)); err != nil {
		log.Fatalf("❌ Error ejecutando migración: %v", err)
	}
	log.Println("✅ Migraciones aplicadas")
}

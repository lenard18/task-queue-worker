package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Status representa el estado de un job en la cola
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusFailed     Status = "failed"
)

// JobType representa el tipo de tarea a ejecutar
type JobType string

const (
	JobSendEmail   JobType = "send_email"   // enviar un email personalizado
	JobSendWelcome JobType = "send_welcome" // email de bienvenida al registrarse
	JobSendReport  JobType = "send_report"  // reporte diario de stats
	JobSendAlert   JobType = "send_alert"   // alerta urgente
)

// Job representa una tarea en la cola de procesamiento
type Job struct {
	ID          int64           `json:"id"`
	Type        JobType         `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      Status          `json:"status"`
	Priority    int             `json:"priority"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	CreatedAt   time.Time       `json:"created_at"`
	ScheduledAt time.Time       `json:"scheduled_at"`
	ProcessedAt *time.Time      `json:"processed_at,omitempty"`
	FailedAt    *time.Time      `json:"failed_at,omitempty"`
	ErrorMsg    *string         `json:"error_msg,omitempty"`
}

// ── Payloads de cada tipo de job ──────────────────────────────────────────────

// EmailPayload contiene los datos para enviar un email
type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// WelcomePayload contiene los datos del email de bienvenida
type WelcomePayload struct {
	To   string `json:"to"`
	Name string `json:"name"`
}

// ReportPayload contiene parámetros para el reporte
type ReportPayload struct {
	To     string `json:"to"`
	Period string `json:"period"` // "daily" | "weekly" | "monthly"
}

// AlertPayload contiene datos de una alerta urgente
type AlertPayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Level   string `json:"level"` // "info" | "warning" | "critical"
}

// ── Repositorio ───────────────────────────────────────────────────────────────

// JobRepository maneja todas las operaciones de base de datos para jobs
type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create inserta un nuevo job en la cola
func (r *JobRepository) Create(jobType JobType, payload interface{}, priority int) (*Job, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO jobs (type, payload, priority)
		VALUES ($1, $2, $3)
		RETURNING id, type, payload, status, priority, attempts, max_attempts,
		          created_at, scheduled_at`

	var job Job
	err = r.db.QueryRow(query, jobType, data, priority).Scan(
		&job.ID, &job.Type, &job.Payload, &job.Status,
		&job.Priority, &job.Attempts, &job.MaxAttempts,
		&job.CreatedAt, &job.ScheduledAt,
	)
	return &job, err
}

// ClaimNext obtiene el próximo job disponible de forma concurrente-segura.
// Usa SELECT FOR UPDATE SKIP LOCKED — la técnica clave de PostgreSQL para colas.
// Si dos workers corren al mismo tiempo, cada uno obtiene un job diferente.
func (r *JobRepository) ClaimNext() (*Job, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT id, type, payload, status, priority, attempts, max_attempts,
		       created_at, scheduled_at
		FROM jobs
		WHERE status = 'pending'
		  AND scheduled_at <= NOW()
		  AND attempts < max_attempts
		ORDER BY priority DESC, created_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED`

	var job Job
	err = tx.QueryRow(query).Scan(
		&job.ID, &job.Type, &job.Payload, &job.Status,
		&job.Priority, &job.Attempts, &job.MaxAttempts,
		&job.CreatedAt, &job.ScheduledAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // cola vacía, está bien
	}
	if err != nil {
		return nil, err
	}

	// Marcar como "processing" para que otro worker no lo tome
	_, err = tx.Exec(`UPDATE jobs SET status='processing', attempts=attempts+1 WHERE id=$1`, job.ID)
	if err != nil {
		return nil, err
	}

	return &job, tx.Commit()
}

// MarkDone marca un job como completado exitosamente
func (r *JobRepository) MarkDone(id int64) error {
	_, err := r.db.Exec(`
		UPDATE jobs SET status='done', processed_at=NOW() WHERE id=$1`, id)
	return err
}

// MarkFailed marca un job como fallido con mensaje de error.
// Si aún tiene intentos disponibles lo devuelve a 'pending' con retraso exponencial.
func (r *JobRepository) MarkFailed(id int64, errMsg string, attempts, maxAttempts int) error {
	if attempts < maxAttempts {
		// Retraso exponencial: 1min, 2min, 4min...
		delaySeconds := 60 * (1 << (attempts - 1))
		_, err := r.db.Exec(`
			UPDATE jobs
			SET status='pending',
			    error_msg=$1,
			    scheduled_at=NOW() + ($2 || ' seconds')::INTERVAL
			WHERE id=$3`, errMsg, delaySeconds, id)
		return err
	}
	// Sin más intentos → fallo definitivo
	_, err := r.db.Exec(`
		UPDATE jobs SET status='failed', error_msg=$1, failed_at=NOW() WHERE id=$2`,
		errMsg, id)
	return err
}

// List devuelve todos los jobs con filtro opcional por status
func (r *JobRepository) List(status string, limit int) ([]*Job, error) {
	query := `
		SELECT id, type, payload, status, priority, attempts, max_attempts,
		       created_at, scheduled_at, processed_at, failed_at, error_msg
		FROM jobs`

	args := []interface{}{}
	if status != "" {
		query += " WHERE status=$1"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		var j Job
		err := rows.Scan(
			&j.ID, &j.Type, &j.Payload, &j.Status, &j.Priority,
			&j.Attempts, &j.MaxAttempts, &j.CreatedAt, &j.ScheduledAt,
			&j.ProcessedAt, &j.FailedAt, &j.ErrorMsg,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &j)
	}
	return jobs, nil
}

// Stats devuelve conteo de jobs por estado
func (r *JobRepository) Stats() (map[string]interface{}, error) {
	rows, err := r.db.Query(`SELECT status, total, avg_seconds FROM jobs_stats`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]interface{}{}
	for rows.Next() {
		var status string
		var total int
		var avgSeconds sql.NullInt64
		if err := rows.Scan(&status, &total, &avgSeconds); err != nil {
			return nil, err
		}
		result[status] = map[string]interface{}{
			"total":       total,
			"avg_seconds": avgSeconds.Int64,
		}
	}
	return result, nil
}

package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lenard18/task-queue-worker/internal/mailer"
	"github.com/lenard18/task-queue-worker/internal/models"
)

// Worker es el motor que procesa jobs de la cola
type Worker struct {
	id       int
	repo     *models.JobRepository
	mailer   *mailer.Mailer
	interval time.Duration
	quit     chan struct{}
}

// Pool gestiona múltiples workers corriendo en paralelo
type Pool struct {
	workers []*Worker
	wg      sync.WaitGroup
}

// NewPool crea un pool de N workers
func NewPool(count int, repo *models.JobRepository, m *mailer.Mailer, pollMs int) *Pool {
	pool := &Pool{}
	interval := time.Duration(pollMs) * time.Millisecond

	for i := 1; i <= count; i++ {
		w := &Worker{
			id:       i,
			repo:     repo,
			mailer:   m,
			interval: interval,
			quit:     make(chan struct{}),
		}
		pool.workers = append(pool.workers, w)
	}
	return pool
}

// Start arranca todos los workers en goroutines independientes
func (p *Pool) Start() {
	log.Printf("🚀 Arrancando pool con %d workers", len(p.workers))
	for _, w := range p.workers {
		p.wg.Add(1)
		go func(worker *Worker) {
			defer p.wg.Done()
			worker.run()
		}(w)
	}
}

// Stop detiene todos los workers limpiamente
func (p *Pool) Stop() {
	log.Println("🛑 Deteniendo workers...")
	for _, w := range p.workers {
		close(w.quit)
	}
	p.wg.Wait()
	log.Println("✅ Todos los workers detenidos")
}

// run es el loop principal de cada worker
// Cada worker revisa la cola, toma un job y lo procesa.
// Si no hay jobs, espera el intervalo configurado.
func (w *Worker) run() {
	log.Printf("👷 Worker #%d iniciado", w.id)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.quit:
			log.Printf("👷 Worker #%d detenido", w.id)
			return
		case <-ticker.C:
			w.processNext()
		}
	}
}

// processNext intenta obtener y procesar el próximo job disponible
func (w *Worker) processNext() {
	job, err := w.repo.ClaimNext()
	if err != nil {
		log.Printf("❌ Worker #%d error obteniendo job: %v", w.id, err)
		return
	}
	if job == nil {
		return // cola vacía, ok
	}

	log.Printf("⚙️  Worker #%d procesando job #%d [%s]", w.id, job.ID, job.Type)
	start := time.Now()

	processErr := w.process(job)

	elapsed := time.Since(start).Milliseconds()

	if processErr != nil {
		log.Printf("❌ Worker #%d job #%d falló en %dms: %v", w.id, job.ID, elapsed, processErr)
		_ = w.repo.MarkFailed(job.ID, processErr.Error(), job.Attempts+1, job.MaxAttempts)
	} else {
		log.Printf("✅ Worker #%d job #%d completado en %dms", w.id, job.ID, elapsed)
		_ = w.repo.MarkDone(job.ID)
	}
}

// process ejecuta la lógica del job según su tipo
func (w *Worker) process(job *models.Job) error {
	switch job.Type {

	case models.JobSendEmail:
		var p models.EmailPayload
		if err := json.Unmarshal(job.Payload, &p); err != nil {
			return fmt.Errorf("payload inválido: %w", err)
		}
		return w.mailer.SendEmail(p.To, p.Subject, p.Body)

	case models.JobSendWelcome:
		var p models.WelcomePayload
		if err := json.Unmarshal(job.Payload, &p); err != nil {
			return fmt.Errorf("payload inválido: %w", err)
		}
		return w.mailer.SendWelcome(p.To, p.Name)

	case models.JobSendReport:
		var p models.ReportPayload
		if err := json.Unmarshal(job.Payload, &p); err != nil {
			return fmt.Errorf("payload inválido: %w", err)
		}
		stats, err := w.repo.Stats()
		if err != nil {
			return fmt.Errorf("error obteniendo stats: %w", err)
		}
		return w.mailer.SendReport(p.To, p.Period, stats)

	case models.JobSendAlert:
		var p models.AlertPayload
		if err := json.Unmarshal(job.Payload, &p); err != nil {
			return fmt.Errorf("payload inválido: %w", err)
		}
		return w.mailer.SendAlert(p.To, p.Message, p.Level)

	default:
		return fmt.Errorf("tipo de job desconocido: %s", job.Type)
	}
}

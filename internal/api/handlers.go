package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lenard18/task-queue-worker/internal/models"
)

type Handler struct {
	repo *models.JobRepository
}

func NewHandler(repo *models.JobRepository) *Handler {
	return &Handler{repo: repo}
}

// ── POST /api/jobs ────────────────────────────────────────────────────────────
// Encola un nuevo job
func (h *Handler) CreateJob(c *gin.Context) {
	var req struct {
		Type     models.JobType `json:"type"     binding:"required"`
		Payload  interface{}    `json:"payload"  binding:"required"`
		Priority int            `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar tipo de job
	validTypes := map[models.JobType]bool{
		models.JobSendEmail:   true,
		models.JobSendWelcome: true,
		models.JobSendReport:  true,
		models.JobSendAlert:   true,
	}
	if !validTypes[req.Type] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tipo inválido",
			"valid_types": []string{
				string(models.JobSendEmail),
				string(models.JobSendWelcome),
				string(models.JobSendReport),
				string(models.JobSendAlert),
			},
		})
		return
	}

	job, err := h.repo.Create(req.Type, req.Payload, req.Priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job encolado exitosamente",
		"job":     job,
	})
}

// ── GET /api/jobs ─────────────────────────────────────────────────────────────
// Lista jobs con filtro opcional por status
func (h *Handler) ListJobs(c *gin.Context) {
	status := c.Query("status") // pending | processing | done | failed
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	jobs, err := h.repo.List(status, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if jobs == nil {
		jobs = []*models.Job{}
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(jobs),
		"jobs":  jobs,
	})
}

// ── GET /api/stats ────────────────────────────────────────────────────────────
// Dashboard de métricas de la cola
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.repo.Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"queue_stats": stats,
		"description": "Sistema de cola de tareas con PostgreSQL y Go",
	})
}

// ── GET /api/health ───────────────────────────────────────────────────────────
// Health check para Render
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "task-queue-worker",
	})
}

// ── POST /api/demo ────────────────────────────────────────────────────────────
// Encola 5 jobs de ejemplo para ver el sistema en acción
func (h *Handler) EnqueueDemo(c *gin.Context) {
	jobs := []struct {
		jobType  models.JobType
		payload  interface{}
		priority int
	}{
		{
			models.JobSendWelcome,
			models.WelcomePayload{To: "nuevo@demo.com", Name: "Leonardo"},
			5,
		},
		{
			models.JobSendEmail,
			models.EmailPayload{
				To:      "cliente@demo.com",
				Subject: "Tu pedido fue confirmado 🛒",
				Body:    "<h2>¡Pedido confirmado!</h2><p>Tu pedido #1234 está en camino.</p>",
			},
			3,
		},
		{
			models.JobSendAlert,
			models.AlertPayload{
				To:      "admin@demo.com",
				Message: "El servidor de pagos respondió lento (>2s)",
				Level:   "warning",
			},
			10, // prioridad alta
		},
		{
			models.JobSendReport,
			models.ReportPayload{To: "manager@demo.com", Period: "daily"},
			1,
		},
		{
			models.JobSendAlert,
			models.AlertPayload{
				To:      "ops@demo.com",
				Message: "Base de datos al 95% de capacidad",
				Level:   "critical",
			},
			10, // máxima prioridad
		},
	}

	var created []*models.Job
	for _, j := range jobs {
		job, err := h.repo.Create(j.jobType, j.payload, j.priority)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		created = append(created, job)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "5 jobs de demo encolados. Los workers los procesarán en segundos.",
		"jobs":    created,
	})
}

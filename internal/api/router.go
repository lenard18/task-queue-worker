package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas las rutas de la API
func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	// CORS simple para poder probar desde el browser
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", h.HealthCheck)

	// Página de bienvenida con instrucciones
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service":     "🔄 Task Queue Worker",
			"description": "Sistema de cola de tareas con Go + PostgreSQL",
			"author":      "Leonardo Gómez — github.com/lenard18",
			"endpoints": gin.H{
				"POST /api/jobs":    "Encolar un nuevo job",
				"GET  /api/jobs":    "Listar jobs (query: ?status=pending|done|failed|processing&limit=50)",
				"GET  /api/stats":   "Ver estadísticas de la cola",
				"POST /api/demo":    "Encolar 5 jobs de ejemplo para ver el sistema en acción",
				"GET  /health":      "Health check",
			},
			"job_types": gin.H{
				"send_email":   gin.H{"payload": gin.H{"to": "email", "subject": "text", "body": "html"}},
				"send_welcome": gin.H{"payload": gin.H{"to": "email", "name": "text"}},
				"send_report":  gin.H{"payload": gin.H{"to": "email", "period": "daily|weekly|monthly"}},
				"send_alert":   gin.H{"payload": gin.H{"to": "email", "message": "text", "level": "info|warning|critical"}},
			},
		})
	})

	api := r.Group("/api")
	{
		api.POST("/jobs", h.CreateJob)
		api.GET("/jobs", h.ListJobs)
		api.GET("/stats", h.GetStats)
		api.POST("/demo", h.EnqueueDemo)
	}

	return r
}

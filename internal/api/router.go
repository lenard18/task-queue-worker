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

	// Página de bienvenida visual en HTML
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(homeHTML))
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

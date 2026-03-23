package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

// NewHealthHandler builds the HTTP handler for health checks.
func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

// Health serves the liveness endpoint.
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.Health())
}

// Ready serves the readiness endpoint.
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.Ready())
}

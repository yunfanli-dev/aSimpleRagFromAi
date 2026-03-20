package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.Health())
}

func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.Ready())
}

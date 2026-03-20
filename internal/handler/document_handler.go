package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/service"
	"github.com/yunfanli-dev/aSimpleRagFromAi/pkg/response"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(service *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) List(c *gin.Context) {
	kbID := c.Param("id")
	response.JSON(c, http.StatusOK, h.service.List(kbID))
}

func (h *DocumentHandler) Create(c *gin.Context) {
	response.JSON(c, http.StatusAccepted, gin.H{
		"message": "document ingest endpoint scaffolded",
	})
}

func (h *DocumentHandler) Reindex(c *gin.Context) {
	response.JSON(c, http.StatusAccepted, gin.H{
		"message":     "document reindex endpoint scaffolded",
		"document_id": c.Param("id"),
	})
}

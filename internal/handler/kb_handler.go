package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/service"
	"github.com/yunfanli-dev/aSimpleRagFromAi/pkg/response"
)

type KnowledgeBaseHandler struct {
	service *service.KnowledgeBaseService
}

func NewKnowledgeBaseHandler(service *service.KnowledgeBaseService) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{service: service}
}

func (h *KnowledgeBaseHandler) List(c *gin.Context) {
	response.JSON(c, http.StatusOK, h.service.List())
}

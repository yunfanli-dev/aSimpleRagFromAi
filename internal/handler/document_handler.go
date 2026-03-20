package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
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
	items, err := h.service.List(c.Request.Context(), kbID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, items)
}

func (h *DocumentHandler) Create(c *gin.Context) {
	var input domain.CreateDocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	input.KnowledgeBaseID = c.Param("id")

	item, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, item)
}

func (h *DocumentHandler) Get(c *gin.Context) {
	item, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		if err == repository.ErrNotFound {
			response.Error(c, http.StatusNotFound, "document not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

func (h *DocumentHandler) ListChunks(c *gin.Context) {
	items, err := h.service.ListChunks(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, items)
}

func (h *DocumentHandler) Reindex(c *gin.Context) {
	response.JSON(c, http.StatusAccepted, gin.H{
		"message":     "document reindex endpoint scaffolded",
		"document_id": c.Param("id"),
	})
}

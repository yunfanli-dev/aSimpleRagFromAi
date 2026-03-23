package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/repository"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/service"
	"github.com/yunfanli-dev/aSimpleRagFromAi/pkg/response"
)

type KnowledgeBaseHandler struct {
	service *service.KnowledgeBaseService
}

// NewKnowledgeBaseHandler builds the HTTP handler for knowledge base operations.
func NewKnowledgeBaseHandler(service *service.KnowledgeBaseService) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{service: service}
}

// Create creates one knowledge base from the request body.
func (h *KnowledgeBaseHandler) Create(c *gin.Context) {
	var input domain.CreateKnowledgeBaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, item)
}

// Get returns one knowledge base by route ID.
func (h *KnowledgeBaseHandler) Get(c *gin.Context) {
	item, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		if err == repository.ErrNotFound {
			response.Error(c, http.StatusNotFound, "knowledge base not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, item)
}

// List returns all knowledge bases.
func (h *KnowledgeBaseHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, items)
}

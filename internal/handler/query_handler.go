package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/domain"
	"github.com/yunfanli-dev/aSimpleRagFromAi/internal/service"
	"github.com/yunfanli-dev/aSimpleRagFromAi/pkg/response"
)

type QueryHandler struct {
	service *service.QueryService
}

func NewQueryHandler(service *service.QueryService) *QueryHandler {
	return &QueryHandler{service: service}
}

func (h *QueryHandler) Query(c *gin.Context) {
	var req domain.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, h.service.Ask(req))
}

func (h *QueryHandler) Debug(c *gin.Context) {
	var req domain.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Debug = true
	response.JSON(c, http.StatusOK, h.service.Ask(req))
}

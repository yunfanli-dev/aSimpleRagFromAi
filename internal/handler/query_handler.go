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

// NewQueryHandler builds the HTTP handler for query endpoints.
func NewQueryHandler(service *service.QueryService) *QueryHandler {
	return &QueryHandler{service: service}
}

// Query executes the standard retrieval-and-answer flow.
func (h *QueryHandler) Query(c *gin.Context) {
	var req domain.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.service.Ask(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, resp)
}

// Debug executes the query flow with debug output enabled.
func (h *QueryHandler) Debug(c *gin.Context) {
	var req domain.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Debug = true
	resp, err := h.service.Ask(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, resp)
}

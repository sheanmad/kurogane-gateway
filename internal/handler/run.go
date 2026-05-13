package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kurogane/gateway/internal/model"
	"github.com/kurogane/gateway/internal/proxy"
)

type RunHandler struct {
	proxy *proxy.EngineProxy
}

func NewRunHandler(p *proxy.EngineProxy) *RunHandler {
	return &RunHandler{proxy: p}
}

func (h *RunHandler) CreateRun(c *gin.Context) {
	var req model.CreateRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	result, err := h.proxy.CreateRun(&req)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.ErrorResponse{
			Error: fmt.Sprintf("engine error: %v", err),
			Code:  http.StatusBadGateway,
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *RunHandler) GetRun(c *gin.Context) {
	id := c.Param("id")

	result, err := h.proxy.GetRun(id)
	if err != nil {
		if err.Error() == "run not found" {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "Run not found",
				Code:  http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusBadGateway, model.ErrorResponse{
			Error: fmt.Sprintf("engine error: %v", err),
			Code:  http.StatusBadGateway,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *RunHandler) ListRuns(c *gin.Context) {
	runs, err := h.proxy.ListRuns()
	if err != nil {
		c.JSON(http.StatusBadGateway, model.ErrorResponse{
			Error: fmt.Sprintf("engine error: %v", err),
			Code:  http.StatusBadGateway,
		})
		return
	}

	c.JSON(http.StatusOK, runs)
}

func (h *RunHandler) StreamRunEvents(c *gin.Context) {
	id := c.Param("id")

	ch, err := h.proxy.StreamRunEvents(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.ErrorResponse{
			Error: fmt.Sprintf("engine error: %v", err),
			Code:  http.StatusBadGateway,
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent(event.EventType, event.Data)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
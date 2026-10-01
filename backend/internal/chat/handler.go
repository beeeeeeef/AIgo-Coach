package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler HTTP 入口
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleChat POST /api/chat
func (h *Handler) HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ChatResponse{
			Result: false,
			Error:  "请求参数错误: " + err.Error(),
		})
		return
	}

	data, err := h.service.Chat(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ChatResponse{
			Result: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{
		Result: true,
		Data:   data,
	})
}

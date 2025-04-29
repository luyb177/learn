package controller

import (
	"github.com/gin-gonic/gin"
	"learn/Grab_seat/service"
)

type SseController struct {
	se service.SseService
}

func NewSseController(se service.SseService) *SseController {
	return &SseController{se: se}
}

// SseEvent 处理 SSE 事件
// @Summary 处理 SSE 事件
// @Description 处理服务器发送事件（SSE）
// @Tags SSE 服务
// @Accept json
// @Produce text/event-stream
// @Success 200 {string} string "SSE 事件流"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /sse [get]
func (sc *SseController) SseEvent(c *gin.Context) {
	sc.se.SSEServer(c)
}

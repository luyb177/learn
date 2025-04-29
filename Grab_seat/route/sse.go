package route

import (
	"github.com/gin-gonic/gin"
	"learn/Grab_seat/controller"
	"learn/Grab_seat/middleware"
)

type SseRoute struct {
	sc *controller.SseController
}

func NewSseRoute(sc *controller.SseController) *SseRoute {
	return &SseRoute{
		sc: sc,
	}
}
func (sr *SseRoute) SseGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	r.GET("/sse", sr.sc.SseEvent)
}

package route

import (
	"github.com/gin-gonic/gin"
	"learn/Grab_seat/controller"
	"learn/Grab_seat/middleware"
)

type GrabRoot struct {
	gc *controller.GrabController
}

func NewGrabRoot(gc *controller.GrabController) *GrabRoot {
	return &GrabRoot{
		gc: gc,
	}
}

func (gr *GrabRoot) GrabGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	Group := r.Group("/grab")
	{
		Group.POST("/send", gr.gc.Send)
	}
}

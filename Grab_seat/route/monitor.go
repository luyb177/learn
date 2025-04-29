package route

import (
	"github.com/gin-gonic/gin"
	"learn/Grab_seat/controller"
	"learn/Grab_seat/middleware"
)

type MonitorRoot struct {
	mc *controller.MonitorController
}

func NewMonitorRoot(mc *controller.MonitorController) *MonitorRoot {
	return &MonitorRoot{
		mc: mc,
	}
}

func (mr *MonitorRoot) NewGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	monitorGroup := r.Group("/monitor")
	{
		monitorGroup.POST("/one", mr.mc.CheckOneSeat)
	}
}

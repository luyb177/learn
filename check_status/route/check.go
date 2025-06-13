package route

import (
	"github.com/gin-gonic/gin"
	"learn/check_status/controller"
	"learn/check_status/middleware"
)

type CheckRoute struct {
	cc *controller.CheckController
}

func NewCheckRoute(cc *controller.CheckController) *CheckRoute {
	return &CheckRoute{cc: cc}
}

func (cr *CheckRoute) NewGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	CheckGroup := r.Group("/check")
	{
		CheckGroup.GET("/get/event", cr.cc.GetEvent)
		CheckGroup.POST("/add/user", cr.cc.AddUser)
		CheckGroup.GET("/delete/user", cr.cc.DeleteUser)
		CheckGroup.POST("/alter/qq", cr.cc.AlterQQ)
		CheckGroup.POST("/set/seat", cr.cc.SetSeatRecord)
		CheckGroup.POST("/get/seat", cr.cc.GetSeatRecord)
		CheckGroup.POST("/alter/seat", cr.cc.AlterSeatRecord)
	}

}

package controller

import (
	"github.com/gin-gonic/gin"
	"learn/Grab_seat/api/request"
	"learn/Grab_seat/api/response"
	"learn/Grab_seat/service"
	"net/http"
)

type GrabController struct {
	gs service.GrabService
}

func NewGrabController(gs service.GrabService) *GrabController {
	return &GrabController{
		gs: gs,
	}
}

// Send 处理发送请求
// @Summary 发送消息
// @Description 发送消息到指定的抓取服务
// @Tags 抓取服务
// @Accept json
// @Produce json
// @Param grab body request.Grab true "抓取请求"
// @Success 200 {object} response.Response "发送成功"
// @Failure 400 {object} response.Response "解析失败或发送失败"
// @Router /grab/send [post]
func (gc *GrabController) Send(c *gin.Context) {
	var grab request.Grab
	if err := c.ShouldBindJSON(&grab); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "解析失败",
		})
		return
	}

	err := gc.gs.SendMsg(&grab)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "发送失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "发送成功，正在为你处理请求",
	})
}

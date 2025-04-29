package controller

import (
	"github.com/gin-gonic/gin"
	"learn/Grab_seat/api/request"
	"learn/Grab_seat/api/response"
	"learn/Grab_seat/service"
	"net/http"
)

type MonitorController struct {
	ms service.MonitorService
}

func NewMonitorController(ms service.MonitorService) *MonitorController {
	return &MonitorController{
		ms: ms,
	}
}

// CheckOneSeat 检测一个座位的状态
// @Summary 检测一个座位的状态
// @Description 检测指定座位的状态
// @Tags 监控服务
// @Accept json
// @Produce json
// @Param grab body request.Grab true "检测请求"
// @Success 200 {object} response.Response "检测成功"
// @Failure 400 {object} response.Response "解析失败或检测失败"
// @Router /monitor/one [post]
func (mc *MonitorController) CheckOneSeat(c *gin.Context) {
	var grab request.Grab
	if err := c.ShouldBindJSON(&grab); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "解析失败",
		})
		return
	}
	res, err := mc.ms.CheckOneSeat(&grab)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "检测失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "检测成功",
		Data:    res,
	})
}

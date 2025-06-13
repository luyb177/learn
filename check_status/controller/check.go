package controller

import (
	"github.com/gin-gonic/gin"
	"learn/check_status/api/request"
	"learn/check_status/api/response"
	"learn/check_status/service"
	"net/http"
	"strconv"
)

type CheckController struct {
	cs service.CheckService
}

func NewCheckController(cs service.CheckService) *CheckController {
	return &CheckController{cs: cs}
}

// GetEvent 分页获取发邮件时间
// @Summary 分页获取发邮件时间
// @Description 分页获取发邮件时间
// @Tags events
// @Accept json
// @Produce json
// @Param pn query int true "页码"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "pn解析失败"
// @Failure 404 {object} response.Response "获取失败"
// @Router /events [get]
func (cc *CheckController) GetEvent(c *gin.Context) {
	pnStr := c.Query("pn")
	pn, err := strconv.Atoi(pnStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "pn解析失败",
		})
		return
	}
	events, err := cc.cs.GetEvent(pn)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "获取失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "获取成功",
		Data:    events,
	})
}

// AddUser 添加用户
// @Summary 添加用户
// @Description 添加新用户
// @Tags users
// @Accept json
// @Produce json
// @Param user body request.User true "用户信息"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "参数解析失败"
// @Router /users [post]
func (cc *CheckController) AddUser(c *gin.Context) {
	var user request.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	err = cc.cs.AddUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "添加失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "添加成功",
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据用户名删除用户
// @Tags users
// @Accept json
// @Produce json
// @Param name query string true "用户名"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "删除失败"
// @Router /users [delete]
func (cc *CheckController) DeleteUser(c *gin.Context) {
	name := c.Query("name")
	err := cc.cs.DeleteUser(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "删除成功",
	})
}

// AlterQQ 修改用户的QQ
// @Summary 修改用户的QQ
// @Description 根据用户信息修改用户的QQ
// @Tags users
// @Accept json
// @Produce json
// @Param user body request.User true "用户信息"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "参数解析失败"
// @Failure 404 {object} response.Response "修改失败"
// @Router /users/qq [post]
func (cc *CheckController) AlterQQ(c *gin.Context) {
	var user request.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	err = cc.cs.AlterQQ(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "修改失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "修改成功",
	})
}

// SetSeatRecord 设置预约消息
// @Summary 设置预约消息
// @Description 根据传入的预约信息设置预约消息
// @Tags seats
// @Accept json
// @Produce json
// @Param grab body request.Grab true "预约信息"
// @Success 200 {object} response.Response "设置预约消息成功"
// @Failure 400 {object} response.Response "参数解析失败 或 设置预约消息失败"
// @Router /set/seat [post]
func (cc *CheckController) SetSeatRecord(c *gin.Context) {
	var grab request.Grab
	err := c.ShouldBindJSON(&grab)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	err = cc.cs.SetSeatRecord(&grab)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "设置预约消息失败",
			Data:    nil,
		})
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "设置预约消息成功",
		Data:    nil,
	})
}

// GetSeatRecord 获取预约记录
// @Summary 获取预约记录
// @Description 根据请求参数获取用户的预约记录信息
// @Tags seats
// @Accept json
// @Produce json
// @Param req body request.GetRecordReq true "查询条件"
// @Success 200 {object} response.Response "获取预约消息成功"
// @Failure 400 {object} response.Response "参数解析失败 或 获取预约消息失败"
// @Router /get/seat [post]
func (cc *CheckController) GetSeatRecord(c *gin.Context) {
	var req request.GetRecordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}
	info, err := cc.cs.GetSeatRecord(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "获取预约消息失败",
			Data:    nil,
		})
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "获取预约消息成功",
		Data:    info,
	})
}

// AlterSeatRecord 修改预约记录
// @Summary 修改预约记录
// @Description 根据传入的预约信息修改已有的预约记录
// @Tags seats
// @Accept json
// @Produce json
// @Param grab body request.Grab true "预约信息"
// @Success 200 {object} response.Response "修改预约消息成功"
// @Failure 400 {object} response.Response "参数解析失败 或 修改预约消息失败"
// @Router /alter/seat [post]
func (cc *CheckController) AlterSeatRecord(c *gin.Context) {
	var grab request.Grab
	err := c.ShouldBindJSON(&grab)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	err = cc.cs.AlterSeatRecord(&grab)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "修改预约消息失败",
		})
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "修改预约消息成功",
	})
}

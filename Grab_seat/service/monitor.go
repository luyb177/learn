package service

import (
	"encoding/json"
	"fmt"
	"io"
	"learn/Grab_seat/api/request"
	"learn/Grab_seat/api/response"
	"learn/Grab_seat/client"
	"learn/Grab_seat/dao"
	"learn/Grab_seat/model"
	"learn/Grab_seat/tool"
	"log"
	"net/http"
	"time"
)

type MonitorService interface {
	CheckOneSeatStatus(grabInfo *model.GrabInfo) bool
	CheckOneSeat(grab *request.Grab) (*response.Res, error)
}

type MonitorServiceImpl struct {
	client *client.Client
	Dao    dao.GrabDAO
}

func NewMonitorServiceImpl(client *client.Client, Dao dao.GrabDAO) *MonitorServiceImpl {
	return &MonitorServiceImpl{
		client: client,
		Dao:    Dao,
	}
}

// CheckSeatStatus 根据时间检测全部座位状态
func (msr *MonitorServiceImpl) CheckSeatStatus(grabInfo *model.GrabInfo, CheckFromNow bool) (*response.Res, error) {
	var temp = *grabInfo
	if CheckFromNow { //检测当前时间往后的座位状态
		now := time.Now()
		nowHour := now.Hour()
		nowMinute := now.Minute()
		temp.FrStart = fmt.Sprintf("%d%%3A%d", nowHour, nowMinute)
		temp.FrEnd = "22%3A00"
	}
	path := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=8&display=fp&md=d&room_id=" + temp.RoomId +
		"&purpose=&selectOpenAty=&cld_name=default&date=" + temp.Date + "&fr_start=" + temp.FrStart + "&fr_end=" + temp.FrEnd + "&act=get_rsv_sta&_=" + temp.TimeMs
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res, err := msr.client.Client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var resP response.Res
	err = json.Unmarshal(body, &resP)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &resP, nil
}

// CheckOneSeatStatusInfo CheckSeatStatusOne 检测一个座位状态
func (msr *MonitorServiceImpl) CheckOneSeatStatusInfo(grabInfo *model.GrabInfo) (*response.Res, error) {

	path := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=8&display=fp&md=d&dev_id=" + grabInfo.DevId + "&room_id=" + grabInfo.RoomId +
		"&purpose=&selectOpenAty=&cld_name=default&date=" + grabInfo.Date + "&fr_start=" + grabInfo.FrStart + "&fr_end=" + grabInfo.FrEnd + "&act=get_rsv_sta&_=" + grabInfo.TimeMs
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res, err := msr.client.Client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var resP response.Res
	err = json.Unmarshal(body, &resP)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &resP, nil
}

func (msr *MonitorServiceImpl) CheckOneSeatStatus(grabInfo *model.GrabInfo) bool {
	res, err := msr.CheckOneSeatStatusInfo(grabInfo)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	for _, v := range res.Data {
		for _, v1 := range v.Ts {
			return v1.Occupy
		}
	}
	// 未被占用
	return false
}

func (msr *MonitorServiceImpl) CheckOneSeat(grab *request.Grab) (*response.Res, error) {
	grabInfo := tool.GetParameters(grab)
	devId, roomId, err := msr.Dao.FindSeatId(grab.Seat)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	grabInfo.DevId = devId
	grabInfo.RoomId = roomId
	res, err := msr.CheckOneSeatStatusInfo(grabInfo)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return res, nil
}

package service

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"learn/check_status/api/request"
	"learn/check_status/client"
	"learn/check_status/dao"
	"learn/check_status/model"
	"learn/check_status/pb"
	"learn/check_status/tool"
	"log"
	"net/url"
	"time"
)

const (
	FirstFloorAtrium = 101699187 // 一楼中庭的roomId
	FirstFloorOpen   = 101699179 // 一楼开敞roomId
	SecondFloorOpen  = 101699189
	SecondFloorBooth = 101699191 // 二楼卡座roomId
)

type CheckService interface {
	GetEvent(pn int) ([]model.Event, error)
	AddUser(user *request.User) error
	DeleteUser(name string) error
	AlterQQ(user *request.User) error
}
type Check struct {
	client *client.Client
	// 新增上下文相关字段
	ctx    context.Context    // 主上下文
	cancel context.CancelFunc // 取消函数
	Dao    dao.CheckDAO
	Mail   *tool.Mail
}

func NewCheck(client *client.Client, Dao dao.CheckDAO, Mail *tool.Mail) *Check {
	// 登录
	client.Login()

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	var check = Check{
		client: client,
		ctx:    ctx,
		cancel: cancel,
		Dao:    Dao,
		Mail:   Mail,
	}
	// 开启定时任务
	go check.StartCheck()
	return &check
}

// StartCheck 所有的定时任务
func (c *Check) StartCheck() {
	cr := cron.New(cron.WithSeconds())
	// 30分钟登录一次
	_, err := cr.AddFunc("@every 30m", c.client.Login)
	if err != nil {
		log.Printf("重新登录失败: %v", err)
		return
	}

	//c.CheckSeat()
	// 每10分钟检查一次
	_, err = cr.AddFunc("@every 10m", c.CheckSeat)
	if err != nil {
		log.Printf("检查座位失败: %v", err)
		return
	}
	cr.Start()

	// 阻塞程序
	select {
	case <-c.ctx.Done():
		cr.Stop()
		c.cancel()
		log.Println("服务停止")
	}
}

// CheckSeat 检查座位的人
// 1、获取当前时间
// 2、判断当前时间是否在8:00-22:00之间
func (c *Check) CheckSeat() {
	// 获取当前时间

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	if now.Hour() >= 8 && now.Hour() <= 22 {
		// 获取座位信息
		c.GetSeatInfo()
	} else {
		log.Println("当前时间不在8:00-22:00之间")
		return
	}
}

// GetSeatInfo 获取座位信息
// 1、检测当前时间-22:00的座位信息
func (c *Check) GetSeatInfo() {
	// 当前时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	date := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	frStart := fmt.Sprintf("%d%%3A%d", now.Hour(), now.Minute())
	frEnd := "22%3A00"
	timeMS := now.UnixNano() / int64(time.Millisecond)

	TwentyTwo := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, now.Location())
	expiredTime := int(TwentyTwo.Sub(now).Seconds())
	// 循环获取四个地方的座位信息
	for i := 0; i < 4; i++ {
		var Res pb.Res
		roomId := c.GetRoomId(i)
		fmt.Println(date, frStart, frEnd)
		fmt.Printf("正在查询%d的座位信息\n", i)
		path := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/device.aspx?byType=devcls&classkind=8&display=fp&md=d&room_id=" + fmt.Sprintf("%d", roomId) +
			"&purpose=&selectOpenAty=&cld_name=default&date=" + date + "&fr_start=" + frStart + "&fr_end=" + frEnd + "&act=get_rsv_sta&_=" + fmt.Sprintf("%d", timeMS)

		res, err := c.client.Client.Get(path)
		if err != nil {
			log.Printf("查询%d的座位信息失败: %v", i, err)
			continue
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("读取body失败")
			continue
		}

		err = protojson.UnmarshalOptions{
			DiscardUnknown: true, // 忽略未定义字段
		}.Unmarshal(body, &Res)
		if err != nil {
			log.Printf("解析body失败")
			continue
		}
		// 查询成功
		if Res.Ret == 1 {
			log.Printf("查询%d的座位信息成功,正在检索中", i)

			for _, info := range Res.Data {
				for _, T := range info.Ts {
					// 检查座位的人名
					// url编码
					owner := url.QueryEscape(T.Owner)

					if c.Dao.IsExist(owner) {
						// 有

						if !c.Dao.IsMarked(owner, info.Title, T.Start, T.End) {
							// 未被标记
							// 添加次数
							c.Dao.AddCount(owner)
							// 获取添加后的次数
							count := c.Dao.GetCount(owner)
							// 获取qq
							qq := c.Dao.GetQQ(owner)
							err := c.Mail.SendEmailByQQEmail(qq, info.RoomName, info.Title, T.Owner, T.Start, T.End, count)
							if err != nil {
								log.Printf("发送邮件失败%v\n", err)
								return
							}
							//标记一下
							//加一个开始的参数 可以避免重复发送邮件 再加一个end更好
							//失效时间是当前-22:00
							err = c.Dao.AddMark(owner, info.Title, T.Start, T.End, time.Duration(expiredTime)*time.Second)
							if err != nil {
								log.Printf("标记失败%v\n", err)
								return
							}
							var event = model.Event{
								Name:      T.Owner,
								StartTime: T.Start,
								EndTime:   T.End,
							}
							err = c.Dao.SaveEvent(&event)
							if err != nil {
								log.Printf("保存到数据库失败：%v\n", err)
							}
						} else {
							// 已标记
							log.Printf("%s在%s到%s的预约已被标记", T.Owner, T.Start, T.End)
						}
					} else {
						log.Printf("%s并不在数据库中", T.Owner)
					}
				}
			}

		} else if Res.Ret == 0 {
			log.Printf("发生错误:%s", Res.Msg)
			continue
		} else if Res.Ret == -1 {
			log.Println("cookie失效，正在为您重新登录")
			c.ReTry()
		}

		// 在for循环中defer会泄露
		res.Body.Close()
	}
}

// ReTry 重新登录并重新获取
func (c *Check) ReTry() {
	c.client.Login()
	c.GetSeatInfo()
}

// GetRoomId 获取roomId
func (c *Check) GetRoomId(i int) int {
	switch i {
	case 0:
		return FirstFloorAtrium
	case 1:
		return FirstFloorOpen
	case 2:
		return SecondFloorOpen
	case 3:
		return SecondFloorBooth
	default:
		return -1
	}
}

// GetEvent 分页拉取发邮箱事件
func (c *Check) GetEvent(pn int) ([]model.Event, error) {
	return c.Dao.GetEvent(pn)
}

// AddUser 添加用户
func (c *Check) AddUser(user *request.User) error {
	// url 编码
	UrlName := url.QueryEscape(user.Name)
	return c.Dao.AddUser(UrlName, user.QQ)
}

// DeleteUser 删除用户
func (c *Check) DeleteUser(name string) error {
	UrlName := url.QueryEscape(name)
	return c.Dao.DeleteUser(UrlName)
}

// AlterQQ 修改用户的QQ
func (c *Check) AlterQQ(user *request.User) error {
	// HSet 在存在的时候会覆盖
	UrlName := url.QueryEscape(user.Name)
	return c.Dao.AddUser(UrlName, user.QQ)
}

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
	SetSeatRecord(grab *request.Grab) error
	GetSeatRecord(req *request.GetRecordReq) (*model.SeatInfo, error)
	AlterSeatRecord(grab *request.Grab) error
}
type Check struct {
	client *client.Client
	// 新增上下文相关字段
	ctx     context.Context    // 主上下文
	cancel  context.CancelFunc // 取消函数
	Dao     dao.CheckDAO
	GrabDao dao.GrabDAO
	Mail    *tool.Mail
}

func NewCheck(client *client.Client, Dao dao.CheckDAO, Mail *tool.Mail, GrabDao dao.GrabDAO) *Check {
	// 登录
	client.Login()

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	var check = Check{
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
		Dao:     Dao,
		GrabDao: GrabDao,
		Mail:    Mail,
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

	// 每1小时获取座位
	_, err = cr.AddFunc("@every 1h", c.GetSeat)
	if err != nil {
		log.Printf("获取座位失败: %v", err)
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
func (c *Check) SetSeatRecord(grab *request.Grab) error {
	info := model.SeatInfo{
		Seat:  grab.Seat,
		Start: grab.Start,
		End:   grab.End,
		Date:  grab.Date,
	}
	return c.Dao.SetSeatRecord(&info)
}

func (c *Check) GetSeatRecord(req *request.GetRecordReq) (*model.SeatInfo, error) {
	return c.Dao.GetSeatRecord(req.Date)
}

func (c *Check) AlterSeatRecord(grab *request.Grab) error {
	{
		info := model.SeatInfo{
			Seat:  grab.Seat,
			Start: grab.Start,
			End:   grab.End,
			Date:  grab.Date,
		}
		return c.Dao.SetSeatRecord(&info)
	}
}

// GetSeat 获取明天座位
// 1、判断时间距离18:00是否小于1小时，若大于1小时则返回，小于1小时则睡眠至17:58
// 2、按照预先设置的座位信息进行预约
func (c *Check) GetSeat() {

	// 第一步：判断时间
	now := time.Now()
	Six := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	sub := Six.Sub(now)
	if sub > 1*time.Hour {
		log.Printf("[INFO] 现在时间是%v,距离18:00还有%v,请等待\n", now, sub)
		return
	}
	// 如果18:00已经过了
	if sub < 0 {
		if now.Before(Six.Add(2 * time.Minute)) {
			log.Printf("[WARN] 已过18:00但仍在容忍范围，立即尝试预约")
			_ = c.BookSeat()
			return
		}
		log.Printf("[INFO] 现在时间是%v,已经错过预约时间%v\n", now, Six)
		return
	}

	// 计算时间睡眠
	sleepDuration := sub - 30*time.Second // 提前30 醒来
	if sleepDuration > 0 {
		// 使用context控制安全休眠
		ctx, cancel := context.WithTimeout(context.Background(), sleepDuration)
		defer cancel()
		select {
		case <-ctx.Done(): // 正常唤醒
		case <-c.ctx.Done(): // 全局上下文取消（优雅退出）
			log.Printf("[INFO] process canceled")
			return
		}

	}

	// 进入检查阶段
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		now = time.Now()
		select {
		case <-ticker.C:
			if now.Equal(Six) || now.After(Six) {
				err := c.BookSeat()
				if err != nil {
					log.Printf("[ERROR] 预约失败%v\n", err)
					return
				}
				log.Println("[INFO] 预约成功")
				return
			}
		case <-c.ctx.Done():
			log.Printf("[INFO] process canceled")
			return

		}
	}

}

func (c *Check) BookSeat() error {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	data := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), tomorrow.Day())
	// 获取座位信息和时间
	info, err := c.Dao.GetSeatRecord(data)
	if err != nil {
		log.Printf("[ERROR] 获取座位信息失败%v\n", err)
		return err
	}

	//获取参数
	devId, roomId, err := c.GrabDao.FindSeatId(info.Seat)
	if err != nil {
		log.Printf("[ERROR] 获取devId和roomId失败%v\n", err)
		return err
	}
	grabInfo := tool.GetParameters(info)
	grabInfo.DevId = devId
	grabInfo.RoomId = roomId

	path := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?dialogid=&dev_id=" + grabInfo.DevId + "&lab_id=&kind_id=&room_id=&type=dev&prop=&test_id=&term=&Vnumber=&classkind=&test_name=&start=" + grabInfo.Start + "&end=" + grabInfo.End + "&start_time=" + grabInfo.StartTime + "&end_time=" + grabInfo.EndTime + "&up_file=&memo=&act=set_resv&_=" + grabInfo.TimeMs

	respon, err := c.client.Client.Get(path)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer respon.Body.Close()
	body, _ := io.ReadAll(respon.Body)
	var res pb.Res
	err = protojson.UnmarshalOptions{
		DiscardUnknown: true, // 忽略未定义字段
	}.Unmarshal(body, &res)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	fmt.Println(res)
	if res.Ret == 1 {
		log.Printf("[INFO] 用户%d预约%s成功", 1, grabInfo.Seat)
		return nil
	} else if res.Ret == -1 {
		// cookie 失效
		c.client.Login()
		// 重新登录
		c.GetSeat()
	} else if res.Ret == 0 {
		log.Println(res.Msg)
	}
	return nil
}

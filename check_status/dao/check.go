package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"learn/check_status/model"
	"time"
)

type CheckDAO interface {
	AddUser(name string, qq string) error
	IsExist(name string) bool
	AddMark(name, seat, start, end string, expireTime time.Duration) error
	IsMarked(name, seat, start, end string) bool
	GetQQ(name string) string
	AddCount(name string)
	GetCount(name string) string
	SaveEvent(event *model.Event) error
	GetEvent(pn int) ([]model.Event, error)
	DeleteUser(name string) error
	GetSeatRecord(date string) (*model.SeatInfo, error)
	SetSeatRecord(info *model.SeatInfo) error
	SetOccupied(date string) error
	GetOccupied(date string) (bool, error)
}

type CheckDAOImpl struct {
	rdb *redis.Client
	db  *gorm.DB
}

func NewCheckDAO(rdb *redis.Client, db *gorm.DB) *CheckDAOImpl {
	return &CheckDAOImpl{
		rdb: rdb,
		db:  db,
	}
}

// AddUser 添加检查的用户
func (dao *CheckDAOImpl) AddUser(name string, qq string) error {
	_, err := dao.rdb.HSet(context.Background(), "name_qq", name, qq).Result()
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser 删除要检查的用户
func (dao *CheckDAOImpl) DeleteUser(name string) error {
	_, err := dao.rdb.HDel(context.Background(), "name_qq", name).Result()
	if err != nil {
		return err
	}
	return nil
}

// IsExist 判断该目标用户是否存在
func (dao *CheckDAOImpl) IsExist(name string) bool {
	return dao.rdb.HExists(context.Background(), "name_qq", name).Val()
}

// AddMark 添加标记
func (dao *CheckDAOImpl) AddMark(name, seat, start, end string, expireTime time.Duration) error {
	_, err := dao.rdb.SetEX(context.Background(), "name:"+name+":"+seat+":"+start+" ~ "+end, 1, expireTime).Result()
	if err != nil {
		return err
	}
	return nil
}

// IsMarked 检查是否被标记
func (dao *CheckDAOImpl) IsMarked(name, seat, start, end string) bool {
	return dao.rdb.Exists(context.Background(), "name:"+name+":"+seat+":"+start+" ~ "+end).Val() == 1
}

// GetQQ 根据name获取qq
func (dao *CheckDAOImpl) GetQQ(name string) string {
	return dao.rdb.HGet(context.Background(), "name_qq", name).Val()
}

// AddCount 添加次数
func (dao *CheckDAOImpl) AddCount(name string) {
	dao.rdb.Incr(context.Background(), "count:"+name).Val()
}

// GetCount 获取次数
func (dao *CheckDAOImpl) GetCount(name string) string {
	return dao.rdb.Get(context.Background(), "count:"+name).Val()
}

// SaveEvent 保存发送邮件事件
func (dao *CheckDAOImpl) SaveEvent(event *model.Event) error {
	result := dao.db.Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetEvent 分页获取发送邮件事件
func (dao *CheckDAOImpl) GetEvent(pn int) ([]model.Event, error) {
	var events []model.Event
	result := dao.db.Model(&model.Event{}).Order("id DESC").Limit(10).Offset((pn - 1) * 10).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetSeatRecord 获取用户设置的座位信息
func (dao *CheckDAOImpl) GetSeatRecord(date string) (*model.SeatInfo, error) {
	userId := 1
	var info model.SeatInfo
	err := dao.rdb.HGetAll(context.Background(), "seat:record:"+fmt.Sprintf("%d:", userId)+date).Scan(&info)
	if err != nil {
		return nil, err
	}
	info.Date = date
	return &info, nil
}

// SetSeatRecord  设置用户预约的座位信息
// 用户只需要设置 座位号、开始时间、结束时间
func (dao *CheckDAOImpl) SetSeatRecord(info *model.SeatInfo) error {
	userId := 1
	err := dao.rdb.HSet(context.Background(), "seat:record:"+fmt.Sprintf("%d:", userId)+info.Date,
		"seat", info.Seat,
		"start", info.Start,
		"end", info.End).Err()
	if err != nil {
		return err
	}
	return nil
}

func (dao *CheckDAOImpl) SetOccupied(date string) error {
	userId := 1
	err := dao.rdb.HSet(context.Background(), "seat:record:"+fmt.Sprintf("%d", userId)+date, "occupied", 1).Err()
	if err != nil {
		return err
	}
	return nil
}

func (dao *CheckDAOImpl) GetOccupied(date string) (bool, error) {
	userId := 1
	occupied, err := dao.rdb.HGet(context.Background(), "seat:record:"+fmt.Sprintf("%d", userId)+date, "occupied").Int()
	if err != nil {
		return false, err
	}
	return occupied == 1, nil
}

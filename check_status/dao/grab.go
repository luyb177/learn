package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type GrabDAO interface {
	FindSeatId(seat string) (string, string, error)
}
type GrabDAOImpl struct {
	rdb *redis.Client
}

func NewGrabDAOImpl(rdb *redis.Client) *GrabDAOImpl {
	return &GrabDAOImpl{
		rdb: rdb,
	}
}

// FindSeatId 通过房间号查找id // dev_id room_id
func (dao *GrabDAOImpl) FindSeatId(seat string) (string, string, error) {
	result := dao.rdb.HMGet(context.Background(), "seat:"+seat, "dev_id", "room_id").Val()
	if len(result) == 0 {
		return "", "", fmt.Errorf("seat not found")
	}
	return result[0].(string), result[1].(string), nil
}

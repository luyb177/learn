package model

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255);not null"`
	StartTime string `gorm:"type:varchar(255);not null"`
	EndTime   string `gorm:"type:varchar(255);not null"`
}

type GrabInfo struct {
	Seat      string `json:"seat"`
	DevId     string `json:"dev_id"`
	RoomId    string `json:"room_id"`
	Date      string `json:"date"`
	Start     string `json:"start"`
	End       string `json:"end"`
	FrStart   string `json:"fr_start"`
	FrEnd     string `json:"fr_end"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	TimeMs    string `json:"time_ms"`
	CheckTime string `json:"check_time"` // tomorrow today
}

type SeatInfo struct {
	Seat  string `redis:"seat"`
	Start string `redis:"start"`
	End   string `redis:"end"`
	Date  string `redis:"date"` // 注意：Date 字段不会保存在 Redis 中
}

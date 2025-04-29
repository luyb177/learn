package model

import (
	"gorm.io/gorm"
)

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

type Content struct {
	gorm.Model
	Seat    string
	Start   string
	End     string
	Status  string
	Content string
}

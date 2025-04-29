package model

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255);not null"`
	StartTime string `gorm:"type:varchar(255);not null"`
	EndTime   string `gorm:"type:varchar(255);not null"`
}

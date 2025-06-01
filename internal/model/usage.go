package model

import (
	"gorm.io/gorm"
)

type Usage struct {
	gorm.Model
	Id     uint   `gorm:"primaryKey"`
	UserId string `gorm:"not null"`
	VnetId string `gorm:""`
	Usage  int64  `gorm:"not null"`
}

func (m *Usage) TableName() string {
	return "usages"
}

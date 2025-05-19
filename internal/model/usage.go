package model

import "gorm.io/gorm"

type Usage struct {
	gorm.Model
	Id      uint   `gorm:"primaryKey"`
	UsageId string `gorm:"unique;not null"`
	UserId  string `gorm:"not null"`
	Date    string `gorm:"not null"`
	Usage   int    `gorm:"not null"`
}

func (m *Usage) TableName() string {
	return "usage"
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type Usage struct {
	Id        uint   `gorm:"primarykey"`
	UserId    string `gorm:"unique;not null"`
	Date      string `gorm:"not null"`
	Usage     int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Usage) TableName() string {
	return "usages"
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type Usages struct {
	Id         uint   `gorm:"primarykey"`
	UserId     string `gorm:"unique;not null"`
	Date       string `gorm:"not null"`
	UsageCount int    `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (u *Usages) TableName() string {
	return "usages"
}

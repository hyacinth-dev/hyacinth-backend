// 此模块定义了用户表的结构体
// 以及与数据库交互的相关方法
// 该模块使用 GORM 作为 ORM 框架

package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uint   `gorm:"primarykey"`
	UserId    string `gorm:"unique;not null"`
	Nickname  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	IsAdmin   bool   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type VNet struct {
	ID           string `gorm:"primaryKey"`
	UserId       string `gorm:"not null"`
	Name         string `gorm:"not null"`
	Enabled      bool   `gorm:"not null"`
	Token        string `gorm:"not null"`
	Password     string `gorm:"not null"`
	IpRange      string `gorm:"not null"`
	EnableDHCP   bool   `gorm:"not null"`
	ClientsLimit int    `gorm:"not null"`
	Clients      int    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (VNet) TableName() string {
	return "vnets"
}

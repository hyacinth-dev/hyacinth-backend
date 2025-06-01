package model

import "gorm.io/gorm"

type Vnet struct {
	gorm.Model
	VnetId        string `gorm:"unique;not null"`
	UserId        string `gorm:"not null"`
	Comment       string
	Enabled       bool   `gorm:"not null"`
	Token         string `gorm:"not null"`
	Password      string `gorm:"not null"`
	IpRange       string `gorm:"not null"`
	EnableDHCP    bool   `gorm:"not null"`
	ClientsLimit  int    `gorm:"not null"`
	ClientsOnline int    `gorm:"not null"`
	NeedUpdate    bool   `gorm:"not null"`
}

func (m *Vnet) TableName() string {
	return "vnets"
}

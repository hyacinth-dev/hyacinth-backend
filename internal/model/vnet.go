package model

import "gorm.io/gorm"

type Vnet struct {
	gorm.Model
}

func (m *Vnet) TableName() string {
    return "vnet"
}

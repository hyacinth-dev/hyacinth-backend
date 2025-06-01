// 此模块定义了用户表的结构体
// 以及与数据库交互的相关方法
// 该模块使用 GORM 作为 ORM 框架

package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId           string     `gorm:"unique;not null"`
	Username         string     `gorm:"not null"`
	Password         string     `gorm:"not null"`
	Email            string     `gorm:"not null"`
	UserGroup        int        `gorm:"not null;default:1"`
	PrivilegeExpiry  *time.Time `gorm:"default:null"`
	RemainingTraffic int64      `gorm:"not null;default:0"`
}

func (u *User) TableName() string {
	return "users"
}

// GetUserGroupName 获取用户组名称
func (u *User) GetUserGroupName() string {
	switch u.UserGroup {
	case 0:
		return "管理员"
	case 1:
		return "普通用户"
	case 2:
		return "青铜用户"
	case 3:
		return "白银用户"
	case 4:
		return "黄金用户"
	default:
		return "未知用户组"
	}
}

// IsPrivilegeExpired 检查特权是否过期
func (u *User) IsPrivilegeExpired() bool {
	if u.PrivilegeExpiry == nil {
		return true
	}
	return time.Now().After(*u.PrivilegeExpiry)
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.UserGroup == 0
}

// IsVip 检查用户是否为VIP用户（青铜及以上等级）
func (u *User) IsVip() bool {
	return u.UserGroup >= 2 && !u.IsPrivilegeExpired()
}

// GetRemainingTrafficMB 获取剩余流量（MB）
func (u *User) GetRemainingTrafficMB() float64 {
	return float64(u.RemainingTraffic) / (1024 * 1024)
}

// GetRemainingTrafficGB 获取剩余流量（GB）
func (u *User) GetRemainingTrafficGB() float64 {
	return float64(u.RemainingTraffic) / (1024 * 1024 * 1024)
}

// FormatRemainingTraffic 格式化剩余流量显示
func (u *User) FormatRemainingTraffic() string {
	gb := u.GetRemainingTrafficGB()
	if gb >= 1024 {
		return fmt.Sprintf("%.2f TB", gb/1024)
	}
	return fmt.Sprintf("%.2f GB", gb)
}

// GetVirtualNetworkLimit 获取用户虚拟网络数量限制
func (u *User) GetVirtualNetworkLimit() int {
	switch u.UserGroup {
	case 0: // 管理员
		return 999 // 无限制
	case 1: // 普通用户
		return 1
	case 2: // 青铜用户
		if u.IsPrivilegeExpired() {
			return 1 // 特权过期回到普通用户限制
		}
		return 3
	case 3: // 白银用户
		if u.IsPrivilegeExpired() {
			return 1 // 特权过期回到普通用户限制
		}
		return 5
	case 4: // 黄金用户
		if u.IsPrivilegeExpired() {
			return 1 // 特权过期回到普通用户限制
		}
		return 10
	default:
		return 1
	}
}

// GetMaxClientsLimitPerVNet 获取用户单个虚拟网络的最大在线人数限制
func (u *User) GetMaxClientsLimitPerVNet() int {
	switch u.UserGroup {
	case 0: // 管理员
		return 999999 // 无限制
	case 1: // 普通用户
		return 3
	case 2: // 青铜用户
		if u.IsPrivilegeExpired() {
			return 3 // 特权过期回到普通用户限制
		}
		return 5
	case 3: // 白银用户
		if u.IsPrivilegeExpired() {
			return 3 // 特权过期回到普通用户限制
		}
		return 10
	case 4: // 黄金用户
		if u.IsPrivilegeExpired() {
			return 3 // 特权过期回到普通用户限制
		}
		return 999999 // 无限制
	default:
		return 3
	}
}

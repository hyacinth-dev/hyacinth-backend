// 在该模块中定义用户相关的请求和响应结构体
// 主要包括注册、登录、获取用户信息、更新用户信息等功能
// 该模块的请求和响应结构体主要用于与前端进行数据交互

package v1

type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required" example:"1234@gmail.com"`
	Password        string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
	IsAdmin     bool   `json:"isAdmin"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Username string `json:"username" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required" example:"123456"`
	NewPassword     string `json:"newPassword" binding:"required" example:"654321"`
}

type GetProfileResponseData struct {
	UserId           string  `json:"userId"`
	Username         string  `json:"username" example:"alan"`
	Email            string  `json:"email" binding:"required,email" example:"1234@gmail.com"`
	UserGroup        int     `json:"userGroup"`
	UserGroupName    string  `json:"userGroupName"`
	PrivilegeExpiry  *string `json:"privilegeExpiry"`
	IsVip            bool    `json:"isVip"`
	ActiveTunnels    int     `json:"activeTunnels"`
	AvailableTraffic string  `json:"availableTraffic"`
	OnlineDevices    int     `json:"onlineDevices"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type GetProfileByAdminRequest struct {
	UserId string `json:"userId" binding:"required" example:"1234"`
}

// PurchasePackageRequest 购买增值服务套餐请求
type PurchasePackageRequest struct {
	PackageType int `json:"packageType" binding:"required,min=2,max=4" example:"2"` // 2=青铜 3=白银 4=黄金
	Duration    int `json:"duration" binding:"min=1,max=12" example:"1"`            // 购买时长（月数），1-12个月
}

// PurchasePackageResponse 购买增值服务套餐响应
type PurchasePackageResponse struct {
	Response
}

// GetUserGroupResponseData 获取用户组信息响应数据
type GetUserGroupResponseData struct {
	UserGroup int `json:"userGroup"`
}

// GetUserGroupResponse 获取用户组信息响应
type GetUserGroupResponse struct {
	Response
	Data GetUserGroupResponseData
}

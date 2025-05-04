// 在该模块中定义用户相关的请求和响应结构体
// 主要包括注册、登录、获取用户信息、更新用户信息等功能
// 该模块的请求和响应结构体主要用于与前端进行数据交互

package v1

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	//添加username,nickname
	Username string `json:"username" example:"alan" binding:"required"`
	Nickname string `json:"nickname" example:"alan" binding:"required"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	//改成email或username
	EmailOrUsername string `json:"emailOrUsername" binding:"required" example:"alan"`
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
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type GetUsageRequest struct {
	Range string `form:"range" binding:"required" example:"month"`
}

type UsageData struct {
	DateOrMonth string `json:"dateOrMonth" example:"2023-10"`
	Usage       int    `json:"usage"`
}

type GetUsageResponseData struct {
	Usages []UsageData `json:"usages"`
}

type GetUsageResponse struct {
	Response
	Data GetUsageResponseData
}
type VNetData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Token        string `json:"token"`
	Password     string `json:"password"`
	IpRange      string `json:"ipRange"`
	EnableDHCP   bool   `json:"enableDHCP"`
	ClientsLimit int    `json:"clientsLimit"`
	Clients      int    `json:"clients"`
}

type GetVNetResponseData struct {
	Networks []VNetData `json:"networks"`
}

type GetVNetResponse struct {
	Response
	Data GetVNetResponseData `json:"data"`
}
type UpdateVNetRequest struct {
	Name         string `json:"name" binding:"required"`
	Enabled      bool   `json:"enabled" binding:"required"`
	Token        string `json:"token" binding:"required"`
	Password     string `json:"password" binding:"required"`
	IpRange      string `json:"ipRange" binding:"required"`
	EnableDHCP   bool   `json:"enableDHCP" binding:"required"`
	ClientsLimit int    `json:"clientsLimit" binding:"required"`
}
type CreateVNetRequest struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Token        string `json:"token"`
	Password     string `json:"password"`
	IpRange      string `json:"ipRange"`
	EnableDHCP   bool   `json:"enableDHCP"`
	ClientsLimit int    `json:"clientsLimit"`
}

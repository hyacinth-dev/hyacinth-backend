package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
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
	Date  string `json:"date"`
	Usage int    `json:"usage"`
}

type GetUsageResponseData struct {
	Usages []UsageData `json:"usages"`
}

type GetUsageResponse struct {
	Response
	Data GetUsageResponseData
}

package v1

type GetUsageRequest struct {
	UserId int64  `form:"user_id" binding:"required" example:"1"`
	Range  string `form:"range" binding:"required" example:"month"`
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

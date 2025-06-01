package v1

type GetUsageRequest struct {
	UserId string `form:"userId" example:"1"`
	VnetId string `form:"vnetId" example:"vnet_123"` // 虚拟网络ID，可选，空值表示所有虚拟网络
	Range  string `form:"range" binding:"required" example:"12months"`
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

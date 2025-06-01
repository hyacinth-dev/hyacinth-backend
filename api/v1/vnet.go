package v1

type VnetProfile struct {
	VnetId       string `json:"vnetId" example:"1234"`
	Comment      string `json:"comment" example:"我的虚拟网络"`
	Enabled      bool   `json:"enabled" example:"true"`
	Token        string `json:"token" example:"1234"`
	Password     string `json:"password" example:"1234"`
	IpRange      string `json:"ipRange" example:"192.168.1.0/24"`
	EnableDHCP   bool   `json:"enableDHCP" example:"true"`
	ClientsLimit int    `json:"clientsLimit" example:"10"`
}

type GetVnetRequest struct {
	Id string `json:"Id" binding:"required" example:"1234"`
}

type GetVnetByVnetIdResponseData struct {
	UserId string `json:"userId" example:"1234"`
	VnetProfile
	ClientsOnline int `json:"clientsOnline" example:"5"`
}

type GetVnetByUserIdResponseItem struct {
	VnetId string `json:"vnetId" example:"1234"`
	VnetProfile
	ClientsOnline int `json:"clientsOnline" example:"5"`
}

type GetVnetResponseData struct {
	Vnets []GetVnetByUserIdResponseItem `json:"vnets"`
}

type GetVnetByVnetIdResponse struct {
	Response
	Data GetVnetByVnetIdResponseData
}

type GetVnetByUserIdResponse struct {
	Response
	Data GetVnetResponseData
}

type UpdateVnetRequest struct {
	VnetProfile
}

type CreateVnetRequest struct {
	VnetProfile
}

type DeleteVnetRequest struct {
	VnetID string `json:"vnetId" binding:"required" example:"1234"`
}

type EnableVnetRequest struct {
	VnetID string `json:"vnetId" binding:"required" example:"1234"`
}

type DisableVnetRequest struct {
	VnetID string `json:"vnetId" binding:"required" example:"1234"`
}

type GetVNetLimitInfoResponseData struct {
	CurrentCount           int `json:"currentCount" example:"2"`
	MaxLimit               int `json:"maxLimit" example:"5"`
	UserGroup              int `json:"userGroup" example:"3"`
	MaxClientsLimitPerVNet int `json:"maxClientsLimitPerVNet" example:"10"`
}

type GetVNetLimitInfoResponse struct {
	Response
	Data GetVNetLimitInfoResponseData `json:"data"`
}

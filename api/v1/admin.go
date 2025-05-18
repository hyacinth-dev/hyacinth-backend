package v1

type GetTotalUsageResponseData struct {
	Total int `json:"total"`
}

type GetTotalUsageResponse struct {
	Response
	Data GetTotalUsageResponseData `json:"data"`
}

type UserUsageData struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	NumNetworks int    `json:"numNetworks"`
	Usage       int    `json:"usage"`
}
type UsagePageItem struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	NumNetworks int    `json:"numNetworks"`
	Usage       int    `json:"usage"`
}

type GetUsagePageResponseData struct {
	Items     []UsagePageItem `json:"items"`
	PageCount int             `json:"pageCount"`
}

type GetUsagePageResponse struct {
	Response
	Data GetUsagePageResponseData `json:"data"`
}

type AdminVNetData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Token        string `json:"token"`
	Password     string `json:"password"`
	IpRange      string `json:"ipRange"`
	EnableDHCP   bool   `json:"enableDHCP"`
	ClientsLimit int    `json:"clientsLimit"`
	Clients      int    `json:"clients"`
	UserID       string `json:"userId"`
	UserName     string `json:"userName"`
}
type AdminGetVNetRequest struct {
	UserID string `json:"userId" binding:"required"`
}
type AdminGetVNetResponseData struct {
	Networks []AdminVNetData `json:"networks"`
}

type AdminGetVNetResponse struct {
	Response
	Data AdminGetVNetResponseData `json:"data"`
}
type AdminUpdateVNetRequest struct {
	Name         string `json:"name" binding:"required"`
	Enabled      bool   `json:"enabled" binding:"required"`
	Token        string `json:"token" binding:"required"`
	Password     string `json:"password" binding:"required"`
	IpRange      string `json:"ipRange" binding:"required"`
	EnableDHCP   bool   `json:"enableDHCP" binding:"required"`
	ClientsLimit int    `json:"clientsLimit" binding:"required"`
}
type AdminCreateVNetRequest struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Token        string `json:"token"`
	Password     string `json:"password"`
	IpRange      string `json:"ipRange"`
	EnableDHCP   bool   `json:"enableDHCP"`
	ClientsLimit int    `json:"clientsLimit"`
}

// 在此模块中定义请求的入口点
// 主要职责是接受请求、解析请求参数、验证请求合法性、调用服务层的函数、处理返回结果

package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	*Handler
	userService  service.UserService
	usageService service.UsageService
	vnetService  service.VnetService
}

func NewUserHandler(
	handler *Handler,
	userService service.UserService,
	usageService service.UsageService,
	vnetService service.VnetService,
) *UserHandler {
	return &UserHandler{
		Handler:      handler,
		userService:  userService,
		usageService: usageService,
		vnetService:  vnetService,
	}
}

// Register godoc
// @Summary 用户注册
// @Schemes
// @Description 目前只支持邮箱登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 校验用户名或邮箱格式
	usernamePattern := `^[a-zA-Z0-9_]{3,20}$`
	emailPattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	if matchedUser, _ := regexp.MatchString(usernamePattern, req.UsernameOrEmail); !matchedUser {
		if matchedEmail, _ := regexp.MatchString(emailPattern, req.UsernameOrEmail); !matchedEmail {
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
			return
		}
	}

	response, err := h.userService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, response)
}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetProfileResponse
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	// 获取用户基本信息
	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 直接调用 vnetService 获取活跃隧道数量
	activeTunnels, err := h.vnetService.GetOnlineTunnels(ctx, userId)
	if err != nil {
		h.logger.WithContext(ctx).Error("vnetService.GetOnlineTunnels error", zap.Error(err))
		// 如果获取失败，设置为0，不影响其他数据的返回
		activeTunnels = 0
	}

	// 调用 vnetService 获取在线设备数量
	onlineDevices, err := h.vnetService.GetOnlineDevicesCount(ctx, userId)
	if err != nil {
		h.logger.WithContext(ctx).Error("vnetService.GetOnlineDevicesCount error", zap.Error(err))
		// 如果获取失败，设置为0，不影响其他数据的返回
		onlineDevices = 0
	}

	// 更新活跃隧道数量和在线设备数量
	user.ActiveTunnels = activeTunnels
	user.OnlineDevices = onlineDevices

	v1.HandleSuccess(ctx, user)
}

// UpdateProfile godoc
// @Summary 修改用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req v1.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 获取当前用户信息
	currentUser, err := h.userService.GetUserByID(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 如果用户名发生了变化，检查新用户名是否已被其他用户使用
	if currentUser.Username != req.Username {
		exists, err := h.userService.CheckUsernameExists(ctx, req.Username, userId)
		if err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
			return
		}
		if exists {
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErrUsernameConflict, nil)
			return
		}

		if err := h.userService.UpdateProfile(ctx, userId, &req); err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
			return
		}
	}

	v1.HandleSuccess(ctx, nil)
}

// PurchasePackage godoc
// @Summary 购买增值服务套餐
// @Schemes
// @Description 购买增值服务套餐，传入套餐号(2=青铜,3=白银,4=黄金)和购买时长(1-12个月)
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.PurchasePackageRequest true "params"
// @Success 200 {object} v1.PurchasePackageResponse
// @Router /user/purchase [post]
func (h *UserHandler) PurchasePackage(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	var req v1.PurchasePackageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 获取当前用户信息
	user, err := h.userService.GetUserByID(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 如果是降级套餐，需要检查当前虚拟网络的最大连接数设置是否超过新套餐限制
	if user.UserGroup > req.PackageType {
		// 获取新套餐的最大在线人数限制
		var newMaxClientsLimit int
		switch req.PackageType {
		case 1: // 普通用户
			newMaxClientsLimit = 3
		case 2: // 青铜用户
			newMaxClientsLimit = 5
		case 3: // 白银用户
			newMaxClientsLimit = 10
		case 4: // 黄金用户
			newMaxClientsLimit = 999999 // 无限制
		default:
			newMaxClientsLimit = 3
		}

		// 获取用户的所有虚拟网络
		vnets, err := h.vnetService.GetVnetByUserId(ctx, userId)
		if err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
			return
		}

		// 检查是否有虚拟网络的最大连接数设置超过新套餐限制
		for _, vnet := range *vnets {
			if vnet.ClientsLimit > newMaxClientsLimit {
				v1.HandleError(ctx, http.StatusBadRequest, v1.ErrVnetClientsLimitExceeded, nil)
				return
			}
		}
	}

	if err := h.userService.PurchasePackage(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *UserHandler) GetUsage(ctx *gin.Context) {
	var req v1.GetUsageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	req.UserId = userId

	usage, err := h.usageService.GetUsage(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, usage)
}

// GetVNetList godoc
// @Summary 获取用户的虚拟网络列表
// @Schemes
// @Description 获取当前用户的所有虚拟网络
// @Tags 虚拟网络模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetVnetByUserIdResponse
// @Router /vnet [get]
func (h *UserHandler) GetVNetList(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	vnets, err := h.vnetService.GetVnetByUserId(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 转换为API响应格式
	var responseItems []v1.GetVnetByUserIdResponseItem
	if vnets != nil {
		for _, vnet := range *vnets {
			responseItems = append(responseItems, v1.GetVnetByUserIdResponseItem{
				VnetId: vnet.VnetId,
				VnetProfile: v1.VnetProfile{
					VnetId:       vnet.VnetId,
					Comment:      vnet.Comment,
					Enabled:      vnet.Enabled,
					Token:        vnet.Token,
					Password:     vnet.Password,
					IpRange:      vnet.IpRange,
					EnableDHCP:   vnet.EnableDHCP,
					ClientsLimit: vnet.ClientsLimit,
				},
				ClientsOnline: vnet.ClientsOnline,
			})
		}
	}

	response := v1.GetVnetResponseData{
		Vnets: responseItems,
	}

	v1.HandleSuccess(ctx, response)
}

// CreateVNet godoc
// @Summary 创建虚拟网络
// @Schemes
// @Description 为当前用户创建新的虚拟网络
// @Tags 虚拟网络模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateVnetRequest true "params"
// @Success 200 {object} v1.Response
// @Router /vnet [post]
func (h *UserHandler) CreateVNet(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	var req v1.CreateVnetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		println("cannot bind json:", err.Error())
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 生成唯一的VnetId
	req.VnetId = generateVnetId(userId)

	// 检查vnet名称是否已存在
	exists, err := h.vnetService.CheckVnetTokenExists(ctx, req.Token, "")
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	if exists {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrVnetTokenAlreadyUse, nil)
		return
	}

	// 获取用户信息用于验证
	user, err := h.userService.GetUserByID(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 检查客户端数量限制
	maxClientsLimit := user.GetMaxClientsLimitPerVNet()
	if req.ClientsLimit > maxClientsLimit {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrVnetClientsLimitExceeded, nil)
		return
	}

	// 检查虚拟网络数量限制
	if req.Enabled {
		// 获取当前运行中的虚拟网络数量
		currentRunningCount, err := h.vnetService.GetRunningVnetCount(ctx, userId)
		if err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
			return
		}

		// 检查是否超过限制
		if currentRunningCount >= user.GetVirtualNetworkLimit() {
			// 如果超过限制，返回错误给前端
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErrVnetLimitExceeded, nil)
			return
		}
	}

	if err := h.vnetService.CreateVnet(ctx, &req, userId); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// UpdateVNet godoc
// @Summary 更新虚拟网络
// @Schemes
// @Description 更新用户的虚拟网络配置
// @Tags 虚拟网络模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param vnetId path string true "VNet ID"
// @Param request body v1.UpdateVnetRequest true "params"
// @Success 200 {object} v1.Response
// @Router /vnet/{vnetId} [put]
func (h *UserHandler) UpdateVNet(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	vnetId := ctx.Param("vnetId")
	if vnetId == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var req v1.UpdateVnetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 设置VnetId
	req.VnetId = vnetId

	// fmt.Println("UpdateVNet request:", req)

	// 验证VNet是否属于当前用户
	vnet, err := h.vnetService.GetVnetByVnetId(ctx, vnetId)
	if err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}
	if vnet.UserId != userId {
		v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
		return
	}

	// 获取用户信息用于验证
	user, err := h.userService.GetUserByID(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 检查客户端数量限制
	maxClientsLimit := user.GetMaxClientsLimitPerVNet()
	if maxClientsLimit != 999999 && req.ClientsLimit > maxClientsLimit {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrVnetClientsLimitExceeded, nil)
		return
	}

	// 检查虚拟网络数量限制（如果要启用网络）
	if req.Enabled && !vnet.Enabled {
		// 获取当前运行中的虚拟网络数量
		currentRunningCount, err := h.vnetService.GetRunningVnetCount(ctx, userId)
		if err != nil {
			v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
			return
		}

		// 检查是否超过限制
		if currentRunningCount >= user.GetVirtualNetworkLimit() {
			// 如果超过限制，返回错误给前端
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErrVnetLimitExceeded, nil)
			return
		}
	}

	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(req.Token, vnetId)
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	exists, err := h.vnetService.CheckVnetTokenExists(ctx, req.Token, vnetId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	if exists {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrVnetTokenAlreadyUse, nil)
		return
	}

	if err := h.vnetService.UpdateVnet(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// DeleteVNet godoc
// @Summary 删除虚拟网络
// @Schemes
// @Description 删除用户的虚拟网络
// @Tags 虚拟网络模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param vnetId path string true "VNet ID"
// @Success 200 {object} v1.Response
// @Router /vnet/{vnetId} [delete]
func (h *UserHandler) DeleteVNet(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	vnetId := ctx.Param("vnetId")
	if vnetId == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 验证VNet是否属于当前用户
	vnet, err := h.vnetService.GetVnetByVnetId(ctx, vnetId)
	if err != nil {
		v1.HandleError(ctx, http.StatusNotFound, v1.ErrNotFound, nil)
		return
	}
	if vnet.UserId != userId {
		v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
		return
	}

	req := v1.DeleteVnetRequest{
		VnetID: vnetId,
	}

	if err := h.vnetService.DeleteVnet(ctx, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// GetUserGroup godoc
// @Summary 获取用户组信息
// @Schemes
// @Description 获取当前用户的组信息，用于商城套餐显示
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetUserGroupResponse
// @Router /user/group [get]
func (h *UserHandler) GetUserGroup(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	// 获取用户基本信息
	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// 构造响应数据
	response := v1.GetUserGroupResponseData{
		UserGroup: user.UserGroup,
	}

	v1.HandleSuccess(ctx, response)
}

// GetVNetLimitInfo godoc
// @Summary 获取用户的虚拟网络限制信息
// @Schemes
// @Description 获取当前用户的虚拟网络限制和使用情况
// @Tags 虚拟网络模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetVNetLimitInfoResponse
// @Router /vnet/limit [get]
func (h *UserHandler) GetVNetLimitInfo(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 获取当前运行中的虚拟网络数量
	currentRunningCount, err := h.vnetService.GetRunningVnetCount(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	// 构造响应数据
	response := v1.GetVNetLimitInfoResponseData{
		CurrentCount:           currentRunningCount,
		MaxLimit:               user.GetVirtualNetworkLimit(),
		UserGroup:              user.UserGroup,
		MaxClientsLimitPerVNet: user.GetMaxClientsLimitPerVNet(),
	}

	v1.HandleSuccess(ctx, response)
}

// 生成VnetId的辅助函数
func generateVnetId(userId string) string {
	// 这里可以使用更复杂的ID生成逻辑
	// 简单起见，使用userId + 时间戳
	return userId + "_vnet_" + fmt.Sprintf("%d", time.Now().Unix())
}

// ChangePassword godoc
// @Summary 修改密码
// @Schemes
// @Description 修改用户密码，需要验证当前密码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ChangePasswordRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user/password [put]
func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req v1.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.ChangePassword(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

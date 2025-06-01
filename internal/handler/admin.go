package handler

import (
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	*Handler
	adminService service.AdminService
	userService  service.UserService
	usageService service.UsageService
	vnetService  service.VnetService
}

func NewAdminHandler(
	handler *Handler,
	adminService service.AdminService,
	userService service.UserService,
	usageService service.UsageService,
	vnetService service.VnetService,
) *AdminHandler {
	return &AdminHandler{
		Handler:      handler,
		adminService: adminService,
		userService:  userService,
		usageService: usageService,
		vnetService:  vnetService,
	}
}

func (h *AdminHandler) GetUsage(ctx *gin.Context) {
	var req v1.GetUsageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	usage, err := h.usageService.GetUsage(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, usage)
}

func (h *AdminHandler) GetUserProfile(ctx *gin.Context) {
	var req v1.GetProfileByAdminRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	user, err := h.userService.GetProfile(ctx, req.UserId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, user)
}

func (h *AdminHandler) GetVnetByUserId(ctx *gin.Context) {
	var req v1.GetVnetRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	vnet, err := h.vnetService.GetVnetByUserId(ctx, req.Id)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, vnet)
}

func (h *AdminHandler) GetVnetByVnetId(ctx *gin.Context) {
	var req v1.GetVnetRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	vnet, err := h.vnetService.GetVnetByVnetId(ctx, req.Id)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, vnet)
}

func (h *AdminHandler) UpdateVnet(ctx *gin.Context) {
	var req v1.UpdateVnetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	err := h.vnetService.UpdateVnet(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *AdminHandler) DeleteVnet(ctx *gin.Context) {
	var req v1.DeleteVnetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	err := h.vnetService.DeleteVnet(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetOnlineDevicesStats godoc
// @Summary 获取在线设备统计
// @Schemes
// @Description 获取指定用户或所有用户的在线设备数量统计
// @Tags 管理员模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param userId query string false "用户ID，为空则获取所有用户的统计"
// @Success 200 {object} v1.GetOnlineDevicesStatsResponse
// @Router /admin/online-devices-stats [get]
func (h *AdminHandler) GetOnlineDevicesStats(ctx *gin.Context) {
	userId := ctx.Query("userId")
	// 如果userId为空，将使用空字符串，这在repository层会被处理为查询所有用户

	onlineDevicesCount, err := h.vnetService.GetOnlineDevicesCount(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	response := map[string]interface{}{
		"userId":             userId,
		"onlineDevicesCount": onlineDevicesCount,
	}

	v1.HandleSuccess(ctx, response)
}

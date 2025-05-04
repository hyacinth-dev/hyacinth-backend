package handler

import (
	"github.com/gin-gonic/gin"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/service"
	"net/http"
)

type AdminHandler struct {
	*Handler
	adminService service.AdminService
}

func NewAdminHandler(handler *Handler, adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		Handler:      handler,
		adminService: adminService,
	}
}

func (h *AdminHandler) GetTotalUsage(ctx *gin.Context) {
	data, err := h.adminService.GetTotalUsage(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, data)
}

func (h *AdminHandler) GetUsagePage(ctx *gin.Context) {
	page := ctx.Param("page")
	const pageSize = 10 // 默认每页条目数量

	data, err := h.adminService.GetUsagePage(ctx, page, pageSize)

	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, data)
}

func (h *AdminHandler) GetUsage(ctx *gin.Context) {

	userId := ctx.Param("id")
	var req v1.GetUsageRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	usage, err := h.adminService.GetUsage(ctx, userId, &req)

	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, usage)
}

func (h *AdminHandler) AdminGetVNet(ctx *gin.Context) {
	var req v1.AdminGetVNetRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	vnet, err := h.adminService.AdminGetVNet(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, vnet)
}

func (h *AdminHandler) AdminCreateVNet(ctx *gin.Context) {
	userId := ctx.Param("USERID")
	var req v1.AdminCreateVNetRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.adminService.AdminCreateVNet(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *AdminHandler) AdminUpdateVNet(ctx *gin.Context) {
	vnetId := ctx.Param("VNETID")
	var req v1.AdminUpdateVNetRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.adminService.AdminUpdateVNet(ctx, vnetId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *AdminHandler) AdminDeleteVNet(ctx *gin.Context) {
	vnetId := ctx.Param("VNETID")

	if err := h.adminService.AdminDeleteVNet(ctx, vnetId); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

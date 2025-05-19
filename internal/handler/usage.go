package handler

import (
	"github.com/gin-gonic/gin"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/service"
	"net/http"
)

type UsageHandler struct {
	*Handler
	usageService service.UsageService
}

func NewUsageHandler(
	handler *Handler,
	usageService service.UsageService,
) *UsageHandler {
	return &UsageHandler{
		Handler:      handler,
		usageService: usageService,
	}
}

// UpdateProfile godoc
// @Summary 获取用量
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.GetUsageRequest true "params"
// @Success 200 {object} v1.Response
// @Router /usage [get]
func (h *UsageHandler) GetUsage(ctx *gin.Context) {
	var req v1.GetUsageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	usage, err := h.usageService.GetAllUsagesByUserId(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, usage)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"hyacinth-backend/internal/service"
)

type VnetHandler struct {
	*Handler
	vnetService service.VnetService
}

func NewVnetHandler(
    handler *Handler,
    vnetService service.VnetService,
) *VnetHandler {
	return &VnetHandler{
		Handler:      handler,
		vnetService: vnetService,
	}
}

func (h *VnetHandler) GetVnet(ctx *gin.Context) {

}

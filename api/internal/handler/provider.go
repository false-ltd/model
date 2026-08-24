package handler

import (
	"net/http"

	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/service"
	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	providerService *service.ProviderService
}

func NewProviderHandler(providerService *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{providerService: providerService}
}

// List godoc
// @Summary 获取服务商列表
// @Description 返回所有 AI 服务商及其模型数量
// @Tags providers
// @Produce json
// @Success 200 {object} model.PagedResponse
// @Failure 500 {object} model.Response
// @Router /api/v1/providers [get]
func (h *ProviderHandler) List(c *gin.Context) {
	providers, count, err := h.providerService.List()
	if err != nil {
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.CountSuccessResponse(providers, count))
}

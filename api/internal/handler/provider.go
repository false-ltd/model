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
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} model.Response
// @Router /api/v1/providers [get]
func (h *ProviderHandler) List(c *gin.Context) {
	providers, count, err := h.providerService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(50001, err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    providers,
		"meta":    gin.H{"count": count},
	})
}

package handler

import (
	"net/http"

	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/service"
	"github.com/gin-gonic/gin"
)

type AtlasHandler struct {
	atlasService *service.AtlasService
}

func NewAtlasHandler(atlasService *service.AtlasService) *AtlasHandler {
	return &AtlasHandler{atlasService: atlasService}
}

// Get godoc
// @Summary 获取模型星图点云
// @Description 返回全部模型的紧凑点位数据（价格、上下文、能力数、提供商），用于交互式星图
// @Tags atlas
// @Produce json
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/atlas [get]
func (h *AtlasHandler) Get(c *gin.Context) {
	data, err := h.atlasService.Get()
	if err != nil {
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(data))
}

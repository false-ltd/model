package handler

import (
	"net/http"
	"strconv"

	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/service"
	"github.com/gin-gonic/gin"
)

type CompareHandler struct {
	modelService *service.ModelService
}

func NewCompareHandler(modelService *service.ModelService) *CompareHandler {
	return &CompareHandler{modelService: modelService}
}

// Compare godoc
// @Summary 模型对比
// @Description 根据多个模型 ID 返回对比数据，保持请求顺序并去重
// @Tags compare
// @Produce json
// @Param ids query string true "模型 ID，逗号分隔 (如 1,2,3)"
// @Success 200 {object} model.PagedResponse
// @Failure 500 {object} model.Response
// @Router /api/v1/compare [get]
func (h *CompareHandler) Compare(c *gin.Context) {
	idsStr := c.Query("ids")

	var ids []uint
	for _, s := range parseCommaSlice(idsStr) {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, uint(n))
	}

	result, err := h.modelService.Compare(ids)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.CountSuccessResponse(result.Data, result.Count))
}

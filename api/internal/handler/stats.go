package handler

import (
	"net/http"

	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/service"
	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// Get godoc
// @Summary 获取统计数据
// @Description 返回 KPI、分层分布、能力统计、模态分布等综合统计
// @Tags stats
// @Produce json
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/stats [get]
func (h *StatsHandler) Get(c *gin.Context) {
	stats, err := h.statsService.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(50001, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(stats))
}

package handler

import (
	"net/http"

	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/service"
	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	syncService *service.SyncService
}

func NewSyncHandler(syncService *service.SyncService) *SyncHandler {
	return &SyncHandler{syncService: syncService}
}

// Trigger godoc
// @Summary 触发数据同步
// @Description 从 models.dev 拉取最新数据，10 分钟内不会重复同步
// @Tags sync
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/sync [post]
func (h *SyncHandler) Trigger(c *gin.Context) {
	result, err := h.syncService.Trigger()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(50001, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(result))
}

// Status godoc
// @Summary 查看同步状态
// @Description 返回最近一次同步时间
// @Tags sync
// @Produce json
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/sync/status [get]
func (h *SyncHandler) Status(c *gin.Context) {
	status, err := h.syncService.GetStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(50001, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(status))
}

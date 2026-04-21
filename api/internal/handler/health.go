package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health godoc
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

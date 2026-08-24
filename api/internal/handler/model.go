package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/service"
	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	modelService *service.ModelService
}

func NewModelHandler(modelService *service.ModelService) *ModelHandler {
	return &ModelHandler{modelService: modelService}
}

// List godoc
// @Summary 获取模型列表
// @Description 分页查询 AI 模型，支持筛选、排序
// @Tags models
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(100)
// @Param q query string false "名称搜索"
// @Param sort query string false "排序字段 (name/cost_input/cost_output/limit_context/release_date)" default(name)
// @Param order query string false "排序方向 (asc/desc)" default(asc)
// @Param providers query string false "服务商过滤，逗号分隔"
// @Param input_types query string false "输入模态过滤，逗号分隔"
// @Param output_types query string false "输出模态过滤，逗号分隔"
// @Param reasoning query bool false "推理模型筛选 (true/false)"
// @Param tool_call query bool false "工具调用筛选 (true/false)"
// @Param vision query bool false "视觉模型筛选 (true/false)"
// @Param attachment query bool false "附件支持筛选 (true/false)"
// @Param open_weights query bool false "开放权重筛选 (true/false)"
// @Param structured_output query bool false "结构化输出筛选 (true/false)"
// @Param temperature query bool false "温度参数筛选 (true/false)"
// @Param free_only query bool false "仅免费模型"
// @Param under_1 query bool false "输入成本低于 $1"
// @Param price_min query number false "最低输入成本"
// @Param price_max query number false "最高输入成本"
// @Param price_output_min query number false "最低输出成本"
// @Param price_output_max query number false "最高输出成本"
// @Success 200 {object} model.PagedResponse
// @Failure 500 {object} model.Response
// @Router /api/v1/models [get]
func (h *ModelHandler) List(c *gin.Context) {
	f := parseModelFilter(c)
	result, err := h.modelService.List(f)
	if err != nil {
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.PagedSuccessResponse(result.Data, result.Total, f.Page, f.PageSize))
}

// Get godoc
// @Summary 获取模型详情
// @Description 根据 ID 获取单个模型详情（含服务商信息）
// @Tags models
// @Produce json
// @Param id path int true "模型 ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /api/v1/models/{id} [get]
func (h *ModelHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(40001, "invalid model ID"))
		return
	}
	m, err := h.modelService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse(40401, "model not found"))
			return
		}
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(m))
}

func parseModelFilter(c *gin.Context) *model.ModelFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	if pageSize < 1 {
		pageSize = 100
	}
	if pageSize > 200 {
		pageSize = 200
	}

	f := &model.ModelFilter{
		Page:        page,
		PageSize:    pageSize,
		Query:       c.Query("q"),
		Sort:        c.DefaultQuery("sort", "name"),
		Order:       c.DefaultQuery("order", "asc"),
		Providers:   parseCommaSlice(c.Query("providers")),
		InputTypes:  parseCommaSlice(c.Query("input_types")),
		OutputTypes: parseCommaSlice(c.Query("output_types")),

		Reasoning:        parseBoolFilter(c, "reasoning"),
		ToolCall:         parseBoolFilter(c, "tool_call"),
		Vision:           parseBoolFilter(c, "vision"),
		Attachment:       parseBoolFilter(c, "attachment"),
		OpenWeights:      parseBoolFilter(c, "open_weights"),
		StructuredOutput: parseBoolFilter(c, "structured_output"),
		Temperature:      parseBoolFilter(c, "temperature"),
		FreeOnly:         parseBoolFilter(c, "free_only"),
		Under1:           parseBoolFilter(c, "under_1"),
	}

	if v := c.Query("price_min"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.PriceMin = &f64
		}
	}
	if v := c.Query("price_max"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.PriceMax = &f64
		}
	}
	if v := c.Query("price_output_min"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.PriceOutputMin = &f64
		}
	}
	if v := c.Query("price_output_max"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.PriceOutputMax = &f64
		}
	}

	return f
}

// parseBoolFilter maps an explicit true/false query value to a *bool filter
// and returns nil when the parameter is absent or unrecognized.
func parseBoolFilter(c *gin.Context, key string) *bool {
	v := c.Query(key)
	switch v {
	case "true":
		b := true
		return &b
	case "false":
		b := false
		return &b
	}
	return nil
}

func parseCommaSlice(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			result = append(result, p)
		}
	}
	return result
}

// internalError logs the cause server-side and responds with a generic
// message so DSN/SQL internals never leak to clients.
func internalError(c *gin.Context, err error) {
	reqID, _ := c.Get("request_id")
	log.Printf("[%v] internal error on %s %s: %v", reqID, c.Request.Method, c.Request.URL.Path, err)
	c.JSON(http.StatusInternalServerError, model.ErrorResponse(50001, "internal server error"))
}

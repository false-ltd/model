package handler

import (
	"net/http"
	"strconv"

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
// @Param reasoning query bool false "仅推理模型"
// @Param tool_call query bool false "仅工具调用模型"
// @Param vision query bool false "仅视觉模型"
// @Param attachment query bool false "仅支持附件的模型"
// @Param free_only query bool false "仅免费模型"
// @Param under_1 query bool false "输入成本低于 $1"
// @Param price_min query number false "最低输入成本"
// @Param price_max query number false "最高输入成本"
// @Success 200 {object} model.PagedResponse
// @Failure 500 {object} model.Response
// @Router /api/v1/models [get]
func (h *ModelHandler) List(c *gin.Context) {
	f := parseModelFilter(c)
	result, err := h.modelService.List(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(50001, err.Error()))
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
		c.JSON(http.StatusNotFound, model.ErrorResponse(40401, "model not found"))
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
	}

	if v := c.Query("reasoning"); v == "true" {
		f.Reasoning = boolPtr(true)
	}
	if v := c.Query("tool_call"); v == "true" {
		f.ToolCall = boolPtr(true)
	}
	if v := c.Query("vision"); v == "true" {
		f.Vision = boolPtr(true)
	}
	if v := c.Query("attachment"); v == "true" {
		f.Attachment = boolPtr(true)
	}
	if v := c.Query("open_weights"); v == "true" {
		f.OpenWeights = boolPtr(true)
	}
	if v := c.Query("structured_output"); v == "true" {
		f.StructuredOutput = boolPtr(true)
	}
	if v := c.Query("temperature"); v == "true" {
		f.Temperature = boolPtr(true)
	}
	if v := c.Query("free_only"); v == "true" {
		f.FreeOnly = boolPtr(true)
	}
	if v := c.Query("under_1"); v == "true" {
		f.Under1 = boolPtr(true)
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

func parseCommaSlice(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			if part != "" {
				result = append(result, part)
			}
			start = i + 1
		}
	}
	return result
}

func boolPtr(b bool) *bool { return &b }

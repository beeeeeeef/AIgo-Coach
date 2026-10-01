package analyze

import (
	"github.com/gin-gonic/gin"
)

// Handler 处理分析相关的 HTTP 请求
type Handler struct {
	service *AnalyzeService // 业务逻辑服务
}

// NewHandler 创建 Handler 实例
func NewHandler(service *AnalyzeService) *Handler {
	return &Handler{
		service: service,
	}
}

// HandleAnalyze 处理 POST /api/analyze 请求
// 这是插件调用的核心接口
func (h *Handler) HandleAnalyze(c *gin.Context) {
	// 1. 解析请求体
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, AnalyzeResponse{
			Result: false,
			Error:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 2. 调用业务层分析代码
	resp, err := h.service.AnalyzeCode(c.Request.Context(), req)
	if err != nil {
		// 业务层返回错误（如 AI 调用失败）
		c.JSON(500, AnalyzeResponse{
			Result: false,
			Error:  err.Error(),
		})
		return
	}

	// 3. 返回成功响应
	c.JSON(200, resp)
}

// HandleAnalyzeMock 保留的 mock 接口（用于测试或降级）
func HandleAnalyzeMock(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, AnalyzeResponse{
			Result: false,
			Error:  err.Error(),
		})
		return
	}

	// 返回固定 mock 数据
	c.JSON(200, AnalyzeResponse{
		Result: true,
		Data: data{
			Summary:    "你现在的思路接近暴力枚举，可以解决小数据，但复杂度较高。",
			Problem:    []string{"当前做法需要两层循环，时间复杂度是 O(n^2)"},
			Hints:      []string{"思考：遍历到 nums[i] 时，你真正需要找的另一个数是什么？"},
			Complexity: "目标方向通常可以做到 O(n) 时间复杂度。",
		},
	})
}

package analyze

import (
	"github.com/gin-gonic/gin"
)

// 假接口
func HandleAnalyzeMock(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, AnalyzeResponse{
			Result: false,
			Error:  err.Error(),
		})
		return
	}
	resq := AnalyzeCode(req)
	c.JSON(200, resq)

}

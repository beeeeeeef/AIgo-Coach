package main

import (
	"aigo-coach/backend/internal/llm" // 导入你写的 Gemini 模块
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// 2. 配置跨域 (CORS) - 允许前端访问
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 测试接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// /chat 路由 ---
	r.POST("/chat", func(c *gin.Context) {
		// 定义请求格式
		type RequestBody struct {
			Code     string `json:"code"`
			Question string `json:"question"`
		}
		var req RequestBody

		// 解析 JSON
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
			return
		}

		// 🔥 关键点：调用 ChatWithGemini
		reply, err := llm.ChatWithGemini(req.Code, req.Question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 返回结果
		c.JSON(http.StatusOK, gin.H{"reply": reply})
	})

	// 启动服务器
	r.Run(":8080")
}

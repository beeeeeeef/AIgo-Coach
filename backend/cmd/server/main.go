package main

import (
	"log"
	"os"

	"aigo-coach/backend/internal/ai"
	"aigo-coach/backend/internal/analyze"
	"aigo-coach/backend/internal/chat"
	"aigo-coach/backend/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到 .env 文件，将从系统环境变量读取配置")
	}

	// 2. DeepSeek API Key
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 未配置 DEEPSEEK_API_KEY 环境变量")
	}

	// 3. PostgreSQL
	databaseURL := os.Getenv("DATABASE_URL")
	sqlDB, err := db.OpenPostgres(databaseURL)
	if err != nil {
		log.Fatal("错误: ", err)
	}
	defer sqlDB.Close()
	log.Println("PostgreSQL 连接成功")

	// 4. AI 客户端
	aiClient := ai.NewDeepSeekClient(apiKey)
	log.Println("DeepSeek 客户端初始化成功")

	// 5. analyze（保留调试）
	analyzeService := analyze.NewAnalyzeService(aiClient)
	analyzeHandler := analyze.NewHandler(analyzeService)

	// 6. chat（插件主路径）
	chatRepo := chat.NewRepository(sqlDB)
	chatService := chat.NewService(chatRepo, aiClient)
	chatHandler := chat.NewHandler(chatService)

	// 7. 路由
	r := gin.Default()
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

	r.POST("/api/chat", chatHandler.HandleChat)
	r.POST("/api/analyze", analyzeHandler.HandleAnalyze)
	r.POST("/api/analyze/mock", analyze.HandleAnalyzeMock)

	log.Println("服务启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}

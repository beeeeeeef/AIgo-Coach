package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//创建一个默认的路由引擎
	r := gin.Default()
	//配置cors（跨域允许）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	//定义一个简单的测试接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "Backend is running!",
		})
	})
	//启动HTTP服务，监听在8080端口
	r.Run(":8080")
}

package com

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitMiddleware() error {
	// 发布模式
	if Config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	Router = gin.New()

	// 日志中间件
	Router.Use(gin.Logger())

	// 错误处理中间件
	Router.Use(gin.Recovery())

	// 鉴权中间件
	Router.Use(authMiddleware())

	// 跨域中间件
	Router.Use(accessControlAllow())

	log.Println("[INFO] 初始化中间件成功")

	return nil
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头鉴权信息
		auth_header := c.GetHeader("Authorization")

		// 匿名请求
		if auth_header == "" {
			c.Next()
			return
		}

		// 鉴权
		token_parts := strings.Split(auth_header, " ")
		if len(token_parts) != 2 || token_parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "鉴权失败",
				"error":   "非法请求令牌",
			})
			return
		}

		// 设置请求令牌
		c.Set("signed_token", token_parts[1])

		c.Next()
	}
}

func accessControlAllow() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")   // 添加 Authorization
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type") // 可选：允许前端访问的响应头
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")                      // 可选：如果需要跨域携带 Cookie

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 204
			return
		}

		c.Next()
	}
}

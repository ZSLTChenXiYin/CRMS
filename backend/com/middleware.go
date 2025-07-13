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
		}

		// 鉴权
		token_parts := strings.Split(auth_header, " ")
		if len(token_parts) != 2 || token_parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "非法请求令牌",
			})
			return
		}

		// 设置请求令牌
		c.Set("signed_token", token_parts[1])

		c.Next()
	}
}

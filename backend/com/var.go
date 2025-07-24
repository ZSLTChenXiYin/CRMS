package com

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config   = &config{} // 配置文件
	Router   *gin.Engine // 路由
	Database *gorm.DB    // 数据库
	Redis    *redis.Client
	Context  = context.Background()
)

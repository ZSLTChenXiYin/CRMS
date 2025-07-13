package com

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Config   = &config{} // 配置文件
	Router   *gin.Engine // 路由
	Database *gorm.DB    // 数据库
)

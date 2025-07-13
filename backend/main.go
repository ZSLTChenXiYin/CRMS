package main

import (
	"fmt"
	"log"

	"github.com/ZSLTChenXiYin/CRMS/crms-backend/api"
	"github.com/ZSLTChenXiYin/CRMS/crms-backend/api/service"
	"github.com/ZSLTChenXiYin/CRMS/crms-backend/com"
)

func init() {
	// 初始化配置
	err := com.InitConfig()
	if err != nil {
		panic(err)
	}

	// 初始化数据库
	err = com.InitDatabase()
	if err != nil {
		panic(err)
	}

	// 初始化中间件
	err = com.InitMiddleware()
	if err != nil {
		panic(err)
	}

	// 初始化路由
	err = api.InitUser()
	if err != nil {
		panic(err)
	}

	err = api.InitAsset()
	if err != nil {
		panic(err)
	}
}

func main() {
	// 自动迁移
	if com.Config.AutoMigrate {
		err := com.Database.AutoMigrate(&service.User{}, &service.Asset{}, &service.UserAssetMapping{}, &service.OperationLog{})
		if err != nil {
			panic(err)
		}

		log.Println("[INFO] 自动迁移成功")
	}

	log.Println("[INFO] 服务启动端口：", com.Config.Port)

	// 启动服务
	err := com.Router.Run(fmt.Sprintf(":%d", com.Config.Port))
	if err != nil {
		panic(err)
	}
}

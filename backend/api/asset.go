package api

import (
	"log"

	"github.com/ZSLTChenXiYin/CRMS/crms-backend/api/controller"
	"github.com/ZSLTChenXiYin/CRMS/crms-backend/com"
)

func InitAsset() error {
	asset_group := com.Router.Group("/asset")
	{
		// 添加个人资产
		asset_group.POST("/add", controller.PostAddAsset)

		// 分享个人资产权限
		asset_group.POST("/share", controller.PostShareAsset)
		// 查询个人分享的资产权限
		asset_group.GET("/share", controller.GetShareAsset)
		// 修改个人分享的资产权限
		asset_group.PUT("/share", controller.PutShareAsset)
		// 撤销个人分享的资产权限
		asset_group.DELETE("/share", controller.DeleteShareAsset)
		// 获取资产列表（个人资产和有权限的非个人资产）
		asset_group.GET("/list", controller.GetAssetList)
		// 查询指定资产信息（个人资产和有权限的非个人资产）
		asset_group.GET("/info", controller.GetAssetInfo)
		// 修改指定资产信息（个人资产和有权限的非个人资产）
		asset_group.PUT("/update", controller.PutAssetInfo)

		// 删除个人资产
		asset_group.DELETE("/delete", controller.DeleteAsset)
	}

	log.Println("[INFO] Asset API初始化成功")

	return nil
}

package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ZSLTChenXiYin/CRMS/crms-backend/api/service"
	"github.com/ZSLTChenXiYin/CRMS/crms-backend/com"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 添加个人资产
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	type: string
//	name: string
//	data: string (JSON格式)
func PostAddAsset(c *gin.Context) {
	// 获取用户信息
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "添加个人资产失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理用户请求
	resource_type := c.PostForm("type")
	name := c.PostForm("name")
	data := datatypes.JSON(c.PostForm("data"))

	switch resource_type {
	case "server":
		// 验证数据格式
		server := service.Server{}
		err := json.Unmarshal(data, &server)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "添加个人资产失败",
				"error":   err.Error(),
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的资源类型",
		})
		return
	}

	// 创建资产
	asset := &service.Asset{
		Type:    resource_type,
		Name:    name,
		Data:    data,
		OwnerID: user_token.UserID,
	}

	asset_service := service.NewAssetService(com.Database)

	err = asset_service.CreateAsset(asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "添加个人资产失败",
			"error":   err.Error(),
		})
		return
	}

	// 创建操作日志
	additional_information := map[string]any{}

	service.RecordUserBasicInformation(c, additional_information)

	// 异步记录操作日志
	go func() {
		operation_log := &service.OperationLog{
			UserID:        user_token.UserID,
			ResourceType:  "asset",
			ResourceID:    asset.ID,
			Action:        "create",
			ActionDetails: "添加个人资产",
			CreatedAt:     time.Now(),
		}

		additional_infomation_json, err := json.Marshal(additional_information)
		if err != nil {
			log.Println("操作日志记录错误:", err)
		} else {
			operation_log.AdditionalInformation = additional_infomation_json
		}

		err = service.NewOperationLogService(com.Database).CreateOperationLog(operation_log)
		if err != nil {
			log.Println("操作日志记录错误:", err)

			json_data, err := json.Marshal(operation_log)
			if err != nil {
				log.Println("操作日志记录错误:", err)
			}

			log.Println("操作日志记录:", string(json_data))
			return
		}
	}()

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "添加个人资产成功",
	})
}

// 删除个人资产
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	asset_id: uint (资产ID)
func DeleteAsset(c *gin.Context) {
	// 获取用户令牌
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "删除个人资产失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理用户请求
	asset_id, err := strconv.ParseUint(c.PostForm("asset_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除个人资产失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取资产
	asset, err := service.NewAssetService(com.Database).GetAssetByID(uint(asset_id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "资产不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "删除个人资产失败",
				"error":   err.Error(),
			})
		}
		return
	}
	// 检查资产是否属于当前用户
	if asset.OwnerID != user_token.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "无权限删除",
		})
	}

	// 删除资产
	err = service.NewAssetService(com.Database).DeleteAssetByID(uint(asset_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除个人资产失败",
			"error":   err.Error(),
		})
		return
	}

	// 删除用户资产映射
	err = service.NewUserAssetService(com.Database).DeleteUserAssetMappingByAssetId(uint(asset_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除个人资产失败",
			"error":   err.Error(),
		})
		return
	}

	// 创建操作日志
	additional_information := map[string]any{}
	service.RecordUserBasicInformation(c, additional_information)

	// 异步记录操作日志
	go func() {
		operation_log := &service.OperationLog{
			UserID:        user_token.UserID,
			ResourceType:  "asset",
			ResourceID:    asset.ID,
			Action:        "delete",
			ActionDetails: "删除个人资产",
			CreatedAt:     time.Now(),
		}

		additional_infomation_json, err := json.Marshal(additional_information)
		if err != nil {
			log.Println("操作日志记录错误:", err)
		} else {
			operation_log.AdditionalInformation = additional_infomation_json
		}

		err = service.NewOperationLogService(com.Database).CreateOperationLog(operation_log)
		if err != nil {
			log.Println("操作日志记录错误:", err)

			json_data, err := json.Marshal(operation_log)
			if err != nil {
				log.Println("操作日志记录错误:", err)
			}

			log.Println("操作日志记录:", string(json_data))
			return
		}
	}()

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "删除个人资产成功",
	})
}

// 获取资产列表（个人资产和有权限的非个人资产）
//
// 请求头参数：
//
//	Authorization: Bearer <token>
func GetAssetList(c *gin.Context) {
	// 获取用户令牌
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "获取资产列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取个人资产
	asset_service := service.NewAssetService(com.Database)

	private_asset_list, err := asset_service.GetAssetByOwnerID(user_token.UserID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "获取资产列表失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 获取用户共享的资产列表
	shared_asset_list, err := service.NewUserAssetService(com.Database).GetUserAssetMappingsByUserId(user_token.UserID)
	var public_asset_list []service.Asset
	if err == nil {
		public_asset_list = make([]service.Asset, len(shared_asset_list))
		for index, mapping := range shared_asset_list {
			public_asset, err := asset_service.GetAssetByID(mapping.AssetID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "获取资产列表失败",
					"error":   err.Error(),
				})
				return
			}
			public_asset_list[index] = *public_asset
		}
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "获取资产列表失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 构建资产列表
	type _AssetList struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		CreateAt int64  `json:"create_at"`
		UpdateAt int64  `json:"update_at"`
	}

	asset_list := make([]_AssetList, len(private_asset_list)+len(public_asset_list))

	for index, asset := range private_asset_list {
		asset_list[index] = _AssetList{
			Type:     asset.Type,
			Name:     asset.Name,
			CreateAt: asset.CreatedAt.Unix(),
			UpdateAt: asset.UpdatedAt.Unix(),
		}
	}
	for index, asset := range public_asset_list {
		asset_list[index+len(private_asset_list)] = _AssetList{
			Type:     asset.Type,
			Name:     asset.Name,
			CreateAt: asset.CreatedAt.Unix(),
			UpdateAt: asset.UpdatedAt.Unix(),
		}
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message":    "获取资产列表成功",
		"asset_list": asset_list,
	})
}

// 查询指定资产信息（个人资产和有权限的非个人资产）
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	asset_id: uint (资产ID)
func GetAssetInfo(c *gin.Context) {
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "查询资产信息失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理用户请求
	asset_id, err := strconv.ParseUint(c.PostForm("asset_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "查询资产信息失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取资产
	asset_service := service.NewAssetService(com.Database)

	asset, err := asset_service.GetAssetByID(uint(asset_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询资产信息失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查资产是否属于当前用户，或者用户是否有权限使用
	if asset.OwnerID != user_token.UserID {
		shared_asset, err := service.NewUserAssetService(com.Database).GetUserAssetMappingsByAssetId(uint(asset_id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "查询资产信息失败",
				"error":   err.Error(),
			})
			return
		}

		permission := false
		for _, mapping := range shared_asset {
			if mapping.UserID == user_token.UserID {
				// 用户有权限查看资产信息
				permission = true
				break
			}
		}

		if !permission {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "无权限查看资产信息",
			})
			return
		}
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "查询资产信息成功",
		"asset":   asset,
	})
}

// 修改指定资产信息（个人资产和有权限的非个人资产）
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	asset_id: uint (资产ID)
//	reset_type: string (重置类型，name或data)
//	new_name: string (新资产名称，如果reset_type为name时必填)
//	new_data: string (新资产数据，如果reset_type为data时必填)
func PutAssetInfo(c *gin.Context) {
	// 获取用户令牌
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "修改资产信息失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理用户请求
	asset_id, err := strconv.ParseUint(c.PostForm("asset_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改资产信息失败",
			"error":   err.Error(),
		})
		return
	}
	reset_type := c.PostForm("reset_type")
	if reset_type != "name" && reset_type != "data" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改资产信息失败",
			"error":   "无效的重置类型",
		})
		return
	}
	new_name := c.PostForm("new_name")
	if new_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "资产名称不能为空",
		})
	}
	new_data := c.PostForm("new_data")
	if new_data == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "资产数据不能为空",
		})
		return
	}

	// 获取资产
	asset_service := service.NewAssetService(com.Database)

	asset, err := asset_service.GetAssetByID(uint(asset_id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改资产信息失败",
			"error":   err.Error(),
		})
		return
	}

	if asset.OwnerID != user_token.UserID {
		shared_asset, err := service.NewUserAssetService(com.Database).GetUserAssetMappingsByAssetId(uint(asset_id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "修改资产信息失败",
				"error":   err.Error(),
			})
			return
		}

		ok := false
		for _, mapping := range shared_asset {
			if mapping.UserID == user_token.UserID {
				if mapping.Permission != "execute" {
					c.JSON(http.StatusForbidden, gin.H{
						"message": "无权限修改资产信息",
					})
					return
				}
				// 用户有权限修改资产信息
				ok = true
				break
			}
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "无权限修改资产信息",
			})
			return
		}
	}

	switch reset_type {
	case "name":
		// 修改资产名称
		err = asset_service.UpdateAssetNameByID(uint(asset_id), new_name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "修改资产名称失败",
				"error":   err.Error(),
			})
		}
	case "data":
		// 修改资产数据
		err = asset_service.UpdateAssetDataByID(uint(asset_id), datatypes.JSON(new_data))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "修改资产数据失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 创建操作日志
	additional_information := map[string]any{}
	service.RecordUserBasicInformation(c, additional_information)

	// 异步记录操作日志
	go func() {
		operation_log := &service.OperationLog{
			UserID:        user_token.UserID,
			ResourceType:  "asset",
			ResourceID:    asset.ID,
			Action:        "update",
			ActionDetails: "修改资产信息",
			CreatedAt:     time.Now(),
		}

		additional_infomation_json, err := json.Marshal(additional_information)
		if err != nil {
			log.Println("操作日志记录错误:", err)
		} else {
			operation_log.AdditionalInformation = additional_infomation_json
		}

		err = service.NewOperationLogService(com.Database).CreateOperationLog(operation_log)
		if err != nil {
			log.Println("操作日志记录错误:", err)

			json_data, err := json.Marshal(operation_log)
			if err != nil {
				log.Println("操作日志记录错误:", err)
			}

			log.Println("操作日志记录:", string(json_data))
			return
		}
	}()

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "修改资产信息成功",
	})
}

// 分享个人资产权限
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	asset_id: uint (资产ID)
//	email: string (分享用户的邮箱)
//	permission: string (分享权限，use或execute)
func PostShareAsset(c *gin.Context) {
	// 获取用户信息
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "分享个人资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理用户请求
	asset_id, err := strconv.ParseUint(c.PostForm("asset_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "分享个人资产权限失败",
			"error":   err.Error(),
		})
		return
	}
	email := c.PostForm("email")
	permission := c.PostForm("permission")

	// 获取分享用户
	user_service := service.NewUserService(com.Database)

	user, err := user_service.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "分享用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "分享个人资产权限失败",
				"error":   err.Error(),
			})
		}
		return
	}

	// 获取分享资产
	asset_service := service.NewAssetService(com.Database)

	switch permission {
	case "use", "execute":
		// 校验资产是否存在
		asset, err := asset_service.GetAssetByID(uint(asset_id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "资产不存在",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "分享个人资产权限失败",
					"error":   err.Error(),
				})
			}
			return
		}

		// 校验资产是否属于当前用户
		if asset.OwnerID != user_token.UserID {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "没有分享权限",
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "分享权限错误",
		})
		return
	}

	// 创建用户资产映射
	err = service.NewUserAssetService(com.Database).CreateUserAssetMapping(user.ID, uint(asset_id), permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "分享个人资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 创建操作日志
	additional_information := map[string]any{}

	service.RecordUserBasicInformation(c, additional_information)

	// 异步记录操作日志
	go func() {
		operation_log := &service.OperationLog{
			UserID:        user_token.UserID,
			ResourceType:  "asset",
			ResourceID:    uint(asset_id),
			Action:        "create",
			ActionDetails: "分享个人资产权限",
			CreatedAt:     time.Now(),
		}

		additional_infomation_json, err := json.Marshal(additional_information)
		if err != nil {
			log.Println("操作日志记录错误:", err)
		} else {
			operation_log.AdditionalInformation = additional_infomation_json
		}

		err = service.NewOperationLogService(com.Database).CreateOperationLog(operation_log)
		if err != nil {
			log.Println("操作日志记录错误:", err)

			json_data, err := json.Marshal(operation_log)
			if err != nil {
				log.Println("操作日志记录错误:", err)
			}

			log.Println("操作日志记录:", string(json_data))
			return
		}
	}()

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "分享个人资产权限成功",
	})
}

// 查询个人分享的资产权限
//
// 请求头参数：
//
//	Authorization: Bearer <token>
func GetShareAsset(c *gin.Context) {
	// 获取用户信息
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "查询个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取用户分享的资产
	user_asset_mappings, err := service.NewUserAssetService(com.Database).GetUserAssetMappingsByUserId(user_token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	type _Mapping struct {
		ID         uint   `json:"id"`
		AssetID    uint   `json:"asset_id"`
		UserID     uint   `json:"user_id"`
		Permission string `json:"permission"`
		CreateAt   int64  `json:"create_at"`
		UpdateAt   int64  `json:"update_at"`
	}

	mappings := make([]_Mapping, len(user_asset_mappings))
	for index, mapping := range user_asset_mappings {
		mappings[index] = _Mapping{
			ID:         mapping.ID,
			AssetID:    mapping.AssetID,
			UserID:     mapping.UserID,
			Permission: mapping.Permission,
			CreateAt:   mapping.CreatedAt.Unix(),
			UpdateAt:   mapping.UpdatedAt.Unix(),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "查询个人分享的资产权限成功",
		"mappings": mappings,
	})
}

// 修改个人分享的资产权限
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	user_id: uint (用户ID)
//	asset_id: uint (资产ID)
//	permission: string (分享权限，use或execute)
func PutShareAsset(c *gin.Context) {
	// 获取用户信息
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "修改个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理用户请求
	user_id, err := strconv.ParseUint(c.PostForm("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}
	asset_id, err := strconv.ParseUint(c.PostForm("asset_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}
	new_permission := c.PostForm("permission")
	if new_permission != "use" && new_permission != "execute" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的分享权限",
		})
		return
	}

	// 获取资产
	asset, err := service.NewAssetService(com.Database).GetAssetByID(uint(asset_id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "资产不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "修改个人分享的资产权限失败",
				"error":   err.Error(),
			})
		}
		return
	}

	// 检查资产是否属于当前用户
	if asset.OwnerID != user_token.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "无权限修改资产信息",
		})
		return
	}

	// 修改用户资产映射权限
	err = service.NewUserAssetService(com.Database).UpdateUserAssetMappingPermissionByUserIDAndAssetID(uint(user_id), uint(asset_id), new_permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 创建操作日志
	additional_information := map[string]any{}
	service.RecordUserBasicInformation(c, additional_information)

	// 异步记录操作日志
	go func() {
		operation_log := &service.OperationLog{
			UserID:        user_token.UserID,
			ResourceType:  "asset",
			ResourceID:    uint(asset_id),
			Action:        "update",
			ActionDetails: "修改个人分享的资产权限",
			CreatedAt:     time.Now(),
		}

		additional_infomation_json, err := json.Marshal(additional_information)
		if err != nil {
			log.Println("操作日志记录错误:", err)
		} else {
			operation_log.AdditionalInformation = additional_infomation_json
		}

		err = service.NewOperationLogService(com.Database).CreateOperationLog(operation_log)
		if err != nil {
			log.Println("操作日志记录错误:", err)

			json_data, err := json.Marshal(operation_log)
			if err != nil {
				log.Println("操作日志记录错误:", err)
			}

			log.Println("操作日志记录:", string(json_data))
			return
		}
	}()

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "修改个人分享的资产权限成功",
	})
}

// 撤销个人分享的资产权限
//
// 请求头参数：
//
//	Authorization: Bearer <token>
//
// 请求体参数：form
//
//	user_id: uint (用户ID)
//	asset_id: uint (资产ID)
func DeleteShareAsset(c *gin.Context) {
	// 获取用户信息
	user_token, err := service.AnalyzeToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "撤销个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理用户请求
	user_id, err := strconv.ParseUint(c.PostForm("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "撤销个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}
	asset_id, err := strconv.ParseUint(c.PostForm("asset_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "撤销个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取资产
	asset_service := service.NewAssetService(com.Database)

	asset, err := asset_service.GetAssetByID(uint(asset_id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "资产不存在",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "撤销个人分享的资产权限失败",
				"error":   err.Error(),
			})
			return
		}
	}

	// 检查资产是否属于当前用户
	if asset.OwnerID != user_token.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "无权限撤销资产信息",
		})
		return
	}

	// 撤销用户资产映射
	err = service.NewUserAssetService(com.Database).DeleteUserAssetMappingByUserIdAndAssetId(uint(user_id), uint(asset_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "撤销个人分享的资产权限失败",
			"error":   err.Error(),
		})
		return
	}

	// 创建操作日志
	additional_information := map[string]any{}
	service.RecordUserBasicInformation(c, additional_information)

	// 异步记录操作日志
	go func() {
		operation_log := &service.OperationLog{
			UserID:        user_token.UserID,
			ResourceType:  "asset",
			ResourceID:    uint(asset_id),
			Action:        "delete",
			ActionDetails: "撤销个人分享的资产权限",
			CreatedAt:     time.Now(),
		}

		additional_infomation_json, err := json.Marshal(additional_information)
		if err != nil {
			log.Println("操作日志记录错误:", err)
		} else {
			operation_log.AdditionalInformation = additional_infomation_json
		}

		err = service.NewOperationLogService(com.Database).CreateOperationLog(operation_log)
		if err != nil {
			log.Println("操作日志记录错误:", err)

			json_data, err := json.Marshal(operation_log)
			if err != nil {
				log.Println("操作日志记录错误:", err)
			}

			log.Println("操作日志记录:", string(json_data))
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "撤销个人分享的资产权限成功",
	})
}

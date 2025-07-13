package api

import (
	"log"

	"github.com/ZSLTChenXiYin/CRMS/crms-backend/api/controller"
	"github.com/ZSLTChenXiYin/CRMS/crms-backend/com"
)

func InitUser() error {
	user_group := com.Router.Group("/user")
	{
		// 注册用户
		user_group.POST("/register", controller.PostRegister)

		// 登录用户
		user_group.POST("/login", controller.PostLogin)
		// 登出用户
		user_group.POST("/logout", controller.PostLogout)

		// 查询用户个人信息
		user_group.GET("/info", controller.GetUserInfo)
		// 修改用户个人信息
		user_group.PUT("/update", controller.PutUserInfo)

		// 修改用户个人密码
		user_group.PUT("/reset-password", controller.PutResetPassword)

		// 删除用户
		user_group.DELETE("/delete", controller.DeleteUser)
	}

	log.Println("[INFO] User API初始化成功")

	return nil
}
